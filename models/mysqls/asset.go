package mysqls

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"math"
	"smart/tools/enums"
	"strconv"
	"time"
)

// 资产其他信息,JSON格式，硬件/主机名/资产类型/虚拟资产/归属地/责任人/责任人邮箱
type Asset struct {
	ID                    int       `gorm:"column:id;primary_key" json:"id"`                            // 主键
	AssetGroupID          int       `gorm:"column:asset_group_id" json:"assetGroupID"`                  // 所属资产组，1表示未分组
	IP                    string    `gorm:"column:ip" json:"ip"`                                        // 资产ip或域名
	IPSegment             string    `gorm:"column:ip_segment" json:"ipSegment"`                         // ip段
	AssetType             int       `gorm:"column:asset_type" json:"assetType"`                         // 资产类型
	IpNum                 int64     `gorm:"column:ip_num" json:"ipNum"`                                 // ip值
	Name                  string    `gorm:"column:name" json:"name"`                                    // 资产名称
	OperateSystem         string    `gorm:"column:operate_system" json:"operateSystem"`                 // 操作系统
	Location              string    `gorm:"column:location" json:"location"`                            // 地理位置
	RiskLevel             int       `gorm:"column:risk_level" json:"riskLevel"`                         // 风险等级，1&2-高危,3-中危,4-低危,0-安全
	RiskNum               string    `gorm:"column:risk_num" json:"riskNum"`                             // 记录资产最近一次任务的风险漏洞情况
	BusinessSystem        string    `gorm:"column:business_system" json:"businessSystem"`               // 业务系统
	ResponsibleDepartment string    `gorm:"column:responsible_department" json:"responsibleDepartment"` // 责任部门
	FilingLevel           int       `gorm:"column:filing_level" json:"filingLevel"`                     // 备案等级，1-无,2-等保一级,3-等保二级,4-等保三级,5-等保四级,6-等保五级
	DeviceWeight          int       `gorm:"column:device_weight" json:"deviceWeight"`                   // 设备权重 1高 2中 3低 4极高 5极低
	TrustLevel            int       `gorm:"column:trust_level" json:"trustLevel"`                       // 可信权重 1 可信 2 未登记
	IsCloudHost           int       `gorm:"column:is_cloud_host" json:"isCloudHost"`                    // 是否是云主机
	Tags                  string    `gorm:"column:tags" json:"tags"`                                    // 标签，多个用英文逗号隔开
	Info                  string    `gorm:"column:info" json:"info"`                                    // 资产其他信息,JSON格式，硬件/主机名/资产类型/虚拟资产/归属地/责任人/责任人邮箱
	IsIgnore              int       `gorm:"column:is_ignore" json:"isIgnore"`                           // 最新资产排序是否忽略，1-不忽略，2-忽略
	AssetChangesType      int       `gorm:"column:asset_changes_type" json:"assetChangesType"`          // 资产变化类型:1-未变化,2-已减少IP，3-新增加IP，4-端口变化IP，5-服务变化IP，6-组件变化IP
	Islive                int       `gorm:"column:islive" json:"islive"`                                // 存活状态:1-存活,2-不存活IP
	CreateTime            time.Time `gorm:"column:create_time" json:"createTime"`                       // 创建时间
	UpdateTime            time.Time `gorm:"column:update_time" json:"updateTime"`                       // 修改时间
	TargetIds             string    `gorm:"column:target_ids" json:"targetIds"`                         // 渗透目标id,多个用英文逗号隔开
	Status                int       `gorm:"column:status" json:"status"`                                // 状态，1-待同步、2-同步中、3-已完成
}

func (a *Asset) TableName() string {
	return "asset"
}

