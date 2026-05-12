package services

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"regexp"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"
)

type ReportVerify struct {
}

// AnalysisTianJingVul 分析启明天境报告
func (r *ReportVerify) AnalysisTianJingVul(ctx context.Context, content string) map[string][]map[string]interface{} {
	itemPattern := regexp.MustCompile(`collapse-content[\s\S]*?/div`)
	cvePattern := regexp.MustCompile(`CVE-[0-9]{4}-[0-9]{4,}`)
	namePattern := regexp.MustCompile(`tableTitle="\d*(.*?)"`)
	riskPattern := regexp.MustCompile(`vertical-m(.*?)<`)
	ipPattern := regexp.MustCompile(`<td>(?:\d{1,3}\.){3}\d{1,3}</td>`)
	descPattern := regexp.MustCompile(`简单描述[\s\S]*?<td>([\s\S]*?)</td>`)
	fixPattern := regexp.MustCompile(`修补建议[\s\S]*?<td>([\s\S]*?)</td>`)
	vulList := itemPattern.FindAllString(content, 300)
	vulMap := make(map[string][]map[string]interface{}, 0)

	for _, vul := range vulList {
		riskStr := riskPattern.FindString(vul)
		if riskStr == "" || strings.Contains(riskStr, ">信息<") {
			continue
		}
		cve := cvePattern.FindString(vul)
		nameList := namePattern.FindStringSubmatch(vul)
		descList := descPattern.FindStringSubmatch(vul)
		fixList := fixPattern.FindStringSubmatch(vul)
		ip := strings.Trim(strings.Trim(ipPattern.FindString(vul), "<td>"), "</td>")
		//name = strings.Join(strings.Split(strings.Trim(name, `"`), " ")[1:], " ")
		if len(nameList) != 2 || len(descList) != 2 || len(fixList) != 2 {
			continue
		}
		var risk int
		if strings.Contains(riskStr, "高危险") {
			risk = enums.VulLibrariesRiskHigh
		} else if strings.Contains(riskStr, "中危险") {
			risk = enums.VulLibrariesRiskMiddle
		} else if strings.Contains(riskStr, "低危险") {
			risk = enums.VulLibrariesRiskLow
		} else {
			continue
		}
		tempMap := map[string]interface{}{
			"name": strings.TrimSpace(nameList[1]),
			"ip":   ip,
			"risk": risk,
			"cve":  cve,
			"desc": strings.TrimSpace(strings.ReplaceAll(descList[1], "\n", "")),
			"fix":  strings.TrimSpace(strings.ReplaceAll(fixList[1], "\n", "")),
		}
		if vulMap[ip] == nil {
			vulMap[ip] = make([]map[string]interface{}, 0)
		}
		vulMap[ip] = append(vulMap[ip], tempMap)
	}
	return vulMap
}

// AnalysisTianJingPort 分析启明天境报告的端口
func (r *ReportVerify) AnalysisTianJingPort(ctx context.Context, content string) map[string]string {
	itemPattern := regexp.MustCompile(`tableTitle="4.1.5 端口服务Top10"[\s\S]*?/table`)
	trPattern := regexp.MustCompile(`<tr>[\s\S]*?</tr>`)
	portPattern := regexp.MustCompile(`<td[\s\S]*?</td>`)
	ipPattern := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)
	portInfo := itemPattern.FindString(content)
	portMap := make(map[string]string, 0)
	trList := trPattern.FindAllString(portInfo, 100)
	if len(trList) < 2 {
		return portMap
	}
	for _, tr := range trList[1:] {
		if !strings.Contains(tr, ">") || !strings.Contains(tr, "<") {
			continue
		}
		port := strings.Split(strings.Split(portPattern.FindString(tr), ">")[1], "<")[0]
		ip := ipPattern.FindString(tr)
		portMap[ip] += port + ","
	}
	return portMap
}

