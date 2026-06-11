package services

import (
	"strings"
	"testing"
)

func TestCheckCommandOutput(t *testing.T) {
	engine := GetBaselineEngine()

	cases := []struct {
		output   string
		expected string
		match    string
		want     bool
	}{
		{"644", "644", "exact", true},
		{"600", "0", "contains", false},
		{"600", "600", "exact", true},
		{"NO_AUDIT_RULES", "faillock", "contains", false},
		{"", "", "exact", true},
		{"", "", "not_contains", true},
		{"/usr/bin/suid", "", "not_contains", false},
		{"root", "root", "exact", true},
		{"password requisite pam_pwquality.so retry=3", `\bpam_pwquality\.so\b`, "regex", true},
	}

	for _, tc := range cases {
		got := engine.CheckCommandOutput(tc.output, tc.expected, tc.match)
		if got != tc.want {
			t.Fatalf("CheckCommandOutput(%q, %q, %q) = %v, want %v", tc.output, tc.expected, tc.match, got, tc.want)
		}
	}
}

func TestMergeRulesWithBuiltin(t *testing.T) {
	builtinBaselineRules = []BaselineRule{
		{Name: "内置规则", Category: 1, OSType: 1, Commands: []string{"echo ok"}},
	}
	dbRules := []BaselineRule{
		{Name: "CIS审计", Category: 6, OSType: 1, Commands: []string{"auditctl -l"}},
	}
	merged := mergeRulesWithBuiltin(dbRules)
	if len(merged) != 2 {
		t.Fatalf("mergeRulesWithBuiltin len = %d, want 2", len(merged))
	}
}

func TestSanitizeBaselineRuleAuthselect(t *testing.T) {
	rule, ok := sanitizeBaselineRule(BaselineRule{
		Name:          "Ensure Active Authselect Profile Includes PAM Modules",
		Commands:      []string{`# grep -P '\b(pam_pwquality\.so|pam_pwhistory\.so|pam_faillock\.so|pam_unix\.so)\b' /etc/authselect/"$(head -1 /etc/authselect/authselect.conf)"/{system,password}-auth`},
		ExpectedValue: `\b(pam_pwquality\.so|pam_pwhistory\.so|pam_faillock\.so|pam_unix\.so)\b`,
		MatchType:     "contains",
	})
	if !ok {
		t.Fatal("expected authselect rule to be kept after sanitize")
	}
	if rule.MatchType != "regex" {
		t.Fatalf("expected regex matchType, got %s", rule.MatchType)
	}
	if len(rule.Commands) != 1 || !strings.Contains(rule.Commands[0], "NOT_APPLICABLE") {
		t.Fatalf("expected rewritten authselect command, got %#v", rule.Commands)
	}
}

func TestSanitizeBaselineRuleSkipsBrokenInstructionalRule(t *testing.T) {
	_, ok := sanitizeBaselineRule(BaselineRule{
		Name:      "Enable authselect",
		Commands:  []string{`if test "$?" -ne 0; then`, `if rpm --quiet --verify pam; then`, `echo "authselect is not used" >&2`},
		MatchType: "contains",
	})
	if ok {
		t.Fatal("expected broken instructional rule to be skipped")
	}
}

func TestSanitizeBaselineRuleKeepsConcretePwqualityRule(t *testing.T) {
	rule, ok := sanitizeBaselineRule(BaselineRule{
		Name:          "设置密码最小长度",
		Commands:      []string{`grep -E '^minlen' /etc/security/pwquality.conf 2>/dev/null || echo 'NOT_CONFIGURED'`},
		ExpectedValue: `minlen = 14`,
		MatchType:     "contains",
	})
	if !ok {
		t.Fatal("expected concrete pwquality rule to be kept")
	}
	if len(rule.Commands) != 1 || !strings.Contains(rule.Commands[0], "pwquality.conf") {
		t.Fatalf("unexpected command after sanitize: %#v", rule.Commands)
	}
}

func TestSanitizeBaselineRuleRewritesGenericSSHConfigRule(t *testing.T) {
	rule, ok := sanitizeBaselineRule(BaselineRule{
		Name:          "禁用 SSH root Login with a Password (Insecure)",
		Commands:      []string{`sshd -T 2>/dev/null | grep -E '^PermitRootLogin\s' || grep -E '^PermitRootLogin\s+' /etc/ssh/sshd_config 2>/dev/null || echo 'NOT_FOUND'`},
		ExpectedValue: `PermitRootLogin prohibit-password`,
		MatchType:     "contains",
	})
	if !ok {
		t.Fatal("expected SSH config rule to be kept")
	}
	if len(rule.Commands) != 1 {
		t.Fatalf("expected one rewritten SSH command, got %#v", rule.Commands)
	}
	if !strings.Contains(rule.Commands[0], "NOT_APPLICABLE") || !strings.Contains(rule.Commands[0], "sshd_config.d") {
		t.Fatalf("expected rewritten SSH command with fallback logic, got %#v", rule.Commands)
	}
}

func TestSanitizeBaselineRuleRepairsRecoverableOptionRule(t *testing.T) {
	rule, ok := sanitizeBaselineRule(BaselineRule{
		Name:          "设置 Password Strength Minimum Different Characters",
		Commands:      []string{`grep -E '^.*' /etc/pam.d/.* 2>/dev/null | grep -E '.*' || echo 'NOT_CONFIGURED'`},
		ExpectedValue: `{option}`,
		MatchType:     "contains",
	})
	if !ok {
		t.Fatal("expected recoverable option rule to be kept")
	}
	if rule.MatchType != "regex" {
		t.Fatalf("expected regex matchType, got %s", rule.MatchType)
	}
	if rule.ExpectedValue != `(^|\s)difok\s*=\s*\d+($|\s)` {
		t.Fatalf("unexpected expected value: %q", rule.ExpectedValue)
	}
	if len(rule.Commands) != 1 || !strings.Contains(rule.Commands[0], "pwquality.conf") || !strings.Contains(rule.Commands[0], "difok") {
		t.Fatalf("unexpected repaired command: %#v", rule.Commands)
	}
}

func TestSanitizeBaselineRuleSkipsUnrecoverableOptionRule(t *testing.T) {
	_, ok := sanitizeBaselineRule(BaselineRule{
		Name:          "使用 Kerberos 安全挂载远程文件系统",
		Commands:      []string{`mount | grep '.*' 2>/dev/null | grep '.*' || echo 'NOT_FOUND'`},
		ExpectedValue: `{option}`,
		MatchType:     "contains",
	})
	if ok {
		t.Fatal("expected unrecoverable option rule to be skipped")
	}
}
