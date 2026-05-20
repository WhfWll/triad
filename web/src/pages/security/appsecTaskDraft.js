import {
  getBuiltinStrategy,
  cloneBuiltinConfig,
  isBuiltinStrategyId,
  loadCustomStrategy
} from './appsecBuiltinStrategies.js'
import { loadBuiltinOverride, mergeScanConfig } from './appsecStrategyStorage.js'

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
    const override = loadBuiltinOverride(strategyId)
    let config = cloneBuiltinConfig(strategyId)
    if (override && override.config) {
      config = mergeScanConfig(config, override.config)
    }
    return {
      isCustom: false,
      strategyMeta: {
        ...builtin,
        name: (override && override.name) || builtin.name,
        desc: (override && override.desc) || builtin.desc,
        icon: (override && override.icon) || builtin.icon
      },
      scanConfig: config
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
    const override = loadBuiltinOverride(strategyId)
    let config = cloneBuiltinConfig(strategyId)
    if (override && override.config) {
      config = mergeScanConfig(config, override.config)
    }
    return {
      isCustom: false,
      strategyMeta: {
        ...builtin,
        name: (override && override.name) || builtin.name,
        desc: (override && override.desc) || builtin.desc,
        icon: (override && override.icon) || builtin.icon
      },
      scanConfig: config
    }
  }

  return null
}
