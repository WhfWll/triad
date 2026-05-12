package services

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

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
	Config     *HostConnConfig
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
	LastUsed   time.Time
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

func connPoolKey(host string, port int, username string) string {
	return fmt.Sprintf("%s:%d@%s", username, port, host)
}

func (m *HostConnManager) GetConnection(ctx context.Context, config *HostConnConfig) (*HostConnection, error) {
	key := connPoolKey(config.Host, config.Port, config.Username)

	m.mu.RLock()
	conn, exists := m.connPool[key]
	m.mu.RUnlock()

	if exists && conn.SSHClient != nil {
		_, _, err := conn.SSHClient.SendRequest("keepalive@openssh.com", true, nil)
		if err == nil {
			conn.LastUsed = time.Now()
			return conn, nil
		}
		log.Warnf("SSH connection to %s expired, reconnecting", key)
		m.Close(key)
	}

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

	conn = &HostConnection{
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

func (m *HostConnManager) ExecuteCommand(ctx context.Context, conn *HostConnection, command string) (string, error) {
	if conn.SSHClient == nil {
		return "", fmt.Errorf("ssh client is nil")
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
