package application

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

// DataSecRuleApp 数据安全检测规则管理
type DataSecRuleApp struct{}

func (a *DataSecRuleApp) ReloadFromDB(ctx context.Context) error {
	return services.ReloadDatasecRulesFromDB(ctx)
}

func (a *DataSecRuleApp) GetRulesFromDB(ctx context.Context) *typespec.DatasecRulesListResp {
	var model mysqls.DatasecRule
	all, err := model.ListAll(ctx)
	if err != nil {
		return &typespec.DatasecRulesListResp{
			BuiltinTotal: services.BuiltinDatasecRuleCount(),
			TargetTotal:  3800,
		}
	}

	dbHits := map[int]int{}
	catHits := map[int]int{}
	for _, r := range all {
		if r.Enabled != 1 {
			continue
		}
		dt := r.DBType
		if dt == 0 {
			for _, t := range []int{1, 2, 3, 4, 5} {
				dbHits[t]++
			}
		} else {
			dbHits[dt]++
		}
		catHits[r.Category]++
	}

	resp := &typespec.DatasecRulesListResp{
		Total:        len(all),
		BuiltinTotal: services.BuiltinDatasecRuleCount(),
		TargetTotal:  3800,
		Rules:        make([]typespec.DatasecRuleListItem, 0, len(all)),
	}

	dbOrder := []int{0, enums.DBSupportTypeMySQL, enums.DBSupportTypePostgreSQL, enums.DBSupportTypeMongoDB, enums.DBSupportTypeRedis, enums.DBSupportTypeCouchDB}
	for _, dt := range dbOrder {
		if c, ok := dbHits[dt]; ok && c > 0 {
			name := "全部数据库"
			if dt > 0 {
				name = enums.BaselineEnum.GetDBTypeName(dt)
			}
			resp.ByDBType = append(resp.ByDBType, typespec.DatasecRulesCountByDBType{
				DBType: dt, DBTypeName: name, Count: c,
			})
		}
	}
	catKeys := make([]int, 0, len(catHits))
	for k := range catHits {
		catKeys = append(catKeys, k)
	}
	sort.Ints(catKeys)
	for _, cat := range catKeys {
		resp.ByCategory = append(resp.ByCategory, typespec.DatasecRulesCountByCategory{
			Category: cat, CategoryName: enums.BaselineEnum.GetDBCheckCategoryName(cat), Count: catHits[cat],
		})
	}

	for _, r := range all {
		var queries []string
		if r.QueriesJSON != "" {
			_ = json.Unmarshal([]byte(r.QueriesJSON), &queries)
		}
		dbName := enums.BaselineEnum.GetDBTypeName(r.DBType)
		if r.DBType == 0 {
			dbName = "全部"
		}
		resp.Rules = append(resp.Rules, typespec.DatasecRuleListItem{
			ID:              r.ID,
			RuleCode:        r.RuleCode,
			Name:            r.Name,
			Description:     r.Description,
			Category:        r.Category,
			CategoryName:    enums.BaselineEnum.GetDBCheckCategoryName(r.Category),
			Risk:            r.Risk,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.Risk),
			DBType:          r.DBType,
			DBTypeName:      dbName,
			ExpectedValue:   r.ExpectedValue,
			MatchType:       r.MatchType,
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			Queries:         queries,
			Enabled:         r.Enabled,
		})
	}
	return resp
}

