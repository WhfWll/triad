package crons

import (
	"context"
	"encoding/json"
	"fmt"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"

	log "github.com/sirupsen/logrus"
)

type TaskGroupOverView struct {
	Name         string
	Risk         int
	CreateTime   string
	Describe     string
	TaskStat     []int
	TargetStat   []int
	VulStat      []int
	ExpStat      []int
	EvidenceStat []int
	VulTypeStat  map[string]interface{}
	TaskTrend    map[string]interface{}
}

func TaskGroupInfoStat() {
	var (
		taskGroupSrv       services.TaskGroup
		taskModel          mysqls.TaskTask
		targetModel        mysqls.TaskTarget
		taskVulModel       mysqls.TaskVul
		taskEvidenceModel  mysqls.TaskEvidence
		taskGroupTaskModel mysqls.TaskGroupTask
		IsFinished         bool = true
	)
	ctx := context.Background()
	// 第一步 获取待统计的任务
	groups, err := taskGroupSrv.GetGroupByIsStat(ctx, enums.TaskGroupIsStatNo)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, group := range groups {
		taskIdList, err := taskGroupTaskModel.GetTaskIdByGroupId(ctx, group.ID)
		if err != nil {
			fmt.Println(err)
		}
		overview := TaskGroupOverView{
			Name:       group.Name,
			Describe:   group.Describe,
			CreateTime: group.CreateTime.String(),
		}
		// 第二步 对 组内的任务 进行统计
		tempUserIdList := make([]int, 0)
		taskList, _, err := taskModel.GetTaskCheckTaskList(ctx, 0, 100000, 0, "", "", "", taskIdList, tempUserIdList)
		var highTask, middleTask, lowTask, safeTask int
		for _, task := range taskList {
			switch task.RiskLevel {
			case enums.TaskRiskHigh:
				highTask += 1
			case enums.TaskRiskMid:
				middleTask += 1
			case enums.TaskRiskLow:
				lowTask += 1
			case enums.TaskRiskSafe:
				safeTask += 1
			}
			if task.Status != enums.TaskStatusFinish && task.Status != enums.TaskStatusPausing {
				IsFinished = false
			}
		}
		if IsFinished {
			err := group.UpdateTaskGroupStatus(ctx, enums.TaskStatusFinish, enums.TaskGroupIsStatYes)
			if err != nil {
				log.Error(err)
			}
		}
		taskStat := []int{highTask, middleTask, lowTask, safeTask}
		if highTask > 0 {
			overview.Risk = enums.TaskRiskHigh
		} else if middleTask > 0 {
			overview.Risk = enums.TaskRiskMid
		} else if lowTask > 0 {
			overview.Risk = enums.TaskRiskLow
		} else {
			overview.Risk = enums.TaskRiskSafe
		}
		overview.TaskStat = taskStat

		// 第三步 对 组内的目标 进行统计
		var highTarget, middleTarget, lowTarget, safeTarget int
		targetList := targetModel.GetTargetsByTaskIds(ctx, taskIdList)
		for _, target := range targetList {
			switch target.RiskLevel {
			case enums.TargetRiskHigh:
				highTarget += 1
			case enums.TargetRiskMid:
				middleTarget += 1
			case enums.TargetRiskLow:
				lowTarget += 1
			case enums.TargetRiskLowNoFound:
				safeTarget += 1
			}
		}
		targetStat := []int{highTarget, middleTarget, lowTarget, safeTarget}
		overview.TargetStat = targetStat

		// 第四步 对 组内的漏洞 进行统计
		var deadVul, highVul, middleVul, lowVul int
		var notVerify, verifySuccess, expSuccess int
		vulTypeMap := make(map[int]int, 0)
		taskVulList := taskVulModel.GetTaskVulByTaskIds(ctx, taskIdList)
		for _, vul := range taskVulList {
			switch vul.Risk {
			case enums.VulLibrariesRiskDead:
				deadVul += 1
			case enums.VulLibrariesRiskHigh:
				highVul += 1
			case enums.VulLibrariesRiskMiddle:
				middleVul += 1
			case enums.VulLibrariesRiskLow:
				lowVul += 1
			}

			switch vul.Status {
			case enums.VulStatusNotVerify:
				notVerify += 1
			case enums.VulStatusVerifySuccess:
				verifySuccess += 1
			case enums.VulStatusVerifyUsed:
				expSuccess += 1
			}

			vulTypeMap[vul.Type] += 1
		}
		vulStat := []int{deadVul, highVul, middleVul, lowVul}
		expStat := []int{expSuccess, verifySuccess, notVerify}
		vulTypeStatX := make([]string, 0)
		vulTypeStatY := make([]int, 0)
		for vtype, number := range vulTypeMap {
			vulTypeStatX = append(vulTypeStatX, enums.ToolsVulnerabilityEnum.GetTypeEnum(vtype))
			vulTypeStatY = append(vulTypeStatY, number)
		}
		overview.VulStat = vulStat
		overview.ExpStat = expStat
		overview.VulTypeStat = map[string]interface{}{
			"x": vulTypeStatX,
			"y": vulTypeStatY,
		}

		// 第五步 对任务组内 证据做下统计
		var remoteControl, dataLeak, loginInfo int
		taskEvidenceList := taskEvidenceModel.GetTaskEvidenceByTaskIds(ctx, taskIdList)
		for _, evidence := range taskEvidenceList {
			switch evidence.RiskType {
			case enums.VulScriptEvidenceTypeCommandExec:
				remoteControl += 1
			case enums.VulScriptEvidenceTypeWeakPass:
				loginInfo += 1
			case enums.VulScriptEvidenceTypeInfoLeak:
				dataLeak += 1
			case enums.VulScriptEvidenceTypeFileLeak:
				dataLeak += 1
			case enums.VulScriptEvidenceTypeData:
				dataLeak += 1
			}
		}
		evidenceStat := []int{remoteControl, dataLeak, loginInfo}
		overview.EvidenceStat = evidenceStat

		// 第六步 对任务组内 任务趋势进行统计
		dateList, startTime, dateFormat := utils.MonthDateList()
		taskTrendStatList := taskModel.GetTaskTrendStat(ctx, startTime, dateFormat, 0, 0)
		taskTrendStatMap := make(map[string]int)
		taskTrendY := make([]int, 0)
		for _, item := range taskTrendStatList {
			taskTrendStatMap[item.Date] = item.Count
		}
		for _, item := range dateList {
			taskTrendY = append(taskTrendY, taskTrendStatMap[item])
		}
		overview.TaskTrend = map[string]interface{}{
			"x": dateList,
			"y": taskTrendY,
		}

		overviewByte, err := json.Marshal(overview)
		err = group.UpdateTaskGroupOverview(ctx, string(overviewByte))
		if err != nil {
			log.Println("updateTaskGroupOverview error: ", err.Error())
		}
	}

}
