package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"
)

// 渗透测试任务管理服务

type TaskTask struct {
}

// 枚举信息 - 执行类型
func (t *TaskTask) ExecuteTypeEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.TaskTaskEnum.AllExecTypeEnum())
}

// 枚举信息 - 风险等级
func (t *TaskTask) RiskLevelEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.TaskTaskEnum.AllRiskEnum())
}

// 枚举信息 - 资产风险类型
func (t *TaskTask) AssetRiskEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.GetRiskEnumMap())
}

// 枚举信息 - 状态
func (t *TaskTask) StatusEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.TaskTaskEnum.AllStatusEnum())
}

// 枚举信息 - 运行时间段
func (t *TaskTask) RuntimePeriodEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.TaskTaskEnum.AllRuntimePeriodEnum())
}

// 是否存活
func (t *TaskTask) IsAliveEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.TargetIsAliveEnum())
}

// 枚举信息 - 周期执行 - 执行计划
func (t *TaskTask) CyclePlanningTypeEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.TaskTaskEnum.AllExecCycleTypeEnum())
}

// 枚举信息 - 周期执行 - 每周执行一次的值
func (t *TaskTask) CyclePlanningTypeWeekValueEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.TaskTaskEnum.AllExecCycleTypeWeekValueEnum())
}

// 枚举信息 - 周期执行 - 每月执行一次的值
func (t *TaskTask) CyclePlanningTypeMonthValueEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.TaskTaskEnum.AllExecCycleTypeMonthValueEnum())
}

// 枚举信息 - 代理模式
func (t *TaskTask) ProxyMode() []typespec.GlobalOptionsItemRes {
	var proxyConfig enums.ProxyConfig
	return toolsSort(proxyConfig.AllProxyConfigProtoEnum())
}

// 枚举信息 - 网站登陆凭证类型
func (t *TaskTask) WebLoginTypeEnum() []typespec.GlobalOptionsItemRes {
	var websiteLogin enums.WebsiteLoginConfig
	return toolsSort(websiteLogin.AllWebsiteLoginType())
}

// StartOneTask 启动一个待开始的任务
func (t *TaskTask) StartOneTask(ctx context.Context, task mysqls.TaskTask) error {
	log.Info("StartOneTask start one task: ", task.ID)
	task.Status = enums.TaskStatusRunning
	err := task.UpdateTaskCheckTask(ctx)
	if err != nil {
		return errors.New("StartOneTask one task: " + err.Error())
	}
	return nil
}

// StartTasks 启动任务
func (t *TaskTask) StartTasks(ctx context.Context, taskIds []int) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	//更新task表
	var task mysqls.TaskTask
	var param = map[string]interface{}{
		"status":      enums.TaskStatusRunning,
		"update_time": time.Now(),
	}
	err := task.UpdateTaskTaskByIds(dCtx, taskIds, param)
	if err != nil {
		return err
	}
	//更新taskInfo表
	var taskInfo mysqls.TaskTaskInfo
	var params = map[string]interface{}{
		"status":            enums.TaskStatusRunning,
		"execute_last_time": time.Now().Format(utils.DateTime),
		"update_time":       time.Now(),
	}
	err = taskInfo.UpdateTaskInfoByTaskIds(dCtx, taskIds, params)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// 获取正在运行的任务个数
func (t *TaskTask) GetRunningNumber(ctx context.Context) int64 {
	var task mysqls.TaskTask
	return task.GetTaskNumberByTaskStatus(ctx, enums.TaskStatusRunning)
}

// GetRunningTask 获取正在运行的任务个数
func (t *TaskTask) GetRunningTask(ctx context.Context) []int64 {
	var task mysqls.TaskTask
	return task.GetTaskIDByTaskStatus(ctx, enums.TaskStatusRunning)
}

// GetTaskNumByStatus 根据状态获取任务数量
func (t *TaskTask) GetTaskNumByStatus(ctx context.Context, status int) ([]int, int64, error) {
	var (
		task   mysqls.TaskTask
		total  int64
		result []int
	)
	taskRes, err := task.GetTaskNumByStatus(ctx, status)
	if err != nil {
		return result, total, err
	}
	for i := 0; i < len(taskRes); i++ {
		result = append(result, taskRes[i].ID)
	}
	total = int64(len(taskRes))
	return result, total, nil
}

