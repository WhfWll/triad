package enums

import "strconv"

// 类型枚举
const (
	DictionaryTypeUser          = 1 // 用户字典
	DictionaryTypePassword      = 2 // 密码字典
	DictionaryTypeWifi          = 3 // wifi字典
	DictionaryTypeWebPathScan   = 4 // web路径爆破字典
	DictionaryTypeSubdomainScan = 5 // 子域名爆破字典
)

// 适用范围枚举
const (
	// 用户字典和密码字典适用范围
	DictionaryServiceWeakPassCommon    = 1
	DictionaryServiceWeakPassSsh       = 2
	DictionaryServiceWeakPassFtp       = 3
	DictionaryServiceWeakPassMemcached = 4
	DictionaryServiceWeakPassMongodb   = 5
	DictionaryServiceWeakPassMssql     = 6
	DictionaryServiceWeakPassMysql     = 7
	DictionaryServiceWeakPassPostgres  = 8
	DictionaryServiceWeakPassRdp       = 9
	DictionaryServiceWeakPassRedis     = 10
	DictionaryServiceWeakPassSmb       = 11
	DictionaryServiceWeakPassTelnet    = 12
	DictionaryServiceWeakPassTomcat    = 13
	DictionaryServiceWeakPassVnc       = 14

	DictionaryServiceWeakPassSnmpv2       = 15
	DictionaryServiceWeakPassSnmpv3Md5    = 16
	DictionaryServiceWeakPassSnmpv3Sha    = 17
	DictionaryServiceWeakPassSnmpv3Sha224 = 18
	DictionaryServiceWeakPassSnmpv3Sha256 = 19
	DictionaryServiceWeakPassSnmpv3Sha384 = 20
	DictionaryServiceWeakPassSnmpv3Sha512 = 21
	DictionaryServiceWeakPassHttp         = 22
	DictionaryServiceWeakPassOracle       = 23
	DictionaryServiceWeakPassLdap         = 24
	DictionaryServiceWeakPassSmtp         = 25
	DictionaryServiceWeakPassPop3         = 26
	DictionaryServiceWeakPassImap         = 27

	// web路径爆破适用范围
	DictionaryServiceWebPathScanCommon = 101
	DictionaryServiceWebPathScanPhp    = 102
	DictionaryServiceWebPathScanAsp    = 103
	DictionaryServiceWebPathScanAspx   = 104
	DictionaryServiceWebPathScanJsp    = 105
	// 子域名爆破适用范围
	DictionaryServiceSubdomainScanCommon = 201
	// wifi适用范围(从301开始写)
)

// LangKeyValServiceID 语言判定关键词 关键词->字典service 获取对应的爆破路径出来
var LangKeyValServiceID = map[string]int{
	"php":  DictionaryServiceWebPathScanPhp,
	"asp":  DictionaryServiceWebPathScanAsp,
	"aspx": DictionaryServiceWebPathScanAspx,
	"jsp":  DictionaryServiceWebPathScanJsp,
}

// 是否默认枚举
const (
	DictionaryDefaultYes = 1 // 是默认的
	DictionaryDefaultNo  = 2 // 不是默认的
)

// 来源枚举
const (
	DictionarySourceSystem = 1 // 系统
	DictionarySourceManual = 2 // 手动
)

// DictionaryTypeEnum 获取字典类型枚举
func DictionaryTypeEnum() map[int]string {
	typeEnum := map[int]string{
		DictionaryTypeUser:          "用户字典",
		DictionaryTypePassword:      "密码字典",
		DictionaryTypeWifi:          "wifi字典",
		DictionaryTypeWebPathScan:   "web路径爆破字典",
		DictionaryTypeSubdomainScan: "子域名爆破字典",
	}
	return typeEnum
}

// GetDictionaryType 获取字典类型的值
func GetDictionaryType(k int) string {
	typeEnum := DictionaryTypeEnum()
	if value, ok := typeEnum[k]; ok {
		return value
	}
	return ""
}

// GetDictionaryTypeArray 字典类型枚举数据
func GetDictionaryTypeArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{
		{
			Value: DictionaryTypeUser,
			Label: GetDictionaryType(DictionaryTypeUser),
		},
		{
			Value: DictionaryTypePassword,
			Label: GetDictionaryType(DictionaryTypePassword),
		},
		{
			Value: DictionaryTypeWebPathScan,
			Label: GetDictionaryType(DictionaryTypeWebPathScan),
		},
		{
			Value: DictionaryTypeSubdomainScan,
			Label: GetDictionaryType(DictionaryTypeSubdomainScan),
		},
	}
	return result
}

