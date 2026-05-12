package invoke

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/config"
)

const ScannerErrorPrefix = "__scanner_error__:"
const defaultScannerProcessTimeout = 2 * time.Hour

type scannerRuntimeConfig struct {
	ProcessTimeoutSeconds int `json:"process_timeout_seconds"`
}

func formatScannerError(stage, target, scannerCmd string, err error) string {
	if scannerCmd != "" {
		return fmt.Sprintf("%s; target=%s; scanner=%s; error=%v", stage, target, scannerCmd, err)
	}
	return fmt.Sprintf("%s; target=%s; error=%v", stage, target, err)
}

func getScannerProcessTimeout() time.Duration {
	var scannerConfig scannerRuntimeConfig
	if err := config.Load("scanner", &scannerConfig); err != nil {
		return defaultScannerProcessTimeout
	}
	if scannerConfig.ProcessTimeoutSeconds <= 0 {
		return defaultScannerProcessTimeout
	}
	return time.Duration(scannerConfig.ProcessTimeoutSeconds) * time.Second
}

// 通用YAML配置模板常量
const configTemplate = `# 通用配置
scanUrl: {{.ScanURL}}
scanIp: {{.ScanIP}}
timeout: 30
maxDepth: 3
concurrency: 1
userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
cookie: {{.Cookie}}
headers:
{{- range $key, $value := .Headers}}
  "{{$key}}": "{{$value}}"
{{- end}}
proxy: "{{.ProxyConfig.URL}}"
webHook: {{.Webhook}}
crawlerMode: {{.CrawlerMode}}
scanMode: {{.ScanMode}}

# 脚本列表配置
scriptList:{{range .Scripts}}
  - "{{.}}"{{end}}

# 模块特定配置
# 无头浏览器爬虫配置（支持录制文件登录）
crawlerx:
  exePath: 
  wsAddress: ws://127.0.0.1:7317
  vue: false
  sensitiveWord: "logout,退出,登出"
  localStorage:
  loginUsername: "{{.LoginUsername}}"                           # 登录用户名
  loginPassword: "{{.LoginPassword}}"                           # 登录密码
  aiDomain: "{{.AiDomain}}"                                     # AI域名地址
  aiKey: "{{.AiKey}}"                                           # AI密钥
  headless: true  # 无头模式运行
  
  # 录制文件登录配置
  recordingFile:
    enabled: {{.RecordingFileConfig.Enabled}}  # 启用录制文件登录
    filePath: "{{.RecordingFileConfig.FilePath}}"  # 录制文件路径

# AI模型配置
aiModel:
  # 通用配置（保留超时与重试，域名与Key用于回退）
  timeout: {{.AiModel.Timeout}}
  maxRetries: {{.AiModel.MaxRetries}}

  # 文本识别模型配置（用于网页结构识别）
  textModel:
    name: "{{.AiModel.TextModel.Name}}"
    enabled: {{.AiModel.TextModel.Enabled}}
    domain: "{{.AiModel.TextModel.Domain}}"
    apiKey: "{{.AiModel.TextModel.ApiKey}}"

  # 视觉模型配置（用于验证码识别）
  visionModel:
    name: "{{.AiModel.VisionModel.Name}}"
    enabled: {{.AiModel.VisionModel.Enabled}} 
    domain: "{{.AiModel.VisionModel.Domain}}"
    apiKey: "{{.AiModel.VisionModel.ApiKey}}"

logConfig:
  enabled: true                           # 是否启用日志记录
  level: "debug"                          # 日志级别: debug, info, warn, error
  format: "text"                          # 日志格式: text, json
  filePath: "scannerlogs/{{.RandomFileName}}"     # 日志文件路径，为空则只输出到控制台
  maxSize: 100                            # 日志文件最大大小(MB)
  maxBackups: 3                           # 保留的日志文件数量
  maxAge: 7                               # 日志文件保留天数
  compress: true                          # 是否压缩旧日志文件
  enableConsole: false                    # 是否同时输出到控制台

portScan:
  enabled: {{.PortScanConfig.Enabled}}
  scanPort: {{.PortScanConfig.ScanPort}}
  protocol: {{.PortScanConfig.Protocol}}
  timeout: {{.PortScanConfig.Timeout}}
  concurrent: {{.PortScanConfig.Concurrent}}
  fingerprint:
    enabled: {{.PortScanConfig.Fingerprint.Enabled}}
    rulesFile: "{{.PortScanConfig.Fingerprint.RulesFile}}"

testMode:
  principleVerify: {{.TestModeConfig.PrincipleVerify}}
  versionMatch: {{.TestModeConfig.VersionMatch}}

# 原理验证脚本配置
principleVerifyScripts:
  enabled: true                               # 是否启用原理验证脚本过滤
{{- if .ScriptListFile}}
  scriptListFile: {{.ScriptListFile}}         # 外部脚本列表文件路径
{{- end}}
{{- if .ScriptNames}}
  scriptNames:                                # 脚本名称列表（分布式环境推荐）
{{- range .ScriptNames}}
    - "{{.}}"
{{- end}}
{{- end}}

# 后渗透配置
postExploit:
  enabled: {{.PostExploitEnable}}                      
  reverseShell:
    listenerHost: "{{.ReverseHost}}"                   
    listenerPort: "{{.ReversePort}}"                           
    shellType: "bash"                               
    timeout: 30                                     
  webshell:
    uploadPath: "/tmp/"                             
    shellFile: "shell.php"                          
    password: "rebeyond"                                

serviceBrute:
  enabled: {{.WeakPassConfig.Enabled}}
  timeout: {{.WeakPassConfig.Timeout}}
  targetsConcurrent: 5
  tasksConcurrent: 3
  guessNum: {{.WeakPassConfig.GuessNum}}
  guessTime: {{.WeakPassConfig.GuessTime}}
  guessRate: {{.WeakPassConfig.GuessRate}}
  enabledServices: [{{range $i, $service := .WeakPassConfig.Services}}{{if $i}}, {{end}}"{{$service}}"{{end}}]
  serviceConfigs:{{range $serviceName, $config := .WeakPassConfig.ServiceConfigs}}
    {{$serviceName}}:
      enabled: {{$config.Enabled}}
      timeout: {{$config.Timeout}}
      usernames: [{{range $i, $username := $config.Usernames}}{{if $i}}, {{end}}"{{$username}}"{{end}}]
      passwords: [{{range $i, $password := $config.Passwords}}{{if $i}}, {{end}}"{{$password}}"{{end}}]{{end}}
`

