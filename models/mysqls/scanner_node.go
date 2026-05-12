package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

type ScannerNode struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Name       string    `gorm:"column:name" json:"name"`              // 节点名称
	IP         string    `gorm:"column:ip" json:"ip"`                  // 节点IP
	Port       string    `gorm:"column:port" json:"port"`              // 节点端口
	Status     int       `gorm:"column:status" json:"status"`          // 状态 0不在线 1在线
	IsDisable  int       `gorm:"column:is_disable" json:"isDisable"`   // 禁用状态 0启用 1禁用
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 发生时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 发生时间
}

// TableName sets insert table name for this struct type
func (y *ScannerNode) TableName() string {
	return "scanner_node"
}

// Get retrieves a list of yakNode from database
func (y *ScannerNode) GetYakNodeList(ctx context.Context, page, limit int, search string) ([]ScannerNode, int64, error) {
	var (
		yakNodeList []ScannerNode
		count       int64
		db          = mysql.FromContext(ctx).Model(&ScannerNode{})
	)

	if search != "" {
		db = db.Where("name LIKE ? OR ip LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&yakNodeList)

	return yakNodeList, count, nil
}

// Get retrieves a single record of yakNode from database
func (y *ScannerNode) GetYakNode(ctx context.Context) (ScannerNode, error) {
	var (
		yakNode ScannerNode
		err     error
		db      = mysql.FromContext(ctx).Model(&ScannerNode{})
	)

	curErr := db.Where("id = ?", y.ID).First(&yakNode).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return yakNode, err
}

// Add persists yakNode to database
func (y *ScannerNode) AddYakNode(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&ScannerNode{})

	if err := db.Create(y).Error; err != nil {
		return err
	}

	return nil
}

// Update changes yakNode by id
func (y *ScannerNode) UpdateYakNode(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&ScannerNode{})

	if err := db.Where("id = ?", y.ID).Updates(y).Error; err != nil {
		return err
	}

	return nil
}

// Delete yakNode by id
func (y *ScannerNode) DeleteYakNode(ctx context.Context, ids []string) error {
	var db = mysql.FromContext(ctx).Model(&ScannerNode{})

	//y.Estate = "deleted"
	y.UpdateTime = time.Now()
	if err := db.Where("id in ?", ids).Delete(y).Error; err != nil {
		return err
	}

	return nil
}

// 设置是否禁用
func (y *ScannerNode) UpdateIsDisable(ctx context.Context, id, isDisbale int) error {
	var db = mysql.FromContext(ctx).Model(&ScannerNode{})

	if err := db.Where("id = ?", id).Update("is_disable", isDisbale).Error; err != nil {
		return err
	}

	return nil
}

// 设置在线离线状态
func (y *ScannerNode) UpdateStatus(ctx context.Context, id, status int) error {
	var db = mysql.FromContext(ctx).Model(&ScannerNode{})

	if err := db.Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// 所有可用节点
func (y *ScannerNode) AllEnbaleNode(ctx context.Context) []ScannerNode {
	var (
		yakNode []ScannerNode
		db      = mysql.FromContext(ctx).Model(&ScannerNode{})
	)

	db.
		Where("status = ?", enums.YakNodeStatusOnline).
		Where("is_disable = ?", enums.YakNodeIsDisableN).
		Find(&yakNode)

	return yakNode
}

// 所有节点
func (y *ScannerNode) AllNode(ctx context.Context) []ScannerNode {
	var (
		yakNode []ScannerNode
		db      = mysql.FromContext(ctx).Model(&ScannerNode{})
	)

	db.Find(&yakNode)

	return yakNode
}

// 依据ip与port获取节点
func (y *ScannerNode) GetByIpPort(ctx context.Context, ip, port string) (ScannerNode, error) {
	var (
		yakNode ScannerNode
		db      = mysql.FromContext(ctx).Model(&ScannerNode{})
	)

	if err := db.
		Where("ip = ?", ip).
		Where("port = ?", port).
		First(&yakNode).Error; err != nil {
		return yakNode, err
	}

	return yakNode, nil
}
