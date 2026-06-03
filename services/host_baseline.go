package services

import (
	"context"
	"fmt"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strings"
	"sync"
	"time"
)

const maxBaselineActualValueLen = 16000

func truncateBaselineField(s string) string {
	if len(s) <= maxBaselineActualValueLen {
		return s
	}
	return s[:maxBaselineActualValueLen-20] + "\n...(truncated)"
}

type HostBaselineChecker struct {
	connManager *HostConnManager
	ruleEngine  *BaselineEngine
}

var (
	globalHostBaselineChecker *HostBaselineChecker
	hostBaselineOnce          sync.Once
)

func GetHostBaselineChecker() *HostBaselineChecker {
	hostBaselineOnce.Do(func() {
		globalHostBaselineChecker = &HostBaselineChecker{
			connManager: GetHostConnManager(),
			ruleEngine:  GetBaselineEngine(),
		}
	})
	return globalHostBaselineChecker
}

type BaselineCheckTask struct {
	TaskID         int
	TargetID       int
	Host           string
	Port           int
	Username       string
	Password       string
	Key            string
	OSType         int
	Transport      int  // 0=auto锛堢敱杩炴帴灞傝В鏋愶級锛?=SSH锛?=WinRM
	WinRMUseHttps  bool
	ScanScene      int // 1=瀹夊叏閰嶇疆鏍告煡 2=涓绘満婕忔礊妫€娴?
	RuleProgress   func(item typespec.BatchRuleResult)
}

type BaselineCheckReport struct {
	TaskID          int
	TargetID        int
	TargetIP        string
	OSType          int
	TotalRules      int
	PassCount       int
	FailCount       int
	ErrorCount      int
	CriticalCount   int
	HighRiskCount   int
	MiddleRiskCount int
	LowRiskCount    int
	ComplianceScore float64
	Results         []mysqls.BaselineCheckResult
	StartTime       time.Time
	EndTime         time.Time
}

