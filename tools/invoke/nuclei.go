package invoke

import (
	"context"
	"fmt"
	"io/ioutil"
)

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
