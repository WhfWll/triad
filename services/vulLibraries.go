package services

import (
	"context"
	"fmt"

	"smart/models/mysqls"
	"smart/tools/enums"
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

// GetAppSecDefaultVulLibraries 应用安全未指定插件时：yak/nuclei（原理验证）+ universal（scriptList）
func (a *VulLibraries) GetAppSecDefaultVulLibraries(ctx context.Context, safeTest bool) ([]mysqls.VulLibraries, error) {
	all, err := (&mysqls.VulLibraries{}).AllVulLibraries(ctx, fmt.Sprintf("status = %d", enums.VulLibrariesStatusSucess))
	if err != nil {
		return nil, err
	}
	out := make([]mysqls.VulLibraries, 0)
	for _, v := range all {
		if v.Status != enums.VulLibrariesStatusSucess || v.Pocname == "" {
			continue
		}
		st := v.ScriptType
		if st != enums.VulScriptTypeYak && st != enums.VulScriptTypeNuclei && st != enums.VulScriptTypeUniversal {
			continue
		}
		if safeTest {
			if v.ExploitImpact == enums.VulScriptExploitImpactRefuseServer ||
				v.ExploitImpact == enums.VulScriptExploitImpactServiceBreakdown ||
				v.ExploitImpact == enums.VulScriptExploitImpactOutAge ||
				v.ExploitImpact == enums.VulScriptExploitImpactServerSlow {
				continue
			}
		}
		out = append(out, v)
	}
	return out, nil
}

// GetPrincipleScanVulLibraries 未指定插件 ID 时加载全部可用的 yak/nuclei 原理验证脚本
func (a *VulLibraries) GetPrincipleScanVulLibraries(ctx context.Context, safeTest bool) ([]mysqls.VulLibraries, error) {
	all, err := (&mysqls.VulLibraries{}).AllVulLibraries(ctx, fmt.Sprintf("status = %d", enums.VulLibrariesStatusSucess))
	if err != nil {
		return nil, err
	}
	out := make([]mysqls.VulLibraries, 0)
	for _, v := range all {
		if v.Status != enums.VulLibrariesStatusSucess {
			continue
		}
		if v.ScriptType != "yak" && v.ScriptType != "nuclei" {
			continue
		}
		if v.Pocname == "" {
			continue
		}
		if safeTest {
			if v.ExploitImpact == enums.VulScriptExploitImpactRefuseServer ||
				v.ExploitImpact == enums.VulScriptExploitImpactServiceBreakdown ||
				v.ExploitImpact == enums.VulScriptExploitImpactOutAge ||
				v.ExploitImpact == enums.VulScriptExploitImpactServerSlow {
				continue
			}
		}
		out = append(out, v)
	}
	return out, nil
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