// 根据ip/操作系统/风险等级/标签/资产名称/业务系统/责任部门/备案等级查询资产
func (a *Asset) GetAssetByIpOperateSystemRiskLevelTags(ctx context.Context, search, ip, operateSystem string, riskLevel int, tags, assetName, businessSystem, responsibleDepartment string, filingLevel int, assetIds []int) []Asset {
	var (
		assetList []Asset
		db        = mysql.FromContext(ctx).Model(&Asset{})
	)
	if len(ip) > 0 {
		db = db.Where("ip = ?", ip)
	}
	if len(operateSystem) > 0 {
		db = db.Where("operate_system = ?", operateSystem)
	}
	if riskLevel > 0 {
		db = db.Where("risk_level = ?", riskLevel)
	}
	if len(businessSystem) > 0 {
		db = db.Where("business_system = ?", businessSystem)
	}
	if len(responsibleDepartment) > 0 {
		db = db.Where("responsible_department = ?", responsibleDepartment)
	}
	if filingLevel > 0 {
		db = db.Where("filing_level = ?", filingLevel)
	}
	if len(assetIds) > 0 {
		db = db.Where("id in ?", assetIds)
	}
	if len(tags) > 0 {
		db = db.Where("tags like ?", "%"+tags+"%")
	}
	if len(assetName) > 0 {
		db = db.Where("name like ?", "%"+assetName+"%")
	}
	if len(search) > 0 {
		db = db.Where("ip like ? or name like ?", "%"+search+"%", "%"+search+"%")
	}
	db.Find(&assetList)
	return assetList
}

// 查询资产总数量
func (a *Asset) GetAssetTotal(ctx context.Context) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Count(&count)
	return count
}

// 按照资产危险等级计算数量
// SELECT risk_level,COUNT(id) as total FROM `asset` GROUP BY risk_level;
type GetAssetTotalByRiskLevel struct {
	RiskLevel int `gorm:"column:risk_level" json:"risk_level"`
	Total     int `gorm:"column:total" json:"total"`
}

