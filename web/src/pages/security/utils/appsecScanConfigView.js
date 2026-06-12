import { getTestModeLabel } from '../appsecTaskLabels.js'
import { PORT_RANGE_MAP } from '../appsecPortRanges.js'

const PORT_RANGE_LABELS = {
  10: 'TOP10 端口',
  20: 'TOP20 端口',
  50: 'TOP50 端口',
  100: 'TOP100 端口',
  500: 'TOP500 端口',
  1000: 'TOP1000 端口',
  65535: '全部端口 (1-65535)',
  0: '自定义端口'
}

const TCP_SCAN_LABELS = {
  1: 'TCP-Connect',
  2: 'TCP SYN',
  3: 'TCP FIN',
  4: 'TCP ACK',
  5: 'TCP NULL',
  6: 'UDP'
}

const CRAWLER_RANGE_LABELS = {
  0: '全域名扫描',
  1: '目标 URL 和子目录'
}

const CRAWLER_SPEED_LABELS = {
  1: '高速',
  2: '中速',
  3: '低速'
}

const CRAWLER_DEPTH_LABELS = {
  0: '不限',
  1: '1',
  2: '2',
  3: '3',
  4: '4',
  5: '5'
}

const CRAWLER_MAX_URL_LABELS = {
  0: '不限',
  100: '100',
  200: '200',
  500: '500',
  1000: '1000',
  5000: '5000',
  10000: '10000',
  100000: '100000'
}

const CRAWLER_REPEAT_LABELS = {
  0: '不限',
  1: 'page / method / query / post 敏感',
  2: 'page / method / query 敏感',
  3: 'page / method / query-name 敏感',
  4: 'page / method 敏感',
  5: 'page 敏感'
}

const PATH_SPEED_LABELS = { 1: '高速', 2: '中速', 3: '低速' }
const PATH_TIME_LABELS = { 0: '不限', 1: '1 min', 3: '3 min', 4: '5 min', 10: '10 min', 30: '30 min', 60: '60 min' }

const WEAK_DICT_LABELS = { 1: '默认字典', 2: '通用字典', 3: '补充字典' }
const WEAK_RATE_LABELS = { 1: '高速', 2: '中速', 3: '低速' }

const ALIVE_PROBE_LABELS = {
  1: 'ICMP-PING',
  2: 'ARP-PING',
  3: 'TCP-PING',
  4: 'UDP-PING',
  5: 'TCP-ACK',
  6: 'TCP-SYN'
}

const PROXY_PROTO_LABELS = { 1: 'HTTP', 2: 'HTTPS', 3: 'SOCKS4', 4: 'SOCKS5' }

const LOGIN_TYPE_LABELS = { 1: 'Header', 2: 'Cookie', 3: '账号密码', 4: '登录序列' }

const LATERAL_STRATEGY_LABELS = {
  same_subnet: '同网段探测',
  neighbor: '邻居发现',
  exclude_current: '排除同网段',
  custom_range: '自定义范围',
  auto_subnet: '自动子网'
}

function pickBlock(cfg, primary, alt) {
  if (!cfg || typeof cfg !== 'object') return {}
  return cfg[primary] || cfg[alt] || {}
}

function yesNo(val) {
  return val ? '是' : '否'
}

function enabled(val) {
  return Boolean(val)
}

function display(val, fallback = '-') {
  if (val === null || val === undefined || val === '') return fallback
  if (typeof val === 'boolean') return yesNo(val)
  return String(val)
}

function labelFrom(map, val, zhField, obj) {
  if (obj && zhField && obj[zhField]) return obj[zhField]
  if (val === null || val === undefined || val === '') return '-'
  const key = Number(val)
  if (Number.isFinite(key) && map[key] !== undefined) return map[key]
  return display(val)
}

function timeoutLabel(sec, zhField, obj) {
  if (obj && obj[zhField]) return obj[zhField]
  if (!sec && sec !== 0) return '-'
  return `${sec}s`
}

function portRangeLabel(ps) {
  if (ps.portScanTypeZh) return ps.portScanTypeZh
  const type = ps.portScanType
  if (type === 0 || type === '0') {
    return ps.scanPort ? `自定义：${ps.scanPort}` : '自定义端口'
  }
  const label = PORT_RANGE_LABELS[type]
  if (label) return label
  return display(type)
}

