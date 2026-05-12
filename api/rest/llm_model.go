package rest

import (
	"context"
	"smart/api/typespec"
	"smart/application"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

// LlmModelList 大模型列表
func LlmModelList(c *gin.Context) {
	var (
		req typespec.LlmModelListReq
		res typespec.LlmModelListRes
		app application.LlmModel
	)
	ctx := server.NewContext(context.Background(), c)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("LlmModelList parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	if err := app.LlmModelList(ctx, &req, &res); err != nil {
		log.Error("LlmModelList parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LlmModelDetail 大模型详情
func LlmModelDetail(c *gin.Context) {
	var (
		req typespec.LlmModelDetailReq
		res typespec.LlmModelDetailRes
		app application.LlmModel
	)
	ctx := server.NewContext(context.Background(), c)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("LlmModelDetail parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	if err := app.LlmModelDetail(ctx, &req, &res); err != nil {
		log.Error("LlmModelDetail parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LlmModelSave 保存大模型（添加/编辑）
func LlmModelSave(c *gin.Context) {
	var (
		req typespec.LlmModelSaveReq
		res typespec.LlmModelSaveRes
		app application.LlmModel
	)
	ctx := server.NewContext(context.Background(), c)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("LlmModelSave parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	if err := app.LlmModelSave(ctx, &req, &res); err != nil {
		log.Error("LlmModelSave parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LlmModelDelete 删除大模型
func LlmModelDelete(c *gin.Context) {
	var (
		req typespec.LlmModelDeleteReq
		res typespec.LlmModelDeleteRes
		app application.LlmModel
	)
	ctx := server.NewContext(context.Background(), c)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("LlmModelDelete parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	if err := app.LlmModelDelete(ctx, &req, &res); err != nil {
		log.Error("LlmModelDelete parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LlmModelSetDefault 设置默认大模型
func LlmModelSetDefault(c *gin.Context) {
	var (
		req typespec.LlmModelSetDefaultReq
		res typespec.LlmModelSetDefaultRes
		app application.LlmModel
	)
	ctx := server.NewContext(context.Background(), c)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("LlmModelSetDefault parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	if err := app.LlmModelSetDefault(ctx, &req, &res); err != nil {
		log.Error("LlmModelSetDefault parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LlmModelEnabledTest 测试大模型是否可用
func LlmModelEnabledTest(c *gin.Context) {
	var (
		req typespec.LlmModelEnabledTestReq
		res typespec.LlmModelEnabledTestRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("LlmModelEnabledTest parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var (
		ctx = server.NewContext(context.Background(), c)
		app application.LlmModel
	)
	if err := app.LlmModelEnabledTest(ctx, &req, &res); err != nil {
		log.Error("LlmModelSetDefault error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// LlmModelEnums llm模型枚举
func LlmModelEnums(c *gin.Context) {
	var (
		res typespec.LlmModelEnumsResp
		app application.LlmModel
	)
	app.LlmModelEnums(&res)
	server.RespSuccess(c, res)
}

// LlmModelEnhancementDetail llm模型增强详情
func LlmModelEnhancementDetail(c *gin.Context) {
	var (
		resp typespec.LlmModelEnhancementDetailResp
		app  application.LlmModel
		ctx  = server.NewContext(context.Background(), c)
	)
	if err := app.LlmModelEnhancementDetail(ctx, &resp); err != nil {
		log.Error("LlmModelEnhancementDetail error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// LlmModelEnhancementEdit llm模型增强切换
func LlmModelEnhancementEdit(c *gin.Context) {
	var (
		req typespec.LlmModelEnhancementEditReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("LlmModelEnhancementEdit parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	var (
		app application.LlmModel
		ctx = server.NewContext(context.Background(), c)
	)
	if err := app.LlmModelEnhancementEdit(ctx, &req); err != nil {
		log.Error("LlmModelEnhancementEdit error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	server.RespSuccess(c, nil)
}
