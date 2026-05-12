package typespec

// 逻辑任务创建
type LogicTaskCreateReq struct {
	Name       string `form:"name" json:"name" binding:"required"`             //任务名称
	TargetUrl  string `form:"targetUrl" json:"targetUrl" binding:"required"`   //测试目标
	Type       int    `form:"type" json:"type" binding:"required"`             //扫描类型
	ScanConfig string `form:"scanConfig" json:"scanConfig" binding:"required"` //扫描配置
}
type LogicTaskCreateResp struct {
	TaskId int `json:"taskId"`
}

// 逻辑任务结束
type LogicTaskStopReq struct {
	Id int `form:"id" json:"id"` //任务id
}
type LogicTaskStopResp struct {
}

// 逻辑任务删除
type LogicTaskDelReq struct {
	Ids string `form:"ids" json:"ids"` //任务id
}
type LogicTaskDelResp struct {
}

// 逻辑任务列表
type LogicTaskListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type LogicTaskListResp struct {
	Total int64                    `json:"total"`
	List  []LogicTaskListRespItems `json:"list"`
}

type LogicTaskListRespItems struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	TypeName   string `json:"typeName"`
	Risk       int    `json:"risk"`
	RiskName   string `json:"riskName"`
	Status     int    `json:"status"`
	StatusName string `json:"statusName"`
	UserID     int    `json:"userID"`
	UpdateTime string `json:"updateTime"`
}

