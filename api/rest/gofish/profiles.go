package gofish

import (
	"smart/api/typespec/gofishtypespec"
	"smart/application/gofish"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
)

// ProfileAll 获取全部发送配置
func (api *GoPhishAPI) ProfileAll(c *gin.Context) {
	var req gofishtypespec.GetListInfoReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}

	var app gofish.GoPhishBiz
	resp, err := app.ProfileGetAllBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ProfileDetail 详情
func (api *GoPhishAPI) ProfileDetail(c *gin.Context) {
	var req gofishtypespec.ProfileIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.ProfileGetDetailBiz(c, req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ProfileCreate 创建
func (api *GoPhishAPI) ProfileCreate(c *gin.Context) {
	var req gofishtypespec.CreateProfileReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.ProfileCreateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ProfileUpdate 更新
func (api *GoPhishAPI) ProfileUpdate(c *gin.Context) {
	var req gofishtypespec.UpdateProfileReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.ProfileUpdateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ProfileDelete 删除
func (api *GoPhishAPI) ProfileDelete(c *gin.Context) {
	var req gofishtypespec.ProfileIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	if err := app.ProfileDeleteBiz(c, req.ID); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ok")
}

// ProfileSendTestEmail 邮件测试
func (api *GoPhishAPI) ProfileSendTestEmail(c *gin.Context) {
	var req gofishtypespec.SendTestEmailReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	if err := app.ProfileSendTestEmailBiz(c, &req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ok")
}
