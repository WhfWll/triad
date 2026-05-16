<template>
  <div class="security-container">
    <div class="main-title" v-if="!embedded">敏感数据发现</div>
    
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
        <el-table-column prop="dbType" label="数据库类型">
          <template slot-scope="scope">
            <span :class="getDBTypeClass(scope.row.dbType)">{{ getDBTypeName(scope.row.dbType) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="dbHost" label="数据库地址" :show-overflow-tooltip="true">
        </el-table-column>
        <el-table-column prop="totalCount" label="发现数据条数">
        </el-table-column>
        <el-table-column prop="highCount" label="高敏感数据">
          <template slot-scope="scope">
            <span class="count-high">{{ scope.row.highCount || 0 }}</span>
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

    <el-dialog title="新建敏感数据扫描任务" :visible.sync="dialogVisible" width="600px">
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
        <el-form-item label="敏感数据类型" prop="dataTypes">
          <el-select v-model="taskForm.dataTypes" multiple placeholder="请选择敏感数据类型">
            <el-option label="身份证号" :value="1"></el-option>
            <el-option label="银行卡号" :value="2"></el-option>
            <el-option label="护照号" :value="3"></el-option>
            <el-option label="手机号" :value="4"></el-option>
            <el-option label="邮箱" :value="5"></el-option>
            <el-option label="地址" :value="6"></el-option>
            <el-option label="出生日期" :value="7"></el-option>
            <el-option label="姓名" :value="8"></el-option>
            <el-option label="Token" :value="9"></el-option>
            <el-option label="证书信息" :value="10"></el-option>
            <el-option label="密码哈希" :value="11"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </span>
    </el-dialog>

    <el-dialog title="扫描结果详情" :visible.sync="detailVisible" width="900px">
      <div v-if="detailData">
        <div class="detail-header">
          <h3>{{ detailData.name }}</h3>
          <p>数据库: {{ getDBTypeName(detailData.dbType) }} - {{ detailData.dbHost }}:{{ detailData.dbPort }}</p>
        </div>
        
        <div class="detail-summary">
          <div class="summary-item">
            <span class="summary-label">总发现条数</span>
            <span class="summary-value total">{{ detailData.totalCount || 0 }}</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">高敏感</span>
            <span class="summary-value high">{{ detailData.highCount || 0 }}</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">中敏感</span>
            <span class="summary-value medium">{{ detailData.mediumCount || 0 }}</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">低敏感</span>
            <span class="summary-value low">{{ detailData.lowCount || 0 }}</span>
          </div>
        </div>

        <div class="detail-section">
          <h4>数据类型分布</h4>
          <div class="type-distribution">
            <div v-for="item in detailData.typeStats" :key="item.dataType" class="type-item">
              <span class="type-name">{{ getDataTypeName(item.dataType) }}</span>
              <span class="type-count">{{ item.count }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <h4>敏感数据详情</h4>
          <el-table :data="detailData.items || []" style="width: 100%">
            <el-table-column prop="tableName" label="表名" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="columnName" label="字段名" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="dataType" label="数据类型">
              <template slot-scope="scope">{{ getDataTypeName(scope.row.dataType) }}</template>
            </el-table-column>
            <el-table-column prop="sensitivityLevel" label="敏感等级">
              <template slot-scope="scope"><span :class="getSensitivityClass(scope.row.sensitivityLevel)">{{ getSensitivityName(scope.row.sensitivityLevel) }}</span></template>
            </el-table-column>
            <el-table-column prop="sampleData" label="示例数据" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="count" label="数量">
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
  name: 'SensitiveData',
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
        dbType: 1,
        dbHost: '',
        dbPort: 3306,
        dbName: '',
        dbUser: '',
        dbPassword: '',
        dataTypes: []
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
      const res = await security.getSensitiveDataList({
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
        dbPassword: '',
        dataTypes: []
      }
    },
    async submitForm() {
      this.$refs.taskForm.validate(async (valid) => {
        if (valid) {
          const res = await security.runSensitiveScan(this.taskForm)
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
    getDBTypeName(type) {
      const map = { 1: 'MySQL', 2: 'PostgreSQL', 3: 'MongoDB', 4: 'Redis', 5: 'CouchDB' }
      return map[type] || '未知'
    },
    getDBTypeClass(type) {
      const map = { 1: 'db-mysql', 2: 'db-postgresql', 3: 'db-mongodb', 4: 'db-redis', 5: 'db-couchdb' }
      return map[type] || 'db-default'
    },
    getStatusName(status) {
      const map = { 1: '等待扫描', 2: '扫描中', 3: '已完成' }
      return map[status] || '未知'
    },
    getStatusClass(status) {
      const map = { 1: 'status-wait', 2: 'status-running', 3: 'status-complete' }
      return map[status] || 'status-default'
    },
    getDataTypeName(type) {
      const map = { 1: '身份证号', 2: '银行卡号', 3: '护照号', 4: '手机号', 5: '邮箱', 6: '地址', 7: '出生日期', 8: '姓名', 9: 'Token', 10: '证书信息', 11: '密码哈希' }
      return map[type] || '未知'
    },
    getSensitivityName(level) {
      const map = { 1: '高敏感', 2: '中敏感', 3: '低敏感' }
      return map[level] || '未知'
    },
    getSensitivityClass(level) {
      const map = { 1: 'sensitivity-high', 2: 'sensitivity-medium', 3: 'sensitivity-low' }
      return map[level] || 'sensitivity-default'
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

.count-high {
  color: #ef4444;
  font-weight: bold;
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

.detail-summary {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.summary-item {
  flex: 1;
  padding: 15px;
  border-radius: 8px;
  text-align: center;
  background: rgba(0, 212, 170, 0.05);
  border: 1px solid rgba(0, 212, 170, 0.2);
}

.summary-label {
  display: block;
  font-size: 12px;
  color: #94a3b8;
  margin-bottom: 5px;
}

.summary-value {
  display: block;
  font-size: 28px;
  font-weight: bold;
}

.summary-value.total { color: #00d4aa; }
.summary-value.high { color: #ef4444; }
.summary-value.medium { color: #f97316; }
.summary-value.low { color: #10b981; }

.detail-section {
  margin-bottom: 20px;
}

.detail-section h4 {
  color: #00d4aa;
  margin-bottom: 15px;
  font-size: 14px;
}

.type-distribution {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.type-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(0, 212, 170, 0.1);
  border-radius: 4px;
}

.type-name {
  color: #94a3b8;
  font-size: 12px;
}

.type-count {
  color: #00d4aa;
  font-weight: bold;
}

.sensitivity-high, .sensitivity-medium, .sensitivity-low {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.sensitivity-high { background: rgba(239, 68, 68, 0.2); color: #ef4444; }
.sensitivity-medium { background: rgba(249, 115, 22, 0.2); color: #f97316; }
.sensitivity-low { background: rgba(16, 185, 129, 0.2); color: #10b981; }
</style>