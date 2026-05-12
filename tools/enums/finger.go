package enums

const (
	// v0.6-新统计
	FingerClassMiddleware               = 1  // 中间件
	FingerClassHoneypot                 = 2  // 蜜罐
	FingerClassMainAppFramework         = 3  // 主流应用框架
	FingerClassCollaborativeOfficeSuite = 4  // 协同办公套件
	FingerClassMainCMS                  = 5  // 主流CMS
	FingerClassNetworkEquipment         = 6  // 网络设备
	FingerClassSafeEquipment            = 7  // 安全设备
	FingerClassDatabaseSer              = 8  // 数据库服务
	FingerClassCodeDevelopment          = 9  // 代码研发
	FingerClassBigDataPlatform          = 10 // 大数据平台
	FingerClassVirtualSer               = 11 // 虚拟化服务
	FingerClassMainThirdSer             = 12 // 主流第三方服务
	FingerClassIOT                      = 13 // IOT
	FingerClassOfficeEquipment          = 14 // 办公设备
	FingerClassMailServer               = 15 // 邮件服务器
	FingerClassCentralControlClass      = 16 // 集权管控类
	FingerClassHostOperatingSystem      = 17 // 主机操作系统
	FingerClassFuncType                 = 18 // 功能类型
	FingerClassUniversalVulDetection    = 19 // 通用漏洞检测
	FingerClassCloudPlatform            = 20 // 云计算类型
	FingerClassWebAPP                   = 21 // web应用类型
	// 老的类型统计内容
	FingerClassContentManagementSystem  = 1  // 内容管理系统
	FingerClassHardwareDevice           = 2  // 硬件设备
	FingerClassDatabaseManagementSystem = 3  // 数据库管理系统
	FingerClassDocumentTool             = 4  // 文档工具
	FingerClassWebComponent             = 5  // web组件
	FingerClassECommerceSystem          = 6  // 电子商务系统
	FingerClassImageLibrary             = 7  // 图片库
	FingerClassMotherboard              = 9  // 主机板
	FingerClassAnalysisSystem           = 10 // 分析系统
	FingerClassJSFramework              = 12 // js框架
	FingerClassIssueTracker             = 13 // 问题追踪器
	FingerClassCommentSystem            = 15 // 评论系统
	FingerClassCaptcha                  = 16 // 验证码
	FingerClassFontScript               = 17 // 字体脚本
	FingerClassWebDevelopmentFramework  = 18 // web开发框架
	FingerClassMiscellaneous            = 19 // 杂项
	FingerClassEditor                   = 20 // 编辑器
	FingerClassLearningManagementSystem = 21 // 学习管理系统
	FingerClassWebServer                = 22 // web服务器
	FingerClassCacheTool                = 23 // 缓存工具
	FingerClassJSGraphicsFramework      = 25 // js图形框架
	FingerClassMobileFramework          = 26 // 移动框架
	FingerClassProgrammingLanguage      = 27 // 编程语言
	FingerClassOperatingSystem          = 28 // 操作系统
	FingerClassSearchEngine             = 29 // 搜索引擎
	FingerClassWebmail                  = 30 // 网页邮件
	FingerClassCdn                      = 31 // cdn
	FingerClassMarketingAutomation      = 32 // 营销自动化
	FingerClassWebServerExtension       = 33 // Web服务器扩展
	FingerClassDatabase                 = 34 // 数据库
	FingerClassManagementAnalysis       = 35 // 管理分析与规划系统
	FingerClassAdvertisingNetwork       = 36 // 广告网络
	FingerClassNetworkDevice            = 37 // 网络设备
	FingerClassMediaServer              = 38 // 媒体服务器
	FingerClassVideoSurveillance        = 39 // 视频监控
	FingerClassPaymentProcessor         = 41 // 付款处理器
	FingerClassTagManagementSystem      = 42 // 标签管理系统
	FingerClassBuildCISystem            = 44 // 构建CI系统
	FingerClassControlSystem            = 45 // 控制系统
	FingerClassRemoteAccessSystem       = 46 // 远程访问系统
	FingerClassDevelopmentTool          = 47 // 开发工具
	FingerClassNetworkStorage           = 48 // 网络存储
	FingerClassFeedReader               = 49 // feed阅读器
	FingerClassDocumentManagementSystem = 50 // 文档管理系统
	FingerClassLoginPageGenerator       = 51 // 登录页生成器
	FingerClassOnlineChat               = 52 // 在线聊天
	FingerClassCustomerManagementSystem = 53 // 客户管理系统
	FingerClassSEO                      = 54 // SEO
	FingerClassBillingSystem            = 55 // 账单系统
	FingerClassEncryptionSoftware       = 56 // 加密软件
	FingerClassStaticWebsiteGenerator   = 57 // 静态网站生成器
	FingerClassUserPortable             = 58 // 用户随身携带
	FingerClassJSLibrary                = 59 // js库
	FingerClassContainer                = 60 // 容器
	FingerClassSaaS                     = 61 // SaaS
	FingerClassPaaS                     = 62 //
	FingerClassIaaS                     = 63 //
	FingerClassBigData                  = 64 // 大数据
	FingerClassLoadBalancer             = 65 // 负载均衡
	FingerClassSecurityDevice           = 66 // 安全设备
	FingerClassApplicationService       = 69 // 应用服务
	FingerClassEmail                    = 70 // 电子邮件
	FingerClassOfficeAuto               = 71 // 办公自动化
	FingerClassOther                    = 72 // 其他

	// 服务分层	硬件层、系统层、服务层、支撑层、应用层
	OtherLayer    = 0 // 未知
	HardwareLayer = 1 // 硬件层
	SystemLayer   = 2 // 系统层
	ServiceLayer  = 3 // 服务层
	SupportLayer  = 4 // 支撑层
	AppLayer      = 5 // 应用层
	// 指纹设备类型 软硬件
	FingerSoft = 0 // 软件
	FingerHard = 1 // 硬件

)

