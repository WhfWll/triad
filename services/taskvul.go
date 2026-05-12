package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"smart/models/mysqls"
	aesEncryption "smart/tools/encryption"
	"smart/tools/enums"
	"smart/tools/invoke"
	"smart/tools/network"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/encryption"
)

// extractLocationFromRequest 从HTTP请求报文中提取完整的URL
func extractLocationFromRequest(scriptResultMap map[string]string, riskData invoke.Risk) string {
	var location string
	if scriptResultMap["location"] != "" {
		location = scriptResultMap["location"]
	} else {
		location = riskData.Url
	}
	if location != "" {
		return location
	}

	if len(riskData.Request) == 0 {
		return ""
	}

	requestStr := string(riskData.Request)
	lines := strings.Split(strings.ReplaceAll(requestStr, "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return ""
	}

	// 解析第一行获取请求路径
	firstLine := strings.Fields(lines[0])
	if len(firstLine) < 2 {
		return ""
	}

	path := firstLine[1]

	// 查找Host头
	var host string
	for _, line := range lines[1:] {
		if strings.HasPrefix(strings.ToLower(line), "host:") {
			host = strings.TrimSpace(line[5:])
			break
		}
	}

	if host == "" {
		return ""
	}

	// 判断协议，默认使用http
	scheme := "http"
	if strings.Contains(strings.ToLower(requestStr), "https") {
		scheme = "https"
	}

	return scheme + "://" + host + path
}

// 检测出的漏洞管理的服务
type TaskVul struct {
}

// AddTaskVul 插入一条检测漏洞
func (t *TaskVul) AddTaskVul(ctx context.Context, taskId int, targetId int, scriptResult ScriptResult, targetResultID int) error {
	var taskVul mysqls.TaskVul
	if scriptResult.Script.DataType == enums.VulDataTypeTwo { //查询数据是否已经存在
		vulRes, err := taskVul.GetTaskVulByPocNameAndTargetId(ctx, targetId, scriptResult.Script.Pocname)
		if err != nil {
			return err
		}
		if vulRes.ID != 0 {
			return nil
		}
	}
	// 销售许可证临时添加

	taskVul.DataType = scriptResult.Script.DataType
	//taskVul.Name = scriptResult.Libraries.Name
	// taskVul.Name = hex.EncodeToString(aesEcb.AesEncryptECB([]byte(scriptResult.Libraries.Name), key))
	taskVul.Name = scriptResult.Libraries.Name

	taskVul.Class = scriptResult.Libraries.Class
	taskVul.Type = scriptResult.Libraries.Type
	taskVul.Risk = scriptResult.Libraries.Risk

	//taskVul.ExploitImpact = scriptResult.Libraries.ExploitImpact
	taskVul.VulID = scriptResult.Libraries.Cve + "," + scriptResult.Libraries.Cnvd + "," + scriptResult.Libraries.Cnnvd
	taskVul.Description = scriptResult.Libraries.Description
	//taskVul.Description = hex.EncodeToString(aesEcb.AesEncryptECB([]byte(scriptResult.Libraries.Description), key))

	taskVul.FixSuggest = scriptResult.Libraries.FixSuggest
	//taskVul.FixSuggest = hex.EncodeToString(aesEcb.AesEncryptECB([]byte(scriptResult.Libraries.FixSuggest), key))

	//taskVul.Location = scriptResult.Result.Location
	taskVul.Location = scriptResult.Result.Location

	//taskVul.Pocname = scriptResult.Script.Pocname
	taskVul.Pocname = scriptResult.Script.Pocname

	taskVul.TargetResultID = targetResultID
	taskVul.VulParam = scriptResult.Script.ScriptParam
	taskVul.TaskID = taskId
	taskVul.TargetID = targetId
	//taskVul.Location = scriptResult.Result.Location

	taskVul.PublishedTime = scriptResult.Libraries.PublishedTime
	taskVul.Status = enums.VulStatusVerifyExist
	if scriptResult.Script.VerifyType == enums.VulScriptVerifyTypeExp { //exp类型漏洞利用成功
		taskVul.Status = enums.VulStatusVerifyExist
	}
	if scriptResult.Script.DataType == enums.VulDataTypeTwo { //待测漏洞是位置漏洞
		taskVul.Status = enums.VulStatusVerifyExist
	}
	taskVul.TestStatus = enums.VulTestStatusNotTest
	taskVul.CreateTime = time.Now()
	taskVul.UpdateTime = time.Now()
	taskVul.TargetUrl = scriptResult.Result.TargetUrl
	//taskVul.TargetUrl = scriptResult.Result.TargetUrl
	taskVul.DecisionVulId = scriptResult.Libraries.VulId
	//处理漏洞截图、请求报文、响应报文
	var (
		scriptDetailTemp   map[string]string
		verMsgTmp          = make([]ScriptResultDetailVulProve, 0)
		snapshottmp        string
		payload            string
		payloadSuccessFlag string
	)
	json.Unmarshal([]byte(scriptResult.Result.Detail), &scriptDetailTemp)
	if v, ok := scriptDetailTemp["screen_data"]; ok { //截图
		snapshottmp = v
		delete(scriptDetailTemp, "screen_data") //剔除结果中的截图数据
	}
	taskVul.Snapshot = snapshottmp
	if v, ok := scriptDetailTemp["payload"]; ok { //payload
		payload, _ = url.QueryUnescape(v)
		scriptDetailTemp["payload"], _ = url.QueryUnescape(v)
	}
	if v, ok := scriptDetailTemp["payload_success_flag"]; ok { //payload_success_flag
		payloadSuccessFlag = v
	}

	if v, ok := scriptDetailTemp["vul_prove"]; ok { //多个报文
		json.Unmarshal([]byte(v), &verMsgTmp)
	}
	if len(verMsgTmp) == 0 {
		verMsgTmp = append(verMsgTmp, ScriptResultDetailVulProve{
			Request:            scriptResult.Result.Request,
			Response:           scriptResult.Result.Response,
			Payload:            payload,
			PayloadSuccessFlag: payloadSuccessFlag,
		})
	}

	verMsgTmpByte, _ := json.Marshal(verMsgTmp)
	taskVul.VerMsg = string(verMsgTmpByte)

	tmpDetail, _ := json.Marshal(scriptDetailTemp)
	taskVul.VulResult = string(tmpDetail)
	err := taskVul.AddTaskVul(ctx)

	if err != nil {
		log.Errorf("[DEBUG-VUL] AddTaskVulByScannerResult 失败: %v", err)
		return errors.New("AddTaskVul error: " + err.Error())
	}
	log.Infof("[DEBUG-VUL] AddTaskVulByScannerResult 成功: ID=%d", taskVul.ID)
	return err
}

// UpdateTaskVulById 修改一条检测漏洞
func (t *TaskVul) UpdateTaskVulById(ctx context.Context, taskVulId, datatype int, status, testStatus int) error {
	var taskVul mysqls.TaskVul
	var param = map[string]interface{}{
		"status":      status,
		"test_status": testStatus,
		"update_time": time.Now(),
	}
	return taskVul.UpdateByIdAndDataType(ctx, taskVulId, datatype, param)
}

// EndTaskVulById 结束待检测漏洞
func (t *TaskVul) EndTaskVulById(ctx context.Context, taskVulId, datatype int, status, testStatus int) error {
	var taskVul mysqls.TaskVul
	vulRes, err := taskVul.GetById(ctx, taskVulId)
	if err != nil {
		return err
	}
	if vulRes.Status == enums.VulStatusVerifyExist || vulRes.TestStatus == enums.VulTestStatusTested || vulRes.TestStatus == enums.VulTestStatusNotTest {
		return nil
	}
	var param = map[string]interface{}{
		"status":      status,
		"test_status": testStatus,
		"update_time": time.Now(),
	}
	return taskVul.UpdateByIdAndDataType(ctx, taskVulId, datatype, param)
}

// VulList 漏洞测试列表及筛选
func (t *TaskVul) VulList(ctx context.Context, taskId, targetId int, vultype, risk int, search string, dataType, page, limit, status int) ([]mysqls.TaskVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.GetTaskVulList(ctx, taskId, targetId, vultype, risk, search, dataType, page, limit, status)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		t.DecryptTaskVul(&list[i])
	}
	return list, count, nil
}

