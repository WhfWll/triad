package application

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"regexp"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/network"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"
)

// 三方工具
type TripartiteTools struct {
}

// TripartiteToolsXrayCreate xray 创建任务
func (t *TripartiteTools) TripartiteToolsXrayCreate(ctx context.Context, req *typespec.TripartiteToolsXRayCreateReq) error {
	// 验证目标填写是否正确
	if ok, err := network.IpSegmentTools.CheckSuccessUrl(req.Target); !ok {
		return err
	}

	// 创建任务
	var tripartiteSrv services.TripartiteTools
	_, err := tripartiteSrv.TripartiteToolsXrayCreate(ctx, req.TaskName, req.Target, req.IsCrawler)
	if err != nil {
		return err
	}

	return nil
}

// TripartiteToolsXrayUpload xray 文件上传
func (t *TripartiteTools) TripartiteToolsXrayUpload(ctx context.Context, req *typespec.TripartiteToolsXRayUploadReq, content []byte) error {
	// 获取文件内容
	var xrayStruct []utils.XRayItem
	err := json.Unmarshal(content, &xrayStruct)
	if err != nil {
		return errors.New("解析xray内容错误:" + err.Error())
	}

	// 事物开启 插入数据与修改数据
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	// 创建任务
	var tripartiteSrv services.TripartiteTools
	xrayId, err := tripartiteSrv.TripartiteToolsXrayCreate(dCtx, req.TaskName, "upload", false)
	if err != nil {
		return err
	}

	// 将结果储存入xray的结果表中
	xrayResultModelDatas := make([]services.TripartiteToolsXrayCreateResultItem, 0)
	for _, item := range xrayStruct {
		requestInfo, err := json.Marshal(item.Detail.Snapshot)
		if err != nil {
			log.Error("上传导入xray服务，解析请求信息失败：" + err.Error())
			continue
		}

		extraInfo, err := json.Marshal(item.Detail.Extra)
		if err != nil {
			log.Error("上传导入xray服务，解析扩展信息失败：" + err.Error())
			continue
		}
		createTime := time.Unix(item.CreateTime/1000, 0)

		xrayResultModelDatas = append(xrayResultModelDatas, services.TripartiteToolsXrayCreateResultItem{
			XrayTaskId:  xrayId,
			Addr:        item.Detail.Addr,
			Payload:     item.Detail.Payload,
			RequestInfo: string(requestInfo),
			Extra:       string(extraInfo),
			Plugin:      item.Plugin,
			CreateTime:  createTime,
			UpdateTime:  createTime,
		})
	}
	// 结果不为0则插入
	if len(xrayResultModelDatas) > 0 {
		err = tripartiteSrv.TripartiteToolsXrayResultCreate(dCtx, xrayResultModelDatas)
		if err != nil {
			return errors.New("上传导入xray调用结果失败：" + err.Error())
		}
	}

	err = tripartiteSrv.TripartiteToolsXrayFinish(dCtx, xrayId, len(xrayResultModelDatas))
	if err != nil {
		return errors.New("上传导入xray finish失败：" + err.Error())
	}

	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// TripartiteToolsXrayDel xray 删除任务
func (t *TripartiteTools) TripartiteToolsXrayDel(ctx context.Context, req *typespec.TripartiteToolsXRayDelReq) error {
	ids := strings.Split(req.XrayIds, ",")
	idInt := make([]int, 0)
	for _, id := range ids {
		idI, err := strconv.Atoi(strings.TrimSpace(id))
		if err == nil {
			idInt = append(idInt, idI)
		}
	}
	if len(ids) == 0 {
		return nil
	}

	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	// 删除任务
	var tripartiteSrv services.TripartiteTools
	err := tripartiteSrv.TripartiteToolsXrayDel(dCtx, idInt)
	if err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// TripartiteToolsXRayPage xray 任务列表
func (t *TripartiteTools) TripartiteToolsXRayPage(ctx context.Context, req *typespec.TripartiteToolsXRayPageReq, res *typespec.TripartiteToolsXRayPageRes) error {
	var tripartiteSrv services.TripartiteTools
	list, total, err := tripartiteSrv.TripartiteToolsXRayPage(ctx, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	res.Page = req.Page
	res.Size = req.Size
	res.Total = total
	for _, item := range list {
		res.List = append(res.List, typespec.TripartiteToolsXRayPageItem{
			Id:         item.ID,
			TaskName:   item.TaskName,
			RiskNum:    item.RiskNum,
			Status:     item.Status,
			StatusEnum: enums.XrayEnum.GetStatusEnum(item.Status),
			CreateTime: item.CreateTime.Format(utils.DateTime),
		})
	}
	return nil
}

// TripartiteToolsXRayDetailPage xray 任务详情列表
func (t *TripartiteTools) TripartiteToolsXRayDetailPage(ctx context.Context, req *typespec.TripartiteToolsXRayDetailPageReq, res *typespec.TripartiteToolsXRayDetailPageRes) error {
	var tripartiteSrv services.TripartiteTools
	list, total, err := tripartiteSrv.TripartiteToolsXRayDetailPage(ctx, req.XrayId, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	res.Page = req.Page
	res.Size = req.Size
	res.Total = total
	for _, item := range list {
		requestAndResponse := make([][]string, 0)
		_ = json.Unmarshal([]byte(item.RequestInfo), &requestAndResponse)

		var paramPosition, paramKey string
		extraMap := make(map[string]interface{})
		_ = json.Unmarshal([]byte(item.Extra), &extraMap)
		if param, ok := extraMap["param"]; ok {
			if paramByte, err := json.Marshal(param); err == nil {
				paramMap := make(map[string]interface{})
				_ = json.Unmarshal(paramByte, &paramMap)
				if _, ok := paramMap["position"]; ok {
					paramPosition = paramMap["position"].(string)
				}
				if _, ok := paramMap["key"]; ok {
					paramKey = paramMap["key"].(string)
				}
			}
		}
		res.List = append(res.List, typespec.TripartiteToolsXRayDetailPageItem{
			Id:                 item.ID,
			Addr:               item.Addr,
			PluginVul:          item.Plugin,
			ParamPosition:      paramPosition,
			ParamKey:           paramKey,
			Payload:            item.Payload,
			RequestAndResponse: requestAndResponse,
			Extra:              item.Extra,
			CreateTime:         item.CreateTime.Format(utils.DateTime),
		})
	}
	return nil
}

// TripartiteToolsBurpsuiteCreate burpsuite 创建任务
func (t *TripartiteTools) TripartiteToolsBurpsuiteCreate(ctx context.Context, req *typespec.TripartiteToolsBurpsuiteCreateReq) error {
	// 逗号分割目标 或 换行分割
	target := make([]string, 0)
	for _, item := range strings.Split(req.Target, ",") {
		if strings.Contains(item, "\n") {
			for _, item1 := range strings.Split(req.Target, "\n") {
				if final := strings.TrimSpace(item1); final != "" {
					target = append(target, final)
				}
			}
		} else {
			if final := strings.TrimSpace(item); final != "" {
				target = append(target, final)
			}
		}
	}
	if len(target) == 0 {
		return errors.New("target数据非法")
	}

	var burpsuiteSrv services.TripartiteTools
	_, err := burpsuiteSrv.TripartiteToolsBurpsuiteCreate(ctx, req.TaskName, strings.Join(target, ","))
	return err
}

// TripartiteToolsBurpsuiteUpload burpsuite 文件上传
func (t *TripartiteTools) TripartiteToolsBurpsuiteUpload(ctx context.Context, req *typespec.TripartiteToolsBurpsuiteUploadReq, content []byte) error {
	// 获取文件内容
	var burpsuiteStruct services.BurpsuiteResultData
	err := json.Unmarshal(content, &burpsuiteStruct)
	if err != nil {
		return errors.New("解析burpsuite内容错误:" + err.Error())
	}

	// 事物开启 插入数据与修改数据
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	// 创建任务
	var burpsuiteSrv services.TripartiteTools
	burpsuiteId, err := burpsuiteSrv.TripartiteToolsBurpsuiteCreate(dCtx, req.TaskName, "upload")
	if err != nil {
		return err
	}

	// 创建详情
	// 风险等级
	finalRisk := 5
	finalResultInsertData := make([]mysqls.BurpsuiteTaskResult, 0)
	for _, originData := range burpsuiteStruct.IssueEvents {
		// 风险等级
		risk := 5
		switch strings.ToLower(originData.Issue.Severity) {
		case "high":
			risk = 2
		case "medium":
			risk = 3
		case "low":
			risk = 4
		case "information":
			risk = 5
		}
		if finalRisk > risk {
			finalRisk = risk
		}

		// 请求与响应信息
		requestResponse, _ := json.Marshal(originData.Issue.Evidence)
		// 原数据结果ID
		originResultId, _ := strconv.Atoi(originData.Id)

		caption := originData.Issue.Caption
		insertionPoint := ""
		if strings.Contains(caption, "[") {
			re := regexp.MustCompile(`\[(.*?)\]`)
			match := re.FindStringSubmatch(caption)
			if len(match) > 1 {
				insertionPoint = match[1]
			}
		}

		finalResultInsertData = append(finalResultInsertData, mysqls.BurpsuiteTaskResult{
			BurpsuiteTaskID:       burpsuiteId,
			OriginResultId:        originResultId,
			Action:                originData.Type,
			IssueType:             originData.Issue.Name,
			Host:                  originData.Issue.Origin,
			Path:                  originData.Issue.Path,
			InsertionPoint:        insertionPoint,
			Severity:              strconv.Itoa(risk),
			Confidence:            originData.Issue.Confidence,
			Describe:              originData.Issue.Description,
			IssueBackground:       originData.Issue.IssueBackground,
			RemediationBackground: originData.Issue.RemediationBackground,
			RequestResponse:       string(requestResponse),
			InternalData:          originData.Issue.InternalData,
			CreateTime:            time.Now(),
			UpdateTime:            time.Now(),
		})

	}
	if len(finalResultInsertData) > 0 {
		err := burpsuiteSrv.TripartiteToolsBurpsuiteResultAdds(dCtx, finalResultInsertData)
		if err != nil {
			return errors.New("burpsuite 插入结果失败 err=" + err.Error())
		}
	}

	// 更新任务状态
	err = burpsuiteSrv.TripartiteToolsBurpsuiteUpdateStatusRiskById(dCtx, burpsuiteId, enums.BurpsuiteStatusDone, finalRisk)
	if err != nil {
		return errors.New("burpsuite 更新状态失败 err=" + err.Error())
	}

	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// TripartiteToolsBurpsuiteDel burpsuite 删除任务
func (t *TripartiteTools) TripartiteToolsBurpsuiteDel(ctx context.Context, req *typespec.TripartiteToolsBurpsuiteDelReq) error {
	ids := strings.Split(req.BurpsuiteIds, ",")
	idInt := make([]int, 0)
	for _, id := range ids {
		idI, err := strconv.Atoi(strings.TrimSpace(id))
		if err == nil {
			idInt = append(idInt, idI)
		}
	}
	if len(ids) == 0 {
		return nil
	}

	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	// 删除任务
	var tripartiteSrv services.TripartiteTools
	err := tripartiteSrv.TripartiteToolsBurpsuiteDel(dCtx, idInt)
	if err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// TripartiteToolsBurpsuitePage burpsuite 任务列表
func (t *TripartiteTools) TripartiteToolsBurpsuitePage(ctx context.Context, req *typespec.TripartiteToolsBurpsuitePageReq, res *typespec.TripartiteToolsBurpsuitePageRes) error {
	var burpsuiteSrv services.TripartiteTools
	list, total, err := burpsuiteSrv.TripartiteToolsBurpsuitePage(ctx, req.Page, req.Size, req.Search)

	if err != nil {
		return err
	}
	res.Page = req.Page
	res.Size = req.Size
	res.Total = total
	for _, item := range list {
		res.List = append(res.List, typespec.TripartiteToolsBurpsuitePageItem{
			Id:         item.ID,
			TaskName:   item.TaskName,
			RiskEnum:   enums.BurpsuiteEnum.GetRiskEnum(item.Risk),
			Status:     item.Status,
			StatusEnum: enums.XrayEnum.GetStatusEnum(item.Status),
			CreateTime: item.CreateTime.Format(utils.DateTime),
		})
	}
	return nil
}

// TripartiteToolsBurpsuiteDetailPage burpsuite 任务详情列表
func (t *TripartiteTools) TripartiteToolsBurpsuiteDetailPage(ctx context.Context, req *typespec.TripartiteToolsBurpsuiteDetailPageReq, res *typespec.TripartiteToolsBurpsuiteDetailPageRes) error {
	var burpsuiteSrv services.TripartiteTools
	list, total, err := burpsuiteSrv.TripartiteToolsBurpsuiteResultPage(ctx, req.BurpsuiteId, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	res.Page = req.Page
	res.Size = req.Size
	res.Total = total

	type requestResponseStruct struct {
		Type   string `json:"type"`
		Detail struct {
			BandFlags []string `json:"band_flags"`
		} `json:"detail"`
		RequestResponse struct {
			Url     string `json:"url"`
			Request []struct {
				Type   string `json:"type"`
				Data   string `json:"data"`
				Length int    `json:"length"`
			} `json:"request"`
			Response []struct {
				Type   string `json:"type"`
				Data   string `json:"data"`
				Length int    `json:"length"`
			} `json:"response"`
			WasRedirectFollowed bool   `json:"was_redirect_followed"`
			RequestTime         string `json:"request_time"`
		} `json:"request_response"`
	}

	for _, item := range list {
		riskEnum := ""
		if item.Severity != "" {
			risk, err := strconv.Atoi(item.Severity)
			if err == nil {
				riskEnum = enums.BurpsuiteEnum.GetRiskEnum(risk)
			}
		}

		requestResponseFinal := make([][]string, 0)
		var requestResponse []requestResponseStruct
		err = json.Unmarshal([]byte(item.RequestResponse), &requestResponse)
		if err == nil {
			for _, reqRes := range requestResponse {
				var reqString string
				var resString string
				for _, originReq := range reqRes.RequestResponse.Request {
					decodeData, err := base64.StdEncoding.DecodeString(originReq.Data)
					if err == nil {
						reqString += string(decodeData) + "\r\n"
					}
				}
				for _, originRes := range reqRes.RequestResponse.Response {
					decodeData, err := base64.StdEncoding.DecodeString(originRes.Data)
					if err == nil {
						resString += string(decodeData) + "\r\n"
					}
				}
				tempData := make([]string, 0)
				tempData = append(tempData, reqString)
				tempData = append(tempData, resString)
				requestResponseFinal = append(requestResponseFinal, tempData)
			}
		}
		res.List = append(res.List, typespec.TripartiteToolsBurpsuiteDetailPageItem{
			Id:              item.ID,
			Action:          item.Action,
			IssueType:       item.IssueType,
			Host:            item.Host,
			Path:            item.Path,
			InsertionPoint:  item.InsertionPoint,
			Confidence:      item.Confidence,
			Time:            item.CreateTime.Format(utils.DateTime),
			RiskEnum:        riskEnum,
			Desc:            item.Describe,
			IssueBackground: item.IssueBackground,
			Fix:             item.RemediationBackground,
			RequestResponse: requestResponseFinal,
		})
	}

	return nil
}

// TripartiteToolsWifiApList wifi 所有在线wifi列表
func (t *TripartiteTools) TripartiteToolsWifiApList(ctx context.Context, req *typespec.TripartiteToolsWifiApListReq, res *typespec.TripartiteToolsWifiApListRes) error {
	var tripartiteToolsSrv services.TripartiteTools
	apList, _ := tripartiteToolsSrv.TripartiteToolsWifiApList(ctx, req.Search, req.StartDate, req.EndDate)

	for _, item := range apList {
		res.List = append(res.List, typespec.TripartiteToolsWifiApListItem{
			Ssid:             item.Ssid,
			LastSignalRssi:   item.LastSignalRssi,
			SsidCryptset:     item.SsidCryptset,
			SsidCryptsetEnum: enums.WifiApInfoEnum.GetSsidCryptsetEnum(item.SsidCryptset),
			SourceMac:        item.SourceMac,
			Manuf:            item.Manuf,
			Carrier:          enums.WifiTaskEnum.GetCarrier(item.Carrierset),
			LastTime:         time.Unix(item.LastTime, 0).Format(utils.DateTime),
		})
	}
	return nil
}

// TripartiteToolsWifiApCreate wifi 创建任务
func (t *TripartiteTools) TripartiteToolsWifiApCreate(ctx context.Context, req *typespec.TripartiteToolsWifiCreateReq, res *typespec.TripartiteToolsWifiCreateRes) error {
	var tripartiteToolSrv services.TripartiteTools
	apInfo := tripartiteToolSrv.TripartiteToolsWifiApGetBySourceMac(ctx, req.SourceMac)
	if apInfo.SourceMac == "" {
		return errors.New("未知的Wifi信息")
	}

	// 创建任务时仅关注字段 mac status ssid task_name
	var wifiCreateData services.WifiCreateData
	wifiCreateData.TaskName = req.TaskName
	wifiCreateData.Mac = apInfo.SourceMac
	wifiCreateData.Status = enums.WifiTaskStatusIni
	wifiCreateData.PasswdSource = 0
	wifiCreateData.PasswdDict = ""
	wifiCreateData.Channel = apInfo.Channel
	wifiCreateData.Encrypt = apInfo.SsidCryptset
	wifiCreateData.Carrier = apInfo.Carrierset
	wifiCreateData.Passwd = ""
	wifiCreateData.StartTime = 0
	wifiCreateData.EndTime = 0
	wifiCreateData.Ssid = apInfo.Ssid
	wifiCreateData.IsSimulate = 0
	wifiCreateData.SimulateDuration = 30 * 60 // 30分钟
	wifiCreateData.IsCrack = 1                // 是否爆破
	wifiCreateData.IsEmbed = 0                // 是否模拟：0：否，1：是

	return tripartiteToolSrv.TripartiteToolsWifiApCreate(ctx, wifiCreateData)
}

// TripartiteToolsWifiPage wifi 列表
func (t *TripartiteTools) TripartiteToolsWifiPage(ctx context.Context, req *typespec.TripartiteToolsWifiPageReq, res *typespec.TripartiteToolsWifiPageRes) error {
	var tripartiteToolSrv services.TripartiteTools
	// 任务列表
	data, total := tripartiteToolSrv.TripartiteToolsWifiPage(ctx, req.Page, req.Size, req.Search)
	res.Total = total
	res.Page = req.Page
	res.Size = req.Size
	if len(data) == 0 {
		return nil
	}

	// 任务ID
	taskId := make([]int, 0)
	for _, item := range data {
		taskId = append(taskId, item.TaskID)
	}

	// 日志列表
	taskLogMaps, err := tripartiteToolSrv.TripartiteToolsWifiLogByTaskId(ctx, taskId)
	if err != nil {
		return err
	}

	for _, item := range data {
		logTemp := make([]map[string]string, 0)
		for _, log := range taskLogMaps[item.TaskID] {
			logTemp = append(logTemp, map[string]string{
				"log":  log.Content,
				"time": time.Unix(log.GenerateTime, 0).Format(utils.DateTime),
			})
		}
		res.List = append(res.List, typespec.TripartiteToolsWifiPageItem{
			TaskId:   item.TaskID,
			TaskName: item.TaskName,
			Ssid:     item.Ssid,
			Encrypt:  enums.WifiTaskEnum.GetCryptset(item.Encrypt),
			Carrier:  enums.WifiTaskEnum.GetCarrier(item.Carrier),
			//Channel:   enums.WifiApInfoEnum.GetHandlerEnum(item.Channel),
			Channel:    strconv.Itoa(item.Channel),
			StartTime:  time.Unix(item.StartTime, 0).Format(utils.DateTime),
			CreateTime: item.CreateTime.Format(utils.DateTime),
			Status:     enums.WifiTaskEnum.GetStatus(item.Status),
			Passwd:     item.Passwd,
			LogList:    logTemp,
		})
	}
	return nil
}

// TripartiteToolsWifiDel wifi 删除
func (t TripartiteTools) TripartiteToolsWifiDel(ctx context.Context, req *typespec.TripartiteToolsWifiDelReq, res *typespec.TripartiteToolsWifiDelRes) error {
	taskIds := strings.Split(req.TaskIds, ",")
	if len(taskIds) > 0 {
		taskIdInt := make([]int, 0)
		for _, id := range taskIds {
			idI, err := strconv.Atoi(id)
			if err == nil {
				taskIdInt = append(taskIdInt, idI)
			}
		}
		if len(taskIdInt) > 0 {
			var tripartiteToolSrv services.TripartiteTools
			return tripartiteToolSrv.TripartiteToolsWifiDel(ctx, taskIdInt)
		}
	}
	return nil
}
