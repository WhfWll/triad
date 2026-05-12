package enums

const (
	//报告类型
	ReportTypeTask          = 1 //综述报告
	ReportTypeTarget        = 2 //目标报告
	ReportTypeLogicTask     = 3 //逻辑漏洞综述报告
	ReportTypeLogicTarget   = 4 //逻辑漏洞目标报告
	ReportTypeVulScanTask   = 5 //漏扫任务报告
	ReportTypeVulScanTarget = 6 //漏扫目标报告

	//报告状态
	ReportStatusWait    = 1 //待生成
	ReportStatusRunning = 2 //生成中
	ReportStatusFinish  = 3 //已生成
	//报告格式
	ReportFormatHtml  = 1 //html
	ReportFormatWord  = 2 //word
	ReportFormatPdf   = 3 //pdf
	ReportFormatExcel = 4 //excel
	ReportFormatCsv   = 5 //CSV
	//报告mapset

	// 报告输出方式
	MergeOutputType    = 1 // MergeOutputType 合并输出
	OneByOneOutputType = 2 // OneByOneOutputType 逐个输出
)

type ReportEnums struct {
}

func (r ReportEnums) ReportTypeEnum() map[int]string {
	enum := map[int]string{
		ReportTypeTask:        "综述报告",
		ReportTypeTarget:      "目标报告",
		ReportTypeLogicTask:   "逻辑漏洞综述报告",
		ReportTypeLogicTarget: "逻辑漏洞目标报告",
	}
	return enum
}

func (r ReportEnums) GetReportType(rtype int) string {
	enum := r.ReportTypeEnum()
	if res, ok := enum[rtype]; ok {
		return res
	}
	return ""
}

func (r ReportEnums) GetReportTypeEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: ReportTypeTask,
		Label: r.GetReportType(ReportTypeTask),
	}, {
		Value: ReportTypeTarget,
		Label: r.GetReportType(ReportTypeTarget),
	}}
	return result
}

func (r ReportEnums) ReportStatusEnum() map[int]string {
	enum := map[int]string{
		ReportStatusWait:    "待生成",
		ReportStatusRunning: "生成中",
		ReportStatusFinish:  "已生成",
	}
	return enum
}

func (r ReportEnums) GetReportStatus(status int) string {
	enum := r.ReportStatusEnum()
	if res, ok := enum[status]; ok {
		return res
	}
	return ""
}

func (r ReportEnums) GetReportStatusEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: ReportStatusWait,
		Label: r.GetReportStatus(ReportStatusWait),
	}, {
		Value: ReportStatusRunning,
		Label: r.GetReportStatus(ReportStatusRunning),
	}, {
		Value: ReportStatusFinish,
		Label: r.GetReportStatus(ReportStatusFinish),
	}}
	return result
}

func (r ReportEnums) ReportFormatEnum() map[int]string {
	enum := map[int]string{
		ReportFormatHtml:  "html",
		ReportFormatWord:  "word",
		ReportFormatPdf:   "pdf",
		ReportFormatExcel: "excel",
		ReportFormatCsv:   "CSV",
	}
	return enum
}

func (r ReportEnums) GetReportFormat(format int) string {
	enum := r.ReportFormatEnum()
	if res, ok := enum[format]; ok {
		return res
	}
	return ""
}

func (r ReportEnums) GetReportFormatEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: ReportFormatHtml,
		Label: r.GetReportFormat(ReportFormatHtml),
	}, {
		Value: ReportFormatWord,
		Label: r.GetReportFormat(ReportFormatWord),
	}, {
		Value: ReportFormatPdf,
		Label: r.GetReportFormat(ReportFormatPdf),
	}, {
		Value: ReportFormatExcel,
		Label: r.GetReportFormat(ReportFormatExcel),
	}, {
		Value: ReportFormatCsv,
		Label: r.GetReportFormat(ReportFormatCsv),
	},
	}
	return result
}
