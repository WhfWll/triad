package crons

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"smart/services"
)

// SystemRunMonitor 系统运行监控
func SystemRunMonitor() {
	ctx := context.Background()
	var (
		logSrv     services.TaskLog
		logInfoSrv services.TaskLogInfo
	)
	taskLogList, err := logSrv.GetFirstFiveTaskLog(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	logIdList := make([]int, 0)
	for _, log := range taskLogList {
		logIdList = append(logIdList, log.ID)
	}
	logInfoList, err := logInfoSrv.GetTaskLogInfoListByLogIds(ctx, logIdList)
	fmt.Println("retrieve the latest 5 log entries: ", len(logInfoList))
	if len(logInfoList) > 0 {
		return
	}
	// 结束smart进程和decision进程，让其重新启动
	log.Println("Detected that the latest 5 log entries are empty; the service may have crashed. Initiating and then terminating the 'smart' and 'decision' processes.")
	killSmartProcess(ctx)
}

// killProcess 结束进程
func killSmartProcess(ctx context.Context) error {
	cmd := exec.Command("pkill", "smart")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
	fmt.Println("pkill openvpn:\n ", string(out))
	return nil
}

// killProcess 结束进程
func killDecisionProcess(ctx context.Context) error {
	cmd := exec.Command("pkill", "decision")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
	fmt.Println("pkill openvpn:\n ", string(out))
	return nil
}
