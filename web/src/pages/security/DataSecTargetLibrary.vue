<template>
  <div class="security-container">
    <div class="main-title">数据库目标库</div>
    <p class="page-intro">保存常用数据库连接，创建扫描任务时可直接选用，无需重复填写账号密码。</p>

    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="openEdit()">新增目标</el-button>
          <el-dropdown trigger="click" @command="onTestCommand">
            <el-button size="small" :loading="batchTestLoading">
              连接测试<i class="el-icon-arrow-down el-icon--right"></i>
            </el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="selected" :disabled="!selectedRows.length">
                测试所选（{{ selectedRows.length }}）
              </el-dropdown-item>
              <el-dropdown-item command="page" :disabled="!tableData.length">测试当前页全部</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <el-dropdown trigger="click" @command="onImportExportCommand">
            <el-button size="small">
              导入/导出<i class="el-icon-arrow-down el-icon--right"></i>
            </el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="import">导入 JSON</el-dropdown-item>
              <el-dropdown-item divided command="export-no-pwd">导出（不含密码）</el-dropdown-item>
              <el-dropdown-item command="export-pwd">导出（含密码）</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
        <div class="serach-condition">
          <el-select v-model="filterDbType" placeholder="库类型" size="small" clearable class="filter-type" @change="loadList">
            <el-option v-for="o in dbTypeOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-input v-model="search" placeholder="搜索名称/地址/分组" size="small" clearable class="search-input" @keydown.enter.native="loadList" />
          <el-button type="primary" size="small" @click="loadList">搜索</el-button>
        </div>
      </div>

      <el-table v-loading="loading" :data="tableData" class="myTable" @selection-change="onSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column prop="name" label="名称" min-width="140" :show-overflow-tooltip="true" />
        <el-table-column prop="groupName" label="分组" width="100" :show-overflow-tooltip="true" />
        <el-table-column label="类型" width="96">
          <template slot-scope="scope">{{ dbTypeName(scope.row.dbType) }}</template>
        </el-table-column>
        <el-table-column prop="dbHost" label="地址" width="140" :show-overflow-tooltip="true" />
        <el-table-column prop="dbPort" label="端口" width="72" />
        <el-table-column prop="dbName" label="库名" width="100" :show-overflow-tooltip="true" />
        <el-table-column prop="dbUser" label="用户" width="100" :show-overflow-tooltip="true" />
        <el-table-column label="密码" width="72">
          <template slot-scope="scope">{{ scope.row.hasPassword ? '已保存' : '-' }}</template>
        </el-table-column>
        <el-table-column prop="updateTime" label="更新时间" width="160" />
        <el-table-column label="操作" width="200" fixed="right">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="testRow(scope.row)">测试</el-link>
            <el-link :underline="false" class="link_primary" @click="openEdit(scope.row)">编辑</el-link>
            <el-link :underline="false" class="link_danger" @click="removeRow(scope.row)">删除</el-link>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        @current-change="onPageChange"
      />
    </div>

    <el-dialog :title="editForm.id ? '编辑目标' : '新增目标'" :visible.sync="dialogVisible" width="560px">
      <el-form ref="editFormRef" :model="editForm" :rules="rules" label-width="88px" size="small">
        <el-form-item label="名称" prop="name">
          <el-input v-model="editForm.name" placeholder="便于识别的名称" />
        </el-form-item>
        <el-form-item label="分组">
          <el-input v-model="editForm.groupName" placeholder="如：生产环境 / 测试库" />
        </el-form-item>
        <el-form-item label="数据库类型" prop="dbType">
          <el-select v-model="editForm.dbType" style="width:100%">
            <el-option v-for="o in dbTypeOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="地址" prop="dbHost">
          <el-input v-model="editForm.dbHost" />
        </el-form-item>
        <el-form-item label="端口">
          <el-input v-model.number="editForm.dbPort" type="number" />
        </el-form-item>
        <el-form-item label="库名">
          <el-input v-model="editForm.dbName" />
        </el-form-item>
        <el-form-item label="用户" prop="dbUser">
          <el-input v-model="editForm.dbUser" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="dbPassword">
          <el-input v-model="editForm.dbPassword" type="password" show-password :placeholder="editForm.id ? '留空则不修改' : '请输入密码'" autocomplete="new-password" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="editForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button :loading="testLoading" @click="testCurrentForm">连接测试</el-button>
        <el-button type="primary" :loading="saveLoading" @click="save">保存</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'
