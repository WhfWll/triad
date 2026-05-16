//go:build ignore

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	"os"
	"smart/services"
	"smart/tools/enums"
	"strconv"
	"time"
)

func main() {
	date := "2023-11-21"
	parsedTime, err := time.Parse(enums.TimeYMDBarLayout, date)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	days := int(parsedTime.Sub(time.Now()).Hours() / 24)
	serialNumber := generateSerialNumberLocal()
	authCode := generateAuthCodeLocal(serialNumber, strconv.Itoa(days))
	err = authLocal(ctx, serialNumber, authCode, date, strconv.Itoa(days))
	if err != nil {
		fmt.Println("授权出错", err)
		return
	}
	fmt.Println("授权完成，请在界面中系统设置中查看详细授权信息")
}

func generateAuthCodeLocal(serialNumber string, days string) string {
	_, err := strconv.Atoi(days)
	if err != nil {
		fmt.Println("请输入正确的天数")
		os.Exit(1)
	}
	tempData, err := buildAuthDataLocal(serialNumber, days)
	if err != nil {
		panic(err)
	}
	var auth services.Auth
	ciphertext := auth.RsaEncrypt(tempData, []byte(enums.SWPubKey))
	return hex.EncodeToString(ciphertext)
}

func generateSerialNumberLocal() string {
	ctx := context.Background()
	var auth services.Auth
	serialNumber, err := auth.GenerateSystemSerialNumber(ctx)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return serialNumber
}

func authLocal(ctx context.Context, serialNumber string, authCode, date, authDays string) error {
	err := config.NewConfig("/opt/laozhi/smart/config.json") //加载配置文件
	if err != nil {
		return err
	}
	mysql.Setup()
	var authService services.Auth
	err = authService.UpdateAuthInfo(ctx, map[string]string{
		"authCode":    authCode,
		"authTime":    time.Now().Format(enums.TimeYMDBarLayout),
		"productID":   serialNumber,
		"productName": "自动化渗透测试系统",
		"authExpTime": date,
		"authDays":    authDays + "天",
		"leftDays":    authDays,
	})
	if err != nil {
		return err
	}
	err = authService.UpdateAuthState(ctx, enums.ProductAuthStateSuccess)
	if err != nil {
		return err
	}
	return nil
}

func buildAuthDataLocal(serialNumber string, days string) ([]byte, error) {
	tempMap := make(map[string]string, 0)
	tempMap["productID"] = serialNumber
	tempMap["generateTime"] = time.Now().Format("2006-01-02")
	tempMap["authDays"] = days
	tempData, err := json.Marshal(tempMap)
	if err != nil {
		return nil, nil
	}
	return tempData, nil
}
