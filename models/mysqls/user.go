package mysqls

import (
	"context"
	"smart/tools/enums"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID                  int       `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
	Estate              string    `gorm:"column:estate;type:varchar(64);default:valid;comment:数据状态valid/deleted;NOT NULL" json:"estate"`
	Type                int       `gorm:"column:type;type:tinyint(4);default:0;comment:账号类型,1超级管理员，2普通用户;NOT NULL" json:"type"`
	Status              int       `gorm:"column:status;type:tinyint(4);default:0;comment:账号状态;NOT NULL" json:"status"`
	IsAlive             int       `gorm:"column:is_alive;type:tinyint(4);default:1;comment:在线状态，1-在线，2-离线;NOT NULL" json:"is_alive"`
	Username            string    `gorm:"column:username;type:varchar(50);comment:账号名称;NOT NULL" json:"username"`
	RealName            string    `gorm:"column:real_name;type:varchar(50);comment:真实姓名;NOT NULL" json:"real_name"`
	Email               string    `gorm:"column:email;type:varchar(50);comment:邮箱;NOT NULL" json:"email"`
	Department          string    `gorm:"column:department;type:varchar(50);comment:所属部门;NOT NULL" json:"department"`
	Password            string    `gorm:"column:password;type:varchar(128);comment:密码;NOT NULL" json:"password"`
	Remark              string    `gorm:"column:remark;type:varchar(128);comment:描述;NOT NULL" json:"remark"`
	OperatorID          int       `gorm:"column:operator_id;type:int(11);default:0;comment:操作人id;NOT NULL" json:"operator_id"`
	CreateTime          time.Time `gorm:"column:create_time;type:datetime;default:1970-01-01 08:00:01;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime          time.Time `gorm:"column:update_time;type:datetime;default:1970-01-01 08:00:01;comment:修改时间;NOT NULL" json:"update_time"`
	LastLoginTime       time.Time `gorm:"column:last_login_time;type:datetime;default:1970-01-01 08:00:01;comment:上次登录时间;NOT NULL" json:"last_login_time"`
	LastOperateTime     time.Time `gorm:"column:last_operate_time;type:datetime;default:1970-01-01 08:00:01;comment:用户最后操作时间;NOT NULL" json:"last_operate_time"`
	AccountExpireTime   time.Time `gorm:"column:account_expire_time;type:datetime;default:1970-01-01 08:00:01;comment:账号过期时间;NOT NULL" json:"account_expire_time"`
	PasswordExpireTime  time.Time `gorm:"column:password_expire_time;type:datetime;default:1970-01-01 08:00:01;comment:秘密过期时间;NOT NULL" json:"password_expire_time"`
	Token               string    `gorm:"column:token;type:varchar(64);comment:API密钥;NOT NULL" json:"token"`
	TokenCreateTime     time.Time `gorm:"column:token_create_time;type:datetime;default:1970-01-01 08:00:01;comment:密钥生成时间;NOT NULL" json:"token_create_time"`
	TokenExpireTime     time.Time `gorm:"column:token_expire_time;type:datetime;default:1970-01-01 08:00:01;comment:密钥过期时间;NOT NULL" json:"token_expire_time"`
	PasswordChangeTime  time.Time `gorm:"column:password_change_time;type:datetime;default:1970-01-01 08:00:01;comment:密码更新时间;NOT NULL" json:"password_change_time"`
	PasswordAlreadyUsed string    `gorm:"column:password_already_used;type:varchar(2000);comment:已经使用过的密码;NOT NULL" json:"password_already_used"`
}

// TableName sets insert table name for this struct type
func (m *User) TableName() string {
	return "user"
}

// GetAllUserList 获取全部用户信息
func (m *User) GetAllUserList(ctx context.Context) ([]User, int64, error) {
	var (
		myusersUserList []User
		count           int64
		db              = mysql.FromContext(ctx).Model(&User{})
	)

	db.Where("status != ? AND estate = ? ", enums.UserStatusDelete, m.Estate).Find(&myusersUserList)
	db.Count(&count)
	return myusersUserList, count, nil
}

