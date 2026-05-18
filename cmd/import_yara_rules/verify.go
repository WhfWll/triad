//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"gitlabee.4dogs.cn/common/config"
	mysqlEnv "gitlabee.4dogs.cn/common/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
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
		os.Exit(1)
	}

	config.NewConfig(configFile)

	dbSettings := make([]mysqlEnv.DatabaseSetting, 1)
	config.Load("mysql", &dbSettings)

	var dsn string
	for _, v := range dbSettings {
		if v.Name == "default" {
			dsn = v.Master
			break
		}
	}
	if dsn == "" && len(dbSettings) > 0 {
		dsn = dbSettings[0].Master
	}

	gormDB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	sqlDB, _ := gormDB.DB()
	defer sqlDB.Close()

	ctx := context.Background()

	var total int64
	gormDB.WithContext(ctx).Table("malware_rule").Count(&total)

	var byCategory []struct {
		Category string
		Count    int64
	}
	gormDB.WithContext(ctx).Table("malware_rule").Select("category, COUNT(*) as count").Group("category").Order("count desc").Scan(&byCategory)

	var byRisk []struct {
		RiskLevel int
		Count     int64
	}
	gormDB.WithContext(ctx).Table("malware_rule").Select("risk_level, COUNT(*) as count").Group("risk_level").Order("risk_level").Scan(&byRisk)

	fmt.Println("===========================================")
	fmt.Println("  YARA 病毒库规则导入验证")
	fmt.Println("===========================================")
	fmt.Printf("  总计: %d 条规则\n\n", total)

	fmt.Println("  按分类统计:")
	for _, c := range byCategory {
		fmt.Printf("    %-20s %d\n", c.Category, c.Count)
	}

	fmt.Println("\n  按风险统计:")
	for _, r := range byRisk {
		level := map[int]string{1: "高危", 2: "中危", 3: "低危"}[r.RiskLevel]
		fmt.Printf("    %-10s %d\n", level, r.Count)
	}
	fmt.Println("===========================================")
}