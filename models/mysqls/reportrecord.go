package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

type Reportrecord struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Name       string    `gorm:"column:name" json:"name"`              // 报告标题
	Type       int       `gorm:"column:type" json:"type"`              // 报告类型，1-统计报告，2-目标报告
	Status     int       `gorm:"column:status" json:"status"`          // 报告状态，1-生成中，2-已生成
	ConfigJSON string    `gorm:"column:config_json" json:"configJSON"` // 报告配置参数
	Format     int       `gorm:"column:format" json:"format"`          // 报告格式，1-html，2-word,3-pdf,4-file,5-csv
	Content    string    `gorm:"column:content" json:"content"`        // 报告内容
	UserID     int       `gorm:"column:user_id" json:"userID"`         // 提交者id
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

// TableName sets insert table name for this struct type
func (r *Reportrecord) TableName() string {
	return "report_record"
}

// GetReportrecordList 报告清单列表及筛选
func (r *Reportrecord) GetReportrecordList(ctx context.Context, field, search string, page, limit int, userIdList []int) ([]Reportrecord, int64, error) {
	var (
		reportrecordList []Reportrecord
		count            int64
		db               = mysql.FromContext(ctx).Model(&Reportrecord{})
		query            string
		args             []interface{}
	)
	if field == "" {
		field = "*"
	}
	query += "1 = 1"
	if len(search) > 0 {
		query += " and name LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if len(userIdList) > 0 {
		query += " and user_id in ?"
		args = append(args, userIdList)
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Select(field).Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&reportrecordList)
	return reportrecordList, count, nil
}

// GetReportrecord 根据id查询一条报告
func (r *Reportrecord) GetReportrecord(ctx context.Context, id int) (Reportrecord, error) {
	var (
		reportrecord Reportrecord
		err          error
		db           = mysql.FromContext(ctx).Model(&Reportrecord{})
	)
	curErr := db.Where("id = ?", id).First(&reportrecord).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return reportrecord, err
}

// AddReportrecord 生成报告
func (r *Reportrecord) AddReportrecord(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&Reportrecord{})
	if err := db.Create(r).Error; err != nil {
		return 0, err
	}
	return r.ID, nil
}

// DeleteReportrecord 根据id删除报告
// ids必须为[]int/[]string
func (r *Reportrecord) DeleteReportrecord(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&Reportrecord{})
	if err := db.Where("id IN ?", ids).Delete(&Reportrecord{}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateContent 更新报告内容
func (r *Reportrecord) UpdateContent(ctx context.Context, id int, content string) error {
	data := make(map[string]interface{})
	data["content"] = content
	data["status"] = enums.ReportStatusFinish

	var db = mysql.FromContext(ctx).Model(&Reportrecord{})
	if err := db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// UpdateStatus 更新报告状态
func (r *Reportrecord) UpdateStatus(ctx context.Context, id int, status int) error {
	var db = mysql.FromContext(ctx).Model(&Reportrecord{})
	if err := db.Where("id = ?", id).
		Update("status", status).
		Error; err != nil {
		return err
	}
	return nil
}

// 获取待生成的报告
func (r *Reportrecord) GetsByStatus(ctx context.Context, status int) []Reportrecord {
	var (
		reportrecord []Reportrecord
		db           = mysql.FromContext(ctx).Model(&Reportrecord{})
	)
	db.Where("status = ?", status).Find(&reportrecord)

	return reportrecord
}
