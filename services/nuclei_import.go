package services

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	sqlitemodel "smart/models/sqlite"
	"smart/tools/enums"
	"strings"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gopkg.in/yaml.v3"
)

var nucleiSeverityToRisk = map[string]int{
	"critical": enums.VulLibrariesRiskDead,
	"high":     enums.VulLibrariesRiskHigh,
	"medium":   enums.VulLibrariesRiskMiddle,
	"low":      enums.VulLibrariesRiskLow,
	"info":     enums.VulLibrariesRiskInfo,
	"unknown":  enums.VulLibrariesRiskInfo,
}

var nucleiCveRegexp = regexp.MustCompile(`CVE-\d{4}-\d+`)

type nucleiTemplateDoc struct {
	ID        string      `yaml:"id"`
	Workflows interface{} `yaml:"workflows"`
	Flow      interface{} `yaml:"flow"`
	Info      struct {
		Name           string      `yaml:"name"`
		Severity       string      `yaml:"severity"`
		Description    string      `yaml:"description"`
		Tags           interface{} `yaml:"tags"`
		Reference      interface{} `yaml:"reference"`
		Remediation    string      `yaml:"remediation"`
		Classification struct {
			CVEID     interface{} `yaml:"cve-id"`
			CNVD      string      `yaml:"cnvd"`
			CNNVD     string      `yaml:"cnnvd"`
			CVSSScore string      `yaml:"cvss-score"`
		} `yaml:"classification"`
	} `yaml:"info"`
}

type nucleiSkipError struct {
	Reason string
}

func (e *nucleiSkipError) Error() string {
	return e.Reason
}

func ensureImportStats(result *VulkitImportResult) *ImportRuleStats {
	if result.Stats == nil {
		result.Stats = &ImportRuleStats{
			BySeverity: map[string]int{},
			ByType:     map[string]int{},
			SkipReason: map[string]int{},
		}
	}
	if result.Stats.BySeverity == nil {
		result.Stats.BySeverity = map[string]int{}
	}
	if result.Stats.ByType == nil {
		result.Stats.ByType = map[string]int{}
	}
	if result.Stats.SkipReason == nil {
		result.Stats.SkipReason = map[string]int{}
	}
	return result.Stats
}

func addSkipReasonStat(result *VulkitImportResult, reason string) {
	stats := ensureImportStats(result)
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "鍏朵粬"
	}
	stats.SkipReason[reason]++
}

func addImportedTemplateStat(result *VulkitImportResult, doc nucleiTemplateDoc, relativePath, name string, tags, refs []string) {
	stats := ensureImportStats(result)
	severity := normalizeNucleiSeverityLabel(doc.Info.Severity)
	stats.BySeverity[severity]++
	stats.ByType[nucleiTypeLabel(mapNucleiTemplateToType(relativePath, tags, name, refs))]++
}

func ImportNucleiTemplatesFromUpload(src io.Reader, filename string) (*VulkitImportResult, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("璇诲彇涓婁紶鏂囦欢澶辫触: %v", err)
	}
	cleaned, err := cleanupLegacyNucleiWorkflowRows()
	if err != nil {
		return nil, fmt.Errorf("娓呯悊鍘嗗彶 workflow 妯℃澘澶辫触: %v", err)
	}

	lowerName := strings.ToLower(strings.TrimSpace(filename))
	switch {
	case strings.HasSuffix(lowerName, ".zip"):
		result, err := importNucleiTemplatesFromZip(data)
		if result != nil {
			result.Cleaned = cleaned
		}
		return result, err
	case strings.HasSuffix(lowerName, ".yaml"), strings.HasSuffix(lowerName, ".yml"):
		if err := importSingleNucleiTemplate(data, filepath.Base(filename)); err != nil {
			if skipErr, ok := err.(*nucleiSkipError); ok {
				result := &VulkitImportResult{Total: 1, Skip: 1, Cleaned: cleaned, Errors: []string{skipErr.Error()}}
				addSkipReasonStat(result, skipErr.Error())
				return result, nil
			}
			return nil, err
		}
		result := &VulkitImportResult{Total: 1, Success: 1, Cleaned: cleaned}
		if doc, ok := parseNucleiTemplateDoc(data); ok {
			addImportedTemplateStat(result, doc, filepath.Base(filename), firstNonEmpty(doc.Info.Name, doc.ID, filepath.Base(filename)), normalizeNucleiStrings(doc.Info.Tags), normalizeNucleiStrings(doc.Info.Reference))
		}
		return result, nil
	default:
		return nil, fmt.Errorf("不支持的 nuclei 模板格式: %s，请上传 .zip、.yaml 或 .yml 文件", filename)
	}
}

