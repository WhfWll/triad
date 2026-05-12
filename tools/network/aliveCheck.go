package network

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ArpPing 使用ARP请求检测主机是否存活
func ArpPing(targetIP string) (bool, error) {
	// 验证IP地址格式
	if net.ParseIP(targetIP) == nil {
		return false, fmt.Errorf("invalid IP address: %s", targetIP)
	}

	// 根据操作系统选择不同的ARP命令
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Windows使用arp命令
		cmd = exec.Command("arp", "-a", targetIP)
	case "linux", "darwin":
		// Linux和macOS使用arping命令
		cmd = exec.Command("arping", "-c", "1", "-W", "1", targetIP)
	default:
		return false, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果命令执行失败，可能是目标主机不存活或命令不存在
		return false, nil
	}

	// 解析输出结果
	outputStr := string(output)
	if runtime.GOOS == "windows" {
		// Windows下检查是否包含MAC地址
		return strings.Contains(outputStr, targetIP) && !strings.Contains(outputStr, "No ARP Entries Found"), nil
	} else {
		// Linux/macOS下检查arping的返回结果
		return strings.Contains(outputStr, "reply") || strings.Contains(outputStr, "Unicast reply"), nil
	}
}

// ArpPingWithTimeout 带超时的ARP ping检测
func ArpPingWithTimeout(targetIP string, timeout time.Duration) (bool, error) {
	resultChan := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		isAlive, err := ArpPing(targetIP)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- isAlive
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return false, err
	case <-time.After(timeout):
		return false, fmt.Errorf("arp ping timeout for %s", targetIP)
	}
}

// BatchArpPing 批量ARP ping检测
func BatchArpPing(targetIPs []string, timeout time.Duration) map[string]bool {
	results := make(map[string]bool)
	resultChan := make(chan struct {
		ip    string
		alive bool
	}, len(targetIPs))

	// 并发检测
	for _, ip := range targetIPs {
		go func(targetIP string) {
			isAlive, err := ArpPingWithTimeout(targetIP, timeout)
			if err != nil {
				isAlive = false
			}
			resultChan <- struct {
				ip    string
				alive bool
			}{ip: targetIP, alive: isAlive}
		}(ip)
	}

	// 收集结果
	for i := 0; i < len(targetIPs); i++ {
		result := <-resultChan
		results[result.ip] = result.alive
	}

	return results
}

// IcmpPing 使用原生Go实现的ICMP请求检测主机是否存活
func IcmpPing(targetIP string) (bool, error) {
	// 验证IP地址格式
	ip := net.ParseIP(targetIP)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", targetIP)
	}

	// 创建ICMP连接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		// 如果无法创建原始套接字（权限不足），回退到系统ping命令
		return IcmpPingFallback(targetIP)
	}
	defer conn.Close()

	// 创建ICMP Echo Request消息
	message := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Hello, World!"),
		},
	}

	// 序列化ICMP消息
	data, err := message.Marshal(nil)
	if err != nil {
		return false, fmt.Errorf("failed to marshal ICMP message: %v", err)
	}

	// 发送ICMP包
	dst := &net.IPAddr{IP: ip}
	start := time.Now()
	_, err = conn.WriteTo(data, dst)
	if err != nil {
		return false, fmt.Errorf("failed to send ICMP packet: %v", err)
	}

	// 设置读取超时
	err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		return false, fmt.Errorf("failed to set read deadline: %v", err)
	}

	// 读取回复
	reply := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return false, nil // 超时，主机不可达
		}
		return false, fmt.Errorf("failed to read ICMP reply: %v", err)
	}

	// 解析ICMP回复
	rm, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return false, fmt.Errorf("failed to parse ICMP reply: %v", err)
	}

	// 检查是否是Echo Reply
	if rm.Type == ipv4.ICMPTypeEchoReply {
		if echoReply, ok := rm.Body.(*icmp.Echo); ok {
			// 验证ID和序列号
			if echoReply.ID == (os.Getpid()&0xffff) && echoReply.Seq == 1 {
				duration := time.Since(start)
				fmt.Printf("ICMP reply from %v: time=%v\n", peer, duration)
				return true, nil
			}
		}
	}

	return false, nil
}

