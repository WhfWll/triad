package services

import (
	"strings"

	"smart/models/mysqls"
	"smart/tools/enums"
)

// DataSecTargetScanMeta 单个目标扫描摘要（版本、基线统计、CVE 命中）
type DataSecTargetScanMeta struct {
	DbVersion     string
	BaselineTotal int
	BaselinePass  int
	BaselineFail  int
	BaselineError int
	CveMatchCount int
}

// SummarizeDBCheckResults 从检查结果汇总目标级扫描摘要
func SummarizeDBCheckResults(list []mysqls.DBCheckResult) DataSecTargetScanMeta {
	meta := DataSecTargetScanMeta{}
	for _, r := range list {
		if IsDatasecMetaCheckResult(r) {
			if r.RuleName == "数据库版本识别" {
				v := strings.TrimSpace(r.ActualValue)
				if v != "" && !strings.EqualFold(v, "unknown") {
					meta.DbVersion = v
				}
			}
			continue
		}
		if IsCveDBCheckResult(r) {
			if r.CheckResult == enums.BaselineCheckResultFail {
				meta.CveMatchCount++
			}
			continue
		}
		meta.BaselineTotal++
		switch r.CheckResult {
		case enums.BaselineCheckResultPass:
			meta.BaselinePass++
		case enums.BaselineCheckResultFail:
			meta.BaselineFail++
		default:
			meta.BaselineError++
		}
	}
	return meta
}

// IsDatasecMetaCheckResult 版本识别、CVE 汇总等元数据项
func IsDatasecMetaCheckResult(r mysqls.DBCheckResult) bool {
	return isVersionInfoResult(r) || r.RuleName == "CVE 版本匹配"
}
