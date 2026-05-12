package typespec

import (
	"smart/tools/enums"
	"smart/tools/utils"
)

// 任务中心 - 渗透任务

// 任务枚举
type TaskTaskEnumRes struct {
	ExecuteType   []GlobalOptionsItemRes `json:"executeType"`
	RiskLevel     []GlobalOptionsItemRes `json:"riskLevel"`
	Status        []GlobalOptionsItemRes `json:"status"`
	RuntimePeriod []GlobalOptionsItemRes `json:"runtimePeriod"`
	Weight        interface{}            `json:"weight"`
	TestIntensity interface{}            `json:"testIntensity"`

	// 周期计划执行时的枚举
	CyclePlanningType       []GlobalOptionsItemRes `json:"cyclePlanningType"`
	CyclePlanningWeekValue  []GlobalOptionsItemRes `json:"cyclePlanningWeekValue"`
	CyclePlanningMonthValue []GlobalOptionsItemRes `json:"cyclePlanningMonthValue"`
	// 网站登陆凭证
	WebLoginType []GlobalOptionsItemRes `json:"webLoginType"`
	//横向移动
	//攻击面

	AttackFaceType interface{} `json:"attackFaceType"`
	// 代理协议类型
	ProxyProto []GlobalOptionsItemRes `json:"proxyProto"`
}

// 网站登陆凭证校验
type TaskTaskWebLoginCheckReq struct {
	TaskCheckTarget string `form:"taskCheckTarget" json:"taskCheckTarget"`
	Target          string `form:"target" json:"target"`
	VerifyType      int    `form:"verifyType" json:"verifyType"`
	VerifyValue     string `form:"verifyValue" json:"verifyValue"`
}
type TaskTaskWebLoginCheckRes struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
}

// 创建
type TaskSaveReq struct {
	TaskTemplateId int              `form:"taskTemplateId" json:"taskTemplateId" binding:"required"`
	Target         string           `form:"target" json:"target" binding:"required"`
	TaskName       string           `form:"taskName" json:"taskName" binding:"required"`
	ExecuteType    int              `form:"executeType" json:"executeType" binding:"required"` // 1即时执行 2定时执行 3周期执行
	ExecuteJson    ExecuteJson      `form:"executeJson" json:"executeJson" binding:"required"`
	Config         enums.ConfigJson `form:"config" json:"config" binding:"required"`
	Weight         int              `form:"weight" json:"weight" binding:"required"`
	UserId         int              `form:"userId" json:"userId"`
	Pid            int              `form:"pid" json:"pid"`                     // 任务父级任务ID，用于任务验证测试功能，默认为0 验证测试任务时必须传值
	SyncReport     string           `form:"syncReport" json:"syncReport"`       // 同步报告
	ExcludeTarget  string           `form:"excludeTarget" json:"excludeTarget"` // 排除目标
}

// 执行方式说明
// 即时执行：以下参数无效
// 定时执行：以下参数仅StartTime有效
// 周期执行：以下参数除StartTime外都有效
type ExecuteJson struct {
	StartTime          string   `form:"startTime" json:"startTime"`                   // 开始时间
	CyclePlanningType  int      `form:"cyclePlanningType" json:"cyclePlanningType"`   // 周期计划类型 1周 2月
	CyclePlanningValue int      `form:"cyclePlanningValue" json:"cyclePlanningValue"` // 周期计划类型 具体周期 （其中周：1周一 2周二... 7周日） （其中月：1一号 2二号...）
	CyclePlanningHour  string   `form:"cyclePlanningHour" json:"cyclePlanningHour"`   // 周期计划类型 具体小时与分钟 00:01 ｜ 00:20等
	EndTime            string   `form:"endTime" json:"endTime"`                       // 结束时间
	RuntimePeriod      []string `form:"runtimePeriod" json:"runtimePeriod" binding:"required"`
}
type TaskSaveResp struct {
	TaskId int `json:"task_id"`
}

