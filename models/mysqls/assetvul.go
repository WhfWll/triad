package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Assetvul struct {
	ID             int       `gorm:"column:id;primary_key" json:"id"`               // 主键
	TaskVulID      int       `gorm:"column:task_vul_id" json:"taskVulID"`           // 所属任务漏洞id
	AssetID        int       `gorm:"column:asset_id" json:"assetID"`                // 所属资产id
	TaskID         int       `gorm:"column:task_id" json:"taskID"`                  // 所属任务id
	TargetID       int       `gorm:"column:target_id" json:"targetID"`              // 所属目标id
	TargetURL      string    `gorm:"column:target_url" json:"targetURL"`            // 测试目标地址
	Pocname        string    `gorm:"column:pocname" json:"pocname"`                 // 漏洞标识
	Name           string    `gorm:"column:name" json:"name"`                       // 漏洞名称
	Class          int       `gorm:"column:class" json:"class"`                     // 漏洞分类
	Type           int       `gorm:"column:type" json:"type"`                       // 漏洞类型
	Risk           int       `gorm:"column:risk" json:"risk"`                       // 风险等级，1-致命/2-高危/3-中危/4-低危/5-信息
	Location       string    `gorm:"column:location" json:"location"`               // 漏洞位置或地址
	Status         int       `gorm:"column:status" json:"status"`                   // 漏洞状态，1-未验证，2-验证成功，3-利用成功，4-已修复，5-流量验证,6-未知,7-存在,8-不存在
	TestStatus     int       `gorm:"column:test_status" json:"testStatus"`          // 测试状态，1-未测试，2-测试中，3-已测试
	ExploitImpact  string    `gorm:"column:exploit_impact" json:"exploitImpact"`    // 利用影响
	VulID          string    `gorm:"column:vul_id" json:"vulID"`                    // 漏洞id
	Description    string    `gorm:"column:description" json:"description"`         // 漏洞描述
	FixSuggest     string    `gorm:"column:fix_suggest" json:"fixSuggest"`          // 修复建议
	PublishedTime  string    `gorm:"column:published_time" json:"publishedTime"`    // 披露时间
	AffectRange    string    `gorm:"column:affect_range" json:"affectRange"`        // 影响范围
	TargetResultID int       `gorm:"column:target_result_id" json:"targetResultID"` // 检测结果id
	VulNumber      string    `gorm:"column:vul_number" json:"vulNumber"`            // 漏洞编号
	VulAddress     string    `gorm:"column:vul_address" json:"vulAddress"`          // 漏洞地址
	RefURL         string    `gorm:"column:ref_url" json:"refURL"`                  // 参考链接
	Cvss           string    `gorm:"column:cvss" json:"cvss"`                       // cvss评分
	VulResult      string    `gorm:"column:vul_result" json:"vulResult"`            // 漏洞结果
	VulParam       string    `gorm:"column:vul_param" json:"vulParam"`              // 漏洞请求参数
	VerMsg         string    `gorm:"column:ver_msg" json:"verMsg"`                  // 验证报文
	IsReplace      int       `gorm:"column:is_replace" json:"isReplace"`            // 是否替换漏洞，1-否，2-是
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`          // 创建时间
	UpdateTime     time.Time `gorm:"column:update_time" json:"updateTime"`          // 修改时间
}

// TableName sets insert table name for this struct type
func (a *Assetvul) TableName() string {
	return "asset_vul"
}

// Get retrieves a list of assetvul from database
func (a *Assetvul) GetAssetvulList(ctx context.Context, page, limit int) ([]Assetvul, int64, error) {
	var (
		assetvulList []Assetvul
		count        int64
		db           = mysql.FromContext(ctx).Model(&Assetvul{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&assetvulList)
	db.Count(&count)

	return assetvulList, count, nil
}

// 根据漏洞名称模糊查询关联的资产id
func (a *Assetvul) GetAssetIdsByName(ctx context.Context, name string) []Assetvul {
	var (
		assetvulList []Assetvul
		db           = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Select("asset_id").Where("name like ?", "%"+name+"%").Distinct("asset_id").Find(&assetvulList)
	return assetvulList
}

// 查询资产漏洞总数
func (a *Assetvul) GetAssetvulTotal(ctx context.Context) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Count(&count)
	return count
}

// 查询大于某个创建时间的资产漏洞总数
func (a *Assetvul) GetAssetvulTotalGtCreatTime(ctx context.Context, createTime string) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Where("create_time >= ?", createTime).Count(&count)
	return count
}

// 按照资产漏洞危险等级计算数量
type GetAssetVulTotalByRisk struct {
	Risk  int `gorm:"column:risk" json:"risk"`
	Total int `gorm:"column:total" json:"total"`
}

func (a *Assetvul) GetAssetVulTotalByRisk(ctx context.Context) []GetAssetVulTotalByRisk {
	var (
		result []GetAssetVulTotalByRisk
		db     = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Select("risk,COUNT(id) as total").Group("risk").Find(&result)
	return result
}

// 漏洞发现趋势
type AssetVulStat struct {
	Date  string `gorm:"column:date" json:"date"`
	Risk  int    `gorm:"column:risk" json:"risk"`
	Count int64  `gorm:"column:count" json:"count"`
}

func (a *Assetvul) GetAssetVulstat(ctx context.Context, startTime, dateFormat string) []AssetVulStat {
	var (
		assetVulStatRes []AssetVulStat
		db              = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	if startTime != "" {
		db = db.Where("update_time > ?", startTime)
	}
	db.Select("DATE_FORMAT(update_time, '" + dateFormat + "') as date,risk, COUNT(id) as count").Group("DATE_FORMAT(update_time, '" + dateFormat + "'),risk").Find(&assetVulStatRes)
	return assetVulStatRes
}

// Get retrieves a single record of assetvul from database
func (a *Assetvul) GetAssetvul(ctx context.Context) (Assetvul, error) {
	var (
		assetvul Assetvul
		err      error
		db       = mysql.FromContext(ctx).Model(&Assetvul{})
	)

	curErr := db.Where("id = ?", a.ID).First(&assetvul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return assetvul, err
}

// Add persists assetvul to database
func (a *Assetvul) AddAssetvul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetvul{})

	if err := db.Create(a).Error; err != nil {
		return err
	}

	return nil
}

// Update changes assetvul by id
func (a *Assetvul) UpdateAssetvul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetvul{})

	if err := db.Where("id = ?", a.ID).Updates(a).Error; err != nil {
		return err
	}

	return nil
}

// 漏洞类型统计
type AssetVulTypeStat struct {
	Type  int `gorm:"column:type" json:"type"`
	Risk  int `gorm:"column:risk" json:"risk"`
	Count int `gorm:"column:count" json:"count"`
}

func (a *Assetvul) GetAssetVulTypeStat(ctx context.Context) []AssetVulTypeStat {
	var (
		assetVulTypeStatList []AssetVulTypeStat
		db                   = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Select("type,risk, Count(id) as count").Group("type,risk").Order("risk,type").Find(&assetVulTypeStatList)
	return assetVulTypeStatList
}

// 批量添加资产信息收集数据
func (a *Assetvul) MultipartInsert(ctx context.Context, list *[]Assetvul) error {
	var db = mysql.FromContext(ctx).Model(&Assetvul{})
	if err := db.Create(list).Error; err != nil {
		return err
	}
	return nil
}

// DeleteByVulId 通过漏洞ID删除漏洞
func (a *Assetvul) DeleteByVulId(ctx context.Context, id []int) error {
	return mysql.FromContext(ctx).Model(&Assetvul{}).Where("id in ?", id).Delete(&Assetvul{}).Error
}

// DeleteByAssetId 通过资产ID删除漏洞
func (a *Assetvul) DeleteByAssetId(ctx context.Context, id []int) error {
	return mysql.FromContext(ctx).Model(&Assetvul{}).Where("asset_id in ?", id).Delete(&Assetvul{}).Error
}

func (a *Assetvul) AllByAssetId(ctx context.Context, assetId int) []Assetvul {
	var list []Assetvul
	mysql.FromContext(ctx).Model(&Assetvul{}).Where("asset_id = ?", assetId).Find(&list)
	return list
}

func (a *Assetvul) UpdateIsRepeatById(ctx context.Context, id, isRepeat int) error {
	return mysql.FromContext(ctx).Model(&Assetvul{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_replace":  isRepeat,
		"update_time": time.Now(),
	}).Error
}

// AllAssetVulByAssetIds 通过资产ID获取所有资产漏洞信息
func (a *Assetvul) AllAssetVulByAssetIds(ctx context.Context, assetIds []int) []Assetvul {
	var (
		list []Assetvul
		db   = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Where("asset_id in ?", assetIds).Find(&list)
	return list
}

// AllAssetVulList 获取所有资产漏洞信息
func (a *Assetvul) AllAssetVulList(ctx context.Context, assetIds []int, search string, page, limit int) ([]Assetvul, int64) {
	var (
		list  []Assetvul
		count int64
		db    = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Where("asset_id in ?", assetIds)
	db.Count(&count)
	if search != "" {
		db = db.Where("name like ?", "%"+search+"%").Or("location LIKE ?", "%"+search+"%")
	}
	db.Limit(limit).Offset(limit * (page - 1)).Find(&list)
	return list, count
}

// GetAssetVulByAssetIds 通过资产ID获取资产漏洞信息
func (a *Assetvul) GetAssetVulByAssetIds(ctx context.Context, assetIds []int, page, limit int) ([]Assetvul, int64) {
	var (
		list  []Assetvul
		db    = mysql.FromContext(ctx).Model(&Assetvul{})
		count int64
	)
	db.Where("asset_id in ?", assetIds)
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&list)
	return list, count
}

func (a *Assetvul) UpdateAssetVulStatus(ctx context.Context, id, status int) error {
	return mysql.FromContext(ctx).Model(&Assetvul{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}).Error
}

func (a *Assetvul) GetAssetVulTypeStatByIds(ctx context.Context, assetIds []int) []AssetVulTypeStat {
	var (
		assetVulTypeStatList []AssetVulTypeStat
		db                   = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Where("asset_id in ?", assetIds).Select("type,risk, Count(id) as count").Group("type,risk").Find(&assetVulTypeStatList)
	return assetVulTypeStatList
}

// GetAssetTargetURLByName 根据漏洞名称模糊查询关联的资产地址
func (a *Assetvul) GetAssetTargetURLByName(ctx context.Context, name string) []Assetvul {
	var (
		assetvulList []Assetvul
		db           = mysql.FromContext(ctx).Model(&Assetvul{})
	)
	db.Select("target_url").Where("name like ?", "%"+name+"%").Distinct("target_url").Find(&assetvulList)
	return assetvulList
}
