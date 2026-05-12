package services

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"smart/api/typespec"
	"smart/models/mysqls"
	"time"

	"gitlabee.4dogs.cn/common/redis"
)

type ToolsManage struct {
}

// ToolsPingCreateReq 辅助工具 - ping - 创建
func (a *ToolsManage) ToolsPingCreate(ctx context.Context, redisClient *redis.RdsClient, ip, token string) {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
	// 构建ping命令
	cmd := exec.CommandContext(ctx, "ping", ip)

	// 获取命令的标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("获取标准输出管道时出错:", err)
		return
	}

	// 开始执行命令
	fmt.Println("ping - 创建任务 - 开始执行命令")
	if err := cmd.Start(); err != nil {
		fmt.Println("执行命令时出错:", err)
		return
	}

	// 逐行读取输出
	scanner := bufio.NewScanner(stdout)
	fmt.Println("ping - 创建任务 - 逐行读取输出")
	for scanner.Scan() {
		closeSign, _ := redisClient.Get(ctx, token+"_stop").Bool()
		if closeSign {
			cancel()
			fmt.Println("ping - 收到停止信号 - 调用cancel完成")
			continue
		}
		line := scanner.Text()
		redisClient.LPush(ctx, token+"_result", line)
		redisClient.Expire(ctx, token+"_result", 600*time.Second)
	}

	fmt.Println("ping - 收到停止信号 - 调用cancel完成 - cmd 协程退出")
	return
}

// ToolsPingResultUnread 辅助工具 - ping - 列表
func (a *ToolsManage) ToolsPingResultUnread(ctx context.Context, redisClient *redis.RdsClient, token string) []string {
	res := make([]string, 0)
	for true {
		line, _ := redisClient.RPop(ctx, token+"_result").Result()
		if line == "" {
			return res
		}
		res = append(res, line)
	}

	return nil
}

// ToolsPingStop 辅助工具 - ping - 停止
func (a *ToolsManage) ToolsPingStop(ctx context.Context, redisClient *redis.RdsClient, token string) error {
	fmt.Println("ping - 发送停止信号")
	redisClient.Set(ctx, token+"_stop", true, 10*time.Second)
	return nil
}

// ToolsTracerouteCreate 辅助工具 - traceroute - 创建
func (a *ToolsManage) ToolsTracerouteCreate(ctx context.Context, redisClient *redis.RdsClient, ip, token string) {
	ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
	// 构建traceroute命令
	cmd := exec.CommandContext(ctx, "traceroute", ip)

	// 获取命令的标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("获取标准输出管道时出错:", err)
		return
	}

	// 开始执行命令
	fmt.Println("traceroute - 创建任务 - 开始执行命令")
	if err := cmd.Start(); err != nil {
		fmt.Println("执行命令时出错:", err)
		return
	}

	// 逐行读取输出
	scanner := bufio.NewScanner(stdout)
	fmt.Println("traceroute - 创建任务 - 逐行读取输出")
	for scanner.Scan() {
		closeSign, _ := redisClient.Get(ctx, token+"_stop").Bool()
		if closeSign {
			cancel()
			fmt.Println("traceroute - 收到停止信号 - 调用cancel完成")
			continue
		}
		line := scanner.Text()
		redisClient.LPush(ctx, token+"_result", line)
		redisClient.Expire(ctx, token+"_result", 600*time.Second)
	}

	fmt.Println("traceroute - 收到停止信号 - 调用cancel完成 - cmd 协程退出")
	return
}

// ToolsTracerouteResultUnread 辅助工具 - traceroute - 列表
func (a *ToolsManage) ToolsTracerouteResultUnread(ctx context.Context, redisClient *redis.RdsClient, token string) []string {
	res := make([]string, 0)
	for true {
		line, _ := redisClient.RPop(ctx, token+"_result").Result()
		if line == "" {
			return res
		}
		res = append(res, line)
	}

	return nil
}

// ToolsTracerouteStop 辅助工具 - traceroute - 停止
func (a *ToolsManage) ToolsTracerouteStop(ctx context.Context, redisClient *redis.RdsClient, token string) error {
	fmt.Println("traceroute - 发送停止信号")
	redisClient.Set(ctx, token+"_stop", true, 10*time.Second)
	return nil
}

func (a *ToolsManage) ToolFileList(ctx context.Context, page, size int) ([]mysqls.Toolfile, int64, error) {
	var toolFileMysqls mysqls.Toolfile
	return toolFileMysqls.GetToolfileList(ctx, page, size)
}

func (a *ToolsManage) GetToolFileByPath(ctx context.Context, filePath string) (mysqls.Toolfile, error) {
	var toolFileMysqls mysqls.Toolfile
	return toolFileMysqls.GetToolfileByPath(ctx, filePath)
}

