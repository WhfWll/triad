package typespec

// BAS规则枚举
type BasRuleEnumRes struct {
	Class       []GlobalOptionsItemRes `json:"class"`
	AttackType  []GlobalOptionsItemRes `json:"attackType"`
	AttackStage []GlobalOptionsItemRes `json:"attackStage"`
	RiskLevel   []GlobalOptionsItemRes `json:"riskLevel"`
	Status      []GlobalOptionsItemRes `json:"status"`
}

// BAS规则列表
type BasRuleListReq struct {
	Page        int    `form:"page" json:"page" binding:"required"`
	Size        int    `form:"size" json:"size" binding:"required"`
	ClassType   int    `form:"classType" json:"classType"`
	Search      string `form:"search" json:"search"` // Bas名称搜索
	RiskLevel   int    `form:"riskLevel" json:"riskLevel"`
	AttackStage int    `form:"attackStage" json:"attackStage"`
	AttackType  string `form:"attackType" json:"attackType"`
	Ids         string `form:"ids" json:"ids"`
}
type BasRuleListRes struct {
	Page  int               `json:"page" binding:"required"`
	Size  int               `json:"size" binding:"required"`
	Total int64             `json:"total"`
	List  []BasRuleListItem `json:"list"`
}
type BasRuleListItem struct {
	Id                             int    `json:"id"`
	Content                        string `json:"content"`
	Name                           string `json:"name"`
	NameZh                         string `json:"nameZh"`
	Hash                           string `json:"hash"`
	ClassType                      string `json:"classType"`
	Protocol                       string `json:"protocol"`
	Keywords                       string `json:"keywords"`
	KeywordsZh                     string `json:"keywordsZh"`
	Description                    string `json:"description"`
	DescriptionZh                  string `json:"descriptionZh"`
	Cve                            string `json:"cve"`
	RawTrafficBeyondIpPacketBase64 string `json:"rawTrafficBeyondIpPacketBase64"`
	RawTrafficBeyondHttpBase64     string `json:"rawTrafficBeyondHttpBase64"`
	AttackMode                     int    `json:"attackMode"`
	AttackModeEnum                 string `json:"attackModeEnum"`
	AttackStage                    int    `json:"attackStage"`
	AttackStageEnum                string `json:"attackStageEnum"`
	RiskLevel                      int    `json:"riskLevel"`
	RiskLevelEnum                  string `json:"riskLevelEnum"`
	AffectTarget                   string `json:"affectTarget"`
	AffectScore                    string `json:"affectScore"`
	RelationAttackMethod           string `json:"relationAttackMethod"`
	RefUrl                         string `json:"refUrl"`
	FixSuggest                     string `json:"fixSuggest"`
	CreateTime                     string `json:"createTime"`
	UpdateTime                     string `json:"updateTime"`
}

type BasRuleInfoReq struct {
	BasRuleId int `form:"basRuleId" json:"basRuleId" binding:"required"`
}
type BasRuleInfoResp struct {
	Id                             int    `json:"id"`
	Content                        string `json:"content"`
	Name                           string `json:"name"`
	NameZh                         string `json:"nameZh"`
	Hash                           string `json:"hash"`
	ClassType                      string `json:"classType"`
	Protocol                       string `json:"protocol"`
	Keywords                       string `json:"keywords"`
	KeywordsZh                     string `json:"keywordsZh"`
	Description                    string `json:"description"`
	DescriptionZh                  string `json:"descriptionZh"`
	Cve                            string `json:"cve"`
	RawTrafficBeyondIpPacketBase64 string `json:"rawTrafficBeyondIpPacketBase64"`
	RawTrafficBeyondHttpBase64     string `json:"rawTrafficBeyondHttpBase64"`
	AttackMode                     int    `json:"attackMode"`
	AttackModeEnum                 string `json:"attackModeEnum"`
	AttackStage                    int    `json:"attackStage"`
	AttackStageEnum                string `json:"attackStageEnum"`
	RiskLevel                      int    `json:"riskLevel"`
	RiskLevelEnum                  string `json:"riskLevelEnum"`
	AffectTarget                   string `json:"affectTarget"`
	AffectScore                    string `json:"affectScore"`
	RelationAttackMethod           string `json:"relationAttackMethod"`
	RefUrl                         string `json:"refUrl"`
	FixSuggest                     string `json:"fixSuggest"`
	CreateTime                     string `json:"createTime"`
	UpdateTime                     string `json:"updateTime"`
}

