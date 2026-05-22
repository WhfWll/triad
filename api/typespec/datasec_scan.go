package typespec

// DataSecDBTargetInput 单个数据库扫描目标（多目标任务）
type DataSecDBTargetInput struct {
	DBHost     string  `json:"dbHost"`
	DBPort     FlexInt `json:"dbPort"`
	DBName     string  `json:"dbName"`
	DBUser     string  `json:"dbUser"`
	DBPassword string  `json:"dbPassword"`
}

// DataSecDBRunReq 数据库基线检查任务创建（与前端 DBSecurity 表单对齐）
type DataSecDBRunReq struct {
	Name             string                 `json:"name"`
	DBType           FlexInt                `json:"dbType"`
	Targets          []DataSecDBTargetInput `json:"targets"`
	LibraryTargetIDs []int                  `json:"libraryTargetIds"`
	ScanSensitive    bool                   `json:"scanSensitive"`
	DataTypes        []int                  `json:"dataTypes,omitempty"`
	ScanAllDB        bool                   `json:"scanAllDb,omitempty"`
	// 兼容旧版单目标字段
	DBHost     string  `json:"dbHost"`
	DBPort     FlexInt `json:"dbPort"`
	DBName     string  `json:"dbName"`
	DBUser     string  `json:"dbUser"`
	DBPassword string  `json:"dbPassword"`
}

// DataSecDBTestConnReq 数据库连接测试（字段与创建任务一致）
type DataSecDBTestConnReq struct {
	DBType     FlexInt `json:"dbType"`
	DBHost     string  `json:"dbHost"`
	DBPort     FlexInt `json:"dbPort"`
	DBName     string  `json:"dbName"`
	DBUser     string  `json:"dbUser"`
	DBPassword string  `json:"dbPassword"`
}

type DataSecDBTestConnResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// DataSecSensitiveRunReq 敏感数据发现任务创建
type DataSecSensitiveRunReq struct {
	Name             string                 `json:"name"`
	DBType           FlexInt                `json:"dbType"`
	Targets          []DataSecDBTargetInput `json:"targets"`
	LibraryTargetIDs []int                  `json:"libraryTargetIds"`
	DataTypes        []int                  `json:"dataTypes"`
	ScanAllDB        bool                   `json:"scanAllDb"`
	// 兼容旧版单目标字段
	DBHost     string  `json:"dbHost"`
	DBPort     FlexInt `json:"dbPort"`
	DBName     string  `json:"dbName"`
	DBUser     string  `json:"dbUser"`
	DBPassword string  `json:"dbPassword"`
}

type DataSecScanListReq struct {
	Page   int    `form:"page"`
	Size   int    `form:"size"`
	Search string `form:"search"`
}

type DataSecScanDetailReq struct {
	ID       string `form:"id" binding:"required"`
	TargetID int    `form:"targetId"`
}

type DataSecScanRunResp struct {
	ID string `json:"id"`
}

// DataSecTargetItem 任务内单个数据库目标摘要
type DataSecTargetItem struct {
	ID              int    `json:"id"`
	TargetURL       string `json:"targetUrl"`
	DBType          int    `json:"dbType"`
	DBHost          string `json:"dbHost"`
	DBPort          int    `json:"dbPort"`
	DBName          string `json:"dbName"`
	Status          int    `json:"status"`
	RiskLevel       int    `json:"riskLevel,omitempty"`
	DbVersion       string `json:"dbVersion,omitempty"`
	BaselineTotal   int    `json:"baselineTotal,omitempty"`
	BaselinePass    int    `json:"baselinePass,omitempty"`
	BaselineFail    int    `json:"baselineFail,omitempty"`
	BaselineError   int    `json:"baselineError,omitempty"`
	CveMatchCount   int    `json:"cveMatchCount,omitempty"`
	CriticalCount   int    `json:"criticalCount,omitempty"`
	HighRiskCount   int    `json:"highRiskCount,omitempty"`
	MiddleRiskCount int    `json:"middleRiskCount,omitempty"`
	LowRiskCount    int    `json:"lowRiskCount,omitempty"`
	TotalCount      int    `json:"totalCount,omitempty"`
	HighCount       int    `json:"highCount,omitempty"`
	MediumCount     int    `json:"mediumCount,omitempty"`
	LowCount        int    `json:"lowCount,omitempty"`
	ErrorMessage    string `json:"errorMessage,omitempty"`
}