const (
	FingerSourceSystem = 1 // 系统自带
	FingerSourceUser   = 2 // 用户添加
)

const (
	FingerTypeWeb     = 1 // web类型
	FingerTypeService = 2 // 服务类型
	FingerTypeDevice  = 3 // 设备类型
)

var FingerEnum finger

type finger struct {
}

// AllNewClass v6.0 新类型统计
func (f *finger) AllNewClass() map[int]string {
	return map[int]string{
		// v0.6-新统计
		FingerClassMiddleware:               "中间件",
		FingerClassHoneypot:                 "蜜罐",
		FingerClassMainAppFramework:         "主流应用框架",
		FingerClassCollaborativeOfficeSuite: "协同办公套件",
		FingerClassMainCMS:                  "主流CMS",
		FingerClassNetworkEquipment:         "网络设备",
		FingerClassSafeEquipment:            "安全设备",
		FingerClassDatabaseSer:              "数据库服务",
		FingerClassCodeDevelopment:          "代码研发",
		FingerClassBigDataPlatform:          "大数据平台",
		FingerClassVirtualSer:               "虚拟化服务",
		FingerClassMainThirdSer:             "主流第三方服务",
		FingerClassIOT:                      "IOT",
		FingerClassOfficeEquipment:          "办公设备",
		FingerClassMailServer:               "邮件服务器",
		FingerClassCentralControlClass:      "集权管控类",
		FingerClassHostOperatingSystem:      "主机操作系统",
		FingerClassFuncType:                 "功能类型",
		FingerClassUniversalVulDetection:    "通用漏洞检测",
		FingerClassOther:                    "其他",
		FingerClassCloudPlatform:            "云计算类型",
		FingerClassWebAPP:                   "web应用类型",
		FingerClassCacheTool:                "缓存工具",
		FingerClassJSGraphicsFramework:      "js图形框架",
		FingerClassMobileFramework:          "移动框架",
		FingerClassProgrammingLanguage:      "编程语言",
		FingerClassOperatingSystem:          "操作系统",
		FingerClassSearchEngine:             "搜索引擎",
		FingerClassWebmail:                  "网页邮件",
		FingerClassCdn:                      "cdn",
		FingerClassMarketingAutomation:      "营销自动化",
		FingerClassManagementAnalysis:       "管理分析与规划系统",
		FingerClassAdvertisingNetwork:       "广告网络",
		FingerClassNetworkDevice:            "网络设备",
		FingerClassMediaServer:              "媒体服务器",
		FingerClassTagManagementSystem:      "标签管理系统",
		FingerClassBuildCISystem:            "构建CI系统",
		FingerClassControlSystem:            "控制系统",
		FingerClassRemoteAccessSystem:       "远程访问系统",
		FingerClassDevelopmentTool:          "开发工具",
		FingerClassNetworkStorage:           "网络存储",
		FingerClassLoadBalancer:             "负载均衡",
		FingerClassApplicationService:       "应用服务",
		FingerClassOfficeAuto:               "办公自动化",
	}
}
func (f *finger) GetNewClass(class int) string {
	if res, ok := f.AllNewClass()[class]; ok {
		return res
	}
	return "暂无"
}

