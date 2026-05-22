<template>
  <el-dialog
    title="从目标库选择"
    :visible.sync="visibleProxy"
    width="760px"
    append-to-body
    @open="loadList"
  >
    <div class="picker-toolbar">
      <el-input v-model="search" placeholder="搜索名称/地址/分组" size="small" clearable class="search-input" @keydown.enter.native="loadList" />
      <el-button size="small" type="primary" @click="loadList">搜索</el-button>
      <router-link to="/datasec/targets" class="link-manage">管理目标库</router-link>
    </div>
    <el-table
      v-loading="loading"
      :data="list"
      class="myTable"
      max-height="360"
      size="small"
      @selection-change="onSelectionChange"
    >
      <el-table-column type="selection" width="48" />
      <el-table-column prop="name" label="名称" min-width="120" :show-overflow-tooltip="true" />
      <el-table-column prop="groupName" label="分组" width="100" :show-overflow-tooltip="true" />
      <el-table-column prop="dbHost" label="地址" width="130" :show-overflow-tooltip="true" />
      <el-table-column prop="dbPort" label="端口" width="72" />
      <el-table-column prop="dbName" label="库名" width="90" :show-overflow-tooltip="true" />
      <el-table-column prop="dbUser" label="用户" width="90" :show-overflow-tooltip="true" />
    </el-table>
    <p v-if="!loading && !list.length" class="empty-hint">暂无匹配目标，请先在「数据库目标库」中添加</p>
    <span slot="footer">
      <el-button size="small" @click="visibleProxy = false">取消</el-button>
      <el-button type="primary" size="small" @click="confirm">添加所选（{{ selected.length }}）</el-button>
    </span>
  </el-dialog>
</template>

<script>
import security from '@/api/security.js'

export default {
  name: 'DataSecTargetPicker',
  props: {
    visible: { type: Boolean, default: false },
    dbType: { type: Number, default: 1 },
    excludeIds: { type: Array, default: () => [] }
  },
  data() {
    return {
      loading: false,
      search: '',
      list: [],
      selected: []
    }
  },
  computed: {
    visibleProxy: {
      get() { return this.visible },
      set(v) { this.$emit('update:visible', v) }
    }
  },
  methods: {
    async loadList() {
      this.loading = true
      try {
        const res = await security.getDatasecTargetList({
          page: 1,
          size: 200,
          dbType: this.dbType,
          search: this.search || undefined
        })
        const rows = (res.code === 200 && res.data && res.data.list) || []
        const exclude = new Set((this.excludeIds || []).map(Number))
        this.list = rows.filter((r) => !exclude.has(Number(r.id)))
      } catch {
        this.list = []
      } finally {
        this.loading = false
      }
    },
    onSelectionChange(rows) {
      this.selected = rows || []
    },
    confirm() {
      if (!this.selected.length) {
        this.$message({ message: '请至少选择一个目标', type: 'warning' })
        return
      }
      this.$emit('pick', this.selected.map((r) => ({
        id: r.id,
        name: r.name,
        dbHost: r.dbHost,
        dbPort: r.dbPort,
        dbName: r.dbName,
        dbUser: r.dbUser
      })))
      this.visibleProxy = false
      this.selected = []
    }
  }
}
</script>

<style lang="less" scoped>
.picker-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  .search-input { width: 240px; }
}
.link-manage {
  margin-left: auto;
  font-size: 13px;
  color: #00d4aa;
}
.empty-hint {
  text-align: center;
  color: #64748b;
  font-size: 13px;
  padding: 16px 0;
}
</style>
