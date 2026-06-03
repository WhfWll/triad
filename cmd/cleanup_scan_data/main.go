package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"

	"smart/models/mysqls"
)

var tablesInOrder = []interface{}{
	// 1. 任务日志详情
	&mysqls.Taskloginfo{},
	// 2. 任务日志
	&mysqls.Tasklog{},
	// 3. 任务检测结果
	&mysqls.TaskResult{},
	// 4. 任务通用结果
	&mysqls.TaskTaskResult{},
	// 5. 任务漏洞
	&mysqls.TaskVul{},
	// 6. 任务证据
	&mysqls.TaskEvidence{},
	// 7. 任务检查结果
	&mysqls.TaskCheckResult{},
	// 8. 任务目标结果
	&mysqls.TaskTargetResult{},
	// 9. 任务信息
	&mysqls.TaskTaskInfo{},
	// 10. 任务目标
	&mysqls.TaskTarget{},
	// 11. 任务（主表）
	&mysqls.TaskTask{},
	// 12. 任务组任务关联
	&mysqls.TaskGroupTask{},
}

var hostTables = []interface{}{
	&mysqls.BaselineCheckResult{},
	&mysqls.HostVulnScan{},
	&mysqls.HostVulnFinding{},
	&mysqls.HostMalwareScan{},
	&mysqls.MalwareCheckResult{},
}

var dataTables = []interface{}{
	&mysqls.DBCheckResult{},
	&mysqls.SensitiveDataResult{},
}

var reportTables = []interface{}{
	&mysqls.SecurityReport{},
	&mysqls.Reportrecord{},
}

var sessionTables = []interface{}{
	&mysqls.RemoteSession{},
}

var tableNameMap = map[string]string{
	"Taskloginfo":        "task_log_info",
	"Tasklog":            "task_log",
	"TaskResult":         "task_result",
	"TaskTaskResult":     "task_task_result",
	"TaskVul":            "task_vul",
	"TaskEvidence":       "task_evidence",
	"TaskCheckResult":    "task_check_result",
	"TaskTargetResult":   "task_target_result",
	"TaskTaskInfo":       "task_task_info",
	"TaskTarget":         "task_target",
	"TaskTask":           "task_task",
	"TaskGroupTask":      "task_group_task",
	"BaselineCheckResult": "baseline_check_result",
	"HostVulnScan":       "host_vuln_scan",
	"HostVulnFinding":    "host_vuln_finding",
	"HostMalwareScan":    "host_malware_scan",
	"MalwareCheckResult": "malware_check_result",
	"DBCheckResult":      "db_check_result",
	"SensitiveDataResult":"sensitive_data_result",
	"SecurityReport":     "security_report",
	"Reportrecord":       "report_record",
	"RemoteSession":      "remote_session",
}

func getTableName(m interface{}) string {
	switch t := m.(type) {
	case *mysqls.Taskloginfo:
		return t.TableName()
	case *mysqls.Tasklog:
		return t.TableName()
	case *mysqls.TaskResult:
		return t.TableName()
	case *mysqls.TaskTaskResult:
		return t.TableName()
	case *mysqls.TaskVul:
		return t.TableName()
	case *mysqls.TaskEvidence:
		return t.TableName()
	case *mysqls.TaskCheckResult:
		return t.TableName()
	case *mysqls.TaskTargetResult:
		return t.TableName()
	case *mysqls.TaskTaskInfo:
		return t.TableName()
	case *mysqls.TaskTarget:
		return t.TableName()
	case *mysqls.TaskTask:
		return t.TableName()
	case *mysqls.TaskGroupTask:
		return t.TableName()
	case *mysqls.BaselineCheckResult:
		return t.TableName()
	case *mysqls.HostVulnScan:
		return t.TableName()
	case *mysqls.HostVulnFinding:
		return t.TableName()
	case *mysqls.HostMalwareScan:
		return t.TableName()
	case *mysqls.MalwareCheckResult:
		return t.TableName()
	case *mysqls.DBCheckResult:
		return t.TableName()
	case *mysqls.SensitiveDataResult:
		return t.TableName()
	case *mysqls.SecurityReport:
		return t.TableName()
	case *mysqls.Reportrecord:
		return t.TableName()
	case *mysqls.RemoteSession:
		return t.TableName()
	default:
		return "unknown"
	}
}

func truncateTable(ctx context.Context, model interface{}) error {
	tableName := getTableName(model)
	log.Infof("  清理表: %s ...", tableName)
	db := mysql.FromContext(ctx).Model(model)
	// 禁用外键检查后删除
	_ = mysql.FromContext(ctx).Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	err := db.Delete(nil).Error
	_ = mysql.FromContext(ctx).Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	if err != nil {
		return fmt.Errorf("清理表 %s 失败: %v", tableName, err)
	}
	return nil
}

