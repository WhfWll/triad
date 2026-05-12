package typespec

import "time"

type LogsEnumRes struct {
	Type interface{} `json:"type"`
}

type LogAuditListReq struct {
	Page      int    `json:"page" form:"page" binding:"required"` //页码
	Size      int    `json:"size" form:"size" binding:"required"` //每页的数量
	Search    string `json:"search" form:"search"`                //搜索日志内容
	LogType   int    `json:"logType" form:"logType"`              //日志类型
	StartTime string `json:"startTime" form:"startTime"`          //开始时间
	EndTime   string `json:"endTime" form:"endTime"`              //结束时间
}

type LogAuditListRes struct {
	List  []LogAuditListItemRes `json:"list"`
	Total int64                 `json:"total"`
}

type LogAuditListItemRes struct {
	Id          int    `json:"id"`
	LogType     int    `json:"logType"`
	LogTypeName string `json:"logTypeName"`
	Content     string `json:"content"`
	Username    string `json:"username"`
	Ip          string `json:"ip"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
}

type LogAuditEmptyRes struct {
}

type LogBackupListReq struct {
	Page int `json:"page" form:"page" binding:"required"` //页码
	Size int `json:"size" form:"size" binding:"required"` //每页的数量
}

type LogBackupListRes struct {
	List  []LogBackupListItemRes `json:"list"`
	Total int64                  `json:"total"`
}

type LogBackupListItemRes struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
	Path       string `json:"path"`
}

type LogBackupDeleteReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type LogBackupDeleteRes struct {
}

type LogBackupDownloadReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type LogBackupDownloadRes struct {
	Path string `json:"path"`
}

type LogBackupConfigReq struct {
	IsOpen   int       `json:"isOpen" form:"isOpen" binding:"required"` //是否开启周期备份，1-开启，2-关闭
	Cycle    int       `json:"cycle" form:"cycle" binding:"required"`   //备份周期，按月计算
	SaveTime time.Time `json:"saveTime" form:"saveTime"`                //创建时间
	RunTime  time.Time `json:"runTime" form:"runTime"`                  //执行时间
}

type LogBackupConfigRes struct {
}

type LogBackupConfigInfoRes struct {
	IsOpen int `json:"isOpen"`
	Cycle  int `json:"cycle"`
}
type SetLogExpirationTimeReq struct {
	ExpirationTime int `json:"expirationTime" form:"expirationTime"`
}

type SetLogExpirationTimeResp struct {
}

type GetLogExpirationTimeResp struct {
	ExpirationTime int `json:"expirationTime"`
}

type LogBackupNowRes struct {
}