func (a *Asset) GetAssetTotalByRiskLevel(ctx context.Context) []GetAssetTotalByRiskLevel {
	var (
		result []GetAssetTotalByRiskLevel
		db     = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Select("risk_level,COUNT(id) as total").Group("risk_level").Find(&result)
	return result
}

// Get retrieves a single record of asset from database
func (a *Asset) GetAsset(ctx context.Context) (Asset, error) {
	var (
		asset Asset
		err   error
		db    = mysql.FromContext(ctx).Model(&Asset{})
	)

	curErr := db.Where("id = ?", a.ID).First(&asset).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return asset, err
}

// 新增一条数据
func (a *Asset) AddAsset(ctx context.Context) (int, error) {
	db := mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Create(a).Error; err != nil {
		return 0, err
	}
	return a.ID, nil
}

// 根据id修改资产信息
func (a *Asset) UpdateAssetById(ctx context.Context, id int, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Where("id = ?", id).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

func (a *Asset) UpdateAssetByIds(ctx context.Context, id any, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Where("id in ?", id).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// 根据ip修改数据
func (a *Asset) UpdateAssetByIp(ctx context.Context, ip string, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Where("ip = ?", ip).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// 查询所有资产组数据
func (a *Asset) GetAllAsset(ctx context.Context, limit, offset int) []Asset {
	var (
		assetList []Asset
		db        = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Limit(limit).Offset(offset).Order("id desc").Find(&assetList)
	return assetList
}

// 依据IP分组 并 返回分组后的IP
func (a *Asset) AllIpByGroup(ctx context.Context) (returnData []string) {
	var asset []struct {
		IP string `gorm:"column:ip" json:"ip"` // 资产ip或域名
	}
	var db = mysql.FromContext(ctx).Model(&Asset{})
	db.Select("ip").Group("ip").Find(&asset)
	for _, item := range asset {
		returnData = append(returnData, item.IP)
	}
	return
}

// 批量添加资产数据
func (a *Asset) MultipartInsert(ctx context.Context, list *[]Asset) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Create(list).Error; err != nil {
		return err
	}
	return nil
}

func (a *Asset) AllByIp(ctx context.Context, ip []string) []Asset {
	var (
		list []Asset
		db   = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Where("ip in ?", ip).Find(&list)
	return list
}

func (a *Asset) DeleteById(ctx context.Context, id []int) error {
	return mysql.FromContext(ctx).Model(&Asset{}).Where("id in ?", id).Delete(&Asset{}).Error
}

func (a *Asset) GetIPsByIds(ctx context.Context, ids []int) ([]string, error) {
	var ips []string
	err := mysql.FromContext(ctx).
		Model(&Asset{}).
		Where("id IN ?", ids).
		Pluck("ip", &ips).
		Error
	return ips, err
}

func (a *Asset) AllAssetByGroupIds(ctx context.Context, groupIds []int) []Asset {
	var (
		list []Asset
		db   = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Where("asset_group_id in ?", groupIds).Find(&list)
	return list
}

func (a *Asset) AllByStatus(ctx context.Context, status int) []Asset {
	var (
		list []Asset
		db   = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Where("status = ?", status).Find(&list)
	return list
}

func (a *Asset) AllAssetByGroupIdsAndSearch(ctx context.Context, groupIds []int, assetIP, systemOp, tags, assetType, isCloudHost, domain string, assetRisk *int, fillingLevel, page, size int, ips []string) ([]Asset, int64) {
	var (
		list  []Asset
		db    = mysql.FromContext(ctx).Model(&Asset{})
		count int64
		query string
		args  []interface{}
	)
	db = db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveYes) + " and ip !='' and asset_changes_type != " + strconv.Itoa(enums.AssetChangeTypeReduce))
	query += "1 = ?"
	args = append(args, 1)
	if assetIP != "" {
		query += " and ip LIKE ?"
		args = append(args, "%"+assetIP+"%")
	}
	if domain != "" {
		query += " and ip LIKE ?"
		args = append(args, "%"+domain+"%")
	}
	if systemOp != "" {
		query += " AND operate_system LIKE ?"
		args = append(args, "%"+systemOp+"%")
	}
	if tags != "" {
		query += " AND tags LIKE ?"
		args = append(args, "%"+tags+"%")
	}
	if assetRisk != nil {
		query += " and risk_level = ?"
		args = append(args, *assetRisk)
	}
	if assetType != "" {
		query += " and asset_type = ?"
		args = append(args, assetType)
	}
	if fillingLevel != 0 {
		query += " and filing_level = ?"
		args = append(args, fillingLevel)
	}
	if len(ips) != 0 {
		query += " and ip IN (?)"
		args = append(args, ips)
	}
	if isCloudHost != "" {
		query += " and is_cloud_host = ?"
		args = append(args, isCloudHost)
	}
	db.Where("asset_group_id in ?", groupIds).Where(query, args...).Count(&count)
	db.Where("asset_group_id in ?", groupIds).Where(query, args...)
	db.Limit(size).Offset(size * (page - 1)).Order("update_time desc").Find(&list)
	return list, count
}

// GetAssetByIps 通过ip列表获取资产
func (a *Asset) GetAssetByIps(ctx context.Context, ips []string) ([]Asset, error) {
	var (
		assets []Asset
		err    error
		db     = mysql.FromContext(ctx).Model(&Asset{})
	)
	curErr := db.Where("ip in ?", ips).Find(&assets).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assets, err
}

func (a *Asset) AddAll(ctx context.Context, datas []Asset) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.CreateInBatches(datas, 1000).Error; err != nil {
		return err
	}
	return nil
}

// GetAssetByIds 通过id列表获取资产
func (a *Asset) GetAssetByIds(ctx context.Context, ids []int) ([]Asset, error) {
	var (
		assets []Asset
		err    error
		db     = mysql.FromContext(ctx).Model(&Asset{})
	)
	curErr := db.Where("id in ?", ids).Find(&assets).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assets, err
}

// 根据忽略状态统计资产变化情况
// SELECT asset_changes_type,COUNT(id) as total FROM `asset` WHERE is_ignore=1  GROUP BY asset_changes_type;
type AssetChangeType struct {
	AssetChangesType int   `gorm:"column:asset_changes_type" json:"assetChangesType"`
	Total            int64 `gorm:"column:total" json:"total"`
}

func (a *Asset) CountAssetChangeTypeByIsIgnore(ctx context.Context, isIgnore int) []AssetChangeType {
	var (
		result []AssetChangeType
		db     = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Select("asset_changes_type,COUNT(id) as total").Where("is_ignore = ?", isIgnore).Group("asset_changes_type").Find(&result)
	return result
}

// 资产变化列表查询
func (a *Asset) GetChangeAssetList(ctx context.Context, isIgnore int, search string, assetChangesType, isLive int, updateTime string, page, limit int) ([]Asset, int64) {
	var (
		list  []Asset
		count int64
		order = "desc"
		db    = mysql.FromContext(ctx).Model(&Asset{})
	)
	if isIgnore != 0 {
		db = db.Where("is_ignore = ?", isIgnore)
	}
	if assetChangesType != 0 {
		db = db.Where("asset_changes_type = ?", assetChangesType)
	}
	if isLive != 0 {
		db = db.Where("islive = ?", isLive)
	}
	if len(search) > 0 {
		db = db.Where("ip like ? or operate_system like ?", "%"+search+"%", "%"+search+"%")
	}
	if len(updateTime) > 0 {
		order = updateTime
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("update_time " + order).Find(&list)
	return list, count
}

func (a *Asset) GetChangeAssetAll(ctx context.Context, isIgnore int, search string, assetChangesType, isLive int) []Asset {
	var (
		list []Asset
		db   = mysql.FromContext(ctx).Model(&Asset{})
	)
	if isIgnore != 0 {
		db = db.Where("is_ignore = ?", isIgnore)
	}
	if assetChangesType != 0 {
		db = db.Where("asset_changes_type = ?", assetChangesType)
	}
	if isLive != 0 {
		db = db.Where("islive = ?", isLive)
	}
	if len(search) > 0 {
		db = db.Where("ip like ? or operate_system like ?", "%"+search+"%", "%"+search+"%")
	}
	db.Find(&list)
	return list
}

// GetAssetScanIPCountByActivityID 依据活动获取IP资产数量
func (a *Asset) GetAssetScanIPCountByActivityID(ctx context.Context, activityID int) (int64, error) {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Asset{})
	)
	db = db.Where("activity_id = ?", activityID)
	db.Count(&count)
	return count, nil
}

// UpdateAssetGroupIDByIds 批量更新资产的资产组ID为默认资产组ID
func (a *Asset) UpdateAssetGroupIDByIds(ctx context.Context, ids []int, assetGroupID int) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Where("id in ?", ids).Updates(map[string]interface{}{
		"asset_group_id": assetGroupID,
	}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateAssetByAssetID 根据id修改资产信息
func (a *Asset) UpdateAssetByAssetID(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Where("id = ?", a.ID).Updates(a).Error; err != nil {
		return err
	}
	if a.IsCloudHost == 0 {
		if err := db.Where("id = ?", a.ID).Updates(map[string]interface{}{
			"is_cloud_host": a.IsCloudHost,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpdateAssetByAssetIP 根据IP修改资产
func (a *Asset) UpdateAssetByAssetIP(ctx context.Context, ip, riskNum, opSys string, riskLevel int) (int, error) {
	var (
		db    = mysql.FromContext(ctx).Model(&Asset{})
		asset Asset
	)
	err := db.Where("ip = ?", ip).First(&asset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	if err := db.Where("ip = ?", ip).Updates(map[string]interface{}{
		"risk_level":         riskLevel,
		"risk_num":           riskNum,
		"update_time":        time.Now(),
		"islive":             enums.AssetIsLiveYes,
		"asset_changes_type": enums.AssetChangeTypeAdd,
		"operate_system":     opSys,
	}).Error; err != nil {
		return 0, err
	}
	return asset.ID, nil
}

// UpsertAssetByAssetIP 根据IP修改资产信息，不存在则创建
func (a *Asset) UpsertAssetByAssetIP(ctx context.Context, ip, riskNum, opSys string, riskLevel int) (int, error) {
	var asset Asset
	err := mysql.FromContext(ctx).Model(&Asset{}).Where("ip = ?", ip).First(&asset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[INFO] UpsertAssetByAssetIP - asset not found, create new | ip=%s | riskLevel=%d | riskNum=%s | opSys=%s",
				ip, riskLevel, riskNum, opSys)
			newAsset := Asset{
				IP:               ip,
				RiskLevel:        riskLevel,
				AssetGroupID:     enums.DefaultAssetGroup,
				RiskNum:          riskNum,
				OperateSystem:    opSys,
				Islive:           enums.AssetIsLiveYes,
				AssetChangesType: enums.AssetChangeTypeAdd,
				CreateTime:       time.Now(),
				UpdateTime:       time.Now(),
			}
			if err := mysql.FromContext(ctx).Model(&Asset{}).Create(&newAsset).Error; err != nil {
				log.Printf("[ERROR] UpsertAssetByAssetIP - failed to create asset | ip=%s | err=%v", ip, err)
				return 0, err
			}
			log.Printf("[INFO] UpsertAssetByAssetIP - asset created successfully | ip=%s | assetID=%d", ip, newAsset.ID)
			return newAsset.ID, nil
		}
		log.Printf("[ERROR] UpsertAssetByAssetIP - query asset failed | ip=%s | err=%v", ip, err)
		return 0, err
	}
	// 存在则更新
	if opSys == "" {
		opSys = asset.OperateSystem
	}
	if err := mysql.FromContext(ctx).Model(&Asset{}).Where("ip = ?", ip).Updates(map[string]interface{}{
		"risk_level":         riskLevel,
		"risk_num":           riskNum,
		"update_time":        time.Now(),
		"islive":             enums.AssetIsLiveYes,
		"asset_changes_type": enums.AssetChangeTypeAdd,
		"operate_system":     opSys,
	}).Error; err != nil {
		log.Printf("[ERROR] UpsertAssetByAssetIP - failed to update asset | ip=%s | err=%v", ip, err)
		return 0, err
	}
	log.Printf("[INFO] UpsertAssetByAssetIP - asset updated successfully | ip=%s | assetID=%d", ip, asset.ID)
	return asset.ID, nil
}

// BathAddAsset 批量新增
func (i *Asset) BathAddAsset(ctx context.Context, assets []Asset) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	// 批量入库做限制 避免Prepared statement contains too many placeholders问题出现[一条语句占位符最多为65535] 默认3000条分页录入
	if len(assets) > 3000 {
		pageSize := 3000                                                 //每次处理的数量
		page := int(math.Ceil(float64(len(assets)) / float64(pageSize))) //分页数
		for num := 1; num <= page; num++ {
			sliceStart := (num - 1) * pageSize
			sliceEnd := sliceStart + pageSize
			if num == page {
				sliceEnd = len(assets)
			}
			newIps := assets[sliceStart:sliceEnd]
			if err := db.Create(&newIps).Error; err != nil {
				return err
			}
		}
	} else {
		if err := db.Create(&assets).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllAssetList 查询所有资产组数据
func (a *Asset) GetAllAssetList(ctx context.Context, startTime, endTime string) ([]Asset, int64) {
	var (
		assetList []Asset
		db        = mysql.FromContext(ctx).Model(&Asset{})
		count     int64
		query     string
		args      []interface{}
	)
	db = db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveYes) + " and ip !='' and asset_changes_type != " + strconv.Itoa(enums.AssetChangeTypeReduce))
	query += "1 = 1"
	if startTime != "" {
		query += " and create_time >= ? "
		args = append(args, startTime)
	}
	if endTime != "" {
		query += " and create_time <= ?"
		args = append(args, endTime)
	}
	db.Where(query, args...).Count(&count)
	db.Order("id desc").Where(query, args...).Find(&assetList)
	return assetList, count
}

// CalNewIPChange 计算新增Ip变化
func (a *Asset) CalNewIPChange(ctx context.Context, startTime, endTime, changeType string) int64 {
	var (
		db    = mysql.FromContext(ctx).Model(&Asset{})
		count int64
		query string
		args  []interface{}
	)
	if changeType == "add" {
		db = db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveYes) + " and asset_changes_type = " + strconv.Itoa(enums.AssetChangeTypeAdd))
	} else if changeType == "reduce" {
		db = db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveNo) + " and asset_changes_type = " + strconv.Itoa(enums.AssetChangeTypeReduce))
	}
	query += "1 = 1"
	if startTime != "" {
		query += " and create_time >= ? "
		args = append(args, startTime)
	}
	if endTime != "" {
		query += " and create_time <= ?"
		args = append(args, endTime)
	}
	db.Where(query, args...).Count(&count)
	return count
}

func (a *Asset) AllAssetByGroupIdsAndSearchNoPage(ctx context.Context, groupIds []int, assetIP, port, service, systemOp, vulName, domain, finger, tags string, assetRisk, assetType, fillingLevel int) ([]Asset, int64) {
	var (
		list  []Asset
		db    = mysql.FromContext(ctx).Model(&Asset{})
		count int64
	)
	db = db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveYes) + " and asset_changes_type != " + strconv.Itoa(enums.AssetChangeTypeReduce))
	if assetIP != "" {
		db = db.Where("ip = ?", assetIP)
	}
	if systemOp != "" {
		db = db.Where("operate_system like ?", "%"+systemOp+"%")
	}
	if tags != "" {
		db = db.Where("tags like ?", "%"+tags+"%")
	}
	if assetRisk != 0 {
		db = db.Where("risk_level = ?", assetRisk)
	}
	if assetType != 0 {
		db = db.Where("asset_type = ?", assetType)
	}
	if fillingLevel != 0 {
		db = db.Where("filing_level = ?", fillingLevel)
	}
	db = db.Where("asset_group_id in ?", groupIds)
	db.Count(&count)
	db.Order("id desc").Find(&list)
	return list, count
}

func (a *Asset) AllAssetByIdsAndSearchNoPage(ctx context.Context, groupIds []int, assetIP, port, service, systemOp, vulName, domain, finger, tags string, assetRisk, assetType, fillingLevel int) ([]Asset, int64) {
	var (
		list  []Asset
		db    = mysql.FromContext(ctx).Model(&Asset{})
		count int64
		query string
		args  []interface{}
	)
	db = db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveYes) + " and asset_changes_type != " + strconv.Itoa(enums.AssetChangeTypeReduce))
	query += "1 = ?"
	args = append(args, 1)
	if assetIP != "" {
		query += " and ip = ?"
		args = append(args, assetIP)
	}
	if systemOp != "" {
		query += " AND operate_system LIKE ?"
		args = append(args, "%"+systemOp+"%")
	}
	if tags != "" {
		query += " AND tags LIKE ?"
		args = append(args, "%"+tags+"%")
	}
	if assetRisk != 0 {
		query += " and risk_level = ?"
		args = append(args, assetRisk)
	}
	if assetType != 0 {
		query += " and asset_type = ?"
		args = append(args, assetType)
	}
	if fillingLevel != 0 {
		query += " and filing_level = ?"
		args = append(args, fillingLevel)
	}
	if len(groupIds) != 0 {
		query += " and id IN (?)"
		args = append(args, groupIds)
	}
	db.Where(query, args...).Count(&count)
	db.Where(query, args...).Order("id desc").Find(&list)
	return list, count
}

