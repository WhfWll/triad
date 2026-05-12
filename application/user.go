package application

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"image/color"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/services"
	"smart/tools/encryption"
	"smart/tools/enums"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/captcha/base64Captcha"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/redis"
)

type User struct {
}

// Logout 退出登录
func (u *User) Logout(ctx context.Context, userID int) error {
	var userService services.User
	if err := userService.UpdateIsAlive(ctx, userID, enums.UserOffline, ""); err != nil {
		return err
	}
	return nil
}

// ChangeUserStatus 修改用户状态
func (u *User) ChangeUserStatus(ctx context.Context, opid int, req *typespec.ChangeUserStatusReq) error {
	var userService services.User
	if errStr := userService.CheckUserStatusParam(ctx, opid, req.UserID, req.Status); errStr != "" {
		return errors.New(errStr)
	}
	if err := userService.UpdateUserStatus(ctx, opid, req.UserID, req.Status); err != nil {
		return err
	}
	return nil
}

// UpdateUserExp 修改用户过期时间
func (u *User) UpdateUserExp(ctx context.Context, opid int, req *typespec.UpdateUserExpReq) error {
	var userService services.User
	if errStr := userService.CheckUserExpireTimeParam(ctx, opid, req.UserID, req.AccountExpireTime); errStr != "" {
		return errors.New(errStr)
	}
	if err := userService.UpdateAccountExpireTime(ctx, opid, req.UserID, req.AccountExpireTime); err != nil {
		return err
	}
	return nil
}

// ResetPassWord 重置密码
func (u *User) ResetPassWord(ctx context.Context, opid int, req *typespec.UserResetPassPWReq) error {
	var userService services.User
	// 判断密码是否正确
	if errStr := userService.CheckUserPassParam(ctx, opid, req.UserID, services.PasswordDecodeString(req.Password)); errStr != "" {
		return errors.New(errStr)
	}
	if err := userService.UpdatePassWord(ctx, opid, req.UserID, services.PasswordEncryption(enums.ResetPassWordInfo), ""); err != nil {
		return err
	}
	return nil
}

// UpdatePassWord 修改密码
func (u *User) UpdatePassWord(ctx context.Context, uid int, req *typespec.UserPassPWReq) error {
	var userService services.User
	req.Password = services.PasswordDecodeString(req.Password)
	req.Repassword = services.PasswordDecodeString(req.Repassword)
	if errStr := userService.CheckUserEditPassParam(ctx, *req, uid); errStr != "" {
		return errors.New(errStr)
	}
	if err := userService.UpdatePassWord(ctx, uid, req.UserID, services.PasswordEncryption(req.Password), req.OldPassword); err != nil {
		return err
	}
	return nil
}

// DelUser 删除用户
func (u *User) DelUser(ctx context.Context, uid int, req *typespec.UserDelReq) error {
	var userService services.User
	if errStr := userService.CheckUserDelParam(ctx, uid, req.UserIds); errStr != "" {
		return errors.New(errStr)
	}
	if err := userService.DelUser(ctx, uid, req.UserIds); err != nil {
		return err
	}
	return nil
}

// UserDetail 用户详情
func (u *User) UserDetail(ctx context.Context, req *typespec.UserDetailReq, resp *typespec.UserListInfo) error {
	var userService services.User
	users, err := userService.GetUserDetail(ctx, req.UserID)
	if err != nil {
		return err
	}
	*resp = typespec.UserListInfo{
		ID:            users.ID,
		Username:      users.Username,
		RoleStr:       enums.RoleChoice[users.Type],
		Role:          users.Type,
		AccountExpire: users.AccountExpireTime.Format(enums.TimeLayout),
		LastTime:      users.LastLoginTime.Format(enums.TimeLayout),
		Status:        users.Status,
		StatusStr:     enums.StatusMap[users.Status],
		Email:         users.Email,
		Department:    users.Department,
		Remark:        users.Remark,
	}
	return nil
}

