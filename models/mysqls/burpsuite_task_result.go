package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type BurpsuiteTaskResult struct {
	ID                    int       `gorm:"column:id;primary_key" json:"id"`                            // 主键
	BurpsuiteTaskID       int       `gorm:"column:burpsuite_task_id" json:"burpsuiteTaskID"`            // burpsuite任务表ID
	OriginResultId        int       `gorm:"column:origin_result_id" json:"originResultId"`              // burpsuite软件生成的任务结果ID
	Action                string    `gorm:"column:action" json:"action"`                                // 对应json结果的外层type
	IssueType             string    `gorm:"column:issue_type" json:"issueType"`                         // 问题类型 取原数据的issue.name
	Host                  string    `gorm:"column:host" json:"host"`                                    // 主机
	Path                  string    `gorm:"column:path" json:"path"`                                    // 路径
	InsertionPoint        string    `gorm:"column:insertion_point" json:"insertionPoint"`               // 插入点
	Severity              string    `gorm:"column:severity" json:"severity"`                            // 严重程度 也是风险等级 取原数据的即可
	Confidence            string    `gorm:"column:confidence" json:"confidence"`                        // confidence 取原数据的即可
	Describe              string    `gorm:"column:describe" json:"describe"`                            // 漏洞描述 取原数据即可
	IssueBackground       string    `gorm:"column:issue_background" json:"issueBackground"`             // 问题背景 取原数据即可
	RemediationBackground string    `gorm:"column:remediation_background" json:"remediationBackground"` // 补救背景 取原数据即可
	RequestResponse       string    `gorm:"column:request_response" json:"requestResponse"`             // 请求与响应信息 取原数据即可
	InternalData          string    `gorm:"column:internal_data" json:"internalData"`                   // 内部数据 取原数据即可
	CreateTime            time.Time `gorm:"column:create_time" json:"createTime"`                       // 创建时间
	UpdateTime            time.Time `gorm:"column:update_time" json:"updateTime"`                       // 修改时间
}

// TableName sets insert table name for this struct type
func (b *BurpsuiteTaskResult) TableName() string {
	return "burpsuite_task_result"
}

// Get retrieves a list of burpsuiteTaskResult from database
func (b *BurpsuiteTaskResult) GetBurpsuiteTaskResultList(ctx context.Context, burpsuiteId, page, limit int, search string) ([]BurpsuiteTaskResult, int64, error) {
	var (
		burpsuiteTaskResultList []BurpsuiteTaskResult
		count                   int64
		db                      = mysql.FromContext(ctx).Model(&BurpsuiteTaskResult{})
	)

	if search != "" {
		db = db.Where("host like ?", "%"+search+"%")
	}
	db = db.Where("burpsuite_task_id = ?", burpsuiteId)
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&burpsuiteTaskResultList)

	return burpsuiteTaskResultList, count, nil
}

// Get retrieves a single record of burpsuiteTaskResult from database
func (b *BurpsuiteTaskResult) GetBurpsuiteTaskResult(ctx context.Context) (BurpsuiteTaskResult, error) {
	var (
		burpsuiteTaskResult BurpsuiteTaskResult
		err                 error
		db                  = mysql.FromContext(ctx).Model(&BurpsuiteTaskResult{})
	)

	curErr := db.Where("id = ?", b.ID).First(&burpsuiteTaskResult).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return burpsuiteTaskResult, err
}

// Add persists burpsuiteTaskResult to database
func (b *BurpsuiteTaskResult) AddBurpsuiteTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BurpsuiteTaskResult{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// Update changes burpsuiteTaskResult by id
func (b *BurpsuiteTaskResult) UpdateBurpsuiteTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BurpsuiteTaskResult{})

	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

// Delete burpsuiteTaskResult by id
func (b *BurpsuiteTaskResult) DeleteBurpsuiteTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BurpsuiteTaskResult{})

	//b.Estate = "deleted"
	b.UpdateTime = time.Now()
	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

// Delete burpsuiteTaskResult by id
func (b *BurpsuiteTaskResult) GetHostAndPathByBurpsuiteId(ctx context.Context, burpsuiteId []int) []BurpsuiteTaskResult {
	var (
		db   = mysql.FromContext(ctx).Model(&BurpsuiteTaskResult{})
		data []BurpsuiteTaskResult
	)

	db.Where("burpsuite_task_id in ?", burpsuiteId).Find(&data)

	return data
}

// AddsBurpsuiteTaskResult
func (b *BurpsuiteTaskResult) AddsBurpsuiteTaskResult(ctx context.Context, datas []BurpsuiteTaskResult) error {
	var db = mysql.FromContext(ctx).Model(&BurpsuiteTaskResult{})

	if err := db.Create(&datas).Error; err != nil {
		return err
	}

	return nil
}

// 通过burpsuiteIds删除
func (x *BurpsuiteTaskResult) DelByBurpsuiteIds(ctx context.Context, burpsuiteIds []int) error {
	var (
		db = mysql.FromContext(ctx).Model(&BurpsuiteTaskResult{})
	)

	if err := db.Where("burpsuite_task_id in ?", burpsuiteIds).Delete(x).Error; err != nil {
		return err
	}

	return nil
}
