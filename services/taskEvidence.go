package services

import (
	"context"
	"smart/models/mysqls"
	"strings"
	"time"
)

type TaskEvidence struct {
}

// GetTaskEvidenceList 返回任务证据列表
func (te *TaskEvidence) GetTaskEvidenceList(ctx context.Context, taskID, targetId, page, limit, riskType int, search string) ([]mysqls.TaskEvidence, int64, error) {
	var taskEvidence = &mysqls.TaskEvidence{
		TaskID:   taskID,
		TargetID: targetId,
	}
	res, count, err := taskEvidence.GetTaskEvidenceList(ctx, taskID, targetId, page, limit, riskType, search)
	if err != nil {
		return res, 0, err
	}
	return res, count, nil
}

// DelTaskEvidences 批量删除任务证据 通过task_id
func (te *TaskEvidence) DelTaskEvidences(ctx context.Context, ids string) error {
	var taskEvidenceModel mysqls.TaskEvidence
	return taskEvidenceModel.DeleteTaskEvidenceIds(ctx, strings.Split(ids, ","))
}

// GetTaskEvidenceInfo 返回任务证据详情
func (te *TaskEvidence) GetTaskEvidenceInfo(ctx context.Context, id int) (mysqls.TaskEvidence, error) {
	var taskEvidence = &mysqls.TaskEvidence{
		ID: id,
	}
	res, err := taskEvidence.GetTaskEvidence(ctx)
	if err != nil {
		return mysqls.TaskEvidence{}, err
	}
	return res, nil
}

// AddTaskEvidenceInfo 添加任务证据信息
func (te *TaskEvidence) AddTaskEvidenceInfo(ctx context.Context, taskID, targetId, riskType int, targetURL, vulName, riskDetail, downloadedFiles string) error {
	var taskEvidence = &mysqls.TaskEvidence{
		Estate:          "valid",
		TaskID:          taskID,
		TargetID:        targetId,
		TargetURL:       targetURL,
		RiskType:        riskType,
		VulName:         vulName,
		RiskDetail:      riskDetail,
		DownloadedFiles: downloadedFiles,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	if err := taskEvidence.AddTaskEvidence(ctx); err != nil {
		return err
	}
	return nil
}

// GetVulEvidenceStat 获取漏洞取证统计
func (te *TaskEvidence) GetVulEvidenceStat(ctx context.Context, uid int, role int) []mysqls.VulEvidenceStat {
	var taskEvidence mysqls.TaskEvidence
	return taskEvidence.GetVulEvidenceStat(ctx, uid, role)
}
