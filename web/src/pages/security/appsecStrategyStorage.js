const CUSTOM_KEY = 'appsec_strategies'
const OVERRIDE_KEY = 'appsec_builtin_overrides'

export function loadCustomStrategies() {
  try {
    const raw = localStorage.getItem(CUSTOM_KEY)
    return raw ? JSON.parse(raw) : []
  } catch {
    return []
  }
}

export function saveCustomStrategies(list) {
  localStorage.setItem(CUSTOM_KEY, JSON.stringify(list))
}

export function loadCustomStrategy(id) {
  return loadCustomStrategies().find(x => x.id === id) || null
}

export function loadBuiltinOverrides() {
  try {
    const raw = localStorage.getItem(OVERRIDE_KEY)
    return raw ? JSON.parse(raw) : {}
  } catch {
    return {}
  }
}

export function loadBuiltinOverride(id) {
  const all = loadBuiltinOverrides()
  return all[id] || null
}

export function saveBuiltinOverride(id, payload) {
  const all = loadBuiltinOverrides()
  all[id] = {
    ...payload,
    updatedAt: Date.now()
  }
  localStorage.setItem(OVERRIDE_KEY, JSON.stringify(all))
}

export function mergeScanConfig(base, patch) {
  const out = JSON.parse(JSON.stringify(base || {}))
  if (!patch) return out
  if (patch.vulIdsConfig) {
    out.vulIdsConfig = [...patch.vulIdsConfig]
  }
  ;['testMode', 'safeTest', 'vulExploit', 'testIntensity'].forEach(k => {
    if (patch[k] !== undefined) out[k] = patch[k]
  })
  if (patch.webCrawler) out.webCrawler = { ...out.webCrawler, ...patch.webCrawler }
  if (patch.portScan) out.portScan = { ...out.portScan, ...patch.portScan }
  if (patch.proxy) out.proxy = { ...out.proxy, ...patch.proxy }
  return out
}
