package gophish

import (
	"testing"
)

// 测试配置
const (
	TestBaseURL = "https://192.168.0.177:3333"
	TestAPIKey  = "9ad173731769cb2a1c58a9d71969703c1ca28c173fe9d0368e3144d43415d8f1"
)

// TestUserManagement 测试用户管理
func TestUserManagement(t *testing.T) {
	client := New(TestBaseURL, TestAPIKey)

	// 测试获取所有用户账户
	t.Run("GetAll", func(t *testing.T) {
		accounts, err := client.UserManagement().GetAll()
		if err != nil {
			t.Logf("获取用户账户失败: %v", err)
		} else {
			t.Logf("找到 %d 个用户账户", len(accounts))
			// 显示每个账户的角色信息
			for i, account := range accounts {
				t.Logf("账户 %d: %s, 角色: %s", i+1, account.Username, account.GetRoleSlug())
			}
		}
	})

	// 测试创建用户账户
	t.Run("Create", func(t *testing.T) {
		newAccount := &CreateUserAccountRequest{
			Username: "testuser21",
			Password: "testpass1231",
			Role:     "user",
			Email:    "testuser21@4dogs.cn",
		}

		createdAccount, err := client.UserManagement().Create(newAccount)
		if err != nil {
			t.Logf("创建用户账户失败: %v", err)
		} else {
			t.Logf("创建用户账户成功，ID: %d, 角色: %s， api_key:%s", createdAccount.ID, createdAccount.GetRoleSlug(), createdAccount.ApiKey)
		}
	})
}

// TestSendingProfiles 测试发送配置文件管理
func TestSendingProfiles(t *testing.T) {
	client := New(TestBaseURL, TestAPIKey)

	// 测试获取所有发送配置文件
	t.Run("GetAll", func(t *testing.T) {
		profiles, err := client.SendingProfiles().GetAll()
		if err != nil {
			t.Logf("获取发送配置文件失败: %v", err)
		} else {
			t.Logf("找到 %d 个发送配置文件", len(profiles))
		}
	})

	// 测试创建发送配置文件
	t.Run("Create", func(t *testing.T) {
		newProfile := &CreateSendingProfileRequest{
			Name:             "测试配置",
			Host:             "smtp.example.com:587",
			FromAddress:      "test@example.com",
			InterfaceType:    "SMTP",
			IgnoreCertErrors: true,
		}

		createdProfile, err := client.SendingProfiles().Create(newProfile)
		if err != nil {
			t.Logf("创建发送配置文件失败: %v", err)
		} else {
			t.Logf("创建发送配置文件成功，ID: %d", createdProfile.ID)
		}
	})
}

// TestTemplates 测试模板管理
func TestTemplates(t *testing.T) {
	client := New(TestBaseURL, TestAPIKey)

	// 测试获取所有模板
	t.Run("GetAll", func(t *testing.T) {
		templates, err := client.Templates().GetAll()
		if err != nil {
			t.Logf("获取模板失败: %v", err)
		} else {
			t.Logf("找到 %d 个模板", len(templates))
		}
	})

	// 测试创建模板
	t.Run("Create", func(t *testing.T) {
		newTemplate := &CreateTemplateRequest{
			Name:    "测试模板",
			Subject: "测试邮件主题",
			HTML:    "<html><body><h1>测试邮件</h1><p>这是一个测试邮件。</p></body></html>",
			Text:    "测试邮件\n这是一个测试邮件。",
		}

		createdTemplate, err := client.Templates().Create(newTemplate)
		if err != nil {
			t.Logf("创建模板失败: %v", err)
		} else {
			t.Logf("创建模板成功，ID: %d", createdTemplate.ID)
		}
	})
}

// TestLandingPages 测试着陆页管理
func TestLandingPages(t *testing.T) {
	client := New(TestBaseURL, TestAPIKey)

	// 测试获取所有着陆页
	t.Run("GetAll", func(t *testing.T) {
		pages, err := client.LandingPages().GetAll()
		if err != nil {
			t.Logf("获取着陆页失败: %v", err)
		} else {
			t.Logf("找到 %d 个着陆页", len(pages))
		}
	})

	// 测试创建着陆页
	t.Run("Create", func(t *testing.T) {
		newPage := &CreateLandingPageRequest{
			Name:               "测试着陆页",
			HTML:               "<html><body><h1>欢迎</h1><form><input type='text' name='username'><input type='password' name='password'></form></body></html>",
			CaptureCredentials: true,
			RedirectURL:        "https://example.com",
		}

		createdPage, err := client.LandingPages().Create(newPage)
		if err != nil {
			t.Logf("创建着陆页失败: %v", err)
		} else {
			t.Logf("创建着陆页成功，ID: %d", createdPage.ID)
		}
	})
}

// TestUsers 测试用户管理
func TestUsers(t *testing.T) {
	client := New(TestBaseURL, TestAPIKey)

	// 测试获取所有用户
	t.Run("GetAll", func(t *testing.T) {
		users, err := client.Users().GetAll()
		if err != nil {
			t.Logf("获取用户失败: %v", err)
		} else {
			t.Logf("找到 %d 个用户", len(users))
		}
	})

	// 测试创建用户
	t.Run("Create", func(t *testing.T) {
		newUser := &CreateUserRequest{
			FirstName: "张",
			LastName:  "三",
			Email:     "zhangsan@example.com",
			Position:  "员工",
		}

		createdUser, err := client.Users().Create(newUser)
		if err != nil {
			t.Logf("创建用户失败: %v", err)
		} else {
			t.Logf("创建用户成功，邮箱: %s", createdUser.Email)
		}
	})
}

