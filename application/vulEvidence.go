package application

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"smart/api/typespec"
	"smart/services"
	"smart/services/reverse_shell"
	"smart/tools/enums"
	"smart/tools/file"
	"smart/tools/sshutils"
	"strconv"
	"strings"
	"time"
)

// VulEvidenceListApp 任务漏洞证据
type VulEvidenceListApp struct {
}

// VulEvidenceList 漏洞取证列表
func (a *VulEvidenceListApp) VulEvidenceList(ctx context.Context, req *typespec.VulEvidenceListReq, resp *typespec.VulEvidenceListRes) error {
	var taskVulEvicenceSrc services.TaskEvidence
	res, count, err := taskVulEvicenceSrc.GetTaskEvidenceList(ctx, req.TaskId, req.TargetId, req.Page, req.Size, req.RiskType, req.Search)
	if err != nil {
		return err
	}
	var vulEvidenceListInfo []typespec.VulEvidenceListInfo
	for _, v := range res {
		vulEvidenceListInfo = append(vulEvidenceListInfo, typespec.VulEvidenceListInfo{
			ID:            v.ID,
			VulName:       v.VulName,
			RiskType:      enums.TaskEvidenceMap[v.RiskType],
			RiskTypeNum:   0,
			TargetUrl:     v.TargetURL,
			CheckResultId: 0,
			TargetId:      "",
		})
	}
	*resp = typespec.VulEvidenceListRes{
		Count:  count,
		Result: vulEvidenceListInfo,
	}
	return nil
}

// DelVulEvidence 删除漏洞取证
func (a *VulEvidenceListApp) DelVulEvidence(ctx context.Context, req *typespec.VulEvidenceDelReq) error {
	var taskVulEvicenceSrc services.TaskEvidence
	return taskVulEvicenceSrc.DelTaskEvidences(ctx, req.IDs)
}

// VulEvidenceInfo 漏洞取证详情
func (a *VulEvidenceListApp) VulEvidenceInfo(ctx context.Context, req *typespec.VulEvidenceInfoReq, resp *typespec.VulEvidenceInfoRes) error {
	var (
		taskVulEvicenceSrc services.TaskEvidence
		riskDetail         []typespec.VulEvidenceInfo
		detail             map[string]interface{}
		downloadFiles      interface{}
	)
	// 0 获取漏洞详情
	res, err := taskVulEvicenceSrc.GetTaskEvidenceInfo(ctx, req.ID)
	if err != nil {
		return err
	}
	// 1 处理风险详情
	if err = json.Unmarshal([]byte(res.RiskDetail), &detail); err != nil {
		return err
	}
	switch res.RiskType {
	case enums.VulScriptEvidenceTypeWeakPass:
		riskDetail = formatWeakPass(detail)
	case enums.VulScriptEvidenceTypeInfoLeak:
		riskDetail = formatInfoLeak(ctx, detail)
	case enums.VulScriptEvidenceTypeFileLeak:
		riskDetail = formatFileLeak(ctx, detail)
	case enums.VulScriptEvidenceTypeData:
		riskDetail = formatDataBase(detail)
	}
	// 2 下载文件列表
	if len(res.DownloadedFiles) != 0 {
		if err = json.Unmarshal([]byte(res.DownloadedFiles), &downloadFiles); err != nil {
			return err
		}
	}
	*resp = typespec.VulEvidenceInfoRes{
		RiskDetail:      riskDetail,
		FileDirectory:   typespec.FileDirectory{},
		DownloadedFiles: downloadFiles,
	}
	return nil
}

// RiskTypeInfoEnum 风险类型枚举
func (a *VulEvidenceListApp) RiskTypeInfoEnum(ctx context.Context, resp *map[int]string) error {
	*resp = enums.TaskEvidenceMap
	return nil
}

