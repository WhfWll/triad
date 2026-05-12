package typespec

// VulEvidenceListReq 漏洞取证列表请求结构体
type VulEvidenceListReq struct {
	TaskId   int    `json:"taskId" form:"taskId" binding:"required"`
	TargetId int    `json:"targetId" form:"targetId"`
	RiskType int    `json:"riskType" form:"riskType"`
	Search   string `json:"search" form:"search"`
	Page     int    `json:"page" form:"page" binding:"required"`
	Size     int    `json:"size" form:"size" binding:"required"`
	TargetID string `json:"targetID" form:"targetID"`
}

// VulEvidenceListRes 漏洞证据列表返回结构体
type VulEvidenceListRes struct {
	Success bool                  `json:"success"`
	Count   int64                 `json:"count"`
	Result  []VulEvidenceListInfo `json:"result"`
}

// VulEvidenceListInfo 漏洞证据列表信息
type VulEvidenceListInfo struct {
	ID            int    `json:"id"`
	VulName       string `json:"vulName"`
	RiskType      string `json:"riskType"`
	RiskTypeNum   int    `json:"riskTypeNum"`
	TargetUrl     string `json:"targetUrl"`
	CheckResultId int    `json:"checkResultID"`
	TargetId      string `json:"targetID"`
}

// VulEvidenceDelReq 漏洞取证删除
type VulEvidenceDelReq struct {
	IDs string `json:"ids" form:"ids" binding:"required"`
}

// VulEvidenceInfoReq 漏洞取证详情请求结构体
type VulEvidenceInfoReq struct {
	ID int `json:"id" form:"id" binding:"required"`
}

// ToCaptureInfoReq 发起抓取信息
type ToCaptureInfoReq struct {
	ID          int `json:"id" form:"id" binding:"required"`
	CaptureType int `json:"captureType" form:"captureType"`
}

// ToCaptureInfoRes 发起抓取动作
type ToCaptureInfoRes struct {
	Cmd string `json:"cmd"`
}

// VulEvidenceInfoRes 漏洞证据详情信息
type VulEvidenceInfoRes struct {
	RiskDetail      []VulEvidenceInfo `json:"riskDetail"`
	FileDirectory   FileDirectory     `json:"fileDirectory"`
	DownloadedFiles interface{}       `json:"downloadedFiles"`
}

// VulEvidenceInfo 漏洞证据详情信息
type VulEvidenceInfo struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
}

// RiskDetail 风险详情
type RiskDetail struct {
	VulName           string `json:"vulName"`
	VulRisk           string `json:"vulRisk"`
	VulTargetUrl      string `json:"vulTargetUrl"`
	RemoteControlType string `json:"remoteControlType"`
	Alive             string `json:"alive"`
	Shell             string `json:"shell"`
}

// FileDirectory 文件目录
type FileDirectory struct {
}

// DownloadedFiles 已下载文件
type DownloadedFiles struct {
	FileName string `json:"fileName"`
	FileSize string `json:"fileSize"`
}

// DelVulEvidenceDownloadedFile 删除漏洞证据下载文件
type DelVulEvidenceDownloadedFile struct {
	ID int `json:"id" form:"id"`
}

// DownLoadVulEvidenceDownloadedFile 下载漏洞证据下载文件
type DownLoadVulEvidenceDownloadedFile struct {
	ID int `json:"id" form:"id"`
}

// DownLoadVulEvidenceDownloadedFileRes 下载漏洞证据文件-返回
type DownLoadVulEvidenceDownloadedFileRes struct {
}

// EvidenceUseReq 证据利用请求结构体
type EvidenceUseReq struct {
	CheckResultId int    `form:"id" json:"id" v:"required#结果ID必填"`
	Cmd           string `form:"cmd" json:"cmd" v:"required#命令必填"`
}

// EvidenceUseRes 证据利用返回结构体
type EvidenceUseRes struct {
	Key    string `json:"key"`
	Result string `json:"result"`
}

type FileManagementReq struct {
	Id   int    `form:"id" json:"id" v:"required#结果ID必填"`
	Path string `form:"path" json:"path" v:"required#命令必填"`
}

type FileManagementResp struct {
	List interface{} `json:"list"`
}

type ShellFileDownloadReq struct {
	Id   int    `form:"id" json:"id" v:"required#结果ID必填"`
	Path string `form:"path" json:"path" v:"required#命令必填"`
}

type ShellFileDownloadResp struct {
	List interface{} `json:"list"`
}

// FileDownloadReq 文件下载
type FileDownloadReq struct {
	RemoteSessionID int    `form:"remoteSessionID" json:"remoteSessionID" binding:"required"`
	FileName        string `form:"fileName" json:"fileName" binding:"required"`
}

// DelFileReq 删除文件
type DelFileReq struct {
	RemoteSessionID int    `form:"remoteSessionID" json:"remoteSessionID" binding:"required"`
	FileName        string `form:"fileName" json:"fileName" binding:"required"`
}

// FileDownloadResp 下载文件返回
type FileDownloadResp struct {
}
