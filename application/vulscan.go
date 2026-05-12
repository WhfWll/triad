package application

import (
	"context"
	"github.com/pkg/errors"
	"smart/api/typespec"
	"smart/client/httpclients"
)

type VulScan struct {
}

// TaskList 任务列表
func (v *VulScan) TaskList(ctx context.Context, req *typespec.VulScanTaskListReq, resp *typespec.VulScanTaskListResp) error {
	// 配置请求参数
	var param httpclients.VulScanTaskListReq
	param.Page = req.Page
	param.Size = req.Size
	param.Risk = req.Risk
	param.Search = req.Search

	// 请求查询
	remoteRes, err := httpclients.VulScanTaskList(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}

	// 组织返回的数据
	resp.Total = remoteRes.Data.Total
	for _, v := range remoteRes.Data.List {
		resp.List = append(resp.List, typespec.VulScanTaskListResItem{
			Id:         v.Id,
			Name:       v.Name,
			Type:       v.Type,
			TypeName:   v.TypeName,
			Risk:       v.Risk,
			RiskName:   v.RiskName,
			Status:     v.Status,
			StatusName: v.StatusName,
			High:       v.High,
			Middle:     v.Middle,
			Low:        v.Low,
			Safe:       v.Safe,
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
		})
	}
	return nil
}

// TaskSave 任务保存
func (v *VulScan) TaskSave(ctx context.Context, req *typespec.VulScanTaskSaveReq, resp *typespec.VulScanTaskSaveResp) error {
	// 配置请求参数
	var param httpclients.VulScanTaskSaveReq
	param.Name = req.Name
	param.Target = req.Target
	param.ToScanPort = req.ToScanPort
	param.OnlyPortScan = req.OnlyPortScan
	// 请求查询
	remoteRes, err := httpclients.VulScanTaskSave(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}
	// 组织返回的数据
	resp.Id = remoteRes.Data.Id
	return nil
}

// TaskStop 任务停止
func (v *VulScan) TaskStop(ctx context.Context, req *typespec.VulScanTaskStopReq, resp *typespec.VulScanTaskStopResp) error {
	// 配置请求参数
	var param httpclients.VulScanTaskStopReq
	param.Id = req.Id
	// 请求查询
	remoteRes, err := httpclients.VulScanTaskStop(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}
	// 组织返回的数据
	return nil
}

// TaskDelete 任务删除
func (v *VulScan) TaskDelete(ctx context.Context, req *typespec.VulScanTaskDeleteReq, resp *typespec.VulScanTaskDeleteResp) error {
	// 配置请求参数
	var param httpclients.VulScanTaskDeleteReq
	param.Ids = req.Ids
	// 请求查询
	remoteRes, err := httpclients.VulScanTaskDelete(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}
	return nil
}

// TargetList 目标列表
func (v *VulScan) TargetList(ctx context.Context, req *typespec.VulScanTargetListReq, resp *typespec.VulScanTargetListResp) error {
	// 配置请求参数
	var param httpclients.VulScanTargetListReq
	param.Page = req.Page
	param.Size = req.Size
	param.TaskId = req.TaskId
	param.Risk = req.Risk
	param.Search = req.Search
	param.IsAlive = req.IsAlive

	// 请求查询
	remoteRes, err := httpclients.VulScanTargetList(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}

	// 组织返回的数据
	resp.Total = remoteRes.Data.Total
	for _, v := range remoteRes.Data.List {
		resp.List = append(resp.List, typespec.VulScanTargetListResItem{
			Id:          v.Id,
			TaskId:      v.TaskId,
			Ip:          v.Ip,
			Target:      v.Target,
			System:      v.System,
			Port:        v.Port,
			Risk:        v.Risk,
			RiskName:    v.RiskName,
			High:        v.High,
			Middle:      v.Middle,
			Low:         v.Low,
			IsAlive:     v.IsAlive,
			IsAliveName: v.IsAliveName,
			Status:      v.Status,
			StatusName:  v.StatusName,
			CreateTime:  v.CreateTime,
			UpdateTime:  v.UpdateTime,
		})
	}
	return nil
}

