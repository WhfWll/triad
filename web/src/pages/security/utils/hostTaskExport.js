function escapeCsvCell(val) {
  const s = val == null ? '' : String(val)
  if (/[",\n\r]/.test(s)) {
    return `"${s.replace(/"/g, '""')}"`
  }
  return s
}

function downloadBlob(filename, content, type) {
  const blob = new Blob([content], { type })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = filename
  a.click()
  URL.revokeObjectURL(a.href)
}

export function exportHostItemsCsv(taskId, mode, items) {
  let headers = []
  let rows = []

  if (mode === 'baseline') {
    headers = ['目标主机', '分类', '检查项', '结果', '风险', '期望值', '实际值', '检查时间', '修复建议']
    rows = (items || []).map(item => [
      item.targetIp || '',
      item.categoryName || '',
      item.ruleName || '',
      item.resultName || '',
      item.riskName || '',
      item.expectedValue || '',
      item.actualValue || '',
      item.checkTime || '',
      item.fixSuggestion || ''
    ])
  } else if (mode === 'vuln') {
    headers = ['目标主机', 'CVE ID', '严重程度', '风险', '影响软件包', '版本', '漏洞标题', '发现时间']
    rows = (items || []).map(item => [
      item.targetIp || '',
      item.cveId || '',
      item.severity || '',
      item.riskName || '',
      item.packageName || '',
      item.packageVersion || '',
      item.title || '',
      item.checkTime || ''
    ])
  } else {
    headers = ['目标主机', '检测类型', '匹配规则', '风险', '文件路径', '进程信息', '发现时间', '描述', '修复建议']
    rows = (items || []).map(item => [
      item.targetIp || '',
      item.checkTypeName || '',
      item.matchRule || '',
      item.riskName || '',
      item.filePath || '',
      item.processInfo || '',
      item.checkTime || '',
      item.description || '',
      item.fixSuggestion || ''
    ])
  }

  const lines = [headers, ...rows].map(cols => cols.map(escapeCsvCell).join(','))
  downloadBlob(`hostsec-${taskId || 'task'}-${mode}-items.csv`, '\uFEFF' + lines.join('\n'), 'text/csv;charset=utf-8')
}

export function buildHostSummaryText(meta) {
  const {
    taskId,
    kindLabel,
    mode,
    checkTime,
    targetCount,
    itemCount,
    statData,
    vulnStat,
    malwareStat
  } = meta

  const lines = [
    '主机安全检查审查摘要',
    '========================',
    `任务批次: ${taskId || '-'}`,
    `任务类型: ${kindLabel || '-'}`,
    `检测模式: ${mode || '-'}`,
    `检测目标数: ${targetCount || 0}`,
    `结果条数: ${itemCount || 0}`,
    `执行时间: ${checkTime || '-'}`,
    `导出时间: ${new Date().toLocaleString()}`
  ]

  if (mode === 'baseline') {
    lines.push(
      `检查项总数: ${statData.totalRules || 0}`,
      `通过: ${statData.passCount || 0}`,
      `不通过: ${statData.failCount || 0}`,
      `异常: ${statData.errorCount || 0}`,
      `跳过: ${statData.skipCount || 0}`,
      `合规通过率: ${statData.effectivePassRate || 0}%`,
      `整体通过率: ${statData.passRate || 0}%`
    )
  } else if (mode === 'vuln') {
    lines.push(
      `扫描包数: ${vulnStat.packages || 0}`,
      `漏洞总数: ${vulnStat.matchedVulns || 0}`,
      `严重: ${vulnStat.critical || 0}`,
      `高危: ${vulnStat.high || 0}`,
      `中危: ${vulnStat.medium || 0}`,
      `低危: ${vulnStat.low || 0}`
    )
  } else {
    lines.push(
      `发现项总数: ${malwareStat.totalFindings || 0}`,
      `严重: ${malwareStat.critical || 0}`,
      `高危: ${malwareStat.high || 0}`,
      `中危: ${malwareStat.medium || 0}`,
      `低危: ${malwareStat.low || 0}`
    )
  }

  return lines.join('\n')
}

export function downloadTextFile(filename, content) {
  downloadBlob(filename, content, 'text/plain;charset=utf-8')
}
