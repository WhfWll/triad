package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/file"
	"os/exec"
	"path"
)

var (
	xrayExecFile, xrayExecConfigFile, outputDir string
)

type XRayItem struct {
	CreateTime int64 `json:"create_time"`
	Detail     struct {
		Addr     string      `json:"addr"`
		Payload  string      `json:"payload"`
		Snapshot [][]string  `json:"snapshot"`
		Extra    interface{} `json:"extra"`
	} `json:"detail"`
	Plugin string `json:"plugin"`
	Target struct {
		Url string `json:"url"`
	} `json:"target"`
}

// 参考文档 https://docs.xray.cool/#/README
func XRayExec(ctx context.Context, target string, isCrawler int) ([]XRayItem, error) {
	// 获取配置文件
	if err := xrayGetConfig(); err != nil {
		return nil, err
	}

	// 获取执行命令
	var command []string
	var outFile string
	if isCrawler == 1 {
		command, outFile = xrayGetCrawlerCommand(target)
	} else {
		command, outFile = xrayGetBaseCommand(target)
	}

	// 执行命令 这里为阻塞读取
	fmt.Println(command)
	if err := xrayExec(ctx, command); err != nil {
		return nil, errors.New("xray cmd执行失败：" + err.Error())
	}

	// 获取执行结果
	content, err := file.ReadFile(outFile)
	if err != nil {
		return nil, errors.New("xray 获取执行结果失败：" + err.Error())
	}

	var res []XRayItem
	if err := json.Unmarshal(content, &res); err != nil {
		return nil, errors.New("xray 获取解析结果失败：" + err.Error() + "jsonData is " + string(content))
	}
	return res, nil
}

// 获取配置文件
func xrayGetConfig() error {
	configStruct := make(map[string]string)
	err := config.Load("xray", &configStruct)
	if err != nil {
		return errors.New("config.json获取xray配置失败：" + err.Error())
	}
	xrayExecFile = configStruct["execFile"]
	if xrayExecFile == "" {
		return errors.New("config.json获取xray配置中的可执行文件execFile不可为空")
	}
	if !file.CheckFileExist(xrayExecFile) {
		return errors.New("config.json获取xray配置中的可执行文件execFile不存在")
	}

	xrayExecConfigFile = configStruct["configFile"]
	if xrayExecConfigFile == "" {
		return errors.New("config.json获取xray配置中的可执行文件的配置文件configFile不可为空")
	}
	if !file.CheckFileExist(xrayExecConfigFile) {
		return errors.New("config.json获取xray配置中的可执行文件的配置文件configFile不存在")
	}

	outputDir = configStruct["outputDir"]
	if outputDir == "" {
		return errors.New("config.json获取xray配置中的结果输出文件夹outputDir不可为空")
	}
	if !file.CheckFileExist(outputDir) {
		return errors.New("config.json获取xray配置中结果输出文件夹outputDir不存在")
	}

	return nil
}

func xrayGetBaseCommand(target string) ([]string, string) {
	outputFile := path.Join(outputDir, uuid.New().String()+".json")
	execStep := []string{
		xrayExecFile,
		"--config", xrayExecConfigFile,
		"webscan",
		"--url", target,
		"--json-output", outputFile,
	}
	return execStep, outputFile
}

func xrayGetCrawlerCommand(target string) ([]string, string) {
	outputFile := path.Join(outputDir, uuid.New().String()+".json")
	execStep := []string{
		xrayExecFile,
		"--config", xrayExecConfigFile,
		"webscan",
		"--basic-crawler", target,
		"--json-output", outputFile,
	}
	return execStep, outputFile
}

func xrayExec(ctx context.Context, command []string) error {
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)

	var out, err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	err1 := cmd.Run()
	if err1 != nil {
		return errors.New(out.String())
	}
	if err.String() != "" {
		return errors.New(err.String())
	}
	return nil
}
