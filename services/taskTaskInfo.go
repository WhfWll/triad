package services

import (
	"context"
	"encoding/json"
	"errors"
	"smart/models/mysqls"
	"smart/models/redises"
	"smart/tools/enums"
	"sort"
	"strconv"
	"strings"
)

// TaskTaskInfo 任务概要信息处理
type TaskTaskInfo struct {
}

// GetTargetStats 对目标进行统计  目标风险统计  目标存活统计  漏洞风险等级统计
func (t *TaskTaskInfo) GetTargetStats(ctx context.Context, taskTask mysqls.TaskTask, overViewStruct *OverView) (bool, error) {
	var (
		taskTarget           mysqls.TaskTarget
		isALLTargetDelFinish = true //是否所有目标都停止或删除
	)
	overViewStruct.Port = make([]*OverViewItems, 0)
	overViewStruct.Service = make([]*OverViewItems, 0)
	overViewStruct.Component = make([]*OverViewItems, 0)
	overViewStruct.OpSys = make([]*OverViewItems, 0)
	overViewStruct.SubDomain = make([]*OverViewItems, 0)
	overViewStruct.UrlTags = make([]*OverViewItems, 0)

	targetList := taskTarget.GetTargetsByTaskId(ctx, taskTask.ID)
	var (
		highTargetNumber   = 0 //高危目标数
		middleTargetNumber = 0 //中危目标数
		lowTargetNumber    = 0 //低危目标数
		safeTargetNumber   = 0 //安全目标数
		allNumber          = 0 //总目标数
		aliveNumber        = 0 //存活目标数
		notAliveNumber     = 0 //不存活目标数
		waitCheckNumber    = 0 //待检测目标数
	)
	for _, target := range targetList {
		if target.Status != enums.TargetStatusFinish {
			isALLTargetDelFinish = false
		}
		allNumber += 1
		//计算风险等级
		if target.RiskLevel == enums.TargetRiskHigh {
			highTargetNumber += 1
		} else if target.RiskLevel == enums.TargetRiskMid {
			middleTargetNumber += 1
		} else if target.RiskLevel == enums.TargetRiskLow {
			lowTargetNumber += 1
		} else if target.RiskLevel == enums.TargetRiskLowNoFound {
			safeTargetNumber += 1
		}
		//计算是否存活
		if target.IsAlive == enums.TargetIsAliveY {
			aliveNumber += 1
		} else if target.IsAlive == enums.TargetIsAliveN && target.Status == enums.TargetStatusFinish {
			notAliveNumber += 1
		} else {
			waitCheckNumber += 1
		}
	}
	overViewStruct.TargetRisk = [4]int{highTargetNumber, middleTargetNumber, lowTargetNumber, safeTargetNumber}
	overViewStruct.TargetNum = []int{aliveNumber, notAliveNumber, waitCheckNumber, allNumber}

	return isALLTargetDelFinish, nil
}

