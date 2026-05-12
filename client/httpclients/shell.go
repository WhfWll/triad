package httpclients

import (
	"context"
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
	"io"
	"strconv"
	"strings"
)

/***************************** Shell - 执行命令 ********************************/
// Shell命令执行请求
type ShellExecuteReq struct {
	SessionId string `json:"session_id" form:"session_id" binding:"required"` // 会话ID
	Command   string `json:"command" form:"command" binding:"required"`       // 执行的命令
}

// Shell命令执行响应
type ShellExecuteRes struct {
	Success bool   `json:"success"` // 执行是否成功
	Message string `json:"message"` // 消息
	Data    struct {
		Command   string `json:"command"`    // 执行的命令
		Result    string `json:"result"`     // 命令执行结果
		SessionId string `json:"session_id"` // 会话ID
		Timestamp string `json:"timestamp"`  // 时间戳
	} `json:"data"` // 响应数据
	Error     string `json:"error"`      // 错误信息
	Timestamp string `json:"timestamp"`  // 响应时间戳
	RequestId string `json:"request_id"` // 请求ID
}

// ExecuteShellCommand Shell命令执行
func ExecuteShellCommand(ctx context.Context, req ShellExecuteReq) (res ShellExecuteRes, err error) {
	h, err := httpclient.NewHttpSend("service_shell", "/api/sessions/execute") // 第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	// 参数获取
	param := make(map[string]interface{})
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}
	err = json.Unmarshal(reqByte, &param)
	if err != nil {
		return
	}

	// 绑定参数
	h.SetSendType("json")
	h.SetBody(param) // 设置POST请求参数

	// 发起请求
	h.SetHeader(map[string]interface{}{"Content-Type": "application/json", "Connection": "close"})
	response, err := h.Post() // 发送POST请求
	if err != nil {
		return
	}

	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

/***************************** Shell - 文件列表 ********************************/
// 文件信息结构体
type FileInfo struct {
	IsDir   bool   `json:"is_dir"`   // 是否为目录
	ModTime string `json:"mod_time"` // 修改时间
	Name    string `json:"name"`     // 文件名
	Size    int64  `json:"size"`     // 文件大小
}

// Shell文件列表请求
type ShellFileListReq struct {
	SessionId string `json:"session_id" form:"session_id" binding:"required"` // 会话ID
	Path      string `json:"path" form:"path" binding:"required"`             // 路径
}

// Shell文件列表响应
type ShellFileListRes struct {
	Success bool   `json:"success"` // 执行是否成功
	Message string `json:"message"` // 消息
	Data    struct {
		Data struct {
			Files []FileInfo `json:"files"` // 文件列表
		} `json:"data"`
		Path      string `json:"path"`       // 路径
		SessionId string `json:"session_id"` // 会话ID
		Success   bool   `json:"success"`    // 成功标识
		Timestamp string `json:"timestamp"`  // 时间戳
	} `json:"data"` // 响应数据
	Error     string `json:"error"`      // 错误信息
	Timestamp string `json:"timestamp"`  // 响应时间戳
	RequestId string `json:"request_id"` // 请求ID
}

