package network

import (
	"github.com/go-ping/ping"
	"net"
	"strconv"
	"sync"
	"time"
)

// 本地网络
func LocalPing() (network, label, delay string, err error) {
	ipAddr := make([]string, 0)
	ipAddr = append(ipAddr, "114.114.114.114")
	ipAddr = append(ipAddr, "180.76.76.76")
	ipAddr = append(ipAddr, "223.5.5.5")
	return Ping(ipAddr)
}

// 网络情况 ipAddr = [xxx.xxx.com | 192.168.0.0]
func Ping(ipAddr []string) (network, label, delay string, err error) {
	//ipAddr := make([]string, 0)
	//ipAddr = append(ipAddr, "114.114.114.114")
	//ipAddr = append(ipAddr, "180.76.76.76")
	//ipAddr = append(ipAddr, "223.5.5.5")

	wg := sync.WaitGroup{}
	var minRtts time.Duration
	for _, ip := range ipAddr {
		wg.Add(1)
		go func(ip string, err *error) {
			pinger, err1 := ping.NewPinger(ip)
			if err1 != nil {
				*err = err1
			}
			pinger.Timeout = 1 * time.Second
			pinger.Count = 1
			pinger.SetPrivileged(true) //该库尝试通过UDP发送"非特权" ping，所以要加上此步，不然报权限错误
			err2 := pinger.Run()       // Blocks until finished.
			if err2 != nil {
				*err = err2
			}
			stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
			minRtts += stats.MinRtt
			wg.Done()
		}(ip, &err)
	}
	wg.Wait()

	minRtts = minRtts / time.Duration(len(ipAddr))

	if minRtts == 0 {
		label = "断网"
		return "0", label, "-", err
	} else if minRtts < 20*time.Millisecond {
		// 10毫秒内
		network = "5"
		label = "优"
	} else if minRtts < 50*time.Millisecond {
		// 100毫秒内
		network = "4"
		label = "良"
	} else if minRtts < 100*time.Millisecond {
		// 300毫秒内
		network = "3"
		label = "中"
	} else if minRtts < 300*time.Millisecond {
		// 500毫秒内
		network = "2"
		label = "中"
	} else {
		network = "1"
		label = "差"
	}

	return network, label, strconv.FormatFloat(float64(minRtts.Microseconds())/1000, 'f', 3, 64) + "ms", err
}

// 获取本机地址
func GetLocalIp() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ips = append(ips, ipNet.IP.String())
		}
	}
	return ips, nil
}
