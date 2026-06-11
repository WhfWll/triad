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
	"smart/tools/data"
	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

const (
	appSecStatusWaiting = 1
	appSecStatusRunning = 2
	appSecStatusDone    = 3
)

// AppSecScan 应用安全 Web 扫描（task_task / task_target / task_vul）
type AppSecScan struct{}

var (
	appSecTaskCancelMu sync.Mutex
	appSecTaskCancels  = make(map[int]context.CancelFunc)
)

func setAppSecTaskCancel(taskID int, cancel context.CancelFunc) {
	appSecTaskCancelMu.Lock()
	appSecTaskCancels[taskID] = cancel
	appSecTaskCancelMu.Unlock()
}

func takeAppSecTaskCancel(taskID int) context.CancelFunc {
	appSecTaskCancelMu.Lock()
	cancel := appSecTaskCancels[taskID]
	delete(appSecTaskCancels, taskID)
	appSecTaskCancelMu.Unlock()
	return cancel
}

type appSecExtendMeta struct {
	AppSecScanType string `json:"appsecScanType"`
	AppType        int    `json:"appType"`
	StrategyID     string `json:"strategyId"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}

func (a *AppSecScan) RunDynamicScan(ctx context.Context, uid int, req *typespec.AppSecScanRunReq) (*typespec.AppSecScanRunResp, error) {
	return a.runScan(ctx, uid, "dyn", enums.TaskTypeAppSecDynamic, req)
}

func (a *AppSecScan) RunAppSpecificScan(ctx context.Context, uid int, req *typespec.AppSecScanRunReq) (*typespec.AppSecScanRunResp, error) {
	return a.runScan(ctx, uid, "app", enums.TaskTypeAppSecApp, req)
}

func (a *AppSecScan) ListDynamicScans(ctx context.Context, uid int, req *typespec.AppSecScanListReq) (*typespec.AppSecScanListResp, error) {
	return a.listFromDB(ctx, uid, enums.TaskTypeAppSecDynamic, "dyn", req)
}

func (a *AppSecScan) ListAppSpecificScans(ctx context.Context, uid int, req *typespec.AppSecScanListReq) (*typespec.AppSecScanListResp, error) {
	return a.listFromDB(ctx, uid, enums.TaskTypeAppSecApp, "app", req)
}

func (a *AppSecScan) GetDynamicScanDetail(ctx context.Context, uid int, id string) (*typespec.AppSecTaskItem, error) {
	return a.getDetailFromDB(ctx, uid, enums.TaskTypeAppSecDynamic, "dyn", id)
}

func (a *AppSecScan) GetAppSpecificScanDetail(ctx context.Context, uid int, id string) (*typespec.AppSecTaskItem, error) {
	return a.getDetailFromDB(ctx, uid, enums.TaskTypeAppSecApp, "app", id)
}

func (a *AppSecScan) StopTask(ctx context.Context, uid int, req *typespec.AppSecTaskStopReq) (*typespec.AppSecTaskStopResp, error) {
	taskID, err := strconv.Atoi(strings.TrimSpace(req.ID))
	if err != nil || taskID <= 0 {
		return nil, fmt.Errorf("任务不存在")
	}
	wantType := enums.TaskTypeAppSecDynamic
	kind := strings.TrimSpace(req.Kind)
	if kind == "app" {
		wantType = enums.TaskTypeAppSecApp
	}

	var taskModel mysqls.TaskTask
	task, err := taskModel.GetTaskCheckTask(ctx, taskID)
	if err != nil || task.ID == 0 || task.TaskType != wantType {
		return nil, fmt.Errorf("任务不存在")
	}
	if uid > 0 && task.UserID != uid {
		return nil, fmt.Errorf("无权操作该任务")
	}

	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, taskID)
	if len(targets) == 0 {
		return nil, fmt.Errorf("任务无扫描目标")
	}

	now := time.Now()
	if cancel := takeAppSecTaskCancel(taskID); cancel != nil {
		cancel()
	}
	var targetModel mysqls.TaskTarget
	var taskLogSrv services.TaskLog
	stopped := false
	for _, target := range targets {
		if target.Status == enums.TargetStatusFinish {
			continue
		}
		services.CancelTargetScan(target.ID)
		meta := parseAppSecExtend(target.ExtendField)
		meta.ErrorMessage = "任务已手动结束"
		metaBytes, _ := json.Marshal(meta)
		_ = targetModel.UpdateTargetById(ctx, target.ID, map[string]interface{}{
			"status":       enums.TargetStatusFinish,
			"extend_field": string(metaBytes),
			"update_time":  now,
			"end_time":     now,
		})
		_ = taskLogSrv.UpdateTaskLogStateByTargetId(ctx, target.ID, enums.TaskStatusFinish)
		stopped = true
	}

	if !stopped && task.Status != enums.TaskStatusRunning && task.Status != enums.TaskStatusBegin && task.Status != enums.TaskStatusPausing {
		return nil, fmt.Errorf("任务未在运行中")
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
	return &typespec.AppSecTaskStopResp{Stopped: true}, nil
}

func (a *AppSecScan) runScan(ctx context.Context, uid int, scanType string, taskType int, req *typespec.AppSecScanRunReq) (*typespec.AppSecScanRunResp, error) {
	targetList, err := parseAppSecTargets(req)
	if err != nil {
		return nil, err
	}

	configJson, err := appSecPayloadToConfigJson(req)
	if err != nil {
		return nil, err
	}
	if _, err := resolveAppSecVulLibraries(ctx, configJson); err != nil {
		return nil, err
	}

	taskID, targets, err := a.createAppSecTask(ctx, uid, taskType, scanType, req, targetList, configJson)
	if err != nil {
		return nil, err
	}

	scanCtx, cancel := context.WithCancel(context.Background())
	setAppSecTaskCancel(taskID, cancel)
	sem := make(chan struct{}, services.GetAppScanConcurrent(ctx))
	for _, row := range targets {
		t := row
		go func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			a.startTargetScan(scanCtx, t, configJson)
		}()
	}

	return &typespec.AppSecScanRunResp{ID: strconv.Itoa(taskID)}, nil
}

func parseAppSecTargets(req *typespec.AppSecScanRunReq) ([]string, error) {
	raw := strings.TrimSpace(req.Target)
	if raw == "" {
		raw = strings.TrimSpace(req.TargetURL)
	}
	if raw == "" {
		return nil, fmt.Errorf("扫描目标不能为空")
	}
	var analysis data.TaskCheckTaskAnalysisTarget
	analysis.AnalysisTarget(raw, "")
	if len(analysis.ErrorTargetList) > 0 {
		return nil, fmt.Errorf("目标格式错误: %s", strings.Join(analysis.ErrorTargetList, "; "))
	}
	if len(analysis.TargetList) == 0 {
		return nil, fmt.Errorf("请至少填写一个有效扫描目标")
	}
	out := make([]string, 0, len(analysis.TargetList))
	seen := make(map[string]struct{})
	for _, t := range analysis.TargetList {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if !strings.HasPrefix(t, "http://") && !strings.HasPrefix(t, "https://") {
			t = "http://" + t
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("请至少填写一个有效扫描目标")
	}
	return out, nil
}

func (a *AppSecScan) createAppSecTask(ctx context.Context, uid, taskType int, scanType string, req *typespec.AppSecScanRunReq, targetList []string, configJson enums.ConfigJson) (int, []mysqls.TaskTarget, error) {
	var taskSrv services.TaskTask
	taskName := firstNonEmpty(req.Name, "未命名任务")
	taskID, err := taskSrv.Save(ctx, uid, 0, taskName, targetList, enums.TaskExecTypeImmediate, "{}", configJson, time.Unix(0, 0), 0, 0)
	if err != nil {
		return 0, nil, err
	}

	ext := appSecExtendMeta{
		AppSecScanType: scanType,
		AppType:        req.AppType,
		StrategyID:     req.Strategy,
	}
	extBytes, _ := json.Marshal(ext)
	cfgJSON := mustJSON(configJson)

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

func (a *AppSecScan) startTargetScan(ctx context.Context, target mysqls.TaskTarget, configJson enums.ConfigJson) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	current, _ := (&mysqls.TaskTarget{}).GetTaskTarget(ctx, target.ID)
	if current.ID > 0 && current.Status == enums.TargetStatusFinish {
		return
	}
	// 与创建任务时一致：未选插件则加载全部可用 yak/nuclei，并写回目标配置供 scanner 使用
	if vulLibs, err := resolveAppSecVulLibraries(ctx, configJson); err != nil {
		log.Errorf("appsec resolve vul libraries: %v", err)
		return
	} else if len(configJson.VulIdsConfig) == 0 && len(vulLibs) > 0 {
		ids := make([]int, 0, len(vulLibs))
		for _, v := range vulLibs {
			ids = append(ids, v.ID)
		}
		configJson.VulIdsConfig = ids
		target.TaskTemplateJSON = mustJSON(configJson)
		_ = target.UpdateTaskTarget(ctx)
	}

	var targetSrv services.TaskTarget
	logId, err := targetSrv.UpdateTargetAndSaveTargetLog(ctx, target.TaskID, target.ID, enums.TargetStatusRunning, enums.TargetIsAliveY, target.TargetURL)
	if err != nil {
		log.Errorf("appsec startTargetScan log error: %v", err)
		return
	}
	services.NewVulnerabilityScanner().StartLocalScan(ctx, logId, target, configJson)
}

func (a *AppSecScan) listFromDB(ctx context.Context, uid, taskType int, scanType string, req *typespec.AppSecScanListReq) (*typespec.AppSecScanListResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 10
	}

	var taskModel mysqls.TaskTask
	tasks, total, err := taskModel.GetAppSecTaskList(ctx, taskType, uid, page, size, strings.TrimSpace(req.Search))
	if err != nil {
		return nil, err
	}

	list := make([]typespec.AppSecTaskItem, 0, len(tasks))
	for _, t := range tasks {
		item, err := a.buildTaskItem(ctx, t, scanType, false)
		if err != nil {
			continue
		}
		list = append(list, *item)
	}
	return &typespec.AppSecScanListResp{List: list, Total: int(total)}, nil
}

func (a *AppSecScan) getDetailFromDB(ctx context.Context, uid, taskType int, scanType, id string) (*typespec.AppSecTaskItem, error) {
	taskID, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil || taskID <= 0 {
		return nil, fmt.Errorf("任务不存在")
	}
	var taskModel mysqls.TaskTask
	task, err := taskModel.GetTaskCheckTask(ctx, taskID)
	if err != nil || task.ID == 0 {
		return nil, fmt.Errorf("任务不存在")
	}
	if task.TaskType != taskType {
		return nil, fmt.Errorf("任务不存在")
	}
	if uid > 0 && task.UserID != uid {
		return nil, fmt.Errorf("无权查看该任务")
	}
	return a.buildTaskItem(ctx, task, scanType, true)
}

func (a *AppSecScan) buildTaskItem(ctx context.Context, task mysqls.TaskTask, scanType string, withDetails bool) (*typespec.AppSecTaskItem, error) {
	targets := (&mysqls.TaskTarget{}).GetTargetsByTaskId(ctx, task.ID)
	var meta appSecExtendMeta
	if len(targets) > 0 {
		meta = parseAppSecExtend(targets[0].ExtendField)
	}

	critical, high, mid, low := countTaskVulnsByRisk(ctx, task.ID)
	vulnTotal := critical + high + mid + low
	pageTotal := 0
	targetItems := make([]typespec.AppSecTargetItem, 0, len(targets))
	for _, t := range targets {
		tc, th, tm, tl := countTaskVulnsByRiskForTarget(ctx, task.ID, t.ID)
		pages := countCrawlerPages(ctx, t.ID)
		pageTotal += pages
		targetItems = append(targetItems, typespec.AppSecTargetItem{
			ID:              t.ID,
			TargetURL:       t.TargetURL,
			Status:          mapTargetStatusToAppSecOrDefault(t.Status),
			PageCount:       pages,
			VulnCount:       tc + th + tm + tl,
			CriticalCount:   tc,
			HighRiskCount:   th,
			MiddleRiskCount: tm,
			LowRiskCount:    tl,
		})
	}

	summary, primaryURL := formatTargetSummary(targets)
	item := &typespec.AppSecTaskItem{
		ID:              strconv.Itoa(task.ID),
		Name:            task.TaskName,
		TargetURL:       primaryURL,
		TargetSummary:   summary,
		TargetCount:     len(targets),
		Targets:         targetItems,
		StrategyID:      meta.StrategyID,
		AppType:         meta.AppType,
		Status:          resolveAppSecTaskStatus(task, targets),
		RiskLevel:       aggregateFrontendRisk(critical, high, mid, low, task.RiskLevel),
		PageCount:       pageTotal,
		VulnCount:       vulnTotal,
		CriticalCount:   critical,
		HighRiskCount:   high,
		MiddleRiskCount: mid,
		LowRiskCount:    low,
		CreateTime:      formatAppSecTime(task.CreateTime),
		ScanTime:        formatAppSecTime(task.UpdateTime),
		ErrorMessage:    meta.ErrorMessage,
		Vulns:           []typespec.AppSecVulnItem{},
		Pages:           []typespec.AppSecPageItem{},
	}

	if withDetails {
		item.Vulns = loadTaskVulns(ctx, task.ID, scanType)
		if len(targets) > 0 {
			item.TargetID = targets[0].ID
			item.ScanConfig = parseAppSecScanConfig(targets[0].TaskTemplateJSON)
			allPages := make([]typespec.AppSecPageItem, 0)
			for _, t := range targets {
				for _, p := range loadCrawlerPages(ctx, t.ID) {
					p.TargetID = t.ID
					allPages = append(allPages, p)
				}
			}
			item.Pages = allPages
		}
	}

	return item, nil
}

func formatTargetSummary(targets []mysqls.TaskTarget) (summary, primary string) {
	if len(targets) == 0 {
		return "", ""
	}
	primary = targets[0].TargetURL
	if len(targets) == 1 {
		return primary, primary
	}
	return fmt.Sprintf("%s 等 %d 个目标", primary, len(targets)), primary
}

func mapTargetStatusToAppSecOrDefault(targetStatus int) int {
	if s, ok := mapTargetStatusToAppSec(targetStatus); ok {
		return s
	}
	return appSecStatusWaiting
}

func countTaskVulnsByRiskForTarget(ctx context.Context, taskID, targetID int) (critical, high, mid, low int) {
	type row struct {
		Risk  int
		Count int64
	}
	var rows []row
	db := mysql.FromContext(ctx).Model(&mysqls.TaskVul{})
	_ = db.Select("risk, count(*) as count").
		Where("task_id = ? AND target_id = ? AND data_type = ?", taskID, targetID, enums.VulDataTypOne).
		Group("risk").Scan(&rows).Error
	for _, r := range rows {
		switch r.Risk {
		case enums.VulLibrariesRiskDead:
			critical += int(r.Count)
		case enums.VulLibrariesRiskHigh:
			high += int(r.Count)
		case enums.VulLibrariesRiskMiddle:
			mid += int(r.Count)
		case enums.VulLibrariesRiskLow:
			low += int(r.Count)
		}
	}
	return
}

func parseAppSecScanConfig(raw string) map[string]interface{} {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil
	}
	return cfg
}

func loadTaskVulns(ctx context.Context, taskID int, scanType string) []typespec.AppSecVulnItem {
	var vulModel mysqls.TaskVul
	list, _, err := vulModel.GetTaskVulList(ctx, taskID, 0, 0, 0, "", enums.VulDataTypOne, 1, 500, 0)
	if err != nil {
		return nil
	}
	out := make([]typespec.AppSecVulnItem, 0, len(list))
	for _, v := range list {
		out = append(out, typespec.AppSecVulnItem{
			TargetID:    v.TargetID,
			TargetURL:   v.TargetUrl,
			Name:        v.Name,
			Type:        v.Type,
			TypeName:    enums.ToolsVulnerabilityEnum.GetTypeEnum(v.Type),
			RiskLevel:   vulLibRiskToFrontend(v.Risk),
			URL:         firstNonEmpty(v.Location, v.VulAddress, v.TargetUrl),
			Description: v.Description,
			Payload:     v.VulParam,
			// 验证报文常含请求响应，前端详情页展示
			Request:    v.VerMsg,
			Suggestion: v.FixSuggest,
			Method:     "GET",
		})
	}
	_ = scanType
	return out
}

func loadCrawlerPages(ctx context.Context, targetID int) []typespec.AppSecPageItem {
	var resultModel mysqls.TaskTaskResult
	list, err := resultModel.GetTaskTaskResultByType(ctx, enums.TaskResultObjTypeInfo, enums.TaskResultSubObjTypeUrl, []string{strconv.Itoa(targetID)})
	if err != nil {
		return nil
	}
	pages := make([]typespec.AppSecPageItem, 0, len(list))
	seen := make(map[string]struct{})
	for _, r := range list {
		url := strings.TrimSpace(r.Field1)
		if url == "" {
			continue
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		pages = append(pages, typespec.AppSecPageItem{URL: url, Name: url, TargetID: targetID})
	}
	return pages
}

func countCrawlerPages(ctx context.Context, targetID int) int {
	return len(loadCrawlerPages(ctx, targetID))
}

func countTaskVulnsByRisk(ctx context.Context, taskID int) (critical, high, mid, low int) {
	type row struct {
		Risk  int
		Count int64
	}
	var rows []row
	db := mysql.FromContext(ctx).Model(&mysqls.TaskVul{})
	_ = db.Select("risk, count(*) as count").
		Where("task_id = ? AND data_type = ?", taskID, enums.VulDataTypOne).
		Group("risk").Scan(&rows).Error
	for _, r := range rows {
		switch r.Risk {
		case enums.VulLibrariesRiskDead:
			critical += int(r.Count)
		case enums.VulLibrariesRiskHigh:
			high += int(r.Count)
		case enums.VulLibrariesRiskMiddle:
			mid += int(r.Count)
		case enums.VulLibrariesRiskLow:
			low += int(r.Count)
		}
	}
	return
}

// resolveAppSecTaskStatus 聚合多目标状态：任一运行中→扫描中，全部完成→已完成，否则等待/进行中
func resolveAppSecTaskStatus(task mysqls.TaskTask, targets []mysqls.TaskTarget) int {
	if len(targets) == 0 {
		return mapTaskStatusToFrontend(task.Status)
	}
	running, waiting, doneCount := 0, 0, 0
	for _, t := range targets {
		s := mapTargetStatusToAppSecOrDefault(t.Status)
		switch s {
		case appSecStatusRunning:
			running++
		case appSecStatusDone:
			doneCount++
		case appSecStatusWaiting:
			waiting++
		}
	}
	if running > 0 {
		return appSecStatusRunning
	}
	if doneCount == len(targets) {
		return appSecStatusDone
	}
	if waiting == len(targets) {
		return appSecStatusWaiting
	}
	if doneCount > 0 {
		return appSecStatusRunning
	}
	return mapTaskStatusToFrontend(task.Status)
}

func mapTargetStatusToAppSec(targetStatus int) (int, bool) {
	switch targetStatus {
	case enums.TargetStatusToTrigger, enums.TargetStatusToBegin:
		return appSecStatusWaiting, true
	case enums.TargetStatusRunning, enums.TargetStatusPausing:
		return appSecStatusRunning, true
	case enums.TargetStatusFinish:
		return appSecStatusDone, true
	default:
		return 0, false
	}
}

func mapTaskStatusToFrontend(status int) int {
	switch status {
	case enums.TaskStatusTrigger, enums.TaskStatusBegin:
		return appSecStatusWaiting
	case enums.TaskStatusRunning:
		return appSecStatusRunning
	case enums.TaskStatusFinish:
		return appSecStatusDone
	default:
		return appSecStatusWaiting
	}
}

func aggregateFrontendRisk(critical, high, mid, low, taskRisk int) int {
	if critical > 0 {
		return 0
	}
	if high > 0 {
		return 1
	}
	if mid > 0 {
		return 2
	}
	if low > 0 {
		return 3
	}
	switch taskRisk {
	case enums.TaskRiskHigh:
		return 1
	case enums.TaskRiskMid:
		return 2
	case enums.TaskRiskLow:
		return 3
	default:
		return 4
	}
}

func parseAppSecExtend(raw string) appSecExtendMeta {
	var meta appSecExtendMeta
	if raw == "" {
		return meta
	}
	_ = json.Unmarshal([]byte(raw), &meta)
	return meta
}

func mustJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func appSecPayloadToConfigJson(req *typespec.AppSecScanRunReq) (enums.ConfigJson, error) {
	raw, err := json.Marshal(req)
	if err != nil {
		return enums.ConfigJson{}, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return enums.ConfigJson{}, err
	}
	rename := map[string]string{
		"webCrawler":   "webCrawlerConfig",
		"portScan":     "portScanConfig",
		"proxy":        "proxyConfig",
		"webPathScan":  "webPathScanConfig",
		"weakPass":     "weakPassConfig",
		"websiteLogin": "websiteLoginConfig",
	}
	for oldKey, newKey := range rename {
		if v, ok := m[oldKey]; ok {
			m[newKey] = v
			delete(m, oldKey)
		}
	}
	delete(m, "name")
	delete(m, "target")
	delete(m, "targetUrl")
	delete(m, "appType")
	delete(m, "strategy")

	b, err := json.Marshal(m)
	if err != nil {
		return enums.ConfigJson{}, err
	}
	var cfg enums.ConfigJson
	if err := json.Unmarshal(b, &cfg); err != nil {
		return enums.ConfigJson{}, err
	}
	if cfg.Mode.Mode == 0 {
		cfg.Mode.Mode = enums.TaskConfigurationModeCommon
	}
	return cfg, nil
}

func resolveAppSecVulLibraries(ctx context.Context, configJson enums.ConfigJson) ([]mysqls.VulLibraries, error) {
	var lib services.VulLibraries
	var out []mysqls.VulLibraries
	var err error
	if len(configJson.VulIdsConfig) > 0 {
		out, err = lib.GetVulLibsByIds(ctx, configJson.VulIdsConfig)
	} else {
		out, err = lib.GetAppSecDefaultVulLibraries(ctx, configJson.SafeTest)
	}
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("没有可用的漏洞检测插件，请在插件页选择或检查漏洞库")
	}
	return out, nil
}

func vulLibRiskToFrontend(risk int) int {
	switch risk {
	case enums.VulLibrariesRiskDead:
		return 0
	case enums.VulLibrariesRiskHigh:
		return 1
	case enums.VulLibrariesRiskMiddle:
		return 2
	case enums.VulLibrariesRiskLow:
		return 3
	default:
		return 4
	}
}

func formatAppSecTime(t time.Time) string {
	if t.IsZero() || t.Year() < 1980 {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
