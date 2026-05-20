<template>
  <div class="security-container nessus-editor">
    <!-- 新建：先选模板 -->
    <template v-if="step === 'pick'">
      <div class="editor-header">
        <div class="breadcrumb">
          <span class="crumb-muted">扫描策略</span>
          <span class="crumb-sep">/</span>
          <span class="crumb-current">新建策略</span>
        </div>
        <a class="back-link" @click.prevent="goBack">
          <i class="el-icon-arrow-left" /> 返回策略列表
        </a>
      </div>

      <div class="list_box pick-box">
        <app-sec-strategy-picker
          :strategies="pickStrategies"
          intro="选择内置模板作为起点，进入下一步后可修改名称、插件与扫描参数。"
          @select="onPickTemplate"
        />
        <div class="pick-footer">
          <el-button size="small" @click="goBack">取消</el-button>
        </div>
      </div>
    </template>

    <!-- 编辑 / 新建第二步 -->
    <template v-else>
      <div class="editor-header">
        <div class="breadcrumb">
          <span class="crumb-muted">扫描策略</span>
          <span class="crumb-sep">/</span>
          <span class="crumb-current">{{ isNew ? '新建策略' : form.name || '编辑策略' }}</span>
        </div>
        <a class="back-link" @click.prevent="goBackOrPick">
          <i class="el-icon-arrow-left" /> {{ isNew ? '重新选择模板' : '返回策略列表' }}
        </a>
      </div>

      <div v-if="templateMeta" class="selected-banner list_box">
        <div class="banner-icon-wrap" :style="iconWrapStyle(baseStrategyId)">
          <span>{{ templateMeta.icon }}</span>
        </div>
        <div class="banner-text">
          <span class="banner-label">基于模板</span>
          <strong>{{ templateMeta.name }}</strong>
          <span class="banner-desc">{{ templateMeta.desc }}</span>
        </div>
      </div>

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
          v-if="strategySections.vuln"
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

            <el-form label-position="top" class="nessus-form">
              <template v-if="activeSection === 'general'">
                <el-form-item label="策略名称" required>
                  <el-input v-model="form.name" placeholder="如：日常 Web 巡检策略" maxlength="50" />
                </el-form-item>
                <el-form-item label="策略描述">
                  <el-input
                    v-model="form.desc"
                    type="textarea"
                    :rows="3"
                    placeholder="简要说明该策略的适用场景"
                    maxlength="200"
                  />
                </el-form-item>
              </template>

              <scan-strategy-config-form
                v-else-if="scanConfig && configSectionKey"
                :config="scanConfig"
                :strategy-id="baseStrategyId"
                :section="configSectionKey"
                page-mode="all"
              />
            </el-form>
          </main>
        </div>

        <div v-if="activeTab === 'plugins'" class="plugins-layout">
          <h2 class="panel-title">漏洞检测插件</h2>
          <p class="panel-hint">选择该策略默认启用的漏洞脚本（对应 vulIdsConfig）。未选择时按「默认全部」执行。</p>
          <scan-strategy-config-form
            v-if="scanConfig"
            ref="pluginsForm"
            :key="pluginsFormKey"
            :config="scanConfig"
            :strategy-id="baseStrategyId"
            section="vuln"
            page-mode="vuln"
          />
        </div>

        <div class="editor-footer">
          <div class="editor-footer-inner">
            <el-button @click="goBackOrPick">取消</el-button>
            <el-button type="primary" :loading="saving" @click="saveStrategy">保存策略</el-button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import ScanStrategyConfigForm from './components/ScanStrategyConfigForm.vue'
import AppSecStrategyPicker from './components/AppSecStrategyPicker.vue'
import { strategyIconWrapStyle } from './appsecStrategyTheme.js'
import {
  BUILTIN_STRATEGIES,
  getBuiltinStrategy,
  cloneBuiltinConfig,
  getStrategySections,
  isBuiltinStrategyId
} from './appsecBuiltinStrategies.js'
import {
  loadCustomStrategies,
  saveCustomStrategies,
  loadCustomStrategy,
  loadBuiltinOverride,
  saveBuiltinOverride,
  mergeScanConfig
} from './appsecStrategyStorage.js'
import { applyScanAssessmentDefaults } from './appsecBuiltinStrategies.js'