// DictionaryServiceEnum 获取字典适用范围枚举
func DictionaryServiceEnum() map[int]string {
	serviceEnum := map[int]string{
		DictionaryServiceWeakPassCommon:       "通用",
		DictionaryServiceWeakPassSsh:          "ssh",
		DictionaryServiceWeakPassFtp:          "ftp",
		DictionaryServiceWeakPassMemcached:    "memcached",
		DictionaryServiceWeakPassMongodb:      "mongodb",
		DictionaryServiceWeakPassMssql:        "mssql",
		DictionaryServiceWeakPassMysql:        "mysql",
		DictionaryServiceWeakPassPostgres:     "postgres",
		DictionaryServiceWeakPassRdp:          "rdp",
		DictionaryServiceWeakPassRedis:        "redis",
		DictionaryServiceWeakPassSmb:          "smb",
		DictionaryServiceWeakPassTelnet:       "telnet",
		DictionaryServiceWeakPassTomcat:       "tomcat",
		DictionaryServiceWeakPassVnc:          "vnc",
		DictionaryServiceWeakPassSnmpv2:       "snmpv2",
		DictionaryServiceWeakPassSnmpv3Md5:    "snmpv3_md5",
		DictionaryServiceWeakPassSnmpv3Sha:    "snmpv3_sha",
		DictionaryServiceWeakPassSnmpv3Sha224: "snmpv3_sha_224",
		DictionaryServiceWeakPassSnmpv3Sha256: "snmpv3_sha_256",
		DictionaryServiceWeakPassSnmpv3Sha384: "snmpv3_sha_384",
		DictionaryServiceWeakPassSnmpv3Sha512: "snmpv3_sha_512",
		DictionaryServiceWeakPassHttp:         "HTTP",

		DictionaryServiceWebPathScanCommon: "通用",
		DictionaryServiceWebPathScanPhp:    "php",
		DictionaryServiceWebPathScanAsp:    "asp",
		DictionaryServiceWebPathScanAspx:   "aspx",
		DictionaryServiceWebPathScanJsp:    "jsp",

		DictionaryServiceSubdomainScanCommon: "通用",
	}
	return serviceEnum
}

// DictionaryServiceEnum 获取字典适用范围枚举
func DictionaryServiceToScripNameEnum(services []string) []string {
	var result = make([]string, 0)
	serviceEnum := map[int]string{
		DictionaryServiceWeakPassCommon:       "bf_weak_common",
		DictionaryServiceWeakPassSsh:          "bf_ssh",
		DictionaryServiceWeakPassFtp:          "bf_ftp",
		DictionaryServiceWeakPassMemcached:    "bf_memcached",
		DictionaryServiceWeakPassMongodb:      "bf_mongodb",
		DictionaryServiceWeakPassMssql:        "bf_mssql",
		DictionaryServiceWeakPassMysql:        "bf_mysql",
		DictionaryServiceWeakPassPostgres:     "bf_postgres",
		DictionaryServiceWeakPassRdp:          "bf_rdp",
		DictionaryServiceWeakPassRedis:        "bf_redis",
		DictionaryServiceWeakPassSmb:          "bf_smb",
		DictionaryServiceWeakPassTelnet:       "bf_telnet",
		DictionaryServiceWeakPassTomcat:       "bf_tomcat",
		DictionaryServiceWeakPassVnc:          "bf_vnc",
		DictionaryServiceWeakPassSnmpv2:       "bf_snmpv2",
		DictionaryServiceWeakPassSnmpv3Md5:    "bf_snmpv3_md5",
		DictionaryServiceWeakPassSnmpv3Sha:    "bf_snmpv3_sha",
		DictionaryServiceWeakPassSnmpv3Sha224: "bf_snmpv3_sha_224",
		DictionaryServiceWeakPassSnmpv3Sha256: "bf_snmpv3_sha_256",
		DictionaryServiceWeakPassSnmpv3Sha384: "bf_snmpv3_sha_384",
		DictionaryServiceWeakPassSnmpv3Sha512: "bf_snmpv3_sha_512",
		DictionaryServiceWebPathScanCommon:    "bf_webpath_common",
		DictionaryServiceWebPathScanPhp:       "bf_php",
		DictionaryServiceWebPathScanAsp:       "bf_asp",
		DictionaryServiceWebPathScanAspx:      "bf_aspx",
		DictionaryServiceWebPathScanJsp:       "bf_jsp",
		DictionaryServiceSubdomainScanCommon:  "bf_subdomain_common",
	}
	for i := 0; i < len(services); i++ {
		tmpSer, _ := strconv.Atoi(services[i])
		if v, ok := serviceEnum[tmpSer]; ok {
			result = append(result, v)
		}
	}
	return result
}

