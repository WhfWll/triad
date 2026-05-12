package crons

import (
	"context"
	log "github.com/sirupsen/logrus"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"time"
)

// TaskReportVerify 定时检测是否有待执行的报告验证任务
func TaskReportVerify() {
	ctx := context.Background()
	var srv services.ReportVerify
	// 第一步 检查目标是否超过了并发限制
	runningTargets, err := srv.GetRunningTargets(ctx, 0)
	srv.HandleTimeoutTargets(ctx, runningTargets)
	if len(runningTargets) >= 10 {
		return
	}
	// 第二步 获取一个待开始的目标
	waitTarget, err := srv.GetOneWaitTarget(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	if waitTarget.ID == 0 {
		return
	}
	go execVerify(ctx, waitTarget, srv)
}

func execVerify(ctx context.Context, waitTarget mysqls.Reportverifytarget, srv services.ReportVerify) {
	// 第一步 进行任务和目标状态更新
	err := srv.UpdateTaskStatus(ctx, waitTarget.TaskId, enums.TaskStatusRunning)
	err = srv.UpdateTargetStatus(ctx, waitTarget.ID, enums.TaskStatusRunning)
	if err != nil {
		return
	}
	//第二步  进行端口扫描调用
	portResultList, rootUrlList, _ := srv.PortScan(ctx, waitTarget.Target, waitTarget.AnalysisData)
	err = srv.ReportVerifyHandlePort(ctx, waitTarget.TaskId, waitTarget.Target, portResultList)
	if len(rootUrlList) == 0 { //如果获取不到root_url的参数，就不再往下走了
		err = srv.UpdateTargetStatus(ctx, waitTarget.ID, enums.TaskStatusFinish) //更新目标和任务的状态
		targets, err := srv.GetWaitAndRunningTarget(ctx, waitTarget.TaskId)
		if err != nil {
			log.Error(err)
		}
		if len(targets) == 0 {
			err = srv.UpdateTaskStatus(ctx, waitTarget.TaskId, enums.TaskStatusFinish)
		}
		return
	}
	vulList, err := srv.GetReportVerifyVulListCanCall(ctx, waitTarget.Target) //获取可以进行验证的漏洞列表
	cveList := make([]string, 0)
	for _, vul := range vulList {
		if vul.Cve == "" {
			continue
		}
		cveList = append(cveList, vul.Cve)
	}
	// 第三步 进行报告验证任务测试
	paramList := make([]map[string]string, 0)
	for _, rootUrl := range rootUrlList {
		paramMap := make(map[string]string, 0)
		paramMap["key"] = "root_url"
		paramMap["value"] = rootUrl
		paramList = append(paramList, paramMap)
	}
	canVerifyList, err := srv.ReportVerifyTask(ctx, waitTarget.TaskId, waitTarget.ID, cveList, paramList)
	if err != nil {
		log.Error(err)
		return
	}
	log.Println("can verify list: ", canVerifyList)
	for {
		time.Sleep(10 * time.Second)
		end, resultMap, err := srv.ReportVerifyTaskResult(ctx, waitTarget.TaskId, waitTarget.ID)
		if err != nil {
			log.Error(err)
		}
		err = srv.ReportVerifyHandleResult(ctx, waitTarget.TaskId, waitTarget.Target, resultMap)
		if end {
			break
		}
	}

	// 第四步 统计任务信息
	taskStats, targetMap, err := srv.ReportVerifyOverView(ctx, waitTarget.TaskId)
	err = srv.UpdateTargetRisk(ctx, waitTarget.TaskId, targetMap)
	err = srv.UpdateTaskOverView(ctx, waitTarget.TaskId, taskStats)
	if err != nil {
		log.Error(err)
	}

	// 第四步 任务和目标状态更新完成
	err = srv.UpdateTargetStatus(ctx, waitTarget.ID, enums.TaskStatusFinish)
	targets, err := srv.GetWaitAndRunningTarget(ctx, waitTarget.TaskId)
	if len(targets) == 0 {
		err = srv.UpdateTaskStatus(ctx, waitTarget.TaskId, enums.TaskStatusFinish)
	}
	if err != nil {
		log.Error(err)
	}
}
