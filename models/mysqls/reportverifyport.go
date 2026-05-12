package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type Reportverifyport struct {
	ID        int    `gorm:"column:id;primary_key" json:"id"`   // 主键
	Target    string `gorm:"column:target" json:"target"`       // 目标
	Port      string `gorm:"column:port" json:"port"`           // 端口
	Scheme    string `gorm:"column:scheme" json:"scheme"`       // 协议
	Service   string `gorm:"column:service" json:"service"`     // 服务
	Component string `gorm:"column:component" json:"component"` // 组件
	TaskId    int    `gorm:"column:task_id" json:"task_id"`     // 任务id
	TargetId  int    `gorm:"column:target_id" json:"target_id"` // 目标id
}

// TableName sets insert table name for this struct type
func (r *Reportverifyport) TableName() string {
	return "report_verify_port"
}

// Get retrieves a list of reportverifyport from database
func (r *Reportverifyport) GetReportverifyportList(ctx context.Context, taskId, page, limit int, search string) ([]Reportverifyport, int64, error) {
	var (
		reportverifyportList []Reportverifyport
		count                int64
		db                   = mysql.FromContext(ctx).Model(&Reportverifyport{})
		query                string
		args                 []interface{}
	)

	query += "task_id = ?"
	args = append(args, taskId)
	if search != "" {
		query += " and port LIKE ? or serivce like ? or component like ?"
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	dbs := db.Where(query, args...)
	dbs.Limit(limit).Offset(limit * (page - 1)).Find(&reportverifyportList)
	dbs.Count(&count)

	return reportverifyportList, count, nil
}

// Get retrieves a single record of reportverifyport from database
func (r *Reportverifyport) GetReportverifyport(ctx context.Context) (Reportverifyport, error) {
	var (
		reportverifyport Reportverifyport
		err              error
		db               = mysql.FromContext(ctx).Model(&Reportverifyport{})
	)

	curErr := db.Where("id = ?", r.ID).First(&reportverifyport).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return reportverifyport, err
}

// Add persists reportverifyport to database
func (r *Reportverifyport) AddReportverifyport(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifyport{})

	if err := db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

// Update changes reportverifyport by id
func (r *Reportverifyport) UpdateReportverifyport(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifyport{})

	if err := db.Where("id = ?", r.ID).Updates(r).Error; err != nil {
		return err
	}

	return nil
}

// Delete reportverifyport by id
func (r *Reportverifyport) DeleteReportverifyport(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifyport{})

	if err := db.Where("id = ?", r.ID).Updates(r).Error; err != nil {
		return err
	}

	return nil
}
