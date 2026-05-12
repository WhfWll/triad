package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"time"
)

// AssetRiskTrend 资产风险趋势
type AssetRiskTrend struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	AssetID    int       `gorm:"column:asset_id" json:"assetID"`       // 资产ID
	IP         string    `gorm:"column:ip" json:"ip"`                  // 资产ip或域名
	RiskLevel  int       `gorm:"column:risk_level" json:"riskLevel"`   // 风险等级，1&2-高危,3-中危,4-低危,0-安全
	RiskNum    string    `gorm:"column:risk_num" json:"riskNum"`       // 记录资产最近一次任务的风险漏洞情况
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	TargetID   int       `gorm:"column:target_id" json:"targetID"`     // 渗透目标id
	TaskID     int       `gorm:"column:task_id" json:"taskID"`         // 任务ID
}

func (art *AssetRiskTrend) TableName() string {
	return "asset_risk_trend"
}

// AddAssetRiskTrendInfo 批量插入资产风险趋势数据
func (art *AssetRiskTrend) AddAssetRiskTrendInfo(ctx context.Context, data []AssetRiskTrend) error {
	db := mysql.FromContext(ctx).Model(&AssetRiskTrend{})
	if len(data) == 0 {
		return nil
	}
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (art *AssetRiskTrend) DeleteByAssetIDs(ctx context.Context, AssetIDs []int) error {
	return mysql.FromContext(ctx).Model(&AssetRiskTrend{}).Where("asset_id in ?", AssetIDs).Delete(&AssetRiskTrend{}).Error
}

// UpdateAssetRiskTrendInfo 根据 ID 更新资产风险趋势数据
func (art *AssetRiskTrend) UpdateAssetRiskTrendInfo(ctx context.Context, id int, updates map[string]interface{}) error {
	db := mysql.FromContext(ctx).Model(&AssetRiskTrend{})
	if err := db.Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// GetAllRiskAssetTrendList 查询所有风险资产趋势数据
func (art *AssetRiskTrend) GetAllRiskAssetTrendList(ctx context.Context, startTime, endTime string) ([]AssetRiskTrend, int64) {
	var (
		assetRiskTrendList []AssetRiskTrend
		db                 = mysql.FromContext(ctx).Model(&AssetRiskTrend{})
		count              int64
		query              string
		args               []interface{}
	)
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
	db.Order("create_time desc").Where(query, args...).Find(&assetRiskTrendList)
	return assetRiskTrendList, count
}

// GetLatestRiskAssetTrendListByIP 查询时间段内每个IP最新的资产风险趋势记录
func (art *AssetRiskTrend) GetLatestRiskAssetTrendListByIP(ctx context.Context, startTime, endTime string) ([]AssetRiskTrend, int64, error) {
	var (
		result []AssetRiskTrend
		db     = mysql.FromContext(ctx)
		count  int64
	)

	// 子查询只保留 ip 和 max_time，避免 ONLY_FULL_GROUP_BY 报错
	subQuery := db.Model(&AssetRiskTrend{}).
		Select("ip, MAX(create_time) AS max_time").
		Where("risk_level != 0").
		Group("ip")

	if startTime != "" {
		subQuery = subQuery.Where("create_time >= ?", startTime)
	}
	if endTime != "" {
		subQuery = subQuery.Where("create_time <= ?", endTime)
	}

	// 主查询 join 回去查完整字段
	mainQuery := db.Model(&AssetRiskTrend{}).
		Joins("JOIN (?) AS latest ON asset_risk_trend.ip = latest.ip AND asset_risk_trend.create_time = latest.max_time", subQuery).
		Order("asset_risk_trend.create_time DESC")

	// 统计数量
	if err := mainQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := mainQuery.Find(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, count, nil
}
