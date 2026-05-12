package typespec

// 节点管理 - 设置是否启用分布式
type NodeSetDistributeReq struct {
	Status int `json:"status" form:"status"`
}
type NodeSetDistributeRes struct {
}

// 节点管理 - 获取是否启用分布式
type NodeIsDistributeRes struct {
	Status int `json:"status"`
}

// 节点管理 - 新增节点
type NodeAddReq struct {
	Name string `form:"name" json:"name" binding:"required"`
	Ip   string `form:"ip" json:"ip" binding:"required"`
	Port string `form:"port" json:"port" binding:"required"`
}
type NodeAddRes struct {
}

// 节点管理 - 编辑节点
type NodeEditReq struct {
	Id   int    `form:"id" json:"id" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
	Ip   string `form:"ip" json:"ip" binding:"required"`
	Port string `form:"port" json:"port" binding:"required"`
}
type NodeEditRes struct {
}

// 节点管理 - 节点详情
type NodeInfoReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type NodeInfoRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

// 节点列表
type NodeListReq struct {
	Page   int    `json:"page" form:"page" binding:"required"` //页码
	Size   int    `json:"size" form:"size" binding:"required"` //每页的数量
	Search string `json:"search" form:"search"`                // 搜索关键词
}
type NodeListRes struct {
	List  []NodeListItem `json:"list"`
	Total int64          `json:"total"`
}
type NodeListItem struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	RunningNum    int    `json:"runningNum"` // 运行任务数
	Status        int    `json:"status"`
	StatusEnum    string `json:"statusEnum"`
	IsDisable     int    `json:"IsDisable"` // 禁用状态 0启用 1禁用
	IsDisableEnum string `json:"isDisableEnum"`
	CreateTime    string `json:"createTime"`
	UpdateTime    string `json:"updateTime"`
}

// 节点删除
type NodeDelReq struct {
	Id string `form:"id" json:"id" binding:"required"`
}
type NodeDelRes struct {
}

// 节点禁用状态设置
type NodeDisOrEnableReq struct {
	Id        int `form:"id" json:"id" binding:"required"`
	IsDisable int `form:"isDisable" json:"isDisable"`
}
type NodeDisOrEnableRes struct {
}

// 获取所有可用节点 - 用于渗透任务时定向渗透使用
type NodeAllEnableReq struct {
}
type NodeAllEnableRes struct {
	List []NodeAllEnableItem `json:"list"`
}
type NodeAllEnableItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
