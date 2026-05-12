package typespec

type VulRiskListReq struct {
	TaskId     int    `json:"taskId" form:"taskId"` // 任务ID
	TargetId   int    `json:"targetId" form:"targetId"`
	Search     string `json:"search" form:"search"`         // 漏洞名称/位置
	Location   string `json:"location" form:"location"`     // 漏洞地址
	IP         string `json:"ip" form:"ip"`                 // IP
	Type       int    `json:"type" form:"type"`             // 漏洞类型
	Risk       int    `json:"risk" form:"risk"`             // 漏洞风险
	Status     int    `json:"status" form:"status"`         // 风险状态
	VerifyType int    `json:"verifyType" form:"verifyType"` // 验证状态
	Page       int    `json:"page" form:"page" binding:"required"`
	Size       int    `json:"size" form:"size" binding:"required"`
}

type VulRiskListResp struct {
	List  []VulRiskListRespItems `json:"list"`
	Total int64                  `json:"total"`
}

type VulRiskListRespItems struct {
	// 新需求需要的内容
	ID             int    `json:"id"`             // 主键ID
	TargetUrl      string `json:"targetUrl"`      // 测试目标
	Location       string `json:"location"`       // 漏洞地址
	Name           string `json:"name"`           // 漏洞名称
	Type           int    `json:"type"`           // 漏洞类型
	TypeName       string `json:"typeName"`       // 漏洞类型名
	RiskLevel      int    `json:"riskLevel"`      // 漏洞等级
	RiskLevelStr   string `json:"riskLevelStr"`   // 漏洞等级名称
	VerifyType     int    `json:"verifyType"`     // 漏洞验证方式
	VerifyTypeName string `json:"verifyTypeName"` // 漏洞验证方式str
	Status         int    `json:"status"`         // 漏洞状态
	StatusName     string `json:"statusName"`     // 漏洞状态名
	FixStatus      int    `json:"fixStatus"`      // 修复状态
	FixStatusStr   string `json:"fixStatusStr"`   // 修复状态名称
	FindTime       string `json:"findTime"`       // 发现时间
	CreateTime     string `json:"createTime"`     // 创建时间
	FixTime        string `json:"fixTime"`        // 修复时间
	// 列表不展示 注释
	//Num              int      `json:"num"`              // 漏洞数量
	//CVSS             string   `json:"cvss"`             // CVSS值
	//Desc             string   `json:"desc"`             // 漏洞描述
	//AssetList        []string `json:"assetList"`        // 相关资产信息
	//FixTimeLength    int64    `json:"fixTimeLength"`    // 修复时长 默认秒
	//FixTimeLengthStr string   `json:"fixTimeLengthStr"` // 修复时长
	//TestStatus       int      `json:"testStatus"`       // 测试状态
	//TestStatusName   string   `json:"testStatusName"`   // 测试状态str
}

// VulRiskAssetList 漏洞相关资产列表
type VulRiskAssetList struct {
	AssetID   int    `json:"assetID"`
	TargetUrl string `json:"targetUrl"` // 测试目标
}

type VulRiskInfoReq struct {
	ID int `json:"id" form:"id" binding:"required"`
}

type VulRiskDeleteReq struct {
	ID int `json:"id" form:"id" binding:"required"`
}

type VulRiskUpdateReq struct {
	IDs    string `json:"ids" form:"ids" binding:"required"`
	Status int    `json:"status" form:"status" binding:"required"`
}

