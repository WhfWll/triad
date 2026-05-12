package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/invoke"
	"smart/tools/network"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// TaskTaskResult 通用漏洞信息管理内容
type TaskTaskResult struct {
}

// HandleResult 加工处理结果信息
func (t *TaskTaskResult) HandleResult(ctx context.Context, scannerLog invoke.ScannerLog, targetId int) ([]map[string]string, error) {
	resultMapList := make([]map[string]string, 0)
	switch scannerLog.Pocname {
	case enums.ScriptNamePortScan: //如果是端口扫描结果
		var portScanResult PortScanResultLog
		// 尝试解析为新的 PortScanResultLog 格式
		err := json.Unmarshal([]byte(scannerLog.Content), &portScanResult)

		// 如果解析成功且包含 FingerPrintResult，则处理新格式
		if err == nil && portScanResult.FingerPrintResult != nil {
			for _, serviceInfo := range portScanResult.FingerPrintResult.ServiceInfo {
				// 构造 MatchResult 以兼容现有逻辑
				portInt, _ := strconv.Atoi(serviceInfo.Port)
				matchResult := MatchResult{
					Target: serviceInfo.Host,
					Port:   portInt,
					Fingerprint: &FingerprintInfo{
						ServiceName: serviceInfo.ServiceName,
						// 可以补充更多字段如果需要
					},
				}

				// 处理服务信息 - 批量处理所有端口
				serviceResults, err := t.handlePortScanServiceList(ctx, matchResult)
				if err != nil {
					log.Error("HandleResult handlePortScanServiceList error: " + err.Error())
				}
				resultMapList = append(resultMapList, serviceResults...)
				// 处理站点信息 - 批量处理所有HTTP/HTTPS端口
				siteResults, err := t.handlePortScanSiteList(ctx, matchResult)
				if err != nil {
					log.Error("HandleResult handlePortScanSiteList error: " + err.Error())
				}
				resultMapList = append(resultMapList, siteResults...)
			}
		} else {
			// 尝试解析为旧的 MatchResult 格式 (兼容分布式节点可能返回的旧格式)
			var matchResult MatchResult
			err := json.Unmarshal([]byte(scannerLog.Content), &matchResult)
			if err != nil {
				// 只有当两种格式都解析失败时才报错
				fmt.Println("Parse port scan result error:", err.Error())
				return resultMapList, nil
			}

			// 如果解析成功，直接使用 matchResult 处理
			if matchResult.Port != 0 {
				// 处理服务信息
				serviceResults, err := t.handlePortScanServiceList(ctx, matchResult)
				if err != nil {
					log.Error("HandleResult handlePortScanServiceList error: " + err.Error())
				}
				resultMapList = append(resultMapList, serviceResults...)
				// 处理站点信息
				siteResults, err := t.handlePortScanSiteList(ctx, matchResult)
				if err != nil {
					log.Error("HandleResult handlePortScanSiteList error: " + err.Error())
				}
				resultMapList = append(resultMapList, siteResults...)
			}
		}

		//case enums.ScriptNameCrawlerx:
		//	itemMap, err := t.handleCrawlerUrl(ctx, detailMap)
		//	if err != nil {
		//		log.Error("HandleResult error: " + err.Error())
		//	}
		//	if itemMap["subObjType"] != "" {
		//		resultMapList = append(resultMapList, itemMap)
		//	}
		//case enums.ScriptNameSubdomain: //如果是子域名爆破检测结果
		//	itemMap, err := t.handleSubdomain(ctx, detailMap)
		//	if err != nil {
		//		log.Error("HandleResult error: " + err.Error())
		//	}
		//	if itemMap["subObjType"] != "" {
		//		resultMapList = append(resultMapList, itemMap)
		//	}
		//case enums.ScriptNameWafDetect: //
		//	itemMap, err := t.handleWafDetect(ctx, detailMap, targetId)
		//	if err != nil {
		//		log.Error("handleWafDetect error: " + err.Error())
		//	}
		//	if itemMap["subObjType"] != "" {
		//		resultMapList = append(resultMapList, itemMap)
		//	}
		//case enums.ScriptNameCdnDetect: //
		//	itemMap, err := t.handleCDN(ctx, detailMap)
		//	if err != nil {
		//		log.Error("handleCDN error: " + err.Error())
		//	}
		//	if itemMap["subObjType"] != "" {
		//		resultMapList = append(resultMapList, itemMap)
		//	}
		//case enums.ScriptNameWhois:
		//	itemMap, err := t.handleWhois(ctx, detailMap)
		//	if err != nil {
		//		log.Error("handleWhois error: " + err.Error())
		//	}
		//	if itemMap["subObjType"] != "" {
		//		resultMapList = append(resultMapList, itemMap)
		//	}
	}
	return resultMapList, nil
}

