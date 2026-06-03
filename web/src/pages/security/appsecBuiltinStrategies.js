/** 内置扫描策略（新建任务第 1 步仅展示此类模板） */

import { loadBuiltinOverride, mergeScanConfig } from './appsecStrategyStorage.js'
import {
  cloneDefaultWeakPass,
  WEAKPASS_DEFAULT_SERVICE_IDS
} from './appsecWeakPassDefaults.js'

export const BUILTIN_STRATEGIES = [
  {
    id: 'builtin-full',
    name: '全漏洞扫描',
    desc: '覆盖所有已知漏洞类型，启用爬虫与端口扫描',
    icon: '🛡'
  },
  {
    id: 'builtin-highrisk',
    name: '高危漏洞扫描',
    desc: '仅高危/严重等级的漏洞脚本',
    icon: '🔴'
  },
  {
    id: 'builtin-web',
    name: 'Web漏洞扫描',
    desc: '专注 Web 应用漏洞，启用爬虫深入抓取',
    icon: '🌐'
  },
  {
    id: 'builtin-weakpass',
    name: '弱口令扫描',
    desc: '检测常见服务的弱口令漏洞',
    icon: '🔑'
  },
  {
    id: 'builtin-component',
    name: '组件漏洞扫描',
    desc: '扫描第三方组件和中间件已知漏洞',
    icon: '📦'
  },
  {
    id: 'builtin-portscan',
    name: '端口扫描',
    desc: '快速发现开放端口与运行服务',
    icon: '🔍'
  }
]

/** 扫描评估默认：仅原理验证，其余关闭（界面已隐藏，保存/建任务时仍写入） */
export const DEFAULT_SCAN_ASSESSMENT = {
  testMode: 'principle',
  safeTest: false,
  vulExploit: false
}

const DEFAULT_WEBSITE_LOGIN = { isOpen: false, list: [] }

function withAssessmentDefaults(cfg) {
  const websiteLogin = cfg.websiteLogin
    ? { ...DEFAULT_WEBSITE_LOGIN, ...cfg.websiteLogin, list: [...(cfg.websiteLogin.list || [])] }
    : JSON.parse(JSON.stringify(DEFAULT_WEBSITE_LOGIN))
  return {
    ...DEFAULT_SCAN_ASSESSMENT,
    websiteLogin: JSON.parse(JSON.stringify(DEFAULT_WEBSITE_LOGIN)),
    weakPass: cloneDefaultWeakPass(),
    ...cfg,
    ...DEFAULT_SCAN_ASSESSMENT,
    websiteLogin,
    weakPass: cfg.weakPass ? cloneDefaultWeakPass(cfg.weakPass) : cloneDefaultWeakPass()
  }
}

