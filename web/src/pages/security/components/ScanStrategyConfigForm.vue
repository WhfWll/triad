<template>
  <div class="scan-strategy-config-form">
    <template v-if="showPanel('vuln')">
      <el-divider v-if="!section" content-position="left">漏洞脚本选择</el-divider>
      <div class="vuln-selector">
        <div class="vuln-toolbar">
          <div class="vuln-filter">
            <el-select v-model="vulnFilter.class" placeholder="漏洞分类" size="small" clearable @change="fetchVulns" style="width: 130px">
              <el-option v-for="item in enumOptions.class" :key="'c-' + item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-select v-model="vulnFilter.type" placeholder="漏洞类型" size="small" clearable @change="fetchVulns" style="width: 130px">
              <el-option v-for="item in enumOptions.type" :key="'t-' + item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-select v-model="vulnFilter.risk" placeholder="风险等级" size="small" clearable @change="fetchVulns" style="width: 120px">
              <el-option v-for="item in enumOptions.risk" :key="'r-' + item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-input v-model="vulnFilter.search" placeholder="搜索漏洞名称或CVE" size="small" clearable style="width: 220px" @keydown.enter.native="fetchVulns" />
            <el-button size="small" @click="fetchVulns">搜索</el-button>
          </div>
          <div class="vuln-actions">
            <span class="vuln-count">已选 <strong>{{ selectedVulnIds.length }}</strong> 个</span>
            <span class="vuln-hint">表头勾选 = 当前页</span>
            <el-button size="small" type="text" :loading="selectAllAllLoading" @click="selectAllAllPages">全选</el-button>
            <el-button size="small" type="text" @click="clearAllVulns">清空</el-button>
          </div>
        </div>
        <p v-if="vulnLoadError" class="vuln-load-error">{{ vulnLoadError }}</p>
        <el-table
          ref="vulnTable"
          v-loading="vulnLoading"
          :data="vulnList"
          row-key="id"
          style="width: 100%"
          class="vuln-table"
          :max-height="vulnTableMaxHeight"
          @selection-change="onVulnSelectionChange"
        >
          <el-table-column type="selection" width="45" reserve-selection />
          <el-table-column prop="name" label="漏洞名称" min-width="200" :show-overflow-tooltip="true" />
          <el-table-column prop="classEnum" label="分类" width="90" />
          <el-table-column prop="typeEnum" label="类型" width="100" />
          <el-table-column prop="riskName" label="风险" width="70">
            <template slot-scope="scope">
              <span :class="riskClass(scope.row.risk)">{{ scope.row.riskName }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="vulNum" label="编号" width="160" :show-overflow-tooltip="true" />
        </el-table>
        <el-pagination
          v-if="vulnTotal > vulnPageSize"
          :page-size="vulnPageSize"
          layout="total, prev, pager, next"
          :total="vulnTotal"
          :current-page="vulnPage"
          @current-change="onVulnPageChange"
          small
          background
        />
      </div>
    </template>

    <template v-if="showPanel('login')">
      <el-form-item label-width="0" class="switch-row">
        <el-switch v-model="websiteLogin.isOpen" active-text="开启" inactive-text="关闭" />
        <span class="form-tip">扫描需登录页面时，可配置 Cookie 或自定义 Header</span>
      </el-form-item>
      <template v-if="websiteLogin.isOpen">
        <div class="login-toolbar">
          <el-button size="small" type="primary" @click="addLoginRow">新增登录</el-button>
        </div>
        <el-table :data="websiteLogin.list" size="small" class="login-table" empty-text="暂无登录凭证">
          <el-table-column prop="target" label="登录地址" min-width="160">
            <template slot-scope="scope">
              <span v-if="!scope.row._editing">{{ scope.row.target || '-' }}</span>
              <el-input v-else v-model="scope.row.target" size="small" placeholder="http://example.com/login" />
            </template>
          </el-table-column>
          <el-table-column prop="verifyType" label="凭证类型" width="130">
            <template slot-scope="scope">
              <span v-if="!scope.row._editing">{{ loginTypeLabel(scope.row.verifyType) }}</span>
              <el-select v-else v-model="scope.row.verifyType" size="small" style="width: 110px">
                <el-option v-for="item in loginTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column prop="verifyValue" label="凭证" min-width="140">
            <template slot-scope="scope">
              <span v-if="!scope.row._editing" class="login-value-preview">{{ scope.row.verifyValue || '-' }}</span>
              <el-input
                v-else
                v-model="scope.row.verifyValue"
                type="textarea"
                :rows="2"
                size="small"
                placeholder="Cookie 或 Header 内容"
              />
            </template>
          </el-table-column>
          <el-table-column prop="verifyStatusZh" label="状态" width="88">
            <template slot-scope="scope">
              <span v-if="scope.row.verifyStatusZh" :class="loginStatusClass(scope.row.verifyStatus)">
                {{ scope.row.verifyStatusZh }}
              </span>
              <span v-else class="login-status-pending">-</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="140">
            <template slot-scope="scope">
              <el-link v-if="scope.row._editing" :underline="false" @click="saveLoginRow(scope)">保存</el-link>
              <el-link v-else :underline="false" @click="editLoginRow(scope)">编辑</el-link>
              <el-link :underline="false" @click="removeLoginRow(scope)">删除</el-link>
            </template>
          </el-table-column>
        </el-table>
      </template>
    </template>

    <template v-if="showPanel('crawler')">
      <el-divider v-if="!section" content-position="left">爬虫配置</el-divider>
      <el-form-item label-width="0" class="switch-row">
        <el-switch v-model="config.webCrawler.isOpen" active-text="开启" inactive-text="关闭" />
      </el-form-item>
      <template v-if="config.webCrawler.isOpen">
        <el-form-item label="爬取深度">
          <el-input-number v-model="config.webCrawler.maxDepth" :min="1" :max="10" :step="1" size="small" />
        </el-form-item>
        <el-form-item label="爬取范围">
          <el-select v-model="config.webCrawler.scanRange" size="small" style="width: 240px">
            <el-option label="全域名扫描" :value="0" />
            <el-option label="目标URL和子目录" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="爬取速度">
          <el-select v-model="config.webCrawler.crawlerSpeed" size="small" style="width: 240px">
            <el-option label="低速" :value="1" />
            <el-option label="中速" :value="2" />
            <el-option label="高速" :value="3" />
          </el-select>
        </el-form-item>
      </template>
    </template>

    <template v-if="showPanel('port')">
      <el-divider v-if="!section" content-position="left">端口扫描配置</el-divider>
      <el-form-item label-width="0" class="switch-row">
        <el-switch v-model="config.portScan.isOpen" active-text="开启" inactive-text="关闭" />
      </el-form-item>
      <template v-if="config.portScan.isOpen">
        <el-form-item label="扫描端口">
          <el-select v-model="portRangeType" size="small" style="width: 240px" @change="onPortRangeChange">
            <el-option label="TOP10 端口" :value="10" />
            <el-option label="TOP20 端口" :value="20" />
            <el-option label="TOP50 端口" :value="50" />
            <el-option label="TOP100 端口（默认）" :value="100" />
            <el-option label="TOP500 端口" :value="500" />
            <el-option label="TOP1000 端口" :value="1000" />
            <el-option label="全部端口 (1-65535)" :value="65535" />
            <el-option label="自定义端口" :value="0" />
          </el-select>
          <el-input
            v-if="portRangeType === 0"
            v-model="config.portScan.scanPort"
            placeholder="如: 80,443,3306,8000-9000"
            size="small"
            style="width: 240px; margin-left: 8px"
          />
        </el-form-item>
        <el-form-item label="扫描方式">
          <el-select v-model="config.portScan.tcpScanType" size="small" style="width: 200px">
            <el-option label="TCP-Connect" :value="1" />
            <el-option label="TCP SYN" :value="2" />
            <el-option label="TCP FIN" :value="3" />
            <el-option label="TCP ACK" :value="4" />
            <el-option label="TCP NULL" :value="5" />
            <el-option label="UDP" :value="6" />
          </el-select>
        </el-form-item>
        <el-form-item label="超时时间">
          <el-select v-model="config.portScan.timeout" size="small" style="width: 160px">
            <el-option label="3s" :value="3" />
            <el-option label="5s" :value="5" />
            <el-option label="10s（默认）" :value="10" />
            <el-option label="20s" :value="20" />
            <el-option label="30s" :value="30" />
            <el-option label="60s" :value="60" />
            <el-option label="120s" :value="120" />
          </el-select>
        </el-form-item>
      </template>
    </template>

    <template v-if="showPanel('weakpass')">
      <el-divider v-if="!section" content-position="left">弱口令扫描</el-divider>
      <el-form-item label-width="0" class="switch-row">
        <el-switch
          v-model="weakPass.isOpen"
          active-text="开启"
          inactive-text="关闭"
        />
        <span class="form-tip">开启后对已选服务进行弱口令爆破检测</span>
      </el-form-item>
      <template v-if="weakPass.isOpen">
        <div class="weakpass-toolbar">
          <span class="weakpass-count">已选 <strong>{{ weakPass.services.length }}</strong> 个协议</span>
          <el-button size="small" type="text" @click="selectAllWeakPassServices">全选</el-button>
          <el-button size="small" type="text" @click="clearWeakPassServices">清空</el-button>
        </div>
        <div v-loading="weakPassServicesLoading" class="weakpass-services">
          <el-checkbox-group v-model="weakPass.services" class="weakpass-check-group">
            <el-checkbox
              v-for="item in weakPassServiceOptions"
              :key="item.value"
              :label="item.value"
            >
              {{ item.label }}
            </el-checkbox>
          </el-checkbox-group>
          <p v-if="!weakPassServicesLoading && !weakPassServiceOptions.length" class="weakpass-empty">
            未加载到协议列表，请刷新页面或检查场景枚举接口
          </p>
        </div>
      </template>
    </template>

    <template v-if="showPanel('advanced')">
      <el-divider v-if="!section" content-position="left">高级配置</el-divider>
      <el-form-item label="代理配置">
        <el-switch v-model="config.proxy.isOpen" active-text="使用代理" inactive-text="关闭" />
        <div v-if="config.proxy.isOpen" class="sub-form" style="margin-top: 10px">
          <el-select v-model="config.proxy.proto" size="small" style="width: 100px">
            <el-option label="HTTP" :value="1" />
            <el-option label="HTTPS" :value="2" />
            <el-option label="SOCKS5" :value="3" />
          </el-select>
          <el-input v-model="config.proxy.addr" placeholder="代理地址" size="small" style="width: 180px" />
          <el-input v-model="config.proxy.port" placeholder="端口" size="small" style="width: 100px" />
        </div>
      </el-form-item>
    </template>
  </div>
</template>

<script>
import { vulnerability } from '@/api/tool.js'
import task from '@/api/task.js'
import scene from '@/api/scene.js'
import { PORT_RANGE_MAP } from '../appsecPortRanges.js'
import { getStrategySections } from '../appsecBuiltinStrategies.js'
import {
  cloneDefaultWeakPass,
  WEAKPASS_SERVICE_FALLBACK
} from '../appsecWeakPassDefaults.js'

export default {
  name: 'ScanStrategyConfigForm',
  props: {
    config: { type: Object, required: true },
    strategyId: { type: String, default: 'builtin-full' },
    /** scan=仅扫描/爬虫/端口；vuln=仅漏洞脚本；all=全部 */
    pageMode: { type: String, default: 'all', validator: v => ['scan', 'vuln', 'all'].includes(v) },
    /** 单面板：login | port | crawler | advanced | vuln */
    section: { type: String, default: '' }
  },
  data() {
    return {
      loginTypeOptions: [],
      enumOptions: { class: [], type: [], risk: [] },
      vulnList: [],
      syncingVulnIds: false,
      restoringVulnSelection: false,
      restoreSelectionToken: 0,
      vulnTotal: 0,
      vulnPage: 1,
      vulnPageSize: 50,
      selectedVulnIds: [],
      vulnFilter: { class: null, type: null, risk: null, search: '' },
      portRangeType: 100,
      vulnLoading: false,
      vulnLoadError: '',
      selectAllAllLoading: false,
      weakPassServiceOptions: [],
      weakPassServicesLoading: false
    }
  },
  computed: {
    sections() {
      return getStrategySections(this.strategyId)
    },
    visibleSections() {
      const base = this.sections
      if (this.pageMode === 'scan') {
        return { ...base, vuln: false }
      }
      if (this.pageMode === 'vuln') {
        return {
          vuln: base.vuln,
          scan: false,
          crawler: false,
          port: false,
          advanced: false,
          weakPass: false
        }
      }
      return base
    },
    vulnTableMaxHeight() {
      if (this.section === 'vuln' || this.pageMode === 'vuln') return 560
      return 360
    },
    websiteLogin() {
      this.ensureWebsiteLogin()
      return this.config.websiteLogin
    },
    weakPass() {
      this.ensureWeakPass()
      return this.config.weakPass
    }
  },
  watch: {
    'config.vulIdsConfig': {
      immediate: true,
      deep: true,
      handler(ids) {
        if (this.syncingVulnIds) return
        const next = [...(ids || [])]
        if (next.length === this.selectedVulnIds.length &&
          next.every((id, i) => String(id) === String(this.selectedVulnIds[i]))) {
          return
        }
        this.selectedVulnIds = next
        if (this.showPanel('vuln') && this.vulnList.length) {
          this.restoreSelection()
        }
      }
    },
    strategyId() {
      this.fetchVulns()
    },
    section(val) {
      if (val === 'vuln' || this.pageMode === 'vuln') {
        this.applyVulnIdsFromConfig()
        this.fetchVulns()
      }
      if (val === 'weakpass') {
        this.fetchWeakPassServices()
      }
    }
  },
  async mounted() {
    this.ensureWebsiteLogin()
    this.ensureWeakPass()
    const tasks = [this.fetchEnums(), this.fetchLoginTypes()]
    if (this.sections.weakPass) tasks.push(this.fetchWeakPassServices())
    await Promise.all(tasks)
    this.applyVulnIdsFromConfig()
    this.guessPortRangeType()
    if (this.showPanel('vuln')) {
      this.fetchVulns()
    }
  },
  methods: {
    ensureWebsiteLogin() {
      if (!this.config) return
      if (!this.config.websiteLogin) {
        this.$set(this.config, 'websiteLogin', { isOpen: false, list: [] })
      } else if (!Array.isArray(this.config.websiteLogin.list)) {
        this.$set(this.config.websiteLogin, 'list', [])
      }
    },
    ensureWeakPass() {
      if (!this.config) return
      if (!this.config.weakPass) {
        this.$set(this.config, 'weakPass', cloneDefaultWeakPass())
      } else if (!Array.isArray(this.config.weakPass.services)) {
        this.$set(this.config.weakPass, 'services', [])
      }
    },
    async fetchWeakPassServices() {
      if (!this.showPanel('weakpass')) return
      this.weakPassServicesLoading = true
      try {
        const res = await scene.getSceneEnum()
        const list = res && res.code === 200 && res.data && res.data.weakPass && res.data.weakPass.services
        if (Array.isArray(list) && list.length) {
          this.weakPassServiceOptions = list
            .filter(item => item && item.value != null && item.label !== '通用')
            .map(item => ({
              value: Number(item.value),
              label: String(item.label)
            }))
        } else {
          this.weakPassServiceOptions = [...WEAKPASS_SERVICE_FALLBACK]
        }
      } catch {
        this.weakPassServiceOptions = [...WEAKPASS_SERVICE_FALLBACK]
      } finally {
        this.weakPassServicesLoading = false
      }
    },
    selectAllWeakPassServices() {
      this.weakPass.services = this.weakPassServiceOptions.map(x => x.value)
    },
    clearWeakPassServices() {
      this.weakPass.services = []
    },
    async fetchLoginTypes() {
      try {
        const res = await task.taskEnum()
        if (res && res.code === 200 && res.data && res.data.webLoginType) {
          this.loginTypeOptions = this.normalizeEnumList(res.data.webLoginType)
          if (this.loginTypeOptions.length && this.websiteLogin.list.length) {
            const def = this.loginTypeOptions[0].value
            this.websiteLogin.list.forEach(row => {
              if (row.verifyType == null || row.verifyType === '') row.verifyType = def
            })
          }
        }
      } catch {
        this.loginTypeOptions = [
          { label: 'Cookie', value: 1 },
          { label: 'Header', value: 2 }
        ]
      }
    },
    loginTypeLabel(value) {
      const hit = this.loginTypeOptions.find(x => x.value === value)
      return (hit && hit.label) || '-'
    },
    loginStatusClass(status) {
      if (status === 1) return 'login-status-ok'
      if (status === 2) return 'login-status-warn'
      return 'login-status-fail'
    },
    addLoginRow() {
      const defType = (this.loginTypeOptions[0] && this.loginTypeOptions[0].value) || 1
      this.websiteLogin.list.push({
        target: '',
        verifyType: defType,
        verifyValue: '',
        verifyStatus: 0,
        verifyStatusZh: '',
        _editing: true
      })
    },
    editLoginRow(scope) {
      this.$set(this.websiteLogin.list[scope.$index], '_editing', true)
    },
    removeLoginRow(scope) {
      this.websiteLogin.list.splice(scope.$index, 1)
    },
    async saveLoginRow(scope) {
      const row = this.websiteLogin.list[scope.$index]
      if (!row.target || !String(row.target).trim()) {
        this.$message.warning('请填写登录地址')
        return
      }
      if (!row.verifyValue || !String(row.verifyValue).trim()) {
        this.$message.warning('请填写凭证内容')
        return
      }
      try {
        const res = await task.websitelogincheck({
          task_check_target: row.target,
          target: row.target,
          verifyType: row.verifyType,
          verifyValue: row.verifyValue
        })
        if (res.code == 200) {
          row.verifyStatus = res.data.statusCode
          row.verifyStatusZh = res.data.status
          row._editing = false
          this.$message.success('凭证校验完成')
        } else {
          this.$message.error(res.msg || '校验失败')
        }
      } catch {
        row._editing = false
        this.$message.warning('校验接口不可用，已保存凭证（未校验）')
      }
    },
    getSelectedVulnIds() {
      return [...this.selectedVulnIds]
    },
    applyVulnIdsFromConfig() {
      const ids = [...((this.config && this.config.vulIdsConfig) || [])]
      this.selectedVulnIds = ids
      if (this.vulnList.length) {
        this.restoreSelection()
      }
    },
    showPanel(key) {
      const areaMap = {
        login: 'login',
        port: 'port',
        crawler: 'crawler',
        advanced: 'advanced',
        vuln: 'vuln',
        weakpass: 'weakPass'
      }
      const area = areaMap[key]
      if (this.section) {
        return this.section === key && this.visibleSections[area]
      }
      return Boolean(this.visibleSections[area])
    },
    normalizeEnumList(raw) {
      if (!raw) return []
      if (Array.isArray(raw)) {
        return raw.map(item => {
          if (item && typeof item === 'object') {
            return {
              label: String(item.label != null ? item.label : item.name || ''),
              value: item.value
            }
          }
          return { label: String(item), value: item }
        })
      }
      if (typeof raw === 'object') {
        return Object.keys(raw).map(k => ({
          label: String(raw[k]),
          value: /^\d+$/.test(k) ? Number(k) : k
        }))
      }
      return []
    },
    async fetchEnums() {
      try {
        const res = await vulnerability.getVulObjectlist()
        if (res && res.code === 200 && res.data) {
          this.enumOptions = {
            class: this.normalizeEnumList(res.data.class),
            type: this.normalizeEnumList(res.data.type),
            risk: this.normalizeEnumList(res.data.risk)
          }
        } else {
          this.useFallbackEnumOptions()
        }
      } catch {
        this.useFallbackEnumOptions()
      }
    },
    useFallbackEnumOptions() {
      this.enumOptions = {
        class: this.normalizeEnumList({
          1: 'Web应用', 2: '数据库', 3: '中间件', 4: '操作系统', 5: '安全设备', 6: '应用软件'
        }),
        type: this.normalizeEnumList({
          1: 'SQL注入', 2: 'XSS', 3: 'SSRF', 4: '命令执行', 5: '文件上传', 6: '文件包含',
          7: '信息泄露', 8: '未授权访问', 9: '弱口令', 10: 'CSRF', 11: 'XXE', 12: '逻辑漏洞',
          13: '远程代码执行', 14: '组件漏洞'
        }),
        risk: this.normalizeEnumList({
          0: '严重', 1: '高危', 2: '中危', 3: '低危', 4: '信息'
        })
      }
    },
    buildVulnQueryParams(page, size) {
      return {
        page,
        size,
        libName: this.vulnFilter.search || undefined,
        libClass: this.vulnFilter.class || undefined,
        libType: this.vulnFilter.type || undefined,
        libRisk: this.vulnFilter.risk !== null && this.vulnFilter.risk !== '' ? this.vulnFilter.risk : undefined
      }
    },
    async fetchVulns() {
      if (!this.showPanel('vuln')) return
      this.vulnLoading = true
      this.vulnLoadError = ''
      try {
        const res = await vulnerability.getObjectData(
          this.buildVulnQueryParams(this.vulnPage, this.vulnPageSize)
        )
        if (res && res.code === 200) {
          this.vulnList = (res.data && res.data.list) || []
          this.vulnTotal = (res.data && res.data.total) || 0
          this.applyVulnIdsFromConfig()
          this.restoreSelection()
        } else {
          this.vulnList = []
          this.vulnTotal = 0
          this.vulnLoadError = (res && res.msg) || '加载漏洞列表失败'
        }
      } catch (e) {
        this.vulnList = []
        this.vulnTotal = 0
        this.vulnLoadError = '无法连接漏洞库接口，请检查网络或登录状态'
        console.error('[ScanStrategyConfigForm] fetchVulns', e)
      } finally {
        this.vulnLoading = false
      }
    },
    restoreSelection() {
      const token = ++this.restoreSelectionToken
      this.$nextTick(() => {
        if (token !== this.restoreSelectionToken) return
        const table = this.$refs.vulnTable
        if (!table || !this.vulnList.length) return
        this.restoringVulnSelection = true
        this.syncingVulnIds = true
        table.clearSelection()
        this.vulnList
          .filter(v => this.isIdSelected(v.id))
          .forEach(row => table.toggleRowSelection(row, true))
        // clearSelection 会在 syncing 结束后异步触发空的 selection-change，误删当前页已选 ID
        this.$nextTick(() => {
          setTimeout(() => {
            if (token !== this.restoreSelectionToken) return
            this.syncingVulnIds = false
            this.restoringVulnSelection = false
          }, 0)
        })
      })
    },
    syncVulIdsToConfig() {
      this.syncingVulnIds = true
      this.config.vulIdsConfig = [...this.selectedVulnIds]
      this.$nextTick(() => {
        this.syncingVulnIds = false
      })
    },
    isIdSelected(id) {
      const sid = String(id)
      return this.selectedVulnIds.some(x => String(x) === sid)
    },
    onVulnSelectionChange(rows) {
      if (this.syncingVulnIds || this.restoringVulnSelection) return
      const pageIdSet = new Set(this.vulnList.map(r => String(r.id)))
      const selectedOnPage = rows.map(r => r.id)
      const kept = this.selectedVulnIds.filter(id => !pageIdSet.has(String(id)))
      this.selectedVulnIds = [...kept, ...selectedOnPage]
      this.syncVulIdsToConfig()
    },
    async selectAllAllPages() {
      if (!this.vulnTotal) {
        this.$message({ message: '当前筛选条件下没有可选漏洞', type: 'warning' })
        return
      }
      try {
        if (this.vulnTotal > 200) {
          await this.$confirm(
            `将勾选符合当前筛选条件的全部 ${this.vulnTotal} 个漏洞（含所有分页），是否继续？`,
            '全选确认',
            { type: 'info', confirmButtonText: '确定', cancelButtonText: '取消' }
          )
        }
      } catch {
        return
      }

      this.selectAllAllLoading = true
      try {
        const pageSize = this.vulnPageSize
        const totalPages = Math.ceil(this.vulnTotal / pageSize)
        const allIds = []
        for (let p = 1; p <= totalPages; p++) {
          const res = await vulnerability.getObjectData(this.buildVulnQueryParams(p, pageSize))
          if (res && res.code === 200) {
            const list = (res.data && res.data.list) || []
            list.forEach(row => {
              if (row && row.id != null) allIds.push(row.id)
            })
          } else {
            throw new Error((res && res.msg) || '加载失败')
          }
        }
        this.selectedVulnIds = [...new Set([...this.selectedVulnIds, ...allIds])]
        this.syncVulIdsToConfig()
        this.restoreSelection()
        this.$message({
          message: `已勾选全部 ${allIds.length} 个漏洞（当前筛选共 ${this.vulnTotal} 条）`,
          type: 'success'
        })
      } catch (e) {
        this.$message({ message: e.message || '全选全部失败，请稍后重试', type: 'error' })
      } finally {
        this.selectAllAllLoading = false
      }
    },
    clearAllVulns() {
      this.selectedVulnIds = []
      this.syncVulIdsToConfig()
      this.restoreSelection()
    },
    onVulnPageChange(page) {
      this.vulnPage = page
      this.fetchVulns()
    },
    onPortRangeChange(val) {
      if (val !== 0) this.config.portScan.scanPort = PORT_RANGE_MAP[val] || ''
    },
    guessPortRangeType() {
      const port = this.config.portScan && this.config.portScan.scanPort
      if (!port) {
        this.portRangeType = 100
        return
      }
      const hit = Object.entries(PORT_RANGE_MAP).find(([, v]) => v === port)
      this.portRangeType = hit ? Number(hit[0]) : 0
    },
    riskClass(risk) {
      return ({ 0: 'risk-critical', 1: 'risk-high', 2: 'risk-medium', 3: 'risk-low', 4: 'risk-info' })[risk] || ''
    }
  }
}
</script>

<style lang="less" scoped>
.switch-row {
  margin-bottom: 16px;
}
.form-tip {
  color: #64748b;
  font-size: 12px;
  margin-left: 12px;
}
.sub-form {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.vuln-selector {
  padding: 16px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  border: 1px solid rgba(0, 212, 170, 0.1);
  margin-bottom: 10px;
}
.vuln-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
  gap: 8px;
}
.vuln-filter {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
}
.vuln-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}
.vuln-count {
  color: #94a3b8;
  font-size: 13px;
}
.vuln-count strong {
  color: #00d4aa;
}
.vuln-hint {
  color: #64748b;
  font-size: 12px;
  margin-right: 4px;
}
.vuln-load-error {
  color: #f87171;
  font-size: 13px;
  margin: 0 0 10px;
}
.vuln-table,
.login-table {
  border: 1px solid rgba(0, 212, 170, 0.1);
  border-radius: 10px;
  overflow: hidden;
}
/deep/ .el-divider__text {
  color: #00d4aa;
  font-weight: 600;
  background: transparent;
}
/deep/ .el-divider {
  background-color: rgba(255, 255, 255, 0.08);
}
/deep/ .vuln-table .el-table__header th {
  background: rgba(0, 0, 0, 0.4) !important;
  color: #94a3b8;
}
/deep/ .vuln-table .el-table__body td {
  background: rgba(0, 0, 0, 0.2) !important;
  color: #cbd5e1;
}
/deep/ .vuln-table .el-table--enable-row-hover .el-table__body tr:hover > td {
  background: rgba(0, 212, 170, 0.05) !important;
}
.login-toolbar {
  margin-bottom: 12px;
}
.login-value-preview {
  display: inline-block;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}
.login-status-ok {
  color: #34d399;
}
.login-status-warn {
  color: #fbbf24;
}
.login-status-fail {
  color: #f87171;
}
.login-status-pending {
  color: #64748b;
}
.weakpass-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}
.weakpass-count {
  color: #94a3b8;
  font-size: 13px;
  margin-right: 8px;
}
.weakpass-count strong {
  color: #00d4aa;
}
.weakpass-services {
  padding: 16px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  border: 1px solid rgba(0, 212, 170, 0.1);
  min-height: 80px;
}
.weakpass-check-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 20px;
}
.weakpass-empty {
  color: #64748b;
  font-size: 13px;
  margin: 0;
}
/deep/ .weakpass-check-group .el-checkbox {
  color: #94a3b8;
  margin-right: 0;
}
/deep/ .weakpass-check-group .el-checkbox__label {
  color: #cbd5e1;
}
/deep/ .login-table .el-table__header th {
  background: rgba(0, 0, 0, 0.4) !important;
  color: #94a3b8;
}
/deep/ .login-table .el-table__body td {
  background: rgba(0, 0, 0, 0.2) !important;
  color: #cbd5e1;
}
</style>
