package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type TaskGroupTask struct {
	ID      int `gorm:"column:id;primary_key" json:"id"` // 主键
	TaskID  int `gorm:"column:task_id" json:"taskID"`    // 任务id
	GroupID int `gorm:"column:group_id" json:"groupID"`  // 任务组id
}

// TableName sets insert table name for this struct type
func (t *TaskGroupTask) TableName() string {
	return "task_group_task"
}

// Get retrieves a list of taskGroupTask from database
func (t *TaskGroupTask) GetTaskGroupTaskList(ctx context.Context, page, limit int) ([]TaskGroupTask, int64, error) {
	var (
		taskGroupTaskList []TaskGroupTask
		count             int64
		db                = mysql.FromContext(ctx).Model(&TaskGroupTask{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&taskGroupTaskList)
	db.Count(&count)

	return taskGroupTaskList, count, nil
}

// Get retrieves a single record of taskGroupTask from database
func (t *TaskGroupTask) GetTaskGroupTask(ctx context.Context) (TaskGroupTask, error) {
	var (
		taskGroupTask TaskGroupTask
		err           error
		db            = mysql.FromContext(ctx).Model(&TaskGroupTask{})
	)

	curErr := db.Where("id = ?", t.ID).First(&taskGroupTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return taskGroupTask, err
}

// Add persists taskGroupTask to database
func (t *TaskGroupTask) AddTaskGroupTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroupTask{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// Update changes taskGroupTask by id
func (t *TaskGroupTask) UpdateTaskGroupTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroupTask{})

	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// Delete taskGroupTask by id
func (t *TaskGroupTask) DeleteTaskGroupTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroupTask{})

	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// GetTaskIdByGroupId 通过组id获取任务id
func (t *TaskGroupTask) GetTaskIdByGroupId(ctx context.Context, groupId int) ([]int, error) {
	var (
		taskGroupTaskList []TaskGroupTask
		db                = mysql.FromContext(ctx).Model(&TaskGroupTask{})
		taskIdList        []int
	)
	db.Where("group_id = ?", groupId).Find(&taskGroupTaskList)

	for _, taskGroup := range taskGroupTaskList {
		taskIdList = append(taskIdList, taskGroup.TaskID)
	}
	return taskIdList, nil
}
