//go:build ignore

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"gitlabee.4dogs.cn/common/config"
	mysqlEnv "gitlabee.4dogs.cn/common/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// 寻找配置文件
	configPaths := []string{
		"../config.json",
		"/opt/laozhi/smart/config.json",
		"D:\\goproject\\smart\\config.json",
		"D:\\goprojects\\smart\\config.json",
		"config.json",
	}

	var configFile string
	for _, p := range configPaths {
		if _, err := os.Stat(p); err == nil {
			configFile = p
			break
		}
	}

	if configFile == "" {
		fmt.Println("❌ 找不到配置文件")
		fmt.Println("请手动执行 SQL 文件: docs/sql/import_yara_rules.sql")
		fmt.Println("使用 MySQL 客户端连接后执行: source docs/sql/import_yara_rules.sql")
		os.Exit(1)
	}

	fmt.Println("加载配置文件:", configFile)
	if err := config.NewConfig(configFile); err != nil {
		fmt.Println("❌ 加载配置文件失败:", err)
		os.Exit(1)
	}

	dbSettings := make([]mysqlEnv.DatabaseSetting, 1)
	if err := config.Load("mysql", &dbSettings); err != nil {
		fmt.Println("❌ 读取 MySQL 配置失败:", err)
		os.Exit(1)
	}

	var defaultDSN string
	for _, v := range dbSettings {
		if v.Name == "default" {
			defaultDSN = v.Master
			break
		}
	}
	if defaultDSN == "" && len(dbSettings) > 0 {
		defaultDSN = dbSettings[0].Master
	}

	if defaultDSN == "" {
		fmt.Println("❌ 未找到 MySQL 连接信息")
		os.Exit(1)
	}

	fmt.Println("连接 MySQL...")
	gormDB, err := gorm.Open(mysql.Open(defaultDSN), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		fmt.Println("❌ 连接 MySQL 失败:", err)
		os.Exit(1)
	}

	sqlDB, _ := gormDB.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	defer sqlDB.Close()

	// 先建表
	ctx := context.Background()
	fmt.Println("创建 malware_rule 表...")
	gormDB.WithContext(ctx).Exec(`CREATE TABLE IF NOT EXISTS malware_rule (
		id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
		name varchar(255) NOT NULL DEFAULT '' COMMENT '规则名称',
		description varchar(512) NOT NULL DEFAULT '' COMMENT '规则描述',
		risk_level int(11) NOT NULL DEFAULT '0' COMMENT '风险等级：1-高危 2-中危 3-低危',
		rule_content text COMMENT 'YARA 规则内容',
		os_type int(11) NOT NULL DEFAULT '0' COMMENT '适用系统类型：1-Linux 2-Windows 3-国产 4-嵌入式 0-通用',
		category varchar(255) NOT NULL DEFAULT '' COMMENT '规则分类，如：挖矿木马、Webshell、APT',
		status int(11) NOT NULL DEFAULT '1' COMMENT '1-启用 0-停用',
		create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
		update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
		PRIMARY KEY (id),
		KEY idx_os_type (os_type),
		KEY idx_status (status)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='病毒库规则表'`)

	// 清空表
	fmt.Println("清空 malware_rule 表...")
	gormDB.WithContext(ctx).Exec("TRUNCATE TABLE malware_rule")

	// 读取 SQL 文件
	sqlFile := "docs/sql/import_yara_rules.sql"
	f, err := os.Open(sqlFile)
	if err != nil {
		fmt.Println("❌ 打开 SQL 文件失败:", err)
		os.Exit(1)
	}
	defer f.Close()

	// 按行读取并批量执行
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var batch []string
	batchSize := 100
	total := 0
	start := time.Now()

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		batch = append(batch, line)

		if strings.HasSuffix(line, ";") && len(batch) >= batchSize {
			sql := strings.Join(batch, "\n")
			if err := gormDB.WithContext(ctx).Exec(sql).Error; err != nil {
				// 尝试逐条执行
				for _, singleSQL := range batch {
					if strings.TrimSpace(singleSQL) == "" {
						continue
					}
					if err2 := gormDB.WithContext(ctx).Exec(singleSQL).Error; err2 != nil {
						fmt.Printf("⚠️ 执行失败 (跳过): %v\n  SQL: %s\n", err2, singleSQL[:min(len(singleSQL), 200)])
					} else {
						total++
					}
				}
			} else {
				total += len(batch)
			}
			batch = nil

			if total%2000 == 0 {
				elapsed := time.Since(start)
				fmt.Printf("已导入 %d 条规则，耗时 %v...\n", total, elapsed.Round(time.Second))
			}
		}
	}

	// 处理剩余
	if len(batch) > 0 {
		sql := strings.Join(batch, "\n")
		if err := gormDB.WithContext(ctx).Exec(sql).Error; err != nil {
			for _, singleSQL := range batch {
				if strings.TrimSpace(singleSQL) == "" {
					continue
				}
				if err2 := gormDB.WithContext(ctx).Exec(singleSQL).Error; err2 != nil {
					fmt.Printf("⚠️ 执行失败 (跳过): %v\n", err2)
				} else {
					total++
				}
			}
		} else {
			total += len(batch)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("✅ 完成！共导入 %d 条病毒库规则，耗时 %v\n", total, elapsed.Round(time.Second))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
