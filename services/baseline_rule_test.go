package services

import "testing"

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