// TestGroups 测试用户组管理
func TestGroups(t *testing.T) {
	client := New(TestBaseURL, TestAPIKey)

	// 测试获取所有用户组
	t.Run("GetAll", func(t *testing.T) {
		groups, err := client.Groups().GetAll()
		if err != nil {
			t.Logf("获取用户组失败: %v", err)
		} else {
			t.Logf("找到 %d 个用户组", groups.Total)
		}
	})

	// 测试创建用户组
	/*t.Run("Create", func(t *testing.T) {
		newGroup := &CreateGroupRequest{
			Name: "测试用户组1",
			Targets: []Target{
				{
					FirstName: "李",
					LastName:  "四1",
					Email:     "lisi1@example.com",
					Position:  "经理",
				},
			},
		}

		createdGroup, err := client.Groups().Create(newGroup)
		if err != nil {
			t.Logf("创建用户组失败: %v", err)
		} else {
			t.Logf("创建用户组成功，ID: %d", createdGroup.ID)
		}
	})*/
}

// TestCampaigns 测试活动管理
func TestCampaigns(t *testing.T) {
	client := New(TestBaseURL, TestAPIKey)

	// 测试获取所有活动
	t.Run("GetAll", func(t *testing.T) {
		campaigns, err := client.Campaigns().GetAll()
		if err != nil {
			t.Logf("获取活动失败: %v", err)
		} else {
			t.Logf("找到 %d 个活动", campaigns.Total)
		}
	})

	// 测试创建活动
	t.Run("Create", func(t *testing.T) {
		// 首先获取必要的资源
		profiles, _ := client.SendingProfiles().GetAll()
		templates, _ := client.Templates().GetAll()
		pages, _ := client.LandingPages().GetAll()

		if len(profiles) > 0 && len(templates) > 0 && len(pages) > 0 {
			newCampaign := &CreateCampaignRequest{
				Name:     "测试活动",
				Template: NameInfo{Name: "ssss"},
				Page:     NameInfo{Name: pages[0].Name},
				Smtp:     NameInfo{Name: profiles[0].Name},
				Url:      "https://gophish.example.com",
			}

			createdCampaign, err := client.Campaigns().Create(newCampaign)
			if err != nil {
				t.Logf("创建活动失败: %v", err)
			} else {
				t.Logf("创建活动成功，ID: %d", createdCampaign.ID)
			}
		} else {
			t.Log("跳过创建活动测试：缺少必要的资源（发送配置文件、模板或着陆页）")
		}
	})
}

// TestIntegration 集成测试 - 测试完整的钓鱼活动流程
func TestIntegration(t *testing.T) {
	client := New(TestBaseURL, TestAPIKey)

	t.Log("开始集成测试...")

	// 1. 创建发送配置文件
	profile := &CreateSendingProfileRequest{
		Name:             "集成测试配置",
		Host:             "smtp.example.com:587",
		FromAddress:      "test@example.com",
		InterfaceType:    "SMTP",
		IgnoreCertErrors: true,
	}
	createdProfile, err := client.SendingProfiles().Create(profile)
	if err != nil {
		t.Logf("创建发送配置文件失败: %v", err)
		return
	}
	t.Logf("✓ 创建发送配置文件成功，ID: %d", createdProfile.ID)

	// 2. 创建模板
	template := &CreateTemplateRequest{
		Name:    "集成测试模板",
		Subject: "集成测试邮件",
		HTML:    "<html><body><h1>集成测试</h1></body></html>",
		Text:    "集成测试邮件",
	}
	createdTemplate, err := client.Templates().Create(template)
	if err != nil {
		t.Logf("创建模板失败: %v", err)
		return
	}
	t.Logf("✓ 创建模板成功，ID: %d", createdTemplate.ID)

	// 3. 创建着陆页
	page := &CreateLandingPageRequest{
		Name:               "集成测试着陆页",
		HTML:               "<html><body><h1>欢迎</h1></body></html>",
		CaptureCredentials: true,
		RedirectURL:        "https://example.com",
	}
	createdPage, err := client.LandingPages().Create(page)
	if err != nil {
		t.Logf("创建着陆页失败: %v", err)
		return
	}
	t.Logf("✓ 创建着陆页成功，ID: %d", createdPage.ID)

	// 4. 创建用户
	user := &CreateUserRequest{
		FirstName: "集成",
		LastName:  "测试",
		Email:     "integration@example.com",
		Position:  "测试员",
	}
	createdUser, err := client.Users().Create(user)
	if err != nil {
		t.Logf("创建用户失败: %v", err)
		return
	}
	t.Logf("✓ 创建用户成功，邮箱: %s", createdUser.Email)

	// 5. 创建用户组
	group := &CreateGroupRequest{
		Name: "集成测试组",
		Targets: []Target{
			{
				Email:     createdUser.Email,
				FirstName: createdUser.FirstName,
				LastName:  createdUser.LastName,
				Position:  createdUser.Position,
			},
		},
	}
	createdGroup, err := client.Groups().Create(group)
	if err != nil {
		t.Logf("创建用户组失败: %v", err)
		return
	}
	t.Logf("✓ 创建用户组成功，ID: %d", createdGroup.ID)

	// 6. 创建活动
	campaign := &CreateCampaignRequest{
		Name:     "集成测试活动",
		Template: NameInfo{Name: createdTemplate.Name},
		Page:     NameInfo{Name: createdPage.Name},
		Smtp:     NameInfo{Name: createdProfile.Name},
		Url:      "https://gophish.example.com",
	}
	createdCampaign, err := client.Campaigns().Create(campaign)
	if err != nil {
		t.Logf("创建活动失败: %v", err)
		return
	}
	t.Logf("✓ 创建活动成功，ID: %d", createdCampaign.ID)

	t.Log("集成测试完成！所有模块都正常工作。")
}
