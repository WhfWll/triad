package typespec

type SecurityReportGenerateReq struct {
	Module string `json:"module" form:"module" binding:"required"`
	TaskID int    `json:"taskId" form:"taskId" binding:"required"`
}

type SecurityReportGenerateResp struct {
	ID int `json:"id"`
}

type SecurityReportListReq struct {
	Page int `json:"page" form:"page" binding:"required"`
	Size int `json:"size" form:"size" binding:"required"`
}

type SecurityReportListItem struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Module     string `json:"module"`
	ModuleName string `json:"moduleName"`
	TaskID     int    `json:"taskId"`
	TaskName   string `json:"taskName"`
	CreateTime string `json:"createTime"`
}

type SecurityReportListResp struct {
	List  []SecurityReportListItem `json:"list"`
	Total int64                    `json:"total"`
}

type SecurityReportDetailReq struct {
	ID int `json:"id" form:"id" binding:"required"`
}

type SecurityReportDetailResp struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Module  string `json:"module"`
	TaskID  int    `json:"taskId"`
	Content string `json:"content"`
}

type SecurityReportDeleteReq struct {
	ID int `json:"id" form:"id" binding:"required"`
}
