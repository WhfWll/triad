// Package chat
// @Author bcy2007  2025/11/25 17:54
package chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CommonClient struct {
	url       string
	apiKey    string
	model     string
	streaming bool
	reason    bool
}

type CommonRequest struct {
	Model          string         `json:"model"`
	Messages       []AiMessage    `json:"messages"`
	Stream         bool           `json:"stream"`
	EnableThinking bool           `json:"enable_thinking"`
	StreamOptions  map[string]any `json:"stream_options,omitempty"`
}

func NewCommonClient(url, apiKey, model string, streaming, reason bool) *CommonClient {
	return &CommonClient{
		url:       url,
		apiKey:    apiKey,
		model:     model,
		streaming: streaming,
		reason:    reason,
	}
}

func (c *CommonClient) headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + c.apiKey,
	}
}

func (c *CommonClient) generateRequest() interface{} {
	req := CommonRequest{
		Model:          c.model,
		Messages:       nil,
		Stream:         c.streaming,
		EnableThinking: c.reason,
	}
	if c.streaming {
		req.StreamOptions = map[string]any{
			"include_usage": true,
		}
	}
	return req
}

func (c *CommonClient) Run(ctx context.Context, callback func(msgType string, msg any), queries ...AiMessage) error {
	commonReq, _ := c.generateRequest().(CommonRequest)
	commonReq.Messages = make([]AiMessage, 0)
	commonReq.Messages = append(commonReq.Messages, queries...)
	commonReqBody, err := json.Marshal(commonReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(commonReqBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	for k, v := range c.headers() {
		req.Header.Set(k, v)
	}
	client := http.Client{
		Timeout: 120 * time.Second,
	}
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		//return fmt.Errorf("API调用失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
		return errors.New(string(bytes.TrimSpace(body)))
	}
	return c.getResponse(resp, callback)
}

func (c *CommonClient) getResponse(resp *http.Response, callback func(msgType string, msg any)) error {
	fmt.Println(resp)
	if c.streaming {
		return c.getStreamResponse(resp, callback)
	} else {
		return c.getNormalResponse(resp, callback)
	}
}

func (c *CommonClient) getStreamResponse(resp *http.Response, callback func(msgType string, msg any)) error {
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	var (
		fullContent strings.Builder
	)
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过空行和非data行
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		// 提取JSON数据
		jsonData := strings.TrimPrefix(line, "data: ")
		if jsonData == "[DONE]" {
			callback(MsgTypeDone, "done")
			break
		}
		var chunk AiStreamChunk
		if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) > 0 {
			reason := chunk.Choices[0].Delta.ReasoningContent
			content := chunk.Choices[0].Delta.Content
			if reason != "" {
				callback(MsgTypeReason, reason)
			}
			if content != "" {
				callback(MsgTypeContent, content)
				fullContent.WriteString(content)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	callback(MsgTypeFull, fullContent.String())
	return nil
}

func (c *CommonClient) getNormalResponse(resp *http.Response, callback func(msgType string, msg any)) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}
	var result AiCommonResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}
	if len(result.Choices) > 0 {
		if result.Choices[0].Message.ReasoningContent != "" {
			callback(MsgTypeReason, result.Choices[0].Message.ReasoningContent)
		}
		if result.Choices[0].Message.Content != "" {
			callback(MsgTypeContent, result.Choices[0].Message.Content)
		}
	}
	return nil
}
