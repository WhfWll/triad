package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type Toolfile struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Name       string    `gorm:"column:name" json:"name"`              // 文件名称
	FileType   string    `gorm:"column:file_type" json:"fileType"`     // 文件类型
	FilePath   string    `gorm:"column:file_path" json:"filePath"`     // 文件路径
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

func (t *Toolfile) TableName() string {
	return "tool_file"
}

func (t *Toolfile) GetToolfileList(ctx context.Context, page, limit int) ([]Toolfile, int64, error) {
	var (
		toolfileList []Toolfile
		count        int64
		db           = mysql.FromContext(ctx).Model(&Toolfile{})
	)
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&toolfileList)
	return toolfileList, count, nil
}

// Get retrieves a single record of toolfile from database
func (t *Toolfile) GetToolfile(ctx context.Context) (Toolfile, error) {
	var (
		toolfile Toolfile
		err      error
		db       = mysql.FromContext(ctx).Model(&Toolfile{})
	)

	curErr := db.Where("id = ?", t.ID).First(&toolfile).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return toolfile, err
}

// GetToolfileByPath retrieves a single record of toolfile by filePath from database
func (t *Toolfile) GetToolfileByPath(ctx context.Context, filePath string) (Toolfile, error) {
	var (
		toolfile Toolfile
		err      error
		db       = mysql.FromContext(ctx).Model(&Toolfile{})
	)

	curErr := db.Where("file_path = ?", filePath).First(&toolfile).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	if curErr == gorm.ErrRecordNotFound {
		return toolfile, curErr
	}

	return toolfile, err
}

func (t *Toolfile) AddToolfile(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Toolfile{})
	if err := db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

// Update changes toolfile by id
func (t *Toolfile) UpdateToolfile(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Toolfile{})

	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// Delete toolfile by id
func (t *Toolfile) DeleteToolfile(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Toolfile{})

	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}
