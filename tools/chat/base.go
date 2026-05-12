// Package chat
// @Author bcy2007  2025/11/25 17:43
package chat

import (
	"context"
	"fmt"
)

const (
	MsgTypeReason  = "reason"
	MsgTypeContent = "content"
	MsgTypeDone    = "done"
	MsgTypeFull    = "full"
)

type AiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AiCommonResponse struct {
	Choices []AiResponseMessage `json:"choices"`
}

type AiStreamChunk struct {
	Choices []struct {
		Delta AiResponse `json:"delta"`
	} `json:"choices"`
}

type AiResponseMessage struct {
	Message AiResponse `json:"message"`
}

type AiResponse struct {
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
	Role             string `json:"role"`
}

type AiInterface interface {
	headers() map[string]string
	generateRequest() interface{}
	Run(ctx context.Context, callback func(msgType string, msg any), queries ...AiMessage) error
}

var (
	DefaultCallBack = func(msgType string, msg any) {
		fmt.Printf("[%s] %v\n", msgType, msg)
	}
	NoneCallBack = func(_ string, _ any) {}
)
