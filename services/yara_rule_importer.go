package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"smart/models/mysqls"
	"strings"
	"time"
)

type ImportedRule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RiskLevel   int    `json:"riskLevel"`
	RuleContent string `json:"ruleContent"`
	OsType      int    `json:"osType"`
	Category    string `json:"category"`
	Status      int    `json:"status"`
}

type ImportResult struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Errors  []string `json:"errors"`
}

type YaraRuleImporter struct{}

func NewYaraRuleImporter() *YaraRuleImporter {
	return &YaraRuleImporter{}
}

func (imp *YaraRuleImporter) ImportFromYara(ctx context.Context, text string) (*ImportResult, error) {
	rules := extractYaraRules(text)
	if len(rules) == 0 {
		return nil, fmt.Errorf("未从文件中解析到任何规则")
	}

	result := &ImportResult{Total: len(rules)}
	now := time.Now()

	for _, rule := range rules {
		name := rule.Name
		if len(name) > 200 {
			name = name[:200]
		}
		desc := rule.Description
		if len(desc) > 500 {
			desc = desc[:500]
		}

		model := &mysqls.MalwareRule{
			Name:        name,
			Description: desc,
			RiskLevel:   rule.RiskLevel,
			RuleContent: rule.RuleContent,
			OsType:      rule.OsType,
			Category:    rule.Category,
			Status:      1,
			CreateTime:  now,
			UpdateTime:  now,
		}
		if err := model.Add(ctx); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("导入规则 %s 失败: %v", rule.Name, err))
			continue
		}
		result.Success++
	}
	return result, nil
}

func (imp *YaraRuleImporter) ImportFromJSON(ctx context.Context, data []byte) (*ImportResult, error) {
	var rules []ImportedRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}
	if len(rules) == 0 {
		return nil, fmt.Errorf("JSON 中未包含任何规则")
	}

	result := &ImportResult{Total: len(rules)}
	now := time.Now()

	for i, rule := range rules {
		name := rule.Name
		if len(name) > 200 {
			name = name[:200]
		}
		desc := rule.Description
		if len(desc) > 500 {
			desc = desc[:500]
		}
		if rule.Status == 0 {
			rule.Status = 1
		}

		model := &mysqls.MalwareRule{
			Name:        name,
			Description: desc,
			RiskLevel:   rule.RiskLevel,
			RuleContent: rule.RuleContent,
			OsType:      rule.OsType,
			Category:    rule.Category,
			Status:      rule.Status,
			CreateTime:  now,
			UpdateTime:  now,
		}
		if err := model.Add(ctx); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("第 %d 条规则导入失败: %v", i+1, err))
			continue
		}
		result.Success++
	}
	return result, nil
}

type parsedYaraRule struct {
	Name        string
	Description string
	RiskLevel   int
	RuleContent string
	OsType      int
	Category    string
}

func extractYaraRules(text string) []parsedYaraRule {
	text = stripYaraComments(text)

	re := regexp.MustCompile(`(?s)((?:private\s+)?rule\s+(\w+)\s*(?::\s*[\w_]+\s*)?\{.*?\})`)
	matches := re.FindAllStringSubmatch(text, -1)

	metaRe := regexp.MustCompile(`(?m)^\s*description\s*=\s*"([^"]*)"`)

	var rules []parsedYaraRule
	for _, m := range matches {
		fullRule := m[1]
		ruleName := m[2]

		desc := ""
		metaMatch := metaRe.FindStringSubmatch(fullRule)
		if len(metaMatch) >= 2 {
			desc = metaMatch[1]
		}

		cat := detectRuleCategory(ruleName, fullRule)
		risk := detectRuleRisk(cat, ruleName)
		osType := 2

		if strings.Contains(fullRule, "elf") || strings.Contains(fullRule, "ELF") {
			osType = 1
		} else if strings.Contains(ruleName, "Linux") || strings.Contains(ruleName, "Mirai") {
			osType = 1
		}

		rules = append(rules, parsedYaraRule{
			Name:        ruleName,
			Description: desc,
			RiskLevel:   risk,
			RuleContent: fullRule,
			OsType:      osType,
			Category:    cat,
		})
	}
	return rules
}