// Get retrieves a list of myusersUser from database
func (m *User) GetUserList(ctx context.Context, page, limit int, search string) ([]User, int64, error) {
	var (
		myusersUserList []User
		count           int64
		db              = mysql.FromContext(ctx).Model(&User{})
	)
	if search != "" {
		db = db.Where("username like ?", "%"+search+"%")
	}
	db.Where("estate = ?", m.Estate).Count(&count)
	db.Limit(limit).Offset(limit*(page-1)).Where("status != ? AND estate = ?", enums.UserStatusDelete, m.Estate).Order("id desc").Find(&myusersUserList)

	return myusersUserList, count, nil
}

// GetAllUserListByType 获取用户信息通过
func (m *User) GetAllUserListByType(ctx context.Context, page, limit int, keyword string) ([]User, int64, error) {
	var (
		myusersUserList []User
		count           int64
		db              = mysql.FromContext(ctx).Model(&User{})
	)
	db.Where("status != ? AND estate = ? AND type = ?", enums.UserStatusDelete, m.Estate, m.Type).Count(&count)
	if keyword != "" {
		db = db.Where("username like ?", "%"+keyword+"%")
	}
	db.Limit(limit).Offset(limit*(page-1)).Where("status != ? AND estate = ? AND type = ?", enums.UserStatusDelete, m.Estate, m.Type).Order("id DESC").Find(&myusersUserList)
	return myusersUserList, count, nil
}

// Get retrieves a single record of myusersUser from database
func (m *User) GetUser(ctx context.Context) (User, error) {
	var (
		myusersUser User
		err         error
		db          = mysql.FromContext(ctx).Model(&User{})
	)

	curErr := db.Where("id = ?", m.ID).First(&myusersUser).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return myusersUser, err
}

// Add persists myusersUser to database
func (m *User) AddUser(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&User{})

	if err := db.Create(m).Error; err != nil {
		return err
	}

	return nil
}

// Update changes myusersUser by id
func (m *User) UpdateUser(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&User{})

	if err := db.Debug().Where("id = ?", m.ID).Updates(m).Error; err != nil {
		return err
	}

	return nil
}

// EditUserInfo 修改用户信息
func (m *User) EditUserInfo(ctx context.Context, updateInfo map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&User{})

	if err := db.Where("id = ?", m.ID).Updates(updateInfo).Error; err != nil {
		return err
	}

	return nil
}

