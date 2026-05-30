/** 应用安全任务结果导出 */

function escapeCsvCell(val) {
  const s = val == null ? '' : String(val)
  if (/[",\n\r]/.test(s)) {
    return `"${s.replace(/"/g, '""')}"`
  }
  return s
}

export function exportVulnsToCsv(task, scanType, { getRiskName }) {
  const vulns = (task && task.vulns) || []
  const headers = ['漏洞名称', '类型', '风险等级', 'URL', '描述', '修复建议']
  const rows = vulns.map(v => [
    v.name,
    v.typeName || '-',
    getRiskName(v.riskLevel),
    v.url || '',
    (v.description || '').replace(/\s+/g, ' ').slice(0, 500),
    (v.suggestion || '').replace(/\s+/g, ' ').slice(0, 500)
  ])
  const lines = [headers, ...rows].map(cols => cols.map(escapeCsvCell).join(','))
  const bom = '\uFEFF'
  const blob = new Blob([bom + lines.join('\n')], { type: 'text/csv;charset=utf-8' })
  const name = `appsec-${task.id || 'task'}-vulns.csv`
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = name
  a.click()
  URL.revokeObjectURL(a.href)
}

export function buildAuditSummaryText(task, scanType, labels) {
  const {
    getStatusName,
    getRiskName,
    getAppTypeName,
    scanTypeLabel
  } = labels
  const lines = [
    '应用安全扫描审查摘要',
    '========================',
    `任务名称: ${task.name || '-'}`,
    `任务 ID: ${task.id || '-'}`,
    `扫描类型: ${scanTypeLabel(scanType)}`,
    `扫描目标: ${task.targetSummary || task.targetUrl || '-'}`,
    ...(task.targetCount > 1 ? [`目标数量: ${task.targetCount}`] : []),
    `状态: ${getStatusName(task.status, scanType)}`,
    `整体风险: ${getRiskName(task.riskLevel)}`,
    `创建时间: ${task.createTime || '-'}`,
    `扫描时间: ${task.scanTime || '-'}`
  ]
  if (scanType === 'app' && task.appType) {
    lines.push(`应用类型: ${getAppTypeName(task.appType)}`)
  }
  if (task.strategyId) {
    lines.push(`扫描策略: ${task.strategyId}`)
  }
  if (scanType === 'dyn') {
    lines.push(`爬取页面: ${task.pageCount || 0}`)
  }
  lines.push(
    `漏洞统计: 严重 ${task.criticalCount || 0} / 高危 ${task.highRiskCount || 0} / 中危 ${task.middleRiskCount || 0} / 低危 ${task.lowRiskCount || 0} / 合计 ${task.vulnCount || 0}`,
    `导出时间: ${new Date().toLocaleString()}`
  )
  return lines.join('\n')
}

export function downloadTextFile(filename, content) {
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = filename
  a.click()
  URL.revokeObjectURL(a.href)
}
