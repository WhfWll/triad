<template>
  <div class="security-container mod-hub">
    <div class="main-title">应用安全检查 · 任务管理</div>
    <p class="page-intro">覆盖需求中的运行时安全测试能力入口：动态扫描、专项应用检测（与历史路由兼容，见下方页签）。</p>
    <el-tabs v-model="subTab" class="hub-inner-tabs" @tab-click="onTabClick">
      <el-tab-pane label="动态扫描" name="dyn" />
      <el-tab-pane label="专项应用检测" name="app" />
    </el-tabs>
    <div class="tab-panel">
      <keep-alive>
        <component :is="activeComp" :embedded="true" />
      </keep-alive>
    </div>
  </div>
</template>

<script>
import DynamicScan from './DynamicScan.vue'
import AppSpecificScan from './AppSpecificScan.vue'

export default {
  name: 'AppSecTaskHub',
  components: { DynamicScan, AppSpecificScan },
  data() {
    return {
      subTab: 'dyn'
    }
  },
  computed: {
    activeComp() {
      return this.subTab === 'dyn' ? 'DynamicScan' : 'AppSpecificScan'
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
    const t = this.$route.query.tab
    if (t === 'app' || t === 'dyn') {
      this.subTab = t
    }
  },
  watch: {
    '$route.query.tab'(t) {
      if (t === 'app' || t === 'dyn') this.subTab = t
    }
  },
  methods: {
    onTabClick(tab) {
      const name = tab && tab.name ? tab.name : this.subTab
      this.$router.replace({ path: '/appsec/tasks', query: { tab: name } }).catch(() => {})
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
