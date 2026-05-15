
&lt;template&gt;
  &lt;div class="security-container"&gt;
    &lt;div class="main-title"&gt;专项应用检测&lt;/div&gt;
    
    &lt;div class="list_box"&gt;
      &lt;div class="search-box"&gt;
        &lt;div class="operationbutton"&gt;
          &lt;el-button type="primary" size="small" @click="btnCreate"&gt;新建检测任务&lt;/el-button&gt;
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
        &lt;el-table-column prop="appType" label="应用类型"&gt;
          &lt;template slot-scope="scope"&gt;
            &lt;span :class="getAppTypeClass(scope.row.appType)"&gt;{{ getAppTypeName(scope.row.appType) }}&lt;/span&gt;
          &lt;/template&gt;
        &lt;/el-table-column&gt;
        &lt;el-table-column prop="targetUrl" label="目标地址" :show-overflow-tooltip="true"&gt;
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

    &lt;el-dialog title="新建专项应用检测任务" :visible.sync="dialogVisible" width="600px"&gt;
      &lt;el-form :model="taskForm" :rules="rules" ref="taskForm" label-width="100px"&gt;
        &lt;el-form-item label="任务名称" prop="name"&gt;
          &lt;el-input v-model="taskForm.name" placeholder="请输入任务名称"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="应用类型" prop="appType"&gt;
          &lt;el-select v-model="taskForm.appType" placeholder="请选择应用类型"&gt;
            &lt;el-option label="万户OA" :value="1"&gt;&lt;/el-option&gt;
            &lt;el-option label="用友NC" :value="2"&gt;&lt;/el-option&gt;
            &lt;el-option label="蓝凌EKP" :value="3"&gt;&lt;/el-option&gt;
            &lt;el-option label="云时空" :value="4"&gt;&lt;/el-option&gt;
            &lt;el-option label="亿赛通" :value="5"&gt;&lt;/el-option&gt;
            &lt;el-option label="D-Link" :value="6"&gt;&lt;/el-option&gt;
            &lt;el-option label="通达OA" :value="7"&gt;&lt;/el-option&gt;
            &lt;el-option label="WordPress" :value="8"&gt;&lt;/el-option&gt;
            &lt;el-option label="ThinkPHP" :value="9"&gt;&lt;/el-option&gt;
            &lt;el-option label="Spring Boot" :value="10"&gt;&lt;/el-option&gt;
            &lt;el-option label="通用CMS" :value="11"&gt;&lt;/el-option&gt;
          &lt;/el-select&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="目标地址" prop="targetUrl"&gt;
          &lt;el-input v-model="taskForm.targetUrl" placeholder="请输入目标URL"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="Cookie" prop="cookie"&gt;
          &lt;el-input type="textarea" v-model="taskForm.cookie" placeholder="可选，登录后的Cookie" :rows="2"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="检测类型" prop="checkTypes"&gt;
          &lt;el-select v-model="taskForm.checkTypes" multiple placeholder="请选择检测类型"&gt;
            &lt;el-option label="远程代码执行" :value="1"&gt;&lt;/el-option&gt;
            &lt;el-option label="未授权访问" :value="2"&gt;&lt;/el-option&gt;
            &lt;el-option label="SQL注入" :value="3"&gt;&lt;/el-option&gt;
            &lt;el-option label="文件上传" :value="4"&gt;&lt;/el-option&gt;
            &lt;el-option label="弱口令" :value="5"&gt;&lt;/el-option&gt;
            &lt;el-option label="XSS" :value="6"&gt;&lt;/el-option&gt;
            &lt;el-option label="SSRF" :value="7"&gt;&lt;/el-option&gt;
            &lt;el-option label="信息泄露" :value="8"&gt;&lt;/el-option&gt;
          &lt;/el-select&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="自定义Header" prop="customHeaders"&gt;
          &lt;el-input type="textarea" v-model="taskForm.customHeaders" placeholder="可选，JSON格式的Header" :rows="2"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
      &lt;/el-form&gt;
      &lt;span slot="footer"&gt;
        &lt;el-button @click="dialogVisible = false"&gt;取消&lt;/el-button&gt;
        &lt;el-button type="primary" @click="submitForm"&gt;确定&lt;/el-button&gt;
      &lt;/span&gt;
    &lt;/el-dialog&gt;

    &lt;el-dialog title="检测结果详情" :visible.sync="detailVisible" width="900px"&gt;
      &lt;div v-if="detailData"&gt;
        &lt;div class="detail-header"&gt;
          &lt;h3&gt;{{ detailData.name }}&lt;/h3&gt;
          &lt;p&gt;目标: {{ detailData.targetUrl }} ({{ getAppTypeName(detailData.appType) }})&lt;/p&gt;
        &lt;/div&gt;
        
        &lt;div class="detail-stats"&gt;
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

        &lt;div class="detail-section"&gt;
          &lt;h4&gt;漏洞详情列表&lt;/h4&gt;
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
            &lt;el-table-column prop="description" label="描述" :show-overflow-tooltip="true"&gt;
            &lt;/el-table-column&gt;
            &lt;el-table-column prop="suggestion" label="修复建议" :show-overflow-tooltip="true"&gt;
            &lt;/el-table-column&gt;
          &lt;/el-table&gt;
        &lt;/div&gt;
      &lt;/div&gt;
      &lt;span slot="footer"&gt;
        &lt;el-button @click="detailVisible = false"&gt;关闭&lt;/el-button&gt;
      &lt;/span&gt;
    &lt;/el-dialog&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script&gt;
