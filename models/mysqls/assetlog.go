package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

type Assetlog struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Type       int       `gorm:"column:type" json:"type"`              // 日志类型，1-创建资产渗透
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
}

// TableName sets insert table name for this struct type
func (a *Assetlog) TableName() string {
	return "asset_log"
}

// Get retrieves a list of assetlog from database
func (a *Assetlog) GetAssetlogList(ctx context.Context, page, limit int) ([]Assetlog, int64, error) {
	var (
		assetlogList []Assetlog
		count        int64
		db           = mysql.FromContext(ctx).Model(&Assetlog{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&assetlogList)
	db.Count(&count)

	return assetlogList, count, nil
}

// Get retrieves a single record of assetlog from database
func (a *Assetlog) GetAssetlog(ctx context.Context) (Assetlog, error) {
	var (
		assetlog Assetlog
		err      error
		db       = mysql.FromContext(ctx).Model(&Assetlog{})
	)

	curErr := db.Where("id = ?", a.ID).First(&assetlog).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return assetlog, err
}

// Add persists assetlog to database
func (a *Assetlog) AddAssetlog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetlog{})

	if err := db.Create(a).Error; err != nil {
		return err
	}

	return nil
}

// Update changes assetlog by id
func (a *Assetlog) UpdateAssetlog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetlog{})

	if err := db.Where("id = ?", a.ID).Updates(a).Error; err != nil {
		return err
	}

	return nil
}

//资产渗透测试趋势
type AssetTaskStat struct {
	Date  string `gorm:"column:date" json:"date"`
	Count int64  `gorm:"column:count" json:"count"`
}

func (a *Assetlog) GetAssetTaskstat(ctx context.Context, startTime, dateFormat string) []AssetTaskStat {
	var (
		assetTaskStatRes []AssetTaskStat
		db               = mysql.FromContext(ctx).Model(&Assetlog{})
	)
	db.Where("type = ?", enums.AssetLogTypeOne)
	if startTime != "" {
		db = db.Where("create_time > ?", startTime)
	}
	db.Select("DATE_FORMAT(create_time, '" + dateFormat + "') as date, COUNT(id) as count").Group("DATE_FORMAT(create_time, '" + dateFormat + "')").Find(&assetTaskStatRes)
	return assetTaskStatRes
}
