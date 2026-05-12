package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"time"
)

type UserGoPhishKey struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	UserID     int       `gorm:"column:user_id" json:"userId"`         // UserId
	GoPhishKey string    `gorm:"column:gophish_key" json:"gophishKey"` // gophish_api_key
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

// TableName sets insert table name for this struct type
func (u *UserGoPhishKey) TableName() string {
	return "user_gophish_key"
}

// GetUserGoPhishKeyInfoByUserID 通过用户ID获取信息
func (u *UserGoPhishKey) GetUserGoPhishKeyInfoByUserID(ctx context.Context, userid int) (UserGoPhishKey, error) {
	var userGroupsInfo UserGoPhishKey
	err := mysql.FromContext(ctx).Model(&UserGoPhishKey{}).Where("user_id =  ? ", userid).Find(&userGroupsInfo).Error
	return userGroupsInfo, err
}

// AddUserGoPhishKey 。
func (u *UserGoPhishKey) AddUserGoPhishKey(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&UserGoPhishKey{})

	if err := db.Create(map[string]interface{}{
		"user_id":     u.UserID,
		"gophish_key": u.GoPhishKey,
		"update_time": time.Now(),
		"create_time": time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}