// 逻辑任务-目标列表
type LogicTargetListReq struct {
	TaskId int    `form:"task_id" json:"task_id" binding:"required"`
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type LogicTargetListResp struct {
	Total int64                      `json:"total"`
	List  []LogicTargetListRespItems `json:"list"`
}

type LogicTargetListRespItems struct {
	Id          int    `json:"id"`
	Type        int    `json:"type"`
	TypeName    string `json:"typeName"`
	TargetUrl   string `json:"targetUrl"`
	Status      int    `json:"status"`
	StatusName  string `json:"statusName"`
	Risk        int    `json:"risk"`
	RiskName    string `json:"riskName"`
	IsAlive     int    `json:"isAlive"`
	IsAliveName string `json:"isAliveName"`
}

// 逻辑任务-漏洞列表
type LogicVulListReq struct {
	TaskId int    `form:"task_id" json:"task_id" binding:"required"`
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type LogicVulListResp struct {
	Total int64                   `json:"total"`
	List  []LogicVulListRespItems `json:"list"`
}

type LogicVulListRespItems struct {
	Id         int    `json:"id"`
	TargetUrl  string `json:"targetUrl"`
	Type       int    `json:"type"`     //漏洞类型
	TypeName   string `json:"typeName"` //漏洞名称
	Risk       int    `json:"risk"`
	RiskName   string `json:"riskName"`
	Location   string `json:"location"`   //漏洞位置
	CreateTime string `json:"createTime"` //发现时间
}

// 逻辑任务-日志列表
type LogicLogListReq struct {
	TaskId int    `form:"task_id" json:"task_id" binding:"required"`
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type LogicLogListResp struct {
	Total int64                   `json:"total"`
	List  []LogicLogListRespItems `json:"list"`
}

type LogicLogListRespItems struct {
	Id          int    `json:"id"`
	TargetURL   string `json:"target_url"`
	Status      int    `json:"status"`
	StatusName  string `json:"statusName"`
	CreateTime  string `json:"createTime"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAlive     int    `json:"isAlive"`
	IsAliveName string `json:"isAliveName"`
}

// 逻辑任务-漏洞详情
type LogicVulInfoReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type LogicVulInfoResp struct {
	Name               string              `json:"name"`
	Risk               int                 `json:"risk"`
	RiskName           string              `json:"riskName"`
	Description        string              `json:"description"`
	FixSuggest         string              `json:"fixSuggest"`
	Location           string              `json:"location"`
	Result             string              `json:"result"`
	PayloadSuccessFlag string              `json:"payload_success_flag"`
	Payload            string              `json:"payload"`
	VerMsg             []map[string]string `json:"verMsg"`
}

// 逻辑任务-日志详情
type LogicLogInfoReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type LogicLogInfoResp struct {
	List []LogicLogInfoRespItems `json:"list"`
}

type LogicLogInfoRespItems struct {
	Id         int    `json:"id"`
	Pocname    string `json:"pocname"`
	CreateTime string `json:"createTime"`
	TargetURL  string `json:"target_url"`
	Result     string `json:"result"`
	TaskId     int    `json:"taskId"`
}

// 逻辑任务-日志详情
type LogicEnumReq struct {
}
type LogicEnumResp struct {
	ScanType    []GlobalOptionsItemRes `json:"scanType"`    //扫描类型  1 越权 2 敏感信息遍历
	CredPattern []GlobalOptionsItemRes `json:"credPattern"` //登录凭证方式  1 cookie 2 header
	WhitePath   string                 `json:"whitePath"`   //路径白名单
	BlackPath   string                 `json:"blackPath"`   //路径黑名单
	FuzzParam   struct {
		Character string `json:"character"` //字符型
		Number    string `json:"number"`    //数字型
	} `json:"fuzzParam"` //fuzz参数
	Response struct {
		JsonKeyword   string `json:"jsonKeyword"`   //json关键字
		NoJsonSwitch  bool   `json:"noJsonSwitch"`  //非json关键字开关
		NoJsonKeyword string `json:"noJsonKeyword"` //非json关键字
	} `json:"response"`
	Crawler struct {
		MaxDeep     []GlobalOptionsItemRes `json:"maxDeep"`     // 爬取深度 0不限
		MaxUrl      []GlobalOptionsItemRes `json:"maxUrl"`      // 最大爬取url数量 0不限
		ScanRange   []GlobalOptionsItemRes `json:"scanRange"`   // 爬虫爬取范围
		Timeout     []GlobalOptionsItemRes `json:"timeout"`     // 爬虫单页面超时时间设置 0不限 单位秒
		FullTimeout []GlobalOptionsItemRes `json:"fullTimeout"` // 爬虫全局超时时间设置 0不限 单位秒
		/** 爬虫结果重复过滤设置
		0 不限，请勿设置爬虫过滤
		1 `newcrawlerx.UnLimitRepeat` 对page，method，query-name，query-value和post-data敏感
		2 `newcrawlerx.LowRepeatLevel` 对page，method，query-name和query-value敏感（默认）
		3 `newcrawlerx.MediumRepeatLevel` 对page，method和query-name敏感
		4 `newcrawlerx.HighRepeatLevel` 对page和method敏感
		5 `newcrawlerx.ExtremeRepeatLevel` 对page敏感
		*/
		ScanRepeat []GlobalOptionsItemRes `json:"scanRepeat"` // 爬虫结果重复过滤设置
		Sensitive  string                 `json:"sensitive"`  // 敏感词
		BlackList  string                 `json:"blackList"`  // 爬虫黑名单
		WhiteList  string                 `json:"whiteList"`  // 爬虫白名单
	} `json:"crawler"` // 爬虫
}

type LogicVulDeleteReq struct {
	Ids string `form:"ids" json:"ids" binding:"required"`
}
type LogicVulDeleteResp struct {
}

type LogicVulTestReq struct {
	VerMsg string `form:"verMsg" json:"verMsg" binding:"required"`
}
type LogicVulTestResp struct {
	RespVerMsg string `form:"respVerMsg" json:"respVerMsg"`
}

type LogicTaskCopyReq struct {
	Id int `form:"id" json:"id" binding:"required"` //任务id
}
type LogicTaskCopyResp struct {
	Id         int    `form:"id" json:"id" binding:"required"`                 //任务id
	Name       string `form:"name" json:"name" binding:"required"`             //任务名称
	Type       int    `form:"type" json:"type" binding:"required"`             //任务类型
	TypeName   string `form:"typeName" json:"typeName" binding:"required"`     //任务类型名称
	TargetUrl  string `form:"targetUrl" json:"targetUrl" binding:"required"`   //任务测试范围
	ScanConfig string `form:"scanConfig" json:"scanConfig" binding:"required"` //扫描配置信息
}

// LogicReportSaveReq 报告保存
type LogicReportSaveReq struct {
	Name       string `json:"name" form:"name" binding:"required"`
	Type       int    `json:"type" form:"type" binding:"required"`
	ConfigJson string `json:"configJson" form:"configJson" binding:"required"`
	Format     int    `json:"format" form:"format" binding:"required"`
	UserId     int    `json:"userId" form:"userId" binding:"required"`
	OutputType int    `json:"outputType" form:"outputType"` // 1合并输出 2逐个输出
	ObjIDName  string `json:"objIDName" form:"objIDName"`   // 逐个输出会用到 为了规避objid和目标对不上的问题
}
type LogicReportSaveResp struct {
}

// LogicFlowBaseListReq 获取流量列表接口
type LogicFlowBaseListReq struct {
	TaskId   int    `json:"taskId" form:"taskId" binding:"required"`
	TargetId int    `json:"targetId" form:"targetId"`
	Page     int    `json:"page" form:"page" binding:"required"`
	Size     int    `json:"size" form:"size" binding:"required"`
	Search   string `json:"search" form:"search"`
}
type LogicFlowBaseListResp struct {
	Total int64                        `json:"total"`
	List  []LogicFlowBaseListRespItems `json:"list"`
}
type LogicFlowBaseListRespItems struct {
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
type LogicFlowBaseInfoReq struct {
	FlowBaseId int `form:"flowBaseId" json:"flowBaseId" binding:"required"`
}
type LogicFlowBaseInfoResp struct {
	ReqHeader  string `json:"reqHeader"`
	RespHeader string `json:"respHeader"`
}

type LogicFlowBaseExportReq struct {
	TaskId int `form:"taskId" json:"taskId" binding:"required"`
}
type LogicFlowBaseExportResp struct {
	List []LogicFlowBaseListRespItems `json:"list"`
}
