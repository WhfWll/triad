package enums

import (
	"encoding/json"
	"strconv"
)

// 配置类型 注意修改此值时同时需要修改 ConfigJson 的json值
const (
	ConfigJsonPortScanKey         = "portScanConfig"         //端口扫描配置
	ConfigJsonWebCrawlerKey       = "webCrawlerConfig"       //动态爬虫配置
	ConfigJsonPathScanKey         = "webPathScanConfig"      //路径爆破配置
	ConfigJsonWeakPass            = "weakPassConfig"         //弱口令配置
	ConfigJsonSubdomainCollectKey = "subdomainCollectConfig" //子域名收集配置
	ConfigJsonWebsiteLoginKey     = "websiteLoginConfig"     //网站登录凭证配置
	ConfigJsonVulIdsKey           = "vulIdsConfig"           // 关联的漏洞ID
	ConfigJsonAliveProbeConfig    = "aliveProbeConfig"       // 关联的漏洞ID
	ConfigSafeTestKey             = "safeTest"               // 安全测试
	ConfigLateralMoveKey          = "lateralMove"            // 横向移动
	ConfigTestIntensityKey        = "testIntensity"          // 测试强度
	ConfigVulExploitKey           = "vulExploit"             // 漏洞利用
	ConfigSceneMapSetKey          = "taskSceneEnumValue"
)

// 解析configJson
func (configJson *ConfigJson) Decode(data map[string]string) (err error) {
	for key, item := range data {
		byteData := []byte(item)

		switch key {
		case ConfigJsonPortScanKey:
			var portScanConfig PortScanConfig
			err = json.Unmarshal(byteData, &portScanConfig)
			configJson.PortScanConfig = portScanConfig
		case ConfigJsonWebCrawlerKey:
			var webCrawlerConfig WebCrawlerConfig
			err = json.Unmarshal(byteData, &webCrawlerConfig)
			configJson.WebCrawlerConfig = webCrawlerConfig
		case ConfigJsonPathScanKey:
			var webPathScanConfig WebPathScanConfig
			err = json.Unmarshal(byteData, &webPathScanConfig)
			configJson.WebPathScanConfig = webPathScanConfig
		case ConfigJsonWeakPass:
			var weakPassConfig WeakPassConfig
			err = json.Unmarshal(byteData, &weakPassConfig)
			configJson.WeakPassConfig = weakPassConfig
		case ConfigJsonSubdomainCollectKey:
			var subdomainCollectConfig SubdomainCollectConfig
			err = json.Unmarshal(byteData, &subdomainCollectConfig)
			configJson.SubdomainCollectConfig = subdomainCollectConfig
		case ConfigJsonWebsiteLoginKey:
			var websiteLoginConfig WebsiteLoginConfig
			err = json.Unmarshal(byteData, &websiteLoginConfig)
			configJson.WebsiteLoginConfig = websiteLoginConfig
		case ConfigJsonVulIdsKey:
			var vulIds []int
			err = json.Unmarshal(byteData, &vulIds)
			configJson.VulIdsConfig = vulIds
		case ConfigJsonAliveProbeConfig:
			var aliveProbeConfig AliveProbeConfig
			err = json.Unmarshal(byteData, &aliveProbeConfig)
			configJson.AliveProbeConfig = aliveProbeConfig
		case ConfigSafeTestKey: //安全测试
			safeTest, err := strconv.ParseBool(string(byteData))
			if err != nil {
				return err
			}
			configJson.SafeTest = safeTest
		case ConfigLateralMoveKey: //横向移动
			var lateralMove LateralMove
			err = json.Unmarshal(byteData, &lateralMove)
			configJson.LateralMove = lateralMove
		case ConfigTestIntensityKey: //测试强度
			testIntensity, err := strconv.Atoi(string(byteData))
			if err != nil {
				return err
			}
			configJson.TestIntensity = testIntensity
		case ConfigVulExploitKey: //漏洞利用
			vulExploit, err := strconv.ParseBool(string(byteData))
			if err != nil {
				return err
			}
			configJson.VulExploit = vulExploit
		}
		if err != nil {
			return err
		}
	}

	return
}

/************** ConfigJson 数据定义 ******************/
/*
ConfigJson
* 定义 taskConfiguration数据表 configJson 结构体
* 谨慎修改，多处使用，迁移到公共仓里面
*/
type ConfigJson struct {
	AliveProbeConfig       AliveProbeConfig       `json:"aliveProbeConfig"`       // 端口扫描配置
	PortScanConfig         PortScanConfig         `json:"portScanConfig"`         // 端口扫描配置
	WebCrawlerConfig       WebCrawlerConfig       `json:"webCrawlerConfig"`       // 动态爬虫配置
	WebPathScanConfig      WebPathScanConfig      `json:"webPathScanConfig"`      // 路径爆破配置
	WeakPassConfig         WeakPassConfig         `json:"weakPassConfig"`         // 弱口令配置
	SubdomainCollectConfig SubdomainCollectConfig `json:"subdomainCollectConfig"` // 子域名收集配置
	WebsiteLoginConfig     WebsiteLoginConfig     `json:"websiteLoginConfig"`     // 网站登录凭证配置
	VulIdsConfig           []int                  `json:"vulIdsConfig"`           // 关联的漏洞ID
	VulRiskLevels          []int                  `json:"vulRiskLevels"`          // 未指定插件时按风险等级筛选（应用安全）
	VulClassLevels         []int                  `json:"vulClassLevels"`         // 未指定插件时按漏洞分类筛选（应用安全）
	TestIntensity          int                    `json:"testIntensity"`          // 测试强度
	VulExploit             bool                   `json:"vulExploit"`             // 漏洞利用
	SafeTest               bool                   `json:"safeTest"`               // 安全测试
	LateralMove            LateralMove            `json:"lateralMove"`            // 横向移动
	Mode                   Mode                   `json:"mode"`                   // 渗透模式
	ProxyConfig            ProxyConfig            `json:"proxyConfig"`            // 代理模式
	TestMode               string                 `json:"testMode"`               // 测试模式
}

/************** mode 数据定义 ******************/
/** 渗透模式
未开启分布式：yak引擎为CMD方式调用
	0通用渗透：使用本机节点
	1定向渗透：不可使用
开启分布式：yak引擎为GRPC方式调用
	0通用渗透：轮询使用可用节点，如节点都不存活，则无法创建任务
	1定向渗透：仅支持分布式：如未开启分布式引擎，则无法选择节点，开启后，可选在线且启用的节点
*/

type DistributeNodeIds []int

func (d *DistributeNodeIds) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	// 如果是数组
	if data[0] == '[' {
		var ids []int
		if err := json.Unmarshal(data, &ids); err != nil {
			return err
		}
		*d = ids
		return nil
	}
	// 如果是单个数字
	var id int
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	*d = []int{id}
	return nil
}

type Mode struct {
	Mode             int               `form:"mode" json:"mode"`                         // 渗透模式
	ModeZh           string            `form:"modeZh" json:"modeZh"`                     // 渗透模式
	DistributeNodeId DistributeNodeIds `form:"distributeNodeId" json:"distributeNodeId"` // 定向渗透节点ID
}

// 渗透模式
const (
	TaskConfigurationModeCommon = 1 // 通用渗透
	TaskConfigurationModeTarget = 2 // 定向渗透
)

func (mode *Mode) AllModeEnum() map[int]string {
	enum := map[int]string{
		TaskConfigurationModeCommon: "通用渗透",
		TaskConfigurationModeTarget: "定向渗透",
	}
	return enum
}

func (mode *Mode) GetModeEnum(k int) string {
	enum := mode.AllModeEnum()

	value, ok := enum[k]
	if ok {
		return value
	}

	return ""
}

