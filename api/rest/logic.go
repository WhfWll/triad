package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
	"smart/tools/errno"
)

// 逻辑漏洞 - 任务新建
func LogicTaskCreate(c *gin.Context) {
	var (
		req  typespec.LogicTaskCreateReq
		resp typespec.LogicTaskCreateResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicTaskCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.TaskCreate(c, &req, &resp)
	if err != nil {
		log.Error("LogicTaskCreate application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 任务结束
func LogicTaskStop(c *gin.Context) {
	var (
		req  typespec.LogicTaskStopReq
		resp typespec.LogicTaskStopResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicTaskCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.TaskStop(c, &req, &resp)
	if err != nil {
		log.Error("LogicTaskStop application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 任务删除
func LogicTaskDel(c *gin.Context) {
	var (
		req  typespec.LogicTaskDelReq
		resp typespec.LogicTaskDelResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicTaskDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.TaskDelete(c, &req, &resp)
	if err != nil {
		log.Error("LogicTaskDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 任务列表
func LogicTaskList(c *gin.Context) {
	var (
		req  typespec.LogicTaskListReq
		resp typespec.LogicTaskListResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.TaskList(c, &req, &resp)
	if err != nil {
		log.Error("LogicTaskList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 目标列表
func LogicTargetList(c *gin.Context) {
	var (
		req  typespec.LogicTargetListReq
		resp typespec.LogicTargetListResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.TargetList(c, &req, &resp)
	if err != nil {
		log.Error("LogicTargetList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 漏洞列表
func LogicVulList(c *gin.Context) {
	var (
		req  typespec.LogicVulListReq
		resp typespec.LogicVulListResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicVulList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.VulList(c, &req, &resp)
	if err != nil {
		log.Error("LogicVulList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 漏洞列表
func LogicLogList(c *gin.Context) {
	var (
		req  typespec.LogicLogListReq
		resp typespec.LogicLogListResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicLogList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.LogList(c, &req, &resp)
	if err != nil {
		log.Error("LogicLogList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 漏洞列表
func LogicVulInfo(c *gin.Context) {
	var (
		req  typespec.LogicVulInfoReq
		resp typespec.LogicVulInfoResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicVulInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.VulInfo(c, &req, &resp)
	if err != nil {
		log.Error("LogicVulInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 漏洞列表
func LogicLogInfo(c *gin.Context) {
	var (
		req  typespec.LogicLogInfoReq
		resp typespec.LogicLogInfoResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicLogInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.LogInfo(c, &req, &resp)
	if err != nil {
		log.Error("LogicLogInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 枚举接口
func LogicEnum(c *gin.Context) {
	var resp typespec.LogicEnumResp

	var app application.Logic
	err := app.LogicEnum(c, &resp)
	if err != nil {
		log.Error("LogicEnum application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 漏洞删除
func LogicVulDelete(c *gin.Context) {
	var (
		req  typespec.LogicVulDeleteReq
		resp typespec.LogicVulDeleteResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("VulDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}
	var app application.Logic
	err = app.VulDelete(c, &req, &resp)
	if err != nil {
		log.Error("LogicEnum application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 报文测试
func LogicVulTest(c *gin.Context) {
	var (
		req  typespec.LogicVulTestReq
		resp typespec.LogicVulTestResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("VulDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}
	var app application.Logic
	err = app.VulTest(c, &req, &resp)
	if err != nil {
		log.Error("LogicEnum application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "目标站点访问异常")
		return
	}
	server.RespSuccess(c, resp)
}

// 逻辑漏洞 - 任务信息
func LogicTaskCopy(c *gin.Context) {
	var (
		req  typespec.LogicTaskCopyReq
		resp typespec.LogicTaskCopyResp
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogicTaskCopy parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.Logic
	err = app.TaskCopy(c, &req, &resp)
	if err != nil {
		log.Error("LogicTaskCopy application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// LogicReportSave 生成报告
func LogicReportSave(c *gin.Context) {
	var (
		req  typespec.LogicReportSaveReq
		resp typespec.LogicReportSaveResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Logic
	if err := app.ReportSave(ctx, &req); err != nil {
		log.Error("ReportSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// LogicFlowBaseList 逻辑漏洞流量信息
func LogicFlowBaseList(c *gin.Context) {
	var (
		req  typespec.LogicFlowBaseListReq
		resp typespec.LogicFlowBaseListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Logic
	if err := app.LogicFlowBaseList(ctx, &req, &resp); err != nil {
		log.Error("ReportSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// LogicFlowBaseInfo 逻辑漏洞流量详情
func LogicFlowBaseInfo(c *gin.Context) {
	var (
		req  typespec.LogicFlowBaseInfoReq
		resp typespec.LogicFlowBaseInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("LogicFlowBaseInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Logic
	if err := app.LogicFlowBaseInfo(ctx, &req, &resp); err != nil {
		log.Error("FlowBaseInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// LogicFlowBaseExport 逻辑漏洞流量导出
func LogicFlowBaseExport(c *gin.Context) {
	var (
		req  typespec.LogicFlowBaseExportReq
		resp typespec.LogicFlowBaseExportResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("LogicFlowBaseExport parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Logic
	if err := app.LogicFlowBaseExport(ctx, &req, &resp); err != nil {
		log.Error("FlowBaseInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
