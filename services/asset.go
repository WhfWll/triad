package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/models/mysqls"
	aesEncryption "smart/tools/encryption"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"
)

type Asset struct {
	aesEcb aesEncryption.AesEcb
}

// 根据ip数值规则查询所有资产数据,并返回map
func (a *Asset) GetAllAssetByIpNumRange(ctx context.Context, ipNumRange string) (map[string]mysqls.Asset, map[string]int, []string, error) {
	var (
		assetModel mysqls.Asset
		ipNum      = make([]mysqls.IpNum, 0)
		assetMap   = make(map[string]mysqls.Asset, 0)
		ipMap      = make(map[string]int, 0)
		ipArray    = make([]string, 0)
		eqIpNum    = make([]mysqls.IpNum, 0)
		neqIpNum   = make([]mysqls.IpNum, 0)
	)
	if len(ipNumRange) == 0 {
		return assetMap, ipMap, ipArray, errors.New("ip num range is empty")
	}
	err := json.Unmarshal([]byte(ipNumRange), &ipNum)
	if err != nil {
		return assetMap, ipMap, ipArray, err
	}

	for i := 0; i < len(ipNum); i++ {
		if ipNum[i].Relation == "eq" || ipNum[i].Relation == "in" {
			eqIpNum = append(eqIpNum, ipNum[i])
		} else if ipNum[i].Relation == "neq" || ipNum[i].Relation == "nin" {
			neqIpNum = append(neqIpNum, ipNum[i])
		}
	}
	assetRes := assetModel.GetAssetByIpNumRange(ctx, eqIpNum, neqIpNum)
	for i := 0; i < len(assetRes); i++ {
		assetMap[assetRes[i].IP] = assetRes[i]
		ipMap[assetRes[i].IP] = 1
		ipArray = append(ipArray, assetRes[i].IP)
	}
	return assetMap, ipMap, ipArray, nil
}

// 查询所有资产端口数据，并返回map
func (a *Asset) GetAllAssetPortByIps(ctx context.Context, ips []string, portRange string) map[string]map[int]mysqls.Assetport {
	var (
		assetPortModel mysqls.Assetport
		assetPortMap   = make(map[string]map[int]mysqls.Assetport, 0)
		ports          = make([]string, 0)
		portrange      = make([][2]string, 0)
	)
	portArray := strings.Split(portRange, ",")
	for i := 0; i < len(portArray); i++ {
		if strings.Contains(portArray[i], "-") {
			tmpPort := strings.Split(portArray[i], "-")
			if len(tmpPort) >= 2 {
				portrange = append(portrange, [2]string{tmpPort[0], tmpPort[1]})
			}
		} else {
			ports = append(ports, portArray[i])
		}
	}
	assetPortRes := assetPortModel.GetAssetportByIps(ctx, ips, ports, portrange)
	for i := 0; i < len(assetPortRes); i++ {
		val, ok := assetPortMap[assetPortRes[i].IP]
		if ok {
			val[assetPortRes[i].Port] = assetPortRes[i]
			assetPortMap[assetPortRes[i].IP] = val
		} else {
			assetPortMap[assetPortRes[i].IP] = map[int]mysqls.Assetport{assetPortRes[i].Port: assetPortRes[i]}
		}
	}
	return assetPortMap
}

