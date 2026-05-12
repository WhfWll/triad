package yak

import (
	"context"
	"errors"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/redis"
	"google.golang.org/grpc"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
	"smart/tools/network"
	"strconv"
)

// 此文件已废弃 BAS执行已由高伟那边接手，使用熊哥开发的agent进行处理
// 可以删除
var YakGrpc yakGrpc

type yakGrpc struct {
}

func (y *yakGrpc) getConnect() (*grpc.ClientConn, error) {
	yakGrpcAddr := make(map[string]string)
	if err := config.Load("yak_grpc_addr", &yakGrpcAddr); err != nil {
		return nil, errors.New("获取配置yak_grpc_addr失败" + err.Error())
	}
	if yakGrpcAddr["host"] == "" || yakGrpcAddr["port"] == "" {
		return nil, errors.New("配置yak_grpc_addr.host或yak_grpc_addr.port不可为空")
	}
	// 验证yak grpc是否可通信
	if !network.TelnetIsOpen(yakGrpcAddr["host"], yakGrpcAddr["port"]) {
		return nil, errors.New("配置yak_grpc_addr.host=" + yakGrpcAddr["host"] + "或yak_grpc_addr.port=" + yakGrpcAddr["port"] + "网络不可达")
	}

	// 创建gRPC连接
	conn, err := grpc.Dial(yakGrpcAddr["host"]+":"+yakGrpcAddr["port"], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.New("无法连接到服务器：" + err.Error())
	}

	return conn, nil
}

// 查询BAS规则库
func (y *yakGrpc) GrpcBasMakerRulesGets(ctx context.Context, keywords ...string) (map[string]*ypb.ChaosMakerRule, error) {
	param := ypb.QueryChaosMakerRuleRequest{
		Pagination: &ypb.Paging{
			Page:    1,
			Limit:   100,
			OrderBy: "",
			Order:   "",
		},
		RuleType: "",
		Keywords: keywords,
	}

	conn, err := y.getConnect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := ypb.NewYakClient(conn)
	res, err := client.QueryChaosMakerRule(ctx, &param)
	if err != nil {
		return nil, errors.New("获取BAS规则库 err：" + err.Error())
	}

	// 处理为返回已名字为标准返回map
	dataMap := make(map[string]*ypb.ChaosMakerRule)
	for _, item := range res.GetData() {
		dataMap[item.Name] = item
	}
	return dataMap, nil
}

// 导入BAS规则
func (y *yakGrpc) GrpcImportBasMakerRules(ctx context.Context, data []ypb.ImportChaosMakerRulesRequest) (string, error) {
	if len(data) == 0 {
		return "", errors.New("规则不可为空")
	}
	conn, err := y.getConnect()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := ypb.NewYakClient(conn)
	for _, item := range data {
		singleData := item
		_, err = client.ImportChaosMakerRules(ctx, &singleData)
		if err != nil {
			return "", errors.New("导入BAS规则库 err：" + err.Error())
		}
	}

	return "", nil
}

// 执行BAS规则
func (y *yakGrpc) ExecuteChaosMakerRuleRequest(ctx context.Context, taskId int, vulInBoxAddr string, groups []string) ([]ypb.ExecResult, error) {
	// agent地址
	add := make([]string, 0)
	add = append(add, vulInBoxAddr)

	// 仅执行指定的规则
	if len(groups) == 0 {
		return nil, errors.New("执行规则groups不可未空")
	}
	group := make([]*ypb.ChaosMakerRuleGroup, 0)
	for _, item := range groups {
		group = append(group, &ypb.ChaosMakerRuleGroup{
			Keywords: item,
		})
	}

	// 构造执行参数
	execChaosMakerRuleParam := ypb.ExecuteChaosMakerRuleRequest{
		Groups:                          group,
		ExtraOverrideDestinationAddress: add,
		// 随机延迟
		Concurrent:             1,
		TrafficDelayMinSeconds: 1,
		TrafficDelayMaxSeconds: 5,
		// 额外重复，如果为 -1 认为是永久重复
		ExtraRepeat: 1,
		// 每组流量之间重复的次数
		GroupGapSeconds: 1,
	}

	conn, err := y.getConnect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 开始执行
	client := ypb.NewYakClient(conn)
	execStream, err := client.ExecuteChaosMakerRule(ctx, &execChaosMakerRuleParam)
	if err != nil {
		return nil, errors.New("无法连接到服务器：%v" + err.Error())
	}

	// 接收执行日志
	redisClient, err := redis.NewClient()
	if err != nil {
		return nil, errors.New("结束BAS任务获取redisClient失败：" + err.Error())
	}
	// 由于Yak未开发完毕，所以拿不到每个规则的执行结束标识，并且任务永远不会停止，只有客户端自己停止了，默认单规则收取20条日志自动结束
	logs := make([]ypb.ExecResult, 0)
	flag := 20 * len(groups)
	for flag > 0 {
		// 是否结束
		if redisClient.Get(ctx, "bas_end_task_"+strconv.Itoa(taskId)).String() == "Y" {
			flag = 0
			continue
		}

		res, err := execStream.Recv()
		if err != nil {
			flag = 0
			continue
		}
		flag--
		logs = append(logs, *res)
	}
	return logs, nil
}
