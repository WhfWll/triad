package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"time"

	"gorm.io/gorm"
)

// AiScenario AI应用场景模型
type AiScenario struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Description string    `gorm:"column:description;type:varchar(500)" json:"description"`
	Icon        string    `gorm:"column:icon;type:varchar(100)" json:"icon"`
	IsEnabled   int       `gorm:"column:is_enabled;type:tinyint;default:1;comment:是否启用：1-启用，2-禁用" json:"is_enabled"`
	LlmModelID  int       `gorm:"column:llm_model_id;type:int;default:0;comment:关联的大模型ID" json:"llm_model_id"`
	Prompt      string    `gorm:"column:prompt;type:text;comment:场景提示词" json:"prompt"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 表名
func (AiScenario) TableName() string {
	return "llm_scenarios"
}

// GetAiScenarioList 获取AI应用场景列表
func (a *AiScenario) GetAiScenarioList(ctx context.Context) ([]AiScenario, error) {
	var scenarios []AiScenario
	err := mysql.GetDB().WithContext(ctx).Order("id ASC").Find(&scenarios).Error
	return scenarios, err
}

// GetAiScenarioByID 根据ID获取AI应用场景
func (a *AiScenario) GetAiScenarioByID(ctx context.Context, id int) (*AiScenario, error) {
	var scenario AiScenario
	err := mysql.GetDB().WithContext(ctx).Where("id = ?", id).First(&scenario).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &scenario, nil
}

// UpdateAiScenario 更新AI应用场景
func (a *AiScenario) UpdateAiScenario(ctx context.Context, scenario *AiScenario) error {
	return mysql.GetDB().WithContext(ctx).Model(scenario).Where("id = ?", scenario.ID).Updates(scenario).Error
}

// InitAiScenarios 初始化AI应用场景数据
func InitAiScenarios() error {
	scenarios := []AiScenario{
		{
			Name:        "验证码识别",
			Description: "使用AI识别各种类型的验证码",
			Icon:        "captcha",
			IsEnabled:   1,
			LlmModelID:  0,
			Prompt:      "请识别图片中的验证码内容，只返回验证码的文字或数字，不要包含其他说明。",
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
		{
			Name:        "网页结构识别",
			Description: "分析网页结构，识别关键元素和内容",
			Icon:        "webpage",
			IsEnabled:   1,
			LlmModelID:  0,
			Prompt:      "请分析网页的HTML结构，识别出主要的功能模块、表单元素、链接等关键信息，并以结构化的方式返回。",
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
		{
			Name:        "源码分析",
			Description: "分析代码结构，发现潜在的安全问题",
			Icon:        "code",
			IsEnabled:   1,
			LlmModelID:  0,
			Prompt:      "请分析提供的源代码，识别潜在的安全漏洞、代码质量问题和改进建议。重点关注SQL注入、XSS、CSRF等常见安全问题。",
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
		{
			Name:        "报告生成",
			Description: "根据扫描结果生成专业的安全报告",
			Icon:        "report",
			IsEnabled:   1,
			LlmModelID:  0,
			Prompt:      "请根据提供的安全扫描结果，生成一份专业的安全评估报告。报告应包含：执行摘要、发现的漏洞详情、风险等级评估、修复建议等内容。",
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
	}

	for _, scenario := range scenarios {
		var existingScenario AiScenario
		err := mysql.GetDB().Where("name = ?", scenario.Name).First(&existingScenario).Error
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建新记录
			if err := mysql.GetDB().Create(&scenario).Error; err != nil {
				return err
			}
		}
	}

	return nil
}