package httpclients

import (
	"context"
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
	"smart/api/typespec"
)

// 决策引擎 - 指纹库

/***************************** 指纹库 - 枚举 ********************************/
type OpenFingerEnumRes struct {
	Code int `json:"code"`
	Data struct {
		Class      []typespec.GlobalOptionsItemRes `json:"class"`      // 指纹分类
		IsDev      []typespec.GlobalOptionsItemRes `json:"isDev"`      // 设备指纹
		Level      []typespec.GlobalOptionsItemRes `json:"level"`      // 设备分层指纹
		SoftOrHard []typespec.GlobalOptionsItemRes `json:"softOrHard"` // 设备软硬件
	} `json:"data"`
	Msg string `json:"msg"`
}

// 指纹 - 枚举
func GetDecisionFingerEnum(ctx context.Context) (res OpenFingerEnumRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/finger/openenum") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

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

/***************************** 指纹库 - 列表 ********************************/
type OpenFingerListReq struct {
	Page       int    `form:"page" json:"page" binding:"required"`
	Size       int    `form:"size" json:"size" binding:"required"`
	Class      string `form:"class" json:"class"`
	FingerName string `form:"fingerName" json:"fingerName"`
}
type OpenFingerListRes struct {
	Code int `json:"code"`
	Data struct {
		Page  int                  `json:"page"`
		Size  int                  `json:"size"`
		Total int64                `json:"total"`
		List  []OpenFingerListItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type OpenFingerListItem struct {
	Id                  int    `json:"id"`
	Name                string `json:"fingerName"`
	Class               int    `json:"fingerClass"`
	ClassEnum           string `json:"fingerClassEnum"`
	Desc                string `json:"desc"`
	Num                 int    `json:"num"`
	Rule                string `json:"rule"`
	IconHash            string `json:"iconHash"`
	CertHash            string `json:"certHash"`
	Level               string `json:"level"`
	AppName             string `json:"appName"`
	AppClassEnum        string `json:"appClassEnum"`
	FingerSoftHardTypes string `json:"fingerSoftHardTypes"` // 指纹软硬
	IsDev               int    `json:"isDev"`
	CPE                 string `json:"cpe"`
	Version             string `json:"version"`
}

// 指纹 - 枚举
func GetDecisionFingerList(ctx context.Context, req OpenFingerListReq) (res OpenFingerListRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/finger/openlist") //第一个参数是配置文件中client中的key，第二个参数uri
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

// UpdateFingerResp 添加/修改指纹返回
type UpdateFingerResp struct {
	Code int `json:"code"`
	Data struct {
		ID int `form:"id" json:"id"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// UpdateFingerInfo 添加/修改指纹信息
func UpdateFingerInfo(ctx context.Context, req typespec.EditFingerReq) (res UpdateFingerResp, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/finger/save")
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
	h.SetBody(param)
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

// FingerTestResResp 指纹详情
type FingerTestResResp struct {
	Code int `json:"code"`
	Data struct {
		typespec.FingerCreateTestResp
	} `json:"data"`
	Msg string `json:"msg"`
}

// TestFingerInfo 测试指纹信息
func TestFingerInfo(ctx context.Context, req typespec.FingerCreateTestReq) (res FingerTestResResp, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/finger/createappointtest")
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
	h.SetBody(param)
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

// FingerTestResultResResp 指纹详情
type FingerTestResultResResp struct {
	Code int `json:"code"`
	Data struct {
		Result string `json:"result"` //调用结果
		End    bool   `json:"end"`    //是否结束判断依据
	} `json:"data"`
	Msg string `json:"msg"`
}

// TestFingerResult 测试指纹信息返回
func TestFingerResult(ctx context.Context, req typespec.FingerTestResultsReq) (res FingerTestResultResResp, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/scriptcallresult")
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
	h.GetUrlBuild(param)
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

// DelFingerInfo 删除指纹信息
func DelFingerInfo(ctx context.Context, id int) (res UpdateFingerResp, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/finger/delete")
	if err != nil {
		return
	}
	// 参数获取
	param := map[string]interface{}{
		"id": id,
	}
	h.SetBody(param)
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

// FingerDetailResp 指纹详情
type FingerDetailResp struct {
	Code int `json:"code"`
	Data struct {
		FingerInfo
	} `json:"data"`
	Msg string `json:"msg"`
}

type FingerInfo struct {
	ID                  int    `form:"id" json:"id"`
	AppClass            int    `form:"appClass" json:"appClass"`
	AppClassEnum        string `form:"appClassEnum" json:"appClassEnum"`
	AppVersion          string `form:"appVersion" json:"appVersion"`
	CnName              string `form:"cnName" json:"cnName"`
	AppName             string `form:"appName" json:"appName"`
	Flag                string `form:"flag" json:"flag"`
	IsDev               int    `form:"isDev" json:"isDev"`
	IsDevEnum           string `form:"isDevEnum" json:"isDevEnum"`
	Desc                string `form:"desc" json:"desc"`
	CreateTime          string `form:"createTime" json:"createTime"`
	FingerSoftHardTypes string `json:"fingerSoftHardTypes"` // 指纹软硬
	LevelID             string `json:"levelID"`             // 指纹分层ID
	Level               string `json:"level"`               // 指纹分层
}

// FingerDetail 指纹信息详情
func FingerDetail(ctx context.Context, id int) (res FingerDetailResp, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/finger/detail")
	if err != nil {
		return
	}
	// 参数获取
	param := map[string]interface{}{
		"id": id,
	}
	h.GetUrlBuild(param)
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
