package application

import (
	"context"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
)

type TaskGroup struct {
}

// Create 任务组新建
func (t *TaskGroup) Create(ctx context.Context, req *typespec.TaskGroupCreateReq, resp *typespec.TaskGroupCreateResp) error {
	var srv services.TaskGroup
	err := srv.Create(ctx, req.Name, req.Describe)
	if err != nil {
		return err
	}
	return nil
}

// List 任务组列表
func (t *TaskGroup) List(ctx context.Context, req *typespec.TaskGroupListReq, resp *typespec.TaskGroupListResp) error {
	var srv services.TaskGroup
	groups, count, err := srv.List(ctx, req.Search, req.Page, req.Size)
	if err != nil {
		return err
	}
	var groupList []typespec.TaskGroupListInfo
	for _, v := range groups {
		groupList = append(groupList, typespec.TaskGroupListInfo{
			ID:         v.ID,
			Name:       v.Name,
			HighNum:    v.HighNum,
			MiddleNum:  v.MiddleNum,
			LowNum:     v.LowNum,
			SafeNum:    v.SafeNum,
			CreateTime: v.CreateTime.Format(enums.TimeLayout),
			UpdateTime: v.CreateTime.Format(enums.TimeLayout),
			StatusNum:  v.Status,
			Status:     enums.GetTargetStatus(v.Status),
			Describe:   v.Describe,
		})
	}
	*resp = typespec.TaskGroupListResp{
		Page:  req.Page,
		Size:  req.Size,
		Total: count,
		List:  groupList,
	}
	return nil
}

// Delete 任务组删除
func (t *TaskGroup) Delete(ctx context.Context, req *typespec.TaskGroupDeleteReq, resp *typespec.TaskGroupDeleteResp) error {
	var srv services.TaskGroup
	for _, id := range strings.Split(req.Id, ",") {
		id = strings.ReplaceAll(id, " ", "")
		if id == "" {
			continue
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		err = srv.Delete(ctx, idInt)
		if err != nil {
			continue
		}
	}
	return nil
}

// GroupBind 任务组 任务与组绑定
func (t *TaskGroup) GroupBind(ctx context.Context, req *typespec.TaskGroupGroupBindReq, resp *typespec.TaskGroupGroupBindResp) error {
	var srv services.TaskGroup
	err := srv.GroupBind(ctx, req.TaskId, req.GroupId)
	if err != nil {
		return err
	}
	return nil
}

// TaskList 任务组 任务列表
func (t *TaskGroup) TaskList(ctx context.Context, req *typespec.TaskGroupTaskListReq, res *typespec.TaskGroupTaskListResp) error {
	var srv services.TaskGroup
	list, total, err := srv.TaskList(ctx, req.GroupId, req.Page, req.Size)
	if err != nil {
		return err
	}
	res.Total = total
	taskId := make([]int, 0, len(list))
	for _, item := range list {
		taskId = append(taskId, item.ID)

		// 目标的风险等级
		targetRiskList := make([]int, 4)
		targetRiskList[0] = item.HigeNum
		targetRiskList[1] = item.MiddleNum
		targetRiskList[2] = item.LowNum
		targetRiskList[3] = item.SafeNum

		res.List = append(res.List, typespec.TaskGroupTaskListItemRes{
			Id:              item.ID,
			TaskName:        item.TaskName,
			ExecuteType:     item.ExecuteType,
			ExecuteTypeName: enums.TaskTaskEnum.ExecTypeEnum(item.ExecuteType),
			RiskLevel:       item.RiskLevel,
			RiskLevelName:   enums.TaskTaskEnum.RiskEnum(item.RiskLevel),
			Status:          item.Status,
			StatusName:      enums.TaskTaskEnum.StatusEnum(item.Status),
			TargetRisk:      targetRiskList,
			CreateTime:      item.CreateTime.Format(utils.DateTime),
			UpdateTime:      item.UpdateTime.Format(utils.DateTime),
		})
	}

	/************** A6项目需要进度条 begin **************/
	var taskTaskInfoSrv services.TaskTaskInfo
	taskOverviewMap := taskTaskInfoSrv.GetOverviewByTaskIds(ctx, taskId)
	for k, item := range res.List {
		switch item.Status {
		case enums.TaskStatusBegin:
			res.List[k].Progress = 0
		case enums.TaskStatusFinish:
			res.List[k].Progress = 100
		default:
			if taskOverview, ok := taskOverviewMap[item.Id]; ok {
				if taskOverview.Progress == nil {
					res.List[k].Progress = 1
				} else {
					res.List[k].Progress = taskOverview.Progress.Value
				}
			} else {
				res.List[k].Progress = 1
			}
		}
	}
	return nil
}

// Overview 任务组 统计信息
func (t *TaskGroup) Overview(ctx context.Context, req *typespec.TaskGroupOverViewReq, res *typespec.TaskGroupOverViewResp) error {
	var srv services.TaskGroup
	group, err := srv.GetGroup(ctx, req.GroupId)
	if err != nil {
		return err
	}
	res.Overview = group.Overview
	return nil
}

// Status 任务组 状态信息
func (t *TaskGroup) Status(ctx context.Context, req *typespec.TaskGroupStatusReq, res *typespec.TaskGroupStatusResp) error {
	var srv services.TaskGroup
	group, err := srv.GetGroup(ctx, req.GroupId)
	if err != nil {
		return err
	}
	res.StatusNumber = group.Status
	res.Status = enums.GetTargetStatus(group.Status)
	return nil
}

// Edit 任务组编辑
func (t *TaskGroup) Edit(ctx context.Context, req *typespec.TaskGroupEditReq, resp *typespec.TaskGroupEditResp) error {
	var srv services.TaskGroup
	err := srv.Edit(ctx, req.Id, req.Name, req.Describe)
	if err != nil {
		return err
	}
	return nil
}
