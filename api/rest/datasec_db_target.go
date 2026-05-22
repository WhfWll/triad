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
	server.RespSuccess(c, resp)
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
	server.RespSuccess(c, resp)
}