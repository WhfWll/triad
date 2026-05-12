package mysqls

import (
	"context"
	"smart/tools/enums"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

type SystemMessage struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      //主键
	Content    string    `gorm:"column:content" json:"content"`        //消息内容
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` //创建时间
	Type       int       `gorm:"column:type" json:"type"`              //消息类型
	UserId     int       `gorm:"column:user_id" json:"userId"`         //操作人id
	Status     int       `gorm:"column:status" json:"status"`          //消息读取状态
}

// TableName sets insert table name for this struct type
func (s *SystemMessage) TableName() string {
	return "system_message"
}

// GetSystemMessageList 获取系统消息列表
func (s *SystemMessage) GetSystemMessageList(ctx context.Context, page, limit, messageType int, search, startTime, endTime string, uid, role int) ([]SystemMessage, int64, error) {
	var (
		systemMessageList []SystemMessage
		count             int64
		db                = mysql.FromContext(ctx).Model(&SystemMessage{})
	)
	if role == enums.UserRoleOrdinary {
		db = db.Where("user_id = ?", uid)
	}
	if messageType != 0 {
		db = db.Where("type = ?", messageType)
	}
	if search != "" {
		db = db.Where("content LIKE ?", "%"+search+"%")
	}
	if startTime != "" && endTime != "" {
		db = db.Where("create_time > ? and create_time < ?", startTime, endTime)
	}
	db.Count(&count)
	db.Limit(limit).Offset((page - 1) * limit).Order("create_time desc").Find(&systemMessageList)
	return systemMessageList, count, nil
}

// SystemMessageAdd 添加消息
func (s *SystemMessage) SystemMessageAdd(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&SystemMessage{})
	if err := db.Create(s).Error; err != nil {
		return err
	}
	return nil
}