// InvokeScanner 调用漏洞扫描器
func InvokeScanner(ctx context.Context, target string, configJson enums.ConfigJson, vulLibraries []mysqls.VulLibraries, callBackFunc func(context.Context, string), customPortScanConfig ...*PortScanConfigResult) {
	yakitServer := NewYakitServer(
		0,
		SetYakitServer_LogHandler(func(level string, info string) {
			//fmt.Println("LEVEL:", level, "INFO: ", info)
			//通过回调函数保存数据
			callBackFunc(ctx, info)
		}),
	)
	yakitServer.Start()
	defer yakitServer.Shutdown()
	// 生成临时pocname列表文件
	tempScriptListFile, err := generatePocNameListFile(ctx, vulLibraries)
	if err != nil {
		log.Errorf("生成临时脚本列表文件失败: %v", err)
		tempScriptListFile = "" // 如果生成失败，使用空字符串
	}

	// 生成端口扫描配置，获取临时指纹文件路径
	var portScanConfig *PortScanConfigResult
	if len(customPortScanConfig) > 0 && customPortScanConfig[0] != nil {
		portScanConfig = customPortScanConfig[0]
	} else {
		portScanConfig = generatePortScanConfig(ctx, configJson)
	}

	configFileName := saveScannerConfig(ctx, yakitServer.Addr(), target, configJson, vulLibraries, tempScriptListFile, portScanConfig)
	var scriptParamList = []string{"scan", "-config", configFileName}
	scannerStarted := false

	defer func() {
		time.Sleep(1 * time.Second)
		os.RemoveAll(configFileName)
		// 清理临时脚本列表文件
		cleanupTempFile(tempScriptListFile)
		// 清理临时指纹规则文件
		if portScanConfig != nil && portScanConfig.TempRulesFile != "" {
			cleanupTempFile(portScanConfig.TempRulesFile)
		}
		if scannerStarted {
			callBackFunc(ctx, "end")
		}
	}()

	processTimeout := getScannerProcessTimeout()
	ctxTimeout, cancel := context.WithTimeout(ctx, processTimeout)
	defer cancel()
	// 获取当前工作目录并组装 scanner.exe 的完整路径
	workDir, err := os.Getwd()
	if err != nil {
		log.Errorf("获取工作目录失败: %v", err)
		callBackFunc(ctx, ScannerErrorPrefix+formatScannerError("prepare scanner command failed", target, "", err))
		return
	}
	scannerCmd := filepath.Join(workDir, "scanner")
	if runtime.GOOS == "windows" {
		scannerCmd = filepath.Join(workDir, "scanner.exe")
	}
	//fmt.Println("[command]:", scannerCmd, scriptParamList)
	log.Infof("invoke scanner process timeout: %s", processTimeout)
	cmd := exec.CommandContext(ctxTimeout, scannerCmd, scriptParamList...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//if cmd.Process != nil {
	//	c.Pid = cmd.Process.Pid
	//}
	err = cmd.Start()
	fmt.Println(cmd.String())
	if err != nil {
		log.Errorf("exec scanner %v failed: %s", configFileName, err)
		callBackFunc(ctx, ScannerErrorPrefix+formatScannerError("start scanner command failed", target, scannerCmd, err))
		return
	}
	scannerStarted = true
	time.Sleep(1 * time.Second)
	if err := cmd.Wait(); err != nil {
		log.Errorf("commmand execute error: %s", err)
		callBackFunc(ctx, ScannerErrorPrefix+formatScannerError("scanner command exited with error", target, scannerCmd, err))
		return
	}
}

// generatePocNameListFile 生成包含用户选定漏洞pocname的临时文件
func generatePocNameListFile(ctx context.Context, vulLibraries []mysqls.VulLibraries) (string, error) {
	// 筛选script_type为yak或nuclei的记录
	pocNames := make([]string, 0)
	for _, lib := range vulLibraries {
		if lib.Status != enums.VulLibrariesStatusSucess {
			continue
		}
		if (lib.ScriptType == "yak" || lib.ScriptType == "nuclei") && lib.Pocname != "" {
			pocNames = append(pocNames, lib.Pocname)
		}
	}
	if len(pocNames) == 0 {
		return "", fmt.Errorf("没有找到符合条件的pocname")
	}
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "principle_verify_scripts_*.txt")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer tempFile.Close()
	// 写入pocname列表，每行一个
	content := strings.Join(pocNames, "\n")
	if _, err := tempFile.WriteString(content); err != nil {
		os.Remove(tempFile.Name()) // 清理临时文件
		return "", fmt.Errorf("写入临时文件失败: %v", err)
	}
	log.Infof("生成临时脚本列表文件: %s，包含 %d 个pocname", tempFile.Name(), len(pocNames))
	return tempFile.Name(), nil
}

// generateRandomFileName 生成随机文件名
func generateRandomFileName() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return fmt.Sprintf("scanner_%d.log", time.Now().Unix())
	}
	return fmt.Sprintf("scanner_%s.log", hex.EncodeToString(bytes))
}

// cleanupTempFile 清理临时文件
func cleanupTempFile(filePath string) {
	if filePath != "" {
		if err := os.Remove(filePath); err != nil {
			log.Errorf("清理临时文件失败: %v", err)
		} else {
			log.Infof("已清理临时文件: %s", filePath)
		}
	}
}

// ReverseIpHost 反向连接配置结构体
type ReverseIpHost struct {
	ReverseType int    `json:"reverseType"`
	ReverseHost string `json:"reverseHost"`
	ReversePort int    `json:"reversePort"`
}

// getReverseHostAndPort 获取反向连接的主机和端口
func getReverseHostAndPort(ctx context.Context) (string, string, error) {
	var reverseHost = "192.168.0.68"
	var reversePort = "6666"

	var mapset mysqls.MapSet
	mapSetRes, err := mapset.GetsByObjKey(ctx, enums.ReverseIpHostMapSetObjKey)

	var objValue ReverseIpHost
	err = json.Unmarshal([]byte(mapSetRes.ObjValue), &objValue)

	if objValue.ReversePort != 0 {
		reversePort = strconv.Itoa(objValue.ReversePort)
	}

	if objValue.ReverseType == enums.TypeCustom && len(objValue.ReverseHost) != 0 { //自定义
		reverseHost = objValue.ReverseHost
	} else { //系统
		reverseHost, err = getReverseIp()
		if err != nil {
			return reverseHost, reversePort, err
		}
	}

	return reverseHost, reversePort, nil
}