// 新增资产和资产端口
func (a *Asset) AddAssetIpHandle(ctx context.Context, ip string, assetScanIp mysqls.Assetscanip, scanPort map[int]mysqls.Assetscanport) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	//新增资产
	var assetModel = mysqls.Asset{
		AssetGroupID:          enums.AssetGroupUngroupedId,
		IP:                    ip,
		IPSegment:             "",
		AssetType:             0,
		IpNum:                 assetScanIp.IpNum,
		Name:                  "",
		OperateSystem:         assetScanIp.Os,
		BusinessSystem:        "",
		ResponsibleDepartment: "",
		FilingLevel:           enums.AssetFilingLevelZero,
		IsCloudHost:           0,
		Tags:                  "",
		Info:                  "",
		IsIgnore:              enums.AssetIsIgnoreNo,
		AssetChangesType:      enums.AssetChangeTypeAdd,
		Islive:                enums.AssetIsLiveYes,
		CreateTime:            time.Now(),
		UpdateTime:            time.Now(),
	}
	_, err := assetModel.AddAsset(dCtx)
	if err != nil {
		return err
	}
	//新增资产端口
	var assetPortModel mysqls.Assetport
	var assetPortData = make([]mysqls.Assetport, 0)
	for _, vale := range scanPort {
		var tmp = mysqls.Assetport{
			IP:         ip,
			Port:       vale.Port,
			Protocol:   vale.Protocol,
			Service:    vale.Service,
			Assembly:   vale.Assembly,
			Remark:     "新开放端口",
			Islive:     enums.AssetIsLiveYes,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		assetPortData = append(assetPortData, tmp)
	}

	if len(assetPortData) > 0 {
		err = assetPortModel.AddAssetportMany(dCtx, assetPortData)
		if err != nil {
			return err
		}
	}
	//修改扫描资产
	var paramsIp = map[string]interface{}{
		"status":             enums.TaskIpStatusFinish,
		"asset_changes_type": enums.TaskIpChangeTypeAdd,
		"update_time":        time.Now(),
	}
	err = assetScanIp.UpdateAssetscanip(dCtx, assetScanIp.ID, paramsIp)
	if err != nil {
		return err
	}
	//修改扫描端口
	var assetScanPortModel mysqls.Assetscanport
	var paramsPort = map[string]interface{}{
		"status":      enums.TaskPortStatusFinish,
		"remark":      "新开放端口",
		"update_time": time.Now(),
	}
	err = assetScanPortModel.UpdateAssetscanportByScanIpId(dCtx, assetScanIp.ID, paramsPort)
	if err != nil {
		return err
	}
	//提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil

}

// 根据忽略状态统计资产变化情况
func (a *Asset) CountAssetChangeType(ctx context.Context, isIgnore int) [7]int64 {
	var (
		assetModel mysqls.Asset
		result     [7]int64
	)
	assetChangeTypeRes := assetModel.CountAssetChangeTypeByIsIgnore(ctx, isIgnore)
	for i := 0; i < len(assetChangeTypeRes); i++ {
		if assetChangeTypeRes[i].AssetChangesType < 8 {
			result[assetChangeTypeRes[i].AssetChangesType] = assetChangeTypeRes[i].Total
		}
	}
	return result
}

// 资产变化列表查询
func (a *Asset) GetChangeAssetList(ctx context.Context, isIgnore int, search string, assetChangesType, isLive int, updateTime string, page, size int) ([]mysqls.Asset, []string, int64) {
	var (
		assetModel mysqls.Asset
		assetIps   = make([]string, 0)
	)
	assetRes, total := assetModel.GetChangeAssetList(ctx, isIgnore, search, assetChangesType, isLive, updateTime, page, size)
	for i := 0; i < len(assetRes); i++ {
		assetIps = append(assetIps, assetRes[i].IP)
	}
	return assetRes, assetIps, total
}

// 资产变化查询所有数据
func (a *Asset) GetChangeAssetAll(ctx context.Context, isIgnore int, search string, assetChangesType, isLive int) ([]mysqls.Asset, []string) {
	var (
		assetModel mysqls.Asset
		assetIps   = make([]string, 0)
	)
	assetRes := assetModel.GetChangeAssetAll(ctx, isIgnore, search, assetChangesType, isLive)
	for i := 0; i < len(assetRes); i++ {
		assetIps = append(assetIps, assetRes[i].IP)
	}
	return assetRes, assetIps
}