import { pickImportFile } from './utils/datasecTargetExport.js'

const DB_TYPES = [
  { value: 1, label: 'MySQL' },
  { value: 2, label: 'PostgreSQL' },
  { value: 3, label: 'MongoDB' },
  { value: 4, label: 'Redis' },
  { value: 5, label: 'CouchDB' }
]

export default {
  name: 'DataSecTargetLibrary',
  data() {
    return {
      loading: false,
      saveLoading: false,
      testLoading: false,
      batchTestLoading: false,
      selectedRows: [],
      tableData: [],
      total: 0,
      page: 1,
      pageSize: 20,
      search: '',
      filterDbType: null,
      dialogVisible: false,
      dbTypeOptions: DB_TYPES,
      editForm: this.emptyForm(),
      rules: {
        dbHost: [{ required: true, message: '请输入地址', trigger: 'blur' }],
        dbUser: [{ required: true, message: '请输入用户', trigger: 'blur' }],
        dbType: [{ required: true, message: '请选择类型', trigger: 'change' }]
      }
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/datasec/targets'
    this.loadList()
  },
  methods: {
    emptyForm() {
      return { id: 0, name: '', groupName: '', dbType: 1, dbHost: '', dbPort: 3306, dbName: '', dbUser: '', dbPassword: '', remark: '' }
    },
    dbTypeName(t) {
      const o = DB_TYPES.find((x) => x.value === t)
      return o ? o.label : '-'
    },
    async loadList() {
      this.loading = true
      try {
        const res = await security.getDatasecTargetList({
          page: this.page,
          size: this.pageSize,
          search: this.search || undefined,
          dbType: this.filterDbType || undefined
        })
        if (res.code === 200 && res.data) {
          this.tableData = res.data.list || []
          this.total = res.data.total || 0
        }
      } finally {
        this.loading = false
      }
    },
    onPageChange(p) {
      this.page = p
      this.loadList()
    },
    onSelectionChange(rows) {
      this.selectedRows = rows || []
    },
    onTestCommand(command) {
      if (command === 'selected') this.batchTestSelected()
      else if (command === 'page') this.batchTestCurrentPage()
    },
    onImportExportCommand(command) {
      if (command === 'import') this.importJson()
      else if (command === 'export-no-pwd') this.exportJson(false)
      else if (command === 'export-pwd') this.exportJson(true)
    },
    buildTestPayload() {
      return {
        id: this.editForm.id || 0,
        dbType: this.editForm.dbType,
        dbHost: (this.editForm.dbHost || '').trim(),
        dbPort: Number(this.editForm.dbPort) || 0,
        dbName: (this.editForm.dbName || '').trim(),
        dbUser: (this.editForm.dbUser || '').trim(),
        dbPassword: this.editForm.dbPassword || ''
      }
    },
    async testCurrentForm() {
      if (!(this.editForm.dbHost || '').trim()) {
        this.$message({ message: '请先填写地址', type: 'warning' })
        return
      }
      if (!(this.editForm.dbUser || '').trim()) {
        this.$message({ message: '请先填写用户名', type: 'warning' })
        return
      }
      if (!this.editForm.id && !this.editForm.dbPassword) {
        this.$message({ message: '新增目标请先填写密码', type: 'warning' })
        return
      }
      this.testLoading = true
      try {
        const res = await security.testDatasecTargetConn(this.buildTestPayload())
        if (res.code === 200 && res.data) {
          this.$message({
            message: res.data.message || (res.data.ok ? '连接成功' : '连接失败'),
            type: res.data.ok ? 'success' : 'error',
            duration: 5000
          })
        } else {
          this.$message({ message: res.msg || '连接测试失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '连接测试失败: ' + (e.message || ''), type: 'error' })
      } finally {
        this.testLoading = false
      }
    },
    async testRow(row) {
      this.batchTestLoading = true
      try {
        const res = await security.batchTestDatasecTargetConn({ ids: [row.id] })
        this.showBatchTestResult(res)
      } finally {
        this.batchTestLoading = false
      }
    },
    async batchTestSelected() {
      await this.runBatchTest(this.selectedRows.map((r) => r.id))
    },
    async batchTestCurrentPage() {
      await this.runBatchTest(this.tableData.map((r) => r.id))
    },
    async runBatchTest(ids) {
      if (!ids.length) {
        this.$message({ message: '没有可测试的目标', type: 'warning' })
        return
      }
      this.batchTestLoading = true
      try {
        const res = await security.batchTestDatasecTargetConn({ ids })
        this.showBatchTestResult(res)
      } catch (e) {
        this.$message({ message: '批量测试失败: ' + (e.message || ''), type: 'error' })
      } finally {
        this.batchTestLoading = false
      }
    },
    showBatchTestResult(res) {
      if (res.code !== 200 || !res.data) {
        this.$message({ message: res.msg || '测试失败', type: 'error' })
        return
      }
      const d = res.data
      const msg = d.fail === 0
        ? `全部 ${d.ok} 个目标连接成功`
        : `成功 ${d.ok} 个，失败 ${d.fail} 个`
      this.$message({ message: msg, type: d.fail === 0 ? 'success' : 'warning', duration: 5000 })
      if (d.fail > 0 && (d.results || []).length) {
        const fails = d.results.filter((r) => !r.ok).slice(0, 3)
        const detail = fails.map((r) => `${r.name || r.dbHost}: ${r.message}`).join('；')
        if (detail) {
          setTimeout(() => {
            this.$message({ message: detail, type: 'error', duration: 8000, showClose: true })
          }, 300)
        }
      }
    },
    openEdit(row) {
      if (row) {
        this.editForm = {
          id: row.id,
          name: row.name,
          groupName: row.groupName,
          dbType: row.dbType,
          dbHost: row.dbHost,
          dbPort: row.dbPort,
          dbName: row.dbName,
          dbUser: row.dbUser,
          dbPassword: '',
          remark: row.remark
        }
      } else {
        this.editForm = this.emptyForm()
      }
      this.dialogVisible = true
      this.$nextTick(() => this.$refs.editFormRef && this.$refs.editFormRef.clearValidate())
    },
    async save() {
      this.$refs.editFormRef.validate(async (valid) => {
        if (!valid) return
        if (!this.editForm.id && !this.editForm.dbPassword) {
          this.$message({ message: '请填写密码', type: 'warning' })
          return
        }
        this.saveLoading = true
        try {
          const res = await security.saveDatasecTarget({ ...this.editForm })
          if (res.code === 200) {
            this.$message({ message: '保存成功', type: 'success' })
            this.dialogVisible = false
            this.loadList()
          } else {
            this.$message({ message: res.msg || '保存失败', type: 'error' })
          }
        } finally {
          this.saveLoading = false
        }
      })
    },
    removeRow(row) {
      this.$confirm(`确定删除目标「${row.name || row.dbHost}」？`, '提示', { type: 'warning' }).then(async () => {
        const res = await security.deleteDatasecTarget({ id: row.id })
        if (res.code === 200) {
          this.$message({ message: '已删除', type: 'success' })
          this.loadList()
        } else {
          this.$message({ message: res.msg || '删除失败', type: 'error' })
        }
      }).catch(() => {})
    },
    async importJson() {
      try {
        const text = await pickImportFile('.json')
        const data = JSON.parse(text)
        const items = Array.isArray(data) ? data : (data.items || data.targets || [])
        const res = await security.importDatasecTargets({ items })
        if (res.code === 200) {
          this.$message({ message: `成功导入 ${(res.data && res.data.imported) || 0} 条`, type: 'success' })
          this.loadList()
        } else {
          this.$message({ message: res.msg || '导入失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '导入失败: ' + (e.message || ''), type: 'error' })
      }
    },
    async exportJson(includePassword) {
      const res = await security.exportDatasecTargets({ includePassword: includePassword ? '1' : '0' })
      if (res.code !== 200 || !res.data) {
        this.$message({ message: res.msg || '导出失败', type: 'error' })
        return
      }
      const blob = new Blob([JSON.stringify(res.data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `datasec-target-library-${Date.now()}.json`
      a.click()
      URL.revokeObjectURL(url)
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

.page-intro {
  color: #94a3b8;
  font-size: 13px;
  margin: 0 0 12px;
}

.filter-type { width: 120px; margin-right: 8px; }
.search-input { width: 220px; margin-right: 8px; }

.operationbutton {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.link_danger { color: #f87171; margin-left: 8px; }
</style>
