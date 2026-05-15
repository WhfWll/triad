import request from '@/utils/request'

const security = {
  // 恶意代码检测
  runMalwareScan: (data) => {
    return request({
      url: '/api/smart/baseline/malware/scan',
      method: 'post',
      data
    })
  },
  getMalwareList: (params) => {
    return request({
      url: '/api/smart/baseline/malware/list',
      method: 'get',
      params
    })
  },
  getMalwareDetail: (params) => {
    return request({
      url: '/api/smart/baseline/malware/detail',
      method: 'get',
      params
    })
  },
  delMalwareTask: (params) => {
    return request({
      url: '/api/smart/baseline/malware/del',
      method: 'get',
      params
    })
  },

  // 数据库安全检查
  runDBCheck: (data) => {
    return request({
      url: '/api/smart/baseline/db/check',
      method: 'post',
      data
    })
  },
  getDBCheckList: (params) => {
    return request({
      url: '/api/smart/baseline/db/list',
      method: 'get',
      params
    })
  },
  getDBCheckDetail: (params) => {
    return request({
      url: '/api/smart/baseline/db/detail',
      method: 'get',
      params
    })
  },
  delDBCheckTask: (params) => {
    return request({
      url: '/api/smart/baseline/db/del',
      method: 'get',
      params
    })
  },

  // 敏感数据发现
  runSensitiveScan: (data) => {
    return request({
      url: '/api/smart/baseline/sensitive/scan',
      method: 'post',
      data
    })
  },
  getSensitiveDataList: (params) => {
    return request({
      url: '/api/smart/baseline/sensitive/list',
      method: 'get',
      params
    })
  },
  getSensitiveDetail: (params) => {
    return request({
      url: '/api/smart/baseline/sensitive/detail',
      method: 'get',
      params
    })
  },
  getSensitiveStat: (params) => {
    return request({
      url: '/api/smart/baseline/sensitive/stat',
      method: 'get',
      params
    })
  },
  delSensitiveTask: (params) => {
    return request({
      url: '/api/smart/baseline/sensitive/del',
      method: 'get',
      params
    })
  },

  // 基线检查
  runBaselineCheck: (data) => {
    return request({
      url: '/api/smart/baseline/check',
      method: 'post',
      data
    })
  },
  getBaselineList: (params) => {
    return request({
      url: '/api/smart/baseline/list',
      method: 'get',
      params
    })
  },
  getBaselineDetail: (params) => {
    return request({
      url: '/api/smart/baseline/detail',
      method: 'get',
      params
    })
  },
  getBaselineStat: (params) => {
    return request({
      url: '/api/smart/baseline/stat',
      method: 'get',
      params
    })
  },
  getBaselineRules: () => {
    return request({
      url: '/api/smart/baseline/rules',
      method: 'get'
    })
  },
  getBaselineEnums: () =&gt; {
    return request({
      url: '/api/smart/baseline/enums',
      method: 'get'
    })
  },

  // 专项应用检测
  runAppSpecificScan: (data) =&gt; {
    return request({
      url: '/api/smart/appspecific/scan',
      method: 'post',
      data
    })
  },
  getAppSpecificList: (params) =&gt; {
    return request({
      url: '/api/smart/appspecific/list',
      method: 'get',
      params
    })
  },
  getAppSpecificDetail: (params) =&gt; {
    return request({
      url: '/api/smart/appspecific/detail',
      method: 'get',
      params
    })
  },
  delAppSpecificTask: (params) =&gt; {
    return request({
      url: '/api/smart/appspecific/del',
      method: 'get',
      params
    })
  },

  // 动态扫描
  runDynamicScan: (data) =&gt; {
    return request({
      url: '/api/smart/dynamicscan/scan',
      method: 'post',
      data
    })
  },
  getDynamicScanList: (params) =&gt; {
    return request({
      url: '/api/smart/dynamicscan/list',
      method: 'get',
      params
    })
  },
  getDynamicScanDetail: (params) =&gt; {
    return request({
      url: '/api/smart/dynamicscan/detail',
      method: 'get',
      params
    })
  },
  delDynamicScanTask: (params) =&gt; {
    return request({
      url: '/api/smart/dynamicscan/del',
      method: 'get',
      params
    })
  }
}

export default security