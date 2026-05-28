package rest

import (
	"context"
	"fmt"

	"smart/api/typespec"
	"smart/application"
	"smart/services"
	"smart/tools/enums"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	username, _ := c.Get("username")
	ip := c.ClientIP()
	var svc services.LogAudit
	target := req.Target
	if target == "" {
		target = req.TargetURL
	}
	_ = svc.LogAuditAdd(ctx, enums.LogAuditTypeOperate,
		fmt.Sprintf("创建了应用安全动态扫描任务，目标: %s", target),
		fmt.Sprintf("%v", username), ip)
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
	username, _ := c.Get("username")
	ip := c.ClientIP()
	var svc services.LogAudit
	target := req.Target
	if target == "" {
		target = req.TargetURL
	}
	_ = svc.LogAuditAdd(ctx, enums.LogAuditTypeOperate,
		fmt.Sprintf("创建了应用安全专项检测任务，目标: %s", target),
		fmt.Sprintf("%v", username), ip)
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
	addOperateAuditf(c, ctx, "查看了应用安全动态扫描详情，任务ID: %s", req.ID)
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
	addOperateAuditf(c, ctx, "查看了应用安全专项检测详情，任务ID: %s", req.ID)
	server.RespSuccess(c, resp)
}
