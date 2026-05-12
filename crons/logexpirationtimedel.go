package crons

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"smart/services"
	"smart/tools/enums"
	tooltime "smart/tools/time"
)

type ExpirationTime struct {
	ExpirationTime int `json:"expirationTime"`
}

func DelExpirationTimeLog() {
	var (
		mapSetService    services.MapSet
		ctx              = context.Background()
		logBackupService services.LogBackup
	)
	//查询日志过期配置
	mapSetValue, err := mapSetService.GetMapValue(ctx, enums.LogExpTimeConfigMapSetObjKey)
	if err != nil || len(mapSetValue) == 0 {
		log.Info("crons DelExpirationTimeLog get config err or empty!")
		return
	}
	var tmpExpiraTime ExpirationTime
	if err = json.Unmarshal([]byte(mapSetValue), &tmpExpiraTime); err != nil {
		log.Info("crons DelExpirationTimeLog json unmarshal err:", err)
		return
	}
	//计算过期时间节点
	delTimeString := tooltime.GetBeforeDayTime(tmpExpiraTime.ExpirationTime)
	//删除小于过期时间节点的数据
	logRes, err := logBackupService.GetLogBackupByLtCreateTime(ctx, delTimeString)
	if err != nil {
		log.Info("crons DelExpirationTimeLog GetLogBackupByLtCreateTime err:", err)
		return
	}
	for i := 0; i < len(logRes); i++ {
		err = logBackupService.DelLogBackupMany(ctx, logRes[i].ID, "./"+logRes[i].Path)
		if err != nil {
			log.Info("crons DelExpirationTimeLog DelLogBackupMany err:", err)
		}
	}

}
