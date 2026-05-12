package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"time"
)

type DBCheckResult struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`
	TaskID          int       `gorm:"column:task_id" json:"taskId"`
	TargetID        int       `gorm:"column:target_id" json:"targetId"`
	TargetIP        string    `gorm:"column:target_ip" json:"targetIp"`
	DBType          int       `gorm:"column:db_type" json:"dbType"`
	DBName          string    `gorm:"column:db_name" json:"dbName"`
	CheckCategory   int       `gorm:"column:check_category" json:"checkCategory"`
	RuleID          int       `gorm:"column:rule_id" json:"ruleId"`
	RuleName        string    `gorm:"column:rule_name" json:"ruleName"`
	CheckResult     int       `gorm:"column:check_result" json:"checkResult"`
	ExpectedValue   string    `gorm:"column:expected_value" json:"expectedValue"`
	ActualValue     string    `gorm:"column:actual_value" json:"actualValue"`
	RiskLevel       int       `gorm:"column:risk_level" json:"riskLevel"`
	FixSuggestion   string    `gorm:"column:fix_suggestion" json:"fixSuggestion"`
	RiskDescription string    `gorm:"column:risk_description" json:"riskDescription"`
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`
}

func (DBCheckResult) TableName() string {
	return "db_check_result"
}

func (d *DBCheckResult) Add(ctx context.Context) error {
	return mysql.FromContext(ctx).Model(d).Create(d).Error
}

func (d *DBCheckResult) BatchAdd(ctx context.Context, list []DBCheckResult) error {
	if len(list) == 0 {
		return nil
	}
	return mysql.FromContext(ctx).Model(d).CreateInBatches(list, 100).Error
}

func (d *DBCheckResult) GetByTargetID(ctx context.Context, targetID int) ([]DBCheckResult, error) {
	var list []DBCheckResult
	err := mysql.FromContext(ctx).Model(d).Where("target_id = ?", targetID).Order("check_category asc").Find(&list).Error
	return list, err
}

func (d *DBCheckResult) GetByTaskID(ctx context.Context, taskID int) ([]DBCheckResult, error) {
	var list []DBCheckResult
	err := mysql.FromContext(ctx).Model(d).Where("task_id = ?", taskID).Order("check_category asc").Find(&list).Error
	return list, err
}

func (d *DBCheckResult) DeleteByTaskID(ctx context.Context, taskID int) error {
	return mysql.FromContext(ctx).Model(d).Where("task_id = ?", taskID).Delete(nil).Error
}
