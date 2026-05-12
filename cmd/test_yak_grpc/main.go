package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. 连接到远程 Yak gRPC 服务
	targetAddr := "192.168.22.142:9919"
	fmt.Printf("Connecting to Yak gRPC at %s...\n", targetAddr)

	conn, err := grpc.Dial(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := ypb.NewYakClient(conn)

	// 2. 准备测试脚本
	// 这是一个模拟 poc-nginxwebui-rce.yak 的简化版，用于测试参数传递
	testScript := `
println("Yak Script Started")

// 尝试从 CLI 参数获取 (这是最常见的方式)
target = cli.String("target")
port = cli.String("port")

println("Retrieved target from cli.String: " + target)
println("Retrieved port from cli.String: " + port)

// 模拟网络请求 (如果 target 是空的，这里会失败)
if target == "" {
    println("Error: Target is empty! Params injection might be failing.")
} else {
    println("Target is valid, attempting to ping...")
    // 简单打印验证通过
    println("Connection and Parameter Injection Successful!")
}
`

	// 3. 构造请求，注入参数
	// 这里设置最终攻击目标
	targetHost := "192.168.22.142"
	targetPort := "80"

	execReq := &ypb.ExecRequest{
		Script: testScript,
		Params: []*ypb.ExecParamItem{
			{Key: "target", Value: targetHost},
			{Key: "port", Value: targetPort},
		},
	}

	fmt.Printf("Executing script with target=%s, port=%s...\n", targetHost, targetPort)

	// 4. 执行并获取结果
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Exec(ctx, execReq)
	if err != nil {
		log.Fatalf("Exec failed: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			// EOF or error
			fmt.Printf("Stream ended: %v\n", err)
			break
		}

		if len(resp.Raw) > 0 {
			fmt.Printf("[Output] %s", string(resp.Raw))
		}
		if resp.IsMessage {
			fmt.Printf("[Message] %s\n", string(resp.Message))
		}
	}
}