// UserList 用户列表
func (u *User) UserList(ctx context.Context, req *typespec.UserListReq, resp *typespec.UserListRes) error {
	var userService services.User
	users, count, err := userService.GetUserList(ctx, req.Search, req.Page, req.Size)
	if err != nil {
		return err
	}
	// 获取用户组信息
	var userIds []int
	for _, v := range users {
		userIds = append(userIds, v.ID)
	}
	var userGroupListSrv services.UserGroupUserList
	groupNameMap, _ := userGroupListSrv.GetUserGroupNamesMap(ctx, userIds)
	groupIdsMap, _ := userGroupListSrv.GetUserGroupIdsMap(ctx, userIds)

	var userList []typespec.UserListInfo
	for _, v := range users {
		gIds := groupIdsMap[v.ID]
		if gIds == nil {
			gIds = []int{}
		}
		userList = append(userList, typespec.UserListInfo{
			ID:            v.ID,
			Username:      v.Username,
			RoleStr:       enums.RoleChoice[v.Type],
			Role:          v.Type,
			AccountExpire: v.AccountExpireTime.Format(enums.TimeYMDBarLayout),
			LastTime:      v.LastLoginTime.Format(enums.TimeYMDHMinLayout),
			Email:         v.Email,
			Department:    v.Department,
			Status:        v.Status,
			StatusStr:     enums.StatusMap[v.Status],
			Remark:        v.Remark,
			GroupStr:      groupNameMap[v.ID],
			GroupIds:      gIds,
			IsAlive:       v.IsAlive,
		})
	}
	*resp = typespec.UserListRes{
		Page:  req.Page,
		Size:  req.Size,
		Total: count,
		List:  userList,
	}
	return nil
}

// UserEnumList 用户枚举列表
func (u *User) UserEnumList(ctx context.Context, req *typespec.UserEnumListReq, resp *typespec.UserEnumListRes) error {
	switch req.Type {
	case enums.UserIdentity:
		resp.Data = enums.RoleChoice
	case enums.UserPresence:
		resp.Data = enums.PresenceMap
	case enums.UserAccountStatus:
		resp.Data = enums.StatusMap
	}
	return nil
}

// 验证码 全局缓存
var captchaStore = base64Captcha.DefaultMemStore

// 验证码 UserLoginCaptcha
func (u *User) UserLoginCaptcha(ctx context.Context, req *typespec.UserLoginCaptchaReq, resp *typespec.UserLoginCaptchaRes) error {

	var captchaMath base64Captcha.DriverMath

	// 验证码宽高
	captchaMath.Height = 40
	captchaMath.Width = 120
	captchaMath.BgColor = &color.RGBA{R: 255, G: 255, B: 255}
	captchaMath.Fonts = []string{"arial.ttf"}
	driverMath := captchaMath.ConvertFonts()

	id, b64s, err := base64Captcha.NewCaptcha(driverMath, captchaStore).Generate()
	if err != nil {
		return err
	}

	// 记录验证码有效期
	cacheClient, err := redis.NewClient()
	if err != nil {
		return err
	}
	cacheClient.Set(ctx, id, "Y", 5*time.Minute)

	resp.CaptchaId = id
	resp.Data = b64s
	return nil
}

