package enums

const (
	TargetStatusToTrigger int = 1
	TargetStatusToBegin   int = 2
	TargetStatusRunning   int = 3
	TargetStatusFinish    int = 4
	TargetStatusPausing   int = 5

	TargetRiskHigh       int = 1 //高危
	TargetRiskMid        int = 2 //中危
	TargetRiskLow        int = 3 //低危
	TargetRiskLowNoFound int = 4 //未发现

	TargetIsAliveN = 1 // 目标未存活
	TargetIsAliveY = 2 // 目标存活

	TargetIsRemoteSessionN = 1 // 目标远程会话 无
	TargetIsRemoteSessionY = 2 // 目标远程会话 有

	TaskOperateCreate  string = "create"  //任务创建消息
	TaskOperateStop    string = "stop"    //任务结束消息
	TaskOperateResume  string = "resume"  //任务恢复消息
	TaskOperatePause   string = "pause"   //任务暂停消息
	TaskOperateAddInfo string = "addInfo" //添加信息
	TaskOperateTestVul string = "testvul" //待测漏洞测试

	//部分用到特殊信息收集脚本名称
	ScriptNamePortScan              string = "port_scan"                       //端口扫描脚本
	ScriptNameFingerPrint           string = "fingerprint"                     //指纹识别脚本
	ScriptNameWebScanner            string = "web_scanner"                     //web扫描器脚本
	ScriptNameHttpBaseInfoGet       string = "http_base_info_get"              //网站基本信息获取脚本
	ScriptNameCrawlerx              string = "crawlerx"                        //爬虫脚本
	ScriptNameWebDirPathScan        string = "web_dir_path_scan"               //web路径爆破脚本
	ScriptNameDeviceOsInfo          string = "device_os_info"                  //操作系统识别脚本
	ScriptNameWafDetect             string = "waf_detect"                      //waf识别脚本
	ScriptNameSubdomain             string = "subdomain_brute"                 //子域名爆破脚本
	ScriptNameCdnDetect             string = "cdn_detect"                      //cdn探测脚本
	ScriptNameWhois                 string = "whois"                           //whois脚本
	ScriptNameSecondDirBrute        string = "second_dir_brute"                // 二级目录爆破
	ScriptNameLogicBeyondPermission string = "over_permission_check"           //越权检测
	ScriptNameLogicTraverseTesting  string = "sensitive_info_traversal"        //信息遍历测试
	ScriptNameLogicUnAuthAccess     string = "logic_unauthorized_access_check" //信息遍历测试
	ScriptNameServiceBruteCheck     string = "service_bruteforce_check"        //服务爆破

	TargetIsScoreToBegin = 1 //评分未开始
	TargetIsScoreRunning = 2 //评分进行中
	TargetIsScoreFinish  = 3 //评分完成
	TargetIsScoreNoneed  = 4 //无需评分

	TargetIsScoreSwitchOn  = 1 //可利用评分打开
	TargetIsScoreSwitchOff = 2 //可利用评分关闭

)

const (
	TaskCheckTargetTypeHost     = 1  //主机
	TaskCheckTargetTypeDatabase = 2  //数据库
	TaskCheckTargetTypeWeb      = 3  //web
	TaskCheckTargetTypeCloud    = 4  //云平台
	TaskCheckTargetTypeBigData  = 5  //大数据
	TaskCheckTargetTypeNet      = 6  //网络设备
	TaskCheckTargetTypeSafe     = 7  //安全设备
	TaskCheckTargetTypeVideo    = 8  //视频监控
	TaskCheckTargetTypeOffice   = 9  //办公自动化
	TaskCheckTargetTypeAp       = 10 //AP
	TaskCheckTargetTypeHc       = 11 //HC
)

func GetTargetStatusEnum() map[int]string {
	enum := map[int]string{
		TargetStatusToTrigger: "待触发",
		TargetStatusToBegin:   "待执行",
		TargetStatusRunning:   "运行中",
		TargetStatusFinish:    "已结束",
		TargetStatusPausing:   "暂停中",
	}
	return enum
}

func GetTargetStatus(status int) string {
	enum := GetTargetStatusEnum()
	if v, ok := enum[status]; ok {
		return v
	}
	return ""
}

func GetTargetRiskEnum() map[int]string {
	enum := map[int]string{
		TargetRiskHigh:       "高危",
		TargetRiskMid:        "中危",
		TargetRiskLow:        "低危",
		TargetRiskLowNoFound: "安全",
	}
	return enum
}

func GetTargetRisk(risk int) string {
	enum := GetTargetRiskEnum()
	if v, ok := enum[risk]; ok {
		return v
	}
	return ""
}

func GetTargetIsRemoteSessionEnum() map[int]string {
	enum := map[int]string{
		TargetIsRemoteSessionN: "无",
		TargetIsRemoteSessionY: "有",
	}
	return enum
}

func GetTargetIsRemoteSession(isRemoteSession int) string {
	enum := GetTargetIsRemoteSessionEnum()
	if v, ok := enum[isRemoteSession]; ok {
		return v
	}
	return ""
}

// TargetIsAliveEnum 目标存活状态枚举
func TargetIsAliveEnum() map[int]string {
	enum := map[int]string{
		TargetIsAliveN: "不存活",
		TargetIsAliveY: "存活",
	}
	return enum
}

// GetTargetIsAlive 目标存活状态的值
func GetTargetIsAlive(k int) string {
	enum := TargetIsAliveEnum()
	if v, ok := enum[k]; ok {
		return v
	}
	return ""
}