func cleanupLegacyNucleiWorkflowRows() (int, error) {
	workflowPathPatterns := []string{
		"%/workflows/%",
		"workflows/%",
		"%-workflow.yaml",
		"%-workflow.yml",
	}
	var cleaned int64

	libArgs := make([]interface{}, 0, 6+len(workflowPathPatterns))
	libArgs = append(libArgs, enums.VulScriptTypeNuclei, enums.VulScriptTypeYak)
	for _, pattern := range workflowPathPatterns {
		libArgs = append(libArgs, pattern)
	}
	libArgs = append(libArgs, "%-workflow", "VULKIT-%-WORKFLOW")
	libResult := mysql.DB.Exec(`DELETE FROM vul_libraries
		WHERE script_type IN (?, ?)
		  AND (
			LOWER(pocname) LIKE ?
			OR LOWER(pocname) LIKE ?
			OR LOWER(pocname) LIKE ?
			OR LOWER(pocname) LIKE ?
			OR LOWER(name) LIKE ?
			OR UPPER(vul_id) LIKE ?
		  )`, libArgs...)
	if libResult.Error != nil {
		return 0, libResult.Error
	}
	cleaned += libResult.RowsAffected

	sqliteCleaned, err := cleanupLegacyWorkflowScriptsFromScannerDB()
	if err != nil {
		return 0, err
	}
	cleaned += int64(sqliteCleaned)

	return int(cleaned), nil
}

func cleanupLegacyWorkflowScriptsFromScannerDB() (int, error) {
	db, err := sqlitemodel.GetScannerDB()
	if err != nil {
		return 0, err
	}
	result := db.Exec(`DELETE FROM vul_scripts
		WHERE type IN (?, ?)
		  AND (
			LOWER(script_name) LIKE ?
			OR LOWER(script_name) LIKE ?
			OR LOWER(script_name) LIKE ?
			OR LOWER(script_name) LIKE ?
			OR LOWER(script_name) LIKE ?
		  )`,
		enums.VulScriptTypeNuclei,
		enums.VulScriptTypeYak,
		"%/workflows/%",
		"workflows/%",
		"%-workflow.yaml",
		"%-workflow.yml",
		"%-workflow",
	)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func ImportNucleiTemplatesFromPath(rootPath string) (*VulkitImportResult, error) {
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return nil, fmt.Errorf("瀵煎叆璺緞涓嶈兘涓虹┖")
	}

	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("璇诲彇瀵煎叆璺緞澶辫触: %v", err)
	}

	if !info.IsDir() {
		f, err := os.Open(rootPath)
		if err != nil {
			return nil, fmt.Errorf("鎵撳紑瀵煎叆鏂囦欢澶辫触: %v", err)
		}
		defer f.Close()
		return ImportNucleiTemplatesFromUpload(f, filepath.Base(rootPath))
	}

	result := &VulkitImportResult{}
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", path, walkErr))
			return nil
		}
		if info == nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", path, err))
			return nil
		}

		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			relPath = filepath.Base(path)
		}
		if err := importSingleNucleiTemplate(content, relPath); err != nil {
			if skipErr, ok := err.(*nucleiSkipError); ok {
				result.Total++
				result.Skip++
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %s", relPath, skipErr.Error()))
				addSkipReasonStat(result, skipErr.Error())
				return nil
			}
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", relPath, err))
			return nil
		}
		result.Total++
		result.Success++
		if doc, ok := parseNucleiTemplateDoc(content); ok {
			addImportedTemplateStat(result, doc, relPath, firstNonEmpty(doc.Info.Name, doc.ID, filepath.Base(relPath)), normalizeNucleiStrings(doc.Info.Tags), normalizeNucleiStrings(doc.Info.Reference))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if result.Total == 0 {
		return nil, fmt.Errorf("鐩綍涓湭鎵惧埌鍙鍏ョ殑 .yaml/.yml nuclei 妯℃澘")
	}
	return result, nil
}

