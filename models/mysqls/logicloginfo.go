package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Logicloginfo struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	TaskID     int       `gorm:"column:task_id" json:"taskID"`         // 任务id
	TargetID   int       `gorm:"column:target_id" json:"targetID"`     // 目标id
	LogID      int       `gorm:"column:log_id" json:"logID"`           // 日志id
	TargetURL  string    `gorm:"column:target_url" json:"targetURL"`   // 测试目标
	Pocname    string    `gorm:"column:pocname" json:"pocname"`        // 漏洞标识
	Result     string    `gorm:"column:result" json:"result"`          // 日志结果
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
}

// TableName sets insert table name for this struct type
func (l *Logicloginfo) TableName() string {
	return "logic_log_info"
}

// Get retrieves a list of logicloginfo from database
func (l *Logicloginfo) GetLogicloginfoList(ctx context.Context, page, limit int) ([]Logicloginfo, int64, error) {
	var (
		logicloginfoList []Logicloginfo
		count            int64
		db               = mysql.FromContext(ctx).Model(&Logicloginfo{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&logicloginfoList)
	db.Count(&count)

	return logicloginfoList, count, nil
}

// Get retrieves a single record of logicloginfo from database
func (l *Logicloginfo) GetLogicloginfo(ctx context.Context) (Logicloginfo, error) {
	var (
		logicloginfo Logicloginfo
		err          error
		db           = mysql.FromContext(ctx).Model(&Logicloginfo{})
	)

	curErr := db.Where("id = ?", l.ID).First(&logicloginfo).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return logicloginfo, err
}

// Add persists logicloginfo to database
func (l *Logicloginfo) AddLogicloginfo(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logicloginfo{})

	if err := db.Create(l).Error; err != nil {
		return err
	}

	return nil
}

// Update changes logicloginfo by id
func (l *Logicloginfo) UpdateLogicloginfo(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logicloginfo{})

	if err := db.Where("id = ?", l.ID).Updates(l).Error; err != nil {
		return err
	}

	return nil
}

// Delete logicloginfo by id
func (l *Logicloginfo) DeleteLogicloginfo(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logicloginfo{})

	if err := db.Where("id = ?", l.ID).Updates(l).Error; err != nil {
		return err
	}
	return nil
}

func (l *Logicloginfo) GetLogicloginfoListByLogId(ctx context.Context, logId int) ([]Logicloginfo, error) {
	var (
		logicloginfoList []Logicloginfo
		count            int64
		db               = mysql.FromContext(ctx).Model(&Logicloginfo{})
	)
	db.Where("log_id = ?", logId).Find(&logicloginfoList)
	db.Count(&count)
	return logicloginfoList, nil
}
