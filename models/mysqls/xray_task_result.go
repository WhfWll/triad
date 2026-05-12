package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type XrayTaskResult struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`        // 主键
	XrayTaskID  int       `gorm:"column:xray_task_id" json:"xrayTaskID"`  // xray任务表ID
	Addr        string    `gorm:"column:addr" json:"addr"`                // 地址
	Payload     string    `gorm:"column:payload" json:"payload"`          // payload
	RequestInfo string    `gorm:"column:request_info" json:"requestInfo"` // 请求与响应信息
	Extra       string    `gorm:"column:extra" json:"extra"`              // 扩展信息
	Plugin      string    `gorm:"column:plugin" json:"plugin"`            // 使用到的插件/漏洞类型
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`   // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`   // 修改时间
}

// TableName sets insert table name for this struct type
func (x *XrayTaskResult) TableName() string {
	return "xray_task_result"
}

// Get retrieves a list of xrayTaskResult from database
func (x *XrayTaskResult) GetXrayTaskResultList(ctx context.Context, xrayId, page, limit int, search string) ([]XrayTaskResult, int64, error) {
	var (
		xrayTaskResultList []XrayTaskResult
		count              int64
		db                 = mysql.FromContext(ctx).Model(&XrayTaskResult{})
	)
	if search != "" {
		db = db.Where("addr like ?", "%"+search+"%")
	}
	db = db.Where("xray_task_id = ?", xrayId)
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&xrayTaskResultList)

	return xrayTaskResultList, count, nil
}

// Get retrieves a single record of xrayTaskResult from database
func (x *XrayTaskResult) GetXrayTaskResult(ctx context.Context) (XrayTaskResult, error) {
	var (
		xrayTaskResult XrayTaskResult
		err            error
		db             = mysql.FromContext(ctx).Model(&XrayTaskResult{})
	)

	curErr := db.Where("id = ?", x.ID).First(&xrayTaskResult).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return xrayTaskResult, err
}

// Add persists xrayTaskResult to database
func (x *XrayTaskResult) AddXrayTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&XrayTaskResult{})

	if err := db.Create(x).Error; err != nil {
		return err
	}

	return nil
}

// Update changes xrayTaskResult by id
func (x *XrayTaskResult) UpdateXrayTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&XrayTaskResult{})

	if err := db.Where("id = ?", x.ID).Updates(x).Error; err != nil {
		return err
	}

	return nil
}

// Delete xrayTaskResult by id
func (x *XrayTaskResult) DeleteXrayTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&XrayTaskResult{})

	//x.Estate = "deleted"
	x.UpdateTime = time.Now()
	if err := db.Where("id = ?", x.ID).Updates(x).Error; err != nil {
		return err
	}

	return nil
}

// 批量添加
func (x *XrayTaskResult) AddsXrayTaskResult(ctx context.Context, data []XrayTaskResult) error {
	var db = mysql.FromContext(ctx).Model(&XrayTaskResult{})

	if err := db.Create(data).Error; err != nil {
		return err
	}

	return nil
}

// 删除
func (x *XrayTaskResult) DelByXrayIds(ctx context.Context, xrayIds []int) error {
	var (
		db = mysql.FromContext(ctx).Model(&XrayTaskResult{})
	)

	if err := db.Where("xray_task_id in ?", xrayIds).Delete(x).Error; err != nil {
		return err
	}

	return nil
}
