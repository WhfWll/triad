package application

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/services"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"
)

type FlowApp struct{}

// FlowTaskEnum 流量分析枚举
func (f *FlowApp) FlowTaskEnum(ctx context.Context, resp *typespec.FlowTaskEnumResp) error {
	var (
		flowEnums enums.Flow
		flowSrv   services.Flow
		logicSrv  services.Logic
		err       error
	)
	resp.Status = flowEnums.GetFlowTaskStatusEnumArray()
	resp.ExpireTime = flowEnums.GetFlowTaskExpireTimeEnumArray()
	resp.NetworkCard, err = flowSrv.FlowTaskNetworkCardEnums(ctx)
	if err != nil {
		return err
	}
	resp.VulRiskLevel = flowEnums.GetFlowRiskLevelEnumArray()
	resp.CredPattern = logicSrv.CredPatternEnum()
	resp.FuzzParam = struct {
		Character string `json:"character"` //字符型
		Number    string `json:"number"`    //数字型
	}{Character: enums.Character, Number: enums.Number}
	resp.Response.JsonKeyword = enums.JsonKeyword
	resp.Response.NoJsonSwitch = true
	resp.Response.NoJsonKeyword = enums.NoJsonKeyword
	resp.VulName = flowSrv.FlowGetVulNameMap(ctx)
	return nil
}

