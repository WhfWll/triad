package httpclients

import (
	"context"
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
)

// 决策引擎 - 节点管理

/***************************** 节点管理 - 节点添加 ********************************/
type OpenYakNodeAddReq struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
	Name string `json:"name"`
}
type OpenYakNodeAddRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// 节点管理 - 节点添加
func AddDecisionYakNode(ctx context.Context, req OpenYakNodeAddReq) (res OpenYakNodeAddRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/system/nodeadd") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 节点管理 - 节点列表 ********************************/
type OpenYakNodeListReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
type OpenYakNodeListRes struct {
	Code int `json:"code"`
	Data struct {
		List  []OpenYakNodeListItem `json:"list"`
		Total int64                 `json:"total"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type OpenYakNodeListItem struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	RunningNum    int    `json:"runningNum"` // 运行任务数
	Status        int    `json:"status"`
	StatusEnum    string `json:"statusEnum"`
	IsDisable     int    `json:"IsDisable"` // 禁用状态 0启用 1禁用
	IsDisableEnum string `json:"isDisableEnum"`
	CreateTime    string `json:"createTime"`
	UpdateTime    string `json:"updateTime"`
}

// 节点管理 - 节点列表
func GetDecisionYakNodeList(ctx context.Context, req OpenYakNodeListReq) (res OpenYakNodeListRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/system/nodelist") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 节点管理 - 节点删除 ********************************/
type OpenYakNodeDelReq struct {
	Id string `json:"id"`
}
type OpenYakNodeDelRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// 节点管理 - 节点删除
func DelDecisionYakNode(ctx context.Context, req OpenYakNodeDelReq) (res OpenYakNodeDelRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/system/nodedel") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 节点管理 - 设置节点是否启用 ********************************/
type OpenYakNodeSetIsDisableReq struct {
	Id        int `json:"id"`
	IsDisable int `json:"isDisable"`
}
type OpenYakNodeSetIsDisableRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// 节点管理 - 设置节点是否启用
func EditDecisionYakNodeIsDisable(ctx context.Context, req OpenYakNodeSetIsDisableReq) (res OpenYakNodeSetIsDisableRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/system/nodedisorenable") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 节点管理 - 获取所有可用节点 ********************************/
type OpenYakNodeAllEnableReq struct {
}
type OpenYakNodeAllEnableRes struct {
	Code int `json:"code"`
	Data struct {
		List []OpenYakNodeAllEnableItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type OpenYakNodeAllEnableItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// 节点管理 - 获取所有可用节点
func GetDecisionYakNodeAllEnable(ctx context.Context, req OpenYakNodeAllEnableReq) (res OpenYakNodeAllEnableRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/system/nodeallenable") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 节点管理 - 设置是否启用分布式 ********************************/
type OpenYakNodeSetDistributeReq struct {
	Status int `json:"status"` // 0禁用 1启用
}
type OpenYakNodeSetDistributeRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// 节点管理 - 设置是否启用分布式
func EditDecisionYakNodeSetDistribute(ctx context.Context, req OpenYakNodeSetDistributeReq) (res OpenYakNodeSetDistributeRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/system/nodesetdistribute") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 节点管理 - 获取是否启用分布式 ********************************/
type OpenYakNodeGetDistributeReq struct {
}
type OpenYakNodeGetDistributeRes struct {
	Code int `json:"code"`
	Data struct {
		Status int `json:"status"` // 0禁用 1启用
	} `json:"data"`
	Msg string `json:"msg"`
}

// 节点管理 - 获取是否启用分布式
func GetDecisionYakNodeGetDistribute(ctx context.Context, req OpenYakNodeGetDistributeReq) (res OpenYakNodeGetDistributeRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/system/nodegetdistribute") //第一个参数是配置文件中client中的key，第二个参数uri
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
