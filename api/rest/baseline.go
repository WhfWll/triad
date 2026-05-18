package rest

import (
	"context"
	"smart/api/typespec"
	"smart/application"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
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

func BaselineBatchCheckRun(c *gin.Context) {
	var req typespec.BaselineBatchCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("BaselineBatchCheckRun param error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.RunBaselineBatchCheckAsync(ctx, &req)
	if err != nil {
		log.Errorf("BaselineBatchCheckRun error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func BaselineBatchCheckProgress(c *gin.Context) {
	taskID := c.Query("taskId")
	if taskID == "" {
		server.RespFail(c, 4000, "参数错误: taskId 必填")
		return
	}
	id, err := strconv.Atoi(taskID)
	if err != nil {
		server.RespFail(c, 4000, "参数错误: taskId 必须为数字")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp := app.GetBatchTaskProgress(ctx, id)
	if resp == nil {
		server.RespFail(c, 4000, "任务不存在")
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

func BaselineTaskList(c *gin.Context) {
	var req typespec.BaselineTaskListReq
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("BaselineTaskList param error: %v", err)
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetBaselineTaskList(ctx, &req)
	if err != nil {
		log.Errorf("BaselineTaskList error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func BaselineTaskTargets(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Query("taskId"))
	if err != nil || taskID <= 0 {
		server.RespFail(c, 4000, "taskId 参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	items, err := app.GetBaselineTaskTargets(ctx, taskID)
	if err != nil {
		log.Errorf("BaselineTaskTargets error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, map[string]interface{}{
		"list":  items,
		"total": len(items),
	})
}

func BaselineRulesList(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp := app.GetBaselineRules(ctx)
	server.RespSuccess(c, resp)
}

func BaselineRulesReload(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	if err := app.ReloadBaselineRulesFromDB(ctx); err != nil {
		log.Errorf("BaselineRulesReload: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, map[string]bool{"ok": true})
}

func BaselineRulesImport(c *gin.Context) {
	var req typespec.BaselineRulesImportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("BaselineRulesImport param error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp := app.ImportBaselineRules(ctx, &req)
	server.RespSuccess(c, resp)
}

// BaselineRulesListFromDB 从数据库获取规则列表
func BaselineRulesListFromDB(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp := app.GetBaselineRulesFromDB(ctx)
	server.RespSuccess(c, resp)
}

// BaselineRuleDetail 获取单条规则详情
func BaselineRuleDetail(c *gin.Context) {
	var req struct {
		ID int `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil || req.ID <= 0 {
		server.RespFail(c, 4000, "参数错误: id 必填")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetBaselineRuleDetail(ctx, req.ID)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// BaselineRuleCreate 新增规则
func BaselineRuleCreate(c *gin.Context) {
	var req typespec.BaselineRuleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("BaselineRuleCreate param error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	if err := app.CreateBaselineRule(ctx, &req); err != nil {
		log.Errorf("BaselineRuleCreate error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, map[string]bool{"ok": true})
}

// BaselineRuleUpdate 编辑规则
func BaselineRuleUpdate(c *gin.Context) {
	var req typespec.BaselineRuleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("BaselineRuleUpdate param error: %v", err)
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	if err := app.UpdateBaselineRule(ctx, &req); err != nil {
		log.Errorf("BaselineRuleUpdate error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, map[string]bool{"ok": true})
}

// BaselineRuleDelete 删除规则
func BaselineRuleDelete(c *gin.Context) {
	var req struct {
		ID int `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil || req.ID <= 0 {
		server.RespFail(c, 4000, "参数错误: id 必填")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	if err := app.DeleteBaselineRule(ctx, req.ID); err != nil {
		log.Errorf("BaselineRuleDelete error: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, map[string]bool{"ok": true})
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

func MalwareScanTaskList(c *gin.Context) {
	var req typespec.MalwareTaskListReq
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("MalwareScanTaskList param error: %v", err)
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.BaselineApp
	resp, err := app.GetMalwareTaskList(ctx, &req)
	if err != nil {
		log.Errorf("MalwareScanTaskList error: %v", err)
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