// GetVulStats 对漏洞进行统计  漏洞验证状态统计  漏洞取证统计
//func (t *TaskTaskInfo) GetVulStats(ctx context.Context, taskTask mysqls.TaskTask, overViewStruct *OverView) error {
//	var taskVul TaskVul
//	var (
//		unVerifyNumber = 0                          //未验证漏洞数
//		verifyNumber   = 0                          //验证存在漏洞数
//		useNumber      = 0                          //利用成功漏洞数
//		risklevel      = enums.TargetRiskLowNoFound //目标等级，默认为4->未发现
//		vulNumArray    [6]int                       //每个等级的数量，元素含义分别为：无漏洞个数/致命漏洞个数/高危漏洞个数/中危漏洞个数/低危漏洞个数/信息漏洞个数
//	)
//
//	// 使用聚合查询获取统计数据，避免加载大文本字段
//	statsData := taskVul.GetVulStatsByTaskId(ctx, taskTask.ID, enums.VulDataTypOne)
//
//	// 处理聚合结果
//	for _, stat := range statsData {
//		// 统计验证状态
//		switch stat.Status {
//		case enums.VulStatusNotVerify:
//			unVerifyNumber += stat.Count
//		case enums.VulStatusVerifySuccess:
//			verifyNumber += stat.Count
//		case enums.VulStatusVerifyUsed:
//			useNumber += stat.Count
//		}
//
//		// 统计风险等级
//		if stat.Status == enums.VulStatusRepairSuccess || stat.Risk == enums.VulLibrariesRiskNot {
//			// 已经修复或risk为0的算安全
//			vulNumArray[5] += stat.Count
//			continue
//		}
//
//		// 确保风险等级在有效范围内
//		if stat.Risk >= 0 && stat.Risk < len(vulNumArray) {
//			vulNumArray[stat.Risk] += stat.Count
//		}
//	}
//
//	// 计算整体风险等级
//	if vulNumArray[1] > 0 || vulNumArray[2] > 0 { //致命或高->高
//		risklevel = 1
//	} else if vulNumArray[3] > 0 { //中->中
//		risklevel = 2
//	} else if vulNumArray[4] > 0 { //低->低
//		risklevel = 3
//	}
//
//	// 填充统计结果
//	overViewStruct.VulExploitImpact = append(overViewStruct.VulExploitImpact, &OverViewItems{Label: "未验证", Value: unVerifyNumber})
//	overViewStruct.VulExploitImpact = append(overViewStruct.VulExploitImpact, &OverViewItems{Label: "验证存在", Value: verifyNumber})
//	overViewStruct.VulExploitImpact = append(overViewStruct.VulExploitImpact, &OverViewItems{Label: "利用成功", Value: useNumber})
//
//	overViewStruct.EvidenceStat = append(overViewStruct.EvidenceStat, &OverViewItems{Label: "远程控制", Value: 0})
//	overViewStruct.EvidenceStat = append(overViewStruct.EvidenceStat, &OverViewItems{Label: "数据泄露", Value: 0})
//	overViewStruct.EvidenceStat = append(overViewStruct.EvidenceStat, &OverViewItems{Label: "登录凭证", Value: 0})
//
//	overViewStruct.VulRisk = []int{vulNumArray[1], vulNumArray[2], vulNumArray[3], vulNumArray[4]}
//	overViewStruct.RiskLevel = risklevel
//
//	return nil
//}

