package services

import (
	"context"
	"fmt"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strings"
)

// ReportStatCard is a single stat card in the report
type ReportStatCard struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Class string `json:"class"`
}

// ReportTableRow is a generic table row (cells are string columns)
type ReportTableRow struct {
	Cells []string `json:"cells"`
}

// HostReportMeta holds data for host security report
type HostReportMeta struct {
	TaskID         int
	TaskName       string
	TaskKind       string
	CheckTime      string
	TargetCount    int
	StatCards      []ReportStatCard
	Targets        []ReportTableRow
	Findings       []ReportTableRow
	FindingColumns []string
}

// AppReportMeta holds data for app security report
type AppReportMeta struct {
	TaskID        int
	TaskName      string
	TaskKind      string
	TargetSummary string
	CheckTime     string
	StatCards     []ReportStatCard
	Vulns         []ReportTableRow
	TargetCount   int
}

// DataReportMeta holds data for data security report
type DataReportMeta struct {
	TaskID            int
	TaskName          string
	KindLabel         string
	TargetSummary     string
	CheckTime         string
	StatCards         []ReportStatCard
	BaselineChecks    []ReportTableRow
	BaselineColumns   []string
	SensitiveFindings []ReportTableRow
	SensitiveColumns  []string
	TargetCount       int
}

// HostSecReportService provides report data for security check modules
type HostSecReportService struct{}

func (s *HostSecReportService) GetHostReportMeta(ctx context.Context, taskID int) *HostReportMeta {
	var baselineDao mysqls.BaselineCheckResult

	targets, err := baselineDao.GetTargetsByTaskID(ctx, taskID)
	if err == nil && len(targets) > 0 {
		return s.buildHostBaselineMeta(ctx, taskID, targets)
	}

	var malwareDao mysqls.HostMalwareScan
	malwareRows, err := malwareDao.ListByTaskID(ctx, taskID)
	if err == nil && len(malwareRows) > 0 {
		return s.buildHostMalwareMeta(ctx, taskID, malwareRows)
	}

	return nil
}

func (s *HostSecReportService) buildHostBaselineMeta(ctx context.Context, taskID int, targets []mysqls.BaselineTaskTargetRow) *HostReportMeta {
	meta := &HostReportMeta{
		TaskID:      taskID,
		TaskName:    fmt.Sprintf("安全配置核查 #%d", taskID),
		TaskKind:    "安全配置核查",
		CheckTime:   "-",
		TargetCount: len(targets),
	}

	var resultDao mysqls.BaselineCheckResult
	results, _ := resultDao.GetByTaskID(ctx, taskID)

	if len(targets) > 0 {
		if len(results) > 0 {
			meta.CheckTime = results[0].CreateTime.Format("2006-01-02 15:04:05")
		}
	}

	pass, fail, errCount := 0, 0, 0
	for _, r := range results {
		switch r.CheckResult {
		case 1:
			pass++
		case 2:
			fail++
		default:
			errCount++
		}
	}

	meta.StatCards = []ReportStatCard{
		{Label: "检查项总数", Value: fmt.Sprintf("%d", len(results)), Class: ""},
		{Label: "通过", Value: fmt.Sprintf("%d", pass), Class: "pass"},
		{Label: "不通过", Value: fmt.Sprintf("%d", fail), Class: "fail"},
		{Label: "异常", Value: fmt.Sprintf("%d", errCount), Class: "medium"},
	}
	if len(results) > 0 {
		rate := float64(pass) / float64(len(results)) * 100
		meta.StatCards = append(meta.StatCards, ReportStatCard{Label: "通过率", Value: fmt.Sprintf("%.1f%%", rate), Class: "rate"})
	}

	meta.Targets = make([]ReportTableRow, 0, len(targets))
	for _, t := range targets {
		osName := fmt.Sprintf("OS#%d", t.OSType)
		meta.Targets = append(meta.Targets, ReportTableRow{Cells: []string{t.TargetIP, osName}})
	}

	meta.FindingColumns = []string{"目标主机", "检查项", "结果", "风险", "期望值", "实际值"}
	for _, r := range results {
		resultName := "未知"
		switch r.CheckResult {
		case 1:
			resultName = "通过"
		case 2:
			resultName = "不通过"
		case 3:
			resultName = "异常"
		case 4:
			resultName = "不适配"
		}
		riskName := "低危"
		switch r.RuleRisk {
		case 1:
			riskName = "高危"
		case 2:
			riskName = "中危"
		case 3:
			riskName = "低危"
		}
		meta.Findings = append(meta.Findings, ReportTableRow{
			Cells: []string{r.TargetIP, r.RuleName, resultName, riskName, r.ExpectedValue, r.ActualValue},
		})
	}

	return meta
}

