package httpclients

import (
	"context"
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
	"smart/api/typespec"
	"smart/tools/enums"
)

/***************************** 漏洞 - 公共参数 ********************************/
type OpenVulScriptResScriptItem struct {
	Id               uint   `json:"id"`               // id
	DataType         uint   `json:"dataType"`         // 1真数据 2假数据
	Pocname          string `json:"pocname"`          // pocname
	Type             string `json:"type"`             // 工具类型 yak | python | mitm | ...
	VerifyType       string `json:"verifyType"`       // 校验类型 poc | exp | ...
	EvidenceType     uint   `json:"evidenceType"`     // 取证类型
	EvidenceTypeEnum string `json:"evidenceTypeEnum"` // 取证类型
}
type OpenVulScriptResLibItem struct {
	Id                int    `json:"id"`                // id
	DataType          int    `json:"dataType"`          // 1真数据 2假数据
	Name              string `json:"name"`              // 漏洞名称
	Risk              int    `json:"risk"`              // 漏洞风险
	RiskEnum          string `json:"riskEnum"`          // 漏洞风险
	Type              int    `json:"type"`              // 漏洞类型
	TypeEnum          string `json:"typeEnum"`          // 漏洞类型
	Class             int    `json:"class"`             // 漏洞分类
	ClassEnum         string `json:"classEnum"`         // 漏洞分类
	PublishedTime     string `json:"publishedTime"`     // 公开时间
	Description       string `json:"description"`       // 漏洞描述
	AffectRange       string `json:"affectRange"`       // 影响范围
	ExploitImpact     int    `json:"exploitImpact"`     // 利用影响
	ExploitImpactEnum string `json:"exploitImpactEnum"` // 利用影响
	FixSuggest        string `json:"fixSuggest"`        // 修复建议
	Cve               string `json:"cve"`               // cve编号
	Cnvd              string `json:"cnvd"`              // cnvd编号
	Cnnvd             string `json:"cnnvd"`             // cnnvd编号
	Component         string `json:"component"`         // 受影响的组建名称
	Status            int    `json:"status"`
	StatusEnum        string `json:"statusEnum"`
	VulId             string `json:"vulId"`
}

/***************************** 漏洞 - 列表 ********************************/
// 决策引擎脚本信息获取
type OpenVulScriptListReq struct {
	Page              uint     `form:"page" json:"page"` // 必填
	Size              uint     `form:"size" json:"size"` // 必填
	LibIds            []int    `form:"libIds" json:"libIds"`
	ScriptVerifyTypes []string `form:"scriptVerifyTypes" json:"scriptVerifyTypes"`
	LibClasses        []int    `form:"libClasses" json:"libClasses"`
	LibTypes          []int    `form:"libTypes" json:"libTypes"`
	LibRisks          []int    `form:"libRisks" json:"libRisks"`
	OperatingSystem   []int    `form:"operatingSystem" json:"operatingSystem"` //操作系统 1:windows, 2:linux, 3:unix, 4:其他
	LibName           string   `form:"libName" json:"libName"`
	Status            []int    `form:"status" json:"status"`
	ExploitImpact     []int    `form:"exploitImpact" json:"exploitImpact"`
	PtOrder           string   `form:"ptOrder" json:"ptOrder"`
}

