package gophish

import (
	"fmt"
	"net/http"
	"smart/tools/time"
)

// LandingPagesAPI 处理着陆页相关的 API 调用
type LandingPagesAPI struct {
	client *Client
}

// LandingPage 表示着陆页
type LandingPage struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	HTML               string `json:"html"`
	CaptureCredentials bool   `json:"capture_credentials"`
	CapturePasswords   bool   `json:"capture_passwords"`
	ModifiedDate       string `json:"modified_date,omitempty"`
	RedirectURL        string `json:"redirect_url"`
}

// CreateLandingPageRequest 创建着陆页的请求
type CreateLandingPageRequest struct {
	Name               string `json:"name"`
	HTML               string `json:"html"`
	CaptureCredentials bool   `json:"capture_credentials"`
	CapturePasswords   bool   `json:"capture_passwords"`
	RedirectURL        string `json:"redirect_url"`
}

// ImportSiteRequest 导入网站的请求
type ImportSiteRequest struct {
	URL              string `json:"url"`
	IncludeResources bool   `json:"include_resources"`
}

// ImportSiteResponse 导入网站的响应
type ImportSiteResponse struct {
	HTML string `json:"html"`
}

// GetAll 获取所有着陆页
func (api *LandingPagesAPI) GetAll() ([]LandingPage, error) {
	resp, err := api.client.Get("/api/pages/")
	if err != nil {
		return nil, err
	}

	var pages []LandingPage
	if err := api.client.ParseResponse(resp, &pages); err != nil {
		return nil, err
	}

	for idx, group := range pages {
		pages[idx].ModifiedDate, _ = time.UTCToBeijingTime(group.ModifiedDate)
	}
	return pages, nil
}

// GetByID 根据 ID 获取着陆页
func (api *LandingPagesAPI) GetByID(id int64) (*LandingPage, error) {
	endpoint := fmt.Sprintf("/api/pages/%d", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var page LandingPage
	if err := api.client.ParseResponse(resp, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// Create 创建新的着陆页
func (api *LandingPagesAPI) Create(req *CreateLandingPageRequest) (*LandingPage, error) {
	resp, err := api.client.Post("/api/pages/", req)
	if err != nil {
		return nil, err
	}

	var page LandingPage
	if err := api.client.ParseResponse(resp, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// Update 更新着陆页
func (api *LandingPagesAPI) Update(id int64, page *LandingPage) (*LandingPage, error) {
	endpoint := fmt.Sprintf("/api/pages/%d", id)
	resp, err := api.client.Put(endpoint, page)
	if err != nil {
		return nil, err
	}

	var updatedPage LandingPage
	if err := api.client.ParseResponse(resp, &updatedPage); err != nil {
		return nil, err
	}

	return &updatedPage, nil
}

// Delete 删除着陆页
func (api *LandingPagesAPI) Delete(id int64) error {
	endpoint := fmt.Sprintf("/api/pages/%d", id)
	resp, err := api.client.Delete(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err := api.client.ParseResponse(resp, &apiResp); err != nil {
			return fmt.Errorf("删除失败，状态码: %d", resp.StatusCode)
		}
		return fmt.Errorf("删除失败: %s", apiResp.Message)
	}

	return nil
}

// ImportSite 导入网站作为着陆页
func (api *LandingPagesAPI) ImportSite(req *ImportSiteRequest) (*ImportSiteResponse, error) {
	resp, err := api.client.Post("/api/import/site", req)
	if err != nil {
		return nil, err
	}

	var response ImportSiteResponse
	if err := api.client.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
