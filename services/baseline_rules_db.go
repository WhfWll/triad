package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"smart/models/mysqls"

	log "github.com/sirupsen/logrus"
)

// InitBaselineRulesFromDB 启动时调用：从 host_baseline_rule 表加载启用的规则到引擎。
// 表不存在或查询失败时仅打日志，不返回致命错误。
func InitBaselineRulesFromDB(ctx context.Context) {
	if err := ReloadBaselineRulesFromDB(ctx); err != nil {
		if isMissingHostBaselineRuleTable(err) {
			log.Infof("host_baseline_rule table not found: %v, rules will be empty until import", err)
			return
		}
		if strings.Contains(err.Error(), "commands_json") {
			log.Errorf("host_baseline_rule invalid data: %v", err)
			return
		}
		log.Warnf("host_baseline_rule load skipped: %v", err)
	}
}

// ReloadBaselineRulesFromDB 从库重新加载规则（导入 SQL 后可调此接口热更新）。表缺失、查询失败或 bad JSON 会返回 error。
func ReloadBaselineRulesFromDB(ctx context.Context) error {
	return GetBaselineEngine().reloadFromDB(ctx)
}

func isMissingHostBaselineRuleTable(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "doesn't exist") || strings.Contains(s, "does not exist") ||
		strings.Contains(s, "no such table") || strings.Contains(s, "unknown table")
}

func (e *BaselineEngine) reloadFromDB(ctx context.Context) error {
	var model mysqls.HostBaselineRule
	rows, err := model.ListAllEnabled(ctx)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		log.Info("host baseline rules: no enabled rows in DB, keep built-in defaults")
		return nil
	}

	rules := make([]BaselineRule, 0, len(rows))
	for _, row := range rows {
		var cmds []string
		if strings.TrimSpace(row.CommandsJSON) != "" {
			if err := json.Unmarshal([]byte(row.CommandsJSON), &cmds); err != nil {
				return fmt.Errorf("rule_code=%d commands_json: %w", row.RuleCode, err)
			}
		}
		mt := row.MatchType
		if mt == "" {
			mt = "contains"
		}
		rules = append(rules, BaselineRule{
			ID:              row.RuleCode,
			Name:            row.Name,
			Description:     row.Description,
			Category:        row.Category,
			Risk:            row.Risk,
			OSType:          row.OSType,
			Commands:        cmds,
			ExpectedValue:   row.ExpectedValue,
			MatchType:       mt,
			FixSuggestion:   row.FixSuggestion,
			RiskDescription: row.RiskDescription,
		})
	}

	merged := mergeRulesWithBuiltin(rules)

	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = merged
	e.rulesMap = make(map[int][]BaselineRule)
	for _, rule := range merged {
		e.rulesMap[rule.OSType] = append(e.rulesMap[rule.OSType], rule)
	}
	log.Infof("host baseline rules loaded from DB: %d enabled rows, %d total after merge", len(rules), len(merged))
	return nil
}
