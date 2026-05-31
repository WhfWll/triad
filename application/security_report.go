package application

import (
	"context"
	"fmt"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type SecurityReportApp struct{}

func (a *SecurityReportApp) Generate(ctx context.Context, req *typespec.SecurityReportGenerateReq, uid int) (*typespec.SecurityReportGenerateResp, error) {
	var title string
	var taskName string
	var htmlContent string

	taskID := req.TaskID.Int()

	switch req.Module {
	case "host":
		t, name, err := a.buildHostReport(ctx, taskID)
		if err != nil {
			return nil, err
		}
		title = name
		taskName = name
		htmlContent = t
	case "app":
		t, name, err := a.buildAppReport(ctx, taskID)
		if err != nil {
			return nil, err
		}
		title = name
		taskName = name
		htmlContent = t
	case "data":
		t, name, err := a.buildDataReport(ctx, taskID)
		if err != nil {
			return nil, err
		}
		title = name
		taskName = name
		htmlContent = t
	default:
		return nil, fmt.Errorf("unknown module: %s", req.Module)
	}

	model := mysqls.SecurityReport{
		Title:      title,
		Module:     req.Module,
		TaskID:     taskID,
		TaskName:   taskName,
		Content:    htmlContent,
		CreateBy:   uid,
		CreateTime: time.Now(),
	}
	if err := model.Add(ctx); err != nil {
		return nil, err
	}
	log.Infof("security report generated: id=%d module=%s taskId=%d", model.ID, req.Module, taskID)
	return &typespec.SecurityReportGenerateResp{ID: model.ID}, nil
}

func (a *SecurityReportApp) List(ctx context.Context, req *typespec.SecurityReportListReq) (*typespec.SecurityReportListResp, error) {
	var model mysqls.SecurityReport
	rows, total, err := model.List(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	list := make([]typespec.SecurityReportListItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, typespec.SecurityReportListItem{
			ID:         r.ID,
			Title:      r.Title,
			Module:     r.Module,
			ModuleName: moduleName(r.Module),
			TaskID:     r.TaskID,
			TaskName:   r.TaskName,
			CreateTime: r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return &typespec.SecurityReportListResp{List: list, Total: total}, nil
}

func (a *SecurityReportApp) Detail(ctx context.Context, req *typespec.SecurityReportDetailReq) (*typespec.SecurityReportDetailResp, error) {
	var model mysqls.SecurityReport
	row, err := model.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &typespec.SecurityReportDetailResp{
		ID:      row.ID,
		Title:   row.Title,
		Module:  row.Module,
		TaskID:  row.TaskID,
		Content: row.Content,
	}, nil
}

func (a *SecurityReportApp) Delete(ctx context.Context, req *typespec.SecurityReportDeleteReq) error {
	var model mysqls.SecurityReport
	return model.DeleteByID(ctx, req.ID)
}

func moduleName(m string) string {
	switch m {
	case "host":
		return "主机安全检查"
	case "app":
		return "应用安全检查"
	case "data":
		return "数据安全检查"
	}
	return m
}

// ---------------------------------------------------------------------------
// HTML report builders
// ---------------------------------------------------------------------------

func (a *SecurityReportApp) buildHostReport(ctx context.Context, taskID int) (html, taskName string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic building host report: %v", r)
		}
	}()

	var svc services.HostSecReportService
	meta := svc.GetHostReportMeta(ctx, taskID)
	if meta == nil {
		return "", "", fmt.Errorf("task %d not found", taskID)
	}

	var b strings.Builder
	b.WriteString(docHeader(meta.TaskName, "主机安全检查"))
	writeStatCards(&b, meta.StatCards)
	writeMetaSection(&b, [][2]string{
		{"任务批次", fmt.Sprintf("#%d", meta.TaskID)},
		{"任务类型", meta.TaskKind},
		{"检测目标数", fmt.Sprintf("%d 个", meta.TargetCount)},
		{"执行时间", meta.CheckTime},
	})
	writeTargetsTable(&b, meta.Targets)
	writeFindingsSection(&b, meta.Findings, meta.FindingColumns)
	b.WriteString(docFooter())
	return b.String(), meta.TaskName, nil
}

func (a *SecurityReportApp) buildAppReport(ctx context.Context, taskID int) (html, taskName string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic building app report: %v", r)
		}
	}()

	var svc services.HostSecReportService
	meta := svc.GetAppReportMeta(ctx, taskID)
	if meta == nil {
		return "", "", fmt.Errorf("task %d not found", taskID)
	}

	var b strings.Builder
	b.WriteString(docHeader(meta.TaskName, "应用安全检查"))
	writeStatCards(&b, meta.StatCards)
	writeMetaSection(&b, [][2]string{
		{"任务名称", meta.TaskName},
		{"任务 ID", fmt.Sprintf("#%d", meta.TaskID)},
		{"扫描类型", meta.TaskKind},
		{"扫描目标", meta.TargetSummary},
		{"扫描时间", meta.CheckTime},
	})
	writeVulnTable(&b, meta.Vulns)
	b.WriteString(docFooter())
	return b.String(), meta.TaskName, nil
}

