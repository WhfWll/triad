//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	dbutils "smart/tools/data"
	"time"

	"github.com/urfave/cli"
	mysqlsEnv "gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"
)

var ServiceCommand = cli.Command{
	Name:  "service",
	Usage: "进行系统服务信息监控",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "all"},
		&cli.BoolFlag{Name: "mysql"},
		&cli.BoolFlag{Name: "redis"},
		//&cli.BoolFlag{Name: "nginx"},
		//&cli.BoolFlag{Name: "go-rod"},
		//&cli.BoolFlag{Name: "msf"},
		//&cli.BoolFlag{Name: "smart"},
		//&cli.BoolFlag{Name: "decision"},
	},
	Action: func(c *cli.Context) error {
		ctx := context.Background()
		all := c.Bool("all")
		mysql := c.Bool("mysql")
		if all {
			checkAll(ctx)
		}
		if mysql {
			err := dbutils.CheckMysql(ctx)
			if err != nil {
				fmt.Println("运行环境异常: ", err.Error())
			}
		}
		return nil
	},
}

var DataSubCommand = cli.Command{
	Name:  "data",
	Usage: "进行系统运行数据信息监控",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "task"},
	},
	Action: func(c *cli.Context) error {
		mysqlsEnv.Setup()
		redis.Setup()

		ctx := context.Background()
		task := c.Bool("task")
		if task {
			fmt.Println(ctx)
			fmt.Println("monitor task data")
		}

		return nil
	},
}

func main() {
	dbutils.InitConfigInfo()
	app := cli.NewApp()
	app.Usage = "smart monitor tools"
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, ServiceCommand)
	app.Commands = append(app.Commands, DataSubCommand)
	err := app.Run(os.Args)
	if err != nil {
		println(err.Error())
		return
	}
}

func checkAll(ctx context.Context) {
	for {
		err := dbutils.CheckMysql(ctx)
		if err != nil {
			fmt.Println("运行环境异常: ", err.Error())
			dbutils.AddLog(ctx, "运行环境异常: "+err.Error())
		}
		err = dbutils.CheckRedis(ctx)
		if err != nil {
			fmt.Println("运行环境异常: ", err.Error())
			mysqlsEnv.Setup()
			dbutils.AddLog(ctx, "运行环境异常: "+err.Error())
		}
		time.Sleep(20 * time.Second)
	}
}
