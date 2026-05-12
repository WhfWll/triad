// Package rest
// @Author bcy2007  2025/12/23 13:14
package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

func ReverseServerStart(c *gin.Context) {
	var req typespec.ReverseServerStartReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	var (
		app application.ReverseServer
		ctx = server.NewContext(context.Background(), c)
	)
	app.Start(ctx, &req)
	server.RespSuccess(c, nil)
}

func ReverseServerStatus(c *gin.Context) {
	var (
		resp typespec.ReverseServerStatusResp
		app  application.ReverseServer
		ctx  = server.NewContext(context.Background(), c)
	)
	app.Status(ctx, &resp)
	server.RespSuccess(c, resp)
}

func ReverseServerStop(c *gin.Context) {
	var (
		app application.ReverseServer
		ctx = server.NewContext(context.Background(), c)
	)
	app.Stop(ctx)
	server.RespSuccess(c, nil)
}

func ReverseServerMessages(c *gin.Context) {
	var (
		req  typespec.ReverseServerMessageReq
		resp typespec.ReverseServerMessageResp
	)
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	var (
		app application.ReverseServer
		ctx = server.NewContext(context.Background(), c)
	)
	app.Messages(ctx, &req, &resp)
	server.RespSuccess(c, resp)
}

func ReverseServerClearMessage(c *gin.Context) {
	var (
		app application.ReverseServer
		ctx = server.NewContext(context.Background(), c)
	)
	app.ClearMessage(ctx)
	server.RespSuccess(c, nil)
}
