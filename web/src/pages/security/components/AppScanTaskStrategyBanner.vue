<template>
  <div v-if="meta" class="strategy-banner">
    <span class="banner-icon">{{ meta.icon }}</span>
    <div>
      <div class="banner-name">{{ meta.name }}</div>
      <div class="banner-desc">{{ meta.desc }}</div>
    </div>
    <el-tag size="mini" type="info">{{ isCustom ? '自定义策略' : '内置策略' }}</el-tag>
  </div>
</template>

<script>
export default {
  name: 'AppScanTaskStrategyBanner',
  props: {
    meta: { type: Object, default: null },
    isCustom: { type: Boolean, default: false }
  }
}
</script>

<style lang="less" scoped>
.strategy-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  margin-bottom: 16px;
  background: rgba(0, 212, 170, 0.06);
  border: 1px solid rgba(0, 212, 170, 0.25);
  border-radius: 8px;
}
.banner-icon {
  font-size: 28px;
}
.banner-name {
  color: #e2e8f0;
  font-weight: 600;
  font-size: 15px;
}
.banner-desc {
  color: #64748b;
  font-size: 12px;
  margin-top: 2px;
}
</style>