// AnalysisNsFocus 分析绿盟vulFocus报告
func (r *ReportVerify) AnalysisNsFocus(ctx context.Context, content string) map[string][]map[string]interface{} {
	itemPattern := regexp.MustCompile(`[even|odd] vuln_[\s\S]*?/table`)
	vulList := itemPattern.FindAllString(content, 300)
	cvePattern := regexp.MustCompile(`CVE-[0-9]{4}-[0-9]{4,}`)
	namePattern := regexp.MustCompile(`<span.*?>(.*?)</span>`)
	riskPattern := regexp.MustCompile(`vuln_.*.gif`)
	ipPattern := regexp.MustCompile(`<td width="80%">([\s\S]*?)<a href=`)
	descPattern := regexp.MustCompile(`详细描述[\s\S]*?<td>([\s\S]*?)</td>`)
	fixPattern := regexp.MustCompile(`解决办法[\s\S]*?<td>([\s\S]*?)</td>`)
	cvssPattern := regexp.MustCompile(`CVSS评分[\s\S]*?<td>([\s\S]*?)</td>`)
	vulMap := make(map[string][]map[string]interface{}, 0)

	for _, vul := range vulList {
		riskStr := strings.Trim(strings.Trim(riskPattern.FindString(vul), "vuln_"), ".gif")
		if riskStr == "" {
			continue
		}
		var risk int
		if strings.Contains(riskStr, "high") {
			risk = enums.VulLibrariesRiskHigh
		} else if strings.Contains(riskStr, "middle") {
			risk = enums.VulLibrariesRiskMiddle
		} else if strings.Contains(riskStr, "low") {
			risk = enums.VulLibrariesRiskLow
		} else {
			continue
		}
		cve := cvePattern.FindString(vul)
		nameList := namePattern.FindStringSubmatch(vul)
		ipList := ipPattern.FindStringSubmatch(vul)
		descList := descPattern.FindStringSubmatch(vul)
		fixList := fixPattern.FindStringSubmatch(vul)
		cvssList := cvssPattern.FindStringSubmatch(vul)
		if len(nameList) != 2 || len(ipList) != 2 || len(descList) != 2 || len(fixList) != 2 || len(cvssList) != 2 {
			continue
		}
		for _, ip := range strings.Split(ipList[1], ";&nbsp") {
			ip = strings.TrimSpace(strings.ReplaceAll(ip, "\n", ""))
			if ip == "" {
				continue
			}
			tempMap := map[string]interface{}{
				"name": nameList[1],
				"ip":   ip,
				"risk": risk,
				"cve":  cve,
				"desc": descList[1],
				"fix":  fixList[1],
				"cvss": cvssList[1],
			}
			if vulMap[ip] == nil {
				vulMap[ip] = make([]map[string]interface{}, 0)
			}
			vulMap[ip] = append(vulMap[ip], tempMap)
		}
	}
	return vulMap
}

// AnalysisNsPort 分析绿盟vulFocus报告中的端口信息
func (r *ReportVerify) AnalysisNsPort(ctx context.Context, vulMap map[string][]map[string]interface{}) map[string]string {
	portMap := make(map[string]string, 0)
	var portScan enums.PortScanConfig
	portScanEnum := portScan.AllPortScanTypeValue()
	for ip, _ := range vulMap {
		portMap[ip] = portScanEnum[enums.TaskConfigurationPortScanTypeTop100]
	}
	return portMap
}

// AnalysisTianJingPort 分析启明天境报告的端口
func (r *ReportVerify) AnalysisTianJingPort2025(ctx context.Context, content string) map[string]string {
	portMap := make(map[string]string, 0)
	// 提取IP地址
	ipPattern := regexp.MustCompile(`主机地址[\s\S]*?<td>([\s\S]*?)</td>`)
	ipMatch := ipPattern.FindStringSubmatch(content)
	if len(ipMatch) < 2 {
		return portMap
	}
	ip := strings.TrimSpace(ipMatch[1])
	// 提取端口信息
	portSectionPattern := regexp.MustCompile(`<div class="title4">[\s\S]*?端口详情[\s\S]*?<table>[\s\S]*?excel专用`)
	portSection := portSectionPattern.FindString(content)
	// 提取每个端口行
	trPattern := regexp.MustCompile(`<tr class="con-item-title">[\s\S]*?<td class="pl-s">([\d]+)</td>[\s\S]*?</tr>`)
	portMatches := trPattern.FindAllStringSubmatch(portSection, -1)
	var ports []string
	tempportMap := make(map[string]bool) // 用于去重的map
	for _, match := range portMatches {
		if len(match) >= 2 {
			port := match[1]
			if port != "0" && !tempportMap[port] { // 排除端口0和重复端口
				ports = append(ports, port)
				tempportMap[port] = true // 标记该端口已添加
			}
		}
	}
	if len(ports) > 0 {
		portMap[ip] = strings.Join(ports, ",")
	}
	return portMap
}

