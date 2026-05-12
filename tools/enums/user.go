package enums

import "time"

const (
	ResetPassWordInfo = "admin@#123" // 重置密码
	UserOnLine        = 1
	UserOffline       = 2
)

// 角色类型
const (
	UserIdentity      = 1 // UserIdentity 身份
	UserPresence      = 2 // UserPresence 在线状态
	UserAccountStatus = 3 // UserAccountStatus 账号状态
)

// 角色
const (
	UserRoleOrdinary      = 1 // 普通用户
	UserRoleAdministrator = 2 // 管理员
	UserRoleAuditor       = 3 // 审计员
	UserRoleSuperAdmin    = 4 // 超级管理员
)

// 状态
const (
	UserStatusDelete               = 0
	UserStatusSuccess              = 1
	UserStatusLoginDisable         = 2
	UserStatusAccountExpireDisable = 3
	UserStatusPassExpireDisable    = 4
	UserStatusSetDisable           = 5
	UserStatusNoLoginDisable       = 6
	PresenceStatusLive             = 0
	PresenceStatusOffline          = 1
)

var RoleChoice = map[int]string{
	UserRoleOrdinary:      "普通用户",
	UserRoleAdministrator: "管理员",
	UserRoleAuditor:       "审计员",
	UserRoleSuperAdmin:    "超级管理员",
}

var StatusMap = map[int]string{
	UserStatusDelete:               "账号已被删除",
	UserStatusSuccess:              "账号正常",
	UserStatusLoginDisable:         "登录失败次数过多禁用",
	UserStatusAccountExpireDisable: "账号有效期超期禁用",
	UserStatusPassExpireDisable:    "密码长期未修改禁用",
	UserStatusSetDisable:           "手动禁用",
	UserStatusNoLoginDisable:       "账号长时间未登录禁用",
}

// PresenceMap 在线状态枚举
var PresenceMap = map[int]string{
	PresenceStatusLive:    "在线",
	PresenceStatusOffline: "离线",
}

// 是否打开测试范围
const (
	IsRangeOpenNo  = 0 //  否
	IsRangeOpenYes = 1 //  是
)
const AdminCacheKey = "smart:user:admin"

var UserEnum user

type user struct {
}

func (user *user) GetRoleChoice(role int) string {
	if res, ok := RoleChoice[role]; ok {
		return res
	}
	return ""
}

func (user *user) GetStatus(status int) string {
	if res, ok := StatusMap[status]; ok {
		return res
	}
	return "当前账号状态未知"
}

func (user *user) GetAccountExpire(accountExpire time.Time) string {
	if accountExpire.IsZero() {
		return "2099-01-01"
	}
	return accountExpire.Format("2006-01-02")
}

func (user *user) GetLastOperateTime(lastOperateTime time.Time) string {
	if lastOperateTime.IsZero() {
		return ""
	}
	return lastOperateTime.Format("2006-01-02")
}
