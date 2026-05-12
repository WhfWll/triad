package httpclients

import (
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
	"smart/api/typespec"
)

/***************************** 字典库 - 通过类型获取所有服务枚举 ********************************/
type OpenDictionaryGetServiceEnumByTypeReq struct {
	Type int `form:"type" json:"type" binding:"required"`
}
type OpenDictionaryGetServiceEnumByTypeRes struct {
	Code int `json:"code"`
	Data struct {
		List []typespec.GlobalOptionsItemRes `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 字典库 - 通过类型获取所有服务枚举
func OpenDictionaryGetServiceEnumByType(req OpenDictionaryGetServiceEnumByTypeReq) (res OpenDictionaryGetServiceEnumByTypeRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opengetserviceenumbytype") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 字典库 - 通过类型获取字典 ********************************/
type OpenDictionaryGetByTypeReq struct {
	Type      []int `form:"type" json:"type"`
	IsDefault int   `form:"isDefault" json:"isDefault"`
}
type OpenDictionaryGetByTypeRes struct {
	Code int `json:"code"`
	Data struct {
		List []OpenDictionaryGetByTypeItemRes `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type OpenDictionaryGetByTypeItemRes struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Service int    `json:"service"`
}

// 字典库 - 依据类型获取每个服务下（注意：默认）的配置 types = 类型 and is_default=默认 group by services
func GetDecisionDictionaryByType(req OpenDictionaryGetByTypeReq) (res OpenDictionaryGetByTypeRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionarybytype") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 字典库 - 通过类型与服务获取字典 ********************************/
type OpenDictionaryGetByTypeAndServiceReq struct {
	Type    []int `form:"type" json:"type"`
	Service int   `form:"service" json:"service"`
}
type OpenDictionaryGetByTypeAndServiceRes struct {
	Code int `json:"code"`
	Data struct {
		List []OpenDictionaryGetByTypeAndServiceItemRes `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type OpenDictionaryGetByTypeAndServiceItemRes struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Service   int    `json:"service"`
	Types     int    `json:"types"`
	IsDefault int    `json:"isDefault"`
}

// 字典库 - 依据类型与服务获取所有数据
func GetDecisionDictionaryByTypeAndService(req OpenDictionaryGetByTypeAndServiceReq) (res OpenDictionaryGetByTypeAndServiceRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionarybytypeandservice") //第一个参数是配置文件中client中的key，第二个参数uri
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

type OpenDictionaryEnumsRes struct {
	Code int `json:"code"`
	Data struct {
		Types   interface{} `json:"types"`
		Service struct {
			WeakPass      interface{} `json:"weakPass"`
			Wifi          interface{} `json:"wifi"`
			WebPathScan   interface{} `json:"webPathScan"`
			SubdomainScan interface{} `json:"subdomainScan"`
		} `json:"service"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// OpenDictionaryEnums 字典库 - 枚举
func OpenDictionaryEnums() (res OpenDictionaryEnumsRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionaryenums") //第一个参数是配置文件中client中的key，第二个参数uri

	//设置请求头
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	//发送请求
	response, err := h.Get()
	if err != nil {
		return
	}
	//解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

type OpenDictionaryListRes struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			Id            int    `json:"id"`
			Sources       int    `json:"sources"`
			Types         int    `json:"types"`
			TypesName     string `json:"typesName"`
			Service       int    `json:"service"`
			ServiceName   string `json:"serviceName"`
			IsDefault     int    `json:"isDefault"`
			IsDefaultName string `json:"isDefaultName"`
			Name          string `json:"name"`
			CreateTime    string `json:"createTime"`
			UpdateTime    string `json:"updateTime"`
		} `json:"list"`
		Total int64 `json:"total"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// OpenDictionaryList 字典库 - 列表页
func OpenDictionaryList(param map[string]interface{}) (res OpenDictionaryListRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionarylist") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	//绑定参数
	h.GetUrlBuild(param) //设置get请求参数，post请求不要使用
	//设置请求头
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	//发送请求
	response, err := h.Get()
	if err != nil {
		return
	}
	//解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

type OpenDictionaryDetailRes struct {
	Code int `json:"code"`
	Data struct {
		Id         int    `json:"id"`
		Sources    int    `json:"sources"`
		Types      int    `json:"types"`
		Service    int    `json:"service"`
		IsDefault  int    `json:"isDefault"`
		Name       string `json:"name"`
		Content    string `json:"content"`
		CreateTime string `json:"createTime"`
		UpdateTime string `json:"updateTime"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// OpenDictionaryDetail 字典库 - 详情页
func OpenDictionaryDetail(param map[string]interface{}) (res OpenDictionaryDetailRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionarydetail") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	//绑定参数
	h.GetUrlBuild(param)
	//设置请求头
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	//发送请求
	response, err := h.Get()
	if err != nil {
		return
	}
	//解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

type OpenDictionaryDeleteRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// OpenDictionaryDelete 字典库 - 删除
func OpenDictionaryDelete(param map[string]interface{}) (res OpenDictionaryDeleteRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionarydelete") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	//绑定参数
	h.GetUrlBuild(param)
	//设置请求头
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	//发送请求
	response, err := h.Get()
	if err != nil {
		return
	}
	//解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

type OpenDictionarySetDefaultRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// OpenDictionarySetDefault 字典库 - 设置默认字典
func OpenDictionarySetDefault(param map[string]interface{}) (res OpenDictionarySetDefaultRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionarysetdefault") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	//绑定参数
	h.GetUrlBuild(param)
	//设置请求头
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	//发送请求
	response, err := h.Get()
	if err != nil {
		return
	}
	//解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

type OpenDictionaryAddOrEditRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// OpenDictionaryAddOrEdit 字典库- 新增或编辑
func OpenDictionaryAddOrEdit(param map[string]interface{}) (res OpenDictionaryAddOrEditRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionaryaddoredit") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	//绑定参数
	h.SetBody(param)
	//设置请求头
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	//发送请求
	response, err := h.Post()
	if err != nil {
		return
	}
	//解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

type OpenDictionaryFieldsByIdsRes struct {
	Code int `json:"code"`
	Data struct {
		List []OpenDictionaryFieldsByIdsResItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type OpenDictionaryFieldsByIdsResItem struct {
	DictName string `json:"dictName"`
	DictType int    `json:"dictType"`
}

// OpenDictionaryFieldsByIds 根据id列表获取字典名称列表
func OpenDictionaryFieldsByIds(param map[string]interface{}) (res OpenDictionaryFieldsByIdsRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/dictionary/opendictionaryfields") //第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	//绑定参数
	h.GetUrlBuild(param)
	//设置请求头
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	//发送请求
	response, err := h.Get()
	if err != nil {
		return
	}
	//解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}
