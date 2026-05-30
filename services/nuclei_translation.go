package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type nucleiZhEntry struct {
	NameZh        string `json:"name_zh"`
	DescriptionZh string `json:"description_zh"`
	FixZh         string `json:"fix_zh"`
}

type nucleiRegexReplacement struct {
	re *regexp.Regexp
	zh string
}

var (
	nucleiZhMapOnce sync.Once
	nucleiZhMap     map[string]nucleiZhEntry

	nucleiSentenceSpacesRE = regexp.MustCompile(`\s+`)
	nucleiHanRE            = regexp.MustCompile(`[\p{Han}]`)
	nucleiAsciiWordRE      = regexp.MustCompile(`[A-Za-z]{4,}`)

	nucleiEnsureDisabledRE         = regexp.MustCompile(`(?i)^ensures?\s+(.+?)\s+is\s+disabled\s+to\s+prevent\s+(.+?)\.?$`)
	nucleiEnsureEnabledRE          = regexp.MustCompile(`(?i)^ensures?\s+(.+?)\s+is\s+enabled\s+to\s+(.+?)\.?$`)
	nucleiTemplateDetectRE         = regexp.MustCompile(`(?i)^this template detects?\s+(.+?)\.?$`)
	nucleiDetectRE                 = regexp.MustCompile(`(?i)^detects?\s+(.+?)\.?$`)
	nucleiCheckWhetherRE           = regexp.MustCompile(`(?i)^checks?\s+whether\s+(.+?)\.?$`)
	nucleiUnrestrictedAccessRE     = regexp.MustCompile(`(?i)^unrestricted\s+(.+?)\s+access(?:\s+(?:in|on)\s+(.+))?$`)
	nucleiRestrictAccessRE         = regexp.MustCompile(`(?i)^restrict\s+(.+?)\s+access(?:\s+(?:in|on)\s+(.+))?$`)
	nucleiNotUsedRE                = regexp.MustCompile(`(?i)^(.+?)\s+is\s+not\s+used$`)
	nucleiNotConfiguredRE          = regexp.MustCompile(`(?i)^(.+?)\s+not\s+configured$`)
	nucleiNotEnabledOnRE           = regexp.MustCompile(`(?i)^(.+?)\s+not\s+enabled(?:\s+on\s+(.+))?$`)
	nucleiRemoveExpiredInRE        = regexp.MustCompile(`(?i)^remove\s+expired\s+(.+?)\s+in\s+(.+)$`)
	nucleiPublicAccessEnabledRE    = regexp.MustCompile(`(?i)^public\s+access\s+to\s+(.+?)\s+-\s+enabled$`)
	nucleiSupportMissingRE         = regexp.MustCompile(`(?i)^(.+?)\s+support\s+for\s+(.+?)\s+-\s+missing$`)
	nucleiForTargetStatusRE        = regexp.MustCompile(`(?i)^(.+?)\s+for\s+(.+?)\s+-\s+(disabled|enabled|missing|unconfigured)$`)
	nucleiForTargetTailStatusRE    = regexp.MustCompile(`(?i)^(.+?)\s+for\s+(.+?)\s+(disabled|enabled|missing|unconfigured)$`)
	nucleiUnconfiguredSuffixRE     = regexp.MustCompile(`(?i)^(.+?)\s+-\s+unconfigured$`)
	nucleiDisabledSuffixRE         = regexp.MustCompile(`(?i)^(.+?)\s+-\s+disabled$`)
	nucleiEnabledSuffixRE          = regexp.MustCompile(`(?i)^(.+?)\s+-\s+enabled$`)
	nucleiMissingSuffixRE          = regexp.MustCompile(`(?i)^(.+?)\s+-\s+missing$`)
	nucleiOutdatedSuffixRE         = regexp.MustCompile(`(?i)^(.+?)\s+-\s+outdated$`)
	nucleiDisabledTailRE           = regexp.MustCompile(`(?i)^(.+?)\s+disabled$`)
	nucleiEnabledTailRE            = regexp.MustCompile(`(?i)^(.+?)\s+enabled$`)
	nucleiEnabledOnRE              = regexp.MustCompile(`(?i)^(.+?)\s+enabled\s+on\s+(.+)$`)
	nucleiRequiresMinLengthRE      = regexp.MustCompile(`(?i)^(.+?)\s+requires\s+minimum\s+length\s+(\d+)\s+or\s+greater$`)
	nucleiRequiresAtLeastRE        = regexp.MustCompile(`(?i)^(.+?)\s+requires\s+at\s*least\s+one\s+(.+?)(?:\s+-\s+unconfigured)?$`)
	nucleiAllowedWithoutPasswordRE = regexp.MustCompile(`(?i)^(.+?)\s+allowed\s+without\s+password$`)
	nucleiCanRedirectRE            = regexp.MustCompile(`(?i)^(.+?)\s+can\s+redirect\s+(.+)$`)
	nucleiAllowsWithoutLoggingOnRE = regexp.MustCompile(`(?i)^(.+?)\s+allows\s+(.+?)\s+without\s+logging\s+on$`)
	nucleiAllowedToServersRE       = regexp.MustCompile(`(?i)^(.+?)\s+to\s+(.+?)\s+allowed$`)
	nucleiDashSuffixRE             = regexp.MustCompile(`^(.+?)\s+-\s+(.+)$`)
	nucleiUpgradeRE                = regexp.MustCompile(`(?i)^upgrade to (.+?)\.?$`)
	nucleiApplyPatchRE             = regexp.MustCompile(`(?i)^apply (?:the )?(?:vendor )?patch(?:es)?\.?$`)
	nucleiDisableRE                = regexp.MustCompile(`(?i)^disable\s+(.+?)\.?$`)
	nucleiRestrictRE               = regexp.MustCompile(`(?i)^restrict\s+(.+?)\.?$`)

	nucleiGlossary = []nucleiRegexReplacement{
		mustNucleiReplacement(`directory browsing detection`, "\u76ee\u5f55\u6d4f\u89c8\u68c0\u6d4b"),
		mustNucleiReplacement(`directory browsing`, "\u76ee\u5f55\u6d4f\u89c8"),
		mustNucleiReplacement(`remote code execution`, "\u8fdc\u7a0b\u4ee3\u7801\u6267\u884c"),
		mustNucleiReplacement(`remote command execution`, "\u8fdc\u7a0b\u547d\u4ee4\u6267\u884c"),
		mustNucleiReplacement(`local privilege escalation`, "\u672c\u5730\u63d0\u6743"),
		mustNucleiReplacement(`privilege escalation`, "\u6743\u9650\u63d0\u5347"),
		mustNucleiReplacement(`authentication bypass`, "\u8ba4\u8bc1\u7ed5\u8fc7"),
		mustNucleiReplacement(`authorization bypass`, "\u6388\u6743\u7ed5\u8fc7"),
		mustNucleiReplacement(`security bypass`, "\u5b89\u5168\u7ed5\u8fc7"),
		mustNucleiReplacement(`sql injection`, "SQL \u6ce8\u5165"),
		mustNucleiReplacement(`cross-site scripting`, "\u8de8\u7ad9\u811a\u672c"),
		mustNucleiReplacement(`cross site scripting`, "\u8de8\u7ad9\u811a\u672c"),
		mustNucleiReplacement(`server-side request forgery`, "\u670d\u52a1\u7aef\u8bf7\u6c42\u4f2a\u9020"),
		mustNucleiReplacement(`xml external entity`, "XML \u5916\u90e8\u5b9e\u4f53"),
		mustNucleiReplacement(`command injection`, "\u547d\u4ee4\u6ce8\u5165"),
		mustNucleiReplacement(`code injection`, "\u4ee3\u7801\u6ce8\u5165"),
		mustNucleiReplacement(`template injection`, "\u6a21\u677f\u6ce8\u5165"),
		mustNucleiReplacement(`local file inclusion`, "\u672c\u5730\u6587\u4ef6\u5305\u542b"),
		mustNucleiReplacement(`remote file inclusion`, "\u8fdc\u7a0b\u6587\u4ef6\u5305\u542b"),
		mustNucleiReplacement(`arbitrary file download`, "\u4efb\u610f\u6587\u4ef6\u4e0b\u8f7d"),
		mustNucleiReplacement(`arbitrary file read`, "\u4efb\u610f\u6587\u4ef6\u8bfb\u53d6"),
		mustNucleiReplacement(`file upload`, "\u6587\u4ef6\u4e0a\u4f20"),
		mustNucleiReplacement(`path traversal`, "\u8def\u5f84\u7a7f\u8d8a"),
		mustNucleiReplacement(`directory traversal`, "\u76ee\u5f55\u7a7f\u8d8a"),
		mustNucleiReplacement(`information disclosure`, "\u4fe1\u606f\u6cc4\u9732"),
		mustNucleiReplacement(`sensitive file disclosure`, "\u654f\u611f\u6587\u4ef6\u6cc4\u9732"),
		mustNucleiReplacement(`default credentials`, "\u9ed8\u8ba4\u51ed\u8bc1"),
		mustNucleiReplacement(`weak password`, "\u5f31\u53e3\u4ee4"),
		mustNucleiReplacement(`unauthorized access`, "\u672a\u6388\u6743\u8bbf\u95ee"),
		mustNucleiReplacement(`open redirect`, "\u5f00\u653e\u91cd\u5b9a\u5411"),
		mustNucleiReplacement(`denial of service`, "\u62d2\u7edd\u670d\u52a1"),
		mustNucleiReplacement(`ufida`, "\u7528\u53cb"),
		mustNucleiReplacement(`landray`, "\u84dd\u51cc"),
		mustNucleiReplacement(`weaver`, "\u6cdb\u5fae"),
		mustNucleiReplacement(`pan micro`, "\u6cdb\u5fae"),
		mustNucleiReplacement(`zentao`, "\u7985\u9053"),
		mustNucleiReplacement(`wuzhicms`, "\u4e94\u6307CMS"),
		mustNucleiReplacement(`hongjing`, "\u5b8f\u666f"),
		mustNucleiReplacement(`xintianqing`, "\u5929\u64ce"),
		mustNucleiReplacement(`human resource management system`, "\u4eba\u529b\u8d44\u6e90\u7ba1\u7406\u7cfb\u7edf"),
		mustNucleiReplacement(`e-office`, "E-Office"),
		mustNucleiReplacement(`e-cology`, "E-Cology"),
		mustNucleiReplacement(`xmirpcservlet`, "XmlRpcServlet"),
		mustNucleiReplacement(`hardening`, "\u52a0\u56fa"),
		mustNucleiReplacement(`misconfiguration`, "\u914d\u7f6e\u9519\u8bef"),
		mustNucleiReplacement(`configuration`, "\u914d\u7f6e"),
		mustNucleiReplacement(`disclosure`, "\u6cc4\u9732"),
		mustNucleiReplacement(`exposure`, "\u66b4\u9732"),
		mustNucleiReplacement(`bypass`, "\u7ed5\u8fc7"),
		mustNucleiReplacement(`detection`, "\u68c0\u6d4b"),
		mustNucleiReplacement(`vulnerability`, "\u6f0f\u6d1e"),
		mustNucleiReplacement(`unrestricted`, "\u672a\u53d7\u9650"),
		mustNucleiReplacement(`access`, "\u8bbf\u95ee"),
		mustNucleiReplacement(`not enabled`, "\u672a\u542f\u7528"),
		mustNucleiReplacement(`not configured`, "\u672a\u914d\u7f6e"),
		mustNucleiReplacement(`unconfigured`, "\u672a\u914d\u7f6e"),
		mustNucleiReplacement(`ssl/tls certificates`, "SSL/TLS \u8bc1\u4e66"),
		mustNucleiReplacement(`public access`, "\u516c\u7f51\u8bbf\u95ee"),
		mustNucleiReplacement(`api server`, "API Server"),
		mustNucleiReplacement(`ack`, "ACK"),
		mustNucleiReplacement(`cluster`, "\u96c6\u7fa4"),
		mustNucleiReplacement(`network policies`, "\u7f51\u7edc\u7b56\u7565"),
		mustNucleiReplacement(`global service`, "\u5168\u5c40\u670d\u52a1"),
		mustNucleiReplacement(`multi-region`, "\u8de8\u533a\u57df"),
		mustNucleiReplacement(`logging`, "\u65e5\u5fd7"),
		mustNucleiReplacement(`powershell`, "PowerShell"),
		mustNucleiReplacement(`os patches`, "OS \u8865\u4e01"),
		mustNucleiReplacement(`ram`, "RAM"),
		mustNucleiReplacement(`store passwords using reversible encryption`, "\u4f7f\u7528\u53ef\u9006\u52a0\u5bc6\u5b58\u50a8\u5bc6\u7801"),
		mustNucleiReplacement(`patches`, "\u8865\u4e01"),
		mustNucleiReplacement(`outdated`, "\u8fc7\u65f6"),
		mustNucleiReplacement(`encryption`, "\u52a0\u5bc6"),
		mustNucleiReplacement(`unattached disks`, "\u672a\u6302\u8f7d\u78c1\u76d8"),
		mustNucleiReplacement(`maximum password retry constraint policy`, "\u6700\u5927\u5bc6\u7801\u91cd\u8bd5\u7ea6\u675f\u7b56\u7565"),
		mustNucleiReplacement(`password policy expiration`, "\u5bc6\u7801\u7b56\u7565\u8fc7\u671f\u65f6\u95f4"),
		mustNucleiReplacement(`password policy reuse`, "\u5bc6\u7801\u91cd\u7528\u7b56\u7565"),
		mustNucleiReplacement(`password policy`, "\u5bc6\u7801\u7b56\u7565"),
		mustNucleiReplacement(`minimum length`, "\u6700\u5c0f\u957f\u5ea6"),
		mustNucleiReplacement(`lowercase`, "\u5c0f\u5199\u5b57\u6bcd"),
		mustNucleiReplacement(`number`, "\u6570\u5b57"),
		mustNucleiReplacement(`symbol`, "\u7b26\u53f7"),
		mustNucleiReplacement(`script block logging`, "Script Block \u65e5\u5fd7"),
		mustNucleiReplacement(`remote desktop connections`, "\u8fdc\u7a0b\u684c\u9762\u8fde\u63a5"),
		mustNucleiReplacement(`remote desktop users`, "\u8fdc\u7a0b\u684c\u9762\u7528\u6237"),
		mustNucleiReplacement(`remote desktop`, "\u8fdc\u7a0b\u684c\u9762"),
		mustNucleiReplacement(`connections`, "\u8fde\u63a5"),
		mustNucleiReplacement(`users`, "\u7528\u6237"),
		mustNucleiReplacement(`drives`, "\u9a71\u52a8\u5668"),
		mustNucleiReplacement(`network level authentication`, "\u7f51\u7edc\u7ea7\u522b\u8eab\u4efd\u9a8c\u8bc1"),
		mustNucleiReplacement(`rdp`, "RDP"),
		mustNucleiReplacement(`allowed`, "\u5141\u8bb8"),
		mustNucleiReplacement(`without password`, "\u65e0\u9700\u5bc6\u7801"),
		mustNucleiReplacement(`without logging on`, "\u672a\u767b\u5f55\u65f6"),
		mustNucleiReplacement(`non-server os`, "\u975e\u670d\u52a1\u5668\u64cd\u4f5c\u7cfb\u7edf"),
		mustNucleiReplacement(`remote assistance`, "\u8fdc\u7a0b\u534f\u52a9"),
		mustNucleiReplacement(`restrict anonymous access`, "\u9650\u5236\u533f\u540d\u8bbf\u95ee"),
		mustNucleiReplacement(`safe dll search mode`, "Safe DLL Search Mode"),
		mustNucleiReplacement(`secure boot`, "Secure Boot"),
		mustNucleiReplacement(`shutdown`, "\u5173\u673a"),
		mustNucleiReplacement(`system`, "\u7cfb\u7edf"),
		mustNucleiReplacement(`unencrypted passwords`, "\u672a\u52a0\u5bc6\u5bc6\u7801"),
		mustNucleiReplacement(`smb servers`, "SMB \u670d\u52a1\u5668"),
	}
)