func importNucleiTemplatesFromZip(data []byte) (*VulkitImportResult, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("瑙ｅ帇 nuclei zip 澶辫触: %v", err)
	}

	result := &VulkitImportResult{}
	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name))
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", f.Name, err))
			continue
		}
		content, readErr := io.ReadAll(rc)
		_ = rc.Close()
		if readErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", f.Name, readErr))
			continue
		}
		if err := importSingleNucleiTemplate(content, f.Name); err != nil {
			if skipErr, ok := err.(*nucleiSkipError); ok {
				result.Total++
				result.Skip++
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %s", f.Name, skipErr.Error()))
				addSkipReasonStat(result, skipErr.Error())
				continue
			}
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", f.Name, err))
			continue
		}
		result.Total++
		result.Success++
		if doc, ok := parseNucleiTemplateDoc(content); ok {
			addImportedTemplateStat(result, doc, f.Name, firstNonEmpty(doc.Info.Name, doc.ID, filepath.Base(f.Name)), normalizeNucleiStrings(doc.Info.Tags), normalizeNucleiStrings(doc.Info.Reference))
		}
	}

	if result.Total == 0 {
		return nil, fmt.Errorf("zip 涓湭鎵惧埌鍙鍏ョ殑 .yaml/.yml nuclei 妯℃澘")
	}
	return result, nil
}

func looksLikeNucleiTemplate(content []byte) bool {
	doc, ok := parseNucleiTemplateDoc(content)
	return ok && strings.TrimSpace(doc.ID) != "" && strings.TrimSpace(doc.Info.Name) != ""
}

func parseNucleiTemplateDoc(content []byte) (nucleiTemplateDoc, bool) {
	var doc nucleiTemplateDoc
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return doc, false
	}
	return doc, true
}

func importSingleNucleiTemplate(content []byte, relativePath string) error {
	var doc nucleiTemplateDoc
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return fmt.Errorf("解析 nuclei 模板失败: %v", err)
	}
	if strings.TrimSpace(doc.ID) == "" && strings.TrimSpace(doc.Info.Name) == "" {
		return fmt.Errorf("不是有效的 nuclei 模板")
	}

	pocname := normalizeNucleiPocname(relativePath, doc.ID)
	name := firstNonEmpty(doc.Info.Name, doc.ID, filepath.Base(pocname))
	tags := normalizeNucleiStrings(doc.Info.Tags)
	refs := normalizeNucleiStrings(doc.Info.Reference)
	if reason, skip := shouldSkipNucleiTemplate(doc, relativePath, name, tags, refs); skip {
		return &nucleiSkipError{Reason: reason}
	}
	cve := firstNonEmpty(normalizeNucleiStrings(doc.Info.Classification.CVEID)...)
	if cve == "" {
		cve = nucleiCveRegexp.FindString(string(content))
	}
	vulID := buildNucleiVulID(pocname, cve)
	localizedName, localizedDescription, localizedFix := localizeNucleiTemplateFields(pocname, doc, tags, refs)

	risk := mapNucleiSeverity(strings.TrimSpace(doc.Info.Severity))
	vulType := mapNucleiTemplateToType(pocname, tags, name, refs)
	vulClass := mapNucleiTemplateToClass(pocname, tags, name)
	now := time.Now()

	result := mysql.DB.Exec(`INSERT INTO vul_libraries (data_type, vul_id, name, cve, risk, type, class, description, affect_range, fix_suggest, component, status, pocname, verify_type, exploit_impact, script_type, poc_or_exp, cnvd, cnnvd, cvss_score, create_time, update_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE name=VALUES(name), cve=VALUES(cve), risk=VALUES(risk), type=VALUES(type), class=VALUES(class),
		description=VALUES(description), affect_range=VALUES(affect_range), fix_suggest=VALUES(fix_suggest), component=VALUES(component),
		status=VALUES(status), verify_type=VALUES(verify_type), exploit_impact=VALUES(exploit_impact), script_type=VALUES(script_type),
		poc_or_exp=VALUES(poc_or_exp), cnvd=VALUES(cnvd), cnnvd=VALUES(cnnvd), cvss_score=VALUES(cvss_score), update_time=VALUES(update_time)`,
		1, vulID, localizedName, cve, risk, vulType, vulClass, localizedDescription,
		deriveNucleiAffectRange(pocname, tags), localizedFix, strings.Join(tags, ","),
		enums.VulLibrariesStatusSucess, pocname, 1, mapNucleiExploitImpact(vulType, tags), enums.VulScriptTypeNuclei,
		enums.VulScriptVerifyTypePoc, strings.TrimSpace(doc.Info.Classification.CNVD), strings.TrimSpace(doc.Info.Classification.CNNVD),
		strings.TrimSpace(doc.Info.Classification.CVSSScore), now, now,
	)
	if result.Error != nil {
		return fmt.Errorf("鍐欏叆 vul_libraries 澶辫触: %v", result.Error)
	}

	if err := saveNucleiScriptToScannerDB(pocname, localizedName, string(content), vulID, now); err != nil {
		return fmt.Errorf("写入 scanner.db 的 vul_scripts 失败: %v", err)
	}

	return nil
}