// 任务列表
type TaskListReq struct {
	Page      int    `form:"page" json:"page" binding:"required"`
	Size      int    `form:"size" json:"size" binding:"required"`
	RiskLevel int    `form:"riskLevel" json:"riskLevel"`
	TaskName  string `form:"taskName" json:"taskName"`
	StartTime string `form:"startTime" json:"startTime"`
	EndTime   string `form:"endTime" json:"endTime"`
}
type TaskListRes struct {
	Total int               `json:"total"`
	List  []TaskListItemRes `json:"list"`
}
type TaskListItemRes struct {
	Id              int    `json:"id"`
	TaskName        string `json:"taskName"`
	ExecuteType     int    `json:"executeType"`
	ExecuteTypeName string `json:"executeTypeName"`
	RiskLevel       int    `json:"riskLevel"`
	RiskLevelName   string `json:"riskLevelName"`
	Status          int    `json:"status"`
	StatusName      string `json:"statusName"`
	TargetRisk      []int  `json:"targetRisk"` // 任务风险等级 下标区分 0高危 1中危 2低危 3安全
	CreateTime      string `json:"createTime"`
	UpdateTime      string `json:"updateTime"`
	Progress        int    `json:"progress"` // A6项目需要进度条，虽然并不准确
}

// 任务删除
type TaskDelReq struct {
	TaskIds string `form:"taskIds" json:"taskIds" binding:"required"`
}
type TaskDelRes struct {
}

// 任务复制
type TaskCopyReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}
type TaskCopyRes struct {
	TaskSaveReq // 使用创建任务的请求 作为响应数据
}

// 任务状态更改 暂停 / 恢复 / 结束
type TaskChangeStateReq struct {
	TaskId  int    `form:"taskId" json:"taskId"`
	Operate string `form:"operate" json:"operate"`
}
type TaskChangeStateRes struct {
}

type GetStateReq struct {
	TaskId int `json:"taskId" form:"taskId" binding:"required"`
}

type GetStateResp struct {
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
}

type OverViewReq struct {
	TaskId int `json:"taskId" form:"taskId" binding:"required"`
}

type OverViewResp struct {
	Id                 int         `json:"id"`
	TaskName           string      `json:"taskName"`
	ExecuteType        int         `json:"executeType"`
	ExecuteTypeName    string      `json:"executeTypeName"`
	RuntimePeriod      string      `json:"runtimePeriod"`
	StartTime          string      `json:"startTime"`
	CyclePlanningType  int         `json:"cyclePlanningType"`
	CyclePlanningValue int         `json:"cyclePlanningValue"`
	CyclePlanningHour  string      `json:"cyclePlanningHour"`
	EndTime            string      `json:"endTime"`
	CyclePlanningName  string      `json:"cyclePlanningName"`
	EndTimeName        string      `json:"endTimeName"`
	RiskLevel          int         `json:"riskLevel"`
	RiskLevelName      string      `json:"riskLevelName"`
	TaskTemplateId     int         `json:"taskTemplateId"`
	TaskTemplateName   string      `json:"taskTemplateName"`
	Status             int         `json:"status"`
	StatusName         string      `json:"statusName"`
	CreateTime         string      `json:"createTime"`
	UpdateTime         string      `json:"updateTime"`
	TargetRisk         [4]int      `json:"targetRisk"`
	TargetNum          []int       `json:"targetNum"`
	VulTotal           int         `json:"vulTotal"`
	VulRisk            []int       `json:"vulRisk"`
	VulExploitImpact   interface{} `json:"vulExploitImpact"`
	EvidenceStat       interface{} `json:"evidenceStat"`
	Port               interface{} `json:"port"`
	Service            interface{} `json:"service"`
	Component          interface{} `json:"component"`
	OpSys              interface{} `json:"opSys"`
	SubDomain          interface{} `json:"subDomain"`
	UrlTags            interface{} `json:"urlTags"`
}

type TaskConfigInfoReq struct {
	TaskId int `json:"taskId" form:"taskId" binding:"required"`
}

type TaskConfigInfoRes struct {
	Priority   int              `json:"priority"`
	PriorityZh string           `json:"priorityZh"`
	Config     enums.ConfigJson `json:"config"`
}

