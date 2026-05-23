package typespec

type DatasecDBTargetListReq struct {
	Page   int    `form:"page"`
	Size   int    `form:"size"`
	DBType int    `form:"dbType"`
	Search string `form:"search"`
}

type DatasecDBTargetListItem struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	GroupName  string `json:"groupName"`
	DBType     int    `json:"dbType"`
	DBHost     string `json:"dbHost"`
	DBPort     int    `json:"dbPort"`
	DBName     string `json:"dbName"`
	DBUser     string `json:"dbUser"`
	HasPassword bool  `json:"hasPassword"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type DatasecDBTargetListResp struct {
	List  []DatasecDBTargetListItem `json:"list"`
	Total int                       `json:"total"`
}

type DatasecDBTargetSaveReq struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	GroupName  string `json:"groupName"`
	DBType     FlexInt `json:"dbType"`
	DBHost     string `json:"dbHost"`
	DBPort     FlexInt `json:"dbPort"`
	DBName     string `json:"dbName"`
	DBUser     string `json:"dbUser"`
	DBPassword string `json:"dbPassword"`
	Remark     string `json:"remark"`
}

type DatasecDBTargetDeleteReq struct {
	ID int `form:"id" binding:"required"`
}

type DatasecDBTargetImportReq struct {
	Items []DatasecDBTargetImportItem `json:"items"`
}

type DatasecDBTargetImportItem struct {
	Name       string  `json:"name"`
	GroupName  string  `json:"groupName"`
	DBType     FlexInt `json:"dbType"`
	DBHost     string  `json:"dbHost"`
	DBPort     FlexInt `json:"dbPort"`
	DBName     string  `json:"dbName"`
	DBUser     string  `json:"dbUser"`
	DBPassword string  `json:"dbPassword"`
	Remark     string  `json:"remark"`
}

type DatasecDBTargetExportResp struct {
	Version int                         `json:"version"`
	Items   []DatasecDBTargetExportItem `json:"items"`
}

type DatasecDBTargetExportItem struct {
	Name       string `json:"name"`
	GroupName  string `json:"groupName"`
	DBType     int    `json:"dbType"`
	DBHost     string `json:"dbHost"`
	DBPort     int    `json:"dbPort"`
	DBName     string `json:"dbName"`
	DBUser     string `json:"dbUser"`
	DBPassword string `json:"dbPassword,omitempty"`
	Remark     string `json:"remark"`
}

type DatasecTaskCloneTargetsReq struct {
	ID   string `form:"id" binding:"required"`
	Kind string `form:"kind"`
}

type DatasecTaskCloneTargetsResp struct {
	DBType        int                    `json:"dbType"`
	Targets       []DataSecDBTargetInput `json:"targets"`
	ScanSensitive bool                   `json:"scanSensitive,omitempty"`
	DataTypes     []int                  `json:"dataTypes,omitempty"`
	ScanAllDb     bool                   `json:"scanAllDb,omitempty"`
}

type DatasecTaskRerunReq struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type DatasecTaskDeleteReq struct {
	ID   string `form:"id" binding:"required"`
	Kind string `form:"kind"`
}

type DatasecDBTargetTestReq struct {
	ID         int     `json:"id"`
	DBType     FlexInt `json:"dbType"`
	DBHost     string  `json:"dbHost"`
	DBPort     FlexInt `json:"dbPort"`
	DBName     string  `json:"dbName"`
	DBUser     string  `json:"dbUser"`
	DBPassword string  `json:"dbPassword"`
}

type DatasecDBTargetBatchTestReq struct {
	IDs []int `json:"ids"`
}

type DatasecDBTargetBatchTestItem struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	DBHost  string `json:"dbHost"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type DatasecDBTargetBatchTestResp struct {
	Total   int                            `json:"total"`
	OK      int                            `json:"ok"`
	Fail    int                            `json:"fail"`
	Results []DatasecDBTargetBatchTestItem `json:"results"`
}
