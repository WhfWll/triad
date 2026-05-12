package enums

const (
	SessionStatusSucc  = 1 //在线
	SessionStatusFail  = 2 //离线
	CaptureTypeWindows = 1 //windows
	CaptureTypeCentos  = 2 //centos
)

// RemoteSessionMap  远程会话
var RemoteSessionMap = map[int]string{
	VulScriptEvidenceTypeCommandExec: "远程控制",
}

// SessionStatusMap  返回状态
var SessionStatusMap = map[int]string{
	SessionStatusSucc: "在线",
	SessionStatusFail: "离线",
}

const (
	WebShellTypeBehinder = "behinder" //webshell
	WebShellTypeGodzilla = "godzilla" //webshell
)
