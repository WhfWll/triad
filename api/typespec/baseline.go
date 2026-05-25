package typespec

type BaselineCheckReq struct {
	TaskID   int    `json:"taskId" form:"taskId"` // 可选；<=0 时由服务端自动生成一次核查批次 ID
	TargetID int    `json:"targetId" form:"targetId"`
	Host     string `json:"host" form:"host" binding:"required"`
	Port     int    `json:"port" form:"port"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"`
	OSType   int    `json:"osType" form:"osType" binding:"required"`
	// Transport 连接方式：0=自动（Windows→WinRM，其它→SSH），1=强制 SSH，2=强制 WinRM
	Transport int `json:"transport" form:"transport"`
	// WinRMUseHttps WinRM 是否走 HTTPS（典型端口 5986；false 为 HTTP 5985）
	WinRMUseHttps bool `json:"winrmUseHttps" form:"winrmUseHttps"`
	// ScanScene 1=安全配置核查 2=主机漏洞检测（与场景 1 共用远程引擎与规则库，仅任务归类不同）
	ScanScene int `json:"scanScene" form:"scanScene"`
}

// BaselineBatchCheckReq 批量核查请求，多个目标共享一个 taskId
type BaselineBatchCheckReq struct {
	TaskID  int                `json:"taskId" form:"taskId"` // 可选；<=0 时由服务端自动生成
	Targets []BaselineCheckReq `json:"targets" binding:"required,min=1"`
}

// BaselineBatchCheckResp 批量核查响应（异步模式：创建后立即返回 taskId）
type BaselineBatchCheckResp struct {
	TaskID int `json:"taskId"`
}

// BaselineBatchTaskProgress 批量任务进度
type BaselineBatchTaskProgress struct {
	TaskID           int                   `json:"taskId"`
	Status           string                `json:"status"` // running / completed / failed
	TotalTargets     int                   `json:"totalTargets"`
	CompletedTargets int                   `json:"completedTargets"`
	Targets          []BatchTargetProgress `json:"targets"`
	CreatedAt        string                `json:"createdAt"`
}

// BatchTargetProgress 单个目标进度
type BatchTargetProgress struct {
	Host    string `json:"host"`
	Status  string `json:"status"` // pending / running / completed / failed
	Message string `json:"message,omitempty"`
}

type BaselineCheckResp struct {
	TaskID     int                 `json:"taskId"`
	TargetIP   string              `json:"targetIP"`
	OSType     int                 `json:"osType"`
	OSTypeName string              `json:"osTypeName"`
	TotalRules int                 `json:"totalRules"`
	PassCount  int                 `json:"passCount"`
	FailCount  int                 `json:"failCount"`
	ErrorCount int                 `json:"errorCount"`
	Results    []BaselineCheckItem `json:"results"`
	StartTime  string              `json:"startTime"`
	EndTime    string              `json:"endTime"`
}

type BaselineCheckItem struct {
	ID              int    `json:"id"`
	TargetID        int    `json:"targetId"`
	TargetIP        string `json:"targetIp"`
	RuleID          int    `json:"ruleId"`
	RuleName        string `json:"ruleName"`
	RuleCategory    int    `json:"ruleCategory"`
	CategoryName    string `json:"categoryName"`
	RiskLevel       int    `json:"riskLevel"`
	RiskName        string `json:"riskName"`
	CheckResult     int    `json:"checkResult"`
	ResultName      string `json:"resultName"`
	ExpectedValue   string `json:"expectedValue"`
	ActualValue     string `json:"actualValue"`
	CheckCommand    string `json:"checkCommand"`
	FixSuggestion   string `json:"fixSuggestion"`
	RiskDescription string `json:"riskDescription"`
	CheckTime       string `json:"checkTime"`
}

type BaselineCheckResultListReq struct {
	TaskID   int `json:"taskId" form:"taskId"`
	TargetID int `json:"targetId" form:"targetId"`
	Page     int `json:"page" form:"page"`
	Size     int `json:"size" form:"size"`
}

type BaselineCheckResultListResp struct {
	List  []BaselineCheckItem `json:"list"`
	Total int64               `json:"total"`
}

type BaselineStatReq struct {
	TaskID int `json:"taskId" form:"taskId"`
}

