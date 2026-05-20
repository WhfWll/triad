package services

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

const createDatasecRuleTableSQL = `CREATE TABLE IF NOT EXISTS datasec_rule (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  rule_code int(11) NOT NULL COMMENT '规则编号',
  name varchar(255) NOT NULL DEFAULT '' COMMENT '规则名称',
  description varchar(512) NOT NULL DEFAULT '' COMMENT '规则描述',
  category int(11) NOT NULL DEFAULT '0' COMMENT '分类',
  risk int(11) NOT NULL DEFAULT '0' COMMENT '风险等级',
  db_type int(11) NOT NULL DEFAULT '0' COMMENT '适用数据库类型',
  queries_json text COMMENT '检查 SQL JSON 数组',
  expected_value varchar(512) NOT NULL DEFAULT '' COMMENT '期望值',
  match_type varchar(32) NOT NULL DEFAULT 'contains' COMMENT '匹配方式',
  fix_suggestion text COMMENT '修复建议',
  risk_description varchar(512) NOT NULL DEFAULT '' COMMENT '风险说明',
  enabled tinyint(4) NOT NULL DEFAULT '1' COMMENT '1=启用',
  create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_rule_code (rule_code),
  KEY idx_db_type_enabled (db_type, enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据安全检测规则表'`

const alterDatasecRuleUTF8 = `ALTER TABLE datasec_rule CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci`

// EnsureDatasecRuleTable 确保 datasec_rule 表存在（开发环境未跑全量 SQL 时自动建表）
func EnsureDatasecRuleTable(ctx context.Context) {
	db := mysql.FromContext(ctx)
	if err := db.Exec(createDatasecRuleTableSQL).Error; err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "already exists") {
			log.Warnf("EnsureDatasecRuleTable: %v", err)
		}
	}
	_ = db.Exec(alterDatasecRuleUTF8).Error
}