// getReverseIp 获取系统IP地址
func getReverseIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var ip string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				if strings.HasPrefix(ip, "172.17") || strings.HasPrefix(ip, "172.18") || strings.HasPrefix(ip, "172.19") || strings.HasPrefix(ip, "169.254") {
					continue
				}
				break
			}
		}
	}
	return ip, nil
}

// ScanConfig 扫描配置结构体
type ScanConfig struct {
	ScanMode    string
	ScanURL     string
	ScanIP      string
	CrawlerMode string
	PostExploit bool
}

// parseScanConfig 解析扫描配置
func parseScanConfig(target string, configJson enums.ConfigJson) *ScanConfig {
	config := &ScanConfig{
		ScanURL: target,
	}

	// 根据target判断扫描模式
	if configJson.TestMode == "portscan" {
		config.ScanMode = "portscan"
		if parsedURL, err := url.Parse(target); err == nil {
			if host, _, err := net.SplitHostPort(parsedURL.Host); err == nil {
				config.ScanIP = host
			} else {
				config.ScanIP = parsedURL.Host
			}
		} else {
			config.ScanIP = target
		}
	} else if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		config.ScanMode = "web-xiaozhi"
		// 从URL中提取host
		if parsedURL, err := url.Parse(target); err == nil {
			// 从Host中去掉端口号，只保留IP
			if host, _, err := net.SplitHostPort(parsedURL.Host); err == nil {
				config.ScanIP = host
			} else {
				// 如果没有端口号，直接使用Host
				config.ScanIP = parsedURL.Host
			}
		} else {
			config.ScanIP = target
		}
	} else {
		config.ScanMode = "all"
		config.ScanIP = target
	}

	// 设置漏洞利用模式
	if configJson.VulExploit {
		config.PostExploit = true
	}

	// 设置爬虫模式
	if configJson.WebCrawlerConfig.IsOpen {
		config.CrawlerMode = "crawler"
	} else {
		config.CrawlerMode = "disable"
	}

	return config
}

// saveScriptContent 保存脚本源码函数
func saveScannerConfig(ctx context.Context, addr, target string, configJson enums.ConfigJson, vulLibraries []mysqls.VulLibraries, tempScriptListFile string, portScanConfig *PortScanConfigResult) string {
	// 解析扫描配置
	scanConfig := parseScanConfig(target, configJson)
	// 从工具参数中提取扫描URL
	if scanConfig.ScanURL == "" && scanConfig.ScanIP == "" {
		return ""
	}
	// 生成测试模式配置
	testModeConfig := generateTestModeConfig(ctx, configJson)
	// 调用弱口令配置生成函数
	weakPassConfig := generateWeakPassConfig(ctx, configJson)
	// 生成录制文件配置
	recordingFileConfig := generateRecordingFileConfig(ctx, target, configJson)
	if recordingFileConfig.Enabled {
		scanConfig.CrawlerMode = enums.ScriptNameCrawlerx
	}
	// 生成代理配置
	proxyConfig := generateProxyConfig(ctx, configJson)

	// 获取反向连接的主机和端口
	reverseHost, reversePort, err := getReverseHostAndPort(ctx)
	if err != nil {
		log.Errorf("获取反向连接配置失败: %v", err)
		// 使用默认值
		reverseHost = "192.168.0.68"
		reversePort = "6666"
	}

	debug, scriptParam, cookie := "false", "", ""
	loginUsername, loginPassword := "", ""
	headers := make(map[string]string)

	for _, script := range vulLibraries {
		if script.ScriptType == enums.VulScriptTypeUniversal {
			scriptParam += script.Pocname + ","
		}
	}
	if configJson.WebsiteLoginConfig.IsOpen {
		for _, node := range configJson.WebsiteLoginConfig.List {
			if strings.Contains(node.Target, target) || strings.Contains(target, node.Target) {
				if node.VerifyType == enums.TaskConfigurationWebsiteLoginCookie {
					cookie = node.VerifyValue
				} else if node.VerifyType == enums.TaskConfigurationWebsiteLoginAccount {
					// 处理账号密码类型，VerifyValue格式为 "用户名/密码"
					scanConfig.CrawlerMode = enums.ScriptNameCrawlerx
					if strings.Contains(node.VerifyValue, "/") {
						parts := strings.SplitN(node.VerifyValue, "/", 2)
						if len(parts) == 2 {
							loginUsername = parts[0]
							loginPassword = parts[1]
						}
					}
				} else if node.VerifyType == enums.TaskConfigurationWebsiteLoginHeader {
					var err error
					headers, err = utils.UrlHandleLogic.FormatHeadersNoHttpHeader(node.VerifyValue)
					if err != nil {
						log.Warnf("解析Header配置失败: %v", err)
					}
				}
			}
		}
	}

	// 加载AI模型配置（来源 llm_scenarios 与 llm_models）
	//aiModelCfg := generateAiModelConfig(ctx)
	aiModelCfg, isEnhanced := generateAiModelConfigNew(ctx)
	if isEnhanced {
		scanConfig.CrawlerMode = enums.ScriptNameCrawlerx
	}
	aiDomain := aiModelCfg.Domain
	aiKey := aiModelCfg.ApiKey

	// 生成配置文件内容
	configContent := fmt.Sprintf(configTemplate)
	tpl := template.Must(template.New("config").Parse(configContent))

	var buf bytes.Buffer
	// 处理脚本列表参数
	scripts := strings.Split(scriptParam, ",")
	if len(scripts) == 0 {
		scripts = []string{""} // 默认脚本
	}
	// 修改模板执行部分
	if err := tpl.Execute(&buf, struct {
		ScanURL             string
		ScanIP              string
		Debug               string
		Webhook             string
		Scripts             []string
		Cookie              string
		Headers             map[string]string
		ScanMode            string
		PostExploitEnable   bool
		RandomFileName      string
		ScriptListFile      string
		ScriptNames         []string
		TestModeConfig      *TestModeResult
		PortScanConfig      *PortScanConfigResult
		WeakPassConfig      *WeakPassConfig
		CrawlerMode         string
		RecordingFileConfig *RecordingFileConfig
		ProxyConfig         *ProxyConfig
		ReverseHost         string
		ReversePort         string
		LoginUsername       string
		LoginPassword       string
		AiDomain            string
		AiKey               string
		AiModel             *AiModelConfig
	}{
		ScanURL:             scanConfig.ScanURL,
		ScanIP:              scanConfig.ScanIP,
		Debug:               debug,
		Webhook:             addr,
		Scripts:             scripts,
		Cookie:              cookie,
		Headers:             headers,
		ScanMode:            scanConfig.ScanMode,
		PostExploitEnable:   scanConfig.PostExploit,
		RandomFileName:      generateRandomFileName(),
		ScriptNames:         make([]string, 0),
		ScriptListFile:      tempScriptListFile,
		TestModeConfig:      testModeConfig,
		PortScanConfig:      portScanConfig,
		WeakPassConfig:      weakPassConfig,
		CrawlerMode:         scanConfig.CrawlerMode,
		RecordingFileConfig: recordingFileConfig,
		ProxyConfig:         proxyConfig,
		ReverseHost:         reverseHost,
		ReversePort:         reversePort,
		LoginUsername:       loginUsername,
		LoginPassword:       loginPassword,
		AiDomain:            aiDomain,
		AiKey:               aiKey,
		AiModel:             aiModelCfg,
	}); err != nil {
		fmt.Println(err)
		log.Errorf("生成配置文件失败: %v", err)
		return ""
	}
	// 创建临时配置文件
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		log.Errorf("创建临时文件失败: %v", err)
		return ""
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		log.Errorf("写入配置文件失败: %v", err)
		return ""
	}
	f.Close()
	fmt.Println(f.Name())
	fmt.Println(string(buf.Bytes()))
	return f.Name()
}