// UserLogin 用户登陆
func (u *User) UserLogin(ctx context.Context, req *typespec.UserLoginReq, resp *typespec.UserLoginRes) (err error) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("user login panic: %v", p)
			err = errors.New("系统内部错误")
		}
	}()
	cacheClient, err := redis.NewClient()
	if err != nil {
		return err
	}
	// 验证码是否过期
	captchaVal := cacheClient.Get(ctx, req.CaptchaId)
	if captchaVal.Val() != "Y" {
		return errors.New("验证码已过期，请点击验证码刷新")
	}
	// 验证码是否正确
	if !captchaStore.Verify(req.CaptchaId, req.CheckCode, true) {
		return errors.New("验证码错误")
	}
	var (
		userSrv services.User // 利用用户名 获取用户信息
		//logAuditService services.LogAudit //增加审计日志-登录日志
		// 处理密码-解密
		aesCbc  encryption.AesCbc
		certVal string
		pbkdf   encryption.Pbkdf2Sha256
	)
	user, err := userSrv.GetForUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	// 判断用户信息是否存在
	if user.ID == 0 {
		return errors.New("用户名或密码错误")
	}
	// 判断账号状态 非成功状态 均返回错误
	if user.Status != enums.UserStatusSuccess {
		return errors.New(enums.UserEnum.GetStatus(user.Status))
	}
	// 判断用户是否被限制访问
	isLimit := cacheClient.Get(ctx, req.Username+"_limit")
	if isLimit.Val() != "" {
		expire := cacheClient.TTL(ctx, req.Username+"_limit")
		return errors.New("账号已被锁定，请于" + expire.String() + "后重试")
	}
	if err = config.Load("aesCbcUserCert", &certVal); err != nil {
		return err
	}
	if certVal == "" {
		return errors.New("服务端密钥缺失")
	}
	hexDePass, err := hex.DecodeString(req.Password)
	if err != nil {
		return errors.New("密码加密协商不一致，禁止登录")
	}
	// 判断密码是否正确
	if ok, _ := pbkdf.PasswordVerify(string(aesCbc.AesDecryptCBC(hexDePass, []byte(certVal))), user.Password); !ok {
		// 验证-密码不正确
		target, err := recodeLoginErr(ctx, user.Username, user.Type)
		switch target {
		case "disable":
			_ = userSrv.UpdateStatusById(ctx, user.ID, enums.UserStatusLoginDisable)
		case "limit":
		}
		return err
	}
	// 生成token
	token, err := userSrv.GenerateToken(user.ID)
	if err != nil {
		return err
	}

	/*
		// 销售许可证使用，判断是否为首次登录或者密码过期，如果是则强制修改密码
		if user.LastLoginTime.Format(enums.TimeYMDBarLayout) == "" || user.LastLoginTime.Format(enums.TimeYMDBarLayout) == "1970-01-01" {
			resp.FirstLogin = true
		}
		var mapSetService services.MapSet
		passwordValidString, err := mapSetService.GetMapValue(ctx, enums.SystemSecurityPasswordValid)
		passwordValid, err := strconv.Atoi(passwordValidString)
		if err != nil {
			passwordValid = 0
		}
		if user.PasswordChangeTime.Add(time.Duration(passwordValid*31*24)*time.Hour).Format(enums.TimeYMDBarLayout) < time.Now().Format(enums.TimeYMDBarLayout) {
			resp.FirstLogin = true
		}
		if ok, _ := pbkdf.PasswordVerify("admin123321", user.Password); ok {
			resp.FirstLogin = true
		}
		if ok, _ := pbkdf.PasswordVerify("luckyadmin123", user.Password); ok {
			resp.FirstLogin = true
		}
	*/

	// 生成jwt token 在全局 response中间件内 登录成功
	resp.Token = token
	resp.Uid = user.ID
	resp.Role = user.Type
	// 修改用户状态为在线
	userSrv.UpdateIsAlive(ctx, resp.Uid, enums.UserOnLine, "")

	cacheClient.HSet(ctx, enums.RedisUserAlive, user.ID, token)
	return nil
}

// 记录登录失败信息
func recodeLoginErr(ctx context.Context, username string, role int) (string, error) {
	// 记录登录错误次数
	cacheClient, err := redis.NewClient()
	if err != nil {
		return "", err
	}
	errNumObj, err := cacheClient.Incr(ctx, username).Result()
	if err != nil {
		return "", err
	}
	errNum := int(errNumObj)

	if errNum == 1 {
		cacheClient.Expire(ctx, username, 5*time.Minute)
	}

	// 全局安全策略
	var mapSet services.MapSet
	strategy, err := mapSet.GetsSystemSecurity(ctx, false)

	if err != nil {
		return "", err
	}

	switch role {
	case enums.UserRoleOrdinary: // 普通用户
		if strategy.SystemSecurityUserLimit == 0 {
			return "", errors.New(`用户名或密码错误`)
		}
		if errNum >= strategy.SystemSecurityUserLimit {
			// 禁用
			cacheClient.Del(ctx, username)
			return "disable", errors.New(fmt.Sprintf(`用户名或密码错误，已错误%d次，账号已被禁用`, errNum))
		}
		return "", errors.New(fmt.Sprintf(`用户名或密码错误，已错误%d次，累计错误%d次后将被禁用`, errNum, strategy.SystemSecurityUserLimit))
	default:
		if strategy.SystemSecurityAdminLimit == 0 || strategy.SystemSecurityBanTime == 0 {
			return "", errors.New(`用户名或密码错误`)
		}
		if errNum >= strategy.SystemSecurityAdminLimit {
			cacheClient.Del(ctx, username)
			// 限制登录时间
			errMsg := fmt.Sprintf(`用户名或密码错误，已错误%d次，账号已被锁定，请于%d分钟后重试`, errNum, strategy.SystemSecurityBanTime)
			cacheClient.Set(ctx, username+"_limit", errMsg, time.Duration(strategy.SystemSecurityBanTime)*time.Minute)
			return "limit", errors.New(errMsg)
		}
		return "", errors.New(fmt.Sprintf(`用户名或密码错误，已错误%d次，累计错误%d次后将锁定%d分钟`, errNum, strategy.SystemSecurityAdminLimit, strategy.SystemSecurityBanTime))
	}
}

