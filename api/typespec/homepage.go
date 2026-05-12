package typespec

type TaskInfoStatReq struct {
	Mode int `json:"mode" form:"mode" binding:"required"` //统计模式 1-周 2-月 3-年
}

type TaskInfoStatRes struct {
	TaskCount           int      `json:"taskCount"`
	LatestWeekTaskCount int      `json:"latestWeekTaskCount"`
	Date                []string `json:"date"`
	Count               []int    `json:"count"`
}

type VulEvidenceStatRes struct {
	FileLeakCount         int `json:"fileLeakCount"`
	InfoLeakCount         int `json:"infoLeakCount"`
	DbCount               int `json:"dbCount"`
	LoginCredentialsCount int `json:"loginCredentialsCount"`
	RemoteControlCount    int `json:"remoteControlCount"`
}

type TargetRiskStatRes struct {
	HighCount       int `json:"highCount"`
	MediumCount     int `json:"mediumCount"`
	LowCount        int `json:"lowCount"`
	SafeCount       int `json:"safeCount"`
	TargetRiskCount int `json:"targetRiskCount"`
}

type TaskVulRiskStatRes struct {
	FatalCount       int `json:"fatalCount"`
	HighCount        int `json:"highCount"`
	MediumCount      int `json:"mediumCount"`
	LowCount         int `json:"lowCount"`
	InfoCount        int `json:"infoCount"`
	TaskVulRiskCount int `json:"taskVulRiskCount"`
}

type ToolInfoStatRes struct {
	VulCount        int `json:"vulCount"`
	FingerCount     int `json:"fingerCount"`
	TaskSceneCount  int `json:"taskSceneCount"`
	DictionaryCount int `json:"dictionaryCount"` // 字典库数量
}

type VulTypeStatReq struct {
	Mode int `json:"mode" form:"mode" binding:"required"` //统计模式 1-周 2-月 3-年
}

type VulTypeStatRes struct {
	VulType  []string `json:"vulType"`
	VulCount []int    `json:"vulCount"`
}

type VulFindTrendStatReq struct {
	Mode int `json:"mode" form:"mode" binding:"required"` //统计模式 1-周 2-月 3-年
}

type VulFindTrendStatRes struct {
	Date  []string `json:"date"`
	Count []int    `json:"count"`
}

type MessageStatRes struct {
	List []MessageStatResItem `json:"list"`
}

type MessageStatResItem struct {
	Id         int    `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	Type       int    `json:"type"`
	TypeEnum   string `json:"TypeEnum"`
	Status     int    `json:"status"`
	StatusEnum string `json:"StatusEnum"`
	UserId     int    `json:"userId"`
}
