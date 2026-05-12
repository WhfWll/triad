package crons

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	aesEncryption "smart/tools/encryption"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
)

// ReportGenerate 生成报告
func ReportGenerate() {
	ctx := context.Background()
	// 报告分类
	var reportSrv services.Report
	taskCate, targetCate, err := reportSrv.GetReportContentEnum(ctx)
	if err != nil {
		log.Error("crontab ReportGenerate 获取报告分类 err：" + err.Error())
		return
	}

	reportRecodes := reportSrv.GetsByStatus(ctx, enums.ReportStatusWait)
	for _, reportRecode := range reportRecodes {
		// 优先更新报告状态
		_ = reportSrv.UpdateStatus(ctx, reportRecode.ID, enums.ReportStatusRunning)
		// 转换报告ID为string 方便下面日志记录
		reportRecodeId := strconv.Itoa(reportRecode.ID)

		// 报告配置 - 配置详情
		configJson := make(map[string]interface{})
		err = json.Unmarshal([]byte(reportRecode.ConfigJSON), &configJson)
		if err != nil {
			log.Error("crontab ReportGenerate " + reportRecodeId + " 解析配置json err：" + err.Error())
			continue
		}
		var (
			objId  float64
			objIDs string
		)
		// 考虑多个目标生成文档
		objIdMid := configJson["objId"]
		if _, ok := objIdMid.(string); ok {
			if !strings.Contains(objIdMid.(string), ",") {
				objId, _ = strconv.ParseFloat(objIdMid.(string), 64)
			} else {
				// 就是多个报告合并输出
				objIDs = objIdMid.(string)
			}
		} else {
			objId = objIdMid.(float64)
		}
		if objId == 0 && objIDs == "" {
			log.Error("crontab ReportGenerate " + reportRecodeId + " 报告配置中未检测到objId")
			continue
		}
		configJsonContentMap := make(map[string]int)
		if configJsonContent, ok := configJson["content"]; ok {
			configJsonContentByte, err := json.Marshal(configJsonContent)
			if err != nil {
				log.Error("crontab ReportGenerate " + reportRecodeId + " 解析报告配置中的content err：" + err.Error())
				continue
			}
			err = json.Unmarshal(configJsonContentByte, &configJsonContentMap)
			if err != nil {
				log.Error("crontab ReportGenerate " + reportRecodeId + " 反向解析报告配置中的content err：" + err.Error())
				continue
			}
		}
		// 根据类型 处理任务报告/目标报告
		switch reportRecode.Type {
		case enums.ReportTypeTask: // 任务报告
			// 任务ID
			if objIDs != "" {
				ids := strings.Split(objIDs, ",")
				var list []typespec.ReportTaskContent
				for _, v := range ids {
					taskId, _ := strconv.Atoi(v)
					info, err := ExecReportTypeTask(ctx, taskId, reportRecodeId, reportRecode, taskCate, configJsonContentMap)
					if err != nil {
						continue
					}
					list = append(list, info)
				}
				taskContentByte, _ := json.Marshal(list)
				if err = reportSrv.ReportUpdateContent(ctx, reportRecode.ID, string(taskContentByte)); err != nil {
					log.Error("crontab ReportGenerate " + reportRecodeId + " 任务报告生成结果 err：" + err.Error())
					continue
				}
			} else {
				taskContent, err := ExecReportTypeTask(ctx, int(objId), reportRecodeId, reportRecode, taskCate, configJsonContentMap)
				if err != nil {
					continue
				}
				taskContentByte, _ := json.Marshal(taskContent)
				if err = reportSrv.ReportUpdateContent(ctx, reportRecode.ID, string(taskContentByte)); err != nil {
					log.Error("crontab ReportGenerate " + reportRecodeId + " 任务报告生成结果 err：" + err.Error())
					continue
				}
			}
		case enums.ReportTypeTarget: // 目标报告
			var taskContentByte []byte
			// 合并输出
			if objIDs != "" {
				ids := strings.Split(objIDs, ",")
				var list []typespec.ReportTargetContent
				for _, v := range ids {
					targetId, _ := strconv.Atoi(v)
					info, err := ExecReportTypeTarget(ctx, targetId, reportRecodeId, reportRecode, targetCate, configJsonContentMap)
					if err != nil {
						continue
					}
					list = append(list, info)
				}
				taskContentByte, _ = json.Marshal(list)
			} else {
				// 目标ID
				info, err := ExecReportTypeTarget(ctx, int(objId), reportRecodeId, reportRecode, targetCate, configJsonContentMap)
				if err != nil {
					continue
				}
				taskContentByte, _ = json.Marshal(info)
			}
			err = reportSrv.ReportUpdateContent(ctx, reportRecode.ID, string(taskContentByte))
			if err != nil {
				log.Error("crontab ReportGenerate " + reportRecodeId + " 目标报告生成结果 err：" + err.Error())
				continue
			}
		case enums.ReportTypeLogicTask: // 逻辑漏洞综述报告
			// 任务ID
			if objIDs != "" {
				ids := strings.Split(objIDs, ",")
				var list []typespec.ReportTaskContent
				for _, v := range ids {
					taskId, _ := strconv.Atoi(v)
					info, err := ExecReportTypeTask(ctx, taskId, reportRecodeId, reportRecode, taskCate, configJsonContentMap)
					if err != nil {
						continue
					}
					list = append(list, info)
				}
				taskContentByte, _ := json.Marshal(list)
				if err = reportSrv.ReportUpdateContent(ctx, reportRecode.ID, string(taskContentByte)); err != nil {
					log.Error("crontab ReportGenerate " + reportRecodeId + " 任务报告生成结果 err：" + err.Error())
					continue
				}
			} else {
				taskContent, err := ExecReportLogicTask(ctx, int(objId), reportRecodeId, reportRecode, taskCate, configJsonContentMap)
				if err != nil {
					continue
				}
				taskContentByte, _ := json.Marshal(taskContent)
				if err = reportSrv.ReportUpdateContent(ctx, reportRecode.ID, string(taskContentByte)); err != nil {
					log.Error("crontab ReportGenerate " + reportRecodeId + " 任务报告生成结果 err：" + err.Error())
					continue
				}
			}
		case enums.ReportTypeLogicTarget: // 逻辑漏洞综述报告
			var taskContentByte []byte
			// 合并输出
			if objIDs != "" {
				ids := strings.Split(objIDs, ",")
				var list []typespec.ReportTargetContent
				for _, v := range ids {
					targetId, _ := strconv.Atoi(v)
					info, err := ExecReportLogicTarget(ctx, targetId, reportRecodeId, reportRecode, targetCate, configJsonContentMap)
					//info, err := ExecReportTypeTarget(ctx, targetId, reportRecodeId, reportRecode, targetCate, configJsonContentMap)
					if err != nil {
						continue
					}
					list = append(list, info)
				}
				taskContentByte, _ = json.Marshal(list)
			} else {
				// 目标ID
				info, err := ExecReportLogicTarget(ctx, int(objId), reportRecodeId, reportRecode, targetCate, configJsonContentMap)
				if err != nil {
					continue
				}
				taskContentByte, _ = json.Marshal(info)
			}
			err = reportSrv.ReportUpdateContent(ctx, reportRecode.ID, string(taskContentByte))
			if err != nil {
				log.Error("crontab ReportGenerate " + reportRecodeId + " 目标报告生成结果 err：" + err.Error())
				continue
			}
		}
	}
	return
}

