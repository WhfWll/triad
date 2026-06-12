package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	sqlitemodel "smart/models/sqlite"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/yaml.v2"
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
	Total   int              `json:"total"`
	Success int              `json:"success"`
	Skip    int              `json:"skip"`
	Cleaned int              `json:"cleaned,omitempty"`
	Errors  []string         `json:"errors"`
	Message string           `json:"message,omitempty"`
	Stats   *ImportRuleStats `json:"stats,omitempty"`
}

type ImportRuleStats struct {
	BySeverity map[string]int `json:"bySeverity,omitempty"`
	ByType     map[string]int `json:"byType,omitempty"`
	SkipReason map[string]int `json:"skipReason,omitempty"`
}

func ImportVulnerabilitiesFromUpload(src io.Reader, filename string) (*VulkitImportResult, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %v", err)
	}

	lowerName := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lowerName, ".zip"):
		cls := classifyZipContents(data)
		switch {
		case cls.NucleiYaml > 0 && cls.YakCount == 0 && cls.VulkitYaml == 0:
			return ImportNucleiTemplatesFromUpload(bytes.NewReader(data), filename)
		case cls.YakCount > 0 && cls.NucleiYaml == 0 && cls.VulkitYaml == 0:
			return importVulnsFromZip(data)
		default:
			return importMixedZip(data)
		}
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

type zipClassification struct {
	YakCount   int
	VulkitYaml int
	NucleiYaml int
}

func classifyZipContents(data []byte) zipClassification {
	var cls zipClassification
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return cls
	}
	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := strings.ToLower(filepath.ToSlash(f.Name))
		ext := strings.ToLower(filepath.Ext(f.Name))
		switch ext {
		case ".yak":
			cls.YakCount++
		case ".yaml", ".yml":
			if strings.Contains(name, "/workflows/") || strings.HasPrefix(name, "workflows/") ||
				strings.Contains(name, "/helpers/") || strings.HasPrefix(name, "helpers/") {
				cls.NucleiYaml++
				continue
			}
			rc, err := f.Open()
			if err != nil {
				cls.VulkitYaml++
				continue
			}
			content, readErr := io.ReadAll(rc)
			rc.Close()
			if readErr != nil {
				cls.VulkitYaml++
				continue
			}
			if looksLikeNucleiTemplate(content) {
				cls.NucleiYaml++
			} else {
				cls.VulkitYaml++
			}
		}
	}
	return cls
}

func importMixedZip(data []byte) (*VulkitImportResult, error) {
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

		rc, err := f.Open()
		if err != nil {
			continue
		}
		content, readErr := io.ReadAll(rc)
		rc.Close()
		if readErr != nil {
			continue
		}

		switch ext {
		case ".yak":
			scriptName := strings.TrimSuffix(filepath.Base(f.Name), ext)
			err = importSingleVuln(content, scriptName, f.Name)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", scriptName, err))
				continue
			}
			result.Success++
			result.Total++
		case ".yaml", ".yml":
			if strings.HasSuffix(f.Name, "-workflow.yaml") || strings.HasSuffix(f.Name, "-workflow.yml") {
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
	}

	if result.Total == 0 {
		if len(result.Errors) > 0 {
			return nil, fmt.Errorf("导入失败: %s", result.Errors[0])
		}
		return nil, fmt.Errorf("zip中未找到可导入的脚本")
	}
	return result, nil
}

func shouldTreatZipAsNuclei(data []byte) bool {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return false
	}

	yakCount := 0
	nucleiTemplateCount := 0
	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := strings.ToLower(filepath.ToSlash(f.Name))
		switch filepath.Ext(name) {
		case ".yak":
			yakCount++
		case ".yaml", ".yml":
			if strings.Contains(name, "/workflows/") || strings.HasPrefix(name, "workflows/") || strings.HasSuffix(name, "-workflow.yaml") || strings.HasSuffix(name, "-workflow.yml") {
				return true
			}
			rc, err := f.Open()
			if err != nil {
				continue
			}
			content, readErr := io.ReadAll(rc)
			_ = rc.Close()
			if readErr != nil {
				continue
			}
			if looksLikeNucleiTemplate(content) {
				nucleiTemplateCount++
			}
		}
	}

	return yakCount == 0 && nucleiTemplateCount > 0
}