// GetNoLiveAssetByIp 查询不存活【软删除】IP
func (a *Asset) GetNoLiveAssetByIp(ctx context.Context, assetIP string) Asset {
	var (
		info Asset
		db   = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveNo) + " and asset_changes_type = " + strconv.Itoa(enums.AssetChangeTypeReduce) + " and ip = '" + assetIP + "'").First(&info)
	return info
}

// GetLiveAssetByIp 查询存活IP
func (a *Asset) GetLiveAssetByIp(ctx context.Context, assetIP string) Asset {
	var (
		info Asset
		db   = mysql.FromContext(ctx).Model(&Asset{})
	)
	db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveYes) + " and ip = '" + assetIP + "'").First(&info)
	return info
}

type IpNum struct {
	Value1   int64  `json:"value1"`
	Value2   int64  `json:"value2"`
	Relation string `json:"relation"`
}

func (a *Asset) GetAssetByIpNumRange(ctx context.Context, eqIpNum []IpNum, neqIpNum []IpNum) []Asset {
	var (
		assets []Asset
		db     = mysql.FromContext(ctx).Model(&Asset{})
		query  string
		args   []interface{}
	)
	if len(eqIpNum) == 0 {
		return assets
	}
	query += "1 = 1"
	for i := 0; i < len(eqIpNum); i++ {
		if i == 0 {
			if eqIpNum[i].Relation == "eq" {
				query += " and (ip_num = ?"
				args = append(args, eqIpNum[i].Value1)
			} else if eqIpNum[i].Relation == "in" {
				query += " and ((ip_num >= ? and ip_num <= ?)"
				args = append(args, eqIpNum[i].Value1, eqIpNum[i].Value2)
			}
		} else {
			if eqIpNum[i].Relation == "eq" {
				query += " or ip_num = ?"
				args = append(args, eqIpNum[i].Value1)
			} else if eqIpNum[i].Relation == "in" {
				query += " or (ip_num >= ? and ip_num <= ?)"
				args = append(args, eqIpNum[i].Value1, eqIpNum[i].Value2)
			}

		}
		if i == len(eqIpNum)-1 {
			query += ")"
		}
	}
	for i := 0; i < len(neqIpNum); i++ {
		if i == 0 {
			if neqIpNum[i].Relation == "neq" {
				query += " and (ip_num != ?"
				args = append(args, neqIpNum[i].Value1)
			} else if neqIpNum[i].Relation == "nin" {
				query += " and ((ip_num < ? or ip_num > ?)"
				args = append(args, neqIpNum[i].Value1, neqIpNum[i].Value2)
			}
		} else {
			if neqIpNum[i].Relation == "neq" {
				query += " and ip_num != ?"
				args = append(args, neqIpNum[i].Value1)
			} else if neqIpNum[i].Relation == "nin" {
				query += " and (ip_num < ? or ip_num > ?)"
				args = append(args, neqIpNum[i].Value1, neqIpNum[i].Value2)
			}
		}
		if i == len(neqIpNum)-1 {
			query += ")"
		}
	}
	db.Where(query, args...).Find(&assets)
	return assets
}

