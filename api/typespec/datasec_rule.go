package typespec

type DatasecRulesStatsResp struct {
	Total        int                            `json:"total"`
	EnabledTotal int                            `json:"enabledTotal"`
	BuiltinTotal int                            `json:"builtinTotal"`
	TargetTotal  int                            `json:"targetTotal"`
	ByDBType     []DatasecRulesCountByDBType    `json:"byDbType"`
	ByCategory   []DatasecRulesCountByCategory  `json:"byCategory"`
}

type DatasecRulesListResp struct {
	Total        int                            `json:"total"`
	EnabledTotal int                            `json:"enabledTotal"`
	BuiltinTotal int                            `json:"builtinTotal"`
	TargetTotal  int                            `json:"targetTotal"`
	ByDBType     []DatasecRulesCountByDBType    `json:"byDbType"`
	ByCategory   []DatasecRulesCountByCategory  `json:"byCategory"`
	Rules        []DatasecRuleListItem          `json:"rules"`
}

type DatasecRulesCountByDBType struct {
	DBType     int    `json:"dbType"`
	DBTypeName string `json:"dbTypeName"`
	Count      int    `json:"count"`
}

type DatasecRulesCountByCategory struct {
	Category     int    `json:"category"`
	CategoryName string `json:"categoryName"`
	Count        int    `json:"count"`
}

type DatasecRuleListItem struct {
	ID              int      `json:"id"`
	RuleCode        int      `json:"ruleCode"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	CategoryName    string   `json:"categoryName"`
	Risk            int      `json:"risk"`
	RiskName        string   `json:"riskName"`
	DBType          int      `json:"dbType"`
	DBTypeName      string   `json:"dbTypeName"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion,omitempty"`
	RiskDescription string   `json:"riskDescription,omitempty"`
	Queries         []string `json:"queries,omitempty"`
	Enabled         int      `json:"enabled"`
}

type DatasecRulesImportReq struct {
	Rules []DatasecRuleImportItem `json:"rules" binding:"required"`
}

type DatasecRuleImportItem struct {
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

type DatasecRulesImportResp struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Skipped int `json:"skipped"`
}

type DatasecRuleCreateReq struct {
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	Risk            int      `json:"risk"`
	DBType          int      `json:"dbType"`
	Queries         []string `json:"queries"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
	Enabled         int      `json:"enabled"`
}

type DatasecRuleUpdateReq struct {
	ID              int      `json:"id" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	Risk            int      `json:"risk"`
	DBType          int      `json:"dbType"`
	Queries         []string `json:"queries"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
	Enabled         int      `json:"enabled"`
}

type DatasecRuleDetailReq struct {
	ID int `form:"id" binding:"required"`
}

type DatasecRuleDeleteReq struct {
	ID int `form:"id" binding:"required"`
}

type DatasecCveImportReq struct {
	Limit int `json:"limit"`
}

type DatasecCveImportPreviewResp struct {
	AvailableInDB int  `json:"availableInDb"`
	TargetTotal   int  `json:"targetTotal"`
	Message       string `json:"message,omitempty"`
}
