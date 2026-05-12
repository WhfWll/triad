package application

type BasReceivHeartBeat struct {
	Ip string `json:"ip"`
}

type BasReceivResult struct {
	Ip  string   `json:"ip"`
	Md5 []string `json:"md5"`
}

type VulEvidenceMapSet struct {
	Match string `json:"match"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
