/** 内置策略卡片主题色（列表页 / 选择器共用） */
export const STRATEGY_THEME = {
  'builtin-full': { accent: '#00d4aa', glow: 'rgba(0, 212, 170, 0.35)' },
  'builtin-highrisk': { accent: '#f87171', glow: 'rgba(248, 113, 113, 0.35)' },
  'builtin-web': { accent: '#38bdf8', glow: 'rgba(56, 189, 248, 0.35)' },
  'builtin-weakpass': { accent: '#fbbf24', glow: 'rgba(251, 191, 36, 0.35)' },
  'builtin-component': { accent: '#a78bfa', glow: 'rgba(167, 139, 250, 0.35)' },
  'builtin-portscan': { accent: '#34d399', glow: 'rgba(52, 211, 153, 0.35)' }
}

const DEFAULT_THEME = { accent: '#00d4aa', glow: 'rgba(0, 212, 170, 0.35)' }

export function getStrategyTheme(id) {
  return STRATEGY_THEME[id] || DEFAULT_THEME
}

export function strategyCardAccentStyle(id) {
  const t = getStrategyTheme(id)
  return {
    '--card-accent': t.accent,
    '--card-glow': t.glow
  }
}

export function strategyIconWrapStyle(id) {
  const t = getStrategyTheme(id)
  return {
    background: `linear-gradient(135deg, ${t.accent}22 0%, ${t.accent}08 100%)`,
    borderColor: `${t.accent}44`,
    boxShadow: `0 8px 24px -8px ${t.glow}`
  }
}