// AnalysisTianJingVul2025 分析启明天境2025版报告
func (r *ReportVerify) AnalysisTianJingVul2025(ctx context.Context, content string) map[string][]map[string]interface{} {
	vulMap := make(map[string][]map[string]interface{}, 0)

	// 提取IP地址
	ipPattern := regexp.MustCompile(`主机地址[\s\S]*?<td>([\s\S]*?)</td>`)
	ipMatch := ipPattern.FindStringSubmatch(content)
	if len(ipMatch) < 2 {
		return vulMap
	}
	ip := strings.TrimSpace(ipMatch[1])

	// 提取各个风险级别的漏洞
	riskSections := []struct {
		pattern string
		risk    int
	}{
		{`<!-- 超危漏洞 -->[\s\S]*?<!-- 高危漏洞 -->`, enums.VulLibrariesRiskDead},
		{`<!-- 高危漏洞 -->[\s\S]*?<!-- 中危漏洞 -->`, enums.VulLibrariesRiskHigh},
		{`<!-- 中危漏洞 -->[\s\S]*?<!-- 低危漏洞 -->`, enums.VulLibrariesRiskMiddle},
		{`<!-- 低危漏洞 -->[\s\S]*?<!-- 信息漏洞 -->`, enums.VulLibrariesRiskLow},
	}

	for _, riskSection := range riskSections {
		sectionPattern := regexp.MustCompile(riskSection.pattern)
		section := sectionPattern.FindString(content)

		// 提取每个漏洞条目
		vulItemPattern := regexp.MustCompile(`<div class="collapse-item">[\s\S]*?<div class="collapse-content">[\s\S]*?</div>\s*</div>`)
		vulItems := vulItemPattern.FindAllString(section, -1)

		for _, item := range vulItems {
			// 提取CVE编号
			cvePattern := regexp.MustCompile(`CVE-\d{4}-\d{4,}`)
			cve := cvePattern.FindString(item)

			// 提取漏洞名称
			namePattern := regexp.MustCompile(`<span class="text-(?:v|h|m|low)">(.*?)</span>`)
			nameMatches := namePattern.FindAllStringSubmatch(item, -1)
			var name string
			if len(nameMatches) == 1 && len(nameMatches[0]) == 2 {
				name = strings.TrimSpace(nameMatches[0][1])
			}

			// 查找漏洞详情
			if name != "" && cve != "" {
				// 查找漏洞的详细描述和修复建议
				var desc, fix string
				descPattern := regexp.MustCompile(`<span class="text-h">` + regexp.QuoteMeta(name) + `</span>[\s\S]*?简单描述[\s\S]*?<td>([\s\S]*?)</td>`)
				descMatches := descPattern.FindStringSubmatch(content)
				if len(descMatches) >= 2 {
					desc = strings.TrimSpace(strings.ReplaceAll(descMatches[1], "\n", ""))
				}

				fixPattern := regexp.MustCompile(`<span class="text-h">` + regexp.QuoteMeta(name) + `</span>[\s\S]*?排查方法及修复建议[\s\S]*?<td>([\s\S]*?)</td>`)
				fixMatches := fixPattern.FindStringSubmatch(content)

				if len(fixMatches) >= 2 {
					fix = strings.TrimSpace(strings.ReplaceAll(fixMatches[1], "\n", ""))
				}

				tempMap := map[string]interface{}{
					"name": name,
					"ip":   ip,
					"risk": riskSection.risk,
					"cve":  cve,
					"desc": desc,
					"fix":  fix,
				}

				if vulMap[ip] == nil {
					vulMap[ip] = make([]map[string]interface{}, 0)
				}
				vulMap[ip] = append(vulMap[ip], tempMap)
			}
		}
	}
	return vulMap
}

