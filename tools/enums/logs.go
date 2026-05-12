package enums

//审计日志类型
const (
	LogAuditTypeLogin   = 1
	LogAuditTypeOperate = 2
	LogAuditTypeWarn    = 3
)

//日志周期备份开启状态
const (
	LogBackupCycleOpen  = 1
	LogBackupCycleClose = 2
)

// LogBackupDir 日志备份目录
const (
	LogBackupDir = "backup/logs/"
)

// LogAuditTypeEnum 获取审计日志类型的映射
func LogAuditTypeEnum() map[int]string {
	typeEnum := map[int]string{
		LogAuditTypeLogin:   "登录日志",
		LogAuditTypeOperate: "操作日志",
		LogAuditTypeWarn:    "告警日志",
	}
	return typeEnum
}

// GetLogAuditType 获取审计日志类型的值
func GetLogAuditType(k int) string {
	typeEnum := LogAuditTypeEnum()
	if value, ok := typeEnum[k]; ok {
		return value
	}
	return ""
}

// GetLogAuditTypeArray 获取审计日志类型的数组
func GetLogAuditTypeArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{
		{
			Value: LogAuditTypeLogin,
			Label: GetLogAuditType(LogAuditTypeLogin),
		},
		{
			Value: LogAuditTypeOperate,
			Label: GetLogAuditType(LogAuditTypeOperate),
		},
		{
			Value: LogAuditTypeWarn,
			Label: GetLogAuditType(LogAuditTypeWarn),
		},
	}
	return result
}
