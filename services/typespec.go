package services

import (
	"github.com/yaklang/yaklang/common/netx"
	"smart/tools/enums"
	"time"
)

type SceneTaskTemplateCreateVulIdsItem struct {
	Id     int    `form:"id" json:"id" binding:"required"`
	Origin string `form:"origin" json:"origin" binding:"required"`
}

// 任务创建消息
type TaskControlMessage struct {
	ObjId           string                        `json:"objId"`     //目标id
	Operate         string                        `json:"operate"`   //对任务的操作
	CallTools       []string                      `json:"callTools"` //执行工具
	TargetUrl       string                        `json:"targetUrl"` //目标url
	CallVuls        []int                         `json:"callVuls"`  //执行工具
	InfoList        []Info                        `json:"infoList"`
	Info            string                        `json:"info"`            //数据冗余字段json,发送过来数据带着task_id/targert_id/log_id等信息，可以在处理结果时候方便回表查询
	TestIntensity   string                        `json:"testIntensity"`   //测试强调
	SafeTest        bool                          `json:"safeTest"`        //是否开启安全测试
	LateralMove     TaskControlMessageLateralMove `json:"lateralMove"`     //横向移动
	VulExploit      bool                          `json:"vulExploit"`      //是否开启漏洞利用
	SafeTestPocname string                        `json:"safeTestPocname"` //执行安全测试的pocname
	SafeTestId      int                           `json:"safeTestId"`      //执行安全测试的漏洞id
	Mode            enums.Mode                    `json:"mode"`            // 任务运行模式
	Proxy           string                        `json:"proxy"`           // 任务代理配置
	TaskId          int                           `json:"taskId"`          // 任务id
}

type TaskControlMessageLateralMove struct {
	Range  string `json:"range"`
	IsOpen bool   `json:"isOpen"`
}

type TaskControlMessageInfo struct {
	TaskId   int `json:"taskId"`
	TargetId int `json:"targetId"`
	LogId    int `json:"logId"`
}

// 任务信息数据
type Info struct {
	Name  string `json:"name"`  //key值
	Value string `json:"value"` //value值
}

/********************* 创建渗透任务时的参数 ***********************/
// 配置数据
type TaskRunConf struct {
	StartTime          string   `json:"startTime"`
	CyclePlanningType  int      `json:"cyclePlanningType"`
	CyclePlanningValue int      `json:"cyclePlanningValue"`
	CyclePlanningHour  string   `json:"cyclePlanningHour"`
	EndTime            string   `json:"endTime"`
	RuntimePeriod      []string `json:"runtimePeriod"`
}

// 脚本结果数据
type ScriptResult struct {
	ObjId      string   `json:"objId"`       // 任务ID (渗透任务时是目标ID即target_id)
	NodeId     string   `json:"nodeId"`      // 结果ID
	FatherId   string   `json:"fatherId"`    // 结果父级ID
	TaskStatus int      `json:"task_status"` // 状态 created 创建  running 运行中  finished 已结束  paused 暂停中
	Status     int      `json:"status"`      // 状态 1开始 2成功 3失败 4结束
	Script     struct { // 脚本信息
		DataType     int    `json:"dataType"`     //漏洞数据类型，1或0是漏洞测试，2-待测漏洞
		Pocname      string `json:"pocname"`      // pocname
		Type         string `json:"type"`         // 工具类型 yak | python | mitm | msf | ...
		VerifyType   string `json:"verifyType"`   // 校验类型 poc | exp | ...
		EvidenceType int    `json:"evidenceType"` // 取证类型
		ScriptParam  string `json:"scriptParam"`  //工具调用参数
		Source       string `json:"source"`       //漏洞来源,vultest-待测漏洞手动触发
		SafeTestId   int    `json:"safeTestId"`   //待测漏洞id
	} `json:"script"`
	Libraries struct { // 漏洞信息
		Id                int    `json:"id"`                // 漏洞Id
		VulId             string `json:"vulId"`             // 决策引擎唯一漏洞Id
		Name              string `json:"name"`              // 漏洞名称
		Risk              int    `json:"risk"`              // 漏洞风险
		Type              int    `json:"type"`              // 漏洞类型
		Class             int    `json:"class"`             // 漏洞分类
		PublishedTime     string `json:"publishedTime"`     // 公开时间
		Description       string `json:"description"`       // 漏洞描述
		AffectRange       string `json:"affectRange"`       // 影响范围
		ExploitImpact     int    `json:"exploitImpact"`     // 利用影响
		ExploitImpactEnum string `json:"exploitImpactEnum"` // 利用影响
		FixSuggest        string `json:"fixSuggest"`        // 修复建议
		Cve               string `json:"cve"`               // cve编号
		Cnvd              string `json:"cnvd"`              // cnvd编号
		Cnnvd             string `json:"cnnvd"`             // cnnvd编号
		Component         string `json:"component"`         // 受影响的组建名称
	} `json:"libraries"`
	Result struct { // 结果数据
		Request   string `json:"request"`   // 请求信息
		Response  string `json:"response"`  // 响应信息
		Location  string `json:"location"`  // 漏洞位置
		Detail    string `json:"detail"`    // 其他描述
		TargetUrl string `json:"targetUrl"` // 其他描述
	} `json:"result"`
}

