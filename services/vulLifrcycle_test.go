package services

import (
	"encoding/json"
	"fmt"
	"smart/api/typespec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_appendLogToTask 测试核心 JSON 拼接逻辑
// 由于 appendLogToTask 是私有方法，测试文件必须和 services 在同一个包下 (package services)
func Test_appendLogToTask(t *testing.T) {
	service := &VulLifecycle{}

	// 模拟日志内容
	logMsg1 := "2025-12-23 14:00:00 Admin 发现漏洞"
	logMsg2 := "2025-12-23 14:10:00 Admin 确认漏洞"
	logMsg3 := "2025-12-23 13:00:00 YTest 发现漏洞 (任务B)"

	t.Run("场景1: 首次插入(原有Content为空)", func(t *testing.T) {
		// 模拟第一次写入
		resultJSON := service.appendLogToTask("", 101, "渗透任务A", logMsg1, "2025-12-23 14:00:00")

		// 验证解析
		var logs []typespec.TaskLifecycleLog
		err := json.Unmarshal([]byte(resultJSON), &logs)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(logs))
		assert.Equal(t, 101, logs[0].TaskID)
		assert.Equal(t, "渗透任务A", logs[0].TaskName)
		assert.Equal(t, 1, len(logs[0].Lifecycle))
		assert.Equal(t, logMsg1, logs[0].Lifecycle[0].Content)

		fmt.Printf("场景1 JSON结果: %s\n", resultJSON)
	})

	t.Run("场景2: 同一任务追加日志(原有Content不为空)", func(t *testing.T) {
		// 先构造初始状态
		initialLogs := []typespec.TaskLifecycleLog{
			{
				TaskID:   101,
				TaskName: "渗透任务A",
				Lifecycle: []typespec.LifecycleItem{
					{Time: "2025-12-23 14:00:00", Content: logMsg1},
				},
			},
		}
		initialJSONBytes, _ := json.Marshal(initialLogs)
		initialJSON := string(initialJSONBytes)

		// 执行追加
		resultJSON := service.appendLogToTask(initialJSON, 101, "渗透任务A", logMsg2, "2025-12-23 14:10:00")

		// 验证
		var logs []typespec.TaskLifecycleLog
		json.Unmarshal([]byte(resultJSON), &logs)

		assert.Equal(t, 1, len(logs), "应该仍然只有一个任务块")
		assert.Equal(t, 2, len(logs[0].Lifecycle), "应该有两条日志")
		assert.Equal(t, logMsg1, logs[0].Lifecycle[0].Content)
		assert.Equal(t, logMsg2, logs[0].Lifecycle[1].Content, "第二条日志应该是追加的内容")

		fmt.Printf("场景2 JSON结果: %s\n", resultJSON)
	})

	t.Run("场景3: 新任务ID写入(追加新任务块)", func(t *testing.T) {
		// 先构造初始状态 (任务A)
		initialLogs := []typespec.TaskLifecycleLog{
			{
				TaskID:   101,
				TaskName: "渗透任务A",
				Lifecycle: []typespec.LifecycleItem{
					{Time: "2025-12-23 14:00:00", Content: logMsg1},
				},
			},
		}
		initialJSONBytes, _ := json.Marshal(initialLogs)
		initialJSON := string(initialJSONBytes)
		// 执行写入 (任务B，ID不同)
		resultJSON := service.appendLogToTask(initialJSON, 202, "渗透任务B", logMsg3, "2025-12-23 13:00:00")
		// 验证
		var logs []typespec.TaskLifecycleLog
		json.Unmarshal([]byte(resultJSON), &logs)

		assert.Equal(t, 2, len(logs), "应该有两个任务块")
		// 检查任务A是否还在
		assert.Equal(t, 101, logs[0].TaskID)
		// 检查任务B是否追加
		assert.Equal(t, 202, logs[1].TaskID)
		assert.Equal(t, "渗透任务B", logs[1].TaskName)
		assert.Equal(t, logMsg3, logs[1].Lifecycle[0].Content)
		fmt.Printf("场景3 JSON结果: %s\n", resultJSON)
	})

	t.Run("场景4: 任务名称更新", func(t *testing.T) {
		// 初始：任务ID 101 叫 "旧名称"
		initialLogs := []typespec.TaskLifecycleLog{
			{
				TaskID:   101,
				TaskName: "旧名称",
				Lifecycle: []typespec.LifecycleItem{
					{Time: "2025-12-23 14:00:00", Content: logMsg1},
				},
			},
		}
		initialJSONBytes, _ := json.Marshal(initialLogs)

		// 执行写入：任务ID 101 叫 "新名称"
		resultJSON := service.appendLogToTask(string(initialJSONBytes), 101, "新名称", logMsg2, "2025-12-23 14:10:00")

		var logs []typespec.TaskLifecycleLog
		json.Unmarshal([]byte(resultJSON), &logs)

		assert.Equal(t, "新名称", logs[0].TaskName, "任务名称应该被更新")
	})
}