// 组装excel数据
func (a *Asset) AssembleChangeAssetList(assetRes []mysqls.Asset, assetPortRes map[string][]mysqls.Assetport) [][]string {
	var result = make([][]string, 0)
	//资产IP,操作系统,资产类型,资产状态变化,存活状态,更新时间,开放端口,协议,服务,组件,端口变化说明
	for i := 0; i < len(assetRes); i++ {
		var tmp = make([][]string, 0)
		v, ok := assetPortRes[assetRes[i].IP]
		if ok {
			for j := 0; j < len(v); j++ {
				var portTmp = []string{
					assetRes[i].IP,
					assetRes[i].OperateSystem,
					enums.AssetEnum.GetAssetTypeName(assetRes[i].AssetType),
					enums.AssetEnum.GetAssetChangeTypeName(assetRes[i].AssetChangesType),
					enums.AssetEnum.GetAssetIsLiveTypeName(assetRes[i].Islive),
					assetRes[i].UpdateTime.Format(enums.TimeLayout),
					strconv.Itoa(v[j].Port),
					v[j].Protocol,
					v[j].Service,
					v[j].Assembly,
					v[j].Remark,
				}
				tmp = append(tmp, portTmp)
			}
		} else {
			tmp = [][]string{{
				assetRes[i].IP,
				assetRes[i].OperateSystem,
				enums.AssetEnum.GetAssetTypeName(assetRes[i].AssetType),
				enums.AssetEnum.GetAssetChangeTypeName(assetRes[i].AssetChangesType),
				enums.AssetEnum.GetAssetIsLiveTypeName(assetRes[i].Islive),
				assetRes[i].UpdateTime.Format(enums.TimeLayout), "", "", "", "", "",
			}}
		}
		result = append(result, tmp...)
	}
	return result
}

// 查找资产开放端口
func (a *Asset) GetAssetPort(ctx context.Context, assetIps []string) map[string][]string {
	var (
		assetPortModel mysqls.Assetport
		result         = make(map[string][]string, 0)
	)
	portRes := assetPortModel.GetAssetPortByIpIslive(ctx, assetIps, enums.TaskIpIsLiveTypeYes)
	for i := 0; i < len(portRes); i++ {
		result[portRes[i].IP] = append(result[portRes[i].IP], strconv.Itoa(portRes[i].Port))
	}
	return result
}

// 查找资产端口数据
func (a *Asset) GetAssetPortData(ctx context.Context, assetIps []string) map[string][]mysqls.Assetport {
	var (
		assetPortModel mysqls.Assetport
		result         = make(map[string][]mysqls.Assetport, 0)
	)
	portRes := assetPortModel.GetAssetPortByIpIslive(ctx, assetIps, enums.TaskIpIsLiveTypeYes)
	for i := 0; i < len(portRes); i++ {
		result[portRes[i].IP] = append(result[portRes[i].IP], portRes[i])
	}
	return result
}

// 根据id忽略资产更新
func (a *Asset) UpdateAssetIgnore(ctx context.Context, assetIds any) error {
	var assetModel mysqls.Asset
	var param = map[string]interface{}{
		"is_ignore":   enums.AssetIsIgnoreYes,
		"update_time": time.Now(),
	}
	return assetModel.UpdateAssetByIds(ctx, assetIds, param)
}

// 1-未变化端口,2-删除端口，3-新增端口，4关闭->开放端口, 5-服务变化端口，6-组件变化端口，7开发->关闭端口
func differencePort(scanPortres map[int]mysqls.Assetscanport, assetPortRes map[int]mysqls.Assetport) (map[int]int, int) {
	var (
		result           = make(map[int]int, 0)
		assetChangesType = enums.TaskIpChangeTypeNo //资产-无变化（默认）
	)
	//新增-只有扫描端口有的
	for k, _ := range scanPortres {
		if _, ok := assetPortRes[k]; !ok {
			result[k] = enums.TaskPortChangeTypeAdd       //端口-新增
			assetChangesType = enums.TaskIpChangeTypePort //资产-端口变化
		}
	}
	//删除-只有资产端口有的
	for k, _ := range assetPortRes {
		if _, ok := scanPortres[k]; !ok {
			result[k] = enums.TaskPortChangeTypeReduce    //端口-减少
			assetChangesType = enums.TaskIpChangeTypePort //资产-端口变化
		}
	}
	//交集
	for k1, v1 := range scanPortres {
		if v2, ok := assetPortRes[k1]; ok {
			if v1.Islive == enums.TaskIpIsLiveTypeYes && v2.Islive == enums.AssetIsLiveNo {
				result[k1] = enums.TaskPortChangeTypePortOpen //端口-关闭调整为打开
				assetChangesType = enums.TaskIpChangeTypePort //资产-端口变化
			} else if v1.Islive == enums.TaskIpIsLiveTypeNo && v2.Islive == enums.AssetIsLiveYes {
				result[k1] = enums.TaskPortChangeTypePortClose //端口-打开调整为关闭
				assetChangesType = enums.TaskIpChangeTypePort  //资产-端口变化
			} else if v1.Service != v2.Service {
				result[k1] = enums.TaskPortChangeTypeService //端口-服务变化
				if assetChangesType == enums.TaskIpChangeTypeNo || assetChangesType == enums.TaskIpChangeTypeAssembly {
					assetChangesType = enums.TaskIpChangeTypeService //资产-服务变化
				}
			} else if v1.Assembly != v2.Assembly {
				result[k1] = enums.TaskPortChangeTypeAssembly //端口-组件变化
				if assetChangesType == enums.TaskIpChangeTypeNo {
					assetChangesType = enums.TaskIpChangeTypeAssembly //资产-组件变化
				}
			} else {
				result[k1] = enums.TaskPortChangeTypeNo //端口-无变化（默认）
			}
		}
	}
	return result, assetChangesType
}

