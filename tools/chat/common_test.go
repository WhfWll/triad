// Package chat
// @Author bcy2007  2025/11/25 18:32
package chat

import (
	"context"
	"fmt"
	"testing"
)

func TestCommonClient(t *testing.T) {
	var (
		url      = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
		key      = "sk-9efc5a54aac84db48dd486d17e6d7891"
		model    = "qwen-plus"
		ctx      = context.Background()
		callback = func(msgType string, msg any) {
			fmt.Printf("[%s] %v\n", msgType, msg)
		}
	)
	client := NewCommonClient(url, key, model, false, false)
	err := client.Run(ctx, callback, AiMessage{
		Role:    "user",
		Content: "请介绍一下你自己",
	})
	if err != nil {
		t.Error(err)
	}
}