type TargetChangeStateReq struct {
	TargetId int    `json:"targetId" form:"targetId" binding:"required"`
	Operate  string `json:"operate" form:"operate" binding:"required"`
}

type TargetChangeStateResp struct{}

type TargetListReq struct {
	TaskId    int    `json:"taskId" form:"taskId" binding:"required"`
	Search    string `json:"search" form:"search"`
	RiskLevel int    `json:"riskLevel" form:"riskLevel" `
	Page      int    `json:"page" form:"page" binding:"required"`
	Size      int    `json:"size" form:"size" binding:"required"`
}

type TargetListResp struct {
	List  []TargetListRespItems `json:"list"`
	Total int64                 `json:"total"`
}

type TargetListRespItems struct {
	Id            int      `json:"id"`
	TargetUrl     string   `json:"targetUrl"`
	OpSys         string   `json:"opSys"`
	OpenPort      []string `json:"openPort"`
	RiskLevel     int      `json:"riskLevel"`
	RiskLevelName string   `json:"riskLevelName"`
	VulNum        []int    `json:"vulNum"`
	Status        int      `json:"status"`
	StatusName    string   `json:"statusName"`
	UseScore      string   `json:"useScore"`
	IsAlive       int      `json:"isAlive"`
	IsAliveName   string   `json:"isAliveName"`
}

type UpdateTargetUseScoreReq struct {
	Data map[int]int `json:"data" form:"data" binding:"required"`
}

type UpdateTargetUseScoreResp struct {
}

type TargetDelReq struct {
	TargetIds string `json:"targetIds" form:"targetIds" binding:"required"`
	TaskId    int    `json:"taskId" form:"taskId" binding:"required"`
}

type TargetDelResp struct{}

type TaskResultListReq struct {
	ObjType    int    `json:"objType" form:"objType" binding:"required"`
	SubObjType string `json:"subObjType" form:"subObjType" binding:"required"`
	ObjId      string `json:"objId" form:"objId" binding:"required"`
	Search     string `json:"search" form:"search"`
	Page       int    `json:"page" form:"page" binding:"required"`
	Size       int    `json:"size" form:"size" binding:"required"`
}

type TaskResultListResp struct {
	List  []map[string]interface{} `json:"list"`
	Total int64                    `json:"total"`
}

type TaskResultUrlTreeReq struct {
	ObjId string `form:"objId" json:"objId" binding:"required"`
}

type TaskResultUrlTreeResp struct {
	TreeData []*utils.TreeNode `json:"treeData"`
}

type TaskResultDelReq struct {
	TaskTaskResultIds string `json:"taskTaskResultIds" form:"taskTaskResultIds" binding:"required"`
	TaskId            int    `json:"taskId" form:"taskId" binding:"required"`
}

type TaskResultDelResp struct{}

type TaskResultDetailReq struct {
	TaskTaskResultId int `json:"taskTaskResultId" form:"taskTaskResultId" binding:"required"`
}
type TaskResultDetailResp map[string]interface{}

type VulListReq struct {
	TaskId   int    `json:"taskId" form:"taskId" binding:"required"`
	TargetId int    `json:"targetId" form:"targetId"`
	Search   string `json:"search" form:"search"`
	Type     int    `json:"type" form:"type"`
	Risk     int    `json:"risk" form:"risk"`
	Status   int    `json:"status" form:"status"`
	DataType int    `json:"dataType" form:"dataType" binding:"required"`
	Page     int    `json:"page" form:"page" binding:"required"`
	Size     int    `json:"size" form:"size" binding:"required"`
}

type VulListResp struct {
	List  []VulListRespItems `json:"list"`
	Total int64              `json:"total"`
}

type VulListRespItems struct {
	Id             int    `json:"id"`
	DataType       int    `json:"dataType"`
	TargetUrl      string `json:"targetUrl"`
	Name           string `json:"name"`
	Type           int    `json:"type"`
	TypeName       string `json:"typeName"`
	Risk           int    `json:"risk"`
	RiskName       string `json:"riskName"`
	Location       string `json:"location"`
	Status         int    `json:"status"`
	StatusName     string `json:"statusName"`
	TestStatus     int    `json:"testStatus"`
	TestStatusName string `json:"testStatusName"`
	IsSnapshot     int    `json:"isSnapshot"`
}

