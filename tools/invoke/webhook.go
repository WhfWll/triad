package invoke

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"smart/tools/network"
	"time"
)

type WebHookServer struct {
	s    *http.Server
	addr string
}

func NewWebHookServerEx(port int, cb func(data interface{})) *WebHookServer {
	//addr := fmt.Sprintf("127.0.0.1:%v", port)
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer writer.WriteHeader(200)
			//fmt.Printf("webhook met req: %v<-%v", addr, request.RemoteAddr)
			reqBytes, err := network.HttpDumpWithBody(request, true)
			if err != nil {
				fmt.Errorf("webhook read request failed: %s", err)
				return
			}

			requestIns, err := network.ReadHTTPRequestEx(bufio.NewReader(bytes.NewBuffer(reqBytes)), true)
			if err != nil {
				fmt.Errorf("re-build webhook request failed: %s", err)
				return
			}

			originUrl := requestIns.URL
			target := fmt.Sprintf("http://%v%v", addr, requestIns.URL.Path)
			requestIns.URL, _ = url.Parse(target)
			if requestIns.URL != nil && requestIns.Body == nil {
				requestIns.Body = http.NoBody
				requestIns.GetBody = func() (io.ReadCloser, error) {
					return http.NoBody, nil
				}
			}

			if requestIns.URL == nil {
				requestIns.URL = originUrl
			}
			cb(requestIns)
		}),
		TLSConfig:         nil,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	server.SetKeepAlivesEnabled(true)
	return &WebHookServer{
		s:    server,
		addr: addr,
	}
}

func (w *WebHookServer) Start() {
	go func() {
		err := w.s.ListenAndServe()
		if err != nil {

		}
	}()
}

func TimeoutContext(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}

func (w *WebHookServer) Shutdown() {
	_ = w.s.Shutdown(TimeoutContext(1 * time.Second))
}

func (w *WebHookServer) Addr() string {
	return fmt.Sprintf("http://%v/webhook", w.addr)
}
