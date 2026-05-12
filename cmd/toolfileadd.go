package main

import (
	"context"
	"flag"
	"fmt"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/file"
	"gitlabee.4dogs.cn/common/mysql"
	"os"
	"path/filepath"
	"smart/services"
	"strings"
)

/*
工具管理-辅助工具-工具库-导入脚本
注：
（1）导入的文件目录或文件，在服务器上不能删除，否则无法提供下载功能；
（2）可重复导入，相同路径的数据不会导入成功；
（3）可导入目录，也可导入单个文件；
*/
func main() {
	var filePath string
	flag.StringVar(&filePath, "filepath", "", "文件夹绝对路径") //扫描文件夹的绝对路径
	flag.Parse()
	if len(filePath) == 0 {
		fmt.Println("filepath参数-文件夹绝对路径不能为空...")
		return
	}
	if strings.HasPrefix(filePath, ".") {
		fmt.Println("filepath参数-必须是绝对路径")
		return
	}
	if !file.CheckFileExist(filePath) {
		fmt.Println("扫描的文件目录不存在...")
		return
	}
	err := config.NewConfig("../config.json") //加载配置文件
	if err != nil {
		return
	}
	mysql.Setup()
	fmt.Println("开始导入...")
	importToolsFile(filePath)
	fmt.Println("导入完成...")
}

func importToolsFile(filePathString string) {
	err := filepath.Walk(filePathString, handleFileInfo)
	if err != nil {
		fmt.Println("遍历目录出错:", err)
	}
}

func handleFileInfo(filePath string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("访问路径时出错:", err)
		return nil
	}
	if strings.TrimPrefix(filepath.Ext(info.Name()), ".") == "zip" {
		ctx := context.Background()
		var (
			toolSrv  services.ToolsManage
			fileType = "py"
		)
		fileNameArray := strings.Split(info.Name(), "_")
		if len(fileNameArray) > 1 {
			tmp := strings.Split(fileNameArray[len(fileNameArray)-1], ".")
			if len(tmp) > 0 {
				fileType = tmp[0]
			}
		}
		err = toolSrv.AddToolFile(ctx, info.Name(), fileType, filePath)
		if err != nil {
			fmt.Println("插入数据库出错:", err)
		}
	} else {
		fmt.Println("只能导入.zip后缀的文件...")
	}
	return nil
}