// ServiceConfig 服务配置结构
type ServiceConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Timeout   int      `yaml:"timeout"`
	Usernames []string `yaml:"usernames"`
	Passwords []string `yaml:"passwords"`
}

// WeakPassConfig 弱口令配置结构
type WeakPassConfig struct {
	Enabled        bool                     `json:"enabled"`        // 是否启用弱口令扫描
	Services       []string                 `json:"services"`       // 启用的服务列表
	Timeout        int                      `json:"timeout"`        // 超时时间
	GuessNum       int                      `json:"guessNum"`       // 猜测次数
	GuessTime      int                      `json:"guessTime"`      // 猜测时间
	GuessRate      int                      `json:"guessRate"`      // 猜测速率
	ServiceConfigs map[string]ServiceConfig `json:"serviceConfigs"` // 各服务的详细配置
}

// PortScanConfigResult 端口扫描配置结果结构体
type PortScanConfigResult struct {
	Enabled       bool              `yaml:"enabled"`
	ScanPort      string            `yaml:"scanPort"`
	Protocol      string            `yaml:"protocol"`
	Timeout       int               `yaml:"timeout"`
	Concurrent    int               `yaml:"concurrent"`
	Fingerprint   FingerprintConfig `yaml:"fingerprint"`
	TempRulesFile string            `yaml:"-"` // 临时指纹规则文件路径，不序列化到YAML
}

// AiModelTextConfig 文本模型配置
type AiModelTextConfig struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
	Domain  string `yaml:"domain"`
	ApiKey  string `yaml:"apiKey"`
}

// AiModelVisionConfig 视觉模型配置
type AiModelVisionConfig struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
	Domain  string `yaml:"domain"`
	ApiKey  string `yaml:"apiKey"`
}

// AiModelConfig AI模型整体配置
type AiModelConfig struct {
	Domain      string              `yaml:"domain"`
	ApiKey      string              `yaml:"apiKey"`
	Timeout     int                 `yaml:"timeout"`
	MaxRetries  int                 `yaml:"maxRetries"`
	TextModel   AiModelTextConfig   `yaml:"textModel"`
	VisionModel AiModelVisionConfig `yaml:"visionModel"`
}

// generateAiModelConfig 从 llm_scenarios 和 llm_models 加载 AI 配置
// 约定：场景中 LlmModelID 指向默认文本模型；验证码识别场景名为 "验证码识别" 指向视觉模型；
// 若未找到或未启用，回退到默认模型；否则使用传入的固定域名与密钥等。
// tobe deleted
func generateAiModelConfig(ctx context.Context) *AiModelConfig {
	var (
		aiCfg = &AiModelConfig{
			Domain:      "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
			ApiKey:      "",
			Timeout:     30,
			MaxRetries:  3,
			TextModel:   AiModelTextConfig{Name: "", Enabled: true, Domain: "", ApiKey: ""},
			VisionModel: AiModelVisionConfig{Name: "", Enabled: true, Domain: "", ApiKey: ""},
		}
		scenarioModel mysqls.AiScenario
		modelModel    mysqls.LlmModel
	)

	// 加载所有场景
	scenarios, err := scenarioModel.GetAiScenarioList(ctx)
	if err != nil {
		log.Warnf("加载AI场景失败: %v", err)
	}
	fmt.Println("11111111111111")
	// 文本识别（网页结构）场景：名称包含 "网页结构识别" 时优先作为文本模型来源
	var textModelID int
	var visionModelID int
	for _, sc := range scenarios {
		if sc.IsEnabled == 1 {
			if strings.Contains(sc.Name, "网页结构") && sc.LlmModelID > 0 {
				textModelID = sc.LlmModelID
			}
			if strings.Contains(sc.Name, "验证码") && sc.LlmModelID > 0 {
				visionModelID = sc.LlmModelID
			}
		}
	}
	fmt.Println("2222222222222222")
	// 回退：若未绑定，尝试使用默认模型
	if textModelID == 0 {
		if def, err := modelModel.GetDefaultLlmModel(ctx, enums.LlmModelTypeText); err == nil && def.ID > 0 {
			textModelID = def.ID
		}
	}
	if visionModelID == 0 {
		if def, err := modelModel.GetDefaultLlmModel(ctx, enums.LlmModelTypeImg); err == nil && def.ID > 0 {
			visionModelID = def.ID
		}
	}
	fmt.Println("33333333333333333", textModelID, visionModelID)
	// 加载文本模型详情
	if textModelID > 0 {
		if m, err := modelModel.GetLlmModelByID(ctx, textModelID); err == nil && m.Status == 1 {
			aiCfg.TextModel.Name = m.ModelID
			aiCfg.TextModel.Enabled = true
			aiCfg.TextModel.Domain = m.ApiUrl
			aiCfg.TextModel.ApiKey = m.ApiKey
		}
	}
	// 加载视觉模型详情
	if visionModelID > 0 {
		if m, err := modelModel.GetLlmModelByID(ctx, visionModelID); err == nil && m.Status == 1 {
			aiCfg.VisionModel.Name = m.ModelID
			aiCfg.VisionModel.Enabled = true
			aiCfg.VisionModel.Domain = m.ApiUrl
			aiCfg.VisionModel.ApiKey = m.ApiKey
		}
	}
	fmt.Println("44444444444444")
	// 若依然为空，填充合理默认（通用回退）
	if aiCfg.ApiKey == "" {
		aiCfg.ApiKey = "sk-927357d4e4e64951a5317bb0efb8b5fb"
	}
	if aiCfg.TextModel.Name == "" {
		aiCfg.TextModel.Name = "qwen-plus"
	}
	if aiCfg.VisionModel.Name == "" {
		aiCfg.VisionModel.Name = "qwen-vl-max"
	}

	// 为每个模型填充回退的 domain/apiKey（当模型未设置时）
	if aiCfg.TextModel.Domain == "" {
		aiCfg.TextModel.Domain = aiCfg.Domain
	}
	if aiCfg.TextModel.ApiKey == "" {
		aiCfg.TextModel.ApiKey = aiCfg.ApiKey
	}
	if aiCfg.VisionModel.Domain == "" {
		aiCfg.VisionModel.Domain = aiCfg.Domain
	}
	if aiCfg.VisionModel.ApiKey == "" {
		aiCfg.VisionModel.ApiKey = aiCfg.ApiKey
	}

	return aiCfg
}

