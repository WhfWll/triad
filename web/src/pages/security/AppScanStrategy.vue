<template>
  <div class="security-container">
    <div class="main-title">应用安全 · 扫描策略</div>
    <p class="page-intro">管理扫描策略，定义漏洞脚本选择、爬虫、端口扫描等配置，供新建任务时复用。</p>

    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <span class="db-stat">共 <strong>{{ strategies.length }}</strong> 个策略</span>
          <el-button type="primary" size="small" icon="el-icon-plus" @click="createStrategy" style="margin-left: 12px">新建策略</el-button>
        </div>
      </div>

      <div class="strategy-grid">
        <div v-for="s in strategies" :key="s.id" class="strategy-card" :class="{ builtin: s.builtin }">
          <div class="card-header">
            <div class="card-icon">{{ s.icon }}</div>
            <div class="card-info">
              <div class="card-name">
                {{ s.name }}
                <el-tag v-if="s.builtin" size="mini" type="info">内置</el-tag>
              </div>
              <div class="card-desc">{{ s.desc }}</div>
            </div>
          </div>
          <div class="card-stats">
            <span class="stat-item">漏洞脚本: <strong>{{ s.vulnCount || 0 }}</strong></span>
            <span class="stat-item">爬虫: <strong>{{ s.config?.webCrawler?.isOpen ? '开启' : '关闭' }}</strong></span>
            <span class="stat-item">端口扫描: <strong>{{ s.config?.portScan?.isOpen ? '开启' : '关闭' }}</strong></span>
          </div>
          <div class="card-actions">
            <el-button size="mini" @click="editStrategy(s)">编辑</el-button>
            <el-button size="mini" type="primary" @click="useStrategy(s)">使用</el-button>
            <el-button v-if="!s.builtin" size="mini" type="danger" plain @click="deleteStrategy(s)">删除</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import security from '@/api/security.js'

const DEFAULT_STRATEGIES = [
  {
    id: 'builtin-full',
    name: '全漏洞扫描',
    desc: '覆盖所有已知漏洞类型，启用爬虫与端口扫描',
    icon: '🛡',
    builtin: true,
    vulnCount: 0,
    config: {
      testMode: 'principle',
      safeTest: true,
      vulExploit: false,
      testIntensity: 3,
      vulIdsConfig: [],
      webCrawler: { isOpen: true, maxDepth: 5, scanRange: 0, crawlerSpeed: 2 },
      portScan: { isOpen: true, scanPort: '21,22,23,80,443,445,3306,8000,8080', tcpScanType: 1, timeout: 10, concurrent: 100 },
      proxy: { isOpen: false }
    }
  },
  {
    id: 'builtin-highrisk',
    name: '高危漏洞扫描',
    desc: '仅高危/严重等级的漏洞脚本',
    icon: '🔴',
    builtin: true,
    vulnCount: 0,
    config: {
      testMode: 'principle',
      safeTest: true,
      vulExploit: false,
      testIntensity: 4,
      vulIdsConfig: [],
      webCrawler: { isOpen: true, maxDepth: 3, scanRange: 0, crawlerSpeed: 2 },
      portScan: { isOpen: true, scanPort: '80,443,3306,8080', tcpScanType: 1, timeout: 10, concurrent: 100 },
      proxy: { isOpen: false }
    }
  },
  {
    id: 'builtin-web',
    name: 'Web漏洞扫描',
    desc: '专注Web应用漏洞检测，启用爬虫深入抓取',
    icon: '🌐',
    builtin: true,
    vulnCount: 0,
    config: {
      testMode: 'principle',
      safeTest: true,
      vulExploit: false,
      testIntensity: 3,
      vulIdsConfig: [],
      webCrawler: { isOpen: true, maxDepth: 6, scanRange: 0, crawlerSpeed: 3 },
      portScan: { isOpen: false, scanPort: '', tcpScanType: 1, timeout: 10, concurrent: 100 },
      proxy: { isOpen: false }
    }
  },
  {
    id: 'builtin-weakpass',
    name: '弱口令扫描',
    desc: '检测常见服务的弱口令漏洞',
    icon: '🔑',
    builtin: true,
    vulnCount: 0,
    config: {
      testMode: 'principle',
      safeTest: false,
      vulExploit: false,
      testIntensity: 2,
      vulIdsConfig: [],
      webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
      portScan: { isOpen: true, scanPort: '21,22,23,3306,3389,5432,6379', tcpScanType: 1, timeout: 10, concurrent: 50 },
      proxy: { isOpen: false }
    }
  },
  {
    id: 'builtin-component',
    name: '组件漏洞扫描',
    desc: '扫描第三方组件和中间件已知漏洞',
    icon: '📦',
    builtin: true,
    vulnCount: 0,
    config: {
      testMode: 'version',
      safeTest: true,
      vulExploit: false,
      testIntensity: 3,
      vulIdsConfig: [],
      webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
      portScan: { isOpen: true, scanPort: '80,443,8080,8000,8443', tcpScanType: 1, timeout: 10, concurrent: 100 },
      proxy: { isOpen: false }
    }
  },
  {
    id: 'builtin-portscan',
    name: '端口扫描',
    desc: '快速发现开放端口与运行服务',
    icon: '🔍',
    builtin: true,
    vulnCount: 0,
    config: {
      testMode: 'principle',
      safeTest: true,
      vulExploit: false,
      testIntensity: 1,
      vulIdsConfig: [],
      webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 },
      portScan: { isOpen: true, scanPort: '1-65535', tcpScanType: 2, timeout: 5, concurrent: 200 },
      proxy: { isOpen: false }
    }
  }
]

