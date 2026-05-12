package rest

import (
	"context"
	"smart/api/typespec"
	"smart/application"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

// TaskInfoStat 首页 - 任务信息统计
func TaskInfoStat(c *gin.Context) {
	var (
		req typespec.TaskInfoStatReq
		res typespec.TaskInfoStatRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskInfoStat parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	// ctx := server.NewContext(context.Background(), c)
	var app application.Homepage
	app.TaskInfoStat(c, &req, &res)

	server.RespSuccess(c, res)
}

// VulEvidenceStat 首页 - 漏洞取证信息统计
func VulEvidenceStat(c *gin.Context) {
	var (
		res typespec.VulEvidenceStatRes
		// ctx = server.NewContext(context.Background(), c)
		app application.Homepage
	)
	app.VulEvidenceStat(c, &res)

	server.RespSuccess(c, res)
}

// TargetRiskStat 首页 - 目标风险统计
func TargetRiskStat(c *gin.Context) {
	var res typespec.TargetRiskStatRes
	// ctx := server.NewContext(context.Background(), c)
	var app application.Homepage
	app.TargetRiskStat(c, &res)

	server.RespSuccess(c, res)
}

// TaskVulRiskStat 首页 - 任务漏洞风险统计
func TaskVulRiskStat(c *gin.Context) {
	var res typespec.TaskVulRiskStatRes
	// ctx := server.NewContext(context.Background(), c)
	var app application.Homepage
	app.TaskVulRiskStat(c, &res)

	server.RespSuccess(c, res)
}

// ToolInfoStat 首页 - 工具信息数量统计
func ToolInfoStat(c *gin.Context) {
	var (
		res typespec.ToolInfoStatRes
		app application.Homepage
		ctx = server.NewContext(context.Background(), c)
	)

	if err := app.ToolInfoStat(ctx, &res); err != nil {
		log.Error("ToolInfoStat parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// VulTypeStat 首页 - 漏洞类型统计
func VulTypeStat(c *gin.Context) {
	var (
		req typespec.VulTypeStatReq
		res typespec.VulTypeStatRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulTypeStat parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	// ctx := server.NewContext(context.Background(), c)
	var app application.Homepage
	if err := app.VulTypeStat(c, &req, &res); err != nil {
		log.Error("VulTypeStat parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// VulFindTrendStat 首页 - 漏洞发现趋势统计
func VulFindTrendStat(c *gin.Context) {
	var (
		req typespec.VulFindTrendStatReq
		res typespec.VulFindTrendStatRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulFindTrendStat parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	// ctx := server.NewContext(context.Background(), c)
	var app application.Homepage
	if err := app.VulFindTrendStat(c, &req, &res); err != nil {
		log.Error("VulTypeStat parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// MessageStat 首页 - 最新消息统计模块
func MessageStat(c *gin.Context) {
	var (
		res typespec.MessageStatRes
		app application.Homepage
		// ctx = server.NewContext(context.Background(), c)
	)
	app.MessageStat(c, &res)

	server.RespSuccess(c, res)
}
