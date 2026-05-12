package typespec

type AssetTreeOverallReq struct {
	Search                string `form:"search" json:"search"`
	Ip                    string `form:"ip" json:"ip"` // ip或域名
	OperateSystem         string `form:"operateSystem" json:"operateSystem"`
	RiskLevel             int    `form:"riskLevel" json:"riskLevel"`
	Tags                  string `form:"tags" json:"tags"`
	AssetName             string `form:"assetName" json:"assetName"`
	BusinessSystem        string `form:"businessSystem" json:"businessSystem"`
	ResponsibleDepartment string `form:"responsibleDepartment" json:"responsibleDepartment"`
	FilingLevel           int    `form:"filingLevel" json:"filingLevel"`
	Port                  string `form:"port" json:"port"`
	Service               string `form:"service" json:"service"`
	Component             string `form:"component" json:"component"`
	VulName               string `form:"vulName" json:"vulName"`
}

type AssetTreeOverallResp struct {
	List []AssetTreeOverallRespItems `json:"list"`
}

type AssetTreeOverallRespItems struct {
	Id    int                         `json:"id"`
	Pid   int                         `json:"pid"`
	Name  string                      `json:"name"`
	Type  int                         `json:"type"`
	Items []AssetTreeOverallRespItems `json:"items"`
}
type AssetGroupResp struct {
	List []AssetGroupRespItems `json:"list"`
}
type AssetGroupRespItems struct {
	Id    int                   `json:"id"`
	Pid   int                   `json:"pid"`
	Name  string                `json:"name"`
	Level int                   `json:"level"`
	Items []AssetGroupRespItems `json:"items"`
}

// AssetTreeAddReq 新增资产组
type AssetTreeAddReq struct {
	Name   string `form:"name" json:"name" binding:"required"`
	Pid    int    `form:"pid" json:"pid"`
	Remark string `form:"remark" json:"remark"`
}

// AssetGroupEditReq 编辑资产组
type AssetGroupEditReq struct {
	Id     int    `form:"id" json:"id" binding:"required"` // 资产组id
	Name   string `form:"name" json:"name"`                // 资产组名称
	Pid    int    `form:"pid" json:"pid"`                  // 上级资产组
	Remark string `form:"remark" json:"remark"`            // 资产组说明
}

// AssetDeleteReq 资产组和资产删除
type AssetDeleteReq struct {
	GroupIds string `form:"groupIds" json:"groupIds"`
	AssetIds string `form:"assetIds" json:"assetIds"`
}

// AssetListReq 子资请求参数
type AssetListReq struct {
	GroupId      int    `form:"groupID" json:"groupID"`              // 资产组
	Search       string `form:"search"`                              // 搜索
	AssetIP      string `form:"assetIP"`                             // 资产IP
	AssetRisk    *int   `form:"assetRisk"`                           // 资产风险
	Port         string `form:"port"`                                // 端口
	Service      string `form:"service"`                             // 服务
	SystemOp     string `form:"systemOp"`                            // 系统
	VulName      string `form:"vulName"`                             // 漏洞名称
	Domain       string `form:"domain"`                              // 域名
	Finger       string `form:"finger"`                              // 指纹
	AssetType    string `form:"assetType"`                           // 资产类型
	FillingLevel int    `form:"fillingLevel"`                        // 等保级别
	Tags         string `form:"tags"`                                // 标签
	IsCloudHost  string `form:"isCloudHost"`                         // 是否是云主机
	Page         int    `form:"page" json:"page" binding:"required"` //任务组页数
	Size         int    `form:"size" json:"size" binding:"required"` //任务组每页大小
}

// AssetListRes 资产返回参数
type AssetListRes struct {
	Count      int          `json:"count"`
	AssetsInfo []AssetsInfo `json:"assetsInfo"`
}

// AssetsInfo 资产信息
type AssetsInfo struct {
	AssetID               int        `json:"assetID"`
	IP                    string     `json:"ip"`                    // 资产IP
	AssetType             string     `json:"assetType"`             // 资产类型
	System                string     `json:"system"`                // 系统
	OpenPort              string     `json:"openPort"`              // 开放端口
	AssetRiskName         string     `json:"assetRiskName"`         // 资产风险名 严重 高 中 低
	AssetGroupName        string     `json:"assetGroupName"`        // 资产组名
	VulStatics            VulStatics `json:"vulStatics"`            // 资产漏洞统计
	TestTime              string     `json:"testTime"`              // 测试/更新时间
	DeviceWeight          string     `json:"deviceWeight"`          // 设备权重 1高 2中 3低 4极高 5极低
	TrustLevel            string     `json:"trustLevel"`            // 可信设备 1 可信 2 未登记
	Location              string     `json:"location"`              // 地理
	ResponsibleDepartment string     `json:"responsibleDepartment"` // 责任部门

	//IsCloudHost    int        `json:"isCloudHost"`    // 云主机
	//AssetStatics   AssetStatics `json:"assetStatics"`   // 资产风险统计
}