function portDetail(ps) {
  const type = Number(ps.portScanType)
  if (type === 0 && ps.scanPort) return ps.scanPort
  if (PORT_RANGE_MAP[type]) {
    const ports = PORT_RANGE_MAP[type]
    return ports.length > 120 ? `${ports.slice(0, 120)}…` : ports
  }
  return ps.scanPort || '-'
}

function joinList(list, zhList) {
  if (Array.isArray(zhList) && zhList.length) return zhList.join('、')
  if (Array.isArray(list) && list.length) return list.join('、')
  return '-'
}

function buildItems(entries) {
  return entries
    .filter(item => item && item.label)
    .map(item => ({
      label: item.label,
      value: item.value != null && item.value !== '' ? item.value : '-',
      full: Boolean(item.full),
      mono: Boolean(item.mono)
    }))
}

export function buildScanConfigSections(cfg) {
  if (!cfg || typeof cfg !== 'object') return { basics: [], sections: [], vulIdsCount: 0, loginRows: [] }

  const wc = pickBlock(cfg, 'webCrawlerConfig', 'webCrawler')
  const ps = pickBlock(cfg, 'portScanConfig', 'portScan')
  const px = pickBlock(cfg, 'proxyConfig', 'proxy')
  const wp = pickBlock(cfg, 'weakPassConfig', 'weakPass')
  const wps = pickBlock(cfg, 'webPathScanConfig', 'webPathScan')
  const ap = pickBlock(cfg, 'aliveProbeConfig', 'aliveProbe')
  const wl = pickBlock(cfg, 'websiteLoginConfig', 'websiteLogin')
  const sd = pickBlock(cfg, 'subdomainCollectConfig', 'subdomainCollect')
  const lm = cfg.lateralMove || {}
  const vulIds = Array.isArray(cfg.vulIdsConfig) ? cfg.vulIdsConfig : []

  const basics = buildItems([
    { label: '测试模式', value: getTestModeLabel(cfg.testMode) },
    { label: '漏洞利用', value: yesNo(cfg.vulExploit) }
  ])

  const sections = [
    {
      key: 'port',
      title: '端口扫描',
      icon: 'el-icon-connection',
      enabled: enabled(ps.isOpen),
      items: buildItems([
        { label: '智能端口扫描', value: yesNo(ps.intelligencePort) },
        { label: '端口范围', value: portRangeLabel(ps) },
        { label: '扫描端口', value: portDetail(ps), full: true, mono: true },
        { label: '扫描方式', value: labelFrom(TCP_SCAN_LABELS, ps.tcpScanType, 'tcpScanTypeZh', ps) },
        { label: '超时时间', value: timeoutLabel(ps.timeout, 'timeoutZh', ps) },
        { label: '并发数', value: labelFrom({}, ps.concurrent, 'concurrentZh', ps) || display(ps.concurrent) }
      ])
    },
    {
      key: 'crawler',
      title: 'Web 爬虫',
      icon: 'el-icon-share',
      enabled: enabled(wc.isOpen),
      items: buildItems([
        { label: '爬取深度', value: labelFrom(CRAWLER_DEPTH_LABELS, wc.maxDepth, 'maxDepthZh', wc) },
        { label: '最大 URL 数', value: labelFrom(CRAWLER_MAX_URL_LABELS, wc.maxUrl, 'maxUrlZh', wc) },
        { label: '爬取范围', value: labelFrom(CRAWLER_RANGE_LABELS, wc.scanRange, 'scanRangeZh', wc) },
        { label: '爬取速度', value: labelFrom(CRAWLER_SPEED_LABELS, wc.crawlerSpeed, null, wc) },
        { label: '单页超时', value: timeoutLabel(wc.timeout, 'timeoutZh', wc) },
        { label: '全局超时', value: timeoutLabel(wc.fullTimeout, 'fullTimeoutZh', wc) },
        { label: 'URL 去重', value: labelFrom(CRAWLER_REPEAT_LABELS, wc.scanRepeat, 'scanRepeatZh', wc) },
        { label: '后缀过滤', value: display(wc.suffixFilter) },
        { label: '白名单', value: display(wc.whiteList), full: true },
        { label: '黑名单', value: display(wc.blackList), full: true }
      ])
    },
    {
      key: 'path',
      title: '路径爆破',
      icon: 'el-icon-folder-opened',
      enabled: enabled(wps.isOpen),
      items: buildItems([
        { label: '智能路径爆破', value: yesNo(wps.isIntelligent) },
        { label: '猜测速率', value: labelFrom(PATH_SPEED_LABELS, wps.guessRate, 'guessRateZh', wps) },
        { label: '猜测时长', value: labelFrom(PATH_TIME_LABELS, wps.guessTimeout, 'guessTimeoutZh', wps) },
        { label: '路径字典', value: joinList(wps.scanDict, wps.dickNames), full: true },
        { label: '标题黑名单', value: display(wps.titleBlack), full: true }
      ])
    },
    {
      key: 'weakpass',
      title: '弱口令扫描',
      icon: 'el-icon-key',
      enabled: enabled(wp.isOpen),
      items: buildItems([
        { label: '检测协议', value: joinList(wp.services, wp.servicesZh), full: true },
        { label: '字典类型', value: labelFrom(WEAK_DICT_LABELS, wp.dictType, 'dictTypeZh', wp) },
        { label: '通用用户字典', value: display(wp.commonUserDictZh || wp.commonUserDict) },
        { label: '通用密码字典', value: display(wp.commonPassDictZh || wp.commonPassDict) },
        { label: '补充账号', value: display(wp.addAccount), full: true },
        { label: '补充密码', value: display(wp.addPass), full: true },
        { label: '仅使用补充字典', value: yesNo(wp.onlyUseAdd) },
        { label: '猜测次数', value: labelFrom({}, wp.guessNum, 'guessNumZh', wp) || display(wp.guessNum) },
        { label: '猜测时长', value: labelFrom(PATH_TIME_LABELS, wp.guessTimeout, 'guessTimeoutZh', wp) },
        { label: '猜测速率', value: labelFrom(WEAK_RATE_LABELS, wp.guessRate, 'guessRateZh', wp) },
        { label: '验证码模式', value: display(wp.captchaMode) }
      ])
    },
    {
      key: 'alive',
      title: '存活探测',
      icon: 'el-icon-aim',
      enabled: enabled(ap.isOpen),
      items: buildItems([
        { label: '探测类型', value: labelFrom(ALIVE_PROBE_LABELS, ap.aliveProbeType, 'aliveProbeTypeZh', ap) },
        { label: '探测端口', value: display(ap.probePort), full: true, mono: true }
      ])
    },
    {
      key: 'subdomain',
      title: '子域名收集',
      icon: 'el-icon-link',
      enabled: enabled(sd.isOpen),
      items: buildItems([
        { label: '子域名字典', value: display(sd.subdomainDictZh || sd.subdomainDict) }
      ])
    },
    {
      key: 'proxy',
      title: '代理配置',
      icon: 'el-icon-position',
      enabled: enabled(px.isOpen),
      items: buildItems([
        { label: '代理协议', value: labelFrom(PROXY_PROTO_LABELS, px.proto) },
        { label: '代理地址', value: display(px.addr) },
        { label: '代理端口', value: display(px.port) },
        { label: '认证', value: px.isAuth ? `是（${display(px.username)}）` : '否' }
      ])
    },
    {
      key: 'lateral',
      title: '横向移动',
      icon: 'el-icon-sort',
      enabled: enabled(lm.isOpen),
      items: buildItems([
        { label: '策略', value: LATERAL_STRATEGY_LABELS[lm.strategy] || display(lm.strategy) },
        { label: '扫描范围', value: display(lm.range), full: true },
        { label: '扫描端口', value: display(lm.ports), full: true, mono: true },
        { label: '超时', value: lm.timeout ? `${lm.timeout}s` : '-' }
      ])
    }
  ]

  const loginRows = enabled(wl.isOpen) && Array.isArray(wl.list)
    ? wl.list.map(row => ({
        target: display(row.target),
        verifyType: row.verifyTypeZh || LOGIN_TYPE_LABELS[row.verifyType] || '-',
        verifyValue: display(row.verifyValue),
        verifyStatus: row.verifyStatusZh || '-'
      }))
    : []

  return {
    basics,
    sections,
    vulIdsCount: vulIds.length,
    loginRows,
    loginEnabled: enabled(wl.isOpen),
    crawlerHeaders: enabled(wc.isOpen) && Array.isArray(wc.headers) ? wc.headers : []
  }
}

export function formatScanConfigJson(cfg) {
  if (!cfg) return '{}'
  try {
    return JSON.stringify(cfg, null, 2)
  } catch {
    return String(cfg)
  }
}