// ExecReportTypeTask 执行处理任务报告
func ExecReportTypeTask(ctx context.Context, taskId int, reportRecodeId string, reportRecode mysqls.Reportrecord, taskCate []services.ReportContentItem, configJsonContentMap map[string]int) (typespec.ReportTaskContent, error) {
	// 任务
	var (
		taskSrv     services.TaskTask
		targetSrv   services.TaskTarget
		taskVulSrv  services.TaskVul
		taskContent typespec.ReportTaskContent
	)
	task, err := taskSrv.GetTaskByTaskId(ctx, taskId)
	if err != nil {
		log.Error("crontab ReportGenerate " + reportRecodeId + " 任务报告获取任务数据 err：" + err.Error())
		return typespec.ReportTaskContent{}, err
	}
	if task.ID == 0 {
		log.Error("crontab ReportGenerate " + reportRecodeId + " 任务报告未知的任务数据")
		return typespec.ReportTaskContent{}, err
	}
	// 任务下所有目标
	targets, _ := targetSrv.GetTargetByTaskId(ctx, taskId)
	// 任务下所有漏洞数据
	taskVuls := taskVulSrv.GetsByTaskId(ctx, taskId, enums.VulDataTypOne)
	taskContent.ReportId = reportRecode.ID
	// 封面
	taskContent.ReportCover = Report.ReportGenerateCover(reportRecode)
	// 导航
	taskContent.Catalog = Report.ReportGenerateCate(taskCate, configJsonContentMap)
	// 任务概述
	taskContent.TaskOverview = Report.ReportGenerateTaskOverview(ctx, &task, &targets, &taskVuls)
	// 信息统计 - 目标风险统计
	targetState := taskContent.TaskOverview.TargetStat
	taskContent.TargetRisk = Report.ReportGenerateTaskInfoTarget(targetState.HighTarget, targetState.MiddleTarget, targetState.LowTarget, targetState.SafeTarget)
	// 信息统计 - 漏洞风险统计
	taskContent.VulRisk = Report.ReportGenerateTaskInfoVulRIsk(&taskVuls)
	// 信息统计 - 漏洞类型统计
	taskContent.VulType = Report.ReportGenerateTaskInfoVulType(&taskVuls)
	// 信息统计 - Top危险漏洞
	taskContent.TopVulRisk = Report.ReportGenerateTaskInfoTop(ctx, &taskVuls)
	// 目标风险
	taskContent.TargetDetails = Report.ReportGenerateTaskTargetRisk(&targets, &taskVuls)
	// 漏洞详情
	taskContent.VulDetails = Report.ReportGenerateTaskVulDetail(&targets, &taskVuls)
	return taskContent, nil
}

// ExecReportTypeTarget 执行处理目标报告
func ExecReportTypeTarget(ctx context.Context, targetId int, reportRecodeId string, reportRecode mysqls.Reportrecord, targetCate []services.ReportContentItem, configJsonContentMap map[string]int) (typespec.ReportTargetContent, error) {
	// 查询目标
	var (
		targetSrv    services.TaskTarget
		taskVulSrv   services.TaskVul
		targetConten typespec.ReportTargetContent
	)
	target, err := targetSrv.GetTargetById(ctx, targetId)
	if err != nil {
		log.Error("crontab ReportGenerate " + reportRecodeId + " 目标报告获取目标数据 err：" + err.Error())
		return typespec.ReportTargetContent{}, err
	}
	if target.ID == 0 {
		log.Error("crontab ReportGenerate " + reportRecodeId + " 目标报告未知的目标数据")
		return typespec.ReportTargetContent{}, err
	}
	// 目标下所有漏洞数据
	taskVuls := taskVulSrv.GetsByTargetId(ctx, targetId, enums.VulDataTypOne)
	targetConten.ReportId = reportRecode.ID
	// 封面
	targetConten.ReportCover = Report.ReportGenerateCover(reportRecode)
	// 导航
	targetConten.Catalog = Report.ReportGenerateCate(targetCate, configJsonContentMap)
	// 报告摘要
	targetConten.TargetOverview = Report.ReportGenerateTargetOverview(&target, &taskVuls)
	// 资产信息
	targetConten.AssetInfo = Report.ReportGenerateTargetAssetInfo(ctx, target.OpSys, target.ID)
	// 漏洞信息
	targetConten.VulInfo = Report.ReportGenerateTargetVulDetail(&[]mysqls.TaskTarget{target}, &taskVuls)
	return targetConten, nil
}

var Report report

type report struct {
}

var LogicReport logicReport

type logicReport struct {
}

// 报告生成 - 首页背景
func (r *report) ReportGenerateCover(reportRecode mysqls.Reportrecord) (res typespec.ReportCoverNode) {
	res.Title = reportRecode.Name
	res.CreateTime = reportRecode.CreateTime.Format(utils.DateOnly)
	res.BackgroundImg = ""
	return res
}

// 报告生成 - 分类
func (r *report) ReportGenerateCate(configJsonContent []services.ReportContentItem, configJsonContentMap map[string]int) (returnData []typespec.CatalogNode) {
	fmt.Println(configJsonContentMap)
	firstTag := 0
	for _, item := range configJsonContent {
		fmt.Println(item.Value)
		if configJsonContentMap[item.Value] == 1 {
			firstTag++
			nextTag := 0
			var catalogNode typespec.CatalogNode
			catalogNode.Id = item.Value
			catalogNode.Name = strconv.Itoa(firstTag) + ". " + item.Label
			catalogNode.IsShow = true
			for _, item2 := range item.Items {
				if configJsonContentMap[item2.Value] == 1 {
					nextTag++
					var catalogNode2 typespec.CatalogNode
					catalogNode2.Id = item2.Value
					catalogNode2.Name = strconv.Itoa(firstTag) + "." + strconv.Itoa(nextTag) + " " + item2.Label
					catalogNode2.IsShow = true
					catalogNode.Catalog = append(catalogNode.Catalog, catalogNode2)
				}
			}
			returnData = append(returnData, catalogNode)
		}
	}
	return
}

// 报告生成 - 任务报告 - 任务概述
func (r *report) ReportGenerateTaskOverview(ctx context.Context, task *mysqls.TaskTask, targets *[]mysqls.TaskTarget, taskVul *[]mysqls.TaskVul) (returnData typespec.TaskOverview) {

	// 任务数据
	returnData.TaskName = task.TaskName
	returnData.Date = task.CreateTime.Format("2006-01-02 15:04") + "至" + task.UpdateTime.Format("2006-01-02 15:04")
	returnData.TaskRiskStr = enums.TaskTaskEnum.RiskEnum(task.RiskLevel)
	returnData.TargetStat.Total = task.TargetNum
	returnData.TargetStat.HighTarget = task.HigeNum
	returnData.TargetStat.MiddleTarget = task.MiddleNum
	returnData.TargetStat.LowTarget = task.LowNum
	returnData.TargetStat.SafeTarget = task.SafeNum

	// 获取场景名称
	var templateSrv services.SceneTaskTemplate
	template, err := templateSrv.GetTaskTemplateById(ctx, task.TaskTemplateID)
	if err == nil {
		returnData.TemplateName = template.TemplateName
	}

	// 获取目标数据
	for _, target := range *targets {
		if target.IsAlive == enums.TargetIsAliveY {
			returnData.TargetStat.LiveTarget += 1
		}
	}

	// 获取漏洞数据
	var (
		vulNumArray [6]int //每个等级的数量，元素含义分别为：无漏洞个数/致命漏洞个数/高危漏洞个数/中危漏洞个数/低危漏洞个数/信息漏洞个数
		vulTotal    int
	)
	for _, vul := range *taskVul {
		switch vul.Status {
		case enums.VulStatusVerifySuccess:
			returnData.VulnVerify.VerifySuccess++
		case enums.VulStatusVerifyUsed:
			returnData.VulnVerify.UseSuccess++
		case enums.VulStatusNotVerify:
			returnData.VulnVerify.RepairSuccess++
		}
		if vul.Status == enums.VulStatusRepairSuccess || vul.Risk == enums.VulLibrariesRiskNot { //已经修复或risk为0的算安全
			vulNumArray[5] += 1
			continue
		}
		vulNumArray[vul.Risk] += 1
	}
	returnData.VulnStat.Total = vulTotal
	returnData.VulnStat.DeadlyNumber = vulNumArray[1]
	returnData.VulnStat.HighNumber = vulNumArray[2]
	returnData.VulnStat.MiddleNumber = vulNumArray[3]
	returnData.VulnStat.LowNumber = vulNumArray[4]

	return
}

// 报告生成 - 任务报告 - 信息统计 - 目标风险统计
func (r *report) ReportGenerateTaskInfoTarget(high, middle, low, safe int) (returnData typespec.TargetRisk) {
	returnData.Total = high + middle + low + safe

	returnData.HighNumber = high
	highNumberRate, _ := utils.MathPercent(returnData.Total, returnData.HighNumber)
	returnData.HighNumberRate = highNumberRate

	returnData.MiddleNumber = middle
	middleNumberRate, _ := utils.MathPercent(returnData.Total, returnData.MiddleNumber)
	returnData.MiddleNumberRate = middleNumberRate

	returnData.LowNumber = low
	lowNumberRate, _ := utils.MathPercent(returnData.Total, returnData.LowNumber)
	returnData.LowNumberRate = lowNumberRate

	returnData.SafeNumber = safe
	safeNumberRate, _ := utils.MathPercent(returnData.Total, returnData.SafeNumber)
	returnData.SafeNumberRate = safeNumberRate
	return
}

