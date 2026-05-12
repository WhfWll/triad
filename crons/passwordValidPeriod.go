package crons

import (
	"context"
	log "github.com/sirupsen/logrus"
	"smart/services"
	"smart/tools/enums"
	"strconv"
	"time"
)

// CheckPasswordValidPeriod 检查密码有效期 包含 1 长期未修改密码检查   2 长期未登录检查
func CheckPasswordValidPeriod() {
	log.Println("begin CheckPasswordValidPeriod")
	ctx := context.Background()
	currentTime := time.Now().Format(enums.TimeYMDBarLayout)
	var (
		userService services.User
		userNames   string
	)
	// 1 获取所有有效的用户
	userList, _, err := userService.GetUserAllValidList(ctx)
	if err != nil {
		log.Error("CheckAccountValidPeriod Get UserList Error: " + err.Error())
	}
	// 2 检查密码过期和账户登陆情况 并 修改账号状态
	var zeroTime time.Time
	var mapSetService services.MapSet
	passwordValidString, err := mapSetService.GetMapValue(ctx, enums.SystemSecurityPasswordValid)
	passwordValid, err := strconv.Atoi(passwordValidString)
	if err != nil {
		passwordValid = 0
	}
	ExpiredUnusedString, err := mapSetService.GetMapValue(ctx, enums.SystemSecurityExpireUnused)
	ExpiredUnused, err := strconv.Atoi(ExpiredUnusedString)
	if err != nil {
		ExpiredUnused = 0
	}
	for _, userVal := range userList {
		if userVal.Status != enums.UserStatusSuccess || userVal.Type == enums.UserRoleSuperAdmin || userVal.Type == enums.UserRoleAuditor { //超级管理员和审计员没有账户过期
			continue
		}
		// 3 如果密码修改时间为空，就设置当前时间为密码修改时间
		if userVal.PasswordChangeTime == zeroTime {
			if err = userService.UpdatePasswordChangeTime(ctx, userVal.ID, time.Now()); err != nil {
				log.Error("CheckPasswordValidPeriod UpdatePasswordChangeTime Error: " + err.Error())
			}
		} else if passwordValid != 0 {
			// 4 检查密码是否超期未修改
			if userVal.PasswordChangeTime.Add(time.Duration(passwordValid*31*24)*time.Hour).Format(enums.TimeYMDBarLayout) < currentTime {
				if err = userService.UpdateUserStatus(ctx, 0, userVal.ID, enums.UserStatusPassExpireDisable); err != nil {
					log.Error("CheckAccountValidPeriod UpdateUserStatus Error: " + err.Error())
				} else {
					userNames += userVal.Username + " "
				}
			}
		}
		// 5 检查账户是否超期未登录
		if ExpiredUnused != 0 {
			if userVal.LastLoginTime.Add(time.Duration(ExpiredUnused*24)*time.Hour).Format(enums.TimeYMDBarLayout) < currentTime {
				if err = userService.UpdateUserStatus(ctx, 0, userVal.ID, enums.UserStatusNoLoginDisable); err != nil {
					log.Error("CheckAccountValidPeriod UpdateUserStatus Error: " + err.Error())
				} else {
					userNames += userVal.Username + " "
				}
			}
		}
	}
	log.Println("CheckPasswordValidPeriod END userList: " + userNames)
	return
}
