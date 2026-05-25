package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

// HostBaselineRule 主机核查规则（可替代内置规则；enabled=1 的行会加载入引擎）
type HostBaselineRule struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`
	RuleCode        int       `gorm:"column:rule_code" json:"ruleCode"` // 与核查结果 rule_id 对应，全局唯一
	Name            string    `gorm:"column:name" json:"name"`
	Description     string    `gorm:"column:description" json:"description"`
	Category        int       `gorm:"column:category" json:"category"`
	Risk            int       `gorm:"column:risk" json:"risk"`
	OSType          int       `gorm:"column:os_type" json:"osType"`
	CommandsJSON    string    `gorm:"column:commands_json" json:"commandsJson"` // JSON 数组，如 ["cmd1","cmd2"]
	ExpectedValue   string    `gorm:"column:expected_value" json:"expectedValue"`
	MatchType       string    `gorm:"column:match_type" json:"matchType"`
	FixSuggestion   string    `gorm:"column:fix_suggestion" json:"fixSuggestion"`
	RiskDescription string    `gorm:"column:risk_description" json:"riskDescription"`
	Enabled         int       `gorm:"column:enabled" json:"enabled"` // 1 启用 0 停用
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (HostBaselineRule) TableName() string {
	return "host_baseline_rule"
}

// ListAllEnabled 返回所有启用的规则，按 rule_code 排序
func (m *HostBaselineRule) ListAllEnabled(ctx context.Context) ([]HostBaselineRule, error) {
	var list []HostBaselineRule
	err := mysql.FromContext(ctx).Model(m).Where("enabled = ?", 1).Order("rule_code asc").Find(&list).Error
	return list, err
}

// ListAll 返回所有规则（含已停用），按 id 降序
func (m *HostBaselineRule) ListAll(ctx context.Context) ([]HostBaselineRule, error) {
	var list []HostBaselineRule
	err := mysql.FromContext(ctx).Model(m).Order("id desc").Find(&list).Error
	return list, err
}

// HostBaselineRuleSummary 列表摘要（不含 commands_json / fix_suggestion / risk_description 等大字段）
type HostBaselineRuleSummary struct {
	ID            int    `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	Description   string `gorm:"column:description"`
	Category      int    `gorm:"column:category"`
	Risk          int    `gorm:"column:risk"`
	OSType        int    `gorm:"column:os_type"`
	ExpectedValue string `gorm:"column:expected_value"`
	MatchType     string `gorm:"column:match_type"`
	Enabled       int    `gorm:"column:enabled"`
}

// ListSummary 返回规则摘要列表，按 id 降序
func (m *HostBaselineRule) ListSummary(ctx context.Context) ([]HostBaselineRuleSummary, error) {
	var list []HostBaselineRuleSummary
	err := mysql.FromContext(ctx).Model(m).
		Select("id, name, description, category, risk, os_type, expected_value, match_type, enabled").
		Order("id desc").
		Find(&list).Error
	return list, err
}

// CountAll 返回规则总数（含已停用）
func (m *HostBaselineRule) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := mysql.FromContext(ctx).Model(m).Count(&count).Error
	return count, err
}

// CountByOSType 按操作系统分组统计
func (m *HostBaselineRule) CountByOSType(ctx context.Context) (map[int]int, error) {
	type row struct {
		OSType int `gorm:"column:os_type"`
		Count  int `gorm:"column:cnt"`
	}
	var rows []row
	err := mysql.FromContext(ctx).Model(m).
		Select("os_type, COUNT(*) AS cnt").
		Group("os_type").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[int]int, len(rows))
	for _, r := range rows {
		out[r.OSType] = r.Count
	}
	return out, nil
}

// CountByCategory 按核查分类分组统计
func (m *HostBaselineRule) CountByCategory(ctx context.Context) (map[int]int, error) {
	type row struct {
		Category int `gorm:"column:category"`
		Count    int `gorm:"column:cnt"`
	}
	var rows []row
	err := mysql.FromContext(ctx).Model(m).
		Select("category, COUNT(*) AS cnt").
		Group("category").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[int]int, len(rows))
	for _, r := range rows {
		out[r.Category] = r.Count
	}
	return out, nil
}

// GetByID 根据主键获取单条规则
func (m *HostBaselineRule) GetByID(ctx context.Context, id int) (*HostBaselineRule, error) {
	var rule HostBaselineRule
	err := mysql.FromContext(ctx).Model(m).Where("id = ?", id).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// FindByName 根据名称查找规则，返回第一条匹配的记录
func (m *HostBaselineRule) FindByName(ctx context.Context, name string) (*HostBaselineRule, error) {
	var rule HostBaselineRule
	err := mysql.FromContext(ctx).Model(m).Where("name = ?", name).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// ExistsByName 检查指定名称的规则是否已存在
func (m *HostBaselineRule) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := mysql.FromContext(ctx).Model(m).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create 新增规则
func (m *HostBaselineRule) Create(ctx context.Context, rule *HostBaselineRule) error {
	now := time.Now()
	rule.CreateTime = now
	rule.UpdateTime = now
	if rule.Enabled == 0 {
		rule.Enabled = 1
	}
	return mysql.FromContext(ctx).Model(m).Create(rule).Error
}

// Update 更新规则
func (m *HostBaselineRule) Update(ctx context.Context, rule *HostBaselineRule) error {
	rule.UpdateTime = time.Now()
	return mysql.FromContext(ctx).Model(m).Where("id = ?", rule.ID).Select("*").Omit("id", "create_time").Updates(rule).Error
}

// Delete 删除规则
func (m *HostBaselineRule) Delete(ctx context.Context, id int) error {
	return mysql.FromContext(ctx).Model(m).Where("id = ?", id).Delete(&HostBaselineRule{}).Error
}
