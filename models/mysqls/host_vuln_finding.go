package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

// HostVulnFinding 主机 CVE 漏洞发现明细
type HostVulnFinding struct {
	ID             int       `gorm:"column:id;primary_key" json:"id"`
	TaskID         int       `gorm:"column:task_id" json:"taskId"`
	TargetID       int       `gorm:"column:target_id" json:"targetId"`
	TargetIP       string    `gorm:"column:target_ip" json:"targetIp"`
	CveID          string    `gorm:"column:cve_id" json:"cveId"`
	Title          string    `gorm:"column:title" json:"title"`
	Severity       string    `gorm:"column:severity" json:"severity"`
	RiskLevel      int       `gorm:"column:risk_level" json:"riskLevel"`
	PackageName    string    `gorm:"column:package_name" json:"packageName"`
	PackageVersion string    `gorm:"column:package_version" json:"packageVersion"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
}

func (HostVulnFinding) TableName() string {
	return "host_vuln_finding"
}

func (m *HostVulnFinding) BatchAdd(ctx context.Context, list []HostVulnFinding) error {
	if len(list) == 0 {
		return nil
	}
	return mysql.FromContext(ctx).Model(m).CreateInBatches(list, 100).Error
}

func (m *HostVulnFinding) DeleteByTaskTarget(ctx context.Context, taskID int, targetIP string) error {
	q := mysql.FromContext(ctx).Model(m).Where("task_id = ?", taskID)
	if targetIP != "" {
		q = q.Where("target_ip = ?", targetIP)
	}
	return q.Delete(nil).Error
}

func (m *HostVulnFinding) ListByTask(ctx context.Context, taskID int, targetIP string) ([]HostVulnFinding, error) {
	q := mysql.FromContext(ctx).Model(m).Where("task_id = ?", taskID)
	if targetIP != "" {
		q = q.Where("target_ip = ?", targetIP)
	}
	var list []HostVulnFinding
	err := q.Order("risk_level DESC, cve_id ASC").Find(&list).Error
	return list, err
}
