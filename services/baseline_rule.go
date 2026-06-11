package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"smart/tools/enums"
	"strings"
	"sync"
)

type BaselineRule struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	Risk            int      `json:"risk"`
	OSType          int      `json:"osType"`
	Commands        []string `json:"commands"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
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
	builtinBaselineRules []BaselineRule
)

// mergeRulesWithBuiltin 将数据库规则与内置基础规则合并，避免仅导入 CIS 审计规则时基础项无法执行。
func mergeRulesWithBuiltin(dbRules []BaselineRule) []BaselineRule {
	if len(builtinBaselineRules) == 0 {
		return dbRules
	}
	existing := make(map[string]bool, len(dbRules))
	for _, r := range dbRules {
		existing[ruleCompositeKey(r.Name, r.Category, r.OSType)] = true
	}
	merged := make([]BaselineRule, 0, len(dbRules)+len(builtinBaselineRules))
	merged = append(merged, dbRules...)
	for _, r := range builtinBaselineRules {
		if !existing[ruleCompositeKey(r.Name, r.Category, r.OSType)] {
			merged = append(merged, r)
		}
	}
	return merged
}

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

		for _, rule := range fileRules {
			rule, ok := sanitizeBaselineRule(rule)
			if !ok {
				continue
			}
			e.rules = append(e.rules, rule)
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

	for _, rule := range rules {
		rule, ok := sanitizeBaselineRule(rule)
		if !ok {
			continue
		}
		e.rules = append(e.rules, rule)
		e.rulesMap[rule.OSType] = append(e.rulesMap[rule.OSType], rule)
	}
	return nil
}

// ImportRules 导入规则列表到引擎（运行时热加载），按 name+category+osType 去重
func (e *BaselineEngine) ImportRules(rules []BaselineRule) (int, int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 构建已有规则唯一键集合
	existing := make(map[string]bool, len(e.rules))
	for _, r := range e.rules {
		key := ruleCompositeKey(r.Name, r.Category, r.OSType)
		existing[key] = true
	}

	skipped := 0
	imported := 0
	for _, rule := range rules {
		rule, ok := sanitizeBaselineRule(rule)
		if !ok {
			skipped++
			continue
		}
		if rule.Name == "" || len(rule.Commands) == 0 {
			skipped++
			continue
		}
		key := ruleCompositeKey(rule.Name, rule.Category, rule.OSType)
		if existing[key] {
			skipped++
			continue
		}
		if rule.OSType < 1 || rule.OSType > 4 {
			rule.OSType = 1
		}
		existing[key] = true
		e.rules = append(e.rules, rule)
		e.rulesMap[rule.OSType] = append(e.rulesMap[rule.OSType], rule)
		imported++
	}
	return imported, skipped
}

// ruleCompositeKey 生成规则唯一键：name|category|osType
func ruleCompositeKey(name string, category, osType int) string {
	return fmt.Sprintf("%s|%d|%d", name, category, osType)
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
		if expected == "" {
			return strings.TrimSpace(output) == ""
		}
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
	if len(output) == 0 || len(pattern) == 0 {
		return false
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(output)
}

func init() {
	builtinBaselineRules = []BaselineRule{
		// SSH配置检查
		{ID: 1, Name: "检查SSH协议版本", Description: "SSH应仅使用协议版本2", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^Protocol' /etc/ssh/sshd_config"}, ExpectedValue: "Protocol 2", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 Protocol 2"},
		{ID: 2, Name: "检查SSH root登录限制", Description: "应禁止root直接通过SSH登录", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PermitRootLogin' /etc/ssh/sshd_config"}, ExpectedValue: "PermitRootLogin no", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 PermitRootLogin no"},
		{ID: 3, Name: "检查SSH空闲超时", Description: "SSH连接空闲超时应不超过600秒", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep 'ClientAliveInterval' /etc/ssh/sshd_config"}, ExpectedValue: "ClientAliveInterval", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 ClientAliveInterval 300"},
		{ID: 4, Name: "检查SSH最大认证尝试", Description: "SSH最大认证尝试次数应不超过3次", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^MaxAuthTries' /etc/ssh/sshd_config"}, ExpectedValue: "MaxAuthTries", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 MaxAuthTries <= 3"},
		{ID: 5, Name: "检查SSH登录超时", Description: "登录超时时间应设置合理", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep 'LoginGraceTime' /etc/ssh/sshd_config"}, ExpectedValue: "LoginGraceTime", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 LoginGraceTime 60"},
		{ID: 6, Name: "检查SSH空密码登录", Description: "应禁止空密码登录", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PermitEmptyPasswords' /etc/ssh/sshd_config"}, ExpectedValue: "PermitEmptyPasswords no", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 PermitEmptyPasswords no"},
		{ID: 7, Name: "检查SSH GSSAPI认证", Description: "未使用的GSSAPI认证应禁用", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^GSSAPIAuthentication' /etc/ssh/sshd_config"}, ExpectedValue: "GSSAPIAuthentication no", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 GSSAPIAuthentication no"},
		{ID: 8, Name: "检查SSH X11转发", Description: "未使用的X11转发应禁用", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^X11Forwarding' /etc/ssh/sshd_config"}, ExpectedValue: "X11Forwarding no", MatchType: "contains", FixSuggestion: "在/etc/ssh/sshd_config中设置 X11Forwarding no"},

		// 密码策略检查
		{ID: 9, Name: "检查密码过期时间", Description: "密码最长使用期限应不超过90天", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PASS_MAX_DAYS' /etc/login.defs"}, ExpectedValue: "PASS_MAX_DAYS", MatchType: "contains", FixSuggestion: "在/etc/login.defs中设置 PASS_MAX_DAYS <= 90"},
		{ID: 10, Name: "检查密码最小长度", Description: "密码最小长度应不低于8位", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PASS_MIN_LEN' /etc/login.defs"}, ExpectedValue: "PASS_MIN_LEN", MatchType: "contains", FixSuggestion: "在/etc/login.defs中设置 PASS_MIN_LEN >= 8"},
		{ID: 11, Name: "检查密码最小有效期", Description: "密码修改间隔应至少1天", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PASS_MIN_DAYS' /etc/login.defs"}, ExpectedValue: "PASS_MIN_DAYS", MatchType: "contains", FixSuggestion: "在/etc/login.defs中设置 PASS_MIN_DAYS >= 1"},
		{ID: 12, Name: "检查密码警告天数", Description: "密码过期前应提前警告", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^PASS_WARN_AGE' /etc/login.defs"}, ExpectedValue: "PASS_WARN_AGE", MatchType: "contains", FixSuggestion: "在/etc/login.defs中设置 PASS_WARN_AGE >= 7"},
		{ID: 13, Name: "检查密码加密算法", Description: "密码应使用SHA512或更高级算法加密", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^ENCRYPT_METHOD' /etc/login.defs"}, ExpectedValue: "SHA512", MatchType: "contains", FixSuggestion: "在/etc/login.defs中设置 ENCRYPT_METHOD SHA512"},

		// 用户权限检查
		{ID: 14, Name: "检查空密码账户", Description: "系统中不应存在空密码账户", Category: enums.BaselineCategoryUserPermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"awk -F: '($2 == \"\") {print $1}' /etc/shadow"}, ExpectedValue: "", MatchType: "exact", FixSuggestion: "为空密码账户设置密码: passwd <username>"},
		{ID: 15, Name: "检查UID为0的账户", Description: "不应存在非root的UID为0账户", Category: enums.BaselineCategoryUserPermission, Risk: enums.BaselineRiskCritical, OSType: enums.BaselineOSTypeLinux, Commands: []string{"awk -F: '$3==0 {print $1}' /etc/passwd"}, ExpectedValue: "root", MatchType: "exact", FixSuggestion: "检查并删除非root的UID为0账户"},
		{ID: 16, Name: "检查sudoers配置", Description: "sudoers文件应限制root权限", Category: enums.BaselineCategoryUserPermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep -E '^[^#].*ALL.*ALL' /etc/sudoers 2>/dev/null || grep -E '^[^#].*ALL.*ALL' /etc/sudoers.d/* 2>/dev/null | head -10"}, ExpectedValue: "", MatchType: "exact", FixSuggestion: "审查sudoers配置，限制ALL权限"},
		{ID: 17, Name: "检查wheel组", Description: "应通过wheel组限制su权限", Category: enums.BaselineCategoryUserPermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^wheel:' /etc/group"}, ExpectedValue: "wheel", MatchType: "contains", FixSuggestion: "配置wheel组并设置/etc/pam.d/su"},

		// 防火墙检查
		{ID: 18, Name: "检查防火墙状态", Description: "系统防火墙应处于开启状态", Category: enums.BaselineCategoryFirewall, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"systemctl is-active firewalld 2>/dev/null || systemctl is-active iptables 2>/dev/null || echo inactive"}, ExpectedValue: "active", MatchType: "contains", FixSuggestion: "启动防火墙服务: systemctl enable --now firewalld"},
		{ID: 19, Name: "检查firewalld服务", Description: "firewalld服务应开机自启", Category: enums.BaselineCategoryFirewall, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"systemctl is-enabled firewalld 2>/dev/null || echo disabled"}, ExpectedValue: "enabled", MatchType: "contains", FixSuggestion: "设置firewalld开机自启: systemctl enable firewalld"},

		// 内核安全检查
		{ID: 20, Name: "检查SELinux状态", Description: "SELinux应处于 enforcing 模式", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"getenforce"}, ExpectedValue: "Enforcing", MatchType: "contains", FixSuggestion: "设置 SELINUX=enforcing 在 /etc/selinux/config"},
		{ID: 21, Name: "检查IP转发", Description: "IP转发应在非网关主机上禁用", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"sysctl net.ipv4.ip_forward"}, ExpectedValue: "net.ipv4.ip_forward = 0", MatchType: "contains", FixSuggestion: "设置 net.ipv4.ip_forward=0 在 /etc/sysctl.conf"},
		{ID: 22, Name: "检查ICMP重定向", Description: "应禁用ICMP重定向接受", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"sysctl net.ipv4.conf.all.accept_redirects"}, ExpectedValue: "= 0", MatchType: "contains", FixSuggestion: "设置 net.ipv4.conf.all.accept_redirects=0"},
		{ID: 23, Name: "检查ICMP源路由", Description: "应禁用ICMP源路由", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"sysctl net.ipv4.conf.all.accept_source_route"}, ExpectedValue: "= 0", MatchType: "contains", FixSuggestion: "设置 net.ipv4.conf.all.accept_source_route=0"},
		{ID: 24, Name: "检查TCP SYN Cookie", Description: "应启用SYN Cookie防止SYN攻击", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"sysctl net.ipv4.tcp_syncookies"}, ExpectedValue: "= 1", MatchType: "contains", FixSuggestion: "设置 net.ipv4.tcp_syncookies=1"},
		{ID: 25, Name: "检查IP spoofing保护", Description: "应启用IP欺骗保护", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"sysctl net.ipv4.conf.all.rp_filter"}, ExpectedValue: "= 1", MatchType: "contains", FixSuggestion: "设置 net.ipv4.conf.all.rp_filter=1"},
		{ID: 26, Name: "检查core dump限制", Description: "应限制core dump", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"ulimit -c", "grep 'hard core' /etc/security/limits.conf 2>/dev/null || echo not set"}, ExpectedValue: "0", MatchType: "contains", FixSuggestion: "在/etc/security/limits.conf中设置 * hard core 0"},
		{ID: 27, Name: "检查IPv6转发", Description: "IPv6转发应禁用", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"sysctl net.ipv6.conf.all.forwarding"}, ExpectedValue: "= 0", MatchType: "contains", FixSuggestion: "设置 net.ipv6.conf.all.forwarding=0"},

		// 文件权限检查
		{ID: 28, Name: "检查/etc/shadow权限", Description: "shadow文件权限应设置为600", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"stat -c %a /etc/shadow"}, ExpectedValue: "600", MatchType: "exact", FixSuggestion: "chmod 600 /etc/shadow"},
		{ID: 29, Name: "检查/etc/passwd权限", Description: "passwd文件权限应设置为644", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"stat -c %a /etc/passwd"}, ExpectedValue: "644", MatchType: "exact", FixSuggestion: "chmod 644 /etc/passwd"},
		{ID: 30, Name: "检查/etc/group权限", Description: "group文件权限应设置为644", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"stat -c %a /etc/group"}, ExpectedValue: "644", MatchType: "exact", FixSuggestion: "chmod 644 /etc/group"},
		{ID: 31, Name: "检查/etc/sudoers权限", Description: "sudoers文件权限应设置为440", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"stat -c %a /etc/sudoers"}, ExpectedValue: "440", MatchType: "exact", FixSuggestion: "chmod 440 /etc/sudoers"},
		{ID: 32, Name: "检查umask设置", Description: "默认umask应设置为027或更严格", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep '^umask' /etc/profile 2>/dev/null || echo not set"}, ExpectedValue: "027", MatchType: "contains", FixSuggestion: "在/etc/profile中设置 umask 027"},
		{ID: 33, Name: "检查不必要的SUID文件", Description: "应审计系统中的SUID文件", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"find / -perm -4000 -type f 2>/dev/null | head -50"}, ExpectedValue: "", MatchType: "not_contains", FixSuggestion: "审计SUID文件，移除不必要的SUID权限"},
		{ID: 34, Name: "检查不必要的SGID文件", Description: "应审计系统中的SGID文件", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"find / -perm -2000 -type f 2>/dev/null | head -50"}, ExpectedValue: "", MatchType: "not_contains", FixSuggestion: "审计SGID文件，移除不必要的SGID权限"},
		{ID: 35, Name: "检查系统文件完整性", Description: "关键系统文件不应被修改", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"rpm -Va 2>/dev/null | head -20 || dpkg --verify 2>/dev/null | head -20 || echo 'package manager not supported'"}, ExpectedValue: "", MatchType: "not_contains", FixSuggestion: "重新安装被篡改的软件包"},
		{ID: 36, Name: "检查/tmp目录权限", Description: "/tmp目录权限应设置合理", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"stat -c %a /tmp"}, ExpectedValue: "1777", MatchType: "exact", FixSuggestion: "chmod 1777 /tmp"},
		{ID: 37, Name: "检查/var/log权限", Description: "日志目录权限应限制", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"stat -c %a /var/log"}, ExpectedValue: "755", MatchType: "exact", FixSuggestion: "chmod 755 /var/log"},

		// 审计日志检查
		{ID: 38, Name: "检查审计服务状态", Description: "审计服务auditd应处于运行状态", Category: enums.BaselineCategoryAuditLog, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"systemctl is-active auditd 2>/dev/null || echo inactive"}, ExpectedValue: "active", MatchType: "contains", FixSuggestion: "启动审计服务: systemctl enable --now auditd"},
		{ID: 39, Name: "检查审计服务开机自启", Description: "auditd服务应开机自启", Category: enums.BaselineCategoryAuditLog, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"systemctl is-enabled auditd 2>/dev/null || echo disabled"}, ExpectedValue: "enabled", MatchType: "contains", FixSuggestion: "设置auditd开机自启: systemctl enable auditd"},
		{ID: 40, Name: "检查审计规则", Description: "应配置必要的审计规则", Category: enums.BaselineCategoryAuditLog, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"auditctl -l 2>/dev/null | wc -l"}, ExpectedValue: "0", MatchType: "not_contains", FixSuggestion: "配置审计规则"},

		// 网络服务检查
		{ID: 41, Name: "检查不必要的网络服务", Description: "应禁用不必要的网络服务", Category: enums.BaselineCategoryNetworkService, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"ss -tuln | grep -E ':22|:80|:443' | wc -l"}, ExpectedValue: "", MatchType: "contains", FixSuggestion: "停止并禁用不必要的服务"},
		{ID: 42, Name: "检查TCP Wrappers配置", Description: "/etc/hosts.allow和hosts.deny应配置", Category: enums.BaselineCategoryNetworkService, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"ls -la /etc/hosts.allow /etc/hosts.deny 2>/dev/null | wc -l"}, ExpectedValue: "2", MatchType: "contains", FixSuggestion: "配置hosts.allow和hosts.deny"},

		// 系统更新检查
		{ID: 43, Name: "检查系统更新", Description: "应定期更新系统", Category: enums.BaselineCategorySystemUpdate, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"yum check-update 2>/dev/null | grep -i security | head -5 || apt list --upgradable 2>/dev/null | grep -i security | head -5 || echo 'no security updates'"}, ExpectedValue: "no security updates", MatchType: "contains", FixSuggestion: "执行系统更新: yum update 或 apt update && apt upgrade"},
		{ID: 44, Name: "检查软件源配置", Description: "软件源应配置正确", Category: enums.BaselineCategorySystemUpdate, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"ls /etc/yum.repos.d/*.repo 2>/dev/null || ls /etc/apt/sources.list.d/*.list 2>/dev/null || echo 'no repos found'"}, ExpectedValue: "", MatchType: "not_contains", FixSuggestion: "配置正确的软件源"},

		// 额外安全检查
		{ID: 45, Name: "检查历史命令记录", Description: "应限制命令历史记录长度", Category: enums.BaselineCategoryOther, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep 'HISTSIZE' /etc/profile 2>/dev/null"}, ExpectedValue: "HISTSIZE", MatchType: "contains", FixSuggestion: "在/etc/profile中设置合理的HISTSIZE"},
		{ID: 46, Name: "检查SSH密钥权限", Description: "SSH密钥文件权限应设置为600", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeLinux, Commands: []string{"find /home -name '.ssh' -type d 2>/dev/null | xargs ls -la 2>/dev/null | head -20"}, ExpectedValue: "", MatchType: "not_contains", FixSuggestion: "设置SSH密钥文件权限为600"},
		{ID: 47, Name: "检查crontab权限", Description: "crontab文件权限应限制", Category: enums.BaselineCategoryOther, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"stat -c %a /etc/crontab 2>/dev/null"}, ExpectedValue: "600", MatchType: "exact", FixSuggestion: "chmod 600 /etc/crontab"},
		{ID: 48, Name: "检查PAM配置", Description: "PAM配置应启用密码复杂度检查", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"grep 'pam_cracklib' /etc/pam.d/passwd 2>/dev/null"}, ExpectedValue: "pam_cracklib", MatchType: "contains", FixSuggestion: "在/etc/pam.d/passwd中添加pam_cracklib模块"},
		{ID: 49, Name: "检查NTP同步", Description: "系统时间应同步", Category: enums.BaselineCategoryOther, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeLinux, Commands: []string{"timedatectl status 2>/dev/null | grep 'System clock synchronized'"}, ExpectedValue: "yes", MatchType: "contains", FixSuggestion: "配置NTP时间同步"},
		{ID: 50, Name: "检查邮件服务", Description: "不应运行不必要的邮件服务", Category: enums.BaselineCategoryNetworkService, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeLinux, Commands: []string{"systemctl is-active postfix 2>/dev/null || systemctl is-active sendmail 2>/dev/null || echo inactive"}, ExpectedValue: "inactive", MatchType: "contains", FixSuggestion: "停止并禁用不必要的邮件服务"},

		// Windows 密码策略检查
		{ID: 101, Name: "检查Windows密码策略-最小长度", Description: "密码最小长度应不低于8位", Category: enums.BaselineCategoryPasswordPolicyWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net accounts | findstr /i 'Minimum password length'"}, ExpectedValue: "8", MatchType: "contains", FixSuggestion: "使用gpedit.msc配置密码策略"},
		{ID: 102, Name: "检查Windows密码策略-最长有效期", Description: "密码最长有效期应不超过90天", Category: enums.BaselineCategoryPasswordPolicyWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net accounts | findstr /i 'Maximum password age'"}, ExpectedValue: "90", MatchType: "contains", FixSuggestion: "使用gpedit.msc配置密码策略"},
		{ID: 103, Name: "检查Windows密码策略-最短有效期", Description: "密码最短有效期应至少1天", Category: enums.BaselineCategoryPasswordPolicyWin, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net accounts | findstr /i 'Minimum password age'"}, ExpectedValue: "1", MatchType: "contains", FixSuggestion: "使用gpedit.msc配置密码策略"},
		{ID: 104, Name: "检查Windows密码策略-密码历史", Description: "应保留至少5个历史密码", Category: enums.BaselineCategoryPasswordPolicyWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net accounts | findstr /i 'Password history'"}, ExpectedValue: "5", MatchType: "contains", FixSuggestion: "使用gpedit.msc配置密码策略"},
		{ID: 105, Name: "检查Windows账户锁定策略", Description: "账户锁定阈值应设置合理", Category: enums.BaselineCategoryPasswordPolicyWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net accounts | findstr /i 'Lockout threshold'"}, ExpectedValue: "5", MatchType: "contains", FixSuggestion: "使用gpedit.msc配置账户锁定策略"},

		// Windows 用户权限检查
		{ID: 111, Name: "检查Windows管理员账户", Description: "应禁用默认管理员账户", Category: enums.BaselineCategoryUserPermissionWin, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net user Administrator | findstr /i 'Account active'"}, ExpectedValue: "No", MatchType: "contains", FixSuggestion: "禁用默认管理员账户"},
		{ID: 112, Name: "检查Windows来宾账户", Description: "应禁用来宾账户", Category: enums.BaselineCategoryUserPermissionWin, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net user Guest | findstr /i 'Account active'"}, ExpectedValue: "No", MatchType: "contains", FixSuggestion: "禁用来宾账户"},
		{ID: 113, Name: "检查Windows空密码账户", Description: "不应存在空密码账户", Category: enums.BaselineCategoryUserPermissionWin, Risk: enums.BaselineRiskCritical, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net user | findstr /v '-----' | findstr /v 'User name'"}, ExpectedValue: "", MatchType: "exact", FixSuggestion: "为所有账户设置密码"},

		// Windows 防火墙检查
		{ID: 121, Name: "检查Windows防火墙状态", Description: "Windows防火墙应开启", Category: enums.BaselineCategoryFirewallWin, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeWindows, Commands: []string{"netsh advfirewall show allprofiles state"}, ExpectedValue: "ON", MatchType: "contains", FixSuggestion: "开启Windows防火墙"},
		{ID: 122, Name: "检查Windows防火墙入站规则", Description: "应限制入站规则", Category: enums.BaselineCategoryFirewallWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"netsh advfirewall show inboundrule name=all | findstr /i 'Allow'"}, ExpectedValue: "", MatchType: "exact", FixSuggestion: "审查并限制入站规则"},

		// Windows 服务检查
		{ID: 131, Name: "检查Windows自动更新服务", Description: "自动更新服务应开启", Category: enums.BaselineCategoryServiceWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"sc query wuauserv | findstr /i 'STATE'"}, ExpectedValue: "RUNNING", MatchType: "contains", FixSuggestion: "开启Windows自动更新服务"},
		{ID: 132, Name: "检查Windows远程桌面服务", Description: "远程桌面服务应禁用", Category: enums.BaselineCategoryServiceWin, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeWindows, Commands: []string{"sc query termservice | findstr /i 'STATE'"}, ExpectedValue: "STOPPED", MatchType: "contains", FixSuggestion: "禁用远程桌面服务"},
		{ID: 133, Name: "检查Windows Telnet服务", Description: "Telnet服务应禁用", Category: enums.BaselineCategoryServiceWin, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeWindows, Commands: []string{"sc query tlntsvr | findstr /i 'STATE'"}, ExpectedValue: "STOPPED", MatchType: "contains", FixSuggestion: "禁用Telnet服务"},
		{ID: 134, Name: "检查Windows FTP服务", Description: "FTP服务应禁用", Category: enums.BaselineCategoryServiceWin, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeWindows, Commands: []string{"sc query ftpsvc | findstr /i 'STATE'"}, ExpectedValue: "STOPPED", MatchType: "contains", FixSuggestion: "禁用FTP服务"},

		// Windows 审计策略检查
		{ID: 141, Name: "检查Windows审计策略-账户登录", Description: "应启用账户登录审计", Category: enums.BaselineCategoryAuditPolicyWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"auditpol /get /category:AccountLogon"}, ExpectedValue: "Success and Failure", MatchType: "contains", FixSuggestion: "使用auditpol配置审计策略"},
		{ID: 142, Name: "检查Windows审计策略-登录事件", Description: "应启用登录事件审计", Category: enums.BaselineCategoryAuditPolicyWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"auditpol /get /category:Logon"}, ExpectedValue: "Success and Failure", MatchType: "contains", FixSuggestion: "使用auditpol配置审计策略"},
		{ID: 143, Name: "检查Windows审计策略-账户管理", Description: "应启用账户管理审计", Category: enums.BaselineCategoryAuditPolicyWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"auditpol /get /category:AccountManagement"}, ExpectedValue: "Success and Failure", MatchType: "contains", FixSuggestion: "使用auditpol配置审计策略"},

		// Windows 安全设置检查
		{ID: 151, Name: "检查Windows UAC设置", Description: "UAC应启用", Category: enums.BaselineCategoryOther, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"reg query HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System /v EnableLUA"}, ExpectedValue: "0x1", MatchType: "contains", FixSuggestion: "启用UAC"},
		{ID: 152, Name: "检查Windows远程访问服务", Description: "Remote Registry服务应禁用", Category: enums.BaselineCategoryOther, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeWindows, Commands: []string{"sc query remoteregistry | findstr /i 'STATE'"}, ExpectedValue: "STOPPED", MatchType: "contains", FixSuggestion: "禁用Remote Registry服务"},
		{ID: 153, Name: "检查Windows共享文件夹", Description: "应限制共享文件夹", Category: enums.BaselineCategoryOther, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeWindows, Commands: []string{"net share"}, ExpectedValue: "", MatchType: "exact", FixSuggestion: "删除不必要的共享文件夹"},

		// 国产操作系统（麒麟/UOS）基础安全检查
		{ID: 201, Name: "检查国产操作系统版本", Description: "确认操作系统版本信息", Category: enums.BaselineCategoryOther, Risk: enums.BaselineRiskInfo, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"cat /etc/os-release | grep PRETTY_NAME"}, ExpectedValue: "", MatchType: "contains", FixSuggestion: "确认操作系统版本"},
		{ID: 202, Name: "检查麒麟系统安全服务", Description: "麒麟安全服务应运行", Category: enums.BaselineCategoryServiceWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"systemctl is-active kylin-security 2>/dev/null || echo inactive"}, ExpectedValue: "active", MatchType: "contains", FixSuggestion: "启动麒麟安全服务"},
		{ID: 203, Name: "检查UOS安全中心", Description: "UOS安全中心应运行", Category: enums.BaselineCategoryServiceWin, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"systemctl is-active uos-security-center 2>/dev/null || echo inactive"}, ExpectedValue: "active", MatchType: "contains", FixSuggestion: "启动UOS安全中心"},
		{ID: 204, Name: "检查国产系统SELinux", Description: "SELinux应处于enforcing模式", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"getenforce"}, ExpectedValue: "Enforcing", MatchType: "contains", FixSuggestion: "设置SELINUX=enforcing"},
		{ID: 205, Name: "检查国产系统防火墙", Description: "防火墙应开启", Category: enums.BaselineCategoryFirewall, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"systemctl is-active firewalld 2>/dev/null || systemctl is-active ufw 2>/dev/null || echo inactive"}, ExpectedValue: "active", MatchType: "contains", FixSuggestion: "启动防火墙服务"},
		{ID: 206, Name: "检查国产系统SSH配置", Description: "SSH应仅使用协议版本2", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"grep '^Protocol' /etc/ssh/sshd_config"}, ExpectedValue: "Protocol 2", MatchType: "contains", FixSuggestion: "设置Protocol 2"},
		{ID: 207, Name: "检查国产系统root登录", Description: "应禁止root直接SSH登录", Category: enums.BaselineCategorySSHConfig, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"grep '^PermitRootLogin' /etc/ssh/sshd_config"}, ExpectedValue: "PermitRootLogin no", MatchType: "contains", FixSuggestion: "设置PermitRootLogin no"},
		{ID: 208, Name: "检查国产系统密码策略", Description: "密码最小长度应不低于8位", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"grep '^PASS_MIN_LEN' /etc/login.defs"}, ExpectedValue: "PASS_MIN_LEN", MatchType: "contains", FixSuggestion: "设置PASS_MIN_LEN >= 8"},
		{ID: 209, Name: "检查国产系统密码过期", Description: "密码最长有效期应不超过90天", Category: enums.BaselineCategoryPasswordPolicy, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"grep '^PASS_MAX_DAYS' /etc/login.defs"}, ExpectedValue: "PASS_MAX_DAYS", MatchType: "contains", FixSuggestion: "设置PASS_MAX_DAYS <= 90"},
		{ID: 210, Name: "检查国产系统空密码账户", Description: "不应存在空密码账户", Category: enums.BaselineCategoryUserPermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"awk -F: '($2 == \"\") {print $1}' /etc/shadow"}, ExpectedValue: "", MatchType: "exact", FixSuggestion: "为空密码账户设置密码"},
		{ID: 211, Name: "检查国产系统UID为0账户", Description: "不应存在非root的UID为0账户", Category: enums.BaselineCategoryUserPermission, Risk: enums.BaselineRiskCritical, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"awk -F: '$3==0 {print $1}' /etc/passwd"}, ExpectedValue: "root", MatchType: "exact", FixSuggestion: "检查并删除非root的UID为0账户"},
		{ID: 212, Name: "检查国产系统文件权限", Description: "/etc/shadow权限应设置为600", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"stat -c %a /etc/shadow"}, ExpectedValue: "600", MatchType: "exact", FixSuggestion: "chmod 600 /etc/shadow"},
		{ID: 213, Name: "检查国产系统审计服务", Description: "auditd服务应运行", Category: enums.BaselineCategoryAuditLog, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"systemctl is-active auditd 2>/dev/null || echo inactive"}, ExpectedValue: "active", MatchType: "contains", FixSuggestion: "启动审计服务"},
		{ID: 214, Name: "检查国产系统内核参数", Description: "IP转发应禁用", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"sysctl net.ipv4.ip_forward"}, ExpectedValue: "= 0", MatchType: "contains", FixSuggestion: "设置net.ipv4.ip_forward=0"},
		{ID: 215, Name: "检查国产系统SYN Cookie", Description: "应启用SYN Cookie", Category: enums.BaselineCategoryKernelSecurity, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"sysctl net.ipv4.tcp_syncookies"}, ExpectedValue: "= 1", MatchType: "contains", FixSuggestion: "设置net.ipv4.tcp_syncookies=1"},
		{ID: 216, Name: "检查国产系统安全更新", Description: "应安装安全更新", Category: enums.BaselineCategorySystemUpdate, Risk: enums.BaselineRiskMiddle, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"apt list --upgradable 2>/dev/null | grep -i security | head -5 || yum check-update 2>/dev/null | grep -i security | head -5 || echo 'no security updates'"}, ExpectedValue: "no security updates", MatchType: "contains", FixSuggestion: "执行系统更新"},
		{ID: 217, Name: "检查国产系统sudoers配置", Description: "sudoers文件权限应设置为440", Category: enums.BaselineCategoryFilePermission, Risk: enums.BaselineRiskHigh, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"stat -c %a /etc/sudoers"}, ExpectedValue: "440", MatchType: "exact", FixSuggestion: "chmod 440 /etc/sudoers"},
		{ID: 218, Name: "检查国产系统NTP同步", Description: "系统时间应同步", Category: enums.BaselineCategoryOther, Risk: enums.BaselineRiskLow, OSType: enums.BaselineOSTypeDomestic, Commands: []string{"timedatectl status 2>/dev/null | grep 'System clock synchronized'"}, ExpectedValue: "yes", MatchType: "contains", FixSuggestion: "配置NTP时间同步"},
	}
	if err := GetBaselineEngine().LoadRules(func() []byte {
		data, _ := json.Marshal(builtinBaselineRules)
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
