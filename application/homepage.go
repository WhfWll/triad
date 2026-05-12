package application

import (
	"context"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"
	"time"
)

type Homepage struct{}

// TaskInfoStat 首页 - 任务信息统计
func (h *Homepage) TaskInfoStat(ctx context.Context, req *typespec.TaskInfoStatReq, res *typespec.TaskInfoStatRes) {
	//处理参数
	var startTime, dateFormat string
	dateList := make([]string, 0)
	switch req.Mode {
	case enums.StatModeWeek:
		dateList, startTime, dateFormat = utils.WeekDateList()
	case enums.StatModeMonth:
		dateList, startTime, dateFormat = utils.MonthDateList()
	case enums.StatModeYear:
		dateList, startTime, dateFormat = utils.YearDateList()
	}

	var taskService services.TaskTask

	uid := 0
	role := 0
	if v := ctx.Value("uid"); v != nil {
		if val, ok := v.(int); ok {
			uid = val
		}
	}
	if v := ctx.Value("role"); v != nil {
		if val, ok := v.(int); ok {
			role = val
		}
	}

	taskTrendStatList := taskService.GetTaskTrendStat(ctx, startTime, dateFormat, uid, role)
	//计算时间为最近一周的开始时间
	latestWeekStartTime := time.Now().AddDate(0, 0, -6).Format(enums.TimeYMDBarLayout)
	taskCount, latestWeekTaskCount := taskService.GetTaskCount(ctx, latestWeekStartTime, uid, role)

	//构造返回数据
	taskTrendStatMap := make(map[string]int)
	for _, item := range taskTrendStatList {
		taskTrendStatMap[item.Date] = item.Count
	}
	for _, item := range dateList {
		res.Count = append(res.Count, taskTrendStatMap[item])
	}
	res.Date = dateList
	res.TaskCount = taskCount
	res.LatestWeekTaskCount = latestWeekTaskCount
}

// TargetRiskStat 目标风险统计
func (h *Homepage) TargetRiskStat(ctx context.Context, res *typespec.TargetRiskStatRes) {
	var targetService services.TaskTarget

	uid := 0
	role := 0
	if v := ctx.Value("uid"); v != nil {
		if val, ok := v.(int); ok {
			uid = val
		}
	}
	if v := ctx.Value("role"); v != nil {
		if val, ok := v.(int); ok {
			role = val
		}
	}

	targetList := targetService.GetTargetRiskStat(ctx, uid, role)
	for _, item := range targetList {
		switch item.RiskLevel {
		case enums.TargetRiskHigh:
			res.HighCount = item.Count
		case enums.TargetRiskMid:
			res.MediumCount = item.Count
		case enums.TargetRiskLow:
			res.LowCount = item.Count
		case enums.TargetRiskLowNoFound:
			res.SafeCount = item.Count
		}
	}
	res.TargetRiskCount = res.HighCount + res.MediumCount + res.LowCount + res.SafeCount
}

// TaskVulRiskStat 首页 - 任务漏洞风险统计
func (h *Homepage) TaskVulRiskStat(ctx context.Context, res *typespec.TaskVulRiskStatRes) {
	var taskVulService services.TaskVul

	uid := 0
	role := 0
	if v := ctx.Value("uid"); v != nil {
		if val, ok := v.(int); ok {
			uid = val
		}
	}
	if v := ctx.Value("role"); v != nil {
		if val, ok := v.(int); ok {
			role = val
		}
	}

	taskVulList := taskVulService.GetTaskVulRiskStat(ctx, uid, role)
	for _, item := range taskVulList {
		switch item.Risk {
		case enums.VulLibrariesRiskDead:
			res.FatalCount = item.Count
		case enums.VulLibrariesRiskHigh:
			res.HighCount = item.Count
		case enums.VulLibrariesRiskMiddle:
			res.MediumCount = item.Count
		case enums.VulLibrariesRiskLow:
			res.LowCount = item.Count
		case enums.VulLibrariesRiskInfo:
			res.InfoCount = item.Count
		}
	}
	res.TaskVulRiskCount = res.FatalCount + res.HighCount + res.MediumCount + res.LowCount
}

