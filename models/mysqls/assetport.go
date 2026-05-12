package mysqls

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"strconv"
	"time"
)

type Assetport struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	IP         string    `gorm:"column:ip" json:"ip"`                  // ip
	Port       int       `gorm:"column:port" json:"port"`              // 端口
	TaskID     int       `gorm:"column:task_id" json:"taskID"`         // 任务ID
	Protocol   string    `gorm:"column:protocol" json:"protocol"`      // 协议
	Service    string    `gorm:"column:service" json:"service"`        // 服务
	Assembly   string    `gorm:"column:assembly" json:"assembly"`      // 组件
	Remark     string    `gorm:"column:remark" json:"remark"`          // 备注
	Islive     int       `gorm:"column:islive" json:"islive"`          // 存活状态:1-存活,2-不存活
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

func (a *Assetport) TableName() string {
	return "asset_port"
}

// 根据ip数组和端口范围限制查询所有数据
func (a *Assetport) GetAssetportByIps(ctx context.Context, ips []string, ports []string, portrange [][2]string) []Assetport {
	var (
		assetportList []Assetport
		db            = mysql.FromContext(ctx).Model(&Assetport{})
	)
	db.Where("ip in ?", ips)
	if len(ports) > 0 && len(portrange) == 0 {
		db.Where("port in ?", ports)
	} else if len(ports) == 0 && len(portrange) > 0 {
		for i := 0; i < len(portrange); i++ {
			if i == 0 {
				db.Where("port >= ? and port <= ?", portrange[i][0], portrange[i][1])
				continue
			}
			db.Or("port >= ? and port <= ?", portrange[i][0], portrange[i][1])
		}
	} else if len(ports) > 0 && len(portrange) > 0 {
		db.Where("port in ?", ports)
		for i := 0; i < len(portrange); i++ {
			db.Or("port >= ? and port <= ?", portrange[i][0], portrange[i][1])
		}
	}
	db.Find(&assetportList)
	return assetportList
}

// Get retrieves a single record of assetport from database
func (a *Assetport) GetAssetport(ctx context.Context) (Assetport, error) {
	var (
		assetport Assetport
		err       error
		db        = mysql.FromContext(ctx).Model(&Assetport{})
	)

	curErr := db.Where("id = ?", a.ID).First(&assetport).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return assetport, err
}

// 新增一条数据
func (a *Assetport) AddAssetport(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetport{})
	if err := db.Create(a).Error; err != nil {
		return err
	}
	return nil
}

// 批量新增数据
func (a *Assetport) AddAssetportMany(ctx context.Context, data []Assetport) error {
	var db = mysql.FromContext(ctx).Model(&Assetport{})
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

// 批量删除数据（根据 IP）
func (a *Assetport) DeleteAssetPortByIPs(ctx context.Context, ips []string) error {
	if len(ips) == 0 {
		return nil // 没有要删的内容，直接返回
	}
	db := mysql.FromContext(ctx).Model(&Assetport{})
	if err := db.Where("ip IN ?", ips).Delete(&Assetport{}).Error; err != nil {
		return err
	}
	return nil
}

// 根据ip修改数据
func (a *Assetport) UpdateAssetportByIp(ctx context.Context, ip string, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetport{})
	if err := db.Where("ip = ?", ip).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据ip和端口修改数据
func (a *Assetport) UpdateAssetportByIpPort(ctx context.Context, ip string, port int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetport{})
	if err := db.Where("ip = ? and port = ?", ip, port).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据扫描ip表ids和存活状态查找数据
func (a *Assetport) GetAssetPortByIpIslive(ctx context.Context, assetIps []string, isLive int) []Assetport {
	var (
		assetPortList []Assetport
		db            = mysql.FromContext(ctx).Model(&Assetport{})
	)
	db.Where("islive = ? and ip in ?", isLive, assetIps).Find(&assetPortList)
	return assetPortList
}

// 根据ip获取端口数据
func (a *Assetport) GetAssetPortByIp(ctx context.Context, ip string) []Assetport {
	var (
		assetPortList []Assetport
		db            = mysql.FromContext(ctx).Model(&Assetport{})
	)
	db.Where("ip = ?", ip).Find(&assetPortList)
	return assetPortList
}

// GetIPsByPort 更具端口查询ip信息
func (a *Assetport) GetIPsByPort(ctx context.Context, port int, service, finger string) []Assetport {
	var (
		assetPortList []Assetport
		db            = mysql.FromContext(ctx).Model(&Assetport{})
		query         string
		args          []interface{}
	)
	db = db.Where("islive = " + strconv.Itoa(enums.TaskIpIsLiveTypeYes))
	query += "1 = ?"
	args = append(args, 1)
	if port != 0 {
		query += " and port = ?"
		args = append(args, port)
	}
	if service != "" {
		query += " AND service LIKE ?"
		args = append(args, "%"+service+"%")
	}
	if finger != "" {
		query += " AND assembly LIKE ?"
		args = append(args, "%"+finger+"%")
	}
	db.Where(query, args...).Find(&assetPortList)
	return assetPortList
}

// GetAssetPortByServiceAndIP 通过服务和IP返回对应去重端口信息
func (a *Assetport) GetAssetPortByServiceAndIP(ctx context.Context, service, ip string) []Assetport {
	var (
		assetPortList []Assetport
		db            = mysql.FromContext(ctx).Model(&Assetport{})
	)
	db.Where("service = ? AND ip = ?", service, ip).Find(&assetPortList)
	return assetPortList
}

// GetAssetPortByFingerAndIP 通过指纹和IP返回对应去重端口信息
func (a *Assetport) GetAssetPortByFingerAndIP(ctx context.Context, finger, ip string) []Assetport {
	var (
		assetPortList []Assetport
		db            = mysql.FromContext(ctx).Model(&Assetport{})
	)
	db.Where("assembly LIKE ? AND ip = ?", "%"+finger+"%", ip).Find(&assetPortList)
	return assetPortList
}

// UpsertAssetPortMany 插入或更新 Assetport 数据（根据 asset_id + port 进行判断）
func (a *Assetport) UpsertAssetPortMany(ctx context.Context, data []Assetport, taskID int) error {
	if len(data) == 0 {
		return nil
	}
	for _, port := range data {
		log.Printf("[DEBUG] Got port -taskID:%d IP: %s | Port: %d | Protocol: %s | Service: %s", taskID, port.IP, port.Port, port.Protocol, port.Service)
		var existing Assetport
		err := mysql.FromContext(ctx).Model(&Assetport{}).Where("ip = ? AND port = ?", port.IP, port.Port).
			First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 不存在就插入
				if err := mysql.FromContext(ctx).Model(&Assetport{}).Create(&port).Error; err != nil {
					return fmt.Errorf("insert failed for IP: %s, port: %d: %w", port.IP, port.Port, err)
				}
			} else {
				// 查询出错
				return fmt.Errorf("query failed for IP: %s, port: %d: %w", port.IP, port.Port, err)
			}
		} else {
			// 存在就更新
			if err := mysql.FromContext(ctx).
				Model(&Assetport{}).
				Where("id = ?", existing.ID).
				Omit("create_time").
				Updates(port).Error; err != nil {
				return fmt.Errorf("update failed for asset_id: %s, port: %d: %w", port.IP, port.Port, err)
			}

		}
	}
	return nil
}
