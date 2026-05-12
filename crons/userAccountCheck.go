package crons

import (
	"context"
	log "github.com/sirupsen/logrus"
	"smart/services"
)

func userAccountCheck() {
	ctx := context.Background()

	go func() {
		var err error
		var user services.User

		// 密码校验
		err = user.CheckPasswordExpire(ctx)
		if err != nil {
			log.Error("定时检测账号-密码校验 " + err.Error())
		}

		// 账号校验
		err = user.CheckAccountExpire(ctx)
		if err != nil {
			log.Error("定时检测账号-账号校验 " + err.Error())
		}

		// 长时间未登陆校验
		err = user.CheckAccountLongNoLogin(ctx)
		if err != nil {
			log.Error("定时检测账号-长时间未登陆校验 " + err.Error())
		}
	}()

}
