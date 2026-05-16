<template>
  <div class="security-container mod-hub">
    <div class="main-title">数据安全检查 · 任务管理</div>
    <p class="page-intro">数据库安全检查与敏感数据发现；与历史「数据库安全」「敏感数据发现」菜单路由兼容。</p>
    <el-tabs v-model="subTab" class="hub-inner-tabs" @tab-click="onTabClick">
      <el-tab-pane label="数据库安全检查" name="db" />
      <el-tab-pane label="敏感数据发现" name="sensitive" />
    </el-tabs>
    <div class="tab-panel">
      <keep-alive>
        <component :is="activeComp" :embedded="true" />
      </keep-alive>
    </div>
  </div>
</template>

<script>
import DBSecurity from './DBSecurity.vue'
import SensitiveData from './SensitiveData.vue'

export default {
  name: 'DataSecTaskHub',
  components: { DBSecurity, SensitiveData },
  data() {
    return {
      subTab: 'db'
    }
  },
  computed: {
    activeComp() {
      return this.subTab === 'db' ? 'DBSecurity' : 'SensitiveData'
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/datasec/tasks'
    const t = this.$route.query.tab
    if (t === 'sensitive' || t === 'db') {
      this.subTab = t
    }
  },
  watch: {
    '$route.query.tab'(t) {
      if (t === 'sensitive' || t === 'db') this.subTab = t
    }
  },
  methods: {
    onTabClick(tab) {
      const name = tab && tab.name ? tab.name : this.subTab
      this.$router.replace({ path: '/datasec/tasks', query: { tab: name } }).catch(() => {})
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

.page-intro {
  color: #94a3b8;
  font-size: 13px;
  margin: 0 0 12px;
  max-width: 900px;
}

.hub-inner-tabs {
  margin-bottom: 8px;
}

.tab-panel {
  margin-top: 4px;
}
</style>
