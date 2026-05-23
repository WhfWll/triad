package application

import (
	"context"
	"fmt"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/enums"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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
)

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
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, t typespec.VulnScanCveTarget) {
			defer wg.Done()
			defer func() { <-sem }()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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

			if err != nil {
				result.Error = err.Error()
				bt.mu.Lock()
				bt.Errors = append(bt.Errors, fmt.Sprintf("%s: %v", t.Host, err))
				bt.Results = append(bt.Results, result)
				bt.Progress++
				bt.mu.Unlock()
				return
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
		}(i, target)
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