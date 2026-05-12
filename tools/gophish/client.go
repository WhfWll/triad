package gophish

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 封装了 Gophish API 的 HTTP 客户端
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient 创建一个新的 Gophish 客户端
func NewClient(baseURL, apiKey string) *Client {
	// 创建自定义的 Transport 来跳过 SSL 证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

// APIResponse 表示 Gophish API 的通用响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// doRequest 执行 HTTP 请求的通用方法
func (c *Client) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.BaseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %v", err)
	}

	return resp, nil
}

// Get 执行 GET 请求
func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.doRequest("GET", endpoint, nil)
}

// Post 执行 POST 请求
func (c *Client) Post(endpoint string, body interface{}) (*http.Response, error) {
	return c.doRequest("POST", endpoint, body)
}

// Put 执行 PUT 请求
func (c *Client) Put(endpoint string, body interface{}) (*http.Response, error) {
	return c.doRequest("PUT", endpoint, body)
}

// Delete 执行 DELETE 请求
func (c *Client) Delete(endpoint string) (*http.Response, error) {
	return c.doRequest("DELETE", endpoint, nil)
}

// ParseResponse 解析 API 响应
func (c *Client) ParseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}
	if resp.StatusCode >= 400 {
		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return fmt.Errorf("解析错误响应失败: %v", err)
		}
		return fmt.Errorf("API 错误: %s", apiResp.Message)
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("解析响应失败: %v", err)
		}
	}

	return nil
}
