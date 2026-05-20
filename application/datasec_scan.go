package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
)

const (
	dataSecScanKindDB        = "db"
	dataSecScanKindSensitive = "sensitive"
	dataSecMaxConcurrentDB   = 5
)

// DataSecScan 数据安全（数据库基线 / 敏感数据发现）
type DataSecScan struct{}

type datasecExtendMeta struct {
	DataSecScanKind string `json:"dataSecScanKind"`
	ErrorMessage    string `json:"errorMessage,omitempty"`
}

type datasecTaskConfig struct {
	DBType     int    `json:"dbType"`
	DBHost     string `json:"dbHost"`
	DBPort     int    `json:"dbPort"`
	DBName     string `json:"dbName"`
	DBUser     string `json:"dbUser"`
	DBPassword string `json:"dbPassword"`
	DataTypes  []int  `json:"dataTypes,omitempty"`
	ScanAllDB  bool   `json:"scanAllDb,omitempty"`
}

func (a *DataSecScan) RunDBCheck(ctx context.Context, uid int, req *typespec.DataSecDBRunReq) (*typespec.DataSecScanRunResp, error) {
	return a.runDataSecTask(ctx, uid, enums.TaskTypeDataSecDB, dataSecScanKindDB, req.Name, datasecTaskConfig{
		DBType: req.DBType, DBHost: strings.TrimSpace(req.DBHost), DBPort: req.DBPort,
		DBName: strings.TrimSpace(req.DBName), DBUser: strings.TrimSpace(req.DBUser), DBPassword: req.DBPassword,
	})
}

func (a *DataSecScan) RunSensitiveScan(ctx context.Context, uid int, req *typespec.DataSecSensitiveRunReq) (*typespec.DataSecScanRunResp, error) {
	return a.runDataSecTask(ctx, uid, enums.TaskTypeDataSecSensitive, dataSecScanKindSensitive, req.Name, datasecTaskConfig{
		DBType: req.DBType, DBHost: strings.TrimSpace(req.DBHost), DBPort: req.DBPort,
		DBName: strings.TrimSpace(req.DBName), DBUser: strings.TrimSpace(req.DBUser), DBPassword: req.DBPassword,
		DataTypes: req.DataTypes, ScanAllDB: req.ScanAllDB,
	})
}

func (a *DataSecScan) runDataSecTask(ctx context.Context, uid, taskType int, scanKind, taskName string, cfg datasecTaskConfig) (*typespec.DataSecScanRunResp, error) {
	if cfg.DBHost == "" {
		return nil, fmt.Errorf("数据库地址不能为空")
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("用户名不能为空")
	}
	if cfg.DBType < 1 {
		return nil, fmt.Errorf("请选择数据库类型")
	}
	if cfg.DBPort <= 0 {
		cfg.DBPort = defaultDataSecDBPort(cfg.DBType)
	}

	targetURL := formatDataSecTargetURL(cfg.DBHost, cfg.DBPort, cfg.DBName)
	taskID, targets, err := a.createDataSecTask(ctx, uid, taskType, scanKind, taskName, []string{targetURL}, cfg)
	if err != nil {
		return nil, err
	}

	scanCtx := context.Background()
	sem := make(chan struct{}, dataSecMaxConcurrentDB)
	for _, row := range targets {
		t := row
		go func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			a.startDataSecTargetScan(scanCtx, t, scanKind)
		}()
	}
	return &typespec.DataSecScanRunResp{ID: strconv.Itoa(taskID)}, nil
}

