package application

import (
	"context"
	"encoding/json"
	"errors"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/services"
	"smart/tools/enums"
	"strings"
)

type Report struct {
}

// ReportEnum 报告枚举
func (r *Report) ReportEnum(ctx context.Context, resp *typespec.ReportEnumResp) (err error) {
	var report services.Report
	//报告类型,报告状态,报告格式枚举
	resp.Type, resp.Status, resp.Format = report.GetReportEnum()
	//报告内容枚举
	resp.TaskContent, resp.TargetContent, err = report.GetReportContentEnum(ctx)
	if err != nil {
		return err
	}
	return nil
}

// ReportList 报告列表及筛选
func (r *Report) ReportList(ctx context.Context, req *typespec.ReportListReq, resp *typespec.ReportListResp) error {
	var (
		report      services.Report
		reportenums enums.ReportEnums
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
	reportRes, count, err := report.ReportList(ctx, req.Search, req.Page, req.Size, userIdList)
	if err != nil {
		return err
	}
	resp.Total = count
	for i := 0; i < len(reportRes); i++ {
		var tmp typespec.ReportListRespItems
		tmp.Id = reportRes[i].ID
		tmp.Name = reportRes[i].Name
		tmp.Type = reportRes[i].Type
		tmp.TypeName = reportenums.GetReportType(reportRes[i].Type)
		tmp.Status = reportRes[i].Status
		tmp.StatusName = reportenums.GetReportStatus(reportRes[i].Status)
		tmp.Format = reportRes[i].Format
		tmp.FormatName = reportenums.GetReportFormat(reportRes[i].Format)
		tmp.CreateTime = "--"
		if reportRes[i].Status == enums.ReportStatusFinish {
			tmp.CreateTime = reportRes[i].UpdateTime.Format(enums.TimeYMDHMinLayout)
		}
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// ReportDownload 报告下载
func (r *Report) ReportDownload(ctx context.Context, req *typespec.ReportDownloadReq, resp *typespec.ReportDownloadResp) error {
	var report services.Report
	reportRes, err := report.ReportDownload(ctx, req.ReportId)
	if err != nil {
		return err
	}
	resp.Id = reportRes.ID
	resp.Name = reportRes.Name
	resp.Type = reportRes.Type
	resp.Status = reportRes.Status
	resp.ConfigJson = reportRes.ConfigJSON
	resp.Format = reportRes.Format
	resp.Content = reportRes.Content
	resp.CreateTime = reportRes.CreateTime.Format(enums.TimeLayout)
	resp.UpdateTime = reportRes.UpdateTime.Format(enums.TimeLayout)
	return nil
}

// ReportDel 报告删除
func (r *Report) ReportDel(ctx context.Context, req *typespec.ReportDelReq) error {
	idArray := strings.Split(req.ReportIds, ",")
	var report services.Report
	return report.ReportDel(ctx, idArray)
}

// ReportSave 保存
func (r *Report) ReportSave(ctx context.Context, req *typespec.ReportSaveReq) error {
	var report services.Report
	if req.Type == enums.ReportTypeVulScanTask || req.Type == enums.ReportTypeVulScanTarget {
		reportId, _ := report.ReportSave(ctx, req.Name, req.Type, req.ConfigJson, req.Format, req.UserId)
		clientReq := httpclients.VulScanReportSaveReq{
			ReportId:   reportId,
			Name:       req.Name,
			Type:       req.Type,
			ConfigJson: req.ConfigJson,
			Format:     req.Format,
			OutputType: req.OutputType,
			ObjIDName:  req.ObjIDName,
		}
		resp, err := httpclients.VulScanReportSave(ctx, clientReq)
		if err != nil {
			return err
		}
		if resp.Code != 200 {
			return errors.New(resp.Msg)
		}
		return nil
	}
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
		report.ReportSave(ctx, "测试目标报告", req.Type, req.ConfigJson, req.Format, req.UserId)
	} else {
		report.ReportSave(ctx, req.Name, req.Type, req.ConfigJson, req.Format, req.UserId)
	}
	return nil
}

// execConfigJson 处理configJson
func execBatchConfigJson(configJson string) map[string]string {
	var configJsonStr map[string]interface{}
	configJsonRes := make(map[string]string)
	json.Unmarshal([]byte(configJson), &configJsonStr)
	objIdMid := configJsonStr["objId"]
	if _, ok := objIdMid.(string); ok {
		objIDs := objIdMid.(string)
		ids := strings.Split(objIDs, ",")
		for _, v := range ids {
			configJsonStr["objId"] = v
			str, _ := json.Marshal(configJsonStr)
			configJsonRes[v] = string(str)
		}
	}
	return configJsonRes
}