func buildNucleiVulID(pocname, cve string) string {
	cve = strings.TrimSpace(cve)
	if cve != "" && len(cve) <= 50 {
		return cve
	}

	base := strings.ToUpper(strings.ReplaceAll(sanitizeForPocname(pocname), ".", "-"))
	if base == "" {
		base = "TEMPLATE"
	}
	candidate := "NUCLEI-" + base
	if len(candidate) <= 50 {
		return candidate
	}

	sum := sha1.Sum([]byte(pocname))
	hashSuffix := strings.ToUpper(hex.EncodeToString(sum[:4]))
	maxBaseLen := 50 - len("NUCLEI-") - 1 - len(hashSuffix)
	if maxBaseLen < 1 {
		maxBaseLen = 1
	}
	if len(base) > maxBaseLen {
		base = base[:maxBaseLen]
	}
	return "NUCLEI-" + base + "-" + hashSuffix
}

func saveNucleiScriptToScannerDB(scriptName, libName, content, vulID string, now time.Time) error {
	db, err := sqlitemodel.GetScannerDB()
	if err != nil {
		return err
	}
	if err := db.Where("script_name = ?", scriptName).Delete(&sqlitemodel.VulScript{}).Error; err != nil {
		return err
	}
	record := sqlitemodel.VulScript{
		UserID:       0,
		ScriptName:   scriptName,
		Type:         enums.VulScriptTypeNuclei,
		LibName:      libName,
		Content:      content,
		VulID:        vulID,
		VerifyType:   enums.VulScriptVerifyTypePoc,
		CreateTime:   now.Format("2006-01-02 15:04:05"),
		UpdateTime:   now.Format("2006-01-02 15:04:05"),
		EvidenceType: "0",
	}
	return db.Create(&record).Error
}

func normalizeNucleiStrings(v interface{}) []string {
	switch t := v.(type) {
	case nil:
		return nil
	case string:
		parts := strings.Split(t, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
		return out
	case []interface{}:
		out := make([]string, 0, len(t))
		for _, item := range t {
			out = append(out, normalizeNucleiStrings(item)...)
		}
		return out
	case []string:
		out := make([]string, 0, len(t))
		for _, item := range t {
			item = strings.TrimSpace(item)
			if item != "" {
				out = append(out, item)
			}
		}
		return out
	default:
		raw, _ := json.Marshal(t)
		var out []string
		_ = json.Unmarshal(raw, &out)
		return out
	}
}

func normalizeNucleiPocname(relativePath, templateID string) string {
	relativePath = filepath.ToSlash(strings.TrimSpace(relativePath))
	relativePath = strings.TrimPrefix(relativePath, "./")
	if relativePath != "" {
		base := filepath.Base(relativePath)
		base = strings.TrimSuffix(base, ".yaml")
		base = strings.TrimSuffix(base, ".yml")
		return base
	}
	templateID = strings.TrimSpace(templateID)
	if templateID == "" {
		templateID = fmt.Sprintf("nuclei-%d", time.Now().UnixNano())
	}
	templateID = strings.TrimSuffix(templateID, ".yaml")
	templateID = strings.TrimSuffix(templateID, ".yml")
	return strings.ReplaceAll(templateID, "\\", "/")
}

func mapNucleiSeverity(severity string) int {
	if risk, ok := nucleiSeverityToRisk[strings.ToLower(strings.TrimSpace(severity))]; ok {
		return risk
	}
	return enums.VulLibrariesRiskInfo
}

func normalizeNucleiSeverityLabel(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "critical":
		return "critical"
	case "high":
		return "high"
	case "medium":
		return "medium"
	case "low":
		return "low"
	case "info":
		return "info"
	default:
		return "unknown"
	}
}

