package services

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"smart/tools/enums"
	"strings"
	"sync"
	"time"

	"github.com/masterzen/winrm"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type HostConnConfig struct {
	Host               string
	Port               int
	Username           string
	Password           string
	PrivateKey         string
	OSType             int
	// Transport 有效值：HostTransportSSH、HostTransportWinRM；0 表示由 GetConnection 按 OSType 推断（兼容旧调用）
	Transport int
	Timeout            time.Duration
	WinRM_PORT         int
	UseHTTPS           bool
	InsecureSkipVerify bool
}

type HostConnManager struct {
	mu       sync.RWMutex
	connPool map[string]*HostConnection
}

type HostConnection struct {
	Config      *HostConnConfig
	SSHClient   *ssh.Client
	SFTPClient  *sftp.Client
	WinRMClient *winrm.Client
	LastUsed    time.Time
}

var globalConnManager *HostConnManager
var connManagerOnce sync.Once

func GetHostConnManager() *HostConnManager {
	connManagerOnce.Do(func() {
		globalConnManager = &HostConnManager{
			connPool: make(map[string]*HostConnection),
		}
	})
	return globalConnManager
}

func connPoolKey(host string, port int, username string, transport int) string {
	return fmt.Sprintf("t%d|%s:%d@%s", transport, username, port, host)
}

func (m *HostConnManager) GetConnection(ctx context.Context, config *HostConnConfig) (*HostConnection, error) {
	transport := config.Transport
	if transport != enums.HostTransportSSH && transport != enums.HostTransportWinRM {
		if config.OSType == enums.BaselineOSTypeWindows {
			transport = enums.HostTransportWinRM
		} else {
			transport = enums.HostTransportSSH
		}
		config.Transport = transport
	}

	if transport == enums.HostTransportWinRM {
		if config.Port <= 0 {
			if config.UseHTTPS {
				config.Port = 5986
			} else {
				config.Port = 5985
			}
		}
	} else {
		if config.Port <= 0 {
			config.Port = 22
		}
	}

	key := connPoolKey(config.Host, config.Port, config.Username, transport)

	m.mu.RLock()
	conn, exists := m.connPool[key]
	m.mu.RUnlock()

	if exists {
		if conn.SSHClient != nil {
			_, _, err := conn.SSHClient.SendRequest("keepalive@openssh.com", true, nil)
			if err == nil {
				conn.LastUsed = time.Now()
				return conn, nil
			}
			log.Warnf("SSH connection to %s expired, reconnecting", key)
		} else if conn.WinRMClient != nil {
			conn.LastUsed = time.Now()
			return conn, nil
		}
		m.Close(key)
	}

	if transport == enums.HostTransportWinRM {
		return m.createWinRMConnection(ctx, config, key)
	}

	return m.createSSHConnection(ctx, config, key)
}

func (m *HostConnManager) createSSHConnection(ctx context.Context, config *HostConnConfig, key string) (*HostConnection, error) {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	if config.Password != "" {
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(config.Password)}
	} else if config.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("parse private key failed: %v", err)
		}
		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		return nil, fmt.Errorf("no auth method provided (password or private key)")
	}

	addr := net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port))
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh dial failed: %v", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Warnf("sftp init failed (non-fatal): %v", err)
	}

	conn := &HostConnection{
		Config:     config,
		SSHClient:  sshClient,
		SFTPClient: sftpClient,
		LastUsed:   time.Now(),
	}

	m.mu.Lock()
	m.connPool[key] = conn
	m.mu.Unlock()

	return conn, nil
}

func (m *HostConnManager) createWinRMConnection(ctx context.Context, config *HostConnConfig, key string) (*HostConnection, error) {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	winrmPort := config.Port
	if winrmPort == 0 {
		if config.UseHTTPS {
			winrmPort = 5986
		} else {
			winrmPort = 5985
		}
	}

	endpoint := winrm.NewEndpoint(
		config.Host,
		winrmPort,
		config.UseHTTPS,
		config.InsecureSkipVerify,
		nil,
		nil,
		nil,
		timeout,
	)

	winrmClient, err := winrm.NewClient(endpoint, config.Username, config.Password)
	if err != nil {
		return nil, fmt.Errorf("winrm client create failed: %v", err)
	}

	conn := &HostConnection{
		Config:      config,
		WinRMClient: winrmClient,
		LastUsed:    time.Now(),
	}

	m.mu.Lock()
	m.connPool[key] = conn
	m.mu.Unlock()

	return conn, nil
}

