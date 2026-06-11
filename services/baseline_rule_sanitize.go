package services

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	baselineHTMLTagRE            = regexp.MustCompile(`<[^>]+>`)
	baselineBracketPlaceholderRE = regexp.MustCompile(`\[[A-Za-z_][A-Za-z0-9_<>/-]*\]`)
	baselineSSHRuleCmdRE         = regexp.MustCompile(`^sshd -T 2>/dev/null \| grep -E '(\^[^']+)' \|\| grep -E '(\^[^']+)' /etc/ssh/sshd_config 2>/dev/null \|\| echo 'NOT_FOUND'$`)
)

func sanitizeBaselineRule(rule BaselineRule) (BaselineRule, bool) {
	rule.Name = strings.TrimSpace(rule.Name)
	rule.Description = strings.TrimSpace(rule.Description)
	rule.ExpectedValue = strings.TrimSpace(rule.ExpectedValue)
	rule.MatchType = normalizeBaselineMatchType(rule.MatchType, rule.ExpectedValue)

	if RuleHasUnresolvedPlaceholder(rule.Commands, rule.ExpectedValue) {
		return rule, false
	}
	if hasTemplateExpectedValue(rule.ExpectedValue) {
		rewritten, ok := rewriteRecoverableOptionRule(rule)
		if !ok {
			return rule, false
		}
		rule = rewritten
	}

	cmds := make([]string, 0, len(rule.Commands))
	for _, raw := range rule.Commands {
		cmd, ok := sanitizeBaselineCommand(raw)
		if !ok {
			continue
		}
		cmds = append(cmds, cmd)
	}
	if len(cmds) == 0 {
		return rule, false
	}
	if isPlaceholderPseudoCheckRule(cmds, rule.ExpectedValue) {
		return rule, false
	}
	if isBrokenBaselineIfBlock(cmds) {
		return rule, false
	}

	rule.Commands = cmds
	rule = rewriteSpecialBaselineRule(rule)
	rule.MatchType = normalizeBaselineMatchType(rule.MatchType, rule.ExpectedValue)
	return rule, true
}

func SanitizeBaselineRuleForImport(rule BaselineRule) (BaselineRule, bool) {
	return sanitizeBaselineRule(rule)
}

func sanitizeBaselineCommand(raw string) (string, bool) {
	orig := strings.TrimSpace(raw)
	if orig == "" {
		return "", false
	}
	if hasUnsupportedBaselinePlaceholder(orig) {
		return "", false
	}

	cmd := stripLeadingBaselineComment(orig)
	cmd = baselineHTMLTagRE.ReplaceAllString(cmd, "")
	cmd = normalizeBaselineCommand(cmd)
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return "", false
	}
	if isInstructionalBaselineText(cmd) {
		return "", false
	}
	return cmd, true
}

func stripLeadingBaselineComment(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	for strings.HasPrefix(cmd, "#") {
		next := strings.TrimSpace(strings.TrimPrefix(cmd, "#"))
		if next == "" {
			return ""
		}
		if !looksLikeShellCommand(next) {
			return cmd
		}
		cmd = next
	}
	return cmd
}

