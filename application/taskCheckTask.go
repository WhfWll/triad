package application

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/data"
	aesEncryption "smart/tools/encryption"
	"smart/tools/enums"
	"smart/tools/network"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yaklang/yaklang/common/utils/lowhttp/poc"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
)

// 任务中心 - 创建任务

type TaskCheckTask struct {
}

// Enum 任务枚举
func (a *TaskCheckTask) Enum(ctx context.Context, res *typespec.TaskTaskEnumRes) error {
	var srv services.TaskTask

	res.Status = srv.StatusEnum()
	res.RiskLevel = srv.RiskLevelEnum()
	res.ExecuteType = srv.ExecuteTypeEnum()
	res.RuntimePeriod = srv.RuntimePeriodEnum()
	res.RuntimePeriod = srv.RuntimePeriodEnum()
	res.Weight = enums.TaskTaskEnum.GetTaskWeightEnumArray()
	res.TestIntensity = enums.TaskTaskEnum.GetTaskTestIntensityEnumArray()
	// 周期计划执行时的枚举
	res.CyclePlanningType = srv.CyclePlanningTypeEnum()
	res.CyclePlanningWeekValue = srv.CyclePlanningTypeWeekValueEnum()
	res.CyclePlanningMonthValue = srv.CyclePlanningTypeMonthValueEnum()
	// 网站登陆凭证
	res.WebLoginType = srv.WebLoginTypeEnum()
	//横向移动

	//攻击面
	var attackFaceType enums.AttackFace
	res.AttackFaceType = attackFaceType.GetAttackFaceTypeEnumArray()
	// 代理协议
	res.ProxyProto = srv.ProxyMode()
	return nil
}

// TaskTaskWebLoginCheck 网站登陆凭证校验
func (a *TaskCheckTask) TaskTaskWebLoginCheck(ctx context.Context, req *typespec.TaskTaskWebLoginCheckReq, res *typespec.TaskTaskWebLoginCheckRes) error {
	// 分析目标
	var analysisTarget data.TaskCheckTaskAnalysisTarget
	analysisTarget.AnalysisTarget(req.Target, "")
	errorTargetList := analysisTarget.ErrorTargetList
	if len(errorTargetList) > 0 {
		return errors.New(strings.Join(errorTargetList, ","))
	}
	targetList := analysisTarget.TargetList

	//登录地址是否在任务测试目标范围内
	existList := make([]int, 0)
	//当登录检测地址没有携带协议时，默认补全协议为http
	fullLoginUrl := utils.UrlHandleLogic.GetFullUrl(req.Target)
	_, loginHostname, loginPort, _, _, err := utils.UrlHandleLogic.ParseUrl(req.Target)
	if err != nil {
		return err
	}
	for _, oneTarget := range targetList {
		_, checkTargetHostname, checkTargetPort, _, _, err := utils.UrlHandleLogic.ParseUrl(oneTarget)
		if err != nil {
			return err
		}
		if checkTargetPort == "" {
			if loginHostname == checkTargetHostname {
				existList = append(existList, 1)
			}
		} else {
			if loginHostname == checkTargetHostname && loginPort == checkTargetPort {
				existList = append(existList, 1)
			}
		}
	}
	if len(existList) == 0 {
		return errors.New("登录地址必须在测试范围内")
	}

	//判断登录凭证是否有效
	resultCode, err := utils.UrlHandleLogic.VerifyStatus(fullLoginUrl, req.VerifyType, req.VerifyValue)
	if err != nil {
		return err
	}

	var websiteLogin enums.WebsiteLoginConfig
	res.Status = websiteLogin.VerifyStatusZh(resultCode)
	res.StatusCode = resultCode

	return nil
}

// Save 创建任务
func (a *TaskCheckTask) Save(ctx context.Context, req *typespec.TaskSaveReq, resp *typespec.TaskSaveResp) error {
	//添加任务创建测试范围校验的逻辑
	var (
		userGroupSrv   services.UserGroupUserList
		analysisTarget data.TaskCheckTaskAnalysisTarget
		scannerNode    services.Node
	)

	if ctx.Value("uid") == nil {
		return errors.New("用户未登录")
	}
	uid := ctx.Value("uid").(int)
	userIdList := []int{uid}
	userGroupList, err := userGroupSrv.GetUserGroupsByIds(ctx, userIdList)
	var testRangeLimit string
	for _, group := range userGroupList {
		if group.IsRangeOpen == enums.IsRangeOpenNo {
			continue
		}
		testRangeLimit += group.Range + "\n"
	}
	//analysisTarget.AnalysisTarget(testRangeLimit)
	//测试范围限制校验开关
	var mapSetServices services.MapSet
	testScopeSwitch, _ := mapSetServices.GetMapValue(ctx, enums.TestScopeSwitchMapSetObjKey)
	var rangeList []string

	// 中电邮储项目添加排除目标
	analysisTarget.AnalysisTarget(testRangeLimit, "")
	rangeList = analysisTarget.TargetList
	analysisTarget.TargetList = make([]string, 0)

	/******** 数据是否合法 ***********/
	// 校验场景是否存在
	var srvTemplate services.SceneTaskTemplate
	template, err := srvTemplate.GetByTemplateId(ctx, req.TaskTemplateId, enums.TaskTemplateStatusSuccess)
	if err != nil {
		return err
	}
	if template.ID == 0 {
		return errors.New("任务场景不存在")
	}
	// 判断config中的请求参数是否合规，并把需要默认值的参数赋值
	configStruct, err := data.TaskCheckTaskConfig.VerifyConfig(req.Config)
	if err != nil {
		return err
	}
	// 空5项目添加
	if req.Pid != 0 {
		configStruct.VulExploit = true
	}
	/******** 数据是否合法 end ***********/

	/******** 渗透模式校验 ***********/
	enableDistribute := scannerNode.SystemNodeIsAvailable(ctx)
	if enableDistribute && req.Config.Mode.Mode == enums.TaskConfigurationModeTarget {
		if len(req.Config.Mode.DistributeNodeId) == 0 {
			return errors.New("请选择引擎节点")
		}
	}
	//// 未开启 则不允许使用定向渗透
	//if remoteRes.Data.Status == 0 && req.Config.Mode.Mode == enums.TaskConfigurationModeTarget {
	//	return errors.New("请前往系统管理启用分布式引擎")
	//}

	/******* 渗透模式校验 end *******/

	/******** 非标准数据处理 ***********/
	taskName := req.TaskName
	taskName = strings.ReplaceAll(taskName, " ", "")
	taskName = strings.ReplaceAll(taskName, "\r", "")
	taskName = strings.ReplaceAll(taskName, "\n", "")

	// 分析目标
	//fmt.Println("11111111111111111111", testScopeSwitch)
	analysisTarget.SkipUrlValidation = testScopeSwitch == enums.TestScopeSwitchClose
	analysisTarget.AnalysisTarget(req.Target, req.ExcludeTarget)
	if testScopeSwitch != enums.TestScopeSwitchClose {
		errorTargetList := analysisTarget.ErrorTargetList
		if len(errorTargetList) > 0 {
			return errors.New(strings.Join(errorTargetList, ","))
		}
	}
	targetList := analysisTarget.TargetList
	//测试目标黑白名单过滤
	var mapSet services.MapSet
	targetList, err = mapSet.TargetIpWhiteBlack(ctx, targetList)
	if err != nil {
		return err
	}
	//一个任务的目标总数不能大于3200
	var taskControl map[string]int
	if err := config.Load("task_control", &taskControl); err != nil {
		log.Error("LogicTaskExec load config err", err)
	}
	var maxScanTargetNumber int
	if taskControl["max_scan_target_number"] > 0 {
		maxScanTargetNumber = taskControl["max_scan_target_number"]
	} else {
		maxScanTargetNumber = 3200
	}
	if len(targetList) > maxScanTargetNumber || len(targetList) == 0 {
		return errors.New("检测目标超过最大个数" + strconv.Itoa(maxScanTargetNumber) + "或开启了测试目标黑白名单检测目标为0")
	}
	/******** 非标准数据处理end ***********/

	//测试范围限制校验
	if len(rangeList) != 0 {
		metRequireList := make([]string, 0)
		for _, item := range rangeList {
			for _, target := range targetList {
				if item == target {
					metRequireList = append(metRequireList, target)
				}
			}
		}
		if len(metRequireList) != len(targetList) {
			return errors.New("检测目标中包含了用户组限制范围以外的目标")
		}
	}

	// 强制开启存活探测
	req.Config.AliveProbeConfig.IsOpen = true

	// 计算任务什么时候运行
	var taskExecNextTime data.TaskExecNextTime
	taskExecNextTime.StartTime, _ = time.ParseInLocation(utils.DateTime, req.ExecuteJson.StartTime, time.Local)
	taskExecNextTime.CyclePlanningType = req.ExecuteJson.CyclePlanningType
	taskExecNextTime.CyclePlanningValue = req.ExecuteJson.CyclePlanningValue
	taskExecNextTime.CyclePlanningHour = req.ExecuteJson.CyclePlanningHour
	taskExecNextTime.EndTime, _ = time.ParseInLocation(utils.DateTime, req.ExecuteJson.EndTime, time.Local)
	nextRuntime, err := taskExecNextTime.Compute(req.ExecuteType)
	if err != nil {
		return err
	}
	// 执行方式对应的参数
	executeJsonByte, _ := json.Marshal(req.ExecuteJson)

	// 创建
	var srvTask services.TaskTask
	if srvTask.CheckRequestIsRepeat(ctx, req.UserId, req.TaskTemplateId, taskName, req.Target) {
		return errors.New("请暂缓提交重复的请求")
	}
	resp.TaskId, err = srvTask.Save(ctx, req.UserId, req.TaskTemplateId, taskName, targetList, req.ExecuteType, string(executeJsonByte), configStruct, nextRuntime, req.Weight, req.Pid)
	if err != nil {
		return err
	}

	return nil
}

// SendTaskInfoToAuditLog 发送任务创建信息到审计日志
func (a *TaskCheckTask) SendTaskInfoToAuditLog(ctx context.Context, req *typespec.TaskSaveReq, err error, ip string) error {
	/******** 非标准数据处理 ***********/
	taskName := req.TaskName
	taskName = strings.ReplaceAll(taskName, " ", "")
	taskName = strings.ReplaceAll(taskName, "\r", "")
	taskName = strings.ReplaceAll(taskName, "\n", "")

	//获取用户
	var userService services.User
	user, _ := userService.GetUserDetail(ctx, req.UserId)

	var logAuditService services.LogAudit
	content := taskName + "任务创建成功"
	if err != nil {
		content = taskName + "任务创建失败"
	}
	return logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeOperate, content, user.Username, ip)
}

// SendTaskInfoToMessage 发送任务创建信息到消息表
func (a *TaskCheckTask) SendTaskInfoToMessage(ctx context.Context, req *typespec.TaskSaveReq, err error) error {
	/******** 非标准数据处理 ***********/
	taskName := req.TaskName
	taskName = strings.ReplaceAll(taskName, " ", "")
	taskName = strings.ReplaceAll(taskName, "\r", "")
	taskName = strings.ReplaceAll(taskName, "\n", "")

	messageType := enums.MessageTypeNotice
	content := taskName + "任务创建成功"
	if err != nil {
		messageType = enums.MessageTypeError
		content = taskName + "任务创建失败"
	}
	var messageService services.SystemMessage
	return messageService.SystemMessageAdd(ctx, content, messageType, req.UserId, enums.MessageStatusUnread)
}

// TaskTaskList 任务列表
func (a *TaskCheckTask) TaskTaskList(ctx context.Context, req *typespec.TaskListReq, res *typespec.TaskListRes) error {
	var srv services.TaskTask
	//添加普通用户只能获取自身所属任务的逻辑
	var userModel services.User
	if ctx.Value("uid") == nil {
		return errors.New("用户未登录")
	}
	uid := ctx.Value("uid").(int)
	user, err := userModel.GetUserForId(ctx, uid)
	userIdList := make([]int, 0)
	if user.Type == enums.UserRoleAuditor {
		return errors.New("审计员只能进行审计日志查看")
	}
	if user.Type == enums.UserRoleOrdinary {
		userIdList = append(userIdList, uid)
	}
	list, total, err := srv.List(ctx, req.Page, req.Size, req.TaskName, req.RiskLevel, req.StartTime, req.EndTime, userIdList)
	if err != nil {
		return err
	}
	res.Total = total
	taskId := make([]int, 0, len(list))
	for _, item := range list {
		taskId = append(taskId, item.ID)

		// 目标的风险等级
		targetRiskList := make([]int, 4)
		targetRiskList[0] = item.HigeNum
		targetRiskList[1] = item.MiddleNum
		targetRiskList[2] = item.LowNum
		targetRiskList[3] = item.SafeNum

		res.List = append(res.List, typespec.TaskListItemRes{
			Id:              item.ID,
			TaskName:        item.TaskName,
			ExecuteType:     item.ExecuteType,
			ExecuteTypeName: enums.TaskTaskEnum.ExecTypeEnum(item.ExecuteType),
			RiskLevel:       item.RiskLevel,
			RiskLevelName:   enums.TaskTaskEnum.RiskEnum(item.RiskLevel),
			Status:          item.Status,
			StatusName:      enums.TaskTaskEnum.StatusEnum(item.Status),
			TargetRisk:      targetRiskList,
			CreateTime:      item.CreateTime.Format(utils.DateTime),
			UpdateTime:      item.UpdateTime.Format(utils.DateTime),
		})
	}

	/************** A6项目需要进度条 begin **************/
	var taskTaskInfoSrv services.TaskTaskInfo
	taskOverviewMap := taskTaskInfoSrv.GetOverviewByTaskIds(ctx, taskId)
	for k, item := range res.List {
		switch item.Status {
		case enums.TaskStatusBegin:
			res.List[k].Progress = 0
		case enums.TaskStatusFinish:
			res.List[k].Progress = 100
		default:
			if taskOverview, ok := taskOverviewMap[item.Id]; ok {
				if taskOverview.Progress == nil {
					res.List[k].Progress = 1
				} else {
					res.List[k].Progress = taskOverview.Progress.Value
				}
			} else {
				res.List[k].Progress = 1
			}
		}
	}
	/************** A6项目需要进度条 end **************/

	return nil
}

