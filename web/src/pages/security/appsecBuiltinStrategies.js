/** 内置扫描策略（新建任务第 1 步仅展示此类模板） */

import { loadBuiltinOverride, mergeScanConfig } from './appsecStrategyStorage.js'

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

const BUILTIN_CONFIGS = {
  'builtin-full': {
    testMode: 'principle',
    safeTest: true,
    vulExploit: false,
    testIntensity: 3,
    vulIdsConfig: [],
    webCrawler: { isOpen: true, maxDepth: 5, scanRange: 0, crawlerSpeed: 2 },
    portScan: { isOpen: true, scanPort: '21,22,23,80,443,445,3306,8000,8080', tcpScanType: 1, timeout: 10, concurrent: 100 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  },
  'builtin-highrisk': {
    testMode: 'principle',
    safeTest: true,
    vulExploit: false,
    testIntensity: 4,
    vulIdsConfig: [],
    webCrawler: { isOpen: true, maxDepth: 3, scanRange: 0, crawlerSpeed: 2 },
    portScan: { isOpen: true, scanPort: '80,443,3306,8080', tcpScanType: 1, timeout: 10, concurrent: 100 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  },
  'builtin-web': {
    testMode: 'principle',
    safeTest: true,
    vulExploit: false,
    testIntensity: 3,
    vulIdsConfig: [],
    webCrawler: { isOpen: true, maxDepth: 6, scanRange: 0, crawlerSpeed: 3 },
    portScan: { isOpen: false, scanPort: '', tcpScanType: 1, timeout: 10, concurrent: 100 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  },
  'builtin-weakpass': {
    testMode: 'principle',
    safeTest: false,
    vulExploit: false,
    testIntensity: 2,
    vulIdsConfig: [],
    webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
    portScan: { isOpen: true, scanPort: '21,22,23,3306,3389,5432,6379', tcpScanType: 1, timeout: 10, concurrent: 50 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  },
  'builtin-component': {
    testMode: 'version',
    safeTest: true,
    vulExploit: false,
    testIntensity: 3,
    vulIdsConfig: [],
    webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
    portScan: { isOpen: true, scanPort: '80,443,8080,8000,8443', tcpScanType: 1, timeout: 10, concurrent: 100 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  },
  'builtin-portscan': {
    testMode: 'principle',
    safeTest: true,
    vulExploit: false,
    testIntensity: 1,
    vulIdsConfig: [],
    webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
    portScan: { isOpen: true, scanPort: '1-65535', tcpScanType: 2, timeout: 5, concurrent: 200 },
    proxy: { isOpen: false, proto: 1, addr: '', port: '' }
  }
}

/** 各策略在任务配置页显示的区块 */
const STRATEGY_SECTIONS = {
  'builtin-full': { vuln: true, scan: true, crawler: true, port: true, advanced: true },
  'builtin-highrisk': { vuln: true, scan: true, crawler: true, port: true, advanced: true },
  'builtin-web': { vuln: true, scan: true, crawler: true, port: false, advanced: true },
  'builtin-weakpass': { vuln: true, scan: true, crawler: false, port: true, advanced: false },
  'builtin-component': { vuln: true, scan: true, crawler: false, port: true, advanced: true },
  'builtin-portscan': { vuln: false, scan: true, crawler: false, port: true, advanced: false }
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
  return STRATEGY_SECTIONS[strategyId] || STRATEGY_SECTIONS['builtin-full']
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
