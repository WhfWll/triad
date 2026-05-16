<template>
  <div class="security-container">
    <div class="main-title" v-if="!embedded">专项应用检测</div>
    
    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="btnCreate">新建检测任务</el-button>
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
        <el-table-column prop="appType" label="应用类型">
          <template slot-scope="scope">
            <span :class="getAppTypeClass(scope.row.appType)">{{ getAppTypeName(scope.row.appType) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="targetUrl" label="目标地址" :show-overflow-tooltip="true">
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

    <el-dialog title="新建专项应用检测任务" :visible.sync="dialogVisible" width="600px">
      <el-form :model="taskForm" :rules="rules" ref="taskForm" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="taskForm.name" placeholder="请输入任务名称"></el-input>
        </el-form-item>
        <el-form-item label="应用类型" prop="appType">
          <el-select v-model="taskForm.appType" placeholder="请选择应用类型">
            <el-option label="万户OA" :value="1"></el-option>
            <el-option label="用友NC" :value="2"></el-option>
            <el-option label="蓝凌EKP" :value="3"></el-option>
            <el-option label="云时空" :value="4"></el-option>
            <el-option label="亿赛通" :value="5"></el-option>
            <el-option label="D-Link" :value="6"></el-option>
            <el-option label="通达OA" :value="7"></el-option>
            <el-option label="WordPress" :value="8"></el-option>
            <el-option label="ThinkPHP" :value="9"></el-option>
            <el-option label="Spring Boot" :value="10"></el-option>
            <el-option label="通用CMS" :value="11"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="目标地址" prop="targetUrl">
          <el-input v-model="taskForm.targetUrl" placeholder="请输入目标URL"></el-input>
        </el-form-item>
        <el-form-item label="Cookie" prop="cookie">
          <el-input type="textarea" v-model="taskForm.cookie" placeholder="可选，登录后的Cookie" :rows="2"></el-input>
        </el-form-item>
        <el-form-item label="检测类型" prop="checkTypes">
          <el-select v-model="taskForm.checkTypes" multiple placeholder="请选择检测类型">
            <el-option label="远程代码执行" :value="1"></el-option>
            <el-option label="未授权访问" :value="2"></el-option>
            <el-option label="SQL注入" :value="3"></el-option>
            <el-option label="文件上传" :value="4"></el-option>
            <el-option label="弱口令" :value="5"></el-option>
            <el-option label="XSS" :value="6"></el-option>
            <el-option label="SSRF" :value="7"></el-option>
            <el-option label="信息泄露" :value="8"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="自定义Header" prop="customHeaders">
          <el-input type="textarea" v-model="taskForm.customHeaders" placeholder="可选，JSON格式的Header" :rows="2"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </span>
    </el-dialog>

    <el-dialog title="检测结果详情" :visible.sync="detailVisible" width="900px">
      <div v-if="detailData">
        <div class="detail-header">
          <h3>{{ detailData.name }}</h3>
          <p>目标: {{ detailData.targetUrl }} ({{ getAppTypeName(detailData.appType) }})</p>
        </div>
        
        <div class="detail-stats">
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

        <div class="detail-section">
          <h4>漏洞详情列表</h4>
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
            <el-table-column prop="description" label="描述" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="suggestion" label="修复建议" :show-overflow-tooltip="true">
            </el-table-column>
          </el-table>
        </div>
      </div>
      <span slot="footer">
        <el-button @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'

export default {
  name: 'AppSpecificScan',
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
      this.$refs.taskForm.validate(async (valid) => {
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
      this.detailVisible = true
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
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

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
</style>
