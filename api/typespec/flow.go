package typespec

type FlowTaskEnumResp struct {
	ExpireTime   interface{} `json:"expireTime"`
	NetworkCard  interface{} `json:"networkCard"`
	Status       interface{} `json:"status"`
	VulRiskLevel interface{} `json:"vulRiskLevel"`
	CredPattern  interface{} `json:"credPattern"`
	FuzzParam    interface{} `json:"fuzzParam"`
	Response     struct {
		JsonKeyword   string `json:"jsonKeyword"`   //json关键字
		NoJsonSwitch  bool   `json:"NoJsonSwitch"`  //非json关键字开关
		NoJsonKeyword string `json:"noJsonKeyword"` //非json关键字
	} `json:"response"`
	VulName interface{} `json:"vulName"`
}

type FlowTaskListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}

type FlowTaskListResp struct {
	Total int64                   `json:"total"`
	List  []FlowTaskListRespItems `json:"list"`
}

type FlowTaskListRespItems struct {
	Id            int    `json:"id"`
	TaskName      string `json:"taskName"`
	RiskLevel     int    `json:"riskLevel"`
	RiskLevelName string `json:"riskLevelName"`
	Status        int    `json:"status"`
	StatusName    string `json:"statusName"`
	VulNum        []int  `json:"vulNum"`
	UpdateTime    string `json:"updateTime"`
}

type FlowTaskDelReq struct {
	FlowTaskIds string `form:"flowTaskIds" json:"flowTaskIds" binding:"required"`
}

type FlowTaskDelResp struct {
}

type FlowTaskAddReq struct {
	TaskName    string `form:"taskName" json:"taskName" binding:"required"`
	Port        string `form:"port" json:"port" binding:"required"`
	TargetUrl   string `form:"targetUrl" json:"targetUrl" binding:"required"`
	NetworkCard string `form:"networkCard" json:"networkCard" binding:"required"`
	ExpireTime  int    `form:"expireTime" json:"expireTime" binding:"required"`
	UserId      int    `form:"userId" json:"userId" binding:"required"`
	OtherConfig string `form:"otherConfig" json:"otherConfig"`
	VulConfig   string `form:"vulConfig" json:"vulConfig"`
}

type OtherConfig struct {
	WaitCred struct {
		Pattern int    `json:"pattern"`
		Value   string `json:"value"`
	} `json:"waitCred"`
	FuzzParam struct {
		Character string `json:"character"`
		Number    string `json:"number"`
	}
	FuzzDict struct {
		Character string `json:"character"`
		Number    string `json:"number"`
	}
	Response struct {
		JsonKeyword   string `json:"jsonKeyword"`
		NoJsonSwitch  bool   `json:"noJsonSwitch"`
		NoJsonKeyword string `json:"noJsonKeyword"`
	}
}

type FlowTaskAddResp struct {
}

type ChangeFlowTaskStatusReq struct {
	FlowTaskId int    `form:"flowTaskId" json:"flowTaskId" binding:"required"`
	Operate    string `form:"operate" json:"operate" binding:"required"`
	UserId     int    `form:"userId" json:"userId" binding:"required"`
}

type ChangeFlowTaskStatusResp struct {
}

type FlowTaskInfoReq struct {
	FlowTaskId int `form:"flowTaskId" json:"flowTaskId" binding:"required"`
}

type FlowTaskInfoResp struct {
	Id             int         `json:"id"`
	TaskName       string      `json:"taskName"`
	NetworkCard    string      `json:"networkCard"`
	Port           string      `json:"port"`
	CreateTime     string      `json:"createTime"`
	UserId         int         `json:"userId"`
	UserName       string      `json:"userName"`
	ExpireTime     int         `json:"expireTime"`
	ExpireTimeName string      `json:"expireTimeName"`
	RiskLevel      int         `json:"riskLevel"`
	RiskLevelName  string      `json:"riskLevelName"`
	Status         int         `json:"status"`
	StatusName     string      `json:"statusName"`
	Target         string      `json:"target"`
	VulNum         []int       `json:"vulNum"`
	TargetNum      []int       `json:"targetNum"`
	OtherConfig    OtherConfig `form:"otherConfig" json:"otherConfig"`
	VulConfig      string      `form:"vulConfig" json:"vulConfig"`
	VulConfigZh    string      `form:"vulConfig_zh" json:"vulConfig_zh"`
}

type HttpsCertReq struct {
	FlowTaskId int `form:"flowTaskId" json:"flowTaskId" binding:"required"`
}

type HttpsCertResp struct {
	Cert string `json:"cert"`
}

type FlowTaskStatusReq struct {
	FlowTaskId int `form:"flowTaskId" json:"flowTaskId" binding:"required"`
}

type FlowTaskStatusResp struct {
	Id         int    `json:"id"`
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
}

