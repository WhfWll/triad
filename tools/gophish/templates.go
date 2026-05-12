package gophish

import (
	"fmt"
	"net/http"
	"smart/tools/time"
)

// TemplatesAPI 处理模板相关的 API 调用
type TemplatesAPI struct {
	client *Client
}

// Template 表示邮件模板
type Template struct {
	ID             int64        `json:"id"`
	Name           string       `json:"name"`
	Subject        string       `json:"subject"`
	Text           string       `json:"text"`
	HTML           string       `json:"html"`
	ModifiedDate   string       `json:"modified_date,omitempty"`
	EnvelopeSender string       `json:"envelope_sender"`
	Attachments    []Attachment `json:"attachments,omitempty"`
}

// Attachment 表示邮件附件
type Attachment struct {
	Content string `json:"content"` // base64 编码的内容
	Type    string `json:"type"`
	Name    string `json:"name"`
}

// CreateTemplateRequest 创建模板的请求
type CreateTemplateRequest struct {
	Name           string       `json:"name"`
	Subject        string       `json:"subject"`
	Text           string       `json:"text"`
	HTML           string       `json:"html"`
	EnvelopeSender string       `json:"envelope_sender"`
	Attachments    []Attachment `json:"attachments,omitempty"`
}

// ImportEmailRequest 导入邮件的请求
type ImportEmailRequest struct {
	Content      string `json:"content"`
	ConvertLinks bool   `json:"convert_links"`
}

// ImportEmailResponse 导入邮件的响应
type ImportEmailResponse struct {
	Text    string `json:"text"`
	HTML    string `json:"html"`
	Subject string `json:"subject"`
}

// GetAll 获取所有模板
func (api *TemplatesAPI) GetAll() ([]Template, error) {
	resp, err := api.client.Get("/api/templates/")
	if err != nil {
		return nil, err
	}

	var templates []Template
	if err := api.client.ParseResponse(resp, &templates); err != nil {
		return nil, err
	}

	for idx, group := range templates {
		templates[idx].ModifiedDate, _ = time.UTCToBeijingTime(group.ModifiedDate)
	}

	return templates, nil
}

// GetByID 根据 ID 获取模板
func (api *TemplatesAPI) GetByID(id int64) (*Template, error) {
	endpoint := fmt.Sprintf("/api/templates/%d", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var template Template
	if err := api.client.ParseResponse(resp, &template); err != nil {
		return nil, err
	}

	return &template, nil
}

// Create 创建新的模板
func (api *TemplatesAPI) Create(req *CreateTemplateRequest) (*Template, error) {
	resp, err := api.client.Post("/api/templates/", req)
	if err != nil {
		return nil, err
	}

	var template Template
	if err := api.client.ParseResponse(resp, &template); err != nil {
		return nil, err
	}

	return &template, nil
}

// Update 更新模板
func (api *TemplatesAPI) Update(id int64, template *Template) (*Template, error) {
	endpoint := fmt.Sprintf("/api/templates/%d", id)
	resp, err := api.client.Put(endpoint, template)
	if err != nil {
		return nil, err
	}

	var updatedTemplate Template
	if err := api.client.ParseResponse(resp, &updatedTemplate); err != nil {
		return nil, err
	}

	return &updatedTemplate, nil
}

// Delete 删除模板
func (api *TemplatesAPI) Delete(id int64) error {
	endpoint := fmt.Sprintf("/api/templates/%d", id)
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

// ImportEmail 导入邮件作为模板
func (api *TemplatesAPI) ImportEmail(req *ImportEmailRequest) (*ImportEmailResponse, error) {
	resp, err := api.client.Post("/api/import/email", req)
	if err != nil {
		return nil, err
	}

	var response ImportEmailResponse
	if err := api.client.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