// UserManageOP 用户管理 - 创建/修改
func (u *User) UserManageOP(ctx context.Context, uid int, req *typespec.UserManageOpReq, res *typespec.UserManageCreateRes) error {
	var userService services.User
	// 修改
	if req.ID != 0 {
		// 参数校验
		if errStr := userService.CheckUserEditParam(ctx, *req, uid); errStr != "" {
			return errors.New(errStr)
		}
		if err := userService.UpdateUser(ctx, uid, req.ID, req.Role, req.Username, req.Email, req.Department, req.Remark, req.GroupIds); err != nil {
			return err
		}
	} else {
		// 参数校验
		req.Password = services.PasswordDecodeString(req.Password)
		req.Repassword = services.PasswordDecodeString(req.Repassword)
		if errStr := userService.CheckUserAddParam(ctx, *req, uid); errStr != "" {
			return errors.New(errStr)
		}
		if err := userService.CreateUser(ctx, req.Username, services.PasswordEncryption(req.Password), req.AccountExpireTime, req.Email, req.GroupIds, req.Department, req.Remark, uid, req.Role); err != nil {
			return err
		}
	}
	return nil
}

// 用户管理 - 更新
func (u *User) UserManageUpdate(ctx context.Context, createUid int, req *typespec.UserManageUpdateReq, res *typespec.UserManageUpdateRes) error {
	var userService services.User
	if err := userService.UserIsAdmin(ctx, createUid); err != nil {
		return err
	}

	// 检查用户是否存在
	user, err := userService.GetUserForId(ctx, req.Id)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New("未知的用户")
	}

	// 检查用户名是否已存在
	// 检查邮箱是否已存在
	ok, err := userService.CheckUserIsAlready(ctx, req.Id, req.Username, req.Email)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("用户名或邮箱已存在，请确认")
	}

	err = userService.UpdateUser(ctx, createUid, req.Id, req.Role, req.Username, req.Email, req.Department, req.Remark, req.GroupIds)
	if err != nil {
		return err
	}
	return nil
}

// LoginByLianTong 免密登录
func (u *User) LoginByLianTong(ctx context.Context, req *typespec.UserLoginLianTongReq, resp *typespec.UserLoginLianTongResp) error {
	// 第一步 根据授权码获取 access_token
	var clientConfig map[string]map[string]interface{}
	if err := config.Load("client", &clientConfig); err != nil {
		return err
	}
	if clientConfig["service_liantong"] == nil {
		return errors.New("免密登录配置缺失")
	}
	err := httpclients.GetTokenValidLianTong(httpclients.GetTokenValidLianTongReq{
		SsoToken: req.SsoToken,
	})
	if err != nil {
		return err
	}
	//第三步 获取本系统用户信息
	var (
		userSrv         services.User
		logAuditService services.LogAudit //增加审计日志-登录日志
	)
	user, err := userSrv.GetForUsername(ctx, req.UserName)
	if err != nil {
		return err
	}
	// 判断用户信息是否存在,如果用户不存在则创建用户
	if user.ID == 0 {
		accountExpire := time.Now().Add(3 * 365 * 24 * time.Hour).Format(enums.TimeLayout)
		err = userSrv.CreateUser(ctx, req.UserName, services.PasswordEncryption(enums.ResetPassWordInfo), accountExpire, "", "", "", "", 1, enums.UserRoleAdministrator)
		if err != nil {
			return err
		}
		user, err = userSrv.GetForUsername(ctx, req.UserName)
		if err != nil {
			return err
		}
	}
	// 判断账号状态 非成功状态 均返回错误
	if user.Status != enums.UserStatusSuccess {
		if err := userSrv.UpdateUserStatus(ctx, 1, user.ID, enums.UserStatusSuccess); err != nil {
			return err
		}
		// 联通免密登录去掉状态校验这一行
		// return errors.New(enums.UserEnum.GetStatus(user.Status))
	}
	// 生成token
	token, err := userSrv.GenerateToken(user.ID)
	if err != nil {
		return err
	}
	// 生成jwt token 在全局 response中间件内 登录成功
	resp.Token = token
	resp.Uid = user.ID
	resp.Role = user.Type
	resp.Username = user.Username
	// 修改用户状态为在线
	userSrv.UpdateIsAlive(ctx, resp.Uid, enums.UserOnLine, "")
	// 保存到审计日志
	if err = logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeLogin, user.Username+"登录成功", user.Username, ""); err != nil {
		return err
	}
	return nil
}