func (a *ToolsManage) AddToolFile(ctx context.Context, fileName, fileType, filePath string) error {
	var toolFileMysqls = mysqls.Toolfile{
		Name:       fileName,
		FileType:   fileType,
		FilePath:   filePath,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	return toolFileMysqls.AddToolfile(ctx)
}

// ToolsFingerList 获取指纹列表
func (a *ToolsManage) ToolsFingerList(ctx context.Context, page, size int, fingerName string, level int, class int) ([]mysqls.Finger, int64, error) {
	var fingerModel mysqls.Finger
	return fingerModel.GetFingerList(ctx, page, size, fingerName, level, class)
}

//
//// GetFingerClassEnumMap 获取指纹分类枚举映射
//func (a *ToolsManage) GetFingerClassEnumMap() map[int]string {
//	return map[int]string{
//		1: "Web应用",
//		2: "操作系统",
//		3: "硬件设备",
//		4: "网络设备",
//		5: "安全设备",
//		6: "数据库",
//		7: "中间件",
//		8: "开发框架",
//		9: "其他",
//	}
//}
//
//// GetFingerLevelEnumMap 获取指纹层级枚举映射
//func (a *ToolsManage) GetFingerLevelEnumMap() map[int]string {
//	return map[int]string{
//		enums.FingerLevelApp:      "应用层",
//		enums.FingerLevelSupport:  "支撑层",
//		enums.FingerLevelService:  "服务层",
//		enums.FingerLevelSystem:   "系统层",
//		enums.FingerLevelHardware: "硬件层",
//		enums.FingerLevelDefault:  "默认层",
//	}
//}

// GetFingerEnums 获取指纹枚举数据
func (a *ToolsManage) GetFingerEnums() map[string][]typespec.GlobalOptionsItemRes {
	return map[string][]typespec.GlobalOptionsItemRes{
		"class": {
			{Value: 1, Label: "中间件"},
			{Value: 2, Label: "蜜罐"},
			{Value: 3, Label: "主流应用框架"},
			{Value: 4, Label: "协同办公套件"},
			{Value: 5, Label: "主流CMS"},
			{Value: 6, Label: "网络设备"},
			{Value: 7, Label: "安全设备"},
			{Value: 8, Label: "数据库服务"},
			{Value: 9, Label: "代码研发"},
			{Value: 10, Label: "大数据平台"},
			{Value: 11, Label: "虚拟化服务"},
			{Value: 12, Label: "主流第三方服务"},
			{Value: 13, Label: "IOT"},
			{Value: 14, Label: "办公设备"},
			{Value: 15, Label: "邮件服务器"},
			{Value: 16, Label: "集权管控类"},
			{Value: 17, Label: "主机操作系统"},
			{Value: 18, Label: "功能类型"},
			{Value: 19, Label: "通用漏洞检测"},
			{Value: 20, Label: "云计算类型"},
			{Value: 21, Label: "web应用类型"},
			{Value: 23, Label: "缓存工具"},
			{Value: 25, Label: "js图形框架"},
			{Value: 26, Label: "移动框架"},
			{Value: 27, Label: "编程语言"},
			{Value: 28, Label: "操作系统"},
			{Value: 29, Label: "搜索引擎"},
			{Value: 30, Label: "网页邮件"},
			{Value: 31, Label: "cdn"},
			{Value: 32, Label: "营销自动化"},
			{Value: 35, Label: "管理分析与规划系统"},
			{Value: 36, Label: "广告网络"},
			{Value: 37, Label: "网络设备"},
			{Value: 38, Label: "媒体服务器"},
			{Value: 42, Label: "标签管理系统"},
			{Value: 44, Label: "构建CI系统"},
			{Value: 45, Label: "控制系统"},
			{Value: 46, Label: "远程访问系统"},
			{Value: 47, Label: "开发工具"},
			{Value: 48, Label: "网络存储"},
			{Value: 65, Label: "负载均衡"},
			{Value: 69, Label: "应用服务"},
			{Value: 71, Label: "办公自动化"},
			{Value: 72, Label: "其他"},
		},
		"isDev": {
			{Value: 1, Label: "否"},
			{Value: 2, Label: "是"},
		},
		"level": {
			{Value: 0, Label: "未分层"},
			{Value: 1, Label: "硬件层"},
			{Value: 2, Label: "系统层"},
			{Value: 3, Label: "服务层"},
			{Value: 4, Label: "支撑层"},
			{Value: 5, Label: "应用层"},
		},
		"softOrHard": {
			{Value: 0, Label: "软件"},
			{Value: 1, Label: "硬件"},
		},
	}
}