// VulnExportItem 对应 webScanner 导出的 vulns_export.json 单条记录
type VulnExportItem struct {
	Pocname         string `json:"pocname"`
	Name            string `json:"name"`
	VulId           string `json:"VulId"`
	VulIdLower      string `json:"vulId"`
	Risk            int    `json:"risk"`
	Type            int    `json:"type"`
	Class           int    `json:"class"`
	OperatingSystem int    `json:"operatingSystem"`
	Description     string `json:"description"`
	FixSuggest      string `json:"fixSuggest"`
	AffectRange     string `json:"affectRange"`
	CVE             string `json:"cve"`
	CNVD            string `json:"cnvd"`
	CNNVD           string `json:"cnnvd"`
	VerifyType      string `json:"verifyType"`
	ExploitImpact   int    `json:"exploitImpact"`
	ScriptType      string `json:"scriptType"`
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

	// 将脚本内容写入 scanner.db 的 vul_scripts 表（供脚本测试功能使用）
	if err := saveYakScriptToScannerDB(cleanName, name, string(content), vulID, now); err != nil {
		log.Warnf("写入 scanner.db 的 vul_scripts 失败: %v (不影响主库导入)", err)
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

func saveYakScriptToScannerDB(scriptName, libName, content, vulID string, now time.Time) error {
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
		Type:         enums.VulScriptTypeYak,
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

type VulKitFingerprint struct {
	Cms      string   `json:"cms"`
	Keyword  []string `json:"keyword"`
	Location string   `json:"location"`
	Method   string   `json:"method"`
}

type FingerprintRuleMethod struct {
	Keywords []KeywordMatcher `yaml:"keywords,omitempty"`
}

type KeywordMatcher struct {
	Product string `yaml:"product"`
	Regexp  string `yaml:"regexp"`
}

var categoryToAppClass = map[string]int{
	"中间件": enums.FingerClassMiddleware, "Web框架": enums.FingerClassMiddleware, "主流应用框架": enums.FingerClassMainAppFramework,
	"CMS": enums.FingerClassMainCMS, "开源程序": enums.FingerClassMainCMS, "电商平台": enums.FingerClassECommerceSystem, "GIS平台": enums.FingerClassAnalysisSystem,
	"办公oa": enums.FingerClassCollaborativeOfficeSuite, "OA": enums.FingerClassCollaborativeOfficeSuite, "企业管理系统": enums.FingerClassCollaborativeOfficeSuite, "协作平台": enums.FingerClassCollaborativeOfficeSuite,
	"ERP":      enums.FingerClassCollaborativeOfficeSuite,
	"数据库大数据组件": enums.FingerClassDatabaseSer,
	"大数据平台":    enums.FingerClassBigDataPlatform, "对象存储": enums.FingerClassBigDataPlatform, "数据分析": enums.FingerClassBigDataPlatform, "数据治理": enums.FingerClassBigDataPlatform,
	"BI平台": enums.FingerClassBigDataPlatform,
	"网络设备": enums.FingerClassNetworkEquipment, "边界网络设备": enums.FingerClassNetworkEquipment, "负载均衡": enums.FingerClassLoadBalancer,
	"网关": enums.FingerClassNetworkEquipment, "NAS": enums.FingerClassNetworkStorage,
	"防火墙": enums.FingerClassSafeEquipment, "VPN": enums.FingerClassSafeEquipment, "VPN网关": enums.FingerClassSafeEquipment, "网络安全": enums.FingerClassSafeEquipment,
	"身份认证": enums.FingerClassSafeEquipment,
	"开发工具": enums.FingerClassCodeDevelopment, "代码托管": enums.FingerClassCodeDevelopment, "低代码平台": enums.FingerClassCodeDevelopment, "开发平台": enums.FingerClassCodeDevelopment,
	"自动化工具": enums.FingerClassCodeDevelopment, "主机面板": enums.FingerClassCodeDevelopment, "ML平台": enums.FingerClassCodeDevelopment, "CI-CD": enums.FingerClassBuildCISystem,
	"AI推理引擎": enums.FingerClassCodeDevelopment, "AI工具链": enums.FingerClassCodeDevelopment,
	"视频监控": enums.FingerClassVideoSurveillance, "监控设备": enums.FingerClassVideoSurveillance, "监控平台": enums.FingerClassVideoSurveillance,
	"邮件系统": enums.FingerClassMailServer,
	"虚拟化":  enums.FingerClassVirtualSer, "云平台虚拟化DevOps": enums.FingerClassVirtualSer, "云存储": enums.FingerClassVirtualSer, "云平台": enums.FingerClassCloudPlatform,
	"服务器管理": enums.FingerClassHostOperatingSystem,
	"即时通讯":  enums.FingerClassOnlineChat, "备份系统": enums.FingerClassOther, "打印管理": enums.FingerClassOfficeEquipment, "文件传输": enums.FingerClassApplicationService,
	"文件共享": enums.FingerClassApplicationService, "服务端口": enums.FingerClassApplicationService, "移动管理": enums.FingerClassApplicationService, "终端管理": enums.FingerClassApplicationService,
	"远程管理": enums.FingerClassRemoteAccessSystem, "IT服务管理": enums.FingerClassApplicationService,
}

var nameAppClassMap = map[string]int{
	"nginx": enums.FingerClassMiddleware, "apache": enums.FingerClassMiddleware, "tomcat": enums.FingerClassMiddleware, "jetty": enums.FingerClassMiddleware,
	"jboss": enums.FingerClassMiddleware, "weblogic": enums.FingerClassMiddleware, "websphere": enums.FingerClassMiddleware, "iis": enums.FingerClassMiddleware,
	"caddy": enums.FingerClassMiddleware, "goahead": enums.FingerClassMiddleware, "tongweb": enums.FingerClassMiddleware, "tengine": enums.FingerClassMiddleware,
	"kafka": enums.FingerClassBigDataPlatform, "flink": enums.FingerClassBigDataPlatform, "apisix": enums.FingerClassMiddleware, "shenyu": enums.FingerClassMiddleware,
	"consul": enums.FingerClassApplicationService, "minio": enums.FingerClassBigDataPlatform,
	"wordpress": enums.FingerClassMainCMS, "joomla": enums.FingerClassMainCMS, "drupal": enums.FingerClassMainCMS, "dedecms": enums.FingerClassMainCMS,
	"phpcms": enums.FingerClassMainCMS, "discuz": enums.FingerClassMainCMS, "ecshop": enums.FingerClassECommerceSystem, "phpmyadmin": enums.FingerClassDatabase,
	"metinfo": enums.FingerClassMainCMS, "empirecms": enums.FingerClassMainCMS, "thinkcmf": enums.FingerClassMainCMS, "74cms": enums.FingerClassMainCMS,
	"zzcms": enums.FingerClassMainCMS, "doccms": enums.FingerClassMainCMS, "confluence": enums.FingerClassMainCMS, "jira": enums.FingerClassIssueTracker,
	"seeyon": enums.FingerClassCollaborativeOfficeSuite, "tongda": enums.FingerClassCollaborativeOfficeSuite, "yonyou": enums.FingerClassCollaborativeOfficeSuite,
	"oa": enums.FingerClassOfficeAuto, "finereport": enums.FingerClassAnalysisSystem, "dzzoffice": enums.FingerClassOfficeAuto, "glpi": enums.FingerClassApplicationService,
	"dolibarr": enums.FingerClassCollaborativeOfficeSuite,
	"cisco":    enums.FingerClassNetworkEquipment, "huawei": enums.FingerClassNetworkEquipment, "h3c": enums.FingerClassNetworkEquipment,
	"juniper": enums.FingerClassNetworkEquipment, "zyxel": enums.FingerClassNetworkEquipment, "router": enums.FingerClassNetworkEquipment,
	"switch": enums.FingerClassNetworkEquipment, "ruijie": enums.FingerClassNetworkEquipment, "fortinet": enums.FingerClassSafeEquipment, "forti": enums.FingerClassSafeEquipment,
	"firewall": enums.FingerClassSafeEquipment, "防火墙": enums.FingerClassSafeEquipment, "ips": enums.FingerClassSafeEquipment, "ids": enums.FingerClassSafeEquipment,
	"waf": enums.FingerClassSafeEquipment, "深信服": enums.FingerClassSafeEquipment, "绿盟": enums.FingerClassSafeEquipment, "安恒": enums.FingerClassSafeEquipment,
	"山石": enums.FingerClassSafeEquipment, "天融信": enums.FingerClassSafeEquipment, "palo alto": enums.FingerClassSafeEquipment, "pan-os": enums.FingerClassSafeEquipment,
	"sonicwall": enums.FingerClassSafeEquipment,
	"mysql":     enums.FingerClassDatabaseSer, "postgresql": enums.FingerClassDatabaseSer, "oracle": enums.FingerClassDatabaseSer, "sqlserver": enums.FingerClassDatabaseSer,
	"mongodb": enums.FingerClassDatabaseSer, "redis": enums.FingerClassDatabaseSer, "elasticsearch": enums.FingerClassBigDataPlatform, "mariadb": enums.FingerClassDatabaseSer,
	"sqlite": enums.FingerClassDatabaseSer, "hive": enums.FingerClassBigDataPlatform,
	"gitlab": enums.FingerClassCodeDevelopment, "jenkins": enums.FingerClassBuildCISystem, "sonarqube": enums.FingerClassCodeDevelopment,
	"xxl-job": enums.FingerClassCodeDevelopment, "yapi": enums.FingerClassCodeDevelopment, "metersphere": enums.FingerClassCodeDevelopment,
	"langflow": enums.FingerClassCodeDevelopment, "backstage": enums.FingerClassCodeDevelopment, "appsmith": enums.FingerClassCodeDevelopment,
	"n8n": enums.FingerClassCodeDevelopment, "marimo": enums.FingerClassCodeDevelopment, "metabase": enums.FingerClassCodeDevelopment,
}

func ImportFingerprintsFromUpload(src io.Reader, filename string) (*VulkitImportResult, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %v", err)
	}

	if strings.HasSuffix(strings.ToLower(filename), ".zip") {
		return importFingerprintsFromZip(data)
	}

	if isNmapFingerprintFile(data, filename) {
		return importNmapFingerprint(data, filename)
	}

	return importFingerprintsFromSingleFile(data, filename)
}

func importFingerprintsFromZip(data []byte) (*VulkitImportResult, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("解压zip失败: %v", err)
	}

	result := &VulkitImportResult{}
	seenProducts := make(map[string]bool)

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if !strings.EqualFold(filepath.Base(f.Name), "fingerprint.json") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			continue
		}

		fpData, _ := io.ReadAll(rc)
		rc.Close()
		fpData = convertGBKToUTF8(fpData)

		productName := extractProductFromZipPath(f.Name)
		categoryName := extractCategoryFromZipPath(f.Name)

		if seenProducts[productName] {
			continue
		}
		seenProducts[productName] = true

		fingerprints := parseFingerprintData(fpData)
		if len(fingerprints) == 0 {
			continue
		}

		err = saveFingerprintToDB(productName, categoryName, fingerprints)
		if err != nil {
			continue
		}
		result.Success++
		result.Total++
	}

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if !strings.EqualFold(f.Name, "fingerprint/finger.json") && !strings.EqualFold(f.Name, "fingerprint\\finger.json") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			continue
		}
		fpData, _ := io.ReadAll(rc)
		rc.Close()

		var extraFingerprints []VulKitFingerprint
		if err := json.Unmarshal(fpData, &extraFingerprints); err != nil {
			continue
		}
		for _, fp := range extraFingerprints {
			pn := normalizeProductName(fp.Cms)
			if seenProducts[pn] {
				continue
			}
			seenProducts[pn] = true
			err := saveFingerprintToDB(pn, "", []VulKitFingerprint{fp})
			if err != nil {
				continue
			}
			result.Success++
			result.Total++
		}
	}

	if result.Total == 0 {
		return nil, fmt.Errorf("在zip中未找到fingerprint.json文件")
	}
	result.Message = fmt.Sprintf("成功导入 %d 个产品指纹", result.Success)
	return result, nil
}