func (a *SecurityReportApp) buildDataReport(ctx context.Context, taskID int) (html, taskName string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic building data report: %v", r)
		}
	}()

	var svc services.HostSecReportService
	meta := svc.GetDataReportMeta(ctx, taskID)
	if meta == nil {
		return "", "", fmt.Errorf("task %d not found", taskID)
	}

	var b strings.Builder
	b.WriteString(docHeader(meta.TaskName, meta.KindLabel))
	writeStatCards(&b, meta.StatCards)
	writeMetaSection(&b, [][2]string{
		{"任务名称", meta.TaskName},
		{"任务 ID", fmt.Sprintf("#%d", meta.TaskID)},
		{"任务类型", meta.KindLabel},
		{"扫描目标", meta.TargetSummary},
		{"检查时间", meta.CheckTime},
	})
	if len(meta.BaselineChecks) > 0 {
		writeSectionTitle(&b, "基线检查项")
		writeTable(&b, meta.BaselineColumns, meta.BaselineChecks)
	}
	if len(meta.SensitiveFindings) > 0 {
		writeSectionTitle(&b, "敏感字段发现")
		writeTable(&b, meta.SensitiveColumns, meta.SensitiveFindings)
	}
	b.WriteString(docFooter())
	return b.String(), meta.TaskName, nil
}

// ---------------------------------------------------------------------------
// HTML template helpers
// ---------------------------------------------------------------------------

