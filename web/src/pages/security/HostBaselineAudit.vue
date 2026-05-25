<template>
  <div class="security-container">
    <div class="main-title">安全配置核查</div>

    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="default" size="small" @click="openRulesLibrary">规则库与分类</el-button>
          <el-button type="primary" size="small" @click="openDialog">新建核查</el-button>
        </div>
      </div>

      <el-table :data="tableData" style="width: 100%" class="myTable">
        <el-table-column prop="targetIp" label="目标主机" :show-overflow-tooltip="true" />
        <el-table-column prop="osTypeName" label="操作系统类型" width="140" />
        <el-table-column prop="totalRules" label="检查项数" width="100" />
        <el-table-column prop="passCount" label="通过" width="80" />
        <el-table-column prop="failCount" label="不通过" width="90" />
        <el-table-column prop="errorCount" label="异常" width="80" />
        <el-table-column prop="checkTime" label="最近核查时间" width="170" />
        <el-table-column label="操作" width="100">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="openDetail(scope.row)">详情</el-link>
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

    <el-dialog title="新建安全配置核查" :visible.sync="dialogVisible" width="700px" custom-class="theme-dialog">
      <p class="form-hint">
        选择<strong>连接方式</strong>：自动模式下 Windows 走 <strong>WinRM</strong>（默认 5985/5986），Linux/国产/嵌入式走
        <strong>SSH</strong>（默认 22）。可添加多个目标主机，共享连接方式配置。
      </p>
      <el-form ref="batchForm" label-width="118px">
        <el-form-item label="连接方式">
          <el-select v-model="batchTransport" placeholder="请选择" style="width: 100%" @change="onBatchTransportChange">
            <el-option v-for="opt in transportOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="isBatchWinRM" label="WinRM HTTPS">
          <el-checkbox v-model="batchWinrmUseHttps" @change="onBatchWinrmHttpsChange">使用 HTTPS（典型端口 5986）</el-checkbox>
        </el-form-item>
      </el-form>

      <div class="targets-header">
        <span class="targets-title">目标列表（{{ targets.length }} 个）</span>
        <el-button type="primary" size="small" icon="el-icon-plus" @click="addTarget">添加目标</el-button>
      </div>

      <div v-for="(t, idx) in targets" :key="idx" class="target-card">
        <div class="target-card-header">
          <span class="target-card-index">目标 #{{ idx + 1 }}</span>
          <el-button type="danger" size="mini" icon="el-icon-delete" circle @click="removeTarget(idx)" />
        </div>
        <el-form :model="t" :rules="targetRules" ref="targetForm" label-width="100px" class="target-form">
          <el-row :gutter="12">
            <el-col :span="12">
              <el-form-item label="目标主机" prop="host">
                <el-input v-model="t.host" placeholder="IP 或域名" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="batchPortLabel" prop="port">
                <el-input-number v-model="t.port" :min="1" :max="65535" controls-position="right" style="width: 100%" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="12">
            <el-col :span="12">
              <el-form-item :label="batchUsernameLabel" prop="username">
                <el-input v-model="t.username" placeholder="登录用户名" autocomplete="off" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="batchPasswordLabel" prop="password">
                <el-input v-model="t.password" type="password" :placeholder="isBatchWinRM ? 'WinRM 密码（必填）' : '密码（与私钥二选一）'" show-password autocomplete="new-password" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="12">
            <el-col :span="12">
              <el-form-item label="操作系统" prop="osType">
                <el-select v-model="t.osType" placeholder="请选择" style="width: 100%">
                  <el-option v-for="opt in osTypes" :key="opt.value" :label="opt.label" :value="opt.value" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item v-if="!isBatchWinRM" label="SSH 私钥" prop="key">
                <el-input v-model="t.key" placeholder="可选：粘贴私钥 PEM" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </div>

      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitBatchForm">开始核查（{{ targets.length }} 个目标）</el-button>
      </span>
    </el-dialog>

    <el-dialog title="漏洞检测规则库" :visible.sync="rulesDialogVisible" width="960px" custom-class="theme-dialog">
      <div v-if="rulesLoading" class="detail-loading">加载中…</div>
      <div v-else-if="rulesPayload">
        <p class="rules-summary">
          当前引擎已加载检测规则 <strong>{{ rulesPayload.total }}</strong> 条（内置规则 + <code>data/baseline</code> 目录下 JSON）。
          需求中的「1200+」为规划容量时，以本页实际条数为准；可通过追加规则文件扩容。
        </p>
        <el-row :gutter="16" class="rules-summary-row">
          <el-col :span="12">
            <div class="rules-panel-title">按操作系统</div>
            <el-table :data="rulesPayload.byOsType || []" size="small" class="myTable" max-height="200">
              <el-table-column prop="osTypeName" label="类型" />
              <el-table-column prop="count" label="规则数" width="90" />
            </el-table>
          </el-col>
          <el-col :span="12">
            <div class="rules-panel-title">按核查分类</div>
            <el-table :data="rulesPayload.byCategory || []" size="small" class="myTable" max-height="200">
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
        <el-table :data="filteredRules" class="myTable" max-height="380" size="small">
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="osTypeName" label="适用 OS" width="110" />
          <el-table-column prop="categoryName" label="分类" width="120" :show-overflow-tooltip="true" />
          <el-table-column prop="riskName" label="风险" width="80" />
          <el-table-column prop="name" label="检查项" min-width="140" :show-overflow-tooltip="true" />
          <el-table-column prop="description" label="说明" min-width="160" :show-overflow-tooltip="true" />
          <el-table-column label="检查命令" min-width="140" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              {{ formatRuleCommands(scope.row.commands) }}
            </template>
          </el-table-column>
        </el-table>
        <p class="rules-footnote">分类管理：当前为只读规则目录（按 OS 与核查维度汇总）；若需启用/禁用或在线编辑规则，需后续增加配置存储与接口。</p>
      </div>
      <div v-else class="detail-empty">暂无数据</div>
      <span slot="footer">
        <el-button @click="rulesDialogVisible = false">关闭</el-button>
      </span>
    </el-dialog>

    <el-dialog title="核查详情" :visible.sync="detailVisible" width="900px" custom-class="theme-dialog">
      <div v-if="detailLoading" class="detail-loading">加载中…</div>
      <div v-else-if="detailRows.length">
        <p class="detail-meta">目标：{{ detailTarget }} · 共 {{ detailRows.length }} 条检查项</p>
        <el-table :data="detailRows" style="width: 100%" class="myTable" max-height="420">
          <el-table-column prop="categoryName" label="分类" width="120" />
          <el-table-column prop="ruleName" label="检查项" :show-overflow-tooltip="true" />
          <el-table-column prop="resultName" label="结果" width="90" />
          <el-table-column prop="riskName" label="风险" width="90" />
          <el-table-column prop="expectedValue" label="期望值" :show-overflow-tooltip="true" />
          <el-table-column prop="actualValue" label="实际值" :show-overflow-tooltip="true" />
          <el-table-column prop="fixSuggestion" label="修复建议" :show-overflow-tooltip="true" />
        </el-table>
      </div>
      <div v-else class="detail-empty">暂无明细</div>
      <span slot="footer">
        <el-button @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>

    <el-dialog title="任务执行中" :visible.sync="progressVisible" width="500px" custom-class="theme-dialog" :close-on-click-modal="false" :show-close="false">
      <div class="progress-body">
        <div class="progress-summary">
          <el-progress :percentage="progressPercent" :status="progressStatus" :stroke-width="16" />
          <p class="progress-text">{{ progressText }}</p>
        </div>
        <div v-for="(t, idx) in progressTargets" :key="idx" class="progress-target-row">
          <div class="progress-target-main">
            <span class="progress-target-host">{{ t.host }}</span>
            <el-tag :type="progressTagType(t.status)" size="mini">{{ progressTagLabel(t.status) }}</el-tag>
          </div>
          <p v-if="t.message" class="progress-target-msg">{{ t.message }}</p>
        </div>
      </div>
      <span slot="footer">
        <el-button v-if="progressDone" type="primary" @click="onProgressDone">查看结果</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'

