<template>
  <div class="security-container">
    <div class="main-title">{{ isDynamic ? '新建扫描任务' : '新建专项检测任务' }}</div>
    <p class="page-intro">选择扫描策略并指定目标，快速创建安全扫描任务。</p>

    <div class="list_box">
      <el-form :model="form" ref="form" label-width="110px" class="task-form">

        <el-divider content-position="left">基础配置</el-divider>

        <el-form-item label="任务名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入任务名称" maxlength="100" />
        </el-form-item>

        <el-form-item label="目标地址" prop="targetUrl">
          <el-input v-model="form.targetUrl" placeholder="支持 URL (http://example.com)、IP (192.168.1.1)、IP段 (192.168.1.1/24)" />
        </el-form-item>

        <template v-if="!isDynamic">
          <el-form-item label="应用类型" prop="appType">
            <el-select v-model="form.appType" placeholder="请选择应用类型" style="width: 300px">
              <el-option label="万户 OA" :value="1" />
              <el-option label="用友 NC" :value="2" />
              <el-option label="蓝凌 EKP" :value="3" />
              <el-option label="云时空" :value="4" />
              <el-option label="亿赛通" :value="5" />
              <el-option label="D-Link" :value="6" />
              <el-option label="通达 OA" :value="7" />
              <el-option label="WordPress" :value="8" />
              <el-option label="ThinkPHP" :value="9" />
              <el-option label="Spring Boot" :value="10" />
              <el-option label="通用 CMS" :value="11" />
            </el-select>
          </el-form-item>
        </template>

        <el-divider content-position="left">扫描策略</el-divider>

        <div class="strategy-cards">
          <div
            v-for="s in strategies"
            :key="s.id"
            class="strategy-card"
            :class="{ active: selectedStrategy === s.id }"
            @click="selectedStrategy = s.id"
          >
            <div class="strategy-icon">{{ s.icon }}</div>
            <div class="strategy-info">
              <div class="strategy-name">{{ s.name }}</div>
              <div class="strategy-desc">{{ s.desc }}</div>
            </div>
          </div>
          <div class="strategy-card strategy-more" @click="goStrategyMgmt">
            <div class="strategy-more-icon">+</div>
            <div class="strategy-info">
              <div class="strategy-name">管理策略</div>
              <div class="strategy-desc">新建或编辑扫描策略</div>
            </div>
          </div>
        </div>

      </el-form>

      <div class="form-actions">
        <el-button @click="goBack">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">创建任务</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import security from '@/api/security.js'

const DEFAULT_STRATEGIES = [
  { id: 'builtin-full', name: '全漏洞扫描', desc: '启用所有脚本 + 爬虫 + 端口', icon: '🛡' },
  { id: 'builtin-highrisk', name: '高危漏洞扫描', desc: '仅高危/严重等级漏洞', icon: '🔴' },
  { id: 'builtin-web', name: 'Web漏洞扫描', desc: 'Web漏洞 + 爬虫深入抓取', icon: '🌐' },
  { id: 'builtin-weakpass', name: '弱口令扫描', desc: '常见服务弱口令检测', icon: '🔑' },
  { id: 'builtin-component', name: '组件漏洞扫描', desc: '第三方组件已知漏洞', icon: '📦' },
  { id: 'builtin-portscan', name: '端口扫描', desc: '发现开放端口与运行服务', icon: '🔍' }
]

