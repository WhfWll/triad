package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"smart/models/mysqls"
	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
)

var (
	datasecRulesMu     sync.RWMutex
	datasecRulesByDB   map[int][]dbRuleDef // key: dbType
	datasecRulesLoaded bool
)

// InitDatasecRulesFromDB 启动时从 datasec_rule 表加载规则到内存
func InitDatasecRulesFromDB(ctx context.Context) {
	if err := ReloadDatasecRulesFromDB(ctx); err != nil {
		if isMissingDatasecRuleTable(err) {
			log.Infof("datasec_rule table not found: %v, DB scan will use builtin rules until import", err)
			return
		}
		log.Warnf("datasec_rule load skipped: %v", err)
	}
}

func isMissingDatasecRuleTable(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "doesn't exist") || strings.Contains(s, "does not exist") ||
		strings.Contains(s, "no such table") || strings.Contains(s, "unknown table")
}

// ReloadDatasecRulesFromDB 热重载数据安全规则
func ReloadDatasecRulesFromDB(ctx context.Context) error {
	var model mysqls.DatasecRule
	rows, err := model.ListAllEnabled(ctx)
	if err != nil {
		return err
	}

	byDB := make(map[int][]dbRuleDef)
	for _, row := range rows {
		var queries []string
		if strings.TrimSpace(row.QueriesJSON) != "" {
			if err := json.Unmarshal([]byte(row.QueriesJSON), &queries); err != nil {
				return fmt.Errorf("rule_code=%d queries_json: %w", row.RuleCode, err)
			}
		}
		mt := row.MatchType
		if mt == "" {
			mt = "contains"
		}
		def := dbRuleDef{
			ID:            row.RuleCode,
			Name:          row.Name,
			Description:   row.Description,
			Category:      row.Category,
			Risk:          row.Risk,
			Queries:       queries,
			ExpectedValue: row.ExpectedValue,
			FixSuggestion: row.FixSuggestion,
			CheckFunc:     buildDBRuleCheckFunc(mt, row.ExpectedValue),
			KnowledgeOnly: mt == "cve_kb",
		}
		types := []int{row.DBType}
		if row.DBType == 0 {
			types = []int{
				enums.DBSupportTypeMySQL,
				enums.DBSupportTypePostgreSQL,
				enums.DBSupportTypeMongoDB,
				enums.DBSupportTypeRedis,
				enums.DBSupportTypeCouchDB,
			}
		}
		for _, dt := range types {
			byDB[dt] = append(byDB[dt], def)
		}
	}

	datasecRulesMu.Lock()
	defer datasecRulesMu.Unlock()
	datasecRulesByDB = byDB
	datasecRulesLoaded = len(rows) > 0
	if datasecRulesLoaded {
		log.Infof("datasec rules loaded from DB: %d enabled rows", len(rows))
	} else {
		log.Info("datasec rules: no enabled rows in DB, using builtin per DB type")
	}
	return nil
}

func buildDBRuleCheckFunc(matchType, expected string) func(string) bool {
	switch matchType {
	case "exact":
		return func(s string) bool { return strings.TrimSpace(s) == expected }
	case "not_contains":
		return func(s string) bool { return !strings.Contains(s, expected) }
	case "empty":
		return func(s string) bool { return strings.TrimSpace(s) == "" }
	case "always":
		return func(s string) bool { return true }
	default:
		return containsCheck(expected)
	}
}

func getDBBaselineRules(dbType int) []dbRuleDef {
	datasecRulesMu.RLock()
	defer datasecRulesMu.RUnlock()
	if datasecRulesLoaded {
		if rules, ok := datasecRulesByDB[dbType]; ok && len(rules) > 0 {
			return rules
		}
	}
	return getBuiltinDBBaselineRules(dbType)
}

func getBuiltinDBBaselineRules(dbType int) []dbRuleDef {
	switch dbType {
	case enums.DBSupportTypeMySQL:
		return getMySQLBaselineRules()
	case enums.DBSupportTypePostgreSQL:
		return getPostgreSQLBaselineRules()
	case enums.DBSupportTypeMongoDB:
		return getMongoDBBaselineRules()
	case enums.DBSupportTypeRedis:
		return getRedisBaselineRules()
	case enums.DBSupportTypeCouchDB:
		return getCouchDBBaselineRules()
	}
	return nil
}

// BuiltinDatasecRuleCount 内置规则总数（各库型之和，用于规则页展示）
func BuiltinDatasecRuleCount() int {
	types := []int{
		enums.DBSupportTypeMySQL,
		enums.DBSupportTypePostgreSQL,
		enums.DBSupportTypeMongoDB,
		enums.DBSupportTypeRedis,
		enums.DBSupportTypeCouchDB,
	}
	n := 0
	for _, t := range types {
		n += len(getBuiltinDBBaselineRules(t))
	}
	return n
}

// ExportBuiltinDatasecRuleItems 导出内置规则为可导入 JSON 结构
func ExportBuiltinDatasecRuleItems() []DatasecRuleExportItem {
	types := []struct {
		dbType int
		rules  []dbRuleDef
	}{
		{enums.DBSupportTypeMySQL, getMySQLBaselineRules()},
		{enums.DBSupportTypePostgreSQL, getPostgreSQLBaselineRules()},
		{enums.DBSupportTypeMongoDB, getMongoDBBaselineRules()},
		{enums.DBSupportTypeRedis, getRedisBaselineRules()},
		{enums.DBSupportTypeCouchDB, getCouchDBBaselineRules()},
	}
	out := make([]DatasecRuleExportItem, 0, 32)
	for _, block := range types {
		for _, r := range block.rules {
			out = append(out, DatasecRuleExportItem{
				Name:            r.Name,
				Description:     r.Description,
				Category:        r.Category,
				Risk:            r.Risk,
				DBType:          block.dbType,
				Queries:         r.Queries,
				ExpectedValue:   r.ExpectedValue,
				MatchType:       inferBuiltinMatchType(r),
				FixSuggestion:   r.FixSuggestion,
				RiskDescription: r.Description,
			})
		}
	}
	return out
}

// DatasecRuleExportItem 内置/导入规则条目
type DatasecRuleExportItem struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	Risk            int      `json:"risk"`
	DBType          int      `json:"dbType"`
	Queries         []string `json:"queries"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
}

func inferBuiltinMatchType(r dbRuleDef) string {
	ev := strings.TrimSpace(r.ExpectedValue)
	if ev == "空结果" || ev == "非空" {
		return "empty"
	}
	if ev == "3306" || ev == "27017" || ev == "不超过1000" || ev == "CouchDB" || ev == "认证配置" {
		return "always"
	}
	return "contains"
}
