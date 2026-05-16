//go:build ignore

package main

import (
	"context"
	"flag"
	"fmt"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"
	"time"
)

func main() {
	var (
		list   string
		obj    string
		handle string
	)
	flag.StringVar(&list, "list", "", "")
	flag.StringVar(&obj, "obj", "task", "")
	flag.StringVar(&handle, "handle", "", "")
	flag.Parse()

	ctx := context.Background()
	err := config.NewConfig("../config.json") //加载配置文件
	if err != nil {
		return
	}
	mysql.Setup()

	if list == "running" {
		listRunning(ctx, obj)
	} else if list == "all" {
		fmt.Println(list)
		listAll(ctx, obj)
	} else if list == "timeout" {
		listTimeout(ctx, obj)
	} else if handle == "handle" {

	} else {
		return
	}
}

func listRunning(ctx context.Context, obj string) {
	if obj == "task" {
		fmt.Println("列出正在运行的任务")
		var srv services.TaskTask
		for _, task := range srv.GetTaskByStatus(ctx, enums.TaskStatusRunning) {
			fmt.Println("id:", task.ID, "taskName:", task.TaskName, "createTime:", task.CreateTime.Format(enums.TimeLayout), "updateTime:", task.UpdateTime.Format(enums.TimeLayout))
		}
	} else if obj == "target" {
		fmt.Println("列出正在运行的目标")
		var srv services.TaskTarget
		for _, target := range srv.GetTargetByStatus(ctx, enums.TargetStatusRunning) {
			fmt.Println("id:", target.ID, "targetUrl:", target.TargetURL, "createTime:", target.CreateTime.Format(enums.TimeLayout), "updateTime:", target.UpdateTime.Format(enums.TimeLayout))
		}
	}
}

func listAll(ctx context.Context, obj string) {
	if obj == "task" {
	}
}

func listTimeout(ctx context.Context, obj string) {
	var db = mysql.FromContext(ctx).Model(&mysqls.TaskTarget{})
	overTime := time.Now().Add(-1 * time.Hour).Format(utils.DateTime)
	var targetList []mysqls.TaskTarget
	err := db.Debug().Where("status = ? and update_time < ?", enums.TargetStatusRunning, overTime).Find(&targetList).Error
	if err != nil {
		fmt.Println(err)
	}
	for _, target := range targetList {
		fmt.Println("id:", target.ID, "targetUrl:", target.TargetURL, "createTime:", target.CreateTime.Format(enums.TimeLayout), "updateTime:", target.UpdateTime.Format(enums.TimeLayout))
	}
}

func handleTimeoutTasks(obj string) {
}

func handleTimeoutTargets(obj string) {
}
