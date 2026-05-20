package rest

import (
	"context"

	"smart/api/typespec"
	"smart/application"

	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
)

func appSecUID(c *gin.Context) int {
	uid, ok := c.Get("uid")
	if !ok {
		return 0
	}
	switch v := uid.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func AppSecDynamicScanRun(c *gin.Context) {
	var req typespec.AppSecScanRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AppSecScan
	resp, err := app.RunDynamicScan(ctx, appSecUID(c), &req)
	if err != nil {
		log.Errorf("AppSecDynamicScanRun: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func AppSecDynamicScanList(c *gin.Context) {
	var req typespec.AppSecScanListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AppSecScan
	resp, err := app.ListDynamicScans(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func AppSecAppSpecificScanRun(c *gin.Context) {
	var req typespec.AppSecScanRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AppSecScan
	resp, err := app.RunAppSpecificScan(ctx, appSecUID(c), &req)
	if err != nil {
		log.Errorf("AppSecAppSpecificScanRun: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func AppSecAppSpecificScanList(c *gin.Context) {
	var req typespec.AppSecScanListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AppSecScan
	resp, err := app.ListAppSpecificScans(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func AppSecDynamicScanDetail(c *gin.Context) {
	var req typespec.AppSecScanDetailReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AppSecScan
	resp, err := app.GetDynamicScanDetail(ctx, appSecUID(c), req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func AppSecAppSpecificScanDetail(c *gin.Context) {
	var req typespec.AppSecScanDetailReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AppSecScan
	resp, err := app.GetAppSpecificScanDetail(ctx, appSecUID(c), req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
