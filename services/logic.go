package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/models/mysqls"
	"smart/models/redises"
	"smart/tools/enums"
	"smart/tools/network"
	"strconv"
	"strings"
	"time"
)

type Logic struct {
}

// TaskSave 逻辑漏洞测试
func (a *Logic) TaskSave(ctx context.Context, name, targetUrl, scanConfig string, scanType, userId int) (int, error) {
	var taskModel = mysqls.LogicTask{
		Name:       name,
		TargetUrl:  targetUrl,
		Status:     enums.LogicTaskStatusBegin,
		UserID:     userId,
		Type:       scanType,
		ScanConfig: scanConfig,
		Risk:       enums.LogicTaskRiskSafe,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	return taskModel.AddLogicTask(ctx)
}

// TargetSave 目标保存
func (a *Logic) TargetSave(ctx context.Context, targetUrl string, taskId, scanType int) (int, error) {
	var targetModel = mysqls.Logictarget{
		TargetURL: targetUrl,
		Status:    enums.LogicTaskStatusBegin,
		TaskID:    taskId,
		Type:      scanType,
		IsAlive:   enums.LogicTargetIsAliveN,
		//ConfigJson: scanConfig,
		Risk:       enums.LogicTaskRiskSafe,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	return targetModel.AddLogictarget(ctx)
}

// Stop 逻辑漏洞 - 任务结束
func (a *Logic) Stop(ctx context.Context, taskId int) error {
	var (
		taskModel   mysqls.LogicTask
		cacheClient redises.RedisHash
		targetModel mysqls.Logictarget
	)
	// 结束运行中目标的 yak进程
	runningTargetList, _ := targetModel.GetLogicTargetListByTargetIdAndStatus(ctx, []int{taskId}, enums.TargetStatusRunning)
	go func() {
		for _, target := range runningTargetList {
			callId, _ := cacheClient.GetHashHGet(ctx, enums.LogicCallIdCacheKey, strconv.Itoa(target.ID))
			if callId != "" {
				err := httpclients.DecisionScriptStop(callId)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()

	err := targetModel.UpdateTargetListByTaskId(ctx, taskId, enums.TargetStatusFinish)
	param := map[string]interface{}{
		"status": enums.LogicTaskStatusFinish,
	}
	err = taskModel.UpdateLogicTaskById(ctx, taskId, param)
	if err != nil {
		return err
	}
	return nil
}

// Delete 逻辑漏洞 - 任务删除
func (a *Logic) Delete(ctx context.Context, id int) error {
	var taskModel = mysqls.LogicTask{
		ID: id,
	}
	err := taskModel.DeleteLogicTask(ctx)
	if err != nil {
		return err
	}
	return nil
}

// List 逻辑漏洞 - 任务列表
func (a *Logic) List(ctx context.Context, page, size int, search string) ([]mysqls.LogicTask, int64, error) {
	var taskModel mysqls.LogicTask
	return taskModel.GetLogicTaskList(ctx, page, size, search)
}

// GetTaskIdsByStatus 逻辑漏洞 - 根据任务状态获取任务id
func (a *Logic) GetTaskIdsByStatus(ctx context.Context, status int) ([]int, error) {
	var taskModel mysqls.LogicTask
	taskIdList := make([]int, 0)
	taskList, err := taskModel.GetLogicTaskListByStatus(ctx, status)
	if err != nil {
		return taskIdList, err
	}
	for _, task := range taskList {
		taskIdList = append(taskIdList, task.ID)
	}
	return taskIdList, nil
}

// GetTargetsByStatus 逻辑漏洞 - 根据运行状态来获取目标列表
func (a *Logic) GetTargetsByStatus(ctx context.Context, status int) ([]mysqls.Logictarget, error) {
	var targetModel mysqls.Logictarget
	targetList, err := targetModel.GetLogicTargetListByStatus(ctx, status)
	if err != nil {
		return targetList, err
	}
	return targetList, nil
}

// BeyondPermCall 逻辑漏洞 - 越权漏洞调用
func (a *Logic) BeyondPermCall(ctx context.Context, url string, config BeyondPermConfig) (string, error) {
	vulParam := make([]map[string]interface{}, 0)
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "targetUrl",
		"value": url,
	})
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "rawCredential_type",
		"value": enums.Logic{}.GetCredPatternName(config.LoginCred.Pattern),
	})
	loginCredValue := removeTrailingNewlinesRegex(config.LoginCred.Value)
	if !strings.Contains(loginCredValue, `: `) {
		loginCredValue = strings.ReplaceAll(loginCredValue, `:`, `: `)
	}
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "rawCredential",
		"value": loginCredValue,
	})
	if config.WaitCred.Value != "" {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "fuzzCredential_type",
			"value": enums.Logic{}.GetCredPatternName(config.WaitCred.Pattern),
		})
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "fuzzCredential",
			"value": config.WaitCred.Value,
		})
	}
	if config.WhitePath != "" {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "path_list",
			"value": config.WhitePath,
		})
	}
	if config.BlackPath != "" {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "black_path_list",
			"value": config.BlackPath,
		})
	}
	if config.Keywords != "" {
		keywords := strings.ReplaceAll(config.Keywords, ",", "\n")
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "over_permission_flag_list",
			"value": keywords,
		})
	}
	if config.Crawler.Range == enums.CrawlerConfigAllDomain {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "scanRange",
			"value": enums.LogicCrawlerConfigAllDomain,
		})
	} else if config.Crawler.Range == enums.CrawlerConfigSubDomain {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "scanRange",
			"value": enums.LogicCrawlerConfigSubDomain,
		})
	}
	if config.Crawler.Depth != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "crawlerMaxDepth",
			"value": strconv.Itoa(config.Crawler.Depth),
		})
	}
	if config.Crawler.MaxLink != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "crawlerMaxUrl",
			"value": strconv.Itoa(config.Crawler.MaxLink),
		})
	}
	if config.Crawler.SingleLink != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "pageTimeout",
			"value": strconv.Itoa(config.Crawler.SingleLink),
		})
	}
	if len(config.Crawler.Sensitive) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "sensitiveWords",
			"value": strings.Join(config.Crawler.Sensitive, ","),
		})
	}
	if len(config.Crawler.WhiteWord) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "whitelist",
			"value": strings.Join(config.Crawler.WhiteWord, ","),
		})
	}
	if len(config.Crawler.BlackWord) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "blacklist",
			"value": strings.Join(config.Crawler.BlackWord, ","),
		})
	}
	callId, err := httpclients.DecisionScriptCall(enums.ScriptNameLogicBeyondPermission, vulParam)
	if err != nil {
		return callId, err
	}
	return callId, nil
}

