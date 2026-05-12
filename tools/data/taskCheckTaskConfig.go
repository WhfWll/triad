package data

import (
	"errors"
	"smart/tools/enums"
	"strconv"
	"strings"
)

var TaskCheckTaskConfig taskCheckTaskConfig

type taskCheckTaskConfig struct {
}

// VerifyConfig 验证任务中的配置
func (c taskCheckTaskConfig) VerifyConfig(configJson enums.ConfigJson) (configStruct enums.ConfigJson, err error) {
	//验证端口扫描配置
	err = c.VerifyPortScanConfig(&configJson)
	if err != nil {
		return configJson, err
	}
	//验证动态爬虫配置
	err = c.VerifyWebCrawlerConfig(&configJson)
	if err != nil {
		return configJson, err
	}
	//验证路径爆破配置
	err = c.VerifyWebPathScanConfig(&configJson)
	if err != nil {
		return configJson, err
	}
	//验证站点登录凭证配置
	err = c.VerifyWebsiteLoginConfig(&configJson)
	if err != nil {
		return configJson, err
	}
	// 代理模式配置
	err = c.VerifyProxyConfig(&configJson)
	if err != nil {
		return configJson, err
	}
	// 验证横向移动配置
	err = c.VerifyLateralMoveConfig(&configJson)
	if err != nil {
		return configJson, err
	}
	return configJson, nil
}

// VerifyPortScanConfig 验证端口扫描配置
func (c taskCheckTaskConfig) VerifyPortScanConfig(configStruct *enums.ConfigJson) error {
	if !configStruct.PortScanConfig.IsOpen {
		return errors.New("端口扫描必须开启")
	}
	if configStruct.PortScanConfig.TCPScanType == 0 {
		return errors.New("TCP扫描类型不能为空")
	}
	if configStruct.PortScanConfig.ScanPort == "" {
		return errors.New("扫描端口不能为空")
	}
	if configStruct.PortScanConfig.Timeout <= 0 {
		configStruct.PortScanConfig.Timeout = enums.TaskConfigurationPortScanTimeoutDefault
	}
	if configStruct.PortScanConfig.Concurrent <= 0 {
		configStruct.PortScanConfig.Concurrent = enums.TaskConfigurationPortScanConcurrentDefault
	}
	if configStruct.PortScanConfig.Timeout > 3600 {
		return errors.New("端口扫描超时时间过大，请确认")
	}
	if configStruct.PortScanConfig.Concurrent > 1000 {
		return errors.New("端口扫描并发数过大，请确认")
	}

	scanPorts := strings.Split(configStruct.PortScanConfig.ScanPort, ",")
	for _, item := range scanPorts {
		checkList := make([]string, 0)
		if strings.Contains(item, "-") {
			itemSlice := strings.Split(item, "-")
			if len(itemSlice) != 2 {
				return errors.New("端口扫描[" + item + "]格式错误，请确认")
			}
			checkList = append(checkList, itemSlice[0])
			checkList = append(checkList, itemSlice[1])
		} else {
			checkList = append(checkList, item)
		}

		for _, port := range checkList {
			portInt, err := strconv.Atoi(port)
			if err != nil {
				return errors.New("端口扫描[" + port + "]格式错误，请确认")
			}
			if portInt < 0 {
				return errors.New("端口扫描[" + port + "]端口过小，请确认")
			}
			if portInt > 65535 {
				return errors.New("端口扫描[" + port + "]端口过大，请确认")
			}
		}
	}
	//var portScanConfig *enums.PortScanConfig
	//configStruct.PortScanConfig.PortScanTypeZh = portScanConfig.PortScanTypeEnum(configStruct.PortScanConfig.PortScanType)
	//configStruct.PortScanConfig.TCPScanTypeZh = portScanConfig.TcpScanTypeEnum(configStruct.PortScanConfig.TCPScanType)
	return nil
}

// VerifyWebCrawlerConfig 验证动态爬虫配置
func (c taskCheckTaskConfig) VerifyWebCrawlerConfig(configStruct *enums.ConfigJson) error {
	if !configStruct.WebCrawlerConfig.IsOpen {
		return nil
	}
	//var webCrawlerConfig *enums.WebCrawlerConfig
	//configStruct.WebCrawlerConfig.MaxDepthZh = webCrawlerConfig.CrawlerDeep(configStruct.WebCrawlerConfig.MaxDepth)
	//configStruct.WebCrawlerConfig.MaxUrlZh = webCrawlerConfig.CrawlerMaxConnect(configStruct.WebCrawlerConfig.MaxUrl)
	//configStruct.WebCrawlerConfig.ScanRangeZh = webCrawlerConfig.CrawlerRange(configStruct.WebCrawlerConfig.ScanRange)
	//configStruct.WebCrawlerConfig.TimeoutZh = webCrawlerConfig.SingleTimeout(configStruct.WebCrawlerConfig.Timeout)
	//configStruct.WebCrawlerConfig.FullTimeoutZh = webCrawlerConfig.FullTimeoutEnum(configStruct.WebCrawlerConfig.FullTimeout)
	//configStruct.WebCrawlerConfig.ScanRepeatZh = webCrawlerConfig.CrawlerRepeat(configStruct.WebCrawlerConfig.ScanRepeat)
	return nil
}

