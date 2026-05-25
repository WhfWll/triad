package services

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

const createHostVulnScanTableSQL = `CREATE TABLE IF NOT EXISTS host_vuln_scan (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  task_id int(11) NOT NULL DEFAULT '0' COMMENT '任务批次',
  target_id int(11) NOT NULL DEFAULT '0' COMMENT '目标序号',
  target_ip varchar(64) NOT NULL DEFAULT '' COMMENT '目标 IP',
  os_type int(11) NOT NULL DEFAULT '0' COMMENT '操作系统类型',
  packages int(11) NOT NULL DEFAULT '0' COMMENT '扫描软件包数',
  matched_vulns int(11) NOT NULL DEFAULT '0' COMMENT '匹配 CVE 数',
  critical int(11) NOT NULL DEFAULT '0' COMMENT '严重',
  high int(11) NOT NULL DEFAULT '0' COMMENT '高危',
  medium int(11) NOT NULL DEFAULT '0' COMMENT '中危',
  low int(11) NOT NULL DEFAULT '0' COMMENT '低危',
  worst_risk_level int(11) NOT NULL DEFAULT '0' COMMENT '最高风险等级',
  scan_status int(11) NOT NULL DEFAULT '1' COMMENT '1=完成 2=异常',
  error_message text COMMENT '扫描失败原因',
  create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '扫描时间',
  PRIMARY KEY (id),
  KEY idx_task_target (task_id, target_ip),
  KEY idx_create_time (create_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机 CVE 漏洞扫描任务汇总'`

const createHostVulnFindingTableSQL = `CREATE TABLE IF NOT EXISTS host_vuln_finding (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  task_id int(11) NOT NULL DEFAULT '0' COMMENT '任务批次',
  target_id int(11) NOT NULL DEFAULT '0' COMMENT '目标序号',
  target_ip varchar(64) NOT NULL DEFAULT '' COMMENT '目标 IP',
  cve_id varchar(32) NOT NULL DEFAULT '' COMMENT 'CVE 编号',
  title varchar(512) NOT NULL DEFAULT '' COMMENT '漏洞标题',
  severity varchar(32) NOT NULL DEFAULT '' COMMENT '严重程度文本',
  risk_level int(11) NOT NULL DEFAULT '0' COMMENT '风险等级',
  package_name varchar(255) NOT NULL DEFAULT '' COMMENT '影响软件包',
  package_version varchar(128) NOT NULL DEFAULT '' COMMENT '软件包版本',
  create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发现时间',
  PRIMARY KEY (id),
  KEY idx_task_target (task_id, target_ip),
  KEY idx_cve (cve_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机 CVE 漏洞发现明细'`

// EnsureHostVulnSchema 启动时确保 CVE 扫描结果表存在。
func EnsureHostVulnSchema(ctx context.Context) {
	db := mysql.GetDB()
	if db == nil {
		return
	}
	for _, sql := range []string{createHostVulnScanTableSQL, createHostVulnFindingTableSQL} {
		if err := db.Exec(sql).Error; err != nil {
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "denied") || strings.Contains(msg, "syntax") {
				log.Warnf("EnsureHostVulnSchema: %v", err)
			}
		}
	}
}
