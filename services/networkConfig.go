package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"smart/tools/enums"
	"smart/tools/file"
	"smart/tools/network"
	"strconv"
	"strings"
)

type NetWorkConfig struct {
}

// DockerComposeConfigUpdate 修改docker-compose配置文件的nginx的访问端口
func (n *NetWorkConfig) DockerComposeConfigUpdate(port string) error {
	//读取原有的yaml文件
	data, err := ioutil.ReadFile(enums.DockerComposeConfigPath)
	if err != nil {
		fmt.Println("read docker-compose config err: ", err)
		return err
	}
	//解析数据到结构体
	var dockerComposeConfig DockerComposeConfig
	if err = yaml.Unmarshal(data, &dockerComposeConfig); err != nil {
		fmt.Println("Unmarshal docker-compose config err: ", err)
		return err
	}
	//修改数据
	ports := make([]string, 0)
	if len(dockerComposeConfig.Services.Nginx.Ports) == 0 {
		return errors.New("nginx配置有误")
	}
	otherPortsList := make([]string, 0)
	if len(dockerComposeConfig.Services.Nginx.Ports) > 1 {
		otherPortsList = dockerComposeConfig.Services.Nginx.Ports[1:]
	}
	updatePortsStr := dockerComposeConfig.Services.Nginx.Ports[0]
	updatePortList := strings.Split(updatePortsStr, ":")
	if len(updatePortList) != 2 {
		return errors.New("nginx配置中端口格式有误")
	}
	updatePortList[0] = port
	ports = append(ports, strings.Join(updatePortList, ":"))
	if len(otherPortsList) > 0 {
		ports = append(ports, otherPortsList...)
	}
	dockerComposeConfig.Services.Nginx.Ports = ports
	//序列化为yaml格式
	outByte, err := yaml.Marshal(dockerComposeConfig)
	if err != nil {
		return err
	}
	//保存修改后的yaml文件
	if err = ioutil.WriteFile(enums.DockerComposeConfigPath, outByte, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// UbuntuNetworkConfigSave 保存ubuntu静态ip相关网络配置
func (n *NetWorkConfig) UbuntuNetworkConfigSave(ip, netmask, gateway, dns, sDns string) error {
	//获取默认网卡名称
	name, err := network.GetDefaultNetInterfaceName()
	if err != nil {
		return err
	}
	//确定配置文件路径
	filenameList, err := file.GetFilesFromDirectory(enums.UbuntuNetworkConfigDir)
	if err != nil {
		return err
	}
	var filePath string
	for _, item := range filenameList {
		if strings.HasSuffix(item, ".yaml") {
			filePath = item
			break
		}
	}
	if filePath == "" {
		filePath = filepath.Join(enums.UbuntuNetworkConfigDir, "50-cloud-init.yaml")
	}
	//转换子网掩码
	netmaskNumber, err := network.SubNetMaskToLen(netmask)
	if err != nil {
		return err
	}
	//组装数据
	addresses := []string{ip + "/" + strconv.Itoa(netmaskNumber)}
	nameServers := []string{dns}
	if sDns != "" {
		nameServers = append(nameServers, sDns)
	}
	data := map[string]interface{}{
		"network": map[string]interface{}{
			"ethernets": map[string]interface{}{
				name: map[string]interface{}{
					"dhcp4":     false,
					"addresses": addresses,
					"gateway4":  gateway,
					"nameservers": map[string]interface{}{
						"addresses": nameServers,
					},
				},
			},
			"version": 2,
		},
	}
	//序列化为yaml数据
	outByte, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	//保存数据到配置文件
	if err = ioutil.WriteFile(filePath, outByte, 644); err != nil {
		return err
	}
	return nil
}

// CentosNetworkConfigSave 保存centos静态ip相关网络配置
func (n *NetWorkConfig) CentosNetworkConfigSave(ip, netmask, gateway, dns, sDns string) error {
	//获取默认网卡名称
	name, err := network.GetDefaultNetInterfaceName()
	if err != nil {
		return err
	}
	//确定配置文件路径
	filenameList, err := file.GetFilesFromDirectory(enums.CentosNetworkConfigDir)
	if err != nil {
		return err
	}
	var filePath string
	for _, item := range filenameList {
		if strings.HasPrefix(item, "ifcfg-e") {
			filePath = item
			break
		}
	}
	if filePath == "" {
		filePath = filepath.Join(enums.CentosNetworkConfigDir, "ifcfg-"+name)
	}
	//组合数据
	originDataByte, err := file.ReadFile(filePath)
	if err != nil {
		return err
	}
	dataList := make([]string, 0)
	for _, line := range strings.Split(string(originDataByte), "\n") {
		if strings.Contains(line, "BOOTPROTO") {
			line = "BOOTPROTO=static"
		} else if strings.Contains(line, "ONBOOT") {
			line = "ONBOOT=yes"
		} else if strings.Contains(line, "IPADDR") {
			line = "IPADDR=" + ip
		} else if strings.Contains(line, "GATEWAY") {
			line = "GATEWAY=" + gateway
		} else if strings.Contains(line, "NETMASK") {
			line = "NETMASK=" + netmask
		} else if strings.Contains(line, "DNS1") {
			line = "DNS1=" + dns
		} else if strings.Contains(line, "DNS2") {
			continue
		}
		dataList = append(dataList, line)
	}
	if sDns != "" {
		dataList = append(dataList, "DNS2="+sDns)
	}
	//格式化数据
	data := strings.Join(dataList, "\n")
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//保存数据
	if err = file.WriteFile(filePath, dataByte); err != nil {
		return err
	}
	return nil
}

// RestartUbuntuNetwork 重启ubuntu的网卡
func (n *NetWorkConfig) RestartUbuntuNetwork() error {
	if _, err := exec.Command("netplan", "apply").Output(); err != nil {
		return err
	}
	return nil
}

// RestartCentosNetwork 重启centos的网卡
func (n *NetWorkConfig) RestartCentosNetwork() error {
	if _, err := exec.Command("systemctl", "restart", "network").Output(); err != nil {
		return err
	}
	return nil
}

// DockerComposeUpNginx 启动编排的容器nginx
func (n *NetWorkConfig) DockerComposeUpNginx() error {
	var err error
	if _, err = exec.Command("docker-compose", "-f", enums.DockerComposeConfigPath, "up", "-d", "nginx").Output(); err != nil {
		return err
	}
	return nil
}
