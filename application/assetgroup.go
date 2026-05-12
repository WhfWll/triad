package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/file"
	"smart/tools/network"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

type AssetGroupApp struct{}

// AssetTree 资产树整体结构
func (a *AssetGroupApp) AssetTree(ctx context.Context, req *typespec.AssetTreeOverallReq, resp *typespec.AssetTreeOverallResp) error {
	var (
		assetTreeSrv services.AssetGroup
		assetRes     []mysqls.Asset
	)
	resp.List = make([]typespec.AssetTreeOverallRespItems, 0)
	assetGroupRes := assetTreeSrv.GetAllAssetGroup(ctx)                             //查询所有的资产组
	tmpAssetGroupRes := assetTreeSrv.HandAssetGroupByAsset(assetGroupRes, assetRes) //合并资产和资产组数据结构
	tmpAssetGroupRes = assetTreeSrv.SortAssetGroupSlice(tmpAssetGroupRes)           //排序
	resp.List = assetTreeSrv.SortAssetAssetGroup(tmpAssetGroupRes, 0)               //组装资产和资产组tree结构
	return nil
}

// AssetGroup 获取资产组结构
func (a *AssetGroupApp) AssetGroup(ctx context.Context, resp *typespec.AssetGroupResp) error {
	var assetSrv services.AssetGroup
	resp.List = assetSrv.SortAssetGroup(assetSrv.GetAllAssetGroup(ctx), 0)
	return nil
}

// AssetGroupAdd 新增资产组
func (a *AssetGroupApp) AssetGroupAdd(ctx context.Context, req *typespec.AssetTreeAddReq) error {
	var assetGroupSrv services.AssetGroup
	level := 0
	if req.Pid != 0 {
		fatherGroup, _ := assetGroupSrv.GetAssetGroupInfo(ctx, req.Pid)
		if fatherGroup.ID == 0 {
			return errors.New("上级组不存在")
		}
		level = fatherGroup.Level + 1
	}
	// 资产组层级最多6层
	if level > 6 {
		return errors.New("资产组层级超过上限！请更换上级。")
	}
	if req.Name == "" {
		return errors.New("资产组名不能为空")
	}
	// 验证IP准确性
	if network.IsSingleIP(req.Name) {
		return errors.New("不能使用IP做为资产组名称")
	}
	// 如果所属上级不同 可以名字相同
	if req.Pid != 0 {
		assetGroupInfo, _ := assetGroupSrv.GetAssetGroupInfoByNamePid(ctx, req.Name, req.Pid)
		if assetGroupInfo.ID != 0 {
			return errors.New("资产组名重复！ 请检查后重新添加！")
		}
	} else {
		assetGroupInfo, _ := assetGroupSrv.GetAssetGroupInfoByName(ctx, req.Name)
		if assetGroupInfo.ID != 0 {
			return errors.New("资产组名重复！ 请检查后重新添加！")
		}
	}
	// 默认资产组不能嵌套资产组
	defaultAssetGroupID := assetGroupSrv.GetAssetGroupDefaultGroupID(ctx)
	if req.Pid == defaultAssetGroupID {
		return errors.New("默认资产组不能嵌套资产组!")
	}
	if err := assetGroupSrv.AssetGroupAdd(ctx, req.Name, req.Remark, req.Pid, level); err != nil {
		return errors.New("任务组创建失败:" + err.Error())
	}
	return nil
}

// AssetGroupEdit 资产组编辑
func (a *AssetGroupApp) AssetGroupEdit(ctx context.Context, req *typespec.AssetGroupEditReq) error {
	var (
		assetGroupSrv services.AssetGroup
		level         int
	)
	level = 0
	if req.Pid != 0 {
		group, _ := assetGroupSrv.GetAssetGroupInfo(ctx, req.Pid)
		if group.ID == 0 {
			return errors.New("归属组不存在")
		}
		level = group.Level + 1
	}
	// 资产组层级最多6层
	if level > 6 {
		return errors.New("资产组层级超过上限！请更换上级。")
	}
	if req.Name == "" {
		return errors.New("资产组名不能为空")
	}
	if req.Pid == req.Id {
		return errors.New("父节点冲突 请重新选择！")
	}
	assetGroupInfo, _ := assetGroupSrv.GetAssetGroupInfoByName(ctx, req.Name)
	if assetGroupInfo.ID != req.Id && assetGroupInfo.ID != 0 {
		return errors.New("资产组名重复！ 请检查后重新修改！")
	}
	err := assetGroupSrv.UpdateAssetGroup(ctx, req.Id, req.Pid, req.Name, req.Remark, level)
	if err != nil {
		return err
	}
	return nil
}

