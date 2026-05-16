package enums

// 基线检查 - OS类型
const (
	BaselineOSTypeLinux    = 1 // Linux/Unix
	BaselineOSTypeWindows  = 2 // Windows
	BaselineOSTypeDomestic = 3 // 国产OS (Kylin/UOS)
	BaselineOSTypeEmbedded = 4 // 嵌入式OS
)

// 主机远程连接方式（SSH / WinRM；用于安全配置核查等）
const (
	HostTransportAuto  = 0 // 按操作系统：Windows→WinRM，其它→SSH
	HostTransportSSH   = 1
	HostTransportWinRM = 2
)

// 主机远程检查场景（写入 baseline_check_result.scan_scene；漏洞检测与配置核查共用规则库）
const (
	HostScanSceneBaseline = 1 // 安全配置核查
	HostScanSceneHostVuln = 2 // 主机漏洞检测
)

// 基线检查 - 检测结果
const (
	BaselineCheckResultPass  = 1 // 通过
	BaselineCheckResultFail  = 2 // 不通过
	BaselineCheckResultError = 3 // 检查出错
	BaselineCheckResultSkip  = 4 // 跳过
)

// 基线检查 - 规则分类
const (
	BaselineCategoryPasswordPolicy    = 1  // 密码策略
	BaselineCategoryUserPermission    = 2  // 用户权限
	BaselineCategoryFirewall          = 3  // 防火墙规则
	BaselineCategoryKernelSecurity    = 4  // 内核安全
	BaselineCategoryFilePermission    = 5  // 文件权限
	BaselineCategoryAuditLog          = 6  // 审计日志
	BaselineCategoryNetworkService    = 7  // 网络服务
	BaselineCategorySystemUpdate      = 8  // 系统更新
	BaselineCategorySSHConfig         = 9  // SSH配置
	BaselineCategoryPasswordPolicyWin = 10 // Windows密码策略
	BaselineCategoryUserPermissionWin = 11 // Windows用户权限
	BaselineCategoryFirewallWin       = 12 // Windows防火墙
	BaselineCategoryServiceWin        = 13 // Windows服务
	BaselineCategoryAuditPolicyWin    = 14 // Windows审计策略
	BaselineCategoryOther             = 99 // 其他
)

// 基线检查 - 风险等级
const (
	BaselineRiskCritical = 0 // 严重
	BaselineRiskHigh     = 1 // 高危
	BaselineRiskMiddle   = 2 // 中危
	BaselineRiskLow      = 3 // 低危
	BaselineRiskInfo     = 4 // 信息
)

// 恶意代码检测 - 检测类型
const (
	MalwareCheckTypeFileIntegrity = 1 // 文件完整性
	MalwareCheckTypeProcess       = 2 // 进程检测
	MalwareCheckTypeNetwork       = 3 // 网络连接
	MalwareCheckTypeWebshell      = 4 // Webshell检测
	MalwareCheckTypeRootkit       = 5 // Rootkit检测
	MalwareCheckTypeSensitiveFile = 6 // 敏感文件
)

// 恶意代码检测 - 风险等级
const (
	MalwareRiskCritical = 1 // 严重
	MalwareRiskHigh     = 2 // 高危
	MalwareRiskMiddle   = 3 // 中危
	MalwareRiskLow      = 4 // 低危
)

// 数据库类型
const (
	DBSupportTypeMySQL      = 1 // MySQL
	DBSupportTypePostgreSQL = 2 // PostgreSQL
	DBSupportTypeMongoDB    = 3 // MongoDB
	DBSupportTypeRedis      = 4 // Redis
	DBSupportTypeCouchDB    = 5 // CouchDB
)

// 敏感数据 - 级别
const (
	SensitiveDataLevelHigh   = 1 // 高敏
	SensitiveDataLevelMiddle = 2 // 中敏
	SensitiveDataLevelLow    = 3 // 低敏
)

// 敏感数据 - 分类
const (
	SensitiveDataTypeIDCard      = 1  // 身份证号
	SensitiveDataTypeBankCard    = 2  // 银行卡号
	SensitiveDataTypePassport    = 3  // 护照号
	SensitiveDataTypePhone       = 4  // 手机号
	SensitiveDataTypeEmail       = 5  // 邮箱
	SensitiveDataTypeAddress     = 6  // 地址
	SensitiveDataTypeBirthDate   = 7  // 出生日期
	SensitiveDataTypeName        = 8  // 姓名
	SensitiveDataTypeToken       = 9  // Token/密钥
	SensitiveDataTypeCertificate = 10 // 证书
	SensitiveDataTypePassword    = 11 // 密码hash
	SensitiveDataTypeOther       = 99 // 其他
)

// 数据库配置核查 - 检查维度
const (
	DBCheckCategoryAuthentication = 1 // 身份认证
	DBCheckCategoryAuthorization  = 2 // 权限控制
	DBCheckCategoryConfigSecure   = 3 // 配置安全
	DBCheckCategoryAuditLog       = 4 // 审计日志
	DBCheckCategoryNetwork        = 5 // 网络安全
	DBCheckCategoryEncryption     = 6 // 加密
)

type baseline struct{}

var BaselineEnum baseline