type OpenVulScriptListRes struct {
	Code int `json:"code"`
	Data struct {
		Total  int64                      `json:"total"`
		LibIds []int                      `json:"libIds"` // 当前满足查询条件的所有lib ID
		List   []OpenVulScriptListResItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type OpenVulScriptListResItem struct {
	Script    []OpenVulScriptResScriptItem `json:"script"`
	Libraries OpenVulScriptResLibItem      `json:"libraries"`
}

/***************************** 漏洞 - 详情 ********************************/
type OpenVulScriptDetailReq struct {
	LibId int `form:"libId" json:"libId"`
}
type OpenVulScriptDetailRes struct {
	Code int `json:"code"`
	Data struct {
		Script    []OpenVulScriptResScriptItem `json:"script"`
		Libraries OpenVulScriptResLibItem      `json:"libraries"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞 - 详情
func GetDecisionScriptDetail(req OpenVulScriptDetailReq) (res OpenVulScriptDetailRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openvulscriptdetail") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 漏洞 - 详情 通过pocname获取 ********************************/
type OpenVulScriptDetailByPocnameReq struct {
	Pocname string `form:"pocname" json:"pocname"`
}
type OpenVulScriptDetailByPocnameRes struct {
	Code int `json:"code"`
	Data struct {
		Script    []OpenVulScriptResScriptItem `json:"script"`
		Libraries OpenVulScriptResLibItem      `json:"libraries"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞 - 详情 通过pocname获取
func GetDecisionScriptDetailByPocname(req OpenVulScriptDetailByPocnameReq) (res OpenVulScriptDetailByPocnameRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openvulscriptdetailbypocname") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 漏洞 - 通过pocnames获取lib表ID 用于初始化模板 ********************************/
type OpenVulLibIdByPocnamesReq struct {
	MatchType string `form:"matchType" json:"matchType" binding:"required"` // 匹配类型 eq完全匹配 like模糊匹配
	Pocnames  string `form:"pocnames" json:"pocnames" binding:"required"`   // 多个逗号分割
}
type OpenVulLibIdByPocnamesRes struct {
	Code int `json:"code"`
	Data struct {
		LibIds []int `json:"libIds"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞 - 通过pocnames获取lib表ID 用于初始化模板
func OpenVulLibIdByPocnames(req OpenVulLibIdByPocnamesReq) (res OpenVulLibIdByPocnamesRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openvullibidbypocnames") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 漏洞 - 通过脚本类型获取lib表ID 用于初始化模板 ********************************/
type OpenVulLibIdByScriptReq struct {
	ScriptType string `form:"scriptType" json:"scriptType" binding:"required"` // 多个逗号分割
}
type OpenVulLibIdByScriptRes struct {
	Code int `json:"code"`
	Data struct {
		LibIds []int `json:"libIds"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞 - 通过脚本类型获取lib表ID 用于初始化模板
func OpenVulLibIdByScript(req OpenVulLibIdByScriptReq) (res OpenVulLibIdByScriptRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openvullibidbyscripttype") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 漏洞 - 修改 ********************************/
type OpenVulScriptEditReq struct {
	LibId         uint   `form:"libId" json:"libId" binding:"required"`
	Risk          uint   `form:"risk" json:"risk"`
	ExploitImpact uint   `form:"exploitImpact" json:"exploitImpact"`
	Description   string `form:"description" json:"description"`
	FixSuggest    string `form:"fixSuggest" json:"fixSuggest"`
	AffectRange   string `form:"affectRange" json:"affectRange"`
}
type OpenVulScriptEditRes struct {
	Code int `json:"code"`
	Data struct {
	} `json:"data"`
	Msg string `json:"msg"`
}

func EditDecisionScript(ctx context.Context, req OpenVulScriptEditReq) (res OpenVulScriptEditRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openvulscriptedit") //第一个参数是配置文件中client中的key，第二个参数uri
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
	h.SetSendType("from")
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

/***************************** 漏洞 - 孤立节点 ********************************/
type OpenVulScriptOrphanReq struct {
	LibIds string `form:"libIds" json:"libIds"`
}
type OpenVulScriptOrphanRes struct {
	Code int `json:"code"`
	Data []struct {
		LibId   uint   `json:"libId"`
		Pocname string `json:"pocname"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 获取参数中漏洞ID 有哪些是孤立节点
func GetDecisionOrphanScript(ctx context.Context, req OpenVulScriptOrphanReq) (res OpenVulScriptOrphanRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openvulscriptorphan") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 漏洞 - 枚举值 ********************************/
type OpenVulScriptOptions struct {
	Code int `json:"code"`
	Data struct {
		VulScriptType             []typespec.GlobalOptionsItemRes `json:"vulScriptType"`
		VulScriptVerifyType       []typespec.GlobalOptionsItemRes `json:"vulScriptVerifyType"`
		VulScriptEvidenceType     []typespec.GlobalOptionsItemRes `json:"vulScriptEvidenceType"`
		VulScriptExploitImpact    []typespec.GlobalOptionsItemRes `json:"vulScriptExploitImpact"`
		VulLibrariesClass         []typespec.GlobalOptionsItemRes `json:"vulLibrariesClass"`
		VulLibrariesType          []typespec.GlobalOptionsItemRes `json:"vulLibrariesType"`
		VulLibrariesRisk          []typespec.GlobalOptionsItemRes `json:"vulLibrariesRisk"`
		VulLibrariesOperateSystem []typespec.GlobalOptionsItemRes `json:"vulLibrariesOperateSystem"`
		VulLibrariesPriority      []typespec.GlobalOptionsItemRes `json:"vulLibrariesPriority"`
		VulLibrariesStatus        []typespec.GlobalOptionsItemRes `json:"vulLibrariesStatus"`
		OperateSystem             []typespec.GlobalOptionsItemRes `json:"operateSystem"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 获取参数中漏洞ID 有哪些是孤立节点
func GetDecisionOptions(ctx context.Context) (res OpenVulScriptOptions, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openoptions") //第一个参数是配置文件中client中的key，第二个参数uri
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

/***************************** 漏洞 - 修改状态 ********************************/
type OpenVulEditStatusReq struct {
	LibId  int `form:"libId" json:"libId" binding:"required"`
	Status int `form:"status" json:"status" binding:"required"`
}
type OpenVulEditStatusRes struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// 获取参数中漏洞ID 有哪些是孤立节点
func EditDecisionVulStatus(ctx context.Context, req OpenVulEditStatusReq) (res OpenVulEditStatusRes, err error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openvuleditstatus") //第一个参数是配置文件中client中的key，第二个参数uri
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
	h.SetSendType("from")
	h.SetBody(param) //设置get请求参数，post请求不要使用

	// 发起请求
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Post() //发送请求
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

type DecisionTaskVulVerifyResp struct {
	Code int `json:"code"`
	Data struct {
		Result  int    `json:"result"`
		Content string `json:"content"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// DecisionTaskVulVerify 任务漏洞信息验证接口
func DecisionTaskVulVerify(pocname string, vulParam []map[string]interface{}) (int, string, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/scriptcallnow")
	if err != nil {
		return 0, "", err
	}
	var param = map[string]interface{}{
		"param":    vulParam,
		"toolName": pocname,
	}
	h.SetBody(param)
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	res, err := h.Post()
	if err != nil {
		return 0, "", err
	}
	// 解析响应
	var tmp DecisionTaskVulVerifyResp
	err = json.Unmarshal(res, &tmp)
	if err != nil {
		return 0, "", err
	}
	if tmp.Data.Result == 2 {
		return enums.VulStatusRepairSuccess, tmp.Data.Content, nil //已经修复了
	}
	return 0, tmp.Data.Content, err
}

type DecisionScriptGraphResp struct {
	Code int `json:"code"`
	Data struct {
		Nodes interface{} `json:"nodes"`
		Links interface{} `json:"links"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// DecisionScriptGraph 知识图谱查询接口
func DecisionScriptGraph(ctx context.Context, vulIds []int) (DecisionScriptGraphResp, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/graphdisplay")
	var tmp DecisionScriptGraphResp
	if err != nil {
		return tmp, err
	}
	var param = map[string]interface{}{
		"vulIds": vulIds,
	}
	h.SetBody(param)
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	res, err := h.Post()
	if err != nil {
		return tmp, err
	}
	// 解析响应
	err = json.Unmarshal(res, &tmp)
	if err != nil {
		return tmp, err
	}

	return tmp, err
}

type DecisionScriptCallResp struct {
	Code int `json:"code"`
	Data struct {
		Result string `json:"result"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// DecisionScriptCall 脚本调用
func DecisionScriptCall(pocname string, vulParam []map[string]interface{}) (string, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/scriptcall")
	if err != nil {
		return "", err
	}
	var param = map[string]interface{}{
		"param":    vulParam,
		"toolName": pocname,
	}
	h.SetBody(param)
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	res, err := h.Post()
	if err != nil {
		return "", err
	}
	// 解析响应
	var tmp DecisionScriptCallResp
	err = json.Unmarshal(res, &tmp)
	if err != nil {
		return "", err
	}
	return tmp.Data.Result, err
}

type DecisionScriptCallResultResp struct {
	Code int `json:"code"`
	Data struct {
		Result string `json:"result"`
		End    bool   `json:"end"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// DecisionScriptCallResult 脚本调用获取结果
func DecisionScriptCallResult(callId string, isClean bool) (bool, string, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/scriptcallresult")
	if err != nil {
		return false, "", err
	}
	h.GetUrlBuild(map[string]interface{}{
		"callId":  callId,
		"isClean": isClean,
	})
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	res, err := h.Get()
	if err != nil {
		return true, "", err
	}
	// 解析响应
	var tmp DecisionScriptCallResultResp
	err = json.Unmarshal(res, &tmp)
	if err != nil {
		return true, "", err
	}
	return tmp.Data.End, tmp.Data.Result, err
}

type DecisionReportVerifyTaskResp struct {
	Code int `json:"code"`
	Data struct {
		ObjId         int                 `json:"objId"`
		SubObjId      int                 `json:"subObjId"`
		CanVerifyList []map[string]string `json:"canVerifyList"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// DecisionReportVerifyTask 发送报告验证任务
func DecisionReportVerifyTask(objId, subObjId int, cveList []string, paramMap []map[string]string) ([]map[string]string, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openreportverifytask")
	if err != nil {
		return nil, err
	}
	var param = map[string]interface{}{
		"cveList":  cveList,
		"param":    paramMap,
		"objId":    objId,
		"subObjId": subObjId,
	}
	h.SetBody(param)
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	res, err := h.Post()
	if err != nil {
		return nil, err
	}
	// 解析响应
	var tmp DecisionReportVerifyTaskResp
	err = json.Unmarshal(res, &tmp)
	if err != nil {
		return nil, err
	}
	return tmp.Data.CanVerifyList, nil
}

type DecisionReportVerifyTaskResultResp struct {
	Code int `json:"code"`
	Data struct {
		End       bool                `json:"end"`
		ResultMap map[string][]string `json:"resultMap"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// DecisionReportVerifyTaskResult 发送报告验证任务
func DecisionReportVerifyTaskResult(objId, subObjId int) (bool, map[string][]string, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/openreportverifytaskresult")
	if err != nil {
		return false, nil, err
	}
	h.GetUrlBuild(map[string]interface{}{
		"objId":    objId,
		"subObjId": subObjId,
	})
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	res, err := h.Get()
	if err != nil {
		return false, nil, err
	}
	// 解析响应
	var tmp DecisionReportVerifyTaskResultResp
	err = json.Unmarshal(res, &tmp)
	if err != nil {
		return false, nil, err
	}
	return tmp.Data.End, tmp.Data.ResultMap, nil
}

// DecisionScriptStop 脚本结束
func DecisionScriptStop(callId string) error {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/vulscripts/scriptcallstop")
	if err != nil {
		return err
	}
	h.GetUrlBuild(map[string]interface{}{
		"callId": callId,
	})
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	res, err := h.Get()
	// 解析响应
	var tmp DecisionScriptCallResp
	err = json.Unmarshal(res, &tmp)
	if err != nil {
		return err
	}
	return err
}
