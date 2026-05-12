package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Logiclog struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	TaskID     int       `gorm:"column:task_id" json:"taskID"`         // 任务id
	TargetID   int       `gorm:"column:target_id" json:"targetID"`     // 目标id
	TargetURL  string    `gorm:"column:target_url" json:"targetURL"`   // 目标地址
	Status     int       `gorm:"column:status" json:"status"`          // 状态
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	StartTime  time.Time `gorm:"column:start_time" json:"startTime"`   // 开始时间
	EndTime    time.Time `gorm:"column:end_time" json:"endTime"`       // 结束时间
	IsAlive    int       `gorm:"column:is_alive" json:"isAlive"`       // 是否存活
}

// TableName sets insert table name for this struct type
func (l *Logiclog) TableName() string {
	return "logic_log"
}

// Get retrieves a list of logiclog from database
func (l *Logiclog) GetLogiclogList(ctx context.Context, taskId, page, limit int, search string) ([]Logiclog, int64, error) {
	var (
		logiclogList []Logiclog
		count        int64
		db           = mysql.FromContext(ctx).Model(&Logiclog{})
	)
	if search != "" {
		db = db.Where("target_url LIKE ?", "%"+search+"%")
	}
	db = db.Where("task_id = ?", taskId)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&logiclogList)
	db.Count(&count)

	return logiclogList, count, nil
}

// Get retrieves a single record of logiclog from database
func (l *Logiclog) GetLogiclog(ctx context.Context) (Logiclog, error) {
	var (
		logiclog Logiclog
		err      error
		db       = mysql.FromContext(ctx).Model(&Logiclog{})
	)

	curErr := db.Where("id = ?", l.ID).First(&logiclog).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return logiclog, err
}

// Add persists logiclog to database
func (l *Logiclog) AddLogiclog(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&Logiclog{})
	if err := db.Create(l).Error; err != nil {
		return l.ID, err
	}
	return l.ID, nil
}

// Update changes logiclog by id
func (l *Logiclog) UpdateLogiclog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logiclog{})

	if err := db.Where("id = ?", l.ID).Updates(l).Error; err != nil {
		return err
	}

	return nil
}

// Delete logiclog by id
func (l *Logiclog) DeleteLogiclog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logiclog{})
	if err := db.Where("id = ?", l.ID).Delete(l).Error; err != nil {
		return err
	}
	return nil
}

// Update changes logictarget by id
func (l *Logiclog) UpdateLogicLogParam(ctx context.Context, targetId int, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Logiclog{})
	if err := db.Where("target_id = ?", targetId).Updates(param).Error; err != nil {
		return err
	}
	return nil
}
