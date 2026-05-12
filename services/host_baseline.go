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
	TaskID   int
	TargetID int
	Host     string
	Port     int
	Username string
	Password string
	Key      string
	OSType   int
}

type BaselineCheckReport struct {
	TaskID      int
	TargetID    int
	TargetIP    string
	OSType      int
	TotalRules  int
	PassCount   int
	FailCount   int
	ErrorCount  int
	Results     []mysqls.BaselineCheckResult
	StartTime   time.Time
	EndTime     time.Time
}

func (c *HostBaselineChecker) RunBaselineCheck(ctx context.Context, task *BaselineCheckTask) (*BaselineCheckReport, error) {
	report := &BaselineCheckReport{
		TaskID:    task.TaskID,
		TargetID:  task.TargetID,
		TargetIP:  task.Host,
		OSType:    task.OSType,
		StartTime: time.Now(),
	}

	connConfig := &HostConnConfig{
		Host:     task.Host,
		Port:     task.Port,
		Username: task.Username,
		Password: task.Password,
		PrivateKey: task.Key,
		OSType:   task.OSType,
		Timeout:  30 * time.Second,
	}

	conn, err := c.connManager.GetConnection(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to host %s failed: %v", task.Host, err)
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
			RuleID:          rule.ID,
			RuleName:        rule.Name,
			RuleCategory:    rule.Category,
			RuleRisk:        rule.Risk,
			ExpectedValue:   rule.ExpectedValue,
			CheckCommand:    strings.Join(rule.Commands, "; "),
			FixSuggestion:   rule.FixSuggestion,
			RiskDescription: rule.RiskDescription,
		}

		allOutput := ""
		checkErr := false
		for _, cmd := range rule.Commands {
			output, err := c.connManager.ExecuteCommand(ctx, conn, cmd)
			if err != nil {
				result.ActualValue = fmt.Sprintf("ERROR: %v", err)
				result.CheckResult = enums.BaselineCheckResultError
				checkErr = true
				break
			}
			allOutput += output + "\n"
		}

		if !checkErr {
			result.ActualValue = strings.TrimSpace(allOutput)
			if c.ruleEngine.CheckCommandOutput(result.ActualValue, rule.ExpectedValue, rule.MatchType) {
				result.CheckResult = enums.BaselineCheckResultPass
				report.PassCount++
			} else {
				result.CheckResult = enums.BaselineCheckResultFail
				report.FailCount++
			}
		} else {
			report.ErrorCount++
		}

		result.CreateTime = time.Now()
		report.Results = append(report.Results, result)
	}

	report.EndTime = time.Now()

	var model mysqls.BaselineCheckResult
	if err := model.DeleteByTaskID(ctx, task.TaskID); err != nil {
		return report, fmt.Errorf("clean old results failed: %v", err)
	}

	if err := model.BatchAdd(ctx, report.Results); err != nil {
		return report, fmt.Errorf("save results failed: %v", err)
	}

	return report, nil
}
