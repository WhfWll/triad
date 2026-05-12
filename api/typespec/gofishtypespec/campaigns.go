package gofishtypespec

// CampaignIDReq 通用按 ID 查询/操作
type CampaignIDReq struct {
	ID int64 `json:"id" form:"id" common:"活动id" v:"required|min:1#活动id最小值为1"`
}

// CreateCampaignReq 创建活动请求
type CreateCampaignReq struct {
	Name     string `json:"name"`
	Template struct {
		Name string `json:"name"`
	} `json:"template"`
	Url  string `json:"url"`
	Page struct {
		Name string `json:"name"`
	} `json:"page"`
	Smtp struct {
		Name string `json:"name"`
	} `json:"smtp"`
	LaunchDate string      `json:"launch_date"`
	SendByDate interface{} `json:"send_by_date"`
	Groups     []struct {
		Name string `json:"name"`
	} `json:"groups"`
}

// UpdateCampaignReq 更新活动请求（与返回结构接近，按需字段）
type UpdateCampaignReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type DeleteByIDReq struct {
	ID int64 `json:"id" form:"id" v:"required|min:1#id最小值为1"`
}

type CampaignListReq struct {
	ActiveStatus string `json:"activeStatus" form:"activeStatus" common:"activeStatus"`
	GetListInfoReq
}
