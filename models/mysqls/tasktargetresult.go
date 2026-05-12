package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type TaskTargetResult struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`       // 主键
	ObjType    int       `gorm:"column:obj_type" json:"objType"`        // 数据类型
	SubObjType string    `gorm:"column:sub_obj_type" json:"subObjType"` // 数据子类型
	ObjID      string    `gorm:"column:obj_id" json:"objID"`            // 数据对象id
	SubObjID   string    `gorm:"column:sub_obj_id" json:"subObjID"`     // 数据子对象id
	Identify   string    `gorm:"column:identify" json:"identify"`       // 数据标识符
	Field1     string    `gorm:"column:field1" json:"field1"`           // 筛选字段1
	Field2     string    `gorm:"column:field2" json:"field2"`           // 筛选字段2
	Field3     string    `gorm:"column:field3" json:"field3"`           // 筛选字段3
	JSONResult string    `gorm:"column:json_result" json:"jsonresult"`  // 各数据类型的json格式结果
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`  // 创建时间
}

// TableName sets insert table name for this struct type
func (t *TaskTargetResult) TableName() string {
	return "task_target_result"
}

// Get retrieves a list of taskTargetResult from database
//ObjIDs必须是一个[]int/[]string数组
func (t *TaskTargetResult) GetTaskTargetResultList(ctx context.Context, ObjIDs interface{}) ([]TaskTargetResult, error) {
	var (
		taskTargetResultList []TaskTargetResult
		db                   = mysql.FromContext(ctx).Model(&TaskTargetResult{})
	)
	db.Where("obj_id IN ?", ObjIDs).Find(&taskTargetResultList)
	return taskTargetResultList, nil
}

// Get retrieves a single record of taskTargetResult from database
func (t *TaskTargetResult) GetTaskTargetResult(ctx context.Context) (TaskTargetResult, error) {
	var (
		taskTargetResult TaskTargetResult
		err              error
		db               = mysql.FromContext(ctx).Model(&TaskTargetResult{})
	)

	curErr := db.Where("id = ?", t.ID).First(&taskTargetResult).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return taskTargetResult, err
}

// Add persists taskTargetResult to database
func (t *TaskTargetResult) AddTaskTargetResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTargetResult{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// Update changes taskTargetResult by id
func (t *TaskTargetResult) UpdateTaskTargetResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskTargetResult{})

	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

//DeleteByIds 通过id删除数据
//ids必须为[]int/[]string
func (t *TaskTargetResult) DeleteByIds(ctx context.Context, ids interface{}) error {
	var db = mysql.FromContext(ctx).Model(&TaskTargetResult{})
	if err := db.Where("id in ?", ids).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

//DeleteByObjIds 通过objids删除数据
//ids必须为[]int/[]string
func (t *TaskTargetResult) DeleteByObjIds(ctx context.Context, objIds any) error {
	var db = mysql.FromContext(ctx).Model(&TaskTargetResult{})
	if err := db.Where("obj_id in ?", objIds).Delete(t).Error; err != nil {
		return err
	}
	return nil
}

// 批量删除 通过task_ids
func (t *TaskTargetResult) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&TaskTargetResult{})
	if err := db.Where("obj_id in ?", taskIds).Delete(t).Error; err != nil {
		return err
	}

	return nil
}
