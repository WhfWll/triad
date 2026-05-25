package services

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type YaraRule struct {
	Name       string
	Tags       []string
	Meta       map[string]string
	StringDefs []YaraStringDef
	Condition  YaraCondition
	OSType     int // 0=通用，与 malware_rule.os_type 一致
}

type YaraStringDef struct {
	ID       string
	Value    string
	HexValue []byte
	IsHex    bool
	IsRegex  bool
	Regex    *regexp.Regexp
	Nocase   bool
	Wide     bool
	Fullword bool
}

type YaraCondition struct {
	Op       string
	Children []YaraCondition
	StringID string
	Count    int
	Offset   int
	IsAt     bool
	InStart  int
	InEnd    int
}

type YaraMatch struct {
	RuleName    string
	Description string
	Severity    string
	StringID    string
	Matched     string
	Offset      int
}

type YaraEngine struct {
	mu    sync.RWMutex
	rules []*YaraRule
}

var (
	globalYaraEngine *YaraEngine
	yaraEngineOnce   sync.Once
)

func GetYaraEngine() *YaraEngine {
	yaraEngineOnce.Do(func() {
		globalYaraEngine = &YaraEngine{}
	})
	return globalYaraEngine
}

func (e *YaraEngine) LoadRules(rulesDir string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules = nil
	var allRules []*YaraRule

	err := filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".yar") {
			return nil
		}
		rules, err := parseYaraFile(path)
		if err != nil {
			return fmt.Errorf("parse %s failed: %v", path, err)
		}
		allRules = append(allRules, rules...)
		return nil
	})
	if err != nil {
		return err
	}

	e.rules = allRules
	return nil
}

func (e *YaraEngine) LoadRulesFromFiles(files []string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules = nil
	var allRules []*YaraRule

	for _, f := range files {
		rules, err := parseYaraFile(f)
		if err != nil {
			return fmt.Errorf("parse %s failed: %v", f, err)
		}
		allRules = append(allRules, rules...)
	}

	e.rules = allRules
	return nil
}

// LoadRulesFromContents 从 YARA 规则文本片段加载（如 malware_rule.rule_content）。
func (e *YaraEngine) LoadRulesFromContents(contents []string, osTypes []int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules = nil
	var allRules []*YaraRule

	for i, text := range contents {
		rules := parseYaraContent(text)
		osType := 0
		if i < len(osTypes) {
			osType = osTypes[i]
		}
		for _, rule := range rules {
			rule.OSType = osType
			allRules = append(allRules, rule)
		}
	}

	e.rules = allRules
	return nil
}

func (e *YaraEngine) Match(data []byte) []YaraMatch {
	return e.MatchForOS(data, 0)
}

// MatchForOS 仅匹配通用规则（OSType=0）或与 osType 一致的规则。
func (e *YaraEngine) MatchForOS(data []byte, osType int) []YaraMatch {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var matches []YaraMatch
	for _, rule := range e.rules {
		if osType > 0 && rule.OSType > 0 && rule.OSType != osType {
			continue
		}
		result := evaluateRule(rule, data)
		if result != nil {
			matches = append(matches, result...)
		}
	}
	return matches
}

func (e *YaraEngine) RuleCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.rules)
}

func parseYaraFile(path string) ([]*YaraRule, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseYaraContent(string(content)), nil
}

