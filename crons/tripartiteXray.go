package crons

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"time"
)

func TripartiteXray() {
	ctx := context.Background()

	// 获取所有待执行的任务
	var tripartiteSrc services.TripartiteTools
	xrayWaits := tripartiteSrc.TripartiteToolsXrayGetsByStatus(ctx, enums.XrayStatusWait)

	for _, xray := range xrayWaits {
		// 更新状态为运行中
		err := tripartiteSrc.TripartiteToolsXrayUpdateStatusById(ctx, xray.ID, enums.XrayStatusRunning)
		if err != nil {
			log.Error("使用异步，更新xray状态失败：" + err.Error() + " xray表ID为：" + strconv.Itoa(xray.ID))
			continue
		}

		go TripartiteXrayProcess(ctx, tripartiteSrc, xray)
	}

}

func TripartiteXrayProcess(ctx context.Context, tripartiteSrc services.TripartiteTools, xray mysqls.XrayTask) {
	riskNum := 0
	defer func() {
		// 更新xray状态为已完成
		err := tripartiteSrc.TripartiteToolsXrayFinish(ctx, xray.ID, riskNum)
		if err != nil {
			log.Error("使用异步，修改xray finish失败：" + err.Error() + " xray表ID为：" + strconv.Itoa(xray.ID))
			return
		}
	}()

	// 调用xray服务
	content, err := utils.XRayExec(ctx, xray.Target, xray.IsCrawler)
	if err != nil {
		log.Error("使用异步，调用xray服务失败：" + err.Error() + " xray表ID为：" + strconv.Itoa(xray.ID))
		return
	}

	// 将结果储存入xray的结果表中
	xrayResultModelDatas := make([]services.TripartiteToolsXrayCreateResultItem, 0)
	for _, item := range content {
		requestInfo, err := json.Marshal(item.Detail.Snapshot)
		if err != nil {
			log.Error("使用异步，调用xray服务，解析请求信息失败：" + err.Error() + " 目标为：" + xray.Target)
			continue
		}

		extraInfo, err := json.Marshal(item.Detail.Extra)
		if err != nil {
			log.Error("使用异步，调用xray服务，解析扩展信息失败：" + err.Error() + " 目标为：" + xray.Target)
			continue
		}
		createTime := time.Unix(item.CreateTime/1000, 0)

		xrayResultModelDatas = append(xrayResultModelDatas, services.TripartiteToolsXrayCreateResultItem{
			XrayTaskId:  xray.ID,
			Addr:        item.Detail.Addr,
			Payload:     item.Detail.Payload,
			RequestInfo: string(requestInfo),
			Extra:       string(extraInfo),
			Plugin:      item.Plugin,
			CreateTime:  createTime,
			UpdateTime:  createTime,
		})
	}
	// 结果不为0则插入
	if len(xrayResultModelDatas) > 0 {
		err = tripartiteSrc.TripartiteToolsXrayResultCreate(ctx, xrayResultModelDatas)
		if err != nil {
			log.Error("使用异步，插入xray调用结果失败：" + err.Error() + " xray表ID为：" + strconv.Itoa(xray.ID))
			return
		}
	}

	riskNum = len(xrayResultModelDatas)
}
