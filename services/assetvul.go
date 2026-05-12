package services

import (
	"context"
	"smart/models/mysqls"
)

type AssetVul struct {
}

// AllAssetVulByAssetIds 获取对应资产id下的所有漏洞信息
func (ag *AssetVul) AllAssetVulByAssetIds(ctx context.Context, assetIds []int) []mysqls.Assetvul {
	assetVulModel := mysqls.Assetvul{}
	return assetVulModel.AllAssetVulByAssetIds(ctx, assetIds)
}

// AllAssetVul 获取所有资产信息
func (ag *AssetVul) AllAssetVul(ctx context.Context, assetIds []int, search string, page, size int) ([]mysqls.Assetvul, int64) {
	assetVulModel := mysqls.Assetvul{}
	return assetVulModel.AllAssetVulList(ctx, assetIds, search, page, size)
}

// DeleteAssetVulByIds 通过漏洞ID删除
func (ag *AssetVul) DeleteAssetVulByIds(ctx context.Context, vulIds []int) error {
	var assetModel mysqls.Assetvul
	return assetModel.DeleteByVulId(ctx, vulIds)
}

// GetAssetVulByAssetIds 获取对应资产id下的漏洞信息
func (ag *AssetVul) GetAssetVulByAssetIds(ctx context.Context, assetIds []int, page, size int) ([]mysqls.Assetvul, int64) {
	assetVulModel := mysqls.Assetvul{}
	return assetVulModel.GetAssetVulByAssetIds(ctx, assetIds, page, size)
}

// UpdateAssetVulStatus 更新资产漏洞状态
func (ag *AssetVul) UpdateAssetVulStatus(ctx context.Context, vulID, status int) error {
	var assetVulModel mysqls.Assetvul
	return assetVulModel.UpdateAssetVulStatus(ctx, vulID, status)
}
