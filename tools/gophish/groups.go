package gophish

import (
	"fmt"
	"net/http"
	"smart/tools/time"
)

// GroupsAPI 处理用户组相关的 API 调用
type GroupsAPI struct {
	client *Client
}

// Group 表示用户组（根据官方文档）
type Group struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Targets      []Target `json:"targets,omitempty"`
	ModifiedDate string   `json:"modified_date"`
	NumTargets   int64    `json:"num_targets,omitempty"`
}

// Target 表示组中的目标用户
type Target struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Position  string `json:"position"`
}

type GroupResponse struct {
	Groups []Group `json:"groups"`
	Total  int64   `json:"total"`
}

// CreateGroupRequest 创建用户组的请求
type CreateGroupRequest struct {
	Name    string   `json:"name"`
	Targets []Target `json:"targets"`
}

// GetAll 获取所有用户组
func (api *GroupsAPI) GetAll() (GroupResponse, error) {
	resp, err := api.client.Get("/api/groups/summary")
	if err != nil {
		return GroupResponse{}, err
	}

	var groups GroupResponse
	if err = api.client.ParseResponse(resp, &groups); err != nil {
		return GroupResponse{}, err
	}

	for idx, group := range groups.Groups {
		groups.Groups[idx].ModifiedDate, _ = time.UTCToBeijingTime(group.ModifiedDate)
	}

	return groups, nil
}

// GetByID 根据 ID 获取用户组
func (api *GroupsAPI) GetByID(id int64) (*Group, error) {
	endpoint := fmt.Sprintf("/api/groups/%d", id)
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var group Group
	if err := api.client.ParseResponse(resp, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// Create 创建新用户组
func (api *GroupsAPI) Create(req *CreateGroupRequest) (*Group, error) {
	resp, err := api.client.Post("/api/groups/", req)
	if err != nil {
		return nil, err
	}

	var group Group
	if err := api.client.ParseResponse(resp, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// Update 更新用户组
func (api *GroupsAPI) Update(id int64, group *Group) (*Group, error) {
	endpoint := fmt.Sprintf("/api/groups/%d", id)
	resp, err := api.client.Put(endpoint, group)
	if err != nil {
		return nil, err
	}

	var updatedGroup Group
	if err := api.client.ParseResponse(resp, &updatedGroup); err != nil {
		return nil, err
	}

	return &updatedGroup, nil
}

// Delete 删除用户组
func (api *GroupsAPI) Delete(id int64) error {
	endpoint := fmt.Sprintf("/api/groups/%d", id)
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