// TaskTaskDel 任务删除
func (a *TaskCheckTask) TaskTaskDel(ctx context.Context, req *typespec.TaskDelReq) error {
	ids := strings.Split(req.TaskIds, ",")
	taskIds := make([]int, 0)
	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		taskIds = append(taskIds, idInt)
	}

	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	// 删除任务相关所有数据
	// 任务表 task_task
	var taskSrv services.TaskTask
	err := taskSrv.DelTaskByIds(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 任务详情表 task_task_info
	err = taskSrv.DelTaskInfoByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 目标表 task_target
	var taskTargetSrv services.TaskTarget
	err = taskTargetSrv.DelTaskInfoByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 检测结果表 task_result
	var taskResultSrv services.TaskResult
	err = taskResultSrv.DelTaskInfoByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 通用记录数据表 task_target_result
	var taskTargetResultSrv services.TaskTargetResult
	err = taskTargetResultSrv.DelTaskInfoByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 通用结果数据表 task_task_result
	var taskTaskResultSrv services.TaskTaskResult
	err = taskTaskResultSrv.DelTaskInfoByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 任务日志表 task_log
	var taskLogSrv services.TaskLog
	err = taskLogSrv.DelTaskInfoByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 任务日志详情表 task_log_info
	var taskLogInfoSrv services.TaskLogInfo
	err = taskLogInfoSrv.DelTaskInfoByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 漏洞表 task_vul
	var taskVulSrv services.TaskVul
	err = taskVulSrv.DelTaskInfoByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 事务提交
	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// TaskTaskCopy 任务复制
func (a *TaskCheckTask) TaskTaskCopy(ctx context.Context, req *typespec.TaskCopyReq, res *typespec.TaskCopyRes) error {
	// 拷贝的逻辑与创建任务一致，仅任务名称前需要增加copy_标识

	// 获取待拷贝的任务
	var taskSrv services.TaskTask
	task, taskInfo, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	if task.ID == 0 {
		return errors.New("未知的任务")
	}
	// 执行方式对应的json值
	var execJson typespec.ExecuteJson
	err = json.Unmarshal([]byte(taskInfo.ExecuteJSON), &execJson)
	if err != nil {
		return err
	}
	// 配置
	var conf enums.ConfigJson
	err = json.Unmarshal([]byte(taskInfo.TaskTemplateJSON), &conf)
	if err != nil {
		return err
	}

	// 获取目标
	var targetSrv services.TaskTarget
	targets, _ := targetSrv.GetTargetByTaskId(ctx, req.TaskId)
	targetSlice := make([]string, 0)
	for _, item := range targets {
		targetSlice = append(targetSlice, item.TargetURL)
	}

	// 组织返回数据
	res.TaskTemplateId = task.TaskTemplateID
	res.Target = strings.Join(targetSlice, ",")
	res.TaskName = "copy_" + task.TaskName
	res.ExecuteType = task.ExecuteType
	res.ExecuteJson = execJson
	res.Config = conf
	res.Weight = task.Weight
	res.UserId = task.UserID
	return nil
}

// TaskTaskChangeState 任务修改状态
func (a *TaskCheckTask) TaskTaskChangeState(ctx context.Context, req *typespec.TaskChangeStateReq) error {
	switch req.Operate {
	case enums.TaskOperatePause: // 暂停
		return a.taskTaskChangeStatePurse(ctx, req.TaskId)
	case enums.TaskOperateResume: // 恢复
		return a.taskTaskChangeStateResume(ctx, req.TaskId)
	case enums.TaskOperateStop: // 停止
		return a.taskTaskChangeStateFinish(ctx, req.TaskId)
	}
	return nil
}

// taskTaskChangeStatePurse 任务修改状态 - 暂停
func (a *TaskCheckTask) taskTaskChangeStatePurse(ctx context.Context, taskId int) error {

	var taskSrv services.TaskTask
	var taskInfoSrv services.TaskTaskInfo
	var taskLogSrv services.TaskLog

	task, _, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, taskId)
	if err != nil {
		return err
	}
	if task.ID == 0 {
		return errors.New("未知的任务")
	}
	if task.Status != enums.TaskStatusRunning {
		return errors.New("该任务非运行中 无需暂停")
	}

	// 获取任务下的所有目标
	// 1 给每个目标发送操作信号
	// 2 修改每个目标对应的状态
	// 3 修改任务对应的状态
	var targetSev services.TaskTarget
	targetList, _ := targetSev.GetTargetByTaskId(ctx, taskId)
	// 任务状态
	for _, target := range targetList {
		// 任务状态 仅运行中的任务允许暂停
		if target.Status == enums.TargetStatusRunning {
			// 通知决策引擎进行暂停操作
			//_ = targetSev.PurseTarget(target.ID)
			// 修改目标状态 为暂停操作
			_ = targetSev.UpdateTargetStateById(ctx, target.ID, enums.TargetStatusPausing)
			// 日志状态
			_ = taskLogSrv.UpdateTaskLogStateByTargetId(ctx, target.ID, enums.TaskStatusPausing)
		}
	}

	// 修改任务状态 为暂停操作
	_ = taskSrv.UpdateTaskStateById(ctx, taskId, enums.TaskStatusPausing)
	_ = taskInfoSrv.UpdateTaskStateByTaskId(ctx, taskId, enums.TaskStatusPausing)
	return nil
}

// taskTaskChangeStateResume 任务修改状态 - 恢复
func (a *TaskCheckTask) taskTaskChangeStateResume(ctx context.Context, taskId int) error {

	var taskSrv services.TaskTask
	var taskInfoSrv services.TaskTaskInfo
	var taskLogSrv services.TaskLog

	task, _, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, taskId)
	if err != nil {
		return err
	}
	if task.ID == 0 {
		return errors.New("未知的任务")
	}
	if task.Status != enums.TaskStatusPausing {
		return errors.New("该任务非暂停中 无需恢复")
	}

	// 获取任务下的所有目标
	// 1 给每个目标发送操作信号
	// 2 修改每个目标对应的状态
	// 3 修改任务对应的状态
	var targetSev services.TaskTarget
	targetList, _ := targetSev.GetTargetByTaskId(ctx, taskId)
	// 任务状态
	for _, target := range targetList {
		// 任务状态 仅暂停中的任务允许恢复
		if target.Status == enums.TargetStatusPausing {
			// 通知决策引擎进行恢复操作
			err = targetSev.ResumeTarget(target.ID)
			if err != nil {
				log.Error("任务修改状态 - 恢复 - 通知决策引擎失败：" + err.Error())
			}
			// 修改目标状态 为恢复操作
			_ = targetSev.UpdateTargetStateById(ctx, target.ID, enums.TargetStatusRunning)
			// 日志状态
			_ = taskLogSrv.UpdateTaskLogStateByTargetId(ctx, target.ID, enums.TargetStatusRunning)
		}
	}

	// 修改任务状态 为恢复操作
	_ = taskSrv.UpdateTaskStateById(ctx, taskId, enums.TaskStatusRunning)
	_ = taskInfoSrv.UpdateTaskStateByTaskId(ctx, taskId, enums.TaskStatusRunning)
	return nil
}

// taskTaskChangeStateFinish 任务修改状态 - 结束
func (a *TaskCheckTask) taskTaskChangeStateFinish(ctx context.Context, taskId int) error {

	var taskSrv services.TaskTask
	var taskInfoSrv services.TaskTaskInfo
	var taskLogSrv services.TaskLog

	task, _, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, taskId)
	if err != nil {
		return err
	}
	if task.ID == 0 {
		return errors.New("未知的任务")
	}
	if task.Status != enums.TaskStatusBegin && task.Status != enums.TaskStatusRunning && task.Status != enums.TaskStatusPausing {
		return errors.New("该任务非待开始或运行中或暂停中 无需结束")
	}

	// 获取任务下的所有目标
	// 1 给每个目标发送操作信号
	// 2 修改每个目标对应的状态
	// 3 修改任务对应的状态
	go func() {
		//改成协程，防止阻塞
		var targetSev services.TaskTarget
		targetList, _ := targetSev.GetTargetByTaskId(ctx, taskId)
		// 任务状态
		for _, target := range targetList {
			// 任务状态 待开始或运行中或暂停中的任务允许停止
			if target.Status == enums.TargetStatusToBegin || target.Status == enums.TargetStatusRunning || target.Status == enums.TargetStatusPausing {
				// 通知决策引擎进行停止操作
				//_ = targetSev.StopTarget(target.ID)
				// 取消本地扫描任务
				services.CancelTargetScan(target.ID)
				// 修改目标状态 为停止操作
				_ = targetSev.UpdateTargetStateById(ctx, target.ID, enums.TargetStatusFinish)
				// 日志状态
				_ = taskLogSrv.UpdateTaskLogStateByTargetId(ctx, target.ID, enums.TargetStatusFinish)
			}
		}
	}()

	// 修改任务状态 为停止操作
	_ = taskSrv.UpdateTaskStateById(ctx, taskId, enums.TaskStatusFinish)
	_ = taskInfoSrv.UpdateTaskStateByTaskId(ctx, taskId, enums.TaskStatusFinish)

	// 提前完成任务同步更新资产信息
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("UpdateAssetInfo panic: %v", r)
			}
		}()
		log.Printf("Active FinishTask Begin UpdateAssetInfo taskID:%d", taskId)
		assetSrv := services.Asset{}
		assetSrv.UpdateAssetInfo(context.Background(), taskId)
	}()
	return nil
}

// OverView 任务概览
func (a *TaskCheckTask) GetState(ctx context.Context, req *typespec.GetStateReq, resp *typespec.GetStateResp) error {
	var taskSrv services.TaskTask
	taskRes, err := taskSrv.GetTaskByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	if taskRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	resp.Status = taskRes.Status
	resp.StatusName = enums.TaskTaskEnum.StatusEnum(taskRes.Status)
	return nil
}

// OverView 任务概览
func (a *TaskCheckTask) OverView(ctx context.Context, req *typespec.OverViewReq, resp *typespec.OverViewResp) error {
	//查询task和taskinfo数据
	var taskSrv services.TaskTask
	taskRes, taskinfoRes, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	if taskRes.ID == 0 || taskinfoRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	//查询场景信息
	var tmplateSrv services.SceneTaskTemplate
	tempRes, err := tmplateSrv.GetTaskTemplateById(ctx, taskRes.TaskTemplateID)
	if err != nil {
		return err
	}
	//组装返回结果
	resp.Id = taskRes.ID
	resp.TaskName = taskRes.TaskName
	resp.ExecuteType = taskRes.ExecuteType
	resp.ExecuteTypeName = enums.TaskTaskEnum.ExecTypeEnum(taskRes.ExecuteType)
	resp.RuntimePeriod = ""
	resp.StartTime = ""
	resp.CyclePlanningType = 0
	resp.CyclePlanningValue = 0
	resp.CyclePlanningHour = ""
	resp.EndTime = ""
	resp.CyclePlanningName = ""
	resp.EndTimeName = ""
	if len(taskinfoRes.ExecuteJSON) > 0 {
		var tmp services.TaskRunConf
		err = json.Unmarshal([]byte(taskinfoRes.ExecuteJSON), &tmp)
		if err != nil {
			return err
		}
		resp.RuntimePeriod = strings.Join(tmp.RuntimePeriod, ",")
		resp.StartTime = tmp.StartTime
		resp.CyclePlanningType = tmp.CyclePlanningType
		resp.CyclePlanningValue = tmp.CyclePlanningValue
		resp.CyclePlanningHour = tmp.CyclePlanningHour
		resp.EndTime = tmp.EndTime
		if taskRes.ExecuteType == enums.TaskExecTypeTiming {
			resp.CyclePlanningName = tmp.StartTime
		} else if taskRes.ExecuteType == enums.TaskExecTypeCycle {
			resp.CyclePlanningName += enums.TaskTaskEnum.ExecCycleTypeEnum(tmp.CyclePlanningType) + " "
			if tmp.CyclePlanningType == enums.TaskExecTypeCycleTypeWeek {
				resp.CyclePlanningName += enums.TaskTaskEnum.ExecCycleTypeWeekValueEnum(tmp.CyclePlanningValue) + " " + tmp.CyclePlanningHour
			} else if tmp.CyclePlanningType == enums.TaskExecTypeCycleTypeMonth {
				resp.CyclePlanningName += enums.TaskTaskEnum.ExecCycleTypeMonthValueEnum(tmp.CyclePlanningValue) + " " + tmp.CyclePlanningHour
			}
		}
		resp.EndTimeName = tmp.EndTime
	}
	resp.RiskLevel = taskRes.RiskLevel
	resp.RiskLevelName = enums.TaskTaskEnum.RiskEnum(taskRes.RiskLevel)
	resp.TaskTemplateId = taskRes.TaskTemplateID
	resp.TaskTemplateName = ""
	if tempRes.ID != 0 {
		resp.TaskTemplateName = tempRes.TemplateName
	}
	resp.Status = taskRes.Status
	resp.StatusName = enums.TaskTaskEnum.StatusEnum(taskRes.Status)
	resp.CreateTime = taskRes.CreateTime.Format(enums.TimeLayout)
	resp.UpdateTime = ""
	if taskRes.Status == enums.TaskStatusFinish { //已完成状态才有结束时间
		resp.UpdateTime = taskRes.UpdateTime.Format(enums.TimeLayout)
	}
	resp.TargetRisk = [4]int{}
	resp.TargetNum = make([]int, 0)
	resp.VulRisk = make([]int, 0)
	resp.VulExploitImpact = make([]services.OverViewItems, 0)
	resp.EvidenceStat = make([]services.OverViewItems, 0)
	resp.Port = make([]services.OverViewItems, 0)
	resp.Service = make([]services.OverViewItems, 0)
	resp.Component = make([]services.OverViewItems, 0)
	resp.OpSys = make([]services.OverViewItems, 0)
	resp.SubDomain = make([]services.OverViewItems, 0)
	resp.UrlTags = make([]services.OverViewItems, 0)
	if len(taskinfoRes.Overview) > 0 {
		var tmp services.OverView
		err = json.Unmarshal([]byte(taskinfoRes.Overview), &tmp)
		if err != nil {
			return err
		}
		resp.TargetRisk = tmp.TargetRisk
		resp.TargetNum = tmp.TargetNum
		resp.VulTotal = tmp.VulTotal
		resp.VulRisk = tmp.VulRisk
		resp.VulExploitImpact = tmp.VulExploitImpact
		resp.EvidenceStat = tmp.EvidenceStat
		resp.Port = tmp.Port
		resp.Service = tmp.Service
		resp.Component = tmp.Component
		resp.OpSys = tmp.OpSys
		resp.SubDomain = tmp.SubDomain
		resp.UrlTags = tmp.UrlTags
	}
	return nil
}

