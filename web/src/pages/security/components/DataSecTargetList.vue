<template>
  <div class="datasec-target-list">
    <div class="targets-header">
      <span class="targets-title">目标列表（{{ totalCount }} 个）</span>
      <div class="targets-actions">
        <el-button type="text" size="small" @click="$emit('pick-library')">从目标库选择</el-button>
        <el-button type="text" size="small" icon="el-icon-plus" @click="openAddDialog">手动添加</el-button>
      </div>
    </div>
    <p class="targets-hint">同一任务可扫描多个数据库实例；类型与敏感规则（如有）对所有目标生效。</p>

    <div v-if="libraryPicks.length" class="library-picks">
      <div class="library-title">目标库（{{ libraryPicks.length }}）</div>
      <div v-for="item in libraryPicks" :key="'lib-' + item.id" class="library-row">
        <span class="lib-tag">库</span>
        <span class="lib-name" :title="item.name">{{ item.name || item.dbHost }}</span>
        <span class="lib-addr">{{ item.dbHost }}:{{ item.dbPort || defaultPort }}</span>
        <el-button type="text" size="mini" class="btn-remove" @click="$emit('remove-library', item.id)">移除</el-button>
      </div>
    </div>

    <div v-if="!innerTargets.length && !libraryPicks.length" class="targets-empty">暂无目标，请从目标库选择或手动添加</div>
    <div v-else-if="innerTargets.length" class="targets-table">
      <div class="targets-table-header">
        <span class="th-index">#</span>
        <span class="th-host">地址</span>
        <span class="th-port">端口</span>
        <span class="th-db">库名</span>
        <span class="th-user">用户</span>
        <span class="th-pwd">密码</span>
        <span class="th-actions">操作</span>
      </div>
      <div v-for="(t, idx) in innerTargets" :key="idx" class="targets-table-row">
        <span class="td-index">{{ idx + 1 }}</span>
        <span class="td-host" :title="t.dbHost">{{ t.dbHost }}</span>
        <span class="td-port">{{ t.dbPort || defaultPort }}</span>
        <span class="td-db" :title="t.dbName">{{ t.dbName || '-' }}</span>
        <span class="td-user" :title="t.dbUser">{{ t.dbUser || '-' }}</span>
        <span class="td-pwd">{{ maskPassword(t.dbPassword) }}</span>
        <span class="td-actions">
          <el-button type="text" size="mini" @click="openEditDialog(idx)">编辑</el-button>
          <el-button type="text" size="mini" class="btn-remove" @click="removeTarget(idx)">删除</el-button>
        </span>
      </div>
    </div>

    <el-dialog
      :title="editIndex >= 0 ? '编辑目标' : '添加目标'"
      :visible.sync="dialogVisible"
      width="520px"
      append-to-body
      @closed="onDialogClosed"
    >
      <el-form ref="targetForm" :model="editForm" :rules="formRules" label-width="72px" size="small">
        <el-form-item label="地址" prop="dbHost">
          <el-input v-model="editForm.dbHost" placeholder="IP 或主机名" />
        </el-form-item>
        <el-form-item label="端口" prop="dbPort">
          <el-input v-model.number="editForm.dbPort" type="number" :placeholder="String(defaultPort)" />
        </el-form-item>
        <el-form-item label="库名" prop="dbName">
          <el-input v-model="editForm.dbName" placeholder="可选" />
        </el-form-item>
        <el-form-item :label="isRedis ? '用户（可选）' : '用户'" prop="dbUser">
          <el-input
            v-model="editForm.dbUser"
            :placeholder="isRedis ? 'Redis 通常无需用户名' : '用户名'"
            autocomplete="off"
          />
        </el-form-item>
        <el-form-item label="密码" prop="dbPassword">
          <el-input v-model="editForm.dbPassword" type="password" placeholder="密码" show-password autocomplete="new-password" />
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button size="small" @click="dialogVisible = false">取消</el-button>
        <el-button size="small" :loading="dialogTestLoading" @click="testCurrentTarget">连接测试</el-button>
        <el-button type="primary" size="small" @click="saveTarget">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import security from '@/api/security.js'
import { isRedisDbType } from '../datasecTaskLabels.js'

const DEFAULT_PORTS = { 1: 3306, 2: 5432, 3: 27017, 4: 6379, 5: 5984 }

export function createEmptyDataSecTarget(dbType = 1) {
  return {
    dbHost: '',
    dbPort: DEFAULT_PORTS[dbType] || 3306,
    dbName: '',
    dbUser: '',
    dbPassword: ''
  }
}

