package crons

import (
	"context"
	"smart/services"
	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
)

// TargetTimeout 任务超时处理
func TargetTimeout() {
	var (
		taskTarget services.TaskTarget
		taskLog    services.TaskLog
	)
	ctx := context.Background()

	// 第一步 获取超时目标ID列表
	// 注意：这里依赖 TaskLog 表的查询逻辑，默认超时时间为1小时(在model中硬编码)
	targetIds, err := taskLog.GetTimeoutTargetIds(ctx)
	if err != nil {
		log.Error("获取超时目标报错: ", err)
	}

	// 第二步 逐个处理超时目标
	for _, targetId := range targetIds {
		log.Infof("正在结束超时检测目标: %d", targetId)

		// 1. 取消本地扫描上下文 (终止本地协程)
		services.CancelTargetScan(targetId)

		// 2. 更新数据库状态为 Finish
		err = taskTarget.UpdateTargetAndLogStateById(ctx, targetId, enums.TargetStatusFinish, enums.TargetStatusFinish)
		if err != nil {
			log.Errorf("更新超时目标状态异常 (TargetID: %d): %v", targetId, err)
		}
	}

	// 第三步 兜底处理 (防止有漏网之鱼，例如 TaskTarget 中有但 TaskLog 中没有的记录)
	err = taskTarget.HandleTimeoutTargets(ctx)
	if err != nil {
		log.Error("兜底处理超时检测目标异常: ", err)
	}

	// 处理超时未结束的检测目标的日志
	err = taskLog.HandleTimeoutTaskLog(ctx)
	if err != nil {
		log.Error("兜底处理超时检测目标日志异常: ", err)
	}
}