func mustNucleiReplacement(pattern, zh string) nucleiRegexReplacement {
	return nucleiRegexReplacement{
		re: regexp.MustCompile(`(?i)` + regexp.QuoteMeta(pattern)),
		zh: zh,
	}
}

func localizeNucleiTemplateFields(pocname string, doc nucleiTemplateDoc, tags, refs []string) (string, string, string) {
	name := strings.TrimSpace(doc.Info.Name)
	desc := strings.TrimSpace(doc.Info.Description)
	fix := strings.TrimSpace(doc.Info.Remediation)

	if entry, ok := lookupNucleiZhEntry(pocname, doc.ID); ok {
		if strings.TrimSpace(entry.NameZh) != "" {
			name = strings.TrimSpace(entry.NameZh)
		}
		if strings.TrimSpace(entry.DescriptionZh) != "" {
			desc = strings.TrimSpace(entry.DescriptionZh)
		}
		if strings.TrimSpace(entry.FixZh) != "" {
			fix = strings.TrimSpace(entry.FixZh)
		}
	}

	name = autoTranslateNucleiName(name)
	desc = autoTranslateNucleiDescription(desc, name, tags)
	fix = autoTranslateNucleiRemediation(fix, name, tags)
	return name, desc, fix
}

func lookupNucleiZhEntry(pocname, templateID string) (nucleiZhEntry, bool) {
	nucleiZhMapOnce.Do(loadNucleiZhMap)

	keys := []string{
		strings.ToLower(filepath.ToSlash(strings.TrimSpace(pocname))),
		strings.ToLower(strings.TrimSpace(templateID)),
		strings.ToLower(filepath.Base(filepath.ToSlash(strings.TrimSpace(pocname)))),
	}
	for _, key := range keys {
		if key == "" {
			continue
		}
		if entry, ok := nucleiZhMap[key]; ok {
			return entry, true
		}
	}
	return nucleiZhEntry{}, false
}

