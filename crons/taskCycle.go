package crons

import (
	"context"
	"encoding/json"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/data"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// 任务运行方式 - 定时执行 / 周期执行
func TaskCycle() {
	var (
		taskSrv     services.TaskTask
		taskInfoSrv services.TaskTaskInfo
	)
	ctx := context.Background()
	currentTime := time.Now().Format(utils.DateTime)

	taskAllow, _ := isAllowRunningTaskAndTarget()
	if !taskAllow {
		log.Error("任务运行方式 - 超出最大任务运行数量，停止检查")
		return
	}

	// 定时任务
	go TaskCycleTimer(taskSrv, taskInfoSrv, ctx, currentTime)

	// 周期任务
	go TaskCycleCycle(taskSrv, taskInfoSrv, ctx, currentTime)
}

// 任务是否允许运行
func isAllowRunningTaskAndTarget() (taskAllow, targetAllow bool) {
	var taskControl TaskControl
	//if err := config.Load("task_control", &taskControl); err != nil {
	//	log.Error("TargetCreate error: get task control failed, please check the configuration")
	//	return
	//}
	var mapSetService services.MapSet
	ctx := context.Background()
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.CurTasksInfoMapSetObjKey)
	if err != nil {
		log.Error("TargetCreate error: get task control failed", err)
		return
	}
	if objValueStr != "" {
		if err = json.Unmarshal([]byte(objValueStr), &taskControl); err != nil {
			log.Error("TargetCreate unmarshal config err", err)
			return
		}
	} else {
		taskControl.MaxRunningTaskNumber = 5
		taskControl.MaxRunningTargetNumber = 10
	}

	var taskTask services.TaskTask
	taskIDs := taskTask.GetRunningTask(ctx)
	if int64(len(taskIDs)) < taskControl.MaxRunningTaskNumber {
		taskAllow = true
	}

	var taskTarget services.TaskTarget
	runningTargetNumber := taskTarget.GetRunningNumber(ctx)
	if int64(runningTargetNumber) < taskControl.MaxRunningTargetNumber {
		targetAllow = true
	}
	return
}

func TaskCycleTimer(taskSrv services.TaskTask, taskInfoSrv services.TaskTaskInfo, ctx context.Context, currentTime string) {
	// 获取待运行的定时任务
	timerTaskInfo, err := taskInfoSrv.GetOneWaitTiming(ctx, currentTime)
	if err != nil {
		log.Error("任务运行方式 - 定时执行：获取任务失败")
		return
	}
	if timerTaskInfo.ID == 0 {
		// log.Info("任务运行方式 - 定时执行未检测到任务")
		return
	}
	log.Info("任务运行方式 - 定时执行检测到待运行任务 task_id=" + strconv.Itoa(timerTaskInfo.TaskID))

	// 修改任务状态为运行中即可，后续逻辑由任务发送处理
	TaskCycleCycleSetRunning(ctx, taskSrv, taskInfoSrv, timerTaskInfo.TaskID)
}

