package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Flowlog struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`       // 主键
	FlowTaskID int       `gorm:"column:flow_task_id" json:"flowTaskID"` // 所属流量分析任务id
	Content    string    `gorm:"column:content" json:"content"`         // 日志内容
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`  // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`  // 修改时间
}

// TableName sets insert table name for this struct type
func (f *Flowlog) TableName() string {
	return "flow_log"
}

// Get retrieves a list of flowlog from database
func (f *Flowlog) GetFlowlogList(ctx context.Context, search string, flowTaskId int, page, limit int) ([]Flowlog, int64, error) {
	var (
		flowlogList []Flowlog
		count       int64
		db          = mysql.FromContext(ctx).Model(&Flowlog{})
		query       string
		args        []interface{}
	)
	query += "flow_task_id = ?"
	args = append(args, flowTaskId)
	if len(search) > 0 {
		query += " and content LIKE ?"
		args = append(args, "%"+search+"%")
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&flowlogList)
	return flowlogList, count, nil
}

// Get retrieves a single record of flowlog from database
func (f *Flowlog) GetFlowlog(ctx context.Context) (Flowlog, error) {
	var (
		flowlog Flowlog
		err     error
		db      = mysql.FromContext(ctx).Model(&Flowlog{})
	)

	curErr := db.Where("id = ?", f.ID).First(&flowlog).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return flowlog, err
}

//AddFlowlog 新增日志
func (f *Flowlog) AddFlowlog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Flowlog{})
	if err := db.Create(f).Error; err != nil {
		return err
	}
	return nil
}

// Update changes flowlog by id
func (f *Flowlog) UpdateFlowlog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Flowlog{})

	if err := db.Where("id = ?", f.ID).Updates(f).Error; err != nil {
		return err
	}

	return nil
}

// 批量删除 通过flow_task_id
func (f *Flowlog) DeleteByFlowTaskIds(ctx context.Context, flowTaskId any) error {
	var db = mysql.FromContext(ctx).Model(&Flowlog{})
	if err := db.Where("flow_task_id IN ?", flowTaskId).Delete(f).Error; err != nil {
		return err
	}
	return nil
}