func (t *TaskTaskResult) HandleCrawlerResult(ctx context.Context, scannerLog invoke.ScannerLog) ([]map[string]string, error) {
	resultMapList := make([]map[string]string, 0)
	if len(scannerLog.Result) == 0 {
		return resultMapList, nil
	}

	crawlerMode := toCrawlerResultString(scannerLog.Result["crawler_mode"])
	targetURL := toCrawlerResultString(scannerLog.Result["target_url"])

	// 优先从 crawl_items 或 items 获取结构化对象（含 method/status_code/title/risk_type/content_type 等）
	var items []interface{}
	if v, ok := scannerLog.Result["crawl_items"].([]interface{}); ok && len(v) > 0 {
		items = v
	} else if v, ok := scannerLog.Result["items"].([]interface{}); ok && len(v) > 0 {
		items = v
	} else if v, ok := scannerLog.Result["crawl_urls"].([]interface{}); ok && len(v) > 0 {
		items = v
	} else {
		return resultMapList, nil
	}

	for _, item := range items {
		detailMap := map[string]string{
			"url":         "",
			"method":      "",
			"status":      "",
			"title":       "",
			"param":       "",
			"bodyLength":  "",
			"riskType":    "",
			"contentType": "",
			"crawlerMode": crawlerMode,
			"targetURL":   targetURL,
		}

		switch v := item.(type) {
		case string:
			detailMap["url"] = strings.TrimSpace(v)
		case map[string]interface{}:
			detailMap["url"] = toCrawlerResultString(v["url"])
			detailMap["method"] = strings.ToUpper(toCrawlerResultString(v["method"]))
			detailMap["status"] = toCrawlerResultString(v["status_code"])
			if detailMap["status"] == "" {
				detailMap["status"] = toCrawlerResultString(v["status"])
			}
			detailMap["title"] = toCrawlerResultString(v["title"])
			detailMap["param"] = toCrawlerResultString(v["params"])
			if detailMap["param"] == "" {
				detailMap["param"] = toCrawlerResultString(v["param"])
			}
			detailMap["bodyLength"] = toCrawlerResultString(v["body_length"])
			if detailMap["bodyLength"] == "" {
				detailMap["bodyLength"] = toCrawlerResultString(v["bodyLength"])
			}
			detailMap["riskType"] = toCrawlerResultString(v["risk_type"])
			detailMap["contentType"] = toCrawlerResultString(v["content_type"])
		default:
			detailMap["url"] = strings.TrimSpace(fmt.Sprint(v))
		}

		if detailMap["url"] == "" {
			continue
		}

		itemMap, err := t.handleCrawlerUrl(ctx, detailMap)
		if err != nil {
			return resultMapList, err
		}
		if itemMap["subObjType"] != "" {
			resultMapList = append(resultMapList, itemMap)
		}
	}

	return resultMapList, nil
}

