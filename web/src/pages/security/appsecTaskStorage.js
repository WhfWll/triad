/** 应用安全扫描任务本地存储（后端接口就绪后可替换为 API） */

const DYN_KEY = 'appsec_dynamic_tasks'
const APP_KEY = 'appsec_app_tasks'

/** 避免重复调度 */
const scheduledKeys = new Set()

function storageKey(type) {
  return type === 'app' ? APP_KEY : DYN_KEY
}

function formatTime(d = new Date()) {
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function loadAll(type) {
  try {
    const raw = localStorage.getItem(storageKey(type))
    return raw ? JSON.parse(raw) : []
  } catch {
    return []
  }
}

function saveAll(type, list) {
  localStorage.setItem(storageKey(type), JSON.stringify(list))
}

function simKey(type, taskId) {
  return `${type}:${taskId}`
}

function completeTaskSimulation(type, taskId) {
  const list = loadAll(type)
  const idx = list.findIndex(t => t.id === taskId)
  if (idx < 0) return
  const task = list[idx]
  if (task.status === 3) return

  const configured = (task.vulnCount && task.scanConfig && task.scanConfig.vulIdsConfig)
    ? task.scanConfig.vulIdsConfig.length
    : task.vulnCount || 0
  const found = configured > 0 ? Math.min(configured, 2 + Math.floor(Math.random() * 4)) : Math.floor(Math.random() * 2)

  list[idx] = {
    ...task,
    status: 3,
    pageCount: 8 + Math.floor(Math.random() * 40),
    vulnCount: found,
    criticalCount: found > 0 ? 1 : 0,
    highRiskCount: found > 1 ? 1 : 0,
    middleRiskCount: found > 2 ? 1 : 0,
    lowRiskCount: Math.max(0, found - 3),
    riskLevel: found > 0 ? (found >= 3 ? 0 : 1) : 4,
    scanTime: formatTime(),
    vulns: [],
    pages: []
  }
  saveAll(type, list)
  scheduledKeys.delete(simKey(type, taskId))
}

/** 模拟扫描过程（无真实引擎时的演示逻辑） */
export function scheduleTaskSimulation(type, taskId) {
  const key = simKey(type, taskId)
  if (scheduledKeys.has(key)) return
  scheduledKeys.add(key)
  const delay = 4000 + Math.floor(Math.random() * 4000)
  setTimeout(() => completeTaskSimulation(type, taskId), delay)
}

/** 将卡在「等待」的旧任务拉起为扫描中 */
function bootstrapPendingTasks(type) {
  const list = loadAll(type)
  let changed = false
  list.forEach(t => {
    if (t.status === 1) {
      t.status = 2
      if (!t.scanTime || t.scanTime === '-') t.scanTime = formatTime()
      changed = true
      scheduleTaskSimulation(type, t.id)
    }
  })
  if (changed) saveAll(type, list)
}

/**
 * @param {'dyn'|'app'} type
 * @param {object} payload 提交参数
 */
export function createAppSecTask(type, payload) {
  const list = loadAll(type)
  const now = formatTime()
  const vulnIds = (payload.vulIdsConfig && payload.vulIdsConfig.length) || 0
  const task = {
    id: `task-${type}-${Date.now()}`,
    name: payload.name || '未命名任务',
    targetUrl: payload.target || payload.targetUrl || '',
    strategyId: payload.strategy || '',
    appType: payload.appType,
    scanConfig: payload,
    status: 2,
    riskLevel: 4,
    pageCount: 0,
    vulnCount: vulnIds,
    criticalCount: 0,
    highRiskCount: 0,
    middleRiskCount: 0,
    lowRiskCount: 0,
    createTime: now,
    scanTime: now,
    vulns: [],
    pages: []
  }
  list.unshift(task)
  saveAll(type, list)
  scheduleTaskSimulation(type, task.id)
  return task
}

export function listAppSecTasks(type, { page = 1, size = 10, search = '' } = {}) {
  bootstrapPendingTasks(type)
  let list = loadAll(type)
  const q = (search || '').trim().toLowerCase()
  if (q) {
    list = list.filter(
      t =>
        (t.name && t.name.toLowerCase().includes(q)) ||
        (t.targetUrl && t.targetUrl.toLowerCase().includes(q))
    )
  }
  const total = list.length
  const start = (Math.max(1, page) - 1) * size
  return {
    list: list.slice(start, start + size),
    total
  }
}

export function getAppSecTask(type, id) {
  return loadAll(type).find(t => t.id === id) || null
}

export function hasRunningAppSecTasks(type) {
  bootstrapPendingTasks(type)
  return loadAll(type).some(t => t.status === 2)
}
