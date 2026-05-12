package gophish

import (
	"fmt"
	"net/http"
	"smart/tools/time"
)

// CampaignsAPI 处理活动相关的 API 调用
type CampaignsAPI struct {
	client *Client
}

// Campaign 表示钓鱼活动
type Campaign struct {
	ID            int64          `json:"id"`
	Name          string         `json:"name"`
	CreatedDate   string         `json:"created_date"`
	LaunchDate    string         `json:"launch_date"`
	SendByDate    string         `json:"send_by_date"`
	CompletedDate string         `json:"completed_date"`
	Template      Template       `json:"template"`
	Page          LandingPage    `json:"page"`
	Status        string         `json:"status"`
	Results       []Results      `json:"results"`
	Timeline      []Timeline     `json:"timeline"`
	Smtp          SendingProfile `json:"smtp"`
	Url           string         `json:"url"`
}

type Results struct {
	Id           string `json:"id"`
	Status       string `json:"status"`
	Ip           string `json:"ip"`
	Latitude     int    `json:"latitude"`
	Longitude    int    `json:"longitude"`
	SendDate     string `json:"send_date"`
	Reported     bool   `json:"reported"`
	ModifiedDate string `json:"modified_date"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Position     string `json:"position"`
}

type ResultResp struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Status   string     `json:"status"`
	Results  []Results  `json:"results"`
	Timeline []Timeline `json:"timeline"`
}

// Result 表示活动结果
type Result struct {
	ID           string  `json:"id"`
	CampaignID   int64   `json:"campaign_id"`
	UserID       int64   `json:"user_id"`
	Email        string  `json:"email"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Status       string  `json:"status"`
	IP           string  `json:"ip"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	UserAgent    string  `json:"user_agent"`
	ModifiedDate string  `json:"modified_date"`
}

type CampaignsResp struct {
	Total     int            `json:"total"`
	Campaigns []CampaignList `json:"campaigns"`
}
type CampaignList struct {
	Id            int    `json:"id"`
	CreatedDate   string `json:"created_date"`
	LaunchDate    string `json:"launch_date"`
	SendByDate    string `json:"send_by_date"`
	CompletedDate string `json:"completed_date"`
	Status        string `json:"status"`
	Name          string `json:"name"`
	Stats         Stats  `json:"stats"`
}
type Stats struct {
	Total         int `json:"total"`
	Sent          int `json:"sent"`
	Opened        int `json:"opened"`
	Clicked       int `json:"clicked"`
	SubmittedData int `json:"submitted_data"`
	EmailReported int `json:"email_reported"`
	Error         int `json:"error"`
}

type Timeline struct {
	CampaignId int    `json:"campaign_id"`
	Email      string `json:"email"`
	Time       string `json:"time"`
	Message    string `json:"message"`
	Details    string `json:"details"`
}

// CreateCampaignRequest 创建活动的请求
type CreateCampaignRequest struct {
	Name       string      `json:"name"`
	Template   NameInfo    `json:"template"`
	Url        string      `json:"url"`
	Page       NameInfo    `json:"page"`
	Smtp       NameInfo    `json:"smtp"`
	LaunchDate string      `json:"launch_date"`
	SendByDate interface{} `json:"send_by_date"`
	Groups     []struct {
		Name string `json:"name"`
	} `json:"groups"`
}

type NameInfo struct {
	Name string `json:"name"`
}

// GetAll 获取所有活动
func (api *CampaignsAPI) GetAll() (CampaignsResp, error) {
	resp, err := api.client.Get("/api/campaigns/summary")
	if err != nil {
		return CampaignsResp{}, err
	}

	var campaigns CampaignsResp
	if err = api.client.ParseResponse(resp, &campaigns); err != nil {
		return CampaignsResp{}, err
	}

	for idx, group := range campaigns.Campaigns {
		campaigns.Campaigns[idx].CreatedDate, _ = time.UTCToBeijingTime(group.CreatedDate)
		campaigns.Campaigns[idx].LaunchDate, _ = time.UTCToBeijingTime(group.LaunchDate)
		campaigns.Campaigns[idx].SendByDate, _ = time.UTCToBeijingTime(group.SendByDate)
		campaigns.Campaigns[idx].CompletedDate, _ = time.UTCToBeijingTime(group.CompletedDate)
	}

	return campaigns, nil
}

// GetByID 根据 ID 获取活动
func (api *CampaignsAPI) GetByID(id int64) (*Campaign, error) {
	endpoint := fmt.Sprintf("/api/campaigns/%d", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var campaign Campaign
	if err = api.client.ParseResponse(resp, &campaign); err != nil {
		return nil, err
	}

	campaign.CreatedDate, _ = time.SmartConvertToBeijing(campaign.CreatedDate)
	campaign.LaunchDate, _ = time.UTCToBeijingTime(campaign.LaunchDate)
	campaign.SendByDate, _ = time.UTCToBeijingTime(campaign.SendByDate)
	campaign.CompletedDate, _ = time.UTCToBeijingTime(campaign.CompletedDate)

	return &campaign, nil
}

// Create 创建新活动
func (api *CampaignsAPI) Create(req *CreateCampaignRequest) (*Campaign, error) {
	resp, err := api.client.Post("/api/campaigns/", req)
	if err != nil {
		return nil, err
	}

	var campaign Campaign
	if err = api.client.ParseResponse(resp, &campaign); err != nil {
		return nil, err
	}

	return &campaign, nil
}

// Update 更新活动
func (api *CampaignsAPI) Update(id int64, campaign *Campaign) (*Campaign, error) {
	endpoint := fmt.Sprintf("/api/campaigns/%d", id)
	resp, err := api.client.Put(endpoint, campaign)
	if err != nil {
		return nil, err
	}

	var updatedCampaign Campaign
	if err = api.client.ParseResponse(resp, &updatedCampaign); err != nil {
		return nil, err
	}

	return &updatedCampaign, nil
}

// Delete 删除活动
func (api *CampaignsAPI) Delete(id int64) error {
	endpoint := fmt.Sprintf("/api/campaigns/%d", id)
	resp, err := api.client.Delete(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err = api.client.ParseResponse(resp, &apiResp); err != nil {
			return fmt.Errorf("删除失败，状态码: %d", resp.StatusCode)
		}
		return fmt.Errorf("删除失败: %s", apiResp.Message)
	}

	return nil
}

// Launch 启动活动
func (api *CampaignsAPI) Launch(id int64) error {
	endpoint := fmt.Sprintf("/api/campaigns/%d/launch", id)
	resp, err := api.client.Post(endpoint, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err = api.client.ParseResponse(resp, &apiResp); err != nil {
			return fmt.Errorf("启动失败，状态码: %d", resp.StatusCode)
		}
		return fmt.Errorf("启动失败: %s", apiResp.Message)
	}

	return nil
}

// Complete 完成活动
func (api *CampaignsAPI) Complete(id int64) error {
	endpoint := fmt.Sprintf("/api/campaigns/%d/complete", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err = api.client.ParseResponse(resp, &apiResp); err != nil {
			return fmt.Errorf("完成失败，状态码: %d", resp.StatusCode)
		}
		return fmt.Errorf("完成失败: %s", apiResp.Message)
	}

	return nil
}

// Results 活动结果
func (api *CampaignsAPI) Results(id int64) (interface{}, error) {
	endpoint := fmt.Sprintf("/api/campaigns/%d/results", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err = api.client.ParseResponse(resp, &apiResp); err != nil {
			return nil, fmt.Errorf("完成失败，状态码: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("完成失败: %s", apiResp.Message)
	}

	var results ResultResp
	if err = api.client.ParseResponse(resp, &results); err != nil {
		return nil, err
	}

	for idx, group := range results.Results {
		results.Results[idx].SendDate, _ = time.UTCToBeijingTime(group.SendDate)
		results.Results[idx].ModifiedDate, _ = time.UTCToBeijingTime(group.ModifiedDate)
	}

	for idx, group := range results.Timeline {
		results.Timeline[idx].Time, _ = time.UTCToBeijingTime(group.Time)
	}

	return results, nil
}
