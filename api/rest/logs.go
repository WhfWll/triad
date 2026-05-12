package rest

import (
	"context"
	"smart/api/typespec"
	"smart/application"
	"smart/tools/errno"
	"smart/tools/file"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

type LogAudit struct {
}

// LogsEnum 日志枚举
func LogsEnum(c *gin.Context) {
	var res typespec.LogsEnumRes
	app := application.Logs{}
	app.LogsEnum(&res)

	server.RespSuccess(c, res)
}

// LogAuditList 审计日志-列表
func LogAuditList(c *gin.Context) {
	var (
		req typespec.LogAuditListReq
		res typespec.LogAuditListRes
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogAuditList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err = app.LogAuditList(ctx, &req, &res)
	if err != nil {
		log.Error("LogAuditList request error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LogAuditEmpty 审计日志-清空
func LogAuditEmpty(c *gin.Context) {
	var res typespec.LogAuditEmptyRes

	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err := app.LogAuditEmpty(ctx)
	if err != nil {
		log.Error("LogAuditEmpty request error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// LogBackupConfig 日志备份配置的新增或修改
func LogBackupConfig(c *gin.Context) {
	var (
		req typespec.LogBackupConfigReq
		res typespec.LogBackupConfigRes
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogBackupConfig parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err = app.LogBackupConfig(ctx, &req)
	if err != nil {
		log.Error("LogBackupConfig request error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LogBackupConfigInfo 日志备份配置信息
func LogBackupConfigInfo(c *gin.Context) {
	var res typespec.LogBackupConfigInfoRes

	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err := app.LogBackupConfigInfo(ctx, &res)
	if err != nil {
		log.Error("LogBackupConfigInfo request error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// 设置日志过期时间
func SetLogExpirationTime(c *gin.Context) {
	var (
		req  typespec.SetLogExpirationTimeReq
		resp typespec.SetLogExpirationTimeResp
	)
	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("SetLogExpirationTime parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}
	if req.ExpirationTime < 0 {
		server.RespFail(c, errno.CheckParamsErr, "日志保留时间不能为负")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err = app.SetLogExpirationTime(ctx, &req)
	if err != nil {
		log.Error("SetLogExpirationTime application error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 查询日志过期时间
func GetLogExpirationTime(c *gin.Context) {
	var resp typespec.GetLogExpirationTimeResp
	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err := app.GetLogExpirationTime(ctx, &resp)
	if err != nil {
		log.Error("GetLogExpirationTime application error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// LogBackupNow 立即备份
func LogBackupNow(c *gin.Context) {
	var res typespec.LogBackupNowRes

	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err := app.LogBackupNow(ctx)
	if err != nil {
		log.Error("LogBackupNow request error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LogBackupList 日志备份列表
func LogBackupList(c *gin.Context) {
	var (
		req typespec.LogBackupListReq
		res typespec.LogBackupListRes
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogBackupList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err = app.LogBackupList(ctx, &req, &res)
	if err != nil {
		log.Error("LogBackupList request error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LogBackupDelete 日志备份删除
func LogBackupDelete(c *gin.Context) {
	var (
		req typespec.LogBackupDeleteReq
		res typespec.LogBackupDeleteRes
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogBackupDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err = app.LogBackupDelete(ctx, &req)
	if err != nil {
		log.Error("LogBackupDelete request error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// LogBackupDownload 日志备份下载
func LogBackupDownload(c *gin.Context) {
	var (
		req typespec.LogBackupDownloadReq
		res typespec.LogBackupDownloadRes
	)

	err := c.ShouldBind(&req)
	if err != nil {
		log.Error("LogBackupDownload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	app := application.Logs{}
	err = app.LogBackupDownload(ctx, &req, &res)
	if err != nil {
		log.Error("LogBackupDownload request error,the fail reason is：" + err.Error())
		server.RespFail(c, errno.CheckParamsErr, err.Error())
		return
	}
	//判断服务器端的文件是否存在（因其他原因导致的文件不存在）
	if !file.CheckPathExist(res.Path) {
		log.Errorf("log backup file (%s) is not exist", res.Path)
		server.RespFail(c, 4000, "备份文件已不存在")
		return
	}
	filepathList := strings.Split(res.Path, "/")
	filename := filepathList[len(filepathList)-1]
	c.Header("Content-Type", "application/octet-stream")             //唤起浏览器的下载
	c.Header("Content-Disposition", "attachment;filename="+filename) //设置文件名
	c.File(res.Path)
}
