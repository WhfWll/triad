<template>
  <div class="security-container mod-hub">
    <p class="page-intro">
      同一套远程连接（SSH / WinRM）能力，新建任务时区分：<strong>安全配置核查</strong>、<strong>主机漏洞检测</strong>（CVE 版本匹配）、<strong>恶意代码检测</strong>（YARA 规则引擎）。
    </p>

    <div class="tab-panel">
    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-tooltip :content="$store.state.systemAuthorized ? '' : '系统未授权，请前往「系统配置 → 系统授权」页面完成授权'" :disabled="$store.state.systemAuthorized" placement="bottom">
            <el-button type="primary" size="small" :disabled="!$store.state.systemAuthorized" @click="openCreateDialog">新建主机检查任务</el-button>
          </el-tooltip>
          <el-button
            size="small"
            :disabled="selectedRows.length === 0"
            :loading="batchDeleteLoading"
            @click="batchDeleteTasks"
          >批量删除</el-button>
        </div>
        <div class="serach-condition">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索目标 IP / 任务批次 / 任务类型"
            size="small"
            clearable
            class="task-search-input"
            @keydown.enter.native="applySearch"
          />
          <el-select v-model="taskTypeFilter" size="small" class="task-type-select" @change="onFilterChange">
            <el-option label="全部任务类型" value="all" />
            <el-option label="安全配置核查" value="baseline" />
            <el-option label="主机漏洞检测" value="vuln" />
            <el-option label="恶意代码检测" value="malware" />
          </el-select>
          <el-button type="primary" size="small" @click="applySearch">搜索</el-button>
          <el-button type="primary" size="small" @click="resetSearch">重置</el-button>
        </div>
      </div>

      <div class="table-scroll-wrap">
      <el-table
        ref="taskTable"
        v-loading="tableLoading"
        :data="displayRows"
        :row-key="taskRowKey"
        class="myTable"
        @selection-change="onSelectionChange"
      >
        <el-table-column type="selection" width="48" reserve-selection />
        <el-table-column prop="kindLabel" label="任务类型" width="130" />
        <el-table-column prop="targetIp" label="目标主机" min-width="110" :show-overflow-tooltip="true" />
        <el-table-column prop="osTypeName" label="操作系统" width="108" />
        <el-table-column label="结果摘要" min-width="200" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            {{ rowSummary(scope.row) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="88" align="center">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.isRunning" :type="runStatusTagType(scope.row.runStatus)" size="mini">{{ scope.row.runStatusLabel }}</el-tag>
            <el-tag v-else-if="(scope.row.source === 'vuln' || scope.row.source === 'malware') && scope.row.scanStatus === 2" type="danger" size="mini" effect="plain">异常</el-tag>
            <el-tag v-else type="success" size="mini" effect="plain">已完成</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="112" align="center" class-name="col-progress">
          <template slot-scope="scope">
            <span v-if="scope.row.isRunning" class="inline-progress">{{ scope.row.progressText }}</span>
            <span v-else class="inline-progress done">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="checkTime" label="时间" width="158" :show-overflow-tooltip="true" />
        <el-table-column label="操作" width="210">
          <template slot-scope="scope">
            <el-link v-if="scope.row.source === 'baseline' || scope.row.source === 'vuln' || scope.row.source === 'malware'" :underline="false" class="link_primary" @click="openDetailPage(scope.row)">详情</el-link>
            <el-link v-else :underline="false" class="link_primary" @click="openDetail(scope.row)">详情</el-link>
            <el-link v-if="scope.row.isRunning" :underline="false" class="link_primary" style="margin-left: 10px" @click="stopTask(scope.row)">结束</el-link>
            <el-link :underline="false" class="link_danger" style="margin-left: 10px" @click="deleteTask(scope.row)">删除</el-link>
          </template>
        </el-table-column>
      </el-table>
      </div>

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

    <el-dialog title="新建主机检查任务" :visible.sync="createVisible" width="700px" custom-class="theme-dialog" :close-on-click-modal="false">
      <el-form :model="taskForm" :rules="createRules" ref="taskForm" label-width="118px">
        <el-form-item label="任务类型" prop="taskKind">
          <el-radio-group v-model="taskForm.taskKind" @change="onTaskKindChange">
            <el-radio label="baseline">安全配置核查</el-radio>
            <el-radio label="vuln">主机漏洞检测（CVE）</el-radio>
            <el-radio label="malware">恶意代码检测</el-radio>
          </el-radio-group>
          <p class="field-hint">
            <template v-if="taskForm.taskKind === 'vuln'">SSH 登录后获取软件版本列表，匹配 CVE 数据库（25.8 万条记录）发现已知漏洞。</template>
            <template v-else-if="taskForm.taskKind === 'malware'">基于 YARA 规则引擎的恶意代码检测，支持 2000+ 条规则（Linux 恶意软件、Webshell、APT、挖矿木马等）。</template>
            <template v-else>对账号、补丁、防火墙、审计等进行基线核查。</template>
          </p>
        </el-form-item>
      </el-form>

      <div class="targets-header">
        <span class="targets-title">目标列表（{{ targets.length }} 个）</span>
        <el-button type="primary" size="small" icon="el-icon-plus" @click="openAddTargetDialog">添加目标</el-button>
      </div>

      <div v-if="targets.length === 0" class="targets-empty">
        <p>暂无目标，请点击上方按钮添加</p>
      </div>

      <div v-else class="targets-table">
        <div class="targets-table-header">
          <span class="th-index">序号</span>
          <span class="th-host">目标主机</span>
          <span class="th-protocol">协议</span>
          <span class="th-port">端口</span>
          <span class="th-username">用户名</span>
          <span class="th-os">操作系统</span>
          <span class="th-actions">操作</span>
        </div>
        <div v-for="(t, idx) in targets" :key="idx" class="targets-table-row">
          <span class="td-index">{{ idx + 1 }}</span>
          <span class="td-host" :title="t.host">{{ t.host }}</span>
          <span class="td-protocol">{{ t.transport === 'winrm' ? 'WinRM' : 'SSH' }}</span>
          <span class="td-port">{{ t.port }}</span>
          <span class="td-username">{{ t.username }}</span>
          <span class="td-os">{{ getOsTypeName(t.osType) }}</span>
          <span class="td-actions">
            <el-button type="text" size="mini" @click="editTarget(idx)">编辑</el-button>
            <el-button type="text" size="mini" style="color: #f56c6c" @click="removeTarget(idx)">删除</el-button>
          </span>
        </div>
      </div>

      <span slot="footer">
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitCreate">开始执行（{{ targets.length }} 个目标）</el-button>
      </span>
    </el-dialog>

    <el-dialog
      :title="editTargetIndex >= 0 ? '编辑目标' : '添加目标'"
      :visible.sync="targetDialogVisible"
      width="580px"
      custom-class="theme-dialog host-target-dialog"
      :close-on-click-modal="false"
      @closed="onTargetDialogClosed"
    >
      <el-form
        ref="editTargetForm"
        :model="editTargetForm"
        :rules="targetRules"
        label-position="top"
        size="small"
        class="host-target-form nessus-form"
      >
        <div class="form-section">
          <div class="section-title">连接信息</div>
          <el-row :gutter="16">
            <el-col :span="14">
              <el-form-item label="目标主机" prop="host">
                <el-input v-model="editTargetForm.host" placeholder="IP 或域名" @input="clearConnTestResult" />
              </el-form-item>
            </el-col>
            <el-col :span="10">
              <el-form-item label="连接协议" prop="transport">
                <el-select v-model="editTargetForm.transport" placeholder="请选择" style="width: 100%" @change="onEditTargetTransportChange">
                  <el-option :value="'ssh'" label="SSH" />
                  <el-option :value="'winrm'" label="WinRM" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="8">
              <el-form-item :label="editTargetPortLabel" prop="port">
                <el-input-number v-model="editTargetForm.port" :min="1" :max="65535" :controls="false" style="width: 100%" @change="clearConnTestResult" />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="操作系统" prop="osType">
                <el-select v-model="editTargetForm.osType" placeholder="请选择" style="width: 100%" @change="onEditTargetOsTypeChange">
                  <el-option v-for="opt in osTypes" :key="opt.value" :label="opt.label" :value="opt.value" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col v-if="editTargetForm.transport === 'winrm'" :span="8">
              <el-form-item label="HTTPS">
                <el-switch v-model="editTargetForm.winrmUseHttps" @change="onEditWinrmHttpsChange" />
                <span class="switch-hint">{{ editTargetForm.winrmUseHttps ? '5986' : '5985' }}</span>
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="form-section">
          <div class="section-title">认证信息</div>
          <p class="form-tip">
            <template v-if="editTargetForm.transport === 'winrm'">WinRM 使用用户名 + 密码认证。</template>
            <template v-else>SSH 支持密码或 PEM 私钥，二者至少填一项。</template>
          </p>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="用户名" prop="username">
                <el-input v-model="editTargetForm.username" placeholder="登录用户名" autocomplete="off" @input="clearConnTestResult" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="密码" prop="password">
                <el-input
                  v-model="editTargetForm.password"
                  type="password"
                  :placeholder="editTargetForm.transport === 'winrm' ? 'WinRM 密码（必填）' : '密码（与私钥二选一）'"
                  show-password
                  autocomplete="new-password"
                  @input="clearConnTestResult"
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item v-if="editTargetForm.transport !== 'winrm'" label="SSH 私钥（可选）">
            <el-input
              v-model="editTargetForm.key"
              type="textarea"
              :rows="3"
              placeholder="粘贴 PEM 格式私钥；留空则使用密码登录"
              @input="clearConnTestResult"
            />
          </el-form-item>
        </div>

        <div v-if="connTestResult" :class="['conn-test-result', connTestResult.ok ? 'ok' : 'fail']">
          <div class="result-line">
            <i :class="connTestResult.ok ? 'el-icon-success' : 'el-icon-error'"></i>
            <span>{{ connTestResult.message }}</span>
          </div>
          <div v-if="connTestResult.detail" class="result-detail">{{ connTestResult.detail }}</div>
        </div>
      </el-form>
      <span slot="footer" class="host-target-footer">
        <el-button @click="targetDialogVisible = false">取消</el-button>
        <el-button :loading="connTestLoading" @click="testTargetConnection">连接测试</el-button>
        <el-button type="primary" @click="saveTarget">保存</el-button>
      </span>
    </el-dialog>

    <el-dialog title="检查详情" :visible.sync="detailVisible" width="900px" custom-class="theme-dialog">
      <div v-if="detailLoading" class="detail-loading">加载中…</div>
      <div v-else-if="detailRows.length">
        <p class="detail-meta">{{ detailMeta }}</p>
        <el-table :data="detailRows" style="width: 100%" class="myTable" max-height="420">
          <template v-if="detailMode === 'baseline'">
            <el-table-column prop="categoryName" label="分类" width="120" />
            <el-table-column prop="ruleName" label="检查项" :show-overflow-tooltip="true" />
            <el-table-column prop="resultName" label="结果" width="90" />
            <el-table-column prop="riskName" label="风险" width="90" />
            <el-table-column prop="expectedValue" label="期望值" :show-overflow-tooltip="true" />
            <el-table-column prop="actualValue" label="实际值" :show-overflow-tooltip="true" />
            <el-table-column prop="fixSuggestion" label="修复建议" :show-overflow-tooltip="true" />
          </template>
          <template v-else-if="detailMode === 'malware'">
            <el-table-column prop="checkTypeName" label="检测类型" width="140" />
            <el-table-column prop="riskName" label="风险" width="90" />
            <el-table-column prop="matchRule" label="匹配规则" width="160" :show-overflow-tooltip="true" />
            <el-table-column prop="filePath" label="文件路径" :show-overflow-tooltip="true" />
            <el-table-column prop="processInfo" label="进程信息" :show-overflow-tooltip="true" />
            <el-table-column prop="description" label="描述" :show-overflow-tooltip="true" />
            <el-table-column prop="fixSuggestion" label="修复建议" :show-overflow-tooltip="true" />
          </template>
        </el-table>
      </div>
      <div v-else class="detail-empty">暂无明细</div>
      <span slot="footer">
        <el-button @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'

export default {
  name: 'HostSecurityHub',
  data() {
    return {
      taskTypeFilter: 'all',
      searchKeyword: '',
      searchApplied: '',
      tableLoading: false,
      tableRows: [],
      totalpage: 0,
      formData: { page: 1 },
      pageSize: 10,
      currentpage: 1,
      createVisible: false,
      submitLoading: false,
      osTypes: [],
      taskForm: {
        taskKind: 'baseline'
      },
      createRules: {
        taskKind: [{ required: true, message: '请选择任务类型', trigger: 'change' }]
      },
      targets: [],
      targetRules: {
        host: [{ required: true, message: '请输入目标主机', trigger: 'blur' }],
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        osType: [{ required: true, message: '请选择操作系统类型', trigger: 'change' }]
      },
      targetDialogVisible: false,
      editTargetIndex: -1,
      connTestLoading: false,
      connTestResult: null,
      editTargetForm: {
        host: '',
        transport: 'ssh',
        port: 22,
        winrmUseHttps: false,
        username: '',
        password: '',
        key: '',
        osType: 1
      },
      detailVisible: false,
      detailLoading: false,
      detailRows: [],
      detailMeta: '',
      detailMode: 'baseline',
      progressTaskId: 0,
      progressPollType: '',
      progressTimer: null,
      listRefreshTimer: null,
      listLoadPromise: null,
      lastListLoadAt: 0,
      listAutoRefreshMs: 6000,
      selectedRows: [],
      batchDeleteLoading: false,
      stopTaskLoading: false
    }
  },
  computed: {
    editTargetPortLabel() {
      return this.editTargetForm.transport === 'winrm' ? 'WinRM 端口' : 'SSH 端口'
    },
    displayRows() {
      if (this.taskTypeFilter === 'all') {
        return this.tableRows
      }
      const kw = (this.searchApplied || '').trim().toLowerCase()
      if (!kw) return this.tableRows
      return this.tableRows.filter(row => {
        const ip = (row.targetIp || '').toLowerCase()
        const kind = (row.kindLabel || '').toLowerCase()
        const taskId = String(row.taskId || '')
        return ip.includes(kw) || kind.includes(kw) || taskId.includes(kw)
      })
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/hostsec/tasks'
    this.loadEnums()
    this.loadData({ force: true })
    this.startListAutoRefresh()
  },
  beforeDestroy() {
    this.stopProgressPoll()
    this.stopListAutoRefresh()
  },
  methods: {
    startListAutoRefresh() {
      this.stopListAutoRefresh()
      this.listRefreshTimer = setInterval(() => {
        if (!this.shouldAutoRefreshList()) return
        this.loadData({ silent: true, minIntervalMs: this.listAutoRefreshMs })
      }, this.listAutoRefreshMs)
    },
    stopListAutoRefresh() {
      if (this.listRefreshTimer) {
        clearInterval(this.listRefreshTimer)
        this.listRefreshTimer = null
      }
    },
    shouldAutoRefreshList() {
      if (this._isDestroyed || this._isBeingDestroyed) return false
      if (this.$route && this.$route.path !== '/hostsec/tasks') return false
      if (this.createVisible || this.targetDialogVisible || this.detailVisible) return false
      return !this.progressTaskId
    },
    createTarget() {
      return {
        host: '',
        transport: 'ssh',
        port: 22,
        winrmUseHttps: false,
        username: '',
        password: '',
        key: '',
        osType: this.osTypes.length ? this.osTypes[0].value : 1
      }
    },
    addTarget() {
      this.targets.push(this.createTarget())
    },
    openAddTargetDialog() {
      this.editTargetIndex = -1
      this.connTestResult = null
      this.editTargetForm = {
        host: '',
        transport: 'ssh',
        port: 22,
        winrmUseHttps: false,
        username: '',
        password: '',
        key: '',
        osType: this.osTypes.length ? this.osTypes[0].value : 1
      }
      this.targetDialogVisible = true
    },
    onTargetDialogClosed() {
      this.connTestResult = null
      this.connTestLoading = false
    },
    clearConnTestResult() {
      this.connTestResult = null
    },
    buildTargetConnPayload(form) {
      return {
        host: (form.host || '').trim(),
        port: form.port,
        username: (form.username || '').trim(),
        password: form.password || '',
        key: form.transport === 'winrm' ? '' : (form.key || ''),
        osType: form.osType,
        transport: form.transport === 'winrm' ? 2 : 1,
        winrmUseHttps: form.transport === 'winrm' ? form.winrmUseHttps : false
      }
    },
    validateTargetForm(showMessage = true) {
      const form = this.editTargetForm
      const warn = (msg) => {
        if (showMessage) this.$message({ message: msg, type: 'warning' })
        return false
      }
      if (!(form.host || '').trim()) return warn('请输入目标主机')
      if (!(form.username || '').trim()) return warn('请输入用户名')
      if (form.transport === 'winrm') {
        if (!form.password) return warn('WinRM 需要填写密码')
      } else if (!form.password && !form.key) {
        return warn('请填写 SSH 密码或私钥')
      }
      return true
    },
    async testTargetConnection() {
      if (!this.validateTargetForm()) return
      this.connTestLoading = true
      this.connTestResult = null
      try {
        const res = await security.testHostConn(this.buildTargetConnPayload(this.editTargetForm))
        if (res.code === 200 && res.data) {
          this.connTestResult = {
            ok: !!res.data.ok,
            message: res.data.message || (res.data.ok ? '连接成功' : '连接失败'),
            detail: res.data.detail || ''
          }
          if (res.data.ok) {
            this.$message({ message: '连接测试通过', type: 'success' })
          }
        } else {
          this.connTestResult = { ok: false, message: res.msg || '连接测试失败', detail: '' }
          this.$message({ message: res.msg || '连接测试失败', type: 'error' })
        }
      } catch (e) {
        this.connTestResult = { ok: false, message: e.message || '连接测试失败', detail: '' }
        this.$message({ message: '连接测试失败: ' + (e.message || ''), type: 'error' })
      } finally {
        this.connTestLoading = false
      }
    },
    onEditTargetTransportChange() {
      this.clearConnTestResult()
      if (this.editTargetForm.transport === 'winrm') {
        this.editTargetForm.port = this.editTargetForm.winrmUseHttps ? 5986 : 5985
        this.editTargetForm.key = ''
      } else {
        this.editTargetForm.port = 22
      }
    },
    onEditWinrmHttpsChange() {
      this.clearConnTestResult()
      if (this.editTargetForm.transport === 'winrm') {
        this.editTargetForm.port = this.editTargetForm.winrmUseHttps ? 5986 : 5985
      }
    },
    onEditTargetOsTypeChange() {
      this.clearConnTestResult()
    },
    editTarget(idx) {
      const t = this.targets[idx]
      this.editTargetIndex = idx
      this.connTestResult = null
      this.editTargetForm = {
        host: t.host,
        transport: t.transport || 'ssh',
        port: t.port,
        winrmUseHttps: t.winrmUseHttps || false,
        username: t.username,
        password: t.password,
        key: t.key,
        osType: t.osType
      }
      this.targetDialogVisible = true
    },
    saveTarget() {
      if (!this.validateTargetForm()) return
      if (this.editTargetIndex >= 0) {
        this.targets[this.editTargetIndex] = { ...this.editTargetForm }
      } else {
        this.targets.push({ ...this.editTargetForm })
      }
      this.targetDialogVisible = false
      this.$message({ message: this.editTargetIndex >= 0 ? '修改成功' : '添加成功', type: 'success' })
    },
    getOsTypeName(osType) {
      const opt = this.osTypes.find(o => o.value === osType)
      return opt ? opt.label : '未知'
    },
    removeTarget(idx) {
      this.targets.splice(idx, 1)
    },
    onTaskKindChange() {
    },
    async loadEnums() {
      const res = await security.getBaselineEnums()
      if (res.code === 200 && res.data) {
        if (res.data.osTypes) this.osTypes = res.data.osTypes
      }
      if (!this.osTypes.length) {
        this.osTypes = [
          { value: 1, label: 'Linux/Unix' },
          { value: 2, label: 'Windows' },
          { value: 3, label: '国产操作系统' },
          { value: 4, label: '嵌入式 OS' }
        ]
      }
    },
    onFilterChange() {
      const ownerTab = this.progressOwnerTab(this.progressPollType)
      if (ownerTab && ownerTab !== this.taskTypeFilter && this.taskTypeFilter !== 'all') {
        this.stopProgressPoll()
        this.progressTaskId = 0
        this.progressPollType = ''
      }
      this.formData.page = 1
      this.currentpage = 1
      this.selectedRows = []
      if (this.$refs.taskTable) {
        this.$refs.taskTable.clearSelection()
      }
      this.loadData({ force: true })
    },
    applySearch() {
      this.searchApplied = (this.searchKeyword || '').trim()
      this.formData.page = 1
      this.currentpage = 1
      this.loadData({ force: true })
    },
    resetSearch() {
      this.searchKeyword = ''
      this.searchApplied = ''
      this.taskTypeFilter = 'all'
      this.formData.page = 1
      this.currentpage = 1
      this.selectedRows = []
      if (this.$refs.taskTable) {
        this.$refs.taskTable.clearSelection()
      }
      this.loadData({ force: true })
    },
    progressOwnerTab(pollType) {
      const map = { baseline: 'baseline', cve: 'vuln', yara: 'malware' }
      return map[pollType] || ''
    },
    tabSource(taskTypeFilter) {
      if (taskTypeFilter === 'baseline' || taskTypeFilter === 'vuln' || taskTypeFilter === 'malware') {
        return taskTypeFilter
      }
      return ''
    },
    filterRunningRowsForTab(rows, taskTypeFilter) {
      const source = this.tabSource(taskTypeFilter)
      if (!source) {
        return (rows || []).filter(r => r.isRunning)
      }
      return (rows || []).filter(r => r.isRunning && r.source === source)
    },
    mergeRowsWithRunningRows(rows, runningRows) {
      const merged = []
      const seen = new Set()
      ;(runningRows || []).forEach(row => {
        const key = this.taskRowKey(row)
        if (seen.has(key)) return
        seen.add(key)
        merged.push(row)
      })
      ;(rows || []).forEach(row => {
        const key = this.taskRowKey(row)
        if (seen.has(key)) return
        seen.add(key)
        merged.push(row)
      })
      return merged
    },
    baselineCheckedCount(r) {
      return (r.passCount || 0) + (r.failCount || 0) + (r.errorCount || 0)
    },
    baselineProgressText(r, done, total) {
      if (typeof done === 'number' && typeof total === 'number' && total > 0) {
        return `${done}/${total}`
      }
      const checked = this.baselineCheckedCount(r)
      const rules = r.totalRules || 0
      if (r.isRunning && rules > 0) {
        return `${checked}/${rules}`
      }
      if (r.isRunning) {
        return '准备中'
      }
      return '—'
    },
    mapBaselineListRow(r, extra = {}) {
      const isRunning = !!r.isRunning
      const row = {
        source: 'baseline',
        kindLabel: r.scanSceneName || '安全配置核查',
        taskId: r.taskId,
        targetIp: r.targetIp,
        scanScene: r.scanScene || 1,
        osTypeName: r.osTypeName,
        totalRules: r.totalRules,
        passCount: r.passCount,
        failCount: r.failCount,
        errorCount: r.errorCount,
        checkTime: r.checkTime,
        isRunning,
        runStatus: isRunning ? 'running' : '',
        runStatusLabel: isRunning ? '执行中' : '',
        progressText: this.baselineProgressText(r),
        ...extra
      }
      return row
    },
    maybeResumeBaselinePoll(rows) {
      if (this.progressTaskId) return
      const running = (rows || []).find(r => r.source === 'baseline' && r.isRunning)
      if (!running) return
      this.progressTaskId = running.taskId
      this.progressPollType = 'baseline'
      this.pollProgressInline()
    },
    async loadData(options = {}) {
      const { silent = false, force = false, minIntervalMs = 2500 } = options
      const now = Date.now()
      if (!force) {
        if (this.listLoadPromise) return this.listLoadPromise
        if (this.lastListLoadAt && now - this.lastListLoadAt < minIntervalMs) {
          return Promise.resolve()
        }
      }
      const runningRows = this.filterRunningRowsForTab(this.tableRows, this.taskTypeFilter)
      if (!silent) {
        this.tableLoading = true
      }
      const loadPromise = (async () => {
        if (this.taskTypeFilter === 'all') {
          await this.loadUnifiedList()
        } else if (this.taskTypeFilter === 'malware') {
          const res = await security.getYaraTaskList({
            page: this.formData.page,
            size: this.pageSize
          })
          if (res.code === 200 && res.data) {
            this.totalpage = res.data.total || 0
            const list = res.data.list || []
            this.tableRows = list.map((r) => ({
              source: 'malware',
              kindLabel: '恶意代码检测',
              taskId: r.taskId,
              targetIp: r.targetIp,
              scanScene: 0,
              isRunning: r.scanStatus === 0,
              runStatus: r.scanStatus === 0 ? 'running' : '',
              runStatusLabel: r.scanStatus === 0 ? (r.scanStatusName || '执行中') : '',
              progressText: r.scanStatus === 0 ? '进行中' : '—',
              osTypeName: r.osTypeName,
              totalFindings: r.totalFindings,
              critical: r.critical,
              high: r.high,
              worstRiskName: r.worstRiskName,
              scanStatus: r.scanStatus,
              scanStatusName: r.scanStatusName,
              errorMessage: r.errorMessage,
              checkTime: r.checkTime
            }))
          } else {
            this.tableRows = []
            this.$message({ message: res.msg || '加载失败', type: 'error' })
          }
        } else if (this.taskTypeFilter === 'vuln') {
          const res = await security.getHostVulnTaskList({
            page: this.formData.page,
            size: this.pageSize
          })
          if (res.code === 200 && res.data) {
            this.totalpage = res.data.total || 0
            this.tableRows = (res.data.list || []).map((r) => ({
              source: 'vuln',
              kindLabel: '主机漏洞检测',
              taskId: r.taskId,
              targetIp: r.targetIp,
              scanScene: 2,
              osTypeName: r.osTypeName,
              packages: r.packages,
              matchedVulns: r.matchedVulns,
              critical: r.critical,
              high: r.high,
              medium: r.medium,
              low: r.low,
              worstRiskName: r.worstRiskName,
              scanStatus: r.scanStatus,
              scanStatusName: r.scanStatusName,
              errorMessage: r.errorMessage,
              checkTime: r.checkTime,
              isRunning: r.scanStatus === 0,
              runStatus: r.scanStatus === 0 ? 'running' : '',
              runStatusLabel: r.scanStatus === 0 ? (r.scanStatusName || '执行中') : '',
              progressText: r.scanStatus === 0 ? '进行中' : '—'
            }))
          } else {
            this.tableRows = []
            this.$message({ message: res.msg || '加载失败', type: 'error' })
          }
        } else {
          const res = await security.getBaselineTaskList({
            page: this.formData.page,
            size: this.pageSize,
            scanScene: 1
          })
          if (res.code === 200 && res.data) {
            this.totalpage = res.data.total || 0
            this.tableRows = (res.data.list || []).map((r) => this.mapBaselineListRow(r))
            this.maybeResumeBaselinePoll(this.tableRows)
          } else {
            this.tableRows = []
            this.$message({ message: res.msg || '加载失败', type: 'error' })
          }
        }
        if (runningRows.length) {
          this.tableRows = this.mergeRowsWithRunningRows(this.tableRows, runningRows)
        }
      })()
      this.listLoadPromise = loadPromise
      try {
        await loadPromise
      } finally {
        this.lastListLoadAt = Date.now()
        this.listLoadPromise = null
        if (!silent) {
          this.tableLoading = false
        }
      }
    },
    async loadUnifiedList() {
      const page = this.formData.page
      const size = this.pageSize
      const fetchSize = Math.max(page * size, 50)
      const kw = (this.searchApplied || '').trim().toLowerCase()
      const [bRes, vRes, mRes] = await Promise.all([
        security.getBaselineTaskList({ page: 1, size: fetchSize, scanScene: 1 }),
        security.getHostVulnTaskList({ page: 1, size: fetchSize }),
        security.getYaraTaskList({ page: 1, size: fetchSize })
      ])
      const merged = []
      let total = 0
      if (bRes.code === 200 && bRes.data && bRes.data.list) {
        total += bRes.data.total || 0
        for (const r of bRes.data.list) {
          merged.push(this.mapBaselineListRow(r, {
            _ts: this.parseCheckTime(r.checkTime)
          }))
        }
      }
      if (vRes.code === 200 && vRes.data && vRes.data.list) {
        total += vRes.data.total || 0
        for (const r of vRes.data.list) {
          merged.push({
            source: 'vuln',
            kindLabel: '主机漏洞检测',
            taskId: r.taskId,
            targetIp: r.targetIp,
            scanScene: 2,
            osTypeName: r.osTypeName,
            packages: r.packages,
            matchedVulns: r.matchedVulns,
            critical: r.critical,
            high: r.high,
            medium: r.medium,
            low: r.low,
            worstRiskName: r.worstRiskName,
            scanStatus: r.scanStatus,
            scanStatusName: r.scanStatusName,
            errorMessage: r.errorMessage,
            checkTime: r.checkTime,
            isRunning: r.scanStatus === 0,
            runStatus: r.scanStatus === 0 ? 'running' : '',
            runStatusLabel: r.scanStatus === 0 ? (r.scanStatusName || '执行中') : '',
            progressText: r.scanStatus === 0 ? '进行中' : '—',
            _ts: this.parseCheckTime(r.checkTime)
          })
        }
      }
      if (mRes.code === 200 && mRes.data && mRes.data.list) {
        total += mRes.data.total || 0
        for (const r of mRes.data.list) {
          merged.push({
            source: 'malware',
            kindLabel: '恶意代码检测',
            taskId: r.taskId,
            targetIp: r.targetIp,
            scanScene: 0,
            osTypeName: r.osTypeName,
            totalFindings: r.totalFindings,
            critical: r.critical,
            high: r.high,
            worstRiskName: r.worstRiskName,
            scanStatus: r.scanStatus,
            scanStatusName: r.scanStatusName,
            errorMessage: r.errorMessage,
            checkTime: r.checkTime,
            isRunning: r.scanStatus === 0,
            runStatus: r.scanStatus === 0 ? 'running' : '',
            runStatusLabel: r.scanStatus === 0 ? (r.scanStatusName || '执行中') : '',
            progressText: r.scanStatus === 0 ? '进行中' : '—',
            _ts: this.parseCheckTime(r.checkTime)
          })
        }
      }
      merged.sort((a, b) => (b._ts || 0) - (a._ts || 0))
      let filtered = merged
      if (kw) {
        filtered = merged.filter(row => {
          const ip = (row.targetIp || '').toLowerCase()
          const kind = (row.kindLabel || '').toLowerCase()
          const taskId = String(row.taskId || '')
          return ip.includes(kw) || kind.includes(kw) || taskId.includes(kw)
        })
        this.totalpage = filtered.length
      } else {
        this.totalpage = total
      }
      const offset = (page - 1) * size
      this.tableRows = filtered.slice(offset, offset + size)
      this.maybeResumeBaselinePoll(this.tableRows)
    },
    parseCheckTime(checkTime) {
      if (!checkTime || checkTime === '—') return 0
      return new Date(String(checkTime).replace(/-/g, '/')).getTime() || 0
    },
    async loadMerged() {
      await this.loadUnifiedList()
    },
    openCreateDialog() {
      this.createVisible = true
      this.taskForm = {
        taskKind: 'baseline'
      }
      this.targets = []
      this.$nextTick(() => {
        if (this.$refs.taskForm) this.$refs.taskForm.clearValidate()
      })
    },
    async submitCreate() {
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
        if (this.taskForm.taskKind !== 'malware') {
          if (t.transport === 'winrm') {
            if (!t.password) {
              this.$message({ message: `目标 #${i + 1} WinRM 需要填写密码`, type: 'warning' })
              return
            }
          } else if (!t.password && !t.key) {
            this.$message({ message: `目标 #${i + 1} 请填写 SSH 密码或私钥`, type: 'warning' })
            return
          }
        } else if (!t.password && !t.key) {
          this.$message({ message: `目标 #${i + 1} 请填写密码或私钥`, type: 'warning' })
          return
        }
      }
      this.submitLoading = true
      const targetSnapshot = this.targets.map(t => ({ ...t }))
      const taskKind = this.taskForm.taskKind
      try {
        if (taskKind === 'malware') {
          const payload = {
            targets: this.targets.map(t => ({
              host: t.host,
              port: t.port,
              username: t.username,
              password: t.password,
              key: t.transport === 'winrm' ? '' : t.key,
              osType: t.osType,
              transport: t.transport === 'winrm' ? 2 : 1,
              winrmUseHttps: t.transport === 'winrm' ? t.winrmUseHttps : false
            }))
          }
          this.createVisible = false
          const res = await security.runYaraBatchScan(payload)
          if (res.code === 200) {
            this.startInlineProgress(res.data.taskId, taskKind, targetSnapshot)
          } else {
            this.$message({ message: res.msg || '创建任务失败', type: 'error' })
          }
        } else if (taskKind === 'vuln') {
          const payload = {
            targets: this.targets.map(t => ({
              host: t.host,
              port: t.port,
              username: t.username,
              password: t.password,
              key: t.transport === 'winrm' ? '' : t.key,
              osType: t.osType,
              transport: t.transport === 'winrm' ? 2 : 1,
              winrmUseHttps: t.transport === 'winrm' ? t.winrmUseHttps : false
            }))
          }
          this.createVisible = false
          const res = await security.runCveBatchScan(payload)
          if (res.code === 200) {
            this.startInlineProgress(res.data.taskId, taskKind, targetSnapshot)
          } else {
            this.$message({ message: res.msg || '创建任务失败', type: 'error' })
          }
        } else {
          const scanScene = 1
          const payload = {
            targets: this.targets.map(t => ({
              host: t.host,
              port: t.port,
              username: t.username,
              password: t.password,
              key: t.transport === 'winrm' ? '' : t.key,
              osType: t.osType,
              transport: t.transport === 'winrm' ? 2 : 1,
              winrmUseHttps: t.transport === 'winrm' ? t.winrmUseHttps : false,
              scanScene
            }))
          }
          this.createVisible = false
          const res = await security.runBaselineBatchCheck(payload)
          if (res.code === 200) {
            this.startInlineProgress(res.data.taskId, taskKind, targetSnapshot)
          } else {
            this.$message({ message: res.msg || '创建任务失败', type: 'error' })
          }
        }
      } finally {
        this.submitLoading = false
      }
    },
    buildRunningRows(taskId, kind, targets) {
      const kindMap = { baseline: '安全配置核查', vuln: '主机漏洞检测', malware: '恶意代码检测' }
      const sourceMap = { baseline: 'baseline', vuln: 'vuln', malware: 'malware' }
      const scanScene = kind === 'vuln' ? 2 : kind === 'baseline' ? 1 : 0
      const total = targets.length || 1
      return targets.map(t => ({
        isRunning: true,
        taskId,
        source: sourceMap[kind],
        kindLabel: kindMap[kind],
        targetIp: t.host,
        osTypeName: this.getOsTypeName(t.osType),
        scanScene,
        runStatus: 'running',
        runStatusLabel: '执行中',
        progressText: kind === 'baseline' ? '准备中' : `0/${total}`,
        progressDone: 0,
        progressTotal: kind === 'baseline' ? 0 : total,
        checkTime: '—',
        runMessage: ''
      }))
    },
    startInlineProgress(taskId, kind, targets) {
      this.stopProgressPoll()
      this.progressTaskId = taskId
      const pollMap = { baseline: 'baseline', vuln: 'cve', malware: 'yara' }
      this.progressPollType = pollMap[kind] || 'baseline'
      this.taskTypeFilter = 'all'
      this.formData.page = 1
      this.currentpage = 1
      this.selectedRows = []
      if (this.$refs.taskTable) {
        this.$refs.taskTable.clearSelection()
      }
      const runningRows = this.buildRunningRows(taskId, kind, targets)
      this.loadData({ silent: true, force: true }).then(() => {
        this.tableRows = this.mergeRowsWithRunningRows(this.tableRows, runningRows)
        this.totalpage = Math.max(this.totalpage, this.tableRows.length)
        this.pollProgressInline()
      })
    },
    updateRunningRowsFromProgress(apiTargets, done, total) {
      if (!this.progressTaskId) return
      this.tableRows = this.tableRows.map(row => {
        if (!row.isRunning || row.taskId !== this.progressTaskId) return row
        const pt = (apiTargets || []).find(t => (t.host || t.targetIp) === row.targetIp)
        let progressDone = done
        let progressTotal = total
        if (this.progressPollType === 'baseline' && pt && pt.totalItems > 0) {
          progressDone = (pt.items || []).length
          progressTotal = pt.totalItems
        }
        const updates = {
          progressDone,
          progressTotal,
          progressText: this.baselineProgressText(row, progressDone, progressTotal),
          isRunning: true,
          runStatus: 'running',
          runStatusLabel: '执行中'
        }
        if (pt) {
          const status = pt.status || (pt.error ? 'failed' : 'completed')
          if (status === 'completed' || status === 'failed') {
            updates.runStatus = status
            updates.runStatusLabel = this.runStatusLabel(status)
          }
          updates.runMessage = pt.message || pt.error || ''
          if (this.progressPollType === 'baseline' && status === 'running') {
            updates.passCount = (pt.items || []).filter(i => i.checkResult === 1).length
            updates.failCount = (pt.items || []).filter(i => i.checkResult === 2).length
            updates.errorCount = (pt.items || []).filter(i => i.checkResult === 3).length
            updates.totalRules = pt.totalItems || row.totalRules
          }
        } else if (done < total) {
          updates.runStatus = 'running'
          updates.runStatusLabel = '执行中'
        }
        return { ...row, ...updates }
      })
    },
    async pollProgressInline() {
      if (!this.progressTaskId) return
      try {
        let completed = false
        if (this.progressPollType === 'baseline') {
          const res = await security.getBaselineBatchProgress({ taskId: this.progressTaskId })
          if (res.code === 200 && res.data) {
            const d = res.data
            this.updateRunningRowsFromProgress(d.targets, d.completedTargets || 0, d.totalTargets || 1)
            completed = d.status === 'completed'
          }
        } else if (this.progressPollType === 'cve') {
          const res = await security.getCveBatchProgress({ taskId: this.progressTaskId })
          if (res.code === 200 && res.data) {
            const d = res.data
            const targets = (d.results || []).map(r => ({
              host: r.targetIp,
              status: r.error ? 'failed' : 'completed',
              error: r.error
            }))
            this.updateRunningRowsFromProgress(targets, d.progress || 0, d.total || 1)
            completed = d.status === 'completed'
          }
        } else if (this.progressPollType === 'yara') {
          const res = await security.getYaraBatchProgress({ taskId: this.progressTaskId })
          if (res.code === 200 && res.data) {
            const d = res.data
            const targets = (d.results || []).map(r => ({
              host: r.targetIp,
              status: r.error ? 'failed' : 'completed',
              error: r.error
            }))
            this.updateRunningRowsFromProgress(targets, d.progress || 0, d.total || 1)
            completed = d.status === 'completed'
          }
        }
        if (completed) {
          this.finishProgressPoll()
          return
        }
      } catch (e) {
        console.error('pollProgressInline error:', e)
      }
      this.progressTimer = setTimeout(() => this.pollProgressInline(), this.progressPollType === 'baseline' ? 2500 : 3000)
    },
    finishProgressPoll() {
      this.stopProgressPoll()
      this.progressTaskId = 0
      this.progressPollType = ''
      this.tableRows = this.tableRows.filter(r => !r.isRunning)
      this.formData.page = 1
      this.currentpage = 1
      this.loadData({ silent: true, force: true })
    },
    stopProgressPoll() {
      if (this.progressTimer) {
        clearTimeout(this.progressTimer)
        this.progressTimer = null
      }
    },
    runStatusTagType(status) {
      if (status === 'completed') return 'success'
      if (status === 'failed') return 'danger'
      if (status === 'running') return 'warning'
      return 'info'
    },
    runStatusLabel(status) {
      if (status === 'completed') return '已完成'
      if (status === 'failed') return '执行失败'
      if (status === 'running') return '执行中'
      return '等待中'
    },
    openDetailPage(row) {
      this.$router.push({
        path: '/hostsec/task-detail',
        query: {
          taskId: row.taskId,
          kindLabel: row.kindLabel || '安全配置核查',
          checkTime: row.checkTime || '',
          source: row.source || (row.scanScene === 2 ? 'vuln' : row.source === 'malware' ? 'malware' : 'baseline')
        }
      })
    },
    rowSummary(row) {
      if (row.source === 'vuln') {
        if (row.scanStatus === 2) return `扫描异常：${row.errorMessage || '未知错误'}`
        return `包 ${row.packages || 0} · 漏洞 ${row.matchedVulns || 0} · 严重 ${row.critical || 0} / 高危 ${row.high || 0}`
      }
      if (row.source === 'malware') {
        if (row.scanStatus === 2) return `扫描异常：${row.errorMessage || '未知错误'}`
        return `发现 ${row.totalFindings || 0} 项 · 严重 ${row.critical || 0} / 高危 ${row.high || 0}`
      }
      if (row.isRunning) {
        const checked = this.baselineCheckedCount(row)
        const total = row.progressTotal || row.totalRules || '—'
        return `执行中 · 已检查 ${checked}/${total} 项`
      }
      return `检查 ${row.totalRules || 0} · 通过 ${row.passCount || 0} / 不通过 ${row.failCount || 0}`
    },
    taskRowKey(row) {
      return `${row.source}-${row.taskId}-${row.targetIp || ''}-${row.scanScene || 0}`
    },
    onSelectionChange(rows) {
      this.selectedRows = rows || []
    },
    toDeleteItem(row) {
      return {
        source: row.source,
        taskId: row.taskId,
        targetIp: row.targetIp || '',
        scanScene: row.scanScene || (row.source === 'vuln' ? 2 : row.source === 'baseline' ? 1 : 0)
      }
    },
    async deleteTask(row) {
      try {
        await this.$confirm(
          `确定删除「${row.kindLabel}」任务记录吗？\n目标：${row.targetIp || '-'} · #${row.taskId}`,
          '删除确认',
          { type: 'warning' }
        )
      } catch {
        return
      }
      await this.doDeleteTasks([row])
    },
    async stopTask(row) {
      if (!row || !row.isRunning) return
      try {
        await this.$confirm(
          `确定结束「${row.kindLabel}」任务吗？\n目标：${row.targetIp || '-'} · #${row.taskId}`,
          '结束确认',
          { type: 'warning' }
        )
      } catch {
        return
      }
      this.stopTaskLoading = true
      try {
        const res = await security.stopHostSecTasks({
          items: [this.toDeleteItem(row)]
        })
        if (res.code === 200) {
          this.$message({ message: '任务已结束', type: 'success' })
          if (this.progressTaskId === row.taskId) {
            this.stopProgressPoll()
            this.progressTaskId = 0
            this.progressPollType = ''
          }
          await this.loadData({ force: true })
        } else {
          this.$message({ message: res.msg || '结束失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: e.message || '结束失败', type: 'error' })
      } finally {
        this.stopTaskLoading = false
      }
    },
    async batchDeleteTasks() {
      if (!this.selectedRows.length) return
      try {
        await this.$confirm(
          `确定删除选中的 ${this.selectedRows.length} 条任务记录吗？此操作不可恢复。`,
          '批量删除确认',
          { type: 'warning' }
        )
      } catch {
        return
      }
      await this.doDeleteTasks(this.selectedRows)
    },
    async doDeleteTasks(rows) {
      this.batchDeleteLoading = true
      try {
        const res = await security.deleteHostSecTasks({
          items: rows.map(r => this.toDeleteItem(r))
        })
        if (res.code === 200) {
          this.$message({ message: `已删除 ${(res.data && res.data.deleted) || rows.length} 条记录`, type: 'success' })
          const deletedKeys = new Set(rows.map(r => this.taskRowKey(r)))
          this.tableRows = this.tableRows.filter(r => !deletedKeys.has(this.taskRowKey(r)))
          this.selectedRows = []
          if (this.$refs.taskTable) {
            this.$refs.taskTable.clearSelection()
          }
          this.loadData({ force: true })
        } else {
          this.$message({ message: res.msg || '删除失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: e.message || '删除失败', type: 'error' })
      } finally {
        this.batchDeleteLoading = false
      }
    },
    async openDetail(row) {
      this.detailVisible = true
      this.detailLoading = true
      this.detailRows = []
      this.detailMeta = `${row.kindLabel} · 目标：${row.targetIp || ''}`
      try {
        if (row.source === 'malware') {
          this.detailMode = 'malware'
          const res = await security.getYaraResultList({ taskId: row.taskId })
          if (res.code === 200) {
            this.detailRows = res.data.list || []
          } else {
            this.$message({ message: res.msg || '加载失败', type: 'error' })
          }
        }
      } finally {
        this.detailLoading = false
      }
    },
    handlecurrentchange(t) {
      this.formData.page = t
      this.currentpage = t
      this.loadData({ force: true })
    },
    handleSizeChange(t) {
      this.formData.page = 1
      this.pageSize = t
      this.loadData({ force: true })
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';
@import './css/appsec-tokens.less';

.page-intro {
  color: #94a3b8;
  font-size: 13px;
  line-height: 1.55;
  margin: 0 0 12px;
  max-width: 900px;
}

.tab-panel {
  margin-top: 4px;
}

.list_box {
  width: 100%;
  max-width: 100%;
  min-width: 0;
  box-sizing: border-box;
}

.search-box {
  flex-wrap: wrap;
  gap: 12px;
  max-width: 100%;
}

.operationbutton {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.serach-condition {
  flex-wrap: wrap;
  flex-shrink: 0;
}

.task-search-input {
  width: 220px;
}

.task-type-select {
  width: 168px;
}

.table-scroll-wrap {
  width: 100%;
  max-width: 100%;
  min-width: 0;
  overflow-x: auto;
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

.targets-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-top: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.targets-title {
  font-size: 14px;
  color: #94a3b8;
}

.target-card {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 12px;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.target-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px dashed rgba(255, 255, 255, 0.06);
}

.target-card-index {
  font-size: 13px;
  color: #06b6d4;
  font-weight: 500;
}

.target-form {
  margin: 0;
}

.inline-progress {
  display: inline-block;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
  color: #e2e8f0;
  font-size: 13px;

  &.done {
    color: rgba(148, 163, 184, 0.5);
  }
}

.run-hint {
  font-size: 12px;
  color: rgba(148, 163, 184, 0.75);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
  max-width: 110px;
}

.targets-empty {
  padding: 24px;
  text-align: center;
  color: #94a3b8;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 8px;
  border: 1px dashed rgba(255, 255, 255, 0.08);
}

.targets-table {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  overflow: hidden;
}

.targets-table-header {
  display: flex;
  background: rgba(0, 0, 0, 0.2);
  padding: 10px 12px;
  font-size: 12px;
  color: #94a3b8;
  font-weight: 500;
}

.targets-table-row {
  display: flex;
  padding: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  align-items: center;
  transition: background-color 0.2s;

  &:hover {
    background: rgba(255, 255, 255, 0.03);
  }

  &:last-child {
    border-bottom: none;
  }
}

.th-index, .td-index {
  width: 40px;
  text-align: center;
}

.th-host, .td-host {
  flex: 1;
  min-width: 120px;
  color: #e2e8f0;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.th-protocol, .td-protocol {
  width: 80px;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

.th-port, .td-port {
  width: 70px;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

.th-username, .td-username {
  width: 110px;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

.th-os, .td-os {
  width: 100px;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

.th-actions, .td-actions {
  width: 100px;
  text-align: center;
  display: inline-flex;
  gap: 2px;
}

.td-actions .el-button {
  padding: 0 6px;
  margin: 0;
}

::v-deep .el-table--enable-row-hover .el-table__body tr:hover > td {
  background-color: rgba(0, 212, 170, 0.08) !important;
}

::v-deep .host-target-dialog {
  .el-dialog__body {
    padding: 8px 24px 4px;
  }
}

.host-target-form {
  .form-section {
    margin-bottom: 18px;
    padding-bottom: 6px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);

    &:last-of-type {
      border-bottom: none;
      margin-bottom: 8px;
    }
  }

  .section-title {
    font-size: 13px;
    font-weight: 600;
    color: #00d4aa;
    margin-bottom: 12px;
  }

  .form-tip {
    font-size: 12px;
    color: #64748b;
    margin: -4px 0 12px;
    line-height: 1.5;
  }

  .switch-hint {
    margin-left: 8px;
    font-size: 12px;
    color: #64748b;
  }

  ::v-deep .el-form-item__label {
    color: #cbd5e1;
    font-weight: 600;
    line-height: 1.4;
    padding-bottom: 6px;
  }

  ::v-deep .el-form-item {
    margin-bottom: 14px;
  }

  ::v-deep .el-textarea__inner {
    font-family: Consolas, 'Courier New', monospace;
    font-size: 12px;
  }
}

.conn-test-result {
  padding: 10px 12px;
  border-radius: @appsec-radius-sm;
  font-size: 13px;
  line-height: 1.5;

  .result-line {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  &.ok {
    background: rgba(16, 185, 129, 0.12);
    border: 1px solid rgba(16, 185, 129, 0.35);
    color: #6ee7b7;
  }

  &.fail {
    background: rgba(239, 68, 68, 0.12);
    border: 1px solid rgba(239, 68, 68, 0.35);
    color: #fca5a5;
  }

  .result-detail {
    padding-left: 20px;
    font-size: 12px;
    opacity: 0.9;
    word-break: break-all;
    font-family: Consolas, 'Courier New', monospace;
  }
}

.host-target-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
