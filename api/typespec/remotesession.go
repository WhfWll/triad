package typespec

// RemoteSessionListReq 会话列表请求结构体
type RemoteSessionListReq struct {
	TaskId   int    `json:"taskId" form:"taskId" binding:"required"`
	RiskType int    `json:"riskType" form:"riskType"`
	Search   string `json:"search" form:"search"`
	Page     int    `json:"page" form:"page" binding:"required"`
	Size     int    `json:"size" form:"size" binding:"required"`
	TargetID int    `json:"targetID" form:"targetID"`
}

// RemoteSessionListRes 会话列表返回结构体
type RemoteSessionListRes struct {
	Success bool                    `json:"success"`
	Count   int64                   `json:"count"`
	Result  []RemoteSessionListInfo `json:"result"`
}

// RemoteSessionListInfo 会话信息
type RemoteSessionListInfo struct {
	ID            int    `json:"id"`
	SessionNum    string `json:"sessionNum"`
	TargetUrl     string `json:"targetUrl"`
	Route         string `json:"route"`
	RemoteControl string `json:"remoteControl"`
	Status        string `json:"status"`
}

// RemoteSessionDelReq 会话删除
type RemoteSessionDelReq struct {
	IDs string `json:"ids" form:"ids" binding:"required"`
}

// RemoteSessionInfoReq 会话详情请求结构体
type RemoteSessionInfoReq struct {
	ID int `json:"id" form:"id" binding:"required"`
}

// RemoteSessionInfoRes 会话详情信息
type RemoteSessionInfoRes struct {
	Detail          []RemoteSessionInfo `json:"detail"`
	FileDirectory   FileDirectory       `json:"fileDirectory"`
	DownloadedFiles interface{}         `json:"downloadedFiles"`
}

// RemoteSessionDirReq 会话列出目录请求结构体
type RemoteSessionDirReq struct {
	ID int `json:"id" form:"id" binding:"required"`
}

// RemoteSessionDirRes 列出目录响应结构体
type RemoteSessionDirRes struct {
	Result []RemoteSessionFile `json:"result"`
}

// RemoteSessionFile 文件信息
type RemoteSessionFile struct {
	FileName   string `json:"fileName"`   //文件名
	FileType   string `json:"fileType"`   //文件类型
	FileSize   string `json:"fileSize"`   //文件大小
	FilePerm   string `json:"filePerm"`   //文件权限
	LastModify string `json:"lastModify"` //最后修改时间
}

// RemoteSessionInfo 详情信息
type RemoteSessionInfo struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
}

type ExceShellManyReq struct {
	RemoteSessionId int    `json:"remoteSessionId" form:"remoteSessionId" binding:"required"`
	CaptureType     int    `json:"captureType" form:"captureType"`
	CaptureInfoIds  string `json:"captureInfoIds" form:"captureInfoIds"`
	FileName        string `json:"fileName" form:"fileName"`
	FilePath        string `json:"filePath" form:"filePath"`
}
type ExceShellManyResp struct {
	Result string `json:"result"`
}