// generateAiModelConfigNew 新的ai参数加载函数 只从 llm_models 加载 AI 配置
func generateAiModelConfigNew(ctx context.Context) (*AiModelConfig, bool) {
	var (
		config      AiModelConfig
		mapSetModel mysqls.MapSet
		llmModel    mysqls.LlmModel
	)
	config.Timeout = 30
	config.MaxRetries = 3
	modelEnhance, err := mapSetModel.GetsByObjKey(ctx, enums.LlmModelEnhancementObjKey)
	if err != nil {
		return &config, false
	}
	modelEnhanceInt, _ := strconv.Atoi(modelEnhance.ObjValue)
	if modelEnhanceInt != enums.LlmModelEnhancementOpen {
		return &config, false
	}
	textModel, err := llmModel.GetDefaultLlmModel(ctx, enums.LlmModelTypeText)
	if err == nil && textModel.ID > 0 {
		config.TextModel = AiModelTextConfig{
			Name:    textModel.ModelName,
			Enabled: true,
			Domain:  textModel.ApiUrl,
			ApiKey:  textModel.ApiKey,
		}
	}
	visionModel, err := llmModel.GetDefaultLlmModel(ctx, enums.LlmModelTypeImg)
	if err == nil && visionModel.ID > 0 {
		config.VisionModel = AiModelVisionConfig{
			Name:    visionModel.ModelName,
			Enabled: true,
			Domain:  visionModel.ApiUrl,
			ApiKey:  visionModel.ApiKey,
		}
	}
	return &config, true
}

// FingerprintConfig 指纹识别配置
type FingerprintConfig struct {
	Enabled   bool   `yaml:"enabled"`
	RulesFile string `yaml:"rulesFile"`
}

func generatePortScanConfig(ctx context.Context, configJson enums.ConfigJson) *PortScanConfigResult {
	// 根据配置设置端口扫描协议
	var protocol string
	if configJson.PortScanConfig.TCPScanType == enums.TaskConfigurationTcpScanTypeSyn {
		protocol = "tcp_syn"
	} else if configJson.PortScanConfig.TCPScanType == enums.TaskConfigurationTcpScanTypeAck {
		protocol = "tcp_ack"
	} else if configJson.PortScanConfig.TCPScanType == enums.TaskConfigurationUDP {
		protocol = "udp"
	} else if configJson.PortScanConfig.TCPScanType == enums.TaskConfigurationTcpScanTypeFin {
		protocol = "tcp_fin"
	} else if configJson.PortScanConfig.TCPScanType == enums.TaskConfigurationTcpScanTypeNull {
		protocol = "tcp_null"
	} else {
		protocol = "tcp"
	}

	var finger mysqls.Finger
	fingerList, err := finger.GetFingerListBySource(ctx, enums.FingerSourceUser)
	if err != nil {
		log.Errorf("获取指纹列表失败: %v", err)
	}
	enable := false
	var rulesFile string
	if len(fingerList) > 0 {
		enable = true
		// 构建指纹内容，确保每条规则以 YAML 列表项形式保存
		var ruleContent string
		for _, item := range fingerList {
			if item.Flag == "" {
				continue
			}
			lines := strings.Split(item.Flag, "\n")
			for i, line := range lines {
				if i == 0 {
					trimmed := strings.TrimLeft(line, " \t")
					if trimmed == "" {
						continue
					}
					if !strings.HasPrefix(trimmed, "- ") {
						line = "- " + trimmed
					} else {
						line = trimmed
					}
				} else {
					if strings.TrimSpace(line) != "" {
						line = "  " + line
					}
				}
				ruleContent += line + "\n"
			}
		}

		// 创建临时指纹规则文件
		tempFile, err := os.CreateTemp("", "fingerprint-rules-*.yml")
		if err != nil {
			log.Errorf("创建临时指纹规则文件失败: %v", err)
			rulesFile = "fingerprint-rules.yml" // 使用默认文件名作为备选
		} else {
			defer tempFile.Close()
			// 写入指纹内容到临时文件
			if _, err := tempFile.WriteString(ruleContent); err != nil {
				log.Errorf("写入指纹规则内容失败: %v", err)
				rulesFile = "fingerprint-rules.yml" // 使用默认文件名作为备选
			} else {
				rulesFile = tempFile.Name()
				// Windows系统需要转换路径分隔符为双反斜杠
				if runtime.GOOS == "windows" {
					rulesFile = strings.ReplaceAll(rulesFile, "\\", "\\\\")
				}
				log.Infof("指纹规则文件已创建: %s, 包含 %d 条规则", rulesFile, len(fingerList))
			}
		}
	} else {
		rulesFile = "fingerprint-rules.yml" // 没有指纹规则时使用默认文件名
	}

	// 构建端口扫描配置结果
	portScanConfig := &PortScanConfigResult{
		Enabled:    configJson.PortScanConfig.IsOpen,
		ScanPort:   configJson.PortScanConfig.ScanPort,
		Protocol:   protocol,
		Timeout:    configJson.PortScanConfig.Timeout,
		Concurrent: configJson.PortScanConfig.Concurrent,
		Fingerprint: FingerprintConfig{
			Enabled:   enable,
			RulesFile: rulesFile,
		},
		TempRulesFile: "", // 初始化为空
	}

	// 如果创建了临时文件，记录路径用于后续清理
	if enable && rulesFile != "fingerprint-rules.yml" {
		portScanConfig.TempRulesFile = rulesFile
	}
	if portScanConfig.Timeout <= 0 {
		portScanConfig.Timeout = enums.TaskConfigurationPortScanTimeoutDefault
	}
	if portScanConfig.Concurrent <= 0 {
		portScanConfig.Concurrent = enums.TaskConfigurationPortScanConcurrentDefault
	}

	// 记录日志
	log.Infof("生成端口扫描配置: 协议=%s, 端口=%s, timeout=%d, concurrent=%d, 指纹规则文件=%s", protocol, configJson.PortScanConfig.ScanPort, portScanConfig.Timeout, portScanConfig.Concurrent, rulesFile)
	return portScanConfig
}

