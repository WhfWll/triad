package services

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

// 三方工具
type TripartiteTools struct {
}

// TripartiteToolsXrayCreate xray 创建任务
func (t *TripartiteTools) TripartiteToolsXrayCreate(ctx context.Context, taskName, target string, isCrawler bool) (int, error) {
	var xrayModel mysqls.XrayTask
	xrayModel.TaskName = taskName
	xrayModel.Target = target
	if isCrawler {
		xrayModel.IsCrawler = 1
	}
	xrayModel.Status = enums.XrayStatusWait
	xrayModel.CreateTime = time.Now()
	xrayModel.UpdateTime = time.Now()
	err := xrayModel.AddXrayTask(ctx)
	if err != nil {
		return 0, err
	}
	return xrayModel.ID, nil
}

// TripartiteToolsXrayDel xray 删除通过IDs
func (t *TripartiteTools) TripartiteToolsXrayDel(ctx context.Context, xrayIds []int) error {
	var xrayModel mysqls.XrayTask
	err := xrayModel.DelByIds(ctx, xrayIds)
	if err != nil {
		return err
	}

	var xrayResultModel mysqls.XrayTaskResult
	err = xrayResultModel.DelByXrayIds(ctx, xrayIds)
	if err != nil {
		return err
	}
	return nil
}

// TripartiteToolsXrayGetsByStatus xray 通过状态获取所有任务
func (t *TripartiteTools) TripartiteToolsXrayGetsByStatus(ctx context.Context, status int) []mysqls.XrayTask {
	var xrayModel mysqls.XrayTask
	return xrayModel.GetsByStatus(ctx, status)
}

// TripartiteToolsXrayUpdateStatusById xray 更新状态
func (t *TripartiteTools) TripartiteToolsXrayUpdateStatusById(ctx context.Context, id, status int) error {
	var xrayModel mysqls.XrayTask
	return xrayModel.UpdateStatusById(ctx, id, status)
}

// TripartiteToolsXrayResultCreate xray 结果数据储存
func (t *TripartiteTools) TripartiteToolsXrayResultCreate(ctx context.Context, result []TripartiteToolsXrayCreateResultItem) error {
	xrayResultModelDatas := make([]mysqls.XrayTaskResult, 0)
	for _, item := range result {
		xrayResultModelDatas = append(xrayResultModelDatas, mysqls.XrayTaskResult{
			XrayTaskID:  item.XrayTaskId,
			Addr:        item.Addr,
			Payload:     item.Payload,
			RequestInfo: item.RequestInfo,
			Extra:       item.Extra,
			Plugin:      item.Plugin,
			CreateTime:  item.CreateTime,
			UpdateTime:  item.UpdateTime,
		})
	}

	if len(xrayResultModelDatas) > 0 {
		var xrayResultModel mysqls.XrayTaskResult
		err := xrayResultModel.AddsXrayTaskResult(ctx, xrayResultModelDatas)
		if err != nil {
			return err
		}
	}
	return nil
}

// TripartiteToolsXrayFinish xray 任务结束
func (t *TripartiteTools) TripartiteToolsXrayFinish(ctx context.Context, xrayId, riskNum int) error {
	// 更新任务表数据
	var xrayModel mysqls.XrayTask
	err := xrayModel.Finish(ctx, xrayId, riskNum)
	if err != nil {
		return err
	}
	return nil
}

// TripartiteToolsXRayPage xray 任务列表
func (t *TripartiteTools) TripartiteToolsXRayPage(ctx context.Context, page, size int, search string) ([]mysqls.XrayTask, int64, error) {
	var xrayModel mysqls.XrayTask
	return xrayModel.GetXrayTaskList(ctx, page, size, search)
}

// TripartiteToolsXRayDetailPage xray 任务详情列表
func (t *TripartiteTools) TripartiteToolsXRayDetailPage(ctx context.Context, xrayId, page, size int, search string) ([]mysqls.XrayTaskResult, int64, error) {
	var xrayModelResult mysqls.XrayTaskResult
	return xrayModelResult.GetXrayTaskResultList(ctx, xrayId, page, size, search)
}

// TripartiteToolsBurpsuiteCreate burpsuite 任务创建
func (t *TripartiteTools) TripartiteToolsBurpsuiteCreate(ctx context.Context, taskName string, target string) (int, error) {
	var burpsuiteModel mysqls.BurpsuiteTask
	burpsuiteModel.TaskName = taskName
	burpsuiteModel.Target = target
	burpsuiteModel.Risk = 0
	burpsuiteModel.IsCrawler = 1
	burpsuiteModel.Status = enums.BurpsuiteStatusWait
	burpsuiteModel.CreateTime = time.Now()
	burpsuiteModel.UpdateTime = time.Now()
	err := burpsuiteModel.AddBurpsuiteTask(ctx)
	return burpsuiteModel.ID, err
}

