package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// VulScanTaskList 漏洞扫描任务列表
func VulScanTaskList(c *gin.Context) {
	var (
		req typespec.VulScanTaskListReq
		res typespec.VulScanTaskListResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.TaskList(ctx, &req, &res); err != nil {
		log.Error("VulScanTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanTaskSave 漏洞扫描任务列表
func VulScanTaskSave(c *gin.Context) {
	var (
		req typespec.VulScanTaskSaveReq
		res typespec.VulScanTaskSaveResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulScanTaskSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.TaskSave(ctx, &req, &res); err != nil {
		log.Error("VulScanTaskSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanTaskStop 漏洞扫描任务结束
func VulScanTaskStop(c *gin.Context) {
	var (
		req typespec.VulScanTaskStopReq
		res typespec.VulScanTaskStopResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTaskStop parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.TaskStop(ctx, &req, &res); err != nil {
		log.Error("VulScanTaskStop parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanTaskDelete 漏洞扫描任务删除
func VulScanTaskDelete(c *gin.Context) {
	var (
		req typespec.VulScanTaskDeleteReq
		res typespec.VulScanTaskDeleteResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTaskDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.TaskDelete(ctx, &req, &res); err != nil {
		log.Error("VulScanTaskDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanTargetList 漏洞扫描目标列表
func VulScanTargetList(c *gin.Context) {
	var (
		req typespec.VulScanTargetListReq
		res typespec.VulScanTargetListResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.TargetList(ctx, &req, &res); err != nil {
		log.Error("VulScanTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanVulList 漏洞扫描漏洞列表
func VulScanVulList(c *gin.Context) {
	var (
		req typespec.VulScanVulListReq
		res typespec.VulScanVulListResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.VulList(ctx, &req, &res); err != nil {
		log.Error("VulScanTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanVulDetail 漏洞扫描漏洞详情
func VulScanVulDetail(c *gin.Context) {
	var (
		req typespec.VulScanVulDetailReq
		res typespec.VulScanVulDetailResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.VulDetail(ctx, &req, &res); err != nil {
		log.Error("VulScanVulDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanCveList 漏洞扫描cve列表
func VulScanCveList(c *gin.Context) {
	var (
		req typespec.VulScanCveListReq
		res typespec.VulScanCveListResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanCveList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.CveList(ctx, &req, &res); err != nil {
		log.Error("VulScanCveList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanCveDetail 漏洞扫描cve详情
func VulScanCveDetail(c *gin.Context) {
	var (
		req typespec.VulScanCveDetailReq
		res typespec.VulScanCveDetailResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.CveDetail(ctx, &req, &res); err != nil {
		log.Error("VulScanVulDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanTaskOverview 任务详情
func VulScanTaskOverview(c *gin.Context) {
	var (
		req typespec.VulScanTaskOverviewReq
		res typespec.VulScanTaskOverviewResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTaskOverview parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.TaskOverview(ctx, &req, &res); err != nil {
		log.Error("VulScanTaskOverview parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// VulScanTaskState 获取任务状态
func VulScanTaskState(c *gin.Context) {
	var (
		req typespec.VulScanTaskGetStateReq
		res typespec.VulScanTaskGetStateResp
	)
	if err := c.Bind(&req); err != nil {
		log.Error("VulScanTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulScan
	if err := app.TaskState(ctx, &req, &res); err != nil {
		log.Error("VulScanVulDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}
