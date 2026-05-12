package httpclients

import (
	"encoding/json"
	"errors"
	"gitlabee.4dogs.cn/common/httpclient"
)

type CreateFlowTaskReq struct {
	FlowTaskId int                      `json:"flowTaskId"`
	Node       string                   `json:"node"`
	MitmHost   string                   `json:"mitmHost"`
	Port       string                   `json:"port"`
	ExpireTime int                      `json:"expireTime"`
	ListenData []CreateFlowTaskReqItems `json:"listenData"`
}

type CreateFlowTaskReqItems struct {
	Host         string `json:"host"`
	FlowTaskId   int    `json:"flowTaskId"`
	FlowTargetId int    `json:"flowTargetId"`
}

type flowResp struct {
	Code int `json:"code"`
	Data struct {
	} `json:"data"`
	Msg string `json:"msg"`
}

// 开始任务
func CreateFlowTask(param map[string]interface{}) error {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/flow/createtask")
	if err != nil {
		return err
	}
	h.SetBody(param) //设置get请求参数，post请求不要使用
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Post() //发送get请求
	if err != nil {
		return err
	}
	var tmpResp flowResp
	err = json.Unmarshal(response, &tmpResp)
	if err != nil {
		return err
	}
	if tmpResp.Code != 200 && len(tmpResp.Msg) > 0 {
		return errors.New(tmpResp.Msg)
	}
	return nil
}

// 操作任务
func OprateFlowTask(param map[string]interface{}) error {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/flow/opratetask")
	if err != nil {
		return err
	}
	h.SetBody(param) //设置get请求参数，post请求不要使用
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Post() //发送get请求
	if err != nil {
		return err
	}
	var tmpResp flowResp
	err = json.Unmarshal(response, &tmpResp)
	if err != nil {
		return err
	}
	if tmpResp.Code != 200 && len(tmpResp.Msg) > 0 {
		return errors.New(tmpResp.Msg)
	}
	return nil
}

type getNodeInfoResp struct {
	Code int                 `json:"code"`
	Data GetNodeInfoRespData `json:"data"`
	Msg  string              `json:"msg"`
}

type GetNodeInfoRespData struct {
	NodeId       string `json:"nodeId"`
	NodeType     string `json:"nodeType"`
	MainIpAddr   string `json:"mainIpAddr"`
	FlowTaskID   string `json:"flowTaskID"`
	FlowServerID string `json:"flowServerID"`
	NodeStatus   string `json:"nodeStatus"`
	CreateTime   string `json:"createTime"`
}

type downloadMITMCertResp struct {
	Code int `json:"code"`
	Data struct {
		Cert string `json:"cert"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 获取节点信息
func GetNodeInfo(param map[string]interface{}) (*GetNodeInfoRespData, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/flow/getnodeinfo")
	if err != nil {
		return &GetNodeInfoRespData{}, err
	}
	// 绑定参数
	h.GetUrlBuild(param) //设置get请求参数，post请求不要使用
	// 发起请求
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Get() //发送get请求
	if err != nil {
		return &GetNodeInfoRespData{}, err
	}
	// 解析响应
	var res getNodeInfoResp
	err = json.Unmarshal(response, &res)
	if err != nil {
		return &GetNodeInfoRespData{}, err
	}
	if res.Code != 200 {
		return &GetNodeInfoRespData{}, errors.New(res.Msg)
	}
	return &res.Data, nil
}

// 证书下载
func DownloadMITMCert() (string, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/flow/httpscert")
	if err != nil {
		return "", err
	}
	// 发起请求
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Get() //发送get请求
	if err != nil {
		return "", err
	}
	var res downloadMITMCertResp
	err = json.Unmarshal(response, &res)
	if err != nil {
		return "", err
	}
	if res.Code != 200 {
		return "", errors.New(res.Msg)
	}
	return res.Data.Cert, nil
}

type GetVulNameByTypeResp struct {
	Code int                  `json:"code"`
	Data GetVulNameByTypeData `json:"data"`
	Msg  string               `json:"msg"`
}

type GetVulNameByTypeData struct {
	Total int `json:"total"`
	List  []struct {
		VulName string `json:"vulName"`
		VulId   string `json:"vulId"`
		Pocname string `json:"pocname"`
	} `json:"list"`
}

// GetVulNameByType 通过漏洞类型获取漏洞名称
func GetVulNameByType(param map[string]interface{}) (GetVulNameByTypeData, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openvulscriptlistbyscripttype")
	if err != nil {
		return GetVulNameByTypeData{}, err
	}
	// 绑定参数
	h.GetUrlBuild(param) //设置get请求参数，post请求不要使用
	// 发起请求
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Get() //发送get请求
	if err != nil {
		return GetVulNameByTypeData{}, err
	}
	// 解析响应
	var res GetVulNameByTypeResp
	err = json.Unmarshal(response, &res)
	if err != nil {
		return GetVulNameByTypeData{}, err
	}
	if res.Code != 200 {
		return GetVulNameByTypeData{}, errors.New(res.Msg)
	}
	return res.Data, nil
}