func (a *DataSecScan) createDataSecTask(ctx context.Context, uid, taskType int, scanKind, taskName string, targetList []string, cfg datasecTaskConfig) (int, []mysqls.TaskTarget, error) {
	var taskSrv services.TaskTask
	taskID, err := taskSrv.Save(ctx, uid, 0, firstNonEmpty(taskName, "未命名任务"), targetList, enums.TaskExecTypeImmediate, "{}", enums.ConfigJson{}, time.Unix(0, 0), 0, 0)
	if err != nil {
		return 0, nil, err
	}

	ext := datasecExtendMeta{DataSecScanKind: scanKind}
	extBytes, _ := json.Marshal(ext)
	cfgJSON := mustJSON(cfg)

	var taskModel mysqls.TaskTask
	_ = taskModel.UpdateTaskTaskByIds(ctx, []int{taskID}, map[string]interface{}{
		"task_type":   taskType,
		"target_num":  len(targetList),
		"status":      enums.TaskStatusRunning,
		"update_time": time.Now(),
	})

	var taskInfo mysqls.TaskTaskInfo
	_ = taskInfo.UpdateTaskInfoByTaskIds(ctx, []int{taskID}, map[string]interface{}{
		"status":      enums.TaskStatusRunning,
		"update_time": time.Now(),
	})

	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, taskID)
	if len(targets) == 0 {
		return 0, nil, fmt.Errorf("创建扫描目标失败")
	}
	for i := range targets {
		targets[i].ExtendField = string(extBytes)
		targets[i].TaskTemplateJSON = cfgJSON
		_ = targets[i].UpdateTaskTarget(ctx)
	}
	return taskID, targets, nil
}

func (a *DataSecScan) startDataSecTargetScan(ctx context.Context, target mysqls.TaskTarget, scanKind string) {
	cfg := parseDatasecConfig(target.TaskTemplateJSON)
	var targetSrv services.TaskTarget
	_, err := targetSrv.UpdateTargetAndSaveTargetLog(ctx, target.TaskID, target.ID, enums.TargetStatusRunning, enums.TargetIsAliveY, target.TargetURL)
	if err != nil {
		log.Errorf("datasec start target log: %v", err)
		a.finishDataSecTarget(ctx, target, scanKind, enums.TaskRiskSafe, err.Error())
		return
	}

	var scanErr error
	var targetRisk int
	switch scanKind {
	case dataSecScanKindDB:
		targetRisk, scanErr = a.runDBBaselineForTarget(ctx, target, cfg)
	case dataSecScanKindSensitive:
		targetRisk, scanErr = a.runSensitiveForTarget(ctx, target, cfg)
	default:
		scanErr = fmt.Errorf("未知扫描类型")
	}

	errMsg := ""
	if scanErr != nil {
		errMsg = scanErr.Error()
		log.Errorf("datasec scan task=%d target=%d: %v", target.TaskID, target.ID, scanErr)
	}
	a.finishDataSecTarget(ctx, target, scanKind, targetRisk, errMsg)
}

func (a *DataSecScan) runDBBaselineForTarget(ctx context.Context, target mysqls.TaskTarget, cfg datasecTaskConfig) (int, error) {
	checker := services.GetDBBaselineChecker()
	task := &services.DBCheckTask{
		TaskID:   target.TaskID,
		TargetID: target.ID,
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		DBType:   cfg.DBType,
		Username: cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	}
	report, err := checker.RunDBCheck(ctx, task)
	if err != nil {
		return enums.TaskRiskSafe, err
	}
	return worstDBCheckTaskRisk(report.Results), nil
}

func (a *DataSecScan) runSensitiveForTarget(ctx context.Context, target mysqls.TaskTarget, cfg datasecTaskConfig) (int, error) {
	finder := services.GetSensitiveDataFinder()
	task := &services.SensitiveDataTask{
		TaskID:    target.TaskID,
		TargetID:  target.ID,
		Host:      cfg.DBHost,
		Port:      cfg.DBPort,
		DBType:    cfg.DBType,
		Username:  cfg.DBUser,
		Password:  cfg.DBPassword,
		DBName:    cfg.DBName,
		ScanAllDB: cfg.ScanAllDB,
	}
	report, err := finder.RunScan(ctx, task)
	if err != nil {
		return enums.TaskRiskSafe, err
	}
	if report.HighCount > 0 {
		return enums.TaskRiskHigh, nil
	}
	if report.MiddleCount > 0 {
		return enums.TaskRiskMid, nil
	}
	if report.LowCount > 0 {
		return enums.TaskRiskLow, nil
	}
	return enums.TaskRiskSafe, nil
}

func (a *DataSecScan) finishDataSecTarget(ctx context.Context, target mysqls.TaskTarget, scanKind string, targetRisk int, errMsg string) {
	ext := parseDatasecExtend(target.ExtendField)
	ext.DataSecScanKind = scanKind
	ext.ErrorMessage = errMsg
	extBytes, _ := json.Marshal(ext)

	params := map[string]interface{}{
		"status":      enums.TargetStatusFinish,
		"risk_level":  mapTaskRiskToTargetRisk(targetRisk),
		"extend_field": string(extBytes),
		"update_time": time.Now(),
	}
	var t mysqls.TaskTarget
	_ = t.UpdateTargetById(ctx, target.ID, params)
	a.maybeFinishDataSecTask(ctx, target.TaskID)
}