var mockDB = make(map[string]*MockVulLifecycleRow)

type MockVulLifecycleRow struct {
	ID       int
	PocName  string
	Name     string
	Location string
	Content  string // 也就是数据库里的 JSON 字符串
	FindNum  int
}

func simulateAppendLogToTask(currentContent string, taskId int, taskName string, newLogMsg string, timeStr string) string {
	var logs []typespec.TaskLifecycleLog
	if currentContent != "" {
		_ = json.Unmarshal([]byte(currentContent), &logs)
	}

	found := false
	newItem := typespec.LifecycleItem{
		Time:    timeStr,
		Content: newLogMsg,
	}

	for i, item := range logs {
		if item.TaskID == taskId {
			logs[i].Lifecycle = append(logs[i].Lifecycle, newItem)
			if taskName != "" {
				logs[i].TaskName = taskName
			}
			found = true
			break
		}
	}

	if !found {
		newBlock := typespec.TaskLifecycleLog{
			TaskID:    taskId,
			TaskName:  taskName,
			Lifecycle: []typespec.LifecycleItem{newItem},
		}
		logs = append(logs, newBlock)
	}

	bytes, _ := json.Marshal(logs) // 正常代码用 Marshal
	// 为了控制台打印好看，这里模拟用 MarshalIndent，实际存库不用 Indent
	// bytes, _ := json.MarshalIndent(logs, "", "  ")
	return string(bytes)
}

// ==========================================
// 4. 模拟 AddVulLifecycle 函数
// ==========================================
func Simulate_AddVulLifecycle(pocName, vulName, location, content, timeStr, taskName string, taskID int) {
	fmt.Printf("\n>>> 正在执行写入请求: TaskID=%d, TaskName=%s, Log=%s\n", taskID, taskName, content)

	// 1. 生成唯一Key，模拟数据库联合主键查询
	uniqueKey := pocName + "|" + vulName + "|" + location

	row, exists := mockDB[uniqueKey]

	if exists {
		// === 情况 A: 记录已存在 (Update) ===
		fmt.Println("   [DB状态] 记录已存在，执行更新逻辑...")

		// 核心：取出旧 Content，追加新 Content
		newContent := simulateAppendLogToTask(row.Content, taskID, taskName, content, timeStr)

		// 模拟 Update 操作
		row.Content = newContent
		row.FindNum++
		fmt.Println("   [DB操作] Update 成功")

	} else {
		// === 情况 B: 记录不存在 (Insert) ===
		fmt.Println("   [DB状态] 记录不存在，执行插入逻辑...")

		// 核心：传入空字符串，生成初始 JSON
		initialContent := simulateAppendLogToTask("", taskID, taskName, content, timeStr)

		// 模拟 Insert 操作
		newRow := &MockVulLifecycleRow{
			ID:       len(mockDB) + 1, // 模拟自增ID
			PocName:  pocName,
			Name:     vulName,
			Location: location,
			Content:  initialContent,
			FindNum:  1,
		}
		mockDB[uniqueKey] = newRow
		fmt.Println("   [DB操作] Insert 成功")
	}

	// 打印当前数据库里存的 JSON 内容
	printDBContent(uniqueKey)
}

// 辅助打印函数
func printDBContent(key string) {
	row := mockDB[key]
	var prettyJSON interface{}
	json.Unmarshal([]byte(row.Content), &prettyJSON)
	prettyBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Printf("   [当前数据库存储的 Content]:\n%s\n", string(prettyBytes))
}

// ==========================================
// 5. 执行测试
// ==========================================
func TestRunSimulation(t *testing.T) {
	// 基础信息
	poc := "thinkphp_rce"
	name := "ThinkPHP远程代码执行"
	loc := "http://example.com"

	// --- 模拟第一次扫描 (任务A) ---
	// 期望：数据库插入一条新记录，JSON里有一个任务块
	Simulate_AddVulLifecycle(poc, name, loc,
		"System 发现漏洞 (首次)",
		"2025-01-01 10:00:00",
		"渗透任务A", 101)

	// --- 模拟第二次扫描 (任务A，同一个任务扫出的) ---
	// 期望：数据库记录更新，JSON里依然是一个任务块，但 lifecycle 数组变长了
	Simulate_AddVulLifecycle(poc, name, loc,
		"System 再次发现漏洞 (复测)",
		"2025-01-01 12:00:00",
		"渗透任务A", 101)

	// --- 模拟第三次操作 (人工操作，或者新任务B) ---
	// 期望：数据库记录更新，JSON里增加了一个新的任务块 (TaskID=202)
	Simulate_AddVulLifecycle(poc, name, loc,
		"Admin 确认漏洞存在",
		"2025-01-02 09:00:00",
		"人工复核", 202)
}
