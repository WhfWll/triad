package services

import (
	"context"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

type TaskGroup struct {
}

// Create 任务组新建
func (t *TaskGroup) Create(ctx context.Context, name, describe string) error {
	var group = mysqls.TaskGroup{
		Name:       name,
		Describe:   describe,
		Status:     enums.TaskStatusRunning,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		IsStat:     enums.TaskGroupIsStatNo,
	}
	return group.AddTaskGroup(ctx)
}

// List 任务组列表
func (t *TaskGroup) List(ctx context.Context, search string, page, limit int) ([]mysqls.TaskGroup, int64, error) {
	var group mysqls.TaskGroup
	return group.GetTaskGroupList(ctx, search, page, limit)
}

// Delete 任务组删除
func (t *TaskGroup) Delete(ctx context.Context, id int) error {
	var group = mysqls.TaskGroup{
		ID: id,
	}
	return group.DeleteTaskGroup(ctx)
}

// GroupBind 任务组绑定
func (t *TaskGroup) GroupBind(ctx context.Context, taskId, groupId int) error {
	var group = mysqls.TaskGroupTask{
		TaskID:  taskId,
		GroupID: groupId,
	}
	return group.AddTaskGroupTask(ctx)
}

// TaskList 任务组内任务列表
func (t *TaskGroup) TaskList(ctx context.Context, groupId, page, size int) ([]mysqls.TaskTask, int, error) {
	var (
		groupModel mysqls.TaskGroupTask
		taskModel  mysqls.TaskTask
		taskList   []mysqls.TaskTask
		total      int64
	)
	taskIdList, err := groupModel.GetTaskIdByGroupId(ctx, groupId)
	if err != nil || len(taskIdList) == 0 {
		return taskList, 0, err
	}
	tempUserIdList := make([]int, 0)
	taskList, total, err = taskModel.GetTaskCheckTaskList(ctx, page, size, 0, "", "", "", taskIdList, tempUserIdList)
	return taskList, int(total), err
}

// GetGroupByIsStat 通过是否统计获取任务组
func (t *TaskGroup) GetGroupByIsStat(ctx context.Context, isStat int) ([]mysqls.TaskGroup, error) {
	var group mysqls.TaskGroup
	return group.GetTaskGroupByIsStat(ctx, isStat)
}

// GetGroup 获取任务组
func (t *TaskGroup) GetGroup(ctx context.Context, id int) (mysqls.TaskGroup, error) {
	var group = mysqls.TaskGroup{
		ID: id,
	}
	return group.GetTaskGroup(ctx)
}

// Edit 任务组编辑
func (t *TaskGroup) Edit(ctx context.Context, id int, name, describe string) error {
	var group mysqls.TaskGroup
	return group.UpdateTaskGroupNameAndDescribe(ctx, id, name, describe)
}
