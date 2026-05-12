// Package typespec
// @Author bcy2007  2025/12/23 11:05
package typespec

type ReverseServerStartReq struct {
	Host string `json:"host" form:"host"`
	Port int    `json:"port" form:"port"`
}

type ReverseServerStatusResp struct {
	Status     bool   `json:"status"`
	ReverseUrl string `json:"reverseUrl"`
	ErrInfo    string `json:"errInfo"`
}

type ReverseServerMessageReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type ReverseServerMessageResp struct {
	Total int                        `json:"total"`
	List  []ReverseServerMessageItem `json:"list"`
}

type ReverseServerMessageItem struct {
	ReverseType string `json:"reverseType"`
	RemoteAddr  string `json:"remoteAddr"`
	Token       string `json:"token"`
	Response    string `json:"response"`
}