// GetVulStats 对漏洞进行统计  漏洞验证状态统计  漏洞取证统计
func (t *TaskTaskInfo) GetVulStats(ctx context.Context, taskTask mysqls.TaskTask, overViewStruct *OverView) error {
	var taskVul mysqls.TaskVul
	var (
		versionMatchNumber = 0                          //低危漏洞数
		vulVerifyNumber    = 0                          //低危漏洞数
		riskLevel          = enums.TargetRiskLowNoFound //目标等级，默认为4->未发现
		vulNumArray        [6]int                       //每个等级的数量，元素含义分别为：无漏洞个数/致命漏洞个数/高危漏洞个数/中危漏洞个数/低危漏洞个数/信息漏洞个数
		unVerifyNumber     = 0                          //未验证漏洞数
		verifyNumber       = 0                          //验证存在漏洞数
		useNumber          = 0                          //利用成功漏洞数
	)
	for _, vul := range taskVul.GetsByTaskId(ctx, taskTask.ID, enums.VulDataTypOne) {
		if vul.TestStatus == enums.VulVerifyTypeVersionMatch {
			versionMatchNumber += 1
		} else {
			vulVerifyNumber += 1
		}
		if vul.Status == enums.VulStatusRepairSuccess || vul.Risk == enums.VulLibrariesRiskNot { //已经修复或risk为0的算安全
			vulNumArray[5] += 1
			continue
		}
		vulNumArray[vul.Risk] += 1
	}
	if vulNumArray[1] > 0 || vulNumArray[2] > 0 { //致命或高->高
		riskLevel = 1
	} else if vulNumArray[3] > 0 { //中->中
		riskLevel = 2
	} else if vulNumArray[4] > 0 { //低->低
		riskLevel = 3
	}

	// 使用聚合查询获取统计数据，避免加载大文本字段
	statsData := taskVul.GetVulStatsByTaskId(ctx, taskTask.ID, enums.VulDataTypOne)

	// 处理聚合结果
	for _, stat := range statsData {
		// 统计验证状态
		switch stat.Status {
		case enums.VulStatusNotVerify:
			unVerifyNumber += stat.Count
		case enums.VulStatusVerifySuccess:
			verifyNumber += stat.Count
		case enums.VulStatusVerifyUsed:
			useNumber += stat.Count
		}
	}

	// 统计漏洞验证类型
	overViewStruct.VerifyType = append(overViewStruct.VerifyType, &OverViewItems{Label: "原理验证", Value: vulVerifyNumber})
	overViewStruct.VerifyType = append(overViewStruct.VerifyType, &OverViewItems{Label: "版本匹配", Value: versionMatchNumber})

	overViewStruct.VulRisk = []int{vulNumArray[1], vulNumArray[2], vulNumArray[3], vulNumArray[4]}
	overViewStruct.RiskLevel = riskLevel

	// 填充统计结果
	overViewStruct.VulExploitImpact = append(overViewStruct.VulExploitImpact, &OverViewItems{Label: "未验证", Value: unVerifyNumber})
	overViewStruct.VulExploitImpact = append(overViewStruct.VulExploitImpact, &OverViewItems{Label: "验证存在", Value: verifyNumber})
	overViewStruct.VulExploitImpact = append(overViewStruct.VulExploitImpact, &OverViewItems{Label: "利用成功", Value: useNumber})

	overViewStruct.EvidenceStat = append(overViewStruct.EvidenceStat, &OverViewItems{Label: "远程控制", Value: 0})
	overViewStruct.EvidenceStat = append(overViewStruct.EvidenceStat, &OverViewItems{Label: "数据泄露", Value: 0})
	overViewStruct.EvidenceStat = append(overViewStruct.EvidenceStat, &OverViewItems{Label: "登录凭证", Value: 0})

	return nil
}

