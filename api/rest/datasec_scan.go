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

func DataSecDBTestConn(c *gin.Context) {
	var req typespec.DataSecDBTestConnReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.TestDBConnection(ctx, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

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
	username, _ := c.Get("username")
	ip := c.ClientIP()
	var svc services.LogAudit
	targetInfo := fmt.Sprintf("%s:%v", req.DBHost, req.DBPort)
	if len(req.Targets) > 0 {
		targetInfo = fmt.Sprintf("%s:%v", req.Targets[0].DBHost, req.Targets[0].DBPort)
	}
	_ = svc.LogAuditAdd(ctx, enums.LogAuditTypeOperate,
		fmt.Sprintf("创建了数据库安全检查任务 %s，目标: %s", req.Name, targetInfo),
		fmt.Sprintf("%v", username), ip)
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
	addOperateAuditf(c, ctx, "查看了数据库安全检查详情，任务ID: %s，目标ID: %d", req.ID, req.TargetID)
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
	username, _ := c.Get("username")
	ip := c.ClientIP()
	var svc services.LogAudit
	targetInfo := fmt.Sprintf("%s:%v", req.DBHost, req.DBPort)
	if len(req.Targets) > 0 {
		targetInfo = fmt.Sprintf("%s:%v", req.Targets[0].DBHost, req.Targets[0].DBPort)
	}
	_ = svc.LogAuditAdd(ctx, enums.LogAuditTypeOperate,
		fmt.Sprintf("创建了敏感数据扫描任务 %s，目标: %s", req.Name, targetInfo),
		fmt.Sprintf("%v", username), ip)
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
	addOperateAuditf(c, ctx, "查看了敏感数据扫描详情，任务ID: %s，目标ID: %d", req.ID, req.TargetID)
	server.RespSuccess(c, resp)
}
