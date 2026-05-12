package application

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"smart/api/typespec"
	"smart/services"
	aesEncryption "smart/tools/encryption"
	"smart/tools/enums"
	tooltime "smart/tools/time"
	"smart/tools/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RiskManageApp struct {
	aesEcb aesEncryption.AesEcb
}

// aesKey 用于解密（可后续抽取配置）
var aesKey = []byte("9876787656785679")

// VulRiskStatistics 漏洞风险统计
func (rms *RiskManageApp) VulRiskStatistics(ctx context.Context, res *typespec.VulRiskStaticsRes) error {
	var taskVulSrv services.RiskManage

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

	vulRes, total, err := taskVulSrv.VulAllDeduplicationList(ctx, 0, uid, role)
	if err != nil {
		return err
	}
	var (
		fixedVulCount int
		fixTimeLength int64
		highRiskCount int
	)
	vulStatusMap := make(map[string]int)
	vulRiskLevelMap := make(map[string]int)
	vulTypeMap := make(map[string]int)
	vulNameMap := make(map[string]*typespec.VulCount)
	for _, v := range vulRes {
		// 漏洞类型
		vulType := v.Type
		vulTypeName := enums.ToolsVulnerabilityEnum.GetTypeEnum(vulType)
		// 漏洞等级
		riskLevel := v.Risk
		riskLevelStr := enums.ToolsVulnerabilityEnum.GetRiskEnum(riskLevel)
		if riskLevel == enums.VulLibrariesRiskHigh || riskLevel == enums.VulLibrariesRiskDead {
			highRiskCount++
		}
		fixStatus := enums.VulUnrepaired
		if v.Status == enums.VulStatusRepairSuccess {
			fixStatus = enums.VulRepaired
			// 漏洞修复时间
			if v.UpdateTime.After(v.CreateTime) {
				duration := v.UpdateTime.Sub(v.CreateTime)
				fixTimeLength = int64(duration.Round(time.Second).Seconds())
			} else {
				fixTimeLength = 0
			}
			fixedVulCount++
		}
		fixStatusStr := enums.ToolsVulnerabilityEnum.GetVulFixStatusEnum(fixStatus)
		vulStatusMap[fixStatusStr]++
		vulRiskLevelMap[riskLevelStr]++
		vulTypeMap[vulTypeName]++
		var decryptedName string
		if utils.IsHexString(v.Name) {
			nameDecodeByte, _ := hex.DecodeString(v.Name)
			decryptedName = string(rms.aesEcb.AesDecryptECB(nameDecodeByte, aesKey))
		} else {
			decryptedName = v.Name
		}
		v.Name = decryptedName
		// 统计 TOP 10 漏洞
		if item, ok := vulNameMap[v.Name]; ok {
			item.Num++
		} else {
			vulNameMap[v.Name] = &typespec.VulCount{
				Name:         v.Name,
				Num:          v.Count,
				RiskLevel:    v.Risk,
				RiskLevelStr: enums.ToolsVulnerabilityEnum.GetRiskEnum(v.Risk),
				CVSS:         v.Cvss,
			}
		}
	}
	if fixedVulCount > 0 {
		averageFixTimeSeconds := fixTimeLength / int64(fixedVulCount)
		res.AverageFixTime = tooltime.FormatDuration(averageFixTimeSeconds)
	} else {
		res.AverageFixTime = "暂无已修复漏洞数据"
	}

	for statusName, count := range vulStatusMap {
		res.StatusStatistics = append(res.StatusStatistics, typespec.VulStatusCount{
			StatusName: statusName,
			Count:      count,
		})
	}
	for riskName, count := range vulRiskLevelMap {
		res.RiskLevelStatistics = append(res.RiskLevelStatistics, typespec.VulRiskLevelCount{
			RiskName: riskName,
			Count:    count,
		})
	}
	sortRiskLevelStatistics(res.RiskLevelStatistics)
	for typeName, count := range vulTypeMap {
		res.TypeStatistics = append(res.TypeStatistics, typespec.VulTypeCount{
			TypeName: typeName,
			Count:    count,
		})
	}
	topVuls := make([]typespec.VulCount, 0, len(vulNameMap))
	for _, v := range vulNameMap {
		topVuls = append(topVuls, *v)
	}
	sort.Slice(topVuls, func(i, j int) bool {
		return topVuls[i].Num > topVuls[j].Num
	})
	if len(topVuls) > 10 {
		res.Top10Vulnerabilities = topVuls[:10]
	} else {
		res.Top10Vulnerabilities = topVuls
	}
	res.TotalVulnerabilities = total
	res.HighRiskCount = highRiskCount
	return nil
}

