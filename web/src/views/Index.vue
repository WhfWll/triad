<template>
  <div v-loading="loading" class="dashboard security-container">
    <div class="welcome-bar">
      <div>
        <h2 class="welcome-title">{{ productName }} · 安全运营概览</h2>
        <p class="welcome-desc">汇总主机 / 应用 / 数据三大安全检查模块的任务执行与规则库现状</p>
      </div>
      <el-button type="primary" size="small" :loading="loading" @click="loadDashboard">刷新</el-button>
    </div>

    <div class="module-row">
      <router-link
        v-for="mod in modules"
        :key="mod.key"
        :to="mod.link"
        class="module-card"
        :class="'mod-' + mod.key"
      >
        <div class="mod-icon"><i :class="mod.icon"></i></div>
        <div class="mod-body">
          <span class="mod-label">{{ mod.label }}</span>
          <span class="mod-value">{{ mod.taskCount }}</span>
          <span class="mod-sub">{{ mod.subText }}</span>
        </div>
      </router-link>
    </div>

    <div class="rules-row">
      <div v-for="rule in ruleStats" :key="rule.key" class="rule-chip">
        <span class="rule-name">{{ rule.label }}</span>
        <span class="rule-count">{{ rule.count }}</span>
      </div>
    </div>

    <div class="charts-grid">
      <div class="chart-card">
        <div class="chart-hd"><h3>主机核查结果（最近任务）</h3></div>
        <div class="chart-bd"><div id="hostPie" class="chart-box"></div></div>
      </div>
      <div class="chart-card">
        <div class="chart-hd"><h3>各模块任务数量</h3></div>
        <div class="chart-bd"><div id="moduleBar" class="chart-box"></div></div>
      </div>
    </div>

    <div class="chart-card recent-card">
      <div class="chart-hd">
        <h3>最近检查记录</h3>
        <span class="recent-hint">合并展示各模块最近一批任务</span>
      </div>
      <div class="chart-bd table-bd">
        <el-table :data="recentTasks" style="width: 100%" class="myTable" empty-text="暂无检查记录，请从各模块发起扫描任务">
          <el-table-column prop="moduleLabel" label="模块" width="100" />
          <el-table-column prop="kindLabel" label="类型" width="130" />
          <el-table-column prop="target" label="目标" min-width="180" :show-overflow-tooltip="true" />
          <el-table-column prop="summary" label="结果摘要" min-width="160" :show-overflow-tooltip="true" />
          <el-table-column prop="checkTime" label="时间" width="170" />
          <el-table-column label="操作" width="90" align="right">
            <template slot-scope="scope">
              <el-link :underline="false" class="link_primary" @click.prevent="goTask(scope.row)">查看</el-link>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <div class="quick-row">
      <router-link v-for="q in quickLinks" :key="q.to" :to="q.to" class="quick-item">
        <i :class="q.icon"></i>
        <span>{{ q.label }}</span>
      </router-link>
    </div>
  </div>
</template>

<script>
var echarts = require('echarts')
import security from '@/api/security.js'
import { PRODUCT_NAME } from '@/config/product.js'

const CHART_THEME = {
  text: '#94a3b8',
  axis: '#64748b',
  split: 'rgba(255,255,255,0.08)',
  primary: '#00d4aa',
  pass: '#34d399',
  fail: '#f87171',
  warn: '#fbbf24',
  info: '#60a5fa',
}