type BaselineStatResp struct {
	TotalRules int     `json:"totalRules"`
	PassCount  int     `json:"passCount"`
	FailCount  int     `json:"failCount"`
	PassRate   float64 `json:"passRate"`
}

// BaselineTaskListReq 按核查批次（task_id）分页列表
type BaselineTaskListReq struct {
	Page      int `json:"page" form:"page"`
	Size      int `json:"size" form:"size"`
	ScanScene int `json:"scanScene" form:"scanScene"` // 0=全部 1=安全配置核查 2=主机漏洞检测
}

type BaselineTaskListItem struct {
	TaskID        int    `json:"taskId"`
	TargetIP      string `json:"targetIp"`
	OSType        int    `json:"osType"`
	OSTypeName    string `json:"osTypeName"`
	ScanScene     int    `json:"scanScene"`
	ScanSceneName string `json:"scanSceneName"`
	TotalRules    int    `json:"totalRules"`
	PassCount     int    `json:"passCount"`
	FailCount     int    `json:"failCount"`
	ErrorCount    int    `json:"errorCount"`
	CheckTime     string `json:"checkTime"`
}

type BaselineTaskListResp struct {
	List  []BaselineTaskListItem `json:"list"`
	Total int64                  `json:"total"`
}

// HostSecTaskDeleteItem 主机安全检查任务删除项（列表一行对应一条记录）
type HostSecTaskDeleteItem struct {
	Source    string `json:"source"` // baseline / vuln / malware
	TaskID    int    `json:"taskId" binding:"required"`
	TargetIP  string `json:"targetIp"`
	ScanScene int    `json:"scanScene"` // baseline/vuln 时：1=配置核查 2=漏洞检测
}

type HostSecTaskDeleteReq struct {
	Items []HostSecTaskDeleteItem `json:"items" binding:"required,min=1"`
}

type HostSecTaskDeleteResp struct {
	Deleted int `json:"deleted"`
}

// BaselineTaskTargetItem 任务的目标列表
type BaselineTaskTargetItem struct {
	TargetID   int    `json:"targetId"`
	TargetIP   string `json:"targetIp"`
	OSType     int    `json:"osType"`
	OSTypeName string `json:"osTypeName"`
	TotalRules int    `json:"totalRules"`
	PassCount  int    `json:"passCount"`
	FailCount  int    `json:"failCount"`
	ErrorCount int    `json:"errorCount"`
}

// BaselineRulesStatsResp 规则库统计（仅数量，供仪表盘等轻量场景）
type BaselineRulesStatsResp struct {
	Total      int                            `json:"total"`
	ByOsType   []BaselineRulesCountByOS       `json:"byOsType"`
	ByCategory []BaselineRulesCountByCategory `json:"byCategory"`
}

// BaselineRulesListResp 规则库：总条数、按操作系统/核查分类汇总、明细（供界面展示与验收对照）
type BaselineRulesListResp struct {
	Total      int                            `json:"total"`
	ByOsType   []BaselineRulesCountByOS       `json:"byOsType"`
	ByCategory []BaselineRulesCountByCategory `json:"byCategory"`
	Rules      []BaselineRuleListItem         `json:"rules"`
}

type BaselineRulesCountByOS struct {
	OSType     int    `json:"osType"`
	OSTypeName string `json:"osTypeName"`
	Count      int    `json:"count"`
}

type BaselineRulesCountByCategory struct {
	Category     int    `json:"category"`
	CategoryName string `json:"categoryName"`
	Count        int    `json:"count"`
}

type BaselineRuleListItem struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	CategoryName    string   `json:"categoryName"`
	Risk            int      `json:"risk"`
	RiskName        string   `json:"riskName"`
	OSType          int      `json:"osType"`
	OSTypeName      string   `json:"osTypeName"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion,omitempty"`
	RiskDescription string   `json:"riskDescription,omitempty"`
	Commands        []string `json:"commands,omitempty"`
}

// BaselineRulesImportReq 规则导入请求
type BaselineRulesImportReq struct {
	Rules []BaselineRuleImportItem `json:"rules" binding:"required"`
}

