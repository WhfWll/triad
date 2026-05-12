package invoke

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"smart/tools/config"
	"smart/tools/enums"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

// Falcon配置模板常量
const falconConfigTemplate = `# 通用配置
general:
  timeout: {{.General.Timeout}}                    # 默认超时时间（秒）
  concurrency: {{.General.Concurrency}}             # 并发数
  userAgent: "{{.General.UserAgent}}"
  proxy: "{{.General.Proxy}}"                      # 代理设置
  verbose: {{.General.Verbose}}                     # 是否启用详细输出
  maxRetries: {{.General.MaxRetries}}               # 最大重试次数
  delayBetween: {{.General.DelayBetween}}           # 请求间隔（毫秒）
  parentScriptID: "{{.General.ParentScriptID}}"    # 父节点ID，用于攻击路径追踪
  grandparentScriptID: "{{.General.GrandparentScriptID}}"    # 父节点ID，用于攻击路径追踪

# 日志配置
logConfig:
  enabled: {{.LogConfig.Enabled}}
  level: "{{.LogConfig.Level}}"                    # debug, info, warn, error
  format: "{{.LogConfig.Format}}"                  # text, json
  filePath: "{{.LogConfig.FilePath}}"
  maxSize: {{.LogConfig.MaxSize}}                   # MB
  maxBackups: {{.LogConfig.MaxBackups}}
  maxAge: {{.LogConfig.MaxAge}}                     # 天
  compress: {{.LogConfig.Compress}}
  enableConsole: {{.LogConfig.EnableConsole}}

# 要使用的漏洞利用器类型
exploiter: "{{.Exploiter}}"

# 目标配置
target:
  # Web漏洞利用目标配置
  web:
    url: "{{.Target.Web.URL}}"                     # Web目标URL
    headers:{{range $key, $value := .Target.Web.Headers}}
      "{{$key}}": "{{$value}}"{{end}}
    cookies: "{{.Target.Web.Cookies}}"
    followRedirects: {{.Target.Web.FollowRedirects}}
    verifySSL: {{.Target.Web.VerifySSL}}
  
  # 主机漏洞利用目标配置
  host:
    ip: "{{.Target.Host.IP}}"                      # 目标主机IP
    port: {{.Target.Host.Port}}                     # 目标端口
    protocol: "{{.Target.Host.Protocol}}"          # 协议类型

# 漏洞利用器配置
exploiters:
  # SQL注入配置
  sqlInjection:
    timeout: {{.Exploiters.SQLInjection.Timeout}}
    parameter: "{{.Exploiters.SQLInjection.Parameter}}"  # 注入参数名

  # 永恒之蓝配置
  eternalBlue:
    enabled: {{.Exploiters.EternalBlue.Enabled}}
    timeout: {{.Exploiters.EternalBlue.Timeout}}
    defaultPort: {{.Exploiters.EternalBlue.DefaultPort}}
    msfConfig:
      host: "{{.Exploiters.EternalBlue.MSFConfig.Host}}"
      port: {{.Exploiters.EternalBlue.MSFConfig.Port}}
      user: "{{.Exploiters.EternalBlue.MSFConfig.User}}"
      pass: "{{.Exploiters.EternalBlue.MSFConfig.Pass}}"
    payloadConfig:
      lhost: "{{.Exploiters.EternalBlue.PayloadConfig.LHost}}"  # 本地监听IP
      lport: {{.Exploiters.EternalBlue.PayloadConfig.LPort}}     # 本地监听端口
      payload: "{{.Exploiters.EternalBlue.PayloadConfig.Payload}}"
    exploitOptions:
      target: {{.Exploiters.EternalBlue.ExploitOptions.Target}}
      maxExploitAttempts: {{.Exploiters.EternalBlue.ExploitOptions.MaxExploitAttempts}}

# 后渗透配置
postExploit:
  enabled: {{.PostExploit.Enabled}}

  # 反向Shell配置
  reverseShell:
    listenerHost: "{{.PostExploit.ReverseShell.ListenerHost}}"
    listenerPort: {{.PostExploit.ReverseShell.ListenerPort}}
    shellType: "{{.PostExploit.ReverseShell.ShellType}}"     # bash, powershell, cmd
    timeout: {{.PostExploit.ReverseShell.Timeout}}
    autoStart: {{.PostExploit.ReverseShell.AutoStart}}

  # Webshell配置
  webshell:
    uploadPath: "{{.PostExploit.Webshell.UploadPath}}"
    shellFile: "{{.PostExploit.Webshell.ShellFile}}"
    password: "{{.PostExploit.Webshell.Password}}"

  # 横向移动配置
  lateral:
    enabled: {{.PostExploit.Lateral.Enabled}}
    scanInternalNetwork: {{.PostExploit.Lateral.ScanInternalNetwork}}
    credentialHarvesting: {{.PostExploit.Lateral.CredentialHarvesting}}
    networkRange: "{{.PostExploit.Lateral.NetworkRange}}"

  # 清理配置
  cleanup:
    enabled: {{.PostExploit.Cleanup.Enabled}}
    autoCleanup: {{.PostExploit.Cleanup.AutoCleanup}}
    cleanupDelay: {{.PostExploit.Cleanup.CleanupDelay}}       # 秒
    removeUploaded: {{.PostExploit.Cleanup.RemoveUploaded}}
    clearLogs: {{.PostExploit.Cleanup.ClearLogs}}

# 输出配置
output:
  format: "{{.Output.Format}}"                     # text, json, xml
  saveResults: {{.Output.SaveResults}}
  resultsPath: "{{.Output.ResultsPath}}"
  includeEvidence: {{.Output.IncludeEvidence}}
  includeRawResponse: {{.Output.IncludeRawResponse}}
  timestampFormat: "{{.Output.TimestampFormat}}"

# 安全配置
security:
  maxConcurrentTargets: {{.Security.MaxConcurrentTargets}}
  rateLimiting:
    enabled: {{.Security.RateLimiting.Enabled}}
    requestsPerSecond: {{.Security.RateLimiting.RequestsPerSecond}}
    burstSize: {{.Security.RateLimiting.BurstSize}}
  requireConfirmation: {{.Security.RequireConfirmation}}     # 执行前需要确认

# 通知配置
notifications:
  enabled: {{.Notifications.Enabled}}
  # Webhook通知
  webhook:
    url: "{{.Notifications.Webhook.URL}}"
    method: "{{.Notifications.Webhook.Method}}"
    headers:
      "Content-Type": "application/json"
`

