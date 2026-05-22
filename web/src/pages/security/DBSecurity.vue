<template>
  <div class="security-container">
    <div class="main-title" v-if="!embedded">数据库安全检查</div>
    
    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="btnCreate">新建检查任务</el-button>
          <router-link to="/datasec/targets"><el-button size="small">目标库管理</el-button></router-link>
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
        <el-table-column prop="dbHost" label="扫描目标" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <span>{{ scope.row.targetSummary || scope.row.dbHost }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="riskLevel" label="风险等级" width="88">
          <template slot-scope="scope">
            <span :class="getRiskClass(scope.row.riskLevel)">{{ getRiskName(scope.row.riskLevel) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="基线" width="96">
          <template slot-scope="scope">
            <span v-if="scope.row.baselineTotal">{{ scope.row.baselineFail || 0 }}/{{ scope.row.baselineTotal }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="cveMatchCount" label="CVE" width="64">
          <template slot-scope="scope">
            <span v-if="scope.row.cveMatchCount">{{ scope.row.cveMatchCount }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="totalCount" label="敏感" width="64">
          <template slot-scope="scope">
            <span v-if="scope.row.totalCount">{{ scope.row.totalCount }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="88">
          <template slot-scope="scope">
            <span :class="getStatusClass(scope.row.status)">{{ getStatusName(scope.row.status) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="160">
        </el-table-column>
        <el-table-column prop="checkTime" label="检查时间" width="160">
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="handleDetail(scope.row)">详情</el-link>
            <el-link :underline="false" class="link_primary" @click="handleRerun(scope.row)">再次检测</el-link>
            <el-link :underline="false" class="link_primary" @click="handleCopyTargets(scope.row)">复制目标</el-link>
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

    <el-dialog title="新建数据库安全检查任务" :visible.sync="dialogVisible" width="720px">
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
        <data-sec-target-list
          v-model="taskForm.targets"
          :db-type="taskForm.dbType"
          :library-picks="taskForm.libraryPicks"
          @pick-library="pickerVisible = true"
          @remove-library="removeLibraryPick"
          @import-db-type="taskForm.dbType = $event"
        />
        <el-form-item label="敏感数据">
          <el-checkbox v-model="taskForm.scanSensitive">同时扫描敏感数据（库表字段）</el-checkbox>
        </el-form-item>
        <el-form-item v-if="taskForm.scanSensitive" label="扫描范围">
          <el-checkbox v-model="taskForm.scanAllDb">扫描实例下全部库（默认仅当前库）</el-checkbox>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button @click="saveTargetsToLibrary">保存到目标库</el-button>
        <el-button :loading="testConnLoading" @click="testAllConnections">测试全部连接</el-button>
        <el-button type="primary" @click="submitForm">开始检查（{{ targetTotal }} 个目标）</el-button>
      </span>
    </el-dialog>

    <data-sec-target-picker
      :visible.sync="pickerVisible"
      :db-type="taskForm.dbType"
      :exclude-ids="taskForm.libraryTargetIds"
      @pick="onPickLibrary"
    />
  </div>
</template>

<script>
import security from '@/api/security.js'
import DataSecTargetList from './components/DataSecTargetList.vue'
import DataSecTargetPicker from './components/DataSecTargetPicker.vue'

export default {
  name: 'DBSecurity',
  components: { DataSecTargetList, DataSecTargetPicker },
  props: {
    embedded: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      dialogVisible: false,
      pickerVisible: false,
      testConnLoading: false,
      multipleSelection: [],
      tableData: [],
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
        targets: [],
        libraryTargetIds: [],
        libraryPicks: [],
        scanSensitive: true,
        scanAllDb: false
      },
      rules: {
        name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
        dbType: [{ required: true, message: '请选择数据库类型', trigger: 'change' }]
      }
    }
  },
  computed: {
    targetTotal() {
      return (this.taskForm.targets || []).length + (this.taskForm.libraryTargetIds || []).length
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
      this.resetTaskForm()
    },
    resetTaskForm() {
      this.taskForm = {
        name: '',
        dbType: 1,
        targets: [],
        libraryTargetIds: [],
        libraryPicks: [],
        scanSensitive: true,
        scanAllDb: false
      }
    },
    onPickLibrary(rows) {
      const ids = new Set(this.taskForm.libraryTargetIds || [])
      const picks = (this.taskForm.libraryPicks || []).slice()
      rows.forEach((r) => {
        if (!ids.has(r.id)) {
          ids.add(r.id)
          picks.push(r)
        }
      })
      this.taskForm.libraryTargetIds = Array.from(ids)
      this.taskForm.libraryPicks = picks
    },
    removeLibraryPick(id) {
      this.taskForm.libraryTargetIds = (this.taskForm.libraryTargetIds || []).filter((x) => x !== id)
      this.taskForm.libraryPicks = (this.taskForm.libraryPicks || []).filter((x) => x.id !== id)
    },
    buildDbPayload() {
      return {
        name: (this.taskForm.name || '').trim(),
        dbType: Number(this.taskForm.dbType) || 1,
        scanSensitive: !!this.taskForm.scanSensitive,
        scanAllDb: !!this.taskForm.scanAllDb,
        libraryTargetIds: this.taskForm.libraryTargetIds || [],
        targets: (this.taskForm.targets || []).map((t) => {
          const port = parseInt(t.dbPort, 10)
          return {
            dbHost: (t.dbHost || '').trim(),
            dbPort: Number.isFinite(port) ? port : 0,
            dbName: (t.dbName || '').trim(),
            dbUser: (t.dbUser || '').trim(),
            dbPassword: t.dbPassword || ''
          }
        })
      }
    },
    buildTargetPayload(t) {
      const port = parseInt(t.dbPort, 10)
      return {
        dbType: Number(this.taskForm.dbType) || 1,
        dbHost: (t.dbHost || '').trim(),
        dbPort: Number.isFinite(port) ? port : 0,
        dbName: (t.dbName || '').trim(),
        dbUser: (t.dbUser || '').trim(),
        dbPassword: t.dbPassword || ''
      }
    },
    async testAllConnections() {
      const targets = this.taskForm.targets || []
      const libIds = this.taskForm.libraryTargetIds || []
      if (!targets.length && !libIds.length) {
        this.$message({ message: '请至少添加一个目标', type: 'warning' })
        return
      }
      this.testConnLoading = true
      let ok = 0
      let fail = 0
      try {
        if (libIds.length) {
          const res = await security.batchTestDatasecTargetConn({ ids: libIds })
          if (res.code === 200 && res.data) {
            ok += res.data.ok || 0
            fail += res.data.fail || 0
          } else {
            fail += libIds.length
          }
        }
        for (let i = 0; i < targets.length; i++) {
          const t = targets[i]
          if (!t.dbHost || !t.dbUser) {
            fail++
            continue
          }
          const res = await security.testDataSecDBConn(this.buildTargetPayload(t))
          if (res.code === 200 && res.data && res.data.ok) ok++
          else fail++
        }
        const msg = fail === 0
          ? `全部 ${ok} 个目标连接成功`
          : `成功 ${ok} 个，失败 ${fail} 个（请检查地址与凭据）`
        this.$message({ message: msg, type: fail === 0 ? 'success' : 'warning', duration: 5000 })
      } catch (e) {
        this.$message({ message: '全部连接测试失败: ' + (e.message || ''), type: 'error' })
      } finally {
        this.testConnLoading = false
      }
    },
    async submitForm() {
      this.$refs.taskForm.validate(async (valid) => {
        if (!valid) return
        const targets = this.taskForm.targets || []
        const libCount = (this.taskForm.libraryTargetIds || []).length
        if (!targets.length && !libCount) {
          this.$message({ message: '请至少添加一个目标', type: 'warning' })
          return
        }
        for (let i = 0; i < targets.length; i++) {
          const t = targets[i]
          if (!(t.dbHost || '').trim() || !(t.dbUser || '').trim() || !t.dbPassword) {
            this.$message({ message: `请完善第 ${i + 1} 个目标的地址、用户名和密码`, type: 'warning' })
            return
          }
        }
        const res = await security.runDBCheck(this.buildDbPayload())
          if (res.code == 200) {
            this.$message({ message: '任务创建成功', type: 'success' })
            this.dialogVisible = false
            this.getData()
          } else {
            this.$message({ message: res.msg, type: 'error' })
          }
      })
    },
    async handleRerun(row) {
      const id = row.id || row.ID
      if (!id) return
      try {
        const res = await security.rerunDataSecTask({ id, kind: 'db', name: `${row.name || '任务'}-再次检测` })
        if (res.code === 200) {
          this.$message({ message: '已创建再次检测任务', type: 'success' })
          this.getData()
        } else {
          this.$message({ message: res.msg || '操作失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '操作失败: ' + (e.message || ''), type: 'error' })
      }
    },
    async handleCopyTargets(row) {
      const id = row.id || row.ID
      if (!id) return
      try {
        const res = await security.cloneDataSecTaskTargets({ id, kind: 'db' })
        if (res.code === 200 && res.data) {
          this.dialogVisible = true
          this.taskForm = {
            name: `${row.name || '任务'}-复制`,
            dbType: res.data.dbType || row.dbType || 1,
            targets: res.data.targets || [],
            libraryTargetIds: [],
            libraryPicks: [],
            scanSensitive: res.data.scanSensitive !== false,
            scanAllDb: !!res.data.scanAllDb
          }
          this.$message({ message: '已复制历史任务目标', type: 'success' })
        } else {
          this.$message({ message: res.msg || '复制失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '复制失败: ' + (e.message || ''), type: 'error' })
      }
    },
    async saveTargetsToLibrary() {
      const targets = this.taskForm.targets || []
      if (!targets.length) {
        this.$message({ message: '请先手动添加目标（目标库条目无需重复保存）', type: 'warning' })
        return
      }
      try {
        const res = await security.saveDatasecTargetsFromTask({
          dbType: this.taskForm.dbType,
          groupName: '',
          targets: targets.map((t) => this.buildTargetPayload(t))
        })
        if (res.code === 200) {
          this.$message({ message: `已保存 ${(res.data && res.data.saved) || 0} 个目标到目标库`, type: 'success' })
        } else {
          this.$message({ message: res.msg || '保存失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '保存失败: ' + (e.message || ''), type: 'error' })
      }
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
    handleDetail(row) {
      const id = row.id || row.ID
      if (!id) {
        this.$message({ message: '无效的任务', type: 'warning' })
        return
      }
      this.$router.push({
        path: '/datasec/task/detail',
        query: { id, kind: 'db', from: 'db' }
      })
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
      const map = {
        1: '身份认证', 2: '权限控制', 3: '配置安全', 4: '审计日志',
        5: '网络安全', 6: '加密', 7: 'SQL 注入', 8: '敏感数据识别'
      }
      return map[category] || '未知'
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

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

.text-muted { color: #64748b; }

</style>