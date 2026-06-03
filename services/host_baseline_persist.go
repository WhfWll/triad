package services

import (
	"context"
	"fmt"
	"time"

	"smart/models/mysqls"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

// BaselineCheckStatusPending 在 baseline_check_result 表中标记一行"任务占位"记录（check_result=0），
// 用于在真实检查结果写入前让任务列表能立即看到该 (task_id, target_ip)。
const BaselineCheckStatusPending = 0

// PersistBaselineCheckPlaceholder 在创建批量核查任务时为每个目标写入一行占位记录，
// 确保任务列表接口（按 task_id+target_ip 聚合）在检查完成前就能看到该任务。
// 真实检查开始时 RunBaselineCheck 会先调用 DeleteByTaskIDAndTargetIP 清理占位行，
// 检查完成后通过 BatchAdd 写入真实规则结果，因此占位行无需单独清理。
func PersistBaselineCheckPlaceholder(ctx context.Context, taskID, targetID int, targetIP string, osType, scanScene int) error {
	row := &mysqls.BaselineCheckResult{
		TaskID:      taskID,
		TargetID:    targetID,
		TargetIP:    targetIP,
		OSType:      osType,
		ScanScene:   scanScene,
		RuleID:      0,
		RuleName:    "任务执行中",
		CheckResult: BaselineCheckStatusPending,
		CreateTime:  time.Now(),
	}

	fmt.Printf(">>> PersistBaselineCheckPlaceholder: inserting row taskID=%d targetIP=%s scanScene=%d checkResult=%d ruleName=%s\n", taskID, targetIP, scanScene, BaselineCheckStatusPending, "任务执行中")

	result := mysql.FromContext(ctx).Model(row).Create(row)
	if result.Error != nil {
		log.Errorf("PersistBaselineCheckPlaceholder: DB insert error: %v", result.Error)
		return result.Error
	}
	fmt.Printf(">>> PersistBaselineCheckPlaceholder: DB insert OK, rowsAffected=%d id=%d\n", result.RowsAffected, row.ID)
	return nil
}