// CaptureInfoEnum 抓取信息枚举
func (a *VulEvidenceListApp) CaptureInfoEnum(ctx context.Context, resp *map[int]string) error {
	*resp = enums.TaskEvidenceCaptureInfoMap
	return nil
}

// ToCaptureInfo 抓取信息
func (a *VulEvidenceListApp) ToCaptureInfo(ctx context.Context, req *typespec.ToCaptureInfoReq, resp *typespec.ToCaptureInfoRes) error {
	if req.CaptureType == 1 {
		resp.Cmd = enums.TaskEvidenceCaptureInfoWindowsExecMap[req.ID]
		return nil
	}
	resp.Cmd = enums.TaskEvidenceCaptureInfoCentosExecMap[req.ID]
	return nil
}

// VulEvidenceUse 利用
func (a *VulEvidenceListApp) VulEvidenceUse(ctx context.Context, req *typespec.EvidenceUseReq, resp *typespec.EvidenceUseRes) error {
	var rsSrv services.RemoteSession
	remoteSession, err := rsSrv.GetRemoteSessionInfo(ctx, req.CheckResultId)
	if err != nil {
		return err
	}
	if remoteSession.RemoteControl == "webshell" {
		wsm := services.Wsm{}
		var tempMap map[string]string
		err = json.Unmarshal([]byte(remoteSession.Detail), &tempMap)
		if err != nil {
			return err
		}
		if tempMap["webshell"] == "" {
			return errors.New("webshell is empty")
		}
		if tempMap["webshellType"] == enums.WebShellTypeGodzilla {
			err = wsm.InitGodzilla(tempMap["webshell"])
			if err != nil {
				return err
			}
			status, err := wsm.PingGodzilla(ctx)
			if err != nil {
				return err
			}
			fmt.Println("godzilla webshell connect status: ", status)
			result, err := wsm.CommandExecGodzilla(ctx, req.Cmd)
			if err != nil {
				return err
			}
			resp.Key = req.Cmd
			resp.Result = result
		} else {
			err = wsm.Init(tempMap["webshell"])
			if err != nil {
				return err
			}
			status, err := wsm.Ping(ctx)
			if err != nil {
				return err
			}
			fmt.Println("behinder webshell connect status: ", status)
			result, err := wsm.CommandExec(ctx, req.Cmd)
			if err != nil {
				return err
			}
			resp.Key = req.Cmd
			resp.Result = result
		}
		return nil
	} else if remoteSession.RemoteControl == "ssh" {
		var tempMap map[string]string
		err = json.Unmarshal([]byte(remoteSession.Detail), &tempMap)
		if err != nil {
			return err
		}
		host := remoteSession.TargetURL
		if strings.Contains(host, "://") {
			u, err := url.Parse(host)
			if err == nil {
				host = u.Hostname()
			}
		} else {
			if h, _, err := net.SplitHostPort(host); err == nil {
				host = h
			}
		}
		port := tempMap["port"]
		if port == "" {
			port = "22"
		}
		username := tempMap["username"]
		password := tempMap["password"]

		result, err := sshutils.ExecCommand(host, port, username, password, req.Cmd)
		if err != nil {
			return err
		}
		resp.Key = req.Cmd
		resp.Result = result
		return nil
	} else if remoteSession.RemoteControl == "反弹shell" || remoteSession.Status == enums.VulScriptEvidenceTypeReverseShell {
		// 直接调用内部反弹shell服务执行命令
		result, err := reverse_shell.GetService().ExecuteCommand(ctx, remoteSession.SessionID, req.Cmd)
		if err != nil {
			fmt.Println("ExecuteCommand error: ", err)
			return err
		}
		resp.Key = req.Cmd
		resp.Result = result
		return nil
	} else if remoteSession.RemoteControl == "远控" {
		var command string
		if req.Cmd == "screenshot" {
			command = rewriteScreenshotCommand()
		} else {
			command = req.Cmd
		}
		commands := []string{command}
		execResult, err := rsSrv.InvokeFalconPostExploit(ctx, remoteSession.SessionID, remoteSession.RemoteControl, commands)
		if err != nil {
			return err
		}
		if strings.Contains(execResult, "Screenshot saved to") {
			var screenshotInfo map[string]string
			screenshotInfo, err = getScreenshotInfo(execResult)
			if err != nil {
				return err
			}
			// save screenshot info into db
			downloadedFilesStr := remoteSession.DownloadedFiles
			var downloadedFiles = make([]map[string]string, 0)
			if downloadedFilesStr != "" {
				err = json.Unmarshal([]byte(downloadedFilesStr), &downloadedFiles)
				if err != nil {
					return err
				}
			}
			fmt.Println("downloadedFiles: ", downloadedFiles)
			downloadedFiles = append(downloadedFiles, screenshotInfo)
			var newDownloadedInfo []byte
			newDownloadedInfo, err = json.Marshal(downloadedFiles)
			if err != nil {
				return err
			}
			fmt.Println("newDownloadedInfo: ", string(newDownloadedInfo))
			remoteSession.DownloadedFiles = string(newDownloadedInfo)
			err = remoteSession.UpdateRemoteSession(ctx)
			if err != nil {
				return err
			}
		}
		resp.Key = req.Cmd
		resp.Result = execResult
	}
	return nil
}