// ListShellFiles Shell文件列表获取
func ListShellFiles(ctx context.Context, req ShellFileListReq) (res ShellFileListRes, err error) {
	h, err := httpclient.NewHttpSend("service_shell", "/api/files/list") // 第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	// 设置查询参数
	queryParams := map[string]interface{}{
		"session_id": req.SessionId,
		"path":       req.Path,
	}
	h.GetUrlBuild(queryParams)

	// 发起请求
	h.SetHeader(map[string]interface{}{"Connection": "close"})
	response, err := h.Get() // 发送GET请求
	if err != nil {
		return
	}

	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

/***************************** Shell - 会话列表 ********************************/
// Shell会话列表请求
type ShellSessionListReq struct {
	TimeFilter string `json:"time_filter" form:"time_filter"` // 时间过滤参数，默认60s
}

// 会话统计信息
type SessionStats struct {
	CommandsExecuted int64  `json:"commands_executed"` // 执行的命令数
	BytesReceived    int64  `json:"bytes_received"`    // 接收的字节数
	BytesSent        int64  `json:"bytes_sent"`        // 发送的字节数
	Uptime           int64  `json:"uptime"`            // 运行时间（纳秒）
	LastCommand      string `json:"last_command"`      // 最后执行的命令
	LastCommandTime  string `json:"last_command_time"` // 最后命令执行时间
	ErrorCount       int    `json:"error_count"`       // 错误计数
}

// 会话元数据
type SessionMetadata struct {
	LastActive  string `json:"last_active"`  // 最后活跃时间
	SessionId   string `json:"session_id"`   // 会话ID
	SessionType string `json:"session_type"` // 会话类型
	StartTime   string `json:"start_time"`   // 开始时间
	Status      string `json:"status"`       // 状态
	Target      string `json:"target"`       // 目标
	Uptime      string `json:"uptime"`       // 运行时间字符串
}

// 会话信息
type SessionInfo struct {
	Id           string          `json:"id"`           // 会话ID
	Type         string          `json:"type"`         // 会话类型
	Status       string          `json:"status"`       // 状态
	Target       string          `json:"target"`       // 目标
	StartTime    string          `json:"start_time"`   // 开始时间
	LastActive   string          `json:"last_active"`  // 最后活跃时间
	Stats        SessionStats    `json:"stats"`        // 统计信息
	Metadata     SessionMetadata `json:"metadata"`     // 元数据
	Capabilities interface{}     `json:"capabilities"` // 能力信息
}

// Shell会话列表响应
type ShellSessionListRes struct {
	Success   bool          `json:"success"`    // 执行是否成功
	Message   string        `json:"message"`    // 消息
	Data      []SessionInfo `json:"data"`       // 会话列表数据
	Error     string        `json:"error"`      // 错误信息
	Timestamp string        `json:"timestamp"`  // 响应时间戳
	RequestId string        `json:"request_id"` // 请求ID
}

// ListShellSessions Shell会话列表获取
func ListShellSessions(ctx context.Context, req ShellSessionListReq) (res ShellSessionListRes, err error) {
	h, err := httpclient.NewHttpSend("service_shell", "/api/sessions") // 第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	// 设置查询参数，默认时间过滤为60s
	timeFilter := req.TimeFilter
	if timeFilter == "" {
		timeFilter = "60s"
	}
	queryParams := map[string]interface{}{
		"time_filter": timeFilter,
	}
	h.GetUrlBuild(queryParams)

	// 发起请求
	h.SetHeader(map[string]interface{}{"Connection": "close"})
	response, err := h.Get() // 发送GET请求
	if err != nil {
		return
	}

	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	return
}

/***************************** Shell - 文件下载 ********************************/
// Shell文件下载请求
type ShellFileDownloadReq struct {
	SessionId  string `json:"session_id" form:"session_id" binding:"required"`   // 会话ID
	RemotePath string `json:"remote_path" form:"remote_path" binding:"required"` // 远程文件路径
}

// Shell文件下载响应
type ShellFileDownloadRes struct {
	Content     []byte            `json:"-"`            // 文件内容（二进制数据）
	Headers     map[string]string `json:"headers"`      // 响应头信息
	FileName    string            `json:"filename"`     // 文件名
	ContentType string            `json:"content_type"` // 内容类型
	Size        int64             `json:"size"`         // 文件大小
}

// DownloadShellFile Shell文件下载
func DownloadShellFile(ctx context.Context, req ShellFileDownloadReq) (res ShellFileDownloadRes, err error) {
	h, err := httpclient.NewHttpSend("service_shell", "/api/files/download") // 第一个参数是配置文件中client中的key，第二个参数uri
	if err != nil {
		return
	}

	// 设置查询参数
	queryParams := map[string]interface{}{
		"session_id":  req.SessionId,
		"remote_path": req.RemotePath,
	}
	h.GetUrlBuild(queryParams)

	// 发起请求
	h.SetHeader(map[string]interface{}{"Connection": "close"})
	// 使用Send方法获取完整的http.Response对象
	resp, err := h.Send("GET")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 获取响应头信息
	headerMap := h.GetHeader(resp)
	res.Headers = headerMap

	// 解析Content-Disposition获取文件名
	if contentDisposition, ok := res.Headers["Content-Disposition"]; ok {
		// 简单解析filename，实际可能需要更复杂的解析
		if idx := strings.Index(contentDisposition, "filename="); idx != -1 {
			filename := contentDisposition[idx+9:] // 跳过"filename="
			filename = strings.Trim(filename, `"`)
			res.FileName = filename
		}
	}

	// 获取内容类型
	if contentType, ok := res.Headers["Content-Type"]; ok {
		res.ContentType = contentType
	}

	// 获取文件大小
	if contentLength, ok := res.Headers["Content-Length"]; ok {
		if size, parseErr := strconv.ParseInt(contentLength, 10, 64); parseErr == nil {
			res.Size = size
		}
	}

	// 文件内容就是响应体的二进制数据
	res.Content = body

	return
}
