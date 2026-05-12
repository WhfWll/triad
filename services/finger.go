package services

import (
	"context"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

type Finger struct{}

// GetFingerClassEnum 指纹分类枚举
func (a *Finger) GetFingerClassEnum() (result []typespec.GlobalOptionsItemRes) {
	return toolsSort(enums.FingerEnum.AllNewClass())
}

// GetFingerIsDevEmum 指纹设备枚举
func (a *Finger) GetFingerIsDevEmum() (result []typespec.GlobalOptionsItemRes) {
	return toolsSort(enums.FingerEnum.AllIsDev())
}

// GetFingerLevelEnum 指纹设备分层枚举
func (a *Finger) GetFingerLevelEnum() (result []typespec.GlobalOptionsItemRes) {
	return toolsSort(enums.FingerEnum.AllFingerLevel())
}

// GetFingerSoftOrHardEnum 指纹设备软硬件枚举
func (a *Finger) GetFingerSoftOrHardEnum() (result []typespec.GlobalOptionsItemRes) {
	return toolsSort(enums.FingerEnum.AllFingerSoftOrHard())
}

// GetFingerList 指纹列表
func (a *Finger) GetFingerList(ctx context.Context, page, limit int, name string, class int) ([]mysqls.Finger, int64, error) {
	var finger mysqls.Finger
	return finger.GetFingerList(ctx, page, limit, name, 0, class)
}

// GetFingerCount 获取指纹总数
func (a *Finger) GetFingerCount(ctx context.Context) (int64, error) {
	var finger mysqls.Finger
	return finger.Count(ctx)
}

// GetAllFinger 获取所有指纹数据
func (a *Finger) GetAllFinger(ctx context.Context) ([]mysqls.Finger, error) {
	var finger mysqls.Finger
	return finger.GetAllFinger(ctx)
}

// GetFingerDetail 指纹详情
func (a *Finger) GetFingerDetail(ctx context.Context, id int) (mysqls.Finger, error) {
	var finger = mysqls.Finger{
		ID: id,
	}
	return finger.GetFinger(ctx)
}

// SaveFinger 保存指纹
func (a *Finger) SaveFinger(ctx context.Context, appName, appVersion, flag, desc, level string, fingerType, source, appClass int) (mysqls.Finger, error) {
	now := time.Now()
	var finger = mysqls.Finger{
		AppVersion: appVersion,
		CnName:     appName,
		AppName:    appName,
		Flag:       flag,
		AppClass:   appClass,
		Desc:       desc,
		Level:      level,
		FingerType: fingerType,
		Source:     source,
		CreateTime: now,
		UpdateTime: &now,
	}
	err := finger.AddFinger(ctx)
	return finger, err
}

// EditFinger 编辑指纹
func (a *Finger) EditFinger(ctx context.Context, id int, appName, appVersion, flag, desc, level string, fingerType, source, appClass int) (mysqls.Finger, error) {
	now := time.Now()
	var finger = mysqls.Finger{
		ID:         id,
		AppVersion: appVersion,
		CnName:     appName,
		AppName:    appName,
		Flag:       flag,
		AppClass:   appClass,
		FingerType: fingerType,
		Source:     source,
		Desc:       desc,
		Level:      level,
		UpdateTime: &now,
	}
	err := finger.UpdateFinger(ctx)
	return finger, err
}

// DeleteFinger 删除指纹
func (a *Finger) DeleteFinger(ctx context.Context, id int) error {
	var finger = mysqls.Finger{
		ID: id,
	}
	err := finger.DeleteFinger(ctx)
	return err
}

// DeleteAllFinger 清理掉所有指纹
func (a *Finger) DeleteAllFinger(ctx context.Context) error {
	var finger mysqls.Finger
	err := finger.DeleteAllFinger(ctx)
	return err
}

//
//// HistoryResults 指纹历史检测结果
//func (a *Finger) HistoryResults(ctx context.Context, page int, size int) ([]mysqls.Fingerresult, int64, error) {
//	var fingerResult mysqls.Fingerresult
//	return fingerResult.GetFingerresultList(ctx, page, size)
//}
//
//// CreateFingerTest 指纹历史检测结果
//func (a *Finger) CreateFingerTest(ctx context.Context, toolParamList []map[string]string) string {
//	toolName := enums.ScriptNameFingerPrint
//	callInfo := invoke.CallInfo{
//		CallId:        utils.GetUUID(),
//		TaskId:        "",
//		ToolName:      toolName,
//		ToolParamList: make([]invoke.ToolParam, 0),
//		Status:        enums.CallStatusCreated,
//	}
//	for _, toolParamMap := range toolParamList {
//		toolParam := invoke.ToolParam{
//			ParamName:  toolParamMap["key"],
//			ParamValue: toolParamMap["value"],
//		}
//		callInfo.ToolParamList = append(callInfo.ToolParamList, toolParam)
//	}
//	go callInfo.Invoke(ctx, a.fingerTestCallBackFunc)
//	return callInfo.CallId
//}
//
//// CreateAppointFingerTest 指定指纹检测
//func (a *Finger) CreateAppointFingerTest(ctx context.Context, toolParamList []map[string]string, fingerID int) string {
//	toolName := enums.ScriptNameFingerPrint
//	callInfo := invoke.CallInfo{
//		CallId:        utils.GetUUID(),
//		TaskId:        "",
//		ToolName:      toolName,
//		ToolParamList: make([]invoke.ToolParam, 0),
//		Status:        enums.CallStatusCreated,
//	}
//	filePath, err := writeFingerTempFile(ctx, callInfo.CallId, fingerID)
//	if err != nil {
//		return err.Error()
//	}
//	for _, toolParamMap := range toolParamList {
//		toolParam := invoke.ToolParam{
//			ParamName:  toolParamMap["key"],
//			ParamValue: toolParamMap["value"],
//		}
//		callInfo.ToolParamList = append(callInfo.ToolParamList, toolParam)
//	}
//	callInfo.ToolParamList = append(callInfo.ToolParamList, invoke.ToolParam{
//		ParamName:  "rules",
//		ParamValue: filePath,
//	})
//	go callInfo.Invoke(ctx, a.fingerTestCallBackFunc)
//	return callInfo.CallId
//}
//
//func writeFingerTempFile(ctx context.Context, callID string, fingerID int) (string, error) {
//	fingerModel := mysqls.Finger{
//		ID: fingerID,
//	}
//	fingerInfo, err := fingerModel.GetFinger(ctx)
//	if err != nil {
//		return "", err
//	}
//	tempDir := os.TempDir()
//	fileName := "finger_test_" + callID + ".json"
//	tmpFilePath := fmt.Sprintf("%s/%s", tempDir, fileName)
//	tmpFile, err := os.Create(tmpFilePath)
//	if err != nil {
//		return "", err
//	}
//	fingerWriteInfo := []FingerInfoJson{{
//		AppClass:   enums.FingerEnum.GetClass(fingerInfo.AppClass),
//		AppName:    fingerInfo.AppName,
//		AppVersion: fingerInfo.AppVersion,
//		CnName:     fingerInfo.CnName,
//		Flag:       fingerInfo.Flag,
//		Level:      fingerInfo.Level,
//	}}
//	jsonData, err := json.MarshalIndent(fingerWriteInfo, "", "  ")
//	if err != nil {
//		return "", err
//	}
//	_, err = tmpFile.Write(jsonData)
//	if err != nil {
//		return "", err
//	}
//	return tmpFilePath, nil
//}
//
//// 获取脚本结果的回调函数
//func (a *Finger) fingerTestCallBackFunc(ctx context.Context, callInfo *invoke.CallInfo, result string) {
//	var list redises.RedisList
//	cacheKey := "decision:call:" + callInfo.CallId
//	utils.Lock(ctx, "lock:"+cacheKey)
//	err := list.SetListRPush(ctx, cacheKey, result)
//	if err != nil {
//		fmt.Println(err)
//	}
//	utils.Release(ctx, "lock:"+cacheKey)
//
//	if !strings.Contains(result, "CreatedAt") || !strings.Contains(result, "Details") {
//		return
//	}
//
//	fmt.Println(result)
//	var tempData map[string]interface{}
//	err = json.Unmarshal([]byte(result), &tempData)
//	var fingerResult mysqls.Fingerresult
//	err = fingerResult.AddFingerresult(ctx)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//// SyncFinger 同步指纹代码
//func (a *Finger) SyncFinger(ctx context.Context, fingerDataList []typespec.ImportFingerData) error {
//	filePath := "/opt/laozhi/decision/" // 你要检查的文件路径
//	if _, err := os.Stat(filePath); os.IsNotExist(err) {
//		filePath = "/opt/huwang/decision/finger/finger_rules.json"
//	} else {
//		filePath = "/opt/laozhi/decision/finger/finger_rules.json"
//	}
//	var syncFingerList []map[string]string
//	appClassEnum := enums.FingerEnum.AllClass()
//	for _, finger := range fingerDataList {
//		syncFingerList = append(syncFingerList, map[string]string{
//			"app_class":   appClassEnum[finger.AppClass],
//			"app_version": finger.AppVersion,
//			"cn_name":     finger.CnName,
//			"app_name":    finger.AppName,
//			"flag":        finger.Flag,
//			"level":       finger.Level,
//		})
//	}
//	fingerData, _ := json.Marshal(syncFingerList)
//	text := pretty.Pretty(fingerData)
//	_, err := os.Stat(filePath)
//	if err != nil {
//		if os.IsNotExist(err) {
//			_, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0666)
//			if err != nil {
//				log.Println("创建文件失败", err)
//			}
//		}
//	}
//	err = ioutil.WriteFile(filePath, text, 0666)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// GetFingerByAppName 通过名称获取指纹
func (a *Finger) GetFingerByAppName(ctx context.Context, appName string) (mysqls.Finger, error) {
	var finger mysqls.Finger
	return finger.GetFingerByAppName(ctx, appName)
}
