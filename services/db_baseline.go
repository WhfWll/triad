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

type DBBaselineChecker struct {
	connManager *DBConnManager
}

var (
	globalDBBaselineChecker *DBBaselineChecker
	dbBaselineOnce          sync.Once
)

func GetDBBaselineChecker() *DBBaselineChecker {
	dbBaselineOnce.Do(func() {
		globalDBBaselineChecker = &DBBaselineChecker{
			connManager: GetDBConnManager(),
		}
	})
	return globalDBBaselineChecker
}

type DBCheckTask struct {
	TaskID   int
	TargetID int
	Host     string
	Port     int
	DBType   int
	Username string
	Password string
	DBName   string
}

type DBCheckReport struct {
	TaskID     int
	TargetID   int
	TargetIP   string
	DBType     int
	TotalRules int
	PassCount  int
	FailCount  int
	Results    []mysqls.DBCheckResult
	StartTime  time.Time
	EndTime    time.Time
}

func (c *DBBaselineChecker) RunDBCheck(ctx context.Context, task *DBCheckTask) (*DBCheckReport, error) {
	report := &DBCheckReport{
		TaskID:    task.TaskID,
		TargetID:  task.TargetID,
		TargetIP:  task.Host,
		DBType:    task.DBType,
		StartTime: time.Now(),
	}

	config := &DBConnConfig{
		DBType:   task.DBType,
		Host:     task.Host,
		Port:     task.Port,
		Username: task.Username,
		Password: task.Password,
		DBName:   task.DBName,
		Timeout:  30 * time.Second,
	}

	conn, err := c.connManager.GetConnection(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connect to db %s failed: %v", task.Host, err)
	}

	rules := getDBBaselineRules(task.DBType)
	report.TotalRules = len(rules)

	for _, rule := range rules {
		select {
		case <-ctx.Done():
			return report, ctx.Err()
		default:
		}

		result := mysqls.DBCheckResult{
			TaskID:          task.TaskID,
			TargetID:        task.TargetID,
			TargetIP:        task.Host,
			DBType:          task.DBType,
			DBName:          task.DBName,
			CheckCategory:   rule.Category,
			RuleID:          rule.ID,
			RuleName:        rule.Name,
			ExpectedValue:   rule.ExpectedValue,
			RiskLevel:       rule.Risk,
			FixSuggestion:   rule.FixSuggestion,
			RiskDescription: rule.Description,
		}

		allOutput := ""
		checkErr := false
		for _, query := range rule.Queries {
			results, err := c.connManager.ExecuteQuery(ctx, conn, query)
			if err != nil {
				result.ActualValue = fmt.Sprintf("ERROR: %v", err)
				result.CheckResult = enums.BaselineCheckResultError
				checkErr = true
				break
			}
			for _, row := range results {
				for k, v := range row {
					allOutput += fmt.Sprintf("%s=%s ", k, v)
				}
			}
			allOutput += "\n"
		}

		if !checkErr {
			result.ActualValue = strings.TrimSpace(allOutput)
			if rule.CheckFunc(result.ActualValue) {
				result.CheckResult = enums.BaselineCheckResultPass
				report.PassCount++
			} else {
				result.CheckResult = enums.BaselineCheckResultFail
				report.FailCount++
			}
		}

		result.CreateTime = time.Now()
		report.Results = append(report.Results, result)
	}

	report.EndTime = time.Now()

	var model mysqls.DBCheckResult
	if err := model.DeleteByTaskID(ctx, task.TaskID); err != nil {
		return report, fmt.Errorf("clean old results failed: %v", err)
	}
	if err := model.BatchAdd(ctx, report.Results); err != nil {
		return report, fmt.Errorf("save results failed: %v", err)
	}

	return report, nil
}

type dbRuleDef struct {
	ID            int
	Name          string
	Description   string
	Category      int
	Risk          int
	Queries       []string
	ExpectedValue string
	FixSuggestion string
	CheckFunc     func(string) bool
}

func containsCheck(expected string) func(string) bool {
	return func(actual string) bool {
		return strings.Contains(actual, expected)
	}
}

func getDBBaselineRules(dbType int) []dbRuleDef {
	switch dbType {
	case enums.DBSupportTypeMySQL:
		return getMySQLBaselineRules()
	case enums.DBSupportTypePostgreSQL:
		return getPostgreSQLBaselineRules()
	case enums.DBSupportTypeMongoDB:
		return getMongoDBBaselineRules()
	case enums.DBSupportTypeRedis:
		return getRedisBaselineRules()
	case enums.DBSupportTypeCouchDB:
		return getCouchDBBaselineRules()
	}
	return nil
}

