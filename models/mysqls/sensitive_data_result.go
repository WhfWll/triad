package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

type SensitiveDataResult struct {
	ID           int       `gorm:"column:id;primary_key" json:"id"`
	TaskID       int       `gorm:"column:task_id" json:"taskId"`
	TargetID     int       `gorm:"column:target_id" json:"targetId"`
	TargetIP     string    `gorm:"column:target_ip" json:"targetIp"`
	DBType       int       `gorm:"column:db_type" json:"dbType"`
	DBName       string    `gorm:"column:db_name" json:"dbName"`
	TableNameStr string    `gorm:"column:table_name" json:"tableName"`
	ColumnName   string    `gorm:"column:column_name" json:"columnName"`
	DataType     int       `gorm:"column:data_type" json:"dataType"`
	DataLevel    int       `gorm:"column:data_level" json:"dataLevel"`
	SampleData   string    `gorm:"column:sample_data" json:"sampleData"`
	MatchRule    string    `gorm:"column:match_rule" json:"matchRule"`
	MatchType    int       `gorm:"column:match_type" json:"matchType"`
	TotalRows    int64     `gorm:"column:total_rows" json:"totalRows"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
}

func (SensitiveDataResult) TableName() string {
	return "sensitive_data_result"
}

func (s *SensitiveDataResult) Add(ctx context.Context) error {
	return mysql.FromContext(ctx).Model(s).Create(s).Error
}

func (s *SensitiveDataResult) BatchAdd(ctx context.Context, list []SensitiveDataResult) error {
	if len(list) == 0 {
		return nil
	}
	return mysql.FromContext(ctx).Model(s).CreateInBatches(list, 100).Error
}

func (s *SensitiveDataResult) GetByTargetID(ctx context.Context, targetID int) ([]SensitiveDataResult, error) {
	var list []SensitiveDataResult
	err := mysql.FromContext(ctx).Model(s).Where("target_id = ?", targetID).Order("data_level asc").Find(&list).Error
	return list, err
}

func (s *SensitiveDataResult) GetByTaskID(ctx context.Context, taskID int) ([]SensitiveDataResult, error) {
	var list []SensitiveDataResult
	err := mysql.FromContext(ctx).Model(s).Where("task_id = ?", taskID).Order("data_level asc").Find(&list).Error
	return list, err
}

func (s *SensitiveDataResult) DeleteByTaskID(ctx context.Context, taskID int) error {
	return mysql.FromContext(ctx).Model(s).Where("task_id = ?", taskID).Delete(nil).Error
}

func (s *SensitiveDataResult) GetStatByTaskID(ctx context.Context, taskID int) (highCount, middleCount, lowCount, total int64, err error) {
	db := mysql.FromContext(ctx).Model(s).Where("task_id = ?", taskID)
	db.Count(&total)
	db.Where("data_level = ?", 1).Count(&highCount)
	db.Where("data_level = ?", 2).Count(&middleCount)
	db.Where("data_level = ?", 3).Count(&lowCount)
	return
}
