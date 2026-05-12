package typespec

type TaskGroupCreateReq struct {
	Name     string `form:"name" json:"name" binding:"required"`         //任务组名称
	Describe string `form:"describe" json:"describe" binding:"required"` //任务组描述
}

type TaskGroupCreateResp struct {
}

// TaskGroupListReq  任务组列表请求数据
type TaskGroupListReq struct {
	Search string `form:"search" json:"search"`                //任务组过滤条件
	Page   int    `form:"page" json:"page" binding:"required"` //任务组页数
	Size   int    `form:"size" json:"size" binding:"required"` //任务组每页大小
}

// TaskGroupListResp 任务组列表响应数据
type TaskGroupListResp struct {
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
	Total int64               `json:"total"`
	List  []TaskGroupListInfo `json:"list"`
}

// TaskGroupListInfo 任务组信息
type TaskGroupListInfo struct {
	ID         int    `json:"id"`          //id
	Name       string `json:"name"`        //名称
	HighNum    int    `json:"high_num"`    //高危任务
	MiddleNum  int    `json:"middle_num"`  //中危任务
	LowNum     int    `json:"low_num"`     //低危任务
	SafeNum    int    `json:"safe_num"`    //安全任务
	UpdateTime string `json:"update_time"` //更新时间
	CreateTime string `json:"create_time"` //创建时间
	StatusNum  int    `json:"status_num"`  //状态码
	Status     string `json:"Status"`      //状态中文
	Describe   string `json:"describe"`    //描述信息
}

// TaskGroupDeleteReq 任务组删除请求
type TaskGroupDeleteReq struct {
	Id string `form:"id" json:"id" binding:"required"` //任务组名称
}

// TaskGroupDeleteResp 任务组删除响应
type TaskGroupDeleteResp struct {
}

// TaskGroupGroupBindReq 任务组 将任务与组绑定 请求
type TaskGroupGroupBindReq struct {
	TaskId  int `form:"task_id" json:"task_id"`   //任务id
	GroupId int `form:"group_id" json:"group_id"` //任务组id
}

// TaskGroupGroupBindResp 任务组 将任务与组绑定 响应
type TaskGroupGroupBindResp struct {
}

// TaskGroupTaskListReq 任务组 将任务与组绑定 请求
type TaskGroupTaskListReq struct {
	GroupId int `form:"group_id" json:"group_id"` //任务组id
	Page    int `form:"page" json:"page"`         //任务组页数
	Size    int `form:"size" json:"size"`         //任务组大小
}

// TaskGroupTaskListResp 任务组 任务列表 响应
type TaskGroupTaskListResp struct {
	Total int                        `json:"total"`
	List  []TaskGroupTaskListItemRes `json:"list"`
}
type TaskGroupTaskListItemRes struct {
	Id              int    `json:"id"`
	TaskName        string `json:"taskName"`
	ExecuteType     int    `json:"executeType"`
	ExecuteTypeName string `json:"executeTypeName"`
	RiskLevel       int    `json:"riskLevel"`
	RiskLevelName   string `json:"riskLevelName"`
	Status          int    `json:"status"`
	StatusName      string `json:"statusName"`
	TargetRisk      []int  `json:"targetRisk"` // 任务风险等级 下标区分 0高危 1中危 2低危 3安全
	CreateTime      string `json:"createTime"`
	UpdateTime      string `json:"updateTime"`
	Progress        int    `json:"progress"` // A6项目需要进度条，虽然并不准确
}

// TaskGroupOverViewReq 任务组 获取统计信息 请求
type TaskGroupOverViewReq struct {
	GroupId int `form:"group_id" json:"group_id"` //任务组id
}

// TaskGroupOverViewResp 任务组 获取统计信息 响应
type TaskGroupOverViewResp struct {
	Overview interface{} `json:"overview"`
}

// TaskGroupStatusReq 任务组 任务状态 请求
type TaskGroupStatusReq struct {
	GroupId int `form:"group_id" json:"group_id"` //任务组id
}

// TaskGroupStatusResp 任务组 任务状态 响应
type TaskGroupStatusResp struct {
	Status       string `json:"status"`
	StatusNumber int    `json:"statusNumber"`
}

type TaskGroupEditReq struct {
	Id       int    `form:"id" json:"id" binding:"required"`             //任务组id
	Name     string `form:"name" json:"name" binding:"required"`         //任务组名称
	Describe string `form:"describe" json:"describe" binding:"required"` //任务组描述
}

type TaskGroupEditResp struct {
}
