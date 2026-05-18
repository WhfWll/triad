package typespec

type VulnScanCveReq struct {
	TaskID        int    `json:"taskId"`
	TargetID      int    `json:"targetId"`
	Host          string `json:"host" binding:"required"`
	Port          int    `json:"port"`
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password"`
	Key           string `json:"key"`
	OSType        int    `json:"osType"`
	Transport     int    `json:"transport"`
	WinRMUseHttps bool   `json:"winrmUseHttps"`
}

type VulnScanCveResp struct {
	TaskID       int               `json:"taskId"`
	TargetIP     string            `json:"targetIp"`
	OSType       int               `json:"osType"`
	OSTypeName   string            `json:"osTypeName"`
	Packages     int               `json:"packages"`
	MatchedVulns int               `json:"matchedVulns"`
	Critical     int               `json:"critical"`
	High         int               `json:"high"`
	Medium       int               `json:"medium"`
	Low          int               `json:"low"`
	StartTime    string            `json:"startTime"`
	EndTime      string            `json:"endTime"`
	Results      []VulnScanCveItem `json:"results"`
}

type VulnScanCveItem struct {
	PackageName    string `json:"packageName"`
	PackageVersion string `json:"packageVersion"`
	Cve            string `json:"cve"`
	Title          string `json:"title"`
	Severity       string `json:"severity"`
	RiskLevel      int    `json:"riskLevel"`
}

type VulnScanCveTarget struct {
	Host          string `json:"host" binding:"required"`
	Port          int    `json:"port"`
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password"`
	Key           string `json:"key"`
	OSType        int    `json:"osType"`
	Transport     int    `json:"transport"`
	WinRMUseHttps bool   `json:"winrmUseHttps"`
}

type VulnScanCveBatchReq struct {
	TaskID   int                  `json:"taskId"`
	Targets  []VulnScanCveTarget `json:"targets" binding:"required"`
}

type VulnScanCveBatchResp struct {
	TaskID int    `json:"taskId"`
	Status string `json:"status"`
}

type VulnScanCveTargetResult struct {
	TargetIP     string            `json:"targetIp"`
	OSType       int               `json:"osType"`
	Packages     int               `json:"packages"`
	MatchedVulns int               `json:"matchedVulns"`
	Critical     int               `json:"critical"`
	High         int               `json:"high"`
	Medium       int               `json:"medium"`
	Low          int               `json:"low"`
	StartTime    string            `json:"startTime"`
	EndTime      string            `json:"endTime"`
	Error        string            `json:"error,omitempty"`
	Results      []VulnScanCveItem `json:"results"`
}

type VulnScanCveProgressReq struct {
	TaskID int `form:"taskId" json:"taskId" binding:"required"`
}

type VulnScanCveProgressResp struct {
	TaskID   int                       `json:"taskId"`
	Status   string                    `json:"status"`
	Progress int                       `json:"progress"`
	Total    int                       `json:"total"`
	Results  []VulnScanCveTargetResult `json:"results"`
	Errors   []string                  `json:"errors"`
}