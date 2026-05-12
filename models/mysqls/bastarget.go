package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type BasTarget struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`     // 主键
	BasTaskID       int       `gorm:"column:bas_task_id" json:"basTaskID"` // bas任务ID
	Addr            string    `gorm:"column:addr" json:"addr"`             // 地址
	Status          int       `gorm:"column:status" json:"status"`
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`            // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`            // 修改时间
	BasTemplateID   int       `gorm:"column:bas_template_id" json:"basTemplateID"`     // 方案ID
	BasTemplateJSON string    `gorm:"column:bas_template_json" json:"basTemplateJSON"` // 方案数据json
}

// TableName sets insert table name for this struct type
func (b *BasTarget) TableName() string {
	return "bas_target"
}

// Get retrieves a list of basTaskTarget from database
func (b *BasTarget) GetBasTaskTargetList(ctx context.Context, taskId, page, limit int, search string) ([]BasTarget, int64, error) {
	var (
		basTaskTargetList []BasTarget
		count             int64
		db                = mysql.FromContext(ctx).Model(&BasTarget{})
	)

	db = db.Where("bas_task_id = ?", taskId)
	if search != "" {
		db = db.Where("addr like ?", "%"+search+"%")
	}

	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&basTaskTargetList)

	return basTaskTargetList, count, nil
}

// Get retrieves a single record of basTaskTarget from database
func (b *BasTarget) GetBasTaskTarget(ctx context.Context) (BasTarget, error) {
	var (
		basTaskTarget BasTarget
		err           error
		db            = mysql.FromContext(ctx).Model(&BasTarget{})
	)

	curErr := db.Where("id = ?", b.ID).First(&basTaskTarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return basTaskTarget, err
}

// Add persists basTaskTarget to database
func (b *BasTarget) AddBasTaskTarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// Update changes basTaskTarget by id
func (b *BasTarget) UpdateBasTaskTarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})

	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

// Delete basTaskTarget by id
func (b *BasTarget) DeleteBasTaskTarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})

	//b.Estate = "deleted"
	b.UpdateTime = time.Now()
	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

func (b *BasTarget) DeleteByTaskId(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})

	//b.Estate = "deleted"
	b.UpdateTime = time.Now()
	if err := db.Where("bas_task_id in ?", taskIds).Delete(b).Error; err != nil {
		return err
	}

	return nil
}
func (b *BasTarget) DeleteById(ctx context.Context, ids []int) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})

	//b.Estate = "deleted"
	b.UpdateTime = time.Now()
	if err := db.Where("id in ?", ids).Delete(b).Error; err != nil {
		return err
	}

	return nil
}

func (b *BasTarget) AddAll(ctx context.Context, datas []BasTarget) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})

	if err := db.Create(datas).Error; err != nil {
		return err
	}

	return nil
}

func (b *BasTarget) GetByTaskIdsAndStatus(ctx context.Context, taskIds any, status int) []BasTarget {
	var (
		basTaskList []BasTarget
		db          = mysql.FromContext(ctx).Model(&BasTarget{})
	)
	db.Where("bas_task_id in ? and status = ?", taskIds, status).Find(&basTaskList)
	return basTaskList
}

func (b *BasTarget) UpdateStatusByTaskId(ctx context.Context, taskId, status int) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})

	if err := db.Where("bas_task_id = ?", taskId).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

//根据主键更新数据
func (b *BasTarget) UpdateByIds(ctx context.Context, ids any, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})
	if err := db.Where("id IN ?", ids).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

func (b *BasTarget) UpdateStatusById(ctx context.Context, id, status int) error {
	var db = mysql.FromContext(ctx).Model(&BasTarget{})

	if err := db.Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (b *BasTarget) GetByTaskIdAndInStatus(ctx context.Context, taskId int, status []int) []BasTarget {
	var (
		basTaskList []BasTarget
		db          = mysql.FromContext(ctx).Model(&BasTarget{})
	)

	db.Where("bas_task_id = ? and status in ?", taskId, status).Find(&basTaskList)

	return basTaskList
}