func (a *VulEvidenceListApp) FileManagement(ctx context.Context, req *typespec.FileManagementReq, resp *typespec.FileManagementResp) error {
	var rsSrv services.RemoteSession
	remoteSession, err := rsSrv.GetRemoteSessionInfo(ctx, req.Id)
	if err != nil {
		return err
	}
	if remoteSession.RemoteControl == "webshell" {
		wsm := services.Wsm{}
		var tempMap map[string]string
		err = json.Unmarshal([]byte(remoteSession.Detail), &tempMap)
		if err != nil {
			return err
		}
		if tempMap["webshell"] == "" {
			return errors.New("webshell is empty")
		}

		var (
			filePathMap = []map[string]string{}
			fileString  string
		)
		if tempMap["webshellType"] == enums.WebShellTypeGodzilla {
			err = wsm.InitGodzilla(tempMap["webshell"])
			if err != nil {
				return err
			}
			status, err := wsm.PingGodzilla(ctx)
			if err != nil {
				return err
			}
			fmt.Println("godzilla webshell connect status: ", status)
			fileString, err = wsm.GodzillaFileList(req.Path)
			if err != nil {
				return err
			}
		} else {
			err = wsm.Init(tempMap["webshell"])
			if err != nil {
				return err
			}
			status, err := wsm.Ping(ctx)
			if err != nil {
				return err
			}
			fmt.Println("behinder webshell connect status: ", status)
			fileString, err = wsm.FileList(req.Path)
			if err != nil {
				return err
			}
		}
		err = json.Unmarshal([]byte(fileString), &filePathMap)
		if err != nil {
			return err
		}
		for i := 0; i < len(filePathMap); i++ {
			for k, v := range filePathMap[i] {
				decodedBytes, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					continue
				}
				filePathMap[i][k] = string(decodedBytes)
			}
		}
		resp.List = filePathMap
		return nil
	} else if remoteSession.RemoteControl == "反弹shell" {
		// Use internal reverse shell service
		path := req.Path
		if path == "" {
			path = "."
		}

		files, err := reverse_shell.GetService().ListFiles(ctx, remoteSession.SessionID, path)
		if err != nil {
			fmt.Println("ListFiles error: ", err)
			return err
		}

		// Build file path list compatible with frontend
		// Format: [{"lastModified":"...","name":"...","perm":"...","size":"...","type":"..."}]
		var fileList []map[string]string
		for _, file := range files {
			fileMap := make(map[string]string)
			fileMap["name"] = file.Name
			if file.IsDir {
				fileMap["type"] = "directory"
			} else {
				fileMap["type"] = "file"
			}
			fileMap["size"] = fmt.Sprintf("%d", file.Size)
			fileMap["lastModified"] = file.ModTime
			fileMap["perm"] = file.Perm

			fileList = append(fileList, fileMap)
		}
		resp.List = fileList
		return nil
	}
	return nil
}