// BaselineRuleImportItem 单条导入规则
type BaselineRuleImportItem struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	Risk            int      `json:"risk"`
	OSType          int      `json:"osType"`
	Commands        []string `json:"commands"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
}

// BaselineRulesImportResp 规则导入响应
type BaselineRulesImportResp struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Skipped int `json:"skipped"`
}

// BaselineRuleCreateReq 新增规则请求
type BaselineRuleCreateReq struct {
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	Risk            int      `json:"risk"`
	OSType          int      `json:"osType"`
	Commands        []string `json:"commands"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
	Enabled         int      `json:"enabled"`
}

// BaselineRuleUpdateReq 编辑规则请求
type BaselineRuleUpdateReq struct {
	ID              int      `json:"id" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	Risk            int      `json:"risk"`
	OSType          int      `json:"osType"`
	Commands        []string `json:"commands"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
	Enabled         int      `json:"enabled"`
}

// BaselineRuleDetailResp 规则详情响应
type BaselineRuleDetailResp struct {
	ID              int      `json:"id"`
	RuleCode        int      `json:"ruleCode"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	CategoryName    string   `json:"categoryName"`
	Risk            int      `json:"risk"`
	RiskName        string   `json:"riskName"`
	OSType          int      `json:"osType"`
	OSTypeName      string   `json:"osTypeName"`
	Commands        []string `json:"commands"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
	Enabled         int      `json:"enabled"`
	CreateTime      string   `json:"createTime"`
	UpdateTime      string   `json:"updateTime"`
}

type MalwareScanReq struct {
	TaskID   int    `json:"taskId" form:"taskId"` // 可选；<=0 时由服务端自动生成
	TargetID int    `json:"targetId" form:"targetId"`
	Host     string `json:"host" form:"host" binding:"required"`
	Port     int    `json:"port" form:"port"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"`
	OSType   int    `json:"osType" form:"osType"`
}

type MalwareScanResp struct {
	TaskID     int                `json:"taskId"`
	TargetIP   string             `json:"targetIP"`
	HasMalware bool               `json:"hasMalware"`
	Results    []MalwareCheckItem `json:"results"`
	StartTime  string             `json:"startTime"`
	EndTime    string             `json:"endTime"`
}

type MalwareCheckItem struct {
	ID            int    `json:"id"`
	CheckType     int    `json:"checkType"`
	CheckTypeName string `json:"checkTypeName"`
	RiskLevel     int    `json:"riskLevel"`
	RiskName      string `json:"riskName"`
	MatchRule     string `json:"matchRule"`
	FilePath      string `json:"filePath"`
	Description   string `json:"description"`
	FixSuggestion string `json:"fixSuggestion"`
	CheckTime     string `json:"checkTime"`
}

type MalwareResultListReq struct {
	TaskID   int `json:"taskId" form:"taskId"`
	TargetID int `json:"targetId" form:"targetId"`
	Page     int `json:"page" form:"page"`
	Size     int `json:"size" form:"size"`
}

type MalwareResultListResp struct {
	List  []MalwareCheckItem `json:"list"`
	Total int64              `json:"total"`
}

// MalwareTaskListReq 恶意代码检测按批次聚合列表
type MalwareTaskListReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type MalwareTaskListItem struct {
	TaskID         int    `json:"taskId"`
	TargetIP       string `json:"targetIp"`
	TotalFindings  int    `json:"totalFindings"`
	WorstRiskLevel int    `json:"worstRiskLevel"`
	WorstRiskName  string `json:"worstRiskName"`
	CheckTime      string `json:"checkTime"`
}

type MalwareTaskListResp struct {
	List  []MalwareTaskListItem `json:"list"`
	Total int64                 `json:"total"`
}

type DBCheckReq struct {
	TaskID   int    `json:"taskId" form:"taskId" binding:"required"`
	TargetID int    `json:"targetId" form:"targetId"`
	Host     string `json:"host" form:"host" binding:"required"`
	Port     int    `json:"port" form:"port"`
	DBType   int    `json:"dbType" form:"dbType" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password"`
	DBName   string `json:"dbName" form:"dbName"`
}

