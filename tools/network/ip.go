package network

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"regexp"
)

func Ipv4ToInt64(ip string) int64 {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		// 尝试当作域名解析
		addrs, err := net.LookupIP(ip)
		if err != nil || len(addrs) == 0 {
			return 0
		}
		// 只取第一个可用 IPv4
		for _, addr := range addrs {
			if ip4 := addr.To4(); ip4 != nil {
				ret := big.NewInt(0)
				ret.SetBytes(ip4)
				return ret.Int64()
			}
		}
		return 0
	}
	ip4 := parsedIP.To4()
	if ip4 == nil {
		return 0
	}
	ret := big.NewInt(0)
	ret.SetBytes(ip4)
	return ret.Int64()
}

func Ipv4toN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func NtoIpv4(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func IsSingleIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil
}

func IsIPRange(ipStr string) bool {
	match, err := regexp.MatchString(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}-\d{1,3}$`, ipStr)
	if err != nil || !match {
		return false
	}
	return true
}

func IsCIDR(ipStr string) bool {
	_, _, err := net.ParseCIDR(ipStr)
	return err == nil
}

func GetCIDRIpRang(ipStr string, subnetBits int) (string, string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", "", errors.New("无效的IP地址....")
	}
	mask := net.CIDRMask(subnetBits, 32)
	network := ip.Mask(mask)
	firstIP := network
	lastIP := net.IPv4(0, 0, 0, 0).To4()
	for i := 0; i < 4; i++ {
		lastIP[i] = network[i] | ^mask[i]
	}
	return firstIP.String(), lastIP.String(), nil
}
