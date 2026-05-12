package application

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/enums"
	"strconv"
	"strings"
)

type ReportVerify struct {
}

// Upload 上传报告逻辑处理
func (r *ReportVerify) Upload(ctx context.Context, req *typespec.ReportVerifyUploadReq, fileName, tempFileName string, resp *typespec.ReportVerifyUploadResp) error {
	var srv services.ReportVerify
	fileType := srv.AnalysisFileType(ctx, fileName)
	var portMap map[string]string
	var vulMap map[string][]map[string]interface{}
	fileInfoMap := make(map[string]string, 0)
	var producer int
	sourceMap := make(map[string]string, 0)
	data, _ := os.ReadFile(tempFileName)

	if fileType == enums.ReportVerifyFileTypeHtml {
		sourceMap["tempFile"] = string(data)
	} else if fileType == enums.ReportVerifyFileTypeZip {
		err := srv.AnalysisZip(ctx, tempFileName, sourceMap)
		if err != nil {
			return err
		}
	}
	for name, content := range sourceMap {
		tempProducer := srv.AnalysisProducer(ctx, content)
		if fileType == enums.ReportVerifyFileTypeUnKnown {
			continue
		}
		if tempProducer == enums.ReportVerifyProducerUnKnown {
			continue
		}
		producer = tempProducer
		if tempProducer == enums.ReportVerifyProducerTianJing {
			portMap = srv.AnalysisTianJingPort(ctx, content)
			vulMap = srv.AnalysisTianJingVul(ctx, content)
		} else if tempProducer == enums.ReportVerifyProducerNsfocus {
			vulMap = srv.AnalysisNsFocus(ctx, content)
			portMap = srv.AnalysisNsPort(ctx, vulMap)
		} else if tempProducer == enums.ReportVerifyProducerTianJing2025 {
			if strings.Contains(name, "主机漏扫报表") {
				continue
			}
			portMap = srv.AnalysisTianJingPort2025(ctx, content)
			vulMap = srv.AnalysisTianJingVul2025(ctx, content)
		}
	}
	if len(vulMap) == 0 {
		return errors.New("报告格式不正确或报告中未检测到漏洞")
	}

	fileInfoMap["fileName"] = fileName
	fileInfoMap["targetNumber"] = strconv.Itoa(len(vulMap))
	fileInfoMap["fileSize"] = strconv.Itoa(len(string(data))/1024) + "k"

	go func() {
		defer func() {
			os.Remove(tempFileName)
		}()
		taskId, err := srv.SaveTask(ctx, req.Name, fileInfoMap, producer, req.UserId)
		if err != nil {
			log.Println(err)
			return
		}
		err = srv.SaveTarget(ctx, taskId, portMap, vulMap)
		if err != nil {
			log.Println(err)
			return
		}
		err = srv.SaveVul(ctx, taskId, vulMap)
		if err != nil {
			log.Println(err)
			return
		}
		taskStats, _, err := srv.ReportVerifyOverView(ctx, taskId)
		if err != nil {
			log.Println(err)
			return
		}
		err = srv.UpdateTaskOverView(ctx, taskId, taskStats)
		if err != nil {
			log.Println(err)
			return
		}
	}()
	return nil
}

