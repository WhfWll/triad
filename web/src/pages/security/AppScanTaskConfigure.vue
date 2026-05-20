<template>
  <div class="security-container nessus-editor">
    <!-- 顶栏：面包屑 + 返回 -->
    <div class="editor-header">
      <div class="breadcrumb">
        <span class="crumb-muted">新建扫描</span>
        <span class="crumb-sep">/</span>
        <span class="crumb-current">{{ strategyMeta ? strategyMeta.name : '扫描任务' }}</span>
      </div>
      <a class="back-link" @click.prevent="goPrev">
        <i class="el-icon-arrow-left" /> 返回扫描策略
      </a>
    </div>

    <!-- 顶栏 Tab：配置 | 插件 -->
    <div class="editor-tabs">
      <button
        type="button"
        class="tab-btn"
        :class="{ active: activeTab === 'settings' }"
        @click="switchTab('settings')"
      >
        配置
      </button>
      <button
        v-if="needsVulnStep"
        type="button"
        class="tab-btn"
        :class="{ active: activeTab === 'plugins' }"
        @click="switchTab('plugins')"
      >
        插件
        <span v-if="pluginCount > 0" class="tab-badge">{{ pluginCount }}</span>
      </button>
    </div>

    <div class="editor-body list_box">
      <!-- ========== 配置 Tab ========== -->
      <div v-if="activeTab === 'settings'" class="settings-layout">
        <aside class="settings-sidebar">
          <button
            v-for="item in sidebarItems"
            :key="item.key"
            type="button"
            class="nav-item"
            :class="{ active: activeSection === item.key }"
            @click="activeSection = item.key"
          >
            {{ item.label }}
          </button>
        </aside>

        <main class="settings-panel">
          <h2 class="panel-title">{{ panelTitle }}</h2>
          <p v-if="panelHint" class="panel-hint">{{ panelHint }}</p>

          <el-form
            ref="form"
            :model="form"
            label-position="top"
            class="nessus-form"
          >
            <!-- 基础配置 -->
            <template v-if="activeSection === 'general'">
              <el-form-item label="名称" required>
                <el-input v-model="form.name" placeholder="请输入任务名称" maxlength="100" />
              </el-form-item>
              <el-form-item label="目标" required>
                <el-input
                  v-model="form.targetUrl"
                  type="textarea"
                  :rows="4"
                  placeholder="支持 URL（http://example.com）、IP（192.168.1.1）、网段（192.168.1.0/24）&#10;每行一个目标，或使用逗号分隔"
                />
                <p class="field-hint">可填写多个扫描目标，每行一个或使用逗号分隔。</p>
              </el-form-item>
              <el-form-item v-if="!isDynamic" label="应用类型" required>
                <el-select v-model="form.appType" placeholder="请选择应用类型" style="width: 100%; max-width: 400px">
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
                <p class="field-hint">专项检测将按应用类型加载对应 PoC 规则。</p>
              </el-form-item>
            </template>

            <!-- 扫描 / 发现 / 高级：复用配置表单 -->
            <scan-strategy-config-form
              v-else-if="scanConfig && configSectionKey"
              :config="scanConfig"
              :strategy-id="effectiveStrategyId"
              :section="configSectionKey"
              page-mode="all"
            />
          </el-form>
        </main>
      </div>

      <!-- ========== 插件 Tab ========== -->
      <div v-if="activeTab === 'plugins'" class="plugins-layout">
        <h2 class="panel-title">漏洞检测插件</h2>
        <p class="panel-hint">选择本次任务要执行的漏洞脚本（对应后端 vulIdsConfig）。未选择时将按策略默认规则执行。</p>
        <scan-strategy-config-form
          v-if="scanConfig"
          ref="pluginsForm"
          :key="'plugins-' + effectiveStrategyId"
          :config="scanConfig"
          :strategy-id="effectiveStrategyId"
          section="vuln"
          page-mode="vuln"
        />
      </div>

      <div class="editor-footer">
        <div class="editor-footer-inner">
          <el-button @click="goPrev">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="submitTask">保存并启动扫描</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import security from '@/api/security.js'
