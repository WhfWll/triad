<template>
  <div class="security-container">
    <div class="main-title" v-if="!embedded">鏁忔劅鏁版嵁鍙戠幇</div>
    
    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <el-tooltip :content="$store.state.systemAuthorized ? '' : '系统未授权，请前往「系统配置 → 系统授权」页面完成授权'" :disabled="$store.state.systemAuthorized" placement="bottom">
            <el-button type="primary" size="small" :disabled="!$store.state.systemAuthorized" @click="btnCreate">鏂板缓鎵弿浠诲姟</el-button>
          </el-tooltip>
          <router-link to="/datasec/targets"><el-button size="small">鐩爣搴撶鐞?/el-button></router-link>
        </div>
        <div class="serach-condition">
          <div class="search-text">
            <el-input placeholder="鎼滅储浠诲姟鍚嶇О" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select" size="small" clearable></el-input>
            <el-button type="primary" size="small" @click="handlesearch">鎼滅储</el-button>
          </div>
          <div>
            <el-button type="primary" size="small" @click="handleReset">閲嶇疆</el-button>
          </div>
        </div>
      </div>

      <el-table :data="tableData" style="width: 100%" class="myTable" @selection-change="handleSelectionChange">
        <el-table-column width="55" type="selection">
        </el-table-column>
        <el-table-column prop="name" label="浠诲姟鍚嶇О" :show-overflow-tooltip="true">
        </el-table-column>
        <el-table-column prop="dbType" label="鏁版嵁搴撶被鍨?>
          <template slot-scope="scope">
            <span :class="getDBTypeClass(scope.row.dbType)">{{ getDBTypeName(scope.row.dbType) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="dbHost" label="鎵弿鐩爣" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <span>{{ scope.row.targetSummary || scope.row.dbHost }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="totalCount" label="鍙戠幇鏁版嵁鏉℃暟">
        </el-table-column>
        <el-table-column prop="highCount" label="楂樻晱鎰熸暟鎹?>
          <template slot-scope="scope">
            <span class="count-high">{{ scope.row.highCount || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="鐘舵€?>
          <template slot-scope="scope">
            <span :class="getStatusClass(scope.row.status)">{{ getStatusName(scope.row.status) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="鍒涘缓鏃堕棿">
        </el-table-column>
        <el-table-column prop="scanTime" label="鎵弿鏃堕棿">
        </el-table-column>
        <el-table-column label="鎿嶄綔" width="220">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="handleDetail(scope.row)">璇︽儏</el-link>
            <el-link :underline="false" class="link_primary" @click="handleRerun(scope.row)">鍐嶆妫€娴?/el-link>
            <el-link :underline="false" class="link_primary" @click="handleCopyTargets(scope.row)">澶嶅埗鐩爣</el-link>
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

    <el-dialog title="鏂板缓鏁忔劅鏁版嵁鎵弿浠诲姟" :visible.sync="dialogVisible" width="720px">
      <el-form :model="taskForm" :rules="rules" ref="taskForm" label-width="100px">
        <el-form-item label="浠诲姟鍚嶇О" prop="name">
          <el-input v-model="taskForm.name" placeholder="璇疯緭鍏ヤ换鍔″悕绉?></el-input>
        </el-form-item>
        <el-form-item label="鏁版嵁搴撶被鍨? prop="dbType">
          <el-select v-model="taskForm.dbType" placeholder="璇烽€夋嫨鏁版嵁搴撶被鍨?>
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
        <el-form-item label="鏁忔劅鏁版嵁绫诲瀷" prop="dataTypes">
          <el-select v-model="taskForm.dataTypes" multiple placeholder="璇烽€夋嫨鏁忔劅鏁版嵁绫诲瀷">
            <el-option label="韬唤璇佸彿" :value="1"></el-option>
            <el-option label="閾惰鍗″彿" :value="2"></el-option>
            <el-option label="鎶ょ収鍙? :value="3"></el-option>
            <el-option label="鎵嬫満鍙? :value="4"></el-option>
            <el-option label="閭" :value="5"></el-option>
            <el-option label="鍦板潃" :value="6"></el-option>
            <el-option label="鍑虹敓鏃ユ湡" :value="7"></el-option>
            <el-option label="濮撳悕" :value="8"></el-option>
            <el-option label="Token" :value="9"></el-option>
            <el-option label="璇佷功淇℃伅" :value="10"></el-option>
            <el-option label="瀵嗙爜鍝堝笇" :value="11"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">鍙栨秷</el-button>
        <el-button @click="saveTargetsToLibrary">淇濆瓨鍒扮洰鏍囧簱</el-button>
        <el-button :loading="testConnLoading" @click="testAllConnections">娴嬭瘯鍏ㄩ儴杩炴帴</el-button>
        <el-button type="primary" @click="submitForm">寮€濮嬫壂鎻忥紙{{ targetTotal }} 涓洰鏍囷級</el-button>
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
        name: [{ required: true, message: '璇疯緭鍏ヤ换鍔″悕绉?, trigger: 'blur' }],
        dbType: [{ required: true, message: '璇烽€夋嫨鏁版嵁搴撶被鍨?, trigger: 'change' }]
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
        this.$message({ message: '璇疯嚦灏戞坊鍔犱竴涓洰鏍?, type: 'warning' })
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
          ? `鍏ㄩ儴 ${ok} 涓洰鏍囪繛鎺ユ垚鍔焋
          : `鎴愬姛 ${ok} 涓紝澶辫触 ${fail} 涓紙璇锋鏌ュ湴鍧€涓庡嚟鎹級`
        this.$message({ message: msg, type: fail === 0 ? 'success' : 'warning', duration: 5000 })
      } catch (e) {
        this.$message({ message: '鍏ㄩ儴杩炴帴娴嬭瘯澶辫触: ' + (e.message || ''), type: 'error' })
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
          this.$message({ message: '璇疯嚦灏戞坊鍔犱竴涓洰鏍?, type: 'warning' })
          return
        }
        for (let i = 0; i < targets.length; i++) {
          const t = targets[i]
          if (!(t.dbHost || '').trim() || !(t.dbUser || '').trim() || !t.dbPassword) {
            this.$message({ message: `璇峰畬鍠勭 ${i + 1} 涓洰鏍囩殑鍦板潃銆佺敤鎴峰悕鍜屽瘑鐮乣, type: 'warning' })
            return
          }
        }
        const res = await security.runSensitiveScan(this.buildDbPayload())
          if (res.code == 200) {
            this.$message({ message: '浠诲姟鍒涘缓鎴愬姛', type: 'success' })
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
        const res = await security.rerunDataSecTask({ id, kind: 'sensitive', name: `${row.name || '浠诲姟'}-鍐嶆妫€娴媊 })
        if (res.code === 200) {
          this.$message({ message: '宸插垱寤哄啀娆℃娴嬩换鍔?, type: 'success' })
          this.getData()
        } else {
          this.$message({ message: res.msg || '鎿嶄綔澶辫触', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '鎿嶄綔澶辫触: ' + (e.message || ''), type: 'error' })
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
            name: `${row.name || '浠诲姟'}-澶嶅埗`,
            dbType: res.data.dbType || row.dbType || 1,
            targets: res.data.targets || [],
            libraryTargetIds: [],
            libraryPicks: [],
            dataTypes: res.data.dataTypes || []
          }
          this.$message({ message: '宸插鍒跺巻鍙蹭换鍔＄洰鏍?, type: 'success' })
        } else {
          this.$message({ message: res.msg || '澶嶅埗澶辫触', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '澶嶅埗澶辫触: ' + (e.message || ''), type: 'error' })
      }
    },
    async saveTargetsToLibrary() {
      const targets = this.taskForm.targets || []
      if (!targets.length) {
        this.$message({ message: '璇峰厛鎵嬪姩娣诲姞鐩爣锛堢洰鏍囧簱鏉＄洰鏃犻渶閲嶅淇濆瓨锛?, type: 'warning' })
        return
      }
      try {
        const res = await security.saveDatasecTargetsFromTask({
          dbType: this.taskForm.dbType,
          groupName: '',
          targets: targets.map((t) => this.buildTargetPayload(t))
        })
        if (res.code === 200) {
          this.$message({ message: `宸蹭繚瀛?${(res.data && res.data.saved) || 0} 涓洰鏍囧埌鐩爣搴揱, type: 'success' })
        } else {
          this.$message({ message: res.msg || '淇濆瓨澶辫触', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '淇濆瓨澶辫触: ' + (e.message || ''), type: 'error' })
      }
    },
    async handleDel(row) {
      this.$confirm('纭鍒犻櫎璇ヤ换鍔★紵', '鎻愮ず', {
        confirmButtonText: '纭畾',
        cancelButtonText: '鍙栨秷',
        type: 'warning'
      }).then(() => {
        this.$message({ message: '鍒犻櫎鎴愬姛', type: 'success' })
        this.getData()
      }).catch(() => {})
    },
    handleDetail(row) {
      const id = row.id || row.ID
      if (!id) {
        this.$message({ message: '鏃犳晥鐨勪换鍔?, type: 'warning' })
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
      return map[type] || '鏈煡'
    },
    getDBTypeClass(type) {
      const map = { 1: 'db-mysql', 2: 'db-postgresql', 3: 'db-mongodb', 4: 'db-redis', 5: 'db-couchdb' }
      return map[type] || 'db-default'
    },
    getStatusName(status) {
      const map = { 1: '绛夊緟鎵弿', 2: '鎵弿涓?, 3: '宸插畬鎴? }
      return map[status] || '鏈煡'
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