// TaskList 任务列表
func (r *ReportVerify) TaskList(ctx context.Context, req *typespec.ReportVerifyTaskListReq, resp *typespec.ReportVerifyTaskListResp) error {
	var srv services.ReportVerify
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
	reportVerifyList, total, err := srv.TaskList(ctx, req.Page, req.Size, req.Risk, req.Producer, req.StartTime, req.EndTime, req.Search, userIdList)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(reportVerifyList); i++ {
		var tmp typespec.ReportVerifyTaskListRespItems
		tmp.Id = reportVerifyList[i].ID
		tmp.Name = reportVerifyList[i].Name
		tmp.ExecuteType = reportVerifyList[i].ExecuteType
		tmp.ExecuteTypeName = enums.TaskTaskEnum.ExecTypeEnum(reportVerifyList[i].ExecuteType)
		tmp.Exp = reportVerifyList[i].Exp
		tmp.Verify = reportVerifyList[i].Verify
		tmp.Failed = reportVerifyList[i].Failed
		tmp.Risk = reportVerifyList[i].Risk
		tmp.RiskName = enums.GetTargetRisk(reportVerifyList[i].Risk)
		tmp.Producer = reportVerifyList[i].Producer
		tmp.ProducerName = enums.ReportVerifyEnum.ReportVerifyProducerEnum(reportVerifyList[i].Producer)
		tmp.UpdateTime = reportVerifyList[i].UpdateTime.Format(enums.TimeLayout)
		tmp.Status = reportVerifyList[i].Status
		tmp.StatusName = enums.TaskTaskEnum.StatusEnum(reportVerifyList[i].Status)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// TaskDetail 任务详情
func (r *ReportVerify) TaskDetail(ctx context.Context, req *typespec.ReportVerifyTaskDetailReq, resp *typespec.ReportVerifyTaskDetailResp) error {
	var (
		srv     services.ReportVerify
		userSrv services.User
	)
	taskDetail, err := srv.TaskDetail(ctx, req.Id)
	if err != nil {
		return err
	}
	resp.Id = taskDetail.ID
	resp.Name = taskDetail.Name
	resp.ExecuteType = taskDetail.ExecuteType
	resp.ExecuteTypeName = enums.TaskTaskEnum.ExecTypeEnum(taskDetail.ExecuteType)
	resp.Exp = taskDetail.Exp
	resp.Verify = taskDetail.Verify
	resp.Failed = taskDetail.Failed
	user, _ := userSrv.GetForId(ctx, taskDetail.User)
	resp.User = user.Username

	resp.Risk = taskDetail.Risk
	resp.RiskName = enums.TaskTaskEnum.RiskEnum(taskDetail.Risk)
	resp.Producer = taskDetail.Producer
	resp.ProducerName = enums.ReportVerifyEnum.ReportVerifyProducerEnum(taskDetail.Producer)
	resp.UpdateTime = taskDetail.UpdateTime.Format(enums.TimeLayout)
	resp.Status = taskDetail.Status
	resp.StatusName = enums.TaskTaskEnum.StatusEnum(taskDetail.Status)
	resp.CreateTime = taskDetail.CreateTime.Format(enums.TimeLayout)
	//文件信息
	var fileInfoMap map[string]string
	err = json.Unmarshal([]byte(taskDetail.Fileinfo), &fileInfoMap)
	resp.FileName = fileInfoMap["fileName"]
	resp.TargetNumber = fileInfoMap["targetNumber"]
	resp.FileSize = fileInfoMap["fileSize"]
	if err != nil {
		return err
	}
	return nil
}

// TargetList 目标列表
func (r *ReportVerify) TargetList(ctx context.Context, req *typespec.ReportVerifyTargetListReq, resp *typespec.ReportVerifyTargetListResp) error {
	var srv services.ReportVerify
	targetList, total, err := srv.TargetList(ctx, req.TaskId, req.Risk, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(targetList); i++ {
		var tmp typespec.ReportVerifyTargetRespItems
		tmp.Id = targetList[i].ID
		tmp.TaskId = targetList[i].TaskId
		tmp.Target = targetList[i].Target
		tmp.Os = targetList[i].Os
		tmp.Risk = targetList[i].Risk
		tmp.RiskName = enums.GetTargetRisk(targetList[i].Risk)
		tmp.Exp = targetList[i].Exp
		tmp.Verify = targetList[i].Verify
		tmp.Failed = targetList[i].Failed
		tmp.UnVerify = targetList[i].Unverify
		tmp.Status = targetList[i].Status
		tmp.StatusName = enums.GetTargetStatus(targetList[i].Status)
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// PortList 端口列表
func (r *ReportVerify) PortList(ctx context.Context, req *typespec.ReportVerifyPortListReq, resp *typespec.ReportVerifyPortListResp) error {
	var srv services.ReportVerify
	portList, total, err := srv.PortList(ctx, req.TaskId, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(portList); i++ {
		var tmp typespec.ReportVerifyPortRespItems
		tmp.Id = portList[i].ID
		tmp.TaskId = portList[i].TaskId
		tmp.Port = portList[i].Port
		tmp.Scheme = portList[i].Scheme
		tmp.Service = portList[i].Service
		tmp.Component = portList[i].Component
		tmp.Target = portList[i].Target
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// VulList 漏洞列表
func (r *ReportVerify) VulList(ctx context.Context, req *typespec.ReportVerifyVulListReq, resp *typespec.ReportVerifyVulListResp) error {
	var srv services.ReportVerify
	vulList, total, err := srv.VulList(ctx, req.TaskId, req.Page, req.Size, req.Risk, req.Status, req.Search)
	if err != nil {
		return err
	}
	resp.Total = total
	for i := 0; i < len(vulList); i++ {
		var tmp typespec.ReportVerifyVulRespItems
		tmp.Id = vulList[i].ID
		tmp.TaskId = vulList[i].TaskId
		tmp.Name = vulList[i].Name
		tmp.Risk = vulList[i].Risk
		tmp.RiskName = enums.ToolsVulnerabilityEnum.GetRiskEnum(vulList[i].Risk)
		tmp.Status = vulList[i].Status
		tmp.StatusName = enums.ReportVerifyEnum.ReportVerifyStatusEnum(vulList[i].Status)
		tmp.Location = vulList[i].Location
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// Enum 枚举接口
func (r *ReportVerify) Enum(ctx context.Context, req *typespec.ReportVerifyEnumReq, resp *typespec.ReportVerifyEnumResp) error {
	// 枚举信息 - 状态
	var (
		rvSrv   services.ReportVerify
		taskSrv services.TaskTask
	)
	resp.ProducerType = rvSrv.ProducerEnum()
	resp.Status = taskSrv.StatusEnum()
	resp.Risk = taskSrv.RiskLevelEnum()
	resp.VulRisk = rvSrv.VulRiskEnum()
	resp.ExecuteType = rvSrv.ExecuteTypeEnum()
	//resp.ExecuteType = taskSrv.ExecuteTypeEnum()
	resp.VulStatus = rvSrv.VulStatusEnum()
	return nil
}

// StatsInfo 统计信息
func (r *ReportVerify) StatsInfo(ctx context.Context, req *typespec.ReportVerifyStatsInfoReq, resp *typespec.ReportVerifyStatsInfoResp) error {
	var srv services.ReportVerify
	overview, err := srv.GetTaskOverView(ctx, req.TaskId)
	if err != nil {
		return nil
	}
	resp.AllVul = overview["allVul"]
	resp.UnVerify = overview["unVerify"]
	resp.Verify = overview["verify"]
	resp.Failed = overview["failed"]
	resp.Exp = overview["exp"]
	resp.HighVul = overview["highVul"]
	resp.MiddleVul = overview["middleVul"]
	resp.LowVul = overview["lowVul"]
	resp.AllTarget = overview["allTarget"]
	resp.HighTarget = overview["highTarget"]
	resp.MiddleTarget = overview["middleTarget"]
	resp.LowTarget = overview["lowTarget"]
	resp.SafeTarget = overview["safeTarget"]
	resp.AliveTarget = overview["aliveTarget"]
	return nil
}

// TaskStop 结束任务
func (r *ReportVerify) TaskStop(ctx context.Context, req *typespec.ReportVerifyTaskStopReq, resp *typespec.ReportVerifyTaskStopResp) error {
	var srv services.ReportVerify
	err := srv.UpdateTaskStatus(ctx, req.TaskId, enums.TaskStatusFinish)
	if err != nil {
		return err
	}
	err = srv.UpdateTargetStatusByTaskId(ctx, req.TaskId, enums.TargetStatusFinish)
	if err != nil {
		return err
	}
	return nil
}

// TaskDelete 删除任务
func (r *ReportVerify) TaskDelete(ctx context.Context, req *typespec.ReportVerifyTaskDeleteReq, resp *typespec.ReportVerifyTaskDeleteResp) error {
	var srv services.ReportVerify
	for _, taskId := range strings.Split(req.TaskId, ",") {
		id, err := strconv.Atoi(taskId)
		if err != nil {
			return err
		}
		err = srv.TaskDelete(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// TargetDelete 删除目标
func (r *ReportVerify) TargetDelete(ctx context.Context, req *typespec.ReportVerifyTargetDeleteReq, resp *typespec.ReportVerifyTargetDeleteResp) error {
	var srv services.ReportVerify
	for _, targetId := range strings.Split(req.TargetId, ",") {
		id, err := strconv.Atoi(targetId)
		if err != nil {
			return err
		}
		err = srv.TargetDelete(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// VulDelete 删除漏洞
func (r *ReportVerify) VulDelete(ctx context.Context, req *typespec.ReportVerifyVulDeleteReq, resp *typespec.ReportVerifyVulDeleteResp) error {
	var srv services.ReportVerify
	for _, vulId := range strings.Split(req.VulId, ",") {
		id, err := strconv.Atoi(vulId)
		if err != nil {
			return err
		}
		err = srv.VulDelete(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// VulDetail 漏洞详情
func (r *ReportVerify) VulDetail(ctx context.Context, req *typespec.ReportVerifyVulDetailReq, resp *typespec.ReportVerifyVulDetailResp) error {
	var srv services.ReportVerify
	vul, err := srv.VulDetail(ctx, req.VulId)
	if err != nil {
		return err
	}
	resp.Id = vul.ID
	resp.Name = vul.Name
	resp.Risk = vul.Risk
	resp.Desc = vul.Desc
	resp.Fix = vul.Fix
	resp.Cve = vul.Cve
	resp.Cnnvd = vul.Cnnvd
	resp.Cvss = vul.Cvss
	return nil
}
