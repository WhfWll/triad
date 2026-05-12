package mysqls

import (
	"context"
	"fmt"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type TaskTarget struct {
	ID               int       `gorm:"column:id;primary_key" json:"id"`                   // 主键
	TaskID           int       `gorm:"column:task_id" json:"taskID"`                      // 所属任务id
	TargetURL        string    `gorm:"column:target_url" json:"targetURL"`                // 测试目标地址
	Status           int       `gorm:"column:status" json:"status"`                       // 目标状态
	Weight           int       `gorm:"column:weight" json:"weight"`                       //优先权重
	RiskLevel        int       `gorm:"column:risk_level" json:"riskLevel"`                // 任务风险等级
	OpSys            string    `gorm:"column:op_sys" json:"opSys"`                        // 操作系统
	IsAlive          int       `gorm:"column:is_alive" json:"isAlive"`                    // 是否存活
	TargetType       int       `gorm:"column:target_type" json:"targetType"`              // 目标类型
	TaskTemplateID   int       `gorm:"column:task_template_id" json:"taskTemplateID"`     // 所选择的任务场景id
	TaskTemplateJSON string    `gorm:"column:task_template_json" json:"taskTemplateJSON"` // 任务场景参数
	IsRemoteSession  int       `gorm:"column:is_remote_session" json:"isRemoteSession"`   // 是否有远程会话:1否，2是
	UserID           int       `gorm:"column:user_id" json:"userID"`                      // 提交者id
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`              // 发生时间
	UpdateTime       time.Time `gorm:"column:update_time" json:"updateTime"`              // 发生时间
	EndTime          time.Time `gorm:"column:end_time" json:"endTime"`                    // 发生时间
	UseScore         int       `gorm:"column:use_score" json:"useScore"`                  // 利用评分
	IsScore          int       `gorm:"column:is_score" json:"IsScore"`                    //利用评分状态
	ExtendField      string    `gorm:"extend_field" json:"extendField"`                   // 其他字段
}

// TableName sets insert table name for this struct type
func (t *TaskTarget) TableName() string {
	return "task_target"
}

// 参数ids必须是一个[]int或[]string数组
func (t *TaskTarget) GetTaskTargetListByIds(ctx context.Context, ids interface{}, status int) []TaskTarget {
	var (
		taskTargetList []TaskTarget
		db             = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	db.Where("id IN ?", ids).Find(&taskTargetList)
	return taskTargetList
}

// GetTaskTarget 根据id查询目标信息
func (t *TaskTarget) GetTaskTarget(ctx context.Context, id int) (TaskTarget, error) {
	var (
		taskTarget TaskTarget
		err        error
		db         = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	curErr := db.Where("id = ?", id).First(&taskTarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskTarget, err
}

// GetByIds 根据ids查询目标信息
func (t *TaskTarget) GetByIds(ctx context.Context, ids []int, status int) []TaskTarget {
	var (
		taskTarget []TaskTarget
		db         = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	db.Where("id in ?", ids).Find(&taskTarget)

	return taskTarget
}

// Add persists taskTarget to database
func (t *TaskTarget) AddTaskTarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// Update changes taskTarget by id
func (t *TaskTarget) UpdateTaskTarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})

	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTargetById 根据id更新目标数据
func (t *TaskTarget) UpdateTargetById(ctx context.Context, id int, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	if err := db.Where("id = ?", id).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTaskTarget 根据id删除目标数据
// ids必须为[]int/[]string数组
func (t *TaskTarget) DeleteTaskTarget(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	if err := db.Where("id IN ?", ids).Delete(&TaskTarget{}).Error; err != nil {
		return err
	}
	return nil
}

// 通过状态获取目标数量
func (t *TaskTarget) GetTargetNumberByStatus(ctx context.Context, taskStatus int) int64 {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	var runningNumber int64
	err := db.Where("status = ?", taskStatus).Count(&runningNumber).Error
	if err != nil {
		fmt.Println(err)
	}
	return runningNumber
}

// GetTargetNumByStatus 通过状态获取目标数量
func (t *TaskTarget) GetTargetNumByStatus(ctx context.Context, taskStatus int) (int64, error) {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	var total int64
	err := db.Where("status = ?", taskStatus).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// 通过状态获取目标数量
func (t *TaskTarget) GetTargetNumberByStatusAndTaskID(ctx context.Context, statusList []int, taskID int) int64 {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})

	var runningNumber int64
	err := db.Where("status in ? and task_id = ?", statusList, taskID).Count(&runningNumber).Error
	if err != nil {
		fmt.Println(err)
	}
	return runningNumber
}

// GetTargetNumberTaskID 通过状态获取目标数量
func (t *TaskTarget) GetTargetNumberTaskID(ctx context.Context, taskID string) int64 {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	var allNumber int64
	err := db.Where("task_id = ?", taskID).Count(&allNumber).Error
	if err != nil {
		fmt.Println(err)
	}
	return allNumber
}

// 通过状态获取一个测试目标
func (t *TaskTarget) GetOneTargetByStatus(ctx context.Context, status int, taskID int64) TaskTarget {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})

	var taskCheckTarget TaskTarget
	db.Where("status = ? and task_id = ?", status, taskID).Limit(1).Find(&taskCheckTarget)

	return taskCheckTarget
}

// 通过状态获取目标数量
func (t *TaskTarget) GetOneWaitTarget(ctx context.Context, status int, taskID []int64) TaskTarget {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})

	var taskCheckTarget TaskTarget
	db.Where("status = ? and task_id in ? ", status, taskID).Order("weight desc,id").Limit(1).Find(&taskCheckTarget)

	return taskCheckTarget
}

// GetTargetsByTaskIdsStatus 根据任务id和状态 获取多条目标数据
func (t *TaskTarget) GetTargetsByTaskIdsStatus(ctx context.Context, taskID any, status int, limit int) []TaskTarget {
	var (
		db              = mysql.FromContext(ctx).Model(&TaskTarget{})
		taskCheckTarget []TaskTarget
	)
	//db.Where("task_id in ? and status = ? and is_score in ?", taskID, status, isScore).Order("weight desc,is_score desc,id").Limit(limit).Find(&taskCheckTarget)
	db.Where("task_id in ? and status = ? ", taskID, status).Order("weight desc,is_score desc,id").Limit(limit).Find(&taskCheckTarget)
	return taskCheckTarget
}

// GetTaskTargetByUrl 根据任务ID和URL获取目标
func (t *TaskTarget) GetTaskTargetByUrl(ctx context.Context, taskID int, url string) (TaskTarget, error) {
	var (
		db         = mysql.FromContext(ctx).Model(&TaskTarget{})
		taskTarget TaskTarget
	)
	err := db.Where("task_id = ? AND target_url = ?", taskID, url).First(&taskTarget).Error
	return taskTarget, err
}

// 批量添加
func (t *TaskTarget) Adds(ctx context.Context, data *[]TaskTarget) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})

	if err := db.Create(data).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTaskCheckTargetStatus 更新检测目标状态
func (t *TaskTarget) UpdateTaskCheckTargetStatus(ctx context.Context, targetID int, status int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	if err := db.Where("id = ?", targetID).Updates(TaskTarget{Status: status, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTargetStatusById 根据id更新检测目标状态
func (t *TaskTarget) UpdateTargetStatusById(ctx context.Context, id int, status int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	if err := db.Where("id = ?", id).Updates(TaskTarget{Status: status, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTargetStatusByIds 根据ids更新检测目标状态
func (t *TaskTarget) UpdateTargetStatusByIds(ctx context.Context, ids any, status int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	if err := db.Where("id IN ?", ids).Updates(TaskTarget{Status: status, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTargetStatsById 根据id更新统计信息
func (t *TaskTarget) UpdateTargetStatsById(ctx context.Context, id, vulNum, risklevel int, vulNumArray [6]int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	var tmp = map[string]interface{}{
		// "vul_num":    vulNum,
		"risk_level": risklevel,
		//"fatal_num":   vulNumArray[1],
		// "high_num":    vulNumArray[2],
		// "middle_num":  vulNumArray[3],
		// "low_num":     vulNumArray[4],
		// "safe_num":    vulNumArray[5],
		"update_time": time.Now(),
	}
	if err := db.Where("id = ?", id).Updates(tmp).Error; err != nil {
		return err
	}
	return nil
}

// GetTargetListByTaskId 根据任务id查询列表
func (t *TaskTarget) GetTargetListByTaskId(ctx context.Context, taskId, riskLevel int, search string, isalive int, page, limit int) ([]TaskTarget, int64, error) {
	var (
		taskTargetList []TaskTarget
		count          int64
		db             = mysql.FromContext(ctx).Model(&TaskTarget{})
		query          string
		args           []interface{}
	)
	//query += "task_id = ? and is_alive = ?"
	//args = append(args, taskId, isalive)
	query += "task_id = ?"
	args = append(args, taskId)
	if riskLevel != 0 {
		query += " and risk_level = ?"
		args = append(args, riskLevel)
	}
	if len(search) > 0 {
		query += " and (target_url LIKE ? or op_sys LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Order("is_alive desc,use_score desc, id desc").Find(&taskTargetList)
	return taskTargetList, count, nil
}

// GetTargetsByTaskId 获取taskId下的所有target
func (t *TaskTarget) GetTargetsByTaskId(ctx context.Context, taskId int) []TaskTarget {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	var targetList []TaskTarget
	db.Where("task_id = ?", taskId).Find(&targetList)
	return targetList
}

// GetTargetsByTaskId 获取taskId下的所有target
func (t *TaskTarget) GetTargetNumberByAliveAndTaskId(ctx context.Context, taskId int) int64 {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	var number int64
	db.Where("is_alive = 1 and task_id = ?", taskId).Count(&number)
	return number
}

// 批量删除 通过task_ids
func (t *TaskTarget) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	if err := db.Where("task_id in ?", taskIds).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

// 通过task_id获取
func (t *TaskTarget) GetByTaskId(ctx context.Context, taskId int) ([]TaskTarget, int64) {
	var (
		taskTargetList []TaskTarget
		count          int64
		db             = mysql.FromContext(ctx).Model(&TaskTarget{})
	)

	db.Count(&count)
	_ = db.Where("task_id = ?", taskId).Order("risk_level").Find(&taskTargetList)
	return taskTargetList, count
}

// UpdateTargetAliveById 根据id更新检测目标存活性
func (t *TaskTarget) UpdateTargetAliveById(ctx context.Context, id int, isAlive int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	if err := db.Where("id = ?", id).Updates(TaskTarget{IsAlive: isAlive, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTargetOpSys 更新操作系统信息
func (t *TaskTarget) UpdateTargetOpSys(ctx context.Context, id int, opSys string) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	if err := db.Where("id = ?", id).Updates(TaskTarget{OpSys: opSys, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

// HandleTimeoutTargetList 通过状态获取目标数量
func (t *TaskTarget) HandleTimeoutTargetList(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	overTime := time.Now().Add(-getTargetTimeout()).Format(utils.DateTime)
	err := db.Where("status = ? and create_time < ?", enums.TargetStatusRunning, overTime).Updates(TaskTarget{Status: enums.TargetStatusFinish, UpdateTime: time.Now()}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetTargetByStatus 通过状态获取测试目标
func (t *TaskTarget) GetTargetByStatus(ctx context.Context, status int) []TaskTarget {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})

	var taskCheckTarget []TaskTarget
	db.Where("status = ? ", status).Find(&taskCheckTarget)

	return taskCheckTarget
}

// GetTargetsByTaskIds 根据任务id获取多条目标数据
func (t *TaskTarget) GetTargetsByTaskIds(ctx context.Context, taskID []int) []TaskTarget {
	var (
		db              = mysql.FromContext(ctx).Model(&TaskTarget{})
		taskCheckTarget []TaskTarget
	)
	db.Where("task_id in ? ", taskID).Find(&taskCheckTarget)
	return taskCheckTarget
}

// GetTargetsByTargetURL 根据目标获取目标数据
func (t *TaskTarget) GetTargetsByTargetURL(ctx context.Context, taskId int, targetUrl string) (TaskTarget, error) {
	var (
		taskTarget TaskTarget
		err        error
		db         = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	curErr := db.Where("task_id = ? and target_url = ?", taskId, targetUrl).First(&taskTarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskTarget, err
}

// GetTargetsByTargetURLLike 根据目标模糊搜索获取目标数据
func (t *TaskTarget) GetTargetsByTargetURLLike(ctx context.Context, taskId int, targetUrl string) (TaskTarget, error) {
	var (
		taskTarget TaskTarget
		err        error
		db         = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	curErr := db.Where("task_id = ? and target_url like ?", taskId, "%"+targetUrl+"%").First(&taskTarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskTarget, err
}

// GetTargetCount 获取目标总数或根据开始时间获取目标总数
func (t *TaskTarget) GetTargetCount(ctx context.Context, startTime string) (int64, int64) {
	var (
		count       int64
		filterCount int64
		db          = mysql.FromContext(ctx).Model(&TaskTarget{})
	)

	db.Count(&count)

	if startTime != "" {
		db.Where("update_time > ?", startTime)
		db.Count(&filterCount)
	}

	return count, filterCount
}

type TargetRiskStat struct {
	RiskLevel int `json:"risk_level"`
	Count     int `json:"count"`
}

// GetTargetRiskStat 目标风险统计
func (t *TaskTarget) GetTargetRiskStat(ctx context.Context, uid int, role int) []TargetRiskStat {
	var (
		targetRiskStatList []TargetRiskStat
		db                 = mysql.FromContext(ctx).Model(&TaskTarget{})
	)

	if role == enums.UserRoleOrdinary {
		db = db.Joins("JOIN task_task ON task_target.task_id = task_task.id").Where("task_task.user_id = ?", uid)
	}

	db.Select("task_target.risk_level, Count(*) as count").Group("task_target.risk_level").Find(&targetRiskStatList)

	return targetRiskStatList
}

// All 获取所有目标
func (t *TaskTarget) All(ctx context.Context, filter string) ([]TaskTarget, error) {
	var (
		taskTargetList []TaskTarget
		db             = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	if filter != "" {
		db.Where(filter).Find(&taskTargetList)
	} else {
		db.Find(&taskTargetList)
	}
	return taskTargetList, nil
}

// GetLatestTarget 获取最新的目标
func (t *TaskTarget) GetLatestTarget(ctx context.Context) (TaskTarget, error) {
	var (
		taskTarget TaskTarget
		db         = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	db.Order("id desc").Limit(1).Find(&taskTarget)
	return taskTarget, nil
}

// GroupByTargetUrl 获取所有目标，并按照目标分组
type allTargetUrl struct {
	ID        int    `gorm:"column:id;primary_key" json:"id"`    // 主键
	TargetURL string `gorm:"column:target_url" json:"targetURL"` // 测试目标地址
}

func (t *TaskTarget) AllTargetUrl(ctx context.Context, status int) (taskTarget []TaskTarget) {
	var (
		db = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	db.Select("id,target_url").Where("status = ?", status).Order("create_time DESC").Find(&taskTarget)
	return
}

// GetTargetsByTaskIdAndStatus 根据任务ID和状态获取目标列表
func (t *TaskTarget) GetTargetsByTaskIdAndStatus(ctx context.Context, taskId, status int) []TaskTarget {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	var targets []TaskTarget
	db.Where("task_id = ? AND status = ?", taskId, status).Find(&targets)
	return targets
}

// GetTargetsByTaskId 根据任务ID和状态获取目标列表
func (t *TaskTarget) GetTargetsByTaskIdAndStatusList(ctx context.Context, taskId int, status []int) []TaskTarget {
	var db = mysql.FromContext(ctx).Model(&TaskTarget{})
	var targets []TaskTarget
	db.Where("task_id = ? AND status in ?", taskId, status).Find(&targets)
	return targets
}

// GetTargetListByTargetUrl 根据target_url查询列表
func (t *TaskTarget) GetTargetListByTargetUrl(ctx context.Context, targetURL string) (TaskTarget, error) {
	var (
		db             = mysql.FromContext(ctx).Model(&TaskTarget{})
		taskTargetList TaskTarget
		query          string
		args           []interface{}
	)
	query += "status = " + strconv.Itoa(enums.TargetStatusFinish) + " and target_url LIKE ? "
	args = append(args, "%"+targetURL+"%")
	dbs := db.Where(query, args...)
	dbs.Order("id desc").First(&taskTargetList)
	return taskTargetList, nil
}

// BatchGetTargetListByTargetUrl 根据 target_url 批量查询最新完成状态的目标列表
// NOTE:
//  1. 由于 SQL 中 IN 条件中元素过多会导致性能下降甚至报错（数据库对 SQL 语句长度有限制）
//  2. 先查询所有符合条件的记录，再在 Go 代码中按 target_url 分组，取每组最大 id 记录作为最新数据 join的方式验证后效果很不好。
func (t *TaskTarget) BatchGetTargetListByTargetUrl(ctx context.Context, targetURLs []string) ([]TaskTarget, error) {
	const batchSize = 200
	db := mysql.FromContext(ctx)
	var allTargets []TaskTarget
	for i := 0; i < len(targetURLs); i += batchSize {
		end := i + batchSize
		if end > len(targetURLs) {
			end = len(targetURLs)
		}
		batch := targetURLs[i:end]

		var batchTargets []TaskTarget
		err := db.Model(&TaskTarget{}).
			Where("status = ? AND target_url IN ?", enums.TargetStatusFinish, batch).
			Find(&batchTargets).Error
		if err != nil {
			return nil, err
		}
		allTargets = append(allTargets, batchTargets...)
	}

	// 用 map 记录每个 target_url 最大 id 的记录
	latestMap := make(map[string]TaskTarget)
	for _, t := range allTargets {
		existing, ok := latestMap[t.TargetURL]
		if !ok || t.ID > existing.ID {
			latestMap[t.TargetURL] = t
		}
	}

	// 转成 slice 返回
	result := make([]TaskTarget, 0, len(latestMap))
	for _, v := range latestMap {
		result = append(result, v)
	}
	return result, nil
}

// GetFinishedTargetURLToIDMap 以 target_url 为 key，ID 为最大已完成的 ID
func (t *TaskTarget) GetFinishedTargetURLToIDMap(ctx context.Context) (map[string]int, error) {
	var results []struct {
		TargetURL string
		ID        int
	}
	db := mysql.FromContext(ctx).Model(&TaskTarget{}).
		Select("target_url, MAX(id) as id").
		Where("status = ?", enums.TargetStatusFinish).
		Group("target_url")
	if err := db.Scan(&results).Error; err != nil {
		return nil, err
	}
	m := make(map[string]int, len(results))
	for _, r := range results {
		m[r.TargetURL] = r.ID
	}
	return m, nil
}

// GetFinishedTargetURLToIDMapByTime 时间范围内 以 target_url 为 key，ID 为最大已完成的 ID
func (t *TaskTarget) GetFinishedTargetURLToIDMapByTime(ctx context.Context, startTime, endTime string) (map[string]int, error) {
	var results []struct {
		TargetURL string
		ID        int
	}
	db := mysql.FromContext(ctx).Model(&TaskTarget{}).
		Select("target_url, MAX(id) as id").
		Where("status = ?", enums.TargetStatusFinish).
		Group("target_url")
	if startTime != "" {
		db.Where("update_time > ?", startTime)
	}
	if endTime != "" {
		db.Where("update_time <= ?", endTime)
	}
	if err := db.Scan(&results).Error; err != nil {
		return nil, err
	}
	m := make(map[string]int, len(results))
	for _, r := range results {
		m[r.TargetURL] = r.ID
	}
	return m, nil
}

// GetTargetsByTaskID 根据任务ID 获取目标数据
func (t *TaskTarget) GetTargetsByTaskID(ctx context.Context, taskId int) ([]TaskTarget, error) {
	var (
		taskTarget []TaskTarget
		err        error
		db         = mysql.FromContext(ctx).Model(&TaskTarget{})
	)
	curErr := db.Select("id,task_id,target_url,op_sys,target_type,is_alive,status,user_id,risk_level").Where("task_id = ? ", taskId).Find(&taskTarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskTarget, err
}