// VulList 漏洞列表
func (v *VulScan) VulList(ctx context.Context, req *typespec.VulScanVulListReq, resp *typespec.VulScanVulListResp) error {
	// 配置请求参数
	var param httpclients.VulScanVulListReq
	param.Page = req.Page
	param.Size = req.Size
	param.TaskId = req.TaskId
	param.Risk = req.Risk
	param.Search = req.Search

	// 请求查询
	remoteRes, err := httpclients.VulScanVulList(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}

	// 组织返回的数据
	resp.Total = remoteRes.Data.Total
	for _, v := range remoteRes.Data.List {
		resp.List = append(resp.List, typespec.VulScanVulListResItem{
			Id:          v.Id,
			TaskID:      v.TaskID,
			TargetID:    v.TargetID,
			Name:        v.Name,
			Port:        v.Port,
			Risk:        v.Risk,
			RiskName:    v.RiskName,
			CreateTime:  v.CreateTime,
			UpdateTime:  v.UpdateTime,
			Ip:          v.Ip,
			Cwe:         v.Cwe,
			PublishDate: v.PublishDate,
			Cve:         v.Cve,
		})
	}
	return nil
}

// VulDetail 漏洞详情
func (v *VulScan) VulDetail(ctx context.Context, req *typespec.VulScanVulDetailReq, resp *typespec.VulScanVulDetailResp) error {
	// 配置请求参数
	var param httpclients.VulScanVulDetailReq
	param.Id = req.Id

	// 请求查询
	remoteRes, err := httpclients.VulScanVulDetail(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}

	// 组织返回的数据
	resp.ID = remoteRes.Data.ID
	resp.TaskID = remoteRes.Data.TaskID
	resp.TargetID = remoteRes.Data.TargetID
	resp.Name = remoteRes.Data.Name
	resp.Cve = remoteRes.Data.Cve
	resp.Port = remoteRes.Data.Port
	resp.Description = remoteRes.Data.Description
	resp.Solution = remoteRes.Data.Solution
	resp.Parameter = remoteRes.Data.Parameter
	resp.Detail = remoteRes.Data.Detail
	resp.Risk = remoteRes.Data.Risk
	resp.RiskName = remoteRes.Data.RiskName
	resp.CreateTime = remoteRes.Data.CreateTime
	resp.UpdateTime = remoteRes.Data.UpdateTime
	resp.Ip = remoteRes.Data.Ip
	resp.Cwe = remoteRes.Data.Cwe
	resp.Vendor = remoteRes.Data.Vendor
	resp.Product = remoteRes.Data.Product
	resp.Version = remoteRes.Data.Version
	resp.Cpes = remoteRes.Data.Cpes
	resp.CvssVersion = remoteRes.Data.CvssVersion
	resp.CvssVector = remoteRes.Data.CvssVector
	resp.PublishDate = remoteRes.Data.PublishDate
	resp.ExploitabilityScore = remoteRes.Data.ExploitabilityScore
	resp.References = remoteRes.Data.References
	return nil
}

// CveList 漏洞列表
func (v *VulScan) CveList(ctx context.Context, req *typespec.VulScanCveListReq, resp *typespec.VulScanCveListResp) error {
	// 配置请求参数
	var param httpclients.VulScanCveListReq
	param.Page = req.Page
	param.Size = req.Size
	param.Search = req.Search

	// 请求查询
	remoteRes, err := httpclients.VulScanCveList(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}

	// 组织返回的数据
	resp.Total = remoteRes.Data.Total
	for _, v := range remoteRes.Data.List {
		resp.List = append(resp.List, typespec.VulScanCveListResItem{
			Id:                v.Id,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
			DeletedAt:         v.DeletedAt,
			Cve:               v.Cve,
			Cwe:               v.Cwe,
			TitleZh:           v.TitleZh,
			DescriptionMain:   v.DescriptionMain,
			DescriptionMainZh: v.DescriptionMainZh,
			Descriptions:      v.Descriptions,
			Vendor:            v.Vendor,
			Product:           v.Product,
			Severity:          v.Severity,
			SeverityName:      v.SeverityName,
			PublishedDate:     v.PublishedDate,
			BaseCvssv2Score:   v.BaseCvssv2Score,
		})
	}
	return nil
}