// generateWeakPassConfig 生成弱口令配置信息
func generateWeakPassConfig(ctx context.Context, configJson enums.ConfigJson) *WeakPassConfig {
	config := &WeakPassConfig{
		ServiceConfigs: make(map[string]ServiceConfig),
	}

	var dictModel mysqls.Dictionary
	if configJson.WeakPassConfig.IsOpen {
		config.Enabled = true
		serviceMap := enums.DictionaryServiceEnum()
		// 设置弱口令爆破参数
		config.Timeout = 300 // 默认超时时间
		config.GuessNum = configJson.WeakPassConfig.GuessNum
		config.GuessTime = configJson.WeakPassConfig.GuessTimeout
		config.GuessRate = configJson.WeakPassConfig.GuessRate

		// 预加载通用字典（如果是通用字典模式且非仅补充模式）
		var commonUsernames, commonPasswords []string
		if !configJson.WeakPassConfig.OnlyUseAdd && configJson.WeakPassConfig.DictType == 2 {
			// 加载通用用户字典
			if configJson.WeakPassConfig.CommonUserDict > 0 {
				if dict, err := dictModel.DictionaryRecord(ctx, configJson.WeakPassConfig.CommonUserDict); err == nil {
					if dict.Content != "" {
						lines := strings.Split(strings.TrimSpace(dict.Content), "\n")
						for _, line := range lines {
							if line = strings.TrimSpace(line); line != "" {
								commonUsernames = append(commonUsernames, line)
							}
						}
					}
				} else {
					log.Warnf("加载通用用户字典失败(ID=%d): %v", configJson.WeakPassConfig.CommonUserDict, err)
				}
			}
			// 加载通用密码字典
			if configJson.WeakPassConfig.CommonPassDict > 0 {
				if dict, err := dictModel.DictionaryRecord(ctx, configJson.WeakPassConfig.CommonPassDict); err == nil {
					if dict.Content != "" {
						lines := strings.Split(strings.TrimSpace(dict.Content), "\n")
						for _, line := range lines {
							if line = strings.TrimSpace(line); line != "" {
								commonPasswords = append(commonPasswords, line)
							}
						}
					}
				} else {
					log.Warnf("加载通用密码字典失败(ID=%d): %v", configJson.WeakPassConfig.CommonPassDict, err)
				}
			}
		}

		// 为每个服务生成独立的配置
		for _, serviceID := range configJson.WeakPassConfig.Services {
			serviceName, exists := serviceMap[serviceID]
			if !exists {
				continue
			}

			config.Services = append(config.Services, serviceName)
			var usernames, passwords []string

			// 1. 加载基础字典（非"仅使用补充字典"模式）
			if !configJson.WeakPassConfig.OnlyUseAdd {
				if configJson.WeakPassConfig.DictType == 2 {
					// --- 通用字典逻辑 ---
					usernames = append(usernames, commonUsernames...)
					passwords = append(passwords, commonPasswords...)
				} else {
					// --- 默认字典逻辑 ---
					// 获取用户名字典
					dictUserNameList, _, _ := dictModel.DictionaryList(ctx, 1, 1000, enums.DictionaryTypeUser, "")
					for _, dict := range dictUserNameList {
						if dict.Service == serviceID && dict.IsDefault == enums.DictionaryDefaultYes {
							// 解析字典内容，按换行符分割
							if dict.Content != "" {
								dictUsers := strings.Split(strings.TrimSpace(dict.Content), "\n")
								for _, user := range dictUsers {
									user = strings.TrimSpace(user)
									if user != "" {
										usernames = append(usernames, user)
									}
								}
							}
						}
					}
					// 获取密码字典
					dictPasswordList, _, _ := dictModel.DictionaryList(ctx, 1, 1000, enums.DictionaryTypePassword, "")
					for _, dict := range dictPasswordList {
						if dict.Service == serviceID && dict.IsDefault == enums.DictionaryDefaultYes {
							// 解析字典内容，按换行符分割
							if dict.Content != "" {
								dictPasswords := strings.Split(strings.TrimSpace(dict.Content), "\n")
								for _, pass := range dictPasswords {
									pass = strings.TrimSpace(pass)
									if pass != "" {
										passwords = append(passwords, pass)
									}
								}
							}
						}
					}
				}
			}

			// 2. 追加补充字典（仅在"仅使用补充字典"模式 或 "补充字典"来源模式下）
			if configJson.WeakPassConfig.OnlyUseAdd || configJson.WeakPassConfig.DictType == 3 {
				if configJson.WeakPassConfig.AddAccount != "" {
					addUsers := strings.Split(strings.TrimSpace(configJson.WeakPassConfig.AddAccount), "\n")
					for _, user := range addUsers {
						user = strings.TrimSpace(user)
						if user != "" {
							usernames = append(usernames, user)
						}
					}
				}
				if configJson.WeakPassConfig.AddPass != "" {
					addPasswords := strings.Split(strings.TrimSpace(configJson.WeakPassConfig.AddPass), "\n")
					for _, pass := range addPasswords {
						pass = strings.TrimSpace(pass)
						if pass != "" {
							passwords = append(passwords, pass)
						}
					}
				}
			}
			// 去重处理
			usernames = removeDuplicates(usernames)
			passwords = removeDuplicates(passwords)
			// 创建服务配置
			config.ServiceConfigs[serviceName] = ServiceConfig{
				Enabled:   true,
				Timeout:   config.Timeout,
				Usernames: usernames,
				Passwords: passwords,
			}
		}
		fmt.Printf("弱口令配置已启用: 服务=%v, 猜测次数=%d, 猜测时间=%d, 猜测速率=%d\n", config.Services, config.GuessNum, config.GuessTime, config.GuessRate)
		for serviceName, serviceConfig := range config.ServiceConfigs {
			fmt.Printf("服务 %s: 用户名数量=%d, 密码数量=%d\n", serviceName, len(serviceConfig.Usernames), len(serviceConfig.Passwords))
		}
	} else {
		config.Enabled = false
		config.Services = []string{}
		config.Timeout = 300
		config.GuessNum = 3
		config.GuessTime = 30
		config.GuessRate = 1
	}

	return config
}

