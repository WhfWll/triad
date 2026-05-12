package typespec

import "time"

// 任务 - 列表
type VulScanTaskListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Risk   int    `form:"risk" json:"risk"`
	Search string `form:"search" json:"search"`
}
type VulScanTaskListResp struct {
	Total int                      `json:"total"`
	List  []VulScanTaskListResItem `json:"list"`
}
type VulScanTaskListResItem struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Type         int    `json:"type"`
	TypeName     string `json:"typeName"`
	Risk         int    `json:"risk"`
	RiskName     string `json:"riskName"`
	Status       int    `json:"status"`
	StatusName   string `json:"statusName"`
	High         int    `json:"high"`
	Middle       int    `json:"middle"`
	Low          int    `json:"low"`
	Safe         int    `json:"safe"`
	CreateTime   string `json:"createTime"`
	UpdateTime   string `json:"updateTime"`
	TargetNumber int    `json:"targetNumber"` //目标数量
}

// 目标 - 列表
type VulScanTargetListReq struct {
	Page    int    `form:"page" json:"page" binding:"required"`
	Size    int    `form:"size" json:"size" binding:"required"`
	TaskId  int    `form:"taskId" json:"taskId" binding:"required"`
	Risk    int    `form:"risk" json:"risk"`
	Search  string `form:"search" json:"search"`
	IsAlive string `form:"isAlive" json:"isAlive"`
}
type VulScanTargetListResp struct {
	Total int                        `json:"total"`
	List  []VulScanTargetListResItem `json:"list"`
}
type VulScanTargetListResItem struct {
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

// 任务 - 保存
type VulScanTaskSaveReq struct {
	Name         string `json:"name" binding:"required"`
	Target       string `json:"target" binding:"required"`
	ToScanPort   string `json:"toScanPort" binding:"required"`
	OnlyPortScan bool   `json:"only_port_scan"` // 仅端口扫描模式
}
type VulScanTaskSaveResp struct {
	Id int `json:"id"`
}

// 任务 - 结束
type VulScanTaskStopReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type VulScanTaskStopResp struct {
}

// 任务 - 结束
type VulScanTaskDeleteReq struct {
	Ids string `form:"ids" json:"ids" binding:"required"`
}
type VulScanTaskDeleteResp struct {
}

// 漏洞 - 列表
type VulScanVulListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	TaskId int    `form:"taskId" json:"taskId" binding:"required"`
	Risk   int    `form:"risk" json:"risk"`
	Search string `form:"search" json:"search"`
}
type VulScanVulListResp struct {
	Total int                     `json:"total"`
	List  []VulScanVulListResItem `json:"list"`
}
type VulScanVulListResItem struct {
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

// 漏洞 - 列表
type VulScanVulDetailReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type VulScanVulDetailResp struct {
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
	Version             string    `json:"version"`              // 版本号
	Cpes                string    `json:"cpes"`                 //
	CvssVersion         string    `json:"cvss_version"`         // cvss版本
	CvssVector          string    `json:"cvss_vector"`          // cvss向量
	PublishDate         string    `json:"publish_date"`         // 发布时间
	ExploitabilityScore string    `json:"exploitability_score"` // 可利用评分
	References          string    `json:"references"`           // 参考链接
}

// cve - 列表
type VulScanCveListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type VulScanCveListResp struct {
	Total int                     `json:"total"`
	List  []VulScanCveListResItem `json:"list"`
}
type VulScanCveListResItem struct {
	Id                int     `json:"id"`
	CreatedAt         string  `json:"created_at"`          // 创建时间
	UpdatedAt         string  `json:"updated_at"`          // 更新时间
	DeletedAt         string  `json:"deleted_at"`          // 删除时间
	Cve               string  `json:"cve"`                 // cve编号
	Cwe               string  `json:"cwe"`                 // cwe类型
	TitleZh           string  `json:"title_zh"`            // 中文名字
	DescriptionMain   string  `json:"description_main"`    // 描述
	DescriptionMainZh string  `json:"description_main_zh"` // 描述
	Descriptions      string  `json:"descriptions"`        // 描述
	Vendor            string  `json:"vendor"`              // 厂商
	Product           string  `json:"product"`             // 厂商
	Severity          string  `json:"severity"`            // 风险等级
	SeverityName      string  `json:"severityName"`        // 风险等级
	PublishedDate     string  `json:"published_date"`      // 发布日期
	BaseCvssv2Score   float64 `json:"baseCvssv2Score"`     // 发布日期
}

// cve - 详情
type VulScanCveDetailReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type VulScanCveDetailResp struct {
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
	Product                 string  `json:"product"`                   // 厂商
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
}

// 任务概览
type VulScanTaskOverviewReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type VulScanTaskOverviewResp struct {
	TaskName   string `form:"taskName" json:"taskName" binding:"required"`
	Risk       int    `form:"risk" json:"risk" binding:"required"`
	RiskName   string `form:"riskName" json:"riskName" binding:"required"`
	TargetRisk []int  `form:"targetRisk" json:"targetRisk" binding:"required"`
	TargetNum  []int  `form:"targetNum" json:"targetNum" binding:"required"`
	VulRisk    []int  `form:"vulRisk" json:"vulRisk" binding:"required"`
	CreateTime string `form:"createTime" json:"createTime" binding:"required"`
	UpdateTime string `form:"updateTime" json:"updateTime" binding:"required"`
	Ports      string `form:"ports" json:"ports" binding:"required"`
}

// 获取任务状态
type VulScanTaskGetStateReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type VulScanTaskGetStateResp struct {
	Status     int    `form:"status" json:"status" binding:"required"`
	StatusName string `form:"statusName" json:"statusName" binding:"required"`
}
