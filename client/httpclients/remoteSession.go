package httpclients

import (
	"encoding/json"
	"gitlabee.4dogs.cn/common/httpclient"
)

type GetRecentShellReq struct {
	Seconds      int `form:"seconds" json:"seconds" binding:"required"`
	EvidenceType int `form:"evidenceType" json:"evidenceType" binding:"required"`
}

type GetRecentShellResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SessionKey    string `json:"sessionKey"`
		SessionID     string `json:"sessionId"`
		Route         string `json:"route"`
		RemoteControl string `json:"remoteControl"`
	}
}

type ExecShellCmdReq struct {
	CmdKey string `form:"cmdKey" json:"cmdKey" binding:"required"`
	Cmd    string `form:"cmd" json:"cmd" binding:"required"`
}

type ExecShellCmdResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Result string `json:"result"`
	}
}

func GetRecentShell(req GetRecentShellReq) (GetRecentShellResp, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/shell/openrecentshell")
	var res GetRecentShellResp
	if err != nil {
		return res, err
	}
	h.GetUrlBuild(map[string]interface{}{
		"seconds":      req.Seconds,
		"evidenceType": req.EvidenceType,
	})
	response, err := h.Get() //发送get请求
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(response, &res) // 解析响应
	if err != nil {
		return res, err
	}
	return res, nil
}

func ExecShellCmd(req ExecShellCmdReq) (ExecShellCmdResp, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/shell/openshellcmd")
	var res ExecShellCmdResp
	if err != nil {
		return res, err
	}
	param := map[string]interface{}{
		"cmdKey": req.CmdKey,
		"cmd":    req.Cmd,
	}
	h.SetBody(param)
	response, err := h.Post() //发送get请求
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(response, &res) // 解析响应
	if err != nil {
		return res, err
	}
	return res, nil
}

func ExecRemoteCmd(req ExecShellCmdReq) (ExecShellCmdResp, error) {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/shell/openremotecmd")
	var res ExecShellCmdResp
	if err != nil {
		return res, err
	}
	param := map[string]interface{}{
		"cmdKey": req.CmdKey,
		"cmd":    req.Cmd,
	}
	h.SetBody(param)
	response, err := h.Post() //发送get请求
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(response, &res) // 解析响应
	if err != nil {
		return res, err
	}
	return res, nil
}
