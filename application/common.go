package application

import (
	"context"
	"errors"
	"smart/api/typespec"
	"smart/client/httpclients"
)

type Common struct {
}

// GlobalOptions 全局枚举接口
func (common Common) GlobalOptions(ctx context.Context, res *typespec.GlobalOptionsRes) error {
	// 获取漏洞服务枚举值
	decisionOptions, err := httpclients.GetDecisionOptions(ctx)
	if err != nil {
		return errors.New("decision服务通信错误：" + err.Error())
	}

	/****************************** 脚本信息 ***************************/
	// 脚本类型
	res.VulScriptType = decisionOptions.Data.VulScriptType

	// 验证类型
	// [value:label]
	res.VulScriptVerifyType = decisionOptions.Data.VulScriptVerifyType

	// 取证类型
	res.VulScriptEvidenceType = decisionOptions.Data.VulScriptEvidenceType

	// 影响力
	res.VulScriptExploitImpact = decisionOptions.Data.VulScriptExploitImpact

	/****************************** 漏洞信息 ***************************/
	// 漏洞分类
	res.VulLibrariesClass = decisionOptions.Data.VulLibrariesClass

	// 漏洞类型
	res.VulLibrariesType = decisionOptions.Data.VulLibrariesType

	// 风险等级
	res.VulLibrariesRisk = decisionOptions.Data.VulLibrariesRisk
	return nil
}