// AnalysisProducer 分析报告厂商
func (r *ReportVerify) AnalysisProducer(ctx context.Context, content string) int {
	if strings.Contains(content, `xmlns:td="http://www.tdymeleaf.org"`) {
		if strings.Contains(content, `<script src="./js/echarts.js"></script>`) || strings.Contains(content, `<script src="../js/common-detail.js"></script>`) {
			return enums.ReportVerifyProducerTianJing2025
		} else {
			return enums.ReportVerifyProducerTianJing
		}
	} else if strings.Contains(content, `绿盟科技&#34;远程安全评估系统&#34;安全评估报告`) {
		return enums.ReportVerifyProducerNsfocus
	}
	return enums.ReportVerifyProducerUnKnown
}

// AnalysisFileType 分析文件类型
func (r *ReportVerify) AnalysisFileType(ctx context.Context, filename string) int {
	if strings.HasSuffix(filename, ".html") {
		return enums.ReportVerifyFileTypeHtml
	} else if strings.HasSuffix(filename, ".zip") {
		return enums.ReportVerifyFileTypeZip
	}
	return enums.ReportVerifyFileTypeUnKnown
}

// SaveTask 保存报告验证任务
func (r *ReportVerify) SaveTask(ctx context.Context, taskName string, fileInfoMap map[string]string, producer, userId int) (int, error) {
	fileInfoByte, _ := json.Marshal(fileInfoMap)
	var task = mysqls.Reportverifytask{
		Name:        taskName,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		ExecuteType: enums.TaskExecTypeImmediate,
		Producer:    producer,
		User:        userId,
		Status:      enums.TaskStatusBegin,
		Risk:        enums.TaskRiskSafe,
		Overview:    "",
		IsStats:     enums.TaskIsStatsNo,
		Fileinfo:    string(fileInfoByte),
		Exp:         0,
		Verify:      0,
		Failed:      0,
		Unverify:    0,
		ExecuteTime: time.Now(),
	}
	err := task.AddReportverifytask(ctx)
	if err != nil {
		return task.ID, err
	}
	return task.ID, nil
}

// SaveTarget 保存报告验证目标
func (r *ReportVerify) SaveTarget(ctx context.Context, taskId int, portMap map[string]string, vulMap map[string][]map[string]interface{}) error {
	for ip, _ := range vulMap {
		var ports string
		if portMap[ip] != "" {
			ports = strings.Trim(portMap[ip], ",")
		} else {
			var portScanConfig enums.PortScanConfig
			portEnum := portScanConfig.AllPortScanTypeValue()
			ports = portEnum[enums.TaskConfigurationPortScanTypeTop1000]
		}
		var target = mysqls.Reportverifytarget{
			TaskId:       taskId,
			Target:       ip,
			Os:           "",
			Risk:         enums.TargetRiskLowNoFound,
			Status:       enums.TargetStatusToBegin,
			AnalysisData: ports,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err := target.AddReportverifytarget(ctx)
		if err != nil {
			log.Error("save target error: " + err.Error())
		}
	}
	return nil
}

func reportVerifyStringField(vul map[string]interface{}, key string) (string, bool) {
	value, ok := vul[key]
	if !ok || value == nil {
		return "", false
	}
	switch v := value.(type) {
	case string:
		return v, true
	default:
		return fmt.Sprint(v), true
	}
}

func reportVerifyRiskField(vul map[string]interface{}) (int, bool) {
	value, ok := vul["risk"]
	if !ok || value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case json.Number:
		risk, err := v.Int64()
		if err != nil {
			return 0, false
		}
		return int(risk), true
	case string:
		risk, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0, false
		}
		return risk, true
	default:
		return 0, false
	}
}

func truncateReportVerifyField(value string) string {
	if len(value) > 996 {
		return value[:996]
	}
	return value
}