func importFingerprintsFromSingleFile(data []byte, filename string) (*VulkitImportResult, error) {
	data = convertGBKToUTF8(data)
	fingerprints := parseFingerprintData(data)
	if len(fingerprints) == 0 {
		return nil, fmt.Errorf("无法解析指纹文件，请确认是VulKit格式的fingerprint.json")
	}

	grouped := make(map[string][]VulKitFingerprint)
	for _, fp := range fingerprints {
		name := normalizeProductName(fp.Cms)
		if name == "" {
			continue
		}
		grouped[name] = append(grouped[name], fp)
	}

	result := &VulkitImportResult{}
	for productName, fps := range grouped {
		err := saveFingerprintToDB(productName, "", fps)
		if err != nil {
			log.Errorf("保存指纹[%s]失败: %v", productName, err)
			result.Errors = append(result.Errors, fmt.Sprintf("保存[%s]失败: %v", productName, err))
			continue
		}
		result.Success++
		result.Total++
	}

	result.Message = fmt.Sprintf("成功导入 %d 个产品指纹", result.Success)
	return result, nil
}

func saveFingerprintToDB(productName, categoryName string, fingerprints []VulKitFingerprint) error {
	var rules []FingerprintRuleMethod
	for _, fp := range fingerprints {
		var keywords []KeywordMatcher
		for _, kw := range fp.Keyword {
			keywords = append(keywords, KeywordMatcher{
				Product: productName,
				Regexp:  kw,
			})
		}
		if len(keywords) > 0 {
			rules = append(rules, FingerprintRuleMethod{Keywords: keywords})
		}
	}
	if len(rules) == 0 {
		return nil
	}

	flagBytes, _ := yaml.Marshal(map[string]interface{}{"methods": rules})
	flagStr := string(flagBytes)

	if len(flagStr) > 1800 {
		mysql.DB.Exec("ALTER TABLE finger MODIFY COLUMN flag LONGTEXT CHARACTER SET utf8mb4 DEFAULT '' COMMENT '匹配内容 指纹规则'")
		mysql.DB.Exec("ALTER TABLE finger MODIFY COLUMN `desc` TEXT CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '指纹描述'")
	}

	appClass := categoryToAppClass[categoryName]
	if appClass == 0 {
		appClass = guessAppClass(productName)
	}

	fingerType := enums.FingerEnum.GetFingerTypeByClass(appClass)
	level := enums.FingerEnum.GetFingerLevelByClass(appClass)

	now := time.Now()
	var count int64
	mysql.DB.Table("finger").Where("app_name = ?", productName).Count(&count)

	desc := fmt.Sprintf("上传导入, 规则数: %d", len(fingerprints))
	if categoryName != "" {
		desc = fmt.Sprintf("来自VulKit, 分类: %s, 规则数: %d", categoryName, len(fingerprints))
	}

	if count > 0 {
		updates := map[string]interface{}{
			"flag":        flagStr,
			"app_class":   appClass,
			"desc":        desc,
			"update_time": now,
		}
		if level != "" {
			updates["level"] = level
		}
		return mysql.DB.Table("finger").Where("app_name = ?", productName).Updates(updates).Error
	}
	return mysql.DB.Table("finger").Create(map[string]interface{}{
		"app_name":    productName,
		"cn_name":     productName,
		"app_version": "",
		"source":      2,
		"level":       level,
		"finger_type": fingerType,
		"app_class":   appClass,
		"flag":        flagStr,
		"desc":        desc,
		"create_time": now,
		"update_time": now,
	}).Error
}

