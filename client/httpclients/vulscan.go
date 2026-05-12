package httpclients

import (
	"context"
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
	"time"
)

/***************************** 漏洞 - 列表 ********************************/
// 漏洞扫描任务列表
type VulScanTaskListReq struct {
	Page   int    `form:"page" json:"page"`     // 必填
	Size   int    `form:"size" json:"size"`     // 必填
	Risk   int    `form:"risk" json:"risk"`     // 必填
	Search string `form:"search" json:"search"` // 必填
}

type VulScanTaskListResp struct {
	Code int `json:"code"`
	Data struct {
		Total int                       `json:"total"`
		List  []VulScanTaskListRespItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type VulScanTaskListRespItem struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	TypeName   string `json:"typeName"`
	Risk       int    `json:"risk"`
	RiskName   string `json:"riskName"`
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
	High       int    `json:"high"`
	Middle     int    `json:"middle"`
	Low        int    `json:"low"`
	Safe       int    `json:"safe"`
	CreateTime string `json:"CreateTime"`
	UpdateTime string `json:"UpdateTime"`
}

// 漏洞扫描任务列表
func VulScanTaskList(ctx context.Context, req VulScanTaskListReq) (res VulScanTaskListResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/task/list") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描目标列表
type VulScanTargetListReq struct {
	Page    int    `form:"page" json:"page"`     // 必填
	Size    int    `form:"size" json:"size"`     // 必填
	TaskId  int    `form:"taskId" json:"taskId"` // 必填
	Risk    int    `form:"risk" json:"risk"`     // 必填
	Search  string `form:"search" json:"search"` // 必填
	IsAlive string `form:"isAlive" json:"isAlive"`
}

type VulScanTargetListResp struct {
	Code int `json:"code"`
	Data struct {
		Total int                         `json:"total"`
		List  []VulScanTargetListRespItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type VulScanTargetListRespItem struct {
	Id          int    `json:"id"`
	TaskId      int    `json:"task_id"`     // 任务id
	Ip          string `json:"ip"`          // ip地址
	Target      string `json:"target"`      // 测试目标
	System      string `json:"system"`      // 操作系统
	Port        string `json:"port"`        // 扫描端口
	Risk        int    `json:"risk"`        // 任务风险
	RiskName    string `json:"riskName"`    // 任务风险名称
	High        int    `json:"high"`        // 高危漏洞数
	Middle      int    `json:"middle"`      // 中危漏洞数
	Low         int    `json:"low"`         // 低危漏洞数
	IsAlive     int    `json:"isAlive"`     // 是否存活
	IsAliveName string `json:"isAliveName"` // 是否存活名称
	Status      int    `json:"status"`      // 目标状态
	StatusName  string `json:"statusName"`  // 任务风险名称
	CreateTime  string `json:"createTime"`  //创建时间
	UpdateTime  string `json:"updateTime"`  //更新时间
}

// 漏洞扫描任务列表
func VulScanTargetList(ctx context.Context, req VulScanTargetListReq) (res VulScanTargetListResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/target/list") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描目标列表
type VulScanTaskSaveReq struct {
	Name         string `json:"name" binding:"required"`
	Target       string `json:"target" binding:"required"`
	ToScanPort   string `json:"toScanPort" binding:"required"`
	OnlyPortScan bool   `json:"only_port_scan"` // 仅端口扫描模式
}

type VulScanTaskSaveResp struct {
	Code int `json:"code"`
	Data struct {
		Id int `json:"id"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞扫描任务保存
func VulScanTaskSave(ctx context.Context, req VulScanTaskSaveReq) (res VulScanTaskSaveResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/task/save") //第一个参数是配置文件中client中的key，第二个参数uri
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

	response, err := h.Post() //发送post请求
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

// 漏洞扫描任务停止
type VulScanTaskStopReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}

type VulScanTaskStopResp struct {
	Code int `json:"code"`
	Data struct {
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞扫描任务保存
func VulScanTaskStop(ctx context.Context, req VulScanTaskStopReq) (res VulScanTaskStopResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/task/stop") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描任务删除
type VulScanTaskDeleteReq struct {
	Ids string `form:"ids" json:"ids" binding:"required"`
}

type VulScanTaskDeleteResp struct {
	Code int `json:"code"`
	Data struct {
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞扫描任务保存
func VulScanTaskDelete(ctx context.Context, req VulScanTaskDeleteReq) (res VulScanTaskDeleteResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/task/delete") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描目标列表
type VulScanVulListReq struct {
	Page   int    `form:"page" json:"page"`     // 必填
	Size   int    `form:"size" json:"size"`     // 必填
	TaskId int    `form:"taskId" json:"taskId"` // 必填
	Risk   int    `form:"risk" json:"risk"`     // 必填
	Search string `form:"search" json:"search"` // 必填
}

type VulScanVulListResp struct {
	Code int `json:"code"`
	Data struct {
		Total int                      `json:"total"`
		List  []VulScanVulListRespItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type VulScanVulListRespItem struct {
	Id          int    `json:"id"`
	TaskID      int    `json:"taskID"`      // 任务id
	TargetID    int    `json:"targetID"`    // 目标id
	Name        string `json:"name"`        // 漏洞名称
	Port        string `json:"port"`        // 风险端口
	Risk        int    `json:"risk"`        // 风险等级
	RiskName    string `json:"riskName"`    // 风险等级名称
	CreateTime  string `json:"createTime"`  // 创建时间
	UpdateTime  string `json:"updateTime"`  // 更新时间
	Ip          string `json:"ip"`          // ip地址
	Cwe         string `json:"cwe"`         // cwe类型
	PublishDate string `json:"publishDate"` // 发布日期
	Cve         string `json:"cve"`         // cve编号
}

// 漏洞扫描漏洞列表
func VulScanVulList(ctx context.Context, req VulScanVulListReq) (res VulScanVulListResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/vul/list") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描漏洞详情
type VulScanVulDetailReq struct {
	Id int `form:"id" json:"id"` // 必填
}
type VulScanVulDetailResp struct {
	Code int `json:"code"`
	Data struct {
		ID                  int       `json:"id"`                   // 主键
		TaskID              int       `json:"taskID"`               // 任务id
		TargetID            int       `json:"targetID"`             // 目标id
		Name                string    `json:"name"`                 // 漏洞名称
		Cve                 string    `json:"cve"`                  // cwe类型
		Port                string    `json:"port"`                 // 风险端口
		Description         string    `json:"description"`          // 漏洞描述
		Solution            string    `json:"solution"`             // 解决方案
		Parameter           string    `json:"parameter"`            // 参数
		Detail              string    `json:"detail"`               // 漏洞详情
		Risk                int       `json:"risk"`                 // 风险等级
		RiskName            string    `json:"riskName"`             // 风险等级中文
		CreateTime          time.Time `json:"createTime"`           // 创建时间
		UpdateTime          time.Time `json:"updateTime"`           // 更新时间
		Ip                  string    `json:"ip"`                   // ip地址
		Cwe                 string    `json:"cwe"`                  // cwe类型
		Vendor              string    `json:"vendor"`               // 厂商
		Product             string    `json:"product"`              // 产品
		Version             string    `json:"version"`              // 产品
		Cpes                string    `json:"cpes"`                 //
		CvssVersion         string    `json:"cvss_version"`         // cvss版本
		CvssVector          string    `json:"cvss_vector"`          // cvss向量
		PublishDate         string    `json:"publish_date"`         // 发布时间
		ExploitabilityScore string    `json:"exploitability_score"` // 可利用评分
		References          string    `json:"references"`           // 参考链接
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞扫描漏洞详情
func VulScanVulDetail(ctx context.Context, req VulScanVulDetailReq) (res VulScanVulDetailResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/vul/detail") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描目标列表
type VulScanCveListReq struct {
	Page   int    `form:"page" json:"page"`     // 必填
	Size   int    `form:"size" json:"size"`     // 必填
	Search string `form:"search" json:"search"` // 必填
}

type VulScanCveListResp struct {
	Code int `json:"code"`
	Data struct {
		Total int                      `json:"total"`
		List  []VulScanCveListRespItem `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type VulScanCveListRespItem struct {
	Id                int     `json:"id"`
	CreatedAt         string  `json:"created_at"`          // 创建时间
	UpdatedAt         string  `json:"updated_at"`          // 更新时间
	DeletedAt         string  `json:"deleted_at"`          // 删除时间
	Cve               string  `json:"cve"`                 // cve编号
	Cwe               string  `json:"cwe"`                 // cwe类型
	ProblemType       string  `json:"problem_type"`        // 问题类型
	References        string  `json:"references"`          // 参考链接
	TitleZh           string  `json:"title_zh"`            // 中文名字
	Solution          string  `json:"solution"`            // 解决方案
	DescriptionMain   string  `json:"description_main"`    // 描述
	DescriptionMainZh string  `json:"description_main_zh"` // 描述
	Descriptions      string  `json:"descriptions"`        // 描述
	Vendor            string  `json:"vendor"`              // 厂商
	Product           string  `json:"product"`             // 厂商
	Severity          string  `json:"severity"`            // 基础严重性
	SeverityName      string  `json:"severityName"`        // 基础严重性
	PublishedDate     string  `json:"published_date"`      // 发布日期
	BaseCvssv2Score   float64 `json:"base_cvssv2_score"`   // 发布日期
}

// 漏洞扫描漏洞列表
func VulScanCveList(ctx context.Context, req VulScanCveListReq) (res VulScanCveListResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/cve/list") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描漏洞详情
type VulScanCveDetailReq struct {
	Id int `form:"id" json:"id"` // 必填
}

type VulScanCveDetailResp struct {
	Code int `json:"code"`
	Data struct {
		Id                      int     `json:"id"`
		CreatedAt               string  `json:"created_at"`                // 创建时间
		UpdatedAt               string  `json:"updated_at"`                // 更新时间
		DeletedAt               string  `json:"deleted_at"`                // 删除时间
		Cve                     string  `json:"cve"`                       // cve编号
		Cwe                     string  `json:"cwe"`                       // cwe类型
		ProblemType             string  `json:"problem_type"`              // 问题类型
		References              string  `json:"references"`                // 参考链接
		TitleZh                 string  `json:"title_zh"`                  // 中文名字
		Solution                string  `json:"solution"`                  // 解决方案
		DescriptionMain         string  `json:"description_main"`          // 描述
		DescriptionMainZh       string  `json:"description_main_zh"`       // 描述
		Descriptions            string  `json:"descriptions"`              // 描述
		Vendor                  string  `json:"vendor"`                    // 厂商
		Product                 string  `json:"product"`                   // 产品
		CpeConfigurations       string  `json:"cpe_configurations"`        // 通用平台枚举
		CvssVersion             string  `json:"cvss_version"`              // cvss版本
		CvssVectorString        string  `json:"cvss_vector_string"`        // 通用漏洞评分系统向量表示法
		AccessVector            string  `json:"access_vector"`             // 访问向量
		AccessComplexity        string  `json:"access_complexity"`         // 访问复杂性
		Authentication          string  `json:"authentication"`            // 认证信息
		ConfidentialityImpact   string  `json:"confidentiality_impact"`    // 机密性影响
		IntegrityImpact         string  `json:"integrity_impact"`          // 完整性影响
		AvailabilityImpact      string  `json:"availability_impact"`       // 可用性影响
		BaseCvssv2Score         float64 `json:"base_cvs_sv2_score"`        // 基础评分
		Severity                string  `json:"severity"`                  // 基础严重性
		SeverityName            string  `json:"severityName"`              // 基础严重性
		ExploitabilityScore     float64 `json:"exploitability_score"`      // 可利用性评分
		ImpactScore             float64 `json:"impact_score"`              // 影响评分
		ObtainAllPrivilege      string  `json:"obtain_all_privilege"`      // 获取所有权限
		ObtainUserPrivilege     string  `json:"obtain_user_privilege"`     // 获取用户权限
		ObtainOtherPrivilege    string  `json:"obtain_other_privilege"`    // 获取其他权限
		UserInteractionRequired string  `json:"user_interaction_required"` // 用户交互要求
		PublishedDate           string  `json:"published_date"`            // 发布日期
		LastModifiedDate        string  `json:"last_modified_date"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞扫描漏洞详情
func VulScanCveDetail(ctx context.Context, req VulScanCveDetailReq) (res VulScanCveDetailResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/cve/detail") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描漏洞详情
type VulScanTaskOverviewReq struct {
	Id int `form:"id" json:"id"` // 必填
}
type VulScanTaskOverviewResp struct {
	Code int `json:"code"`
	Data struct {
		TaskName   string `form:"taskName" json:"taskName" binding:"required"`
		Risk       int    `form:"risk" json:"risk" binding:"required"`
		RiskName   string `form:"riskName" json:"riskName" binding:"required"`
		TargetRisk []int  `form:"targetRisk" json:"targetRisk" binding:"required"`
		TargetNum  []int  `form:"targetNum" json:"targetNum" binding:"required"`
		VulRisk    []int  `form:"vulRisk" json:"vulRisk" binding:"required"`
		CreateTime string `form:"createTime" json:"createTime" binding:"required"`
		UpdateTime string `form:"updateTime" json:"updateTime" binding:"required"`
		Ports      string `form:"ports" json:"ports" binding:"required"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞扫描漏洞详情
func VulScanTaskOverview(ctx context.Context, req VulScanTaskOverviewReq) (res VulScanTaskOverviewResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/task/overview") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描漏洞详情
type VulScanTaskStateReq struct {
	Id int `form:"id" json:"id"` // 必填
}
type VulScanTaskStateResp struct {
	Code int `json:"code"`
	Data struct {
		Status     int    `form:"status" json:"status" binding:"required"`
		StatusName string `form:"statusName" json:"statusName" binding:"required"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞扫描漏洞状态
func VulScanTaskState(ctx context.Context, req VulScanTaskStateReq) (res VulScanTaskStateResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/task/getstate") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描获取报告内容
type VulScanReportContentReq struct {
	Id int `form:"id" json:"id"` // 必填
}
type VulScanReportContentResp struct {
	Code int `json:"code"`
	Data struct {
		Status  int    `json:"status"`  // 状态
		Content string `json:"content"` // 内容
	} `json:"data"`
	Msg string `json:"msg"`
}

// 漏洞扫描获取报告内容
func VulScanReportContent(ctx context.Context, req VulScanReportContentReq) (res VulScanReportContentResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/report/content") //第一个参数是配置文件中client中的key，第二个参数uri
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

// 漏洞扫描获取报告生成
type VulScanReportSaveReq struct {
	ReportId   int    `json:"reportId" form:"reportId" binding:"required"`
	Name       string `json:"name" form:"name" binding:"required"`
	Type       int    `json:"type" form:"type" binding:"required"`
	ConfigJson string `json:"configJson" form:"configJson" binding:"required"`
	Format     int    `json:"format" form:"format" binding:"required"`
	OutputType int    `json:"outputType" form:"outputType"` // 1合并输出 2逐个输出
	ObjIDName  string `json:"objIDName" form:"objIDName"`   // 逐个输出会用到 为了规避objid和目标对不上的问题
}
type VulScanReportSaveResp struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

// 漏洞扫描获取报告生成
func VulScanReportSave(ctx context.Context, req VulScanReportSaveReq) (res VulScanReportSaveResp, err error) {
	h, err := httpclient.NewHttpSend("service_vulscan", "/vulscan/report/save") //第一个参数是配置文件中client中的key，第二个参数uri
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
