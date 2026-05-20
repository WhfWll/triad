package typespec

// AppSecScanRunReq 应用安全扫描任务创建（与前端 scanConfig 字段对齐）
type AppSecScanRunReq struct {
	Name       string `json:"name"`
	Target     string `json:"target"`
	TargetURL  string `json:"targetUrl"`
	AppType    int    `json:"appType"`
	Strategy   string `json:"strategy"`
	TestMode   string `json:"testMode"`
	SafeTest   bool   `json:"safeTest"`
	VulExploit bool   `json:"vulExploit"`
	TestIntensity int `json:"testIntensity"`
	VulIdsConfig  []int `json:"vulIdsConfig"`

	WebsiteLogin map[string]interface{} `json:"websiteLogin"`
	WebCrawler map[string]interface{} `json:"webCrawler"`
	PortScan   map[string]interface{} `json:"portScan"`
	Proxy      map[string]interface{} `json:"proxy"`
	WebPathScan map[string]interface{} `json:"webPathScan"`
	WeakPass   map[string]interface{} `json:"weakPass"`
}

type AppSecScanListReq struct {
	Page   int    `form:"page"`
	Size   int    `form:"size"`
	Search string `form:"search"`
}

type AppSecScanDetailReq struct {
	ID string `form:"id" binding:"required"`
}

type AppSecTargetItem struct {
	ID              int    `json:"id"`
	TargetURL       string `json:"targetUrl"`
	Status          int    `json:"status"`
	PageCount       int    `json:"pageCount"`
	VulnCount       int    `json:"vulnCount"`
	CriticalCount   int    `json:"criticalCount"`
	HighRiskCount   int    `json:"highRiskCount"`
	MiddleRiskCount int    `json:"middleRiskCount"`
	LowRiskCount    int    `json:"lowRiskCount"`
}

type AppSecVulnItem struct {
	TargetID    int    `json:"targetId"`
	TargetURL   string `json:"targetUrl,omitempty"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	RiskLevel   int    `json:"riskLevel"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Payload     string `json:"payload"`
	Method      string `json:"method,omitempty"`
	Request     string `json:"request,omitempty"`
	Response    string `json:"response,omitempty"`
	Suggestion  string `json:"suggestion,omitempty"`
}

type AppSecPageItem struct {
	TargetID int    `json:"targetId,omitempty"`
	URL      string `json:"url"`
	Name     string `json:"name"`
}

type AppSecTaskItem struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	TargetURL       string                 `json:"targetUrl"`
	TargetSummary   string                 `json:"targetSummary,omitempty"`
	TargetCount     int                    `json:"targetCount,omitempty"`
	TargetID        int                    `json:"targetId,omitempty"`
	Targets         []AppSecTargetItem     `json:"targets,omitempty"`
	StrategyID      string                 `json:"strategyId"`
	AppType         int                    `json:"appType,omitempty"`
	ScanConfig      map[string]interface{} `json:"scanConfig,omitempty"`
	Status          int                    `json:"status"`
	RiskLevel       int              `json:"riskLevel"`
	PageCount       int              `json:"pageCount"`
	VulnCount       int              `json:"vulnCount"`
	CriticalCount   int              `json:"criticalCount"`
	HighRiskCount   int              `json:"highRiskCount"`
	MiddleRiskCount int              `json:"middleRiskCount"`
	LowRiskCount    int              `json:"lowRiskCount"`
	CreateTime      string           `json:"createTime"`
	ScanTime        string           `json:"scanTime"`
	ErrorMessage    string           `json:"errorMessage,omitempty"`
	Vulns           []AppSecVulnItem `json:"vulns"`
	Pages           []AppSecPageItem `json:"pages"`
}

type AppSecScanListResp struct {
	List  []AppSecTaskItem `json:"list"`
	Total int              `json:"total"`
}

type AppSecScanRunResp struct {
	ID string `json:"id"`
}
