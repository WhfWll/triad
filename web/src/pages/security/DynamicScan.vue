<template>
  <div class="security-container">
    <div class="main-title" v-if="!embedded">动态扫描</div>
    
    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="btnCreate">新建扫描任务</el-button>
        </div>
        <div class="serach-condition">
          <div class="search-text">
            <el-input placeholder="搜索任务名称" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select" size="small" clearable></el-input>
            <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
          </div>
          <div>
            <el-button type="primary" size="small" @click="handleReset">重置</el-button>
          </div>
        </div>
      </div>

      <el-table :data="tableData" style="width: 100%" class="myTable" @selection-change="handleSelectionChange">
        <el-table-column width="55" type="selection">
        </el-table-column>
        <el-table-column prop="name" label="任务名称" :show-overflow-tooltip="true">
        </el-table-column>
        <el-table-column prop="targetUrl" label="目标地址" :show-overflow-tooltip="true">
        </el-table-column>
        <el-table-column prop="pageCount" label="爬取页面数">
        </el-table-column>
        <el-table-column prop="vulnCount" label="发现漏洞数">
        </el-table-column>
        <el-table-column prop="riskLevel" label="风险等级">
          <template slot-scope="scope">
            <span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template slot-scope="scope">
            <span :class="getStatusClass(scope.row.status)">{{ getStatusName(scope.row.status) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间">
        </el-table-column>
        <el-table-column prop="scanTime" label="扫描时间">
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="handleDetail(scope.row)">详情</el-link>
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
        @size-change="handleSizeChange">
      </el-pagination>
    </div>

    <el-dialog title="新建动态扫描任务" :visible.sync="dialogVisible" width="650px">
      <el-form :model="taskForm" :rules="rules" ref="taskForm" label-width="120px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="taskForm.name" placeholder="请输入任务名称"></el-input>
        </el-form-item>
        <el-form-item label="目标地址" prop="targetUrl">
          <el-input v-model="taskForm.targetUrl" placeholder="请输入目标URL"></el-input>
        </el-form-item>
        <el-form-item label="扫描模式" prop="scanMode">
          <el-radio-group v-model="taskForm.scanMode">
            <el-radio :label="1">CrawlerX爬虫</el-radio>
            <el-radio :label="2">MITM代理模式</el-radio>
            <el-radio :label="3">两者结合</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="最大深度" prop="maxDepth" v-if="taskForm.scanMode !== 2">
          <el-input-number v-model="taskForm.maxDepth" :min="1" :max="10" :step="1"></el-input-number>
        </el-form-item>
        <el-form-item label="并发数" prop="concurrency">
          <el-input-number v-model="taskForm.concurrency" :min="1" :max="20" :step="1"></el-input-number>
        </el-form-item>
        <el-form-item label="Cookie" prop="cookie">
          <el-input type="textarea" v-model="taskForm.cookie" placeholder="可选，登录后的Cookie" :rows="2"></el-input>
        </el-form-item>
        <el-form-item label="登录配置" prop="loginConfig">
          <div class="login-config">
            <el-checkbox v-model="taskForm.enableLogin">启用自动登录</el-checkbox>
            <div v-if="taskForm.enableLogin" class="login-form">
              <el-input v-model="taskForm.loginUrl" placeholder="登录页面URL" size="small" style="margin-bottom: 10px"></el-input>
              <el-input v-model="taskForm.loginUsername" placeholder="用户名" size="small" style="margin-bottom: 10px"></el-input>
              <el-input type="password" v-model="taskForm.loginPassword" placeholder="密码" size="small"></el-input>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="检测类型" prop="checkTypes">
          <el-select v-model="taskForm.checkTypes" multiple placeholder="请选择检测类型">
            <el-option label="SQL注入" :value="1"></el-option>
            <el-option label="XSS" :value="2"></el-option>
            <el-option label="SSRF" :value="3"></el-option>
            <el-option label="XXE" :value="4"></el-option>
            <el-option label="命令注入" :value="5"></el-option>
            <el-option label="文件包含" :value="6"></el-option>
            <el-option label="文件上传" :value="7"></el-option>
            <el-option label="CSRF" :value="8"></el-option>
            <el-option label="信息泄露" :value="9"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="代理配置" prop="proxyConfig">
          <div class="proxy-config">
            <el-checkbox v-model="taskForm.enableProxy">使用代理</el-checkbox>
            <div v-if="taskForm.enableProxy" class="proxy-form">
              <el-select v-model="taskForm.proxyType" placeholder="代理类型" size="small" style="margin-bottom: 10px; width: 100%">
                <el-option label="HTTP" :value="1"></el-option>
                <el-option label="HTTPS" :value="2"></el-option>
                <el-option label="SOCKS5" :value="3"></el-option>
              </el-select>
              <el-input v-model="taskForm.proxyHost" placeholder="代理地址" size="small" style="margin-bottom: 10px"></el-input>
              <el-input-number v-model="taskForm.proxyPort" :min="1" :max="65535" :step="1" size="small" placeholder="端口"></el-input-number>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </span>
    </el-dialog>

    <el-dialog title="扫描结果详情" :visible.sync="detailVisible" width="1000px">
      <div v-if="detailData">
        <div class="detail-header">
          <h3>{{ detailData.name }}</h3>
          <p>目标: {{ detailData.targetUrl }}</p>
        </div>
        
        <div class="detail-stats">
          <div class="stat-item pages">
            <span class="stat-label">爬取页面数</span>
            <span class="stat-value">{{ detailData.pageCount || 0 }}</span>
          </div>
          <div class="stat-item critical">
            <span class="stat-label">严重漏洞</span>
            <span class="stat-value">{{ detailData.criticalCount || 0 }}</span>
          </div>
          <div class="stat-item high">
            <span class="stat-label">高危漏洞</span>
            <span class="stat-value">{{ detailData.highRiskCount || 0 }}</span>
          </div>
          <div class="stat-item medium">
            <span class="stat-label">中危漏洞</span>
            <span class="stat-value">{{ detailData.middleRiskCount || 0 }}</span>
          </div>
          <div class="stat-item low">
            <span class="stat-label">低危漏洞</span>
            <span class="stat-value">{{ detailData.lowRiskCount || 0 }}</span>
          </div>
        </div>

        <el-tabs v-model="activeTab" class="detail-tabs">
          <el-tab-pane label="漏洞列表" name="vulns">
            <el-table :data="detailData.vulns || []" style="width: 100%">
              <el-table-column prop="name" label="漏洞名称" :show-overflow-tooltip="true">
              </el-table-column>
              <el-table-column prop="type" label="漏洞类型" :show-overflow-tooltip="true">
                <template slot-scope="scope">{{ getVulnTypeName(scope.row.type) }}</template>
              </el-table-column>
              <el-table-column prop="riskLevel" label="风险等级">
                <template slot-scope="scope">
                  <span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="url" label="漏洞URL" :show-overflow-tooltip="true">
              </el-table-column>
              <el-table-column label="操作" width="120">
                <template slot-scope="scope">
                  <el-button type="text" size="small" @click="showVulnDetail(scope.row)">查看详情</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="站点地图" name="sitemap">
            <div class="sitemap-tree">
              <el-tree
                :data="siteMapTree"
                :props="treeProps"
                :expand-on-click-node="false"
                node-key="path">
              </el-tree>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
      <span slot="footer">
        <el-button @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>

    <el-dialog title="漏洞详情" :visible.sync="vulnDetailVisible" width="700px">
      <div v-if="currentVuln">
        <h4 style="color: #00d4aa; margin-bottom: 15px">{{ currentVuln.name }}</h4>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="漏洞类型">{{ getVulnTypeName(currentVuln.type) }}</el-descriptions-item>
          <el-descriptions-item label="风险等级">
            <span :class="getRiskClass(currentVuln.riskLevel)" style="padding: 4px 12px">{{ getRiskName(currentVuln.riskLevel) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="漏洞URL">{{ currentVuln.url }}</el-descriptions-item>
          <el-descriptions-item label="请求方式">{{ currentVuln.method || 'GET' }}</el-descriptions-item>
        </el-descriptions>
        
        <div class="detail-section">
          <h4 style="color: #00d4aa; margin-top: 20px; margin-bottom: 10px">请求报文</h4>
          <pre class="code-block">{{ currentVuln.request || '暂无' }}</pre>
        </div>

        <div class="detail-section">
          <h4 style="color: #00d4aa; margin-bottom: 10px">响应报文</h4>
          <pre class="code-block">{{ currentVuln.response || '暂无' }}</pre>
        </div>

        <div class="detail-section">
          <h4 style="color: #00d4aa; margin-bottom: 10px">修复建议</h4>
          <p style="color: #94a3b8; line-height: 1.6">{{ currentVuln.suggestion || '暂无' }}</p>
        </div>
      </div>
      <span slot="footer">
        <el-button @click="vulnDetailVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'

export default {
  name: 'DynamicScan',
  props: {
    embedded: {
      type: Boolean,
      default: false
    }
  },
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
      this.$refs.taskForm.validate(async (valid) => {
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
      this.$confirm('确认删除该任务？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$message({ message: '删除成功', type: 'success' })
        this.getData()
      }).catch(() => {})
    },
    async handleDetail(row) {
      this.detailData = row
      this.siteMapTree = this.buildSiteMapTree(row.pages || [])
      this.detailVisible = true
      this.activeTab = 'vulns'
    },
    showVulnDetail(vuln) {
      this.currentVuln = vuln
      this.vulnDetailVisible = true
    },
    buildSiteMapTree(pages) {
      const tree = []
      const map = {}
      pages.forEach(page => {
        const url = page.url || page
        const parts = url.replace(/^https?:\/\//, '').split('/').filter(Boolean)
        let current = tree
        parts.forEach((part, index) => {
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
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

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
  background: #1a1a2e;
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
</style>