// GetOneWaitingTask 获取正在运行的任务个数
func (t *TaskTask) GetOneWaitingTask(ctx context.Context) mysqls.TaskTask {
	var taskTask mysqls.TaskTask
	return taskTask.GetOneTaskByStatus(ctx, enums.TaskExecTypeImmediate, enums.TaskStatusBegin)
}

// GetWaitingTaskByLimit 获取等待执行的任务
func (t *TaskTask) GetWaitingTaskByLimit(ctx context.Context, limit int) ([]int, error) {
	var (
		taskTask mysqls.TaskTask
		result   []int
	)
	taskRes, err := taskTask.GetTasksByExecuteTypeStatusLimit(ctx, enums.TaskExecTypeImmediate, enums.TaskStatusBegin, limit)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(taskRes); i++ {
		result = append(result, taskRes[i].ID)
	}
	return result, nil
}

// GetRunningTaskAndTarget 根据任务状态和目标状态获取任务/目标数量和任务id
func (t *TaskTask) GetRunningTaskAndTarget(ctx context.Context, taskState int, targetState int) ([]int, int64, int64, error) {
	var (
		task   mysqls.TaskTask
		target mysqls.TaskTarget
		result []int
	)
	taskRes, err := task.GetTaskNumByStatus(ctx, taskState)
	if err != nil {
		return nil, 0, 0, err
	}
	for i := 0; i < len(taskRes); i++ {
		result = append(result, taskRes[i].ID)
	}
	targetTotal, err := target.GetTargetNumByStatus(ctx, targetState)
	if err != nil {
		return nil, 0, 0, err
	}
	return result, int64(len(taskRes)), targetTotal, nil
}