// ToolInfoStat 工具信息统计
func (h *Homepage) ToolInfoStat(ctx context.Context, res *typespec.ToolInfoStatRes) error {
	var (
		fingerService     services.Finger
		vulScriptService  services.VulScripts
		sceneService      services.SceneTaskTemplate
		dictionaryService services.Dictionary
	)

	// 获取指纹总数
	fingerCount, err := fingerService.GetFingerCount(ctx)
	if err != nil {
		return err
	}
	res.FingerCount = int(fingerCount)

	// 获取漏洞脚本总数
	vulCount, err := vulScriptService.GetVulScriptCount(ctx)
	if err != nil {
		return err
	}
	res.VulCount = int(vulCount)

	// 获取场景任务总数
	res.TaskSceneCount = sceneService.GetTaskSceneCount(ctx)

	// 获取字典库总数
	dictionaryCount, err := dictionaryService.GetDictionaryCount(ctx)
	if err != nil {
		return err
	}
	res.DictionaryCount = int(dictionaryCount)

	return nil
}

// VulTypeStat 任务漏洞类型统计
func (h *Homepage) VulTypeStat(ctx context.Context, req *typespec.VulTypeStatReq, res *typespec.VulTypeStatRes) error {
	var taskVulService services.TaskVul

	uid := 0
	role := 0
	if v := ctx.Value("uid"); v != nil {
		if val, ok := v.(int); ok {
			uid = val
		}
	}
	if v := ctx.Value("role"); v != nil {
		if val, ok := v.(int); ok {
			role = val
		}
	}

	vulTypeList, vulCountList, err := taskVulService.TaskVulTypeStat(ctx, req.Mode, uid, role)
	res.VulType = vulTypeList
	res.VulCount = vulCountList
	return err
}

// VulFindTrendStat 任务漏洞发现趋势统计
func (h *Homepage) VulFindTrendStat(ctx context.Context, req *typespec.VulFindTrendStatReq, res *typespec.VulFindTrendStatRes) error {
	var taskVulService services.TaskVul

	uid := 0
	role := 0
	if v := ctx.Value("uid"); v != nil {
		if val, ok := v.(int); ok {
			uid = val
		}
	}
	if v := ctx.Value("role"); v != nil {
		if val, ok := v.(int); ok {
			role = val
		}
	}

	dateList, countList, err := taskVulService.TaskVulFindTrendStat(ctx, req.Mode, uid, role)
	res.Date, res.Count = dateList, countList
	return err
}

// VulEvidenceStat 首页 - 漏洞取证
func (h *Homepage) VulEvidenceStat(ctx context.Context, res *typespec.VulEvidenceStatRes) {
	var taskEvidenceService services.TaskEvidence

	uid := 0
	role := 0
	if v := ctx.Value("uid"); v != nil {
		if val, ok := v.(int); ok {
			uid = val
		}
	}
	if v := ctx.Value("role"); v != nil {
		if val, ok := v.(int); ok {
			role = val
		}
	}

	taskEvidenceList := taskEvidenceService.GetVulEvidenceStat(ctx, uid, role)
	for _, item := range taskEvidenceList {
		switch item.RiskType {
		case enums.VulScriptEvidenceTypeWeakPass:
			res.LoginCredentialsCount = item.RiskCount
		case enums.VulScriptEvidenceTypeInfoLeak:
			res.InfoLeakCount = item.RiskCount
		case enums.VulScriptEvidenceTypeFileLeak:
			res.FileLeakCount = item.RiskCount
		case enums.VulScriptEvidenceTypeData:
			res.DbCount = item.RiskCount
		}
	}
	var remoteSessionService services.RemoteSession
	res.RemoteControlCount = remoteSessionService.GetRemoteSessionCount(ctx, uid, role)
}

// MessageStat 最新消息统计模块
func (h *Homepage) MessageStat(ctx context.Context, res *typespec.MessageStatRes) {
	var messageService services.SystemMessage

	uid := 0
	role := 0
	if v := ctx.Value("uid"); v != nil {
		if val, ok := v.(int); ok {
			uid = val
		}
	}
	if v := ctx.Value("role"); v != nil {
		if val, ok := v.(int); ok {
			role = val
		}
	}

	//最新的10条消息
	messageList, _, _ := messageService.GetSystemMessageList(ctx, 1, 5, 0, "", "", "", uid, role)
	for _, item := range messageList {
		res.List = append(res.List, typespec.MessageStatResItem{
			Id:         item.ID,
			Content:    item.Content,
			Type:       item.Type,
			TypeEnum:   enums.MessageTypeMap[item.Type],
			CreateTime: item.CreateTime.Format(enums.TimeLayout),
			Status:     item.Status,
			StatusEnum: enums.MessageStatusMap[item.Status],
			UserId:     item.UserId,
		})
	}
}
