package typespec

type GetAssetEnumsResp struct {
	AssetChangeType interface{} `json:"assetChangeType"`
	AssetIsLiveType interface{} `json:"assetIsLiveType"`
	DeviceWeight    interface{} `json:"deviceWeight"`
	LoginProtocol   interface{} `json:"loginProtocol"`
	TrustLevel      interface{} `json:"trustLevel"`
}

type GetChangeAssetStatResp struct {
	AddIp            int64 `json:"addIp"`
	ReduceIp         int64 `json:"reduceIp"`
	PortChangeIp     int64 `json:"portChangeIp"`
	ServiceChangeIp  int64 `json:"serviceChangeIp"`
	AssemblyChangeIp int64 `json:"assemblyChangeIp"`
}

type GetChangeAssetListReq struct {
	Search           string `form:"search" json:"search"`
	AssetChangesType int    `form:"assetChangesType" json:"assetChangesType"`
	IsLive           int    `form:"isLive" json:"isLive"`
	UpdateTime       string `form:"updateTime" json:"updateTime"`
	Page             int    `form:"page" json:"page" v:"required#页数缺失"`
	Size             int    `form:"size" json:"size" v:"required#页大小缺失"`
}

type GetChangeAssetListResp struct {
	Total int64                        `json:"total"`
	List  []GetChangeAssetListRespList `json:"list"`
}
type GetChangeAssetListRespList struct {
	Id                   int    `json:"id"`
	Ip                   string `json:"ip"`
	OperateSystem        string `json:"operateSystem"`
	PortOpen             string `json:"portOpen"`
	AssetChangesType     int    `json:"assetChangesType"`
	AssetChangesTypeName string `json:"assetChangesTypeName"`
	IsLive               int    `json:"isLive"`
	IsLiveName           string `json:"isLiveName"`
	UpdateTime           string `json:"updateTime"`
}

type ChangeAssetDelReq struct {
	AssetIds string `form:"assetIds" json:"assetIds" v:"required#资产ID缺失"`
}

type ChangeAssetDelResp struct {
}

type ChangeAssetExportReq struct {
	Search           string `form:"search" json:"search"`
	AssetChangesType int    `form:"assetChangesType" json:"assetChangesType"`
	IsLive           int    `form:"isLive" json:"isLive"`
}

// AssetCollectReq 资产收集请求结构体
type AssetCollectReq struct {
	ActivityID int `form:"activityID" json:"activityID"`
}

// AssetCollectRes 资产收集返回结构体
type AssetCollectRes struct {
	IPTotal       int64 `json:"iPTotal"`       // 总资产IP
	UnsafeIPTotal int64 `json:"unsafeIPTotal"` // 不安全IP
	ScanAsset     int64 `json:"scanAsset"`     // 资产探测
}

// AssetSummarizeRes 资产综述返回结构体
type AssetSummarizeRes struct {
	AssetTotal             int64                 `json:"assetTotal"`             // 资产总数
	SafeLoopholeTotal      int64                 `json:"safeLoopholeTotal"`      // 发现安全漏洞-数量
	ConfigErrTotal         int64                 `json:"configErrTotal"`         // 发现配置问题-数量
	AssetRiskStat          []AssetRiskStatistics `json:"assetRiskStat"`          // 资产风险统计
	RiskStatics            []AssetRiskStatistics `json:"riskStatics"`            // 风险类型统计
	AssetRiskLevelTrendRes []AssetRiskTrend      `json:"assetRiskLevelTrendRes"` // 资产风险分布趋势
	AssetRiskTypeTrendRes  []AssetRiskTypeTrend  `json:"assetRiskTypeTrendRes"`  // 资产类型分布趋势
	RecentDangerAsset      []RiskAsset           `json:"recentDangerAsset"`      // 近期危险资产
	// AssetRiskTrendRes      []AssetTypeTrend        `json:"assetRiskTrendRes"`      // 资产风险趋势
	// 老的
	YoYLastWeekAssetTotal     int64               `json:"yoyLastWeekAssetTotal"`     // 同比上周新增资产总数
	NewAddIpNum               int64               `json:"newAddIpNum"`               // 新增IP
	TkNewAddIpNum             int64               `json:"tkNewAddIpNum"`             // 本周新增IP
	NewReduceIpNum            int64               `json:"newReduceIpNum"`            // 减少IP
	TkNewReduceIpNum          int64               `json:"tkNewReduceIpNum"`          // 本周减少IP
	AssetTypeStatics          []AssetRiskTrend    `json:"assetTypeStatics"`          // 资产类型统计
	AssetTrendChangeRes       AssetTrendChangeRes `json:"assetTrendChangeRes"`       // 资产变化趋势
	AssetReduceTrendChangeRes AssetTrendChangeRes `json:"assetReduceTrendChangeRes"` // 资产变化趋势
}

