/** 应用安全任务列表/详情共用文案与工具 */

export const RISK_NAMES = { 0: '严重', 1: '高危', 2: '中危', 3: '低危', 4: '信息' }
export const RISK_CLASSES = {
  0: 'risk-critical',
  1: 'risk-high',
  2: 'risk-medium',
  3: 'risk-low',
  4: 'risk-info'
}

export const DYN_STATUS_NAMES = { 1: '等待扫描', 2: '扫描中', 3: '已完成' }
export const APP_STATUS_NAMES = { 1: '等待检测', 2: '检测中', 3: '已完成' }
export const STATUS_CLASSES = {
  1: 'status-wait',
  2: 'status-running',
  3: 'status-complete'
}

export const DYN_VULN_TYPES = {
  1: 'SQL注入',
  2: 'XSS',
  3: 'SSRF',
  4: 'XXE',
  5: '命令注入',
  6: '文件包含',
  7: '文件上传',
  8: 'CSRF',
  9: '信息泄露'
}

export const APP_VULN_TYPES = {
  1: '远程代码执行',
  2: '未授权访问',
  3: 'SQL注入',
  4: '文件上传',
  5: '弱口令',
  6: 'XSS',
  7: 'SSRF',
  8: '信息泄露'
}

export const APP_TYPE_NAMES = {
  1: '万户 OA',
  2: '用友 NC',
  3: '蓝凌 EKP',
  4: '云时空',
  5: '亿赛通',
  6: 'D-Link',
  7: '通达 OA',
  8: 'WordPress',
  9: 'ThinkPHP',
  10: 'Spring Boot',
  11: '通用 CMS',
  12: '泛微 OA',
  13: '致远 OA',
  14: '金蝶云星空',
  15: 'Apache Struts'
}

export function getRiskName(risk) {
  return RISK_NAMES[risk] || '未知'
}

export function getRiskClass(risk) {
  return RISK_CLASSES[risk] || 'risk-default'
}

export function getStatusName(status, scanType) {
  const map = scanType === 'app' ? APP_STATUS_NAMES : DYN_STATUS_NAMES
  return map[status] || '未知'
}

export function getStatusClass(status) {
  return STATUS_CLASSES[status] || 'status-default'
}

export function getVulnTypeName(type, scanType) {
  const map = scanType === 'app' ? APP_VULN_TYPES : DYN_VULN_TYPES
  return map[type] || '未知'
}

export function getAppTypeName(type) {
  return APP_TYPE_NAMES[type] || '未知'
}

export function buildSiteMapTree(pages) {
  const tree = []
  const map = {}
  ;(pages || []).forEach(page => {
    const url = (page && page.url) || page || ''
    if (!url) return
    const parts = String(url)
      .replace(/^https?:\/\//, '')
      .split('/')
      .filter(Boolean)
    let current = tree
    parts.forEach((part, index) => {
      const path = parts.slice(0, index + 1).join('/')
      if (!map[path]) {
        const node = { name: part, path, children: [] }
        map[path] = node
        current.push(node)
      }
      current = map[path].children
    })
  })
  return tree
}

export function scanTypeLabel(scanType) {
  return scanType === 'app' ? '专项应用检测' : '动态扫描'
}

export const TEST_MODE_NAMES = {
  principle: '原理验证',
  version: '版本匹配',
  '1': '原理验证',
  '2': '版本匹配'
}

export function getTestModeLabel(mode) {
  if (!mode) return '-'
  const key = String(mode).toLowerCase()
  if (TEST_MODE_NAMES[key]) return TEST_MODE_NAMES[key]
  if (key.includes('principle') && key.includes('version')) return '原理验证 + 版本匹配'
  if (key.includes('principle')) return '原理验证'
  if (key.includes('version')) return '版本匹配'
  return mode
}
