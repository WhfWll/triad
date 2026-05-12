package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Logictarget struct {
	ID         int       `gorm:"column:id" json:"id"`                  // 主键
	TargetURL  string    `gorm:"column:target_url" json:"targetURL"`   // 目标地址
	TaskID     int       `gorm:"column:task_id" json:"taskID"`         // 任务id
	Status     int       `gorm:"column:status" json:"status"`          // 目标状态
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 更新时间
	Type       int       `gorm:"column:type" json:"type"`              // 逻辑漏洞测试类型
	ConfigJson string    `gorm:"column:configJson" json:"configJson"`  // 扫描配置信息
	Risk       int       `gorm:"column:risk" json:"risk"`              // 风险等级
	IsAlive    int       `gorm:"column:is_alive" json:"isAlive"`       // 是否存活
}

// TableName sets insert table name for this struct type
func (l *Logictarget) TableName() string {
	return "logic_target"
}

// Get retrieves a list of logictarget from database
func (l *Logictarget) GetLogictargetList(ctx context.Context, taskId, page, limit int, search string) ([]Logictarget, int64, error) {
	var (
		logictargetList []Logictarget
		count           int64
		db              = mysql.FromContext(ctx).Model(&Logictarget{})
	)
	if search != "" {
		db = db.Where("target_url LIKE ?", "%"+search+"%")
	}
	db.Where("task_id = ?", taskId).Limit(limit).Offset(limit * (page - 1)).Find(&logictargetList)
	db.Count(&count)

	return logictargetList, count, nil
}

// Get retrieves a single record of logictarget from database
func (l *Logictarget) GetLogictarget(ctx context.Context) (Logictarget, error) {
	var (
		logictarget Logictarget
		err         error
		db          = mysql.FromContext(ctx).Model(&Logictarget{})
	)

	curErr := db.Where("id = ?", l.ID).First(&logictarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return logictarget, err
}

// Add persists logictarget to database
func (l *Logictarget) AddLogictarget(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&Logictarget{})
	if err := db.Create(l).Error; err != nil {
		return 0, err
	}
	return l.ID, nil
}

// Update changes logictarget by id
func (l *Logictarget) UpdateLogictarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logictarget{})
	if err := db.Where("id = ?", l.ID).Updates(l).Error; err != nil {
		return err
	}
	return nil
}

// Delete logictarget by id
func (l *Logictarget) DeleteLogictarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logictarget{})

	l.UpdateTime = time.Now()
	if err := db.Where("id = ?", l.ID).Updates(l).Error; err != nil {
		return err
	}

	return nil
}

// Get retrieves a list of logictarget from database
func (l *Logictarget) GetLogicTargetListByTargetIdAndStatus(ctx context.Context, taskIdList []int, status int) ([]Logictarget, error) {
	var (
		logictargetList []Logictarget
		db              = mysql.FromContext(ctx).Model(&Logictarget{})
	)
	db.Where("task_id in ? and status = ?", taskIdList, status).Find(&logictargetList)
	return logictargetList, nil
}

// Update changes logictarget by id
func (l *Logictarget) UpdateLogicTargetParam(ctx context.Context, targetId int, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Logictarget{})
	if err := db.Where("id = ?", targetId).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// Get retrieves a list of logictarget from database
func (l *Logictarget) GetLogicTargetListByStatus(ctx context.Context, status int) ([]Logictarget, error) {
	var (
		logictargetList []Logictarget
		db              = mysql.FromContext(ctx).Model(&Logictarget{})
	)
	db.Where("status = ?", status).Find(&logictargetList)
	return logictargetList, nil
}

// Get retrieves a list of logictarget from database
func (l *Logictarget) GetLogicTargetListByTaskIdAndStatus(ctx context.Context, taskId int, status []int) ([]Logictarget, error) {
	var (
		logictargetList []Logictarget
		db              = mysql.FromContext(ctx).Model(&Logictarget{})
	)
	db.Where("task_id = ? and status in ?", taskId, status).Find(&logictargetList)
	return logictargetList, nil
}

// Get retrieves a list of logictarget from database
func (l *Logictarget) UpdateTargetListByTaskId(ctx context.Context, taskId, status int) error {
	var db = mysql.FromContext(ctx).Model(&Logictarget{})
	db.Where("task_id = ?", taskId).Update("status", status)
	return nil
}
