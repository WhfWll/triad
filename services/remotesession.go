package services

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"smart/client/httpclients"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/invoke"
	"strconv"
	"strings"
	"time"
)

type RemoteSession struct {
}

// GetRemoteSessionList 返回列表
func (re *RemoteSession) GetRemoteSessionList(ctx context.Context, taskID, targetId, page, limit, riskType int, search string) ([]mysqls.RemoteSession, int64, error) {
	var remoteSession = &mysqls.RemoteSession{
		TaskID: taskID,
	}
	res, count, err := remoteSession.GetRemoteSessionList(ctx, taskID, targetId, page, limit, search)
	if err != nil {
		return res, 0, err
	}
	return res, count, nil
}

// DelGetRemoteSessions 批量删除
func (re *RemoteSession) DelGetRemoteSessions(ctx context.Context, ids string) error {
	var remoteSessionModel mysqls.RemoteSession
	return remoteSessionModel.DeleteRemoteSessionIds(ctx, strings.Split(ids, ","))
}

// GetRemoteSessionInfo 返回会话详情
func (re *RemoteSession) GetRemoteSessionInfo(ctx context.Context, id int) (mysqls.RemoteSession, error) {
	var remoteSessionModel = &mysqls.RemoteSession{
		ID: id,
	}
	res, err := remoteSessionModel.GetRemoteSession(ctx)
	if err != nil {
		return mysqls.RemoteSession{}, err
	}
	return res, nil
}

// AddRemoteSessionInfo 添加会话信息
func (re *RemoteSession) AddRemoteSessionInfo(ctx context.Context, taskID, targetId, evidenceType int, targetURL, sessionId, remoteControl string, detailMap map[string]interface{}) error {
	log.Println("添加一条远控信息", taskID, targetId, evidenceType, targetURL, detailMap)
	var route, sessionKey string
	route = "0"
	if evidenceType == enums.VulScriptEvidenceTypeWebShell {
		detailMap["shellType"] = "webshell"
		sessionId = strconv.FormatInt(time.Now().Unix(), 10)
		route = "1"
		remoteControl = "webshell"
		sessionKey = ""
	}
	detailByte, _ := json.Marshal(detailMap)
	var remoteSessionModel = &mysqls.RemoteSession{
		Estate:          "valid",
		TaskID:          taskID,
		TargetID:        targetId,
		SessionID:       sessionId,
		TargetURL:       targetURL,
		Route:           route,
		RemoteControl:   remoteControl,
		Detail:          string(detailByte),
		Status:          1,
		DownloadedFiles: "",
		UserID:          0,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		SessionKey:      sessionKey,
	}
	if err := remoteSessionModel.AddRemoteSession(ctx); err != nil {
		log.Println("AddRemoteSessionInfo AddRemoteSession error: ", err)
		log.Error("AddRemoteSessionInfo AddRemoteSession error: ", err)
		return err
	}
	return nil
}

// GetRemoteSessionDownLoadFileList 获取远程文件列表
func (re *RemoteSession) GetRemoteSessionDownLoadFileList(ctx context.Context, remoteSessionID int) (downloadFiles []map[string]string, err error) {
	var remoteSessionModel = &mysqls.RemoteSession{ID: remoteSessionID}
	res, err := remoteSessionModel.GetRemoteSession(ctx)
	if err != nil {
		return
	}
	if len(res.DownloadedFiles) != 0 {
		if err = json.Unmarshal([]byte(res.DownloadedFiles), &downloadFiles); err != nil {
			return
		}
	}
	return downloadFiles, nil
}

// UpdateRemoteSessionDownLoadFileList 更新远程文件列表
func (re *RemoteSession) UpdateRemoteSessionDownLoadFileList(ctx context.Context, remoteSessionID int, newDownloadFileByte string) error {
	var remoteSessionModel = &mysqls.RemoteSession{ID: remoteSessionID, DownloadedFiles: newDownloadFileByte, UpdateTime: time.Now()}
	return remoteSessionModel.UpdateRemoteSession(ctx)
}

// GetShellCmdMany 获取抓取信息命令
func (re *RemoteSession) GetShellCmdMany(captureType int, CaptureInfoIds string) []string {
	var reuslt = make([]string, 0)
	captureInfoArray := strings.Split(CaptureInfoIds, ",")
	for i := 0; i < len(captureInfoArray); i++ {
		tmpId, _ := strconv.Atoi(captureInfoArray[i])
		if captureType == enums.CaptureTypeWindows {
			reuslt = append(reuslt, enums.TaskEvidenceCaptureInfoWindowsExecMap[tmpId])
			continue
		}
		reuslt = append(reuslt, enums.TaskEvidenceCaptureInfoCentosExecMap[tmpId])
	}

	return reuslt
}

// ExceShellMany 批量执行shell命令
func (re *RemoteSession) ExceShellMany(sessionkey string, cmdShellArray []string) string {
	var result string
	for i := 0; i < len(cmdShellArray); i++ {
		result += cmdShellArray[i] + "\n"
		execResult, err := httpclients.ExecShellCmd(httpclients.ExecShellCmdReq{
			CmdKey: sessionkey,
			Cmd:    cmdShellArray[i],
		})
		if err != nil {
			result += err.Error() + "\n"
			continue
		}
		result += execResult.Data.Result + "\n"
	}
	return result
}