// 报告生成 - 任务报告 - 信息统计 - 漏洞风险统计
func (r *report) ReportGenerateTaskInfoVulRIsk(taskVul *[]mysqls.TaskVul) (returnData []typespec.VulRisk) {
	var deadVul typespec.VulRisk
	deadVul.RiskType = "致命漏洞"

	var highVul typespec.VulRisk
	highVul.RiskType = "高危漏洞"

	var middleVul typespec.VulRisk
	middleVul.RiskType = "中危漏洞"

	var lowVul typespec.VulRisk
	lowVul.RiskType = "低危漏洞"

	var total int
	for _, vul := range *taskVul {
		switch vul.Risk {
		case enums.VulLibrariesRiskDead: // 致命
			total++
			deadVul.Total++
			switch vul.Status {
			case enums.VulStatusNotVerify: // 未验证
				deadVul.RepairSuccess++
			case enums.VulStatusVerifySuccess: // 验证成功
				deadVul.VerifySuccess++
			case enums.VulStatusVerifyUsed: // 利用成功
				deadVul.UseSuccess++
			}
		case enums.VulLibrariesRiskHigh: // 高危
			total++
			highVul.Total++
			switch vul.Status {
			case enums.VulStatusNotVerify: // 未验证
				highVul.RepairSuccess++
			case enums.VulStatusVerifySuccess: // 验证成功
				highVul.VerifySuccess++
			case enums.VulStatusVerifyUsed: // 利用成功
				highVul.UseSuccess++
			}
		case enums.VulLibrariesRiskMiddle: // 中危
			total++
			middleVul.Total++
			switch vul.Status {
			case enums.VulStatusNotVerify: // 未验证
				middleVul.RepairSuccess++
			case enums.VulStatusVerifySuccess: // 验证成功
				middleVul.VerifySuccess++
			case enums.VulStatusVerifyUsed: // 利用成功
				middleVul.UseSuccess++
			}
		case enums.VulLibrariesRiskLow: // 低危
			total++
			lowVul.Total++
			switch vul.Status {
			case enums.VulStatusNotVerify: // 未验证
				lowVul.RepairSuccess++
			case enums.VulStatusVerifySuccess: // 验证成功
				lowVul.VerifySuccess++
			case enums.VulStatusVerifyUsed: // 利用成功
				lowVul.UseSuccess++
			}
		}
	}

	returnData = append(returnData, deadVul)
	returnData = append(returnData, highVul)
	returnData = append(returnData, middleVul)
	returnData = append(returnData, lowVul)

	// 计算统计与每种类型下的百分比
	countVerifySuccess := 0
	countrepairVerify := 0
	countuseSuccess := 0
	countTotal := 0
	for k, item := range returnData {
		// 每个漏洞分类的占比
		numRate, _ := utils.MathPercent(total, item.Total)
		returnData[k].Percent = numRate

		// 统计信息
		countVerifySuccess += item.VerifySuccess
		countrepairVerify += item.RepairSuccess
		countuseSuccess += item.UseSuccess
		countTotal += item.Total
	}

	// 补充统计信息至列表
	returnData = append(returnData, typespec.VulRisk{
		RiskType:      "统计",
		VerifySuccess: countVerifySuccess,
		RepairSuccess: countrepairVerify,
		UseSuccess:    countuseSuccess,
		Total:         countTotal,
		Percent:       "100%",
	})
	return
}

// 报告生成 - 任务报告 - 信息统计 - 漏洞类型统计
func (r *report) ReportGenerateTaskInfoVulType(taskVul *[]mysqls.TaskVul) (returnData []typespec.VulType) {
	type tempStruct struct {
		VulnType  string      `json:"vulnType"`  // 任务报告--漏洞类型--漏洞类型
		Total     int         `json:"total"`     // 任务报告--漏洞类型--数量
		TargetIds map[int]int `json:"targetIds"` // 任务报告--漏洞类型--目标数量
	}
	var totalType int

	vulTypeMap := make(map[int]tempStruct)
	for _, item := range *taskVul {
		totalType++
		if vulType, ok := vulTypeMap[item.Type]; ok {
			// 总数++
			vulType.Total++
			// 影响目标数量
			if _, okTarget := vulType.TargetIds[item.TargetID]; okTarget {
				vulType.TargetIds[item.TargetID]++
			} else {
				vulType.TargetIds[item.TargetID] = 1
			}
			vulTypeMap[item.Type] = vulType
		} else {
			var temp tempStruct
			temp.VulnType = enums.ToolsVulnerabilityEnum.GetTypeEnum(item.Type)
			temp.Total = 1
			temp.TargetIds = make(map[int]int)
			temp.TargetIds[item.TargetID] = 1
			vulTypeMap[item.Type] = temp
		}
	}

	// 转换数据
	for _, item := range vulTypeMap {
		rate, _ := utils.MathPercent(totalType, item.Total)

		returnData = append(returnData, typespec.VulType{
			VulnType:     item.VulnType,
			Total:        item.Total,
			Percent:      rate,
			TargetNumber: len(item.TargetIds),
		})
	}
	// 按数量从高到低排序
	for i := 0; i < len(returnData)-1; i++ {
		for j := 0; j < len(returnData)-1-i; j++ {
			if returnData[j].Total < returnData[j+1].Total {
				returnData[j], returnData[j+1] = returnData[j+1], returnData[j]
			}
		}
	}
	return
}

// 报告生成 - 任务报告 - 信息统计 - Top危险漏洞
func (r *report) ReportGenerateTaskInfoTop(ctx context.Context, taskVul *[]mysqls.TaskVul) (returnData []typespec.TopVulRisk) {
	type tempStruct struct {
		Name      string
		Risk      int
		Num       int
		TargetIds map[int]int
	}

	// 数据库查询时已按风险等级排序，无需再次排序

	// 只取前10个漏洞，重复不算，但出现次数需要++
	num := 10
	margeData := make(map[string]tempStruct)
	targetIds := make([]int, 0)

	// 销售许可证临时添加
	key := []byte("9876787656785679")
	aesEcb := aesEncryption.AesEcb{}

	for _, item := range *taskVul {
		targetIds = append(targetIds, item.TargetID)
		if marge, ok := margeData[item.DecisionVulId]; ok {
			marge.Num++
			marge.TargetIds[item.TargetID] = item.TargetID
			margeData[item.DecisionVulId] = marge
		} else {
			// 首次出现数量++
			num--

			var temp tempStruct
			//temp.Name = item.Name
			nameDecodeByte, _ := hex.DecodeString(item.Name)
			temp.Name = string(aesEcb.AesDecryptECB(nameDecodeByte, key))
			if temp.Name == "" {
				temp.Name = item.Name
			}

			temp.Risk = item.Risk
			temp.Num = 1
			temp.TargetIds = make(map[int]int)
			temp.TargetIds[item.TargetID] = item.TargetID
			margeData[item.DecisionVulId] = temp
		}

		if num == 0 {
			break
		}
	}

	// 获取目标地址
	var targetSrv services.TaskTarget
	targets := targetSrv.GetTargetByIds(ctx, targetIds, 0)

	// 将数据分组并转换为切片以便排序
	tempSlice := make([]tempStruct, 0, len(margeData))
	for _, item := range margeData {
		tempSlice = append(tempSlice, item)
	}

	// 按风险等级排序：致命(1) > 高危(2) > 中危(3) > 低危(4)
	for i := 0; i < len(tempSlice)-1; i++ {
		for j := 0; j < len(tempSlice)-1-i; j++ {
			if tempSlice[j].Risk > tempSlice[j+1].Risk {
				tempSlice[j], tempSlice[j+1] = tempSlice[j+1], tempSlice[j]
			}
		}
	}

	// 将排序后的数据转换为返回格式
	for _, item := range tempSlice {
		var temp typespec.TopVulRisk
		temp.VulName = item.Name
		temp.Risk = enums.ToolsVulnerabilityEnum.GetRiskEnum(item.Risk)
		temp.Number = item.Num

		affectTargets := make([]string, 0)
		for _, targetId := range item.TargetIds {
			if target, ok := targets[targetId]; ok {
				affectTargets = append(affectTargets, target.TargetURL)
			}
		}
		temp.AffectTargets = strings.Join(affectTargets, ",")
		returnData = append(returnData, temp)
	}
	return
}

