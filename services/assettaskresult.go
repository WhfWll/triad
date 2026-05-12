package services

import (
	"context"
	"smart/models/mysqls"
)

type AssetTaskResult struct {
}

// GetAllAssetTaskResult 获取所有资产返回信息
func (atr *AssetTaskResult) GetAllAssetTaskResult(ctx context.Context, assetID int) ([]mysqls.Assettaskresult, error) {
	var assetTaskResultModel mysqls.Assettaskresult
	taskResultList := assetTaskResultModel.GetAllTaskPortResultByAssetId(ctx, assetID)
	return taskResultList, nil
}

// GetAssetTaskResult 获取资产返回信息
func (atr *AssetTaskResult) GetAssetTaskResult(ctx context.Context, assetID, page, size int) ([]mysqls.Assettaskresult, error) {
	var assetTaskResultModel mysqls.Assettaskresult
	taskResultList := assetTaskResultModel.GetTaskPortResultByAssetId(ctx, assetID, page, size)
	return taskResultList, nil
}

// GetAssetTaskOpenPort 获取所有资产返回端口信息，并去重
func (atr *AssetTaskResult) GetAssetTaskOpenPort(ctx context.Context, assetID int) string {
	var (
		assetOpenPort        string
		assetTaskResultModel mysqls.Assettaskresult
	)
	uniquePorts := make(map[string]struct{})
	taskResultList := assetTaskResultModel.GetAllTaskPortResultByAssetId(ctx, assetID)
	for _, v := range taskResultList {
		uniquePorts[v.Field2] = struct{}{}
	}
	for port := range uniquePorts {
		assetOpenPort += port + ","
	}
	if len(assetOpenPort) > 0 {
		assetOpenPort = assetOpenPort[:len(assetOpenPort)-1]
	}
	return assetOpenPort
}