// IcmpPingFallback 回退到系统ping命令的实现
func IcmpPingFallback(targetIP string) (bool, error) {
	// 根据操作系统选择不同的ping命令
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Windows使用ping命令
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", targetIP)
	case "linux", "darwin":
		// Linux和macOS使用ping命令
		cmd = exec.Command("ping", "-c", "1", "-W", "1", targetIP)
	default:
		return false, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果命令执行失败，可能是目标主机不存活
		return false, nil
	}

	// 解析输出结果
	outputStr := string(output)
	if runtime.GOOS == "windows" {
		// Windows下检查是否收到回复
		return strings.Contains(outputStr, "TTL="), nil
	} else {
		// Linux/macOS下检查ping的返回结果
		return strings.Contains(outputStr, "1 received") || strings.Contains(outputStr, "1 packets received"), nil
	}
}

// IcmpPingAdvanced 高级ICMP ping实现，支持更多选项
func IcmpPingAdvanced(targetIP string, timeout time.Duration, count int) ([]time.Duration, error) {
	// 验证IP地址格式
	ip := net.ParseIP(targetIP)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", targetIP)
	}

	// 创建ICMP连接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to create ICMP connection: %v", err)
	}
	defer conn.Close()

	var results []time.Duration
	pid := os.Getpid() & 0xffff

	for i := 0; i < count; i++ {
		// 创建ICMP Echo Request消息
		message := &icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   pid,
				Seq:  i + 1,
				Data: generateRandomData(32), // 32字节随机数据
			},
		}

		// 序列化ICMP消息
		data, err := message.Marshal(nil)
		if err != nil {
			continue
		}

		// 发送ICMP包
		dst := &net.IPAddr{IP: ip}
		start := time.Now()
		_, err = conn.WriteTo(data, dst)
		if err != nil {
			continue
		}

		// 设置读取超时
		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			continue
		}

		// 读取回复
		reply := make([]byte, 1500)
		n, _, err := conn.ReadFrom(reply)
		if err != nil {
			continue
		}

		// 解析ICMP回复
		rm, err := icmp.ParseMessage(int(ipv4.ICMPTypeEchoReply), reply[:n])
		if err != nil {
			continue
		}

		// 检查是否是正确的Echo Reply
		if rm.Type == ipv4.ICMPTypeEchoReply {
			if echoReply, ok := rm.Body.(*icmp.Echo); ok {
				if echoReply.ID == pid && echoReply.Seq == i+1 {
					duration := time.Since(start)
					results = append(results, duration)
				}
			}
		}

		// 等待一段时间再发送下一个包
		if i < count-1 {
			time.Sleep(1 * time.Second)
		}
	}

	return results, nil
}

// generateRandomData 生成指定长度的随机数据
func generateRandomData(length int) []byte {
	data := make([]byte, length)
	_, err := rand.Read(data)
	if err != nil {
		// 如果随机数生成失败，使用时间戳填充
		timestamp := time.Now().UnixNano()
		binary.LittleEndian.PutUint64(data[:8], uint64(timestamp))
		for i := 8; i < length; i++ {
			data[i] = byte(i)
		}
	}
	return data
}

// TcpPing 使用TCP连接检测主机端口是否存活
func TcpPing(targetIP string, port int) (bool, error) {
	// 验证IP地址格式
	if net.ParseIP(targetIP) == nil {
		return false, fmt.Errorf("invalid IP address: %s", targetIP)
	}

	// 尝试建立TCP连接
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", targetIP, port), 3*time.Second)
	if err != nil {
		return false, nil
	}
	defer conn.Close()
	return true, nil
}

