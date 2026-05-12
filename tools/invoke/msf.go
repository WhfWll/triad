package invoke

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
	"gitlabee.4dogs.cn/common/config"
	"io/ioutil"
	"net/http"
	"regexp"
	"smart/models/redises"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"
)

func (c *CallInfo) invokeMsfScript(ctx context.Context, scriptType string, scriptContent string, callBackFunc func(context.Context, *CallInfo, string)) {
	token, err := c.Auth(ctx)
	if err != nil {
		log.Println(err)
		callBackFunc(ctx, c, err.Error())
	}
	id, err := c.CreateConsole(ctx, token)
	defer c.DestroyConsole(ctx, token, id)
	sendData, vulReg, rhost := c.BuildSendData(ctx, scriptContent)
	c.WriteConsole(ctx, token, id, sendData)
	err = c.ReadConsoleHandleResult(ctx, token, rhost, id, vulReg, callBackFunc)
	if err != nil {
		log.Error(err)
	}
	callBackFunc(ctx, c, "end")
}

func (c *CallInfo) Auth(ctx context.Context) (string, error) {
	apiUsername := "msf"
	apiPassword := "4dogs.cn"
	options := []interface{}{"auth.login", apiUsername, apiPassword}
	body, err := c.sendCommand(ctx, options)
	if err != nil {
		return "", err
	}
	var bodyMap map[string]interface{}
	err = msgpack.Unmarshal(body, &bodyMap)
	if err != nil {
		return "", err
	}
	if bodyMap["token"] == nil {
		return "", errors.New("无法获取到token")
	}
	return string(bodyMap["token"].([]uint8)), nil
}

func (c *CallInfo) RunModule(ctx context.Context, token string, execType, pluginName, lHost, lPort, rHost, rPort string) (string, error) {
	options := []interface{}{"module.execute", token, execType, pluginName}
	tempMap := make(map[string]string, 0)
	if lHost != "" {
		tempMap["LHOST"] = lHost
	}
	if lPort != "" {
		tempMap["LPORT"] = lPort
	}
	if rHost != "" {
		tempMap["RHOST"] = rHost
	}
	if rPort != "" {
		tempMap["RPORT"] = rPort
	}
	options = append(options, tempMap)
	body, err := c.sendCommand(ctx, options)
	if err != nil {
		return "", nil
	}
	log.Println(body)
	return "", nil
}

func (c *CallInfo) sendCommand(ctx context.Context, options []interface{}) ([]byte, error) {
	var (
		clientMsf map[string]map[string]interface{}
		apiUrl    string
	)
	err := config.Load("client", &clientMsf)
	if err == nil {
		if clientMsf["service_msf"] != nil {
			serviceMsf := clientMsf["service_msf"]
			apiUrl = serviceMsf["base_uri"].(string)
		}
	}
	if apiUrl == "" {
		apiUrl = "http://127.0.0.1:55553/api"
	}
	requestBody, err := msgpack.Marshal(options)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "binary/message-pack")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func (c *CallInfo) ListSessions(ctx context.Context, token string) (map[string]map[string]string, error) {
	options := []interface{}{"session.list", token}
	body, err := c.sendCommand(ctx, options)
	var bodyMap map[int]map[string]interface{}
	err = msgpack.Unmarshal(body, &bodyMap)
	sessionMap := make(map[string]map[string]string, 0)
	if err != nil {
		return sessionMap, err
	}
	for sid, sMap := range bodyMap {
		tempMap := make(map[string]string, 0)
		for key, value := range sMap {
			if key == "session_port" {
				tempMap[key] = strconv.Itoa(int(value.(uint16)))
			} else if key == "workspace" || key == "platform" || key == "session_host" || key == "routes" || key == "target_host" {
				tempMap[key] = value.(string)
			} else {
				tempMap[key] = string(value.([]uint8))
			}
		}
		sessionMap[strconv.Itoa(sid)] = tempMap
	}
	return sessionMap, nil
}

func (c *CallInfo) CreateConsole(ctx context.Context, token string) (string, error) {
	options := []interface{}{"console.create", token}
	body, err := c.sendCommand(ctx, options)
	if err != nil {
		return "", err
	}
	var bodyMap map[string]interface{}
	err = msgpack.Unmarshal(body, &bodyMap)
	if err != nil {
		return "", err
	}
	return bodyMap["id"].(string), nil
}

