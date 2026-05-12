package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Flowtarget struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`       // 主键
	FlowTaskID int       `gorm:"column:flow_task_id" json:"flowTaskID"` // 所属流量分析任务id
	TargetURL  string    `gorm:"column:target_url" json:"targetURL"`    // 测试目标地址
	Status     int       `gorm:"column:status" json:"status"`           // 状态，1-禁用，2-正常
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`  // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`  // 修改时间
}

// TableName sets insert table name for this struct type
func (f *Flowtarget) TableName() string {
	return "flow_target"
}

// Get retrieves a list of flowtarget from database
func (f *Flowtarget) GetFlowtargetList(ctx context.Context, flowTaskId int) ([]Flowtarget, error) {
	var (
		flowtargetList []Flowtarget
		db             = mysql.FromContext(ctx).Model(&Flowtarget{})
	)
	db.Where("flow_task_id = ?", flowTaskId).Find(&flowtargetList)
	return flowtargetList, nil
}

// Get retrieves a single record of flowtarget from database
func (f *Flowtarget) GetFlowtarget(ctx context.Context) (Flowtarget, error) {
	var (
		flowtarget Flowtarget
		err        error
		db         = mysql.FromContext(ctx).Model(&Flowtarget{})
	)

	curErr := db.Where("id = ?", f.ID).First(&flowtarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return flowtarget, err
}

// Add persists flowtarget to database
func (f *Flowtarget) AddFlowtarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Flowtarget{})

	if err := db.Create(f).Error; err != nil {
		return err
	}

	return nil
}

// Update changes flowtarget by id
func (f *Flowtarget) UpdateFlowtarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Flowtarget{})

	if err := db.Where("id = ?", f.ID).Updates(f).Error; err != nil {
		return err
	}

	return nil
}

// DeleteFlowtskByIds 批量删除
func (f *Flowtarget) DeleteFlowtargetByTaskIds(ctx context.Context, flowTaskIds any) error {
	var db = mysql.FromContext(ctx).Model(&Flowtarget{})
	if err := db.Where("flow_task_id in ?", flowTaskIds).Delete(f).Error; err != nil {
		return err
	}
	return nil
}

// AddManyFlowtarget 批量新增
func (f *Flowtarget) AddManyFlowtarget(ctx context.Context, data []Flowtarget) error {
	var db = mysql.FromContext(ctx).Model(&Flowtarget{})
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
