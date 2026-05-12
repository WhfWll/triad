package crons

import (
	"context"
	log "github.com/sirupsen/logrus"
	"smart/client/httpclients"
	"smart/services"
	"smart/tools/enums"
)

// VulScanReportUpdate 漏扫报告更新
func VulScanReportUpdate() {
	ctx := context.Background()
	// 报告分类
	var reportSrv services.Report
	reportRecodes := reportSrv.GetsByStatus(ctx, enums.ReportStatusRunning)
	for _, report := range reportRecodes {
		if report.Type != enums.ReportTypeVulScanTask && report.Type != enums.ReportTypeVulScanTarget {
			continue
		}
		req := httpclients.VulScanReportContentReq{
			Id: report.ID,
		}
		resp, err := httpclients.VulScanReportContent(ctx, req)
		if err != nil {
			continue
		}
		if resp.Data.Status != enums.ReportStatusFinish {
			continue
		}
		err = report.UpdateContent(ctx, report.ID, resp.Data.Content)
		if err != nil {
			log.Println(err)
		}
	}
}