// InvokeFalcon 调用falcon漏洞利用器
func InvokeFalcon(ctx context.Context, target string, configJson enums.ConfigJson, scannerLog *ScannerLog, callBackFunc func(context.Context, string)) {
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

	// 生成falcon配置文件
	configFileName := saveFalconConfig(ctx, yakitServer.Addr(), target, configJson, scannerLog)
	if configFileName == "" {
		callBackFunc(ctx, "生成falcon配置文件失败")
		return
	}

	var scriptParamList = []string{"exploit", "--config", configFileName}
	defer func() {
		time.Sleep(1 * time.Second)
		os.RemoveAll(configFileName)
		callBackFunc(ctx, "end")
	}()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 1800*time.Second) // 30分钟超时
	defer cancel()

	// 获取当前工作目录并组装 falcon.exe 的完整路径
	workDir, err := os.Getwd()
	if err != nil {
		log.Errorf("获取工作目录失败: %v", err)
		return
	}

	falconCmd := filepath.Join(workDir, "falcon")
	if runtime.GOOS == "windows" {
		falconCmd = filepath.Join(workDir, "falcon.exe")
	}

	fmt.Printf("[falcon command]: %s %s\n", falconCmd, strings.Join(scriptParamList, " "))
	cmd := exec.CommandContext(ctxTimeout, falconCmd, scriptParamList...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Errorf("exec falcon %v failed: %s", configFileName, err)
		callBackFunc(ctx, fmt.Sprintf("falcon执行失败: %s", err))
		return
	}

	callBackFunc(ctx, "falcon漏洞利用器开始执行")
	time.Sleep(1 * time.Second)

	if err := cmd.Wait(); err != nil {
		log.Errorf("falcon command execute error: %s", err)
		callBackFunc(ctx, fmt.Sprintf("falcon执行错误: %s", err))
		return
	}

	callBackFunc(ctx, "falcon漏洞利用器执行完成")
}

// generateRandomFileName 生成随机文件名
func generateFalconRandomFileName() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("falcon_%d.yaml", time.Now().Unix())
	}
	return fmt.Sprintf("falcon_%s.yaml", hex.EncodeToString(bytes))
}

