package network

import (
	"net"
	"net/url"
	"strings"
)

// 解析url的函数
func ParseUrl(target string) (string, string, string, error) {
	if !strings.HasPrefix(target, "http") {
		target = "http://" + target
	}
	u, err := url.Parse(target)
	if err != nil {
		return "", "", "", err
	}
	splitArray := strings.Split(u.Host, ":")
	host := splitArray[0]
	var port string
	if len(splitArray) == 1 {
		port = "80"
		if u.Scheme == "https" {
			port = "443"
		}
	} else {
		port = splitArray[1]
	}

	return u.Scheme, host, port, nil
}

// 解析域名的函数
func ResolveDomain(name string) (string, error) {
	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

// ParseUrl2 解析url的函数
func ParseUrl2(target string) *url.URL {
	if !strings.HasPrefix(target, "http") && !strings.HasPrefix(target, "data:image") {
		target = "http://" + target
	}

	u, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	return u
}
