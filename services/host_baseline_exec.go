package services

import (
	"strings"
)

// normalizeBaselineCommand 清理 CIS/手册类规则中带有的 shell 提示符前缀。
func normalizeBaselineCommand(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	for strings.HasPrefix(cmd, "$ ") {
		cmd = strings.TrimSpace(cmd[2:])
	}
	if strings.HasPrefix(cmd, "$") {
		cmd = strings.TrimSpace(cmd[1:])
	}
	return cmd
}

func hasUnresolvedPlaceholder(text string) bool {
	return strings.Contains(strings.ToLower(text), "placeholder_value")
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
