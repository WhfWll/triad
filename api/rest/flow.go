package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// FlowTaskEnum 流量分析任务枚举
func FlowTaskEnum(c *gin.Context) {
	var resp typespec.FlowTaskEnumResp
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowTaskEnum(ctx, &resp); err != nil {
		log.Error("FlowTaskEnum application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowTaskList 流量分析任务列表
func FlowTaskList(c *gin.Context) {
	var (
		req  typespec.FlowTaskListReq
		resp typespec.FlowTaskListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowTaskList(c, &req, &resp); err != nil {
		log.Error("FlowTaskList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowTaskDel 流量分析任务删除
func FlowTaskDel(c *gin.Context) {
	var (
		req  typespec.FlowTaskDelReq
		resp typespec.FlowTaskDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowTaskDel(ctx, &req); err != nil {
		log.Error("FlowTaskDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowTaskAdd 流量分析创建
func FlowTaskAdd(c *gin.Context) {
	var (
		req  typespec.FlowTaskAddReq
		resp typespec.FlowTaskAddResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskAdd parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowTaskAdd(ctx, &req); err != nil {
		log.Error("FlowTaskAdd application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ChangeFlowTaskStatus 操作流量分析
func ChangeFlowTaskStatus(c *gin.Context) {
	var (
		req  typespec.ChangeFlowTaskStatusReq
		resp typespec.ChangeFlowTaskStatusResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ChangeFlowTaskStatus parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.ChangeFlowTaskStatus(ctx, &req); err != nil {
		log.Error("ChangeFlowTaskStatus application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowTaskInfo 任务详情
func FlowTaskInfo(c *gin.Context) {
	var (
		req  typespec.FlowTaskInfoReq
		resp typespec.FlowTaskInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowTaskInfo(ctx, &req, &resp); err != nil {
		log.Error("FlowTaskInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// HttpsCert https证书下载
func HttpsCert(c *gin.Context) {
	var (
		req  typespec.HttpsCertReq
		resp typespec.HttpsCertResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("HttpsCert parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.HttpsCert(ctx, &req, &resp); err != nil {
		log.Error("HttpsCert application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowTaskStatus 任务状态
func FlowTaskStatus(c *gin.Context) {
	var (
		req  typespec.FlowTaskStatusReq
		resp typespec.FlowTaskStatusResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskStatus parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowTaskStatus(ctx, &req, &resp); err != nil {
		log.Error("FlowTaskStatus application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowRiskList 漏洞列表
func FlowRiskList(c *gin.Context) {
	var (
		req  typespec.FlowRiskListReq
		resp typespec.FlowRiskListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskVulList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowRiskList(ctx, &req, &resp); err != nil {
		log.Error("FlowTaskVulList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowRiskInfo 漏洞列表
func FlowRiskInfo(c *gin.Context) {
	var (
		req  typespec.FlowRiskInfoReq
		resp typespec.FlowRiskInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskVulInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowRiskInfo(ctx, &req, &resp); err != nil {
		log.Error("FlowTaskVulInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowRiskDel 漏洞删除
func FlowRiskDel(c *gin.Context) {
	var (
		req  typespec.FlowRiskDelReq
		resp typespec.FlowRiskDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskVulList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowRiskDel(ctx, &req); err != nil {
		log.Error("FlowTaskVulList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowBaseList 被动流量列表
func FlowBaseList(c *gin.Context) {
	var (
		req  typespec.FlowBaseListReq
		resp typespec.FlowBaseListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowBaseList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowBaseList(ctx, &req, &resp); err != nil {
		log.Error("FlowBaseList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowBaseInfo 被动流量详情
func FlowBaseInfo(c *gin.Context) {
	var (
		req  typespec.FlowBaseInfoReq
		resp typespec.FlowBaseInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowBaseInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowBaseInfo(ctx, &req, &resp); err != nil {
		log.Error("FlowBaseInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowBaseDel 被动流量删除
func FlowBaseDel(c *gin.Context) {
	var (
		req  typespec.FlowBaseDelReq
		resp typespec.FlowBaseDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowBaseDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowBaseDel(ctx, &req); err != nil {
		log.Error("FlowBaseDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowLogInfo 被动流量日志查询
func FlowLogInfo(c *gin.Context) {
	var (
		req  typespec.FlowLogInfoReq
		resp typespec.FlowLogInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowLogInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowLogInfo(ctx, &req, &resp); err != nil {
		log.Error("FlowLogInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowLogDel 被动流量日志清除
func FlowLogDel(c *gin.Context) {
	var (
		req  typespec.FlowLogDelReq
		resp typespec.FlowLogDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowLogDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowLogDel(ctx, &req); err != nil {
		log.Error("FlowLogDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowTaskEdit 流量分析保存
func FlowTaskEdit(c *gin.Context) {
	var (
		req  typespec.FlowTaskEditReq
		resp typespec.FlowTaskEditResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowTaskEdit(ctx, &req); err != nil {
		log.Error("FlowTaskAdd application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FlowTaskExport 被动流量 流量导出
func FlowTaskExport(c *gin.Context) {
	var (
		req  typespec.FlowTaskExportReq
		resp typespec.FlowTaskExportResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskExport parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.FlowApp
	if err := app.FlowTaskExport(ctx, &req, &resp); err != nil {
		log.Error("FlowTaskAdd application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