func stripYaraComments(text string) string {
	var result strings.Builder
	inBlock := false
	i := 0
	for i < len(text) {
		if !inBlock && i+1 < len(text) && text[i] == '/' && text[i+1] == '/' {
			end := strings.Index(text[i:], "\n")
			if end < 0 {
				break
			}
			i += end + 1
			result.WriteByte('\n')
			continue
		}
		if !inBlock && i+1 < len(text) && text[i] == '/' && text[i+1] == '*' {
			inBlock = true
			i += 2
			continue
		}
		if inBlock && i+1 < len(text) && text[i] == '*' && text[i+1] == '/' {
			inBlock = false
			i += 2
			continue
		}
		if !inBlock {
			result.WriteByte(text[i])
		}
		i++
	}
	return result.String()
}

func detectRuleCategory(ruleName, fullRule string) string {
	upper := strings.ToUpper(ruleName)
	full := strings.ToUpper(fullRule)

	if strings.Contains(upper, "WEBSHELL") || strings.Contains(upper, "WSHELL") || strings.Contains(full, ": WEBSHELL") {
		return "Webshell"
	}
	if strings.Contains(upper, "MALD") || strings.Contains(full, "MALDOC") || strings.Contains(full, "CVE-2017-") || strings.Contains(full, "CVE-2018-") {
		return "恶意文档"
	}
	if strings.Contains(upper, "EK_") || strings.Contains(full, "EXPLOIT_KIT") {
		return "漏洞利用工具包"
	}
	if strings.Contains(upper, "EMAIL_") || strings.Contains(full, "PHISHING") || strings.Contains(full, "SCAM") {
		return "钓鱼邮件"
	}
	if strings.Contains(upper, "ANTIDEBUG") || strings.Contains(upper, "ANTIVM") || strings.Contains(full, "ANTI_DEBUG") {
		return "反调试/反虚拟机"
	}
	if strings.Contains(upper, "CRYPTO") || strings.Contains(full, "CRYPTO_SIGNATURES") {
		return "加密检测"
	}
	if strings.Contains(full, "CVE-20") {
		return "CVE漏洞利用"
	}
	if strings.Contains(upper, "PACKER") || strings.Contains(upper, "PEID") || strings.Contains(full, "COMPILER") {
		return "加壳检测"
	}
	if strings.Contains(upper, "ANDROID") || strings.Contains(full, "MOBILE") {
		return "移动恶意软件"
	}
	if strings.Contains(upper, "APT_") || strings.Contains(upper, "APT") {
		return "恶意软件"
	}
	if strings.Contains(upper, "MALW_") || strings.Contains(upper, "TROJAN") || strings.Contains(upper, "BACKDOOR") || strings.Contains(upper, "RAT_") {
		return "恶意软件"
	}
	if strings.Contains(full, "CAPABILIT") {
		return "能力检测"
	}
	return "通用检测"
}

func detectRuleRisk(category, ruleName string) int {
	switch category {
	case "Webshell", "恶意软件", "CVE漏洞利用", "漏洞利用工具包", "钓鱼邮件", "恶意文档", "移动恶意软件":
		return 1
	case "反调试/反虚拟机":
		return 2
	case "加壳检测", "加密检测", "能力检测":
		return 3
	default:
		upper := strings.ToUpper(ruleName)
		if strings.Contains(upper, "APT") || strings.Contains(upper, "MALW") || strings.Contains(upper, "TROJAN") {
			return 1
		}
		return 2
	}
}

type YaraRuleExportItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RiskLevel   int    `json:"riskLevel"`
	RuleContent string `json:"ruleContent"`
	OsType      int    `json:"osType"`
	Category    string `json:"category"`
	Status      int    `json:"status"`
}

func ExportRulesToJSON(rules []mysqls.MalwareRule) ([]byte, error) {
	var items []YaraRuleExportItem
	for _, r := range rules {
		items = append(items, YaraRuleExportItem{
			Name:        r.Name,
			Description: r.Description,
			RiskLevel:   r.RiskLevel,
			RuleContent: r.RuleContent,
			OsType:      r.OsType,
			Category:    r.Category,
			Status:      r.Status,
		})
	}
	return json.MarshalIndent(items, "", "  ")
}
