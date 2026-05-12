package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type MapSet struct {
	ID       int    `gorm:"column:id;primary_key" json:"id"`  // 主键
	Estate   string `gorm:"column:estate" json:"estate"`      // 数据状态valid/deleted
	ObjKey   string `gorm:"column:obj_key" json:"objKey"`     // 字典键key
	ObjValue string `gorm:"column:obj_value" json:"objValue"` // 字典值value
	Content  string `gorm:"column:content" json:"content"`    // 备注说明
}

// TableName sets insert table name for this struct type
func (m *MapSet) TableName() string {
	return "map_set"
}

// Get retrieves a list of mapSet from database
func (m *MapSet) GetMapSetList(ctx context.Context, page, limit int) ([]MapSet, int64, error) {
	var (
		mapSetList []MapSet
		count      int64
		db         = mysql.FromContext(ctx).Model(&MapSet{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&mapSetList)
	db.Count(&count)

	return mapSetList, count, nil
}

// Get retrieves a single record of mapSet from database
func (m *MapSet) GetMapSet(ctx context.Context) (MapSet, error) {
	var (
		mapSet MapSet
		err    error
		db     = mysql.FromContext(ctx).Model(&MapSet{})
	)

	curErr := db.Where("id = ?", m.ID).First(&mapSet).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return mapSet, err
}

// Add persists mapSet to database
func (m *MapSet) AddMapSet(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&MapSet{})

	if err := db.Create(m).Error; err != nil {
		return err
	}

	return nil
}

// Update changes mapSet by id
func (m *MapSet) UpdateMapSet(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&MapSet{})

	if err := db.Where("id = ?", m.ID).Updates(m).Error; err != nil {
		return err
	}

	return nil
}

// Delete mapSet by id
func (m *MapSet) DeleteMapSet(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&MapSet{})

	m.Estate = "deleted"
	//m.UpdateTime = time.Now()
	if err := db.Where("id = ?", m.ID).Updates(m).Error; err != nil {
		return err
	}

	return nil
}

// 获取单个map
func (m *MapSet) GetsByObjKey(ctx context.Context, objKey string) (MapSet, error) {
	var (
		mapSet MapSet
		err    error
		db     = mysql.FromContext(ctx).Model(&MapSet{})
	)

	if err = db.Where("obj_key = ?", objKey).First(&mapSet).Error; err != nil && err != gorm.ErrRecordNotFound {
		return MapSet{}, err
	}

	return mapSet, nil
}

// 获取系统管理 - 系统配置 - 安全设置
func (m *MapSet) GetsLikeLeftObjKey(ctx context.Context, likeObjKey string) []MapSet {
	var (
		mapSet []MapSet
		db     = mysql.FromContext(ctx).Model(&MapSet{})
	)

	db.Where("obj_key like ?", likeObjKey+"%").Find(&mapSet)

	return mapSet
}

// UpdateObjValueByObjKey 通过key对value更新
func (m *MapSet) UpdateObjValueByObjKey(ctx context.Context, objKey, objValue string) error {
	var db = mysql.FromContext(ctx).Model(&MapSet{})
	if err := db.Where("obj_key = ?", objKey).Updates(MapSet{ObjValue: objValue}).Error; err != nil {
		return err
	}
	return nil
}

// ListByObjKeys 通过objKey数组获取MapSet列表
func (m *MapSet) ListByObjKeys(ctx context.Context, objKeys any) []MapSet {
	var (
		mapSetList []MapSet
		db         = mysql.FromContext(ctx).Model(&MapSet{})
	)

	db.Where("obj_key in ?", objKeys).Find(&mapSetList)

	return mapSetList
}
