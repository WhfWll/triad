package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

type SecurityReport struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`
	Title      string    `gorm:"column:title" json:"title"`
	Module     string    `gorm:"column:module" json:"module"`
	TaskID     int       `gorm:"column:task_id" json:"taskId"`
	TaskName   string    `gorm:"column:task_name" json:"taskName"`
	Content    string    `gorm:"column:content" json:"content"`
	CreateBy   int       `gorm:"column:create_by" json:"createBy"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
}

func (SecurityReport) TableName() string {
	return "security_report"
}

func (m *SecurityReport) Add(ctx context.Context) error {
	return mysql.FromContext(ctx).Model(m).Create(m).Error
}

func (m *SecurityReport) DeleteByID(ctx context.Context, id int) error {
	return mysql.FromContext(ctx).Model(m).Where("id = ?", id).Delete(nil).Error
}

func (m *SecurityReport) GetByID(ctx context.Context, id int) (*SecurityReport, error) {
	var row SecurityReport
	err := mysql.FromContext(ctx).Model(m).Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

type SecurityReportListRow struct {
	ID         int       `gorm:"column:id" json:"id"`
	Title      string    `gorm:"column:title" json:"title"`
	Module     string    `gorm:"column:module" json:"module"`
	TaskID     int       `gorm:"column:task_id" json:"taskId"`
	TaskName   string    `gorm:"column:task_name" json:"taskName"`
	CreateBy   int       `gorm:"column:create_by" json:"createBy"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
}

func (m *SecurityReport) List(ctx context.Context, page, size int) ([]SecurityReportListRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	var total int64
	q := mysql.FromContext(ctx).Model(m).Order("create_time DESC")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	var rows []SecurityReportListRow
	if err := q.Offset(offset).Limit(size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
