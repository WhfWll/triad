package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"smart/tools/file"
	"strconv"
	"strings"
)

// 分布式yak引擎节点在线安装

// makeself可执行文件地址
const makeSelfDir = "/home/dogs/build/makeself-2.5.0/makeself.sh"

// 服务器制作目录
const makeDir = "/opt/upgrade/yak_grpc_node/build"

// 服务器制作完成存放的目录
const makeBinDir = "/opt/upgrade/yak_grpc_node/bin"
const makeDoneDir = makeBinDir + "/yakGrpcNodeInstall.bin"

// awvs目录
const awvsDir = "/opt/laozhi/decision/scanner"

// finger指纹地址
const fingerFileName = "/opt/laozhi/decision/finger/finger_rules.json"

// 安装文件地址
const installFileName = makeDir + "/install.sh"

// yaklang 下载地址
const downloadYaklangUrl = "https://yaklang.oss-accelerate.aliyuncs.com/yak/latest/yak_linux_amd64"

// 使用makeself打包后，部署安装时执行的bash文件
const installSh = `#!/bin/bash

# 指纹 - 源文件路径
fingerSourceFile="./finger_rules.json"
# 指纹 - 目标目录
fingerTargetFolder="/opt/laozhi/decision/finger"
# 创建目录 - 如果不存在的话
mkdir -p "$fingerTargetFolder"
# 删除服务上的文件
echo "删除服务上的finger文件$fingerTargetFolder/finger_rules.json"
rm -rf "$fingerTargetFolder/finger_rules.json"
# 移动文件到目标文件夹
mv "$fingerSourceFile" "$fingerTargetFolder"
echo "安装指纹文件成功"

# webscanner - 源目录
webscannerSourceDir="./scanner"
# webscanner - 目标目录
webscannerTargetDir="/opt/laozhi/decision"
# 删除服务上的文件
echo "删除服务上webscanner的文件夹$webscannerTargetDir/scanner"
rm -rf "$webscannerTargetDir/scanner"
# 创建目录 - 如果不存在的话
mkdir -p "$webscannerTargetDir"
# 移动文件到目标文件夹
mv "$webscannerSourceDir" "$webscannerTargetDir"
echo "安装scanner插件成功"

# yak - 源目录
yakSourceDir="./yak"
# webscanner - 目标目录
yakTargetDir="/opt/laozhi/yak"
# 删除服务上的文件
echo "删除服务上的yak文件夹$yakTargetDir"
rm -rf "$yakTargetDir"
# 创建目录 - 如果不存在的话
mkdir -p "$yakTargetDir"
chmod -R 777 "$yakTargetDir"
# 移动文件到目标文件夹
mv "$yakSourceDir" "$yakTargetDir"
echo "安装分布式引擎成功"

# yak - 引擎启动
yakRunCmd="sudo nohup $yakTargetDir/yak grpc --host 0.0.0.0 --port 9919 --project-db $yakTargetDir/project.db --profile-db $yakTargetDir/profile.db > /dev/null 2>&1 &"
echo "准备启动分布式引擎"
eval "$yakRunCmd"
echo "引擎启动完成"`

