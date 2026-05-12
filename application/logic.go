package application

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/services"
	"smart/tools/data"
	"smart/tools/enums"
	"strconv"
	"strings"
)

type Logic struct {
}

// 逻辑漏洞 - 任务新建
func (a *Logic) TaskCreate(ctx context.Context, req *typespec.LogicTaskCreateReq, resp *typespec.LogicTaskCreateResp) error {
	var (
		srv            services.Logic
		analysisTarget data.TaskCheckTaskAnalysisTarget
	)

	if ctx.Value("uid") == nil {
		return errors.New("用户未登录")
	}
	uid := ctx.Value("uid").(int)
	if req.Type == enums.LogicTypeBeyondPermission {
		var scanConfig services.BeyondPermConfig
		err := json.Unmarshal([]byte(req.ScanConfig), &scanConfig)
		if err != nil {
			return err
		}
	} else if req.Type == enums.LogicTypeTraverseTesting {
		var scanConfig services.SensInfoConfig
		err := json.Unmarshal([]byte(req.ScanConfig), &scanConfig)
		if err != nil {
			return err
		}
	} else if req.Type == enums.LogicTypeUnAuthAccess {
		var scanConfig services.UnAuthAccessConfig
		err := json.Unmarshal([]byte(req.ScanConfig), &scanConfig)
		if err != nil {
			return err
		}
	} else {
		return errors.New("不支持的扫描类型")
	}

	analysisTarget.AnalysisTarget(req.TargetUrl, "")
	errorTargetList := analysisTarget.ErrorTargetList
	if len(errorTargetList) > 0 {
		return errors.New(strings.Join(errorTargetList, ","))
	}

	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	targetList := analysisTarget.TargetList
	taskId, err := srv.TaskSave(dCtx, req.Name, req.TargetUrl, req.ScanConfig, req.Type, uid)
	if err != nil {
		return err
	}
	for _, targetUrl := range targetList {
		_, err := srv.TargetSave(dCtx, targetUrl, taskId, req.Type)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		log.Error("logic task save transaction commit err " + err.Error())
	}
	resp.TaskId = taskId
	return nil
}

// 逻辑漏洞 - 任务结束
func (a *Logic) TaskStop(ctx context.Context, req *typespec.LogicTaskStopReq, resp *typespec.LogicTaskStopResp) error {
	var srv services.Logic
	err := srv.Stop(ctx, req.Id)
	if err != nil {
		return err
	}
	return nil
}

// 逻辑漏洞 - 任务删除
func (a *Logic) TaskDelete(ctx context.Context, req *typespec.LogicTaskDelReq, resp *typespec.LogicTaskDelResp) error {
	var srv services.Logic
	for _, id := range strings.Split(req.Ids, ",") {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		err = srv.Delete(ctx, idInt)
		if err != nil {
			return err
		}
	}
	return nil
}

