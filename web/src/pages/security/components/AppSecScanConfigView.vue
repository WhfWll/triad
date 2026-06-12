<template>
  <div class="scan-config-view">
    <div class="config-panel is-enabled">
      <div class="config-panel-head">
        <div class="panel-title"><i class="el-icon-setting"></i> 基础策略</div>
      </div>
      <div class="cfg-kv-grid cfg-kv-grid--basics">
        <div
          v-for="item in basicItems"
          :key="item.label"
          class="cfg-kv-card"
        >
          <span class="cfg-kv-label">{{ item.label }}</span>
          <span class="cfg-kv-value" :class="{ mono: item.mono }">{{ item.value }}</span>
        </div>
      </div>
    </div>

    <div class="cfg-cap-strip">
      <div class="cfg-cap-strip-title">扫描能力开关</div>
      <div class="config-capability-row">
        <span
          v-for="cap in capabilityChips"
          :key="cap.key"
          class="cfg-cap-chip"
          :class="{ on: cap.on }"
        >
          <i :class="cap.icon"></i> {{ cap.label }}
        </span>
      </div>
    </div>

    <div
      v-for="section in view.sections"
      :key="section.key"
      class="config-panel"
      :class="section.enabled ? 'is-enabled' : 'is-disabled'"
    >
      <div class="config-panel-head">
        <div class="panel-title">
          <i :class="section.icon"></i> {{ section.title }}
        </div>
        <span class="cfg-status-tag" :class="section.enabled ? 'on' : 'off'">
          {{ section.enabled ? '已启用' : '未启用' }}
        </span>
      </div>
      <template v-if="section.enabled">
        <div class="cfg-kv-grid">
          <div
            v-for="item in section.items"
            :key="item.label"
            class="cfg-kv-card"
            :class="{ 'full-width': item.full }"
          >
            <span class="cfg-kv-label">{{ item.label }}</span>
            <div v-if="shouldShowPortTags(item)" class="cfg-port-tags">
              <span v-for="tag in portTags(item.value)" :key="tag" class="cfg-port-tag">{{ tag }}</span>
              <span v-if="portTagsOverflow(item.value)" class="cfg-port-more">
                +{{ portTagsOverflow(item.value) }} 个
              </span>
            </div>
            <span v-else class="cfg-kv-value" :class="{ mono: item.mono }">{{ item.value }}</span>
          </div>
        </div>
      </template>
      <p v-else class="config-disabled-hint">未启用此扫描能力</p>
    </div>

    <div class="config-panel" :class="view.loginEnabled ? 'is-enabled' : 'is-disabled'">
      <div class="config-panel-head">
        <div class="panel-title"><i class="el-icon-user"></i> 网站登录凭证</div>
        <span class="cfg-status-tag" :class="view.loginEnabled ? 'on' : 'off'">
          {{ view.loginEnabled ? '已启用' : '未启用' }}
        </span>
      </div>
      <div v-if="view.loginEnabled && view.loginRows.length" class="config-table-wrap">
        <el-table :data="view.loginRows" class="myTable" size="small">
          <el-table-column prop="target" label="登录地址" min-width="200" :show-overflow-tooltip="true" />
          <el-table-column prop="verifyType" label="凭证类型" width="100" />
          <el-table-column prop="verifyValue" label="凭证内容" min-width="180" :show-overflow-tooltip="true" />
          <el-table-column prop="verifyStatus" label="状态" width="100" />
        </el-table>
      </div>
      <p v-else-if="view.loginEnabled" class="config-disabled-hint">已开启登录凭证，但未配置条目</p>
      <p v-else class="config-disabled-hint">未启用网站登录凭证</p>
    </div>

    <div class="config-panel is-enabled">
      <div class="config-panel-head">
        <div class="panel-title"><i class="el-icon-cpu"></i> 漏洞检测插件</div>
        <span class="cfg-count-badge">{{ vulIdsLabel }}</span>
      </div>
      <p class="config-vuln-hint">
        任务执行时将按所选插件库规则进行漏洞检测；插件 ID 列表过长时不逐项展开，仅展示数量。
      </p>
    </div>

    <div v-if="view.crawlerHeaders.length" class="config-panel is-enabled">
      <div class="config-panel-head">
        <div class="panel-title"><i class="el-icon-s-operation"></i> 爬虫请求头</div>
      </div>
      <div class="config-table-wrap">
        <el-table :data="view.crawlerHeaders" class="myTable" size="small">
          <el-table-column prop="key" label="Header" min-width="160" />
          <el-table-column prop="value" label="Value" min-width="240" :show-overflow-tooltip="true" />
        </el-table>
      </div>
    </div>

    <div class="config-raw-wrap">
      <button type="button" class="config-raw-toggle" @click="showRaw = !showRaw">
        <i :class="showRaw ? 'el-icon-arrow-up' : 'el-icon-arrow-down'"></i>
        {{ showRaw ? '收起原始 JSON' : '查看原始配置 JSON' }}
      </button>
      <pre v-if="showRaw" class="code-block config-json">{{ rawJson }}</pre>
    </div>
  </div>
</template>

<script>
import { buildScanConfigSections, formatScanConfigJson } from '../utils/appsecScanConfigView.js'

const PORT_TAG_LIMIT = 18

export default {
  name: 'AppSecScanConfigView',
  props: {
    config: { type: Object, default: null },
    strategyId: { type: String, default: '' }
  },
  data() {
    return {
      showRaw: false
    }
  },
  computed: {
    view() {
      return buildScanConfigSections(this.config)
    },
    basicItems() {
      const items = [...(this.view.basics || [])]
      if (this.strategyId) {
        items.push({ label: '扫描策略', value: this.strategyId, mono: true })
      }
      return items
    },
    rawJson() {
      return formatScanConfigJson(this.config)
    },
    vulIdsLabel() {
      const n = this.view.vulIdsCount
      return n > 0 ? `已选 ${n} 个` : '未指定（按策略默认）'
    },
    capabilityChips() {
      const map = {}
      ;(this.view.sections || []).forEach(s => {
        map[s.key] = s.enabled
      })
      return [
        { key: 'port', label: '端口扫描', icon: 'el-icon-connection', on: map.port },
        { key: 'crawler', label: 'Web 爬虫', icon: 'el-icon-share', on: map.crawler },
        { key: 'path', label: '路径爆破', icon: 'el-icon-folder-opened', on: map.path },
        { key: 'weakpass', label: '弱口令', icon: 'el-icon-key', on: map.weakpass },
        { key: 'proxy', label: '代理', icon: 'el-icon-position', on: map.proxy }
      ]
    }
  },
  methods: {
    shouldShowPortTags(item) {
      if (!item || !item.mono || !item.value || item.value === '-') return false
      return String(item.value).includes(',')
    },
    portTags(value) {
      return String(value || '')
        .split(',')
        .map(s => s.trim())
        .filter(Boolean)
        .slice(0, PORT_TAG_LIMIT)
    },
    portTagsOverflow(value) {
      const all = String(value || '')
        .split(',')
        .map(s => s.trim())
        .filter(Boolean)
      return all.length > PORT_TAG_LIMIT ? all.length - PORT_TAG_LIMIT : 0
    }
  }
}
</script>

<style lang="less" scoped>
@import '../css/appsec-scan-config-view.less';
</style>