func getMySQLBaselineRules() []dbRuleDef {
	return []dbRuleDef{
		{ID: 1, Name: "检查root密码强度", Description: "root账户应设置强密码", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskHigh, Queries: []string{"SELECT user, authentication_string FROM mysql.user WHERE user='root' AND authentication_string=''"}, ExpectedValue: "空结果", FixSuggestion: "为root账户设置强密码", CheckFunc: func(s string) bool { return strings.Contains(s, "NULL") || s == "" }},
		{ID: 2, Name: "检查匿名账户", Description: "应删除匿名账户", Category: enums.DBCheckCategoryAuthorization, Risk: enums.BaselineRiskHigh, Queries: []string{"SELECT user, host FROM mysql.user WHERE user='' OR user='anonymous'"}, ExpectedValue: "空结果", FixSuggestion: "删除匿名账户: DROP USER ''@'localhost'", CheckFunc: func(s string) bool { return s == "" }},
		{ID: 3, Name: "检查远程root登录", Description: "应限制root远程登录", Category: enums.DBCheckCategoryAuthorization, Risk: enums.BaselineRiskHigh, Queries: []string{"SELECT user, host FROM mysql.user WHERE user='root' AND host='%'"}, ExpectedValue: "空结果", FixSuggestion: "删除远程root: DROP USER 'root'@'%'", CheckFunc: func(s string) bool { return s == "" }},
		{ID: 4, Name: "检查本地文件读取", Description: "应禁用local-infile", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskMiddle, Queries: []string{"SHOW VARIABLES LIKE 'local_infile'"}, ExpectedValue: "OFF", FixSuggestion: "在my.cnf中设置 local-infile=0", CheckFunc: containsCheck("OFF")},
		{ID: 5, Name: "检查端口", Description: "应使用非默认端口或限制访问", Category: enums.DBCheckCategoryNetwork, Risk: enums.BaselineRiskLow, Queries: []string{"SHOW VARIABLES LIKE 'port'"}, ExpectedValue: "3306", FixSuggestion: "考虑更改默认端口", CheckFunc: func(s string) bool { return true }},
		{ID: 6, Name: "检查log配置", Description: "应开启错误日志", Category: enums.DBCheckCategoryAuditLog, Risk: enums.BaselineRiskMiddle, Queries: []string{"SHOW VARIABLES LIKE 'log_error'"}, ExpectedValue: "非空", FixSuggestion: "在my.cnf中设置 log_error", CheckFunc: func(s string) bool { return !strings.Contains(s, "NULL") && s != "" }},
		{ID: 7, Name: "检查general_log", Description: "建议开启通用查询日志便于审计", Category: enums.DBCheckCategoryAuditLog, Risk: enums.BaselineRiskLow, Queries: []string{"SHOW VARIABLES LIKE 'general_log'"}, ExpectedValue: "ON", FixSuggestion: "在my.cnf中设置 general_log=ON", CheckFunc: containsCheck("ON")},
		{ID: 8, Name: "检查SSL配置", Description: "建议开启SSL连接", Category: enums.DBCheckCategoryEncryption, Risk: enums.BaselineRiskMiddle, Queries: []string{"SHOW VARIABLES LIKE 'have_ssl'"}, ExpectedValue: "YES", FixSuggestion: "配置SSL证书并开启SSL", CheckFunc: containsCheck("YES")},
		{ID: 9, Name: "检查max_connections", Description: "最大连接数应在合理范围", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskLow, Queries: []string{"SHOW VARIABLES LIKE 'max_connections'"}, ExpectedValue: "不超过1000", FixSuggestion: "根据实际需要调整max_connections", CheckFunc: func(s string) bool { return true }},
		{ID: 10, Name: "检查skip_grant_tables", Description: "不应跳过权限表", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskCritical, Queries: []string{"SHOW VARIABLES LIKE 'skip_grant_tables'"}, ExpectedValue: "OFF", FixSuggestion: "确保skip_grant_tables为OFF", CheckFunc: containsCheck("OFF")},
	}
}

