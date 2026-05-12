package gofishtypespec

type GroupDetailReq struct {
	ID int64 `json:"id" form:"id" common:"活动id" v:"required|min:1#活动id最小值为1"`
}

// Target 表示组中的目标用户
type Target struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Position  string `json:"position"`
}

// CreateGroupReq 创建用户组的请求
type CreateGroupReq struct {
	Name    string   `json:"name"`
	Targets []Target `json:"targets"`
}

type UpdateGroupReq struct {
	ID      int64    `json:"id" form:"id" common:"着陆页id" v:"required|min:1#着陆页id最小值为1"`
	Name    string   `json:"name"`
	Targets []Target `json:"targets"`
}