type VulInfoReq struct {
	TaskVulId int `json:"taskVulId" form:"taskVulId" binding:"required"`
}

type VulInfoResp struct {
	Id            int                 `json:"id"`
	TaskId        int                 `json:"taskId"`
	TargetId      int                 `json:"targetId"`
	TargetUrl     string              `json:"targetUrl"`
	Pocname       string              `json:"pocname"`
	Name          string              `json:"name"`
	Type          int                 `json:"type"`
	TypeName      string              `json:"typeName"`
	Risk          int                 `json:"risk"`
	RiskName      string              `json:"riskName"`
	VulNumber     string              `json:"vulNumber"`
	Cvss          string              `json:"cvss"`
	PublishedTime string              `json:"publishedTime"`
	ExploitImpact string              `json:"exploitImpact"`
	Description   string              `json:"description"`
	FixSuggest    string              `json:"fixSuggest"`
	RefUrl        string              `json:"refUrl"`
	VulAddress    string              `json:"vulAddress"`
	VulResult     string              `json:"vulResult"`
	VulParam      string              `json:"vulParam"`
	VerMsg        []VulInfoRespVerMsg `json:"verMsg"`
	Location      string              `json:"location"`
	Status        int                 `json:"status"`
	StatusName    string              `json:"statusName"`
	VulId         string              `json:"vulId"`
}

type VulInfoRespVerMsg struct {
	Request            string `json:"request"`
	Response           string `json:"response"`
	Payload            string `json:"payload"`
	PayloadSuccessFlag string `json:"payload_success_flag"`
}

type GetVulSnapshotReq struct {
	TaskVulId int `json:"taskVulId" form:"taskVulId" binding:"required"`
}

type GetVulSnapshotResp struct {
	Snapshot string `json:"snapshot"`
}

type VulDelReq struct {
	TaskVulIds string `json:"taskVulIds" form:"taskVulIds" binding:"required"`
	TaskId     int    `json:"taskId" form:"taskId" binding:"required"`
}

type VulDelResp struct{}

type VulTestReq struct {
	TaskVulId int    `json:"taskVulId" form:"taskVulId" binding:"required"`
	VerMsg    string `json:"verMsg" form:"verMsg" binding:"required"`
}

type VulTestResp struct {
	RespVerMsg string `json:"respVerMsg"`
}

type VulVerifyReq struct {
	TaskId    int    `json:"taskId" form:"taskId" binding:"required"`
	TaskVulId int    `json:"taskVulId" form:"taskVulId" binding:"required"`
	Pocname   string `json:"pocname" form:"pocname" binding:"required"`
	VulParam  string `json:"vulParam" form:"vulParam" binding:"required"`
}

type VulVerifyResp struct{}

type TestVulTestReq struct {
	TaskVulIds string `json:"taskVulIds" form:"taskVulIds" binding:"required"`
	TaskId     int    `json:"taskId" form:"taskId" binding:"required"`
}

type TestVulTestResp struct{}

type AttackLinkReq struct {
	TaskId int `json:"taskId" form:"taskId" binding:"required"`
}

type AttackLinkResp struct {
	Nodes []map[string]interface{} `json:"nodes"`
	Edges []map[string]interface{} `json:"edges"`
}

type LogListReq struct {
	TaskId int    `json:"taskId" form:"taskId" binding:"required"`
	Search string `json:"search" form:"search"`
	Page   int    `json:"page" form:"page" binding:"required"`
	Size   int    `json:"size" form:"size" binding:"required"`
}

type LogListResp struct {
	List  []LogListRespItems `json:"list"`
	Total int64              `json:"total"`
}

