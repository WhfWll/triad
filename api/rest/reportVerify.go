package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"io"
	"os"
	"smart/api/typespec"
	"smart/application"
)

// ReportVerifyUpload 上传报告
func ReportVerifyUpload(c *gin.Context) {
	var (
		req  typespec.ReportVerifyUploadReq
		resp typespec.ReportVerifyUploadResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	uid, exist := c.Get("uid")
	if !exist {
		server.RespFail(c, 4000, "参数错误,错误原因："+"用户未登录")
		return
	}
	req.UserId = uid.(int)
	ctx := server.NewContext(context.Background(), c)

	formFile, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Info("ReportVerifyUpload form file err：" + err.Error())
		server.RespFail(c, 4000, "上传文件出错,错误原因："+err.Error())
		return
	}
	defer formFile.Close()

	// 创建一个临时文件以保存上传的文件
	tempFile, err := os.CreateTemp("", header.Filename)
	if err != nil {
		server.RespFail(c, 4000, "文件类型错误")
		return
	}
	defer func() {
		tempFile.Close()
	}()

	// 将上传的文件保存到临时文件中
	_, err = io.Copy(tempFile, formFile)
	if err != nil {
		return
	}

	var app application.ReportVerify
	if err := app.Upload(ctx, &req, header.Filename, tempFile.Name(), &resp); err != nil {
		log.Error("ReportVerifyUpload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyTaskList 任务列表
func ReportVerifyTaskList(c *gin.Context) {
	var (
		req  typespec.ReportVerifyTaskListReq
		resp typespec.ReportVerifyTaskListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.TaskList(c, &req, &resp); err != nil {
		log.Error("ReportVerifyTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyTaskDetail 任务详情
func ReportVerifyTaskDetail(c *gin.Context) {
	var (
		req  typespec.ReportVerifyTaskDetailReq
		resp typespec.ReportVerifyTaskDetailResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.TaskDetail(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyTargetList 目标列表
func ReportVerifyTargetList(c *gin.Context) {
	var (
		req  typespec.ReportVerifyTargetListReq
		resp typespec.ReportVerifyTargetListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.TargetList(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyTargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyPortList 端口列表
func ReportVerifyPortList(c *gin.Context) {
	var (
		req  typespec.ReportVerifyPortListReq
		resp typespec.ReportVerifyPortListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyPortList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.PortList(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyPortList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyVulList 漏洞列表
func ReportVerifyVulList(c *gin.Context) {
	var (
		req  typespec.ReportVerifyVulListReq
		resp typespec.ReportVerifyVulListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyPortList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.VulList(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyPortList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyEnum 枚举接口
func ReportVerifyEnum(c *gin.Context) {
	var (
		req  typespec.ReportVerifyEnumReq
		resp typespec.ReportVerifyEnumResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyEnum parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.Enum(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyEnum parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyStatsInfo 展示报告验证统计信息
func ReportVerifyStatsInfo(c *gin.Context) {
	var (
		req  typespec.ReportVerifyStatsInfoReq
		resp typespec.ReportVerifyStatsInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyStatsInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.StatsInfo(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyStatsInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyTaskStop 任务结束
func ReportVerifyTaskStop(c *gin.Context) {
	var (
		req  typespec.ReportVerifyTaskStopReq
		resp typespec.ReportVerifyTaskStopResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyStatsInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.TaskStop(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyStatsInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyTaskDelete 删除任务
func ReportVerifyTaskDelete(c *gin.Context) {
	var (
		req  typespec.ReportVerifyTaskDeleteReq
		resp typespec.ReportVerifyTaskDeleteResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyStatsInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.TaskDelete(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyStatsInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyTargetDelete 删除目标
func ReportVerifyTargetDelete(c *gin.Context) {
	var (
		req  typespec.ReportVerifyTargetDeleteReq
		resp typespec.ReportVerifyTargetDeleteResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyTargetDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.TargetDelete(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyTargetDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyVulDelete 删除漏洞
func ReportVerifyVulDelete(c *gin.Context) {
	var (
		req  typespec.ReportVerifyVulDeleteReq
		resp typespec.ReportVerifyVulDeleteResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyTargetDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.VulDelete(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyTargetDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReportVerifyVulDetail 漏洞详情
func ReportVerifyVulDetail(c *gin.Context) {
	var (
		req  typespec.ReportVerifyVulDetailReq
		resp typespec.ReportVerifyVulDetailResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReportVerifyTargetDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.ReportVerify
	if err := app.VulDetail(ctx, &req, &resp); err != nil {
		log.Error("ReportVerifyTargetDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