func (c *CallInfo) DestroyConsole(ctx context.Context, token interface{}, id string) (bool, error) {
	options := []interface{}{"console.destroy", token, id}
	body, err := c.sendCommand(ctx, options)
	if err != nil {
		return true, err
	}
	var bodyMap map[string]interface{}
	err = msgpack.Unmarshal(body, &bodyMap)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *CallInfo) WriteConsole(ctx context.Context, token string, id, data string) (bool, error) {
	options := []interface{}{"console.write", token, id, data}
	body, err := c.sendCommand(ctx, options)
	if err != nil {
		return false, err
	}
	var bodyMap map[string]interface{}
	err = msgpack.Unmarshal(body, &bodyMap)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *CallInfo) ReadConsole(ctx context.Context, token string, id string, callBackFunc func(context.Context, *CallInfo, string)) (string, error) {
	options := []interface{}{"console.read", token, id}
	var data string
	for {
		time.Sleep(1 * time.Second)
		body, err := c.sendCommand(ctx, options)
		if err != nil {
			break
		}
		var bodyMap map[string]interface{}
		err = msgpack.Unmarshal(body, &bodyMap)
		if err != nil {
			continue
		}
		if bodyMap["data"] == nil {
			break
		}
		tempData := string(bodyMap["data"].([]uint8))
		for _, value := range strings.Split(tempData, "\n") {
			if strings.HasPrefix(value, "[") && !strings.HasPrefix(value, "[!]") {
				data += value + "\n"
			}
		}
		if bodyMap["busy"].(bool) == false {
			break
		}
	}
	return data, nil
}

