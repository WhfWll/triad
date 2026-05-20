<template>
  <div class="security-container">
    <div class="main-title" v-if="!embedded">专项应用检测</div>

    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="btnCreate">新建检测任务</el-button>
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
        <el-table-column prop="appType" label="应用类型">
          <template slot-scope="scope">
            <span :class="getAppTypeClass(scope.row.appType)">{{ getAppTypeName(scope.row.appType) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="扫描目标" min-width="220" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <span>{{ scope.row.targetSummary || scope.row.targetUrl }}</span>
            <el-tag v-if="scope.row.targetCount > 1" size="mini" type="info" class="target-count-tag">
              {{ scope.row.targetCount }}个
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="vulnCount" label="发现漏洞数" />
        <el-table-column prop="riskLevel" label="风险等级">
          <template slot-scope="scope">
            <span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template slot-scope="scope">
            <span :class="getStatusClass(scope.row.status)">{{ getStatusName(scope.row.status, 'app') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" />
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
import {
  getRiskName,
  getRiskClass,
  getStatusName,
  getStatusClass,
  getAppTypeName
} from './appsecTaskLabels.js'

const APP_TYPE_CLASSES = {
  1: 'app-wanhui',
  2: 'app-yongyou',
  3: 'app-lanling',
  4: 'app-yunshikong',
  5: 'app-yisaitong',
  6: 'app-dlink',
  7: 'app-tongda',
  8: 'app-wordpress',
  9: 'app-thinkphp',
  10: 'app-springboot',
  11: 'app-generic'
}

export default {
  name: 'AppSpecificScan',
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
    getAppTypeName,
    getAppTypeClass(type) {
      return APP_TYPE_CLASSES[type] || 'app-default'
    },
    async getData() {
      const res = await security.getAppSpecificList({
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
      this.$router.push({ path: '/appsec/task/new', query: { type: 'app' } })
    },
    handleDetail(row) {
      this.$router.push({
        path: '/appsec/task/detail',
        query: { id: row.id, type: 'app' }
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

.app-wanhui, .app-yongyou, .app-lanling, .app-yunshikong, .app-yisaitong, .app-dlink, .app-tongda, .app-wordpress, .app-thinkphp, .app-springboot, .app-generic, .app-default {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.app-wanhui { background: rgba(245, 101, 101, 0.2); color: #f56565; }
.app-yongyou { background: rgba(97, 175, 254, 0.2); color: #61affe; }
.app-lanling { background: rgba(67, 153, 52, 0.2); color: #439934; }
.app-yunshikong { background: rgba(231, 162, 13, 0.2); color: #e7a20d; }
.app-yisaitong { background: rgba(125, 86, 205, 0.2); color: #7d56cd; }
.app-dlink { background: rgba(59, 130, 246, 0.2); color: #3b82f6; }
.app-tongda { background: rgba(249, 115, 22, 0.2); color: #f97316; }
.app-wordpress { background: rgba(212, 180, 68, 0.2); color: #d4b444; }
.app-thinkphp { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.app-springboot { background: rgba(106, 176, 76, 0.2); color: #6ab04c; }
.app-generic { background: rgba(148, 163, 184, 0.2); color: #94a3b8; }

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
