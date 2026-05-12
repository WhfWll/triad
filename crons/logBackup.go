package crons

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"smart/services"
	"smart/tools/enums"
	"time"
)

// LogBackupCron 定时日志备份
func LogBackupCron() {
	log.Info("日志周期备份检测中...")
	ctx := context.Background()
	//获取日志备份配置数据
	var mapSetService services.MapSet
	objValue, err := mapSetService.GetMapValue(ctx, enums.LogBackupConfigMapSetObjKey)
	if err != nil {
		log.Error("LogBackupCron GetMapValue err：" + err.Error())
		return
	}
	if objValue == "" {
		log.Error("LogBackupCron GetMapValue Data Not Exist")
		return
	}
	//解析日志备份配置数据到结构体中
	var logBackupConfigMapSet services.LogBackupConfigMapSet
	if err = json.Unmarshal([]byte(objValue), &logBackupConfigMapSet); err != nil {
		log.Error("LogBackupCron Json Unmarshal err：" + err.Error())
		return
	}
	//是否开启周期日志备份
	if logBackupConfigMapSet.IsOpen != enums.LogBackupCycleOpen {
		log.Info("LogBackupCron cycle backup is not open")
		return
	}
	//根据执行日期作判断
	currentTime := time.Now()
	if currentTime.Before(logBackupConfigMapSet.RunTime) {
		log.Info("LogBackupCron runtime not arrived")
		return
	}
	/* ========== 执行日志备份 ========== */
	var logBackupService services.LogBackup
	if err := logBackupService.LogBackupNow(ctx); err != nil {
		log.Error("LogBackupCron Json Unmarshal err：" + err.Error())
		return
	}
	/* ========== 执行日志备份end ========== */

	/* ========== 更新日志备份配置 ========== */
	//更新ObjValue值
	logBackupConfigMapSet.RunTime = currentTime.AddDate(0, logBackupConfigMapSet.Cycle, 0)
	valueByte, err := json.Marshal(logBackupConfigMapSet)
	if err != nil {
		log.Error("LogBackupCron Json Marshal err：" + err.Error())
		return
	}
	valueStr := string(valueByte)
	//更新下次执行日期到MapSet
	err = mapSetService.Create(ctx, enums.LogBackupConfigMapSetObjKey, valueStr, enums.LogBackupConfigMapSetContent)
	if err != nil {
		log.Error("LogBackupCron save config err：" + err.Error())
		return
	}
	/* ========== 更新日志备份配置end ========== */
}