// VulRiskList 任务详情 - 风险漏洞列表
func (rms *RiskManageApp) VulRiskList(ctx context.Context, req *typespec.VulRiskListReq, resp *typespec.VulRiskListResp) error {
	var taskVulSrv services.RiskManage

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

	vulRes, total, err := taskVulSrv.GetTaskDeduplicationByLocationList(ctx, req.TaskId, req.TargetId, req.Type, req.Risk, req.Search, req.IP, req.Location, req.Page, req.Size, req.Status, req.VerifyType, uid, role)
	if err != nil {
		return err
	}
	resp.Total = total
	resp.List = make([]typespec.VulRiskListRespItems, 0)
	for _, v := range vulRes {
		var vulInfo typespec.VulRiskListRespItems
		vulInfo.ID = v.ID
		vulInfo.Name = v.Name
		vulInfo.Location = v.Location
		// 漏洞类型
		vulInfo.Type = v.Type
		vulInfo.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(v.Type)
		// 漏洞风险等级
		vulInfo.RiskLevel = v.Risk
		vulInfo.RiskLevelStr = enums.ToolsVulnerabilityEnum.GetRiskEnum(v.Risk)
		vulInfo.FixTime = v.UpdateTime.Format(enums.TimeLayout)
		// 漏洞状态
		vulInfo.FixStatus = enums.VulUnrepaired
		if v.Status == enums.VulStatusRepairSuccess {
			vulInfo.FixStatus = enums.VulRepaired
		}
		vulInfo.FixStatusStr = enums.ToolsVulnerabilityEnum.GetVulFixStatusEnum(vulInfo.FixStatus)
		vulInfo.CreateTime = v.CreateTime.Format(enums.TimeLayout)
		vulInfo.FindTime = v.CreateTime.Format(enums.TimeLayout)
		vulInfo.Status = v.Status
		vulInfo.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(v.Status)
		vulInfo.Type = v.Type
		vulInfo.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(v.Type)
		vulInfo.VerifyType = v.TestStatus
		vulInfo.VerifyTypeName = enums.GetTaskVulVerifyTypeEnum(v.TestStatus)
		resp.List = append(resp.List, vulInfo)
	}
	return nil
}

// VulRiskDel 漏洞删除
func (rms *RiskManageApp) VulRiskDel(ctx context.Context, req *typespec.VulRiskDelReq) error {
	var (
		taskVulSrv    services.TaskVul
		taskCheckTask TaskCheckTask
	)
	taskIDs, _ := taskVulSrv.GetVulTaskIDByIDs(ctx, strings.Split(req.IDs, ","))
	for taskID, VulIDs := range taskIDs {
		taskCheckTask.VulDel(ctx, &typespec.VulDelReq{
			TaskVulIds: VulIDs,
			TaskId:     taskID,
		})
	}
	return nil
}

