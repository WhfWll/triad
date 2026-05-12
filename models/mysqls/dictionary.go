package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

type Dictionary struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Sources    int       `gorm:"column:sources" json:"sources"`        // 来源，1系统，2手动
	Service    int       `gorm:"column:service" json:"service"`        // 适用范围
	Name       string    `gorm:"column:name" json:"name"`              // 字典名称
	Types      int       `gorm:"column:types" json:"types"`            // 字典类型：1-用户字典，2-密码字典，3-wifi，4-路径字典，5-子域名字典
	IsDefault  int       `gorm:"column:is_default" json:"isDefault"`   // 是否是默认，1-是，2-否
	Content    string    `gorm:"column:content" json:"content"`        // 字典内容,换行符隔开
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

// TableName sets insert table name for this struct type
func (d *Dictionary) TableName() string {
	return "dictionary"
}

// DictionaryList Get retrieves a list of dictionary from database
func (d *Dictionary) DictionaryList(ctx context.Context, page, limit, types int, name string) ([]Dictionary, int64, error) {
	var (
		dictionaryList []Dictionary
		count          int64
		db             = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	if name != "" {
		db.Where(`name like ?`, "%"+name+"%")
	}
	if types != 0 {
		db.Where("types = ?", types)
	}

	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&dictionaryList)

	return dictionaryList, count, nil
}

// DictionaryRecord Get retrieves a single record of dictionary from database
func (d *Dictionary) DictionaryRecord(ctx context.Context, dictId int) (Dictionary, error) {
	var (
		dictionary Dictionary
		err        error
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	curErr := db.Where("id = ?", dictId).First(&dictionary).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return dictionary, err
}

// AddDictionary Add persists dictionary to database
func (d *Dictionary) AddDictionary(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Dictionary{})

	if err := db.Create(d).Error; err != nil {
		return err
	}

	return nil
}

// UpdateDictionary Update changes dictionary by id
func (d *Dictionary) UpdateDictionary(ctx context.Context, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Dictionary{})

	if err := db.Where("id = ?", param["id"]).Updates(param).Error; err != nil {
		return err
	}

	return nil
}

// DeleteDictionary Delete dictionary by id
func (d *Dictionary) DeleteDictionary(ctx context.Context, dictId int) error {
	var db = mysql.FromContext(ctx).Model(&Dictionary{})

	if err := db.Where("id = ?", dictId).Delete(d).Error; err != nil {
		return err
	}

	return nil
}

// Count 获取字典总数
func (d *Dictionary) Count(ctx context.Context) (int64, error) {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Dictionary{})
	)
	err := db.Count(&count).Error
	return count, err
}

// MultiDeleteDictionary multi delete by ids
// ids必须为[]int/[]string数组
func (d *Dictionary) MultiDeleteDictionary(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&Dictionary{})

	if err := db.Where("id in ?", ids).Delete(d).Error; err != nil {
		return err
	}

	return nil
}

// CancelDefaultDictionary 取消同类型和同服务的默认字典
func (d *Dictionary) CancelDefaultDictionary(ctx context.Context, types, service int) error {
	var db = mysql.FromContext(ctx).Model(&Dictionary{})

	if err := db.Where("service = ?", service).
		Where("types = ?", types).
		Update("is_default", enums.DictionaryDefaultNo).Error; err != nil {
		return err
	}

	return nil
}

// SetDefaultDictionary 根据id设置默认字典
func (d *Dictionary) SetDefaultDictionary(ctx context.Context, dictId int) error {
	var db = mysql.FromContext(ctx).Model(&Dictionary{})

	if err := db.Where("id = ?", dictId).
		Update("is_default", enums.DictionaryDefaultYes).Error; err != nil {
		return err
	}

	return nil
}

