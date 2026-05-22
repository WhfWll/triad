package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"smart/models/mysqls"
	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
)

const (
	cveScanMaxRecords = 800
	cveScanMaxMatches = 50
	cveRuleIDBase     = 100000
)

// runCVEVersionChecks 基于实例版本与 CVE 库 CPE 做在线漏洞匹配
func (c *DBBaselineChecker) runCVEVersionChecks(ctx context.Context, task *DBCheckTask, conn *DBConnection, version string) []mysqls.DBCheckResult {
	if strings.TrimSpace(version) == "" {
		return []mysqls.DBCheckResult{
			c.buildVersionInfoResult(task, "", enums.BaselineCheckResultSkip, "未能识别数据库版本，已跳过 CVE 版本匹配"),
		}
	}

	cveDB := GetCveDB()
	if cveDB == nil || !cveDB.IsAvailable() {
		return []mysqls.DBCheckResult{
			c.buildVersionInfoResult(task, version, enums.BaselineCheckResultSkip, "CVE 库不可用，已跳过版本漏洞匹配"),
		}
	}

	records, err := cveDB.QueryDatasecCVEsForDBType(task.DBType, cveScanMaxRecords)
	if err != nil {
		log.Warnf("datasec cve query dbType=%d: %v", task.DBType, err)
		return []mysqls.DBCheckResult{
			c.buildVersionInfoResult(task, version, enums.BaselineCheckResultError, "CVE 查询失败: "+err.Error()),
		}
	}

	out := []mysqls.DBCheckResult{
		c.buildVersionInfoResult(task, version, enums.BaselineCheckResultPass, "当前实例版本"),
	}

	matchCount := 0
	for i, rec := range records {
		if matchCount >= cveScanMaxMatches {
			break
		}
		ok, _ := cveDB.MatchCpe(version, rec.CpeConfigurations)
		if !ok {
			continue
		}
		matchCount++
		title := strings.TrimSpace(rec.TitleZh)
		if title == "" {
			title = rec.Cve
		}
		risk := cveSeverityToBaselineRisk(rec.Severity)
		if risk == enums.BaselineRiskInfo {
			risk = enums.BaselineRiskMiddle
		}
		out = append(out, mysqls.DBCheckResult{
			TaskID:          task.TaskID,
			TargetID:        task.TargetID,
			TargetIP:        task.Host,
			DBType:          task.DBType,
			DBName:          task.DBName,
			CheckCategory:   inferDatasecCategoryFromCve(title, rec.TitleZh),
			RuleID:          cveRuleIDBase + i,
			RuleName:        fmt.Sprintf("[%s] %s", rec.Cve, TruncateUTF8Bytes(title, 200)),
			CheckResult:     enums.BaselineCheckResultFail,
			ExpectedValue:   "已修复版本 / 不受影响",
			ActualValue:     fmt.Sprintf("version=%s severity=%s product=%s", version, rec.Severity, rec.Product),
			RiskLevel:       risk,
			FixSuggestion:   fmt.Sprintf("当前版本 %s 命中 %s，请升级至官方修复版本或应用安全补丁。", version, rec.Cve),
			RiskDescription: TruncateUTF8Bytes(fmt.Sprintf("CVE 版本匹配：%s；%s", rec.Cve, title), 512),
			CreateTime:      time.Now(),
		})
	}

	if matchCount == 0 {
		out = append(out, mysqls.DBCheckResult{
			TaskID:          task.TaskID,
			TargetID:        task.TargetID,
			TargetIP:        task.Host,
			DBType:          task.DBType,
			DBName:          task.DBName,
			CheckCategory:   enums.DBCheckCategoryConfigSecure,
			RuleID:          cveRuleIDBase,
			RuleName:        "CVE 版本匹配",
			CheckResult:     enums.BaselineCheckResultPass,
			ExpectedValue:   "无已知 CVE 命中",
			ActualValue:     fmt.Sprintf("version=%s scanned=%d", version, len(records)),
			RiskLevel:       enums.BaselineRiskInfo,
			FixSuggestion:   "建议定期重新扫描并关注数据库安全公告。",
			RiskDescription: fmt.Sprintf("在 %d 条 CVE 记录中未发现与版本 %s 匹配的漏洞", len(records), version),
			CreateTime:      time.Now(),
		})
	}
	return out
}

