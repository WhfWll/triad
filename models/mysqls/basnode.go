package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Basnode struct {
	ID           int       `gorm:"column:id;primary_key" json:"id"`          // 主键
	Name         string    `gorm:"column:name" json:"name"`                  // 节点名称
	IP           string    `gorm:"column:ip" json:"ip"`                      // IP
	OnlineStatus int       `gorm:"column:online_status" json:"onlineStatus"` // 在线状态 1-在线 2-离线
	Status       int       `gorm:"column:status" json:"status"`              // 节点状态 1-启用 2-禁用
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`     // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`     // 修改时间
}

// TableName sets insert table name for this struct type
func (b *Basnode) TableName() string {
	return "bas_node"
}

// 根据id和状态获取node信息
func (b *Basnode) GetBasnodesByIds(ctx context.Context, ids any, onlinStatus, status int) ([]Basnode, error) {
	var (
		basnodeList []Basnode
		db          = mysql.FromContext(ctx).Model(&Basnode{})
	)
	db.Where("id IN ? and online_status = ? and status = ?", ids, onlinStatus, status).Find(&basnodeList)
	return basnodeList, nil
}

// 根据ip获取node信息
func (b *Basnode) GetBasNodeByIp(ctx context.Context, ip string) (Basnode, error) {
	var (
		basnode Basnode
		err     error
		db      = mysql.FromContext(ctx).Model(&Basnode{})
	)
	curErr := db.Where("ip = ?", ip).First(&basnode).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return basnode, err
}

// 新增一条节点数据
func (b *Basnode) AddBasnode(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Basnode{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// 更新节点数据
func (b *Basnode) UpdateBasnode(ctx context.Context, id int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Basnode{})
	if err := db.Where("id = ?", id).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

//更新多条数据
func (b *Basnode) UpdateBasnodeByIds(ctx context.Context, ids any, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Basnode{})
	if err := db.Where("id in ?", ids).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 节点列表
func (b *Basnode) GetBasNodeList(ctx context.Context, page, limit int, search string) ([]Basnode, int64, error) {
	var (
		basTaskList []Basnode
		count       int64
		db          = mysql.FromContext(ctx).Model(&Basnode{})
	)
	if search != "" {
		db = db.Where("name like ?", "%"+search+"%")
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id DESC").Find(&basTaskList)

	return basTaskList, count, nil
}

//根据状态和在线状态查询agent数据
func (b *Basnode) GetBasNodeListByStatusAndOnlinestatus(ctx context.Context, status, onlineStatus int) []Basnode {
	var (
		basTaskList []Basnode
		db          = mysql.FromContext(ctx).Model(&Basnode{})
	)
	db.Where("status = ? and online_status = ?", status, onlineStatus).Order("id DESC").Find(&basTaskList)
	return basTaskList
}

// 根据在线状态和更新时间查询agent数据
func (b *Basnode) GetBasAgentByOnlineStatusAndUpdatetime(ctx context.Context, onlineStatus int, updateTime string) []Basnode {
	var (
		basTaskList []Basnode
		db          = mysql.FromContext(ctx).Model(&Basnode{})
	)
	db.Where("online_status = ? and update_time <= ?", onlineStatus, updateTime).Find(&basTaskList)
	return basTaskList
}