export default {
  name: 'AppScanStrategy',
  data() {
    return {
      builtinStrategies: DEFAULT_STRATEGIES,
      customStrategies: []
    }
  },
  computed: {
    strategies() {
      return [...this.builtinStrategies, ...this.customStrategies]
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
    this.loadCustomStrategies()
  },
  methods: {
    loadCustomStrategies() {
      try {
        const saved = localStorage.getItem('appsec_strategies')
        if (saved) this.customStrategies = JSON.parse(saved)
      } catch { this.customStrategies = [] }
    },
    saveCustomStrategies() {
      localStorage.setItem('appsec_strategies', JSON.stringify(this.customStrategies))
    },
    createStrategy() {
      this.$router.push('/appsec/strategy/new')
    },
    editStrategy(s) {
      if (s.builtin) {
        this.$router.push({ path: '/appsec/strategy/edit', query: { id: s.id } })
      } else {
        this.$router.push({ path: '/appsec/strategy/edit', query: { id: s.id } })
      }
    },
    useStrategy(s) {
      this.$router.push({ path: '/appsec/task/new', query: { strategyId: s.id, strategyBuiltin: s.builtin ? '1' : '0' } })
    },
    deleteStrategy(s) {
      this.$confirm(`确认删除策略「${s.name}」？`, '提示', { type: 'warning' }).then(() => {
        this.customStrategies = this.customStrategies.filter(x => x.id !== s.id)
        this.saveCustomStrategies()
        this.$message({ message: '已删除', type: 'success' })
      }).catch(() => {})
    }
  }
}
</script>

<style lang="less" scoped>
.page-intro {
  color: #94a3b8;
  font-size: 13px;
  margin: 0 0 16px;
  max-width: 900px;
}

.strategy-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-top: 16px;
}

.strategy-card {
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 20px;
  transition: all 0.25s;
}

.strategy-card:hover {
  border-color: rgba(0, 212, 170, 0.3);
  background: rgba(0, 0, 0, 0.45);
}

.strategy-card.builtin {
  border-left: 3px solid rgba(0, 212, 170, 0.5);
}

.card-header {
  display: flex;
  gap: 14px;
  margin-bottom: 14px;
}

.card-icon {
  font-size: 32px;
  line-height: 1;
  flex-shrink: 0;
}

.card-info {
  flex: 1;
  min-width: 0;
}

.card-name {
  color: #e2e8f0;
  font-weight: 600;
  font-size: 15px;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-desc {
  color: #64748b;
  font-size: 12px;
  line-height: 1.4;
}

.card-stats {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 14px;
  padding: 10px 0;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.stat-item {
  color: #94a3b8;
  font-size: 12px;
}

.stat-item strong {
  color: #cbd5e1;
}

.card-actions {
  display: flex;
  gap: 8px;
}
</style>
