package typespec

import "time"

// AuthInfoReq 授权信息请求结构体
type AuthInfoReq struct {
}

// AuthInfoRes 授权信息返回
type AuthInfoRes struct {
	ProductName string `json:"productName"`
	ProductID   string `json:"productID"`   //
	AuthCode    string `json:"authCode"`    // 产品序列号
	AuthTime    string `json:"authTime"`    // 授权日期
	AuthExpTime string `json:"authExpTime"` // 授权过期时间
	AuthDays    string `json:"authDays"`    // 授权时长
	LeftDays    string `json:"leftDays"`    // 剩余时间
	Status      bool   `json:"status"`
}

// AuthSaveReq 授权操作请求结构体
type AuthSaveReq struct {
	AuthCode string `form:"authCode" json:"authCode"`
}

type TargetIpSaveReq struct {
	IsOpen    int    `form:"isOpen" json:"isOpen" binding:"required"`
	Type      int    `form:"type" json:"type" binding:"required"`
	WhiteList string `form:"whiteList" json:"whiteList"`
	BlackList string `form:"blackList" json:"blackList"`
}

type TargetIpSaveResp struct{}

type TargetIpListResp struct {
	IsOpen    int    `json:"isOpen"`
	Type      int    `json:"type"`
	WhiteList string `json:"whiteList"`
	BlackList string `json:"blackList"`
}

type GetReverseIpHostResp struct {
	ReverseType int    `json:"reverseType"`
	ReverseHost string `json:"reverseHost"`
	ReversePort int    `json:"reversePort"`
}

type ReverseIpHostSaveReq struct {
	ReverseType int    `form:"reverseType" json:"reverseType" binding:"required"`
	ReverseHost string `form:"reverseHost" json:"reverseHost"`
	ReversePort int    `form:"reversePort" json:"reversePort" binding:"required"`
}

type ReverseIpHostSaveResp struct{}

type SystemConfigBackupListReq struct {
	Page int `json:"page" form:"page" binding:"required"` //页码
	Size int `json:"size" form:"size" binding:"required"` //每页的数量
}

type SystemConfigBackupListRes struct {
	List  []SystemConfigBackupListItemRes `json:"list"`
	Total int64                           `json:"total"`
}

type SystemConfigBackupListItemRes struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
	Path       string `json:"path"`
}

type LSystemConfigBackupListItemRes struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
	Path       string `json:"path"`
}

type SystemConfigBackupDeleteReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type SystemConfigBackupDeleteRes struct{}

type SystemConfigBackupDownloadReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type SystemConfigBackupDownloadRes struct {
	Path string `json:"path"`
}

type SystemConfigBackupConfigReq struct {
	IsOpen   int       `json:"isOpen" form:"isOpen" binding:"required"` //是否开启周期备份，1-开启，2-关闭
	Cycle    int       `json:"cycle" form:"cycle" binding:"required"`   //备份周期，按月计算
	SaveTime time.Time `json:"saveTime" form:"saveTime"`                //创建时间
	RunTime  time.Time `json:"runTime" form:"runTime"`                  //执行时间
}

type SystemConfigBackupConfigRes struct{}

type SystemConfigBackupConfigInfoRes struct {
	IsOpen int `json:"isOpen"`
	Cycle  int `json:"cycle"`
}

type SystemConfigBackupNowRes struct{}

type SystemConfigBackupRestoreReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type SystemConfigBackupRestoreRes struct{}

type SystemSettingIpWhiteSaveReq struct {
	IsOpen int    `json:"isOpen" form:"isOpen" binding:"required" doc:"是否开启，1 - 开启 2 - 关闭"`
	Ip     string `json:"ip" form:"ip" binding:"required" doc:"白名单ip或者ip段，多个时用换行符或英文逗号分割"`
}

type SystemSettingIpWhiteSaveRes struct{}

type SystemSettingIpWhiteInfoRes struct {
	IsOpen int    `json:"isOpen"`
	Ip     string `json:"ip"`
}