func loadNucleiZhMap() {
	nucleiZhMap = map[string]nucleiZhEntry{}
	for _, candidate := range []string{
		filepath.Join("data", "nuclei_zh_map.json"),
		filepath.Join("..", "data", "nuclei_zh_map.json"),
	} {
		data, err := os.ReadFile(candidate)
		if err != nil {
			continue
		}
		var parsed map[string]nucleiZhEntry
		if json.Unmarshal(data, &parsed) != nil {
			continue
		}
		for key, entry := range parsed {
			nucleiZhMap[strings.ToLower(filepath.ToSlash(strings.TrimSpace(key)))] = entry
		}
		return
	}
}

func autoTranslateNucleiName(name string) string {
	if translated, ok := translateNucleiTitleByPattern(name); ok {
		name = translated
	}
	name = translateSecurityText(name)
	name = cleanupZhText(name)
	if name == "" {
		return "\u672a\u547d\u540d\u6f0f\u6d1e"
	}
	return name
}

func translateNucleiTitleByPattern(name string) (string, bool) {
	name = strings.TrimSpace(name)
	switch {
	case nucleiUnrestrictedAccessRE.MatchString(name):
		m := nucleiUnrestrictedAccessRE.FindStringSubmatch(name)
		subject := cleanupZhText(translateSecurityText(m[1]))
		scope := cleanupZhText(translateSecurityText(m[2]))
		if scope != "" {
			return scope + " \u4e0a\u7684 " + subject + " \u672a\u53d7\u9650\u8bbf\u95ee", true
		}
		return subject + " \u672a\u53d7\u9650\u8bbf\u95ee", true
	case nucleiRestrictAccessRE.MatchString(name):
		m := nucleiRestrictAccessRE.FindStringSubmatch(name)
		subject := cleanupZhText(translateSecurityText(m[1]))
		scope := cleanupZhText(translateSecurityText(m[2]))
		if scope != "" {
			return "\u9650\u5236 " + scope + " \u4e0a\u7684 " + subject + " \u8bbf\u95ee", true
		}
		return "\u9650\u5236 " + subject + " \u8bbf\u95ee", true
	case nucleiNotUsedRE.MatchString(name):
		m := nucleiNotUsedRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u672a\u4f7f\u7528"), true
	case nucleiNotConfiguredRE.MatchString(name):
		m := nucleiNotConfiguredRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u672a\u914d\u7f6e"), true
	case nucleiNotEnabledOnRE.MatchString(name):
		m := nucleiNotEnabledOnRE.FindStringSubmatch(name)
		subject := cleanupZhText(translateSecurityText(m[1]))
		scope := cleanupZhText(translateSecurityText(m[2]))
		if scope != "" {
			return scope + " \u672a\u542f\u7528 " + subject, true
		}
		return subject + " \u672a\u542f\u7528", true
	case nucleiRemoveExpiredInRE.MatchString(name):
		m := nucleiRemoveExpiredInRE.FindStringSubmatch(name)
		target := cleanupZhText(translateSecurityText(m[1]))
		scope := cleanupZhText(translateSecurityText(m[2]))
		return scope + " \u4e2d\u5b58\u5728\u8fc7\u671f\u7684 " + target, true
	case nucleiPublicAccessEnabledRE.MatchString(name):
		m := nucleiPublicAccessEnabledRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u5df2\u5f00\u542f\u516c\u7f51\u8bbf\u95ee"), true
	case nucleiSupportMissingRE.MatchString(name):
		m := nucleiSupportMissingRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		right := cleanupZhText(translateSecurityText(m[2]))
		return cleanupZhText(left + " \u7f3a\u5c11 " + right + " \u652f\u6301"), true
	case nucleiEnabledOnRE.MatchString(name):
		m := nucleiEnabledOnRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		right := cleanupZhText(translateSecurityText(m[2]))
		return cleanupZhText("\u5728 " + right + " \u4e0a\u542f\u7528\u4e86 " + left), true
	case nucleiRequiresMinLengthRE.MatchString(name):
		m := nucleiRequiresMinLengthRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		return cleanupZhText(left + " \u8981\u6c42\u6700\u5c0f\u957f\u5ea6\u4e3a " + m[2] + " \u6216\u66f4\u5927"), true
	case nucleiRequiresAtLeastRE.MatchString(name):
		m := nucleiRequiresAtLeastRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		right := cleanupZhText(translateSecurityText(m[2]))
		return cleanupZhText(left + " \u8981\u6c42\u81f3\u5c11\u5305\u542b\u4e00\u4e2a " + right), true
	case nucleiForTargetStatusRE.MatchString(name):
		m := nucleiForTargetStatusRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		right := cleanupZhText(translateSecurityText(m[2]))
		status := map[string]string{
			"disabled":     "\u5df2\u7981\u7528",
			"enabled":      "\u5df2\u542f\u7528",
			"missing":      "\u7f3a\u5931",
			"unconfigured": "\u672a\u914d\u7f6e",
		}[strings.ToLower(m[3])]
		return cleanupZhText(right + " \u7684 " + left + " " + status), true
	case nucleiForTargetTailStatusRE.MatchString(name):
		m := nucleiForTargetTailStatusRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		right := cleanupZhText(translateSecurityText(m[2]))
		status := map[string]string{
			"disabled":     "\u5df2\u7981\u7528",
			"enabled":      "\u5df2\u542f\u7528",
			"missing":      "\u7f3a\u5931",
			"unconfigured": "\u672a\u914d\u7f6e",
		}[strings.ToLower(m[3])]
		return cleanupZhText(right + " \u7684 " + left + " " + status), true
	case nucleiUnconfiguredSuffixRE.MatchString(name):
		m := nucleiUnconfiguredSuffixRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u672a\u914d\u7f6e"), true
	case nucleiAllowedWithoutPasswordRE.MatchString(name):
		m := nucleiAllowedWithoutPasswordRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		return cleanupZhText(left + " \u5141\u8bb8\u65e0\u9700\u5bc6\u7801"), true
	case nucleiCanRedirectRE.MatchString(name):
		m := nucleiCanRedirectRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		right := cleanupZhText(translateSecurityText(m[2]))
		return cleanupZhText(left + " \u53ef\u91cd\u5b9a\u5411 " + right), true
	case nucleiAllowsWithoutLoggingOnRE.MatchString(name):
		m := nucleiAllowsWithoutLoggingOnRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		right := cleanupZhText(translateSecurityText(m[2]))
		return cleanupZhText(left + " \u5141\u8bb8\u5728\u672a\u767b\u5f55\u65f6" + right), true
	case nucleiAllowedToServersRE.MatchString(name):
		m := nucleiAllowedToServersRE.FindStringSubmatch(name)
		left := cleanupZhText(translateSecurityText(m[1]))
		right := cleanupZhText(translateSecurityText(m[2]))
		return cleanupZhText("\u5141\u8bb8\u5411 " + right + " \u53d1\u9001 " + left), true
	case nucleiDisabledSuffixRE.MatchString(name):
		m := nucleiDisabledSuffixRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u5df2\u7981\u7528"), true
	case nucleiEnabledSuffixRE.MatchString(name):
		m := nucleiEnabledSuffixRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u5df2\u542f\u7528"), true
	case nucleiMissingSuffixRE.MatchString(name):
		m := nucleiMissingSuffixRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u7f3a\u5931"), true
	case nucleiOutdatedSuffixRE.MatchString(name):
		m := nucleiOutdatedSuffixRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u8fc7\u65f6"), true
	case nucleiDisabledTailRE.MatchString(name):
		m := nucleiDisabledTailRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u5df2\u7981\u7528"), true
	case nucleiEnabledTailRE.MatchString(name):
		m := nucleiEnabledTailRE.FindStringSubmatch(name)
		return cleanupZhText(translateSecurityText(m[1]) + " \u5df2\u542f\u7528"), true
	case nucleiDashSuffixRE.MatchString(name):
		m := nucleiDashSuffixRE.FindStringSubmatch(name)
		left := strings.TrimSpace(m[1])
		right := strings.TrimSpace(m[2])
		rightZh := cleanupZhText(translateSecurityText(right))
		if rightZh != "" && rightZh != right {
			return left + " - " + rightZh, true
		}
	}
	return "", false
}

