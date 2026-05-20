package mysqls

import (
	"context"
	"fmt"
	"smart/tools/enums"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type TaskTask struct {
	ID             int       `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
	TaskName       string    `gorm:"column:task_name;type:varchar(50);comment:任务名称;NOT NULL" json:"task_name"`
	RiskLevel      int       `gorm:"column:risk_level;type:int(11);default:0;comment:该任务风险等级;NOT NULL" json:"risk_level"`
	Status         int       `gorm:"column:status;type:int(11);default:0;comment:任务运行状态" json:"status"`
	Weight         int       `gorm:"column:weight" json:"weight"` //优先权重
	IsStats        int       `gorm:"column:is_stats;type:tinyint(4);default:1;comment:是否需要重新统计，1否，2是" json:"is_stats"`
	TaskType       int       `gorm:"column:task_type;type:int(11);default:0;comment:任务类型;NOT NULL" json:"task_type"`
	ExecuteType    int       `gorm:"column:execute_type;type:tinyint(4);default:1;comment:执行方法:1-即时执行,2-定时执行,3-周期执行,4-监控执行;NOT NULL" json:"execute_type"`
	TaskTemplateID int       `gorm:"column:task_template_id;type:int(11);default:0;comment:所选择的任务场景id" json:"task_template_id"`
	TargetNum      int       `gorm:"column:target_num;type:int(11);default:0;comment:该任务下目标数量;NOT NULL" json:"target_num"`
	HigeNum        int       `gorm:"column:hige_num;type:int(11);default:0;comment:高危目标个数" json:"hige_num"`
	MiddleNum      int       `gorm:"column:middle_num;type:int(11);default:0;comment:中危目标个数" json:"middle_num"`
	LowNum         int       `gorm:"column:low_num;type:int(11);default:0;comment:低危目标个数" json:"low_num"`
	SafeNum        int       `gorm:"column:safe_num;type:int(11);default:0;comment:安全目标个数" json:"safe_num"`
	UserID         int       `gorm:"column:user_id;type:int(11);default:0;comment:提交者id;NOT NULL" json:"user_id"`
	Pid            int       `gorm:"column:pid;type:int(11);default:0;comment:任务父级任务ID，用于任务验证测试;NOT NULL" json:"pid"`
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;default:1970-01-01 08:00:01;comment:发生时间;NOT NULL" json:"create_time"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;default:1970-01-01 08:00:01;comment:发生时间;NOT NULL" json:"update_time"`
}

// TableName sets insert table name for this struct type
func (t *TaskTask) TableName() string {
	return "task_task"
}

// Get retrieves a list of taskCheckTask from database
func (t *TaskTask) GetTaskCheckTaskList(ctx context.Context, page, limit, riskLevel int, taskName string, startTime, endTime string, taskIdList, userIdList []int) ([]TaskTask, int64, error) {
	var (
		taskCheckTaskList []TaskTask
		count             int64
		db                = mysql.FromContext(ctx).Model(&TaskTask{})
	)

	// 筛选任务名称
	if taskName != "" {
		db = db.Where("task_name like ?", "%"+taskName+"%")
	}
	// 筛选风险等级
	if riskLevel != 0 {
		db = db.Where("risk_level = ?", riskLevel)
	}
	// 任务开始时间与结束时间筛选
	if startTime != "" && endTime != "" {
		db = db.Where("update_time > ?", startTime).Where("update_time <= ?", endTime)
	}
	// 任务id筛选
	if len(taskIdList) != 0 {
		db = db.Where("id in ?", taskIdList)
	}
	// 用户id筛选
	if len(userIdList) != 0 {
		db = db.Where("user_id in ?", userIdList)
	}

	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("create_time desc").Find(&taskCheckTaskList)

	return taskCheckTaskList, count, nil
}

// Get retrieves a single record of taskCheckTask from database
func (t *TaskTask) GetTaskCheckTask(ctx context.Context, taskId int) (TaskTask, error) {
	var (
		taskCheckTask TaskTask
		err           error
		db            = mysql.FromContext(ctx).Model(&TaskTask{})
	)
	curErr := db.Where("id = ?", taskId).First(&taskCheckTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskCheckTask, err
}

// Add persists taskCheckTask to database
func (t *TaskTask) AddTaskCheckTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// Update changes taskCheckTask by id
func (t *TaskTask) UpdateTaskCheckTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})

	if err := db.Debug().Where("id = ?", t.ID).Updates(map[string]interface{}{"task_name": t.TaskName, "risk_level": t.RiskLevel, "status": t.Status, "is_stats": t.IsStats, "task_type": t.TaskType, "execute_type": t.ExecuteType, "task_template_id": t.TaskTemplateID, "target_num": t.TargetNum, "hige_num": t.HigeNum, "middle_num": t.MiddleNum, "low_num": t.LowNum, "safe_num": t.SafeNum, "user_id": t.UserID, "create_time": t.CreateTime, "update_time": t.UpdateTime}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTaskTaskByIds 批量更新任务表   code opt???
func (t *TaskTask) UpdateTaskTaskByIds(ctx context.Context, ids any, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})
	if err := db.Where("id IN ?", ids).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTaskTaskIsStats 更新是否需要重新统计字段
