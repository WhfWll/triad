export function kindLabel(kind) {
  return kind === 'sensitive' ? '敏感数据发现' : '数据库安全检查'
}

export const DB_TYPE_REDIS = 4

export function isRedisDbType(type) {
  return Number(type) === DB_TYPE_REDIS
}

export function getDBTypeName(type) {
  const map = { 1: 'MySQL', 2: 'PostgreSQL', 3: 'MongoDB', 4: 'Redis', 5: 'CouchDB' }
  return map[type] || '未知'
}

export function getDBTypeClass(type) {
  const map = { 1: 'db-mysql', 2: 'db-postgresql', 3: 'db-mongodb', 4: 'db-redis', 5: 'db-couchdb' }
  return map[type] || 'db-default'
}

export function getRiskName(risk) {
  const map = { 0: '严重', 1: '高危', 2: '中危', 3: '低危', 4: '信息' }
  return map[risk] || '未知'
}

export function getRiskClass(risk) {
  const map = { 0: 'risk-critical', 1: 'risk-high', 2: 'risk-medium', 3: 'risk-low', 4: 'risk-info' }
  return map[risk] || 'risk-default'
}

export function getStatusName(status) {
  const map = { 1: '等待中', 2: '进行中', 3: '已完成' }
  return map[status] || '未知'
}

export function getStatusClass(status) {
  const map = { 1: 'status-wait', 2: 'status-running', 3: 'status-complete' }
  return map[status] || 'status-default'
}

export function getCategoryName(category) {
  const map = {
    1: '身份认证', 2: '权限控制', 3: '配置安全', 4: '审计日志',
    5: '网络安全', 6: '加密', 7: 'SQL 注入', 8: '敏感数据识别'
  }
  return map[category] || '未知'
}

export function getDataTypeName(type) {
  const map = {
    1: '身份证号', 2: '银行卡号', 3: '护照号', 4: '手机号', 5: '邮箱', 6: '地址',
    7: '出生日期', 8: '姓名', 9: 'Token', 10: '证书信息', 11: '密码哈希'
  }
  return map[type] || '未知'
}

export function getSensitivityName(level) {
  const map = { 1: '高敏感', 2: '中敏感', 3: '低敏感' }
  return map[level] || '未知'
}

export function getSensitivityClass(level) {
  const map = { 1: 'sensitivity-high', 2: 'sensitivity-medium', 3: 'sensitivity-low' }
  return map[level] || 'sensitivity-default'
}
