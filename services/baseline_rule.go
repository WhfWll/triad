package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"smart/tools/enums"
	"sync"
)

type BaselineRule struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Category       int      `json:"category"`
	Risk           int      `json:"risk"`
	OSType         int      `json:"osType"`
	Commands       []string `json:"commands"`
	ExpectedValue  string   `json:"expectedValue"`
	MatchType      string   `json:"matchType"`
	FixSuggestion  string   `json:"fixSuggestion"`
	RiskDescription string  `json:"riskDescription"`
}

type BaselineEngineConfig struct {
	RulesDir string
}

type BaselineEngine struct {
	mu       sync.RWMutex
	rules    []BaselineRule
	rulesMap map[int][]BaselineRule
	config   *BaselineEngineConfig
}

var (
	globalBaselineEngine *BaselineEngine
	baselineEngineOnce   sync.Once
)

func GetBaselineEngine() *BaselineEngine {
	baselineEngineOnce.Do(func() {
		globalBaselineEngine = &BaselineEngine{
			rulesMap: make(map[int][]BaselineRule),
			config:   &BaselineEngineConfig{RulesDir: "data/baseline"},
		}
	})
	return globalBaselineEngine
}

func (e *BaselineEngine) Init(ctx context.Context, config *BaselineEngineConfig) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if config != nil {
		e.config = config
	}

	rulesDir := e.config.RulesDir
	if rulesDir == "" {
		rulesDir = "data/baseline"
	}

	entries, err := os.ReadDir(rulesDir)
	if err != nil {
		return fmt.Errorf("read rules dir %s failed: %v", rulesDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".json" {
			continue
		}

		filePath := filepath.Join(rulesDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read rule file %s failed: %v", filePath, err)
		}

		var fileRules []BaselineRule
		if err := json.Unmarshal(data, &fileRules); err != nil {
			return fmt.Errorf("parse rule file %s failed: %v", filePath, err)
		}

		e.rules = append(e.rules, fileRules...)
		for _, rule := range fileRules {
			e.rulesMap[rule.OSType] = append(e.rulesMap[rule.OSType], rule)
		}
	}

	return nil
}

func (e *BaselineEngine) LoadRules(data []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var rules []BaselineRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return err
	}

	e.rules = append(e.rules, rules...)
	for _, rule := range rules {
		e.rulesMap[rule.OSType] = append(e.rulesMap[rule.OSType], rule)
	}
	return nil
}

func (e *BaselineEngine) GetRulesByOSType(osType int) []BaselineRule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.rulesMap[osType]
}

func (e *BaselineEngine) GetAllRules() []BaselineRule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]BaselineRule, len(e.rules))
	copy(result, e.rules)
	return result
}

func (e *BaselineEngine) GetRulesByCategory(osType int, category int) []BaselineRule {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var result []BaselineRule
	for _, rule := range e.rulesMap[osType] {
		if rule.Category == category {
			result = append(result, rule)
		}
	}
	return result
}

func (e *BaselineEngine) CheckCommandOutput(output string, expected string, matchType string) bool {
	if matchType == "" || matchType == "contains" {
		return containsOutput(output, expected)
	} else if matchType == "exact" {
		return output == expected
	} else if matchType == "regex" {
		return matchRegex(output, expected)
	} else if matchType == "not_contains" {
		return !containsOutput(output, expected)
	}
	return false
}

func (e *BaselineEngine) GetAllLinuxRules() []BaselineRule {
	return e.GetRulesByOSType(enums.BaselineOSTypeLinux)
}

func (e *BaselineEngine) GetAllWindowsRules() []BaselineRule {
	return e.GetRulesByOSType(enums.BaselineOSTypeWindows)
}

func (e *BaselineEngine) GetAllDomesticRules() []BaselineRule {
	return e.GetRulesByOSType(enums.BaselineOSTypeDomestic)
}

func (e *BaselineEngine) GetAllEmbeddedRules() []BaselineRule {
	return e.GetRulesByOSType(enums.BaselineOSTypeEmbedded)
}

func containsOutput(output, expected string) bool {
	return len(expected) == 0 || containsIgnoreSpace(output, expected)
}

func containsIgnoreSpace(output, expected string) bool {
	oi := 0
	ei := 0
	for oi < len(output) && ei < len(expected) {
		if output[oi] == ' ' || output[oi] == '\t' || output[oi] == '\n' || output[oi] == '\r' {
			oi++
			continue
		}
		if expected[ei] == ' ' || expected[ei] == '\t' || expected[ei] == '\n' || expected[ei] == '\r' {
			ei++
			continue
		}
		if output[oi] != expected[ei] {
			return false
		}
		oi++
		ei++
	}
	return ei == len(expected)
}

func matchRegex(output, pattern string) bool {
	return len(output) > 0 && len(pattern) > 0
}

