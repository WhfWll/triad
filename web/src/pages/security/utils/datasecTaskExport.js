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

export function exportDatasecRowsCsv(taskId, exportType, rows, labels) {
  let headers = []
  let dataRows = []

  if (exportType === 'checks') {
    headers = ['检查类别', '规则名称', '风险等级', '检查结果', '期望值', '实际值', '描述', '修复建议']
    dataRows = (rows || []).map(row => [
      labels.getCategoryName(row.category),
      row.ruleName || '',
      labels.getRiskName(row.riskLevel),
      row.result || '',
      row.expectedValue || '',
      row.actualValue || '',
      row.description || '',
      row.suggestion || ''
    ])
  } else if (exportType === 'cve') {
    headers = ['CVE / 规则', '风险等级', '结果', '版本信息', '描述', '修复建议']
    dataRows = (rows || []).map(row => [
      row.ruleName || '',
      labels.getRiskName(row.riskLevel),
      row.result || '',
      row.actualValue || '',
      row.description || '',
      row.suggestion || ''
    ])
  } else {
    headers = ['表名', '字段名', '数据类型', '敏感等级', '样例数据', '数量']
    dataRows = (rows || []).map(row => [
      row.tableName || '',
      row.columnName || '',
      labels.getDataTypeName(row.dataType),
      labels.getSensitivityName(row.sensitivityLevel),
      row.sampleData || '',
      row.count || 0
    ])
  }

  const lines = [headers, ...dataRows].map(cols => cols.map(escapeCsvCell).join(','))
  downloadBlob(`datasec-${taskId || 'task'}-${exportType}.csv`, '\uFEFF' + lines.join('\n'), 'text/csv;charset=utf-8')
}

export function buildDatasecSummaryText(meta, labels) {
  const {
    task,
    kind,
    targetCount,
    timeLabel,
    scanSummary,
    cveFailCount,
    findingCount
  } = meta

  const lines = [
    '数据安全检查审查摘要',
    '========================',
    `任务名称: ${(task && task.name) || '-'}`,
    `任务 ID: ${(task && task.id) || '-'}`,
    `任务类型: ${kind === 'db' ? '数据库安全检查' : '敏感数据扫描'}`,
    `数据库类型: ${labels.getDBTypeName((task && task.dbType) || 0)}`,
    `扫描目标数: ${targetCount || 0}`,
    `${kind === 'db' ? '检查时间' : '扫描时间'}: ${timeLabel || '-'}`,
    `导出时间: ${new Date().toLocaleString()}`
  ]

  if (kind === 'db') {
    lines.push(
      `基线检查: 共 ${scanSummary.baselineTotal || 0} 项 / 通过 ${scanSummary.baselinePass || 0} / 不通过 ${scanSummary.baselineFail || 0} / 异常 ${scanSummary.baselineError || 0}`,
      `CVE 命中: ${cveFailCount || 0} 条`,
      `风险统计: 严重 ${(task && task.criticalCount) || 0} / 高危 ${(task && task.highRiskCount) || 0} / 中危 ${(task && task.middleRiskCount) || 0} / 低危 ${(task && task.lowRiskCount) || 0}`
    )
  }

  lines.push(
    `敏感字段: ${findingCount || 0} 条`,
    `敏感等级: 高 ${(task && task.highCount) || 0} / 中 ${(task && task.mediumCount) || 0} / 低 ${(task && task.lowCount) || 0}`
  )

  return lines.join('\n')
}

export function downloadTextFile(filename, content) {
  downloadBlob(filename, content, 'text/plain;charset=utf-8')
}
