package gophish

import (
	"fmt"
	"net/http"
	"smart/tools/time"
)

// SendingProfilesAPI 处理发送配置文件相关的 API 调用
type SendingProfilesAPI struct {
	client *Client
}

// SendingProfile 表示发送配置文件
type SendingProfile struct {
	ID               int64    `json:"id"`
	Name             string   `json:"name"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	Host             string   `json:"host"`
	InterfaceType    string   `json:"interface_type"`
	FromAddress      string   `json:"from_address"`
	IgnoreCertErrors bool     `json:"ignore_cert_errors"`
	ModifiedDate     string   `json:"modified_date"`
	Headers          []Header `json:"headers"`
}

// Header 表示邮件头
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateSendingProfileRequest 创建发送配置文件的请求
type CreateSendingProfileRequest struct {
	Name             string   `json:"name"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	Host             string   `json:"host"`
	InterfaceType    string   `json:"interface_type"`
	FromAddress      string   `json:"from_address"`
	IgnoreCertErrors bool     `json:"ignore_cert_errors"`
	Headers          []Header `json:"headers,omitempty"`
}

type SendTestEmailReq struct {
	Template  NameInfo `json:"template"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Position  string   `json:"position"`
	Url       string   `json:"url"`
	Smtp      Smtp     `json:"smtp"`
	Page      NameInfo `json:"page"`
}

type Smtp struct {
	Name             string   `json:"name"`
	FromAddress      string   `json:"from_address"`
	Host             string   `json:"host"`
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	IgnoreCertErrors bool     `json:"ignore_cert_errors"`
	Headers          []Header `json:"headers"`
}

// GetAll 获取所有发送配置文件
func (api *SendingProfilesAPI) GetAll() ([]SendingProfile, error) {
	resp, err := api.client.Get("/api/smtp/")
	if err != nil {
		return nil, err
	}

	var profiles []SendingProfile
	if err := api.client.ParseResponse(resp, &profiles); err != nil {
		return nil, err
	}

	for idx, group := range profiles {
		profiles[idx].ModifiedDate, _ = time.UTCToBeijingTime(group.ModifiedDate)
	}

	return profiles, nil
}

// GetByID 根据 ID 获取发送配置文件
func (api *SendingProfilesAPI) GetByID(id int64) (*SendingProfile, error) {
	endpoint := fmt.Sprintf("/api/smtp/%d", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var profile SendingProfile
	if err := api.client.ParseResponse(resp, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// Create 创建新的发送配置文件
func (api *SendingProfilesAPI) Create(req *CreateSendingProfileRequest) (*SendingProfile, error) {
	resp, err := api.client.Post("/api/smtp/", req)
	if err != nil {
		return nil, err
	}

	var profile SendingProfile
	if err := api.client.ParseResponse(resp, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// Update 更新发送配置文件
func (api *SendingProfilesAPI) Update(id int64, profile *SendingProfile) (*SendingProfile, error) {
	endpoint := fmt.Sprintf("/api/smtp/%d", id)
	resp, err := api.client.Put(endpoint, profile)
	if err != nil {
		return nil, err
	}

	var updatedProfile SendingProfile
	if err := api.client.ParseResponse(resp, &updatedProfile); err != nil {
		return nil, err
	}

	return &updatedProfile, nil
}

// Delete 删除发送配置文件
func (api *SendingProfilesAPI) Delete(id int64) error {
	endpoint := fmt.Sprintf("/api/smtp/%d", id)
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

// SendTestEmail 发送测试邮件
func (api *SendingProfilesAPI) SendTestEmail(req *SendTestEmailReq) error {
	endpoint := fmt.Sprintf("/api/util/send_test_email")
	resp, err := api.client.Post(endpoint, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err := api.client.ParseResponse(resp, &apiResp); err != nil {
			return fmt.Errorf("发送失败，状态码: %d", resp.StatusCode)
		}
		return fmt.Errorf("发送失败: %s", apiResp.Message)
	}

	return nil
}
