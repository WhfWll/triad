package mysqls

import (
	"context"
	"encoding/json"
	"smart/tools/enums"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

type LogAudit struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      //主键
	LogType    int       `gorm:"column:log_type" json:"logType"`       //日志类型，1-登录日志，2-操作日志，3-告警日志
	Content    string    `gorm:"column:content" json:"content"`        //日志内容
	Username   string    `gorm:"column:username" json:"username"`      //操作用户
	Ip         string    `gorm:"column:ip" json:"ip"`                  //登录ip
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` //创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` //修改时间
}

// TableName sets insert table name for this struct type
func (l *LogAudit) TableName() string {
	return "log_audit"
}

// LogAuditList 审计日志列表
func (l *LogAudit) LogAuditList(ctx context.Context, page, limit, logType int, search, startTime, endTime string) ([]LogAudit, int64, error) {
	var (
		logAuditList []LogAudit
		count        int64
		db           = mysql.FromContext(ctx).Model(&LogAudit{})
	)
	if search != "" {
		db.Where("content like ?", "%"+search+"%")
	}
	if logType != 0 {
		db.Where("log_type = ?", logType)
	}
	if startTime != "" {
		db.Where("update_time >= ?", startTime)
	}
	if endTime != "" {
		db.Where("update_time <= ?", endTime)
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&logAuditList)
	return logAuditList, count, nil
}

// LogAuditEmpty 清空审计日志
func (l *LogAudit) LogAuditEmpty(ctx context.Context) error {
	db := mysql.FromContext(ctx).Model(&LogAudit{})

	expirationDays := 180
	var mapSetModel MapSet
	mapSet, err := mapSetModel.GetsByObjKey(ctx, enums.LogExpTimeConfigMapSetObjKey)
	if err == nil && mapSet.ObjValue != "" {
		var cfg struct {
			ExpirationTime int `json:"expirationTime"`
		}
		if e := json.Unmarshal([]byte(mapSet.ObjValue), &cfg); e == nil && cfg.ExpirationTime >= 0 {
			expirationDays = cfg.ExpirationTime
		}
	}

	cutoff := time.Now().AddDate(0, 0, -expirationDays).Format(enums.TimeLayout)
	if err := db.Where("create_time < ?", cutoff).Delete(l).Error; err != nil {
		return err
	}
	return nil
}

// LogAuditAdd 新增审计日志
func (l *LogAudit) LogAuditAdd(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&LogAudit{})
	if err := db.Create(l).Error; err != nil {
		return err
	}
	return nil
}

// LogAuditAll 审计日志所有的数据
func (l *LogAudit) LogAuditAll(ctx context.Context) ([]LogAudit, error) {
	var (
		logAuditList []LogAudit
		db           = mysql.FromContext(ctx).Model(&LogAudit{})
	)
	db.Where("1 = 1").Find(&logAuditList)
	return logAuditList, nil
}
