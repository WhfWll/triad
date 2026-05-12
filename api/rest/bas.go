package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gitlabee.4dogs.cn/common/redis"
	"gitlabee.4dogs.cn/common/server"
	"path"
	"smart/api/typespec"
	"smart/application"
	"smart/tools/file"
	"strings"
	"time"
)

// BAS需求

// 规则导入
func BasRuleImport(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因1："+err.Error())
		return
	}
	if file == nil {
		server.RespFail(c, 4000, "请上传文件")
		return
	}
	filename := path.Join("/var/tmp/", file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因2："+err.Error())
		return
	}

	excel, err := excelize.OpenFile(filename)
	if err != nil {
		server.RespFail(c, 4000, "参数错误,错误原因3："+err.Error())
		return
	}
	defer excel.Close()

	ctx := server.NewContext(context.Background(), c)

	var app application.Bas
	if err := app.BasRuleImport(ctx, excel); err != nil {
		log.Error("BasRuleImport parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, struct{}{})
}

// 规则枚举
func BasRuleEnum(c *gin.Context) {
	var (
		res typespec.BasRuleEnumRes
	)

	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	app.BasRuleEnum(ctx, &res)
	server.RespSuccess(c, res)
}

// 规则查询列表
func BasRuleGet(c *gin.Context) {
	var (
		req typespec.BasRuleListReq
		res typespec.BasRuleListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasRuleGet parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasRuleGet(ctx, &req, &res); err != nil {
		log.Error("BasRuleGet parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 规则详情
func BasRuleInfo(c *gin.Context) {
	var (
		req typespec.BasRuleInfoReq
		res typespec.BasRuleInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasRuleInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasRuleInfo(ctx, &req, &res); err != nil {
		log.Error("BasRuleInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 规则编辑
func BasRuleEdit(c *gin.Context) {
	var (
		req typespec.BasRuleEditReq
		res typespec.BasRuleEditRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasRuleEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	if req.AttackMode == 0 {
		server.RespFail(c, 4000, "请选择攻击方式")
		return
	}
	if req.AttackStage == 0 {
		server.RespFail(c, 4000, "请选择攻击阶段")
		return
	}
	if req.RiskLevel == 0 {
		server.RespFail(c, 4000, "请选择影响级别")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasRuleEdit(ctx, &req, &res); err != nil {
		log.Error("BasRuleEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 剧本集创建
func BasTemplateCreate(c *gin.Context) {
	var (
		req typespec.BasTemplateCreateReq
		res typespec.BasTemplateCreateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasTemplateCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasTemplateCreate(ctx, &req); err != nil {
		log.Error("BasTemplateCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 剧本集列表
func BasTemplateGet(c *gin.Context) {
	var (
		req typespec.BasTemplateListReq
		res typespec.BasTemplateListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasTemplateGet parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasTemplateGet(ctx, &req, &res); err != nil {
		log.Error("BasTemplateGet parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 剧本集详情
func BasGetTemplateById(c *gin.Context) {
	var (
		req typespec.BasTemplateGetReq
		res typespec.BasTemplateGetRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasGetTemplateById parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasGetTemplateById(ctx, &req, &res); err != nil {
		log.Error("BasGetTemplateById parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 剧本集删除
func BasDelTemplateById(c *gin.Context) {
	var (
		req typespec.BasTemplateDelReq
		res typespec.BasTemplateDelRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasDelTemplateById parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasDelTemplateById(ctx, &req); err != nil {
		log.Error("BasDelTemplateById parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 剧本集设置是否默认
func BasTemplateSetDefault(c *gin.Context) {
	var (
		req typespec.BasTemplateSetDefaultReq
		res typespec.BasTemplateSetDefaultRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasTemplateSetDefault parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasTemplateSetDefault(ctx, &req); err != nil {
		log.Error("BasTemplateSetDefault parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BAS任务 检测agent是否在线
func BasAgentIsOnline(c *gin.Context) {
	var (
		req typespec.BasAgentIsOnlineReq
		res typespec.BasAgentIsOnlineRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasAgentIsOnline parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasAgentIsOnline(ctx, &req); err != nil {
		log.Error("BasAgentIsOnline parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BAS任务列表
func BasGetTask(c *gin.Context) {
	var (
		req typespec.BasTaskGetReq
		res typespec.BasTaskGetRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasGetTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasGetTask(ctx, &req, &res); err != nil {
		log.Error("BasGetTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BAS任务结束
func BasEndTaskById(c *gin.Context) {
	var (
		req typespec.BasTaskEndReq
		res typespec.BasTaskEndRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasEndTaskById parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasEndTaskById(ctx, &req); err != nil {
		log.Error("BasEndTaskById parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BAS任务删除
func BasDelTask(c *gin.Context) {
	var (
		req typespec.BasTaskDelReq
		res typespec.BasTaskDelRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasDelTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasDelTask(ctx, &req); err != nil {
		log.Error("BasDelTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BAS任务详情
func BasGetTaskTargetPage(c *gin.Context) {
	var (
		req typespec.BasTaskDetailReq
		res typespec.BasTaskDetailRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasGetTaskTargetPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasGetTaskTargetPage(ctx, &req, &res); err != nil {
		log.Error("BasGetTaskTargetPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BAS任务目标日志
func BasGetTargetLogs(c *gin.Context) {
	var (
		req typespec.BasTargetLogReq
		res typespec.BasTargetLogRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasGetTargetLogs parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasGetTargetLogs(ctx, &req, &res); err != nil {
		log.Error("BasGetTargetLogs parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BAS任务目标删除
func BasGetTargetDel(c *gin.Context) {
	var (
		req typespec.BasTargetDelReq
		res typespec.BasTargetDelRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasGetTargetDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasTargetDel(ctx, &req); err != nil {
		log.Error("BasGetTargetDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// BAS枚举接口
func BasEnum(c *gin.Context) {
	var resp typespec.BasEnumRes
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasEnum(ctx, &resp); err != nil {
		log.Error("BasEnum application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// BAS任务创建
func BasCreateTask(c *gin.Context) {
	var (
		req  typespec.BasTaskCreateReq
		resp typespec.BasTaskCreateResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasCreateTask parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	if len(req.BasNodeIds) == 0 {
		server.RespFail(c, 4000, "所选攻击目标不能为空...")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	uid, ok := c.Get("uid")
	if !ok {
		return
	}
	var app application.Bas
	if err := app.BasCreateTask(ctx, &req, uid.(int)); err != nil {
		log.Error("BasCreateTask application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// bas心跳及检测结果接收
func BasReceivResult(c *gin.Context) {
	var (
		req  typespec.BasReceivResultReq
		resp typespec.BasReceivResultResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasReceivResult parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasReceivResult(ctx, &req); err != nil {
		log.Error("BasReceivResult application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 漏洞测试统计
func BasVulStat(c *gin.Context) {
	var (
		req  typespec.BasVulStatReq
		resp typespec.BasVulStatResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasVulStat parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasVulStat(ctx, &req, &resp); err != nil {
		log.Error("BasVulStat application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 漏洞测试列表
func BasVulList(c *gin.Context) {
	var (
		req  typespec.BasVulListReq
		resp typespec.BasVulListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasVulList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasVulList(ctx, &req, &resp); err != nil {
		log.Error("BasVulList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// 漏洞测试删除
func BasVulDel(c *gin.Context) {
	var (
		req  typespec.BasVulDelReq
		resp typespec.BasVulDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasVulDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasVulDel(ctx, &req); err != nil {
		log.Error("BasVulDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// Bas agent下载
func BasAgentDownload(c *gin.Context) {
	var (
		req typespec.BasAgentDownloadReq
		res typespec.BasAgentDownloadRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasAgentList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	if req.GetTempToken {
		tempToken := uuid.New().String()
		// 获取临时token
		cacheClient, _ := redis.NewClient()
		cacheClient.Set(c, tempToken, "Y", 1*time.Hour)
		res.TempToken = tempToken
		server.RespSuccess(c, res)
		return
	}

	// Agent 所在的固定位置
	agentFile := "/opt/laozhi/smart/bas-agent-linux"
	if req.Platform == "windows" {
		agentFile = "/opt/laozhi/smart/bas-agent-windows.exe"
	}

	//判断服务器端的文件是否存在（因其他原因导致的文件不存在）
	if !file.CheckPathExist(agentFile) {
		log.Errorf("agent file (%s) is not exist", agentFile)
		server.RespFail(c, 4000, "BAS agent文件不存在")
		return
	}
	filepathList := strings.Split(agentFile, "/")
	filename := filepathList[len(filepathList)-1]
	c.Header("Content-Type", "application/octet-stream")             //唤起浏览器的下载
	c.Header("Content-Disposition", "attachment;filename="+filename) //设置文件名
	c.File(agentFile)
}

// Bas agent列表
func BasAgentList(c *gin.Context) {
	var (
		req typespec.BasAgentListReq
		res typespec.BasAgentListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasAgentList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasAgentList(ctx, &req, &res); err != nil {
		log.Error("BasAgentList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// 可用节点列表
func BasAgentLive(c *gin.Context) {
	var res typespec.BasAgentLiveResp
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasAgentLive(ctx, &res); err != nil {
		log.Error("BasAgentLive application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// Bas agent列表
func BasAgentStatusEdit(c *gin.Context) {
	var (
		req typespec.BasAgentStatusEditReq
		res typespec.BasAgentStatusEditRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("BasAgentStatusEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.Bas
	if err := app.BasAgentStatusEdit(ctx, &req, &res); err != nil {
		log.Error("BasAgentStatusEdit application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}