/************** PortScanConfig 数据定义 ******************/

// AliveProbeConfig 探活配置
type AliveProbeConfig struct {
	IsOpen           bool   `json:"isOpen"`           //是否开启 1:开启, 2:不开启
	AliveProbeType   int    `json:"aliveProbeType"`   //存活探测类型value
	AliveProbeTypeZh string `json:"aliveProbeTypeZh"` //端口扫描类型label
	ProbePort        string `json:"probePort"`        //扫描的端口
}

// PortScanConfig 端口扫描配置
type PortScanConfig struct {
	IsOpen           bool   `json:"isOpen"`           //是否开启 1:开启, 2:不开启
	IntelligencePort bool   `json:"intelligencePort"` //是否开启智能端口扫描
	TCPScanType      int    `json:"tcpScanType"`      //tcp扫描类型value
	TCPScanTypeZh    string `json:"tcpScanTypeZh"`    //tcp扫描类型label
	PortScanType     int    `json:"portScanType"`     //端口扫描类型value
	PortScanTypeZh   string `json:"portScanTypeZh"`   //端口扫描类型label
	Timeout          int    `json:"timeout"`          //端口扫描超时时间(秒)
	TimeoutZh        string `json:"timeoutZh"`        //端口扫描超时时间
	Concurrent       int    `json:"concurrent"`       //端口扫描并发数
	ConcurrentZh     string `json:"concurrentZh"`     //端口扫描并发数
	ScanPort         string `json:"scanPort"`         //扫描的端口
}

/************** WebCrawlerConfig 数据定义 ******************/
// WebCrawlerConfig 动态爬虫配置
type WebCrawlerConfig struct {
	IsOpen     bool   `json:"isOpen"`     //是否开启动态爬虫配置
	MaxDepth   int    `json:"maxDepth"`   // 爬取深度 0不限
	MaxDepthZh string `json:"maxDepthZh"` // 爬取深度
	MaxUrl     int    `json:"maxUrl"`     // 最大爬取url数量 0不限
	MaxUrlZh   string `json:"maxUrlZh"`   // 最大爬取url数量
	/**
	`newcrawlerx.AllDomainScan` 表示爬取全域名 （默认0）
	`newcrawlerx.SubMenuScan` 表示爬取目标URL和子目录
	*/
	ScanRange       int    `json:"scanRange"`     // 爬虫爬取范围
	ScanRangeZh     string `json:"scanRangeZh"`   // 爬虫爬取范围
	Timeout         int    `json:"timeout"`       // 爬虫单页面超时时间设置 单位秒
	TimeoutZh       string `json:"timeoutZh"`     // 爬虫单页面超时时间设置
	FullTimeout     int    `json:"fullTimeout"`   // 爬虫全局超时时间设置 单位秒
	FullTimeoutZh   string `json:"fullTimeoutZh"` // 爬虫全局超时时间设置
	CrawlerSpeedKey int    `json:"crawlerSpeed"`  // 爬虫爬取速度
	/** 爬虫结果重复过滤设置
	0 不限，请勿设置爬虫过滤
	1 `newcrawlerx.UnLimitRepeat` 对page，method，query-name，query-value和post-data敏感
	2 `newcrawlerx.LowRepeatLevel` 对page，method，query-name和query-value敏感（默认）
	3 `newcrawlerx.MediumRepeatLevel` 对page，method和query-name敏感
	4 `newcrawlerx.HighRepeatLevel` 对page和method敏感
	5 `newcrawlerx.ExtremeRepeatLevel` 对page敏感
	*/
	ScanRepeat   int                       `json:"scanRepeat"` // 爬虫结果重复过滤设置
	ScanRepeatZh string                    `json:"scanRepeatZh"`
	BlackList    string                    `json:"blackList"`    // 爬虫黑名单参数设置 多个英文逗号分割
	WhiteList    string                    `json:"whiteList"`    // 爬虫白名单参数设置 多个英文逗号分割
	Headers      []WebCrawlerHeadersConfig `json:"headers"`      // 爬虫request的header设置
	SuffixFilter string                    `json:"suffixFilter"` //后缀过滤
	LocalStorage LocalStorage              `json:"localStorage"` //localStorage
}
type WebCrawlerHeadersConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LocalStorage struct {
	IsOpen bool               `json:"isOpen"`
	List   []LocalStoragelist `json:"list"`
}

type LocalStoragelist struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

/************** WebPathScanConfig 数据定义 ******************/
// WebPathScanConfig 路径爆破配置
type WebPathScanConfig struct {
	IsOpen         bool     `json:"isOpen"`         //是否开启路径爆破配置
	GuessRate      int      `json:"guessRate"`      //猜测速率
	GuessRateZh    string   `json:"guessRateZh"`    //猜测速率
	GuessTimeout   int      `json:"guessTimeout"`   //猜测时长
	GuessTimeoutZh string   `json:"guessTimeoutZh"` //猜测时长
	TitleBlack     string   `json:"titleBlack"`     //排除标题黑名单
	ScanDict       []int    `json:"scanDict"`       //路径字典列表 注意：这里记录的是字典库ID
	DickNames      []string `json:"dickNames"`      //路径字典名称列表
	IsIntelligent  bool     `json:"isIntelligent"`  // 是否打开智能路径爆破
}

/************** WeakPassConfig 数据定义 ******************/
// WeakPassConfig 弱口令配置
type WeakPassConfig struct {
	IsOpen     bool     `json:"isOpen"`     // 是否开启弱口令配置
	Services   []int    `json:"services"`   // 服务 注意：这里记录的是服务ID，也就是service字段的值，由于types已经默认
	ServicesZh []string `json:"servicesZh"` // 服务名称列表
	/** 字典类型 1默认字典 2通用字典 3补充字典
	=1 需要查询 types=弱口令 and is_default=1 group by service 的content
	=2 需要查询 【用户字典：types=弱口令 and id=CommonUserDict】 【密码字典：types=弱口令 and id=CommonPassDict】
	=3 如果OnlyUseAdd=true 则仅使用AddAccount & AddPass 否则将AddAccount & AddPass分别追加到=1默认字典的数据中
	*/
	DictType         int    `json:"dictType"`
	DictTypeZh       string `json:"dictTypeZh"`
	CommonUserDict   int    `json:"commonUserDict"`   // 通用用户字典 注意：这里记录的是字典库的主键ID
	CommonUserDictZh string `json:"commonUserDictZh"` // 通用用户字典
	CommonPassDict   int    `json:"commonPassDict"`   // 通用密码字典 注意：这里记录的是字典库的主键ID
	CommonPassDictZh string `json:"commonPassDictZh"` // 通用密码字典
	AddAccount       string `json:"addAccount"`       // 补充账号
	AddPass          string `json:"addPass"`          // 补充密码
	OnlyUseAdd       bool   `json:"onlyUseAdd"`       // 仅使用补充字典
	GuessNum         int    `json:"guessNum"`         // 猜测次数
	GuessNumZh       string `json:"guessNumZh"`       // 猜测次数
	GuessTimeout     int    `json:"guessTimeout"`     // 猜测时间
	GuessTimeoutZh   string `json:"guessTimeoutZh"`   // 猜测时间
	GuessRate        int    `json:"guessRate"`        // 猜测速率
	GuessRateZh      string `json:"guessRateZh"`      // 猜测速率
	CaptchaMode      string `json:"captchaMode"`      // 带验证码爆破 字母数字common_alphanumeric 简单算数common_arithmetic
}