type DataSecDBListItem struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	DBType          int                   `json:"dbType"`
	DBHost          string                `json:"dbHost"`
	DBPort          int                   `json:"dbPort"`
	DBName          string                `json:"dbName"`
	TargetSummary   string                `json:"targetSummary,omitempty"`
	TargetCount     int                   `json:"targetCount,omitempty"`
	Targets         []DataSecTargetItem   `json:"targets,omitempty"`
	RiskLevel       int                   `json:"riskLevel"`
	Status          int                   `json:"status"`
	CreateTime      string                `json:"createTime"`
	CheckTime       string                `json:"checkTime"`
	CriticalCount   int                   `json:"criticalCount,omitempty"`
	HighRiskCount   int                   `json:"highRiskCount,omitempty"`
	MiddleRiskCount int                   `json:"middleRiskCount,omitempty"`
	LowRiskCount    int                   `json:"lowRiskCount,omitempty"`
	BaselineTotal   int                   `json:"baselineTotal,omitempty"`
	BaselineFail    int                   `json:"baselineFail,omitempty"`
	CveMatchCount   int                   `json:"cveMatchCount,omitempty"`
	ScanSensitive   bool                  `json:"scanSensitive,omitempty"`
	TotalCount      int                   `json:"totalCount,omitempty"`
	HighCount       int                   `json:"highCount,omitempty"`
	MediumCount     int                   `json:"mediumCount,omitempty"`
	LowCount        int                   `json:"lowCount,omitempty"`
	Items           []DataSecDBDetailItem `json:"items,omitempty"`
	TypeStats       []DataSecSensitiveTypeStatItem `json:"typeStats,omitempty"`
	SensitiveItems  []DataSecSensitiveDetailItem   `json:"sensitiveItems,omitempty"`
}

type DataSecDBDetailItem struct {
	TargetID    int    `json:"targetId,omitempty"`
	Category    int    `json:"category"`
	RiskLevel   int    `json:"riskLevel"`
	Result      string `json:"result"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
	RuleName      string `json:"ruleName,omitempty"`
	ActualValue   string `json:"actualValue,omitempty"`
	ExpectedValue string `json:"expectedValue,omitempty"`
	IsCve         bool   `json:"isCve,omitempty"`
}

type DataSecDBListResp struct {
	List  []DataSecDBListItem `json:"list"`
	Total int                 `json:"total"`
}

type DataSecSensitiveListItem struct {
	ID            string                         `json:"id"`
	Name          string                         `json:"name"`
	DBType        int                            `json:"dbType"`
	DBHost        string                         `json:"dbHost"`
	DBPort        int                            `json:"dbPort"`
	DBName        string                         `json:"dbName"`
	TargetSummary string                         `json:"targetSummary,omitempty"`
	TargetCount   int                            `json:"targetCount,omitempty"`
	Targets       []DataSecTargetItem            `json:"targets,omitempty"`
	TotalCount    int                            `json:"totalCount"`
	HighCount    int                            `json:"highCount"`
	MediumCount  int                            `json:"mediumCount"`
	LowCount     int                            `json:"lowCount"`
	Status       int                            `json:"status"`
	CreateTime   string                         `json:"createTime"`
	ScanTime     string                         `json:"scanTime"`
	TypeStats    []DataSecSensitiveTypeStatItem `json:"typeStats,omitempty"`
	Items        []DataSecSensitiveDetailItem   `json:"items,omitempty"`
}

type DataSecSensitiveTypeStatItem struct {
	DataType int `json:"dataType"`
	Count    int `json:"count"`
}

type DataSecSensitiveDetailItem struct {
	TargetID         int    `json:"targetId,omitempty"`
	TableName        string `json:"tableName"`
	ColumnName       string `json:"columnName"`
	DataType         int    `json:"dataType"`
	SensitivityLevel int    `json:"sensitivityLevel"`
	SampleData       string `json:"sampleData"`
	Count            int    `json:"count"`
}

type DataSecSensitiveListResp struct {
	List  []DataSecSensitiveListItem `json:"list"`
	Total int                        `json:"total"`
}