func (baseline) GetOSTypeName(osType int) string {
	m := map[int]string{
		BaselineOSTypeLinux:    "Linux/Unix",
		BaselineOSTypeWindows:  "Windows",
		BaselineOSTypeDomestic: "国产操作系统",
		BaselineOSTypeEmbedded: "嵌入式OS",
	}
	if v, ok := m[osType]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetHostScanSceneName(scene int) string {
	switch scene {
	case HostScanSceneHostVuln:
		return "主机漏洞检测"
	case HostScanSceneBaseline:
		return "安全配置核查"
	default:
		return "安全配置核查"
	}
}

func (baseline) GetHostTransportName(t int) string {
	m := map[int]string{
		HostTransportAuto:  "自动（按操作系统）",
		HostTransportSSH:  "SSH",
		HostTransportWinRM: "WinRM",
	}
	if v, ok := m[t]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetCheckResultName(result int) string {
	m := map[int]string{
		BaselineCheckResultPass:  "通过",
		BaselineCheckResultFail:  "不通过",
		BaselineCheckResultError: "错误",
		BaselineCheckResultSkip:  "跳过",
	}
	if v, ok := m[result]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetCategoryName(category int) string {
	m := map[int]string{
		BaselineCategoryPasswordPolicy:    "密码策略",
		BaselineCategoryUserPermission:    "用户权限",
		BaselineCategoryFirewall:          "防火墙规则",
		BaselineCategoryKernelSecurity:    "内核安全",
		BaselineCategoryFilePermission:    "文件权限",
		BaselineCategoryAuditLog:          "审计日志",
		BaselineCategoryNetworkService:    "网络服务",
		BaselineCategorySystemUpdate:      "系统更新",
		BaselineCategorySSHConfig:         "SSH配置",
		BaselineCategoryPasswordPolicyWin: "密码策略",
		BaselineCategoryUserPermissionWin: "用户权限",
		BaselineCategoryFirewallWin:       "防火墙",
		BaselineCategoryServiceWin:        "服务配置",
		BaselineCategoryAuditPolicyWin:    "审计策略",
		BaselineCategoryOther:             "其他",
	}
	if v, ok := m[category]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetBaselineRiskName(risk int) string {
	m := map[int]string{
		BaselineRiskCritical: "严重",
		BaselineRiskHigh:     "高危",
		BaselineRiskMiddle:   "中危",
		BaselineRiskLow:      "低危",
		BaselineRiskInfo:     "信息",
	}
	if v, ok := m[risk]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetMalwareRiskName(risk int) string {
	m := map[int]string{
		MalwareRiskCritical: "严重",
		MalwareRiskHigh:     "高危",
		MalwareRiskMiddle:   "中危",
		MalwareRiskLow:      "低危",
	}
	if v, ok := m[risk]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetDBTypeName(dbType int) string {
	m := map[int]string{
		DBSupportTypeMySQL:      "MySQL",
		DBSupportTypePostgreSQL: "PostgreSQL",
		DBSupportTypeMongoDB:    "MongoDB",
		DBSupportTypeRedis:      "Redis",
		DBSupportTypeCouchDB:    "CouchDB",
	}
	if v, ok := m[dbType]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetDBCheckCategoryName(c int) string {
	m := map[int]string{
		DBCheckCategoryAuthentication: "身份认证",
		DBCheckCategoryAuthorization:  "权限控制",
		DBCheckCategoryConfigSecure:   "配置安全",
		DBCheckCategoryAuditLog:       "审计日志",
		DBCheckCategoryNetwork:        "网络安全",
		DBCheckCategoryEncryption:     "加密",
	}
	if v, ok := m[c]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetSensitiveDataLevelName(level int) string {
	m := map[int]string{
		SensitiveDataLevelHigh:   "高敏",
		SensitiveDataLevelMiddle: "中敏",
		SensitiveDataLevelLow:    "低敏",
	}
	if v, ok := m[level]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetSensitiveDataTypeName(t int) string {
	m := map[int]string{
		SensitiveDataTypeIDCard:      "身份证号",
		SensitiveDataTypeBankCard:    "银行卡号",
		SensitiveDataTypePassport:    "护照号",
		SensitiveDataTypePhone:       "手机号",
		SensitiveDataTypeEmail:       "邮箱",
		SensitiveDataTypeAddress:     "地址",
		SensitiveDataTypeBirthDate:   "出生日期",
		SensitiveDataTypeName:        "姓名",
		SensitiveDataTypeToken:       "Token/密钥",
		SensitiveDataTypeCertificate: "证书",
		SensitiveDataTypePassword:    "密码hash",
		SensitiveDataTypeOther:       "其他",
	}
	if v, ok := m[t]; ok {
		return v
	}
	return "未知"
}

func (baseline) GetMalwareCheckTypeName(t int) string {
	m := map[int]string{
		MalwareCheckTypeFileIntegrity: "文件完整性",
		MalwareCheckTypeProcess:       "进程检测",
		MalwareCheckTypeNetwork:       "网络连接",
		MalwareCheckTypeWebshell:      "Webshell检测",
		MalwareCheckTypeRootkit:       "Rootkit检测",
		MalwareCheckTypeSensitiveFile: "敏感文件",
	}
	if v, ok := m[t]; ok {
		return v
	}
	return "未知"
}