// removeDuplicates 去除字符串切片中的重复元素
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// TestModeResult 测试模式配置结果结构体
type TestModeResult struct {
	PrincipleVerify bool `yaml:"principleVerify"`
	VersionMatch    bool `yaml:"versionMatch"`
}

// RecordingFileConfig 录制文件配置结构体
type RecordingFileConfig struct {
	Enabled  bool   `yaml:"enabled"`  // 是否启用录制文件登录
	FilePath string `yaml:"filePath"` // 录制文件路径
}

// ProxyConfig 代理配置结构体
type ProxyConfig struct {
	Enabled bool   `yaml:"enabled"` // 是否启用代理
	URL     string `yaml:"url"`     // 代理URL
}

// generateProxyConfig 生成代理配置
func generateProxyConfig(ctx context.Context, configJson enums.ConfigJson) *ProxyConfig {
	config := &ProxyConfig{
		Enabled: false,
		URL:     "",
	}

	// 检查是否启用了代理配置
	if configJson.ProxyConfig.IsOpen {
		config.Enabled = true

		// 根据代理协议类型构建代理URL
		var proxyURL string
		// 将Port字符串转换为整数
		port := 0
		if portInt, err := strconv.Atoi(configJson.ProxyConfig.Port); err == nil {
			port = portInt
		}

		switch configJson.ProxyConfig.Proto {
		case enums.ProxyConfigProtoHTTP:
			proxyURL = fmt.Sprintf("http://%s:%d", configJson.ProxyConfig.Addr, port)
		case enums.ProxyConfigProtoHTTPS:
			proxyURL = fmt.Sprintf("https://%s:%d", configJson.ProxyConfig.Addr, port)
		case enums.ProxyConfigProtoSOCKS4:
			// TCP和UDP都使用socks5协议
			proxyURL = fmt.Sprintf("socks4://%s:%d", configJson.ProxyConfig.Addr, port)
		case enums.ProxyConfigProtoSOCKS5:
			proxyURL = fmt.Sprintf("socks5://%s:%d", configJson.ProxyConfig.Addr, port)
		default:
			// 默认使用HTTP代理
			proxyURL = fmt.Sprintf("http://%s:%d", configJson.ProxyConfig.Addr, port)
		}

		// 如果启用了认证，添加认证信息
		if configJson.ProxyConfig.IsAuth {
			// 解析URL并添加认证信息
			if parsedURL, err := url.Parse(proxyURL); err == nil {
				parsedURL.User = url.UserPassword(configJson.ProxyConfig.Username, configJson.ProxyConfig.Password)
				proxyURL = parsedURL.String()
			}
		}

		config.URL = proxyURL
		log.Infof("已生成代理配置: %s", proxyURL)
	}

	return config
}

// generateRecordingFileConfig 生成录制文件配置
func generateRecordingFileConfig(ctx context.Context, target string, configJson enums.ConfigJson) *RecordingFileConfig {
	config := &RecordingFileConfig{
		Enabled:  false,
		FilePath: "recording_login_session.json",
	}

	// 检查是否启用了网站登录配置
	if !configJson.WebsiteLoginConfig.IsOpen {
		return config
	}

	// 查找匹配的LoginSequence配置
	for _, node := range configJson.WebsiteLoginConfig.List {
		// 检查目标是否匹配
		if strings.Contains(node.Target, target) || strings.Contains(target, node.Target) {
			// 检查是否为LoginSequence类型
			if node.VerifyType == enums.TaskConfigurationWebsiteLoginLoginSequence {
				config.Enabled = true

				// 创建临时录制文件
				tempFile, err := os.CreateTemp("", "recording_login_session_*.json")
				if err != nil {
					log.Errorf("创建临时录制文件失败: %v", err)
					return config
				}
				defer tempFile.Close()

				// 将verifyValue保存为录制文件
				recordingFilePath := tempFile.Name()
				if err := os.WriteFile(recordingFilePath, []byte(node.VerifyValue), 0644); err != nil {
					log.Errorf("保存录制文件失败: %v", err)
					config.Enabled = false
					return config
				}

				config.FilePath = recordingFilePath
				log.Infof("已生成录制文件: %s", recordingFilePath)
				break
			}
		}
	}

	return config
}

// generateTestModeConfig 生成测试模式配置
func generateTestModeConfig(ctx context.Context, configJson enums.ConfigJson) *TestModeResult {
	var principleVerify, versionMatch bool

	// 根据TestMode设置测试模式参数
	switch configJson.TestMode {
	case "1":
		principleVerify = true
		versionMatch = false
	case "2":
		principleVerify = false
		versionMatch = true
	case "1,2":
		principleVerify = true
		versionMatch = true
	default:
		principleVerify = false
		versionMatch = false
	}

	testModeConfig := &TestModeResult{
		PrincipleVerify: principleVerify,
		VersionMatch:    versionMatch,
	}

	// 记录日志
	log.Infof("生成测试模式配置: TestMode=%s, 原理验证=%t, 版本匹配=%t",
		configJson.TestMode, principleVerify, versionMatch)

	return testModeConfig
}

