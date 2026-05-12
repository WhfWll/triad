package data

import (
	"context"
	"errors"
	"fmt"
	"os"
	"smart/models/mysqls"
	"strings"
	"sync"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"gitlabee.4dogs.cn/common/config"
	mysqlsEnv "gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitConfigInfo() {
	var configFile string
	if _, err := os.Stat("../config.json"); err == nil {
		configFile = "../config.json"
	} else if _, err = os.Stat("/opt/laozhi/smart/config.json"); err == nil {
		configFile = "/opt/laozhi/smart/config.json"
	} else if _, err = os.Stat("D:\\goproject\\smart\\config.json"); err == nil {
		configFile = "D:\\goproject\\smart\\config.json"
	} else if _, err = os.Stat("D:\\goprojects\\smart\\config.json"); err == nil {
		configFile = "D:\\goprojects\\smart\\config.json"
	} else {
		fmt.Println("找不到配置文件")
		return
	}
	fmt.Println("加载配置文件: ", configFile)
	err := config.NewConfig(configFile) //加载配置文件
	if err != nil {
		fmt.Println(err)
	}
}

// 检测mysql是否正常
func CheckMysql(ctx context.Context) error {
	dbSettings := make([]mysqlsEnv.DatabaseSetting, 1)
	err := config.Load("mysql", &dbSettings) //加载数据库配置
	if err != nil {
		return errors.New("mysql 连接异常: " + err.Error())
	}
	var defaultDBSetting *mysqlsEnv.DatabaseSetting
	for _, v := range dbSettings {
		if v.Name == "default" {
			defaultDBSetting = &v
			break
		}
	}
	if defaultDBSetting == nil {
		defaultDBSetting = &dbSettings[0]
		dbSettings = dbSettings[1:]
	}
	gormDB, err := gorm.Open(mysql.Open(defaultDBSetting.Master), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return errors.New("mysql 连接异常: " + err.Error())
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return errors.New("mysql 数据库连接异常: " + err.Error())
	}
	err = sqlDB.Close()
	if err != nil {
		return errors.New("mysql 数据库关闭异常: " + err.Error())
	}
	fmt.Println("mysql 运行正常")
	return nil
}

// 检测系统环境
func CheckRedis(ctx context.Context) error {
	var mapredisSetting = make(map[string]redis.Settings, 0)
	if err := config.Load("redis", &mapredisSetting); err != nil {
		return errors.New("redis 配置异常: " + err.Error())
	}
	if _, ok := mapredisSetting["default"]; !ok {
		return errors.New("redis 配置异常: ")
	}
	var dbMutex sync.RWMutex
	for _, v := range mapredisSetting {
		if v.Default != "" {
			dbMutex.Lock()
			defer dbMutex.Unlock()
			opt, err := goredis.ParseURL(v.Default)
			if err != nil {
				return errors.New("redis 关闭异常: " + err.Error())
			}
			redisClient := goredis.NewClient(opt)
			err = redisClient.Close()
			if err != nil {
				return errors.New("redis 关闭异常: " + err.Error())
			}
		}
	}
	fmt.Println("redis 运行正常")
	return nil
}

// 检测数据库结构
func CheckDbStruct(ctx context.Context, dbStructFile string) {
	dbSettings := make([]mysqlsEnv.DatabaseSetting, 1)
	err := config.Load("mysql", &dbSettings) //加载数据库配置
	if err != nil {
		fmt.Println(errors.New("mysql 连接异常: " + err.Error()))
	}
	var defaultDBSetting *mysqlsEnv.DatabaseSetting
	for _, v := range dbSettings {
		if v.Name == "default" {
			defaultDBSetting = &v
			break
		}
	}

	if defaultDBSetting == nil {
		defaultDBSetting = &dbSettings[0]
		dbSettings = dbSettings[1:]
	}
	targetDecisionDBUrl := strings.ReplaceAll(defaultDBSetting.Master, "smart", "decision")
	targetSmartDBUrl := defaultDBSetting.Master
	err = RestoreDbStruct(ctx, dbStructFile, targetDecisionDBUrl, targetSmartDBUrl)
	if err != nil {
		fmt.Println(err)
	}
}

// 添加日志
func AddLog(ctx context.Context, error string) error {
	var logAuditModel = mysqls.LogAudit{
		LogType:    1,
		Content:    error,
		Username:   "admin",
		Ip:         "127.0.0.1",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	return logAuditModel.LogAuditAdd(ctx)
}