func (t *TaskTaskInfo) GetInfoStats(ctx context.Context, taskTask mysqls.TaskTask, overViewStruct *OverView) error {
	var taskTaskResult mysqls.TaskTaskResult
	resultList, err := taskTaskResult.GetTaskTaskResultListByTaskId(ctx, taskTask.ID)
	if err != nil {
		return errors.New("GetInfoStats error: " + err.Error())
	}
	portMap := make(map[string]int, 0)
	serviceMap := make(map[string]int, 0)
	componentMap := make(map[string]int, 0)
	domainMap := make(map[string]int, 0)
	for _, result := range resultList {
		if result.ObjType != enums.TaskResultObjTypeInfo {
			continue
		}
		//统计 端口 服务 组件
		if result.SubObjType == enums.TaskResultSubObjTypeService {
			if result.Field2 != "" {
				if portMap[result.Field2] == 0 {
					portMap[result.Field2] = 1
				} else {
					portMap[result.Field2] += 1
				}
			}
			if result.Field3 != "" {
				if serviceMap[result.Field3] == 0 {
					serviceMap[result.Field3] = 1
				} else {
					serviceMap[result.Field3] += 1
				}
			}
			if result.Field4 != "" {
				for _, component := range strings.Split(result.Field4, "/") {
					if component == "" {
						continue
					}
					if componentMap[component] == 0 {
						componentMap[component] = 1
					} else {
						componentMap[component] += 1
					}
				}
			}
		}
		//统计url
		if result.SubObjType == enums.TaskResultSubObjTypeUrl {

		}
		//统计域名
		if result.SubObjType == enums.TaskResultSubObjTypeSubdomain {
			if result.Field1 != "" {
				if domainMap[result.Field1] == 0 {
					domainMap[result.Field1] = 1
				} else {
					domainMap[result.Field1] += 1
				}
			}
		}
	}

	// 限制各字段的最大数量
	const (
		maxPortCount      = 20 // 端口最多20个
		maxServiceCount   = 15 // 服务最多15个
		maxComponentCount = 15 // 组件最多15个
		maxSubDomainCount = 10 // 子域名最多10个
	)

	// 转换为切片并按值排序（取最常见的）
	type itemCount struct {
		label string
		value int
	}

	// 处理端口数据
	var portItems []itemCount
	for label, value := range portMap {
		portItems = append(portItems, itemCount{label: label, value: value})
	}
	sort.Slice(portItems, func(i, j int) bool {
		return portItems[i].value > portItems[j].value
	})
	for i, item := range portItems {
		if i >= maxPortCount {
			break
		}
		overViewStruct.Port = append(overViewStruct.Port, &OverViewItems{Label: item.label, Value: item.value})
	}

	// 处理服务数据
	var serviceItems []itemCount
	for label, value := range serviceMap {
		serviceItems = append(serviceItems, itemCount{label: label, value: value})
	}
	sort.Slice(serviceItems, func(i, j int) bool {
		return serviceItems[i].value > serviceItems[j].value
	})
	for i, item := range serviceItems {
		if i >= maxServiceCount {
			break
		}
		overViewStruct.Service = append(overViewStruct.Service, &OverViewItems{Label: item.label, Value: item.value})
	}

	// 处理组件数据
	var componentItems []itemCount
	for label, value := range componentMap {
		componentItems = append(componentItems, itemCount{label: label, value: value})
	}
	sort.Slice(componentItems, func(i, j int) bool {
		return componentItems[i].value > componentItems[j].value
	})
	for i, item := range componentItems {
		if i >= maxComponentCount {
			break
		}
		overViewStruct.Component = append(overViewStruct.Component, &OverViewItems{Label: item.label, Value: item.value})
	}

	// 处理子域名数据
	var domainItems []itemCount
	for label, value := range domainMap {
		domainItems = append(domainItems, itemCount{label: label, value: value})
	}
	sort.Slice(domainItems, func(i, j int) bool {
		return domainItems[i].value > domainItems[j].value
	})
	for i, item := range domainItems {
		if i >= maxSubDomainCount {
			break
		}
		overViewStruct.SubDomain = append(overViewStruct.SubDomain, &OverViewItems{Label: item.label, Value: item.value})
	}

	return nil
}

func (t *TaskTaskInfo) GetProgressStats(ctx context.Context, taskTask mysqls.TaskTask, overViewStruct *OverView) error {
	var taskLogInfo mysqls.Taskloginfo
	// 获取所有有日志的目标的日志统计
	logResult, err := taskLogInfo.GetTaskLogInfoNumberByTaskId(ctx, taskTask.ID)
	if err != nil {
		return err
	}
	// 转为Map方便查找
	logMap := make(map[int]int)
	for _, item := range logResult {
		logMap[item.TargetId] = item.Num
	}

	// 获取任务下所有目标（用于判断状态）
	var taskTarget mysqls.TaskTarget
	targetList := taskTarget.GetTargetsByTaskId(ctx, taskTask.ID)

	totalTargets := len(targetList)
	if totalTargets == 0 {
		return nil
	}

	var progress int
	for _, target := range targetList {
		// 1. 如果目标已完成（状态为Finish），视为100%进度
		// 注意：不存活的目标最终状态也会是Finish，所以只判断Status即可
		if target.Status == enums.TargetStatusFinish {
			progress += 100
			continue
		}

		// 2. 如果目标正在运行，根据日志数量估算进度
		// 只有Running状态才计算日志进度，避免Begin状态误算
		if target.Status == enums.TargetStatusRunning {
			if num, ok := logMap[target.ID]; ok {
				var progressTemp int
				if num < 10 {
					progressTemp = num / 2
				} else if num < 100 {
					progressTemp = 10/2 + (num-10)/3
				} else if num < 300 {
					progressTemp = 10/2 + 90/3 + (num-100)/5
				} else if num < 1000 {
					progressTemp = 10/2 + 90/3 + 200/5 + (num-300)/30
				} else {
					progressTemp = 99
				}
				progress += progressTemp
			}
		}
		// 3. 其他状态（如Begin）进度为0
	}

	val := progress / totalTargets
	if val == 0 && progress > 0 { // 有进度但不足1%时，显示1%
		val = 1
	}
	overViewStruct.Progress = &OverViewItems{Label: "progress", Value: val}

	return nil
}

