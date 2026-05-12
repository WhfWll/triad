package services

import (
	"context"
	"encoding/json"
	"fmt"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/network"
	"strings"

	log "github.com/sirupsen/logrus"
)

// TaskControl 任务控制配置
type TaskControl struct {
	MaxRunningTaskNumber   int64 `json:"curTasks"`
	MaxRunningTargetNumber int64 `json:"curIp"`
}

// TaskListener 任务监听器
type TaskListener struct {
	taskTask          TaskTask
	taskTarget        TaskTarget
	taskConfiguration TaskConfiguration
	vulScanner        *VulnerabilityScanner
}

// NewTaskListener 创建新的任务监听器
func NewTaskListener() *TaskListener {
	return &TaskListener{
		vulScanner: NewVulnerabilityScanner(),
	}
}

// StartTaskListener 启动任务监听器
func (tl *TaskListener) StartTaskListener(ctx context.Context) {
	var taskControl TaskControl

	// 第一步：加载配置
	//if err := config.Load("task_control", &taskControl); err != nil {
	//	log.Error("TaskListener load config err", err)
	//	return
	//}
	var mapSetService MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.CurTasksInfoMapSetObjKey)
	if err != nil {
		log.Error("TaskListener load config from db err", err)
		return
	}
	if objValueStr != "" {
		if err = json.Unmarshal([]byte(objValueStr), &taskControl); err != nil {
			log.Error("TaskListener unmarshal config err", err)
			return
		}
	} else {
		taskControl.MaxRunningTaskNumber = 5
		taskControl.MaxRunningTargetNumber = 10
	}

	// 第二步：查询目前运行中的任务和目标
	taskIds, taskNum, targetNum, err := tl.taskTask.GetRunningTaskAndTarget(ctx, enums.TaskStatusRunning, enums.TargetStatusRunning)
	if err != nil {
		log.Error("TaskListener GetRunningTaskAndTarget err", err)
		return
	}

	if taskNum >= taskControl.MaxRunningTaskNumber && targetNum >= taskControl.MaxRunningTargetNumber {
		// 如果运行数量大于配置设置，直接返回
		return
	}

	// 第三步：启动任务
	if taskNum < taskControl.MaxRunningTaskNumber {
		startNum := taskControl.MaxRunningTaskNumber - taskNum
		newTaskIds, err := tl.taskTask.GetWaitingTaskByLimit(ctx, int(startNum))
		if err != nil {
			log.Error("TaskListener GetWaitingTaskByLimit err", err)
			return
		}
		err = tl.taskTask.StartTasks(ctx, newTaskIds)
		if err != nil {
			log.Error("TaskListener StartTasks err", err)
			return
		}
		taskIds = append(taskIds, newTaskIds...)
	}

	// 第四步：启动目标
	if len(taskIds) == 0 {
		return
	}

	if targetNum < taskControl.MaxRunningTargetNumber {
		startNum := taskControl.MaxRunningTargetNumber - targetNum
		targetRes := tl.taskTarget.GetTargetsByTaskIdsStatus(ctx, taskIds, enums.TargetStatusToBegin, int(startNum))
		if len(targetRes) == 0 {
			return
		}

		for i := 0; i < len(targetRes); i++ {
			target := targetRes[i]
			// 启动目标处理
			go tl.processTarget(ctx, target)
		}
	}
}