// AllFingerLevel 指纹层级
func (f *finger) AllFingerLevel() map[int]string {
	return map[int]string{
		OtherLayer:    "未分层",
		HardwareLayer: "硬件层",
		SystemLayer:   "系统层",
		ServiceLayer:  "服务层",
		SupportLayer:  "支撑层",
		AppLayer:      "应用层",
	}
}

// AllFingerType 所有指纹类型
func (f *finger) AllFingerType() map[int]string {
	return map[int]string{
		FingerTypeWeb:     "web",
		FingerTypeService: "服务",
		FingerTypeDevice:  "设备",
	}
}

// AllSource 所有来源类型
func (f *finger) AllSource() map[int]string {
	return map[int]string{
		FingerSourceSystem: "系统自带",
		FingerSourceUser:   "用户添加",
	}
}

// AllFingerSoftOrHard 指纹设备软硬件
func (f *finger) AllFingerSoftOrHard() map[int]string {
	return map[int]string{
		FingerSoft: "软件",
		FingerHard: "硬件",
	}
}

// GetFingerSoftOrHard 判断指纹是 软件/硬件
func (f *finger) GetFingerSoftOrHard(typeID int) string {
	if res, ok := f.AllFingerSoftOrHard()[typeID]; ok {
		return res
	}
	return "软件"
}

// GetFingerLevel 获取指纹等级
func (f *finger) GetFingerLevel(level int) string {
	if res, ok := f.AllFingerLevel()[level]; ok {
		return res
	}
	return "暂无"
}

