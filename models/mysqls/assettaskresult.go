package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

type Assettaskresult struct {
	ID           int       `gorm:"column:id;primary_key" json:"id"`           // 主键
	TaskResultID int       `gorm:"column:task_result_id" json:"taskResultID"` // 所属任务结果id
	AssetID      int       `gorm:"column:asset_id" json:"assetID"`            // 所属资产id
	SubObjType   string    `gorm:"column:sub_obj_type" json:"subObjType"`     // 数据子类型
	ObjID        string    `gorm:"column:obj_id" json:"objID"`                // 数据对象id
	SubObjID     string    `gorm:"column:sub_obj_id" json:"subObjID"`         // 数据子对象id
	Identify     string    `gorm:"column:identify" json:"identify"`           // 数据标识符,用来修改时方便定位数据的
	Field1       string    `gorm:"column:field1" json:"field1"`               // 筛选字段1
	Field2       string    `gorm:"column:field2" json:"field2"`               // 筛选字段2
	Field3       string    `gorm:"column:field3" json:"field3"`               // 筛选字段3
	Field4       string    `gorm:"column:field4" json:"field4"`               // 筛选字段4
	JSONResult   string    `gorm:"column:json_result" json:"jsonresult"`      // 各数据类型的json格式结果
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`      // 创建时间
}

// TableName sets insert table name for this struct type
func (a *Assettaskresult) TableName() string {
	return "asset_task_result"
}

// Get retrieves a list of assettaskresult from database
func (a *Assettaskresult) GetAssettaskresultList(ctx context.Context, page, limit int) ([]Assettaskresult, int64, error) {
	var (
		assettaskresultList []Assettaskresult
		count               int64
		db                  = mysql.FromContext(ctx).Model(&Assettaskresult{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&assettaskresultList)
	db.Count(&count)

	return assettaskresultList, count, nil
}

// 根据端口服务和组件查询并返回所属资产id
func (a *Assettaskresult) GetAssetIdsByPortServiceComponent(ctx context.Context, port, service, component string, assetIds []int) []Assettaskresult {
	var (
		assettaskresultList []Assettaskresult
		db                  = mysql.FromContext(ctx).Model(&Assettaskresult{})
	)
	if len(port) > 0 {
		db = db.Where("field2 = ?", port)
	}
	if len(service) > 0 {
		db = db.Where("field3 = ?", service)
	}
	if len(component) > 0 {
		db = db.Where("field4 like ?", "%"+component+"%")
	}
	if len(assetIds) > 0 {
		db = db.Where("asset_id in ?", assetIds)
	}
	db.Select("asset_id").Distinct("asset_id").Find(&assettaskresultList)
	return assettaskresultList
}

// Get retrieves a single record of assettaskresult from database
func (a *Assettaskresult) GetAssettaskresult(ctx context.Context) (Assettaskresult, error) {
	var (
		assettaskresult Assettaskresult
		err             error
		db              = mysql.FromContext(ctx).Model(&Assettaskresult{})
	)

	curErr := db.Where("id = ?", a.ID).First(&assettaskresult).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assettaskresult, err
}

// Add persists assettaskresult to database
func (a *Assettaskresult) AddAssettaskresult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assettaskresult{})

	if err := db.Create(a).Error; err != nil {
		return err
	}

	return nil
}

// Update changes assettaskresult by id
func (a *Assettaskresult) UpdateAssettaskresult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assettaskresult{})
	if err := db.Where("id = ?", a.ID).Updates(a).Error; err != nil {
		return err
	}
	return nil
}

// 服务类型统计
type AssetServiceTypeStat struct {
	ServiceType string `gorm:"column:service_type" json:"service_type"`
	Count       int    `gorm:"column:count" json:"count"`
}

func (a *Assettaskresult) GetAssetServiceTypeStat(ctx context.Context) []AssetServiceTypeStat {
	var (
		assetVulTypeStatList []AssetServiceTypeStat
		db                   = mysql.FromContext(ctx).Model(&Assettaskresult{})
	)
	db.Select("field3 as service_type , Count(id) as count").Where("sub_obj_type = ?", enums.AssetTaskResultSubObjTypeService).Group("field3").Find(&assetVulTypeStatList)
	return assetVulTypeStatList
}

// 批量添加资产信息收集数据
func (a *Assettaskresult) MultipartInsert(ctx context.Context, list *[]Assettaskresult) error {
	var db = mysql.FromContext(ctx).Model(&Assettaskresult{})
	if err := db.Create(list).Error; err != nil {
		return err
	}
	return nil
}

func (a *Assettaskresult) DeleteById(ctx context.Context, id []int) error {
	return mysql.FromContext(ctx).Model(&Assettaskresult{}).Where("asset_id in ?", id).Delete(&Assettaskresult{}).Error
}

func (a *Assettaskresult) AllByAssetId(ctx context.Context, assetId int) []Assettaskresult {
	var list []Assettaskresult
	mysql.FromContext(ctx).Model(&Assettaskresult{}).Where("asset_id = ?", assetId).Find(&list)
	return list
}

func (a *Assettaskresult) UpdatePortServiceComponentById(ctx context.Context, id int, port, service, component string) error {
	return mysql.FromContext(ctx).Model(&Assettaskresult{}).Where("id = ?", id).Updates(map[string]interface{}{
		"field2":      port,
		"field3":      service,
		"field4":      component,
		"create_time": time.Now(),
	}).Error
}

// GetAllTaskPortResultByAssetId 获取所有资产任务返回 端口列表 默认子对象类型 1_1 为端口列表
func (a *Assettaskresult) GetAllTaskPortResultByAssetId(ctx context.Context, assetId int) []Assettaskresult {
	var list []Assettaskresult
	mysql.FromContext(ctx).Model(&Assettaskresult{}).Where("asset_id = ? and sub_obj_type = '1_1'", assetId).Find(&list)
	return list
}

// GetTaskPortResultByAssetId 获取资产任务返回 端口列表 默认子对象类型 1_1 为端口列表
func (a *Assettaskresult) GetTaskPortResultByAssetId(ctx context.Context, assetId, page, limit int) []Assettaskresult {
	var list []Assettaskresult
	mysql.FromContext(ctx).Model(&Assettaskresult{}).Where("asset_id = ? and sub_obj_type = '1_1'", assetId).Limit(limit).Offset(limit * (page - 1)).Find(&list)
	return list
}

// Get retrieves a single record of assettaskresult from database
func (a *Assettaskresult) GetAssettaskresultByAssetId(ctx context.Context, assetId int) ([]Assettaskresult, error) {
	var (
		assettaskresultList []Assettaskresult
		err                 error
		db                  = mysql.FromContext(ctx).Model(&Assettaskresult{})
	)
	curErr := db.Where("asset_id = ?", assetId).Find(&assettaskresultList).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assettaskresultList, err
}

// Update changes assettaskresult by id
func (a *Assettaskresult) UpdateAssettaskresultById(ctx context.Context, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Assettaskresult{})
	if err := db.Where("id = ?", a.ID).Updates(param).Error; err != nil {
		return err
	}
	return nil
}