// processTarget 处理单个目标
func (tl *TaskListener) processTarget(ctx context.Context, target mysqls.TaskTarget) {
	// 修改目标状态,添加目标日志,返回必要信息
	logId, err := tl.taskTarget.UpdateTargetAndSaveTargetLog(ctx, target.TaskID, target.ID, enums.TargetStatusRunning, target.IsAlive, target.TargetURL)
	if err != nil {
		log.Info("TaskListener UpdateTargetAndSaveTargetlog err：", err)
		return
	}

	var configJson enums.ConfigJson
	err = json.Unmarshal([]byte(target.TaskTemplateJSON), &configJson)
	if err != nil {
		fmt.Println("TaskListener json.Unmarshal err：", err)
		return
	}

	// 如果没有配置漏洞ID，则从任务配置中获取
	if configJson.VulIdsConfig == nil {
		taskModel, _ := tl.taskTask.GetTaskByTaskId(ctx, target.TaskID)
		configResult, _ := tl.taskConfiguration.GetConfigInfoByIdAndConfigKey(ctx, taskModel.TaskTemplateID, enums.ConfigJsonVulIdsKey)
		tempVulIds := make([]int, 0)
		err = json.Unmarshal([]byte(configResult.ConfigJson), &tempVulIds)
		if err == nil {
			configJson.VulIdsConfig = tempVulIds
		}
	}

	// 启动存活探测
	go tl.performAliveProbe(ctx, target, configJson)

	// 启动扫描任务
	if configJson.Mode.Mode == enums.TaskConfigurationModeTarget {
		// 如果是定向渗透，进行分布式扫描
		go func() {
			if err := tl.vulScanner.StartDistributedScan(ctx, logId, target, configJson); err != nil {
				log.Errorf("start distributed scan failed: %v", err)
				var taskLogInfo TaskLogInfo
				_ = taskLogInfo.AddTaskLogInfo(ctx, target.TaskID, target.ID, logId, target.TargetURL, "distributed_scan", "start distributed scan failed: "+err.Error())
				tl.HandleScanComplete(ctx, target)
			}
		}()
	} else {
		// 启动本地扫描器
		go tl.vulScanner.StartLocalScan(ctx, logId, target, configJson)
	}
}

// performAliveProbe 执行存活探测
func (tl *TaskListener) performAliveProbe(ctx context.Context, target mysqls.TaskTarget, configJson enums.ConfigJson) {
	// 先解析URL获取主机名
	url := network.ParseUrl2(target.TargetURL)
	host := url.Host
	// 如果主机名包含端口号，去除端口号部分
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	// 解析域名获取IP
	ip, err := network.ResolveDomain(host)
	if err != nil {
		log.Info("TaskListener ResolveDomain err：", err)
		// 如果解析失败，将存活状态设置为false
		err = tl.taskTarget.UpdateTargetAlive(ctx, target.ID, enums.TargetIsAliveN)
		if err != nil {
			log.Info("TaskListener UpdateTargetAndSaveTargetlog err：", err)
		}
		return
	}
	// 使用解析后的IP地址进行存活探测
	isAlive, err := tl.taskTarget.AliveProbe(ip, configJson)
	if err != nil {
		log.Info("TaskListener PingIpIsLive err：", err)
	}
	err = tl.taskTarget.UpdateTargetAlive(ctx, target.ID, isAlive)
	if err != nil {
		log.Info("TaskListener UpdateTargetAndSaveTargetlog err：", err)
	}
}

// HandleScanComplete 处理扫描完成逻辑
func (tl *TaskListener) HandleScanComplete(ctx context.Context, target mysqls.TaskTarget) {
	// 更新目标状态为已完成
	err := tl.taskTarget.UpdateTargetStateById(ctx, target.ID, enums.TargetStatusFinish)
	if err != nil {
		log.Error("更新目标状态失败: ", err)
	}

	// 更新任务日志状态
	var taskLog TaskLog
	err = taskLog.UpdateTaskLogStateByTargetId(ctx, target.ID, enums.TargetStatusFinish)
	if err != nil {
		log.Error("更新任务日志状态失败: ", err)
	}

	// 异步检查任务下的所有目标是否都已完成，如果是则更新任务状态
	go func() {
		// 检查该任务下是否还有运行中的目标
		runningTargets := tl.taskTarget.GetTargetsByTaskIdAndStatusList(ctx, target.TaskID, []int{enums.TargetStatusToBegin, enums.TargetStatusRunning})
		if len(runningTargets) == 0 {
			// 所有目标都已完成，更新任务状态
			tl.taskTask.FinishTask(ctx, target.TaskID)
		}
	}()
}