func (m *HostConnManager) ExecuteCommand(ctx context.Context, conn *HostConnection, command string) (string, error) {
	if conn.WinRMClient != nil {
		return m.executeWinRMCommand(ctx, conn, command)
	}

	if conn.SSHClient == nil {
		return "", fmt.Errorf("no valid connection (SSH or WinRM)")
	}

	session, err := conn.SSHClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("create ssh session failed: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("execute command failed: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

func (m *HostConnManager) executeWinRMCommand(ctx context.Context, conn *HostConnection, command string) (string, error) {
	var stdout, stderr bytes.Buffer

	done := make(chan error)
	go func() {
		defer close(done)
		_, err := conn.WinRMClient.Run(command, &stdout, &stderr)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-done:
		output := stdout.String()
		if stderr.Len() > 0 {
			output += "\nSTDERR: " + stderr.String()
		}
		if err != nil {
			return output, fmt.Errorf("winrm execute command failed: %v", err)
		}
		return output, nil
	}
}

func (m *HostConnManager) ExecuteCommands(ctx context.Context, conn *HostConnection, commands []string) (map[string]string, error) {
	results := make(map[string]string, len(commands))
	for _, cmd := range commands {
		output, err := m.ExecuteCommand(ctx, conn, cmd)
		if err != nil {
			results[cmd] = fmt.Sprintf("ERROR: %v", err)
		} else {
			results[cmd] = output
		}
	}
	return results, nil
}

func (m *HostConnManager) Close(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, exists := m.connPool[key]; exists {
		if conn.SFTPClient != nil {
			conn.SFTPClient.Close()
		}
		if conn.SSHClient != nil {
			conn.SSHClient.Close()
		}
		delete(m.connPool, key)
	}
}

// HostTestConnResult 连接测试结果。
type HostTestConnResult struct {
	OK            bool
	Message       string
	TransportName string
	Detail        string
}

// TestConnection 验证 SSH/WinRM 凭据是否可用（测试后关闭连接，不长期占用连接池）。
func (m *HostConnManager) TestConnection(ctx context.Context, config *HostConnConfig) *HostTestConnResult {
	if config == nil || strings.TrimSpace(config.Host) == "" {
		return &HostTestConnResult{OK: false, Message: "请填写目标主机"}
	}
	if strings.TrimSpace(config.Username) == "" {
		return &HostTestConnResult{OK: false, Message: "请填写用户名"}
	}

	transport := config.Transport
	if transport != enums.HostTransportSSH && transport != enums.HostTransportWinRM {
		if config.OSType == enums.BaselineOSTypeWindows {
			transport = enums.HostTransportWinRM
		} else {
			transport = enums.HostTransportSSH
		}
		config.Transport = transport
	}

	if transport == enums.HostTransportWinRM {
		if config.Password == "" {
			return &HostTestConnResult{OK: false, Message: "WinRM 需要填写密码"}
		}
	} else if config.Password == "" && config.PrivateKey == "" {
		return &HostTestConnResult{OK: false, Message: "请填写 SSH 密码或私钥"}
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}
	config.Timeout = timeout

	testCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := m.GetConnection(testCtx, config)
	if err != nil {
		return &HostTestConnResult{OK: false, Message: err.Error()}
	}
	defer m.closeByConfig(config)

	transportName := enums.BaselineEnum.GetHostTransportName(config.Transport)
	cmd := "echo triad_ok"
	if transport == enums.HostTransportWinRM {
		cmd = "hostname"
	}
	output, err := m.ExecuteCommand(testCtx, conn, cmd)
	if err != nil {
		return &HostTestConnResult{OK: false, Message: err.Error(), TransportName: transportName}
	}

	port := config.Port
	if port <= 0 {
		if transport == enums.HostTransportWinRM {
			if config.UseHTTPS {
				port = 5986
			} else {
				port = 5985
			}
		} else {
			port = 22
		}
	}
	detail := strings.TrimSpace(output)
	if detail == "" {
		detail = "命令执行成功"
	}
	return &HostTestConnResult{
		OK:            true,
		Message:       fmt.Sprintf("连接成功（%s · %s:%d）", transportName, config.Host, port),
		TransportName: transportName,
		Detail:        detail,
	}
}

func (m *HostConnManager) closeByConfig(config *HostConnConfig) {
	transport := config.Transport
	if transport != enums.HostTransportSSH && transport != enums.HostTransportWinRM {
		if config.OSType == enums.BaselineOSTypeWindows {
			transport = enums.HostTransportWinRM
		} else {
			transport = enums.HostTransportSSH
		}
	}
	port := config.Port
	if transport == enums.HostTransportWinRM {
		if port <= 0 {
			if config.UseHTTPS {
				port = 5986
			} else {
				port = 5985
			}
		}
	} else if port <= 0 {
		port = 22
	}
	m.Close(connPoolKey(config.Host, port, config.Username, transport))
}

func (m *HostConnManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, conn := range m.connPool {
		if conn.SFTPClient != nil {
			conn.SFTPClient.Close()
		}
		if conn.SSHClient != nil {
			conn.SSHClient.Close()
		}
		delete(m.connPool, key)
	}
}

func (m *HostConnManager) CleanIdleConnections(maxIdle time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for key, conn := range m.connPool {
		if now.Sub(conn.LastUsed) > maxIdle {
			if conn.SFTPClient != nil {
				conn.SFTPClient.Close()
			}
			if conn.SSHClient != nil {
				conn.SSHClient.Close()
			}
			delete(m.connPool, key)
		}
	}
}