func parseYaraContent(raw string) []*YaraRule {
	text := stripComments(raw)

	var rules []*YaraRule
	for {
		text = strings.TrimSpace(text)
		if text == "" {
			break
		}

		if strings.HasPrefix(text, "include") {
			idx := strings.Index(text, "\n")
			if idx < 0 {
				break
			}
			text = text[idx+1:]
			continue
		}

		if !strings.HasPrefix(text, "rule ") && !strings.HasPrefix(text, "private rule ") {
			idx := strings.Index(text, "\n")
			if idx < 0 {
				break
			}
			text = text[idx+1:]
			continue
		}

		rule, remaining := parseOneRule(text)
		if rule != nil {
			rules = append(rules, rule)
		}
		text = remaining
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

func parseOneRule(text string) (*YaraRule, string) {
	rule := &YaraRule{
		Meta: make(map[string]string),
	}

	rest := text

	if strings.HasPrefix(rest, "private ") {
		rest = rest[len("private "):]
	}

	rest = rest[len("rule "):]
	rest = strings.TrimLeft(rest, " \t\n")

	nameEnd := strings.IndexAny(rest, " \t\n:{")
	if nameEnd < 0 {
		return nil, text
	}

	ruleName := rest[:nameEnd]
	rule.Name = ruleName
	rest = rest[nameEnd:]

	colonIdx := strings.Index(rest, ":")
	if colonIdx >= 0 {
		tagPart := rest[colonIdx+1:]
		tagEnd := strings.IndexAny(tagPart, " \t\n{")
		if tagEnd > 0 {
			tagStr := strings.TrimSpace(tagPart[:tagEnd])
			if tagStr != "" {
				rule.Tags = strings.Fields(tagStr)
			}
			rest = tagPart[tagEnd:]
		} else {
			rest = tagPart
		}
	}

	braceIdx := strings.Index(rest, "{")
	if braceIdx < 0 {
		return nil, text
	}
	rest = rest[braceIdx+1:]

	depth := 1
	for i := 0; i < len(rest); i++ {
		if rest[i] == '{' {
			depth++
		} else if rest[i] == '}' {
			depth--
			if depth == 0 {
				body := rest[:i]
				remaining := rest[i+1:]
				parseRuleBody(body, rule)
				return rule, remaining
			}
		}
	}

	return nil, text
}

func parseRuleBody(body string, rule *YaraRule) {
	lines := strings.Split(body, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "meta:" {
			currentSection = "meta"
			continue
		} else if line == "strings:" {
			currentSection = "strings"
			continue
		} else if line == "condition:" {
			currentSection = "condition"
			continue
		}

		switch currentSection {
		case "meta":
			if eqIdx := strings.Index(line, "="); eqIdx > 0 {
				key := strings.TrimSpace(line[:eqIdx])
				val := strings.TrimSpace(line[eqIdx+1:])
				val = strings.Trim(val, "\"")
				rule.Meta[key] = val
			}
		case "strings":
			if sd := parseStringDef(line); sd != nil {
				rule.StringDefs = append(rule.StringDefs, *sd)
			}
		case "condition":
			cond := parseCondition(line)
			if cond != nil {
				rule.Condition = *cond
			}
		}
	}
}

func parseStringDef(line string) *YaraStringDef {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "//") {
		return nil
	}

	if !strings.HasPrefix(line, "$") {
		return nil
	}

	sd := &YaraStringDef{}

	spaceIdx := strings.IndexAny(line, " \t=")
	if spaceIdx < 0 {
		return nil
	}

	sd.ID = strings.TrimSpace(line[:spaceIdx])
	rest := strings.TrimSpace(line[spaceIdx+1:])
	if strings.HasPrefix(rest, "=") {
		rest = strings.TrimSpace(rest[1:])
	}

	if strings.HasPrefix(rest, "/") && strings.HasSuffix(rest, "/") {
		sd.IsRegex = true
		sd.Value = rest[1 : len(rest)-1]
		return sd
	}

	if strings.HasPrefix(rest, "{") {
		hexEnd := strings.Index(rest, "}")
		if hexEnd > 0 {
			hexStr := rest[1:hexEnd]
			hexStr = strings.NewReplacer(" ", "", "-", "").Replace(hexStr)
			sd.IsHex = true
			sd.HexValue = decodeHex(hexStr)
		}
		return sd
	}

	if strings.HasPrefix(rest, "\"") {
		endQuote := strings.LastIndex(rest, "\"")
		if endQuote > 0 {
			sd.Value = rest[1:endQuote]
			modifiers := strings.Fields(rest[endQuote+1:])
			for _, m := range modifiers {
				switch strings.ToLower(m) {
				case "nocase":
					sd.Nocase = true
				case "wide":
					sd.Wide = true
				case "fullword":
					sd.Fullword = true
				case "ascii":
				}
			}
		}
	}

	return sd
}

func decodeHex(s string) []byte {
	var result []byte
	for i := 0; i+1 < len(s); i += 2 {
		b := byte(0)
		fmt.Sscanf(s[i:i+2], "%02x", &b)
		result = append(result, b)
	}
	return result
}

func parseCondition(line string) *YaraCondition {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	line = strings.TrimSuffix(line, ";")

	if strings.HasPrefix(line, "all of them") {
		return &YaraCondition{Op: "all"}
	}
	if strings.HasPrefix(line, "any of them") {
		return &YaraCondition{Op: "any"}
	}
	if strings.HasPrefix(line, "all of ") {
		return &YaraCondition{Op: "all"}
	}
	if strings.HasPrefix(line, "any of ") {
		return &YaraCondition{Op: "any"}
	}

	if strings.Contains(line, " and ") {
		parts := strings.Split(line, " and ")
		cond := &YaraCondition{Op: "and"}
		for _, p := range parts {
			if c := parseSimpleCondition(strings.TrimSpace(p)); c != nil {
				cond.Children = append(cond.Children, *c)
			}
		}
		return cond
	}

	if strings.Contains(line, " or ") {
		parts := strings.Split(line, " or ")
		cond := &YaraCondition{Op: "or"}
		for _, p := range parts {
			if c := parseSimpleCondition(strings.TrimSpace(p)); c != nil {
				cond.Children = append(cond.Children, *c)
			}
		}
		return cond
	}

	if c := parseSimpleCondition(line); c != nil {
		return c
	}

	return &YaraCondition{Op: "true"}
}

func parseSimpleCondition(expr string) *YaraCondition {
	expr = strings.TrimSpace(expr)

	if expr == "true" {
		return &YaraCondition{Op: "true"}
	}
	if expr == "false" {
		return &YaraCondition{Op: "false"}
	}

	if strings.HasPrefix(expr, "#") {
		parts := strings.Fields(expr)
		if len(parts) >= 3 && parts[1] == ">=" {
			varID := parts[0]
			count := 0
			fmt.Sscanf(parts[2], "%d", &count)
			return &YaraCondition{Op: "count", StringID: varID, Count: count}
		}
	}

	if strings.HasPrefix(expr, "$") {
		if strings.Contains(expr, " at ") {
			parts := strings.Fields(expr)
			if len(parts) >= 3 {
				offset := 0
				fmt.Sscanf(parts[2], "%d", &offset)
				return &YaraCondition{Op: "at", StringID: parts[0], Offset: offset, IsAt: true}
			}
		}
		if strings.Contains(expr, " in ") {
			parts := strings.Fields(expr)
			if len(parts) >= 4 {
				start, end := 0, 0
				fmt.Sscanf(parts[2], "%d", &start)
				fmt.Sscanf(parts[3], "%d", &end)
				return &YaraCondition{Op: "in", StringID: parts[0], InStart: start, InEnd: end}
			}
		}
		if strings.HasPrefix(expr, "$") {
			return &YaraCondition{Op: "string", StringID: expr}
		}
	}

	if strings.HasPrefix(expr, "not ") {
		inner := strings.TrimSpace(expr[4:])
		if c := parseSimpleCondition(inner); c != nil {
			return &YaraCondition{Op: "not", Children: []YaraCondition{*c}}
		}
	}

	if match := regexp.MustCompile(`^(\d+)\s+of\s+them$`).FindStringSubmatch(expr); len(match) > 0 {
		count := 0
		fmt.Sscanf(match[1], "%d", &count)
		return &YaraCondition{Op: "count_of", Count: count}
	}

	return nil
}

func evaluateRule(rule *YaraRule, data []byte) []YaraMatch {
	stringMatches := make(map[string][]YaraStringMatch)

	for i := range rule.StringDefs {
		matches := matchStringDef(&rule.StringDefs[i], data)
		if len(matches) > 0 {
			stringMatches[rule.StringDefs[i].ID] = matches
		}
	}

	if !evaluateCondition(&rule.Condition, stringMatches, len(rule.StringDefs)) {
		return nil
	}

	var matches []YaraMatch
	severity := rule.Meta["severity"]
	if severity == "" {
		severity = rule.Meta["score"]
	}
	desc := rule.Meta["description"]
	if desc == "" {
		desc = rule.Name
	}

	for _, sd := range rule.StringDefs {
		if ms, ok := stringMatches[sd.ID]; ok {
			for _, m := range ms {
				matches = append(matches, YaraMatch{
					RuleName:    rule.Name,
					Description: desc,
					Severity:    severity,
					StringID:    sd.ID,
					Matched:     string(m.Value),
					Offset:      m.Offset,
				})
			}
		}
	}

	if len(matches) == 0 {
		matches = append(matches, YaraMatch{
			RuleName:    rule.Name,
			Description: desc,
			Severity:    severity,
		})
	}

	return matches
}

type YaraStringMatch struct {
	Value  []byte
	Offset int
}

func matchStringDef(sd *YaraStringDef, data []byte) []YaraStringMatch {
	if sd.IsRegex {
		if sd.Regex == nil {
			re, err := regexp.Compile(sd.Value)
			if err != nil {
				return nil
			}
			sd.Regex = re
		}
		locs := sd.Regex.FindAllIndex(data, -1)
		var matches []YaraStringMatch
		for _, loc := range locs {
			matches = append(matches, YaraStringMatch{
				Value:  data[loc[0]:loc[1]],
				Offset: loc[0],
			})
		}
		return matches
	}

	if sd.IsRegex {
		return nil
	}

	if sd.IsHex && len(sd.HexValue) > 0 {
		return bytesFindAll(data, sd.HexValue)
	}

	if sd.Value != "" {
		searchData := data
		searchStr := sd.Value

		if sd.Nocase {
			searchData = toLowerBytes(data)
			searchStr = strings.ToLower(searchStr)
		}

		var matches []YaraStringMatch
		searchBytes := []byte(searchStr)
		rawBytes := []byte(sd.Value)

		offset := 0
		for {
			idx := indexOf(searchData, searchBytes, offset)
			if idx < 0 {
				break
			}

			if sd.Fullword {
				if !isFullWord(data, idx, len(rawBytes)) {
					offset = idx + 1
					continue
				}
			}

			matches = append(matches, YaraStringMatch{
				Value:  rawBytes,
				Offset: idx,
			})
			offset = idx + 1
		}

		if sd.Wide {
			wideBytes := toWideBytes(searchBytes)
			offset = 0
			for {
				idx := indexOf(searchData, wideBytes, offset)
				if idx < 0 {
					break
				}
				matches = append(matches, YaraStringMatch{
					Value:  rawBytes,
					Offset: idx,
				})
				offset = idx + 1
			}
		}

		return matches
	}

	return nil
}

func evaluateCondition(cond *YaraCondition, matches map[string][]YaraStringMatch, totalStrings int) bool {
	if cond == nil {
		return true
	}

	switch cond.Op {
	case "true":
		return true
	case "false":
		return false
	case "string":
		_, ok := matches[cond.StringID]
		return ok
	case "not":
		if len(cond.Children) > 0 {
			return !evaluateCondition(&cond.Children[0], matches, totalStrings)
		}
		return false
	case "and":
		for _, c := range cond.Children {
			if !evaluateCondition(&c, matches, totalStrings) {
				return false
			}
		}
		return true
	case "or":
		for _, c := range cond.Children {
			if evaluateCondition(&c, matches, totalStrings) {
				return true
			}
		}
		return false
	case "all":
		for _, sd := range cond.Children {
			if !evaluateCondition(&sd, matches, totalStrings) {
				return false
			}
		}
		if len(cond.Children) == 0 {
			for _, ms := range matches {
				if len(ms) == 0 {
					return false
				}
			}
			return true
		}
		return true
	case "any":
		for _, sd := range cond.Children {
			if evaluateCondition(&sd, matches, totalStrings) {
				return true
			}
		}
		if len(cond.Children) == 0 {
			return len(matches) > 0
		}
		return false
	case "count":
		ms, ok := matches[cond.StringID]
		if !ok {
			return false
		}
		return len(ms) >= cond.Count
	case "count_of":
		matchedCount := 0
		for _, ms := range matches {
			if len(ms) > 0 {
				matchedCount++
			}
		}
		return matchedCount >= cond.Count
	case "at":
		ms, ok := matches[cond.StringID]
		if !ok {
			return false
		}
		for _, m := range ms {
			if m.Offset == cond.Offset {
				return true
			}
		}
		return false
	case "in":
		ms, ok := matches[cond.StringID]
		if !ok {
			return false
		}
		for _, m := range ms {
			if m.Offset >= cond.InStart && m.Offset <= cond.InEnd {
				return true
			}
		}
		return false
	}

	return false
}

func indexOf(data, pattern []byte, start int) int {
	if start >= len(data) {
		return -1
	}
	data = data[start:]
	idx := bytesIndexOf(data, pattern)
	if idx < 0 {
		return -1
	}
	return start + idx
}

func bytesIndexOf(data, pattern []byte) int {
	if len(pattern) == 0 {
		return 0
	}
	if len(pattern) > len(data) {
		return -1
	}
	for i := 0; i <= len(data)-len(pattern); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if data[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func bytesFindAll(data, pattern []byte) []YaraStringMatch {
	var matches []YaraStringMatch
	offset := 0
	for {
		idx := bytesIndexOf(data[offset:], pattern)
		if idx < 0 {
			break
		}
		matches = append(matches, YaraStringMatch{
			Value:  pattern,
			Offset: offset + idx,
		})
		offset += idx + 1
	}
	return matches
}

func toLowerBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		if b >= 'A' && b <= 'Z' {
			result[i] = b + 32
		} else {
			result[i] = b
		}
	}
	return result
}

func toWideBytes(data []byte) []byte {
	result := make([]byte, len(data)*2)
	for i, b := range data {
		result[i*2] = b
		result[i*2+1] = 0
	}
	return result
}

func isFullWord(data []byte, start, length int) bool {
	if start > 0 {
		prev := data[start-1]
		if isAlphaNum(prev) {
			return false
		}
	}
	end := start + length
	if end < len(data) {
		next := data[end]
		if isAlphaNum(next) {
			return false
		}
	}
	return true
}

func isAlphaNum(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}
