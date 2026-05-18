<template>
  <div class="security-container">
    <div class="main-title">
      <el-link :underline="false" class="link-back" @click="$router.push('/hostsec/tasks')">
        <i class="el-icon-arrow-left"></i> 返回
      </el-link>
      任务详情 · {{ kindLabel }}
      <span class="title-sub">#{{ taskId }}</span>
    </div>

    <el-tabs v-model="activeTab" class="detail-tabs">
      <el-tab-pane label="概况" name="overview">
        <div v-loading="statLoading" class="stat-cards">
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
        <el-table v-loading="targetsLoading" :data="targets" style="width: 100%" class="myTable">
          <el-table-column type="index" label="序号" width="60" />
          <el-table-column prop="targetIp" label="目标主机" width="160" />
          <el-table-column prop="osTypeName" label="操作系统" width="120" />
          <el-table-column prop="totalRules" label="检查项数" width="100" />
          <el-table-column prop="passCount" label="通过" width="72" />
          <el-table-column prop="failCount" label="不通过" width="82" />
          <el-table-column prop="errorCount" label="异常" width="72" />
          <el-table-column label="操作" width="90">
            <template slot-scope="scope">
              <el-link :underline="false" class="link_primary" @click="filterByTarget(scope.row)">查看结果</el-link>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="核查项" name="items">
        <div class="items-toolbar">
          <el-select v-model="itemFilter.targetIp" placeholder="按目标筛选" size="small" clearable @change="onItemFilterChange" style="width: 160px">
            <el-option v-for="t in targets" :key="t.targetIp" :value="t.targetIp" :label="t.targetIp" />
          </el-select>
          <el-select v-model="itemFilter.result" placeholder="按结果筛选" size="small" clearable @change="onItemFilterChange" style="width: 120px; margin-left: 8px">
            <el-option value="pass" label="通过" />
            <el-option value="fail" label="不通过" />
          </el-select>
          <span class="items-count">共 {{ filteredItems.length }} 项</span>
        </div>
        <el-table v-loading="itemsLoading" :data="filteredItems" style="width: 100%" class="myTable" max-height="560">
          <el-table-column prop="targetIp" label="目标主机" width="130" />
          <el-table-column prop="categoryName" label="分类" width="120" />
          <el-table-column prop="ruleName" label="检查项" min-width="160" :show-overflow-tooltip="true" />
          <el-table-column prop="resultName" label="结果" width="80" />
          <el-table-column prop="riskName" label="风险" width="80" />
          <el-table-column prop="expectedValue" label="期望值" min-width="140" :show-overflow-tooltip="true" />
          <el-table-column prop="actualValue" label="实际值" min-width="140" :show-overflow-tooltip="true" />
          <el-table-column prop="fixSuggestion" label="修复建议" min-width="160" :show-overflow-tooltip="true" />
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="日志" name="logs">
        <div v-if="logs.length === 0" class="logs-empty">暂无日志信息</div>
        <div v-else class="logs-list">
          <div v-for="(log, idx) in logs" :key="idx" class="log-item">
            <span class="log-time">{{ log.time }}</span>
            <span :class="['log-level', 'level-' + log.level]">{{ log.levelLabel }}</span>
            <span class="log-msg">{{ log.message }}</span>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import security from '@/api/security.js'

export default {
  name: 'HostTaskDetail',
  data() {
    return {
      taskId: 0,
      kindLabel: '安全配置核查',
      checkTime: '',
      activeTab: 'overview',
      statLoading: false,
      statData: {
        totalRules: 0,
        passCount: 0,
        failCount: 0,
        passRate: 0
      },
      targetsLoading: false,
      targets: [],
      itemsLoading: false,
      allItems: [],
      itemFilter: {
        targetIp: '',
        result: ''
      },
      logs: []
    }
  },
  computed: {
    filteredItems() {
      let list = this.allItems
      if (this.itemFilter.targetIp) {
        list = list.filter(i => i.targetIp === this.itemFilter.targetIp)
      }
      if (this.itemFilter.result === 'pass') {
        list = list.filter(i => i.checkResult === 1)
      } else if (this.itemFilter.result === 'fail') {
        list = list.filter(i => i.checkResult === 2)
      }
      return list
    }
  },
  created() {
    this.taskId = parseInt(this.$route.query.taskId) || 0
    this.kindLabel = this.$route.query.kindLabel || '安全配置核查'
    this.checkTime = this.$route.query.checkTime || ''
    if (this.taskId > 0) {
      this.loadStat()
      this.loadTargets()
      this.loadItems()
    }
  },
  methods: {
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
    filterByTarget(row) {
      this.itemFilter.targetIp = row.targetIp
      this.itemFilter.result = ''
      this.activeTab = 'items'
    },
    onItemFilterChange() {
    }
  }
}
</script>

<style lang="less" scoped>
.main-title {
  font-size: 18px;
  font-weight: 600;
  color: #e2e8f0;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.link-back {
  color: #94a3b8;
  font-size: 14px;
  &:hover { color: #00d4aa; }
}
.title-sub {
  font-size: 14px;
  color: #64748b;
  font-weight: 400;
}

.detail-tabs {
  ::v-deep .el-tabs__item {
    color: #94a3b8;
    font-size: 14px;
    &.is-active { color: #00d4aa; }
  }
  ::v-deep .el-tabs__active-bar { background-color: #00d4aa; }
  ::v-deep .el-tabs__nav-wrap::after { background-color: rgba(255,255,255,0.06); }
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

.targets-toolbar, .items-toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}
.targets-count, .items-count {
  font-size: 12px;
  color: #94a3b8;
  margin-left: 12px;
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
</style>
