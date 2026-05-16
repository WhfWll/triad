<template>
  <div class="security-container">
    <div class="main-title">主机安全检查 · 任务管理</div>
    <p class="page-intro">
      同一套远程连接（SSH / WinRM）能力，新建任务时区分：<strong>安全配置核查</strong>、<strong>主机漏洞检测</strong>（与核查共用规则库与引擎，仅任务归类不同）、<strong>恶意代码检测</strong>。
    </p>

    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="default" size="small" @click="$router.push('/hostsec/rules')">检测规则</el-button>
          <el-button type="primary" size="small" @click="openCreateDialog">新建主机检查任务</el-button>
        </div>
      </div>

      <el-tabs v-model="listTab" class="hub-tabs" @tab-click="onListTabChange">
        <el-tab-pane label="全部（最近）" name="all" />
        <el-tab-pane label="安全配置核查" name="baseline" />
        <el-tab-pane label="主机漏洞检测" name="vuln" />
        <el-tab-pane label="恶意代码检测" name="malware" />
      </el-tabs>

      <el-table v-loading="tableLoading" :data="tableRows" style="width: 100%" class="myTable">
        <el-table-column prop="kindLabel" label="任务类型" width="130" />
        <el-table-column prop="targetIp" label="目标主机" :show-overflow-tooltip="true" />
        <el-table-column v-if="listTab !== 'malware'" prop="osTypeName" label="操作系统" width="120" />
        <el-table-column v-if="listTab === 'malware'" prop="totalFindings" label="发现项数" width="100" />
        <el-table-column v-if="listTab !== 'malware'" prop="totalRules" label="检查项数" width="100" />
        <el-table-column v-if="listTab !== 'malware'" prop="passCount" label="通过" width="72" />
        <el-table-column v-if="listTab !== 'malware'" prop="failCount" label="不通过" width="82" />
        <el-table-column v-if="listTab !== 'malware'" prop="errorCount" label="异常" width="72" />
        <el-table-column v-if="listTab === 'malware'" prop="worstRiskName" label="最高风险" width="100" />
        <el-table-column prop="checkTime" label="时间" width="170" />
        <el-table-column label="操作" width="90">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="openDetail(scope.row)">详情</el-link>
          </template>
        </el-table-column>
      </el-table>

      <p v-if="listTab === 'all'" class="merge-hint">以上为各类型最近一批记录的合并视图；精确分页请在上方子类页签中查看。</p>

      <el-pagination
        v-if="listTab !== 'all'"
        :page-size="pageSize"
        background
        layout="total, prev, pager, next, sizes, jumper"
        :total="totalpage"
        :current-page="currentpage"
        @current-change="handlecurrentchange"
        @size-change="handleSizeChange"
      />
    </div>

    <el-dialog title="新建主机检查任务" :visible.sync="createVisible" width="700px" custom-class="theme-dialog" :close-on-click-modal="false">
      <el-form :model="taskForm" :rules="createRules" ref="taskForm" label-width="118px">
        <el-form-item label="任务类型" prop="taskKind">
          <el-radio-group v-model="taskForm.taskKind" @change="onTaskKindChange">
            <el-radio label="baseline">安全配置核查</el-radio>
            <el-radio label="vuln">主机漏洞检测</el-radio>
            <el-radio label="malware">恶意代码检测</el-radio>
          </el-radio-group>
          <p class="field-hint">
            <template v-if="taskForm.taskKind === 'vuln'">与「安全配置核查」共用规则库与检测引擎，便于验收「漏洞规则」类指标；结果在列表中归类为漏洞检测。</template>
            <template v-else-if="taskForm.taskKind === 'malware'">登录后执行恶意行为/特征类脚本检测；当前实现以 SSH 自动连接为主（与操作系统选项一致）。</template>
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

    <el-dialog :title="editTargetIndex >= 0 ? '编辑目标' : '添加目标'" :visible.sync="targetDialogVisible" width="500px" custom-class="theme-dialog">
      <el-form :model="editTargetForm" :rules="targetRules" ref="editTargetForm" label-width="90px">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="目标主机" prop="host">
              <el-input v-model="editTargetForm.host" placeholder="IP 或域名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="连接协议" prop="transport">
              <el-select v-model="editTargetForm.transport" placeholder="请选择" style="width: 100%" @change="onEditTargetTransportChange">
                <el-option :value="'ssh'" label="SSH" />
                <el-option :value="'winrm'" label="WinRM" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item :label="editTargetPortLabel" prop="port">
              <el-input-number v-model="editTargetForm.port" :min="1" :max="65535" :controls="false" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="HTTPS" v-if="editTargetForm.transport === 'winrm'">
              <el-switch v-model="editTargetForm.winrmUseHttps" @change="onEditWinrmHttpsChange" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item :label="editTargetUsernameLabel" prop="username">
              <el-input v-model="editTargetForm.username" placeholder="登录用户名" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="editTargetPasswordLabel" prop="password">
              <el-input v-model="editTargetForm.password" type="password" :placeholder="editTargetForm.transport === 'winrm' ? 'WinRM 密码（必填）' : '密码（与私钥二选一）'" show-password autocomplete="new-password" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="操作系统" prop="osType">
              <el-select v-model="editTargetForm.osType" placeholder="请选择" style="width: 100%" @change="onEditTargetOsTypeChange">
                <el-option v-for="opt in osTypes" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item v-if="editTargetForm.transport !== 'winrm'" label="SSH 私钥" prop="key">
              <el-input v-model="editTargetForm.key" type="textarea" :rows="2" placeholder="可选：粘贴私钥 PEM" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <span slot="footer">
        <el-button @click="targetDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTarget">保存</el-button>
      </span>
    </el-dialog>

    <el-dialog title="任务执行中" :visible.sync="progressVisible" width="500px" custom-class="theme-dialog" :close-on-click-modal="false" :show-close="false">
      <div class="progress-body">
        <div class="progress-summary">
          <el-progress :percentage="progressPercent" :status="progressStatus" :stroke-width="16" />
          <p class="progress-text">{{ progressText }}</p>
        </div>
        <div v-for="(t, idx) in progressTargets" :key="idx" class="progress-target-row">
          <span class="progress-target-host">{{ t.host }}</span>
          <el-tag :type="progressTagType(t.status)" size="mini">{{ progressTagLabel(t.status) }}</el-tag>
        </div>
      </div>
      <span slot="footer">
        <el-button v-if="progressDone" type="primary" @click="onProgressDone">查看结果</el-button>
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
          <template v-else>
            <el-table-column prop="checkTypeName" label="检测类型" width="140" />
            <el-table-column prop="riskName" label="风险" width="90" />
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
      listTab: 'all',
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
      targets: [],
      targetRules: {
        host: [{ required: true, message: '请输入目标主机', trigger: 'blur' }],
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        osType: [{ required: true, message: '请选择操作系统类型', trigger: 'change' }]
      },
      targetDialogVisible: false,
      editTargetIndex: -1,
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
    editTargetPortLabel() {
      return this.editTargetForm.transport === 'winrm' ? 'WinRM 端口' : 'SSH 端口'
    },
    editTargetUsernameLabel() {
      return this.editTargetForm.transport === 'winrm' ? '用户名' : 'SSH 用户名'
    },
    editTargetPasswordLabel() {
      return this.editTargetForm.transport === 'winrm' ? '密码' : 'SSH 密码'
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/hostsec/tasks'
    this.loadEnums()
    this.loadData()
  },
  methods: {
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
    onEditTargetTransportChange() {
      if (this.editTargetForm.transport === 'winrm') {
        this.editTargetForm.port = this.editTargetForm.winrmUseHttps ? 5986 : 5985
        this.editTargetForm.key = ''
      } else {
        this.editTargetForm.port = 22
      }
    },
    onEditWinrmHttpsChange() {
      if (this.editTargetForm.transport === 'winrm') {
        this.editTargetForm.port = this.editTargetForm.winrmUseHttps ? 5986 : 5985
      }
    },
    editTarget(idx) {
      const t = this.targets[idx]
      this.editTargetIndex = idx
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
      if (!this.editTargetForm.host) {
        this.$message({ message: '请输入目标主机', type: 'warning' })
        return
      }
      if (!this.editTargetForm.username) {
        this.$message({ message: '请输入用户名', type: 'warning' })
        return
      }
      if (this.editTargetForm.transport === 'winrm') {
        if (!this.editTargetForm.password) {
          this.$message({ message: 'WinRM 需要填写密码', type: 'warning' })
          return
        }
      } else if (!this.editTargetForm.password && !this.editTargetForm.key) {
        this.$message({ message: '请填写 SSH 密码或私钥', type: 'warning' })
        return
      }
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
    onListTabChange() {
      this.formData.page = 1
      this.currentpage = 1
      this.loadData()
    },
    async loadData() {
      this.tableLoading = true
      try {
        if (this.listTab === 'all') {
          await this.loadMerged()
        } else if (this.listTab === 'malware') {
          const res = await security.getMalwareTaskList({
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
              totalFindings: r.totalFindings,
              worstRiskName: r.worstRiskName,
              checkTime: r.checkTime
            }))
          } else {
            this.tableRows = []
            this.$message({ message: res.msg || '加载失败', type: 'error' })
          }
        } else {
          const scanScene = this.listTab === 'vuln' ? 2 : 1
          const res = await security.getBaselineTaskList({
            page: this.formData.page,
            size: this.pageSize,
            scanScene
          })
          if (res.code === 200 && res.data) {
            this.totalpage = res.data.total || 0
            this.tableRows = (res.data.list || []).map((r) => ({
              source: 'baseline',
              kindLabel: r.scanSceneName || '安全配置核查',
              taskId: r.taskId,
              targetIp: r.targetIp,
              osTypeName: r.osTypeName,
              totalRules: r.totalRules,
              passCount: r.passCount,
              failCount: r.failCount,
              errorCount: r.errorCount,
              checkTime: r.checkTime
            }))
          } else {
            this.tableRows = []
            this.$message({ message: res.msg || '加载失败', type: 'error' })
          }
        }
      } finally {
        this.tableLoading = false
      }
    },
    async loadMerged() {
      const take = 25
      const [bRes, mRes] = await Promise.all([
        security.getBaselineTaskList({ page: 1, size: take, scanScene: 0 }),
        security.getMalwareTaskList({ page: 1, size: take })
      ])
      const merged = []
      if (bRes.code === 200 && bRes.data && bRes.data.list) {
        for (const r of bRes.data.list) {
          merged.push({
            source: 'baseline',
            kindLabel: r.scanSceneName || '安全配置核查',
            taskId: r.taskId,
            targetIp: r.targetIp,
            osTypeName: r.osTypeName,
            totalRules: r.totalRules,
            passCount: r.passCount,
            failCount: r.failCount,
            errorCount: r.errorCount,
            checkTime: r.checkTime,
            _ts: new Date(r.checkTime.replace(/-/g, '/')).getTime()
          })
        }
      }
      if (mRes.code === 200 && mRes.data && mRes.data.list) {
        for (const r of mRes.data.list) {
          merged.push({
            source: 'malware',
            kindLabel: '恶意代码检测',
            taskId: r.taskId,
            targetIp: r.targetIp,
            totalFindings: r.totalFindings,
            worstRiskName: r.worstRiskName,
            checkTime: r.checkTime,
            _ts: new Date(r.checkTime.replace(/-/g, '/')).getTime()
          })
        }
      }
      merged.sort((a, b) => (b._ts || 0) - (a._ts || 0))
      this.tableRows = merged.slice(0, take)
      this.totalpage = this.tableRows.length
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
      try {
        if (this.taskForm.taskKind === 'malware') {
          for (const t of this.targets) {
            const res = await security.runMalwareScan({
              host: t.host,
              port: t.port,
              username: t.username,
              password: t.password,
              key: t.key,
              osType: t.osType
            })
            if (res.code !== 200) {
              this.$message({ message: `目标 ${t.host} 执行失败: ${res.msg || '未知错误'}`, type: 'error' })
            }
          }
          this.$message({ message: `恶意代码检测完成，共 ${this.targets.length} 个目标`, type: 'success' })
          this.createVisible = false
          this.formData.page = 1
          this.currentpage = 1
          this.loadData()
        } else {
          const scanScene = this.taskForm.taskKind === 'vuln' ? 2 : 1
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
            this.startProgressPolling(res.data.taskId)
          } else {
            this.$message({ message: res.msg || '创建任务失败', type: 'error' })
          }
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
            this.loadData()
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
      if (status === 'failed') return '失败'
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
      this.detailMeta = `${row.kindLabel} · 目标：${row.targetIp || ''}`
      try {
        if (row.source === 'malware') {
          this.detailMode = 'malware'
          const res = await security.getMalwareList({ taskId: row.taskId })
          if (res.code === 200) {
            this.detailRows = res.data.list || []
          } else {
            this.$message({ message: res.msg || '加载失败', type: 'error' })
          }
        } else {
          this.detailMode = 'baseline'
          const res = await security.getBaselineList({ taskId: row.taskId })
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
      this.loadData()
    },
    handleSizeChange(t) {
      this.formData.page = 1
      this.pageSize = t
      this.loadData()
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

.hub-tabs {
  margin-bottom: 12px;
}

.merge-hint {
  font-size: 12px;
  color: rgba(148, 163, 184, 0.85);
  margin: 10px 0 0;
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.progress-target-row:last-child {
  border-bottom: none;
}

.progress-target-host {
  font-size: 13px;
  color: #e2e8f0;
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
</style>