// 报告生成 - 任务报告 - 目标风险
func (r *report) ReportGenerateTaskTargetRisk(targets *[]mysqls.TaskTarget, taskVul *[]mysqls.TaskVul) (returnData []typespec.TargetDetails) {
	// 先将所有漏洞下的目标状态做排序，并取等级最高的1个 等级计算方式：利用成功 > 验证成功 > 未验证
	sortData := make(map[int]int) // {targetId:status}
	for _, vul := range *taskVul {
		// 过滤状态大于利用成功的
		if vul.Status > enums.VulStatusVerifyUsed {
			continue
		}
		if status, ok := sortData[vul.TargetID]; ok {
			if status < vul.Status {
				sortData[vul.TargetID] = vul.Status
			}
		} else {
			sortData[vul.TargetID] = vul.Status
		}
	}
	var targetIds = make([]int, 0)
	for i := 0; i < len(*targets); i++ {
		targetIds = append(targetIds, (*targets)[i].ID)
	}
	var (
		taskVulSrv services.TaskVul
		ctx        = context.Background()
	)

	_, _, vulNumArrayMap, err := taskVulSrv.GetTargetStatsBytargetIds(ctx, targetIds)
	if err != nil {
		return
	}
	for _, item := range *targets {
		// 过滤掉不存活的目标
		if item.IsAlive == enums.TargetIsAliveN {
			continue
		}
		var (
			tempfatal  = 0
			temphigh   = 0
			tempmiddle = 0
			templow    = 0
		)
		if v, ok := vulNumArrayMap[item.ID]; ok {
			tempfatal = v[1]
			temphigh = v[2]
			tempmiddle = v[3]
			templow = v[4]
		}
		vulStatus := enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(sortData[item.ID])
		if vulStatus == "" {
			vulStatus = "未发现漏洞"
		}
		returnData = append(returnData, typespec.TargetDetails{
			Target: item.TargetURL,
			Risk:   enums.GetTargetRisk(item.RiskLevel),

			DeadlyNumber: tempfatal,
			HighNumber:   temphigh,
			MiddleNumber: tempmiddle,
			LowNumber:    templow,

			VulStatus: vulStatus,
		})

	}
	return
}

// 报告生成 - 任务报告 - 漏洞详情
func (r *report) ReportGenerateTaskVulDetail(targets *[]mysqls.TaskTarget, taskVul *[]mysqls.TaskVul) (returnData []typespec.VulDetails) {
	targetMaps := make(map[int]string)
	for _, item := range *targets {
		targetMaps[item.ID] = item.TargetURL
	}

	// 销售许可证临时添加
	key := []byte("9876787656785679")
	aesEcb := aesEncryption.AesEcb{}

	// 使用map进行漏洞去重和统计
	type tempStruct struct {
		Num        int
		targetUrls map[string]string
		Status     int
		Data       mysqls.TaskVul
	}
	vulMaps := make(map[string]tempStruct) // {DecisionVulId:{num:0,targetUrls:{},status:0,data:{}}}

	for _, item := range *taskVul {
		if vul, ok := vulMaps[item.DecisionVulId]; ok {
			// 漏洞已存在，增加出现次数
			vul.Num++
			// 添加影响目标
			vul.targetUrls[targetMaps[item.TargetID]] = ""
			// 更新状态（取最高状态）
			if item.Status <= enums.VulStatusVerifyUsed && item.Status > vul.Status {
				vul.Status = item.Status
			}
			vulMaps[item.DecisionVulId] = vul
		} else {
			// 新漏洞
			var temp tempStruct
			temp.Num = 1
			temp.Data = item
			temp.Status = item.Status
			temp.targetUrls = make(map[string]string)
			temp.targetUrls[targetMaps[item.TargetID]] = ""
			vulMaps[item.DecisionVulId] = temp
		}
	}

	// 先将 vulMaps 转换为切片，以便排序
	type vulMapItem struct {
		key string
		vul tempStruct
	}

	vulSlice := make([]vulMapItem, 0, len(vulMaps))
	for key, vul := range vulMaps {
		vulSlice = append(vulSlice, vulMapItem{key: key, vul: vul})
	}

	// 按照风险等级排序：致命(1) > 高危(2) > 中危(3) > 低危(4)
	// 使用冒泡排序按风险等级排序
	for i := 0; i < len(vulSlice)-1; i++ {
		for j := 0; j < len(vulSlice)-1-i; j++ {
			if vulSlice[j].vul.Data.Risk > vulSlice[j+1].vul.Data.Risk {
				vulSlice[j], vulSlice[j+1] = vulSlice[j+1], vulSlice[j]
			}
		}
	}

	// 转换数据返回
	for _, item := range vulSlice {
		vul := item.vul
		// 构建影响目标列表
		affectTargets := make([]string, 0)
		for targetUrl, _ := range vul.targetUrls {
			if targetUrl != "" { // 过滤空的目标URL
				affectTargets = append(affectTargets, targetUrl)
			}
		}

		cveCnvdCnnvd := strings.Split(vul.Data.VulID, ",")
		cve := cveCnvdCnnvd[0]
		if cve == "" {
			cve = "—"
		}

		// 处理公开日期
		publishDate := vul.Data.PublishedTime
		if publishDate == "" {
			publishDate = "—"
		}

		// 处理影响范围
		affectRange := vul.Data.AffectRange
		if affectRange == "" {
			affectRange = "—"
		}

		// 处理参考链接
		refLink := vul.Data.RefUrl
		if refLink == "" {
			refLink = "—"
		}

		// 处理 Name
		vulName := vul.Data.Name
		if utils.IsHexString(vul.Data.Name) {
			nameDecodeByte, _ := hex.DecodeString(vul.Data.Name)
			vulName = string(aesEcb.AesDecryptECB(nameDecodeByte, key))
		}

		// 处理 VulResult
		vulResult := vul.Data.VulResult
		if utils.IsHexString(vul.Data.VulResult) {
			resultDecodeByte, _ := hex.DecodeString(vul.Data.VulResult)
			vulResult = string(aesEcb.AesDecryptECB(resultDecodeByte, key))
		}

		// 处理 Location
		location := vul.Data.Location
		if utils.IsHexString(vul.Data.Location) {
			locationDecodeByte, _ := hex.DecodeString(vul.Data.Location)
			location = string(aesEcb.AesDecryptECB(locationDecodeByte, key))
		}
		// HTML特殊字符处理
		location = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(location, `<`, `&lt;`), `>`, `&gt;`), `prompt(1)`, `Prompt(1)`)

		returnData = append(returnData, typespec.VulDetails{
			VulName:       vulName,
			Risk:          enums.ToolsVulnerabilityEnum.GetRiskEnum(vul.Data.Risk),
			VulStatus:     enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(vul.Status),
			Number:        vul.Num,
			Type:          enums.ToolsVulnerabilityEnum.GetTypeEnum(vul.Data.Type),
			Cve:           cve,
			PublishDate:   publishDate,
			Describe:      vul.Data.Description,
			Res:           vulResult,
			Location:      location,
			Fix:           vul.Data.FixSuggest,
			AffectRange:   affectRange,
			AffectTargets: strings.Join(affectTargets, ","),
			Link:          refLink,
			VerMsg: typespec.VerMsg{
				Request:  "",
				Response: "",
				Payload:  "",
			},
		})
	}
	return
}

