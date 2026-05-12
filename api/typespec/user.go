package typespec

// UserDetailReq 用户详情请求结构体
type UserDetailReq struct {
	UserID int `form:"userId" json:"userId"`
}

// UserListReq 用户列表请求结构体
type UserListReq struct {
	Page   int    `form:"page" json:"page" v:"required|integer|min:1#页数必填"`
	Size   int    `form:"size" json:"size" v:"required#数量必填"`
	Search string `form:"search" json:"search"`
}

// UserListRes 用户列表返回
type UserListRes struct {
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Total int64          `json:"total"`
	List  []UserListInfo `json:"list"`
}

// UserListInfo 用户列表信息
type UserListInfo struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	RoleStr       string `json:"roleStr"`
	Role          int    `json:"role"`
	AccountExpire string `json:"accountExpire"`
	LastTime      string `json:"lastTime"`
	Status        int    `json:"status"`
	StatusStr     string `json:"statusStr"`
	Email         string `json:"email"`
	Department    string `json:"department"`
	Remark        string `json:"remark"`
	GroupStr      string `json:"groupStr"`
	GroupIds      []int  `json:"groupIds"`
	IsAlive       int    `json:"isAlive"`
}

// UserEnumListReq 用户枚举列表请求结构体
type UserEnumListReq struct {
	Type int `form:"type"` // Type类型
}

// UserEnumListRes 用户枚举列表返回
type UserEnumListRes struct {
	Data interface{} `json:"data"`
}

// 登录验证吗
type UserLoginCaptchaReq struct {
}

type UserLoginCaptchaRes struct {
	CaptchaId string `json:"captchaId"`
	Data      string `json:"data"`
}

// UserLoginReq 用户登录
type UserLoginReq struct {
	CaptchaId string `form:"captchaId" json:"captchaId" binding:"required"`
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	CheckCode string `form:"checkCode" json:"checkCode" binding:"required"`
	IP        string `json:"ip"`
}
type UserLoginRes struct {
	Token      string `json:"token"`
	Uid        int    `json:"uid"`
	Role       int    `json:"role"`
	FirstLogin bool   `json:"firstLogin"`
}

// 用户管理 - 列表
type UserManageListReq struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
}
type UserManageListRes struct {
	Page  int               `json:"page"`
	Size  int               `json:"size"`
	Total int               `json:"total"`
	List  []UserListItemRes `json:"list"`
}
type UserListItemRes struct {
	Id            int        `json:"id"`
	Username      string     `json:"username"`
	RoleStr       string     `json:"roleStr"`
	Role          int        `json:"role"`
	Email         string     `json:"email"`
	AccountExpire string     `json:"accountExpire"`
	LastTime      string     `json:"lastTime"`
	GroupStr      string     `json:"groupStr"`
	GroupArr      [][]string `json:"groupArr"`
	Status        int        `json:"status"`
	StatusStr     string     `json:"statusStr"`
	Department    string     `json:"department"`
	Remark        string     `json:"remark"`
}

// UserManageOpReq 用户管理 - 创建
type UserManageOpReq struct {
	Username          string `form:"username" json:"username" binding:"required"`
	Role              int    `form:"role" json:"role" binding:"required"`
	Password          string `form:"password" json:"password"`
	Repassword        string `form:"repassword" json:"repassword"`
	Email             string `form:"email" json:"email"`
	AccountExpireTime string `form:"accountExpireTime" json:"accountExpireTime"`
	GroupIds          string `form:"groupIds" json:"groupIds"`
	Department        string `form:"department" json:"department"`
	Remark            string `form:"remark" json:"remark"`
	ID                int    `form:"id" json:"id"` // 修改使用
}

type UserManageCreateRes struct {
}