type LogListRespItems struct {
	Id          int    `json:"id"`
	TaskId      int    `json:"taskId"`
	TargetId    int    `json:"targetId"`
	TargetUrl   string `json:"targetUrl"`
	Status      int    `json:"status"`
	StatusName  string `json:"statusName"`
	IsAlive     int    `json:"isAlive"`
	IsAliveName string `json:"isAliveName"`
	CreateTime  string `json:"createTime"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

type LogInfoReq struct {
	TaskLogId int `json:"taskLogId" form:"taskLogId" binding:"required"`
}

type LogInfoResp struct {
	List  []LogInfoRespItems `json:"list"`
	Total int64              `json:"total"`
}

type LogInfoRespItems struct {
	Id         int    `json:"id"`
	TaskId     int    `json:"taskId"`
	TargetId   int    `json:"targetId"`
	TargetUrl  string `json:"targetUrl"`
	Pocname    string `json:"pocname"`
	Result     string `json:"result"`
	CreateTime string `json:"createTime"`
}

// 创建
type ApiTaskTaskCreateReq struct {
	TaskTemplateId int              `form:"taskTemplateId" json:"taskTemplateId" binding:"required"`
	Target         string           `form:"target" json:"target" binding:"required"`
	TaskName       string           `form:"taskName" json:"taskName" binding:"required"`
	ExecuteType    int              `form:"executeType" json:"executeType" binding:"required"` // 1即时执行 2定时执行 3周期执行
	ExecuteJson    ApiExecuteJson   `form:"executeJson" json:"executeJson"`
	Config         enums.ConfigJson `form:"config" json:"config"`
}

type ApiExecuteJson struct {
	StartTime          string   `form:"startTime" json:"startTime"`                   // 开始时间
	CyclePlanningType  int      `form:"cyclePlanningType" json:"cyclePlanningType"`   // 周期计划类型 1周 2月
	CyclePlanningValue int      `form:"cyclePlanningValue" json:"cyclePlanningValue"` // 周期计划类型 具体周期 （其中周：1周一 2周二... 7周日） （其中月：1一号 2二号...）
	CyclePlanningHour  string   `form:"cyclePlanningHour" json:"cyclePlanningHour"`   // 周期计划类型 具体小时与分钟 00:01 ｜ 00:20等
	EndTime            string   `form:"endTime" json:"endTime"`                       // 结束时间
	RuntimePeriod      []string `form:"runtimePeriod" json:"runtimePeriod"`
}

type TaskProgressReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"` //
}

type TaskProgressResp struct {
	Progress string `form:"progress" json:"progress" ` //
}

type ApiVulListReq struct {
	TaskId int    `json:"taskId" form:"taskId" binding:"required"`
	Search string `json:"search" form:"search"`
	Type   int    `json:"type" form:"type"`
	Risk   int    `json:"risk" form:"risk"`
	Page   int    `json:"page" form:"page" binding:"required"`
	Size   int    `json:"size" form:"size" binding:"required"`
}

type ApiVulListResp struct {
	List  []ApiVulListRespItems `json:"list"`
	Total int64                 `json:"total"`
}

type ApiVulListRespItems struct {
	Id            int    `json:"id"`
	TargetUrl     string `json:"targetUrl"`
	Name          string `json:"name"`
	Type          int    `json:"type"`
	TypeName      string `json:"typeName"`
	Risk          int    `json:"risk"`
	RiskName      string `json:"riskName"`
	Location      string `json:"location"`
	Status        int    `json:"status"`
	StatusName    string `json:"statusName"`
	TaskId        int    `json:"taskId"`
	TargetId      int    `json:"targetId"`
	Pocname       string `json:"pocname"`
	Cve           string `json:"cve"`
	Cvss          string `json:"cvss"`
	PublishedTime string `json:"publishedTime"`
	ExploitImpact string `json:"exploitImpact"`
	Description   string `json:"description"`
	FixSuggest    string `json:"fixSuggest"`
	RefUrl        string `json:"refUrl"`
	VulAddress    string `json:"vulAddress"`
	VulResult     string `json:"vulResult"`
	VulParam      string `json:"vulParam"`
	VerMsg        string `json:"verMsg"`
	RespVerMsg    string `json:"respVerMsg"`
	Class         int    `json:"class"`
	Ip            string `json:"ip"`
	Port          string `json:"port"`
}

