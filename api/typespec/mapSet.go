package typespec

// MapSetSystemSecurityReq 系统配置 修改安全设置请求
type MapSetSystemSecurityReq struct {
	SystemSecurityPasswordRank          int `form:"systemSecurityPasswordRank" json:"systemSecurityPasswordRank" dc:"密码等级(1=低：长度不低于8位，可为全字母或数字)(2=中：长度不低于8位，包含字母、符号、数字中的两类)(3=高：长度不低于8位，包含字母、符号、数字中的三类)" binding:"required"` // 密码强度 默认1
	SystemSecurityPasswordCycle         int `form:"systemSecurityPasswordCycle" json:"systemSecurityPasswordCycle" dc:"密码有效期 单位月 0不限"`                                                                                       // 密码有效期 单位 月 0不限
	SystemSecurityLoginTimeout          int `form:"systemSecurityLoginTimeout" json:"systemSecurityLoginTimeout" dc:"登录超时（多长时间未操作自动退出）单位分 0不限"`                                                                              // 用户长时间为操作自动退出系统 单位 分 0不限
	SystemSecurityAccountNoLoginDisable int `form:"systemSecurityAccountNoLoginDisable" json:"systemSecurityAccountNoLoginDisable" dc:"账号长时间未登录系统（到达时间直接禁用）单位月 0不限"`                                                         // 账号长时间未登录系统禁用账号 单位 月 0不限
	SystemSecurityUserLimit             int `form:"systemSecurityUserLimit" json:"systemSecurityUserLimit" dc:"普通用户限制登录次数（到达次数直接禁用）单位次 0不限"`                                                                                 // 普通账号连续登陆出错禁止登陆 单位 次 0 不限
	SystemSecurityAdminLimit            int `form:"systemSecurityAdminLimit" json:"systemSecurityAdminLimit" dc:"除用户员外限制登录次数（到达次数限制登录时间具体时间由systemSecurityAdminBanTime定义）单位次 0不限"`                                           // 管理/审计员连续登陆出错限制次数 单位 次 0不限
	SystemSecurityBanTime               int `form:"systemSecurityBanTime" json:"systemSecurityBanTime" dc:"除普通用户外禁止登录时间 单位小时 0不禁止"`                                                                                          // 管理/审计员连续登陆出错限制-禁止登陆时长 单位 小时 0不限
	SystemSecurityPasswordValid         int `form:"systemSecurityPasswordValid" json:"systemSecurityPasswordValid" dc:"密码有效期 单位月 0不限"`                                                                                       // 管理/审计员连续登陆出错限制-禁止登陆时长 单位 小时 0不限
	SystemSecurityExpireUnused          int `form:"systemSecurityExpireUnused" json:"systemSecurityExpireUnused" dc:"超期未使用 单位月 0不限"`                                                                                         // 管理/审计员连续登陆出错限制-禁止登陆时长 单位 小时 0不限
}
type MapSetSystemSecurityRes struct {
}

// SystemSecurityInfoReq 配置请求结构体
type SystemSecurityInfoReq struct {
}

// SecurityInfo 安全信息
type SecurityInfo struct {
	PasswordRank  int `json:"password_rank"`
	PasswordCycle int `json:"password_cycle"`
	LoginTimeout  int `json:"login_timeout"`
	UserLimit     int `json:"user_limit"`
	AdminLimit    int `json:"admin_limit"`
	BanTime       int `json:"ban_time"`
	PasswordValid int `json:"password_valid"`
	ExpireUnused  int `json:"expire_unused"`
}

// WarningInfo 告警信息
type WarningInfo struct {
	CpuThreshold    float64 `json:"cpu_threshold"`
	MemoryThreshold float64 `json:"memory_threshold"`
	FlowThreshold   float64 `json:"flow_threshold"`
	DiskThreshold   float64 `json:"disk_threshold"`
	IsOpen          bool    `json:"is_open"`
}

// SystemSecurityInfoRes 配置请求返回
type SystemSecurityInfoRes struct {
	Success  bool         `json:"success"`
	Security SecurityInfo `json:"security"`
	Warning  WarningInfo  `json:"warning"`
}