// UdpPing 使用UDP数据包检测主机端口是否存活
// UdpPingAdvanced 更准确的UDP端口检测
func UdpPingAdvanced(targetIP string, port int) (bool, error) {
	// 验证IP地址格式
	if net.ParseIP(targetIP) == nil {
		return false, fmt.Errorf("invalid IP address: %s", targetIP)
	}

	// 根据常见端口发送特定的探测包
	var probeData []byte
	switch port {
	case 53: // DNS
		// 发送DNS查询包
		probeData = []byte{0x12, 0x34, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}
	case 123: // NTP
		// 发送NTP请求包
		probeData = []byte{0x1b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 161: // SNMP
		// 发送SNMP GetRequest
		probeData = []byte{0x30, 0x26, 0x02, 0x01, 0x00, 0x04, 0x06, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0xa0, 0x19, 0x02, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x00, 0x02, 0x01, 0x00, 0x30, 0x0b, 0x30, 0x09, 0x06, 0x05, 0x2b, 0x06, 0x01, 0x02, 0x01, 0x05, 0x00}
	case 631: // IPP/CUPS
		// 使用IPP GetPrinterAttributes请求
		ippRequest := "\x02\x00" + // IPP version 2.0
			"\x00\x0b" + // Get-Printer-Attributes operation
			"\x00\x00\x00\x01" + // request-id
			"\x01" + // begin-attribute-group-tag (operation-attributes-tag)
			"\x47" + // charset attribute
			"\x00\x12attributes-charset" +
			"\x00\x05utf-8" +
			"\x48" + // natural-language attribute
			"\x00\x1battributes-natural-language" +
			"\x00\x05en-us" +
			"\x45" + // uri attribute
			"\x00\x0bprinter-uri" +
			"\x00\x14ipp://localhost/printers/" +
			"\x03" // end-of-attributes-tag
		probeData = []byte(ippRequest)
	case 69: // TFTP
		// TFTP读请求
		probeData = []byte{0x00, 0x01, 0x74, 0x65, 0x73, 0x74, 0x00, 0x6f, 0x63, 0x74, 0x65, 0x74, 0x00}
	default:
		// 通用探测数据
		probeData = []byte("UDP_PROBE")
	}

	// 创建UDP连接
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", targetIP, port), 6*time.Second)
	if err != nil {
		return false, nil
	}
	defer conn.Close()

	// 发送探测数据
	_, err = conn.Write(probeData)
	if err != nil {
		return false, nil
	}

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// 超时：对于UDP来说，这通常意味着端口可能开放但服务不响应
			// 或者端口关闭。我们采用保守策略，认为端口关闭
			return false, nil
		}
		// 其他错误（如连接被拒绝）
		return false, nil
	}

	// 收到回复，端口开放
	if n > 0 {
		return true, nil
	}

	return false, nil
}

// TcpSynPing 使用TCP SYN包检测主机端口是否存活（需要原始套接字权限）
func TcpSynPing(targetIP string, port int) (bool, error) {
	// 验证IP地址格式
	if net.ParseIP(targetIP) == nil {
		return false, fmt.Errorf("invalid IP address: %s", targetIP)
	}

	// 在Windows上，我们使用标准的TCP连接但立即关闭来模拟SYN扫描
	if runtime.GOOS == "windows" {
		return TcpSynPingWindows(targetIP, port)
	}

	// 在Linux/Unix系统上使用原始套接字（需要root权限）
	return TcpSynPingUnix(targetIP, port)
}

// TcpSynPingWindows Windows下的TCP SYN检测实现
func TcpSynPingWindows(targetIP string, port int) (bool, error) {
	// 尝试建立连接但立即关闭
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", targetIP, port), 1*time.Second)
	if err != nil {
		// 检查错误类型，连接被拒绝也表示主机存活
		if strings.Contains(err.Error(), "connection refused") {
			return true, nil
		}
		return false, nil
	}
	defer conn.Close()
	return true, nil
}

// TcpSynPingUnix Unix系统下的TCP SYN检测实现
func TcpSynPingUnix(targetIP string, port int) (bool, error) {
	// 这里提供一个简化的实现，实际的SYN扫描需要原始套接字
	// 由于权限限制，我们使用connect()系统调用的非阻塞模式
	return TcpSynPingWindows(targetIP, port) // 使用相同的逻辑
}

// TcpAckPing 使用TCP ACK包检测主机是否存活（需要原始套接字权限）
func TcpAckPing(targetIP string, port int) (bool, error) {
	// 验证IP地址格式
	if net.ParseIP(targetIP) == nil {
		return false, fmt.Errorf("invalid IP address: %s", targetIP)
	}

	// TCP ACK扫描主要用于绕过防火墙，这里提供一个简化实现
	// 实际的ACK扫描需要构造原始TCP包
	return TcpAckPingSimple(targetIP, port)
}

// TcpAckPingSimple 简化的TCP ACK检测实现
func TcpAckPingSimple(targetIP string, port int) (bool, error) {
	// 尝试连接，如果收到RST响应则表示主机存活
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", targetIP, port), 1*time.Second)
	if err != nil {
		// 连接被重置或拒绝都表示主机存活
		if strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "connection reset") {
			return true, nil
		}
		return false, nil
	}
	defer conn.Close()
	return true, nil
}

// 带超时的各种ping方法
func IcmpPingWithTimeout(targetIP string, timeout time.Duration) (bool, error) {
	resultChan := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		isAlive, err := IcmpPing(targetIP)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- isAlive
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return false, err
	case <-time.After(timeout):
		return false, fmt.Errorf("icmp ping timeout for %s", targetIP)
	}
}