// AssetListDelReq 资产列表删除
type AssetListDelReq struct {
	IDs string `form:"ids"`
}

// VulStatics 漏洞统计
type VulStatics struct {
	DeadlyVul     int `json:"deadlyVul"`     // 致命漏洞
	HighRiskVul   int `json:"highRiskVul"`   // 高危漏洞
	MediumRiskVul int `json:"mediumRiskVul"` // 中危漏洞
	LowRiskVul    int `json:"lowRiskVul"`    // 低危漏洞
}

// AssetStatics 资产统计
type AssetStatics struct {
	HighRiskAsset   int `json:"highRiskAsset"`
	MiddleRiskAsset int `json:"middleRiskAsset"`
	LowRiskAsset    int `json:"lowRiskAsset"`
	SafeAsset       int `json:"safeAsset"`
	TotalAsset      int `json:"totalAsset"`
}

// AssetAddReq 资产添加
type AssetAddReq struct {
	IP                    string `form:"ip" json:"ip" binding:"required"`                     // ip或域名Y
	AssetGroupID          int    `form:"assetGroupID" json:"assetGroupID" binding:"required"` // 所属资产组IDY
	Name                  string `form:"name" json:"name"`                                    // 资产名Y
	OpSys                 string `form:"opSys" json:"opSys"`                                  // 操作系统Y
	SystemAdmin           string `form:"systemAdmin" json:"systemAdmin"`                      // 系统管理员Y
	Tags                  string `form:"tags" json:"tags"`                                    // 资产标签Y
	ResponsibleDepartment string `form:"responsibleDepartment" json:"responsibleDepartment"`  // 责任部门
	// 不展示的
	DeviceWeight         int    `form:"deviceWeight" json:"deviceWeight"`                 // 设备权重 1高 2中 3低 4极高 5极低
	TrustLevel           int    `form:"trustLevel" json:"trustLevel"`                     // 可信设备 1 可信 2 未登记
	IsCloudHost          int    `form:"isCloudHost" json:"isCloudHost"`                   // 是否是云主机
	EquipmentForm        string `form:"equipmentForm" json:"equipmentForm"`               // 设备形态
	EqualProtectionLevel int    `form:"equalProtectionLevel" json:"equalProtectionLevel"` // 等保等级 0无、1一级、2二级、3三级、4四级、5五级
	// 资产登录信息
	User     string `form:"user" json:"user"` // 登录用户名
	Port     int    `form:"port" json:"port"` // 登录端口
	Protocol int    `form:"protocol"`         // 连接协议，如 ssh、rdp、telnet
	Password string `form:"password"`         // 登录密码
	// 以前的资产默认字段
	IPSegment           string `form:"ipSegment" json:"ipSegment"`                     // IP段
	AssetType           int    `form:"assetType" json:"assetType"`                     // 资产类型
	BaseSoftwareName    string `form:"baseSoftwareName" json:"baseSoftwareName"`       // 基础软件名称
	BaseSoftwareVersion string `form:"baseSoftwareVersion" json:"baseSoftwareVersion"` // 基础软件版本
	BaseHardwareName    string `form:"baseHardwareName" json:"baseHardwareName"`       // 硬件
	Purpose             string `form:"purpose" json:"purpose"`                         // 用途
	SystemName          string `form:"systemName" json:"systemName"`                   // 系统名称
	SystemOp            string `form:"systemOp" json:"systemOp"`                       // 系统运维人员
	DeploymentLocation  string `form:"deploymentLocation" json:"deploymentLocation"`
}

// AssetEditReq 资产修改
type AssetEditReq struct {
	ID int `form:"id" json:"id" binding:"required"` // 资产组id
	AssetAddReq
}

// AssetDetailReq 资产详情请求参数
type AssetDetailReq struct {
	ID         int `form:"id"`
	Page       int `json:"page" form:"page"`
	Size       int `json:"size" form:"size"`
	SelectType int `json:"selectType" form:"selectType"  binding:"required"` // 检索类型 1资产信息 2 漏洞信息 3 配置安全 4 管理信息
}

