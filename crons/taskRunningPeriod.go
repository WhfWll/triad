package crons

import (
	"context"
	"encoding/json"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/redis"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/data"
	"smart/tools/enums"
	"strconv"
	"time"
)

// 任务运行时间段校验
/*
定时监控渗透任务运行时间段
运行中的任务，校验其允许的运行时间段，如未在允许的时间内，则将其暂停 并记录哪些任务是自动停止的
自动停止的任务，校验其允许运行的时间，如在允许时间内，则将其恢复执行
*/
const redisNotRunningPeriod = "task_not_runtime_period"

type checkTaskRuntimePeriod struct {
}

// 定时执行
func TaskRunningPeriod() {
	fmt.Println("定时监控渗透任务运行时间段")
	ctx := context.Background()
	var taskRuntimePeriod checkTaskRuntimePeriod
	// 暂停校验
	taskRuntimePeriod.purse(ctx)
	// 恢复校验
	taskRuntimePeriod.resume(ctx)
}

// 校验是否在运行段内
func (taskRuntimePeriod checkTaskRuntimePeriod) purse(ctx context.Context) {
	var taskInfoSrv services.TaskTaskInfo
	taskInfo := taskInfoSrv.GetRunningTaskInfo(ctx)

	// 没有要处理的任务
	if len(taskInfo) == 0 {
		return
	}
	taskRuntimePeriod.process(ctx, taskInfo)
}

// 校验自动暂停的任务，需要恢复
func (taskRuntimePeriod checkTaskRuntimePeriod) resume(ctx context.Context) {
	var (
		taskInfoSrv services.TaskTaskInfo
	)
	redisClient, err := redis.NewClient()
	if err != nil {
		log.Error("任务运行时间段校验 redisClient获取失败")
		return
	}

	list := redisClient.ZRange(ctx, redisNotRunningPeriod, 0, -1).Val()

	taskIds := make([]int, 0)
	for _, item := range list {
		taskId, err := strconv.Atoi(item)
		if err != nil {
			// 出错了，删掉当前数据
			redisClient.ZRem(ctx, redisNotRunningPeriod, taskId)
		}
		taskIds = append(taskIds, taskId)
	}

	// 只有暂停中的任务才会继续执行
	// 记录最终需要处理的数据
	finalTaskInfo := make([]mysqls.TaskTaskInfo, 0)
	taskInfos := taskInfoSrv.GetByTaskId(ctx, taskIds)
	for _, item := range taskInfos {
		if item.Status != enums.TaskStatusPausing {
			// 非暂停，无需处理，删除缓存
			redisClient.ZRem(ctx, redisNotRunningPeriod, item.TaskID)
			continue
		}
		finalTaskInfo = append(finalTaskInfo, item)
	}

	// 没有要处理的任务
	if len(finalTaskInfo) == 0 {
		return
	}
	taskRuntimePeriod.process(ctx, finalTaskInfo)
}

func (taskRuntimePeriod checkTaskRuntimePeriod) process(ctx context.Context, taskInfo []mysqls.TaskTaskInfo) {
	redisClient, err := redis.NewClient()
	if err != nil {
		log.Error("任务运行时间段校验 redisClient获取失败")
		return
	}

	var (
		targetSev        services.TaskTarget
		taskSrv          services.TaskTask
		taskInfoSrv      services.TaskTaskInfo
		taskLogSrv       services.TaskLog
		taskRunningCheck data.TaskRuntimeCheck
	)

	for _, item := range taskInfo {
		// 获取runtime_period
		var execJson typespec.ExecuteJson
		err := json.Unmarshal([]byte(item.ExecuteJSON), &execJson)
		if err != nil {
			log.Error("task_id=" + strconv.Itoa(item.TaskID) + "解析runtime_period失败：" + err.Error())
			continue
		}

		//fmt.Println("任务ID", configTaskMap[int(item.Id)].Id)
		//fmt.Println("状态", configTaskMap[int(item.Id)].TaskStatus)
		//fmt.Println(isAllowRun)
		//fmt.Println("---------------")
		// 是否在运行时间段内
		if taskRunningCheck.CheckIsAllowRunning(execJson.RuntimePeriod) {
			// 允许运行
			if item.Status == enums.TaskStatusPausing {
				// 暂停的任务。需要恢复
				// 1 恢复所有目标 并 通知决策引擎
				// 2 恢复所有任务

				// 1
				targetList, _ := targetSev.GetTargetByTaskId(ctx, item.TaskID)
				for _, target := range targetList {
					// 任务状态 仅暂停中的任务允许恢复
					if target.Status == enums.TargetStatusPausing {
						// 通知决策引擎进行恢复操作
						_ = targetSev.ResumeTarget(target.ID)
						// 修改目标状态 为恢复操作
						_ = targetSev.UpdateTargetStateById(ctx, target.ID, enums.TargetStatusRunning)
						// 日志状态
						_ = taskLogSrv.UpdateTaskLogStateByTargetId(ctx, target.ID, enums.TargetStatusRunning)
					}
				}

				// 2
				_ = taskSrv.UpdateTaskStateById(ctx, item.TaskID, enums.TaskStatusRunning)
				_ = taskInfoSrv.UpdateTaskStateByTaskId(ctx, item.TaskID, enums.TaskStatusRunning)

				// 已恢复的任务，删掉缓存
				redisClient.ZRem(ctx, redisNotRunningPeriod, item.TaskID)
				// 当前任务已处理，跳过后面的时间校验
				break
			}
		} else {
			// 禁止运行
			if item.Status == enums.TaskStatusRunning {
				// 运行中的任务。需要暂停，且记录暂停的任务，方便恢复
				// 暂停
				targetList, _ := targetSev.GetTargetByTaskId(ctx, item.TaskID)
				// 任务状态
				for _, target := range targetList {
					// 任务状态 仅运行中的任务允许暂停
					if target.Status == enums.TargetStatusRunning {
						// 通知决策引擎进行暂停操作
						_ = targetSev.PurseTarget(target.ID)
						// 修改目标状态 为暂停操作
						_ = targetSev.UpdateTargetStateById(ctx, target.ID, enums.TargetStatusPausing)
						// 日志状态
						_ = taskLogSrv.UpdateTaskLogStateByTargetId(ctx, target.ID, enums.TargetStatusPausing)
					}
				}

				// 修改任务状态 为暂停操作
				_ = taskSrv.UpdateTaskStateById(ctx, item.TaskID, enums.TaskStatusPausing)
				_ = taskInfoSrv.UpdateTaskStateByTaskId(ctx, item.TaskID, enums.TaskStatusPausing)

				// 记录自动暂停的任务ID
				redisClient.ZAdd(ctx, redisNotRunningPeriod, &redis2.Z{
					Score:  float64(time.Now().Unix()),
					Member: item.TaskID,
				})
				break
			}
		}
	}
}