// GetDictionaryService 获取字典适用范围的值
func GetDictionaryService(k int) string {
	serviceEnum := DictionaryServiceEnum()
	if value, ok := serviceEnum[k]; ok {
		return value
	}
	return ""
}

// GetDictionaryServiceWeakPassArray 弱口令爆破枚举数据
func GetDictionaryServiceWeakPassArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{
		{
			Value: DictionaryServiceWeakPassCommon,
			Label: GetDictionaryService(DictionaryServiceWeakPassCommon),
		},
		{
			Value: DictionaryServiceWeakPassSsh,
			Label: GetDictionaryService(DictionaryServiceWeakPassSsh),
		},
		{
			Value: DictionaryServiceWeakPassFtp,
			Label: GetDictionaryService(DictionaryServiceWeakPassFtp),
		},
		{
			Value: DictionaryServiceWeakPassMemcached,
			Label: GetDictionaryService(DictionaryServiceWeakPassMemcached),
		},
		{
			Value: DictionaryServiceWeakPassMongodb,
			Label: GetDictionaryService(DictionaryServiceWeakPassMongodb),
		},
		{
			Value: DictionaryServiceWeakPassMssql,
			Label: GetDictionaryService(DictionaryServiceWeakPassMssql),
		},
		{
			Value: DictionaryServiceWeakPassMysql,
			Label: GetDictionaryService(DictionaryServiceWeakPassMysql),
		},
		{
			Value: DictionaryServiceWeakPassPostgres,
			Label: GetDictionaryService(DictionaryServiceWeakPassPostgres),
		},
		{
			Value: DictionaryServiceWeakPassRdp,
			Label: GetDictionaryService(DictionaryServiceWeakPassRdp),
		},
		{
			Value: DictionaryServiceWeakPassRedis,
			Label: GetDictionaryService(DictionaryServiceWeakPassRedis),
		},
		{
			Value: DictionaryServiceWeakPassSmb,
			Label: GetDictionaryService(DictionaryServiceWeakPassSmb),
		},
		{
			Value: DictionaryServiceWeakPassTelnet,
			Label: GetDictionaryService(DictionaryServiceWeakPassTelnet),
		},
		{
			Value: DictionaryServiceWeakPassTomcat,
			Label: GetDictionaryService(DictionaryServiceWeakPassTomcat),
		},
		{
			Value: DictionaryServiceWeakPassVnc,
			Label: GetDictionaryService(DictionaryServiceWeakPassVnc),
		},
		{
			Value: DictionaryServiceWeakPassSnmpv2,
			Label: GetDictionaryService(DictionaryServiceWeakPassSnmpv2),
		},
		{
			Value: DictionaryServiceWeakPassSnmpv3Md5,
			Label: GetDictionaryService(DictionaryServiceWeakPassSnmpv3Md5),
		},
		{
			Value: DictionaryServiceWeakPassSnmpv3Sha,
			Label: GetDictionaryService(DictionaryServiceWeakPassSnmpv3Sha),
		},
		{
			Value: DictionaryServiceWeakPassSnmpv3Sha224,
			Label: GetDictionaryService(DictionaryServiceWeakPassSnmpv3Sha224),
		},
		{
			Value: DictionaryServiceWeakPassSnmpv3Sha256,
			Label: GetDictionaryService(DictionaryServiceWeakPassSnmpv3Sha256),
		},
		{
			Value: DictionaryServiceWeakPassSnmpv3Sha384,
			Label: GetDictionaryService(DictionaryServiceWeakPassSnmpv3Sha384),
		},
		{
			Value: DictionaryServiceWeakPassSnmpv3Sha512,
			Label: GetDictionaryService(DictionaryServiceWeakPassSnmpv3Sha512),
		},
		{
			Value: DictionaryServiceWeakPassHttp,
			Label: GetDictionaryService(DictionaryServiceWeakPassHttp),
		},
	}
	return result
}

