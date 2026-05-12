package services

import (
	"context"
	"encoding/json"
	"fmt"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

type VulLifecycle struct{}

// GetList 获取生命周期列表（分页+搜索）
func (v *VulLifecycle) GetList(ctx context.Context, page, limit int, pocName, name string) ([]mysqls.VulLifecycle, int64, error) {
	var vulLifecycleModel mysqls.VulLifecycle
	return vulLifecycleModel.GetList(ctx, page, limit, pocName, name)
}

// GetDetail 获取单个详情（通过联合主键）
func (v *VulLifecycle) GetDetail(ctx context.Context, pocName, name, location string) (mysqls.VulLifecycle, error) {
	var vulLifecycleModel mysqls.VulLifecycle
	return vulLifecycleModel.GetVulLifecycleDetail(ctx, pocName, name, location)
}

// UpdateContentByUniqueKey 更新周期信息（通过联合主键）
func (v *VulLifecycle) UpdateContentByUniqueKey(ctx context.Context, pocName, name, location, content string) error {
	var vulLifecycleModel mysqls.VulLifecycle
	return vulLifecycleModel.UpdateContentByUniqueKey(ctx, pocName, name, location, content)
}

// UpdateContentById 更新周期信息（通过ID）
// 场景：前端列表页编辑时，通常传 ID
func (v *VulLifecycle) UpdateContentById(ctx context.Context, id int, content string) error {
	var vulLifecycleModel mysqls.VulLifecycle
	return vulLifecycleModel.UpdateContentById(ctx, id, content)
}

// AddVulLifecycle 添加漏洞周期任务
func (v *VulLifecycle) AddVulLifecycle(ctx context.Context, pocName, vulName, location, content, taskName string, taskID int) (int, error) {
	var vulLifecycleModel mysqls.VulLifecycle
	now := time.Now().Format(enums.TimeLayout)
	// 判断漏洞是否存在 如果存在就更新content即可
	vulLifecycleInfo, _ := vulLifecycleModel.GetVulLifecycleDetail(ctx, pocName, vulName, location)
	if vulLifecycleInfo.ID != 0 {
		// 拼接新的content update
		vulLifecycleModel.UpdateContentByUniqueKey(ctx, pocName, vulName, location, v.appendLogToTask(vulLifecycleInfo.Content, taskID, taskName, content, now))
	} else {
		vulLifecycleModel.PocName = pocName
		vulLifecycleModel.VulName = vulName
		vulLifecycleModel.Location = location
		vulLifecycleModel.Content = v.appendLogToTask("", taskID, taskName, content, now)
		vulLifecycleModel.FindNum = 1
		if _, err := vulLifecycleModel.Add(ctx); err != nil {
			return 0, err
		}
	}
	return 0, nil
}

// ConstructDiscoveryContent 构造【发现漏洞】的日志内容
// 格式：时间 + 用户 【user】 创建【task_name】，发现目标资产【target_url】存在【vul_name】，状态为 【status】
func (v *VulLifecycle) ConstructDiscoveryContent(username, taskName, targetUrl, vulName string, status int) string {
	if taskName == "" {
		taskName = "未知任务"
	}
	if username == "" {
		username = "未知用户"
	}
	return fmt.Sprintf("用户 %s 创建%s，发现目标资产%s存在%s漏洞，状态为 %s",
		username, taskName, targetUrl, vulName, enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(status))
}

// ConstructStatusChangeContent 构造【状态变更】的日志内容
// 格式：时间 + 用户 【user】 更改漏洞状态为 【vul_status】
func (v *VulLifecycle) ConstructStatusChangeContent(username string, status int) string {
	if username == "" {
		username = "未知用户"
	}
	return fmt.Sprintf("用户 %s 更改漏洞状态为 %s",
		username, enums.ToolsVulnerabilityEnum.GetTaskVulStatusEnum(status))
}

// 辅助方法：将 Content 字符串解析为结构体列表
func parseContent(content string) []typespec.TaskLifecycleLog {
	var logs []typespec.TaskLifecycleLog
	if content == "" {
		return logs
	}
	_ = json.Unmarshal([]byte(content), &logs)
	return logs
}

// 辅助方法：将结构体列表转回 JSON 字符串
func formatContent(logs []typespec.TaskLifecycleLog) string {
	bytes, _ := json.Marshal(logs)
	return string(bytes)
}

// appendLogToTask 将一条新日志追加到指定 TaskID 的记录下
func (v *VulLifecycle) appendLogToTask(currentContent string, taskId int, taskName string, newLogMsg, timeStr string) string {
	logs := parseContent(currentContent)
	newItem := typespec.LifecycleItem{
		Time:    timeStr,
		Content: newLogMsg,
	}
	found := false
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
	return formatContent(logs)
}
