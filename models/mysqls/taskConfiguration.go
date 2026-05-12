package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type TaskConfiguration struct {
	ID         uint      `gorm:"column:id;primary_key" json:"id"`
	ObjId      int       `gorm:"column:obj_id" json:"objId"`
	ConfigKey  string    `gorm:"column:config_key" json:"configKey"`
	ConfigJson string    `gorm:"column:config_json" json:"configJson"`
	UserId     int       `gorm:"column:user_id" json:"userId"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName sets insert table name for this struct type
func (t *TaskConfiguration) TableName() string {
	return "task_configuration"
}

// Get retrieves a list of taskConfiguration from database
func (t *TaskConfiguration) GetTaskConfigurationList(ctx context.Context, page, limit int) ([]TaskConfiguration, int64, error) {
	var (
		taskConfigurationList []TaskConfiguration
		count                 int64
		db                    = mysql.FromContext(ctx).Model(&TaskConfiguration{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&taskConfigurationList)
	db.Count(&count)

	return taskConfigurationList, count, nil
}

// Get retrieves a single record of taskConfiguration from database
func (t *TaskConfiguration) GetTaskConfiguration(ctx context.Context) (TaskConfiguration, error) {
	var (
		taskConfiguration TaskConfiguration
		err               error
		db                = mysql.FromContext(ctx).Model(&TaskConfiguration{})
	)

	curErr := db.Where("id = ?", t.ID).First(&taskConfiguration).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return taskConfiguration, err
}

// Add persists taskConfiguration to database
func (t *TaskConfiguration) AddTaskConfiguration(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskConfiguration{})

	t.CreateTime = time.Now()
	t.UpdateTime = time.Now()
	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// GetTaskConfigByTaskId retrieves task configuration by task ID
func (t *TaskConfiguration) GetTaskConfigByTaskId(ctx context.Context, taskId int) ([]TaskConfiguration, error) {
	var list []TaskConfiguration
	err := mysql.FromContext(ctx).Model(&TaskConfiguration{}).Where("obj_id = ?", taskId).Find(&list).Error
	return list, err
}

// Update changes taskConfiguration by id
func (t *TaskConfiguration) UpdateTaskConfiguration(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskConfiguration{})

	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// Delete taskConfiguration by id
func (t *TaskConfiguration) DeleteTaskConfiguration(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskConfiguration{})

	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

// 批量添加
func (t *TaskConfiguration) AddAll(ctx context.Context, list []TaskConfiguration) error {
	var db = mysql.FromContext(ctx).Model(&TaskConfiguration{})
	if err := db.Create(list).Error; err != nil {
		return err
	}
	return nil
}

// 通过obj_id获取所有配置
func (t *TaskConfiguration) AllByObjId(ctx context.Context, objId int, configKey string) []TaskConfiguration {
	var (
		taskConfiguration []TaskConfiguration
		db                = mysql.FromContext(ctx).Model(&TaskConfiguration{})
	)

	if configKey != "" {
		db = db.Where("config_key = ?", configKey)
	}

	db.Where("obj_id = ?", objId).Find(&taskConfiguration)

	return taskConfiguration
}

// Get retrieves a single record of taskConfiguration from database
func (t *TaskConfiguration) GetTaskConfigurationByIdAndConfigKey(ctx context.Context, objId int, configKey string) (TaskConfiguration, error) {
	var (
		taskConfiguration TaskConfiguration
		err               error
		db                = mysql.FromContext(ctx).Model(&TaskConfiguration{})
	)

	curErr := db.Where("obj_id = ? and config_key = ?", objId, configKey).First(&taskConfiguration).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return taskConfiguration, err
}
