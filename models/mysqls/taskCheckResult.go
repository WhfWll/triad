package mysqls

//
//import (
//	"context"
//	"gitlabee.4dogs.cn/common/mysql"
//	"gorm.io/gorm"
//	"time"
//)
//
//type TaskCheckresult struct {
//	ID         uint      `gorm:"column:id;primary_key" json:"id"`      // 主键
//	FatherID   string    `gorm:"column:father_id" json:"fatherID"`     // 父节点id
//	NodeID     string    `gorm:"column:node_id" json:"nodeID"`         // 节点id
//	TargetID   string    `gorm:"column:target_id" json:"targetID"`     // 目标id
//	Location   string    `gorm:"column:location" json:"location"`      // 漏洞位置
//	Pocname    string    `gorm:"column:pocname" json:"pocname"`        // 漏洞标识
//	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
//	Result     string    `gorm:"column:result" json:"result"`          // 漏洞结果
//	Request    string    `gorm:"column:request" json:"request"`        // 请求报文
//	Response   string    `gorm:"column:response" json:"response"`      // 响应报文
//}
//
//// TableName sets insert table name for this struct type
//func (t *TaskCheckresult) TableName() string {
//	return "task_checkresult"
//}
//
//// Get retrieves a list of taskCheckresult from database
//func (t *TaskCheckresult) GetTaskCheckresultList(ctx context.Context, page, limit int) ([]TaskCheckresult, int64, error) {
//	var (
//		taskCheckresultList []TaskCheckresult
//		count               int64
//		db                  = mysql.FromContext(ctx).Model(&TaskCheckresult{})
//	)
//
//	db.Limit(limit).Offset(limit * (page - 1)).Find(&taskCheckresultList)
//	db.Count(&count)
//
//	return taskCheckresultList, count, nil
//}
//
//// Get retrieves a single record of taskCheckresult from database
//func (t *TaskCheckresult) GetTaskCheckresult(ctx context.Context) (TaskCheckresult, error) {
//	var (
//		taskCheckresult TaskCheckresult
//		err             error
//		db              = mysql.FromContext(ctx).Model(&TaskCheckresult{})
//	)
//
//	curErr := db.Where("id = ?", t.ID).First(&taskCheckresult).Error
//	if curErr != nil && curErr != gorm.ErrRecordNotFound {
//		err = curErr
//	}
//
//	return taskCheckresult, err
//}
//
//// Add persists taskCheckresult to database
//func (t *TaskCheckresult) AddTaskCheckresult(ctx context.Context) error {
//	var db = mysql.FromContext(ctx).Model(&TaskCheckresult{})
//
//	if err := db.Create(t).Error; err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Update changes taskCheckresult by id
//func (t *TaskCheckresult) UpdateTaskCheckresult(ctx context.Context) error {
//	var db = mysql.FromContext(ctx).Model(&TaskCheckresult{})
//
//	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Delete taskCheckresult by id
//func (t *TaskCheckresult) DeleteTaskCheckresult(ctx context.Context) error {
//	var db = mysql.FromContext(ctx).Model(&TaskCheckresult{})
//
//	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
//		return err
//	}
//
//	return nil
//}