type AddTargetReq struct {
	TaskId int    `form:"taskId" json:"taskId" binding:"required"`
	Target string `form:"target" json:"target" binding:"required"`
	UserId int    `form:"userId" json:"userId" binding:"required"`
}

type AddTargetResp struct{}

type AddAttackFaceReq struct {
	TaskId         int    `form:"taskId" json:"taskId" binding:"required"`
	AttackFaceType int    `form:"attackFaceType" json:"attackFaceType" binding:"required"`
	Target         string `form:"target" json:"target" binding:"required"`
	Params         string `form:"params" json:"params"`
	UserId         int    `form:"userId" json:"userId" binding:"required"`
}

type AddAttackFacePort struct {
	Port int `json:"port"`
}

type AddAttackFacePath struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type AddAttackFaceLoginCred struct {
	Ip    string `json:"ip"`
	Port  string `json:"port"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AddAttackFaceResp struct{}

type AddVulReq struct {
	TaskId  int    `form:"taskId" json:"taskId" binding:"required"`
	RootUrl string `form:"taskId" json:"rootUrl" binding:"required"`
	Pocname string `form:"pocname" json:"pocname"`
	UserId  int    `form:"userId" json:"userId" binding:"required"`
}

type AddVulResp struct{}

// 目标路径图
type TaskTargetMapReq struct {
	TargetId int `form:"targetId" json:"targetId" binding:"required"`
}
type TaskTargetMapRes struct {
	State            []TaskTargetMapState `json:"state"`
	Edg              []TaskTargetMapEdg   `json:"edg"`
	Network          TaskTargetMapNetwork `json:"network"`
	Risk             TaskTargetMapRisk    `json:"risk"`
	Progress         string               `json:"progress"`
	RemoteControlNum string               `json:"remoteControlNum"`
}
type TaskTargetMapState struct {
	Id       string `json:"id"`
	Label    string `json:"label"`
	Status   string `json:"status"`
	Success  bool   `json:"success"`
	Hide     bool   `json:"hide"`
	Pocname  string `json:"pocname"`
	Class    string `json:"class"`
	DetailId string `json:"detailId"`
}
type TaskTargetMapEdg struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Hide  bool   `json:"hide"`
}
type TaskTargetMapRisk struct {
	Deadly int `json:"deadly"`
	High   int `json:"high"`
	Middle int `json:"middle"`
	Low    int `json:"low"`
}
type TaskTargetMapNetwork struct {
	LocalDelay    string `json:"local_delay"`
	LocalQuality  string `json:"local_quality"`
	Network       string `json:"network"`
	TargetDelay   string `json:"target_delay"`
	TargetQuality string `json:"target_quality"`
}

// 目标路径图 - 详情 - 请求
type TaskTargetMapNodeDetailReq struct {
	DetailId string `form:"detailId" json:"detailId" binding:"required"`
}

// 目标路径图 - 详情 - 响应
type TaskTargetMapNodeDetailRes struct {
	Type string      `json:"type"` // 当前数据响应类型
	Data interface{} `json:"data"` // 对应响应类型的数据
}

// 目标路径图 - 详情 - 响应 - 根节点
type TaskTargetMapNodeDetailRoot struct {
	Title        string `json:"title"`
	TemplateName string `json:"templateName"`
	PortScan     string `json:"portScan"`
}

// 目标路径图 - 详情 - 响应 - 人造端口扫描节点
type TaskTargetMapNodeDetailPortScan struct {
	Title string                            `json:"title"`
	Ports []TaskTargetMapNodeDetailPortItem `json:"ports"`
}
type TaskTargetMapNodeDetailPortItem struct {
	Port    string `json:"port"`
	Service string `json:"service"`
}

// 目标路径图 - 详情 - 响应 - 端口节点
type TaskTargetMapNodeDetailPort struct {
	Title   string `json:"title"`
	Service string `json:"service"`
	Result  string `json:"result"`
}

// 目标路径图 - 详情 - 响应 - 漏洞节点
type TaskTargetMapNodeDetailVul struct {
	Title       string `json:"title"`
	Status      int    `json:"status"`
	StatusEnum  string `json:"statusEnum"`
	Type        int    `json:"type"`
	TypeEnum    string `json:"typeEnum"`
	Risk        int    `json:"risk"`
	RiskEnum    string `json:"riskEnum"`
	Description string `json:"Description"`
	FixSuggest  string `json:"fixSuggest"`
	Link        string `json:"link"`
	Location    string `json:"location"`
	Result      string `json:"result"`
}

// 三方 - 任务数据导出 - Three标识
type TaskThreeExportReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}
type TaskThreeExportRes struct {
	VulnFound struct {
		DataSource   DataThreeSource   `json:"data-source"`
		ScanTaskInfo ScanThreeTaskInfo `json:"scan-task-info"`
		HostList     []ThreeHostList   `json:"host-list"`
		VulnList     []ThreeVulnList   `json:"vuln-list"`
	} `json:"VulnFound"`
}

// 三方 - 任务数据导出 - 系统信息
type DataThreeSource struct {
	Name        string `json:"name"`
	ProductName string `json:"product-name"`
	VendorName  string `json:"vendor-name"`
	VersionInfo string `json:"version-info"`
	CreateTime  string `json:"create-time"`
}

// 三方 - 任务数据导出 - 任务信息
type ScanThreeTaskInfo struct {
	TargetIPS string `json:"Target-IPS"`
	BeginTime string `json:"begin-time"`
	EndTime   string `json:"end-time"`
	ScannerIP string `json:"scanner-IP"`
}

// 三方 - 任务数据导出 - 任务目标信息列表
type ThreeHostList struct {
	Hostip   string          `json:"hostip"`
	Ostype   string          `json:"ostype"`
	PortList []ThreePortList `json:"port-list"`
}

// 三方 - 任务数据导出 - 任务目标信息列表 - 端口列表
type ThreePortList struct {
	Portno       int    `json:"portno"`
	Prototype    int    `json:"prototype"`
	ProtocolName string `json:"protocol-name"`
}

// 三方 - 任务数据导出 - 任务下漏洞信息
type ThreeVulnList struct {
	VulnId                string `json:"vulnId"`
	VulnName              string `json:"vulnName"`
	ShortDesc             string `json:"shortDesc"`
	FullDesc              string `json:"fullDesc"`
	RiskLevelValue        string `json:"riskLevelValue"`
	Platforms             string `json:"platforms"`
	AssetIp               string `json:"assetIp"`
	RepairAdvice          string `json:"repairAdvice"`
	CncveTag              string `json:"cncveTag"`
	CveTag                string `json:"cveTag"`
	CnnvdTag              string `json:"cnnvdTag"`
	CnvdTag               string `json:"cnvdTag"`
	CvssScore             string `json:"cvssScore"`
	VulnDisclosureTime    string `json:"vulnDisclosureTime"`
	VulnDisclosureTimeStr string `json:"vulnDisclosureTimeStr"`
	SecurityValue         string `json:"securityValue"`
	PrincipleValue        string `json:"principleValue"`
	Protocol              string `json:"protocol"`
	Port                  string `json:"port"`
	VptValue              string `json:"vptValue"`
	StatusValue           string `json:"statusValue"`
	Payload               string `json:"payload"`
	OriginalRequest       string `json:"originalRequest"`
	ResponseInfo          string `json:"responseInfo"`
}

// 任务详情 - 拓扑图
type TaskTopologyMapReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}
type TaskTopologyMapRes struct {
	Status           int    `json:"status"`
	StatusEnum       string `json:"statusEnum"`
	RemoteControlNum int    `json:"remoteControlNum"`
	RiskTarget       struct {
		High int `json:"high"`
		Mid  int `json:"mid"`
		Low  int `json:"low"`
		Safe int `json:"safe"`
	} `json:"riskTarget"`
	CheckTarget struct {
		Finish  int `json:"finish"`
		Running int `json:"running"`
		Wait    int `json:"wait"`
		Fail    int `json:"fail"`
	} `json:"checkTarget"`
	State []TaskTargetMapState `json:"state"`
	Edg   []TaskTargetMapEdg   `json:"edg"`
}

// 任务详情 - 拓扑图节点详情
type TaskTopologyMapNodeDetailReq struct {
	DetailId string `form:"detailId" json:"detailId" binding:"required"`
}
type TaskTopologyMapNodeDetailRes struct {
	TargetUrl        string `json:"targetUrl"`
	RiskEnum         string `json:"riskEnum"`
	VerifyEnum       string `json:"verifyEnum"`
	Fen              string `json:"fen"`
	ApplicationLevel string `json:"applicationLevel"` // 业务层
	SupportLevel     string `json:"supportLevel"`     // 支撑层
	ServiceLevel     string `json:"serviceLevel"`     // 服务层
	SystemLevel      string `json:"systemLevel"`      // 系统层
	HardwareLevel    string `json:"hardwareLevel"`    // 硬件层
}

// 任务详情 - 拓扑图节点详情 - 漏洞取证
type TaskTopologyMapNodeDetaiVulQuzhengReq struct {
	DetailId string `form:"detailId" json:"detailId" binding:"required"`
	Page     int    `form:"page" json:"page" binding:"required"`
	Size     int    `form:"size" json:"size" binding:"required"`
}
type TaskTopologyMapNodeDetaiVulQuzhengRes struct {
	Total int                                         `json:"total"`
	List  []TaskTopologyMapNodeDetaiVulQuzhengResItem `json:"list"`
}
type TaskTopologyMapNodeDetaiVulQuzhengResItem struct {
	Class string `json:"class"`
	Name  string `json:"name"`
}

// 任务详情 - 拓扑图节点详情 - 远程控制
type TaskTopologyMapNodeDetaiVulKongzhiReq struct {
	DetailId string `form:"detailId" json:"detailId" binding:"required"`
	Page     int    `form:"page" json:"page" binding:"required"`
	Size     int    `form:"size" json:"size" binding:"required"`
}
type TaskTopologyMapNodeDetaiVulKongzhiRes struct {
	Total int                                         `json:"total"`
	List  []TaskTopologyMapNodeDetaiVulKongzhiResItem `json:"list"`
}
type TaskTopologyMapNodeDetaiVulKongzhiResItem struct {
	Class  string `json:"class"`
	Status string `json:"status"`
}

// 获取平台下所有发现的漏洞数据
type TaskAllVulByPageReq struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
}
type TaskAllVulByPageRes struct {
	Total int64                 `json:"total"`
	List  []ApiVulListRespItems `json:"list"`
}

// CheckVulResp 漏洞检测返回
type CheckVulResp struct {
	FalseAlarm bool   `json:"falseAlarm"` // 漏洞是否误报
	Content    string `json:"content"`    // 原因
}

// QianJiAIVulResult  qianji ai漏洞结果返回结构体
type QianJiAIVulResult struct {
	Location           string `json:"location"`   //地址
	VulAddress         string `json:"vulAddress"` //地址
	PayloadSuccessFlag string `json:"payload_success_flag"`
	PocName            string `json:"pocname"`       // 漏洞名
	AiVulName          string `json:"aiVulName"`     // ai返回漏洞名
	FixSuggest         string `json:"fixSuggest"`    // 修复建议
	ReferenceLink      string `json:"referenceLink"` // 参考链接
	Result             string `json:"result"`        // 漏洞结果
	Description        string `json:"description"`   // 漏洞描述
	Cvss               string `json:"cvss"`          // 评分
	VulNumber          string `json:"vulNumber"`     // 漏洞编号
	PublishedTime      string `json:"publishedTime"` // 披露时间
	ExploitImpact      string `json:"exploitImpact"` // 漏洞影响
}

// QianJiAIDescription qianji ai漏洞结果描述结构体
type QianJiAIDescription struct {
	IsSuccess     bool   `json:"is_success"`
	Details       string `json:"details"`
	Request       string `json:"request"`
	Response      string `json:"response"`
	Runtimeid     string `json:"runtimeid"`
	Vuladdress    string `json:"vuladdress"`
	FixSuggest    string `json:"fix_suggest"`
	ExploitImpact string `json:"exploit_impact"`
	Cvss          string `json:"cvss"`
	VulType       string `json:"vul_type"`
}
