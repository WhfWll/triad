package httpclients

import (
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
)

// QianjiVulCheck 千机-漏洞信息检测
type QianjiVulCheck struct {
	TaskID string `json:"TaskID"`
	Status bool   `json:"status"`
	Result struct {
		IsSuccess bool   `json:"is_success"` // 如果是false证明是误报
		Details   string `json:"details"`    // 误报详情解释
		Request   string `json:"request"`
		Response  string `json:"response"`
		Runtimeid string `json:"runtimeid"`
	} `json:"result"`
}

// VulInfoCheck 漏洞信息检测【是否是误报】
func VulInfoCheck(httpFlow, vulInfo string) (falseAlarm bool, content string, err error) {
	h, err := httpclient.NewHttpSend("service_qianji", "/v1/VulCheck")
	if err != nil {
		return
	}
	param := map[string]interface{}{
		"HTTPFlow": httpFlow,
		"VulInfo":  vulInfo,
	}
	h.SetBody(param)
	response, err := h.Post()
	if err != nil {
		return
	}
	var res QianjiVulCheck
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}
	// qianji接口IsSuccess是false代表是误报
	if res.Result.IsSuccess == false {
		falseAlarm = true
		content = res.Result.Details
	}
	return
}
