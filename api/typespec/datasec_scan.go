package typespec

// DataSecDBRunReq 数据库基线检查任务创建（与前端 DBSecurity 表单对齐）
type DataSecDBRunReq struct {
	Name       string `json:"name"`
	DBType     int    `json:"dbType"`
	DBHost     string `json:"dbHost"`
	DBPort     int    `json:"dbPort"`
	DBName     string `json:"dbName"`
	DBUser     string `json:"dbUser"`
	DBPassword string `json:"dbPassword"`
}

// DataSecSensitiveRunReq 敏感数据发现任务创建
type DataSecSensitiveRunReq struct {
	Name       string `json:"name"`
	DBType     int    `json:"dbType"`
	DBHost     string `json:"dbHost"`
	DBPort     int    `json:"dbPort"`
	DBName     string `json:"dbName"`
	DBUser     string `json:"dbUser"`
	DBPassword string `json:"dbPassword"`
	DataTypes  []int  `json:"dataTypes"`
	ScanAllDB  bool   `json:"scanAllDb"`
}

type DataSecScanListReq struct {
	Page   int    `form:"page"`
	Size   int    `form:"size"`
	Search string `form:"search"`
}

type DataSecScanDetailReq struct {
	ID string `form:"id" binding:"required"`
}

type DataSecScanRunResp struct {
	ID string `json:"id"`
}

type DataSecDBListItem struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	DBType          int                   `json:"dbType"`
	DBHost          string                `json:"dbHost"`
	DBPort          int                   `json:"dbPort"`
	DBName          string                `json:"dbName"`
	RiskLevel       int                   `json:"riskLevel"`
	Status          int                   `json:"status"`
	CreateTime      string                `json:"createTime"`
	CheckTime       string                `json:"checkTime"`
	CriticalCount   int                   `json:"criticalCount,omitempty"`
	HighRiskCount   int                   `json:"highRiskCount,omitempty"`
	MiddleRiskCount int                   `json:"middleRiskCount,omitempty"`
	LowRiskCount    int                   `json:"lowRiskCount,omitempty"`
	Items           []DataSecDBDetailItem `json:"items,omitempty"`
}

type DataSecDBDetailItem struct {
	Category    int    `json:"category"`
	RiskLevel   int    `json:"riskLevel"`
	Result      string `json:"result"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}

type DataSecDBListResp struct {
	List  []DataSecDBListItem `json:"list"`
	Total int                 `json:"total"`
}

type DataSecSensitiveListItem struct {
	ID           string                         `json:"id"`
	Name         string                         `json:"name"`
	DBType       int                            `json:"dbType"`
	DBHost       string                         `json:"dbHost"`
	DBPort       int                            `json:"dbPort"`
	DBName       string                         `json:"dbName"`
	TotalCount   int                            `json:"totalCount"`
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
