package typespec

type GraphDataNodes struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	SymbolSize string `json:"symbolSize"`
}

type GraphDataLinks struct {
	Source int `json:"source"`
	Target int `json:"target"`
}
