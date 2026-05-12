package rest

import (
	"context"
	"os"
	"path/filepath"
	"smart/api/typespec"
	"smart/application"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/file"
	"smart/tools/network"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

// AuthInfo 授权信息
func AuthInfo(c *gin.Context) {
	var (
		req typespec.AuthInfoReq
		res typespec.AuthInfoRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AuthInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var systemManage application.SystemManage
	if err := systemManage.AuthInfo(ctx, &req, &res); err != nil {
		log.Error("AuthInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// AuthSave 授权
func AuthSave(c *gin.Context) {
	var (
		req typespec.AuthSaveReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AuthSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var systemManage application.SystemManage
	if err := systemManage.AuthSave(ctx, &req); err != nil {
		log.Error("AuthSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "auth success")
}

// GenerateProductID 生成产品id
func GenerateProductID(c *gin.Context) {
	ctx := server.NewContext(context.Background(), c)
	var systemManage application.SystemManage
	if err := systemManage.GenerateProductID(ctx); err != nil {
		log.Error("GenerateProductID parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "GenerateProductID success")
}

// 测试目标黑白名单更新
func TargetIpSave(c *gin.Context) {
	var (
		req  typespec.TargetIpSaveReq
		resp typespec.TargetIpSaveResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TargetIpSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//参数逻辑校验
	if req.IsOpen == enums.TargetIpWhiteBlackIsOpenOn {
		if (req.Type == enums.TargetIpWhiteBlackTypeWhite && len(req.WhiteList) == 0) || (req.Type == enums.TargetIpWhiteBlackTypeBlack && len(req.BlackList) == 0) {
			server.RespFail(c, 4000, "开启状态,名单不能为空")
			return
		}
	}
	ctx := server.NewContext(context.Background(), c)
	var systemManage application.SystemManage
	if err := systemManage.TargetIpSave(ctx, &req); err != nil {
		log.Error("TargetIpSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// SystemManualRollback 系统管理 - 升级还原 - 手动回滚
func SystemManualRollback(c *gin.Context) {
	var req typespec.SystemManualRollbackReq
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBind(&req); err != nil {
			server.RespFail(c, 4000, "参数错误: "+err.Error())
			return
		}
	}

	var app application.SystemManage
	ctx := server.NewContext(context.Background(), c)
	if err := app.SystemManualRollback(ctx, &req); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "回滚任务已开始，请查询进度")
}

// TargetIpList 测试目标黑白名单查询
func TargetIpList(c *gin.Context) {
	var resp typespec.TargetIpListResp
	ctx := server.NewContext(context.Background(), c)
	var systemManage application.SystemManage
	if err := systemManage.TargetIpList(ctx, &resp); err != nil {
		log.Error("TargetIpList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// GetReverseIpHost 远程监听查询
func GetReverseIpHost(c *gin.Context) {
	var resp typespec.GetReverseIpHostResp
	ctx := server.NewContext(context.Background(), c)
	var systemManage application.SystemManage
	if err := systemManage.GetReverseIpHost(ctx, &resp); err != nil {
		log.Error("GetReverseIpHost parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// ReverseIpHostSave 远程监听修改
func ReverseIpHostSave(c *gin.Context) {
	var (
		req  typespec.ReverseIpHostSaveReq
		resp typespec.ReverseIpHostSaveResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("ReverseIpHostSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//参数逻辑校验
	if req.ReverseType == enums.TypeCustom && len(req.ReverseHost) == 0 { //自定义
		server.RespFail(c, 4000, "自定义类型,监听IP不能为空")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var systemManage application.SystemManage
	if err := systemManage.ReverseIpHostSave(ctx, &req); err != nil {
		log.Error("TargetIpSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// SystemConfigBackupConfigSave 系统管理-配置备份-保存配置
func SystemConfigBackupConfigSave(c *gin.Context) {
	var (
		req typespec.SystemConfigBackupConfigReq
		res typespec.SystemConfigBackupConfigRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemConfigBackupConfig parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemConfigBackupConfigSave(ctx, &req); err != nil {
		log.Error("SystemConfigBackupConfig parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemConfigBackupConfigInfo 系统管理-配置备份-配置信息
func SystemConfigBackupConfigInfo(c *gin.Context) {
	var res typespec.SystemConfigBackupConfigInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemConfigBackupConfigInfo(ctx, &res); err != nil {
		log.Error("SystemConfigBackupConfigInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemConfigBackupList 系统管理-配置备份-列表
func SystemConfigBackupList(c *gin.Context) {
	var (
		req typespec.SystemConfigBackupListReq
		res typespec.SystemConfigBackupListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemConfigBackupList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemConfigBackupList(ctx, &req, &res); err != nil {
		log.Error("SystemConfigBackupList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemConfigBackupDownload 系统管理-配置备份-下载
func SystemConfigBackupDownload(c *gin.Context) {
	var (
		req typespec.SystemConfigBackupDownloadReq
		res typespec.SystemConfigBackupDownloadRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemConfigBackupDownload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemConfigBackupDownload(ctx, &req, &res); err != nil {
		log.Error("SystemConfigBackupDownload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	//判断服务器端的文件是否存在（因其他原因导致的文件不存在）
	if !file.CheckPathExist(res.Path) {
		log.Errorf("system config backup file (%s) is not exist", res.Path)
		server.RespFail(c, 4000, "备份文件已不存在")
		return
	}
	filepathList := strings.Split(res.Path, "/")
	filename := filepathList[len(filepathList)-1]
	c.Header("Content-Type", "application/octet-stream")             //唤起浏览器的下载
	c.Header("Content-Disposition", "attachment;filename="+filename) //设置文件名
	c.File(res.Path)
}

// SystemConfigBackupDelete 系统管理-配置备份-删除
func SystemConfigBackupDelete(c *gin.Context) {
	var (
		req typespec.SystemConfigBackupDeleteReq
		res typespec.SystemConfigBackupDeleteRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemConfigBackupDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemConfigBackupDelete(ctx, &req); err != nil {
		log.Error("SystemConfigBackupDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemConfigBackupNow 系统管理-配置备份-立即备份
func SystemConfigBackupNow(c *gin.Context) {
	var res typespec.SystemConfigBackupNowRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemConfigBackupNow(ctx); err != nil {
		log.Error("SystemConfigBackupNow parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemConfigBackupRestore 系统管理-配置备份-恢复
func SystemConfigBackupRestore(c *gin.Context) {
	var (
		req typespec.SystemConfigBackupRestoreReq
		res typespec.SystemConfigBackupRestoreRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemConfigBackupRestore parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemConfigBackupRestore(ctx, &req); err != nil {
		log.Error("SystemConfigBackupRestore parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingIpWhiteSave 系统管理 - 系统设置 - 系统访问白名单 - 保存
func SystemSettingIpWhiteSave(c *gin.Context) {
	var (
		req typespec.SystemSettingIpWhiteSaveReq
		res typespec.SystemSettingIpWhiteSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemSettingIpWhiteSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	//格式化ip内容为英文逗号分割
	if strings.Contains(req.Ip, "\n") {
		req.Ip = strings.Join(strings.Split(req.Ip, "\n"), ",")
	}
	//验证ip、ip段
	for _, item := range strings.Split(req.Ip, ",") {
		if !network.IpSegmentTools.VerifyNetmaskIpSegment(item) && !network.IpSegmentTools.VerifyCrossbarIpSegment(item) && !network.IpSegmentTools.VerifyIp(item) {
			log.Errorf("SystemSettingIpWhiteSave parameter error,the fail reason is：%s格式错误", item)
			server.RespFail(c, 4000, item+"格式错误")
			return
		}
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingIpWhiteSave(ctx, &req); err != nil {
		log.Error("SystemSettingIpWhiteSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingIpWhiteInfo 系统管理 - 系统设置 - 系统访问白名单 - 信息
func SystemSettingIpWhiteInfo(c *gin.Context) {
	var res typespec.SystemSettingIpWhiteInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingIpWhiteInfo(ctx, &res); err != nil {
		log.Error("SystemSettingIpWhiteInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingSyslogSave 系统管理 - 系统设置 - syslog服务 - 保存
func SystemSettingSyslogSave(c *gin.Context) {
	var (
		req typespec.SystemSettingSyslogSaveReq
		res typespec.SystemSettingSyslogSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemSettingSyslogSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//校验参数的格式
	if !network.IpSegmentTools.VerifyIp(req.Ip) {
		server.RespFail(c, 4000, "ip格式错误")
		return
	}
	if !network.IpSegmentTools.VerifyPort(req.Port) {
		server.RespFail(c, 4000, "端口号范围0-65535")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingSyslogSave(ctx, &req); err != nil {
		log.Error("SystemSettingSyslogSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingSyslogInfo 系统管理 - 系统设置 - syslog服务 - 信息
func SystemSettingSyslogInfo(c *gin.Context) {
	var res typespec.SystemSettingSyslogInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingSyslogInfo(ctx, &res); err != nil {
		log.Error("SystemSettingSyslogInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingMailSave 系统设置 - 邮箱配置 - 保存
func SystemSettingMailSave(c *gin.Context) {
	var (
		req typespec.SystemSettingMailSaveReq
		res typespec.SystemSettingMailSaveRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemSettingMailSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	//校验参数的格式
	//if !network.IpSegmentTools.VerifyIp(req.Address) {
	//	server.RespFail(c, 4000, "ip格式错误")
	//	return
	//}
	if !network.IpSegmentTools.VerifyPort(req.Port) {
		server.RespFail(c, 4000, "端口号范围0-65535")
		return
	}
	if !network.IpSegmentTools.VerifyMail(req.Username) {
		server.RespFail(c, 4000, "邮箱账号格式错误")
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingMailSave(ctx, &req); err != nil {
		log.Error("SystemSettingMailSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// SystemSettingMailInfo 系统设置 - 邮箱配置 - 信息
func SystemSettingMailInfo(c *gin.Context) {
	var res typespec.SystemSettingMailInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingMailInfo(ctx, &res); err != nil {
		log.Error("SystemSettingMailInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingMailVerify 系统设置 - 邮箱配置 - 验证
func SystemSettingMailVerify(c *gin.Context) {
	var (
		req typespec.SystemSettingMailSaveReq
		res typespec.SystemSettingMailSaveRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemSettingMailVerify parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingMailVerify(ctx, &req); err != nil {
		log.Error("SystemSettingMailVerify parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// SystemSettingNetworkConfigSave 系统管理 - 系统设置 - 网络配置 - 保存
func SystemSettingNetworkConfigSave(c *gin.Context) {
	var (
		req typespec.SystemSettingNetworkConfigSaveReq
		res typespec.SystemSettingNetworkConfigSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemSettingNetworkConfigSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//校验参数的格式
	if !network.IpSegmentTools.VerifyIp(req.Ip) {
		server.RespFail(c, 4000, "ip地址格式不正确")
		return
	}
	if !network.IpSegmentTools.VerifyIp(req.Mask) {
		server.RespFail(c, 4000, "子网掩码格式不正确")
		return
	}
	if !network.IpSegmentTools.VerifyIp(req.Gateway) {
		server.RespFail(c, 4000, "默认网关格式不正确")
		return
	}
	if !network.IpSegmentTools.VerifyIp(req.DnsServer) {
		server.RespFail(c, 4000, "DNS服务器格式不正确")
		return
	}
	if req.StandbyDnsServer != "" && !network.IpSegmentTools.VerifyIp(req.StandbyDnsServer) {
		server.RespFail(c, 4000, "备用DNS服务器格式不正确")
		return
	}
	if !network.IpSegmentTools.VerifyPort(req.WebPort) {
		server.RespFail(c, 4000, "端口号范围0-65535")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingNetworkConfigSave(ctx, &req); err != nil {
		log.Error("SystemSettingNetworkConfigSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingNetworkConfigInfo 系统管理 - 系统设置 - 网络配置 - 信息
func SystemSettingNetworkConfigInfo(c *gin.Context) {
	var res typespec.SystemSettingNetworkConfigInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingNetworkConfigInfo(ctx, &res); err != nil {
		log.Error("SystemSettingNetworkConfigInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingMonitorWarnInfo 系统管理 - 系统设置 - 系统监控告警 - 信息
func SystemSettingMonitorWarnInfo(c *gin.Context) {
	var res typespec.SystemMonitorWarnInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingMonitorWarnInfo(ctx, &res); err != nil {
		log.Error("SystemSettingMonitorWarnInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSettingMonitorWarnSave 系统管理 - 系统设置 - 系统监控告警 - 保存
func SystemSettingMonitorWarnSave(c *gin.Context) {
	var (
		req typespec.SystemMonitorWarnSaveReq
		res typespec.SystemMonitorWarnSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemSettingMonitorWarnSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemSettingMonitorWarnSave(ctx, &req); err != nil {
		log.Error("SystemSettingMonitorWarnSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BusinessSettingTcpBlindTestSave 系统管理 - 业务设置 - Tcp盲测平台 - 保存
func BusinessSettingTcpBlindTestSave(c *gin.Context) {
	var (
		req typespec.TcpBlindTestSaveReq
		res typespec.TcpBlindTestSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BusinessSettingTcpBlindTestSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//参数逻辑校验
	if req.Type == enums.TypeCustom && len(req.Host) == 0 { //自定义
		server.RespFail(c, 4000, "自定义类型,监听IP不能为空")
		return
	}
	//校验参数格式
	if req.Type == enums.TypeCustom && !network.IpSegmentTools.VerifyIp(req.Host) {
		server.RespFail(c, 4000, "ip格式错误")
		return
	}
	if !network.IpSegmentTools.VerifyPort(req.Port) {
		server.RespFail(c, 4000, "端口号范围0-65535")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.BusinessSettingTcpBlindTestSave(ctx, &req); err != nil {
		log.Error("BusinessSettingTcpBlindTestSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BusinessSettingTcpBlindTestInfo 系统管理 - 业务设置 - Tcp盲测平台 - 信息
func BusinessSettingTcpBlindTestInfo(c *gin.Context) {
	var res typespec.TcpBlindTestInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.BusinessSettingTcpBlindTestInfo(ctx, &res); err != nil {
		log.Error("BusinessSettingTcpBlindTestInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BusinessSettingHttpBlindTestSave 系统管理 - 业务设置 - http盲测平台 - 保存
func BusinessSettingHttpBlindTestSave(c *gin.Context) {
	var (
		req typespec.HttpBlindTestSaveReq
		res typespec.HttpBlindTestSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BusinessSettingHttpBlindTestSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//参数逻辑校验
	if req.Type == enums.TypeCustom && len(req.Host) == 0 { //自定义
		server.RespFail(c, 4000, "自定义类型,监听IP不能为空")
		return
	}
	//校验参数格式
	if req.Type == enums.TypeCustom && !network.IpSegmentTools.VerifyIp(req.Host) {
		server.RespFail(c, 4000, "ip格式错误")
		return
	}
	if !network.IpSegmentTools.VerifyPort(req.Port) {
		server.RespFail(c, 4000, "端口号范围0-65535")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.BusinessSettingHttpBlindTestSave(ctx, &req); err != nil {
		log.Error("BusinessSettingHttpBlindTestSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BusinessSettingHttpBlindTestInfo 系统管理 - 业务设置 - http盲测平台 - 信息
func BusinessSettingHttpBlindTestInfo(c *gin.Context) {
	var res typespec.HttpBlindTestInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.BusinessSettingHttpBlindTestInfo(ctx, &res); err != nil {
		log.Error("BusinessSettingHttpBlindTestInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BusinessSettingDnsBlindTestSave 系统管理 - 业务设置 - Dns盲测平台 - 保存
func BusinessSettingDnsBlindTestSave(c *gin.Context) {
	var (
		req typespec.DnsBlindTestSaveReq
		res typespec.DnsBlindTestSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BusinessSettingDnsBlindTestSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//参数逻辑校验
	if req.Type == enums.TypeCustom && len(req.Domain) == 0 { //自定义
		server.RespFail(c, 4000, "自定义类型,监听IP不能为空")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.BusinessSettingDnsBlindTestSave(ctx, &req); err != nil {
		log.Error("BusinessSettingDnsBlindTestSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BusinessSettingDnsBlindTestInfo 系统管理 - 业务设置 - Dns盲测平台 - 信息
func BusinessSettingDnsBlindTestInfo(c *gin.Context) {
	var res typespec.DnsBlindTestInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.BusinessSettingDnsBlindTestInfo(ctx, &res); err != nil {
		log.Error("BusinessSettingDnsBlindTestInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BusinessSettingIcmpBlindTestSave 系统管理 - 业务设置 - Icmp盲测平台 - 保存
func BusinessSettingIcmpBlindTestSave(c *gin.Context) {
	var (
		req typespec.IcmpBlindTestSaveReq
		res typespec.IcmpBlindTestSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BusinessSettingIcmpBlindTestSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//参数逻辑校验
	if req.Type == enums.TypeCustom && len(req.Host) == 0 { //自定义
		server.RespFail(c, 4000, "自定义类型,监听IP不能为空")
		return
	}
	//校验参数格式
	if req.Type == enums.TypeCustom && !network.IpSegmentTools.VerifyIp(req.Host) {
		server.RespFail(c, 4000, "ip格式错误")
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.BusinessSettingIcmpBlindTestSave(ctx, &req); err != nil {
		log.Error("BusinessSettingIcmpBlindTestSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BusinessSettingIcmpBlindTestInfo 系统管理 - 业务设置 - Icmp盲测平台 - 信息
func BusinessSettingIcmpBlindTestInfo(c *gin.Context) {
	var res typespec.IcmpBlindTestInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.BusinessSettingIcmpBlindTestInfo(ctx, &res); err != nil {
		log.Error("BusinessSettingIcmpBlindTestInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// CurTasksInfo 业务设置 - 任务并发配置 - 信息
func CurTasksInfo(c *gin.Context) {
	var res typespec.CurTasksInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.CurTasksInfo(ctx, &res); err != nil {
		log.Error("CurTasksInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// CurTasksSave 业务设置 - 任务并发配置 - 保存
func CurTasksSave(c *gin.Context) {
	var (
		req  typespec.CurTasksSaveReq
		resp typespec.CurTasksSaveResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("CurTasksSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.CurTasksSave(ctx, &req); err != nil {
		log.Error("CurTasksSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// UseScoreInfo 业务设置 - 可以利用评分 - 信息
func UseScoreInfo(c *gin.Context) {
	var res typespec.UseScoreInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.UseScoreInfo(ctx, &res); err != nil {
		log.Error("UseScoreInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// UseScoreSave 业务设置 - 可以利用评分 - 保存
func UseScoreSave(c *gin.Context) {
	var (
		req typespec.UseScoreSaveReq
		res typespec.UseScoreSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UseScoreSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.UseScoreSave(ctx, &req); err != nil {
		log.Error("UseScoreSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TestScopeInfo 业务设置 - 测试范围校验开关 - 信息
func TestScopeInfo(c *gin.Context) {
	var res typespec.TestScopeInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.TestScopeInfo(ctx, &res); err != nil {
		log.Error("TestScopeInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TestScopeSave 业务设置 - 测试范围校验开关 - 保存
func TestScopeSave(c *gin.Context) {
	var (
		req typespec.TestScopeSaveReq
		res typespec.TestScopeSaveRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TestScopeSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.TestScopeSave(ctx, &req); err != nil {
		log.Error("TestScopeSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// CpuInfo 系统管理 - 系统监控 - cpu
func CpuInfo(c *gin.Context) {
	var res typespec.CpuInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.CpuInfo(ctx, &res); err != nil {
		log.Error("CpuInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// MemoryInfo 系统管理 - 系统监控 - 内存
func MemoryInfo(c *gin.Context) {
	var res typespec.MemoryInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.MemoryInfo(ctx, &res); err != nil {
		log.Error("MemoryInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// DiskInfo 系统管理 - 系统监控 - 磁盘
func DiskInfo(c *gin.Context) {
	var res typespec.DiskInfoRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.DiskInfo(ctx, &res); err != nil {
		log.Error("DiskInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// RouteList 系统管理 - 系统设置 - 路由配置 - 路由列表
func RouteList(c *gin.Context) {
	var (
		res typespec.SystemRouteListRes
		app application.SystemManage
	)
	if err := app.RouteList(&res); err != nil {
		log.Error("RouteList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// RouteAdd 系统管理 - 系统设置 - 路由配置 - 增加路由
func RouteAdd(c *gin.Context) {
	var (
		req typespec.SystemRouteAddReq
		res typespec.SystemRouteAddRes
		app application.SystemManage
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("RouteAdd parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	//校验参数格式
	if !network.IpSegmentTools.VerifyIp(req.Ip) {
		server.RespFail(c, 4000, "ip格式错误")
		return
	}
	if !network.IpSegmentTools.VerifyIp(req.Netmask) {
		server.RespFail(c, 4000, "掩码格式错误")
		return
	}
	if !network.IpSegmentTools.VerifyIp(req.Gateway) {
		server.RespFail(c, 4000, "网关格式错误")
		return
	}
	if err := app.RouteAdd(&req); err != nil {
		log.Error("RouteAdd parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// RouteDelete 系统管理 - 系统设置 - 路由配置 - 删除路由
func RouteDelete(c *gin.Context) {
	var (
		req typespec.SystemRouteDeleteReq
		res typespec.SystemRouteDeleteRes
		app application.SystemManage
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("RouteDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	if err := app.RouteDelete(&req); err != nil {
		log.Error("RouteDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemVersion 系统管理 - 升级还原 - 系统版本信息
func SystemVersion(c *gin.Context) {
	var res typespec.SystemVersionRes
	ctx := server.NewContext(context.Background(), c)
	var app application.SystemManage
	if err := app.SystemVersion(ctx, &res); err != nil {
		log.Error("SystemVersion parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// GetUpgradeStatus 系统管理 - 升级还原 - 获取升级进度
func GetUpgradeStatus(c *gin.Context) {
	status := services.GetUpgradeStatus()
	server.RespSuccess(c, status)
}

// UploadUpgradeFile 系统管理 - 升级还原 - 上传系统更新文件 (新版)
func UploadUpgradeFile(c *gin.Context) {
	// 1. 获取上传文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		server.RespFail(c, 4000, "获取文件失败: "+err.Error())
		return
	}

	// 2. 基础校验
	if !strings.HasSuffix(fileHeader.Filename, ".zip") {
		server.RespFail(c, 4000, "仅支持 .zip 格式的升级包")
		return
	}

	// 3. 准备保存路径
	// 每次升级都建议使用一个新的临时目录，或者清理旧目录
	// 这里沿用 enums.SystemUpgradeFileDir，但建议确保它是干净的
	if err := os.RemoveAll(enums.SystemUpgradeFileDir); err != nil {
		log.Errorf("清理升级目录失败: %v", err)
		server.RespFail(c, 4000, "清理升级目录失败")
		return
	}
	if err := file.CreateDir(enums.SystemUpgradeFileDir); err != nil {
		server.RespFail(c, 4000, "创建升级目录失败")
		return
	}
	savePath := filepath.Join(enums.SystemUpgradeFileDir, fileHeader.Filename)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		log.Errorf("保存升级包失败: %v", err)
		server.RespFail(c, 4000, "保存升级包失败")
		return
	}

	mgr := services.NewUpgradeManager()
	manifest, err := mgr.PrepareUpgrade(c.Request.Context(), savePath)
	if err != nil {
		log.Errorf("升级包校验失败: %v", err)
		server.RespFail(c, 4000, "升级包校验失败: "+err.Error())
		return
	}

	// 辅助转换中文描述
	var typeDesc, scopeDesc string
	switch manifest.Type {
	case services.UpgradeTypeSystem:
		typeDesc = "系统升级包"
	case services.UpgradeTypeVuln:
		typeDesc = "漏洞/脚本库升级包"
	default:
		typeDesc = "未知类型"
	}

	if manifest.Type == services.UpgradeTypeVuln {
		switch manifest.VulnScope {
		case services.VulnScopePrinciple:
			scopeDesc = "快速更新包"
		case services.VulnScopeFull:
			scopeDesc = "全量更新包"
		case services.VulnScopeScript:
			scopeDesc = "仅脚本库更新"
		case services.VulnScopeMixed:
			scopeDesc = "混合模式更新"
		default:
			scopeDesc = "未知范围"
		}
	}

	// 返回 manifest 信息供前端确认
	server.RespSuccess(c, gin.H{
		"filename":    fileHeader.Filename,
		"version":     manifest.Version,
		"type":        manifest.Type,
		"typeDesc":    typeDesc, // 新增类型中文描述
		"vulnScope":   manifest.VulnScope,
		"scopeDesc":   scopeDesc, // 新增范围中文描述
		"description": manifest.Description,
		"needRestart": manifest.NeedRestart,
		"buildTime":   manifest.BuildTime,
	})
}

// ConfirmUpgrade 确认升级
func ConfirmUpgrade(c *gin.Context) {
	var req typespec.ConfirmUpgradeReq
	if err := c.ShouldBind(&req); err != nil {
		server.RespFail(c, 4000, "参数错误: "+err.Error())
		return
	}

	zipPath := filepath.Join(enums.SystemUpgradeFileDir, req.Filename)
	if !file.CheckPathExist(zipPath) {
		server.RespFail(c, 4000, "升级包不存在或已过期，请重新上传")
		return
	}

	mgr := services.NewUpgradeManager()
	if err := mgr.ExecuteUpgrade(c.Request.Context(), zipPath); err != nil {
		log.Errorf("启动升级失败: %v", err)
		server.RespFail(c, 4000, "启动升级失败: "+err.Error())
		return
	}

	server.RespSuccess(c, "升级任务已启动")
}

// GenerateToken 系统管理 - API秘钥 — 生成秘钥
func GenerateToken(c *gin.Context) {
	var (
		req  typespec.GenerateTokenReq
		resp typespec.GenerateTokenResp
		app  application.SystemManage
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("GenerateToken parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	if err := app.GeneratePlatformToken(ctx, &req, &resp); err != nil {
		log.Error("GenerateToken parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TokenList 系统管理 - API秘钥 — 用户token列表
func TokenList(c *gin.Context) {
	var (
		req  typespec.TokenListReq
		resp typespec.TokenListResp
		app  application.SystemManage
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("TokenList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	if err := app.TokenList(ctx, &req, &resp); err != nil {
		log.Error("GenerateToken parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// SystemNodeDownload 节点管理 - 安装包下载
func SystemNodeDownload(c *gin.Context) {
	var (
		req typespec.SystemNodeDownloadReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemNodeDownload parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 获取当前目录
	pwd, err := os.Getwd()
	if err != nil {
		log.Error("SystemNodeDownload get pwd error: " + err.Error())
	} else {
		log.Info("SystemNodeDownload current directory: " + pwd)
	}

	filename := "scanner"
	if req.Os == "windows" {
		filename = "scanner.exe"
	}

	// 优先在当前目录查找
	downloadPath := filepath.Join(pwd, filename)
	if !file.CheckPathExist(downloadPath) {
		// 其次在项目目录查找
		downloadPath = filepath.Join(enums.SystemUpgradeProjectDir, "smart", filename)
		if !file.CheckPathExist(downloadPath) {
			server.RespFail(c, 4000, "agent文件不存在")
			return
		}
	}

	c.Header("Content-Type", "application/octet-stream")             //唤起浏览器的下载
	c.Header("Content-Disposition", "attachment;filename="+filename) //设置文件名
	c.File(downloadPath)
}

// TokenDelete 系统管理 - API秘钥 — 删除token接口
func TokenDelete(c *gin.Context) {
	var (
		req  typespec.TokenDeleteReq
		resp typespec.TokenDeleteResp
		app  application.SystemManage
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("TokenList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	if err := app.TokenDelete(ctx, &req, &resp); err != nil {
		log.Error("GenerateToken parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