func (s *HostSecReportService) buildHostMalwareMeta(ctx context.Context, taskID int, rows []mysqls.HostMalwareScan) *HostReportMeta {
	meta := &HostReportMeta{
		TaskID:      taskID,
		TaskName:    fmt.Sprintf("恶意代码检测 #%d", taskID),
		TaskKind:    "恶意代码检测（YARA）",
		TargetCount: len(rows),
	}

	totalCritical, totalHigh, totalMedium, totalLow := 0, 0, 0, 0
	for _, r := range rows {
		totalCritical += r.Critical
		totalHigh += r.High
		totalMedium += r.Medium
		totalLow += r.Low
		if meta.CheckTime == "" || r.CreateTime.Format("2006-01-02 15:04:05") > meta.CheckTime {
			meta.CheckTime = r.CreateTime.Format("2006-01-02 15:04:05")
		}
	}
	if meta.CheckTime == "" {
		meta.CheckTime = "-"
	}

	meta.StatCards = []ReportStatCard{
		{Label: "目标总数", Value: fmt.Sprintf("%d", len(rows)), Class: ""},
		{Label: "严重", Value: fmt.Sprintf("%d", totalCritical), Class: "critical"},
		{Label: "高危", Value: fmt.Sprintf("%d", totalHigh), Class: "high"},
		{Label: "中危", Value: fmt.Sprintf("%d", totalMedium), Class: "medium"},
		{Label: "低危", Value: fmt.Sprintf("%d", totalLow), Class: "low"},
	}

	meta.Targets = make([]ReportTableRow, 0, len(rows))
	for _, r := range rows {
		osName := fmt.Sprintf("OS#%d", r.OSType)
		meta.Targets = append(meta.Targets, ReportTableRow{Cells: []string{r.TargetIP, osName}})
	}

	var checkResultDao mysqls.MalwareCheckResult
	results, _ := checkResultDao.GetByTaskID(ctx, taskID)

	meta.FindingColumns = []string{"目标主机", "匹配规则", "风险", "文件路径"}
	for _, r := range results {
		riskName := "低危"
		switch r.RiskLevel {
		case 1:
			riskName = "严重"
		case 2:
			riskName = "高危"
		case 3:
			riskName = "中危"
		}
		meta.Findings = append(meta.Findings, ReportTableRow{
			Cells: []string{r.TargetIP, r.MatchRule, riskName, r.FilePath},
		})
	}

	return meta
}

func (s *HostSecReportService) GetAppReportMeta(ctx context.Context, taskID int) *AppReportMeta {
	meta := &AppReportMeta{
		TaskID:        taskID,
		TaskKind:      "应用安全检查",
		TargetSummary: "-",
		CheckTime:     "-",
	}

	// 1. 任务基本信息
	var taskDao mysqls.TaskTask
	task, err := taskDao.GetTaskCheckTask(ctx, taskID)
	if err == nil {
		meta.TaskName = task.TaskName
		if !task.CreateTime.IsZero() {
			meta.CheckTime = task.CreateTime.Format("2006-01-02 15:04:05")
		}
	} else {
		meta.TaskName = fmt.Sprintf("应用安全检查 #%d", taskID)
	}

	// 2. 目标信息
	var targetDao mysqls.TaskTarget
	targets := targetDao.GetTargetsByTaskId(ctx, taskID)
	meta.TargetCount = len(targets)
	urls := make([]string, 0, len(targets))
	for _, t := range targets {
		if t.TargetURL != "" {
			urls = append(urls, t.TargetURL)
		}
	}
	if len(urls) > 0 {
		meta.TargetSummary = strings.Join(urls, ", ")
	}

	// 3. 漏洞统计数据
	var vulnDao mysqls.TaskVul
	stats := vulnDao.GetVulStatsByTaskId(ctx, taskID, enums.VulDataTypOne)
	critical, high, mid, low := 0, 0, 0, 0
	for _, s := range stats {
		switch s.Risk {
		case enums.VulLibrariesRiskDead:
			critical += s.Count
		case enums.VulLibrariesRiskHigh:
			high += s.Count
		case enums.VulLibrariesRiskMiddle:
			mid += s.Count
		case enums.VulLibrariesRiskLow:
			low += s.Count
		}
	}
	total := critical + high + mid + low
	meta.StatCards = []ReportStatCard{
		{Label: "严重", Value: fmt.Sprintf("%d", critical), Class: "critical"},
		{Label: "高危", Value: fmt.Sprintf("%d", high), Class: "high"},
		{Label: "中危", Value: fmt.Sprintf("%d", mid), Class: "medium"},
		{Label: "低危", Value: fmt.Sprintf("%d", low), Class: "low"},
		{Label: "漏洞合计", Value: fmt.Sprintf("%d", total), Class: ""},
	}

	// 4. 漏洞详情列表
	if total > 0 {
		vulns := vulnDao.GetsByTaskId(ctx, taskID, enums.VulDataTypOne)
		meta.Vulns = make([]ReportTableRow, 0, len(vulns))
		for _, v := range vulns {
			riskName := "低危"
			switch v.Risk {
			case enums.VulLibrariesRiskDead:
				riskName = "严重"
			case enums.VulLibrariesRiskHigh:
				riskName = "高危"
			case enums.VulLibrariesRiskMiddle:
				riskName = "中危"
			}
			typeName := enums.ToolsVulnerabilityEnum.GetTypeEnum(v.Type)
			url := v.Location
			if url == "" {
				url = v.VulAddress
			}
			if url == "" {
				url = v.TargetUrl
			}
			meta.Vulns = append(meta.Vulns, ReportTableRow{
				Cells: []string{v.Name, typeName, riskName, url},
			})
		}
	}

	return meta
}

func (s *HostSecReportService) GetDataReportMeta(ctx context.Context, taskID int) *DataReportMeta {
	meta := &DataReportMeta{
		TaskID:        taskID,
		TaskName:      fmt.Sprintf("数据安全检查 #%d", taskID),
		KindLabel:     "数据安全检查",
		TargetSummary: "-",
		CheckTime:     "-",
		StatCards: []ReportStatCard{
			{Label: "严重", Value: "0", Class: "critical"},
			{Label: "高危", Value: "0", Class: "high"},
			{Label: "中危", Value: "0", Class: "medium"},
			{Label: "低危", Value: "0", Class: "low"},
		},
	}
	return meta
}