// LoginBySiYuan 四院免密登录
func (u *User) LoginBySiYuan(ctx context.Context, req *typespec.UserLoginSiYuanReq, resp *typespec.UserLoginSiYuanResp) error {
	// 第一步 根据授权码获取 access_token
	var clientConfig map[string]map[string]interface{}
	if err := config.Load("client", &clientConfig); err != nil {
		return err
	}
	if clientConfig["service_siyuan"] == nil {
		return errors.New("免密登录配置缺失")
	}
	uid, err := httpclients.GetTokenValidSiYuan(httpclients.GetTokenValidSiYuanReq{
		Token:        req.Token,
		PlatformCode: clientConfig["service_siyuan"]["platformCode"].(string),
	})
	if err != nil {
		return err
	}
	//第三步 获取本系统用户信息
	var (
		userSrv         services.User
		logAuditService services.LogAudit //增加审计日志-登录日志
	)
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return err
	}
	log.Println("uidInt", uidInt)
	user, err := userSrv.GetUserForId(ctx, uidInt)
	if err != nil {
		return err
	}
	log.Println("user", user)
	// 生成token
	token, err := userSrv.GenerateToken(user.ID)
	if err != nil {
		return err
	}
	// 生成jwt token 在全局 response中间件内 登录成功
	resp.Token = token
	resp.Uid = user.ID
	resp.Role = user.Type
	resp.Username = user.Username
	log.Println("resp", resp)
	// 修改用户状态为在线
	userSrv.UpdateIsAlive(ctx, resp.Uid, enums.UserOnLine, "")
	// 保存到审计日志
	if err = logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeLogin, user.Username+"登录成功", user.Username, ""); err != nil {
		return err
	}
	return nil
}

// LoginByApiToken 免密登录
func (u *User) LoginByApiToken(ctx context.Context, req *typespec.LoginByApiTokenReq, resp *typespec.LoginByApiTokenResp) error {
	// 第一步 根据授权码获取 access_token
	//第三步 获取本系统用户信息
	var (
		userSrv         services.User
		logAuditService services.LogAudit //增加审计日志-登录日志
	)
	user, err := userSrv.GetUserForToken(ctx, req.Token)
	log.Println("user", user)
	// 生成token
	token, err := userSrv.GenerateToken(user.ID)
	if err != nil {
		return err
	}
	// 生成jwt token 在全局 response中间件内 登录成功
	resp.Token = token
	resp.Uid = user.ID
	resp.Role = user.Type
	resp.Username = user.Username
	log.Println("resp", resp)
	// 修改用户状态为在线
	userSrv.UpdateIsAlive(ctx, resp.Uid, enums.UserOnLine, "")
	// 保存到审计日志
	if err = logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeLogin, user.Username+"登录成功", user.Username, ""); err != nil {
		return err
	}
	return nil
}

