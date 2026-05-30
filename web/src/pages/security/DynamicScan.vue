<template>
  <div class="security-container">
    <div class="main-title" v-if="!embedded">动态扫描</div>

    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="btnCreate">新建扫描任务</el-button>
        </div>
        <div class="serach-condition">
          <div class="search-text">
            <el-input
              placeholder="搜索任务名称"
              @keydown.enter.native="handlesearch"
              v-model="formData.search"
              class="input-with-select"
              size="small"
              clearable
            />
            <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
          </div>
          <div>
            <el-button type="primary" size="small" @click="handleReset">重置</el-button>
          </div>
        </div>
      </div>

      <el-table :data="tableData" style="width: 100%" class="myTable" @selection-change="handleSelectionChange">
        <el-table-column width="55" type="selection" />
        <el-table-column prop="name" label="任务名称" :show-overflow-tooltip="true" />
        <el-table-column label="扫描目标" width="220" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <span>{{ scope.row.targetSummary || scope.row.targetUrl }}</span>
            <el-tag v-if="scope.row.targetCount > 1" size="mini" type="info" class="target-count-tag">
              {{ scope.row.targetCount }}个
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pageCount" label="爬取页面数" width="100" />
        <el-table-column prop="vulnCount" label="发现漏洞数" width="100" />
        <el-table-column prop="riskLevel" label="风险等级">
          <template slot-scope="scope">
            <span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template slot-scope="scope">
            <span :class="getStatusClass(scope.row.status)">{{ getStatusName(scope.row.status, 'dyn') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="scanTime" label="扫描时间" />
        <el-table-column label="操作" width="100">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="handleDetail(scope.row)">详情</el-link>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :page-size="pageSize"
        background
        layout="total, prev, pager, next, sizes, jumper"
        :total="totalpage"
        :current-page="currentpage"
        @current-change="handlecurrentchange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script>
import security from '@/api/security.js'
import { getRiskName, getRiskClass, getStatusName, getStatusClass } from './appsecTaskLabels.js'

export default {
  name: 'DynamicScan',
  props: {
    embedded: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      multipleSelection: [],
      tableData: [],
      formData: {
        search: '',
        page: 1
      },
      pageSize: 10,
      currentpage: 1,
      totalpage: 0,
      pollTimer: null
    }
  },
  mounted() {
    this.getData()
  },
  activated() {
    if (this.embedded) {
      this.getData()
      this.startPollIfNeeded()
    }
  },
  deactivated() {
    this.stopPoll()
  },
  beforeDestroy() {
    this.stopPoll()
  },
  methods: {
    getRiskName,
    getRiskClass,
    getStatusName,
    getStatusClass,
    async getData() {
      const res = await security.getDynamicScanList({
        page: this.formData.page,
        size: this.pageSize,
        search: this.formData.search
      })
      if (res.code == 200) {
        this.tableData = (res.data && res.data.list) || []
        this.totalpage = (res.data && res.data.total) || 0
        this.startPollIfNeeded()
      } else {
        this.$message({ message: res.msg, type: 'error' })
      }
    },
    startPollIfNeeded() {
      this.stopPoll()
      if (!this.embedded) return
      if (!this.tableData.some(r => r.status === 2)) return
      this.pollTimer = setInterval(() => {
        this.getData()
        if (!this.tableData.some(r => r.status === 2)) this.stopPoll()
      }, 3000)
    },
    stopPoll() {
      if (this.pollTimer) {
        clearInterval(this.pollTimer)
        this.pollTimer = null
      }
    },
    btnCreate() {
      this.$router.push({ path: '/appsec/task/new', query: { type: 'dyn' } })
    },
    handleDetail(row) {
      this.$router.push({
        path: '/appsec/task/detail',
        query: { id: row.id, type: 'dyn' }
      })
    },
    handlesearch() {
      this.formData.page = 1
      this.getData()
      this.currentpage = 1
    },
    handleReset() {
      this.formData.page = 1
      this.formData.search = ''
      this.getData()
      this.currentpage = 1
    },
    handlecurrentchange(t) {
      this.formData.page = t
      this.getData()
      this.currentpage = t
    },
    handleSizeChange(t) {
      this.formData.page = 1
      this.pageSize = t
      this.getData()
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

.risk-critical, .risk-high, .risk-medium, .risk-low, .risk-info {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.risk-critical { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
.risk-high { background: rgba(249, 115, 22, 0.2); color: #f97316; }
.risk-medium { background: rgba(234, 179, 8, 0.2); color: #eab308; }
.risk-low { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.risk-info { background: rgba(148, 163, 184, 0.2); color: #94a3b8; }

.status-wait, .status-running, .status-complete {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-wait { background: rgba(234, 179, 8, 0.2); color: #eab308; }
.status-running { background: rgba(59, 130, 246, 0.2); color: #3b82f6; }
.status-complete { background: rgba(16, 185, 129, 0.2); color: #10b981; }
</style>
