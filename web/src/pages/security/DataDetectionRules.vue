<template>
  <div class="security-container">
    <div class="main-title">数据安全检查 · 检测规则</div>
    <div class="list_box rules-placeholder">
      <p class="lead">对应需求「数据安全检查」中的规则库与分类管理指标。</p>
      <ul>
        <li>规划内置数据安全检测规则 <strong>3800+</strong> 条（身份认证、权限、配置安全、SQL 注入、敏感数据识别等）。</li>
        <li>数据库基线核查的<strong>运行时规则</strong>随任务执行；本页将用于展示「已启用规则集 / 分类 / 版本」及与 MySQL、PostgreSQL、MongoDB、Redis、CouchDB 等类型的映射。</li>
        <li>任务入口见 <router-link to="/datasec/tasks">任务管理</router-link>（数据库安全检查、敏感数据发现）。</li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DataDetectionRules',
  mounted() {
    this.$store.state.activefirstMenu = '/datasec/rules'
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

.rules-placeholder {
  padding: 8px 4px 24px;
  color: #94a3b8;
  font-size: 14px;
  line-height: 1.7;
  max-width: 880px;
}

.lead {
  color: #cbd5e1;
  margin-bottom: 16px;
}

ul {
  padding-left: 20px;
  margin: 0;
}

a {
  color: #00d4aa;
}
</style>
