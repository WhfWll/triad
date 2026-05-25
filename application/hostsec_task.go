package application

import (
	"context"
	"fmt"

	"smart/api/typespec"
	"smart/models/mysqls"
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