// TraverseTestingCall 逻辑漏洞 - 信息遍历测试
func (a *Logic) TraverseTestingCall(ctx context.Context, url string, config SensInfoConfig) (string, error) {
	vulParam := make([]map[string]interface{}, 0)
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "targetUrl",
		"value": url,
	})
	if config.LoginCred.Pattern == enums.LogicCredPatternCookie {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "rawCookie",
			"value": config.LoginCred.Value,
		})
	} else {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "rawHeaders",
			"value": config.LoginCred.Value,
		})
	}
	if config.FuzzParam.Character != "" {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "string_param_list",
			"value": config.FuzzParam.Character,
		})
	}
	if config.FuzzParam.Number != "" {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "digital_param_list",
			"value": config.FuzzParam.Number,
		})
	}
	if config.FuzzDict.Character != "" {
		character := strings.ReplaceAll(config.FuzzDict.Character, ",", "\n")
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "string_dic_content",
			"value": character,
		})
	}
	if config.FuzzDict.Number != "" {
		number := strings.ReplaceAll(config.FuzzDict.Number, ",", "\n")
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "digital_dic_content",
			"value": number,
		})
	}
	if config.Response.JsonKeyword != "" {
		jsonKeyword := strings.ReplaceAll(config.Response.JsonKeyword, ",", "\n")
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "json_rsp_keyword_dic",
			"value": jsonKeyword,
		})
	}
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "not_json_check",
		"value": strconv.FormatBool(config.Response.NoJsonSwitch),
	})
	if config.Response.NoJsonKeyword != "" {
		keyword := strings.ReplaceAll(config.Response.NoJsonKeyword, ",", "\n")
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "not_json_rsp_keyword_dic",
			"value": keyword,
		})
	}
	if config.Crawler.Range == enums.CrawlerConfigAllDomain {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "scanRange",
			"value": enums.LogicCrawlerConfigAllDomain,
		})
	} else if config.Crawler.Range == enums.CrawlerConfigSubDomain {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "scanRange",
			"value": enums.LogicCrawlerConfigSubDomain,
		})
	}
	if config.Crawler.Depth != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "crawlerMaxDepth",
			"value": strconv.Itoa(config.Crawler.Depth),
		})
	}
	if config.Crawler.MaxLink != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "crawlerMaxUrl",
			"value": strconv.Itoa(config.Crawler.MaxLink),
		})
	}
	if config.Crawler.SingleLink != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "pageTimeout",
			"value": strconv.Itoa(config.Crawler.SingleLink),
		})
	}
	if len(config.Crawler.Sensitive) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "sensitiveWords",
			"value": strings.Join(config.Crawler.Sensitive, ","),
		})
	}
	if len(config.Crawler.WhiteWord) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "whitelist",
			"value": strings.Join(config.Crawler.WhiteWord, ","),
		})
	}
	if len(config.Crawler.BlackWord) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "blacklist",
			"value": strings.Join(config.Crawler.BlackWord, ","),
		})
	}
	callId, err := httpclients.DecisionScriptCall(enums.ScriptNameLogicTraverseTesting, vulParam)
	if err != nil {
		return callId, err
	}
	return callId, nil
}

