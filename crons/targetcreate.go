package crons

import (
	"context"
	"encoding/json"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/network"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TaskControl struct {
	MaxRunningTaskNumber   int64 `json:"curTasks"`
	MaxRunningTargetNumber int64 `json:"curIp"`
}

// TargetCreate 主要的任务创建和管理函数
func TargetCreate() {
	var (
		ctx               = context.Background()
		taskControl       TaskControl
		taskTask          services.TaskTask
		taskTarget        services.TaskTarget
		taskConfiguration services.TaskConfiguration
		mapSetService     services.MapSet
	)

	// 初始化三大模块
	vulnerabilityScanner := services.NewVulnerabilityScanner()

	//第一步：加载配置
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.CurTasksInfoMapSetObjKey)
	if err != nil {
		log.Error("TargetCreate load config from db err", err)
		return
	}
	if objValueStr != "" {
		if err = json.Unmarshal([]byte(objValueStr), &taskControl); err != nil {
			log.Error("TargetCreate unmarshal config err", err)
			return
		}
	} else {
		// 默认配置
		taskControl.MaxRunningTaskNumber = 5
		taskControl.MaxRunningTargetNumber = 10
	}
	//第二步：查询目前运行中的任务和目标
	taskIds, taskNum, targetNum, err := taskTask.GetRunningTaskAndTarget(ctx, enums.TaskStatusRunning, enums.TargetStatusRunning)
	if err != nil {
		log.Error("TargetCreate GetRunningTaskAndTarget err", err)
		return
	}
	if taskNum >= taskControl.MaxRunningTaskNumber && targetNum >= taskControl.MaxRunningTargetNumber { //如果运行数量大于配置设置，直接返回
		return
	}
	//第三步：启动任务
	if taskNum < taskControl.MaxRunningTaskNumber {
		startNum := taskControl.MaxRunningTaskNumber - taskNum //启动个数
		newTaskIds, err := taskTask.GetWaitingTaskByLimit(ctx, int(startNum))
		if err != nil {
			log.Error("TargetCreate GetWaitingTaskByLimit err", err)
			return
		}
		err = taskTask.StartTasks(ctx, newTaskIds)
		if err != nil {
			log.Error("TargetCreate StartTasks err", err)
			return
		}
		taskIds = append(taskIds, newTaskIds...)
	}
	//第四步：启动目标
	if len(taskIds) == 0 {
		return
	}
	if targetNum < taskControl.MaxRunningTargetNumber {
		startNum := taskControl.MaxRunningTargetNumber - targetNum //启动个数
		targetRes := taskTarget.GetTargetsByTaskIdsStatus(ctx, taskIds, enums.TargetStatusToBegin, int(startNum))
		if len(targetRes) == 0 {
			return
		}

		for i := 0; i < len(targetRes); i++ {
			currentTarget := targetRes[i]
			var configJson enums.ConfigJson
			err = json.Unmarshal([]byte(currentTarget.TaskTemplateJSON), &configJson)
			if err != nil {
				log.Error("TargetCreate unmarshal configJson err", err)
				continue
			}

			// 存活探测
			var isAlive int = enums.TargetIsAliveN
			host := currentTarget.TargetURL
			parsedURL := network.ParseUrl2(host)
			if parsedURL != nil {
				host = parsedURL.Host
			}
			if strings.Contains(host, ":") {
				host = strings.Split(host, ":")[0]
			}

			ip, err := network.ResolveDomain(host)
			if err != nil {
				log.Infof("TargetCreate: failed to resolve domain %s: %v, skipping", host, err)
				_ = taskTarget.UpdateStatus(ctx, currentTarget.ID, enums.TargetStatusFinish)
				continue
			}

			alive, err := network.IcmpPing(ip)
			if err != nil {
				log.Infof("TargetCreate: IcmpPing for %s (%s) failed: %v, skipping", currentTarget.TargetURL, ip, err)
			}
			if alive {
				isAlive = enums.TargetIsAliveY
			}

			// 如果没有配置漏洞ID，则从任务配置中获取
			if configJson.VulIdsConfig == nil {
				taskModel, _ := taskTask.GetTaskByTaskId(ctx, currentTarget.TaskID)
				configResult, _ := taskConfiguration.GetConfigInfoByIdAndConfigKey(ctx, taskModel.TaskTemplateID, enums.ConfigJsonVulIdsKey)
				tempVulIds := make([]int, 0)
				err = json.Unmarshal([]byte(configResult.ConfigJson), &tempVulIds)
				if err == nil {
					configJson.VulIdsConfig = tempVulIds
				}
			}

			// 使用任务监听模块处理目标
			go func(target mysqls.TaskTarget, config enums.ConfigJson, aliveStatus int) {
				// 修改目标状态,添加目标日志,返回必要信息
				logId, err := taskTarget.UpdateTargetAndSaveTargetLog(ctx, target.TaskID, target.ID, enums.TargetStatusRunning, aliveStatus, target.TargetURL)
				if err != nil {
					log.Info("TargetCreate UpdateTargetAndSaveTargetlog err：", err)
					return
				}
				// 根据配置模式启动扫描
				if config.Mode.Mode == enums.TaskConfigurationModeTarget {
					// 如果是定向渗透，进行分布式扫描
					go func() {
						err := vulnerabilityScanner.StartDistributedScan(ctx, logId, target, config)
						if err != nil {
							var taskLogInfo services.TaskLogInfo
							_ = taskLogInfo.AddTaskLogInfo(ctx, target.TaskID, target.ID, logId, target.TargetURL, "distributed_scan", "start distributed scan failed: "+err.Error())
							if updateErr := taskTarget.UpdateTargetAndLogStateById(ctx, target.ID, enums.TargetStatusFinish, enums.TargetStatusFinish); updateErr != nil {
								log.Errorf("update target status after distributed scan start failed: %v", updateErr)
							}
							runningTargets := taskTarget.GetTargetsByTaskIdAndStatusList(ctx, target.TaskID, []int{enums.TargetStatusToBegin, enums.TargetStatusRunning})
							if len(runningTargets) == 0 {
								taskTask.FinishTask(ctx, target.TaskID)
							}
							log.Errorf("启动分布式扫描失败: %v", err)
						}
					}()
				} else {
					// 启动本地扫描器
					go vulnerabilityScanner.StartLocalScan(ctx, logId, target, config)
				}
			}(currentTarget, configJson, isAlive)
		}
	}
}