// GetDictionaryServiceWebPathScanArray web路径爆破枚举数据
func GetDictionaryServiceWebPathScanArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{
		{
			Value: DictionaryServiceWebPathScanCommon,
			Label: GetDictionaryService(DictionaryServiceWebPathScanCommon),
		},
		{
			Value: DictionaryServiceWebPathScanPhp,
			Label: GetDictionaryService(DictionaryServiceWebPathScanPhp),
		},
		{
			Value: DictionaryServiceWebPathScanAsp,
			Label: GetDictionaryService(DictionaryServiceWebPathScanAsp),
		},
		{
			Value: DictionaryServiceWebPathScanAspx,
			Label: GetDictionaryService(DictionaryServiceWebPathScanAspx),
		},
		{
			Value: DictionaryServiceWebPathScanJsp,
			Label: GetDictionaryService(DictionaryServiceWebPathScanJsp),
		},
	}
	return result
}

// GetDictionaryServiceSubdomainScanArray 子域名爆破枚举数据
func GetDictionaryServiceSubdomainScanArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{
		{
			Value: DictionaryServiceSubdomainScanCommon,
			Label: GetDictionaryService(DictionaryServiceSubdomainScanCommon),
		},
	}
	return result
}

// DictionaryDefaultEnum 获取字典是否默认枚举
func DictionaryDefaultEnum() map[int]string {
	defaultEnum := map[int]string{
		DictionaryDefaultYes: "是",
		DictionaryDefaultNo:  "否",
	}
	return defaultEnum
}

// GetDictionaryDefault 获取字典是否默认的值
func GetDictionaryDefault(k int) string {
	defaultEnum := DictionaryDefaultEnum()
	if value, ok := defaultEnum[k]; ok {
		return value
	}
	return ""
}

// DictionarySourceEnum 获取字典来源枚举
func DictionarySourceEnum() map[int]string {
	sourceEnum := map[int]string{
		DictionarySourceSystem: "系统",
		DictionarySourceManual: "手动",
	}
	return sourceEnum
}

// GetDictionarySource 获取字典来源的值
func GetDictionarySource(k int) string {
	sourceEnum := DictionarySourceEnum()
	if value, ok := sourceEnum[k]; ok {
		return value
	}
	return ""
}

// AllDictionaryServiceWeakPassEnum 获取弱口令字典适用范围枚举
func AllDictionaryServiceWeakPassEnum(hasCommon bool) map[int]string {
	serviceEnum := map[int]string{
		DictionaryServiceWeakPassSsh:          "ssh",
		DictionaryServiceWeakPassFtp:          "ftp",
		DictionaryServiceWeakPassMemcached:    "memcached",
		DictionaryServiceWeakPassMongodb:      "mongodb",
		DictionaryServiceWeakPassMssql:        "mssql",
		DictionaryServiceWeakPassMysql:        "mysql",
		DictionaryServiceWeakPassPostgres:     "postgres",
		DictionaryServiceWeakPassRdp:          "rdp",
		DictionaryServiceWeakPassRedis:        "redis",
		DictionaryServiceWeakPassSmb:          "smb",
		DictionaryServiceWeakPassTelnet:       "telnet",
		DictionaryServiceWeakPassTomcat:       "tomcat",
		DictionaryServiceWeakPassVnc:          "vnc",
		DictionaryServiceWeakPassSnmpv2:       "snmpv2",
		DictionaryServiceWeakPassSnmpv3Md5:    "snmpv3_md5",
		DictionaryServiceWeakPassSnmpv3Sha:    "snmpv3_sha",
		DictionaryServiceWeakPassSnmpv3Sha224: "snmpv3_sha_224",
		DictionaryServiceWeakPassSnmpv3Sha256: "snmpv3_sha_256",
		DictionaryServiceWeakPassSnmpv3Sha384: "snmpv3_sha_384",
		DictionaryServiceWeakPassSnmpv3Sha512: "snmpv3_sha_512",
		DictionaryServiceWeakPassHttp:         "HTTP",
		DictionaryServiceWeakPassOracle:       "oracle",
		DictionaryServiceWeakPassLdap:         "ldap",
		DictionaryServiceWeakPassSmtp:         "smtp",
		DictionaryServiceWeakPassPop3:         "pop3",
		DictionaryServiceWeakPassImap:         "imap",
	}
	if hasCommon {
		serviceEnum[DictionaryServiceWeakPassCommon] = "通用"
	}
	return serviceEnum
}

// 字典类型
//const (
//	TaskConfigurationWeakPassDictTypeDefault = 1 // 默认字典
//	TaskConfigurationWeakPassDictTypeCommon  = 2 // 通用字典
//	TaskConfigurationWeakPassDictTypeAdd     = 3 // 补充字典
//)
