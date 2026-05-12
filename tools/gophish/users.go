package gophish

import (
	"fmt"
	"net/http"
)

// UsersAPI 处理用户相关的 API 调用
type UsersAPI struct {
	client *Client
}

// User 表示用户（实际上是目标用户，与 Target 结构相同）
type User struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Position  string `json:"position"`
}

// CreateUserRequest 创建用户的请求
type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Position  string `json:"position"`
}

// GetAll 获取所有用户
func (api *UsersAPI) GetAll() ([]User, error) {
	resp, err := api.client.Get("/api/users/")
	if err != nil {
		return nil, err
	}

	var users []User
	if err := api.client.ParseResponse(resp, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID 根据 ID 获取用户
func (api *UsersAPI) GetByID(id int64) (*User, error) {
	endpoint := fmt.Sprintf("/api/users/%d", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var user User
	if err := api.client.ParseResponse(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Create 创建新用户
func (api *UsersAPI) Create(req *CreateUserRequest) (*User, error) {
	resp, err := api.client.Post("/api/users/", req)
	if err != nil {
		return nil, err
	}

	var user User
	if err := api.client.ParseResponse(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Update 更新用户
func (api *UsersAPI) Update(id int64, user *User) (*User, error) {
	endpoint := fmt.Sprintf("/api/users/%d", id)
	resp, err := api.client.Put(endpoint, user)
	if err != nil {
		return nil, err
	}

	var updatedUser User
	if err := api.client.ParseResponse(resp, &updatedUser); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// Delete 删除用户
func (api *UsersAPI) Delete(id int64) error {
	endpoint := fmt.Sprintf("/api/users/%d", id)
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