// DecryptTaskVul 解密TaskVul中的加密字段（兼容旧数据）
func (t *TaskVul) DecryptTaskVul(vul *mysqls.TaskVul) {
	aesEcb := aesEncryption.AesEcb{}
	// Helper to decrypt a single field
	decryptField := func(val *string) {
		if val != nil && utils.IsHexString(*val) {
			decoded, err := hex.DecodeString(*val)
			if err == nil {
				// aesKey is defined in assetgroup.go (package private)
				decrypted := aesEcb.AesDecryptECB(decoded, aesKey)
				// If decryption is successful (no panic, though AesDecryptECB usually doesn't return error), use it.
				*val = string(decrypted)
			}
		}
	}

	decryptField(&vul.Name)
	decryptField(&vul.Pocname)
	decryptField(&vul.TargetUrl)
	decryptField(&vul.Location)
	decryptField(&vul.Description)
	decryptField(&vul.FixSuggest)
	decryptField(&vul.VulResult)
	decryptField(&vul.VerMsg)
}

// SecondScanVulList 二次扫描漏洞列表
func (t *TaskVul) SecondScanVulList(ctx context.Context, taskId []int, targetId int, vultype, risk int, search string, dataType, page, limit, status int) ([]mysqls.TaskVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.SecondScanGetTaskVulList(ctx, taskId, targetId, vultype, risk, search, dataType, page, limit, status)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		t.DecryptTaskVul(&list[i])
	}
	return list, count, nil
}