// 报告生成 - 目标报告 - 漏洞详情
func (r *report) ReportGenerateTargetVulDetail(targets *[]mysqls.TaskTarget, taskVul *[]mysqls.TaskVul) (returnData []typespec.VulDetails) {
	targetMaps := make(map[int]string)
	for _, item := range *targets {
		targetMaps[item.ID] = item.TargetURL
	}

	// 销售许可证临时添加
	key := []byte("9876787656785679")
	aesEcb := aesEncryption.AesEcb{}

	// 最多显示20个漏洞的详细报文
	maxVerMsgCount := 20
	currentVerMsgCount := 0
	// 单个报文最大长度限制（字符数）
	maxMsgLength := 1500

	// 转换数据返回
	for _, vul := range *taskVul {
		cveCnvdCnnvd := strings.Split(vul.VulID, ",")
		cve := cveCnvdCnnvd[0]
		if cve == "" {
			cve = "-"
		}

		// 处理公开日期
		publishDate := vul.PublishedTime
		if publishDate == "" {
			publishDate = "-"
		}

		// 处理影响范围
		affectRange := vul.AffectRange
		if affectRange == "" {
			affectRange = "-"
		}

		// 处理参考链接
		refLink := vul.RefUrl
		if refLink == "" {
			refLink = "-"
		}

		// 处理 Name
		vulName := vul.Name
		if utils.IsHexString(vul.Name) {
			nameDecodeByte, _ := hex.DecodeString(vul.Name)
			vulName = string(aesEcb.AesDecryptECB(nameDecodeByte, key))
		}

		// 处理 VulResult
		vulResult := vul.VulResult
		if utils.IsHexString(vul.VulResult) {
			resultDecodeByte, _ := hex.DecodeString(vul.VulResult)
			vulResult = string(aesEcb.AesDecryptECB(resultDecodeByte, key))
		}

		// 处理 Location
		location := vul.Location
		if utils.IsHexString(vul.Location) {
			locationDecodeByte, _ := hex.DecodeString(vul.Location)
			location = string(aesEcb.AesDecryptECB(locationDecodeByte, key))
		}
		// HTML特殊字符处理
		location = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(location, `<`, `&lt;`), `>`, `&gt;`), `prompt(1)`, `Prompt(1)`)

		// 处理 VerMsg - 只显示前20个漏洞的详细报文
		var verMsg []typespec.VerMsg
		verMsgInfo := typespec.VerMsg{}

		// 只显示前20个漏洞的详细报文
		if currentVerMsgCount >= maxVerMsgCount {
			// 超过限制数量的漏洞不返回详细报文，避免报告过大
			verMsgInfo.Request = "报文过长已省略（仅显示前" + fmt.Sprintf("%d", maxVerMsgCount) + "个漏洞的详细报文）"
			verMsgInfo.Response = "报文过长已省略（仅显示前" + fmt.Sprintf("%d", maxVerMsgCount) + "个漏洞的详细报文）"
			verMsgInfo.Payload = ""
			verMsgInfo.PayloadSuccessFlag = ""
		} else {
			// 正常处理报文
			if utils.IsHexString(vul.VerMsg) {
				verMsgByte, _ := hex.DecodeString(vul.VerMsg)
				verMsgStr := string(aesEcb.AesDecryptECB(verMsgByte, key))
				_ = json.Unmarshal([]byte(verMsgStr), &verMsg)
			} else {
				_ = json.Unmarshal([]byte(vul.VerMsg), &verMsg)
			}

			// 只显示第一个验证信息
			if len(verMsg) >= 1 {
				verMsgInfo = verMsg[0]
			}

			// 对报文长度进行限制
			if len(verMsgInfo.Request) > maxMsgLength {
				verMsgInfo.Request = verMsgInfo.Request[:maxMsgLength] + "...(报文过长已截断)"
			}
			if len(verMsgInfo.Response) > maxMsgLength {
				verMsgInfo.Response = verMsgInfo.Response[:maxMsgLength] + "...(报文过长已截断)"
			}
			if len(verMsgInfo.Payload) > maxMsgLength {
				verMsgInfo.Payload = verMsgInfo.Payload[:maxMsgLength] + "...(报文过长已截断)"
			}

			// 增加已处理的报文计数
			currentVerMsgCount++
		}

		returnData = append(returnData, typespec.VulDetails{
			VulName:     vulName,
			Risk:        enums.ToolsVulnerabilityEnum.GetRiskEnum(vul.Risk),
			VulStatus:   enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(vul.Status),
			Type:        enums.ToolsVulnerabilityEnum.GetTypeEnum(vul.Type),
			Cve:         cve,
			PublishDate: publishDate,
			Describe:    vul.Description,
			Res:         vulResult,
			Location:    location,
			Fix:         vul.FixSuggest,
			AffectRange: affectRange,
			Link:        refLink,
			VerMsg: typespec.VerMsg{
				Request:            verMsgInfo.Request,
				Response:           verMsgInfo.Response,
				Payload:            verMsgInfo.Payload,
				PayloadSuccessFlag: verMsgInfo.PayloadSuccessFlag,
			},
		})
	}
	return
}

// 报告生成 - 目标报告 - 报告摘要
func (r *report) ReportGenerateTargetOverview(target *mysqls.TaskTarget, taskVul *[]mysqls.TaskVul) (returnData typespec.TargetOverview) {
	returnData.Target = target.TargetURL
	returnData.Risk = enums.GetTargetRisk(target.RiskLevel)
	var (
		taskVulRes services.TaskVul
		ctx        = context.Background()
	)
	total, _, vulNumArray, err := taskVulRes.GetTargetStats(ctx, target.ID)
	if err != nil {
		return
	}
	returnData.VulnStat.Total = total
	returnData.VulnStat.DeadlyNumber = vulNumArray[1]
	returnData.VulnStat.HighNumber = vulNumArray[2]
	returnData.VulnStat.MiddleNumber = vulNumArray[3]
	returnData.VulnStat.LowNumber = vulNumArray[4]

	// 获取漏洞数据
	for _, vul := range *taskVul {
		switch vul.Status {
		case enums.VulStatusVerifySuccess:
			returnData.VulnVerify.VerifySuccess++
		case enums.VulStatusVerifyUsed:
			returnData.VulnVerify.UseSuccess++
		case enums.VulStatusNotVerify:
			returnData.VulnVerify.RepairSuccess++
		}
	}

	returnData.CreateDate = target.CreateTime.Format(utils.DateOnly) + "至" + target.UpdateTime.Format(utils.DateOnly)
	return
}

// 报告生成 - 目标报告 - 资产信息
func (r *report) ReportGenerateTargetAssetInfo(ctx context.Context, opSys string, targetId int) (returnData typespec.AssetReportInfo) {
	// 资产信息从 task_task_result 表里取数据
	// 其中 组建/指纹条件 obj_type=1(信息) sub_obj_type=1_2(站点) 中的fingerprint
	// 其中 服务条件 obj_type=1(信息) sub_obj_type=1_1(服务) 中的port/protocol
	// 其中 ip/url条件 obj_type=1(信息) sub_obj_type=1_1(服务) 中的ip
	// 其中 系统条件 target表中op_sys字段

	// 查询结果表
	var taskTaskResultSrv services.TaskTaskResult
	// 服务数据
	portProtocol := make([]string, 0)
	ip := ""
	serviceInfos := taskTaskResultSrv.GetTaskInfoByTargetIdAndSubObjType(ctx, targetId, enums.TaskResultSubObjTypeService)
	for _, result := range serviceInfos {
		var service string
		jsonResultMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(result.JSONResult), &jsonResultMap)
		if err == nil && jsonResultMap["service"] != "" {
			service = jsonResultMap["service"].(string)
		}
		if service == "" {
			service = "-"
		}

		portProtocol = append(portProtocol, result.Field2+"/"+service)
		ip = result.Field1
	}

	// 站点数据
	// 指纹数据 有可能重复，所以这里去重
	fingerMap := make(map[string]string)
	finger := make([]string, 0)
	hostInfos := taskTaskResultSrv.GetTaskInfoByTargetIdAndSubObjType(ctx, targetId, enums.TaskResultSubObjTypeSite)
	for _, result := range hostInfos {
		if result.Field3 != "" && fingerMap[result.Field3] == "" {
			finger = append(finger, result.Field3)
			fingerMap[result.Field3] = result.Field3
		}
	}

	returnData.Component = strings.Join(finger, "、")
	returnData.Service = strings.Join(portProtocol, "、")
	returnData.IpOrUrl = ip
	returnData.System = opSys
	return
}

