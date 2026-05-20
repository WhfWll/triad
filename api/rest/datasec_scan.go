package rest

import (
	"context"

	"smart/api/typespec"
	"smart/application"

	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
)

func DataSecDBCheckRun(c *gin.Context) {
	var req typespec.DataSecDBRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.RunDBCheck(ctx, appSecUID(c), &req)
	if err != nil {
		log.Errorf("DataSecDBCheckRun: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DataSecDBCheckList(c *gin.Context) {
	var req typespec.DataSecScanListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.ListDBChecks(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DataSecDBCheckDetail(c *gin.Context) {
	var req typespec.DataSecScanDetailReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.GetDBCheckDetail(ctx, appSecUID(c), req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DataSecSensitiveScanRun(c *gin.Context) {
	var req typespec.DataSecSensitiveRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.RunSensitiveScan(ctx, appSecUID(c), &req)
	if err != nil {
		log.Errorf("DataSecSensitiveScanRun: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DataSecSensitiveScanList(c *gin.Context) {
	var req typespec.DataSecScanListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.ListSensitiveScans(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DataSecSensitiveScanDetail(c *gin.Context) {
	var req typespec.DataSecScanDetailReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.GetSensitiveScanDetail(ctx, appSecUID(c), req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