func docHeader(title, moduleName string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;background:#1f2937;color:#e5e7eb;padding:32px 20px}
.report-wrap{width:min(1200px,100%%);margin:0 auto;background:#111827;border-radius:12px;box-shadow:0 4px 6px -1px rgba(0,0,0,0.3);padding:32px 40px}
.report-header{border-bottom:2px solid #00d4aa;padding-bottom:20px;margin-bottom:28px}
.report-header h1{font-size:22px;font-weight:700;color:#f9fafb}
.report-header .meta{font-size:13px;color:#9ca3af;margin-top:6px}
.report-header .meta span{margin-right:16px}
.section-title{font-size:16px;font-weight:600;color:#f9fafb;margin:28px 0 14px;padding-left:12px;border-left:3px solid #00d4aa}
.stat-row{display:flex;gap:12px;margin-bottom:20px;flex-wrap:wrap}
.stat-card{flex:1;min-width:100px;background:#1f2937;border:1px solid #374151;border-radius:8px;padding:16px 20px;text-align:center}
.stat-card .stat-value{font-size:26px;font-weight:700;line-height:1.3;color:#f9fafb}
.stat-card .stat-label{font-size:12px;color:#9ca3af;margin-top:4px}
.stat-card.critical .stat-value{color:#ef4444}
.stat-card.high .stat-value{color:#f97316}
.stat-card.medium .stat-value{color:#eab308}
.stat-card.low .stat-value{color:#22c55e}
.stat-card.pass .stat-value{color:#22c55e}
.stat-card.fail .stat-value{color:#ef4444}
.stat-card.rate .stat-value{color:#3b82f6}
.info-grid{display:grid;grid-template-columns:1fr 1fr;gap:12px;background:#1f2937;border:1px solid #374151;border-radius:8px;padding:16px 20px;margin-bottom:20px}
.info-item .label{font-size:12px;color:#9ca3af}
.info-item .value{font-size:14px;color:#e5e7eb;font-weight:500;word-break:break-all}
table{width:100%%;border-collapse:collapse;margin-bottom:20px;font-size:13px}
th{background:#1f2937;color:#9ca3af;font-weight:600;padding:10px 12px;text-align:left;border-bottom:2px solid #374151;white-space:normal;word-break:break-word}
td{padding:10px 12px;border-bottom:1px solid #374151;color:#d1d5db;vertical-align:top;white-space:normal;word-break:break-word;overflow-wrap:anywhere}
tr:hover td{background:#1f2937}
.report-table-findings-six{table-layout:fixed}
.report-table-findings-six th:nth-child(1),.report-table-findings-six td:nth-child(1){width:14%%}
.report-table-findings-six th:nth-child(2),.report-table-findings-six td:nth-child(2){width:42%%}
.report-table-findings-six th:nth-child(3),.report-table-findings-six td:nth-child(3){width:8%%;white-space:nowrap;word-break:keep-all;overflow-wrap:normal}
.report-table-findings-six th:nth-child(4),.report-table-findings-six td:nth-child(4){width:8%%;white-space:nowrap;word-break:keep-all;overflow-wrap:normal}
.report-table-findings-six th:nth-child(5),.report-table-findings-six td:nth-child(5){width:11%%;white-space:normal;word-break:break-word;overflow-wrap:anywhere}
.report-table-findings-six th:nth-child(6),.report-table-findings-six td:nth-child(6){width:17%%;white-space:normal;word-break:break-word;overflow-wrap:anywhere}
.report-table-targets td:nth-child(1),.report-table-targets th:nth-child(1){width:12%%}
.report-table-targets td:nth-child(2),.report-table-targets th:nth-child(2){width:50%%}
.report-table-targets td:nth-child(3),.report-table-targets th:nth-child(3){width:38%%}
.report-table-db-baseline{table-layout:fixed}
.report-table-db-baseline th:nth-child(1),.report-table-db-baseline td:nth-child(1){width:10%%;white-space:normal;word-break:break-word}
.report-table-db-baseline th:nth-child(2),.report-table-db-baseline td:nth-child(2){width:10%%;white-space:normal;word-break:break-word}
.report-table-db-baseline th:nth-child(3),.report-table-db-baseline td:nth-child(3){width:20%%;white-space:normal;word-break:break-word}
.report-table-db-baseline th:nth-child(4),.report-table-db-baseline td:nth-child(4){width:8%%;white-space:nowrap;word-break:keep-all}
.report-table-db-baseline th:nth-child(5),.report-table-db-baseline td:nth-child(5){width:8%%;white-space:nowrap;word-break:keep-all}
.report-table-db-baseline th:nth-child(6),.report-table-db-baseline td:nth-child(6){width:12%%;white-space:normal;word-break:break-word}
.report-table-db-baseline th:nth-child(7),.report-table-db-baseline td:nth-child(7){width:12%%;white-space:normal;word-break:break-word}
.report-table-db-baseline th:nth-child(8),.report-table-db-baseline td:nth-child(8){width:20%%;white-space:normal;word-break:break-word}
.report-table-sensitive-data{table-layout:fixed}
.report-table-sensitive-data th:nth-child(1),.report-table-sensitive-data td:nth-child(1){width:10%%;white-space:normal;word-break:break-word}
.report-table-sensitive-data th:nth-child(2),.report-table-sensitive-data td:nth-child(2){width:10%%;white-space:normal;word-break:break-word}
.report-table-sensitive-data th:nth-child(3),.report-table-sensitive-data td:nth-child(3){width:15%%;white-space:normal;word-break:break-word}
.report-table-sensitive-data th:nth-child(4),.report-table-sensitive-data td:nth-child(4){width:15%%;white-space:normal;word-break:break-word}
.report-table-sensitive-data th:nth-child(5),.report-table-sensitive-data td:nth-child(5){width:10%%;white-space:normal;word-break:break-word}
.report-table-sensitive-data th:nth-child(6),.report-table-sensitive-data td:nth-child(6){width:8%%;white-space:nowrap;word-break:keep-all}
.report-table-sensitive-data th:nth-child(7),.report-table-sensitive-data td:nth-child(7){width:15%%;white-space:normal;word-break:break-word}
.report-table-sensitive-data th:nth-child(8),.report-table-sensitive-data td:nth-child(8){width:17%%;white-space:normal;word-break:break-word}
@media (max-width: 960px){.report-wrap{padding:24px 20px}.report-table-findings-six,.report-table-db-baseline,.report-table-sensitive-data{table-layout:auto}.report-table-findings-six th,.report-table-findings-six td,.report-table-db-baseline th,.report-table-db-baseline td,.report-table-sensitive-data th,.report-table-sensitive-data td{white-space:normal;word-break:break-word;overflow-wrap:anywhere}}
.report-footer{text-align:center;font-size:12px;color:#6b7280;border-top:1px solid #374151;padding-top:16px;margin-top:28px}
.risk-critical{color:#ef4444;font-weight:600}
.risk-high{color:#f97316;font-weight:600}
.risk-medium{color:#eab308;font-weight:600}
.risk-low{color:#22c55e}
</style>
</head>
<body>
<div class="report-wrap">
<div class="report-header">
<h1>%s</h1>
<div class="meta"><span>模块: %s</span><span>生成时间: %s</span></div>
</div>`, title, title, moduleName, time.Now().Format("2006-01-02 15:04:05"))
}

func docFooter() string {
	return fmt.Sprintf(`<div class="report-footer">由安全检查系统自动生成 · %s</div></div></body></html>`, time.Now().Format("2006-01-02 15:04:05"))
}

func writeStatCards(b *strings.Builder, cards []services.ReportStatCard) {
	if len(cards) == 0 {
		return
	}
	b.WriteString(`<div class="stat-row">`)
	for _, c := range cards {
		cls := ""
		switch c.Class {
		case "critical":
			cls = "critical"
		case "high":
			cls = "high"
		case "medium":
			cls = "medium"
		case "low":
			cls = "low"
		case "pass":
			cls = "pass"
		case "fail":
			cls = "fail"
		case "rate":
			cls = "rate"
		}
		b.WriteString(fmt.Sprintf(`<div class="stat-card %s"><div class="stat-value">%s</div><div class="stat-label">%s</div></div>`, cls, c.Value, c.Label))
	}
	b.WriteString(`</div>`)
}

func writeMetaSection(b *strings.Builder, items [][2]string) {
	b.WriteString(`<div class="section-title">任务信息</div><div class="info-grid">`)
	for _, item := range items {
		b.WriteString(fmt.Sprintf(`<div class="info-item"><div class="label">%s</div><div class="value">%s</div></div>`, item[0], item[1]))
	}
	b.WriteString(`</div>`)
}

func writeSectionTitle(b *strings.Builder, title string) {
	b.WriteString(fmt.Sprintf(`<div class="section-title">%s</div>`, title))
}

func writeTargetsTable(b *strings.Builder, targets []services.ReportTableRow) {
	if len(targets) == 0 {
		return
	}
	writeSectionTitle(b, "检测目标")
	b.WriteString(`<table><thead><tr><th>序号</th><th>目标主机</th><th>操作系统</th></tr></thead><tbody>`)
	for i, t := range targets {
		cells := t.Cells
		ip := ""
		os := ""
		if len(cells) > 0 {
			ip = cells[0]
		}
		if len(cells) > 1 {
			os = cells[1]
		}
		b.WriteString(fmt.Sprintf(`<tr><td>%d</td><td>%s</td><td>%s</td></tr>`, i+1, ip, os))
	}
	b.WriteString(`</tbody></table>`)
}

func writeFindingsSection(b *strings.Builder, findings []services.ReportTableRow, columns []string) {
	if len(findings) == 0 {
		return
	}
	writeSectionTitle(b, "检测结果")
	writeTable(b, columns, findings)
}

func writeVulnTable(b *strings.Builder, vulns []services.ReportTableRow) {
	if len(vulns) == 0 {
		return
	}
	writeSectionTitle(b, "漏洞列表")
	b.WriteString(`<table><thead><tr><th>序号</th><th>漏洞名称</th><th>类型</th><th>风险</th><th>URL</th></tr></thead><tbody>`)
	for i, v := range vulns {
		cells := v.Cells
		name := ""
		typ := ""
		risk := ""
		url := ""
		if len(cells) > 0 {
			name = cells[0]
		}
		if len(cells) > 1 {
			typ = cells[1]
		}
		if len(cells) > 2 {
			risk = cells[2]
		}
		if len(cells) > 3 {
			url = cells[3]
		}
		riskCls := "risk-medium"
		switch risk {
		case "严重":
			riskCls = "risk-critical"
		case "高危":
			riskCls = "risk-high"
		case "低危":
			riskCls = "risk-low"
		}
		b.WriteString(fmt.Sprintf(`<tr><td>%d</td><td>%s</td><td>%s</td><td class="%s">%s</td><td>%s</td></tr>`, i+1, name, typ, riskCls, risk, url))
	}
	b.WriteString(`</tbody></table>`)
}

func writeTable(b *strings.Builder, headers []string, rows []services.ReportTableRow) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}
	tableClass := reportTableClass(headers)
	if tableClass == "" {
		b.WriteString(`<table><thead><tr>`)
	} else {
		b.WriteString(fmt.Sprintf(`<table class="%s"><thead><tr>`, tableClass))
	}
	for _, h := range headers {
		b.WriteString(fmt.Sprintf(`<th>%s</th>`, h))
	}
	b.WriteString(`</tr></thead><tbody>`)
	for _, row := range rows {
		b.WriteString(`<tr>`)
		for _, cell := range row.Cells {
			b.WriteString(fmt.Sprintf(`<td>%s</td>`, cell))
		}
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
}

func reportTableClass(headers []string) string {
	if len(headers) == 6 {
		return "report-table-findings-six"
	}
	if len(headers) == 8 {
		if headers[0] == "目标" && headers[1] == "检查类别" && headers[2] == "规则名称" {
			return "report-table-db-baseline"
		}
		if headers[0] == "目标" && headers[1] == "库名" && headers[2] == "表名" {
			return "report-table-sensitive-data"
		}
	}
	return ""
}
