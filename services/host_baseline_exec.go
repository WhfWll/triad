package services

import (
	"strings"
)

// normalizeBaselineCommand 清理 CIS/手册类规则中带有的 shell 提示符前缀，并去掉 sudo（SSH 非交互无法输入密码）。
func normalizeBaselineCommand(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	for strings.HasPrefix(cmd, "$ ") {
		cmd = strings.TrimSpace(cmd[2:])
	}
	if strings.HasPrefix(cmd, "$") {
		cmd = strings.TrimSpace(cmd[1:])
	}
	cmd = stripSudoPrefix(cmd)
	return cmd
}

func stripSudoPrefix(cmd string) string {
	for {
		lower := strings.ToLower(cmd)
		switch {
		case strings.HasPrefix(lower, "sudo -n "):
			cmd = strings.TrimSpace(cmd[8:])
		case strings.HasPrefix(lower, "sudo "):
			cmd = strings.TrimSpace(cmd[5:])
		default:
			return cmd
		}
	}
}

func hasUnresolvedPlaceholder(text string) bool {
	return strings.Contains(strings.ToLower(text), "placeholder_value")
}

// RuleHasUnresolvedPlaceholder 判断规则命令或期望值是否仍含 CIS 占位符（不适配）。
func RuleHasUnresolvedPlaceholder(commands []string, expectedValue string) bool {
	if hasUnresolvedPlaceholder(expectedValue) {
		return true
	}
	for _, cmd := range commands {
		if hasUnresolvedPlaceholder(cmd) {
			return true
		}
	}
	return false
}

// isBaselineExecutionError 判断输出是否表示命令未能正常执行（与合规「不通过」区分）。
func isBaselineExecutionError(output string) bool {
	output = strings.TrimSpace(output)
	if output == "" {
		return false
	}
	lower := strings.ToLower(output)
	indicators := []string{
		"command not found",
		"syntax error",
		"bash:",
		"/bin/sh:",
		"/bin/bash:",
		"zsh:",
		"ksh:",
		"sudo:",
		"a terminal is required",
		"password is required",
		"not allowed to run sudo",
		"permission denied",
		"cannot execute",
		"no such file or directory",
		"is a directory",
		"error:",
		"not a valid identifier",
	}
	for _, s := range indicators {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}
