package main

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"

	"smart/models/mysqls"
	"smart/services"
	"smart/tools/encryption"
	"smart/tools/enums"
	"smart/tools/utils"
)

func main() {
	var reportID int
	var outputPath string
	flag.IntVar(&reportID, "id", 0, "Report ID to export")
	flag.StringVar(&outputPath, "out", "vul_report.csv", "Output CSV file path")
	flag.Parse()

	if reportID == 0 {
		fmt.Println("Please provide a report ID using -id")
		os.Exit(1)
	}

	// Load Config
	configPath := "../../config.json"
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

	// Parse ConfigJSON to get Task ID(s)
	configJson := make(map[string]interface{})
	err = json.Unmarshal([]byte(reportRecode.ConfigJSON), &configJson)
	if err != nil {
		log.Error("Parse ConfigJSON err: " + err.Error())
		return
	}

	var taskIds []int
	objIdMid := configJson["objId"]
	if _, ok := objIdMid.(string); ok {
		if strings.Contains(objIdMid.(string), ",") {
			strs := strings.Split(objIdMid.(string), ",")
			for _, s := range strs {
				id, _ := strconv.Atoi(s)
				taskIds = append(taskIds, id)
			}
		} else {
			id, _ := strconv.ParseFloat(objIdMid.(string), 64)
			taskIds = append(taskIds, int(id))
		}
	} else {
		if val, ok := objIdMid.(float64); ok {
			taskIds = append(taskIds, int(val))
		}
	}

	if len(taskIds) == 0 {
		fmt.Println("No task IDs found in report config.")
		return
	}

	// Prepare CSV Writer
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()
	// Add BOM for Excel to open UTF-8 correctly
	file.Write([]byte("\xEF\xBB\xBF"))

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	header := []string{"Risk Level", "Vulnerability Type", "Vulnerability Name", "Target URL", "Poc Name", "Risk Score"}
	writer.Write(header)

	var taskVulSrv services.TaskVul
	aesEcb := encryption.AesEcb{}
	// Key for decryption (from reportGenerate.go)
	key := []byte("9876787656785679")

	riskEnumMap := enums.GetRiskEnumMap()
	typeEnumMap := enums.AllTypeEnumMap()

	for _, taskId := range taskIds {
		fmt.Printf("Processing Task ID: %d\n", taskId)
		// Get Vuls
		// Assuming enums.VulDataTypOne is correct as used in reportGenerate.go
		taskVuls := taskVulSrv.GetsByTaskId(ctx, taskId, 1) // 1 = VulDataTypOne

		for _, vul := range taskVuls {
			// Decrypt Name
			vulName := vul.Name
			if utils.IsHexString(vul.Name) {
				nameDecodeByte, _ := hex.DecodeString(vul.Name)
				decrypted := aesEcb.AesDecryptECB(nameDecodeByte, key)
				if len(decrypted) > 0 {
					// Check if decrypted string is empty or garbled?
					// Usually AesDecryptECB handles padding.
					vulName = string(decrypted)
					if vulName == "" {
						vulName = vul.Name
					}
				}
			}

			// Map Risk
			riskStr := riskEnumMap[vul.Risk]
			if riskStr == "" {
				riskStr = strconv.Itoa(vul.Risk)
			}

			// Map Type
			typeStr := typeEnumMap[vul.Type]
			if typeStr == "" {
				typeStr = strconv.Itoa(vul.Type)
			}

			// Write Row
			row := []string{
				riskStr,
				typeStr,
				vulName,
				vul.TargetUrl,
				vul.Pocname,
				strconv.Itoa(vul.Risk),
			}
			writer.Write(row)
		}
	}

	fmt.Printf("CSV export completed: %s\n", outputPath)
}
