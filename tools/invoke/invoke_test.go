package invoke

//
//import (
//	"context"
//	"fmt"
//	"gitlabee.4dogs.cn/common/config"
//	"gitlabee.4dogs.cn/common/mq"
//	"gitlabee.4dogs.cn/common/mysql"
//	"gitlabee.4dogs.cn/common/neo4j"
//	"gitlabee.4dogs.cn/common/redis"
//	"testing"
//)
//
//func TestMain(m *testing.M) {
//	fmt.Println("初始化")
//	configFilePath := "../config.json"
//	err := config.NewConfig(configFilePath)
//	if err != nil {
//		fmt.Println("new config err:", err)
//		return
//	}
//	mysql.Setup()
//	redis.Setup()
//	neo4j.Setup()
//	mq.RabbitMQSetup()
//	m.Run()
//}
//
//// 测试调用web扫描工具
//func TestInvokeWebScanner(t *testing.T) {
//	ctx := context.Background()
//	toolParamList := []ToolParam{
//		ToolParam{ParamName: "root_url", ParamValue: "http://192.168.0.96"},
//	}
//	callInfo := CallInfo{
//		CallId:        "ddd",
//		ToolName:      "web_scanner",
//		ToolParamList: toolParamList,
//	}
//	callInfo.Invoke(ctx, callBackFunc)
//}
//
////脚本执行回调函数
//func callBackFunc(ctx context.Context, callInfo *CallInfo, result string) {
//	fmt.Println(callInfo.CallId)
//	fmt.Println(callInfo.ToolName)
//	fmt.Println(result)
//}
//
////测试msfapi授权功能
//func TestMsfAuth(t *testing.T) {
//	ctx := context.Background()
//	var callInfo CallInfo
//	token, err := callInfo.Auth(ctx)
//	if err != nil {
//		t.Logf("msf auth error: %s", err.Error())
//		t.FailNow()
//	}
//	if token == "" {
//		t.Logf("msf auth token error")
//	}
//}
//
//// 测试调用msf脚本
//func TestInvokeMsfScript(t *testing.T) {
//	ctx := context.Background()
//	toolParamList := []ToolParam{
//		ToolParam{ParamName: "ip", ParamValue: "192.168.88.99"},
//		ToolParam{ParamName: "reverseHost", ParamValue: "192.168.3.20"},
//		ToolParam{ParamName: "smb_port", ParamValue: "445"},
//	}
//	c := CallInfo{
//		CallId:        "ddd",
//		ToolName:      "ms17_010_eternalblue",
//		ToolParamList: toolParamList,
//	}
//	c.Invoke(ctx, callBackFunc)
//}