func (f *finger) AllClass() map[int]string {
	return map[int]string{
		FingerClassContentManagementSystem:  "内容管理系统",
		FingerClassHardwareDevice:           "硬件设备",
		FingerClassDatabaseManagementSystem: "数据库管理系统",
		FingerClassDocumentTool:             "文档工具",
		FingerClassWebComponent:             "web组件",
		FingerClassECommerceSystem:          "电子商务系统",
		FingerClassImageLibrary:             "图片库",
		FingerClassMotherboard:              "主机板",
		FingerClassAnalysisSystem:           "分析系统",
		FingerClassJSFramework:              "js框架",
		FingerClassIssueTracker:             "问题追踪器",
		FingerClassCommentSystem:            "评论系统",
		FingerClassCaptcha:                  "验证码",
		FingerClassFontScript:               "字体脚本",
		FingerClassWebDevelopmentFramework:  "web开发框架",
		FingerClassMiscellaneous:            "杂项",
		FingerClassEditor:                   "编辑器",
		FingerClassLearningManagementSystem: "学习管理系统",
		FingerClassWebServer:                "web服务器",
		FingerClassCacheTool:                "缓存工具",
		FingerClassJSGraphicsFramework:      "js图形框架",
		FingerClassMobileFramework:          "移动框架",
		FingerClassProgrammingLanguage:      "编程语言",
		FingerClassOperatingSystem:          "操作系统",
		FingerClassSearchEngine:             "搜索引擎",
		FingerClassWebmail:                  "网页邮件",
		FingerClassCdn:                      "cdn",
		FingerClassMarketingAutomation:      "营销自动化",
		FingerClassWebServerExtension:       "Web服务器扩展",
		FingerClassDatabase:                 "数据库",
		FingerClassManagementAnalysis:       "管理分析与规划系统",
		FingerClassAdvertisingNetwork:       "广告网络",
		FingerClassNetworkDevice:            "网络设备",
		FingerClassMediaServer:              "媒体服务器",
		FingerClassVideoSurveillance:        "视频监控",
		FingerClassPaymentProcessor:         "付款处理器",
		FingerClassTagManagementSystem:      "标签管理系统",
		FingerClassBuildCISystem:            "构建CI系统",
		FingerClassControlSystem:            "控制系统",
		FingerClassRemoteAccessSystem:       "远程访问系统",
		FingerClassDevelopmentTool:          "开发工具",
		FingerClassNetworkStorage:           "网络存储",
		FingerClassFeedReader:               "feed阅读器",
		FingerClassDocumentManagementSystem: "文档管理系统",
		FingerClassLoginPageGenerator:       "登录页生成器",
		FingerClassOnlineChat:               "在线聊天",
		FingerClassCustomerManagementSystem: "客户管理系统",
		FingerClassSEO:                      "SEO",
		FingerClassBillingSystem:            "账单系统",
		FingerClassEncryptionSoftware:       "加密软件",
		FingerClassStaticWebsiteGenerator:   "静态网站生成器",
		FingerClassUserPortable:             "用户随身携带",
		FingerClassJSLibrary:                "js库",
		FingerClassContainer:                "容器",
		FingerClassSaaS:                     "SaaS",
		FingerClassPaaS:                     "PasS",
		FingerClassIaaS:                     "IaaS",
		FingerClassBigData:                  "大数据",
		FingerClassLoadBalancer:             "负载均衡",
		FingerClassSecurityDevice:           "安全设备",
		FingerClassApplicationService:       "应用服务",
		FingerClassEmail:                    "电子邮件",
		FingerClassOfficeAuto:               "办公自动化",
		FingerClassOther:                    "其他",
	}
}
func (f *finger) GetClass(class int) string {
	if res, ok := f.AllClass()[class]; ok {
		return res
	}
	return "暂无"
}

// 获取指纹枚举值，因为导入的文件是字符串

func (f *finger) GetClassEnumInt(class string) int {
	allClass := f.AllClass()
	for k, v := range allClass {
		if v == class {
			return k
		}
	}
	return FingerClassOther
}

// 是否设备指纹
const (
	FingerIsDevN = 1
	FingerIsDevY = 2
)

func (f *finger) AllIsDev() map[int]string {
	return map[int]string{
		FingerIsDevN: "否",
		FingerIsDevY: "是",
	}
}
func (f *finger) GetIsDev(isDev int) string {
	return f.AllIsDev()[isDev]
}
func (f *finger) GetIsDevEnumInt(isDev string) int {
	allIsDev := f.AllIsDev()
	for k, v := range allIsDev {
		if v == isDev {
			return k
		}
	}
	return FingerIsDevN
}

const (
	FingerFromSourceDomestic = iota + 1
	FingerFromSourceForeign
)

var FingerFromSourceMap = map[int]string{
	FingerFromSourceDomestic: "国产",
	FingerFromSourceForeign:  "非国产",
}

func (f *finger) GetFromSource(source int) string {
	return FingerFromSourceMap[source]
}

func (f *finger) FromSourceEnums() interface{} {
	return []struct {
		Label string `json:"label"`
		Value int    `json:"value"`
	}{
		{
			Label: "国产",
			Value: FingerFromSourceDomestic,
		},
		{
			Label: "非国产",
			Value: FingerFromSourceForeign,
		},
	}
}