// saveFalconConfig 保存falcon配置文件
func saveFalconConfig(ctx context.Context, yakitServerAddr string, target string, configJson enums.ConfigJson, scannerLog *ScannerLog) string {
	// 生成falcon配置
	falconConfig := generateFalconConfig(yakitServerAddr, target, configJson, scannerLog)
	// 生成配置文件内容
	tpl := template.Must(template.New("falconConfig").Parse(falconConfigTemplate))

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, falconConfig); err != nil {
		log.Errorf("生成falcon配置文件失败: %v", err)
		return ""
	}

	// 创建临时配置文件
	f, err := os.CreateTemp("", "falcon-config-*.yaml")
	if err != nil {
		log.Errorf("创建临时文件失败: %v", err)
		return ""
	}

	if _, err := f.Write(buf.Bytes()); err != nil {
		log.Errorf("写入配置文件失败: %v", err)
		return ""
	}
	f.Close()

	fmt.Printf("[falcon config]: %s\n", f.Name())
	fmt.Println(string(buf.Bytes()))
	return f.Name()
}

// FalconConfig falcon配置结构体
type FalconConfig struct {
	General       GeneralConfig       `yaml:"general"`
	LogConfig     LogConfig           `yaml:"logConfig"`
	Exploiter     string              `yaml:"exploiter"`
	Target        TargetConfig        `yaml:"target"`
	Exploiters    ExploitersConfig    `yaml:"exploiters"`
	PostExploit   PostExploitConfig   `yaml:"postExploit"`
	Output        OutputConfig        `yaml:"output"`
	Security      SecurityConfig      `yaml:"security"`
	Notifications NotificationsConfig `yaml:"notifications"`
}

type GeneralConfig struct {
	Timeout             int    `yaml:"timeout"`
	Concurrency         int    `yaml:"concurrency"`
	UserAgent           string `yaml:"userAgent"`
	Proxy               string `yaml:"proxy"`
	Verbose             bool   `yaml:"verbose"`
	MaxRetries          int    `yaml:"maxRetries"`
	DelayBetween        int    `yaml:"delayBetween"`
	ParentScriptID      string `yaml:"parentScriptID"`
	GrandparentScriptID string `yaml:"grandparentScriptID"`
}

type LogConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Level         string `yaml:"level"`
	Format        string `yaml:"format"`
	FilePath      string `yaml:"filePath"`
	MaxSize       int    `yaml:"maxSize"`
	MaxBackups    int    `yaml:"maxBackups"`
	MaxAge        int    `yaml:"maxAge"`
	Compress      bool   `yaml:"compress"`
	EnableConsole bool   `yaml:"enableConsole"`
}

type TargetConfig struct {
	Web  WebTargetConfig  `yaml:"web"`
	Host HostTargetConfig `yaml:"host"`
}

type WebTargetConfig struct {
	URL             string            `yaml:"url"`
	Headers         map[string]string `yaml:"headers"`
	Cookies         string            `yaml:"cookies"`
	FollowRedirects bool              `yaml:"followRedirects"`
	VerifySSL       bool              `yaml:"verifySSL"`
}

type HostTargetConfig struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

type ExploitersConfig struct {
	SQLInjection SQLInjectionConfig `yaml:"sqlInjection"`
	EternalBlue  EternalBlueConfig  `yaml:"eternalBlue"`
}

type SQLInjectionConfig struct {
	Timeout   int    `yaml:"timeout"`
	Parameter string `yaml:"parameter"`
}

type EternalBlueConfig struct {
	Enabled        bool                 `yaml:"enabled"`
	Timeout        int                  `yaml:"timeout"`
	DefaultPort    int                  `yaml:"defaultPort"`
	MSFConfig      MSFConfig            `yaml:"msfConfig"`
	PayloadConfig  PayloadConfig        `yaml:"payloadConfig"`
	ExploitOptions ExploitOptionsConfig `yaml:"exploitOptions"`
}

type MSFConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type PayloadConfig struct {
	LHost   string `yaml:"lhost"`
	LPort   int    `yaml:"lport"`
	Payload string `yaml:"payload"`
}

type ExploitOptionsConfig struct {
	Target             int `yaml:"target"`
	MaxExploitAttempts int `yaml:"maxExploitAttempts"`
}