func (c *CallInfo) ReadConsoleHandleResult(ctx context.Context, token, rhost, id, vulReg string, callBackFunc func(context.Context, *CallInfo, string)) error {
	options := []interface{}{"console.read", token, id}
	var data string
	var sessionData string
	// 第一步 通过msf api 获取上个命令执行结果
	for {
		time.Sleep(1 * time.Second)
		body, err := c.sendCommand(ctx, options)
		if err != nil {
			break
		}
		var bodyMap map[string]interface{}
		err = msgpack.Unmarshal(body, &bodyMap)
		if err != nil {
			continue
		}
		if bodyMap["data"] == nil {
			break
		}
		tempData := string(bodyMap["data"].([]uint8))
		for _, value := range strings.Split(tempData, "\n") {
			if !strings.HasPrefix(value, "[") || strings.HasPrefix(value, "[!]") {
				continue
			}
			callBackFunc(ctx, c, value)
			//if strings.Contains(value, "failed") || strings.Contains(value, "aborted") {
			//	data = ""
			//}
			if strings.HasPrefix(value, "[+]") {
				data += value
			}
			if strings.HasPrefix(value, "[*]") {
				data += value
			}
			if strings.Contains(strings.ToLower(value), "session") {
				data += strings.ToLower(value)
				sessionData += strings.ToLower(value)
			}
		}
		if bodyMap["busy"].(bool) == false {
			break
		}
	}
	// 第二步 通过正则分析调用结果
	//if data != "" && (strings.Contains(data, "vulnerable") || strings.Contains(data, "session")) {
	if vulReg == "" {
		vulReg = `session ([0-9]+)`
	}
	re := regexp.MustCompile(vulReg)
	regResult := re.FindStringSubmatch(data)
	if len(regResult) > 0 {
		//ipPattern := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
		locationPattern := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`)
		matchLocation := locationPattern.FindStringSubmatch(data)
		//matchIp := ipPattern.FindStringSubmatch(data)

		tempDetails := make(map[string]string, 0)
		if len(regResult) > 1 {
			tempDetails["reg_result"] = regResult[1]
		} else {
			tempDetails["reg_result"] = regResult[0]
		}
		if matchLocation != nil {
			tempDetails["location"] = matchLocation[0]
		}
		tempDetails["target"] = rhost
		//if matchIp != nil {
		//tempDetails["target"] = matchIp[0]
		//}
		sessionPattern := regexp.MustCompile(`session ([0-9]+)`)
		matchSession := sessionPattern.FindStringSubmatch(sessionData)
		if matchSession != nil {
			tempDetails["session"] = matchSession[1]
			var (
				redisHash   redises.RedisHash
				redisCommon redises.RedisCommon
			)
			sessionKey := enums.DecisionShell
			cmdKey := enums.RemoteShellPreKey + strconv.FormatInt(time.Now().Unix(), 10)
			shellInfoMap := map[string]string{
				"createTime": time.Now().String(),
				"shellType":  "remote",
				"sessionId":  matchSession[1],
			}
			shellInfoMapByte, _ := json.Marshal(shellInfoMap)
			err := redisHash.SetHashHset(ctx, sessionKey, cmdKey, string(shellInfoMapByte))
			redisCommon.Expire(ctx, sessionKey, time.Duration(24*time.Hour))
			if err != nil {
				log.Error(err)
			}
		}
		tempDetailsByte, _ := json.Marshal(tempDetails)
		tempMap := map[string]string{"CreatedAt": "1689782484", "Hash": "1689782484", "Details": string(tempDetailsByte), "Severity": "2"}
		tempByte, _ := json.Marshal(tempMap)
		callBackFunc(ctx, c, string(tempByte))
	}
	return nil
}

// BuildSendData 构造发送数据
func (c *CallInfo) BuildSendData(ctx context.Context, scriptContent string) (data string, vulReg string, rhost string) {
	for _, line := range strings.Split(scriptContent, "\n") {
		if strings.Contains(line, "vul_reg") {
			vulReg = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, "#", ""), "vul_reg:", ""))
		}
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		line = strings.ToLower(line)
		if strings.Contains(line, "rhost") || strings.Contains(line, "lhost") || strings.Contains(line, "rport") || strings.Contains(line, "lport") {
			continue
		}
		data += line + "\n"
	}
	for _, param := range c.ToolParamList {
		if param.ParamName == "ip" {
			data += "set RHOSTS " + param.ParamValue + "\n"
			rhost = param.ParamValue
		} else if param.ParamName == "reverseHost" {
			data += "set LHOST " + param.ParamValue + "\n"
		} else if param.ParamName == "lPort" {
			data += "set LPORT " + param.ParamValue + "\n"
		} else if strings.HasSuffix(param.ParamName, "port") {
			data += "set RPORT " + param.ParamValue + "\n"
		} else {
			data += "set " + param.ParamName + " " + param.ParamValue + "\n"
		}
	}
	data += "run -z\n"
	return
}

// HandleResult 处理返回结果数据
func (c *CallInfo) HandleResult(ctx context.Context, dataString string) (string, error) {
	successData := ""
	for _, value := range strings.Split(dataString, "\n") {
		if strings.HasPrefix(value, "[+]") {
			successData += value + "\n"
		}
	}
	if successData != "" {
		tempMap := map[string]string{"CreatedAt": "1689782484", "Details": strings.Trim(dataString, "\n")}
		tempByte, err := json.Marshal(tempMap)
		if err != nil {
			return "", err
		}
		return string(tempByte), nil
	} else {
		return dataString, nil
	}
}

func (c *CallInfo) WriteSession(ctx context.Context, token, sid, data string) error {
	options := []interface{}{"session.meterpreter_write", token, sid, data + "\n"}
	body, err := c.sendCommand(ctx, options)
	var bodyMap interface{}
	err = msgpack.Unmarshal(body, &bodyMap)
	if err != nil {
		return err
	}
	return nil
}

func (c *CallInfo) ReadSession(ctx context.Context, token string, sid string) (interface{}, error) {
	//time.Sleep(3000 * time.Millisecond)
	//options := []interface{}{"session.meterpreter_read", token, sid + "\n"}
	//body, err := c.sendCommand(ctx, options)
	//var bodyMap map[string]interface{}
	//err = msgpack.Unmarshal(body, &bodyMap)
	//return string(bodyMap["data"].([]uint8)), err
	for i := 0; i < 5; i++ {
		time.Sleep(400 * time.Millisecond)
		options := []interface{}{"session.meterpreter_read", token, sid + "\n"}
		body, err := c.sendCommand(ctx, options)
		var bodyMap map[string]interface{}
		err = msgpack.Unmarshal(body, &bodyMap)
		if err != nil {
			return "", err
		}
		result := bodyMap["data"].([]uint8)
		if len(result) > 0 {
			return string(result), nil
		}
	}
	return "", nil
}
