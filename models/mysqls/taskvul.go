package mysqls

import (
	"context"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type TaskVul struct {
	ID             int       `gorm:"column:id;primary_key" json:"id"`               // 主键
	DataType       int       `gorm:"column:data_type" json:"dataType"`              // 数据类型
	TaskID         int       `gorm:"column:task_id" json:"taskID"`                  // 所属任务id
	TargetID       int       `gorm:"column:target_id" json:"targetID"`              // 所属目标id
	TargetUrl      string    `gorm:"column:target_url" json:"targetUrl"`            // 测试目标地址
	Pocname        string    `gorm:"column:pocname" json:"pocname"`                 // 漏洞标识
	Name           string    `gorm:"column:name" json:"name"`                       // 漏洞名称
	Class          int       `gorm:"column:class" json:"class"`                     // 漏洞分类
	Type           int       `gorm:"column:type" json:"type"`                       // 漏洞类型
	Risk           int       `gorm:"column:risk" json:"risk"`                       // 风险等级
	Location       string    `gorm:"column:location" json:"location"`               // 漏洞位置或地址
	Status         int       `gorm:"column:status" json:"status"`                   // 漏洞状态
	TestStatus     int       `gorm:"column:test_status" json:"testStatus"`          // 测试状态
	ExploitImpact  string    `gorm:"column:exploit_impact" json:"exploitImpact"`    // 利用影响
	VulID          string    `gorm:"column:vul_id" json:"vulID"`                    // 漏洞id
	Description    string    `gorm:"column:description" json:"description"`         // 漏洞描述
	FixSuggest     string    `gorm:"column:fix_suggest" json:"fixSuggest"`          // 修复建议
	PublishedTime  string    `gorm:"column:published_time" json:"publishedTime"`    // 披露时间
	AffectRange    string    `gorm:"column:affect_range" json:"affectRange"`        // 影响范围
	TargetResultID int       `gorm:"column:target_result_id" json:"targetResultID"` // 检测结果id
	VulNumber      string    `gorm:"column:vul_number" json:"vulNumber"`            // 漏洞编号
	VulAddress     string    `gorm:"column:vul_address" json:"vulAddress"`          // 漏洞地址
	RefUrl         string    `gorm:"column:ref_url" json:"refUrl"`                  // 参考连接
	Cvss           string    `gorm:"column:cvss" json:"cvss"`                       // cvss评分
	VulResult      string    `gorm:"column:vul_result" json:"vulResult"`            // 漏洞结果
	VulParam       string    `gorm:"column:vul_param" json:"vulParam"`              // 漏洞请求参数
	VerMsg         string    `gorm:"column:ver_msg" json:"verMsg"`                  // 验证报文
	DecisionVulId  string    `gorm:"column:decision_vul_id" json:"decisionVulId"`   // 决策引擎唯一漏洞ID
	Snapshot       string    `gorm:"column:snapshot" json:"snapshot"`               // 截图
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`          // 创建时间
	UpdateTime     time.Time `gorm:"column:update_time" json:"updateTime"`          // 修改时间
	PatchUrl       string    `gorm:"column:patch_url" json:"patch_url"`             // 补丁地址
}

// TableName sets insert table name for this struct type
func (t *TaskVul) TableName() string {
	return "task_vul"
}

// GetTaskVulList 任务漏洞信息列表及筛选
func (t *TaskVul) GetTaskVulList(ctx context.Context, taskId, targetId, vultype, risk int, search string, dataType, page, limit, status int) ([]TaskVul, int64, error) {
	var (
		taskVulList []TaskVul
		count       int64
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
		query       string
		args        []interface{}
	)
	if targetId > 0 {
		query += "target_id = " + strconv.Itoa(targetId) + " and "
	}
	query += "task_id = ? and data_type = ?"
	args = append(args, taskId, dataType)
	if vultype > 0 {
		query += " and type = ? "
		args = append(args, vultype)
	}
	if risk > 0 {
		query += " and risk = ? "
		args = append(args, risk)
	}
	if len(search) > 0 {
		query += " and (target_url LIKE ? or name LIKE ? or location LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status > 0 {
		query += " and status = ? "
		args = append(args, status)
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Order("risk").Offset(limit * (page - 1)).Find(&taskVulList)

	return taskVulList, count, nil
}

// GetTaskVulDetail 通过 target_id, pocname, location 获取详情
func (t *TaskVul) GetTaskVulDetail(ctx context.Context, targetId int, pocName, location string) (TaskVul, error) {
	var (
		taskVul TaskVul
		err     error
		db      = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	curErr := db.Where("target_id = ? AND pocname = ? AND location = ?", targetId, pocName, location).First(&taskVul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskVul, err
}

// UpdateTaskVul 更新漏洞信息
func (t *TaskVul) UpdateTaskVul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}
	return nil
}

// SecondScanGetTaskVulList 二次扫描漏洞列表
func (t *TaskVul) SecondScanGetTaskVulList(ctx context.Context, taskId []int, targetId, vultype, risk int, search string, dataType, page, limit, status int) ([]TaskVul, int64, error) {
	var (
		taskVulList []TaskVul
		count       int64
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
		query       string
		args        []interface{}
	)
	if targetId > 0 {
		query += "target_id = " + strconv.Itoa(targetId) + " and "
	}
	query += "task_id in ? and data_type = ?"
	args = append(args, taskId, dataType)
	if vultype > 0 {
		query += " and type = ? "
		args = append(args, vultype)
	}
	if risk > 0 {
		query += " and risk = ? "
		args = append(args, risk)
	}
	if len(search) > 0 {
		query += " and (target_url LIKE ? or name LIKE ? or location LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status > 0 {
		query += " and status = ? "
		args = append(args, status)
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Order("risk").Offset(limit * (page - 1)).Find(&taskVulList)

	return taskVulList, count, nil
}

// GetTaskVulListByTargetIds 根据目标id获取任务漏洞信息
// targetIds 必须是一个[]int/[]string数组
func (t *TaskVul) GetTaskVulListByTargetIds(ctx context.Context, targetIds any, dataType int) ([]TaskVul, error) {
	var (
		taskVulList []TaskVul
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	// 按风险等级从高到低排序：致命(1)、高危(2)、中危(3)、低危(4)
	db.Select("id,target_id,type,status,risk").Where("target_id IN ? and data_type = ?", targetIds, dataType).Order("risk ASC").Find(&taskVulList)
	return taskVulList, nil
}

// GetTaskVulListByIds 根据ids获取任务漏洞信息
// ids 必须是一个[]int/[]string数组
func (t *TaskVul) GetTaskVulListByIds(ctx context.Context, ids any) ([]TaskVul, error) {
	var (
		taskVulList []TaskVul
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	db.Where("id IN ?", ids).Find(&taskVulList)
	return taskVulList, nil
}

// GetTaskVul 根据id查询一条数据
func (t *TaskVul) GetTaskVul(ctx context.Context, id int) (TaskVul, error) {
	var (
		taskVul TaskVul
		err     error
		db      = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	curErr := db.Where("id = ?", id).First(&taskVul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskVul, err
}

// GetTaskVul 根据id查询一条数据
func (t *TaskVul) GetTaskVulByPocNameAndTargetId(ctx context.Context, targetId int, pocName string) (TaskVul, error) {
	var (
		taskVul TaskVul
		err     error
		db      = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	curErr := db.Where("target_id = ? and pocname = ?", targetId, pocName).First(&taskVul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskVul, err
}

// Add persists taskVul to database
func (t *TaskVul) AddTaskVul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// UpdateVerMsgById 根据id修改请求报文和响应报文
func (t *TaskVul) UpdateVerMsgById(ctx context.Context, id int, verMsg string) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})
	var tmpData = map[string]interface{}{
		"ver_msg":     verMsg,
		"update_time": time.Now(),
	}
	if err := db.Where("id = ?", id).Updates(tmpData).Error; err != nil {
		return err
	}
	return nil
}

// UpdateVerMsgById 根据id修改请求报文和响应报文
func (t *TaskVul) UpdateStatusById(ctx context.Context, id int, status int) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})
	var tmpData = map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}
	if err := db.Where("id = ?", id).Updates(tmpData).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTestStatusByIds 根据id修改测试状态
func (t *TaskVul) UpdateTestStatusByIds(ctx context.Context, ids any, testStatus int) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})
	var tmpData = map[string]interface{}{
		"test_status": testStatus,
		"update_time": time.Now(),
	}
	if err := db.Where("id IN ?", ids).Updates(tmpData).Error; err != nil {
		return err
	}
	return nil
}

// UpdateVerMsgById 根据id修改请求报文和响应报文
func (t *TaskVul) UpdateByIdAndDataType(ctx context.Context, id int, dataType int, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})
	if err := db.Where("id = ? and data_type = ?", id, dataType).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTaskVul 删除漏洞测试
// ids必须为[]int/[]string数组
func (t *TaskVul) DeleteTaskVul(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})
	if err := db.Where("id IN ?", ids).Delete(&TaskVul{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTaskVulByTargetIds 根据目标ids删除漏洞测试
// ids必须为[]int/[]string数组
func (t *TaskVul) DeleteTaskVulByTargetIds(ctx context.Context, targetIds any) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})
	if err := db.Where("target_id IN ?", targetIds).Delete(&TaskVul{}).Error; err != nil {
		return err
	}
	return nil
}

// 批量删除 通过task_ids
func (t *TaskVul) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&TaskVul{})
	if err := db.Where("task_id in ?", taskIds).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

// GetsByTaskId 依据任务ID获取所有
func (t *TaskVul) GetsByTaskId(ctx context.Context, taskId int, dataType int) []TaskVul {
	var (
		taskVulList []TaskVul
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	// 按风险等级从高到低排序：致命(1)、高危(2)、中危(3)、低危(4)
	db.Where("task_id = ?  and data_type = ?", taskId, dataType).Order("risk ASC").Find(&taskVulList)
	return taskVulList
}

// GetsByTaskId 依据任务ID获取所有
func (t *TaskVul) GetByTargetResultId(ctx context.Context, targetResultId int) (TaskVul, error) {
	var (
		taskVulList TaskVul
		err         error
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	if curErr := db.Where("target_result_id = ?", targetResultId).First(&taskVulList).Error; curErr != nil {
		err = curErr
	}
	return taskVulList, err
}

// GetByTargetResultIds 依据结果IDs获取所有
func (t *TaskVul) GetByTargetResultIds(ctx context.Context, targetResultId []int, targetId int) []TaskVul {
	var (
		taskVulList []TaskVul
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	db.Where("target_result_id in ? and target_id = ?", targetResultId, targetId).Find(&taskVulList)
	return taskVulList
}

// GetsByTargetId 依据目标ID获取所有
func (t *TaskVul) GetsByTargetId(ctx context.Context, targetId int, dataType int) []TaskVul {
	var (
		taskVulList []TaskVul
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	// 按风险等级从高到低排序：致命(1)、高危(2)、中危(3)、低危(4)
	db.Where("target_id = ? and data_type = ?", targetId, dataType).Order("risk ASC").Find(&taskVulList)
	return taskVulList
}

// GetsByTaskId 依据任务ID获取所有
func (t *TaskVul) GetById(ctx context.Context, id int) (TaskVul, error) {
	var (
		taskVulList TaskVul
		err         error
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	if curErr := db.Where("id = ?", id).First(&taskVulList).Error; curErr != nil {
		err = curErr
	}
	return taskVulList, err
}

// GetTaskVulCount 获取任务漏洞总数或根据开始时间获取任务漏洞总数
func (t *TaskVul) GetTaskVulCount(ctx context.Context, startTime string) (int64, int64) {
	var (
		count       int64
		filterCount int64
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)

	db.Count(&count)

	if startTime != "" {
		db.Where("update_time > ?", startTime)
		db.Count(&filterCount)
	}

	return count, filterCount
}

type TaskVulRiskStat struct {
	Risk  int `json:"risk"`
	Count int `json:"count"`
}

// GetTaskVulRiskStat 任务漏洞风险统计
func (t *TaskVul) GetTaskVulRiskStat(ctx context.Context, uid int, role int) []TaskVulRiskStat {
	var (
		taskVulRiskStatList []TaskVulRiskStat
		db                  = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	if role == enums.UserRoleOrdinary {
		db = db.Joins("JOIN task_task ON task_vul.task_id = task_task.id").Where("task_task.user_id = ?", uid)
	}
	db.Select("task_vul.risk, Count(*) as count").Group("task_vul.risk").Find(&taskVulRiskStatList)

	return taskVulRiskStatList
}

type TaskVulTypeStat struct {
	Type  int `json:"type"`
	Count int `json:"count"`
}

// GetTaskVulTypeStat 获取漏洞类型统计
func (t *TaskVul) GetTaskVulTypeStat(ctx context.Context, startTime string, uid int, role int) []TaskVulTypeStat {
	var (
		taskVulTypeStatList []TaskVulTypeStat
		db                  = mysql.FromContext(ctx).Model(&TaskVul{})
	)

	if role == enums.UserRoleOrdinary {
		db = db.Joins("JOIN task_task ON task_vul.task_id = task_task.id").Where("task_task.user_id = ?", uid)
	}

	if startTime != "" {
		db = db.Where("task_vul.create_time > ?", startTime)
	}

	db.Select("task_vul.type, Count(*) as count").Group("task_vul.type").Find(&taskVulTypeStatList)

	return taskVulTypeStatList
}

type TaskVulTrendStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetTaskVulFindTrendStat 获取任务漏洞发现趋势统计 startTime控制查询范围 dateFormat控制日期的格式化输出
func (t *TaskVul) GetTaskVulFindTrendStat(ctx context.Context, startTime, dateFormat string, uid int, role int) []TaskVulTrendStat {
	var (
		taskVulTrendStatList []TaskVulTrendStat
		db                   = mysql.FromContext(ctx).Model(&TaskVul{})
	)

	if role == enums.UserRoleOrdinary {
		db = db.Joins("JOIN task_task ON task_vul.task_id = task_task.id").Where("task_task.user_id = ?", uid)
	}

	if startTime != "" {
		db = db.Where("task_vul.create_time > ?", startTime)
	}

	db.Select("DATE_FORMAT(task_vul.create_time, '" + dateFormat + "') as date, COUNT(*) as count").
		Group("DATE_FORMAT(task_vul.create_time, '" + dateFormat + "')").
		Find(&taskVulTrendStatList)

	return taskVulTrendStatList
}

// All 获取所有目标
func (t *TaskVul) All(ctx context.Context, filter string) ([]TaskVul, error) {
	var (
		taskVulList []TaskVul
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	if filter != "" {
		db.Where(filter).Find(&taskVulList)
	} else {
		db.Find(&taskVulList)
	}
	return taskVulList, nil
}

func (t *TaskVul) AllByPage(ctx context.Context, page, size int) ([]TaskVul, int64) {
	var (
		count       int64
		taskVulList []TaskVul
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	db.Count(&count)
	db = db.Order("create_time desc").Offset((page - 1) * size).Limit(size).Find(&taskVulList)
	return taskVulList, count
}

// GetTaskVulByTaskIds 通过task_id获取任务漏洞
func (t *TaskVul) GetTaskVulByTaskIds(ctx context.Context, taskIds []int) []TaskVul {
	var (
		taskVulList []TaskVul
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	// 按风险等级从高到低排序：致命(1)、高危(2)、中危(3)、低危(4)
	db.Where("task_id in ?", taskIds).Order("risk ASC").Find(&taskVulList)
	return taskVulList
}

// VulStatsResult 漏洞统计结果
type VulStatsResult struct {
	Status int `json:"status"`
	Risk   int `json:"risk"`
	Count  int `json:"count"`
}

// GetVulStatsByTaskId 获取漏洞统计（聚合查询，避免加载大文本字段）
func (t *TaskVul) GetVulStatsByTaskId(ctx context.Context, taskId int, dataType int) []VulStatsResult {
	var (
		statsResult []VulStatsResult
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	// 使用聚合查询，只统计数量，不加载description、fix_suggest等大文本字段
	db.Select("status, risk, count(*) as count").
		Where("task_id = ? and data_type = ?", taskId, dataType).
		Group("status, risk").
		Find(&statsResult)
	return statsResult
}

// GetAllTaskVulList 全部任务漏洞信息列表
func (t *TaskVul) GetAllTaskVulList(ctx context.Context) ([]TaskVul, int64, error) {
	var (
		taskVulList []TaskVul
		count       int64
		db          = mysql.FromContext(ctx).Model(&TaskVul{})
	)
	db.Count(&count)
	db.Where("risk != 0").Order("risk asc").Find(&taskVulList)
	return taskVulList, count, nil
}

// GetTaskVulListByIP 通过ip信息任务漏洞信息列表及筛选
func (t *TaskVul) GetTaskVulListByIP(ctx context.Context, ip string) ([]TaskVul, int64, error) {
	var (
		taskVulList []TaskVul
		count       int64
		db          = mysql.FromContext(ctx).Model(&TaskVul{}).Select("id,target_url,name,type,risk,location,status,verify_type,test_status,create_time")
	)
	if ip != "" {
		db = db.Where("target_url LIKE ?", "%"+ip+"%")
	}
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("risk ASC").Find(&taskVulList).Error; err != nil {
		return nil, 0, err
	}
	return taskVulList, count, nil
}

type TargetRiskResultCount struct {
	TargetId int `json:"target_id"`
	Count    int `json:"count"`
}

func (t *TaskVul) GetTargetIdsRiskResultCount(ctx context.Context, taskId int, targetIds []int, risk int) ([]TargetRiskResultCount, error) {
	var (
		db     = mysql.FromContext(ctx).Model(&TaskVul{}).Debug()
		result []TargetRiskResultCount
		err    error
	)
	if len(targetIds) == 0 {
		return result, nil
	}
	db = db.Select("target_id, Count(*) as count")
	if risk > 0 {
		db = db.Where("risk <= ?", risk).Where("risk > 0")
	}
	err = db.Where("task_id = ?", taskId).Where("target_id in ?", targetIds).Group("target_id").Find(&result).Error
	return result, err
}

// GetTaskVulListByTargetID 通过targetID信息任务漏洞信息列表及筛选
func (t *TaskVul) GetTaskVulListByTargetID(ctx context.Context, targetID int) ([]TaskVul, int64, error) {
	var (
		taskVulList []TaskVul
		count       int64
		db          = mysql.FromContext(ctx).Model(&TaskVul{}).Select("id,target_url,name,type,risk,location,status,verify_type,test_status,create_time")
	)
	if targetID != 0 {
		db = db.Where("target_id = ?", targetID)
	}
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("risk ASC").Find(&taskVulList).Error; err != nil {
		return nil, 0, err
	}
	return taskVulList, count, nil
}

// GetTaskVulInfoByTargetID 通过targetID信息任务漏洞信息列表及筛选
func (t *TaskVul) GetTaskVulInfoByTargetID(ctx context.Context, targetID int) (TaskVul, int64, error) {
	var (
		taskVulInfo TaskVul
		count       int64
		db          = mysql.FromContext(ctx).Model(&TaskVul{}).Select("id,target_url,name,type,risk,location,status,verify_type,test_status,create_time")
	)
	if targetID != 0 {
		db = db.Where("target_id = ?", targetID)
	}
	if err := db.Count(&count).Error; err != nil {
		return TaskVul{}, 0, err
	}
	if err := db.Order("risk ASC").First(&taskVulInfo).Error; err != nil {
		return TaskVul{}, 0, err
	}
	return taskVulInfo, count, nil
}

// GetSimpleTaskVulList 获取任务漏洞列表（不去重，展示原始明细）
func (t *TaskVul) GetSimpleTaskVulList(ctx context.Context, vulName, ip string, page, limit int) ([]TaskVul, int64, error) {
	var (
		vulList    []TaskVul // 直接使用原始模型接收
		count      int64
		db         = mysql.FromContext(ctx)
		queryWhere string
		argsWhere  []interface{}
	)
	queryWhere += "1 = 1"
	if ip != "" {
		queryWhere += " and (target_url LIKE ? )"
		argsWhere = append(argsWhere, "%"+ip+"%")
	}
	if vulName != "" {
		queryWhere += " and (name LIKE ? )"
		argsWhere = append(argsWhere, "%"+vulName+"%")
	}
	query := db.Model(&TaskVul{}).Where(queryWhere, argsWhere...)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.
		Order("id DESC").
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&vulList).Error; err != nil {
		return nil, 0, err
	}
	return vulList, count, nil
}

type DeduplicatedVul struct {
	TaskVul              // 嵌入 TaskVul 的字段
	Count         int    `gorm:"column:count"`
	TargetUrlList string `gorm:"column:target_url_list"`
}

// NOTE: 当 pocname (POC名称)、name (漏洞名称)、type (类型)、risk (风险等级) 这四个字段都相同时，这就是同一个漏洞
// GetTaskDeduplicationVulList 任务漏洞去重信息列表及筛选
func (t *TaskVul) GetTaskDeduplicationVulList(ctx context.Context, taskId, targetId, vultype, risk int, search, ip string, dataType, page, limit, status int) ([]DeduplicatedVul, int64, error) {
	var (
		deduplicatedList []DeduplicatedVul // 接收结果的列表
		count            int64
		db               = mysql.FromContext(ctx)
		queryWhere       string
		argsWhere        []interface{}
	)
	queryWhere += "1 = 1"
	if taskId != 0 {
		queryWhere += " and task_id = ? "
		argsWhere = append(argsWhere, taskId)
	}
	if dataType != 0 {
		queryWhere += " and data_type = ? "
		argsWhere = append(argsWhere, dataType)
	}
	if targetId > 0 {
		queryWhere += " and target_id = ? "
		argsWhere = append(argsWhere, targetId)
	}
	if vultype > 0 {
		queryWhere += " and type = ? "
		argsWhere = append(argsWhere, vultype)
	}
	if risk > 0 {
		queryWhere += " and risk = ? "
		argsWhere = append(argsWhere, risk)
	}
	if ip != "" {
		queryWhere += " and (target_url LIKE ? )"
		argsWhere = append(argsWhere, "%"+ip+"%")
	}
	if len(search) > 0 {
		queryWhere += " and (target_url LIKE ? or name LIKE ? or location LIKE ?)"
		argsWhere = append(argsWhere, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status > 0 {
		queryWhere += " and status = ? "
		argsWhere = append(argsWhere, status)
	}
	// group：按漏洞指纹(pocname+name+type+risk)分组，计算出每组的统计信息和最新ID
	// MAX(id): 获取该组最新的记录ID
	// GROUP_CONCAT: 聚合该组所有受影响的 URL，用于前端展示影响范围
	// COUNT(*): 计算该组漏洞发现的总次数
	subQuery := db.Model(&TaskVul{}).
		Select("MAX(id) as max_id, GROUP_CONCAT(DISTINCT target_url) AS target_url_list, COUNT(*) AS count, pocname, name, type, risk").
		Where(queryWhere, argsWhere...).
		Group("pocname, name, type, risk")
	selectFields := "tv_1.id, tv_1.pocname, tv_1.name, tv_1.type, tv_1.risk, tv_1.cvss, tv_1.status, tv_1.description, tv_1.create_time, tv_1.update_time, latest_vul.count, latest_vul.target_url_list"
	// 通过 JOIN 子查询的结果，获取最新那条记录(max_id)的完整详情(desc, time等)
	mainQuery := db.Model(&TaskVul{}).
		Table("task_vul AS tv_1").
		Select(selectFields).
		Joins(`JOIN (?) AS latest_vul
             ON tv_1.pocname = latest_vul.pocname
             AND tv_1.name = latest_vul.name
             AND tv_1.type = latest_vul.type
             AND tv_1.risk = latest_vul.risk
             AND tv_1.id = latest_vul.max_id`, subQuery)
	if err := mainQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	// 通过order id desc 保证最新的漏洞显示在最前面
	if err := mainQuery.
		Limit(limit).
		Offset(limit * (page - 1)).
		Order("tv_1.id DESC").
		Find(&deduplicatedList).Error; err != nil {
		return nil, 0, err
	}
	return deduplicatedList, count, nil
}

// GetAllTaskDeduplicationVulList 所有漏洞去重信息列表及筛选
func (t *TaskVul) GetAllTaskDeduplicationVulList(ctx context.Context, taskID int, uid, role int) ([]DeduplicatedVul, int64, error) {
	var (
		deduplicatedList []DeduplicatedVul // 接收结果的列表
		count            int64
		db               = mysql.FromContext(ctx)
	)
	subQuery := db.Model(&TaskVul{}).
		Select("MAX(id) as max_id, GROUP_CONCAT(DISTINCT target_url) AS target_url_list, COUNT(*) AS count, pocname, name, type, risk")

	// 普通用户只能看自己的数据
	if role == enums.UserRoleOrdinary {
		subQuery = subQuery.Joins("JOIN task_task ON task_vul.task_id = task_task.id").Where("task_task.user_id = ?", uid)
	}

	subQuery = subQuery.Group("pocname, name, type, risk")

	selectFields := "tv_1.id, tv_1.pocname, tv_1.name, tv_1.type, tv_1.risk, tv_1.cvss, tv_1.status, tv_1.description, tv_1.create_time, tv_1.update_time, latest_vul.count, latest_vul.target_url_list"
	mainQuery := db.Model(&TaskVul{}).
		Table("task_vul AS tv_1").
		Select(selectFields).
		Joins(`JOIN (?) AS latest_vul
             ON tv_1.pocname = latest_vul.pocname
             AND tv_1.name = latest_vul.name
             AND tv_1.type = latest_vul.type
             AND tv_1.risk = latest_vul.risk
             AND tv_1.id = latest_vul.max_id`, subQuery)

	if err := mainQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := mainQuery.
		Order("tv_1.id DESC").
		Find(&deduplicatedList).Error; err != nil {
		return nil, 0, err
	}
	return deduplicatedList, count, nil
}

// GetDuplicateTaskIds 获取指定漏洞指纹的重复记录 TaskID 列表
func (t *TaskVul) GetDuplicateTaskIds(ctx context.Context, taskId, targetId, vultype, risk int, search, ip string, dataType, status int,
	currentPocName, currentName string, currentType, currentRisk int) ([]int64, int64, error) {
	var (
		taskIds    []int64 // 接收 TaskID 的切片
		count      int64
		db         = mysql.FromContext(ctx)
		queryWhere string
		argsWhere  []interface{}
	)
	queryWhere += "1 = 1"

	if taskId != 0 {
		queryWhere += " and task_id = ? "
		argsWhere = append(argsWhere, taskId)
	}
	if dataType != 0 {
		queryWhere += " and data_type = ? "
		argsWhere = append(argsWhere, dataType)
	}
	if targetId > 0 {
		queryWhere += " and target_id = ? "
		argsWhere = append(argsWhere, targetId)
	}
	if vultype > 0 {
		queryWhere += " and type = ? "
		argsWhere = append(argsWhere, vultype)
	}
	if risk > 0 {
		queryWhere += " and risk = ? "
		argsWhere = append(argsWhere, risk)
	}
	if ip != "" {
		queryWhere += " and (target_url LIKE ? )"
		argsWhere = append(argsWhere, "%"+ip+"%")
	}
	if len(search) > 0 {
		queryWhere += " and (target_url LIKE ? or name LIKE ? or location LIKE ?)"
		argsWhere = append(argsWhere, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status > 0 {
		queryWhere += " and status = ? "
		argsWhere = append(argsWhere, status)
	}
	// 特定的漏洞名称、POC、类型和风险等级
	queryWhere += " AND pocname = ? AND name = ? AND type = ? AND risk = ?"
	argsWhere = append(argsWhere, currentPocName, currentName, currentType, currentRisk)
	err := db.Model(&TaskVul{}).
		Where(queryWhere, argsWhere...).
		Order("id DESC").
		Pluck("task_id", &taskIds).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}
	return taskIds, count, nil
}

// GetTaskDeduplicationByLocationList 按照 name, pocname, location 进行去重
// 当漏洞名称、POC名称和具体漏洞位置（location）都一致时，视为同一个漏洞进行合并
func (t *TaskVul) GetTaskDeduplicationByLocationList(ctx context.Context, taskId, targetId, vultype, risk int, search, ip, location string, dataType, page, limit, status, verifyType, uid, role int) ([]DeduplicatedVul, int64, error) {
	var (
		deduplicatedList []DeduplicatedVul
		count            int64
		db               = mysql.FromContext(ctx)
		argsWhere        []interface{}
		queryBuilder     strings.Builder
	)

	queryBuilder.WriteString("1 = 1")
	if taskId != 0 {
		queryBuilder.WriteString(" AND task_vul.task_id = ?")
		argsWhere = append(argsWhere, taskId)
	}
	if dataType != 0 {
		queryBuilder.WriteString(" AND task_vul.data_type = ?")
		argsWhere = append(argsWhere, dataType)
	}
	if targetId > 0 {
		queryBuilder.WriteString(" AND task_vul.target_id = ?")
		argsWhere = append(argsWhere, targetId)
	}
	if vultype > 0 {
		queryBuilder.WriteString(" AND task_vul.type = ?")
		argsWhere = append(argsWhere, vultype)
	}
	if risk > 0 {
		queryBuilder.WriteString(" AND task_vul.risk = ?")
		argsWhere = append(argsWhere, risk)
	}
	if status > 0 {
		queryBuilder.WriteString(" AND task_vul.status = ?")
		argsWhere = append(argsWhere, status)
	}
	if verifyType > 0 {
		queryBuilder.WriteString(" AND task_vul.test_status = ?")
		argsWhere = append(argsWhere, verifyType)
	}
	if ip != "" {
		queryBuilder.WriteString(" AND task_vul.target_url LIKE ?")
		argsWhere = append(argsWhere, "%"+ip+"%")
	}
	if location != "" {
		queryBuilder.WriteString(" AND task_vul.location LIKE ?")
		argsWhere = append(argsWhere, "%"+location+"%")
	}
	if len(search) > 0 {
		queryBuilder.WriteString(" AND task_vul.name LIKE ?")
		argsWhere = append(argsWhere, "%"+search+"%")
	}
	whereSQL := queryBuilder.String()

	subQueryDB := db.Model(&TaskVul{})
	// 普通用户只能看自己的数据
	if role == enums.UserRoleOrdinary {
		subQueryDB = subQueryDB.Joins("JOIN task_task ON task_vul.task_id = task_task.id").Where("task_task.user_id = ?", uid)
	}

	subQuery := subQueryDB.
		Select("MAX(task_vul.id) as max_id, GROUP_CONCAT(DISTINCT task_vul.target_url) AS target_url_list, COUNT(*) AS count").
		Where(whereSQL, argsWhere...).
		Group("task_vul.name, task_vul.pocname, task_vul.location")
	countResult := db.Table("(?) as tmp", subQuery)
	if err := countResult.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return []DeduplicatedVul{}, 0, nil
	}
	selectFields := "tv_1.id, tv_1.pocname, tv_1.name, tv_1.type, tv_1.risk, tv_1.status, tv_1.location, tv_1.test_status, tv_1.create_time, tv_1.update_time, latest_vul.count, latest_vul.target_url_list"
	err := db.Model(&TaskVul{}).
		Table("task_vul AS tv_1").
		Select(selectFields).
		Joins("JOIN (?) AS latest_vul ON tv_1.id = latest_vul.max_id", subQuery).
		Order("tv_1.id DESC").
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&deduplicatedList).Error
	if err != nil {
		return nil, 0, err
	}

	return deduplicatedList, count, nil
}
