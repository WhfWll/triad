package config

// PocExploitConfig POC漏洞利用配置
type PocExploitConfig struct {
	PocName       string // POC名称
	ExploiterType string // 漏洞利用器类型
	HostPort      int    // 主机端口
	Protocol      string // 协议类型
}

// 支持漏洞利用的POC配置列表（统一配置）
var exploitablePocConfigs = []PocExploitConfig{
	{
		PocName:       "mitm_sql_check",
		ExploiterType: "sql_injection",
		HostPort:      0,
		Protocol:      "http",
	},
	{
		PocName:       "ms17010_eternalblue_check",
		ExploiterType: "eternal_blue",
		HostPort:      445,
		Protocol:      "smb",
	},
	{
		PocName:       "apache_2.4.49_path_traversal_and_file_disclosure",
		ExploiterType: "apache_cve_2021_41773",
		HostPort:      0,
		Protocol:      "http",
	},
}

// GetPocExploitConfig 根据POC名称获取漏洞利用配置
func GetPocExploitConfig(pocName string) *PocExploitConfig {
	for _, config := range exploitablePocConfigs {
		if config.PocName == pocName {
			return &config
		}
	}
	return nil
}

// IsExploitablePoc 检查POC名称是否支持漏洞利用
func IsExploitablePoc(pocName string) bool {
	return GetPocExploitConfig(pocName) != nil
}
