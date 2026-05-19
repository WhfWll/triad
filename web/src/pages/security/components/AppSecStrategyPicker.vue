<template>
  <div class="strategy-picker">
    <p v-if="intro" class="picker-intro">{{ intro }}</p>
    <div class="strategy-cards">
      <button
        v-for="s in strategies"
        :key="s.id"
        type="button"
        class="strategy-card"
        :class="{ selected: selectedId === s.id }"
        :style="cardAccentStyle(s.id)"
        @click="$emit('select', s)"
      >
        <div class="strategy-icon-wrap" :style="iconWrapStyle(s.id)">
          <span class="strategy-icon">{{ s.icon }}</span>
        </div>
        <div class="strategy-info">
          <div class="strategy-name-row">
            <span class="strategy-name">{{ s.name }}</span>
            <span v-if="s.builtin" class="mini-badge">内置</span>
          </div>
          <p class="strategy-desc">{{ s.desc }}</p>
        </div>
        <span v-if="showArrow" class="strategy-arrow">→</span>
      </button>
    </div>
  </div>
</template>

<script>
import { strategyCardAccentStyle, strategyIconWrapStyle } from '../appsecStrategyTheme.js'

export default {
  name: 'AppSecStrategyPicker',
  props: {
    strategies: { type: Array, default: () => [] },
    selectedId: { type: String, default: '' },
    intro: { type: String, default: '' },
    showArrow: { type: Boolean, default: true }
  },
  methods: {
    cardAccentStyle: strategyCardAccentStyle,
    iconWrapStyle: strategyIconWrapStyle
  }
}
</script>

<style lang="less" scoped>
.picker-intro {
  color: #94a3b8;
  font-size: 13px;
  margin: 0 0 18px;
  max-width: 720px;
  line-height: 1.5;
}

.strategy-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 14px;
}

.strategy-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  width: 100%;
  padding: 16px 18px;
  text-align: left;
  cursor: pointer;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  background: linear-gradient(165deg, rgba(30, 32, 48, 0.9) 0%, rgba(18, 20, 32, 0.95) 100%);
  transition: transform 0.22s ease, border-color 0.22s ease, box-shadow 0.22s ease;
  border-left: 3px solid var(--card-accent, #00d4aa);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, var(--card-accent), transparent);
    opacity: 0.7;
  }

  &:hover {
    border-color: rgba(255, 255, 255, 0.14);
    box-shadow: 0 10px 28px -12px var(--card-glow);
    transform: translateY(-2px);
  }

  &.selected {
    border-color: var(--card-accent);
    background: rgba(0, 212, 170, 0.06);
    box-shadow: 0 0 0 1px var(--card-accent) inset;
  }
}

.strategy-icon-wrap {
  flex-shrink: 0;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid;
  border-radius: 12px;
}

.strategy-icon {
  font-size: 24px;
  line-height: 1;
}

.strategy-info {
  flex: 1;
  min-width: 0;
}

.strategy-name-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 4px;
}

.strategy-name {
  color: #e2e8f0;
  font-weight: 600;
  font-size: 15px;
}

.mini-badge {
  font-size: 11px;
  padding: 1px 6px;
  color: #5eead4;
  background: rgba(0, 212, 170, 0.12);
  border: 1px solid rgba(0, 212, 170, 0.25);
  border-radius: 4px;
}

.strategy-desc {
  margin: 0;
  color: #64748b;
  font-size: 12px;
  line-height: 1.45;
}

.strategy-arrow {
  color: var(--card-accent, #00d4aa);
  font-size: 18px;
  align-self: center;
  flex-shrink: 0;
}
</style>
