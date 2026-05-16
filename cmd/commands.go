//go:build ignore

package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func main() {
	/*
		//脚本加载配置的方法
		//方法一：手动指定配置文件位置
		err := config.NewConfig("../config.json") //加载配置文件
		if err != nil {
			return
		}
		mysql.Setup()
		redis.Setup()
		//方法二：单独通过传参的方式加载组件
		logger.SetupByParam("../../log", "log", 7)                             //初始化日志
		mysql.SetupByParam("root:123456@tcp(192.168.0.190:3306)/smartdefense") //初始化mysql
		redis.SetupByParam(":CM$ecThreat%232015@192.168.0.187:6379/0")         //初始化redis

	*/
	app := cli.NewApp()
	app.Usage = "smart manage tools"
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, VersionSubCommand)
	app.Commands = append(app.Commands, InstallSubCommand)
	app.Commands = append(app.Commands, YakSubCommand)
	err := app.Run(os.Args)
	if err != nil {
		println(err.Error())
		return
	}
}

var InstallSubCommand = cli.Command{
	Name:  "install",
	Usage: "安装 Yak/Install Yak  (Add to ENV PATH)",
	Action: func(c *cli.Context) error {
		return nil
	},
}

var VersionSubCommand = cli.Command{
	Name:  "version",
	Usage: "smart version",
	Action: func(c *cli.Context) error {
		version := "V0.1.0.230628β"
		fmt.Println("Version: " + version)
		return nil
	},
}

var YakSubCommand = cli.Command{
	Name:  "yak",
	Usage: "执行yak调用",
	Action: func(c *cli.Context) error {
		return nil
	},
}