// SaveVul 保存报告验证漏洞
func (r *ReportVerify) SaveVul(ctx context.Context, taskId int, vulMap map[string][]map[string]interface{}) error {
	for ip, vuls := range vulMap {
		for _, vul := range vuls {
			name, ok := reportVerifyStringField(vul, "name")
			if !ok || strings.TrimSpace(name) == "" {
				log.Warnf("SaveVul skip invalid vul, missing name, ip: %s", ip)
				continue
			}
			risk, ok := reportVerifyRiskField(vul)
			if !ok {
				log.Warnf("SaveVul use default risk, invalid risk value, ip: %s, name: %s", ip, name)
				risk = enums.VulLibrariesRiskInfo
			}
			cve, _ := reportVerifyStringField(vul, "cve")
			desc, _ := reportVerifyStringField(vul, "desc")
			fix, _ := reportVerifyStringField(vul, "fix")
			var rvVul = mysqls.Reportverifyvul{
				TaskId:   taskId,
				Name:     strings.TrimSpace(name),
				Risk:     risk,
				Status:   enums.ReportVerifyStatusUnVerify,
				Location: ip,
				Cve:      strings.TrimSpace(cve),
				Desc:     truncateReportVerifyField(desc),
				Fix:      truncateReportVerifyField(fix),
				//Cvss:     vul["cvss"].(string),
			}
			err := rvVul.AddReportverifyvul(ctx)
			if err != nil {
				log.Error("SaveVul: ", err)
			}
		}
	}
	return nil
}

// SavePort 保存报告验证端口
func (r *ReportVerify) SavePort(ctx context.Context, content string) error {
	return nil
}

// TaskList 报告验证任务列表
func (r *ReportVerify) TaskList(ctx context.Context, page, size, risk, producer int, startTime, endTime, search string, userIdList []int) ([]mysqls.Reportverifytask, int64, error) {
	var model mysqls.Reportverifytask
	list, total, err := model.GetReportverifytaskList(ctx, page, size, risk, producer, startTime, endTime, search, userIdList)
	return list, total, err
}

// TaskDetail 报告验证任务详情
func (r *ReportVerify) TaskDetail(ctx context.Context, id int) (mysqls.Reportverifytask, error) {
	var model = mysqls.Reportverifytask{
		ID: id,
	}
	return model.GetReportverifytask(ctx)
}

// TargetList 报告验证目标列表
func (r *ReportVerify) TargetList(ctx context.Context, taskId, risk, page, size int, search string) ([]mysqls.Reportverifytarget, int64, error) {
	var model mysqls.Reportverifytarget
	list, total, err := model.GetReportverifytargetList(ctx, taskId, risk, page, size, search)
	return list, total, err
}

// PortList 报告验证端口列表
func (r *ReportVerify) PortList(ctx context.Context, taskId, page, size int, search string) ([]mysqls.Reportverifyport, int64, error) {
	var model mysqls.Reportverifyport
	list, total, err := model.GetReportverifyportList(ctx, taskId, page, size, search)
	return list, total, err
}

// VulList 报告验证漏洞列表
func (r *ReportVerify) VulList(ctx context.Context, taskId, page, size, risk, status int, search string) ([]mysqls.Reportverifyvul, int64, error) {
	var model mysqls.Reportverifyvul
	list, total, err := model.GetReportverifyvulList(ctx, taskId, page, size, risk, status, search)
	return list, total, err
}

// GetOneWaitTarget 获取一个正在等待的报告验证任务
func (r *ReportVerify) GetOneWaitTarget(ctx context.Context) (mysqls.Reportverifytarget, error) {
	var model mysqls.Reportverifytarget
	return model.GetOneWaitTarget(ctx)
}

// ProducerEnum 获取报告验证任务
func (r *ReportVerify) ProducerEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.ReportVerifyEnum.AllReportVerifyTypeEnum())
}

// VulStatusEnum 获取报告获取状态
func (r *ReportVerify) VulStatusEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.ReportVerifyEnum.AllReportVerifyStatusEnum())
}

// ExecuteTypeEnum 执行类型枚举
func (r *ReportVerify) ExecuteTypeEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(map[int]string{
		enums.TaskExecTypeImmediate: "即时任务",
	})
}

// ExecuteTypeEnum 执行类型枚举
func (r *ReportVerify) VulRiskEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(map[int]string{
		enums.VulLibrariesRiskDead:   "致命",
		enums.VulLibrariesRiskHigh:   "高危",
		enums.VulLibrariesRiskMiddle: "中危",
		enums.VulLibrariesRiskLow:    "低危",
	})
}