// Bas 规则编辑
type BasRuleEditReq struct {
	Id                   int    `form:"id" json:"id" binding:"required"`
	AffectTarget         string `form:"affectTarget" json:"affectTarget" binding:"required"`
	AttackMode           int    `form:"attackMode" json:"attackMode"`
	AttackStage          int    `form:"attackStage" json:"attackStage"`
	RiskLevel            int    `form:"riskLevel" json:"riskLevel"`
	AffectScore          string `form:"affectScore" json:"affectScore" binding:"required"`
	RelationAttackMethod string `form:"relationAttackMethod" json:"relationAttackMethod" binding:"required"`
	FixSuggest           string `form:"fixSuggest" json:"fixSuggest" binding:"required"`
	RefUrl               string `form:"refUrl" json:"refUrl" binding:"required"`
}
type BasRuleEditRes struct {
}

// 剧本集创建
type BasTemplateCreateReq struct {
	Id      int    `form:"id" json:"id"`
	Name    string `form:"name" json:"name" binding:"required"`       // 方案名称
	Desc    string `form:"desc" json:"desc" binding:"required"`       // 方案描述
	RuleIds []int  `form:"ruleIds" json:"ruleIds" binding:"required"` // 规则ID集合
}
type BasTemplateCreateRes struct {
}

// 剧本集 列表
type BasTemplateListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"` // 剧本名搜索
}
type BasTemplateListRes struct {
	Page  int                   `form:"page" json:"page" binding:"required"`
	Size  int                   `form:"size" json:"size" binding:"required"`
	Total int64                 `json:"total"`
	List  []BasTemplateListItem `json:"list"`
}
type BasTemplateListItem struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	IsDefault     int    `json:"isDefault"`
	IsDefaultEnum string `json:"isDefaultEnum"`
	CreateTime    string `json:"createTime"`
}

// 剧本详情
type BasTemplateGetReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type BasTemplateGetRes struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	RuleIds    []int  `json:"ruleIds"`
	UpdateTime string `json:"updateTime"`
	CreateTime string `json:"createTime"`
}

// 剧本删除
type BasTemplateDelReq struct {
	Id string `form:"id" json:"id" binding:"required"`
}
type BasTemplateDelRes struct {
}

// 剧本集设置是否默认
type BasTemplateSetDefaultReq struct {
	Id int `form:"id" json:"id"`
}
type BasTemplateSetDefaultRes struct {
}

// agent 是否在线
type BasAgentIsOnlineReq struct {
	Host string `form:"host" json:"host" binding:"required"`
	Port string `form:"port" json:"port" binding:"required"`
}
type BasAgentIsOnlineRes struct {
}

// bas 任务创建
type BasTaskCreateReq struct {
	Name          string `form:"name" json:"name" binding:"required"`
	BasTemplateId int    `form:"basTemplateId" json:"basTemplateId" binding:"required"`
	BasNodeIds    []int  `form:"basNodeIds" json:"basNodeIds" binding:"required"`
}

type BasTaskCreateResp struct {
}

// bas 任务列表
type BasTaskGetReq struct {
	Page      int    `form:"page" json:"page" binding:"required"`
	Size      int    `form:"size" json:"size" binding:"required"`
	RiskLevel int    `form:"riskLevel" json:"riskLevel"`
	Search    string `form:"search" json:"search"` // 剧本名搜索
}
type BasTaskGetRes struct {
	Page  int              `json:"page"`
	Size  int              `json:"size"`
	Total int64            `json:"total"`
	List  []BasTaskGetItem `json:"list"`
}
type BasTaskGetItem struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	TemplateId    int    `json:"templateId"`
	TemplateName  string `json:"templateName"`
	CreateTime    string `json:"createTime"`
	RiskLevel     int    `json:"riskLevel"`
	RiskLevelEnum string `json:"riskLevelEnum"`
	Status        int    `json:"status"`
	StatusEnum    string `json:"statusEnum"`
	UserId        int    `json:"userId"`
	UserName      string `json:"userName"`
}

// BAS 任务结束
type BasTaskEndReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type BasTaskEndRes struct {
}

// BAS 任务删除
type BasTaskDelReq struct {
	Id string `form:"id" json:"id" binding:"required"`
}
type BasTaskDelRes struct {
}