func nucleiTypeLabel(vulType int) string {
	switch vulType {
	case enums.VulLibrariesTypeSqlInj:
		return "sql-injection"
	case enums.VulLibrariesTypeXss:
		return "xss"
	case enums.VulLibrariesTypeCsrf:
		return "csrf"
	case enums.VulLibrariesTypeCodeExec:
		return "rce"
	case enums.VulLibrariesTypeSsrf:
		return "ssrf"
	case enums.VulLibrariesTypeLfi:
		return "lfi"
	case enums.VulLibrariesTypeRfi:
		return "rfi"
	case enums.VulLibrariesTypeFileRead:
		return "file-read"
	case enums.VulLibrariesTypeFileUpload:
		return "file-upload"
	case enums.VulLibrariesTypePathTraversal:
		return "path-traversal"
	case enums.VulLibrariesTypeWeakPass:
		return "weak-password"
	case enums.VulLibrariesTypeUnauthAccess:
		return "unauthorized"
	case enums.VulLibrariesTypeInfoDisclosure:
		return "info-disclosure"
	case enums.VulLibrariesTypeDos:
		return "dos"
	default:
		return "other"
	}
}

func shouldSkipNucleiTemplate(doc nucleiTemplateDoc, relativePath, name string, tags, refs []string) (string, bool) {
	severity := strings.ToLower(strings.TrimSpace(doc.Info.Severity))
	pathLower := strings.ToLower(filepath.ToSlash(strings.TrimSpace(relativePath)))

	if severity == "" || severity == "info" || severity == "unknown" {
		return "已跳过：信息收集/未知严重级别模板", true
	}
	if isNucleiWorkflowTemplate(doc, pathLower) {
		return "已跳过：工作流/辅助模板", true
	}
	if strings.Contains(pathLower, "/helpers/") || strings.HasPrefix(pathLower, "helpers/") || pathLower == "helpers" ||
		strings.HasSuffix(pathLower, "-helper.yaml") || strings.HasSuffix(pathLower, "-helper.yml") ||
		strings.HasSuffix(strings.ToLower(strings.TrimSpace(doc.ID)), "-helper") {
		return "已跳过：工作流/辅助模板", true
	}
	return "", false
}

func isNucleiWorkflowTemplate(doc nucleiTemplateDoc, pathLower string) bool {
	idLower := strings.ToLower(strings.TrimSpace(doc.ID))
	if doc.Workflows != nil || doc.Flow != nil {
		return true
	}
	if strings.Contains(pathLower, "/workflows/") || strings.HasPrefix(pathLower, "workflows/") || pathLower == "workflows" {
		return true
	}
	if strings.HasSuffix(pathLower, "-workflow.yaml") || strings.HasSuffix(pathLower, "-workflow.yml") {
		return true
	}
	return strings.HasSuffix(idLower, "-workflow")
}

