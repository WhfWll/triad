package rest

import (
	"context"
	"strconv"
	"strings"

	"smart/api/typespec"
	"smart/application"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
)

func DataSecDBTargetList(c *gin.Context) {
	var req typespec.DatasecDBTargetListReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.NewDatasecDBTargetApp()
	resp, err := app.List(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func DataSecDBTargetSave(c *gin.Context) {
	var req typespec.DatasecDBTargetSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.NewDatasecDBTargetApp()
	if err := app.Save(ctx, appSecUID(c), &req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	if req.ID > 0 {
		addOperateAuditf(c, ctx, "编辑了数据安全目标库，目标ID: %d，名称: %s", req.ID, req.Name)
	} else {
		addOperateAuditf(c, ctx, "新增了数据安全目标库，名称: %s，目标: %s:%v/%s", req.Name, req.DBHost, req.DBPort, req.DBName)
	}
	server.RespSuccess(c, nil)
}

func DataSecDBTargetDelete(c *gin.Context) {
	var req typespec.DatasecDBTargetDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.NewDatasecDBTargetApp()
	if err := app.Delete(ctx, appSecUID(c), req.ID); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "删除了数据安全目标库，目标ID: %d", req.ID)
	server.RespSuccess(c, nil)
}

func DataSecDBTargetImport(c *gin.Context) {
	var req typespec.DatasecDBTargetImportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.NewDatasecDBTargetApp()
	n, err := app.Import(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "导入了数据安全目标库，导入数量: %d", n)
	server.RespSuccess(c, gin.H{"imported": n})
}

func DataSecDBTargetExport(c *gin.Context) {
	idsStr := strings.TrimSpace(c.Query("ids"))
	includePassword := c.Query("includePassword") == "1" || c.Query("includePassword") == "true"
	var ids []int
	if idsStr != "" {
		for _, part := range strings.Split(idsStr, ",") {
			if id, err := strconv.Atoi(strings.TrimSpace(part)); err == nil && id > 0 {
				ids = append(ids, id)
			}
		}
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.NewDatasecDBTargetApp()
	resp, err := app.Export(ctx, appSecUID(c), ids, includePassword)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "导出了数据安全目标库，指定ID数量: %d，包含密码: %t，导出数量: %d", len(ids), includePassword, len(resp.Items))
	server.RespSuccess(c, resp)
}

func DataSecTaskCloneTargets(c *gin.Context) {
	var req typespec.DatasecTaskCloneTargetsReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.CloneTaskTargets(ctx, appSecUID(c), req.ID, req.Kind)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "复制了数据安全历史任务目标，任务ID: %s，类型: %s，目标数量: %d", req.ID, req.Kind, len(resp.Targets))
	server.RespSuccess(c, resp)
}

func DataSecTaskRerun(c *gin.Context) {
	var req typespec.DatasecTaskRerunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	resp, err := app.RerunTask(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "发起了数据安全再次检测，原任务ID: %s，类型: %s，新任务ID: %s", req.ID, req.Kind, resp.ID)
	server.RespSuccess(c, resp)
}

func DataSecTaskDelete(c *gin.Context) {
	var req typespec.DatasecTaskDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.DataSecScan
	if err := app.DeleteTask(ctx, appSecUID(c), req.ID, req.Kind); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "删除了数据安全任务，任务ID: %s，类型: %s", req.ID, req.Kind)
	server.RespSuccess(c, nil)
}

func DataSecSaveTargetsToLibrary(c *gin.Context) {
	var req struct {
		DBType    int                             `json:"dbType"`
		GroupName string                          `json:"groupName"`
		Targets   []typespec.DataSecDBTargetInput `json:"targets"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.NewDatasecDBTargetApp()
	n, err := app.SaveTargetsFromTask(ctx, appSecUID(c), req.DBType, req.Targets, req.GroupName)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "将任务目标保存到了数据安全目标库，分组: %s，保存数量: %d", req.GroupName, n)
	server.RespSuccess(c, gin.H{"saved": n})
}

func DataSecDBTargetTestConn(c *gin.Context) {
	var req typespec.DatasecDBTargetTestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.NewDatasecDBTargetApp()
	resp, err := app.TestConnection(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "测试了数据安全目标库连接，目标ID: %d，目标: %s:%v/%s，结果: %t", req.ID, req.DBHost, req.DBPort, req.DBName, resp.OK)
	server.RespSuccess(c, resp)
}

func DataSecDBTargetBatchTestConn(c *gin.Context) {
	var req typespec.DatasecDBTargetBatchTestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.NewDatasecDBTargetApp()
	resp, err := app.BatchTestConnection(ctx, appSecUID(c), &req)
	if err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	addOperateAuditf(c, ctx, "批量测试了数据安全目标库连接，目标数量: %d，成功: %d，失败: %d", len(req.IDs), resp.OK, resp.Fail)
	server.RespSuccess(c, resp)
}
