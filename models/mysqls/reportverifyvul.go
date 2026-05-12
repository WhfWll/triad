package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
)

type Reportverifyvul struct {
	ID       int    `gorm:"column:id;primary_key" json:"id"`   // 主键
	TaskId   int    `gorm:"column:task_id" json:"task_id"`     // 任务id
	TargetId int    `gorm:"column:target_id" json:"target_id"` // 目标id
	Name     string `gorm:"column:name" json:"name"`           // 漏洞名称
	Risk     int    `gorm:"column:risk" json:"risk"`           // 风险等级
	Status   int    `gorm:"column:status" json:"status"`       // 漏洞状态
	Location string `gorm:"column:location" json:"location"`   // 漏洞位置
	Cve      string `gorm:"column:cve" json:"cve"`             // cve
	Cnvd     string `gorm:"column:Cnvd" json:"Cnvd"`           // cnvd
	Cnnvd    string `gorm:"column:cnnvd" json:"cnnvd"`         // cnnvd
	Desc     string `gorm:"column:desc" json:"desc"`           // 漏洞描述
	Fix      string `gorm:"column:fix" json:"fix"`             // 修复建议
	Cvss     string `gorm:"column:cvss" json:"cvss"`           // cvss评分
}

// TableName sets insert table name for this struct type
func (r *Reportverifyvul) TableName() string {
	return "report_verify_vul"
}

// Get retrieves a list of reportverifyvul from database
func (r *Reportverifyvul) GetReportverifyvulList(ctx context.Context, taskId, page, limit, risk, status int, search string) ([]Reportverifyvul, int64, error) {
	var (
		reportverifyvulList []Reportverifyvul
		count               int64
		db                  = mysql.FromContext(ctx).Model(&Reportverifyvul{})
		query               string
		args                []interface{}
	)
	query += "task_id = ?"
	args = append(args, taskId)
	if risk != 0 {
		query += " and risk = ?"
		args = append(args, risk)
	}
	if status != 0 {
		query += " and status = ?"
		args = append(args, status)
	}
	if search != "" {
		query += " and (name LIKE ? or location like ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Find(&reportverifyvulList)
	return reportverifyvulList, count, nil
}

// Get retrieves a single record of reportverifyvul from database
func (r *Reportverifyvul) GetReportverifyvul(ctx context.Context) (Reportverifyvul, error) {
	var (
		reportverifyvul Reportverifyvul
		err             error
		db              = mysql.FromContext(ctx).Model(&Reportverifyvul{})
	)

	curErr := db.Where("id = ?", r.ID).First(&reportverifyvul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return reportverifyvul, err
}

// Add persists reportverifyvul to database
func (r *Reportverifyvul) AddReportverifyvul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifyvul{})

	if err := db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

// Update changes reportverifyvul by id
func (r *Reportverifyvul) UpdateReportverifyvul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifyvul{})

	if err := db.Where("id = ?", r.ID).Updates(r).Error; err != nil {
		return err
	}

	return nil
}

// Delete reportverifyvul by id
func (r *Reportverifyvul) DeleteReportverifyvul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifyvul{})

	if err := db.Where("id = ?", r.ID).Delete(r).Error; err != nil {
		return err
	}

	return nil
}

// 获取可以在小智上进行验证的脚本
func (r *Reportverifyvul) GetReportVerifyVulListCanCall(ctx context.Context, target string) ([]Reportverifyvul, error) {
	var (
		reportverifyvulList []Reportverifyvul
		db                  = mysql.FromContext(ctx).Model(&Reportverifyvul{})
	)
	db.Where(`location = ? and status = ? and cve != ""`, target, enums.ReportVerifyStatusUnVerify).Find(&reportverifyvulList)
	return reportverifyvulList, nil
}

// UpdateReportverifyvulStatus
func (r *Reportverifyvul) UpdateReportverifyvulStatus(ctx context.Context, taskId int, target, cve string, status int) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifyvul{})
	if err := db.Where("task_id = ? and location =? and cve = ?", taskId, target, cve).Updates(Reportverifyvul{Status: status}).Error; err != nil {
		return err
	}
	return nil
}

// 获取可以在小智上进行验证的脚本
func (r *Reportverifyvul) GetReportVerifyVulListByStatus(ctx context.Context, taskId int) ([]Reportverifyvul, error) {
	var (
		reportverifyvulList []Reportverifyvul
		db                  = mysql.FromContext(ctx).Model(&Reportverifyvul{})
	)
	db.Where(`task_id = ?`, taskId).Find(&reportverifyvulList)
	return reportverifyvulList, nil
}
