<template>
  <div class="security-container task-detail-page">
    <div class="detail-topbar list_box">
      <el-link :underline="false" class="link-back" @click="goBack">
        <i class="el-icon-arrow-left"></i> 返回
      </el-link>
      <div class="topbar-main" v-if="task">
        <h1 class="task-title">{{ task.name }}</h1>
        <p class="task-meta">
          <span>{{ typeLabel }}</span>
          <span class="dot">·</span>
          <span>{{ task.targetSummary || task.targetUrl }}</span>
          <span v-if="task.targetCount > 1" class="dot">·</span>
          <span v-if="task.targetCount > 1">{{ task.targetCount }} 个目标</span>
          <span v-if="scanType === 'app' && task.appType" class="dot">·</span>
          <span v-if="scanType === 'app' && task.appType">{{ getAppTypeName(task.appType) }}</span>
        </p>
        <div class="task-tags">
          <span :class="getStatusClass(task.status)">{{ getStatusName(task.status, scanType) }}</span>
          <span :class="getRiskClass(task.riskLevel)">整体 {{ getRiskName(task.riskLevel) }}</span>
        </div>
      </div>
      <div class="topbar-actions">
        <el-button size="small" :loading="loading" @click="loadDetail">刷新</el-button>
        <el-button size="small" :disabled="!taskId" @click="generateReport">生成报告</el-button>
        <el-dropdown trigger="click" @command="onExportCommand">
          <el-button type="primary" size="small" :disabled="!task">
            导出 <i class="el-icon-arrow-down el-icon--right"></i>
          </el-button>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item command="csv">漏洞列表 CSV</el-dropdown-item>
            <el-dropdown-item command="summary">审查摘要 TXT</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </div>

    <div v-loading="loading" class="detail-body">
      <template v-if="task">
        <el-tabs v-model="activeTab" class="detail-tabs list_box" @tab-click="onTabClick">
          <!-- 概览 -->
          <el-tab-pane label="概览" name="overview">
            <div class="stat-row">
              <div v-if="scanType === 'dyn'" class="stat-card pages">
                <span class="stat-label">爬取页面</span>
                <span class="stat-value">{{ task.pageCount || 0 }}</span>
              </div>
              <div class="stat-card critical">
                <span class="stat-label">严重</span>
                <span class="stat-value">{{ task.criticalCount || 0 }}</span>
              </div>
              <div class="stat-card high">
                <span class="stat-label">高危</span>
                <span class="stat-value">{{ task.highRiskCount || 0 }}</span>
              </div>
              <div class="stat-card medium">
                <span class="stat-label">中危</span>
                <span class="stat-value">{{ task.middleRiskCount || 0 }}</span>
              </div>
              <div class="stat-card low">
                <span class="stat-label">低危</span>
                <span class="stat-value">{{ task.lowRiskCount || 0 }}</span>
              </div>
              <div class="stat-card total">
                <span class="stat-label">漏洞合计</span>
                <span class="stat-value">{{ task.vulnCount || (task.vulns && task.vulns.length) || 0 }}</span>
              </div>
            </div>
            <div class="info-panel inner-panel">
              <div class="panel-title">任务信息</div>
              <div class="info-grid">
                <div class="info-item">
                  <span class="label">任务 ID</span>
                  <span class="value mono">{{ task.id }}</span>
                </div>
                <div class="info-item">
                  <span class="label">扫描类型</span>
                  <span class="value">{{ typeLabel }}</span>
                </div>
                <div class="info-item" v-if="scanType === 'app'">
                  <span class="label">应用类型</span>
                  <span class="value">{{ task.appType ? getAppTypeName(task.appType) : '-' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">创建时间</span>
                  <span class="value">{{ task.createTime || '-' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">扫描时间</span>
                  <span class="value">{{ task.scanTime || '-' }}</span>
                </div>
                <div class="info-item" v-if="targetList.length">
                  <span class="label">扫描目标数</span>
                  <span class="value">{{ targetList.length }} 个</span>
                </div>
                <div class="info-item full-width" v-if="task.targetSummary && task.targetCount > 1">
                  <span class="label">目标摘要</span>
                  <span class="value">{{ task.targetSummary }}</span>
                </div>
                <div class="info-item" v-if="task.strategyId">
                  <span class="label">扫描策略</span>
                  <span class="value">{{ task.strategyId }}</span>
                </div>
                <div class="info-item">
                  <span class="label">检测插件数</span>
                  <span class="value">{{ pluginCountLabel }}</span>
                </div>
                <div class="info-item">
                  <span class="label">测试模式</span>
                  <span class="value">{{ configSummary.testMode }}</span>
                </div>
              </div>
              <p v-if="task.errorMessage" class="error-banner">{{ task.errorMessage }}</p>
            </div>
            <div class="info-panel inner-panel">
              <div class="panel-title">扫描能力摘要</div>
              <div class="capability-tags">
                <el-tag v-if="configSummary.crawler" size="small" type="success">Web 爬虫</el-tag>
                <el-tag v-if="configSummary.portScan" size="small" type="success">端口扫描</el-tag>
                <el-tag v-if="configSummary.proxy" size="small">代理</el-tag>
                <el-tag v-if="configSummary.safeTest" size="small">安全测试模式</el-tag>
                <el-tag v-if="!configSummary.crawler && !configSummary.portScan" size="small" type="info">仅漏洞检测</el-tag>
              </div>
            </div>
          </el-tab-pane>

          <!-- 扫描目标 -->
          <el-tab-pane v-if="targetList.length" name="targets">
            <span slot="label">扫描目标 <span class="tab-count">({{ targetList.length }})</span></span>
            <div class="targets-toolbar">
              <span class="targets-hint">点击行可筛选该目标的漏洞与站点资产；双击行快速查看漏洞。</span>
              <el-button v-if="selectedTargetId" size="small" @click="clearTargetFilter">查看全部目标</el-button>
            </div>
            <el-table
              :data="targetList"
              class="myTable targets-table"
              highlight-current-row
              :current-row-key="selectedTargetId"
              row-key="id"
              @row-click="onTargetRowClick"
              @row-dblclick="onTargetRowDblclick"
            >
              <el-table-column prop="targetUrl" label="目标地址" min-width="240" :show-overflow-tooltip="true" />
              <el-table-column prop="status" label="状态" width="100">
                <template slot-scope="scope">
                  <span :class="getStatusClass(scope.row.status)">{{ getStatusName(scope.row.status, scanType) }}</span>
                </template>
              </el-table-column>
              <el-table-column v-if="scanType === 'dyn'" prop="pageCount" label="页面" width="72" />
              <el-table-column prop="vulnCount" label="漏洞" width="72" />
              <el-table-column prop="criticalCount" label="严重" width="64" />
              <el-table-column prop="highRiskCount" label="高危" width="64" />
              <el-table-column label="操作" width="120">
                <template slot-scope="scope">
                  <el-link :underline="false" class="link_primary" @click.stop="selectTargetAndTab(scope.row, 'vulns')">
                    漏洞
                  </el-link>
                  <el-link
                    v-if="scanType === 'dyn'"
                    :underline="false"
                    class="link_primary"
                    @click.stop="selectTargetAndTab(scope.row, 'assets')"
                  >
                    资产
                  </el-link>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <!-- 漏洞 -->
          <el-tab-pane name="vulns">
            <span slot="label">漏洞 <span v-if="task.vulnCount" class="tab-count">({{ filteredVulns.length }})</span></span>
            <p v-if="selectedTargetId && selectedTargetLabel" class="target-filter-banner">
              当前仅显示目标：<strong>{{ selectedTargetLabel }}</strong>
              <el-link type="primary" :underline="false" @click="clearTargetFilter">显示全部</el-link>
            </p>
            <div class="vuln-filters">
              <el-select
                v-if="hasMultipleTargets"
                v-model="selectedTargetId"
                placeholder="按目标筛选"
                size="small"
                clearable
                class="filter-select target-filter-select"
                @change="onTargetFilterChange"
              >
                <el-option v-for="t in targetList" :key="t.id" :label="t.targetUrl" :value="t.id" />
              </el-select>
              <el-input v-model="vulnFilter.keyword" placeholder="搜索漏洞名称或 URL" size="small" clearable class="filter-input" />
              <el-select v-model="vulnFilter.risk" placeholder="风险等级" size="small" clearable class="filter-select">
                <el-option label="严重" :value="0" />
                <el-option label="高危" :value="1" />
                <el-option label="中危" :value="2" />
                <el-option label="低危" :value="3" />
                <el-option label="信息" :value="4" />
              </el-select>
              <el-button size="small" @click="resetVulnFilter">重置筛选</el-button>
            </div>
            <div class="vuln-layout">
              <div class="vuln-list-pane">
                <el-table
                  :data="pagedVulns"
                  class="myTable"
                  highlight-current-row
                  :current-row-key="selectedVulnKey"
                  row-key="_rowKey"
                  max-height="560"
                  @current-change="onVulnSelect"
                >
                  <el-table-column prop="name" label="漏洞名称" min-width="180" :show-overflow-tooltip="true" />
                  <el-table-column prop="typeName" label="类型" width="110" :show-overflow-tooltip="true" />

                  <el-table-column prop="riskLevel" label="风险" width="88">
                    <template slot-scope="scope">
                      <span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span>
                    </template>
                  </el-table-column>
                  <el-table-column
                    v-if="hasMultipleTargets && !selectedTargetId"
                    prop="targetUrl"
                    label="目标"
                    min-width="160"
                    :show-overflow-tooltip="true"
                  />
                  <el-table-column prop="url" label="URL" min-width="200" :show-overflow-tooltip="true" />
                </el-table>
                <el-pagination
                  v-if="filteredVulns.length > 0"
                  background
                  class="items-pager"
                  layout="total, prev, pager, next, sizes, jumper"
                  :total="filteredVulns.length"
                  :page-size="vulnPageSize"
                  :current-page="vulnPage"
                  :page-sizes="[20, 50, 100, 200]"
                  @current-change="onVulnPageChange"
                  @size-change="onVulnPageSizeChange"
                />
                <p v-if="!filteredVulns.length" class="empty-hint">暂无匹配的漏洞记录。</p>
              </div>
              <div class="vuln-detail-pane">
                <template v-if="selectedVuln">
                  <h3 class="vuln-detail-title">{{ selectedVuln.name }}</h3>
                  <el-descriptions :column="1" border size="small" class="vuln-desc">
                    <el-descriptions-item label="漏洞类型">{{ selectedVuln.typeName || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="风险等级">
                      <span :class="getRiskClass(selectedVuln.riskLevel)">{{ getRiskName(selectedVuln.riskLevel) }}</span>
                    </el-descriptions-item>
                    <el-descriptions-item label="漏洞 URL">{{ selectedVuln.url || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="请求方式">{{ selectedVuln.method || 'GET' }}</el-descriptions-item>
                  </el-descriptions>
                  <div class="block-section">
                    <h4>漏洞描述</h4>
                    <p class="text-block">{{ selectedVuln.description || '暂无' }}</p>
                  </div>
                  <div class="block-section ver-msg-section">
                    <h4>验证报文</h4>
                    <template v-if="verMsgPairs.length">
                      <div v-for="(pair, idx) in verMsgPairs" :key="idx" class="ver-msg-pair">
                        <div v-if="verMsgPairs.length > 1" class="ver-msg-index">#{{ idx + 1 }}</div>
                        <div class="ver-msg-block">
                          <h5>请求报文</h5>
                          <pre class="code-block">{{ formatHttpText(pair.request) }}</pre>
                        </div>
                        <div class="ver-msg-block">
                          <h5>响应报文</h5>
                          <pre class="code-block">{{ formatHttpText(pair.response) }}</pre>
                        </div>
                      </div>
                    </template>
                    <pre v-else class="code-block">{{ requestDisplayFallback }}</pre>
                  </div>
                  <div class="block-section">
                    <h4>修复建议</h4>
                    <p class="text-block">{{ selectedVuln.suggestion || '暂无' }}</p>
                  </div>
                </template>
                <p v-else class="empty-hint pane-hint">在左侧选择一条漏洞查看完整详情</p>
              </div>
            </div>
          </el-tab-pane>

          <!-- 站点资产 -->
          <el-tab-pane v-if="scanType === 'dyn'" label="站点资产" name="assets">
            <p v-if="hasMultipleTargets && !selectedTargetId" class="empty-hint assets-hint">
              多目标任务请先在「扫描目标」页点击某一目标，再查看该目标的站点树与页面列表。
            </p>
            <p v-else-if="selectedTargetLabel" class="target-filter-banner">
              站点资产：{{ selectedTargetLabel }}
              <el-link v-if="hasMultipleTargets" type="primary" :underline="false" @click="clearTargetFilter">全部目标</el-link>
            </p>
            <div class="assets-split">
              <div class="sitemap-pane">
                <div class="sitemap-toolbar">站点树 · 共 {{ flatPageCount }} 个页面</div>
                <div class="sitemap-tree-wrap">
                  <el-tree
                    v-if="siteMapTree.length"
                    :data="siteMapTree"
                    :props="treeProps"
                    default-expand-all
                    node-key="path"
                  />
                  <p v-else class="empty-hint">暂无爬取页面数据</p>
                </div>
              </div>
              <div class="pages-table-pane">
                <el-table :data="pageRows" class="myTable" max-height="520" size="small">
                  <el-table-column type="index" label="#" width="50" />
                  <el-table-column prop="url" label="页面 URL" min-width="280" :show-overflow-tooltip="true" />
                </el-table>
              </div>
            </div>
          </el-tab-pane>

          <!-- 扫描配置 -->
          <el-tab-pane label="扫描配置" name="config">
            <div class="config-summary info-grid">
              <div class="info-item">
                <span class="label">测试模式</span>
                <span class="value">{{ configSummary.testMode }}</span>
              </div>
              <div class="info-item">
                <span class="label">安全测试</span>
                <span class="value">{{ configSummary.safeTest ? '是' : '否' }}</span>
              </div>
              <div class="info-item">
                <span class="label">漏洞利用</span>
                <span class="value">{{ configSummary.vulExploit ? '是' : '否' }}</span>
              </div>
              <div class="info-item">
                <span class="label">插件 ID 数量</span>
                <span class="value">{{ configSummary.vulIdsCount }}</span>
              </div>
              <div class="info-item">
                <span class="label">Web 爬虫</span>
                <span class="value">{{ configSummary.crawler ? '已启用' : '未启用' }}</span>
              </div>
              <div class="info-item">
                <span class="label">端口扫描</span>
                <span class="value">{{ configSummary.portScan ? '已启用' : '未启用' }}</span>
              </div>
            </div>
            <div class="config-json-wrap">
              <div class="panel-title">完整配置快照（TaskTemplateJSON）</div>
              <pre class="code-block config-json">{{ scanConfigJson }}</pre>
            </div>
          </el-tab-pane>

          <!-- 执行日志 -->
          <el-tab-pane label="执行日志" name="logs">
            <app-sec-task-logs :task-id="taskId" />
          </el-tab-pane>

          <!-- 报告导出 -->
          <el-tab-pane label="报告与导出" name="export">
            <div class="export-panel">
              <p class="export-desc">导出本任务的审查结果，用于留存审计记录或二次分析。</p>
              <div class="export-actions">
                <el-button type="primary" @click="exportCsv">导出漏洞列表 (CSV)</el-button>
                <el-button @click="exportSummary">导出审查摘要 (TXT)</el-button>
                <el-button @click="copySummary">复制摘要到剪贴板</el-button>
              </div>
              <div class="export-preview">
                <div class="panel-title">摘要预览</div>
                <pre class="code-block">{{ auditSummaryPreview }}</pre>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </template>
    </div>
  </div>
</template>

<script>
import security from '@/api/security.js'
import AppSecTaskLogs from './components/AppSecTaskLogs.vue'
import {
  exportVulnsToCsv,
  buildAuditSummaryText,
  downloadTextFile
} from './utils/appsecTaskExport.js'
import {
  getRiskName,
  getRiskClass,
  getStatusName,
  getStatusClass,
  getAppTypeName,
  getTestModeLabel,
  buildSiteMapTree,
  scanTypeLabel
} from './appsecTaskLabels.js'

const VALID_TABS = ['overview', 'targets', 'vulns', 'assets', 'config', 'logs', 'export']

export default {
  name: 'AppSecTaskDetail',
  components: { AppSecTaskLogs },
  data() {
    return {
      loading: false,
      task: null,
      activeTab: 'overview',
      selectedVuln: null,
      selectedVulnKey: null,
      siteMapTree: [],
      treeProps: { label: 'name', children: 'children' },
      pollTimer: null,
      vulnFilter: {
        keyword: '',
        risk: null
      },
      selectedTargetId: null,
      vulnPage: 1,
      vulnPageSize: 20
    }
  },
  computed: {
    scanType() {
      const t = this.$route.query.type
      return t === 'app' ? 'app' : 'dyn'
    },
    taskId() {
      return this.$route.query.id || ''
    },
    typeLabel() {
      return scanTypeLabel(this.scanType)
    },
    targetList() {
      if (!this.task) return []
      if (this.task.targets && this.task.targets.length) return this.task.targets
      if (this.task.targetUrl) {
        return [
          {
            id: this.task.targetId,
            targetUrl: this.task.targetUrl,
            status: this.task.status,
            pageCount: this.task.pageCount,
            vulnCount: this.task.vulnCount,
            criticalCount: this.task.criticalCount,
            highRiskCount: this.task.highRiskCount,
            middleRiskCount: this.task.middleRiskCount,
            lowRiskCount: this.task.lowRiskCount
          }
        ]
      }
      return []
    },
    hasMultipleTargets() {
      return this.targetList.length > 1
    },
    selectedTargetLabel() {
      if (!this.selectedTargetId) return ''
      const tid = Number(this.selectedTargetId)
      const t = this.targetList.find(x => Number(x.id) === tid)
      return (t && t.targetUrl) || ''
    },
    filteredPages() {
      const pages = (this.task && this.task.pages) || []
      if (!this.selectedTargetId) {
        if (this.hasMultipleTargets) return []
        return pages
      }
      return pages.filter(p => this.pageBelongsToTarget(p, this.selectedTargetId))
    },
    flatPageCount() {
      if (this.selectedTargetId) return this.filteredPages.length
      if (!this.task) return 0
      if (this.task.pageCount) return this.task.pageCount
      return (this.task.pages && this.task.pages.length) || 0
    },
    pageRows() {
      return this.filteredPages
    },
    filteredVulns() {
      const list = (this.task && this.task.vulns) || []
      const kw = (this.vulnFilter.keyword || '').trim().toLowerCase()
      const risk = this.vulnFilter.risk
      const tid = this.selectedTargetId
      return list.filter(v => {
        if (!this.vulnBelongsToTarget(v, tid)) return false
        if (risk !== null && risk !== '' && v.riskLevel !== risk) return false
        if (!kw) return true
        const name = (v.name || '').toLowerCase()
        const url = (v.url || '').toLowerCase()
        const turl = (v.targetUrl || '').toLowerCase()
        return name.includes(kw) || url.includes(kw) || turl.includes(kw)
      })
    },
    configSummary() {
      const cfg = (this.task && this.task.scanConfig) || {}
      const wc = cfg.webCrawlerConfig || cfg.webCrawler || {}
      const ps = cfg.portScanConfig || cfg.portScan || {}
      const px = cfg.proxyConfig || cfg.proxy || {}
      const ids = cfg.vulIdsConfig
      return {
        testMode: getTestModeLabel(cfg.testMode),
        safeTest: Boolean(cfg.safeTest),
        vulExploit: Boolean(cfg.vulExploit),
        crawler: Boolean(wc.isOpen),
        portScan: Boolean(ps.isOpen),
        proxy: Boolean(px.isOpen),
        vulIdsCount: Array.isArray(ids) ? ids.length : 0
      }
    },
    pluginCountLabel() {
      const n = this.configSummary.vulIdsCount
      return n > 0 ? `${n} 个已选插件` : '未指定（按策略默认/全量规则）'
    },
    scanConfigJson() {
      const cfg = this.task && this.task.scanConfig
      if (!cfg) return '{}'
      try {
        return JSON.stringify(cfg, null, 2)
      } catch {
        return String(cfg)
      }
    },
    auditSummaryPreview() {
      if (!this.task) return ''
      return buildAuditSummaryText(this.task, this.scanType, {
        getStatusName,
        getRiskName,
        getAppTypeName,
        scanTypeLabel
      })
    },
    verMsgPairs() {
      if (!this.selectedVuln) return []
      return this.parseVerMsg(this.selectedVuln.request)
    },
    requestDisplayFallback() {
      if (!this.selectedVuln) return '暂无'
      const raw = this.selectedVuln.request
      if (!raw) return '暂无'
      return this.formatHttpText(raw)
    },
    pagedVulns() {
      const list = this.filteredVulns
      const start = (this.vulnPage - 1) * this.vulnPageSize
      const end = start + this.vulnPageSize
      return list.slice(start, end)
    }
  },
  watch: {
    '$route.query.id'() {
      this.loadDetail()
    },
    '$route.query.tab'(tab) {
      this.applyRouteTab(tab)
    },
    '$route.query.targetId'(id) {
      this.applyRouteTargetId(id)
    },
    selectedTargetId() {
      this.rebuildSiteMap()
    },
    filteredVulns(list) {
      this.vulnPage = 1
      if (!list.length) {
        this.selectedVuln = null
        this.selectedVulnKey = null
        return
      }
      const keep = list.find(v => v._rowKey === this.selectedVulnKey)
      if (!keep) this.onVulnSelect(list[0])
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
    this.applyRouteTab(this.$route.query.tab)
    this.applyRouteTargetId(this.$route.query.targetId)
    this.loadDetail()
  },
  beforeDestroy() {
    this.stopPoll()
  },
  methods: {
    getRiskName,
    getRiskClass,
    getStatusName,
    getStatusClass,
    getAppTypeName,
    parseVerMsg(raw) {
      if (raw == null || raw === '') return []
      if (Array.isArray(raw)) {
        return raw.map(item => this.normalizeVerMsgItem(item)).filter(p => p.request || p.response)
      }
      if (typeof raw === 'object') {
        const one = this.normalizeVerMsgItem(raw)
        return one.request || one.response ? [one] : []
      }
      const text = String(raw).trim()
      if (!text) return []
      try {
        const parsed = JSON.parse(text)
        return this.parseVerMsg(parsed)
      } catch {
        return [{ request: text, response: '' }]
      }
    },
    normalizeVerMsgItem(item) {
      if (!item || typeof item !== 'object') {
        return { request: '', response: '' }
      }
      return {
        request: item.request || item.Request || '',
        response: item.response || item.Response || ''
      }
    },
    formatHttpText(text) {
      if (text == null || text === '') return '暂无'
      return String(text)
        .replace(/\\r\\n/g, '\n')
        .replace(/\\n/g, '\n')
        .replace(/\\r/g, '\n')
        .replace(/\r\n/g, '\n')
    },
    normalizeTargetUrl(url) {
      if (!url) return ''
      let u = String(url).trim().toLowerCase()
      if (u.endsWith('/')) u = u.replace(/\/+$/, '')
      return u
    },
    vulnBelongsToTarget(v, tid) {
      if (!tid) return true
      const want = Number(tid)
      const vid = Number(v.targetId != null ? v.targetId : v.targetID)
      if (Number.isFinite(vid) && vid > 0 && vid === want) return true
      const t = this.targetList.find(x => Number(x.id) === want)
      if (!t) return false
      if (this.targetList.length === 1) return true
      const tt = this.normalizeTargetUrl(t.targetUrl)
      const vt = this.normalizeTargetUrl(v.targetUrl || v.url)
      if (vt && tt && (vt === tt || vt.startsWith(tt) || tt.startsWith(vt))) return true
      return false
    },
    pageBelongsToTarget(p, tid) {
      if (!tid) return true
      const want = Number(tid)
      const pid = Number(p.targetId != null ? p.targetId : p.targetID)
      if (Number.isFinite(pid) && pid > 0 && pid === want) return true
      if (this.targetList.length === 1) return true
      return false
    },
    applyRouteTab(tab) {
      if (!tab) return
      if (tab === 'sitemap') {
        this.activeTab = 'assets'
        return
      }
      if (VALID_TABS.includes(tab)) {
        this.activeTab = tab
      }
    },
    onTabClick() {
      const q = { ...this.$route.query, tab: this.activeTab }
      if (this.activeTab === 'overview') delete q.tab
      this.$router.replace({ path: this.$route.path, query: q }).catch(() => {})
    },
    applyRouteTargetId(id) {
      if (!id) {
        this.selectedTargetId = null
        return
      }
      const n = parseInt(id, 10)
      this.selectedTargetId = Number.isFinite(n) && n > 0 ? n : null
    },
    syncTargetIdToRoute() {
      const q = { ...this.$route.query }
      if (this.selectedTargetId) q.targetId = String(this.selectedTargetId)
      else delete q.targetId
      this.$router.replace({ path: this.$route.path, query: q }).catch(() => {})
    },
    onTargetFilterChange() {
      this.syncTargetIdToRoute()
    },
    clearTargetFilter() {
      this.selectedTargetId = null
      this.syncTargetIdToRoute()
    },
    onTargetRowClick(row) {
      if (!row || row.id == null) return
      this.selectedTargetId = Number(row.id)
      this.syncTargetIdToRoute()
    },
    onTargetRowDblclick(row) {
      this.selectTargetAndTab(row, 'vulns')
    },
    selectTargetAndTab(row, tab) {
      if (!row || row.id == null) return
      this.selectedTargetId = Number(row.id)
      this.activeTab = tab
      const q = { ...this.$route.query, tab, targetId: String(row.id) }
      this.$router.replace({ path: this.$route.path, query: q }).catch(() => {})
    },
    rebuildSiteMap() {
      const pages = this.filteredPages
      this.siteMapTree = buildSiteMapTree(pages)
    },
    goBack() {
      this.$router.push({ path: '/appsec/tasks', query: { tab: this.scanType } })
    },
    resetVulnFilter() {
      this.vulnFilter = { keyword: '', risk: null }
    },
    onExportCommand(cmd) {
      if (cmd === 'csv') this.exportCsv()
      else if (cmd === 'summary') this.exportSummary()
    },
    exportCsv() {
      if (!this.task) return
      exportVulnsToCsv(this.task, this.scanType, {
        getRiskName: this.getRiskName
      })
      this.$message.success('漏洞列表已导出')
    },
    exportSummary() {
      if (!this.task) return
      downloadTextFile(`appsec-${this.task.id}-summary.txt`, this.auditSummaryPreview)
      this.$message.success('审查摘要已导出')
    },
    async copySummary() {
      try {
        await navigator.clipboard.writeText(this.auditSummaryPreview)
        this.$message.success('已复制到剪贴板')
      } catch {
        this.$message.warning('复制失败，请手动选择摘要文本复制')
      }
    },
    async generateReport() {
      if (!this.taskId) return
      try {
        const res = await security.generateSecurityReport({ module: 'app', taskId: Number(this.taskId) })
        if (res.code === 200 && res.data) {
          this.$message({ message: '报告已生成，前往报告中心查看', type: 'success' })
        } else {
          this.$message({ message: res.msg || '生成报告失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '生成报告失败: ' + (e.message || ''), type: 'error' })
      }
    },
    async loadDetail() {
      if (!this.taskId) {
        this.$message.warning('缺少任务 ID')
        this.goBack()
        return
      }
      this.loading = true
      try {
        const api =
          this.scanType === 'app' ? security.getAppSpecificScanDetail : security.getDynamicScanDetail
        const res = await api({ id: this.taskId })
        if (res.code == 200 && res.data) {
          this.applyTask(res.data)
          this.startPollIfNeeded()
        } else {
          this.$message.error(res.msg || '加载失败')
        }
      } catch {
        this.$message.error('加载任务详情失败')
      } finally {
        this.loading = false
      }
    },
    applyTask(data) {
      const defaultTargetId =
        data.targetId || (data.targets && data.targets[0] && data.targets[0].id) || null
      const vulns = (data.vulns || []).map((v, i) => {
        const targetId = v.targetId != null ? v.targetId : v.targetID
        return {
          ...v,
          targetId: targetId != null && targetId !== '' ? Number(targetId) : defaultTargetId,
          _rowKey: String(i)
        }
      })
      this.task = { ...data, vulns }
      if (this.$route.query.targetId) {
        this.applyRouteTargetId(this.$route.query.targetId)
      }
      this.rebuildSiteMap()
      if (this.filteredVulns.length) {
        const keep = this.filteredVulns.find(v => v._rowKey === this.selectedVulnKey)
        this.onVulnSelect(keep || this.filteredVulns[0])
      } else {
        this.selectedVuln = null
        this.selectedVulnKey = null
      }
    },
    onVulnSelect(row) {
      if (!row) return
      this.selectedVuln = row
      this.selectedVulnKey = row._rowKey
    },
    startPollIfNeeded() {
      this.stopPoll()
      if (!this.task || this.task.status !== 2) return
      this.pollTimer = setInterval(() => this.loadDetailQuiet(), 3000)
    },
    async loadDetailQuiet() {
      if (!this.taskId) return
      try {
        const api =
          this.scanType === 'app' ? security.getAppSpecificScanDetail : security.getDynamicScanDetail
        const res = await api({ id: this.taskId })
        if (res.code == 200 && res.data) {
          this.applyTask(res.data)
          if (res.data.status !== 2) this.stopPoll()
        }
      } catch {
        /* ignore */
      }
    },
    stopPoll() {
      if (this.pollTimer) {
        clearInterval(this.pollTimer)
        this.pollTimer = null
      }
    },
    onVulnPageChange(page) {
      this.vulnPage = page
    },
    onVulnPageSizeChange(size) {
      this.vulnPageSize = size
      this.vulnPage = 1
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';
@import './css/appsec-task-detail.less';
</style>
