package services

import (
	"context"
	"smart/models/mysqls"
)

// 通用记录数据表

type TaskTargetResult struct {
}

func (t *TaskTargetResult) DelTaskInfoByTaskId(ctx context.Context, taskIds []int) error {
	var taskTargetResultModel mysqls.TaskTargetResult
	return taskTargetResultModel.DeleteByTaskIds(ctx, taskIds)
}
