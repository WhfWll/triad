package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

const (
	HostVulnScanStatusOK    = 1
	HostVulnScanStatusError = 2
)

// HostVulnScan 主机 CVE 扫描任务汇总（一行 = 一次任务 × 一个目标）
type HostVulnScan struct {
	ID             int       `gorm:"column:id;primary_key" json:"id"`
	TaskID         int       `gorm:"column:task_id" json:"taskId"`
	TargetID       int       `gorm:"column:target_id" json:"targetId"`
	TargetIP       string    `gorm:"column:target_ip" json:"targetIp"`
	OSType         int       `gorm:"column:os_type" json:"osType"`
	Packages       int       `gorm:"column:packages" json:"packages"`
	MatchedVulns   int       `gorm:"column:matched_vulns" json:"matchedVulns"`
	Critical       int       `gorm:"column:critical" json:"critical"`
	High           int       `gorm:"column:high" json:"high"`
	Medium         int       `gorm:"column:medium" json:"medium"`
	Low            int       `gorm:"column:low" json:"low"`
	WorstRiskLevel int       `gorm:"column:worst_risk_level" json:"worstRiskLevel"`
	ScanStatus     int       `gorm:"column:scan_status" json:"scanStatus"`
	ErrorMessage   string    `gorm:"column:error_message" json:"errorMessage"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
}

func (HostVulnScan) TableName() string {
	return "host_vuln_scan"
}

func (m *HostVulnScan) Save(ctx context.Context, row *HostVulnScan) error {
	if err := m.DeleteByTaskTarget(ctx, row.TaskID, row.TargetIP); err != nil {
		return err
	}
	row.CreateTime = time.Now()
	return mysql.FromContext(ctx).Model(m).Create(row).Error
}

func (m *HostVulnScan) DeleteByTaskTarget(ctx context.Context, taskID int, targetIP string) error {
	q := mysql.FromContext(ctx).Model(m).Where("task_id = ?", taskID)
	if targetIP != "" {
		q = q.Where("target_ip = ?", targetIP)
	}
	return q.Delete(nil).Error
}

func (m *HostVulnScan) CountAll(ctx context.Context) (int64, error) {
	var n int64
	err := mysql.FromContext(ctx).Model(m).Count(&n).Error
	return n, err
}

func (m *HostVulnScan) ListPaged(ctx context.Context, page, size int) ([]HostVulnScan, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size
	var list []HostVulnScan
	err := mysql.FromContext(ctx).Model(m).Order("create_time DESC").Offset(offset).Limit(size).Find(&list).Error
	return list, err
}

func (m *HostVulnScan) ListByTaskID(ctx context.Context, taskID int) ([]HostVulnScan, error) {
	var list []HostVulnScan
	err := mysql.FromContext(ctx).Model(m).Where("task_id = ?", taskID).Order("target_ip ASC").Find(&list).Error
	return list, err
}