// DeleteUser 软删除用户
func (m *User) DeleteUser(ctx context.Context, ids []string) error {
	var db = mysql.FromContext(ctx).Model(&User{})
	if err := db.Where(" id IN ?", ids).Updates(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *User) GetByIds(ctx context.Context, ids []int) ([]User, error) {
	var (
		myusersUser []User
		err         error
		db          = mysql.FromContext(ctx).Model(&User{})
	)

	curErr := db.Where("id in ?", ids).Find(&myusersUser).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return myusersUser, err
}

func (m *User) GetUserListByIds(ctx context.Context, ids []string) ([]User, error) {
	var (
		myusersUser []User
		err         error
		db          = mysql.FromContext(ctx).Model(&User{})
	)
	curErr := db.Where("id in ?", ids).Find(&myusersUser).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return myusersUser, err
}

// 依据username获取用户
func (m *User) GetByUsername(ctx context.Context, username string) (User, error) {
	var (
		myusersUser User
		err         error
		db          = mysql.FromContext(ctx).Model(&User{})
	)
	curErr := db.Where("username = ? and estate = ?", username, "valid").First(&myusersUser).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return myusersUser, err
}

// GetByUsernameOrEmail 依据username or email获取用户
func (m *User) GetByUsernameOrEmail(ctx context.Context, username, email string) (User, error) {
	var (
		myusersUser User
		err         error
		db          = mysql.FromContext(ctx).Model(&User{})
	)

	curErr := db.Where("(username = ? or email = ?) and estate = ?", username, email, "valid").First(&myusersUser).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return myusersUser, err
}

// 验证用户
func (m *User) CheckUserIsAlready(ctx context.Context, uid int, username, email string) (bool, error) {
	var (
		myusersUser User
		err         error
		db          = mysql.FromContext(ctx).Model(&User{})
	)
	if err = db.Where("id != ? and (username = ? or email = ?) and estate = ?", uid, username, email, "valid").Find(&myusersUser).Error; err != nil {
		return false, err
	}

	if myusersUser.ID > 0 {
		return true, nil
	}
	return false, nil
}

// UpdateStatusById 设置用户状态
func (m *User) UpdateStatusById(ctx context.Context, id, status int) error {
	var db = mysql.FromContext(ctx).Model(&User{})
	if err := db.Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (m *User) UpdateLastOperateById(ctx context.Context, id int) error {
	var db = mysql.FromContext(ctx).Model(&User{})

	if err := db.Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_time":   time.Now(),
		"last_operate_time": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

// 获取最后操作时间小于某个时间的角色用户
func (m *User) GetsByRoleAndLastOperateTime(ctx context.Context, role int, expireTime int64, status int) []User {
	var (
		myusersUser []User
		db          = mysql.FromContext(ctx).Model(&User{})
	)

	// 普通用户
	db.Where("type = ?", role).
		Where("status = ?", status).
		// 最后操作时间
		Where("last_operate_time < ?", expireTime).
		Find(&myusersUser)

	return myusersUser
}

// 获取账号有效期小于某个时间的角色用户
func (m *User) GetsByRoleAndAccountTime(ctx context.Context, role int, expireTime int64, status int) []User {
	var (
		myusersUser []User
		db          = mysql.FromContext(ctx).Model(&User{})
	)

	// 普通用户
	db.Where("type = ?", role).
		Where("status = ?", status).
		Where("account_expire_time < ?", expireTime).
		Find(&myusersUser)

	return myusersUser
}

// 获取密码有效期小于某个时间的角色用户
func (m *User) GetsByRoleAndPasswordTime(ctx context.Context, role int, expireTime int64, status int) []User {
	var (
		myusersUser []User
		db          = mysql.FromContext(ctx).Model(&User{})
	)

	// 普通用户
	db.Where("type = ?", role).
		Where("status = ?", status).
		Where("password_expire_time < ?", expireTime).
		Find(&myusersUser)

	return myusersUser
}

// GetUserByToken 通过token获取用户
func (m *User) GetUserByToken(ctx context.Context, token string) (User, error) {
	var (
		myusersUser User
		err         error
		db          = mysql.FromContext(ctx).Model(&User{})
	)
	curErr := db.Where("token = ?", token).First(&myusersUser).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return myusersUser, err
}

// UpdatePlatformToken 更新平台token
func (m *User) UpdatePlatformToken(ctx context.Context, username, token string) error {
	var (
		err error
		db  = mysql.FromContext(ctx).Model(&User{})
	)
	curErr := db.Debug().Where("username = ?", username).Updates(User{Token: token, TokenCreateTime: time.Now()}).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return err
}

// GetUserListByToken 获取token列表
func (m *User) GetUserListByToken(ctx context.Context, search string, page, limit int) ([]User, int64, error) {
	var (
		err      error
		db       = mysql.FromContext(ctx).Model(&User{})
		count    int64
		userList []User
	)
	db.Where("token != '' and estate = ?", m.Estate).Count(&count)
	db.Limit(limit).Offset(limit*(page-1)).Where("token != '' and status != ? AND estate = ?", enums.UserStatusDelete, m.Estate).Order("id desc").Find(&userList)
	return userList, count, err
}

// UpdatePasswordChangeTime 更新密码更新时间
func (m *User) UpdatePasswordChangeTime(ctx context.Context, id int, changeTime time.Time) error {
	var (
		err error
		db  = mysql.FromContext(ctx).Model(&User{})
	)
	curErr := db.Debug().Where("id = ?", id).Updates(User{PasswordChangeTime: changeTime}).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return err
}

// TokenDelete 秘钥删除接口
func (m *User) TokenDelete(ctx context.Context, username, token string) error {
	db := mysql.FromContext(ctx).Model(&User{})
	updateData := map[string]interface{}{"token": ""}
	curErr := db.Where("token = ?", token).Updates(updateData).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		return curErr
	}
	return nil
}
