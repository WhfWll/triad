package services

import (
	"context"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

type Node struct {
}

// 节点管理 - 依据ip与端口查询节点
func (sm *Node) SystemNodeGetByIpPort(ctx context.Context, ip, port string) (mysqls.ScannerNode, error) {
	var yakNodeModel mysqls.ScannerNode
	return yakNodeModel.GetByIpPort(ctx, ip, port)
}

// 节点管理 - 依据id
func (sm *Node) SystemNodeGetById(ctx context.Context, id int) (mysqls.ScannerNode, error) {
	var yakNodeModel mysqls.ScannerNode
	yakNodeModel.ID = id
	return yakNodeModel.GetYakNode(ctx)
}

// 节点管理 - 新增节点
func (sm *Node) SystemNodeAdd(ctx context.Context, name, ip, port string) error {
	var yakNodeModel mysqls.ScannerNode
	yakNodeModel.Name = name
	yakNodeModel.IP = ip
	yakNodeModel.Port = port
	yakNodeModel.Status = enums.YakNodeStatusOnline
	yakNodeModel.IsDisable = enums.YakNodeIsDisableN
	yakNodeModel.CreateTime = time.Now()
	yakNodeModel.UpdateTime = time.Now()
	yakNodeModel.AddYakNode(ctx)
	return nil
}

// 节点管理 - 编辑节点
func (sm *Node) SystemNodeEdit(ctx context.Context, id int, name, ip, port string) error {
	var yakNodeModel mysqls.ScannerNode
	yakNodeModel.ID = id
	yakNodeModel.Name = name
	yakNodeModel.IP = ip
	yakNodeModel.Port = port
	yakNodeModel.UpdateTime = time.Now()
	return yakNodeModel.UpdateYakNode(ctx)
}

// 节点管理 - 节点列表
func (sm *Node) SystemNodeList(ctx context.Context, page, size int, search string) ([]mysqls.ScannerNode, int64, error) {
	var yakNodeModel mysqls.ScannerNode
	return yakNodeModel.GetYakNodeList(ctx, page, size, search)
}

// 节点管理 - 删除节点
func (sm *Node) SystemNodeDel(ctx context.Context, ids []string) error {
	var yakNodeModel mysqls.ScannerNode
	return yakNodeModel.DeleteYakNode(ctx, ids)
}

// 节点管理 - 设置禁用｜启用节点
func (sm *Node) SystemNodeDisOrEnable(ctx context.Context, id, isDisbale int) error {
	var yakNodeModel mysqls.ScannerNode
	return yakNodeModel.UpdateIsDisable(ctx, id, isDisbale)
}

// 节点管理 - 设置状态
func (sm *Node) SystemNodeSetStatus(ctx context.Context, id, status int) error {
	var yakNodeModel mysqls.ScannerNode
	return yakNodeModel.UpdateStatus(ctx, id, status)
}

// 节点管理 - 所有可用节点
func (sm *Node) SystemNodeAllEnable(ctx context.Context) []mysqls.ScannerNode {
	var yakNodeModel mysqls.ScannerNode
	return yakNodeModel.AllEnbaleNode(ctx)
}

// 节点管理 - 所有节点
func (sm *Node) SystemNodeAll(ctx context.Context) []mysqls.ScannerNode {
	var yakNodeModel mysqls.ScannerNode
	return yakNodeModel.AllNode(ctx)
}

// 节点管理 - 是否有可用节点 bool
func (sm *Node) SystemNodeIsAvailable(ctx context.Context) bool {
	var yakNodeModel mysqls.ScannerNode
	yakNodeList := yakNodeModel.AllNode(ctx)

	isAvailable := false
	for _, item := range yakNodeList {
		// 在线且启用状态
		if item.Status == enums.YakNodeStatusOnline && item.IsDisable == enums.YakNodeIsDisableN {
			isAvailable = true
		}
	}
	return isAvailable
}