export default {
  name: 'Index',
  data() {
    return {
      productName: PRODUCT_NAME,
      loading: false,
      modules: [
        { key: 'host', label: '主机安全', icon: 'el-icon-monitor', link: '/hostsec/tasks', taskCount: 0, subText: '任务 0 · 不通过 0' },
        { key: 'app', label: '应用安全', icon: 'el-icon-mobile-phone', link: '/appsec/tasks', taskCount: 0, subText: '任务 0 · 漏洞 0' },
        { key: 'data', label: '数据安全', icon: 'el-icon-document', link: '/datasec/tasks', taskCount: 0, subText: '任务 0 · 目标库 0' },
      ],
      ruleStats: [
        { key: 'hostRule', label: '主机核查规则', count: '—' },
        { key: 'cve', label: 'CVE 漏洞库', count: '—' },
        { key: 'malware', label: '病毒库规则', count: '—' },
        { key: 'datasec', label: '数据检测规则', count: '—' },
      ],
      hostCheckStat: { pass: 0, fail: 0, error: 0 },
      moduleBarData: [],
      recentTasks: [],
      chartHost: null,
      chartModule: null,
    }
  },
  computed: {
    quickLinks() {
      return [
        { to: '/hostsec/tasks', icon: 'el-icon-monitor', label: '新建主机检查' },
        { to: '/appsec/task/new', icon: 'el-icon-mobile-phone', label: '新建应用扫描' },
        { to: '/datasec/tasks', icon: 'el-icon-document', label: '新建数据检查' },
        { to: '/hostsec/rules', icon: 'el-icon-setting', label: '主机规则' },
        { to: '/datasec/rules', icon: 'el-icon-tickets', label: '数据规则' },
        { to: '/systemsetting', icon: 'el-icon-s-tools', label: '系统配置' },
      ]
    },
  },
  mounted() {
    this.loadDashboard()
    window.addEventListener('resize', this.handleResize)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.handleResize)
    if (this.chartHost) this.chartHost.dispose()
    if (this.chartModule) this.chartModule.dispose()
  },
  methods: {
    parseTime(str) {
      if (!str) return 0
      return new Date(String(str).replace(/-/g, '/')).getTime() || 0
    },
    pickTotal(res) {
      if (res && res.code === 200 && res.data) return res.data.total || 0
      return 0
    },
    async loadDashboard() {
      this.loading = true
      try {
        const [
          baselineRes,
          yaraRes,
          dynRes,
          appRes,
          dbRes,
          hostRulesRes,
          cveRes,
          malwareRulesRes,
          datasecRulesRes,
          targetsRes,
        ] = await Promise.all([
          security.getBaselineTaskList({ page: 1, size: 40, scanScene: 0 }).catch(() => null),
          security.getYaraTaskList({ page: 1, size: 15 }).catch(() => null),
          security.getDynamicScanList({ page: 1, size: 15 }).catch(() => null),
          security.getAppSpecificList({ page: 1, size: 15 }).catch(() => null),
          security.getDBCheckList({ page: 1, size: 15 }).catch(() => null),
          security.getBaselineRulesStats().catch(() => null),
          security.getCveDBInfo().catch(() => null),
          security.getMalwareRuleList({ page: 1, size: 1 }).catch(() => null),
          security.getDatasecRulesStats().catch(() => null),
          security.getDatasecTargetList({ page: 1, size: 1 }).catch(() => null),
        ])

        const baselineTotal = this.pickTotal(baselineRes)
        const yaraTotal = this.pickTotal(yaraRes)
        const dynTotal = this.pickTotal(dynRes)
        const appTotal = this.pickTotal(appRes)
        const dbTotal = this.pickTotal(dbRes)
        const targetTotal = this.pickTotal(targetsRes)

        let pass = 0
        let fail = 0
        let error = 0
        const baselineList = (baselineRes && baselineRes.data && baselineRes.data.list) || []
        baselineList.forEach((r) => {
          pass += r.passCount || 0
          fail += r.failCount || 0
          error += r.errorCount || 0
        })
        this.hostCheckStat = { pass, fail, error }

        let appVuln = 0
        const dynList = (dynRes && dynRes.data && dynRes.data.list) || []
        const appList = (appRes && appRes.data && appRes.data.list) || []
        ;[...dynList, ...appList].forEach((r) => {
          appVuln += r.vulnCount || 0
        })

        this.modules = [
          {
            key: 'host',
            label: '主机安全',
            icon: 'el-icon-monitor',
            link: '/hostsec/tasks',
            taskCount: baselineTotal + yaraTotal,
            subText: `基线 ${baselineTotal} · 恶意代码 ${yaraTotal} · 不通过 ${fail}`,
          },
          {
            key: 'app',
            label: '应用安全',
            icon: 'el-icon-mobile-phone',
            link: '/appsec/tasks',
            taskCount: dynTotal + appTotal,
            subText: `动态 ${dynTotal} · 专项 ${appTotal} · 漏洞 ${appVuln}`,
          },
          {
            key: 'data',
            label: '数据安全',
            icon: 'el-icon-document',
            link: '/datasec/tasks',
            taskCount: dbTotal,
            subText: `数据库任务 ${dbTotal} · 目标库 ${targetTotal}`,
          },
        ]

        const hostRuleTotal = (hostRulesRes && hostRulesRes.data && hostRulesRes.data.total) || 0
        const cveTotal = (cveRes && cveRes.data && cveRes.data.totalRecords) || 0
        const malwareTotal = this.pickTotal(malwareRulesRes)
        const datasecTotal = (datasecRulesRes && datasecRulesRes.data && (datasecRulesRes.data.enabledTotal ?? datasecRulesRes.data.total)) || 0

        this.ruleStats = [
          { key: 'hostRule', label: '主机核查规则', count: this.fmtNum(hostRuleTotal) },
          { key: 'cve', label: 'CVE 漏洞库', count: this.fmtNum(cveTotal) },
          { key: 'malware', label: '病毒库规则', count: this.fmtNum(malwareTotal) },
          { key: 'datasec', label: '数据检测规则', count: this.fmtNum(datasecTotal) },
        ]

        this.moduleBarData = [
          { name: '主机基线', value: baselineTotal },
          { name: '恶意代码', value: yaraTotal },
          { name: '动态扫描', value: dynTotal },
          { name: '专项应用', value: appTotal },
          { name: '数据库', value: dbTotal },
        ]

        this.recentTasks = this.buildRecentTasks(baselineList, yaraRes, dynList, appList, dbRes)
        this.$nextTick(() => this.renderCharts())
      } finally {
        this.loading = false
      }
    },
    fmtNum(n) {
      const v = Number(n)
      if (!v) return '0'
      return v >= 1000 ? v.toLocaleString() : String(v)
    },
    buildRecentTasks(baselineList, yaraRes, dynList, appList, dbRes) {
      const rows = []
      baselineList.forEach((r) => {
        rows.push({
          module: 'host',
          moduleLabel: '主机安全',
          kindLabel: r.scanSceneName || '安全配置核查',
          target: r.targetIp,
          summary: `通过 ${r.passCount || 0} / 不通过 ${r.failCount || 0}`,
          checkTime: r.checkTime,
          _ts: this.parseTime(r.checkTime),
          route: {
            path: '/hostsec/task-detail',
            query: {
              taskId: r.taskId,
              kindLabel: r.scanSceneName || '安全配置核查',
              checkTime: r.checkTime || '',
            },
          },
        })
      })
      const yaraList = (yaraRes && yaraRes.data && yaraRes.data.list) || []
      yaraList.forEach((r) => {
        rows.push({
          module: 'host',
          moduleLabel: '主机安全',
          kindLabel: '恶意代码检测',
          target: r.targetIp,
          summary: `发现 ${r.totalFindings || 0} 项 · ${r.worstRiskName || '—'}`,
          checkTime: r.checkTime,
          _ts: this.parseTime(r.checkTime),
          route: {
            path: '/hostsec/task-detail',
            query: {
              taskId: r.taskId,
              kindLabel: '恶意代码检测',
              checkTime: r.checkTime || '',
              source: 'malware',
            },
          },
        })
      })
      dynList.forEach((r) => {
        rows.push({
          module: 'app',
          moduleLabel: '应用安全',
          kindLabel: '动态扫描',
          target: r.targetSummary || r.name || '—',
          summary: `漏洞 ${r.vulnCount || 0} · 高危 ${r.highRiskCount || 0}`,
          checkTime: r.scanTime || r.createTime || '—',
          _ts: this.parseTime(r.scanTime || r.createTime),
          route: { path: '/appsec/task/detail', query: { id: r.id, type: 'dyn' } },
        })
      })
      appList.forEach((r) => {
        rows.push({
          module: 'app',
          moduleLabel: '应用安全',
          kindLabel: '专项应用检测',
          target: r.targetSummary || r.name || '—',
          summary: `漏洞 ${r.vulnCount || 0} · 高危 ${r.highRiskCount || 0}`,
          checkTime: r.scanTime || r.createTime || '—',
          _ts: this.parseTime(r.scanTime || r.createTime),
          route: { path: '/appsec/task/detail', query: { id: r.id, type: 'app' } },
        })
      })
      const dbList = (dbRes && dbRes.data && dbRes.data.list) || []
      dbList.forEach((r) => {
        rows.push({
          module: 'data',
          moduleLabel: '数据安全',
          kindLabel: '数据库安全检查',
          target: r.targetSummary || r.name || '—',
          summary: `基线不通过 ${r.baselineFail || 0} · CVE ${r.cveMatchCount || 0}`,
          checkTime: r.checkTime || r.createTime || '—',
          _ts: this.parseTime(r.checkTime || r.createTime),
          route: { path: '/datasec/task/detail', query: { id: r.id } },
        })
      })
      rows.sort((a, b) => (b._ts || 0) - (a._ts || 0))
      return rows.slice(0, 12)
    },
    goTask(row) {
      if (row.route) this.$router.push(row.route)
    },
    handleResize() {
      if (this.chartHost) this.chartHost.resize()
      if (this.chartModule) this.chartModule.resize()
    },
    chartTooltip(extra) {
      return Object.assign({
        backgroundColor: 'rgba(26, 26, 46, 0.95)',
        borderColor: 'rgba(0, 212, 170, 0.2)',
        textStyle: { color: '#e2e8f0' },
      }, extra || {})
    },
    renderCharts() {
      this.renderHostPie()
      this.renderModuleBar()
    },
    renderHostPie() {
      const el = document.getElementById('hostPie')
      if (!el) return
      if (this.chartHost) this.chartHost.dispose()
      this.chartHost = echarts.init(el)
      const s = this.hostCheckStat
      const data = [
        { name: '通过', value: s.pass, itemStyle: { color: CHART_THEME.pass } },
        { name: '不通过', value: s.fail, itemStyle: { color: CHART_THEME.fail } },
        { name: '异常', value: s.error, itemStyle: { color: CHART_THEME.warn } },
      ].filter((d) => d.value > 0)
      if (!data.length) {
        data.push({ name: '暂无数据', value: 1, itemStyle: { color: 'rgba(148,163,184,0.2)' } })
      }
      this.chartHost.setOption({
        tooltip: this.chartTooltip({ trigger: 'item', formatter: '{b}: {c} ({d}%)' }),
        legend: {
          bottom: 0,
          textStyle: { color: CHART_THEME.text, fontSize: 12 },
        },
        series: [{
          type: 'pie',
          radius: ['42%', '68%'],
          center: ['50%', '45%'],
          label: { color: CHART_THEME.text, fontSize: 12 },
          data,
        }],
      })
    },
    renderModuleBar() {
      const el = document.getElementById('moduleBar')
      if (!el) return
      if (this.chartModule) this.chartModule.dispose()
      this.chartModule = echarts.init(el)
      const names = this.moduleBarData.map((d) => d.name)
      const values = this.moduleBarData.map((d) => d.value)
      this.chartModule.setOption({
        tooltip: this.chartTooltip({ trigger: 'axis' }),
        grid: { left: 12, right: 12, top: 24, bottom: 8, containLabel: true },
        xAxis: {
          type: 'category',
          data: names,
          axisLabel: { color: CHART_THEME.axis, fontSize: 11, interval: 0 },
          axisLine: { lineStyle: { color: CHART_THEME.split } },
          axisTick: { show: false },
        },
        yAxis: {
          type: 'value',
          minInterval: 1,
          axisLabel: { color: CHART_THEME.axis },
          splitLine: { lineStyle: { color: CHART_THEME.split, type: 'dashed' } },
        },
        series: [{
          type: 'bar',
          data: values,
          barMaxWidth: 36,
          itemStyle: {
            color: CHART_THEME.primary,
            borderRadius: [4, 4, 0, 0],
          },
        }],
      })
    },
  },
}
</script>

