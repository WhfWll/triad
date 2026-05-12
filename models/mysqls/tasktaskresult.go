package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

// TaskTaskResult 通用结果数据表
type TaskTaskResult struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`       // 主键
	ObjType    int       `gorm:"column:obj_type" json:"objType"`        // 数据类型：1-攻击面，2-证据，3-会话
	SubObjType string    `gorm:"column:sub_obj_type" json:"subObjType"` // 数据子类型，1-1:开放端口，1-2:敏感路径，1-3:URL,1-4:子域名，1-5：邮箱
	ObjID      string    `gorm:"column:obj_id" json:"objID"`            // 数据对象id
	SubObjID   string    `gorm:"column:sub_obj_id" json:"subObjID"`     // 数据子对象id
	Identify   string    `gorm:"column:identify" json:"identify"`       // 数据标识符
	Field1     string    `gorm:"column:field1" json:"field1"`           // 筛选字段1
	Field2     string    `gorm:"column:field2" json:"field2"`           // 筛选字段2
	Field3     string    `gorm:"column:field3" json:"field3"`           // 筛选字段3
	Field4     string    `gorm:"column:field4" json:"field4"`           // 筛选字段4
	JSONResult string    `gorm:"column:json_result" json:"jsonresult"`  // 各数据类型的json格式结果
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`  // 创建时间
}

// TableName sets insert table name for this struct type
func (t *TaskTaskResult) TableName() string {
	return "task_task_result"
}

// GetTaskTaskResultList 任务信息列表及筛选
func (t *TaskTaskResult) GetTaskTaskResultList(ctx context.Context, objType int, subObjType string, objId string, search string, page, limit int) ([]TaskTaskResult, int64, error) {
	var (
		taskTaskResultList []TaskTaskResult
		count              int64
		db                 = mysql.FromContext(ctx).Model(&TaskTaskResult{})
		query              string
		args               []interface{}
	)
	query += "obj_type = ? and sub_obj_type = ? and obj_id = ?"
	args = append(args, objType, subObjType, objId)
	if len(search) > 0 {
		query += " and (field1 LIKE ? or field2 LIKE ? or field3 LIKE ? or field4 LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Find(&taskTaskResultList)
	return taskTaskResultList, count, nil
}

// GetTaskTaskResultListBySubObjIds 根据SubObjIDs获取数据
// subObjIds必须是一个[]int/[]string数组
func (t *TaskTaskResult) GetTaskTaskResultListBySubObjIds(ctx context.Context, subObjIds interface{}) ([]TaskTaskResult, error) {
	var (
		taskTaskResultList []TaskTaskResult
		db                 = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)
	db.Where("sub_obj_id IN ?", subObjIds).Find(&taskTaskResultList)
	return taskTaskResultList, nil
}

// GetTaskTargetResultByType 根据subobjid及类型获取数据
func (t *TaskTaskResult) GetTaskTaskResultByType(ctx context.Context, objType int, subObjType string, subobjID any) ([]TaskTaskResult, error) {
	var (
		taskTaskResultList []TaskTaskResult
		db                 = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)
	db.Where("obj_type = ? and sub_obj_type = ? and sub_obj_id IN ?", objType, subObjType, subobjID).Find(&taskTaskResultList)
	return taskTaskResultList, nil
}

// Get retrieves a single record of taskTaskResult from database
func (t *TaskTaskResult) GetTaskTaskResult(ctx context.Context) (TaskTaskResult, error) {
	var (
		taskTaskResult TaskTaskResult
		err            error
		db             = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)

	curErr := db.Where("id = ?", t.ID).First(&taskTaskResult).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return taskTaskResult, err
}

// Add persists taskTaskResult to database
func (t *TaskTaskResult) AddTaskTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	if err := db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

// Update changes taskTaskResult by id
func (t *TaskTaskResult) UpdateTaskTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskResult{})

	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTaskTaskResult 根据id删除信息收集数据
