package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type BurpsuiteTask struct {
	ID           int       `gorm:"column:id;primary_key" json:"id"`           // 主键
	OriginTaskId int       `gorm:"column:origin_task_id" json:"originTaskId"` // burpsuite软件生成的任务ID
	TaskName     string    `gorm:"column:task_name" json:"taskName"`          // 任务名称
	Target       string    `gorm:"column:target" json:"target"`               // 目标
	Risk         int       `gorm:"column:risk" json:"risk"`                   // 风险等级
	IsCrawler    int       `gorm:"column:is_crawler" json:"isCrawler"`        // 是否开启爬虫 0为开始 1默认开启（调用API后它是默认开启的，不能更改）
	Status       int       `gorm:"column:status" json:"status"`               // 状态 1待运行 2运行中 3已完成
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`      // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`      // 修改时间
}

// TableName sets insert table name for this struct type
func (b *BurpsuiteTask) TableName() string {
	return "burpsuite_task"
}

// Get retrieves a list of burpsuiteTask from database
func (b *BurpsuiteTask) GetBurpsuiteTaskList(ctx context.Context, page, limit int, search string) ([]BurpsuiteTask, int64, error) {
	var (
		burpsuiteTaskList []BurpsuiteTask
		count             int64
		db                = mysql.FromContext(ctx).Model(&BurpsuiteTask{})
	)

	if search != "" {
		db = db.Where("task_name like ?", "%"+search+"%")
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("create_time DESC").Find(&burpsuiteTaskList)

	return burpsuiteTaskList, count, nil
}

// Get retrieves a single record of burpsuiteTask from database
func (b *BurpsuiteTask) GetBurpsuiteTask(ctx context.Context) (BurpsuiteTask, error) {
	var (
		burpsuiteTask BurpsuiteTask
		err           error
		db            = mysql.FromContext(ctx).Model(&BurpsuiteTask{})
	)

	curErr := db.Where("id = ?", b.ID).First(&burpsuiteTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return burpsuiteTask, err
}

// Add persists burpsuiteTask to database
func (b *BurpsuiteTask) AddBurpsuiteTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BurpsuiteTask{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// Update changes burpsuiteTask by id
func (b *BurpsuiteTask) UpdateBurpsuiteTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BurpsuiteTask{})

	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

// Delete burpsuiteTask by id
func (b *BurpsuiteTask) DeleteBurpsuiteTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BurpsuiteTask{})

	//b.Estate = "deleted"
	b.UpdateTime = time.Now()
	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

// 更新任务结果
func (x *BurpsuiteTask) GetsByStatus(ctx context.Context, status int) []BurpsuiteTask {
	var (
		res []BurpsuiteTask
		db  = mysql.FromContext(ctx).Model(&BurpsuiteTask{})
	)

	db.Where("status = ?", status).Find(&res)

	return res
}

// 更新任务状态
func (x *BurpsuiteTask) UpdateRunning(ctx context.Context, id, originTaskId, status int) error {
	var (
		db = mysql.FromContext(ctx).Model(&BurpsuiteTask{})
	)

	data := map[string]interface{}{
		"origin_task_id": originTaskId,
		"status":         status,
	}

	if err := db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// 更新任务状态
func (x *BurpsuiteTask) UpdateStatusById(ctx context.Context, id, status int) error {
	var (
		db = mysql.FromContext(ctx).Model(&BurpsuiteTask{})
	)

	if err := db.Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// 更新任务风险
func (x *BurpsuiteTask) UpdateRiskById(ctx context.Context, id, risk int) error {
	var (
		db = mysql.FromContext(ctx).Model(&BurpsuiteTask{})
	)

	if err := db.Where("id = ?", id).Update("risk", risk).Error; err != nil {
		return err
	}

	return nil
}

// 更新任务状态
func (x *BurpsuiteTask) UpdateStatusRiskById(ctx context.Context, id, status, risk int) error {
	var (
		db = mysql.FromContext(ctx).Model(&BurpsuiteTask{})
	)
	data := map[string]interface{}{
		"status": status,
		"risk":   risk,
	}
	if err := db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// 通过IDs删除
func (x *BurpsuiteTask) DelByIds(ctx context.Context, ids []int) error {
	var (
		db = mysql.FromContext(ctx).Model(&BurpsuiteTask{})
	)

	if err := db.Where("id in ?", ids).Delete(x).Error; err != nil {
		return err
	}

	return nil
}