// AllByIp 依据资产IP获取资产
func (a *Asset) AllByIp(ctx context.Context, ip []string) []mysqls.Asset {
	var assetModel mysqls.Asset
	return assetModel.AllByIp(ctx, ip)
}

// 批量添加资产
func (a *Asset) MultipartInsertAsset(ctx context.Context, list *[]mysqls.Asset) error {
	var assetModel mysqls.Asset
	return assetModel.MultipartInsert(ctx, list)
}

// 批量添加资产 - 信息收集
func (a *Asset) MultipartInsertAssetTaskresule(ctx context.Context, list *[]mysqls.Assettaskresult) error {
	var assetTaskResultModel mysqls.Assettaskresult
	return assetTaskResultModel.MultipartInsert(ctx, list)
}

// 批量添加资产 - 漏洞信息
func (a *Asset) MultipartInsertAssetVul(ctx context.Context, list *[]mysqls.Assetvul) error {
	var assetVulModel mysqls.Assetvul
	return assetVulModel.MultipartInsert(ctx, list)
}

// 更新资产同步状态
func (a *Asset) UpdateStatusDoneByIds(ctx context.Context, Ids []int) error {
	var assetModel mysqls.Asset
	return assetModel.UpdateStatusDoneByIds(ctx, Ids)
}

// UpdateAssetUpdateTime 更新资产-更新时间
func (a *Asset) UpdateAssetUpdateTime(ctx context.Context, Ids []int) error {
	var assetModel mysqls.Asset
	return assetModel.UpdateAssetUpdateTimeInfo(ctx, Ids)
}

// UpdateAssetInfo 更新资产-信息
func (a *Asset) UpdateAssetInfo(ctx context.Context, taskID int) {
	var (
		taskTargetMysql mysqls.TaskTarget
		assetPortData   []mysqls.Assetport
	)
	now := time.Now()
	// 通过taskID拿到测试目标target_url
	taskTargetList, err := taskTargetMysql.GetTargetsByTaskID(ctx, taskID)
	if err != nil {
		log.Printf("[ERROR] UpdateAssetInfo - taskID: %d | GetTargetsByTaskID failed: %v", taskID, err)
		return
	}
	for _, target := range taskTargetList {
		if target.TargetURL == "" {
			log.Printf("[WARN] UpdateAssetInfo - Empty TargetURL for targetID: %d", target.ID)
			continue
		}
		if target.IsAlive == enums.TargetIsAliveN {
			log.Printf("[WARN] UpdateAssetInfo - targetUrl is No Alive: %s", target.TargetURL)
			continue
		}
		// 记录原始URL
		originalURL := target.TargetURL
		// 处理URL信息
		target.TargetURL = utils.ExtractHost(target.TargetURL)
		log.Printf("[INFO] TargetID: %d, Original URL: %s, Processed URL: %s", target.ID, originalURL, target.TargetURL)
		// 通过target_url 找到对应的资产 然后更新风险状态
		assetID, err := a.processAssetRisk(ctx, target, taskID, now)
		if err != nil {
			log.Printf("[ERROR] UpdateAssetInfo - processAssetRisk failed | targetID: %d | err: %v", target.ID, err)
			continue
		}
		// 存资产漏洞信息
		if err := a.addAssetVul(ctx, target, assetID); err != nil {
			log.Printf("[ERROR] UpdateAssetInfo - addAssetVul failed | targetID: %d | assetID: %d | err: %v", target.ID, assetID, err)
		}
		// 查询资产端口信息
		ports := a.collectAssetPorts(ctx, target, now)
		assetPortData = append(assetPortData, ports...)
	}
	// 存储资产端口等信息
	a.saveAssetPorts(ctx, taskID, assetPortData)
}