// BuildScanConfig 构建扫描配置YAML字符串
func BuildScanConfig(ctx context.Context, targetURL string, configJson enums.ConfigJson, vulLibraries []mysqls.VulLibraries) (string, error) {
	// 解析扫描配置
	scanConfig := parseScanConfig(targetURL, configJson)

	// 生成临时pocname列表文件
	tempScriptListFile, err := generatePocNameListFile(ctx, vulLibraries)
	if err != nil {
		log.Errorf("生成临时脚本列表文件失败: %v", err)
		tempScriptListFile = "" // 如果生成失败，使用空字符串
	}

	// 生成端口扫描配置，获取临时指纹文件路径
	portScanConfig := generatePortScanConfig(ctx, configJson)

	// 生成测试模式配置
	testModeConfig := generateTestModeConfig(ctx, configJson)

	// 调用弱口令配置生成函数
	weakPassConfig := generateWeakPassConfig(ctx, configJson)

	// 生成录制文件配置
	recordingFileConfig := generateRecordingFileConfig(ctx, targetURL, configJson)
	// 生成代理配置
	proxyConfig := generateProxyConfig(ctx, configJson)

	// 获取反向连接的主机和端口
	reverseHost, reversePort, err := getReverseHostAndPort(ctx)
	if err != nil {
		log.Errorf("获取反向连接配置失败: %v", err)
		// 使用默认值
		reverseHost = "192.168.0.68"
		reversePort = "6666"
	}

	// 使用通用配置模板常量
	// configTemplate 已在文件顶部定义为常量

	debug, scriptParam, cookie := "false", "", ""
	loginUsername, loginPassword := "", ""
	headers := make(map[string]string)

	for _, script := range vulLibraries {
		if script.ScriptType == enums.VulScriptTypeUniversal {
			scriptParam += script.Pocname + ","
		}
	}
	if configJson.WebsiteLoginConfig.IsOpen {
		for _, node := range configJson.WebsiteLoginConfig.List {
			if strings.Contains(node.Target, targetURL) || strings.Contains(targetURL, node.Target) {
				if node.VerifyType == enums.TaskConfigurationWebsiteLoginCookie {
					cookie = node.VerifyValue
				} else if node.VerifyType == enums.TaskConfigurationWebsiteLoginAccount {
					// 处理账号密码类型，VerifyValue格式为 "用户名/密码"
					scanConfig.CrawlerMode = enums.ScriptNameCrawlerx
					if strings.Contains(node.VerifyValue, "/") {
						parts := strings.SplitN(node.VerifyValue, "/", 2)
						if len(parts) == 2 {
							loginUsername = parts[0]
							loginPassword = parts[1]
						}
					}
				} else if node.VerifyType == enums.TaskConfigurationWebsiteLoginHeader {
					var err error
					headers, err = utils.UrlHandleLogic.FormatHeadersNoHttpHeader(node.VerifyValue)
					if err != nil {
						log.Warnf("解析Header配置失败: %v", err)
					}
				}
			}
		}
	}

	// 生成配置文件内容
	tpl := template.Must(template.New("config").Parse(configTemplate))

	var buf bytes.Buffer
	// 处理脚本列表参数
	scripts := strings.Split(scriptParam, ",")
	if len(scripts) == 0 {
		scripts = []string{""} // 默认脚本
	}

	// 提取脚本名称列表（用于分布式环境）
	var scriptNames []string
	for _, script := range vulLibraries {
		if script.ScriptType == enums.VulScriptTypeYak || script.ScriptType == enums.VulScriptTypeNuclei {
			scriptNames = append(scriptNames, script.Pocname)
		}
	}

	// 执行模板
	//aiModelCfg := generateAiModelConfig(ctx)
	aiModelCfg, isEnhanced := generateAiModelConfigNew(ctx)
	if isEnhanced {
		scanConfig.CrawlerMode = enums.ScriptNameCrawlerx
	}
	aiDomain := aiModelCfg.Domain
	aiKey := aiModelCfg.ApiKey

	if err = tpl.Execute(&buf, struct {
		ScanURL             string
		ScanIP              string
		Debug               string
		Webhook             string
		Scripts             []string
		Cookie              string
		Headers             map[string]string
		ScanMode            string
		PostExploitEnable   bool
		RandomFileName      string
		ScriptListFile      string
		ScriptNames         []string
		TestModeConfig      *TestModeResult
		PortScanConfig      *PortScanConfigResult
		WeakPassConfig      *WeakPassConfig
		CrawlerMode         string
		RecordingFileConfig *RecordingFileConfig
		ProxyConfig         *ProxyConfig
		ReverseHost         string
		ReversePort         string
		LoginUsername       string
		LoginPassword       string
		AiDomain            string
		AiKey               string
		AiModel             *AiModelConfig
	}{
		ScanURL:             scanConfig.ScanURL,
		ScanIP:              scanConfig.ScanIP,
		Debug:               debug,
		Webhook:             "",
		Scripts:             scripts,
		Cookie:              cookie,
		Headers:             headers,
		ScanMode:            scanConfig.ScanMode,
		PostExploitEnable:   scanConfig.PostExploit,
		RandomFileName:      generateRandomFileName(),
		ScriptListFile:      tempScriptListFile,
		ScriptNames:         scriptNames,
		TestModeConfig:      testModeConfig,
		PortScanConfig:      portScanConfig,
		WeakPassConfig:      weakPassConfig,
		CrawlerMode:         scanConfig.CrawlerMode,
		RecordingFileConfig: recordingFileConfig,
		ProxyConfig:         proxyConfig,
		ReverseHost:         reverseHost,
		ReversePort:         reversePort,
		LoginUsername:       loginUsername,
		LoginPassword:       loginPassword,
		AiDomain:            aiDomain,
		AiKey:               aiKey,
		AiModel:             aiModelCfg,
	}); err != nil {
		log.Errorf("生成配置文件失败: %v", err)
		return "", err
	}

	return buf.String(), nil
}
