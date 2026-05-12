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
)

// AuthStatusUri 需要校验授权状态的uri
var AuthStatusUri = map[string]bool{
	"/secops/auth/save": true,
	"/secops/auth/info": true,
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
