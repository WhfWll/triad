package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
	"smart/tools/utils"
)

// UserDetail 用户详情
func UserDetail(c *gin.Context) {
	var (
		req typespec.UserDetailReq
		res typespec.UserListInfo
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.UserDetail(ctx, &req, &res); err != nil {
		log.Error("UserDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// UserList 用户列表
func UserList(c *gin.Context) {
	var (
		req typespec.UserListReq
		res typespec.UserListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.UserList(ctx, &req, &res); err != nil {
		log.Error("UserList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// UserEnumList 用户枚举列表
func UserEnumList(c *gin.Context) {
	var (
		req typespec.UserEnumListReq
		res typespec.UserEnumListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserEnumList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.UserEnumList(ctx, &req, &res); err != nil {
		log.Error("UserEnumList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// UserLoginCaptcha 登陆验证码
func UserLoginCaptcha(c *gin.Context) {
	var (
		req typespec.UserLoginCaptchaReq
		res typespec.UserLoginCaptchaRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLoginCaptcha parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)

	var app application.User
	if err := app.UserLoginCaptcha(ctx, &req, &res); err != nil {
		log.Error("UserLoginCaptcha parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// UserLogin 登陆
func UserLogin(c *gin.Context) {
	var (
		req typespec.UserLoginReq
		res typespec.UserLoginRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	req.IP = utils.GetClientIp(c)
	if err := app.UserLogin(ctx, &req, &res); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	c.Set("uid", res.Uid)
	server.RespSuccess(c, res)
}

// UserManageOp 用户管理 - 用户创建/更新
func UserManageOp(c *gin.Context) {
	var (
		req typespec.UserManageOpReq
		res typespec.UserManageCreateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserManageOp parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 获取操作用户的ID
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("UserManageOp parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.UserManageOP(ctx, uid.(int), &req, &res); err != nil {
		log.Error("UserManageOp parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "操作成功")
}

// DelUser 删除用户
func DelUser(c *gin.Context) {
	var (
		req typespec.UserDelReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserManageOp parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 获取操作用户的ID
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("UserManageOp parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.DelUser(ctx, uid.(int), &req); err != nil {
		log.Error("UserManageOp parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "UserManageOp success")
}

// UpdatePassWord 修改密码
func UpdatePassWord(c *gin.Context) {
	var (
		req typespec.UserPassPWReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UpdatePassWord parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 获取操作用户的ID
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("UpdatePassWord parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.UpdatePassWord(ctx, uid.(int), &req); err != nil {
		log.Error("UpdatePassWord parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "UpdatePassWord success")
}

// ResetPassWord 重置密码
func ResetPassWord(c *gin.Context) {
	var (
		req typespec.UserResetPassPWReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ResetPassWord parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 获取操作用户的ID
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("ResetPassWord parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.ResetPassWord(ctx, uid.(int), &req); err != nil {
		log.Error("ResetPassWord parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ResetPassWord success")
}

// UpdateUserExp 修改过期时间
func UpdateUserExp(c *gin.Context) {
	var (
		req typespec.UpdateUserExpReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UpdateUserExp parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 获取操作用户的ID
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("UpdateUserExp parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.UpdateUserExp(ctx, uid.(int), &req); err != nil {
		log.Error("UpdateUserExp parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "UpdateUserExp success")
}

// ChangeUserStatus 修改用户状态
func ChangeUserStatus(c *gin.Context) {
	var (
		req typespec.ChangeUserStatusReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ChangeUserStatus parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 获取操作用户的ID
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("ChangeUserStatus parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.ChangeUserStatus(ctx, uid.(int), &req); err != nil {
		log.Error("ChangeUserStatus parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ChangeUserStatus success")
}

// LogOut 推出
func LogOut(c *gin.Context) {
	var (
		req typespec.LogOutReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("LogOut parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 获取操作用户的ID
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("LogOut parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.Logout(ctx, uid.(int)); err != nil {
		log.Error("LogOut parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "LogOut success")
}

// LoginByLianTong 联通靶场导调免密登录
func LoginByLianTong(c *gin.Context) {
	var (
		req typespec.UserLoginLianTongReq
		res typespec.UserLoginLianTongResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.LoginByLianTong(ctx, &req, &res); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	c.Set("uid", res.Uid)
	server.RespSuccess(c, res)
}

// LoginBySiYuan 四院免密登录项目
func LoginBySiYuan(c *gin.Context) {
	var (
		req typespec.UserLoginSiYuanReq
		res typespec.UserLoginSiYuanResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.LoginBySiYuan(ctx, &req, &res); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	c.Set("uid", res.Uid)
	server.RespSuccess(c, res)
}

// LoginByApiToken 中测yakit项目登录
func LoginByApiToken(c *gin.Context) {
	var (
		req typespec.LoginByApiTokenReq
		res typespec.LoginByApiTokenResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.LoginByApiToken(ctx, &req, &res); err != nil {
		log.Error("LoginByApiToken parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	c.Set("uid", res.Uid)
	server.RespSuccess(c, res)
}

// LoginByHTYZ 航天运载项目登录
func LoginByHTYZ(c *gin.Context) {
	var (
		req typespec.LoginByHTYZReq
		res typespec.LoginByHTYZResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.LoginByHTYZ(ctx, &req, &res); err != nil {
		log.Error("LoginByHTYZ parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	c.Set("uid", res.Uid)
	server.RespSuccess(c, res)
}

// LoginByGAYS 公安一所免密登录
func LoginByGAYS(c *gin.Context) {
	var (
		req typespec.LoginByGAYSReq
		res typespec.LoginByGAYSResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.LoginByGAYS(ctx, &req, &res); err != nil {
		log.Error("LoginByHTYZ parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	c.Set("uid", res.Uid)
	server.RespSuccess(c, res)
}

// PasswordCheck 密码检查
func PasswordCheck(c *gin.Context) {
	var (
		req typespec.PasswordCheckReq
		res typespec.PasswordCheckResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.PasswordCheck(ctx, req.UserId, &req, &res); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// LoginBy15Suo 15所项目免密登录接口
func LoginBy15Suo(c *gin.Context) {
	var (
		req typespec.UserLogin15SuoReq
		res typespec.UserLogin15SuoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.User
	if err := app.LoginBy15Suo(ctx, &req, &res); err != nil {
		log.Error("UserLogin parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	c.Set("uid", res.Uid)
	server.RespSuccess(c, res)
}
