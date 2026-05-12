package crons

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"smart/services"
	"smart/tools/enums"
	tooltime "smart/tools/time"
	"strings"
	"time"
)

func BasExec() {
	var (
		basSrv       services.Bas
		ctx          = context.Background()
		basTargetIds = make([]int, 0)
	)
	//查询任务数据
	basTaskIds, err := basSrv.GetBasTaskIdsByStatus(ctx, enums.BasTaskStatusWait)
	if err != nil {
		log.Info("BasExec GetBasTaskByStatus err:", err)
		return
	}
	if len(basTaskIds) == 0 {
		return
	}
	//查询目标数据
	basTargetRes := basSrv.GetBasTargetByTaskIdsAndStatus(ctx, basTaskIds, enums.BasTaskStatusWait)
	if len(basTargetRes) == 0 {
		return
	}
	for i := 0; i < len(basTargetRes); i++ {
		basTargetIds = append(basTargetIds, basTargetRes[i].ID)
		var templateData []int
		err = json.Unmarshal([]byte(basTargetRes[i].BasTemplateJSON), &templateData)
		if err != nil || len(templateData) == 0 {
			log.Info("BasExec BasTemplateJSON json Unmarshal err or ruleIds is empty...")
			continue
		}
		basRuleRes := basSrv.GetBasRuleGetByIds(ctx, templateData) //查询bas规则数据
		if len(basRuleRes) == 0 {
			log.Info("BasExec GetBasRuleGetByIds not find rule data")
			continue
		}
		//生成临时文件
		tmpFilePath, err := basSrv.SaveBasRuleFilePath(basRuleRes)
		if err != nil {
			log.Info("BasExec SaveBasRuleFilePath err:", err)
			continue
		}
		//请求发送任务
		go basSrv.SendBasTask(ctx, basTargetRes[i].BasTaskID, basTargetRes[i].ID, basTargetRes[i].Addr, tmpFilePath, sendBasTaskHandle)
		//存入日志
		logData := fmt.Sprintf("发送bas任务成功,发送目标:%s,发送规则id为:%s", basTargetRes[i].Addr, basTargetRes[i].BasTemplateJSON)
		err = basSrv.AddBasLog(ctx, basTargetRes[i].BasTaskID, basTargetRes[i].ID, logData)
		if err != nil {
			log.Info("SendBasTask SendBasTaskHandle AddBasLog err:", err)
		}
	}
	//更改状态
	err = basSrv.UpdateBasTaskTargetStatus(ctx, basTaskIds, basTargetIds, enums.BasTaskStatusRunning)
	if err != nil {
		log.Info("BasExec UpdateBasTaskTargetStatus err:", err)
	}
	time.Sleep(180 * time.Second)
	err = basSrv.UpdateBasTaskTargetStatus(ctx, basTaskIds, basTargetIds, enums.BasTaskStatusDone)
	if err != nil {
		log.Info("BasExec UpdateBasTaskTargetStatus err:", err)
	}
}

//任务回调处理函数
func sendBasTaskHandle(ctx context.Context, basTaskId int, basTargetId int, ip string, data string) {
	var basSrv services.Bas
	//存入日志
	err := basSrv.AddBasLog(ctx, basTaskId, basTargetId, data)
	if err != nil {
		log.Info("SendBasTask SendBasTaskHandle AddBasLog err:", err)
	}
	//存入漏洞测试
	if strings.Contains(data, "\"type\":\"bas\"") {
		var (
			tmpData services.BaSendMsg
			ruleIds = make([]int, 0)
		)
		err = json.Unmarshal([]byte(data), &tmpData)
		if err != nil {
			log.Info("SendBasTask json Unmarshal err:", err)
			return
		}
		for i := 0; i < len(tmpData.Content); i++ {
			ruleIds = append(ruleIds, tmpData.Content[i].RuleId)
		}
		basRuleRes := basSrv.GetBasRulesByIds(ctx, ruleIds) //查询bas规则
		basVulRes := basSrv.OrderBasSendData(basTaskId, basTargetId, ip, tmpData.Content, basRuleRes)
		err = basSrv.SaveBasVul(ctx, basVulRes)
		if err != nil {
			log.Info("SendBasTaskSaveBasVul err:", err)
		}
	}
}

func BasNodeOnlineStatus() {
	var (
		basSrv     services.Bas
		ctx        = context.Background()
		beforeTime = tooltime.GetBeforeSecondTime(30) //30秒没有更新的数据下线处理
	)

	//根据上下状态和时间查询basnode数据
	basNodeIds := basSrv.GetBasAgentByOnlineStatusAndUpdatetime(ctx, enums.BasNodeOnlineStatusOnline, beforeTime)
	//更新上下线状态
	if len(basNodeIds) > 0 {
		var params = map[string]interface{}{"online_status": enums.BasNodeOnlineStatusOffline}
		err := basSrv.UpdateBasNodeByIds(ctx, basNodeIds, params)
		if err != nil {
			log.Info("BasNodeOnlineStatus UpdateBasNodeByIds err:", err)
		}
	}
}
