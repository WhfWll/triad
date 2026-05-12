package services

import (
	"context"
	"fmt"
	"regexp"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strings"
	"sync"
	"time"
)

type SensitiveDataFinder struct {
	connManager *DBConnManager
}

var (
	globalSensitiveDataFinder *SensitiveDataFinder
	sensitiveDataOnce         sync.Once
)

type SensitiveRule struct {
	Type       int
	Level      int
	Name       string
	Pattern    *regexp.Regexp
	FieldNames []string
}

func GetSensitiveDataFinder() *SensitiveDataFinder {
	sensitiveDataOnce.Do(func() {
		globalSensitiveDataFinder = &SensitiveDataFinder{
			connManager: GetDBConnManager(),
		}
	})
	return globalSensitiveDataFinder
}

type SensitiveDataTask struct {
	TaskID    int
	TargetID  int
	Host      string
	Port      int
	DBType    int
	Username  string
	Password  string
	DBName    string
	ScanAllDB bool
}

type SensitiveDataReport struct {
	TaskID      int
	TargetID    int
	TargetIP    string
	DBType      int
	HighCount   int
	MiddleCount int
	LowCount    int
	Results     []mysqls.SensitiveDataResult
	StartTime   time.Time
	EndTime     time.Time
}

var sensitiveRules = []SensitiveRule{
	{Type: enums.SensitiveDataTypeIDCard, Level: enums.SensitiveDataLevelHigh, Name: "身份证号", Pattern: regexp.MustCompile(`[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]`),
		FieldNames: []string{"id_card", "idcard", "identity_card", "身份证", "sfz", "身份证号", "identity_no", "id_number"}},
	{Type: enums.SensitiveDataTypeBankCard, Level: enums.SensitiveDataLevelHigh, Name: "银行卡号", Pattern: regexp.MustCompile(`\b(62|60|58|56|55|54|53|52|51|50|49|48|47|46|45|44|43|42|41|40)\d{13,18}\b`),
		FieldNames: []string{"bank_card", "bankcard", "card_no", "card_number", "银行卡", "银行卡号", "bank_account", "account_no"}},
	{Type: enums.SensitiveDataTypePassport, Level: enums.SensitiveDataLevelHigh, Name: "护照号", Pattern: regexp.MustCompile(`[PpEe]\d{8}`),
		FieldNames: []string{"passport", "护照", "passport_no"}},
	{Type: enums.SensitiveDataTypePhone, Level: enums.SensitiveDataLevelMiddle, Name: "手机号", Pattern: regexp.MustCompile(`1[3-9]\d{9}`),
		FieldNames: []string{"phone", "mobile", "tel", "cellphone", "手机号", "手机", "联系电话", "mobile_phone", "telephone"}},
	{Type: enums.SensitiveDataTypeEmail, Level: enums.SensitiveDataLevelMiddle, Name: "邮箱", Pattern: regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		FieldNames: []string{"email", "mail", "邮箱", "e_mail", "mailbox"}},
	{Type: enums.SensitiveDataTypeAddress, Level: enums.SensitiveDataLevelMiddle, Name: "地址", Pattern: nil,
		FieldNames: []string{"address", "addr", "地址", "住址", "家庭住址", "详细地址", "residence", "location"}},
	{Type: enums.SensitiveDataTypeBirthDate, Level: enums.SensitiveDataLevelMiddle, Name: "出生日期", Pattern: regexp.MustCompile(`\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])`),
		FieldNames: []string{"birth", "birthday", "出生日期", "出生年月", "birth_date", "born_date"}},
	{Type: enums.SensitiveDataTypeName, Level: enums.SensitiveDataLevelLow, Name: "姓名", Pattern: nil,
		FieldNames: []string{"name", "姓名", "user_name", "username", "full_name", "real_name", "真实姓名", "nickname"}},
}

