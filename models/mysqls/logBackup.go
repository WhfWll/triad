package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type LogBackup struct {
	ID   int    `gorm:"column:id;primary_key" json:"id"` //主键
	Name string `gorm:"column:name" json:"name"`         //备份文件名称
	//Content    string    `gorm:"column:content" json:"content"`        //日志内容
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` //备份时间
	Path       string    `gorm:"path" json:"path"`                     //备份日志文件路径
}

// TableName sets insert table name for this struct type
func (l *LogBackup) TableName() string {
	return "log_backup"
}

// LogBackupList 备份日志列表
func (l *LogBackup) LogBackupList(ctx context.Context, page, limit int) ([]LogBackup, int64, error) {
	var (
		logBackupList []LogBackup
		count         int64
		db            = mysql.FromContext(ctx).Model(&LogBackup{})
	)
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&logBackupList)
	return logBackupList, count, nil
}

// LogBackupAdd 新增备份日志
func (l *LogBackup) LogBackupAdd(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&LogBackup{})
	if err := db.Create(l).Error; err != nil {
		return err
	}
	return nil
}

// LogBackupDelete 删除备份日志
func (l *LogBackup) LogBackupDelete(ctx context.Context, id int) error {
	var db = mysql.FromContext(ctx).Model(&LogBackup{})
	if err := db.Where("id = ?", id).Delete(l).Error; err != nil {
		return err
	}
	return nil
}

// LogBackupRecord 单条日志备份信息
func (l *LogBackup) LogBackupRecord(ctx context.Context, id int) (LogBackup, error) {
	var (
		logBackup LogBackup
		err       error
		db        = mysql.FromContext(ctx).Model(&LogBackup{})
	)
	curErr := db.Where("id = ?", id).First(&logBackup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return logBackup, err
}

//查询小于某个创建时间的数据
func (l *LogBackup) GetLogBackupByLtCreateTime(ctx context.Context, creatTime string) ([]LogBackup, error) {
	var (
		logBackupList []LogBackup
		db            = mysql.FromContext(ctx).Model(&LogBackup{})
	)
	db.Where("create_time <= ?", creatTime).Find(&logBackupList)
	return logBackupList, nil
}