func (c *HostBaselineChecker) RunBaselineCheck(ctx context.Context, task *BaselineCheckTask) (*BaselineCheckReport, error) {
	scene := task.ScanScene
	if scene <= 0 {
		scene = enums.HostScanSceneBaseline
	}
	fmt.Printf(">>> RunBaselineCheck: start taskID=%d target=%s osType=%d port=%d user=%s transport=%d\n",
		task.TaskID, task.Host, task.OSType, task.Port, task.Username, task.Transport)

	report := &BaselineCheckReport{
		TaskID:    task.TaskID,
		TargetID:  task.TargetID,
		TargetIP:  task.Host,
		OSType:    task.OSType,
		StartTime: time.Now(),
	}

	connConfig := &HostConnConfig{
		Host:       task.Host,
		Port:       task.Port,
		Username:   task.Username,
		Password:   task.Password,
		PrivateKey: task.Key,
		OSType:     task.OSType,
		Transport:  task.Transport,
		UseHTTPS:   task.WinRMUseHttps,
		Timeout:    30 * time.Second,
	}

	fmt.Printf(">>> RunBaselineCheck: connecting to %s:%d ...\n", task.Host, task.Port)
	conn, err := c.connManager.GetConnection(ctx, connConfig)
	if err != nil {
		fmt.Printf(">>> RunBaselineCheck: CONNECTION FAILED: %v\n", err)
		return nil, fmt.Errorf("connect to host %s failed: %v", task.Host, err)
	}
	fmt.Printf(">>> RunBaselineCheck: connected OK\n")

	rules := c.ruleEngine.GetRulesByOSType(task.OSType)
	if len(rules) == 0 {
		fmt.Printf(">>> RunBaselineCheck: no rules for osType=%d, trying fallback\n", task.OSType)
		rules = c.ruleEngine.GetAllRules()
		var filtered []BaselineRule
		for _, rule := range rules {
			if rule.OSType == 0 || rule.OSType == task.OSType {
				filtered = append(filtered, rule)
			}
		}
		rules = filtered
	}
	fmt.Printf(">>> RunBaselineCheck: total rules loaded: %d\n", len(rules))

	report.TotalRules = len(rules)

	var cleanModel mysqls.BaselineCheckResult
	if err := cleanModel.DeleteByTaskTargetAndScene(ctx, task.TaskID, task.Host, scene); err != nil {
		fmt.Printf(">>> RunBaselineCheck: DeleteByTaskTargetAndScene error: %v\n", err)
		return report, fmt.Errorf("clean old results failed: %v", err)
	}

	totalWeight := 0.0
	deduction := 0.0

	emitProgress := func(result mysqls.BaselineCheckResult) {
		if task.RuleProgress == nil {
			return
		}
		task.RuleProgress(typespec.BatchRuleResult{
			RuleID:        result.RuleID,
			RuleName:      result.RuleName,
			CheckResult:   result.CheckResult,
			ResultName:    enums.BaselineEnum.GetCheckResultName(result.CheckResult),
			ExpectedValue: result.ExpectedValue,
			ActualValue:   summarizeBaselineLogValue(result.ActualValue),
			Time:          result.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	for idx, rule := range rules {
		fmt.Printf(">>> RunBaselineCheck: checking rule[%d] id=%d name=%q category=%d risk=%d osType=%d\n",
			idx, rule.ID, rule.Name, rule.Category, rule.Risk, rule.OSType)
		select {
		case <-ctx.Done():
			fmt.Printf(">>> RunBaselineCheck: context cancelled\n")
			return report, ctx.Err()
		default:
		}

		result := mysqls.BaselineCheckResult{
			TaskID:          task.TaskID,
			TargetID:        task.TargetID,
			TargetIP:        task.Host,
			OSType:          task.OSType,
			ScanScene:       scene,
			RuleID:          rule.ID,
			RuleName:        rule.Name,
			RuleCategory:    rule.Category,
			RuleRisk:        rule.Risk,
			ExpectedValue:   rule.ExpectedValue,
			CheckCommand:    strings.Join(rule.Commands, "; "),
			FixSuggestion:   rule.FixSuggestion,
			RiskDescription: rule.RiskDescription,
		}

		ruleWeight := getRiskWeight(rule.Risk)
		totalWeight += ruleWeight

		if ruleHasPlaceholder(rule) {
			fmt.Printf(">>> RunBaselineCheck: rule[%d] has placeholder, skipping\n", idx)
			result.ActualValue = "规则含未替换占位符，当前环境不适配"
			result.CheckResult = enums.BaselineCheckResultSkip
			result.CreateTime = time.Now()
			report.Results = append(report.Results, result)
			if err := result.Add(ctx); err != nil {
				fmt.Printf(">>> RunBaselineCheck: rule[%d] INSERT error (skip): %v\n", idx, err)
			}
			emitProgress(result)
			continue
		}

		fmt.Printf(">>> RunBaselineCheck: rule[%d] has %d command(s), first cmd: %q\n", idx, len(rule.Commands),
			func(cmds []string) string { if len(cmds) > 0 { return cmds[0] }; return "" }(rule.Commands))

		allOutput := ""
		execFailed := false
		for ci, cmd := range rule.Commands {
			cmdToRun := normalizeBaselineCommand(cmd)
			if cmdToRun == "" {
				fmt.Printf(">>> RunBaselineCheck: rule[%d] cmd[%d] empty after normalize, skip\n", idx, ci)
				continue
			}
			fmt.Printf(">>> RunBaselineCheck: rule[%d] executing cmd[%d]: %s\n", idx, ci, cmdToRun[:min(len(cmdToRun), 80)])
			rawOutput, err := c.connManager.ExecuteCommand(ctx, conn, cmdToRun)
			output := strings.TrimSpace(rawOutput)
			if err != nil && output == "" {
				fmt.Printf(">>> RunBaselineCheck: rule[%d] cmd[%d] error+no output: %v\n", idx, ci, err)
				allOutput += "\n"
				continue
			}
			allOutput += output + "\n"
			if isBaselineExecutionError(output) {
				fmt.Printf(">>> RunBaselineCheck: rule[%d] cmd[%d] execution error output\n", idx, ci)
				execFailed = true
				break
			}
		}

		result.ActualValue = truncateBaselineField(strings.TrimSpace(allOutput))
		fmt.Printf(">>> RunBaselineCheck: rule[%d] actualValue len=%d, execFailed=%v\n", idx, len(result.ActualValue), execFailed)

		if execFailed || isBaselineExecutionError(result.ActualValue) {
			fmt.Printf(">>> RunBaselineCheck: rule[%d] -> ERROR\n", idx)
			result.CheckResult = enums.BaselineCheckResultError
			report.ErrorCount++
			deduction += ruleWeight * 0.5
		} else if c.ruleEngine.CheckCommandOutput(result.ActualValue, rule.ExpectedValue, rule.MatchType) {
			fmt.Printf(">>> RunBaselineCheck: rule[%d] -> PASS\n", idx)
			result.CheckResult = enums.BaselineCheckResultPass
			report.PassCount++
		} else {
			fmt.Printf(">>> RunBaselineCheck: rule[%d] -> FAIL (expected=%q matchType=%s)\n", idx, rule.ExpectedValue, rule.MatchType)
			result.CheckResult = enums.BaselineCheckResultFail
			report.FailCount++
			deduction += ruleWeight
			c.updateRiskCount(report, rule.Risk)
		}

		result.CreateTime = time.Now()
		report.Results = append(report.Results, result)
		if err := result.Add(ctx); err != nil {
			fmt.Printf(">>> RunBaselineCheck: rule[%d] INSERT error: %v\n", idx, err)
		}
		emitProgress(result)
	}

	if totalWeight > 0 {
		report.ComplianceScore = ((totalWeight - deduction) / totalWeight) * 100
	} else {
		report.ComplianceScore = 100.0
	}

	report.EndTime = time.Now()
	fmt.Printf(">>> RunBaselineCheck: done. totalRules=%d pass=%d fail=%d err=%d skip=%d score=%.1f\n",
		report.TotalRules, report.PassCount, report.FailCount, report.ErrorCount,
		len(report.Results)-report.PassCount-report.FailCount-report.ErrorCount, report.ComplianceScore)

	return report, nil
}

func summarizeBaselineLogValue(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
	s = strings.ReplaceAll(s, "\n", " | ")
	if len(s) <= 240 {
		return s
	}
	return s[:237] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ruleHasPlaceholder(rule BaselineRule) bool {
	return RuleHasUnresolvedPlaceholder(rule.Commands, rule.ExpectedValue)
}

func getRiskWeight(risk int) float64 {
	switch risk {
	case enums.BaselineRiskCritical:
		return 5.0
	case enums.BaselineRiskHigh:
		return 3.0
	case enums.BaselineRiskMiddle:
		return 2.0
	case enums.BaselineRiskLow:
		return 1.0
	case enums.BaselineRiskInfo:
		return 0.5
	default:
		return 1.0
	}
}

func (c *HostBaselineChecker) updateRiskCount(report *BaselineCheckReport, risk int) {
	switch risk {
	case enums.BaselineRiskCritical:
		report.CriticalCount++
	case enums.BaselineRiskHigh:
		report.HighRiskCount++
	case enums.BaselineRiskMiddle:
		report.MiddleRiskCount++
	case enums.BaselineRiskLow:
		report.LowRiskCount++
	}
}


