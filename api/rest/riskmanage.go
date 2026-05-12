package rest

import (
	"smart/api/typespec"
	"smart/application"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

// VulRiskStatistics 漏洞风险统计
func VulRiskStatistics(c *gin.Context) {
	var resp typespec.VulRiskStaticsRes
	// ctx := server.NewContext(context.Background(), c)
	var app application.RiskManageApp
	if err := app.VulRiskStatistics(c, &resp); err != nil {
		server.RespFail(c, 400, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulRiskList 风险管理-漏洞风险-列表
func VulRiskList(c *gin.Context) {
	var (
		req  typespec.VulRiskListReq
		resp typespec.VulRiskListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulRiskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	// ctx := server.NewContext(context.Background(), c)
	var app application.RiskManageApp
	if err := app.VulRiskList(c, &req, &resp); err != nil {
		log.Error("VulRiskList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// DelVulRisk 风险管理-漏洞风险-删除
func DelVulRisk(c *gin.Context) {
	var req typespec.VulRiskDelReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("DelVulRisk parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	// ctx := server.NewContext(context.Background(), c)
	var app application.RiskManageApp
	if err := app.VulRiskDel(c, &req); err != nil {
		log.Error("VulRiskDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "success")
}

// UpdateVulRisk 风险管理-漏洞风险-更新
func UpdateVulRisk(c *gin.Context) {
	var req typespec.VulRiskUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UpdateVulRisk parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	// ctx := server.NewContext(context.Background(), c)
	var app application.RiskManageApp
	if err := app.VulRiskUpdate(c, &req); err != nil {
		log.Error("UpdateVulRisk application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "")
}

// VulRiskDetail 风险管理-漏洞风险-详情
func VulRiskDetail(c *gin.Context) {
	var (
		req  typespec.VulRiskDetailReq
		resp typespec.VulRiskInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulRiskDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	// ctx := server.NewContext(context.Background(), c)
	var app application.RiskManageApp
	if err := app.VulRiskDetail(c, &req, &resp); err != nil {
		log.Error("VulRiskDetail application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulRiskTest 风险管理-漏洞风险-测试
func VulRiskTest(c *gin.Context) {
	var (
		req  typespec.VulTestReq
		resp typespec.VulTestResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulRiskTest parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	// ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.VulTest(c, &req, &resp); err != nil {
		log.Error("VulRiskTest application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulRiskVerify 风险管理-漏洞风险-验证
func VulRiskVerify(c *gin.Context) {
	var req typespec.VulRiskVerifyReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulRiskVulVerify parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误!!")
		return
	}
	// ctx := server.NewContext(context.Background(), c)
	var app application.RiskManageApp
	if err := app.VulRiskVulVerify(c, &req); err != nil {
		log.Error("VulRiskVulVerify application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "success")
}

// RiskManageEnum 风险管理枚举
func RiskManageEnum(c *gin.Context) {
	var res typespec.RiskManageEnumsRes
	// ctx := server.NewContext(context.Background(), c)
	var app application.RiskManageApp
	if err := app.RiskManageEnums(c, &res); err != nil {
		log.Error("RiskManageEnum parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}
