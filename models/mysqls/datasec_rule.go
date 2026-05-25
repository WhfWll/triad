package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

// DatasecRule 数据安全检测规则（数据库基线 / SQL 注入 / 敏感数据等）
type DatasecRule struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`
	RuleCode        int       `gorm:"column:rule_code" json:"ruleCode"`
	Name            string    `gorm:"column:name" json:"name"`
	Description     string    `gorm:"column:description" json:"description"`
	Category        int       `gorm:"column:category" json:"category"`
	Risk            int       `gorm:"column:risk" json:"risk"`
	DBType          int       `gorm:"column:db_type" json:"dbType"`
	QueriesJSON     string    `gorm:"column:queries_json" json:"queriesJson"`
	ExpectedValue   string    `gorm:"column:expected_value" json:"expectedValue"`
	MatchType       string    `gorm:"column:match_type" json:"matchType"`
	FixSuggestion   string    `gorm:"column:fix_suggestion" json:"fixSuggestion"`
	RiskDescription string    `gorm:"column:risk_description" json:"riskDescription"`
	Enabled         int       `gorm:"column:enabled" json:"enabled"`
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (DatasecRule) TableName() string {
	return "datasec_rule"
}

func (m *DatasecRule) ListAllEnabled(ctx context.Context) ([]DatasecRule, error) {
	var list []DatasecRule
	err := mysql.FromContext(ctx).Model(m).Where("enabled = ?", 1).Order("rule_code asc").Find(&list).Error
	return list, err
}

func (m *DatasecRule) ListAll(ctx context.Context) ([]DatasecRule, error) {
	var list []DatasecRule
	err := mysql.FromContext(ctx).Model(m).Order("id desc").Find(&list).Error
	return list, err
}

// DatasecRuleSummary 列表摘要（不含 queries_json / fix_suggestion / risk_description 等大字段）
type DatasecRuleSummary struct {
	ID            int    `gorm:"column:id"`
	RuleCode      int    `gorm:"column:rule_code"`
	Name          string `gorm:"column:name"`
	Description   string `gorm:"column:description"`
	Category      int    `gorm:"column:category"`
	Risk          int    `gorm:"column:risk"`
	DBType        int    `gorm:"column:db_type"`
	ExpectedValue string `gorm:"column:expected_value"`
	MatchType     string `gorm:"column:match_type"`
	Enabled       int    `gorm:"column:enabled"`
}

// ListSummary 返回规则摘要列表，按 id 降序
func (m *DatasecRule) ListSummary(ctx context.Context) ([]DatasecRuleSummary, error) {
	var list []DatasecRuleSummary
	err := mysql.FromContext(ctx).Model(m).
		Select("id, rule_code, name, description, category, risk, db_type, expected_value, match_type, enabled").
		Order("id desc").
		Find(&list).Error
	return list, err
}

// CountAll 返回规则总数（含已停用）
func (m *DatasecRule) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := mysql.FromContext(ctx).Model(m).Count(&count).Error
	return count, err
}

func (m *DatasecRule) GetByID(ctx context.Context, id int) (*DatasecRule, error) {
	var rule DatasecRule
	err := mysql.FromContext(ctx).Model(m).Where("id = ?", id).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (m *DatasecRule) Create(ctx context.Context, rule *DatasecRule) error {
	now := time.Now()
	rule.CreateTime = now
	rule.UpdateTime = now
	if rule.Enabled == 0 {
		rule.Enabled = 1
	}
	return mysql.FromContext(ctx).Model(m).Create(rule).Error
}

func (m *DatasecRule) Update(ctx context.Context, rule *DatasecRule) error {
	rule.UpdateTime = time.Now()
	return mysql.FromContext(ctx).Model(m).Where("id = ?", rule.ID).Select("*").Omit("id", "create_time").Updates(rule).Error
}

func (m *DatasecRule) Delete(ctx context.Context, id int) error {
	return mysql.FromContext(ctx).Model(m).Where("id = ?", id).Delete(&DatasecRule{}).Error
}

func (m *DatasecRule) CountEnabled(ctx context.Context) (int64, error) {
	var n int64
	err := mysql.FromContext(ctx).Model(m).Where("enabled = ?", 1).Count(&n).Error
	return n, err
}

// CountByCategoryEnabled 按分类统计启用规则数
func (m *DatasecRule) CountByCategoryEnabled(ctx context.Context) (map[int]int, error) {
	type row struct {
		Category int `gorm:"column:category"`
		Count    int `gorm:"column:cnt"`
	}
	var rows []row
	err := mysql.FromContext(ctx).Model(m).
		Select("category, COUNT(*) AS cnt").
		Where("enabled = ?", 1).
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

// CountByDBTypeEnabled 按数据库类型统计启用规则数（db_type=0 计入全部 5 类）
func (m *DatasecRule) CountByDBTypeEnabled(ctx context.Context) (map[int]int, error) {
	type row struct {
		DBType int `gorm:"column:db_type"`
		Count  int `gorm:"column:cnt"`
	}
	var rows []row
	err := mysql.FromContext(ctx).Model(m).
		Select("db_type, COUNT(*) AS cnt").
		Where("enabled = ?", 1).
		Group("db_type").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[int]int)
	universal := 0
	for _, r := range rows {
		if r.DBType == 0 {
			universal = r.Count
		} else {
			out[r.DBType] = r.Count
		}
	}
	if universal > 0 {
		for _, t := range []int{1, 2, 3, 4, 5} {
			out[t] += universal
		}
	}
	return out, nil
}