// VulInfo 漏洞测试详情
func (t *TaskVul) VulInfo(ctx context.Context, taskVulId int) (mysqls.TaskVul, error) {
	var taskVul mysqls.TaskVul
	return taskVul.GetTaskVul(ctx, taskVulId)
}

// GetVulTargetIdByIds 根据id获取目标id
func (t *TaskVul) GetVulTargetIdByIds(ctx context.Context, ids any) ([]int, error) {
	var (
		taskVul mysqls.TaskVul
		result  = make([]int, 0)
	)
	vulRes, err := taskVul.GetTaskVulListByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(vulRes); i++ {
		result = append(result, vulRes[i].TargetID)
	}
	return result, nil
}

// GetVulByIds 根据id获取漏洞数据
func (t *TaskVul) GetVulByIds(ctx context.Context, ids any) ([]mysqls.TaskVul, error) {
	var taskVul mysqls.TaskVul
	vulRes, err := taskVul.GetTaskVulListByIds(ctx, ids)
	return vulRes, err
}

// GetVulMapByIDs 根据id获取漏洞数据，返回 map[int]TaskVul
func (t *TaskVul) GetVulMapByIDs(ctx context.Context, ids any) (map[int]mysqls.TaskVul, error) {
	var taskVulModel mysqls.TaskVul
	vulList, err := taskVulModel.GetTaskVulListByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	vulMap := make(map[int]mysqls.TaskVul, len(vulList))
	for _, v := range vulList {
		vulMap[v.ID] = v
	}
	return vulMap, nil
}

// GetVulByIds 根据id获取漏洞数据
func (t *TaskVul) TestVulCreatDecisionParams(vulRes []mysqls.TaskVul, targetlogmap map[int]int) ([]TaskControlMessage, error) {
	var result = make([]TaskControlMessage, 0)
	//根据目标查询日志id
	for i := 0; i < len(vulRes); i++ {
		var msgStruct = TaskControlMessage{
			ObjId:           strconv.Itoa(vulRes[i].TargetID),
			Operate:         enums.TaskOperateTestVul,
			CallTools:       make([]string, 0),
			TargetUrl:       vulRes[i].TargetUrl,
			CallVuls:        make([]int, 0),
			InfoList:        make([]Info, 0),
			Info:            "",
			TestIntensity:   "",
			SafeTestPocname: vulRes[i].Pocname,
			SafeTestId:      vulRes[i].ID,
		}
		//info
		var infomsg = TaskControlMessageInfo{
			TaskId:   vulRes[i].TaskID,
			TargetId: vulRes[i].TargetID,
		}
		if v, ok := targetlogmap[vulRes[i].TargetID]; ok {
			infomsg.LogId = v
		}
		infostring, err := json.Marshal(infomsg)
		if err != nil {
			continue
		}
		msgStruct.Info = string(infostring)
		result = append(result, msgStruct)
	}
	return result, nil
}

// TestVulUpdateStatus 修改测试漏洞的状态
func (t *TaskVul) TestVulUpdateStatus(ctx context.Context, vulRes []mysqls.TaskVul, testStatus int) error {
	var (
		ids     = make([]int, 0)
		taskVul mysqls.TaskVul
	)
	for i := 0; i < len(vulRes); i++ {
		ids = append(ids, vulRes[i].ID)
	}
	err := taskVul.UpdateTestStatusByIds(ctx, ids, testStatus)
	if err != nil {
		return err
	}
	return nil
}

// VulDel 删除漏洞测试
func (t *TaskVul) VulDel(ctx context.Context, ids any) error {
	var taskVul mysqls.TaskVul
	return taskVul.DeleteTaskVul(ctx, ids)
}

// Base64StdDecodeString //base64解密
func (t *TaskVul) Base64StdDecodeString(msg string) (string, error) {
	resMsg, err := encryption.Base64StdDecodeString(msg)
	return string(resMsg), err
}

func (t *TaskVul) Base64StdEncodeToString(msg string) string {
	return encryption.Base64StdEncodeToString([]byte(msg))
}