// 创建
func (t *TaskTask) Save(ctx context.Context, uid, taskTemplateId int, taskName string, targetList []string, executeType int,
	executeJson string, configStruct enums.ConfigJson, nextRuntime time.Time, weight, pid int) (int, error) {
	configStrByte, _ := json.Marshal(configStruct)

	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	//创建任务
	var taskModel mysqls.TaskTask
	taskModel.TaskName = taskName
	taskModel.RiskLevel = enums.TaskRiskSafe // 默认安全
	taskModel.Status = enums.TaskStatusBegin // 默认待执行
	taskModel.Weight = weight
	taskModel.TaskType = enums.TaskTypeMultipleTask
	taskModel.ExecuteType = executeType
	taskModel.TaskTemplateID = taskTemplateId
	taskModel.TargetNum = len(targetList)
	taskModel.SafeNum = taskModel.TargetNum
	taskModel.UserID = uid
	taskModel.Pid = pid
	taskModel.CreateTime = time.Now()
	taskModel.UpdateTime = time.Now()
	taskModel.IsStats = enums.TaskIsStatsYes
	err := taskModel.AddTaskCheckTask(dCtx)
	if err != nil {
		return 0, err
	}

	// 创建任务详情
	var taskInfo mysqls.TaskTaskInfo
	taskInfo.TaskID = taskModel.ID
	taskInfo.TaskName = taskModel.TaskName
	taskInfo.TaskType = taskModel.TaskType
	taskInfo.Status = enums.TaskStatusBegin
	taskInfo.Weight = weight
	taskInfo.CheckTarget = strings.Join(targetList, ",")
	taskInfo.ExecuteType = executeType
	taskInfo.ExecuteJSON = executeJson
	taskInfo.ExecuteLastTime = time.Unix(0, 0) // 默认1970年
	taskInfo.ExecuteNextTime = nextRuntime
	taskInfo.TaskTemplateID = taskTemplateId
	taskInfo.TaskTemplateJSON = string(configStrByte)
	taskInfo.Overview = ""
	taskInfo.UserID = uid
	taskInfo.CreateTime = time.Now()
	taskInfo.UpdateTime = time.Now()
	err = taskInfo.AddTaskTaskInfo(dCtx)
	if err != nil {
		return 0, err
	}

	//查询可利用评分开关
	var (
		isScore     = enums.TargetIsScoreNoneed
		mapSetMysql mysqls.MapSet
	)
	mapSetRes, err := mapSetMysql.GetsByObjKey(ctx, enums.UseScoreSwitchMapSetObjKey)
	if err == nil && len(mapSetRes.ObjValue) > 0 {
		useScoreSwitch, _ := strconv.Atoi(mapSetRes.ObjValue)
		if useScoreSwitch == enums.TargetIsScoreSwitchOn {
			isScore = enums.TargetIsScoreRunning
		}
	}

	//创建目标
	targets := make([]mysqls.TaskTarget, 0)
	for _, targetUrl := range targetList {
		//目标入库
		targets = append(targets, mysqls.TaskTarget{
			TaskID:           taskModel.ID,
			TargetURL:        targetUrl,
			Status:           enums.TargetStatusToBegin, // 默认待开始
			Weight:           weight,
			RiskLevel:        enums.TargetRiskLowNoFound,    // 默认未发现
			IsAlive:          enums.TargetIsAliveN,          // 默认未存活
			TargetType:       enums.TaskCheckTargetTypeHost, // 主机类型
			TaskTemplateID:   taskTemplateId,
			TaskTemplateJSON: string(configStrByte),
			UserID:           uid,
			IsRemoteSession:  enums.TargetIsRemoteSessionN,
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
			EndTime:          time.Unix(0, 0), // 默认1970年
			UseScore:         0,
			IsScore:          isScore,
		})
	}
	var targetModel mysqls.TaskTarget
	fmt.Println(len(targets)/3200 + 1)
	for i := 0; i < len(targets)/3200+1; i++ {
		var tempList []mysqls.TaskTarget
		if i == len(targets)/3200 {
			tempList = targets[i*3200:]
		} else {
			tempList = targets[i*3200 : (i+1)*3200-1]
		}
		err = targetModel.Adds(dCtx, &tempList)
	}
	if err != nil {
		return 0, err
	}

	//请求决策引擎进行可利用评分
	// if isScore != enums.TargetIsScoreNoneed {
	// 	var targetUrlMap []targetUrlMapItems
	// 	for i := 0; i < len(targets); i++ {
	// 		var tmp = targetUrlMapItems{
	// 			TargetId:  targets[i].ID,
	// 			TargetUrl: targets[i].TargetURL,
	// 		}
	// 		targetUrlMap = append(targetUrlMap, tmp)
	// 	}
	// 	//请求决策引擎计算利用评分
	// 	var params = map[string]interface{}{
	// 		"taskId":       taskModel.ID,
	// 		"targetUrlMap": targetUrlMap,
	// 		"toScanPort":   configStruct.PortScanConfig.ScanPort,
	// 		"tcpScanType":  configStruct.PortScanConfig.TCPScanType,
	// 	}
	// 	err = httpclients.GetUseScore(params)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// }

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return taskModel.ID, nil
}

