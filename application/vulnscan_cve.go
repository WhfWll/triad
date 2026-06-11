package application

import (
	"context"
	"fmt"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
)

type VulnScanCveApp struct{}

type cveBatchTask struct {
	TaskID   int
	Req      *typespec.VulnScanCveBatchReq
	Status   string
	Progress int
	Total    int
	Results  []typespec.VulnScanCveTargetResult
	Errors   []string
	StartAt  time.Time
	Done     bool
	mu       sync.Mutex
}

var (
	cveBatchTasks   = make(map[int]*cveBatchTask)
	cveBatchTasksMu sync.Mutex
	cveTaskCounter  int
	cveBatchCancels = make(map[string]context.CancelFunc)
)

func cveBatchTargetKey(taskID int, targetIP string) string {
	return fmt.Sprintf("%d|%s", taskID, targetIP)
}

func setCveBatchCancel(taskID int, targetIP string, cancel context.CancelFunc) {
	cveBatchTasksMu.Lock()
	cveBatchCancels[cveBatchTargetKey(taskID, targetIP)] = cancel
	cveBatchTasksMu.Unlock()
}

func takeCveBatchCancel(taskID int, targetIP string) context.CancelFunc {
	key := cveBatchTargetKey(taskID, targetIP)
	cveBatchTasksMu.Lock()
	cancel := cveBatchCancels[key]
	delete(cveBatchCancels, key)
	cveBatchTasksMu.Unlock()
	return cancel
}

func updateHostCveBatchStopped(taskID int, targetIP string) {
	cveBatchTasksMu.Lock()
	bt := cveBatchTasks[taskID]
	cveBatchTasksMu.Unlock()
	if bt == nil {
		return
	}
	bt.mu.Lock()
	defer bt.mu.Unlock()
	bt.Results = append(bt.Results, typespec.VulnScanCveTargetResult{
		TargetIP: targetIP,
		Error:    "任务已手动结束",
	})
	bt.Progress++
	if bt.Progress >= bt.Total {
		bt.Status = "completed"
		bt.Done = true
	}
}

func (a *VulnScanCveApp) RunCveScan(ctx context.Context, req *typespec.VulnScanCveReq) (*typespec.VulnScanCveResp, error) {
	scanner := services.GetHostVulnScanner()
	if scanner == nil {
		return nil, fmt.Errorf("CVE scanner not available")
	}

	task := &services.VulnScanTask{
		TaskID:        req.TaskID,
		TargetID:      req.TargetID,
		Host:          req.Host,
		Port:          req.Port,
		Username:      req.Username,
		Password:      req.Password,
		Key:           req.Key,
		OSType:        req.OSType,
		Transport:     req.Transport,
		WinRMUseHttps: req.WinRMUseHttps,
	}

	report, err := scanner.RunVulnScan(ctx, task)
	if err != nil {
		return nil, err
	}

	dbCtx := mysql.NewContext(ctx, mysql.GetDB())
	if persistErr := services.PersistHostVulnScanResults(dbCtx, report, nil); persistErr != nil {
		log.Errorf("PersistHostVulnScanResults failed: %v", persistErr)
	}

	resp := &typespec.VulnScanCveResp{
		TaskID:       report.TaskID,
		TargetIP:     report.TargetIP,
		OSType:       report.OSType,
		OSTypeName:   enums.BaselineEnum.GetOSTypeName(report.OSType),
		Packages:     report.Packages,
		MatchedVulns: report.MatchedVulns,
		Critical:     report.CriticalCount,
		High:         report.HighCount,
		Medium:       report.MediumCount,
		Low:          report.LowCount,
		StartTime:    report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:      report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.VulnScanCveItem{
			PackageName:    r.PackageName,
			PackageVersion: r.PackageVersion,
			Cve:            r.Cve,
			Title:          r.Title,
			Severity:       r.Severity,
			RiskLevel:      r.RiskLevel,
		})
	}

	return resp, nil
}

func (a *VulnScanCveApp) RunBatchCveScan(ctx context.Context, req *typespec.VulnScanCveBatchReq) (*typespec.VulnScanCveBatchResp, error) {
	cveBatchTasksMu.Lock()
	cveTaskCounter++
	taskID := cveTaskCounter
	if req.TaskID > 0 {
		taskID = req.TaskID
	}
	bt := &cveBatchTask{
		TaskID:  taskID,
		Req:     req,
		Status:  "running",
		Total:   len(req.Targets),
		StartAt: time.Now(),
	}
	cveBatchTasks[taskID] = bt
	cveBatchTasksMu.Unlock()

	go a.runBatchCveScan(taskID, req)

	return &typespec.VulnScanCveBatchResp{
		TaskID: taskID,
		Status: "running",
	}, nil
}

