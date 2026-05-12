package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"smart/tools/utils"
	"time"
)

type TaskTaskInfo struct {
	ID               int       `gorm:"column:id;primary_key" json:"id"`  // 主键
	TaskID           int       `gorm:"column:task_id" json:"taskID"`     // 所属任务id
	TaskName         string    `gorm:"column:task_name" json:"taskName"` // 任务名称
	Status           int       `gorm:"column:status;type:int(11);default:0;comment:任务运行状态" json:"status"`
	Weight           int       `gorm:"column:weight" json:"weight"`                       //优先权重
	TaskType         int       `gorm:"column:task_type" json:"taskType"`                  // 任务类型
	CheckTarget      string    `gorm:"column:check_target" json:"checkTarget"`            // 用户填写的检测目标
	ExecuteType      int       `gorm:"column:execute_type" json:"executeType"`            // 执行方法 1即时执行、2定时执行 、3周期执行、4监控执行
	ExecuteLastTime  time.Time `gorm:"column:execute_last_time" json:"executeLastTime"`   // 任务上次执行时间
	ExecuteNextTime  time.Time `gorm:"column:execute_next_time" json:"executeNextTime"`   // 任务下次执行时间
	ExecuteJSON      string    `gorm:"column:execute_json" json:"executeJSON"`            // 执行方式具体参数
	TaskTemplateID   int       `gorm:"column:task_template_id" json:"taskTemplateID"`     // 所选择的任务场景id
	TaskTemplateJSON string    `gorm:"column:task_template_json" json:"taskTemplateJSON"` // 任务场景参数
	Overview         string    `gorm:"column:overview" json:"overview"`                   // 概览统计
	UserID           int       `gorm:"column:user_id" json:"userID"`                      // 提交者id
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`              // 发生时间
	UpdateTime       time.Time `gorm:"column:update_time" json:"updateTime"`              // 发生时间
}

// TableName sets insert table name for this struct type
func (t *TaskTaskInfo) TableName() string {
	return "task_task_info"
}

// Get retrieves a list of taskTaskInfo from database
func (t *TaskTaskInfo) GetTaskTaskInfoList(ctx context.Context, page, limit int) ([]TaskTaskInfo, int64, error) {
	var (
		taskTaskInfoList []TaskTaskInfo
		count            int64
		db               = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&taskTaskInfoList)
	db.Count(&count)

	return taskTaskInfoList, count, nil
}

// GetTaskTaskInfoByTaskId 根据taskId查询任务详情信息
func (t *TaskTaskInfo) GetTaskTaskInfoByTaskId(ctx context.Context, taskId int) (TaskTaskInfo, error) {
	var (
		taskTaskInfo TaskTaskInfo
		err          error
		db           = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	)
	curErr := db.Where("task_id = ?", taskId).First(&taskTaskInfo).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskTaskInfo, err
}

// Add persists taskTaskInfo to database
func (t *TaskTaskInfo) AddTaskTaskInfo(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskInfo{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// 修改任务状态
func (t *TaskTaskInfo) UpdateStatusByTaskId(ctx context.Context, taskId, status int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskInfo{})

	data := make(map[string]interface{})
	data["status"] = status
	if status == enums.TaskStatusRunning {
		data["execute_last_time"] = time.Now().Format(utils.DateTime)
	}
	data["update_time"] = time.Now()
	if err := db.Where("task_id = ?", taskId).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTaskInfoByTaskIds 根据任务id更新任务详情信息 code opt???
func (t *TaskTaskInfo) UpdateTaskInfoByTaskIds(ctx context.Context, taskId any, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	if err := db.Where("task_id IN ?", taskId).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// Delete taskTaskInfo by id
func (t *TaskTaskInfo) DeleteTaskTaskInfo(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskInfo{})

	//t.Estate = "deleted"
	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTaskTaskInfoOverview 更新统计信息
func (t *TaskTaskInfo) UpdateTaskTaskInfoOverview(ctx context.Context, taskId int, overview string) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	if err := db.Where("task_id = ?", taskId).Update("overview", overview).Error; err != nil {
		return err
	}
	return nil
}

// 批量删除 通过task_ids
func (t *TaskTaskInfo) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	if err := db.Where("task_id in ?", taskIds).Delete(t).Error; err != nil {
		return err
	}
	return nil
}

// 批量获取 通过task_ids
func (t *TaskTaskInfo) GetByTaskIds(ctx context.Context, taskIds []int) []TaskTaskInfo {
	var (
		data []TaskTaskInfo
		db   = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	)

	db.Where("task_id in ?", taskIds).Find(&data)
	return data
}

// GetOneWaitTiming 依据允许运行的时间获取一个定时任务
func (t *TaskTaskInfo) GetOneWaitTiming(ctx context.Context, runDate string) (TaskTaskInfo, error) {
	var (
		taskCheckTask TaskTaskInfo
		err           error
		db            = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	)
	curErr := db.
		Where("status = ?", enums.TaskStatusBegin).
		Where("execute_type = ?", enums.TaskExecTypeTiming).
		Where("execute_next_time <= ?", runDate).Order("weight desc,id").
		First(&taskCheckTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskCheckTask, err
}

// GetOneWaitCycle 依据允许运行的时间获取一个周期任务
func (t *TaskTaskInfo) GetOneWaitCycle(ctx context.Context, runDate string) (TaskTaskInfo, error) {
	var (
		taskCheckTask TaskTaskInfo
		err           error
		db            = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	)
	curErr := db.
		Where("status = ?", enums.TaskStatusBegin).
		Where("execute_type = ?", enums.TaskExecTypeCycle).
		Where("execute_next_time <= ?", runDate).Order("weight desc,id").
		First(&taskCheckTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskCheckTask, err
}

// Update changes taskTaskInfo by id
func (t *TaskTaskInfo) UpdateExecNextTimeByTaskId(ctx context.Context, taskId int, execNextTime string) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskInfo{})

	if err := db.Where("task_id = ?", taskId).Update("execute_next_time", execNextTime).Error; err != nil {
		return err
	}

	return nil
}

// 通过任务状态获取任务
func (t *TaskTaskInfo) GetTaskByTaskStatus(ctx context.Context, taskStatus int) []TaskTaskInfo {
	var (
		taskTask []TaskTaskInfo
		db       = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	)

	db.Where("status = ?", taskStatus).Find(&taskTask)
	return taskTask
}

// All 获取所有目标
func (t *TaskTaskInfo) All(ctx context.Context, filter string) ([]TaskTaskInfo, error) {
	var (
		taskInfoList []TaskTaskInfo
		db           = mysql.FromContext(ctx).Model(&TaskTaskInfo{})
	)
	if filter != "" {
		db.Where(filter).Find(&taskInfoList)
	} else {
		db.Find(&taskInfoList)
	}
	return taskInfoList, nil
}