func main() {
	// 校验对应文件是否存在，是否可以制作
	if err := checkFile(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// 先清空原制作目录
	fmt.Println("清空制作目录重新制作 " + makeDir)
	if err := os.RemoveAll(makeDir + "/*"); err != nil {
		fmt.Println(err.Error())
		return
	}

	// 移动文件到制作目录中
	// awvs目录
	fmt.Println("移动awvs目录到制作目录中 cp -r " + awvsDir + " " + makeDir)
	if err := copyDirectory(awvsDir, makeDir+"/scanner", true); err != nil {
		fmt.Println(err.Error())
		return
	}
	// finger指纹地址
	fmt.Println("移动finger指纹文件到制作目录中 cp " + fingerFileName + " " + makeDir + "/finger_rules.json")
	if err := copyFile(fingerFileName, makeDir+"/finger_rules.json"); err != nil {
		fmt.Println(err.Error())
		return
	}
	// 使用makeself打包后，部署安装时执行的bash文件
	fmt.Println("使用makeself打包后，部署安装时执行的bash文件写入操作 write to " + installFileName)
	if err := writeFile(installFileName, installSh); err != nil {
		fmt.Println(err.Error())
		return
	}
	// 下载最新yak 到安装包内
	fmt.Println("下载最新yak 到安装包内 write to " + makeDir + "/yak")
	if err := downloadYaklang(makeDir + "/yak"); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("开始使用makeself打包")
	// 执行makeself进行打包 sudo makeself.sh --target /tmp/94eb5a1cbed058a9 /home/dogs/20230915 test.bin 描述信息 ./update.sh
	makeCmd := make([]string, 0)
	// 指定打包后的文件，在目标服务器上执行时，解压至哪个目录下
	makeCmd = append(makeCmd, "--target")
	makeCmd = append(makeCmd, "/tmp/94eb5a1cbed058a9")
	// 指定当前要打包哪个目录
	makeCmd = append(makeCmd, makeDir)
	// 指定打包后可执行文件的名称
	makeCmd = append(makeCmd, makeDoneDir)
	// 解压时允许输出的内容
	makeCmd = append(makeCmd, "解压安装包")
	// 指定在目标服务器上，解压后执行的安装程序
	makeCmd = append(makeCmd, "./install.sh")
	// 执行打包
	err := exec.Command(makeSelfDir, makeCmd...).Run()
	if err != nil {
		fmt.Println("makeself打包执行" + makeSelfDir + " " + strings.Join(makeCmd, " ") + " err:" + err.Error())
		return
	}

	fmt.Println("打包完成，可执行文件地址：" + makeDoneDir)
	fmt.Println("done")
}

// 校验对应文件是否存在，是否可以制作
func checkFile() error {
	// 校验makeself是否存在
	fmt.Println("校验makeself是否存在")
	if !file.CheckPathExist(makeSelfDir) {
		return errors.New("制作脚本" + makeSelfDir + "不存在，请确认服务器上是否有此文件")
	}
	// 校验服务器制作目录是否存在
	fmt.Println("校验服务器制作目录是否存在")
	if !file.CheckPathExist(makeDir) {
		return errors.New("制作目录" + makeDir + "不存在，请在服务器上创建此目录")
	}
	// 校验awvs目录是否存在
	fmt.Println("校验awvs目录是否存在")
	if !file.CheckPathExist(awvsDir) {
		return errors.New("awvs目录" + awvsDir + "不存在，请确认服务器上是否有此目录")
	}
	// 校验finger指纹文件是否存在
	fmt.Println("校验finger指纹文件是否存在")
	if !file.CheckPathExist(fingerFileName) {
		return errors.New("finger指纹地址" + fingerFileName + "不存在，请确认服务器上是否有此文件")
	}

	if !file.CheckPathExist(makeBinDir) {
		err := os.MkdirAll(makeBinDir, os.ModePerm)
		if err != nil {
			return errors.New("服务器" + makeBinDir + "不存在，尝试创建失败：" + err.Error())
		}
	}
	return nil
}

func copyDirectory(src, dst string, isDir bool) error {
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	defer os.Chmod(dst, 0755) // Set destination directory permissions to match source directory permissions

	// 如果目标是目录且目录不存在，那就先创建
	if isDir && !file.CheckPathExist(dst) {
		if err = os.MkdirAll(dst, os.ModePerm); err != nil {
			return err
		}
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		if file.IsDir() {
			err = copyDirectory(srcPath, dstPath, true)
			if err != nil {
				return errors.New("写文件失败:" + err.Error())
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return errors.New("写文件失败:" + err.Error())
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	fmt.Println("File copied successfully:", dst)
	return nil
}

func writeFile(dst string, content string) error {
	fd, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return errors.New("打开安装文件失败:" + err.Error())
	}
	defer fd.Close()

	_, err = fd.WriteString(content)
	if err != nil {
		return errors.New("写如安装文件失败:" + err.Error())
	}
	return nil
}

func downloadYaklang(dst string) error {
	response, err := http.Get(downloadYaklangUrl)
	if err != nil {
		return errors.New("下载yaklang文件时出现错误:" + err.Error())
	}
	defer response.Body.Close()

	// 检查响应状态码
	if response.StatusCode != http.StatusOK {
		return errors.New("下载yaklang文件时出现错误:" + strconv.Itoa(response.StatusCode))
	}

	// 创建本地文件以保存下载的内容
	yakfile, err := os.Create(dst)
	if err != nil {
		return errors.New("创建yaklang文件时出现错误:" + err.Error())
	}
	defer yakfile.Close()

	// 将响应体中的内容复制到本地文件
	_, err = io.Copy(yakfile, response.Body)
	if err != nil {
		return errors.New("复制yaklang文件内容时出现错误:" + err.Error())
	}

	if err := os.Chmod(dst, os.ModePerm); err != nil {
		return errors.New("下载yaklang文件后给权限出现错误:" + err.Error())
	}
	return nil
}
