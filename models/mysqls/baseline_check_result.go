package mysqls

import (
	"context"
	"fmt"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

type BaselineCheckResult struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`
	TaskID          int       `gorm:"column:task_id" json:"taskId"`
	TargetID        int       `gorm:"column:target_id" json:"targetId"`
	TargetIP        string    `gorm:"column:target_ip" json:"targetIp"`
	OSType          int       `gorm:"column:os_type" json:"osType"`
	ScanScene       int       `gorm:"column:scan_scene" json:"scanScene"`
	RuleID          int       `gorm:"column:rule_id" json:"ruleId"`
	RuleName        string    `gorm:"column:rule_name" json:"ruleName"`
	RuleCategory    int       `gorm:"column:rule_category" json:"ruleCategory"`
	RuleRisk        int       `gorm:"column:rule_risk" json:"ruleRisk"`
	CheckResult     int       `gorm:"column:check_result" json:"checkResult"`
	ExpectedValue   string    `gorm:"column:expected_value" json:"expectedValue"`
	ActualValue     string    `gorm:"column:actual_value" json:"actualValue"`
	CheckCommand    string    `gorm:"column:check_command" json:"checkCommand"`
	FixSuggestion   string    `gorm:"column:fix_suggestion" json:"fixSuggestion"`
	RiskDescription string    `gorm:"column:risk_description" json:"riskDescription"`
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`
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
	err := mysql.FromContext(ctx).Model(b).
		Where("target_id = ?", targetID).
		Where("check_result <> ?", 0).
		Order("rule_category asc").
		Find(&list).Error
	return list, err
}

func (b *BaselineCheckResult) GetByTaskID(ctx context.Context, taskID int) ([]BaselineCheckResult, error) {
	var list []BaselineCheckResult
	err := mysql.FromContext(ctx).Model(b).
		Where("task_id = ?", taskID).
		Where("check_result <> ?", 0).
		Order("rule_category asc").
		Find(&list).Error
	return list, err
}

func (b *BaselineCheckResult) DeleteByTaskID(ctx context.Context, taskID int) error {
	return mysql.FromContext(ctx).Model(b).Where("task_id = ?", taskID).Delete(nil).Error
}

func (b *BaselineCheckResult) DeleteByTaskIDAndTargetIP(ctx context.Context, taskID int, targetIP string) error {
	return mysql.FromContext(ctx).Model(b).Where("task_id = ? AND target_ip = ?", taskID, targetIP).Delete(nil).Error
}

func (b *BaselineCheckResult) DeleteByTaskTargetAndScene(ctx context.Context, taskID int, targetIP string, scanScene int) error {
	q := mysql.FromContext(ctx).Model(b).Where("task_id = ? AND target_ip = ?", taskID, targetIP)
	if scanScene > 0 {
		q = q.Where("scan_scene = ?", scanScene)
	}
	return q.Delete(nil).Error
}

func (b *BaselineCheckResult) CountByTaskTargetAndScene(ctx context.Context, taskID int, targetIP string, scanScene int) (int64, error) {
	var total int64
	q := mysql.FromContext(ctx).Model(b).Where("task_id = ? AND target_ip = ?", taskID, targetIP)
	if scanScene > 0 {
		q = q.Where("scan_scene = ?", scanScene)
	}
	err := q.Count(&total).Error
	return total, err
}

func (b *BaselineCheckResult) StopPendingByTaskTargetAndScene(ctx context.Context, taskID int, targetIP string, scanScene int, actualValue string) (int64, error) {
	q := mysql.FromContext(ctx).Model(b).
		Where("task_id = ? AND target_ip = ? AND check_result = ?", taskID, targetIP, 0)
	if scanScene > 0 {
		q = q.Where("scan_scene = ?", scanScene)
	}
	res := q.Updates(map[string]interface{}{
		"check_result": 4,
		"rule_name":    "任务已手动结束",
		"actual_value": actualValue,
		"create_time":  time.Now(),
	})
	return res.RowsAffected, res.Error
}

func (b *BaselineCheckResult) GetStatByTaskID(ctx context.Context, taskID int) (passCount, failCount, total int64, err error) {
	stat, err := b.GetDetailStatByTaskID(ctx, taskID)
	if err != nil || stat == nil {
		return
	}
	return stat.PassCount, stat.FailCount, stat.TotalRules, nil
}

// BaselineDetailStatRow 任务级核查统计（详情概况）
type BaselineDetailStatRow struct {
	TotalRules   int64 `gorm:"column:total_rules"`
	PassCount    int64 `gorm:"column:pass_count"`
	FailCount    int64 `gorm:"column:fail_count"`
	ErrorCount   int64 `gorm:"column:error_count"`
	SkipCount    int64 `gorm:"column:skip_count"`
	FailCritical int64 `gorm:"column:fail_critical"`
	FailHigh     int64 `gorm:"column:fail_high"`
	FailMiddle   int64 `gorm:"column:fail_middle"`
	FailLow      int64 `gorm:"column:fail_low"`
}