// 任务概览
type OverView struct {
	TargetRisk       [4]int           `json:"targetRisk"`
	TargetNum        []int            `json:"targetNum"`
	VulTotal         int              `json:"vulTotal"`
	VulRisk          []int            `json:"vulRisk"`
	RiskLevel        int              `json:"riskLevel"`
	VulExploitImpact []*OverViewItems `json:"vulExploitImpact"`
	EvidenceStat     []*OverViewItems `json:"evidenceStat"`
	Port             []*OverViewItems `json:"port"`
	Service          []*OverViewItems `json:"service"`
	Component        []*OverViewItems `json:"component"`
	OpSys            []*OverViewItems `json:"opSys"`
	SubDomain        []*OverViewItems `json:"subDomain"`
	UrlTags          []*OverViewItems `json:"urlTags"`
	Progress         *OverViewItems   `json:"progress"`
	VerifyType       []*OverViewItems `json:"verifyType"`
}
type OverViewItems struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type ReportContentMapSet struct {
	TaskContent   []ReportContentItem `json:"taskContent"`
	TargetContent []ReportContentItem `json:"targetContent"`
}
type ReportContentItem struct {
	Value string               `json:"value"`
	Label string               `json:"label"`
	Items []*ReportContentItem `json:"items"`
}

type TripartiteToolsXrayCreateResultItem struct {
	XrayTaskId  int       `json:"xrayTaskId"`
	Addr        string    `json:"addr"`
	Payload     string    `json:"payload"`
	RequestInfo string    `json:"requestInfo"`
	Extra       string    `json:"extra"`
	Plugin      string    `json:"plugin"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

// wifi创建参数
type WifiCreateData struct {
	TaskName         string `json:"taskName"`
	Mac              string `json:"mac"`
	Channel          int    `json:"channel"`
	Encrypt          int    `json:"encrypt"`
	Carrier          int    `json:"carrier"`
	Status           int    `json:"status"`
	Ssid             string `json:"ssid"`
	Passwd           string `json:"passwd"`
	PasswdSource     int    `json:"passwdSource"`
	StartTime        int64  `json:"startTime"`
	EndTime          int64  `json:"endTime"`
	IsSimulate       int    `json:"isSimulate"`
	SimulateDuration int    `json:"simulateDuration"`
	IsCrack          int    `json:"isCrack"`
	PasswdDict       string `json:"passwdDict"`
	IsEmbed          int    `json:"IsEmbed"`
}
type TargetIpWhiteBlackMapSet struct {
	IsOpen      int      `json:"isOpen"`
	Type        int      `json:"type"`
	WhiteList   string   `json:"whiteList"`
	BlackList   string   `json:"blackList"`
	IpListArray []string `json:"ipListArray"`
}

type LogBackupConfigMapSet struct {
	IsOpen   int       `json:"isOpen"`
	Cycle    int       `json:"cycle"`
	SaveTime time.Time `json:"saveTime"`
	RunTime  time.Time `json:"runTime"`
}

type SystemConfigBackupConfigMapSet struct {
	IsOpen   int       `json:"isOpen"`
	Cycle    int       `json:"cycle"`
	SaveTime time.Time `json:"saveTime"`
	RunTime  time.Time `json:"runTime"`
}

type ReverseIpHost struct {
	ReverseType int    `json:"reverseType"`
	ReverseHost string `json:"reverseHost"`
	ReversePort int    `json:"reversePort"`
}

// 三方工具 - burpsuite 请求http后的结果数据
type BurpsuiteResultData struct {
	ScanMetrics struct {
		CurrentUrl                  string `json:"current_url"`
		CrawlRequestsMade           int    `json:"crawl_requests_made"`
		CrawlNetworkErrors          int    `json:"crawl_network_errors"`
		CrawlUniqueLocationsVisited int    `json:"crawl_unique_locations_visited"`
		CrawlRequestsQueued         int    `json:"crawl_requests_queued"`
		AuditQueueItemsCompleted    int    `json:"audit_queue_items_completed"`
		AuditQueueItemsWaiting      int    `json:"audit_queue_items_waiting"`
		AuditRequestsMade           int    `json:"audit_requests_made"`
		AuditNetworkErrors          int    `json:"audit_network_errors"`
		IssueEvents                 int    `json:"issue_events"`
		CrawlAndAuditCaption        string `json:"crawl_and_audit_caption"`
		CrawlAndAuditProgress       int    `json:"crawl_and_audit_progress"`
		TotalElapsedTime            int    `json:"total_elapsed_time"`
	} `json:"scan_metrics"`
	IssueEvents []IssueEvents `json:"issue_events"`
	TaskId      string        `json:"task_id"`
	ScanStatus  string        `json:"scan_status"`
	Message     string        `json:"message"`
	ErrorCode   int           `json:"error_code"`
}
type IssueEvents struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Issue struct {
		Name                  string `json:"name"`
		TypeIndex             int    `json:"type_index"`
		SerialNumber          string `json:"serial_number"`
		Origin                string `json:"origin"`
		Path                  string `json:"path"`
		Severity              string `json:"severity"`
		Confidence            string `json:"confidence"`
		Description           string `json:"description"`
		IssueBackground       string `json:"issue_background"`
		RemediationBackground string `json:"remediation_background"`
		Caption               string `json:"caption"`
		Evidence              []struct {
			Type   string `json:"type"`
			Detail struct {
				BandFlags []string `json:"band_flags"`
			} `json:"detail"`
			RequestResponse struct {
				Url     string `json:"url"`
				Request []struct {
					Type   string `json:"type"`
					Data   string `json:"data"`
					Length int    `json:"length"`
				} `json:"request"`
				Response []struct {
					Type   string `json:"type"`
					Data   string `json:"data,omitempty"`
					Length int    `json:"length"`
				} `json:"response"`
				WasRedirectFollowed bool   `json:"was_redirect_followed"`
				RequestTime         string `json:"request_time"`
			} `json:"request_response"`
		} `json:"evidence"`
		InternalData string `json:"internal_data"`
	} `json:"issue"`
}

type SystemSettingIpWhiteMapSet struct {
	IsOpen int    `json:"isOpen"`
	Ip     string `json:"ip"`
}

type SystemSettingSyslogMapSet struct {
	IsOpen int    `json:"isOpen"`
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Types  string `json:"types"`
}

type SystemSettingMailMapSet struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Encrypt  string `json:"encrypt"`
}

type SystemSettingNetworkConfigMapSet struct {
	Ip               string `json:"ip"`
	Mask             string `json:"mask"`
	Gateway          string `json:"gateway"`
	DnsServer        string `json:"dnsServer"`
	StandbyDnsServer string `json:"standbyDnsServer"`
	WebPort          int    `json:"webPort"`
}

type MonitorWarnMapSet struct {
	IsOpen     int `json:"isOpen"`
	CpuWarn    int `json:"cpuWarn"`
	MemoryWarn int `json:"memoryWarn"`
	DiskWarn   int `json:"diskWarn"`
	FlowWarn   int `json:"flowWarn"`
}

type TcpBlindTestMapSet struct {
	Type int    `json:"type"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HttpBlindTestMapSet struct {
	Type int    `json:"type"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DnsBlindTestMapSet struct {
	Type   int    `json:"type"`
	Domain string `json:"domain"`
}

type IcmpBlindTestMapSet struct {
	Type int    `json:"type"`
	Host string `json:"host"`
}

type MonitorCpuMemoryMapSet struct {
	List []MonitorCpuMemoryItem `json:"list"`
}

type MonitorCpuMemoryItem struct {
	XData string `json:"xData"`
	YData int    `json:"yData"`
}

type MonitorDiskMapSet struct {
	Free        int `json:"free"`
	Used        int `json:"used"`
	Total       int `json:"total"`
	FreePercent int `json:"freePercent"`
	UsedPercent int `json:"usedPercent"`
}

type FlowResultNetflowMsg struct {
	Hash            string `json:"hash"`
	Host            string `json:"host"`
	Ip              string `json:"ip"`
	Method          string `json:"method"`
	Protocol        string `json:"protocol"`
	ReqHeader       string `json:"reqHeader"`
	RespCode        string `json:"respCode"`
	RespContent     string `json:"respContent"`
	RespContentType string `json:"respContentType"`
	RespHeader      string `json:"respHeader"`
	RespTitle       string `json:"respTitle"`
	Url             string `json:"url"`
}

type FlowResultTagsMsg struct {
	Hash   string `json:"hash"`
	Method string `json:"method"`
	Tags   string `json:"tags"`
	Url    string `json:"url"`
}

type FlowResultRiskMsg struct {
	Title           string `json:"title"`
	RiskLevel       int    `json:"riskLevel"`
	Host            string `json:"host"`
	Ip              string `json:"ip"`
	Port            string `json:"port"`
	Hash            string `json:"hash"`
	RiskTypeVerbose string `json:"riskTypeVerbose"`
	Parameter       string `json:"parameter"`
	Request         string `json:"request"`
	Response        string `json:"response"`
	Detail          string `json:"detail"`
}

type SystemVersionMapSet struct {
	CurrentVersion string `json:"currentVersion"` //当前系统版本
	//LatestVersion  string `json:"latestVersion"`
	//FindTime       string `json:"findTime"`
	UpdateTime    string `json:"updateTime"`    //系统更新时间
	VulUpdateTime string `json:"vulUpdateTime"` //工具库更新时间
	VulVersion    string `json:"vulVersion"`    //工具库版本
	SysFileName   string `json:"sysFileName"`   // 系统更新文件
	VulFileName   string `json:"VulFileName"`   //工具库更新文件
	LastSystemVersion string `json:"lastSystemVersion"` // 上次系统版本
	LastVulVersion    string `json:"lastVulVersion"`    // 上次工具库版本
}

// BAS 导入规则时的所有数据
type BasImportRules struct {
	Id                                 int       `json:"id"`
	Content                            string    `json:"content"`
	Name                               string    `json:"name"`
	NameZh                             string    `json:"nameZh"`
	ClassType                          string    `json:"classType"`
	Protocal                           string    `json:"protocal"`
	Keywords                           string    `json:"keywords"`
	KeywordsZH                         string    `json:"keywords_zh"`
	Description                        string    `json:"description"`
	DescriptionZH                      string    `json:"description_zh"`
	Cve                                string    `json:"cve"`
	RawTrafficBeyondIpPacketBase64Name string    `json:"rawTrafficBeyondIpPacketBase64Name"`
	RawTrafficBeyondHttpBase64Name     string    `json:"rawTrafficBeyondHttpBase64Name"`
	Hash                               string    `json:"hash"`
	CreateTime                         time.Time `json:"createTime"`
	UpdateTime                         time.Time `json:"updateTime"`
}

// BAS 创建目标日志
type BasTaskLogCreate struct {
	BasTaskId       int    `json:"basTaskId"`
	BasTaskTargetId int    `json:"basTaskTargetId"`
	Content         string `json:"content"`
}

type DockerComposeConfig struct {
	Networks interface{}           `json:"networks" yaml:"networks"`
	Services DockerComposeServices `json:"services" yaml:"services"`
	Version  string                `json:"version" yaml:"version"`
}
type DockerComposeServices struct {
	Db       interface{}                `json:"db" yaml:"db"`
	GoRod    interface{}                `json:"go-rod" yaml:"go-rod"`
	Neo4j    interface{}                `json:"neo4j" yaml:"neo4j"`
	Nginx    DockerComposeServicesNginx `json:"nginx" yaml:"nginx"`
	Rabbitmq interface{}                `json:"rabbitmq" yaml:"rabbitmq"`
	Redis    interface{}                `json:"redis" yaml:"redis"`
	Msf      interface{}                `json:"msf" yaml:"msf"`
}

type DockerComposeServicesNginx struct {
	Command       interface{} `json:"command" yaml:"command"`
	ContainerName interface{} `json:"container_name" yaml:"container_name"`
	Image         interface{} `json:"image" yaml:"image"`
	Networks      interface{} `json:"networks" yaml:"networks"`
	Ports         []string    `json:"ports" yaml:"ports"`
	Restart       interface{} `json:"restart" yaml:"restart"`
	Volumes       interface{} `json:"volumes" yaml:"volumes"`
}

type VulResult struct {
	Password string `json:"password"`
	Service  string `json:"service"`
	Target   string `json:"target"`
	Username string `json:"username"`
	Content  string `json:"content"`
	DbType   string `json:"dbType"`
}

// 指纹识别结果结构体
type FingerPrintResult struct {
	Port       string       `json:"port"`
	TargetInfo []TargetInfo `json:"target_info"`
}

type TargetInfo struct {
	AppClass   string `json:"app_class"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	CnName     string `json:"cn_name"`
	Flag       string `json:"flag"`
	IsDev      bool   `json:"is_dev"`
	Level      string `json:"level"` // 层级
}

type targetUrlMapItems struct {
	TargetId  int    `json:"targetId"`
	TargetUrl string `json:"targetUrl"`
}

type RemoteSessionFile struct {
	FileName string `json:"fileName"`
	FileSize string `json:"fileSize"`
	FilePath string `json:"filePath"`
}

type BasRuleFile struct {
	RuleId  int    `json:"ruleId"`
	Content string `json:"content"`
}

type BaSendMsg struct {
	Content []BaSendMsgContent `json:"content"`
	Type    string             `json:"type"`
}

type BaSendMsgContent struct {
	RuleId int      `json:"ruleId"`
	Md5    []string `json:"md5"`
}

type ScriptResultDetailVulProve struct {
	Request            string `json:"request"`
	Response           string `json:"response"`
	Payload            string `json:"payload"`
	PayloadSuccessFlag string `json:"payload_success_flag"`
}

// AssetImportData 资产导入数据格式
type AssetImportData struct {
	Id          int       `json:"id"`          // 域名/IP
	DomainIp    string    `json:"domainIp"`    // 域名/IP
	OpSys       string    `json:"opSys"`       // 操作系统
	Hardware    string    `json:"hardware"`    //硬件
	Hostname    string    `json:"hostname"`    //主机名
	Group       string    `json:"group"`       //资产组
	AssetName   string    `json:"assetName"`   //资产名称
	Business    string    `json:"business"`    //业务系统
	AssetType   string    `json:"assetType"`   //资产类型
	Virtual     string    `json:"virtual"`     //虚拟资产
	RecordLevel int       `json:"recordLevel"` //备案等级
	Location    string    `json:"location"`    //归属地
	Department  string    `json:"department"`  //部门
	Person      string    `json:"person"`      //责任人
	Email       string    `json:"email"`       //邮箱
	Tags        string    `json:"tags"`        //标签
	CreateTime  time.Time `json:"createTime"`  //创建时间
	UpdateTime  time.Time `json:"updateTime"`  //标签
}

// 越权配置信息解析
type BeyondPermConfig struct {
	LoginCred struct {
		Pattern int    `json:"pattern"`
		Value   string `json:"value"`
	} `json:"loginCred"`
	WaitCred struct {
		Pattern int    `json:"pattern"`
		Value   string `json:"value"`
	} `json:"waitCred"`
	WhitePath string `json:"whitePath"`
	BlackPath string `json:"blackPath"`
	Keywords  string `json:"keywords"`
	Crawler   struct {
		Range      int      `json:"range"`
		Depth      int      `json:"depth"`
		MaxLink    int      `json:"maxLink"`
		SingleLink int      `json:"singleLink"`
		Sensitive  []string `json:"sensitive"`
		BlackWord  []string `json:"blackWord"`
		WhiteWord  []string `json:"whiteWord"`
	} `json:"crawler"`
}

// 敏感信息遍历 配置信息
type SensInfoConfig struct {
	LoginCred struct {
		Pattern int    `json:"pattern"`
		Value   string `json:"value"`
	} `json:"loginCred"`
	FuzzParam struct {
		Character string `json:"character"`
		Number    string `json:"number"`
	}
	FuzzDict struct {
		Character string `json:"character"`
		Number    string `json:"number"`
	}
	Response struct {
		JsonKeyword   string `json:"jsonKeyword"`
		NoJsonSwitch  bool   `json:"noJsonSwitch"`
		NoJsonKeyword string `json:"noJsonKeyword"`
	}
	Crawler struct {
		Range      int      `json:"range"`
		Depth      int      `json:"depth"`
		MaxLink    int      `json:"maxLink"`
		SingleLink int      `json:"singleLink"`
		Sensitive  []string `json:"sensitive"`
		BlackWord  []string `json:"blackWord"`
		WhiteWord  []string `json:"whiteWord"`
	} `json:"crawler"`
}

// 未授权访问 配置信息解析
type UnAuthAccessConfig struct {
	LoginCred struct {
		Pattern int    `json:"pattern"`
		Value   string `json:"value"`
	} `json:"loginCred"`
	WhitePath        string `json:"whitePath"`
	BlackPath        string `json:"blackPath"`
	Keywords         string `json:"keywords"`
	CredIdentifyList string `json:"credIdentifyList"`
	Crawler          struct {
		Range      int      `json:"range"`
		Depth      int      `json:"depth"`
		MaxLink    int      `json:"maxLink"`
		SingleLink int      `json:"singleLink"`
		Sensitive  []string `json:"sensitive"`
		BlackWord  []string `json:"blackWord"`
		WhiteWord  []string `json:"whiteWord"`
	} `json:"crawler"`
}

type PortState string

func (p *PortState) String() string {
	return string(*p)
}

var (
	OPEN    PortState = "open"
	CLOSED  PortState = "closed"
	UNKNOWN PortState = "unknown"
)

type MatchResult struct {
	Target      string           `json:"target"`
	Port        int              `json:"port"`
	State       PortState        `json:"state"`
	Reason      string           `json:"reason"`
	Fingerprint *FingerprintInfo `json:"fingerprint"`
}

type PortScanResultLog struct {
	FingerPrintResult *PortScanFingerPrintResult `json:"finger_print_result"`
}

type PortScanFingerPrintResult struct {
	ServiceInfo       []*PortScanServiceInfo `json:"service_info"`
	TLSInspectResults interface{}            `json:"TLSInspectResults"`
}

type PortScanServiceInfo struct {
	Host        string      `json:"host"`
	Port        string      `json:"port"`
	ServiceName string      `json:"service_name"`
	ProductName string      `json:"product_name"`
	Version     string      `json:"version"`
	IsWeb       int         `json:"is_web"`
	CpeList     interface{} `json:"cpe_list"`
}

type FingerprintInfo struct {
	IP               string            `json:"ip"`
	Port             int               `json:"port"`
	Proto            TransportProto    `json:"proto"`
	ServiceName      string            `json:"service_name"`
	ProductVerbose   string            `json:"product_verbose"`
	Info             string            `json:"info"`
	Version          string            `json:"version"`
	Hostname         string            `json:"hostname"`
	OperationVerbose string            `json:"operation_verbose"`
	DeviceType       string            `json:"device_type"`
	CPEs             []string          `json:"cpes"`
	Raw              string            `json:"raw"`
	Banner           string            `json:"banner"`
	CPEFromUrls      map[string][]*CPE `json:"cpe_from_urls"`
	HttpFlows        []*HTTPFlow       `json:"http_flows"`

	// tls info for fill...
	TLSInspectResults []*netx.TLSInspectResult
}

type TransportProto string

var (
	TCP TransportProto = "tcp"
	UDP TransportProto = "udp"
)

type HTTPFlow struct {
	StatusCode     int    `json:"status_code"`
	IsHTTPS        bool   `json:"is_https"`
	RequestHeader  []byte `json:"request_header"`
	RequestBody    []byte `json:"request_body"`
	ResponseHeader []byte `json:"response_header"`
	ResponseBody   []byte `json:"response_body"`
	CPEs           []*CPE `json:"cp_es"`
}

type CPE struct {
	Part     string `yaml:"part,omitempty" json:"part"`
	Vendor   string `yaml:"vendor,omitempty" json:"vendor"`
	Product  string `yaml:"product,omitempty" json:"product"`
	Version  string `yaml:"version,omitempty" json:"version"`
	Update   string `yaml:"update,omitempty" json:"update"`
	Edition  string `yaml:"edition,omitempty" json:"edition"`
	Language string `yaml:"language,omitempty" json:"language"`
}