// ExecReportLogicTask 生成未生成的报告
func ExecReportLogicTask(ctx context.Context, taskId int, reportRecodeId string, reportRecode mysqls.Reportrecord, taskCate []services.ReportContentItem, configJsonContentMap map[string]int) (typespec.ReportTaskContent, error) {
	// 任务
	var (
		logicSrv    services.Logic
		taskContent typespec.ReportTaskContent
	)
	task, err := logicSrv.GetTaskById(ctx, taskId)
	if err != nil {
		log.Error("crontab ReportGenerate " + reportRecodeId + " 任务报告获取任务数据 err：" + err.Error())
		return typespec.ReportTaskContent{}, err
	}
	if task.ID == 0 {
		log.Error("crontab ReportGenerate " + reportRecodeId + " 任务报告未知的任务数据")
		return typespec.ReportTaskContent{}, err
	}
	// 任务下所有目标
	targets, _ := logicSrv.GetTargetByTaskId(ctx, taskId)

	// 任务下所有漏洞数据
	taskVuls, _ := logicSrv.GetVulByTaskId(ctx, taskId)
	taskContent.ReportId = reportRecode.ID
	// 封面
	taskContent.ReportCover = LogicReport.ReportGenerateCover(reportRecode)
	// 导航
	taskContent.Catalog = LogicReport.ReportGenerateCate(taskCate, configJsonContentMap)
	// 任务概述
	taskContent.TaskOverview = LogicReport.ReportGenerateTaskOverview(ctx, &task, &targets, &taskVuls)
	// 信息统计 - 目标风险统计
	targetState := taskContent.TaskOverview.TargetStat
	taskContent.TargetRisk = LogicReport.ReportGenerateTaskInfoTarget(targetState.HighTarget, targetState.MiddleTarget, targetState.LowTarget, targetState.SafeTarget)
	// 信息统计 - 漏洞风险统计
	taskContent.VulRisk = LogicReport.ReportGenerateTaskInfoVulRIsk(&taskVuls)
	// 信息统计 - 漏洞类型统计
	taskContent.VulType = LogicReport.ReportGenerateTaskInfoVulType(&taskVuls)
	// 信息统计 - Top危险漏洞
	taskContent.TopVulRisk = LogicReport.ReportGenerateTaskInfoTop(ctx, &taskVuls)
	// 目标风险
	taskContent.TargetDetails = LogicReport.ReportGenerateTaskTargetRisk(&targets, &taskVuls)
	// 漏洞详情
	taskContent.VulDetails = LogicReport.ReportGenerateTaskVulDetail(&targets, &taskVuls)
	return taskContent, nil
}

// ExecReportLogicTarget 执行生成逻辑漏洞目标报告
func ExecReportLogicTarget(ctx context.Context, targetId int, reportRecodeId string, reportRecode mysqls.Reportrecord, targetCate []services.ReportContentItem, configJsonContentMap map[string]int) (typespec.ReportTargetContent, error) {
	var (
		targetConten typespec.ReportTargetContent
		logicSrv     services.Logic
	)
	target, err := logicSrv.GetTargetById(ctx, targetId)
	if err != nil {
		log.Error("crontab ReportGenerate " + reportRecodeId + " 目标报告获取目标数据 err：" + err.Error())
		return typespec.ReportTargetContent{}, err
	}
	if target.ID == 0 {
		log.Error("crontab ReportGenerate " + reportRecodeId + " 目标报告未知的目标数据")
		return typespec.ReportTargetContent{}, err
	}
	// 目标下所有漏洞数据
	taskVuls, _ := logicSrv.GetVulByTargetId(ctx, targetId)

	targetConten.ReportId = reportRecode.ID
	// 封面
	targetConten.ReportCover = LogicReport.ReportGenerateCover(reportRecode)
	// 导航
	targetConten.Catalog = LogicReport.ReportGenerateCate(targetCate, configJsonContentMap)
	// 报告摘要
	targetConten.TargetOverview = LogicReport.ReportGenerateTargetOverview(&target, &taskVuls)
	// 资产信息
	targetConten.AssetInfo = LogicReport.ReportGenerateTargetAssetInfo(ctx, "", target.ID)
	// 漏洞信息
	targetConten.VulInfo = LogicReport.ReportGenerateTargetVulDetail(&[]mysqls.Logictarget{target}, &taskVuls)
	return targetConten, nil
}

// 报告生成 - 首页背景
func (l *logicReport) ReportGenerateCover(reportRecode mysqls.Reportrecord) (res typespec.ReportCoverNode) {
	res.Title = reportRecode.Name
	res.CreateTime = reportRecode.CreateTime.Format(utils.DateOnly)
	res.BackgroundImg = ""
	return res
}

// 报告生成 - 分类
func (l *logicReport) ReportGenerateCate(configJsonContent []services.ReportContentItem, configJsonContentMap map[string]int) (returnData []typespec.CatalogNode) {
	fmt.Println(configJsonContentMap)
	firstTag := 0
	for _, item := range configJsonContent {
		fmt.Println(item.Value)
		if configJsonContentMap[item.Value] == 1 {
			firstTag++
			nextTag := 0
			var catalogNode typespec.CatalogNode
			catalogNode.Id = item.Value
			catalogNode.Name = strconv.Itoa(firstTag) + ". " + item.Label
			catalogNode.IsShow = true
			for _, item2 := range item.Items {
				if configJsonContentMap[item2.Value] == 1 {
					nextTag++
					var catalogNode2 typespec.CatalogNode
					catalogNode2.Id = item2.Value
					catalogNode2.Name = strconv.Itoa(firstTag) + "." + strconv.Itoa(nextTag) + " " + item2.Label
					catalogNode2.IsShow = true
					catalogNode.Catalog = append(catalogNode.Catalog, catalogNode2)
				}
			}
			returnData = append(returnData, catalogNode)
		}
	}
	return
}

// 报告生成 - 任务报告 - 任务概述
func (l *logicReport) ReportGenerateTaskOverview(ctx context.Context, task *mysqls.LogicTask, targets *[]mysqls.Logictarget, taskVul *[]mysqls.Logicvul) (returnData typespec.TaskOverview) {

	// 任务数据
	returnData.TaskName = task.Name
	returnData.Date = task.CreateTime.Format("2006-01-02 15:04") + "至" + task.UpdateTime.Format("2006-01-02 15:04")
	returnData.TaskRiskStr = enums.TaskTaskEnum.RiskEnum(task.Risk)
	returnData.TargetStat.Total = task.TargetNum
	returnData.TargetStat.HighTarget = task.HighNum
	returnData.TargetStat.MiddleTarget = task.MiddleNum
	returnData.TargetStat.LowTarget = task.LowNum
	returnData.TargetStat.SafeTarget = task.SafeNum

	// 获取场景名称
	returnData.TemplateName = "逻辑漏洞内置扫描场景"

	// 获取目标数据
	for _, target := range *targets {
		if target.IsAlive == enums.TargetIsAliveY {
			returnData.TargetStat.LiveTarget += 1
		}
	}

	// 获取漏洞数据
	var (
		vulNumArray [6]int //每个等级的数量，元素含义分别为：无漏洞个数/致命漏洞个数/高危漏洞个数/中危漏洞个数/低危漏洞个数/信息漏洞个数
		vulTotal    int
	)
	for _, vul := range *taskVul {
		switch vul.Status {
		case enums.VulStatusVerifySuccess:
			returnData.VulnVerify.VerifySuccess++
		case enums.VulStatusVerifyUsed:
			returnData.VulnVerify.UseSuccess++
		case enums.VulStatusNotVerify:
			returnData.VulnVerify.RepairSuccess++
		}
		if vul.Status == enums.VulStatusRepairSuccess || vul.Risk == enums.VulLibrariesRiskNot { //已经修复或risk为0的算安全
			vulNumArray[5] += 1
			continue
		}
		vulNumArray[vul.Risk] += 1
	}
	returnData.VulnStat.Total = vulTotal
	returnData.VulnStat.DeadlyNumber = vulNumArray[1]
	returnData.VulnStat.HighNumber = vulNumArray[2]
	returnData.VulnStat.MiddleNumber = vulNumArray[3]
	returnData.VulnStat.LowNumber = vulNumArray[4]

	return
}

// 报告生成 - 任务报告 - 信息统计 - 目标风险统计
func (l *logicReport) ReportGenerateTaskInfoTarget(high, middle, low, safe int) (returnData typespec.TargetRisk) {
	returnData.Total = high + middle + low + safe

	returnData.HighNumber = high
	highNumberRate, _ := utils.MathPercent(returnData.Total, returnData.HighNumber)
	returnData.HighNumberRate = highNumberRate

	returnData.MiddleNumber = middle
	middleNumberRate, _ := utils.MathPercent(returnData.Total, returnData.MiddleNumber)
	returnData.MiddleNumberRate = middleNumberRate

	returnData.LowNumber = low
	lowNumberRate, _ := utils.MathPercent(returnData.Total, returnData.LowNumber)
	returnData.LowNumberRate = lowNumberRate

	returnData.SafeNumber = safe
	safeNumberRate, _ := utils.MathPercent(returnData.Total, returnData.SafeNumber)
	returnData.SafeNumberRate = safeNumberRate
	return
}

