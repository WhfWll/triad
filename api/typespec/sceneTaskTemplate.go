package typespec

import "smart/tools/enums"

// 场景管理 - 任务场景

// 场景枚举
type SceneEnumsRes struct {
	AliveProbe struct {
		AliveProbeType      []GlobalOptionsItemRes `json:"aliveProbeType"`      // 探活方式
		AliveProbePortRange []GlobalOptionsItemRes `json:"aliveProbePortRange"` // 探活端口范围
	} `json:"aliveProbe"` // 端口扫描
	PortScan struct {
		TcpScanType    []GlobalOptionsItemRes `json:"tcpScanType"`    // TCP扫描方式
		PortRange      []GlobalOptionsItemRes `json:"portRange"`      // 端口范围
		PortRangeValue []GlobalOptionsItemRes `json:"portRangeValue"` // 端口值
		Timeout        []GlobalOptionsItemRes `json:"timeout"`        // 端口扫描超时时间
		Concurrent     []GlobalOptionsItemRes `json:"concurrent"`     // 端口扫描并发数
	} `json:"portScan"` // 端口扫描
	Crawler struct {
		MaxDeep []GlobalOptionsItemRes `json:"maxDeep"` // 爬取深度 0不限
		MaxUrl  []GlobalOptionsItemRes `json:"maxUrl"`  // 最大爬取url数量 0不限
		/**
		`newcrawlerx.AllDomainScan` 表示爬取全域名 （默认0）
		`newcrawlerx.SubMenuScan` 表示爬取目标URL和子目录
		*/
		ScanRange   []GlobalOptionsItemRes `json:"scanRange"`   // 爬虫爬取范围
		Timeout     []GlobalOptionsItemRes `json:"timeout"`     // 爬虫单页面超时时间设置 0不限 单位秒
		FullTimeout []GlobalOptionsItemRes `json:"fullTimeout"` // 爬虫全局超时时间设置 0不限 单位秒
		/** 爬虫结果重复过滤设置
		0 不限，请勿设置爬虫过滤
		1 `newcrawlerx.UnLimitRepeat` 对page，method，query-name，query-value和post-data敏感
		2 `newcrawlerx.LowRepeatLevel` 对page，method，query-name和query-value敏感（默认）
		3 `newcrawlerx.MediumRepeatLevel` 对page，method和query-name敏感
		4 `newcrawlerx.HighRepeatLevel` 对page和method敏感
		5 `newcrawlerx.ExtremeRepeatLevel` 对page敏感
		*/
		ScanRepeat   []GlobalOptionsItemRes `json:"scanRepeat"`   // 爬虫结果重复过滤设置
		CrawlerSpeed []GlobalOptionsItemRes `json:"crawlerSpeed"` // 爬取速度
		BlackList    string                 `json:"blackList"`    // 爬虫黑名单
		WhiteList    string                 `json:"whiteList"`    // 爬虫白名单
	} `json:"crawler"` // 爬虫
	WebPathScan struct {
		Speed      []GlobalOptionsItemRes `json:"speed"`      // 猜测速率
		Times      []GlobalOptionsItemRes `json:"times"`      // 猜测时长
		ScanDict   []GlobalOptionsItemRes `json:"scanDict"`   // 路径字典
		TitleBlack string                 `json:"titleBlack"` // 排除标题黑名单
	} `json:"webPathScan"` // Web路径爆破
	WeakPass struct {
		Services       []GlobalOptionsItemRes `json:"services"`       // 服务 注意：这里返回的value是服务ID 也就是service字段 可以理解为类型
		Type           []GlobalOptionsItemRes `json:"type"`           // 字典类型
		CommonUserDict []GlobalOptionsItemRes `json:"commonUserDict"` // 通用用户字典 注意：这里返回的value是字典的ID，因为通用用户字典会有很多数据，如果返回服务将无法找到用户选择的哪条字典数据
		CommonPassDict []GlobalOptionsItemRes `json:"commonPassDict"` // 通用密码字典 注意：这里返回的value是字典的ID，因为通用用户字典会有很多数据，如果返回服务将无法找到用户选择的哪条字典数据
		OnlyUseAdd     bool                   `json:"onlyUseAdd"`     // 仅使用补充字典
		GuessNum       []GlobalOptionsItemRes `json:"guessNum"`       // 猜测次数
		GuessTimeout   []GlobalOptionsItemRes `json:"guessTimeout"`   // 猜测时间
		GuessRate      []GlobalOptionsItemRes `json:"guessRate"`      // 猜测速率
	} `json:"weakPass"` // 弱口令
	SubdomainDictCollect []GlobalOptionsItemRes `json:"subdomainDictCollect"` // 子域名收集
	LateralMove          struct {
	} `json:"lateralMove"`
}

