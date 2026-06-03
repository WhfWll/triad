/** 弱口令扫描配置默认值（与后端 enums.WeakPassConfig 字段对齐） */

export const DEFAULT_WEAK_PASS = {
  isOpen: false,
  services: [],
  dictType: 1,
  commonUserDict: 0,
  commonPassDict: 0,
  addAccount: '',
  addPass: '',
  onlyUseAdd: false,
  guessNum: 3,
  guessTimeout: 300,
  guessRate: 3,
  captchaMode: 'common_arithmetic'
}

/** 弱口令策略模板默认启用的爆破服务（不含「通用」） */
export const WEAKPASS_DEFAULT_SERVICE_IDS = [2, 3, 5, 6, 7, 8, 9, 10, 11, 12]

/** 场景枚举接口不可用时的服务列表兜底 */
export const WEAKPASS_SERVICE_FALLBACK = [
  { value: 2, label: 'ssh' },
  { value: 3, label: 'ftp' },
  { value: 4, label: 'memcached' },
  { value: 5, label: 'mongodb' },
  { value: 6, label: 'mssql' },
  { value: 7, label: 'mysql' },
  { value: 8, label: 'postgres' },
  { value: 9, label: 'rdp' },
  { value: 10, label: 'redis' },
  { value: 11, label: 'smb' },
  { value: 12, label: 'telnet' },
  { value: 13, label: 'tomcat' },
  { value: 14, label: 'vnc' },
  { value: 22, label: 'http' }
]

export function cloneDefaultWeakPass(overrides = {}) {
  const services = overrides.services != null ? [...overrides.services] : [...DEFAULT_WEAK_PASS.services]
  return {
    ...DEFAULT_WEAK_PASS,
    ...overrides,
    services
  }
}