// AssetGroupDel 资产组删除
func (a *AssetGroupApp) AssetGroupDel(ctx context.Context, req *typespec.AssetDeleteReq) error {
	var (
		assetGroupSrv services.AssetGroup
		assetId       []int
	)
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	groupId := strings.TrimSpace(req.GroupIds)
	defaultAssetGroupID := assetGroupSrv.GetAssetGroupDefaultGroupID(ctx)
	deleteGroupIdSlice := make([]int, 0)
	if groupId != "" {
		groupIdStringSlice := strings.Split(groupId, ",")
		for _, assetGroupID := range groupIdStringSlice {
			groupIdInt, _ := strconv.Atoi(assetGroupID)
			// ！!!不能删除默认组!!!
			if groupIdInt == defaultAssetGroupID {
				return errors.New("包含默认资产组,请重新选择!")
			}
			assetGroupAssetIds := assetGroupSrv.RecursionGetAllSubAssetGroupId(ctx, groupIdInt)
			deleteGroupIdSlice = append(deleteGroupIdSlice, assetGroupAssetIds...)
		}
		if len(deleteGroupIdSlice) > 0 {
			assets := assetGroupSrv.AllAssetByGroupIds(dCtx, deleteGroupIdSlice)
			for _, item := range assets {
				assetId = append(assetId, item.ID)
			}
		}
	}
	// 删除资产组
	if len(deleteGroupIdSlice) > 0 {
		err := assetGroupSrv.DeleteByGroupIds(dCtx, deleteGroupIdSlice)
		if err != nil {
			return err
		}
	}
	// 将资产移动到默认组中
	if len(assetId) > 0 {
		err := assetGroupSrv.UpdateAssetGroupID(dCtx, assetId, defaultAssetGroupID)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// AssetDel 资产删除
func (a *AssetGroupApp) AssetDel(ctx context.Context, req *typespec.AssetDeleteReq) error {
	var assetGroupSrv services.AssetGroup
	assetIdIntSlice := make([]int, 0)
	assetId := strings.TrimSpace(req.AssetIds)
	if assetId != "" {
		assetIdStringSlice := strings.Split(assetId, ",")
		for _, item := range assetIdStringSlice {
			assetIdInt, err := strconv.Atoi(item)
			if err != nil {
				return errors.New("资产ID中存在非数值类数据，请确认")
			}
			assetIdIntSlice = append(assetIdIntSlice, assetIdInt)
		}
	}
	if len(assetIdIntSlice) > 0 {
		err := assetGroupSrv.DeleteAssetByAssetIds(ctx, assetIdIntSlice)
		if err != nil {
			return err
		}
	}
	return nil
}

// SelectAllAssetList 全选所有资产
func (a *AssetGroupApp) SelectAllAssetList(ctx context.Context, req *typespec.SelectAllAssetReq, res *typespec.SelectAllAssetRes) error {
	var (
		assetGroupSrv   services.AssetGroup
		groupIdIntSlice []int
	)
	// 获取资产组下所有子资产信息 和 路径
	subAssetGroupList, _, err := assetGroupSrv.GetSubAssetGroupList(ctx, req.GroupId)
	if err != nil {
		return err
	}
	groupIdIntSlice = []int{req.GroupId}
	for _, v := range subAssetGroupList {
		groupIdIntSlice = append(groupIdIntSlice, v.ID)
	}
	// 获取资产相关信息
	allAsset, _ := assetGroupSrv.AllAssetByGroupIdsAndSearch(ctx, groupIdIntSlice, req.AssetIP, req.Port, req.Service, req.SystemOp, req.VulName, req.Domain, req.Finger, req.Tags, req.AssetType, req.IsCloudHost, req.AssetRisk, req.FillingLevel, 1, 10000)
	for _, v := range allAsset {
		res.AssetIDs = append(res.AssetIDs, v.ID)
		res.AssetIPs = append(res.AssetIPs, v.IP)
	}
	return nil
}

// AssetList 资产列表
func (a *AssetGroupApp) AssetList(ctx context.Context, req *typespec.AssetListReq, res *typespec.AssetListRes) error {
	var (
		assetGroupSrv   services.AssetGroup
		assetPortSrv    services.AssetPort
		groupIdIntSlice []int
	)
	// 获取资产组下所有子资产信息 和 路径
	subAssetGroupList, _, err := assetGroupSrv.GetSubAssetGroupList(ctx, req.GroupId)
	if err != nil {
		return err
	}
	groupIdIntSlice = []int{req.GroupId}
	for _, v := range subAssetGroupList {
		groupIdIntSlice = append(groupIdIntSlice, v.ID)
	}
	// 获取资产相关信息
	allAsset, count := assetGroupSrv.AllAssetByGroupIdsAndSearch(ctx, groupIdIntSlice, req.AssetIP, req.Port, req.Service, req.SystemOp, req.VulName, req.Domain, req.Finger, req.Tags, req.AssetType, req.IsCloudHost, req.AssetRisk, req.FillingLevel, req.Page, req.Size)
	res.AssetsInfo = make([]typespec.AssetsInfo, 0)
	for _, v := range allAsset {
		if v.IP == "" {
			continue
		}
		riskLevel, riskCountArr := services.GetAssetRiskLevel(ctx, v.IP)
		// 获取资产端口信息
		openPort := assetPortSrv.GetAssetPortForTarget(ctx, v.IP)
		// 获取资产组信息
		assetGroupInfo, _ := assetGroupSrv.GetAssetGroupInfo(ctx, v.AssetGroupID)
		res.AssetsInfo = append(res.AssetsInfo, typespec.AssetsInfo{
			AssetID:               v.ID,
			IP:                    v.IP,
			System:                v.OperateSystem,
			OpenPort:              openPort,
			Location:              v.Location,
			AssetGroupName:        assetGroupInfo.Name,
			ResponsibleDepartment: v.ResponsibleDepartment,
			DeviceWeight:          enums.AssetEnum.GetDeviceWeightName(v.DeviceWeight),
			TrustLevel:            enums.AssetEnum.GetTrustLevelName(v.TrustLevel),
			AssetType:             enums.AssetEnum.GetAssetTypeName(v.AssetType),
			AssetRiskName:         enums.GetAssetRisk(riskLevel),
			VulStatics: typespec.VulStatics{
				DeadlyVul:     riskCountArr[1],
				HighRiskVul:   riskCountArr[2],
				MediumRiskVul: riskCountArr[3],
				LowRiskVul:    riskCountArr[4],
			},
			TestTime: v.UpdateTime.Format(enums.TimeLayout),
		})
	}
	res.Count = int(count)
	return nil
}

// staticsAssets 统计资产信息
func staticsAssets(assetList []mysqls.Asset) (map[int]int, []int) {
	countMap := make(map[int]int)
	var uniqueIDs []int
	for _, asset := range assetList {
		ctx := context.Background()
		level, _ := services.GetAssetRiskLevel(ctx, asset.IP)
		countMap[level]++
		uniqueIDs = append(uniqueIDs, asset.ID)
	}
	return countMap, removeDuplicates(uniqueIDs)
}

// staticsAssetVul 统计资产漏洞信息
func staticsAssetVul(assetList []mysqls.Assetvul) (map[int]int, map[string]int) {
	countMap := make(map[int]int)
	countTypeMap := make(map[string]int)
	for _, asset := range assetList {
		countMap[asset.Risk]++
		countTypeMap[enums.ToolsVulnerabilityEnum.GetClassEnum(asset.Type)]++
	}
	return countMap, countTypeMap
}

// removeDuplicates 切片去重
func removeDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

// checkAssetRisk 检查资产组风险
func checkAssetRisk(staticsAssetsMap map[int]int) string {
	if _, ok := staticsAssetsMap[enums.HighRiskAsset]; ok {
		return "高危"
	} else if _, ok := staticsAssetsMap[enums.MiddleRiskAsset]; ok {
		return "中危"
	} else if _, ok := staticsAssetsMap[enums.LowRiskAsset]; ok {
		return "低危"
	} else if _, ok := staticsAssetsMap[enums.SafeAsset]; ok {
		return "安全"
	}
	return "未知"
}

// AddAsset 新增资产
func (a *AssetGroupApp) AddAsset(ctx context.Context, req *typespec.AssetAddReq) error {
	var (
		assetGroupSrv services.AssetGroup
		assetConnSrv  services.AssetConnectionService
	)
	// 验证IP准确性
	if !network.IsSingleIP(req.IP) && !network.IsValidDomainFormat(req.IP) {
		return errors.New("资产IP/域名 格式错误")
	}
	info := assetGroupSrv.OpAssetInfoJson(req.BaseSoftwareName, req.BaseSoftwareVersion, req.BaseHardwareName, req.Purpose, req.EquipmentForm, req.DeploymentLocation, req.SystemName, req.SystemAdmin, req.SystemOp)
	assetInfo := assetGroupSrv.GetNoLiveAssetByIp(ctx, req.IP)
	if assetInfo.ID != 0 {
		assetGroupSrv.PhysicsDeleteAssetByAssetId(ctx, []int{assetInfo.ID})
	}
	assetLiveInfo := assetGroupSrv.GetLiveAssetByIp(ctx, req.IP)
	if assetLiveInfo.ID != 0 {
		return errors.New("新增资产IP已存在, 请重新输入!")
	}
	assetID, err := assetGroupSrv.AssetAdd(ctx, req.IP, req.Name, req.OpSys, "", req.Tags, info, req.ResponsibleDepartment, req.AssetGroupID, 0, req.EqualProtectionLevel, req.IsCloudHost, req.DeviceWeight, req.TrustLevel)
	if err != nil {
		return err
	}
	if req.Protocol != 0 || req.User != "" || req.Password != "" || req.Port != 0 {
		if err := assetConnSrv.CreateAssetConnection(ctx, assetID, req.Port, req.Protocol, req.IP, req.User, req.Password); err != nil {
			log.Println("CreateAssetConnection", err.Error())
		}
	}
	return nil
}

// EditAsset 修改资产
func (a *AssetGroupApp) EditAsset(ctx context.Context, req *typespec.AssetEditReq) error {
	var (
		assetGroupSrv services.AssetGroup
		assetConnSrv  services.AssetConnectionService
		assetPortSrv  services.AssetPort
	)
	assetInfo, err := assetGroupSrv.GetAssetByID(ctx, req.ID)
	if err != nil || assetInfo.ID == 0 {
		return errors.New("找不到数据...")
	}
	if assetInfo.IP != req.IP {
		assetPortRes := assetPortSrv.GetAssetPortByIp(ctx, assetInfo.IP)
		if len(assetPortRes) > 0 {
			return errors.New("该ip存在端口数据，无法修改ip字段,请删除后再添加...")
		}
	}
	// 验证IP准确性
	if !network.IsSingleIP(req.IP) && !network.IsValidDomainFormat(req.IP) {
		return errors.New("资产IP/域名 格式错误")
	}
	info := assetInfo.Info
	if req.BaseSoftwareName != "" || req.BaseSoftwareVersion != "" || req.Purpose != "" || req.EquipmentForm != "" || req.DeploymentLocation != "" || req.SystemOp != "" || req.SystemAdmin != "" {
		info = assetGroupSrv.OpAssetInfoJson(req.BaseSoftwareName, req.BaseSoftwareVersion, req.BaseHardwareName, req.Purpose, req.EquipmentForm, req.DeploymentLocation, req.SystemName, req.SystemAdmin, req.SystemOp)
	}
	if err := assetGroupSrv.UpdateAssetById(ctx, req.IP, req.OpSys, req.IPSegment, info, req.Tags, req.ResponsibleDepartment, req.Name, req.AssetGroupID, req.AssetType, req.EqualProtectionLevel, req.ID, req.IsCloudHost); err != nil {
		return err
	}
	if req.Protocol != 0 || req.User != "" || req.Password != "" || req.Port != 0 {
		assetConnSrv.UpdateAssetConnection(ctx, req.IP, req.ID, req.Port, req.Protocol, req.User, req.Password)
	}
	return nil
}

// GetAssetOpenPort 获取资产开放端口
func (a *AssetGroupApp) GetAssetOpenPort(ctx context.Context, ip string) []typespec.ScanTaskPortListList {
	var (
		targetSrv services.TaskTarget // 获取渗透目标数据
		ids       []int
	)
	target, _ := targetSrv.TargetListByTargetUrl(ctx, ip)
	targetRes, _ := targetSrv.GetTargetByTaskId(ctx, target.TaskID)
	for i := 0; i < len(targetRes); i++ {
		if targetRes[i].TargetURL == "" || !strings.Contains(targetRes[i].TargetURL, ip) {
			continue
		}
		ids = append(ids, targetRes[i].ID)
	}
	taskResultRes, _ := targetSrv.GetAssetTargetOpenPort(ctx, ids)
	portMap := make(map[int]typespec.ScanTaskPortListList)
	for i := 0; i < len(taskResultRes); i++ {
		// 1_1 1_2
		var parsed struct {
			Port        string `json:"port"`
			Protocol    string `json:"protocol"`
			Service     string `json:"service"`
			Component   string `json:"component"`
			Title       string `json:"title"`
			Fingerprint string `json:"fingerprint"`
		}
		if err := json.Unmarshal([]byte(taskResultRes[i].JSONResult), &parsed); err != nil {
			continue
		}
		port, _ := strconv.Atoi(parsed.Port)
		newEntry := typespec.ScanTaskPortListList{
			Id:         taskResultRes[i].ID,
			Port:       port,
			Protocol:   parsed.Protocol,
			Service:    parsed.Service,
			Assembly:   parsed.Component,
			Title:      parsed.Title,
			Remark:     "",
			CreateTime: taskResultRes[i].CreateTime.Format(enums.TimeLayout),
		}
		// 如果已有该端口，进行字段补全
		if existing, ok := portMap[port]; ok {
			portMap[port] = mergePortEntries(existing, newEntry)
		} else {
			portMap[port] = newEntry
		}
	}
	var assetPortList []typespec.ScanTaskPortListList
	for _, v := range portMap {
		assetPortList = append(assetPortList, v)
	}
	return assetPortList
}

// AssetDetail 资产详情
func (a *AssetGroupApp) AssetDetail(ctx context.Context, req *typespec.AssetDetailReq, resp *typespec.AssetDetail) error {
	var (
		assetGroupSrv services.AssetGroup
		assetPortList []typespec.ScanTaskPortListList
		riskApp       RiskManageApp
		riskVulResp   typespec.VulRiskListResp
	)
	assetInfo, err := assetGroupSrv.GetAssetByID(ctx, req.ID)
	if assetInfo.IP == "" {
		return errors.New("非法参数")
	}
	if err != nil {
		return err
	}
	if assetInfo.ID == 0 {
		return errors.New("查询资产信息不存在。")
	}
	if req.SelectType == 1 {
		resp.Info = a.GetAssetOpenPort(ctx, assetInfo.IP)
	}
	if req.SelectType == 2 {
		if req.Page == 0 || req.Size == 0 {
			return errors.New("参数错误")
		}
		// 漏洞信息
		riskApp.VulRiskListByIP(ctx, &typespec.VulRiskListReq{
			IP:   assetInfo.IP,
			Page: req.Page,
			Size: req.Size,
		}, &riskVulResp)
		resp.VulInfo = riskVulResp
	}
	// 管理信息
	assetGroupInfo, _ := assetGroupSrv.GetAssetGroupInfo(ctx, assetInfo.AssetGroupID)
	groupName := assetGroupInfo.Name
	if name, _ := assetGroupSrv.CheckAssetGroupName(ctx, assetGroupInfo.ID); name != "" {
		groupName = name
	}
	jsonInfo := assetGroupSrv.GetAssetJsonInfo(assetInfo.Info)
	riskLevel, riskCountArr := services.GetAssetRiskLevel(ctx, assetInfo.IP)
	// 获取资产连接信息
	var assetConnSrv services.AssetConnectionService
	assetConnectionList, _ := assetConnSrv.GetConnectionsByIPWithDecryptedPassword(ctx, assetInfo.IP)
	var connections []typespec.AssetConnectionDTO
	for _, c := range assetConnectionList {
		connections = append(connections, typespec.AssetConnectionDTO{
			Username: c.Username,
			Password: c.Password,
			Port:     c.Port,
			Protocol: enums.ConnEnum.GetConnMethodEnum(c.Protocol),
		})
	}
	resp.ManageInfo = typespec.AssetInfo{
		IP:                      assetInfo.IP,
		AssetName:               assetInfo.Name,
		RiskLevelStr:            enums.GetAssetRisk(riskLevel),
		Score:                   services.CalculateAssetScore(ctx, riskCountArr, riskLevel, assetInfo.IP),
		AssetGroupID:            assetInfo.AssetGroupID,
		AssetGroupName:          groupName,
		OpSys:                   assetInfo.OperateSystem,
		ResponsibleDepartment:   assetInfo.ResponsibleDepartment,
		IPSegment:               assetInfo.IPSegment,
		AssetType:               assetInfo.AssetType,
		BaseSoftwareName:        jsonInfo.SoftwareName,
		BaseSoftwareVersion:     jsonInfo.SoftwareVersion,
		BaseHardwareName:        jsonInfo.HardwareName,
		Purpose:                 jsonInfo.Purpose,
		EquipmentForm:           jsonInfo.EquipmentForm,
		DeploymentLocation:      jsonInfo.DeploymentLocation,
		EqualProtectionLevel:    assetInfo.FilingLevel,
		EqualProtectionLevelStr: enums.AssetEnum.GetAssetFilingLevel(assetInfo.FilingLevel),
		DeviceWeight:            enums.AssetEnum.GetDeviceWeightName(assetInfo.DeviceWeight),
		TrustLevel:              enums.AssetEnum.GetTrustLevelName(assetInfo.TrustLevel),
		Location:                assetInfo.Location,
		IsCloudHost:             assetInfo.IsCloudHost,
		SystemName:              jsonInfo.SystemName,
		SystemAdmin:             jsonInfo.SystemAdmin,
		SystemOp:                jsonInfo.SystemOperate,
		Tags:                    assetInfo.Tags,
		PortList:                assetPortList,
		Connections:             connections,
	}
	return nil
}

// GetAssetGroupEnums 资产组枚举
func (a *AssetGroupApp) GetAssetGroupEnums(ctx context.Context, resp *typespec.GetAssetGroupEnumsResp) error {
	resp.AssetType = enums.AssetEnum.AssetType()
	resp.FillingLevel = enums.AssetEnum.GetAssetFilingLevelEnumArray()
	resp.AssetRisk = enums.AssetEnum.AssetRiskType()
	resp.DeviceWeight = enums.AssetEnum.DeviceWeight()
	resp.TrustLevel = enums.AssetEnum.TrustLevel()
	var assetGroup services.AssetGroup
	resp.AssetGroup = assetGroup.AssetGroupEnums(ctx)
	resp.LoginProtocol = enums.ConnEnum.GetConnMethodEnums()
	return nil
}

// ReadXlsx 读取文件信息
func (a *AssetGroupApp) ReadXlsx(ctx context.Context, c *gin.Context) (string, error) {
	var (
		fileTools     file.Excel
		assetGroupSrv services.AssetGroup
		repeatIp      string
	)
	_, lists, err := fileTools.ImportData(c)
	if err != nil {
		return "", err
	}
	// IP 归属资产组 资产名 操作系统  管理员   标签 部门
	for _, v := range lists {
		// 如果列数不足7列，进行补全
		for len(v) < 7 {
			v = append(v, "")
		}

		ip := v[0] // 资产IP
		if ip == "" {
			continue
		}
		assetGroupID := enums.DefaultAssetGroup
		assetGroupName := v[1] // 归属资产组
		if assetGroupName != "" {
			assetGroupID = getAssetGroupInfo(ctx, assetGroupName)
		}
		// 通过名称获取资产组ID 如果没有就创建资产组
		name := v[2]                  // 资产名称
		sysOP := v[3]                 // 操作系统
		systemAdmin := v[4]           // 管理员
		tags := v[5]                  // 标签
		responsibleDepartment := v[6] // 部门
		//equipmentForm := v[3]                                               // 设备形态
		//deviceWeight := enums.AssetEnum.GetDeviceWeightNameID(v[4])         // 设备权重
		//equalProtectionLevel := enums.AssetEnum.GetAssetFilingLevelID(v[5]) // 等保等级
		assetType := 0
		jsonInfo := assetGroupSrv.OpAssetInfoJson("", "", "", "", "", "", "", systemAdmin, "")
		// 如果ip存在且被删除过 恢复存活状态
		assetInfo := assetGroupSrv.GetNoLiveAssetByIp(ctx, ip)
		if assetInfo.ID != 0 {
			assetGroupSrv.PhysicsDeleteAssetByAssetId(ctx, []int{assetInfo.ID})
		}
		_, err = assetGroupSrv.AssetAdd(ctx, ip, name, sysOP, "", tags, jsonInfo, responsibleDepartment, assetGroupID, assetType, 0, 0, 0, 0)
		if err != nil {
			repeatIp += ip + ","
		}
	}
	if repeatIp == "" {
		return "", nil
	}
	msg := ""
	if strings.HasSuffix(repeatIp, ",") {
		repeatIp = repeatIp[:len(repeatIp)-1]
		msg = "导入资产中存在重复资产：" + repeatIp
	}
	return msg, nil
}

func getAssetGroupInfo(ctx context.Context, assetGroupName string) int {
	var agModel mysqls.Assetgroup
	assetGroupID, _ := agModel.GetOrCreateAssetGroupByPath(ctx, assetGroupName)
	return assetGroupID
	//if !strings.Contains(assetGroupName, "/") {
	//assGroupInfo, _ := assetGroup.GetAssetGroupInfoByName(ctx, assetGroupName)
	//if assGroupInfo.ID != 0 {
	//	return assGroupInfo.ID
	//} else {
	//	id, _ := assetGroup.AssetGroupAddReturnID(ctx, assetGroupName)
	//	return id
	//}
	//} else {
	//	assetGroupID, _ := agModel.GetOrCreateAssetGroupByPath(ctx, assetGroupName)
	//	return assetGroupID
	//}
	//return 0
}

// ExportXlsx 导出文件信息
func (a *AssetGroupApp) ExportXlsx(ctx context.Context, c *gin.Context) error {
	var (
		fileTools file.Excel
		//assetGroupSrv services.AssetGroup
	)
	_, lists, err := fileTools.ImportData(c)
	if err != nil {
		return err
	}
	for _, v := range lists {
		assetGroupID, _ := strconv.Atoi(v[1])
		assetType, _ := strconv.Atoi(v[4])
		equalProtectionLevel, _ := strconv.Atoi(v[10])
		fmt.Println(assetGroupID, assetType, equalProtectionLevel)
		//assetGroupSrv.AssetAdd(ctx, v[0],v[1])
	}
	return nil
}

// AssetExportList 资产列表
func (a *AssetGroupApp) AssetExportList(ctx context.Context, req *typespec.AssetExportListReq, res *typespec.AssetExportListRes) error {
	var (
		assetGroupSrv   services.AssetGroup
		assetPortSrv    services.AssetPort
		groupIdIntSlice []int
		assetIds        []int
		allAsset        []mysqls.Asset
		count           int64
	)
	// 获取资产组下所有子资产信息 和路径
	if req.GroupId != 0 {
		subAssetGroupList, _, err := assetGroupSrv.GetSubAssetGroupList(ctx, req.GroupId)
		if err != nil {
			return err
		}
		// 获取所有资产组 对应的资产信息 以及统计
		groupIdIntSlice = []int{req.GroupId}
		for _, v := range subAssetGroupList {
			groupIdIntSlice = append(groupIdIntSlice, v.ID)
		}
		allAsset, count = assetGroupSrv.AllAssetListByGroupIdsAndSearch(ctx, groupIdIntSlice, req.AssetIP, req.Port, req.Service, req.SystemOp, req.VulName, req.Domain, req.Finger, req.Tags, req.AssetRisk, req.AssetType, req.FillingLevel)
	}
	if req.IDs != "" && req.IDs != "all" {
		ids := strings.Split(req.IDs, ",")
		for _, v := range ids {
			id, _ := strconv.Atoi(v)
			assetIds = append(assetIds, id)
		}
		allAsset, count = assetGroupSrv.AllAssetByIdsAndSearch(ctx, assetIds, req.AssetIP, req.Port, req.Service, req.SystemOp, req.VulName, req.Domain, req.Finger, req.Tags, req.AssetRisk, req.AssetType, req.FillingLevel)
	}
	if req.IDs == "all" {
		allAsset, count = assetGroupSrv.AllAssetByIdsAndSearch(ctx, []int{}, req.AssetIP, req.Port, req.Service, req.SystemOp, req.VulName, req.Domain, req.Finger, req.Tags, req.AssetRisk, req.AssetType, req.FillingLevel)
	}
	res.Count = int(count)
	// 如果检索服务/指纹值返回对应的端口和服务
	for _, v := range allAsset {
		// 获取资产端口信息
		jsonInfo := assetGroupSrv.GetAssetJsonInfo(v.Info)
		var openPort string
		assetGroupInfo, _ := assetGroupSrv.GetAssetGroupInfo(ctx, v.AssetGroupID)
		openPortList := assetPortSrv.GetAssetPortList(ctx, []string{v.IP})
		var (
			service  string
			assembly string
			protocol string
		)
		for _, v := range openPortList {
			service += v.Service + ","
			assembly += v.Assembly + ","
			protocol += v.Protocol + ","
		}
		if req.ExportType == enums.SpecialAssetExportType {
			if req.Service != "" {
				openPort = assetPortSrv.GetAssetPortByService(ctx, req.Service, v.IP)
			}
			service = req.Service
			//if req.Finger != "" {
			//	openPort = assetPortSrv.GetAssetPortByFinger(ctx, req.Finger, v.IP)
			//}
		} else {
			openPort = assetPortSrv.GetAssetPortInfo(ctx, []string{v.IP})
		}
		riskLevel, _ := services.GetAssetRiskLevel(ctx, v.IP)
		res.AssetsExportInfo = append(res.AssetsExportInfo, typespec.AssetsExportInfo{
			IP:                    v.IP,
			AssetName:             v.Name,
			OpSys:                 v.OperateSystem,
			AssetGroupName:        assetGroupInfo.Name,
			OpenPort:              openPort,
			Service:               service,
			AssetRiskName:         enums.ToolsVulnerabilityEnum.GetRiskEnum(riskLevel),
			ResponsibleDepartment: v.ResponsibleDepartment,
			Location:              v.Location,
			TestTime:              v.UpdateTime.Format(enums.TimeLayout),
			SystemAdmin:           jsonInfo.SystemAdmin,
			Tags:                  v.Tags,
			//AssetType:           enums.AssetEnum.GetAssetTypeName(v.AssetType),
			//IsCloudHost:         true,
			//EqualProtectionName: enums.AssetEnum.GetAssetFilingLevel(v.FilingLevel),
			//EquipmentForm: jsonInfo.EquipmentForm,
			//DeviceWeight: enums.AssetEnum.GetDeviceWeightName(v.DeviceWeight),
			//TrustLevel:  enums.AssetEnum.GetTrustLevelName(v.TrustLevel),
		})
	}
	return nil
}

// AssetGroupDetail 资产组详情
func (a *AssetGroupApp) AssetGroupDetail(ctx context.Context, req *typespec.AssetGroupDetailReq, resp *typespec.AssetGroupDetail) error {
	var (
		assetGroupSrv services.AssetGroup
		pidName       string
	)
	assetGroupInfo, _ := assetGroupSrv.GetAssetGroupInfo(ctx, req.Id)
	if assetGroupInfo.Pid != 0 {
		res, _ := assetGroupSrv.GetAssetGroupInfo(ctx, assetGroupInfo.Pid)
		pidName = res.Name
	}
	*resp = typespec.AssetGroupDetail{
		ID:      assetGroupInfo.ID,
		Name:    assetGroupInfo.Name,
		Pid:     assetGroupInfo.Pid,
		PidName: pidName,
		Remark:  assetGroupInfo.Remark,
	}
	return nil
}

// mergePortEntries 补全两个相同 port 的记录字段
func mergePortEntries(a, b typespec.ScanTaskPortListList) typespec.ScanTaskPortListList {
	// 只在 a 中字段为空时才使用 b 的值
	if a.Protocol == "" && b.Protocol != "" {
		a.Protocol = b.Protocol
	}
	if a.Service == "" && b.Service != "" {
		a.Service = b.Service
	}
	if a.Assembly == "" && b.Assembly != "" {
		a.Assembly = b.Assembly
	}
	// 保留较新的时间（如果你愿意）
	if b.CreateTime > a.CreateTime {
		a.CreateTime = b.CreateTime
	}
	// Remark 可以按需补充合并逻辑
	return a
}

// WritePort 写入端口
func WritePort(ctx context.Context) error {
	var (
		assetGroupSrv services.AssetGroup
	)

	// 获取资产相关信息
	allAsset := assetGroupSrv.GetAllAsset(ctx)
	var assetPortData []mysqls.Assetport
	for _, v := range allAsset {
		if v.IP == "" {
			continue
		}
		var (
			targetSrv services.TaskTarget // 获取渗透目标数据
			ids       []int
		)
		ip := v.IP
		target, _ := targetSrv.TargetListByTargetUrl(ctx, ip)
		targetRes, _ := targetSrv.GetTargetByTaskId(ctx, target.TaskID)
		for i := 0; i < len(targetRes); i++ {
			if targetRes[i].TargetURL == "" || targetRes[i].TargetURL != ip {
				continue
			}
			ids = append(ids, targetRes[i].ID)
		}
		taskResultRes, _ := targetSrv.GetAssetTargetOpenPort(ctx, ids)
		portMap := make(map[int]typespec.ScanTaskPortListList)
		for i := 0; i < len(taskResultRes); i++ {
			// 1_1 1_2
			var parsed struct {
				Port        string `json:"port"`
				Protocol    string `json:"protocol"`
				Service     string `json:"service"`
				Component   string `json:"component"`
				Title       string `json:"title"`
				Fingerprint string `json:"fingerprint"`
			}
			if err := json.Unmarshal([]byte(taskResultRes[i].JSONResult), &parsed); err != nil {
				continue
			}
			port, _ := strconv.Atoi(parsed.Port)
			newEntry := typespec.ScanTaskPortListList{
				Id:         taskResultRes[i].ID,
				Port:       port,
				Protocol:   parsed.Protocol,
				Service:    parsed.Service,
				Assembly:   parsed.Component,
				Remark:     "",
				CreateTime: taskResultRes[i].CreateTime.Format(enums.TimeLayout),
			}
			// 如果已有该端口，进行字段补全
			if existing, ok := portMap[port]; ok {
				portMap[port] = mergePortEntries(existing, newEntry)
			} else {
				portMap[port] = newEntry
			}
		}
		now := time.Now()
		for _, v := range portMap {
			assetPortData = append(assetPortData, mysqls.Assetport{
				IP:         target.TargetURL,
				Port:       v.Port,
				Protocol:   v.Protocol,
				Service:    v.Service,
				Assembly:   v.Assembly,
				Remark:     "新开放端口",
				Islive:     enums.AssetIsLiveYes,
				CreateTime: now,
				UpdateTime: now,
			})
		}
	}
	var assetPortModel mysqls.Assetport
	err := assetPortModel.AddAssetportMany(ctx, assetPortData)
	if err != nil {
		return err
	}
	return nil
}

// AssetConnList 资产连接方式列表
func (a *AssetGroupApp) AssetConnList(ctx context.Context, req *typespec.AssetConnListReq, res *typespec.AssetConnListRes) error {
	var assetConnSrv services.AssetConnectionService
	assetConnList, count, err := assetConnSrv.GetAssetConnectionsList(ctx, req.IP, req.Port, req.Protocol, req.Page, req.Size)
	if err != nil {
		return err
	}
	res.Count = count
	res.AssetConnInfo = make([]typespec.AssetConnInfo, 0)
	for _, v := range assetConnList {
		res.AssetConnInfo = append(res.AssetConnInfo, typespec.AssetConnInfo{
			ID:          v.ID,
			IP:          v.IP,
			Port:        v.Port,
			Protocol:    v.Protocol,
			ProtocolStr: enums.ConnEnum.GetConnMethodEnum(v.Protocol),
			User:        v.Username,
			Pass:        v.Password,
		})
	}
	return nil
}