// AssetRiskStatistics 资产风险统计
type AssetRiskStatistics struct {
	AssetType string `json:"assetType"`
	Count     int    `json:"count"`
}

// AssetRiskTrend 资产风险分布趋势
type AssetRiskTrend struct {
	Date        string `json:"date"`
	Fatal       int    `json:"fatal"`
	High        int    `json:"high"`
	Medium      int    `json:"medium"`
	Low         int    `json:"low"`
	Safe        int    `json:"safe"`
	FatalAsset  int    `json:"fatalAsset"`
	HighAsset   int    `json:"highAsset"`
	MediumAsset int    `json:"mediumAsset"`
	LowAsset    int    `json:"lowAsset"`
	SafeAsset   int    `json:"safeAsset"`
	TotalAsset  int    `json:"totalAsset"`
	TotalVul    int    `json:"totalVul"`
}

type AssetRiskTypeTrend struct {
	Date   string `json:"date"`
	UnSafe int    `json:"unSafe"`
	Safe   int    `json:"safe"`
}

type RiskAsset struct {
	ID            int        `json:"id"`
	IP            string     `json:"ip"`
	AssetTypeName string     `json:"assetTypeName"` // 业务类型
	Os            string     `json:"os"`            // 操作系统
	OpenPort      string     `json:"openPort"`      // 开放端口
	VulStatics    VulStatics `json:"vulStatics"`    // 资产漏洞统计
	RiskLevel     string     `json:"riskLevel"`
	Time          string     `json:"time"`
}

type AssetTypeTrend struct {
	Date      string `json:"date"`
	AssetType string `json:"assetType"`
	Count     int    `json:"count"`
}

// RecentDangerAssetInfo 近期危险资产
type RecentDangerAssetInfo struct {
	IP           string                 `json:"ip"`           // ip或域名
	OpSys        string                 `json:"opSys"`        // 操作系统
	AssetType    int                    `json:"assetType"`    // 资产类型
	PortList     []ScanTaskPortListList `json:"portList"`     // 端口信息
	RiskLevel    int                    `json:"riskLevel"`    // 风险等级
	RiskLevelStr string                 `json:"riskLevelStr"` // 风险等级-string
	VulInfo      VulStatics             `json:"vulInfo"`      // 漏洞 高中低
}

// AssetTypeStatistics 资产类型统计
type AssetTypeStatistics struct {
	AssetType string
	Count     int
}

// AssetSummarizeReq 资产变化趋势请求结构体
type AssetSummarizeReq struct {
	StartTime string `form:"startTime"` // 开始时间
	EndTime   string `form:"endTime"`   // 结束时间
	TimeType  int    `form:"timeType"`  // 类型选择 1 最近7天 2 本月 3 本年
}

// AssetTrendChangeRes 资产变化趋势
type AssetTrendChangeRes struct {
	Name  []string  `json:"name"`
	Value []float64 `json:"value"`
}

// CycleAssetRiskStatistics 周期资产风险统计
type CycleAssetRiskStatistics struct {
	Time                     string `json:"time"`
	TotalAsset               int    `json:"totalAsset"`
	AssetRiskTrendStatistics []AssetRiskTrendStatistics
}

// AssetRiskTrendStatistics 资产风险统计
type AssetRiskTrendStatistics struct {
	AssetType  int `json:"assetType"`
	VulCount   int `json:"vulCount"`
	AssetCount int `json:"assetCount"`
}