func toCrawlerResultString(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

// handlePortScanServiceList 批量处理端口扫描服务信息
func (t *TaskTaskResult) handlePortScanServiceList(ctx context.Context, fingerResult MatchResult) ([]map[string]string, error) {
	resultMapList := make([]map[string]string, 0)
	itemMap := make(map[string]string)
	itemMap["subObjType"] = enums.TaskResultSubObjTypeService
	itemMap["field1"] = fingerResult.Target
	itemMap["field2"] = strconv.Itoa(fingerResult.Port)

	var services string
	jsonResultMap := make(map[string]string, 0)
	var components string
	var componentMap = make(map[string]bool)

	if fingerResult.Fingerprint != nil {
		replaceReg, _ := regexp.Compile(`\[.+?]`)
		fingerPrintStr := replaceReg.ReplaceAllString(fingerResult.Fingerprint.ServiceName, "")
		for _, item := range strings.Split(fingerPrintStr, "/") {
			isFingerOrSystem := false
			for _, service := range ServiceList {
				if item == service {
					isFingerOrSystem = true
					services = item
					break
				}
			}
			for _, system := range SystemList {
				if item == system {
					isFingerOrSystem = true
					break
				}
			}
			if isFingerOrSystem {
				continue
			}
			if !componentMap[item] {
				componentMap[item] = true
			}
		}
	}
	for k := range componentMap {
		components += "/" + k
	}
	itemMap["field3"] = services
	itemMap["field4"] = strings.Trim(components, "/")

	jsonResultMap["ip"] = fingerResult.Target
	jsonResultMap["protocol"] = "tcp"
	jsonResultMap["port"] = strconv.Itoa(fingerResult.Port)
	jsonResultMap["service"] = itemMap["field3"]
	jsonResultMap["component"] = itemMap["field4"]
	jsonResultMap["banner"] = ""
	jsonResultMap["title"] = ""
	jsonResultMap["response"] = ""
	jsonResultMap["cert"] = ""

	if fingerResult.Fingerprint != nil {
		jsonResultMap["banner"] = fingerResult.Fingerprint.Banner
	}

	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)

	resultMapList = append(resultMapList, itemMap)
	return resultMapList, nil
}

// handlePortScanSiteList 批量处理端口扫描站点信息
func (t *TaskTaskResult) handlePortScanSiteList(ctx context.Context, fingerResult MatchResult) ([]map[string]string, error) {
	resultMapList := make([]map[string]string, 0)
	// 只处理HTTP/HTTPS服务
	if fingerResult.Fingerprint == nil || !strings.Contains(fingerResult.Fingerprint.ServiceName, "http") {
		return resultMapList, nil
	}
	itemMap := make(map[string]string, 0)
	scheme := "http"
	if strings.Contains(fingerResult.Fingerprint.ServiceName, "https") {
		scheme = "https"
	}

	jsonResultMap := make(map[string]string, 0)
	itemMap["subObjType"] = enums.TaskResultSubObjTypeSite
	itemMap["field1"] = scheme + "://" + fingerResult.Target + ":" + strconv.Itoa(fingerResult.Port)
	itemMap["field2"] = strconv.Itoa(fingerResult.Port)
	itemMap["field3"] = fingerResult.Fingerprint.ServiceName
	itemMap["field4"] = ""

	response := ""
	jsonResultMap["target"] = itemMap["field1"]
	jsonResultMap["title"] = ""
	jsonResultMap["protocol"] = "tcp"
	jsonResultMap["port"] = strconv.Itoa(fingerResult.Port)
	jsonResultMap["service"] = fingerResult.Fingerprint.ServiceName
	jsonResultMap["response"] = response
	jsonResultMap["fingerprint"] = fingerResult.Fingerprint.ServiceName
	jsonResultMap["cert"] = ""
	jsonResultMap["screenshot"] = ""
	jsonResultMap["tag"] = ""
	jsonResultMap["waf"] = ""

	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)

	resultMapList = append(resultMapList, itemMap)
	return resultMapList, nil
}

