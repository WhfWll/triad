package services

import (
	"context"
	"smart/models/mysqls"
)

type RiskManage struct{}

// VulAllList 所有漏洞测试列表及筛选
func (rm *RiskManage) VulAllList(ctx context.Context) ([]mysqls.TaskVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.GetAllTaskVulList(ctx)
	if err != nil {
		return nil, 0, err
	}
	var taskVulSrv TaskVul
	for i := range list {
		taskVulSrv.DecryptTaskVul(&list[i])
	}
	return list, count, nil
}

// VulListByIP 漏洞测试列表 通过ip筛选
func (rm *RiskManage) VulListByIP(ctx context.Context, ip string) ([]mysqls.TaskVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.GetTaskVulListByIP(ctx, ip)
	if err != nil {
		return nil, 0, err
	}
	var taskVulSrv TaskVul
	for i := range list {
		taskVulSrv.DecryptTaskVul(&list[i])
	}
	return list, count, nil
}

func (rm *RiskManage) GetTargetIdsRiskResultCount(ctx context.Context, taskId int, targetIds []int, risk int) ([]mysqls.TargetRiskResultCount, error) {
	var taskVul mysqls.TaskVul
	return taskVul.GetTargetIdsRiskResultCount(ctx, taskId, targetIds, risk)
}

// VulListByTargetID 漏洞测试列表 通过targetID筛选
func (rm *RiskManage) VulListByTargetID(ctx context.Context, targetID int) ([]mysqls.TaskVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.GetTaskVulListByTargetID(ctx, targetID)
	if err != nil {
		return nil, 0, err
	}
	var taskVulSrv TaskVul
	for i := range list {
		taskVulSrv.DecryptTaskVul(&list[i])
	}
	return list, count, nil
}

// VulInfoByTargetID 漏洞测试信息 通过targetID筛选
func (rm *RiskManage) VulInfoByTargetID(ctx context.Context, targetID int) (mysqls.TaskVul, int64, error) {
	var taskVul mysqls.TaskVul
	info, count, err := taskVul.GetTaskVulInfoByTargetID(ctx, targetID)
	if err != nil {
		return info, 0, err
	}
	var taskVulSrv TaskVul
	taskVulSrv.DecryptTaskVul(&info)
	return info, count, nil
}

// VulDeduplicationList 去重漏洞测试列表及筛选
func (rm *RiskManage) VulDeduplicationList(ctx context.Context, taskId, targetId int, vulType, risk int, search, ip string, page, limit, status int) ([]mysqls.DeduplicatedVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.GetTaskDeduplicationVulList(ctx, taskId, targetId, vulType, risk, search, ip, 1, page, limit, status)
	if err != nil {
		return nil, 0, err
	}
	var taskVulSrv TaskVul
	for i := range list {
		taskVulSrv.DecryptTaskVul(&list[i].TaskVul)
	}
	return list, count, nil
}

// GetTaskDeduplicationByLocationList 去重漏洞测试列表及筛选 按照漏洞名和地址去重
func (rm *RiskManage) GetTaskDeduplicationByLocationList(ctx context.Context, taskId, targetId, vulType, risk int, search, ip, location string, page, limit, status, verifyType, uid, role int) ([]mysqls.DeduplicatedVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.GetTaskDeduplicationByLocationList(ctx, taskId, targetId, vulType, risk, search, ip, location, 1, page, limit, status, verifyType, uid, role)
	if err != nil {
		return nil, 0, err
	}
	var taskVulSrv TaskVul
	for i := range list {
		taskVulSrv.DecryptTaskVul(&list[i].TaskVul)
	}
	return list, count, nil
}

// VulSimpleVulList 简单漏洞测试列表及筛选
func (rm *RiskManage) VulSimpleVulList(ctx context.Context, name, ip string, page, limit int) ([]mysqls.TaskVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.GetSimpleTaskVulList(ctx, name, ip, page, limit)
	if err != nil {
		return nil, 0, err
	}
	var taskVulSrv TaskVul
	for i := range list {
		taskVulSrv.DecryptTaskVul(&list[i])
	}
	return list, count, nil
}

// VulAllDeduplicationList 去全部重漏洞测试列表及筛选
func (rm *RiskManage) VulAllDeduplicationList(ctx context.Context, taskID int, uid, role int) ([]mysqls.DeduplicatedVul, int64, error) {
	var taskVul mysqls.TaskVul
	list, count, err := taskVul.GetAllTaskDeduplicationVulList(ctx, taskID, uid, role)
	if err != nil {
		return nil, 0, err
	}
	var taskVulSrv TaskVul
	for i := range list {
		taskVulSrv.DecryptTaskVul(&list[i].TaskVul)
	}
	return list, count, nil
}
