package services

import (
	"fmt"
	"strings"

	"smart/tools/enums"
)

// DatasecCveImportPreview CVE 库可导入为数据安全规则的数量预览
type DatasecCveImportPreview struct {
	AvailableInDB int `json:"availableInDb"`
	TargetTotal   int `json:"targetTotal"`
}

const datasecCveTargetTotal = 3800

// datasecCveWhereClause 从 default-cve.db 筛选与数据库相关的 CVE
const datasecCveWhereClause = `(
  ',' || product || ',' LIKE '%,mysql,%' OR ',' || product || ',' LIKE '%,mariadb,%' OR
  ',' || product || ',' LIKE '%,postgresql,%' OR ',' || product || ',' LIKE '%,postgres,%' OR
  ',' || product || ',' LIKE '%,mongodb,%' OR ',' || product || ',' LIKE '%,redis,%' OR
  ',' || product || ',' LIKE '%,couchdb,%' OR ',' || product || ',' LIKE '%,sqlite,%' OR
  ',' || product || ',' LIKE '%,mssql,%' OR ',' || product || ',' LIKE '%,sql_server,%' OR
  ',' || vendor || ',' LIKE '%,oracle,%' OR ',' || product || ',' LIKE '%,db2,%' OR
  ',' || product || ',' LIKE '%,cassandra,%' OR ',' || product || ',' LIKE '%,elasticsearch,%' OR
  ',' || product || ',' LIKE '%,influxdb,%' OR ',' || product || ',' LIKE '%,memcached,%' OR
  ',' || product || ',' LIKE '%,tidb,%' OR ',' || product || ',' LIKE '%,clickhouse,%' OR
  LOWER(product) LIKE '%database%' OR LOWER(COALESCE(title_zh,'')) LIKE '%数据库%' OR
  LOWER(COALESCE(description_main_zh, description_main, '')) LIKE '%sql注入%' OR
  LOWER(COALESCE(description_main_zh, description_main, '')) LIKE '%sql injection%'
)`

// PreviewDatasecCveImport 统计 CVE 库中可导入条数
func PreviewDatasecCveImport() (*DatasecCveImportPreview, error) {
	cveDB := GetCveDB()
	if cveDB == nil || !cveDB.IsAvailable() {
		return &DatasecCveImportPreview{TargetTotal: datasecCveTargetTotal}, fmt.Errorf("CVE 库不可用，请确认 data/default-cve.db 存在")
	}
	var n int64
	err := cveDB.db.QueryRow("SELECT COUNT(DISTINCT cve) FROM cves WHERE " + datasecCveWhereClause).Scan(&n)
	if err != nil {
		return nil, err
	}
	return &DatasecCveImportPreview{AvailableInDB: int(n), TargetTotal: datasecCveTargetTotal}, nil
}

// BuildDatasecRulesFromCve 将 CVE 记录转换为可写入 datasec_rule 的导入项（知识库型，不执行 SQL）
func BuildDatasecRulesFromCve(limit int) ([]DatasecRuleExportItem, error) {
	cveDB := GetCveDB()
	if cveDB == nil || !cveDB.IsAvailable() {
		return nil, fmt.Errorf("CVE 库不可用")
	}
	if limit <= 0 {
		limit = datasecCveTargetTotal
	}

	query := fmt.Sprintf(`SELECT DISTINCT cve,
		COALESCE(NULLIF(title_zh, ''), cve),
		COALESCE(severity, ''),
		COALESCE(product, ''),
		COALESCE(vendor, ''),
		COALESCE(NULLIF(description_main_zh, ''), description_main, ''),
		COALESCE(solution, '')
		FROM cves WHERE %s
		ORDER BY CASE UPPER(severity)
			WHEN 'CRITICAL' THEN 0 WHEN 'HIGH' THEN 1 WHEN 'MEDIUM' THEN 2 WHEN 'LOW' THEN 3 ELSE 4 END,
			cve
		LIMIT ?`, datasecCveWhereClause)

	rows, err := cveDB.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]DatasecRuleExportItem, 0, limit)
	for rows.Next() {
		var cveID, title, severity, product, vendor, desc, solution string
		if err := rows.Scan(&cveID, &title, &severity, &product, &vendor, &desc, &solution); err != nil {
			continue
		}
		title = strings.TrimSpace(title)
		if title == "" {
			title = cveID
		}
		name := TruncateUTF8Bytes(fmt.Sprintf("[%s] %s", cveID, title), 255)
		fix := strings.TrimSpace(SanitizeUTF8(solution))
		if fix == "" {
			fix = "请参考 CVE 公告升级数据库或应用官方补丁；数据库安全检查任务已支持基于实例版本的 CVE 在线匹配。"
		}
		out = append(out, DatasecRuleExportItem{
			Name:            name,
			Description:     TruncateUTF8Bytes(firstNonEmptyStr(desc, title), 512),
			Category:        inferDatasecCategoryFromCve(title, desc),
			Risk:            cveSeverityToBaselineRisk(severity),
			DBType:          inferDatasecDBType(product, vendor),
			Queries:         []string{fmt.Sprintf("-- CVE-KB:%s", cveID)},
			ExpectedValue:   cveID,
			MatchType:       "cve_kb",
			FixSuggestion:   fix,
			RiskDescription: TruncateUTF8Bytes(fmt.Sprintf("CVE 知识库：%s；影响产品 %s", cveID, product), 512),
		})
	}
	return out, nil
}