func (t *TaskTask) UpdateTaskTaskIsStats(ctx context.Context, id int, isStats int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})
	if err := db.Where("id = ?", id).Updates(TaskTask{IsStats: isStats}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTaskTaskIsStats 更新是否需要重新统计字段
func (t *TaskTask) UpdateTaskStateIsStatsById(ctx context.Context, id int, status, isStats int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})
	if err := db.Where("id = ?", id).Updates(TaskTask{
		IsStats:    isStats,
		Status:     status,
		UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

// Delete taskCheckTask by id
func (t *TaskTask) DeleteTaskCheckTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})

	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// 批量删除
func (t *TaskTask) DeleteByIds(ctx context.Context, ids []int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})
	if err := db.Where("id in ?", ids).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

// 通过任务状态获取任务运行个数
func (t *TaskTask) GetTaskNumberByTaskStatus(ctx context.Context, taskStatus int) int64 {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})

	var runningNumber int64
	err := db.Where("task_status = ?", taskStatus).Count(&runningNumber).Error
	if err != nil {
		fmt.Println(err)
	}
	return runningNumber
}

// 通过任务状态获取任务运行个数
func (t *TaskTask) GetTaskIDByTaskStatus(ctx context.Context, taskStatus int) []int64 {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})

	var taskIDs []int64
	err := db.Select("id").Where("status = ?", taskStatus).Find(&taskIDs).Error
	if err != nil {
		fmt.Println(err)
	}
	return taskIDs
}

// GetTaskNumByStatus 根据任务状态获取任务数量
func (t *TaskTask) GetTaskNumByStatus(ctx context.Context, taskStatus int) ([]TaskTask, error) {
	var (
		taskCheckTaskList []TaskTask
		db                = mysql.FromContext(ctx).Model(&TaskTask{})
		err               error
	)
	db.Where("status = ?", taskStatus).Find(&taskCheckTaskList)
	return taskCheckTaskList, err
}

// 通过任务状态获取一个任务
func (t *TaskTask) GetOneTaskByStatus(ctx context.Context, executeType int, taskStatus int) TaskTask {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})

	var taskTask TaskTask
	db.Where("execute_type = ?", executeType).Where("status = ?", taskStatus).Order("weight desc,id").Find(&taskTask)
	return taskTask
}

// GetTasksByExecuteTypeStatusLimit 通过执行方式/状态获取多个任务数据
func (t *TaskTask) GetTasksByExecuteTypeStatusLimit(ctx context.Context, executeType, taskStatus, limit int) ([]TaskTask, error) {
	var (
		taskCheckTaskList []TaskTask
		db                = mysql.FromContext(ctx).Model(&TaskTask{})
		err               error
	)
	db.Where("execute_type = ? and status = ?", executeType, taskStatus).Order("weight desc,id").Limit(limit).Find(&taskCheckTaskList)
	return taskCheckTaskList, err
}

// Update changes taskCheckTask by id
func (t *TaskTask) UpdateTaskCheckTaskStatusByID(ctx context.Context, taskID int, status int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})
	if err := db.Where("id = ?", taskID).Updates(map[string]interface{}{"status": status, "update_time": time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

// GetTaskByIsStats 通过是否需要统计获取任务
func (t *TaskTask) GetTaskByIsStats(ctx context.Context, isStat int) []TaskTask {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})
	var taskTaskList []TaskTask
	db.Where("is_stats = ?", isStat).Find(&taskTaskList)
	return taskTaskList
}

// UpdateTaskStatsByTaskId 通过id更新任务是否统计字段
func (t *TaskTask) UpdateTaskStatsByTaskId(ctx context.Context, taskId int, isStat int) {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})
	db.Where("id = ?", taskId).Update("is_stats", isStat)
}