func (c *DBBaselineChecker) buildVersionInfoResult(task *DBCheckTask, version string, result int, desc string) mysqls.DBCheckResult {
	actual := version
	if actual == "" {
		actual = "unknown"
	}
	return mysqls.DBCheckResult{
		TaskID:          task.TaskID,
		TargetID:        task.TargetID,
		TargetIP:        task.Host,
		DBType:          task.DBType,
		DBName:          task.DBName,
		CheckCategory:   enums.DBCheckCategoryConfigSecure,
		RuleID:          99999,
		RuleName:        "数据库版本识别",
		CheckResult:     result,
		ExpectedValue:   "可识别版本",
		ActualValue:     actual,
		RiskLevel:       enums.BaselineRiskInfo,
		FixSuggestion:   "",
		RiskDescription: desc,
		CreateTime:      time.Now(),
	}
}

// QueryDatasecCVEsForDBType 按库型查询可用于版本匹配的 CVE 记录
func (d *CveDB) QueryDatasecCVEsForDBType(dbType int, limit int) ([]CveRecord, error) {
	if !d.IsAvailable() {
		return nil, fmt.Errorf("cve db unavailable")
	}
	if limit <= 0 {
		limit = cveScanMaxRecords
	}
	filter := datasecDBTypeProductFilter(dbType)
	if filter == "" {
		return nil, nil
	}
	query := fmt.Sprintf(`SELECT cve, COALESCE(title_zh, ''), COALESCE(severity, ''), COALESCE(product, ''), COALESCE(vendor, ''), COALESCE(cpe_configurations, '{}')
		FROM cves WHERE %s AND (%s)
		ORDER BY CASE UPPER(severity)
			WHEN 'CRITICAL' THEN 0 WHEN 'HIGH' THEN 1 WHEN 'MEDIUM' THEN 2 WHEN 'LOW' THEN 3 ELSE 4 END,
			cve
		LIMIT ?`, datasecCveWhereClause, filter)

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CveRecord
	for rows.Next() {
		var r CveRecord
		var cpeStr string
		if err := rows.Scan(&r.Cve, &r.TitleZh, &r.Severity, &r.Product, &r.Vendor, &cpeStr); err != nil {
			continue
		}
		r.CpeConfigurations = []byte(cpeStr)
		r.SeverityLevel = severityToLevel(r.Severity)
		results = append(results, r)
	}
	return results, nil
}

func datasecDBTypeProductFilter(dbType int) string {
	switch dbType {
	case enums.DBSupportTypeMySQL:
		return `',' || product || ',' LIKE '%,mysql,%' OR ',' || product || ',' LIKE '%,mariadb,%' OR ',' || product || ',' LIKE '%,percona,%'`
	case enums.DBSupportTypePostgreSQL:
		return `',' || product || ',' LIKE '%,postgresql,%' OR ',' || product || ',' LIKE '%,postgres,%'`
	case enums.DBSupportTypeMongoDB:
		return `',' || product || ',' LIKE '%,mongodb,%'`
	case enums.DBSupportTypeRedis:
		return `',' || product || ',' LIKE '%,redis,%'`
	case enums.DBSupportTypeCouchDB:
		return `',' || product || ',' LIKE '%,couchdb,%'`
	default:
		return ""
	}
}

func IsCveDBCheckResult(r mysqls.DBCheckResult) bool {
	if r.RuleID >= cveRuleIDBase {
		return true
	}
	return strings.Contains(r.RuleName, "CVE-") || strings.Contains(r.RiskDescription, "CVE 版本匹配")
}

func isVersionInfoResult(r mysqls.DBCheckResult) bool {
	return r.RuleID == 99999 && r.RuleName == "数据库版本识别"
}