<style lang="less" scoped>
@import '../pages/bas/css/bas-list-page.less';

.dashboard {
  height: 100%;
  overflow-y: auto;
  padding: 4px 0 24px;
  box-sizing: border-box;
}

.welcome-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 20px 24px;
  background: linear-gradient(135deg, rgba(0, 212, 170, 0.12), rgba(124, 58, 237, 0.08));
  border: 1px solid rgba(0, 212, 170, 0.15);
  border-radius: 12px;
}

.welcome-title {
  margin: 0 0 6px;
  font-size: 20px;
  font-weight: 600;
  color: #e2e8f0;
}

.welcome-desc {
  margin: 0;
  font-size: 13px;
  color: #94a3b8;
}

.module-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 16px;
}

.module-card {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px;
  background: #1a1a2e;
  border: 1px solid rgba(0, 212, 170, 0.08);
  border-radius: 12px;
  text-decoration: none;
  transition: all 0.25s;

  &:hover {
    border-color: rgba(0, 212, 170, 0.28);
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
  }
}

.mod-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  color: #fff;
  flex-shrink: 0;
}

.mod-host .mod-icon { background: linear-gradient(135deg, #00d4aa, #00b894); }
.mod-app .mod-icon { background: linear-gradient(135deg, #7c3aed, #6d28d9); }
.mod-data .mod-icon { background: linear-gradient(135deg, #3b82f6, #2563eb); }

.mod-body {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.mod-label {
  font-size: 13px;
  color: #94a3b8;
  margin-bottom: 4px;
}

.mod-value {
  font-size: 28px;
  font-weight: 700;
  color: #e2e8f0;
  line-height: 1.2;
}

.mod-sub {
  margin-top: 6px;
  font-size: 12px;
  color: #64748b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.rules-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.rule-chip {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(0, 212, 170, 0.08);
  border-radius: 10px;
}

.rule-name {
  font-size: 13px;
  color: #94a3b8;
}

.rule-count {
  font-size: 18px;
  font-weight: 600;
  color: #00d4aa;
}

.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
}

.chart-card {
  background: #1a1a2e;
  border: 1px solid rgba(0, 212, 170, 0.08);
  border-radius: 12px;
  overflow: hidden;

  &:hover {
    border-color: rgba(0, 212, 170, 0.15);
  }
}

.recent-card {
  margin-bottom: 16px;
}

.chart-hd {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(0, 212, 170, 0.06);

  h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 600;
    color: rgba(226, 232, 240, 0.85);
  }
}

.recent-hint {
  font-size: 12px;
  color: #64748b;
}

.chart-bd {
  padding: 12px 16px 16px;
}

.table-bd {
  padding: 0 16px 16px;
}

.chart-box {
  height: 220px;
}

.quick-row {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 12px;
}

.quick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(0, 212, 170, 0.08);
  border-radius: 10px;
  text-decoration: none;
  color: #94a3b8;
  font-size: 12px;
  transition: all 0.2s;

  i {
    font-size: 20px;
    color: #00d4aa;
  }

  &:hover {
    border-color: rgba(0, 212, 170, 0.25);
    color: #e2e8f0;
    background: rgba(0, 212, 170, 0.06);
  }
}

@media (max-width: 1200px) {
  .module-row,
  .rules-row {
    grid-template-columns: 1fr 1fr;
  }
  .charts-grid {
    grid-template-columns: 1fr;
  }
  .quick-row {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