// 列表
type SceneTaskTemplateListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
}
type SceneTaskTemplateListRes struct {
	Total int64                          `json:"total"`
	List  []SceneTaskTemplateListResItem `json:"list"`
}
type SceneTaskTemplateListResItem struct {
	Id            int    `json:"id"`
	TemplateName  string `json:"templateName"`  //场景模版名字
	Describe      string `json:"describe"`      //描述
	IsDefault     int    `json:"isDefault"`     //场景模版是否默认
	IsDefaultName string `json:"isDefaultName"` //场景模版是否默认
	UserId        int    `json:"userId"`        //场景模版创建人
	UserName      string `json:"userName"`      //场景模版创建人
	CreateTime    string `json:"createTime"`    //
}
type SceneTaskTemplateListResItemSource struct {
	Label string `json:"label"` //下拉菜单中的标签名
	Value uint   `json:"value"` //下拉菜单中标签名对应的数字
}

// 创建/编辑任务场景模版
type SceneTaskTemplateCreateReq struct {
	TaskTemplateId int              `form:"taskTemplateId" json:"taskTemplateId"`                // 场景模版Id
	TemplateName   string           `form:"templateName" json:"templateName" binding:"required"` // 场景模版名字
	Describe       string           `form:"describe" json:"describe"`                            // 场景描述
	Config         enums.ConfigJson `form:"config" json:"config" binding:"required"`
	UserId         int              `form:"userId" json:"userId"` //创建人id
}
type SceneTaskTemplateCreateRes struct {
	TaskTemplateId int `json:"task_template_id"`
}

// 详情
type SceneTaskTemplateDetailReq struct {
	TaskTemplateId int `form:"taskTemplateId" json:"taskTemplateId" binding:"required"`
}
type SceneTaskTemplateDetailRes struct {
	Id           int              `json:"id"`
	TemplateName string           `json:"templateName"`
	Describe     string           `json:"describe"`
	IsDefault    uint             `json:"isDefault"`
	UserId       int              `json:"userId"`
	Config       enums.ConfigJson `json:"config"`
	CreateTime   string           `json:"createTime"`
}

// 任务场景 - 拷贝
type SceneTaskTemplateCopyReq struct {
	TaskTemplateId int `form:"taskTemplateId" json:"taskTemplateId" binding:"required"`
	UserId         int `form:"userId" json:"userId"` //创建人id
}
type SceneTaskTemplateCopyRes struct {
	TaskTemplateId int `json:"taskTemplateId"`
}

// 任务场景 - 设置默认
type SceneTaskTemplateSetDefaultReq struct {
	TaskTemplateId int `form:"taskTemplateId" json:"taskTemplateId" binding:"required"`
}
type SceneTaskTemplateSetDefaultRes struct {
}

// 任务场景 - 删除
type SceneTaskTemplateDelReq struct {
	TaskTemplateIds string `form:"taskTemplateIds" json:"taskTemplateIds" binding:"required"`
}
type SceneTaskTemplateDelRes struct {
}

// 任务场景 - 任务中的模板下拉菜单
type SceneTaskTemplateToCheckTaskOptionsReq struct {
}
type SceneTaskTemplateToCheckTaskOptionsRes struct {
	TaskTemplate []TemplateParamsNode `json:"taskTemplate"`
	//TaskReportTemplate   []TemplateParamsNode `json:"taskReportTemplate"`
	//TargetReportTemplate []TemplateParamsNode `json:"targetReportTemplate"`
}
type TemplateParamsNode struct {
	Name      string `json:"name"`
	IsDefault uint   `json:"isDefault"`
	Id        int    `json:"id"`
	Describe  string `json:"describe"`
}

// 知识图谱
type GraphReq struct {
	TaskTemplateId int `form:"taskTemplateId" json:"taskTemplateId" binding:"required"`
}
type GraphRes struct {
	Nodes interface{} `json:"nodes"`
	Links interface{} `json:"links"`
}
