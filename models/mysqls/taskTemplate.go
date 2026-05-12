package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"
)

type TaskTemplate struct {
	ID           int       `gorm:"column:id;primary_key" json:"id"`
	Estate       string    `gorm:"column:estate" json:"estate"`
	TemplateName string    `gorm:"column:template_name" json:"templateName"`
	Describe     string    `gorm:"column:describe" json:"describe"`
	IsDefault    uint      `gorm:"column:is_default" json:"isDefault"`
	UserID       int       `gorm:"column:user_id" json:"userId"`
	Source       uint      `gorm:"column:source" json:"source"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName sets insert table name for this struct type
func (t *TaskTemplate) TableName() string {
	return "task_template"
}

// Get retrieves a list of taskTemplate from database
func (t *TaskTemplate) GetTaskTemplateList(ctx context.Context, page, limit int, templateName string, source []int, Estate string) ([]TaskTemplate, int64, error) {
	var (
		taskTemplateList []TaskTemplate
		count            int64
		db               = mysql.FromContext(ctx).Model(&TaskTemplate{})
	)
	if templateName != "" {
		db = db.Where("template_name like ?", "%"+templateName+"%")
	}

	if Estate != "" {
		db = db.Where("estate = ?", Estate)
	}

	if len(source) > 0 {
		sourceStr := make([]string, 0)
		for _, item := range source {
			sourceStr = append(sourceStr, strconv.Itoa(item))
		}
		db = db.Where("source in (?)", strings.Join(sourceStr, ","))
	}

	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id DESC").Find(&taskTemplateList)

	return taskTemplateList, count, nil
}

// Get retrieves a single record of taskTemplate from database
func (t *TaskTemplate) GetTaskTemplate(ctx context.Context, id int, estate string) (TaskTemplate, error) {
	var (
		taskTemplate TaskTemplate
		err          error
		db           = mysql.FromContext(ctx).Model(&TaskTemplate{})
	)
	if estate != "" {
		db = db.Where("estate = ?", estate)
	}
	curErr := db.Where("id = ?", id).First(&taskTemplate).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskTemplate, err
}

func (t *TaskTemplate) AllTaskTemplateByIds(ctx context.Context, ids []int) ([]TaskTemplate, error) {
	var (
		taskTemplate []TaskTemplate
		err          error
		db           = mysql.FromContext(ctx).Model(&TaskTemplate{})
	)

	curErr := db.Where("id in ?", ids).Find(&taskTemplate).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return taskTemplate, err
}

// Add persists taskTemplate to database
func (t *TaskTemplate) AddTaskTemplate(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTemplate{})

	t.CreateTime = time.Now()
	t.UpdateTime = time.Now()
	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// Update changes taskTemplate by id
func (t *TaskTemplate) UpdateTaskTemplate(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTemplate{})

	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// Delete taskTemplate by id
func (t *TaskTemplate) DeleteTaskTemplate(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTemplate{})

	//t.Estate = "deleted"
	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

func (t *TaskTemplate) AllTaskTemplate(ctx context.Context, field string, estate string) ([]TaskTemplate, int64, error) {
	var (
		taskTemplateList []TaskTemplate
		count            int64
		db               = mysql.FromContext(ctx).Model(&TaskTemplate{})
	)

	if field == "" {
		field = "*"
	}

	if estate != "" {
		db = db.Where("estate = ?", estate)
	}

	db.Select(field).Order("id DESC").Find(&taskTemplateList)
	db.Count(&count)

	return taskTemplateList, count, nil
}

func (t *TaskTemplate) GetByName(ctx context.Context, templateName string) (TaskTemplate, error) {
	var (
		taskTemplateList TaskTemplate
		db               = mysql.FromContext(ctx).Model(&TaskTemplate{})
	)

	db.Where("template_name = ?", templateName).First(&taskTemplateList)

	return taskTemplateList, nil
}

func (t *TaskTemplate) CancelDefault(ctx context.Context) error {
	db := mysql.FromContext(ctx).Model(&TaskTemplate{})
	if err := db.Where("1 = 1").Update("is_default", enums.TaskTemplateIsDefaultN).Error; err != nil {
		return err
	}
	return nil
}

func (t *TaskTemplate) SetDefault(ctx context.Context, templateId int) error {
	db := mysql.FromContext(ctx).Model(&TaskTemplate{})
	if err := db.Where("id = ?", templateId).Update("is_default", enums.TaskTemplateIsDefaultY).Error; err != nil {
		return err
	}
	return nil
}

func (t *TaskTemplate) UpdateStatusByIds(ctx context.Context, ids []int, estate string) error {
	db := mysql.FromContext(ctx).Model(&TaskTemplate{})
	if err := db.Where("id in ?", ids).Update("estate", estate).Error; err != nil {
		return err
	}
	return nil
}

// GetTaskTemplateById 根据id查询模板信息（包括软删除）
func (t *TaskTemplate) GetTaskTemplateById(ctx context.Context, id int) (TaskTemplate, error) {
	var (
		taskTemplate TaskTemplate
		err          error
		db           = mysql.FromContext(ctx).Model(&TaskTemplate{})
	)
	curErr := db.Where("id = ?", id).First(&taskTemplate).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskTemplate, err
}

// GetTaskSceneCount 获取任务场景数量
func (t *TaskTemplate) GetTaskSceneCount(ctx context.Context, estate string) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&TaskTemplate{})
	)
	db.Where("estate", estate)
	db.Count(&count)
	return count
}
