package services

import (
	"strings"
	"unicode/utf8"
)

// SanitizeUTF8 去掉非法 UTF-8 字节，避免 MySQL 1366
func SanitizeUTF8(s string) string {
	return strings.ToValidUTF8(s, "")
}

// TruncateRunes 按字符数截断（非按字节）
func TruncateRunes(s string, maxRunes int) string {
	s = SanitizeUTF8(s)
	if maxRunes <= 0 {
		return ""
	}
	rs := []rune(s)
	if len(rs) <= maxRunes {
		return s
	}
	return string(rs[:maxRunes]) + "…"
}

// TruncateUTF8Bytes 按 UTF-8 字节上限截断，不切断多字节字符
func TruncateUTF8Bytes(s string, maxBytes int) string {
	s = SanitizeUTF8(s)
	if maxBytes <= 0 {
		return ""
	}
	if len(s) <= maxBytes {
		return s
	}
	const ellipsis = "…"
	room := maxBytes - len(ellipsis)
	if room < 1 {
		return ellipsis
	}
	for room > 0 && !utf8.RuneStart(s[room]) {
		room--
	}
	return s[:room] + ellipsis
}

// SanitizeDatasecRuleText 导入前裁剪各字段至表列长度
func SanitizeDatasecRuleText(name, desc, fix, riskDesc, expected string) (string, string, string, string, string) {
	return TruncateUTF8Bytes(name, 255),
		TruncateUTF8Bytes(desc, 512),
		SanitizeUTF8(fix),
		TruncateUTF8Bytes(riskDesc, 512),
		TruncateUTF8Bytes(expected, 512)
}
