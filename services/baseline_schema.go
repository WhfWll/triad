package services

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

const alterBaselineCheckResultColumnsSQL = `ALTER TABLE baseline_check_result
  MODIFY COLUMN expected_value text COMMENT '期望值',
  MODIFY COLUMN actual_value text COMMENT '实际值'`

// EnsureBaselineCheckResultSchema 启动时扩展核查结果字段，避免命令输出过长写入失败。
func EnsureBaselineCheckResultSchema(ctx context.Context) {
	db := mysql.GetDB()
	if db == nil {
		return
	}
	if err := db.Exec(alterBaselineCheckResultColumnsSQL).Error; err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "doesn't exist") || strings.Contains(msg, "does not exist") {
			log.Infof("baseline_check_result table not found, skip column alter")
			return
		}
		log.Warnf("EnsureBaselineCheckResultSchema: %v", err)
	}
}
