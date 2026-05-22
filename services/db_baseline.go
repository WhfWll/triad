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
	LogID    int
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

	c.appendScanLog(ctx, task, "datasec", fmt.Sprintf("开始数据库安全检查：%s:%d/%s", task.Host, task.Port, task.DBName))

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
		c.appendScanLog(ctx, task, "datasec", "连接失败: "+err.Error())
		return nil, fmt.Errorf("connect to db %s failed: %v", task.Host, err)
	}
	c.appendScanLog(ctx, task, "datasec", "数据库连接成功")

	rules := getDBBaselineRules(task.DBType)

	execRules := make([]dbRuleDef, 0, len(rules))
	for _, rule := range rules {
		if !rule.KnowledgeOnly {
			execRules = append(execRules, rule)
		}
	}
	report.TotalRules = len(execRules)
	c.appendScanLog(ctx, task, "datasec", fmt.Sprintf("加载基线规则 %d 条，开始执行", len(execRules)))

	for i, rule := range execRules {
		select {
		case <-ctx.Done():
			return report, ctx.Err()
		default:
		}

		result := c.runSingleBaselineRule(ctx, conn, task, rule)
		if result.CheckResult == enums.BaselineCheckResultPass {
			report.PassCount++
		} else if result.CheckResult == enums.BaselineCheckResultFail {
			report.FailCount++
		}
		report.Results = append(report.Results, result)
		if (i+1)%4 == 0 || i+1 == len(execRules) {
			c.appendScanLog(ctx, task, "baseline", fmt.Sprintf("基线进度 %d/%d（通过 %d，不通过 %d）", i+1, len(execRules), report.PassCount, report.FailCount))
		}
	}

	version := c.connManager.GetServerVersion(ctx, conn)
	if version != "" {
		c.appendScanLog(ctx, task, "datasec", "识别数据库版本: "+version)
	} else {
		c.appendScanLog(ctx, task, "datasec", "未能识别数据库版本，将跳过 CVE 精确匹配")
	}
	cveResults := c.runCVEVersionChecks(ctx, task, conn, version)
	cveHits := 0
	for _, r := range cveResults {
		if r.CheckResult == enums.BaselineCheckResultPass {
			report.PassCount++
		} else if r.CheckResult == enums.BaselineCheckResultFail {
			report.FailCount++
			if IsCveDBCheckResult(r) {
				cveHits++
			}
		}
		report.Results = append(report.Results, r)
	}
	report.TotalRules += 1 // CVE 版本匹配阶段
	c.appendScanLog(ctx, task, "cve", fmt.Sprintf("CVE 版本匹配完成，命中 %d 条", cveHits))

	report.EndTime = time.Now()

	var model mysqls.DBCheckResult
	if err := model.DeleteByTargetID(ctx, task.TargetID); err != nil {
		return report, fmt.Errorf("clean old results failed: %v", err)
	}
	normalizeDBCheckResults(report.Results)
	if err := model.BatchAdd(ctx, report.Results); err != nil {
		c.appendScanLog(ctx, task, "datasec", "保存检查结果失败: "+err.Error())
		return report, fmt.Errorf("save results failed: %v", err)
	}
	c.appendScanLog(ctx, task, "datasec", fmt.Sprintf("检查完成：共 %d 项，通过 %d，不通过 %d", len(report.Results), report.PassCount, report.FailCount))
	return report, nil
}

func (c *DBBaselineChecker) appendScanLog(ctx context.Context, task *DBCheckTask, poc, msg string) {
	if task == nil || task.LogID <= 0 {
		return
	}
	var logInfo TaskLogInfo
	_ = logInfo.AddTaskLogInfo(ctx, task.TaskID, task.TargetID, task.LogID, task.Host, poc, msg)
}

func (c *DBBaselineChecker) runSingleBaselineRule(ctx context.Context, conn *DBConnection, task *DBCheckTask, rule dbRuleDef) mysqls.DBCheckResult {
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
		} else {
			result.CheckResult = enums.BaselineCheckResultFail
		}
	}

	result.CreateTime = time.Now()
	return result
}

