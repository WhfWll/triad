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
