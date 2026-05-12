package gophish

import (
	"fmt"
	"net/http"
)

// UserManagementAPI 处理用户管理相关的 API 调用
type UserManagementAPI struct {
	client *Client
}

// UserAccount 表示用户账户（参考官方文档）
type UserAccount struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Role         Role   `json:"role"`
	ModifiedDate string `json:"modified_date"`
	// 以下字段并非文档固定返回，但保留用于创建/更新时传参
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`

	// 创建后的返回
	ApiKey string `json:"api_key"`
}

// Role 表示账户角色（参考官方文档）
type Role struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// GetRoleSlug 获取角色 slug
func (u *UserAccount) GetRoleSlug() string { return u.Role.Slug }

// GetApiKey 获取 ApiKey
func (u *UserAccount) GetApiKey() string { return u.ApiKey }

// CreateUserAccountRequest 创建用户账户的请求
type CreateUserAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

// UpdateUserAccountRequest 更新用户账户的请求
type UpdateUserAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// GetAll 获取所有用户账户
func (api *UserManagementAPI) GetAll() ([]UserAccount, error) {
	resp, err := api.client.Get("/api/users/")
	if err != nil {
		return nil, err
	}

	var accounts []UserAccount
	if err := api.client.ParseResponse(resp, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// GetByID 根据 ID 获取用户账户
func (api *UserManagementAPI) GetByID(id int64) (*UserAccount, error) {
	endpoint := fmt.Sprintf("/api/users/%d", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var account UserAccount
	if err := api.client.ParseResponse(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// Create 创建新用户账户
func (api *UserManagementAPI) Create(req *CreateUserAccountRequest) (*UserAccount, error) {
	resp, err := api.client.Post("/api/users/", req)
	if err != nil {
		return nil, err
	}

	var account UserAccount
	if err := api.client.ParseResponse(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// Update 更新用户账户
func (api *UserManagementAPI) Update(id int64, req *UpdateUserAccountRequest) (*UserAccount, error) {
	endpoint := fmt.Sprintf("/api/users/%d", id)
	resp, err := api.client.Put(endpoint, req)
	if err != nil {
		return nil, err
	}

	var account UserAccount
	if err := api.client.ParseResponse(resp, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// Delete 删除用户账户
func (api *UserManagementAPI) Delete(id int64) error {
	endpoint := fmt.Sprintf("/api/users/%d", id)
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

// ChangePassword 修改用户密码（按官方文档通过 PUT /api/users/:id 设置 password 字段）
func (api *UserManagementAPI) ChangePassword(id int64, newPassword string) error {
	endpoint := fmt.Sprintf("/api/users/%d", id)
	req := map[string]string{"password": newPassword}

	resp, err := api.client.Put(endpoint, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp APIResponse
		if err := api.client.ParseResponse(resp, &apiResp); err != nil {
			return fmt.Errorf("修改密码失败，状态码: %d", resp.StatusCode)
		}
		return fmt.Errorf("修改密码失败: %s", apiResp.Message)
	}

	return nil
}
