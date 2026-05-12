package network

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"runtime"
	"strings"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

// ConvertByte2String 处理字符串乱码问题
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

type RouteListRes struct {
	List []RouteListItemRes
}

type RouteListItemRes struct {
	Ip      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

// MacRouteQuery 查询路由表，过滤掉了ipv6
func MacRouteQuery() (RouteListRes, error) {
	name := "netstat"
	arg := "-nr"
	queryCmd := exec.Command(name, arg)
	var routeList RouteListRes
	data, err := queryCmd.Output()
	if err != nil {
		fmt.Printf("exec query command err: %s\n", err)
		return routeList, err
	}
	output := string(data)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Destination") || strings.Contains(line, "::") {
			continue
		}
		line = strings.ReplaceAll(line, "default", "0.0.0.0")
		fields := strings.Fields(line)
		if len(fields) != 5 && len(fields) != 4 {
			continue
		}
		var ip, netmask string
		dstArray := strings.Split(fields[0], "/")
		if len(dstArray) == 2 {
			ip = dstArray[0]
			netmask = dstArray[1]
		} else {
			ip = fields[0]
			netmask = ip
		}
		routeList.List = append(routeList.List, RouteListItemRes{
			Ip:      ip,
			Netmask: netmask,
			Gateway: fields[1],
		})
	}
	return routeList, nil
}

// MacRouteAdd 增加路由
func MacRouteAdd(ip, netmask, gateway string) error {
	name := "route"
	args := []string{"-n", "add", "-net", ip, "-netmask", netmask, gateway}
	addCmd := exec.Command(name, args...)
	data, err := addCmd.Output()
	if err != nil {
		fmt.Printf("mac exec add command err: %s\n", err)
		return err
	}
	output := string(data)
	if strings.Contains(output, "File exists") {
		return errors.New("路由已存在")
	}
	return nil
}

// MacRouteDelete 删除路由
func MacRouteDelete(ip, netmask, gateway string) error {
	name := "route"
	args := []string{"-n", "delete", "-net", ip, "-netmask", netmask, gateway}
	addCmd := exec.Command(name, args...)
	_, err := addCmd.Output()
	if err != nil {
		fmt.Printf("exec delete command err: %s\n", err)
		return err
	}
	return nil
}

// WindowsRouteQuery 查询路由表，过滤掉了ipv6
func WindowsRouteQuery() (RouteListRes, error) {
	name := "route"
	arg := "print"
	queryCmd := exec.Command(name, arg)
	var routeList RouteListRes
	data, err := queryCmd.Output()
	if err != nil {
		fmt.Printf("exec query command err: %s\n", err)
		return routeList, err
	}
	output := ConvertByte2String(data, GB18030)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		//过滤标头和ipv6，接口列表
		if strings.Contains(line, "Destination") || strings.Contains(line, "网关") || strings.Contains(line, "::") || strings.Contains(line, "......") {
			continue
		}
		//替换 在链路上为0.0.0.0
		line = strings.ReplaceAll(line, "在链路上", "0.0.0.0")
		fields := strings.Fields(line)
		//过滤杂乱信息
		if len(fields) != 5 && len(fields) != 4 {
			continue
		}
		routeList.List = append(routeList.List, RouteListItemRes{
			Ip:      fields[0],
			Netmask: fields[1],
			Gateway: fields[2],
		})
	}
	return routeList, nil
}

// WindowsRouteAdd 增加路由
func WindowsRouteAdd(ip, netmask, gateway string) error {
	name := "route"
	args := []string{"-p", "add", ip, "mask", netmask, gateway}
	addCmd := exec.Command(name, args...)
	data, err := addCmd.Output()
	if err != nil {
		fmt.Printf("exec add command err: %s\n", err)
		return err
	}
	output := string(data)
	if strings.Contains(output, "File exists") {
		return errors.New("路由已存在")
	}
	return nil
}

// WindowsRouteDelete 删除路由
func WindowsRouteDelete(ip, netmask, gateway string) error {
	name := "route"
	args := []string{"-p", "delete", ip, "mask", netmask, gateway}
	addCmd := exec.Command(name, args...)
	_, err := addCmd.Output()
	if err != nil {
		fmt.Printf("exec delete command err: %s\n", err)
		return err
	}
	return nil
}

// LinuxRouteQuery 查询路由表
func LinuxRouteQuery() (RouteListRes, error) {
	name := "route"
	arg := "-n"
	queryCmd := exec.Command(name, arg)
	var routeList RouteListRes
	data, err := queryCmd.Output()
	if err != nil {
		fmt.Printf("exec query command err: %s\n", err)
		return routeList, err
	}
	output := string(data)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Destination") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 8 {
			continue
		}
		routeList.List = append(routeList.List, RouteListItemRes{
			Ip:      fields[0],
			Gateway: fields[1],
			Netmask: fields[2],
		})
	}
	return routeList, nil
}

// LinuxRouteAdd 增加路由
func LinuxRouteAdd(ip, netmask, gateway string) error {
	name := "route"
	args := []string{"add", "-net", ip, "netmask", netmask, "gw", gateway}
	addCmd := exec.Command(name, args...)
	data, err := addCmd.Output()
	if err != nil {
		fmt.Printf("exec add command err: %s\n", err)
		return err
	}
	output := string(data)
	if strings.Contains(output, "File exists") {
		return errors.New("路由已存在")
	}
	return nil
}

// LinuxRouteDelete 删除路由
func LinuxRouteDelete(ip, netmask, gateway string) error {
	name := "route"
	args := []string{"del", "-net", ip, "netmask", netmask, "gw", gateway}
	addCmd := exec.Command(name, args...)
	_, err := addCmd.Output()
	if err != nil {
		fmt.Printf("exec delete command err: %s\n", err)
		return err
	}
	return nil
}

type SystemRoute struct {
}

// RouteQuery 查询路由表
func (s SystemRoute) RouteQuery() (RouteListRes, error) {
	var (
		err       error
		routeList RouteListRes
	)
	if runtime.GOOS == "linux" {
		routeList, err = LinuxRouteQuery()
	} else if runtime.GOOS == "windows" {
		routeList, err = WindowsRouteQuery()
	} else {
		routeList, err = MacRouteQuery()
	}
	return routeList, err
}

// RouteAdd 增加路由
func (s SystemRoute) RouteAdd(ip, netmask, gateway string) error {
	var err error
	if runtime.GOOS == "linux" {
		err = LinuxRouteAdd(ip, netmask, gateway)
	} else if runtime.GOOS == "windows" {
		err = WindowsRouteAdd(ip, netmask, gateway)
	} else {
		err = MacRouteAdd(ip, netmask, gateway)
	}
	return err
}

// RouteDelete 删除路由
func (s SystemRoute) RouteDelete(ip, netmask, gateway string) error {
	var err error
	if runtime.GOOS == "linux" {
		err = LinuxRouteDelete(ip, netmask, gateway)
	} else if runtime.GOOS == "windows" {
		err = WindowsRouteDelete(ip, netmask, gateway)
	} else {
		err = MacRouteDelete(ip, netmask, gateway)
	}
	return err
}
