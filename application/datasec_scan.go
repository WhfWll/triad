package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
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
)

// DataSecScan 数据安全（数据库基线 / 敏感数据发现）
type DataSecScan struct{}

var (
	dataSecTaskCancelMu sync.Mutex
	dataSecTaskCancels  = make(map[int]context.CancelFunc)
)

func setDataSecTaskCancel(taskID int, cancel context.CancelFunc) {
	dataSecTaskCancelMu.Lock()
	dataSecTaskCancels[taskID] = cancel
	dataSecTaskCancelMu.Unlock()
}

func takeDataSecTaskCancel(taskID int) context.CancelFunc {
	dataSecTaskCancelMu.Lock()
	cancel := dataSecTaskCancels[taskID]
	delete(dataSecTaskCancels, taskID)
	dataSecTaskCancelMu.Unlock()
	return cancel
}

type datasecExtendMeta struct {
	DataSecScanKind string `json:"dataSecScanKind"`
	ErrorMessage    string `json:"errorMessage,omitempty"`
}

type datasecTaskConfig struct {
	DBType        int    `json:"dbType"`
	DBHost        string `json:"dbHost"`
	DBPort        int    `json:"dbPort"`
	DBName        string `json:"dbName"`
	DBUser        string `json:"dbUser"`
	DBPassword    string `json:"dbPassword"`
	ScanSensitive bool   `json:"scanSensitive,omitempty"`
	DataTypes     []int  `json:"dataTypes,omitempty"`
	ScanAllDB     bool   `json:"scanAllDb,omitempty"`
}

func (a *DataSecScan) RunDBCheck(ctx context.Context, uid int, req *typespec.DataSecDBRunReq) (*typespec.DataSecScanRunResp, error) {
	configs, err := a.buildDBRunConfigs(ctx, uid, req)
	if err != nil {
		return nil, err
	}
	return a.runDataSecTask(ctx, uid, enums.TaskTypeDataSecDB, dataSecScanKindDB, req.Name, configs)
}

func (a *DataSecScan) RunSensitiveScan(ctx context.Context, uid int, req *typespec.DataSecSensitiveRunReq) (*typespec.DataSecScanRunResp, error) {
	configs, err := a.buildSensitiveRunConfigs(ctx, uid, req)
	if err != nil {
		return nil, err
	}
	return a.runDataSecTask(ctx, uid, enums.TaskTypeDataSecSensitive, dataSecScanKindSensitive, req.Name, configs)
}

func (a *DataSecScan) buildDBRunConfigs(ctx context.Context, uid int, req *typespec.DataSecDBRunReq) ([]datasecTaskConfig, error) {
	configs, err := a.mergeRunTargetConfigs(ctx, uid, req.DBType.Int(), req.Targets, req.LibraryTargetIDs,
		legacyDataSecTargetInputs(req.DBHost, req.DBPort, req.DBName, req.DBUser, req.DBPassword), req.DataTypes, req.ScanAllDB)
	if err != nil {
		return nil, err
	}
	if req.ScanSensitive {
		for i := range configs {
			configs[i].ScanSensitive = true
			if len(req.DataTypes) > 0 {
				configs[i].DataTypes = req.DataTypes
			}
			configs[i].ScanAllDB = req.ScanAllDB
		}
	}
	return configs, nil
}

func (a *DataSecScan) buildSensitiveRunConfigs(ctx context.Context, uid int, req *typespec.DataSecSensitiveRunReq) ([]datasecTaskConfig, error) {
	return a.mergeRunTargetConfigs(ctx, uid, req.DBType.Int(), req.Targets, req.LibraryTargetIDs,
		legacyDataSecTargetInputs(req.DBHost, req.DBPort, req.DBName, req.DBUser, req.DBPassword), req.DataTypes, req.ScanAllDB)
}

