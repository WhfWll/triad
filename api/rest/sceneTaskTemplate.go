package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// 场景枚举
func SceneEnums(c *gin.Context) {
	var (
		res typespec.SceneEnumsRes
	)

	var app application.SceneTaskTemplate
	if err := app.SceneEnums(&res); err != nil {
		log.Error("SceneEnums parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SceneTaskTemplateList 任务场景列表
func SceneTaskTemplateList(c *gin.Context) {
	var (
		req typespec.SceneTaskTemplateListReq
		res typespec.SceneTaskTemplateListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SceneTaskTemplateList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SceneTaskTemplate
	if err := app.List(ctx, &req, &res); err != nil {
		log.Error("SceneTaskTemplateList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SceneTaskTemplateSave 创建更新
func SceneTaskTemplateSave(c *gin.Context) {
	var (
		req typespec.SceneTaskTemplateCreateReq
		res typespec.SceneTaskTemplateCreateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SceneTaskTemplateCreateOrUpdate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	/*
		if len(req.Config.VulIdsConfig) == 0 {
			server.RespFail(c, 4000, "参数错误,错误原因：请至少选择一个漏洞")
			return
		}
	*/
	var app application.SceneTaskTemplate
	if err := app.CreateOrUpdate(ctx, &req, &res); err != nil {
		log.Error("SceneTaskTemplateCreateOrUpdate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SceneTaskTemplateDetail 任务场景详情
func SceneTaskTemplateDetail(c *gin.Context) {
	var (
		req typespec.SceneTaskTemplateDetailReq
		res typespec.SceneTaskTemplateDetailRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SceneTaskTemplateDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SceneTaskTemplate
	if err := app.Detail(ctx, &req, &res); err != nil {
		log.Error("SceneTaskTemplateDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SceneTaskTemplateCopy 拷贝
func SceneTaskTemplateCopy(c *gin.Context) {
	var (
		req typespec.SceneTaskTemplateCopyReq
		res typespec.SceneTaskTemplateCopyRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SceneTaskTemplateCopy parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SceneTaskTemplate
	if err := app.Copy(ctx, &req, &res); err != nil {
		log.Error("SceneTaskTemplateCopy parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SceneTaskTemplateSetDefault 设置默认
func SceneTaskTemplateSetDefault(c *gin.Context) {
	var (
		req typespec.SceneTaskTemplateSetDefaultReq
		res typespec.SceneTaskTemplateSetDefaultRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SceneTaskTemplateSetDefault parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SceneTaskTemplate
	if err := app.SetDefault(ctx, &req); err != nil {
		log.Error("SceneTaskTemplateSetDefault parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)

}

// 删除
func SceneTaskTemplateDel(c *gin.Context) {
	var (
		req typespec.SceneTaskTemplateDelReq
		res typespec.SceneTaskTemplateDelRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SceneTaskTemplateDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SceneTaskTemplate
	if err := app.Del(ctx, &req); err != nil {
		log.Error("SceneTaskTemplateDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 任务中的模板下拉菜单
func SceneTaskTemplateToTaskTemplateOptions(c *gin.Context) {
	var (
		res typespec.SceneTaskTemplateToCheckTaskOptionsRes
	)
	ctx := server.NewContext(context.Background(), c)
	var app application.SceneTaskTemplate
	if err := app.ToTaskTemplateOptions(ctx, &res); err != nil {
		log.Error("SceneTaskTemplateDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// Graph 知识图谱显示场景
func Graph(c *gin.Context) {
	var (
		req typespec.GraphReq
		res typespec.GraphRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("Graph parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SceneTaskTemplate
	if err := app.Graph(ctx, &req, &res); err != nil {
		log.Error("SceneTaskTemplateDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}