func autoTranslateNucleiDescription(desc, localizedName string, tags []string) string {
	desc = strings.TrimSpace(desc)
	switch {
	case desc == "":
		return buildGenericChineseDescription(localizedName, tags)
	case nucleiEnsureDisabledRE.MatchString(desc):
		m := nucleiEnsureDisabledRE.FindStringSubmatch(desc)
		return cleanupZhText("\u786e\u4fdd\u7981\u7528 " + translateSecurityText(m[1]) + "\uff0c\u4ee5\u9632\u6b62 " + translateSecurityText(m[2]) + "\u3002")
	case nucleiEnsureEnabledRE.MatchString(desc):
		m := nucleiEnsureEnabledRE.FindStringSubmatch(desc)
		return cleanupZhText("\u786e\u4fdd\u542f\u7528 " + translateSecurityText(m[1]) + "\uff0c\u4ee5\u4fbf " + translateSecurityText(m[2]) + "\u3002")
	case nucleiTemplateDetectRE.MatchString(desc):
		m := nucleiTemplateDetectRE.FindStringSubmatch(desc)
		return cleanupZhText("\u8be5\u6a21\u677f\u7528\u4e8e\u68c0\u6d4b " + translateSecurityText(m[1]) + "\u3002")
	case nucleiDetectRE.MatchString(desc):
		m := nucleiDetectRE.FindStringSubmatch(desc)
		return cleanupZhText("\u7528\u4e8e\u68c0\u6d4b " + translateSecurityText(m[1]) + "\u3002")
	case nucleiCheckWhetherRE.MatchString(desc):
		m := nucleiCheckWhetherRE.FindStringSubmatch(desc)
		return cleanupZhText("\u7528\u4e8e\u68c0\u67e5\u662f\u5426 " + translateSecurityText(m[1]) + "\u3002")
	default:
		desc = cleanupZhText(translateSecurityText(desc))
		if looksStillEnglish(desc) {
			return buildGenericChineseDescription(localizedName, tags)
		}
		return desc
	}
}

