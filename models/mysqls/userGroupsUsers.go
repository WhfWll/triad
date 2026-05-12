package mysqls

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type UserGroupsUsers struct {
	ID           int       `gorm:"column:id;primary_key" json:"id"`           // 主键
	UserGroupsID int       `gorm:"column:user_groups_id" json:"userGroupsID"` // 用户组ID
	UserID       int       `gorm:"column:user_id" json:"userID"`              // 用户ID
	SubmitUserID int       `gorm:"column:submit_user_id" json:"submitUserID"` // 提交人用户ID
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`      // 创建时间
}

// TableName sets insert table name for this struct type
func (u *UserGroupsUsers) TableName() string {
	return "user_groups_users"
}

// Get retrieves a list of userGroupsUsers from database
func (u *UserGroupsUsers) GetUserGroupsUsersList(ctx context.Context, page, limit int) ([]UserGroupsUsers, int64, error) {
	var (
		userGroupsUsersList []UserGroupsUsers
		count               int64
		db                  = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&userGroupsUsersList)
	db.Count(&count)

	return userGroupsUsersList, count, nil
}

// Get retrieves a single record of userGroupsUsers from database
func (u *UserGroupsUsers) GetUserGroupsUsers(ctx context.Context) (UserGroupsUsers, error) {
	var (
		userGroupsUsers UserGroupsUsers
		err             error
		db              = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	)

	curErr := db.Where("id = ?", u.ID).First(&userGroupsUsers).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return userGroupsUsers, err
}

// Add persists userGroupsUsers to database
func (u *UserGroupsUsers) AddUserGroupsUsers(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&UserGroupsUsers{})

	if err := db.Create(u).Error; err != nil {
		return err
	}

	return nil
}

// Update changes userGroupsUsers by id
func (u *UserGroupsUsers) UpdateUserGroupsUsers(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&UserGroupsUsers{})

	if err := db.Where("id = ?", u.ID).Updates(u).Error; err != nil {
		return err
	}

	return nil
}

// Delete userGroupsUsers by id
func (u *UserGroupsUsers) DeleteUserGroupsUsers(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&UserGroupsUsers{})

	if err := db.Where("id = ?", u.ID).Updates(u).Error; err != nil {
		return err
	}

	return nil
}

// 通过用户IDs 查询
func (u *UserGroupsUsers) GetByUserIds(ctx context.Context, userIds []int) []UserGroupsUsers {
	var (
		userGroupsUsers []UserGroupsUsers
		db              = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	)

	userIdStr := make([]string, 0)
	for _, item := range userIds {
		userIdStr = append(userIdStr, strconv.Itoa(item))
	}

	db.Where("id in (?)", strings.Join(userIdStr, ",")).Find(&userGroupsUsers)
	return userGroupsUsers
}

// 根据用户 ID 更新状态
func (u *UserGroupsUsers) DeleteByUserId(ctx context.Context, userId int) error {
	var db = mysql.FromContext(ctx).Model(&UserGroupsUsers{})

	if err := db.Where("user_id = ?", userId).Delete(&u).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserGroupsUsers) AddsUserGroupsUsers(ctx context.Context, datas []UserGroupsUsers) error {
	var db = mysql.FromContext(ctx).Model(&UserGroupsUsers{})

	if err := db.Create(datas).Error; err != nil {
		return err
	}

	return nil
}

// GetUserGroupsUserListByQuery 通过query获取用户组列表
func (u *UserGroupsUsers) GetUserGroupsUserListByQuery(ctx context.Context, query string) ([]UserGroupsUsers, error) {
	var (
		userGroupsList []UserGroupsUsers
		db             = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	)
	db.Raw(query).Find(&userGroupsList)
	return userGroupsList, nil
}

// GetUserGroupsAllUserList 获取用户组列表
func (u *UserGroupsUsers) GetUserGroupsAllUserList(ctx context.Context) ([]UserGroupsUsers, error) {
	var (
		userGroupsList []UserGroupsUsers
		db             = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	)
	db.Find(&userGroupsList)
	return userGroupsList, nil
}

// DeleteUserGroupsUsersByGroupID 通过用户组ID删除用户组
func (u *UserGroupsUsers) DeleteUserGroupsUsersByGroupID(ctx context.Context, groupID int) error {
	var db = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	if err := db.Where("user_groups_id = ?", groupID).Delete(u).Error; err != nil {
		return err
	}
	return nil
}

// GetUserGroupsListByUserIDs 通过userids获取用户组列表
func (u *UserGroupsUsers) GetUserGroupsListByUserIDs(ctx context.Context, ids []int) ([]UserGroupsUsers, error) {
	var (
		userGroupsList []UserGroupsUsers
		db             = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	)
	db.Where(" user_id IN ? ", ids).Find(&userGroupsList)
	return userGroupsList, nil
}

func (u *UserGroupsUsers) GetUserGroupsListByUserGroupIDs(ctx context.Context, ids []int) ([]UserGroupsUsers, error) {
	var (
		userGroupsList []UserGroupsUsers
		db             = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	)
	db.Where(" user_groups_id IN ? ", ids).Find(&userGroupsList)
	return userGroupsList, nil
}

// AddUserGroupsUsersByString 添加用户与组的关联(批量)
func (u *UserGroupsUsers) AddUserGroupsUsersByString(ctx context.Context, userID int, groupIds string) error {
	var db = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	// 解析 groupIds
	groupIDs := strings.Split(groupIds, ",")
	var users []UserGroupsUsers
	for _, groupID := range groupIDs {
		if groupID == "" {
			continue
		}
		gID, err := strconv.Atoi(groupID)
		if err != nil {
			continue
		}
		users = append(users, UserGroupsUsers{
			UserID:       userID,
			UserGroupsID: gID,
			CreateTime:   time.Now(),
		})
	}
	if len(users) > 0 {
		if err := db.Create(&users).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpdateUserGroupsUsersByString 更新用户与组的关联
func (u *UserGroupsUsers) UpdateUserGroupsUsersByString(ctx context.Context, userID int, groupIds string) error {
	var db = mysql.FromContext(ctx).Model(&UserGroupsUsers{})
	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除原有关联
		if err := tx.Where("user_id = ?", userID).Delete(&UserGroupsUsers{}).Error; err != nil {
			return err
		}
		// 2. 添加新关联
		if groupIds != "" {
			groupIDs := strings.Split(groupIds, ",")
			var users []UserGroupsUsers
			for _, groupID := range groupIDs {
				if groupID == "" {
					continue
				}
				gID, err := strconv.Atoi(groupID)
				if err != nil {
					continue
				}
				users = append(users, UserGroupsUsers{
					UserID:       userID,
					UserGroupsID: gID,
					CreateTime:   time.Now(),
				})
			}
			if len(users) > 0 {
				if err := tx.Create(&users).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
