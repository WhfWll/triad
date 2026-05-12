package gofishtypespec

type PageIDReq struct {
	ID int64 `json:"id" form:"id" common:"着陆页id" v:"required|min:1#着陆页id最小值为1"`
}

type CreatePageReq struct {
	Name               string `json:"name"`
	HTML               string `json:"html"`
	CaptureCredentials bool   `json:"capture_credentials"`
	CapturePasswords   bool   `json:"capture_passwords"`
	RedirectURL        string `json:"redirect_url"`
}

type UpdatePageReq struct {
	ID                 int64  `json:"id" form:"id" common:"着陆页id" v:"required|min:1#着陆页id最小值为1"`
	Name               string `json:"name"`
	HTML               string `json:"html"`
	CaptureCredentials bool   `json:"capture_credentials"`
	CapturePasswords   bool   `json:"capture_passwords"`
	RedirectURL        string `json:"redirect_url"`
}

type ImportSiteReq struct {
	URL              string `json:"url"`
	IncludeResources bool   `json:"include_resources"`
}
