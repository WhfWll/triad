package services

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/file"
	"smart/tools/utils"
	"strconv"
	"time"
)

type LogBackup struct {
}

// LogBackupAdd 新增备份日志
func (l *LogBackup) LogBackupAdd(ctx context.Context, name, path string) error {
	var logBackupModel = mysqls.LogBackup{
		Name:       name,
		Path:       path,
		CreateTime: time.Now(),
	}
	return logBackupModel.LogBackupAdd(ctx)
}

// LogBackupDelete 删除备份日志
func (l *LogBackup) LogBackupDelete(ctx context.Context, id int) error {
	var logBackupModel mysqls.LogBackup
	return logBackupModel.LogBackupDelete(ctx, id)
}

// LogBackupList 备份日志列表
func (l *LogBackup) LogBackupList(ctx context.Context, page, limit int) ([]mysqls.LogBackup, int64, error) {
	var logBackupModel mysqls.LogBackup
	return logBackupModel.LogBackupList(ctx, page, limit)
}

// LogBackupRecord 单条备份日志信息
func (l *LogBackup) LogBackupRecord(ctx context.Context, id int) (mysqls.LogBackup, error) {
	var logBackupModel mysqls.LogBackup
	return logBackupModel.LogBackupRecord(ctx, id)
}

// LogBackupNow 日志备份
func (l *LogBackup) LogBackupNow(ctx context.Context) error {
	//获取所有的审计日志
	var logAuditModel mysqls.LogAudit
	logAuditList, err := logAuditModel.LogAuditAll(ctx)
	if err != nil {
		return err
	}
	data := make([][]string, 0)
	data = append(data, []string{"id", "logType", "content", "username", "ip", "createTime", "updateTime"})
	for _, logAudit := range logAuditList {
		data = append(data, []string{
			strconv.Itoa(logAudit.ID),
			strconv.Itoa(logAudit.LogType),
			logAudit.Content,
			logAudit.Username,
			logAudit.Ip,
			logAudit.CreateTime.Format(enums.TimeLayout),
			logAudit.UpdateTime.Format(enums.TimeLayout),
		})
	}
	//备份日志
	if !file.CheckPathExist(enums.LogBackupDir) {
		if err = file.CreateDir(enums.LogBackupDir); err != nil {
			return err
		}
	}
	name := "日志备份文件_" + utils.DatetimeNumberStr()
	filepath := enums.LogBackupDir + name + ".csv"
	err = file.WriteCsv(data, filepath)
	if err != nil {
		return err
	}
	//日志入备份日志库
	var logBackupModel = mysqls.LogBackup{
		Name:       name,
		Path:       filepath,
		CreateTime: time.Now(),
	}
	if err = logBackupModel.LogBackupAdd(ctx); err != nil {
		return err
	}
	//清空所有的审计日志
	if err = logAuditModel.LogAuditEmpty(ctx); err != nil {
		return err
	}
	return nil
}

//查询小于某个创建时间的数据
func (l *LogBackup) GetLogBackupByLtCreateTime(ctx context.Context, creatTime string) ([]mysqls.LogBackup, error) {
	var logBackupModel mysqls.LogBackup
	return logBackupModel.GetLogBackupByLtCreateTime(ctx, creatTime)
}

//删除日志和日志文件
func (l *LogBackup) DelLogBackupMany(ctx context.Context, id int, path string) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	var logBackupModel mysqls.LogBackup
	//删除文件
	if file.CheckPathExist(path) {
		err := file.RemoveFile(path)
		if err != nil {
			return err
		}
	}
	//删除数据
	err := logBackupModel.LogBackupDelete(dCtx, id)
	if err != nil {
		return err
	}
	//提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
