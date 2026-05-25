package services

import (
	"context"
	"fmt"
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
	Transport      int  // 0=auto（由连接层解析），1=SSH，2=WinRM
	WinRMUseHttps  bool
	ScanScene      int // 1=安全配置核查 2=主机漏洞检测
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

	conn, err := c.connManager.GetConnection(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to host %s failed: %v", task.Host, err)
	}

	var cleanModel mysqls.BaselineCheckResult
	if err := cleanModel.DeleteByTaskIDAndTargetIP(ctx, task.TaskID, task.Host); err != nil {
		return nil, fmt.Errorf("clean old results failed: %v", err)
	}

	rules := c.ruleEngine.GetRulesByOSType(task.OSType)
	if len(rules) == 0 {
		rules = c.ruleEngine.GetAllRules()
		var filtered []BaselineRule
		for _, rule := range rules {
			if rule.OSType == 0 || rule.OSType == task.OSType {
				filtered = append(filtered, rule)
			}
		}
		rules = filtered
	}

	report.TotalRules = len(rules)

	totalWeight := 0.0
	deduction := 0.0

	for _, rule := range rules {
		select {
		case <-ctx.Done():
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
			result.ActualValue = "规则含未替换占位符，当前环境不适配"
			result.CheckResult = enums.BaselineCheckResultSkip
			result.CreateTime = time.Now()
			report.Results = append(report.Results, result)
			continue
		}

		allOutput := ""
		execFailed := false
		for _, cmd := range rule.Commands {
			cmdToRun := normalizeBaselineCommand(cmd)
			if cmdToRun == "" {
				continue
			}
			rawOutput, err := c.connManager.ExecuteCommand(ctx, conn, cmdToRun)
			output := strings.TrimSpace(rawOutput)
			if err != nil && output == "" {
				allOutput += "\n"
				continue
			}
			allOutput += output + "\n"
			if isBaselineExecutionError(output) {
				execFailed = true
				break
			}
		}

		result.ActualValue = truncateBaselineField(strings.TrimSpace(allOutput))
		if execFailed || isBaselineExecutionError(result.ActualValue) {
			result.CheckResult = enums.BaselineCheckResultError
			report.ErrorCount++
			deduction += ruleWeight * 0.5
		} else if c.ruleEngine.CheckCommandOutput(result.ActualValue, rule.ExpectedValue, rule.MatchType) {
			result.CheckResult = enums.BaselineCheckResultPass
			report.PassCount++
		} else {
			result.CheckResult = enums.BaselineCheckResultFail
			report.FailCount++
			deduction += ruleWeight
			c.updateRiskCount(report, rule.Risk)
		}

		result.CreateTime = time.Now()
		report.Results = append(report.Results, result)
	}

	if totalWeight > 0 {
		report.ComplianceScore = ((totalWeight - deduction) / totalWeight) * 100
	} else {
		report.ComplianceScore = 100.0
	}

	report.EndTime = time.Now()

	if err := cleanModel.BatchAdd(ctx, report.Results); err != nil {
		return report, fmt.Errorf("save results failed: %v", err)
	}

	return report, nil
}

func ruleHasPlaceholder(rule BaselineRule) bool {
	if hasUnresolvedPlaceholder(rule.ExpectedValue) {
		return true
	}
	for _, cmd := range rule.Commands {
		if hasUnresolvedPlaceholder(cmd) {
			return true
		}
	}
	return false
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
