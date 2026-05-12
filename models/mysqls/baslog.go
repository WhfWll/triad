package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type BasLog struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`         // 主键
	BasTaskID   int       `gorm:"column:bas_task_id" json:"basTaskID"`     // bas任务ID
	BasTargetId int       `gorm:"column:bas_target_id" json:"BasTargetId"` // bas任务目标agent ID
	Content     string    `gorm:"column:content" json:"content"`           // 日志内容
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`    // 创建时间
}

// TableName sets insert table name for this struct type
func (b *BasLog) TableName() string {
	return "bas_log"
}

// Get retrieves a list of basTaskTargetLog from database
func (b *BasLog) GetBasTaskTargetLogList(ctx context.Context, page, limit int) ([]BasLog, int64, error) {
	var (
		basTaskTargetLogList []BasLog
		count                int64
		db                   = mysql.FromContext(ctx).Model(&BasLog{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&basTaskTargetLogList)
	db.Count(&count)

	return basTaskTargetLogList, count, nil
}

// Get retrieves a single record of basTaskTargetLog from database
func (b *BasLog) GetBasTaskTargetLog(ctx context.Context) (BasLog, error) {
	var (
		basTaskTargetLog BasLog
		err              error
		db               = mysql.FromContext(ctx).Model(&BasLog{})
	)

	curErr := db.Where("id = ?", b.ID).First(&basTaskTargetLog).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return basTaskTargetLog, err
}

// 新增一条bas日志
func (b *BasLog) AddBasLog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasLog{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// Update changes basTaskTargetLog by id
func (b *BasLog) UpdateBasTaskTargetLog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasLog{})

	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

func (b *BasLog) AddAll(ctx context.Context, datas []BasLog) error {
	var db = mysql.FromContext(ctx).Model(&BasLog{})

	if err := db.Create(datas).Error; err != nil {
		return err
	}

	return nil
}

func (b *BasLog) AllByTargetId(ctx context.Context, targetId int) []BasLog {
	var db = mysql.FromContext(ctx).Model(&BasLog{})
	var basTaskTargetLog []BasLog
	db.Where("bas_target_id = ?", targetId).Find(&basTaskTargetLog)
	return basTaskTargetLog
}

func (b *BasLog) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&BasLog{})
	if err := db.Where("bas_task_id in ?", taskIds).Delete(&b).Error; err != nil {
		return err
	}
	return nil
}

func (b *BasLog) DeleteByTargetIds(ctx context.Context, targetIds []int) error {
	var db = mysql.FromContext(ctx).Model(&BasLog{})
	if err := db.Where("bas_target_id in ?", targetIds).Delete(&b).Error; err != nil {
		return err
	}
	return nil
}