// DecryptionVerMsg 解密报文
func (t *TaskVul) DecryptionVerMsg(msg string) (string, string, map[string]interface{}, map[string]interface{}, error) {
	var (
		scheme    = "http"
		method    = "GET"
		path      = "/"
		host      = ""
		space     = false
		body      = ""
		url       = ""
		headerMap = make(map[string]interface{}, 0)
		bodyMap   = make(map[string]interface{}, 0)
	)
	httpsnum := strings.Count(msg, "https")
	if httpsnum > 0 {
		scheme = "https"
	}
	strings.Replace(msg, "\r\n", "\n", -1)
	msgArray := strings.Split(msg, "\n")
	for i := 0; i < len(msgArray); i++ {
		msgLineArray := strings.Split(msgArray[i], ":")
		if i == 0 {
			//msgZeroArray := strings.Split(msgLineArray[0], " ")
			msgZeroArray := msgLineArray[0]
			if len(msgZeroArray) < 2 {
				return method, url, headerMap, bodyMap, errors.New("报文格式错误")
			}
			if strings.HasPrefix(strings.ToLower(msgZeroArray), "get") {
				method = "GET"
			} else if strings.HasPrefix(strings.ToLower(msgZeroArray), "post") {
				method = "POST"
			} else if strings.HasPrefix(strings.ToLower(msgZeroArray), "put") {
				method = "PUT"
			} else if strings.HasPrefix(strings.ToLower(msgZeroArray), "delete") {
				method = "DELETE"
			} else if strings.HasPrefix(strings.ToLower(msgZeroArray), "patch") {
				method = "PATCH"
			}
			//method = msgZeroArray[0]
			path = strings.Replace(msgZeroArray, method, "", 1)
			path = strings.TrimSpace(strings.Replace(path, "HTTP/1.1", "", -1))
			continue
		}
		if msgLineArray[0] == "Host" {
			if len(msgLineArray) >= 3 {
				host = msgLineArray[1] + ":" + msgLineArray[2]
			} else {
				host = msgLineArray[1]
			}
			continue
		}
		if len(msgLineArray) == 1 && len(msgLineArray[0]) == 0 {
			space = true
			continue
		}
		if space && len(msgLineArray) > 0 && len(msgLineArray[0]) > 0 {
			body = msgLineArray[0]
			continue
		}
		if len(msgLineArray) == 2 {
			headerMap[msgLineArray[0]] = strings.TrimSpace(msgLineArray[1])
		}
	}
	if len(host) == 0 {
		return method, url, headerMap, bodyMap, errors.New("报文格式错误")
	}
	url = scheme + "://" + strings.TrimSpace(strings.Trim(host, "\r")) + path
	if len(body) > 0 {
		err := json.Unmarshal([]byte(body), &bodyMap)
		if err != nil {
			bodyMap = decryptionBody(body)
		}
	}
	return method, url, headerMap, bodyMap, nil
}

func decryptionBody(body string) map[string]interface{} {
	var result = make(map[string]interface{}, 0)
	bodyArray := strings.Split(body, "&")
	for i := 0; i < len(bodyArray); i++ {
		tmp := strings.Split(bodyArray[i], "=")
		if len(tmp) > 1 {
			result[tmp[0]] = tmp[1]
		}
	}
	return result
}

// DecryptionVerMsg 解密报文
func (t *TaskVul) BuildRespVerMsg(proto, state string, header map[string]string, body string) (string, error) {
	resp_str := fmt.Sprintf("%s %s ", proto, state)
	for k, v := range header {
		resp_str += k + ":" + v + "\r\n"
	}
	resp_str += "\r\n" + body
	return resp_str, nil
}

// UpdateVerMsgById 根据id修改请求报文和响应报文
func (t *TaskVul) UpdateVerMsgById(ctx context.Context, id int, verMsg, respVerMsg string, originalVerMsg string) error {
	var (
		originVermsgTemp   []ScriptResultDetailVulProve
		vermsgTemp         []ScriptResultDetailVulProve
		payload            string
		payloadSuccessFlag string
	)
	json.Unmarshal([]byte(originalVerMsg), &originVermsgTemp)
	if len(originVermsgTemp) > 0 {
		payload = originVermsgTemp[0].Payload
		payloadSuccessFlag = originVermsgTemp[0].PayloadSuccessFlag
	}
	vermsgTemp = append(vermsgTemp, ScriptResultDetailVulProve{
		Request:            verMsg,
		Response:           respVerMsg,
		Payload:            payload,
		PayloadSuccessFlag: payloadSuccessFlag,
	})
	vermsgTempByte, _ := json.Marshal(vermsgTemp)
	var taskVul mysqls.TaskVul
	key := []byte("9876787656785679")
	aesEcb := aesEncryption.AesEcb{}
	verMsgTempString := hex.EncodeToString(aesEcb.AesEncryptECB(vermsgTempByte, key))
	return taskVul.UpdateVerMsgById(ctx, id, verMsgTempString)
	//return taskVul.UpdateVerMsgById(ctx, id, string(vermsgTempByte))
}

// UpdateStatusById 根据id修改漏洞状态
func (t *TaskVul) UpdateStatusById(ctx context.Context, id int, status int) error {
	var taskVul mysqls.TaskVul
	return taskVul.UpdateStatusById(ctx, id, status)
}

// 批量删除 通过task_id
func (t *TaskVul) DelTaskInfoByTaskId(ctx context.Context, taskIds []int) error {
	var taskVulModel mysqls.TaskVul
	return taskVulModel.DeleteByTaskIds(ctx, taskIds)
}