func inferDatasecDBType(product, vendor string) int {
	p := "," + strings.ToLower(product) + ","
	v := "," + strings.ToLower(vendor) + ","
	switch {
	case strings.Contains(p, ",mysql,") || strings.Contains(p, ",mariadb,") || strings.Contains(p, ",percona,"):
		return enums.DBSupportTypeMySQL
	case strings.Contains(p, ",postgresql,") || strings.Contains(p, ",postgres,"):
		return enums.DBSupportTypePostgreSQL
	case strings.Contains(p, ",mongodb,"):
		return enums.DBSupportTypeMongoDB
	case strings.Contains(p, ",redis,"):
		return enums.DBSupportTypeRedis
	case strings.Contains(p, ",couchdb,"):
		return enums.DBSupportTypeCouchDB
	case strings.Contains(p, ",sqlite,"):
		return 0
	default:
		if strings.Contains(v, ",oracle,") && (strings.Contains(p, "mysql") || strings.Contains(p, "database")) {
			return enums.DBSupportTypeMySQL
		}
		return 0
	}
}

func inferDatasecCategoryFromCve(title, desc string) int {
	s := strings.ToLower(title + " " + desc)
	switch {
	case strings.Contains(s, "sql注入") || strings.Contains(s, "sql injection") || strings.Contains(s, "sqli"):
		return enums.DBCheckCategorySQLInjection
	case strings.Contains(s, "敏感") || strings.Contains(s, "泄露") || strings.Contains(s, "泄漏") || strings.Contains(s, "disclosure"):
		return enums.DBCheckCategorySensitiveData
	case strings.Contains(s, "认证") || strings.Contains(s, "密码") || strings.Contains(s, "password") || strings.Contains(s, "login") || strings.Contains(s, "auth"):
		return enums.DBCheckCategoryAuthentication
	case strings.Contains(s, "权限") || strings.Contains(s, "privilege") || strings.Contains(s, "authorization"):
		return enums.DBCheckCategoryAuthorization
	case strings.Contains(s, "审计") || strings.Contains(s, "audit") || strings.Contains(s, "log"):
		return enums.DBCheckCategoryAuditLog
	case strings.Contains(s, "加密") || strings.Contains(s, "ssl") || strings.Contains(s, "tls") || strings.Contains(s, "cipher"):
		return enums.DBCheckCategoryEncryption
	case strings.Contains(s, "网络") || strings.Contains(s, "network") || strings.Contains(s, "bind"):
		return enums.DBCheckCategoryNetwork
	default:
		return enums.DBCheckCategoryConfigSecure
	}
}

func cveSeverityToBaselineRisk(severity string) int {
	switch strings.ToUpper(strings.TrimSpace(severity)) {
	case "CRITICAL":
		return enums.BaselineRiskCritical
	case "HIGH":
		return enums.BaselineRiskHigh
	case "MEDIUM":
		return enums.BaselineRiskMiddle
	case "LOW":
		return enums.BaselineRiskLow
	default:
		return enums.BaselineRiskInfo
	}
}

func firstNonEmptyStr(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