func TcpPingWithTimeout(targetIP string, port int, timeout time.Duration) (bool, error) {
	resultChan := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		isAlive, err := TcpPing(targetIP, port)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- isAlive
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return false, err
	case <-time.After(timeout):
		return false, fmt.Errorf("tcp ping timeout for %s:%d", targetIP, port)
	}
}

func UdpPingWithTimeout(targetIP string, port int, timeout time.Duration) (bool, error) {
	resultChan := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		isAlive, err := UdpPingAdvanced(targetIP, port)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- isAlive
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return false, err
	case <-time.After(timeout):
		return false, fmt.Errorf("udp ping timeout for %s:%d", targetIP, port)
	}
}

func TcpSynPingWithTimeout(targetIP string, port int, timeout time.Duration) (bool, error) {
	resultChan := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		isAlive, err := TcpSynPing(targetIP, port)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- isAlive
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return false, err
	case <-time.After(timeout):
		return false, fmt.Errorf("tcp syn ping timeout for %s:%d", targetIP, port)
	}
}

func TcpAckPingWithTimeout(targetIP string, port int, timeout time.Duration) (bool, error) {
	resultChan := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		isAlive, err := TcpAckPing(targetIP, port)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- isAlive
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return false, err
	case <-time.After(timeout):
		return false, fmt.Errorf("tcp ack ping timeout for %s:%d", targetIP, port)
	}
}

// parsePorts 解析端口字符串，支持单个端口或逗号分隔的多个端口
func parsePorts(portStr string) ([]int, error) {
	var ports []int
	portStrs := strings.Split(portStr, ",")

	for _, ps := range portStrs {
		ps = strings.TrimSpace(ps)
		if ps == "" {
			continue
		}

		port, err := strconv.Atoi(ps)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", ps)
		}

		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("port out of range: %d", port)
		}

		ports = append(ports, port)
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("no valid ports specified")
	}

	return ports, nil
}

