package crons

import (
	"context"
	log "github.com/sirupsen/logrus"
	"smart/services"
	"smart/tools/enums"
	"time"
)

// CheckAccountValidPeriod 检查账号有效期 包含 1 账号有效期结束 2 长期未修改密码
func CheckAccountValidPeriod() {
	log.Println("begin checkAccountValidPeriod")
	ctx := context.Background()
	currentTime := time.Now()
	var (
		userService services.User
		userNames   string
	)
	// 1 获取所有有效的用户
	userList, _, err := userService.GetUserAllValidList(ctx)
	if err != nil {
		log.Error("CheckAccountValidPeriod Get UserList Error: " + err.Error())
		return
	}
	// 2 检查账户过期情况 并 修改账号状态
	for _, userVal := range userList {
		if userVal.Status != enums.UserStatusSuccess || userVal.Type == enums.UserRoleSuperAdmin || userVal.Type == enums.UserRoleAuditor { //超级管理员和审计员没有账户过期
			continue
		}
		if currentTime.After(userVal.AccountExpireTime) {
			// 更新账户状态
			if err = userService.UpdateUserStatus(ctx, 0, userVal.ID, enums.UserStatusAccountExpireDisable); err != nil {
				log.Error("CheckAccountValidPeriod UpdateUserStatus Error: " + err.Error())
				continue
			} else {
				userNames += userVal.Username + " "
			}
		}
		// todo 密码过期 07-15版本 不做
		//if userVal.PasswordExpireTime.Format(enums.TimeYMDBarLayout) > currentTime {
		//	// 更新账户状态
		//	continue
		//}
	}
	log.Println("checkAccountValidPeriod END userList: " + userNames)
	return
}