// UnAuthAccessCall 逻辑漏洞 - 未授权访问调用
func (a *Logic) UnAuthAccessCall(ctx context.Context, url string, config UnAuthAccessConfig) (string, error) {
	vulParam := make([]map[string]interface{}, 0)
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "targetUrl",
		"value": url,
	})
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "rawCredential_type",
		"value": enums.Logic{}.GetCredPatternName(config.LoginCred.Pattern),
	})
	loginCredValue := removeTrailingNewlinesRegex(config.LoginCred.Value)
	if !strings.Contains(loginCredValue, `: `) {
		loginCredValue = strings.ReplaceAll(loginCredValue, `:`, `: `)
	}
	vulParam = append(vulParam, map[string]interface{}{
		"key":   "rawCredential",
		"value": loginCredValue,
	})
	if config.WhitePath != "" {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "path_list",
			"value": config.WhitePath,
		})
	}
	if config.BlackPath != "" {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "black_path_list",
			"value": config.BlackPath,
		})
	}
	if config.Keywords != "" {
		keywords := strings.ReplaceAll(config.Keywords, ",", "\n")
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "over_permission_flag_list",
			"value": keywords,
		})
	}
	if config.CredIdentifyList != "" {
		credIdentifyList := strings.ReplaceAll(config.CredIdentifyList, ",", "\n")
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "Credential_identification_list",
			"value": credIdentifyList,
		})
	}
	if config.Crawler.Range == enums.CrawlerConfigAllDomain {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "scanRange",
			"value": enums.LogicCrawlerConfigAllDomain,
		})
	} else if config.Crawler.Range == enums.CrawlerConfigSubDomain {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "scanRange",
			"value": enums.LogicCrawlerConfigSubDomain,
		})
	}
	if config.Crawler.Depth != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "crawlerMaxDepth",
			"value": strconv.Itoa(config.Crawler.Depth),
		})
	}
	if config.Crawler.MaxLink != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "crawlerMaxUrl",
			"value": strconv.Itoa(config.Crawler.MaxLink),
		})
	}
	if config.Crawler.SingleLink != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "pageTimeout",
			"value": strconv.Itoa(config.Crawler.SingleLink),
		})
	}
	if len(config.Crawler.Sensitive) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "sensitiveWords",
			"value": strings.Join(config.Crawler.Sensitive, ","),
		})
	}
	if len(config.Crawler.WhiteWord) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "whitelist",
			"value": strings.Join(config.Crawler.WhiteWord, ","),
		})
	}
	if len(config.Crawler.BlackWord) != 0 {
		vulParam = append(vulParam, map[string]interface{}{
			"key":   "blacklist",
			"value": strings.Join(config.Crawler.BlackWord, ","),
		})
	}
	callId, err := httpclients.DecisionScriptCall(enums.ScriptNameLogicUnAuthAccess, vulParam)
	if err != nil {
		return callId, err
	}
	return callId, nil
}

