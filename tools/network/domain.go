package network

import "regexp"

var domainRegex = regexp.MustCompile(`^(?i)[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z]{2,})+$`)

// IsValidDomainFormat 检测域名合法性
func IsValidDomainFormat(domain string) bool {
	return domainRegex.MatchString(domain)
}
