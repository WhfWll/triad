package httpclients

import (
	"context"
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
)

// 决策引擎 - Ip域名绑定

/***************************** Ip域名绑定 - 创建 ********************************/
type OpenIpHostBindCreateReq struct {
	Id    int    `form:"id" json:"id"` // 编辑时必传
	Ip    string `form:"ip" json:"ip" binding:"required"`
	Hosts string `form:"hosts" json:"hosts" binding:"required"`
}
type OpenIpHostBindCreateRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// ip或域名绑定 - 列表
func PostDecisionIpHostBindCreate(ctx context.Context, req OpenIpHostBindCreateReq) (res OpenIpHostBindCreateRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/iphostbind/opencreate") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	// 参数获取
	param := make(map[string]interface{})
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}
	err = json.Unmarshal(reqByte, &param)
	if err != nil {
		return
	}

	// 绑定参数
	h.SetBody(param) //设置get请求参数，post请求不要使用

	// 发起请求
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Post() //发送get请求
	if err != nil {
		return
	}

	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

/***************************** Ip域名绑定 - 列表 ********************************/
type OpenIpHostBindListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type OpenIpHostBindListRes struct {
	Code int `json:"code"`
	Data struct {
		Page  int                      `json:"page"`
		Size  int                      `json:"size"`
		Total int64                    `json:"total"`
		List  []OpenIpHostBindListItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type OpenIpHostBindListItem struct {
	Id    int    `json:"id"`
	Ip    string `json:"ip"`
	Hosts string `json:"hosts"`
}

// ip或域名绑定 - 列表
func GetDecisionIpHostBindList(ctx context.Context, req OpenIpHostBindListReq) (res OpenIpHostBindListRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/iphostbind/openlist") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	// 参数获取
	param := make(map[string]interface{})
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}
	err = json.Unmarshal(reqByte, &param)
	if err != nil {
		return
	}

	// 绑定参数
	h.GetUrlBuild(param) //设置get请求参数，post请求不要使用

	// 发起请求
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Get() //发送get请求
	if err != nil {
		return
	}

	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

/***************************** Ip域名绑定 - 删除 ********************************/
type OpenIpHostBindDelReq struct {
	Id string `form:"id" json:"id"`
}
type OpenIpHostBindDelRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// ip或域名绑定 - 列表
func PostDecisionIpHostBindDel(ctx context.Context, req OpenIpHostBindDelReq) (res OpenIpHostBindDelRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/iphostbind/opendel") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	// 参数获取
	param := make(map[string]interface{})
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}
	err = json.Unmarshal(reqByte, &param)
	if err != nil {
		return
	}

	// 绑定参数
	h.SetBody(param) //设置get请求参数，post请求不要使用

	// 发起请求
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Post() //发送get请求
	if err != nil {
		return
	}

	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}
