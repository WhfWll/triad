package sshutils

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// ExecCommand executes a command on a remote SSH server using password authentication.
func ExecCommand(host string, port string, user string, password string, cmd string) (string, error) {
	client, err := getSSHClient(host, port, user, password)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Execute the command
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		if len(output) > 0 {
			// If there is output but also an error (e.g. exit status 1), return both
			return string(output), fmt.Errorf("command execution failed: %v, output: %s", err, string(output))
		}
		return "", fmt.Errorf("failed to run command: %v", err)
	}

	return string(output), nil
}

// UploadFile uploads a local file to a remote server via SFTP.
func UploadFile(host, port, user, password, localPath, remotePath string) error {
	client, err := getSSHClient(host, port, user, password)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create SFTP client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create sftp client: %v", err)
	}
	defer sftpClient.Close()

	// Open local file
	srcFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer srcFile.Close()

	// Create remote file
	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %v", err)
	}
	defer dstFile.Close()

	// Copy file contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	// Ensure permissions are set (e.g., executable for binaries)
	if err := sftpClient.Chmod(remotePath, 0755); err != nil {
		return fmt.Errorf("failed to chmod remote file: %v", err)
	}

	return nil
}

func getSSHClient(host, port, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, should verify host key
		Timeout:         10 * time.Second,
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", net.JoinHostPort(host, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	return client, nil
}