// TaskConfigInfo 渗透任务 - 任务配置信息
func (a *TaskCheckTask) TaskConfigInfo(ctx context.Context, req *typespec.TaskConfigInfoReq, res *typespec.TaskConfigInfoRes) error {
	var taskSrv services.TaskTask
	task, taskInfo, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	if taskInfo.ID == 0 {
		return errors.New("任务不存在")
	}

	var conf enums.ConfigJson
	err = json.Unmarshal([]byte(taskInfo.TaskTemplateJSON), &conf)
	if err != nil {
		return err
	}

	// 优先使用主表 task_task 的权重，确保数据一致性
	res.Priority = task.Weight
	res.PriorityZh = enums.TaskTaskEnum.GetTaskWeightEnum(task.Weight)

	//处理config中的中文返回
	var aliveProbeConfig *enums.AliveProbeConfig
	conf.AliveProbeConfig.AliveProbeTypeZh = aliveProbeConfig.AliveProbeTypeEnum(conf.AliveProbeConfig.AliveProbeType)
	//端口扫描
	var portScanConfig *enums.PortScanConfig
	conf.PortScanConfig.PortScanTypeZh = portScanConfig.PortScanTypeEnum(conf.PortScanConfig.PortScanType)
	conf.PortScanConfig.TCPScanTypeZh = portScanConfig.TcpScanTypeEnum(conf.PortScanConfig.TCPScanType)
	conf.PortScanConfig.TimeoutZh = portScanConfig.PortScanTimeoutEnum(conf.PortScanConfig.Timeout)
	conf.PortScanConfig.ConcurrentZh = portScanConfig.PortScanConcurrentEnum(conf.PortScanConfig.Concurrent)
	//动态爬虫
	var webCrawlerConfig *enums.WebCrawlerConfig
	conf.WebCrawlerConfig.MaxDepthZh = webCrawlerConfig.CrawlerDeep(conf.WebCrawlerConfig.MaxDepth)
	conf.WebCrawlerConfig.MaxUrlZh = webCrawlerConfig.CrawlerMaxConnect(conf.WebCrawlerConfig.MaxUrl)
	conf.WebCrawlerConfig.ScanRangeZh = webCrawlerConfig.CrawlerRange(conf.WebCrawlerConfig.ScanRange)
	conf.WebCrawlerConfig.TimeoutZh = webCrawlerConfig.SingleTimeout(conf.WebCrawlerConfig.Timeout)
	conf.WebCrawlerConfig.FullTimeoutZh = webCrawlerConfig.FullTimeoutEnum(conf.WebCrawlerConfig.FullTimeout)
	conf.WebCrawlerConfig.ScanRepeatZh = webCrawlerConfig.CrawlerRepeat(conf.WebCrawlerConfig.ScanRepeat)
	//web路径爆破
	var webPathScanConfig *enums.WebPathScanConfig
	conf.WebPathScanConfig.GuessRateZh = webPathScanConfig.WebPathScanSpeed(conf.WebPathScanConfig.GuessRate)
	conf.WebPathScanConfig.GuessTimeoutZh = webPathScanConfig.WebPathScanTime(conf.WebPathScanConfig.GuessTimeout)
	//口令爆破
	var weakPassConfig enums.WeakPassConfig
	for _, serviceInt := range conf.WeakPassConfig.Services {
		conf.WeakPassConfig.ServicesZh = append(conf.WeakPassConfig.ServicesZh, enums.GetDictionaryService(serviceInt))
	}
	conf.WeakPassConfig.DictTypeZh = weakPassConfig.WeakPassDictType(conf.WeakPassConfig.DictType)
	conf.WeakPassConfig.GuessNumZh = weakPassConfig.WeakPassGuessNumber(conf.WeakPassConfig.GuessNum)
	conf.WeakPassConfig.GuessTimeoutZh = weakPassConfig.WeakPassGuessTime(conf.WeakPassConfig.GuessTimeout)
	conf.WeakPassConfig.GuessRateZh = weakPassConfig.WeakPassGuessRate(conf.WeakPassConfig.GuessRate)
	//站点登录凭证
	var websiteLoginConfig *enums.WebsiteLoginConfig
	loginConfigNodeList := make([]enums.LoginConfigNode, 0)
	for _, item := range conf.WebsiteLoginConfig.List {
		if item.Scheme == "" {
			item.Scheme = "http"
		}
		loginConfigNodeList = append(loginConfigNodeList, enums.LoginConfigNode{
			Target:         item.Target,
			VerifyType:     item.VerifyType,
			VerifyTypeZh:   websiteLoginConfig.VerifyStatusZh(item.VerifyStatus),
			VerifyStatus:   item.VerifyStatus,
			VerifyStatusZh: websiteLoginConfig.WebsiteLoginType(item.VerifyType),
			VerifyValue:    item.VerifyValue,
			Scheme:         item.Scheme,
		})
	}
	conf.WebsiteLoginConfig.List = loginConfigNodeList
	//渗透模式
	var mode enums.Mode
	conf.Mode.ModeZh = mode.GetModeEnum(conf.Mode.Mode)
	//处理字典部分
	ids := make([]int, 0)
	if conf.WeakPassConfig.CommonUserDict != 0 {
		ids = append(ids, conf.WeakPassConfig.CommonUserDict)
		// 填充通用用户字典中文名称
		var dictModel mysqls.Dictionary
		dictInfo, _ := dictModel.DictionaryRecord(ctx, conf.WeakPassConfig.CommonUserDict)
		if dictInfo.ID != 0 {
			conf.WeakPassConfig.CommonUserDictZh = dictInfo.Name
		}
	}
	if conf.WeakPassConfig.CommonPassDict != 0 {
		ids = append(ids, conf.WeakPassConfig.CommonPassDict)
		// 填充通用密码字典中文名称
		var dictModel mysqls.Dictionary
		dictInfo, _ := dictModel.DictionaryRecord(ctx, conf.WeakPassConfig.CommonPassDict)
		if dictInfo.ID != 0 {
			conf.WeakPassConfig.CommonPassDictZh = dictInfo.Name
		}
	}
	if len(conf.WebPathScanConfig.ScanDict) > 0 {
		ids = append(ids, conf.WebPathScanConfig.ScanDict...)
	}
	if conf.SubdomainCollectConfig.SubdomainDict != 0 {
		ids = append(ids, conf.SubdomainCollectConfig.SubdomainDict)
		// 填充子域名字典中文名称
		var dictModel mysqls.Dictionary
		dictInfo, _ := dictModel.DictionaryRecord(ctx, conf.SubdomainCollectConfig.SubdomainDict)
		if dictInfo.ID != 0 {
			conf.SubdomainCollectConfig.SubdomainDictZh = dictInfo.Name
		}
	}
	res.Config = conf
	return nil
}

// TargetChangeState 修改测试目标状态
func (a *TaskCheckTask) TargetChangeState(ctx context.Context, req *typespec.TargetChangeStateReq) error {
	//查询目标状态
	var targetSrv services.TaskTarget
	targetRes, err := targetSrv.GetTargetById(ctx, req.TargetId)
	if err != nil {
		return err
	}
	if targetRes.Status != enums.TargetStatusRunning {
		return errors.New("只有运行中状态的目标才能执行结束操作")
	}
	//请求决策引擎停止任务
	//err = targetSrv.StopTarget(req.TargetId)
	//if err != nil {
	//	return err
	//}
	// 取消本地扫描任务
	services.CancelTargetScan(req.TargetId)
	//修改目标状态和目标日志为结束
	err = targetSrv.UpdateTargetAndLogStateById(ctx, req.TargetId, enums.TargetStatusFinish, enums.TargetStatusFinish)
	if err != nil {
		return err
	}
	//更新任务表重新统计
	var taskSrv services.TaskTask
	err = taskSrv.UpdateIsStats(ctx, targetRes.TaskID, enums.TaskIsStatsYes)
	if err != nil {
		return err
	}
	return nil
}

