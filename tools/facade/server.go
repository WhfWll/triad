// Package facade
// @Author bcy2007  2025/12/22 16:26
package facade

import (
	"context"
	"fmt"
	"smart/tools/facade/base"
	"strings"
	"sync"
	"time"
)

type facadeMessage struct {
	FacadeType string
	RemoteAddr string
	Token      string
	Response   string
}

var (
	lock           = new(sync.RWMutex)
	localServer    *base.FacadeServer
	serverCancel   func()
	historyMessage []facadeMessage
	errInfo        string
)

func NewFacadeServer(host string, port int, timeout int) {
	lock.Lock()
	defer lock.Unlock()
	if localServer != nil {
		localServer.CancelServe()
	}
	if serverCancel != nil {
		serverCancel()
	}
	historyMessage = make([]facadeMessage, 0)
	errInfo = ""
	server := base.NewFacadeServer(host, port)
	server.OnHandle(func(n *base.Notification) {
		lock.Lock()
		defer lock.Unlock()
		historyMessage = append(historyMessage, facadeMessage{
			FacadeType: n.Type,
			RemoteAddr: n.RemoteAddr,
			Token:      n.Token,
			Response:   n.ResponseInfo,
		})
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	localServer = server
	serverCancel = cancel
	// 启动 server 并监听超时或取消
	go func(s *base.FacadeServer, cancelFunc func()) {
		err := s.ServeWithContext(ctx)
		if err != nil {
			if !strings.Contains(err.Error(), "closed") && !strings.Contains(err.Error(), "deadline exceeded") {
				errInfo = fmt.Sprintf("FacadeServer ServeWithContext 启动失败: %v", err)
			}
			// 如果启动失败，立即取消 context，让监听超时的 goroutine 也能及时退出
			cancelFunc()
			// 清理 server 引用，避免状态不一致
			lock.Lock()
			defer lock.Unlock()
			if localServer == s {
				localServer = nil
				serverCancel = nil
			}
		}
	}(localServer, cancel)
	// 监听超时或取消，把对应的 server 引用清空，避免旧实例残留
	go func(s *base.FacadeServer) {
		<-ctx.Done()
		fmt.Println("close!")
		lock.Lock()
		defer lock.Unlock()
		if localServer == s {
			localServer = nil
			serverCancel = nil
		}
	}(localServer)
}

func CloseFacadeServer() {
	lock.Lock()
	defer lock.Unlock()
	if localServer != nil {
		localServer.CancelServe()
	}
	if serverCancel != nil {
		serverCancel()
	}
	localServer = nil
	serverCancel = nil
	errInfo = ""
}

func ClearMessage() {
	lock.Lock()
	defer lock.Unlock()
	historyMessage = make([]facadeMessage, 0)
	errInfo = ""
}

func ReadMessage(page, size int) ([]facadeMessage, int) {
	lock.Lock()
	defer lock.Unlock()
	return getPageSize(historyMessage, page, size)
}

func GetFacadeServerStatus() (bool, string) {
	lock.Lock()
	defer lock.Unlock()
	if localServer == nil {
		return false, errInfo
	}
	return true, fmt.Sprintf("%s:%d", localServer.Host, localServer.Port)
}

func getPageSize[T any](data []T, page, size int) ([]T, int) {
	var (
		total = len(data)
	)
	if page <= 0 || size <= 0 {
		return nil, total
	}
	// 从后向前计算：第1页是最新的数据（数组末尾）
	// end 从数组末尾开始计算
	end := total - (page-1)*size
	start := total - page*size
	if end <= 0 {
		return nil, total
	}
	if start < 0 {
		start = 0
	}
	// 直接反向复制，避免先切片再反转，性能更好
	result := make([]T, end-start)
	for i := 0; i < len(result); i++ {
		result[i] = data[end-1-i]
	}
	return result, total
}