// GetTasksByStatus 通过任务状态获取任务
func (t *TaskTask) GetTasksByStatus(ctx context.Context, taskStatus int) []TaskTask {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})

	var taskTask []TaskTask
	db.Where("status = ?", taskStatus).Find(&taskTask)
	return taskTask
}

// UpdateTargetRiskById 通过id更新目标风险统计
func (t *TaskTask) UpdateTargetRiskById(ctx context.Context, taskId, highNum, middleNum, lowNum, safeNum, riskLevel int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTask{})
	if err := db.Where("id = ?", taskId).Updates(
		map[string]interface{}{
			"hige_num":   highNum,
			"middle_num": middleNum,
			"low_num":    lowNum,
			"safe_num":   safeNum,
			"risk_level": riskLevel,
		}).Error; err != nil {
		return err
	}
	return nil
}

// GetTaskCount 获取任务总数或根据开始时间获取任务总数
func (t *TaskTask) GetTaskCount(ctx context.Context, startTime string, uid int, role int) (int64, int64) {
	var (
		count       int64
		filterCount int64
		db          = mysql.FromContext(ctx).Model(&TaskTask{})
	)

	if role == enums.UserRoleOrdinary {
		db = db.Where("user_id = ?", uid)
	}

	db.Count(&count)

	if startTime != "" {
		db = db.Where("update_time > ?", startTime)
		db.Count(&filterCount)
	}

	return count, filterCount
}

type TaskTrendStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetTaskTrendStat 获取任务趋势统计
func (t *TaskTask) GetTaskTrendStat(ctx context.Context, startTime, dateFormat string, uid int, role int) []TaskTrendStat {
	var (
		TaskTrendStatList []TaskTrendStat
		db                = mysql.FromContext(ctx).Model(&TaskTask{})
	)

	if role == enums.UserRoleOrdinary {
		db = db.Where("user_id = ?", uid)
	}

	db.Select("DATE_FORMAT(update_time, '"+dateFormat+"') as date, COUNT(*) as count").
		Where("update_time > ?", startTime).
		Group("DATE_FORMAT(update_time, '" + dateFormat + "')").
		Find(&TaskTrendStatList)

	return TaskTrendStatList
}

// All 获取所有任务
func (t *TaskTask) All(ctx context.Context, filter string) ([]TaskTask, error) {
	var (
		taskCheckTaskList []TaskTask
		db                = mysql.FromContext(ctx).Model(&TaskTask{})
	)
	if filter != "" {
		db.Where(filter).Find(&taskCheckTaskList)
	} else {
		db.Find(&taskCheckTaskList)
	}
	return taskCheckTaskList, nil
}

// GetLatestTask 获取最新的任务
func (t *TaskTask) GetLatestTask(ctx context.Context) (TaskTask, error) {
	var (
		taskCheckTask TaskTask
		db            = mysql.FromContext(ctx).Model(&TaskTask{})
	)
	db.Order("id desc").Limit(1).Find(&taskCheckTask)
	return taskCheckTask, nil
}

// GetTaskTrendStat 获取任务趋势统计
func (t *TaskTask) GetTaskTrendStatByTaskIds(ctx context.Context, startTime, dateFormat string, taskIdList []int) []TaskTrendStat {
	var (
		TaskTrendStatList []TaskTrendStat
		db                = mysql.FromContext(ctx).Model(&TaskTask{})
	)

	db.Select("DATE_FORMAT(update_time, '"+dateFormat+"') as date, COUNT(*) as count").
		Where("update_time > ? and task_id in ?", startTime, taskIdList).
		Group("DATE_FORMAT(update_time, '" + dateFormat + "')").
		Find(&TaskTrendStatList)

	return TaskTrendStatList
}

// GetAppSecTaskList 应用安全任务列表（按 task_type 筛选）
func (t *TaskTask) GetAppSecTaskList(ctx context.Context, taskType, userID, page, limit int, search string) ([]TaskTask, int64, error) {
	var (
		list  []TaskTask
		count int64
		db    = mysql.FromContext(ctx).Model(&TaskTask{}).Where("task_type = ?", taskType)
	)
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if search != "" {
		db = db.Where("task_name LIKE ?", "%"+search+"%")
	}
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	err := db.Order("create_time DESC").Limit(limit).Offset((page - 1) * limit).Find(&list).Error
	return list, count, err
}