func autoTranslateNucleiRemediation(fix, localizedName string, tags []string) string {
	fix = strings.TrimSpace(fix)
	switch {
	case fix == "":
		return buildGenericChineseFix(localizedName, tags)
	case nucleiUpgradeRE.MatchString(fix):
		m := nucleiUpgradeRE.FindStringSubmatch(fix)
		return cleanupZhText("\u8bf7\u5347\u7ea7\u5230 " + translateSecurityText(m[1]) + "\u3002")
	case nucleiApplyPatchRE.MatchString(fix):
		return "\u8bf7\u5c3d\u5feb\u5e94\u7528\u5382\u5546\u8865\u4e01\u5e76\u5b8c\u6210\u590d\u6d4b\u3002"
	case nucleiDisableRE.MatchString(fix):
		m := nucleiDisableRE.FindStringSubmatch(fix)
		return cleanupZhText("\u8bf7\u7981\u7528 " + translateSecurityText(m[1]) + "\u3002")
	case nucleiRestrictRE.MatchString(fix):
		m := nucleiRestrictRE.FindStringSubmatch(fix)
		return cleanupZhText("\u8bf7\u9650\u5236 " + translateSecurityText(m[1]) + "\u3002")
	default:
		fix = cleanupZhText(translateSecurityText(fix))
		if looksStillEnglish(fix) {
			return buildGenericChineseFix(localizedName, tags)
		}
		return fix
	}
}