// CveDetail cve详情
func (v *VulScan) CveDetail(ctx context.Context, req *typespec.VulScanCveDetailReq, resp *typespec.VulScanCveDetailResp) error {
	// 配置请求参数
	var param httpclients.VulScanCveDetailReq
	param.Id = req.Id

	// 请求查询
	remoteRes, err := httpclients.VulScanCveDetail(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}

	resp.Id = remoteRes.Data.Id
	resp.CreatedAt = remoteRes.Data.CreatedAt
	resp.UpdatedAt = remoteRes.Data.UpdatedAt
	resp.DeletedAt = remoteRes.Data.DeletedAt
	resp.Cve = remoteRes.Data.Cve
	resp.Cwe = remoteRes.Data.Cwe
	resp.TitleZh = remoteRes.Data.TitleZh
	resp.DescriptionMain = remoteRes.Data.DescriptionMain
	resp.DescriptionMainZh = remoteRes.Data.DescriptionMainZh
	resp.Descriptions = remoteRes.Data.Descriptions
	resp.Vendor = remoteRes.Data.Vendor
	resp.Product = remoteRes.Data.Product
	resp.Severity = remoteRes.Data.Severity
	resp.SeverityName = remoteRes.Data.SeverityName
	resp.CpeConfigurations = remoteRes.Data.CpeConfigurations
	resp.CvssVersion = remoteRes.Data.CvssVersion
	resp.CvssVectorString = remoteRes.Data.CvssVectorString
	resp.AccessVector = remoteRes.Data.AccessVector
	resp.AccessComplexity = remoteRes.Data.AccessComplexity
	resp.Authentication = remoteRes.Data.Authentication
	resp.ConfidentialityImpact = remoteRes.Data.ConfidentialityImpact
	resp.IntegrityImpact = remoteRes.Data.IntegrityImpact
	resp.AvailabilityImpact = remoteRes.Data.AvailabilityImpact
	resp.BaseCvssv2Score = remoteRes.Data.BaseCvssv2Score
	resp.PublishedDate = remoteRes.Data.PublishedDate
	resp.Solution = remoteRes.Data.Solution
	resp.References = remoteRes.Data.References
	return nil
}

// TaskOverview 任务概览
func (v *VulScan) TaskOverview(ctx context.Context, req *typespec.VulScanTaskOverviewReq, resp *typespec.VulScanTaskOverviewResp) error {
	// 配置请求参数
	var param httpclients.VulScanTaskOverviewReq
	param.Id = req.Id

	// 请求查询
	remoteRes, err := httpclients.VulScanTaskOverview(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}
	// 组织返回的数据
	resp.TaskName = remoteRes.Data.TaskName
	resp.Risk = remoteRes.Data.Risk
	resp.RiskName = remoteRes.Data.RiskName
	resp.TargetRisk = remoteRes.Data.TargetRisk
	resp.TargetNum = remoteRes.Data.TargetNum
	resp.VulRisk = remoteRes.Data.VulRisk
	resp.CreateTime = remoteRes.Data.CreateTime
	resp.UpdateTime = remoteRes.Data.UpdateTime
	resp.Ports = remoteRes.Data.Ports
	return nil
}

// TaskState 任务状态
func (v *VulScan) TaskState(ctx context.Context, req *typespec.VulScanTaskGetStateReq, resp *typespec.VulScanTaskGetStateResp) error {
	// 配置请求参数
	var param httpclients.VulScanTaskStateReq
	param.Id = req.Id

	// 请求查询
	remoteRes, err := httpclients.VulScanTaskState(ctx, param)
	if err != nil {
		return err
	}
	if remoteRes.Code != 200 {
		// 请求出错
		return errors.New(remoteRes.Msg)
	}
	// 组织返回的数据
	resp.Status = remoteRes.Data.Status
	resp.StatusName = remoteRes.Data.StatusName
	return nil
}