// processAssetRisk 记录资产风险信息
func (a *Asset) processAssetRisk(ctx context.Context, target mysqls.TaskTarget, taskID int, now time.Time) (int, error) {
	asset := mysqls.Asset{}
	assetRiskTrend := mysqls.AssetRiskTrend{}
	riskLevel, riskCountArrStr := GetAssetRiskLevelByTargetID(ctx, target.ID)
	log.Printf("[INFO] UpsertAssetInfo - targetURL: %s | riskLevel: %d | riskNum: %s", target.TargetURL, riskLevel, riskCountArrStr)
	assetID, err := asset.UpsertAssetByAssetIP(ctx, target.TargetURL, riskCountArrStr, target.OpSys, riskLevel)
	if err != nil {
		log.Printf("[ERROR] UpdateAssetByAssetIP - ip: %s | err: %v", target.TargetURL, err)
		return 0, err
	}
	err = assetRiskTrend.AddAssetRiskTrendInfo(ctx, []mysqls.AssetRiskTrend{
		{
			AssetID:    assetID,
			IP:         target.TargetURL,
			RiskLevel:  riskLevel,
			RiskNum:    riskCountArrStr,
			CreateTime: now,
			TargetID:   target.ID,
			TaskID:     taskID,
		},
	})
	if err != nil {
		log.Printf("[ERROR] AddAssetRiskTrendInfo - ip: %s | err: %v", target.TargetURL, err)
		return 0, err
	}
	return assetID, nil
}

func (a *Asset) collectAssetPorts(ctx context.Context, target mysqls.TaskTarget, now time.Time) []mysqls.Assetport {
	var assetPortData []mysqls.Assetport
	targetSrv := TaskTarget{}
	taskResultRes, err := targetSrv.GetAssetTargetOpenPort(ctx, []int{target.ID})
	if err != nil {
		log.Printf("[ERROR] GetAssetTargetOpenPort - targetID: %d | err: %v", target.ID, err)
		return nil
	}
	for _, vv := range taskResultRes {
		var parsed struct {
			Port      string `json:"port"`
			Protocol  string `json:"protocol"`
			Service   string `json:"service"`
			Component string `json:"component"`
		}
		if err := json.Unmarshal([]byte(vv.JSONResult), &parsed); err != nil {
			log.Printf("[WARN] JSON unmarshal failed for targetID: %d | err: %v", target.ID, err)
			continue
		}
		port, err := strconv.Atoi(parsed.Port)
		if err != nil || port <= 0 || port > 65535 {
			log.Printf("[WARN] Invalid port in result for targetID: %d | raw port: %s", target.ID, parsed.Port)
			continue
		}
		assetPortData = append(assetPortData, mysqls.Assetport{
			IP:         target.TargetURL,
			Port:       port,
			Protocol:   parsed.Protocol,
			TaskID:     target.TaskID,
			Service:    parsed.Service,
			Assembly:   parsed.Component,
			Remark:     "新开放端口",
			Islive:     enums.AssetIsLiveYes,
			CreateTime: now,
			UpdateTime: now,
		})
	}
	return assetPortData
}

func (a *Asset) saveAssetPorts(ctx context.Context, taskID int, assetPortData []mysqls.Assetport) {
	if len(assetPortData) == 0 {
		log.Printf("[INFO] UpdateAssetInfo - No valid asset ports found for taskID: %d", taskID)
		return
	}
	assetPortMysql := mysqls.Assetport{}
	if err := assetPortMysql.UpsertAssetPortMany(ctx, assetPortData, taskID); err != nil {
		log.Printf("[ERROR] AddAssetportMany failed - taskID: %d | err: %v", taskID, err)
		return
	}
	log.Printf("[INFO] UpdateAssetInfo - Inserted %d asset ports for taskID: %d", len(assetPortData), taskID)
}

