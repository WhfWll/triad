package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
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
	ID   string `yaml:"id"`
	Info struct {
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

func ImportNucleiTemplatesFromUpload(src io.Reader, filename string) (*VulkitImportResult, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %v", err)
	}

	lowerName := strings.ToLower(strings.TrimSpace(filename))
	switch {
	case strings.HasSuffix(lowerName, ".zip"):
		return importNucleiTemplatesFromZip(data)
	case strings.HasSuffix(lowerName, ".yaml"), strings.HasSuffix(lowerName, ".yml"):
		if err := importSingleNucleiTemplate(data, filepath.Base(filename)); err != nil {
			return nil, err
		}
		return &VulkitImportResult{Total: 1, Success: 1}, nil
	default:
		return nil, fmt.Errorf("不支持的 nuclei 模板格式: %s，请上传 .zip、.yaml 或 .yml 文件", filename)
	}
}

func ImportNucleiTemplatesFromPath(rootPath string) (*VulkitImportResult, error) {
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return nil, fmt.Errorf("导入路径不能为空")
	}

	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("读取导入路径失败: %v", err)
	}

	if !info.IsDir() {
		f, err := os.Open(rootPath)
		if err != nil {
			return nil, fmt.Errorf("打开导入文件失败: %v", err)
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
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", relPath, err))
			return nil
		}
		result.Total++
		result.Success++
		return nil
	})
	if err != nil {
		return nil, err
	}
	if result.Total == 0 {
		return nil, fmt.Errorf("目录中未找到可导入的 .yaml/.yml nuclei 模板")
	}
	return result, nil
}

func importNucleiTemplatesFromZip(data []byte) (*VulkitImportResult, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("解压 nuclei zip 失败: %v", err)
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
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", f.Name, err))
			continue
		}
		result.Total++
		result.Success++
	}

	if result.Total == 0 {
		return nil, fmt.Errorf("zip 中未找到可导入的 .yaml/.yml nuclei 模板")
	}
	return result, nil
}

func looksLikeNucleiTemplate(content []byte) bool {
	var doc nucleiTemplateDoc
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return false
	}
	return strings.TrimSpace(doc.ID) != "" && strings.TrimSpace(doc.Info.Name) != ""
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
	cve := firstNonEmpty(normalizeNucleiStrings(doc.Info.Classification.CVEID)...)
	if cve == "" {
		cve = nucleiCveRegexp.FindString(string(content))
	}
	vulID := cve
	if vulID == "" {
		vulID = "NUCLEI-" + strings.ToUpper(strings.ReplaceAll(sanitizeForPocname(pocname), ".", "-"))
	}
	if len(vulID) > 120 {
		vulID = vulID[:120]
	}

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
		1, vulID, name, cve, risk, vulType, vulClass, strings.TrimSpace(doc.Info.Description),
		deriveNucleiAffectRange(pocname, tags), strings.TrimSpace(doc.Info.Remediation), strings.Join(tags, ","),
		enums.VulLibrariesStatusSucess, pocname, 1, mapNucleiExploitImpact(vulType, tags), enums.VulScriptTypeNuclei,
		enums.VulScriptVerifyTypePoc, strings.TrimSpace(doc.Info.Classification.CNVD), strings.TrimSpace(doc.Info.Classification.CNNVD),
		strings.TrimSpace(doc.Info.Classification.CVSSScore), now, now,
	)
	if result.Error != nil {
		return fmt.Errorf("写入 vul_libraries 失败: %v", result.Error)
	}

	contentStr := string(content)
	vulScriptResult := mysql.DB.Exec(`INSERT INTO vul_scripts (data_type, user_id, script_name, type, lib_name, content, vul_id, verify_type, create_time, update_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE script_name=VALUES(script_name), lib_name=VALUES(lib_name), content=VALUES(content), verify_type=VALUES(verify_type), update_time=VALUES(update_time)`,
		1, 0, pocname, enums.VulScriptTypeNuclei, name, contentStr, vulID, enums.VulScriptVerifyTypePoc, now, now,
	)
	if vulScriptResult.Error != nil {
		return fmt.Errorf("写入 vul_scripts 失败: %v", vulScriptResult.Error)
	}

	return nil
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
		return relativePath
	}
	templateID = strings.TrimSpace(templateID)
	if templateID == "" {
		templateID = fmt.Sprintf("nuclei-%d", time.Now().UnixNano())
	}
	if !strings.HasSuffix(strings.ToLower(templateID), ".yaml") && !strings.HasSuffix(strings.ToLower(templateID), ".yml") {
		templateID += ".yaml"
	}
	return strings.ReplaceAll(templateID, "\\", "/")
}

func mapNucleiSeverity(severity string) int {
	if risk, ok := nucleiSeverityToRisk[strings.ToLower(strings.TrimSpace(severity))]; ok {
		return risk
	}
	return enums.VulLibrariesRiskInfo
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
