package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

type Reportverifytarget struct {
	ID           int       `gorm:"column:id;primary_key" json:"id"`           // 主键
	TaskId       int       `gorm:"column:task_id" json:"task_id"`             // 任务id
	Target       string    `gorm:"column:target" json:"target"`               // 目标
	Os           string    `gorm:"column:os" json:"os"`                       // 操作系统
	Risk         int       `gorm:"column:risk" json:"risk"`                   // 风险等级
	IsAlive      int       `gorm:"column:is_alive" json:"is_alive"`           // 是否存活
	Exp          int       `gorm:"column:exp" json:"exp"`                     // 利用成功数
	Verify       int       `gorm:"column:verify" json:"verify"`               // 验证成功数
	Failed       int       `gorm:"column:failed" json:"failed"`               // 验证失败数
	Unverify     int       `gorm:"column:unverify" json:"unverify"`           // 未能验证数
	Status       int       `gorm:"column:status" json:"status"`               // 运行状态
	AnalysisData string    `gorm:"column:analysis_data" json:"analysis_data"` // 中间分析数据
	CreateTime   time.Time `gorm:"column:create_time" json:"create_time"`     // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`     // 更新时间
}

// TableName sets insert table name for this struct type
func (r *Reportverifytarget) TableName() string {
	return "report_verify_target"
}

// Get retrieves a list of reportverifytarget from database
func (r *Reportverifytarget) GetReportverifytargetList(ctx context.Context, taskId, risk, page, limit int, search string) ([]Reportverifytarget, int64, error) {
	var (
		reportverifytargetList []Reportverifytarget
		count                  int64
		db                     = mysql.FromContext(ctx).Model(&Reportverifytarget{})
		query                  string
		args                   []interface{}
	)

	query += "task_id = ?"
	args = append(args, taskId)
	if risk != 0 {
		query += " and risk = ?"
		args = append(args, risk)
	}
	if search != "" {
		query += " and target like ?"
		args = append(args, "%"+search+"%")
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Find(&reportverifytargetList)

	return reportverifytargetList, count, nil
}

// Get retrieves a single record of reportverifytarget from database
func (r *Reportverifytarget) GetReportverifytarget(ctx context.Context) (Reportverifytarget, error) {
	var (
		reportverifytarget Reportverifytarget
		err                error
		db                 = mysql.FromContext(ctx).Model(&Reportverifytarget{})
	)

	curErr := db.Where("id = ?", r.ID).First(&reportverifytarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return reportverifytarget, err
}

// Add persists reportverifytarget to database
func (r *Reportverifytarget) AddReportverifytarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytarget{})

	if err := db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

// Update changes reportverifytarget by id
func (r *Reportverifytarget) UpdateReportverifytarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytarget{})

	if err := db.Where("id = ?", r.ID).Updates(r).Error; err != nil {
		return err
	}

	return nil
}

// Delete reportverifytarget by id
func (r *Reportverifytarget) DeleteReportverifytarget(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytarget{})

	if err := db.Where("id = ?", r.ID).Delete(r).Error; err != nil {
		return err
	}

	return nil
}

func (r *Reportverifytarget) GetOneWaitTarget(ctx context.Context) (Reportverifytarget, error) {
	var (
		reportverifytarget Reportverifytarget
		err                error
		db                 = mysql.FromContext(ctx).Model(&Reportverifytarget{})
	)
	curErr := db.Where("status = ?", enums.TargetStatusToBegin).First(&reportverifytarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return reportverifytarget, err
}

// 更新目标状态
func (r *Reportverifytarget) UpdateReportverifytargetStatus(ctx context.Context, id, status int) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytarget{})
	if err := db.Where("id = ?", id).Updates(Reportverifytarget{Status: status}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Reportverifytarget) GetTargetsByStatus(ctx context.Context, taskId int, status []int) ([]Reportverifytarget, error) {
	var (
		reportverifytarget []Reportverifytarget
		err                error
		db                 = mysql.FromContext(ctx).Model(&Reportverifytarget{})
		query              string
		args               []interface{}
	)
	query += "status in ?"
	args = append(args, status)
	if taskId != 0 {
		query += " and task_id = ?"
		args = append(args, taskId)
	}
	curErr := db.Where(query, args...).Find(&reportverifytarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return reportverifytarget, err
}

// Update changes reportverifytarget by id
func (r *Reportverifytarget) UpdateReportverifytargetRisk(ctx context.Context, taskId int, target string, risk, unVerify, verify, failed, exp int) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytarget{})
	if err := db.Where("task_id = ? and target = ?", taskId, target).Updates(Reportverifytarget{
		Risk:     risk,
		Unverify: unVerify,
		Verify:   verify,
		Failed:   failed,
		Exp:      exp,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Reportverifytarget) GetAllTargets(ctx context.Context, taskId int) ([]Reportverifytarget, error) {
	var (
		reportverifytarget []Reportverifytarget
		err                error
		db                 = mysql.FromContext(ctx).Model(&Reportverifytarget{})
	)
	curErr := db.Where("task_id = ?", taskId).Find(&reportverifytarget).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return reportverifytarget, err
}

// 更新目标状态
func (r *Reportverifytarget) UpdateReportverifytargetStatusByTaskId(ctx context.Context, taskId, status int) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytarget{})
	if err := db.Where("task_id = ?", taskId).Updates(Reportverifytarget{Status: status}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteReportverifytargetByTaskId 删除报告验证目标通过任务id
func (r *Reportverifytarget) DeleteReportverifytargetByTaskId(ctx context.Context, taskId int) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytarget{})
	if err := db.Where("task_id = ?", taskId).Delete(r).Error; err != nil {
		return err
	}
	return nil
}
