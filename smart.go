package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"smart/crons"
	"smart/routers"
	"smart/services"
	"smart/services/reverse_shell"
	"smart/tools/enums"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/logger"
	"gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"

	log "github.com/sirupsen/logrus"
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
	// Windows 开发：config 指向 /opt/laozhi/ 时，若 CVE 库在当前工程目录则回退到 cwd
	if wd, err := os.Getwd(); err == nil {
		cveInWd := filepath.Join(wd, "data", "default-cve.db")
		if _, err := os.Stat(cveInWd); err == nil {
			if _, err := os.Stat(filepath.Join(enums.SystemUpgradeProjectDir, "data", "default-cve.db")); os.IsNotExist(err) {
				enums.InitSystemPaths(wd + string(os.PathSeparator))
			}
		}
	}

	err = logger.Setup()
	if err != nil {
		fmt.Println("setup logger err:", err)
		return
	}
	mysql.Setup()

	services.InitCveDB()
	services.EnsureHostVulnSchema(context.Background())
	services.EnsureHostMalwareSchema(context.Background())
	// YARA 规则库较大，后台加载避免阻塞 HTTP 服务启动
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("malware YARA rules: background loading panicked: %v", r)
			}
		}()
		ctx := mysql.NewContext(context.Background(), mysql.GetDB())
		log.Info("malware YARA rules: loading from DB in background...")
		services.InitMalwareRulesFromDB(ctx)
	}()
	// 从数据库加载基线规则（不再从 JSON 文件加载）
	services.InitBaselineRulesFromDB(context.Background())
	services.EnsureBaselineCheckResultSchema(context.Background())
	services.EnsureDatasecRuleTable(context.Background())
	services.InitDatasecRulesFromDB(context.Background())
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
