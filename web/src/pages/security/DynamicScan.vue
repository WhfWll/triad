
&lt;template&gt;
  &lt;div class="security-container"&gt;
    &lt;div class="main-title"&gt;动态扫描&lt;/div&gt;
    
    &lt;div class="list_box"&gt;
      &lt;div class="search-box"&gt;
        &lt;div class="operationbutton"&gt;
          &lt;el-button type="primary" size="small" @click="btnCreate"&gt;新建扫描任务&lt;/el-button&gt;
        &lt;/div&gt;
        &lt;div class="serach-condition"&gt;
          &lt;div class="search-text"&gt;
            &lt;el-input placeholder="搜索任务名称" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select" size="small" clearable&gt;&lt;/el-input&gt;
            &lt;el-button type="primary" size="small" @click="handlesearch"&gt;搜索&lt;/el-button&gt;
          &lt;/div&gt;
          &lt;div&gt;
            &lt;el-button type="primary" size="small" @click="handleReset"&gt;重置&lt;/el-button&gt;
          &lt;/div&gt;
        &lt;/div&gt;
      &lt;/div&gt;

      &lt;el-table :data="tableData" style="width: 100%" class="myTable" @selection-change="handleSelectionChange"&gt;
        &lt;el-table-column width="55" type="selection"&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="name" label="任务名称" :show-overflow-tooltip="true"&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="targetUrl" label="目标地址" :show-overflow-tooltip="true"&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="pageCount" label="爬取页面数"&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="vulnCount" label="发现漏洞数"&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="riskLevel" label="风险等级"&gt;
          &lt;template slot-scope="scope"&gt;
            &lt;span :class="getRiskClass(scope.row.riskLevel)"&gt;{{ getRiskName(scope.row.riskLevel) }}&lt;/span&gt;
          &lt;/template&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="status" label="状态"&gt;
          &lt;template slot-scope="scope"&gt;
            &lt;span :class="getStatusClass(scope.row.status)"&gt;{{ getStatusName(scope.row.status) }}&lt;/span&gt;
          &lt;/template&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="createTime" label="创建时间"&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="scanTime" label="扫描时间"&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column label="操作"&gt;
          &lt;template slot-scope="scope"&gt;
            &lt;el-link :underline="false" class="link_primary" @click="handleDetail(scope.row)"&gt;详情&lt;/el-link&gt;
            &lt;el-link :underline="false" class="link_danger" @click="handleDel(scope.row)"&gt;删除&lt;/el-link&gt;
          &lt;/template&gt;
        &lt;/el-table-column&gt;
      &lt;/el-table&gt;

      &lt;el-pagination
        :page-size="pageSize"
        background
        layout="total, prev, pager, next, sizes, jumper"
        :total="totalpage"
        :current-page="currentpage"
        @current-change="handlecurrentchange"
        @size-change="handleSizeChange"&gt;
      &lt;/el-pagination&gt;
    &lt;/div&gt;

    &lt;el-dialog title="新建动态扫描任务" :visible.sync="dialogVisible" width="650px"&gt;
      &lt;el-form :model="taskForm" :rules="rules" ref="taskForm" label-width="120px"&gt;
        &lt;el-form-item label="任务名称" prop="name"&gt;
          &lt;el-input v-model="taskForm.name" placeholder="请输入任务名称"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="目标地址" prop="targetUrl"&gt;
          &lt;el-input v-model="taskForm.targetUrl" placeholder="请输入目标URL"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="扫描模式" prop="scanMode"&gt;
          &lt;el-radio-group v-model="taskForm.scanMode"&gt;
            &lt;el-radio :label="1"&gt;CrawlerX爬虫&lt;/el-radio&gt;
            &lt;el-radio :label="2"&gt;MITM代理模式&lt;/el-radio&gt;
            &lt;el-radio :label="3"&gt;两者结合&lt;/el-radio&gt;
          &lt;/el-radio-group&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="最大深度" prop="maxDepth" v-if="taskForm.scanMode !== 2"&gt;
          &lt;el-input-number v-model="taskForm.maxDepth" :min="1" :max="10" :step="1"&gt;&lt;/el-input-number&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="并发数" prop="concurrency"&gt;
          &lt;el-input-number v-model="taskForm.concurrency" :min="1" :max="20" :step="1"&gt;&lt;/el-input-number&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="Cookie" prop="cookie"&gt;
          &lt;el-input type="textarea" v-model="taskForm.cookie" placeholder="可选，登录后的Cookie" :rows="2"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="登录配置" prop="loginConfig"&gt;
          &lt;div class="login-config"&gt;
            &lt;el-checkbox v-model="taskForm.enableLogin"&gt;启用自动登录&lt;/el-checkbox&gt;
            &lt;div v-if="taskForm.enableLogin" class="login-form"&gt;
              &lt;el-input v-model="taskForm.loginUrl" placeholder="登录页面URL" size="small" style="margin-bottom: 10px"&gt;&lt;/el-input&gt;
              &lt;el-input v-model="taskForm.loginUsername" placeholder="用户名" size="small" style="margin-bottom: 10px"&gt;&lt;/el-input&gt;
              &lt;el-input type="password" v-model="taskForm.loginPassword" placeholder="密码" size="small"&gt;&lt;/el-input&gt;
            &lt;/div&gt;
          &lt;/div&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="检测类型" prop="checkTypes"&gt;
          &lt;el-select v-model="taskForm.checkTypes" multiple placeholder="请选择检测类型"&gt;
            &lt;el-option label="SQL注入" :value="1"&gt;&lt;/el-option&gt;
            &lt;el-option label="XSS" :value="2"&gt;&lt;/el-option&gt;
            &lt;el-option label="SSRF" :value="3"&gt;&lt;/el-option&gt;
            &lt;el-option label="XXE" :value="4"&gt;&lt;/el-option&gt;
            &lt;el-option label="命令注入" :value="5"&gt;&lt;/el-option&gt;
            &lt;el-option label="文件包含" :value="6"&gt;&lt;/el-option&gt;
            &lt;el-option label="文件上传" :value="7"&gt;&lt;/el-option&gt;
            &lt;el-option label="CSRF" :value="8"&gt;&lt;/el-option&gt;
            &lt;el-option label="信息泄露" :value="9"&gt;&lt;/el-option&gt;
          &lt;/el-select&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="代理配置" prop="proxyConfig"&gt;
          &lt;div class="proxy-config"&gt;
            &lt;el-checkbox v-model="taskForm.enableProxy"&gt;使用代理&lt;/el-checkbox&gt;
            &lt;div v-if="taskForm.enableProxy" class="proxy-form"&gt;
              &lt;el-select v-model="taskForm.proxyType" placeholder="代理类型" size="small" style="margin-bottom: 10px; width: 100%"&gt;
                &lt;el-option label="HTTP" :value="1"&gt;&lt;/el-option&gt;
                &lt;el-option label="HTTPS" :value="2"&gt;&lt;/el-option&gt;
                &lt;el-option label="SOCKS5" :value="3"&gt;&lt;/el-option&gt;
              &lt;/el-select&gt;
              &lt;el-input v-model="taskForm.proxyHost" placeholder="代理地址" size="small" style="margin-bottom: 10px"&gt;&lt;/el-input&gt;
              &lt;el-input-number v-model="taskForm.proxyPort" :min="1" :max="65535" :step="1" size="small" placeholder="端口"&gt;&lt;/el-input-number&gt;
            &lt;/div&gt;
          &lt;/div&gt;
        &lt;/el-form-item&gt;
      &lt;/el-form&gt;
      &lt;span slot="footer"&gt;
        &lt;el-button @click="dialogVisible = false"&gt;取消&lt;/el-button&gt;
        &lt;el-button type="primary" @click="submitForm"&gt;确定&lt;/el-button&gt;
      &lt;/span&gt;
    &lt;/el-dialog&gt;

    &lt;el-dialog title="扫描结果详情" :visible.sync="detailVisible" width="1000px"&gt;
      &lt;div v-if="detailData"&gt;
        &lt;div class="detail-header"&gt;
          &lt;h3&gt;{{ detailData.name }}&lt;/h3&gt;
          &lt;p&gt;目标: {{ detailData.targetUrl }}&lt;/p&gt;
        &lt;/div&gt;
        
        &lt;div class="detail-stats"&gt;
          &lt;div class="stat-item pages"&gt;
            &lt;span class="stat-label"&gt;爬取页面数&lt;/span&gt;
            &lt;span class="stat-value"&gt;{{ detailData.pageCount || 0 }}&lt;/span&gt;
          &lt;/div&gt;
          &lt;div class="stat-item critical"&gt;
            &lt;span class="stat-label"&gt;严重漏洞&lt;/span&gt;
            &lt;span class="stat-value"&gt;{{ detailData.criticalCount || 0 }}&lt;/span&gt;
          &lt;/div&gt;
          &lt;div class="stat-item high"&gt;
            &lt;span class="stat-label"&gt;高危漏洞&lt;/span&gt;
            &lt;span class="stat-value"&gt;{{ detailData.highRiskCount || 0 }}&lt;/span&gt;
          &lt;/div&gt;
          &lt;div class="stat-item medium"&gt;
            &lt;span class="stat-label"&gt;中危漏洞&lt;/span&gt;
            &lt;span class="stat-value"&gt;{{ detailData.middleRiskCount || 0 }}&lt;/span&gt;
          &lt;/div&gt;
          &lt;div class="stat-item low"&gt;
            &lt;span class="stat-label"&gt;低危漏洞&lt;/span&gt;
            &lt;span class="stat-value"&gt;{{ detailData.lowRiskCount || 0 }}&lt;/span&gt;
          &lt;/div&gt;
        &lt;/div&gt;

        &lt;el-tabs v-model="activeTab" class="detail-tabs"&gt;
          &lt;el-tab-pane label="漏洞列表" name="vulns"&gt;
            &lt;el-table :data="detailData.vulns || []" style="width: 100%"&gt;
              &lt;el-table-column prop="name" label="漏洞名称" :show-overflow-tooltip="true"&gt;
              &lt;/el-table-column&gt;
              &lt;el-table-column prop="type" label="漏洞类型" :show-overflow-tooltip="true"&gt;
                &lt;template slot-scope="scope"&gt;{{ getVulnTypeName(scope.row.type) }}&lt;/template&gt;
              &lt;/el-table-column&gt;
              &lt;el-table-column prop="riskLevel" label="风险等级"&gt;
                &lt;template slot-scope="scope"&gt;
                  &lt;span :class="getRiskClass(scope.row.riskLevel)"&gt;{{ getRiskName(scope.row.riskLevel) }}&lt;/span&gt;
                &lt;/template&gt;
              &lt;/el-table-column&gt;
              &lt;el-table-column prop="url" label="漏洞URL" :show-overflow-tooltip="true"&gt;
              &lt;/el-table-column&gt;
              &lt;el-table-column label="操作" width="120"&gt;
                &lt;template slot-scope="scope"&gt;
                  &lt;el-button type="text" size="small" @click="showVulnDetail(scope.row)"&gt;查看详情&lt;/el-button&gt;
                &lt;/template&gt;
              &lt;/el-table-column&gt;
            &lt;/el-table&gt;
          &lt;/el-tab-pane&gt;
          &lt;el-tab-pane label="站点地图" name="sitemap"&gt;
            &lt;div class="sitemap-tree"&gt;
              &lt;el-tree
                :data="siteMapTree"
                :props="treeProps"
                :expand-on-click-node="false"
                node-key="path"&gt;
              &lt;/el-tree&gt;
            &lt;/div&gt;
          &lt;/el-tab-pane&gt;
        &lt;/el-tabs&gt;
      &lt;/div&gt;
      &lt;span slot="footer"&gt;
        &lt;el-button @click="detailVisible = false"&gt;关闭&lt;/el-button&gt;
      &lt;/span&gt;
    &lt;/el-dialog&gt;

    &lt;el-dialog title="漏洞详情" :visible.sync="vulnDetailVisible" width="700px"&gt;
      &lt;div v-if="currentVuln"&gt;
        &lt;h4 style="color: #00d4aa; margin-bottom: 15px"&gt;{{ currentVuln.name }}&lt;/h4&gt;
        &lt;el-descriptions :column="1" border&gt;
          &lt;el-descriptions-item label="漏洞类型"&gt;{{ getVulnTypeName(currentVuln.type) }}&lt;/el-descriptions-item&gt;
          &lt;el-descriptions-item label="风险等级"&gt;
            &lt;span :class="getRiskClass(currentVuln.riskLevel)" style="padding: 4px 12px"&gt;{{ getRiskName(currentVuln.riskLevel) }}&lt;/span&gt;
          &lt;/el-descriptions-item&gt;
          &lt;el-descriptions-item label="漏洞URL"&gt;{{ currentVuln.url }}&lt;/el-descriptions-item&gt;
          &lt;el-descriptions-item label="请求方式"&gt;{{ currentVuln.method || 'GET' }}&lt;/el-descriptions-item&gt;
        &lt;/el-descriptions&gt;
        
        &lt;div class="detail-section"&gt;
          &lt;h4 style="color: #00d4aa; margin-top: 20px; margin-bottom: 10px"&gt;请求报文&lt;/h4&gt;
          &lt;pre class="code-block"&gt;{{ currentVuln.request || '暂无' }}&lt;/pre&gt;
        &lt;/div&gt;

        &lt;div class="detail-section"&gt;
          &lt;h4 style="color: #00d4aa; margin-bottom: 10px"&gt;响应报文&lt;/h4&gt;
          &lt;pre class="code-block"&gt;{{ currentVuln.response || '暂无' }}&lt;/pre&gt;
        &lt;/div&gt;

        &lt;div class="detail-section"&gt;
          &lt;h4 style="color: #00d4aa; margin-bottom: 10px"&gt;修复建议&lt;/h4&gt;
          &lt;p style="color: #94a3b8; line-height: 1.6"&gt;{{ currentVuln.suggestion || '暂无' }}&lt;/p&gt;
        &lt;/div&gt;
      &lt;/div&gt;
      &lt;span slot="footer"&gt;
        &lt;el-button @click="vulnDetailVisible = false"&gt;关闭&lt;/el-button&gt;
      &lt;/span&gt;
    &lt;/el-dialog&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script&gt;