/************** SubdomainCollectConfig 数据定义 ******************/
// SubdomainCollectConfig 子域名收集配置
type SubdomainCollectConfig struct {
	IsOpen          bool   `json:"isOpen"`          //是否开启子域名收集配置
	SubdomainDict   int    `json:"subdomainDict"`   //子域名字典  注意：这里记录的是字典库ID
	SubdomainDictZh string `json:"subdomainDictZh"` //子域名字典
}

/************* WebsiteLoginConfig 站点登录凭证配置 ****************/
// WebsiteLoginConfig 站点登录凭证配置
type WebsiteLoginConfig struct {
	IsOpen bool              `json:"isOpen"` //是否开启站点登录凭证配置
	List   []LoginConfigNode `json:"list"`   //站点登录信息列表
}

// LoginConfigNode 站点登录信息
type LoginConfigNode struct {
	Target         string `json:"target"`         //登录地址
	VerifyType     int    `json:"verifyType"`     //凭证类型
	VerifyTypeZh   string `json:"verifyTypeZh"`   //凭证类型
	VerifyValue    string `json:"verifyValue"`    //凭证
	VerifyStatus   int    `json:"verifyStatus"`   //状态（1/2/3）
	VerifyStatusZh string `json:"verifyStatusZh"` //状态（认证成功/访问成功/访问失败）
	Scheme         string `json:"scheme"`         //协议
}

type LateralMove struct {
	Range    string `json:"range"`
	IsOpen   bool   `json:"isOpen"`
	Strategy string `json:"strategy"` // 策略: auto_subnet | custom_range
	Ports    string `json:"ports"`    // 扫描端口
	Timeout  int    `json:"timeout"`  // 超时时间(秒)
}

/************** ProxyConfig 代理定义 ******************/
// ProxyConfig 代理定义
type ProxyConfig struct {
	IsOpen   bool   `json:"isOpen"`   // 是否开启 1:开启, 2:不开启
	Addr     string `json:"addr"`     // 代理地址
	Port     string `json:"port"`     // 代理端口
	Proto    int    `json:"proto"`    // 代理协议 配置在任务枚举文件中，属于任务的子关系
	IsAuth   bool   `json:"isAuth"`   // 是否开启认证
	Username string `json:"username"` // 认证账号
	Password string `json:"password"` // 认证密码
}

// TestModeConfig 测试模式
type TestModeConfig struct{}

type PenetrationModeConfig struct{}

/**********************************************/
/************* 以下为场景枚举配置 ****************/
/**********************************************/
/************* 端口扫描 **************/
// TCP扫描方式
const (
	TaskConfigurationTcpScanTypeConnect = 1
	TaskConfigurationTcpScanTypeSyn     = 2
	TaskConfigurationTcpScanTypeFin     = 3
	TaskConfigurationTcpScanTypeAck     = 4
	TaskConfigurationTcpScanTypeNull    = 5
	TaskConfigurationUDP                = 6
	TaskConfigurationTcpNull            = 7
	TaskConfigurationTcpScanTypeDefault = 1 // 默认选择 - 下标
)

const (
	TaskConfigurationAliveProbeICMP   = 1
	TaskConfigurationAliveProbeArp    = 2
	TaskConfigurationAliveProbeTCP    = 3
	TaskConfigurationAliveProbeUDP    = 4
	TaskConfigurationAliveProbeTCPACK = 5
	TaskConfigurationAliveProbeTCPSYN = 6
)

func (aliveProbeConfig *AliveProbeConfig) AllAliveProbeTypeEnum() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationAliveProbeICMP:   "ICMP-PING",
		TaskConfigurationAliveProbeArp:    "ARP-PING",
		TaskConfigurationAliveProbeTCP:    "TCP-PING",
		TaskConfigurationAliveProbeUDP:    "UDP-PING",
		TaskConfigurationAliveProbeTCPACK: "TCP-ACK",
		TaskConfigurationAliveProbeTCPSYN: "TCP-SYN",
	}
	return enum, TaskConfigurationAliveProbeICMP
}
func (aliveProbeConfig *AliveProbeConfig) AliveProbeTypeEnum(aliveProbeType int) string {
	allEnums, _ := aliveProbeConfig.AllAliveProbeTypeEnum()
	if res, ok := allEnums[aliveProbeType]; ok {
		return res
	}
	return ""
}

func (aliveProbeConfig *AliveProbeConfig) AllAliveProbePortRangeEnum() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationAliveProbeArp:    "",
		TaskConfigurationAliveProbeICMP:   "",
		TaskConfigurationAliveProbeTCP:    "22,23,80,443,8000,8080",
		TaskConfigurationAliveProbeUDP:    "53,123,161,514",
		TaskConfigurationAliveProbeTCPACK: "22,23,80,443,8000,8080",
		TaskConfigurationAliveProbeTCPSYN: "22,23,80,443,8000,8080",
	}
	return enum, TaskConfigurationAliveProbeArp
}

func (PortScanConfig *PortScanConfig) AllTcpScanTypeEnum() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationTcpScanTypeConnect: "TCP-Connect",
		TaskConfigurationTcpScanTypeSyn:     "TCP SYN",
		TaskConfigurationTcpScanTypeFin:     "TCP FIN",
		TaskConfigurationTcpScanTypeAck:     "TCP ACK",
		TaskConfigurationTcpScanTypeNull:    "TCP NULL",
		TaskConfigurationUDP:                "UDP",
	}
	return enum, TaskConfigurationTcpScanTypeDefault
}
func (PortScanConfig *PortScanConfig) TcpScanTypeEnum(tcpScanType int) string {
	allEnums, _ := PortScanConfig.AllTcpScanTypeEnum()
	if res, ok := allEnums[tcpScanType]; ok {
		return res
	}
	return ""
}

// 端口范围
const (
	TaskConfigurationPortScanTypeTop10   = 10    //TOP10端口
	TaskConfigurationPortScanTypeTop20   = 20    //TOP20端口
	TaskConfigurationPortScanTypeTop50   = 50    //TOP50端口
	TaskConfigurationPortScanTypeTop100  = 100   //TOP100端口
	TaskConfigurationPortScanTypeTop500  = 500   //TOP500端口
	TaskConfigurationPortScanTypeTop1000 = 1000  //TOP1000端口
	TaskConfigurationPortScanTypeAll     = 65535 //全部端口
	TaskConfigurationPortScanTypeCustom  = 0     //自定义端口

	TaskConfigurationPortScanTypeDefault = 100 // 默认选择 - 下标
)

// 端口范围选项
func (PortScanConfig *PortScanConfig) AllPortScanTypeEnum() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationPortScanTypeTop10:   "TOP10端口",
		TaskConfigurationPortScanTypeTop20:   "TOP20端口",
		TaskConfigurationPortScanTypeTop50:   "TOP50端口",
		TaskConfigurationPortScanTypeTop100:  "TOP100端口",
		TaskConfigurationPortScanTypeTop500:  "TOP500端口",
		TaskConfigurationPortScanTypeTop1000: "TOP1000端口",
		TaskConfigurationPortScanTypeAll:     "全部端口",
		TaskConfigurationPortScanTypeCustom:  "自定义端口",
	}
	return enum, TaskConfigurationPortScanTypeDefault
}
func (PortScanConfig *PortScanConfig) PortScanTypeEnum(portScan int) string {
	enum, _ := PortScanConfig.AllPortScanTypeEnum()
	if res, ok := enum[portScan]; ok {
		return res
	}
	return ""
}

