package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/invoke"
	"strconv"
	"time"
)

// 任务检测结果

type TaskResult struct {
}

// AddTaskResult 插入一条检测结果
func (t *TaskResult) AddTaskResult(ctx context.Context, target mysqls.TaskTarget, scannerLog invoke.ScannerLog, vulLibId, risk, checked int, name, result string) (int, error) {
	var taskResult mysqls.TaskResult
	taskResult.TaskID = target.TaskID
	taskResult.TargetID = target.ID
	taskResult.NodeID = scannerLog.ScriptExecutionID
	taskResult.FatherID = scannerLog.ParentScriptID
	taskResult.Pocname = scannerLog.Pocname
	taskResult.DecisionLibId = vulLibId
	taskResult.VulName = name
	taskResult.Risk = risk
	taskResult.Result = result
	taskResult.Checked = checked
	taskResult.CreateTime = time.Now()
	err := taskResult.AddTaskResult(ctx)
	if err != nil {
		return 0, errors.New("AddTaskResult error: " + err.Error())
	}
	return taskResult.ID, nil
}

// UpdateTaskResult 依据任务ID与目标ID与pocname 更新结果表的状态，漏洞结果，请求、响应报文，是否检测出来
func (t *TaskResult) UpdateTaskResult(ctx context.Context, scriptResult ScriptResult) (int, error) {
	targetId, err := strconv.Atoi(scriptResult.ObjId)
	if err != nil {
		return 0, errors.New("UpdateTaskResult error: " + err.Error())
	}

	if scriptResult.NodeId == "" {
		return 0, errors.New("UpdateTaskResult error: nodeId不可为空，其中targetId为" + scriptResult.ObjId)
	}

	var taskResultModel mysqls.TaskResult
	taskResult, _ := taskResultModel.GetByTargetIdAndNodeId(ctx, targetId, scriptResult.NodeId)
	if taskResult.ID == 0 {
		return 0, nil
		// 这个报错看起来是多余的，因为如果没找到数据，那么下面的代码就不会执行，所以就注释掉了
		//return 0, errors.New("UpdateTaskResult error: 依据targetId与nodeID未查询到结果数据，其中targetId为" + strconv.Itoa(targetId) + " nodeId为" + scriptResult.NodeId)
	}

	// 如果是[端口扫描 | web路径爆破]的结果，需要追加，而不是覆盖，用于路径图展示
	switch scriptResult.Script.Pocname {
	case "port_scan", "web_dir_path_scan":
		fmt.Println("返回多条数据:"+scriptResult.Script.Pocname, scriptResult.Result.Detail)
		// 返回的result
		if scriptResult.Result.Detail != "" {
			// 解析返回的result
			resResultMap := make(map[string]string, 0)
			if err = json.Unmarshal([]byte(scriptResult.Result.Detail), &resResultMap); err == nil {
				fmt.Println("返回多条数据:" + scriptResult.Script.Pocname + " 解析成功")
				// 当前已储存的数据
				finalSaveResultData := make([]map[string]string, 0)
				if taskResult.Result != "" {
					json.Unmarshal([]byte(taskResult.Result), &finalSaveResultData)
					finalSaveResultData = append(finalSaveResultData, resResultMap)
				} else {
					finalSaveResultData = append(finalSaveResultData, resResultMap)
				}

				finalResultString, _ := json.Marshal(finalSaveResultData)
				fmt.Println("返回多条数据:" + scriptResult.Script.Pocname + " 追加成功：" + string(finalResultString))

				taskResult.Result = string(finalResultString)
			}
		}
	default:
		taskResult.Result = scriptResult.Result.Detail
	}

	taskResult.Request = scriptResult.Result.Request
	taskResult.Response = scriptResult.Result.Response
	taskResult.Checked = enums.TaskResultCheckedY
	taskResult.UpdateTaskResult(ctx)

	return taskResult.ID, nil
}

// DelTaskInfoByTaskId 删除任务检测结果表 通过task_id
func (t *TaskResult) DelTaskInfoByTaskId(ctx context.Context, taskIds []int) error {
	var taskResultModel mysqls.TaskResult
	return taskResultModel.DeleteByTaskIds(ctx, taskIds)
}

// GetTaskResultList 返回漏洞取证结果信息
func (t *TaskResult) GetTaskResultList(ctx context.Context, taskID, riskType, page, size int, search, targetIDs string) ([]mysqls.TaskResult, int64, error) {
	var taskResult mysqls.TaskResult
	res, count, err := taskResult.GetTaskResultList(ctx, taskID, riskType, page, size, search, targetIDs)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

// AllTaskResultByTargetId 依据目标ID获取所有结果数据
func (t *TaskResult) AllTaskResultByTargetId(ctx context.Context, targetId int, pocname string) []mysqls.TaskResult {
	var taskResult mysqls.TaskResult
	return taskResult.AllByTargetId(ctx, targetId, pocname)
}

// GetTaskResultById 依据目标ID获取所有结果数据
func (t *TaskResult) GetTaskResultById(ctx context.Context, resultId int) (mysqls.TaskResult, error) {
	var taskResult mysqls.TaskResult
	taskResult.ID = resultId
	return taskResult.GetTaskResult(ctx)
}

// UpdateTaskResultRisk updates the risk level of a task result by its node ID
func (t *TaskResult) UpdateTaskResultRisk(ctx context.Context, nodeID string, risk int) error {
	var taskResult mysqls.TaskResult
	err := taskResult.UpdateRiskByNodeID(ctx, nodeID, risk)
	if err != nil {
		return errors.New("UpdateTaskResultRisk error: " + err.Error())
	}
	return nil
}
