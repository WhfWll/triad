package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type Finger struct {
	ID           int        `gorm:"column:id;primary_key" json:"id"`          // 主键
	AppClass     int        `gorm:"column:app_class" json:"appClass"`         // 应用分类
	AppVersion   string     `gorm:"column:app_version" json:"appVersion"`     // 应用版本
	CnName       string     `gorm:"column:cn_name" json:"cnName"`             // 厂商名称
	AppName      string     `gorm:"column:app_name" json:"appName"`           // 应用名称
	Flag         string     `gorm:"column:flag" json:"flag"`                  // 匹配内容
	Source       int        `gorm:"column:source" json:"source"`              // 系统来源 1 系统自带	 2 用户添加
	FingerType   int        `gorm:"column:finger_type" json:"finger_type"`    // 指纹类型  1 web指纹  2 系统指纹  3 设备指纹
	Desc         string     `gorm:"column:desc" json:"desc"`                  // 描述
	Level        string     `gorm:"column:level" json:"level"`                // 分层级别
	CreateTime   time.Time  `gorm:"column:create_time" json:"createTime"`     // 发生时间
	UpdateTime   *time.Time `gorm:"column:update_time" json:"update_time"`    // 更新时间
	FingerSource int        `gorm:"column:finger_source" json:"fingerSource"` // 指纹来源（是否国产化）
}

// TableName sets insert table name for this struct type
func (f *Finger) TableName() string {
	return "finger"
}

// GetFingerList 获取指纹列表
func (f *Finger) GetFingerList(ctx context.Context, page, limit int, name string, level int, class int) ([]Finger, int64, error) {
	var (
		fingerList []Finger
		count      int64
		db         = mysql.FromContext(ctx).Model(&Finger{})
	)
	if name != "" {
		db.Where(`app_name like ? or cn_name like ?`, "%"+name+"%", "%"+name+"%")
	}
	if level != 0 {
		db.Where(`level = ?`, level)
	}
	if class != 0 {
		db.Where(`app_class = ?`, class)
	}
	db.Order("update_time desc")
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&fingerList)
	return fingerList, count, nil
}

// GetAllFinger 获取指纹列表
func (f *Finger) GetAllFinger(ctx context.Context) ([]Finger, error) {
	var (
		fingerList []Finger
		db         = mysql.FromContext(ctx).Model(&Finger{})
	)
	db.Where("id>=0").Find(&fingerList)
	return fingerList, nil
}

// GetFinger 根据id获取一条指纹数据
func (f *Finger) GetFinger(ctx context.Context) (Finger, error) {
	var (
		finger Finger
		err    error
		db     = mysql.FromContext(ctx).Model(&Finger{})
	)
	curErr := db.Where("id = ?", f.ID).First(&finger).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return finger, err
}

// Add persists finger to database
func (f *Finger) AddFinger(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Finger{})

	if err := db.Create(f).Error; err != nil {
		return err
	}

	return nil
}

// Update changes finger by id
func (f *Finger) UpdateFinger(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Finger{})

	if err := db.Where("id = ?", f.ID).Updates(f).Error; err != nil {
		return err
	}

	return nil
}

// Delete finger by id
func (f *Finger) DeleteFinger(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Finger{})

	if err := db.Where("id = ?", f.ID).Delete(f).Error; err != nil {
		return err
	}

	return nil
}

// Count 获取指纹总数
func (f *Finger) Count(ctx context.Context) (int64, error) {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Finger{})
	)
	err := db.Count(&count).Error
	return count, err
}

// DeleteAllFinger 清理掉所有指纹数据
func (f *Finger) DeleteAllFinger(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Finger{})
	if err := db.Exec("DELETE FROM finger").Error; err != nil {
		return err
	}
	return nil
}

// 通过app_name获取指纹
func (f *Finger) LikesFingerForAppName(ctx context.Context, appName string) ([]Finger, error) {
	var (
		finger []Finger
		err    error
		db     = mysql.FromContext(ctx).Model(&Finger{})
	)

	curErr := db.Distinct("app_name").Where("app_name like ?", "%"+appName+"%").Find(&finger).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return finger, err
}

// GetFingerCount 获取指纹数量
func (f *Finger) GetFingerCount(ctx context.Context) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Finger{})
	)
	db.Count(&count)
	return count
}

// GetFingerPartFields 获取指纹的部分字段
func (f *Finger) GetFingerPartFields(ctx context.Context, fields string) ([]Finger, error) {
	var (
		fingerList []Finger
		db         = mysql.FromContext(ctx).Model(&Finger{})
	)
	if fields == "" {
		fields = "*"
	}
	db.Select(fields).Find(&fingerList)
	return fingerList, nil
}

// GetFingerByAppName 获取指纹数据按照appName
func (f *Finger) GetFingerByAppName(ctx context.Context, appName string) (Finger, error) {
	var (
		finger Finger
		err    error
		db     = mysql.FromContext(ctx).Model(&Finger{})
	)
	curErr := db.Where("app_name = ?", appName).First(&finger).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return finger, err
}

// GetFingerListBySource 获取指纹列表
func (f *Finger) GetFingerListBySource(ctx context.Context, source int) ([]Finger, error) {
	var (
		fingerList []Finger
		db         = mysql.FromContext(ctx).Model(&Finger{})
	)
	db.Where("source = ?", source).Find(&fingerList)
	return fingerList, nil
}
