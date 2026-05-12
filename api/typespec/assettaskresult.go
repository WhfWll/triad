package typespec

// AssetTaskResultJsonRes 资产任务返回json结构体
type AssetTaskResultJsonRes struct {
	Banner    string `json:"banner"`    //
	Cert      string `json:"cert"`      // 证书
	Component string `json:"component"` // 组件
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Protocol  string `json:"protocol"` //协议类型
	Response  string `json:"response"` // 端口返回
	Service   string `json:"service"`  // 服务
	Title     string `json:"title"`
	HtmlTitle string `json:"htmlTitle"`
}
