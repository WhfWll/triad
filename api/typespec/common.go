package typespec

// 全局 选项内容
type GlobalOptionsRes struct {
	VulScriptType                    []GlobalOptionsItemRes `json:"vul_script_type"`
	VulScriptVerifyType              []GlobalOptionsItemRes `json:"vul_script_verify_type"`
	VulScriptEvidenceType            []GlobalOptionsItemRes `json:"vul_script_evidence_type"`
	VulScriptExploitImpact           []GlobalOptionsItemRes `json:"vul_script_exploit_impact"`
	VulLibrariesClass                []GlobalOptionsItemRes `json:"vul_libraries_class"`
	VulLibrariesType                 []GlobalOptionsItemRes `json:"vul_libraries_type"`
	VulLibrariesRisk                 []GlobalOptionsItemRes `json:"vul_libraries_risk"`
	TaskCheckTaskAttackDecisionModel []GlobalOptionsItemRes `json:"task_check_task_attack_decision_model"`
	TaskCheckTaskExecType            []GlobalOptionsItemRes `json:"task_check_task_exec_type"`
	TaskCheckTaskInvokeApi           []GlobalOptionsItemRes `json:"task_check_task_invoke_api"`
	TaskCheckTaskPortScanType        []GlobalOptionsItemRes `json:"task_check_task_port_scan_type"`
	TaskCheckTaskRuntime             []GlobalOptionsItemRes `json:"task_check_task_runtime"`
	TaskCheckTaskStatus              []GlobalOptionsItemRes `json:"task_check_task_status"`
	TaskCheckTaskType                []GlobalOptionsItemRes `json:"task_check_task_type"`
	TaskCheckTaskWebsiteLoginType    []GlobalOptionsItemRes `json:"task_check_task_website_login_type"`
}
type GlobalOptionsItemRes struct {
	Value     interface{} `json:"value"`
	Label     interface{} `json:"label"`
	IsDefault bool        `json:"isDefault,omitempty"`
}