func (t *TaskTaskInfo) GetOpSysStats(ctx context.Context, taskTask mysqls.TaskTask, overViewStruct *OverView) error {
	var taskTarget mysqls.TaskTarget
	targetList := taskTarget.GetTargetsByTaskId(ctx, taskTask.ID)
	opSysMap := make(map[string]int, 0)
	for _, target := range targetList {
		if target.OpSys == "" {
			continue
		}
		opSysMap[target.OpSys] += 1
	}

	// 限制操作系统类型最多10个
	const maxOpSysCount = 10

	type itemCount struct {
		label string
		value int
	}

	var opSysItems []itemCount
	for label, value := range opSysMap {
		opSysItems = append(opSysItems, itemCount{label: label, value: value})
	}

	// 按数量排序，取最常见的
	sort.Slice(opSysItems, func(i, j int) bool {
		return opSysItems[i].value > opSysItems[j].value
	})

	for i, item := range opSysItems {
		if i >= maxOpSysCount {
			break
		}
		overViewStruct.OpSys = append(overViewStruct.OpSys, &OverViewItems{Label: item.label, Value: item.value})
	}

	return nil
}

// GetTaskCache 获取任务缓存
func (t *TaskTaskInfo) GetTaskCache(ctx context.Context, taskId int) (map[string]string, error) {
	preKey := "smart:task:" + strconv.Itoa(taskId)
	var hash redises.RedisHash
	taskCache, err := hash.GetHashHGetAll(ctx, preKey)
	if err != nil {
		return taskCache, errors.New("GetTaskCache error: " + err.Error())
	}
	if len(taskCache) == 0 {
		var taskTarget mysqls.TaskTarget
		for _, one := range taskTarget.GetTargetsByTaskId(ctx, taskId) {
			detailMap := make(map[string]int)
			detailMap["isAlive"] = one.IsAlive
			detailMap["status"] = one.Status
			detailMap["risk"] = one.RiskLevel
			detailByte, _ := json.Marshal(detailMap)
			taskCache[strconv.Itoa(one.ID)] = string(detailByte)
		}
	}
	return taskCache, nil
}

func (t *TaskTaskInfo) UpdateTaskCache(ctx context.Context, taskId int, taskCache map[string]string) error {
	preKey := "smart:task:" + strconv.Itoa(taskId)
	var hash redises.RedisHash
	var err error
	for key, value := range taskCache {
		err = hash.SetHashHset(ctx, preKey, key, value)
	}
	if err != nil {
		return errors.New("UpdateTaskCache error: " + err.Error())
	}
	return nil
}

// UpdateTaskOverview 更新概览信息
func (t *TaskTaskInfo) UpdateTaskOverview(ctx context.Context, taskTask mysqls.TaskTask, overViewStruct OverView) error {
	var taskTaskInfo mysqls.TaskTaskInfo
	overview, err := json.Marshal(overViewStruct)
	err = taskTaskInfo.UpdateTaskTaskInfoOverview(ctx, taskTask.ID, string(overview))
	if err != nil {
		return errors.New("UpdateTaskOverview error: " + err.Error())
	}
	return nil
}

