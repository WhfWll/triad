package enums

const (
	VulVerifyTypePrincipleVerification = 1 //原理验证
	VulVerifyTypeVersionMatch          = 2 //版本匹配
)

// GetTaskVulVerifyTypeEnum 获取任务漏洞验证方式枚举
func GetTaskVulVerifyTypeEnum(state int) string {
	enums := map[int]string{
		VulVerifyTypeVersionMatch:          "版本匹配", //版本匹配
		VulVerifyTypePrincipleVerification: "原理验证", //原理验证
	}
	if res, ok := enums[state]; ok {
		return res
	}
	return ""
}

func AllTaskVulVerifyTypeEnum() map[int]string {
	enum := map[int]string{
		VulVerifyTypeVersionMatch:          "版本匹配", //版本匹配
		VulVerifyTypePrincipleVerification: "原理验证", //原理验证
	}
	return enum
}

// GetTaskVulVerifyTypeArray 获取漏洞类型
func GetTaskVulVerifyTypeArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{
		{
			Value: VulVerifyTypeVersionMatch,
			Label: GetTaskVulVerifyTypeEnum(VulVerifyTypeVersionMatch),
		},
		{
			Value: VulVerifyTypePrincipleVerification,
			Label: GetTaskVulVerifyTypeEnum(VulVerifyTypePrincipleVerification),
		},
	}
	return result
}

