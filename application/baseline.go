package application

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

type BaselineApp struct{}

func normalizeBaselineCheckReq(req *typespec.BaselineCheckReq) {
	if req.TaskID <= 0 {
		id := int(time.Now().UnixNano() % 2000000000)
		if id <= 0 {
			id = -id + 1
		}
		req.TaskID = id
	}
	if req.ScanScene != enums.HostScanSceneHostVuln {
		req.ScanScene = enums.HostScanSceneBaseline
	}
}

func normalizeMalwareScanReq(req *typespec.MalwareScanReq) {
	if req.Port <= 0 {
		req.Port = 22
	}
	if req.TaskID <= 0 {
		id := int(time.Now().UnixNano() % 2000000000)
		if id <= 0 {
			id = -id + 1
		}
		req.TaskID = id
	}
	if req.OSType <= 0 {
		req.OSType = enums.BaselineOSTypeLinux
	}
}

func (a *BaselineApp) RunBaselineCheck(ctx context.Context, req *typespec.BaselineCheckReq) (*typespec.BaselineCheckResp, error) {
	normalizeBaselineCheckReq(req)
	checker := services.GetHostBaselineChecker()
	task := &services.BaselineCheckTask{
		TaskID:        req.TaskID,
		TargetID:      req.TargetID,
		Host:          req.Host,
		Port:          req.Port,
		Username:      req.Username,
		Password:      req.Password,
		Key:           req.Key,
		OSType:        req.OSType,
		Transport:     req.Transport,
		WinRMUseHttps: req.WinRMUseHttps,
		ScanScene:     req.ScanScene,
	}

	report, err := checker.RunBaselineCheck(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := &typespec.BaselineCheckResp{
		TaskID:     report.TaskID,
		TargetIP:   report.TargetIP,
		OSType:     report.OSType,
		OSTypeName: enums.BaselineEnum.GetOSTypeName(report.OSType),
		TotalRules: report.TotalRules,
		PassCount:  report.PassCount,
		FailCount:  report.FailCount,
		ErrorCount: report.ErrorCount,
		StartTime:  report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:    report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.BaselineCheckItem{
			ID:              r.ID,
			RuleID:          r.RuleID,
			RuleName:        r.RuleName,
			RuleCategory:    r.RuleCategory,
			CategoryName:    enums.BaselineEnum.GetCategoryName(r.RuleCategory),
			RiskLevel:       r.RuleRisk,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.RuleRisk),
			CheckResult:     r.CheckResult,
			ResultName:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			ExpectedValue:   r.ExpectedValue,
			ActualValue:     r.ActualValue,
			CheckCommand:    r.CheckCommand,
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			CheckTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

// batchTaskManager 管理批量任务的进度状态
var (
	batchTasksMu sync.RWMutex
	batchTasks   = make(map[int]*typespec.BaselineBatchTaskProgress)
)

func registerBatchTask(taskID int, targets []typespec.BaselineCheckReq) {
	progress := &typespec.BaselineBatchTaskProgress{
		TaskID:           taskID,
		Status:           "running",
		TotalTargets:     len(targets),
		CompletedTargets: 0,
		Targets:          make([]typespec.BatchTargetProgress, len(targets)),
		CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
	}
	for i, t := range targets {
		progress.Targets[i] = typespec.BatchTargetProgress{
			Host:   t.Host,
			Status: "pending",
		}
	}
	batchTasksMu.Lock()
	batchTasks[taskID] = progress
	batchTasksMu.Unlock()
}

func updateBatchTargetStatus(taskID int, idx int, status, message string) {
	batchTasksMu.RLock()
	p, ok := batchTasks[taskID]
	batchTasksMu.RUnlock()
	if !ok {
		return
	}
	p.Targets[idx].Status = status
	if message != "" {
		p.Targets[idx].Message = message
	}
	if status == "completed" || status == "failed" {
		p.CompletedTargets++
	}
	if p.CompletedTargets >= p.TotalTargets {
		p.Status = "completed"
	}
}

// RunBaselineBatchCheckAsync 异步批量多目标核查，创建任务后立即返回
func (a *BaselineApp) RunBaselineBatchCheckAsync(ctx context.Context, req *typespec.BaselineBatchCheckReq) (*typespec.BaselineBatchCheckResp, error) {
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("targets is empty")
	}

	batchTaskID := req.TaskID
	if batchTaskID <= 0 {
		batchTaskID = int(time.Now().UnixNano() % 2000000000)
		if batchTaskID <= 0 {
			batchTaskID = -batchTaskID + 1
		}
	}

	registerBatchTask(batchTaskID, req.Targets)

	go a.runBatchChecks(context.Background(), batchTaskID, req.Targets)

	return &typespec.BaselineBatchCheckResp{TaskID: batchTaskID}, nil
}

func (a *BaselineApp) runBatchChecks(ctx context.Context, batchTaskID int, targets []typespec.BaselineCheckReq) {
	maxConcurrency := 10
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for i, target := range targets {
		target := target
		target.TaskID = batchTaskID
		if target.ScanScene <= 0 {
			target.ScanScene = enums.HostScanSceneBaseline
		}
		wg.Add(1)
		go func(idx int, t typespec.BaselineCheckReq) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			updateBatchTargetStatus(batchTaskID, idx, "running", "")
			_, err := a.RunBaselineCheck(ctx, &t)
			if err != nil {
				log.Errorf("runBatchChecks: target %s failed: %v", t.Host, err)
				updateBatchTargetStatus(batchTaskID, idx, "failed", err.Error())
				return
			}
			updateBatchTargetStatus(batchTaskID, idx, "completed", "")
		}(i, target)
	}

	wg.Wait()
	log.Infof("runBatchChecks: batch task %d completed", batchTaskID)
}

// GetBatchTaskProgress 获取批量任务进度
func (a *BaselineApp) GetBatchTaskProgress(ctx context.Context, taskID int) *typespec.BaselineBatchTaskProgress {
	batchTasksMu.RLock()
	defer batchTasksMu.RUnlock()
	p, ok := batchTasks[taskID]
	if !ok {
		return nil
	}
	return p
}

func (a *BaselineApp) GetBaselineResults(ctx context.Context, req *typespec.BaselineCheckResultListReq) (*typespec.BaselineCheckResultListResp, error) {
	var model mysqls.BaselineCheckResult
	var list []mysqls.BaselineCheckResult
	var err error

	if req.TargetID > 0 {
		list, err = model.GetByTargetID(ctx, req.TargetID)
	} else {
		list, err = model.GetByTaskID(ctx, req.TaskID)
	}
	if err != nil {
		return nil, err
	}

	resp := &typespec.BaselineCheckResultListResp{
		Total: int64(len(list)),
	}
	for _, r := range list {
		resp.List = append(resp.List, typespec.BaselineCheckItem{
			ID:              r.ID,
			RuleID:          r.RuleID,
			RuleName:        r.RuleName,
			RuleCategory:    r.RuleCategory,
			CategoryName:    enums.BaselineEnum.GetCategoryName(r.RuleCategory),
			RiskLevel:       r.RuleRisk,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.RuleRisk),
			CheckResult:     r.CheckResult,
			ResultName:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			ExpectedValue:   r.ExpectedValue,
			ActualValue:     r.ActualValue,
			CheckCommand:    r.CheckCommand,
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			CheckTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) GetBaselineStat(ctx context.Context, req *typespec.BaselineStatReq) (*typespec.BaselineStatResp, error) {
	var model mysqls.BaselineCheckResult
	pass, fail, total, err := model.GetStatByTaskID(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}

	rate := 0.0
	if total > 0 {
		rate = float64(pass) / float64(total) * 100
	}

	return &typespec.BaselineStatResp{
		TotalRules: int(total),
		PassCount:  int(pass),
		FailCount:  int(fail),
		PassRate:   rate,
	}, nil
}

func (a *BaselineApp) GetBaselineTaskList(ctx context.Context, req *typespec.BaselineTaskListReq) (*typespec.BaselineTaskListResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 10
	}
	var model mysqls.BaselineCheckResult
	scanScene := req.ScanScene
	total, err := model.CountDistinctTasks(ctx, scanScene)
	if err != nil {
		return nil, err
	}
	rows, err := model.ListGroupedByTask(ctx, page, size, scanScene)
	if err != nil {
		return nil, err
	}
	resp := &typespec.BaselineTaskListResp{Total: total}
	for _, r := range rows {
		resp.List = append(resp.List, typespec.BaselineTaskListItem{
			TaskID:        r.TaskID,
			TargetIP:      r.TargetIP,
			OSType:        r.OSType,
			OSTypeName:    enums.BaselineEnum.GetOSTypeName(r.OSType),
			ScanScene:     r.ScanScene,
			ScanSceneName: enums.BaselineEnum.GetHostScanSceneName(r.ScanScene),
			TotalRules:    int(r.TotalRules),
			PassCount:     int(r.PassCount),
			FailCount:     int(r.FailCount),
			ErrorCount:    int(r.ErrCount),
			CheckTime:     r.LastTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) ReloadBaselineRulesFromDB(ctx context.Context) error {
	return services.ReloadBaselineRulesFromDB(ctx)
}

// ImportBaselineRules 导入规则到引擎并写入数据库
func (a *BaselineApp) ImportBaselineRules(ctx context.Context, req *typespec.BaselineRulesImportReq) *typespec.BaselineRulesImportResp {
	// 先写入数据库，再从数据库重载到引擎，保证引擎与DB始终同步
	var model mysqls.HostBaselineRule

	// 查询当前最大 rule_code，用于生成唯一编号
	var maxRuleCode int
	mysql.FromContext(ctx).Model(&model).Select("COALESCE(MAX(rule_code), 0)").Scan(&maxRuleCode)

	// 构建已有规则唯一键集合（name+category+osType），用于数据库去重
	dbExisting := make(map[string]bool)
	var allRules []mysqls.HostBaselineRule
	if err := mysql.FromContext(ctx).Model(&model).Select("name, category, os_type").Find(&allRules).Error; err == nil {
		for _, r := range allRules {
			key := fmt.Sprintf("%s|%d|%d", r.Name, r.Category, r.OSType)
			dbExisting[key] = true
		}
	}

	dbSuccess := 0
	dbSkipped := 0
	nextCode := maxRuleCode
	for _, item := range req.Rules {
		key := fmt.Sprintf("%s|%d|%d", item.Name, item.Category, item.OSType)
		if dbExisting[key] {
			dbSkipped++
			continue
		}
		cmdsJSON := "[]"
		if len(item.Commands) > 0 {
			b, _ := json.Marshal(item.Commands)
			cmdsJSON = string(b)
		}
		nextCode++
		rule := &mysqls.HostBaselineRule{
			RuleCode:        nextCode,
			Name:            item.Name,
			Description:     item.Description,
			Category:        item.Category,
			Risk:            item.Risk,
			OSType:          item.OSType,
			CommandsJSON:    cmdsJSON,
			ExpectedValue:   item.ExpectedValue,
			MatchType:       item.MatchType,
			FixSuggestion:   item.FixSuggestion,
			RiskDescription: item.RiskDescription,
			Enabled:         1,
		}
		if err := model.Create(ctx, rule); err != nil {
			log.Errorf("ImportBaselineRules: 写入数据库失败 name=%s err=%v", item.Name, err)
			dbSkipped++
		} else {
			dbExisting[key] = true
			dbSuccess++
		}
	}

	// 从数据库重载到引擎，保证引擎与DB完全一致
	if err := services.ReloadBaselineRulesFromDB(ctx); err != nil {
		log.Errorf("ImportBaselineRules: 重载引擎失败 err=%v", err)
	}

	return &typespec.BaselineRulesImportResp{
		Total:   len(req.Rules),
		Success: dbSuccess,
		Skipped: dbSkipped,
	}
}

// GetBaselineRulesFromDB 从数据库获取规则列表
func (a *BaselineApp) GetBaselineRulesFromDB(ctx context.Context) *typespec.BaselineRulesListResp {
	var model mysqls.HostBaselineRule
	all, err := model.ListAll(ctx)
	if err != nil {
		return &typespec.BaselineRulesListResp{Total: 0}
	}

	osHits := map[int]int{}
	catHits := map[int]int{}
	for _, r := range all {
		osHits[r.OSType]++
		catHits[r.Category]++
	}

	resp := &typespec.BaselineRulesListResp{
		Total: len(all),
		Rules: make([]typespec.BaselineRuleListItem, 0, len(all)),
	}
	osOrder := []int{
		enums.BaselineOSTypeLinux,
		enums.BaselineOSTypeWindows,
		enums.BaselineOSTypeDomestic,
		enums.BaselineOSTypeEmbedded,
	}
	for _, os := range osOrder {
		resp.ByOsType = append(resp.ByOsType, typespec.BaselineRulesCountByOS{
			OSType:     os,
			OSTypeName: enums.BaselineEnum.GetOSTypeName(os),
			Count:      osHits[os],
		})
	}

	catKeys := make([]int, 0, len(catHits))
	for k := range catHits {
		catKeys = append(catKeys, k)
	}
	sort.Ints(catKeys)
	for _, cat := range catKeys {
		resp.ByCategory = append(resp.ByCategory, typespec.BaselineRulesCountByCategory{
			Category:     cat,
			CategoryName: enums.BaselineEnum.GetCategoryName(cat),
			Count:        catHits[cat],
		})
	}

	for _, r := range all {
		var cmds []string
		if r.CommandsJSON != "" {
			json.Unmarshal([]byte(r.CommandsJSON), &cmds)
		}
		resp.Rules = append(resp.Rules, typespec.BaselineRuleListItem{
			ID:              r.ID,
			Name:            r.Name,
			Description:     r.Description,
			Category:        r.Category,
			CategoryName:    enums.BaselineEnum.GetCategoryName(r.Category),
			Risk:            r.Risk,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.Risk),
			OSType:          r.OSType,
			OSTypeName:      enums.BaselineEnum.GetOSTypeName(r.OSType),
			ExpectedValue:   r.ExpectedValue,
			MatchType:       r.MatchType,
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			Commands:        cmds,
		})
	}
	return resp
}

// GetBaselineRuleDetail 从数据库获取单条规则详情
func (a *BaselineApp) GetBaselineRuleDetail(ctx context.Context, id int) (*typespec.BaselineRuleDetailResp, error) {
	var model mysqls.HostBaselineRule
	rule, err := model.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	var cmds []string
	if rule.CommandsJSON != "" {
		json.Unmarshal([]byte(rule.CommandsJSON), &cmds)
	}
	return &typespec.BaselineRuleDetailResp{
		ID:              rule.ID,
		RuleCode:        rule.RuleCode,
		Name:            rule.Name,
		Description:     rule.Description,
		Category:        rule.Category,
		CategoryName:    enums.BaselineEnum.GetCategoryName(rule.Category),
		Risk:            rule.Risk,
		RiskName:        enums.BaselineEnum.GetBaselineRiskName(rule.Risk),
		OSType:          rule.OSType,
		OSTypeName:      enums.BaselineEnum.GetOSTypeName(rule.OSType),
		Commands:        cmds,
		ExpectedValue:   rule.ExpectedValue,
		MatchType:       rule.MatchType,
		FixSuggestion:   rule.FixSuggestion,
		RiskDescription: rule.RiskDescription,
		Enabled:         rule.Enabled,
		CreateTime:      rule.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:      rule.UpdateTime.Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateBaselineRule 新增规则到数据库
func (a *BaselineApp) CreateBaselineRule(ctx context.Context, req *typespec.BaselineRuleCreateReq) error {
	cmdsJSON := "[]"
	if len(req.Commands) > 0 {
		b, _ := json.Marshal(req.Commands)
		cmdsJSON = string(b)
	}
	rule := &mysqls.HostBaselineRule{
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		Risk:            req.Risk,
		OSType:          req.OSType,
		CommandsJSON:    cmdsJSON,
		ExpectedValue:   req.ExpectedValue,
		MatchType:       req.MatchType,
		FixSuggestion:   req.FixSuggestion,
		RiskDescription: req.RiskDescription,
		Enabled:         req.Enabled,
	}
	var model mysqls.HostBaselineRule
	return model.Create(ctx, rule)
}

// UpdateBaselineRule 更新数据库中的规则
func (a *BaselineApp) UpdateBaselineRule(ctx context.Context, req *typespec.BaselineRuleUpdateReq) error {
	cmdsJSON := "[]"
	if len(req.Commands) > 0 {
		b, _ := json.Marshal(req.Commands)
		cmdsJSON = string(b)
	}
	rule := &mysqls.HostBaselineRule{
		ID:              req.ID,
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		Risk:            req.Risk,
		OSType:          req.OSType,
		CommandsJSON:    cmdsJSON,
		ExpectedValue:   req.ExpectedValue,
		MatchType:       req.MatchType,
		FixSuggestion:   req.FixSuggestion,
		RiskDescription: req.RiskDescription,
		Enabled:         req.Enabled,
	}
	var model mysqls.HostBaselineRule
	return model.Update(ctx, rule)
}

// DeleteBaselineRule 删除数据库中的规则
func (a *BaselineApp) DeleteBaselineRule(ctx context.Context, id int) error {
	var model mysqls.HostBaselineRule
	return model.Delete(ctx, id)
}

func (a *BaselineApp) GetBaselineRules(ctx context.Context) *typespec.BaselineRulesListResp {
	engine := services.GetBaselineEngine()
	all := engine.GetAllRules()
	sort.Slice(all, func(i, j int) bool {
		if all[i].OSType != all[j].OSType {
			return all[i].OSType < all[j].OSType
		}
		if all[i].Category != all[j].Category {
			return all[i].Category < all[j].Category
		}
		return all[i].ID < all[j].ID
	})

	osHits := map[int]int{}
	catHits := map[int]int{}
	for _, r := range all {
		osHits[r.OSType]++
		catHits[r.Category]++
	}

	resp := &typespec.BaselineRulesListResp{
		Total: len(all),
		Rules: make([]typespec.BaselineRuleListItem, 0, len(all)),
	}
	osOrder := []int{
		enums.BaselineOSTypeLinux,
		enums.BaselineOSTypeWindows,
		enums.BaselineOSTypeDomestic,
		enums.BaselineOSTypeEmbedded,
	}
	for _, os := range osOrder {
		resp.ByOsType = append(resp.ByOsType, typespec.BaselineRulesCountByOS{
			OSType:     os,
			OSTypeName: enums.BaselineEnum.GetOSTypeName(os),
			Count:      osHits[os],
		})
	}

	catKeys := make([]int, 0, len(catHits))
	for k := range catHits {
		catKeys = append(catKeys, k)
	}
	sort.Ints(catKeys)
	for _, cat := range catKeys {
		resp.ByCategory = append(resp.ByCategory, typespec.BaselineRulesCountByCategory{
			Category:     cat,
			CategoryName: enums.BaselineEnum.GetCategoryName(cat),
			Count:        catHits[cat],
		})
	}

	for _, r := range all {
		resp.Rules = append(resp.Rules, typespec.BaselineRuleListItem{
			ID:              r.ID,
			Name:            r.Name,
			Description:     r.Description,
			Category:        r.Category,
			CategoryName:    enums.BaselineEnum.GetCategoryName(r.Category),
			Risk:            r.Risk,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.Risk),
			OSType:          r.OSType,
			OSTypeName:      enums.BaselineEnum.GetOSTypeName(r.OSType),
			ExpectedValue:   r.ExpectedValue,
			MatchType:       r.MatchType,
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			Commands:        r.Commands,
		})
	}
	return resp
}

func (a *BaselineApp) RunMalwareScan(ctx context.Context, req *typespec.MalwareScanReq) (*typespec.MalwareScanResp, error) {
	normalizeMalwareScanReq(req)
	scanner := services.GetMalwareScanner()
	task := &services.MalwareScanTask{
		TaskID:   req.TaskID,
		TargetID: req.TargetID,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Key:      req.Key,
		OSType:   req.OSType,
	}

	report, err := scanner.RunMalwareScan(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := &typespec.MalwareScanResp{
		TaskID:     report.TaskID,
		TargetIP:   report.TargetIP,
		HasMalware: report.HasMalware,
		StartTime:  report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:    report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.MalwareCheckItem{
			ID:            r.ID,
			CheckType:     r.CheckType,
			CheckTypeName: enums.BaselineEnum.GetMalwareCheckTypeName(r.CheckType),
			RiskLevel:     r.RiskLevel,
			RiskName:      enums.BaselineEnum.GetMalwareRiskName(r.RiskLevel),
			MatchRule:     r.MatchRule,
			FilePath:      r.FilePath,
			Description:   r.Description,
			FixSuggestion: r.FixSuggestion,
			CheckTime:     r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (a *BaselineApp) GetMalwareResults(ctx context.Context, req *typespec.MalwareResultListReq) (*typespec.MalwareResultListResp, error) {
	var model mysqls.MalwareCheckResult
	var list []mysqls.MalwareCheckResult
	var err error

	if req.TargetID > 0 {
		list, err = model.GetByTargetID(ctx, req.TargetID)
	} else {
		list, err = model.GetByTaskID(ctx, req.TaskID)
	}
	if err != nil {
		return nil, err
	}

	resp := &typespec.MalwareResultListResp{
		Total: int64(len(list)),
	}
	for _, r := range list {
		resp.List = append(resp.List, typespec.MalwareCheckItem{
			ID:            r.ID,
			CheckType:     r.CheckType,
			CheckTypeName: enums.BaselineEnum.GetMalwareCheckTypeName(r.CheckType),
			RiskLevel:     r.RiskLevel,
			RiskName:      enums.BaselineEnum.GetMalwareRiskName(r.RiskLevel),
			MatchRule:     r.MatchRule,
			FilePath:      r.FilePath,
			Description:   r.Description,
			FixSuggestion: r.FixSuggestion,
			CheckTime:     r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) GetMalwareTaskList(ctx context.Context, req *typespec.MalwareTaskListReq) (*typespec.MalwareTaskListResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 10
	}
	var model mysqls.MalwareCheckResult
	total, err := model.CountDistinctMalwareTasks(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := model.ListMalwareGroupedByTask(ctx, page, size)
	if err != nil {
		return nil, err
	}
	resp := &typespec.MalwareTaskListResp{Total: total}
	for _, r := range rows {
		resp.List = append(resp.List, typespec.MalwareTaskListItem{
			TaskID:         r.TaskID,
			TargetIP:       r.TargetIP,
			TotalFindings:  int(r.TotalFindings),
			WorstRiskLevel: r.WorstRiskLevel,
			WorstRiskName:  enums.BaselineEnum.GetMalwareRiskName(r.WorstRiskLevel),
			CheckTime:      r.LastTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) RunDBCheck(ctx context.Context, req *typespec.DBCheckReq) (*typespec.DBCheckResp, error) {
	checker := services.GetDBBaselineChecker()
	task := &services.DBCheckTask{
		TaskID:   req.TaskID,
		TargetID: req.TargetID,
		Host:     req.Host,
		Port:     req.Port,
		DBType:   req.DBType,
		Username: req.Username,
		Password: req.Password,
		DBName:   req.DBName,
	}

	report, err := checker.RunDBCheck(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := &typespec.DBCheckResp{
		TaskID:     report.TaskID,
		TargetIP:   report.TargetIP,
		DBType:     report.DBType,
		DBTypeName: enums.BaselineEnum.GetDBTypeName(report.DBType),
		TotalRules: report.TotalRules,
		PassCount:  report.PassCount,
		FailCount:  report.FailCount,
		StartTime:  report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:    report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.DBCheckItem{
			ID:              r.ID,
			RuleID:          r.RuleID,
			RuleName:        r.RuleName,
			CheckCategory:   r.CheckCategory,
			CategoryName:    enums.BaselineEnum.GetDBCheckCategoryName(r.CheckCategory),
			CheckResult:     r.CheckResult,
			ResultName:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			ExpectedValue:   r.ExpectedValue,
			ActualValue:     r.ActualValue,
			RiskLevel:       r.RiskLevel,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.RiskLevel),
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			CheckTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (a *BaselineApp) GetDBCheckResults(ctx context.Context, req *typespec.DBCheckResultListReq) (*typespec.DBCheckResultListResp, error) {
	var model mysqls.DBCheckResult
	var list []mysqls.DBCheckResult
	var err error

	if req.TargetID > 0 {
		list, err = model.GetByTargetID(ctx, req.TargetID)
	} else {
		list, err = model.GetByTaskID(ctx, req.TaskID)
	}
	if err != nil {
		return nil, err
	}

	resp := &typespec.DBCheckResultListResp{
		Total: int64(len(list)),
	}
	for _, r := range list {
		resp.List = append(resp.List, typespec.DBCheckItem{
			ID:              r.ID,
			RuleID:          r.RuleID,
			RuleName:        r.RuleName,
			CheckCategory:   r.CheckCategory,
			CategoryName:    enums.BaselineEnum.GetDBCheckCategoryName(r.CheckCategory),
			CheckResult:     r.CheckResult,
			ResultName:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			ExpectedValue:   r.ExpectedValue,
			ActualValue:     r.ActualValue,
			RiskLevel:       r.RiskLevel,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.RiskLevel),
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			CheckTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) RunSensitiveDataScan(ctx context.Context, req *typespec.SensitiveDataScanReq) (*typespec.SensitiveDataScanResp, error) {
	finder := services.GetSensitiveDataFinder()
	task := &services.SensitiveDataTask{
		TaskID:    req.TaskID,
		TargetID:  req.TargetID,
		Host:      req.Host,
		Port:      req.Port,
		DBType:    req.DBType,
		Username:  req.Username,
		Password:  req.Password,
		DBName:    req.DBName,
		ScanAllDB: req.ScanAllDB,
	}

	report, err := finder.RunScan(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := &typespec.SensitiveDataScanResp{
		TaskID:      report.TaskID,
		TargetIP:    report.TargetIP,
		DBType:      report.DBType,
		DBTypeName:  enums.BaselineEnum.GetDBTypeName(report.DBType),
		HighCount:   report.HighCount,
		MiddleCount: report.MiddleCount,
		LowCount:    report.LowCount,
		StartTime:   report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:     report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.SensitiveDataItem{
			ID:            r.ID,
			DBName:        r.DBName,
			TableName:     r.TableNameStr,
			ColumnName:    r.ColumnName,
			DataType:      r.DataType,
			DataTypeName:  enums.BaselineEnum.GetSensitiveDataTypeName(r.DataType),
			DataLevel:     r.DataLevel,
			DataLevelName: enums.BaselineEnum.GetSensitiveDataLevelName(r.DataLevel),
			MatchRule:     r.MatchRule,
			SampleData:    r.SampleData,
			CreateTime:    r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (a *BaselineApp) GetSensitiveDataResults(ctx context.Context, req *typespec.SensitiveDataListReq) (*typespec.SensitiveDataListResp, error) {
	var model mysqls.SensitiveDataResult
	var list []mysqls.SensitiveDataResult
	var err error

	if req.TargetID > 0 {
		list, err = model.GetByTargetID(ctx, req.TargetID)
	} else {
		list, err = model.GetByTaskID(ctx, req.TaskID)
	}
	if err != nil {
		return nil, err
	}

	resp := &typespec.SensitiveDataListResp{
		Total: int64(len(list)),
	}
	for _, r := range list {
		resp.List = append(resp.List, typespec.SensitiveDataItem{
			ID:            r.ID,
			DBName:        r.DBName,
			TableName:     r.TableNameStr,
			ColumnName:    r.ColumnName,
			DataType:      r.DataType,
			DataTypeName:  enums.BaselineEnum.GetSensitiveDataTypeName(r.DataType),
			DataLevel:     r.DataLevel,
			DataLevelName: enums.BaselineEnum.GetSensitiveDataLevelName(r.DataLevel),
			MatchRule:     r.MatchRule,
			SampleData:    r.SampleData,
			CreateTime:    r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) GetSensitiveDataStat(ctx context.Context, req *typespec.SensitiveDataStatReq) (*typespec.SensitiveDataStatResp, error) {
	var model mysqls.SensitiveDataResult
	high, middle, low, total, err := model.GetStatByTaskID(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}
	return &typespec.SensitiveDataStatResp{
		HighCount:   int(high),
		MiddleCount: int(middle),
		LowCount:    int(low),
		TotalCount:  int(total),
	}, nil
}

func (a *BaselineApp) GetEnums(ctx context.Context) interface{} {
	return map[string]interface{}{
		"osTypes": []map[string]interface{}{
			{"value": enums.BaselineOSTypeLinux, "label": enums.BaselineEnum.GetOSTypeName(enums.BaselineOSTypeLinux)},
			{"value": enums.BaselineOSTypeWindows, "label": enums.BaselineEnum.GetOSTypeName(enums.BaselineOSTypeWindows)},
			{"value": enums.BaselineOSTypeDomestic, "label": enums.BaselineEnum.GetOSTypeName(enums.BaselineOSTypeDomestic)},
			{"value": enums.BaselineOSTypeEmbedded, "label": enums.BaselineEnum.GetOSTypeName(enums.BaselineOSTypeEmbedded)},
		},
		"hostTransports": []map[string]interface{}{
			{"value": enums.HostTransportAuto, "label": enums.BaselineEnum.GetHostTransportName(enums.HostTransportAuto)},
			{"value": enums.HostTransportSSH, "label": enums.BaselineEnum.GetHostTransportName(enums.HostTransportSSH)},
			{"value": enums.HostTransportWinRM, "label": enums.BaselineEnum.GetHostTransportName(enums.HostTransportWinRM)},
		},
		"dbTypes": []map[string]interface{}{
			{"value": enums.DBSupportTypeMySQL, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypeMySQL)},
			{"value": enums.DBSupportTypePostgreSQL, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypePostgreSQL)},
			{"value": enums.DBSupportTypeMongoDB, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypeMongoDB)},
			{"value": enums.DBSupportTypeRedis, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypeRedis)},
			{"value": enums.DBSupportTypeCouchDB, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypeCouchDB)},
		},
		"checkCategories": []map[string]interface{}{
			{"value": enums.BaselineCategoryPasswordPolicy, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryPasswordPolicy)},
			{"value": enums.BaselineCategoryUserPermission, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryUserPermission)},
			{"value": enums.BaselineCategoryFirewall, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryFirewall)},
			{"value": enums.BaselineCategoryKernelSecurity, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryKernelSecurity)},
			{"value": enums.BaselineCategoryAuditLog, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryAuditLog)},
			{"value": enums.BaselineCategoryNetworkService, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryNetworkService)},
			{"value": enums.BaselineCategorySSHConfig, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategorySSHConfig)},
		},
		"dbCheckCategories": []map[string]interface{}{
			{"value": enums.DBCheckCategoryAuthentication, "label": enums.BaselineEnum.GetDBCheckCategoryName(enums.DBCheckCategoryAuthentication)},
			{"value": enums.DBCheckCategoryAuthorization, "label": enums.BaselineEnum.GetDBCheckCategoryName(enums.DBCheckCategoryAuthorization)},
			{"value": enums.DBCheckCategoryConfigSecure, "label": enums.BaselineEnum.GetDBCheckCategoryName(enums.DBCheckCategoryConfigSecure)},
			{"value": enums.DBCheckCategoryAuditLog, "label": enums.BaselineEnum.GetDBCheckCategoryName(enums.DBCheckCategoryAuditLog)},
		},
		"sensitiveDataLevels": []map[string]interface{}{
			{"value": enums.SensitiveDataLevelHigh, "label": enums.BaselineEnum.GetSensitiveDataLevelName(enums.SensitiveDataLevelHigh)},
			{"value": enums.SensitiveDataLevelMiddle, "label": enums.BaselineEnum.GetSensitiveDataLevelName(enums.SensitiveDataLevelMiddle)},
			{"value": enums.SensitiveDataLevelLow, "label": enums.BaselineEnum.GetSensitiveDataLevelName(enums.SensitiveDataLevelLow)},
		},
		"checkResults": []map[string]interface{}{
			{"value": enums.BaselineCheckResultPass, "label": enums.BaselineEnum.GetCheckResultName(enums.BaselineCheckResultPass)},
			{"value": enums.BaselineCheckResultFail, "label": enums.BaselineEnum.GetCheckResultName(enums.BaselineCheckResultFail)},
			{"value": enums.BaselineCheckResultError, "label": enums.BaselineEnum.GetCheckResultName(enums.BaselineCheckResultError)},
		},
	}
}
