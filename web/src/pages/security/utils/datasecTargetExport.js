export function exportTargetsJson(targets, dbType) {
  const payload = {
    version: 1,
    dbType: Number(dbType) || 1,
    targets: (targets || []).map((t) => ({
      dbHost: t.dbHost || '',
      dbPort: Number(t.dbPort) || 0,
      dbName: t.dbName || '',
      dbUser: t.dbUser || '',
      dbPassword: t.dbPassword || ''
    }))
  }
  const blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `datasec-targets-${Date.now()}.json`
  a.click()
  URL.revokeObjectURL(url)
}

export function parseTargetsImportFile(text) {
  const data = JSON.parse(text)
  let targets = []
  let dbType = 1
  if (Array.isArray(data)) {
    targets = data
  } else if (data && Array.isArray(data.targets)) {
    targets = data.targets
    dbType = data.dbType || dbType
  } else if (data && Array.isArray(data.items)) {
    targets = data.items.map((item) => ({
      dbHost: item.dbHost,
      dbPort: item.dbPort,
      dbName: item.dbName,
      dbUser: item.dbUser,
      dbPassword: item.dbPassword
    }))
    dbType = data.dbType || dbType
  } else {
    throw new Error('无法识别的 JSON 格式')
  }
  return {
    dbType,
    targets: targets.map((t) => ({
      dbHost: (t.dbHost || '').trim(),
      dbPort: Number(t.dbPort) || 0,
      dbName: (t.dbName || '').trim(),
      dbUser: (t.dbUser || '').trim(),
      dbPassword: t.dbPassword || ''
    })).filter((t) => t.dbHost && t.dbUser)
  }
}

export function pickImportFile(accept = '.json') {
  return new Promise((resolve, reject) => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = accept
    input.onchange = () => {
      const file = input.files && input.files[0]
      if (!file) {
        reject(new Error('未选择文件'))
        return
      }
      const reader = new FileReader()
      reader.onload = () => resolve(String(reader.result || ''))
      reader.onerror = () => reject(new Error('读取文件失败'))
      reader.readAsText(file)
    }
    input.click()
  })
}
