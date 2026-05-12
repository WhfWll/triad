package main

import (
	"context"
	"fmt"
	"os"
	"smart/services"

	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"
)

func main() {
	// 尝试加载配置文件，默认假设在当前目录下或上级目录
	configFile := "config.json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 尝试上级目录
		configFile = "../config.json"
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			// 尝试上上级目录
			configFile = "../../config.json"
		}
	}

	fmt.Printf("Using config file: %s\n", configFile)

	// 注意：config.NewConfig 在 smart.go 中没有参数，但可能支持可变参数。
	// 这里假设如果不传参数默认读取 config.json。如果传入参数则读取指定文件。
	// 根据 cmd/commands.go 的用法，传入路径是可行的。
	if err := config.NewConfig(configFile); err != nil {
		fmt.Printf("Load config failed: %v\n", err)
		// 如果加载失败，尝试不传参数（默认行为）
		fmt.Println("Retrying with default config location...")
		if err := config.NewConfig(); err != nil {
			fmt.Printf("Load default config failed: %v\n", err)
			return
		}
	}

	// 初始化数据库
	mysql.Setup()
	redis.Setup()

	ctx := context.Background()
	authSrv := &services.Auth{}

	// 1. 获取 Linux 系统序列号 (通过 dmidecode)
	fmt.Println("--------------------------------------------------")
	fmt.Println("Attempting to generate system serial number (requires dmidecode)...")
	sysSerial, err := authSrv.GenerateSystemSerialNumber(ctx)
	if err != nil {
		fmt.Printf("Error generating system serial number: %v\n", err)
	} else {
		fmt.Printf("System Serial Number (Calculated): %s\n", sysSerial)
	}

	// 2. 获取数据库中的授权信息
	fmt.Println("--------------------------------------------------")
	fmt.Println("Attempting to retrieve auth info from database...")
	authInfo, err := authSrv.GetAuthInfo(ctx)
	if err != nil {
		fmt.Printf("Error retrieving auth info from DB: %v\n", err)
	} else {
		fmt.Printf("Auth Info from DB:\n")
		if val, ok := authInfo["productID"]; ok {
			fmt.Printf("  Product ID (Serial in DB): %s\n", val)
		} else {
			fmt.Println("  Product ID not found in auth info")
		}
		
		if val, ok := authInfo["authCode"]; ok {
			fmt.Printf("  Auth Code: %s\n", val)
		}
		
		if val, ok := authInfo["authTime"]; ok {
			fmt.Printf("  Auth Time: %s\n", val)
		}
	}
	fmt.Println("--------------------------------------------------")
}