func (a *Asset) UpdateStatusDoneByIds(ctx context.Context, ids []int) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Where("id in ?", ids).Updates(map[string]interface{}{
		"status":     enums.AssetStatusDone,
		"target_ids": "",
	}).Error; err != nil {
		return err
	}
	return nil
}

// CountAllRiskAssets 查询所有风险资产信息数量
func (a *Asset) CountAllRiskAssets(ctx context.Context) (count int64) {
	db := mysql.FromContext(ctx).Model(&Asset{}).
		Where("islive = ? AND asset_changes_type != ? AND risk_level NOT IN ?", enums.AssetIsLiveYes, enums.AssetChangeTypeReduce, []int{0, 4})
	db.Count(&count)
	return count
}

// CountAllAssets 查询所有资产的数量（可选时间范围）
func (a *Asset) CountAllAssets(ctx context.Context) (count int64) {
	db := mysql.FromContext(ctx).Model(&Asset{}).
		Where("islive = ? AND asset_changes_type != ?", enums.AssetIsLiveYes, enums.AssetChangeTypeReduce)
	db.Count(&count)
	return count
}

// GetAssetByIP 通过ip获取资产信息
func (a *Asset) GetAssetByIP(ctx context.Context) (Asset, error) {
	var (
		asset Asset
		err   error
		db    = mysql.FromContext(ctx).Model(&Asset{})
	)

	curErr := db.Where("ip = ?", a.IP).First(&asset).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return asset, err
}

// UpdateAssetUpdateTimeInfo	更新资产信息
func (a *Asset) UpdateAssetUpdateTimeInfo(ctx context.Context, ids []int) error {
	var db = mysql.FromContext(ctx).Model(&Asset{})
	if err := db.Where("id in ?", ids).Updates(map[string]interface{}{
		"update_time": time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetAllRiskAsset 获取风险资产
func (a *Asset) GetAllRiskAsset(ctx context.Context) ([]Asset, int64) {
	var (
		assetList []Asset
		count     int64
		db        = mysql.FromContext(ctx).Model(&Asset{}).Select("id,ip,risk_level,create_time,update_time,risk_num")
	)
	db.Where("islive = " + strconv.Itoa(enums.AssetIsLiveYes) + " and ip !='' and asset_changes_type != " + strconv.Itoa(enums.AssetChangeTypeReduce)).Count(&count)
	db.Where("risk_level != 0 ").Order("id desc").Find(&assetList)
	return assetList, count
}
