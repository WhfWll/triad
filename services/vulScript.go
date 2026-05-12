package services

import (
	"context"
	"smart/models/mysqls"
	"smart/tools/enums"
)

type VulScripts struct{}

// GetVulIdsByIds 根据脚本id获取漏洞id
func (a *VulScripts) GetVulIdsByIds(ctx context.Context, vulIds []string) []mysqls.VulScripts {
	var script mysqls.VulScripts
	return script.GetScriptNamesByVulId(ctx, vulIds)
}

// GetVulScriptCount 获取脚本总数
func (a *VulScripts) GetVulScriptCount(ctx context.Context) (int64, error) {
	var script mysqls.VulLibraries
	return script.Count(ctx)
}

// 脚本 - 列表（分页）
func (a *VulScripts) OpenVulScriptList(ctx context.Context, page, size int, libIds []int, libName string, libRisks []int, libClasses []int, libTypes []int,
	scriptVerifyTypes []string, exploitImpact []int, operateSystem []int, status []int, ptOrder string) ([]int, []mysqls.VulLibraries, int64, error) {
	// 重置掉默认值
	if len(libIds) == 1 && libIds[0] == 0 {
		libIds = make([]int, 0)
	}
	if len(libRisks) == 1 && libRisks[0] == 0 {
		libRisks = make([]int, 0)
	}
	if len(libClasses) == 1 && libClasses[0] == 0 {
		libClasses = make([]int, 0)
	}
	if len(libTypes) == 1 && libTypes[0] == 0 {
		libTypes = make([]int, 0)
	}
	if len(operateSystem) == 1 && operateSystem[0] == 0 {
		operateSystem = make([]int, 0)
	}
	if len(exploitImpact) == 1 && exploitImpact[0] == 0 {
		exploitImpact = make([]int, 0)
	}
	if len(scriptVerifyTypes) == 1 && scriptVerifyTypes[0] == "" {
		scriptVerifyTypes = make([]string, 0)
	}
	if len(status) == 1 && status[0] == 0 {
		status = make([]int, 0)
	}

	// 查询验证类型
	var scriptModel mysqls.VulScripts
	vulIds := make([]string, 0)
	if len(scriptVerifyTypes) > 0 {
		vulIds = scriptModel.AllVulIdByVerifyType(ctx, scriptVerifyTypes)
		if len(vulIds) == 0 {
			return nil, nil, 0, nil
		}
	}

	var vulLib mysqls.VulLibraries
	//漏洞id列表
	vulIdList, err := vulLib.GetVulIdList(ctx, libIds, vulIds, libName, libTypes, libClasses, libRisks, exploitImpact, operateSystem, ptOrder, status, enums.VulVerifyTypePrincipleVerification)
	if err != nil {
		return nil, nil, 0, err
	}
	//漏洞列表页查询数据
	vulLibList, count, err := vulLib.GetVulLibList(ctx, page, size, enums.VulVerifyTypePrincipleVerification, libIds, vulIds, libName, libTypes, libClasses, libRisks, exploitImpact, operateSystem, ptOrder, status, "")
	if err != nil {
		return nil, nil, 0, err
	}
	return vulIdList, vulLibList, count, nil
}
