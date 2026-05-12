package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"smart/api/typespec"
	"smart/models/mysqls"
	aesEncryption "smart/tools/encryption"
	"smart/tools/enums"
	"smart/tools/network"
	"smart/tools/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/redis"
)

type AssetGroup struct {
	aesEcb aesEncryption.AesEcb
}

// aesKey 用于解密（可后续抽取配置）
var aesKey = []byte("9876787656785679")

// GetAllAssetGroup 查询所有资产组数据
func (ag *AssetGroup) GetAllAssetGroup(ctx context.Context) []mysqls.Assetgroup {
	var AssetGroupModel mysqls.Assetgroup
	return AssetGroupModel.GetAllAssetgroup(ctx)
}

// GetAssetByIpOperateSystemRiskLevelTags 根据ip/操作系统/风险等级/标签/资产名称/业务系统/责任部门/备案等级查询资产
func (ag *AssetGroup) GetAssetByIpOperateSystemRiskLevelTags(ctx context.Context, search, ip, operateSystem string, riskLevel int, tags, assetName, businessSystem, responsibleDepartment string, filingLevel int, assetIds []int) []mysqls.Asset {
	var AssetModel mysqls.Asset
	return AssetModel.GetAssetByIpOperateSystemRiskLevelTags(ctx, search, ip, operateSystem, riskLevel, tags, assetName, businessSystem, responsibleDepartment, filingLevel, assetIds)
}

// FilterAssetGroupByAsset 根据资产对资产组进行过滤
func (ag *AssetGroup) FilterAssetGroupByAsset(assetGroupData []mysqls.Assetgroup, assetData []mysqls.Asset) []mysqls.Assetgroup {
	var (
		assetGroupMap = make(map[int]mysqls.Assetgroup, 0)
		mergedMap     = make(map[int]mysqls.Assetgroup, 0)
		result        []mysqls.Assetgroup
	)
	for i := 0; i < len(assetGroupData); i++ {
		assetGroupMap[assetGroupData[i].ID] = assetGroupData[i]
	}
	for i := 0; i < len(assetData); i++ {
		tmpMergedMap := ag.GetALLAssetGroupSuperiorData(assetGroupMap, assetData[i].AssetGroupID)
		for k, v := range tmpMergedMap {
			mergedMap[k] = v
		}
	}
	for _, v := range mergedMap {
		result = append(result, v)
	}
	return result
}

// GetALLAssetGroupSuperiorData 查找某个资产的所有上级资产组
func (ag *AssetGroup) GetALLAssetGroupSuperiorData(assetGroupData map[int]mysqls.Assetgroup, id int) map[int]mysqls.Assetgroup {
	var result = make(map[int]mysqls.Assetgroup, 0)
	for id != 0 {
		node, ok := assetGroupData[id]
		if !ok {
			break
		}
		result[node.ID] = node
		id = node.Pid
	}
	return result
}

// GetAllAsset 查询所有资产数据,资产数量可能会比较多，采用分批次查询，每批次查询10000条
func (ag *AssetGroup) GetAllAsset(ctx context.Context) []mysqls.Asset {
	var (
		AssetGroupModel mysqls.Asset
		result          = make([]mysqls.Asset, 0)
		batchSize       = 10000
		offset          = 0
	)
	for {
		tmpAssetRes := AssetGroupModel.GetAllAsset(ctx, batchSize, offset)
		if len(tmpAssetRes) == 0 {
			break
		}
		result = append(result, tmpAssetRes...)
		offset += batchSize
	}
	return result
}

// HandAssetGroupByAsset 合并资产和资产组数据结构
func (ag *AssetGroup) HandAssetGroupByAsset(assetGroupData []mysqls.Assetgroup, assetData []mysqls.Asset) []typespec.AssetTreeOverallRespItems {
	var result []typespec.AssetTreeOverallRespItems
	for _, v := range assetGroupData {
		result = append(result, typespec.AssetTreeOverallRespItems{Id: v.ID, Pid: v.Pid, Name: v.Name, Type: enums.AssetTreeNodeTypeOne})
	}
	for i := 0; i < len(assetData); i++ {
		result = append(result, typespec.AssetTreeOverallRespItems{Id: assetData[i].ID, Pid: assetData[i].AssetGroupID, Name: assetData[i].IP, Type: enums.AssetTreeNodeTypeTwo})
	}
	return result
}

// SortAssetGroupSlice 排序，先按照节点类型排序（资产组>资产），类型相同按照id(大->小)
func (ag *AssetGroup) SortAssetGroupSlice(assetGroupTreeData []typespec.AssetTreeOverallRespItems) []typespec.AssetTreeOverallRespItems {
	less := func(i, j int) bool {
		if assetGroupTreeData[i].Type == assetGroupTreeData[j].Type {
			return assetGroupTreeData[i].Id > assetGroupTreeData[j].Id
		}
		return assetGroupTreeData[i].Type < assetGroupTreeData[j].Type
	}
	sort.Slice(assetGroupTreeData, less)
	return assetGroupTreeData
}

// SortAssetAssetGroup 组装资产及资产组tree结构
func (ag *AssetGroup) SortAssetAssetGroup(assetGroupData []typespec.AssetTreeOverallRespItems, pid int) []typespec.AssetTreeOverallRespItems {
	var resultArr = make([]typespec.AssetTreeOverallRespItems, 0)
	for _, v := range assetGroupData {
		if v.Pid == pid {
			var child = make([]typespec.AssetTreeOverallRespItems, 0)
			if v.Type == enums.AssetTreeNodeTypeOne {
				child = ag.SortAssetAssetGroup(assetGroupData, v.Id)
			}
			node := typespec.AssetTreeOverallRespItems{Id: v.Id, Pid: v.Pid, Name: v.Name, Type: v.Type, Items: child}
			resultArr = append(resultArr, node)
		}
	}
	return resultArr
}

// SortAssetGroup 组装资产组tree结构
func (ag *AssetGroup) SortAssetGroup(assetGroupData []mysqls.Assetgroup, pid int) []typespec.AssetGroupRespItems {
	var resultArr = make([]typespec.AssetGroupRespItems, 0)
	for _, v := range assetGroupData {
		if v.ID == enums.AssetGroupUngroupedId { //未分类分组忽略
			continue
		}
		if v.Pid == pid {
			child := ag.SortAssetGroup(assetGroupData, v.ID)
			node := typespec.AssetGroupRespItems{Id: v.ID, Pid: v.Pid, Name: v.Name, Level: v.Level, Items: child}
			resultArr = append(resultArr, node)
		}
	}
	return resultArr
}