// BatchTcpPing 批量TCP端口检测
func BatchTcpPing(targetIP string, ports []int, timeout time.Duration) map[int]bool {
	results := make(map[int]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			isAlive, _ := TcpPingWithTimeout(targetIP, p, timeout)
			mu.Lock()
			results[p] = isAlive
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return results
}

// BatchUdpPing 批量UDP端口检测
func BatchUdpPing(targetIP string, ports []int, timeout time.Duration) map[int]bool {
	results := make(map[int]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			isAlive, _ := UdpPingWithTimeout(targetIP, p, timeout)
			mu.Lock()
			results[p] = isAlive
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return results
}

// BatchTcpSynPing 批量TCP SYN端口检测
func BatchTcpSynPing(targetIP string, ports []int, timeout time.Duration) map[int]bool {
	results := make(map[int]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			isAlive, _ := TcpSynPingWithTimeout(targetIP, p, timeout)
			mu.Lock()
			results[p] = isAlive
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return results
}

// BatchTcpAckPing 批量TCP ACK端口检测
func BatchTcpAckPing(targetIP string, ports []int, timeout time.Duration) map[int]bool {
	results := make(map[int]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			isAlive, _ := TcpAckPingWithTimeout(targetIP, p, timeout)
			mu.Lock()
			results[p] = isAlive
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return results
}

// HandleTcpPortScan 处理TCP端口扫描逻辑
func HandleTcpPortScan(targetIP, portStr string, timeout time.Duration) (bool, error) {
	ports, err := parsePorts(portStr)
	if err != nil {
		fmt.Printf("Error parsing ports: %v\n", err)
		return false, err
	}

	hasAlivePort := false

	if len(ports) == 1 {
		// 单端口检测
		port := ports[0]
		isAlive, err := TcpPingWithTimeout(targetIP, port, timeout)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return false, err
		}
		hasAlivePort = isAlive
		if isAlive {
			fmt.Printf("%s:%d is alive (TCP)\n", targetIP, port)
		} else {
			fmt.Printf("%s:%d is not reachable (TCP)\n", targetIP, port)
		}
	} else {
		// 多端口检测
		fmt.Printf("TCP scanning %s on ports: %s\n", targetIP, portStr)
		results := BatchTcpPing(targetIP, ports, timeout)

		aliveCount := 0
		for _, port := range ports {
			if results[port] {
				fmt.Printf("%s:%d is alive (TCP)\n", targetIP, port)
				aliveCount++
				hasAlivePort = true
			} else {
				fmt.Printf("%s:%d is not reachable (TCP)\n", targetIP, port)
			}
		}
		fmt.Printf("\nSummary: %d/%d ports are alive\n", aliveCount, len(ports))
	}

	return hasAlivePort, nil
}

// HandleUdpPortScan 处理UDP端口扫描逻辑
func HandleUdpPortScan(targetIP, portStr string, timeout time.Duration) (bool, error) {
	ports, err := parsePorts(portStr)
	if err != nil {
		fmt.Printf("Error parsing ports: %v\n", err)
		return false, err
	}

	hasAlivePort := false

	if len(ports) == 1 {
		// 单端口检测
		port := ports[0]
		isAlive, err := UdpPingWithTimeout(targetIP, port, timeout)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return false, err
		}
		hasAlivePort = isAlive
		if isAlive {
			fmt.Printf("%s:%d is alive (UDP)\n", targetIP, port)
		} else {
			fmt.Printf("%s:%d is not reachable (UDP)\n", targetIP, port)
		}
	} else {
		// 多端口检测
		fmt.Printf("UDP scanning %s on ports: %s\n", targetIP, portStr)
		results := BatchUdpPing(targetIP, ports, timeout)

		aliveCount := 0
		for _, port := range ports {
			if results[port] {
				fmt.Printf("%s:%d is alive (UDP)\n", targetIP, port)
				aliveCount++
				hasAlivePort = true
			} else {
				fmt.Printf("%s:%d is not reachable (UDP)\n", targetIP, port)
			}
		}
		fmt.Printf("\nSummary: %d/%d ports are alive\n", aliveCount, len(ports))
	}

	return hasAlivePort, nil
}

// HandleTcpSynPortScan 处理TCP SYN端口扫描逻辑
func HandleTcpSynPortScan(targetIP, portStr string, timeout time.Duration) (bool, error) {
	ports, err := parsePorts(portStr)
	if err != nil {
		fmt.Printf("Error parsing ports: %v\n", err)
		return false, err
	}

	hasAlivePort := false

	if len(ports) == 1 {
		// 单端口检测
		port := ports[0]
		isAlive, err := TcpSynPingWithTimeout(targetIP, port, timeout)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return false, err
		}
		hasAlivePort = isAlive
		if isAlive {
			fmt.Printf("%s:%d is alive (TCP SYN)\n", targetIP, port)
		} else {
			fmt.Printf("%s:%d is not reachable (TCP SYN)\n", targetIP, port)
		}
	} else {
		// 多端口检测
		fmt.Printf("TCP SYN scanning %s on ports: %s\n", targetIP, portStr)
		results := BatchTcpSynPing(targetIP, ports, timeout)

		aliveCount := 0
		for _, port := range ports {
			if results[port] {
				fmt.Printf("%s:%d is alive (TCP SYN)\n", targetIP, port)
				aliveCount++
				hasAlivePort = true
			} else {
				fmt.Printf("%s:%d is not reachable (TCP SYN)\n", targetIP, port)
			}
		}
		fmt.Printf("\nSummary: %d/%d ports are alive\n", aliveCount, len(ports))
	}

	return hasAlivePort, nil
}

// HandleTcpAckPortScan 处理TCP ACK端口扫描
// 返回值：(bool, error) - 只要有一个端口存活就返回true，否则返回false
func HandleTcpAckPortScan(targetIP, portStr string, timeout time.Duration) (bool, error) {
	ports, err := parsePorts(portStr)
	if err != nil {
		return false, fmt.Errorf("error parsing ports: %v", err)
	}

	if len(ports) == 1 {
		// 单端口检测
		port := ports[0]
		isAlive, err := TcpAckPingWithTimeout(targetIP, port, timeout)
		if err != nil {
			return false, fmt.Errorf("error: %v", err)
		}
		if isAlive {
			fmt.Printf("%s:%d is alive (TCP ACK)\n", targetIP, port)
			return true, nil
		} else {
			fmt.Printf("%s:%d is not reachable (TCP ACK)\n", targetIP, port)
			return false, nil
		}
	} else {
		// 多端口检测
		fmt.Printf("TCP ACK scanning %s on ports: %s\n", targetIP, portStr)
		results := BatchTcpAckPing(targetIP, ports, 3*time.Second)

		aliveCount := 0
		hasAlive := false
		for _, port := range ports {
			if results[port] {
				fmt.Printf("%s:%d is alive (TCP ACK)\n", targetIP, port)
				aliveCount++
				hasAlive = true
			} else {
				fmt.Printf("%s:%d is not reachable (TCP ACK)\n", targetIP, port)
			}
		}
		fmt.Printf("\nSummary: %d/%d ports are alive\n", aliveCount, len(ports))
		return hasAlive, nil
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: aliveCheck <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  arp <ip>                    - ARP ping single IP")
		fmt.Println("  icmp <ip>                   - ICMP ping single IP")
		fmt.Println("  tcp <ip> <port[,port...]>   - TCP ping single IP with port(s)")
		fmt.Println("  udp <ip> <port[,port...]>   - UDP ping single IP with port(s)")
		fmt.Println("  syn <ip> <port[,port...]>   - TCP SYN ping single IP with port(s)")
		fmt.Println("  ack <ip> <port[,port...]>   - TCP ACK ping single IP with port(s)")
		fmt.Println("Examples:")
		fmt.Println("  aliveCheck tcp 192.168.1.1 80")
		fmt.Println("  aliveCheck tcp 192.168.1.1 80,443,8080")
		fmt.Println("  aliveCheck udp 192.168.1.1 53,123")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "arp":
		if len(os.Args) < 3 {
			fmt.Println("Usage: aliveCheck arp <ip>")
			os.Exit(1)
		}
		targetIP := os.Args[2]
		isAlive, err := ArpPingWithTimeout(targetIP, 3*time.Second)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if isAlive {
			fmt.Printf("%s is alive (ARP)\n", targetIP)
		} else {
			fmt.Printf("%s is not reachable (ARP)\n", targetIP)
		}

	case "icmp":
		if len(os.Args) < 3 {
			fmt.Println("Usage: aliveCheck icmp <ip>")
			os.Exit(1)
		}
		targetIP := os.Args[2]
		isAlive, err := IcmpPingWithTimeout(targetIP, 3*time.Second)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if isAlive {
			fmt.Printf("%s is alive (ICMP)\n", targetIP)
		} else {
			fmt.Printf("%s is not reachable (ICMP)\n", targetIP)
		}

	case "tcp":
		if len(os.Args) < 4 {
			fmt.Println("Usage: aliveCheck tcp <ip> <port[,port...]>")
			fmt.Println("Examples:")
			fmt.Println("  aliveCheck tcp 192.168.1.1 80")
			fmt.Println("  aliveCheck tcp 192.168.1.1 80,443,8080")
			os.Exit(1)
		}
		targetIP := os.Args[2]
		portStr := os.Args[3]

		// 调用提取的函数
		_, err := HandleTcpPortScan(targetIP, portStr, 3*time.Second)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "udp":
		if len(os.Args) < 4 {
			fmt.Println("Usage: aliveCheck udp <ip> <port[,port...]>")
			fmt.Println("Examples:")
			fmt.Println("  aliveCheck udp 192.168.1.1 53")
			fmt.Println("  aliveCheck udp 192.168.1.1 53,123,161")
			os.Exit(1)
		}
		targetIP := os.Args[2]
		portStr := os.Args[3]

		// 调用提取的函数
		_, err := HandleUdpPortScan(targetIP, portStr, 3*time.Second)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "syn":
		if len(os.Args) < 4 {
			fmt.Println("Usage: aliveCheck syn <ip> <port[,port...]>")
			fmt.Println("Examples:")
			fmt.Println("  aliveCheck syn 192.168.1.1 80")
			fmt.Println("  aliveCheck syn 192.168.1.1 22,80,443")
			os.Exit(1)
		}
		targetIP := os.Args[2]
		portStr := os.Args[3]

		// 调用提取的函数
		_, err := HandleTcpSynPortScan(targetIP, portStr, 3*time.Second)
		if err != nil {
			os.Exit(1)
		}

	case "ack":
		if len(os.Args) < 4 {
			fmt.Println("Usage: aliveCheck ack <ip> <port>")
			os.Exit(1)
		}
		targetIP := os.Args[2]
		portStr := os.Args[3]

		_, err := HandleTcpAckPortScan(targetIP, portStr, 3*time.Second)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

	// 保持原有的ping命令兼容性
	case "ping":
		if len(os.Args) < 3 {
			fmt.Println("Usage: aliveCheck ping <ip>")
			os.Exit(1)
		}
		targetIP := os.Args[2]
		isAlive, err := ArpPingWithTimeout(targetIP, 3*time.Second)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if isAlive {
			fmt.Printf("%s is alive\n", targetIP)
		} else {
			fmt.Printf("%s is not reachable\n", targetIP)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

}