// LoginByHTYZ 航天运载免密登录
func (u *User) LoginByHTYZ(ctx context.Context, req *typespec.LoginByHTYZReq, resp *typespec.LoginByHTYZResp) error {
	// 第一步 根据授权码获取 access_token
	var clientConfig map[string]map[string]interface{}
	if err := config.Load("client", &clientConfig); err != nil {
		return err
	}
	//if clientConfig["service_htyz"] == nil {
	//	return errors.New("免密登录配置缺失")
	//}
	//configData := clientConfig["service_htyz"]
	//if configData["k"] == nil {
	//	return errors.New("配置文件缺少 other_login.k")
	//}
	//mySignKeyBytes, err := base64.RawURLEncoding.DecodeString(configData["k"].(string))
	//if err != nil {
	//	return errors.New("三方登录 secret 解析错误")
	//}
	//// 创建Token结构体
	//tokenObj, err := jwtv5.Parse(req.Token, func(token *jwtv5.Token) (interface{}, error) {
	//	return mySignKeyBytes, nil
	//})
	//if err != nil {
	//	if strings.Contains(err.Error(), "token is expired") {
	//		// 过期 token
	//		if configData["check_token_expired"] == "Y" {
	//			// 需要校验
	//			return err
	//		}
	//	} else {
	//		return err
	//	}
	//}
	//data, ok := tokenObj.Claims.(jwtv5.MapClaims)
	//if !ok {
	//	return errors.New("token 解析失败")
	//}
	//payload, err := utils.CheckOtherLoginPayloadParam(data)
	//if err != nil {
	//	return err
	//}

	//第三步 获取本系统用户信息
	var (
		userSrv         services.User
		logAuditService services.LogAudit //增加审计日志-登录日志
	)
	// 检查用户是否存在 不存在则返回失败
	user, _ := userSrv.GetForUsername(ctx, "admin")
	if user.ID == 0 {
		// 系统内没有这个用户
		return errors.New("用户" + "admin" + "不在当前系统")
	}
	log.Println("user", user)
	// 生成token
	token, err := userSrv.GenerateToken(user.ID)
	if err != nil {
		return err
	}
	// 生成jwt token 在全局 response中间件内 登录成功
	resp.Token = token
	resp.Uid = user.ID
	resp.Role = user.Type
	resp.Username = user.Username
	log.Println("resp", resp)
	// 修改用户状态为在线
	userSrv.UpdateIsAlive(ctx, resp.Uid, enums.UserOnLine, "")
	// 保存到审计日志
	if err = logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeLogin, user.Username+"登录成功", user.Username, ""); err != nil {
		return err
	}
	return nil
}

// LoginByGAYS 公安一所免密登录
func (u *User) LoginByGAYS(ctx context.Context, req *typespec.LoginByGAYSReq, resp *typespec.LoginByGAYSResp) error {
	// 第一步 根据授权码获取 access_token
	var clientConfig map[string]map[string]interface{}
	if err := config.Load("client", &clientConfig); err != nil {
		return err
	}
	if clientConfig["service_gays"] == nil {
		return errors.New("免密登录配置缺失")
	}
	username, err := httpclients.GetUserByGays(httpclients.GetUserGaysReq{
		Authorization: req.Token,
	})
	if err != nil {
		return err
	}
	//第三步 获取本系统用户信息
	var (
		userSrv         services.User
		logAuditService services.LogAudit //增加审计日志-登录日志
	)
	user, err := userSrv.GetForUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		accountExpire := time.Now().Add(3 * 365 * 24 * time.Hour).Format(enums.TimeLayout)
		err = userSrv.CreateUser(ctx, username, "", accountExpire, "", "", "", "", 1, enums.UserRoleAdministrator)
		if err != nil {
			return err
		}
		user, err = userSrv.GetForUsername(ctx, username)
		if err != nil {
			return err
		}
	} else {
		userSrv.UpdateUserStatus(ctx, 1, user.ID, enums.UserStatusSuccess)
	}
	log.Println("user", user)
	// 生成token
	token, err := userSrv.GenerateToken(user.ID)
	if err != nil {
		return err
	}
	// 生成jwt token 在全局 response中间件内 登录成功
	resp.Token = token
	resp.Uid = user.ID
	resp.Role = user.Type
	resp.Username = user.Username
	log.Println("resp", resp)
	// 修改用户状态为在线
	userSrv.UpdateIsAlive(ctx, resp.Uid, enums.UserOnLine, "")
	// 保存到审计日志
	if err = logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeLogin, user.Username+"登录成功", user.Username, ""); err != nil {
		return err
	}
	return nil
}