func (a *DataSecScan) maybeFinishDataSecTask(ctx context.Context, taskID int) {
	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, taskID)
	if len(targets) == 0 {
		return
	}
	allDone := true
	worst := enums.TaskRiskSafe
	for _, t := range targets {
		if t.Status != enums.TargetStatusFinish {
			allDone = false
			break
		}
		if t.RiskLevel < worst {
			worst = mapTargetRiskToTaskRisk(t.RiskLevel)
		}
	}
	if !allDone {
		return
	}
	var taskModel mysqls.TaskTask
	_ = taskModel.UpdateTaskTaskByIds(ctx, []int{taskID}, map[string]interface{}{
		"status":      enums.TaskStatusFinish,
		"risk_level":  worst,
		"update_time": time.Now(),
	})
	var taskInfo mysqls.TaskTaskInfo
	_ = taskInfo.UpdateTaskInfoByTaskIds(ctx, []int{taskID}, map[string]interface{}{
		"status":      enums.TaskStatusFinish,
		"update_time": time.Now(),
	})
}

func (a *DataSecScan) ListDBChecks(ctx context.Context, uid int, req *typespec.DataSecScanListReq) (*typespec.DataSecDBListResp, error) {
	return a.listDBFromDB(ctx, uid, req)
}

func (a *DataSecScan) GetDBCheckDetail(ctx context.Context, uid int, id string) (*typespec.DataSecDBListItem, error) {
	taskID, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil || taskID <= 0 {
		return nil, fmt.Errorf("任务不存在")
	}
	var taskModel mysqls.TaskTask
	task, err := taskModel.GetTaskCheckTask(ctx, taskID)
	if err != nil || task.ID == 0 || task.TaskType != enums.TaskTypeDataSecDB {
		return nil, fmt.Errorf("任务不存在")
	}
	if uid > 0 && task.UserID != uid {
		return nil, fmt.Errorf("无权查看该任务")
	}
	item, err := a.buildDBListItem(ctx, task, true)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (a *DataSecScan) listDBFromDB(ctx context.Context, uid int, req *typespec.DataSecScanListReq) (*typespec.DataSecDBListResp, error) {
	page, size := normalizeDataSecPage(req.Page, req.Size)
	var taskModel mysqls.TaskTask
	tasks, total, err := taskModel.GetAppSecTaskList(ctx, enums.TaskTypeDataSecDB, uid, page, size, strings.TrimSpace(req.Search))
	if err != nil {
		return nil, err
	}
	list := make([]typespec.DataSecDBListItem, 0, len(tasks))
	for _, t := range tasks {
		item, err := a.buildDBListItem(ctx, t, false)
		if err != nil {
			continue
		}
		list = append(list, item)
	}
	return &typespec.DataSecDBListResp{List: list, Total: int(total)}, nil
}

func (a *DataSecScan) buildDBListItem(ctx context.Context, task mysqls.TaskTask, withDetails bool) (typespec.DataSecDBListItem, error) {
	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, task.ID)
	cfg := datasecTaskConfig{}
	if len(targets) > 0 {
		cfg = parseDatasecConfig(targets[0].TaskTemplateJSON)
	}

	critical, high, mid, low := countDBCheckRisksByTask(ctx, task.ID)
	item := typespec.DataSecDBListItem{
		ID:              strconv.Itoa(task.ID),
		Name:            task.TaskName,
		DBType:          cfg.DBType,
		DBHost:          cfg.DBHost,
		DBPort:          cfg.DBPort,
		DBName:          cfg.DBName,
		RiskLevel:       aggregateFrontendRisk(critical, high, mid, low, task.RiskLevel),
		Status:          resolveAppSecTaskStatus(task, targets),
		CreateTime:      formatAppSecTime(task.CreateTime),
		CheckTime:       formatAppSecTime(task.UpdateTime),
		CriticalCount:   critical,
		HighRiskCount:   high,
		MiddleRiskCount: mid,
		LowRiskCount:    low,
	}
	if withDetails {
		item.Items = loadDBCheckDetailItems(ctx, task.ID)
	}
	return item, nil
}

