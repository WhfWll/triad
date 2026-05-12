package crons

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"smart/services"
	"smart/tools/enums"
	"time"
)

// SystemConfigBackupCron 定时系统配置备份
func SystemConfigBackupCron() {
	log.Info("系统配置周期备份检测中...")
	ctx := context.Background()
	//获取系统配置备份配置数据
	var mapSetService services.MapSet
	objValue, err := mapSetService.GetMapValue(ctx, enums.SystemConfigBackupConfigMapSetObjKey)
	if err != nil {
		log.Error("SystemConfigBackupCron GetMapValue err：" + err.Error())
		return
	}
	if objValue == "" {
		log.Error("SystemConfigBackupCron GetMapValue Data Not Exist")
		return
	}
	//解析日志备份配置数据到结构体中
	var systemConfigBackupConfigMapSet services.SystemConfigBackupConfigMapSet
	if err := json.Unmarshal([]byte(objValue), &systemConfigBackupConfigMapSet); err != nil {
		log.Error("SystemConfigBackupCron Json Unmarshal err：" + err.Error())
		return
	}
	//是否开启周期日志备份
	if systemConfigBackupConfigMapSet.IsOpen != enums.ConfigOpen {
		log.Info("SystemConfigBackupCron cycle backup is not open")
		return
	}
	//根据执行日期作判断
	currentTime := time.Now()
	if currentTime.Before(systemConfigBackupConfigMapSet.RunTime) {
		log.Info("SystemConfigBackupCron runtime not arrived")
		return
	}

	/* ========== 执行系统配置备份 ========== */
	var systemConfigBackupService services.SystemConfigBackup
	if err := systemConfigBackupService.SystemConfigBackupNow(ctx); err != nil {
		log.Error("SystemConfigBackupCron runtime err：" + err.Error())
		return
	}
	/* ========== 执行系统配置备份end ========== */

	/* ========== 更新系统配置备份配置 ========== */
	//更新ObjValue值
	systemConfigBackupConfigMapSet.RunTime = currentTime.AddDate(0, systemConfigBackupConfigMapSet.Cycle, 0)
	valueByte, err := json.Marshal(systemConfigBackupConfigMapSet)
	if err != nil {
		log.Error("SystemConfigBackupCron Json Marshal err：" + err.Error())
		return
	}
	//更新下次执行日期到MapSet
	err = mapSetService.Create(ctx, enums.SystemConfigBackupConfigMapSetObjKey, string(valueByte), enums.SystemConfigBackupConfigMapSetContent)
	if err != nil {
		log.Error("SystemConfigBackupCron save config err：" + err.Error())
		return
	}
	/* ========== 更新系统配置备份配置end ========== */
}
