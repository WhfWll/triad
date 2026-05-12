package crons

import (
	"context"
	"encoding/json"
	"smart/models/redises"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/network"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func LogicTaskExec() {
	var (
		srv         services.Logic
		ctx         = context.Background()
		taskControl TaskControl
	)
	//第一步：加载配置
	//if err := config.Load("task_control", &taskControl); err != nil {
	//	log.Error("LogicTaskExec load config err", err)
	//	return
	//}
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.CurTasksInfoMapSetObjKey)
	if err != nil {
		log.Error("LogicTaskExec load config from db err", err)
		return
	}
	if objValueStr != "" {
		if err = json.Unmarshal([]byte(objValueStr), &taskControl); err != nil {
			log.Error("LogicTaskExec unmarshal config err", err)
			return
		}
	} else {
		taskControl.MaxRunningTaskNumber = 5
		taskControl.MaxRunningTargetNumber = 10
	}
	//第二步 查询待执行目标
	targetList, _ := srv.GetTargetsByStatus(ctx, enums.LogicTaskStatusBegin)
	if len(targetList) == 0 {
		return
	}
	//第三步 获取正在运行的目标
	runningTargetList, _ := srv.GetTargetsByStatus(ctx, enums.LogicTaskStatusRunning)
	delta := taskControl.MaxRunningTargetNumber - int64(len(runningTargetList))
	if delta <= 0 {
		return
	}
	//第四步 遍历待执行目标进行测试
	for _, target := range targetList[:delta] {
		if target.ID == 0 {
			continue
		}
		task, err := srv.GetTaskById(ctx, target.TaskID)
		var callId string
		if target.Type == enums.LogicTypeBeyondPermission {
			// 进行越权漏洞测试
			var beyondPermConfig services.BeyondPermConfig
			err = json.Unmarshal([]byte(task.ScanConfig), &beyondPermConfig)
			if err != nil {
				log.Error("LogicTaskExec json unmarshal error:", err.Error())
			}
			callId, err = srv.BeyondPermCall(ctx, target.TargetURL, beyondPermConfig)
			if err != nil {
				log.Error("LogicTaskExec BeyondPermissionCall error:", err.Error())
			}
		} else if target.Type == enums.LogicTypeTraverseTesting {
			// 进行信息遍历漏洞测试
			var sensInfoConfig services.SensInfoConfig
			err = json.Unmarshal([]byte(task.ScanConfig), &sensInfoConfig)
			if err != nil {
				log.Error("LogicTaskExec json unmarshal error:", err.Error())
			}
			callId, err = srv.TraverseTestingCall(ctx, target.TargetURL, sensInfoConfig)
			if err != nil {
				log.Error("LogicTaskExec BeyondPermissionCall error:", err.Error())
			}
		} else if target.Type == enums.LogicTypeUnAuthAccess {
			// 进行信息遍历漏洞测试
			var unAuthAccessConfig services.UnAuthAccessConfig
			err = json.Unmarshal([]byte(task.ScanConfig), &unAuthAccessConfig)
			if err != nil {
				log.Error("LogicTaskExec json unmarshal error:", err.Error())
			}
			callId, err = srv.UnAuthAccessCall(ctx, target.TargetURL, unAuthAccessConfig)
			if err != nil {
				log.Error("LogicTaskExec BeyondPermissionCall error:", err.Error())
			}
		} else {
			log.Println("find not support script type", target.Type)
			continue
		}
		// 设置缓存信息
		var cacheClient redises.RedisHash
		err = cacheClient.SetHashHset(ctx, enums.LogicCallIdCacheKey, strconv.Itoa(target.ID), callId)
		// 保存日志信息
		logId, err := srv.LogSave(ctx, target.TargetURL, task.ID, target.ID)
		if err != nil {
			log.Error("LogicTaskExec LogSave error:", err.Error())
		}
		// 监听处理越权漏洞的检测结果
		go srv.HandleResult(ctx, callId, task.ID, target.ID, logId, target.Type, target.TargetURL)
		//更新任务状态
		err = srv.UpdateTaskParam(ctx, task.ID, map[string]interface{}{"status": enums.LogicTaskStatusRunning, "update_time": time.Now()})
		if err != nil {
			log.Error("LogicTaskExec UpdateLogicTaskById error:", err.Error())
		}
		//更新目标状态
		err = srv.UpdateTargetParam(ctx, target.ID, map[string]interface{}{"status": enums.TargetStatusRunning, "update_time": time.Now()})
		if err != nil {
			log.Error("LogicTaskExec UpdateTargetStatus error:", err.Error())
		}
		//更新日志状态
		err = srv.UpdateLogStatus(ctx, target.ID, map[string]interface{}{"status": enums.TargetStatusRunning})
		if err != nil {
			log.Error("LogicTaskExec update status error:", err.Error())
		}
		//更新存活状态
		go UpdateIsAlive(ctx, target.TargetURL, target.ID, srv)
	}
}

func UpdateIsAlive(ctx context.Context, url string, targetId int, srv services.Logic) {
	_, host, port, err := network.ParseUrl(url)
	if err != nil {
		return
	}
	isAlive := network.TelnetIsOpen(host, port)
	if isAlive {
		err = srv.UpdateTargetParam(ctx, targetId, map[string]interface{}{"is_alive": enums.TargetIsAliveY, "update_time": time.Now()})
		if err != nil {
			log.Error("LogicTaskExec UpdateTargetStatus error:", err.Error())
		}
		err = srv.UpdateLogStatus(ctx, targetId, map[string]interface{}{"is_alive": enums.TargetIsAliveY})
		if err != nil {
			log.Error("LogicTaskExec update status error:", err.Error())
		}
	}
}
