package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

// AssetConnection 资产连接信息
type AssetConnection struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	AssetID    int       `gorm:"column:asset_id" json:"assetId"`       // 关联资产ID
	IP         string    `gorm:"column:ip" json:"ip"`                  // 关联资产IP
	Port       int       `gorm:"column:port" json:"port"`              // 端口
	Protocol   int       `gorm:"column:protocol" json:"protocol"`      // 协议/连接方式，如 ssh、rdp、mysql
	Username   string    `gorm:"column:username" json:"username"`      // 用户名
	Password   string    `gorm:"column:password" json:"-"`             // 密码（加密存储，不直接返回）
	Remark     string    `gorm:"column:remark" json:"remark"`          // 备注
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

// TableName 表名
func (ac *AssetConnection) TableName() string {
	return "asset_connection"
}

// CreateAssetConnection 新增
func (ac *AssetConnection) CreateAssetConnection(ctx context.Context) error {
	return mysql.FromContext(ctx).Model(&AssetConnection{}).Create(ac).Error
}

// UpdateAssetConnection 更新
func (ac *AssetConnection) UpdateAssetConnection(ctx context.Context) error {
	return mysql.FromContext(ctx).Model(&AssetConnection{}).Where("id = ?", ac.ID).Updates(ac).Error
}

// UpsertByIP 有则更新，无则新增（基于 IP 唯一）
func (ac *AssetConnection) UpsertByIP(ctx context.Context) error {
	db := mysql.FromContext(ctx).Model(&AssetConnection{})
	res := db.Where("ip = ?", ac.IP).Updates(map[string]interface{}{
		"port":        ac.Port,
		"protocol":    ac.Protocol,
		"username":    ac.Username,
		"password":    ac.Password,
		"remark":      ac.Remark,
		"update_time": time.Now(),
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		ac.CreateTime = time.Now()
		ac.UpdateTime = time.Now()
		return db.Create(ac).Error
	}
	return nil
}

// DeleteAssetConnection 删除
func (ac *AssetConnection) DeleteAssetConnection(ctx context.Context, id int) error {
	return mysql.FromContext(ctx).Where("id = ?", id).Delete(&AssetConnection{}).Error
}

// GetByAssetID 查询某资产的连接信息
func (ac *AssetConnection) GetByAssetID(ctx context.Context, assetID int) ([]AssetConnection, error) {
	var list []AssetConnection
	err := mysql.FromContext(ctx).Model(ac).Where("asset_id = ?", assetID).Find(&list).Error
	return list, err
}

// GetByAssetIP 通过IP查询某资产的连接信息
func (ac *AssetConnection) GetByAssetIP(ctx context.Context, ip string) ([]AssetConnection, error) {
	var list []AssetConnection
	err := mysql.FromContext(ctx).Model(ac).Where("ip = ?", ip).Find(&list).Error
	return list, err
}

// GetByID 查询单条
func (ac *AssetConnection) GetByID(ctx context.Context, id int) (*AssetConnection, error) {
	var conn AssetConnection
	err := mysql.FromContext(ctx).Model(ac).Where("id = ?", id).First(&conn).Error
	return &conn, err
}

// GetByProtocol 按照协议拿到写一下所有连接方式
func (ac *AssetConnection) GetByProtocol(ctx context.Context, protocol int) ([]AssetConnection, error) {
	var list []AssetConnection
	err := mysql.FromContext(ctx).Model(ac).Where("protocol = ?", protocol).Find(&list).Error
	return list, err
}

// GetAll 拿到全部连接信息
func (ac *AssetConnection) GetAll(ctx context.Context) ([]AssetConnection, error) {
	var list []AssetConnection
	err := mysql.FromContext(ctx).Model(ac).Find(&list).Error
	return list, err
}

// GetAssetConnList 拿到连接信息
func (ac *AssetConnection) GetAssetConnList(ctx context.Context, ip string, port, protocol, page, size int) ([]AssetConnection, int64, error) {
	var (
		list  []AssetConnection
		count int64
	)
	db := mysql.FromContext(ctx).Model(ac)
	db.Count(&count)
	if ip != "" {
		db = db.Where("ip = ?", ip)
	}
	if port > 0 {
		db = db.Where("port = ?", port)
	}
	if protocol > 0 {
		db = db.Where("protocol = ?", protocol)
	}
	err := db.Limit(size).Offset(size * (page - 1)).Find(&list).Error
	return list, count, err
}
