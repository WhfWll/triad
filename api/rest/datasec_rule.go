package rest

import (
	"context"
	"strconv"

	"smart/api/typespec"
	"smart/application"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

func DataSecRulesList(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	server.RespSuccess(c, app.GetRulesFromDB(ctx))
}

func DataSecRulesStats(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	server.RespSuccess(c, app.GetRulesStatsFromDB(ctx))
}

func DataSecRulesReload(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	if err := app.ReloadFromDB(ctx); err != nil {
		log.Errorf("DataSecRulesReload: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAudit(c, ctx, "重载了数据安全检测规则")
	server.RespSuccess(c, nil)
}

func DataSecRulesImport(c *gin.Context) {
	var req typespec.DatasecRulesImportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	resp := app.ImportRules(ctx, &req)
	addOperateAuditf(c, ctx, "导入了数据安全检测规则，导入条数: %d", len(req.Rules))
	server.RespSuccess(c, resp)
}

func DataSecRulesImportBuiltin(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	resp := app.ImportBuiltinRules(ctx)
	addOperateAudit(c, ctx, "导入了内置数据安全检测规则")
	server.RespSuccess(c, resp)
}

func DataSecRuleCreate(c *gin.Context) {
	var req typespec.DatasecRuleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	if err := app.CreateRule(ctx, &req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "新增了数据安全检测规则，名称: %s", req.Name)
	server.RespSuccess(c, nil)
}

func DataSecRuleUpdate(c *gin.Context) {
	var req typespec.DatasecRuleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	if err := app.UpdateRule(ctx, &req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "编辑了数据安全检测规则，规则ID: %d", req.ID)
	server.RespSuccess(c, nil)
}

func DataSecRuleDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	if err := app.DeleteRule(ctx, id); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "删除了数据安全检测规则，规则ID: %d", id)
	server.RespSuccess(c, nil)
}

func DataSecCveImportPreview(c *gin.Context) {
	var app application.DataSecRuleApp
	resp, err := app.PreviewCveImport()
	if err != nil {
		// 仍返回 data 便于前端展示条数/路径提示
		if resp != nil {
			server.RespSuccess(c, resp)
			return
		}
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DataSecRulesImportFromCve(c *gin.Context) {
	var req typespec.DatasecCveImportReq
	_ = c.ShouldBindJSON(&req)
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	resp, err := app.ImportFromCve(ctx, req.Limit)
	if err != nil {
		log.Errorf("DataSecRulesImportFromCve: %v", err)
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "从CVE库导入了数据安全检测规则，限制条数: %d", req.Limit)
	server.RespSuccess(c, resp)
}

func DataSecRuleDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecRuleApp
	resp, err := app.GetRuleDetail(ctx, id)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "查看了数据安全检测规则详情，规则ID: %d", id)
	server.RespSuccess(c, resp)
}