// FindFileWindows msf查找文件
/*返回格式化：
Found 1 result...
    c:\Windows\System32\zh-CN\zipfldr.dll.mui (13312 bytes)
*/
func (re *RemoteSession) FindFileWindows(sessionkey string, fileName, filePath string) ([]RemoteSessionFile, error) {
	var result = make([]RemoteSessionFile, 0)
	execResult, err := httpclients.ExecShellCmd(httpclients.ExecShellCmdReq{
		CmdKey: sessionkey,
		Cmd:    "search -d " + filePath + " -f " + fileName,
	})
	if err != nil {
		return result, err
	}
	if len(execResult.Data.Result) > 0 {
		newString := strings.Split(execResult.Data.Result, "\n")
		if len(newString) > 1 && strings.Contains(newString[0], "Found") && strings.Contains(newString[0], "result") {
			for i := 1; i < len(newString); i++ {
				tmpString := strings.Split(newString[i], "(")
				if len(tmpString) > 1 {
					var tmp RemoteSessionFile
					tmp.FileSize = strings.Replace(tmpString[1], ")", "", 1)
					tmp.FilePath = strings.Replace(strings.TrimSpace(tmpString[0]), "\\", "/", -1) //window路径需要转换
					filePathArray := strings.Split(tmp.FilePath, "/")
					tmp.FileName = filePathArray[len(filePathArray)-1]
					result = append(result, tmp)
				}
			}
		}
	}
	return result, nil
}

//DownloadFileWindows 下载文件
/*
下载到 /home/msf/download/会话id/ 路径会挂载到宿主机上/opt/laozhi/msf上
下载成功显示：
[-] stdapi_fs_stat: Operation failed: The system cannot find the file specified.
下载失败显示：
Downloading: c:/Windows/System32/zh-CN/zipfldr.dll.mui -> /home/msf/download/decision:remote:shell:1695004129/zipfldr.dll.mui
*/
func (re *RemoteSession) DownloadFileWindows(sessionkey string, filePath []RemoteSessionFile) ([]RemoteSessionFile, error) {
	var result = make([]RemoteSessionFile, 0)
	for i := 0; i < len(filePath); i++ {
		downPathString := "/home/msf/download/" + sessionkey + "/"                              //下载路径
		hostPathString := enums.SystemUpgradeProjectDir + "msf/download/" + sessionkey + "/" + filePath[i].FileName //宿主机路径

		execResult, err := httpclients.ExecShellCmd(httpclients.ExecShellCmdReq{
			CmdKey: sessionkey,
			Cmd:    "download " + filePath[i].FilePath + " " + downPathString,
		})
		if err != nil {
			continue
		}
		if len(execResult.Data.Result) > 0 && strings.Contains(execResult.Data.Result, "Downloading") {
			filePath[i].FilePath = hostPathString
			result = append(result, filePath[i])
		}
	}
	return result, nil
}

// SaveDownloadFiles 保存文件下载数据
func (re *RemoteSession) SaveDownloadFiles(ctx context.Context, RemoteSessionId int, downloadedFiles string, filePath []RemoteSessionFile) error {
	var rsMysql mysqls.RemoteSession
	if len(downloadedFiles) > 0 {
		var tmpFile = make([]RemoteSessionFile, 0)
		json.Unmarshal([]byte(downloadedFiles), &tmpFile)
		filePath = append(filePath, tmpFile...)
	}
	filePathByte, err := json.Marshal(filePath)
	if err != nil {
		return err
	}
	var param = map[string]interface{}{
		"downloaded_files": string(filePathByte),
	}
	return rsMysql.UpdateRemoteSessionById(ctx, RemoteSessionId, param)
}

// GetRemoteSessionCount 获取远程会话总数
func (re *RemoteSession) GetRemoteSessionCount(ctx context.Context, uid int, role int) int {
	var remoteSessionModel mysqls.RemoteSession
	return int(remoteSessionModel.GetRemoteSessionCount(ctx, uid, role))
}

// GetRemoteSessionById 根据id获取远程会话
func (re *RemoteSession) GetRemoteSessionById(ctx context.Context, id int) (mysqls.RemoteSession, error) {
	var rsMysql mysqls.RemoteSession
	return rsMysql.GetRemoteSessionById(ctx, id)
}

// InvokeFalconPostExploit 调用falcon后渗透功能
func (re *RemoteSession) InvokeFalconPostExploit(ctx context.Context, sessionId, remoteControl string, commands []string) (string, error) {
	controlType := "remote_control"
	if remoteControl == "远控" {
		controlType = "remote_control" // 默认值
	} else if remoteControl == "webshell" {
		controlType = "webshell" // 后渗透
	} else if remoteControl == "反弹shell" {
		controlType = "reverse_shell" // 后渗透
	}
	// 使用channel来等待回调结果
	resultChan := make(chan string, 1)

	handlePostExploitResult := func(ctx context.Context, result string) {
		var scannerLog invoke.ScannerLog
		err := json.Unmarshal([]byte(result), &scannerLog)
		if err != nil {
			fmt.Println("unmarshal error: ", err)
			resultChan <- ""
			return
		}
		fmt.Println("postResult: ", scannerLog.Content)
		resultChan <- scannerLog.Content
	}

	// 异步执行falcon post exploit
	go func() {
		invoke.InvokeFalconPostExploitWithSession(ctx, sessionId, controlType, commands, handlePostExploitResult)
	}()

	// 等待回调结果
	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
