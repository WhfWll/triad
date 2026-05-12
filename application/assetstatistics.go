package application

import (
	"context"
	log "github.com/sirupsen/logrus"
	"smart/api/typespec"
	"smart/services"
	"sort"
	"sync"
	"time"
)

// AssetCollectApp 资产收集统计
type AssetCollectApp struct{}

func (aca *AssetCollectApp) AssetSummarizeTS(ctx context.Context, req *typespec.AssetSummarizeReq, resp *typespec.AssetSummarizeRes) error {
	var (
		assetGroupSrv services.AssetGroup
		wg            sync.WaitGroup
	)
	// 定义结果变量
	var (
		assetTotalCount, assetLastChangeCount                        int64
		newAddIP, thisWeekNewAddIP, newReduceIP, thisWeekNewReduceIP int64
		assetRiskStat                                                []typespec.AssetRiskStatistics
		riskAssetCount                                               int
		riskStatics                                                  []typespec.AssetRiskStatistics
		assetRiskTrend                                               []typespec.AssetRiskTrend
		assetTypeTrend                                               []typespec.AssetRiskTypeTrend
		recentDangerAssets                                           []typespec.RiskAsset
	)

	start := time.Now()

	// 1. 基础并发任务（无依赖）
	wg.Add(5)

	go func() {
		defer wg.Done()
		t := time.Now()
		assetTotalCount, assetLastChangeCount, newAddIP, thisWeekNewAddIP, newReduceIP, thisWeekNewReduceIP =
			assetGroupSrv.GetAssetStaticsInfo(ctx)
		log.Printf("[耗时] GetAssetStaticsInfo: %v", time.Since(t))
	}()

	// 统计资产风险统计
	go func() {
		defer wg.Done()
		t := time.Now()
		stat, _, riskAssetNum := assetGroupSrv.AssetRiskStatics(ctx)
		assetRiskStat = stat
		riskAssetCount = riskAssetNum
		log.Printf("[耗时] AssetRiskStatics: %v", time.Since(t))
	}()

	// 资产风险类型统计
	go func() {
		defer wg.Done()
		t := time.Now()
		riskStaticsTemp, _ := assetGroupSrv.RiskTypeStaticsTS(ctx)
		riskStatics = riskStaticsTemp
		// 4s
		log.Printf("[耗时] RiskTypeStaticsTS: %v", time.Since(t))
	}()

	// 资产风险分布趋势统计
	go func() {
		defer wg.Done()
		t := time.Now()
		assetRiskTrend, _ = assetGroupSrv.GetCycleAssetRiskForLastNDays(ctx, 15)
		for _, v := range assetRiskTrend {
			unsafe := v.FatalAsset + v.HighAsset + v.MediumAsset + v.LowAsset
			assetTypeTrend = append(assetTypeTrend, typespec.AssetRiskTypeTrend{
				Date:   v.Date,
				Safe:   v.TotalAsset - unsafe,
				UnSafe: unsafe,
			})
		}
		log.Printf("[耗时] GetCycleAssetRiskForLastNDays + 转换: %v", time.Since(t))
	}()

	// 最近危险资产
	go func() {
		defer wg.Done()
		t2 := time.Now()
		recentDangerAssets = assetGroupSrv.LastRiskAssetsTS(ctx)
		// 4s
		log.Printf("[耗时] LastRiskAssetsTS: %v", time.Since(t2))
	}()

	wg.Wait() // 等待以上任务完成（确保 ips 有值）
	log.Printf("[阶段耗时] 并发查询合计耗时: %v", time.Since(start))

	// 返回数据
	*resp = typespec.AssetSummarizeRes{
		AssetTotal:             assetTotalCount,
		SafeLoopholeTotal:      int64(riskAssetCount),
		YoYLastWeekAssetTotal:  assetLastChangeCount,
		NewAddIpNum:            newAddIP,
		TkNewAddIpNum:          thisWeekNewAddIP,
		NewReduceIpNum:         newReduceIP,
		TkNewReduceIpNum:       thisWeekNewReduceIP,
		AssetRiskStat:          assetRiskStat,
		RiskStatics:            riskStatics,
		AssetRiskLevelTrendRes: assetRiskTrend,
		AssetRiskTypeTrendRes:  assetTypeTrend,
		RecentDangerAsset:      recentDangerAssets,
	}
	log.Printf("[总耗时] AssetSummarizeTS 总耗时: %v", time.Since(start))
	return nil
}

// AssetSummarize 资产综述
func (aca *AssetCollectApp) AssetSummarize(ctx context.Context, req *typespec.AssetSummarizeReq, resp *typespec.AssetSummarizeRes) error {
	var (
		assetGroupSrv services.AssetGroup
	)
	assetTotalCount, assetLastChangeCount, newAddIP, thisWeekNewAddIP, newReduceIP, thisWeekNewReduceIP := assetGroupSrv.GetAssetStaticsInfo(ctx)
	// 获取资产风险
	assetRiskStat, _, vulAssetCount := assetGroupSrv.AssetRiskStatics(ctx)
	sortAssetRiskStatistics(assetRiskStat)
	riskStatics, _ := assetGroupSrv.RiskTypeStatics(ctx)
	*resp = typespec.AssetSummarizeRes{
		AssetTotal:             assetTotalCount, // 资产总数
		SafeLoopholeTotal:      int64(vulAssetCount),
		YoYLastWeekAssetTotal:  assetLastChangeCount,
		NewAddIpNum:            newAddIP,
		TkNewAddIpNum:          thisWeekNewAddIP,
		NewReduceIpNum:         newReduceIP,
		TkNewReduceIpNum:       thisWeekNewReduceIP,
		AssetRiskStat:          assetRiskStat,                     // 资产风险统计
		RiskStatics:            riskStatics,                       // 风险类型统计
		AssetRiskLevelTrendRes: assetGroupSrv.AssetRiskTrend(ctx), // 资产风险分布趋势
		AssetRiskTypeTrendRes:  assetGroupSrv.AssetTypeTrend(ctx), // 资产风险趋势
		RecentDangerAsset:      assetGroupSrv.LastRiskAssets(ctx), // 近期危险资产
		//AssetTrendChangeRes:       assetGroupSrv.AssetTrendChange(ctx, req.TimeType, req.StartTime, req.EndTime, "add"),
		//AssetReduceTrendChangeRes: assetGroupSrv.AssetTrendChange(ctx, req.TimeType, req.StartTime, req.EndTime, "reduce"),
	}
	return nil
}

func sortAssetRiskStatistics(stats []typespec.AssetRiskStatistics) {
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].AssetType < stats[j].AssetType
	})
}

// AssetVulStatistics 统计资产漏洞
// NOTE: 过去15天创建的资产 当天最新的漏洞数据总数
func AssetVulStatistics() {
	// 1 先拿到过去15天某一天 以及之前的资产
	// 2 拿到这一天这些资产最新的已完成扫描任务
	// 3 入库保存 这天资产 严重 高 中 低 4类
	// 4 同时 严重 高 中 低算风险主机 其他都是安全
}