// 端口范围选项后联动的端口数据
func (PortScanConfig *PortScanConfig) AllPortScanTypeValue() map[int]string {
	enum := map[int]string{
		TaskConfigurationPortScanTypeTop10:   "21,22,23,80,443,445,3306,8000,8080,8088",
		TaskConfigurationPortScanTypeTop20:   "21,22,23,80,443,445,3306,7000-7002,8000-8003,8080-8083,8088,9200",
		TaskConfigurationPortScanTypeTop50:   "21,22,23,80,88,106,110,111,113,119,135,139,143,144,179,199,389,427,1521,1630,1158,443,445,888,777,999,1070,1080,1090,3306,7000-7003,8000-8003,8008,8080-8083,8088,9000-9002,8090,9200,9300",
		TaskConfigurationPortScanTypeTop100:  "7,9,13,21,22,23,25,26,37,53,79,80,81,88,106,110,111,113,119,135,139,143,144,179,199,389,427,443,444,445,465,513,514,515,543,544,548,554,587,631,646,873,888,990,993,995,1025,1026,1027,1028,1029,1080,1110,1433,1443,1521,1720,1723,1755,1900,2000,2001,2049,2121,2181,2717,3000,3128,3306,3389,3986,4899,5000,5009,5051,5060,5101,5190,5357,5432,5555,5631,5666,5800,5900,6000,6001,6646,7000,7001,7002,7003,7004,7005,7070,8000,8008,8009,8080,8081,8443,8888,9100,9999,10000,11211,32768,49152,49153,49154,49155,49156,49157,8088,9090,8090,8001,82,9080,8082,8089,9000,8002,89,8083,8200,90,8086,801,8011,8085,9001,9200,8100,8012,85,8084,8070,8091,8003,99,7777,8010,8028,8087,83,808,38888,8181,800,18080,8099,8899,86,8360,8300,8800,8180,3505,9002,8053,1000,7080,8989,28017,9060,8006,41516,880,8484,6677,8016,84,7200,9085,5555,8280,1980,8161,9091,7890,8060,6080,8880,8020,889,8881,9081,7007,8004,38501,1010,17,19,255,1024,1030,1041,1048,1049,1053,1054,1056,1064,1065,1801,2103,2107,2967,3001,3703,5001,5050,6004,8031,10010,10250,10255,6888,87,91,92,98,1081,1082,1118,1888,2008,2020,2100,2375,3008,6648,6868,7008,7071,7074,7078,7088,7680,7687,7688,8018,8030,8038,8042,8044,8046,8048,8069,8092,8093,8094,8095,8096,8097,8098,8101,8108,8118,8172,8222,8244,8258,8288,8448,8834,8838,8848,8858,8868,8879,8983,9008,9010,9043,9082,9083,9084,9086,9087,9088,9089,9092,9093,9094,9095,9096,9097,9098,9099,9443,9448,9800,9981,9986,9988,9998,10001,10002,10004,10008,12018,12443,14000,16080,18000,18001,18002,18004,18008,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018",
		TaskConfigurationPortScanTypeTop500:  "7,8,9,13,17,19,20,21,22,23,25,26,37,53,60,65,66,70,77,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,103,106,110,111,113,114,119,122,132,133,135,139,143,144,171,179,180,188,199,200,206,208,211,235,255,268,280,299,302,308,321,381,389,403,421,423,427,442,443,444,445,447,465,511,513,514,515,517,522,543,544,548,554,587,591,610,631,646,666,688,701,770,778,800,801,802,803,804,805,806,808,809,811,812,866,873,877,880,888,889,925,955,983,990,993,995,999,1000,1001,1005,1010,1024,1025,1026,1027,1028,1029,1030,1039,1041,1042,1048,1049,1053,1054,1056,1064,1065,1080,1081,1082,1085,1088,1100,1107,1108,1110,1118,1122,1123,1128,1158,1180,1182,1200,1212,1213,1234,1300,1301,1313,1356,1433,1443,1500,1521,1550,1666,1680,1700,1720,1722,1723,1755,1790,1800,1801,1818,1863,1888,1900,1933,1949,1979,1980,1982,2000,2001,2005,2006,2007,2008,2009,2010,2011,2012,2013,2014,2015,2020,2046,2049,2051,2060,2070,2080,2093,2100,2103,2107,2110,2121,2125,2181,2222,2301,2348,2375,2382,2480,2490,2517,2521,2585,2717,2808,2886,2901,2967,3000,3001,3008,3012,3013,3030,3050,3080,3128,3216,3220,3306,3312,3333,3380,3389,3443,3456,3465,3503,3505,3535,3580,3588,3600,3606,3668,3690,3703,3721,3880,3938,3986,4000,4001,4016,4040,4300,4389,4430,4433,4440,5000,5001,5002,5003,5009,5013,5050,5600,5601,5631,5632,5644,5655,5656,5666,5678,5800,5811,5881,5887,5888,5898,5900,5902,5966,6000,6001,6002,6003,6004,6006,6010,6060,6080,6088,6090,6101,6118,6170,6180,6198,6226,6259,6379,6388,6886,6888,6889,6890,6920,6969,6988,7000,7001,7002,7003,7004,7005,7006,7007,7008,7009,7010,7011,7012,7017,7018,7020,7021,7022,7028,7031,7041,7048,7050,7060,7070,7071,7074,7078,7080,7081,7084,7086,7088,7094,7100,7101,7102,7108,7111,7117,7123,7129,7171,7180,7200,7201,7202,7215,7272,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8016,8018,8019,8020,8021,8022,8023,8024,8025,8026,8027,8028,8029,8030,8031,8032,8033,8035,8036,8037,8038,8039,8040,8041,8042,8043,8044,8045,8046,8048,8050,8051,8053,8055,8056,8057,8058,8060,8061,8062,8064,8065,8066,8069,8070,8071,8073,8077,8078,8079,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8102,8103,8104,8172,8180,8181,8182,8183,8184,8186,8188,8288,8300,8308,8322,8333,8341,8343,8360,8380,8580,8582,8585,8600,8601,8610,8649,8660,9200,9201",
		TaskConfigurationPortScanTypeTop1000: "7,8,9,13,17,19,20,21,22,23,25,26,37,53,60,65,66,70,77,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,103,106,110,111,113,114,119,122,132,133,135,139,143,144,171,179,180,188,199,200,206,208,211,235,255,268,280,299,302,308,321,381,389,403,421,423,427,442,443,444,445,447,465,511,513,514,515,517,522,543,544,548,554,587,591,610,631,646,666,688,701,770,778,800,801,802,803,804,805,806,808,809,811,812,866,873,877,880,888,889,925,955,983,990,993,995,999,1000,1001,1005,1010,1024,1025,1026,1027,1028,1029,1030,1039,1041,1042,1048,1049,1053,1054,1056,1064,1065,1080,1081,1082,1085,1088,1100,1107,1108,1110,1118,1122,1123,1128,1158,1180,1182,1200,1212,1213,1234,1300,1301,1313,1356,1433,1443,1500,1521,1550,1666,1680,1700,1720,1722,1723,1755,1790,1800,1801,1818,1863,1888,1900,1933,1949,1979,1980,1982,2000,2001,2005,2006,2007,2008,2009,2010,2011,2012,2013,2014,2015,2020,2046,2049,2051,2060,2070,2080,2093,2100,2103,2107,2110,2121,2125,2181,2222,2301,2348,2375,2382,2480,2490,2517,2521,2585,2717,2808,2886,2901,2967,3000,3001,3008,3012,3013,3030,3050,3080,3128,3216,3220,3306,3312,3333,3380,3389,3443,3456,3465,3503,3505,3535,3580,3588,3600,3606,3668,3690,3703,3721,3880,3938,3986,4000,4001,4016,4040,4300,4389,4430,4433,4440,4443,4567,4848,4850,4899,5000,5001,5002,5003,5009,5013,5050,5051,5060,5080,5081,5098,5100,5101,5155,5156,5190,5200,5201,5203,5233,5255,5256,5280,5357,5432,5544,5552,5555,5561,5600,5601,5631,5632,5644,5655,5656,5666,5678,5800,5811,5881,5887,5888,5898,5900,5902,5966,6000,6001,6002,6003,6004,6006,6010,6060,6080,6088,6090,6101,6118,6170,6180,6198,6226,6259,6379,6388,6443,6510,6543,6546,6565,6587,6600,6602,6603,6611,6646,6648,6666,6677,6680,6688,6699,6778,6789,6800,6801,6842,6868,6869,6879,6886,6888,6889,6890,6920,6969,6988,7000,7001,7002,7003,7004,7005,7006,7007,7008,7009,7010,7011,7012,7017,7018,7020,7021,7022,7028,7031,7041,7048,7050,7060,7070,7071,7074,7078,7080,7081,7084,7086,7088,7094,7100,7101,7102,7108,7111,7117,7123,7129,7171,7180,7200,7201,7202,7215,7272,7288,7321,7330,7380,7443,7500,7567,7680,7687,7688,7700,7702,7703,7709,7711,7713,7742,7776,7777,7778,7788,7791,7801,7856,7888,7890,7899,7900,7909,7915,7921,7925,7942,7943,7979,7999,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8016,8018,8019,8020,8021,8022,8023,8024,8025,8026,8027,8028,8029,8030,8031,8032,8033,8035,8036,8037,8038,8039,8040,8041,8042,8043,8044,8045,8046,8048,8050,8051,8053,8055,8056,8057,8058,8060,8061,8062,8064,8065,8066,8069,8070,8071,8073,8077,8078,8079,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8102,8103,8104,8108,8111,8112,8118,8119,8122,8123,8130,8133,8136,8144,8161,8168,8172,8180,8181,8182,8183,8184,8186,8188,8189,8190,8191,8192,8193,8196,8197,8200,8213,8220,8222,8244,8258,8260,8280,8282,8283,8288,8300,8308,8322,8333,8341,8343,8360,8380,8381,8382,8383,8384,8390,8399,8400,8401,8402,8443,8445,8448,8477,8480,8481,8484,8500,8567,8580,8582,8585,8600,8601,8610,8649,8660,8666,8680,8686,8688,8700,8710,8720,8735,8780,8781,8787,8788,8799,8800,8801,8802,8806,8808,8809,8810,8813,8822,8834,8838,8839,8844,8848,8858,8860,8864,8866,8868,8877,8879,8880,8881,8885,8886,8887,8888,8889,8890,8891,8892,8895,8898,8899,8900,8902,8910,8912,8913,8955,8956,8972,8980,8983,8987,8988,8989,8990,8991,8997,8999,9000,9001,9002,9003,9004,9005,9006,9007,9008,9009,9010,9011,9012,9013,9014,9015,9020,9022,9025,9030,9031,9036,9039,9043,9050,9053,9060,9061,9070,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9101,9110,9111,9112,9131,9180,9182,9190,9191,9200,9201,9212,9231,9300,9301,9302,9437,9443,9448,9494,9500,9504,9507,9527,9595,9666,9696,9704,9800,9845,9876,9888,9889,9898,9900,9901,9909,9910,9912,9914,9918,9919,9980,9981,9986,9988,9990,9991,9992,9995,9997,9998,9999,10000,10001,10002,10003,10004,10007,10008,10009,10010,10016,10021,10025,10038,10040,10051,10066,10068,10080,10082,10086,10087,10088,10089,10099,10118,10154,10250,10255,10333,11000,11001,11080,11158,11211,11324,11347,11362,11366,11372,11381,12001,12018,12333,12345,12443,12881,13333,13382,13988,14000,14007,15000,15004,15018,15580,15672,15693,15801,15888,16080,16788,17000,17003,17095,17777,18000,18001,18002,18004,18008,18060,18080,18081,18082,18085,18088,18090,18098,18103,18264,18801,18803,18880,18881,18888,19001,19010,19045,19080,19101,19244,20000,20001,20021,20022,20046,20052,20140,20151,20153,20720,20806,20808,21000,21080,21245,21501,21502,22222,22343,22580,23352,23454,25006,25024,27000,27017,28017,28018,28080,28099,28214,28280,28780,30000,30001,30058,30082,30088,30551,31188,31945,32766,32768,34440,38000,38080,38086,38443,38501,38517,38888,40000,40069,40080,40310,41516,42424,43651,45149,45177,47078,47088,47583,48080,49152,49153,49154,49155,49156,49157,49705,49960,50000,50030,50045,50060,50070,50075,50080,50090,50240,51106,55351,55858,57880,58000,58031,58060,58080,58898,59009,59777,59999,60010,60022,60101,60465,61081,61999,65000,65001,65055,65129,65486,65493,65533,65535",
		TaskConfigurationPortScanTypeAll:     "0-65535",
		TaskConfigurationPortScanTypeCustom:  "",
	}
	return enum
}

