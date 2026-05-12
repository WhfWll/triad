package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"smart/client/grpcclients"
	"smart/client/grpcclients/pb"
)

func main() {
	// 连接到gRPC服务器
	client, err := grpcclients.NewScannerClient("localhost:8089")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	fmt.Println("=== gRPC扫描节点测试 ===")

	// 2. 使用已完成的测试任务进行增量获取测试
	fmt.Println("\n2. 使用已完成的测试任务进行增量获取测试...")
	taskID := "test-task-002" // 使用已知的完成任务
	target := "http://192.168.3.3:8770"

	// 2. 启动扫描任务
	fmt.Println("\n2. 启动扫描任务...")
	config := map[string]interface{}{
		"scanUrl":     "http://192.168.3.3:8770/",
		"scanIp":      "192.168.3.3",
		"timeout":     300,
		"concurrency": 1,
		"debug":       true,
		"scanMode":    "all",
		"scripts":     []string{"mitm_sql_check", "mitm_xss_check", "mitm_gitinfo_check"},
		"portScan": map[string]string{
			"scanPort": "8770",
			"protocol": "tcp",
			"timeout":  "10",
		},
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		log.Printf("构造扫描配置失败: %v", err)
		return
	}

	scanResp, err := client.StartScan(taskID, target, string(configBytes))
	if err != nil {
		log.Printf("启动扫描失败: %v", err)
		return
	}
	fmt.Printf("扫描任务已启动: %s\n", scanResp.Message)
	fmt.Printf("任务ID: %s\n", scanResp.TaskId)

	fmt.Printf("测试任务ID: %s\n", taskID)
	fmt.Printf("目标: %s\n", target)

	// 3. 测试增量获取功能（模拟实时监控过程）
	fmt.Println("\n3. 模拟实时监控过程中的增量获取...")

	//var totalResultsReceived int32 = 0
	var totalLogsReceived int32 = 0

	// 模拟多次增量获取（第一次应该获取所有数据，后续应该为0）
	for i := 0; i < 30; i++ {
		fmt.Printf("\n[第%d次增量获取]\n", i+1)

		results, err := client.GetScanResultsIncremental(taskID)
		if err != nil {
			log.Printf("获取增量结果失败: %v", err)
			continue
		}

		newResultCount := len(results.DetectionResults)
		newLogCount := len(results.Logs)

		if newResultCount > 0 {
			fmt.Printf("  🔍 发现新检测结果: %d 个\n", newResultCount)
			displayIncrementalResults(results)
			//totalResultsReceived += int32(newResultCount)
		}

		if newLogCount > 0 {
			fmt.Printf("  📝 发现新日志: %d 条\n", newLogCount)
			displayIncrementalLogs(results.Logs)
			totalLogsReceived += int32(newLogCount)
		}

		if newResultCount == 0 && newLogCount == 0 {
			fmt.Printf("  ⏳ 暂无新数据\n")
		}

		time.Sleep(1 * time.Second)
	}

	// 4. 获取最终扫描结果
	//fmt.Println("\n4. 获取最终扫描结果...")
	//finalResults, err := client.GetScanResults(taskID)
	//
	//if err != nil {
	//	log.Printf("获取最终结果失败: %v", err)
	//} else {
	//	fmt.Printf("✅ 任务完成！总计发现: %d 个结果\n", finalResults.TotalCount)
	//	displayFinalStructuredResults(finalResults)
	//}
	//
	//// 5. 数据传输量对比测试
	//fmt.Println("\n5. 数据传输量对比测试...")
	//performDataTransferComparison(client, taskID, totalResultsReceived, totalLogsReceived)
	//
	//fmt.Println("=== 测试完成 ===")
}

// performDataTransferComparison 执行数据传输量对比测试
func performDataTransferComparison(client *grpcclients.ScannerClient, taskID string, totalResultsReceived, totalLogsReceived int32) {
	fmt.Printf("📊 数据传输量对比分析:\n")
	fmt.Printf("  通过增量模式累计获取: 漏洞结果 %d 个, 日志 %d 条\n", totalResultsReceived, totalLogsReceived)

	// 获取全量数据作为对比基准
	fmt.Printf("  正在获取全量数据作为对比基准...\n")
	fullResults, err := client.GetScanResults(taskID)
	if err != nil {
		log.Printf("获取全量结果失败: %v", err)
		return
	}

	totalResults := int32(len(fullResults.DetectionResults))
	totalLogs := int32(len(fullResults.Logs))

	fmt.Printf("  全量模式获取: 漏洞结果 %d 个, 日志 %d 条\n", totalResults, totalLogs)

	// 计算数据传输效率
	if totalResults > 0 || totalLogs > 0 {
		// 计算重复传输的数据量（如果使用全量模式多次获取）
		estimatedFullModeTransfers := int32(15) // 假设监控过程中进行了15次全量获取
		redundantResults := (estimatedFullModeTransfers - 1) * totalResults
		redundantLogs := (estimatedFullModeTransfers - 1) * totalLogs

		fmt.Printf("\n📈 效率分析（假设监控期间进行 %d 次数据获取）:\n", estimatedFullModeTransfers)
		fmt.Printf("  增量模式传输: 漏洞结果 %d 个, 日志 %d 条\n", totalResultsReceived, totalLogsReceived)
		fmt.Printf("  全量模式传输: 漏洞结果 %d 个, 日志 %d 条\n", estimatedFullModeTransfers*totalResults, estimatedFullModeTransfers*totalLogs)

		if redundantResults > 0 {
			resultSavings := float64(redundantResults) / float64(estimatedFullModeTransfers*totalResults) * 100
			fmt.Printf("  漏洞结果传输节省: %.1f%%\n", resultSavings)
		}

		if redundantLogs > 0 {
			logSavings := float64(redundantLogs) / float64(estimatedFullModeTransfers*totalLogs) * 100
			fmt.Printf("  日志传输节省: %.1f%%\n", logSavings)
		}

		// 测试再次增量获取（应该返回0条新数据）
		fmt.Printf("\n🔄 验证增量获取一致性...\n")
		incrementalResults, err := client.GetScanResultsIncremental(taskID)
		if err != nil {
			log.Printf("增量获取验证失败: %v", err)
		} else {
			newResults := len(incrementalResults.DetectionResults)
			newLogs := len(incrementalResults.Logs)
			fmt.Printf("  再次增量获取: 漏洞结果 %d 个, 日志 %d 条\n", newResults, newLogs)
			if newResults == 0 && newLogs == 0 {
				fmt.Printf("  ✅ 增量获取机制工作正常，无重复数据传输\n")
			} else {
				fmt.Printf("  ⚠️  检测到额外数据，可能有新的扫描活动\n")
			}
		}
	} else {
		fmt.Printf("  📝 本次扫描未产生数据，无法进行传输量对比\n")
	}
}

