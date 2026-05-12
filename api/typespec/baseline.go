package typespec

type BaselineCheckReq struct {
	TaskID   int    `json:"taskId" form:"taskId" binding:"required"`
	TargetID int    `json:"targetId" form:"targetId"`
	Host     string `json:"host" form:"host" binding:"required"`
	Port     int    `json:"port" form:"port"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"`
	OSType   int    `json:"osType" form:"osType" binding:"required"`
}

type BaselineCheckResp struct {
	TaskID     int                  `json:"taskId"`
	TargetIP   string               `json:"targetIP"`
	OSType     int                  `json:"osType"`
	OSTypeName string               `json:"osTypeName"`
	TotalRules int                  `json:"totalRules"`
	PassCount  int                  `json:"passCount"`
	FailCount  int                  `json:"failCount"`
	ErrorCount int                  `json:"errorCount"`
	Results    []BaselineCheckItem  `json:"results"`
	StartTime  string               `json:"startTime"`
	EndTime    string               `json:"endTime"`
}

type BaselineCheckItem struct {
	ID              int    `json:"id"`
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
	TotalRules int `json:"totalRules"`
	PassCount  int `json:"passCount"`
	FailCount  int `json:"failCount"`
	PassRate   float64 `json:"passRate"`
}

type MalwareScanReq struct {
	TaskID   int    `json:"taskId" form:"taskId" binding:"required"`
	TargetID int    `json:"targetId" form:"targetId"`
	Host     string `json:"host" form:"host" binding:"required"`
	Port     int    `json:"port" form:"port"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"`
	OSType   int    `json:"osType" form:"osType"`
}

type MalwareScanResp struct {
	TaskID     int               `json:"taskId"`
	TargetIP   string            `json:"targetIP"`
	HasMalware bool              `json:"hasMalware"`
	Results    []MalwareCheckItem `json:"results"`
	StartTime  string            `json:"startTime"`
	EndTime    string            `json:"endTime"`
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
	TaskID      int                     `json:"taskId"`
	TargetIP    string                  `json:"targetIP"`
	DBType      int                     `json:"dbType"`
	DBTypeName  string                  `json:"dbTypeName"`
	HighCount   int                     `json:"highCount"`
	MiddleCount int                     `json:"middleCount"`
	LowCount    int                     `json:"lowCount"`
	Results     []SensitiveDataItem     `json:"results"`
	StartTime   string                  `json:"startTime"`
	EndTime     string                  `json:"endTime"`
}

type SensitiveDataItem struct {
	ID           int    `json:"id"`
	DBName       string `json:"dbName"`
	TableName    string `json:"tableName"`
	ColumnName   string `json:"columnName"`
	DataType     int    `json:"dataType"`
	DataTypeName string `json:"dataTypeName"`
	DataLevel    int    `json:"dataLevel"`
	DataLevelName string `json:"dataLevelName"`
	MatchRule    string `json:"matchRule"`
	SampleData   string `json:"sampleData"`
	CreateTime   string `json:"createTime"`
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