// ShellFileDownload 通过shell下载文件
func (a *VulEvidenceListApp) ShellFileDownload(ctx context.Context, req *typespec.ShellFileDownloadReq, resp *typespec.ShellFileDownloadResp) error {
	var rsSrv services.RemoteSession
	remoteSession, err := rsSrv.GetRemoteSessionInfo(ctx, req.Id)
	if err != nil {
		return err
	}
	if remoteSession.RemoteControl == "webshell" {
		wsm := services.Wsm{}
		var tempMap map[string]string
		err = json.Unmarshal([]byte(remoteSession.Detail), &tempMap)
		if err != nil {
			return err
		}
		if tempMap["webshell"] == "" {
			return errors.New("webshell is empty")
		}

		var (
			filePathMap = []map[string]string{}
			fileString  string
		)
		if tempMap["webshellType"] == enums.WebShellTypeGodzilla {
			err = wsm.InitGodzilla(tempMap["webshell"])
			if err != nil {
				return err
			}
			status, err := wsm.PingGodzilla(ctx)
			if err != nil {
				return err
			}
			fmt.Println("godzilla webshell connect status: ", status)
			fileString, err = wsm.GodzillaFileList(req.Path)
			if err != nil {
				return err
			}
		} else {
			err = wsm.Init(tempMap["webshell"])
			if err != nil {
				return err
			}
			status, err := wsm.Ping(ctx)
			if err != nil {
				return err
			}
			fmt.Println("behinder webshell connect status: ", status)
			fileString, err = wsm.FileDownload(req.Path)
			if err != nil {
				return err
			}
		}
		err = json.Unmarshal([]byte(fileString), &filePathMap)
		if err != nil {
			return err
		}
		for i := 0; i < len(filePathMap); i++ {
			for k, v := range filePathMap[i] {
				decodedBytes, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					continue
				}
				filePathMap[i][k] = string(decodedBytes)
			}
		}
		resp.List = filePathMap
		return nil
	}
	return nil
}

// formatWeakPass 初始化弱口令返回
func formatWeakPass(detail map[string]interface{}) []typespec.VulEvidenceInfo {
	return []typespec.VulEvidenceInfo{
		{
			Title: "漏洞名称",
			Value: detail["vulName"],
		}, {
			Title: "漏洞风险",
			Value: detail["vulRisk"],
		}, {
			Title: "漏洞地址",
			Value: detail["vulTargetUrl"],
		}, {
			Title: "服务",
			Value: detail["service"],
		}, {
			Title: "账号",
			Value: detail["user"],
		}, {
			Title: "密码",
			Value: detail["password"],
		},
	}
}

// formatInfoLeak 初始化信息泄漏返回
func formatInfoLeak(ctx context.Context, detail map[string]interface{}) []typespec.VulEvidenceInfo {
	var (
		mapSetService services.MapSet
		infoType      string
		infoValue     string
	)
	mapSetValue, err := mapSetService.GetMapValue(ctx, enums.VulEvidenceInfoLeakMapSetObjKey)
	if err == nil && len(mapSetValue) != 0 {
		if v, ok := detail["vulName"]; ok {
			vulName := v.(string)
			var vulEvidenceMapSet []VulEvidenceMapSet
			if err = json.Unmarshal([]byte(mapSetValue), &vulEvidenceMapSet); err == nil {
				infoType, infoValue = handleMapSet(vulName, vulEvidenceMapSet)
			}
		}
	}
	return []typespec.VulEvidenceInfo{
		{
			Title: "漏洞名称",
			Value: detail["vulName"],
		}, {
			Title: "漏洞风险",
			Value: detail["vulRisk"],
		}, {
			Title: "漏洞地址",
			Value: detail["vulUrl"],
		}, {
			Title: "信息类型",
			Value: infoType,
		}, {
			Title: "数据价值",
			Value: infoValue,
		}, {
			Title: "泄漏信息",
			Value: detail["content"],
		},
	}
}

