package rest

import (
	"context"
	"smart/api/typespec"
	"smart/application"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

func VulnScanCveRun(c *gin.Context) {
	var req typespec.VulnScanCveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("VulnScanCveRun parameter error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.VulnScanCveApp
	resp, err := app.RunCveScan(ctx, &req)
	if err != nil {
		log.Errorf("VulnScanCveRun failed: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func VulnScanCveBatchRun(c *gin.Context) {
	var req typespec.VulnScanCveBatchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("VulnScanCveBatchRun parameter error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.VulnScanCveApp
	resp, err := app.RunBatchCveScan(ctx, &req)
	if err != nil {
		log.Errorf("VulnScanCveBatchRun failed: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func VulnScanCveProgress(c *gin.Context) {
	var req typespec.VulnScanCveProgressReq
	if err := c.Bind(&req); err != nil {
		log.Errorf("VulnScanCveProgress parameter error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.VulnScanCveApp
	resp, err := app.GetBatchProgress(ctx, &req)
	if err != nil {
		log.Errorf("VulnScanCveProgress failed: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}