// AssetDetail 资产详情
type AssetDetail struct {
	Info       []ScanTaskPortListList `json:"assetInfo"`  // 资产信息
	VulInfo    VulRiskListResp        `json:"vulInfo"`    // 漏洞信息
	ConfigInfo []ConfigVulListItem    `json:"configInfo"` // 配置安全
	ManageInfo AssetInfo              `json:"manageInfo"` // 管理信息
}

// AssetInfo 资产信息
type AssetInfo struct {
	IP                      string                 `json:"ip"`                             // ip或域名
	AssetName               string                 `json:"assetName"`                      // 资产名称
	AssetGroupID            int                    `json:"assetGroupID"`                   // 所属资产组ID
	AssetGroupName          string                 `json:"assetGroupName"`                 // 所属资产组ID
	OpSys                   string                 `json:"opSys"`                          // 操作系统
	IPSegment               string                 `json:"ipSegment"`                      // IP段
	AssetType               int                    `json:"assetType"`                      // 资产类型
	BaseSoftwareName        string                 `json:"baseSoftwareName"`               // 基础软件名称
	BaseSoftwareVersion     string                 `json:"baseSoftwareVersion"`            // 基础软件版本
	BaseHardwareName        string                 `json:"baseHardwareName"`               // 硬件
	Purpose                 string                 `json:"purpose"`                        // 用途
	EquipmentForm           string                 `json:"equipmentForm"`                  // 设备形态
	DeploymentLocation      string                 `json:"deploymentLocation"`             // 部署位置
	EqualProtectionLevel    int                    `json:"equalProtectionLevel"`           // 等保等级 0无、1一级、2二级、3三级、4四级、5五级
	EqualProtectionLevelStr string                 `json:"equalProtectionLevelStr"`        // 等保等级 0无、1一级、2二级、3三级、4四级、5五级
	IsCloudHost             int                    `form:"isCloudHost" json:"isCloudHost"` // 是否是云主机
	DeviceWeight            string                 `json:"deviceWeight"`                   // 设备权重 1高 2中 3低 4极高 5极低
	TrustLevel              string                 `json:"trustLevel"`                     // 可信设备 1 可信 2 未登记
	Location                string                 `json:"location"`                       // 地理未知
	SystemName              string                 `json:"systemName"`                     // 系统名称
	SystemAdmin             string                 `json:"systemAdmin"`                    // 系统管理员
	SystemOp                string                 `json:"systemOp"`                       // 系统运维人员
	Tags                    string                 `json:"tags"`                           // 资产标签
	PortList                []ScanTaskPortListList `json:"portList"`                       // 资产信息
	RiskLevelStr            string                 `json:"riskLevelStr"`                   // 风险等级
	Score                   string                 `json:"score"`                          // 评分
	Connections             []AssetConnectionDTO   `json:"connections"`                    // 资产连接方式列表（简化）
	ResponsibleDepartment   string                 `json:"responsibleDepartment"`          // 责任部门
}

// AssetConnectionDTO 资产连接方式
type AssetConnectionDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// GetAssetGroupEnumsResp 资产组枚举
type GetAssetGroupEnumsResp struct {
	AssetType     interface{} `json:"assetType"`    // 资产类型
	FillingLevel  interface{} `json:"fillingLevel"` // 等保等级
	AssetRisk     interface{} `json:"assetRisk"`    // 资产风险
	DeviceWeight  interface{} `json:"deviceWeight"`
	TrustLevel    interface{} `json:"trustLevel"`
	AssetGroup    interface{} `json:"assetGroup"`
	LoginProtocol interface{} `json:"loginProtocol"`
}

// AssetImportReq 资产详情请求参数
type AssetImportReq struct {
	File interface{} `form:"file"`
}

// AssetExportListReq 子资请求参数
type AssetExportListReq struct {
	IDs          string `form:"ids" json:"ids"`
	GroupId      int    `form:"groupID" json:"groupID"`
	Search       string `form:"search"`
	AssetIP      string `form:"assetIP"`      // 资产IP
	AssetRisk    int    `form:"assetRisk"`    // 资产风险
	Port         string `form:"port"`         // 端口
	Service      string `form:"service"`      // 服务
	SystemOp     string `form:"systemOp"`     // 系统
	VulName      string `form:"vulName"`      // 漏洞名称
	Domain       string `form:"domain"`       // 域名
	Finger       string `form:"finger"`       // 指纹
	AssetType    int    `form:"assetType"`    // 资产类型
	FillingLevel int    `form:"fillingLevel"` // 等保级别
	Tags         string `form:"tags"`         // 标签
	ExportType   int    `form:"exportType"`   // 导出类型 1 常规导出 2 特殊导出
}