func guessAppClass(name string) int {
	lower := strings.ToLower(name)
	for kw, class := range nameAppClassMap {
		if strings.Contains(lower, kw) {
			return class
		}
	}
	return enums.FingerClassOther
}

func guessAppClassWithLine(name, line string) int {
	productName := extractNmapProductName(line)
	hints := extractNmapRuleHints(line)

	for _, candidate := range []string{productName, hints, name} {
		if candidate == "" {
			continue
		}
		if appClass := guessAppClass(candidate); appClass != enums.FingerClassOther {
			return appClass
		}
	}

	if isNmapIotFingerprint(name, productName, hints) {
		return enums.FingerClassIOT
	}

	return enums.FingerClassOther
}

func parseFingerprintData(data []byte) []VulKitFingerprint {
	var arr []VulKitFingerprint
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr
	}
	var single VulKitFingerprint
	if err := json.Unmarshal(data, &single); err == nil && single.Cms != "" {
		return []VulKitFingerprint{single}
	}
	return nil
}

func convertGBKToUTF8(data []byte) []byte {
	if utf8.Valid(data) {
		return data
	}
	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Data, err := decoder.Bytes(data)
	if err == nil && utf8.Valid(utf8Data) {
		return utf8Data
	}
	return data
}

func normalizeProductName(name string) string {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return name
	}
	return name
}

