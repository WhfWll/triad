package mysqls

import (
	"context"
	"fmt"
	"smart/tools/enums"
	"smart/tools/utils"
	"time"

	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type Tasklog struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	TaskID     int       `gorm:"column:task_id" json:"taskID"`         // 所属任务id
	TargetID   int       `gorm:"column:target_id" json:"targetID"`     // 所属目标id
	TargetURL  string    `gorm:"column:target_url" json:"targetURL"`   // 测试目标地址
	Status     int       `gorm:"column:status" json:"status"`          // 状态
	IsAlive    int       `gorm:"column:is_alive" json:"isAlive"`       //是否存活
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	StartTime  time.Time `gorm:"column:start_time" json:"startTime"`   // 开始时间
	EndTime    time.Time `gorm:"column:end_time" json:"endTime"`       // 结束时间
}

// TableName sets insert table name for this struct type
func (t *Tasklog) TableName() string {
	return "task_log"
}

// GetTasklogList 测试日志列表及筛选
func (t *Tasklog) GetTasklogList(ctx context.Context, taskId int, search string, page, limit int) ([]Tasklog, int64, error) {
	var (
		tasklogList []Tasklog
		count       int64
		db          = mysql.FromContext(ctx).Model(&Tasklog{})
		query       string
		args        []interface{}
	)
	query += "task_id = ?"
	args = append(args, taskId)
	if len(search) > 0 {
		query += " and target_url LIKE ?"
		args = append(args, "%"+search+"%")
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&tasklogList)
	return tasklogList, count, nil
}

// GetTasklogListByTargetIds 根据目标ids获取测试日志
// 参数targetIds必须是一个[]int或者[]string
func (t *Tasklog) GetTasklogListByTargetIds(ctx context.Context, targetIds interface{}) ([]Tasklog, error) {
	var (
		tasklogList []Tasklog
		db          = mysql.FromContext(ctx).Model(&Tasklog{})
	)
	db.Where("target_id IN ?", targetIds).Find(&tasklogList)
	return tasklogList, nil
}

// Get retrieves a single record of tasklog from database
func (t *Tasklog) GetTasklog(ctx context.Context) (Tasklog, error) {
	var (
		tasklog Tasklog
		err     error
		db      = mysql.FromContext(ctx).Model(&Tasklog{})
	)

	curErr := db.Where("id = ?", t.ID).First(&tasklog).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return tasklog, err
}

// Add persists tasklog to database
func (t *Tasklog) AddTasklog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Tasklog{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTasklogByTargetIds 根据TargetIds删除数据
// targetIds必须为[]int/[]string
func (t *Tasklog) DeleteTasklogByTargetIds(ctx context.Context, targetIds any) error {
	var db = mysql.FromContext(ctx).Model(&Tasklog{})
	if err := db.Where("target_id IN ?", targetIds).Delete(t).Error; err != nil {
		return err
	}
	return nil
}

// Get retrieves a single record of tasklog from database
func (t *Tasklog) GetTaskLogByTargetId(ctx context.Context) (Tasklog, error) {
	var (
		tasklog Tasklog
		err     error
		db      = mysql.FromContext(ctx).Model(&Tasklog{})
	)

	curErr := db.Where("target_id = ?", t.TargetID).First(&tasklog).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return tasklog, err
}

// 批量删除 通过task_ids
func (t *Tasklog) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&Tasklog{})
	if err := db.Where("task_id in ?", taskIds).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTaskLogStatus 更新任务日志状态
func (t *Tasklog) UpdateTaskLogStatus(ctx context.Context, status int) error {
	var db = mysql.FromContext(ctx).Model(&Tasklog{})
	var data = map[string]interface{}{
		"status": status,
	}
	if status == enums.TaskStatusFinish {
		data["end_time"] = time.Now()
	}
	if err := db.Where("id = ?", t.ID).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTaskLogStateByTargetId 更新任务日志状态
func (t *Tasklog) UpdateTaskLogStateByTargetId(ctx context.Context, targetId, status int) error {
	var db = mysql.FromContext(ctx).Model(&Tasklog{})
	var data = map[string]interface{}{
		"status": status,
	}
	if status == enums.TaskStatusFinish {
		data["end_time"] = time.Now()
	}
	if err := db.Where("target_id = ?", targetId).Updates(data).Debug().Error; err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// getTargetTimeout 从 scanner 配置读取目标超时时间，默认 2 小时
func getTargetTimeout() time.Duration {
	var cfg struct {
		ProcessTimeoutSeconds int `json:"process_timeout_seconds"`
	}
	if err := config.Load("scanner", &cfg); err != nil {
		return 2 * time.Hour
	}
	if cfg.ProcessTimeoutSeconds <= 0 {
		return 2 * time.Hour
	}
	return time.Duration(cfg.ProcessTimeoutSeconds) * time.Second
}

// GetTimeoutTargetIds 通过状态获取目标数量
func (t *Tasklog) GetTimeoutTargetIds(ctx context.Context) ([]int, error) {
	var db = mysql.FromContext(ctx).Model(&Tasklog{})
	overTime := time.Now().Add(-getTargetTimeout()).Format(utils.DateTime)
	targetIds := make([]int, 0)
	err := db.Select("target_id").Where("status = ? and create_time < ?", enums.TargetStatusRunning, overTime).Find(&targetIds).Error
	if err != nil {
		return targetIds, err
	}
	return targetIds, nil
}

// HandleTimeoutTaskLog 处理超时的检测目标
func (t *Tasklog) HandleTimeoutTaskLog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Tasklog{})
	overTime := time.Now().Add(-getTargetTimeout()).Format(utils.DateTime)
	err := db.Where("status = ? and create_time < ?", enums.TargetStatusRunning, overTime).Updates(Tasklog{Status: enums.TargetStatusFinish}).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateTaskLogAliveByTargetId 更新任务日志存活状态
func (t *Tasklog) UpdateTaskLogAliveByTargetId(ctx context.Context, targetId, isAlive int) error {
	var db = mysql.FromContext(ctx).Model(&Tasklog{})
	var data = map[string]interface{}{
		"is_alive": isAlive,
	}
	if err := db.Where("target_id = ?", targetId).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// All 获取所有目标
func (t *Tasklog) All(ctx context.Context, filter string) ([]Tasklog, error) {
	var (
		taskLogList []Tasklog
		db          = mysql.FromContext(ctx).Model(&Tasklog{})
	)
	if filter != "" {
		db.Where(filter).Find(&taskLogList)
	} else {
		db.Find(&taskLogList)
	}
	return taskLogList, nil
}

// Get retrieves a single record of tasklog from database
func (t *Tasklog) GetFirstFiveTaskLog(ctx context.Context) ([]Tasklog, error) {
	var (
		taskLogList []Tasklog
		err         error
		db          = mysql.FromContext(ctx).Model(&Tasklog{})
	)
	curErr := db.Order("id desc").Limit(5).Find(&taskLogList).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskLogList, err
}