// VulRiskUpdate 漏洞更新
func (rms *RiskManageApp) VulRiskUpdate(ctx context.Context, req *typespec.VulRiskUpdateReq) error {
	var (
		taskVulSrv      services.TaskVul
		vulLifecycleSrv services.VulLifecycle
		taskSrv         services.TaskTaskInfo
		userSrv         services.User
	)
	ids := strings.Split(req.IDs, ",")
	if enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(req.Status) == "" {
		return errors.New("修改状态不合规")
	}
	taskVulMapList, _ := taskVulSrv.GetVulMapByIDs(ctx, ids)
	for _, idStr := range ids {
		idInt, _ := strconv.Atoi(idStr)
		if err := taskVulSrv.UpdateStatusById(ctx, idInt, req.Status); err == nil {
			taskInfo, _ := taskSrv.GetTaskInfoByTaskId(ctx, taskVulMapList[idInt].TaskID)
			userInfo, _ := userSrv.GetUserDetail(ctx, taskInfo.UserID)
			content := vulLifecycleSrv.ConstructStatusChangeContent(userInfo.Username, req.Status)
			vulLifecycleSrv.AddVulLifecycle(ctx, taskVulMapList[idInt].Pocname, taskVulMapList[idInt].Name, taskVulMapList[idInt].Location, content, taskInfo.TaskName, taskVulMapList[idInt].TaskID)
		}

	}
	return nil
}

// VulRiskDetail 漏洞详情
func (rms *RiskManageApp) VulRiskDetail(ctx context.Context, req *typespec.VulRiskDetailReq, resp *typespec.VulRiskInfoResp) error {
	var (
		taskVulSrv          services.TaskVul
		taskLifecycleLogSrv services.VulLifecycle
	)
	taskVulRes, err := taskVulSrv.VulInfo(ctx, req.ID)
	if err != nil {
		return err
	}
	if taskVulRes.ID == 0 {
		return errors.New("找不到该漏洞！")
	}
	resp.Id = taskVulRes.ID
	resp.TaskId = taskVulRes.TaskID
	resp.TargetId = taskVulRes.TargetID
	resp.TargetUrl = taskVulRes.TargetUrl
	resp.Pocname = taskVulRes.Pocname
	resp.Name = taskVulRes.Name
	resp.Type = taskVulRes.Type
	resp.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(taskVulRes.Type)
	resp.Risk = taskVulRes.Risk
	resp.RiskName = enums.ToolsVulnerabilityEnum.GetRiskEnum(taskVulRes.Risk)
	resp.PublishedTime = taskVulRes.PublishedTime
	t, err := time.Parse("2006-01-02 15:04:05", taskVulRes.PublishedTime)
	if err == nil {
		resp.PublishedTime = t.Format("2006-01-02")
	}
	resp.ExploitImpact = taskVulRes.ExploitImpact
	// Description 和 FixSuggest 不需要解密
	resp.Description = taskVulRes.Description
	resp.FixSuggest = taskVulRes.FixSuggest
	resp.RefUrl = taskVulRes.RefUrl
	resp.VulAddress = taskVulRes.VulAddress
	resp.Location = taskVulRes.Location
	resp.PatchUrl = taskVulRes.PatchUrl
	// 处理 VulResult
	resp.VulResult = taskVulRes.VulResult
	resp.VulParam = taskVulRes.VulParam
	// 处理 Location
	resp.Location = taskVulRes.Location
	resp.Status = taskVulRes.Status
	resp.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(taskVulRes.Status)
	resp.VulId = taskVulRes.VulID
	var tmpVerMsg = make([]typespec.VulInfoRespVerMsg, 0)
	json.Unmarshal([]byte(taskVulRes.VerMsg), &tmpVerMsg)
	resp.VerMsg = tmpVerMsg
	// 获取漏洞周期信息
	taskLifeCycleInfo, err := taskLifecycleLogSrv.GetDetail(ctx, taskVulRes.Pocname, taskVulRes.Name, taskVulRes.Location)
	resp.VulLifecycleInfo = make([]typespec.TaskLifecycleLog, 0)
	if err == nil && taskLifeCycleInfo.Content != "" {
		var allLogs []typespec.TaskLifecycleLog
		if err := json.Unmarshal([]byte(taskLifeCycleInfo.Content), &allLogs); err == nil {
			sort.Slice(allLogs, func(i, j int) bool {
				return allLogs[i].TaskID > allLogs[j].TaskID
			})
			resp.VulLifecycleInfo = allLogs
		}
	}
	return nil
}