func extractProductFromZipPath(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i, part := range parts {
		if part == "product" && i+2 < len(parts) {
			return normalizeProductName(parts[i+2])
		}
	}
	parent := filepath.Base(filepath.Dir(path))
	return normalizeProductName(parent)
}

func extractCategoryFromZipPath(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i, part := range parts {
		if part == "product" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func isNmapFingerprintFile(data []byte, filename string) bool {
	if strings.HasSuffix(strings.ToLower(filename), ".nfp") {
		return true
	}

	content := string(data)
	return strings.Contains(content, "Nmap service detection probe list") ||
		strings.Contains(content, "# Nmap service detection") ||
		(strings.Contains(content, "Probe ") && strings.Contains(content, "match "))
}

func importNmapFingerprint(data []byte, filename string) (*VulkitImportResult, error) {
	result := &VulkitImportResult{}

	now := time.Now()
	desc := fmt.Sprintf("Nmap 服务检测指纹库 - %s", filename)
	source := 2

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("开始事务失败: %v", tx.Error)
	}

	var skippedEmpty, skippedUnsupported int
	for lineNo, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			skippedEmpty++
			continue
		}
		if !strings.HasPrefix(line, "match ") && !strings.HasPrefix(line, "softmatch ") {
			skippedUnsupported++
			continue
		}

		serviceName := extractNmapServiceName(line)
		productName := extractNmapProductName(line)
		appVersion := extractNmapVersion(line)
		if serviceName == "" {
			serviceName = "unknown"
		}
		appName := productName
		if appName == "" {
			appName = serviceName
		}
		if len([]rune(appName)) > 200 {
			appName = string([]rune(appName)[:200])
		}
		cnName := appName
		if len([]rune(cnName)) > 200 {
			cnName = string([]rune(cnName)[:200])
		}

		appClass := guessAppClassWithLine(serviceName, line)
		fingerType := enums.FingerEnum.GetFingerTypeByClass(appClass)
		level := enums.FingerEnum.GetFingerLevelByClass(appClass)

		err := tx.Table("finger").Create(map[string]interface{}{
			"app_name":    appName,
			"cn_name":     cnName,
			"app_version": appVersion,
			"source":      source,
			"level":       level,
			"finger_type": fingerType,
			"app_class":   appClass,
			"flag":        line,
			"desc":        desc,
			"create_time": now,
			"update_time": now,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("导入nmap指纹失败: 第%d行: %v", lineNo+1, err)
		}
		result.Success++
		result.Total++
	}

	if result.Total == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("未找到可导入的nmap match/softmatch指纹规则，跳过空行/注释%d行，跳过非规则%d行", skippedEmpty, skippedUnsupported)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交nmap指纹失败: %v", err)
	}
	result.Skip = skippedEmpty + skippedUnsupported
	result.Message = fmt.Sprintf("成功导入 %d 条 Nmap 指纹，跳过 %d 行（空行/注释%d，非规则%d）", result.Success, result.Skip, skippedEmpty, skippedUnsupported)
	return result, nil
}