// 报告生成 - 任务报告 - 信息统计 - 漏洞风险统计
func (l *logicReport) ReportGenerateTaskInfoVulRIsk(taskVul *[]mysqls.Logicvul) (returnData []typespec.VulRisk) {
	var deadVul typespec.VulRisk
	deadVul.RiskType = "致命漏洞"

	var highVul typespec.VulRisk
	highVul.RiskType = "高危漏洞"

	var middleVul typespec.VulRisk
	middleVul.RiskType = "中危漏洞"

	var lowVul typespec.VulRisk
	lowVul.RiskType = "低危漏洞"

	var total int
	for _, vul := range *taskVul {
		switch vul.Risk {
		case enums.VulLibrariesRiskDead: // 致命
			total++
			deadVul.Total++
			switch vul.Status {
			case enums.VulStatusNotVerify: // 未验证
				deadVul.RepairSuccess++
			case enums.VulStatusVerifySuccess: // 验证成功
				deadVul.VerifySuccess++
			case enums.VulStatusVerifyUsed: // 利用成功
				deadVul.UseSuccess++
			}
		case enums.VulLibrariesRiskHigh: // 高危
			total++
			highVul.Total++
			switch vul.Status {
			case enums.VulStatusNotVerify: // 未验证
				highVul.RepairSuccess++
			case enums.VulStatusVerifySuccess: // 验证成功
				highVul.VerifySuccess++
			case enums.VulStatusVerifyUsed: // 利用成功
				highVul.UseSuccess++
			}
		case enums.VulLibrariesRiskMiddle: // 中危
			total++
			middleVul.Total++
			switch vul.Status {
			case enums.VulStatusNotVerify: // 未验证
				middleVul.RepairSuccess++
			case enums.VulStatusVerifySuccess: // 验证成功
				middleVul.VerifySuccess++
			case enums.VulStatusVerifyUsed: // 利用成功
				middleVul.UseSuccess++
			}
		case enums.VulLibrariesRiskLow: // 低危
			total++
			lowVul.Total++
			switch vul.Status {
			case enums.VulStatusNotVerify: // 未验证
				lowVul.RepairSuccess++
			case enums.VulStatusVerifySuccess: // 验证成功
				lowVul.VerifySuccess++
			case enums.VulStatusVerifyUsed: // 利用成功
				lowVul.UseSuccess++
			}
		}
	}

	returnData = append(returnData, deadVul)
	returnData = append(returnData, highVul)
	returnData = append(returnData, middleVul)
	returnData = append(returnData, lowVul)

	// 计算统计与每种类型下的百分比
	countVerifySuccess := 0
	countrepairVerify := 0
	countuseSuccess := 0
	countTotal := 0
	for k, item := range returnData {
		// 每个漏洞分类的占比
		numRate, _ := utils.MathPercent(total, item.Total)
		returnData[k].Percent = numRate

		// 统计信息
		countVerifySuccess += item.VerifySuccess
		countrepairVerify += item.RepairSuccess
		countuseSuccess += item.UseSuccess
		countTotal += item.Total
	}

	// 补充统计信息至列表
	returnData = append(returnData, typespec.VulRisk{
		RiskType:      "统计",
		VerifySuccess: countVerifySuccess,
		RepairSuccess: countrepairVerify,
		UseSuccess:    countuseSuccess,
		Total:         countTotal,
		Percent:       "100%",
	})
	return
}

// 报告生成 - 任务报告 - 信息统计 - 漏洞类型统计
func (l *logicReport) ReportGenerateTaskInfoVulType(taskVul *[]mysqls.Logicvul) (returnData []typespec.VulType) {
	type tempStruct struct {
		VulnType  string      `json:"vulnType"`  // 任务报告--漏洞类型--漏洞类型
		Total     int         `json:"total"`     // 任务报告--漏洞类型--数量
		TargetIds map[int]int `json:"targetIds"` // 任务报告--漏洞类型--目标数量
	}
	var totalType int

	vulTypeMap := make(map[int]tempStruct)
	for _, item := range *taskVul {
		totalType++
		if vulType, ok := vulTypeMap[item.Type]; ok {
			// 总数++
			vulType.Total++
			// 影响目标数量
			if _, okTarget := vulType.TargetIds[item.TargetID]; okTarget {
				vulType.TargetIds[item.TargetID]++
			} else {
				vulType.TargetIds[item.TargetID] = 1
			}
			vulTypeMap[item.Type] = vulType
		} else {
			var temp tempStruct
			temp.VulnType = enums.Logic{}.AllTypeEnum()[item.Type]
			temp.Total = 1
			temp.TargetIds = make(map[int]int)
			temp.TargetIds[item.TargetID] = 1
			vulTypeMap[item.Type] = temp
		}
	}

	// 转换数据
	for _, item := range vulTypeMap {
		rate, _ := utils.MathPercent(totalType, item.Total)

		returnData = append(returnData, typespec.VulType{
			VulnType:     item.VulnType,
			Total:        item.Total,
			Percent:      rate,
			TargetNumber: len(item.TargetIds),
		})
	}
	return
}

// 报告生成 - 任务报告 - 信息统计 - Top危险漏洞
func (l *logicReport) ReportGenerateTaskInfoTop(ctx context.Context, taskVul *[]mysqls.Logicvul) (returnData []typespec.TopVulRisk) {
	type tempStruct struct {
		Name      string
		Risk      int
		Num       int
		TargetIds map[int]int
	}

	// 数据库查询时已按风险等级排序，无需再次排序

	// 只取前10个漏洞，重复不算，但出现次数需要++
	num := 10
	margeData := make(map[string]tempStruct)
	targetIds := make([]int, 0)
	for _, item := range *taskVul {
		targetIds = append(targetIds, item.TargetID)
		if marge, ok := margeData[item.DecisionVulID]; ok {
			marge.Num++
			marge.TargetIds[item.TargetID] = item.TargetID
			margeData[item.DecisionVulID] = marge
		} else {
			// 首次出现数量++
			num--

			var temp tempStruct
			temp.Name = item.Name
			temp.Risk = item.Risk
			temp.Num = 1
			temp.TargetIds = make(map[int]int)
			temp.TargetIds[item.TargetID] = item.TargetID
			margeData[item.DecisionVulID] = temp
		}

		if num == 0 {
			break
		}
	}

	// 获取目标地址
	var targetSrv services.TaskTarget
	targets := targetSrv.GetTargetByIds(ctx, targetIds, 0)

	// 将数据分组并转换为切片以便排序
	tempSlice := make([]tempStruct, 0, len(margeData))
	for _, item := range margeData {
		tempSlice = append(tempSlice, item)
	}

	// 按风险等级排序：致命(1) > 高危(2) > 中危(3) > 低危(4)
	for i := 0; i < len(tempSlice)-1; i++ {
		for j := 0; j < len(tempSlice)-1-i; j++ {
			if tempSlice[j].Risk > tempSlice[j+1].Risk {
				tempSlice[j], tempSlice[j+1] = tempSlice[j+1], tempSlice[j]
			}
		}
	}

	// 将排序后的数据转换为返回格式
	for _, item := range tempSlice {
		var temp typespec.TopVulRisk
		temp.VulName = item.Name
		temp.Risk = enums.ToolsVulnerabilityEnum.GetRiskEnum(item.Risk)
		temp.Number = item.Num

		affectTargets := make([]string, 0)
		for _, targetId := range item.TargetIds {
			if target, ok := targets[targetId]; ok {
				affectTargets = append(affectTargets, target.TargetURL)
			}
		}
		temp.AffectTargets = strings.Join(affectTargets, ",")
		returnData = append(returnData, temp)
	}
	return
}

