package services

import (
	"context"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

type LogAudit struct {
}

// LogAuditTypeEnum 审计日志类型的数组
func (l *LogAudit) LogAuditTypeEnum() interface{} {
	typeEnum := enums.GetLogAuditTypeArray()
	return typeEnum
}

// LogAuditList 审计日志列表
func (l *LogAudit) LogAuditList(ctx context.Context, page, limit, logType int, search, startTime, endTime string) ([]mysqls.LogAudit, int64, error) {
	var logAuditModel mysqls.LogAudit
	return logAuditModel.LogAuditList(ctx, page, limit, logType, search, startTime, endTime)
}

// LogAuditEmpty 清空审计日志
func (l *LogAudit) LogAuditEmpty(ctx context.Context) error {
	var logAuditModel mysqls.LogAudit
	return logAuditModel.LogAuditEmpty(ctx)
}

// LogAuditAdd 新增审计日志
func (l *LogAudit) LogAuditAdd(ctx context.Context, logType int, content, username, ip string) error {
	var logAuditModel = mysqls.LogAudit{
		LogType:    logType,
		Content:    content,
		Username:   username,
		Ip:         ip,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	return logAuditModel.LogAuditAdd(ctx)
}

// LogAuditAll 审计日志所有的数据
func (l *LogAudit) LogAuditAll(ctx context.Context) ([]mysqls.LogAudit, error) {
	var logAuditModel mysqls.LogAudit
	return logAuditModel.LogAuditAll(ctx)
}
