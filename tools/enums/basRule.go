package enums

var BasRuleEnum basRule

type basRule struct {
}

// 规则分类
var (
	BasRuleClassAttemptedAdmin         = 1
	BasRuleClassAttemptedDos           = 2
	BasRuleClassAttemptedRecon         = 3
	BasRuleClassAttemptedUser          = 4
	BasRuleClassBadUnknown             = 5
	BasRuleClassCoinMining             = 6
	BasRuleClassDenialOfService        = 7
	BasRuleClassDomainC2               = 8
	BasRuleClassExploitKit             = 9
	BasRuleClassMiscAttack             = 10
	BasRuleClassNetworkScan            = 11
	BasRuleClassProtocolCommandDecode  = 12
	BasRuleClassShellcodeDetect        = 13
	BasRuleClassSuccessfulAdmin        = 14
	BasRuleClassSuccessfulDos          = 15
	BasRuleClassTargetedActivity       = 16
	BasRuleClassTrojanActivity         = 17
	BasRuleClassWebApplicationActivity = 18
	BasRuleClassWebApplicationAttack   = 19
)

func (b *basRule) AllClassEnum() map[int]string {
	res := map[int]string{
		BasRuleClassAttemptedAdmin:         "attempted-admin",
		BasRuleClassAttemptedDos:           "attempted-dos",
		BasRuleClassAttemptedRecon:         "attempted-recon",
		BasRuleClassAttemptedUser:          "attempted-user",
		BasRuleClassBadUnknown:             "bad-unknown",
		BasRuleClassCoinMining:             "coin-mining",
		BasRuleClassDenialOfService:        "denial-of-service",
		BasRuleClassDomainC2:               "domain-c2",
		BasRuleClassExploitKit:             "exploit-kit",
		BasRuleClassMiscAttack:             "misc-attack",
		BasRuleClassNetworkScan:            "network-scan",
		BasRuleClassProtocolCommandDecode:  "protocol-command-decode",
		BasRuleClassShellcodeDetect:        "shellcode-detect",
		BasRuleClassSuccessfulAdmin:        "successful-admin",
		BasRuleClassSuccessfulDos:          "successful-dos",
		BasRuleClassTargetedActivity:       "targeted-activity",
		BasRuleClassTrojanActivity:         "trojan-activity",
		BasRuleClassWebApplicationActivity: "web-application-activity",
		BasRuleClassWebApplicationAttack:   "web-application-attack",
	}

	return res
}

func (b *basRule) GetClassEnum(status int) string {
	return b.AllClassEnum()[status]
}

// 风险等级 1-高危、2-中危、3-低危、4-安全
var (
	BasRuleRiskHigh = 1
	BasRuleRiskMid  = 2
	BasRuleRiskLow  = 3
	BasRuleRiskSafe = 4
)

func (b *basRule) AllBasRuleRisk() map[int]string {
	return map[int]string{
		BasRuleRiskHigh: "高危",
		BasRuleRiskMid:  "中危",
		BasRuleRiskLow:  "低危",
		BasRuleRiskSafe: "安全",
	}
}

func (b *basRule) GetBasRuleRisk(risk int) string {
	return b.AllBasRuleRisk()[risk]
}

// 攻击阶段 信息收集、漏洞探测、漏洞检测、漏洞利用、后渗透
var (
	BasRuleAttackStageInfoCollect  = 1
	BasRuleAttackStageVulDetection = 2
	BasRuleAttackStageVulCheck     = 3
	BasRuleAttackStageVulUse       = 4
	BasRuleAttackStageBack         = 5
)

func (b *basRule) AllBasRuleAttackStage() map[int]string {
	return map[int]string{
		BasRuleAttackStageInfoCollect:  "信息收集",
		BasRuleAttackStageVulDetection: "漏洞探测",
		BasRuleAttackStageVulCheck:     "漏洞检测",
		BasRuleAttackStageVulUse:       "漏洞利用",
		BasRuleAttackStageBack:         "后渗透",
	}
}

func (b *basRule) GetBasRuleAttackStage(attackStage int) string {
	return b.AllBasRuleAttackStage()[attackStage]
}

const (
	BasResultStatusFailed  = 1
	BasResultStatusSuccess = 2
)

func (b *basRule) AllBasResultStatus() map[int]string {
	return map[int]string{
		BasResultStatusFailed:  "失败",
		BasResultStatusSuccess: "成功",
	}
}
