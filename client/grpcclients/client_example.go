package grpcclients

import (
	"context"
	"fmt"
	"log"
	"time"

	"smart/client/grpcclients/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ScannerClient gRPC客户端示例 - 简化版本
type ScannerClient struct {
	conn   *grpc.ClientConn
	client pb.ScannerServiceClient
}

// NewScannerClient 创建新的客户端
func NewScannerClient(serverAddr string) (*ScannerClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("连接gRPC服务器失败: %v", err)
	}

	client := pb.NewScannerServiceClient(conn)
	return &ScannerClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close 关闭连接
func (c *ScannerClient) Close() error {
	return c.conn.Close()
}

// HealthCheck 健康检查
func (c *ScannerClient) HealthCheck() (*pb.HealthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.HealthCheck(ctx, &pb.HealthRequest{})
}

// StartScan 启动扫描
func (c *ScannerClient) StartScan(taskID, target string, configJSON string) (*pb.ScanResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.ScanRequest{
		TaskId:     taskID,
		Target:     target,
		ConfigJson: configJSON,
	}

	return c.client.ExecuteScan(ctx, req)
}

// GetScanStatus 获取扫描状态
func (c *ScannerClient) GetScanStatus(taskID string) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.GetScanStatus(ctx, &pb.StatusRequest{
		TaskId: taskID,
	})
}

// GetScanResults 获取扫描结果
func (c *ScannerClient) GetScanResults(taskID string) (*pb.ResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.GetScanResults(ctx, &pb.ResultRequest{
		TaskId: taskID,
	})
}

// GetScanResultsIncremental 获取增量扫描结果
func (c *ScannerClient) GetScanResultsIncremental(taskID string) (*pb.ResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.GetScanResults(ctx, &pb.ResultRequest{
		TaskId:      taskID,
		Incremental: true,
	})
}

// GetScanResultsPaginated 获取分页扫描结果
func (c *ScannerClient) GetScanResultsPaginated(taskID string, offset, limit int32) (*pb.ResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.client.GetScanResults(ctx, &pb.ResultRequest{
		TaskId: taskID,
		Offset: offset,
		Limit:  limit,
	})
}

// 使用示例
func ExampleUsage() {
	// 连接到gRPC服务器
	client, err := NewScannerClient("localhost:50051")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 健康检查
	health, err := client.HealthCheck()
	if err != nil {
		log.Printf("健康检查失败: %v", err)
	} else {
		fmt.Printf("节点状态: %s, 存活: %v\n", health.Status, health.Alive)
	}

	// 启动扫描
	taskID := fmt.Sprintf("task-%d", time.Now().Unix())
	configYAML := `timeout: 300
concurrency: 10
scanMode: web`

	scanResp, err := client.StartScan(taskID, "http://example.com", configYAML)
	if err != nil {
		log.Printf("启动扫描失败: %v", err)
		return
	}
	fmt.Printf("扫描任务已启动: %s\n", scanResp.Message)

	// 轮询扫描状态
	for {
		status, err := client.GetScanStatus(taskID)
		if err != nil {
			log.Printf("获取状态失败: %v", err)
			break
		}

		fmt.Printf("任务状态: %s, 进度: %d%%, 消息: %s\n",
			status.Status, status.Progress, status.Message)

		if status.Status == "completed" || status.Status == "failed" || status.Status == "cancelled" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	// 获取扫描结果
	results, err := client.GetScanResults(taskID)
	if err != nil {
		log.Printf("获取结果失败: %v", err)
		return
	}

	fmt.Printf("扫描结果: %s\n", results.ResultsJson)
	fmt.Printf("总计: %d 个结果\n", results.TotalCount)
}