// 创建
func (t *TaskTask) AddTarget(ctx context.Context, targetList []string, taskRes mysqls.TaskTask, taskinfoRes mysqls.TaskTaskInfo, userId int) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	//修改任务
	var param = map[string]interface{}{
		"status":      enums.TaskStatusBegin, // 待执行
		"is_stats":    enums.TaskIsStatsYes,
		"target_num":  taskRes.TargetNum + len(targetList),
		"safe_num":    taskRes.SafeNum + len(targetList),
		"update_time": time.Now(),
	}
	err := taskRes.UpdateTaskTaskByIds(dCtx, []int{taskRes.ID}, param)
	if err != nil {
		return err
	}
	// 修改任务详情
	var checkTarget = taskinfoRes.CheckTarget
	if len(checkTarget) > 0 {
		checkTarget = checkTarget + "," + strings.Join(targetList, ",")
	} else {
		checkTarget = strings.Join(targetList, ",")
	}
	var params = map[string]interface{}{
		"status":       enums.TaskStatusBegin,
		"check_target": checkTarget,
		"update_time":  time.Now(),
	}
	err = taskinfoRes.UpdateTaskInfoByTaskIds(dCtx, []int{taskRes.ID}, params)
	if err != nil {
		return err
	}
	//创建目标
	targets := make([]mysqls.TaskTarget, 0)
	for _, targetUrl := range targetList {
		//目标入库
		targets = append(targets, mysqls.TaskTarget{
			TaskID:           taskRes.ID,
			TargetURL:        targetUrl,
			Status:           enums.TargetStatusToBegin, // 默认待开始
			Weight:           taskRes.Weight,
			RiskLevel:        enums.TargetRiskLowNoFound,    // 默认未发现
			IsAlive:          enums.TargetIsAliveN,          // 默认未存活
			TargetType:       enums.TaskCheckTargetTypeHost, // 主机类型
			TaskTemplateID:   taskinfoRes.TaskTemplateID,
			TaskTemplateJSON: taskinfoRes.TaskTemplateJSON,
			UserID:           userId,
			IsRemoteSession:  enums.TargetIsRemoteSessionN,
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
			EndTime:          time.Unix(0, 0), // 默认1970年
			UseScore:         0,
			IsScore:          enums.TargetIsScoreNoneed,
			ExtendField:      "",
		})
	}
	var targetModel mysqls.TaskTarget
	err = targetModel.Adds(dCtx, &targets)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// 列表
func (t *TaskTask) List(ctx context.Context, page, size int, taskName string, riskLevel int, startTime, endTime string, userIdList []int) ([]mysqls.TaskTask, int, error) {
	var model mysqls.TaskTask
	tempTaskIdList := make([]int, 0)
	list, total, err := model.GetTaskCheckTaskList(ctx, page, size, riskLevel, taskName, startTime, endTime, tempTaskIdList, userIdList)
	return list, int(total), err
}

// FinishTask 完成一个正在运行的任务
func (t *TaskTask) FinishTask(ctx context.Context, taskID int) {
	var task mysqls.TaskTask
	err := task.UpdateTaskCheckTaskStatusByID(ctx, taskID, enums.TaskStatusFinish)
	if err != nil {
		log.Errorf("task finish error: %s", err)
	}
	// 更新资产信息
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("UpdateAssetInfo panic: %v", r)
			}
		}()
		log.Printf("FinishTask Begin UpdateAssetInfo taskID:%d", taskID)
		assetSrv := Asset{}
		assetSrv.UpdateAssetInfo(context.Background(), taskID)
	}()
}

// GetTaskByTaskId 根据任务id获取任务
func (t *TaskTask) GetTaskByTaskId(ctx context.Context, taskId int) (mysqls.TaskTask, error) {
	var (
		taskModels mysqls.TaskTask
		taskResult mysqls.TaskTask
	)
	taskResult, err := taskModels.GetTaskCheckTask(ctx, taskId)
	if err != nil {
		return taskResult, err
	}
	return taskResult, nil
}

// GetTaskTaskinfoByTaskId 根据任务id获取任务和任务详情信息
func (t *TaskTask) GetTaskTaskinfoByTaskId(ctx context.Context, taskId int) (mysqls.TaskTask, mysqls.TaskTaskInfo, error) {
	var (
		taskModels     mysqls.TaskTask
		taskInfoModels mysqls.TaskTaskInfo
		taskResult     mysqls.TaskTask
		taskInfoResult mysqls.TaskTaskInfo
	)
	taskResult, err := taskModels.GetTaskCheckTask(ctx, taskId)
	if err != nil {
		return taskResult, taskInfoResult, err
	}
	taskInfoResult, err = taskInfoModels.GetTaskTaskInfoByTaskId(ctx, taskId)
	if err != nil {
		return taskResult, taskInfoResult, err
	}
	return taskResult, taskInfoResult, nil
}

// 删除任务
func (t *TaskTask) DelTaskByIds(ctx context.Context, ids []int) error {
	var taskModel mysqls.TaskTask
	return taskModel.DeleteByIds(ctx, ids)
}

// 删除任务详情
func (t *TaskTask) DelTaskInfoByTaskId(ctx context.Context, taskIds []int) error {
	var taskInfoModel mysqls.TaskTaskInfo
	return taskInfoModel.DeleteByTaskIds(ctx, taskIds)
}

