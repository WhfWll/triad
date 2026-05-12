package services

import (
	"context"
	"fmt"
	"testing"

	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"
)

func TestMain(m *testing.M) {
	fmt.Println("初始化")
	configFilePath := "../config.json"
	err := config.NewConfig(configFilePath)
	if err != nil {
		fmt.Println("new config err:", err)
		return
	}
	mysql.Setup()
	redis.Setup()
	m.Run()
}

func TestGetAuthInfo(t *testing.T) {
	a := Auth{}
	ctx := context.Background()
	authInfo, err := a.GetAuthInfo(ctx)
	if err != nil {
		t.Logf("get authInfo err: %s", err)
		t.FailNow()
	}
	if authInfo["productID"] == "" {
		t.Logf("get authInfo err: %s", err)
		t.FailNow()
	}
}

//func TestGenerateSystemSerialNumber(t *testing.T) {
//	a := Auth{}
//	ctx := context.Background()
//	serialNumber, err := a.GenerateSystemSerialNumber(ctx)
//	if err != nil {
//		t.Logf("get system serial number err: %s", err)
//		t.FailNow()
//	}
//	fmt.Println(serialNumber)
//}
