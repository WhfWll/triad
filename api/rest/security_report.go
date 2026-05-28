package rest

import (
	"context"
	"smart/api/typespec"
	"smart/application"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

func SecurityReportGenerate(c *gin.Context) {
	var req typespec.SecurityReportGenerateReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	uid := 0
	if v := ctx.Value("uid"); v != nil {
		uid = v.(int)
	}
	var app application.SecurityReportApp
	resp, err := app.Generate(ctx, &req, uid)
	if err != nil {
		log.Errorf("SecurityReportGenerate failed: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "生成了安全检查报告，模块: %s，任务ID: %d，报告ID: %d", req.Module, req.TaskID, resp.ID)
	server.RespSuccess(c, resp)
}

func SecurityReportList(c *gin.Context) {
	var req typespec.SecurityReportListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SecurityReportApp
	resp, err := app.List(ctx, &req)
	if err != nil {
		log.Errorf("SecurityReportList failed: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func SecurityReportDetail(c *gin.Context) {
	var req typespec.SecurityReportDetailReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SecurityReportApp
	resp, err := app.Detail(ctx, &req)
	if err != nil {
		log.Errorf("SecurityReportDetail failed: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "查看了安全检查报告详情，报告ID: %d", req.ID)
	server.RespSuccess(c, resp)
}

func SecurityReportDelete(c *gin.Context) {
	var req typespec.SecurityReportDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SecurityReportApp
	if err := app.Delete(ctx, &req); err != nil {
		log.Errorf("SecurityReportDelete failed: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "删除了安全检查报告，报告ID: %d", req.ID)
	server.RespSuccess(c, gin.H{})
}