// TripartiteToolsBurpsuiteDel burpsuite 删除通过IDs
func (t *TripartiteTools) TripartiteToolsBurpsuiteDel(ctx context.Context, burpsuiteIds []int) error {
	var burpsuiteModel mysqls.BurpsuiteTask
	err := burpsuiteModel.DelByIds(ctx, burpsuiteIds)
	if err != nil {
		return err
	}

	var burpsuiteResultModel mysqls.BurpsuiteTaskResult
	err = burpsuiteResultModel.DelByBurpsuiteIds(ctx, burpsuiteIds)
	if err != nil {
		return err
	}
	return nil
}

// TripartiteToolsBurpsuitePage burpsuite 任务列表
func (t *TripartiteTools) TripartiteToolsBurpsuitePage(ctx context.Context, page, size int, search string) ([]mysqls.BurpsuiteTask, int64, error) {
	var burpsuiteModel mysqls.BurpsuiteTask
	return burpsuiteModel.GetBurpsuiteTaskList(ctx, page, size, search)
}

// TripartiteToolsBurpsuiteGetsByStatus xray 通过状态获取所有任务
func (t *TripartiteTools) TripartiteToolsBurpsuiteGetsByStatus(ctx context.Context, status int) []mysqls.BurpsuiteTask {
	var xrayModel mysqls.BurpsuiteTask
	return xrayModel.GetsByStatus(ctx, status)
}

// TripartiteToolsBurpsuiteUpdateRunning burpsuite 更新状态
func (t *TripartiteTools) TripartiteToolsBurpsuiteUpdateRunning(ctx context.Context, id, originTaskId, status int) error {
	var xrayModel mysqls.BurpsuiteTask
	return xrayModel.UpdateRunning(ctx, id, originTaskId, status)
}

// TripartiteToolsBurpsuiteUpdateStatusById burpsuite 更新状态
func (t *TripartiteTools) TripartiteToolsBurpsuiteUpdateStatusById(ctx context.Context, id, status int) error {
	var xrayModel mysqls.BurpsuiteTask
	return xrayModel.UpdateStatusById(ctx, id, status)
}

// TripartiteToolsBurpsuiteUpdateStatusById burpsuite 更新风险等级
func (t *TripartiteTools) TripartiteToolsBurpsuiteUpdateRiskById(ctx context.Context, id, risk int) error {
	var xrayModel mysqls.BurpsuiteTask
	return xrayModel.UpdateRiskById(ctx, id, risk)
}

// TripartiteToolsBurpsuiteUpdateStatusRiskById burpsuite 更新状态
func (t *TripartiteTools) TripartiteToolsBurpsuiteUpdateStatusRiskById(ctx context.Context, id, status, risk int) error {
	var xrayModel mysqls.BurpsuiteTask
	return xrayModel.UpdateStatusRiskById(ctx, id, status, risk)
}

// TripartiteToolsBurpsuiteResGetsHostAndPathByBurpsuiteId burpsuite 获取所有运行中的主机与路径 [burpsuite_id:origin_result_id:{},burpsuite_id:origin_result_id:{}]
func (t *TripartiteTools) TripartiteToolsBurpsuiteResGetsHostAndPathByBurpsuiteId(ctx context.Context, burpsuiteId []int) map[int]map[int]mysqls.BurpsuiteTaskResult {
	var xrayModelResult mysqls.BurpsuiteTaskResult
	data := xrayModelResult.GetHostAndPathByBurpsuiteId(ctx, burpsuiteId)
	returnData := make(map[int]map[int]mysqls.BurpsuiteTaskResult)
	for _, item := range data {
		if _, ok := returnData[item.BurpsuiteTaskID]; ok {
			returnData[item.BurpsuiteTaskID][item.OriginResultId] = item
		} else {
			burpsuiteResult := make(map[int]mysqls.BurpsuiteTaskResult)
			burpsuiteResult[item.OriginResultId] = item
			returnData[item.BurpsuiteTaskID] = burpsuiteResult
		}
	}
	return returnData
}

// TripartiteToolsBurpsuiteResultAdds burpsuite 添加多条数据
func (t *TripartiteTools) TripartiteToolsBurpsuiteResultAdds(ctx context.Context, datas []mysqls.BurpsuiteTaskResult) error {
	var xrayModelResult mysqls.BurpsuiteTaskResult
	return xrayModelResult.AddsBurpsuiteTaskResult(ctx, datas)
}