// addAssetVul 资产漏洞信息
func (a *Asset) addAssetVul(ctx context.Context, target mysqls.TaskTarget, assetID int) error {
	var (
		asset      Asset
		taskVulSrv TaskVul
	)
	if assetID == 0 {
		return errors.New("asset id is empty")
	}
	taskVulMap := taskVulSrv.AllByTargetIds(ctx, []int{target.ID})
	insertAssetVulDataList := make([]mysqls.Assetvul, 0)
	log.Printf("[INFO] addAssetVul - targetURL: %s | vulNum: %d ", target.TargetURL, len(taskVulMap))
	for _, taskVul := range taskVulMap[target.ID] {
		// 解密 TargetUrl
		var (
			pocName   string
			name      string
			targetUrl string
			location  string
			vulResult string
			verMsg    string
		)
		if utils.IsHexString(taskVul.Pocname) {
			pocNameDecodeByte, _ := hex.DecodeString(taskVul.Pocname)
			pocName = string(a.aesEcb.AesDecryptECB(pocNameDecodeByte, aesKey))
		}
		if utils.IsHexString(taskVul.Name) {
			nameDecodeByte, _ := hex.DecodeString(taskVul.Name)
			name = string(a.aesEcb.AesDecryptECB(nameDecodeByte, aesKey))
		}
		if utils.IsHexString(taskVul.TargetUrl) {
			targetUrlDecodeByte, _ := hex.DecodeString(taskVul.TargetUrl)
			targetUrl = string(a.aesEcb.AesDecryptECB(targetUrlDecodeByte, aesKey))
		}
		if utils.IsHexString(taskVul.Location) {
			locationDecodeByte, _ := hex.DecodeString(taskVul.Location)
			location = string(a.aesEcb.AesDecryptECB(locationDecodeByte, aesKey))
		}
		if utils.IsHexString(taskVul.VulResult) {
			vulResultDecodeByte, _ := hex.DecodeString(taskVul.VulResult)
			vulResult = string(a.aesEcb.AesDecryptECB(vulResultDecodeByte, aesKey))
		}
		if utils.IsHexString(taskVul.VerMsg) {
			verMsgDecodeByte, _ := hex.DecodeString(taskVul.VerMsg)
			verMsg = string(a.aesEcb.AesDecryptECB(verMsgDecodeByte, aesKey))
		}
		insertAssetVulDataList = append(insertAssetVulDataList, mysqls.Assetvul{
			TaskVulID:      taskVul.ID,
			AssetID:        assetID,
			TargetURL:      targetUrl,
			TaskID:         taskVul.TaskID,
			TargetID:       taskVul.TargetID,
			Pocname:        pocName,
			Name:           name,
			Class:          taskVul.Class,
			Type:           taskVul.Type,
			Risk:           taskVul.Risk,
			Location:       location,
			Status:         taskVul.Status,
			TestStatus:     taskVul.TestStatus,
			ExploitImpact:  taskVul.ExploitImpact,
			VulID:          taskVul.VulID,
			Description:    taskVul.Description,
			PublishedTime:  taskVul.PublishedTime,
			AffectRange:    taskVul.AffectRange,
			TargetResultID: taskVul.TargetResultID,
			VulNumber:      taskVul.VulNumber,
			VulAddress:     taskVul.VulAddress,
			RefURL:         taskVul.RefUrl,
			Cvss:           taskVul.Cvss,
			VulResult:      vulResult,
			VulParam:       taskVul.VulParam,
			VerMsg:         verMsg,
			IsReplace:      enums.AssetIsReplaceN,
			CreateTime:     time.Now(),
			UpdateTime:     time.Now(),
		})
	}
	if len(insertAssetVulDataList) != 0 {
		if err := asset.MultipartInsertAssetVul(ctx, &insertAssetVulDataList); err != nil {
			log.Println("[INFO] addAssetVul - MultipartInsertAssetVul Err " + err.Error())
			return err
		}
	}
	return nil
}
