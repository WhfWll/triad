package crons

import (
	"context"
	"smart/client/grpcclients"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
)

// 1 分布式yak节点存活状态校验
// 2 如所有节点都不可用时【状态=下线；是否禁用=1】，自动关闭分布式功能
func scannerNodeIsAlive() {
	ctx := context.Background()

	// 获取所有分布式节点
	var yakNodeSrv services.Node
	yakNodeList := yakNodeSrv.SystemNodeAll(ctx)

	// 1 分布式yak节点存活状态校验
	scannerNodeCheckIsAlive(ctx, yakNodeSrv, yakNodeList)
	// 2 如所有节点都不可用时【状态=下线；是否禁用=1】，自动关闭分布式功能
	//scannerNodeCheckDistribute(ctx, yakNodeSrv, yakNodeList)
}

// 1 分布式yak节点存活状态校验
func scannerNodeCheckIsAlive(ctx context.Context, yakNodeSrv services.Node, yakNodeList []mysqls.ScannerNode) {
	for _, node := range yakNodeList {
		// 存活状态
		go func(item mysqls.ScannerNode) {
			client, err := grpcclients.NewScannerClient(item.IP + ":" + item.Port)
			if err != nil {
				if item.Status == enums.YakNodeStatusOnline {
					_ = yakNodeSrv.SystemNodeSetStatus(ctx, item.ID, enums.YakNodeStatusOffline)
				}
				return
			}
			defer client.Close()
			healthResp, err := client.HealthCheck()

			if healthResp.GetAlive() {
				// 如果状态为离线，更新为在线
				if item.Status == enums.YakNodeStatusOffline {
					_ = yakNodeSrv.SystemNodeSetStatus(ctx, item.ID, enums.YakNodeStatusOnline)
				}
			} else {
				// 不存活
				// 如果状态为
				if item.Status == enums.YakNodeStatusOnline {
					_ = yakNodeSrv.SystemNodeSetStatus(ctx, item.ID, enums.YakNodeStatusOffline)
				}
			}
		}(node)

	}
}

// 2 如所有节点都不可用时【状态=下线；是否禁用=1】，自动关闭分布式功能
func scannerNodeCheckDistribute(ctx context.Context, yakNodeSrv services.Node, yakNodeList []mysqls.ScannerNode) {
	if !yakNodeSrv.SystemNodeIsAvailable(ctx) {
		// 需要关闭分布式
		var mapSetSrv services.MapSet
		_ = mapSetSrv.UpdateMapValue(ctx, enums.MapSetYakNodeObjKey, enums.MapSetYakNodeValueN)
	}
}
