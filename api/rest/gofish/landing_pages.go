package gofish

import (
	"smart/api/typespec/gofishtypespec"
	"smart/application/gofish"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
)

// PageAll 获取全部着陆页
func (api *GoPhishAPI) PageAll(c *gin.Context) {
	var req gofishtypespec.GetListInfoReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.PageGetAllBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// PageDetail 着陆页详情
func (api *GoPhishAPI) PageDetail(c *gin.Context) {
	var req gofishtypespec.PageIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.PageGetDetailBiz(c, req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// PageCreate 创建着陆页
func (api *GoPhishAPI) PageCreate(c *gin.Context) {
	var req gofishtypespec.CreatePageReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.PageCreateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// PageUpdate 更新着陆页
func (api *GoPhishAPI) PageUpdate(c *gin.Context) {
	var req gofishtypespec.UpdatePageReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.PageUpdateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// PageDelete 删除着陆页
func (api *GoPhishAPI) PageDelete(c *gin.Context) {
	var req gofishtypespec.PageIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	if err := app.PageDeleteBiz(c, req.ID); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ok")
}

// PageImportSite 导入网站
func (api *GoPhishAPI) PageImportSite(c *gin.Context) {
	var req gofishtypespec.ImportSiteReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.PageImportSiteBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