func looksLikeShellCommand(cmd string) bool {
	lower := strings.ToLower(strings.TrimSpace(cmd))
	if lower == "" {
		return false
	}
	prefixes := []string{
		"grep ", "find ", "ls ", "cat ", "awk ", "sed ", "rpm ", "dpkg ", "ss ",
		"iptables ", "ip6tables ", "nft ", "ufw ", "timedatectl ", "passwd ", "chage ",
		"systemctl ", "test ", "[ ", "echo ", "stat ", "nmcli ", "auditctl ", "getenforce ",
		"sysctl ", "ulimit ", "lsof ", "more ", "chkstat ", "if ", "head ", "tail ",
		"wc ", "readlink ", "realpath ", "file ", "mount ", "journalctl ", "authselect ",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return strings.ContainsAny(cmd, "|/$(){}\\'")
}

func isInstructionalBaselineText(cmd string) bool {
	lower := strings.ToLower(strings.TrimSpace(cmd))
	switch {
	case strings.HasPrefix(lower, "if a listed service is not required"):
		return true
	case strings.HasPrefix(lower, "if any occurrences of"):
		return true
	case strings.HasPrefix(lower, "password substack system-auth"):
		return true
	case strings.HasPrefix(lower, "fw_services_accept_ext="):
		return true
	default:
		return false
	}
}

func hasUnsupportedBaselinePlaceholder(text string) bool {
	if hasUnresolvedPlaceholder(text) {
		return true
	}
	if baselineHTMLTagRE.MatchString(text) {
		return true
	}
	return baselineBracketPlaceholderRE.MatchString(text)
}

func hasTemplateExpectedValue(expected string) bool {
	expected = strings.TrimSpace(strings.ToLower(expected))
	if expected == "" {
		return false
	}
	switch expected {
	case "{option}", "{value}", "{username}", "{user}", "{path}", "{}":
		return true
	default:
		return false
	}
}

func isPlaceholderPseudoCheckRule(cmds []string, expected string) bool {
	if !hasTemplateExpectedValue(expected) {
		return false
	}
	for _, cmd := range cmds {
		lower := strings.ToLower(strings.TrimSpace(cmd))
		if strings.Contains(lower, "grep -e '^.*' /etc/pam.d/.*") &&
			strings.Contains(lower, "grep -e '.*'") {
			return true
		}
	}
	return false
}

func isBrokenBaselineIfBlock(cmds []string) bool {
	hasIf := false
	hasFi := false
	for _, cmd := range cmds {
		lower := strings.ToLower(strings.TrimSpace(cmd))
		if strings.HasPrefix(lower, "if ") && strings.HasSuffix(lower, " then") {
			hasIf = true
		}
		if lower == "fi" || strings.HasSuffix(lower, "; fi") {
			hasFi = true
		}
	}
	return hasIf && !hasFi
}

func normalizeBaselineMatchType(matchType, expected string) string {
	mt := strings.ToLower(strings.TrimSpace(matchType))
	if mt == "" {
		mt = "contains"
	}
	if mt == "contains" && looksLikeRegexPattern(expected) {
		return "regex"
	}
	return mt
}

func looksLikeRegexPattern(expected string) bool {
	expected = strings.TrimSpace(expected)
	if expected == "" {
		return false
	}
	regexMarkers := []string{`\b`, `\h`, `.*`, `.+`, `^`, `$`, `[`, `]`, `(`, `)`, `|`}
	for _, marker := range regexMarkers {
		if strings.Contains(expected, marker) {
			return true
		}
	}
	return false
}

func rewriteSpecialBaselineRule(rule BaselineRule) BaselineRule {
	if strings.Contains(strings.ToLower(rule.Name), "active authselect profile includes pam modules") && len(rule.Commands) > 0 {
		pattern := rule.ExpectedValue
		rule.Commands = []string{fmt.Sprintf(`profile=$(head -1 /etc/authselect/authselect.conf 2>/dev/null || true); if [ -n "$profile" ] && [ -f "/etc/authselect/$profile/system-auth" ] && [ -f "/etc/authselect/$profile/password-auth" ]; then grep -P '%s' "/etc/authselect/$profile"/{system,password}-auth; else echo 'NOT_APPLICABLE'; fi`, pattern)}
		rule.MatchType = "regex"
	}
	if len(rule.Commands) == 1 {
		if rewritten, ok := rewriteGenericSSHConfigRule(rule.Commands[0]); ok {
			rule.Commands = []string{rewritten}
		}
	}
	return rule
}

func rewriteRecoverableOptionRule(rule BaselineRule) (BaselineRule, bool) {
	name := strings.ToLower(strings.TrimSpace(rule.Name))
	switch {
	case strings.Contains(name, "password strength minimum digit characters"):
		return rewritePamPasswordPolicyRule(rule, "dcredit", `(^|\s)dcredit\s*=\s*-\d+($|\s)`), true
	case strings.Contains(name, "password strength minimum different characters"):
		return rewritePamPasswordPolicyRule(rule, "difok", `(^|\s)difok\s*=\s*\d+($|\s)`), true
	case strings.Contains(name, "password strength minimum lowercase characters"):
		return rewritePamPasswordPolicyRule(rule, "lcredit", `(^|\s)lcredit\s*=\s*-\d+($|\s)`), true
	case strings.Contains(name, "password minimum length"):
		return rewritePamPasswordPolicyRule(rule, "minlen", `(^|\s)minlen\s*=\s*\d+($|\s)`), true
	case strings.Contains(name, "password strength minimum special characters"):
		return rewritePamPasswordPolicyRule(rule, "ocredit", `(^|\s)ocredit\s*=\s*-\d+($|\s)`), true
	case strings.Contains(name, "password retry limit"):
		return rewritePamPasswordPolicyRule(rule, "retry", `(^|\s)retry\s*=\s*\d+($|\s)`), true
	case strings.Contains(name, "password strength minimum uppercase characters"):
		return rewritePamPasswordPolicyRule(rule, "ucredit", `(^|\s)ucredit\s*=\s*-\d+($|\s)`), true
	case strings.Contains(name, "limit password reuse") || strings.Contains(name, "限制密码重用"):
		return rewritePamPasswordPolicyRule(rule, "remember", `(^|\s)remember\s*=\s*\d+($|\s)`), true
	case strings.Contains(name, "delay after failed logon attempts") || strings.Contains(name, "登录失败后强制延迟"):
		return rewritePamPasswordPolicyRule(rule, "delay", `(^|\s)delay\s*=\s*\d+($|\s)`), true
	case strings.Contains(name, "pam_wheel for su authentication") || strings.Contains(name, "pam_wheel 进行 su 身份验证"):
		return rewritePamSuWheelRule(rule, false), true
	case strings.Contains(name, "pam_wheel with group parameter") || strings.Contains(name, "带组参数的 pam_wheel"):
		return rewritePamSuWheelRule(rule, true), true
	case strings.Contains(name, "sudo env_reset"):
		return rewriteSudoDefaultsRule(rule, "env_reset", `(^|\s)env_reset(\s|,|$)`), true
	case strings.Contains(name, "sudo ignore_dot"):
		return rewriteSudoDefaultsRule(rule, "ignore_dot", `(^|\s)ignore_dot(\s|,|$)`), true
	case strings.Contains(name, "sudo noexec"):
		return rewriteSudoDefaultsRule(rule, "noexec", `(^|\s)noexec(\s|,|$)`), true
	case strings.Contains(name, "sudo passwd_timeout"):
		return rewriteSudoDefaultsRule(rule, "passwd_timeout", `(^|\s)passwd_timeout\s*=\s*\d+(\s|,|$)`), true
	case strings.Contains(name, "sudo requiretty"):
		return rewriteSudoDefaultsRule(rule, "requiretty", `(^|\s)requiretty(\s|,|$)`), true
	case strings.Contains(name, "sudo umask"):
		return rewriteSudoDefaultsRule(rule, "umask", `(^|\s)umask\s*=\s*0?[0-7]{3,4}(\s|,|$)`), true
	case strings.Contains(name, "sudo use_pty"):
		return rewriteSudoDefaultsRule(rule, "use_pty", `(^|\s)use_pty(\s|,|$)`), true
	case strings.Contains(name, "sudo logfile"):
		return rewriteSudoDefaultsRule(rule, "logfile", `(^|\s)logfile\s*=\s*[^,\s]+(\s|,|$)`), true
	default:
		return rule, false
	}
}

func rewritePamPasswordPolicyRule(rule BaselineRule, key, expected string) BaselineRule {
	rule.ExpectedValue = expected
	rule.MatchType = "regex"
	rule.Commands = []string{fmt.Sprintf(`files=""; for f in /etc/security/pwquality.conf /etc/security/pwquality.conf.d/*.conf /etc/pam.d/system-auth /etc/pam.d/password-auth /etc/pam.d/common-password; do [ -f "$f" ] && files="$files $f"; done; if [ -z "$files" ]; then echo 'NOT_APPLICABLE'; else grep -Eh '\b%s\s*=\s*[^[:space:]]+' $files 2>/dev/null | tail -n 1 || echo 'NOT_CONFIGURED'; fi`, key)}
	return rule
}

func rewritePamSuWheelRule(rule BaselineRule, requireGroup bool) BaselineRule {
	rule.MatchType = "regex"
	if requireGroup {
		rule.ExpectedValue = `pam_wheel\.so.*group=`
	} else {
		rule.ExpectedValue = `pam_wheel\.so`
	}
	rule.Commands = []string{`files=""; for f in /etc/pam.d/su /etc/pam.d/su-l /etc/pam.d/common-auth; do [ -f "$f" ] && files="$files $f"; done; if [ -z "$files" ]; then echo 'NOT_APPLICABLE'; else grep -Eh 'pam_wheel\.so.*' $files 2>/dev/null | tail -n 1 || echo 'NOT_CONFIGURED'; fi`}
	return rule
}

func rewriteSudoDefaultsRule(rule BaselineRule, key, expected string) BaselineRule {
	rule.ExpectedValue = expected
	rule.MatchType = "regex"
	rule.Commands = []string{fmt.Sprintf(`files=""; for f in /etc/sudoers /etc/sudoers.d/*; do [ -f "$f" ] && files="$files $f"; done; if [ -z "$files" ]; then echo 'NOT_APPLICABLE'; else grep -Eh '^[[:space:]]*Defaults\b.*%s.*' $files 2>/dev/null | tail -n 1 || echo 'NOT_CONFIGURED'; fi`, key)}
	return rule
}

func rewriteGenericSSHConfigRule(cmd string) (string, bool) {
	m := baselineSSHRuleCmdRE.FindStringSubmatch(strings.TrimSpace(cmd))
	if len(m) != 3 {
		return "", false
	}
	livePattern := m[1]
	filePattern := m[2]
	return fmt.Sprintf(`sshd_bin=''; if command -v sshd >/dev/null 2>&1; then sshd_bin=$(command -v sshd); elif [ -x /usr/sbin/sshd ]; then sshd_bin=/usr/sbin/sshd; fi; if [ -n "$sshd_bin" ]; then out=$("$sshd_bin" -T 2>/dev/null | grep -E '%s' | head -n 1); if [ -n "$out" ]; then printf '%%s\n' "$out"; exit 0; fi; fi; conf_files=""; if [ -f /etc/ssh/sshd_config ]; then conf_files="/etc/ssh/sshd_config"; fi; if ls /etc/ssh/sshd_config.d/*.conf >/dev/null 2>&1; then conf_files="$conf_files $(ls /etc/ssh/sshd_config.d/*.conf 2>/dev/null)"; fi; if [ -n "$conf_files" ]; then out=$(grep -E '%s' $conf_files 2>/dev/null | tail -n 1); if [ -n "$out" ]; then printf '%%s\n' "$out"; exit 0; fi; fi; echo 'NOT_APPLICABLE'`, livePattern, filePattern), true
}
