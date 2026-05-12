package typespec

type ReportEnumResp struct {
	Type          interface{} `json:"type"`
	Status        interface{} `json:"status"`
	Format        interface{} `json:"format"`
	TaskContent   interface{} `json:"taskContent"`
	TargetContent interface{} `json:"targetContent"`
}

type ReportListReq struct {
	Search string `json:"search" form:"search"`
	Page   int    `json:"page" form:"page" binding:"required"`
	Size   int    `json:"size" form:"size" binding:"required"`
}

type ReportListResp struct {
	List  []ReportListRespItems `json:"list"`
	Total int64                 `json:"total"`
}

type ReportListRespItems struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	TypeName   string `json:"typeName"`
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
	Format     int    `json:"format"`
	FormatName string `json:"formatName"`
	CreateTime string `json:"createTime"`
}

type ReportDownloadReq struct {
	ReportId int `json:"reportId" form:"reportId" binding:"required"`
}

type ReportDownloadResp struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	Status     int    `json:"status"`
	ConfigJson string `json:"configJson"`
	Format     int    `json:"format"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type ReportDelReq struct {
	ReportIds string `json:"reportIds" form:"reportIds" binding:"required"`
}

type ReportDelResp struct {
}

// ReportSaveReq 报告下载
type ReportSaveReq struct {
	Name       string `json:"name" form:"name" binding:"required"`
	Type       int    `json:"type" form:"type" binding:"required"`
	ConfigJson string `json:"configJson" form:"configJson" binding:"required"`
	Format     int    `json:"format" form:"format" binding:"required"`
	UserId     int    `json:"userId" form:"userId" binding:"required"`
	OutputType int    `json:"outputType" form:"outputType"` // 1合并输出 2逐个输出
	ObjIDName  string `json:"objIDName" form:"objIDName"`   // 逐个输出会用到 为了规避objid和目标对不上的问题
}

type ReportSaveResp struct {
}

// 生成报告接口
type ReportGenerateReq struct {
	ReportId int `form:"reportId" json:"reportId" binding:"required"`
}

/*******************************************************/
/***********************报告内容*************************/
/*******************************************************/
// 报告封面
type ReportCoverNode struct {
	Title         string `json:"title"`
	CreateTime    string `json:"createTime"`
	BackgroundImg string `json:"backgroundImg"`
}

// 报告分类信息
type CatalogNode struct {
	Name    string        `json:"name"`
	Id      string        `json:"id"`
	IsShow  bool          `json:"isShow"`
	Catalog []CatalogNode `json:"catalog"`
}

/***************任务报告***********************/
// 任务报告 导航
type ReportTaskContent struct {
	ReportId      int             `json:"reportId"`      // 任务报告--id
	ReportCover   ReportCoverNode `json:"reportCover"`   // 任务报告--封面信息
	Catalog       []CatalogNode   `json:"catalogParent"` // 任务报告--目录信息
	TaskOverview  TaskOverview    `json:"taskOverview"`  // 任务报告--任务概述
	TargetRisk    TargetRisk      `json:"targetRisk"`    // 任务报告--目标风险
	VulRisk       []VulRisk       `json:"vulRisk"`       // 任务报告--漏洞风险
	VulType       []VulType       `json:"vulType"`       // 任务报告--漏洞类型
	TopVulRisk    []TopVulRisk    `json:"topVulRisk"`    // 任务报告--Top危险漏洞
	TargetDetails []TargetDetails `json:"targetDetails"` // 任务报告--目标风险
	VulDetails    []VulDetails    `json:"vulDetails"`    // 任务报告--漏洞详情
}

// 任务报告--任务概述
type TaskOverview struct {
	TaskName    string `json:"taskName"`    // 任务报告--任务概述--任务名称
	TaskRiskStr string `json:"taskRiskStr"` // 任务报告--任务概述--任务风险
	TargetStat  struct {
		Total        int `json:"total"`        // 任务报告--任务概述--目标分布--总目标数量
		LiveTarget   int `json:"liveTarget"`   // 任务报告--任务概述--目标分布--存活目标数量
		HighTarget   int `json:"HighTarget"`   // 任务报告--任务概述--目标分布--高危目标数量
		MiddleTarget int `json:"middleTarget"` // 任务报告--任务概述--目标分布--中危目标数量
		LowTarget    int `json:"lowTarget"`    // 任务报告--任务概述--目标分布--低危目标数量
		SafeTarget   int `json:"safeTarget"`   // 任务报告--任务概述--目标分布--安全目标数量
	} `json:"targetStat"` // 任务报告--任务概述--目标分布
	VulnStat struct {
		Total        int `json:"total"`        // 任务报告--任务概述--漏洞分布--总漏洞数量
		DeadlyNumber int `json:"deadlyNumber"` // 任务报告--任务概述--漏洞分布--致命漏洞数量
		HighNumber   int `json:"highNumber"`   // 任务报告--任务概述--漏洞分布--高危漏洞数量
		MiddleNumber int `json:"middleNumber"` // 任务报告--任务概述--漏洞分布--中危漏洞数量
		LowNumber    int `json:"lowNumber"`    // 任务报告--任务概述--漏洞分布--低危漏洞数量
	} `json:"vulnStat"` // 任务报告--任务概述--漏洞分布
	VulnVerify struct {
		VerifySuccess int `json:"verifySuccess"` // 任务报告--任务概述--漏洞验证--验证存在数量
		UseSuccess    int `json:"useSuccess"`    // 任务报告--任务概述--漏洞验证--利用成功数量
		RepairSuccess int `json:"repairSuccess"` // 任务报告--任务概述--漏洞验证--未验证数量
	} `json:"vulnVerify"` // 任务报告--任务概述--漏洞验证
	TemplateName string `json:"templateName"` // 任务报告--任务概述--任务场景
	Date         string `json:"date"`         // 任务报告--任务概述--测试时间
}

// 任务报告--目标风险
type TargetRisk struct {
	Total            int    `json:"total"`            // 任务报告--目标风险--目标总数
	HighNumber       int    `json:"highNumber"`       // 任务报告--目标风险--高危目标数量
	HighNumberRate   string `json:"highNumberRate"`   // 任务报告--目标风险--高危目标数量 占比
	MiddleNumber     int    `json:"MiddleNumber"`     // 任务报告--目标风险--中危目标数量
	MiddleNumberRate string `json:"MiddleNumberRate"` // 任务报告--目标风险--中危目标数量 占比
	LowNumber        int    `json:"lowNumber"`        // 任务报告--目标风险--低危目标数量
	LowNumberRate    string `json:"lowNumberRate"`    // 任务报告--目标风险--低危目标数量 占比
	SafeNumber       int    `json:"safeNumber"`       // 任务报告--目标风险--安全危目标数量
	SafeNumberRate   string `json:"safeNumberRate"`   // 任务报告--目标风险--安全危目标数量 占比
}

// 任务报告--风险统计
type VulRisk struct {
	RiskType      string `json:"riskType"`      // 任务报告--风险统计--风险类型
	VerifySuccess int    `json:"verifySuccess"` // 任务报告--风险统计--验证成功
	RepairSuccess int    `json:"repairSuccess"` // 任务报告--风险统计--未验证
	UseSuccess    int    `json:"useSuccess"`    // 任务报告--风险统计--利用成功
	Total         int    `json:"total"`         // 任务报告--风险统计--漏洞总数
	Percent       string `json:"percent"`       // 任务报告--风险统计--占比
}

// 任务报告--漏洞类型
type VulType struct {
	VulnType     string `json:"vulnType"`     // 任务报告--漏洞类型--漏洞类型
	Total        int    `json:"total"`        // 任务报告--漏洞类型--数量
	Percent      string `json:"percent"`      // 任务报告--漏洞类型--占比
	TargetNumber int    `json:"targetNumber"` // 任务报告--漏洞类型--目标数量
}

// 任务报告--Top危险漏洞
type TopVulRisk struct {
	VulName       string `json:"vulName"`       // 任务报告--Top危险漏洞--漏洞名称
	Risk          string `json:"risk"`          // 任务报告--Top危险漏洞--漏洞风险
	Number        int    `json:"number"`        // 任务报告--Top危险漏洞--出现次数
	AffectTargets string `json:"affectTargets"` // 任务报告--Top危险漏洞--影响目标
}

// 任务报告--目标风险
type TargetDetails struct {
	Target       string `json:"target"`       // 任务报告--目标风险--目标
	Risk         string `json:"risk"`         // 任务报告--目标风险--风险
	DeadlyNumber int    `json:"deadlyNumber"` // 任务报告--目标风险--致命漏洞数量
	HighNumber   int    `json:"highNumber"`   // 任务报告--目标风险--高危漏洞数量
	MiddleNumber int    `json:"middleNumber"` // 任务报告--目标风险--中危漏洞数量
	LowNumber    int    `json:"lowNumber"`    // 任务报告--目标风险--低危漏洞数量
	VulStatus    string `json:"vulStatus"`    // 任务报告--目标风险--漏洞状态
}

// 任务报告--漏洞详情
type VulDetails struct {
	VulName       string `json:"vulName"`          // 任务|目标报告--漏洞详情--漏洞名称
	Risk          string `json:"risk"`             // 任务|目标报告--漏洞详情--漏洞风险
	VulStatus     string `json:"vulStatus"`        // 任务|目标报告--漏洞详情--漏洞状态
	Number        int    `json:"number"`           // 任务|目标报告--漏洞详情--出现次数
	Type          string `json:"type"`             // 任务|目标报告--漏洞详情--漏洞类型
	Cve           string `json:"cve"`              // 任务|目标报告--漏洞详情--CVE
	PublishDate   string `json:"publishDate"`      // 任务|目标报告--漏洞详情--漏洞公开日期
	Describe      string `json:"describe"`         // 任务|目标报告--漏洞详情--漏洞描述
	Res           string `json:"res"`              // 任务|目标报告--漏洞详情--漏洞结果 仅目标报告需要
	Fix           string `json:"fix"`              // 任务|目标报告--漏洞详情--漏洞修复建议
	AffectRange   string `json:"affectRange"`      // 任务|目标报告--漏洞详情--漏洞影响范围
	AffectTargets string `json:"AffectTargets"`    // 任务|目标报告--漏洞详情--漏洞影响目标
	Location      string `json:"location"`         // 任务|目标报告--漏洞详情--漏洞位置 仅目标报告需要
	Link          string `json:"link"`             // 任务|目标报告--漏洞详情--参考连接
	VerMsg        VerMsg `json:"verMsg,omitempty"` // 任务|目标报告-验证报文
}

// VerMsg .任务|目标报告-验证报文
type VerMsg struct {
	Request            string `json:"request"`
	Response           string `json:"response"`
	Payload            string `json:"payload"`
	PayloadSuccessFlag string `json:"payload_success_flag"`
}

/*************************目标报告***********************/
// 目标报告
type ReportTargetContent struct {
	ReportId       int             `json:"reportId"`       // 目标报告--id
	ReportCover    ReportCoverNode `json:"reportCover"`    // 目标报告--封面信息
	Catalog        []CatalogNode   `json:"catalogParent"`  // 目标报告--目录信息
	TargetOverview TargetOverview  `json:"targetOverview"` // 目标报告--报告摘要
	AssetInfo      AssetReportInfo `json:"assetInfo"`      // 目标报告--资产信息
	VulInfo        []VulDetails    `json:"vulInfo"`        // 目标报告--漏洞信息
}

// 目标报告--报告摘要
type TargetOverview struct {
	Target   string `json:"target"`
	Risk     string `json:"risk"`
	VulnStat struct {
		Total        int `json:"total"`        // 目标报告--报告摘要--漏洞分布--总漏洞数量
		DeadlyNumber int `json:"deadlyNumber"` // 目标报告--报告摘要--漏洞分布--致命漏洞数量
		HighNumber   int `json:"highNumber"`   // 目标报告--报告摘要--漏洞分布--高危漏洞数量
		MiddleNumber int `json:"middleNumber"` // 目标报告--报告摘要--漏洞分布--中危漏洞数量
		LowNumber    int `json:"lowNumber"`    // 目标报告--报告摘要--漏洞分布--低危漏洞数量
	} `json:"vulnStat"` // 目标报告--报告摘要--漏洞分布
	VulnVerify struct {
		VerifySuccess int `json:"verifySuccess"` // 目标报告--任务概述--漏洞验证--验证存在数量
		UseSuccess    int `json:"useSuccess"`    // 目标报告--任务概述--漏洞验证--利用成功数量
		RepairSuccess int `json:"repairSuccess"` // 目标报告--任务概述--漏洞验证--未验证数量
	} `json:"vulnVerify"` // 目标报告--任务概述--漏洞验证
	CreateDate string `json:"createDate"` // 目标报告--任务概述--测试时间
}

// 目标报告--资产信息
type AssetReportInfo struct {
	Component string `json:"component"` // 目标报告--资产信息--组建
	Service   string `json:"service"`   // 目标报告--资产信息--服务
	IpOrUrl   string `json:"ipOrUrl"`   // 目标报告--资产信息--ip/url
	System    string `json:"system"`    // 目标报告--资产信息--系统
}
