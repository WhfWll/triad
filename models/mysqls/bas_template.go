package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"time"
)

type BasTemplate struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Name       string    `gorm:"column:name" json:"name"`              // 名称
	Desc       string    `gorm:"column:desc" json:"desc"`              // 方案描述
	RuleIds    string    `gorm:"column:rule_ids" json:"ruleIds"`       // 规则ID
	IsDefault  int       `gorm:"column:is_default" json:"isDefault"`   // 是否默认方案
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

// TableName sets insert table name for this struct type
func (b *BasTemplate) TableName() string {
	return "bas_template"
}

// Get retrieves a list of basTemplate from database
func (b *BasTemplate) GetBasTemplateList(ctx context.Context, page, limit int, search string) ([]BasTemplate, int64, error) {
	var (
		basTemplateList []BasTemplate
		count           int64
		db              = mysql.FromContext(ctx).Model(&BasTemplate{})
	)

	if search != "" {
		db = db.Where("name like ?", "%"+search+"%")
	}

	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id DESC").Find(&basTemplateList)

	return basTemplateList, count, nil
}

// Get retrieves a single record of basTemplate from database
func (b *BasTemplate) GetBasTemplate(ctx context.Context) (BasTemplate, error) {
	var (
		basTemplate BasTemplate
		err         error
		db          = mysql.FromContext(ctx).Model(&BasTemplate{})
	)

	curErr := db.Where("id = ?", b.ID).First(&basTemplate).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return basTemplate, err
}

func (b *BasTemplate) GetByIds(ctx context.Context, ids []int) []BasTemplate {
	var (
		basTemplate []BasTemplate
		db          = mysql.FromContext(ctx).Model(&BasTemplate{})
	)

	db.Where("id in ?", ids).Find(&basTemplate)

	return basTemplate
}

// Add persists basTemplate to database
func (b *BasTemplate) AddBasTemplate(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasTemplate{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// Update changes basTemplate by id
func (b *BasTemplate) UpdateBasTemplate(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasTemplate{})

	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

// Delete basTemplate by id
func (b *BasTemplate) DeleteByIds(ctx context.Context, ids []int) error {
	var db = mysql.FromContext(ctx).Model(&BasTemplate{})

	if err := db.Where("id in ?", ids).Delete(b).Error; err != nil {
		return err
	}

	return nil
}

// Get retrieves a single record of basTemplate from database
func (b *BasTemplate) GetByName(ctx context.Context, name string) (BasTemplate, error) {
	var (
		basTemplate BasTemplate
		err         error
		db          = mysql.FromContext(ctx).Model(&BasTemplate{})
	)

	curErr := db.Where("name = ?", name).First(&basTemplate).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return basTemplate, err
}

// 取消默认
func (b *BasTemplate) CancelDefault(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasTemplate{})

	if err := db.Where("is_default = ?", enums.BasTemplateIsDefaultY).Update("is_default", enums.BasTemplateIsDefaultN).Error; err != nil {
		return err
	}

	return nil
}

// 设置默认
func (b *BasTemplate) SetDefault(ctx context.Context, id int) error {
	var db = mysql.FromContext(ctx).Model(&BasTemplate{})

	if err := db.Where("id = ?", id).Update("is_default", enums.BasTemplateIsDefaultY).Error; err != nil {
		return err
	}

	return nil
}
