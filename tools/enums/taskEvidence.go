package enums

// TaskEvidenceMap  任务证据
var TaskEvidenceMap = map[int]string{
	VulScriptEvidenceTypeWeakPass: "登录凭证",
	VulScriptEvidenceTypeInfoLeak: "敏感数据",
	VulScriptEvidenceTypeFileLeak: "敏感文件",
	VulScriptEvidenceTypeData:     "数据库",
}

// RemoteSessionEnumMap 会话map
var RemoteSessionEnumMap = map[int]string{
	VulScriptEvidenceTypeReverseShell:  "反弹shell",
	VulScriptEvidenceTypeRemoteControl: "远控",
	VulScriptEvidenceTypeWebShell:      "webshell",
}

// 横向
var HengxiangMap = map[int]string{
	VulScriptEvidenceTypeHengxiang: "横向",
}

// TaskEvidenceCaptureInfoMap  任务证据抓取信息
var TaskEvidenceCaptureInfoMap = map[int]string{
	1: "用户信息",
	2: "网卡信息",
	3: "系统版本",
	4: "进程信息",
	5: "端口信息",
	6: "环境变量",
	7: "屏幕截图",
}

// TaskEvidenceCaptureInfoCentosExecMap  centos任务证据抓取信息执行
var TaskEvidenceCaptureInfoCentosExecMap = map[int]string{
	1: "cat /etc/passwd",
	2: "ifconfig",
	3: "cat /etc/os-release",
	4: "ps aux",
	5: "netstat -tunpl",
	6: "printenv",
	7: "screenshot",
	8: "exit",
}

// TaskEvidenceCaptureInfoWindowsExecMap  windows任务证据抓取信息执行
var TaskEvidenceCaptureInfoWindowsExecMap = map[int]string{
	1: "getuid",
	2: "ipconfig",
	3: "sysinfo",
	4: "ps",
	5: "netstat -ano",
	6: "getenv path",
	7: "screenshot",
	8: "exit",
}