const OS_WINDOWS = 2
const T_AUTO = 0
const T_SSH = 1
const T_WINRM = 2

export default {
  name: 'HostBaselineAudit',
  data() {
    return {
      dialogVisible: false,
      detailVisible: false,
      submitLoading: false,
      detailLoading: false,
      tableData: [],
      formData: {
        page: 1
      },
      pageSize: 10,
      currentpage: 1,
      totalpage: 0,
      osTypes: [],
      transportOptions: [],
      batchTransport: T_AUTO,
      batchWinrmUseHttps: false,
      targets: [],
      targetRules: {
        host: [{ required: true, message: '请输入目标主机', trigger: 'blur' }],
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        osType: [{ required: true, message: '请选择操作系统类型', trigger: 'change' }]
      },
      detailRows: [],
      detailTarget: '',
      rulesDialogVisible: false,
      rulesLoading: false,
      rulesPayload: null,
      ruleFilterOs: null,
      ruleFilterCat: null,
      ruleKeyword: '',
      progressVisible: false,
      progressTaskId: 0,
      progressTargets: [],
      progressPercent: 0,
      progressText: '',
      progressDone: false,
      progressTimer: null
    }
  },
  computed: {
    progressStatus() {
      if (!this.progressVisible) return ''
      if (this.progressDone) return 'success'
      return ''
    },
    isBatchWinRM() {
      const t = this.batchTransport
      if (t === T_WINRM) return true
      return false
    },
    batchPortLabel() {
      return this.isBatchWinRM ? '远程端口' : 'SSH 端口'
    },
    batchUsernameLabel() {
      return this.isBatchWinRM ? '用户名' : 'SSH 用户名'
    },
    batchPasswordLabel() {
      return this.isBatchWinRM ? '密码' : 'SSH 密码'
    },
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
        const blob = [r.name, r.description, r.fixSuggestion, r.riskDescription, this.formatRuleCommands(r.commands)]
          .filter(Boolean)
          .join(' ')
          .toLowerCase()
        return blob.includes(kw)
      })
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/hostsec/tasks'
    this.loadEnums()
    this.getData()
  },
  methods: {
    createTarget() {
      const port = this.isBatchWinRM ? (this.batchWinrmUseHttps ? 5986 : 5985) : 22
      return {
        host: '',
        port: port,
        username: '',
        password: '',
        key: '',
        osType: this.osTypes.length ? this.osTypes[0].value : 1
      }
    },
    addTarget() {
      this.targets.push(this.createTarget())
    },
    removeTarget(idx) {
      this.targets.splice(idx, 1)
    },
    onBatchTransportChange() {
      this.targets.forEach(t => {
        t.port = this.isBatchWinRM ? (this.batchWinrmUseHttps ? 5986 : 5985) : 22
      })
    },
    onBatchWinrmHttpsChange() {
      this.targets.forEach(t => {
        t.port = this.isBatchWinRM ? (this.batchWinrmUseHttps ? 5986 : 5985) : 22
      })
    },
    formatRuleCommands(cmds) {
      if (!cmds || !cmds.length) return '—'
      return cmds.join('；')
    },
    async openRulesLibrary() {
      this.rulesDialogVisible = true
      this.rulesLoading = true
      this.rulesPayload = null
      this.ruleFilterOs = null
      this.ruleFilterCat = null
      this.ruleKeyword = ''
      try {
        const res = await security.getBaselineRules()
        if (res.code === 200 && res.data) {
          this.rulesPayload = res.data
        } else {
          this.$message({ message: res.msg || '加载规则失败', type: 'error' })
        }
      } finally {
        this.rulesLoading = false
      }
    },
    async loadEnums() {
      const res = await security.getBaselineEnums()
      if (res.code === 200 && res.data) {
        if (res.data.osTypes) this.osTypes = res.data.osTypes
        if (res.data.hostTransports) this.transportOptions = res.data.hostTransports
      }
      if (!this.osTypes.length) {
        this.osTypes = [
          { value: 1, label: 'Linux/Unix' },
          { value: 2, label: 'Windows' },
          { value: 3, label: '国产操作系统' },
          { value: 4, label: '嵌入式 OS' }
        ]
      }
      if (!this.transportOptions.length) {
        this.transportOptions = [
          { value: T_AUTO, label: '自动（按操作系统）' },
          { value: T_SSH, label: 'SSH' },
          { value: T_WINRM, label: 'WinRM' }
        ]
      }
    },
    async getData() {
        const res = await security.getBaselineTaskList({
        page: this.formData.page,
        size: this.pageSize,
        scanScene: 1
      })
      if (res.code === 200) {
        this.tableData = res.data.list || []
        this.totalpage = res.data.total || 0
      } else {
        this.$message({ message: res.msg || '加载失败', type: 'error' })
      }
    },
    openDialog() {
      this.dialogVisible = true
      this.batchTransport = T_AUTO
      this.batchWinrmUseHttps = false
      this.targets = [this.createTarget()]
    },
    submitBatchForm() {
      if (!this.targets.length) {
        this.$message({ message: '请至少添加一个目标主机', type: 'warning' })
        return
      }
      for (let i = 0; i < this.targets.length; i++) {
        const t = this.targets[i]
        if (!t.host) {
          this.$message({ message: `目标 #${i + 1} 请输入主机地址`, type: 'warning' })
          return
        }
        if (!t.username) {
          this.$message({ message: `目标 #${i + 1} 请输入用户名`, type: 'warning' })
          return
        }
        if (this.isBatchWinRM) {
          if (!t.password) {
            this.$message({ message: `目标 #${i + 1} WinRM 需要填写密码`, type: 'warning' })
            return
          }
        } else if (!t.password && !t.key) {
          this.$message({ message: `目标 #${i + 1} 请填写 SSH 密码或私钥`, type: 'warning' })
          return
        }
      }
      this.submitLoading = true
      try {
        const payload = {
          targets: this.targets.map(t => ({
            host: t.host,
            port: t.port,
            username: t.username,
            password: t.password,
            key: this.isBatchWinRM ? '' : t.key,
            osType: t.osType,
            transport: this.batchTransport,
            winrmUseHttps: this.isBatchWinRM ? this.batchWinrmUseHttps : false,
            scanScene: 1
          }))
        }
        this.dialogVisible = false
        this.submitBatch(payload)
      } catch (e) {
        this.submitLoading = false
      }
    },
    async submitBatch(payload) {
      try {
        const res = await security.runBaselineBatchCheck(payload)
        if (res.code === 200) {
          this.startProgressPolling(res.data.taskId)
        } else {
          this.$message({ message: res.msg || '创建任务失败', type: 'error' })
        }
      } finally {
        this.submitLoading = false
      }
    },
    startProgressPolling(taskId) {
      this.progressTaskId = taskId
      this.progressTargets = []
      this.progressPercent = 0
      this.progressText = '任务已创建，正在后台执行…'
      this.progressDone = false
      this.progressVisible = true
      this.pollProgress()
    },
    async pollProgress() {
      if (!this.progressVisible || !this.progressTaskId) return
      try {
        const res = await security.getBaselineBatchProgress({ taskId: this.progressTaskId })
        if (res.code === 200 && res.data) {
          const d = res.data
          this.progressTargets = d.targets || []
          const total = d.totalTargets || 1
          const done = d.completedTargets || 0
          this.progressPercent = Math.round((done / total) * 100)
          this.progressText = `进度：${done}/${total} 个目标完成`
          if (d.status === 'completed') {
            this.progressDone = true
            this.progressPercent = 100
            this.progressText = `全部完成（${total} 个目标）`
            this.formData.page = 1
            this.currentpage = 1
            this.getData()
            return
          }
        }
      } catch (e) {
        console.error('pollProgress error:', e)
      }
      this.progressTimer = setTimeout(() => this.pollProgress(), 2000)
    },
    progressTagType(status) {
      if (status === 'completed') return 'success'
      if (status === 'failed') return 'danger'
      if (status === 'running') return 'warning'
      return 'info'
    },
    progressTagLabel(status) {
      if (status === 'completed') return '已完成'
      if (status === 'failed') return '执行失败'
      if (status === 'running') return '执行中'
      return '等待中'
    },
    onProgressDone() {
      this.progressVisible = false
      if (this.progressTimer) {
        clearTimeout(this.progressTimer)
        this.progressTimer = null
      }
    },
    async openDetail(row) {
      this.detailVisible = true
      this.detailLoading = true
      this.detailRows = []
      this.detailTarget = row.targetIp || ''
      try {
        const res = await security.getBaselineList({ taskId: row.taskId })
        if (res.code === 200) {
          this.detailRows = res.data.list || []
        } else {
          this.$message({ message: res.msg || '加载失败', type: 'error' })
        }
      } finally {
        this.detailLoading = false
      }
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
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

.form-hint {
  color: #94a3b8;
  font-size: 13px;
  line-height: 1.55;
  margin: 0 0 16px;
  padding: 10px 12px;
  background: rgba(0, 212, 170, 0.06);
  border-radius: 6px;
  border: 1px solid rgba(0, 212, 170, 0.15);
}

.field-hint {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: rgba(148, 163, 184, 0.75);
}

.detail-meta {
  color: #94a3b8;
  margin-bottom: 12px;
  font-size: 14px;
}

.detail-loading,
.detail-empty {
  color: #94a3b8;
  padding: 24px;
  text-align: center;
}

.rules-summary {
  color: #94a3b8;
  font-size: 13px;
  line-height: 1.55;
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

.rules-footnote {
  margin: 12px 0 0;
  font-size: 12px;
  color: rgba(148, 163, 184, 0.85);
  line-height: 1.5;
}

.targets-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 16px 0 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.15);
}

.targets-title {
  font-size: 14px;
  font-weight: 600;
  color: #cbd5e1;
}

.target-card {
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid rgba(148, 163, 184, 0.12);
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 12px;
}

.target-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.target-card-index {
  font-size: 13px;
  font-weight: 500;
  color: #00d4aa;
}

.target-form .el-form-item {
  margin-bottom: 12px;
}

.target-form .el-form-item:last-child {
  margin-bottom: 0;
}

.progress-body {
  padding: 8px 0;
}

.progress-summary {
  margin-bottom: 20px;
  text-align: center;
}

.progress-text {
  margin-top: 10px;
  font-size: 14px;
  color: #94a3b8;
}

.progress-target-row {
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.progress-target-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.progress-target-msg {
  margin: 6px 0 0;
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.5;
}

.progress-target-row:last-child {
  border-bottom: none;
}

.progress-target-host {
  font-size: 13px;
  color: #e2e8f0;
}
</style>
