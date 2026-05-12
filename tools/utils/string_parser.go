package utils

import (
	"log"
	"net"
	"net/url"
	"strings"
)

// ExtractHost 处理target_url转成合法的资产信息 资产只保留ip/域名
func ExtractHost(rawURL string) string {
	if ip := net.ParseIP(rawURL); ip != nil {
		return rawURL
	}
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("[ERROR] Failed to parse URL: %s, error: %v", rawURL, err)
		return ""
	}
	if parsedURL.Host == "" {
		return ""
	}
	host, _, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return parsedURL.Host
	}
	return strings.Trim(host, "[]")
}