const (
	TaskConfigurationPortScanTimeout3       = 3
	TaskConfigurationPortScanTimeout5       = 5
	TaskConfigurationPortScanTimeout10      = 10
	TaskConfigurationPortScanTimeout20      = 20
	TaskConfigurationPortScanTimeout30      = 30
	TaskConfigurationPortScanTimeout60      = 60
	TaskConfigurationPortScanTimeout120     = 120
	TaskConfigurationPortScanTimeoutDefault = 10
)

func (PortScanConfig *PortScanConfig) AllPortScanTimeoutEnum() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationPortScanTimeout3:   "3s",
		TaskConfigurationPortScanTimeout5:   "5s",
		TaskConfigurationPortScanTimeout10:  "10s",
		TaskConfigurationPortScanTimeout20:  "20s",
		TaskConfigurationPortScanTimeout30:  "30s",
		TaskConfigurationPortScanTimeout60:  "60s",
		TaskConfigurationPortScanTimeout120: "120s",
	}
	return enum, TaskConfigurationPortScanTimeoutDefault
}

func (PortScanConfig *PortScanConfig) PortScanTimeoutEnum(timeout int) string {
	enum, _ := PortScanConfig.AllPortScanTimeoutEnum()
	if res, ok := enum[timeout]; ok {
		return res
	}
	if timeout > 0 {
		return strconv.Itoa(timeout) + "s"
	}
	return ""
}

const (
	TaskConfigurationPortScanConcurrent1       = 1
	TaskConfigurationPortScanConcurrent5       = 5
	TaskConfigurationPortScanConcurrent10      = 10
	TaskConfigurationPortScanConcurrent20      = 20
	TaskConfigurationPortScanConcurrent50      = 50
	TaskConfigurationPortScanConcurrent100     = 100
	TaskConfigurationPortScanConcurrentDefault = 10
)

func (PortScanConfig *PortScanConfig) AllPortScanConcurrentEnum() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationPortScanConcurrent1:   "1",
		TaskConfigurationPortScanConcurrent5:   "5",
		TaskConfigurationPortScanConcurrent10:  "10",
		TaskConfigurationPortScanConcurrent20:  "20",
		TaskConfigurationPortScanConcurrent50:  "50",
		TaskConfigurationPortScanConcurrent100: "100",
	}
	return enum, TaskConfigurationPortScanConcurrentDefault
}

func (PortScanConfig *PortScanConfig) PortScanConcurrentEnum(concurrent int) string {
	enum, _ := PortScanConfig.AllPortScanConcurrentEnum()
	if res, ok := enum[concurrent]; ok {
		return res
	}
	if concurrent > 0 {
		return strconv.Itoa(concurrent)
	}
	return ""
}