import security from '@/api/security.js'

export default {
  name: 'AppSpecificScan',
  data() {
    return {
      dialogVisible: false,
      detailVisible: false,
      multipleSelection: [],
      tableData: [],
      detailData: {},
      formData: {
        search: '',
        page: 1
      },
      pageSize: 10,
      currentpage: 1,
      totalpage: 0,
      taskForm: {
        name: '',
        appType: 1,
        targetUrl: '',
        cookie: '',
        checkTypes: [],
        customHeaders: ''
      },
      rules: {
        name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
        appType: [{ required: true, message: '请选择应用类型', trigger: 'change' }],
        targetUrl: [{ required: true, message: '请输入目标URL', trigger: 'blur' }]
      }
    }
  },
  mounted() {
    this.getData()
  },
  methods: {
    async getData() {
      const res = await security.getAppSpecificList({
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
        appType: 1,
        targetUrl: '',
        cookie: '',
        checkTypes: [],
        customHeaders: ''
      }
    },
    async submitForm() {
      this.$refs.taskForm.validate(async (valid) =&gt; {
        if (valid) {
          const res = await security.runAppSpecificScan(this.taskForm)
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
      const res = await security.delAppSpecificTask({ id: row.id })
      if (res.code == 200) {
        this.$message({ message: '删除成功', type: 'success' })
        this.getData()
      } else {
        this.$message({ message: res.msg, type: 'error' })
      }
    },
    async handleDetail(row) {
      const res = await security.getAppSpecificDetail({ id: row.id })
      if (res.code == 200) {
        this.detailData = res.data
        this.detailVisible = true
      } else {
        this.$message({ message: res.msg, type: 'error' })
      }
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
    getAppTypeName(type) {
      const map = { 1: '万户OA', 2: '用友NC', 3: '蓝凌EKP', 4: '云时空', 5: '亿赛通', 6: 'D-Link', 7: '通达OA', 8: 'WordPress', 9: 'ThinkPHP', 10: 'Spring Boot', 11: '通用CMS' }
      return map[type] || '未知'
    },
    getAppTypeClass(type) {
      const map = { 1: 'app-wanhui', 2: 'app-yongyou', 3: 'app-lanling', 4: 'app-yunshikong', 5: 'app-yisaitong', 6: 'app-dlink', 7: 'app-tongda', 8: 'app-wordpress', 9: 'app-thinkphp', 10: 'app-springboot', 11: 'app-generic' }
      return map[type] || 'app-default'
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
      const map = { 1: '等待检测', 2: '检测中', 3: '已完成' }
      return map[status] || '未知'
    },
    getStatusClass(status) {
      const map = { 1: 'status-wait', 2: 'status-running', 3: 'status-complete' }
      return map[status] || 'status-default'
    },
    getVulnTypeName(type) {
      const map = { 1: '远程代码执行', 2: '未授权访问', 3: 'SQL注入', 4: '文件上传', 5: '弱口令', 6: 'XSS', 7: 'SSRF', 8: '信息泄露' }
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

.app-wanhui, .app-yongyou, .app-lanling, .app-yunshikong, .app-yisaitong, .app-dlink, .app-tongda, .app-wordpress, .app-thinkphp, .app-springboot, .app-generic, .app-default {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.app-wanhui { background: rgba(245, 101, 101, 0.2); color: #f56565; }
.app-yongyou { background: rgba(97, 175, 254, 0.2); color: #61affe; }
.app-lanling { background: rgba(67, 153, 52, 0.2); color: #439934; }
.app-yunshikong { background: rgba(231, 162, 13, 0.2); color: #e7a20d; }
.app-yisaitong { background: rgba(125, 86, 205, 0.2); color: #7d56cd; }
.app-dlink { background: rgba(59, 130, 246, 0.2); color: #3b82f6; }
.app-tongda { background: rgba(249, 115, 22, 0.2); color: #f97316; }
.app-wordpress { background: rgba(212, 180, 68, 0.2); color: #d4b444; }
.app-thinkphp { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.app-springboot { background: rgba(106, 176, 76, 0.2); color: #6ab04c; }
.app-generic { background: rgba(148, 163, 184, 0.2); color: #94a3b8; }

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

.stat-item.critical .stat-value { color: #ef4444; }
.stat-item.high .stat-value { color: #f97316; }
.stat-item.medium .stat-value { color: #eab308; }
.stat-item.low .stat-value { color: #10b981; }

.detail-section {
  margin-bottom: 20px;
}

.detail-section h4 {
  color: #00d4aa;
  margin-bottom: 15px;
  font-size: 14px;
}
&lt;/style&gt;