// HandleResult 逻辑漏洞 - 漏洞结果处理
func (a *Logic) HandleResult(ctx context.Context, callId string, taskId, targetId, logId, vulType int, targetUrl string) {
	err := a.LogInfoSave(ctx, taskId, targetId, logId, targetUrl, "", "开始进行逻辑漏洞任务测试")
	for {
		//每隔5秒获取下检测结果
		time.Sleep(5 * time.Second)
		end, callResult, err := httpclients.DecisionScriptCallResult(callId, true)

		if err != nil {
			return
		}

		//处理检测结果
		for _, item := range strings.Split(callResult, "\n") {
			var (
				requestData, responseData string
				tempJsonData              map[string]interface{}
			)
			err = json.Unmarshal([]byte(strings.TrimSpace(item)), &tempJsonData)
			if err != nil {
				log.Println(err)
				continue
			}
			// 处理pocname
			var pocname string
			if tempJsonData["pocname"] != nil {
				pocname = tempJsonData["pocname"].(string)
			} else {
				pocname = "crawlerx"
			}
			//处理扫描报文
			if tempJsonData["request"] != nil {
				requestData = tempJsonData["request"].(string)
				delete(tempJsonData, "request")
			}
			if tempJsonData["response"] != nil {
				responseData = tempJsonData["response"].(string)
				delete(tempJsonData, "response")
			}
			// 添加爬虫扫描路径信息结果
			if pocname == enums.ScriptNameCrawlerx {
				err = a.AddLogicFlowBaseByJsonData(ctx, taskId, targetId, tempJsonData)
				if err != nil {
					log.Error(err)
				}
			}
			tempDataByte, err := json.Marshal(tempJsonData)
			//对每条结果先保存到日志中
			err = a.LogInfoSave(ctx, taskId, targetId, logId, targetUrl, pocname, string(tempDataByte))
			if err != nil {
				log.Println(err)
			}
			if pocname == enums.ScriptNameCrawlerx {
				continue
			}
			// 从缓存中获取结果
			location := ""
			if tempJsonData["location"] != nil {
				location = tempJsonData["location"].(string)
			}
			var (
				cacheClient  redises.RedisHash
				commonClient redises.RedisCommon
			)
			url := network.ParseUrl2(location)
			urlKey := url.Scheme + url.Host + url.Path
			if tempJsonData["logic_unauthorized_access_fuzz_method"] != nil {
				urlKey += tempJsonData["logic_unauthorized_access_fuzz_method"].(string)
			}
			fmt.Println("去重key", urlKey)
			cacheResult, err := cacheClient.GetHashHGet(ctx, enums.LogicResultCacheKeyPreKey+strconv.Itoa(targetId), urlKey)
			if cacheResult != "" {
				var tempValue map[string]interface{}
				err = json.Unmarshal([]byte(cacheResult), &tempValue)
				if tempJsonData["payload"] != nil {
					tempValue["payload"] = strings.Trim(tempValue["payload"].(string), `"`) + "," + strings.Trim(tempJsonData["payload"].(string), `""`)
				}
				tempValueByte, _ := json.Marshal(tempValue)
				err = cacheClient.SetHashHset(ctx, enums.LogicResultCacheKeyPreKey+strconv.Itoa(targetId), urlKey, string(tempValueByte))
				err = commonClient.Expire(ctx, enums.LogicResultCacheKeyPreKey+strconv.Itoa(targetId), 10*time.Minute)
				param := map[string]interface{}{"vul_result": string(tempValueByte)}
				var logicVul mysqls.Logicvul
				err = logicVul.UpdateLogicvulParam(ctx, tempValue["vulId"].(string), param)
				continue
			}
			// 保存未授权漏洞
			libRisk, vulId := 0, 0
			if tempJsonData["vul_name"] != nil && tempJsonData["vul_name"] == "未授权" {
				vulId, err = a.VulSave(ctx, taskId, targetId, enums.LogicTypeUnAuthAccess, enums.LogicVulRiskHigh, enums.ReportVerifyStatusUnVerify, enums.UnAuthVulId, pocname, enums.UnAuthVulName, enums.UnAuthVulDesc, enums.UnAuthVulFixSuggest, location, targetUrl, string(tempDataByte), requestData, responseData)
				if err != nil {
					log.Println(err)
				}
				libRisk = enums.LogicVulRiskHigh
			} else {
				// 保存到漏洞结果表中
				decisionDetailReq := httpclients.OpenVulScriptDetailByPocnameReq{Pocname: pocname}
				vulInfo, err := httpclients.GetDecisionScriptDetailByPocname(decisionDetailReq)
				if err != nil {
					log.Println(err)
				}
				lib := vulInfo.Data.Libraries
				vulId, err = a.VulSave(ctx, taskId, targetId, vulType, lib.Risk, enums.ReportVerifyStatusUnVerify, lib.VulId, pocname, lib.Name, lib.Description, lib.FixSuggest, location, targetUrl, string(tempDataByte), requestData, responseData)
				libRisk = lib.Risk
			}
			if libRisk > 0 && libRisk <= enums.LogicVulRiskLow {
				risk := libRisk
				if risk > 1 {
					risk = libRisk - 1
				}
				err = a.UpdateTargetParam(ctx, targetId, map[string]interface{}{"risk": risk, "update_time": time.Now()})
				err = a.UpdateTaskParam(ctx, taskId, map[string]interface{}{"risk": risk, "high_num": 1, "update_time": time.Now()})
			}
			tempJsonData["vulId"] = strconv.Itoa(vulId)
			tempDataByte, _ = json.Marshal(tempJsonData)
			err = cacheClient.SetHashHset(ctx, enums.LogicResultCacheKeyPreKey+strconv.Itoa(targetId), urlKey, string(tempDataByte))
		}
		if end {
			break
		}
	}
	err = a.LogInfoSave(ctx, taskId, targetId, logId, targetUrl, "", "任务测试结束")
	//更新日志状态为已完成
	err = a.UpdateLogStatus(ctx, targetId, map[string]interface{}{"status": enums.TargetStatusFinish, "end_time": time.Now()})
	if err != nil {
		log.Error("BeyondPermHandleResult update status error:", err.Error())
	}
	//更新目标状态为已完成
	err = a.UpdateTargetParam(ctx, targetId, map[string]interface{}{"status": enums.TargetStatusFinish})
	if err != nil {
		log.Error("LogicTaskExec UpdateTargetStatus error:", err.Error())
	}
	// 判断任务结束条件
	runningTarget, err := a.GetTargetsByTaskIdAndStatus(ctx, taskId, []int{enums.LogicTaskStatusTrigger, enums.LogicTaskStatusBegin, enums.LogicTaskStatusRunning})
	if len(runningTarget) == 0 {
		//更新任务状态为已完成
		err = a.UpdateTaskParam(ctx, taskId, map[string]interface{}{"status": enums.LogicTaskStatusFinish, "update_time": time.Now()})
		if err != nil {
			log.Error("LogicTaskExec UpdateTargetStatus error:", err.Error())
		}
	}
	// 清理缓存信息
	var cacheClient redises.RedisHash
	err = cacheClient.HDel(ctx, enums.LogicCallIdCacheKey, callId)
}