// ids必须为[]int/[]string
func (t *TaskTaskResult) DeleteTaskTaskResult(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	if err := db.Where("id IN ?", ids).Delete(&TaskTaskResult{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTaskTaskResultBySubobjids 根据subobjids删除信息收集数据
// subobjIds[]int/[]string
func (t *TaskTaskResult) DeleteTaskTaskResultBySubobjids(ctx context.Context, subobjIds any) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	if err := db.Where("sub_obj_id IN ?", subobjIds).Delete(&TaskTaskResult{}).Error; err != nil {
		return err
	}
	return nil
}

// 批量删除 通过task_ids
func (t *TaskTaskResult) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	if err := db.Where("obj_id in ?", taskIds).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

// GetTaskTaskResultListByTaskId 任务信息列表及筛选
func (t *TaskTaskResult) GetTaskTaskResultListByTaskId(ctx context.Context, objId int) ([]TaskTaskResult, error) {
	var (
		taskTaskResultList []TaskTaskResult
		db                 = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)
	db.Where("obj_id = ?", objId).Find(&taskTaskResultList)
	return taskTaskResultList, nil
}

// GetTaskInfoByTargetIdAndSubObjType 任务信息列表及筛选
func (t *TaskTaskResult) GetTaskInfoByTargetIdAndSubObjType(ctx context.Context, targetId int, subObjType string) []TaskTaskResult {
	var (
		taskTaskResultList []TaskTaskResult
		db                 = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)
	db.
		Where("obj_type = ?", enums.TaskResultObjTypeInfo).
		Where("sub_obj_type = ?", subObjType).
		Where("sub_obj_id = ?", targetId).
		Find(&taskTaskResultList)
	return taskTaskResultList
}

// GetByTaskIdAndSubObjType 任务信息列表及筛选
func (t *TaskTaskResult) GetByTaskIdAndSubObjType(ctx context.Context, taskId int, subObjType string) []TaskTaskResult {
	var (
		taskTaskResultList []TaskTaskResult
		db                 = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)
	db.
		Where("obj_type = ?", enums.TaskResultObjTypeInfo).
		Where("sub_obj_type = ?", subObjType).
		Where("obj_id = ?", taskId).
		Find(&taskTaskResultList)
	return taskTaskResultList
}

// GetByTaskIdAndSubObjID 通过subobjID获取任务信息结果
func (t *TaskTaskResult) GetByTaskIdAndSubObjID(ctx context.Context, objId, objType string) TaskTaskResult {
	var (
		taskTaskResultList TaskTaskResult
		db                 = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)
	db.
		Where("obj_id = ?", objId).
		Where("sub_obj_type = ?", objType).
		First(&taskTaskResultList)
	return taskTaskResultList
}

// All 获取所有结果信息
func (t *TaskTaskResult) All(ctx context.Context, filter string) ([]TaskTaskResult, error) {
	var (
		taskResultList []TaskTaskResult
		db             = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)
	if filter != "" {
		db.Where(filter).Find(&taskResultList)
	} else {
		db.Find(&taskResultList)
	}
	return taskResultList, nil
}

// InfoStatsResult 信息统计结果
type InfoStatsResult struct {
	SubObjType int    `json:"sub_obj_type"`
	Field      string `json:"field"`
	Count      int    `json:"count"`
}

// GetInfoStatsByTaskId 获取任务信息统计（聚合查询）
func (t *TaskTaskResult) GetInfoStatsByTaskId(ctx context.Context, taskId int) map[int][]InfoStatsResult {
	var (
		results = make(map[int][]InfoStatsResult)
		db      = mysql.FromContext(ctx).Model(&TaskTaskResult{})
	)

	// 统计端口
	var portStats []InfoStatsResult
	db.Select("field2 as field, COUNT(*) as count").
		Where("task_id = ? AND obj_type = ? AND sub_obj_type = ? AND field2 != ''",
			taskId, enums.TaskResultObjTypeInfo, enums.TaskResultSubObjTypeService).
		Group("field2").
		Order("count DESC").
		Limit(20).
		Find(&portStats)
	results[1] = portStats

	// 统计服务
	var serviceStats []InfoStatsResult
	db.Select("field3 as field, COUNT(*) as count").
		Where("task_id = ? AND obj_type = ? AND sub_obj_type = ? AND field3 != ''",
			taskId, enums.TaskResultObjTypeInfo, enums.TaskResultSubObjTypeService).
		Group("field3").
		Order("count DESC").
		Limit(15).
		Find(&serviceStats)
	results[2] = serviceStats

	// 统计子域名
	var domainStats []InfoStatsResult
	db.Select("field1 as field, COUNT(*) as count").
		Where("task_id = ? AND obj_type = ? AND sub_obj_type = ? AND field1 != ''",
			taskId, enums.TaskResultObjTypeInfo, enums.TaskResultSubObjTypeSubdomain).
		Group("field1").
		Order("count DESC").
		Limit(10).
		Find(&domainStats)
	results[3] = domainStats

	return results
}
