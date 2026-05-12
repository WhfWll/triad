package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// TaskGroupCreate 任务组新建
func TaskGroupCreate(c *gin.Context) {
	var (
		req  typespec.TaskGroupCreateReq
		resp typespec.TaskGroupCreateResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskGroupCreateReq parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskGroup
	if err := app.Create(ctx, &req, &resp); err != nil {
		log.Error("TaskGroupCreateReq parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskGroupList 任务组列表
func TaskGroupList(c *gin.Context) {
	var (
		req  typespec.TaskGroupListReq
		resp typespec.TaskGroupListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskGroupCreateReq parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskGroup
	if err := app.List(ctx, &req, &resp); err != nil {
		log.Error("TaskGroupCreateReq parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskGroupDelete 任务组删除
func TaskGroupDelete(c *gin.Context) {
	var (
		req  typespec.TaskGroupDeleteReq
		resp typespec.TaskGroupDeleteResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskGroupCreateReq parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskGroup
	if err := app.Delete(ctx, &req, &resp); err != nil {
		log.Error("TaskGroupCreateReq parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskGroupGroupBind 任务组绑定
func TaskGroupGroupBind(c *gin.Context) {
	var (
		req  typespec.TaskGroupGroupBindReq
		resp typespec.TaskGroupGroupBindResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskGroupCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskGroup
	if err := app.GroupBind(ctx, &req, &resp); err != nil {
		log.Error("TaskGroupCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskGroupTaskList 任务组内任务列表
func TaskGroupTaskList(c *gin.Context) {
	var (
		req  typespec.TaskGroupTaskListReq
		resp typespec.TaskGroupTaskListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskGroupCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskGroup
	if err := app.TaskList(ctx, &req, &resp); err != nil {
		log.Error("TaskGroupCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskGroupOverView 任务组 任务统计
func TaskGroupOverView(c *gin.Context) {
	var (
		req  typespec.TaskGroupOverViewReq
		resp typespec.TaskGroupOverViewResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskGroupCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskGroup
	if err := app.Overview(ctx, &req, &resp); err != nil {
		log.Error("TaskGroupCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskGroupStatus 任务组 任务状态
func TaskGroupStatus(c *gin.Context) {
	var (
		req  typespec.TaskGroupStatusReq
		resp typespec.TaskGroupStatusResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskGroupCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskGroup
	if err := app.Status(ctx, &req, &resp); err != nil {
		log.Error("TaskGroupCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskGroupEdit 任务组编辑
func TaskGroupEdit(c *gin.Context) {
	var (
		req  typespec.TaskGroupEditReq
		resp typespec.TaskGroupEditResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskGroupCreateReq parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskGroup
	if err := app.Edit(ctx, &req, &resp); err != nil {
		log.Error("TaskGroupCreateReq parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