/************* 爬虫配置 **************/
// 爬取范围
// 	0 `newcrawlerx.AllDomainScan` 表示爬取全域名 （默认0）
//	1 `newcrawlerx.SubMenuScan` 表示爬取目标URL和子目录
const (
	TaskConfigurationCrawlerScanRangeAllHost   = 0
	TaskConfigurationCrawlerScanRangeUrlAndDir = 1

	TaskConfigurationCrawlerScanRangeDefault = 0 // 默认选择 - 下标
)

func (WebCrawlerConfig *WebCrawlerConfig) AllCrawlerScanRange() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationCrawlerScanRangeAllHost:   "爬取全域名",
		TaskConfigurationCrawlerScanRangeUrlAndDir: "爬取目标URL和子目录",
	}
	return enum, TaskConfigurationCrawlerScanRangeDefault
}

func (WebCrawlerConfig *WebCrawlerConfig) CrawlerRange(crawler int) string {
	enum, _ := WebCrawlerConfig.AllCrawlerScanRange()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

// 爬取深度
const (
	TaskConfigurationCrawlerDeep0 = 0
	TaskConfigurationCrawlerDeep1 = 1
	TaskConfigurationCrawlerDeep2 = 2
	TaskConfigurationCrawlerDeep3 = 3
	TaskConfigurationCrawlerDeep4 = 4
	TaskConfigurationCrawlerDeep5 = 5

	TaskConfigurationCrawlerDefault = 3 // 默认选择 - 下标
)

func (WebCrawlerConfig *WebCrawlerConfig) AllCrawlerDeep() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationCrawlerDeep1: "1",
		TaskConfigurationCrawlerDeep2: "2",
		TaskConfigurationCrawlerDeep3: "3",
		TaskConfigurationCrawlerDeep4: "4",
		TaskConfigurationCrawlerDeep5: "5",
		TaskConfigurationCrawlerDeep0: "不限",
	}
	return enum, TaskConfigurationCrawlerDefault
}
func (WebCrawlerConfig *WebCrawlerConfig) CrawlerDeep(crawler int) string {
	enum, _ := WebCrawlerConfig.AllCrawlerDeep()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

// 连接总数
const (
	TaskConfigurationCrawlerMaxUrl0      = 0
	TaskConfigurationCrawlerMaxUrl100    = 100
	TaskConfigurationCrawlerMaxUrl200    = 200
	TaskConfigurationCrawlerMaxUrl500    = 500
	TaskConfigurationCrawlerMaxUrl1000   = 1000
	TaskConfigurationCrawlerMaxUrl5000   = 5000
	TaskConfigurationCrawlerMaxUrl10000  = 10000
	TaskConfigurationCrawlerMaxUrl100000 = 100000

	TaskConfigurationCrawlerMaxUrlDefault = 1000 // 默认选择 - 下标
)

func (WebCrawlerConfig *WebCrawlerConfig) AllCrawlerMaxUrl() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationCrawlerMaxUrl100:    "100",
		TaskConfigurationCrawlerMaxUrl200:    "200",
		TaskConfigurationCrawlerMaxUrl500:    "500",
		TaskConfigurationCrawlerMaxUrl1000:   "1000",
		TaskConfigurationCrawlerMaxUrl5000:   "5000",
		TaskConfigurationCrawlerMaxUrl10000:  "10000",
		TaskConfigurationCrawlerMaxUrl100000: "100000",
		TaskConfigurationCrawlerMaxUrl0:      "不限",
	}
	return enum, TaskConfigurationCrawlerMaxUrlDefault
}
func (WebCrawlerConfig *WebCrawlerConfig) CrawlerMaxConnect(crawler int) string {
	enum, _ := WebCrawlerConfig.AllCrawlerMaxUrl()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

// url去重
// 0 不限，请勿设置爬虫过滤
// 1 `newcrawlerx.UnLimitRepeat` 对page，method，query-name，query-value和post-data敏感
// 2 `newcrawlerx.LowRepeatLevel` 对page，method，query-name和query-value敏感（默认）
// 3 `newcrawlerx.MediumRepeatLevel` 对page，method和query-name敏感
// 4 `newcrawlerx.HighRepeatLevel` 对page和method敏感
// 5 `newcrawlerx.ExtremeRepeatLevel` 对page敏感
const (
	TaskConfigurationCrawlerUrlRepeatNo     = 0
	TaskConfigurationCrawlerUrlRepeatLevel1 = 1
	TaskConfigurationCrawlerUrlRepeatLevel2 = 2
	TaskConfigurationCrawlerUrlRepeatLevel3 = 3
	TaskConfigurationCrawlerUrlRepeatLevel4 = 4
	TaskConfigurationCrawlerUrlRepeatLevel5 = 5

	TaskConfigurationCrawlerUrlRepeatLevelDefault = 3 // 默认选择 - 下标
)

func (WebCrawlerConfig *WebCrawlerConfig) AllCrawlerRepeat() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationCrawlerUrlRepeatLevel1: "对page，method，query-name，query-value和post-data敏感",
		TaskConfigurationCrawlerUrlRepeatLevel2: "对page，method，query-name和query-value敏感（默认）",
		TaskConfigurationCrawlerUrlRepeatLevel3: "对page，method和query-name敏感",
		TaskConfigurationCrawlerUrlRepeatLevel4: "对page和method敏感",
		TaskConfigurationCrawlerUrlRepeatLevel5: "对page敏感",
		TaskConfigurationCrawlerUrlRepeatNo:     "不限",
	}
	return enum, TaskConfigurationCrawlerUrlRepeatLevelDefault
}
func (WebCrawlerConfig *WebCrawlerConfig) CrawlerRepeat(crawler int) string {
	enum, _ := WebCrawlerConfig.AllCrawlerRepeat()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

func (WebCrawlerConfig *WebCrawlerConfig) AllCrawlerSpeed() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationCrawlerSpeedHigh:   "高速",
		TaskConfigurationCrawlerSpeedMiddle: "中速",
		TaskConfigurationCrawlerSpeedLow:    "低速",
	}
	return enum, TaskConfigurationCrawlerSpeedMiddle
}
func (WebCrawlerConfig *WebCrawlerConfig) CrawlerSpeed(crawler int) string {
	enum, _ := WebCrawlerConfig.AllCrawlerRepeat()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

const (
	TaskConfigurationCrawlerSpeedHigh   = 1 // 高速
	TaskConfigurationCrawlerSpeedMiddle = 2 // 中速
	TaskConfigurationCrawlerSpeedLow    = 3 // 低速
)

// 单链接超时
const (
	TaskConfigurationCrawlerSingleTimeout2   = 2
	TaskConfigurationCrawlerSingleTimeout5   = 5
	TaskConfigurationCrawlerSingleTimeout10  = 10
	TaskConfigurationCrawlerSingleTimeout20  = 20
	TaskConfigurationCrawlerSingleTimeout30  = 30
	TaskConfigurationCrawlerSingleTimeout60  = 60
	TaskConfigurationCrawlerSingleTimeout120 = 120

	TaskConfigurationCrawlerSingleTimeoutDefault = 10 // 默认选择 - 下标
)

func (WebCrawlerConfig *WebCrawlerConfig) AllSingleTimeout() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationCrawlerSingleTimeout2:   "2s",
		TaskConfigurationCrawlerSingleTimeout5:   "5s",
		TaskConfigurationCrawlerSingleTimeout10:  "10s",
		TaskConfigurationCrawlerSingleTimeout20:  "20s",
		TaskConfigurationCrawlerSingleTimeout30:  "30s",
		TaskConfigurationCrawlerSingleTimeout60:  "60s",
		TaskConfigurationCrawlerSingleTimeout120: "120s",
	}
	return enum, TaskConfigurationCrawlerSingleTimeoutDefault
}
func (WebCrawlerConfig *WebCrawlerConfig) SingleTimeout(crawler int) string {
	enum, _ := WebCrawlerConfig.AllSingleTimeout()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

// 单链接超时
const (
	TaskConfigurationCrawlerFullTimeout60   = 60
	TaskConfigurationCrawlerFullTimeout120  = 120
	TaskConfigurationCrawlerFullTimeout300  = 300
	TaskConfigurationCrawlerFullTimeout600  = 600
	TaskConfigurationCrawlerFullTimeout1200 = 1200
	TaskConfigurationCrawlerFullTimeout1800 = 1800
	TaskConfigurationCrawlerFullTimeout3600 = 3600

	TaskConfigurationCrawlerFullTimeoutDefault = 1800
)

func (WebCrawlerConfig *WebCrawlerConfig) AllFullTimeoutEnum() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationCrawlerFullTimeout60:   "60s",
		TaskConfigurationCrawlerFullTimeout120:  "120s",
		TaskConfigurationCrawlerFullTimeout300:  "300s",
		TaskConfigurationCrawlerFullTimeout600:  "600s",
		TaskConfigurationCrawlerFullTimeout1200: "1200s",
		TaskConfigurationCrawlerFullTimeout1800: "1800s",
		TaskConfigurationCrawlerFullTimeout3600: "3600s",
	}
	return enum, TaskConfigurationCrawlerFullTimeoutDefault
}
func (WebCrawlerConfig *WebCrawlerConfig) FullTimeoutEnum(crawler int) string {
	enum, _ := WebCrawlerConfig.AllFullTimeoutEnum()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

/************* Web路径爆破配置 **************/
// 猜测速率
const (
	TaskConfigurationWebPathScanSpeedHigh   = 1
	TaskConfigurationWebPathScanSpeedMiddle = 2
	TaskConfigurationWebPathScanSpeedLow    = 3

	TaskConfigurationWebPathScanSpeedDefault = 2
)

func (WebPathScanConfig *WebPathScanConfig) AllWebPathScanSpeed() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationWebPathScanSpeedHigh:   "高速",
		TaskConfigurationWebPathScanSpeedMiddle: "中速",
		TaskConfigurationWebPathScanSpeedLow:    "低速",
	}
	return enum, TaskConfigurationWebPathScanSpeedDefault
}
func (WebPathScanConfig *WebPathScanConfig) WebPathScanSpeed(crawler int) string {
	enum, _ := WebPathScanConfig.AllWebPathScanSpeed()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

// 猜测时长
const (
	TaskConfigurationWebPathScanTime0  = 0
	TaskConfigurationWebPathScanTime1  = 1
	TaskConfigurationWebPathScanTime3  = 3
	TaskConfigurationWebPathScanTime5  = 4
	TaskConfigurationWebPathScanTime10 = 10
	TaskConfigurationWebPathScanTime30 = 30
	TaskConfigurationWebPathScanTime60 = 60

	TaskConfigurationWebPathScanTimeDefault = 10
)

func (WebPathScanConfig *WebPathScanConfig) AllWebPathScanTime() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationWebPathScanTime1:  "1min",
		TaskConfigurationWebPathScanTime3:  "3min",
		TaskConfigurationWebPathScanTime5:  "5min",
		TaskConfigurationWebPathScanTime10: "10min",
		TaskConfigurationWebPathScanTime30: "30min",
		TaskConfigurationWebPathScanTime60: "60min",
		TaskConfigurationWebPathScanTime0:  "不限",
	}
	return enum, TaskConfigurationWebPathScanTimeDefault
}
func (WebPathScanConfig *WebPathScanConfig) WebPathScanTime(crawler int) string {
	enum, _ := WebPathScanConfig.AllWebPathScanTime()
	if res, ok := enum[crawler]; ok {
		return res
	}
	return ""
}

