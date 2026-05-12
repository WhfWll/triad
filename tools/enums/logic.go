package enums

const (
	_                      = iota
	LogicTaskStatusTrigger //待触发
	LogicTaskStatusBegin   //待执行
	LogicTaskStatusRunning //运行中
	LogicTaskStatusFinish  //已结束
	LogicTaskStatusPausing //暂停中
)

const (
	_                         = iota
	LogicTypeBeyondPermission //越权漏洞
	LogicTypeTraverseTesting  //敏感信息遍历
	LogicTypeUnAuthAccess     //未授权访问
)

const (
	_                      = iota
	LogicCredPatternCookie //越权漏洞
	LogicCredPatternHeader //敏感信息遍历
)

const (
	LogicTaskRiskHigh = 1 //高危
	LogicTaskRiskMid  = 2 //中危
	LogicTaskRiskLow  = 3 //低危
	LogicTaskRiskSafe = 4 //安全
)

const (
	LogicVulRiskDead = 1 //高危
	LogicVulRiskHigh = 2 //高危
	LogicVulRiskMid  = 3 //高危
	LogicVulRiskLow  = 4 //高危
)

type Logic struct {
}

func (l Logic) AllTypeEnum() map[int]string {
	enum := map[int]string{
		LogicTypeBeyondPermission: "越权漏洞",
		LogicTypeTraverseTesting:  "敏感信息遍历",
		LogicTypeUnAuthAccess:     "未授权访问",
	}
	return enum
}

func (l Logic) GetScanTypeName(scanType int) string {
	return l.AllTypeEnum()[scanType]
}

func (l Logic) AllCredPatternEnum() map[int]string {
	enum := map[int]string{
		LogicCredPatternCookie: "cookie",
		LogicCredPatternHeader: "header",
	}
	return enum
}

func (l Logic) GetCredPatternName(pattern int) string {
	return l.AllCredPatternEnum()[pattern]
}

func (l Logic) AllRiskEnum() map[int]string {
	enum := map[int]string{
		LogicTaskRiskHigh: "高危",
		LogicTaskRiskMid:  "中危",
		LogicTaskRiskLow:  "低危",
		LogicTaskRiskSafe: "安全",
	}
	return enum
}

func (l Logic) GetVulRiskName(risk int) string {
	return l.AllVulRiskEnum()[risk]
}

func (l Logic) AllVulRiskEnum() map[int]string {
	enum := map[int]string{
		LogicVulRiskDead: "致命",
		LogicVulRiskHigh: "高危",
		LogicVulRiskMid:  "中危",
		LogicVulRiskLow:  "低危",
	}
	return enum
}

func (l Logic) GetRiskName(risk int) string {
	return l.AllRiskEnum()[risk]
}

func (l Logic) AllStatusEnum() map[int]string {
	enum := map[int]string{
		LogicTaskStatusTrigger: "待触发",
		LogicTaskStatusBegin:   "待开始",
		LogicTaskStatusRunning: "运行中",
		LogicTaskStatusFinish:  "已完成",
	}
	return enum
}

func (l Logic) GetStatusName(status int) string {
	return l.AllStatusEnum()[status]
}

const (
	WhitePath     = "address,edit,id,info,list,member,mem,profile,select,service,user,update"
	BlackPath     = "infoleak,jarheads,sqli_id,truman,wide"
	Character     = "name"
	Number        = "account,id,no,number,user_code"
	JsonKeyword   = "account,add,article,bank,business,birth,card,city,company,content,driver,education,email"
	NoJsonKeyword = "姓名,性别,手机,住址,邮箱"
)

const (
	LogicTargetIsAliveN = 1 // 目标未存活
	LogicTargetIsAliveY = 2 // 目标存活
)

func (l Logic) AllIsAliveEnum() map[int]string {
	enum := map[int]string{
		LogicTargetIsAliveN: "不存活",
		LogicTargetIsAliveY: "存活",
	}
	return enum
}

func (l Logic) GetIsAliveName(isAlive int) string {
	return l.AllIsAliveEnum()[isAlive]
}

const (
	LogicCallIdCacheKey       = "smart:logic:task"
	LogicResultCacheKeyPreKey = "smart:logic:result:"
)

const (
	CrawlerConfigAllDomain = 0 //全域名扫描
	CrawlerConfigSubDomain = 1 //子域名扫描
)
const (
	LogicCrawlerConfigAllDomain = "AllDomainScan"
	LogicCrawlerConfigSubDomain = "SubMenuScan"
)

const (
	UnAuthVulName       = "未授权"
	UnAuthVulDesc       = "未授权漏洞，也常被称为权限绕过或认证绕过漏洞，是一种安全缺陷，允许攻击者在没有提供有效凭证或权限的情况下，访问或执行目标系统、应用、服务中的受保护资源。这类漏洞的存在使得系统的安全性大打折扣，因为它破坏了基本的身份验证机制，使得恶意用户能够如同合法用户一样操作，从而可能窃取敏感信息、篡改数据、执行非法操作等。"
	UnAuthVulFixSuggest = `确保所有对外提供的接口和服务在访问前都需要经过身份验证。\n 使用角色基础的访问控制（RBAC）或属性基础的访问控制（ABAC），确保每个用户或服务账户只能访问其被明确授权的资源。\n应用最小权限原则，确保每个组件、服务或用户只拥有完成其职责所需的最小权限集合。`
	UnAuthVulId         = "4dogs_1000000"
)