// 用户管理 - 更新
type UserManageUpdateReq struct {
	Id         int    `form:"id" json:"id" v:"required#用户ID缺失"`
	Username   string `form:"username" json:"username" binding:"required"`
	Role       int    `form:"role" json:"role" binding:"required"`
	Email      string `form:"email" json:"email" binding:"required"`
	GroupIds   string `form:"groupIds" json:"groupIds" binding:"required"`
	Department string `form:"department" json:"department"`
	Remark     string `form:"remark" json:"remark"`
}
type UserManageUpdateRes struct {
}

// UserDelReq 用户删除
type UserDelReq struct {
	UserIds string `form:"userIds" json:"userIds"`
}

// UserPassPWReq 修改密码
type UserPassPWReq struct {
	UserID      int    `form:"userId" json:"userId"`
	OldPassword string `form:"oldpassword" json:"oldpassword"`
	Password    string `form:"password" json:"password"`
	Repassword  string `form:"repassword" json:"repassword"`
	OperatorID  int    `form:"operatorId" json:"operatorId"`
}

// UserResetPassPWReq 重置密码
type UserResetPassPWReq struct {
	UserID     int    `form:"userId" json:"userId"`
	OperatorID int    `form:"operatorId" json:"operatorId"`
	Password   string `form:"password" json:"password"`
}

// UpdateUserExpReq 更新过期时间
type UpdateUserExpReq struct {
	UserID            int    `form:"userId" json:"userId"`
	AccountExpireTime string `form:"accountExpireTime" json:"accountExpireTime"`
	OperatorID        int    `form:"operatorId" json:"operatorId"`
}

// ChangeUserStatusReq 修改用户状态
type ChangeUserStatusReq struct {
	UserID     int `form:"userId" json:"userId"`
	Status     int `form:"status" json:"status"`
	OperatorID int `form:"operatorId" json:"operatorId"`
}

// LogOutReq 退出登录
type LogOutReq struct {
}

// UserLoginLianTongReq 用户登录
type UserLoginLianTongReq struct {
	SsoToken string `form:"sso_token" json:"sso_token" binding:"required"`
	UserName string `form:"username" json:"username" binding:"required"`
}
type UserLoginLianTongResp struct {
	Token    string `json:"token"`
	Uid      int    `json:"uid"`
	Role     int    `json:"role"`
	Username string `json:"username"`
}

// UserLoginSiYuanReq 用户登录
type UserLoginSiYuanReq struct {
	Token string `form:"token" json:"token" binding:"required"`
}
type UserLoginSiYuanResp struct {
	Token    string `json:"token"`
	Uid      int    `json:"uid"`
	Role     int    `json:"role"`
	Username string `json:"username"`
}

// 中测yakit登录
type LoginByApiTokenReq struct {
	Token string `form:"token" json:"token" binding:"required"`
}
type LoginByApiTokenResp struct {
	Token    string `json:"token"`
	Uid      int    `json:"uid"`
	Role     int    `json:"role"`
	Username string `json:"username"`
}

// 航天运载免密登录
type LoginByHTYZReq struct {
	Token string `form:"token" json:"token" binding:"required"`
}
type LoginByHTYZResp struct {
	Token    string `json:"token"`
	Uid      int    `json:"uid"`
	Role     int    `json:"role"`
	Username string `json:"username"`
}

// 公安一所免密登录
type LoginByGAYSReq struct {
	Token string `form:"token" json:"token" binding:"required"`
}
type LoginByGAYSResp struct {
	Token    string `json:"token"`
	Uid      int    `json:"uid"`
	Role     int    `json:"role"`
	Username string `json:"username"`
}

// PasswordCheckReq 用户登录
type PasswordCheckReq struct {
	UserId int `form:"userId" json:"userId" binding:"required"`
}
type PasswordCheckResp struct {
}

// UserLogin15SuoReq 用户登录
type UserLogin15SuoReq struct {
	Code string `form:"code" json:"code" binding:"required"`
}

type UserLogin15SuoResp struct {
	Token    string `json:"token"`
	Uid      int    `json:"uid"`
	Role     int    `json:"role"`
	Username string `json:"username"`
}