func truncateTableSession(ctx context.Context, model interface{}) error {
	// DELETE 不使用 Exec, 用 model delete
	return truncateTable(ctx, model)
}

func main() {
	var confirm bool
	flag.BoolVar(&confirm, "confirm", false, "确认执行清理操作")
	flag.Parse()

	if !confirm {
		fmt.Println("========================================")
		fmt.Println("  扫描运行数据清理工具")
		fmt.Println("========================================")
		fmt.Println()
		fmt.Println("⚠️  此操作将删除所有扫描任务产生的运行数据，包括：")
		fmt.Println("   - 所有扫描任务及目标")
		fmt.Println("   - 所有检查结果（基线/CVE/恶意代码/数据库/敏感数据）")
		fmt.Println("   - 所有任务日志和日志详情")
		fmt.Println("   - 所有检测结果、漏洞记录、证据")
		fmt.Println("   - 所有安全报告和报告记录")
		fmt.Println("   - 所有远程会话记录")
		fmt.Println()
		fmt.Println("保留的数据包括：")
		fmt.Println("   - 用户和用户组")
		fmt.Println("   - 漏洞库和规则库")
		fmt.Println("   - 系统配置和字典")
		fmt.Println("   - 任务模板/场景")
		fmt.Println("   - 资产信息和目标库")
		fmt.Println("   - 审计日志")
		fmt.Println()
		fmt.Println("用法: go run cmd/cleanup_scan_data/main.go -confirm")
		fmt.Println()
		os.Exit(0)
	}

	fmt.Println("Loading config...")
	configPath := "../../config.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config.json"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Errorf("config.json not found")
		os.Exit(1)
	}

	err := config.NewConfig(configPath)
	if err != nil {
		log.Errorf("load config failed: %v", err)
		os.Exit(1)
	}
	mysql.Setup()

	ctx := context.Background()
	startTime := time.Now()

	log.Info("开始清理扫描运行数据...")

	// ============================================================
	// 1. 安全报告
	// ============================================================
	log.Info("[1/5] 清理安全报告数据...")
	for _, m := range reportTables {
		if err := truncateTable(ctx, m); err != nil {
			log.Error(err)
		}
	}

	// ============================================================
	// 2. 主机安全结果
	// ============================================================
	log.Info("[2/5] 清理主机安全检查结果...")
	for _, m := range hostTables {
		if err := truncateTable(ctx, m); err != nil {
			log.Error(err)
		}
	}

	// ============================================================
	// 3. 数据安全结果
	// ============================================================
	log.Info("[3/5] 清理数据安全检查结果...")
	for _, m := range dataTables {
		if err := truncateTable(ctx, m); err != nil {
			log.Error(err)
		}
	}

	// ============================================================
	// 4. 远程会话
	// ============================================================
	log.Info("[4/5] 清理远程会话记录...")
	for _, m := range sessionTables {
		if err := truncateTable(ctx, m); err != nil {
			log.Error(err)
		}
	}

	// ============================================================
	// 5. 扫描任务（按外键依赖顺序）
	// ============================================================
	log.Info("[5/5] 清理扫描任务数据...")
	for _, m := range tablesInOrder {
		if err := truncateTable(ctx, m); err != nil {
			log.Error(err)
		}
	}

	elapsed := time.Since(startTime)
	log.Infof("清理完成！耗时: %v", elapsed)

	// 验证清理结果
	log.Info("验证清理结果...")
	verifyCleanup(ctx)

	fmt.Println("========================================")
	fmt.Println("  ✅ 扫描运行数据清理完毕")
	fmt.Println("========================================")
}

func verifyCleanup(ctx context.Context) {
	tables := []struct {
		name  string
		model interface{}
	}{
		{"task_task", &mysqls.TaskTask{}},
		{"task_target", &mysqls.TaskTarget{}},
		{"task_result", &mysqls.TaskResult{}},
		{"task_vul", &mysqls.TaskVul{}},
		{"task_log", &mysqls.Tasklog{}},
		{"baseline_check_result", &mysqls.BaselineCheckResult{}},
		{"host_vuln_scan", &mysqls.HostVulnScan{}},
		{"host_malware_scan", &mysqls.HostMalwareScan{}},
		{"db_check_result", &mysqls.DBCheckResult{}},
		{"sensitive_data_result", &mysqls.SensitiveDataResult{}},
		{"security_report", &mysqls.SecurityReport{}},
		{"report_record", &mysqls.Reportrecord{}},
	}

	for _, t := range tables {
		var count int64
		mysql.FromContext(ctx).Model(t.model).Count(&count)
		if count > 0 {
			log.Warnf("  ⚠️ 表 %s 仍有 %d 条记录", t.name, count)
		} else {
			log.Infof("  ✅ 表 %s 已清空", t.name)
		}
	}
}