// 逻辑漏洞 - 任务列表
func (a *Logic) TaskList(ctx context.Context, req *typespec.LogicTaskListReq, resp *typespec.LogicTaskListResp) error {
	var srv services.Logic
	//查询数据
	taskList, total, err := srv.List(ctx, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(taskList); i++ {
		var temp typespec.LogicTaskListRespItems
		temp.Id = taskList[i].ID
		temp.Name = taskList[i].Name
		temp.Type = taskList[i].Type
		temp.TypeName = enums.Logic{}.GetScanTypeName(taskList[i].Type)
		temp.Risk = taskList[i].Risk
		temp.RiskName = enums.Logic{}.GetRiskName(taskList[i].Risk)
		temp.Status = taskList[i].Status
		temp.StatusName = enums.Logic{}.GetStatusName(taskList[i].Status)
		temp.UserID = taskList[i].UserID
		temp.UpdateTime = taskList[i].UpdateTime.Format(enums.TimeLayout)
		resp.List = append(resp.List, temp)
	}
	return nil
}

// 逻辑漏洞 - 目标列表
func (a *Logic) TargetList(ctx context.Context, req *typespec.LogicTargetListReq, resp *typespec.LogicTargetListResp) error {
	var srv services.Logic
	//查询数据
	taskList, total, err := srv.TargetList(ctx, req.TaskId, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(taskList); i++ {
		var temp typespec.LogicTargetListRespItems
		temp.Id = taskList[i].ID
		temp.TargetUrl = taskList[i].TargetURL
		temp.Status = taskList[i].Status
		temp.StatusName = enums.Logic{}.GetStatusName(taskList[i].Status)
		temp.Risk = taskList[i].Risk
		temp.RiskName = enums.Logic{}.GetRiskName(taskList[i].Risk)
		temp.Type = taskList[i].Type
		temp.TypeName = enums.Logic{}.GetScanTypeName(taskList[i].Type)
		temp.IsAlive = taskList[i].IsAlive
		temp.IsAliveName = enums.Logic{}.GetIsAliveName(taskList[i].IsAlive)
		resp.List = append(resp.List, temp)
	}
	return nil
}

// 逻辑漏洞 - 漏洞列表
func (a *Logic) VulList(ctx context.Context, req *typespec.LogicVulListReq, resp *typespec.LogicVulListResp) error {
	var srv services.Logic
	//查询数据
	vulList, total, err := srv.VulList(ctx, req.TaskId, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(vulList); i++ {
		var temp typespec.LogicVulListRespItems
		temp.Id = vulList[i].ID
		temp.TargetUrl = vulList[i].TargetUrl
		temp.Type = vulList[i].Type
		temp.TypeName = enums.Logic{}.GetScanTypeName(vulList[i].Type)
		temp.Risk = vulList[i].Risk
		temp.RiskName = enums.Logic{}.GetVulRiskName(vulList[i].Risk)
		temp.Location = vulList[i].Location
		temp.CreateTime = vulList[i].CreateTime.Format(enums.TimeLayout)
		resp.List = append(resp.List, temp)
	}
	return nil
}

// 逻辑漏洞 - 日志列表
func (a *Logic) LogList(ctx context.Context, req *typespec.LogicLogListReq, resp *typespec.LogicLogListResp) error {
	var srv services.Logic
	//查询数据
	vulList, total, err := srv.LogList(ctx, req.TaskId, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(vulList); i++ {
		var temp typespec.LogicLogListRespItems
		temp.Id = vulList[i].ID
		temp.TargetURL = vulList[i].TargetURL
		temp.IsAlive = vulList[i].IsAlive
		temp.IsAliveName = enums.Logic{}.GetIsAliveName(vulList[i].IsAlive)
		temp.Status = vulList[i].Status
		temp.StatusName = enums.Logic{}.GetStatusName(vulList[i].Status)
		temp.CreateTime = vulList[i].CreateTime.Format(enums.TimeLayout)
		temp.StartTime = vulList[i].StartTime.Format(enums.TimeLayout)
		temp.EndTime = vulList[i].EndTime.Format(enums.TimeLayout)
		resp.List = append(resp.List, temp)
	}
	return nil
}

// 逻辑漏洞 - 日志列表
func (a *Logic) VulInfo(ctx context.Context, req *typespec.LogicVulInfoReq, resp *typespec.LogicVulInfoResp) error {
	var srv services.Logic
	//查询数据
	vulInfo, err := srv.GetVulById(ctx, req.Id)
	if err != nil {
		return err
	}
	resp.Name = vulInfo.Name
	resp.Risk = vulInfo.Risk
	resp.RiskName = enums.Logic{}.GetVulRiskName(vulInfo.Risk)
	resp.Description = vulInfo.Description
	resp.FixSuggest = vulInfo.FixSuggest
	resp.Location = vulInfo.Location
	// 获取 success_payload_flag
	var tempVulMap map[string]interface{}
	err = json.Unmarshal([]byte(vulInfo.VulResult), &tempVulMap)
	if tempVulMap["payload_success_flag"] != nil {
		resp.PayloadSuccessFlag = tempVulMap["payload_success_flag"].(string)
	}
	if tempVulMap["payload"] != nil {
		resp.Payload = strings.Trim(tempVulMap["payload"].(string), `\"`)
		delete(tempVulMap, "payload")
	}
	vulResultByte, _ := json.Marshal(tempVulMap)
	resp.Result = string(vulResultByte)
	// 获取响应报文
	var tempMap []map[string]string
	err = json.Unmarshal([]byte(vulInfo.VerMsg), &tempMap)
	if err == nil {
		resp.VerMsg = tempMap
	}
	return nil
}

// 逻辑漏洞 - 日志详情
func (a *Logic) LogInfo(ctx context.Context, req *typespec.LogicLogInfoReq, resp *typespec.LogicLogInfoResp) error {
	var srv services.Logic
	//查询数据
	logInfoList, err := srv.GetLogById(ctx, req.Id)
	if err != nil {
		return err
	}
	for _, logInfo := range logInfoList {
		item := typespec.LogicLogInfoRespItems{
			Id:         logInfo.ID,
			Result:     logInfo.Result,
			CreateTime: logInfo.CreateTime.Format(enums.TimeLayout),
			Pocname:    logInfo.Pocname,
			TargetURL:  logInfo.TargetURL,
			TaskId:     logInfo.TaskID,
		}
		resp.List = append(resp.List, item)
	}
	return nil
}

// 逻辑漏洞 - 日志详情
func (a *Logic) LogicEnum(ctx context.Context, resp *typespec.LogicEnumResp) error {
	var (
		srv      services.Logic
		sceneSrv services.SceneTaskTemplate
	)
	resp.ScanType = srv.ScanTypeEnum()
	resp.CredPattern = srv.CredPatternEnum()
	resp.WhitePath = enums.WhitePath
	resp.BlackPath = enums.BlackPath
	resp.FuzzParam = struct {
		Character string `json:"character"` //字符型
		Number    string `json:"number"`    //数字型
	}{Character: enums.Character, Number: enums.Number}

	resp.Response.JsonKeyword = enums.JsonKeyword
	resp.Response.NoJsonSwitch = true
	resp.Response.NoJsonKeyword = enums.NoJsonKeyword
	// 爬虫
	resp.Crawler.MaxDeep = sceneSrv.CrawlerDeepEnum()
	resp.Crawler.MaxUrl = sceneSrv.CrawlerMaxUrlEnum()
	resp.Crawler.ScanRange = sceneSrv.CrawlerScanRangeEnum()
	resp.Crawler.Timeout = sceneSrv.CrawlerSingleTimeout()
	resp.Crawler.FullTimeout = sceneSrv.CrawlerFullTimeout()
	resp.Crawler.ScanRepeat = sceneSrv.CrawlerRepeatEnum()
	resp.Crawler.Sensitive = sceneSrv.CrawlerSensitive()
	resp.Crawler.BlackList = sceneSrv.CrawlerWhiteList()
	resp.Crawler.WhiteList = sceneSrv.CrawlerBlackList()
	return nil
}

// 逻辑漏洞 - 任务删除
func (a *Logic) VulDelete(ctx context.Context, req *typespec.LogicVulDeleteReq, resp *typespec.LogicVulDeleteResp) error {
	var srv services.Logic
	for _, id := range strings.Split(req.Ids, ",") {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		err = srv.VulDelete(ctx, idInt)
		if err != nil {
			return err
		}
	}
	return nil
}

// 逻辑漏洞 - 漏洞测试
func (a *Logic) VulTest(ctx context.Context, req *typespec.LogicVulTestReq, resp *typespec.LogicVulTestResp) error {
	var taskVulSrv services.TaskVul
	//解析请求报文
	verMsg, err := taskVulSrv.Base64StdDecodeString(req.VerMsg) //base64解密
	if err != nil {
		return err
	}
	method, url, headerMap, bodyMap, err := taskVulSrv.DecryptionVerMsg(verMsg) //解析字符串结构
	if err != nil {
		return err
	}
	//请求报文
	proto, status, header, bodyresult, err := httpclients.VerMsg(method, url, headerMap, bodyMap) //发送请求
	if err != nil {
		return err
	}
	respbody, err := taskVulSrv.BuildRespVerMsg(proto, status, header, bodyresult) //构建返回报文字符串
	if err != nil {
		return err
	}
	resp.RespVerMsg = taskVulSrv.Base64StdEncodeToString(respbody)
	return nil
}

// 逻辑漏洞 - 任务详情
func (a *Logic) TaskCopy(ctx context.Context, req *typespec.LogicTaskCopyReq, resp *typespec.LogicTaskCopyResp) error {
	var srv services.Logic
	task, err := srv.GetTaskById(ctx, req.Id)
	if err != nil {
		return err
	}
	resp.Id = task.ID
	resp.Name = task.Name
	resp.Type = task.Type
	resp.TypeName = enums.Logic{}.GetScanTypeName(task.Type)
	resp.TargetUrl = task.TargetUrl
	resp.ScanConfig = task.ScanConfig
	return nil
}

// ReportSave 保存
func (a *Logic) ReportSave(ctx context.Context, req *typespec.LogicReportSaveReq) error {
	var report services.Report
	// 逐个输出
	if req.Type == 2 && req.OutputType == enums.OneByOneOutputType {
		var objIDName map[string]string
		json.Unmarshal([]byte(req.ObjIDName), &objIDName)
		batchConfigJson := execBatchConfigJson(req.ConfigJson)
		for k, v := range batchConfigJson {
			name := objIDName[k]
			report.ReportSave(ctx, name+"  测试目标报告", req.OutputType, v, req.Format, req.UserId)
		}
		return nil
	}
	// 合并输出
	if req.Type == 2 && req.OutputType == enums.MergeOutputType {
		report.ReportSave(ctx, "逻辑漏洞报告", req.Type, req.ConfigJson, req.Format, req.UserId)
	} else {
		report.ReportSave(ctx, req.Name, req.Type, req.ConfigJson, req.Format, req.UserId)
	}
	return nil
}

// LogicFlowBaseList 逻辑漏洞流量列表接口
func (a *Logic) LogicFlowBaseList(ctx context.Context, req *typespec.LogicFlowBaseListReq, resp *typespec.LogicFlowBaseListResp) error {
	var srv services.Logic
	flowBaseRes, total, err := srv.FlowBaseList(ctx, req.Search, req.TaskId, req.Page, req.Size)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(flowBaseRes); i++ {
		var tmp typespec.LogicFlowBaseListRespItems
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

// LogicFlowBaseInfo 逻辑漏洞流量详情
func (a *Logic) LogicFlowBaseInfo(ctx context.Context, req *typespec.LogicFlowBaseInfoReq, resp *typespec.LogicFlowBaseInfoResp) error {
	var flowSrv services.Logic
	//查询数据
	flowBaseRes, err := flowSrv.LogicFlowBaseInfo(ctx, req.FlowBaseId)
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

// LogicFlowBaseExport 逻辑漏洞流量导出
func (a *Logic) LogicFlowBaseExport(ctx context.Context, req *typespec.LogicFlowBaseExportReq, resp *typespec.LogicFlowBaseExportResp) error {
	var srv services.Logic
	flowBaseRes, _, err := srv.FlowBaseList(ctx, "", req.TaskId, 1, 100000)
	if err != nil {
		return err
	}
	for i := 0; i < len(flowBaseRes); i++ {
		var tmp typespec.LogicFlowBaseListRespItems
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