import security from '@/api/security.js'

export default {
  name: 'DynamicScan',
  data() {
    return {
      dialogVisible: false,
      detailVisible: false,
      vulnDetailVisible: false,
      activeTab: 'vulns',
      multipleSelection: [],
      tableData: [],
      detailData: {},
      currentVuln: null,
      siteMapTree: [],
      treeProps: {
        label: 'name',
        children: 'children'
      },
      formData: {
        search: '',
        page: 1
      },
      pageSize: 10,
      currentpage: 1,
      totalpage: 0,
      taskForm: {
        name: '',
        targetUrl: '',
        scanMode: 1,
        maxDepth: 3,
        concurrency: 5,
        cookie: '',
        enableLogin: false,
        loginUrl: '',
        loginUsername: '',
        loginPassword: '',
        checkTypes: [],
        enableProxy: false,
        proxyType: 1,
        proxyHost: '',
        proxyPort: 8080
      },
      rules: {
        name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
        targetUrl: [{ required: true, message: '请输入目标URL', trigger: 'blur' }]
      }
    }
  },
  mounted() {
    this.getData()
  },
  methods: {
    async getData() {
      const res = await security.getDynamicScanList({
        page: this.formData.page,
        size: this.pageSize,
        search: this.formData.search
      })
      if (res.code == 200) {
        this.tableData = res.data.list
        this.totalpage = res.data.total
      } else {
        this.$message({ message: res.msg, type: 'error' })
      }
    },
    btnCreate() {
      this.dialogVisible = true
      this.taskForm = {
        name: '',
        targetUrl: '',
        scanMode: 1,
        maxDepth: 3,
        concurrency: 5,
        cookie: '',
        enableLogin: false,
        loginUrl: '',
        loginUsername: '',
        loginPassword: '',
        checkTypes: [],
        enableProxy: false,
        proxyType: 1,
        proxyHost: '',
        proxyPort: 8080
      }
    },
    async submitForm() {
      this.$refs.taskForm.validate(async (valid) =&gt; {
        if (valid) {
          const res = await security.runDynamicScan(this.taskForm)
          if (res.code == 200) {
            this.$message({ message: '任务创建成功', type: 'success' })
            this.dialogVisible = false
            this.getData()
          } else {
            this.$message({ message: res.msg, type: 'error' })
          }
        }
      })
    },
    async handleDel(row) {
      const res = await security.delDynamicScanTask({ id: row.id })
      if (res.code == 200) {
        this.$message({ message: '删除成功', type: 'success' })
        this.getData()
      } else {
        this.$message({ message: res.msg, type: 'error' })
      }
    },
    async handleDetail(row) {
      const res = await security.getDynamicScanDetail({ id: row.id })
      if (res.code == 200) {
        this.detailData = res.data
        this.siteMapTree = this.buildSiteMapTree(res.data.pages || [])
        this.detailVisible = true
        this.activeTab = 'vulns'
      } else {
        this.$message({ message: res.msg, type: 'error' })
      }
    },
    showVulnDetail(vuln) {
      this.currentVuln = vuln
      this.vulnDetailVisible = true
    },
    buildSiteMapTree(pages) {
      const tree = []
      const map = {}
      pages.forEach(page =&gt; {
        const url = page.url || page
        const parts = url.replace(/^https?:\/\//, '').split('/').filter(Boolean)
        let current = tree
        parts.forEach((part, index) =&gt; {
          const path = parts.slice(0, index + 1).join('/')
          if (!map[path]) {
            const node = {
              name: part,
              path: path,
              children: []
            }
            map[path] = node
            current.push(node)
          }
          current = map[path].children
        })
      })
      return tree
    },
    handlesearch() {
      this.formData.page = 1
      this.getData()
      this.currentpage = 1
    },
    handleReset() {
      this.formData.page = 1
      this.formData.search = ''
      this.getData()
      this.currentpage = 1
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
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    getRiskName(risk) {
      const map = { 0: '严重', 1: '高危', 2: '中危', 3: '低危', 4: '信息' }
      return map[risk] || '未知'
    },
    getRiskClass(risk) {
      const map = { 0: 'risk-critical', 1: 'risk-high', 2: 'risk-medium', 3: 'risk-low', 4: 'risk-info' }
      return map[risk] || 'risk-default'
    },
    getStatusName(status) {
      const map = { 1: '等待扫描', 2: '扫描中', 3: '已完成' }
      return map[status] || '未知'
    },
    getStatusClass(status) {
      const map = { 1: 'status-wait', 2: 'status-running', 3: 'status-complete' }
      return map[status] || 'status-default'
    },
    getVulnTypeName(type) {
      const map = { 1: 'SQL注入', 2: 'XSS', 3: 'SSRF', 4: 'XXE', 5: '命令注入', 6: '文件包含', 7: '文件上传', 8: 'CSRF', 9: '信息泄露' }
      return map[type] || '未知'
    }
  }
}
&lt;/script&gt;

&lt;style lang="less" scoped&gt;
.security-container {
  padding: 20px;
}

.main-title {
  font-size: 18px;
  font-weight: bold;
  color: #00d4aa;
  margin-bottom: 20px;
}

.list_box {
  background: #1a1d24;
  border-radius: 8px;
  padding: 24px;
}

.search-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.operationbutton {
  display: flex;
  gap: 10px;
}

.serach-condition {
  display: flex;
  gap: 15px;
}

.search-text {
  display: flex;
  align-items: center;
  gap: 10px;
}

.myTable {
  background: #1a1d24;
}

.el-table__header-wrapper,
.el-table__body-wrapper {
  background: #1a1d24;
}

.el-table th,
.el-table td {
  color: #94a3b8;
  border-bottom: 1px solid #2d3748;
}

.el-table--enable-row-hover .el-table__body tr:hover&gt;td {
  background-color: rgba(0, 212, 170, 0.1);
}

.link_primary {
  color: #00d4aa;
  margin-right: 15px;
}

.link_danger {
  color: #ef4444;
}

.risk-critical, .risk-high, .risk-medium, .risk-low, .risk-info {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.risk-critical { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
.risk-high { background: rgba(249, 115, 22, 0.2); color: #f97316; }
.risk-medium { background: rgba(234, 179, 8, 0.2); color: #eab308; }
.risk-low { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.risk-info { background: rgba(148, 163, 184, 0.2); color: #94a3b8; }

.status-wait, .status-running, .status-complete {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-wait { background: rgba(234, 179, 8, 0.2); color: #eab308; }
.status-running { background: rgba(59, 130, 246, 0.2); color: #3b82f6; }
.status-complete { background: rgba(16, 185, 129, 0.2); color: #10b981; }

.detail-header {
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #2d3748;
}

.detail-header h3 {
  color: #00d4aa;
  margin-bottom: 5px;
}

.detail-header p {
  color: #94a3b8;
}

.detail-stats {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.stat-item {
  flex: 1;
  padding: 15px;
  border-radius: 8px;
  text-align: center;
}

.stat-item.pages { background: rgba(0, 212, 170, 0.1); border: 1px solid rgba(0, 212, 170, 0.3); }
.stat-item.critical { background: rgba(239, 68, 68, 0.1); border: 1px solid rgba(239, 68, 68, 0.3); }
.stat-item.high { background: rgba(249, 115, 22, 0.1); border: 1px solid rgba(249, 115, 22, 0.3); }
.stat-item.medium { background: rgba(234, 179, 8, 0.1); border: 1px solid rgba(234, 179, 8, 0.3); }
.stat-item.low { background: rgba(16, 185, 129, 0.1); border: 1px solid rgba(16, 185, 129, 0.3); }

.stat-label {
  display: block;
  font-size: 12px;
  color: #94a3b8;
  margin-bottom: 5px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: bold;
}

.stat-item.pages .stat-value { color: #00d4aa; }
.stat-item.critical .stat-value { color: #ef4444; }
.stat-item.high .stat-value { color: #f97316; }
.stat-item.medium .stat-value { color: #eab308; }
.stat-item.low .stat-value { color: #10b981; }

.detail-tabs {
  margin-top: 20px;
}

.sitemap-tree {
  max-height: 400px;
  overflow-y: auto;
  background: rgba(0, 0, 0, 0.1);
  padding: 15px;
  border-radius: 4px;
}

.code-block {
  background: #1a1d24;
  padding: 15px;
  border-radius: 4px;
  color: #94a3b8;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  overflow-x: auto;
  max-height: 300px;
  overflow-y: auto;
}

.login-config, .proxy-config {
  width: 100%;
}

.login-form, .proxy-form {
  margin-top: 10px;
  padding-left: 20px;
}
&lt;/style&gt;