// BAS 任务详情
type BasTaskDetailReq struct {
	Id     int    `form:"id" json:"id" binding:"required"`
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"` // 剧本名搜索
}
type BasTaskDetailRes struct {
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
	Total int64               `json:"total"`
	List  []BasTaskDetailItem `json:"list"`
}
type BasTaskDetailItem struct {
	Id            int    `json:"id"`
	Addr          string `json:"addr"`
	RiskLevel     int    `json:"riskLevel"`
	RiskLevelEnum string `json:"riskLevelEnum"`
	HighNum       int    `json:"highNum"`
	MidNum        int    `json:"midNum"`
	LowNum        int    `json:"lowNum"`
	SafeNum       int    `json:"safeNum"`
	Status        int    `json:"status"`
	StatusEnum    string `json:"statusEnum"`
	Create        string `json:"create"`
}

// BAS 任务目标日志
type BasTargetLogReq struct {
	Id int `form:"id" json:"id" binding:"required"`
}
type BasTargetLogRes struct {
	Content []string `json:"content"`
}

// BAS 任务目标删除
type BasTargetDelReq struct {
	Id string `form:"id" json:"id" binding:"required"`
}
type BasTargetDelRes struct {
}

type BasEnumRes struct {
	TaskStatus interface{} `json:"taskStatus"`
	RiskLevel  interface{} `json:"riskLevel"`
	VulStatus  interface{} `json:"vulStatus"`
}

type BasReceivResultReq struct {
	Type    string `form:"type" json:"type" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

type BasReceivResultResp struct {
}

type BasVulStatReq struct {
	BasTaskId int `form:"basTaskId" json:"basTaskId" binding:"required"`
}

type BasVulStatResp struct {
	Status    [2]int64 `json:"status"`
	RiskLevel [4]int64 `json:"riskLevel"`
}

type BasVulListReq struct {
	BasTaskId   int    `form:"basTaskId" json:"basTaskId" binding:"required"`
	RiskLevel   int    `form:"riskLevel" json:"riskLevel"`
	AttackStage int    `form:"attackStage" json:"attackStage"`
	Search      string `form:"search" json:"search"`
	Page        int    `form:"page" json:"page" binding:"required"`
	Size        int    `form:"size" json:"size" binding:"required"`
	AttackMode  string `form:"attackMode" json:"attackMode"`
	Status      int    `form:"status" json:"status"`
}

type BasVulListResp struct {
	Total int64                 `json:"total"`
	List  []BasVulListRespItems `json:"list"`
}

type BasVulListRespItems struct {
	Id              int    `json:"id"`
	RuleID          int    `json:"ruleId"`
	Addr            string `json:"addr"`
	RuleName        string `json:"ruleName"`
	AttackMode      int    `json:"attackMode"`
	AttackModeName  string `json:"attackModeName"`
	AttackStage     int    `json:"attackStage"`
	AttackStageName string `json:"attackStageName"`
	RiskLevel       int    `json:"riskLevel"`
	RiskLevelName   string `json:"riskLevelName"`
	Status          int    `json:"status"`
	StatusName      string `json:"statusName"`
}

type BasVulDelReq struct {
	BasVulIds string `form:"basVulIds" json:"basVulIds" biding:"required"`
	BasTaskId int    `form:"basTaskId" json:"basTaskId" biding:"required"`
}

type BasVulDelResp struct{}

// Bas agent下载
type BasAgentDownloadReq struct {
	Platform     string `form:"platform" json:"platform" biding:"required"`
	GetTempToken bool   `form:"getTempToken" json:"getTempToken"`
}
type BasAgentDownloadRes struct {
	TempToken string `form:"tempToken" json:"tempToken"`
}

// Bas agent列表
type BasAgentListReq struct {
	Search string `form:"search" json:"search"`
	Page   int    `form:"page" json:"page" biding:"required"`
	Size   int    `form:"size" json:"size" biding:"required"`
}
type BasAgentListRes struct {
	Total int                `json:"total"`
	List  []BasAgentListItem `json:"list"`
}
type BasAgentListItem struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Ip               string `json:"ip"`
	OnlineStatus     int    `json:"onlineStatus"`
	OnlineStatusEnum string `json:"onlineStatusEnum"`
	Status           int    `json:"status"`
	StatusEnum       string `json:"statusEnum"`
}

type BasAgentLiveResp struct {
	List []BasAgentLiveRespItem `json:"list"`
}

type BasAgentLiveRespItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

// Bas Agent状态修改
type BasAgentStatusEditReq struct {
	Id     int `form:"id" json:"id" binding:"required"`
	Status int `form:"status" json:"status" binding:"required"`
}
type BasAgentStatusEditRes struct {
}