// AssetExportListRes 资产返回参数
type AssetExportListRes struct {
	Count            int                `json:"count"`
	AssetsExportInfo []AssetsExportInfo `json:"assetsInfo"`
}

// AssetsExportInfo 资产信息
type AssetsExportInfo struct {
	IP                    string `json:"资产IP"`
	AssetGroupName        string `json:"资产组名"`
	AssetName             string `json:"资产名称"` // 资产名称
	OpSys                 string `json:"系统"`   // 操作系统
	AssetType             string `json:"资产类型"` // 资产类型
	EquipmentForm         string `json:"设备形态"` // 设备形态
	EqualProtectionName   string `json:"等保等级"` // 等保等级 0无、1一级、2二级、3三级、4四级、5五级
	SystemAdmin           string `json:"管理员"`  // 系统管理员
	DeviceWeight          string `json:"设备权重"` // 设备权重 1高 2中 3低 4极高 5极低
	TrustLevel            string `json:"可信设备"` // 可信设备 1 可信 2 未登记
	Tags                  string `json:"标签"`   // 资产标签
	OpenPort              string `json:"开放端口"`
	AssetRiskName         string `json:"资产风险等级"`
	ResponsibleDepartment string `json:"部门"`
	IsCloudHost           bool   `json:"云主机"`
	TestTime              string `json:"测试时间"`
	Service               string `json:"服务"`
	Location              string `json:"地理位置"` // 地理位置
}

// AssetGroupDetailReq 资产组详情
type AssetGroupDetailReq struct {
	Id int `form:"id" json:"id" binding:"required"` // 资产组id
}

// AssetGroupDetail 资产组详情
type AssetGroupDetail struct {
	ID      int    `form:"id" json:"id"` // 资产组id
	Name    string `form:"name" json:"name"`
	Pid     int    `form:"pid" json:"pid"`
	PidName string `form:"pidName" json:"pidName"`
	Remark  string `form:"remark" json:"remark"`
}

// SelectAllAssetRes 所有资产返回参数
type SelectAllAssetRes struct {
	AssetIDs []int    `json:"assetIDs"`
	AssetIPs []string `json:"assetIPs"`
}

// SelectAllAssetReq 子资请求参数
type SelectAllAssetReq struct {
	GroupId      int    `form:"groupID" json:"groupID"` // 资产组
	Search       string `form:"search"`                 // 搜索
	AssetIP      string `form:"assetIP"`                // 资产IP
	AssetRisk    *int   `form:"assetRisk"`              // 资产风险
	Port         string `form:"port"`                   // 端口
	Service      string `form:"service"`                // 服务
	SystemOp     string `form:"systemOp"`               // 系统
	VulName      string `form:"vulName"`                // 漏洞名称
	Domain       string `form:"domain"`                 // 域名
	Finger       string `form:"finger"`                 // 指纹
	AssetType    string `form:"assetType"`              // 资产类型
	FillingLevel int    `form:"fillingLevel"`           // 等保级别
	Tags         string `form:"tags"`                   // 标签
	IsCloudHost  string `form:"isCloudHost"`            // 是否是云主机
}

// AssetConnListReq 资产连接请求参数
type AssetConnListReq struct {
	IP       string `form:"ip"`
	Port     int    `form:"port" json:"port"`         // 端口
	Protocol int    `form:"protocol" json:"protocol"` // 协议类型
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
}

// AssetConnListRes 资产返回参数
type AssetConnListRes struct {
	Count         int64           `json:"count"`
	AssetConnInfo []AssetConnInfo `json:"assetConnlist"`
}

// AssetConnInfo 资产连接信息
type AssetConnInfo struct {
	ID          int    `json:"id"`
	IP          string `json:"ip"`          // ip或域名
	Port        int    `json:"port"`        // 端口
	Protocol    int    `json:"protocol"`    // 协议ID
	ProtocolStr string `json:"protocolStr"` // 协议名称
	User        string `json:"user"`        // 登录用户名
	Pass        string `json:"password"`    // 登录密码
}