type DBCheckResp struct {
	TaskID     int           `json:"taskId"`
	TargetIP   string        `json:"targetIP"`
	DBType     int           `json:"dbType"`
	DBTypeName string        `json:"dbTypeName"`
	TotalRules int           `json:"totalRules"`
	PassCount  int           `json:"passCount"`
	FailCount  int           `json:"failCount"`
	Results    []DBCheckItem `json:"results"`
	StartTime  string        `json:"startTime"`
	EndTime    string        `json:"endTime"`
}

type DBCheckItem struct {
	ID              int    `json:"id"`
	RuleID          int    `json:"ruleId"`
	RuleName        string `json:"ruleName"`
	CheckCategory   int    `json:"checkCategory"`
	CategoryName    string `json:"categoryName"`
	CheckResult     int    `json:"checkResult"`
	ResultName      string `json:"resultName"`
	ExpectedValue   string `json:"expectedValue"`
	ActualValue     string `json:"actualValue"`
	RiskLevel       int    `json:"riskLevel"`
	RiskName        string `json:"riskName"`
	FixSuggestion   string `json:"fixSuggestion"`
	RiskDescription string `json:"riskDescription"`
	CheckTime       string `json:"checkTime"`
}

type DBCheckResultListReq struct {
	TaskID   int `json:"taskId" form:"taskId"`
	TargetID int `json:"targetId" form:"targetId"`
	Page     int `json:"page" form:"page"`
	Size     int `json:"size" form:"size"`
}

type DBCheckResultListResp struct {
	List  []DBCheckItem `json:"list"`
	Total int64         `json:"total"`
}

type SensitiveDataScanReq struct {
	TaskID    int    `json:"taskId" form:"taskId" binding:"required"`
	TargetID  int    `json:"targetId" form:"targetId"`
	Host      string `json:"host" form:"host" binding:"required"`
	Port      int    `json:"port" form:"port"`
	DBType    int    `json:"dbType" form:"dbType" binding:"required"`
	Username  string `json:"username" form:"username" binding:"required"`
	Password  string `json:"password" form:"password"`
	DBName    string `json:"dbName" form:"dbName"`
	ScanAllDB bool   `json:"scanAllDb" form:"scanAllDb"`
}

type SensitiveDataScanResp struct {
	TaskID      int                 `json:"taskId"`
	TargetIP    string              `json:"targetIP"`
	DBType      int                 `json:"dbType"`
	DBTypeName  string              `json:"dbTypeName"`
	HighCount   int                 `json:"highCount"`
	MiddleCount int                 `json:"middleCount"`
	LowCount    int                 `json:"lowCount"`
	Results     []SensitiveDataItem `json:"results"`
	StartTime   string              `json:"startTime"`
	EndTime     string              `json:"endTime"`
}

type SensitiveDataItem struct {
	ID            int    `json:"id"`
	DBName        string `json:"dbName"`
	TableName     string `json:"tableName"`
	ColumnName    string `json:"columnName"`
	DataType      int    `json:"dataType"`
	DataTypeName  string `json:"dataTypeName"`
	DataLevel     int    `json:"dataLevel"`
	DataLevelName string `json:"dataLevelName"`
	MatchRule     string `json:"matchRule"`
	SampleData    string `json:"sampleData"`
	CreateTime    string `json:"createTime"`
}

type SensitiveDataListReq struct {
	TaskID   int `json:"taskId" form:"taskId"`
	TargetID int `json:"targetId" form:"targetId"`
	Page     int `json:"page" form:"page"`
	Size     int `json:"size" form:"size"`
}

type SensitiveDataListResp struct {
	List  []SensitiveDataItem `json:"list"`
	Total int64               `json:"total"`
}

type SensitiveDataStatReq struct {
	TaskID int `json:"taskId" form:"taskId"`
}

type SensitiveDataStatResp struct {
	HighCount   int `json:"highCount"`
	MiddleCount int `json:"middleCount"`
	LowCount    int `json:"lowCount"`
	TotalCount  int `json:"totalCount"`
}

// HostTestConnResp 主机 SSH/WinRM 连接测试结果
type HostTestConnResp struct {
	OK            bool   `json:"ok"`
	Message       string `json:"message"`
	TransportName string `json:"transportName,omitempty"`
	Detail        string `json:"detail,omitempty"`
}
