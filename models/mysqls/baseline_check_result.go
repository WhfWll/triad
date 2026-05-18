package mysqls

import (
	"context"
	"fmt"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

type BaselineCheckResult struct {
	ID               int       `gorm:"column:id;primary_key" json:"id"`
	TaskID           int       `gorm:"column:task_id" json:"taskId"`
	TargetID         int       `gorm:"column:target_id" json:"targetId"`
	TargetIP         string    `gorm:"column:target_ip" json:"targetIp"`
	OSType           int       `gorm:"column:os_type" json:"osType"`
	ScanScene        int       `gorm:"column:scan_scene" json:"scanScene"`
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

// BaselineTaskGroupRow 按 task_id 聚合的一行（用于历史列表）
type BaselineTaskGroupRow struct {
	TaskID     int       `gorm:"column:task_id"`
	TargetIP   string    `gorm:"column:target_ip"`
	OSType     int       `gorm:"column:os_type"`
	ScanScene  int       `gorm:"column:scan_scene"`
	LastTime   time.Time `gorm:"column:last_time"`
	TotalRules int64     `gorm:"column:total_rules"`
	PassCount  int64     `gorm:"column:pass_count"`
	FailCount  int64     `gorm:"column:fail_count"`
	ErrCount   int64     `gorm:"column:err_count"`
}

func (b *BaselineCheckResult) CountDistinctTasks(ctx context.Context, scanScene int) (int64, error) {
	var total int64
	tpl := "SELECT COUNT(*) FROM (SELECT 1 FROM baseline_check_result %sGROUP BY task_id, target_ip) AS t"
	var err error
	if scanScene > 0 {
		err = mysql.FromContext(ctx).Raw(fmt.Sprintf(tpl, "WHERE scan_scene = ? "), scanScene).Scan(&total).Error
	} else {
		err = mysql.FromContext(ctx).Raw(fmt.Sprintf(tpl, "")).Scan(&total).Error
	}
	return total, err
}

func (b *BaselineCheckResult) ListGroupedByTask(ctx context.Context, page, size, scanScene int) ([]BaselineTaskGroupRow, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size
	var rows []BaselineTaskGroupRow
	q := mysql.FromContext(ctx).Model(b)
	if scanScene > 0 {
		q = q.Where("scan_scene = ?", scanScene)
	}
	err := q.Select(`task_id,
			target_ip,
			MAX(os_type) AS os_type,
			MAX(scan_scene) AS scan_scene,
			MAX(create_time) AS last_time,
			COUNT(*) AS total_rules,
			SUM(CASE WHEN check_result = 1 THEN 1 ELSE 0 END) AS pass_count,
			SUM(CASE WHEN check_result = 2 THEN 1 ELSE 0 END) AS fail_count,
			SUM(CASE WHEN check_result = 3 THEN 1 ELSE 0 END) AS err_count`).
		Group("task_id, target_ip").
		Order("last_time DESC").
		Offset(offset).
		Limit(size).
		Scan(&rows).Error
	return rows, err
}

// BaselineTaskTargetRow 任务中的目标聚合行
type BaselineTaskTargetRow struct {
	TargetID   int       `gorm:"column:target_id"`
	TargetIP   string    `gorm:"column:target_ip"`
	OSType     int       `gorm:"column:os_type"`
	TotalRules int64     `gorm:"column:total_rules"`
	PassCount  int64     `gorm:"column:pass_count"`
	FailCount  int64     `gorm:"column:fail_count"`
	ErrCount   int64     `gorm:"column:err_count"`
}

func (b *BaselineCheckResult) GetTargetsByTaskID(ctx context.Context, taskID int) ([]BaselineTaskTargetRow, error) {
	var rows []BaselineTaskTargetRow
	err := mysql.FromContext(ctx).Model(b).
		Select(`target_id,
			target_ip,
			MAX(os_type) AS os_type,
			COUNT(*) AS total_rules,
			SUM(CASE WHEN check_result = 1 THEN 1 ELSE 0 END) AS pass_count,
			SUM(CASE WHEN check_result = 2 THEN 1 ELSE 0 END) AS fail_count,
			SUM(CASE WHEN check_result = 3 THEN 1 ELSE 0 END) AS err_count`).
		Where("task_id = ?", taskID).
		Group("target_id, target_ip").
		Order("target_ip ASC").
		Scan(&rows).Error
	return rows, err
}