func TaskCycleCycle(taskSrv services.TaskTask, taskInfoSrv services.TaskTaskInfo, ctx context.Context, currentTime string) {
	// 周期任务运行时 将当前任务设置为运行中，并创建一个新的周期任务
	// 获取周期任务逻辑
	// 1 任务类型为 enums.TaskExecTypeCycle
	// 2 下次运行时间小于当前时间
	// 3 任务结束时间大于当前时间
	// 将上述逻辑拿到的任务设置为运行中 task表与task_info表
	// 任务设置完 计算下次运行时间 并创建一个新的任务(前提是还可以运行下次，如果超过最大运行时间则将下次运行时间改为2099年)

	cycleTaskInfo, err := taskInfoSrv.GetOneWaitCycle(ctx, currentTime)
	if err != nil {
		log.Error("任务运行方式 - 周期执行task_id=" + strconv.Itoa(cycleTaskInfo.TaskID) + "：获取任务失败" + err.Error())
		return
	}
	if cycleTaskInfo.ID == 0 {
		// log.Info("任务运行方式 - 周期执行未检测到任务")
		return
	}
	log.Info("任务运行方式 - 周期执行检测到待运行任务 task_id=" + strconv.Itoa(cycleTaskInfo.TaskID))

	var execJson typespec.ExecuteJson
	err = json.Unmarshal([]byte(cycleTaskInfo.ExecuteJSON), &execJson)
	if err != nil {
		log.Error("任务运行方式 - 周期执行task_id=" + strconv.Itoa(cycleTaskInfo.TaskID) + "：解析运行json错误" + err.Error())
		return
	}
	endTime, err := time.Parse(utils.DateTime, execJson.EndTime)
	if err != nil {
		log.Error("任务运行方式 - 周期执行task_id=" + strconv.Itoa(cycleTaskInfo.TaskID) + "：解析运行时间错误" + err.Error())
		return
	}
	currentTimeObj, _ := time.Parse(utils.DateTime, currentTime)

	if endTime.Unix() < currentTimeObj.Unix() {
		// 任务已结束，不需要再次运行，将任务下次运行时间改为2099年
		log.Info("任务运行方式 - 周期执行task_id=" + strconv.Itoa(cycleTaskInfo.TaskID) + "：无需下次运行")
		_ = taskInfoSrv.UpdateExecNextTime(ctx, cycleTaskInfo.TaskID, "2099-01-01 00:00:00")
		return
	}

	// 将上述逻辑拿到的任务设置为运行中
	TaskCycleCycleSetRunning(ctx, taskSrv, taskInfoSrv, cycleTaskInfo.TaskID)

	// 计算下次运行时间
	var taskExecNextTime data.TaskExecNextTime
	taskExecNextTime.StartTime, _ = time.ParseInLocation(utils.DateTime, execJson.StartTime, time.Local)
	taskExecNextTime.CyclePlanningType = execJson.CyclePlanningType
	taskExecNextTime.CyclePlanningValue = execJson.CyclePlanningValue
	taskExecNextTime.CyclePlanningHour = execJson.CyclePlanningHour
	taskExecNextTime.EndTime, _ = time.ParseInLocation(utils.DateTime, execJson.EndTime, time.Local)
	nextRuntime, err := taskExecNextTime.Compute(cycleTaskInfo.ExecuteType)
	if err != nil {
		log.Error("任务运行方式 - 周期执行task_id=" + strconv.Itoa(cycleTaskInfo.TaskID) + "：计算下次运行时间错误：" + err.Error())
		return
	}
	// 结束时间如果小于下次运行时间 不需要再次运行，将任务下次运行时间改为2099年
	if endTime.Unix() < nextRuntime.Unix() {
		log.Info("任务运行方式 - 周期执行task_id=" + strconv.Itoa(cycleTaskInfo.TaskID) + "：无需下次运行")
		_ = taskInfoSrv.UpdateExecNextTime(ctx, cycleTaskInfo.TaskID, "2099-01-01 00:00:00")
		return
	}

	// 创建一个新的周期任务
	// 获取目标列表
	var targetSrv services.TaskTarget
	targets, _ := targetSrv.GetTargetByTaskId(ctx, cycleTaskInfo.TaskID)
	targetList := make([]string, 0, len(targets))
	for _, item := range targets {
		targetList = append(targetList, item.TargetURL)
	}
	// 获取配置
	var conf enums.ConfigJson
	err = json.Unmarshal([]byte(cycleTaskInfo.TaskTemplateJSON), &conf)
	if err != nil {
		log.Error("任务运行方式 - 周期执行task_id=" + strconv.Itoa(cycleTaskInfo.TaskID) + "：获取配置错误：" + err.Error())
		return
	}
	taskId, err := taskSrv.Save(ctx, 1, cycleTaskInfo.TaskTemplateID, cycleTaskInfo.TaskName, targetList, cycleTaskInfo.ExecuteType, cycleTaskInfo.ExecuteJSON, conf, nextRuntime, cycleTaskInfo.Weight, 0)
	if err != nil {
		log.Error("任务运行方式 - 周期任务task_id=" + strconv.Itoa(taskId) + "：创建失败：" + err.Error())
		return
	}
	log.Info("任务运行方式 - 周期自动创建成功task_id=" + strconv.Itoa(taskId))
}

func TaskCycleCycleSetRunning(ctx context.Context, taskSrv services.TaskTask, taskInfoSrv services.TaskTaskInfo, taskId int) {
	_ = taskSrv.UpdateTaskStateById(ctx, taskId, enums.TaskStatusRunning)
	_ = taskInfoSrv.StartOneTaskInfo(ctx, taskId)
}
