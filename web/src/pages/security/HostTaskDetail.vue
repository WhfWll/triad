<template>
  <div class="security-container task-detail-page">
    <div class="detail-topbar list_box">
      <el-link :underline="false" class="link-back" @click="$router.push('/hostsec/tasks')">
        <i class="el-icon-arrow-left"></i> 返回
      </el-link>
      <div class="topbar-main">
        <h1 class="task-title">
          任务详情 · {{ kindLabel }}
          <span class="title-sub">#{{ taskId }}</span>
        </h1>
      </div>
      <div class="topbar-actions">
        <el-button size="small" :loading="detailLoading" @click="loadAll">刷新</el-button>
        <el-button size="small" :disabled="!taskId" @click="generateReport">生成报告</el-button>
        <el-dropdown trigger="click" @command="onExportCommand">
          <el-button type="primary" size="small" :disabled="!allItems.length">
            导出 <i class="el-icon-arrow-down el-icon--right"></i>
          </el-button>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item command="csv">结果列表 CSV</el-dropdown-item>
            <el-dropdown-item command="summary">审查摘要 TXT</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </div>

    <el-tabs v-model="activeTab" class="detail-tabs list_box">
      <el-tab-pane label="概况" name="overview">
        <div v-if="isBaselineMode" v-loading="statLoading" class="stat-cards">
          <div class="stat-card">
            <div class="stat-value">{{ statData.totalRules }}</div>
            <div class="stat-label">检查项总数</div>
          </div>
          <div class="stat-card stat-pass">
            <div class="stat-value">{{ statData.passCount }}</div>
            <div class="stat-label">通过</div>
          </div>
          <div class="stat-card stat-fail">
            <div class="stat-value">{{ statData.failCount }}</div>
            <div class="stat-label">不通过</div>
          </div>
          <div class="stat-card stat-rate">
            <div class="stat-value">{{ statData.passRate }}%</div>
            <div class="stat-label">通过率</div>
          </div>
        </div>
        <div v-else-if="isVulnMode" v-loading="statLoading" class="stat-cards">
          <div class="stat-card">
            <div class="stat-value">{{ vulnStat.packages }}</div>
            <div class="stat-label">扫描包数</div>
          </div>
          <div class="stat-card stat-fail">
            <div class="stat-value">{{ vulnStat.matchedVulns }}</div>
            <div class="stat-label">漏洞总数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ vulnStat.critical }}</div>
            <div class="stat-label">严重</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ vulnStat.high }}</div>
            <div class="stat-label">高危</div>
          </div>
        </div>
        <div v-else-if="isMalwareMode" v-loading="statLoading" class="stat-cards">
          <div class="stat-card">
            <div class="stat-value">{{ malwareStat.totalFindings }}</div>
            <div class="stat-label">发现项总数</div>
          </div>
          <div class="stat-card stat-fail">
            <div class="stat-value">{{ malwareStat.critical }}</div>
            <div class="stat-label">严重</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ malwareStat.high }}</div>
            <div class="stat-label">高危</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ malwareStat.medium }}</div>
            <div class="stat-label">中危</div>
          </div>
        </div>

        <div v-loading="statLoading" class="info-section">
          <div class="section-title">任务信息</div>
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">任务批次</span>
              <span class="info-value">{{ taskId }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">任务类型</span>
              <span class="info-value">{{ kindLabel }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">检测目标数</span>
              <span class="info-value">{{ targets.length }} 个</span>
            </div>
            <div class="info-item">
              <span class="info-label">执行时间</span>
              <span class="info-value">{{ checkTime }}</span>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="检测目标" name="targets">
        <div class="targets-toolbar">
          <span class="targets-count">共 {{ targets.length }} 个目标</span>
        </div>
        <el-table v-loading="targetsLoading" :data="targets" style="width: 100%" class="myTable targets-table">
          <el-table-column type="index" label="序号" width="60" />
          <el-table-column prop="targetIp" label="目标主机" min-width="180" :show-overflow-tooltip="true" />
          <el-table-column prop="osTypeName" label="操作系统" width="120" />
          <template v-if="isBaselineMode">
            <el-table-column prop="totalRules" label="检查项数" width="100" />
            <el-table-column prop="passCount" label="通过" width="72" />
            <el-table-column prop="failCount" label="不通过" width="82" />
            <el-table-column prop="errorCount" label="异常" width="72" />
          </template>
          <template v-else-if="isVulnMode">
            <el-table-column prop="packages" label="扫描包数" width="100" />
            <el-table-column prop="matchedVulns" label="漏洞数" width="88" />
            <el-table-column prop="critical" label="严重" width="72" />
            <el-table-column prop="high" label="高危" width="72" />
            <el-table-column prop="scanStatusName" label="状态" width="88" />
          </template>
          <template v-else-if="isMalwareMode">
            <el-table-column prop="totalFindings" label="发现项数" width="100" />
            <el-table-column prop="critical" label="严重" width="72" />
            <el-table-column prop="high" label="高危" width="72" />
            <el-table-column prop="worstRiskName" label="最高风险" width="100" />
            <el-table-column prop="scanStatusName" label="状态" width="88" />
          </template>
          <el-table-column label="操作" width="180" align="right">
            <template slot-scope="scope">
              <el-link :underline="false" class="link_primary" @click="filterByTarget(scope.row)">查看结果</el-link>
              <el-link :underline="false" class="link_primary" style="margin-left: 12px" @click="viewTargetLogs(scope.row)">查看日志</el-link>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane :label="itemsTabLabel" name="items">
        <div class="items-toolbar">
          <el-select v-model="itemFilter.targetIp" placeholder="按目标筛选" size="small" clearable @change="onItemFilterChange" style="width: 160px">
            <el-option v-for="t in targets" :key="t.targetIp" :value="t.targetIp" :label="t.targetIp" />
          </el-select>
          <el-select v-if="isBaselineMode" v-model="itemFilter.result" placeholder="按结果筛选" size="small" clearable @change="onItemFilterChange" style="width: 120px; margin-left: 8px">
            <el-option value="pass" label="通过" />
            <el-option value="fail" label="不通过" />
            <el-option value="error" label="异常" />
            <el-option value="skip" label="不适配" />
          </el-select>
          <el-select v-else-if="isVulnMode" v-model="itemFilter.severity" placeholder="按风险筛选" size="small" clearable @change="onItemFilterChange" style="width: 120px; margin-left: 8px">
            <el-option :value="4" label="严重" />
            <el-option :value="3" label="高危" />
            <el-option :value="2" label="中危" />
            <el-option :value="1" label="低危" />
          </el-select>
          <template v-else-if="isMalwareMode">
            <el-select v-model="itemFilter.checkType" placeholder="按检测类型" size="small" clearable @change="onItemFilterChange" style="width: 130px; margin-left: 8px">
              <el-option :value="4" label="Webshell" />
              <el-option :value="2" label="进程检测" />
              <el-option :value="3" label="网络连接" />
              <el-option :value="5" label="Rootkit" />
              <el-option :value="6" label="敏感文件" />
              <el-option :value="1" label="文件完整性" />
            </el-select>
            <el-select v-model="itemFilter.severity" placeholder="按风险筛选" size="small" clearable @change="onItemFilterChange" style="width: 120px; margin-left: 8px">
              <el-option :value="1" label="严重" />
              <el-option :value="2" label="高危" />
              <el-option :value="3" label="中危" />
              <el-option :value="4" label="低危" />
            </el-select>
          </template>
          <span class="items-count">共 {{ filteredItems.length }} 项</span>
        </div>
        <el-table v-if="isBaselineMode" v-loading="itemsLoading" :data="pagedItems" style="width: 100%" class="myTable" max-height="560">
          <el-table-column prop="targetIp" label="目标主机" min-width="130" />
          <el-table-column prop="categoryName" label="分类" width="120" />
          <el-table-column prop="ruleName" label="检查项" min-width="180" :show-overflow-tooltip="true" />
          <el-table-column prop="resultName" label="结果" width="80" />
          <el-table-column prop="riskName" label="风险" width="80" />
          <el-table-column label="操作" width="72" align="right">
            <template slot-scope="scope">
              <el-link :underline="false" class="link_primary" @click="openItemDetail(scope.row)">详情</el-link>
            </template>
          </el-table-column>
        </el-table>
        <el-table v-else-if="isVulnMode" v-loading="itemsLoading" :data="pagedItems" style="width: 100%" class="myTable" max-height="560">
          <el-table-column prop="targetIp" label="目标主机" min-width="130" />
          <el-table-column prop="cveId" label="CVE ID" width="140" />
          <el-table-column prop="severity" label="严重程度" width="100" />
          <el-table-column prop="riskName" label="风险" width="80" />
          <el-table-column prop="packageName" label="影响软件包" min-width="140" :show-overflow-tooltip="true" />
          <el-table-column prop="packageVersion" label="版本" width="120" :show-overflow-tooltip="true" />
          <el-table-column prop="title" label="漏洞标题" min-width="200" :show-overflow-tooltip="true" />
          <el-table-column label="操作" width="72" align="right">
            <template slot-scope="scope">
              <el-link :underline="false" class="link_primary" @click="openItemDetail(scope.row)">详情</el-link>
            </template>
          </el-table-column>
        </el-table>
        <el-table v-else v-loading="itemsLoading" :data="pagedItems" style="width: 100%" class="myTable" max-height="560">
          <el-table-column prop="targetIp" label="目标主机" min-width="130" />
          <el-table-column prop="checkTypeName" label="检测类型" width="110" />
          <el-table-column prop="matchRule" label="匹配规则" min-width="160" :show-overflow-tooltip="true" />
          <el-table-column prop="riskName" label="风险" width="80" />
          <el-table-column label="文件路径" min-width="160" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.filePath">{{ scope.row.filePath }}</span>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column label="进程信息" min-width="140" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.processInfo">{{ scope.row.processInfo }}</span>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="72" align="right">
            <template slot-scope="scope">
              <el-link :underline="false" class="link_primary" @click="openItemDetail(scope.row)">详情</el-link>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-if="filteredItems.length > 0"
          background
          class="items-pager"
          layout="total, prev, pager, next, sizes, jumper"
          :total="filteredItems.length"
          :page-size="itemPageSize"
          :current-page="itemPage"
          :page-sizes="[20, 50, 100, 200]"
          @current-change="onItemPageChange"
          @size-change="onItemPageSizeChange"
        />
        <el-dialog
          :title="itemDetailTitle"
          :visible.sync="itemDetailVisible"
          width="860px"
          custom-class="theme-dialog"
          @closed="onItemDetailClosed"
        >
          <div v-if="itemDetail && isMalwareMode" class="item-detail-body">
            <div class="item-detail-head">
              <h3 class="item-detail-title">{{ itemDetail.matchRule || '恶意发现' }}</h3>
              <div class="item-detail-tags">
                <el-tag size="mini" type="danger">{{ itemDetail.riskName }}</el-tag>
                <el-tag size="mini" type="info">{{ itemDetail.checkTypeName }}</el-tag>
              </div>
            </div>
            <div class="info-grid item-detail-grid">
              <div class="info-item">
                <span class="info-label">目标主机</span>
                <span class="info-value">{{ itemDetail.targetIp || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">发现时间</span>
                <span class="info-value">{{ itemDetail.checkTime || '-' }}</span>
              </div>
              <div class="info-item full-width" v-if="itemDetail.filePath">
                <span class="info-label">文件路径</span>
                <span class="info-value">{{ itemDetail.filePath }}</span>
              </div>
              <div class="info-item full-width" v-if="itemDetail.processInfo">
                <span class="info-label">进程信息</span>
                <span class="info-value">{{ itemDetail.processInfo }}</span>
              </div>
            </div>
            <div class="item-detail-block">
              <div class="block-label">描述</div>
              <pre class="block-content">{{ itemDetail.description || '-' }}</pre>
            </div>
            <div class="item-detail-block">
              <div class="block-label">修复建议</div>
              <pre class="block-content">{{ itemDetail.fixSuggestion || '-' }}</pre>
            </div>
          </div>
          <div v-else-if="itemDetail && isVulnMode" class="item-detail-body">
            <div class="item-detail-head">
              <h3 class="item-detail-title">{{ itemDetail.cveId }}</h3>
              <div class="item-detail-tags">
                <el-tag size="mini" type="danger">{{ itemDetail.severity || itemDetail.riskName }}</el-tag>
              </div>
            </div>
            <div class="info-grid item-detail-grid">
              <div class="info-item">
                <span class="info-label">目标主机</span>
                <span class="info-value">{{ itemDetail.targetIp || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">影响软件包</span>
                <span class="info-value">{{ itemDetail.packageName || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">软件版本</span>
                <span class="info-value">{{ itemDetail.packageVersion || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">发现时间</span>
                <span class="info-value">{{ itemDetail.checkTime || '-' }}</span>
              </div>
            </div>
            <div class="item-detail-block">
              <div class="block-label">漏洞标题</div>
              <pre class="block-content">{{ itemDetail.title || '-' }}</pre>
            </div>
          </div>
          <div v-else-if="itemDetail && isBaselineMode" class="item-detail-body">
            <div class="item-detail-head">
              <h3 class="item-detail-title">{{ itemDetail.ruleName }}</h3>
              <div class="item-detail-tags">
                <el-tag size="mini" :type="itemResultTagType(itemDetail.checkResult)">{{ itemDetail.resultName }}</el-tag>
                <el-tag size="mini" type="info">{{ itemDetail.riskName }}</el-tag>
              </div>
            </div>
            <div class="info-grid item-detail-grid">
              <div class="info-item">
                <span class="info-label">目标主机</span>
                <span class="info-value">{{ itemDetail.targetIp || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">分类</span>
                <span class="info-value">{{ itemDetail.categoryName || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">检查时间</span>
                <span class="info-value">{{ itemDetail.checkTime || '-' }}</span>
              </div>
            </div>
            <div class="item-detail-block">
              <div class="block-label">检查命令</div>
              <pre class="block-content">{{ itemDetail.checkCommand || '-' }}</pre>
            </div>
            <div class="item-detail-block">
              <div class="block-label">期望值</div>
              <pre class="block-content">{{ itemDetail.expectedValue || '-' }}</pre>
            </div>
            <div class="item-detail-block">
              <div class="block-label">实际值</div>
              <pre class="block-content">{{ itemDetail.actualValue || '-' }}</pre>
            </div>
            <div class="item-detail-block">
              <div class="block-label">修复建议</div>
              <pre class="block-content">{{ itemDetail.fixSuggestion || '-' }}</pre>
            </div>
            <div v-if="itemDetail.riskDescription" class="item-detail-block">
              <div class="block-label">风险说明</div>
              <pre class="block-content">{{ itemDetail.riskDescription }}</pre>
            </div>
          </div>
          <span slot="footer">
            <el-button @click="itemDetailVisible = false">关闭</el-button>
          </span>
        </el-dialog>
      </el-tab-pane>

      <el-tab-pane label="日志" name="logs">
        <div v-if="logsLoading" v-loading="true" class="logs-loading" />
        <div v-else-if="targets.length === 0" class="logs-empty">暂无日志信息</div>
        <template v-else>
          <div class="logs-toolbar">
            <el-input
              v-model="logSearch"
              placeholder="搜索目标 IP"
              size="small"
              clearable
              class="log-search-input"
              @keydown.enter.native="onLogSearch"
            />
            <el-button type="primary" size="small" @click="onLogSearch">搜索</el-button>
            <el-button size="small" @click="resetLogSearch">重置</el-button>
            <span class="logs-count">共 {{ filteredTargetLogRows.length }} 个目标日志</span>
          </div>
          <el-table v-loading="logsLoading" :data="filteredTargetLogRows" style="width: 100%" class="myTable" max-height="520">
            <el-table-column type="index" label="序号" width="60" />
            <el-table-column prop="targetIp" label="目标主机" min-width="160" :show-overflow-tooltip="true" />
            <el-table-column prop="osTypeName" label="操作系统" width="120" />
            <template v-if="isBaselineMode">
              <el-table-column prop="totalRules" label="检查项" width="88" />
              <el-table-column prop="passCount" label="通过" width="72" />
              <el-table-column prop="failCount" label="不通过" width="82" />
              <el-table-column prop="errorCount" label="异常" width="72" />
            </template>
            <template v-else-if="isVulnMode">
              <el-table-column prop="packages" label="扫描包数" width="100" />
              <el-table-column prop="matchedVulns" label="漏洞数" width="88" />
              <el-table-column prop="critical" label="严重" width="72" />
              <el-table-column prop="high" label="高危" width="72" />
            </template>
            <template v-else-if="isMalwareMode">
              <el-table-column prop="totalFindings" label="发现项" width="88" />
              <el-table-column prop="critical" label="严重" width="72" />
              <el-table-column prop="high" label="高危" width="72" />
            </template>
            <el-table-column prop="logCount" label="日志条数" width="96" />
            <el-table-column label="操作" width="120" align="right">
              <template slot-scope="scope">
                <el-link :underline="false" class="link_primary" @click="openTargetLogDetail(scope.row)">日志详情</el-link>
              </template>
            </el-table-column>
          </el-table>
        </template>

        <el-dialog
          :title="logDetailTitle"
          :visible.sync="logDetailVisible"
          width="900px"
          custom-class="theme-dialog"
          @closed="onLogDetailClosed"
        >
          <div class="log-detail-toolbar">
            <el-select v-model="logDetailLevel" placeholder="按级别筛选" size="small" clearable style="width: 120px">
              <el-option value="info" label="INFO" />
              <el-option value="warn" label="WARN" />
              <el-option value="error" label="ERROR" />
            </el-select>
            <span class="logs-count">共 {{ detailLogLines.length }} 条</span>
          </div>
          <div v-if="detailLogLines.length === 0" class="logs-empty">该目标暂无日志明细</div>
          <div v-else class="logs-list log-detail-list">
            <div v-for="(log, idx) in detailLogLines" :key="idx" class="log-item">
              <span class="log-time">{{ log.time }}</span>
              <span :class="['log-level', 'level-' + log.level]">{{ log.levelLabel }}</span>
              <span class="log-msg">{{ log.message }}</span>
            </div>
          </div>
          <span slot="footer">
            <el-button @click="logDetailVisible = false">关闭</el-button>
          </span>
        </el-dialog>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import security from '@/api/security.js'
import { buildHostSummaryText, downloadTextFile, exportHostItemsCsv } from './utils/hostTaskExport.js'

export default {
  name: 'HostTaskDetail',
  data() {
    return {
      taskId: 0,
      kindLabel: '安全配置核查',
      taskSource: 'baseline',
      checkTime: '',
      activeTab: 'overview',
      statLoading: false,
      statData: {
        totalRules: 0,
        passCount: 0,
        failCount: 0,
        passRate: 0
      },
      vulnStat: {
        packages: 0,
        matchedVulns: 0,
        critical: 0,
        high: 0,
        medium: 0,
        low: 0
      },
      malwareStat: {
        totalFindings: 0,
        critical: 0,
        high: 0,
        medium: 0,
        low: 0
      },
      targetsLoading: false,
      targets: [],
      itemsLoading: false,
      allItems: [],
      itemFilter: {
        targetIp: '',
        result: '',
        severity: '',
        checkType: ''
      },
      itemPage: 1,
      itemPageSize: 50,
      itemDetailVisible: false,
      itemDetail: null,
      logsByTarget: {},
      logsLoading: false,
      logSearch: '',
      logSearchApplied: '',
      logDetailVisible: false,
      logDetailTitle: '',
      logDetailTargetIp: '',
      logDetailLevel: ''
    }
  },
  computed: {
    isBaselineMode() {
      return this.taskSource === 'baseline'
    },
    isVulnMode() {
      return this.taskSource === 'vuln'
    },
    isMalwareMode() {
      return this.taskSource === 'malware'
    },
    itemsTabLabel() {
      if (this.isVulnMode) return '漏洞列表'
      if (this.isMalwareMode) return '发现项'
      return '核查项'
    },
    itemDetailTitle() {
      if (this.isVulnMode) return '漏洞详情'
      if (this.isMalwareMode) return '恶意发现详情'
      return '核查项详情'
    },
    detailLoading() {
      return this.statLoading || this.targetsLoading || this.itemsLoading || this.logsLoading
    },
    targetLogRows() {
      return this.targets.map(target => ({
        ...target,
        logCount: (this.logsByTarget[target.targetIp] || []).length
      }))
    },
    filteredTargetLogRows() {
      if (!this.logSearchApplied) return this.targetLogRows
      const keyword = this.logSearchApplied.toLowerCase()
      return this.targetLogRows.filter(row => (row.targetIp || '').toLowerCase().includes(keyword))
    },
    detailLogLines() {
      if (!this.logDetailTargetIp) return []
      let list = this.logsByTarget[this.logDetailTargetIp] || []
      if (this.logDetailLevel) {
        list = list.filter(log => log.level === this.logDetailLevel)
      }
      return list
    },
    filteredItems() {
      let list = this.allItems
      if (this.itemFilter.targetIp) {
        list = list.filter(i => i.targetIp === this.itemFilter.targetIp)
      }
      if (this.isVulnMode) {
        if (this.itemFilter.severity !== '' && this.itemFilter.severity !== null && this.itemFilter.severity !== undefined) {
          list = list.filter(i => i.riskLevel === this.itemFilter.severity)
        }
        return list
      }
      if (this.isMalwareMode) {
        if (this.itemFilter.severity !== '' && this.itemFilter.severity !== null && this.itemFilter.severity !== undefined) {
          list = list.filter(i => i.riskLevel === this.itemFilter.severity)
        }
        if (this.itemFilter.checkType !== '' && this.itemFilter.checkType !== null && this.itemFilter.checkType !== undefined) {
          list = list.filter(i => i.checkType === this.itemFilter.checkType)
        }
        return list
      }
      if (this.itemFilter.result === 'pass') {
        list = list.filter(i => i.checkResult === 1)
      } else if (this.itemFilter.result === 'fail') {
        list = list.filter(i => i.checkResult === 2)
      } else if (this.itemFilter.result === 'error') {
        list = list.filter(i => i.checkResult === 3)
      } else if (this.itemFilter.result === 'skip') {
        list = list.filter(i => i.checkResult === 4)
      }
      return list
    },
    pagedItems() {
      const start = (this.itemPage - 1) * this.itemPageSize
      return this.filteredItems.slice(start, start + this.itemPageSize)
    },
    taskModeLabel() {
      if (this.isVulnMode) return 'vuln'
      if (this.isMalwareMode) return 'malware'
      return 'baseline'
    },
    auditSummaryPreview() {
      return buildHostSummaryText({
        taskId: this.taskId,
        kindLabel: this.kindLabel,
        mode: this.taskModeLabel,
        checkTime: this.checkTime,
        targetCount: this.targets.length,
        itemCount: this.allItems.length,
        statData: this.statData,
        vulnStat: this.vulnStat,
        malwareStat: this.malwareStat
      })
    }
  },
  created() {
    this.taskId = parseInt(this.$route.query.taskId) || 0
    this.kindLabel = this.$route.query.kindLabel || '安全配置核查'
    this.taskSource = this.$route.query.source || this.detectTaskSource(this.kindLabel)
    this.checkTime = this.$route.query.checkTime || ''
    if (this.taskId > 0) {
      this.loadAll()
    }
  },
  methods: {
    detectTaskSource(kindLabel) {
      const label = kindLabel || ''
      if (label.indexOf('漏洞') >= 0) return 'vuln'
      if (label.indexOf('恶意') >= 0) return 'malware'
      return 'baseline'
    },
    async loadAll() {
      this.logsLoading = true
      try {
        if (this.isVulnMode) {
          await Promise.all([this.loadVulnStat(), this.loadVulnTargets(), this.loadVulnFindings()])
          this.buildVulnLogs()
        } else if (this.isMalwareMode) {
          await Promise.all([this.loadMalwareStat(), this.loadMalwareTargets(), this.loadMalwareFindings()])
          this.buildMalwareLogs()
        } else {
          await Promise.all([this.loadStat(), this.loadTargets(), this.loadItems()])
          this.buildLogs()
        }
      } finally {
        this.logsLoading = false
      }
    },
    async generateReport() {
      const module = this.isVulnMode ? 'host' : this.isMalwareMode ? 'host' : 'host'
      const taskId = this.taskId
      if (!taskId) return
      try {
        const res = await security.generateSecurityReport({ module, taskId })
        if (res.code === 200 && res.data) {
          this.$message({ message: '报告已生成，前往报告中心查看', type: 'success' })
        } else {
          this.$message({ message: res.msg || '生成报告失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '生成报告失败: ' + (e.message || ''), type: 'error' })
      }
    },
    async loadVulnStat() {
      this.statLoading = true
      try {
        const res = await security.getHostVulnStat({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.vulnStat = {
            packages: res.data.packages || 0,
            matchedVulns: res.data.matchedVulns || 0,
            critical: res.data.critical || 0,
            high: res.data.high || 0,
            medium: res.data.medium || 0,
            low: res.data.low || 0
          }
        }
      } finally {
        this.statLoading = false
      }
    },
    async loadVulnTargets() {
      this.targetsLoading = true
      try {
        const res = await security.getHostVulnTargets({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.targets = res.data.list || []
        }
      } finally {
        this.targetsLoading = false
      }
    },
    async loadVulnFindings() {
      this.itemsLoading = true
      try {
        const res = await security.getHostVulnFindings({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.allItems = res.data.list || []
        }
      } finally {
        this.itemsLoading = false
      }
    },
    async loadMalwareStat() {
      this.statLoading = true
      try {
        const res = await security.getHostMalwareStat({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.malwareStat = {
            totalFindings: res.data.totalFindings || 0,
            critical: res.data.critical || 0,
            high: res.data.high || 0,
            medium: res.data.medium || 0,
            low: res.data.low || 0
          }
        }
      } finally {
        this.statLoading = false
      }
    },
    async loadMalwareTargets() {
      this.targetsLoading = true
      try {
        const res = await security.getHostMalwareTargets({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.targets = res.data.list || []
        }
      } finally {
        this.targetsLoading = false
      }
    },
    async loadMalwareFindings() {
      this.itemsLoading = true
      try {
        const res = await security.getHostMalwareFindings({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.allItems = res.data.list || []
        }
      } finally {
        this.itemsLoading = false
      }
    },
    buildMalwareLogs() {
      const logsByTarget = {}
      const baseTime = this.checkTime || '-'
      for (const target of this.targets) {
        const lines = []
        lines.push({
          time: target.checkTime || baseTime,
          level: 'info',
          levelLabel: 'INFO',
          message: `开始 YARA 恶意代码扫描：${target.targetIp}（${target.osTypeName || '未知系统'}）`
        })
        if (target.scanStatus === 2) {
          lines.push({
            time: target.checkTime || baseTime,
            level: 'error',
            levelLabel: 'ERROR',
            message: `扫描失败：${target.errorMessage || '未知错误'}`
          })
        } else {
          lines.push({
            time: target.checkTime || baseTime,
            level: target.totalFindings > 0 ? 'warn' : 'info',
            levelLabel: target.totalFindings > 0 ? 'WARN' : 'INFO',
            message: `扫描完成：发现 ${target.totalFindings || 0} 项（严重 ${target.critical || 0} / 高危 ${target.high || 0}）`
          })
          const findings = this.allItems.filter(i => i.targetIp === target.targetIp)
          for (const f of findings.slice(0, 20)) {
            lines.push({
              time: f.checkTime || baseTime,
              level: 'warn',
              levelLabel: 'WARN',
              message: `${f.matchRule} · ${f.checkTypeName}${f.filePath ? ' · ' + f.filePath : ''}`
            })
          }
          if (findings.length > 20) {
            lines.push({
              time: target.checkTime || baseTime,
              level: 'info',
              levelLabel: 'INFO',
              message: `另有 ${findings.length - 20} 条发现项，请见「发现项」`
            })
          }
        }
        logsByTarget[target.targetIp] = lines
      }
      this.logsByTarget = logsByTarget
    },
    buildVulnLogs() {
      const logsByTarget = {}
      const baseTime = this.checkTime || '-'
      for (const target of this.targets) {
        const lines = []
        lines.push({
          time: target.checkTime || baseTime,
          level: 'info',
          levelLabel: 'INFO',
          message: `开始 CVE 漏洞扫描：${target.targetIp}（${target.osTypeName || '未知系统'}）`
        })
        if (target.scanStatus === 2) {
          lines.push({
            time: target.checkTime || baseTime,
            level: 'error',
            levelLabel: 'ERROR',
            message: `扫描失败：${target.errorMessage || '未知错误'}`
          })
        } else {
          lines.push({
            time: target.checkTime || baseTime,
            level: target.matchedVulns > 0 ? 'warn' : 'info',
            levelLabel: target.matchedVulns > 0 ? 'WARN' : 'INFO',
            message: `扫描完成：包 ${target.packages || 0} 个，漏洞 ${target.matchedVulns || 0} 个（严重 ${target.critical || 0} / 高危 ${target.high || 0}）`
          })
          const findings = this.allItems.filter(i => i.targetIp === target.targetIp)
          for (const f of findings.slice(0, 20)) {
            lines.push({
              time: f.checkTime || baseTime,
              level: 'warn',
              levelLabel: 'WARN',
              message: `${f.cveId} · ${f.packageName} ${f.packageVersion}`
            })
          }
          if (findings.length > 20) {
            lines.push({
              time: target.checkTime || baseTime,
              level: 'info',
              levelLabel: 'INFO',
              message: `另有 ${findings.length - 20} 条漏洞，请见「漏洞列表」`
            })
          }
        }
        logsByTarget[target.targetIp] = lines
      }
      this.logsByTarget = logsByTarget
    },
    buildLogs() {
      const baseTime = this.checkTime || this.earliestCheckTime()
      const logsByTarget = {}

      if (this.targets.length === 0 && this.allItems.length === 0) {
        this.logsByTarget = {}
        return
      }

      for (const target of this.targets) {
        const lines = []
        lines.push({
          time: baseTime,
          level: 'info',
          levelLabel: 'INFO',
          message: `开始核查目标 ${target.targetIp}（${target.osTypeName || '未知系统'}）`
        })
        lines.push({
          time: baseTime,
          level: 'info',
          levelLabel: 'INFO',
          message: `核查完成：检查项 ${target.totalRules}，通过 ${target.passCount}，不通过 ${target.failCount}，异常 ${target.errorCount || 0}`
        })

        const failedItems = this.allItems.filter(i => i.targetIp === target.targetIp && i.checkResult === 2)
        for (const item of failedItems) {
          lines.push({
            time: item.checkTime || baseTime,
            level: 'warn',
            levelLabel: 'WARN',
            message: `${item.ruleName} 不通过：期望 ${item.expectedValue || '-'}，实际 ${item.actualValue || '-'}`
          })
        }

        const errorItems = this.allItems.filter(i => i.targetIp === target.targetIp && i.checkResult === 3)
        for (const item of errorItems) {
          lines.push({
            time: item.checkTime || baseTime,
            level: 'error',
            levelLabel: 'ERROR',
            message: `${item.ruleName} 执行异常：${item.actualValue || '未知错误'}`
          })
        }

        const issueCount = failedItems.length + errorItems.length
        lines.push({
          time: baseTime,
          level: issueCount > 0 ? 'warn' : 'info',
          levelLabel: issueCount > 0 ? 'WARN' : 'INFO',
          message: `目标核查结束，共 ${target.totalRules} 项，问题 ${issueCount} 项`
        })

        logsByTarget[target.targetIp] = lines
      }

      this.logsByTarget = logsByTarget
    },
    earliestCheckTime() {
      if (!this.allItems.length) return '-'
      return this.allItems.reduce((earliest, item) => {
        if (!item.checkTime) return earliest
        if (!earliest || item.checkTime < earliest) return item.checkTime
        return earliest
      }, '')
    },
    async loadStat() {
      this.statLoading = true
      try {
        const res = await security.getBaselineStat({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.statData = {
            totalRules: res.data.totalRules || 0,
            passCount: res.data.passCount || 0,
            failCount: res.data.failCount || 0,
            passRate: (res.data.passRate || 0).toFixed(1)
          }
        }
      } finally {
        this.statLoading = false
      }
    },
    async loadTargets() {
      this.targetsLoading = true
      try {
        const res = await security.getBaselineTaskTargets({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.targets = res.data.list || []
        }
      } finally {
        this.targetsLoading = false
      }
    },
    async loadItems() {
      this.itemsLoading = true
      try {
        const res = await security.getBaselineList({ taskId: this.taskId })
        if (res.code === 200 && res.data) {
          this.allItems = res.data.list || []
        }
      } finally {
        this.itemsLoading = false
      }
    },
    onExportCommand(cmd) {
      if (cmd === 'csv') {
        this.exportCsv()
      } else if (cmd === 'summary') {
        this.exportSummary()
      }
    },
    exportCsv() {
      if (!this.allItems.length) {
        this.$message.warning('暂无可导出的结果')
        return
      }
      exportHostItemsCsv(this.taskId, this.taskModeLabel, this.filteredItems)
      this.$message.success('结果列表已导出')
    },
    exportSummary() {
      downloadTextFile(`hostsec-${this.taskId}-summary.txt`, this.auditSummaryPreview)
      this.$message.success('审查摘要已导出')
    },
    filterByTarget(row) {
      this.itemFilter.targetIp = row.targetIp
      this.itemFilter.result = ''
      this.itemFilter.severity = ''
      this.itemFilter.checkType = ''
      this.itemPage = 1
      this.activeTab = 'items'
    },
    viewTargetLogs(row) {
      this.activeTab = 'logs'
      this.openTargetLogDetail(row)
    },
    openTargetLogDetail(row) {
      this.logDetailTargetIp = row.targetIp
      this.logDetailTitle = `${row.targetIp} · 执行日志`
      this.logDetailLevel = ''
      this.logDetailVisible = true
    },
    onLogDetailClosed() {
      this.logDetailTargetIp = ''
      this.logDetailLevel = ''
    },
    onLogSearch() {
      this.logSearchApplied = (this.logSearch || '').trim()
    },
    resetLogSearch() {
      this.logSearch = ''
      this.logSearchApplied = ''
    },
    onItemFilterChange() {
      this.itemPage = 1
    },
    onItemPageChange(page) {
      this.itemPage = page
    },
    onItemPageSizeChange(size) {
      this.itemPageSize = size
      this.itemPage = 1
    },
    openItemDetail(row) {
      this.itemDetail = row
      this.itemDetailVisible = true
    },
    onItemDetailClosed() {
      this.itemDetail = null
    },
    itemResultTagType(checkResult) {
      if (checkResult === 1) return 'success'
      if (checkResult === 3) return 'warning'
      if (checkResult === 4) return 'info'
      return 'danger'
    },
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';
@import './css/appsec-task-detail.less';

.task-detail-page {
  .title-sub {
    font-size: 14px;
    color: #64748b;
    font-weight: 400;
    margin-left: 6px;
  }

  .detail-tabs {
    ::v-deep .el-table {
      width: 100% !important;
    }

    ::v-deep .el-table__header-wrapper table,
    ::v-deep .el-table__body-wrapper table {
      width: 100% !important;
    }
  }
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.stat-cards {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}
.stat-card {
  flex: 1;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 8px;
  padding: 20px;
  text-align: center;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #e2e8f0;
  line-height: 1.2;
}
.stat-label {
  font-size: 13px;
  color: #94a3b8;
  margin-top: 6px;
}
.stat-pass .stat-value { color: #34d399; }
.stat-fail .stat-value { color: #f87171; }
.stat-rate .stat-value { color: #60a5fa; }

.info-section {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 8px;
  padding: 20px;
}
.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #e2e8f0;
  margin-bottom: 16px;
}
.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.info-label {
  font-size: 12px;
  color: #64748b;
}
.info-value {
  font-size: 14px;
  color: #e2e8f0;
}

.targets-toolbar, .items-toolbar, .logs-toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
  gap: 8px;
}
.log-search-input {
  width: 220px;
}
.log-detail-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}
.targets-count, .items-count, .logs-count {
  font-size: 12px;
  color: #94a3b8;
  margin-left: auto;
}
.items-pager {
  margin-top: 12px;
  text-align: right;
}

.item-detail-body {
  max-height: 70vh;
  overflow-y: auto;
}

.item-detail-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.item-detail-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #e2e8f0;
  line-height: 1.4;
  flex: 1;
}

.item-detail-tags {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.item-detail-grid {
  margin-bottom: 16px;
}

.item-detail-block {
  margin-bottom: 14px;

  .block-label {
    font-size: 12px;
    color: #64748b;
    margin-bottom: 6px;
  }

  .block-content {
    margin: 0;
    padding: 10px 12px;
    background: rgba(0, 0, 0, 0.25);
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 6px;
    font-size: 12px;
    line-height: 1.6;
    color: #94a3b8;
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 200px;
    overflow-y: auto;
    font-family: 'Consolas', 'Courier New', monospace;
  }
}

.logs-list {
  background: rgba(0,0,0,0.2);
  border-radius: 8px;
  padding: 12px;
  max-height: 560px;
  overflow-y: auto;
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 12px;
}
.log-item {
  display: flex;
  gap: 12px;
  padding: 6px 0;
  border-bottom: 1px solid rgba(255,255,255,0.04);
  &:last-child { border-bottom: none; }
}
.log-time {
  color: #64748b;
  white-space: nowrap;
  min-width: 140px;
}
.log-level {
  min-width: 50px;
  font-weight: 600;
  &.level-info { color: #60a5fa; }
  &.level-warn { color: #fbbf24; }
  &.level-error { color: #f87171; }
}
.log-detail-list {
  max-height: 480px;
}
.log-msg {
  color: #94a3b8;
  word-break: break-all;
}
.logs-empty {
  text-align: center;
  padding: 40px;
  color: #64748b;
  font-size: 14px;
}
.logs-loading {
  min-height: 120px;
}
</style>