func extractNmapServiceName(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return ""
	}
	return fields[1]
}

func extractNmapProductName(line string) string {
	return extractNmapFieldValue(line, 'p')
}

func extractNmapVersion(line string) string {
	return extractNmapFieldValue(line, 'v')
}

func extractNmapRuleHints(line string) string {
	var parts []string
	for _, field := range []byte{'p', 'i', 'd', 'o', 'h'} {
		if value := extractNmapFieldValue(line, field); value != "" {
			parts = append(parts, value)
		}
	}
	return strings.Join(parts, " ")
}

func extractNmapFieldValue(line string, field byte) string {
	versionInfo := extractNmapVersionInfoSection(line)
	for i := 0; i+2 < len(versionInfo); i++ {
		if versionInfo[i] != field {
			continue
		}
		if i > 0 && !unicode.IsSpace(rune(versionInfo[i-1])) {
			continue
		}
		delimiter := versionInfo[i+1]
		if unicode.IsLetter(rune(delimiter)) || unicode.IsDigit(rune(delimiter)) || unicode.IsSpace(rune(delimiter)) {
			continue
		}

		var builder strings.Builder
		escaped := false
		for j := i + 2; j < len(versionInfo); j++ {
			ch := versionInfo[j]
			if escaped {
				builder.WriteByte(ch)
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == delimiter {
				return strings.TrimSpace(builder.String())
			}
			builder.WriteByte(ch)
		}
	}
	return ""
}

func extractNmapVersionInfoSection(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return line
	}

	serviceName := fields[1]
	servicePos := strings.Index(line, serviceName)
	if servicePos < 0 {
		return line
	}

	i := servicePos + len(serviceName)
	for i < len(line) && unicode.IsSpace(rune(line[i])) {
		i++
	}
	if i+1 >= len(line) || (line[i] != 'm' && line[i] != 'r') {
		return line
	}

	delimiter := line[i+1]
	i += 2
	escaped := false
	for i < len(line) {
		ch := line[i]
		if escaped {
			escaped = false
			i++
			continue
		}
		if ch == '\\' {
			escaped = true
			i++
			continue
		}
		if ch == delimiter {
			i++
			break
		}
		i++
	}

	for i < len(line) && !unicode.IsSpace(rune(line[i])) {
		i++
	}
	for i < len(line) && unicode.IsSpace(rune(line[i])) {
		i++
	}

	if i >= len(line) {
		return ""
	}
	return line[i:]
}