// HandlePortScan 处理端口扫描信息
func (t *TaskTaskResult) handlePortScanService(ctx context.Context, detailMap map[string]string) (map[string]string, error) {
	itemMap := make(map[string]string)
	itemMap["subObjType"] = enums.TaskResultSubObjTypeService
	itemMap["field1"] = detailMap["host"]
	itemMap["field2"] = detailMap["port"]
	var services string
	jsonResultMap := make(map[string]string, 0)
	var components string
	var componentMap = make(map[string]bool)
	replaceReg, _ := regexp.Compile(`\[.+?]`)
	fingerPrintStr := replaceReg.ReplaceAllString(detailMap["fingerPrint"], "")
	for _, item := range strings.Split(fingerPrintStr, "/") {
		isFingerOrSystem := false
		//temp := strings.Split(item, "[")
		//if len(temp) == 2 {
		//	item = temp[0]
		//}
		for _, service := range ServiceList {
			if item == service {
				isFingerOrSystem = true
				services = item
				break
			}
		}
		for _, system := range SystemList {
			if item == system {
				isFingerOrSystem = true
				break
			}
		}
		if isFingerOrSystem {
			continue
		}
		if !componentMap[item] {
			componentMap[item] = true
		}
	}
	for k := range componentMap {
		components += "/" + k
	}
	itemMap["field3"] = services
	itemMap["field4"] = strings.Trim(components, "/")

	jsonResultMap["ip"] = detailMap["host"]
	jsonResultMap["protocol"] = "tcp"
	jsonResultMap["port"] = detailMap["port"]
	jsonResultMap["service"] = itemMap["field3"]
	jsonResultMap["component"] = itemMap["field4"]
	jsonResultMap["banner"] = ""
	jsonResultMap["title"] = detailMap["htmlTitle"]
	jsonResultMap["response"] = ""
	jsonResultMap["cert"] = ""
	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)
	return itemMap, nil
}

// HandlePortScan 处理端口扫描信息
func (t *TaskTaskResult) handleFingerPrintService(ctx context.Context, detailMap map[string]string, targetId int) (map[string]string, error) {
	var ttrModel mysqls.TaskTaskResult
	taskTaskResultList := ttrModel.GetTaskInfoByTargetIdAndSubObjType(ctx, targetId, enums.TaskResultSubObjTypeService)
	itemMap := make(map[string]string, 0)
	var fingerPrintResult FingerPrintResult
	err := json.Unmarshal([]byte(detailMap["app_info"]), &fingerPrintResult)
	if err != nil {
		return itemMap, nil
	}
	for _, result := range taskTaskResultList {
		if result.Field2 != fingerPrintResult.Port {
			continue
		}
		itemMap["subObjType"] = enums.TaskResultSubObjTypeService
		itemMap["field1"] = result.Field1
		itemMap["field2"] = result.Field2
		itemMap["field3"] = result.Field3
		itemMap["field4"] = result.Field4
		var jsonResultMap map[string]string
		err = json.Unmarshal([]byte(result.JSONResult), &jsonResultMap)
		if err != nil {
			fmt.Println(err)
		}
		// get components && save in componentMap
		var componentMap = make(map[string]bool)
		var componentStr = jsonResultMap["component"]
		components := strings.Split(componentStr, "/")
		for _, component := range components {
			componentMap[component] = true
		}
		for _, targetInfo := range fingerPrintResult.TargetInfo {
			isSystem := false
			for _, system := range SystemList {
				if strings.ToLower(targetInfo.AppName) == system {
					isSystem = true
					break
				}
			}
			if isSystem {
				continue
			}
			// repeat filter by componentMap
			if !componentMap[targetInfo.AppName] {
				componentMap[targetInfo.AppName] = true
				componentStr = componentStr + "/" + targetInfo.AppName
			}
		}
		jsonResultMap["component"] = componentStr
		jsonResultMapByte, _ := json.Marshal(jsonResultMap)
		itemMap["jsonResult"] = string(jsonResultMapByte)
		err = t.UpdateTaskTaskResult(ctx, result.ID, result.ObjID, result.SubObjID, itemMap)
		if err != nil {
			return itemMap, err
		}
	}
	return itemMap, nil
}

