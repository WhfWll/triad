//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	rulesDir := `D:\shellprojects\rules-master`
	outputPath := `docs\sql\import_yara_rules.sql`
	if len(os.Args) > 1 {
		outputPath = os.Args[1]
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建输出文件失败: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	w := func(format string, a ...interface{}) {
		fmt.Fprintf(outFile, format, a...)
	}

	var allInserts []string
	id := 0

	filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".yar") {
			return nil
		}

		rel, _ := filepath.Rel(rulesDir, path)
		category := detectCategory(rel)
		osType := detectOsType(rel)

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)
		rules := extractRules(text)

		for _, rule := range rules {
			id++
			name := rule.name
			if len(name) > 200 {
				name = name[:200]
			}
			desc := rule.description
			if len(desc) > 500 {
				desc = desc[:500]
			}
			desc = escapeSQL(desc)
			name = escapeSQL(name)
			ruleContent := escapeSQL(rule.fullText)

			riskLevel := detectRiskLevel(category, rule.name)

			insert := fmt.Sprintf("INSERT INTO `malware_rule` (`name`, `description`, `risk_level`, `rule_content`, `os_type`, `category`, `status`) VALUES ('%s', '%s', %d, '%s', %d, '%s', 1);",
				name, desc, riskLevel, ruleContent, osType, category)
			allInserts = append(allInserts, insert)

			if id%1000 == 0 {
				fmt.Fprintf(os.Stderr, "已处理 %d 条规则...\n", id)
			}
		}
		return nil
	})

	w("-- ===========================================\n")
	w("-- YARA 病毒库规则导入脚本\n")
	w("-- 共计 %d 条规则\n", id)
	w("-- 生成时间: 2026-05-17\n")
	w("-- ===========================================\n\n")

	for _, s := range allInserts {
		w("%s\n", s)
	}

	fmt.Fprintf(os.Stderr, "完成！共计 %d 条规则已生成 SQL，输出到 %s\n", id, outputPath)
}

type ruleInfo struct {
	name        string
	description string
	fullText    string
}

func extractRules(text string) []ruleInfo {
	text = stripComments(text)

	var rules []ruleInfo
	metaRe := regexp.MustCompile(`(?m)^\s*description\s*=\s*"([^"]*)"`)

	pos := 0
	for pos < len(text) {
		// 跳过空白和无关行
		trimPos := strings.IndexFunc(text[pos:], func(r rune) bool { return r != ' ' && r != '\t' && r != '\n' && r != '\r' })
		if trimPos < 0 {
			break
		}
		pos += trimPos

		remaining := text[pos:]

		// 跳过 include 指令
		if strings.HasPrefix(remaining, "include") {
			nl := strings.Index(remaining, "\n")
			if nl < 0 {
				break
			}
			pos += nl + 1
			continue
		}

		// 检查是否为 rule 定义
		rulePrefix := "rule "
		privatePrefix := "private "
		isPrivate := false
		if strings.HasPrefix(remaining, privatePrefix) {
			isPrivate = true
			remaining = remaining[len(privatePrefix):]
			pos += len(privatePrefix)
		}

		if !strings.HasPrefix(remaining, rulePrefix) {
			// 不是 rule，跳到下一行
			nl := strings.Index(text[pos:], "\n")
			if nl < 0 {
				break
			}
			pos += nl + 1
			continue
		}
		ruleNameStart := pos
		pos += len(rulePrefix) // 跳过 "rule "

		// 提取规则名：直到遇到空格/换行/tab/:/{
		nameStart := pos
		for pos < len(text) {
			ch := text[pos]
			if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == ':' || ch == '{' {
				break
			}
			pos++
		}
		ruleName := text[nameStart:pos]

		// 跳过 tags 部分（如果有 : 分隔）
		for pos < len(text) && (text[pos] == ' ' || text[pos] == '\t') {
			pos++
		}
		if pos < len(text) && text[pos] == ':' {
			pos++ // 跳过 :
			// 读取 tags 直到找到 {
			for pos < len(text) {
				ch := text[pos]
				if ch == '{' {
					break
				}
				pos++
			}
		}

		// 跳过空白到 {
		for pos < len(text) && (text[pos] == ' ' || text[pos] == '\t' || text[pos] == '\n' || text[pos] == '\r') {
			pos++
		}
		if pos >= len(text) || text[pos] != '{' {
			continue
		}

		// 基于大括号深度找到匹配的 }
		braceDepth := 0
		for pos < len(text) {
			if text[pos] == '{' {
				braceDepth++
			} else if text[pos] == '}' {
				braceDepth--
				if braceDepth == 0 {
					pos++ // 包含最后的 }
					break
				}
			}
			pos++
		}

		if braceDepth != 0 {
			// 没有找到匹配的 }，说明文件不完整，跳过
			break
		}

		// 提取完整规则文本
		fullText := text[ruleNameStart:pos]
		if isPrivate {
			fullText = "private " + fullText
		}

		desc := ""
		metaMatch := metaRe.FindStringSubmatch(fullText)
		if len(metaMatch) >= 2 {
			desc = metaMatch[1]
		}

		if len(ruleName) > 200 {
			ruleName = ruleName[:200]
		}

		rules = append(rules, ruleInfo{
			name:        ruleName,
			description: desc,
			fullText:    fullText,
		})
	}

	return rules
}

func stripComments(text string) string {
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

func detectCategory(relPath string) string {
	relPath = filepath.ToSlash(relPath)
	parts := strings.Split(relPath, "/")
	for _, p := range parts {
		switch {
		case p == "malware":
			return "恶意软件"
		case p == "webshells":
			return "Webshell"
		case p == "maldocs":
			return "恶意文档"
		case p == "exploit_kits":
			return "漏洞利用工具包"
		case p == "email":
			return "钓鱼邮件"
		case p == "antidebug_antivm":
			return "反调试/反虚拟机"
		case p == "crypto":
			return "加密检测"
		case p == "cve_rules":
			return "CVE漏洞利用"
		case p == "packers":
			return "加壳检测"
		case p == "capabilities":
			return "能力检测"
		case p == "utils":
			return "检测工具"
		case strings.Contains(p, "mobile") || strings.Contains(p, "Android"):
			return "移动恶意软件"
		case p == "deprecated":
			return "已废弃规则"
		}
	}
	return "通用检测"
}

func detectOsType(relPath string) int {
	relPath = filepath.ToSlash(relPath)
	parts := strings.Split(relPath, "/")
	for _, p := range parts {
		lower := strings.ToLower(p)
		if lower == "android" || lower == "mobile_malware" || strings.Contains(lower, "mobile") {
			return 3
		}
		if strings.Contains(lower, "linux") {
			return 1
		}
	}
	return 2
}

func detectRiskLevel(category, ruleName string) int {
	switch category {
	case "Webshell", "恶意软件", "CVE漏洞利用":
		return 1
	case "漏洞利用工具包", "钓鱼邮件", "恶意文档":
		return 1
	case "反调试/反虚拟机":
		return 2
	case "加壳检测", "加密检测", "能力检测", "检测工具":
		return 3
	case "移动恶意软件":
		return 1
	default:
		upper := strings.ToUpper(ruleName)
		if strings.Contains(upper, "APT") || strings.Contains(upper, "MALW") || strings.Contains(upper, "TROJAN") {
			return 1
		}
		return 2
	}
}

func escapeSQL(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}