// DictionaryByIdList 根据id列表和系统字典/默认字典，查询是否有不可删字典
// ids必须为[]int/[]string数组
func (d *Dictionary) DictionaryByIdList(ctx context.Context, ids any) (Dictionary, error) {
	var (
		dictionary Dictionary
		err        error
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	curErr := db.Where("id in ?", ids).
		Where("sources = ? or is_default = ?", enums.DictionarySourceSystem, enums.DictionaryDefaultYes).
		First(&dictionary).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return dictionary, err
}

// DictionaryByTypesAndServiceAndName 根据types、name和service查询字典
func (d *Dictionary) DictionaryByTypesAndServiceAndName(ctx context.Context, types, service int, name string) (Dictionary, error) {
	var (
		dictionary Dictionary
		err        error
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	curErr := db.Where("types = ?", types).
		Where("service = ?", service).
		Where("name = ?", name).
		First(&dictionary).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return dictionary, err
}

// DictionaryByTypesAndServiceAndIsDefault 根据types、is_default和service查询字典
func (d *Dictionary) DictionaryByTypesAndServiceAndIsDefault(ctx context.Context, types, service int) (Dictionary, error) {
	var (
		dictionary Dictionary
		err        error
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	curErr := db.Where("types = ?", types).
		Where("service = ?", service).
		Where("is_default = ?", enums.DictionaryDefaultYes).
		First(&dictionary).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return dictionary, err
}

func (d *Dictionary) GetByTypeAndIsDefault(ctx context.Context, typeInt []int, isDefault int) []Dictionary {
	var (
		dictionary []Dictionary
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	if isDefault != 0 {
		db = db.Where("is_default = ?", isDefault)
	}

	db = db.Where("types in ?", typeInt)
	db.Find(&dictionary)
	return dictionary
}

// GetByTypeAndService 一个service有多条数据
func (d *Dictionary) GetByTypeAndService(ctx context.Context, typeInt []int, service int) []Dictionary {
	var (
		dictionary []Dictionary
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	db.
		Where("types in ?", typeInt).
		Where("service = ?", service).
		Find(&dictionary)
	return dictionary
}

// GetDictByServiceAndIsDefault 通过服务和是否默认获取字典
func (d *Dictionary) GetDictByServiceAndIsDefault(ctx context.Context, service int, isDefault int) []Dictionary {
	var (
		dictionary []Dictionary
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)
	db.Where("service = ? and is_default = ?", service, isDefault).Find(&dictionary)
	return dictionary
}

// GetDictById 通过id获取字典
func (d *Dictionary) GetDictById(ctx context.Context, id string) (Dictionary, error) {
	var (
		dictionary Dictionary
		err        error
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	curErr := db.Where("id = ?", id).First(&dictionary).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return dictionary, err
}

// GetDictionaryListByIds 根据id列表获取字典列表
func (d *Dictionary) GetDictionaryListByIds(ctx context.Context, ids []int) []Dictionary {
	var (
		dictionaryList []Dictionary
		db             = mysql.FromContext(ctx).Model(&Dictionary{})
	)

	db.Where("id in ?", ids).Find(&dictionaryList)
	return dictionaryList
}

// GetLangServices 通过获取语言字典
func (d *Dictionary) GetLangServices(ctx context.Context, serviceID int) (Dictionary, error) {
	var (
		dictionary Dictionary
		err        error
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)
	curErr := db.Where("types = ? AND service = ?",
		enums.DictionaryTypeWebPathScan,
		serviceID).First(&dictionary).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return dictionary, err
}

// GetDictByTypeAndIsDefault 通过类型和是否默认获取字典
func (d *Dictionary) GetDictByTypeAndIsDefault(ctx context.Context, types int, isDefault int) []Dictionary {
	var (
		dictionary []Dictionary
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)
	db.Where("types = ? and is_default = ?", types, isDefault).Find(&dictionary)
	return dictionary
}

// All 获取所有字典数据
func (d *Dictionary) All(ctx context.Context) ([]Dictionary, error) {
	var (
		dictionary []Dictionary
		db         = mysql.FromContext(ctx).Model(&Dictionary{})
	)
	db.Find(&dictionary)
	return dictionary, nil
}
