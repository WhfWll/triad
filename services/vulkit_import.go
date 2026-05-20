package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

var (
	vulkitCveRegexp      = regexp.MustCompile(`CVE-\d{4}-\d+`)
	vulkitSeverityToRisk = map[string]int{
		"critical": 1,
		"high":     2,
		"medium":   3,
		"low":      4,
		"info":     5,
	}
	vulkitRiskTitleRegexp    = regexp.MustCompile(`risk\.title\(\s*"([^"]*)"\s*\)`)
	vulkitRiskSevRegexp      = regexp.MustCompile(`risk\.severity\(\s*"([^"]*)"\s*\)`)
	vulkitRiskDescRegexp     = regexp.MustCompile(`risk\.description\(\s*"([^"]*)"\s*\)`)
	vulkitRiskSolutionRegexp = regexp.MustCompile(`risk\.solution\(\s*"([^"]*)"\s*\)`)
	vulkitRiskTypeRegexp     = regexp.MustCompile(`risk\.type\(\s*"([^"]*)"\s*\)`)
	vulkitCVEDetailRegexp    = regexp.MustCompile(`"cve"\s*:\s*"([^"]*)"`)
	vulkitRangeDetailRegexp  = regexp.MustCompile(`"影响范围"\s*:\s*"([^"]*)"`)
	vulkitCompDetailRegexp   = regexp.MustCompile(`"影响组件"\s*:\s*"([^"]*)"`)
)

type VulkitImportResult struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Skip    int      `json:"skip"`
	Errors  []string `json:"errors"`
}

func ImportVulnerabilitiesFromUpload(src io.Reader, filename string) (*VulkitImportResult, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %v", err)
	}

	lowerName := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lowerName, ".zip"):
		return importVulnsFromZip(data)
	case strings.HasSuffix(lowerName, ".json"):
		return importVulnsFromJson(data)
	case strings.HasSuffix(lowerName, ".yak"):
		scriptName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		if err := importSingleVuln(data, scriptName, filename); err != nil {
			return nil, err
		}
		return &VulkitImportResult{Total: 1, Success: 1}, nil
	default:
		return nil, fmt.Errorf("不支持的文件格式: %s，请上传 .zip、.json 或 .yak 文件", filename)
	}
}

// VulnExportItem 对应 webScanner 导出的 vulns_export.json 单条记录
type VulnExportItem struct {
	Pocname          string `json:"pocname"`
	Name             string `json:"name"`
	VulId            string `json:"VulId"`
	VulIdLower       string `json:"vulId"`
	Risk             int    `json:"risk"`
	Type             int    `json:"type"`
	Class            int    `json:"class"`
	OperatingSystem  int    `json:"operatingSystem"`
	Description      string `json:"description"`
	FixSuggest       string `json:"fixSuggest"`
	AffectRange      string `json:"affectRange"`
	CVE              string `json:"cve"`
	CNVD             string `json:"cnvd"`
	CNNVD            string `json:"cnnvd"`
	VerifyType       string `json:"verifyType"`
	ExploitImpact    int    `json:"exploitImpact"`
	ScriptType       string `json:"scriptType"`
}

