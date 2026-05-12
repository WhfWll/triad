package data

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestTaskCheckTaskAnalysisTarget_AnalysisTarget(t *testing.T) {
	tests := []struct {
		name            string
		checkTarget     string
		excludeTarget   string
		wantTargetList  []string
		wantErrorList   []string
		shouldHaveError bool
	}{
		{
			name:        "Valid Numeric Domain",
			checkTarget: "123.com",
			wantTargetList: []string{
				"123.com",
			},
			shouldHaveError: false,
		},
		{
			name:        "Mixed Valid and Invalid Format with Dash",
			checkTarget: "192.168.3.1,192.168.3-3",
			wantTargetList: []string{
				"192.168.3.1",
			},
			shouldHaveError: true,
		},
		{
			name:            "Numeric Domain-like with Dash (No Letters)",
			checkTarget:     "10.1.1-1",
			wantTargetList:  []string{},
			shouldHaveError: true,
		},
		{
			name:        "Valid Domain with Dash",
			checkTarget: "my-site.com",
			wantTargetList: []string{
				"my-site.com",
			},
			shouldHaveError: false,
		},
		{
			name:        "Valid IP",
			checkTarget: "192.168.1.1",
			wantTargetList: []string{
				"192.168.1.1",
			},
			shouldHaveError: false,
		},
		{
			name:        "Valid Domain",
			checkTarget: "example.com",
			wantTargetList: []string{
				"example.com",
			},
			shouldHaveError: false,
		},
		{
			name:        "Valid URL with Params",
			checkTarget: "http://example.com/test?id=1&name=test",
			wantTargetList: []string{
				"http://example.com/test?id=1&name=test",
			},
			shouldHaveError: false,
		},
		{
			name:        "Valid Subdomain",
			checkTarget: "sub.example.com",
			wantTargetList: []string{
				"sub.example.com",
			},
			shouldHaveError: false,
		},
		{
			name:        "Valid CIDR",
			checkTarget: "192.168.1.0/30",
			wantTargetList: []string{
				// Implementation filters out network (.0) and broadcast (.3) addresses
				"192.168.1.1", "192.168.1.2",
			},
			shouldHaveError: false,
		},
		{
			name:        "Valid IP Range",
			checkTarget: "192.168.1.1-3",
			wantTargetList: []string{
				"192.168.1.1", "192.168.1.2", "192.168.1.3",
			},
			shouldHaveError: false,
		},
		{
			name:            "Invalid String",
			checkTarget:     "aaaaaaaa",
			wantTargetList:  []string{},
			shouldHaveError: true,
		},
		{
			name:            "Invalid IP Format (Octet > 255)",
			checkTarget:     "192.168.0.1211111",
			wantTargetList:  []string{},
			shouldHaveError: true,
		},
		{
			name:            "Invalid Target with Port",
			checkTarget:     "192.168.0.127:8000",
			wantTargetList:  []string{},
			shouldHaveError: true,
		},
		{
			name:            "Multiple Invalid Targets",
			checkTarget:     "aaaf\nfffffff\ncccca\nssss",
			wantTargetList:  []string{},
			shouldHaveError: true,
		},
		{
			name:          "Mixed Valid and Invalid with Exclusion",
			checkTarget:   "192.168.1.1\naaaaaaaa\n192.168.1.2",
			excludeTarget: "192.168.1.2",
			wantTargetList: []string{
				"192.168.1.1",
			},
			shouldHaveError: true,
		},
		{
			name:            "Invalid Target with Chinese Comma",
			checkTarget:     "192.168.3.3，3",
			wantTargetList:  []string{},
			shouldHaveError: true,
		},
		{
			name:            "Invalid Target like IP but end with char",
			checkTarget:     "192.168.3.a",
			wantTargetList:  []string{},
			shouldHaveError: true,
		},
		{
			name:        "Mixed Valid and Invalid Format",
			checkTarget: "192.168.3.1,192.168.3-3",
			wantTargetList: []string{
				"192.168.3.1",
			},
			shouldHaveError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TaskCheckTaskAnalysisTarget{}
			a.AnalysisTarget(tt.checkTarget, tt.excludeTarget)

			// Check ErrorTargetList
			if tt.shouldHaveError {
				if len(a.ErrorTargetList) == 0 {
					t.Errorf("AnalysisTarget() expected error but got none")
				}
				// If specific error checking is needed, add it here
				// For now, we just verify that an invalid string produces an error
				if strings.Contains(tt.checkTarget, "aaaaaaaa") {
					found := false
					for _, errStr := range a.ErrorTargetList {
						if strings.Contains(errStr, "aaaaaaaa") {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("AnalysisTarget() error list missing expected error for invalid input")
					}
				}

				if tt.name == "Multiple Invalid Targets" {
					expectedMsg := "存在非法测试目标：aaaf、fffffff、cccca、ssss"
					found := false
					for _, errStr := range a.ErrorTargetList {
						if errStr == expectedMsg {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("AnalysisTarget() error list missing expected aggregated error. Got: %v, Want: %v", a.ErrorTargetList, expectedMsg)
					}
				}
			} else {
				if len(a.ErrorTargetList) > 0 {
					t.Errorf("AnalysisTarget() unexpected errors: %v", a.ErrorTargetList)
				}
			}

			// Check TargetList
			// Sort both lists for comparison
			sort.Strings(a.TargetList)
			sort.Strings(tt.wantTargetList)

			if !reflect.DeepEqual(a.TargetList, tt.wantTargetList) {
				// Handle nil vs empty slice difference
				if len(a.TargetList) == 0 && len(tt.wantTargetList) == 0 {
					return
				}
				t.Errorf("AnalysisTarget() TargetList = %v, want %v", a.TargetList, tt.wantTargetList)
			}
		})
	}
}
