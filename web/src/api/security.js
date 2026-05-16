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

  // 数据库安全检查 (后端路由: /smart/db/*)
  runDBCheck: (data) => {
    return axios.postJson('/smart/db/check', data)
  },
  getDBCheckList: (params) => {
    return axios.get('/smart/db/result', params)
  },

  // 敏感数据发现 (后端路由: /smart/sensitive/*)
  runSensitiveScan: (data) => {
    return axios.postJson('/smart/sensitive/scan', data)
  },
  getSensitiveDataList: (params) => {
    return axios.get('/smart/sensitive/result', params)
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

  // 专项应用检测 (后端暂无接口)
  getAppSpecificList: () => {
    return Promise.resolve({ code: 200, data: { list: [], total: 0 } })
  },
  runAppSpecificScan: () => {
    return Promise.resolve({ code: 200, msg: '功能开发中' })
  },

  // 动态扫描 (后端暂无接口)
  getDynamicScanList: () => {
    return Promise.resolve({ code: 200, data: { list: [], total: 0 } })
  },
  runDynamicScan: () => {
    return Promise.resolve({ code: 200, msg: '功能开发中' })
  }
}

export default security