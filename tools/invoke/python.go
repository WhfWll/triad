package invoke

//
//import (
//	"bufio"
//	"context"
//	"fmt"
//	log "github.com/sirupsen/logrus"
//	"io"
//	"os"
//	"os/exec"
//	"strings"
//	"time"
//)
//
//// invokePythonScript 进行yak脚本调用
//func (c *CallInfo) invokePythonScript(ctx context.Context, scriptType string, scriptContent string, callBackFunc func(context.Context, *CallInfo, string)) {
//	params := make([]string, 0)
//	for _, param := range c.ToolParamList {
//		params = append(params, "--"+param.ParamName)
//		//params = append(params, strings.ReplaceAll(param.ParamValue, " ", ""))
//		// 中电邮储项目中使用
//		params = append(params, strings.TrimSpace(param.ParamValue))
//	}
//	scriptFileName := c.saveScriptContent(ctx, scriptType, scriptContent)
//	var scriptParamList = []string{scriptFileName}
//	scriptParamList = append(scriptParamList, params...)
//	defer func() {
//		time.Sleep(1 * time.Second)
//		os.RemoveAll(scriptFileName)
//		callBackFunc(ctx, c, "end")
//	}()
//
//	ctxTimeout, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
//	defer cancel()
//	// 执行脚本
//	fmt.Println("[command]:", "python3", scriptParamList)
//	cmd := exec.CommandContext(ctxTimeout, "python3", scriptParamList...)
//	stdout, _ := cmd.StdoutPipe()
//	cmd.Stderr = cmd.Stdout
//	cmd.Start()
//	if cmd.Process != nil {
//		c.Pid = cmd.Process.Pid
//	}
//
//	buf := bufio.NewReader(stdout)
//	go func() {
//		for {
//			output, _, err := buf.ReadLine()
//			if err != nil {
//				if err == io.EOF {
//					break
//				}
//			}
//			if output == nil {
//				continue
//			}
//			result := strings.Trim(strings.Trim(strings.Trim(string(output), "\n"), "\r"), " ")
//			fmt.Println("result.........")
//			fmt.Println(result)
//			callBackFunc(ctx, c, result)
//		}
//	}()
//	time.Sleep(1 * time.Second)
//	if err := cmd.Wait(); err != nil {
//		log.Errorf("commmand execute error: %s", err)
//		return
//	}
//}