// FlowTaskList 流量分析任务列表
func (f *FlowApp) FlowTaskList(ctx context.Context, req *typespec.FlowTaskListReq, resp *typespec.FlowTaskListResp) error {
	var (
		flowSrv   services.Flow
		flowenums enums.Flow
		flowEnums enums.Flow
		ids       = make([]int, 0)
	)
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
	//查询数据
	flowTaskRes, total, err := flowSrv.FlowTaskList(ctx, req.Search, req.Page, req.Size, userIdList)
	if err != nil {
		return err
	}
	resp.Total = total
	//统计风险等级
	for i := 0; i < len(flowTaskRes); i++ {
		ids = append(ids, flowTaskRes[i].ID)
	}
	vulMap, risklevelMap, err := flowSrv.FlowRiskStataByTaskIds(ctx, ids)
	//组装返回结果
	for i := 0; i < len(flowTaskRes); i++ {
		var (
			tmp      typespec.FlowTaskListRespItems
			riskl    int
			vulArray = make([]int, 5)
		)

		if v, ok := vulMap[flowTaskRes[i].ID]; ok {
			vulArray = v
		}
		if v, ok := risklevelMap[flowTaskRes[i].ID]; ok {
			riskl = v
		}
		tmp.Id = flowTaskRes[i].ID
		tmp.TaskName = flowTaskRes[i].TaskName
		tmp.RiskLevel = riskl
		tmp.RiskLevelName = flowEnums.FlowRiskLevelToTargetRiskLevelName(riskl)
		tmp.Status = flowTaskRes[i].Status
		tmp.StatusName = flowenums.FlowTaskStatusEnum(flowTaskRes[i].Status)
		tmp.VulNum = vulArray
		tmp.UpdateTime = flowTaskRes[i].UpdateTime.Format(enums.TimeLayout)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// FlowTaskDel 流量分析任务删除
func (f *FlowApp) FlowTaskDel(ctx context.Context, req *typespec.FlowTaskDelReq) error {
	var flowSrv services.Flow
	//查询任务过滤任务,只删除结束的任务
	ids := strings.Split(req.FlowTaskIds, ",")
	newIds, err := flowSrv.GetFlowtaskListByIdsAndStatus(ctx, enums.FlowTaskStatusDone, ids)
	if err != nil {
		return err
	}
	if len(newIds) != len(ids) {
		return errors.New("只有结束的任务才可以被删除,请重新选择...")
	}
	//执行删除
	err = flowSrv.FlowTaskDel(ctx, newIds)
	if err != nil {
		return err
	}
	return nil
}

// FlowTaskAdd 创建流量分析任务
func (f *FlowApp) FlowTaskAdd(ctx context.Context, req *typespec.FlowTaskAddReq) error {
	var flowSrv services.Flow
	//查询节点状态
	nodeRes, err := flowSrv.GetNodeInfo("node3")
	if err != nil {
		return err
	}
	if len(nodeRes.FlowTaskID) != 0 || len(nodeRes.NodeStatus) != 0 {
		return errors.New("节点正在被其他任务占用或节点IP不可用!")
	}
	//保存数据
	targetUrlArray := strings.Split(req.TargetUrl, ",")
	err = flowSrv.FlowTaskSave(ctx, req.TaskName, req.Port, targetUrlArray, req.NetworkCard, req.ExpireTime, req.UserId, req.OtherConfig, req.VulConfig)
	if err != nil {
		return err
	}
	return nil
}

// ChangeFlowTaskStatus 结束流量分析任务
func (f *FlowApp) ChangeFlowTaskStatus(ctx context.Context, req *typespec.ChangeFlowTaskStatusReq) error {
	var flowSrv services.Flow
	//查询数据
	flowTaskRes, err := flowSrv.FlowTaskInfo(ctx, req.FlowTaskId)
	if err != nil || flowTaskRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	//判断状态是否在运行中
	if flowTaskRes.Status != enums.FlowTaskStatusExec {
		return errors.New("状态非运行中")
	}
	//修改状态为已结束
	err = flowSrv.UpdateFlowTask(ctx, req.FlowTaskId, map[string]interface{}{"status": enums.FlowTaskStatusDone})
	if err != nil {
		return err
	}
	//请求决策引擎
	//flowTargetData := make([]mysqls.Flowtarget, 0)
	err = flowSrv.SendDecisionOprate("", req.Operate, req.FlowTaskId)
	if err != nil {
		return err
	}
	return nil
}

// FlowTaskInfo 流量分析详情
func (f *FlowApp) FlowTaskInfo(ctx context.Context, req *typespec.FlowTaskInfoReq, resp *typespec.FlowTaskInfoResp) error {
	var (
		flowSrv   services.Flow
		userSrv   services.User
		flowEnums enums.Flow
	)
	//查询任务数据
	flowTaskRes, err := flowSrv.FlowTaskInfo(ctx, req.FlowTaskId)
	if err != nil {
		return err
	}
	if flowTaskRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	//查询目标数据
	targertIds, targetUrls, err := flowSrv.FlowTargetInfoByTaskId(ctx, req.FlowTaskId)
	if err != nil {
		return err
	}
	//查询漏洞数量
	vulNum, targetNum, riskLevel, err := flowSrv.FlowRiskStata(ctx, req.FlowTaskId, targertIds)
	if err != nil {
		return err
	}
	//查询用户数据
	userRes, err := userSrv.GetUserDetail(ctx, flowTaskRes.UserID)
	if err != nil {
		return err
	}
	resp.Id = flowTaskRes.ID
	resp.TaskName = flowTaskRes.TaskName
	resp.NetworkCard = flowTaskRes.NetwordCard
	resp.Port = flowTaskRes.Port
	resp.CreateTime = flowTaskRes.CreateTime.Format(enums.TimeLayout)
	resp.UserId = flowTaskRes.UserID
	resp.UserName = userRes.Username
	resp.ExpireTime = flowTaskRes.ExpireTime
	resp.ExpireTimeName = flowEnums.FlowTaskExpireTimeEnum(flowTaskRes.ExpireTime)
	resp.RiskLevel = riskLevel
	resp.RiskLevelName = flowEnums.FlowRiskLevelToTargetRiskLevelName(riskLevel)
	resp.Status = flowTaskRes.Status
	resp.StatusName = flowEnums.FlowTaskStatusEnum(flowTaskRes.Status)
	resp.Target = strings.Join(targetUrls, ",")
	resp.VulNum = vulNum
	resp.TargetNum = targetNum
	var otherConfig typespec.OtherConfig
	err = json.Unmarshal([]byte(flowTaskRes.OtherConfig), &otherConfig)
	resp.OtherConfig = otherConfig

	// 获取被动流量任务的漏洞配置信息
	param := map[string]interface{}{
		"scriptType": "mitm",
		"page":       "1",
		"size":       "1000",
	}
	res, _ := httpclients.GetVulNameByType(param)
	vulConfigZh := ""
	for _, vul := range res.List {
		for _, pocname := range strings.Split(flowTaskRes.VulConfig, ",") {
			if strings.TrimSpace(pocname) == "" {
				continue
			}
			if vul.Pocname == pocname {
				vulConfigZh += vul.VulName + ","
			}
		}
	}
	resp.VulConfig = flowTaskRes.VulConfig
	resp.VulConfigZh = strings.Trim(vulConfigZh, ",")
	return nil
}

// HttpsCert 证书下载
func (f *FlowApp) HttpsCert(ctx context.Context, req *typespec.HttpsCertReq, resp *typespec.HttpsCertResp) error {
	var flowSrv services.Flow
	//只有运行中的被动流量才可以下载证书
	flowTaskRes, err := flowSrv.FlowTaskInfo(ctx, req.FlowTaskId)
	if err != nil {
		return err
	}
	if flowTaskRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	if flowTaskRes.Status != enums.FlowTaskStatusExec {
		return errors.New("请确认任务是否正在运行中...")
	}
	cert, err := httpclients.DownloadMITMCert()
	if err != nil {
		return err
	}
	resp.Cert = cert
	return nil
}

// FlowTaskStatus 任务状态
func (f *FlowApp) FlowTaskStatus(ctx context.Context, req *typespec.FlowTaskStatusReq, resp *typespec.FlowTaskStatusResp) error {
	var (
		flowSrv   services.Flow
		flowEnums enums.Flow
	)
	//查询任务数据
	flowTaskRes, err := flowSrv.FlowTaskInfo(ctx, req.FlowTaskId)
	if err != nil {
		return err
	}
	if flowTaskRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	resp.Id = flowTaskRes.ID
	resp.Status = flowTaskRes.Status
	resp.StatusName = flowEnums.FlowTaskStatusEnum(flowTaskRes.Status)
	return nil
}

// FlowRiskList 漏洞列表
func (f *FlowApp) FlowRiskList(ctx context.Context, req *typespec.FlowRiskListReq, resp *typespec.FlowRiskListResp) error {
	var (
		flowSrv   services.Flow
		flowEnums enums.Flow
	)

	flowRiskRes, total, err := flowSrv.FlowRiskList(ctx, req.Search, req.RiskLevel, req.FlowTaskId, req.Page, req.Size)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(flowRiskRes); i++ {
		var tmp typespec.FlowRiskListRespItems
		tmp.Id = flowRiskRes[i].ID
		tmp.Title = flowRiskRes[i].Title
		tmp.Ip = flowRiskRes[i].IP
		tmp.Host = flowRiskRes[i].Host
		tmp.RiskLevel = flowRiskRes[i].RiskLevel
		tmp.RiskLevelName = flowEnums.FlowRiskLevelEnum(flowRiskRes[i].RiskLevel)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// FlowRiskInfo 漏洞详情
func (f *FlowApp) FlowRiskInfo(ctx context.Context, req *typespec.FlowRiskInfoReq, resp *typespec.FlowRiskInfoResp) error {
	var (
		flowSrv   services.Flow
		flowEnums enums.Flow
	)
	flowRiskRes, err := flowSrv.FlowRiskInfo(ctx, req.FlowRiskId)
	if err != nil {
		return err
	}
	if flowRiskRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	resp.Id = flowRiskRes.ID
	resp.Ip = flowRiskRes.IP
	resp.Port = flowRiskRes.Port
	resp.Host = flowRiskRes.Host
	resp.RiskLevel = flowRiskRes.RiskLevel
	resp.RiskLevelName = flowEnums.FlowRiskLevelEnum(flowRiskRes.RiskLevel)
	resp.CreateTime = flowRiskRes.CreateTime.Format(enums.TimeLayout)
	resp.UpdateTime = flowRiskRes.UpdateTime.Format(enums.TimeLayout)
	resp.Hash = flowRiskRes.Hash
	resp.RiskTypeVerbose = flowRiskRes.RiskTypeVerbose
	resp.Origin = "漏洞检测"
	resp.Status = "已验证"
	resp.Url = flowRiskRes.Host
	resp.Parameter = flowRiskRes.Parameter
	decodedName, err := url.QueryUnescape(flowRiskRes.Request)
	resp.Request = decodedName
	resp.Response = flowRiskRes.Response
	resp.Detail = flowRiskRes.Detail
	resp.Token = "-"
	// 获取 success_payload_flag
	var tempVulMap map[string]interface{}
	err = json.Unmarshal([]byte(flowRiskRes.Detail), &tempVulMap)
	if tempVulMap["payload_success_flag"] != nil {
		resp.PayloadSuccessFlag = tempVulMap["payload_success_flag"].(string)
		delete(tempVulMap, "payload_success_flag")
	}
	if tempVulMap["payload"] != nil {
		resp.Payload = tempVulMap["payload"].(string)
		delete(tempVulMap, "payload")
	}
	if tempVulMap["request"] != nil {
		delete(tempVulMap, "request")
	}
	if tempVulMap["response"] != nil {
		delete(tempVulMap, "response")
	}
	detailByte, err := json.Marshal(tempVulMap)
	if err == nil {
		resp.Detail = string(detailByte)
	}
	return nil
}

// FlowRiskDel 漏洞删除
func (f *FlowApp) FlowRiskDel(ctx context.Context, req *typespec.FlowRiskDelReq) error {
	idArray := strings.Split(req.FlowRiskIds, ",")
	var flowSrv services.Flow
	return flowSrv.FlowRiskDel(ctx, idArray)
}

// FlowBaseList 被动流量列表
func (f *FlowApp) FlowBaseList(ctx context.Context, req *typespec.FlowBaseListReq, resp *typespec.FlowBaseListResp) error {
	var flowSrv services.Flow
	flowBaseRes, total, err := flowSrv.FlowBaseList(ctx, req.Search, req.FlowTaskId, req.Page, req.Size)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(flowBaseRes); i++ {
		var tmp typespec.FlowBaseListRespItems
		tmp.Id = flowBaseRes[i].ID
		tmp.Ip = flowBaseRes[i].IP
		tmp.Method = flowBaseRes[i].Method
		tmp.CreateTime = flowBaseRes[i].CreateTime.Format(enums.TimeLayout)
		tmp.RespContentType = flowBaseRes[i].RespContentType
		tmp.RespCode = flowBaseRes[i].RespCode
		tmp.RespTitle = flowBaseRes[i].RespTitle
		tmp.Tags = flowBaseRes[i].Tags
		tmp.Host = flowBaseRes[i].Host
		tmp.Url = flowBaseRes[i].Url
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// FlowBaseInfo 被动流量详情
func (f *FlowApp) FlowBaseInfo(ctx context.Context, req *typespec.FlowBaseInfoReq, resp *typespec.FlowBaseInfoResp) error {
	var flowSrv services.Flow
	//查询数据
	flowBaseRes, err := flowSrv.FlowBaseInfo(ctx, req.FlowBaseId)
	if err != nil {
		return err
	}
	if flowBaseRes.ID == 0 {
		return errors.New("找不到数据！")
	}
	resp.ReqHeader = flowBaseRes.ReqHeader
	resp.RespHeader = flowBaseRes.RespHeader
	return nil
}

// FlowBaseDel 被动流量删除
func (f *FlowApp) FlowBaseDel(ctx context.Context, req *typespec.FlowBaseDelReq) error {
	idArray := strings.Split(req.FlowBaseIds, ",")
	var flowSrv services.Flow
	return flowSrv.FlowBaseDel(ctx, idArray)
}

// FlowLogInfo 被动流量详情
func (f *FlowApp) FlowLogInfo(ctx context.Context, req *typespec.FlowLogInfoReq, resp *typespec.FlowLogInfoResp) error {
	var flowSrv services.Flow
	//查询数据
	flowLogRes, total, err := flowSrv.FlowLogInfo(ctx, req.Search, req.FlowTaskId, req.Page, req.Size)
	if err != nil {
		return err
	}
	resp.Total = total
	//返回结果
	for i := 0; i < len(flowLogRes); i++ {
		var tmp typespec.FlowLogInfoRespItems
		tmp.Id = flowLogRes[i].ID
		tmp.Content = flowLogRes[i].Content
		tmp.CreateTime = flowLogRes[i].CreateTime.Format(enums.TimeLayout)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// FlowLogDel 被动流量日志清除
func (f *FlowApp) FlowLogDel(ctx context.Context, req *typespec.FlowLogDelReq) error {
	var flowSrv services.Flow
	return flowSrv.FlowLogDel(ctx, req.FlowTaskId)
}

// FlowTaskEdit 流量分析任务信息编辑
func (f *FlowApp) FlowTaskEdit(ctx context.Context, req *typespec.FlowTaskEditReq) error {
	var flowSrv services.Flow
	//查询节点状态
	nodeRes, err := flowSrv.GetNodeInfo("node3")
	if err != nil {
		return err
	}
	if nodeRes.FlowTaskID == strconv.Itoa(req.FlowTaskId) {
		err = flowSrv.SendDecisionOprate("", "stop", req.FlowTaskId)
		if err != nil {
			return err
		}
	} else if len(nodeRes.FlowTaskID) != 0 || len(nodeRes.NodeStatus) != 0 {
		return errors.New("节点正在被其他任务占用或节点IP不可用!")
	}
	time.Sleep(1500 * time.Millisecond)
	targetUrlArray := strings.Split(req.TargetUrl, ",")
	_, err = flowSrv.FlowTaskEdit(ctx, req.FlowTaskId, req.TaskName, req.Port, targetUrlArray, req.NetworkCard, req.ExpireTime, req.UserId, req.OtherConfig, req.VulConfig)
	if err != nil {
		return err
	}
	err = flowSrv.UpdateFlowTask(ctx, req.FlowTaskId, map[string]interface{}{"status": enums.FlowTaskStatusExec})
	return nil
}

// FlowTaskExport 任务流量信息导出
func (f *FlowApp) FlowTaskExport(ctx context.Context, req *typespec.FlowTaskExportReq, resp *typespec.FlowTaskExportResp) error {
	var srv services.Flow
	flowBaseRes, _, err := srv.FlowBaseList(ctx, "", req.TaskId, 1, 100000)
	if err != nil {
		return err
	}
	for i := 0; i < len(flowBaseRes); i++ {
		var tmp typespec.FlowBaseListRespItems
		tmp.Id = flowBaseRes[i].ID
		tmp.Ip = flowBaseRes[i].IP
		tmp.Method = flowBaseRes[i].Method
		tmp.CreateTime = flowBaseRes[i].CreateTime.Format(enums.TimeLayout)
		tmp.RespContentType = flowBaseRes[i].RespContentType
		tmp.RespCode = flowBaseRes[i].RespCode
		tmp.RespTitle = flowBaseRes[i].RespTitle
		tmp.Tags = flowBaseRes[i].Tags
		tmp.Host = flowBaseRes[i].Host
		tmp.Url = flowBaseRes[i].Url
		resp.List = append(resp.List, tmp)
	}
	return nil
}