func buildGenericChineseDescription(localizedName string, tags []string) string {
	name := cleanupZhText(localizedName)
	if name == "" {
		name = "\u8be5\u6a21\u677f"
	}
	return "\u7528\u4e8e\u68c0\u6d4b " + name + " \u76f8\u5173\u95ee\u9898\u3002"
}

func buildGenericChineseFix(localizedName string, tags []string) string {
	switch {
	case strings.Contains(localizedName, "\u76ee\u5f55\u6d4f\u89c8"):
		return "\u8bf7\u5173\u95ed\u76ee\u5f55\u6d4f\u89c8\u529f\u80fd\uff0c\u4ec5\u5f00\u653e\u5fc5\u8981\u76ee\u5f55\u4e0e\u6587\u4ef6\u8bbf\u95ee\uff0c\u5e76\u6309\u6700\u5c0f\u66b4\u9732\u539f\u5219\u8c03\u6574\u7ad9\u70b9\u914d\u7f6e\u3002"
	case strings.Contains(localizedName, "SQL \u6ce8\u5165"):
		return "\u8bf7\u5347\u7ea7\u5230\u5b89\u5168\u7248\u672c\uff0c\u4f7f\u7528\u53c2\u6570\u5316\u67e5\u8be2\uff0c\u5e76\u52a0\u5f3a\u8f93\u5165\u6821\u9a8c\u4e0e\u6743\u9650\u63a7\u5236\u3002"
	case strings.Contains(localizedName, "\u8de8\u7ad9\u811a\u672c"):
		return "\u8bf7\u5347\u7ea7\u5230\u5b89\u5168\u7248\u672c\uff0c\u5e76\u5bf9\u8f93\u51fa\u5185\u5bb9\u8fdb\u884c\u8f6c\u4e49\u4e0e\u8fc7\u6ee4\uff0c\u5b8c\u5584\u8f93\u5165\u6821\u9a8c\u3002"
	case strings.Contains(localizedName, "\u672a\u6388\u6743\u8bbf\u95ee"), strings.Contains(localizedName, "\u8ba4\u8bc1\u7ed5\u8fc7"), strings.Contains(localizedName, "\u6388\u6743\u7ed5\u8fc7"):
		return "\u8bf7\u5347\u7ea7\u5230\u5b89\u5168\u7248\u672c\uff0c\u5f3a\u5316\u8ba4\u8bc1\u9274\u6743\u903b\u8f91\uff0c\u5e76\u9650\u5236\u654f\u611f\u63a5\u53e3\u5bf9\u5916\u66b4\u9732\u3002"
	case strings.Contains(localizedName, "\u8fdc\u7a0b\u4ee3\u7801\u6267\u884c"), strings.Contains(localizedName, "\u8fdc\u7a0b\u547d\u4ee4\u6267\u884c"), strings.Contains(localizedName, "\u547d\u4ee4\u6ce8\u5165"), strings.Contains(localizedName, "\u4ee3\u7801\u6ce8\u5165"):
		return "\u8bf7\u5c3d\u5feb\u5347\u7ea7\u5230\u5b89\u5168\u7248\u672c\uff0c\u5e94\u7528\u5b98\u65b9\u8865\u4e01\uff0c\u5e76\u9650\u5236\u5371\u9669\u63a5\u53e3\u4e0e\u6267\u884c\u80fd\u529b\u3002"
	default:
		return "\u5efa\u8bae\u6839\u636e\u5b98\u65b9\u5b89\u5168\u516c\u544a\u3001\u8865\u4e01\u8bf4\u660e\u6216\u4ea7\u54c1\u52a0\u56fa\u6307\u5357\u8fdb\u884c\u4fee\u590d\uff0c\u5e76\u5728\u4fee\u590d\u540e\u5b8c\u6210\u590d\u6d4b\u3002"
	}
}

func translateSecurityText(text string) string {
	out := strings.TrimSpace(text)
	out = strings.ReplaceAll(out, "Access Analyzer", "@@AA@@")
	out = strings.ReplaceAll(out, "access analyzer", "@@AA@@")
	for _, item := range nucleiGlossary {
		out = item.re.ReplaceAllString(out, item.zh)
	}
	out = strings.ReplaceAll(out, "@@AA@@", "Access Analyzer")
	return cleanupZhText(out)
}

func cleanupZhText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	text = nucleiSentenceSpacesRE.ReplaceAllString(text, " ")
	text = strings.ReplaceAll(text, "'s", " 的")
	text = strings.ReplaceAll(text, "’s", " 的")
	replacer := strings.NewReplacer(
		" ,", ",",
		" .", ".",
		" :", ":",
		" ;", ";",
		" )", ")",
		"( ", "(",
		" - ", " - ",
	)
	text = replacer.Replace(text)
	return strings.TrimSpace(text)
}

func looksStillEnglish(text string) bool {
	if text == "" || nucleiHanRE.MatchString(text) {
		return false
	}
	return len(nucleiAsciiWordRE.FindAllString(text, -1)) >= 2
}
