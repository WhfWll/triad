package crons

import (
	"context"
	"smart/services"
	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
)

func TaskInfoStat() {
	var (
		taskTask     services.TaskTask
		taskTaskInfo services.TaskTaskInfo
	)
	ctx := context.Background()
	//第一步 获取等待进行统计的任务
	for _, waitStatTask := range taskTask.GetWaitStatTask(ctx) {
		var overViewStruct services.OverView
		//第二步 获取目标及目标风险等级的统计
		isALLTargetDelFinish, err := taskTaskInfo.GetTargetStats(ctx, waitStatTask, &overViewStruct)
		if err != nil {
			log.Error("TaskInfoStat GetTargetStats error: " + err.Error())
		}
		//第三步 进行漏洞统计
		err = taskTaskInfo.GetVulStats(ctx, waitStatTask, &overViewStruct)
		if err != nil {
			log.Error("TaskInfoStat GetVulStats error: " + err.Error())
		}
		//第四步 进行信息统计
		err = taskTaskInfo.GetInfoStats(ctx, waitStatTask, &overViewStruct)
		if err != nil {
			log.Error("TaskInfoStat GetInfoStats error: " + err.Error())
		}
		//第四步 进行操作系统统计
		err = taskTaskInfo.GetOpSysStats(ctx, waitStatTask, &overViewStruct)
		if err != nil {
			log.Error("TaskInfoStat GetInfoStats error: " + err.Error())
		}
		//第四步 通过日志进行进度统计
		err = taskTaskInfo.GetProgressStats(ctx, waitStatTask, &overViewStruct)
		if err != nil {
			log.Error("TaskInfoStat GetInfoStats error: " + err.Error())
		}

		//第五步 更新统计信息
		err = taskTaskInfo.UpdateTaskOverview(ctx, waitStatTask, overViewStruct)
		if err != nil {
			log.Error("TaskInfoStat UpdateTaskOverview error: " + err.Error())
		}
		err = taskTask.UpdateTargetRiskById(ctx, waitStatTask.ID, overViewStruct.TargetRisk, overViewStruct.RiskLevel)
		if err != nil {
			log.Error("TaskInfoStat UpdateTaskRiskById error: " + err.Error())
		}
		//第五步 如果所有目标都停止或删除
		if isALLTargetDelFinish == true && waitStatTask.Status != enums.TaskStatusFinish {
			taskTask.FinishTask(ctx, waitStatTask.ID)
		}
		// 第六步 如果任务已经结束或者处于暂停中，那么关闭统计标识
		if isALLTargetDelFinish || waitStatTask.Status == enums.TaskStatusFinish || waitStatTask.Status == enums.TaskStatusPausing {
			err = taskTaskInfo.UpdateTaskStats(ctx, waitStatTask, enums.TaskIsStatsNo)
			if err != nil {
				log.Error("TaskInfoStat UpdateTaskStats error: " + err.Error())
			}
		}
	}
}
