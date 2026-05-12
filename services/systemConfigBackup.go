package services

import (
	"context"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/file"
	"smart/tools/utils"
	"strconv"
	"time"
)

type SystemConfigBackup struct {
}

// SystemConfigBackupAdd 新增系统配置备份
func (s *SystemConfigBackup) SystemConfigBackupAdd(ctx context.Context, name, path string) error {
	var systemConfigBackupModel = mysqls.SystemConfigBackup{
		Name:       name,
		Path:       path,
		CreateTime: time.Now(),
	}
	return systemConfigBackupModel.SystemConfigBackupAdd(ctx)
}

// SystemConfigBackupDelete 删除系统配置备份
func (s *SystemConfigBackup) SystemConfigBackupDelete(ctx context.Context, id int) error {
	var systemConfigBackupModel mysqls.SystemConfigBackup
	return systemConfigBackupModel.SystemConfigBackupDelete(ctx, id)
}

// SystemConfigBackupList 系统配置备份列表
func (s *SystemConfigBackup) SystemConfigBackupList(ctx context.Context, page, limit int) ([]mysqls.SystemConfigBackup, int64, error) {
	var systemConfigBackupModel mysqls.SystemConfigBackup
	return systemConfigBackupModel.SystemConfigBackupList(ctx, page, limit)
}

// SystemConfigBackupRecord 单条系统配置备份信息
func (s *SystemConfigBackup) SystemConfigBackupRecord(ctx context.Context, id int) (mysqls.SystemConfigBackup, error) {
	var systemConfigBackupModel mysqls.SystemConfigBackup
	return systemConfigBackupModel.SystemConfigBackupRecord(ctx, id)
}

// SystemConfigBackupNow 系统配置备份
func (s *SystemConfigBackup) SystemConfigBackupNow(ctx context.Context) error {
	//获取所有的系统配置
	var mapSetModel mysqls.MapSet
	mapSetKeys := []string{
		enums.SystemConfigBackupConfigMapSetObjKey,
		enums.LogBackupConfigMapSetObjKey,
		enums.NetworkConfigMapSetObjKey,
		enums.SyslogConfigMapSetObjKey,
		enums.SystemAccessIpWhiteMapSetObjKey,
	} //ToDo 后续追加安全策略、系统告警、邮件服务
	mapSetList := mapSetModel.ListByObjKeys(ctx, mapSetKeys)

	//组合数据
	dataList := make([][]string, 0)
	dataList = append(dataList, []string{"id", "estate", "obj_key", "obj_value", "content"})
	for _, mapSet := range mapSetList {
		dataList = append(dataList, []string{
			strconv.Itoa(mapSet.ID),
			mapSet.Estate,
			mapSet.ObjKey,
			mapSet.ObjValue,
			mapSet.Content,
		})
	}

	//系统配置备份
	if !file.CheckPathExist(enums.SystemConfigBackupDir) {
		if err := file.CreateDir(enums.SystemConfigBackupDir); err != nil {
			return err
		}
	}
	name := "配置备份文件_" + utils.DatetimeNumberStr()
	filepath := enums.SystemConfigBackupDir + name + ".csv"
	err := file.WriteCsv(dataList, filepath)
	if err != nil {
		return err
	}

	//记录保存到系统配置备份表
	var systemConfigBackupModel = mysqls.SystemConfigBackup{
		Name:       name,
		Path:       filepath,
		CreateTime: time.Now(),
	}
	if err = systemConfigBackupModel.SystemConfigBackupAdd(ctx); err != nil {
		return err
	}

	return nil
}
