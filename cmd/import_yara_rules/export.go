//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"smart/models/mysqls"

	"gitlabee.4dogs.cn/common/config"
	mysqlEnv "gitlabee.4dogs.cn/common/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type exportItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RiskLevel   int    `json:"riskLevel"`
	RuleContent string `json:"ruleContent"`
	OsType      int    `json:"osType"`
	Category    string `json:"category"`
	Status      int    `json:"status"`
}

func main() {
	outputPath := "docs/yara_rules.json"
	if len(os.Args) > 1 {
		outputPath = os.Args[1]
	}

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
		fmt.Println("找不到配置文件")
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

	var rules []mysqls.MalwareRule
	gormDB.WithContext(ctx).Model(&mysqls.MalwareRule{}).Order("id").Find(&rules)

	fmt.Fprintf(os.Stderr, "读取到 %d 条规则\n", len(rules))

	items := make([]exportItem, 0, len(rules))
	for _, r := range rules {
		items = append(items, exportItem{
			Name:        r.Name,
			Description: r.Description,
			RiskLevel:   r.RiskLevel,
			RuleContent: r.RuleContent,
			OsType:      r.OsType,
			Category:    r.Category,
			Status:      r.Status,
		})
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON 序列化失败: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "导出完成！输出文件: %s (%.2f MB, %d 条规则)\n", outputPath, float64(len(data))/1024/1024, len(items))
}