func init() {
	defaultBaselineRules := []BaselineRule{
		{ID: 1, Name: "检查SSH协议版本", Description: "SSH应仅使用协议版本2", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^Protocol' /etc/ssh/sshd_config"}, ExpectedValue: "Protocol 2", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 Protocol 2"},
		{ID: 2, Name: "检查SSH root登录限制", Description: "应禁止root直接通过SSH登录", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PermitRootLogin' /etc/ssh/sshd_config"}, ExpectedValue: "PermitRootLogin no", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 PermitRootLogin no"},
		{ID: 3, Name: "检查密码过期时间", Description: "密码最长使用期限应不超过90天", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PASS_MAX_DAYS' /etc/login.defs"}, ExpectedValue: "PASS_MAX_DAYS", MatchType: "contains", FixSuggestion: "在/etc/login.defs中设置 PASS_MAX_DAYS <= 90"},
		{ID: 4, Name: "检查密码最小长度", Description: "密码最小长度应不低于8位", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PASS_MIN_LEN' /etc/login.defs"}, ExpectedValue: "PASS_MIN_LEN", MatchType: "contains", FixSuggestion: "在/etc/login.defs中设置 PASS_MIN_LEN >= 8"},
		{ID: 5, Name: "检查SELinux状态", Description: "SELinux应处于 enforcing 模式", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"getenforce"}, ExpectedValue: "Enforcing", MatchType: "contains", FixSuggestion: "设置 SELINUX=enforcing 在 /etc/selinux/config"},
		{ID: 6, Name: "检查防火墙状态", Description: "系统防火墙应处于开启状态", Category: enums.BaselineCategoryFirewall, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"systemctl is-active firewalld 2>/dev/null || systemctl is-active iptables 2>/dev/null || echo inactive"}, ExpectedValue: "active", MatchType: "contains", FixSuggestion: "启动防火墙服务: systemctl enable --now firewalld"},
		{ID: 7, Name: "检查审计服务状态", Description: "审计服务auditd应处于运行状态", Category: enums.BaselineCategoryAuditLog, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"systemctl is-active auditd 2>/dev/null || echo inactive"}, ExpectedValue: "active", MatchType: "contains", FixSuggestion: "启动审计服务: systemctl enable --now auditd"},
		{ID: 8, Name: "检查空密码账户", Description: "系统中不应存在空密码账户", Category: enums.BaselineCategoryUserPermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"awk -F: '($2 == \"\") {print $1}' /etc/shadow"}, ExpectedValue: "", MatchType: "exact", FixSuggestion: "为空密码账户设置密码: passwd <username>"},
		{ID: 9, Name: "检查/etc/shadow权限", Description: "shadow文件权限应设置为600", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"stat -c %a /etc/shadow"}, ExpectedValue: "0", MatchType: "contains", FixSuggestion: "chmod 600 /etc/shadow"},
		{ID: 10, Name: "检查IP转发", Description: "IP转发应在非网关主机上禁用", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"sysctl net.ipv4.ip_forward"}, ExpectedValue: "net.ipv4.ip_forward = 0", MatchType: "contains", FixSuggestion: "设置 net.ipv4.ip_forward=0 在 /etc/sysctl.conf"},
		{ID: 11, Name: "检查ICMP重定向", Description: "应禁用ICMP重定向接受", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"sysctl net.ipv4.conf.all.accept_redirects"}, ExpectedValue: "= 0", MatchType: "contains", FixSuggestion: "设置 net.ipv4.conf.all.accept_redirects=0"},
		{ID: 12, Name: "检查core dump限制", Description: "应限制core dump", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"ulimit -c", "grep 'hard core' /etc/security/limits.conf 2>/dev/null || echo not set"}, ExpectedValue: "0", MatchType: "contains", FixSuggestion: "在/etc/security/limits.conf中设置 * hard core 0"},
		{ID: 13, Name: "检查SSH空闲超时", Description: "SSH连接空闲超时应不超过600秒", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep 'ClientAliveInterval' /etc/ssh/sshd_config"}, ExpectedValue: "ClientAliveInterval", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 ClientAliveInterval 300"},
		{ID: 14, Name: "检查SSH最大认证尝试", Description: "SSH最大认证尝试次数应不超过3次", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^MaxAuthTries' /etc/ssh/sshd_config"}, ExpectedValue: "MaxAuthTries", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 MaxAuthTries <= 3"},
		{ID: 15, Name: "检查umask设置", Description: "默认umask应设置为027或更严格", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^umask' /etc/profile 2>/dev/null || echo not set"}, ExpectedValue: "027", MatchType: "contains", FixSuggestion: "在/etc/profile中设置 umask 027"},
		{ID: 16, Name: "检查不必要的SUID文件", Description: "应审计系统中的SUID文件", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"find / -perm -4000 -type f 2>/dev/null | head -50"}, ExpectedValue: "", MatchType: "not_contains", FixSuggestion: "审计SUID文件，移除不必要的SUID权限"},
		{ID: 17, Name: "检查密码加密算法", Description: "密码应使用SHA512或更高级算法加密", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^ENCRYPT_METHOD' /etc/login.defs"}, ExpectedValue: "SHA512", MatchType: "contains", FixSuggestion: "在/etc/login.defs中设置 ENCRYPT_METHOD SHA512"},
		{ID: 18, Name: "检查系统文件完整性", Description: "关键系统文件不应被修改", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"rpm -Va 2>/dev/null | head -20 || dpkg --verify 2>/dev/null | head -20 || echo 'package manager not supported'"}, ExpectedValue: "", MatchType: "not_contains", FixSuggestion: "重新安装被篡改的软件包"},
	}
	if err := GetBaselineEngine().LoadRules(func() []byte {
		data, _ := json.Marshal(defaultBaselineRules)
		return data
	}()); err != nil {
		fmt.Printf("load default baseline rules warning: %v\n", err)
	}
}

func (e *BaselineEngine) GetRuleByID(id int) (*BaselineRule, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, rule := range e.rules {
		if rule.ID == id {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("rule not found: %d", id)
}