// GetWaitStatTask 获取待统计任务
func (t *TaskTask) GetWaitStatTask(ctx context.Context) []mysqls.TaskTask {
	var task mysqls.TaskTask
	return task.GetTaskByIsStats(ctx, enums.TaskIsStatsYes)
}

// UpdateIsStats 更新是否重新统计字段
func (t *TaskTask) UpdateIsStats(ctx context.Context, id int, isStats int) error {
	var taskInfoModel mysqls.TaskTask
	return taskInfoModel.UpdateTaskTaskIsStats(ctx, id, isStats)
}

// UpdateTaskStateById 修改目标状态
func (t *TaskTask) UpdateTaskStateById(ctx context.Context, id int, state int) error {
	var taskTask mysqls.TaskTask
	return taskTask.UpdateTaskCheckTaskStatusByID(ctx, id, state)
}

// UpdateTaskStateById 修改目标状态
func (t *TaskTask) UpdateTaskStateIsStatsById(ctx context.Context, id int, state int, isStats int) error {
	var taskTask mysqls.TaskTask
	return taskTask.UpdateTaskStateIsStatsById(ctx, id, state, isStats)
}

func (t *TaskTask) UpdateTargetRiskById(ctx context.Context, id int, targetRisk [4]int, riskLevel int) error {
	var taskTask mysqls.TaskTask
	taskTask.UpdateTargetRiskById(ctx, id, targetRisk[0], targetRisk[1], targetRisk[2], targetRisk[3], riskLevel)
	return nil
}

// UpdateTaskRisk 修改任务风险等级
func (t *TaskTask) UpdateTaskRisk(ctx context.Context, id, risk int) error {
	var taskTaskModel mysqls.TaskTask
	taskTask, err := taskTaskModel.GetTaskCheckTask(ctx, id)

	if risk == enums.VulLibrariesRiskDead || risk == enums.VulLibrariesRiskHigh {
		if taskTask.RiskLevel > enums.TaskRiskHigh {
			taskTask.RiskLevel = enums.TaskRiskHigh
		}
	} else if risk == enums.VulLibrariesRiskMiddle {
		if taskTask.RiskLevel > enums.TaskRiskMid {
			taskTask.RiskLevel = enums.TaskRiskMid
		}
	} else if risk == enums.VulLibrariesRiskLow {
		if taskTask.RiskLevel > enums.TaskRiskLow {
			taskTask.RiskLevel = enums.TaskRiskLow
		}
	}
	err = taskTask.UpdateTaskCheckTask(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetTaskByStatus 通过状态获取任务
func (t *TaskTask) GetTaskByStatus(ctx context.Context, status int) []mysqls.TaskTask {
	var task mysqls.TaskTask
	return task.GetTasksByStatus(ctx, status)
}

// GetTaskCount 获取任务总数或根据开始时间获取任务总数
func (t *TaskTask) GetTaskCount(ctx context.Context, startTime string, uid int, role int) (int, int) {
	var taskModel mysqls.TaskTask
	total, filterTotal := taskModel.GetTaskCount(ctx, startTime, uid, role)
	return int(total), int(filterTotal)
}

// GetTaskTrendStat 获取任务趋势统计
func (t *TaskTask) GetTaskTrendStat(ctx context.Context, startTime, dateFormat string, uid int, role int) []mysqls.TaskTrendStat {
	var taskModel mysqls.TaskTask
	return taskModel.GetTaskTrendStat(ctx, startTime, dateFormat, uid, role)
}

// CheckRequestIsRepeat 检测是否重复请求
func (t *TaskTask) CheckRequestIsRepeat(ctx context.Context, userId, templateId int, taskName, target string) bool {
	redisSrv, err := redis.NewClient()
	if err != nil {
		return false
	}
	cacheKey := strconv.Itoa(userId) + strconv.Itoa(templateId) + taskName + target
	if redisSrv.Get(ctx, cacheKey).Val() == "" {
		cmd := redisSrv.Set(ctx, cacheKey, "1", time.Duration(2*time.Second))
		fmt.Println(cmd.Val())
		return false
	}
	return true
}