// PasswordCheck 免密登录
func (u *User) PasswordCheck(ctx context.Context, uid any, req *typespec.PasswordCheckReq, resp *typespec.PasswordCheckResp) error {
	var userSrv services.User
	user, err := userSrv.GetUserForId(ctx, uid.(int))
	if err != nil {
		return err
	}
	// 第一步 检查是否为默认口令
	var pbkdf encryption.Pbkdf2Sha256
	if ok, _ := pbkdf.PasswordVerify("admin123321", user.Password); ok {
		return errors.New("请及时更新系统默认口令")
	}
	if ok, _ := pbkdf.PasswordVerify("luckyadmin123", user.Password); ok {
		return errors.New("请及时更新系统默认口令")
	}
	// 第二步 检查密码有效期
	var mapSetService services.MapSet
	passwordValidString, err := mapSetService.GetMapValue(ctx, enums.SystemSecurityPasswordValid)
	passwordValid, err := strconv.Atoi(passwordValidString)
	if err != nil {
		passwordValid = 0
	}
	passwordExpireTime := user.PasswordChangeTime.Add(time.Duration(passwordValid*31*24) * time.Hour)
	today := time.Now().Format(enums.TimeYMDBarLayout)
	if passwordExpireTime.Format(enums.TimeYMDBarLayout) < today {
		return errors.New("密码已过期，请及时更新密码有效期")
	}
	passwordValidAlert := 7 // 在7天内密码会过期将进行提醒
	passwordAlertTime := passwordExpireTime.Add(-time.Duration(passwordValidAlert*24) * time.Hour)
	if passwordAlertTime.Format(enums.TimeYMDBarLayout) <= today {
		return errors.New("密码即将过期，请及时更新密码有效期")
	}
	return nil
}

// LoginBy15Suo 免密登录
func (u *User) LoginBy15Suo(ctx context.Context, req *typespec.UserLogin15SuoReq, resp *typespec.UserLogin15SuoResp) error {
	// 第一步 根据授权码获取 access_token
	var clientConfig map[string]map[string]interface{}
	if err := config.Load("client", &clientConfig); err != nil {
		return err
	}
	if clientConfig["service_15suo"] == nil {
		return errors.New("免密登录配置缺失")
	}
	service15Suo := clientConfig["service_15suo"]
	accessToken, err := httpclients.GetAccessToken(httpclients.GetAccessTokenReq{
		GrantType:    "authorization_code",
		ClientId:     service15Suo["client_id"].(string),
		ClientSecret: service15Suo["client_secret"].(string),
		Code:         req.Code,
	})
	if err != nil {
		return err
	}
	fmt.Println("55555555555555: ", accessToken)
	// 第二步 根据access_token获取授权平台用户信息
	res, err := httpclients.GetUserByAccessToken(httpclients.GetUser15suoReq{AccessToken: accessToken})
	if err != nil {
		return err
	}
	//第三步 获取本系统用户信息
	var (
		userSrv         services.User
		logAuditService services.LogAudit //增加审计日志-登录日志
	)
	user, err := userSrv.GetForUsername(ctx, res.Data.UserName)
	if err != nil {
		return err
	}
	// 判断用户信息是否存在,如果用户不存在则创建用户
	if user.ID == 0 {
		accountExpire := time.Now().Add(3 * 365 * 24 * time.Hour).Format(enums.TimeLayout)
		err = userSrv.CreateUser(ctx, res.Data.UserName, "", accountExpire, "", "", "", "", 1, enums.UserRoleAdministrator)
		if err != nil {
			return err
		}
		user, err = userSrv.GetForUsername(ctx, res.Data.UserName)
		if err != nil {
			return err
		}
	}
	// 判断账号状态 非成功状态 均返回错误
	if user.Status != enums.UserStatusSuccess {
		return errors.New(enums.UserEnum.GetStatus(user.Status))
	}

	// 生成token
	token, err := userSrv.GenerateToken(user.ID)
	if err != nil {
		return err
	}
	// 生成jwt token 在全局 response中间件内 登录成功
	resp.Token = token
	resp.Uid = user.ID
	resp.Role = user.Type
	resp.Username = user.Username
	// 修改用户状态为在线
	userSrv.UpdateIsAlive(ctx, resp.Uid, enums.UserOnLine, "")
	// 保存到审计日志
	if err = logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeLogin, user.Username+"登录成功", user.Username, ""); err != nil {
		return err
	}
	return nil
}
