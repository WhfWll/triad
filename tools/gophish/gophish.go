package gophish

// Gophish 是主要的客户端结构体，包含所有模块
type Gophish struct {
	client *Client
}

// New 创建一个新的 Gophish 客户端实例
func New(baseURL, apiKey string) *Gophish {
	return &Gophish{
		client: NewClient(baseURL, apiKey),
	}
}

// SendingProfiles 返回发送配置文件模块
func (g *Gophish) SendingProfiles() *SendingProfilesAPI {
	return &SendingProfilesAPI{client: g.client}
}

// Templates 返回模板模块
func (g *Gophish) Templates() *TemplatesAPI {
	return &TemplatesAPI{client: g.client}
}

// LandingPages 返回着陆页模块
func (g *Gophish) LandingPages() *LandingPagesAPI {
	return &LandingPagesAPI{client: g.client}
}

// Users 返回用户模块
func (g *Gophish) Users() *UsersAPI {
	return &UsersAPI{client: g.client}
}

// Groups 返回用户组模块
func (g *Gophish) Groups() *GroupsAPI {
	return &GroupsAPI{client: g.client}
}

// Campaigns 返回活动模块
func (g *Gophish) Campaigns() *CampaignsAPI {
	return &CampaignsAPI{client: g.client}
}

// UserManagement 返回用户管理模块
func (g *Gophish) UserManagement() *UserManagementAPI {
	return &UserManagementAPI{client: g.client}
}
