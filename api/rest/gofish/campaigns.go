package gofish

import (
	"smart/api/typespec/gofishtypespec"
	"smart/application/gofish"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
)

// CampaignAll 获取全部活动
func (api *GoPhishAPI) CampaignAll(c *gin.Context) {
	var req gofishtypespec.CampaignListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.CampaignGetAllBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// CampaignDetail 获取活动详情
func (api *GoPhishAPI) CampaignDetail(c *gin.Context) {
	var req gofishtypespec.CampaignIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.CampaignGetDetailBiz(c, req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// CampaignCreate 创建活动
func (api *GoPhishAPI) CampaignCreate(c *gin.Context) {
	var req gofishtypespec.CreateCampaignReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.CampaignCreateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// CampaignLaunch 启动活动
func (api *GoPhishAPI) CampaignLaunch(c *gin.Context) {
	var req gofishtypespec.CampaignIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	if err := app.CampaignLaunchBiz(c, req.ID); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ok")
}

// CampaignComplete 完成活动
func (api *GoPhishAPI) CampaignComplete(c *gin.Context) {
	var req gofishtypespec.CampaignIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	if err := app.CampaignCompleteBiz(c, req.ID); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ok")
}

// CampaignUpdate 更新活动
func (api *GoPhishAPI) CampaignUpdate(c *gin.Context) {
	var req gofishtypespec.UpdateCampaignReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.CampaignUpdateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// CampaignDelete 删除活动
func (api *GoPhishAPI) CampaignDelete(c *gin.Context) {
	var req gofishtypespec.DeleteByIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	if err := app.CampaignDeleteBiz(c, req.ID); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ok")
}

// CampaignResult 活动结果
func (api *GoPhishAPI) CampaignResult(c *gin.Context) {
	var req gofishtypespec.CampaignIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.CampaignResultBiz(c, req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