// VulTypeOption 漏洞类型选项
type VulTypeOption struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// GetVulTypeArray 获取所有漏洞类型数组
func GetVulTypeArray() []VulTypeOption {
	typeEnum := map[int]string{
		//VulLibrariesTypeNot:                     "无",
		VulLibrariesTypeXss:                     "跨站脚本",
		VulLibrariesTypeCsrf:                    "跨站请求伪造",
		VulLibrariesTypeSqlInj:                  "Sql注入",
		VulLibrariesTypeLdapInj:                 "ldap注入",
		VulLibrariesTypeSmtpInj:                 "邮件命令注入",
		VulLibrariesTypeNullByteInj:             "空字节注入",
		VulLibrariesTypeCrlfInj:                 "CRLF注入",
		VulLibrariesTypeSsiInj:                  "Ssi注入",
		VulLibrariesTypeXpathInj:                "Xpath注入",
		VulLibrariesTypeXmlInj:                  "Xml注入",
		VulLibrariesTypeXqueryInj:               "Xquery 注入",
		VulLibrariesTypeCodeExec:                "代码执行",
		VulLibrariesTypeRfi:                     "远程文件包含",
		VulLibrariesTypeLfi:                     "本地文件包含",
		VulLibrariesTypeFuncAbuse:               "功能函数滥用",
		VulLibrariesTypeBruteForce:              "暴力破解",
		VulLibrariesTypeBufferOverflow:          "缓冲区溢出",
		VulLibrariesTypeSpoofing:                "内容欺骗",
		VulLibrariesTypeCredentialPrediction:    "证书预测",
		VulLibrariesTypeSessionPrediction:       "会话预测",
		VulLibrariesTypeDos:                     "拒绝服务",
		VulLibrariesTypeFinger:                  "指纹识别",
		VulLibrariesTypeFormatString:            "格式化字符串",
		VulLibrariesTypeHttpResponseSmuggling:   "http响应伪造",
		VulLibrariesTypeHttpResponseSplitting:   "http响应拆分",
		VulLibrariesTypeHttpRequestSplitting:    "http请求拆分",
		VulLibrariesTypeHttpRequestSmuggling:    "http请求伪造",
		VulLibrariesTypeHpp:                     "http参数污染",
		VulLibrariesTypeIntOverflow:             "整数溢出",
		VulLibrariesTypeResLocation:             "可预测资源定位",
		VulLibrariesTypeSessionFixation:         "会话固定",
		VulLibrariesTypeRedirect:                "url重定向",
		VulLibrariesTypePrivilegeEscalation:     "权限提升",
		VulLibrariesTypeResolveError:            "解析错误",
		VulLibrariesTypeFileCreating:            "任意文件创建",
		VulLibrariesTypeFileDownload:            "任意文件下载",
		VulLibrariesTypeFileDeletion:            "任意文件删除",
		VulLibrariesTypeFileRead:                "任意文件读取",
		VulLibrariesTypeBakFileFound:            "备份文件发现",
		VulLibrariesTypeDbFound:                 "数据库发现",
		VulLibrariesTypeDirListing:              "目录遍历",
		VulLibrariesTypeDirTraversal:            "目录穿越",
		VulLibrariesTypeFileUpload:              "文件上传",
		VulLibrariesTypeLoginBypass:             "登录绕过",
		VulLibrariesTypeWeakPass:                "弱密码",
		VulLibrariesTypeRemotePassChange:        "远程密码修改",
		VulLibrariesTypeCodeDisclosure:          "代码泄漏",
		VulLibrariesTypePathDisclosure:          "路径泄漏",
		VulLibrariesTypeInfoDisclosure:          "信息泄漏",
		VulLibrariesTypeSecBypass:               "安全模式绕过",
		VulLibrariesTypeMal:                     "挂马",
		VulLibrariesTypeBlackLink:               "暗链",
		VulLibrariesTypeBackdoor:                "后门",
		VulLibrariesTypeInsecureCookieHandling:  "不安全的Cookie",
		VulLibrariesTypeShellCode:               "shellcode",
		VulLibrariesTypeVariableCoverage:        "变量覆盖",
		VulLibrariesTypeInjectingMalwareCodes:   "恶意代码注入",
		VulLibrariesTypeUploadFiles:             "文件上传",
		VulLibrariesTypeLocalOverflow:           "本地溢出",
		VulLibrariesTypePathTraversal:           "目录穿越",
		VulLibrariesTypeUnauthAccess:            "未授权访问",
		VulLibrariesTypeRemoteOverflow:          "远程溢出",
		VulLibrariesTypeMitm:                    "中间人攻击",
		VulLibrariesTypeOutOfMemory:             "内存溢出",
		VulLibrariesTypeBufferOverRead:          "缓冲区越界读",
		VulLibrariesTypeBackupFileFound:         "备份文件泄漏",
		VulLibrariesTypeUaf:                     "释放后使用",
		VulLibrariesTypeDnsHijacking:            "DNS劫持",
		VulLibrariesTypeImproperInputValidation: "不正确的输入校验",
		VulLibrariesTypeUxss:                    "通用型XSS",
		VulLibrariesTypeSsrf:                    "服务器端请求伪造",
		VulLibrariesTypeCommandExec:             "命令执行",
		VulLibrariesTypeCORS:                    "跨域资源共享",
		VulLibrariesTypeOther:                   "其他",
		VulLibrariesTypeLogic:                   "逻辑漏洞",
		VulLibrariesTypeXiuTan:                  "嗅探欺骗",
		VulLibrariesTypeShenTou:                 "渗透测试",
		VulLibrariesTypeXinXi:                   "信息收集",
		VulLibrariesTypePortScan:                "端口扫描",
		VulLibrariesTypeBrute:                   "暴力破解",
		VulLibrariesTypeVulScan:                 "漏洞扫描",
		VulLibrariesTypeVulExp:                  "漏洞利用",
		VulLibrariesTypeJiaJie:                  "加密解密",
		VulLibrariesTypeZhuabao:                 "抓包改包",
		VulLibrariesTypeInject:                  "注入工具",
		VulLibrariesTypeTiQuan:                  "提权",
		VulLibrariesTypeDirScan:                 "目录扫描",
		VulLibrariesTypeMianSha:                 "免杀辅助",
		VulLibrariesTypeNiXiang:                 "逆向",
	}

	result := make([]VulTypeOption, 0, len(typeEnum))
	for k, v := range typeEnum {
		result = append(result, VulTypeOption{
			Value: k,
			Label: v,
		})
	}
	return result
}
