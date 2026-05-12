package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// AssetTree 资产树
func AssetTree(c *gin.Context) {
	var (
		req  typespec.AssetTreeOverallReq
		resp typespec.AssetTreeOverallResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetTreeOverall parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetTree(ctx, &req, &resp); err != nil {
		log.Error("AssetTreeOverall application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AssetGroupList 资产组列表
func AssetGroupList(c *gin.Context) {
	var resp typespec.AssetGroupResp
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetGroup(ctx, &resp); err != nil {
		log.Error("AssetGroupList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AssetGroupAdd 添加资产组
func AssetGroupAdd(c *gin.Context) {
	var req typespec.AssetTreeAddReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetGroupAdd parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetGroupAdd(ctx, &req); err != nil {
		log.Error("AssetGroupAdd application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "success")
}

// AssetGroupEdit 修改资产组
func AssetGroupEdit(c *gin.Context) {
	var req typespec.AssetGroupEditReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetGroupEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetGroupEdit(ctx, &req); err != nil {
		log.Error("AssetGroupEdit application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "success")
}

// AssetGroupDel 删除资产组
func AssetGroupDel(c *gin.Context) {
	var req typespec.AssetDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetGroupDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetGroupDel(ctx, &req); err != nil {
		log.Error("AssetGroupDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "success")
}

// AssetList 资产列表
func AssetList(c *gin.Context) {
	var (
		req  typespec.AssetListReq
		resp typespec.AssetListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetList(ctx, &req, &resp); err != nil {
		log.Error("AssetList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// SelectAllAsset 资产列表
func SelectAllAsset(c *gin.Context) {
	var (
		req  typespec.SelectAllAssetReq
		resp typespec.SelectAllAssetRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SelectAllAsset parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.SelectAllAssetList(ctx, &req, &resp); err != nil {
		log.Error("SelectAllAsset application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AssetAdd 资产-资产添加
func AssetAdd(c *gin.Context) {
	var req typespec.AssetAddReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetVulTestEditStatus parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AddAsset(ctx, &req); err != nil {
		log.Error("AssetAdd add error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "success")
}

// AssetEdit 资产-资产修改
func AssetEdit(c *gin.Context) {
	var req typespec.AssetEditReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.EditAsset(ctx, &req); err != nil {
		log.Error("AssetEdit add error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "success")
}

// AssetDel 删除资产
func AssetDel(c *gin.Context) {
	var req typespec.AssetDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetDel(ctx, &req); err != nil {
		log.Error("AssetDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "success")
}

// AssetDetail 资产详情
func AssetDetail(c *gin.Context) {
	var (
		req  typespec.AssetDetailReq
		resp typespec.AssetDetail
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetDetail(ctx, &req, &resp); err != nil {
		log.Error("AssetDetail application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// Import 资产导入
func Import(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	res, err := app.ReadXlsx(ctx, c)
	if err != nil {
		log.Error("Import application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	if res != "" {
		server.RespFail(c, 201, res)
		return
	}
	server.RespSuccess(c, "success")
}

// Export 资产导出
func Export(c *gin.Context) {
	var (
		req  typespec.AssetExportListReq
		resp typespec.AssetExportListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetExportList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetExportList(ctx, &req, &resp); err != nil {
		log.Error("AssetExportList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// GetAssetGroupEnums 资产组枚举
func GetAssetGroupEnums(c *gin.Context) {
	var (
		resp typespec.GetAssetGroupEnumsResp
		app  application.AssetGroupApp
	)
	ctx := server.NewContext(context.Background(), c)
	if err := app.GetAssetGroupEnums(ctx, &resp); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AssetGroupInfo 资产组详情
func AssetGroupInfo(c *gin.Context) {
	var (
		req  typespec.AssetGroupDetailReq
		resp typespec.AssetGroupDetail
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetGroupInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	err := app.AssetGroupDetail(ctx, &req, &resp)
	if err != nil {
		log.Error("AssetGroupInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AssetConnList 资产连接列表
func AssetConnList(c *gin.Context) {
	var (
		req  typespec.AssetConnListReq
		resp typespec.AssetConnListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AssetConnList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetGroupApp
	if err := app.AssetConnList(ctx, &req, &resp); err != nil {
		log.Error("AssetConnList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