// AssetGroupAdd 添加任务组
func (ag *AssetGroup) AssetGroupAdd(ctx context.Context, name, remark string, pid, level int) error {
	assetModel := mysqls.Assetgroup{
		Name:       name,
		Remark:     remark,
		Pid:        pid,
		Level:      level,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err := assetModel.AddAssetgroup(ctx)
	if err != nil {
		return err
	}
	return nil
}

// AssetGroupAddReturnID 创建资产返回ID
func (ag *AssetGroup) AssetGroupAddReturnID(ctx context.Context, name string) (int, error) {
	assetModel := mysqls.Assetgroup{
		Name:       name,
		Pid:        0,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	id, err := assetModel.AddAndReturnID(ctx)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetAssetGroupInfo 获取资产组信息
func (ag *AssetGroup) GetAssetGroupInfo(ctx context.Context, groupID int) (mysqls.Assetgroup, error) {
	assetModel := mysqls.Assetgroup{
		ID: groupID,
	}
	assetGroupInfo, err := assetModel.GetAssetgroup(ctx)
	if err != nil {
		return mysqls.Assetgroup{}, err
	}
	return assetGroupInfo, nil
}

// CheckAssetGroupName 查看资产组名称
func (ag *AssetGroup) CheckAssetGroupName(ctx context.Context, groupID int) (string, error) {
	assetModel := mysqls.Assetgroup{
		ID: groupID,
	}
	groupName, err := assetModel.GetAssetGroupBreadcrumb(ctx, groupID)
	if err != nil {
		return "", err
	}
	return groupName, nil
}

// GetAssetGroupInfoByName 通过资产组名 获取资产组信息
func (ag *AssetGroup) GetAssetGroupInfoByName(ctx context.Context, groupName string) (mysqls.Assetgroup, error) {
	assetModel := mysqls.Assetgroup{
		Name: groupName,
	}
	assetGroupInfo, err := assetModel.GetAssetGroupByName(ctx)
	if err != nil {
		return mysqls.Assetgroup{}, err
	}
	return assetGroupInfo, nil
}

// GetAssetGroupInfoByNamePid 通过资产组名和父ID 获取资产组信息
func (ag *AssetGroup) GetAssetGroupInfoByNamePid(ctx context.Context, groupName string, pid int) (mysqls.Assetgroup, error) {
	assetModel := mysqls.Assetgroup{
		Name: groupName,
		Pid:  pid,
	}
	assetGroupInfo, err := assetModel.GetAssetGroupByNamePid(ctx)
	if err != nil {
		return mysqls.Assetgroup{}, err
	}
	return assetGroupInfo, nil
}

// UpdateAssetGroup 更新资产组
func (ag *AssetGroup) UpdateAssetGroup(ctx context.Context, id, pid int, name, remark string, level int) error {
	assetModel := mysqls.Assetgroup{
		ID:         id,
		Name:       name,
		Remark:     remark,
		Pid:        pid,
		Level:      level,
		UpdateTime: time.Now(),
	}
	err := assetModel.UpdateAssetgroup(ctx)
	if err != nil {
		return err
	}
	return nil
}

// DeleteByGroupIds 删除分组，依据分组ID
func (ag *AssetGroup) DeleteByGroupIds(ctx context.Context, groupIds []int) error {
	var assetGroup mysqls.Assetgroup
	return assetGroup.DeleteById(ctx, groupIds)
}

// RecursionGetAllSubAssetGroupId 依据获取某个资产下的所有资产ID
func (ag *AssetGroup) RecursionGetAllSubAssetGroupId(ctx context.Context, assetID int) (assetGroupIds []int) {
	var (
		assetGroup mysqls.Assetgroup
		list       []mysqls.Assetgroup
		err        error
	)
	assetGroupIds = append(assetGroupIds, assetID)
	list, err = assetGroup.GetAssetGroupInfoByPid(ctx, assetID)
	if err != nil || len(list) == 0 {
		return []int{assetID}
	}
	for _, info := range list {
		assetGroupIds = append(assetGroupIds, ag.RecursionGetAllSubAssetGroupId(ctx, info.ID)...)
	}
	return
}

// RecursionGetAllAssetGroupId 依据获取某个资产上的所有资产ID
func (ag *AssetGroup) RecursionGetAllAssetGroupId(ctx context.Context, assetID int) (assetGroupIds []int) {
	var assetGroup mysqls.Assetgroup
	info, err := assetGroup.GetAssetGroupInfo(ctx, assetID)
	if err != nil || info.ID == 0 {
		return []int{assetID}
	}
	assetGroupIds = append(assetGroupIds, info.ID)
	if info.Pid != 0 {
		assetGroupIds = append(assetGroupIds, ag.RecursionGetAllAssetGroupId(ctx, info.Pid)...)
	}
	return
}

// DeleteAssetByAssetIds 删除资产 通过资产ID
func (a *AssetGroup) DeleteAssetByAssetIds(ctx context.Context, assetIds []int) error {
	assetModel := mysqls.Asset{}
	return assetModel.UpdateAssetByIds(ctx, assetIds, map[string]interface{}{
		"asset_changes_type": enums.AssetChangeTypeReduce,
		"islive":             enums.AssetIsLiveNo,
		"update_time":        time.Now(),
	})
}

// DirectlyDeleteAssetByAssetIds 直接 删除资产 通过资产ID
func (a *AssetGroup) DirectlyDeleteAssetByAssetIds(ctx context.Context, assetIds []int) error {
	assetModel := mysqls.Asset{}
	assetRiskTrend := mysqls.AssetRiskTrend{}
	assetPort := mysqls.Assetport{}
	if err := assetModel.DeleteById(ctx, assetIds); err == nil {
		assetRiskTrend.DeleteByAssetIDs(ctx, assetIds)
		assetIp, _ := assetModel.GetIPsByIds(ctx, assetIds)
		assetPort.DeleteAssetPortByIPs(ctx, assetIp)
	}
	return nil
}

// PhysicsDeleteAssetByAssetId 物理删除资产 通过资产ID
func (a *AssetGroup) PhysicsDeleteAssetByAssetId(ctx context.Context, id []int) error {
	assetModel := mysqls.Asset{}
	return assetModel.DeleteById(ctx, id)
}

// UpdateAssetGroupID  批量更新资产组ID
func (a *AssetGroup) UpdateAssetGroupID(ctx context.Context, ids []int, assetGroupID int) error {
	var assetModel mysqls.Asset
	return assetModel.UpdateAssetGroupIDByIds(ctx, ids, assetGroupID)
}

// AllAssetByGroupIds 依据分组获取所有IP
func (a *AssetGroup) AllAssetByGroupIds(ctx context.Context, groupIds []int) []mysqls.Asset {
	var assetModel mysqls.Asset
	return assetModel.AllAssetByGroupIds(ctx, groupIds)
}

// GetAssetGroupDefaultGroupID  获取资产组默认资产组ID
func (a *AssetGroup) GetAssetGroupDefaultGroupID(ctx context.Context) int {
	var assetGroupModel mysqls.Assetgroup
	info, _ := assetGroupModel.GetDefaultAssetGroupInfo(ctx)
	return info.ID
}

// GetSubAssetGroupList 获取子资产组信息和资产组路径
func (ag *AssetGroup) GetSubAssetGroupList(ctx context.Context, groupID int) ([]mysqls.Assetgroup, string, error) {
	assetModel := mysqls.Assetgroup{}
	allAssetGroupList := assetModel.GetAllAssetgroup(ctx)
	info := findSubGroups(allAssetGroupList, groupID)
	return info, "", nil
}

// findSubGroups 递归处理子资产组相关信息 返回父节点下的所有子信息
func findSubGroups(assetGroups []mysqls.Assetgroup, parentID int) []mysqls.Assetgroup {
	var result []mysqls.Assetgroup
	for _, group := range assetGroups {
		if group.Pid == parentID {
			result = append(result, group)
			result = append(result, findSubGroups(assetGroups, group.ID)...)
		}
	}
	return result
}

// AllAssetByGroupIdsAndSearch 检索所有资产信息 分页 search
func (ag *AssetGroup) AllAssetByGroupIdsAndSearch(ctx context.Context, groupIds []int, assetIP, port, service, systemOp, vulName, domain, finger, tags, assetType, isCloudHost string, assetRisk *int, fillingLevel, page, size int) ([]mysqls.Asset, int64) {
	var (
		ips        []string
		assetModel mysqls.Asset
		ipSet      = make(map[string]struct{})
	)
	if port != "" || service != "" || finger != "" {
		var assetPortModel mysqls.Assetport
		portInt, _ := strconv.Atoi(port)
		assetPortLists := assetPortModel.GetIPsByPort(ctx, portInt, service, finger)
		if len(assetPortLists) != 0 {
			for _, v := range assetPortLists {
				ipSet[v.IP] = struct{}{}
			}
		}
	}
	if vulName != "" {
		var assVulModel mysqls.Assetvul
		assetVulList := assVulModel.GetAssetTargetURLByName(ctx, vulName)
		for _, v := range assetVulList {
			if v.TargetURL != "" {
				ipSet[v.TargetURL] = struct{}{}
			}
		}
	}
	if len(ipSet) != 0 {
		ips = make([]string, 0, len(ipSet))
		for ip := range ipSet {
			ips = append(ips, ip)
		}
	}
	return assetModel.AllAssetByGroupIdsAndSearch(ctx, groupIds, assetIP, systemOp, tags, assetType, isCloudHost, domain, assetRisk, fillingLevel, page, size, ips)
}

// AllAssetByIdsAndSearch 通过资产ID 检索所有资产信息 分页 search
func (ag *AssetGroup) AllAssetByIdsAndSearch(ctx context.Context, ids []int, assetIP, port, service, systemOp, vulName, domain, finger, tags string, assetRisk, assetType, fillingLevel int) ([]mysqls.Asset, int64) {
	var assetModel mysqls.Asset
	return assetModel.AllAssetByIdsAndSearchNoPage(ctx, ids, assetIP, port, service, systemOp, vulName, domain, finger, tags, assetRisk, assetType, fillingLevel)
}

// AllAssetListByGroupIdsAndSearch 检索所有资产信息
func (ag *AssetGroup) AllAssetListByGroupIdsAndSearch(ctx context.Context, groupIds []int, assetIP, port, service, systemOp, vulName, domain, finger, tags string, assetRisk, assetType, fillingLevel int) ([]mysqls.Asset, int64) {
	var assetModel mysqls.Asset
	return assetModel.AllAssetByGroupIdsAndSearchNoPage(ctx, groupIds, assetIP, port, service, systemOp, vulName, domain, finger, tags, assetRisk, assetType, fillingLevel)
}

// AssetAdd 资产添加
func (ag *AssetGroup) AssetAdd(ctx context.Context, ip, name, opSys, ipSegment, tags, info, responsibleDepartment string, assetGroupID, assetType, equalProtectionLevel, isCloudHost, deviceWeight, trustLevel int) (int, error) {
	if equalProtectionLevel == 0 {
		equalProtectionLevel = 1
	}
	var assetModel = mysqls.Asset{
		IP:                    ip,
		Name:                  name,
		IpNum:                 network.Ipv4ToInt64(ip),
		AssetGroupID:          assetGroupID,
		OperateSystem:         opSys,
		IPSegment:             ipSegment,
		AssetType:             assetType,
		IsCloudHost:           isCloudHost,
		DeviceWeight:          deviceWeight,
		TrustLevel:            trustLevel,
		Location:              "",
		FilingLevel:           equalProtectionLevel,
		ResponsibleDepartment: responsibleDepartment,
		Info:                  info,
		Tags:                  tags,
		IsIgnore:              enums.AssetIsIgnoreNo,
		AssetChangesType:      enums.AssetChangeTypeAdd,
		Islive:                enums.AssetIsLiveYes,
		CreateTime:            time.Now(),
		UpdateTime:            time.Now(),
	}
	return assetModel.AddAsset(ctx)
}

// UpdateAssetById 通过id更新资产信息
func (ag *AssetGroup) UpdateAssetById(ctx context.Context, ip, opSys, ipSegment, info, tags, responsibleDepartment, assetName string, assetGroupID, assetType, equalProtectionLevel, assetID, isCloudHost int) error {
	assetModel := mysqls.Asset{
		ID:                    assetID,
		IP:                    ip,
		IpNum:                 network.Ipv4ToInt64(ip),
		Name:                  assetName,
		AssetGroupID:          assetGroupID,
		OperateSystem:         opSys,
		IPSegment:             ipSegment,
		AssetType:             assetType,
		Info:                  info,
		IsCloudHost:           isCloudHost,
		ResponsibleDepartment: responsibleDepartment,
		FilingLevel:           equalProtectionLevel,
		IsIgnore:              enums.AssetIsIgnoreNo,
		Islive:                enums.AssetIsLiveYes,
		Tags:                  tags,
		UpdateTime:            time.Now(),
	}
	return assetModel.UpdateAssetByAssetID(ctx)
}

// AssetInfo 资产信息Json结构体
type AssetInfo struct {
	SoftwareName       string `json:"softwareName"`       // 基础软件名
	SoftwareVersion    string `json:"softwareVersion"`    // 基础软版本
	HardwareName       string `json:"hardwareName"`       // 基础硬件名称
	Purpose            string `json:"purpose"`            // 用途
	EquipmentForm      string `json:"equipmentForm"`      // 设备形态
	DeploymentLocation string `json:"deploymentLocation"` // 部署位置
	SystemName         string `json:"systemName"`         // 系统名称
	SystemAdmin        string `json:"systemAdmin"`        // 系统管理员
	SystemOperate      string `json:"systemOperate"`      // 系统运维人员
}

// OpAssetInfoJson 处理资产列表Info_Json信息
func (ag *AssetGroup) OpAssetInfoJson(baseSoftwareName, baseSoftwareVersion, baseHardwareName, purpose, equipmentForm, deploymentLocation, systemName, systemAdmin, systemOp string) string {
	assetInfo := AssetInfo{
		SoftwareName:       baseSoftwareName,
		SoftwareVersion:    baseSoftwareVersion,
		HardwareName:       baseHardwareName,
		Purpose:            purpose,
		EquipmentForm:      equipmentForm,
		DeploymentLocation: deploymentLocation,
		SystemName:         systemName,
		SystemAdmin:        systemAdmin,
		SystemOperate:      systemOp,
	}
	assetInfoByte, _ := json.Marshal(assetInfo)
	return string(assetInfoByte)
}

// GetAssetByID 通过id获取资产信息
func (ag *AssetGroup) GetAssetByID(ctx context.Context, assetID int) (mysqls.Asset, error) {
	assetModel := mysqls.Asset{
		ID: assetID,
	}
	return assetModel.GetAsset(ctx)
}

// GetAssetJsonInfo 获取资产Info_Json信息
func (ag *AssetGroup) GetAssetJsonInfo(info string) (assetInfo AssetInfo) {
	json.Unmarshal([]byte(info), &assetInfo)
	return
}

// GetAssetStaticsInfo 获取资产统计信息
func (ag *AssetGroup) GetAssetStaticsInfo(ctx context.Context) (assetTotalCount, assetLastChangeCount, newAddIP, thisWeekNewAddIP, newReduceIP, thisWeekNewReduceIP int64) {
	assetModel := mysqls.Asset{}
	// 资产总数
	_, assetTotalCount = assetModel.GetAllAssetList(ctx, "", "")
	// 同比上周
	now := time.Now()
	// 计算上周同一天的日期
	last7WeekSameDay := now.AddDate(0, 0, -7)
	last14WeekSameDay := now.AddDate(0, 0, -14)
	_, assetLastCount := assetModel.GetAllAssetList(ctx, last14WeekSameDay.Format(enums.TimeLayout), last7WeekSameDay.Format(enums.TimeLayout))
	_, assetNowCount := assetModel.GetAllAssetList(ctx, last7WeekSameDay.Format(enums.TimeLayout), now.Format(enums.TimeLayout))
	if assetNowCount-assetLastCount > 0 {
		assetLastChangeCount = assetNowCount - assetLastCount
	}
	var weekStart time.Time
	if now.Weekday() == time.Sunday {
		weekStart = now.AddDate(0, 0, -6)
	} else {
		weekStart = now.AddDate(0, 0, -int(now.Weekday())+1)
	}
	// 新增IP
	newAddIP = assetModel.CalNewIPChange(ctx, "", "", "add")
	// 本周新增IP
	thisWeekNewAddIP = assetModel.CalNewIPChange(ctx, weekStart.Format(enums.TimeToDayLayout), now.Format(enums.TimeLayout), "add")
	// 减少IP
	newReduceIP = assetModel.CalNewIPChange(ctx, "", "", "reduce")
	// 本周减少IP
	thisWeekNewReduceIP = assetModel.CalNewIPChange(ctx, weekStart.Format(enums.TimeToDayLayout), now.Format(enums.TimeLayout), "reduce")
	return
}

// AssetTypeStatics 资产类型统计
func (ag *AssetGroup) AssetTypeStatics(ctx context.Context) ([]typespec.AssetRiskStatistics, []string) {
	var assetModel mysqls.Asset
	assetList, _ := assetModel.GetAllAssetList(ctx, "", "")
	typeCount := make(map[int]int)
	ips := make([]string, 0, len(assetList))
	for _, asset := range assetList {
		ips = append(ips, asset.IP)
		typeCount[asset.AssetType]++
	}
	// 将统计结果放入结构体切片中
	var statistics []typespec.AssetRiskStatistics
	for assetType, count := range typeCount {
		statistics = append(statistics, typespec.AssetRiskStatistics{
			AssetType: enums.AssetEnum.GetAssetTypeName(assetType),
			Count:     count,
		})
	}
	return statistics, ips
}

// RiskTypeStatics 风险类型统计
func (ag *AssetGroup) RiskTypeStatics(ctx context.Context) ([]typespec.AssetRiskStatistics, []string) {
	var (
		assetModel mysqls.Asset
		taskTarget TaskTarget
		taskVul    mysqls.TaskVul
	)
	assetList, _ := assetModel.GetAllAssetList(ctx, "", "")
	vulTypeCount := make(map[string]int)
	ips := make([]string, 0, len(assetList))
	for _, asset := range assetList {
		ips = append(ips, asset.IP)
	}
	for _, asset := range assetList {
		targetRes, _ := taskTarget.TargetListByTargetUrl(ctx, asset.IP)
		vulRes, _ := taskVul.GetTaskVulListByTargetIds(ctx, []int{targetRes.ID}, enums.VulDataTypOne)
		for _, v := range vulRes {
			vulTypeCount[enums.ToolsVulnerabilityEnum.GetTypeEnum(v.Type)]++
		}
	}
	var statistics []typespec.AssetRiskStatistics
	for vulType, count := range vulTypeCount {
		statistics = append(statistics, typespec.AssetRiskStatistics{
			AssetType: vulType,
			Count:     count,
		})
	}
	return statistics, ips
}

// RiskTypeStaticsTS 风险类型统计
func (ag *AssetGroup) RiskTypeStaticsTS(ctx context.Context) ([]typespec.AssetRiskStatistics, []string) {
	var (
		assetModel mysqls.Asset
		taskVul    mysqls.TaskVul
		taskTarget TaskTarget
	)
	assetList, _ := assetModel.GetAllAssetList(ctx, "", "")
	vulTypeCount := make(map[string]int)
	ipList := make([]string, 0, len(assetList))
	ipToAssetMap := make(map[string]mysqls.Asset)
	for _, asset := range assetList {
		ipList = append(ipList, asset.IP)
		ipToAssetMap[asset.IP] = asset
	}
	// 批量获取 target
	targetList, _ := taskTarget.BatchTargetListByTargetUrl(ctx, ipList)
	targetIDs := make([]int, 0, len(targetList))
	for _, t := range targetList {
		targetIDs = append(targetIDs, t.ID)
	}
	// 批量获取漏洞信息
	vulList, err := taskVul.GetTaskVulListByTargetIds(ctx, targetIDs, enums.VulDataTypOne)
	if err != nil {
		return nil, nil
	}
	vulTypeCount = make(map[string]int)
	for _, vul := range vulList {
		vulType := enums.ToolsVulnerabilityEnum.GetTypeEnum(vul.Type)
		vulTypeCount[vulType]++
	}
	statistics := make([]typespec.AssetRiskStatistics, 0, len(vulTypeCount))
	for vulType, count := range vulTypeCount {
		statistics = append(statistics, typespec.AssetRiskStatistics{
			AssetType: vulType,
			Count:     count,
		})
	}
	finalIPs := make([]string, 0, len(targetList))
	for _, t := range targetList {
		finalIPs = append(finalIPs, t.TargetURL)
	}
	return statistics, finalIPs
}

// AssetRiskStatics 资产风险统计
func (ag *AssetGroup) AssetRiskStatics(ctx context.Context) ([]typespec.AssetRiskStatistics, int, int) {
	var (
		assetModel     mysqls.Asset
		taskTarget     TaskTarget
		taskVulSrv     TaskVul
		totalVuls      int
		statistics     []typespec.AssetRiskStatistics
		riskAssetCount int
	)
	assetList, _ := assetModel.GetAllAssetList(ctx, "", "")
	// 获取所有已完成的 target_url -> max(ID)
	targetMap, _ := taskTarget.GetFinishedTargetURLToIDMap(ctx)
	var targetIDs []int
	ipToTargetID := make(map[string]int)
	for _, asset := range assetList {
		for targetURL, id := range targetMap {
			if strings.Contains(targetURL, asset.IP) {
				if existID, ok := ipToTargetID[asset.IP]; !ok || id > existID {
					ipToTargetID[asset.IP] = id
				}
			}
		}
	}
	for _, id := range ipToTargetID {
		targetIDs = append(targetIDs, id)
	}
	// 一次性获取统计数据
	_, _, vulNumArrayMap, _ := taskVulSrv.GetTargetStatsBytargetIdsForAsset(ctx, targetIDs)
	// 使用 map 来统计每种资产类型的数量
	targetURLs := make([]string, 0, len(assetList))
	riskCount := make(map[int]int, 5)
	for _, asset := range assetList {
		targetID := ipToTargetID[asset.IP]
		// 统计风险等级
		riskLevel := asset.RiskLevel
		if riskLevel == enums.FatalRiskAsset {
			riskLevel = enums.HighRiskAsset
		}
		riskCount[riskLevel]++
		// 统计漏洞数量
		if vulNum, ok := vulNumArrayMap[targetID]; ok {
			totalVuls += vulNum[1] + vulNum[2] + vulNum[3] + vulNum[4]
		}
		targetURLs = append(targetURLs, asset.IP)
		if asset.RiskLevel != enums.SafeAsset {
			riskAssetCount++
		}
	}
	allRiskLevels := []int{
		enums.HighRiskAsset,
		enums.MiddleRiskAsset,
		enums.LowRiskAsset,
		enums.SafeAsset,
	}
	for _, level := range allRiskLevels {
		if _, exists := riskCount[level]; !exists {
			riskCount[level] = 0
		}
	}
	for _, level := range allRiskLevels {
		statistics = append(statistics, typespec.AssetRiskStatistics{
			AssetType: enums.AssetEnum.GetAssetRiskTypeName(level),
			Count:     riskCount[level],
		})
	}
	return statistics, totalVuls, riskAssetCount
}

// AssetRiskStaticsInHomePage 首页定制的资产风险统计 需要所有的风险等级
func (ag *AssetGroup) AssetRiskStaticsInHomePage(ctx context.Context) map[int]int {
	var (
		assetModel mysqls.Asset
		taskTarget TaskTarget
		taskVulSrv TaskVul
		totalVuls  int
	)
	assetList, _ := assetModel.GetAllAssetList(ctx, "", "")
	// 获取所有已完成的 target_url -> max(ID)
	targetMap, _ := taskTarget.GetFinishedTargetURLToIDMap(ctx)
	var targetIDs []int
	ipToTargetID := make(map[string]int)
	for _, asset := range assetList {
		for targetURL, id := range targetMap {
			if strings.Contains(targetURL, asset.IP) {
				if existID, ok := ipToTargetID[asset.IP]; !ok || id > existID {
					ipToTargetID[asset.IP] = id
				}
			}
		}
	}
	for _, id := range ipToTargetID {
		targetIDs = append(targetIDs, id)
	}
	// 一次性获取统计数据
	_, _, vulNumArrayMap, _ := taskVulSrv.GetTargetStatsBytargetIdsForAsset(ctx, targetIDs)
	// 使用 map 来统计每种资产类型的数量
	targetURLs := make([]string, 0, len(assetList))
	riskCount := make(map[int]int, 5)
	for _, asset := range assetList {
		targetID := ipToTargetID[asset.IP]
		// 统计风险等级
		riskLevel := asset.RiskLevel
		if riskLevel == enums.FatalRiskAsset {
			riskLevel = enums.HighRiskAsset
		}
		riskCount[riskLevel]++
		// 统计漏洞数量
		if vulNum, ok := vulNumArrayMap[targetID]; ok {
			totalVuls += vulNum[1] + vulNum[2] + vulNum[3] + vulNum[4]
		}
		targetURLs = append(targetURLs, asset.IP)
	}
	allRiskLevels := []int{
		enums.HighRiskAsset,
		enums.MiddleRiskAsset,
		enums.LowRiskAsset,
		enums.SafeAsset,
	}
	for _, level := range allRiskLevels {
		if _, exists := riskCount[level]; !exists {
			riskCount[level] = 0
		}
	}
	return riskCount
}

// AssetRiskTrend 最近12天资产风险分布趋势
func (ag *AssetGroup) AssetRiskTrend(ctx context.Context) []typespec.AssetRiskTrend {
	var (
		assetModel mysqls.Asset
		result     []typespec.AssetRiskTrend
	)
	for i := 11; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i)
		start := day.Format("2006-01-02") + " 00:00:00"
		end := day.Format("2006-01-02") + " 23:59:59"
		assetList, _ := assetModel.GetAllAssetList(ctx, start, end)
		var fatal, high, medium, low, safe int
		for _, asset := range assetList {
			level, _ := GetAssetRiskLevel(ctx, asset.IP)
			switch level {
			case enums.FatalRiskAsset: // 严重
				fatal++
			case enums.HighRiskAsset: // 高
				high++
			case enums.MiddleRiskAsset: // 中
				medium++
			case enums.LowRiskAsset: // 低
				low++
			case enums.SafeAsset: // 安全
				safe++
			}
		}
		result = append(result, typespec.AssetRiskTrend{
			Date:   day.Format("2006-01-02"),
			Fatal:  fatal,
			High:   high,
			Medium: medium,
			Low:    low,
			Safe:   safe,
		})
	}
	return result
}

// AssetTypeTrend 最近12天资产风险趋势
func (ag *AssetGroup) AssetTypeTrend(ctx context.Context) []typespec.AssetRiskTypeTrend {
	var (
		assetModel mysqls.Asset
		result     []typespec.AssetRiskTypeTrend
	)
	for i := 11; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i)
		start := day.Format("2006-01-02") + " 00:00:00"
		end := day.Format("2006-01-02") + " 23:59:59"
		assetList, _ := assetModel.GetAllAssetList(ctx, start, end)
		var safe int
		for _, asset := range assetList {
			level, _ := GetAssetRiskLevel(ctx, asset.IP)
			switch level {
			case enums.SafeAsset:
				safe++
			}
		}
		result = append(result, typespec.AssetRiskTypeTrend{
			Date:   day.Format("2006-01-02"),
			UnSafe: len(assetList) - safe,
			Safe:   safe,
		})
	}
	return result
}

