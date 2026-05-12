package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// 字典管理

// CreateSystemSecurity 系统管理 - 系统配置 - 安全配置
func CreateSystemSecurity(c *gin.Context) {
	var (
		req typespec.MapSetSystemSecurityReq
		res typespec.MapSetSystemSecurityRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("CreateSystemSecurity parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.MapSet
	if err := app.CreateSystemSecurity(ctx, &req, &res); err != nil {
		log.Error("CreateSystemSecurity parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// SystemSecurityInfo 安全配置信息
func SystemSecurityInfo(c *gin.Context) {
	var (
		req typespec.SystemSecurityInfoReq
		res typespec.SystemSecurityInfoRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("SystemSecurityInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.MapSet
	if err := app.GetSystemSecurityInfo(ctx, &req, &res); err != nil {
		log.Error("SystemSecurityInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}
