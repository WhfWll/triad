package enums

const SWPrvKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDIVx4xd/pxtcDoScjigT4vnRk6LBZL3gdzaHExyyHkVNA8gZV4
IDqCTjUDB97vVJoIArXupGeHHNMfUnJvT5jen3J8TC92RCMuSay1Bfu+3bh+SbCH
W2EPSxzdSDb57yn2d3AKtOFURkLa8xDAheZqpK2a5MIkkl6BVhsaOilqewIDAQAB
AoGAZrSy1+2ISU9CbFOXVvenJ8XELxx2+cID09iRX1OiNmp8ruhH9mOfWzo41yrr
0YpvxPeOyZ8jLBNM8NvVqtcFqQLYd4PYQMNUPfK6nkVs0LXZfdifqU6Zxaj+DwJa
SwMYWCt1wLCBOw4lgmw+7lkjMo6z2Egts90kQx4CLs/pZIECQQDQbsNaDJd5iG62
A7Ln6FpZ06BZGM+fNDP7b0ajjn4z7w65qvZmJVrzOZ7/laqX36VoNAGXZ4Mnznpl
FnKwYOphAkEA9g+TXTULuBymryG9QInAwHeui68IrsG1ktqPCSooJm0Xxq9sbshS
JFrN4bOp3pfBs5hfxlYqSVsZn1erl1BaWwJBAK/F+PCj1cokCFlu3R09kZRXJ857
Yfw8penQeZ3MuRlK7Pwe9RRHRGABo9ieevMBJBiwYvcv0CdttUIyoB2mXYECQQDH
hBGlfiGmg/TUBLOD5S6Z4XFyadbMfN1R2k4ozDoKDmM9A3kUyvFv8QEHhbqzdrHl
giQGmk9nc6ru/RNxegIJAkA6qYCTEGIWpHLO+Nf9AVgiEMjYRpJ4811YYIEYPi0i
FzHmmCc8ornj4/GXdacNYw6gB2xJ1qh6/BV2QZGxpCL0
-----END RSA PRIVATE KEY-----
`

const SWPubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDIVx4xd/pxtcDoScjigT4vnRk6
LBZL3gdzaHExyyHkVNA8gZV4IDqCTjUDB97vVJoIArXupGeHHNMfUnJvT5jen3J8
TC92RCMuSay1Bfu+3bh+SbCHW2EPSxzdSDb57yn2d3AKtOFURkLa8xDAheZqpK2a
5MIkkl6BVhsaOilqewIDAQAB
-----END PUBLIC KEY-----
`
const (
	ProductAuthStateMapSetKey  = "productAuthState"
	ProductAuthInfoMapSetKey   = "productAuthInfo"
	ProductAuthRecordMapSetKey = "productAuthRecord"
	ProductIDMapSetKey         = "productID"
	ProductAuthStateSuccess    = "1"
	ProductAuthStateFailed     = "0"
	// ProductSoftwareDisplayVersion 系统授权页展示的软件版本号
	ProductSoftwareDisplayVersion = "V0.1.1.260603"
)

// AuthStatusUri 需要校验授权状态的uri
// 当系统未授权时，访问以下接口将被拦截并提示"系统未授权"
// 注意: 授权系统自身的接口(/smart/system/authsave, /smart/system/authinfo)不在此列，
//
//	以便用户可以在未授权状态下完成授权操作
var AuthStatusUri = map[string]bool{
	// 渗透任务
	"/smart/task/save":              true,
	"/smart/task/apisave":           true,
	"/smart/task/copy":              true,
	"/smart/task/addtarget":         true,
	"/smart/task/addattackface":     true,
	"/smart/task/addvul":            true,
	"/smart/task/flowtaskadd":       true,
	"/smart/task/vultest":           true,
	"/smart/task/vulverify":         true,
	"/smart/task/asyncvulverify":    true,
	// 任务组
	"/smart/taskgroup/create":       true,
	"/smart/taskgroup/groupbind":    true,
	// 场景管理
	"/smart/scene/save":             true,
	"/smart/scene/copy":             true,
	// 安全配置核查
	"/smart/baseline/check":         true,
	"/smart/baseline/check/batch":   true,
	// 数据安全
	"/smart/datasec/db/run":         true,
	"/smart/datasec/sensitive/run":  true,
	"/smart/datasec/task/rerun":     true,
	// 应用安全
	"/smart/appsec/dynamic/save":    true,
	"/smart/appsec/appspecific/run": true,
	// 漏洞扫描
	"/smart/vulscan/tasksave":       true,
	// 逻辑漏洞
	"/smart/logic/taskcreate":       true,
	// BAS
	"/smart/bas/taskcreate":         true,
	// 第三方工具
	"/smart/tripartite/xraysave":        true,
	"/smart/tripartite/burpsuitesave":   true,
	"/smart/tripartite/wificreate":      true,
	// 安全报告
	"/smart/security/report/generate": true,
	// 资产中心
	"/smart/asset/add":              true,
	"/smart/asset/import":           true,
}

const (
	ProductAuthInfoNewMapSetKey    = "productAuthInfoNewMapSet"
	ProductAuthInfoNewMapSetInfo   = "新产品授权信息"
	ProductAuthRecordNewMapSetKey  = "productAuthRecordNewMapSet"
	ProductAuthRecordNewMapSetInfo = "新产品授权记录"
)

type auth struct {
}

var AuthEnum auth

const (
	AuthStatusNotAuthed = iota + 1
	AuthStatusAuthed
	AuthStatusAuthOutOfDate
	AuthStatusAuthError
)

var AuthStatusNameMap = map[int]string{
	AuthStatusNotAuthed:     "未授权",
	AuthStatusAuthed:        "已授权",
	AuthStatusAuthOutOfDate: "授权已过期",
	AuthStatusAuthError:     "授权信息错误",
}

func (a *auth) GetAuthStatusName(status int) string {
	if desc, ok := AuthStatusNameMap[status]; ok {
		return desc
	}
	return "未知授权状态"
}

const (
	AuthInfoDefault        = "base"
	AuthInfoBaselineCheck  = "baseline"
	AuthInfoWebsiteMonitor = "website"
)

const AuthDate = "060102"

// AuthInfo 数据库中存储的激活数据
type AuthInfo struct {
	ProductName string            `json:"productName"`
	ProductID   string            `json:"productID"`   //
	AuthCode    string            `json:"authCode"`    // 产品序列号
	AuthTime    map[string]string `json:"authTime"`    // 授权日期
	AuthExpTime map[string]string `json:"authExpTime"` // 授权过期时间
	AuthDays    map[string]int    `json:"authDays"`    // 授权时长
	LeftDays    map[string]int    `json:"leftDays"`    // 剩余时间
}

// AuthData 用户用于激活的激活数据
type AuthData struct {
	Machine  string `json:"machine"`
	Time     string `json:"time"`
	Days     int    `json:"days"`
	Base     bool   `json:"base"`
	Baseline bool   `json:"baseline"`
	Website  bool   `json:"website"`
}
