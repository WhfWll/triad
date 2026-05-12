package gofish

import (
	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec/gofishtypespec"
	"smart/application/gofish"
)

// GetGroupAll .
func (api *GoPhishAPI) GetGroupAll(c *gin.Context) {
	var req gofishtypespec.GetListInfoReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}

	var app gofish.GoPhishBiz
	resp, err := app.GetGroupAllBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// GetGroupDetail .
func (api *GoPhishAPI) GetGroupDetail(c *gin.Context) {
	var app gofish.GoPhishBiz
	var req gofishtypespec.GroupDetailReq

	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	resp, err := app.GetGroupDetailBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, resp)
}

// GroupCreate .
func (api *GoPhishAPI) GroupCreate(c *gin.Context) {
	var app gofish.GoPhishBiz
	var req gofishtypespec.CreateGroupReq

	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	resp, err := app.GroupCreateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, resp)
}

// GroupUpdate .
func (api *GoPhishAPI) GroupUpdate(c *gin.Context) {
	var app gofish.GoPhishBiz
	var req gofishtypespec.UpdateGroupReq

	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	resp, err := app.GroupUpdateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, resp)
}

// GroupDelete .
func (api *GoPhishAPI) GroupDelete(c *gin.Context) {
	var app gofish.GoPhishBiz
	var req gofishtypespec.GroupDetailReq

	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	err := app.GroupDeleteBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, "ok")
}