// PortScan 端口扫描调用
func (r *ReportVerify) PortScan(ctx context.Context, ip, toScanPort string) ([]map[string]string, []string, error) {
	vulParam := make([]map[string]interface{}, 0)
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "ip",
		"value": ip,
	})
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "toScanPort",
		"value": toScanPort,
	})
	callId, err := httpclients.DecisionScriptCall(enums.ScriptNamePortScan, vulParam)
	var tempList []string
	if err != nil {
		return nil, tempList, err
	}
	var portScanResult string
	for {
		time.Sleep(5 * time.Second)
		end, callResult, err := httpclients.DecisionScriptCallResult(callId, false)
		if err != nil {
			return nil, tempList, err
		}
		if !end {
			continue
		}
		portScanResult = callResult
		break
	}
	portResultList := make([]map[string]string, 0)
	for _, result := range strings.Split(portScanResult, "\n") {
		var (
			tempMap map[string]string
			rootUrl string
		)
		err = json.Unmarshal([]byte(result), &tempMap)
		if err != nil {
			continue
		}
		portResultList = append(portResultList, tempMap)
		if tempMap["http_port"] != "" {
			rootUrl = "http://" + ip + ":" + tempMap["http_port"]
		} else if tempMap["https_port"] != "" {
			rootUrl = "https://" + ip + ":" + tempMap["https_port"]
		} else {
			continue
		}
		tempList = append(tempList, rootUrl)
	}
	return portResultList, tempList, nil
}

// UpdateTargetStatus 更新目标状态
func (r *ReportVerify) UpdateTargetStatus(ctx context.Context, id, status int) error {
	var model mysqls.Reportverifytarget
	return model.UpdateReportverifytargetStatus(ctx, id, status)
}

// UpdateTaskStatus 更新任务状态
func (r *ReportVerify) UpdateTaskStatus(ctx context.Context, id, status int) error {
	var model mysqls.Reportverifytask
	return model.UpdateReportverifytaskStatus(ctx, id, status)
}

// GetReportVerifyVulListCanCall 获取报告验证漏洞列表
func (r *ReportVerify) GetReportVerifyVulListCanCall(ctx context.Context, target string) ([]mysqls.Reportverifyvul, error) {
	var model mysqls.Reportverifyvul
	return model.GetReportVerifyVulListCanCall(ctx, target)
}

// ReportVerifyTask 进行报告验证任务
func (r *ReportVerify) ReportVerifyTask(ctx context.Context, objId, subObjId int, cveList []string, paramMap []map[string]string) ([]map[string]string, error) {
	return httpclients.DecisionReportVerifyTask(objId, subObjId, cveList, paramMap)
}

// ReportVerifyTaskResult 获取报告验证结果
func (r *ReportVerify) ReportVerifyTaskResult(ctx context.Context, objId, subObjId int) (bool, map[string][]string, error) {
	return httpclients.DecisionReportVerifyTaskResult(objId, subObjId)
}

