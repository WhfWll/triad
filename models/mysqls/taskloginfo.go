package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Taskloginfo struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	TaskLogID  int       `gorm:"column:task_log_id" json:"taskLogID"`  // 所属任务日志id
	TaskID     int       `gorm:"column:task_id" json:"taskID"`         // 所属任务id
	TargetID   int       `gorm:"column:target_id" json:"targetID"`     // 所属目标id
	TargetURL  string    `gorm:"column:target_url" json:"targetURL"`   // 测试目标地址
	Pocname    string    `gorm:"column:pocname" json:"pocname"`        // pocname
	Result     string    `gorm:"column:result" json:"result"`          // 日志内容
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
}

// TableName sets insert table name for this struct type
func (t *Taskloginfo) TableName() string {
	return "task_log_info"
}

// GetTaskloginfoListTaskIds 根据目标ids获取日志详情数据
// 参数targetIds必须是一个[]int或[]string数组
func (t *Taskloginfo) GetTaskloginfoListTaskIds(ctx context.Context, targetIds interface{}) ([]Taskloginfo, error) {
	var (
		taskloginfoList []Taskloginfo
		db              = mysql.FromContext(ctx).Model(&Taskloginfo{})
	)
	db.Where("target_id IN ?", targetIds).Find(&taskloginfoList)
	return taskloginfoList, nil
}

// GetTaskloginfoListById 查询测试日志详情
func (t *Taskloginfo) GetTaskloginfoListById(ctx context.Context, taskLogId int) ([]Taskloginfo, int64, error) {
	var (
		taskloginfoList []Taskloginfo
		count           int64
		db              = mysql.FromContext(ctx).Model(&Taskloginfo{})
	)

	db.Count(&count)
	db.Where("task_log_id = ?", taskLogId).Find(&taskloginfoList)
	return taskloginfoList, count, nil
}

// Get retrieves a single record of taskloginfo from database
func (t *Taskloginfo) GetTaskloginfo(ctx context.Context) (Taskloginfo, error) {
	var (
		taskloginfo Taskloginfo
		err         error
		db          = mysql.FromContext(ctx).Model(&Taskloginfo{})
	)

	curErr := db.Where("id = ?", t.ID).First(&taskloginfo).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return taskloginfo, err
}

// Add persists taskloginfo to database
func (t *Taskloginfo) AddTaskloginfo(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Taskloginfo{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTaskloginfoByTargetIds 通过TargetIds删除数据
// targetIds必须为[]int/[]string
func (t *Taskloginfo) DeleteTaskloginfoByTargetIds(ctx context.Context, targetIds any) error {
	var db = mysql.FromContext(ctx).Model(&Taskloginfo{})
	if err := db.Where("target_id IN ?", targetIds).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

// 批量删除 通过task_ids
func (t *Taskloginfo) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&Taskloginfo{})
	if err := db.Where("task_id in ?", taskIds).Delete(t).Error; err != nil {
		return err
	}

	return nil
}

type GetTaskLogInfoNumberByTaskIdRes struct {
	Num      int `json:"num"`
	TargetId int `json:"targetId"`
}

func (t *Taskloginfo) GetTaskLogInfoNumberByTaskId(ctx context.Context, taskId int) (result []GetTaskLogInfoNumberByTaskIdRes, err error) {
	var (
		db = mysql.FromContext(ctx).Model(&Taskloginfo{})
	)
	db.Select("count(id) as num,target_id").Where("task_id = ?", taskId).Group("target_id").Find(&result)
	return result, nil
}

func (t *Taskloginfo) GetTaskLogInfoListByLogIds(ctx context.Context, logIds []int) ([]Taskloginfo, error) {
	var (
		db              = mysql.FromContext(ctx).Model(&Taskloginfo{})
		taskLogInfoList []Taskloginfo
	)
	db.Where("task_log_id in ?", logIds).Find(&taskLogInfoList)
	return taskLogInfoList, nil
}