type FlowRiskListReq struct {
	FlowTaskId int    `form:"flowTaskId" json:"flowTaskId" binding:"required"`
	Page       int    `form:"page" json:"page" binding:"required"`
	Size       int    `form:"size" json:"size" binding:"required"`
	Search     string `form:"search" json:"search"`
	RiskLevel  int    `form:"riskLevel" json:"riskLevel"`
}

type FlowRiskListResp struct {
	Total int64                   `json:"total"`
	List  []FlowRiskListRespItems `json:"list"`
}

type FlowRiskListRespItems struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Ip            string `json:"ip"`
	Host          string `json:"host"`
	RiskLevel     int    `json:"riskLevel"`
	RiskLevelName string `json:"riskLevelName"`
}

type FlowRiskInfoReq struct {
	FlowRiskId int `form:"flowRiskId" json:"flowRiskId" binding:"required"`
}
type FlowRiskInfoResp struct {
	Id                 int    `json:"id"`
	Ip                 string `json:"ip"`
	Port               string `json:"port"`
	Host               string `json:"host"`
	RiskLevel          int    `json:"riskLevel"`
	RiskLevelName      string `json:"riskLevelName"`
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`
	Hash               string `json:"hash"`
	RiskTypeVerbose    string `json:"riskTypeVerbose"`
	Origin             string `json:"origin"`
	Status             string `json:"status"`
	Url                string `json:"url"`
	Parameter          string `json:"parameter"`
	Request            string `json:"request"`
	Response           string `json:"response"`
	Detail             string `json:"detail"`
	Token              string `json:"token"`
	Payload            string `json:"payload"`
	PayloadSuccessFlag string `json:"payload_success_flag"`
}

type FlowRiskDelReq struct {
	FlowRiskIds string `form:"flowRiskIds" json:"flowRiskIds" binding:"required"`
}

type FlowRiskDelResp struct {
}

type FlowBaseListReq struct {
	FlowTaskId int    `form:"flowTaskId" json:"flowTaskId" binding:"required"`
	Search     string `form:"search" json:"search"`
	Page       int    `form:"page" json:"page" binding:"required"`
	Size       int    `form:"size" json:"size" binding:"required"`
}

type FlowBaseListResp struct {
	Total int64                   `json:"total"`
	List  []FlowBaseListRespItems `json:"list"`
}
type FlowBaseListRespItems struct {
	Id              int    `json:"id"`
	Ip              string `json:"ip"`
	Method          string `json:"method"`
	CreateTime      string `json:"createTime"`
	RespContentType string `json:"respContentType"`
	RespCode        string `json:"respCode"`
	RespTitle       string `json:"respTitle"`
	Tags            string `json:"tags"`
	Host            string `json:"host"`
	Url             string `json:"url"`
}

type FlowBaseInfoReq struct {
	FlowBaseId int `form:"flowBaseId" json:"flowBaseId" binding:"required"`
}
type FlowBaseInfoResp struct {
	ReqHeader  string `json:"reqHeader"`
	RespHeader string `json:"respHeader"`
}

type FlowBaseDelReq struct {
	FlowBaseIds string `form:"flowBaseIds" json:"flowBaseIds" binding:"required"`
}
type FlowBaseDelResp struct {
}

type FlowLogInfoReq struct {
	FlowTaskId int    `form:"flowTaskId" json:"flowTaskId" binding:"required"`
	Search     string `form:"search" json:"search"`
	Page       int    `form:"page" json:"page" binding:"required"`
	Size       int    `form:"size" json:"size" binding:"required"`
}
type FlowLogInfoResp struct {
	Total int64                  `json:"total"`
	List  []FlowLogInfoRespItems `json:"list"`
}
type FlowLogInfoRespItems struct {
	Id         int    `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
}

type FlowLogDelReq struct {
	FlowTaskId int `form:"flowTaskId" json:"flowTaskId" binding:"required"`
}
type FlowLogDelResp struct {
}

type FlowTaskEditReq struct {
	TaskName    string `form:"taskName" json:"taskName" binding:"required"`
	Port        string `form:"port" json:"port" binding:"required"`
	TargetUrl   string `form:"targetUrl" json:"targetUrl" binding:"required"`
	NetworkCard string `form:"networkCard" json:"networkCard" binding:"required"`
	ExpireTime  int    `form:"expireTime" json:"expireTime" binding:"required"`
	UserId      int    `form:"userId" json:"userId" binding:"required"`
	OtherConfig string `form:"otherConfig" json:"otherConfig"`
	VulConfig   string `form:"vulConfig" json:"vulConfig"`
	FlowTaskId  int    `form:"flowTaskId" json:"flowTaskId"`
}

type FlowTaskEditResp struct {
}

type FlowTaskExportReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}

type FlowTaskExportResp struct {
	List []FlowBaseListRespItems `json:"list"`
}