// 报告生成 - 任务报告 - 目标风险
func (l *logicReport) ReportGenerateTaskTargetRisk(targets *[]mysqls.Logictarget, taskVul *[]mysqls.Logicvul) (returnData []typespec.TargetDetails) {
	// 先将所有漏洞下的目标状态做排序，并取等级最高的1个 等级计算方式：利用成功 > 验证成功 > 未验证
	sortData := make(map[int]int) // {targetId:status}
	for _, vul := range *taskVul {
		// 过滤状态大于利用成功的
		if vul.Status > enums.VulStatusVerifyUsed {
			continue
		}
		if status, ok := sortData[vul.TargetID]; ok {
			if status < vul.Status {
				sortData[vul.TargetID] = vul.Status
			}
		} else {
			sortData[vul.TargetID] = vul.Status
		}
	}
	var targetIds = make([]int, 0)
	for i := 0; i < len(*targets); i++ {
		targetIds = append(targetIds, (*targets)[i].ID)
	}
	//todo
	var (
		logicSrv services.Logic
		ctx      = context.Background()
	)

	_, _, vulNumArrayMap, err := logicSrv.GetTargetStatsBytargetIds(ctx, targetIds)
	if err != nil {
		return
	}
	for _, item := range *targets {
		var (
			tempfatal  = 0
			temphigh   = 0
			tempmiddle = 0
			templow    = 0
		)
		if v, ok := vulNumArrayMap[item.ID]; ok {
			tempfatal = v[1]
			temphigh = v[2]
			tempmiddle = v[3]
			templow = v[4]
		}
		returnData = append(returnData, typespec.TargetDetails{
			Target: item.TargetURL,
			Risk:   enums.GetTargetRisk(item.Risk),

			DeadlyNumber: tempfatal,
			HighNumber:   temphigh,
			MiddleNumber: tempmiddle,
			LowNumber:    templow,

			VulStatus: enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(sortData[item.ID]),
		})

	}
	return
}

// 报告生成 - 任务报告 - 漏洞详情
func (l *logicReport) ReportGenerateTaskVulDetail(targets *[]mysqls.Logictarget, taskVul *[]mysqls.Logicvul) (returnData []typespec.VulDetails) {
	targetMaps := make(map[int]string)
	for _, item := range *targets {
		targetMaps[item.ID] = item.TargetURL
	}

	type tempStruct struct {
		Num        int
		targetUrls map[string]string
		Status     int
		Data       mysqls.Logicvul
	}
	vulMaps := make(map[string]tempStruct) // {pocname:{num:0,data:{}}}
	for _, item := range *taskVul {
		if vul, ok := vulMaps[item.DecisionVulID]; ok {
			vul.Num++
			vul.targetUrls[targetMaps[item.TargetID]] = ""
			if item.Status <= enums.VulStatusVerifyUsed && item.Status > vul.Status {
				vul.Status = item.Status
			}
			vulMaps[item.DecisionVulID] = vul
		} else {
			var temp tempStruct
			temp.Num = 1
			temp.Data = item
			temp.Status = item.Status
			temp.targetUrls = make(map[string]string)
			temp.targetUrls[targetMaps[item.TargetID]] = ""
			vulMaps[item.DecisionVulID] = temp
		}
	}

	// 转换数据返回
	for _, vul := range vulMaps {
		affectTargets := make([]string, 0)
		for targetUrls, _ := range vul.targetUrls {
			affectTargets = append(affectTargets, targetUrls)
		}
		cveCnvdCnnvd := strings.Split(vul.Data.VulID, ",")
		cve := cveCnvdCnnvd[0]
		if cve == "" {
			cve = "-"
		}
		returnData = append(returnData, typespec.VulDetails{
			VulName:       vul.Data.Name,
			Risk:          enums.ToolsVulnerabilityEnum.GetRiskEnum(vul.Data.Risk),
			VulStatus:     enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(vul.Status),
			Number:        vul.Num,
			Type:          "逻辑漏洞",
			Cve:           cve,
			PublishDate:   "-",
			Describe:      vul.Data.Description,
			Res:           vul.Data.VulResult,
			Location:      vul.Data.Location,
			Fix:           vul.Data.FixSuggest,
			AffectRange:   "-",
			AffectTargets: strings.Join(affectTargets, ","),
			Link:          "-",
		})
	}
	return
}

// 报告生成 - 任务报告 - 漏洞详情
func (l *logicReport) ReportGenerateTargetVulDetail(targets *[]mysqls.Logictarget, taskVul *[]mysqls.Logicvul) (returnData []typespec.VulDetails) {
	for _, vul := range *taskVul {
		cveCnvdCnnvd := strings.Split(vul.VulID, ",")
		cve := cveCnvdCnnvd[0]
		if cve == "" {
			cve = "-"
		}
		returnData = append(returnData, typespec.VulDetails{
			VulName:     vul.Name,
			Risk:        enums.ToolsVulnerabilityEnum.GetRiskEnum(vul.Risk),
			VulStatus:   enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(vul.Status),
			Number:      1,
			Type:        "逻辑漏洞",
			Cve:         cve,
			PublishDate: "-",
			Describe:    vul.Description,
			Res:         vul.VulResult,
			Location:    vul.Location,
			Fix:         vul.FixSuggest,
			AffectRange: "-",
			Link:        "-",
		})
	}
	return
}

// 报告生成 - 目标报告 - 报告摘要
func (l *logicReport) ReportGenerateTargetOverview(target *mysqls.Logictarget, taskVul *[]mysqls.Logicvul) (returnData typespec.TargetOverview) {
	returnData.Target = target.TargetURL
	returnData.Risk = enums.GetTargetRisk(target.Risk)
	var (
		srv services.Logic
		ctx = context.Background()
	)
	total, _, vulNumArray, err := srv.GetTargetStats(ctx, target.ID)
	if err != nil {
		return
	}
	returnData.VulnStat.Total = total
	returnData.VulnStat.DeadlyNumber = vulNumArray[1]
	returnData.VulnStat.HighNumber = vulNumArray[2]
	returnData.VulnStat.MiddleNumber = vulNumArray[3]
	returnData.VulnStat.LowNumber = vulNumArray[4]

	// 获取漏洞数据
	for _, vul := range *taskVul {
		switch vul.Status {
		case enums.VulStatusVerifySuccess:
			returnData.VulnVerify.VerifySuccess++
		case enums.VulStatusVerifyUsed:
			returnData.VulnVerify.UseSuccess++
		case enums.VulStatusNotVerify:
			returnData.VulnVerify.RepairSuccess++
		}
	}

	returnData.CreateDate = target.CreateTime.Format(utils.DateOnly) + "至" + target.UpdateTime.Format(utils.DateOnly)
	return
}

// 报告生成 - 目标报告 - 资产信息
func (l *logicReport) ReportGenerateTargetAssetInfo(ctx context.Context, opSys string, targetId int) (returnData typespec.AssetReportInfo) {
	// 资产信息从 task_task_result 表里取数据
	// 其中 组建/指纹条件 obj_type=1(信息) sub_obj_type=1_2(站点) 中的fingerprint
	// 其中 服务条件 obj_type=1(信息) sub_obj_type=1_1(服务) 中的port/protocol
	// 其中 ip/url条件 obj_type=1(信息) sub_obj_type=1_1(服务) 中的ip
	// 其中 系统条件 target表中op_sys字段

	// 查询结果表
	var taskTaskResultSrv services.TaskTaskResult
	// 服务数据
	portProtocol := make([]string, 0)
	ip := ""
	serviceInfos := taskTaskResultSrv.GetTaskInfoByTargetIdAndSubObjType(ctx, targetId, enums.TaskResultSubObjTypeService)
	for _, result := range serviceInfos {
		var service string
		jsonResultMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(result.JSONResult), &jsonResultMap)
		if err == nil && jsonResultMap["service"] != "" {
			service = jsonResultMap["service"].(string)
		}
		if service == "" {
			service = "-"
		}

		portProtocol = append(portProtocol, result.Field2+"/"+service)
		ip = result.Field1
	}

	// 站点数据
	// 指纹数据 有可能重复，所以这里去重
	fingerMap := make(map[string]string)
	finger := make([]string, 0)
	hostInfos := taskTaskResultSrv.GetTaskInfoByTargetIdAndSubObjType(ctx, targetId, enums.TaskResultSubObjTypeSite)
	for _, result := range hostInfos {
		if result.Field3 != "" && fingerMap[result.Field3] == "" {
			finger = append(finger, result.Field3)
			fingerMap[result.Field3] = result.Field3
		}
	}

	returnData.Component = strings.Join(finger, "、")
	returnData.Service = strings.Join(portProtocol, "、")
	returnData.IpOrUrl = ip
	returnData.System = opSys
	return
}
