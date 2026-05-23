import axios from '@/axios/http'

const security = {
  // 恶意代码检测 (后端路由: /smart/malware/*)
  runMalwareScan: (data) => {
    return axios.postJson('/smart/malware/scan', data)
  },
  getMalwareList: (params) => {
    return axios.get('/smart/malware/result', params)
  },
  getMalwareTaskList: (params) => {
    return axios.get('/smart/malware/tasks', params)
  },

  // 数据库安全检查（任务化 /smart/datasec/db/*，旧 /smart/db/* 仍可用）
  testDataSecDBConn: (data) => {
    return axios.postJson('/smart/datasec/db/test-conn', data)
  },
  runDBCheck: (data) => {
    return axios.postJson('/smart/datasec/db/run', data)
  },
  getDBCheckList: (params) => {
    return axios.get('/smart/datasec/db/list', params)
  },
  getDBCheckDetail: (params) => {
    return axios.get('/smart/datasec/db/detail', params)
  },
  rerunDataSecTask: (data) => {
    return axios.postJson('/smart/datasec/task/rerun', data)
  },
  deleteDataSecTask: (params) => {
    return axios.get('/smart/datasec/task/delete', params)
  },
  cloneDataSecTaskTargets: (params) => {
    return axios.get('/smart/datasec/task/clone-targets', params)
  },

  // 数据库目标库
  getDatasecTargetList: (params) => {
    return axios.get('/smart/datasec/target/list', params)
  },
  saveDatasecTarget: (data) => {
    return axios.postJson('/smart/datasec/target/save', data)
  },
  deleteDatasecTarget: (params) => {
    return axios.get('/smart/datasec/target/delete', params)
  },
  importDatasecTargets: (data) => {
    return axios.postJson('/smart/datasec/target/import', data)
  },
  exportDatasecTargets: (params) => {
    return axios.get('/smart/datasec/target/export', params)
  },
  saveDatasecTargetsFromTask: (data) => {
    return axios.postJson('/smart/datasec/target/save-from-task', data)
  },
  testDatasecTargetConn: (data) => {
    return axios.postJson('/smart/datasec/target/test-conn', data)
  },
  batchTestDatasecTargetConn: (data) => {
    return axios.postJson('/smart/datasec/target/batch-test-conn', data)
  },

  // 敏感数据发现（任务化 /smart/datasec/sensitive/*）
  runSensitiveScan: (data) => {
    return axios.postJson('/smart/datasec/sensitive/run', data)
  },
  getSensitiveDataList: (params) => {
    return axios.get('/smart/datasec/sensitive/list', params)
  },
  getSensitiveScanDetail: (params) => {
    return axios.get('/smart/datasec/sensitive/detail', params)
  },

  // 数据安全检测规则 (/smart/datasec/rules/*)
  getDatasecRules: () => {
    return axios.get('/smart/datasec/rules')
  },
  reloadDatasecRules: () => {
    return axios.postJson('/smart/datasec/rules/reload', {})
  },
  importDatasecRules: (data) => {
    return axios.postJson('/smart/datasec/rules/import', data)
  },
  importDatasecBuiltinRules: () => {
    return axios.postJson('/smart/datasec/rules/import-builtin', {})
  },
  previewDatasecCveImport: () => {
    return axios.get('/smart/datasec/rules/cve-preview')
  },
  importDatasecRulesFromCve: (data) => {
    return axios.postJson('/smart/datasec/rules/import-cve', data || {})
  },
  getDatasecRuleDetail: (params) => {
    return axios.get('/smart/datasec/rule/detail', params)
  },
  createDatasecRule: (data) => {
    return axios.postJson('/smart/datasec/rule/create', data)
  },
  updateDatasecRule: (data) => {
    return axios.postJson('/smart/datasec/rule/update', data)
  },
  deleteDatasecRule: (params) => {
    return axios.get('/smart/datasec/rule/delete', params)
  },
  getSensitiveStat: (params) => {
    return axios.get('/smart/sensitive/stat', params)
  },

  // 基线检查 (后端路由: /smart/baseline/*)
  runBaselineCheck: (data) => {
    return axios.postJson('/smart/baseline/check', data)
  },
  runBaselineBatchCheck: (data) => {
    return axios.postJson('/smart/baseline/check/batch', data)
  },
  getBaselineBatchProgress: (params) => {
    return axios.get('/smart/baseline/check/progress', params)
  },
  getBaselineList: (params) => {
    return axios.get('/smart/baseline/result', params)
  },
  getBaselineStat: (params) => {
    return axios.get('/smart/baseline/stat', params)
  },
  getBaselineTaskList: (params) => {
    return axios.get('/smart/baseline/tasks', params)
  },
  getBaselineTaskTargets: (params) => {
    return axios.get('/smart/baseline/task/targets', params)
  },
  getBaselineRules: () => {
    return axios.get('/smart/baseline/rules')
  },
  getBaselineRulesFromDB: () => {
    return axios.get('/smart/baseline/rules/db')
  },
  getBaselineRuleDetail: (params) => {
    return axios.get('/smart/baseline/rule/detail', params)
  },
  createBaselineRule: (data) => {
    return axios.postJson('/smart/baseline/rule/create', data)
  },
  updateBaselineRule: (data) => {
    return axios.postJson('/smart/baseline/rule/update', data)
  },
  deleteBaselineRule: (params) => {
    return axios.get('/smart/baseline/rule/delete', params)
  },
  importBaselineRules: (data) => {
    return axios.postJson('/smart/baseline/rules/import', data)
  },
  getBaselineEnums: () => {
    return axios.get('/smart/baseline/enums')
  },

  // 专项应用检测（后端调用 scanner.exe）
  getAppSpecificList: (params) => {
    return axios.get('/smart/appsec/appspecific/list', params)
  },
  runAppSpecificScan: (data) => {
    return axios.postJson('/smart/appsec/appspecific/run', data)
  },

  // 动态扫描（后端调用 scanner.exe）
  getDynamicScanList: (params) => {
    return axios.get('/smart/appsec/dynamic/list', params)
  },
  runDynamicScan: (data) => {
    return axios.postJson('/smart/appsec/dynamic/run', data)
  },
  getDynamicScanDetail: (params) => {
    return axios.get('/smart/appsec/dynamic/detail', params)
  },
  getAppSpecificScanDetail: (params) => {
    return axios.get('/smart/appsec/appspecific/detail', params)
  },

  // CVE漏洞扫描 (后端路由: /smart/vulnscan/cve/*)
  runCveScan: (data) => {
    return axios.postJson('/smart/vulnscan/cve/run', data)
  },
  runCveBatchScan: (data) => {
    return axios.postJson('/smart/vulnscan/cve/batch', data)
  },
  getCveBatchProgress: (params) => {
    return axios.get('/smart/vulnscan/cve/progress', params)
  },

  // YARA恶意代码检测 (后端路由: /smart/malware/yara/*)
  runYaraScan: (data) => {
    return axios.postJson('/smart/malware/yara/run', data)
  },
  runYaraBatchScan: (data) => {
    return axios.postJson('/smart/malware/yara/batch', data)
  },
  getYaraBatchProgress: (params) => {
    return axios.get('/smart/malware/yara/progress', params)
  },
  getYaraResultList: (params) => {
    return axios.get('/smart/malware/yara/result', params)
  },
  getYaraTaskList: (params) => {
    return axios.get('/smart/malware/yara/tasks', params)
  },

  // 病毒库规则管理 (后端路由: /smart/malware/rule/*)
  getMalwareRuleList: (params) => {
    return axios.get('/smart/malware/rules', params)
  },
  getMalwareRuleDetail: (params) => {
    return axios.get('/smart/malware/rule/detail', params)
  },
  createMalwareRule: (data) => {
    return axios.postJson('/smart/malware/rule/create', data)
  },
  updateMalwareRule: (data) => {
    return axios.postJson('/smart/malware/rule/update', data)
  },
  deleteMalwareRule: (params) => {
    return axios.get('/smart/malware/rule/delete', params)
  },
  importMalwareRules: (data) => {
    return axios.postFormData('/smart/malware/rule/import', data)
  },

  // CVE漏洞库查询 (后端路由: /smart/cvedb/*)
  getCveDBInfo: (params) => {
    return axios.get('/smart/cvedb/info', params)
  },
  queryCveDB: (params) => {
    return axios.get('/smart/cvedb/query', params)
  },

  // 扫描策略管理 (后端路由: /smart/strategy/*)
  getStrategyList: (params) => {
    return axios.get('/smart/strategy/list', params)
  },
  getStrategyDetail: (params) => {
    return axios.get('/smart/strategy/detail', params)
  },
  createStrategy: (data) => {
    return axios.postJson('/smart/strategy/create', data)
  },
  updateStrategy: (data) => {
    return axios.postJson('/smart/strategy/update', data)
  },
  deleteStrategy: (params) => {
    return axios.get('/smart/strategy/delete', params)
  }
}

export default security