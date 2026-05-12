package invoke

import (
	"context"
	"fmt"
	"io/ioutil"
)

// buildPortScanParam 进行yak脚本调用
func (c *CallInfo) buildPortScanParam(ctx context.Context, scriptParamList []string, fname string, params []string) ([]string, []string) {
	fport, err := ioutil.TempFile("", "distributed-portscancode-*.yak")
	tempFileList := make([]string, 0)
	if err != nil {
		fmt.Println(err)
		return scriptParamList, tempFileList
	}
	fport.WriteString(PORTSCAN_PLUGIN_TRIGGER_CODE)
	fport.Close()
	scriptParamList = append(scriptParamList, fport.Name())
	scriptParamList = append(scriptParamList, "--script-name")
	scriptParamList = append(scriptParamList, fname)
	scriptParamList = append(scriptParamList, params...)
	tempFileList = append(tempFileList, fport.Name())
	return scriptParamList, tempFileList
}

const PORTSCAN_PLUGIN_TRIGGER_CODE = `yakit.AutoInitYakit()
log.setLevel("info")

target = cli.String("ip", cli.setRequired(true))
port = cli.String("port", cli.setRequired(false))
name = cli.String("script-name", cli.setRequired(true))
cli.check()

println("TARGET: %v PORT: %v", target, port)
# input your yak code
res, err := servicescan.Scan(
    target, port, servicescan.active(false), servicescan.maxProbes(3), servicescan.probeTimeout(3),
    servicescan.databaseCache(true),
)
if err != nil {
    yakit.Error("服务扫描失败：%v", err)
    die(err)
}

hookManager := hook.NewManager()
err = hook.LoadYakitPluginFromConsole(
    hookManager, 
    name, 
    "handle")
if err != nil {
    yakit.Error("加载 Yak 插件失败：%v", err)
    die("no plugin loaded")
}

yakit.Info("开始执行服务扫描插件：%v", name)
for result = range res {
    if result.IsOpen() {
        yakit.Info("OPEN：%v", str.HostPort(result.Target, result.Port))
    } else {
        yakit.Info("CLOSED: %v", str.HostPort(result.Target, result.Port))
    }
    yakit.Info("扫描完成：%v，准备执行插件: %v", str.HostPort(result.Target, result.Port), name)
    log.info("扫描完成：%v，准备执行插件: %v", str.HostPort(result.Target, result.Port), name)
    hookManager.CallPluginKeyByName(name, "handle", result)
}
`