// UpdateTaskParam 逻辑漏洞 - 更新任务参数
func (a *Logic) UpdateTaskParam(ctx context.Context, taskId int, param map[string]interface{}) error {
	var taskModel mysqls.LogicTask
	err := taskModel.UpdateLogicTaskById(ctx, taskId, param)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTargetParam 逻辑漏洞 - 更新目标参数
func (a *Logic) UpdateTargetParam(ctx context.Context, targetId int, param map[string]interface{}) error {
	var targetModel mysqls.Logictarget
	err := targetModel.UpdateLogicTargetParam(ctx, targetId, param)
	if err != nil {
		return err
	}
	return nil
}

// UpdateLogStatus 逻辑漏洞 - 更新日志状态
func (a *Logic) UpdateLogStatus(ctx context.Context, targetId int, param map[string]interface{}) error {
	var logModel mysqls.Logiclog
	//param := map[string]interface{}{"status": status, "end_time": time.Now()}
	err := logModel.UpdateLogicLogParam(ctx, targetId, param)
	if err != nil {
		return err
	}
	return nil
}

// TargetList 逻辑漏洞 - 任务列表
func (a *Logic) TargetList(ctx context.Context, taskId, page, size int, search string) ([]mysqls.Logictarget, int64, error) {
	var targetModel mysqls.Logictarget
	return targetModel.GetLogictargetList(ctx, taskId, page, size, search)
}

// VulList 逻辑漏洞 - 漏洞列表
func (a *Logic) VulList(ctx context.Context, taskId, page, size int, search string) ([]mysqls.Logicvul, int64, error) {
	var vulModel mysqls.Logicvul
	return vulModel.GetLogicvulList(ctx, taskId, page, size, search)
}

// VulList 逻辑漏洞 - 漏洞列表
func (a *Logic) LogList(ctx context.Context, taskId, page, size int, search string) ([]mysqls.Logiclog, int64, error) {
	var logModel mysqls.Logiclog
	return logModel.GetLogiclogList(ctx, taskId, page, size, search)
}

// GetVulById 逻辑漏洞 - 通过id获取漏洞详情
func (a *Logic) GetVulById(ctx context.Context, id int) (mysqls.Logicvul, error) {
	var vulModel = mysqls.Logicvul{
		ID: id,
	}
	return vulModel.GetLogicvul(ctx)
}

// GetLogById 逻辑漏洞 - 通过id获取日志详情
func (a *Logic) GetLogById(ctx context.Context, id int) ([]mysqls.Logicloginfo, error) {
	var logModel mysqls.Logicloginfo
	return logModel.GetLogicloginfoListByLogId(ctx, id)
}

// 枚举信息 - 扫描类型
func (a *Logic) ScanTypeEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.Logic{}.AllTypeEnum())
}

