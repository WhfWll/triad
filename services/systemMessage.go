package services

import (
	"context"
	"smart/models/mysqls"
	"time"
)

type SystemMessage struct{}

// GetSystemMessageList 获取系统消息列表
func (s *SystemMessage) GetSystemMessageList(ctx context.Context, page, limit, messageType int, search, startTime, endTime string, uid, role int) ([]mysqls.SystemMessage, int64, error) {
	var systemMessageModel mysqls.SystemMessage
	return systemMessageModel.GetSystemMessageList(ctx, page, limit, messageType, search, startTime, endTime, uid, role)
}

// SystemMessageAdd 添加消息
func (s *SystemMessage) SystemMessageAdd(ctx context.Context, content string, messageType, userId, status int) error {
	var systemMessageModel = mysqls.SystemMessage{
		Content:    content,
		Type:       messageType,
		UserId:     userId,
		Status:     status,
		CreateTime: time.Now(),
	}
	return systemMessageModel.SystemMessageAdd(ctx)
}