// displayIncrementalLogs 显示增量获取的日志
func displayIncrementalLogs(logs []string) {
	fmt.Printf("    📝 新增日志条目:\n")
	for i, line := range logs {
		isPortScan := strings.Contains(line, "port_scan") ||
			strings.Contains(line, "PortScan") ||
			strings.Contains(line, "服务扫描")

		prefix := "        "
		if isPortScan {
			prefix = "    [端口扫描] "
		}
		fmt.Printf("%s[%d] %s\n", prefix, i+1, line)
	}
	fmt.Println()
}

// displayIncrementalResults 显示增量结果
func displayIncrementalResults(results *pb.ResultResponse) {
	fmt.Printf("📊 增量检测结果 (%d个):\n", len(results.DetectionResults))

	for i, result := range results.DetectionResults {
		isPortScan := strings.Contains(strings.ToLower(result.ScriptName), "port") ||
			strings.Contains(result.Details, "port_scan")

		tag := "漏洞"
		if isPortScan {
			tag = "端口扫描"
		}

		fmt.Printf("    🔍 [%d] %s: %s\n", i+1, tag, result.ScriptName)
		fmt.Printf("        URL: %s\n", result.Url)
		if result.Parameter != "" {
			fmt.Printf("        参数: %s\n", result.Parameter)
		}
		if result.Payload != "" {
			fmt.Printf("        载荷: %s\n", result.Payload)
		}
		if result.Details != "" {
			fmt.Printf("        详情: %s\n", result.Details)
		}
		fmt.Println()
	}
}

// displayNewStructuredResults 显示新发现的结构化结果（保留用于兼容性）
func displayNewStructuredResults(results *pb.ResultResponse, lastCount int32) {
	// 只显示新增的检测结果
	fmt.Printf("📊 新增检测结果:\n")

	if int(lastCount) < len(results.DetectionResults) {
		newResults := results.DetectionResults[lastCount:]
		for i, result := range newResults {
			fmt.Printf("    🔍 [%d] 漏洞: %s\n", int(lastCount)+i+1, result.ScriptName)
			fmt.Printf("        URL: %s\n", result.Url)
			if result.Parameter != "" {
				fmt.Printf("        参数: %s\n", result.Parameter)
			}
			if result.Payload != "" {
				fmt.Printf("        载荷: %s\n", result.Payload)
			}
			if result.Details != "" {
				fmt.Printf("        详情: %s\n", result.Details)
			}
			fmt.Println()
		}
	}
}

// displayFinalStructuredResults 显示最终完整的结构化结果
func displayFinalStructuredResults(results *pb.ResultResponse) {
	fmt.Printf("📋 扫描详细信息:\n")
	if results.TaskInfo != nil {
		fmt.Printf("  任务ID: %s\n", results.TaskInfo.TaskId)
		fmt.Printf("  目标: %s\n", results.TaskInfo.Target)
		fmt.Printf("  开始时间: %s\n", results.TaskInfo.StartTime)
		if results.TaskInfo.EndTime != "" {
			fmt.Printf("  结束时间: %s\n", results.TaskInfo.EndTime)
		}
		if results.TaskInfo.WebhookUrl != "" {
			fmt.Printf("  Webhook地址: %s\n", results.TaskInfo.WebhookUrl)
		}
	}

	if len(results.DetectionResults) > 0 {
		fmt.Printf("\n📊 发现的漏洞 (%d个):\n", len(results.DetectionResults))
		for i, result := range results.DetectionResults {
			fmt.Printf("  [%d] %s\n", i+1, result.ScriptName)
			fmt.Printf("      URL: %s\n", result.Url)
			if result.Parameter != "" {
				fmt.Printf("      参数: %s\n", result.Parameter)
			}
			if result.Payload != "" {
				fmt.Printf("      载荷: %s\n", result.Payload)
			}
			if result.PayloadSuccessFlag != "" {
				fmt.Printf("      成功标识: %s\n", result.PayloadSuccessFlag)
			}
			if result.Details != "" {
				fmt.Printf("      详情: %s\n", result.Details)
			}
			fmt.Println()
		}
	} else {
		fmt.Printf("\n✅ 未发现安全漏洞\n")
	}

	if len(results.Logs) > 0 {
		fmt.Printf("\n📝 扫描日志 (%d条):\n", len(results.Logs))
		for i, log := range results.Logs {
			fmt.Printf("  [%d] %s\n", i+1, log)
		}
	}
}
