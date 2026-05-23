<template>
  <div class="security-container">
    <p class="page-intro">
      规则来源：<strong>数据库 host_baseline_rule 表</strong>。支持新增、编辑、删除规则，导入的规则也会写入该表。
      <router-link class="link-to-tasks" to="/hostsec/tasks">前往任务管理</router-link>
    </p>

    <div class="list_box">
      <div v-if="rulesLoading" class="detail-loading">加载中…</div>
      <div v-else-if="rulesPayload">
        <div class="rules-header">
          <p class="rules-summary">
            已加载规则 <strong>{{ rulesPayload.total }}</strong> 条。
          </p>
          <div class="rules-header-actions">
            <el-button size="small" icon="el-icon-plus" type="success" @click="openCreate">新增规则</el-button>
            <el-button type="primary" size="small" icon="el-icon-upload2" @click="importDialogVisible = true">导入规则</el-button>
          </div>
        </div>
        <el-row :gutter="16" class="rules-summary-row">
          <el-col :span="12">
            <div class="rules-panel-title">按操作系统</div>
            <el-table :data="rulesPayload.byOsType || []" size="small" class="myTable" max-height="220">
              <el-table-column prop="osTypeName" label="类型" />
              <el-table-column prop="count" label="规则数" width="90" />
            </el-table>
          </el-col>
          <el-col :span="12">
            <div class="rules-panel-title">按核查分类</div>
            <el-table :data="rulesPayload.byCategory || []" size="small" class="myTable" max-height="220">
              <el-table-column prop="categoryName" label="分类" :show-overflow-tooltip="true" />
              <el-table-column prop="count" label="规则数" width="90" />
            </el-table>
          </el-col>
        </el-row>
        <div class="rules-toolbar">
          <el-select v-model="ruleFilterOs" clearable placeholder="全部操作系统" size="small" style="width: 160px">
            <el-option v-for="o in rulesOsFilterOptions" :key="o.value" :label="o.label" :value="o.value" />
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
          <el-table-column prop="osTypeName" label="适用 OS" width="110" />
          <el-table-column prop="categoryName" label="分类" width="120" :show-overflow-tooltip="true" />
          <el-table-column prop="riskName" label="风险" width="80" />
          <el-table-column prop="name" label="检查项" min-width="140" :show-overflow-tooltip="true" />
          <el-table-column prop="description" label="说明" min-width="160" :show-overflow-tooltip="true" />
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
      </div>
      <div v-else class="detail-empty">暂无数据</div>
    </div>

    <el-dialog title="规则详情（只读）" :visible.sync="detailVisible" width="720px" custom-class="theme-dialog" @closed="detailRule = null">
      <div v-if="detailRule" class="rule-detail">
        <p><span class="k">ID</span>{{ detailRule.id }}</p>
        <p><span class="k">检查项</span>{{ detailRule.name }}</p>
        <p><span class="k">适用 OS</span>{{ detailRule.osTypeName }}</p>
        <p><span class="k">分类</span>{{ detailRule.categoryName }}</p>
        <p><span class="k">风险</span>{{ detailRule.riskName }}</p>
        <p v-if="detailRule.matchType"><span class="k">匹配方式</span>{{ detailRule.matchType }}</p>
        <p><span class="k">说明</span>{{ detailRule.description || '—' }}</p>
        <p v-if="detailRule.riskDescription"><span class="k">风险说明</span>{{ detailRule.riskDescription }}</p>
        <p><span class="k">期望值</span><code class="mono">{{ detailRule.expectedValue || '—' }}</code></p>
        <p><span class="k">检查命令</span></p>
        <pre class="cmd-block">{{ formatCommands(detailRule.commands) }}</pre>
        <p><span class="k">修复建议</span>{{ detailRule.fixSuggestion || '—' }}</p>
      </div>
      <span slot="footer">
        <el-button type="primary" @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>

    <el-dialog :title="formMode === 'create' ? '新增规则' : '编辑规则'" :visible.sync="formDialogVisible" width="680px" custom-class="theme-dialog" @closed="resetForm">
      <div class="rule-form">
        <el-form ref="ruleForm" :model="ruleForm" label-width="100px" size="small">
          <el-form-item label="规则名称" required>
            <el-input v-model="ruleForm.name" placeholder="请输入规则名称" />
          </el-form-item>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="操作系统">
                <el-select v-model="ruleForm.osType" style="width: 100%">
                  <el-option :value="1" label="Linux 服务器" />
                  <el-option :value="2" label="Windows" />
                  <el-option :value="3" label="Linux 终端" />
                  <el-option :value="4" label="网络设备" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="风险等级">
                <el-select v-model="ruleForm.risk" style="width: 100%">
                  <el-option :value="1" label="高危" />
                  <el-option :value="2" label="中危" />
                  <el-option :value="3" label="低危" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="核查分类">
                <el-select v-model="ruleForm.category" style="width: 100%">
                  <el-option :value="1" label="身份鉴别" />
                  <el-option :value="2" label="访问控制" />
                  <el-option :value="3" label="安全审计" />
                  <el-option :value="4" label="入侵防范" />
                  <el-option :value="5" label="资源控制" />
                  <el-option :value="6" label="密码安全" />
                  <el-option :value="7" label="网络连接" />
                  <el-option :value="8" label="文件权限" />
                  <el-option :value="9" label="日志配置" />
                  <el-option :value="10" label="系统更新" />
                  <el-option :value="11" label="服务管理" />
                  <el-option :value="12" label="内核参数" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="匹配方式">
                <el-select v-model="ruleForm.matchType" style="width: 100%">
                  <el-option value="contains" label="包含（contains）" />
                  <el-option value="exact" label="精确（exact）" />
                  <el-option value="regex" label="正则（regex）" />
                  <el-option value="not_contains" label="不包含（not_contains）" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="期望值">
            <el-input v-model="ruleForm.expectedValue" placeholder="命令执行后续期望匹配的值" />
          </el-form-item>
          <el-form-item label="检查命令">
            <el-input
              v-model="ruleForm.commandsText"
              type="textarea"
              :rows="4"
              placeholder="每行一条命令，如：&#10;grep '^PASS_MAX_DAYS' /etc/login.defs&#10;grep '^PASS_MIN_DAYS' /etc/login.defs"
            />
          </el-form-item>
          <el-form-item label="规则描述">
            <el-input v-model="ruleForm.description" type="textarea" :rows="3" placeholder="规则描述" />
          </el-form-item>
          <el-form-item label="风险说明">
            <el-input v-model="ruleForm.riskDescription" type="textarea" :rows="2" placeholder="风险说明" />
          </el-form-item>
          <el-form-item label="修复建议">
            <el-input v-model="ruleForm.fixSuggestion" type="textarea" :rows="2" placeholder="修复建议" />
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer">
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="formSaving" @click="saveRule">保存</el-button>
      </span>
    </el-dialog>

    <el-dialog title="导入规则" :visible.sync="importDialogVisible" width="600px" custom-class="theme-dialog" @closed="resetImport">
      <div class="import-body">
        <p class="import-tip">选择 JSON 规则文件导入，文件格式要求：</p>
        <pre class="import-format">[
  {
    "name": "检查项名称",
    "description": "描述",
    "category": 1,
    "risk": 1,
    "osType": 1,
    "commands": ["command1", "command2"],
    "expectedValue": "期望值",
    "matchType": "contains",
    "fixSuggestion": "修复建议",
    "riskDescription": "风险说明"
  }
]</pre>
        <div class="import-file-row">
          <input ref="fileInput" type="file" accept=".json" @change="onFileChange" class="import-file-input" />
          <el-button size="small" @click="$refs.fileInput.click()">选择文件</el-button>
          <span class="import-file-name">{{ importFileName || '未选择文件' }}</span>
        </div>
        <div v-if="importPreview" class="import-preview">
          <p>解析到 <strong>{{ importPreview.length }}</strong> 条规则</p>
          <el-table :data="importPreview.slice(0, 5)" size="small" class="myTable" max-height="200">
            <el-table-column prop="name" label="名称" min-width="140" :show-overflow-tooltip="true" />
            <el-table-column prop="osType" label="OS类型" width="80" />
            <el-table-column prop="risk" label="风险" width="60" />
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
  name: 'HostDetectionRules',
  data() {
    return {
      rulesLoading: false,
      rulesPayload: null,
      ruleFilterOs: null,
      ruleFilterCat: null,
      ruleKeyword: '',
      detailVisible: false,
      detailRule: null,
      importDialogVisible: false,
      importing: false,
      importFileName: '',
      importPreview: null,
      importRawData: null,
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
        osType: 1,
        commandsText: '',
        expectedValue: '',
        matchType: 'contains',
        fixSuggestion: '',
        riskDescription: ''
      }
    }
  },
  computed: {
    rulesOsFilterOptions() {
      const rows = (this.rulesPayload && this.rulesPayload.byOsType) || []
      return rows.map((r) => ({ value: r.osType, label: r.osTypeName }))
    },
    rulesCatFilterOptions() {
      const rows = (this.rulesPayload && this.rulesPayload.byCategory) || []
      return rows.map((r) => ({ value: r.category, label: r.categoryName }))
    },
    filteredRules() {
      const rows = (this.rulesPayload && this.rulesPayload.rules) || []
      const os = this.ruleFilterOs
      const cat = this.ruleFilterCat
      const kw = (this.ruleKeyword || '').trim().toLowerCase()
      return rows.filter((r) => {
        if (os !== null && os !== undefined && os !== '' && r.osType !== os) return false
        if (cat !== null && cat !== undefined && cat !== '' && r.category !== cat) return false
        if (!kw) return true
        const blob = [r.name, r.description, r.fixSuggestion, r.riskDescription].filter(Boolean).join(' ').toLowerCase()
        return blob.includes(kw)
      })
    },
    pagedRules() {
      const start = (this.currentPage - 1) * this.pageSize
      return this.filteredRules.slice(start, start + this.pageSize)
    }
  },
  watch: {
    ruleFilterOs() { this.currentPage = 1 },
    ruleFilterCat() { this.currentPage = 1 },
    ruleKeyword() { this.currentPage = 1 }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/hostsec/rules'
    this.loadRules()
  },
  methods: {
    async loadRules() {
      this.rulesLoading = true
      this.rulesPayload = null
      try {
        const res = await security.getBaselineRulesFromDB()
        if (res.code === 200 && res.data) {
          this.rulesPayload = res.data
        } else {
          this.$message({ message: res.msg || '加载失败', type: 'error' })
        }
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
      this.ruleForm = {
        id: 0,
        name: '',
        description: '',
        category: 1,
        risk: 1,
        osType: 1,
        commandsText: '',
        expectedValue: '',
        matchType: 'contains',
        fixSuggestion: '',
        riskDescription: ''
      }
      this.formDialogVisible = true
    },
    openEdit(row) {
      this.formMode = 'edit'
      this.ruleForm = {
        id: row.id,
        name: row.name || '',
        description: row.description || '',
        category: row.category || 1,
        risk: row.risk || 1,
        osType: row.osType || 1,
        commandsText: (row.commands || []).join('\n'),
        expectedValue: row.expectedValue || '',
        matchType: row.matchType || 'contains',
        fixSuggestion: row.fixSuggestion || '',
        riskDescription: row.riskDescription || ''
      }
      this.formDialogVisible = true
    },
    confirmDelete(row) {
      this.$confirm(`确定删除规则「${row.name}」？删除后不可恢复。`, '确认删除', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          const res = await security.deleteBaselineRule({ id: row.id })
          if (res.code === 200) {
            this.$message({ message: '删除成功', type: 'success' })
            this.loadRules()
          } else {
            this.$message({ message: res.msg || '删除失败', type: 'error' })
          }
        } catch (err) {
          this.$message({ message: '删除请求失败: ' + (err.message || ''), type: 'error' })
        }
      }).catch(() => {})
    },
    async saveRule() {
      if (!this.ruleForm.name) {
        this.$message({ message: '规则名称不能为空', type: 'warning' })
        return
      }
      this.formSaving = true
      try {
        const commands = this.ruleForm.commandsText
          .split('\n')
          .map((s) => s.trim())
          .filter(Boolean)
        const payload = {
          name: this.ruleForm.name,
          description: this.ruleForm.description,
          category: this.ruleForm.category,
          risk: this.ruleForm.risk,
          osType: this.ruleForm.osType,
          commands: commands,
          expectedValue: this.ruleForm.expectedValue,
          matchType: this.ruleForm.matchType,
          fixSuggestion: this.ruleForm.fixSuggestion,
          riskDescription: this.ruleForm.riskDescription
        }
        let res
        if (this.formMode === 'create') {
          res = await security.createBaselineRule(payload)
        } else {
          payload.id = this.ruleForm.id
          res = await security.updateBaselineRule(payload)
        }
        if (res.code === 200) {
          this.$message({
            message: this.formMode === 'create' ? '新增成功' : '保存成功',
            type: 'success'
          })
          this.formDialogVisible = false
          this.loadRules()
        } else {
          this.$message({ message: res.msg || '操作失败', type: 'error' })
        }
      } catch (err) {
        this.$message({ message: '操作失败: ' + (err.message || ''), type: 'error' })
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
        osType: 1,
        commandsText: '',
        expectedValue: '',
        matchType: 'contains',
        fixSuggestion: '',
        riskDescription: ''
      }
    },
    formatCommands(cmds) {
      if (!cmds || !cmds.length) return '—'
      return cmds.join('\n---\n')
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
            this.$message({ message: '文件格式错误：需要 JSON 数组', type: 'error' })
            return
          }
          this.importRawData = data
          this.importPreview = data
        } catch (err) {
          this.$message({ message: 'JSON 解析失败: ' + err.message, type: 'error' })
          this.importPreview = null
          this.importRawData = null
        }
      }
      reader.readAsText(file)
    },
    onPageSizeChange(size) {
      this.pageSize = size
      this.currentPage = 1
    },
    onPageChange(page) {
      this.currentPage = page
    },
    resetImport() {
      this.importFileName = ''
      this.importPreview = null
      this.importRawData = null
      if (this.$refs.fileInput) this.$refs.fileInput.value = ''
    },
    async doImport() {
      if (!this.importRawData || !this.importRawData.length) return
      this.importing = true
      try {
        const res = await security.importBaselineRules({ rules: this.importRawData })
        if (res.code === 200 && res.data) {
          this.$message({
            message: `导入完成：共 ${res.data.total} 条，成功 ${res.data.success} 条，跳过 ${res.data.skipped} 条`,
            type: 'success',
            duration: 5000
          })
          this.importDialogVisible = false
          this.loadRules()
        } else {
          this.$message({ message: res.msg || '导入失败', type: 'error' })
        }
      } catch (err) {
        this.$message({ message: '导入请求失败: ' + (err.message || ''), type: 'error' })
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

.rules-summary {
  color: #94a3b8;
  font-size: 13px;
  margin: 0 0 12px;
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

.detail-loading,
.detail-empty {
  color: #94a3b8;
  padding: 24px;
  text-align: center;
}

.rule-detail {
  color: #cbd5e1;
  font-size: 14px;
  line-height: 1.65;
}

.rule-detail p {
  margin: 0 0 10px;
}

.rule-detail .k {
  display: inline-block;
  min-width: 88px;
  color: #94a3b8;
  margin-right: 8px;
}

.mono {
  display: inline-block;
  margin-top: 4px;
  padding: 6px 10px;
  background: rgba(15, 23, 42, 0.6);
  border-radius: 4px;
  font-size: 13px;
  word-break: break-all;
  white-space: pre-wrap;
}

.cmd-block {
  margin: 0 0 12px;
  padding: 12px;
  background: rgba(15, 23, 42, 0.75);
  border-radius: 6px;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
  color: #e2e8f0;
  max-height: 220px;
  overflow: auto;
}

.rules-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.rules-header .rules-summary {
  margin: 0;
}

.rules-header-actions {
  display: flex;
  gap: 8px;
}

.import-body {
  color: #cbd5e1;
  font-size: 14px;
}

.import-tip {
  margin: 0 0 10px;
  color: #94a3b8;
}

.import-format {
  margin: 0 0 16px;
  padding: 12px;
  background: rgba(15, 23, 42, 0.75);
  border-radius: 6px;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
  color: #e2e8f0;
  max-height: 200px;
  overflow: auto;
}

.import-file-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.import-file-input {
  display: none;
}

.import-file-name {
  color: #94a3b8;
  font-size: 13px;
}

.import-preview {
  margin-top: 12px;
  padding: 12px;
  background: rgba(15, 23, 42, 0.4);
  border-radius: 6px;
}

.import-preview p {
  margin: 0 0 8px;
  color: #94a3b8;
}

.import-more {
  margin-top: 8px !important;
  font-size: 12px;
  color: #64748b;
}

.rules-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.rule-form {
  color: #cbd5e1;
}

.link_danger {
  color: #f56c6c;
}

.link_danger:hover {
  color: #f89898;
}
</style>