func importVulnsFromJson(data []byte) (*VulkitImportResult, error) {
	var items []VulnExportItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("JSON 中未包含漏洞记录")
	}

	result := &VulkitImportResult{}
	for _, item := range items {
		if item.Pocname == "" {
			continue
		}

		now := time.Now()
		exploitImpact := item.ExploitImpact
		if exploitImpact == 0 {
			exploitImpact = 2
		}
		cve := item.CVE
		vulID := item.VulId
		if vulID == "" {
			vulID = item.VulIdLower
		}
		if vulID == "" {
			vulID = cve
		}
		if vulID == "" {
			vulID = "BUILTIN-" + strings.ToUpper(strings.ReplaceAll(item.Pocname, " ", "-"))
		}
		risk := item.Risk
		if risk == 0 {
			risk = 3
		}
		vulType := item.Type
		if vulType == 0 {
			vulType = 74
		}
		vulClass := item.Class
		if vulClass == 0 {
			vulClass = 1
		}
		scriptType := item.ScriptType
		if scriptType == "" {
			scriptType = "universal"
		}
		verifyType := 1
		if item.VerifyType != "" {
			if v, err := strconv.Atoi(item.VerifyType); err == nil && v > 0 {
				verifyType = v
			}
		}

		dbResult := mysql.DB.Exec(`INSERT INTO vul_libraries (data_type, vul_id, name, cve, risk, type, class, description, affect_range, fix_suggest, status, pocname, verify_type, exploit_impact, script_type, cnvd, cnnvd, operating_system, create_time, update_time)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), cve=VALUES(cve), risk=VALUES(risk), type=VALUES(type), class=VALUES(class),
			description=VALUES(description), affect_range=VALUES(affect_range), fix_suggest=VALUES(fix_suggest), status=VALUES(status),
			verify_type=VALUES(verify_type), exploit_impact=VALUES(exploit_impact), script_type=VALUES(script_type),
			cnvd=VALUES(cnvd), cnnvd=VALUES(cnnvd), operating_system=VALUES(operating_system), update_time=VALUES(update_time)`,
			1, vulID, item.Name, cve, risk, vulType, vulClass, item.Description, item.AffectRange, item.FixSuggest,
			2, item.Pocname, verifyType, exploitImpact, scriptType, item.CNVD, item.CNNVD, item.OperatingSystem, now, now,
		)
		if dbResult.Error != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", item.Pocname, dbResult.Error))
			continue
		}
		result.Success++
		result.Total++
	}

	if result.Total == 0 {
		if len(result.Errors) > 0 {
			return nil, fmt.Errorf("导入失败: %s", result.Errors[0])
		}
		return nil, fmt.Errorf("JSON 中无有效漏洞记录（需包含 pocname 字段）")
	}
	return result, nil
}

func importVulnsFromZip(data []byte) (*VulkitImportResult, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("解压zip失败: %v", err)
	}

	result := &VulkitImportResult{}

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name))
		if ext != ".yak" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		scriptName := strings.TrimSuffix(filepath.Base(f.Name), ext)

		rc, err := f.Open()
		if err != nil {
			continue
		}
		content, _ := io.ReadAll(rc)
		rc.Close()

		err = importSingleVuln(content, scriptName, f.Name)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", scriptName, err))
			continue
		}
		result.Success++
		result.Total++
	}

	if result.Total == 0 {
		return nil, fmt.Errorf("在zip中未找到 .yak、.yaml 或 .yml 文件")
	}
	return result, nil
}

func importSingleVuln(content []byte, scriptName, path string) error {
	contentStr := string(content)
	name := vulkitExtractFirstMatch(vulkitRiskTitleRegexp, contentStr)
	if name == "" {
		name = scriptName
	}

	cleanName := sanitizeForPocname(scriptName)

	isExp := strings.HasPrefix(scriptName, "exp-")
	isPoc := strings.HasPrefix(scriptName, "poc-")
	pocOrExp := ""
	vulnVt := 1
	switch {
	case isExp:
		pocOrExp = "exp"
	case isPoc:
		pocOrExp = "poc"
	}

	vulnCategory := extractVulnCategory(path)
	vulnClass := mapVulnCategoryToClass(vulnCategory)

	sevStr := vulkitExtractFirstMatch(vulkitRiskSevRegexp, contentStr)
	riskLevel := vulkitSeverityToRisk[strings.ToLower(sevStr)]
	if riskLevel == 0 {
		riskLevel = 3
	}

	description := vulkitExtractFirstMatch(vulkitRiskDescRegexp, contentStr)
	fixSuggest := vulkitExtractFirstMatch(vulkitRiskSolutionRegexp, contentStr)
	riskType := vulkitExtractFirstMatch(vulkitRiskTypeRegexp, contentStr)
	vulType := mapRiskTypeToVulType(riskType)

	cve := vulkitExtractFirstMatch(vulkitCVEDetailRegexp, contentStr)
	if cve == "" || !vulkitCveRegexp.MatchString(cve) {
		if m := vulkitCveRegexp.FindStringSubmatch(contentStr); len(m) > 0 {
			cve = m[0]
		} else {
			cve = ""
		}
	}

	affectRange := vulkitExtractFirstMatch(vulkitRangeDetailRegexp, contentStr)
	component := vulkitExtractFirstMatch(vulkitCompDetailRegexp, contentStr)

	now := time.Now()

	vulID := cve
	if vulID == "" {
		shortName := sanitizeForPocname(scriptName)
		if len(shortName) > 40 {
			shortName = shortName[:40]
		}
		vulID = "VULKIT-" + strings.ToUpper(strings.ReplaceAll(shortName, " ", "-"))
	}

	result := mysql.DB.Exec(`INSERT INTO vul_libraries (data_type, vul_id, name, cve, risk, type, class, description, affect_range, fix_suggest, component, status, pocname, verify_type, script_type, poc_or_exp, create_time, update_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE name=VALUES(name), cve=VALUES(cve), risk=VALUES(risk), type=VALUES(type), class=VALUES(class),
		description=VALUES(description), affect_range=VALUES(affect_range), fix_suggest=VALUES(fix_suggest), component=VALUES(component),
		status=VALUES(status), verify_type=VALUES(verify_type), script_type=VALUES(script_type), poc_or_exp=VALUES(poc_or_exp), update_time=VALUES(update_time)`,
		1, vulID, name, cve, riskLevel, vulType, vulnClass, description, affectRange, fixSuggest, component, 2, cleanName, vulnVt, "yak", pocOrExp, now, now,
	)
	if result.Error != nil {
		return fmt.Errorf("写入 vul_libraries 失败: %v", result.Error)
	}

	return nil
}