// 枚举信息 - 认证凭证
func (a *Logic) CredPatternEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.Logic{}.AllCredPatternEnum())
}

// 枚举信息 - 认证凭证
func (a *Logic) FuzzParamEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.Logic{}.AllCredPatternEnum())
}

// VulDelete 逻辑漏洞 - 漏洞删除
func (a *Logic) VulDelete(ctx context.Context, id int) error {
	var vulModel = mysqls.Logicvul{
		ID: id,
	}
	err := vulModel.DeleteLogicvul(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetTaskById 逻辑漏洞 - 根据id获取任务
func (a *Logic) GetTaskById(ctx context.Context, id int) (mysqls.LogicTask, error) {
	var taskModel = mysqls.LogicTask{
		ID: id,
	}
	return taskModel.GetLogicTask(ctx)
}

// LogSave 日志保存
func (a *Logic) LogSave(ctx context.Context, targetUrl string, taskId, targetId int) (int, error) {
	var logModel = mysqls.Logiclog{
		TaskID:     taskId,
		TargetID:   targetId,
		TargetURL:  targetUrl,
		Status:     enums.LogicTaskStatusBegin,
		CreateTime: time.Now(),
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		IsAlive:    enums.LogicTargetIsAliveN,
	}
	return logModel.AddLogiclog(ctx)
}

// LogInfoSave 日志详情保存
func (a *Logic) LogInfoSave(ctx context.Context, taskId, targetId, logId int, targetUrl, pocname, result string) error {
	var logInfoModel = mysqls.Logicloginfo{
		TaskID:     taskId,
		TargetID:   targetId,
		LogID:      logId,
		TargetURL:  targetUrl,
		Pocname:    pocname,
		Result:     result,
		CreateTime: time.Now(),
	}
	return logInfoModel.AddLogicloginfo(ctx)
}

// GetTargetsByTaskIdAndStatus 逻辑漏洞 - 根据任务id和状态来获取目标
func (a *Logic) GetTargetsByTaskIdAndStatus(ctx context.Context, taskId int, status []int) ([]mysqls.Logictarget, error) {
	var targetModel mysqls.Logictarget
	targetList, err := targetModel.GetLogicTargetListByTaskIdAndStatus(ctx, taskId, status)
	if err != nil {
		return targetList, err
	}
	return targetList, nil
}

// VulSave 逻辑漏洞 - 漏洞保存
func (a *Logic) VulSave(ctx context.Context, taskId, targetId, vulType, risk, status int, vulId, pocname, name, description, fixSuggest, location, targetUrl, result, request, response string) (int, error) {
	VerMsg := make([]map[string]string, 0)
	requestByte, _ := base64.StdEncoding.DecodeString(request)
	responseByte, _ := base64.StdEncoding.DecodeString(response)
	tempMsg := map[string]string{
		"request":  string(requestByte),
		"response": string(responseByte),
	}
	VerMsg = append(VerMsg, tempMsg)
	VerMsgByte, _ := json.Marshal(VerMsg)
	var vulModel = mysqls.Logicvul{
		TaskID:        taskId,
		TargetID:      targetId,
		Pocname:       pocname,
		Name:          name,
		Class:         "逻辑漏洞",
		Type:          vulType,
		Risk:          risk,
		Location:      location,
		Description:   description,
		FixSuggest:    fixSuggest,
		VulParam:      "",
		VulResult:     result,
		VerMsg:        string(VerMsgByte),
		DecisionVulID: vulId,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		TargetUrl:     targetUrl,
		Status:        status,
		VulID:         vulId,
	}
	return vulModel.AddLogicvul(ctx)
}

// ReportSave 报告保存
func (a *Logic) ReportSave(ctx context.Context, name string, retype int, configjson string, format int, userId int) error {
	var reportMysqls = mysqls.Reportrecord{
		Name:       name,
		Type:       retype,
		Status:     enums.ReportStatusWait,
		ConfigJSON: configjson,
		Format:     format,
		Content:    "",
		UserID:     userId,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	_, err := reportMysqls.AddReportrecord(ctx)
	return err
}

// GetTargetByTaskId 根据任务id获取目标信息
func (a *Logic) GetTargetByTaskId(ctx context.Context, taskId int) ([]mysqls.Logictarget, error) {
	var targetModel mysqls.Logictarget
	statusEnum := []int{enums.LogicTaskStatusBegin, enums.LogicTaskStatusRunning, enums.LogicTaskStatusFinish, enums.LogicTaskStatusPausing}
	list, err := targetModel.GetLogicTargetListByTaskIdAndStatus(ctx, taskId, statusEnum)
	return list, err
}

// GetVulByTaskId 逻辑漏洞 - 通过任务id获取漏洞
func (a *Logic) GetVulByTaskId(ctx context.Context, taskId int) ([]mysqls.Logicvul, error) {
	var vulModel mysqls.Logicvul
	return vulModel.GetLogicVulByTaskId(ctx, taskId)
}

// GetTargetById 根据目标id获取目标信息
func (a *Logic) GetTargetById(ctx context.Context, targetId int) (mysqls.Logictarget, error) {
	var targetModel = mysqls.Logictarget{
		ID: targetId,
	}
	return targetModel.GetLogictarget(ctx)
}

// GetVulByTargetId 逻辑漏洞 - 通过目标id获取漏洞
func (a *Logic) GetVulByTargetId(ctx context.Context, targetId int) ([]mysqls.Logicvul, error) {
	var vulModel mysqls.Logicvul
	return vulModel.GetLogicVulByTargetId(ctx, targetId)
}

// GetTargetStats 计算目标的统计信息
func (a *Logic) GetTargetStats(ctx context.Context, targetId int) (int, int, [6]int, error) {
	var (
		logicVul    mysqls.Logicvul
		vulNum      int                          //漏洞总数
		risklevel   = enums.TargetRiskLowNoFound //目标等级，默认为4->未发现
		vulNumArray [6]int                       //每个等级的数量，元素含义分别为：无漏洞个数/致命漏洞个数/高危漏洞个数/中危漏洞个数/低危漏洞个数/信息漏洞个数
	)
	vulRes, err := logicVul.GetLogicVulByTargetId(ctx, targetId)
	vulNum = len(vulRes) //漏洞总数
	if err != nil {
		return vulNum, risklevel, vulNumArray, err
	}
	for i := 0; i < len(vulRes); i++ {
		if vulRes[i].Status == enums.VulStatusRepairSuccess || vulRes[i].Risk == enums.VulLibrariesRiskNot { //已经修复或risk为0的算安全
			vulNumArray[5] += 1
			continue
		}
		vulNumArray[vulRes[i].Risk] += 1
	}
	//计算目标的漏洞等级
	if vulNumArray[1] > 0 || vulNumArray[2] > 0 { //致命或高->高
		risklevel = 1
	} else if vulNumArray[3] > 0 { //中->中
		risklevel = 2
	} else if vulNumArray[4] > 0 { //低->低
		risklevel = 3
	}
	return vulNum, risklevel, vulNumArray, nil
}

// GetTargetStatsBytargetIds 计算目标的统计信息
func (a *Logic) GetTargetStatsBytargetIds(ctx context.Context, targetIds []int) (map[int]int, map[int]int, map[int][6]int, error) {
	var (
		taskVul     mysqls.Logicvul
		vulNum      = make(map[int]int, 0)    //漏洞总数
		risklevel   = make(map[int]int, 0)    //目标等级，默认为4->未发现
		vulNumArray = make(map[int][6]int, 0) //每个等级的数量，元素含义分别为：无漏洞个数/致命漏洞个数/高危漏洞个数/中危漏洞个数/低危漏洞个数/信息漏洞个数
	)
	//初始化
	for i := 0; i < len(targetIds); i++ {
		vulNum[targetIds[i]] = 0
		risklevel[targetIds[i]] = enums.TargetRiskLowNoFound
		vulNumArray[targetIds[i]] = [6]int{}
	}

	vulRes, err := taskVul.GetLogicVulByTargetIds(ctx, targetIds)
	if err != nil {
		return vulNum, risklevel, vulNumArray, err
	}
	for i := 0; i < len(vulRes); i++ {
		vulNum[vulRes[i].TargetID] += 1
		tmp := vulNumArray[vulRes[i].TargetID]
		if vulRes[i].Status == enums.VulStatusRepairSuccess || vulRes[i].Risk == enums.VulLibrariesRiskNot { //已经修复或risk为0的算安全
			tmp[5] += 1
		} else {
			tmp[vulRes[i].Risk] += 1
		}
		vulNumArray[vulRes[i].TargetID] = tmp
	}
	for k, v := range vulNumArray {
		if v[1] > 0 || v[2] > 0 { //致命或高->高
			risklevel[k] = 1
		} else if v[3] > 0 { //中->中
			risklevel[k] = 2
		} else if v[4] > 0 { //低->低
			risklevel[k] = 3
		}
	}
	return vulNum, risklevel, vulNumArray, nil
}

func removeTrailingNewlinesRegex(s string) string {
	re := regexp.MustCompile(`\n+$`)
	return re.ReplaceAllString(s, "")
}

// FlowBaseList 被动流量列表
func (a *Logic) FlowBaseList(ctx context.Context, search string, flowTaskId int, page, size int) ([]mysqls.LogicFlowBase, int64, error) {
	var lfModel mysqls.LogicFlowBase
	return lfModel.GetFlowbaseList(ctx, search, flowTaskId, page, size)
}

// AddLogicFlowBaseByJsonData 添加被动流量
func (a *Logic) AddLogicFlowBaseByJsonData(ctx context.Context, taskId, targetId int, tempJsonData map[string]interface{}) error {
	var method, respCode, url, title, tags, ip, respType, request, response string
	var err error
	if tempJsonData["method"] != nil {
		method = tempJsonData["method"].(string)
	}
	if tempJsonData["statusCode"] != nil {
		respCode = strconv.FormatFloat(tempJsonData["statusCode"].(float64), 'f', -1, 64)
	}
	if tempJsonData["location"] != nil {
		url = tempJsonData["location"].(string)
		ip, err = extractIPFromURL(tempJsonData["location"].(string))
	}
	if tempJsonData["title"] != nil {
		title = tempJsonData["title"].(string)
	}
	if tempJsonData["tags"] != nil {
		tags = tempJsonData["tags"].(string)
	}
	if tempJsonData["requestBody"] != nil {
		request = tempJsonData["requestBody"].(string)
		delete(tempJsonData, "requestBody")
	}
	if tempJsonData["responseHeaders"] != nil && tempJsonData["responseBody"] != nil && tempJsonData["statusCode"] != nil {
		headers := tempJsonData["responseHeaders"].(map[string]interface{})
		body := tempJsonData["responseBody"].(string)
		statusCode := int(tempJsonData["statusCode"].(float64))

		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode)))
		for key, value := range headers {
			buffer.WriteString(key)
			buffer.WriteString(": ")
			buffer.WriteString(value.(string))
			buffer.WriteString("\r\n")
		}
		buffer.WriteString("\r\n")
		buffer.WriteString(body)
		response = buffer.String()
		delete(tempJsonData, "responseHeaders")
		delete(tempJsonData, "responseBody")
	}
	if err != nil {
		log.Println(err)
	}
	return a.AddLogicFlowBase(ctx, taskId, targetId, method, respCode, url, title, tags, ip, respType, request, response)
}