// LastRiskAssets 最近新增的危险资产（高危或中危）
func (ag *AssetGroup) LastRiskAssets(ctx context.Context) []typespec.RiskAsset {
	var (
		assetModel mysqls.Asset
		result     []typespec.RiskAsset
		ap         AssetPort
	)
	assetList, _ := assetModel.GetAllAssetList(ctx, "", "")
	type assetWithRisk struct {
		asset     mysqls.Asset
		riskLevel int
	}
	var riskyAssets []assetWithRisk
	for _, asset := range assetList {
		riskLevel, _ := GetAssetRiskLevel(ctx, asset.IP)
		if riskLevel != enums.SafeAsset {
			riskyAssets = append(riskyAssets, assetWithRisk{asset: asset, riskLevel: riskLevel})
		}
	}
	sort.Slice(riskyAssets, func(i, j int) bool {
		return riskyAssets[i].asset.UpdateTime.After(riskyAssets[j].asset.UpdateTime)
	})
	limit := 10
	if len(riskyAssets) < 10 {
		limit = len(riskyAssets)
	}
	for _, item := range riskyAssets[:limit] {
		openPort := ap.GetAssetPortInfo(ctx, []string{item.asset.IP})
		riskLevel, riskCountArr := GetAssetRiskLevel(ctx, item.asset.IP)
		if riskCountArr[1] == 0 && riskCountArr[2] == 0 && riskCountArr[3] == 0 && riskCountArr[4] == 0 {
			continue
		}
		result = append(result, typespec.RiskAsset{
			ID:            item.asset.ID,
			IP:            item.asset.IP,
			OpenPort:      openPort,
			Os:            item.asset.OperateSystem,
			AssetTypeName: enums.AssetEnum.GetAssetTypeName(item.asset.AssetType),
			RiskLevel:     enums.GetAssetRisk(riskLevel),
			VulStatics: typespec.VulStatics{
				DeadlyVul:     riskCountArr[1],
				HighRiskVul:   riskCountArr[2],
				MediumRiskVul: riskCountArr[3],
				LowRiskVul:    riskCountArr[4],
			},
			Time: item.asset.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return result
}

// LastRiskAssetsTS 最近新增的危险资产（高危或中危） 调试
func (ag *AssetGroup) LastRiskAssetsTS(ctx context.Context) []typespec.RiskAsset {
	var (
		assetModel mysqls.Asset
		result     []typespec.RiskAsset
		ap         AssetPort
		taskTarget TaskTarget
		taskVulSrv TaskVul
	)
	assetList, _ := assetModel.GetAllAssetList(ctx, "", "")
	// 获取所有已完成目标映射（target_url -> max(id)）
	targetMap, _ := taskTarget.GetFinishedTargetURLToIDMap(ctx)
	ipToTargetID := make(map[string]int)
	for _, asset := range assetList {
		for targetURL, id := range targetMap {
			if strings.Contains(targetURL, asset.IP) {
				if existID, ok := ipToTargetID[asset.IP]; !ok || id > existID {
					ipToTargetID[asset.IP] = id
				}
			}
		}
	}

	// 获取目标 ID 列表
	var targetIDs []int
	for _, id := range ipToTargetID {
		targetIDs = append(targetIDs, id)
	}
	_, riskLevelMap, vulNumMap, _ := taskVulSrv.GetTargetStatsBytargetIdsForAsset(ctx, targetIDs)
	type assetWithRisk struct {
		asset     mysqls.Asset
		riskLevel int
		vulMap    [6]int
	}
	var riskyAssets []assetWithRisk
	for _, asset := range assetList {
		targetID := ipToTargetID[asset.IP]
		if targetID == 0 {
			continue
		}
		riskLevel := enums.SafeAsset
		if v, ok := riskLevelMap[targetID]; ok {
			riskLevel = v
		}
		if riskLevel == enums.FatalRiskAsset {
			riskLevel = enums.HighRiskAsset // 严重归入高危
		}
		if riskLevel == enums.SafeAsset {
			continue
		}
		vulMap := [6]int{}
		if v, ok := vulNumMap[targetID]; ok {
			vulMap = v
		}
		if vulMap[1] == 0 && vulMap[2] == 0 && vulMap[3] == 0 && vulMap[4] == 0 {
			continue
		}
		riskyAssets = append(riskyAssets, assetWithRisk{
			asset:     asset,
			riskLevel: riskLevel,
			vulMap:    vulMap,
		})
	}
	sort.Slice(riskyAssets, func(i, j int) bool {
		return riskyAssets[i].asset.UpdateTime.After(riskyAssets[j].asset.UpdateTime)
	})
	limit := 10
	if len(riskyAssets) < 10 {
		limit = len(riskyAssets)
	}
	for _, item := range riskyAssets[:limit] {
		openPort := ap.GetAssetPortForTarget(ctx, item.asset.IP)
		riskLevel, riskCountArr := GetAssetRiskLevel(ctx, item.asset.IP)
		if riskCountArr[1] == 0 && riskCountArr[2] == 0 && riskCountArr[3] == 0 && riskCountArr[4] == 0 {
			continue
		}
		result = append(result, typespec.RiskAsset{
			ID:            item.asset.ID,
			IP:            item.asset.IP,
			OpenPort:      openPort,
			Os:            item.asset.OperateSystem,
			AssetTypeName: enums.AssetEnum.GetAssetTypeName(item.asset.AssetType),
			RiskLevel:     enums.GetAssetRisk(riskLevel),
			VulStatics: typespec.VulStatics{
				DeadlyVul:     riskCountArr[1],
				HighRiskVul:   riskCountArr[2],
				MediumRiskVul: riskCountArr[3],
				LowRiskVul:    riskCountArr[4],
			},
			Time: item.asset.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return result
}

// AssetTrendChange 资产趋势变化
func (ag *AssetGroup) AssetTrendChange(ctx context.Context, timeType int, startTime, endTime, trendType string) typespec.AssetTrendChangeRes {
	var (
		trendChangeRes typespec.AssetTrendChangeRes
		assetModel     mysqls.Asset
		names          []string
		values         []float64
	)
	if timeType == 0 {
		timeType = 1
	}
	var start, endT time.Time
	if startTime == "" || endTime == "" {
		start = time.Now().AddDate(0, 0, -14)
		endT = time.Now()
	}
	// 根据不同的时间类型和时间范围选择查询条件
	switch timeType {
	case 1: // 最近14天
		for date := start; date.Before(endT); date = date.AddDate(0, 0, 1) {
			// 计算每一天新增的 IP 数量
			startTimes := date.Format(enums.TimeToDayLayout)
			endTimes := date.AddDate(0, 0, 1).Format(enums.TimeToDayLayout)
			count := assetModel.CalNewIPChange(ctx, startTimes, endTimes, trendType)
			names = append(names, date.Format("2006-01-02"))
			values = append(values, float64(count))
		}
	case 2: // 本月最近几天
		start, _ := time.Parse("2006-01-02", startTime)
		end, _ := time.Parse("2006-01-02", endTime)
		end = end.AddDate(0, 0, 1) // 将结束时间调整到后一天
		for date := start; date.Before(end); date = date.AddDate(0, 0, 1) {
			// 计算每一天新增的 IP 数量
			count := assetModel.CalNewIPChange(ctx, date.Format("2006-01-02"), date.AddDate(0, 0, 1).Format("2006-01-02"), "add")
			names = append(names, date.Format("2006-01-02"))
			values = append(values, float64(count))
		}
	case 4: // 本年最近几个月
		start, _ := time.Parse("2006-01-02", startTime)
		end, _ := time.Parse("2006-01-02", endTime)
		end = end.AddDate(0, 0, 1) // 将结束时间调整到后一天
		for date := start; date.Before(end); date = date.AddDate(0, 1, 0) {
			// 计算每一个月新增的 IP 数量
			count := assetModel.CalNewIPChange(ctx, date.Format("2006-01-02"), date.AddDate(0, 1, 0).Format("2006-01-02"), "add")
			names = append(names, date.Format("2006-01"))
			values = append(values, float64(count))
		}
	}
	// 将结果存入返回值
	trendChangeRes.Name = names
	trendChangeRes.Value = values
	return trendChangeRes
}

// GetNoLiveAssetByIp 通过IP检索不存活资产信息
func (ag *AssetGroup) GetNoLiveAssetByIp(ctx context.Context, ip string) mysqls.Asset {
	var AssetModel mysqls.Asset
	return AssetModel.GetNoLiveAssetByIp(ctx, ip)
}

// GetLiveAssetByIp 通过IP检索存活资产信息
func (ag *AssetGroup) GetLiveAssetByIp(ctx context.Context, ip string) mysqls.Asset {
	var AssetModel mysqls.Asset
	return AssetModel.GetLiveAssetByIp(ctx, ip)
}

// UpdateAssetByIdLive 更改IP状态存活
func (ag *AssetGroup) UpdateAssetByIdLive(ctx context.Context, assetID int) error {
	assetModel := mysqls.Asset{
		ID:               assetID,
		IsIgnore:         enums.AssetIsIgnoreNo,
		AssetChangesType: enums.AssetChangeTypeAdd,
		Islive:           enums.AssetIsLiveYes,
		UpdateTime:       time.Now(),
	}
	return assetModel.UpdateAssetByAssetID(ctx)
}

// AssetLevelOneStatics 一级资产统计
func (ag *AssetGroup) AssetLevelOneStatics(ctx context.Context) []typespec.AssetTypeStatistics {
	var (
		assetGroupModel mysqls.Assetgroup
		assetModel      mysqls.Asset
	)
	typeCount := make(map[string]int)
	allLevelOneAssetGroup, _ := assetGroupModel.GetAssetGroupInfoByPid(ctx, 0)
	for _, v := range allLevelOneAssetGroup {
		subAssetGroupList, _, _ := ag.GetSubAssetGroupList(ctx, v.ID)
		// 获取所有资产组 对应的资产信息 以及统计
		groupIdIntSlice := []int{v.ID}
		for _, v := range subAssetGroupList {
			groupIdIntSlice = append(groupIdIntSlice, v.ID)
		}
		_, count := assetModel.AllAssetByGroupIdsAndSearch(ctx, groupIdIntSlice, "", "", "", "", "", "", nil, 0, 0, 0, []string{})
		typeCount[v.Name] = int(count)
	}
	// 将统计结果放入结构体切片中
	var statistics []typespec.AssetTypeStatistics
	for name, count := range typeCount {
		statistics = append(statistics, typespec.AssetTypeStatistics{
			AssetType: name,
			Count:     count,
		})
	}
	// 按 count 值从大到小排序
	sort.Slice(statistics, func(i, j int) bool {
		return statistics[i].Count > statistics[j].Count
	})
	return statistics
}

// GetAssetRiskStatics 获取资产风险统计信息
func (ag *AssetGroup) GetAssetRiskStatics(ctx context.Context) (assetTotalCount, riskAssetTotalCount int64) {
	assetModel := mysqls.Asset{}
	// 资产总数
	assetTotalCount = assetModel.CountAllAssets(ctx)
	// 风险资产总数
	riskAssetTotalCount = assetModel.CountAllRiskAssets(ctx)
	return
}

// baseRiskScore 基础风险值
var baseRiskScore = map[int]int{
	enums.VulLibrariesRiskLow:    30,
	enums.VulLibrariesRiskMiddle: 60,
	enums.VulLibrariesRiskHigh:   80,
	enums.VulLibrariesRiskDead:   90,
	enums.VulLibrariesRiskNot:    0,
	enums.VulLibrariesRiskInfo:   0,
}

// riskLevelConfig 配置每种风险等级的加分规则
type riskLevelConfig struct {
	Multiplier float64
	MaxAdd     int
	MaxCount   int
}

var riskConfigMap = map[int]riskLevelConfig{
	enums.VulLibrariesRiskLow:    {Multiplier: 0.1, MaxAdd: 70, MaxCount: 70},
	enums.VulLibrariesRiskMiddle: {Multiplier: 0.3, MaxAdd: 40, MaxCount: 40},
	enums.VulLibrariesRiskHigh:   {Multiplier: 1.0, MaxAdd: 20, MaxCount: 20},
	enums.VulLibrariesRiskDead:   {Multiplier: 2.0, MaxAdd: 10, MaxCount: 10},
}

// CalculateAssetScore 计算资产得分（基础 + 额外）
func CalculateAssetScore(ctx context.Context, ipRiskMap [6]int, riskLevel int, ip string) string {
	var additionalNum int
	log.Printf("[评分计算] IP: %s 的漏洞风险等级分布: %+v", ip, ipRiskMap)
	if riskLevel == enums.VulLibrariesRiskNot || riskLevel == enums.VulLibrariesRiskInfo {
		log.Printf("[评分计算] IP: %s 未命中任何漏洞，得分为 0分", ip)
		return "0分"
	}
	// 获取风险配置
	cfg, ok := riskConfigMap[riskLevel]
	if ok {
		count := ipRiskMap[riskLevel]
		log.Printf("[评分计算] 风险等级: %d 命中数量: %d，阈值 MaxCount: %d，Multiplier: %.2f，MaxAdd: %d",
			riskLevel, count, cfg.MaxCount, cfg.Multiplier, cfg.MaxAdd)
		if count < cfg.MaxCount {
			rawAdd := float64(count) * cfg.Multiplier
			if rawAdd > float64(cfg.MaxAdd) {
				rawAdd = float64(cfg.MaxAdd)
				log.Printf("[评分计算] 加分超过最大限制，限制为 %.0f", rawAdd)
			}
			additionalNum = int(rawAdd)
			log.Printf("[评分计算] 最终附加分数: %d", additionalNum)
		} else {
			log.Printf("[评分计算] 命中数量 %d >= MaxCount %d，附加分为 0", count, cfg.MaxCount)
		}
	} else {
		log.Printf("[评分计算] 风险等级 %d 不存在于配置中，跳过附加分", riskLevel)
	}
	// 获取基础分
	base, ok := baseRiskScore[riskLevel]
	if !ok {
		base = 0
	}
	log.Printf("[评分计算] 风险等级 %d 的基础分: %d", riskLevel, base)
	score := base + additionalNum
	if score > 100 {
		log.Printf("[评分计算] 分数 %d 超过最大上限 100，重置为 100", score)
		score = 100
	}
	log.Printf("[评分计算] IP: %s 最终得分: %d分", ip, score)
	return strconv.Itoa(score) + "分"
}

// GetAssetRiskLevel 获取资产最近一条完成任务的风险等级和漏洞分布
func GetAssetRiskLevel(ctx context.Context, ip string) (int, [6]int) {
	var (
		taskVulSrv   TaskVul
		targetSrv    TaskTarget
		riskCountArr [6]int
		defaultLevel = enums.SafeAsset
	)
	// 查询目标资产信息
	target, err := targetSrv.TargetListByTargetUrl(ctx, ip)
	if err != nil || target.ID == 0 {
		return defaultLevel, riskCountArr
	}
	// 获取该目标关联的风险等级和漏洞数量数组
	_, riskLevelMap, vulNumMap, err := taskVulSrv.GetTargetStatsBytargetIdsForAsset(ctx, []int{target.ID})
	if err != nil {
		return defaultLevel, riskCountArr
	}
	// 获取风险等级
	riskLevel := enums.SafeAsset
	if v, ok := riskLevelMap[target.ID]; ok {
		riskLevel = v
	}
	if riskLevel == enums.FatalRiskAsset {
		riskLevel = enums.HighRiskAsset
	}
	// 获取风险等级对应漏洞数量分布
	if v, ok := vulNumMap[target.ID]; ok {
		riskCountArr = v
	}
	return riskLevel, riskCountArr
}

// GetAssetRiskLevelAll 获取资产全部风险等级
func GetAssetRiskLevelAll(ctx context.Context, ip string) (int, map[int]int) {
	var (
		taskVulSrv RiskManage
		aesEcb     aesEncryption.AesEcb
	)
	ipRiskMap := make(map[int]int, 4)
	// 通过vul_target表获取数据
	vulRes, _, _ := taskVulSrv.VulListByIP(ctx, "")
	for i := 0; i < len(vulRes); i++ {
		var tmp typespec.VulRiskListRespItems
		if utils.IsHexString(vulRes[i].TargetUrl) {
			targetUrlDecodeByte, _ := hex.DecodeString(vulRes[i].TargetUrl)
			tmp.TargetUrl = string(aesEcb.AesDecryptECB(targetUrlDecodeByte, []byte("9876787656785679")))
		} else {
			tmp.TargetUrl = vulRes[i].TargetUrl
		}
		if strings.Contains(tmp.TargetUrl, ip) {
			ipRiskMap[vulRes[i].Risk]++
		}
	}
	minRisk := 0
	for risk := range ipRiskMap {
		if minRisk == 0 || risk < minRisk {
			minRisk = risk
		}
	}
	if minRisk == 0 {
		minRisk = enums.SafeAsset
	}
	return minRisk, ipRiskMap
}

// AssetGroupEnums 查询所有资产组数据
func (ag *AssetGroup) AssetGroupEnums(ctx context.Context) interface{} {
	// 查询所有资产组
	var agModel mysqls.Assetgroup
	assetGroups := agModel.GetAllAssetgroup(ctx)

	// 构建枚举结果
	result := make([]struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}, 0, len(assetGroups))

	for _, ag := range assetGroups {
		result = append(result, struct {
			Value int    `json:"value"`
			Label string `json:"label"`
		}{
			Value: ag.ID,
			Label: ag.Name,
		})
	}

	return result
}

// GetCycleAssetRiskForLastNDays 获取最近 N 天的资产风险统计数据 【依照任务检测时间为纬度的统计】
func (ag *AssetGroup) GetCycleAssetRiskForLastNDays(ctx context.Context, days int) ([]typespec.AssetRiskTrend, int) {
	var (
		assetModel          mysqls.Asset
		assetRiskTrendModel mysqls.AssetRiskTrend
		trendList           []typespec.AssetRiskTrend
	)
	now := time.Now()
	for i := 0; i < days; i++ {
		day := now.AddDate(0, 0, -i)
		dateStr := day.Format("2006-01-02")
		var trend typespec.AssetRiskTrend
		startTime := dateStr + " 00:00:00"
		endTime := dateStr + " 23:59:59"
		assetRiskCondition := make(map[int]int, 5)
		// 1 先获取 在时间断前创建的资产数量
		lastAssetRiskList, lastAssetRiskCount := assetModel.GetAllAssetList(ctx, "", endTime)
		if i == 0 {
			for _, v := range lastAssetRiskList {
				assetRiskCondition[v.RiskLevel]++
			}
		} else {
			// 2 获取资产趋势信息 在时间断内 扫描的资产信息
			latestByIP := make(map[string]mysqls.AssetRiskTrend)
			assetRiskTrendList, _ := assetRiskTrendModel.GetAllRiskAssetTrendList(ctx, startTime, endTime)
			for _, v := range assetRiskTrendList {
				// 如果这个 IP 还没存，或者当前记录比已有记录更新，就替换
				if exist, ok := latestByIP[v.IP]; !ok || v.CreateTime.After(exist.CreateTime) {
					latestByIP[v.IP] = v
				}
			}
			for _, v := range latestByIP {
				assetRiskCondition[v.RiskLevel]++
				log.Println("Use record:", v.ID, v.IP, v.CreateTime, "RiskLevel:", v.RiskLevel, " sql startTime ", startTime, " endTime ", endTime)
			}
		}
		trend = typespec.AssetRiskTrend{
			Date:        dateStr,
			HighAsset:   assetRiskCondition[enums.HighRiskAsset] + assetRiskCondition[enums.FatalRiskAsset],
			MediumAsset: assetRiskCondition[enums.MiddleRiskAsset],
			LowAsset:    assetRiskCondition[enums.LowRiskAsset],
			SafeAsset:   int(lastAssetRiskCount) - assetRiskCondition[enums.FatalRiskAsset] - assetRiskCondition[enums.HighRiskAsset] - assetRiskCondition[enums.MiddleRiskAsset] - assetRiskCondition[enums.LowRiskAsset],
			TotalAsset:  int(lastAssetRiskCount),
		}
		trendList = append(trendList, trend)
	}
	slices.Reverse(trendList)
	return trendList, 10
}

// GetCycleTaskAssetRiskForLastNDays 获取最近 N 天的资产风险统计数据 【依照任务检测时间为纬度的统计】
func (ag *AssetGroup) GetCycleTaskAssetRiskForLastNDays(ctx context.Context, days int) ([]typespec.AssetRiskTrend, int) {
	var (
		assetGroup = AssetGroup{}
		trendList  []typespec.AssetRiskTrend
		totalVul   int
	)
	now := time.Now()
	redisClient, _ := redis.NewClient()
	for i := 0; i < days; i++ {
		day := now.AddDate(0, 0, -i)
		dateStr := day.Format("2006-01-02")
		cacheKey := fmt.Sprintf("asset_risk_trend:%s", dateStr)
		var trend typespec.AssetRiskTrend
		// 优先从 Redis 获取
		val, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			jsonErr := json.Unmarshal([]byte(val), &trend)
			if jsonErr == nil {
				trendList = append(trendList, trend)
				totalVul += trend.TotalVul
				continue
			}
			log.Printf("缓存数据解析失败，继续查询数据库: %v", jsonErr)
		}
		startTime := dateStr + " 00:00:00"
		endTime := dateStr + " 23:59:59"
		stat, total := assetGroup.CycleAssetRiskStatics(ctx, startTime, endTime)
		trend = typespec.AssetRiskTrend{
			Date:       dateStr,
			Fatal:      0,
			High:       0,
			Medium:     0,
			Low:        0,
			Safe:       0,
			TotalAsset: stat.TotalAsset,
			TotalVul:   total,
		}
		for _, item := range stat.AssetRiskTrendStatistics {
			switch item.AssetType {
			case enums.FatalRiskAsset:
				trend.Fatal += item.VulCount
				trend.FatalAsset += item.AssetCount
			case enums.HighRiskAsset:
				trend.High += item.VulCount
				trend.HighAsset += item.AssetCount
			case enums.MiddleRiskAsset:
				trend.Medium += item.VulCount
				trend.MediumAsset += item.AssetCount
			case enums.LowRiskAsset:
				trend.Low += item.VulCount
				trend.LowAsset += item.AssetCount
			case enums.SafeAsset:
				trend.Safe += item.VulCount
				trend.SafeAsset += item.AssetCount
			}
		}
		// NOTE：需求比较特殊 暂时不加严重 都归入高危资产中统计计算
		trend.HighAsset = trend.HighAsset + trend.FatalAsset

		// 缓存非今日数据
		if i > 0 {
			cacheBytes, _ := json.Marshal(trend)
			redisClient.Set(ctx, cacheKey, cacheBytes, 7*24*time.Hour)
		}
		trendList = append(trendList, trend)
	}
	slices.Reverse(trendList)
	return trendList, 10
}

// CycleAssetRiskStatics 周期资产风险统计
func (ag *AssetGroup) CycleAssetRiskStatics(ctx context.Context, startTime, endTime string) (typespec.CycleAssetRiskStatistics, int) {
	var (
		assetModel mysqls.Asset
		taskTarget TaskTarget
		taskVulSrv TaskVul
		totalVuls  int
		statistics []typespec.AssetRiskTrendStatistics
	)
	log.Printf("周期统计开始：时间范围 %s ~ %s", startTime, endTime)
	assetList, _ := assetModel.GetAllAssetList(ctx, "", endTime)
	log.Printf("共获取到 %d 个资产", len(assetList))
	// 获取 时间范围内 已完成的 target_url -> max(ID)
	targetMap, _ := taskTarget.GetFinishedTargetURLToIDMapByTime(ctx, startTime, endTime)
	log.Printf("任务映射表获取完成，共 %d 个目标", len(targetMap))
	var targetIDs []int
	ipToTargetID := make(map[string]int)
	for _, asset := range assetList {
		for targetURL, id := range targetMap {
			if strings.Contains(targetURL, asset.IP) {
				if existID, ok := ipToTargetID[asset.IP]; !ok || id > existID {
					ipToTargetID[asset.IP] = id
				}
			}
		}
	}
	for _, id := range ipToTargetID {
		targetIDs = append(targetIDs, id)
	}
	// 一次性获取统计数据
	_, riskLevelMap, vulNumArrayMap, _ := taskVulSrv.GetTargetStatsBytargetIdsForAsset(ctx, targetIDs)
	// 使用 map 来统计每种资产类型的数量
	targetURLs := make([]string, 0, len(assetList))
	riskCount := map[int]int{
		enums.FatalRiskAsset:  0,
		enums.HighRiskAsset:   0,
		enums.MiddleRiskAsset: 0,
		enums.LowRiskAsset:    0,
		enums.SafeAsset:       0,
	}
	riskAssetCount := make(map[int]int, 5)
	for _, asset := range assetList {
		targetID := ipToTargetID[asset.IP]
		// 统计风险等级
		riskLevel := enums.SafeAsset
		if v, ok := riskLevelMap[targetID]; ok {
			riskLevel = v
		}
		// 统计漏洞数量
		vulNum, ok := vulNumArrayMap[targetID]
		vulTotal := 0
		if ok {
			// 漏洞计数
			riskCount[enums.FatalRiskAsset] += vulNum[1]
			riskCount[enums.HighRiskAsset] += vulNum[2]
			riskCount[enums.MiddleRiskAsset] += vulNum[3]
			riskCount[enums.LowRiskAsset] += vulNum[4]
			vulTotal = vulNum[1] + vulNum[2] + vulNum[3] + vulNum[4]
			totalVuls += vulTotal
			riskAssetCount[riskLevel]++
		} else {
			// 不存在漏洞数据也计为资产数量
			riskAssetCount[riskLevel]++
		}
		// 输出详细日志
		log.Printf(
			"[资产:%s] 匹配任务ID:%d 风险等级:%s 漏洞数[Fatal:%d High:%d Medium:%d Low:%d] 总漏洞数:%d",
			asset.IP,
			targetID,
			enums.GetAssetRisk(riskLevel),
			vulNum[1], vulNum[2], vulNum[3], vulNum[4],
			vulTotal,
		)
		targetURLs = append(targetURLs, asset.IP)
	}
	for assetRisk, count := range riskCount {
		statistics = append(statistics, typespec.AssetRiskTrendStatistics{
			AssetType:  assetRisk,
			VulCount:   count,
			AssetCount: riskAssetCount[assetRisk],
		})
		log.Printf("风险类型:%s -> 漏洞总数:%d 资产数:%d ", enums.GetAssetRisk(assetRisk), count, riskAssetCount[assetRisk])
	}
	log.Printf("周期统计完成，总漏洞数: %d，总资产数: %d", totalVuls, len(assetList))
	return typespec.CycleAssetRiskStatistics{
		TotalAsset:               len(assetList),
		AssetRiskTrendStatistics: statistics,
	}, totalVuls
}

// GetAssetRiskLevelByTargetID 通过targetID获取资产最近一条完成任务的风险等级和漏洞分布
func GetAssetRiskLevelByTargetID(ctx context.Context, targetID int) (int, string) {
	var (
		taskVulSrv   TaskVul
		riskCountArr [6]int
		defaultLevel = enums.SafeAsset
	)
	// 获取该目标关联的风险等级和漏洞数量数组
	_, riskLevelMap, vulNumMap, err := taskVulSrv.GetTargetStatsBytargetIdsForAsset(ctx, []int{targetID})
	if err != nil {
		return defaultLevel, ""
	}
	// 获取风险等级
	riskLevel := enums.SafeAsset
	if v, ok := riskLevelMap[targetID]; ok {
		riskLevel = v
	}
	if riskLevel == enums.FatalRiskAsset {
		riskLevel = enums.HighRiskAsset
	}
	// 获取风险等级对应漏洞数量分布
	if v, ok := vulNumMap[targetID]; ok {
		riskCountArr = v
	}
	riskCountArrStr, _ := json.Marshal(riskCountArr)
	return riskLevel, string(riskCountArrStr)
}

// GetAllRiskAsset 查询风险资产情况
func (ag *AssetGroup) GetAllRiskAsset(ctx context.Context) (typespec.CycleAssetRiskStatistics, int64) {
	var (
		AssetGroupModel mysqls.Asset
		statistics      []typespec.AssetRiskTrendStatistics
	)
	riskCount := map[int]int{
		enums.FatalRiskAsset:  0,
		enums.HighRiskAsset:   0,
		enums.MiddleRiskAsset: 0,
		enums.LowRiskAsset:    0,
		enums.SafeAsset:       0,
	}
	riskAssetList, allAssetCount := AssetGroupModel.GetAllRiskAsset(ctx)
	for _, v := range riskAssetList {
		riskCount[v.RiskLevel]++
	}
	for assetRisk, count := range riskCount {
		statistics = append(statistics, typespec.AssetRiskTrendStatistics{
			AssetType:  assetRisk,
			AssetCount: count,
		})
	}
	return typespec.CycleAssetRiskStatistics{
		TotalAsset:               int(allAssetCount),
		AssetRiskTrendStatistics: statistics,
	}, 0
}
