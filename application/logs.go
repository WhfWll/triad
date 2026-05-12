package application

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/enums"
	"strings"
	"time"
)

type Logs struct {
}

// LogsEnum 日志枚举
func (l *Logs) LogsEnum(res *typespec.LogsEnumRes) {
	var logAuditService services.LogAudit
	res.Type = logAuditService.LogAuditTypeEnum()
}

// LogAuditList 审计日志-列表
func (l *Logs) LogAuditList(ctx context.Context, req *typespec.LogAuditListReq, res *typespec.LogAuditListRes) error {
	var logAuditService services.LogAudit
	logAuditList, total, err := logAuditService.LogAuditList(ctx, req.Page, req.Size, req.LogType, req.Search, req.StartTime, req.EndTime)
	if err != nil {
		return err
	}

	res.Total = total
	for _, item := range logAuditList {
		res.List = append(res.List, typespec.LogAuditListItemRes{
			Id:          item.ID,
			LogType:     item.LogType,
			LogTypeName: enums.GetLogAuditType(item.LogType),
			Content:     item.Content,
			Username:    item.Username,
			Ip:          item.Ip,
			CreateTime:  item.CreateTime.Format(enums.TimeLayout),
			UpdateTime:  item.UpdateTime.Format(enums.TimeLayout),
		})
	}
	return nil
}

// LogAuditEmpty 审计日志-清空
func (l *Logs) LogAuditEmpty(ctx context.Context) error {
	var logAuditService services.LogAudit
	return logAuditService.LogAuditEmpty(ctx)
}

// LogBackupConfig 日志管理-日志备份-配置
func (l *Logs) LogBackupConfig(ctx context.Context, req *typespec.LogBackupConfigReq) error {
	var mapSetService services.MapSet
	currentTime := time.Now()
	if req.SaveTime.IsZero() {
		req.SaveTime = currentTime
	}
	if req.RunTime.IsZero() {
		req.RunTime = currentTime.AddDate(0, req.Cycle, 0)
	}
	objValueStruct := services.LogBackupConfigMapSet{
		IsOpen:   req.IsOpen,
		Cycle:    req.Cycle,
		SaveTime: req.SaveTime,
		RunTime:  req.RunTime,
	}
	dataByte, err := json.Marshal(objValueStruct)
	if err != nil {
		return err
	}
	objValue := string(dataByte)

	return mapSetService.Create(ctx, enums.LogBackupConfigMapSetObjKey, objValue, enums.LogBackupConfigMapSetContent)
}

// LogBackupConfigInfo 日志管理-日志备份-配置信息
func (l *Logs) LogBackupConfigInfo(ctx context.Context, res *typespec.LogBackupConfigInfoRes) error {
	var mapSetService services.MapSet
	mapSetValue, err := mapSetService.GetMapValue(ctx, enums.LogBackupConfigMapSetObjKey)
	if err != nil {
		return err
	}
	if mapSetValue == "" {
		return nil
	}
	var logBackupConfigMapSet services.LogBackupConfigMapSet
	if err = json.Unmarshal([]byte(mapSetValue), &logBackupConfigMapSet); err != nil {
		return err
	}
	res.Cycle = logBackupConfigMapSet.Cycle
	res.IsOpen = logBackupConfigMapSet.IsOpen
	return nil
}

// 设置日志过期时间
func (l *Logs) SetLogExpirationTime(ctx context.Context, req *typespec.SetLogExpirationTimeReq) error {
	var mapSetService services.MapSet
	dataByte, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return mapSetService.Create(ctx, enums.LogExpTimeConfigMapSetObjKey, string(dataByte), enums.LogExpTimeConfigMapSetContent)
}

// 查询日志过期时间
func (l *Logs) GetLogExpirationTime(ctx context.Context, resp *typespec.GetLogExpirationTimeResp) error {
	var mapSetService services.MapSet
	mapSetValue, err := mapSetService.GetMapValue(ctx, enums.LogExpTimeConfigMapSetObjKey)
	if err != nil {
		return err
	}
	if len(mapSetValue) == 0 {
		resp.ExpirationTime = 0
		return nil
	}
	if err = json.Unmarshal([]byte(mapSetValue), resp); err != nil {
		return err
	}
	return nil
}

// LogBackupList 日志管理-日志备份-列表
func (l *Logs) LogBackupList(ctx context.Context, req *typespec.LogBackupListReq, res *typespec.LogBackupListRes) error {
	var logBackupService services.LogBackup
	logBackupList, total, err := logBackupService.LogBackupList(ctx, req.Page, req.Size)
	if err != nil {
		return err
	}
	res.Total = total
	for _, logBackup := range logBackupList {
		res.List = append(res.List, typespec.LogBackupListItemRes{
			Id:         logBackup.ID,
			Name:       logBackup.Name,
			Path:       logBackup.Path,
			CreateTime: logBackup.CreateTime.Format(enums.TimeLayout),
		})
	}
	return nil
}

// LogBackupDownload 日志管理-日志备份-下载
func (l *Logs) LogBackupDownload(ctx context.Context, req *typespec.LogBackupDownloadReq, res *typespec.LogBackupDownloadRes) error {
	var logBackupService services.LogBackup
	logBackup, err := logBackupService.LogBackupRecord(ctx, req.Id)
	if err != nil {
		return err
	}
	if logBackup.Path == "" {
		return errors.New("下载文件不存在")
	}

	// 校验文件路径是否合法
	cleanPath := filepath.Clean(logBackup.Path)
	cleanBase := filepath.Clean(enums.LogBackupDir)
	if !strings.HasPrefix(cleanPath, cleanBase) {
		return errors.New("非法文件路径")
	}

	res.Path = logBackup.Path
	return nil
}

// LogBackupDelete 日志管理-日志备份-删除
func (l *Logs) LogBackupDelete(ctx context.Context, req *typespec.LogBackupDeleteReq) error {
	var logBackupService services.LogBackup
	return logBackupService.LogBackupDelete(ctx, req.Id)
}

// LogBackupNow 日志管理-日志备份-立即备份
func (l *Logs) LogBackupNow(ctx context.Context) error {
	var logBackupService services.LogBackup
	return logBackupService.LogBackupNow(ctx)
}
