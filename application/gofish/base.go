package gofish

import (
	"context"
	"fmt"
	"smart/models/mysqls"
	"smart/tools/gophish"
	"smart/tools/utils"
	"strings"

	"gitlabee.4dogs.cn/common/config"
)

// GoPhishBiz .
type GoPhishBiz struct {
	GoPhishAPIKey
}

type GoPhishAPIKey struct{}

// getGOPhishServiceInfo .获取gohish 的服务信息
func getGOPhishServiceInfo() (*config.HttpConfig, error) {
	return config.GetClintConfig("service_gophish")
}

// GetGoPhishUserAPiKeyInfo .获取用户的对应的id
func (key *GoPhishAPIKey) GetGoPhishUserAPiKeyInfo(ctx context.Context, userId int, username string) (string, error) {
	goPhishAPIKey := ""
	userGoPhishKeyModel := mysqls.UserGoPhishKey{}
	info, err := userGoPhishKeyModel.GetUserGoPhishKeyInfoByUserID(ctx, userId)
	if err != nil {
		return "", err
	}

	if info.ID == 0 {
		// 创建一个新的账号
		goPhishAPIKey, err = key.CreateUserGoPhishKey(ctx, userId, username)
		if err != nil {
			return "", err
		}
	} else {
		goPhishAPIKey = info.GoPhishKey
	}

	return goPhishAPIKey, nil
}

// CreateUserGoPhishKey 创建个新用户
func (key *GoPhishAPIKey) CreateUserGoPhishKey(ctx context.Context, userId int, username string) (string, error) {
	password, err := utils.GeneratePassword(utils.DefaultConfig)
	if err != nil {
		return "", err
	}

	// 获取管理员使用的apikey
	goPhishAPIKey, _ := key.GetGoPhishUserAPiKeyInfo(ctx, 0, "")

	hc, err := getGOPhishServiceInfo()
	if err != nil {
		return "", err
	}

	client := gophish.New(hc.BaseUri, goPhishAPIKey)
	newAccount := &gophish.CreateUserAccountRequest{
		Username: fmt.Sprintf("%v%v", username, userId),
		Password: fmt.Sprintf("%s%d", password, userId),
		Role:     "user",
		Email:    username + "@4dogs.cn",
	}

	createdAccount, err := client.UserManagement().Create(newAccount)
	if err != nil {
		return "", err
	}

	userGoPhishKeyModel := mysqls.UserGoPhishKey{
		UserID:     userId,
		GoPhishKey: createdAccount.GetApiKey(),
	}

	err = userGoPhishKeyModel.AddUserGoPhishKey(ctx)

	return createdAccount.GetApiKey(), err
}

// 分页方法
func (biz *GoPhishBiz) paginateGroups(groups []gophish.Group, page, size int) *gophish.GroupResponse {
	// 计算总数
	total := int64(len(groups))

	// 如果没有数据，直接返回空响应
	if total == 0 {
		return &gophish.GroupResponse{
			Groups: []gophish.Group{},
			Total:  total,
		}
	}

	// 计算分页索引
	startIndex := (page - 1) * size
	endIndex := startIndex + size

	// 边界检查
	if startIndex >= len(groups) {
		// 超出范围，返回空列表
		return &gophish.GroupResponse{
			Groups: []gophish.Group{},
			Total:  total,
		}
	}

	if endIndex > len(groups) {
		endIndex = len(groups)
	}

	// 分页数据
	pagedGroups := groups[startIndex:endIndex]

	return &gophish.GroupResponse{
		Groups: pagedGroups,
		Total:  total,
	}
}

// 辅助函数：搜索过滤
func (biz *GoPhishBiz) filterGroupsBySearch(groups []gophish.Group, search string) []gophish.Group {
	if search == "" {
		return groups
	}

	filtered := make([]gophish.Group, 0)
	searchLower := strings.ToLower(search)

	for _, group := range groups {
		// 根据实际需求调整搜索字段，这里假设搜索组名
		if strings.Contains(strings.ToLower(group.Name), searchLower) {
			filtered = append(filtered, group)
		}
	}
	return filtered
}

// 更通用的版本 - 返回分页后的原始类型切片
func paginate[T any](items []T, page, size int) ([]T, int64) {
	// 计算总数
	total := int64(len(items))

	// 如果没有数据，直接返回空响应
	if total == 0 {
		return []T{}, total
	}

	// 计算分页索引
	startIndex := (page - 1) * size
	endIndex := startIndex + size

	// 边界检查
	if startIndex >= len(items) {
		return []T{}, total
	}

	if endIndex > len(items) {
		endIndex = len(items)
	}

	return items[startIndex:endIndex], total
}

// 通用搜索过滤方法
func filterItemsBySearch[T any](items []T, search string, searchFunc func(T, string) bool) []T {
	if search == "" {
		return items
	}

	filtered := make([]T, 0)
	searchLower := strings.ToLower(search)

	for _, item := range items {
		if searchFunc(item, searchLower) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// 通用搜索过滤方法
func filterItemsByStatus[T any](items []T, status string, searchFunc func(T, string) bool) []T {
	if status == "" {
		return items
	}

	filtered := make([]T, 0)
	statusLower := strings.ToLower(status)
	for _, item := range items {
		if searchFunc(item, statusLower) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// 特定类型的搜索函数
func (biz *GoPhishBiz) groupSearchFunc(group gophish.Group, searchLower string) bool {
	return strings.Contains(strings.ToLower(group.Name), searchLower)
}

func (biz *GoPhishBiz) landingPageSearchFunc(landingPage gophish.LandingPage, searchLower string) bool {
	return strings.Contains(strings.ToLower(landingPage.Name), searchLower)
}

func (biz *GoPhishBiz) sendingProfileSearchFunc(profile gophish.SendingProfile, searchLower string) bool {
	return strings.Contains(strings.ToLower(profile.Name), searchLower)
}

func (biz *GoPhishBiz) templatesSearchFunc(template gophish.Template, searchLower string) bool {
	return strings.Contains(strings.ToLower(template.Name), searchLower)
}

func (biz *GoPhishBiz) campaignSearchFunc(campaign gophish.CampaignList, searchLower string) bool {
	return strings.Contains(strings.ToLower(campaign.Name), searchLower)
}

func (biz *GoPhishBiz) campaignStatusFunc(campaign gophish.CampaignList, status string) bool {
	return strings.Contains(strings.ToLower(campaign.Status), status)
}