// AddLogicFlowBase 添加被动流量
func (a *Logic) AddLogicFlowBase(ctx context.Context, taskId, targetId int, method, respCode, url, title, tags, ip, respType, request, response string) error {
	var lfModel = mysqls.LogicFlowBase{
		TaskID:          taskId,
		TargetID:        targetId,
		Method:          method,
		RespCode:        respCode,
		Url:             url,
		RespTitle:       title,
		Tags:            tags,
		IP:              ip,
		RespContentType: respType,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		Protocol:        1,
		ReqHeader:       request,
		RespHeader:      response,
	}
	return lfModel.AddFlowbase(ctx)
}

// extractIPFromURL 从url中提取ip
func extractIPFromURL(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	// 提取主机名
	host := parsedURL.Hostname()

	// 尝试解析为主机地址
	ip := net.ParseIP(host)
	if ip == nil {
		return "", fmt.Errorf("hostname %s is not a valid IP address", host)
	}

	return ip.String(), nil
}

// extractContentType 提取响应类型
func extractContentType(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	// 提取主机名
	host := parsedURL.Hostname()

	// 尝试解析为主机地址
	ip := net.ParseIP(host)
	if ip == nil {
		return "", fmt.Errorf("hostname %s is not a valid IP address", host)
	}

	return ip.String(), nil
}

// LogicFlowBaseInfo 逻辑漏洞流量详情
func (a *Logic) LogicFlowBaseInfo(ctx context.Context, flowBaseId int) (mysqls.LogicFlowBase, error) {
	var flowBaseMysql mysqls.LogicFlowBase
	return flowBaseMysql.GetFlowbase(ctx, flowBaseId)
}
