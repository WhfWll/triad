<template>
  <div class="security-container nessus-editor">
    <div class="editor-header">
      <div class="breadcrumb">
        <span class="crumb-muted">{{ isDynamic ? '新建扫描' : '新建专项检测' }}</span>
        <span class="crumb-sep">/</span>
        <span class="crumb-current">选择策略</span>
      </div>
      <a class="back-link" @click.prevent="goBack">
        <i class="el-icon-arrow-left" /> 返回任务列表
      </a>
    </div>

    <div class="list_box pick-box">
      <app-sec-strategy-picker
        :strategies="strategies"
        intro="选择扫描策略模板，下一步在配置页填写目标、扫描参数与插件。"
        @select="selectStrategy"
      />
      <div class="pick-footer">
        <router-link to="/appsec/strategy" class="link-muted">管理自定义扫描策略</router-link>
        <el-button size="small" @click="goBack">取消</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import AppSecStrategyPicker from './components/AppSecStrategyPicker.vue'
import { BUILTIN_STRATEGIES } from './appsecBuiltinStrategies.js'

export default {
  name: 'AppScanTaskNew',
  components: { AppSecStrategyPicker },
  computed: {
    isDynamic() {
      return this.$route.query.type !== 'app'
    },
    scanType() {
      return this.isDynamic ? 'dyn' : 'app'
    },
    strategies() {
      return BUILTIN_STRATEGIES.map(s => ({ ...s, builtin: true }))
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
  },
  methods: {
    selectStrategy(s) {
      this.$router.push({
        path: '/appsec/task/configure',
        query: { strategy: s.id, type: this.scanType }
      })
    },
    goBack() {
      this.$router.push({ path: '/appsec/tasks', query: { tab: this.scanType } })
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';
@import './css/appsec-strategy-editor.less';

.pick-box {
  padding: 24px;
  margin-top: 16px;
}

.pick-footer {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.link-muted {
  color: #64748b;
  font-size: 13px;
  text-decoration: none;
}
.link-muted:hover {
  color: #00d4aa;
}
</style>