/************* Web弱口令配置 **************/
// 字典类型
const (
	TaskConfigurationWeakPassDictTypeDefault = 1 // 默认字典
	TaskConfigurationWeakPassDictTypeCommon  = 2 // 通用字典
	TaskConfigurationWeakPassDictTypeAdd     = 3 // 补充字典
)

func (WeakPassConfig *WeakPassConfig) AllWeakPassDictType() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationWeakPassDictTypeDefault: "默认字典",
		TaskConfigurationWeakPassDictTypeCommon:  "通用字典",
		TaskConfigurationWeakPassDictTypeAdd:     "补充字典",
	}
	return enum, TaskConfigurationWeakPassDictTypeDefault
}
func (WeakPassConfig *WeakPassConfig) WeakPassDictType(dictType int) string {
	enum, _ := WeakPassConfig.AllWeakPassDictType()
	if res, ok := enum[dictType]; ok {
		return res
	}
	return ""
}

// 猜测次数
const (
	TaskConfigurationWeakPassGuessNumber0     = 0
	TaskConfigurationWeakPassGuessNumber1     = 1
	TaskConfigurationWeakPassGuessNumber3     = 3
	TaskConfigurationWeakPassGuessNumber5     = 5
	TaskConfigurationWeakPassGuessNumber10    = 10
	TaskConfigurationWeakPassGuessNumber20    = 20
	TaskConfigurationWeakPassGuessNumber30    = 30
	TaskConfigurationWeakPassGuessNumber50    = 50
	TaskConfigurationWeakPassGuessNumber100   = 100
	TaskConfigurationWeakPassGuessNumber1000  = 1000
	TaskConfigurationWeakPassGuessNumber5000  = 5000
	TaskConfigurationWeakPassGuessNumber10000 = 10000
)

func (WeakPassConfig *WeakPassConfig) AllWeakPassGuessNumber() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationWeakPassGuessNumber0:     "0",
		TaskConfigurationWeakPassGuessNumber1:     "1",
		TaskConfigurationWeakPassGuessNumber3:     "3",
		TaskConfigurationWeakPassGuessNumber5:     "5",
		TaskConfigurationWeakPassGuessNumber10:    "10",
		TaskConfigurationWeakPassGuessNumber20:    "20",
		TaskConfigurationWeakPassGuessNumber30:    "30",
		TaskConfigurationWeakPassGuessNumber50:    "50",
		TaskConfigurationWeakPassGuessNumber100:   "100",
		TaskConfigurationWeakPassGuessNumber1000:  "1000",
		TaskConfigurationWeakPassGuessNumber5000:  "5000",
		TaskConfigurationWeakPassGuessNumber10000: "10000",
	}
	return enum, TaskConfigurationWeakPassGuessNumber1000
}
func (WeakPassConfig *WeakPassConfig) WeakPassGuessNumber(guessNumber int) string {
	enum, _ := WeakPassConfig.AllWeakPassGuessNumber()
	if res, ok := enum[guessNumber]; ok {
		return res
	}
	return ""
}

// 猜测时间 单位分钟
const (
	TaskConfigurationWeakPassGuessTime0  = 0 // 不限
	TaskConfigurationWeakPassGuessTime1  = 1
	TaskConfigurationWeakPassGuessTime3  = 3
	TaskConfigurationWeakPassGuessTime5  = 5
	TaskConfigurationWeakPassGuessTime10 = 10
	TaskConfigurationWeakPassGuessTime30 = 30
	TaskConfigurationWeakPassGuessTime60 = 60
)