// formatFileLeak 初始化文件泄漏返回
func formatFileLeak(ctx context.Context, detail map[string]interface{}) []typespec.VulEvidenceInfo {
	var (
		mapSetService services.MapSet
		infoType      string
		infoValue     string
	)
	mapSetValue, err := mapSetService.GetMapValue(ctx, enums.VulEvidenceFileLeakMapSetObjKey)
	if err == nil && len(mapSetValue) != 0 {
		if v, ok := detail["vulName"]; ok {
			vulName := v.(string)
			var vulEvidenceMapSet []VulEvidenceMapSet
			if err = json.Unmarshal([]byte(mapSetValue), &vulEvidenceMapSet); err == nil {
				infoType, infoValue = handleMapSet(vulName, vulEvidenceMapSet)
			}
		}
	}
	return []typespec.VulEvidenceInfo{
		{
			Title: "漏洞名称",
			Value: detail["vulName"],
		}, {
			Title: "漏洞风险",
			Value: detail["vulRisk"],
		}, {
			Title: "漏洞地址",
			Value: detail["vulTargetUrl"],
		}, {
			Title: "文件类型",
			Value: infoType,
		}, {
			Title: "数据价值",
			Value: infoValue,
		}, {
			Title: "泄漏文件",
			Value: detail["content"],
		},
	}
}

// formatDataBase 初始化数据
func formatDataBase(detail map[string]interface{}) []typespec.VulEvidenceInfo {
	return []typespec.VulEvidenceInfo{
		{
			Title: "漏洞名称",
			Value: detail["vulName"],
		}, {
			Title: "漏洞风险",
			Value: detail["vulRisk"],
		}, {
			Title: "漏洞地址",
			Value: detail["vulUrl"],
		}, {
			Title: "数据库类型",
			Value: detail["dbType"],
		}, {
			Title: "数据库信息",
			Value: detail["content"],
		},
	}
}

func handleMapSet(name string, mapSetData []VulEvidenceMapSet) (string, string) {
	for i := 0; i < len(mapSetData); i++ {
		if strings.Contains(name, mapSetData[i].Match) {
			return mapSetData[i].Type, mapSetData[i].Value
		}
	}
	return "", ""
}

func rewriteScreenshotCommand() string {
	return fmt.Sprintf("screenshot -p /home/msf/tmp/screenshot_%v.png", time.Now().Unix())
}

func getScreenshotInfo(screenshotResult string) (map[string]string, error) {
	var err error
	var result = make(map[string]string)
	resultsSlices := strings.Split(screenshotResult, "Screenshot saved to: ")
	fmt.Println("resultsSlices: ", resultsSlices)
	if len(resultsSlices) < 2 {
		return result, fmt.Errorf("cannot get screenshot path from result: %s", screenshotResult)
	}
	// get item info
	filePath := strings.TrimSpace(resultsSlices[1])
	paths := strings.Split(filePath, "/")
	var fileName string
	if len(paths) > 1 {
		fileName = paths[len(paths)-1]
	}
	realFilePath := fmt.Sprintf("/opt/laozhi/msf/tmp/%s", fileName)
	fmt.Println("realFilePath: ", realFilePath)
	var fileSize int64
	fileSize, err = file.GetFileSize(realFilePath)
	if err != nil {
		return result, err
	}
	result["fileName"] = fileName
	result["filePath"] = realFilePath
	result["fileSize"] = strconv.FormatInt(fileSize, 10)
	fmt.Println("result: ", result)
	return result, nil
}
