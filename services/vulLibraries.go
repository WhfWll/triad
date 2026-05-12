package services

import (
	"context"
	"smart/models/mysqls"
)

type VulLibraries struct{}

// GetVulLibByPocName 根据脚本名称获取漏洞和脚本数据
func (a *VulLibraries) GetVulLibByPocName(ctx context.Context, scriptName string) (mysqls.VulLibraries, error) {
	var lib mysqls.VulLibraries
	return lib.GetVulLibrariesByPocname(ctx, scriptName)
}

// GetVulLibsByIds 根据脚本id获取漏洞和脚本数据
func (a *VulLibraries) GetVulLibsByIds(ctx context.Context, vulIds []int) ([]mysqls.VulLibraries, error) {
	var lib mysqls.VulLibraries
	return lib.AllVulLibrariesForIds(ctx, vulIds)
}

// GetVulIdsByIds 根据脚本id获取漏洞id
func (a *VulLibraries) GetVulIdsByIds(ctx context.Context, vulIds []int) ([]string, error) {
	var lib mysqls.VulLibraries
	return lib.GetVulIdsByIds(ctx, vulIds)
}

// GetVulLibList 获取漏洞库分页列表
func (a *VulLibraries) GetVulLibList(ctx context.Context, page, limit, verifyType int, libIds []int, vulIds []string, search string,
	libTypes, libClasses, libRisks, exploitImpact, operateSystem []int, ptOrder string, status []int, code string) ([]mysqls.VulLibraries, int64, error) {
	var lib mysqls.VulLibraries
	return lib.GetVulLibList(ctx, page, limit, verifyType, libIds, vulIds, search, libTypes, libClasses, libRisks, exploitImpact, operateSystem, ptOrder, status, code)
}

// GetVulLibCount 漏洞库总数
func (a *VulLibraries) GetVulLibCount(ctx context.Context) (int64, error) {
	var lib mysqls.VulLibraries
	return lib.GetVulLibCount(ctx)
}

// UpdateVulLibraries 更新漏洞库
func (a *VulLibraries) UpdateVulLibraries(ctx context.Context, lib mysqls.VulLibraries) error {
	return lib.UpdateVulLibraries(ctx)
}

// UpdateStatus 更新漏洞状态
func (a *VulLibraries) UpdateStatus(ctx context.Context, id, status int) error {
	var lib mysqls.VulLibraries
	return lib.UpdateStatusById(ctx, id, status)
}
