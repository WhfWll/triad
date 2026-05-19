<template>
  <div class="task-steps">
    <div v-for="(s, i) in steps" :key="s.key" class="step-item" :class="{ active: current === s.key, done: stepIndex(s.key) < stepIndex(current) }">
      <span class="step-num">{{ i + 1 }}</span>
      <span class="step-label">{{ s.label }}</span>
      <span v-if="i < steps.length - 1" class="step-line" />
    </div>
  </div>
</template>

<script>
export default {
  name: 'AppScanTaskStepBar',
  props: {
    current: { type: String, required: true },
    showPlugins: { type: Boolean, default: true }
  },
  computed: {
    steps() {
      const list = [
        { key: 'strategy', label: '选择策略' },
        { key: 'configure', label: '任务与扫描配置' }
      ]
      if (this.showPlugins) {
        list.push({ key: 'plugins', label: '漏洞脚本选择' })
      }
      return list
    }
  },
  methods: {
    stepIndex(key) {
      return this.steps.findIndex(s => s.key === key)
    }
  }
}
</script>

<style lang="less" scoped>
.task-steps {
  display: flex;
  align-items: center;
  gap: 0;
  margin-bottom: 20px;
  flex-wrap: wrap;
}
.step-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #64748b;
  font-size: 13px;
}
.step-item.active {
  color: #00d4aa;
}
.step-item.done {
  color: #94a3b8;
}
.step-num {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 1px solid currentColor;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}
.step-item.active .step-num {
  background: rgba(0, 212, 170, 0.15);
  border-color: #00d4aa;
}
.step-line {
  width: 32px;
  height: 1px;
  background: rgba(255, 255, 255, 0.12);
  margin: 0 12px;
}
</style>
