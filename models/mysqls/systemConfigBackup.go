package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type SystemConfigBackup struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      //主键
	Name       string    `gorm:"column:name" json:"name"`              //备份文件名称
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` //备份时间
	Path       string    `gorm:"path" json:"path"`                     //备份文件路径
}

// TableName sets insert table name for this struct type
func (s *SystemConfigBackup) TableName() string {
	return "system_config_backup"
}

// SystemConfigBackupList 系统配置备份列表
func (s *SystemConfigBackup) SystemConfigBackupList(ctx context.Context, page, limit int) ([]SystemConfigBackup, int64, error) {
	var (
		systemConfigBackupList []SystemConfigBackup
		count                  int64
		db                     = mysql.FromContext(ctx).Model(&SystemConfigBackup{})
	)
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&systemConfigBackupList)
	return systemConfigBackupList, count, nil
}

// SystemConfigBackupAdd 新增系统配置备份
func (s *SystemConfigBackup) SystemConfigBackupAdd(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&SystemConfigBackup{})
	if err := db.Create(s).Error; err != nil {
		return err
	}
	return nil
}

// SystemConfigBackupDelete 删除系统配置备份
func (s *SystemConfigBackup) SystemConfigBackupDelete(ctx context.Context, id int) error {
	var db = mysql.FromContext(ctx).Model(&SystemConfigBackup{})
	if err := db.Where("id = ?", id).Delete(s).Error; err != nil {
		return err
	}
	return nil
}

// SystemConfigBackupRecord 单条系统配置备份信息
func (s *SystemConfigBackup) SystemConfigBackupRecord(ctx context.Context, id int) (SystemConfigBackup, error) {
	var (
		systemConfigBackup SystemConfigBackup
		err                error
		db                 = mysql.FromContext(ctx).Model(&SystemConfigBackup{})
	)
	curErr := db.Where("id = ?", id).First(&systemConfigBackup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return systemConfigBackup, err
}
