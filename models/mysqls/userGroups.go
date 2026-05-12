package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type UserGroups struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`         // 主键
	Name        string    `gorm:"column:name" json:"name"`                 // 组名称
	Pid         int       `gorm:"column:pid" json:"pid"`                   // 父级ID
	Status      int       `gorm:"column:status" json:"status"`             // 组状态 0:删除 1:正常
	UserID      int       `gorm:"column:user_id" json:"userID"`            // 创建者id
	IsRangeOpen int       `gorm:"column:is_range_open" json:"isRangeOpen"` // 是否开启测试范围,0否，1是  此处状态前后端已对接成了0和1，为了在任务中添加测试范围限制功能更新成了对接完的状态字段
	Range       string    `gorm:"column:range" json:"range"`               // 测试范围
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`    // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`    // 修改时间
}

// TableName sets insert table name for this struct type
func (u *UserGroups) TableName() string {
	return "user_groups"
}

// Get retrieves a list of userGroups from database
func (u *UserGroups) GetUserGroupsList(ctx context.Context, page, limit int) ([]UserGroups, int64, error) {
	var (
		userGroupsList []UserGroups
		count          int64
		db             = mysql.FromContext(ctx).Model(&UserGroups{})
	)
	db.Where("status = ? ", 1).Count(&count)
	db.Limit(limit).Offset(limit*(page-1)).Where("status = ? ", 1).Order("id DESC").Find(&userGroupsList)
	return userGroupsList, count, nil
}

// GetAllUserGroupsList 获取所有有效的用户自信息
func (u *UserGroups) GetAllUserGroupsList(ctx context.Context) ([]UserGroups, error) {
	var (
		userGroupsList []UserGroups
		db             = mysql.FromContext(ctx).Model(&UserGroups{})
	)
	db.Where("status != ? ", 0).Find(&userGroupsList)
	return userGroupsList, nil
}

// GetUserGroupsListByQuery 通过query获取用户组列表
func (u *UserGroups) GetUserGroupsListByQuery(ctx context.Context, query string) ([]UserGroups, error) {
	var (
		userGroupsList []UserGroups
		db             = mysql.FromContext(ctx).Model(&UserGroups{})
	)
	db.Raw(query).Find(&userGroupsList)
	return userGroupsList, nil
}

// GetUserGroupsListByIDs 通过query获取用户组列表
func (u *UserGroups) GetUserGroupsListByIDs(ctx context.Context, ids []int) ([]UserGroups, error) {
	var (
		userGroupsList []UserGroups
		db             = mysql.FromContext(ctx).Model(&UserGroups{})
	)
	db.Where("status != 0 AND id IN ? ", ids).Find(&userGroupsList)
	return userGroupsList, nil
}

// GetUserGroupsInfoByID 通过ID获取用户信息
func (u *UserGroups) GetUserGroupsInfoByID(ctx context.Context, id int) (UserGroups, error) {
	var (
		userGroupsInfo UserGroups
		db             = mysql.FromContext(ctx).Model(&UserGroups{})
	)
	db.Where("id =  ? ", id).First(&userGroupsInfo)
	return userGroupsInfo, nil
}

// Get retrieves a single record of userGroups from database
func (u *UserGroups) GetUserGroups(ctx context.Context) (UserGroups, error) {
	var (
		userGroups UserGroups
		err        error
		db         = mysql.FromContext(ctx).Model(&UserGroups{})
	)

	curErr := db.Where("id = ?", u.ID).First(&userGroups).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return userGroups, err
}

// Add persists userGroups to database
func (u *UserGroups) AddUserGroups(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&UserGroups{})

	if err := db.Create(map[string]interface{}{
		"name":          u.Name,
		"pid":           u.Pid,
		"range":         u.Range,
		"is_range_open": u.IsRangeOpen,
		"user_id":       u.UserID,
		"status":        u.Status,
		"update_time":   u.UpdateTime,
		"create_time":   u.CreateTime}).Error; err != nil {
		return err
	}
	return nil
}

// Update changes userGroups by id
func (u *UserGroups) UpdateUserGroups(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&UserGroups{})

	if err := db.Where("id = ?", u.ID).Updates(map[string]interface{}{"name": u.Name, "pid": u.Pid, "range": u.Range, "is_range_open": u.IsRangeOpen, "update_time": u.UpdateTime}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateUserGroupsStatus 更新用户组状态
func (u *UserGroups) UpdateUserGroupsStatus(ctx context.Context, ids []string) error {
	var db = mysql.FromContext(ctx).Model(&UserGroups{})
	if err := db.Where("id IN ?", ids).Updates(map[string]interface{}{"status": u.Status, "update_time": u.UpdateTime}).Error; err != nil {
		return err
	}
	return nil
}

// Delete userGroups by id
func (u *UserGroups) DeleteUserGroups(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&UserGroups{})

	//u.Estate = "deleted"
	u.UpdateTime = time.Now()
	if err := db.Where("id = ?", u.ID).Updates(u).Error; err != nil {
		return err
	}

	return nil
}

// 所有组
func (u *UserGroups) All(ctx context.Context) []UserGroups {
	var userGroups []UserGroups
	mysql.FromContext(ctx).Model(&UserGroups{}).Find(&userGroups)
	return userGroups
}
