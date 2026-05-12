package crons

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	"io"
	"net/http"
	"regexp"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"
)

var (
	sendTaskUrl, getResultUrl string
)

func TripartiteBurpsuite() {
	ctx := context.Background()

	// 获取配置
	configMap := make(map[string]string)
	if err := config.Load("burpsuite", &configMap); err != nil {
		log.Error("使用异步，调用burpsuite 获取burpsuite的配置失败：" + err.Error())
		return
	}
	if sendTaskUrl = configMap["sendTaskUrl"]; sendTaskUrl == "" {
		log.Error("使用异步，调用burpsuite 获取burpsuite的sendTaskUrl配置为空")
		return
	}
	if getResultUrl = configMap["getResultUrl"]; getResultUrl == "" {
		log.Error("使用异步，调用burpsuite 获取burpsuite的getResultUrl配置为空")
		return
	}

	// 开启任务
	go TripartiteBurpsuiteStart(ctx)

	go TripartiteBurpsuiteGetResult(ctx)
}

// 开启任务
func TripartiteBurpsuiteStart(ctx context.Context) {
	// 获取所有待执行的任务
	var tripartiteSrc services.TripartiteTools
	burpsuiteWaits := tripartiteSrc.TripartiteToolsBurpsuiteGetsByStatus(ctx, enums.BurpsuiteStatusWait)

	for _, burpsuite := range burpsuiteWaits {
		targets := strings.Split(burpsuite.Target, ",")
		data := make(map[string]interface{})
		data["urls"] = targets
		bytesData, _ := json.Marshal(data)
		resp, _ := http.Post(sendTaskUrl, "application/json", bytes.NewReader(bytesData))
		if resp == nil {
			continue
		}
		defer resp.Body.Close()
		taskId := resp.Header.Get("Location")
		taskId = strings.TrimSpace(taskId)
		var taskOriginIdInt int
		var err error
		if taskId == "" {
			log.Error("使用异步，调用burpsuite 发送任务后转换原任务ID为空，burpsuite表ID=" + strconv.Itoa(burpsuite.ID))
		} else {
			taskOriginIdInt, err = strconv.Atoi(taskId)
			if err != nil {
				log.Error("使用异步，调用burpsuite 发送任务后转换原任务ID失败：" + taskId + " err：" + err.Error())
			}
		}

		// 更新状态为运行中 且更新原任务ID
		err = tripartiteSrc.TripartiteToolsBurpsuiteUpdateRunning(ctx, burpsuite.ID, taskOriginIdInt, enums.XrayStatusRunning)
		if err != nil {
			log.Error("使用异步，更新burpsuite状态失败：" + err.Error() + " burpsuite表ID为：" + strconv.Itoa(burpsuite.ID))
			continue
		}
	}
}