// ReportVerifyHandleResult 进行验证任务结果处理
func (r *ReportVerify) ReportVerifyHandleResult(ctx context.Context, taskId int, target string, resultMap map[string][]string) error {
	var model mysqls.Reportverifyvul
	var err error
	for cve, result := range resultMap {
		if len(result) > 0 {
			err = model.UpdateReportverifyvulStatus(ctx, taskId, target, cve, enums.ReportVerifyStatusVerify)
		} else {
			err = model.UpdateReportverifyvulStatus(ctx, taskId, target, cve, enums.ReportVerifyStatusFailed)
		}
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

// ReportVerifyHandlePort 进行验证任务结果处理
func (r *ReportVerify) ReportVerifyHandlePort(ctx context.Context, taskId int, target string, portResultList []map[string]string) error {
	for _, portResult := range portResultList {
		var model = mysqls.Reportverifyport{
			Target:    target,
			Port:      portResult["port"],
			Scheme:    "tcp",
			Service:   portResult["service"],
			Component: portResult["fingerprint"],
			TaskId:    taskId,
		}
		err := model.AddReportverifyport(ctx)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

// GetWaitAndRunningTarget 获取待开始和运行中的目标
func (r *ReportVerify) GetWaitAndRunningTarget(ctx context.Context, taskId int) ([]mysqls.Reportverifytarget, error) {
	var model mysqls.Reportverifytarget
	status := []int{enums.TargetStatusToTrigger, enums.TargetStatusToBegin, enums.TargetStatusRunning}
	return model.GetTargetsByStatus(ctx, taskId, status)
}

// ReportVerifyOverView 报告验证漏洞统计功能
func (r *ReportVerify) ReportVerifyOverView(ctx context.Context, taskId int) (map[string]int, map[string]map[string]int, error) {
	var model mysqls.Reportverifyvul
	vulList, err := model.GetReportVerifyVulListByStatus(ctx, taskId)
	if err != nil {
		return nil, nil, err
	}
	taskStats := make(map[string]int)
	targetMap := make(map[string]map[string]int, 0)
	for _, vul := range vulList {
		taskStats["allVul"] += 1
		if targetMap[vul.Location] == nil {
			targetMap[vul.Location] = make(map[string]int, 0)
		}
		if vul.Status == enums.ReportVerifyStatusUnVerify {
			taskStats["unVerify"] += 1
			targetMap[vul.Location]["unVerify"] += 1
		} else if vul.Status == enums.ReportVerifyStatusVerify {
			taskStats["verify"] += 1
			targetMap[vul.Location]["verify"] += 1
		} else if vul.Status == enums.ReportVerifyStatusFailed {
			taskStats["failed"] += 1
			targetMap[vul.Location]["failed"] += 1
		} else if vul.Status == enums.ReportVerifyStatusExp {
			taskStats["exp"] += 1
			targetMap[vul.Location]["exp"] += 1
		}
		if vul.Status != enums.ReportVerifyStatusVerify && vul.Status != enums.ReportVerifyStatusExp {
			continue
		}
		if vul.Risk <= enums.VulLibrariesRiskHigh {
			taskStats["risk"] = enums.TaskRiskHigh
			taskStats["highVul"] += 1
			targetMap[vul.Location]["risk"] += enums.TargetRiskHigh
		} else if vul.Risk == enums.VulLibrariesRiskMiddle {
			taskStats["risk"] = enums.TaskRiskMid
			taskStats["middleVul"] += 1
			targetMap[vul.Location]["risk"] += enums.TargetRiskMid
		} else if vul.Risk == enums.VulLibrariesRiskLow {
			taskStats["risk"] = enums.TaskRiskLow
			taskStats["lowVul"] += 1
			targetMap[vul.Location]["risk"] += enums.TargetRiskLow
		} else {
			taskStats["risk"] = enums.TaskRiskSafe
			targetMap[vul.Location]["risk"] += enums.TargetRiskLowNoFound
		}
	}
	return taskStats, targetMap, nil
}

// UpdateTargetRisk 更新目标风险等级
func (r *ReportVerify) UpdateTargetRisk(ctx context.Context, taskId int, targetMap map[string]map[string]int) error {
	var model mysqls.Reportverifytarget
	for target, item := range targetMap {
		err := model.UpdateReportverifytargetRisk(ctx, taskId, target, item["risk"], item["unVerify"], item["verify"], item["failed"], item["exp"])
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

// UpdateTaskOverView 更新任务统计信息
func (r *ReportVerify) UpdateTaskOverView(ctx context.Context, taskId int, taskStats map[string]int) error {
	var (
		taskModel   mysqls.Reportverifytask
		targetModel mysqls.Reportverifytarget
	)
	err := taskModel.UpdateReportverifytaskRisk(ctx, taskId, taskStats["risk"], taskStats["unVerify"], taskStats["verify"], taskStats["failed"], taskStats["exp"])
	if err != nil {
		return err
	}

	allTarget, err := targetModel.GetAllTargets(ctx, taskId)
	for _, target := range allTarget {
		taskStats["allTarget"] += 1
		if target.IsAlive == 1 {
			taskStats["aliveTarget"] += 1
		}
		if target.Risk == enums.TargetRiskHigh {
			taskStats["highTarget"] += 1
		} else if target.Risk == enums.TargetRiskMid {
			taskStats["middleTarget"] += 1
		} else if target.Risk == enums.TargetRiskLow {
			taskStats["lowTarget"] += 1
		} else if target.Risk == enums.TargetRiskLowNoFound {
			taskStats["safeTarget"] += 1
		}
	}

	overviewByte, err := json.Marshal(taskStats)
	err = taskModel.UpdateReportverifytaskOverView(ctx, taskId, string(overviewByte))
	if err != nil {
		return err
	}
	return nil
}

// GetTaskOverView 获取任务统计信息
func (r *ReportVerify) GetTaskOverView(ctx context.Context, taskId int) (map[string]int, error) {
	var model = mysqls.Reportverifytask{
		ID: taskId,
	}
	task, err := model.GetReportverifytask(ctx)
	if err != nil {
		return nil, err
	}
	var overview map[string]int
	err = json.Unmarshal([]byte(task.Overview), &overview)
	if err != nil {
		return overview, err
	}
	return overview, nil
}

// TaskStop 结束任务
func (r *ReportVerify) TaskStop(ctx context.Context, taskId int) (map[string]int, error) {
	var model = mysqls.Reportverifytask{
		ID: taskId,
	}
	task, err := model.GetReportverifytask(ctx)
	if err != nil {
		return nil, err
	}
	var overview map[string]int
	err = json.Unmarshal([]byte(task.Overview), &overview)
	if err != nil {
		return overview, err
	}
	return overview, nil
}

// UpdateTargetStatusByTaskId 按任务id更新目标状态
func (r *ReportVerify) UpdateTargetStatusByTaskId(ctx context.Context, taskId, status int) error {
	var model mysqls.Reportverifytarget
	return model.UpdateReportverifytargetStatusByTaskId(ctx, taskId, status)
}

// TaskDelete 按任务id更新目标状态
func (r *ReportVerify) TaskDelete(ctx context.Context, taskId int) error {
	var model = mysqls.Reportverifytask{
		ID: taskId,
	}
	var targetModel mysqls.Reportverifytarget
	targetModel.DeleteReportverifytargetByTaskId(ctx, taskId)
	return model.DeleteReportverifytask(ctx)
}

// TargetDelete 删除目标
func (r *ReportVerify) TargetDelete(ctx context.Context, targetId int) error {
	var model = mysqls.Reportverifytarget{
		ID: targetId,
	}
	return model.DeleteReportverifytarget(ctx)
}

// GetRunningTargets 获取正在运行报告验证目标
func (r *ReportVerify) GetRunningTargets(ctx context.Context, taskId int) ([]mysqls.Reportverifytarget, error) {
	var model mysqls.Reportverifytarget
	status := []int{enums.TaskStatusRunning}
	return model.GetTargetsByStatus(ctx, taskId, status)
}

// VulDelete 删除漏洞
func (r *ReportVerify) VulDelete(ctx context.Context, vulId int) error {
	var model = mysqls.Reportverifyvul{
		ID: vulId,
	}
	return model.DeleteReportverifyvul(ctx)
}

// VulDetail 漏洞详情
func (r *ReportVerify) VulDetail(ctx context.Context, vulId int) (mysqls.Reportverifyvul, error) {
	var model = mysqls.Reportverifyvul{
		ID: vulId,
	}
	return model.GetReportverifyvul(ctx)
}

// AnalysisZip 分析zip压缩包文件
func (r *ReportVerify) AnalysisZip(ctx context.Context, tempFileName string, sourceMap map[string]string) error {
	reader, err := zip.OpenReader(tempFileName)
	if err != nil {
		return err
	}
	for _, z := range reader.File {
		if !strings.HasSuffix(z.Name, ".html") {
			continue
		}
		zf, err := z.Open()
		if err != nil {
			continue
		}
		fileByte, err := io.ReadAll(zf)
		if err != nil {
			continue
		}
		if len(fileByte) == 0 {
			continue
		}
		sourceMap[z.Name] = string(fileByte)
		zf.Close()
	}
	return nil
}

// HandleTimeoutTargets 处理超时目标
func (r *ReportVerify) HandleTimeoutTargets(ctx context.Context, runningTargets []mysqls.Reportverifytarget) {
	for _, target := range runningTargets {
		if int(time.Now().Sub(target.CreateTime).Hours()) > 1 {
			err := target.UpdateReportverifytargetStatus(ctx, target.ID, enums.TargetStatusFinish)
			if err != nil {
				log.Error("HandleTimeoutTargets: ", err)
			}
		}
	}
}
