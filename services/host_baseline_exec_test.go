package services

import "testing"

func TestNormalizeBaselineCommand(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{`$ sudo grep -i "FAIL_DELAY" /etc/login.defs`, `grep -i "FAIL_DELAY" /etc/login.defs`},
		{"grep '^Protocol' /etc/ssh/sshd_config", "grep '^Protocol' /etc/ssh/sshd_config"},
		{"  $ echo ok  ", "echo ok"},
		{"sudo -n cat /etc/shadow", "cat /etc/shadow"},
		{"# grep pam_pwquality /etc/pam.d/system-auth", "grep pam_pwquality /etc/pam.d/system-auth"},
	}
	for _, tc := range cases {
		got := normalizeBaselineCommand(tc.in)
		if got != tc.want {
			t.Fatalf("normalizeBaselineCommand(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIsBaselineExecutionError(t *testing.T) {
	if !isBaselineExecutionError("bash: $: command not found") {
		t.Fatal("expected command not found to be execution error")
	}
	if isBaselineExecutionError("NO_AUDIT_RULES") {
		t.Fatal("NO_AUDIT_RULES should not be execution error")
	}
	if isBaselineExecutionError("") {
		t.Fatal("empty output should not be execution error")
	}
	if isBaselineExecutionError("FAIL_DELAY 5") {
		t.Fatal("normal grep output should not be execution error")
	}
	if !isBaselineExecutionError("sudo: a terminal is required to read the password") {
		t.Fatal("sudo password prompt should be execution error")
	}
	if isBaselineExecutionError("NOT_APPLICABLE") {
		t.Fatal("NOT_APPLICABLE should not be execution error")
	}
}

func TestHasUnresolvedPlaceholder(t *testing.T) {
	if !hasUnresolvedPlaceholder("grep placeholder_value /etc/login.defs") {
		t.Fatal("expected placeholder detected")
	}
	if hasUnresolvedPlaceholder("grep FAIL_DELAY /etc/login.defs") {
		t.Fatal("expected no placeholder")
	}
}

func TestIsBaselineNotApplicableOutput(t *testing.T) {
	if !isBaselineNotApplicableOutput("NOT_APPLICABLE") {
		t.Fatal("expected NOT_APPLICABLE to be detected")
	}
	if isBaselineNotApplicableOutput("not configured") {
		t.Fatal("unexpected not applicable match")
	}
}