func (a *VulnScanCveApp) runBatchCveScan(taskID int, req *typespec.VulnScanCveBatchReq) {
	cveBatchTasksMu.Lock()
	bt := cveBatchTasks[taskID]
	cveBatchTasksMu.Unlock()

	if bt == nil {
		return
	}

	scanner := services.GetHostVulnScanner()
	if scanner == nil {
		bt.mu.Lock()
		bt.Status = "failed"
		bt.Errors = append(bt.Errors, "CVE scanner not available")
		bt.Done = true
		bt.mu.Unlock()
		return
	}

	sem := make(chan struct{}, services.GetHostScanConcurrent(context.Background()))
	var wg sync.WaitGroup

	for i, target := range req.Targets {
		targetCtx, cancel := context.WithCancel(context.Background())
		setCveBatchCancel(taskID, target.Host, cancel)
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, t typespec.VulnScanCveTarget, runCtx context.Context, stop context.CancelFunc) {
			defer wg.Done()
			defer stop()
			defer takeCveBatchCancel(taskID, t.Host)
			defer func() { <-sem }()

			ctx, cancel := context.WithTimeout(runCtx, 5*time.Minute)
			defer cancel()

			task := &services.VulnScanTask{
				TaskID:        taskID,
				TargetID:      idx + 1,
				Host:          t.Host,
				Port:          t.Port,
				Username:      t.Username,
				Password:      t.Password,
				Key:           t.Key,
				OSType:        t.OSType,
				Transport:     t.Transport,
				WinRMUseHttps: t.WinRMUseHttps,
			}

			report, err := scanner.RunVulnScan(ctx, task)
			result := typespec.VulnScanCveTargetResult{
				TargetIP: t.Host,
				OSType:   t.OSType,
			}

			dbCtx := mysql.NewContext(context.Background(), mysql.GetDB())
			if err != nil {
				if ctx.Err() == context.Canceled {
					bt.mu.Lock()
					bt.Results = append(bt.Results, typespec.VulnScanCveTargetResult{
						TargetIP: t.Host,
						Error:    "任务已手动结束",
					})
					bt.Progress++
					bt.mu.Unlock()
					_ = services.PersistHostVulnScanResults(dbCtx, &services.VulnScanReport{
						TaskID:   taskID,
						TargetID: idx + 1,
						TargetIP: t.Host,
						OSType:   t.OSType,
						EndTime:  time.Now(),
					}, fmt.Errorf("任务已手动结束"))
					return
				}
				result.Error = err.Error()
				stub := &services.VulnScanReport{
					TaskID:   taskID,
					TargetID: idx + 1,
					TargetIP: t.Host,
					OSType:   t.OSType,
					EndTime:  time.Now(),
				}
				if persistErr := services.PersistHostVulnScanResults(dbCtx, stub, err); persistErr != nil {
					log.Errorf("PersistHostVulnScanResults failed for %s: %v", t.Host, persistErr)
				}
				bt.mu.Lock()
				bt.Errors = append(bt.Errors, fmt.Sprintf("%s: %v", t.Host, err))
				bt.Results = append(bt.Results, result)
				bt.Progress++
				bt.mu.Unlock()
				return
			}

			if persistErr := services.PersistHostVulnScanResults(dbCtx, report, nil); persistErr != nil {
				log.Errorf("PersistHostVulnScanResults failed for %s: %v", t.Host, persistErr)
			}

			result.Packages = report.Packages
			result.MatchedVulns = report.MatchedVulns
			result.Critical = report.CriticalCount
			result.High = report.HighCount
			result.Medium = report.MediumCount
			result.Low = report.LowCount
			result.StartTime = report.StartTime.Format("2006-01-02 15:04:05")
			result.EndTime = report.EndTime.Format("2006-01-02 15:04:05")

			for _, r := range report.Results {
				result.Results = append(result.Results, typespec.VulnScanCveItem{
					PackageName:    r.PackageName,
					PackageVersion: r.PackageVersion,
					Cve:            r.Cve,
					Title:          r.Title,
					Severity:       r.Severity,
					RiskLevel:      r.RiskLevel,
				})
			}

			bt.mu.Lock()
			bt.Results = append(bt.Results, result)
			bt.Progress++
			bt.mu.Unlock()

			log.Infof("CVE scan completed for %s: %d vulns found", t.Host, report.MatchedVulns)
		}(i, target, targetCtx, cancel)
	}

	wg.Wait()

	bt.mu.Lock()
	bt.Status = "completed"
	bt.Done = true
	bt.mu.Unlock()
}

