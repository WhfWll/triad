package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// ReportEnum 报告枚举
func ReportEnum(c *gin.Context) {
	var resp typespec.ReportEnumResp
	ctx := server.NewContext(context.Background(), c)
	var app application.Report
	if err := app.ReportEnum(ctx, &resp); err != nil {
		log.Error("ReportEnum parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportList 报告列表及筛选
func ReportList(c *gin.Context) {
	var (
		req  typespec.ReportListReq
		resp typespec.ReportListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//ctx := server.NewContext(context.Background(), c)
	var app application.Report
	if err := app.ReportList(c, &req, &resp); err != nil {
		log.Error("ReportList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportDownload 报告下载
func ReportDownload(c *gin.Context) {
	var (
		req  typespec.ReportDownloadReq
		resp typespec.ReportDownloadResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportDownload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Report
	if err := app.ReportDownload(ctx, &req, &resp); err != nil {
		log.Error("ReportDownload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportDel 报告删除
func ReportDel(c *gin.Context) {
	var (
		req  typespec.ReportDelReq
		resp typespec.ReportDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Report
	if err := app.ReportDel(ctx, &req); err != nil {
		log.Error("ReportDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportSave 生成报告
func ReportSave(c *gin.Context) {
	var (
		req  typespec.ReportSaveReq
		resp typespec.ReportSaveResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Report
	if err := app.ReportSave(ctx, &req); err != nil {
		log.Error("ReportSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