type SystemSettingSyslogSaveReq struct {
	IsOpen int    `json:"isOpen" form:"isOpen" binding:"required" doc:"是否开启，1-开启 2-关闭"`
	Ip     string `json:"ip" form:"ip" binding:"required" doc:"syslog服务器地址"`
	Port   int    `json:"port" form:"port" binding:"required" doc:"syslog服务器端口"`
	Types  string `json:"types" form:"types" binding:"required" doc:"发送日志类型 '1'-审计日志 '2'-调试日志 '3'-告警日志 '4'-报错日志"`
}

type SystemSettingSyslogSaveRes struct{}

type SystemSettingSyslogInfoRes struct {
	IsOpen int    `json:"isOpen"`
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Types  string `json:"types"`
}

type SystemSettingMailSaveReq struct {
	Address  string `json:"address" form:"address" binding:"required" doc:"邮箱服务器地址"`
	Port     int    `json:"port" form:"port" binding:"required" doc:"邮箱服务器端口"`
	Username string `json:"username" form:"username" binding:"required" doc:"邮箱账号"`
	Password string `json:"password" form:"password" binding:"required" doc:"邮箱密码"`
	Encrypt  string `json:"encrypt" form:"encrypt" doc:"加密方式"`
}

type SystemSettingMailSaveRes struct{}

type SystemSettingMailInfoRes struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SystemSettingNetworkConfigSaveReq struct {
	Ip               string `json:"ip" form:"ip" binding:"required" doc:"IP地址"`
	Mask             string `json:"mask" form:"mask" binding:"required" doc:"子网掩码"`
	Gateway          string `json:"gateway" form:"gateway" binding:"required" doc:"默认网关"`
	DnsServer        string `json:"dnsServer" form:"dnsServer" binding:"required" doc:"DNS服务器"`
	StandbyDnsServer string `json:"standbyDnsServer" form:"standbyDnsServer" doc:"备用DNS服务器"`
	WebPort          int    `json:"webPort" form:"webPort" binding:"required" doc:"业务端口"`
}

type SystemSettingNetworkConfigSaveRes struct{}

type SystemSettingNetworkConfigInfoRes struct {
	Ip               string `json:"ip"`
	Mask             string `json:"mask"`
	Gateway          string `json:"gateway"`
	DnsServer        string `json:"dnsServer"`
	StandbyDnsServer string `json:"standbyDnsServer"`
	WebPort          int    `json:"webPort"`
}

type TcpBlindTestSaveReq struct {
	Type int    `form:"type" json:"type" binding:"required" doc:"盲测平台类型 1-使用系统盲测平台 2-使用自定义盲测平台"`
	Host string `form:"host" json:"host" doc:"监听ip"`
	Port int    `form:"port" json:"port" binding:"required" doc:"监听端口"`
}
type TcpBlindTestSaveRes struct{}