type VulRiskInfoResp struct {
	Id               int                 `json:"id"`
	TaskId           int                 `json:"taskId"`       // 任务ID
	TargetId         int                 `json:"targetId"`     // 目标ID
	Name             string              `json:"name"`         // 漏洞名称
	TargetUrl        string              `json:"targetUrl"`    // 测试目标/所属资产
	Risk             int                 `json:"risk"`         // 漏洞风险
	RiskName         string              `json:"riskName"`     // 漏洞风险名称
	Description      string              `json:"description"`  // 漏洞描述
	FixSuggest       string              `json:"fixSuggest"`   // 修复建议
	VerMsg           []VulInfoRespVerMsg `json:"verMsg"`       // 请求报文
	VulLifecycleInfo []TaskLifecycleLog  `json:"vulLifecycle"` // 漏洞生命周期
	// 展示不展示的内容
	Type          int    `json:"type"`          // 漏洞类型
	TypeName      string `json:"typeName"`      // 漏洞类型名称
	Pocname       string `json:"pocname"`       // poc名称
	PublishedTime string `json:"publishedTime"` // 披露时间
	ExploitImpact string `json:"exploitImpact"`
	RefUrl        string `json:"refUrl"`
	VulAddress    string `json:"vulAddress"`
	VulResult     string `json:"vulResult"`
	VulParam      string `json:"vulParam"`
	Location      string `json:"location"`
	Status        int    `json:"status"`
	StatusName    string `json:"statusName"`
	VulId         string `json:"vulId"`
	PatchUrl      string `json:"patch_url"` // 补丁地址
}

// TaskLifecycleLog 对应数据库 content 字段中的单个任务块
type TaskLifecycleLog struct {
	TaskID    int             `json:"task_id"`
	TaskName  string          `json:"task_name"`
	Lifecycle []LifecycleItem `json:"lifecycle"` // 使用数组保证日志顺序
}

type LifecycleItem struct {
	Time    string `json:"time"`    // 时间单独存
	Content string `json:"content"` // 描述内容单独存
}

type VulRiskInfoRespVerMsg struct {
	Request            string `json:"request"`
	Response           string `json:"response"`
	Payload            string `json:"payload"`
	PayloadSuccessFlag string `json:"payload_success_flag"`
}

type VulRiskDelReq struct {
	IDs string `json:"ids" form:"ids" binding:"required"`
}

type VulRiskDetailReq struct {
	ID int `json:"id" form:"id"`
}

type VulRiskTestReq struct {
	ID     int    `json:"id" form:"id" binding:"required"`
	VerMsg string `json:"verMsg" form:"verMsg" binding:"required"`
}

type VulRiskTestResp struct {
	RespVerMsg string `json:"respVerMsg"`
}

type VulRiskVerifyReq struct {
	ID       int    `json:"id" form:"id" binding:"required"`
	TaskId   int    `json:"taskId" form:"taskId"`
	Pocname  string `json:"pocname" form:"pocname"`
	VulParam string `json:"vulParam" form:"vulParam"`
}

// VulRiskStaticsRes 漏洞风险统计结果
type VulRiskStaticsRes struct {
	// 漏洞状态统计
	StatusStatistics []VulStatusCount `json:"statusStatistics"`

	// 漏洞平均修复时间统计
	AverageFixTime string `json:"averageFixTime"`

	// TOP10 漏洞统计
	Top10Vulnerabilities []VulCount `json:"top10Vulnerabilities"`

	// 漏洞等级统计
	RiskLevelStatistics []VulRiskLevelCount `json:"riskLevelStatistics"`

	// 漏洞类型统计
	TypeStatistics []VulTypeCount `json:"typeStatistics"`

	TotalVulnerabilities int64 `json:"totalVulnerabilities"` // 总漏洞数
	HighRiskCount        int   `json:"highRiskCount"`        // 高危漏洞数
}

// VulStatusCount 漏洞状态
type VulStatusCount struct {
	StatusName string `json:"statusName"`
	Count      int    `json:"count"`
}

// VulRiskLevelCount 漏洞等级计数
type VulRiskLevelCount struct {
	RiskName string `json:"riskName"`
	Count    int    `json:"count"`
}

// VulTypeCount 漏洞类型计数
type VulTypeCount struct {
	TypeName string `json:"typeName"`
	Count    int    `json:"count"`
}

