package services

import (
	"context"
	"errors"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

// 日志管理 - 日志管理

type TaskLog struct {
}

// GetTaskLog 查询
func (t *TaskLog) GetTaskLog(ctx context.Context, targetId int) (mysqls.Tasklog, error) {
	var taskLog = &mysqls.Tasklog{
		TargetID: targetId,
	}
	tl, err := taskLog.GetTaskLogByTargetId(ctx)
	if err != nil {
		return tl, errors.New("GetLogId error: " + err.Error())
	}
	return tl, nil
}

// GetManyTaskLogsByTargets 根据目标id获取日志数据,并返回对应关系map
func (t *TaskLog) GetManyTaskLogsByTargets(ctx context.Context, targetIds any) (map[int]int, error) {
	var (
		taskLog mysqls.Tasklog
		result  = make(map[int]int, 0)
	)
	logRes, err := taskLog.GetTasklogListByTargetIds(ctx, targetIds)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(logRes); i++ {
		result[logRes[i].TargetID] = logRes[i].ID
	}
	return result, nil
}

// 添加一条检测日志
func (t *TaskLog) AddTaskLog(ctx context.Context, taskID, targetID, status int, targetURL string) error {
	var taskLog = mysqls.Tasklog{
		TaskID:     taskID,
		TargetID:   targetID,
		TargetURL:  targetURL,
		Status:     status,
		CreateTime: time.Now(),
		StartTime:  time.Now(),
		EndTime:    time.Now(),
	}
	err := taskLog.AddTasklog(ctx)
	if err != nil {
		return err
	}
	return nil
}

// FinishTaskLog 结束任务日志
func (t *TaskLog) FinishTaskLog(ctx context.Context, logId int) error {
	var taskLog = mysqls.Tasklog{
		ID: logId,
	}
	err := taskLog.UpdateTaskLogStatus(ctx, enums.TaskStatusFinish)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTaskLogStateByTargetId 结束任务日志
func (t *TaskLog) UpdateTaskLogStateByTargetId(ctx context.Context, targetId, status int) error {
	var taskLog mysqls.Tasklog
	err := taskLog.UpdateTaskLogStateByTargetId(ctx, targetId, status)
	if err != nil {
		return err
	}
	return nil
}

// LogList 测试日志列表及筛选
func (t *TaskLog) LogList(ctx context.Context, taskId int, search string, page, limit int) ([]mysqls.Tasklog, int64, error) {
	var tasklog mysqls.Tasklog
	return tasklog.GetTasklogList(ctx, taskId, search, page, limit)
}

// LogInfo 测试日志详情
func (t *TaskLog) LogInfo(ctx context.Context, taskLogId int) ([]mysqls.Taskloginfo, int64, error) {
	var tasklog mysqls.Taskloginfo
	return tasklog.GetTaskloginfoListById(ctx, taskLogId)
}

// 批量删除 通过task_id
func (t *TaskLog) DelTaskInfoByTaskId(ctx context.Context, taskIds []int) error {
	var taskLogModel mysqls.Tasklog
	return taskLogModel.DeleteByTaskIds(ctx, taskIds)
}

// GetTimeoutTargetIds 处理超时的检测目标日志
func (t *TaskLog) GetTimeoutTargetIds(ctx context.Context) ([]int, error) {
	var taskLog mysqls.Tasklog
	return taskLog.GetTimeoutTargetIds(ctx)
}

// HandleTimeoutTaskLog 处理超时的检测目标日志
func (t *TaskLog) HandleTimeoutTaskLog(ctx context.Context) error {
	var taskLog mysqls.Tasklog
	err := taskLog.HandleTimeoutTaskLog(ctx)
	return err
}

// UpdateTaskLogAliveByTargetId 更新日志存活通过
func (t *TaskLog) UpdateTaskLogAliveByTargetId(ctx context.Context, targetId, isAlive int) error {
	var taskLog mysqls.Tasklog
	err := taskLog.UpdateTaskLogAliveByTargetId(ctx, targetId, isAlive)
	if err != nil {
		return err
	}
	return nil
}

// GetFirstTenTaskLog 获取前十条日志
func (t *TaskLog) GetFirstFiveTaskLog(ctx context.Context) ([]mysqls.Tasklog, error) {
	var taskLogModel mysqls.Tasklog
	return taskLogModel.GetFirstFiveTaskLog(ctx)
}