const BUILTIN_CONFIGS = {
  'builtin-full': withAssessmentDefaults({
    vulIdsConfig: [],
    webCrawler: { isOpen: true, maxDepth: 5, scanRange: 0, crawlerSpeed: 2 },
    portScan: { isOpen: true, scanPort: '21,22,23,80,443,445,3306,8000,8080', tcpScanType: 1, timeout: 10, concurrent: 100 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  }),
  'builtin-highrisk': withAssessmentDefaults({
    vulIdsConfig: [],
    webCrawler: { isOpen: true, maxDepth: 3, scanRange: 0, crawlerSpeed: 2 },
    portScan: { isOpen: true, scanPort: '80,443,3306,8080', tcpScanType: 1, timeout: 10, concurrent: 100 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  }),
  'builtin-web': withAssessmentDefaults({
    vulIdsConfig: [],
    webCrawler: { isOpen: true, maxDepth: 6, scanRange: 0, crawlerSpeed: 3 },
    portScan: { isOpen: false, scanPort: '', tcpScanType: 1, timeout: 10, concurrent: 100 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  }),
  'builtin-weakpass': withAssessmentDefaults({
    vulIdsConfig: [],
    webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
    portScan: { isOpen: true, scanPort: '21,22,23,3306,3389,5432,6379', tcpScanType: 1, timeout: 10, concurrent: 50 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' },
    weakPass: cloneDefaultWeakPass({
      isOpen: true,
      services: [...WEAKPASS_DEFAULT_SERVICE_IDS]
    })
  }),
  'builtin-component': withAssessmentDefaults({
    vulIdsConfig: [],
    webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
    portScan: { isOpen: true, scanPort: '80,443,8080,8000,8443', tcpScanType: 1, timeout: 10, concurrent: 100 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  }),
  'builtin-portscan': withAssessmentDefaults({
    vulIdsConfig: [],
    webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
    portScan: { isOpen: true, scanPort: '1-65535', tcpScanType: 2, timeout: 5, concurrent: 200 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  })
}

/** 各策略在任务配置页显示的区块 */
const STRATEGY_SECTIONS = {
  'builtin-full': { vuln: true, scan: true, crawler: true, port: true, advanced: true, login: true },
  'builtin-highrisk': { vuln: true, scan: true, crawler: true, port: true, advanced: true, login: true },
  'builtin-web': { vuln: true, scan: true, crawler: true, port: false, advanced: true, login: true },
  'builtin-weakpass': { vuln: true, scan: true, crawler: false, port: true, advanced: false, login: true, weakPass: true },
  'builtin-component': { vuln: true, scan: true, crawler: false, port: true, advanced: true, login: true },
  'builtin-portscan': { vuln: false, scan: true, crawler: false, port: true, advanced: false, login: false }
}

export function getBuiltinStrategy(id) {
  return BUILTIN_STRATEGIES.find(s => s.id === id) || null
}

export function isBuiltinStrategyId(id) {
  return Boolean(BUILTIN_CONFIGS[id])
}

export function cloneBuiltinConfig(id) {
  const cfg = BUILTIN_CONFIGS[id] || BUILTIN_CONFIGS['builtin-full']
  return JSON.parse(JSON.stringify(cfg))
}

export function getStrategySections(strategyId) {
  const s = STRATEGY_SECTIONS[strategyId] || STRATEGY_SECTIONS['builtin-full']
  return { login: false, ...s }
}

/** 统一扫描评估默认值（界面已移除扫描评估页） */
export function applyScanAssessmentDefaults(config) {
  if (!config) return config
  Object.assign(config, DEFAULT_SCAN_ASSESSMENT)
  delete config.testIntensity
  if (!config.weakPass) {
    config.weakPass = cloneDefaultWeakPass()
  } else {
    config.weakPass = cloneDefaultWeakPass(config.weakPass)
  }
  if (!config.websiteLogin) {
    config.websiteLogin = JSON.parse(JSON.stringify(DEFAULT_WEBSITE_LOGIN))
  } else if (!Array.isArray(config.websiteLogin.list)) {
    config.websiteLogin.list = []
  } else {
    config.websiteLogin.list = config.websiteLogin.list.map(row => {
      const { _editing, ...rest } = row || {}
      return rest
    })
  }
  return config
}

/** 供策略管理页使用的完整内置策略列表（含 config，合并本地覆盖） */
export function getBuiltinStrategiesWithConfig() {
  return BUILTIN_STRATEGIES.map(s => {
    const config = cloneBuiltinConfig(s.id)
    const override = loadBuiltinOverride(s.id)
    const merged = override && override.config ? mergeScanConfig(config, override.config) : config
    const vulnCount = (merged.vulIdsConfig && merged.vulIdsConfig.length) || 0
    return {
      ...s,
      builtin: true,
      name: (override && override.name) || s.name,
      desc: (override && override.desc) || s.desc,
      icon: (override && override.icon) || s.icon,
      vulnCount,
      config: merged
    }
  })
}

export { loadCustomStrategy } from './appsecStrategyStorage.js'