// VulCount 漏洞数 (用于TOP10)
type VulCount struct {
	Name         string `json:"name"`
	Num          int    `json:"num"`
	RiskLevel    int    `json:"riskLevel"`
	RiskLevelStr string `json:"riskLevelStr"`
	CVSS         string `json:"cvss"`
}

// VulTypeStatisticsInfo 漏洞类型统计
type VulTypeStatisticsInfo struct {
	AssetType string
	Count     int
}

// VulRiskStatisticsInfo 漏洞风险统计
type VulRiskStatisticsInfo struct {
	AssetType  string
	Count      int
	Percentage string
}

// ConfigVulListReq 配置风险请求结构体
type ConfigVulListReq struct {
	Search   string `json:"search" form:"search"`     // 检索检测项名称
	IP       string `json:"ip" form:"ip"`             // 搜索资产IP
	CheckRes int    `json:"checkRes" form:"checkRes"` // 检测结果
	Sort     int    `json:"sort" form:"sort"`         // 检测时间排序 1 desc 2 asc
	Page     int    `json:"page" form:"page" binding:"required"`
	Size     int    `json:"size" form:"size" binding:"required"`
}

type ConfigVulListResp struct {
	List  []ConfigVulListItem `json:"list"`
	Count int                 `json:"count"`
}

// BaseLineScriptContent 基线脚本内容
type BaseLineScriptContent struct {
	Name        string `json:"name"`
	Standards   string `json:"standards"`
	RiskWarning string `json:"riskWarning"`
}

// ConfigVulListItem 风险配置列表返回
type ConfigVulListItem struct {
	ID               int    `json:"id"`               // 主键ID
	IP               string `json:"ip"`               // 资产IP
	ItemCheckedName  string `json:"itemCheckedName"`  // 检测项名
	ItemCheckedDesc  string `json:"itemCheckedDesc"`  // 检测项描述
	ItemCategoryName string `json:"itemCategoryName"` // 检测分类
	CheckRes         int    `json:"checkRes"`         // 检测结果,1-通过 2-不通过 3-错误
	CheckResMsg      string `json:"checkResMsg"`      // 脚本日志校验结果说明
	CheckTime        string `json:"checkTime"`        // 检测时间
	// other
	//CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`     // 创建时间
	//ItemCategoryName string    `json:"itemCategoryName"`                         // 分类
	//OutputLog        string    `gorm:"column:output_log" json:"outputLog"`       // 日志
	//ReinforceBtn     int       `gorm:"column:reinforce_btn" json:"reinforceBtn"` // 是否展示加固按钮,1-不展示 2-展示
	//RollbackBtn      int       `gorm:"column:rollback_btn" json:"rollbackBtn"`   // 是否展示回退按钮,1-不展示 2-展示
	//UpdateTime       time.Time `gorm:"column:update_time" json:"updateTime"`     // 修改时间
}

// ConfigVulRiskStaticsRes 配置风险统计返回
type ConfigVulRiskStaticsRes struct {
	ConfigNoPassStatistics []ConfigNoPassStatisticsInfo `json:"configNoPassStatistics"` // 未通过检测项统计
	ConfigNoPassNum        int                          `json:"configNoPassNum"`        // 未通过项
	ConfigPassNum          int                          `json:"configPassNum"`          // 通过项
}

// ConfigNoPassStatisticsInfo 未通过检测项统计
type ConfigNoPassStatisticsInfo struct {
	Type  string
	Count int
}

type ConfigRiskDelReq struct {
	IDs string `json:"ids" form:"ids" binding:"required"`
}

// RiskManageEnumsRes 风险管理返回结构体
type RiskManageEnumsRes struct {
	VulVerifyType interface{} `json:"vulVerifyType"`
	VulRiskStatus interface{} `json:"vulRiskStatus"`
	VulType       interface{} `json:"vulType"` // 漏洞类型
	VulRisk       interface{} `json:"vulRisk"` // 漏洞风险
}