//// VerifyWeakPassConfig 弱口令爆破配置
//func (c taskCheckTaskConfig) VerifyWeakPassConfig(configStruct *enums.ConfigJson) error {
//	var weakPassConfig enums.WeakPassConfig
//	for _, serviceInt := range configStruct.WeakPassConfig.Services {
//		configStruct.WeakPassConfig.ServicesZh = append(configStruct.WeakPassConfig.ServicesZh, enums.GetDictionaryService(serviceInt))
//	}
//	configStruct.WeakPassConfig.DictTypeZh = weakPassConfig.WeakPassDictType(configStruct.WeakPassConfig.DictType)
//	configStruct.WeakPassConfig.GuessNumZh = weakPassConfig.WeakPassGuessNumber(configStruct.WeakPassConfig.GuessNum)
//	configStruct.WeakPassConfig.GuessTimeoutZh = weakPassConfig.WeakPassGuessTime(configStruct.WeakPassConfig.GuessTimeout)
//	configStruct.WeakPassConfig.GuessRateZh = weakPassConfig.WeakPassGuessRate(configStruct.WeakPassConfig.GuessRate)
//	return nil
//}

// VerifyWebPathScanConfig 验证路径爆破配置
func (c taskCheckTaskConfig) VerifyWebPathScanConfig(configStruct *enums.ConfigJson) error {
	if !configStruct.WebPathScanConfig.IsOpen {
		return nil
	}
	if configStruct.WebPathScanConfig.GuessRate == 0 {
		configStruct.WebPathScanConfig.GuessRate = 1
	}
	//if len(configStruct.WebPathScanConfig.ScanDict) == 0 {
	//	return errors.New("路径爆破-路径字典不能为空")
	//}
	//var webPathScanConfig *enums.WebPathScanConfig
	//configStruct.WebPathScanConfig.GuessRateZh = webPathScanConfig.WebPathScanSpeed(configStruct.WebPathScanConfig.GuessRate)
	//configStruct.WebPathScanConfig.GuessTimeoutZh = webPathScanConfig.WebPathScanTime(configStruct.WebPathScanConfig.GuessTimeout)
	return nil
}

// VerifyWebsiteLoginConfig 验证站点登录凭证配置
func (c taskCheckTaskConfig) VerifyWebsiteLoginConfig(configStruct *enums.ConfigJson) error {
	if !configStruct.WebsiteLoginConfig.IsOpen {
		return nil
	}
	for _, v := range configStruct.WebsiteLoginConfig.List {
		if v.Target == "" {
			return errors.New("站点登录凭证的登录地址不能为空")
		}
		if !strings.HasPrefix(v.Target, "http") {
			v.Target = "http" + "://" + v.Target
		}
		if v.VerifyType == 0 {
			return errors.New("站点登录凭证的凭证类型不能为空")
		}
		if v.VerifyValue == "" {
			return errors.New("站点登录凭证的凭证不能为空")
		}
		if v.VerifyStatus == 0 {
			return errors.New("站点登录凭证的状态不能为空")
		}
	}
	//var websiteLoginConfig *enums.WebsiteLoginConfig
	//for _, item := range configStruct.WebsiteLoginConfig.List {
	//	item.VerifyStatusZh = websiteLoginConfig.VerifyStatusZh(item.VerifyStatus)
	//	item.VerifyTypeZh = websiteLoginConfig.WebsiteLoginType(item.VerifyType)
	//}
	return nil
}

// VerifyProxyConfig 验证代理模式的配置
func (c taskCheckTaskConfig) VerifyProxyConfig(configStruct *enums.ConfigJson) error {
	if configStruct.ProxyConfig.IsOpen == true {
		if configStruct.ProxyConfig.Addr == "" {
			return errors.New("代理模式：代理IP不可为空")
		}
		if configStruct.ProxyConfig.Port == "" {
			return errors.New("代理模式：代理端口不可为空")
		}

		if configStruct.ProxyConfig.IsAuth == true {
			if configStruct.ProxyConfig.Username == "" {
				return errors.New("代理模式：认证账号不可为空")
			}
			if configStruct.ProxyConfig.Password == "" {
				return errors.New("代理模式：认证密码不可为空")
			}
		}
	}
	return nil
}

// VerifyLateralMoveConfig 验证横向移动配置
func (c taskCheckTaskConfig) VerifyLateralMoveConfig(configStruct *enums.ConfigJson) error {
	if !configStruct.LateralMove.IsOpen {
		return nil
	}
	if configStruct.LateralMove.Timeout < 0 {
		return errors.New("横向移动：超时时间不能小于0")
	}
	strategy := strings.TrimSpace(configStruct.LateralMove.Strategy)
	if strategy != "" &&
		strategy != enums.LateralMoveStrategySameSubnet &&
		strategy != enums.LateralMoveStrategyNeighbor &&
		strategy != enums.LateralMoveStrategyExcludeCurrent &&
		strategy != enums.LateralMoveStrategyCustomRange &&
		strategy != "auto_subnet" { // 保留 "auto_subnet" 兼容旧数据，视作 same_subnet
		return errors.New("横向移动：策略不合法")
	}
	if strategy == enums.LateralMoveStrategyCustomRange {
		if strings.TrimSpace(configStruct.LateralMove.Range) == "" {
			return errors.New("横向移动：自定义范围不能为空")
		}
	}
	ports := strings.TrimSpace(configStruct.LateralMove.Ports)
	if ports != "" {
		items := strings.Split(ports, ",")
		for _, it := range items {
			p := strings.TrimSpace(it)
			if p == "" {
				return errors.New("横向移动：端口格式错误")
			}
			val, err := strconv.Atoi(p)
			if err != nil {
				return errors.New("横向移动：端口格式错误")
			}
			if val <= 0 || val > 65535 {
				return errors.New("横向移动：端口范围不合法")
			}
		}
	}
	return nil
}
