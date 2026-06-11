<template>
  <div class="security-container mod-hub">
    <p class="page-intro">应用安全任务管理当前仅保留动态扫描能力，专项应用检测入口已下线。</p>
    <div class="tab-panel">
      <DynamicScan :embedded="true" />
    </div>
  </div>
</template>

<script>
import DynamicScan from './DynamicScan.vue'

export default {
  name: 'AppSecTaskHub',
  components: { DynamicScan },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
    if (this.$route.query.tab && this.$route.query.tab !== 'dyn') {
      this.$router.replace({ path: '/appsec/tasks' }).catch(() => {})
    }
  },
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

.tab-panel {
  margin-top: 4px;
}
</style>