type PostExploitConfig struct {
	Enabled      bool               `yaml:"enabled"`
	ReverseShell ReverseShellConfig `yaml:"reverseShell"`
	Webshell     WebshellConfig     `yaml:"webshell"`
	Lateral      LateralConfig      `yaml:"lateral"`
	Cleanup      CleanupConfig      `yaml:"cleanup"`
}

type ReverseShellConfig struct {
	ListenerHost string `yaml:"listenerHost"`
	ListenerPort int    `yaml:"listenerPort"`
	ShellType    string `yaml:"shellType"`
	Timeout      int    `yaml:"timeout"`
	AutoStart    bool   `yaml:"autoStart"`
}

type WebshellConfig struct {
	UploadPath string `yaml:"uploadPath"`
	ShellFile  string `yaml:"shellFile"`
	Password   string `yaml:"password"`
}

type LateralConfig struct {
	Enabled              bool   `yaml:"enabled"`
	ScanInternalNetwork  bool   `yaml:"scanInternalNetwork"`
	CredentialHarvesting bool   `yaml:"credentialHarvesting"`
	NetworkRange         string `yaml:"networkRange"`
}

type CleanupConfig struct {
	Enabled        bool `yaml:"enabled"`
	AutoCleanup    bool `yaml:"autoCleanup"`
	CleanupDelay   int  `yaml:"cleanupDelay"`
	RemoveUploaded bool `yaml:"removeUploaded"`
	ClearLogs      bool `yaml:"clearLogs"`
}

type OutputConfig struct {
	Format             string `yaml:"format"`
	SaveResults        bool   `yaml:"saveResults"`
	ResultsPath        string `yaml:"resultsPath"`
	IncludeEvidence    bool   `yaml:"includeEvidence"`
	IncludeRawResponse bool   `yaml:"includeRawResponse"`
	TimestampFormat    string `yaml:"timestampFormat"`
}

type SecurityConfig struct {
	MaxConcurrentTargets int                `yaml:"maxConcurrentTargets"`
	RateLimiting         RateLimitingConfig `yaml:"rateLimiting"`
	RequireConfirmation  bool               `yaml:"requireConfirmation"`
}

type RateLimitingConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerSecond int  `yaml:"requestsPerSecond"`
	BurstSize         int  `yaml:"burstSize"`
}

type NotificationsConfig struct {
	Enabled bool          `yaml:"enabled"`
	Webhook WebhookConfig `yaml:"webhook"`
}

type WebhookConfig struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
}

