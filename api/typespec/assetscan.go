package typespec

type ScanEnumsResp struct {
	ExecuteTypeSave    interface{} `json:"executeTypeSave"`
	ExecuteTypeSearch  interface{} `json:"executeTypeSearch"`
	ExecCycleType      interface{} `json:"execCycleType"`
	ExecCycleTypeWeek  interface{} `json:"execCycleTypeWeek"`
	ExecCycleTypeMonth interface{} `json:"execCycleTypeMonth"`
	PortScanType       interface{} `json:"portScanType"`
	TaskStatus         interface{} `json:"taskStatus"`
	AssetChangeType    interface{} `json:"assetChangeType"`
	AssetIsLiveType    interface{} `json:"assetIsLiveType"`
	AssetStatus        interface{} `json:"assetStatus"`
}

type ScanAddTaskReq struct {
	ActivityId    int                       `form:"activityId" json:"activityId"`
	Name          string                    `form:"name" json:"name" binding:"required"`
	IpRange       string                    `form:"ipRange" json:"ipRange" binding:"required"`
	IpExcludRange string                    `form:"ipExcludRange" json:"ipExcludRange"`
	PortType      string                    `form:"portType" json:"portType" binding:"required"`
	PortRange     string                    `form:"portRange" json:"portRange" binding:"required"`
	UserId        int                       `form:"userId" json:"userId"`
	ExecuteType   int                       `form:"executeType" json:"executeType" binding:"required"`
	ExecuteJson   ScanAddTaskReqExecuteJson `form:"executeJson" json:"executeJson"`
}

type ScanAddTaskReqExecuteJson struct {
	ExecCycleType  int    `form:"execCycleType" json:"execCycleType"`
	ExecCycleValue int    `form:"execCycleValue" json:"execCycleValue"`
	ExecCycleHour  string `form:"execCycleHour" json:"execCycleHour"`
	StartTime      string `form:"startTime" json:"startTime"`
	EndTime        string `form:"endTime" json:"endTime"`
}

type ScanAddTaskResp struct{}

type ScanDelTaskReq struct {
	TaskIds string `form:"taskIds" json:"taskIds" binding:"required"`
}

type ScanDelTaskResp struct{}

type ScanStopTaskReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}

type ScanStopTaskResp struct{}

type GetScanTaskInfoReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}

type GetScanTaskInfoResp struct {
	Id              int         `json:"id"`
	ActivityId      int         `form:"activityId" json:"activityId"`
	Name            string      `form:"name" json:"name"`
	IpRange         string      `form:"ipRange" json:"ipRange"`
	IpExcludRange   string      `form:"ipExcludRange" json:"ipExcludRange"`
	PortType        string      `form:"portType" json:"portType"`
	PortRange       string      `form:"portRange" json:"portRange"`
	UserId          int         `form:"userId" json:"userId"`
	ExecuteType     int         `form:"executeType" json:"executeType"`
	ExecuteTypeName string      `json:"executeTypeName"`
	ExecuteJson     interface{} `form:"executeJson" json:"executeJson"`
	Status          int         `json:"status"`
	StatusName      string      `json:"statusName"`
	CreateTime      string      `json:"createTime"`
	AssetChanges    interface{} `json:"assetChanges"`
}

type AssetChanges struct {
	LiveIp     int64 `json:"liveIp"`
	AddIp      int64 `json:"addIp"`
	ReduceIp   int64 `json:"reduceIp"`
	PortIp     int64 `json:"portIp"`
	ServiceIp  int64 `json:"serviceIp"`
	AssemblyIp int64 `json:"assemblyIp"`
}

type ScanTaskListReq struct {
	ActivityId  int    `form:"activityId" json:"activityId"`
	Search      string `form:"search" json:"search"`
	ExecuteType string `form:"executeType" json:"executeType"`
	CreateTime  string `form:"createTime" json:"createTime"`
	Status      string `form:"status" json:"status"`
	Page        int    `form:"page" json:"page" binding:"required"`
	Size        int    `form:"size" json:"size" binding:"required"`
}

type ScanTaskListResp struct {
	Total int64                  `json:"total"`
	List  []ScanTaskListRespList `json:"list"`
}

type ScanTaskListRespList struct {
	Id              int         `json:"id"`
	Name            string      `json:"name"`
	ExecuteType     int         `json:"executeType"`
	ExecuteTypeName string      `json:"executeTypeName"`
	StatResult      interface{} `json:"statResult"`
	Status          int         `json:"status"`
	StatusName      string      `json:"statusName"`
	CreateTime      string      `json:"createTime"`
}

type ScanTaskExportReq struct {
	TaskIds string `form:"taskIds" json:"taskIds" binding:"required"`
}

type ScanTaskIpListReq struct {
	TaskId           int    `form:"taskId" json:"taskId" binding:"required"`
	Search           string `form:"search" json:"search"`
	IsLive           int    `form:"isLive" json:"isLive"`
	Status           string `form:"status" json:"status"`
	AssetChangesType int    `form:"assetChangesType" json:"assetChangesType"`
	CreateTime       string `form:"createTime" json:"createTime"`
	Page             int    `form:"page" json:"page" binding:"required"`
	Size             int    `form:"size" json:"size" binding:"required"`
}
type ScanTaskIpListResp struct {
	Total int64                    `json:"total"`
	List  []ScanTaskIpListRespList `json:"list"`
}
type ScanTaskIpListRespList struct {
	Id                   int    `json:"id"`
	Ip                   string `json:"ip"`
	Os                   string `json:"os"`
	PortOpen             string `json:"portOpen"`
	AssetChangesType     int    `json:"assetChangesType"`
	AssetChangesTypeName string `json:"assetChangesTypeName"`
	IsLive               int    `json:"isLive"`
	IsLiveName           string `json:"isLiveName"`
	Status               int    `json:"status"`
	StatusName           string `json:"statusName"`
	CreateTime           string `json:"createTime"`
}

type ScanTaskIpExportReq struct {
	TaskId           int    `form:"taskId" json:"taskId" binding:"required"`
	Search           string `form:"search" json:"search"`
	IsLive           int    `form:"isLive" json:"isLive"`
	Status           string `form:"status" json:"status"`
	AssetChangesType int    `form:"assetChangesType" json:"assetChangesType"`
}

type ScanTaskIpDelReq struct {
	ScanIpIds string `form:"scanIpIds" json:"scanIpIds" binding:"required"`
}

type ScanTaskIpDelResp struct{}

type ScanTaskPortListReq struct {
	ScanIpId int `form:"scanIpId" json:"scanIpId" binding:"required"`
	Page     int `form:"page" json:"page" binding:"required"`
	Size     int `form:"size" json:"size" binding:"required"`
}

type ScanTaskPortListResp struct {
	Total int64                  `json:"total"`
	List  []ScanTaskPortListList `json:"list"`
}
type ScanTaskPortListList struct {
	Id         int    `json:"id"`
	Port       int    `json:"port"`       // 端口
	Protocol   string `json:"protocol"`   // 协议
	Service    string `json:"service"`    // 服务
	Assembly   string `json:"assembly"`   // 指纹
	Remark     string `json:"remark"`     // 备注/标题
	CreateTime string `json:"createTime"` // 测试时间
	Title      string `json:"title"`      // 标题
}