import ScanStrategyConfigForm from './components/ScanStrategyConfigForm.vue'
import { getStrategySections } from './appsecBuiltinStrategies.js'
import {
  saveTaskDraft,
  loadTaskDraft,
  clearTaskDraft,
  initStrategyFromRoute
} from './appsecTaskDraft.js'
import { applyScanAssessmentDefaults } from './appsecBuiltinStrategies.js'

const PANEL_META = {
  general: { title: '基础配置', hint: '任务名称、扫描目标及应用类型（专项检测）。' },
  port: { title: '端口扫描', hint: '配置扫描端口范围、方式与超时。' },
  crawler: { title: '爬虫配置', hint: '配置 Web 爬取深度、范围与速度。' },
  advanced: { title: '代理设置', hint: '通过 HTTP/HTTPS/SOCKS5 代理发起扫描。' },
  login: { title: '登录凭证', hint: '为需登录的 Web 目标配置 Cookie 或 Header。' }
}

export default {
  name: 'AppScanTaskConfigure',
  components: { ScanStrategyConfigForm },
  data() {
    return {
      submitting: false,
      strategyMeta: null,
      isCustom: false,
      scanConfig: null,
      activeTab: 'settings',
      activeSection: 'general',
      form: {
        name: '',
        targetUrl: '',
        appType: 1
      }
    }
  },
  computed: {
    strategyId() {
      return this.$route.query.strategy || ''
    },
    effectiveStrategyId() {
      if (this.isCustom && this.strategyMeta && this.strategyMeta.baseStrategyId) {
        return this.strategyMeta.baseStrategyId
      }
      return this.strategyId
    },
    isDynamic() {
      return this.$route.query.type !== 'app'
    },
    scanType() {
      return this.isDynamic ? 'dyn' : 'app'
    },
    strategySections() {
      return getStrategySections(this.effectiveStrategyId)
    },
    needsVulnStep() {
      return this.strategySections.vuln
    },
    pluginCount() {
      return (this.scanConfig && this.scanConfig.vulIdsConfig && this.scanConfig.vulIdsConfig.length) || 0
    },
    sidebarItems() {
      const s = this.strategySections
      const items = [{ key: 'general', label: '基础配置' }]
      if (s.port) items.push({ key: 'port', label: '端口扫描' })
      if (s.crawler) items.push({ key: 'crawler', label: '爬虫配置' })
      if (s.advanced) items.push({ key: 'advanced', label: '代理设置' })
      if (s.login) items.push({ key: 'login', label: '登录凭证' })
      return items
    },
    configSectionKey() {
      const map = { port: 'port', crawler: 'crawler', advanced: 'advanced', login: 'login' }
      return map[this.activeSection] || ''
    },
    panelTitle() {
      if (this.activeTab === 'plugins') return '漏洞检测插件'
      return (PANEL_META[this.activeSection] && PANEL_META[this.activeSection].title) || '配置'
    },
    panelHint() {
      if (this.activeTab === 'plugins') return ''
      return (PANEL_META[this.activeSection] && PANEL_META[this.activeSection].hint) || ''
    }
  },
  watch: {
    '$route.query.tab': {
      immediate: true,
      handler(tab) {
        if (tab === 'plugins' && this.needsVulnStep) {
          this.activeTab = 'plugins'
        } else {
          this.activeTab = 'settings'
        }
      }
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
    if (!this.initPage()) {
      this.$message.warning('请先选择扫描策略')
      this.goPrev()
    }
  },
  methods: {
    initPage() {
      const id = this.strategyId
      if (!id) return false

      const init = initStrategyFromRoute(id, this.$route.query.custom === '1')
      if (!init) return false

      this.isCustom = init.isCustom
      this.strategyMeta = init.strategyMeta
      this.scanConfig = init.scanConfig

      const draft = loadTaskDraft()
      if (draft && draft.strategyId === id) {
        this.form = { ...this.form, ...draft.form }
        this.scanConfig = JSON.parse(JSON.stringify(draft.scanConfig))
      }
      return true
    },
    switchTab(tab) {
      if (tab === 'plugins' && !this.needsVulnStep) return
      this.activeTab = tab
      const q = { ...this.$route.query }
      if (tab === 'plugins') {
        q.tab = 'plugins'
      } else {
        delete q.tab
      }
      this.$router.replace({ path: this.$route.path, query: q }).catch(() => {})
      if (this.scanConfig) this.persistDraft()
    },
    buildDraft() {
      return {
        strategyId: this.strategyId,
        scanType: this.scanType,
        isCustom: this.isCustom,
        strategyMeta: this.strategyMeta,
        form: { ...this.form },
        scanConfig: JSON.parse(JSON.stringify(this.scanConfig))
      }
    },
    persistDraft() {
      if (this.scanConfig) saveTaskDraft(this.buildDraft())
    },
    countTargets(text) {
      if (!text) return 0
      const parts = String(text)
        .split(/[\n,，;；]+/)
        .map(s => s.trim())
        .filter(Boolean)
      return new Set(parts).size
    },
    validateForm() {
      if (!this.form.name) {
        this.$message({ message: '请输入任务名称', type: 'warning' })
        this.activeTab = 'settings'
        this.activeSection = 'general'
        this.syncTabRoute('settings')
        return false
      }
      if (!this.form.targetUrl || !String(this.form.targetUrl).trim()) {
        this.$message({ message: '请输入扫描目标', type: 'warning' })
        this.activeTab = 'settings'
        this.activeSection = 'general'
        this.syncTabRoute('settings')
        return false
      }
      return true
    },
    syncTabRoute(tab) {
      const q = { ...this.$route.query }
      if (tab === 'plugins') {
        q.tab = 'plugins'
      } else {
        delete q.tab
      }
      this.$router.replace({ path: this.$route.path, query: q }).catch(() => {})
    },
    goPrev() {
      clearTaskDraft()
      this.$router.push({ path: '/appsec/task/new', query: { type: this.scanType } })
    },
    goBack() {
      clearTaskDraft()
      this.$router.push({ path: '/appsec/tasks', query: { tab: this.scanType } })
    },
    syncPluginSelectionFromForm() {
      const form = this.$refs.pluginsForm
      if (form && typeof form.getSelectedVulnIds === 'function') {
        const ids = form.getSelectedVulnIds()
        if (this.scanConfig) {
          this.$set(this.scanConfig, 'vulIdsConfig', ids)
        }
      }
    },
    async submitTask() {
      if (!this.validateForm()) return
      if (this.needsVulnStep) {
        this.syncPluginSelectionFromForm()
      }
      this.persistDraft()
      const draft = this.buildDraft()
      if (draft.scanConfig) applyScanAssessmentDefaults(draft.scanConfig)
      this.submitting = true
      try {
        const api = this.isDynamic ? security.runDynamicScan : security.runAppSpecificScan
        const res = await api({
          name: draft.form.name,
          target: draft.form.targetUrl,
          appType: draft.form.appType,
          strategy: draft.strategyId,
          ...draft.scanConfig
        })
        if (res.code == 200) {
          clearTaskDraft()
          const n = this.countTargets(draft.form.targetUrl)
          const msg = n > 1 ? `任务已创建，共 ${n} 个扫描目标` : '任务创建成功'
          this.$message({ message: msg, type: 'success' })
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
@import '../bas/css/bas-list-page.less';
@import './css/appsec-strategy-editor.less';

.field-hint {
  font-size: 12px;
  color: #64748b;
  margin: 8px 0 0;
  line-height: 1.4;
}
</style>