// generateFalconConfig 生成falcon配置
func generateFalconConfig(yakitServerAddr string, target string, configJson enums.ConfigJson, scannerLog *ScannerLog) *FalconConfig {
	// 解析目标URL和IP
	var parameter string
	if scannerLog.Result["location"] != nil {
		target = scannerLog.Result["location"].(string)
	}
	if scannerLog.Result["parameter"] != nil {
		parameter = scannerLog.Result["parameter"].(string)
	}

	webURL := target
	hostIP := target
	hostPort := 80
	protocol := "tcp"
	// 根据target类型设置不同的配置
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		webURL = target
		// 从URL中提取IP
		if strings.Contains(target, "://") {
			parts := strings.Split(target, "://")
			if len(parts) > 1 {
				urlPart := parts[1]
				if strings.Contains(urlPart, "/") {
					urlPart = strings.Split(urlPart, "/")[0]
				}
				if strings.Contains(urlPart, ":") {
					hostIP = strings.Split(urlPart, ":")[0]
				} else {
					hostIP = urlPart
				}
			}
		}
	} else {
		hostIP = target
		webURL = fmt.Sprintf("http://%s", target)
	}

	// 根据POC名称获取漏洞利用配置
	exploiterType := "" // 默认
	if pocConfig := config.GetPocExploitConfig(scannerLog.Pocname); pocConfig != nil {
		exploiterType = pocConfig.ExploiterType
		if pocConfig.HostPort > 0 {
			hostPort = pocConfig.HostPort
		}
		if pocConfig.Protocol != "" {
			protocol = pocConfig.Protocol
		}
	}

	// 确定漏洞利用器类型
	// 获取Cookie配置
	cookie := ""
	if configJson.WebsiteLoginConfig.IsOpen {
		for _, node := range configJson.WebsiteLoginConfig.List {
			if strings.Contains(node.Target, target) || strings.Contains(target, node.Target) {
				if node.VerifyType == enums.TaskConfigurationWebsiteLoginCookie {
					cookie = node.VerifyValue
				}
			}
		}
	}
	// 横向移动配置
	lateralEnabled, scanInternalNetwork, credentialHarvesting := false, false, false
	networkRange := "192.168.1.0/24"
	reverseHost, _ := getReverseIp()
	msfApi := "127.0.0.1"

	if configJson.LateralMove.IsOpen {
		lateralEnabled = true
		// 根据配置设置横向移动参数
		scanInternalNetwork = false
		credentialHarvesting = true
		networkRange = "192.168.4.1/24"
	}

	return &FalconConfig{
		General: GeneralConfig{
			Timeout:             30,
			Concurrency:         5,
			UserAgent:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			Proxy:               "",
			Verbose:             false,
			MaxRetries:          3,
			DelayBetween:        100,
			ParentScriptID:      scannerLog.ScriptExecutionID,
			GrandparentScriptID: scannerLog.ParentScriptID,
		},
		LogConfig: LogConfig{
			Enabled:       true,
			Level:         "info",
			Format:        "text",
			FilePath:      fmt.Sprintf("logs/falcon_%s.log", generateFalconRandomFileName()),
			MaxSize:       100,
			MaxBackups:    5,
			MaxAge:        30,
			Compress:      true,
			EnableConsole: true,
		},
		Exploiter: exploiterType,
		Target: TargetConfig{
			Web: WebTargetConfig{
				URL: webURL,
				Headers: map[string]string{
					"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
					"Accept-Language": "en-US,en;q=0.5",
					"Accept-Encoding": "gzip, deflate",
				},
				Cookies:         cookie,
				FollowRedirects: true,
				VerifySSL:       false,
			},
			Host: HostTargetConfig{
				IP:       hostIP,
				Port:     hostPort,
				Protocol: protocol,
			},
		},
		Exploiters: ExploitersConfig{
			SQLInjection: SQLInjectionConfig{
				Timeout:   45,
				Parameter: parameter, // 默认参数名
			},
			EternalBlue: EternalBlueConfig{
				Enabled:     true,
				Timeout:     120,
				DefaultPort: 445,
				MSFConfig: MSFConfig{
					Host: msfApi,
					Port: 55553,
					User: "msf",
					Pass: "password",
				},
				PayloadConfig: PayloadConfig{
					LHost:   reverseHost,
					LPort:   4444,
					Payload: "windows/x64/meterpreter/reverse_tcp",
				},
				ExploitOptions: ExploitOptionsConfig{
					Target:             0,
					MaxExploitAttempts: 3,
				},
			},
		},
		PostExploit: PostExploitConfig{
			Enabled: configJson.VulExploit, // 根据配置决定是否启用后渗透
			ReverseShell: ReverseShellConfig{
				ListenerHost: "192.168.1.100",
				ListenerPort: 4444,
				ShellType:    "bash",
				Timeout:      300,
				AutoStart:    false,
			},
			Webshell: WebshellConfig{
				UploadPath: "/var/www/html/",
				ShellFile:  "shell.php",
				Password:   "admin123",
			},
			Lateral: LateralConfig{
				Enabled:              lateralEnabled,
				ScanInternalNetwork:  scanInternalNetwork,
				CredentialHarvesting: credentialHarvesting,
				NetworkRange:         networkRange,
			},
			Cleanup: CleanupConfig{
				Enabled:        true,
				AutoCleanup:    false,
				CleanupDelay:   300,
				RemoveUploaded: true,
				ClearLogs:      false,
			},
		},
		Output: OutputConfig{
			Format:             "text",
			SaveResults:        true,
			ResultsPath:        "results/",
			IncludeEvidence:    true,
			IncludeRawResponse: false,
			TimestampFormat:    "2006-01-02 15:04:05",
		},
		Security: SecurityConfig{
			MaxConcurrentTargets: 10,
			RateLimiting: RateLimitingConfig{
				Enabled:           true,
				RequestsPerSecond: 5,
				BurstSize:         10,
			},
			RequireConfirmation: true,
		},
		Notifications: NotificationsConfig{
			Enabled: true,
			Webhook: WebhookConfig{
				URL:    yakitServerAddr,
				Method: "POST",
			},
		},
	}
}
