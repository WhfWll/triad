package network

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var IpSegmentTools ipSegmentTools

type ipSegmentTools struct {
}

// VerifyIp 验证ip格式
func (i ipSegmentTools) VerifyIp(target string) bool {
	rule := `^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$`
	reObj := regexp.MustCompile(rule)
	if reObj != nil && len(reObj.FindAllStringSubmatch(target, -1)) > 0 {
		return true
	}
	return false
}

// VerifyPort 判断端口号是否在0~65535范围内
func (i ipSegmentTools) VerifyPort(port int) bool {
	if port >= 0 && port <= 65535 {
		return true
	}
	return false
}

// VerifyNetmaskIpSegment 验证子网掩码型ip段
func (i ipSegmentTools) VerifyNetmaskIpSegment(target string) bool {
	rule := `^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])/(3[0-2]|[1-2]?[0-9])$`
	reObj := regexp.MustCompile(rule)
	if reObj != nil && len(reObj.FindAllStringSubmatch(target, -1)) > 0 {
		return true
	}
	return false
}

// VerifyCrossbarIpSegment 验证带横杠型ip段
func (i ipSegmentTools) VerifyCrossbarIpSegment(target string) bool {
	rule := `^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])-(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$`
	reObj := regexp.MustCompile(rule)
	if reObj != nil && len(reObj.FindAllString(target, -1)) > 0 {
		return true
	}
	return false
}

// DealCIDR 子网掩码型ip段的逻辑解析
func (i ipSegmentTools) DealCIDR(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []string
	//在循环里创建的所有函数变量共享相同的变量。
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); i.IpTools(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func (i ipSegmentTools) IpTools(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// DealHyphen 带横杠型ip段的逻辑解析
func (i ipSegmentTools) DealHyphen(s string) ([]string, error) {
	tmp := strings.Split(s, ".")
	if len(tmp) == 4 {
		ipRangeTmp := strings.Split(tmp[3], "-")
		var ips []string
		tail, _ := strconv.Atoi(ipRangeTmp[1])
		for head, _ := strconv.Atoi(ipRangeTmp[0]); head <= tail; head++ {
			ips = append(ips, tmp[0]+"."+tmp[1]+"."+tmp[2]+"."+strconv.Itoa(head))
		}
		return ips, nil
	} else {
		return nil, errors.New("格式错误:" + s)
	}
}

// HandleNetmaskIpSegment 处理子网掩码型ip段，获取ip列表
func (i ipSegmentTools) HandleNetmaskIpSegment(target string) ([]string, error) {
	targetList := make([]string, 0)
	ipsTmp, errTmp := i.DealCIDR(target)
	if errTmp != nil {
		return targetList, errTmp
	}
	targetList = append(targetList, ipsTmp...)
	return targetList, nil
}

// HandleCrossbarIpSegment 处理带横杠型ip段，获取ip列表
func (i ipSegmentTools) HandleCrossbarIpSegment(target string) ([]string, error) {
	targetList := make([]string, 0)
	ipsTmp, errTmp := i.DealHyphen(target)
	if errTmp != nil {
		return targetList, errTmp
	}
	targetList = append(targetList, ipsTmp...)
	return targetList, nil
}

// CheckSuccessUrl 验证地址是否合法
func (i ipSegmentTools) CheckSuccessUrl(input string) (bool, error) {
	// 1. 优先验证标准 IPv4 (Strict check)
	if i.VerifyIp(input) {
		return true, nil
	}

	ipv6Re := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
	if ipv6Re.MatchString(input) {
		return true, nil
	}

	httpRe := regexp.MustCompile(`^(http|https)://[^\s]+$`)
	if httpRe.MatchString(input) {
		uu, err := url.Parse(input)
		if err != nil || uu.Scheme == "" || uu.Host == "" {
			err = errors.New("http地址错误")
			return false, err
		} else {
			return true, nil
		}
	} else {
		// 检测是否有http
		if !strings.HasPrefix(input, "http://") {
			// 额外校验：如果包含端口（冒号），且不是 IPv6（前面已通过 regex 检查），则认为是带端口的 IP 或域名，拒绝
			if strings.Contains(input, ":") {
				return false, errors.New("不支持带端口的目标")
			}

			inputWithScheme := "http://" + input
			uu, err := url.Parse(inputWithScheme)
			if err != nil || uu.Scheme == "" || uu.Host == "" {
				err = errors.New("http地址错误")
				return false, err
			} else {
				// 校验 Host 是否合法：必须包含 "." 或者是 "localhost"
				if !strings.Contains(uu.Host, ".") && uu.Host != "localhost" {
					return false, errors.New("非法域名格式")
				}

				// 拦截常见的私有IP误写，如 192.168.3.a
				if strings.HasPrefix(uu.Host, "192.168.") {
					return false, errors.New("非法IP地址: 192.168开头但格式不正确")
				}

				// 校验 Host 字符是否合法（只允许字母、数字、点、短横线）
				hasLetter := false
				for _, r := range uu.Host {
					if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-') {
						return false, errors.New("非法域名格式:包含非法字符")
					}
					if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
						hasLetter = true
					}
				}
				// 如果不是纯数字IP（前面已校验），且不含字母，则既不是域名也不是IP，拦截（例如 10.1.1-1）
				if !hasLetter {
					return false, errors.New("非法域名格式:必须包含字母")
				}

				// 额外校验：防止形如 192.168.0.1211111 这种被误判为域名的非法 IP
				// 如果 Host 全部由数字和点组成，但没通过 VerifyIp（前面已检查），则认为是非法 IP
				isNumeric := true
				for _, r := range uu.Host {
					if r != '.' && (r < '0' || r > '9') {
						isNumeric = false
						break
					}
				}
				if isNumeric {
					return false, errors.New("非法IP地址")
				}

				return true, nil
			}
		}
		err := errors.New("地址非ipv4或ipv6或http 错误")
		return false, err
	}
}

// VerifyMail 验证邮箱地址（支持中文）
func (i ipSegmentTools) VerifyMail(target string) bool {
	rule := "^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	reObj := regexp.MustCompile(rule)
	if reObj != nil && len(reObj.FindAllString(target, -1)) > 0 {
		return true
	}
	return false
}
