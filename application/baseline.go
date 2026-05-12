package application

import (
	"context"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
)

type BaselineApp struct{}

func (a *BaselineApp) RunBaselineCheck(ctx context.Context, req *typespec.BaselineCheckReq) (*typespec.BaselineCheckResp, error) {
	checker := services.GetHostBaselineChecker()
	task := &services.BaselineCheckTask{
		TaskID:   req.TaskID,
		TargetID: req.TargetID,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Key:      req.Key,
		OSType:   req.OSType,
	}

	report, err := checker.RunBaselineCheck(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := &typespec.BaselineCheckResp{
		TaskID:     report.TaskID,
		TargetIP:   report.TargetIP,
		OSType:     report.OSType,
		OSTypeName: enums.BaselineEnum.GetOSTypeName(report.OSType),
		TotalRules: report.TotalRules,
		PassCount:  report.PassCount,
		FailCount:  report.FailCount,
		ErrorCount: report.ErrorCount,
		StartTime:  report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:    report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.BaselineCheckItem{
			ID:              r.ID,
			RuleID:          r.RuleID,
			RuleName:        r.RuleName,
			RuleCategory:    r.RuleCategory,
			CategoryName:    enums.BaselineEnum.GetCategoryName(r.RuleCategory),
			RiskLevel:       r.RuleRisk,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.RuleRisk),
			CheckResult:     r.CheckResult,
			ResultName:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			ExpectedValue:   r.ExpectedValue,
			ActualValue:     r.ActualValue,
			CheckCommand:    r.CheckCommand,
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			CheckTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (a *BaselineApp) GetBaselineResults(ctx context.Context, req *typespec.BaselineCheckResultListReq) (*typespec.BaselineCheckResultListResp, error) {
	var model mysqls.BaselineCheckResult
	var list []mysqls.BaselineCheckResult
	var err error

	if req.TargetID > 0 {
		list, err = model.GetByTargetID(ctx, req.TargetID)
	} else {
		list, err = model.GetByTaskID(ctx, req.TaskID)
	}
	if err != nil {
		return nil, err
	}

	resp := &typespec.BaselineCheckResultListResp{
		Total: int64(len(list)),
	}
	for _, r := range list {
		resp.List = append(resp.List, typespec.BaselineCheckItem{
			ID:              r.ID,
			RuleID:          r.RuleID,
			RuleName:        r.RuleName,
			RuleCategory:    r.RuleCategory,
			CategoryName:    enums.BaselineEnum.GetCategoryName(r.RuleCategory),
			RiskLevel:       r.RuleRisk,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.RuleRisk),
			CheckResult:     r.CheckResult,
			ResultName:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			ExpectedValue:   r.ExpectedValue,
			ActualValue:     r.ActualValue,
			CheckCommand:    r.CheckCommand,
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			CheckTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) GetBaselineStat(ctx context.Context, req *typespec.BaselineStatReq) (*typespec.BaselineStatResp, error) {
	var model mysqls.BaselineCheckResult
	pass, fail, total, err := model.GetStatByTaskID(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}

	rate := 0.0
	if total > 0 {
		rate = float64(pass) / float64(total) * 100
	}

	return &typespec.BaselineStatResp{
		TotalRules: int(total),
		PassCount:  int(pass),
		FailCount:  int(fail),
		PassRate:   rate,
	}, nil
}

func (a *BaselineApp) GetBaselineRules(ctx context.Context) interface{} {
	engine := services.GetBaselineEngine()
	return engine.GetAllRules()
}

func (a *BaselineApp) RunMalwareScan(ctx context.Context, req *typespec.MalwareScanReq) (*typespec.MalwareScanResp, error) {
	scanner := services.GetMalwareScanner()
	task := &services.MalwareScanTask{
		TaskID:   req.TaskID,
		TargetID: req.TargetID,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Key:      req.Key,
		OSType:   req.OSType,
	}

	report, err := scanner.RunMalwareScan(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := &typespec.MalwareScanResp{
		TaskID:     report.TaskID,
		TargetIP:   report.TargetIP,
		HasMalware: report.HasMalware,
		StartTime:  report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:    report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.MalwareCheckItem{
			ID:            r.ID,
			CheckType:     r.CheckType,
			CheckTypeName: enums.BaselineEnum.GetMalwareCheckTypeName(r.CheckType),
			RiskLevel:     r.RiskLevel,
			RiskName:      enums.BaselineEnum.GetBaselineRiskName(r.RiskLevel),
			MatchRule:     r.MatchRule,
			FilePath:      r.FilePath,
			Description:   r.Description,
			FixSuggestion: r.FixSuggestion,
			CheckTime:     r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (a *BaselineApp) GetMalwareResults(ctx context.Context, req *typespec.MalwareResultListReq) (*typespec.MalwareResultListResp, error) {
	var model mysqls.MalwareCheckResult
	var list []mysqls.MalwareCheckResult
	var err error

	if req.TargetID > 0 {
		list, err = model.GetByTargetID(ctx, req.TargetID)
	} else {
		list, err = model.GetByTaskID(ctx, req.TaskID)
	}
	if err != nil {
		return nil, err
	}

	resp := &typespec.MalwareResultListResp{
		Total: int64(len(list)),
	}
	for _, r := range list {
		resp.List = append(resp.List, typespec.MalwareCheckItem{
			ID:            r.ID,
			CheckType:     r.CheckType,
			CheckTypeName: enums.BaselineEnum.GetMalwareCheckTypeName(r.CheckType),
			RiskLevel:     r.RiskLevel,
			RiskName:      enums.BaselineEnum.GetBaselineRiskName(r.RiskLevel),
			MatchRule:     r.MatchRule,
			FilePath:      r.FilePath,
			Description:   r.Description,
			FixSuggestion: r.FixSuggestion,
			CheckTime:     r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) RunDBCheck(ctx context.Context, req *typespec.DBCheckReq) (*typespec.DBCheckResp, error) {
	checker := services.GetDBBaselineChecker()
	task := &services.DBCheckTask{
		TaskID:   req.TaskID,
		TargetID: req.TargetID,
		Host:     req.Host,
		Port:     req.Port,
		DBType:   req.DBType,
		Username: req.Username,
		Password: req.Password,
		DBName:   req.DBName,
	}

	report, err := checker.RunDBCheck(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := &typespec.DBCheckResp{
		TaskID:     report.TaskID,
		TargetIP:   report.TargetIP,
		DBType:     report.DBType,
		DBTypeName: enums.BaselineEnum.GetDBTypeName(report.DBType),
		TotalRules: report.TotalRules,
		PassCount:  report.PassCount,
		FailCount:  report.FailCount,
		StartTime:  report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:    report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.DBCheckItem{
			ID:              r.ID,
			RuleID:          r.RuleID,
			RuleName:        r.RuleName,
			CheckCategory:   r.CheckCategory,
			CategoryName:    enums.BaselineEnum.GetDBCheckCategoryName(r.CheckCategory),
			CheckResult:     r.CheckResult,
			ResultName:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			ExpectedValue:   r.ExpectedValue,
			ActualValue:     r.ActualValue,
			RiskLevel:       r.RiskLevel,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.RiskLevel),
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			CheckTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (a *BaselineApp) GetDBCheckResults(ctx context.Context, req *typespec.DBCheckResultListReq) (*typespec.DBCheckResultListResp, error) {
	var model mysqls.DBCheckResult
	var list []mysqls.DBCheckResult
	var err error

	if req.TargetID > 0 {
		list, err = model.GetByTargetID(ctx, req.TargetID)
	} else {
		list, err = model.GetByTaskID(ctx, req.TaskID)
	}
	if err != nil {
		return nil, err
	}

	resp := &typespec.DBCheckResultListResp{
		Total: int64(len(list)),
	}
	for _, r := range list {
		resp.List = append(resp.List, typespec.DBCheckItem{
			ID:              r.ID,
			RuleID:          r.RuleID,
			RuleName:        r.RuleName,
			CheckCategory:   r.CheckCategory,
			CategoryName:    enums.BaselineEnum.GetDBCheckCategoryName(r.CheckCategory),
			CheckResult:     r.CheckResult,
			ResultName:      enums.BaselineEnum.GetCheckResultName(r.CheckResult),
			ExpectedValue:   r.ExpectedValue,
			ActualValue:     r.ActualValue,
			RiskLevel:       r.RiskLevel,
			RiskName:        enums.BaselineEnum.GetBaselineRiskName(r.RiskLevel),
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			CheckTime:       r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) RunSensitiveDataScan(ctx context.Context, req *typespec.SensitiveDataScanReq) (*typespec.SensitiveDataScanResp, error) {
	finder := services.GetSensitiveDataFinder()
	task := &services.SensitiveDataTask{
		TaskID:    req.TaskID,
		TargetID:  req.TargetID,
		Host:      req.Host,
		Port:      req.Port,
		DBType:    req.DBType,
		Username:  req.Username,
		Password:  req.Password,
		DBName:    req.DBName,
		ScanAllDB: req.ScanAllDB,
	}

	report, err := finder.RunScan(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := &typespec.SensitiveDataScanResp{
		TaskID:      report.TaskID,
		TargetIP:    report.TargetIP,
		DBType:      report.DBType,
		DBTypeName:  enums.BaselineEnum.GetDBTypeName(report.DBType),
		HighCount:   report.HighCount,
		MiddleCount: report.MiddleCount,
		LowCount:    report.LowCount,
		StartTime:   report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:     report.EndTime.Format("2006-01-02 15:04:05"),
	}

	for _, r := range report.Results {
		resp.Results = append(resp.Results, typespec.SensitiveDataItem{
			ID:            r.ID,
			DBName:        r.DBName,
			TableName:     r.TableNameStr,
			ColumnName:    r.ColumnName,
			DataType:      r.DataType,
			DataTypeName:  enums.BaselineEnum.GetSensitiveDataTypeName(r.DataType),
			DataLevel:     r.DataLevel,
			DataLevelName: enums.BaselineEnum.GetSensitiveDataLevelName(r.DataLevel),
			MatchRule:     r.MatchRule,
			SampleData:    r.SampleData,
			CreateTime:    r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (a *BaselineApp) GetSensitiveDataResults(ctx context.Context, req *typespec.SensitiveDataListReq) (*typespec.SensitiveDataListResp, error) {
	var model mysqls.SensitiveDataResult
	var list []mysqls.SensitiveDataResult
	var err error

	if req.TargetID > 0 {
		list, err = model.GetByTargetID(ctx, req.TargetID)
	} else {
		list, err = model.GetByTaskID(ctx, req.TaskID)
	}
	if err != nil {
		return nil, err
	}

	resp := &typespec.SensitiveDataListResp{
		Total: int64(len(list)),
	}
	for _, r := range list {
		resp.List = append(resp.List, typespec.SensitiveDataItem{
			ID:            r.ID,
			DBName:        r.DBName,
			TableName:     r.TableNameStr,
			ColumnName:    r.ColumnName,
			DataType:      r.DataType,
			DataTypeName:  enums.BaselineEnum.GetSensitiveDataTypeName(r.DataType),
			DataLevel:     r.DataLevel,
			DataLevelName: enums.BaselineEnum.GetSensitiveDataLevelName(r.DataLevel),
			MatchRule:     r.MatchRule,
			SampleData:    r.SampleData,
			CreateTime:    r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func (a *BaselineApp) GetSensitiveDataStat(ctx context.Context, req *typespec.SensitiveDataStatReq) (*typespec.SensitiveDataStatResp, error) {
	var model mysqls.SensitiveDataResult
	high, middle, low, total, err := model.GetStatByTaskID(ctx, req.TaskID)
	if err != nil {
		return nil, err
	}
	return &typespec.SensitiveDataStatResp{
		HighCount:   int(high),
		MiddleCount: int(middle),
		LowCount:    int(low),
		TotalCount:  int(total),
	}, nil
}

func (a *BaselineApp) GetEnums(ctx context.Context) interface{} {
	return map[string]interface{}{
		"osTypes": []map[string]interface{}{
			{"value": enums.BaselineOSTypeLinux, "label": enums.BaselineEnum.GetOSTypeName(enums.BaselineOSTypeLinux)},
			{"value": enums.BaselineOSTypeWindows, "label": enums.BaselineEnum.GetOSTypeName(enums.BaselineOSTypeWindows)},
			{"value": enums.BaselineOSTypeDomestic, "label": enums.BaselineEnum.GetOSTypeName(enums.BaselineOSTypeDomestic)},
			{"value": enums.BaselineOSTypeEmbedded, "label": enums.BaselineEnum.GetOSTypeName(enums.BaselineOSTypeEmbedded)},
		},
		"dbTypes": []map[string]interface{}{
			{"value": enums.DBSupportTypeMySQL, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypeMySQL)},
			{"value": enums.DBSupportTypePostgreSQL, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypePostgreSQL)},
			{"value": enums.DBSupportTypeMongoDB, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypeMongoDB)},
			{"value": enums.DBSupportTypeRedis, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypeRedis)},
			{"value": enums.DBSupportTypeCouchDB, "label": enums.BaselineEnum.GetDBTypeName(enums.DBSupportTypeCouchDB)},
		},
		"checkCategories": []map[string]interface{}{
			{"value": enums.BaselineCategoryPasswordPolicy, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryPasswordPolicy)},
			{"value": enums.BaselineCategoryUserPermission, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryUserPermission)},
			{"value": enums.BaselineCategoryFirewall, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryFirewall)},
			{"value": enums.BaselineCategoryKernelSecurity, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryKernelSecurity)},
			{"value": enums.BaselineCategoryAuditLog, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryAuditLog)},
			{"value": enums.BaselineCategoryNetworkService, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategoryNetworkService)},
			{"value": enums.BaselineCategorySSHConfig, "label": enums.BaselineEnum.GetCategoryName(enums.BaselineCategorySSHConfig)},
		},
		"dbCheckCategories": []map[string]interface{}{
			{"value": enums.DBCheckCategoryAuthentication, "label": enums.BaselineEnum.GetDBCheckCategoryName(enums.DBCheckCategoryAuthentication)},
			{"value": enums.DBCheckCategoryAuthorization, "label": enums.BaselineEnum.GetDBCheckCategoryName(enums.DBCheckCategoryAuthorization)},
			{"value": enums.DBCheckCategoryConfigSecure, "label": enums.BaselineEnum.GetDBCheckCategoryName(enums.DBCheckCategoryConfigSecure)},
			{"value": enums.DBCheckCategoryAuditLog, "label": enums.BaselineEnum.GetDBCheckCategoryName(enums.DBCheckCategoryAuditLog)},
		},
		"sensitiveDataLevels": []map[string]interface{}{
			{"value": enums.SensitiveDataLevelHigh, "label": enums.BaselineEnum.GetSensitiveDataLevelName(enums.SensitiveDataLevelHigh)},
			{"value": enums.SensitiveDataLevelMiddle, "label": enums.BaselineEnum.GetSensitiveDataLevelName(enums.SensitiveDataLevelMiddle)},
			{"value": enums.SensitiveDataLevelLow, "label": enums.BaselineEnum.GetSensitiveDataLevelName(enums.SensitiveDataLevelLow)},
		},
		"checkResults": []map[string]interface{}{
			{"value": enums.BaselineCheckResultPass, "label": enums.BaselineEnum.GetCheckResultName(enums.BaselineCheckResultPass)},
			{"value": enums.BaselineCheckResultFail, "label": enums.BaselineEnum.GetCheckResultName(enums.BaselineCheckResultFail)},
			{"value": enums.BaselineCheckResultError, "label": enums.BaselineEnum.GetCheckResultName(enums.BaselineCheckResultError)},
		},
	}
}
