package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Reportverifytask struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`        // 主键
	Name        string    `gorm:"column:name" json:"name"`                // 任务名称
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`   // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`   // 更新时间
	ExecuteType int       `gorm:"column:execute_type" json:"executeType"` // 执行方式
	Producer    int       `gorm:"column:producer" json:"producer"`        // 厂商
	User        int       `gorm:"column:user" json:"user"`                // 提交者
	Status      int       `gorm:"column:status" json:"status"`            // 状态
	Risk        int       `gorm:"column:risk" json:"risk"`                // 风险等级
	Overview    string    `gorm:"column:overview" json:"overview"`        // 统计信息
	IsStats     int       `gorm:"column:is_stats" json:"isStats"`         // 是否统计
	Fileinfo    string    `gorm:"column:fileinfo" json:"fileinfo"`        // 文件信息
	Exp         int       `gorm:"column:exp" json:"exp"`                  // 利用成功
	Verify      int       `gorm:"column:verify" json:"verify"`            // 验证成功
	Failed      int       `gorm:"column:failed" json:"failed"`            // 验证失败
	Unverify    int       `gorm:"column:unverify" json:"unverify"`        // 未能验证
	ExecuteTime time.Time `gorm:"column:execute_time" json:"executeTime"` // 执行时间
}

// TableName sets insert table name for this struct type
func (r *Reportverifytask) TableName() string {
	return "report_verify_task"
}

// Get retrieves a list of reportverifytask from database
func (r *Reportverifytask) GetReportverifytaskList(ctx context.Context, page, limit, risk, producer int, startTime, endTime, search string, userIdList []int) ([]Reportverifytask, int64, error) {
	var (
		reportverifytaskList []Reportverifytask
		count                int64
		db                   = mysql.FromContext(ctx).Model(&Reportverifytask{})
		query                string
		args                 []interface{}
	)

	query += "name like ?"
	args = append(args, "%"+search+"%")
	if risk != 0 {
		query += " and risk = ?"
		args = append(args, risk)
	}
	if producer != 0 {
		query += " and producer = ?"
		args = append(args, producer)
	}
	if startTime != "" {
		query += " and update_time > ?"
		args = append(args, startTime)
	}
	if endTime != "" {
		query += " and update_time < ?"
		args = append(args, endTime)
	}
	if len(userIdList) != 0 {
		query += " and user in ?"
		args = append(args, userIdList)
	}
	if query != "" {
		db.Where(query, args...)
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&reportverifytaskList)

	return reportverifytaskList, count, nil
}

// Get retrieves a single record of reportverifytask from database
func (r *Reportverifytask) GetReportverifytask(ctx context.Context) (Reportverifytask, error) {
	var (
		reportverifytask Reportverifytask
		err              error
		db               = mysql.FromContext(ctx).Model(&Reportverifytask{})
	)

	curErr := db.Where("id = ?", r.ID).First(&reportverifytask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return reportverifytask, err
}

// Add persists reportverifytask to database
func (r *Reportverifytask) AddReportverifytask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytask{})

	if err := db.Create(r).Error; err != nil {
		return err
	}
	return nil
}

// Update changes reportverifytask by id
func (r *Reportverifytask) UpdateReportverifytask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytask{})

	if err := db.Where("id = ?", r.ID).Updates(r).Error; err != nil {
		return err
	}

	return nil
}

// Delete reportverifytask by id
func (r *Reportverifytask) DeleteReportverifytask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytask{})

	r.UpdateTime = time.Now()
	if err := db.Where("id = ?", r.ID).Delete(r).Error; err != nil {
		return err
	}

	return nil
}

// Update changes reportverifytask by id
func (r *Reportverifytask) UpdateReportverifytaskStatus(ctx context.Context, id, status int) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytask{})

	if err := db.Where("id = ?", id).Updates(Reportverifytask{Status: status}).Error; err != nil {
		return err
	}

	return nil
}

// Update changes reportverifytask by id
func (r *Reportverifytask) UpdateReportverifytaskRisk(ctx context.Context, id, risk, unVerify, verify, failed, exp int) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytask{})
	if err := db.Where("id = ?", id).Updates(Reportverifytask{
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

// Update changes reportverifytask by id
func (r *Reportverifytask) UpdateReportverifytaskOverView(ctx context.Context, id int, overview string) error {
	var db = mysql.FromContext(ctx).Model(&Reportverifytask{})
	if err := db.Where("id = ?", id).Updates(Reportverifytask{
		Overview: overview,
	}).Error; err != nil {
		return err
	}
	return nil
}