func getPostgreSQLBaselineRules() []dbRuleDef {
	return []dbRuleDef{
		{ID: 1, Name: "检查密码加密", Description: "密码应使用scram-sha-256加密", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskHigh, Queries: []string{"SHOW password_encryption"}, ExpectedValue: "scram-sha-256", FixSuggestion: "设置 password_encryption = 'scram-sha-256'", CheckFunc: containsCheck("scram")},
		{ID: 2, Name: "检查监听地址", Description: "不应监听所有地址", Category: enums.DBCheckCategoryNetwork, Risk: enums.BaselineRiskHigh, Queries: []string{"SHOW listen_addresses"}, ExpectedValue: "localhost", FixSuggestion: "设置 listen_addresses = 'localhost'", CheckFunc: containsCheck("localhost")},
		{ID: 3, Name: "检查日志记录", Description: "应开启日志记录", Category: enums.DBCheckCategoryAuditLog, Risk: enums.BaselineRiskMiddle, Queries: []string{"SHOW logging_collector"}, ExpectedValue: "on", FixSuggestion: "设置 logging_collector = on", CheckFunc: containsCheck("on")},
		{ID: 4, Name: "检查SSL", Description: "建议开启SSL", Category: enums.DBCheckCategoryEncryption, Risk: enums.BaselineRiskMiddle, Queries: []string{"SHOW ssl"}, ExpectedValue: "on", FixSuggestion: "配置SSL证书", CheckFunc: containsCheck("on")},
		{ID: 5, Name: "检查log_statement", Description: "应记录所有DDL操作", Category: enums.DBCheckCategoryAuditLog, Risk: enums.BaselineRiskLow, Queries: []string{"SHOW log_statement"}, ExpectedValue: "ddl", FixSuggestion: "设置 log_statement = 'ddl'", CheckFunc: func(s string) bool { return !containsCheck("none")(s) }},
	}
}

func getMongoDBBaselineRules() []dbRuleDef {
	return []dbRuleDef{
		{ID: 1, Name: "检查认证启用", Description: "MongoDB应开启认证", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskHigh, Queries: []string{"/_isMaster"}, ExpectedValue: "认证配置", FixSuggestion: "在mongod.conf中设置 security.authorization: enabled", CheckFunc: func(s string) bool { return true }},
		{ID: 2, Name: "检查绑定IP", Description: "不应绑定0.0.0.0", Category: enums.DBCheckCategoryNetwork, Risk: enums.BaselineRiskHigh, Queries: []string{"/_isMaster"}, ExpectedValue: "localhost", FixSuggestion: "设置 bindIp: 127.0.0.1", CheckFunc: func(s string) bool { return true }},
		{ID: 3, Name: "检查端口", Description: "应使用非默认端口", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskLow, Queries: []string{"/_isMaster"}, ExpectedValue: "27017", FixSuggestion: "考虑更改默认端口", CheckFunc: func(s string) bool { return true }},
	}
}

func getRedisBaselineRules() []dbRuleDef {
	return []dbRuleDef{
		{ID: 1, Name: "检查requirepass", Description: "Redis应设置密码", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskHigh, Queries: []string{"CONFIG GET requirepass"}, ExpectedValue: "非空", FixSuggestion: "在redis.conf中设置 requirepass <password>", CheckFunc: func(s string) bool { return !strings.Contains(s, "") && !strings.Contains(s, "NULL") }},
		{ID: 2, Name: "检查绑定地址", Description: "不应绑定0.0.0.0", Category: enums.DBCheckCategoryNetwork, Risk: enums.BaselineRiskHigh, Queries: []string{"CONFIG GET bind"}, ExpectedValue: "127.0.0.1", FixSuggestion: "设置 bind 127.0.0.1", CheckFunc: containsCheck("127.0.0.1")},
		{ID: 3, Name: "检查保护模式", Description: "应开启保护模式", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskMiddle, Queries: []string{"CONFIG GET protected-mode"}, ExpectedValue: "yes", FixSuggestion: "设置 protected-mode yes", CheckFunc: containsCheck("yes")},
		{ID: 4, Name: "检查rename命令", Description: "危险命令应重命名或禁用", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskHigh, Queries: []string{"CONFIG GET rename-command FLUSHALL", "CONFIG GET rename-command FLUSHDB", "CONFIG GET rename-command CONFIG"}, ExpectedValue: "非空", FixSuggestion: "在redis.conf中重命名危险命令", CheckFunc: func(s string) bool { return true }},
	}
}

func getCouchDBBaselineRules() []dbRuleDef {
	return []dbRuleDef{
		{ID: 1, Name: "检查admin密码", Description: "CouchDB应设置admin密码", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskHigh, Queries: []string{"/_session"}, ExpectedValue: "认证配置", FixSuggestion: "配置admin密码", CheckFunc: func(s string) bool { return true }},
		{ID: 2, Name: "检查端口", Description: "默认端口应为5984", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskLow, Queries: []string{"/"}, ExpectedValue: "CouchDB", FixSuggestion: "无需修改", CheckFunc: func(s string) bool { return true }},
		{ID: 3, Name: "检查CORS配置", Description: "应限制CORS来源", Category: enums.DBCheckCategoryNetwork, Risk: enums.BaselineRiskMiddle, Queries: []string{"/_node/_local/_config/httpd/enable_cors"}, ExpectedValue: "false", FixSuggestion: "禁用CORS或限制来源", CheckFunc: func(s string) bool { return true }},
	}
}
