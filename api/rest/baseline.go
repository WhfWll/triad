package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

func BaselineCheckRun(c *gin.Context) {
	var req typespec.BaselineCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("BaselineCheckRun param error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.RunBaselineCheck(ctx, &req)
	if err != nil {
		log.Errorf("BaselineCheckRun error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func BaselineCheckResultList(c *gin.Context) {
	var req typespec.BaselineCheckResultListReq
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("BaselineCheckResultList param error: %v", err)
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetBaselineResults(ctx, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func BaselineCheckStat(c *gin.Context) {
	var req typespec.BaselineStatReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetBaselineStat(ctx, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func BaselineRulesList(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp := app.GetBaselineRules(ctx)
	server.RespSuccess(c, resp)
}

func MalwareScanRun(c *gin.Context) {
	var req typespec.MalwareScanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("MalwareScanRun param error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.RunMalwareScan(ctx, &req)
	if err != nil {
		log.Errorf("MalwareScanRun error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func MalwareResultList(c *gin.Context) {
	var req typespec.MalwareResultListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetMalwareResults(ctx, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DBCheckRun(c *gin.Context) {
	var req typespec.DBCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("DBCheckRun param error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.RunDBCheck(ctx, &req)
	if err != nil {
		log.Errorf("DBCheckRun error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DBCheckResultList(c *gin.Context) {
	var req typespec.DBCheckResultListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetDBCheckResults(ctx, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func SensitiveDataScanRun(c *gin.Context) {
	var req typespec.SensitiveDataScanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("SensitiveDataScanRun param error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.RunSensitiveDataScan(ctx, &req)
	if err != nil {
		log.Errorf("SensitiveDataScanRun error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func SensitiveDataResultList(c *gin.Context) {
	var req typespec.SensitiveDataListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetSensitiveDataResults(ctx, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func SensitiveDataStat(c *gin.Context) {
	var req typespec.SensitiveDataStatReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetSensitiveDataStat(ctx, &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func BaselineEnums(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp := app.GetEnums(ctx)
	server.RespSuccess(c, resp)
}