export default {
  name: 'DataSecTargetList',
  props: {
    value: {
      type: Array,
      default: () => []
    },
    dbType: {
      type: Number,
      default: 1
    },
    requirePassword: {
      type: Boolean,
      default: true
    },
    libraryPicks: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      dialogVisible: false,
      dialogTestLoading: false,
      editIndex: -1,
      editForm: createEmptyDataSecTarget(1)
    }
  },
  computed: {
    isRedis() {
      return isRedisDbType(this.dbType)
    },
    formRules() {
      return {
        dbHost: [{ required: true, message: '请输入地址', trigger: 'blur' }],
        dbUser: this.isRedis ? [] : [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        dbPassword: this.requirePassword && !this.isRedis
          ? [{ required: true, message: '请输入密码', trigger: 'blur' }]
          : []
      }
    },
    innerTargets: {
      get() {
        return this.value
      },
      set(v) {
        this.$emit('input', v)
      }
    },
    defaultPort() {
      return DEFAULT_PORTS[this.dbType] || 3306
    },
    totalCount() {
      return (this.innerTargets || []).length + (this.libraryPicks || []).length
    }
  },
  methods: {
    maskPassword(pwd) {
      if (!pwd) return '-'
      if (pwd.length <= 2) return '**'
      return pwd.slice(0, 1) + '*'.repeat(Math.min(pwd.length - 1, 6))
    },
    openAddDialog() {
      this.editIndex = -1
      this.editForm = createEmptyDataSecTarget(this.dbType)
      this.dialogVisible = true
      this.$nextTick(() => {
        if (this.$refs.targetForm) this.$refs.targetForm.clearValidate()
      })
    },
    openEditDialog(idx) {
      const t = this.innerTargets[idx]
      if (!t) return
      this.editIndex = idx
      this.editForm = { ...t }
      this.dialogVisible = true
      this.$nextTick(() => {
        if (this.$refs.targetForm) this.$refs.targetForm.clearValidate()
      })
    },
    saveTarget() {
      this.$refs.targetForm.validate((valid) => {
        if (!valid) return
        const row = {
          dbHost: (this.editForm.dbHost || '').trim(),
          dbPort: Number(this.editForm.dbPort) || this.defaultPort,
          dbName: (this.editForm.dbName || '').trim(),
          dbUser: (this.editForm.dbUser || '').trim(),
          dbPassword: this.editForm.dbPassword || ''
        }
        const next = this.innerTargets.slice()
        if (this.editIndex >= 0) {
          next.splice(this.editIndex, 1, row)
        } else {
          next.push(row)
        }
        this.$emit('input', next)
        this.dialogVisible = false
      })
    },
    removeTarget(idx) {
      const next = this.innerTargets.slice()
      next.splice(idx, 1)
      this.$emit('input', next)
    },
    onDialogClosed() {
      this.editIndex = -1
      this.dialogTestLoading = false
      this.editForm = createEmptyDataSecTarget(this.dbType)
    },
    buildTestPayload() {
      const port = parseInt(this.editForm.dbPort, 10)
      return {
        dbType: Number(this.dbType) || 1,
        dbHost: (this.editForm.dbHost || '').trim(),
        dbPort: Number.isFinite(port) && port > 0 ? port : this.defaultPort,
        dbName: (this.editForm.dbName || '').trim(),
        dbUser: (this.editForm.dbUser || '').trim(),
        dbPassword: this.editForm.dbPassword || ''
      }
    },
    async testCurrentTarget() {
      if (!(this.editForm.dbHost || '').trim()) {
        this.$message({ message: '请先填写地址', type: 'warning' })
        return
      }
      if (!this.isRedis && !(this.editForm.dbUser || '').trim()) {
        this.$message({ message: '请先填写用户名', type: 'warning' })
        return
      }
      if (this.requirePassword && !this.isRedis && !this.editForm.dbPassword) {
        this.$message({ message: '请先填写密码', type: 'warning' })
        return
      }
      this.dialogTestLoading = true
      try {
        const res = await security.testDataSecDBConn(this.buildTestPayload())
        if (res.code === 200 && res.data) {
          this.$message({
            message: res.data.message || (res.data.ok ? '连接成功' : '连接失败'),
            type: res.data.ok ? 'success' : 'error',
            duration: 5000
          })
        } else {
          this.$message({ message: res.msg || '连接测试失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '连接测试失败: ' + (e.message || ''), type: 'error' })
      } finally {
        this.dialogTestLoading = false
      }
    }
  }
}
</script>

<style lang="less" scoped>
.datasec-target-list {
  margin-bottom: 8px;
}

.targets-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
  flex-wrap: wrap;
  gap: 8px;
}

.targets-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
}

.library-picks {
  margin-bottom: 10px;
  border: 1px solid rgba(59, 130, 246, 0.25);
  border-radius: 6px;
  padding: 8px 10px;
  background: rgba(59, 130, 246, 0.06);
}

.library-title {
  font-size: 12px;
  color: #94a3b8;
  margin-bottom: 6px;
}

.library-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #e2e8f0;
  padding: 4px 0;
}

.lib-tag {
  background: rgba(59, 130, 246, 0.25);
  color: #60a5fa;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 11px;
}

.lib-name { font-weight: 600; }
.lib-addr { color: #94a3b8; flex: 1; }

.targets-title {
  font-size: 13px;
  color: #cbd5e1;
  font-weight: 600;
}

.targets-hint {
  font-size: 12px;
  color: #64748b;
  margin: 0 0 10px;
}

.targets-empty {
  padding: 14px;
  text-align: center;
  color: #64748b;
  border: 1px dashed rgba(148, 163, 184, 0.35);
  border-radius: 6px;
  font-size: 13px;
}

.targets-table {
  border: 1px solid rgba(0, 212, 170, 0.12);
  border-radius: 6px;
  overflow: hidden;
  font-size: 12px;
}

.targets-table-header,
.targets-table-row {
  display: grid;
  grid-template-columns: 36px 1.4fr 64px 1fr 88px 72px 96px;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
}

.targets-table-header {
  background: rgba(0, 0, 0, 0.25);
  color: #94a3b8;
  font-weight: 600;
}

.targets-table-row {
  border-top: 1px solid rgba(148, 163, 184, 0.12);
  color: #e2e8f0;

  &:hover {
    background: rgba(0, 212, 170, 0.04);
  }
}

.td-host,
.td-db,
.td-user {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.td-pwd {
  color: #94a3b8;
  font-family: monospace;
}

.td-actions {
  white-space: nowrap;
}

.btn-remove {
  color: #f87171 !important;
  padding: 0 4px;
}
</style>
