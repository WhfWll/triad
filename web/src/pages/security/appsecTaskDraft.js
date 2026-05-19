import {
  getBuiltinStrategy,
  cloneBuiltinConfig,
  isBuiltinStrategyId,
  loadCustomStrategy
} from './appsecBuiltinStrategies.js'

const DRAFT_KEY = 'appsec_task_draft'

export function saveTaskDraft(draft) {
  sessionStorage.setItem(DRAFT_KEY, JSON.stringify(draft))
}

export function loadTaskDraft() {
  try {
    const raw = sessionStorage.getItem(DRAFT_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

export function clearTaskDraft() {
  sessionStorage.removeItem(DRAFT_KEY)
}

export function draftRouteQuery(draft) {
  if (!draft) return {}
  return {
    strategy: draft.strategyId,
    type: draft.scanType === 'app' ? 'app' : 'dyn',
    ...(draft.isCustom ? { custom: '1' } : {})
  }
}

export function initStrategyFromRoute(strategyId, isCustomQuery) {
  if (!strategyId) return null

  if (!isCustomQuery && isBuiltinStrategyId(strategyId)) {
    const builtin = getBuiltinStrategy(strategyId)
    if (!builtin) return null
    return {
      isCustom: false,
      strategyMeta: { ...builtin },
      scanConfig: cloneBuiltinConfig(strategyId)
    }
  }

  const custom = loadCustomStrategy(strategyId)
  if (custom) {
    return {
      isCustom: true,
      strategyMeta: {
        id: custom.id,
        name: custom.name,
        desc: custom.desc || '',
        icon: custom.icon || '⚙',
        baseStrategyId: custom.baseStrategyId || 'builtin-full'
      },
      scanConfig: JSON.parse(JSON.stringify(custom.config || cloneBuiltinConfig('builtin-full')))
    }
  }

  if (isBuiltinStrategyId(strategyId)) {
    const builtin = getBuiltinStrategy(strategyId)
    if (!builtin) return null
    return {
      isCustom: false,
      strategyMeta: { ...builtin },
      scanConfig: cloneBuiltinConfig(strategyId)
    }
  }

  return null
}
