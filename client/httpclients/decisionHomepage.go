package httpclients

import (
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
)

type DecisionToolInfoStatRes struct {
	Code int `json:"code"`
	Data struct {
		VulCount    int `json:"vulCount"`
		FingerCount int `json:"fingerCount"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// DecisionToolInfoStat 工具信息统计 - 漏洞数量、指纹数量
func DecisionToolInfoStat() (DecisionToolInfoStatRes, error) {
	var res DecisionToolInfoStatRes
	h, err := httpclient.NewHttpSend("service_decision", "/decision/homepage/opentoolinfostat")
	if err != nil {
		return res, err
	}
	//设置请求头
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	//发送请求
	response, err := h.Get()
	if err != nil {
		return res, err
	}
	//解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
