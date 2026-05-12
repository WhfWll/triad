package application

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/redis"
	"gorm.io/gorm"
	"smart/api/typespec"
	"smart/client/grpcclients"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
)

type Node struct {
}

// 节点管理 - 是否启用分布式 - 要求，必须有一个存活节点，如节点挂掉，自动关闭分布式
func (sm *Node) SystemNodeSetDistribute(ctx context.Context, req *typespec.NodeSetDistributeReq) error {
	value := enums.MapSetYakNodeValueN
	if req.Status == 1 {
		// 设置为开启，需要校验是否有可用节点，无可用节点时，不允许开启
		var nodeSrv services.Node
		if !nodeSrv.SystemNodeIsAvailable(ctx) {
			// 没有可用节点
			return errors.New("所有节点不可用，禁止开启分布式")
		}

		value = enums.MapSetYakNodeValueY
	}

	// 获取是否存在配置
	var mapSetSrv services.MapSet
	isEnable, err := mapSetSrv.GetMapValue(ctx, enums.MapSetYakNodeObjKey)
	if err != nil {
		return err
	}

	// 未查到数据，需要新增
	switch isEnable {
	case "": // 无配置，需要添加
		return mapSetSrv.Create(ctx, enums.MapSetYakNodeObjKey, value, enums.MapSetYakNodeContent)
	default: // 存在数据，进行修改
		return mapSetSrv.UpdateMapValue(ctx, enums.MapSetYakNodeObjKey, value)
	}
}

// 节点管理 - 是否启用分布式 - 获取节点状态
func (sm *Node) SystemNodeGetDistribute(ctx context.Context, res *typespec.NodeIsDistributeRes) error {
	// 获取是否存在配置
	var mapSetSrv services.MapSet
	isEnable, err := mapSetSrv.GetMapValue(ctx, enums.MapSetYakNodeObjKey)
	if err == gorm.ErrRecordNotFound {
		res.Status, _ = strconv.Atoi(enums.MapSetYakNodeValueN)
		return nil
	}

	res.Status, _ = strconv.Atoi(isEnable)
	return nil
}

// 节点管理 - 新增节点
func (sm *Node) SystemNodeAdd(ctx context.Context, req *typespec.NodeAddReq, res *typespec.NodeAddRes) error {
	// 校验目标节点是否在线，尝试使用grpc通信
	serverAddr := fmt.Sprintf("%s:%s", req.Ip, req.Port)

	client, err := grpcclients.NewScannerClient(serverAddr)
	if err != nil {
		log.Errorf("创建客户端失败: %v", err)
		return errors.New("目标无法通信，请确认目标IP或端口是否正确")
	}
	defer client.Close()

	healthResp, err := client.HealthCheck()
	if err != nil {
		log.Errorf("健康检查失败: %v", err)
		return errors.New("目标无法通信，请确认目标IP或端口是否正确")
	}

	if healthResp == nil || healthResp.Status != "healthy" {
		return errors.New("目标状态不健康，请确认目标服务器是否正常")
	}

	//if !yak.YakGrpc.IsAlive(ctx, req.Ip, req.Port) {
	//	return errors.New("目标无法通信，请确认目标IP或端口是否正确")
	//}

	var nodeSrv services.Node

	// 节点是否已添加
	yakNode, _ := nodeSrv.SystemNodeGetByIpPort(ctx, req.Ip, req.Port)
	if yakNode.ID != 0 {
		return errors.New("节点ip与端口已存在，请确认")
	}

	nodeSrv.SystemNodeAdd(ctx, req.Name, req.Ip, req.Port)
	return nil
}

