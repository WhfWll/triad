package typespec

// 三方工具

// xray 创建
type TripartiteToolsXRayCreateReq struct {
	TaskName  string `form:"taskName" json:"taskName" binding:"required,max=50"`
	Target    string `form:"target" json:"target" binding:"required,max=50"`
	IsCrawler bool   `json:"isCrawler"`
}
type TripartiteToolsXRayCreateRes struct {
}

// xray 上传
type TripartiteToolsXRayUploadReq struct {
	TaskName string `form:"taskName" json:"taskName" binding:"required,max=50"`
}
type TripartiteToolsXRayUploadRes struct {
}

// xray 删除任务
type TripartiteToolsXRayDelReq struct {
	XrayIds string `form:"xrayIds" json:"xrayIds" binding:"required"`
}
type TripartiteToolsXRayDelRes struct {
}

// xray 获取列表
type TripartiteToolsXRayPageReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type TripartiteToolsXRayPageRes struct {
	Page  int                           `json:"page"`
	Size  int                           `json:"size"`
	Total int64                         `json:"total"`
	List  []TripartiteToolsXRayPageItem `json:"list"`
}
type TripartiteToolsXRayPageItem struct {
	Id         int    `json:"id"`
	TaskName   string `json:"taskName"`
	RiskNum    int    `json:"riskNum"`
	Status     int    `json:"status"`
	StatusEnum string `json:"statusEnum"`
	CreateTime string `json:"createTime"`
}

// xray 获取详情列表
type TripartiteToolsXRayDetailPageReq struct {
	XrayId int    `form:"xrayId" json:"xrayId" binding:"required"`
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type TripartiteToolsXRayDetailPageRes struct {
	Page  int                                 `json:"page"`
	Size  int                                 `json:"size"`
	Total int64                               `json:"total"`
	List  []TripartiteToolsXRayDetailPageItem `json:"list"`
}
type TripartiteToolsXRayDetailPageItem struct {
	Id                 int        `json:"id"`
	Addr               string     `json:"addr"`
	PluginVul          string     `json:"pluginVul"`
	ParamPosition      string     `json:"paramPosition"`
	ParamKey           string     `json:"paramKey"`
	Payload            string     `json:"payload"`
	RequestAndResponse [][]string `json:"requestAndResponse"`
	Extra              string     `json:"extra"`
	CreateTime         string     `json:"createTime"`
}

// burpsuite 创建任务
type TripartiteToolsBurpsuiteCreateReq struct {
	TaskName string `form:"taskName" json:"taskName" binding:"required,max=50"`
	Target   string `form:"target" json:"target" binding:"required,max=5120"`
}
type TripartiteToolsBurpsuiteCreateRes struct {
}

// burpsuite 上传
type TripartiteToolsBurpsuiteUploadReq struct {
	TaskName string `form:"taskName" json:"taskName" binding:"required,max=50"`
}
type TripartiteToolsBurpsuiteUploadRes struct {
}

// burpsuite 删除任务
type TripartiteToolsBurpsuiteDelReq struct {
	BurpsuiteIds string `form:"burpsuiteIds" json:"burpsuiteIds" binding:"required"`
}
type TripartiteToolsBurpsuiteDelRes struct {
}

// burpsuite 任务列表
type TripartiteToolsBurpsuitePageReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type TripartiteToolsBurpsuitePageRes struct {
	Page  int                                `json:"page"`
	Size  int                                `json:"size"`
	Total int64                              `json:"total"`
	List  []TripartiteToolsBurpsuitePageItem `json:"list"`
}
type TripartiteToolsBurpsuitePageItem struct {
	Id         int    `json:"id"`
	TaskName   string `json:"taskName"`
	RiskEnum   string `json:"riskEnum"`
	Status     int    `json:"status"`
	StatusEnum string `json:"statusEnum"`
	CreateTime string `json:"createTime"`
}

// burpsuite 任务详情列表
type TripartiteToolsBurpsuiteDetailPageReq struct {
	BurpsuiteId int    `form:"burpsuiteId" json:"burpsuiteId" binding:"required"`
	Page        int    `form:"page" json:"page" binding:"required"`
	Size        int    `form:"size" json:"size" binding:"required"`
	Search      string `form:"search" json:"search"`
}
type TripartiteToolsBurpsuiteDetailPageRes struct {
	Page  int                                      `json:"page"`
	Size  int                                      `json:"size"`
	Total int64                                    `json:"total"`
	List  []TripartiteToolsBurpsuiteDetailPageItem `json:"list"`
}
type TripartiteToolsBurpsuiteDetailPageItem struct {
	Id              int        `json:"id"`
	Action          string     `json:"action"`
	IssueType       string     `json:"issueType"`
	Host            string     `json:"host"`
	Path            string     `json:"path"`
	InsertionPoint  string     `json:"insertionPoint"`
	Confidence      string     `json:"confidence"`
	Time            string     `json:"time"`
	RiskEnum        string     `json:"riskEnum"`
	Desc            string     `json:"desc"`
	IssueBackground string     `json:"issueBackground"`
	Fix             string     `json:"fix"`
	RequestResponse [][]string `json:"requestResponse"`
}

// wifi 所有在线wifi列表
type TripartiteToolsWifiApListReq struct {
	Search    string `form:"search" json:"search"`
	StartDate string `form:"startDate" json:"startDate"`
	EndDate   string `form:"endDate" json:"endDate"`
}
type TripartiteToolsWifiApListRes struct {
	List []TripartiteToolsWifiApListItem `json:"list"`
}
type TripartiteToolsWifiApListItem struct {
	Ssid             string `json:"ssid"`
	LastSignalRssi   int    `json:"lastSignalRssi"`
	SsidCryptset     int    `json:"ssidCryptset"`
	SsidCryptsetEnum string `json:"ssidCryptsetEnum"`
	SourceMac        string `json:"sourceMac"`
	LastTime         string `json:"lastTime"`
	Manuf            string `json:"manuf"`
	Carrier          string `json:"carrier"`
}

// wifi 创建任务
type TripartiteToolsWifiCreateReq struct {
	TaskName  string `form:"taskName" json:"taskName" binding:"required,max=50"`
	SourceMac string `form:"sourceMac" json:"sourceMac" binding:"required"`
}
type TripartiteToolsWifiCreateRes struct {
}

// wifi 任务列表
type TripartiteToolsWifiPageReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"` // wifi名称 ssid
}
type TripartiteToolsWifiPageRes struct {
	Page  int                           `json:"page"`
	Size  int                           `json:"size"`
	Total int64                         `json:"total"`
	List  []TripartiteToolsWifiPageItem `json:"list"`
}
type TripartiteToolsWifiPageItem struct {
	TaskId     int                 `json:"taskId"`
	TaskName   string              `json:"taskName"`
	Ssid       string              `json:"ssid"`
	Encrypt    string              `json:"encrypt"`
	Carrier    string              `json:"carrier"`
	Channel    string              `json:"channel"`
	StartTime  string              `json:"startTime"`
	CreateTime string              `json:"createTime"`
	Status     string              `json:"status"`
	Passwd     string              `json:"passwd"`
	LogList    []map[string]string `json:"logList"`
}

// wifi 任务详情
type TripartiteToolsWifiDelReq struct {
	TaskIds string `form:"taskIds" json:"taskIds" binding:"required"`
}
type TripartiteToolsWifiDelRes struct {
}