// TripartiteToolsBurpsuitePage burpsuite 任务详情列表
func (t *TripartiteTools) TripartiteToolsBurpsuiteResultPage(ctx context.Context, burpsuiteId, page, size int, search string) ([]mysqls.BurpsuiteTaskResult, int64, error) {
	var burpsuiteModel mysqls.BurpsuiteTaskResult
	return burpsuiteModel.GetBurpsuiteTaskResultList(ctx, burpsuiteId, page, size, search)
}

// TripartiteToolsWifiApList wifi 所有在线wifi列表
func (t *TripartiteTools) TripartiteToolsWifiApList(ctx context.Context, search, startDate, endDate string) ([]mysqls.WifiApInfo, int64) {
	var wifiApInfoModel mysqls.WifiApInfo
	return wifiApInfoModel.ApList(ctx, search, startDate, endDate)
}

// TripartiteToolsWifiApGetBySourceMac wifi 通过mac地址获取单个wifi信息
func (t *TripartiteTools) TripartiteToolsWifiApGetBySourceMac(ctx context.Context, mac string) mysqls.WifiApInfo {
	var wifiApInfoModel mysqls.WifiApInfo
	return wifiApInfoModel.GetBySourceMac(ctx, mac)
}

// TripartiteToolsWifiApCreate wifi 创建任务
func (t *TripartiteTools) TripartiteToolsWifiApCreate(ctx context.Context, data WifiCreateData) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	var wifiTaskModel mysqls.WifiTask
	wifiTaskModel.TaskName = data.TaskName
	wifiTaskModel.Mac = data.Mac
	wifiTaskModel.Channel = data.Channel
	wifiTaskModel.Encrypt = enums.WifiTaskEnum.RecoverCryptset(data.Encrypt)
	wifiTaskModel.Carrier = data.Carrier
	wifiTaskModel.Status = data.Status
	wifiTaskModel.Ssid = data.Ssid
	wifiTaskModel.Passwd = data.Passwd
	wifiTaskModel.PasswdSource = data.PasswdSource
	wifiTaskModel.StartTime = data.StartTime
	wifiTaskModel.EndTime = data.EndTime
	wifiTaskModel.IsSimulate = data.IsSimulate
	wifiTaskModel.SimulateDuration = data.SimulateDuration
	wifiTaskModel.IsCrack = data.IsCrack
	wifiTaskModel.PasswdDict = data.PasswdDict
	wifiTaskModel.IsEmbed = data.IsEmbed
	wifiTaskModel.CreateTime = time.Now()
	wifiTaskModel.UpdateTime = time.Now()
	err := wifiTaskModel.AddWifiTask(dCtx)
	if err != nil {
		return err
	}

	// 查询最新插入的ID
	taskId, err := wifiTaskModel.GetLastTaskID(dCtx)
	if err != nil {
		return err
	}

	// 写创建日志
	var wifiTaskLogModel mysqls.WifiTaskLog
	wifiTaskLogModel.ApMac = data.Mac
	wifiTaskLogModel.TaskID = taskId
	wifiTaskLogModel.Content = "任务目标" + data.Ssid + "（" + data.Mac + "）任务创建成功"
	wifiTaskLogModel.GenerateTime = time.Now().Unix()
	err = wifiTaskLogModel.AddWifiTaskLog(dCtx)
	if err != nil {
		return err
	}

	// 提交事物
	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// TripartiteToolsWifiPage wifi 列表
func (t *TripartiteTools) TripartiteToolsWifiPage(ctx context.Context, page, size int, search string) ([]mysqls.WifiTask, int64) {
	var wifiTaskModel mysqls.WifiTask
	return wifiTaskModel.GetWifiTaskList(ctx, page, size, search)
}

// TripartiteToolsWifiLogByTaskId wifi 日志
func (t *TripartiteTools) TripartiteToolsWifiLogByTaskId(ctx context.Context, taskIds []int) (map[int][]mysqls.WifiTaskLog, error) {
	var wifiTaskLogModel mysqls.WifiTaskLog
	data, err := wifiTaskLogModel.GetWifiTaskLogByTaskId(ctx, taskIds)
	returnData := make(map[int][]mysqls.WifiTaskLog)
	for _, item := range data {
		tempLog := returnData[item.TaskID]
		tempLog = append(tempLog, item)
		returnData[item.TaskID] = tempLog
	}

	return returnData, err
}

// wifi 删除任务
func (t *TripartiteTools) TripartiteToolsWifiDel(ctx context.Context, taskIds []int) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	// 删任务
	var wifiTaskModel mysqls.WifiTask
	err := wifiTaskModel.DeleteByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}
	// 删日志
	var wifiTaskLogModel mysqls.WifiTaskLog
	err = wifiTaskLogModel.DeleteByTaskId(dCtx, taskIds)
	if err != nil {
		return err
	}

	// 提交事物
	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