var ServiceList = []string{"http", "https", "ssh", "ftp", "smb", "rdp", "redis", "mysql", "oracle", "mongodb", "mssql", "http_server", "domain", "tcp", "udp", "microsoft-ds", "netbios-ssn", "msrpc"}
var SystemList = []string{"ubuntu", "debian", "linux", "windows", "debian_linux", "linux_kernel", "unix", "centos"}

// HandlePortScan 处理端口扫描信息
func (t *TaskTaskResult) handlePortScanSite(ctx context.Context, detailMap map[string]string) (map[string]string, error) {
	itemMap := make(map[string]string, 0)
	jsonResultMap := make(map[string]string, 0)
	itemMap["subObjType"] = enums.TaskResultSubObjTypeSite
	itemMap["field1"] = detailMap["service"] + "://" + detailMap["host"] + ":" + detailMap["port"]
	itemMap["field2"] = detailMap["port"]
	itemMap["field3"] = detailMap["fingerprint"]
	itemMap["field4"] = ""

	// 先注释掉此处代码，此处代码会进行页面请求， 会产生阻塞后续信息消费的问题
	//response, err := t.getHttpResponse(ctx, itemMap["field1"])
	//if err != nil {
	//	log.Error(err.Error())
	//}
	response := ""

	jsonResultMap["target"] = itemMap["field1"]
	jsonResultMap["title"] = detailMap["htmlTitle"]
	jsonResultMap["protocol"] = "tcp"
	jsonResultMap["port"] = detailMap["port"]
	jsonResultMap["service"] = detailMap["service"]
	jsonResultMap["response"] = response
	jsonResultMap["fingerprint"] = detailMap["fingerprint"]
	jsonResultMap["cert"] = ""
	jsonResultMap["screenshot"] = ""
	jsonResultMap["tag"] = ""
	jsonResultMap["waf"] = ""
	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)
	return itemMap, nil
}

func (t *TaskTaskResult) getHttpResponse(ctx context.Context, url string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(1 * time.Second)}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	headers := "HTTP/1.1 " + resp.Status + "\n"
	for key, value := range resp.Header {
		headers += key + ": "
		for _, val := range value {
			headers += val + ","
		}
		headers = strings.Trim(headers, ",") + "\n"
	}
	return headers, err
}

// handleCrawlerUrl 处理爬虫发现的URL信息
func (t *TaskTaskResult) handleCrawlerUrl(ctx context.Context, detailMap map[string]string) (map[string]string, error) {
	itemMap := make(map[string]string, 0)
	itemMap["subObjType"] = enums.TaskResultSubObjTypeUrl
	itemMap["field1"] = detailMap["url"]
	itemMap["field2"] = detailMap["method"]
	itemMap["field3"] = detailMap["status"]
	itemMap["field4"] = ""

	jsonResultMap := make(map[string]string, 0)
	jsonResultMap["url"] = detailMap["url"]
	jsonResultMap["method"] = detailMap["method"]
	jsonResultMap["status"] = detailMap["status"]
	jsonResultMap["title"] = detailMap["title"]
	jsonResultMap["param"] = detailMap["param"]
	jsonResultMap["bodyLength"] = detailMap["bodyLength"]
	jsonResultMap["contentType"] = detailMap["contentType"]

	// 优先使用 scanner 返回的 risk_type，其次使用 URL 启发式判断
	jsonResultMap["riskType"] = detailMap["riskType"]
	if jsonResultMap["riskType"] == "" {
		if strings.Contains(detailMap["url"], "login") {
			jsonResultMap["riskType"] = "登录入口"
		} else if strings.Contains(detailMap["url"], "upload") || strings.Contains(detailMap["url"], "evil") || strings.Contains(detailMap["url"], "ckstr") {
			jsonResultMap["riskType"] = "上传入口"
		}
	}

	jsonResultMap["source"] = "爬虫"
	if detailMap["crawlerMode"] != "" {
		jsonResultMap["crawlerMode"] = detailMap["crawlerMode"]
	}
	if detailMap["targetURL"] != "" {
		jsonResultMap["targetURL"] = detailMap["targetURL"]
	}
	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)
	return itemMap, nil
}

