package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"io/ioutil"
	"smart/api/typespec"
	"smart/application"
	"smart/tools/enums"
	"smart/tools/errno"
)

// RemoteSessionList 远程会话列表及筛选
func RemoteSessionList(c *gin.Context) {
	var (
		req  typespec.RemoteSessionListReq
		resp typespec.RemoteSessionListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("RemoteSessionList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.RemoteSessionApp
	if err := app.RemoteSessionList(ctx, &req, &resp); err != nil {
		log.Error("RemoteSessionList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// RemoteSessionDel 远程会话列表删除
func RemoteSessionDel(c *gin.Context) {
	var req typespec.RemoteSessionDelReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("RemoteSessionDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.RemoteSessionApp
	if err := app.DelRemoteSession(ctx, &req); err != nil {
		log.Error("RemoteSessionDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "del success")
}

// RemoteSessionInfo 远程会话详情信息
func RemoteSessionInfo(c *gin.Context) {
	var (
		req  typespec.RemoteSessionInfoReq
		resp typespec.RemoteSessionInfoRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("RemoteSessionInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.RemoteSessionApp
	if err := app.RemoteSessionInfo(ctx, &req, &resp); err != nil {
		log.Error("RemoteSessionInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// RemoteSessionDir 远程会话列出目录
func RemoteSessionDir(c *gin.Context) {
	var (
		req  typespec.RemoteSessionDirReq
		resp typespec.RemoteSessionDirRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("RemoteSessionDir parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//ctx := server.NewContext(context.Background(), c)
	//var app application.RemoteSessionApp
	//if err := app.RemoteSessionDir(ctx, &req, &resp); err != nil {
	//	log.Error("RemoteSessionDir parameter error,the fail reason is：" + err.Error())
	//	server.RespFail(c, 4000, err.Error())
	//	return
	//}
	server.RespSuccess(c, resp)
}

// CaptureInfoEnum 漏洞取证 抓取信息枚举
func CaptureInfoEnum(c *gin.Context) {
	var resp map[int]string
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.CaptureInfoEnum(ctx, &resp); err != nil {
		log.Error("CaptureInfoEnum parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ToCaptureInfo 抓取信息
func ToCaptureInfo(c *gin.Context) {
	var (
		req  typespec.ToCaptureInfoReq
		resp typespec.ToCaptureInfoRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ToCaptureInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.ToCaptureInfo(ctx, &req, &resp); err != nil {
		log.Error("ToCaptureInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// BreakShell 退出
func BreakShell(c *gin.Context) {
	var resp typespec.ToCaptureInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.ToCaptureInfo(ctx, &typespec.ToCaptureInfoReq{ID: 7}, &resp); err != nil {
		log.Error("ToCaptureInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// FileDownload 文件下载
func FileDownload(c *gin.Context) {
	var (
		req  typespec.FileDownloadReq
		resp typespec.FileDownloadResp
	)
	if err := c.Bind(&req); err != nil {
		log.Info("FileDownload ReqParams Verification ERR：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var fileApp application.RemoteSessionApp
	fileName, filepath, err := fileApp.FileDownload(ctx, &req, &resp)
	if err != nil {
		log.Info("FileDownload ERR：" + err.Error())
		server.RespFail(c, errno.RequestErr, err.Error())
		return
	}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		server.RespFail(c, 4000, "read file fail!")
	}
	server.RespDownload(c, fileName, data)
}

// DelFile 删除文件下载
func DelFile(c *gin.Context) {
	var (
		req  typespec.DelFileReq
		resp typespec.FileDownloadResp
	)
	if err := c.Bind(&req); err != nil {
		log.Info("DelFile ReqParams Verification ERR：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var fileApp application.RemoteSessionApp
	if err := fileApp.DelFile(ctx, &req, &resp); err != nil {
		log.Info("DelFile ERR：" + err.Error())
		server.RespFail(c, errno.RequestErr, err.Error())
		return
	}
	server.RespSuccess(c, "")
}

// ExceShellMany 批量收集信息
func ExceShellMany(c *gin.Context) {
	var (
		req  typespec.ExceShellManyReq
		resp typespec.ExceShellManyResp
	)
	if err := c.Bind(&req); err != nil {
		log.Info("ExceShellMany ReqParams Verification ERR：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}
	if len(req.CaptureInfoIds) == 0 && len(req.FileName) == 0 {
		server.RespFail(c, errno.CheckParamsErr, "抓取信息和文件收集参数不能同时为空...")
		return
	}
	if len(req.FilePath) == 0 && req.CaptureType == enums.CaptureTypeCentos {
		req.FilePath = "/"
	} else if len(req.FilePath) == 0 && req.CaptureType == enums.CaptureTypeWindows {
		req.FilePath = "c:/"
	}
	ctx := server.NewContext(context.Background(), c)
	var remoteApp application.RemoteSessionApp
	if err := remoteApp.ExceShellMany(ctx, &req, &resp); err != nil {
		log.Info("ExceShellMany ExceShellMany err：" + err.Error())
		server.RespFail(c, errno.RequestErr, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func FileManagement(c *gin.Context) {
	var (
		req  typespec.FileManagementReq
		resp typespec.FileManagementResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FileManagement parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.FileManagement(ctx, &req, &resp); err != nil {
		log.Error("FileManagement parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
func ShellFileDownload(c *gin.Context) {
	var (
		req  typespec.ShellFileDownloadReq
		resp typespec.ShellFileDownloadResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ShellFileDownload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.ShellFileDownload(ctx, &req, &resp); err != nil {
		log.Error("FileManagement parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
