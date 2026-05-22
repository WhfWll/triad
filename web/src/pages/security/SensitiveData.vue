<template>
  <div class="security-container">
    <div class="main-title" v-if="!embedded">敏感数据发现</div>
    
    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="btnCreate">新建扫描任务</el-button>
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
        <el-table-column label="操作" width="220">
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

    <el-dialog title="新建敏感数据扫描任务" :visible.sync="dialogVisible" width="720px">
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
        <el-button @click="saveTargetsToLibrary">保存到目标库</el-button>
        <el-button :loading="testConnLoading" @click="testAllConnections">测试全部连接</el-button>
        <el-button type="primary" @click="submitForm">开始扫描（{{ targetTotal }} 个目标）</el-button>
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
  name: 'SensitiveData',
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
        dataTypes: []
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
      this.resetTaskForm()
    },
    resetTaskForm() {
      this.taskForm = {
        name: '',
        dbType: 1,
        targets: [],
        libraryTargetIds: [],
        libraryPicks: [],
        dataTypes: []
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
        }),
        dataTypes: this.taskForm.dataTypes || [],
        scanAllDb: false
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
        const res = await security.runSensitiveScan(this.buildDbPayload())
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
        const res = await security.rerunDataSecTask({ id, kind: 'sensitive', name: `${row.name || '任务'}-再次检测` })
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
        const res = await security.cloneDataSecTaskTargets({ id, kind: 'sensitive' })
        if (res.code === 200 && res.data) {
          this.dialogVisible = true
          this.taskForm = {
            name: `${row.name || '任务'}-复制`,
            dbType: res.data.dbType || row.dbType || 1,
            targets: res.data.targets || [],
            libraryTargetIds: [],
            libraryPicks: [],
            dataTypes: res.data.dataTypes || []
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
        query: { id, kind: 'sensitive', from: 'sensitive' }
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
    getStatusName(status) {
      const map = { 1: '等待扫描', 2: '扫描中', 3: '已完成' }
      return map[status] || '未知'
    },
    getStatusClass(status) {
      const map = { 1: 'status-wait', 2: 'status-running', 3: 'status-complete' }
      return map[status] || 'status-default'
    },
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

</style>