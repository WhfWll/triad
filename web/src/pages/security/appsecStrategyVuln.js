import { vulnerability } from '@/api/tool.js'

/** 与后端 enums.VulLibrariesRisk* 一致 */
export const VUL_RISK_DEAD = 1
export const VUL_RISK_HIGH = 2

/** 与后端 enums.VulLibrariesClass* 一致 */
export const VUL_CLASS_WEB = 1
export const VUL_CLASS_WEB_CMS = 5
export const VUL_CLASS_WEB_FRAMEWORK = 6
export const VUL_CLASS_WEB_COMPONENT = 7

/** 内置策略默认插件筛选规则（vulIdsConfig 为空时按此解析） */
export const STRATEGY_VULN_RULES = {
  'builtin-highrisk': {
    vulRiskLevels: [VUL_RISK_DEAD, VUL_RISK_HIGH],
    label: '致命+高危'
  },
  'builtin-web': {
    vulClassLevels: [VUL_CLASS_WEB, VUL_CLASS_WEB_CMS, VUL_CLASS_WEB_FRAMEWORK, VUL_CLASS_WEB_COMPONENT],
    label: 'Web 相关'
  }
}

export function getStrategyVulnRule(strategyId, config) {
  const base = STRATEGY_VULN_RULES[strategyId] || null
  const riskLevels = (config && config.vulRiskLevels && config.vulRiskLevels.length)
    ? [...config.vulRiskLevels]
    : (base && base.vulRiskLevels ? [...base.vulRiskLevels] : [])
  const classLevels = (config && config.vulClassLevels && config.vulClassLevels.length)
    ? [...config.vulClassLevels]
    : (base && base.vulClassLevels ? [...base.vulClassLevels] : [])
  const label = (base && base.label) || ''
  if (!riskLevels.length && !classLevels.length) return null
  return { vulRiskLevels: riskLevels, vulClassLevels: classLevels, label }
}

export function getStrategyVulnRiskLevels(strategyId, config) {
  const rule = getStrategyVulnRule(strategyId, config)
  return rule && rule.vulRiskLevels.length ? rule.vulRiskLevels : null
}

export function getStrategyVulnRuleLabel(strategyId, config) {
  const rule = getStrategyVulnRule(strategyId, config)
  return (rule && rule.label) || ''
}

async function fetchVulnIdsByQueryParam(paramKey, values, pageSize = 200) {
  const ids = new Set()
  for (const value of values) {
    let page = 1
    let totalPages = 1
    while (page <= totalPages) {
      const params = { page, size: pageSize, [paramKey]: value }
      const res = await vulnerability.getObjectData(params)
      if (!res || res.code !== 200) break
      const list = (res.data && res.data.list) || []
      const total = (res.data && res.data.total) || 0
      totalPages = Math.max(1, Math.ceil(total / pageSize))
      list.forEach(row => {
        if (row && row.id != null) ids.add(row.id)
      })
      page += 1
    }
  }
  return [...ids]
}

export function fetchVulnIdsByRiskLevels(levels, pageSize = 200) {
  return fetchVulnIdsByQueryParam('libRisk', levels, pageSize)
}

export function fetchVulnIdsByClassLevels(levels, pageSize = 200) {
  return fetchVulnIdsByQueryParam('libClass', levels, pageSize)
}

export async function fetchVulnIdsByRule(rule, pageSize = 200) {
  if (!rule) return []
  const idSets = []
  if (rule.vulRiskLevels && rule.vulRiskLevels.length) {
    idSets.push(new Set(await fetchVulnIdsByRiskLevels(rule.vulRiskLevels, pageSize)))
  }
  if (rule.vulClassLevels && rule.vulClassLevels.length) {
    idSets.push(new Set(await fetchVulnIdsByClassLevels(rule.vulClassLevels, pageSize)))
  }
  if (!idSets.length) return []
  if (idSets.length === 1) return [...idSets[0]]
  let result = idSets[0]
  for (let i = 1; i < idSets.length; i += 1) {
    result = new Set([...result].filter(id => idSets[i].has(id)))
  }
  return [...result]
}

export async function resolveStrategyVulnIds(strategyId, config) {
  if (config && Array.isArray(config.vulIdsConfig) && config.vulIdsConfig.length) {
    return [...config.vulIdsConfig]
  }
  const rule = getStrategyVulnRule(strategyId, config)
  if (!rule) return []
  return fetchVulnIdsByRule(rule)
}