export default {
  name: 'AppScanTaskNew',
  data() {
    return {
      submitting: false,
      strategies: DEFAULT_STRATEGIES,
      selectedStrategy: 'builtin-full',
      form: {
        name: '',
        targetUrl: '',
        appType: 1
      }
    }
  },
  computed: {
    isDynamic() {
      return this.$route.query.type !== 'app'
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
    this.loadStrategies()
  },
  methods: {
    loadStrategies() {
      try {
        const raw = localStorage.getItem('appsec_strategies')
        if (raw) {
          const custom = JSON.parse(raw)
          this.strategies = [...DEFAULT_STRATEGIES, ...custom.map(s => ({ id: s.id, name: s.name, desc: s.desc, icon: s.icon || '⚙' }))]
        }
      } catch {}
      const sid = this.$route.query.strategyId
      if (sid && this.strategies.find(s => s.id === sid)) {
        this.selectedStrategy = sid
      }
    },
    getStrategyConfig(id) {
      const builtinMap = {
        'builtin-full': { testMode: 'principle', safeTest: true, vulExploit: false, testIntensity: 3, vulIdsConfig: [], webCrawler: { isOpen: true, maxDepth: 5, scanRange: 0, crawlerSpeed: 2 }, portScan: { isOpen: true, scanPort: '21,22,23,80,443,445,3306,8000,8080', tcpScanType: 1, timeout: 10, concurrent: 100 }, proxy: { isOpen: false } },
        'builtin-highrisk': { testMode: 'principle', safeTest: true, vulExploit: false, testIntensity: 4, vulIdsConfig: [], webCrawler: { isOpen: true, maxDepth: 3, scanRange: 0, crawlerSpeed: 2 }, portScan: { isOpen: true, scanPort: '80,443,3306,8080', tcpScanType: 1, timeout: 10, concurrent: 100 }, proxy: { isOpen: false } },
        'builtin-web': { testMode: 'principle', safeTest: true, vulExploit: false, testIntensity: 3, vulIdsConfig: [], webCrawler: { isOpen: true, maxDepth: 6, scanRange: 0, crawlerSpeed: 3 }, portScan: { isOpen: false, scanPort: '', tcpScanType: 1, timeout: 10, concurrent: 100 }, proxy: { isOpen: false } },
        'builtin-weakpass': { testMode: 'principle', safeTest: false, vulExploit: false, testIntensity: 2, vulIdsConfig: [], webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 }, portScan: { isOpen: true, scanPort: '21,22,23,3306,3389,5432,6379', tcpScanType: 1, timeout: 10, concurrent: 50 }, proxy: { isOpen: false } },
        'builtin-component': { testMode: 'version', safeTest: true, vulExploit: false, testIntensity: 3, vulIdsConfig: [], webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 }, portScan: { isOpen: true, scanPort: '80,443,8080,8000,8443', tcpScanType: 1, timeout: 10, concurrent: 100 }, proxy: { isOpen: false } },
        'builtin-portscan': { testMode: 'principle', safeTest: true, vulExploit: false, testIntensity: 1, vulIdsConfig: [], webCrawler: { isOpen: false, maxDepth: 1, scanRange: 0, crawlerSpeed: 1 }, portScan: { isOpen: true, scanPort: '1-65535', tcpScanType: 2, timeout: 5, concurrent: 200 }, proxy: { isOpen: false } }
      }
      if (builtinMap[id]) return builtinMap[id]
      try {
        const raw = localStorage.getItem('appsec_strategies')
        if (raw) {
          const list = JSON.parse(raw)
          const s = list.find(x => x.id === id)
          if (s && s.config) return s.config
        }
      } catch {}
      return builtinMap['builtin-full']
    },
    goStrategyMgmt() {
      this.$router.push('/appsec/strategy')
    },
    goBack() {
      this.$router.push({ path: '/appsec/tasks', query: { tab: this.isDynamic ? 'dyn' : 'app' } })
    },
    async submitForm() {
      if (!this.form.name) { this.$message({ message: '请输入任务名称', type: 'warning' }); return }
      if (!this.form.targetUrl) { this.$message({ message: '请输入目标地址', type: 'warning' }); return }
      this.submitting = true
      try {
        const config = this.getStrategyConfig(this.selectedStrategy)
        const api = this.isDynamic ? security.runDynamicScan : security.runAppSpecificScan
        const res = await api({
          name: this.form.name,
          target: this.form.targetUrl,
          appType: this.form.appType,
          strategy: this.selectedStrategy,
          ...config
        })
        if (res.code == 200) {
          this.$message({ message: '任务创建成功', type: 'success' })
          this.goBack()
        } else {
          this.$message({ message: res.msg || '创建失败', type: 'error' })
        }
      } catch {
        this.$message({ message: '创建失败，请稍后重试', type: 'error' })
      } finally {
        this.submitting = false
      }
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
.task-form { max-width: 860px; padding: 10px 0; }
.form-actions {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  text-align: center;
}
.form-actions .el-button { min-width: 120px; }

.strategy-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}
.strategy-card {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 14px;
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.25s;
}
.strategy-card:hover {
  border-color: rgba(0, 212, 170, 0.4);
  background: rgba(0, 0, 0, 0.5);
}
.strategy-card.active {
  border-color: #00d4aa;
  background: rgba(0, 212, 170, 0.08);
  box-shadow: 0 0 12px rgba(0, 212, 170, 0.15);
}
.strategy-icon { font-size: 24px; line-height: 1; flex-shrink: 0; }
.strategy-info { flex: 1; min-width: 0; }
.strategy-name { color: #e2e8f0; font-weight: 600; font-size: 13px; margin-bottom: 2px; }
.strategy-desc { color: #64748b; font-size: 11px; line-height: 1.3; }
.strategy-more {
  justify-content: center;
  align-items: center;
  border-style: dashed;
}
.strategy-more-icon {
  font-size: 28px;
  color: #64748b;
  line-height: 1;
}
/deep/ .el-divider__text { color: #00d4aa; font-weight: 600; background: transparent; }
/deep/ .el-divider { background-color: rgba(255, 255, 255, 0.08); }
</style>