func vulkitExtractFirstMatch(re *regexp.Regexp, s string) string {
	m := re.FindStringSubmatch(s)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func sanitizeForPocname(name string) string {
	var result []byte
	for _, b := range []byte(name) {
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '-' || b == '_' || b == '.' {
			result = append(result, b)
		}
	}
	s := strings.Trim(string(result), "-._")
	if s == "" {
		return "imported_" + fmt.Sprintf("%d", time.Now().Unix())
	}
	return s
}

func extractVulnCategory(path string) string {
	dir := filepath.Dir(path)
	parts := strings.Split(dir, string(filepath.Separator))
	if len(parts) > 0 {
		candidate := parts[len(parts)-1]
		if !vulkitCveRegexp.MatchString(candidate) {
			return candidate
		}
		if len(parts) > 1 {
			return parts[len(parts)-2]
		}
	}
	return ""
}

func mapVulnCategoryToClass(category string) int {
	switch strings.ToLower(category) {
	case "web", "web应用", "web应用漏洞":
		return 1
	case "cms", "oa", "办公oa", "办公应用", "企业应用":
		return 2
	case "network", "网络设备", "网络设备漏洞":
		return 3
	case "middleware", "中间件", "web框架", "框架":
		return 4
	case "database", "数据库", "大数据":
		return 5
	case "security", "安全设备":
		return 6
	case "os", "系统", "操作系统":
		return 7
	case "iot", "物联网", "工控":
		return 8
	default:
		return 1
	}
}

func mapRiskTypeToVulType(riskType string) int {
	switch strings.ToLower(riskType) {
	case "xss", "uxss":
		return 1
	case "csrf":
		return 2
	case "sqli", "sql-injection", "sql_injection":
		return 3
	case "rce", "code-exec", "code_exec":
		return 12
	case "rfi", "remote-file-include", "remote_file_include":
		return 13
	case "lfi", "local-file-include", "local_file_include":
		return 14
	case "dos", "denial-of-service":
		return 21
	case "fileread", "file-read", "file_read", "arbitrary-file-read":
		return 38
	case "path-traversal", "pathtraversal", "path_traversal", "dir-traversal", "dirtraversal":
		return 42
	case "fileupload", "file-upload", "file_upload", "arbitrary-file-upload":
		return 43
	case "bypass", "login-bypass":
		return 44
	case "weak-pass", "weakpass", "weak_password", "weak-password", "brute", "bruteforce", "brute-force":
		return 45
	case "infoleak", "info-leak", "info_leak", "information-disclosure", "information_disclosure":
		return 49
	case "backdoor", "back-door", "shell":
		return 53
	case "unauth", "unauthorized", "unauth-access":
		return 61
	case "ssrf":
		return 71
	case "command-exec", "command_exec":
		return 72
	case "overflow", "buffer-overflow", "buffer_overflow":
		return 17
	default:
		return 74
	}
}

func init() {
	log.Info("vulkit_import service loaded")
}