// 获取任务结果
func TripartiteBurpsuiteGetResult(ctx context.Context) {
	// 获取所有运行中的任务
	var tripartiteSrc services.TripartiteTools
	burpsuiteRunning := tripartiteSrc.TripartiteToolsBurpsuiteGetsByStatus(ctx, enums.BurpsuiteStatusRunning)
	// 获取所有运行中的任务结果
	burpsuiteIds := make([]int, 0, len(burpsuiteRunning))
	for _, item := range burpsuiteRunning {
		burpsuiteIds = append(burpsuiteIds, item.ID)
	}
	burpsuiteResults := tripartiteSrc.TripartiteToolsBurpsuiteResGetsHostAndPathByBurpsuiteId(ctx, burpsuiteIds)

	// 最终需要储存的结果数据
	finalResultInsertData := make([]mysqls.BurpsuiteTaskResult, 0)
	finalUpdateStatusAndRisk := make(map[int]map[string]int)
	for _, item := range burpsuiteRunning {
		resp, _ := http.Get(getResultUrl + strconv.Itoa(item.OriginTaskId))
		if resp == nil {
			continue
		}
		defer resp.Body.Close()
		dataResByte, _ := io.ReadAll(resp.Body)
		var dataStruct services.BurpsuiteResultData
		if err := json.Unmarshal(dataResByte, &dataStruct); err != nil {
			log.Error("使用异步，调用burpsuite 获取任务结果转换json失败：" + err.Error())
			continue
		}

		// 任务状态
		if _, ok := finalUpdateStatusAndRisk[item.ID]; !ok {
			finalUpdateStatusAndRisk[item.ID] = make(map[string]int)
		}
		if dataStruct.ScanStatus == "succeeded" {
			finalUpdateStatusAndRisk[item.ID]["status"] = enums.BurpsuiteStatusDone
		}

		// 已经存在的结果数据
		resultData := burpsuiteResults[item.ID]
		// 解析内部数据 [host+path:{}]
		//fmt.Println("解析数据：", dataStruct.IssueEvents, string(dataResByte))
		for _, originData := range dataStruct.IssueEvents {
			// 风险等级
			risk := 5
			switch strings.ToLower(originData.Issue.Severity) {
			case "high":
				risk = 2
			case "medium":
				risk = 3
			case "low":
				risk = 4
			case "information":
				risk = 5
			}
			if finalUpdateStatusAndRisk[item.ID]["risk"] == 0 || finalUpdateStatusAndRisk[item.ID]["risk"] > risk {
				finalUpdateStatusAndRisk[item.ID]["risk"] = risk
			}

			// 请求与响应信息
			requestResponse, _ := json.Marshal(originData.Issue.Evidence)
			// 原数据结果ID
			originResultId, _ := strconv.Atoi(originData.Id)

			// 如果数据不存在需要新增
			if _, ok := resultData[originResultId]; !ok {
				caption := originData.Issue.Caption
				insertionPoint := ""
				if strings.Contains(caption, "[") {
					re := regexp.MustCompile(`\[(.*?)\]`)
					match := re.FindStringSubmatch(caption)
					if len(match) > 1 {
						insertionPoint = match[1]
					}
				}

				finalResultInsertData = append(finalResultInsertData, mysqls.BurpsuiteTaskResult{
					BurpsuiteTaskID:       item.ID,
					OriginResultId:        originResultId,
					Action:                originData.Type,
					IssueType:             originData.Issue.Name,
					Host:                  originData.Issue.Origin,
					Path:                  originData.Issue.Path,
					InsertionPoint:        insertionPoint,
					Severity:              strconv.Itoa(risk),
					Confidence:            originData.Issue.Confidence,
					Describe:              originData.Issue.Description,
					IssueBackground:       originData.Issue.IssueBackground,
					RemediationBackground: originData.Issue.RemediationBackground,
					RequestResponse:       string(requestResponse),
					InternalData:          originData.Issue.InternalData,
					CreateTime:            time.Now(),
					UpdateTime:            time.Now(),
				})
			}
		}
	}

	// 事物开启 插入数据与修改数据
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	for burpsuiteId, value := range finalUpdateStatusAndRisk {
		if status, ok := value["status"]; ok {
			err := tripartiteSrc.TripartiteToolsBurpsuiteUpdateStatusById(dCtx, burpsuiteId, status)
			if err != nil {
				log.Error("使用异步，调用burpsuite 后插入结果失败 err=" + err.Error())
				return
			}
		}

		if risk, ok := value["risk"]; ok {
			err := tripartiteSrc.TripartiteToolsBurpsuiteUpdateRiskById(dCtx, burpsuiteId, risk)
			if err != nil {
				log.Error("使用异步，调用burpsuite 后插入结果失败 err=" + err.Error())
				return
			}
		}
	}
	if len(finalResultInsertData) > 0 {
		err := tripartiteSrc.TripartiteToolsBurpsuiteResultAdds(dCtx, finalResultInsertData)
		if err != nil {
			log.Error("使用异步，调用burpsuite 后插入结果失败 err=" + err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("使用异步，调用burpsuite 后插入结果失败 err=" + err.Error())
	}
}