type TcpBlindTestInfoRes struct {
	Type int    `json:"type"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HttpBlindTestSaveReq struct {
	Type int    `form:"type" json:"type" binding:"required" doc:"盲测平台类型 1-使用系统盲测平台 2-使用自定义盲测平台"`
	Host string `form:"host" json:"host" doc:"监听ip"`
	Port int    `form:"port" json:"port" binding:"required" doc:"监听端口"`
}

type HttpBlindTestSaveRes struct{}

type HttpBlindTestInfoRes struct {
	Type int    `json:"type"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DnsBlindTestSaveReq struct {
	Type   int    `form:"type" json:"type" binding:"required" doc:"盲测平台类型 1-使用系统盲测平台 2-使用自定义盲测平台"`
	Domain string `form:"domain" json:"domain" doc:"解释域名"`
}

type DnsBlindTestSaveRes struct{}

type DnsBlindTestInfoRes struct {
	Type   int    `json:"type"`
	Domain string `json:"domain"`
}

type IcmpBlindTestSaveReq struct {
	Type int    `form:"type" json:"type" binding:"required" doc:"盲测平台类型 1-使用系统盲测平台 2-使用自定义盲测平台"`
	Host string `form:"host" json:"host" doc:"监听ip"`
}

type IcmpBlindTestSaveRes struct{}

type IcmpBlindTestInfoRes struct {
	Type int    `json:"type"`
	Host string `json:"host"`
}

type CurTasksInfoRes struct {
	CurIp    int `json:"curIp"`
	CurTasks int `json:"curTasks"`
}

type CurTasksSaveReq struct {
	CurIp    int `form:"curIp" json:"curIp" binding:"required"`
	CurTasks int `form:"curTasks" json:"curTasks" binding:"required"`
}

type CurTasksSaveResp struct {
}

type CpuInfoRes struct {
	List interface{} `json:"list"`
}

type MemoryInfoRes struct {
	List interface{} `json:"list"`
}

type DiskInfoRes struct {
	Free        int `json:"free"`
	Used        int `json:"used"`
	Total       int `json:"total"`
	FreePercent int `json:"freePercent"`
	UsedPercent int `json:"usedPercent"`
}

type SystemRouteListRes struct {
	List []SystemRouteListItemRes `json:"list"`
}

type SystemRouteListItemRes struct {
	Ip      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

type SystemRouteAddReq struct {
	Ip      string `json:"ip" form:"ip" binding:"required"`
	Netmask string `json:"netmask" form:"netmask" binding:"required"`
	Gateway string `json:"gateway" form:"gateway" binding:"required"`
}

type SystemRouteAddRes struct{}

type SystemRouteDeleteReq struct {
	Ip      string `json:"ip" form:"ip" binding:"required"`
	Netmask string `json:"netmask" form:"netmask" binding:"required"`
	Gateway string `json:"gateway" form:"gateway" binding:"required"`
}
type SystemRouteDeleteRes struct{}

type SystemVersionRes struct {
	CurrentVersion    string `json:"currentVersion"`    //当前系统版本
	UpdateTime        string `json:"updateTime"`        //更新时间
	VulUpdateTime     string `json:"vulUpdateTime"`     //漏洞库更新时间
	VulVersion        string `json:"vulVersion"`        //工具库版本
	SysFileName       string `json:"sysFileName"`       //系统更新文件
	VulFileName       string `json:"vulFileName"`       //工具更新文件
	LastSystemVersion string `json:"lastSystemVersion"` //上次系统版本
	LastVulVersion    string `json:"lastVulVersion"`    //上次工具库版本
}

type UploadUpgradeFileRes struct {
	Progress    int    `json:"progress"`
	ZipFilename string `json:"zipFilename"`
}

type SystemOffUpgradeRes struct{}

type SystemManualRollbackReq struct {
	Type string `json:"type" form:"type"` // "SYSTEM" or "VULN"
}

type ConfirmUpgradeReq struct {
	Filename string `json:"filename" binding:"required"`
}

type GenerateTokenReq struct {
	UserName string `json:"username" form:"username" binding:"required"`
}

type GenerateTokenResp struct {
	Token string `json:"token" form:"token" `
}

type TokenListReq struct {
	Page int `json:"page" form:"page" binding:"required"` //页码
	Size int `json:"size" form:"size" binding:"required"` //每页的数量
}

type SystemNodeDownloadReq struct {
	Os string `json:"os" form:"os"` //操作系统
}

type TokenListResp struct {
	List  []TokenListItemRes `json:"list"`
	Total int64              `json:"total"`
}

type TokenListItemRes struct {
	Username   string `json:"username"`
	Token      string `json:"token"`
	CreateTime string `json:"createTime"`
}

/************* 节点管理 ***************/
// 获取 是否启用分布式
type SystemNodeIsDistributeReq struct {
}
type SystemNodeIsDistributeRes struct {
	Status int `json:"status"`
}

// 设置 是否启用分布式
type SystemNodeSetDistributeReq struct {
	Status int `form:"status" json:"status"`
}
type SystemNodeSetDistributeRes struct {
}

// 添加节点
type SystemNodeAddReq struct {
	Ip   string `form:"ip" json:"ip" binding:"required"`
	Port string `form:"port" json:"port" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
}
type SystemNodeAddRes struct {
}

// 节点列表
type SystemNodeListReq struct {
	Page int `json:"page" form:"page" binding:"required"` //页码
	Size int `json:"size" form:"size" binding:"required"` //每页的数量
}
type SystemNodeListRes struct {
	List  []SystemNodeListItem `json:"list"`
	Total int64                `json:"total"`
}
type SystemNodeListItem struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	RunningNum    int    `json:"runningNum"` // 运行任务数
	Status        int    `json:"status"`
	StatusEnum    string `json:"statusEnum"`
	IsDisable     int    `json:"IsDisable"` // 禁用状态 0启用 1禁用
	IsDisableEnum string `json:"isDisableEnum"`
	CreateTime    string `json:"createTime"`
	UpdateTime    string `json:"updateTime"`
}

// 节点删除
type SystemNodeDelReq struct {
	Id string `form:"id" json:"id" binding:"required"`
}
type SystemNodeDelRes struct {
}

// 节点禁用状态设置
type SystemNodeDisOrEnableReq struct {
	Id        int `form:"id" json:"id" binding:"required"`
	IsDisable int `form:"isDisable" json:"isDisable"`
}
type SystemNodeDisOrEnableRes struct {
}

// 获取所有可用节点 - 用于渗透任务时定向渗透使用
type SystemNodeAllEnableRes struct {
	List []SystemNodeAllEnableItem `json:"list"`
}
type SystemNodeAllEnableItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

/************* 节点管理 end ***************/

type SystemMessageListReq struct {
	Page        int    `json:"page" form:"page" binding:"required"` //页码
	Size        int    `json:"size" form:"size" binding:"required"` //每页的数量
	Search      string `json:"search" form:"search"`
	MessageType int    `json:"messageType" form:"messageType"`
	StartTime   string `json:"startTime" form:"startTime"`
	EndTime     string `json:"endTime" form:"endTime"`
}

type SystemMessageListRes struct {
	List []SystemMessageListResItem `json:"list"`
}

type SystemMessageListResItem struct {
	Id              int    `json:"id"`
	Content         string `json:"content"`
	MessageType     int    `json:"messageType"`
	MessageTypeEnum string `json:"messageTypeEnum"`
	Status          int    `json:"status"`
	StatusEnum      string `json:"StatusEnum"`
	CreateTime      string `json:"createTime"`
	UserId          int    `json:"userId"`
}

type SystemMonitorWarnSaveReq struct {
	IsOpen     int `json:"isOpen" form:"isOpen" binding:"required"`
	CpuWarn    int `json:"cpuWarn" form:"cpuWarn" binding:"required"`
	MemoryWarn int `json:"memoryWarn" form:"memoryWarn" binding:"required"`
	DiskWarn   int `json:"diskWarn" form:"diskWarn" binding:"required"`
	FlowWarn   int `json:"flowWarn" form:"flowWarn" binding:"required"`
}

type SystemMonitorWarnSaveRes struct{}

type SystemMonitorWarnInfoRes struct {
	IsOpen     int `json:"isOpen"`
	CpuWarn    int `json:"cpuWarn"`
	MemoryWarn int `json:"memoryWarn"`
	DiskWarn   int `json:"diskWarn"`
	FlowWarn   int `json:"flowWarn"`
}

type UseScoreInfoRes struct {
	IsOpen int `json:"isOpen"`
}

type UseScoreSaveReq struct {
	IsOpen int `json:"isOpen" form:"isOpen" binding:"required"`
}

type UseScoreSaveRes struct{}

type TestScopeInfoRes struct {
	IsOpen int `json:"isOpen"`
}

type TestScopeSaveReq struct {
	IsOpen int `json:"isOpen" form:"isOpen" binding:"required"`
}

type TestScopeSaveRes struct{}

type TokenDeleteReq struct {
	Username string `json:"username" form:"username" binding:"required"` //用户名
	Token    string `json:"token" form:"token" binding:"required"`       //秘钥值
}

type TokenDeleteResp struct {
}
