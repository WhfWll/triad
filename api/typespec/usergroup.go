package typespec

import "smart/tools/data"

type CommonPageReq struct {
	Page int `json:"page" form:"page" v:"required|integer|min:1#页数必填"`
	Size int `json:"size" form:"size" v:"required#数量必填"`
}

// UserGroupListReq 用户组 - 组列表
type UserGroupListReq struct {
	CommonPageReq
}
type UserGroupListRes struct {
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
	Total int                    `json:"total"`
	List  []UserGroupListItemRes `json:"list"`
}
type UserGroupListItemRes struct {
	Id         int             `json:"id"`
	Name       string          `json:"name"`
	Number     int             `json:"number"`
	Pid        int             `json:"pid"`
	PidStr     string          `json:"pid_str"`
	PidArr     [][]interface{} `json:"pid_arr"`
	CreateTime string          `json:"create_time"`
	RangeOpen  int             `json:"range_open"`
	Range      string          `json:"range"`
}

// UserGroupCreateReq 用户组 - 创建
type UserGroupCreateReq struct {
	Name      string `json:"name" form:"name" v:"required#组名称必填"`
	Pid       int    `json:"pid" form:"pid"`
	RangeOpen int    `json:"range_open" form:"range_open"`
	Range     string `json:"range" form:"range"`
}

// UserGroupCreateRes 用户组 - 创建
type UserGroupCreateRes struct {
}

// UserGroupSelectReq 用户组 - 级联选择器
type UserGroupSelectReq struct {
}

// UserGroupSelectRes 用户组 - 级联选择器
type UserGroupSelectRes struct {
	List []*data.TreeList `json:"list"`
}

// UserGroupUpdateReq 用户组 - 编辑
type UserGroupUpdateReq struct {
	Id        int    `json:"id" form:"id" v:"required#组ID缺失"`
	Name      string `json:"name" form:"name" v:"required#组名称必填"`
	Pid       int    `json:"pid" form:"pid"`
	RangeOpen int    `json:"range_open" form:"range_open"`
	Range     string `json:"range" form:"range"`
}

// UserGroupUpdateRes 用户组 - 编辑
type UserGroupUpdateRes struct {
}

// UserGroupUpdateStatusReq 用户组 - 修改状态
type UserGroupUpdateStatusReq struct {
	// 英文逗号分割
	Ids    string `json:"ids" form:"ids" v:"required#组ID缺失"`
	Status int    `json:"status" form:"status" v:"required|in:0,1#状态缺失|状态仅支持0或1"`
}

// UserGroupUpdateStatusRes 用户组 - 修改状态
type UserGroupUpdateStatusRes struct {
}

// UserGroupUserPreselectionReq 用户组 - 组内成员 - 预选列表
type UserGroupUserPreselectionReq struct {
	CommonPageReq
	Id      int    `json:"id" form:"id" v:"required#组ID缺失"`
	Keyword string `json:"keyword" form:"keyword"`
}

// UserGroupUserPreselectionRes 用户组 - 组内成员 - 预选列表
type UserGroupUserPreselectionRes struct {
	Page  int                                `json:"page"`
	Size  int                                `json:"size"`
	Total int                                `json:"total"`
	List  []UserGroupUserPreselectionItemRes `json:"list"`
}
type UserGroupUserPreselectionItemRes struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Role      int    `json:"role"`
	RoleStr   string `json:"role_str"`
	GroupList string `json:"group_list"`
	Status    int    `json:"status"`
	StatusStr string `json:"status_str"`
	Selected  bool   `json:"selected"`
}

// UserGroupUserAlreadyReq 用户组 - 组内成员 - 已选列表
type UserGroupUserAlreadyReq struct {
	Id int `json:"id"  form:"id" v:"required#组ID缺失"`
}

// UserGroupUserAlreadyRes 用户组 - 组内成员 - 已选列表
type UserGroupUserAlreadyRes struct {
	List []UserGroupUserAlreadyItemRes `json:"list"`
}
type UserGroupUserAlreadyItemRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// UserGroupUserRelationReq 用户组 - 组内成员 - 关联
type UserGroupUserRelationReq struct {
	GroupId int    `json:"group_id" form:"group_id" v:"required#组ID缺失"`
	UserIds string `json:"user_ids" form:"user_ids" v:"required#用户IDs缺失"`
}

// UserGroupUserRelationRes 用户组 - 组内成员 - 关联
type UserGroupUserRelationRes struct {
}
