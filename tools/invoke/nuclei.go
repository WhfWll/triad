package invoke

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

const nucleiTestTimeout = 300 * time.Second

// buildNucleiParam 进行nuceli脚本调用
func (c *CallInfo) buildNucleiParam(ctx context.Context, scriptParamList []string, fname string, params []string) ([]string, []string) {
	fnuclei, err := ioutil.TempFile("", "distributed-nucleicode-*.yak")
	tempFileList := make([]string, 0)
	if err != nil {
		fmt.Println(err)
		return scriptParamList, tempFileList
	}
	fnuclei.WriteString(nucleiExecutor)
	fnuclei.Close()
	tempFileList = append(tempFileList, fnuclei.Name())
	scriptParamList = append(scriptParamList, fnuclei.Name())
	scriptParamList = append(scriptParamList, "--pocFile")
	scriptParamList = append(scriptParamList, fname)
	scriptParamList = append(scriptParamList, params...)
	return scriptParamList, tempFileList
}

// ExecuteNucleiScriptForTest 通过 yak 包装器执行 nuclei 模板测试
func ExecuteNucleiScriptForTest(ctx context.Context, nucleiTemplatePath string, params []string) (string, error) {
	fnuclei, err := ioutil.TempFile("", "nuclei-executor-*.yak")
	if err != nil {
		return "", fmt.Errorf("创建nuclei执行脚本失败: %w", err)
	}
	defer os.Remove(fnuclei.Name())

	if _, err = fnuclei.WriteString(nucleiExecutor); err != nil {
		return "", fmt.Errorf("写入nuclei执行脚本失败: %w", err)
	}
	if err = fnuclei.Close(); err != nil {
		return "", fmt.Errorf("关闭nuclei执行脚本失败: %w", err)
	}

	var (
		results    []string
		resultLock sync.Mutex
	)

	yakitServer := NewYakitServer(
		0,
		SetYakitServer_LogHandler(func(level string, info string) {
			if info == "end" {
				return
			}
			resultLock.Lock()
			results = append(results, info)
			resultLock.Unlock()
		}),
	)
	yakitServer.Start()
	defer yakitServer.Shutdown()

	scriptParamList := []string{
		fnuclei.Name(),
		"--yakit-webhook", yakitServer.Addr(),
		"--pocFile", nucleiTemplatePath,
	}
	scriptParamList = append(scriptParamList, params...)

	ctxTimeout, cancel := context.WithTimeout(ctx, nucleiTestTimeout)
	defer cancel()

	yakCmd := "yak"
	if runtime.GOOS == "windows" {
		yakCmd = "yak.exe"
	}

	cmd := exec.CommandContext(ctxTimeout, yakCmd, scriptParamList...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("YAKIT_HOME=%v", os.Getenv("YAKIT_HOME")))

	fmt.Printf("[TestScript][Nuclei] 执行命令: %s %s\n", yakCmd, strings.Join(scriptParamList, " "))

	cmdOutput, err := cmd.CombinedOutput()

	// 等待 webhook 回调处理完成
	time.Sleep(1 * time.Second)

	resultLock.Lock()
	output := strings.Join(results, "\n")
	resultLock.Unlock()

	if len(cmdOutput) > 0 {
		if output != "" {
			output += "\n"
		}
		output += string(cmdOutput)
	}

	if err != nil {
		if ctxTimeout.Err() == context.DeadlineExceeded {
			return output, fmt.Errorf("脚本执行超时（300秒）")
		}
		if strings.TrimSpace(output) == "" {
			return string(cmdOutput), fmt.Errorf("命令执行出错: %w", err)
		}
		output += fmt.Sprintf("\n命令执行出错: %v", err)
	}

	if strings.TrimSpace(output) == "" {
		output = "[信息] Nuclei扫描完成，未发现漏洞"
	}
	return output, nil
}

var nucleiExecutor = `yakit.AutoInitYakit()
log.setLevel("info")
root_url := cli.String("root_url")
ip := cli.String("ip")
target = ""
if root_url != ""{
    target = root_url
}else if ip !=""{
    target = ip
}else{
    yakit.Info("info","缺少必要的参数")
}

pocFile := cli.String("pocFile", cli.setRequired(true))
isWorkflow = cli.Bool("isWorkflow")
debug = cli.Bool("debug")
proxy = cli.StringSlice("proxy")
cli.check()

client = yakit.NewClient(cli.String("yakit-webhook"))
cli.check()
yakit.Info("当前默认 yakit-webhook 为：%v", "")
yakit.Info("info", "参数检查成功。")

if (debug) {
    log.setLevel("debug")
}

if debug {
    yakit.Info("info", "构建基础扫描参数：调试模式")
} else {
    yakit.Info("info", "未开启调试模式")
}

fileRaw = io.ReadFile(pocFile)[0]
yakit.Info("info", "开始针对目标："+target+"，进行漏洞检测")
r, err := nuclei.Scan(target, nuclei.rawTemplate(fileRaw))
die(err)

for a = range r {
    yakit.Info("success", "监测到漏洞/风险["+a.PocName+"]："+a.Severity+" from: "+a.Target+"")
    //yakit.Info("json", a.RawJson)
	//yakit.Info(nuclei.PocVulToRisk(a))
    risk.NewRisk(
        a.Target, 
        risk.severity(a.severity),
        risk.title(a.PocName),
        risk.payload(a.payload),
        risk.request(a.details["request"]),
        risk.response(a.details["response"]),
        // risk.details(a.details)
    )
}

yakit.Info("end", "进程正常结束")
`