func (f *SensitiveDataFinder) RunScan(ctx context.Context, task *SensitiveDataTask) (*SensitiveDataReport, error) {
	report := &SensitiveDataReport{
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
		Timeout:  60 * time.Second,
	}

	conn, err := f.connManager.GetConnection(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connect to db failed: %v", err)
	}

	databases := []string{task.DBName}
	if task.ScanAllDB {
		if dbs, err := f.connManager.GetDatabases(ctx, conn); err == nil {
			databases = dbs
		}
	}

	for _, dbName := range databases {
		select {
		case <-ctx.Done():
			return report, ctx.Err()
		default:
		}

		tables, err := f.connManager.GetTables(ctx, conn, dbName)
		if err != nil {
			continue
		}

		for _, tableName := range tables {
			if strings.HasPrefix(tableName, "_") || strings.HasPrefix(tableName, "pg_") {
				continue
			}

			columns, err := f.connManager.GetColumns(ctx, conn, dbName, tableName)
			if err != nil || columns == nil {
				continue
			}

			for _, col := range columns {
				colName := col["COLUMN_NAME"]
				if colName == "" {
					colName = col["column_name"]
				}
				dataType := col["DATA_TYPE"]
				if dataType == "" {
					dataType = col["data_type"]
				}

				if f.isTextType(dataType) {
					for _, rule := range sensitiveRules {
						if !f.matchFieldName(colName, rule.FieldNames) {
							continue
						}

						result := mysqls.SensitiveDataResult{
							TaskID:       task.TaskID,
							TargetID:     task.TargetID,
							TargetIP:     task.Host,
							DBType:       task.DBType,
							DBName:       dbName,
							TableNameStr: tableName,
							ColumnName:   colName,
							DataType:     rule.Type,
							DataLevel:    rule.Level,
							MatchRule:    rule.Name,
							MatchType:    1,
							CreateTime:   time.Now(),
						}

						sampleQuery := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` IS NOT NULL LIMIT 5",
							colName, dbName, tableName, colName)
						if task.DBType == 2 {
							sampleQuery = fmt.Sprintf("SELECT \"%s\" FROM \"%s\" WHERE \"%s\" IS NOT NULL LIMIT 5",
								colName, tableName, colName)
						}

						sampleResult, err := f.connManager.ExecuteQuery(ctx, conn, sampleQuery)
						if err == nil && len(sampleResult) > 0 {
							for _, row := range sampleResult {
								if val, ok := row[colName]; ok {
									result.SampleData = truncateString(val, 100)
									if rule.Pattern != nil && !rule.Pattern.MatchString(val) {
										result.SampleData = ""
										break
									}
									if result.SampleData != "" {
										break
									}
								}
							}
						}

						if result.SampleData != "" || rule.Pattern == nil {
							if rule.Level == enums.SensitiveDataLevelHigh {
								report.HighCount++
							} else if rule.Level == enums.SensitiveDataLevelMiddle {
								report.MiddleCount++
							} else {
								report.LowCount++
							}
							report.Results = append(report.Results, result)
						}
					}
				}
			}
		}
	}

	report.EndTime = time.Now()

	var model mysqls.SensitiveDataResult
	if err := model.DeleteByTaskID(ctx, task.TaskID); err != nil {
		return report, fmt.Errorf("clean old results failed: %v", err)
	}
	if err := model.BatchAdd(ctx, report.Results); err != nil {
		return report, fmt.Errorf("save results failed: %v", err)
	}

	return report, nil
}

func (f *SensitiveDataFinder) isTextType(dataType string) bool {
	textTypes := []string{"char", "varchar", "text", "longtext", "mediumtext", "tinytext", "json", "blob", "mediumblob", "longblob"}
	dt := strings.ToLower(dataType)
	for _, t := range textTypes {
		if strings.Contains(dt, t) {
			return true
		}
	}
	return false
}

func (f *SensitiveDataFinder) matchFieldName(colName string, patterns []string) bool {
	lower := strings.ToLower(colName)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
