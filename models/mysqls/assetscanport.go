package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Assetscanport struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`                  // 主键
	ActivityID      int       `gorm:"column:activity_id" json:"activityID"`             // 所属活动id
	AssetScanTaskID int       `gorm:"column:asset_scan_task_id" json:"assetScanTaskID"` // 所属任务id
	AssetScanIPID   int       `gorm:"column:asset_scan_ip_id" json:"assetScanIPID"`     // 所属IP表id
	Status          int       `gorm:"column:status" json:"status"`                      // 状态
	IP              string    `gorm:"column:ip" json:"ip"`                              // ip
	Port            int       `gorm:"column:port" json:"port"`                          // 端口
	Protocol        string    `gorm:"column:protocol" json:"protocol"`                  // 协议
	Service         string    `gorm:"column:service" json:"service"`                    // 服务
	Assembly        string    `gorm:"column:assembly" json:"assembly"`                  // 组件
	Remark          string    `gorm:"column:remark" json:"remark"`                      // 变化说明
	Islive          int       `gorm:"column:islive" json:"islive"`                      // 存活状态:1-存活,2-不存活
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`             // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`             // 修改时间
	UserID          int       `gorm:"column:user_id" json:"userID"`                     // 提交者id
}

func (a *Assetscanport) TableName() string {
	return "asset_scan_port"
}

// 列表查询
func (a *Assetscanport) GetAssetscanportList(ctx context.Context, scanIpId, page, limit int) ([]Assetscanport, int64, error) {
	var (
		assetscanportList []Assetscanport
		count             int64
		db                = mysql.FromContext(ctx).Model(&Assetscanport{})
	)
	db.Where("asset_scan_ip_id = ?", scanIpId)
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&assetscanportList)
	return assetscanportList, count, nil
}

// 根据任务id查询所有扫描端口数据
func (a *Assetscanport) GetAllAssetscanportByTaskId(ctx context.Context, taskId int) []Assetscanport {
	var (
		assetscanportList []Assetscanport
		db                = mysql.FromContext(ctx).Model(&Assetscanport{})
	)
	db.Where("asset_scan_task_id = ?", taskId).Find(&assetscanportList)
	return assetscanportList
}

// 根据任务id查询所有扫描端口数据
func (a *Assetscanport) GetAssetscanportCountByTaskId(ctx context.Context, taskId int) (int64, error) {
	var (
		db    = mysql.FromContext(ctx).Model(&Assetscanport{})
		err   error
		count int64
	)
	curErr := db.Where("asset_scan_task_id = ?", taskId).Count(&count).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return count, err
}

// 根据任务id和状态查询所有扫描端口数据
func (a *Assetscanport) GetAssetscanportCountByTaskIdStatus(ctx context.Context, taskId int, statusArray []int) (int64, error) {
	var (
		db    = mysql.FromContext(ctx).Model(&Assetscanport{})
		err   error
		count int64
	)
	curErr := db.Where("asset_scan_task_id = ? and status IN ?", taskId, statusArray).Count(&count).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return count, err
}

// 根据状态查询数据
func (a *Assetscanport) GetAssetscanportByStatus(ctx context.Context, status int) ([]Assetscanport, error) {
	var (
		assetscanportList []Assetscanport
		err               error
		db                = mysql.FromContext(ctx).Model(&Assetscanport{})
	)
	curErr := db.Where("status = ?", status).Find(&assetscanportList).Error
	if curErr != nil {
		err = curErr
	}
	return assetscanportList, err
}

// 根据扫描ip表ids和存活状态查找数据
func (a *Assetscanport) GetScanTaskPortByIpIslive(ctx context.Context, scanIpIds []int, isLive int) []Assetscanport {
	var (
		assetscanportList []Assetscanport
		db                = mysql.FromContext(ctx).Model(&Assetscanport{})
	)
	db.Where("islive = ? and asset_scan_ip_id in ?", isLive, scanIpIds).Find(&assetscanportList)
	return assetscanportList
}

// 根据id查询数据
func (a *Assetscanport) GetAssetscanport(ctx context.Context, id int) (Assetscanport, error) {
	var (
		assetscanport Assetscanport
		err           error
		db            = mysql.FromContext(ctx).Model(&Assetscanport{})
	)
	curErr := db.Where("id = ?", id).First(&assetscanport).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assetscanport, err
}

// 新增一条数据
func (a *Assetscanport) AddAssetscanport(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Create(a).Error; err != nil {
		return err
	}
	return nil
}

// 批量新增数据
func (a *Assetscanport) AddAssetscanportMany(ctx context.Context, datas []Assetscanport) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Create(datas).Error; err != nil {
		return err
	}
	return nil
}

// 根据id修改数据
func (a *Assetscanport) UpdateAssetscanport(ctx context.Context, id int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Where("id = ?", id).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据id和状态修改数据
func (a *Assetscanport) UpdateAssetscanportByIdStatus(ctx context.Context, id, status int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Where("id = ? and status = ?", id, status).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据任务id修改数据
func (a *Assetscanport) UpdateAssetscanportByTaskId(ctx context.Context, taskId int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Where("asset_scan_task_id = ?", taskId).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据扫描ip的id修改数据
func (a *Assetscanport) UpdateAssetscanportByScanIpId(ctx context.Context, assetScanIpId int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Where("asset_scan_ip_id = ?", assetScanIpId).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据扫描ip和端口的id修改数据
func (a *Assetscanport) UpdateAssetscanportByScanIpIdPort(ctx context.Context, assetScanIpId int, port int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Where("asset_scan_ip_id = ? and port = ?", assetScanIpId, port).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据状态和更新时间修改数据
func (a *Assetscanport) UpdateScanTaskPortByStatusUpdateTime(ctx context.Context, status int, updateTimme string, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Where("status = ? and update_time <= ?", status, updateTimme).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据任务id删除扫描ip
func (a *Assetscanport) DeleteScanTaskPortByTaskIds(ctx context.Context, taskIds any) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Where("asset_scan_task_id IN ?", taskIds).Delete(&Assetscanport{}).Error; err != nil {
		return err
	}
	return nil
}

// 根据扫描任务ip的id删除扫描ip
func (a *Assetscanport) DeleteScanTaskPortByScanIpIds(ctx context.Context, scanIpIds any) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanport{})
	if err := db.Where("asset_scan_ip_id IN ?", scanIpIds).Delete(&Assetscanport{}).Error; err != nil {
		return err
	}
	return nil
}