func (WeakPassConfig *WeakPassConfig) AllWeakPassGuessTime() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationWeakPassGuessTime1:  "1min",
		TaskConfigurationWeakPassGuessTime3:  "3min",
		TaskConfigurationWeakPassGuessTime5:  "5min",
		TaskConfigurationWeakPassGuessTime10: "10min",
		TaskConfigurationWeakPassGuessTime30: "30min",
		TaskConfigurationWeakPassGuessTime60: "60min",
		TaskConfigurationWeakPassGuessTime0:  "不限",
	}
	return enum, TaskConfigurationWeakPassGuessTime5
}
func (WeakPassConfig *WeakPassConfig) WeakPassGuessTime(guessNumber int) string {
	enum, _ := WeakPassConfig.AllWeakPassGuessTime()
	if res, ok := enum[guessNumber]; ok {
		return res
	}
	return ""
}

// 猜测速率
const (
	TaskConfigurationWeakPassGuessRateHigh = 1
	TaskConfigurationWeakPassGuessRateMid  = 2
	TaskConfigurationWeakPassGuessRateLow  = 3
)

func (WeakPassConfig *WeakPassConfig) AllWeakPassGuessRate() (map[int]string, int) {
	enum := map[int]string{
		TaskConfigurationWeakPassGuessRateHigh: "高速",
		TaskConfigurationWeakPassGuessRateMid:  "中速",
		TaskConfigurationWeakPassGuessRateLow:  "低速",
	}
	return enum, TaskConfigurationWeakPassGuessRateMid
}
func (WeakPassConfig *WeakPassConfig) WeakPassGuessRate(guessRate int) string {
	enum, _ := WeakPassConfig.AllWeakPassGuessRate()
	if res, ok := enum[guessRate]; ok {
		return res
	}
	return ""
}

/************* 网站登陆凭证 ****************/
// 凭证类型
const (
	TaskConfigurationWebsiteLoginHeader        = 1
	TaskConfigurationWebsiteLoginCookie        = 2
	TaskConfigurationWebsiteLoginAccount       = 3
	TaskConfigurationWebsiteLoginLoginSequence = 4
)

func (WebLogin *WebsiteLoginConfig) AllWebsiteLoginType() map[int]string {
	enum := map[int]string{
		TaskConfigurationWebsiteLoginHeader:  "Header",
		TaskConfigurationWebsiteLoginCookie:  "Cookie",
		TaskConfigurationWebsiteLoginAccount: "账号密码",
		// TaskConfigurationWebsiteLoginLoginSequence: "登录序列",
	}
	return enum
}
func (WebLogin *WebsiteLoginConfig) WebsiteLoginType(types int) string {
	enum := WebLogin.AllWebsiteLoginType()
	if res, ok := enum[types]; ok {
		return res
	}
	return ""
}

// 网站登陆验证状态
const (
	TaskConfigurationWebsiteLoginVerify  = 1
	TaskConfigurationWebsiteLoginSuccess = 2
	TaskConfigurationWebsiteLoginFail    = 3
)

func (WebLogin *WebsiteLoginConfig) VerifyStatusEnum() map[int]string {
	enum := map[int]string{
		TaskConfigurationWebsiteLoginVerify:  "认证成功",
		TaskConfigurationWebsiteLoginSuccess: "访问成功",
		TaskConfigurationWebsiteLoginFail:    "访问失败",
	}
	return enum
}

func (WebLogin *WebsiteLoginConfig) VerifyStatusZh(k int) string {
	enum := WebLogin.VerifyStatusEnum()

	value, ok := enum[k]
	if ok {
		return value
	}

	return ""
}

// 横向移动
const (
	LateralMoveJumpNumZero  = 0
	LateralMoveJumpNumOne   = 1
	LateralMoveJumpNumTwo   = 2
	LateralMoveJumpNumThree = 3
	LateralMoveJumpNumFour  = 4
	LateralMoveJumpNumFive  = 5
	LateralMoveJumpNumSix   = 6
	LateralMoveJumpNumSeven = 7
	LateralMoveJumpNumEight = 8
)

// 横向移动策略
const (
	LateralMoveStrategySameSubnet     = "same_subnet"     // 同网段探测
	LateralMoveStrategyNeighbor       = "neighbor"        // 邻居发现
	LateralMoveStrategyExcludeCurrent = "exclude_current" // 排除同网段
	LateralMoveStrategyCustomRange    = "custom_range"    // 自定义范围
)

// 供给面类型
type AttackFace struct{}

const (
	AttackFaceTypePort      = 1 //端口
	AttackFaceTypePath      = 2 //敏感路径
	AttackFaceTypeLoginCred = 3 //文件上传
)

func (a *AttackFace) AttackFaceTypeEnum() map[int]string {
	enum := map[int]string{
		AttackFaceTypePort:      "端口",
		AttackFaceTypePath:      "敏感路径",
		AttackFaceTypeLoginCred: "登录凭证",
	}
	return enum
}

func (a *AttackFace) GetAttackFaceTypeEnum(aftype int) string {
	enum := a.AttackFaceTypeEnum()
	if res, ok := enum[aftype]; ok {
		return res
	}
	return ""
}

func (a *AttackFace) GetAttackFaceTypeEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: AttackFaceTypePort,
		Label: a.GetAttackFaceTypeEnum(AttackFaceTypePort),
	}, {
		Value: AttackFaceTypePath,
		Label: a.GetAttackFaceTypeEnum(AttackFaceTypePath),
	}, {
		Value: AttackFaceTypeLoginCred,
		Label: a.GetAttackFaceTypeEnum(AttackFaceTypeLoginCred),
	}}
	return result
}

// 代理模式协议枚举 ProxyConfig
// 代理模式协议枚举 ProxyConfig
const (
	ProxyConfigProtoHTTP   = 1
	ProxyConfigProtoHTTPS  = 2
	ProxyConfigProtoSOCKS4 = 3
	ProxyConfigProtoSOCKS5 = 4
)

// 注意这里增加其他协议时需要遵守一定格式比如
// socks5://127.0.0.1那么枚举值必须是socks5
// http://127.0.0.1那么枚举值必须是http 不允许有其他字段
func (proxy ProxyConfig) AllProxyConfigProtoEnum() map[int]string {
	return map[int]string{
		ProxyConfigProtoHTTP:   "http",
		ProxyConfigProtoHTTPS:  "https",
		ProxyConfigProtoSOCKS4: "socks4",
		ProxyConfigProtoSOCKS5: "socks5",
	}
}

func (proxy ProxyConfig) DecisionProxyConfigProtoEnum() map[int]string {
	return map[int]string{
		ProxyConfigProtoHTTP:   "http",
		ProxyConfigProtoHTTPS:  "https",
		ProxyConfigProtoSOCKS4: "socks4",
		ProxyConfigProtoSOCKS5: "socks5",
	}
}

func (proxy ProxyConfig) GetProxyConfigProtoEnum(proto int) string {
	return proxy.DecisionProxyConfigProtoEnum()[proto]
}

// 网站登陆验证状态
const (
	Sensitive = "Update Profile,LogOut"
	WhiteList = ""
	BlackList = ""
)

// 测试模式配置 ProxyConfig
const (
	TestModeConfigPrinciple = 1
	TestModeConfigVersion   = 2
)

const (
	PenetrationModeConfigCommon    = 1
	PenetrationModeConfigDirection = 2
)

func (testMode TestModeConfig) AllTestModeConfigEnum() map[int]string {
	return map[int]string{
		TestModeConfigPrinciple: "原理验证",
		TestModeConfigVersion:   "版本匹配",
	}
}
func (penetrationMode PenetrationModeConfig) AllPenetrationModeConfigEnum() map[int]string {
	return map[int]string{
		PenetrationModeConfigCommon:    "通用渗透",
		PenetrationModeConfigDirection: "定向渗透",
	}
}