func (a *DataSecRuleApp) ImportRules(ctx context.Context, req *typespec.DatasecRulesImportReq) *typespec.DatasecRulesImportResp {
	var model mysqls.DatasecRule
	var maxRuleCode int
	mysql.FromContext(ctx).Model(&model).Select("COALESCE(MAX(rule_code), 0)").Scan(&maxRuleCode)

	existing := make(map[string]bool)
	var all []mysqls.DatasecRule
	if err := mysql.FromContext(ctx).Model(&model).Select("name, category, db_type").Find(&all).Error; err == nil {
		for _, r := range all {
			existing[datasecRuleDedupKey(r.Name, r.Category, r.DBType)] = true
		}
	}

	success, skipped := 0, 0
	nextCode := maxRuleCode
	for _, item := range req.Rules {
		if item.Name == "" {
			skipped++
			continue
		}
		key := datasecRuleDedupKey(item.Name, item.Category, item.DBType)
		if existing[key] {
			skipped++
			continue
		}
		qJSON := "[]"
		if len(item.Queries) > 0 {
			b, _ := json.Marshal(item.Queries)
			qJSON = string(b)
		}
		mt := item.MatchType
		if mt == "" {
			mt = "contains"
		}
		nextCode++
		name, desc, fix, riskDesc, expected := services.SanitizeDatasecRuleText(
			item.Name, item.Description, item.FixSuggestion, item.RiskDescription, item.ExpectedValue,
		)
		rule := &mysqls.DatasecRule{
			RuleCode:        nextCode,
			Name:            name,
			Description:     desc,
			Category:        item.Category,
			Risk:            item.Risk,
			DBType:          item.DBType,
			QueriesJSON:     qJSON,
			ExpectedValue:   expected,
			MatchType:       mt,
			FixSuggestion:   fix,
			RiskDescription: riskDesc,
			Enabled:         1,
		}
		if err := model.Create(ctx, rule); err != nil {
			log.Errorf("ImportDatasecRules: %v", err)
			skipped++
			continue
		}
		existing[key] = true
		success++
	}
	_ = services.ReloadDatasecRulesFromDB(ctx)
	return &typespec.DatasecRulesImportResp{Total: len(req.Rules), Success: success, Skipped: skipped}
}

func (a *DataSecRuleApp) ImportBuiltinRules(ctx context.Context) *typespec.DatasecRulesImportResp {
	items := services.ExportBuiltinDatasecRuleItems()
	importItems := make([]typespec.DatasecRuleImportItem, 0, len(items))
	for _, it := range items {
		importItems = append(importItems, typespec.DatasecRuleImportItem{
			Name: it.Name, Description: it.Description, Category: it.Category, Risk: it.Risk,
			DBType: it.DBType, Queries: it.Queries, ExpectedValue: it.ExpectedValue, MatchType: it.MatchType,
			FixSuggestion: it.FixSuggestion, RiskDescription: it.RiskDescription,
		})
	}
	return a.ImportRules(ctx, &typespec.DatasecRulesImportReq{Rules: importItems})
}

func (a *DataSecRuleApp) CreateRule(ctx context.Context, req *typespec.DatasecRuleCreateReq) error {
	var model mysqls.DatasecRule
	var maxRuleCode int
	mysql.FromContext(ctx).Model(&model).Select("COALESCE(MAX(rule_code), 0)").Scan(&maxRuleCode)
	qJSON := queriesToJSON(req.Queries)
	mt := req.MatchType
	if mt == "" {
		mt = "contains"
	}
	enabled := req.Enabled
	if enabled == 0 {
		enabled = 1
	}
	rule := &mysqls.DatasecRule{
		RuleCode: maxRuleCode + 1, Name: req.Name, Description: req.Description,
		Category: req.Category, Risk: req.Risk, DBType: req.DBType,
		QueriesJSON: qJSON, ExpectedValue: req.ExpectedValue, MatchType: mt,
		FixSuggestion: req.FixSuggestion, RiskDescription: req.RiskDescription, Enabled: enabled,
	}
	if err := model.Create(ctx, rule); err != nil {
		return err
	}
	return services.ReloadDatasecRulesFromDB(ctx)
}

