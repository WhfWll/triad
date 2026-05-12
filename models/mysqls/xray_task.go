package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

type XrayTask struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	TaskName   string    `gorm:"column:task_name" json:"taskName"`     // 任务名称
	Target     string    `gorm:"column:target" json:"target"`          // 目标
	RiskNum    int       `gorm:"column:risk_num" json:"riskNum"`       // 风险数量
	IsCrawler  int       `gorm:"column:is_crawler" json:"isCrawler"`   // 是否开启爬虫 0为开始 1开启
	Status     int       `gorm:"column:status" json:"status"`          // 状态
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

// TableName sets insert table name for this struct type
func (x *XrayTask) TableName() string {
	return "xray_task"
}

// Get retrieves a list of xrayTask from database
func (x *XrayTask) GetXrayTaskList(ctx context.Context, page, limit int, search string) ([]XrayTask, int64, error) {
	var (
		xrayTaskList []XrayTask
		count        int64
		db           = mysql.FromContext(ctx).Model(&XrayTask{})
	)

	if search != "" {
		db = db.Where("task_name like ?", "%"+search+"%")
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("create_time DESC").Find(&xrayTaskList)

	return xrayTaskList, count, nil
}

// Get retrieves a single record of xrayTask from database
func (x *XrayTask) GetXrayTask(ctx context.Context) (XrayTask, error) {
	var (
		xrayTask XrayTask
		err      error
		db       = mysql.FromContext(ctx).Model(&XrayTask{})
	)

	curErr := db.Where("id = ?", x.ID).First(&xrayTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return xrayTask, err
}

// Add persists xrayTask to database
func (x *XrayTask) AddXrayTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&XrayTask{})

	if err := db.Create(x).Error; err != nil {
		return err
	}

	return nil
}

// Update changes xrayTask by id
func (x *XrayTask) UpdateXrayTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&XrayTask{})

	if err := db.Where("id = ?", x.ID).Updates(x).Error; err != nil {
		return err
	}

	return nil
}

// Delete xrayTask by id
func (x *XrayTask) DeleteXrayTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&XrayTask{})

	//x.Estate = "deleted"
	x.UpdateTime = time.Now()
	if err := db.Where("id = ?", x.ID).Updates(x).Error; err != nil {
		return err
	}

	return nil
}

// 更新任务结果
func (x *XrayTask) Finish(ctx context.Context, xrayId, riskNum int) error {
	var db = mysql.FromContext(ctx).Model(&XrayTask{})

	data := map[string]interface{}{
		"status":   enums.XrayStatusDone,
		"risk_num": riskNum,
	}
	x.UpdateTime = time.Now()
	if err := db.Where("id = ?", xrayId).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// 更新任务结果
func (x *XrayTask) GetsByStatus(ctx context.Context, status int) []XrayTask {
	var (
		res []XrayTask
		db  = mysql.FromContext(ctx).Model(&XrayTask{})
	)

	db.Where("status = ?", status).Find(&res)

	return res
}

// 更新任务状态
func (x *XrayTask) UpdateStatusById(ctx context.Context, id, status int) error {
	var (
		db = mysql.FromContext(ctx).Model(&XrayTask{})
	)

	if err := db.Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// 删除
func (x *XrayTask) DelByIds(ctx context.Context, ids []int) error {
	var (
		db = mysql.FromContext(ctx).Model(&XrayTask{})
	)

	if err := db.Where("id in ?", ids).Delete(x).Error; err != nil {
		return err
	}

	return nil
}
