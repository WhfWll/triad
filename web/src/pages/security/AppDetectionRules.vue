<template>
  <div class="security-container">
    <div class="main-title">应用安全检查 · 检测规则</div>
    <div class="list_box rules-placeholder">
      <p class="lead">对应需求「应用安全检查」中的规则库与分类管理指标。</p>
      <ul>
        <li>规划内置应用漏洞检测规则 <strong>7000+</strong> 条（RCE、未授权、SQL 注入、文件上传、弱口令等），与任务执行引擎解耦，可独立版本化发布。</li>
        <li>当前前端任务能力见 <router-link to="/appsec/tasks">任务管理</router-link>；规则目录页将对接后端「规则包 / 分类 / 版本」接口后在此展示。</li>
        <li>通用漏洞与指纹资产可先在 <router-link to="/vulnerability">工具管理 · 漏洞库</router-link> 维护，便于与应用专项规则联动。</li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AppDetectionRules',
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/rules'
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
