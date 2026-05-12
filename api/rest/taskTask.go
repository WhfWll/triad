package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
	"smart/tools/data"
	"smart/tools/enums"
	"smart/tools/utils"
	"strings"
	"time"
	"unicode/utf8"
)

// 任务中心 - 渗透任务

// TaskTaskEnum 创建任务枚举
func TaskTaskEnum(c *gin.Context) {
	var res typespec.TaskTaskEnumRes

	// 调用聚合层
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.Enum(ctx, &res); err != nil {
		log.Error("TaskTaskEnum application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TaskTaskWebLoginCheck 网站登陆凭证校验
func TaskTaskWebLoginCheck(c *gin.Context) {
	var (
		req typespec.TaskTaskWebLoginCheckReq
		res typespec.TaskTaskWebLoginCheckRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTaskWebLoginCheck parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	// 调用聚合层
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTaskWebLoginCheck(ctx, &req, &res); err != nil {
		log.Error("TaskTaskWebLoginCheck application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TaskSave 创建
func TaskSave(c *gin.Context) {
	var (
		req  typespec.TaskSaveReq
		resp typespec.TaskSaveResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	// 任务名称长度校验 最多50个字符，超过截取后面的字符
	if utf8.RuneCountInString(req.TaskName) > 50 {
		taskNameRune := []rune(req.TaskName)
		req.TaskName = string(taskNameRune[:50])
	}
	// 运行时间段 校验
	var runtimeCheck data.TaskRuntimeCheck
	runtimeCheck.StartTime = req.ExecuteJson.StartTime
	runtimeCheck.CyclePlanningType = req.ExecuteJson.CyclePlanningType
	runtimeCheck.CyclePlanningValue = req.ExecuteJson.CyclePlanningValue
	runtimeCheck.CyclePlanningHour = req.ExecuteJson.CyclePlanningHour
	runtimeCheck.EndTime = req.ExecuteJson.EndTime
	runtimeCheck.RuntimePeriod = req.ExecuteJson.RuntimePeriod
	if err := runtimeCheck.CheckRuntime(req.ExecuteType); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}
	// 调用聚合层
	//ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	ip := utils.GetClientIp(c)
	err := app.Save(c, &req, &resp)
	if aErr := app.SendTaskInfoToAuditLog(c, &req, err, ip); aErr != nil {
		err = aErr
	}
	if mErr := app.SendTaskInfoToMessage(c, &req, err); mErr != nil {
		err = mErr
	}
	if err != nil {
		log.Error("TaskSave application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskTaskList 列表
func TaskTaskList(c *gin.Context) {
	var (
		req typespec.TaskListReq
		res typespec.TaskListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTaskList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	if req.StartTime != "" {
		if _, err := time.Parse(utils.DateTime, req.StartTime); err != nil {
			server.RespFail(c, 4000, "参数错误,错误原因：任务开始时间错误："+err.Error())
			return
		}
	}
	if req.EndTime != "" {
		if _, err := time.Parse(utils.DateTime, req.EndTime); err != nil {
			server.RespFail(c, 4000, "参数错误,错误原因：任务结束时间错误："+err.Error())
			return
		}
	}

	// 调用聚合层
	//ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTaskList(c, &req, &res); err != nil {
		log.Error("TaskTaskList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TaskTaskDel 删除
func TaskTaskDel(c *gin.Context) {
	var (
		req typespec.TaskDelReq
		res typespec.TaskDelRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTaskDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	// 调用聚合层
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTaskDel(ctx, &req); err != nil {
		log.Error("TaskTaskDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TaskTaskCopy 任务拷贝 注意这里不是创建，仅获取需要拷贝的任务，将数据返回
func TaskTaskCopy(c *gin.Context) {
	var (
		req typespec.TaskCopyReq
		res typespec.TaskCopyRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTaskCopy parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	// 调用聚合层
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTaskCopy(ctx, &req, &res); err != nil {
		log.Error("TaskTaskCopy application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TaskTaskChangeState 修改任务状态
func TaskTaskChangeState(c *gin.Context) {
	var (
		req typespec.TaskChangeStateReq
		res typespec.TaskChangeStateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTaskChangeState parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	if req.Operate != enums.TaskOperatePause && req.Operate != enums.TaskOperateResume && req.Operate != enums.TaskOperateStop {
		server.RespFail(c, 4000, "参数错误,错误原因：未知的操作")
		return
	}

	// 调用聚合层
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTaskChangeState(ctx, &req); err != nil {
		log.Error("TaskTaskChangeState application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// GetState 获取任务状态
func GetState(c *gin.Context) {
	var (
		req  typespec.GetStateReq
		resp typespec.GetStateResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("GetState parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.GetState(ctx, &req, &resp); err != nil {
		log.Error("GetState application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// OverView 任务概览
func OverView(c *gin.Context) {
	var (
		req  typespec.OverViewReq
		resp typespec.OverViewResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("OverView parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.OverView(ctx, &req, &resp); err != nil {
		log.Error("OverView application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskConfigInfo 渗透任务 - 任务配置信息
func TaskConfigInfo(c *gin.Context) {
	var (
		req typespec.TaskConfigInfoReq
		res typespec.TaskConfigInfoRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskConfigInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskConfigInfo(ctx, &req, &res); err != nil {
		log.Error("TaskConfigInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TargetChangeState 修改测试目标状态
func TargetChangeState(c *gin.Context) {
	var (
		req  typespec.TargetChangeStateReq
		resp typespec.TargetChangeStateResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TargetChangeState parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	if req.Operate != "stop" {
		server.RespFail(c, 4000, "不支持的操纵")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TargetChangeState(ctx, &req); err != nil {
		log.Error("TargetChangeState application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TargetList 测试目标列表及筛选
func TargetList(c *gin.Context) {
	var (
		req  typespec.TargetListReq
		resp typespec.TargetListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TargetList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TargetList(ctx, &req, &resp); err != nil {
		log.Error("TargetList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// UpdateTargetUseScore 批量修改目标的利用评分和状态
func UpdateTargetUseScore(c *gin.Context) {
	var (
		req  typespec.UpdateTargetUseScoreReq
		resp typespec.UpdateTargetUseScoreResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UpdateTargetUseScore parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.UpdateTargetUseScore(ctx, &req); err != nil {
		log.Error("UpdateTargetUseScore application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TargetDel 测试目标删除
func TargetDel(c *gin.Context) {
	var (
		req  typespec.TargetDelReq
		resp typespec.TargetDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TargetDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TargetDel(ctx, &req); err != nil {
		log.Error("TargetDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskResultList 信息收集列表及筛选
func TaskResultList(c *gin.Context) {
	var (
		req  typespec.TaskResultListReq
		resp typespec.TaskResultListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskResultList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskResultList(ctx, &req, &resp); err != nil {
		log.Error("TaskResultList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskResultUrlTree url树  销售许可证使用
func TaskResultUrlTree(c *gin.Context) {
	var (
		req  typespec.TaskResultUrlTreeReq
		resp typespec.TaskResultUrlTreeResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskResultList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskResultUrlTree(ctx, &req, &resp); err != nil {
		log.Error("TaskResultList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskResultDel 信息收集删除
func TaskResultDel(c *gin.Context) {
	var (
		req  typespec.TaskResultDelReq
		resp typespec.TaskResultDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskResultDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskResultDel(ctx, &req); err != nil {
		log.Error("TaskResultDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskResultDetail 信息收集查看详情接口
func TaskResultDetail(c *gin.Context) {
	var (
		req  typespec.TaskResultDetailReq
		resp typespec.TaskResultDetailResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskResultDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskResultDetail(ctx, &req, &resp); err != nil {
		log.Error("TaskResultDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulList 漏洞测试列表及筛选
func VulList(c *gin.Context) {
	var (
		req  typespec.VulListReq
		resp typespec.VulListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.VulList(ctx, &req, &resp); err != nil {
		log.Error("VulList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulInfo 漏洞测试详情
func VulInfo(c *gin.Context) {
	var (
		req  typespec.VulInfoReq
		resp typespec.VulInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.VulInfo(ctx, &req, &resp); err != nil {
		log.Error("VulInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// GetSnapshot 查看漏洞截图
func GetVulSnapshot(c *gin.Context) {
	var (
		req  typespec.GetVulSnapshotReq
		resp typespec.GetVulSnapshotResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("GetVulSnapshot parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.GetVulSnapshot(ctx, &req, &resp); err != nil {
		log.Error("GetVulSnapshot application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulDel 漏洞测试删除
func VulDel(c *gin.Context) {
	var (
		req  typespec.VulDelReq
		resp typespec.VulDelResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulDel parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.VulDel(ctx, &req); err != nil {
		log.Error("VulDel application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulTest 漏洞测试测试
func VulTest(c *gin.Context) {
	var (
		req  typespec.VulTestReq
		resp typespec.VulTestResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulTestResp parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.VulTestByYak(ctx, &req, &resp); err != nil {
		log.Error("VulTestResp application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// VulVerify 漏洞测试验证
func VulVerify(c *gin.Context) {
	var (
		req  typespec.VulVerifyReq
		resp typespec.VulVerifyResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulVerify parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.VulVerify(ctx, &req); err != nil {
		log.Error("VulVerify application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AsyncVulVerify 异步漏洞验证按钮
func AsyncVulVerify(c *gin.Context) {
	var (
		req  typespec.VulVerifyReq
		resp typespec.VulVerifyResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulVerify parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.VulVerify(ctx, &req); err != nil {
		log.Error("VulVerify application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TestVulTest 待测漏洞测试
func TestVulTest(c *gin.Context) {
	var (
		req  typespec.TestVulTestReq
		resp typespec.TestVulTestResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TestVulTest parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	taskVulIds := strings.Split(req.TaskVulIds, ",")
	if len(taskVulIds) == 0 {
		server.RespFail(c, 4000, "漏洞id不能为空")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TestVulTest(ctx, &req, taskVulIds); err != nil {
		log.Error("TestVulTest application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AttackLink 攻击链路图
func AttackLink(c *gin.Context) {
	var (
		req  typespec.AttackLinkReq
		resp typespec.AttackLinkResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("FlowTaskExport parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.AttackLink(ctx, &req, &resp); err != nil {
		log.Error("AttackLink application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// LogList 测试日志列表及筛选
func LogList(c *gin.Context) {
	var (
		req  typespec.LogListReq
		resp typespec.LogListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("LogList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.LogList(ctx, &req, &resp); err != nil {
		log.Error("LogList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// LogInfo 测试日志详情
func LogInfo(c *gin.Context) {
	var (
		req  typespec.LogInfoReq
		resp typespec.LogInfoResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("LogInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.LogInfo(ctx, &req, &resp); err != nil {
		log.Error("LogInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskTaskApiSave 任务创建，作为第三方接口被调用使用
func TaskTaskApiSave(c *gin.Context) {
	var (
		req typespec.ApiTaskTaskCreateReq
		res typespec.TaskSaveResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTaskSave parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	// 校验用户是否存在
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("TaskTaskSave parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}

	// 运行时间段 校验
	var runtimeCheck data.TaskRuntimeCheck
	if req.ExecuteJson.RuntimePeriod != nil {
		runtimeCheck.StartTime = req.ExecuteJson.StartTime
		runtimeCheck.CyclePlanningType = req.ExecuteJson.CyclePlanningType
		runtimeCheck.CyclePlanningValue = req.ExecuteJson.CyclePlanningValue
		runtimeCheck.CyclePlanningHour = req.ExecuteJson.CyclePlanningHour
		runtimeCheck.EndTime = req.ExecuteJson.EndTime
		runtimeCheck.RuntimePeriod = req.ExecuteJson.RuntimePeriod
	} else {
		runtimeCheck.StartTime = ""
		runtimeCheck.CyclePlanningType = 0
		runtimeCheck.CyclePlanningValue = 0
		runtimeCheck.CyclePlanningHour = ""
		runtimeCheck.EndTime = ""
		tempList := make([]string, 0)
		runtimeEnum := enums.TaskTaskEnum.AllRuntimePeriodEnum()
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod1])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod2])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod3])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod4])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod5])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod6])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod7])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod8])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod9])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod10])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod11])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod12])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod13])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod14])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod15])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod16])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod17])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod18])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod19])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod20])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod21])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod22])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod23])
		tempList = append(tempList, runtimeEnum[enums.TaskRuntimePeriod24])
		runtimeCheck.RuntimePeriod = tempList
	}
	if err := runtimeCheck.CheckRuntime(req.ExecuteType); err != nil {
		server.RespFail(c, 4000, err.Error())
		return
	}

	// 调用聚合层
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.ApiSave(ctx, &req, &res, uid.(int), runtimeCheck); err != nil {
		log.Error("TaskTaskSave application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// TaskProgress 任务进度接口
func TaskProgress(c *gin.Context) {
	var (
		req typespec.TaskProgressReq
		res typespec.TaskProgressResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskProgress parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.GetTaskProgress(ctx, &req, &res); err != nil {
		log.Error("TaskTaskSave application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// ApiVulList 漏洞测试列表及筛选
func ApiVulList(c *gin.Context) {
	var (
		req  typespec.ApiVulListReq
		resp typespec.ApiVulListResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.ApiVulList(ctx, &req, &resp); err != nil {
		log.Error("VulList application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AddTarget 动态添加目标
func AddTarget(c *gin.Context) {
	var (
		req  typespec.AddTargetReq
		resp typespec.AddTargetResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AddTarget parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.AddTarget(ctx, &req); err != nil {
		log.Error("AddTarget application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AddAttackFace 动态添加攻击面
func AddAttackFace(c *gin.Context) {
	var (
		req  typespec.AddAttackFaceReq
		resp typespec.AddAttackFaceResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AddAttackFace parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.AddAttackFace(ctx, &req); err != nil {
		log.Error("AddAttackFace application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// AddVul 动态添加漏洞
func AddVul(c *gin.Context) {
	var (
		req  typespec.AddVulReq
		resp typespec.AddVulResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("AddVul parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	if len(req.Pocname) == 0 { //没有脚本的直接返回
		server.RespFail(c, 4000, "该漏洞没有脚本,请尝试其他漏洞!")
		return
	}
	// 处理目标
	if okHttp := strings.HasPrefix(req.RootUrl, "http"); !okHttp {
		req.RootUrl = "http://" + req.RootUrl
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.AddVul(ctx, &req); err != nil {
		log.Error("AddVul application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskTargetMap 任务路径图
func TaskTargetMap(c *gin.Context) {
	var (
		req  typespec.TaskTargetMapReq
		resp typespec.TaskTargetMapRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTargetMap parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTargetMap(ctx, &req, &resp); err != nil {
		log.Error("TaskTargetMap application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskTargetMapNodeDetail 任务路径图 - 节点详情
func TaskTargetMapNodeDetail(c *gin.Context) {
	var (
		req  typespec.TaskTargetMapNodeDetailReq
		resp typespec.TaskTargetMapNodeDetailRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTargetMapNodeDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTargetMapNodeDetail(ctx, &req, &resp); err != nil {
		log.Error("TaskTargetMapNodeDetail application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskThreeExport 三方数据导出
func TaskThreeExport(c *gin.Context) {
	var (
		req  typespec.TaskThreeExportReq
		resp typespec.TaskThreeExportRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskThreeExport parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskThreeExport(ctx, &req, &resp); err != nil {
		log.Error("TaskThreeExport application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskTopologyMap 拓扑图
func TaskTopologyMap(c *gin.Context) {
	var (
		req  typespec.TaskTopologyMapReq
		resp typespec.TaskTopologyMapRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTopologyMap parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTopologyMap(ctx, &req, &resp); err != nil {
		log.Error("TaskTopologyMap application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// TaskTopologyMapNodeDetail 拓扑图节点详情
func TaskTopologyMapNodeDetail(c *gin.Context) {
	var (
		req  typespec.TaskTopologyMapNodeDetailReq
		resp typespec.TaskTopologyMapNodeDetailRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskTopologyMapNodeDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.TaskTopologyMapNodeDetail(ctx, &req, &resp); err != nil {
		log.Error("TaskTopologyMapNodeDetail application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

func TaskAllTaskVulByPage(c *gin.Context) {
	var (
		req  typespec.TaskAllVulByPageReq
		resp typespec.TaskAllVulByPageRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("TaskAllTaskVulByPage parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.AllTaskVulByPage(ctx, &req, &resp); err != nil {
		log.Error("TaskAllTaskVulByPage application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}

// CheckVul 漏洞检测
func CheckVul(c *gin.Context) {
	var (
		req  typespec.VulInfoReq
		resp typespec.CheckVulResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("VulInfo parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.TaskCheckTask
	if err := app.CheckVulInfo(ctx, &req, &resp); err != nil {
		log.Error("VulInfo application error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, resp)
}
