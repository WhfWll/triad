<template>
  <div class="appsec-task-logs">
    <div class="logs-toolbar">
      <el-input
        v-model="search"
        placeholder="搜索目标 URL"
        size="small"
        clearable
        class="search-input"
        @keydown.enter.native="handleSearch"
      />
      <el-button type="primary" size="small" @click="handleSearch">搜索</el-button>
      <el-button size="small" @click="handleReset">重置</el-button>
    </div>

    <el-table v-loading="loading" :data="logTargets" class="myTable" max-height="420">
      <el-table-column prop="targetUrl" label="测试目标" min-width="220" :show-overflow-tooltip="true" />
      <el-table-column prop="statusName" label="状态" width="100" />
      <el-table-column prop="isAliveName" label="存活" width="88" />
      <el-table-column prop="createTime" label="创建时间" width="160" />
      <el-table-column prop="startTime" label="开始时间" width="160" />
      <el-table-column prop="endTime" label="结束时间" width="160" />
      <el-table-column label="操作" width="88" fixed="right">
        <template slot-scope="scope">
          <el-button type="text" size="small" @click="openLogDetail(scope.row)">日志</el-button>
        </template>
      </el-table-column>
    </el-table>
    <p v-if="!loading && !logTargets.length" class="empty-hint">暂无执行日志</p>

    <el-pagination
      v-if="total > 0"
      background
      layout="total, prev, pager, next"
      :total="total"
      :page-size="pageSize"
      :current-page="page"
      class="logs-pager"
      @current-change="onPageChange"
    />

    <el-dialog
      :title="detailTitle"
      :visible.sync="detailVisible"
      width="900px"
      custom-class="theme-dialog"
      @closed="onDetailClosed"
    >
      <div v-loading="detailLoading" class="log-detail-body">
        <ul v-if="logLines.length" class="log-line-list">
          <li v-for="(item, index) in logLines" :key="index">
            <span class="log-time">[{{ item.createTime }}]</span>
            <span v-if="item.pocname" class="log-poc">{{ item.pocname }}:</span>
            <span class="log-text">{{ item.result }}</span>
          </li>
        </ul>
        <p v-else class="empty-hint">该目标暂无扫描过程日志</p>
      </div>
      <span slot="footer">
        <el-button @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { task } from '@/api/task.js'

export default {
  name: 'AppSecTaskLogs',
  props: {
    taskId: { type: [String, Number], required: true }
  },
  data() {
    return {
      loading: false,
      search: '',
      page: 1,
      pageSize: 10,
      total: 0,
      logTargets: [],
      detailVisible: false,
      detailLoading: false,
      detailTitle: '',
      logLines: [],
      pollTimer: null,
      currentLogId: null
    }
  },
  watch: {
    taskId: {
      immediate: true,
      handler() {
        this.page = 1
        this.fetchList()
      }
    }
  },
  beforeDestroy() {
    this.stopPoll()
  },
  methods: {
    async fetchList() {
      if (!this.taskId) return
      this.loading = true
      try {
        const res = await task.getLoglist({
          taskId: Number(this.taskId),
          search: this.search || undefined,
          page: this.page,
          size: this.pageSize
        })
        if (res.code == 200) {
          this.logTargets = (res.data && res.data.list) || []
          this.total = (res.data && res.data.total) || 0
        }
      } catch {
        this.$message.error('加载执行日志失败')
      } finally {
        this.loading = false
      }
    },
    handleSearch() {
      this.page = 1
      this.fetchList()
    },
    handleReset() {
      this.search = ''
      this.page = 1
      this.fetchList()
    },
    onPageChange(p) {
      this.page = p
      this.fetchList()
    },
    openLogDetail(row) {
      this.detailTitle = row.targetUrl || '扫描日志'
      this.currentLogId = row.id
      this.detailVisible = true
      this.fetchLogInfo()
      this.stopPoll()
      if (row.status === 3) {
        this.pollTimer = setInterval(() => this.fetchLogInfo(), 5000)
      }
    },
    async fetchLogInfo() {
      if (!this.currentLogId) return
      this.detailLoading = true
      try {
        const res = await task.loginfo({ taskLogId: this.currentLogId })
        if (res.code == 200) {
          this.logLines = (res.data && res.data.list) || []
        }
      } catch {
        /* ignore */
      } finally {
        this.detailLoading = false
      }
    },
    onDetailClosed() {
      this.stopPoll()
      this.currentLogId = null
      this.logLines = []
    },
    stopPoll() {
      if (this.pollTimer) {
        clearInterval(this.pollTimer)
        this.pollTimer = null
      }
    }
  }
}
</script>

<style lang="less" scoped>
.appsec-task-logs {
  .logs-toolbar {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;
    flex-wrap: wrap;

    .search-input {
      width: 280px;
    }
  }

  .logs-pager {
    margin-top: 12px;
    text-align: right;
  }

  .empty-hint {
    color: #64748b;
    font-size: 13px;
    text-align: center;
    padding: 24px 0;
  }

  .log-line-list {
    list-style: none;
    margin: 0;
    padding: 0;
    max-height: 480px;
    overflow: auto;

    li {
      font-size: 13px;
      line-height: 1.6;
      color: #94a3b8;
      padding: 4px 0;
      border-bottom: 1px solid rgba(255, 255, 255, 0.04);
      word-break: break-all;
    }

    .log-time {
      color: #64748b;
      margin-right: 6px;
    }

    .log-poc {
      color: #00d4aa;
      margin-right: 4px;
    }
  }
}
</style>
