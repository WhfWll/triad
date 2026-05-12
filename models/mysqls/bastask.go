package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type BasTask struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`                 // 主键
	Name            string    `gorm:"column:name" json:"name"`                         // 名称
	BasTemplateID   int       `gorm:"column:bas_template_id" json:"basTemplateID"`     // 方案ID
	BasTemplateJSON string    `gorm:"column:bas_template_json" json:"basTemplateJSON"` // 方案数据json
	RiskLevel       int       `gorm:"column:risk_level" json:"riskLevel"`              // 任务风险等级，1-高危、2-中危、3-低危、4-安全
	Status          int       `gorm:"column:status" json:"status"`
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
	User            int       `gorm:"column:user" json:"user"`              // 创建者
}

// TableName sets insert table name for this struct type
func (b *BasTask) TableName() string {
	return "bas_task"
}

// Get retrieves a list of basTask from database
func (b *BasTask) GetBasTaskList(ctx context.Context, page, limit, riskLevel int, search string) ([]BasTask, int64, error) {
	var (
		basTaskList []BasTask
		count       int64
		db          = mysql.FromContext(ctx).Model(&BasTask{})
	)
	if search != "" {
		db = db.Where("name like ?", "%"+search+"%")
	}
	if riskLevel != 0 {
		db = db.Where("risk_level = ?", riskLevel)
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id DESC").Find(&basTaskList)

	return basTaskList, count, nil
}

//根据id查询一条数据
func (b *BasTask) GetBasTaskById(ctx context.Context, id int) (BasTask, error) {
	var (
		basTask BasTask
		err     error
		db      = mysql.FromContext(ctx).Model(&BasTask{})
	)
	curErr := db.Where("id = ?", id).First(&basTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return basTask, err
}

// 根据状态获取bas任务数据
func (b *BasTask) GetBasTaskByStatus(ctx context.Context, status int) ([]BasTask, error) {
	var (
		basTaskList []BasTask
		db          = mysql.FromContext(ctx).Model(&BasTask{})
	)
	db.Where("status = ?", status).Order("id DESC").Find(&basTaskList)
	return basTaskList, nil
}

// Add persists basTask to database
func (b *BasTask) AddBasTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasTask{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// Update changes basTask by id
func (b *BasTask) UpdateBasTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasTask{})

	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

// Delete basTask by id
func (b *BasTask) DeleteById(ctx context.Context, ids []int) error {
	var db = mysql.FromContext(ctx).Model(&BasTask{})

	//b.Estate = "deleted"
	b.UpdateTime = time.Now()
	if err := db.Where("id in ?", ids).Delete(b).Error; err != nil {
		return err
	}

	return nil
}

func (b *BasTask) UpdateStatusById(ctx context.Context, id, status int) error {
	var db = mysql.FromContext(ctx).Model(&BasTask{})

	if err := db.Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// 根据主键更新数据
func (b *BasTask) UpdateByIds(ctx context.Context, ids any, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&BasTask{})
	if err := db.Where("id IN ?", ids).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

func (b *BasTask) UpdateById(ctx context.Context, id int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&BasTask{})
	if err := db.Where("id = ?", id).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

func (b *BasTask) GetByIds(ctx context.Context, ids []int) []BasTask {
	var (
		basTask []BasTask
		db      = mysql.FromContext(ctx).Model(&BasTask{})
	)

	db.Where("id in ?", ids).Find(&basTask)

	return basTask
}
