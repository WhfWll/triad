package typespec

type ReportVerifyUploadReq struct {
	Name        string `form:"name" json:"name"`
	Producer    string `form:"producer" json:"producer"`
	ExecuteType string `form:"executeType" json:"executeType"`
	UserId      int    `json:"userId" form:"userId"`
}

type ReportVerifyUploadResp struct {
}

type ReportVerifyTaskListReq struct {
	Search    string `json:"search" form:"search"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
	Risk      int    `json:"risk" form:"risk"`
	Producer  int    `json:"producer" form:"producer"`
	Page      int    `json:"page" form:"page" binding:"required"`
	Size      int    `json:"size" form:"size" binding:"required"`
}

type ReportVerifyTaskListResp struct {
	List  []ReportVerifyTaskListRespItems `json:"list"`
	Total int64                           `json:"total"`
}

type ReportVerifyTaskListRespItems struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	ExecuteType     int    `json:"executeType"`
	ExecuteTypeName string `json:"executeTypeName"`
	Exp             int    `json:"exp"`
	Verify          int    `json:"verify"`
	Failed          int    `json:"failed"`
	Risk            int    `json:"risk"`
	RiskName        string `json:"riskName"`
	Producer        int    `json:"producer"`
	ProducerName    string `json:"producerName"`
	UpdateTime      string `json:"updateTime"`
	Status          int    `json:"status"`
	StatusName      string `json:"statusName"`
}

type ReportVerifyTaskDetailReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type ReportVerifyTaskDetailResp struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	ExecuteType     int    `json:"executeType"`
	ExecuteTypeName string `json:"executeTypeName"`
	Exp             int    `json:"exp"`
	Verify          int    `json:"verify"`
	Failed          int    `json:"failed"`
	User            string `json:"user"`
	Risk            int    `json:"risk"`
	RiskName        string `json:"riskName"`
	Producer        int    `json:"producer"`
	ProducerName    string `json:"producerName"`
	UpdateTime      string `json:"updateTime"`
	Status          int    `json:"status"`
	StatusName      string `json:"statusName"`
	FileName        string `json:"fileName"`
	TargetNumber    string `json:"targetNumber"`
	FileSize        string `json:"fileSize"`
	CreateTime      string `json:"createTime"`
}

type ReportVerifyTargetListReq struct {
	TaskId int    `json:"taskId" form:"taskId" binding:"required"`
	Risk   int    `json:"risk" form:"risk" `
	Page   int    `json:"page" form:"page" binding:"required"`
	Size   int    `json:"size" form:"size" binding:"required"`
	Search string `json:"search" form:"search"`
}

type ReportVerifyTargetListResp struct {
	List  []ReportVerifyTargetRespItems `json:"list"`
	Total int64                         `json:"total"`
}

type ReportVerifyTargetRespItems struct {
	Id         int    `json:"id"`
	TaskId     int    `json:"TaskId"`
	Target     string `json:"name"`
	Os         string `json:"os"`
	Risk       int    `json:"risk"`
	RiskName   string `json:"riskName"`
	Exp        int    `json:"exp"`
	Verify     int    `json:"verify"`
	Failed     int    `json:"failed"`
	UnVerify   int    `json:"unVerify"`
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
}

type ReportVerifyPortListReq struct {
	TaskId int    `json:"taskId" form:"taskId" binding:"required"`
	Page   int    `json:"page" form:"page" binding:"required"`
	Size   int    `json:"size" form:"size" binding:"required"`
	Search string `json:"search" form:"search"`
}

type ReportVerifyPortListResp struct {
	List  []ReportVerifyPortRespItems `json:"list"`
	Total int64                       `json:"total"`
}

type ReportVerifyPortRespItems struct {
	Id        int    `json:"id"`
	TaskId    int    `json:"TaskId"`
	Port      string `json:"port"`
	Scheme    string `json:"scheme"`
	Service   string `json:"service"`
	Component string `json:"component"`
	Target    string `json:"target"`
}

type ReportVerifyVulListReq struct {
	TaskId int    `json:"taskId" form:"taskId" binding:"required"`
	Page   int    `json:"page" form:"page" binding:"required"`
	Size   int    `json:"size" form:"size" binding:"required"`
	Risk   int    `json:"risk" form:"risk" `
	Status int    `json:"status" form:"status"`
	Search string `json:"search" form:"search"`
}

type ReportVerifyVulListResp struct {
	List  []ReportVerifyVulRespItems `json:"list"`
	Total int64                      `json:"total"`
}

type ReportVerifyVulRespItems struct {
	Id         int    `json:"id"`
	TaskId     int    `json:"taskId"`
	Name       string `json:"name"`
	Risk       int    `json:"risk"`
	RiskName   string `json:"riskName"`
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
	Location   string `json:"location"`
}

type ReportVerifyEnumReq struct {
}

type ReportVerifyEnumResp struct {
	ProducerType []GlobalOptionsItemRes `json:"producerType"`
	Status       []GlobalOptionsItemRes `json:"status"`
	Risk         []GlobalOptionsItemRes `json:"risk"`
	VulRisk      []GlobalOptionsItemRes `json:"vulRisk"`
	ExecuteType  []GlobalOptionsItemRes `json:"executeType"`
	VulStatus    []GlobalOptionsItemRes `json:"vulStatus"`
}

type ReportVerifyStatsInfoReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}

type ReportVerifyStatsInfoResp struct {
	AllVul       int `json:"allVul"`       //所有漏洞数
	UnVerify     int `json:"unVerify"`     //未验证漏洞数
	Verify       int `json:"verify"`       //验证成功漏洞数
	Failed       int `json:"failed"`       //验证失败漏洞数
	Exp          int `json:"exp"`          //利用成功漏洞数
	HighVul      int `json:"highVul"`      //高危漏洞数
	MiddleVul    int `json:"middleVul"`    //中危漏洞数
	LowVul       int `json:"lowVul"`       //低危漏洞数
	AllTarget    int `json:"allTarget"`    //所有目标漏洞数
	HighTarget   int `json:"highTarget"`   //高危目标数
	MiddleTarget int `json:"middleTarget"` //中危目标数
	LowTarget    int `json:"lowTarget"`    //低危目标数
	SafeTarget   int `json:"safeTarget"`   //安全目标数
	AliveTarget  int `json:"aliveTarget"`  //存活目标数
}

type ReportVerifyTaskStopReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}

type ReportVerifyTaskStopResp struct {
}

type ReportVerifyTaskDeleteReq struct {
	TaskId string `form:"taskId" json:"taskId" binding:"required"`
}

type ReportVerifyTaskDeleteResp struct {
}

type ReportVerifyTargetDeleteReq struct {
	TargetId string `form:"targetId" json:"targetId" binding:"required"`
}

type ReportVerifyTargetDeleteResp struct {
}

type ReportVerifyVulDeleteReq struct {
	VulId string `form:"vulId" json:"vulId" binding:"required"`
}

type ReportVerifyVulDeleteResp struct {
}

type ReportVerifyVulDetailReq struct {
	VulId int `form:"vulId" json:"vulId" binding:"required"`
}

type ReportVerifyVulDetailResp struct {
	Id    int    `json:"vulId"`
	Name  string `json:"name"`
	Risk  int    `json:"risk"`
	Desc  string `json:"desc"`
	Fix   string `json:"fix"`
	Cve   string `json:"cve"`
	Cnnvd string `json:"cnnvd"`
	Cvss  string `json:"cvss"`
}
