package main

import (
	"fmt"
	"os"
	"path/filepath"
	"smart/crons"
	"smart/routers"
	"smart/services/reverse_shell"
	"smart/tools/enums"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/logger"
	"gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"
)

func main() {
	// 初始化系统路径
	if exe, err := os.Executable(); err == nil {
		// 尝试根据可执行文件位置推断
		// 假设标准部署结构: /path/to/project/smart/smart (binary)
		// 则 ProjectDir 为 /path/to/project/
		dir := filepath.Dir(exe)
		if filepath.Base(dir) == "smart" {
			enums.InitSystemPaths(filepath.Dir(dir))
		}
	}

	err := config.NewConfig()
	if err != nil {
		fmt.Println("new config err:", err)
		return
	}

	// 尝试从配置文件加载 project_dir
	var projectDir string
	if err := config.Load("project_dir", &projectDir); err == nil && projectDir != "" {
		enums.InitSystemPaths(projectDir)
	}

	err = logger.Setup()
	if err != nil {
		fmt.Println("setup logger err:", err)
		return
	}
	mysql.Setup()
	redis.Setup()
	crons.Start()

	// 启动反弹shell服务
	if err := reverse_shell.GetService().Start(6666); err != nil {
		fmt.Printf("Failed to start reverse shell service: %v\n", err)
	}

	gin.SetMode(config.GetGinMode())
	r := routers.RegisterRoute()
	err = r.Run(config.GetHost() + ":" + config.GetPort())
	fmt.Println(err)
}
