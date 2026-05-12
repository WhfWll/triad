package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type TaskGroup struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Name       string    `gorm:"column:name" json:"name"`              // 任务组名称
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 更新时间
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	Status     int       `gorm:"column:status" json:"status"`          // 任务状态
	Describe   string    `gorm:"column:describe" json:"describe"`      // 任务组描述
	HighNum    int       `gorm:"column:high_num" json:"high_num"`      // 高危任务数量
	MiddleNum  int       `gorm:"column:middle_num" json:"middle_num"`  // 中危任务数量
	LowNum     int       `gorm:"column:low_num" json:"low_num"`        // 低危任务数量
	SafeNum    int       `gorm:"column:safe_num" json:"safe_num"`      // 安全任务数量
	Overview   string    `gorm:"column:overview" json:"overview"`      // 统计信息
	IsStat     int       `gorm:"column:is_stat" json:"is_stat"`        // 是否统计
	//RiskLevel  int       `gorm:"column:risk_level" json:"risk_level"`  //风险等级
}

// TableName sets insert table name for this struct type
func (t *TaskGroup) TableName() string {
	return "task_group"
}

// Get retrieves a list of taskGroup from database
func (t *TaskGroup) GetTaskGroupList(ctx context.Context, search string, page, limit int) ([]TaskGroup, int64, error) {
	var (
		taskGroupList []TaskGroup
		count         int64
		db            = mysql.FromContext(ctx).Model(&TaskGroup{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&taskGroupList)
	db.Count(&count)

	return taskGroupList, count, nil
}

// Get retrieves a single record of taskGroup from database
func (t *TaskGroup) GetTaskGroup(ctx context.Context) (TaskGroup, error) {
	var (
		taskGroup TaskGroup
		err       error
		db        = mysql.FromContext(ctx).Model(&TaskGroup{})
	)
	curErr := db.Where("id = ?", t.ID).First(&taskGroup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskGroup, err
}

// Add persists taskGroup to database
func (t *TaskGroup) AddTaskGroup(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroup{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// Update changes taskGroup by id
func (t *TaskGroup) UpdateTaskGroup(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroup{})

	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// Delete taskGroup by id
func (t *TaskGroup) DeleteTaskGroup(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroup{})

	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

// GetTaskGroupByIsStat 按是否统计获取任务组
func (t *TaskGroup) GetTaskGroupByIsStat(ctx context.Context, isStat int) ([]TaskGroup, error) {
	var (
		taskGroupList []TaskGroup
		err           error
		db            = mysql.FromContext(ctx).Model(&TaskGroup{})
	)
	curErr := db.Where("is_stat = ?", isStat).Find(&taskGroupList).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskGroupList, err
}

// UpdateTaskGroupOverview 更新任务组统计信息
func (t *TaskGroup) UpdateTaskGroupOverview(ctx context.Context, overview string) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroup{})
	var param = map[string]interface{}{
		"overview": overview,
	}
	if err := db.Where("id = ?", t.ID).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTaskGroupStatus 更新任务组状态
func (t *TaskGroup) UpdateTaskGroupStatus(ctx context.Context, status, isStat int) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroup{})
	var param = map[string]interface{}{
		"status":  status,
		"is_stat": isStat,
	}
	if err := db.Where("id = ?", t.ID).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTaskGroupNameAndDescribe 更新任务组名称和描述
func (t *TaskGroup) UpdateTaskGroupNameAndDescribe(ctx context.Context, id int, name, describe string) error {
	var db = mysql.FromContext(ctx).Model(&TaskGroup{})

	if err := db.Where("id = ?", id).Updates(TaskGroup{Name: name, Describe: describe}).Error; err != nil {
		return err
	}

	return nil
}