const PANEL_META = {
  general: { title: '基本信息', hint: '策略名称与说明，便于在列表中识别。' },
  port: { title: '端口扫描', hint: '配置扫描端口范围、方式与超时。' },
  login: { title: '登录凭证', hint: '为需登录的 Web 目标配置 Cookie 或 Header，扫描时自动携带。' },
  crawler: { title: '爬虫配置', hint: '配置 Web 爬取深度、范围与速度。' },
  advanced: { title: '代理设置', hint: '通过 HTTP/HTTPS/SOCKS5 代理发起扫描。' }
}

export default {
  name: 'AppScanStrategyEdit',
  components: { ScanStrategyConfigForm, AppSecStrategyPicker },
  data() {
    return {
      step: 'edit',
      saving: false,
      activeTab: 'settings',
      activeSection: 'general',
      baseStrategyId: 'builtin-full',
      templateMeta: null,
      scanConfig: null,
      form: {
        name: '',
        desc: '',
        icon: '🛡'
      }
    }
  },
  computed: {
    isNew() {
      const id = this.$route.query.id
      return !id || id === 'new'
    },
    pickStrategies() {
      return BUILTIN_STRATEGIES.map(s => ({ ...s, builtin: true }))
    },
    strategySections() {
      return getStrategySections(this.baseStrategyId)
    },
    pluginCount() {
      return (this.scanConfig && this.scanConfig.vulIdsConfig && this.scanConfig.vulIdsConfig.length) || 0
    },
    sidebarItems() {
      const s = this.strategySections
      const items = [{ key: 'general', label: '基本信息' }]
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
    },
    pluginsFormKey() {
      const id = this.$route.query.id || 'new'
      const n = (this.scanConfig && this.scanConfig.vulIdsConfig && this.scanConfig.vulIdsConfig.length) || 0
      return `plugins-${id}-${n}`
    }
  },
  watch: {
    '$route.query.tab': {
      immediate: true,
      handler(tab) {
        if (this.step !== 'edit') return
        if (tab === 'plugins' && this.strategySections.vuln) {
          this.activeTab = 'plugins'
        } else if (tab !== 'plugins') {
          this.activeTab = 'settings'
        }
      }
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
    if (this.isNew) {
      const tpl = this.$route.query.template
      const meta = tpl && getBuiltinStrategy(tpl)
      if (meta) {
        this.onPickTemplate({ ...meta, id: tpl, builtin: true })
        return
      }
      this.step = 'pick'
      return
    }
    this.step = 'edit'
    this.loadStrategy()
  },
  methods: {
    iconWrapStyle: strategyIconWrapStyle,
    onPickTemplate(s) {
      this.baseStrategyId = s.id
      this.templateMeta = { name: s.name, desc: s.desc, icon: s.icon }
      this.form.name = ''
      this.form.desc = s.desc || ''
      this.form.icon = s.icon || '⚙'
      this.scanConfig = cloneBuiltinConfig(s.id)
      this.activeTab = 'settings'
      this.activeSection = 'general'
      this.step = 'edit'
      this.$router.replace({ path: '/appsec/strategy/new', query: { template: s.id } }).catch(() => {})
    },
    loadStrategy() {
      const id = this.$route.query.id
      if (!id || id === 'new') return

      const custom = loadCustomStrategy(id)
      if (custom) {
        this.baseStrategyId = custom.baseStrategyId || custom.templateId || 'builtin-full'
        this.form.name = custom.name
        this.form.desc = custom.desc || ''
        this.form.icon = custom.icon || '⚙'
        this.scanConfig = JSON.parse(JSON.stringify(custom.config || cloneBuiltinConfig(this.baseStrategyId)))
        this.templateMeta = getBuiltinStrategy(this.baseStrategyId) || {
          name: custom.name,
          desc: custom.desc,
          icon: custom.icon
        }
        return
      }

      if (isBuiltinStrategyId(id)) {
        const meta = getBuiltinStrategy(id)
        if (!meta) return
        const override = loadBuiltinOverride(id)
        this.baseStrategyId = id
        this.form.name = (override && override.name) || meta.name
        this.form.desc = (override && override.desc) || meta.desc
        this.form.icon = (override && override.icon) || meta.icon
        let config = cloneBuiltinConfig(id)
        if (override && override.config) {
          config = mergeScanConfig(config, override.config)
        }
        this.scanConfig = config
        this.templateMeta = meta
      }
    },
    syncPluginSelectionFromForm() {
      const form = this.$refs.pluginsForm
      if (form && typeof form.getSelectedVulnIds === 'function') {
        const ids = form.getSelectedVulnIds()
        if (!this.scanConfig) return
        this.$set(this.scanConfig, 'vulIdsConfig', ids)
      }
    },
    switchTab(tab) {
      if (tab === 'plugins' && !this.strategySections.vuln) return
      this.activeTab = tab
      const q = { ...this.$route.query }
      if (tab === 'plugins') {
        q.tab = 'plugins'
      } else {
        delete q.tab
      }
      this.$router.replace({ path: this.$route.path, query: q }).catch(() => {})
    },
    goBack() {
      this.$router.push('/appsec/strategy')
    },
    goBackOrPick() {
      if (this.isNew && this.step === 'edit') {
        this.step = 'pick'
        this.$router.replace({ path: '/appsec/strategy/new' }).catch(() => {})
        return
      }
      this.goBack()
    },
    async saveStrategy() {
      if (!this.form.name || !String(this.form.name).trim()) {
        this.$message({ message: '请输入策略名称', type: 'warning' })
        this.activeTab = 'settings'
        this.activeSection = 'general'
        return
      }
      if (!this.scanConfig) {
        this.$message({ message: '策略配置未加载', type: 'error' })
        return
      }

      this.syncPluginSelectionFromForm()
      applyScanAssessmentDefaults(this.scanConfig)

      this.saving = true
      try {
        const routeId = this.$route.query.id
        const configSnapshot = JSON.parse(JSON.stringify(this.scanConfig))
        const vulnCount = (configSnapshot.vulIdsConfig && configSnapshot.vulIdsConfig.length) || 0

        if (isBuiltinStrategyId(routeId)) {
          saveBuiltinOverride(routeId, {
            name: this.form.name.trim(),
            desc: (this.form.desc || '').trim(),
            icon: this.form.icon || '⚙',
            config: configSnapshot
          })
          this.$message({ message: '内置策略已保存（含插件选择）', type: 'success' })
          this.goBack()
          return
        }

        let strategies = loadCustomStrategies()
        const saveAsNew = this.isNew
        const id = saveAsNew ? `custom-${Date.now()}` : routeId
        const existing = strategies.findIndex(x => x.id === id)

        const data = {
          id,
          name: this.form.name.trim(),
          desc: (this.form.desc || '').trim(),
          icon: this.form.icon || '⚙',
          builtin: false,
          baseStrategyId: this.baseStrategyId,
          vulnCount,
          config: configSnapshot
        }

        if (existing >= 0) {
          strategies[existing] = data
        } else {
          strategies.push(data)
        }

        saveCustomStrategies(strategies)
        this.$message({ message: '策略保存成功', type: 'success' })
        if (saveAsNew) {
          this.$router.replace({ path: '/appsec/strategy/edit', query: { id } }).catch(() => {})
        } else {
          this.goBack()
        }
      } catch {
        this.$message({ message: '保存失败', type: 'error' })
      } finally {
        this.saving = false
      }
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';
@import './css/appsec-strategy-editor.less';

.pick-box {
  padding: 24px;
  margin-top: 16px;
}

.pick-footer {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  justify-content: flex-end;
}

.selected-banner {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 16px;
  padding: 14px 20px;
  border: 1px solid rgba(0, 212, 170, 0.15);
  background: rgba(0, 212, 170, 0.04);
}

.banner-icon-wrap {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid;
  border-radius: 12px;
  font-size: 22px;
  flex-shrink: 0;
}

.banner-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;

  strong {
    color: #e2e8f0;
    font-size: 15px;
  }
}

.banner-label {
  font-size: 11px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.banner-desc {
  font-size: 12px;
  color: #94a3b8;
}
</style>
