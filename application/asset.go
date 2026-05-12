package application

import (
	"context"
	"errors"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"
)

type AssetApp struct{}

// 资产变化统计
func (a *AssetApp) GetAssetEnums(ctx context.Context, resp *typespec.GetAssetEnumsResp) error {
	resp.AssetChangeType = enums.AssetEnum.AssetChangeType()
	resp.AssetIsLiveType = enums.AssetEnum.AssetIsLiveType()
	resp.DeviceWeight = enums.AssetEnum.DeviceWeight()
	resp.LoginProtocol = enums.LoginProtocolTools.ProtocolList()
	resp.TrustLevel = enums.AssetEnum.TrustLevel()
	return nil
}

// 近期变化资产列表
func (a *AssetApp) GetChangeAssetList(ctx context.Context, req *typespec.GetChangeAssetListReq, resp *typespec.GetChangeAssetListResp) error {
	var (
		assetSrv services.Asset
		openPort = make(map[string][]string, 0)
	)
	assetRes, assetIps, total := assetSrv.GetChangeAssetList(ctx, enums.AssetIsIgnoreNo, req.Search, req.AssetChangesType, req.IsLive, req.UpdateTime, req.Page, req.Size)
	if len(assetIps) > 0 {
		openPort = assetSrv.GetAssetPort(ctx, assetIps) //开放端口
	}
	resp.Total = total
	resp.List = make([]typespec.GetChangeAssetListRespList, 0)
	for i := 0; i < len(assetRes); i++ {
		var tmp = typespec.GetChangeAssetListRespList{
			Id:                   assetRes[i].ID,
			Ip:                   assetRes[i].IP,
			OperateSystem:        assetRes[i].OperateSystem,
			PortOpen:             "",
			AssetChangesType:     assetRes[i].AssetChangesType,
			AssetChangesTypeName: enums.AssetEnum.GetAssetChangeTypeName(assetRes[i].AssetChangesType),
			IsLive:               assetRes[i].Islive,
			IsLiveName:           enums.AssetEnum.GetAssetIsLiveTypeName(assetRes[i].Islive),
			UpdateTime:           assetRes[i].UpdateTime.Format(enums.TimeLayout),
		}
		if v, ok := openPort[assetRes[i].IP]; ok {
			tmp.PortOpen = strings.Join(v, ",")
		}
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// 近期变化资产删除
func (a *AssetApp) ChangeAssetDel(ctx context.Context, req *typespec.ChangeAssetDelReq) error {
	ids := strings.Split(req.AssetIds, ",")
	var assetSrv services.Asset
	return assetSrv.UpdateAssetIgnore(ctx, ids)
}

// AssetVulFindSync 资产漏洞信息 - 同步
func (a *AssetApp) AssetVulFindSync(ctx context.Context, ids string) error {
	var (
		asset         services.Asset
		taskTargetSrv services.TaskTarget // 获取渗透目标数据
		taskVulSrv    services.TaskVul    // 获取目标漏洞数据
	)
	idSlice := strings.Split(ids, ",")
	targetIds := make([]int, 0)
	for _, item := range idSlice {
		idI, err := strconv.Atoi(item)
		if err == nil {
			targetIds = append(targetIds, idI)
		}
	}
	// 获取资产信息
	taskTargetMap := taskTargetSrv.GetTargetByIds(ctx, targetIds, enums.TargetStatusFinish)
	if len(taskTargetMap) == 0 {
		return errors.New("未查询到任何已结束的目标资产")
	}
	assetIDsByTargetID := make(map[int]int, 0)
	for _, item := range taskTargetMap {
		assetModel := mysqls.Asset{
			AssetGroupID:  1,
			IP:            utils.GetHostname(item.TargetURL),
			OperateSystem: item.OpSys,
			FilingLevel:   enums.AssetFilingLevelZero,
			TargetIds:     strconv.Itoa(item.ID),
			Status:        enums.AssetStatusSync,
			IsIgnore:      enums.AssetIsIgnoreNo,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}
		assetRes, err := assetModel.GetAssetByIP(ctx)
		if err != nil {
			return err
		}
		if assetRes.ID != 0 {
			assetIDsByTargetID[item.ID] = assetRes.ID
			continue
		}
		assetID, err := assetModel.AddAsset(ctx)
		if err != nil {
			return err
		}
		assetIDsByTargetID[item.ID] = assetID
	}
	if len(assetIDsByTargetID) == 0 {
		return nil
	}
	// 2 资产漏洞数据同步
	taskVulMap := taskVulSrv.AllByTargetIds(ctx, targetIds)
	insertAssetVulDatas := make([]mysqls.Assetvul, 0)
	for targetId, assetId := range assetIDsByTargetID {
		for _, taskVul := range taskVulMap[targetId] {
			insertAssetVulDatas = append(insertAssetVulDatas, mysqls.Assetvul{
				TaskVulID:      taskVul.ID,
				AssetID:        assetId,
				TaskID:         taskVul.TaskID,
				TargetID:       taskVul.TargetID,
				TargetURL:      taskVul.TargetUrl,
				Pocname:        taskVul.Pocname,
				Name:           taskVul.Name,
				Class:          taskVul.Class,
				Type:           taskVul.Type,
				Risk:           taskVul.Risk,
				Location:       taskVul.Location,
				Status:         taskVul.Status,
				TestStatus:     taskVul.TestStatus,
				ExploitImpact:  taskVul.ExploitImpact,
				VulID:          taskVul.VulID,
				Description:    taskVul.Description,
				FixSuggest:     taskVul.FixSuggest,
				PublishedTime:  taskVul.PublishedTime,
				AffectRange:    taskVul.AffectRange,
				TargetResultID: taskVul.TargetResultID,
				VulNumber:      taskVul.VulNumber,
				VulAddress:     taskVul.VulAddress,
				RefURL:         taskVul.RefUrl,
				Cvss:           taskVul.Cvss,
				VulResult:      taskVul.VulResult,
				VulParam:       taskVul.VulParam,
				VerMsg:         taskVul.VerMsg,
				IsReplace:      enums.AssetIsReplaceN,
				CreateTime:     time.Now(),
				UpdateTime:     time.Now(),
			})
		}
	}
	if len(insertAssetVulDatas) != 0 {
		if err := asset.MultipartInsertAssetVul(ctx, &insertAssetVulDatas); err != nil {
			return err
		}
	}
	return nil
}
