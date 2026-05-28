<template>
  <div class="security-container task-detail-page datasec-detail-page">
    <div class="detail-topbar list_box">
      <el-link :underline="false" class="link-back" @click="goBack">
        <i class="el-icon-arrow-left"></i> 返回
      </el-link>
      <div class="topbar-main" v-if="task">
        <h1 class="task-title">{{ task.name }}</h1>
        <p class="task-meta">
          <span>{{ typeLabel }}</span>
          <span class="dot">·</span>
          <span>{{ getDBTypeName(task.dbType) }}</span>
          <span class="dot">·</span>
          <span>{{ task.targetSummary || `${task.dbHost}:${task.dbPort}` }}</span>
        </p>
        <div class="task-tags">
          <span :class="getStatusClass(task.status)">{{ getStatusName(task.status) }}</span>
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
            <el-dropdown-item v-if="isDb" command="checks">基线检查 CSV</el-dropdown-item>
            <el-dropdown-item v-if="isDb" command="cve">CVE 列表 CSV</el-dropdown-item>
            <el-dropdown-item v-if="showSensitiveTabs" command="findings">敏感字段 CSV</el-dropdown-item>
            <el-dropdown-item command="summary">审查摘要 TXT</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </div>

    <div v-loading="loading" class="detail-body">
      <template v-if="task">
        <el-tabs v-model="activeTab" class="detail-tabs list_box" @tab-click="onTabClick">
          <el-tab-pane label="概览" name="overview">
            <div v-if="isDb" class="stat-row">
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
            </div>
            <div v-if="isDb && hasSensitiveScan" class="stat-row stat-row-secondary">
              <div class="stat-card total">
                <span class="stat-label">敏感发现</span>
                <span class="stat-value">{{ task.totalCount || 0 }}</span>
              </div>
              <div class="stat-card high">
                <span class="stat-label">高敏感</span>
                <span class="stat-value">{{ task.highCount || 0 }}</span>
              </div>
              <div class="stat-card medium">
                <span class="stat-label">中敏感</span>
                <span class="stat-value">{{ task.mediumCount || 0 }}</span>
              </div>
              <div class="stat-card low">
                <span class="stat-label">低敏感</span>
                <span class="stat-value">{{ task.lowCount || 0 }}</span>
              </div>
            </div>
            <div v-if="!isDb" class="stat-row">
              <div class="stat-card total">
                <span class="stat-label">发现条数</span>
                <span class="stat-value">{{ task.totalCount || 0 }}</span>
              </div>
              <div class="stat-card high">
                <span class="stat-label">高敏感</span>
                <span class="stat-value">{{ task.highCount || 0 }}</span>
              </div>
              <div class="stat-card medium">
                <span class="stat-label">中敏感</span>
                <span class="stat-value">{{ task.mediumCount || 0 }}</span>
              </div>
              <div class="stat-card low">
                <span class="stat-label">低敏感</span>
                <span class="stat-value">{{ task.lowCount || 0 }}</span>
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
                  <span class="label">任务类型</span>
                  <span class="value">{{ typeLabel }}</span>
                </div>
                <div class="info-item">
                  <span class="label">数据库类型</span>
                  <span :class="getDBTypeClass(task.dbType)">{{ getDBTypeName(task.dbType) }}</span>
                </div>
                <div class="info-item">
                  <span class="label">扫描目标</span>
                  <span class="value">{{ task.targetSummary || `${task.dbHost}:${task.dbPort}` }}</span>
                </div>
                <div class="info-item" v-if="targetList.length > 1">
                  <span class="label">目标数量</span>
                  <span class="value">{{ targetList.length }} 个</span>
                </div>
                <div class="info-item" v-if="targetList.length <= 1">
                  <span class="label">数据库名</span>
                  <span class="value">{{ task.dbName || '-' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">创建时间</span>
                  <span class="value">{{ task.createTime || '-' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">{{ isDb ? '检查时间' : '扫描时间' }}</span>
                  <span class="value">{{ timeLabel || '-' }}</span>
                </div>
                <div class="info-item" v-if="isDb">
                  <span class="label">实例版本</span>
                  <span class="value mono">{{ displayDbVersion }}</span>
                </div>
                <div class="info-item" v-if="isDb">
                  <span class="label">基线检查</span>
                  <span class="value">共 {{ scanSummary.baselineTotal }} 项 · 通过 {{ scanSummary.baselinePass }} · 不通过 {{ scanSummary.baselineFail }}<template v-if="scanSummary.baselineError"> · 异常 {{ scanSummary.baselineError }}</template></span>
                </div>
                <div class="info-item" v-if="isDb">
                  <span class="label">CVE 命中</span>
                  <span class="value">{{ scanSummary.cveMatchCount }} 条</span>
                </div>
                <div class="info-item" v-if="hasSensitiveScan">
                  <span class="label">敏感数据</span>
                  <span class="value">发现 {{ task.totalCount || 0 }} 条 · 高 {{ task.highCount || 0 }} · 中 {{ task.mediumCount || 0 }} · 低 {{ task.lowCount || 0 }}</span>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane v-if="targetList.length" name="targets">
            <span slot="label">扫描目标 <span class="tab-count">({{ targetList.length }})</span></span>
            <div class="targets-toolbar">
              <span class="targets-hint">点击行筛选下方检查/漏洞/敏感字段结果；再次点击取消筛选。</span>
            </div>
            <el-table
              :data="targetList"
              class="myTable targets-table"
              max-height="480"
              size="small"
              highlight-current-row
              @row-click="onTargetRowClick"
            >
              <el-table-column prop="targetUrl" label="目标" min-width="180" :show-overflow-tooltip="true" />
              <el-table-column v-if="isDb" label="库类型" width="96">
                <template slot-scope="scope">{{ getDBTypeName(scope.row.dbType) }}</template>
              </el-table-column>
              <el-table-column prop="dbHost" label="地址" width="130" :show-overflow-tooltip="true" />
              <el-table-column prop="dbPort" label="端口" width="72" />
              <el-table-column v-if="isDb" prop="dbName" label="库名" width="100" :show-overflow-tooltip="true" />
              <el-table-column v-if="isDb" label="版本" width="120" :show-overflow-tooltip="true">
                <template slot-scope="scope">{{ scope.row.dbVersion || '未识别' }}</template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="88">
                <template slot-scope="scope">
                  <span :class="getStatusClass(scope.row.status)">{{ getStatusName(scope.row.status) }}</span>
                </template>
              </el-table-column>
              <el-table-column v-if="isDb" label="基线" width="120">
                <template slot-scope="scope">
                  {{ scope.row.baselineTotal || 0 }} 项
                  <span v-if="scope.row.baselineFail" class="fail-hint">/ 不通过 {{ scope.row.baselineFail }}</span>
                </template>
              </el-table-column>
              <el-table-column v-if="isDb" label="CVE" width="72">
                <template slot-scope="scope">{{ scope.row.cveMatchCount || 0 }}</template>
              </el-table-column>
              <el-table-column v-if="isDb" label="风险项" width="88">
                <template slot-scope="scope">
                  {{ (scope.row.criticalCount || 0) + (scope.row.highRiskCount || 0) + (scope.row.middleRiskCount || 0) + (scope.row.lowRiskCount || 0) }}
                </template>
              </el-table-column>
              <el-table-column v-if="isDb && hasSensitiveScan" label="敏感" width="72">
                <template slot-scope="scope">{{ scope.row.totalCount || 0 }}</template>
              </el-table-column>
              <el-table-column v-if="!isDb" label="发现" width="88">
                <template slot-scope="scope">{{ scope.row.totalCount || 0 }}</template>
              </el-table-column>
              <el-table-column v-if="isDb" label="备注" min-width="140" :show-overflow-tooltip="true">
                <template slot-scope="scope">{{ scope.row.errorMessage || '-' }}</template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane v-if="isDb" name="checks">
            <span slot="label">基线检查项 <span class="tab-count">({{ scanSummary.baselineTotal }})</span></span>
            <div class="check-filters">
              <el-select v-if="targetList.length > 1" v-model="selectedTargetId" placeholder="全部目标" size="small" clearable class="filter-select target-filter">
                <el-option label="全部目标" :value="null" />
                <el-option v-for="t in targetList" :key="t.id" :label="t.targetUrl || t.dbHost" :value="t.id" />
              </el-select>
              <el-input v-model="checkFilter.keyword" placeholder="搜索描述/建议" size="small" clearable class="filter-input" />
              <el-select v-model="checkFilter.category" placeholder="分类" size="small" clearable class="filter-select">
                <el-option v-for="c in categoryOptions" :key="c.value" :label="c.label" :value="c.value" />
              </el-select>
              <el-select v-model="checkFilter.risk" placeholder="风险" size="small" clearable class="filter-select">
                <el-option label="严重" :value="0" />
                <el-option label="高危" :value="1" />
                <el-option label="中危" :value="2" />
                <el-option label="低危" :value="3" />
              </el-select>
              <el-button size="small" @click="resetCheckFilter">重置</el-button>
            </div>
            <el-table :data="filteredBaselineChecks" class="myTable" max-height="560" size="small">
              <el-table-column prop="category" label="检查类别" width="120">
                <template slot-scope="scope">{{ getCategoryName(scope.row.category) }}</template>
              </el-table-column>
              <el-table-column prop="ruleName" label="规则" min-width="160" :show-overflow-tooltip="true" />
              <el-table-column prop="riskLevel" label="风险" width="88">
                <template slot-scope="scope">
                  <span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="result" label="结果" width="100" />
              <el-table-column prop="actualValue" label="实际值" min-width="140" :show-overflow-tooltip="true" />
              <el-table-column prop="description" label="描述" min-width="180" :show-overflow-tooltip="true" />
              <el-table-column label="操作" width="72" fixed="right">
                <template slot-scope="scope">
                  <el-link :underline="false" class="link_primary" @click="openCheckDetail(scope.row)">详情</el-link>
                </template>
              </el-table-column>
            </el-table>
            <p v-if="!filteredBaselineChecks.length" class="empty-hint">{{ baselineEmptyHint }}</p>
          </el-tab-pane>

          <el-tab-pane v-if="isDb" name="cve">
            <span slot="label">版本漏洞(CVE) <span class="tab-count">({{ cveFailCount }})</span></span>
            <p v-if="dbVersion" class="cve-version-hint">探测版本：<strong>{{ dbVersion }}</strong> · 扫描 CVE 库并做 CPE 版本匹配</p>
            <p v-else-if="displayDbVersion !== '未识别'" class="cve-version-hint">探测版本：<strong>{{ displayDbVersion }}</strong> · 扫描 CVE 库并做 CPE 版本匹配</p>
            <p v-else class="cve-version-hint">未能识别数据库版本，CVE 匹配可能已跳过</p>
            <div class="check-filters">
              <el-select v-if="targetList.length > 1" v-model="selectedTargetId" placeholder="全部目标" size="small" clearable class="filter-select target-filter">
                <el-option label="全部目标" :value="null" />
                <el-option v-for="t in targetList" :key="t.id" :label="t.targetUrl || t.dbHost" :value="t.id" />
              </el-select>
              <el-input v-model="cveFilter.keyword" placeholder="搜索 CVE/描述" size="small" clearable class="filter-input" />
              <el-button size="small" @click="cveFilter.keyword = ''">重置</el-button>
            </div>
            <el-table :data="filteredCveItems" class="myTable" max-height="560" size="small">
              <el-table-column prop="ruleName" label="CVE / 规则" min-width="220" :show-overflow-tooltip="true" />
              <el-table-column prop="riskLevel" label="风险" width="88">
                <template slot-scope="scope">
                  <span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="result" label="结果" width="100" />
              <el-table-column prop="actualValue" label="版本信息" min-width="160" :show-overflow-tooltip="true" />
              <el-table-column prop="suggestion" label="修复建议" min-width="180" :show-overflow-tooltip="true" />
              <el-table-column label="操作" width="72" fixed="right">
                <template slot-scope="scope">
                  <el-link :underline="false" class="link_primary" @click="openCveDetail(scope.row)">详情</el-link>
                </template>
              </el-table-column>
            </el-table>
            <p v-if="!filteredCveItems.length" class="empty-hint">未发现与当前版本匹配的 CVE</p>
          </el-tab-pane>

          <el-tab-pane v-if="showSensitiveTabs" name="distribution">
            <span slot="label">类型分布</span>
            <div v-if="typeStats.length" class="type-distribution">
              <div v-for="item in typeStats" :key="item.dataType" class="type-chip">
                {{ getDataTypeName(item.dataType) }}<strong>{{ item.count }}</strong>
              </div>
            </div>
            <p v-else class="empty-hint">暂无类型分布数据</p>
          </el-tab-pane>

          <el-tab-pane v-if="showSensitiveTabs" name="findings">
            <span slot="label">敏感字段 <span class="tab-count">({{ filteredFindings.length }})</span></span>
            <div class="check-filters">
              <el-select v-if="targetList.length > 1" v-model="selectedTargetId" placeholder="全部目标" size="small" clearable class="filter-select target-filter">
                <el-option label="全部目标" :value="null" />
                <el-option v-for="t in targetList" :key="t.id" :label="t.targetUrl || t.dbHost" :value="t.id" />
              </el-select>
              <el-input v-model="findingFilter.keyword" placeholder="搜索目标/库名/位置/样例" size="small" clearable class="filter-input" />
              <el-select v-model="findingFilter.dataType" placeholder="数据类型" size="small" clearable class="filter-select">
                <el-option v-for="d in dataTypeOptions" :key="d.value" :label="d.label" :value="d.value" />
              </el-select>
              <el-select v-model="findingFilter.level" placeholder="敏感等级" size="small" clearable class="filter-select">
                <el-option label="高敏感" :value="1" />
                <el-option label="中敏感" :value="2" />
                <el-option label="低敏感" :value="3" />
              </el-select>
              <el-button size="small" @click="resetFindingFilter">重置</el-button>
            </div>
            <el-table :data="filteredFindings" class="myTable" max-height="560" size="small">
              <el-table-column v-if="targetList.length > 1" prop="targetLabel" label="目标" min-width="180" :show-overflow-tooltip="true" />
              <el-table-column prop="dbName" label="库名" min-width="120" :show-overflow-tooltip="true" />
              <el-table-column prop="location" label="位置" min-width="220" :show-overflow-tooltip="true" />
              <el-table-column prop="dataType" label="数据类型" width="110">
                <template slot-scope="scope">{{ getDataTypeName(scope.row.dataType) }}</template>
              </el-table-column>
              <el-table-column prop="sensitivityLevel" label="敏感等级" width="100">
                <template slot-scope="scope">
                  <span :class="getSensitivityClass(scope.row.sensitivityLevel)">
                    {{ getSensitivityName(scope.row.sensitivityLevel) }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column prop="sampleData" label="样例" min-width="180" :show-overflow-tooltip="true" />
              <el-table-column prop="count" label="数量" width="72" />
              <el-table-column label="操作" width="72" fixed="right">
                <template slot-scope="scope">
                  <el-link :underline="false" class="link_primary" @click="openFindingDetail(scope.row)">详情</el-link>
                </template>
              </el-table-column>
            </el-table>
            <p v-if="!filteredFindings.length" class="empty-hint">暂无敏感字段记录</p>
          </el-tab-pane>

          <el-tab-pane name="logs">
            <span slot="label">执行日志</span>
            <app-sec-task-logs :task-id="taskId" />
          </el-tab-pane>
        </el-tabs>
      </template>
      <p v-else-if="!loading" class="empty-hint list_box">任务不存在或加载失败</p>
    </div>

    <el-dialog :title="detailDialogTitle" :visible.sync="detailDialogVisible" width="720px" custom-class="theme-dialog">
      <div v-if="detailDialogRow" class="check-detail-panel">
        <div v-for="item in detailDialogFields" :key="item.label" class="detail-field">
          <span class="detail-label">{{ item.label }}</span>
          <span class="detail-value" :class="{ mono: item.mono }">{{ item.value || '-' }}</span>
        </div>
      </div>
      <span slot="footer">
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'
import AppSecTaskLogs from './components/AppSecTaskLogs.vue'
import {
  kindLabel,
  getDBTypeName,
  getDBTypeClass,
  getRiskName,
  getRiskClass,
  getStatusName,
  getStatusClass,
  getCategoryName,
  getDataTypeName,
  getSensitivityName,
  getSensitivityClass
} from './datasecTaskLabels.js'
import { buildDatasecSummaryText, downloadTextFile, exportDatasecRowsCsv } from './utils/datasecTaskExport.js'

const VALID_TABS_SENSITIVE = ['overview', 'targets', 'distribution', 'findings', 'logs']

export default {
  name: 'DataSecTaskDetail',
  components: { AppSecTaskLogs },
  data() {
    return {
      loading: false,
      task: null,
      activeTab: 'overview',
      selectedTargetId: null,
      checkFilter: { keyword: '', category: null, risk: null },
      cveFilter: { keyword: '' },
      findingFilter: { keyword: '', dataType: null, level: null },
      detailDialogVisible: false,
      detailDialogTitle: '',
      detailDialogRow: null,
      detailDialogKind: '',
      categoryOptions: [
        { value: 1, label: '身份认证' },
        { value: 2, label: '权限控制' },
        { value: 3, label: '配置安全' },
        { value: 4, label: '审计日志' },
        { value: 5, label: '网络安全' },
        { value: 6, label: '加密' },
        { value: 7, label: 'SQL 注入' },
        { value: 8, label: '敏感数据识别' }
      ],
      dataTypeOptions: [
        { value: 1, label: '身份证号' }, { value: 2, label: '银行卡号' }, { value: 3, label: '护照号' },
        { value: 4, label: '手机号' }, { value: 5, label: '邮箱' }, { value: 6, label: '地址' },
        { value: 7, label: '出生日期' }, { value: 8, label: '姓名' }, { value: 9, label: 'Token' },
        { value: 10, label: '证书信息' }, { value: 11, label: '密码哈希' }
      ]
    }
  },
  computed: {
    taskId() {
      return this.$route.query.id || ''
    },
    kind() {
      return this.$route.query.kind === 'sensitive' ? 'sensitive' : 'db'
    },
    isDb() {
      return this.kind === 'db'
    },
    hasSensitiveScan() {
      if (!this.task) return false
      if (!this.isDb) return true
      return !!(this.task.scanSensitive || this.task.totalCount || (this.task.sensitiveItems && this.task.sensitiveItems.length))
    },
    showSensitiveTabs() {
      return this.hasSensitiveScan
    },
    validTabs() {
      if (!this.isDb) return VALID_TABS_SENSITIVE
      const tabs = ['overview', 'targets', 'checks', 'cve']
      if (this.hasSensitiveScan) {
        tabs.push('distribution', 'findings')
      }
      tabs.push('logs')
      return tabs
    },
    typeLabel() {
      return kindLabel(this.kind)
    },
    timeLabel() {
      if (!this.task) return ''
      return this.isDb ? this.task.checkTime : this.task.scanTime
    },
    targetList() {
      return (this.task && this.task.targets) || []
    },
    checkItems() {
      return (this.task && this.task.items) || []
    },
    scopedCheckItems() {
      if (!this.selectedTargetId) return this.checkItems
      return this.checkItems.filter((r) => Number(r.targetId) === Number(this.selectedTargetId))
    },
    dbVersion() {
      const rows = this.scopedCheckItems.filter((r) => r.ruleName === '数据库版本识别')
      const row = rows[0] || this.checkItems.find((r) => r.ruleName === '数据库版本识别')
      return row && row.actualValue && row.actualValue !== 'unknown' ? row.actualValue : ''
    },
    displayDbVersion() {
      if (this.dbVersion) return this.dbVersion
      const fromTarget = (this.targetList || []).map((t) => t.dbVersion).find((v) => v)
      return fromTarget || '未识别'
    },
    scanSummary() {
      const list = this.targetList || []
      if (!list.length) {
        return {
          baselineTotal: this.baselineItems.length,
          baselinePass: this.baselineItems.filter((r) => r.result === '通过').length,
          baselineFail: this.baselineItems.filter((r) => r.result === '不通过').length,
          baselineError: this.baselineItems.filter((r) => r.result !== '通过' && r.result !== '不通过').length,
          cveMatchCount: this.cveFailCount
        }
      }
      return list.reduce(
        (acc, t) => ({
          baselineTotal: acc.baselineTotal + (t.baselineTotal || 0),
          baselinePass: acc.baselinePass + (t.baselinePass || 0),
          baselineFail: acc.baselineFail + (t.baselineFail || 0),
          baselineError: acc.baselineError + (t.baselineError || 0),
          cveMatchCount: acc.cveMatchCount + (t.cveMatchCount || 0)
        }),
        { baselineTotal: 0, baselinePass: 0, baselineFail: 0, baselineError: 0, cveMatchCount: 0 }
      )
    },
    baselineEmptyHint() {
      if (this.scanSummary.baselineTotal > 0) return '当前筛选条件下暂无检查项'
      return '暂无检查项记录，请重新执行扫描或在「执行日志」中查看详情'
    },
    baselineItems() {
      return this.scopedCheckItems.filter((r) => !r.isCve && r.ruleName !== '数据库版本识别' && r.ruleName !== 'CVE 版本匹配')
    },
    cveItems() {
      return this.scopedCheckItems.filter((r) => r.isCve && r.result === '不通过')
    },
    cveFailCount() {
      return this.checkItems.filter((r) => r.isCve && r.result === '不通过').length
    },
    typeStats() {
      if (!this.selectedTargetId) {
        return (this.task && this.task.typeStats) || []
      }
      const counts = {}
      this.scopedFindingItems.forEach((r) => {
        counts[r.dataType] = (counts[r.dataType] || 0) + 1
      })
      return Object.keys(counts).map((k) => ({ dataType: Number(k), count: counts[k] }))
    },
    findingItems() {
      if (!this.task) return []
      if (this.isDb) return this.task.sensitiveItems || []
      return this.task.items || []
    },
    scopedFindingItems() {
      if (!this.selectedTargetId) return this.findingItems
      return this.findingItems.filter((r) => Number(r.targetId) === Number(this.selectedTargetId))
    },
    filteredBaselineChecks() {
      const kw = (this.checkFilter.keyword || '').trim().toLowerCase()
      return this.baselineItems.filter((row) => {
        if (this.checkFilter.category != null && this.checkFilter.category !== '' && row.category !== this.checkFilter.category) {
          return false
        }
        if (this.checkFilter.risk != null && this.checkFilter.risk !== '' && row.riskLevel !== this.checkFilter.risk) {
          return false
        }
        if (!kw) return true
        const blob = [row.description, row.suggestion, row.result, getCategoryName(row.category)].join(' ').toLowerCase()
        return blob.includes(kw)
      })
    },
    filteredCveItems() {
      const kw = (this.cveFilter.keyword || '').trim().toLowerCase()
      if (!kw) return this.cveItems
      return this.cveItems.filter((row) => {
        const blob = [row.ruleName, row.description, row.actualValue, row.suggestion].join(' ').toLowerCase()
        return blob.includes(kw)
      })
    },
    filteredFindings() {
      const kw = (this.findingFilter.keyword || '').trim().toLowerCase()
      return this.scopedFindingItems.filter((row) => {
        if (this.findingFilter.dataType != null && this.findingFilter.dataType !== '' && row.dataType !== this.findingFilter.dataType) {
          return false
        }
        if (this.findingFilter.level != null && this.findingFilter.level !== '' && row.sensitivityLevel !== this.findingFilter.level) {
          return false
        }
        if (!kw) return true
        const blob = [row.targetLabel, row.dbName, row.location, row.tableName, row.columnName, row.sampleData, getDataTypeName(row.dataType)].join(' ').toLowerCase()
        return blob.includes(kw)
      })
    },
    detailDialogFields() {
      const row = this.detailDialogRow
      if (!row) return []
      if (this.detailDialogKind === 'check') {
        return [
          { label: '检查类别', value: getCategoryName(row.category) },
          { label: '规则名称', value: row.ruleName },
          { label: '风险等级', value: getRiskName(row.riskLevel) },
          { label: '检查结果', value: row.result },
          { label: '期望值', value: row.expectedValue },
          { label: '实际值', value: row.actualValue, mono: true },
          { label: '描述', value: row.description },
          { label: '修复建议', value: row.suggestion }
        ]
      }
      if (this.detailDialogKind === 'cve') {
        return [
          { label: 'CVE / 规则', value: row.ruleName },
          { label: '风险等级', value: getRiskName(row.riskLevel) },
          { label: '结果', value: row.result },
          { label: '版本信息', value: row.actualValue, mono: true },
          { label: '描述', value: row.description },
          { label: '修复建议', value: row.suggestion }
        ]
      }
      if (this.detailDialogKind === 'finding') {
        return [
          { label: '目标', value: row.targetLabel, mono: true },
          { label: '库名', value: row.dbName, mono: true },
          { label: '位置', value: row.location, mono: true },
          { label: '表名', value: row.tableName, mono: true },
          { label: '字段名', value: row.columnName, mono: true },
          { label: '数据类型', value: getDataTypeName(row.dataType) },
          { label: '敏感等级', value: getSensitivityName(row.sensitivityLevel) },
          { label: '样例数据', value: row.sampleData, mono: true },
          { label: '数量', value: String(row.count || 0) }
        ]
      }
      return []
    },
    auditSummaryPreview() {
      return buildDatasecSummaryText(
        {
          task: this.task,
          kind: this.kind,
          targetCount: this.targetList.length,
          timeLabel: this.timeLabel,
          scanSummary: this.scanSummary,
          cveFailCount: this.cveFailCount,
          findingCount: this.findingItems.length
        },
        {
          getDBTypeName: this.getDBTypeName
        }
      )
    }
  },
  watch: {
    '$route.query.id'() {
      this.syncFromRoute()
      this.loadDetail()
    },
    '$route.query.kind'() {
      this.syncFromRoute()
      this.loadDetail()
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/datasec/tasks'
    this.syncFromRoute()
    this.loadDetail()
  },
  methods: {
    getDBTypeName,
    getDBTypeClass,
    getRiskName,
    getRiskClass,
    getStatusName,
    getStatusClass,
    getCategoryName,
    getDataTypeName,
    getSensitivityName,
    getSensitivityClass,
    onTargetRowClick(row) {
      const id = row && row.id
      if (!id) return
      this.selectedTargetId = this.selectedTargetId === id ? null : id
    },
    syncFromRoute() {
      const tab = this.$route.query.tab
      if (tab && this.validTabs.includes(tab)) {
        this.activeTab = tab
      } else {
        this.activeTab = 'overview'
      }
    },
    onTabClick() {
      this.$router.replace({
        query: { ...this.$route.query, tab: this.activeTab }
      }).catch(() => {})
    },
    goBack() {
      const from = this.$route.query.from || (this.isDb ? 'db' : 'sensitive')
      this.$router.push({ path: '/datasec/tasks', query: { tab: from } })
    },
    async loadDetail() {
      if (!this.taskId) {
        this.$message({ message: '缺少任务 ID', type: 'warning' })
        return
      }
      this.loading = true
      try {
        const api = this.isDb ? security.getDBCheckDetail : security.getSensitiveScanDetail
        const res = await api({ id: this.taskId })
        if (res.code === 200 && res.data) {
          this.task = res.data
          if (this.selectedTargetId && !(this.task.targets || []).some((t) => t.id === this.selectedTargetId)) {
            this.selectedTargetId = null
          }
        } else {
          this.task = null
          this.$message({ message: res.msg || '加载失败', type: 'error' })
        }
      } catch (e) {
        this.task = null
        this.$message({ message: '加载失败: ' + (e.message || ''), type: 'error' })
      } finally {
        this.loading = false
      }
    },
    resetCheckFilter() {
      this.checkFilter = { keyword: '', category: null, risk: null }
    },
    resetFindingFilter() {
      this.findingFilter = { keyword: '', dataType: null, level: null }
    },
    openCheckDetail(row) {
      this.detailDialogKind = 'check'
      this.detailDialogTitle = '基线检查详情'
      this.detailDialogRow = row
      this.detailDialogVisible = true
    },
    openCveDetail(row) {
      this.detailDialogKind = 'cve'
      this.detailDialogTitle = 'CVE 详情'
      this.detailDialogRow = row
      this.detailDialogVisible = true
    },
    openFindingDetail(row) {
      this.detailDialogKind = 'finding'
      this.detailDialogTitle = '敏感字段详情'
      this.detailDialogRow = row
      this.detailDialogVisible = true
    },
    onExportCommand(command) {
      if (command === 'checks') {
        this.exportChecksCsv()
      } else if (command === 'cve') {
        this.exportCveCsv()
      } else if (command === 'findings') {
        this.exportFindingsCsv()
      } else if (command === 'summary') {
        this.exportSummary()
      }
    },
    exportChecksCsv() {
      if (!this.filteredBaselineChecks.length) {
        this.$message.warning('暂无可导出的基线检查项')
        return
      }
      exportDatasecRowsCsv(this.taskId, 'checks', this.filteredBaselineChecks, {
        getCategoryName: this.getCategoryName,
        getRiskName: this.getRiskName,
        getDataTypeName: this.getDataTypeName,
        getSensitivityName: this.getSensitivityName
      })
      this.$message.success('基线检查项已导出')
    },
    exportCveCsv() {
      if (!this.filteredCveItems.length) {
        this.$message.warning('暂无可导出的 CVE 结果')
        return
      }
      exportDatasecRowsCsv(this.taskId, 'cve', this.filteredCveItems, {
        getCategoryName: this.getCategoryName,
        getRiskName: this.getRiskName,
        getDataTypeName: this.getDataTypeName,
        getSensitivityName: this.getSensitivityName
      })
      this.$message.success('CVE 列表已导出')
    },
    exportFindingsCsv() {
      if (!this.filteredFindings.length) {
        this.$message.warning('暂无可导出的敏感字段')
        return
      }
      exportDatasecRowsCsv(this.taskId, 'findings', this.filteredFindings, {
        getCategoryName: this.getCategoryName,
        getRiskName: this.getRiskName,
        getDataTypeName: this.getDataTypeName,
        getSensitivityName: this.getSensitivityName
      })
      this.$message.success('敏感字段已导出')
    },
    exportSummary() {
      downloadTextFile(`datasec-${this.taskId}-summary.txt`, this.auditSummaryPreview)
      this.$message.success('审查摘要已导出')
    },
    async generateReport() {
      if (!this.taskId) return
      const module = this.isDb ? 'data' : 'data'
      try {
        const res = await security.generateSecurityReport({ module, taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.$message({ message: '报告已生成，前往报告中心查看', type: 'success' })
        } else {
          this.$message({ message: res.msg || '生成报告失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '生成报告失败: ' + (e.message || ''), type: 'error' })
      }
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';
@import './css/datasec-task-detail.less';

.db-mysql, .db-postgresql, .db-mongodb, .db-redis, .db-couchdb {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}
.db-mysql { background: rgba(245, 101, 101, 0.2); color: #f56565; }
.db-postgresql { background: rgba(97, 175, 254, 0.2); color: #61affe; }
.db-mongodb { background: rgba(67, 153, 52, 0.2); color: #439934; }
.db-redis { background: rgba(231, 162, 13, 0.2); color: #e7a20d; }
.db-couchdb { background: rgba(125, 86, 205, 0.2); color: #7d56cd; }

.status-wait, .status-running, .status-complete {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}
.status-wait { background: rgba(234, 179, 8, 0.2); color: #eab308; }
.status-running { background: rgba(59, 130, 246, 0.2); color: #3b82f6; }
.status-complete { background: rgba(16, 185, 129, 0.2); color: #10b981; }

.risk-critical, .risk-high, .risk-medium, .risk-low, .risk-info {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}
.risk-critical { background: rgba(220, 38, 38, 0.25); color: #ef4444; }
.risk-high { background: rgba(234, 88, 12, 0.25); color: #f97316; }
.risk-medium { background: rgba(234, 179, 8, 0.25); color: #eab308; }
.risk-low { background: rgba(59, 130, 246, 0.25); color: #3b82f6; }
.risk-info { background: rgba(148, 163, 184, 0.25); color: #94a3b8; }

.sensitivity-high { color: #ef4444; font-weight: 600; }
.sensitivity-medium { color: #eab308; font-weight: 600; }
.sensitivity-low { color: #3b82f6; }

.fail-hint { color: #f97316; margin-left: 4px; }.stat-row-secondary {
  margin-top: 8px;
}

.check-detail-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-field {
  display: grid;
  grid-template-columns: 88px 1fr;
  gap: 12px;
  align-items: start;
}

.detail-label {
  color: #94a3b8;
  font-size: 13px;
}

.detail-value {
  color: #e2e8f0;
  font-size: 13px;
  line-height: 1.6;
  word-break: break-all;
  white-space: pre-wrap;
}

.detail-value.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

</style>