// handleCrawlerUrl 处理端口扫描信息
func (t *TaskTaskResult) handleWebDIrPathScanUrl(ctx context.Context, detailMap map[string]string) (map[string]string, error) {
	itemMap := make(map[string]string, 0)
	itemMap["subObjType"] = enums.TaskResultSubObjTypeUrl
	itemMap["field1"] = detailMap["url"]
	itemMap["field2"] = detailMap["method"]
	itemMap["field3"] = detailMap["status"]
	itemMap["field4"] = ""

	jsonResultMap := make(map[string]string, 0)
	jsonResultMap["url"] = detailMap["url"]
	jsonResultMap["method"] = detailMap["method"]
	jsonResultMap["status"] = detailMap["status"]
	jsonResultMap["title"] = detailMap["title"]
	jsonResultMap["param"] = detailMap["param"]
	jsonResultMap["riskType"] = ""
	jsonResultMap["bodyLength"] = detailMap["bodyLength"]
	jsonResultMap["source"] = "路径爆破"
	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)
	return itemMap, nil
}

// handleCDN 处理cdn信息
func (t *TaskTaskResult) handleCDN(ctx context.Context, detailMap map[string]string) (map[string]string, error) {
	itemMap := make(map[string]string, 0)
	itemMap["field1"] = detailMap["domain"]
	itemMap["field2"] = detailMap["exists"]
	itemMap["field3"] = detailMap["resovleIP"]
	itemMap["field4"] = ""
	itemMap["subObjType"] = enums.TaskResultSubObjTypeCDN
	jsonResultMap := make(map[string]string, 0)
	jsonResultMap["domain"] = detailMap["domain"]
	jsonResultMap["exists"] = detailMap["exists"]
	jsonResultMap["resovleIP"] = detailMap["resovleIP"]
	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)
	return itemMap, nil
}

// handleWhoIs 处理whois信息
func (t *TaskTaskResult) handleWhois(ctx context.Context, detailMap map[string]string) (map[string]string, error) {
	itemMap := make(map[string]string, 0)
	itemMap["subObjType"] = ""
	itemMap["field1"] = detailMap["dataList"]
	itemMap["field2"] = ""
	itemMap["field3"] = ""
	itemMap["field4"] = ""

	jsonResultMap := make(map[string]string, 0)
	jsonResultMap["dataList"] = detailMap["dataList"]
	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)
	return itemMap, nil
}

// handleSubdomain 处理子域名爆破信息
func (t *TaskTaskResult) handleSubdomain(ctx context.Context, detailMap map[string]string) (map[string]string, error) {
	itemMap := make(map[string]string, 0)
	itemMap["subObjType"] = enums.TaskResultSubObjTypeSubdomain
	itemMap["field1"] = detailMap["domain"]
	itemMap["field2"] = detailMap["ip"]
	itemMap["field3"] = ""
	itemMap["field4"] = ""

	jsonResultMap := make(map[string]string, 0)
	jsonResultMap["target"] = detailMap["target"]
	jsonResultMap["subDomain"] = detailMap["domain"]
	jsonResultMap["recordType"] = detailMap["record_type"]
	jsonResultMap["recordValue"] = detailMap["record_value"]
	jsonResultMap["relateIP"] = detailMap["ip"]
	jsonResultMap["source"] = "子域名爆破"
	jsonResultMap["cdn"] = ""
	jsonResultMapByte, _ := json.Marshal(jsonResultMap)
	itemMap["jsonResult"] = string(jsonResultMapByte)
	return itemMap, nil
}

