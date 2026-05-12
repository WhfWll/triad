package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Logicvul struct {
	ID            int       `gorm:"column:id;primary_key" json:"id"`             // 主键
	TaskID        int       `gorm:"column:task_id" json:"taskID"`                // 任务id
	TargetID      int       `gorm:"column:target_id" json:"targetID"`            // 目标id
	Pocname       string    `gorm:"column:pocname" json:"pocname"`               // 漏洞标识
	Name          string    `gorm:"column:name" json:"name"`                     // 漏洞名称
	Class         string    `gorm:"column:class" json:"class"`                   // 漏洞分类
	Type          int       `gorm:"column:type" json:"type"`                     // 漏洞类型
	Risk          int       `gorm:"column:risk" json:"risk"`                     // 风险等级
	Location      string    `gorm:"column:location" json:"location"`             // 漏洞位置
	Description   string    `gorm:"column:description" json:"description"`       // 漏洞描述
	FixSuggest    string    `gorm:"column:fix_suggest" json:"fixSuggest"`        // 修复建议
	VulParam      string    `gorm:"column:vul_param" json:"vulParam"`            // 漏洞参数
	VulResult     string    `gorm:"column:vul_result" json:"vulResult"`          // 漏洞结果
	VerMsg        string    `gorm:"column:ver_msg" json:"verMsg"`                // 验证报文
	DecisionVulID string    `gorm:"column:decision_vul_id" json:"decisionVulID"` // 漏洞id
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`        // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time" json:"updateTime"`        // 更新时间
	TargetUrl     string    `gorm:"column:target_url" json:"target_url"`         // 测试目标
	Status        int       `gorm:"column:status" json:"status"`                 // 漏洞状态
	VulID         string    `gorm:"column:vul_id" json:"vulID"`                  // 漏洞id
}

// TableName sets insert table name for this struct type
func (l *Logicvul) TableName() string {
	return "logic_vul"
}

// Get retrieves a list of logicvul from database
func (l *Logicvul) GetLogicvulList(ctx context.Context, taskId, page, limit int, search string) ([]Logicvul, int64, error) {
	var (
		logicvulList []Logicvul
		count        int64
		db           = mysql.FromContext(ctx).Model(&Logicvul{})
	)
	if search != "" {
		db = db.Where("name LIKE ?", "%"+search+"%")
	}
	db.Where("task_id = ?", taskId)
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&logicvulList)

	return logicvulList, count, nil
}

// Get retrieves a single record of logicvul from database
func (l *Logicvul) GetLogicvul(ctx context.Context) (Logicvul, error) {
	var (
		logicvul Logicvul
		err      error
		db       = mysql.FromContext(ctx).Model(&Logicvul{})
	)

	curErr := db.Where("id = ?", l.ID).First(&logicvul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return logicvul, err
}

// Add persists logicvul to database
func (l *Logicvul) AddLogicvul(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&Logicvul{})

	if err := db.Create(l).Error; err != nil {
		return l.ID, err
	}
	return l.ID, nil
}

// Update changes logicvul by id
func (l *Logicvul) UpdateLogicvul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logicvul{})

	if err := db.Where("id = ?", l.ID).Updates(l).Error; err != nil {
		return err
	}

	return nil
}

// Delete logicvul by id
func (l *Logicvul) DeleteLogicvul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Logicvul{})

	l.UpdateTime = time.Now()
	if err := db.Where("id = ?", l.ID).Delete(l).Error; err != nil {
		return err
	}

	return nil
}

// Update changes logicvul by id
func (l *Logicvul) UpdateLogicvulParam(ctx context.Context, id string, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Logicvul{})
	if err := db.Where("id = ?", id).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// Get retrieves a single record of logicvul from database
func (l *Logicvul) GetLogicVulByTaskId(ctx context.Context, taskId int) ([]Logicvul, error) {
	var (
		logicvul []Logicvul
		err      error
		db       = mysql.FromContext(ctx).Model(&Logicvul{})
	)
	curErr := db.Where("task_id = ?", taskId).Find(&logicvul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return logicvul, err
}

// Get retrieves a single record of logicvul from database
func (l *Logicvul) GetLogicVulByTargetId(ctx context.Context, targetId int) ([]Logicvul, error) {
	var (
		logicvul []Logicvul
		err      error
		db       = mysql.FromContext(ctx).Model(&Logicvul{})
	)
	curErr := db.Where("target_id = ?", targetId).Find(&logicvul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return logicvul, err
}

// Get retrieves a single record of logicvul from database
func (l *Logicvul) GetLogicVulByTargetIds(ctx context.Context, targetIds []int) ([]Logicvul, error) {
	var (
		logicvul []Logicvul
		err      error
		db       = mysql.FromContext(ctx).Model(&Logicvul{})
	)
	curErr := db.Where("target_id in ?", targetIds).Find(&logicvul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return logicvul, err
}