var nmapIotKeywords = []string{
	"iot", "m2m", "modbus", "mqtt", "coap", "bacnet", "tuya",
	"dnp3", "s7comm", "slmp", "iec104", "opcua", "ethernet ip",
	"ethernetip", "profinet", "canopen", "profibus", "profisafe", "melsec",
	"scada", "plc", "hmi", "dcs", "rtu", "softplc",
	"fx3u", "fx5u", "q series", "simatic", "s7 300", "s7 400", "s7 1200", "s7 1500",
	"allen bradley", "rockwell", "beckhoff", "yokogawa", "woodward", "bosch rexroth", "ge fanuc",
}

func isNmapIotFingerprint(parts ...string) bool {
	normalized := normalizeNmapText(strings.Join(parts, " "))
	if normalized == "" {
		return false
	}
	for _, keyword := range nmapIotKeywords {
		if containsNmapKeyword(normalized, keyword) {
			return true
		}
	}
	return false
}

func normalizeNmapText(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))
	lastSpace := true
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			lastSpace = false
			continue
		}
		if !lastSpace {
			builder.WriteByte(' ')
			lastSpace = true
		}
	}
	return strings.TrimSpace(builder.String())
}

func containsNmapKeyword(normalizedText, keyword string) bool {
	normalizedKeyword := normalizeNmapText(keyword)
	if normalizedKeyword == "" {
		return false
	}
	return strings.Contains(" "+normalizedText+" ", " "+normalizedKeyword+" ")
}