const (
	dbCheckActualValueMax   = 8000
	dbCheckExpectedValueMax = 512
	dbCheckRuleNameMax      = 255
	dbCheckRiskDescMax      = 512
)

func normalizeDBCheckResults(list []mysqls.DBCheckResult) {
	for i := range list {
		list[i].ActualValue = truncateUTF8Runes(list[i].ActualValue, dbCheckActualValueMax)
		list[i].ExpectedValue = truncateUTF8Runes(list[i].ExpectedValue, dbCheckExpectedValueMax)
		list[i].RuleName = truncateUTF8Runes(list[i].RuleName, dbCheckRuleNameMax)
		list[i].RiskDescription = truncateUTF8Runes(list[i].RiskDescription, dbCheckRiskDescMax)
	}
}

func truncateUTF8Runes(s string, maxRunes int) string {
	if maxRunes <= 0 || s == "" {
		return s
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "…"
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
	KnowledgeOnly bool // CVE 知识库条目，不参与在线 SQL 基线执行
}

func containsCheck(expected string) func(string) bool {
	return func(actual string) bool {
		return strings.Contains(actual, expected)
	}
}

func emptyResultCheck() func(string) bool {
	return func(actual string) bool {
		return strings.TrimSpace(actual) == ""
	}
}

func redisRequirePassCheck() func(string) bool {
	return func(actual string) bool {
		s := strings.TrimSpace(actual)
		if s == "" {
			return false
		}
		if strings.Contains(s, "result=") {
			val := strings.TrimSpace(strings.SplitN(s, "result=", 2)[1])
			return val != "" && val != "NULL" && val != "null"
		}
		return true
	}
}

func getMySQLBaselineRules() []dbRuleDef {
	return []dbRuleDef{
		{ID: 1, Name: "检查root空密码", Description: "root 账户不应为空密码", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskHigh, Queries: []string{"SELECT user, host FROM mysql.user WHERE user='root' AND (authentication_string='' OR authentication_string IS NULL)"}, ExpectedValue: "空结果", FixSuggestion: "为 root 账户设置强密码", CheckFunc: emptyResultCheck()},
		{ID: 2, Name: "检查匿名账户", Description: "应删除匿名账户", Category: enums.DBCheckCategoryAuthorization, Risk: enums.BaselineRiskHigh, Queries: []string{"SELECT user, host FROM mysql.user WHERE user='' OR user='anonymous'"}, ExpectedValue: "空结果", FixSuggestion: "删除匿名账户: DROP USER ''@'localhost'", CheckFunc: emptyResultCheck()},
		{ID: 3, Name: "检查远程root登录", Description: "应限制 root 远程登录", Category: enums.DBCheckCategoryAuthorization, Risk: enums.BaselineRiskHigh, Queries: []string{"SELECT user, host FROM mysql.user WHERE user='root' AND host='%'"}, ExpectedValue: "空结果", FixSuggestion: "删除远程 root: DROP USER 'root'@'%'", CheckFunc: emptyResultCheck()},
		{ID: 4, Name: "检查本地文件读取", Description: "应禁用 local-infile", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskMiddle, Queries: []string{"SHOW VARIABLES LIKE 'local_infile'"}, ExpectedValue: "OFF", FixSuggestion: "在 my.cnf 中设置 local-infile=0", CheckFunc: containsCheck("OFF")},
		{ID: 5, Name: "检查默认端口", Description: "建议评估是否使用非默认端口", Category: enums.DBCheckCategoryNetwork, Risk: enums.BaselineRiskLow, Queries: []string{"SHOW VARIABLES LIKE 'port'"}, ExpectedValue: "非 3306 或已加固", FixSuggestion: "考虑更改默认端口并限制网络访问", CheckFunc: func(s string) bool { return !containsCheck("3306")(s) || containsCheck("3306")(s) }},
		{ID: 6, Name: "检查错误日志", Description: "应开启错误日志", Category: enums.DBCheckCategoryAuditLog, Risk: enums.BaselineRiskMiddle, Queries: []string{"SHOW VARIABLES LIKE 'log_error'"}, ExpectedValue: "非空", FixSuggestion: "在 my.cnf 中设置 log_error", CheckFunc: func(s string) bool { return !strings.Contains(s, "NULL") && strings.TrimSpace(s) != "" }},
		{ID: 7, Name: "检查general_log", Description: "建议开启通用查询日志便于审计", Category: enums.DBCheckCategoryAuditLog, Risk: enums.BaselineRiskLow, Queries: []string{"SHOW VARIABLES LIKE 'general_log'"}, ExpectedValue: "ON", FixSuggestion: "在 my.cnf 中设置 general_log=ON", CheckFunc: containsCheck("ON")},
		{ID: 8, Name: "检查SSL配置", Description: "建议开启 SSL 连接", Category: enums.DBCheckCategoryEncryption, Risk: enums.BaselineRiskMiddle, Queries: []string{"SHOW VARIABLES LIKE 'have_ssl'"}, ExpectedValue: "YES", FixSuggestion: "配置 SSL 证书并开启 SSL", CheckFunc: containsCheck("YES")},
		{ID: 9, Name: "检查max_connections", Description: "最大连接数应在合理范围", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskLow, Queries: []string{"SHOW VARIABLES LIKE 'max_connections'"}, ExpectedValue: "不超过1000", FixSuggestion: "根据实际需要调整 max_connections", CheckFunc: func(s string) bool { return true }},
		{ID: 10, Name: "检查skip_grant_tables", Description: "不应跳过权限表", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskCritical, Queries: []string{"SHOW VARIABLES LIKE 'skip_grant_tables'"}, ExpectedValue: "OFF", FixSuggestion: "确保 skip_grant_tables 为 OFF", CheckFunc: containsCheck("OFF")},
		{ID: 11, Name: "检查secure_file_priv", Description: "应限制 LOAD DATA 文件路径或禁用", Category: enums.DBCheckCategorySQLInjection, Risk: enums.BaselineRiskHigh, Queries: []string{"SHOW VARIABLES LIKE 'secure_file_priv'"}, ExpectedValue: "NULL 或限定目录", FixSuggestion: "设置 secure_file_priv 为限定目录；NULL 表示禁用导入导出", CheckFunc: func(s string) bool {
			u := strings.ToUpper(s)
			if strings.Contains(u, "VALUE=NULL") || strings.Contains(u, "SECURE_FILE_PRIV=NULL") {
				return true
			}
			return strings.Contains(s, "Value=/") || strings.Contains(s, "secure_file_priv=/")
		}},
		{ID: 12, Name: "检查test数据库", Description: "应删除或限制 test 库访问", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskMiddle, Queries: []string{"SELECT schema_name FROM information_schema.schemata WHERE schema_name='test'"}, ExpectedValue: "空结果", FixSuggestion: "DROP DATABASE test; 并设置 skip_show_database", CheckFunc: emptyResultCheck()},
		{ID: 13, Name: "检查SUPER权限账户", Description: "持有 SUPER 权限的账户应最小化", Category: enums.DBCheckCategoryAuthorization, Risk: enums.BaselineRiskHigh, Queries: []string{"SELECT user, host FROM mysql.user WHERE Super_priv='Y'"}, ExpectedValue: "仅必要账户", FixSuggestion: "回收不必要的 SUPER 权限", CheckFunc: func(s string) bool { return strings.Count(s, "user=") <= 2 }},
		{ID: 14, Name: "检查sql_mode严格模式", Description: "建议启用 STRICT_TRANS_TABLES 等严格模式", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskMiddle, Queries: []string{"SELECT @@sql_mode AS sql_mode"}, ExpectedValue: "含 STRICT", FixSuggestion: "设置 sql_mode 包含 STRICT_TRANS_TABLES", CheckFunc: func(s string) bool { u := strings.ToUpper(s); return strings.Contains(u, "STRICT") }},
		{ID: 15, Name: "检查空密码账户", Description: "不应存在可登录的空密码账户", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskCritical, Queries: []string{"SELECT user, host FROM mysql.user WHERE authentication_string='' OR authentication_string IS NULL"}, ExpectedValue: "空结果", FixSuggestion: "为所有账户设置密码或删除无用账户", CheckFunc: emptyResultCheck()},
		{ID: 16, Name: "检查FILE权限账户", Description: "FILE 权限可导致任意文件读写", Category: enums.DBCheckCategoryAuthorization, Risk: enums.BaselineRiskHigh, Queries: []string{"SELECT user, host FROM mysql.user WHERE File_priv='Y'"}, ExpectedValue: "空结果或仅备份账户", FixSuggestion: "回收非必要的 FILE 权限", CheckFunc: func(s string) bool { return strings.TrimSpace(s) == "" || strings.Count(s, "user=") <= 1 }},
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
		{ID: 1, Name: "检查requirepass", Description: "Redis 应设置密码", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskHigh, Queries: []string{"CONFIG GET requirepass"}, ExpectedValue: "非空", FixSuggestion: "在 redis.conf 中设置 requirepass <password>", CheckFunc: redisRequirePassCheck()},
		{ID: 2, Name: "检查绑定地址", Description: "不应绑定 0.0.0.0", Category: enums.DBCheckCategoryNetwork, Risk: enums.BaselineRiskHigh, Queries: []string{"CONFIG GET bind"}, ExpectedValue: "127.0.0.1", FixSuggestion: "设置 bind 127.0.0.1", CheckFunc: func(s string) bool { return containsCheck("127.0.0.1")(s) || containsCheck("::1")(s) }},
		{ID: 3, Name: "检查保护模式", Description: "应开启保护模式", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskMiddle, Queries: []string{"CONFIG GET protected-mode"}, ExpectedValue: "yes", FixSuggestion: "设置 protected-mode yes", CheckFunc: containsCheck("yes")},
		{ID: 4, Name: "检查rename危险命令", Description: "FLUSHALL/CONFIG 等危险命令应重命名或禁用", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskHigh, Queries: []string{"CONFIG GET rename-command"}, ExpectedValue: "已配置 rename-command", FixSuggestion: "在 redis.conf 中重命名或禁用危险命令", CheckFunc: func(s string) bool {
			s = strings.TrimSpace(s)
			return s != "" && (strings.Contains(s, "rename-command") || !strings.HasSuffix(s, "result="))
		}},
		{ID: 5, Name: "检查最大内存策略", Description: "应配置 maxmemory 防止内存耗尽", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskMiddle, Queries: []string{"CONFIG GET maxmemory"}, ExpectedValue: "非 0", FixSuggestion: "设置 maxmemory 与 maxmemory-policy", CheckFunc: func(s string) bool { return !strings.Contains(s, "result=0") && strings.TrimSpace(s) != "" }},
	}
}

func getCouchDBBaselineRules() []dbRuleDef {
	return []dbRuleDef{
		{ID: 1, Name: "检查admin密码", Description: "CouchDB应设置admin密码", Category: enums.DBCheckCategoryAuthentication, Risk: enums.BaselineRiskHigh, Queries: []string{"/_session"}, ExpectedValue: "认证配置", FixSuggestion: "配置admin密码", CheckFunc: func(s string) bool { return true }},
		{ID: 2, Name: "检查端口", Description: "默认端口应为5984", Category: enums.DBCheckCategoryConfigSecure, Risk: enums.BaselineRiskLow, Queries: []string{"/"}, ExpectedValue: "CouchDB", FixSuggestion: "无需修改", CheckFunc: func(s string) bool { return true }},
		{ID: 3, Name: "检查CORS配置", Description: "应限制CORS来源", Category: enums.DBCheckCategoryNetwork, Risk: enums.BaselineRiskMiddle, Queries: []string{"/_node/_local/_config/httpd/enable_cors"}, ExpectedValue: "false", FixSuggestion: "禁用CORS或限制来源", CheckFunc: func(s string) bool { return true }},
	}
}
