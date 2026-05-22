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

const createSQL = `CREATE TABLE IF NOT EXISTS sensitive_data_result (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  task_id int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  target_id int(11) NOT NULL DEFAULT '0' COMMENT '目标id',
  target_ip varchar(50) NOT NULL DEFAULT '' COMMENT '目标IP',
  db_type int(11) NOT NULL DEFAULT '0' COMMENT '数据库类型',
  db_name varchar(128) NOT NULL DEFAULT '' COMMENT '数据库名',
  table_name varchar(128) NOT NULL DEFAULT '' COMMENT '表名',
  column_name varchar(128) NOT NULL DEFAULT '' COMMENT '字段名',
  data_type int(11) NOT NULL DEFAULT '0' COMMENT '数据类型',
  data_level int(11) NOT NULL DEFAULT '0' COMMENT '数据等级',
  sample_data varchar(512) NOT NULL DEFAULT '' COMMENT '样本数据',
  match_rule varchar(512) NOT NULL DEFAULT '' COMMENT '匹配规则',
  match_type int(11) NOT NULL DEFAULT '0' COMMENT '匹配类型',
  total_rows bigint(20) NOT NULL DEFAULT '0' COMMENT '总行数',
  create_time datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (id),
  KEY idx_task_id (task_id),
  KEY idx_target_id (target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='敏感数据发现结果'`

func main() {
	configPaths := []string{"config.json", "../config.json", "D:\\goproject\\triad\\config.json"}
	var configFile string
	for _, p := range configPaths {
		if _, err := os.Stat(p); err == nil {
			configFile = p
			break
		}
	}
	if configFile == "" {
		fmt.Println("找不到 config.json，请手动执行 docs/sql/sensitive_data_result.sql")
		os.Exit(1)
	}
	if err := config.NewConfig(configFile); err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}
	var dbSettings []mysqlEnv.DatabaseSetting
	if err := config.Load("mysql", &dbSettings); err != nil {
		fmt.Printf("读取 MySQL 配置失败: %v\n", err)
		os.Exit(1)
	}
	dsn := ""
	for _, v := range dbSettings {
		if v.Name == "default" {
			dsn = v.Master
			break
		}
	}
	if dsn == "" && len(dbSettings) > 0 {
		dsn = dbSettings[0].Master
	}
	if dsn == "" {
		fmt.Println("未找到 MySQL 连接信息")
		os.Exit(1)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		os.Exit(1)
	}
	if err := db.WithContext(context.Background()).Exec(createSQL).Error; err != nil {
		fmt.Printf("建表失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("sensitive_data_result 表创建成功")
}
