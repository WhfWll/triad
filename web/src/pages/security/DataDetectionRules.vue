<template>
  <div class="security-container">
    <p class="page-intro">
      规则存储于 <strong>datasec_rule</strong> 表。库中有启用规则时，数据库安全检查任务将优先使用库规则；否则使用内置规则（约 {{ builtinHint }} 条）。
      需求指标 <strong>3800+</strong> 条：可从本地 <strong>data/default-cve.db</strong> 导入数据库相关 CVE（约 {{ cvePreview.availableInDb || '—' }} 条可用）。
      CVE 规则为<strong>知识库条目</strong>；任务扫描时会<strong>自动探测版本并在线匹配 CVE</strong>（最多展示 50 条命中）。
      <router-link class="link-to-tasks" to="/datasec/tasks">前往任务管理</router-link>
    </p>

    <div class="list_box">
      <div v-if="rulesLoading" class="detail-loading">加载中…</div>
      <div v-else-if="rulesPayload">
        <div class="rules-header">
          <div class="rules-summary-block">
            <p class="rules-summary">
              库中规则 <strong>{{ enabledCount }}</strong> 条（共 {{ rulesPayload.total }} 条记录） ·
              内置参考 <strong>{{ rulesPayload.builtinTotal || 0 }}</strong> 条 ·
              规划目标 <strong>{{ rulesPayload.targetTotal || 3800 }}</strong> 条
            </p>
            <el-progress
              :percentage="progressPercent"
              :stroke-width="10"
              :show-text="true"
              :format="() => `${enabledCount} / ${rulesPayload.targetTotal || 3800}`"
              class="target-progress"
            />
          </div>
          <div class="rules-header-actions">
            <el-button size="small" icon="el-icon-plus" type="success" @click="openCreate">新增规则</el-button>
            <el-button size="small" @click="importBuiltin">同步内置规则</el-button>
            <el-button
              size="small"
              type="warning"
              :loading="cveImporting"
              :disabled="cvePreviewLoading || cveAvailable <= 0"
              @click="importFromCve"
            >从 CVE 库导入</el-button>
            <el-button type="primary" size="small" icon="el-icon-upload2" @click="importDialogVisible = true">导入 JSON</el-button>
          </div>
        </div>

        <el-row :gutter="16" class="rules-summary-row">
          <el-col :span="12">
            <div class="rules-panel-title">按数据库类型</div>
            <el-table :data="rulesPayload.byDbType || []" size="small" class="myTable" max-height="220">
              <el-table-column prop="dbTypeName" label="类型" />
              <el-table-column prop="count" label="规则数" width="90" />
            </el-table>
          </el-col>
          <el-col :span="12">
            <div class="rules-panel-title">按检测分类</div>
            <el-table :data="rulesPayload.byCategory || []" size="small" class="myTable" max-height="220">
              <el-table-column prop="categoryName" label="分类" :show-overflow-tooltip="true" />
              <el-table-column prop="count" label="规则数" width="90" />
            </el-table>
          </el-col>
        </el-row>

        <div class="rules-toolbar">
          <el-select v-model="ruleFilterDb" clearable placeholder="全部数据库" size="small" style="width: 160px">
            <el-option v-for="o in rulesDbFilterOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-select v-model="ruleFilterCat" clearable placeholder="全部分类" size="small" style="width: 200px; margin-left: 8px">
            <el-option v-for="c in rulesCatFilterOptions" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
          <el-input
            v-model="ruleKeyword"
            clearable
            size="small"
            placeholder="搜索名称 / 描述 / 修复建议"
            style="width: 260px; margin-left: 8px"
          />
        </div>

        <el-table :data="pagedRules" class="myTable" max-height="480" size="small">
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="dbTypeName" label="数据库" width="110" />
          <el-table-column prop="categoryName" label="分类" width="120" :show-overflow-tooltip="true" />
          <el-table-column prop="riskName" label="风险" width="80" />
          <el-table-column prop="name" label="检查项" min-width="140" :show-overflow-tooltip="true" />
          <el-table-column prop="enabled" label="状态" width="72">
            <template slot-scope="scope">
              <span :class="scope.row.enabled === 1 ? 'tag-on' : 'tag-off'">{{ scope.row.enabled === 1 ? '启用' : '停用' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template slot-scope="scope">
              <el-link :underline="false" class="link_primary" @click="openDetail(scope.row)">详情</el-link>
              <el-link :underline="false" class="link_primary" style="margin-left: 10px" @click="openEdit(scope.row)">编辑</el-link>
              <el-link :underline="false" class="link_danger" style="margin-left: 10px" @click="confirmDelete(scope.row)">删除</el-link>
            </template>
          </el-table-column>
        </el-table>
        <div class="rules-pagination">
          <el-pagination
            background
            layout="total, sizes, prev, pager, next"
            :total="filteredRules.length"
            :page-size="pageSize"
            :page-sizes="[20, 50, 100, 200]"
            :current-page.sync="currentPage"
            @size-change="onPageSizeChange"
            @current-change="onPageChange"
          />
        </div>
        <p class="rules-footnote">
          导入 JSON 字段：<code>name</code>、<code>category</code>、<code>risk</code>、<code>dbType</code>（0=全部，1~5=MySQL/PostgreSQL/MongoDB/Redis/CouchDB）、
          <code>queries</code>（SQL 数组）、<code>expectedValue</code>、<code>matchType</code>（contains / exact / empty / always）。
        </p>
      </div>
      <div v-else class="detail-empty">
        <p>暂无规则数据。请先执行 SQL 创建 <code>datasec_rule</code> 表，或点击「同步内置规则」/「导入规则」。</p>
        <el-button type="primary" size="small" style="margin-top: 12px" @click="importBuiltin">同步内置规则</el-button>
        <el-button size="small" style="margin-top: 12px; margin-left: 8px" @click="importDialogVisible = true">导入规则</el-button>
      </div>
    </div>

    <el-dialog title="规则详情" :visible.sync="detailVisible" width="720px" custom-class="theme-dialog" @closed="detailRule = null">
      <div v-if="detailRule" class="rule-detail">
        <p><span class="k">ID</span>{{ detailRule.id }} / 编号 {{ detailRule.ruleCode }}</p>
        <p><span class="k">检查项</span>{{ detailRule.name }}</p>
        <p><span class="k">数据库</span>{{ detailRule.dbTypeName }}</p>
        <p><span class="k">分类</span>{{ detailRule.categoryName }}</p>
        <p><span class="k">风险</span>{{ detailRule.riskName }}</p>
        <p><span class="k">匹配</span>{{ detailRule.matchType }} · 期望：{{ detailRule.expectedValue || '—' }}</p>
        <p><span class="k">说明</span>{{ detailRule.description || '—' }}</p>
        <p><span class="k">检查语句</span></p>
        <pre class="cmd-block">{{ formatQueries(detailRule.queries) }}</pre>
        <p><span class="k">修复建议</span>{{ detailRule.fixSuggestion || '—' }}</p>
      </div>
      <span slot="footer">
        <el-button type="primary" @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>

    <el-dialog :title="formMode === 'create' ? '新增规则' : '编辑规则'" :visible.sync="formDialogVisible" width="680px" custom-class="theme-dialog" @closed="resetForm">
      <div class="rule-form">
        <el-form label-width="100px" size="small">
          <el-form-item label="规则名称" required>
            <el-input v-model="ruleForm.name" placeholder="检查项名称" />
          </el-form-item>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="数据库">
                <el-select v-model="ruleForm.dbType" style="width: 100%">
                  <el-option :value="0" label="全部类型" />
                  <el-option :value="1" label="MySQL" />
                  <el-option :value="2" label="PostgreSQL" />
                  <el-option :value="3" label="MongoDB" />
                  <el-option :value="4" label="Redis" />
                  <el-option :value="5" label="CouchDB" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="风险等级">
                <el-select v-model="ruleForm.risk" style="width: 100%">
                  <el-option :value="0" label="严重" />
                  <el-option :value="1" label="高危" />
                  <el-option :value="2" label="中危" />
                  <el-option :value="3" label="低危" />
                  <el-option :value="4" label="信息" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="检测分类">
                <el-select v-model="ruleForm.category" style="width: 100%">
                  <el-option v-for="c in categoryOptions" :key="c.value" :label="c.label" :value="c.value" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="匹配方式">
                <el-select v-model="ruleForm.matchType" style="width: 100%">
                  <el-option value="contains" label="包含" />
                  <el-option value="exact" label="精确" />
                  <el-option value="not_contains" label="不包含" />
                  <el-option value="empty" label="结果为空" />
                  <el-option value="always" label="始终通过" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="期望值">
            <el-input v-model="ruleForm.expectedValue" placeholder="匹配参考值说明" />
          </el-form-item>
          <el-form-item label="检查语句">
            <el-input
              v-model="ruleForm.queriesText"
              type="textarea"
              :rows="4"
              placeholder="每行一条 SQL 或检查命令"
            />
          </el-form-item>
          <el-form-item label="规则描述">
            <el-input v-model="ruleForm.description" type="textarea" :rows="2" />
          </el-form-item>
          <el-form-item label="修复建议">
            <el-input v-model="ruleForm.fixSuggestion" type="textarea" :rows="2" />
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer">
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="formSaving" @click="saveRule">保存</el-button>
      </span>
    </el-dialog>

    <el-dialog title="导入规则" :visible.sync="importDialogVisible" width="640px" custom-class="theme-dialog" @closed="resetImport">
      <div class="import-body">
        <p class="import-tip">选择 JSON 文件批量导入（支持大文件，建议分批）。重复项（名称+分类+数据库）将自动跳过。</p>
        <pre class="import-format">[
  {
    "name": "检查匿名账户",
    "description": "应删除匿名账户",
    "category": 2,
    "risk": 1,
    "dbType": 1,
    "queries": ["SELECT user FROM mysql.user WHERE user=''"],
    "expectedValue": "空结果",
    "matchType": "empty",
    "fixSuggestion": "DROP USER ''@'localhost'"
  }
]</pre>
        <p class="import-cat-hint">分类：1身份认证 2权限 3配置 4审计 5网络 6加密 7SQL注入 8敏感数据</p>
        <div class="import-file-row">
          <input ref="fileInput" type="file" accept=".json" @change="onFileChange" class="import-file-input" />
          <el-button size="small" @click="$refs.fileInput.click()">选择文件</el-button>
          <span class="import-file-name">{{ importFileName || '未选择文件' }}</span>
        </div>
        <div v-if="importPreview" class="import-preview">
          <p>解析到 <strong>{{ importPreview.length }}</strong> 条规则</p>
          <el-table :data="importPreview.slice(0, 5)" size="small" class="myTable" max-height="200">
            <el-table-column prop="name" label="名称" min-width="140" :show-overflow-tooltip="true" />
            <el-table-column prop="dbType" label="库类型" width="80" />
            <el-table-column prop="category" label="分类" width="60" />
          </el-table>
          <p v-if="importPreview.length > 5" class="import-more">...还有 {{ importPreview.length - 5 }} 条</p>
        </div>
      </div>
      <span slot="footer">
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="importing" :disabled="!importPreview" @click="doImport">确认导入</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'

export default {
  name: 'DataDetectionRules',
  data() {
    return {
      rulesLoading: false,
      rulesPayload: null,
      ruleFilterDb: null,
      ruleFilterCat: null,
      ruleKeyword: '',
      detailVisible: false,
      detailRule: null,
      importDialogVisible: false,
      importing: false,
      importFileName: '',
      importPreview: null,
      importRawData: null,
      cveImporting: false,
      cvePreviewLoading: true,
      cvePreview: { availableInDb: 0, targetTotal: 3800 },
      currentPage: 1,
      pageSize: 50,
      formDialogVisible: false,
      formMode: 'create',
      formSaving: false,
      ruleForm: {
        id: 0,
        name: '',
        description: '',
        category: 1,
        risk: 1,
        dbType: 1,
        queriesText: '',
        expectedValue: '',
        matchType: 'contains',
        fixSuggestion: '',
        riskDescription: ''
      },
      categoryOptions: [
        { value: 1, label: '身份认证' },
        { value: 2, label: '权限控制' },
        { value: 3, label: '配置安全' },
        { value: 4, label: '审计日志' },
        { value: 5, label: '网络安全' },
        { value: 6, label: '加密' },
        { value: 7, label: 'SQL 注入' },
        { value: 8, label: '敏感数据识别' }
      ]
    }
  },
  computed: {
    builtinHint() {
      return (this.rulesPayload && this.rulesPayload.builtinTotal) || 24
    },
    enabledCount() {
      const rows = (this.rulesPayload && this.rulesPayload.rules) || []
      return rows.filter((r) => r.enabled === 1).length
    },
    progressPercent() {
      const target = (this.rulesPayload && this.rulesPayload.targetTotal) || 3800
      if (!target) return 0
      return Math.min(100, Math.round((this.enabledCount / target) * 100))
    },
    rulesDbFilterOptions() {
      const rows = (this.rulesPayload && this.rulesPayload.byDbType) || []
      return rows.map((r) => ({ value: r.dbType, label: r.dbTypeName }))
    },
    rulesCatFilterOptions() {
      const rows = (this.rulesPayload && this.rulesPayload.byCategory) || []
      return rows.map((r) => ({ value: r.category, label: r.categoryName }))
    },
    filteredRules() {
      const rows = (this.rulesPayload && this.rulesPayload.rules) || []
      const db = this.ruleFilterDb
      const cat = this.ruleFilterCat
      const kw = (this.ruleKeyword || '').trim().toLowerCase()
      return rows.filter((r) => {
        if (db !== null && db !== undefined && db !== '' && r.dbType !== db) return false
        if (cat !== null && cat !== undefined && cat !== '' && r.category !== cat) return false
        if (!kw) return true
        const blob = [r.name, r.description, r.fixSuggestion, r.riskDescription].filter(Boolean).join(' ').toLowerCase()
        return blob.includes(kw)
      })
    },
    pagedRules() {
      const start = (this.currentPage - 1) * this.pageSize
      return this.filteredRules.slice(start, start + this.pageSize)
    },
    cveAvailable() {
      return Number(this.cvePreview && this.cvePreview.availableInDb) || 0
    }
  },
  watch: {
    ruleFilterDb() { this.currentPage = 1 },
    ruleFilterCat() { this.currentPage = 1 },
    ruleKeyword() { this.currentPage = 1 }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/datasec/rules'
    this.loadRules()
    this.loadCvePreview()
  },
  methods: {
    async loadCvePreview() {
      this.cvePreviewLoading = true
      try {
        const res = await security.previewDatasecCveImport()
        if (res.code === 200 && res.data) {
          this.cvePreview = res.data
          if (this.cveAvailable <= 0 && res.data.message) {
            this.$message({ message: res.data.message, type: 'warning', duration: 5000 })
          }
        } else if (res.data) {
          this.cvePreview = res.data
          this.$message({ message: res.msg || res.data.message || 'CVE 库不可用', type: 'warning', duration: 5000 })
        } else {
          this.$message({ message: res.msg || 'CVE 库预览失败', type: 'warning' })
        }
      } catch (e) {
        this.$message({ message: 'CVE 库预览失败: ' + (e.message || ''), type: 'warning' })
      } finally {
        this.cvePreviewLoading = false
      }
    },
    async loadRules() {
      this.rulesLoading = true
      try {
        const res = await security.getDatasecRules()
        if (res.code === 200 && res.data) {
          this.rulesPayload = res.data
        } else {
          this.rulesPayload = null
          this.$message({ message: res.msg || '加载失败', type: 'error' })
        }
      } catch (e) {
        this.rulesPayload = null
        this.$message({ message: '加载失败: ' + (e.message || ''), type: 'error' })
      } finally {
        this.rulesLoading = false
      }
    },
    openDetail(row) {
      this.detailRule = row
      this.detailVisible = true
    },
    openCreate() {
      this.formMode = 'create'
      this.resetForm()
      this.formDialogVisible = true
    },
    openEdit(row) {
      this.formMode = 'edit'
      this.ruleForm = {
        id: row.id,
        name: row.name || '',
        description: row.description || '',
        category: row.category || 1,
        risk: row.risk != null ? row.risk : 1,
        dbType: row.dbType != null ? row.dbType : 1,
        queriesText: (row.queries || []).join('\n'),
        expectedValue: row.expectedValue || '',
        matchType: row.matchType || 'contains',
        fixSuggestion: row.fixSuggestion || '',
        riskDescription: row.riskDescription || ''
      }
      this.formDialogVisible = true
    },
    confirmDelete(row) {
      this.$confirm(`确定删除规则「${row.name}」？`, '确认删除', { type: 'warning' })
        .then(async () => {
          const res = await security.deleteDatasecRule({ id: row.id })
          if (res.code === 200) {
            this.$message({ message: '删除成功', type: 'success' })
            this.loadRules()
          } else {
            this.$message({ message: res.msg || '删除失败', type: 'error' })
          }
        })
        .catch(() => {})
    },
    async saveRule() {
      if (!this.ruleForm.name) {
        this.$message({ message: '规则名称不能为空', type: 'warning' })
        return
      }
      this.formSaving = true
      try {
        const queries = this.ruleForm.queriesText.split('\n').map((s) => s.trim()).filter(Boolean)
        const payload = {
          name: this.ruleForm.name,
          description: this.ruleForm.description,
          category: this.ruleForm.category,
          risk: this.ruleForm.risk,
          dbType: this.ruleForm.dbType,
          queries,
          expectedValue: this.ruleForm.expectedValue,
          matchType: this.ruleForm.matchType,
          fixSuggestion: this.ruleForm.fixSuggestion,
          riskDescription: this.ruleForm.riskDescription,
          enabled: 1
        }
        const res =
          this.formMode === 'create'
            ? await security.createDatasecRule(payload)
            : await security.updateDatasecRule({ ...payload, id: this.ruleForm.id })
        if (res.code === 200) {
          this.$message({ message: '保存成功', type: 'success' })
          this.formDialogVisible = false
          this.loadRules()
        } else {
          this.$message({ message: res.msg || '保存失败', type: 'error' })
        }
      } finally {
        this.formSaving = false
      }
    },
    resetForm() {
      this.ruleForm = {
        id: 0,
        name: '',
        description: '',
        category: 1,
        risk: 1,
        dbType: 1,
        queriesText: '',
        expectedValue: '',
        matchType: 'contains',
        fixSuggestion: '',
        riskDescription: ''
      }
    },
    formatQueries(q) {
      if (!q || !q.length) return '—'
      return q.join('\n---\n')
    },
    async importFromCve() {
      if (this.cvePreviewLoading) {
        this.$message({ message: '正在检测 CVE 库，请稍候再试', type: 'info' })
        return
      }
      let avail = this.cveAvailable
      if (avail <= 0) {
        await this.loadCvePreview()
        avail = this.cveAvailable
      }
      if (avail <= 0) {
        this.$message({
          message: (this.cvePreview && this.cvePreview.message) || 'CVE 库不可用，请确认 data/default-cve.db 存在并已重启后端',
          type: 'warning',
          duration: 5000
        })
        return
      }
      const hint = `将从 CVE 库导入最多 ${Math.min(avail, 3800)} 条数据库相关 CVE 为检测规则（知识库型，不参与 SQL 执行）。已存在的将跳过。`
      try {
        await this.$confirm(hint, '从 CVE 库导入', { type: 'info' })
      } catch (_) {
        return
      }
      this.cveImporting = true
      try {
        const res = await security.importDatasecRulesFromCve({ limit: 3800 })
        if (res.code === 200 && res.data) {
          this.$message({
            message: `CVE 导入完成：成功 ${res.data.success} 条，跳过 ${res.data.skipped} 条`,
            type: 'success',
            duration: 6000
          })
          this.loadRules()
        } else {
          this.$message({ message: res.msg || '导入失败', type: 'error' })
        }
      } finally {
        this.cveImporting = false
      }
    },
    async importBuiltin() {
      try {
        const res = await security.importDatasecBuiltinRules()
        if (res.code === 200 && res.data) {
          this.$message({
            message: `内置规则同步：成功 ${res.data.success} 条，跳过 ${res.data.skipped} 条`,
            type: 'success',
            duration: 5000
          })
          this.loadRules()
        } else {
          this.$message({ message: res.msg || '同步失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: e.message || '请求失败', type: 'error' })
      }
    },
    onFileChange(e) {
      const file = e.target.files[0]
      if (!file) return
      this.importFileName = file.name
      const reader = new FileReader()
      reader.onload = (ev) => {
        try {
          const data = JSON.parse(ev.target.result)
          if (!Array.isArray(data)) {
            this.$message({ message: '需要 JSON 数组', type: 'error' })
            return
          }
          this.importRawData = data
          this.importPreview = data
        } catch (err) {
          this.$message({ message: 'JSON 解析失败: ' + err.message, type: 'error' })
          this.importPreview = null
        }
      }
      reader.readAsText(file)
    },
    resetImport() {
      this.importFileName = ''
      this.importPreview = null
      this.importRawData = null
      if (this.$refs.fileInput) this.$refs.fileInput.value = ''
    },
    onPageSizeChange(size) {
      this.pageSize = size
      this.currentPage = 1
    },
    onPageChange(page) {
      this.currentPage = page
    },
    async doImport() {
      if (!this.importRawData || !this.importRawData.length) return
      this.importing = true
      try {
        const res = await security.importDatasecRules({ rules: this.importRawData })
        if (res.code === 200 && res.data) {
          this.$message({
            message: `导入完成：共 ${res.data.total} 条，成功 ${res.data.success} 条，跳过 ${res.data.skipped} 条`,
            type: 'success',
            duration: 6000
          })
          this.importDialogVisible = false
          this.loadRules()
        } else {
          this.$message({ message: res.msg || '导入失败', type: 'error' })
        }
      } finally {
        this.importing = false
      }
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

.link-to-tasks {
  margin-left: 10px;
  color: #00d4aa;
}

.rules-header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.rules-summary-block {
  flex: 1;
  min-width: 280px;
}

.rules-summary {
  color: #94a3b8;
  font-size: 13px;
  margin: 0 0 8px;
}

.target-progress {
  max-width: 420px;
}

.rules-header-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.rules-summary-row {
  margin-bottom: 14px;
}

.rules-panel-title {
  font-size: 13px;
  color: #cbd5e1;
  margin-bottom: 8px;
}

.rules-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 10px;
}

.rules-pagination {
  margin-top: 12px;
  text-align: right;
}

.rules-footnote {
  margin-top: 12px;
  font-size: 12px;
  color: #64748b;
  line-height: 1.5;
}

.detail-loading,
.detail-empty {
  color: #94a3b8;
  padding: 24px;
  text-align: center;
}

.rule-detail {
  color: #cbd5e1;
  font-size: 14px;
  line-height: 1.6;
}

.rule-detail .k {
  display: inline-block;
  width: 88px;
  color: #94a3b8;
}

.cmd-block {
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid rgba(148, 163, 184, 0.2);
  border-radius: 8px;
  padding: 10px;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 200px;
  overflow: auto;
}

.tag-on {
  color: #10b981;
}
.tag-off {
  color: #94a3b8;
}

.import-body {
  color: #cbd5e1;
}

.import-tip,
.import-cat-hint {
  font-size: 13px;
  color: #94a3b8;
  margin: 0 0 10px;
}

.import-format {
  background: rgba(15, 23, 42, 0.85);
  border-radius: 8px;
  padding: 10px;
  font-size: 11px;
  max-height: 160px;
  overflow: auto;
  margin-bottom: 12px;
}

.import-file-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.import-file-input {
  display: none;
}

.import-file-name {
  font-size: 13px;
  color: #94a3b8;
}

.import-more {
  font-size: 12px;
  color: #64748b;
  margin-top: 6px;
}

code {
  color: #00d4aa;
  font-size: 12px;
}
</style>