// 节点管理 - 编辑节点
func (sm *Node) SystemNodeEdit(ctx context.Context, req *typespec.NodeEditReq, res *typespec.NodeEditRes) error {
	var nodeSrv services.Node

	// 1. 检查节点是否存在
	oldNode, err := nodeSrv.SystemNodeGetById(ctx, req.Id)
	if err != nil {
		return errors.New("节点不存在")
	}
	if oldNode.ID == 0 {
		return errors.New("节点不存在")
	}

	// 2. 校验目标节点是否在线，尝试使用grpc通信
	serverAddr := fmt.Sprintf("%s:%s", req.Ip, req.Port)

	client, err := grpcclients.NewScannerClient(serverAddr)
	if err != nil {
		log.Errorf("创建客户端失败: %v", err)
		return errors.New("目标无法通信，请确认目标IP或端口是否正确")
	}
	defer client.Close()

	healthResp, err := client.HealthCheck()
	if err != nil {
		log.Errorf("健康检查失败: %v", err)
		return errors.New("目标无法通信，请确认目标IP或端口是否正确")
	}

	if healthResp == nil || healthResp.Status != "healthy" {
		return errors.New("目标状态不健康，请确认目标服务器是否正常")
	}

	// 3. 检查IP端口是否冲突（排除自己）
	existNode, _ := nodeSrv.SystemNodeGetByIpPort(ctx, req.Ip, req.Port)
	if existNode.ID != 0 && existNode.ID != req.Id {
		return errors.New("该IP与端口已被其他节点使用，请确认")
	}

	// 4. 执行更新
	return nodeSrv.SystemNodeEdit(ctx, req.Id, req.Name, req.Ip, req.Port)
}

// 节点管理 - 节点详情
func (sm *Node) SystemNodeInfo(ctx context.Context, req *typespec.NodeInfoReq, res *typespec.NodeInfoRes) error {
	var nodeSrv services.Node
	node, err := nodeSrv.SystemNodeGetById(ctx, req.Id)
	if err != nil {
		return err
	}
	if node.ID == 0 {
		return errors.New("节点不存在")
	}

	res.Id = node.ID
	res.Name = node.Name
	res.Ip = node.IP
	res.Port = node.Port
	return nil
}

// 节点管理 - 节点列表
func (sm *Node) SystemNodeList(ctx context.Context, req *typespec.NodeListReq, res *typespec.NodeListRes) error {
	redisClient, err := redis.NewClient()
	if err != nil {
		return errors.New("redis服务获取失败:" + err.Error())
	}

	var nodeSrv services.Node
	list, total, err := nodeSrv.SystemNodeList(ctx, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	res.Total = total

	for _, item := range list {
		runningNum, _ := redisClient.Get(ctx, "decision_running_number_"+item.IP+":"+item.Port).Int()
		res.List = append(res.List, typespec.NodeListItem{
			Id:            item.ID,
			Name:          item.Name,
			Ip:            item.IP,
			Port:          item.Port,
			RunningNum:    runningNum,
			Status:        item.Status,
			StatusEnum:    enums.YakNode.GetStatusEnum(item.Status),
			IsDisable:     item.IsDisable,
			IsDisableEnum: enums.YakNode.GetIsDisableEnum(item.IsDisable),
			CreateTime:    item.CreateTime.Format(utils.DateTime),
			UpdateTime:    item.UpdateTime.Format(utils.DateTime),
		})
	}
	return nil
}

// 节点管理 - 删除节点
func (sm *Node) SystemNodeDel(ctx context.Context, req *typespec.NodeDelReq, res *typespec.NodeDelRes) error {
	ids := strings.Split(req.Id, ",")
	var nodeSrv services.Node
	return nodeSrv.SystemNodeDel(ctx, ids)
}

// 节点管理 - 禁用｜启用节点
func (sm *Node) SystemNodeDisOrEnable(ctx context.Context, req *typespec.NodeDisOrEnableReq, res *typespec.NodeDisOrEnableRes) error {
	var nodeSrv services.Node
	return nodeSrv.SystemNodeDisOrEnable(ctx, req.Id, req.IsDisable)
}

// 节点管理 - 所有可用节点
func (sm *Node) SystemNodeAllEnable(ctx context.Context, req *typespec.NodeAllEnableReq, res *typespec.NodeAllEnableRes) error {
	// 只有开启了分布式，才能选择节点
	//var mapSetSrv services.MapSet
	//distributeStatus, err := mapSetSrv.GetMapValue(ctx, enums.MapSetYakNodeObjKey)
	//if err != nil {
	//	return err
	//}
	//if distributeStatus != enums.MapSetYakNodeValueY {
	//	return errors.New("请先开启分布式功能")
	//}

	var nodeSrv services.Node
	list := nodeSrv.SystemNodeAllEnable(ctx)
	if len(list) == 0 {
		return errors.New("无可用节点，请重新开启节点或先添加节点")
	}

	for _, item := range list {
		res.List = append(res.List, typespec.NodeAllEnableItem{
			Id:   item.ID,
			Name: item.Name,
		})
	}
	return nil
}