func (t *TaskTaskInfo) UpdateTaskStats(ctx context.Context, taskTask mysqls.TaskTask, isStat int) error {
	var taskTaskModel mysqls.TaskTask
	taskTaskModel.UpdateTaskStatsByTaskId(ctx, taskTask.ID, isStat)
	return nil
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// UpdateTaskStateById
func (t *TaskTaskInfo) UpdateTaskStateByTaskId(ctx context.Context, taskId, status int) error {
	var taskInfo mysqls.TaskTaskInfo
	return taskInfo.UpdateStatusByTaskId(ctx, taskId, status)
}

// UpdateStatus 更改状态通过task_id
func (t *TaskTaskInfo) FinishOneTaskInfo(ctx context.Context, taskId int) error {
	var taskInfo mysqls.TaskTaskInfo
	return taskInfo.UpdateStatusByTaskId(ctx, taskId, enums.TaskStatusFinish)
}

// UpdateStatus 更改状态通过task_id
func (t *TaskTaskInfo) StartOneTaskInfo(ctx context.Context, taskId int) error {
	var taskInfo mysqls.TaskTaskInfo
	return taskInfo.UpdateStatusByTaskId(ctx, taskId, enums.TaskStatusRunning)
}

// GetOneWaitTiming 依据允许运行的时间获取一个定时任务
func (t *TaskTaskInfo) GetOneWaitTiming(ctx context.Context, runDate string) (mysqls.TaskTaskInfo, error) {
	var taskTask mysqls.TaskTaskInfo
	return taskTask.GetOneWaitTiming(ctx, runDate)
}

// GetOneWaitCycle 依据允许运行的时间获取一个周期任务
func (t *TaskTaskInfo) GetOneWaitCycle(ctx context.Context, runDate string) (mysqls.TaskTaskInfo, error) {
	var taskTask mysqls.TaskTaskInfo
	return taskTask.GetOneWaitCycle(ctx, runDate)
}

// UpdateExecNextTime 更新下次运行时间
func (t *TaskTaskInfo) UpdateExecNextTime(ctx context.Context, taskId int, runDate string) error {
	var taskInfo mysqls.TaskTaskInfo
	return taskInfo.UpdateExecNextTimeByTaskId(ctx, taskId, runDate)
}

// GetRunningTaskInfo 获取正在运行的任务
func (t *TaskTaskInfo) GetRunningTaskInfo(ctx context.Context) []mysqls.TaskTaskInfo {
	var taskInfo mysqls.TaskTaskInfo
	return taskInfo.GetTaskByTaskStatus(ctx, enums.TaskStatusRunning)
}

// GetByTaskId 依据任务ID获取任务
func (t *TaskTaskInfo) GetByTaskId(ctx context.Context, taskId []int) []mysqls.TaskTaskInfo {
	var taskInfo mysqls.TaskTaskInfo
	return taskInfo.GetByTaskIds(ctx, taskId)
}

// GetTaskConfigInfo 获取任务配置信息
func (t *TaskTaskInfo) GetTaskConfigInfo(ctx context.Context, taskId int) (mysqls.TaskTaskInfo, error) {
	var taskInfoModel mysqls.TaskTaskInfo
	return taskInfoModel.GetTaskTaskInfoByTaskId(ctx, taskId)
}

func (t *TaskTaskInfo) GetOverviewByTaskIds(ctx context.Context, taskId []int) (overViewStruct map[int]OverView) {
	overViewStruct = make(map[int]OverView)
	var taskInfoModel mysqls.TaskTaskInfo
	taskInfos := taskInfoModel.GetByTaskIds(ctx, taskId)
	for _, item := range taskInfos {
		var overViewStructTemp OverView
		_ = json.Unmarshal([]byte(item.Overview), &overViewStructTemp)
		overViewStruct[item.TaskID] = overViewStructTemp
	}
	return
}

// GetTaskInfoByTaskId 依据任务ID获取任务
func (t *TaskTaskInfo) GetTaskInfoByTaskId(ctx context.Context, taskId int) (mysqls.TaskTaskInfo, error) {
	var taskInfo mysqls.TaskTaskInfo
	return taskInfo.GetTaskTaskInfoByTaskId(ctx, taskId)
}
