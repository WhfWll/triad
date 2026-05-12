package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Assetscanip struct {
	ID               int       `gorm:"column:id;primary_key" json:"id"`                   // 主键
	ActivityID       int       `gorm:"column:activity_id" json:"activityID"`              // 所属活动id
	AssetScanTaskID  int       `gorm:"column:asset_scan_task_id" json:"assetScanTaskID"`  // 所属任务id
	IP               string    `gorm:"column:ip" json:"ip"`                               // ip
	PortRange        string    `gorm:"column:port_range" json:"portRange"`                // 端口扫描范围
	IpNum            int64     `gorm:"column:ip_num" json:"ipNum"`                        // ip值
	Os               string    `gorm:"column:os" json:"os"`                               // 操作系统
	Status           int       `gorm:"column:status" json:"status"`                       // 目标状态
	AssetChangesType int       `gorm:"column:asset_changes_type" json:"assetChangesType"` // 资产变化类型:1-未变化,2-已减少IP，3-新增加IP，4-端口变化IP，5-服务变化IP，6-组件变化IP
	Islive           int       `gorm:"column:islive" json:"islive"`                       // 存活状态:1-存活,2-不存活IP
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`              // 创建时间
	UpdateTime       time.Time `gorm:"column:update_time" json:"updateTime"`              // 修改时间
	UserID           int       `gorm:"column:user_id" json:"userID"`                      // 提交者id
	CmdPid           string    `gorm:"column:cmd_pid" json:"cmdPid"`                      // CmdPid
}

func (a *Assetscanip) TableName() string {
	return "asset_scan_ip"
}

// 列表查询
func (a *Assetscanip) GetAssetscanipList(ctx context.Context, taskId int, search string, status []string, islive, assetChangesType int,
	createTime string, page, limit int) ([]Assetscanip, int64, error) {
	var (
		assetscanipList []Assetscanip
		count           int64
		db              = mysql.FromContext(ctx).Model(&Assetscanip{})
		order           = "desc"
	)
	db = db.Where("asset_scan_task_id = ?", taskId)
	if islive > 0 {
		db = db.Where("islive = ?", islive)
	}
	if assetChangesType > 0 {
		db = db.Where("asset_changes_type = ?", assetChangesType)
	}
	if len(status) > 0 {
		db = db.Where("status in ?", status)
	}
	if len(search) > 0 {
		db = db.Where("ip like ? or os like ?", "%"+search+"%", "%"+search+"%")
	}
	if len(createTime) > 0 {
		order = createTime
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id " + order).Find(&assetscanipList)
	return assetscanipList, count, nil
}

// 列表查询
func (a *Assetscanip) GetAssetscanipAll(ctx context.Context, taskId int, search string, status []string, islive, assetChangesType int) []Assetscanip {
	var (
		assetscanipList []Assetscanip
		db              = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db = db.Where("asset_scan_task_id = ?", taskId)
	if islive > 0 {
		db = db.Where("islive = ?", islive)
	}
	if assetChangesType > 0 {
		db = db.Where("asset_changes_type = ?", assetChangesType)
	}
	if len(status) > 0 {
		db = db.Where("status in ?", status)
	}
	if len(search) > 0 {
		db = db.Where("ip like ? or os like ?", "%"+search+"%", "%"+search+"%")
	}
	db.Find(&assetscanipList)
	return assetscanipList
}

// 根据任务id查询所有扫描ip
func (a *Assetscanip) GetAssetscanipAllByTaskIds(ctx context.Context, taskIds any) []Assetscanip {
	var (
		assetscanipList []Assetscanip
		db              = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db.Where("asset_scan_task_id IN ?", taskIds).Find(&assetscanipList)
	return assetscanipList
}

// 根据任务id查询所有数据
func (a *Assetscanip) GetAllAssetscanipByTaskId(ctx context.Context, taskId int) []Assetscanip {
	var (
		assetscanipList []Assetscanip
		db              = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db.Where("asset_scan_task_id = ?", taskId).Find(&assetscanipList)
	return assetscanipList
}

// 根据任务id查询所有数据总数
func (a *Assetscanip) GetAssetscanipCountByTaskId(ctx context.Context, taskId int) (int64, error) {
	var (
		db    = mysql.FromContext(ctx).Model(&Assetscanip{})
		err   error
		count int64
	)
	curErr := db.Where("asset_scan_task_id = ?", taskId).Count(&count).Error
	if curErr != nil {
		err = curErr
	}
	return count, err
}

// 根据任务id查询所有数据总数
func (a *Assetscanip) GetAssetscanipCountByTaskIdStatus(ctx context.Context, taskId int, statusArray []int) (int64, error) {
	var (
		db    = mysql.FromContext(ctx).Model(&Assetscanip{})
		err   error
		count int64
	)
	curErr := db.Where("asset_scan_task_id = ? and status IN ?", taskId, statusArray).Count(&count).Error
	if curErr != nil {
		err = curErr
	}
	return count, err
}

// 根据任务id和存活状态查询总数
func (a *Assetscanip) CountScanIpByTaskIdIslive(ctx context.Context, taskId, islive int) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db.Where("asset_scan_task_id = ? and islive = ?", taskId, islive).Count(&count)
	return count
}

// 根据任务id统计资产变化情况
// SELECT asset_changes_type,COUNT(id) as total FROM `asset_scan_ip` where asset_scan_task_id =2 GROUP BY asset_changes_type;
type ScanIpChangeType struct {
	AssetChangesType int   `gorm:"column:asset_changes_type" json:"assetChangesType"`
	Total            int64 `gorm:"column:total" json:"total"`
}

func (a *Assetscanip) CountScanIpChangeTypeByTaskId(ctx context.Context, taskId int) []ScanIpChangeType {
	var (
		result []ScanIpChangeType
		db     = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db.Select("asset_changes_type,COUNT(id) as total").Where("asset_scan_task_id = ?", taskId).Group("asset_changes_type").Find(&result)
	return result
}

// 根据状态查询数据
func (a *Assetscanip) GetScanTaskIpByStatus(ctx context.Context, status int) ([]Assetscanip, error) {
	var (
		assetscanipList []Assetscanip
		db              = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db.Where("status = ?", status).Find(&assetscanipList)
	return assetscanipList, nil
}

// 根据状态查询一条数据
func (a *Assetscanip) GetScanTaskIpByStatusFirst(ctx context.Context, status int) (Assetscanip, error) {
	var (
		assetscanip Assetscanip
		err         error
		db          = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	curErr := db.Where("status = ?", status).First(&assetscanip).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assetscanip, err
}

// 根据ids和状态查询数据
func (a *Assetscanip) GetScanTaskIpByIdsStatus(ctx context.Context, ids any, status int) []Assetscanip {
	var (
		assetscanipList []Assetscanip
		db              = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db.Where(" status = ? and id in ?", status, ids).Find(&assetscanipList)
	return assetscanipList
}

// 根据id查询一条数据
func (a *Assetscanip) GetAssetscanip(ctx context.Context, id int) (Assetscanip, error) {
	var (
		assetscanip Assetscanip
		err         error
		db          = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	curErr := db.Where("id = ?", id).First(&assetscanip).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assetscanip, err
}

// 插入一条数据
func (a *Assetscanip) AddAssetscanip(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanip{})
	if err := db.Create(a).Error; err != nil {
		return err
	}
	return nil
}

// 修改数据
func (a *Assetscanip) UpdateAssetscanip(ctx context.Context, id int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanip{})
	if err := db.Where("id = ?", id).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据id和状态修改数据
func (a *Assetscanip) UpdateAssetscanipByIdStatus(ctx context.Context, id, status int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanip{})
	if err := db.Where("id = ? and status = ?", id, status).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据任务id修改数据
func (a *Assetscanip) UpdateAssetscanipByTaskId(ctx context.Context, taskId int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanip{})
	if err := db.Where("asset_scan_task_id = ?", taskId).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据任务id批量修改数据
func (a *Assetscanip) UpdateAssetscanipByTaskIdMany(ctx context.Context, scanIpIds []int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanip{})
	if err := db.Where("id IN ?", scanIpIds).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据状态和更新时间修改数据
func (a *Assetscanip) UpdateScanTaskIpByStatusUpdateTime(ctx context.Context, status int, updateTime string, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanip{})
	if err := db.Where("status = ? and update_time <= ?", status, updateTime).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 根据状态和更新时间查询数据
func (a *Assetscanip) GetScanTaskIpByStatusUpdateTime(ctx context.Context, status int, updateTime string) []Assetscanip {
	var (
		assetscanipList []Assetscanip
		db              = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db.Where("status = ? and update_time <= ?", status, updateTime).Find(&assetscanipList)
	return assetscanipList
}

// 根据id删除扫描ip
func (a *Assetscanip) DeleteScanTaskIpByIds(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanip{})
	if err := db.Where("id IN ?", ids).Delete(&Assetscanip{}).Error; err != nil {
		return err
	}
	return nil
}

// 根据任务id删除扫描ip
func (a *Assetscanip) DeleteScanTaskIpByTaskIds(ctx context.Context, taskIds any) error {
	var db = mysql.FromContext(ctx).Model(&Assetscanip{})
	if err := db.Where("asset_scan_task_id IN ?", taskIds).Delete(&Assetscanip{}).Error; err != nil {
		return err
	}
	return nil
}

// GetAllAssetScanIPByActivityId 根据活动id查询所有数据
func (a *Assetscanip) GetAllAssetScanIPByActivityId(ctx context.Context, activityID int) []Assetscanip {
	var (
		assetscanipList []Assetscanip
		db              = mysql.FromContext(ctx).Model(&Assetscanip{})
	)
	db.Where("activity_id = ?", activityID).Find(&assetscanipList)
	return assetscanipList
}