// handleWafDetect 处理waf检测结果信息
func (t *TaskTaskResult) handleWafDetect(ctx context.Context, detailMap map[string]string, targetId int) (map[string]string, error) {
	if detailMap["waf_result"] == "" {
		return nil, nil
	}
	var ttrModel mysqls.TaskTaskResult
	taskTaskResultList := ttrModel.GetTaskInfoByTargetIdAndSubObjType(ctx, targetId, enums.TaskResultSubObjTypeSite)
	itemMap := make(map[string]string, 0)
	_, _, port, err := network.ParseUrl(detailMap["target"])
	if err != nil {
		return nil, err
	}
	for _, result := range taskTaskResultList {
		if result.Field2 != port {
			continue
		}
		itemMap["subObjType"] = enums.TaskResultSubObjTypeSite
		itemMap["field1"] = result.Field1
		itemMap["field2"] = result.Field2
		itemMap["field3"] = result.Field3
		itemMap["field4"] = result.Field4
		var jsonResultMap map[string]string
		err := json.Unmarshal([]byte(result.JSONResult), &jsonResultMap)
		if err != nil {
			fmt.Println(err)
		}
		jsonResultMap["waf"] = detailMap["waf_result"]
		jsonResultMapByte, _ := json.Marshal(jsonResultMap)
		itemMap["jsonResult"] = string(jsonResultMapByte)
	}
	return itemMap, nil
}

// AddTaskTaskResult 插入一条检测漏洞
func (t *TaskTaskResult) AddTaskTaskResult(ctx context.Context, objId, subObjID string, resultMapList []map[string]string) error {
	// 1. Group results by subObjType
	resultsByType := make(map[string][]map[string]string)
	for _, res := range resultMapList {
		typ := res["subObjType"]
		resultsByType[typ] = append(resultsByType[typ], res)
	}

	targetIDInt, _ := strconv.Atoi(subObjID)

	// 2. Process each type
	for typ, results := range resultsByType {
		// Fetch existing results for this target and type to prevent duplicates
		var model mysqls.TaskTaskResult
		existingList := model.GetTaskInfoByTargetIdAndSubObjType(ctx, targetIDInt, typ)

		// Create a map for quick lookup
		existingMap := make(map[string]bool)
		for _, existing := range existingList {
			// Unique key: Field1 + "|" + Field2 (e.g., IP|Port)
			key := existing.Field1 + "|" + existing.Field2
			existingMap[key] = true
		}

		for _, resultMap := range results {
			if resultMap["field1"] == "" {
				continue
			}

			// Check if already exists
			key := resultMap["field1"] + "|" + resultMap["field2"]
			if existingMap[key] {
				// Optional: Update existing record if needed, but for now just skip to avoid duplicates
				continue
			}

			var taskTaskResult = mysqls.TaskTaskResult{
				ObjType:    enums.TaskResultObjTypeInfo,
				SubObjType: resultMap["subObjType"],
				ObjID:      objId,
				SubObjID:   subObjID,
				Field1:     resultMap["field1"],
				Field2:     resultMap["field2"],
				Field3:     resultMap["field3"],
				Field4:     resultMap["field4"],
				JSONResult: resultMap["jsonResult"],
				CreateTime: time.Now(),
			}
			err := taskTaskResult.AddTaskTaskResult(ctx)
			if err != nil {
				log.Error("AddTaskTaskResult error: " + err.Error())
			} else {
				// Add to map to prevent duplicates within the same batch
				existingMap[key] = true
			}
		}
	}
	return nil
}

// UpdateTaskTaskResult 更新一条检测漏洞
func (t *TaskTaskResult) UpdateTaskTaskResult(ctx context.Context, id int, objId, subObjID string, resultMap map[string]string) error {
	var taskTaskResult = mysqls.TaskTaskResult{
		ID:         id,
		ObjType:    enums.TaskResultObjTypeInfo,
		SubObjType: resultMap["subObjType"],
		ObjID:      objId,
		SubObjID:   subObjID,
		Field1:     resultMap["field1"],
		Field2:     resultMap["field2"],
		Field3:     resultMap["field3"],
		Field4:     resultMap["field4"],
		JSONResult: resultMap["jsonResult"],
		CreateTime: time.Now(),
	}
	err := taskTaskResult.UpdateTaskTaskResult(ctx)
	if err != nil {
		log.Error("AddTaskTaskResult error: " + err.Error())
	}
	return nil
}

