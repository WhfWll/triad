package invoke

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

type YakitServer struct {
	port   int
	server *WebHookServer

	// handleProgress
	progressHandler func(id string, progress float64)
	logHandler      func(level string, info string)
}

func SetYakitServer_ProgressHandler(h func(id string, progress float64)) func(s *YakitServer) {
	return func(s *YakitServer) {
		s.progressHandler = h
	}
}

func SetYakitServer_LogHandler(h func(level string, info string)) func(s *YakitServer) {
	return func(s *YakitServer) {
		s.logHandler = h
	}
}
func NewYakitServer(port int, opts ...func(server *YakitServer)) *YakitServer {
	if port <= 0 {
		port = GetRandomAvailableTCPPort()
	}

	s := &YakitServer{
		port: port,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.server = NewWebHookServerEx(port, func(data interface{}) {
		switch ret := data.(type) {
		case *http.Request:
			if ret == nil {
				return
			}
			if ret.RemoteAddr != "" {
				fmt.Println("remote addr: %s", ret.RemoteAddr)
			}

			if ret.Body != nil {
				raw, _ := ioutil.ReadAll(ret.Body)
				if raw != nil {
					s.handleRaw(raw)
				}
			}
		}
	})
	return s
}

func (s *YakitServer) Start() {
	s.server.Start()
	return
}

func (s *YakitServer) Addr() string {
	if s.server == nil {
		return ""
	}
	return s.server.Addr()
}

func (s *YakitServer) Shutdown() {
	if s.server == nil {
		return
	}
	s.server.Shutdown()
}

type YakitMessage struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type YakitProgress struct {
	Id       string  `json:"id"`
	Progress float64 `json:"progress"`
}
type YakitLog struct {
	Level     string `json:"level"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

func (s *YakitServer) handleRaw(raw []byte) {
	var msg YakitMessage
	_ = json.Unmarshal(raw, &msg)
	switch strings.ToLower(msg.Type) {
	case "progress", "prog":
		if s.progressHandler == nil {
			return
		}
		var prog YakitProgress
		err := json.Unmarshal(msg.Content, &prog)
		if err != nil {
			fmt.Errorf("unmarshal progress failed: %s", err)
			return
		}
		s.progressHandler(prog.Id, prog.Progress)
	case "log":
		if s.logHandler == nil {
			return
		}
		var logInfo YakitLog
		err := json.Unmarshal(msg.Content, &logInfo)
		if err != nil {
			fmt.Errorf("unmarshal log failed: %s", err)
			return
		}
		s.logHandler(logInfo.Level, logInfo.Data)
	}
}

func GetRandomAvailableTCPPort() int {
RESET:
	randPort := 55000 + rand.Intn(10000)
	if !IsTCPPortOpen("127.0.0.1", randPort) && IsTCPPortAvailable(randPort) {
		return randPort
	} else {
		goto RESET
	}
}

func IsTCPPortOpen(host string, p int) bool {
	dialer := net.Dialer{}
	dialer.Timeout = 10 * time.Second
	conn, err := dialer.Dial("tcp", HostPort(host, p))
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
func IsTCPPortAvailable(p int) bool {
	return IsPortAvailable("0.0.0.0", p)
}
func IsPortAvailable(host string, p int) bool {
	lis, err := net.Listen("tcp", HostPort(host, p))
	if err != nil {
		return false
	}
	_ = lis.Close()
	return true
}
func HostPort(host string, port interface{}) string {
	return fmt.Sprintf("%v:%v", ParseHostToAddrString(host), port)
}
func ParseHostToAddrString(host string) string {
	ip := net.ParseIP(host)
	if ip == nil {
		return host
	}

	if ret := ip.To4(); ret == nil {
		return fmt.Sprintf("[%v]", ip.String())
	}

	return host
}
