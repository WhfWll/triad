package services

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MailManage struct{}

// SendEmail 发送邮件
func (m *MailManage) SendEmail(recipients, attachmentPaths []string, username, password, smtpHost string, smtpPort int, tcpTimeout time.Duration) error {
	var (
		err     error
		message string
	)

	// 设置发送者的显示名称
	senderDisplayName := "自动化渗透测试平台"
	// 邮件主题和内容
	subject := "渗透任务验证报告"
	body := "hi，all！\r\n渗透任务验证报告结果请查看附件."

	//验证核心参数
	if username == "" || password == "" || smtpHost == "" || smtpPort == 0 || len(recipients) == 0 {
		return errors.New("邮箱服务参数/收件人不能为空")
	}
	if tcpTimeout == 0 {
		tcpTimeout = 5 * time.Second
	}

	// 设置邮件消息
	if len(attachmentPaths) == 0 {
		message = m.CreateMessage(senderDisplayName, username, subject, body, recipients)
	} else {
		message, err = m.CreateMessageWithAttachment(senderDisplayName, username, subject, body, recipients, attachmentPaths)
		if err != nil {
			fmt.Println("设置邮件消息失败: ", err)
			return err
		}
	}

	// 发送邮件（加密端口: 465/587，不加密端口: 25）
	if smtpPort == 25 {
		err = m.SendEmailNormal(username, password, smtpHost, message, smtpPort, tcpTimeout)
	} else {
		err = m.SendEmailWithTls(username, password, smtpHost, message, smtpPort, tcpTimeout)
	}

	if err != nil {
		fmt.Println("邮件发送失败:", err)
		return err
	}

	fmt.Println("邮件发送成功")
	return nil
}

// CreateMessage 创建邮件消息
func (m *MailManage) CreateMessage(senderDisplayName, senderEmail, subject, body string, recipients []string) string {
	// 创建邮件消息（收件人为多个时，用英文的逗号和空格分割）
	message := fmt.Sprintf("From: "+senderDisplayName+" <%s>\r\n", senderEmail)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(recipients, ", "))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "\r\n" + body

	return message
}

// CreateMessageWithAttachment 创建带附件的邮件消息
func (m *MailManage) CreateMessageWithAttachment(senderDisplayName, senderEmail, subject, body string, recipients, attachmentPaths []string) (string, error) {
	// 创建邮件消息
	var msg bytes.Buffer

	// 设置主体
	io.WriteString(&msg, "From: "+senderDisplayName+" <"+senderEmail+">\r\n")
	io.WriteString(&msg, "To: "+strings.Join(recipients, ", ")+"\r\n")
	io.WriteString(&msg, "Subject: "+subject+"\r\n")
	io.WriteString(&msg, "MIME-version: 1.0;\r\n")
	io.WriteString(&msg, "Content-Type: multipart/mixed; boundary=foobar\r\n\r\n")

	// 写入正文
	io.WriteString(&msg, "--foobar\r\n")
	io.WriteString(&msg, "Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
	io.WriteString(&msg, body+"\r\n")

	// 写入附件
	for _, attachmentPath := range attachmentPaths {
		attachmentFile, err := os.Open(attachmentPath)
		if err != nil {
			return "", err
		}
		defer attachmentFile.Close()

		io.WriteString(&msg, "--foobar\r\n")
		io.WriteString(&msg, "Content-Type: application/octet-stream\r\n")
		io.WriteString(&msg, "Content-Disposition: attachment; filename=\""+filepath.Base(attachmentPath)+"\"\r\n\r\n")
		io.Copy(&msg, attachmentFile)
		io.WriteString(&msg, "\r\n")
	}

	// 写入结束标志
	io.WriteString(&msg, "--foobar--\r\n")

	return msg.String(), nil
}

// SendEmailWithTls 安全策略的邮件发送
func (m *MailManage) SendEmailWithTls(username, password, smtpHost, message string, smtpPort int, tcpTimeout time.Duration) error {
	// 设置SMTP客户端配置
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// 设置Dialer，包含超时设置
	dialer := &net.Dialer{
		Timeout: tcpTimeout,
	}

	// 使用DialWithDialer设置连接超时时间
	conn, err := tls.DialWithDialer(dialer, "tcp", smtpHost+":"+strconv.Itoa(smtpPort), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return errors.New("DialWithDialer err: " + err.Error())
	}
	defer conn.Close()

	// 连接到smtp服务器
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return errors.New("NewClient err: " + err.Error())
	}
	defer client.Quit()

	// 认证
	err = client.Auth(auth)
	if err != nil {
		return errors.New("Auth err: " + err.Error())
	}

	// 设置发件人
	err = client.Mail(username)
	if err != nil {
		return errors.New("Set sender err: " + err.Error())
	}

	// 设置收件人
	for _, recipient := range m.GetRecipients(message) {
		err = client.Rcpt(recipient)
		if err != nil {
			return errors.New("Set rcpt err: " + err.Error())
		}
	}

	// 获取数据写入器
	data, err := client.Data()
	if err != nil {
		return errors.New("Write data err: " + err.Error())
	}
	defer data.Close()

	// 写入邮件消息
	_, err = data.Write([]byte(message))
	if err != nil {
		return errors.New("Write message err: " + err.Error())
	}

	return nil
}

// SendEmailNormal 不带安全策略的邮件发送
func (m *MailManage) SendEmailNormal(username, password, smtpHost, message string, smtpPort int, tcpTimeout time.Duration) error {
	// 设置SMTP客户端配置
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// 使用DialTimeout设置连接超时时间
	conn, err := net.DialTimeout("tcp", smtpHost+":"+strconv.Itoa(smtpPort), tcpTimeout)
	if err != nil {
		return errors.New("DialTimeout err: " + err.Error())
	}
	defer conn.Close()

	// 连接到smtp服务器
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return errors.New("Connect smtp service err: " + err.Error())
	}
	defer client.Quit()

	// 开始TLS握手
	err = client.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return errors.New("StartTLS err: " + err.Error())
	}

	// 认证
	err = client.Auth(auth)
	if err != nil {
		return errors.New("Auth err: " + err.Error())
	}

	// 设置发件人
	err = client.Mail(username)
	if err != nil {
		return errors.New("Set sender err: " + err.Error())
	}

	// 设置收件人
	for _, recipient := range m.GetRecipients(message) {
		err = client.Rcpt(recipient)
		if err != nil {
			return errors.New("Set recipient err: " + err.Error())
		}
	}

	// 获取数据写入器
	data, err := client.Data()
	if err != nil {
		return errors.New("Get data writer err: " + err.Error())
	}
	defer data.Close()

	// 写入邮件消息
	_, err = data.Write([]byte(message))
	if err != nil {
		return errors.New("Write email message err: " + err.Error())
	}

	return nil
}

// GetRecipients 获取邮件消息中的收件人
func (m *MailManage) GetRecipients(message string) []string {
	recipients := make([]string, 0)
	lines := strings.Split(message, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "To: ") {
			recipients = strings.Split(line[4:], ", ")
			return recipients
		}
	}
	return recipients
}