func (a *DataSecScan) ListSensitiveScans(ctx context.Context, uid int, req *typespec.DataSecScanListReq) (*typespec.DataSecSensitiveListResp, error) {
	page, size := normalizeDataSecPage(req.Page, req.Size)
	var taskModel mysqls.TaskTask
	tasks, total, err := taskModel.GetAppSecTaskList(ctx, enums.TaskTypeDataSecSensitive, uid, page, size, strings.TrimSpace(req.Search))
	if err != nil {
		return nil, err
	}
	list := make([]typespec.DataSecSensitiveListItem, 0, len(tasks))
	for _, t := range tasks {
		item, err := a.buildSensitiveListItem(ctx, t, false)
		if err != nil {
			continue
		}
		list = append(list, item)
	}
	return &typespec.DataSecSensitiveListResp{List: list, Total: int(total)}, nil
}

func (a *DataSecScan) GetSensitiveScanDetail(ctx context.Context, uid int, id string) (*typespec.DataSecSensitiveListItem, error) {
	taskID, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil || taskID <= 0 {
		return nil, fmt.Errorf("任务不存在")
	}
	var taskModel mysqls.TaskTask
	task, err := taskModel.GetTaskCheckTask(ctx, taskID)
	if err != nil || task.ID == 0 || task.TaskType != enums.TaskTypeDataSecSensitive {
		return nil, fmt.Errorf("任务不存在")
	}
	if uid > 0 && task.UserID != uid {
		return nil, fmt.Errorf("无权查看该任务")
	}
	item, err := a.buildSensitiveListItem(ctx, task, true)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (a *DataSecScan) buildSensitiveListItem(ctx context.Context, task mysqls.TaskTask, withDetails bool) (typespec.DataSecSensitiveListItem, error) {
	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, task.ID)
	cfg := datasecTaskConfig{}
	if len(targets) > 0 {
		cfg = parseDatasecConfig(targets[0].TaskTemplateJSON)
	}
	high, mid, low, total := countSensitiveByTask(ctx, task.ID)
	item := typespec.DataSecSensitiveListItem{
		ID:          strconv.Itoa(task.ID),
		Name:        task.TaskName,
		DBType:      cfg.DBType,
		DBHost:      cfg.DBHost,
		DBPort:      cfg.DBPort,
		DBName:      cfg.DBName,
		TotalCount:  total,
		HighCount:   high,
		MediumCount: mid,
		LowCount:    low,
		Status:      resolveAppSecTaskStatus(task, targets),
		CreateTime:  formatAppSecTime(task.CreateTime),
		ScanTime:    formatAppSecTime(task.UpdateTime),
	}
	if withDetails {
		results, _ := (&mysqls.SensitiveDataResult{}).GetByTaskID(ctx, task.ID)
		item.TypeStats = buildSensitiveTypeStats(results)
		item.Items = buildSensitiveDetailItems(results)
	}
	return item, nil
}

func loadDBCheckDetailItems(ctx context.Context, taskID int) []typespec.DataSecDBDetailItem {
	list, err := (&mysqls.DBCheckResult{}).GetByTaskID(ctx, taskID)
	if err != nil {
		return nil
	}
	out := make([]typespec.DataSecDBDetailItem, 0, len(list))
	for _, r := range list {
		risk := r.RiskLevel
		if r.CheckResult == enums.BaselineCheckResultPass {
			risk = enums.BaselineRiskInfo
		}
		out = append(out, typespec.DataSecDBDetailItem{
			Category:    r.CheckCategory,
			RiskLevel:   risk,
			Result:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			Description: firstNonEmpty(r.RiskDescription, r.RuleName),
			Suggestion:  r.FixSuggestion,
		})
	}
	return out
}

func countDBCheckRisksByTask(ctx context.Context, taskID int) (critical, high, mid, low int) {
	list, err := (&mysqls.DBCheckResult{}).GetByTaskID(ctx, taskID)
	if err != nil {
		return
	}
	for _, r := range list {
		if r.CheckResult != enums.BaselineCheckResultFail {
			continue
		}
		switch r.RiskLevel {
		case enums.BaselineRiskCritical:
			critical++
		case enums.BaselineRiskHigh:
			high++
		case enums.BaselineRiskMiddle:
			mid++
		case enums.BaselineRiskLow:
			low++
		}
	}
	return
}

