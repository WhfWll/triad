package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"time"
)

type BaselineCheckResult struct {
	ID               int       `gorm:"column:id;primary_key" json:"id"`
	TaskID           int       `gorm:"column:task_id" json:"taskId"`
	TargetID         int       `gorm:"column:target_id" json:"targetId"`
	TargetIP         string    `gorm:"column:target_ip" json:"targetIp"`
	OSType           int       `gorm:"column:os_type" json:"osType"`
	RuleID           int       `gorm:"column:rule_id" json:"ruleId"`
	RuleName         string    `gorm:"column:rule_name" json:"ruleName"`
	RuleCategory     int       `gorm:"column:rule_category" json:"ruleCategory"`
	RuleRisk         int       `gorm:"column:rule_risk" json:"ruleRisk"`
	CheckResult      int       `gorm:"column:check_result" json:"checkResult"`
	ExpectedValue    string    `gorm:"column:expected_value" json:"expectedValue"`
	ActualValue      string    `gorm:"column:actual_value" json:"actualValue"`
	CheckCommand     string    `gorm:"column:check_command" json:"checkCommand"`
	FixSuggestion    string    `gorm:"column:fix_suggestion" json:"fixSuggestion"`
	RiskDescription  string    `gorm:"column:risk_description" json:"riskDescription"`
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`
}

func (BaselineCheckResult) TableName() string {
	return "baseline_check_result"
}

func (b *BaselineCheckResult) Add(ctx context.Context) error {
	return mysql.FromContext(ctx).Model(b).Create(b).Error
}

func (b *BaselineCheckResult) BatchAdd(ctx context.Context, list []BaselineCheckResult) error {
	if len(list) == 0 {
		return nil
	}
	return mysql.FromContext(ctx).Model(b).CreateInBatches(list, 100).Error
}

func (b *BaselineCheckResult) GetByTargetID(ctx context.Context, targetID int) ([]BaselineCheckResult, error) {
	var list []BaselineCheckResult
	err := mysql.FromContext(ctx).Model(b).Where("target_id = ?", targetID).Order("rule_category asc").Find(&list).Error
	return list, err
}

func (b *BaselineCheckResult) GetByTaskID(ctx context.Context, taskID int) ([]BaselineCheckResult, error) {
	var list []BaselineCheckResult
	err := mysql.FromContext(ctx).Model(b).Where("task_id = ?", taskID).Order("rule_category asc").Find(&list).Error
	return list, err
}

func (b *BaselineCheckResult) DeleteByTaskID(ctx context.Context, taskID int) error {
	return mysql.FromContext(ctx).Model(b).Where("task_id = ?", taskID).Delete(nil).Error
}

func (b *BaselineCheckResult) GetStatByTaskID(ctx context.Context, taskID int) (passCount, failCount, total int64, err error) {
	db := mysql.FromContext(ctx).Model(b).Where("task_id = ?", taskID)
	db.Count(&total)
	db.Where("check_result = ?", 1).Count(&passCount)
	db.Where("check_result = ?", 2).Count(&failCount)
	return
}

func (b *BaselineCheckResult) GetStatByTargetID(ctx context.Context, targetID int) (passCount, failCount, total int64, err error) {
	db := mysql.FromContext(ctx).Model(b).Where("target_id = ?", targetID)
	db.Count(&total)
	db.Where("check_result = ?", 1).Count(&passCount)
	db.Where("check_result = ?", 2).Count(&failCount)
	return
}