func (b *BaselineCheckResult) GetDetailStatByTaskID(ctx context.Context, taskID int) (*BaselineDetailStatRow, error) {
	var row BaselineDetailStatRow
	err := mysql.FromContext(ctx).Model(b).
		Select(`SUM(CASE WHEN check_result <> 0 THEN 1 ELSE 0 END) AS total_rules,
			SUM(CASE WHEN check_result = 1 THEN 1 ELSE 0 END) AS pass_count,
			SUM(CASE WHEN check_result = 2 THEN 1 ELSE 0 END) AS fail_count,
			SUM(CASE WHEN check_result = 3 THEN 1 ELSE 0 END) AS error_count,
			SUM(CASE WHEN check_result = 4 THEN 1 ELSE 0 END) AS skip_count,
			SUM(CASE WHEN check_result = 2 AND rule_risk = 0 THEN 1 ELSE 0 END) AS fail_critical,
			SUM(CASE WHEN check_result = 2 AND rule_risk = 1 THEN 1 ELSE 0 END) AS fail_high,
			SUM(CASE WHEN check_result = 2 AND rule_risk = 2 THEN 1 ELSE 0 END) AS fail_middle,
			SUM(CASE WHEN check_result = 2 AND rule_risk = 3 THEN 1 ELSE 0 END) AS fail_low`).
		Where("task_id = ?", taskID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

type BaselineFailCategoryRow struct {
	Category int   `gorm:"column:rule_category"`
	Count    int64 `gorm:"column:cnt"`
}

func (b *BaselineCheckResult) GetTopFailCategoriesByTaskID(ctx context.Context, taskID int, limit int) ([]BaselineFailCategoryRow, error) {
	if limit <= 0 {
		limit = 5
	}
	var rows []BaselineFailCategoryRow
	err := mysql.FromContext(ctx).Model(b).
		Select("rule_category, COUNT(*) AS cnt").
		Where("task_id = ? AND check_result = ?", taskID, 2).
		Group("rule_category").
		Order("cnt DESC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
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
	Pending    int64     `gorm:"column:pending"`
}

func (b *BaselineCheckResult) CountDistinctTasks(ctx context.Context, scanScene int) (int64, error) {
	var total int64
	tpl := "SELECT COUNT(*) FROM (SELECT 1 FROM baseline_check_result %sGROUP BY task_id, target_ip) AS t"
	var err error
	if scanScene > 0 {
		err = mysql.FromContext(ctx).Raw(fmt.Sprintf(tpl, "WHERE scan_scene = ? "), scanScene).Scan(&total).Error
		fmt.Printf(">>> CountDistinctTasks: scanScene=%d total=%d err=%v\n", scanScene, total, err)
	} else {
		err = mysql.FromContext(ctx).Raw(fmt.Sprintf(tpl, "")).Scan(&total).Error
		fmt.Printf(">>> CountDistinctTasks: scanScene=all total=%d err=%v\n", total, err)
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
			SUM(CASE WHEN check_result = 0 THEN 1 ELSE 0 END) AS pending,
			SUM(CASE WHEN check_result <> 0 THEN 1 ELSE 0 END) AS total_rules,
			SUM(CASE WHEN check_result = 1 THEN 1 ELSE 0 END) AS pass_count,
			SUM(CASE WHEN check_result = 2 THEN 1 ELSE 0 END) AS fail_count,
			SUM(CASE WHEN check_result = 3 THEN 1 ELSE 0 END) AS err_count`).
		Group("task_id, target_ip").
		Order("last_time DESC").
		Offset(offset).
		Limit(size).
		Scan(&rows).Error
	fmt.Printf(">>> ListGroupedByTask: page=%d size=%d scanScene=%d offset=%d rows=%d err=%v\n", page, size, scanScene, offset, len(rows), err)
	for i, r := range rows {
		fmt.Printf(">>> ListGroupedByTask: row[%d] taskID=%d targetIP=%s totalRules=%d pass=%d fail=%d err=%d lastTime=%v\n",
			i, r.TaskID, r.TargetIP, r.TotalRules, r.PassCount, r.FailCount, r.ErrCount, r.LastTime)
	}
	return rows, err
}

// BaselineTaskTargetRow 任务中的目标聚合行
type BaselineTaskTargetRow struct {
	TargetID   int    `gorm:"column:target_id"`
	TargetIP   string `gorm:"column:target_ip"`
	OSType     int    `gorm:"column:os_type"`
	TotalRules int64  `gorm:"column:total_rules"`
	PassCount  int64  `gorm:"column:pass_count"`
	FailCount  int64  `gorm:"column:fail_count"`
	ErrCount   int64  `gorm:"column:err_count"`
}

func (b *BaselineCheckResult) GetTargetsByTaskID(ctx context.Context, taskID int) ([]BaselineTaskTargetRow, error) {
	var rows []BaselineTaskTargetRow
	err := mysql.FromContext(ctx).Model(b).
		Select(`MAX(target_id) AS target_id,
			target_ip,
			MAX(os_type) AS os_type,
			SUM(CASE WHEN check_result <> 0 THEN 1 ELSE 0 END) AS total_rules,
			SUM(CASE WHEN check_result = 1 THEN 1 ELSE 0 END) AS pass_count,
			SUM(CASE WHEN check_result = 2 THEN 1 ELSE 0 END) AS fail_count,
			SUM(CASE WHEN check_result = 3 THEN 1 ELSE 0 END) AS err_count`).
		Where("task_id = ?", taskID).
		Group("target_ip").
		Order("target_ip ASC").
		Scan(&rows).Error
	return rows, err
}
