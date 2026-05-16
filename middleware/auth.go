package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"smart/services"
	"smart/tools/auth"
	"smart/tools/enums"
	"smart/tools/errno"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/redis"
	"gitlabee.4dogs.cn/common/server"
)

func LoginAuth(c *gin.Context) {
	// 验证系统是否授权
	var mapSetSer services.MapSet
	ctx := context.Background()
	if isAuth := mapSetSer.GetProductAuthState(ctx); isAuth == false {
		if enums.AuthStatusUri[c.Request.URL.Path] {
			fail(c, errno.LoginAuthErr, errors.New(enums.SystemNoAuth))
		}
	}
	// 但前请求地址
	uri := c.Request.URL.Path
	// 是否在白名单
	ok, err := isWhiteApi(uri)
	if err != nil {
		fail(c, errno.LoginAuthErr, err)
		return
	}
	if ok {
		c.Next()
		return
	}

	// 缓存驱动 redis
	cacheClient, err := redis.NewClient()
	if err != nil {
		fail(c, errno.LoginAuthErr, errors.New("redis驱动获取失败"))
		return
	}

	// 是否有临时token
	if tempToken, _ := c.GetQuery("temp_token"); tempToken != "" {
		cacheRes := cacheClient.Get(c, tempToken)
		if tempRes, _ := cacheRes.Result(); tempRes == "Y" {
			// 删除临时token
			cacheClient.Del(c, tempToken)
			c.Next()
			return
		}
	}

	/**
	token机制
	当前系统有两种token：
		1 平台token，在系统管理 - 系统配置 - API密钥处生成
			用途：
				三方平台调用小智接口时携带
				yak脚本调用接口时携带
			携带位置：
				header头（platform-token:xxxxxx）
		2 用户token，登陆时生成，后续每次请求接口都将更新token JWT
			用途：
				用户访问小智各个接口
			携带位置：
				header头 (Authorization:xxx.xxx.xxx)
	使用逻辑：
		优先检查平台token（脚本或三方调用时调用量不可控，优先使用这个，省的每次都校验用户token）
	*/
	platformToken := c.GetHeader("platform-token")

	if platformToken != "" {
		// 优先获取缓存内token
		tokenVar := cacheClient.Get(c, enums.RedisPlatformToken+platformToken)
		if err != nil {
			fail(c, errno.LoginAuthErr, err)
			return
		}
		uid, err := tokenVar.Int()
		if err != nil {
			//fail(c, err)
			//return
		}
		if uid == 0 {
			// 防止缓存击穿
			isLock := cacheClient.SetNX(c, enums.RedisPlatformTokenLock+platformToken, "lock", 1)
			if isLock.Err() == nil {
				// 拿到锁
				// 从数据库获取token
				var userModel services.User
				user, _ := userModel.GetUserForToken(c, platformToken)

				if user.ID == 0 {
					// 缓存起来
					cacheClient.Set(c, enums.RedisPlatformToken+platformToken, -1, 5*time.Second)
					fail(c, errno.LoginAuthErr, errors.New("未知的平台token，请确认"))
					return
				}

				uid = user.ID
				cacheClient.Set(c, enums.RedisPlatformToken+platformToken, uid, 10*time.Hour)
			} else {
				// 提示认证限流
				fail(c, errno.LoginAuthErr, errors.New("平台权限认证失败，请重试"))
				return
			}
		} else if uid == -1 {
			fail(c, errno.LoginAuthErr, errors.New("未知的平台token，请确认"))
			return
		}
		c.Set("uid", uid)
	} else {
		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			fail(c, errno.LoginAuthErr, errors.New("请登录"))
			return
		}

		var jwtDriver auth.JwtFactory
		uid, err := jwtDriver.Auth(token)
		if err != nil {
			fail(c, errno.LoginAuthErr, errors.New("令牌错误，请重新登录获取令牌"))
			return
		}

		if uid == 0 {
			fail(c, errno.LoginAuthErr, errors.New("令牌错误，请重新登录获取令牌"))
			return
		}

		// 验证当前用户是否存在
		var userModel services.User
		user, err := userModel.GetUserForId(c, uid)
		if err != nil {
			fail(c, errno.LoginAuthErr, errors.New("令牌错误，请重新登录获取令牌..."))
			return
		}

		if user.ID == 0 {
			fail(c, errno.LoginAuthErr, errors.New("令牌错误，请重新登录获取令牌...."))
			return
		}

		//业务权限设置，用户组实现后，此处删掉
		bussauth, err := businessAuth(c, c.Request.URL.Path, user.Type)
		if !bussauth || err != nil {
			fail(c, errno.BusinessAuthErr, errors.New("没有该业务模块的权限"))
			return
		}

		// 状态非正常账号
		if user.Status != enums.UserStatusSuccess {
			fail(c, errno.LoginAuthErr, errors.New(enums.UserEnum.GetStatus(user.Status)))
			return
		}

		// 设置账号最新操作时间
		// 拿到锁才更新
		isLock := cacheClient.SetNX(c, enums.RedisUserUpdateOperateTimeLock+strconv.Itoa(user.ID), "Y", 1*time.Second)
		if isLock.Err() == nil {
			err = userModel.UpdateLastOperateById(c, uid)
			if err != nil {
				fail(c, errno.LoginAuthErr, errors.New("最后操作时间更新失败"))
				return
			}
		}

		enableSingleLogin := false
		config.Load("enable_single_login", &enableSingleLogin)
		if enableSingleLogin {
			tokenCmd := cacheClient.HGet(c, enums.RedisUserAlive, strconv.Itoa(uid))
			if tokenCmd.Val() == "" || tokenCmd.Val() != token {
				fail(c, errno.LoginAuthErr, errors.New("该用户未登录"))
				return
			}
		}

		// 仅设置请求上下文，禁止全局缓存，防止污染
		c.Set("uid", uid)
		c.Set("username", user.Username)
		c.Set("role", user.Type)
	}

	// 允许继续
	c.Next()
}

// 是否在api白名单内
func isWhiteApi(uri string) (bool, error) {
	// 获取白名单
	whiteList := make([]string, 0)
	if err := config.Load("white_list", &whiteList); err != nil {
		return false, errors.New("白名单配置获取失败")
	}

	// 检测是否在白名单内
	for _, white := range whiteList {
		if white == uri {
			return true, nil
		}
	}
	return false, nil
}

func fail(c *gin.Context, code int, err error) {
	server.RespFail(c, code, err.Error())
	c.Abort()
	return
}

// 业务权限限制
func businessAuth(ctx context.Context, uri string, userType int) (bool, error) {
	// 超级管理员（type=4）拥有所有权限
	if userType == enums.UserRoleSuperAdmin {
		return true, nil
	}

	var (
		mapSetSer services.MapSet
	)
	value, err := mapSetSer.GetMapValue(ctx, enums.BusinessUserTypeAuthMapSetObjKey)
	if err != nil {
		return false, err
	}
	var bussinessAuthMap map[string]map[string]int
	err = json.Unmarshal([]byte(value), &bussinessAuthMap)
	if err != nil {
		return false, err
	}
	if v, ok := bussinessAuthMap[uri]; ok {
		if _, ok1 := v[strconv.Itoa(userType)]; ok1 {
			return true, nil
		}
		return false, nil
	}
	return true, nil
}
