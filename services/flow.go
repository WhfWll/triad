package services

import (
	"context"
	"encoding/json"
	"errors"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/client/httpclients"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strings"
	"time"
)

type Flow struct{}

// FlowTaskList 流量分析任务列表
func (f *Flow) FlowTaskNetworkCardEnums(ctx context.Context) (interface{}, error) {
	var mapSet mysqls.MapSet
	mapSetRes, err := mapSet.GetsByObjKey(ctx, enums.FlowTaskNetworkCardMapSetKey)
	if err != nil {
		return nil, err
	}
	if len(mapSetRes.ObjValue) == 0 {
		return nil, errors.New("找不到网卡配置枚举")
	}
	var tmp []interface{}
	err = json.Unmarshal([]byte(mapSetRes.ObjValue), &tmp)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

// FlowTaskList 流量分析任务列表
func (f *Flow) FlowTaskList(ctx context.Context, search string, page, size int, userIdList []int) ([]mysqls.Flowtask, int64, error) {
	var flowMysqls mysqls.Flowtask
	return flowMysqls.GetFlowtaskList(ctx, search, page, size, userIdList)
}

// FlowTaskDel 查询是结束状态的任务id
func (f *Flow) GetFlowtaskListByIdsAndStatus(ctx context.Context, status int, ids any) ([]int, error) {
	var (
		flowMysqls mysqls.Flowtask
		result     = make([]int, 0)
	)
	flowTaskRes, err := flowMysqls.GetFlowtaskListByIdsAndStatus(ctx, status, ids)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(flowTaskRes); i++ {
		result = append(result, flowTaskRes[i].ID)
	}
	return result, nil
}

// FlowTaskDel 流量分析任务删除
func (f *Flow) FlowTaskDel(ctx context.Context, ids any) error {
	var (
		flowTaskMysqls   mysqls.Flowtask
		flowTargetMysqls mysqls.Flowtarget
		flowRiskMysqls   mysqls.Flowrisk
		flowBaseMysqls   mysqls.Flowbase
		flowLogMysqls    mysqls.Flowlog
	)
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	//flow_task
	err := flowTaskMysqls.DeleteFlowtaskByIds(dCtx, ids)
	if err != nil {
		return err
	}
	//flow_target
	err = flowTargetMysqls.DeleteFlowtargetByTaskIds(dCtx, ids)
	if err != nil {
		return err
	}
	//flow_risk
	err = flowRiskMysqls.DeleteFlowriskByTaskIds(dCtx, ids)
	if err != nil {
		return err
	}
	//flow_base
	err = flowBaseMysqls.DeleteFlowBaseByTaskIds(dCtx, ids)
	if err != nil {
		return err
	}
	//flow_log
	err = flowLogMysqls.DeleteByFlowTaskIds(dCtx, ids)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// FlowTaskSave 流量分析任务创建
func (f *Flow) FlowTaskSave(ctx context.Context, taskName, port string, targetUrl []string, networkCard string, expireTime, userId int, otherConfig, vulConfig string) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	//任务表
	var flowTaskMysqls = mysqls.Flowtask{
		TaskName:    taskName,
		NodeID:      "",
		NetwordCard: networkCard,
		Port:        port,
		ExpireTime:  expireTime,
		Status:      enums.FlowTaskStatusExec,
		UserID:      userId,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		OtherConfig: otherConfig,
		VulConfig:   vulConfig,
	}
	err := flowTaskMysqls.AddFlowtask(dCtx)
	if err != nil {
		return err
	}
	//任务目标表
	var (
		flowTargetMysqls mysqls.Flowtarget
		flowTargetData   = make([]mysqls.Flowtarget, 0)
	)
	for i := 0; i < len(targetUrl); i++ {
		var tmpflowTargetData = mysqls.Flowtarget{
			FlowTaskID: flowTaskMysqls.ID,
			TargetURL:  targetUrl[i],
			Status:     enums.FlowTargetStatusTwo,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		flowTargetData = append(flowTargetData, tmpflowTargetData)
	}
	err = flowTargetMysqls.AddManyFlowtarget(dCtx, flowTargetData)
	if err != nil {
		return err
	}
	err = sendDecisionCreate(&flowTaskMysqls, flowTargetData, otherConfig, vulConfig)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// SendDecisionCreate 执行任务
func sendDecisionCreate(task *mysqls.Flowtask, target []mysqls.Flowtarget, otherConfig, vulConfig string) error {
	//构建请求参数
	var lisData = make([]httpclients.CreateFlowTaskReqItems, 0)
	for i := 0; i < len(target); i++ {
		var tmp = httpclients.CreateFlowTaskReqItems{
			Host:         target[i].TargetURL,
			FlowTaskId:   target[i].FlowTaskID,
			FlowTargetId: target[i].ID,
		}
		lisData = append(lisData, tmp)
	}
	var param = map[string]interface{}{
		"flowTaskId":  task.ID,
		"node":        "",
		"mitmHost":    task.NetwordCard,
		"port":        task.Port,
		"listenData":  lisData,
		"expireTime":  task.ExpireTime,
		"otherConfig": otherConfig,
		"vulConfig":   vulConfig,
	}
	//请求决策引擎启动服务
	err := httpclients.CreateFlowTask(param)
	if err != nil {
		return err
	}
	return nil
}

// FlowTaskInfo 流量分析任务详情
func (f *Flow) FlowTaskInfo(ctx context.Context, flowTaskId int) (mysqls.Flowtask, error) {
	var flowMysqls mysqls.Flowtask
	return flowMysqls.GetFlowtask(ctx, flowTaskId)
}

// UpdateFlowTask 修改流量分析任务
func (f *Flow) UpdateFlowTask(ctx context.Context, flowTaskId int, param map[string]interface{}) error {
	var flowMysqls mysqls.Flowtask
	return flowMysqls.UpdateFlowtask(ctx, flowTaskId, param)
}

// SendDecisionOprate 操作任务
func (f *Flow) SendDecisionOprate(node string, oprate string, flowTaskId int) error {
	var param = map[string]interface{}{
		"node":       node,
		"oprate":     oprate,
		"flowTaskId": flowTaskId,
	}
	err := httpclients.OprateFlowTask(param)
	if err != nil {
		return err
	}
	return nil
}

// GetNodeMinIp 获取node数据mainIpAddr
func (f *Flow) GetNodeMinIp(node string) (string, error) {
	var param = map[string]interface{}{
		"node": node,
	}
	res, err := httpclients.GetNodeInfo(param)
	if err != nil {
		return "", err
	}
	return res.MainIpAddr, nil
}

// GetNodeInfo 获取node数据mainIpAddr
func (f *Flow) GetNodeInfo(node string) (*httpclients.GetNodeInfoRespData, error) {
	var param = map[string]interface{}{
		"node": node,
	}
	return httpclients.GetNodeInfo(param)
}

// FlowTargetInfoByTaskId 获取流量目标id数组和目标数组
func (f *Flow) FlowTargetInfoByTaskId(ctx context.Context, flowTaskId int) (map[int]int, []string, error) {
	var (
		flowMysqls mysqls.Flowtarget
		ids        = make(map[int]int, 0)
		urls       = make([]string, 0)
	)
	flowTaskRes, err := flowMysqls.GetFlowtargetList(ctx, flowTaskId)
	if err != nil {
		return ids, urls, err
	}
	for i := 0; i < len(flowTaskRes); i++ {
		ids[flowTaskRes[i].ID] = 1
		urls = append(urls, flowTaskRes[i].TargetURL)
	}
	return ids, urls, nil
}

// FlowTargetInfoByTaskId 获取流量目标id数组和目标数组
func (f *Flow) FlowRiskStata(ctx context.Context, flowTaskId int, targetIds map[int]int) ([]int, []int, int, error) {
	var (
		flowRiskMysqls mysqls.Flowrisk
		riskNum        = make([]int, 5)
		riskNumMap     = make(map[int]int, 0)
		targetNum      = make([]int, 4)
		targetNumMap   = make(map[int]int, 0)
		riskLevel      = enums.FlowRiskLevelSafe
	)
	//初始化参数
	for k, _ := range targetIds {
		targetNumMap[k] = enums.FlowRiskLevelSafe
	}
	flowRiskRes, err := flowRiskMysqls.GetFlowriskListByTaskId(ctx, flowTaskId)
	if err != nil {
		return nil, nil, 0, err
	}
	for i := 0; i < len(flowRiskRes); i++ {
		riskNumMap[flowRiskRes[i].RiskLevel] += 1
		if _, ok := targetIds[flowRiskRes[i].FlowTargetID]; ok { //在目标列表中
			if flowRiskRes[i].RiskLevel < targetNumMap[flowRiskRes[i].FlowTargetID] {
				targetNumMap[flowRiskRes[i].FlowTargetID] = flowRiskRes[i].RiskLevel
			}
		}
		if flowRiskRes[i].RiskLevel < riskLevel {
			riskLevel = flowRiskRes[i].RiskLevel
		}
	}
	for k, v := range riskNumMap {
		riskNum[k-1] = v
	}
	for _, v := range targetNumMap {
		if v < 3 {
			targetNum[0] += 1
		}
		if v >= 3 && v <= 5 {
			targetNum[v-2] += 1
		}
	}
	return riskNum, targetNum, riskLevel, nil
}

// FlowTargetInfoByTaskId 获取流量目标id数组和目标数组
func (f *Flow) FlowRiskStataByTaskIds(ctx context.Context, flowTaskId []int) (map[int][]int, map[int]int, error) {
	var (
		flowRiskMysqls mysqls.Flowrisk
		riskNum        = make(map[int][]int, 0)
		riskNumMap     = make(map[int]map[int]int, 0)
		riskLevel      = make(map[int]int, 0)
	)
	//初始化参数
	for i := 0; i < len(flowTaskId); i++ {
		riskNum[flowTaskId[i]] = make([]int, 5)
		riskLevel[flowTaskId[i]] = enums.FlowRiskLevelSafe
		riskNumMap[flowTaskId[i]] = make(map[int]int, 5)
	}
	flowRiskRes, err := flowRiskMysqls.GetFlowriskListByTaskIds(ctx, flowTaskId)
	if err != nil {
		return nil, nil, err
	}
	for i := 0; i < len(flowRiskRes); i++ {
		if v, ok := riskNumMap[flowRiskRes[i].FlowTaskID]; ok {
			riskNumMap[flowRiskRes[i].FlowTaskID][flowRiskRes[i].RiskLevel] = v[flowRiskRes[i].RiskLevel] + 1
		}
		if v, ok := riskLevel[flowRiskRes[i].FlowTaskID]; ok {
			if flowRiskRes[i].RiskLevel < v {
				riskLevel[flowRiskRes[i].FlowTaskID] = flowRiskRes[i].RiskLevel
			}
		}
	}
	for k, v := range riskNumMap {
		for k1, v1 := range v {
			riskNum[k][k1-1] = v1
		}
	}
	return riskNum, riskLevel, nil
}

// FlowRiskList 漏洞列表
func (f *Flow) FlowRiskList(ctx context.Context, search string, riskLevel int, flowTaskId int, page, size int) ([]mysqls.Flowrisk, int64, error) {
	var flowRiskMysql mysqls.Flowrisk
	return flowRiskMysql.GetFlowriskList(ctx, search, riskLevel, flowTaskId, page, size)
}

// FlowRiskInfo 漏洞详情
func (f *Flow) FlowRiskInfo(ctx context.Context, flowRiskId int) (mysqls.Flowrisk, error) {
	var flowRiskMysql mysqls.Flowrisk
	return flowRiskMysql.GetFlowrisk(ctx, flowRiskId)
}

// FlowBaseDel 漏洞删除
func (f *Flow) FlowRiskDel(ctx context.Context, ids any) error {
	var flowBaseMysql mysqls.Flowrisk
	return flowBaseMysql.DeleteFlowrisk(ctx, ids)
}

// FlowTaskBaseList 被动流量列表
func (f *Flow) FlowBaseList(ctx context.Context, search string, flowTaskId int, page, size int) ([]mysqls.Flowbase, int64, error) {
	var flowBaseMysql mysqls.Flowbase
	return flowBaseMysql.GetFlowbaseList(ctx, search, flowTaskId, page, size)
}

// FlowBaseInfo 被动流量详情
func (f *Flow) FlowBaseInfo(ctx context.Context, flowBaseId int) (mysqls.Flowbase, error) {
	var flowBaseMysql mysqls.Flowbase
	return flowBaseMysql.GetFlowbase(ctx, flowBaseId)
}

// FlowBaseDel 被动流量删除
func (f *Flow) FlowBaseDel(ctx context.Context, ids any) error {
	var flowBaseMysql mysqls.Flowbase
	return flowBaseMysql.DeleteFlowbase(ctx, ids)
}

// FlowBaseAdd 新增被动流量
func (f *Flow) FlowBaseAdd(ctx context.Context, flowTaskId, flowTargetID int, flownetresult *FlowResultNetflowMsg) error {
	protocol := 1
	flownetresult.Protocol = strings.ToLower(flownetresult.Protocol)
	ishttps := strings.Contains(flownetresult.Protocol, "https")
	if ishttps {
		protocol = 2
	}
	var flowBaseMysql = mysqls.Flowbase{
		FlowTaskID:      flowTaskId,
		FlowTargetID:    flowTargetID,
		Hash:            flownetresult.Hash,
		YakID:           0,
		Host:            flownetresult.Host,
		Url:             flownetresult.Url,
		IP:              flownetresult.Ip,
		Method:          flownetresult.Method,
		Protocol:        protocol,
		RespTitle:       flownetresult.RespTitle,
		RespCode:        flownetresult.RespCode,
		RespContentType: flownetresult.RespContentType,
		ReqHeader:       flownetresult.ReqHeader,
		RespHeader:      flownetresult.RespHeader,
		RespContent:     flownetresult.RespContent,
		LikeField:       flownetresult.Host + flownetresult.RespTitle,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	return flowBaseMysql.AddFlowbase(ctx)
}

// FlowRiskAdd 新增被动流量漏洞
func (f *Flow) FlowRiskAdd(ctx context.Context, flowTaskId, flowTargetID int, flownetresult *FlowResultRiskMsg) error {
	var flowRiskMysql = mysqls.Flowrisk{
		FlowTaskID:      flowTaskId,
		FlowTargetID:    flowTargetID,
		Hash:            flownetresult.Hash,
		YakID:           0,
		Host:            flownetresult.Host,
		IP:              flownetresult.Ip,
		IPInteger:       "",
		Port:            flownetresult.Port,
		Title:           flownetresult.Title,
		RiskType:        flownetresult.RiskTypeVerbose,
		RiskTypeVerbose: flownetresult.RiskTypeVerbose,
		Payload:         "",
		Detail:          flownetresult.Detail,
		RiskLevel:       flownetresult.RiskLevel,
		Request:         flownetresult.Request,
		Response:        flownetresult.Response,
		Parameter:       flownetresult.Parameter,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	return flowRiskMysql.AddFlowrisk(ctx)
}

// FlowTaskLogDel 清空流量任务日志
func (f *Flow) FlowLogDel(ctx context.Context, flowTaskId int) error {
	var flowLogMysqls mysqls.Flowlog
	return flowLogMysqls.DeleteByFlowTaskIds(ctx, []int{flowTaskId})
}

// FlowTaskLogInfo 查询流量任务日志
func (f *Flow) FlowLogInfo(ctx context.Context, search string, flowTaskId int, page, size int) ([]mysqls.Flowlog, int64, error) {
	var flowLogMysqls mysqls.Flowlog
	return flowLogMysqls.GetFlowlogList(ctx, search, flowTaskId, page, size)
}

// FlowLogAdd 新增流量任务日志
func (f *Flow) FlowLogAdd(ctx context.Context, flowTaskId int, content string) error {
	var flowLogMysqls = mysqls.Flowlog{
		FlowTaskID: flowTaskId,
		Content:    content,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	return flowLogMysqls.AddFlowlog(ctx)
}

// FlowLogAdd 新增流量任务日志
func (f *Flow) FlowTagsAdd(ctx context.Context, flowTaskId int, tagResult *FlowResultTagsMsg) error {
	var (
		flowBaseMysql mysqls.Flowbase
		tags          = make([]string, 0)
	)
	if len(tagResult.Tags) > 0 {
		tagResult.Tags = strings.Replace(tagResult.Tags, "[", "", -1)
		tagResult.Tags = strings.Replace(tagResult.Tags, "]", "", -1)
		tags = strings.Split(tagResult.Tags, " ")
	}
	return flowBaseMysql.GetFlowBaseByTaskIdsTargetIdUrlMethod(ctx, flowTaskId, tagResult.Hash, tagResult.Method, tagResult.Url, tags)
}

// FlowGetVulNameMap 获取漏洞名称映射
func (f *Flow) FlowGetVulNameMap(ctx context.Context) interface{} {
	param := map[string]interface{}{
		"scriptType": "mitm",
		"page":       "1",
		"size":       "1000",
	}
	res, _ := httpclients.GetVulNameByType(param)
	vulNameList := make([]map[string]interface{}, 0)
	for _, vul := range res.List {
		if vul.VulName == "越权漏洞" || vul.VulName == "遍历测试" || vul.VulName == "未授权漏洞" {
			vulNameList = append(vulNameList, map[string]interface{}{
				"value":     vul.Pocname,
				"label":     vul.VulName,
				"isDefault": true,
			})
		} else if strings.HasPrefix(vul.VulName, "逻辑") {
			vulNameList = append(vulNameList, map[string]interface{}{
				"value":     vul.Pocname,
				"label":     vul.VulName,
				"isDefault": true,
			})
		} else {
			//vulNameList = append(vulNameList, map[string]interface{}{
			//	"value": vul.Pocname,
			//	"label": vul.VulName,
			//})
		}
	}
	return vulNameList
}

// FlowTaskEdit 流量分析任务编辑
func (f *Flow) FlowTaskEdit(ctx context.Context, flowTaskId int, taskName, port string, targetUrl []string, networkCard string, expireTime, userId int, otherConfig, vulConfig string) ([]mysqls.Flowtarget, error) {
	// 第一步 启动事务
	var (
		flowTargetMysqls mysqls.Flowtarget
		flowTargetData   = make([]mysqls.Flowtarget, 0)
	)

	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	// 第二步 更新任务信息
	param := map[string]interface{}{
		"task_name":    taskName,
		"netword_card": networkCard,
		"port":         port,
		"expire_time":  expireTime,
		"update_time":  time.Now(),
		"other_config": otherConfig,
		"vul_config":   vulConfig,
	}
	var flowMysqls mysqls.Flowtask
	err := flowMysqls.UpdateFlowtask(dCtx, flowTaskId, param)
	if err != nil {
		return flowTargetData, err
	}
	//第三步 获取被动流量任务

	flowTaskMysqls, err := flowMysqls.GetFlowtask(dCtx, flowTaskId)
	// 第四步 删除上个任务的目标
	err = flowTargetMysqls.DeleteFlowtargetByTaskIds(dCtx, []int{flowTaskMysqls.ID})
	// 第五步 添加新的目标
	for i := 0; i < len(targetUrl); i++ {
		var tmpflowTargetData = mysqls.Flowtarget{
			FlowTaskID: flowTaskMysqls.ID,
			TargetURL:  targetUrl[i],
			Status:     enums.FlowTargetStatusTwo,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		flowTargetData = append(flowTargetData, tmpflowTargetData)
	}
	err = flowTargetMysqls.AddManyFlowtarget(dCtx, flowTargetData)
	if err != nil {
		return flowTargetData, err
	}
	// 第六步 发送给决策引擎被动流量任务开启
	err = sendDecisionCreate(&flowTaskMysqls, flowTargetData, otherConfig, vulConfig)
	if err != nil {
		return flowTargetData, err
	}
	// 第七步 提交事务
	if err := tx.Commit().Error; err != nil { //提交事务
		return flowTargetData, err
	}
	return flowTargetData, nil
}
