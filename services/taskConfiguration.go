package services

import (
	"context"
	"smart/models/mysqls"
)

type TaskConfiguration struct {
}

// 枚举信息 - 执行类型
func (t *TaskConfiguration) GetConfigInfoByIdAndConfigKey(ctx context.Context, templateId int, configKey string) (mysqls.TaskConfiguration, error) {
	var taskConfiguration mysqls.TaskConfiguration
	return taskConfiguration.GetTaskConfigurationByIdAndConfigKey(ctx, templateId, configKey)
}