func (a *VulnScanCveApp) GetBatchProgress(ctx context.Context, req *typespec.VulnScanCveProgressReq) (*typespec.VulnScanCveProgressResp, error) {
	cveBatchTasksMu.Lock()
	bt := cveBatchTasks[req.TaskID]
	cveBatchTasksMu.Unlock()

	if bt == nil {
		return nil, fmt.Errorf("task %d not found", req.TaskID)
	}

	bt.mu.Lock()
	defer bt.mu.Unlock()

	resp := &typespec.VulnScanCveProgressResp{
		TaskID:   bt.TaskID,
		Status:   bt.Status,
		Progress: bt.Progress,
		Total:    bt.Total,
		Results:  bt.Results,
		Errors:   bt.Errors,
	}

	if bt.Done {
		resp.Status = bt.Status
	}

	return resp, nil
}

func hostVulnRiskName(level int) string {
	if level <= 0 {
		return "—"
	}
	return enums.BaselineEnum.GetBaselineRiskName(level)
}

func hostVulnScanStatusName(status int) string {
	if status == mysqls.HostVulnScanStatusError {
		return "异常"
	}
	return "已完成"
}

func (a *VulnScanCveApp) GetTaskList(ctx context.Context, req *typespec.HostVulnTaskListReq) (*typespec.HostVulnTaskListResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 10
	}
	var model mysqls.HostVulnScan
	total, err := model.CountAll(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := model.ListPaged(ctx, page, size)
	if err != nil {
		return nil, err
	}
	resp := &typespec.HostVulnTaskListResp{Total: int(total)}
	for _, r := range rows {
		resp.List = append(resp.List, typespec.HostVulnTaskListItem{
			TaskID:         r.TaskID,
			TargetIP:       r.TargetIP,
			OSType:         r.OSType,
			OSTypeName:     enums.BaselineEnum.GetOSTypeName(r.OSType),
			Packages:       r.Packages,
			MatchedVulns:   r.MatchedVulns,
			Critical:       r.Critical,
			High:           r.High,
			Medium:         r.Medium,
			Low:            r.Low,
			WorstRiskLevel: r.WorstRiskLevel,
			WorstRiskName:  hostVulnRiskName(r.WorstRiskLevel),
			ScanStatus:     r.ScanStatus,
			ScanStatusName: hostVulnScanStatusName(r.ScanStatus),
			ErrorMessage:   r.ErrorMessage,
			CheckTime:      r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *VulnScanCveApp) GetTaskStat(ctx context.Context, req *typespec.HostVulnStatReq) (*typespec.HostVulnStatResp, error) {
	var model mysqls.HostVulnScan
	rows, err := model.ListByTaskID(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}
	resp := &typespec.HostVulnStatResp{TaskID: req.TaskID, TargetCount: len(rows)}
	for _, r := range rows {
		resp.Packages += r.Packages
		resp.MatchedVulns += r.MatchedVulns
		resp.Critical += r.Critical
		resp.High += r.High
		resp.Medium += r.Medium
		resp.Low += r.Low
	}
	return resp, nil
}

func (a *VulnScanCveApp) GetTaskTargets(ctx context.Context, taskID int) ([]typespec.HostVulnTaskListItem, error) {
	var model mysqls.HostVulnScan
	rows, err := model.ListByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	items := make([]typespec.HostVulnTaskListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, typespec.HostVulnTaskListItem{
			TaskID:         r.TaskID,
			TargetIP:       r.TargetIP,
			OSType:         r.OSType,
			OSTypeName:     enums.BaselineEnum.GetOSTypeName(r.OSType),
			Packages:       r.Packages,
			MatchedVulns:   r.MatchedVulns,
			Critical:       r.Critical,
			High:           r.High,
			Medium:         r.Medium,
			Low:            r.Low,
			WorstRiskLevel: r.WorstRiskLevel,
			WorstRiskName:  hostVulnRiskName(r.WorstRiskLevel),
			ScanStatus:     r.ScanStatus,
			ScanStatusName: hostVulnScanStatusName(r.ScanStatus),
			ErrorMessage:   r.ErrorMessage,
			CheckTime:      r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return items, nil
}

func (a *VulnScanCveApp) GetFindings(ctx context.Context, req *typespec.HostVulnFindingListReq) (*typespec.HostVulnFindingListResp, error) {
	var model mysqls.HostVulnFinding
	list, err := model.ListByTask(ctx, req.TaskID, req.TargetIP)
	if err != nil {
		return nil, err
	}
	resp := &typespec.HostVulnFindingListResp{Total: len(list)}
	for _, r := range list {
		resp.List = append(resp.List, typespec.HostVulnFindingItem{
			ID:             r.ID,
			TargetIP:       r.TargetIP,
			CveID:          r.CveID,
			Title:          r.Title,
			Severity:       r.Severity,
			RiskLevel:      r.RiskLevel,
			RiskName:       hostVulnRiskName(r.RiskLevel),
			PackageName:    r.PackageName,
			PackageVersion: r.PackageVersion,
			CheckTime:      r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}
