<template>
  <div class="security-container">
    <div class="main-title">数据库安全检查</div>
    
    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="btnCreate">新建检查任务</el-button>
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
        <el-table-column prop="dbType" label="数据库类型">
          <template slot-scope="scope">
            <span :class="getDBTypeClass(scope.row.dbType)">{{ getDBTypeName(scope.row.dbType) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="dbHost" label="数据库地址" :show-overflow-tooltip="true">
        </el-table-column>
        <el-table-column prop="dbName" label="数据库名称" :show-overflow-tooltip="true">
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
        <el-table-column prop="checkTime" label="检查时间">
        </el-table-column>
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="handleDetail(scope.row)">详情</el-link>
            <el-link :underline="false" class="link_danger" @click="handleDel(scope.row)">删除</el-link>
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

    <el-dialog title="新建数据库安全检查任务" :visible.sync="dialogVisible" width="600px">
      <el-form :model="taskForm" :rules="rules" ref="taskForm" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="taskForm.name" placeholder="请输入任务名称"></el-input>
        </el-form-item>
        <el-form-item label="数据库类型" prop="dbType">
          <el-select v-model="taskForm.dbType" placeholder="请选择数据库类型">
            <el-option label="MySQL" :value="1"></el-option>
            <el-option label="PostgreSQL" :value="2"></el-option>
            <el-option label="MongoDB" :value="3"></el-option>
            <el-option label="Redis" :value="4"></el-option>
            <el-option label="CouchDB" :value="5"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="数据库地址" prop="dbHost">
          <el-input v-model="taskForm.dbHost" placeholder="请输入数据库地址"></el-input>
        </el-form-item>
        <el-form-item label="端口" prop="dbPort">
          <el-input type="number" v-model="taskForm.dbPort" placeholder="请输入端口号"></el-input>
        </el-form-item>
        <el-form-item label="数据库名称" prop="dbName">
          <el-input v-model="taskForm.dbName" placeholder="请输入数据库名称"></el-input>
        </el-form-item>
        <el-form-item label="用户名" prop="dbUser">
          <el-input v-model="taskForm.dbUser" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="dbPassword">
          <el-input type="password" v-model="taskForm.dbPassword" placeholder="请输入密码"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </span>
    </el-dialog>

    <el-dialog title="检查结果详情" :visible.sync="detailVisible" width="800px">
      <div v-if="detailData">
        <div class="detail-header">
          <h3>{{ detailData.name }}</h3>
          <p>数据库: {{ getDBTypeName(detailData.dbType) }} - {{ detailData.dbHost }}:{{ detailData.dbPort }}</p>
        </div>
        <div class="detail-stats">
          <div class="stat-item critical">
            <span class="stat-label">严重风险</span>
            <span class="stat-value">{{ detailData.criticalCount || 0 }}</span>
          </div>
          <div class="stat-item high">
            <span class="stat-label">高危风险</span>
            <span class="stat-value">{{ detailData.highRiskCount || 0 }}</span>
          </div>
          <div class="stat-item medium">
            <span class="stat-label">中危风险</span>
            <span class="stat-value">{{ detailData.middleRiskCount || 0 }}</span>
          </div>
          <div class="stat-item low">
            <span class="stat-label">低危风险</span>
            <span class="stat-value">{{ detailData.lowRiskCount || 0 }}</span>
          </div>
        </div>
        <el-table :data="detailData.items || []" style="width: 100%">
          <el-table-column prop="category" label="检查类别">
            <template slot-scope="scope">{{ getCategoryName(scope.row.category) }}</template>
          </el-table-column>
          <el-table-column prop="riskLevel" label="风险等级">
            <template slot-scope="scope"><span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span></template>
          </el-table-column>
          <el-table-column prop="result" label="检查结果" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="description" label="描述" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="suggestion" label="修复建议" :show-overflow-tooltip="true">
          </el-table-column>
        </el-table>
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
  name: 'DBSecurity',
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
        dbType: 1,
        dbHost: '',
        dbPort: 3306,
        dbName: '',
        dbUser: '',
        dbPassword: ''
      },
      rules: {
        name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
        dbType: [{ required: true, message: '请选择数据库类型', trigger: 'change' }],
        dbHost: [{ required: true, message: '请输入数据库地址', trigger: 'blur' }],
        dbUser: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        dbPassword: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      }
    }
  },
  mounted() {
    this.getData()
  },
  methods: {
    async getData() {
      const res = await security.getDBCheckList({
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
        dbType: 1,
        dbHost: '',
        dbPort: 3306,
        dbName: '',
        dbUser: '',
        dbPassword: ''
      }
    },
    async submitForm() {
      this.$refs.taskForm.validate(async (valid) => {
        if (valid) {
          const res = await security.runDBCheck(this.taskForm)
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
      const res = await security.delDBCheckTask({ id: row.id })
      if (res.code == 200) {
        this.$message({ message: '删除成功', type: 'success' })
        this.getData()
      } else {
        this.$message({ message: res.msg, type: 'error' })
      }
    },
    async handleDetail(row) {
      const res = await security.getDBCheckDetail({ id: row.id })
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
    getDBTypeName(type) {
      const map = { 1: 'MySQL', 2: 'PostgreSQL', 3: 'MongoDB', 4: 'Redis', 5: 'CouchDB' }
      return map[type] || '未知'
    },
    getDBTypeClass(type) {
      const map = { 1: 'db-mysql', 2: 'db-postgresql', 3: 'db-mongodb', 4: 'db-redis', 5: 'db-couchdb' }
      return map[type] || 'db-default'
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
      const map = { 1: '等待检查', 2: '检查中', 3: '已完成' }
      return map[status] || '未知'
    },
    getStatusClass(status) {
      const map = { 1: 'status-wait', 2: 'status-running', 3: 'status-complete' }
      return map[status] || 'status-default'
    },
    getCategoryName(category) {
      const map = { 1: '认证安全', 2: '授权管理', 3: '配置安全', 4: '审计日志', 5: '网络安全', 6: '加密保护' }
      return map[category] || '未知'
    }
  }
}
</script>

<style lang="less" scoped>
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

.el-table--enable-row-hover .el-table__body tr:hover>td {
  background-color: rgba(0, 212, 170, 0.1);
}

.link_primary {
  color: #00d4aa;
  margin-right: 15px;
}

.link_danger {
  color: #ef4444;
}

.db-mysql, .db-postgresql, .db-mongodb, .db-redis, .db-couchdb, .db-default {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.db-mysql { background: rgba(245, 101, 101, 0.2); color: #f56565; }
.db-postgresql { background: rgba(97, 175, 254, 0.2); color: #61affe; }
.db-mongodb { background: rgba(67, 153, 52, 0.2); color: #439934; }
.db-redis { background: rgba(231, 162, 13, 0.2); color: #e7a20d; }
.db-couchdb { background: rgba(125, 86, 205, 0.2); color: #7d56cd; }

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
</style>