// GetTargetStats 计算目标的统计信息
func (t *TaskVul) GetTargetStats(ctx context.Context, targetId int) (int, int, [6]int, error) {
	var (
		taskVul     mysqls.TaskVul
		vulNum      int                          //漏洞总数
		risklevel   = enums.TargetRiskLowNoFound //目标等级，默认为4->未发现
		vulNumArray [6]int                       //每个等级的数量，元素含义分别为：无漏洞个数/致命漏洞个数/高危漏洞个数/中危漏洞个数/低危漏洞个数/信息漏洞个数
	)
	vulRes, err := taskVul.GetTaskVulListByTargetIds(ctx, []int{targetId}, enums.VulDataTypOne)
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
func (t *TaskVul) GetTargetStatsBytargetIds(ctx context.Context, targetIds []int) (map[int]int, map[int]int, map[int][6]int, error) {
	var (
		taskVul     mysqls.TaskVul
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

	vulRes, err := taskVul.GetTaskVulListByTargetIds(ctx, targetIds, enums.VulDataTypOne)
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

// GetsByTaskId 依据任务ID 获取所有漏洞信息
func (t *TaskVul) GetsByTaskId(ctx context.Context, taskId int, dataType int) []mysqls.TaskVul {
	var taskVul mysqls.TaskVul
	list := taskVul.GetsByTaskId(ctx, taskId, dataType)
	for i := range list {
		t.DecryptTaskVul(&list[i])
	}
	return list
}

// GetsByTargetId 依据任务ID 获取所有漏洞信息
func (t *TaskVul) GetsByTargetId(ctx context.Context, targetId int, dataType int) []mysqls.TaskVul {
	var taskVul mysqls.TaskVul
	list := taskVul.GetsByTargetId(ctx, targetId, dataType)
	for i := range list {
		t.DecryptTaskVul(&list[i])
	}
	return list
}

func (t *TaskVul) GetParamIpPort(vulParam string, domainMap map[string]string) (string, string) {
	params := make([]map[string]string, 0)
	err := json.Unmarshal([]byte(vulParam), &params)
	if err != nil {
		return "", ""
	}

	var ip, port string
	for _, param := range params {
		if strings.Contains(param["key"], "ip") {
			ip = param["value"]
		}
		if strings.Contains(param["key"], "port") {
			port = param["value"]
		}
	}
	if ip != "" && port != "" {
		return ip, port
	}

	for _, param := range params {
		if strings.Contains(param["key"], "url") {
			_, host, port, err := network.ParseUrl(param["value"])
			if err != nil {
				return "", ""
			}
			var ip string
			if domainMap[host] == "" {
				ip, err = network.ResolveDomain(host)
				domainMap[host] = ip
			} else {
				ip = domainMap[host]
			}
			if err != nil {
				return "", ""
			}
			return ip, port
		}
	}
	return "", ""
}

func (t *TaskVul) GetParamLocation(vulParam string) string {
	params := make([]map[string]string, 0)
	err := json.Unmarshal([]byte(vulParam), &params)
	if err != nil {
		return ""
	}
	for _, param := range params {
		if strings.Contains(param["key"], "url") {
			return param["value"]
		}
	}
	var ip, port string
	for _, param := range params {
		if strings.Contains(param["key"], "ip") {
			ip = param["value"]
		}
		if strings.Contains(param["key"], "port") {
			port = param["value"]
		}
	}
	if port != "" {
		return ip + ":" + port
	} else {
		return ip
	}
}

// GetByTargetResultId 通过结果表ID获取
func (t *TaskVul) GetByTargetResultId(ctx context.Context, targetResultId int) (mysqls.TaskVul, error) {
	var taskVulModel mysqls.TaskVul
	return taskVulModel.GetByTargetResultId(ctx, targetResultId)
}

// GetByTargetResultId 通过结果表ID获取
func (t *TaskVul) GetByTargetResultIds(ctx context.Context, targetResultIds []int, targetId int) []mysqls.TaskVul {
	var taskVulModel mysqls.TaskVul
	return taskVulModel.GetByTargetResultIds(ctx, targetResultIds, targetId)

}

// GetTaskVulCount 获取任务漏洞总数或根据开始时间获取任务漏洞总数
func (t *TaskVul) GetTaskVulCount(ctx context.Context, startTime string) (int, int) {
	var taskVulModel mysqls.TaskVul
	total, filterTotal := taskVulModel.GetTaskVulCount(ctx, startTime)
	return int(total), int(filterTotal)
}

// GetTaskVulRiskStat 任务漏洞风险统计
func (t *TaskVul) GetTaskVulRiskStat(ctx context.Context, uid int, role int) []mysqls.TaskVulRiskStat {
	var taskVulModel mysqls.TaskVul
	return taskVulModel.GetTaskVulRiskStat(ctx, uid, role)
}

//// GetTargetVerifyStatus 计算与获取任务的  --  验证状态
//func (t *TaskVul) GetTargetVerifyStatus(ctx context.Context, targetIds []int) string {
//	var (
//		taskVul mysqls.TaskVul
//	)
//
//	vulRes, err := taskVul.GetTaskVulListByTargetIds(ctx, targetIds, enums.VulDataTypOne)
//	if err != nil {
//		return ""
//	}
//	verifyNum := enums.VulStatusNotVerify
//	for _, item := range vulRes {
//		if item.Status == enums.VulStatusVerifyUsed || item.Status == enums.VulStatusVerifySuccess {
//			if item.Status > verifyNum {
//				verifyNum = item.Status
//			}
//		}
//	}
//	return enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(verifyNum)
//}

// GetTargetVerifyStatus 计算与获取任务的  --  验证状态
func (t *TaskVul) GetTargetVerifyStatus(ctx context.Context, targetIds []int) string {
	var (
		taskVul mysqls.TaskVul
	)

	vulRes, err := taskVul.GetTaskVulListByTargetIds(ctx, targetIds, enums.VulDataTypOne)
	if err != nil {
		return ""
	}
	verifyNum := enums.VulStatusNotVerify
	for _, item := range vulRes {
		if item.Status == enums.VulStatusVerifyExist || item.Status == enums.VulStatusConfirmExist {
			if item.Status > verifyNum {
				verifyNum = item.Status
			}
		}
	}
	return enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(verifyNum)
}

// TaskVulTypeStat 任务漏洞类型统计
func (t *TaskVul) TaskVulTypeStat(ctx context.Context, mode int, uid int, role int) ([]string, []int, error) {
	var (
		startTime   string
		currentTime = time.Now()
	)
	switch mode {
	case enums.StatModeWeek:
		startTime = currentTime.AddDate(0, 0, -7).String()
	case enums.StatModeMonth:
		startTime = currentTime.AddDate(0, -1, 0).String()
	case enums.StatModeYear:
		startTime = currentTime.AddDate(-1, 0, 0).String()
	}

	var taskVulModel mysqls.TaskVul
	taskVulTypeStatList := taskVulModel.GetTaskVulTypeStat(ctx, startTime, uid, role)
	vulTypeList := make([]string, 0)
	vulCountList := make([]int, 0)
	for _, item := range taskVulTypeStatList {
		vulTypeList = append(vulTypeList, enums.ToolsVulnerabilityEnum.GetTypeEnum(item.Type))
		vulCountList = append(vulCountList, item.Count)
	}
	return vulTypeList, vulCountList, nil
}

// TaskVulFindTrendStat 获取任务漏洞发现趋势统计
func (t *TaskVul) TaskVulFindTrendStat(ctx context.Context, mode int, uid int, role int) ([]string, []int, error) {
	//处理参数
	var startTime, dateFormat string
	dateList := make([]string, 0)
	switch mode {
	case enums.StatModeWeek:
		dateList, startTime, dateFormat = utils.WeekDateList()
	case enums.StatModeMonth:
		dateList, startTime, dateFormat = utils.MonthDateList()
	case enums.StatModeYear:
		dateList, startTime, dateFormat = utils.YearDateList()
	}

	var taskVulService mysqls.TaskVul
	taskVulFindTrendStatList := taskVulService.GetTaskVulFindTrendStat(ctx, startTime, dateFormat, uid, role)
	taskVulFindTrendStatMap := make(map[string]int)
	for _, item := range taskVulFindTrendStatList {
		taskVulFindTrendStatMap[item.Date] = item.Count
	}
	countList := make([]int, 0)
	for _, item := range dateList {
		countList = append(countList, taskVulFindTrendStatMap[item])
	}
	return dateList, countList, nil
}

// 依据目标IDs获取所有目标信息收集数据，并返回map
func (t *TaskVul) AllByTargetIds(ctx context.Context, targetIds []int) (returnData map[int][]mysqls.TaskVul) {
	var taskVulModel mysqls.TaskVul
	taskVul, _ := taskVulModel.GetTaskVulListByTargetIds(ctx, targetIds, enums.VulDataTypOne)
	returnData = make(map[int][]mysqls.TaskVul)
	for _, item := range taskVul {
		// 如果不存在，则初始化
		if _, ok := returnData[item.TargetID]; !ok {
			returnData[item.TargetID] = make([]mysqls.TaskVul, 0)
		}
		returnData[item.TargetID] = append(returnData[item.TargetID], item)
	}
	return
}

// 获取平台发现的所有漏洞 分页 依据发现时间倒叙
func (t *TaskVul) AllByPage(ctx context.Context, pageNum, pageSize int) (returnData []mysqls.TaskVul, count int64) {
	var model mysqls.TaskVul
	returnData, count = model.AllByPage(ctx, pageNum, pageSize)
	return
}

// GetVulStatsByTaskId 获取漏洞统计数据（聚合查询）
func (t *TaskVul) GetVulStatsByTaskId(ctx context.Context, taskId int, dataType int) []mysqls.VulStatsResult {
	var taskVul mysqls.TaskVul
	return taskVul.GetVulStatsByTaskId(ctx, taskId, dataType)
}

// AddTaskVulByScannerResult 插入一条检测漏洞
func (t *TaskVul) AddTaskVulByScannerResult(ctx context.Context, target mysqls.TaskTarget, riskData invoke.Risk, scriptResultMap map[string]string, vulLib mysqls.VulLibraries, targetResultID int, dataType int) error {
	log.Infof("[DEBUG-VUL] AddTaskVulByScannerResult: 准备插入漏洞, TargetID=%d, VulName=%s, DataType=%d", target.ID, vulLib.Name, dataType)
	var taskVul mysqls.TaskVul

	// 销售许可证临时添加

	taskVul.DataType = dataType
	//taskVul.Name = scriptResult.Libraries.Name
	// taskVul.Name = hex.EncodeToString(aesEcb.AesEncryptECB([]byte(vulLib.Name), key))

	taskVul.Name = vulLib.Name

	taskVul.Class = vulLib.Class
	taskVul.Type = vulLib.Type
	taskVul.Risk = vulLib.Risk

	taskVul.VulID = vulLib.Cve + "," + vulLib.Cnvd + "," + vulLib.Cnnvd
	taskVul.Description = vulLib.Description

	taskVul.FixSuggest = vulLib.FixSuggest

	location := extractLocationFromRequest(scriptResultMap, riskData)
	taskVul.Location = location

	//taskVul.Pocname = scriptResult.Script.Pocname
	taskVul.Pocname = scriptResultMap["pocname"]

	taskVul.TargetResultID = targetResultID
	//taskVul.VulParam = scriptResult.Script.ScriptParam
	taskVul.TaskID = target.TaskID
	taskVul.TargetID = target.ID
	//taskVul.Location = scriptResult.Result.Location

	taskVul.PublishedTime = vulLib.PublishedTime
	taskVul.Cvss = vulLib.CvssScore

	taskVul.Status = enums.VulStatusVerifyExist
	taskVul.TestStatus = enums.VulTestStatusPrinciple
	if riskData.ScriptType == "version_match" {
		taskVul.TestStatus = enums.VulTestStatusVersionMatch
		taskVul.Status = enums.VulStatusNotVerify
	}
	taskVul.CreateTime = time.Now()
	taskVul.UpdateTime = time.Now()
	taskVul.TargetUrl = target.TargetURL
	//taskVul.TargetUrl = scriptResult.Result.TargetUrl
	taskVul.DecisionVulId = vulLib.VulID
	//处理漏洞截图、请求报文、响应报文
	var (
		scriptDetailTemp   map[string]string
		verMsgTmp          = make([]ScriptResultDetailVulProve, 0)
		snapshottmp        string
		payload            string
		payloadSuccessFlag string
	)
	json.Unmarshal([]byte(riskData.Details), &scriptDetailTemp)
	if v, ok := scriptDetailTemp["screen_data"]; ok { //截图
		snapshottmp = v
		delete(scriptDetailTemp, "screen_data") //剔除结果中的截图数据
	}
	taskVul.Snapshot = snapshottmp
	if v, ok := scriptDetailTemp["payload"]; ok { //payload
		payload, _ = url.QueryUnescape(v)
		scriptDetailTemp["payload"], _ = url.QueryUnescape(v)
	}
	if v, ok := scriptDetailTemp["payload_success_flag"]; ok { //payload_success_flag
		payloadSuccessFlag = v
	}
	if _, ok := scriptDetailTemp["request"]; ok { //在脚本结果中移除request
		delete(scriptDetailTemp, "request")
	}
	if _, ok := scriptDetailTemp["response"]; ok { //在脚本结果中移除response
		delete(scriptDetailTemp, "response")
	}

	if v, ok := scriptDetailTemp["vul_prove"]; ok { //多个报文
		json.Unmarshal([]byte(v), &verMsgTmp)
	}
	if len(verMsgTmp) == 0 {
		verMsgTmp = append(verMsgTmp, ScriptResultDetailVulProve{
			Request:            string(riskData.Request),
			Response:           string(riskData.Response),
			Payload:            payload,
			PayloadSuccessFlag: payloadSuccessFlag,
		})
	}

	verMsgTmpByte, _ := json.Marshal(verMsgTmp)
	taskVul.VerMsg = string(verMsgTmpByte)

	tmpDetail, _ := json.Marshal(scriptDetailTemp)
	taskVul.VulResult = string(tmpDetail)

	taskVul.PatchUrl = vulLib.PatchUrl
	taskVul.VulNumber = utils.JoinNonEmpty(vulLib.Cve, vulLib.Cnvd, vulLib.Cnnvd, vulLib.Cncve, vulLib.Bugtraq)
	taskVul.Cvss = vulLib.CvssScore

	// 检查漏洞是否存在
	existVul, err := taskVul.GetTaskVulDetail(ctx, target.ID, taskVul.Pocname, taskVul.Location)
	if err != nil {
		return errors.New("GetTaskVulDetail error: " + err.Error())
	}

	if existVul.ID > 0 {
		// 存在则更新
		taskVul.ID = existVul.ID
		taskVul.CreateTime = existVul.CreateTime
		err = taskVul.UpdateTaskVul(ctx)
		if err != nil {
			return errors.New("UpdateTaskVul error: " + err.Error())
		}
	} else {
		// 不存在则新增
		err = taskVul.AddTaskVul(ctx)
		if err != nil {
			return errors.New("AddTaskVul error: " + err.Error())
		}
	}

	// 发现漏洞记录周期信息
	var (
		vulLifecycleSrv VulLifecycle
		taskSrv         TaskTaskInfo
		userSrv         User
	)
	// 通过task 获取创建的任务名称 和 用户名
	taskInfo, err := taskSrv.GetTaskInfoByTaskId(ctx, taskVul.TaskID)
	if err != nil {
		taskInfo.TaskName = "未知任务"
	}
	userInfo, err := userSrv.GetUserDetail(ctx, taskInfo.UserID)
	if err != nil {
		userInfo.Username = "未知用户"
	}
	content := vulLifecycleSrv.ConstructDiscoveryContent(userInfo.Username, taskInfo.TaskName, taskVul.TargetUrl, taskVul.Name, taskVul.Status)
	vulLifecycleSrv.AddVulLifecycle(ctx, taskVul.Pocname, taskVul.Name, taskVul.Location, content, taskInfo.TaskName, taskVul.TaskID)

	return nil
}

// GetVulTaskIDByIDs 根据id获取任务id
func (t *TaskVul) GetVulTaskIDByIDs(ctx context.Context, ids any) (map[int]string, error) {
	var taskVul mysqls.TaskVul
	vulRes, err := taskVul.GetTaskVulListByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	res := make(map[int][]string)
	for _, v := range vulRes {
		taskID := v.TaskID
		vulIDStr := strconv.Itoa(v.ID)
		res[taskID] = append(res[taskID], vulIDStr)
	}
	finalRes := make(map[int]string, len(res))
	for taskID, idList := range res {
		finalRes[taskID] = strings.Join(idList, ",")
	}
	return finalRes, nil
}

// GetTargetStatsBytargetIdsForAsset 计算资产目标的统计信息
func (t *TaskVul) GetTargetStatsBytargetIdsForAsset(ctx context.Context, targetIds []int) (map[int]int, map[int]int, map[int][6]int, error) {
	var (
		taskVul     mysqls.TaskVul
		vulNum      = make(map[int]int, 0)    //漏洞总数
		risklevel   = make(map[int]int, 0)    //目标等级，默认为4->未发现
		vulNumArray = make(map[int][6]int, 0) //每个等级的数量，元素含义分别为：无漏洞个数/致命漏洞个数/高危漏洞个数/中危漏洞个数/低危漏洞个数/信息漏洞个数
	)
	//初始化
	for i := 0; i < len(targetIds); i++ {
		vulNum[targetIds[i]] = 0
		risklevel[targetIds[i]] = enums.SafeAsset
		vulNumArray[targetIds[i]] = [6]int{}
	}
	vulRes, err := taskVul.GetTaskVulListByTargetIds(ctx, targetIds, enums.VulDataTypOne)
	if err != nil || len(vulRes) == 0 {
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
		if v[1] > 0 || v[2] > 0 { // 严重
			risklevel[k] = enums.HighRiskAsset
		} else if v[3] > 0 { //中
			risklevel[k] = enums.MiddleRiskAsset
		} else if v[4] > 0 { //低->低
			risklevel[k] = enums.LowRiskAsset
		}
	}
	return vulNum, risklevel, vulNumArray, nil
}

// GetTaskDeduplicationByLocationList 按照 name, pocname, location 进行去重
func (t *TaskVul) GetTaskDeduplicationByLocationList(ctx context.Context, taskId, targetId, vultype, risk int, search, ip, location string, page, limit, status, verifyType, uid, role int) ([]mysqls.DeduplicatedVul, int64, error) {
	var taskVul mysqls.TaskVul
	return taskVul.GetTaskDeduplicationByLocationList(ctx, taskId, targetId, vultype, risk, search, ip, location, enums.VulDataTypOne, page, limit, status, verifyType, uid, role)
}
