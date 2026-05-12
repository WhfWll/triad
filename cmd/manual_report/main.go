package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"

	"smart/api/typespec"
	"smart/crons"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
)

func main() {
	var reportID int
	flag.IntVar(&reportID, "id", 0, "Report ID to generate")
	flag.Parse()

	if reportID == 0 {
		fmt.Println("Please provide a report ID using -id")
		os.Exit(1)
	}

	// Load Config
	// Try loading from ../../config.json (assuming running from cmd/manual_report/)
	configPath := "../../config.json"
	// Check if file exists, if not try default config.json in current directory
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config.json"
	}

	fmt.Printf("Loading config from: %s\n", configPath)
	err := config.NewConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	mysql.Setup()

	ctx := context.Background()

	// Get Report Record
	db := mysql.FromContext(ctx).Model(&mysqls.Reportrecord{})
	var reportRecode mysqls.Reportrecord
	if err := db.Where("id = ?", reportID).First(&reportRecode).Error; err != nil {
		log.Errorf("Failed to find report with ID %d: %v", reportID, err)
		return
	}

	fmt.Printf("Found report: %s (Type: %d)\n", reportRecode.Name, reportRecode.Type)

	if reportRecode.Type != enums.ReportTypeTask {
		fmt.Printf("Report type is %d. This tool currently only supports ReportTypeTask (%d).\n", reportRecode.Type, enums.ReportTypeTask)
		return
	}

	var reportSrv services.Report
	taskCate, _, err := reportSrv.GetReportContentEnum(ctx)
	if err != nil {
		log.Error("GetReportContentEnum err: " + err.Error())
		return
	}

	// Parse ConfigJSON
	configJson := make(map[string]interface{})
	err = json.Unmarshal([]byte(reportRecode.ConfigJSON), &configJson)
	if err != nil {
		log.Error("Parse ConfigJSON err: " + err.Error())
		return
	}

	var (
		objId  float64
		objIDs string
	)
	objIdMid := configJson["objId"]
	if _, ok := objIdMid.(string); ok {
		if !strings.Contains(objIdMid.(string), ",") {
			objId, _ = strconv.ParseFloat(objIdMid.(string), 64)
		} else {
			objIDs = objIdMid.(string)
		}
	} else {
		if val, ok := objIdMid.(float64); ok {
			objId = val
		}
	}

	configJsonContentMap := make(map[string]int)
	if configJsonContent, ok := configJson["content"]; ok {
		configJsonContentByte, err := json.Marshal(configJsonContent)
		if err != nil {
			log.Error("Marshal config content err: " + err.Error())
			return
		}
		err = json.Unmarshal(configJsonContentByte, &configJsonContentMap)
		if err != nil {
			log.Error("Unmarshal config content map err: " + err.Error())
			return
		}
	}

	reportRecodeId := strconv.Itoa(reportRecode.ID)
	fmt.Println("Generating report content...")

	// Logic from crons/reportGenerate.go
	if objIDs != "" {
		ids := strings.Split(objIDs, ",")
		var list []typespec.ReportTaskContent
		for _, v := range ids {
			taskId, _ := strconv.Atoi(v)
			info, err := crons.ExecReportTypeTask(ctx, taskId, reportRecodeId, reportRecode, taskCate, configJsonContentMap)
			if err != nil {
				log.Error("ExecReportTypeTask error for task " + v + ": " + err.Error())
				continue
			}
			list = append(list, info)
		}
		taskContentByte, _ := json.Marshal(list)
		if err = reportSrv.ReportUpdateContent(ctx, reportRecode.ID, string(taskContentByte)); err != nil {
			log.Error("ReportUpdateContent err: " + err.Error())
			return
		}
	} else {
		taskContent, err := crons.ExecReportTypeTask(ctx, int(objId), reportRecodeId, reportRecode, taskCate, configJsonContentMap)
		if err != nil {
			log.Error("ExecReportTypeTask err: " + err.Error())
			return
		}
		taskContentByte, _ := json.Marshal(taskContent)
		fmt.Printf("Generated JSON Content:\n%s\n", string(taskContentByte))
		fmt.Printf("Generated JSON Size: %.2f KB\n", float64(len(taskContentByte))/1024)
		if err = reportSrv.ReportUpdateContent(ctx, reportRecode.ID, string(taskContentByte)); err != nil {
			log.Error("ReportUpdateContent err: " + err.Error())
			return
		}
	}

	// Update Report Status to Finish
	if err := reportSrv.UpdateStatus(ctx, reportRecode.ID, enums.ReportStatusFinish); err != nil {
		log.Error("UpdateStatus err: " + err.Error())
		return
	}

	fmt.Println("Report generation completed successfully.")
}
