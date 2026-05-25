<template>
  <div class="security-container report-center">
    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" :disabled="!selectedRows.length" @click="batchDelete">批量删除</el-button>
        </div>
        <div class="serach-condition">
          <div class="page-hint">各安全检查模块生成的报告统一归集于此，支持预览和下载。</div>
        </div>
      </div>

      <el-table
        ref="table"
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        class="myTable"
        @selection-change="onSelectionChange"
      >
        <el-table-column type="selection" width="48" />
        <el-table-column prop="title" label="报告标题" min-width="200" :show-overflow-tooltip="true" />
        <el-table-column prop="moduleName" label="所属模块" width="130" />
        <el-table-column prop="taskName" label="关联任务" min-width="140" :show-overflow-tooltip="true" />
        <el-table-column prop="createTime" label="生成时间" width="168" />
        <el-table-column label="操作" width="180" align="right">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="previewReport(scope.row)">预览</el-link>
            <el-link :underline="false" class="link_primary" style="margin-left:12px" @click="downloadReport(scope.row)">下载</el-link>
            <el-popover
              placement="top"
              width="160"
              :ref="`del-${scope.row.id}`"
              popper-class="delButton_popper"
            >
              <p><i class="el-icon-warning"></i> 确定删除吗？</p>
              <div style="text-align:right;margin-top:8px">
                <el-button size="mini" @click="closeDelPopover(scope.row.id)">取消</el-button>
                <el-button size="mini" type="primary" @click="confirmDelete(scope.row)">确定</el-button>
              </div>
              <el-link slot="reference" :underline="false" class="link_danger" style="margin-left:12px">删除</el-link>
            </el-popover>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        background
        layout="total, prev, pager, next, sizes, jumper"
        :total="total"
        :page-size="pageSize"
        :current-page="currentPage"
        :page-sizes="[10, 20, 50, 100]"
        @current-change="onPageChange"
        @size-change="onSizeChange"
      />
    </div>

    <el-dialog
      title="报告预览"
      :visible.sync="previewVisible"
      width="960px"
      custom-class="theme-dialog"
      top="5vh"
      @closed="onPreviewClosed"
    >
      <div v-if="previewLoading" v-loading="true" class="preview-loading" />
      <iframe v-else-if="previewHtml" :srcdoc="previewHtml" class="preview-iframe" frameborder="0" />
      <p v-else class="empty-hint">暂无报告内容</p>
      <span slot="footer">
        <el-button type="primary" v-if="previewRow" @click="downloadReport(previewRow)">下载 HTML</el-button>
        <el-button @click="previewVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'

export default {
  name: 'ReportCenter',
  data() {
    return {
      loading: false,
      tableData: [],
      total: 0,
      currentPage: 1,
      pageSize: 20,
      selectedRows: [],
      previewVisible: false,
      previewLoading: false,
      previewHtml: '',
      previewRow: null
    }
  },
  created() {
    this.loadList()
  },
  methods: {
    async loadList() {
      this.loading = true
      try {
        const res = await security.getSecurityReportList({ page: this.currentPage, size: this.pageSize })
        if (res.code === 200 && res.data) {
          this.tableData = res.data.list || []
          this.total = res.data.total || 0
        }
      } finally {
        this.loading = false
      }
    },
    onSelectionChange(rows) {
      this.selectedRows = rows
    },
    onPageChange(page) {
      this.currentPage = page
      this.loadList()
    },
    onSizeChange(size) {
      this.pageSize = size
      this.currentPage = 1
      this.loadList()
    },
    async previewReport(row) {
      this.previewRow = row
      this.previewVisible = true
      this.previewLoading = true
      this.previewHtml = ''
      try {
        const res = await security.getSecurityReportDetail({ id: row.id })
        if (res.code === 200 && res.data) {
          this.previewHtml = res.data.content
        } else {
          this.$message.warning('加载报告内容失败')
        }
      } finally {
        this.previewLoading = false
      }
    },
    downloadReport(row) {
      this.previewRow = row
      this.downloadHtml(row)
    },
    async downloadHtml(row) {
      let html = this.previewHtml
      if (!html && this.previewVisible && this.previewRow && this.previewRow.id === row.id) {
        html = this.previewHtml
      }
      if (!html) {
        try {
          const res = await security.getSecurityReportDetail({ id: row.id })
          if (res.code === 200 && res.data) {
            html = res.data.content
          }
        } catch { /* ignore */ }
      }
      if (!html) {
        this.$message.warning('报告内容为空')
        return
      }
      const blob = new Blob([html], { type: 'text/html;charset=utf-8' })
      const a = document.createElement('a')
      a.href = URL.createObjectURL(blob)
      a.download = `${row.title || 'report'}.html`
      a.click()
      URL.revokeObjectURL(a.href)
      this.$message.success('报告已下载')
    },
    closeDelPopover(id) {
      const ref = this.$refs[`del-${id}`]
      if (ref && ref.length > 0) ref[0].doClose()
      else if (ref) ref.doClose()
    },
    async confirmDelete(row) {
      try {
        const res = await security.deleteSecurityReport({ id: row.id })
        if (res.code === 200) {
          this.$message.success('报告已删除')
          this.loadList()
        }
      } finally {
        this.closeDelPopover(row.id)
      }
    },
    async batchDelete() {
      if (!this.selectedRows.length) return
      try {
        for (const row of this.selectedRows) {
          await security.deleteSecurityReport({ id: row.id })
        }
        this.$message.success(`已删除 ${this.selectedRows.length} 份报告`)
        this.selectedRows = []
        this.loadList()
      } catch { /* */ }
    },
    onPreviewClosed() {
      this.previewHtml = ''
      this.previewRow = null
    }
  }
}
</script>

<style scoped lang="less">
@import '../bas/css/bas-list-page.less';

.report-center {
  height: 100%;

  .page-hint {
    font-size: 13px;
    color: #64748b;
    line-height: 32px;
  }

  .preview-loading {
    height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .preview-iframe {
    width: 100%;
    height: 70vh;
    border: 1px solid #e2e8f0;
    border-radius: 4px;
  }
}
</style>
