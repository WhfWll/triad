package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/go-ping/ping"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func IsEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)
	if result {
		return true
	} else {
		return false
	}
}
func PingIpIsLive(ip string) (res bool, err error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return false, err
	}
	pinger.Timeout = 1 * time.Second
	pinger.Count = 1
	pinger.SetPrivileged(true) //该库尝试通过UDP发送"非特权" ping，所以要加上此步，不然报权限错误
	err = pinger.Run()         // Blocks until finished.
	if err != nil {
		return false, err
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	minRtts := stats.MinRtt

	if minRtts == 0 {
		return false, nil
	}
	return true, nil
}

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 计算两个值的百分比 低位数 / 高位数
// hight 高位数
// low 低位数
func MathPercent(high, low int) (string, float64) {
	if low == 0 || high == 0 {
		return "0%", 0
	}
	res := (float64(low) / float64(high)) * 10000
	// 省略小数后面
	resInt := int(res)
	// 计算百分比，留两位
	resFloat := float64(resInt) / 100
	return strconv.FormatFloat(resFloat, 'f', -1, 64) + "%", resFloat
}

// 检查当前系统是否有某个命令
func SystemHasCommand(ctx context.Context, name string) bool {
	cancelCtx, _ := context.WithTimeout(ctx, 2*time.Second)
	cmd := exec.CommandContext(cancelCtx, "which", name)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	if strings.TrimSpace(out.String()) == "" {
		return false
	}
	return true
}

func IsKeyboardSorted(input string) bool {
	input = strings.ToLower(input)
	chars := strings.Split(input, "")
	sort.Strings(chars)
	sortedString := strings.Join(chars, "")
	return input == sortedString
}

// GetHostname 解析一个地址，仅返回主机地址
func GetHostname(host string) string {
	host = strings.TrimSpace(host)
	if strings.HasPrefix(host, "http://") {
		host = strings.TrimLeft(host, "http://")
	}
	if strings.HasPrefix(host, "https://") {
		host = strings.TrimLeft(host, "https://")
	}
	// 去协议与端口
	if indexPortPoint := strings.Index(host, ":"); indexPortPoint >= 0 {
		host = host[:indexPortPoint]
	}
	return host
}

// IsHexString 判断字符串是否为16进制编码
func IsHexString(s string) bool {
	if len(s) == 0 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil && len(s)%2 == 0 && len(s) >= 2
}

func JoinNonEmpty(strs ...string) string {
	var nonEmpty []string
	for _, s := range strs {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, ",")
}