func (a *DataSecRuleApp) UpdateRule(ctx context.Context, req *typespec.DatasecRuleUpdateReq) error {
	var model mysqls.DatasecRule
	row, err := model.GetByID(ctx, req.ID)
	if err != nil || row == nil {
		return fmt.Errorf("规则不存在")
	}
	mt := req.MatchType
	if mt == "" {
		mt = "contains"
	}
	enabled := req.Enabled
	if enabled == 0 {
		enabled = 1
	}
	row.Name = req.Name
	row.Description = req.Description
	row.Category = req.Category
	row.Risk = req.Risk
	row.DBType = req.DBType
	row.QueriesJSON = queriesToJSON(req.Queries)
	row.ExpectedValue = req.ExpectedValue
	row.MatchType = mt
	row.FixSuggestion = req.FixSuggestion
	row.RiskDescription = req.RiskDescription
	row.Enabled = enabled
	if err := model.Update(ctx, row); err != nil {
		return err
	}
	return services.ReloadDatasecRulesFromDB(ctx)
}

func (a *DataSecRuleApp) DeleteRule(ctx context.Context, id int) error {
	var model mysqls.DatasecRule
	if err := model.Delete(ctx, id); err != nil {
		return err
	}
	return services.ReloadDatasecRulesFromDB(ctx)
}

func (a *DataSecRuleApp) GetRuleDetail(ctx context.Context, id int) (*typespec.DatasecRuleListItem, error) {
	var model mysqls.DatasecRule
	row, err := model.GetByID(ctx, id)
	if err != nil || row == nil {
		return nil, fmt.Errorf("规则不存在")
	}
	var queries []string
	if row.QueriesJSON != "" {
		_ = json.Unmarshal([]byte(row.QueriesJSON), &queries)
	}
	dbName := enums.BaselineEnum.GetDBTypeName(row.DBType)
	if row.DBType == 0 {
		dbName = "全部"
	}
	return &typespec.DatasecRuleListItem{
		ID: row.ID, RuleCode: row.RuleCode, Name: row.Name, Description: row.Description,
		Category: row.Category, CategoryName: enums.BaselineEnum.GetDBCheckCategoryName(row.Category),
		Risk: row.Risk, RiskName: enums.BaselineEnum.GetBaselineRiskName(row.Risk),
		DBType: row.DBType, DBTypeName: dbName, ExpectedValue: row.ExpectedValue, MatchType: row.MatchType,
		FixSuggestion: row.FixSuggestion, RiskDescription: row.RiskDescription, Queries: queries, Enabled: row.Enabled,
	}, nil
}

func datasecRuleDedupKey(name string, category, dbType int) string {
	return fmt.Sprintf("%s|%d|%d", name, category, dbType)
}

func queriesToJSON(queries []string) string {
	if len(queries) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(queries)
	return string(b)
}

func (a *DataSecRuleApp) PreviewCveImport() (*typespec.DatasecCveImportPreviewResp, error) {
	p, err := services.PreviewDatasecCveImport()
	if err != nil {
		return &typespec.DatasecCveImportPreviewResp{
			TargetTotal: 3800,
			Message:   err.Error() + "（路径: " + services.CveDBResolvedPath() + "）",
		}, err
	}
	msg := fmt.Sprintf("CVE 库中约 %d 条数据库相关漏洞可导入为规则（目标 %d 条）", p.AvailableInDB, p.TargetTotal)
	return &typespec.DatasecCveImportPreviewResp{
		AvailableInDB: p.AvailableInDB,
		TargetTotal:   p.TargetTotal,
		Message:       msg,
	}, nil
}

func (a *DataSecRuleApp) ImportFromCve(ctx context.Context, limit int) (*typespec.DatasecRulesImportResp, error) {
	items, err := services.BuildDatasecRulesFromCve(limit)
	if err != nil {
		return nil, err
	}
	importItems := make([]typespec.DatasecRuleImportItem, 0, len(items))
	for _, it := range items {
		importItems = append(importItems, typespec.DatasecRuleImportItem{
			Name: it.Name, Description: it.Description, Category: it.Category, Risk: it.Risk,
			DBType: it.DBType, Queries: it.Queries, ExpectedValue: it.ExpectedValue, MatchType: it.MatchType,
			FixSuggestion: it.FixSuggestion, RiskDescription: it.RiskDescription,
		})
	}
	return a.ImportRules(ctx, &typespec.DatasecRulesImportReq{Rules: importItems}), nil
}