// TaskResultList 信息收集列表及筛选
func (t *TaskTaskResult) TaskResultList(ctx context.Context, objType int, subObjType string, objId string, search string, page, size int) ([]mysqls.TaskTaskResult, int64, error) {
	var taskResult mysqls.TaskTaskResult
	return taskResult.GetTaskTaskResultList(ctx, objType, subObjType, objId, search, page, size)
}

// CheckCNDTaskResultBySubObjID  通过ID检查信息收集-站点-cdn信息
func (t *TaskTaskResult) CheckCNDTaskResultBySubObjID(ctx context.Context, objId, domain string) string {
	var taskResult mysqls.TaskTaskResult
	cdnRes := taskResult.GetByTaskIdAndSubObjID(ctx, objId, enums.TaskResultSubObjTypeCDN)
	if cdnRes.ID != 0 && cdnRes.Field2 != "false" {
		if strings.Contains(domain, cdnRes.Field1) {
			return cdnRes.Field3
		}
	}
	return ""
}

// DelTargetByIds 删除目标
func (t *TaskTaskResult) DelTaskResult(ctx context.Context, ids any) error {
	var taskResult mysqls.TaskTaskResult
	return taskResult.DeleteTaskTaskResult(ctx, ids)
}

// TaskResultDetail 获取任务结果详情
func (t *TaskTaskResult) TaskResultDetail(ctx context.Context, id int) (mysqls.TaskTaskResult, error) {
	var taskResult = mysqls.TaskTaskResult{
		ID: id,
	}
	return taskResult.GetTaskTaskResult(ctx)
}

// 批量删除 通过task_id
func (t *TaskTaskResult) DelTaskInfoByTaskId(ctx context.Context, taskIds []int) error {
	var taskTaskResultModel mysqls.TaskTaskResult
	return taskTaskResultModel.DeleteByTaskIds(ctx, taskIds)
}

// 依据任务ID与数据子类型获取 信息收集类 数据
func (t *TaskTaskResult) GetTaskInfoByTargetIdAndSubObjType(ctx context.Context, targetId int, subObjType string) []mysqls.TaskTaskResult {
	var taskTaskResultModel mysqls.TaskTaskResult
	return taskTaskResultModel.GetTaskInfoByTargetIdAndSubObjType(ctx, targetId, subObjType)
}

// 依据任务ID与数据子类型获取 信息收集类 数据
func (t *TaskTaskResult) GetByTaskIdAndSubObjType(ctx context.Context, taskId int, subObjType string) []mysqls.TaskTaskResult {
	var taskTaskResultModel mysqls.TaskTaskResult
	return taskTaskResultModel.GetByTaskIdAndSubObjType(ctx, taskId, subObjType)
}

// 依据目标IDs获取所有目标信息收集数据，并返回map
func (t *TaskTaskResult) AllByTargetIds(ctx context.Context, targetIds []int) (returnData map[string][]mysqls.TaskTaskResult) {
	var taskTaskResultModel mysqls.TaskTaskResult
	taskTaskResult, _ := taskTaskResultModel.GetTaskTaskResultByType(ctx, enums.TaskResultObjTypeInfo, enums.TaskResultSubObjTypeService, targetIds)
	returnData = make(map[string][]mysqls.TaskTaskResult)
	for _, item := range taskTaskResult {
		// 如果不存在，则初始化
		if _, ok := returnData[item.SubObjID]; !ok {
			returnData[item.SubObjID] = make([]mysqls.TaskTaskResult, 0)
		}
		returnData[item.SubObjID] = append(returnData[item.SubObjID], item)
	}
	return
}