// TargetList 测试目标列表及筛选
func (a *TaskCheckTask) TargetList(ctx context.Context, req *typespec.TargetListReq, resp *typespec.TargetListResp) error {
	var (
		targetSrv  services.TaskTarget
		taskVulSrv services.TaskVul
	)
	//查询target列表
	targetRes, total, err := targetSrv.TargetList(ctx, req.TaskId, req.RiskLevel, req.Search, req.Page, req.Size)
	if err != nil {
		return err
	}
	//查询开发端口
	var ids []int
	for i := 0; i < len(targetRes); i++ {
		ids = append(ids, targetRes[i].ID)
	}
	openRes, err := targetSrv.GetTargetOpenPort(ctx, ids)
	if err != nil {
		return err
	}
	//统计漏洞数量
	_, riskLevleMap, vulNumArrayMap, err := taskVulSrv.GetTargetStatsBytargetIds(ctx, ids)
	//组装返回结果
	resp.Total = total
	resp.List = make([]typespec.TargetListRespItems, 0)
	for i := 0; i < len(targetRes); i++ {
		var tmp typespec.TargetListRespItems
		tmp.Id = targetRes[i].ID
		tmp.TargetUrl = targetRes[i].TargetURL
		tmp.OpSys = targetRes[i].OpSys
		tmp.OpenPort = []string{}
		if v, ok := openRes[strconv.Itoa(targetRes[i].ID)]; ok {
			tmp.OpenPort = v
		}
		tmp.RiskLevel = enums.TargetRiskLowNoFound
		if v, ok := riskLevleMap[targetRes[i].ID]; ok {
			tmp.RiskLevel = v
		}
		tmp.RiskLevelName = enums.GetTargetRisk(tmp.RiskLevel)
		tmp.VulNum = []int{}
		if v, ok := vulNumArrayMap[targetRes[i].ID]; ok {
			tmp.VulNum = []int{v[1], v[2], v[3], v[4], v[5]}
		}
		tmp.Status = targetRes[i].Status
		tmp.StatusName = enums.GetTargetStatus(targetRes[i].Status)
		tmp.UseScore = strconv.FormatFloat(float64(targetRes[i].UseScore)/100, 'f', 2, 64)
		tmp.IsAlive = targetRes[i].IsAlive
		tmp.IsAliveName = enums.GetTargetIsAlive(targetRes[i].IsAlive)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// UpdateTargetUseScore 批量修改目标的利用评分和状态
func (a *TaskCheckTask) UpdateTargetUseScore(ctx context.Context, req *typespec.UpdateTargetUseScoreReq) error {
	var targetSrv services.TaskTarget
	return targetSrv.UpdateTargetUseScore(ctx, req.Data)
}

// TargetDel 测试目标删除
func (a *TaskCheckTask) TargetDel(ctx context.Context, req *typespec.TargetDelReq) error {
	targetIdArray := strings.Split(req.TargetIds, ",")
	var targetSrv services.TaskTarget
	//删除相关表
	err := targetSrv.DelTargetByIds(ctx, targetIdArray)
	if err != nil {
		return err
	}
	//更新任务表重新统计
	var taskSrv services.TaskTask
	err = taskSrv.UpdateIsStats(ctx, req.TaskId, enums.TaskIsStatsYes)
	if err != nil {
		return err
	}
	return nil
}

// TaskResultList 信息收集列表及筛选
func (a *TaskCheckTask) TaskResultList(ctx context.Context, req *typespec.TaskResultListReq, resp *typespec.TaskResultListResp) error {
	//查询数据
	var taskResultSrv services.TaskTaskResult
	taskResultRes, count, err := taskResultSrv.TaskResultList(ctx, req.ObjType, req.SubObjType, req.ObjId, req.Search, req.Page, req.Size)
	if err != nil {
		return err
	}
	//组装返回结果
	resp.Total = count
	resp.List = make([]map[string]interface{}, 0)
	for i := 0; i < len(taskResultRes); i++ {
		var tmp = make(map[string]interface{}, 0)
		tmp["id"] = taskResultRes[i].ID
		tmp["cdn"] = taskResultSrv.CheckCNDTaskResultBySubObjID(ctx, taskResultRes[i].ObjID, taskResultRes[i].Field1)
		json.Unmarshal([]byte(taskResultRes[i].JSONResult), &tmp)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// TaskResultUrlTree url树状图
func (a *TaskCheckTask) TaskResultUrlTree(ctx context.Context, req *typespec.TaskResultUrlTreeReq, resp *typespec.TaskResultUrlTreeResp) error {
	//查询数据
	var taskResultSrv services.TaskTaskResult
	taskResultRes, _, err := taskResultSrv.TaskResultList(ctx, 1, "1_3", req.ObjId, "", 0, 10000)
	if err != nil {
		return err
	}
	//组装返回结果
	paths := make([][]string, len(taskResultRes))
	for i, item := range taskResultRes {
		paths[i] = utils.ParsePath(item.Field1)
	}
	root := utils.BuildTree(paths)
	tree := utils.ConvertTreeToSlice(root)
	resp.TreeData = tree
	return nil
}

// TaskResultDel 信息收集删除
func (a *TaskCheckTask) TaskResultDel(ctx context.Context, req *typespec.TaskResultDelReq) error {
	idArray := strings.Split(req.TaskTaskResultIds, ",")
	//删除数据
	var taskResultSrv services.TaskTaskResult
	err := taskResultSrv.DelTaskResult(ctx, idArray)
	if err != nil {
		return err
	}
	//更新任务表重新统计
	var taskSrv services.TaskTask
	err = taskSrv.UpdateIsStats(ctx, req.TaskId, enums.TaskIsStatsYes)
	if err != nil {
		return err
	}
	return nil
}

// TaskResultDetail 信息收集查询详情
func (a *TaskCheckTask) TaskResultDetail(ctx context.Context, req *typespec.TaskResultDetailReq, resp *typespec.TaskResultDetailResp) error {
	var taskResultSrv services.TaskTaskResult
	taskTaskResult, err := taskResultSrv.TaskResultDetail(ctx, req.TaskTaskResultId)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(taskTaskResult.JSONResult), resp)
	if err != nil {
		return err
	}
	(*resp)["cdn"] = taskResultSrv.CheckCNDTaskResultBySubObjID(ctx, taskTaskResult.ObjID, taskTaskResult.Field1)
	return nil
}

// VulList 漏洞测试列表及筛选
func (a *TaskCheckTask) VulList(ctx context.Context, req *typespec.VulListReq, resp *typespec.VulListResp) error {
	// 查询任务数据 用于判断是否为二次验证任务（二次验证：验证与父级任务的漏洞是否修复与是否有新增漏洞）
	var taskSrv services.TaskTask
	task, err := taskSrv.GetTaskByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}

	if task.Pid == 0 {
		// 普通任务
		return a.VulListDefault(ctx, req, resp)
	} else {
		// 二次验证测试任务
		return a.VulListRepeatCheck(ctx, task.Pid, req, resp)
	}
}

// 任务详情 - 普通任务漏洞列表
func (a *TaskCheckTask) VulListDefault(ctx context.Context, req *typespec.VulListReq, resp *typespec.VulListResp) error {
	//查询数据
	var taskVulSrv services.TaskVul
	// 销售许可证临时添加
	key := []byte("9876787656785679")
	aesEcb := aesEncryption.AesEcb{}
	req.Search = strings.TrimSpace(req.Search)
	search := ""
	if req.Search != "" {
		search = req.Search
	}
	vulRes, total, err := taskVulSrv.VulList(ctx, req.TaskId, req.TargetId, req.Type, req.Risk, search, req.DataType, req.Page, req.Size, req.Status)
	if err != nil {
		return err
	}

	//组装返回结果
	resp.Total = total
	resp.List = make([]typespec.VulListRespItems, 0)
	for i := 0; i < len(vulRes); i++ {
		var tmp typespec.VulListRespItems
		tmp.Id = vulRes[i].ID
		tmp.DataType = vulRes[i].DataType

		// 处理 TargetUrl
		if utils.IsHexString(vulRes[i].TargetUrl) {
			targetUrlDecodeByte, _ := hex.DecodeString(vulRes[i].TargetUrl)
			tmp.TargetUrl = string(aesEcb.AesDecryptECB(targetUrlDecodeByte, key))
		} else {
			tmp.TargetUrl = vulRes[i].TargetUrl
		}

		// 处理 Name
		if utils.IsHexString(vulRes[i].Name) {
			nameDecodeByte, _ := hex.DecodeString(vulRes[i].Name)
			tmp.Name = string(aesEcb.AesDecryptECB(nameDecodeByte, key))
		} else {
			tmp.Name = vulRes[i].Name
		}

		tmp.Type = vulRes[i].Type
		tmp.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(vulRes[i].Type)
		tmp.Risk = vulRes[i].Risk
		tmp.RiskName = enums.ToolsVulnerabilityEnum.GetRiskEnum(vulRes[i].Risk)

		// 处理 Location
		if utils.IsHexString(vulRes[i].Location) {
			locationDecodeByte, _ := hex.DecodeString(vulRes[i].Location)
			tmp.Location = string(aesEcb.AesDecryptECB(locationDecodeByte, key))
		} else {
			tmp.Location = vulRes[i].Location
		}

		tmp.Status = vulRes[i].Status
		tmp.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(vulRes[i].Status)

		tmp.TestStatus = vulRes[i].TestStatus
		tmp.TestStatusName = enums.ToolsVulnerabilityEnum.GetVulTestStatusEnum(vulRes[i].TestStatus)
		tmp.IsSnapshot = enums.VulSnapshotNo
		if len(vulRes[i].Snapshot) > 0 {
			tmp.IsSnapshot = enums.VulSnapshotYes
		}

		//if req.Search != "" {
		//	if !strings.Contains(tmp.TargetUrl, req.Search) && !strings.Contains(tmp.Name, req.Search) && !strings.Contains(tmp.Location, req.Search) {
		//		continue
		//	}
		//}

		// todo 兜底千机
		if tmp.Name == "qianji" {
			var resp typespec.VulInfoResp
			a.VulInfo(ctx, &typespec.VulInfoReq{TaskVulId: tmp.Id}, &resp)
			tmp.Name = resp.Name
			tmp.TypeName = resp.TypeName
			tmp.Location = resp.VulAddress
			tmp.RiskName = resp.RiskName
			tmp.Risk = resp.Risk
		}
		resp.List = append(resp.List, tmp)
	}
	return nil
}

/*
* 任务详情 - 二次验证任务漏洞列表
* 获取父级任务的漏洞数据，与当前任务进行对比，已父级漏洞为准，除原有漏洞状态之外，增加2个状态 已修复 新增加
* 	@param pid 父级任务ID
 */
func (a *TaskCheckTask) VulListRepeatCheck(ctx context.Context, pid int, req *typespec.VulListReq, resp *typespec.VulListResp) error {
	// 查询父级漏洞数据
	var taskVulSrv services.TaskVul
	taskIdList := []int{pid, req.TaskId}
	vulParentRes, vulParentTotal, err := taskVulSrv.SecondScanVulList(ctx, taskIdList, req.TargetId, req.Type, req.Risk, req.Search, req.DataType, req.Page, req.Size, req.Status)
	if err != nil {
		return err
	}

	// 查询当前任务下的所有漏洞
	vulCurrentRes := taskVulSrv.GetsByTaskId(ctx, req.TaskId, enums.VulDataTypOne)
	if err != nil {
		return err
	}

	// 销售许可证临时添加
	key := []byte("9876787656785679")
	aesEcb := aesEncryption.AesEcb{}

	// 组合任务数据进行验证，唯一漏洞寻找逻辑 （target_id + pocname + name + location）
	vulCurrentMaps := make(map[string]mysqls.TaskVul)
	// 复制当前检测到的漏洞，用于下面逻辑中进行删除，所剩的元素就是需要添加的元素
	recodeIsAddVul := make(map[string]mysqls.TaskVul)
	for _, item := range vulCurrentRes {
		// 处理加密数据用于生成 singleKey
		targetUrl := item.TargetUrl
		name := item.Name
		location := item.Location

		if utils.IsHexString(item.TargetUrl) {
			targetUrlDecodeByte, _ := hex.DecodeString(item.TargetUrl)
			targetUrl = string(aesEcb.AesDecryptECB(targetUrlDecodeByte, key))
		}
		if utils.IsHexString(item.Name) {
			nameDecodeByte, _ := hex.DecodeString(item.Name)
			name = string(aesEcb.AesDecryptECB(nameDecodeByte, key))
		}
		if utils.IsHexString(item.Location) {
			locationDecodeByte, _ := hex.DecodeString(item.Location)
			location = string(aesEcb.AesDecryptECB(locationDecodeByte, key))
		}

		singleKey := strings.ReplaceAll(targetUrl+item.Pocname+name+location, " ", "")
		vulCurrentMaps[singleKey] = item
		recodeIsAddVul[singleKey] = item
	}

	//组装返回结果
	resp.Total = vulParentTotal
	list := make([]typespec.VulListRespItems, 0)
	for _, vulRes := range vulParentRes {
		var tmp typespec.VulListRespItems
		tmp.Id = vulRes.ID
		tmp.DataType = vulRes.DataType

		// 处理 TargetUrl
		if utils.IsHexString(vulRes.TargetUrl) {
			targetUrlDecodeByte, _ := hex.DecodeString(vulRes.TargetUrl)
			tmp.TargetUrl = string(aesEcb.AesDecryptECB(targetUrlDecodeByte, key))
		} else {
			tmp.TargetUrl = vulRes.TargetUrl
		}

		// 处理 Name
		if utils.IsHexString(vulRes.Name) {
			nameDecodeByte, _ := hex.DecodeString(vulRes.Name)
			tmp.Name = string(aesEcb.AesDecryptECB(nameDecodeByte, key))
		} else {
			tmp.Name = vulRes.Name
		}

		tmp.Type = vulRes.Type
		tmp.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(vulRes.Type)
		tmp.Risk = vulRes.Risk
		tmp.RiskName = enums.ToolsVulnerabilityEnum.GetRiskEnum(vulRes.Risk)

		// 处理 Location
		if utils.IsHexString(vulRes.Location) {
			locationDecodeByte, _ := hex.DecodeString(vulRes.Location)
			tmp.Location = string(aesEcb.AesDecryptECB(locationDecodeByte, key))
		} else {
			tmp.Location = vulRes.Location
		}

		tmp.IsSnapshot = enums.VulSnapshotNo
		if len(vulRes.Snapshot) > 0 {
			tmp.IsSnapshot = enums.VulSnapshotYes
		}

		// 验证当前漏洞的实际状态
		// 使用解密后的数据生成 singleKey
		singleKey := strings.ReplaceAll(tmp.TargetUrl+vulRes.Pocname+tmp.Name+tmp.Location, " ", "")
		if _, ok := vulCurrentMaps[singleKey]; !ok {
			// 说明此次验证结果内，没有检测到此漏洞，需要标记为已修复
			tmp.Status = enums.VulStatusRepairSuccess
			tmp.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(tmp.Status)
		} else {
			tmp.Status = vulRes.Status
			tmp.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(tmp.Status)
		}
		delete(recodeIsAddVul, singleKey)

		tmp.TestStatus = vulRes.TestStatus
		tmp.TestStatusName = enums.ToolsVulnerabilityEnum.GetVulTestStatusEnum(vulRes.TestStatus)

		// todo 兜底千机
		if tmp.Name == "qianji" {
			var resp typespec.VulInfoResp
			a.VulInfo(ctx, &typespec.VulInfoReq{TaskVulId: tmp.Id}, &resp)
			tmp.Name = resp.Name
			tmp.TypeName = resp.TypeName
			tmp.Location = resp.VulAddress
			tmp.RiskName = resp.RiskName
			tmp.Risk = resp.Risk
		}

		list = append(list, tmp)
	}

	resp.List = append(resp.List, list...)
	return nil
}

// VulInfo 漏洞测试详情
func (a *TaskCheckTask) VulInfo(ctx context.Context, req *typespec.VulInfoReq, resp *typespec.VulInfoResp) error {
	//查询数据
	var taskVulSrv services.TaskVul
	taskVulRes, err := taskVulSrv.VulInfo(ctx, req.TaskVulId)
	if err != nil {
		return err
	}
	if taskVulRes.ID == 0 {
		return errors.New("找不到该漏洞！")
	}

	// 销售许可证临时添加
	key := []byte("9876787656785679")
	aesEcb := aesEncryption.AesEcb{}

	//组装返回结果
	resp.Id = taskVulRes.ID
	resp.TaskId = taskVulRes.TaskID
	resp.TargetId = taskVulRes.TargetID

	// 处理 TargetUrl
	if utils.IsHexString(taskVulRes.TargetUrl) {
		targetUrlDecodeByte, _ := hex.DecodeString(taskVulRes.TargetUrl)
		resp.TargetUrl = string(aesEcb.AesDecryptECB(targetUrlDecodeByte, key))
	} else {
		resp.TargetUrl = taskVulRes.TargetUrl
	}

	// 处理 Pocname
	if utils.IsHexString(taskVulRes.Pocname) {
		pocnameDecodeByte, _ := hex.DecodeString(taskVulRes.Pocname)
		resp.Pocname = string(aesEcb.AesDecryptECB(pocnameDecodeByte, key))
	} else {
		resp.Pocname = taskVulRes.Pocname
	}

	// 处理 Name
	if utils.IsHexString(taskVulRes.Name) {
		nameDecodeByte, _ := hex.DecodeString(taskVulRes.Name)
		resp.Name = string(aesEcb.AesDecryptECB(nameDecodeByte, key))
	} else {
		resp.Name = taskVulRes.Name
	}

	resp.Type = taskVulRes.Type
	resp.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(taskVulRes.Type)
	resp.Risk = taskVulRes.Risk
	resp.RiskName = enums.ToolsVulnerabilityEnum.GetRiskEnum(taskVulRes.Risk)
	resp.VulNumber = taskVulRes.VulNumber
	resp.Cvss = taskVulRes.Cvss
	resp.PublishedTime = taskVulRes.PublishedTime
	resp.ExploitImpact = taskVulRes.ExploitImpact

	// Description 和 FixSuggest 不需要解密
	resp.Description = taskVulRes.Description
	resp.FixSuggest = taskVulRes.FixSuggest
	resp.RefUrl = taskVulRes.RefUrl
	resp.VulAddress = taskVulRes.VulAddress

	// 处理 VulResult
	if utils.IsHexString(taskVulRes.VulResult) {
		vulResultDecodeByte, _ := hex.DecodeString(taskVulRes.VulResult)
		resp.VulResult = string(aesEcb.AesDecryptECB(vulResultDecodeByte, key))
	} else {
		resp.VulResult = taskVulRes.VulResult
	}

	resp.VulParam = taskVulRes.VulParam

	// 处理 Location
	if utils.IsHexString(taskVulRes.Location) {
		locationDecodeByte, _ := hex.DecodeString(taskVulRes.Location)
		resp.Location = string(aesEcb.AesDecryptECB(locationDecodeByte, key))
	} else {
		resp.Location = taskVulRes.Location
	}

	resp.Status = taskVulRes.Status
	resp.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(taskVulRes.Status)
	resp.VulId = taskVulRes.VulID

	// 处理 VerMsg
	var tmpVerMsg = make([]typespec.VulInfoRespVerMsg, 0)
	if utils.IsHexString(taskVulRes.VerMsg) {
		verMsgDecodeByte, _ := hex.DecodeString(taskVulRes.VerMsg)
		json.Unmarshal(aesEcb.AesDecryptECB(verMsgDecodeByte, key), &tmpVerMsg)
	} else {
		json.Unmarshal([]byte(taskVulRes.VerMsg), &tmpVerMsg)
	}
	resp.VerMsg = tmpVerMsg
	// todo 兜底qianji ai特殊返回 情况特殊 先写到这里
	if resp.Pocname == "qianji" {
		var qianjiResult typespec.QianJiAIVulResult
		json.Unmarshal([]byte(resp.VulResult), &qianjiResult)
		resp.Location = qianjiResult.Location
		resp.Name = qianjiResult.AiVulName
		resp.FixSuggest = qianjiResult.FixSuggest
		resp.RefUrl = qianjiResult.ReferenceLink
		resp.VulResult = qianjiResult.Result
		resp.Description = qianjiResult.Description
		if qianjiResult.Description != "" {
			var qianjiDescInfo typespec.QianJiAIDescription
			json.Unmarshal([]byte(qianjiResult.Description), &qianjiDescInfo)
			resp.Description = qianjiDescInfo.Details
			if len(resp.VerMsg) == 1 {
				if resp.VerMsg[0].Response == `""` {
					resp.VerMsg[0].Response = qianjiDescInfo.Response
				}
				if resp.VerMsg[0].Request == `""` {
					resp.VerMsg[0].Request = qianjiDescInfo.Request
				}
			}
			resp.Cvss = qianjiDescInfo.Cvss
			// todo jixu说:cvss评分设置一下，高于7.5的是高 高于8.5的是致命 6-7.5 中危 低于6 低危
			cvss, err := strconv.ParseFloat(qianjiDescInfo.Cvss, 64)
			if err != nil {
				resp.RiskName = "未知"
				resp.Risk = enums.TargetRiskLowNoFound
			} else {
				resp.RiskName = "低危"
				resp.Risk = enums.TargetRiskLow
				if cvss >= 8.5 {
					resp.RiskName = "致命"
					resp.Risk = enums.TargetRiskHigh
				} else if cvss >= 7.5 {
					resp.RiskName = "高危"
					resp.Risk = enums.TargetRiskHigh
				} else if cvss > 6.0 {
					resp.RiskName = "中危"
					resp.Risk = enums.TargetRiskMid
				}
			}
			resp.VulAddress = qianjiDescInfo.Vuladdress
			resp.ExploitImpact = qianjiDescInfo.ExploitImpact
			resp.TypeName = qianjiDescInfo.VulType
		}
		resp.VulNumber = qianjiResult.VulNumber
	}
	return nil
}

// GetVulSnapshot 查看漏洞截图
func (a *TaskCheckTask) GetVulSnapshot(ctx context.Context, req *typespec.GetVulSnapshotReq, resp *typespec.GetVulSnapshotResp) error {
	var taskVulSrv services.TaskVul
	taskVulRes, err := taskVulSrv.VulInfo(ctx, req.TaskVulId)
	if err != nil {
		return err
	}
	if taskVulRes.ID == 0 || len(taskVulRes.Snapshot) == 0 {
		return errors.New("该漏洞没有截图数据...")
	}
	resp.Snapshot = taskVulRes.Snapshot
	return nil
}

// VulDel 漏洞测试删除
func (a *TaskCheckTask) VulDel(ctx context.Context, req *typespec.VulDelReq) error {
	idArray := strings.Split(req.TaskVulIds, ",")
	var (
		taskVulSrv    services.TaskVul
		taskTargetSrv services.TaskTarget
	)
	//查询被删除数据的目标id
	targetIds, err := taskVulSrv.GetVulTargetIdByIds(ctx, idArray)
	if err != nil {
		return err
	}
	//删除数据
	err = taskVulSrv.VulDel(ctx, idArray)
	if err != nil {
		return err
	}
	//更新目标表统计
	targetIds = data.ArrayIntUnique(targetIds) //targetIds去重
	for i := 0; i < len(targetIds); i++ {
		vulNum, risklevel, vulnumArray, err := taskVulSrv.GetTargetStats(ctx, targetIds[i])
		if err != nil {
			return err
		}
		err = taskTargetSrv.UpdateTargetStats(ctx, targetIds[i], vulNum, risklevel, vulnumArray)
		if err != nil {
			return err
		}
	}
	//更新任务表重新统计
	var taskSrv services.TaskTask
	err = taskSrv.UpdateIsStats(ctx, req.TaskId, enums.TaskIsStatsYes)
	if err != nil {
		return err
	}
	return nil
}

// VulTest 漏洞测试测试
func (a *TaskCheckTask) VulTest(ctx context.Context, req *typespec.VulTestReq, resp *typespec.VulTestResp) error {
	var taskVulSrv services.TaskVul
	//查询漏洞测试数据
	vulRes, err := taskVulSrv.VulInfo(ctx, req.TaskVulId)
	if err != nil || vulRes.ID == 0 {
		return errors.New("找不到该数据...")
	}
	//解析请求报文
	verMsg, err := taskVulSrv.Base64StdDecodeString(req.VerMsg) //base64解密
	if err != nil {
		return err
	}
	method, url, headerMap, body, err := taskVulSrv.DecryptionVerMsg(verMsg) //解析字符串结构
	//todo
	fmt.Println(headerMap)
	if err != nil {
		return err
	}
	//请求报文
	proto, status, header, bodyresult, err := httpclients.VerMsg(method, url, headerMap, body) //发送请求
	if err != nil {
		return err
	}
	respbody, err := taskVulSrv.BuildRespVerMsg(proto, status, header, bodyresult) //构建返回报文字符串
	if err != nil {
		return err
	}
	//修改task_vul表的报文数据
	err = taskVulSrv.UpdateVerMsgById(ctx, req.TaskVulId, verMsg, respbody, vulRes.VerMsg)
	if err != nil {
		return err
	}
	resp.RespVerMsg = taskVulSrv.Base64StdEncodeToString(respbody)
	return nil
}

// VulTestByYak 漏洞测试测试 通过yak库进行发包
func (a *TaskCheckTask) VulTestByYak(ctx context.Context, req *typespec.VulTestReq, resp *typespec.VulTestResp) error {
	var taskVulSrv services.TaskVul
	//查询漏洞测试数据
	vulRes, err := taskVulSrv.VulInfo(ctx, req.TaskVulId)
	if err != nil || vulRes.ID == 0 {
		return errors.New("找不到该数据...")
	}
	//解析请求报文
	verMsg, err := taskVulSrv.Base64StdDecodeString(req.VerMsg) //base64解密
	if err != nil {
		return err
	}
	opts := make([]poc.PocConfigOption, 0)
	opts = append(opts, poc.WithTimeout(20))
	httpsnum := strings.Count(verMsg, "https")
	if httpsnum > 0 {
		opts = append(opts, poc.WithForceHTTPS(true))
	}
	response, _, err := poc.HTTP(verMsg)
	if err != nil {
		return err
	}
	bodyStr := string(response)
	err = taskVulSrv.UpdateVerMsgById(ctx, req.TaskVulId, verMsg, bodyStr, vulRes.VerMsg)
	if err != nil {
		return err
	}
	resp.RespVerMsg = taskVulSrv.Base64StdEncodeToString(bodyStr)
	return nil
}

// VulVerify 漏洞测试验证
func (a *TaskCheckTask) VulVerify(ctx context.Context, req *typespec.VulVerifyReq) error {
	var (
		taskVulSrv services.TaskVul
		//taskTargetSrv services.TaskTarget
	)
	//查询漏洞信息
	vulInfo, err := taskVulSrv.VulInfo(ctx, req.TaskVulId)
	if err != nil {
		return err
	}
	//请求决策引擎重新验证漏洞
	var tmpVulParam []map[string]interface{}
	err = json.Unmarshal([]byte(req.VulParam), &tmpVulParam)
	if err != nil {
		return err
	}
	isVerify, _, err := httpclients.DecisionTaskVulVerify(req.Pocname, tmpVulParam)
	if err != nil {
		return err
	}
	//如果返回的验证为已修复，并且原来状态不是已经修复
	if isVerify == enums.VulStatusRepairSuccess && vulInfo.Status != enums.VulStatusRepairSuccess {
		err = taskVulSrv.UpdateStatusById(ctx, req.TaskVulId, isVerify)
		if err != nil {
			return err
		}
		//更新目标表统计数据
		//vulNum, risklevel, vulnumArray, err := taskVulSrv.GetTargetStats(ctx, vulInfo.TargetID)
		//if err != nil {
		//	return err
		//}
		//err = taskTargetSrv.UpdateTargetStats(ctx, vulInfo.TargetID, vulNum, risklevel, vulnumArray)
		//if err != nil {
		//	return err
		//}
		//更新任务表重新统计
		var taskSrv services.TaskTask
		err = taskSrv.UpdateIsStats(ctx, req.TaskId, enums.TaskIsStatsYes)
		if err != nil {
			return err
		}
	}
	return nil
}

// TestVulTest 待测漏洞测试
func (a *TaskCheckTask) TestVulTest(ctx context.Context, req *typespec.TestVulTestReq, taskVulIds []string) error {
	var (
		taskSrv       services.TaskTask
		taskTargetSrv services.TaskTarget
		taskLogtSrv   services.TaskLog
	)
	taskRes, err := taskSrv.GetTaskByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	if taskRes.ID == 0 {
		return errors.New("找不到任务数据！")
	}
	//查询任务状态
	if taskRes.Status != enums.TaskStatusRunning && taskRes.Status != enums.TaskStatusFinish {
		return errors.New("只有运行中或已结束的任务才可以执行测试漏洞！")
	}
	//组装漏洞数据
	var taskVulSrv services.TaskVul
	vulRes, err := taskVulSrv.GetVulByIds(ctx, taskVulIds)
	if err != nil {
		return err
	}
	if len(vulRes) == 0 {
		return errors.New("找不到漏洞数据...")
	}
	var targetIds = make([]int, 0)
	for i := 0; i < len(vulRes); i++ {
		if vulRes[i].TestStatus == enums.VulTestStatusNotTest { //只有未测试的漏洞可以提交测试
			targetIds = append(targetIds, vulRes[i].TargetID)
		}
	}
	if len(targetIds) != len(vulRes) {
		return errors.New("已经提交测试的漏洞不可以重复测试...")
	}
	//查询日志id
	targetlogmapRes, err := taskLogtSrv.GetManyTaskLogsByTargets(ctx, targetIds)
	if err != nil || len(targetlogmapRes) == 0 {
		return errors.New("查询该任务日志失败！")
	}
	//发送到决策引擎
	taskMsgArray, err := taskVulSrv.TestVulCreatDecisionParams(vulRes, targetlogmapRes)
	if err != nil {
		return err
	}
	for i := 0; i < len(taskMsgArray); i++ {
		err = taskTargetSrv.SendDecisionMsg(taskMsgArray[i])
		if err != nil {
			log.Info("TestVulTest SendDecisionMsg err:", err)
			continue
		}
		time.Sleep(100 * time.Millisecond)
	}
	//更改测试状态测试状态
	err = taskVulSrv.TestVulUpdateStatus(ctx, vulRes, enums.VulTestStatusTesting)
	if err != nil {
		return err
	}
	//更新任务表重新统计
	err = taskSrv.UpdateTaskStateIsStatsById(ctx, req.TaskId, enums.TaskStatusRunning, enums.TaskIsStatsYes)
	if err != nil {
		return err
	}
	//更新测试目标状态
	err = taskTargetSrv.UpdateTargetStateByIds(ctx, targetIds, enums.TargetStatusRunning)
	if err != nil {
		return err
	}
	return nil
}

// AttackLink 攻击链路图
func (a *TaskCheckTask) AttackLink(ctx context.Context, req *typespec.AttackLinkReq, resp *typespec.AttackLinkResp) error {
	var taskTargetSrv services.TaskTarget
	targetList, _ := taskTargetSrv.GetTargetByTaskId(ctx, req.TaskId)
	nodes := make([]map[string]interface{}, 0)
	edges := make([]map[string]interface{}, 0)
	initNode := map[string]interface{}{
		"id":   0,
		"name": "初始节点",
		"type": "image",
		"size": "100",
		"style": map[string]string{
			"src": "/img2/1.png",
		},
	}
	nodes = append(nodes, initNode)
	for _, target := range targetList {
		png := "/img2/4.png"
		if strings.Contains(target.OpSys, "win") {
			png = "/img2/5.png"
		} else if strings.Contains(target.OpSys, "kylin") {
			png = "/img2/2.png"
		} else if target.RiskLevel > 4 && strings.Contains(target.OpSys, "win") {
			png = "/img2/3.png"
		}
		node := map[string]interface{}{
			"id":   strconv.Itoa(target.ID),
			"name": utils.GetHostname(target.TargetURL),
			"type": "image",
			"size": "100",
			"style": map[string]string{
				"src": png,
			},
		}
		nodes = append(nodes, node)
		source := "0"
		if target.ExtendField != "" {
			source = target.ExtendField
		}
		edge := map[string]interface{}{
			"id":     strconv.Itoa(target.ID),
			"source": source,
			"target": strconv.Itoa(target.ID),
		}
		edges = append(edges, edge)
	}
	resp.Nodes = nodes
	resp.Edges = edges
	return nil
}

// LogList 测试日志列表及筛选
func (a *TaskCheckTask) LogList(ctx context.Context, req *typespec.LogListReq, resp *typespec.LogListResp) error {
	//查询数据
	var logsrv services.TaskLog
	logRes, total, err := logsrv.LogList(ctx, req.TaskId, req.Search, req.Page, req.Size)
	if err != nil {
		return err
	}
	//组装返回结果
	resp.Total = total
	resp.List = make([]typespec.LogListRespItems, 0)
	for i := 0; i < len(logRes); i++ {
		var tmp typespec.LogListRespItems
		tmp.Id = logRes[i].ID
		tmp.TaskId = logRes[i].TaskID
		tmp.TargetId = logRes[i].TargetID
		tmp.TargetUrl = logRes[i].TargetURL
		tmp.Status = logRes[i].Status
		tmp.StatusName = enums.GetTargetStatus(logRes[i].Status)
		tmp.IsAlive = logRes[i].IsAlive
		tmp.IsAliveName = enums.GetTargetIsAlive(logRes[i].IsAlive)
		tmp.CreateTime = logRes[i].CreateTime.Format(enums.TimeLayout)
		tmp.StartTime = logRes[i].StartTime.Format(enums.TimeLayout)
		tmp.EndTime = logRes[i].EndTime.Format(enums.TimeLayout)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// LogInfo 测试日志详情
func (a *TaskCheckTask) LogInfo(ctx context.Context, req *typespec.LogInfoReq, resp *typespec.LogInfoResp) error {
	//查询数据
	var logsrv services.TaskLog
	loginfores, total, err := logsrv.LogInfo(ctx, req.TaskLogId)
	if err != nil {
		return err
	}
	//组装返回结果
	resp.Total = total
	resp.List = make([]typespec.LogInfoRespItems, 0)
	for i := 0; i < len(loginfores); i++ {
		var tmp typespec.LogInfoRespItems
		tmp.Id = loginfores[i].ID
		tmp.TaskId = loginfores[i].TaskID
		tmp.TargetId = loginfores[i].TargetID
		tmp.TargetUrl = loginfores[i].TargetURL
		tmp.Pocname = loginfores[i].Pocname
		tmp.Result = loginfores[i].Result
		tmp.CreateTime = loginfores[i].CreateTime.Format(enums.TimeLayout)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// ApiSave 通过api调用创建任务
func (a *TaskCheckTask) ApiSave(ctx context.Context, req *typespec.ApiTaskTaskCreateReq, resp *typespec.TaskSaveResp, uid int, runtimeCheck data.TaskRuntimeCheck) error {
	/******** 数据是否合法 ***********/
	// 校验场景是否存在
	var srvTemplate services.SceneTaskTemplate
	template, err := srvTemplate.GetByTemplateId(ctx, req.TaskTemplateId, enums.TaskTemplateStatusSuccess)
	if err != nil {
		return err
	}
	if template.ID == 0 {
		return errors.New("任务场景不存在")
	}
	taskConfigList := srvTemplate.GetTaskConfigById(ctx, req.TaskTemplateId)

	configJson := enums.ConfigJson{}
	if req.Config.PortScanConfig.ScanPort == "" {
		for _, taskConfig := range taskConfigList {
			if taskConfig.ConfigKey == enums.ConfigJsonPortScanKey {
				var tempConfigJson enums.PortScanConfig
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.PortScanConfig = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigJsonWebCrawlerKey {
				var tempConfigJson enums.WebCrawlerConfig
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.WebCrawlerConfig = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigJsonPathScanKey {
				var tempConfigJson enums.WebPathScanConfig
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.WebPathScanConfig = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigJsonWeakPass {
				var tempConfigJson enums.WeakPassConfig
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.WeakPassConfig = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigJsonSubdomainCollectKey {
				var tempConfigJson enums.SubdomainCollectConfig
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.SubdomainCollectConfig = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigJsonWebsiteLoginKey {
				var tempConfigJson enums.WebsiteLoginConfig
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.WebsiteLoginConfig = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigJsonVulIdsKey {
				var tempConfigJson []int
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.VulIdsConfig = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigJsonAliveProbeConfig {
				var tempConfigJson enums.AliveProbeConfig
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.AliveProbeConfig = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigSafeTestKey {
				var tempConfigJson bool
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.SafeTest = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigLateralMoveKey {
				var tempConfigJson enums.LateralMove
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.LateralMove = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigTestIntensityKey {
				var tempConfigJson int
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.TestIntensity = tempConfigJson
			} else if taskConfig.ConfigKey == enums.ConfigVulExploitKey {
				var tempConfigJson bool
				err := json.Unmarshal([]byte(taskConfig.ConfigJson), &tempConfigJson)
				if err != nil {
					return err
				}
				configJson.VulExploit = tempConfigJson
			}
		}
	} else {
		configJson = req.Config
	}

	// 判断config中的请求参数是否合规，并把需要默认值的参数赋值
	configStruct, err := data.TaskCheckTaskConfig.VerifyConfig(configJson)
	if err != nil {
		return err
	}
	/******** 数据是否合法 end ***********/

	/******** 非标准数据处理 ***********/
	taskName := req.TaskName
	taskName = strings.ReplaceAll(taskName, " ", "")
	taskName = strings.ReplaceAll(taskName, "\r", "")
	taskName = strings.ReplaceAll(taskName, "\n", "")

	// 分析目标
	var analysisTarget data.TaskCheckTaskAnalysisTarget
	analysisTarget.AnalysisTarget(req.Target, "")
	errorTargetList := analysisTarget.ErrorTargetList
	if len(errorTargetList) > 0 {
		return errors.New(strings.Join(errorTargetList, ","))
	}
	targetList := analysisTarget.TargetList
	//一个任务的目标总数不能大于3200
	if len(targetList) > 3200 {
		return errors.New("检测目标超过最大个数: 3200")
	}
	/******** 非标准数据处理end ***********/

	// 计算任务什么时候运行
	var taskExecNextTime data.TaskExecNextTime
	taskExecNextTime.StartTime, _ = time.ParseInLocation(utils.DateTime, runtimeCheck.StartTime, time.Local)
	taskExecNextTime.CyclePlanningType = runtimeCheck.CyclePlanningType
	taskExecNextTime.CyclePlanningValue = runtimeCheck.CyclePlanningValue
	taskExecNextTime.CyclePlanningHour = runtimeCheck.CyclePlanningHour
	taskExecNextTime.EndTime, _ = time.ParseInLocation(utils.DateTime, runtimeCheck.EndTime, time.Local)
	nextRuntime, err := taskExecNextTime.Compute(enums.TaskExecTypeImmediate)
	if err != nil {
		return err
	}
	// 执行方式对应的参数
	executeJsonByte, _ := json.Marshal(runtimeCheck)

	// 创建
	var srvTask services.TaskTask
	resp.TaskId, err = srvTask.Save(ctx, uid, req.TaskTemplateId, taskName, targetList, enums.TaskExecTypeImmediate, string(executeJsonByte), configStruct, nextRuntime, enums.TaskWeightMid, 0)

	if err != nil {
		return err
	}

	return nil
}

// GetTaskProgress 通过api调用创建任务
func (a *TaskCheckTask) GetTaskProgress(ctx context.Context, req *typespec.TaskProgressReq, resp *typespec.TaskProgressResp) error {
	/******** 数据是否合法 ***********/

	var (
		taskTaskInfoModel mysqls.TaskTaskInfo
		taskTaskModel     mysqls.TaskTask
	)
	taskTask, err := taskTaskModel.GetTaskCheckTask(ctx, req.TaskId)
	if taskTask.ID == 0 {
		return errors.New("该任务不存在")
	}
	if taskTask.Status == enums.TaskStatusBegin {
		resp.Progress = "0%"
		return nil
	} else if taskTask.Status == enums.TaskStatusFinish {
		resp.Progress = "100%"
		return nil
	}
	taskTaskInfo, err := taskTaskInfoModel.GetTaskTaskInfoByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	var overViewStruct services.OverView
	if taskTaskInfo.Overview == "" {
		resp.Progress = "1%"
		return nil
	}
	err = json.Unmarshal([]byte(taskTaskInfo.Overview), &overViewStruct)
	if err != nil {
		return err
	}
	if overViewStruct.Progress == nil {
		resp.Progress = "1%"
		return nil
	}

	resp.Progress = strconv.Itoa(overViewStruct.Progress.Value) + "%"
	return nil
}

// ApiVulList 漏洞测试列表及筛选
func (a *TaskCheckTask) ApiVulList(ctx context.Context, req *typespec.ApiVulListReq, resp *typespec.ApiVulListResp) error {
	//查询数据
	var taskVulSrv services.TaskVul
	vulRes, total, err := taskVulSrv.VulList(ctx, req.TaskId, 0, req.Type, req.Risk, req.Search, enums.VulDataTypOne, req.Page, req.Size, 0)
	if err != nil {
		return err
	}

	// 销售许可证临时添加
	key := []byte("9876787656785679")
	aesEcb := aesEncryption.AesEcb{}

	//组装返回结果
	resp.Total = total
	resp.List = make([]typespec.ApiVulListRespItems, 0)
	domainMap := make(map[string]string, 0)
	for i := 0; i < len(vulRes); i++ {
		var tmp typespec.ApiVulListRespItems
		tmp.Id = vulRes[i].ID
		//tmp.TargetUrl = vulRes[i].TargetUrl
		//tmp.Name = vulRes[i].Name
		nameDecodeByte, _ := hex.DecodeString(vulRes[i].Name)
		tmp.Name = string(aesEcb.AesDecryptECB(nameDecodeByte, key))
		targetUrlByte, _ := hex.DecodeString(vulRes[i].TargetUrl)
		tmp.TargetUrl = string(aesEcb.AesDecryptECB(targetUrlByte, key))

		tmp.Type = vulRes[i].Type
		tmp.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(vulRes[i].Type)
		tmp.Risk = vulRes[i].Risk
		tmp.RiskName = enums.ToolsVulnerabilityEnum.GetRiskEnum(vulRes[i].Risk)
		//tmp.Location = vulRes[i].Location
		locationByte, _ := hex.DecodeString(vulRes[i].Location)
		tmp.Location = string(aesEcb.AesDecryptECB(locationByte, key))
		tmp.Status = vulRes[i].Status
		tmp.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(vulRes[i].Status)

		tmp.TaskId = vulRes[i].TaskID
		tmp.TargetId = vulRes[i].TargetID
		//tmp.Pocname = vulRes[i].Pocname
		pocnameByte, _ := hex.DecodeString(vulRes[i].Pocname)
		tmp.Pocname = string(aesEcb.AesDecryptECB(pocnameByte, key))

		tmp.Cve = strings.Split(vulRes[i].VulID, ",")[0]
		tmp.Cvss = vulRes[i].Cvss
		tmp.PublishedTime = vulRes[i].PublishedTime
		tmp.ExploitImpact = vulRes[i].ExploitImpact
		tmp.Description = vulRes[i].Description
		tmp.FixSuggest = vulRes[i].FixSuggest
		tmp.RefUrl = vulRes[i].RefUrl
		tmp.VulAddress = vulRes[i].VulAddress
		//tmp.VulResult = vulRes[i].VulResult
		vulResultByte, _ := hex.DecodeString(vulRes[i].VulResult)
		tmp.VulResult = string(aesEcb.AesDecryptECB(vulResultByte, key))

		tmp.VulParam = vulRes[i].VulParam
		//tmp.VerMsg = vulRes[i].VerMsg
		verMsgByte, _ := hex.DecodeString(vulRes[i].VerMsg)
		tmp.VerMsg = string(aesEcb.AesDecryptECB(verMsgByte, key))

		tmp.Class = vulRes[i].Class
		//if tmp.Class == enums.VulLibrariesClassWindowsSystem || tmp.Class == enums.VulLibrariesClassLinuxOrUnixSystem || tmp.Class == enums.VulLibrariesClassSystem {
		tmp.Ip, tmp.Port = taskVulSrv.GetParamIpPort(vulRes[i].VulParam, domainMap)
		//}
		if tmp.Location == "" {
			tmp.Location = taskVulSrv.GetParamLocation(vulRes[i].VulParam)
		}

		resp.List = append(resp.List, tmp)
	}
	return nil
}

// AddTarget 动态添加目标
func (a *TaskCheckTask) AddTarget(ctx context.Context, req *typespec.AddTargetReq) error {
	//查询任务
	var taskSrv services.TaskTask
	taskRes, taskinfoRes, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	if taskRes.ID == 0 || taskinfoRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	//解析配置
	var conf enums.ConfigJson
	err = json.Unmarshal([]byte(taskinfoRes.TaskTemplateJSON), &conf)
	if err != nil {
		return errors.New("任务配置参数解析错误！")
	}
	//解析目标ip
	var analysisTarget data.TaskCheckTaskAnalysisTarget
	analysisTarget.AnalysisTarget(req.Target, "")
	errorTargetList := analysisTarget.ErrorTargetList
	if len(errorTargetList) > 0 {
		return errors.New(strings.Join(errorTargetList, ","))
	}
	targetList := analysisTarget.TargetList
	//测试目标黑白名单过滤
	var mapSet services.MapSet
	targetList, err = mapSet.TargetIpWhiteBlack(ctx, targetList)
	if err != nil {
		return err
	}
	//一个任务的目标总数不能大于3200
	if len(targetList) > 3200 || len(targetList) == 0 {
		return errors.New("检测目标超过最大个数3200或开启了测试目标黑白名单检测目标为0")
	}
	err = taskSrv.AddTarget(ctx, targetList, taskRes, taskinfoRes, req.UserId)
	if err != nil {
		return err
	}
	return nil
}

// AddAttackFace 动态添加攻击面
func (a *TaskCheckTask) AddAttackFace(ctx context.Context, req *typespec.AddAttackFaceReq) error {
	urlObj, err := url.Parse(req.Target)
	if err != nil {
		return err
	}
	hostArray := strings.Split(urlObj.Host, ":")
	if len(hostArray) == 0 {
		return errors.New("漏洞路径错误！")
	}
	if len(hostArray) > 1 {
		return errors.New("带端口的目标无法添加攻击面...")
	}
	//解析参数
	var (
		port      int
		infoType  string
		infoUrl   string
		infoIp    string
		infoPort  string
		infoValue string
	)

	if req.AttackFaceType == enums.AttackFaceTypePort {
		var attackfaceport typespec.AddAttackFacePort
		err := json.Unmarshal([]byte(req.Params), &attackfaceport)
		if err != nil {
			return err
		}
		if attackfaceport.Port == 0 {
			return errors.New("添加端口供给面，端口参数不能为空")
		}
		port = attackfaceport.Port
	} else if req.AttackFaceType == enums.AttackFaceTypePath {
		var attackfacepath typespec.AddAttackFacePath
		err := json.Unmarshal([]byte(req.Params), &attackfacepath)
		if err != nil {
			return err
		}
		if attackfacepath.Url == "" || attackfacepath.Type == "" {
			return errors.New("动态添加敏感路径，路径参数不能为空")
		}
		infoType = attackfacepath.Type
		infoUrl = attackfacepath.Url
	} else if req.AttackFaceType == enums.AttackFaceTypeLoginCred {
		var attackfacelogincred typespec.AddAttackFaceLoginCred
		err := json.Unmarshal([]byte(req.Params), &attackfacelogincred)
		if err != nil {
			return err
		}
		if attackfacelogincred.Ip == "" || attackfacelogincred.Port == "" || attackfacelogincred.Type == "" {
			return errors.New("动态添加敏感路径，路径参数不能为空")
		}
		infoIp = attackfacelogincred.Ip
		infoPort = attackfacelogincred.Port
		infoType = attackfacelogincred.Type
		infoValue = attackfacelogincred.Value
	}

	//查询任务信息
	var (
		taskSrv       services.TaskTask
		taskTargetSrv services.TaskTarget
		taskLogSrv    services.TaskLog
		msg           services.TaskControlMessage
	)
	taskRes, taskinfoRes, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	if taskRes.ID == 0 || taskinfoRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	if taskRes.Status == enums.TaskStatusTrigger || taskRes.Status == enums.TaskStatusBegin || taskRes.Status == enums.TaskStatusPausing {
		return errors.New("待开始和暂停中的任务不能添加攻击面!")
	}
	//查询任务目标
	targetRes, err := taskTargetSrv.GetTargetsByTargetURL(ctx, req.TaskId, req.Target)
	if err != nil || targetRes.ID == 0 {
		return errors.New("查询该任务目标失败！")
	}
	//查询日志表
	taskLogRes, err := taskLogSrv.GetTaskLog(ctx, targetRes.ID)
	if err != nil || taskLogRes.ID == 0 {
		return errors.New("查询该任务日志失败！")
	}
	if req.AttackFaceType == enums.AttackFaceTypePort { // 动态添加端口扫描
		msg, err = taskTargetSrv.CreateAddInfoMsg(targetRes, port, taskLogRes.ID)
	} else if req.AttackFaceType == enums.AttackFaceTypePath { //动态添加敏感路径
		msg, err = taskTargetSrv.CreateAddInfoMsgPath(targetRes, infoType, infoUrl, taskLogRes.ID)
	} else if req.AttackFaceType == enums.AttackFaceTypeLoginCred { //动态添加登录凭证敏感路径
		msg, err = taskTargetSrv.CreateAddInfoMsgLoginCred(targetRes, infoIp, infoPort, infoType, infoValue, taskLogRes.ID)
	}
	//发送到决策引擎
	if err != nil {
		return err
	}
	err = taskTargetSrv.SendDecisionMsg(msg)
	if err != nil {
		return err
	}
	//更新任务表重新统计
	err = taskSrv.UpdateTaskStateIsStatsById(ctx, req.TaskId, enums.TaskStatusRunning, enums.TaskIsStatsYes)
	if err != nil {
		return err
	}
	//更新测试目标状态
	err = taskTargetSrv.UpdateTargetStateByIds(ctx, []int{targetRes.ID}, enums.TargetStatusRunning)
	if err != nil {
		return err
	}
	//更新测试目标状态
	err = taskLogSrv.UpdateTaskLogStateByTargetId(ctx, targetRes.ID, enums.TargetStatusRunning)
	if err != nil {
		return err
	}
	return nil
}

// AddVul 动态添加漏洞
func (a *TaskCheckTask) AddVul(ctx context.Context, req *typespec.AddVulReq) error {
	var (
		taskSrv        services.TaskTask
		taskTargetSrv  services.TaskTarget
		taskLogSrv     services.TaskLog
		taskLogInfoSrv services.TaskLogInfo
	)
	//解析漏洞路径
	urlObj, err := url.Parse(req.RootUrl)
	if err != nil {
		return err
	}
	hostArray := strings.Split(urlObj.Host, ":")
	if len(hostArray) == 0 {
		return errors.New("漏洞路径错误！")
	}
	//查询任务信息
	taskRes, taskinfoRes, err := taskSrv.GetTaskTaskinfoByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}
	if taskRes.ID == 0 || taskinfoRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	if taskRes.Status == enums.TaskStatusTrigger || taskRes.Status == enums.TaskStatusBegin || taskRes.Status == enums.TaskStatusPausing {
		return errors.New("待开始和暂停中的任务不能添加漏洞!")
	}
	//查询任务目标
	targetRes, err := taskTargetSrv.GetTargetsByTargetURLLike(ctx, req.TaskId, hostArray[0])
	if err != nil || targetRes.ID == 0 {
		return errors.New("查询该任务目标失败！")
	}
	//查询日志表
	taskLogRes, err := taskLogSrv.GetTaskLog(ctx, targetRes.ID)
	if err != nil || taskLogRes.ID == 0 {
		return errors.New("查询该任务日志失败！")
	}
	//请求决策引擎验证漏洞
	var tmpVulParam = []map[string]interface{}{
		map[string]interface{}{
			"key":   "root_url",
			"value": req.RootUrl,
		},
	}
	_, content, err := httpclients.DecisionTaskVulVerify(req.Pocname, tmpVulParam)
	if err != nil {
		return err
	}
	content = fmt.Sprintf(`动态添加漏洞,目标:%s,pocname:%s,漏洞结果:%s`, targetRes.TargetURL, req.Pocname, content)
	//存入日志
	err = taskLogInfoSrv.AddTaskLogInfo(ctx, req.TaskId, targetRes.ID, taskLogRes.ID, targetRes.TargetURL, req.Pocname, content)
	if err != nil {
		return err
	}
	return nil
}

// 目标路径图
func (a *TaskCheckTask) TaskTargetMap(ctx context.Context, req *typespec.TaskTargetMapReq, res *typespec.TaskTargetMapRes) error {
	// 获取任务数据
	var (
		taskTargetSrv services.TaskTarget
		taskVulSrv    services.TaskVul
	)
	target, err := taskTargetSrv.GetTargetById(ctx, req.TargetId)
	if err != nil {
		return err
	}
	// 组合风险等级
	_, _, vulNumArray, err := taskVulSrv.GetTargetStats(ctx, req.TargetId)
	if err != nil {
		return err
	}
	res.Risk.Deadly = vulNumArray[1]
	res.Risk.High = vulNumArray[2]
	res.Risk.Middle = vulNumArray[3]
	res.Risk.Low = vulNumArray[4]
	// 任务进度
	res.Progress = enums.GetTargetStatus(target.Status)
	// 是否有远控
	res.RemoteControlNum = enums.GetTargetIsRemoteSession(target.IsRemoteSession)
	// 获取网络状态
	// 1 获取网络状态 -> 本地网
	_, localQuality, localDelay, err := network.LocalPing()
	if err == nil {
		res.Network.LocalQuality = localQuality
		res.Network.LocalDelay = localDelay
	}
	// 2 获取网络状态 -> 目标网络连接情况
	_, host, _, err := network.ParseUrl(target.TargetURL)
	if err == nil {
		inAddr := make([]string, 0)
		inAddr = append(inAddr, host)
		network, targetQuality, targetDelay, err := network.Ping(inAddr)
		if err == nil {
			res.Network.Network = network
			res.Network.TargetQuality = targetQuality
			res.Network.TargetDelay = targetDelay
		}
	}
	/************ 组织路径图 ***************/
	// 获取目标下的所有结果
	var taskResultSrv services.TaskResult
	results := taskResultSrv.AllTaskResultByTargetId(ctx, req.TargetId, "")
	if len(results) == 0 {
		return nil
	}
	// 造根节点
	res.State = append(res.State, typespec.TaskTargetMapState{
		Id:       "0",
		Label:    target.TargetURL,
		Status:   "end",
		Success:  true,
		Hide:     false,
		Pocname:  "",
		Class:    "info",
		DetailId: "root_0_" + strconv.Itoa(target.ID),
	})

	// 漏洞节点是否在简化路径时隐藏 如果vul表中有数据则不隐藏否则隐藏
	resultIds := make([]int, 0)
	for _, item := range results {
		resultIds = append(resultIds, item.ID)
	}
	taskVuls := taskVulSrv.GetByTargetResultIds(ctx, resultIds, req.TargetId)
	taskVulsMap := make(map[int]bool)
	for _, item := range taskVuls {
		taskVulsMap[item.TargetResultID] = true
	}

	// 遍历结果数据，构建节点
	for _, item := range results {
		// 该节点是否显示
		hide := false
		if _, ok := taskVulsMap[item.ID]; !ok {
			// 特殊节点不隐藏
			if item.Pocname != "http_base_info_get" &&
				item.Pocname != "fingerprint" &&
				item.Pocname != "crawlerx" &&
				item.Pocname != "tempmitm" &&
				item.Pocname != "cdn_detect" &&
				item.Pocname != "web_scanner" &&
				item.Pocname != "waf_detect" &&
				item.Pocname != "whois" &&
				item.Pocname != "port_scan" &&
				item.Pocname != "principle_verify" &&
				item.Pocname != "attacker_ip_determination" {
				hide = true
			}
		}
		// 设置节点样式类
		class := "info"
		if item.Checked == enums.TaskResultCheckedY {
			if strings.Contains(item.VulName, "弱口令") {
				class = "high"
			} else {
				switch item.Risk {
				case enums.VulLibrariesRiskDead:
					class = "deadly"
				case enums.VulLibrariesRiskHigh:
					class = "high"
				case enums.VulLibrariesRiskMiddle:
					class = "middle"
				case enums.VulLibrariesRiskLow:
					class = "low"
				}
			}
		}
		// 设置节点标签
		labelName := item.Pocname
		if item.VulName != "" {
			labelName = item.VulName
		}
		// 设置详情ID后缀
		detailIdSuffix := "vul_"
		if item.Pocname == "cdn_detect" {
			detailIdSuffix = "cdn_"
		} else if item.Pocname == "web_dir_path_scan" {
			detailIdSuffix = "webdirpathscan_"
		} else if item.Pocname == enums.ScriptNameWhois {
			detailIdSuffix = "whois_"
		} else if item.Pocname == "port_scan" {
			detailIdSuffix = "portScan_"
		} else if item.Pocname == "principle_verify" {
			detailIdSuffix = "principle_"
		} else if item.Pocname == "attacker_ip_determination" {
			detailIdSuffix = "attacker_"
		}
		// 添加节点
		res.State = append(res.State, typespec.TaskTargetMapState{
			Id:       item.NodeID,
			Label:    labelName,
			Status:   "end",
			Success:  true,
			Hide:     hide,
			Pocname:  item.Pocname,
			Class:    class,
			DetailId: detailIdSuffix + strconv.Itoa(item.ID) + "_" + strconv.Itoa(item.TargetID),
		})
		// 处理父子关系，构建连接
		if item.FatherID != "" && item.FatherID != "0" {
			// 检测是否有多个上级
			fatherIds := make([]string, 0)
			if strings.Contains(item.FatherID, ",") {
				fatherIds = strings.Split(item.FatherID, ",")
			} else {
				fatherIds = append(fatherIds, item.FatherID)
			}
			for _, fatherId := range fatherIds {
				res.Edg = append(res.Edg, typespec.TaskTargetMapEdg{
					Start: fatherId,
					End:   item.NodeID,
					Hide:  hide,
				})
			}
		} else if item.FatherID == "0" {
			// 直接连接到根节点
			res.Edg = append(res.Edg, typespec.TaskTargetMapEdg{
				Start: "0",
				End:   item.NodeID,
				Hide:  hide,
			})
		}
	}
	return nil
}

// TaskTargetMapNodeDetail 路径节点详情
func (a *TaskCheckTask) TaskTargetMapNodeDetail(ctx context.Context, req *typespec.TaskTargetMapNodeDetailReq, res *typespec.TaskTargetMapNodeDetailRes) error {
	detailId := strings.Split(req.DetailId, "_")
	// 结果ID
	resultId, err := strconv.Atoi(detailId[1])
	if err != nil {
		return err
	}
	// 目标ID
	targetId, err := strconv.Atoi(detailId[2])
	if err != nil {
		return err
	}
	switch detailId[0] {
	case "root":
		// 根节点
		return a.taskTargetMapNodeDetailRoot(ctx, targetId, res)
	case "portScan":
		// 端口扫描节点
		return a.taskTargetMapNodeDetailPortScan(ctx, targetId, res)
	case "port":
		// 端口节点 会有地4个元素：端口
		if len(detailId) != 4 {
			return errors.New("端口节点，缺少端口参数")
		}
		return a.taskTargetMapNodeDetailPort(ctx, resultId, detailId[3], res)
	case "vul":
		// 漏洞节点
		return a.taskTargetMapNodeDetailVul(ctx, resultId, res)
	case "cdn":
		return a.taskTargetMapNodeDetailResult(ctx, "cdn", resultId, res)
	case "webdirpathscan":
		return a.taskTargetMapNodeDetailTaskResult(ctx, "webdirpathscan", resultId, res)
	case "whois":
		return a.taskTargetMapNodeDetailWhois(ctx, resultId, res)
	case "principle":
		return a.taskTargetMapNodeDetailPrinciple(ctx, resultId, res)
	}
	return nil
}

// taskTargetMapNodeDetailRoot 路径节点详情 - 根节点
func (a *TaskCheckTask) taskTargetMapNodeDetailRoot(ctx context.Context, targetId int, res *typespec.TaskTargetMapNodeDetailRes) error {
	var taskTargetSrv services.TaskTarget
	target, err := taskTargetSrv.GetTargetById(ctx, targetId)
	if err != nil {
		return err
	}

	var configStruct enums.ConfigJson
	err = json.Unmarshal([]byte(target.TaskTemplateJSON), &configStruct)
	if err != nil {
		return err
	}

	var taskTemplateSrv services.SceneTaskTemplate
	taskTemplate, err := taskTemplateSrv.GetTaskTemplateById(ctx, target.TaskTemplateID)
	if err != nil {
		return err
	}

	var rootNode typespec.TaskTargetMapNodeDetailRoot
	rootNode.Title = target.TargetURL
	rootNode.TemplateName = taskTemplate.TemplateName
	rootNode.PortScan = configStruct.PortScanConfig.ScanPort

	res.Type = "root"
	res.Data = rootNode
	return nil
}

// taskTargetMapNodeDetailPortScan 路径节点详情 - 端口扫描节点
func (a *TaskCheckTask) taskTargetMapNodeDetailPortScan(ctx context.Context, targetId int, res *typespec.TaskTargetMapNodeDetailRes) error {

	var portScanNode typespec.TaskTargetMapNodeDetailPortScan
	portScanNode.Title = "端口扫描"
	var taskResultSrv services.TaskResult
	portScanResult := taskResultSrv.AllTaskResultByTargetId(ctx, targetId, "port_scan")
	for _, item := range portScanResult {
		portResultMap := make([]map[string]string, 0)
		json.Unmarshal([]byte(item.Result), &portResultMap)
		for _, port := range portResultMap {
			portScanNode.Ports = append(portScanNode.Ports, typespec.TaskTargetMapNodeDetailPortItem{
				Port:    port["port"],
				Service: port["service"],
			})
		}

	}

	res.Type = "portScan"
	res.Data = portScanNode
	return nil
}

// taskTargetMapNodeDetailPort 路径节点详情 - 端口节点
func (a *TaskCheckTask) taskTargetMapNodeDetailPort(ctx context.Context, resultId int, port string, res *typespec.TaskTargetMapNodeDetailRes) error {
	var portNode typespec.TaskTargetMapNodeDetailPort
	var taskResultSrv services.TaskResult
	portResult, err := taskResultSrv.GetTaskResultById(ctx, resultId)
	if err != nil {
		return err
	}

	// 解析端口扫描结果
	portResultMap := make([]map[string]string, 0)
	json.Unmarshal([]byte(portResult.Result), &portResultMap)

	// 标题为端口
	portNode.Title = port
	for _, ports := range portResultMap {
		if ports["port"] == port {
			// 服务
			portNode.Service = ports["service"]
			// 结果数据
			portStr, _ := json.Marshal(ports)
			portNode.Result = string(portStr)
		}
	}
	res.Type = "port"
	res.Data = portNode
	return nil
}

func (a *TaskCheckTask) taskTargetMapNodeDetailResult(ctx context.Context, nodeType string, resultId int, res *typespec.TaskTargetMapNodeDetailRes) error {
	// 获取结果数据
	var taskResultSrv services.TaskResult
	taskResult, err := taskResultSrv.GetTaskResultById(ctx, resultId)
	if err != nil {
		return err
	}
	var tmpData = map[string]interface{}{
		"title":     "CDN",
		"resovleIP": "",
	}
	if len(taskResult.Result) != 0 {
		err = json.Unmarshal([]byte(taskResult.Result), &tmpData)
		if err != nil {
			return err
		}
	}
	res.Type = nodeType
	res.Data = tmpData
	return nil
}

func (a *TaskCheckTask) taskTargetMapNodeDetailTaskResult(ctx context.Context, nodeType string, resultId int, res *typespec.TaskTargetMapNodeDetailRes) error {
	res.Type = nodeType

	// 获取结果数据
	var taskResultSrv services.TaskResult
	taskResult, err := taskResultSrv.GetTaskResultById(ctx, resultId)
	if err != nil {
		return err
	}

	if taskResult.Result != "" {
		resultMap := make([]map[string]string, 0)
		if err := json.Unmarshal([]byte(taskResult.Result), &resultMap); err != nil {
			return err
		}
		res.Data = resultMap
	}

	return nil
}

// taskTargetMapNodeDetailVul 路径节点详情 - 漏洞节点
func (a *TaskCheckTask) taskTargetMapNodeDetailVul(ctx context.Context, resultId int, res *typespec.TaskTargetMapNodeDetailRes) error {
	// 获取结果数据
	var taskResultSrv services.TaskResult
	taskResult, err := taskResultSrv.GetTaskResultById(ctx, resultId)
	if err != nil {
		return err
	}

	var taskVulSrv services.TaskVul
	taskVul, _ := taskVulSrv.GetByTargetResultId(ctx, taskResult.ID)

	status := 1
	name := "--"
	link := "--"
	location := "--"
	result := "--"
	fix := "--"
	desc := "--"
	risk := 0
	typee := 0
	if taskVul.ID != 0 {
		status = taskVul.Status
		if taskVul.RefUrl != "" {
			link = taskVul.RefUrl
		}
		if taskVul.VulAddress != "" {
			location = taskVul.VulAddress
		}
		if taskVul.VulResult != "" {
			result = taskVul.VulResult
		}
		if taskVul.FixSuggest != "" {
			fix = taskVul.FixSuggest
		}
		if taskVul.Description != "" {
			desc = taskVul.Description
		}
		if taskVul.Risk != 0 {
			risk = taskVul.Risk
		}
		if taskVul.Type != 0 {
			typee = taskVul.Type
		}
		if taskVul.Name != "" {
			name = taskVul.Name
		}
	} else {
		// 未查询到具体的执行结果，直接使用任务结果中的数据
		returnData := map[string]interface{}{
			"title": taskResult.VulName,
		}
		res.Type = "infoCollect"
		res.Data = returnData
		return nil
	}

	var vulNode typespec.TaskTargetMapNodeDetailVul
	vulNode.Title = name
	vulNode.Status = status
	vulNode.StatusEnum = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(status)
	vulNode.Type = typee
	vulNode.TypeEnum = enums.ToolsVulnerabilityEnum.GetTypeEnum(typee)
	vulNode.Risk = risk
	vulNode.RiskEnum = enums.ToolsVulnerabilityEnum.GetRiskEnum(risk)
	vulNode.Description = desc
	vulNode.FixSuggest = fix
	vulNode.Link = link
	vulNode.Location = location
	vulNode.Result = result

	res.Type = "vul"
	res.Data = vulNode
	return nil
}

// taskTargetMapNodeDetailPrinciple 路径节点详情 - 组件扫描节点
func (a *TaskCheckTask) taskTargetMapNodeDetailPrinciple(ctx context.Context, resultId int, res *typespec.TaskTargetMapNodeDetailRes) error {
	// 获取结果数据
	var taskResultSrv services.TaskResult
	taskResult, err := taskResultSrv.GetTaskResultById(ctx, resultId)
	if err != nil {
		return err
	}
	// 组织返回数据
	returnData := map[string]interface{}{
		"title": taskResult.VulName,
	}
	res.Type = "principle"
	res.Data = returnData

	return nil
}

// taskTargetMapNodeDetailWhois 路径节点详情 - whois节点
func (a *TaskCheckTask) taskTargetMapNodeDetailWhois(ctx context.Context, resultId int, res *typespec.TaskTargetMapNodeDetailRes) error {
	// 获取结果数据
	var taskResultSrv services.TaskResult
	taskResult, err := taskResultSrv.GetTaskResultById(ctx, resultId)
	if err != nil {
		return err
	}

	// 处理结果数据
	dataList := make([]map[string]string, 0)
	if taskResult.Result != "" {
		resultData := make(map[string]string, 0)
		if err = json.Unmarshal([]byte(taskResult.Result), &resultData); err != nil {
			return err
		}
		dataListJson := resultData["dataList"]
		if err = json.Unmarshal([]byte(dataListJson), &dataList); err != nil {
			return err
		}
	}

	// 组织返回数据
	returnData := map[string]interface{}{
		"title":    "whois",
		"dataList": dataList,
	}
	res.Type = "whois"
	res.Data = returnData

	return nil
}

// TaskThreeExport 三方数据导出
func (a *TaskCheckTask) TaskThreeExport(ctx context.Context, req *typespec.TaskThreeExportReq, res *typespec.TaskThreeExportRes) error {
	// 获取本机地址
	localIps, err := network.GetLocalIp()
	if err != nil {
		return err
	}
	var localIp string
	if len(localIps) > 0 {
		localIp = localIps[0]
	}

	// 获取任务信息
	var taskSrv services.TaskTask
	task, err := taskSrv.GetTaskByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}

	// 获取目标信息
	var taskTargetSrv services.TaskTarget
	targets, _ := taskTargetSrv.GetTargetByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}

	// 获取端口情况
	var taskTaskResultSrv services.TaskTaskResult
	taskPorts := taskTaskResultSrv.GetByTaskIdAndSubObjType(ctx, req.TaskId, enums.TaskResultSubObjTypeService)

	// 获取漏洞信息
	var taskVulSrv services.TaskVul
	tasks := taskVulSrv.GetsByTaskId(ctx, req.TaskId, enums.VulDataTypOne)

	// 组合系统信息
	res.VulnFound.DataSource = typespec.DataThreeSource{
		Name:        task.TaskName,
		ProductName: "自动化渗透测试平台",
		VendorName:  "", // 不知道是啥
		VersionInfo: "v1.3.0",
		CreateTime:  time.Now().Format(utils.DateTime),
	}

	// 获取所有目标信息
	targetUrls := make([]string, 0)
	// 目标信息
	var tempHostList []typespec.ThreeHostList
	for _, item := range targets {
		targetUrls = append(targetUrls, item.TargetURL)

		// 目标信息 - 端口
		var tempPors []typespec.ThreePortList
		for _, port := range taskPorts {
			if port.SubObjID == strconv.Itoa(item.ID) {
				finalPortInt, _ := strconv.Atoi(port.Field2)
				tempPors = append(tempPors, typespec.ThreePortList{
					Portno:       finalPortInt,
					Prototype:    6,
					ProtocolName: port.Field3,
				})
			}
		}
		tempHostList = append(tempHostList, typespec.ThreeHostList{
			Hostip:   item.TargetURL,
			Ostype:   item.OpSys,
			PortList: tempPors,
		})
	}

	// 组合任务信息
	res.VulnFound.ScanTaskInfo = typespec.ScanThreeTaskInfo{
		TargetIPS: strings.Join(targetUrls, ","),
		BeginTime: task.CreateTime.Format(utils.DateTime),
		EndTime:   task.UpdateTime.Format(utils.DateTime),
		ScannerIP: localIp,
	}

	// 组合目标信息
	res.VulnFound.HostList = tempHostList

	// 组合漏洞信息
	for _, item := range tasks {
		// 获取 cve,cnvd,cnnvd
		cveCnvdCnnvd := strings.Split(item.VulID, ",")
		var cve, cnvd, cnnvd string
		for k, v := range cveCnvdCnnvd {
			switch k {
			case 0:
				cve = v
			case 1:
				cnvd = v
			case 2:
				cnnvd = v
			}
		}

		// 获取端口
		vulResult := make(map[string]string)
		json.Unmarshal([]byte(item.VulResult), &vulResult)

		res.VulnFound.VulnList = append(res.VulnFound.VulnList, typespec.ThreeVulnList{
			VulnId:                strconv.Itoa(item.ID),
			VulnName:              item.Name,
			ShortDesc:             item.Description,
			FullDesc:              item.Description,
			RiskLevelValue:        enums.GetTargetRisk(item.Risk),
			Platforms:             item.AffectRange,
			AssetIp:               item.TargetUrl,
			RepairAdvice:          item.FixSuggest,
			CncveTag:              "", // CNCVE编码
			CveTag:                cve,
			CnnvdTag:              cnnvd,
			CnvdTag:               cnvd,
			CvssScore:             item.Cvss,
			VulnDisclosureTime:    item.PublishedTime,
			VulnDisclosureTimeStr: item.PublishedTime,
			SecurityValue:         "", // 是否安全扫描
			PrincipleValue:        "", // 是否为原理扫描
			Protocol:              "", // 漏洞协议
			Port:                  vulResult["port"],
			VptValue:              "", // vpt值
			StatusValue:           enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(item.Status),
			Payload:               item.VerMsg,
			OriginalRequest:       item.VulParam,
		})
	}
	return nil
}

// 拓扑图
func (a *TaskCheckTask) TaskTopologyMap(ctx context.Context, req *typespec.TaskTopologyMapReq, res *typespec.TaskTopologyMapRes) error {
	// 根节点为任务名称
	// 2层节点为目标
	// 3层节点为横向的结果，如果有的话
	// 4层节点为横向的横向结果，如果有话
	// 5与4逻辑一样

	// 第一层
	// 获取任务数据
	var taskSrv services.TaskTask
	task, err := taskSrv.GetTaskByTaskId(ctx, req.TaskId)
	if err != nil {
		return err
	}

	root := typespec.TaskTargetMapState{
		Id:       "0",
		Label:    task.TaskName,
		Status:   "end",
		Success:  true,
		Hide:     false,
		Pocname:  "",
		Class:    "info",
		DetailId: "root_0",
	}
	res.State = append(res.State, root)

	// 第二层 【中间包含第3，4，...层，依据extend_field.hengxiang_pid 关联】
	var targetSrv services.TaskTarget
	target, _ := targetSrv.GetTargetByTaskId(ctx, req.TaskId)
	targetIds := make([]int, 0)
	for _, item := range target {
		targetIds = append(targetIds, item.ID)
	}

	// 任务状态
	res.Status = task.Status
	res.StatusEnum = enums.TaskTaskEnum.StatusEnum(task.Status)

	var taskVulSrv services.TaskVul
	_, riskLevleMap, _, _ := taskVulSrv.GetTargetStatsBytargetIds(ctx, targetIds)

	for _, item := range target {
		// 计算 Class
		class := "info"
		switch riskLevleMap[item.ID] {
		case enums.TargetRiskHigh:
			class = "high"
		case enums.TargetRiskMid:
			class = "middle"
		case enums.TargetRiskLow:
			class = "low"
		case enums.TargetRiskLowNoFound:
			class = "info"
		}

		targetId := strconv.Itoa(item.ID)
		res.State = append(res.State, typespec.TaskTargetMapState{
			Id:       targetId,
			Label:    item.TargetURL,
			Status:   "end",
			Success:  true,
			Hide:     false,
			Pocname:  "",
			Class:    class,
			DetailId: "target_" + targetId,
		})

		extendField := make(map[string]interface{})
		json.Unmarshal([]byte(item.ExtendField), &extendField)
		if pidInterface, ok := extendField["hengxiang_pid"]; ok {
			res.Edg = append(res.Edg, typespec.TaskTargetMapEdg{
				Start: strconv.Itoa(int(pidInterface.(float64))),
				End:   targetId,
				Hide:  false,
			})
		} else {
			res.Edg = append(res.Edg, typespec.TaskTargetMapEdg{
				Start: "0",
				End:   targetId,
				Hide:  false,
			})
		}

		// 远控数量
		if item.IsRemoteSession == enums.TargetIsRemoteSessionY {
			res.RemoteControlNum += 1
		}

		// 目标运行状态
		switch item.Status {
		case enums.TargetStatusFinish:
			res.CheckTarget.Finish++
		case enums.TargetStatusRunning:
			res.CheckTarget.Running++
		case enums.TargetStatusToBegin:
			res.CheckTarget.Wait++
		}
		if item.IsAlive == enums.TargetIsAliveN {
			res.CheckTarget.Fail++
		}

		// 目标风险等级
		switch riskLevleMap[item.ID] {
		case enums.TargetRiskHigh:
			res.RiskTarget.High++
		case enums.TargetRiskMid:
			res.RiskTarget.Mid++
		case enums.TargetRiskLow:
			res.RiskTarget.Low++
		case enums.TargetRiskLowNoFound:
			res.RiskTarget.Safe++
		}
	}

	return nil
}

func (a *TaskCheckTask) TaskTopologyMapNodeCheckParam(detailId string) (int, error) {
	detailIds := strings.Split(detailId, "_")
	if len(detailIds) != 2 {
		return 0, errors.New("detailId参数错误")
	}
	targetId, err := strconv.Atoi(detailIds[1])
	if err != nil {
		return 0, errors.New("目标ID必须是number")
	}
	if detailIds[0] != "target" {
		return 0, errors.New("非target节点不可点击")
	}
	return targetId, nil
}

// 拓扑图节点 - 详情
func (a *TaskCheckTask) TaskTopologyMapNodeDetail(ctx context.Context, req *typespec.TaskTopologyMapNodeDetailReq, res *typespec.TaskTopologyMapNodeDetailRes) error {
	targetId, err := a.TaskTopologyMapNodeCheckParam(req.DetailId)
	if err != nil {
		return err
	}

	var targetSrv services.TaskTarget
	target, err := targetSrv.GetTargetById(ctx, targetId)
	if err != nil {
		return err
	}

	var taskVulSrv services.TaskVul
	_, riskLevleMap, _, _ := taskVulSrv.GetTargetStatsBytargetIds(ctx, []int{targetId})

	// 目标地址
	res.TargetUrl = target.TargetURL
	// 风险等级
	res.RiskEnum = enums.GetTargetRisk(riskLevleMap[targetId])
	// 目标验证状态
	res.VerifyEnum = taskVulSrv.GetTargetVerifyStatus(ctx, []int{targetId})
	// 主机评分
	res.Fen = strconv.FormatFloat(float64(target.UseScore)/100, 'f', 2, 64)
	// 分层数据，取fingerprint脚本的返回结果，依据里面的level字段获取
	levelApp := make([]string, 0)
	levelSupport := make([]string, 0)
	levelService := make([]string, 0)
	levelSystem := make([]string, 0)
	levelHardware := make([]string, 0)
	var taskResult services.TaskResult
	fingerList := taskResult.AllTaskResultByTargetId(ctx, targetId, enums.ScriptNameFingerPrint)
	for _, item := range fingerList {
		resultMap := make(map[string]string)
		json.Unmarshal([]byte(item.Result), &resultMap)
		if appInfo, ok := resultMap["app_info"]; ok {
			var appInfoMap services.FingerPrintResult
			json.Unmarshal([]byte(appInfo), &appInfoMap)
			for _, info := range appInfoMap.TargetInfo {
				// 5 应用层 表示目标资产具体运行的业务应用，例如:
				// 4 支撑层 表示目标资产的开发语言/框架
				// 3 服务层 表示目标资产运行的服务
				// 2 系统层 表示目标资产运行的系统
				// 1 硬件层 表示目标资产运行的硬件设备或安防设备
				// 9 默认为9 表示未分层指纹
				level, _ := strconv.Atoi(info.Level)
				switch level {
				case enums.FingerLevelApp:
					levelApp = append(levelApp, info.AppName)
				case enums.FingerLevelSupport:
					levelSupport = append(levelSupport, info.AppName)
				case enums.FingerLevelService:
					levelService = append(levelService, info.AppName)
				case enums.FingerLevelSystem:
					levelSystem = append(levelSystem, info.AppName)
				case enums.FingerLevelHardware:
					levelHardware = append(levelHardware, info.AppName)
				}
			}
		}
	}
	res.ApplicationLevel = strings.Join(levelApp, "、")
	res.SupportLevel = strings.Join(levelSupport, "、")
	res.ServiceLevel = strings.Join(levelService, "、")
	res.SystemLevel = strings.Join(levelSystem, "、")
	res.HardwareLevel = strings.Join(levelHardware, "、")
	return nil
}

// 获取平台下发现的所有漏洞信息
func (a *TaskCheckTask) AllTaskVulByPage(ctx context.Context, req *typespec.TaskAllVulByPageReq, res *typespec.TaskAllVulByPageRes) error {
	var srv services.TaskVul
	list, total := srv.AllByPage(ctx, req.Page, req.Size)
	res.Total = total

	var taskVulSrv services.TaskVul
	domainMap := make(map[string]string, 0)
	for _, item := range list {
		var tmp typespec.ApiVulListRespItems
		tmp.Id = item.ID
		tmp.TargetUrl = item.TargetUrl
		tmp.Name = item.Name
		tmp.Type = item.Type
		tmp.TypeName = enums.ToolsVulnerabilityEnum.GetTypeEnum(item.Type)
		tmp.Risk = item.Risk
		tmp.RiskName = enums.ToolsVulnerabilityEnum.GetRiskEnum(item.Risk)
		tmp.Location = item.Location
		tmp.Status = item.Status
		tmp.StatusName = enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(item.Status)

		tmp.TaskId = item.TaskID
		tmp.TargetId = item.TargetID
		tmp.Pocname = item.Pocname
		tmp.Cve = strings.Split(item.VulID, ",")[0]
		tmp.Cvss = item.Cvss
		tmp.PublishedTime = item.PublishedTime
		tmp.ExploitImpact = item.ExploitImpact
		tmp.Description = item.Description
		tmp.FixSuggest = item.FixSuggest
		tmp.RefUrl = item.RefUrl
		tmp.VulAddress = item.VulAddress
		tmp.VulResult = item.VulResult
		tmp.VulParam = item.VulParam
		tmp.VerMsg = item.VerMsg
		tmp.Class = item.Class
		tmp.Ip, tmp.Port = taskVulSrv.GetParamIpPort(item.VulParam, domainMap)
		if tmp.Location == "" {
			tmp.Location = taskVulSrv.GetParamLocation(item.VulParam)
		}

		res.List = append(res.List, tmp)
	}
	return nil
}

// CheckVulInfo 漏洞检测
func (a *TaskCheckTask) CheckVulInfo(ctx context.Context, req *typespec.VulInfoReq, resp *typespec.CheckVulResp) error {
	var (
		taskVulSrv services.TaskVul
		tmpVerMsg  = make([]typespec.VulInfoRespVerMsg, 0)
		isSuccess  bool
		content    string
		aesEcb     aesEncryption.AesEcb
	)
	taskVulRes, err := taskVulSrv.VulInfo(ctx, req.TaskVulId)
	if err != nil {
		return err
	}
	if taskVulRes.ID == 0 {
		return errors.New("找不到该漏洞！")
	}
	// 解密获取请求报文
	verMsgDecodeByte, err := hex.DecodeString(taskVulRes.VerMsg)
	json.Unmarshal(aesEcb.AesDecryptECB(verMsgDecodeByte, []byte("9876787656785679")), &tmpVerMsg)
	if len(tmpVerMsg) >= 1 {
		isSuccess, content, err = httpclients.VulInfoCheck(tmpVerMsg[0].Request, taskVulRes.Description)
		if err != nil {
			return err
		}
	}
	resp.FalseAlarm = isSuccess
	resp.Content = content
	return nil
}
