package application

import (
	"context"
	"fmt"

	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
)

// DeleteHostSecTasks 删除主机安全检查任务记录（列表行级：taskId + targetIp + 类型）
func (a *BaselineApp) DeleteHostSecTasks(ctx context.Context, req *typespec.HostSecTaskDeleteReq) (*typespec.HostSecTaskDeleteResp, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items is empty")
	}

	deleted := 0
	seen := make(map[string]bool, len(req.Items))

	for _, item := range req.Items {
		if item.TaskID <= 0 {
			continue
		}
		key := fmt.Sprintf("%s|%d|%s|%d", item.Source, item.TaskID, item.TargetIP, item.ScanScene)
		if seen[key] {
			continue
		}
		seen[key] = true

		switch item.Source {
		case "malware":
			if item.TargetIP == "" {
				return nil, fmt.Errorf("targetIp required for malware delete")
			}
			markYaraBatchTargetDeleted(item.TaskID, item.TargetIP)
			var scanModel mysqls.HostMalwareScan
			if err := scanModel.DeleteByTaskTarget(ctx, item.TaskID, item.TargetIP); err != nil {
				return nil, err
			}
			var model mysqls.MalwareCheckResult
			if err := model.DeleteByTaskIDAndTargetIP(ctx, item.TaskID, item.TargetIP); err != nil {
				return nil, err
			}
		case "vuln":
			if item.TargetIP == "" {
				return nil, fmt.Errorf("targetIp required for vuln delete")
			}
			var scanModel mysqls.HostVulnScan
			if err := scanModel.DeleteByTaskTarget(ctx, item.TaskID, item.TargetIP); err != nil {
				return nil, err
			}
			var findingModel mysqls.HostVulnFinding
			if err := findingModel.DeleteByTaskTarget(ctx, item.TaskID, item.TargetIP); err != nil {
				return nil, err
			}
		case "baseline":
			scene := item.ScanScene
			if scene <= 0 {
				scene = enums.HostScanSceneBaseline
			}
			if item.TargetIP == "" {
				return nil, fmt.Errorf("targetIp required for baseline delete")
			}
			var model mysqls.BaselineCheckResult
			if err := model.DeleteByTaskTargetAndScene(ctx, item.TaskID, item.TargetIP, scene); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported source: %s", item.Source)
		}
		deleted++
	}

	return &typespec.HostSecTaskDeleteResp{Deleted: deleted}, nil
}

func (a *BaselineApp) StopHostSecTasks(ctx context.Context, req *typespec.HostSecTaskStopReq) (*typespec.HostSecTaskStopResp, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items is empty")
	}

	stopped := 0
	seen := make(map[string]bool, len(req.Items))
	for _, item := range req.Items {
		if item.TaskID <= 0 || item.TargetIP == "" {
			continue
		}
		key := fmt.Sprintf("%s|%d|%s|%d", item.Source, item.TaskID, item.TargetIP, item.ScanScene)
		if seen[key] {
			continue
		}
		seen[key] = true

		switch item.Source {
		case "baseline":
			scene := item.ScanScene
			if scene <= 0 {
				scene = enums.HostScanSceneBaseline
			}
			if cancel := takeBaselineBatchCancel(item.TaskID, item.TargetIP, scene); cancel != nil {
				cancel()
			}
			updateHostBaselineBatchStopped(item.TaskID, item.TargetIP, scene)
			var resultModel mysqls.BaselineCheckResult
			if _, err := resultModel.StopPendingByTaskTargetAndScene(ctx, item.TaskID, item.TargetIP, scene, "任务已手动结束"); err != nil {
				return nil, err
			}
			total, err := resultModel.CountByTaskTargetAndScene(ctx, item.TaskID, item.TargetIP, scene)
			if err != nil {
				return nil, err
			}
			if total == 0 {
				row := &mysqls.BaselineCheckResult{
					TaskID:          item.TaskID,
					TargetIP:        item.TargetIP,
					ScanScene:       scene,
					RuleID:          0,
					RuleName:        "任务已手动结束",
					CheckResult:     4,
					ActualValue:     "任务已手动结束",
					FixSuggestion:   "",
					RiskDescription: "",
				}
				if err := row.Add(ctx); err != nil {
					return nil, err
				}
			}
		case "vuln":
			if cancel := takeCveBatchCancel(item.TaskID, item.TargetIP); cancel != nil {
				cancel()
			}
			updateHostCveBatchStopped(item.TaskID, item.TargetIP)
			_ = services.PersistHostVulnScanResults(ctx, &services.VulnScanReport{
				TaskID:   item.TaskID,
				TargetIP: item.TargetIP,
			}, fmt.Errorf("任务已手动结束"))
		case "malware":
			if cancel := takeYaraBatchCancel(item.TaskID, item.TargetIP); cancel != nil {
				cancel()
			}
			updateHostYaraBatchStopped(item.TaskID, item.TargetIP)
			_ = services.PersistHostMalwareScanResults(ctx, &services.YaraMalwareReport{
				TaskID:   item.TaskID,
				TargetIP: item.TargetIP,
			}, fmt.Errorf("任务已手动结束"))
		default:
			return nil, fmt.Errorf("unsupported source: %s", item.Source)
		}
		stopped++
	}
	return &typespec.HostSecTaskStopResp{Stopped: stopped}, nil
}
