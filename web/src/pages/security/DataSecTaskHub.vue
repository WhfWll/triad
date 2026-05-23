<template>
  <div class="security-container mod-hub">
    <p class="page-intro">新建任务时默认同时执行基线检查与敏感数据发现；详情页可查看配置风险、CVE 与敏感字段结果。</p>
    <div class="tab-panel">
      <d-b-security :embedded="true" />
    </div>
  </div>
</template>

<script>
import DBSecurity from './DBSecurity.vue'

export default {
  name: 'DataSecTaskHub',
  components: { DBSecurity },
  mounted() {
    this.$store.state.activefirstMenu = '/datasec/tasks'
    if (this.$route.query.tab === 'sensitive') {
      this.$router.replace({ path: '/datasec/tasks' }).catch(() => {})
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

.tab-panel {
  margin-top: 4px;
}
</style>
