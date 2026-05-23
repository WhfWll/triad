<template>
  <div class="security-container">
    <p class="page-intro">
      CVE 漏洞库管理。当前使用 <strong>default-cve.db</strong>（SQLite），包含 <strong>{{ dbInfo.totalRecords || '—' }}</strong> 条 CVE 记录。
      <router-link class="link-to-tasks" to="/hostsec/tasks">前往任务管理</router-link>
    </p>

    <div class="list_box">
      <div v-if="loading" class="detail-loading">加载中…</div>
      <div v-else>
        <div class="rules-header">
          <p class="rules-summary">
            CVE 漏洞库状态：
            <el-tag :type="dbInfo.isAvailable ? 'success' : 'danger'" size="small">
              {{ dbInfo.isAvailable ? '可用' : '不可用' }}
            </el-tag>
            <span v-if="dbInfo.isAvailable" style="margin-left: 12px;">
              总计 <strong>{{ dbInfo.totalRecords }}</strong> 条漏洞记录
            </span>
          </p>
          <div class="rules-header-actions">
            <el-button size="small" icon="el-icon-refresh" @click="loadData">刷新</el-button>
          </div>
        </div>

        <div class="rules-toolbar">
          <el-input
            v-model="keyword"
            clearable
            size="small"
            placeholder="搜索 CVE ID / 产品名 / 标题"
            style="width: 320px"
            @keyup.enter.native="doSearch"
          />
          <el-button size="small" type="primary" icon="el-icon-search" style="margin-left: 8px" @click="doSearch">搜索</el-button>
        </div>

        <el-table :data="tableData" class="myTable" max-height="520" size="small" v-loading="searchLoading">
          <el-table-column prop="cveId" label="CVE ID" width="140">
            <template slot-scope="scope">
              <el-link :underline="false" class="link_primary" @click="openDetail(scope.row)">{{ scope.row.cveId }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="severity" label="严重程度" width="100">
            <template slot-scope="scope">
              <el-tag :type="severityTagType(scope.row.severity)" size="mini">{{ scope.row.severity || '—' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="product" label="影响产品" width="160" :show-overflow-tooltip="true" />
          <el-table-column prop="vendor" label="厂商" width="120" :show-overflow-tooltip="true" />
          <el-table-column prop="title" label="标题" min-width="200" :show-overflow-tooltip="true" />
          <el-table-column prop="description" label="描述" min-width="240" :show-overflow-tooltip="true" />
        </el-table>

        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          :page-size="pageSize"
          :page-sizes="[20, 50, 100]"
          :current-page.sync="currentPage"
          style="margin-top: 12px"
          @size-change="onSizeChange"
          @current-change="onPageChange"
        />
      </div>
    </div>

    <el-dialog title="CVE 详情" :visible.sync="detailVisible" width="700px" custom-class="theme-dialog">
      <div v-if="detailRow" class="rule-detail">
        <p><span class="k">CVE ID</span><el-tag :type="severityTagType(detailRow.severity)" size="small">{{ detailRow.cveId }}</el-tag></p>
        <p><span class="k">严重程度</span>{{ detailRow.severity || '—' }}</p>
        <p><span class="k">影响产品</span>{{ detailRow.product || '—' }}</p>
        <p><span class="k">厂商</span>{{ detailRow.vendor || '—' }}</p>
        <p><span class="k">标题</span>{{ detailRow.title || '—' }}</p>
        <p><span class="k">描述</span>{{ detailRow.description || '—' }}</p>
      </div>
      <span slot="footer">
        <el-button type="primary" @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'

export default {
  name: 'VulnDetectionRules',
  data() {
    return {
      loading: false,
      searchLoading: false,
      dbInfo: {},
      keyword: '',
      tableData: [],
      total: 0,
      currentPage: 1,
      pageSize: 20,
      detailVisible: false,
      detailRow: null
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      this.loading = true
      try {
        const infoRes = await security.getCveDBInfo()
        if (infoRes.code === 200) {
          this.dbInfo = infoRes.data || {}
        }
        await this.doSearch()
      } finally {
        this.loading = false
      }
    },
    async doSearch() {
      this.searchLoading = true
      try {
        const res = await security.queryCveDB({
          keyword: this.keyword,
          page: this.currentPage,
          size: this.pageSize
        })
        if (res.code === 200) {
          this.tableData = res.data.list || []
          this.total = res.data.total || 0
        }
      } finally {
        this.searchLoading = false
      }
    },
    onSizeChange(val) {
      this.pageSize = val
      this.currentPage = 1
      this.doSearch()
    },
    onPageChange(val) {
      this.currentPage = val
      this.doSearch()
    },
    openDetail(row) {
      this.detailRow = row
      this.detailVisible = true
    },
    severityTagType(severity) {
      if (!severity) return 'info'
      const s = severity.toUpperCase()
      if (s === 'CRITICAL' || s === 'HIGH') return 'danger'
      if (s === 'MEDIUM') return 'warning'
      if (s === 'LOW') return 'info'
      return 'info'
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

.page-intro {
  color: #94a3b8;
  font-size: 13px;
  line-height: 1.55;
  margin: 0 0 16px;
  max-width: 960px;
}
.rules-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.rules-summary {
  font-size: 14px;
  color: #e2e8f0;
  margin: 0;
}
.rules-header-actions {
  display: flex;
  gap: 8px;
}
.rules-toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}
.rule-detail p {
  margin: 8px 0;
  font-size: 14px;
  color: #e2e8f0;
}
.rule-detail .k {
  display: inline-block;
  width: 100px;
  color: #94a3b8;
  margin-right: 8px;
}
</style>