// VulRiskVulVerify 漏洞风险验证
func (rms *RiskManageApp) VulRiskVulVerify(ctx context.Context, req *typespec.VulRiskVerifyReq) error {
	var (
		taskVulSrv    services.TaskVul
		taskCheckTask TaskCheckTask
	)
	taskVulRes, _ := taskVulSrv.VulInfo(ctx, req.ID)
	return taskCheckTask.VulVerify(ctx, &typespec.VulVerifyReq{
		TaskId:    taskVulRes.TaskID,
		TaskVulId: req.ID,
		Pocname:   taskVulRes.Pocname,
		VulParam:  taskVulRes.VulParam,
	})
}

func contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func getCheckResName(checkRes int) string {
	var result = map[int]string{
		1: "通过",
		2: "不通过",
		3: "错误",
	}
	return result[checkRes]
}

// RiskManageEnums 风险管理枚举
func (rms *RiskManageApp) RiskManageEnums(ctx context.Context, res *typespec.RiskManageEnumsRes) error {
	res.VulVerifyType = enums.GetTaskVulVerifyTypeArray()
	res.VulRiskStatus = enums.ToolsVulnerabilityEnum.AllAssetVulArray()
	res.VulType = enums.GetVulTypeArray()
	res.VulRisk = enums.GetRiskEnumArray()
	return nil
}

// VulRiskListByIP 通过IP检索漏洞信息
func (rms *RiskManageApp) VulRiskListByIP(ctx context.Context, req *typespec.VulRiskListReq, resp *typespec.VulRiskListResp) error {
	var taskVulSrv services.RiskManage

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

	vulRes, total, err := taskVulSrv.GetTaskDeduplicationByLocationList(ctx, 0, 0, 0, 0, "", req.IP, "", req.Page, req.Size, 0, 0, uid, role)
	if err != nil {
		return err
	}
	allList := make([]typespec.VulRiskListRespItems, 0)
	for i := 0; i < len(vulRes); i++ {
		var tmp typespec.VulRiskListRespItems
		tmp.ID = vulRes[i].ID
		tmp.TargetUrl = vulRes[i].TargetUrl
		tmp.Name = vulRes[i].Name
		tmp.Type = vulRes[i].Type
		tmp.Location = vulRes[i].Location
		tmp.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(vulRes[i].Type)
		tmp.RiskLevel = vulRes[i].Risk
		tmp.RiskLevelStr = enums.ToolsVulnerabilityEnum.GetRiskEnum(vulRes[i].Risk)
		tmp.VerifyType = vulRes[i].TestStatus
		tmp.VerifyTypeName = enums.GetTaskVulVerifyTypeEnum(vulRes[i].TestStatus)
		tmp.Status = vulRes[i].Status
		tmp.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(vulRes[i].Status)
		tmp.FindTime = vulRes[i].CreateTime.Format(enums.TimeLayout)
		allList = append(allList, tmp)
	}
	resp.List = allList
	resp.Total = total
	return nil
}

// sortRiskLevelStatistics 对风险等级统计切片进行排序。
// 排序规则：优先按指定的风险等级顺序 (致命 > 高危 > 中危 > 低危 > 信息 > 无) 降序，
// 如果风险等级相同，则按 Count 降序。
// NOTICE: 返回的数据格式-兼容前端展示
func sortRiskLevelStatistics(stats []typespec.VulRiskLevelCount) {
	// 顺序：致命 > 高危 > 中危 > 低危 > 信息 > 无
	var riskPriority = map[string]int{
		"致命": 6,
		"高危": 5,
		"中危": 4,
		"低危": 3,
		"信息": 2,
		"无":  1,
	}
	// 2. 使用 sort.Slice 进行排序
	sort.Slice(stats, func(i, j int) bool {
		itemI := stats[i]
		itemJ := stats[j]
		priorityI := riskPriority[itemI.RiskName]
		priorityJ := riskPriority[itemJ.RiskName]
		if priorityI != priorityJ {
			return priorityI > priorityJ
		}
		return itemI.Count > itemJ.Count
	})
}
