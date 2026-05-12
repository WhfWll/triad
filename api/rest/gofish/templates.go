package gofish

import (
	"smart/api/typespec/gofishtypespec"
	"smart/application/gofish"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
)

// TemplateAll 获取全部模板
func (api *GoPhishAPI) TemplateAll(c *gin.Context) {
	var req gofishtypespec.GetListInfoReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.TemplateGetAllBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TemplateDetail 模板详情
func (api *GoPhishAPI) TemplateDetail(c *gin.Context) {
	var req gofishtypespec.TemplateIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.TemplateGetDetailBiz(c, req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TemplateCreate 创建模板
func (api *GoPhishAPI) TemplateCreate(c *gin.Context) {
	var req gofishtypespec.CreateTemplateReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.TemplateCreateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TemplateUpdate 更新模板
func (api *GoPhishAPI) TemplateUpdate(c *gin.Context) {
	var req gofishtypespec.UpdateTemplateReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.TemplateUpdateBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TemplateDelete 删除模板
func (api *GoPhishAPI) TemplateDelete(c *gin.Context) {
	var req gofishtypespec.TemplateIDReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	if err := app.TemplateDeleteBiz(c, req.ID); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "ok")
}

// TemplateImportEmail 导入邮件
func (api *GoPhishAPI) TemplateImportEmail(c *gin.Context) {
	var req gofishtypespec.ImportEmailReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var app gofish.GoPhishBiz
	resp, err := app.TemplateImportEmailBiz(c, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