func worstDBCheckTaskRisk(results []mysqls.DBCheckResult) int {
	worst := enums.TaskRiskSafe
	for _, r := range results {
		if r.CheckResult != enums.BaselineCheckResultFail {
			continue
		}
		tr := baselineRiskToTaskRisk(r.RiskLevel)
		if tr < worst {
			worst = tr
		}
	}
	return worst
}

func baselineRiskToTaskRisk(risk int) int {
	switch risk {
	case enums.BaselineRiskCritical, enums.BaselineRiskHigh:
		return enums.TaskRiskHigh
	case enums.BaselineRiskMiddle:
		return enums.TaskRiskMid
	case enums.BaselineRiskLow:
		return enums.TaskRiskLow
	default:
		return enums.TaskRiskSafe
	}
}

func mapTaskRiskToTargetRisk(taskRisk int) int {
	switch taskRisk {
	case enums.TaskRiskHigh:
		return enums.TargetRiskHigh
	case enums.TaskRiskMid:
		return enums.TargetRiskMid
	case enums.TaskRiskLow:
		return enums.TargetRiskLow
	default:
		return enums.TargetRiskLowNoFound
	}
}

func mapTargetRiskToTaskRisk(targetRisk int) int {
	switch targetRisk {
	case enums.TargetRiskHigh:
		return enums.TaskRiskHigh
	case enums.TargetRiskMid:
		return enums.TaskRiskMid
	case enums.TargetRiskLow:
		return enums.TaskRiskLow
	default:
		return enums.TaskRiskSafe
	}
}

func countSensitiveByTask(ctx context.Context, taskID int) (high, mid, low, total int) {
	list, err := (&mysqls.SensitiveDataResult{}).GetByTaskID(ctx, taskID)
	if err != nil {
		return
	}
	total = len(list)
	for _, r := range list {
		switch r.DataLevel {
		case enums.SensitiveDataLevelHigh:
			high++
		case enums.SensitiveDataLevelMiddle:
			mid++
		default:
			low++
		}
	}
	return
}

func buildSensitiveTypeStats(results []mysqls.SensitiveDataResult) []typespec.DataSecSensitiveTypeStatItem {
	counts := make(map[int]int)
	for _, r := range results {
		counts[r.DataType]++
	}
	out := make([]typespec.DataSecSensitiveTypeStatItem, 0, len(counts))
	for dt, c := range counts {
		out = append(out, typespec.DataSecSensitiveTypeStatItem{DataType: dt, Count: c})
	}
	return out
}

func buildSensitiveDetailItems(results []mysqls.SensitiveDataResult) []typespec.DataSecSensitiveDetailItem {
	out := make([]typespec.DataSecSensitiveDetailItem, 0, len(results))
	for _, r := range results {
		cnt := 1
		if r.TotalRows > 0 {
			cnt = int(r.TotalRows)
		}
		out = append(out, typespec.DataSecSensitiveDetailItem{
			TableName:        r.TableNameStr,
			ColumnName:       r.ColumnName,
			DataType:         r.DataType,
			SensitivityLevel: r.DataLevel,
			SampleData:       r.SampleData,
			Count:            cnt,
		})
	}
	return out
}

func parseDatasecConfig(raw string) datasecTaskConfig {
	var cfg datasecTaskConfig
	if strings.TrimSpace(raw) == "" {
		return cfg
	}
	_ = json.Unmarshal([]byte(raw), &cfg)
	return cfg
}

func parseDatasecExtend(raw string) datasecExtendMeta {
	var meta datasecExtendMeta
	if raw == "" {
		return meta
	}
	_ = json.Unmarshal([]byte(raw), &meta)
	return meta
}

func formatDataSecTargetURL(host string, port int, dbName string) string {
	if dbName != "" {
		return fmt.Sprintf("%s:%d/%s", host, port, dbName)
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func defaultDataSecDBPort(dbType int) int {
	switch dbType {
	case enums.DBSupportTypePostgreSQL:
		return 5432
	case enums.DBSupportTypeMongoDB:
		return 27017
	case enums.DBSupportTypeRedis:
		return 6379
	case enums.DBSupportTypeCouchDB:
		return 5984
	default:
		return 3306
	}
}

func normalizeDataSecPage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	return page, size
}
