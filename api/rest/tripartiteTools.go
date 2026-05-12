package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"io"
	"smart/api/typespec"
	"smart/application"
)

// 三方工具

// TripartiteToolsXrayCreate xray 创建任务
func TripartiteToolsXrayCreate(c *gin.Context) {
	var (
		req typespec.TripartiteToolsXRayCreateReq
		res typespec.TripartiteToolsXRayCreateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsXrayCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsXrayCreate(ctx, &req); err != nil {
		log.Error("TripartiteToolsXrayCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsXrayUpload xray 导入
func TripartiteToolsXrayUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Error("TripartiteToolsXrayUpload upload file error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "文件上传错误,错误原因："+err.Error())
		return
	}
	fileFd, err := file.Open()
	if err != nil {
		log.Error("TripartiteToolsXrayUpload Open file error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "文件上传错误,错误原因："+err.Error())
		return
	}
	content, err := io.ReadAll(fileFd)
	if err != nil {
		log.Error("TripartiteToolsXrayUpload upload file read error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "文件上传读取错误,错误原因："+err.Error())
		return
	}

	var (
		req typespec.TripartiteToolsXRayUploadReq
		res typespec.TripartiteToolsXRayUploadRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsXrayUpload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsXrayUpload(ctx, &req, content); err != nil {
		log.Error("TripartiteToolsXrayUpload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsXrayDel xray 删除任务
func TripartiteToolsXrayDel(c *gin.Context) {
	var (
		req typespec.TripartiteToolsXRayDelReq
		res typespec.TripartiteToolsXRayDelRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsXrayDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsXrayDel(ctx, &req); err != nil {
		log.Error("TripartiteToolsXrayDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsXRayPage xray 任务列表
func TripartiteToolsXRayPage(c *gin.Context) {
	var (
		req typespec.TripartiteToolsXRayPageReq
		res typespec.TripartiteToolsXRayPageRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsXRayPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsXRayPage(ctx, &req, &res); err != nil {
		log.Error("TripartiteToolsXRayPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)

}

// TripartiteToolsXRayDetailPage xray 任务详情列表
func TripartiteToolsXRayDetailPage(c *gin.Context) {

	var (
		req typespec.TripartiteToolsXRayDetailPageReq
		res typespec.TripartiteToolsXRayDetailPageRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsXRayDetailPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsXRayDetailPage(ctx, &req, &res); err != nil {
		log.Error("TripartiteToolsXRayDetailPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsBurpsuiteCreate burpsuite 创建任务
func TripartiteToolsBurpsuiteCreate(c *gin.Context) {
	var (
		req typespec.TripartiteToolsBurpsuiteCreateReq
		res typespec.TripartiteToolsBurpsuiteCreateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsBurpsuiteCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsBurpsuiteCreate(ctx, &req); err != nil {
		log.Error("TripartiteToolsBurpsuiteCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsBurpsuiteUpload burpsuite 导入
func TripartiteToolsBurpsuiteUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Error("TripartiteToolsBurpsuiteUpload upload file error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "文件上传错误,错误原因："+err.Error())
		return
	}
	fileFd, err := file.Open()
	if err != nil {
		log.Error("TripartiteToolsBurpsuiteUpload Open file error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "文件上传错误,错误原因："+err.Error())
		return
	}
	content, err := io.ReadAll(fileFd)
	if err != nil {
		log.Error("TripartiteToolsBurpsuiteUpload upload file read error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "文件上传读取错误,错误原因："+err.Error())
		return
	}

	var (
		req typespec.TripartiteToolsBurpsuiteUploadReq
		res typespec.TripartiteToolsBurpsuiteUploadRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsBurpsuiteUpload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsBurpsuiteUpload(ctx, &req, content); err != nil {
		log.Error("TripartiteToolsBurpsuiteUpload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsBurpsuiteDel burpsuite 删除任务
func TripartiteToolsBurpsuiteDel(c *gin.Context) {
	var (
		req typespec.TripartiteToolsBurpsuiteDelReq
		res typespec.TripartiteToolsBurpsuiteDelRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsBurpsuiteDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsBurpsuiteDel(ctx, &req); err != nil {
		log.Error("TripartiteToolsBurpsuiteDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsBurpsuitePage burpsuite 任务列表
func TripartiteToolsBurpsuitePage(c *gin.Context) {
	var (
		req typespec.TripartiteToolsBurpsuitePageReq
		res typespec.TripartiteToolsBurpsuitePageRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsBurpsuitePage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsBurpsuitePage(ctx, &req, &res); err != nil {
		log.Error("TripartiteToolsBurpsuitePage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsBurpsuiteDetailPage burpsuite 任务详情列表
func TripartiteToolsBurpsuiteDetailPage(c *gin.Context) {
	var (
		req typespec.TripartiteToolsBurpsuiteDetailPageReq
		res typespec.TripartiteToolsBurpsuiteDetailPageRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsBurpsuiteDetailPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsBurpsuiteDetailPage(ctx, &req, &res); err != nil {
		log.Error("TripartiteToolsBurpsuiteDetailPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsWifiApList wifi 所有在线wifi列表
func TripartiteToolsWifiApList(c *gin.Context) {
	var (
		req typespec.TripartiteToolsWifiApListReq
		res typespec.TripartiteToolsWifiApListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsWifiApList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsWifiApList(ctx, &req, &res); err != nil {
		log.Error("TripartiteToolsWifiApList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsWifiCreate wifi 创建
func TripartiteToolsWifiCreate(c *gin.Context) {
	var (
		req typespec.TripartiteToolsWifiCreateReq
		res typespec.TripartiteToolsWifiCreateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsWifiApCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsWifiApCreate(ctx, &req, &res); err != nil {
		log.Error("TripartiteToolsWifiApCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsWifiPage wifi 列表
func TripartiteToolsWifiPage(c *gin.Context) {
	var (
		req typespec.TripartiteToolsWifiPageReq
		res typespec.TripartiteToolsWifiPageRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsWifiPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsWifiPage(ctx, &req, &res); err != nil {
		log.Error("TripartiteToolsWifiPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TripartiteToolsWifiDel wifi 任务删除
func TripartiteToolsWifiDel(c *gin.Context) {
	var (
		req typespec.TripartiteToolsWifiDelReq
		res typespec.TripartiteToolsWifiDelRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TripartiteToolsWifiDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TripartiteTools
	if err := app.TripartiteToolsWifiDel(ctx, &req, &res); err != nil {
		log.Error("TripartiteToolsWifiDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}
