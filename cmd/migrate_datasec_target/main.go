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

const createSQL = `CREATE TABLE IF NOT EXISTS datasec_db_target (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  user_id int(11) NOT NULL DEFAULT '0' COMMENT '所属用户',
  name varchar(128) NOT NULL DEFAULT '' COMMENT '目标名称',
  group_name varchar(64) NOT NULL DEFAULT '' COMMENT '分组',
  db_type int(11) NOT NULL DEFAULT '1' COMMENT '1=MySQL 2=PG 3=Mongo 4=Redis 5=CouchDB',
  db_host varchar(255) NOT NULL DEFAULT '' COMMENT '地址',
  db_port int(11) NOT NULL DEFAULT '0' COMMENT '端口',
  db_name varchar(128) NOT NULL DEFAULT '' COMMENT '库名',
  db_user varchar(128) NOT NULL DEFAULT '' COMMENT '用户名',
  db_password varchar(512) NOT NULL DEFAULT '' COMMENT '密码(AES加密base64)',
  remark varchar(512) NOT NULL DEFAULT '' COMMENT '备注',
  create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_user_db_type (user_id, db_type),
  KEY idx_user_group (user_id, group_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据安全-数据库扫描目标库'`

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
		fmt.Println("找不到 config.json，请手动执行 docs/sql/datasec_db_target.sql")
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
	fmt.Println("datasec_db_target 表创建成功")
}