func (a *DataSecScan) mergeRunTargetConfigs(ctx context.Context, uid, dbType int, targets []typespec.DataSecDBTargetInput, libraryIDs []int, legacy []typespec.DataSecDBTargetInput, dataTypes []int, scanAllDB bool) ([]datasecTaskConfig, error) {
	items := targets
	if len(items) == 0 {
		items = legacy
	}
	configs, err := normalizeDataSecTargetConfigs(dbType, items, nil, dataTypes, scanAllDB)
	if err != nil && len(libraryIDs) == 0 {
		return nil, err
	}
	if configs == nil {
		configs = []datasecTaskConfig{}
	}
	if len(libraryIDs) > 0 {
		libConfigs, libErr := NewDatasecDBTargetApp().ResolveTargets(ctx, uid, dbType, libraryIDs)
		if libErr != nil {
			return nil, libErr
		}
		configs = append(configs, libConfigs...)
	}
	if len(configs) == 0 {
		return nil, fmt.Errorf("请至少添加一个数据库目标")
	}
	// re-dedupe merged configs
	seen := make(map[string]struct{})
	out := make([]datasecTaskConfig, 0, len(configs))
	for _, cfg := range configs {
		key := fmt.Sprintf("%s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		cfg.DBType = dbType
		if dbType > 0 {
			cfg.DBType = dbType
		}
		cfg.DataTypes = dataTypes
		cfg.ScanAllDB = scanAllDB
		out = append(out, cfg)
	}
	return out, nil
}

func legacyDataSecTargetInputs(host string, port typespec.FlexInt, dbName, user, password string) []typespec.DataSecDBTargetInput {
	if strings.TrimSpace(host) == "" {
		return nil
	}
	return []typespec.DataSecDBTargetInput{{
		DBHost: host, DBPort: port, DBName: dbName, DBUser: user, DBPassword: password,
	}}
}

func normalizeDataSecTargetConfigs(dbType int, targets, legacy []typespec.DataSecDBTargetInput, dataTypes []int, scanAllDB bool) ([]datasecTaskConfig, error) {
	items := targets
	if len(items) == 0 {
		items = legacy
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("请至少添加一个数据库目标")
	}
	if dbType < 1 {
		return nil, fmt.Errorf("请选择数据库类型")
	}
	out := make([]datasecTaskConfig, 0, len(items))
	seen := make(map[string]struct{})
	for i, item := range items {
		host := strings.TrimSpace(item.DBHost)
		user := strings.TrimSpace(item.DBUser)
		if host == "" {
			return nil, fmt.Errorf("第 %d 个目标：数据库地址不能为空", i+1)
		}
		if user == "" && dbType != enums.DBSupportTypeRedis {
			return nil, fmt.Errorf("第 %d 个目标：用户名不能为空", i+1)
		}
		port := item.DBPort.Int()
		if port <= 0 {
			port = defaultDataSecDBPort(dbType)
		}
		key := fmt.Sprintf("%s:%d/%s", host, port, strings.TrimSpace(item.DBName))
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, datasecTaskConfig{
			DBType:     dbType,
			DBHost:     host,
			DBPort:     port,
			DBName:     strings.TrimSpace(item.DBName),
			DBUser:     user,
			DBPassword: item.DBPassword,
			DataTypes:  dataTypes,
			ScanAllDB:  scanAllDB,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("请至少添加一个有效数据库目标")
	}
	return out, nil
}

func dataSecConnFromDBRunReq(req *typespec.DataSecDBRunReq) datasecTaskConfig {
	items := req.Targets
	if len(items) == 0 && strings.TrimSpace(req.DBHost) != "" {
		items = legacyDataSecTargetInputs(req.DBHost, req.DBPort, req.DBName, req.DBUser, req.DBPassword)
	}
	if len(items) > 0 {
		configs, _ := normalizeDataSecTargetConfigs(req.DBType.Int(), items, nil, nil, false)
		if len(configs) > 0 {
			return configs[0]
		}
	}
	return datasecTaskConfig{}
}

func dataSecConnFromTestReq(req *typespec.DataSecDBTestConnReq) datasecTaskConfig {
	return datasecTaskConfig{
		DBType: req.DBType.Int(), DBHost: strings.TrimSpace(req.DBHost), DBPort: req.DBPort.Int(),
		DBName: strings.TrimSpace(req.DBName), DBUser: strings.TrimSpace(req.DBUser), DBPassword: req.DBPassword,
	}
}

// TestDBConnection 测试数据库连接是否可用
func (a *DataSecScan) TestDBConnection(ctx context.Context, req *typespec.DataSecDBTestConnReq) (*typespec.DataSecDBTestConnResp, error) {
	cfg := dataSecConnFromTestReq(req)
	if cfg.DBHost == "" {
		return &typespec.DataSecDBTestConnResp{OK: false, Message: "请填写数据库地址"}, nil
	}
	if cfg.DBUser == "" && cfg.DBType != enums.DBSupportTypeRedis {
		return &typespec.DataSecDBTestConnResp{OK: false, Message: "请填写用户名"}, nil
	}
	if cfg.DBType < 1 {
		return &typespec.DataSecDBTestConnResp{OK: false, Message: "请选择数据库类型"}, nil
	}
	if cfg.DBPort <= 0 {
		cfg.DBPort = defaultDataSecDBPort(cfg.DBType)
	}
	config := &services.DBConnConfig{
		DBType: cfg.DBType, Host: cfg.DBHost, Port: cfg.DBPort,
		Username: cfg.DBUser, Password: cfg.DBPassword, DBName: cfg.DBName,
		Timeout: 10 * time.Second,
	}
	_, err := services.GetDBConnManager().GetConnection(ctx, config)
	if err != nil {
		return &typespec.DataSecDBTestConnResp{OK: false, Message: err.Error()}, nil
	}
	dbName := enums.BaselineEnum.GetDBTypeName(cfg.DBType)
	return &typespec.DataSecDBTestConnResp{
		OK:      true,
		Message: fmt.Sprintf("连接成功：%s %s:%d/%s", dbName, cfg.DBHost, cfg.DBPort, cfg.DBName),
	}, nil
}

func (a *DataSecScan) runDataSecTask(ctx context.Context, uid, taskType int, scanKind, taskName string, configs []datasecTaskConfig) (*typespec.DataSecScanRunResp, error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("请至少添加一个数据库目标")
	}

	taskID, targets, err := a.createDataSecTask(ctx, uid, taskType, scanKind, taskName, configs)
	if err != nil {
		return nil, err
	}

	scanCtx, cancel := context.WithCancel(context.Background())
	setDataSecTaskCancel(taskID, cancel)
	sem := make(chan struct{}, services.GetDataScanConcurrent(ctx))
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

func (a *DataSecScan) createDataSecTask(ctx context.Context, uid, taskType int, scanKind, taskName string, configs []datasecTaskConfig) (int, []mysqls.TaskTarget, error) {
	targetList := make([]string, 0, len(configs))
	for _, cfg := range configs {
		targetList = append(targetList, formatDataSecTargetURL(cfg.DBHost, cfg.DBPort, cfg.DBName))
	}

	var taskSrv services.TaskTask
	taskID, err := taskSrv.Save(ctx, uid, 0, firstNonEmpty(taskName, "未命名任务"), targetList, enums.TaskExecTypeImmediate, "{}", enums.ConfigJson{}, time.Unix(0, 0), 0, 0)
	if err != nil {
		return 0, nil, err
	}

	ext := datasecExtendMeta{DataSecScanKind: scanKind}
	extBytes, _ := json.Marshal(ext)

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
	cfgByURL := make(map[string]datasecTaskConfig, len(configs))
	for _, cfg := range configs {
		cfgByURL[formatDataSecTargetURL(cfg.DBHost, cfg.DBPort, cfg.DBName)] = cfg
	}
	for i := range targets {
		cfg := configs[0]
		if c, ok := cfgByURL[targets[i].TargetURL]; ok {
			cfg = c
		} else if i < len(configs) {
			cfg = configs[i]
		}
		targets[i].ExtendField = string(extBytes)
		targets[i].TaskTemplateJSON = mustJSON(cfg)
		_ = targets[i].UpdateTaskTarget(ctx)
	}
	return taskID, targets, nil
}

func (a *DataSecScan) startDataSecTargetScan(ctx context.Context, target mysqls.TaskTarget, scanKind string) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	cfg := parseDatasecConfig(target.TaskTemplateJSON)
	var targetSrv services.TaskTarget
	logId, err := targetSrv.UpdateTargetAndSaveTargetLog(ctx, target.TaskID, target.ID, enums.TargetStatusRunning, enums.TargetIsAliveY, target.TargetURL)
	if err != nil {
		log.Errorf("datasec start target log: %v", err)
		a.finishDataSecTarget(ctx, target, scanKind, enums.TaskRiskSafe, err.Error(), 0)
		return
	}

	var scanErr error
	var targetRisk int
	switch scanKind {
	case dataSecScanKindDB:
		targetRisk, scanErr = a.runDBBaselineForTarget(ctx, target, cfg, logId)
		if cfg.ScanSensitive {
			sensRisk, sensErr := a.runSensitiveForTarget(ctx, target, cfg, logId)
			targetRisk = worstCombinedTaskRisk(targetRisk, sensRisk)
			if sensErr != nil {
				if scanErr != nil {
					scanErr = fmt.Errorf("基线: %v; 敏感数据: %v", scanErr, sensErr)
				} else {
					scanErr = sensErr
				}
			}
		}
	case dataSecScanKindSensitive:
		targetRisk, scanErr = a.runSensitiveForTarget(ctx, target, cfg, logId)
	default:
		scanErr = fmt.Errorf("未知扫描类型")
	}

	errMsg := ""
	if scanErr != nil {
		errMsg = scanErr.Error()
		log.Errorf("datasec scan task=%d target=%d: %v", target.TaskID, target.ID, scanErr)
		a.appendDataSecScanLog(ctx, target, logId, "datasec", "扫描失败: "+errMsg)
	}
	a.finishDataSecTarget(ctx, target, scanKind, targetRisk, errMsg, logId)
}

func (a *DataSecScan) appendDataSecScanLog(ctx context.Context, target mysqls.TaskTarget, logId int, poc, msg string) {
	if logId <= 0 {
		return
	}
	cfg := parseDatasecConfig(target.TaskTemplateJSON)
	host := cfg.DBHost
	if host == "" {
		host = target.TargetURL
	}
	var logInfo services.TaskLogInfo
	_ = logInfo.AddTaskLogInfo(ctx, target.TaskID, target.ID, logId, host, poc, msg)
}

func (a *DataSecScan) runDBBaselineForTarget(ctx context.Context, target mysqls.TaskTarget, cfg datasecTaskConfig, logId int) (int, error) {
	checker := services.GetDBBaselineChecker()
	task := &services.DBCheckTask{
		TaskID:   target.TaskID,
		TargetID: target.ID,
		LogID:    logId,
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

func (a *DataSecScan) runSensitiveForTarget(ctx context.Context, target mysqls.TaskTarget, cfg datasecTaskConfig, logId int) (int, error) {
	a.appendDataSecScanLog(ctx, target, logId, "datasec", fmt.Sprintf("开始敏感数据扫描：%s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName))
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
		a.appendDataSecScanLog(ctx, target, logId, "datasec", "敏感数据扫描失败: "+err.Error())
		return enums.TaskRiskSafe, err
	}
	a.appendDataSecScanLog(ctx, target, logId, "datasec", fmt.Sprintf("敏感数据扫描完成，发现 %d 条", len(report.Results)))
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

func (a *DataSecScan) finishDataSecTarget(ctx context.Context, target mysqls.TaskTarget, scanKind string, targetRisk int, errMsg string, logId int) {
	ext := parseDatasecExtend(target.ExtendField)
	ext.DataSecScanKind = scanKind
	ext.ErrorMessage = errMsg
	extBytes, _ := json.Marshal(ext)

	params := map[string]interface{}{
		"status":       enums.TargetStatusFinish,
		"risk_level":   mapTaskRiskToTargetRisk(targetRisk),
		"extend_field": string(extBytes),
		"update_time":  time.Now(),
	}
	var t mysqls.TaskTarget
	_ = t.UpdateTargetById(ctx, target.ID, params)
	if logId > 0 {
		var taskLogSrv services.TaskLog
		_ = taskLogSrv.FinishTaskLog(ctx, logId)
	}
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
	takeDataSecTaskCancel(taskID)
}

// CloneTaskTargets 从历史任务复制扫描目标（含凭据，仅任务所属用户可用）
func (a *DataSecScan) CloneTaskTargets(ctx context.Context, uid int, id, kind string) (*typespec.DatasecTaskCloneTargetsResp, error) {
	taskID, _, err := a.resolveDataSecTask(ctx, uid, id, kind)
	if err != nil {
		return nil, err
	}
	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, taskID)
	if len(targets) == 0 {
		return nil, fmt.Errorf("任务无扫描目标")
	}
	dbType := 0
	out := make([]typespec.DataSecDBTargetInput, 0, len(targets))
	var dataTypes []int
	scanAllDB := false
	scanSensitive := false
	for _, t := range targets {
		cfg := parseDatasecConfig(t.TaskTemplateJSON)
		if dbType == 0 {
			dbType = cfg.DBType
		}
		if cfg.ScanSensitive {
			scanSensitive = true
		}
		if len(cfg.DataTypes) > 0 {
			dataTypes = cfg.DataTypes
		}
		if cfg.ScanAllDB {
			scanAllDB = true
		}
		out = append(out, typespec.DataSecDBTargetInput{
			DBHost:     cfg.DBHost,
			DBPort:     typespec.FlexInt(cfg.DBPort),
			DBName:     cfg.DBName,
			DBUser:     cfg.DBUser,
			DBPassword: cfg.DBPassword,
		})
	}
	return &typespec.DatasecTaskCloneTargetsResp{
		DBType:        dbType,
		Targets:       out,
		ScanSensitive: scanSensitive,
		DataTypes:     dataTypes,
		ScanAllDb:     scanAllDB,
	}, nil
}

// RerunTask 使用历史任务同一批目标再次执行扫描
func (a *DataSecScan) RerunTask(ctx context.Context, uid int, req *typespec.DatasecTaskRerunReq) (*typespec.DataSecScanRunResp, error) {
	clone, err := a.CloneTaskTargets(ctx, uid, req.ID, req.Kind)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = fmt.Sprintf("再次检测-%s", time.Now().Format("20060102-150405"))
	}
	if req.Kind == dataSecScanKindSensitive || req.Kind == "sensitive" {
		runReq := &typespec.DataSecSensitiveRunReq{
			Name:      name,
			DBType:    typespec.FlexInt(clone.DBType),
			Targets:   clone.Targets,
			DataTypes: clone.DataTypes,
			ScanAllDB: clone.ScanAllDb,
		}
		return a.RunSensitiveScan(ctx, uid, runReq)
	}
	runReq := &typespec.DataSecDBRunReq{
		Name:          name,
		DBType:        typespec.FlexInt(clone.DBType),
		Targets:       clone.Targets,
		ScanSensitive: clone.ScanSensitive,
		DataTypes:     clone.DataTypes,
		ScanAllDB:     clone.ScanAllDb,
	}
	return a.RunDBCheck(ctx, uid, runReq)
}

// DeleteTask 删除数据安全扫描任务及其结果
func (a *DataSecScan) DeleteTask(ctx context.Context, uid int, id, kind string) error {
	if strings.TrimSpace(kind) == "" {
		kind = dataSecScanKindDB
	}
	taskID, _, err := a.resolveDataSecTask(ctx, uid, id, kind)
	if err != nil {
		return err
	}
	var taskModel mysqls.TaskTask
	task, err := taskModel.GetTaskCheckTask(ctx, taskID)
	if err != nil || task.ID == 0 {
		return fmt.Errorf("任务不存在")
	}
	if task.Status == enums.TaskStatusRunning {
		return fmt.Errorf("任务运行中，无法删除")
	}

	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, taskID)
	targetIDs := make([]int, 0, len(targets))
	for _, t := range targets {
		if t.Status == enums.TargetStatusRunning {
			return fmt.Errorf("任务运行中，无法删除")
		}
		targetIDs = append(targetIDs, t.ID)
	}

	_ = (&mysqls.DBCheckResult{}).DeleteByTaskID(ctx, taskID)
	_ = (&mysqls.SensitiveDataResult{}).DeleteByTaskID(ctx, taskID)

	if len(targetIDs) > 0 {
		var targetSrv services.TaskTarget
		if err := targetSrv.DelTargetByIds(ctx, targetIDs); err != nil {
			return err
		}
	} else {
		_ = (&mysqls.TaskTarget{}).DeleteByTaskIds(ctx, []int{taskID})
	}

	taskIDs := []int{taskID}
	var taskLogSrv services.TaskLog
	_ = taskLogSrv.DelTaskInfoByTaskId(ctx, taskIDs)
	var taskLogInfoSrv services.TaskLogInfo
	_ = taskLogInfoSrv.DelTaskInfoByTaskId(ctx, taskIDs)
	var taskTaskResultSrv services.TaskTaskResult
	_ = taskTaskResultSrv.DelTaskInfoByTaskId(ctx, taskIDs)
	var taskTaskSrv services.TaskTask
	_ = taskTaskSrv.DelTaskInfoByTaskId(ctx, taskIDs)
	return taskTaskSrv.DelTaskByIds(ctx, taskIDs)
}

func (a *DataSecScan) StopTask(ctx context.Context, uid int, req *typespec.DatasecTaskStopReq) error {
	kind := strings.TrimSpace(req.Kind)
	if kind == "" {
		kind = dataSecScanKindDB
	}
	taskID, _, err := a.resolveDataSecTask(ctx, uid, req.ID, kind)
	if err != nil {
		return err
	}

	var taskModel mysqls.TaskTask
	task, err := taskModel.GetTaskCheckTask(ctx, taskID)
	if err != nil || task.ID == 0 {
		return fmt.Errorf("任务不存在")
	}
	if task.Status != enums.TaskStatusRunning && task.Status != enums.TaskStatusBegin && task.Status != enums.TaskStatusPausing {
		return fmt.Errorf("任务未在运行中")
	}

	if cancel := takeDataSecTaskCancel(taskID); cancel != nil {
		cancel()
	}

	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, taskID)
	var targetModel mysqls.TaskTarget
	var taskLogSrv services.TaskLog
	now := time.Now()
	for _, target := range targets {
		if target.Status == enums.TargetStatusFinish {
			continue
		}
		meta := parseDatasecExtend(target.ExtendField)
		meta.ErrorMessage = "任务已手动结束"
		metaBytes, _ := json.Marshal(meta)
		_ = targetModel.UpdateTargetById(ctx, target.ID, map[string]interface{}{
			"status":       enums.TargetStatusFinish,
			"extend_field": string(metaBytes),
			"update_time":  now,
			"end_time":     now,
		})
		_ = taskLogSrv.UpdateTaskLogStateByTargetId(ctx, target.ID, enums.TaskStatusFinish)
	}

	_ = taskModel.UpdateTaskTaskByIds(ctx, []int{taskID}, map[string]interface{}{
		"status":      enums.TaskStatusFinish,
		"update_time": now,
	})
	var taskInfo mysqls.TaskTaskInfo
	_ = taskInfo.UpdateTaskInfoByTaskIds(ctx, []int{taskID}, map[string]interface{}{
		"status":      enums.TaskStatusFinish,
		"update_time": now,
	})
	return nil
}

func (a *DataSecScan) resolveDataSecTask(ctx context.Context, uid int, id, kind string) (int, int, error) {
	taskID, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil || taskID <= 0 {
		return 0, 0, fmt.Errorf("任务不存在")
	}
	var taskModel mysqls.TaskTask
	task, err := taskModel.GetTaskCheckTask(ctx, taskID)
	if err != nil || task.ID == 0 {
		return 0, 0, fmt.Errorf("任务不存在")
	}
	wantType := enums.TaskTypeDataSecDB
	if kind == dataSecScanKindSensitive || kind == "sensitive" {
		wantType = enums.TaskTypeDataSecSensitive
	}
	if task.TaskType != wantType {
		return 0, 0, fmt.Errorf("任务类型不匹配")
	}
	if uid > 0 && task.UserID != uid {
		return 0, 0, fmt.Errorf("无权操作该任务")
	}
	return taskID, task.TaskType, nil
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
	summary, _ := formatTargetSummary(targets)

	critical, high, mid, low := countDBCheckRisksByTask(ctx, task.ID)
	sHigh, sMid, sLow, sTotal := countSensitiveByTask(ctx, task.ID)
	baselineTotal, baselineFail, cveMatch := summarizeDBBaselineByTask(ctx, task.ID)
	item := typespec.DataSecDBListItem{
		ID:              strconv.Itoa(task.ID),
		Name:            task.TaskName,
		DBType:          cfg.DBType,
		DBHost:          cfg.DBHost,
		DBPort:          cfg.DBPort,
		DBName:          cfg.DBName,
		TargetSummary:   summary,
		TargetCount:     len(targets),
		RiskLevel:       aggregateFrontendRisk(critical, high, mid, low, task.RiskLevel),
		Status:          resolveAppSecTaskStatus(task, targets),
		CreateTime:      formatAppSecTime(task.CreateTime),
		CheckTime:       formatAppSecTime(task.UpdateTime),
		CriticalCount:   critical,
		HighRiskCount:   high,
		MiddleRiskCount: mid,
		LowRiskCount:    low,
		BaselineTotal:   baselineTotal,
		BaselineFail:    baselineFail,
		CveMatchCount:   cveMatch,
		ScanSensitive:   cfg.ScanSensitive || sTotal > 0,
		TotalCount:      sTotal,
		HighCount:       sHigh,
		MediumCount:     sMid,
		LowCount:        sLow,
	}
	if len(targets) > 1 {
		item.DBHost = summary
	}
	if withDetails {
		item.Targets = buildDataSecDBTargetItems(ctx, targets)
		item.Items = loadDBCheckDetailItems(ctx, task.ID, 0)
		if item.ScanSensitive || sTotal > 0 {
			results, _ := (&mysqls.SensitiveDataResult{}).GetByTaskID(ctx, task.ID)
			item.TypeStats = buildSensitiveTypeStats(results)
			item.SensitiveItems = buildSensitiveDetailItems(results, buildSensitiveTargetLabelMap(targets))
		}
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
	summary, _ := formatTargetSummary(targets)
	high, mid, low, total := countSensitiveByTask(ctx, task.ID)
	item := typespec.DataSecSensitiveListItem{
		ID:            strconv.Itoa(task.ID),
		Name:          task.TaskName,
		DBType:        cfg.DBType,
		DBHost:        cfg.DBHost,
		DBPort:        cfg.DBPort,
		DBName:        cfg.DBName,
		TargetSummary: summary,
		TargetCount:   len(targets),
		TotalCount:    total,
		HighCount:     high,
		MediumCount:   mid,
		LowCount:      low,
		Status:        resolveAppSecTaskStatus(task, targets),
		CreateTime:    formatAppSecTime(task.CreateTime),
		ScanTime:      formatAppSecTime(task.UpdateTime),
	}
	if len(targets) > 1 {
		item.DBHost = summary
	}
	if withDetails {
		item.Targets = buildDataSecSensitiveTargetItems(ctx, targets)
		results, _ := (&mysqls.SensitiveDataResult{}).GetByTaskID(ctx, task.ID)
		item.TypeStats = buildSensitiveTypeStats(results)
		item.Items = buildSensitiveDetailItems(results, buildSensitiveTargetLabelMap(targets))
	}
	return item, nil
}

func buildDataSecDBTargetItems(ctx context.Context, targets []mysqls.TaskTarget) []typespec.DataSecTargetItem {
	out := make([]typespec.DataSecTargetItem, 0, len(targets))
	for _, t := range targets {
		cfg := parseDatasecConfig(t.TaskTemplateJSON)
		ext := parseDatasecExtend(t.ExtendField)
		c, h, m, l := countDBCheckRisksByTarget(ctx, t.ID)
		sh, sm, sl, stotal := countSensitiveByTarget(ctx, t.ID)
		checkRows, _ := (&mysqls.DBCheckResult{}).GetByTargetID(ctx, t.ID)
		meta := services.SummarizeDBCheckResults(checkRows)
		out = append(out, typespec.DataSecTargetItem{
			ID:              t.ID,
			TargetURL:       t.TargetURL,
			DBType:          cfg.DBType,
			DBHost:          cfg.DBHost,
			DBPort:          cfg.DBPort,
			DBName:          cfg.DBName,
			Status:          mapTargetStatusToAppSecOrDefault(t.Status),
			RiskLevel:       aggregateFrontendRisk(c, h, m, l, mapTargetRiskToTaskRisk(t.RiskLevel)),
			DbVersion:       meta.DbVersion,
			BaselineTotal:   meta.BaselineTotal,
			BaselinePass:    meta.BaselinePass,
			BaselineFail:    meta.BaselineFail,
			BaselineError:   meta.BaselineError,
			CveMatchCount:   meta.CveMatchCount,
			CriticalCount:   c,
			HighRiskCount:   h,
			MiddleRiskCount: m,
			LowRiskCount:    l,
			TotalCount:      stotal,
			HighCount:       sh,
			MediumCount:     sm,
			LowCount:        sl,
			ErrorMessage:    ext.ErrorMessage,
		})
	}
	return out
}

func buildDataSecSensitiveTargetItems(ctx context.Context, targets []mysqls.TaskTarget) []typespec.DataSecTargetItem {
	out := make([]typespec.DataSecTargetItem, 0, len(targets))
	for _, t := range targets {
		cfg := parseDatasecConfig(t.TaskTemplateJSON)
		ext := parseDatasecExtend(t.ExtendField)
		high, mid, low, total := countSensitiveByTarget(ctx, t.ID)
		out = append(out, typespec.DataSecTargetItem{
			ID:           t.ID,
			TargetURL:    t.TargetURL,
			DBType:       cfg.DBType,
			DBHost:       cfg.DBHost,
			DBPort:       cfg.DBPort,
			DBName:       cfg.DBName,
			Status:       mapTargetStatusToAppSecOrDefault(t.Status),
			TotalCount:   total,
			HighCount:    high,
			MediumCount:  mid,
			LowCount:     low,
			ErrorMessage: ext.ErrorMessage,
		})
	}
	return out
}

func loadDBCheckDetailItems(ctx context.Context, taskID int, targetID int) []typespec.DataSecDBDetailItem {
	list := loadDBCheckResultsByTask(ctx, taskID)
	out := make([]typespec.DataSecDBDetailItem, 0, len(list))
	for _, r := range list {
		if targetID > 0 && r.TargetID != targetID {
			continue
		}
		risk := r.RiskLevel
		if r.CheckResult == enums.BaselineCheckResultPass {
			risk = enums.BaselineRiskInfo
		}
		out = append(out, typespec.DataSecDBDetailItem{
			TargetID:      r.TargetID,
			Category:      r.CheckCategory,
			RiskLevel:     risk,
			Result:        enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			Description:   firstNonEmpty(r.RiskDescription, r.RuleName),
			Suggestion:    r.FixSuggestion,
			RuleName:      r.RuleName,
			ActualValue:   r.ActualValue,
			ExpectedValue: r.ExpectedValue,
			IsCve:         services.IsCveDBCheckResult(r),
		})
	}
	return out
}

func summarizeDBBaselineByTask(ctx context.Context, taskID int) (baselineTotal, baselineFail, cveMatch int) {
	list := loadDBCheckResultsByTask(ctx, taskID)
	meta := services.SummarizeDBCheckResults(list)
	return meta.BaselineTotal, meta.BaselineFail, meta.CveMatchCount
}

func loadDBCheckResultsByTask(ctx context.Context, taskID int) []mysqls.DBCheckResult {
	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, taskID)
	if len(targets) > 0 {
		out := make([]mysqls.DBCheckResult, 0)
		for _, t := range targets {
			rows, err := (&mysqls.DBCheckResult{}).GetByTargetID(ctx, t.ID)
			if err != nil {
				continue
			}
			out = append(out, rows...)
		}
		if len(out) > 0 {
			return out
		}
	}
	list, _ := (&mysqls.DBCheckResult{}).GetByTaskID(ctx, taskID)
	return list
}

func countDBCheckRisksByTarget(ctx context.Context, targetID int) (critical, high, mid, low int) {
	list, err := (&mysqls.DBCheckResult{}).GetByTargetID(ctx, targetID)
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

func countSensitiveByTarget(ctx context.Context, targetID int) (high, mid, low, total int) {
	list, err := (&mysqls.SensitiveDataResult{}).GetByTargetID(ctx, targetID)
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

func countDBCheckRisksByTask(ctx context.Context, taskID int) (critical, high, mid, low int) {
	list := loadDBCheckResultsByTask(ctx, taskID)
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

func worstCombinedTaskRisk(a, b int) int {
	if a == 0 {
		a = enums.TaskRiskSafe
	}
	if b == 0 {
		b = enums.TaskRiskSafe
	}
	if a < b {
		return a
	}
	return b
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

func buildSensitiveTargetLabelMap(targets []mysqls.TaskTarget) map[int]string {
	out := make(map[int]string, len(targets))
	for _, t := range targets {
		label := strings.TrimSpace(t.TargetURL)
		if label == "" {
			cfg := parseDatasecConfig(t.TaskTemplateJSON)
			label = formatDataSecTargetURL(cfg.DBHost, cfg.DBPort, cfg.DBName)
		}
		out[t.ID] = label
	}
	return out
}

func buildSensitiveLocation(dbName, tableName, columnName string) string {
	parts := make([]string, 0, 3)
	if strings.TrimSpace(dbName) != "" {
		parts = append(parts, strings.TrimSpace(dbName))
	}
	if strings.TrimSpace(tableName) != "" {
		parts = append(parts, strings.TrimSpace(tableName))
	}
	if strings.TrimSpace(columnName) != "" {
		parts = append(parts, strings.TrimSpace(columnName))
	}
	return strings.Join(parts, " / ")
}

func buildSensitiveDetailItems(results []mysqls.SensitiveDataResult, targetLabels map[int]string) []typespec.DataSecSensitiveDetailItem {
	out := make([]typespec.DataSecSensitiveDetailItem, 0, len(results))
	for _, r := range results {
		cnt := 1
		if r.TotalRows > 0 {
			cnt = int(r.TotalRows)
		}
		location := buildSensitiveLocation(r.DBName, r.TableNameStr, r.ColumnName)
		out = append(out, typespec.DataSecSensitiveDetailItem{
			TargetID:         r.TargetID,
			TargetLabel:      targetLabels[r.TargetID],
			DBName:           r.DBName,
			TableName:        r.TableNameStr,
			ColumnName:       r.ColumnName,
			Location:         location,
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
