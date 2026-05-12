package application

import (
	"context"
	"encoding/json"
	"errors"
	"gitlabee.4dogs.cn/common/file"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/enums"
)

// RemoteSessionApp 远程会话
type RemoteSessionApp struct {
}

// RemoteSessionList 会话列表
func (a *RemoteSessionApp) RemoteSessionList(ctx context.Context, req *typespec.RemoteSessionListReq, resp *typespec.RemoteSessionListRes) error {
	var remoteSessionSrc services.RemoteSession
	res, count, err := remoteSessionSrc.GetRemoteSessionList(ctx, req.TaskId, req.TargetID, req.Page, req.Size, req.RiskType, req.Search)
	if err != nil {
		return err
	}
	var vulEvidenceListInfo []typespec.RemoteSessionListInfo
	for _, v := range res {
		vulEvidenceListInfo = append(vulEvidenceListInfo, typespec.RemoteSessionListInfo{
			ID:            v.ID,
			SessionNum:    v.SessionID,
			TargetUrl:     v.TargetURL,
			RemoteControl: v.RemoteControl,
			Route:         v.Route,
			Status:        enums.SessionStatusMap[v.Status],
		})
	}
	*resp = typespec.RemoteSessionListRes{
		Count:  count,
		Result: vulEvidenceListInfo,
	}
	return nil
}

// DelRemoteSession 删除会话
func (a *RemoteSessionApp) DelRemoteSession(ctx context.Context, req *typespec.RemoteSessionDelReq) error {
	var remoteSessionSrc services.RemoteSession
	return remoteSessionSrc.DelGetRemoteSessions(ctx, req.IDs)
}

// RemoteSessionInfo 详情
func (a *RemoteSessionApp) RemoteSessionInfo(ctx context.Context, req *typespec.RemoteSessionInfoReq, resp *typespec.RemoteSessionInfoRes) error {
	var (
		remoteSessionSrc services.RemoteSession
		detail           map[string]interface{}
		downloadFiles    interface{}
	)
	// 0 获取漏洞详情
	res, err := remoteSessionSrc.GetRemoteSessionInfo(ctx, req.ID)
	if err != nil {
		return err
	}
	// 1 处理风险详情
	if err = json.Unmarshal([]byte(res.Detail), &detail); err != nil {
		return err
	}
	detail["alive"] = enums.SessionStatusMap[res.Status]
	// 2 下载文件列表
	if len(res.DownloadedFiles) != 0 {
		if err = json.Unmarshal([]byte(res.DownloadedFiles), &downloadFiles); err != nil {
			return err
		}
	}
	*resp = typespec.RemoteSessionInfoRes{
		Detail:          formatCommandExec(detail),
		FileDirectory:   typespec.FileDirectory{},
		DownloadedFiles: downloadFiles,
	}
	return nil
}

// formatCommandExec 初始化远程控制返回
func formatCommandExec(detail map[string]interface{}) (res []typespec.RemoteSessionInfo) {
	return []typespec.RemoteSessionInfo{
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
			Title: "远控类型",
			Value: detail["shellType"],
		}, {
			Title: "在线状态",
			Value: detail["alive"],
		},
	}
}

// FileDownload 文件下载
func (a *RemoteSessionApp) FileDownload(ctx context.Context, req *typespec.FileDownloadReq, resp *typespec.FileDownloadResp) (string, string, error) {
	var remoteSessionSrc services.RemoteSession
	fileList, _ := remoteSessionSrc.GetRemoteSessionDownLoadFileList(ctx, req.RemoteSessionID)
	for _, v := range fileList {
		if v["fileName"] == req.FileName {
			isExit := file.CheckFileExist(v["filePath"])
			if !isExit {
				return "", "", errors.New("找不到文件，请重新导出")
			} else {
				return req.FileName, v["filePath"], nil
			}
			break
		}
	}
	return "", "", errors.New("下载文件名有误")
}

// DelFile 文件下载删除
func (a *RemoteSessionApp) DelFile(ctx context.Context, req *typespec.DelFileReq, resp *typespec.FileDownloadResp) error {
	var (
		remoteSessionSrc services.RemoteSession
		delFileName      string
	)
	newDownloadFile := make([]map[string]string, 0)
	fileList, _ := remoteSessionSrc.GetRemoteSessionDownLoadFileList(ctx, req.RemoteSessionID)
	for _, v := range fileList {
		if v["fileName"] == req.FileName {
			isExit := file.CheckFileExist(v["filePath"])
			if isExit {
				err := file.RemoveFile(v["filePath"])
				if err == nil {
					delFileName = req.FileName
				}
			}
			break
		} else {
			newDownloadFile = append(newDownloadFile, v)
		}
	}
	if delFileName != "" {
		// 更新数据库
		newDownloadFileByte, _ := json.Marshal(newDownloadFile)
		remoteSessionSrc.UpdateRemoteSessionDownLoadFileList(ctx, req.RemoteSessionID, string(newDownloadFileByte))
	}
	return nil
}

// ExceShellMany 批量收集信息
func (a *RemoteSessionApp) ExceShellMany(ctx context.Context, req *typespec.ExceShellManyReq, resp *typespec.ExceShellManyResp) error {
	var (
		rsSrv       services.RemoteSession
		filePathMap = make([]services.RemoteSessionFile, 0)
	)
	//查询数据
	rsRes, err := rsSrv.GetRemoteSessionInfo(ctx, req.RemoteSessionId)
	if err != nil || rsRes.ID == 0 {
		return errors.New("找不到数据...")
	}
	if rsRes.Status != enums.SessionStatusSucc {
		return errors.New("远控目标不在线，无法请求...")
	}
	//处理信息收集
	if len(req.CaptureInfoIds) > 0 {
		cmdShellArray := rsSrv.GetShellCmdMany(req.CaptureType, req.CaptureInfoIds) //转换命令
		resp.Result = rsSrv.ExceShellMany(rsRes.SessionKey, cmdShellArray)          //执行shell
	}
	//处理文件收集
	if len(req.FileName) > 0 && req.CaptureType == enums.CaptureTypeWindows { //msf
		tmpFilePathMap, err := rsSrv.FindFileWindows(rsRes.SessionKey, req.FileName, req.FilePath) //查找文件
		if err != nil {
			return err
		}
		filePathMap, err = rsSrv.DownloadFileWindows(rsRes.SessionKey, tmpFilePathMap)
		if err != nil {
			return err
		}
	}
	//修改文件下载数据
	if len(filePathMap) > 0 {
		err = rsSrv.SaveDownloadFiles(ctx, rsRes.ID, rsRes.DownloadedFiles, filePathMap)
		if err != nil {
			return err
		}
	}
	return nil
}
