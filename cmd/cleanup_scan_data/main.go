package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

// 与 application/system.go runtimeTablesToClean 保持一致（以实际库表为准）
var runtimeTablesToClean = []string{
	"security_report", "report_record",
	"report_verify_port", "report_verify_target", "report_verify_task", "report_verify_vul",
	"baseline_check_result", "host_vul_scan", "host_vul_finding",
	"host_malware_scan", "malware_check_result",
	"db_check_result", "sensitive_data_result",
	"remote_session",
	"flow_log", "flow_risk", "flow_target", "flow_task", "flow_base",
	"asset_log", "asset_task_result", "asset_vul", "asset_port", "asset_risk_trend",
	"task_log_info", "task_log", "task_result", "task_task_result",
	"task_vul", "task_evidence", "task_target_result", "task_task_info",
	"task_target", "task_configuration", "task_task",
}

func isMissingTableErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "1146") ||
		strings.Contains(msg, "doesn't exist") ||
		strings.Contains(msg, "does not exist") ||
		strings.Contains(msg, "no such table")
}

func clearTable(db *gorm.DB, tableName string) error {
	log.Infof("  清理表: %s ...", tableName)
	truncateSQL := "TRUNCATE TABLE `" + tableName + "`"
	if err := db.Exec(truncateSQL).Error; err != nil {
		if isMissingTableErr(err) {
			log.Warnf("  表 %s 不存在，跳过", tableName)
			return nil
		}
		log.Warnf("  TRUNCATE %s 失败，改用 DELETE: %v", tableName, err)
		if err := db.Exec("DELETE FROM `" + tableName + "`").Error; err != nil {
			if isMissingTableErr(err) {
				log.Warnf("  表 %s 不存在，跳过", tableName)
				return nil
			}
			return fmt.Errorf("清理表 %s 失败: %v", tableName, err)
		}
	}
	return nil
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
	db := mysql.FromContext(ctx)

	log.Info("开始清理扫描运行数据...")
	_ = db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	defer func() {
		_ = db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	}()

	var failed int
	for _, name := range runtimeTablesToClean {
		if err := clearTable(db, name); err != nil {
			log.Error(err)
			failed++
		}
	}

	elapsed := time.Since(startTime)
	if failed > 0 {
		log.Errorf("清理完成，%d 张表失败，耗时: %v", failed, elapsed)
		os.Exit(1)
	}
	log.Infof("清理完成！耗时: %v", elapsed)
	fmt.Println("========================================")
	fmt.Println("  ✅ 扫描运行数据清理完毕")
	fmt.Println("========================================")
}
