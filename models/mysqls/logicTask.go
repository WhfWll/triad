package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type LogicTask struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Name       string    `gorm:"column:name" json:"name"`              // 名称
	TargetUrl  string    `gorm:"column:target_url" json:"target_url"`  // 目标
	Status     int       `gorm:"column:status" json:"status"`          // 状态
	TargetNum  int       `gorm:"column:target_num" json:"targetNum"`   // 目标数量
	UserID     int       `gorm:"column:user_id" json:"userID"`         // 用户id
	ScanConfig string    `gorm:"column:scan_config" json:"scanConfig"` // 扫描配置
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 更新时间
	Type       int       `gorm:"column:type" json:"type"`              // 扫描类型
	Risk       int       `gorm:"column:risk" json:"risk"`              // 风险等级
	HighNum    int       `gorm:"column:high_num" json:"high_num"`      // 高危目标
	MiddleNum  int       `gorm:"column:middle_num" json:"middle_num"`  // 中危目标
	LowNum     int       `gorm:"column:low_num" json:"low_num"`        // 低危目标
	SafeNum    int       `gorm:"column:safe_num" json:"safe_num"`      // 安全目标
}

// TableName sets insert table name for this struct type
func (l *LogicTask) TableName() string {
	return "logic_task"
}

// Get retrieves a list of logicTask from database
func (l *LogicTask) GetLogicTaskList(ctx context.Context, page, limit int, search string) ([]LogicTask, int64, error) {
	var (
		logicTaskList []LogicTask
		count         int64
		db            = mysql.FromContext(ctx).Model(&LogicTask{})
	)
	if search != "" {
		db = db.Where("name like ?", "%"+search+"%").Or("target_url LIKE ?", "%"+search+"%")
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("create_time desc").Find(&logicTaskList)
	return logicTaskList, count, nil
}

// Get retrieves a single record of logicTask from database
func (l *LogicTask) GetLogicTask(ctx context.Context) (LogicTask, error) {
	var (
		logicTask LogicTask
		err       error
		db        = mysql.FromContext(ctx).Model(&LogicTask{})
	)
	curErr := db.Where("id = ?", l.ID).First(&logicTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return logicTask, err
}

// Add persists logicTask to database
func (l *LogicTask) AddLogicTask(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&LogicTask{})
	if err := db.Create(l).Error; err != nil {
		return l.ID, err
	}
	return l.ID, nil
}

// Update changes logicTask by id
func (l *LogicTask) UpdateLogicTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&LogicTask{})

	if err := db.Where("id = ?", l.ID).Updates(l).Error; err != nil {
		return err
	}

	return nil
}

// Delete logicTask by id
func (l *LogicTask) DeleteLogicTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&LogicTask{})

	l.UpdateTime = time.Now()
	if err := db.Where("id = ?", l.ID).Delete(l).Error; err != nil {
		return err
	}

	return nil
}

// 根据id修改逻辑任务信息
func (l *LogicTask) UpdateLogicTaskById(ctx context.Context, id int, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&LogicTask{})
	if err := db.Where("id = ?", id).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// Get retrieves a single record of logicTask from database
func (l *LogicTask) GetLogicTaskListByStatus(ctx context.Context, status int) ([]LogicTask, error) {
	var (
		logicTaskList []LogicTask
		err           error
		db            = mysql.FromContext(ctx).Model(&LogicTask{})
	)
	curErr := db.Where("status = ?", status).Find(&logicTaskList).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return logicTaskList, err
}