func mapNucleiTemplateToClass(pocname string, tags []string, name string) int {
	lower := strings.ToLower(strings.Join(append([]string{pocname, name}, tags...), ","))
	switch {
	case containsAny(lower, "redis", "mongodb", "mysql", "postgres", "postgresql", "elasticsearch", "couchdb"):
		return enums.VulLibrariesClassDatabase
	case containsAny(lower, "struts", "spring", "springboot", "thinkphp", "django", "flask", "shiro", "fastjson"):
		return enums.VulLibrariesClassWebFrameWork
	case containsAny(lower, "weblogic", "jboss", "tomcat", "nginx", "apache", "jenkins", "gitlab", "confluence", "jira"):
		return enums.VulLibrariesClassAppServer
	case containsAny(lower, "wordpress", "drupal", "joomla", "discuz", "ecshop", "dedecms"):
		return enums.VulLibrariesClassWebCmsSystem
	case containsAny(lower, "router", "switch", "firewall", "vpn", "fortinet", "f5", "cisco", "juniper", "h3c", "huawei"):
		return enums.VulLibrariesClassNetDriver
	case containsAny(lower, "windows", "linux", "unix", "ubuntu", "centos"):
		return enums.VulLibrariesClassSystem
	default:
		return enums.VulLibrariesClassWeb
	}
}

func mapNucleiTemplateToType(pocname string, tags []string, name string, refs []string) int {
	lower := strings.ToLower(strings.Join(append([]string{pocname, name}, append(tags, refs...)...), ","))
	switch {
	case containsAny(lower, "sql", "sqli", "sql-injection", "sql_injection"):
		return enums.VulLibrariesTypeSqlInj
	case containsAny(lower, "xss", "uxss"):
		return enums.VulLibrariesTypeXss
	case containsAny(lower, "csrf"):
		return enums.VulLibrariesTypeCsrf
	case containsAny(lower, "rce", "remote-code-execution", "code-exec", "command-exec", "cmd-injection", "command injection"):
		return enums.VulLibrariesTypeCodeExec
	case containsAny(lower, "ssrf"):
		return enums.VulLibrariesTypeSsrf
	case containsAny(lower, "lfi", "local-file-include"):
		return enums.VulLibrariesTypeLfi
	case containsAny(lower, "rfi", "remote-file-include"):
		return enums.VulLibrariesTypeRfi
	case containsAny(lower, "file-read", "arbitrary-file-read", "download", "readfile"):
		return enums.VulLibrariesTypeFileRead
	case containsAny(lower, "upload", "file-upload"):
		return enums.VulLibrariesTypeFileUpload
	case containsAny(lower, "path-traversal", "dir-traversal", "directory-traversal", "traversal"):
		return enums.VulLibrariesTypePathTraversal
	case containsAny(lower, "default-login", "weak-password", "weakpass", "password", "login"):
		return enums.VulLibrariesTypeWeakPass
	case containsAny(lower, "unauth", "unauthorized", "auth-bypass", "authentication-bypass"):
		return enums.VulLibrariesTypeUnauthAccess
	case containsAny(lower, "exposure", "disclosure", "leak", "panel", "config", "backup"):
		return enums.VulLibrariesTypeInfoDisclosure
	case containsAny(lower, "dos", "denial-of-service"):
		return enums.VulLibrariesTypeDos
	default:
		return enums.VulLibrariesTypeOther
	}
}

func mapNucleiExploitImpact(vulType int, tags []string) int {
	if vulType == enums.VulLibrariesTypeDos || containsAny(strings.ToLower(strings.Join(tags, ",")), "dos", "denial-of-service") {
		return enums.VulScriptExploitImpactRefuseServer
	}
	return enums.VulScriptExploitImpactNoImpact
}

func deriveNucleiAffectRange(pocname string, tags []string) string {
	candidates := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" || containsAny(strings.ToLower(tag), "http", "https", "network", "tcp", "ssl", "tls", "misc", "workflow", "fuzz") {
			continue
		}
		candidates = append(candidates, tag)
	}
	if len(candidates) == 0 {
		for _, part := range strings.Split(strings.ToLower(filepath.ToSlash(pocname)), "/") {
			if part == "" || containsAny(part, "http", "network", "cves", "vulnerabilities", "misconfiguration", "default-logins", "workflows") {
				continue
			}
			candidates = append(candidates, part)
		}
	}
	candidates = uniqueStrings(candidates)
	if len(candidates) > 6 {
		candidates = candidates[:6]
	}
	return strings.Join(candidates, ", ")
}

func containsAny(s string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(s, keyword) {
			return true
		}
	}
	return false
}

func uniqueStrings(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" || slices.Contains(out, item) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func firstNonEmpty(items ...string) string {
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			return item
		}
	}
	return ""
}
