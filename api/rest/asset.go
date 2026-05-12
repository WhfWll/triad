package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// GetAssetEnums 资产枚举
func GetAssetEnums(c *gin.Context) {
	var (
		resp typespec.GetAssetEnumsResp
		app  application.AssetApp
	)
	ctx := server.NewContext(context.Background(), c)
	if err := app.GetAssetEnums(ctx, &resp); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 近期变化资产列表
func GetChangeAssetList(c *gin.Context) {
	var (
		req  typespec.GetChangeAssetListReq
		resp typespec.GetChangeAssetListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetApp
	if err := app.GetChangeAssetList(ctx, &req, &resp); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 近期变化资产删除
func ChangeAssetDel(c *gin.Context) {
	var (
		req  typespec.ChangeAssetDelReq
		resp typespec.ChangeAssetDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetApp
	if err := app.ChangeAssetDel(ctx, &req); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// Statistics 资产统计
func Statistics(c *gin.Context) {
	var (
		req  typespec.AssetSummarizeReq
		resp typespec.AssetSummarizeRes
	)
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.AssetCollectApp
	if err := app.AssetSummarizeTS(ctx, &req, &resp); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
