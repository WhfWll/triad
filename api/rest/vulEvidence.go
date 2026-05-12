package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// VulEvidenceList 漏洞取证列表及筛选
func VulEvidenceList(c *gin.Context) {
	var (
		req  typespec.VulEvidenceListReq
		resp typespec.VulEvidenceListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulEvidenceList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.VulEvidenceList(ctx, &req, &resp); err != nil {
		log.Error("VulEvidenceList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulEvidenceDel 漏洞取证列表删除
func VulEvidenceDel(c *gin.Context) {
	var req typespec.VulEvidenceDelReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulEvidenceDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.DelVulEvidence(ctx, &req); err != nil {
		log.Error("DelVulEvidence parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "del success")
}

// VulEvidenceInfo 漏洞取证详情信息
func VulEvidenceInfo(c *gin.Context) {
	var (
		req  typespec.VulEvidenceInfoReq
		resp typespec.VulEvidenceInfoRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulEvidenceInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.VulEvidenceInfo(ctx, &req, &resp); err != nil {
		log.Error("VulEvidenceInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// RiskTypeInfoEnum 风险类型信息枚举
func RiskTypeInfoEnum(c *gin.Context) {
	var resp map[int]string
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.RiskTypeInfoEnum(ctx, &resp); err != nil {
		log.Error("CaptureInfoEnum parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulEvidenceUse 漏洞取证利用
func VulEvidenceUse(c *gin.Context) {
	var (
		req  typespec.EvidenceUseReq
		resp typespec.EvidenceUseRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulEvidenceUse parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.VulEvidenceListApp
	if err := app.VulEvidenceUse(ctx, &req, &resp); err != nil {
		log.Error("VulEvidenceInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
