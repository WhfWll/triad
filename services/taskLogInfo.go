package services

import (
	"context"
	log "github.com/sirupsen/logrus"
	"smart/models/mysqls"
	"time"
)

// 用户管理 - 用户管理

type TaskLogInfo struct {
}

// AddTaskLogInfo 插入一条检测日志
func (t *TaskLogInfo) AddTaskLogInfo(ctx context.Context, taskId, targetId, logId int, targetURL string, pocname string, result string) error {
	var taskLogInfo = mysqls.Taskloginfo{
		TaskLogID:  logId,
		TaskID:     taskId,
		TargetID:   targetId,
		TargetURL:  targetURL,
		Pocname:    pocname,
		Result:     result,
		CreateTime: time.Now(),
	}
	err := taskLogInfo.AddTaskloginfo(ctx)
	if err != nil {
		log.Errorf("add check result error: %s", err)
		return err
	}
	return nil
}

// 批量删除 通过task_id
func (t *TaskLogInfo) DelTaskInfoByTaskId(ctx context.Context, taskIds []int) error {
	var taskLogInfoModel mysqls.Taskloginfo
	return taskLogInfoModel.DeleteByTaskIds(ctx, taskIds)
}

// GetTaskLogInfoListByLogIds 通过logIds获取检测日志
func (t *TaskLogInfo) GetTaskLogInfoListByLogIds(ctx context.Context, logIds []int) ([]mysqls.Taskloginfo, error) {
	var taskLogInfo mysqls.Taskloginfo
	return taskLogInfo.GetTaskLogInfoListByLogIds(ctx, logIds)
}
