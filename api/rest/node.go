package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// 节点管理 - 是否启用分布式 - 要求，必须有一个存活节点，如节点挂掉，自动关闭分布式
func SystemNodeSetDistribute(c *gin.Context) {
	var (
		req  typespec.NodeSetDistributeReq
		resp typespec.NodeSetDistributeRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeSetDistribute parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeSetDistribute(ctx, &req); err != nil {
		log.Error("SystemNodeSetDistribute parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 节点管理 - 是否启用分布式 - 获取节点状态
func SystemNodeGetDistribute(c *gin.Context) {
	var (
		resp typespec.NodeIsDistributeRes
	)

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeGetDistribute(ctx, &resp); err != nil {
		log.Error("SystemNodeGetDistribute parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 节点管理 - 新增节点
func SystemNodeAdd(c *gin.Context) {
	var (
		req  typespec.NodeAddReq
		resp typespec.NodeAddRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeAdd parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeAdd(ctx, &req, &resp); err != nil {
		log.Error("SystemNodeAdd parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 节点管理 - 编辑节点
func SystemNodeEdit(c *gin.Context) {
	var (
		req  typespec.NodeEditReq
		resp typespec.NodeEditRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeEdit(ctx, &req, &resp); err != nil {
		log.Error("SystemNodeEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 节点管理 - 节点详情
func SystemNodeInfo(c *gin.Context) {
	var (
		req  typespec.NodeInfoReq
		resp typespec.NodeInfoRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeInfo(ctx, &req, &resp); err != nil {
		log.Error("SystemNodeInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 节点管理 - 节点列表
func SystemNodeList(c *gin.Context) {
	var (
		req  typespec.NodeListReq
		resp typespec.NodeListRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeList(ctx, &req, &resp); err != nil {
		log.Error("SystemNodeList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 节点管理 - 删除节点
func SystemNodeDel(c *gin.Context) {
	var (
		req  typespec.NodeDelReq
		resp typespec.NodeDelRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeDel(ctx, &req, &resp); err != nil {
		log.Error("SystemNodeDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 节点管理 - 禁用｜启用节点
func SystemNodeDisOrEnable(c *gin.Context) {
	var (
		req  typespec.NodeDisOrEnableReq
		resp typespec.NodeDisOrEnableRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeDisOrEnable parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeDisOrEnable(ctx, &req, &resp); err != nil {
		log.Error("SystemNodeDisOrEnable parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 节点管理 - 所有可用节点
func SystemNodeAllEnable(c *gin.Context) {
	var (
		req  typespec.NodeAllEnableReq
		resp typespec.NodeAllEnableRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeAllEnable parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.Node
	if err := app.SystemNodeAllEnable(ctx, &req, &resp); err != nil {
		log.Error("SystemNodeAllEnable parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
