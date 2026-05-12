package services

import (
	"context"
	"smart/models/mysqls"
	"sort"
	"strconv"
	"strings"
)

type AssetPort struct{}

// GetAssetPortInfo 查询
func (ap *AssetPort) GetAssetPortInfo(ctx context.Context, ips []string) string {
	var (
		AssetPortModel mysqls.Assetport
		ports          string
	)
	list := AssetPortModel.GetAssetportByIps(ctx, ips, nil, nil)
	for _, v := range list {
		ports += strconv.Itoa(v.Port) + ","
	}
	return strings.TrimSuffix(ports, ",")
}

// GetAssetPortList 列表
func (ap *AssetPort) GetAssetPortList(ctx context.Context, ips []string) []mysqls.Assetport {
	var AssetPortModel mysqls.Assetport
	return AssetPortModel.GetAssetportByIps(ctx, ips, nil, nil)
}

// 根据ip获取端口数据
func (ap *AssetPort) GetAssetPortByIp(ctx context.Context, ip string) []mysqls.Assetport {
	var AssetPortModel mysqls.Assetport
	return AssetPortModel.GetAssetPortByIp(ctx, ip)
}

// GetAssetPortByService 根据服务返回服务下的端口信息
func (ap *AssetPort) GetAssetPortByService(ctx context.Context, service, ip string) string {
	var (
		AssetPortModel mysqls.Assetport
		ports          string
	)
	list := AssetPortModel.GetAssetPortByServiceAndIP(ctx, service, ip)

	for _, v := range list {
		ports += strconv.Itoa(v.Port) + ","
	}
	return strings.TrimSuffix(ports, ",")
}

// GetAssetPortByFinger 根据服务返回指纹下的端口信息
func (ap *AssetPort) GetAssetPortByFinger(ctx context.Context, finger, ip string) string {
	var (
		AssetPortModel mysqls.Assetport
		ports          string
	)
	list := AssetPortModel.GetAssetPortByFingerAndIP(ctx, finger, ip)

	for _, v := range list {
		ports += strconv.Itoa(v.Port) + ","
	}
	return strings.TrimSuffix(ports, ",")
}

// GetAssetPortForTarget 查询资产端口
func (ap *AssetPort) GetAssetPortForTarget(ctx context.Context, ip string) string {
	var (
		targetSrv TaskTarget
		ids       []int
		portSet   = make(map[string]struct{}) // 用于去重
	)
	target, _ := targetSrv.TargetListByTargetUrl(ctx, ip)
	targetRes, _ := targetSrv.GetTargetByTaskId(ctx, target.TaskID)
	for i := 0; i < len(targetRes); i++ {
		if targetRes[i].TargetURL == "" || !strings.Contains(targetRes[i].TargetURL, ip) {
			continue
		}
		ids = append(ids, targetRes[i].ID)
	}
	openRes, _ := targetSrv.GetTargetOpenPort(ctx, ids)
	// 去重收集端口
	for _, portList := range openRes {
		for _, port := range portList {
			portSet[port] = struct{}{}
		}
	}
	var uniquePorts []string
	for port := range portSet {
		uniquePorts = append(uniquePorts, port)
	}
	sort.Strings(uniquePorts) // 保证返回有序，便于前端展示
	return strings.Join(uniquePorts, ",")
}

// AddAssetPortByTaskInfo 通过任务更新资产端口信息
func (ap *AssetPort) AddAssetPortByTaskInfo(ctx context.Context, taskID int) string {
	//var (
	//	targetSrv TaskTarget
	//	ids       []int
	//	portSet   = make(map[string]struct{}) // 用于去重
	//	newPorts  []mysqls.Assetport
	//)
	//targetRes, _ := targetSrv.GetTargetByTaskId(ctx, taskID)
	//for i := 0; i < len(targetRes); i++ {
	//	if targetRes[i].TargetURL == "" {
	//		continue
	//	}
	//	ids = append(ids, targetRes[i].ID)
	//}
	//
	//openRes, _ := targetSrv.GetAssetTargetOpenPort(ctx, ids)
	//
	//// 2. 获取开放端口信息
	//openRes, _ := targetSrv.GetTargetOpenPort(ctx, ids)
	//
	//// 3. 遍历生成 Assetport 实例
	//for targetID, portList := range openRes {
	//	ip := ipMap[targetID]
	//	for _, portStr := range portList {
	//		// 端口 JSON 解析
	//		var parsed struct {
	//			Port      string `json:"port"`
	//			Protocol  string `json:"protocol"`
	//			Service   string `json:"service"`
	//			Component string `json:"component"`
	//		}
	//		if err := json.Unmarshal([]byte(portStr), &parsed); err != nil {
	//			continue
	//		}
	//
	//		portInt, err := strconv.Atoi(parsed.Port)
	//		if err != nil || portInt <= 0 {
	//			continue
	//		}
	//
	//		// ip+port 去重
	//		key := fmt.Sprintf("%s:%d", ip, portInt)
	//		if _, exists := portSet[key]; exists {
	//			continue
	//		}
	//		portSet[key] = struct{}{}
	//
	//		newPorts = append(newPorts, Assetport{
	//			IP:         ip,
	//			Port:       portInt,
	//			Protocol:   parsed.Protocol,
	//			Service:    parsed.Service,
	//			Assembly:   parsed.Component,
	//			Islive:     1,
	//			CreateTime: time.Now(),
	//			UpdateTime: time.Now(),
	//		})
	//	}
	//}
	//
	//// 4. 批量写入数据库
	//if len(newPorts) > 0 {
	//	db := mysql.FromContext(ctx)
	//	if err := db.Model(&Assetport{}).Create(&newPorts).Error; err != nil {
	//		log.Printf("Insert asset ports error: %v", err)
	//	}
	//}
	//// 5. 返回端口列表用于展示
	//var allPorts []string
	//for k := range portSet {
	//	parts := strings.Split(k, ":")
	//	if len(parts) == 2 {
	//		allPorts = append(allPorts, parts[1])
	//	}
	//}
	//sort.Strings(allPorts)
	//return strings.Join(allPorts, ",")
	return ""
}
