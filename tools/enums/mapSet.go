package enums

type SystemsSystemSecurity struct {
	SystemSecurityPasswordRank          int `json:"systemSecurityPasswordRank"`          // 密码等级
	SystemSecurityPasswordCycle         int `json:"systemSecurityPasswordCycle"`         // 密码周期
	SystemSecurityLoginTimeout          int `json:"systemSecurityLoginTimeout"`          // 登录超时（多长时间未操作自动退出）
	SystemSecurityUserLimit             int `json:"systemSecurityUserLimit"`             // 普通用户限制登录次数（到达次数直接禁用）
	SystemSecurityAdminLimit            int `json:"systemSecurityAdminLimit"`            // 用户员/审计员限制登录次数（到达次数限制登录时间）
	SystemSecurityBanTime               int `json:"systemSecurityBanTime"`               // 用户员/审计员禁止登录时间
	SystemSecurityAccountNoLoginDisable int `json:"systemSecurityAccountNoLoginDisable"` // 账户多长时间无登录禁用
	SystemSecurityPasswordValid         int `json:"systemSecurityPasswordValid"`         // 密码有效期时长设置
	SystemSecurityExpireUnused          int `json:"systemSecurityExpireUnused"`          // 密码有效期时长设置
}

/******************* 密码强度 *******************/
// 低：长度不低于8位，可为全字母或数字
const SystemSecurityPassRankLow = 1

// 中：长度不低于8位，包含字母、符号、数字中的两类
const SystemSecurityPassRankMiddle = 2

// 高：长度不低于8位，包含字母、符号、数字中的三类
const SystemSecurityPassRankHigh = 3

func (SystemSecurityEnum SystemsSystemSecurity) GetPassRandEnum(rank int) string {
	passRankChoice := map[int]string{
		SystemSecurityPassRankLow:    "低",
		SystemSecurityPassRankMiddle: "中",
		SystemSecurityPassRankHigh:   "高",
	}

	if res, ok := passRankChoice[rank]; ok {
		return res
	}

	return ""
}

const (
	ReportContentMapSetKey                = "reportContent"      //报告内容枚举
	TargetIpWhiteBlackMapSetObjKey        = "targetIpWhiteBlack" //测试目标黑白名单
	TargetIpWhiteBlackMapSetContent       = "测试目标黑白名单"
	LogBackupConfigMapSetObjKey           = "logBackupConfig"
	LogBackupConfigMapSetContent          = "日志备份配置"
	LogExpTimeConfigMapSetObjKey          = "logExpTimeConfig"
	LogExpTimeConfigMapSetContent         = "日志过期时间"
	SystemConfigBackupConfigMapSetObjKey  = "systemConfigBackupConfig"
	SystemConfigBackupConfigMapSetContent = "系统配置备份配置"
	ReverseIpHostMapSetObjKey             = "reverseIpHost"
	ReverseIpHostMapSetContent            = "远程监听"
	SystemAccessIpWhiteMapSetObjKey       = "systemAccessIpWhite"
	SystemAccessIpWhiteMapSetContent      = "系统访问白名单"
	SyslogConfigMapSetObjKey              = "syslogConfig"
	SyslogConfigMapSetContent             = "syslog配置"
	MailConfigMapSetObjKey                = "mailConfig"
	MailConfigMapSetContent               = "邮箱配置"
	NetworkConfigMapSetObjKey             = "networkConfig"
	NetworkConfigMapSetContent            = "网络配置"
	TcpBlindTestMapSetObjKey              = "tcpBlindTest"
	TcpBlindTestContent                   = "TCP盲测平台"
	HttpBlindTestMapSetObjKey             = "httpBlindTest"
	HttpBlindTestContent                  = "http盲测平台"
	DnsBlindTestMapSetObjKey              = "dnsBlindTest"
	DnsBlindTestContent                   = "DNS盲测平台"
	IcmpBlindTestMapSetObjKey             = "icmpBlindTest"
	IcmpBlindTestContent                  = "ICMP盲测平台"
	MonitorCpuMapSetObjKey                = "monitorCpu"
	MonitorCpuMapSetContent               = "最近一小时CPU使用率"
	MonitorMemoryMapSetObjKey             = "monitorMemory"
	MonitorMemoryMapSetContent            = "最近一小时内存使用率"
	MonitorDiskMapSetObjKey               = "monitorDisk"
	MonitorDiskMapSetContent              = "硬盘使用情况"
	FlowTaskNetworkCardMapSetKey          = "flowTaskNetworkCard"
	SystemVersionMapSetObjKey             = "systemVersion"
	SystemVersionMapSetContent            = "系统版本信息"
	BusinessUserTypeAuthMapSetObjKey      = "businessUserTypeAuth"
	BusinessUserTypeAuthMapSetContent     = "业务模块和用户类型对应关系"
	CurTasksInfoMapSetObjKey              = "curTasksInfo"
	CurTasksInfoMapSetContent             = "任务并发配置"
	MonitorWarnMapSetObjKey               = "monitorWarn"
	MonitorWarnMapSetContent              = "系统监控告警"
	UseScoreSwitchMapSetObjKey            = "useScoreSwitch"
	UseScoreSwitchMapSetContent           = "可利用评分开关，1-打开，2-关闭"
	TestScopeSwitchMapSetObjKey           = "testScopeSwitch"
	TestScopeSwitchMapSetContent          = "测试范围校验开关，1-打开，2-关闭"
	TestScopeSwitchOpen                   = "1"
	TestScopeSwitchClose                  = "2"
	VulEvidenceInfoLeakMapSetObjKey       = "vulEvidenceInfoLeak"
	VulEvidenceFileLeakMapSetObjKey       = "vulEvidenceFileLeak"
	SystemSecurityPasswordValid           = "systemSecurityPasswordValid" //密码有效期枚举值
	SystemSecurityExpireUnused            = "systemSecurityExpireUnused"
	AssetDetectProgressMapSetObjKey       = "assetDetectProgress"
)

const (
	MapSetEstateValid   = "valid"   //MapSet为有效状态
	MapSetEstateDeleted = "deleted" //MapSet为已删除状态
)

const (
	FlowTaskYakitScriptMapSetObjKey  = "flowTaskYakitScript" //被动流量服务yak脚本
	FlowTaskYakitScriptMapSetContent = "被动流量服务yak脚本"

	MapSetYakNodeObjKey                 = "yakNodeDistribute"
	MapSetYakNodeContent                = "是否开启分布式引擎"
	MapSetYakNodeValueN                 = "0"                       // 未开启分布式引擎
	MapSetYakNodeValueY                 = "1"                       // 开启分布式引擎
	FlowTaskContentRulesMapSetObjKey    = "flowTaskContentRules"    //被动流量内容规则
	UseScoreFingerBlackListMapSetObjKey = "useScoreFingerBlackList" //利用评分指纹黑名单
)

// llm use
const (
	LlmModelEnhancementObjKey  = "llmModelEnhancement"
	LlmModelEnhancementContent = "开启模型增强"
)
