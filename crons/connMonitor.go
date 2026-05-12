package crons

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"
	"gorm.io/gorm"
	"smart/models/redises"
)

// MonitorDBConnection 监控数据库连接状态
func MonitorDBConnection() {
	ctx := context.Background()

	// 检查Redis连接
	var redisCommon redises.RedisCommon
	_, err := redisCommon.Ping(ctx)
	if err != nil {
		log.Errorf("Redis连接异常: %v", err)
		// 尝试重新连接Redis
		redis.Setup()
	}

	// 检查MySQL连接
	if db := GetDB(); db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Errorf("获取MySQL原生连接失败: %v", err)
		} else {
			// 检查连接是否有效
			err = sqlDB.Ping()
			if err != nil {
				log.Errorf("MySQL连接异常: %v", err)
				// 尝试重新连接
				if err := reconnectMySQL(sqlDB); err != nil {
					log.Errorf("MySQL重连失败: %v", err)
				} else {
					log.Info("MySQL重连成功")
				}
			}
		}
	}

}

// reconnectMySQL 重新连接MySQL
func reconnectMySQL(sqlDB *sql.DB) error {
	// 关闭现有连接
	if err := sqlDB.Close(); err != nil {
		log.Warnf("关闭旧的MySQL连接时出错: %v", err)
	}
	mysql.Setup()
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return mysql.GetDB()
}
