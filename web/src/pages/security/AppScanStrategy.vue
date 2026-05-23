<template>
  <div class="security-container strategy-page">
    <header class="page-header">
      <div class="page-header-text">
        <p class="page-intro">
          管理扫描策略，定义漏洞脚本、爬虫与端口扫描等配置，供新建任务时一键复用。
        </p>
      </div>
    </header>

    <div class="list_box">
      <div class="toolbar">
        <div class="toolbar-left">
          <span class="stat-pill">
            <i class="el-icon-folder-opened" />
            共 <strong>{{ strategies.length }}</strong> 个策略
          </span>
          <span class="stat-pill muted">
            <i class="el-icon-lock" />
            内置 {{ builtinCount }}
          </span>
        </div>
        <el-button type="primary" size="small" icon="el-icon-plus" @click="createStrategy">
          新建策略
        </el-button>
      </div>

      <div class="strategy-grid">
        <article
          v-for="s in strategies"
          :key="s.id"
          class="strategy-card"
          :class="{ builtin: s.builtin, custom: !s.builtin }"
          :style="cardAccentStyle(s.id)"
        >
          <div class="card-top">
            <div class="card-icon-wrap" :style="iconWrapStyle(s.id)">
              <span class="card-icon-char">{{ s.icon }}</span>
            </div>
            <div class="card-head-main">
              <div class="card-title-row">
                <h3 class="card-name">{{ s.name }}</h3>
                <span v-if="s.builtin" class="badge badge-builtin">内置</span>
                <span v-else class="badge badge-custom">自定义</span>
              </div>
              <p class="card-desc">{{ s.desc }}</p>
            </div>
          </div>

          <div class="card-chips">
            <span class="chip">
              <i class="el-icon-document" />
              {{ vulnLabel(s) }}
            </span>
            <span class="chip" :class="{ on: s.config?.webCrawler?.isOpen }">
              <i class="el-icon-connection" />
              爬虫 {{ s.config?.webCrawler?.isOpen ? '开' : '关' }}
            </span>
            <span class="chip" :class="{ on: s.config?.portScan?.isOpen }">
              <i class="el-icon-position" />
              端口 {{ s.config?.portScan?.isOpen ? '开' : '关' }}
            </span>
          </div>

          <footer class="card-footer">
            <el-button size="small" class="btn-ghost" @click="editStrategy(s)">编辑</el-button>
            <el-button
              v-if="!s.builtin"
              size="small"
              type="text"
              class="btn-delete"
              @click="deleteStrategy(s)"
            >
              删除
            </el-button>
            <el-button size="small" type="primary" class="btn-use" @click="useStrategy(s)">
              使用此策略
            </el-button>
          </footer>
        </article>
      </div>

      <p v-if="!customStrategies.length" class="page-tip">
        <i class="el-icon-info" />
        内置策略可直接编辑保存为副本；也可点击「新建策略」创建完全自定义的扫描策略。
      </p>
    </div>
  </div>
</template>

<script>
import { getBuiltinStrategiesWithConfig } from './appsecBuiltinStrategies.js'
import { strategyCardAccentStyle, strategyIconWrapStyle } from './appsecStrategyTheme.js'

export default {
  name: 'AppScanStrategy',
  data() {
    return {
      builtinStrategies: getBuiltinStrategiesWithConfig(),
      customStrategies: []
    }
  },
  computed: {
    strategies() {
      return [...this.builtinStrategies, ...this.customStrategies]
    },
    builtinCount() {
      return this.builtinStrategies.length
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/tasks'
    this.loadCustomStrategies()
  },
  methods: {
    cardAccentStyle: strategyCardAccentStyle,
    iconWrapStyle: strategyIconWrapStyle,
    vulnLabel(s) {
      const n = s.vulnCount || (s.config?.vulIdsConfig && s.config.vulIdsConfig.length) || 0
      if (n > 0) return `漏洞脚本 ${n}`
      return s.builtin ? '漏洞脚本 · 默认全部' : '漏洞脚本 · 未指定'
    },
    loadCustomStrategies() {
      try {
        const saved = localStorage.getItem('appsec_strategies')
        if (saved) this.customStrategies = JSON.parse(saved)
      } catch {
        this.customStrategies = []
      }
    },
    saveCustomStrategies() {
      localStorage.setItem('appsec_strategies', JSON.stringify(this.customStrategies))
    },
    createStrategy() {
      this.$router.push('/appsec/strategy/new')
    },
    editStrategy(s) {
      this.$router.push({ path: '/appsec/strategy/edit', query: { id: s.id } })
    },
    useStrategy(s) {
      this.$router.push({
        path: '/appsec/task/configure',
        query: {
          strategy: s.id,
          type: 'dyn',
          ...(s.builtin ? {} : { custom: '1' })
        }
      })
    },
    deleteStrategy(s) {
      this.$confirm(`确认删除策略「${s.name}」？`, '提示', { type: 'warning' }).then(() => {
        this.customStrategies = this.customStrategies.filter(x => x.id !== s.id)
        this.saveCustomStrategies()
        this.$message({ message: '已删除', type: 'success' })
      }).catch(() => {})
    }
  }
}
</script>

<style lang="less" scoped>
@import '../bas/css/bas-list-page.less';

.strategy-page {
  .page-header {
    margin-bottom: 20px;
  }

  .page-title {
    margin: 0 0 8px;
    font-size: 20px;
    font-weight: 600;
    color: #e2e8f0;
    letter-spacing: 0.02em;
  }

  .page-intro {
    margin: 0;
    max-width: 720px;
    font-size: 13px;
    line-height: 1.6;
    color: #94a3b8;
  }

  .list_box {
    padding: 20px 24px 28px;
    border: 1px solid rgba(0, 212, 170, 0.08);
  }

  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 22px;
    padding-bottom: 18px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  }

  .toolbar-left {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .stat-pill {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    font-size: 13px;
    color: #cbd5e1;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 20px;

    strong {
      color: #00d4aa;
      font-weight: 600;
    }

    i {
      color: #00d4aa;
      font-size: 14px;
    }

    &.muted {
      color: #94a3b8;

      i {
        color: #64748b;
      }
    }
  }

  .strategy-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 18px;
  }

  .strategy-card {
    position: relative;
    display: flex;
    flex-direction: column;
    min-height: 200px;
    padding: 0;
    overflow: hidden;
    background: linear-gradient(165deg, rgba(30, 32, 48, 0.95) 0%, rgba(18, 20, 32, 0.98) 100%);
    border: 1px solid rgba(255, 255, 255, 0.07);
    border-radius: 12px;
    transition: transform 0.22s ease, border-color 0.22s ease, box-shadow 0.22s ease;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 3px;
      background: linear-gradient(90deg, var(--card-accent), transparent);
      opacity: 0.85;
    }

    &:hover {
      transform: translateY(-3px);
      border-color: rgba(255, 255, 255, 0.12);
      box-shadow: 0 12px 32px -12px var(--card-glow);
    }

    &.builtin::after {
      content: '';
      position: absolute;
      top: 12px;
      right: 12px;
      width: 64px;
      height: 64px;
      border-radius: 50%;
      background: radial-gradient(circle, var(--card-glow) 0%, transparent 70%);
      pointer-events: none;
      opacity: 0.5;
    }
  }

  .card-top {
    display: flex;
    gap: 14px;
    padding: 20px 20px 12px;
  }

  .card-icon-wrap {
    flex-shrink: 0;
    width: 52px;
    height: 52px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid;
    border-radius: 14px;
  }

  .card-icon-char {
    font-size: 26px;
    line-height: 1;
  }

  .card-head-main {
    flex: 1;
    min-width: 0;
  }

  .card-title-row {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 6px;
  }

  .card-name {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #f1f5f9;
    line-height: 1.3;
  }

  .badge {
    padding: 2px 8px;
    font-size: 11px;
    font-weight: 500;
    border-radius: 4px;
    line-height: 1.4;
  }

  .badge-builtin {
    color: #5eead4;
    background: rgba(0, 212, 170, 0.12);
    border: 1px solid rgba(0, 212, 170, 0.25);
  }

  .badge-custom {
    color: #93c5fd;
    background: rgba(59, 130, 246, 0.12);
    border: 1px solid rgba(59, 130, 246, 0.25);
  }

  .card-desc {
    margin: 0;
    font-size: 13px;
    line-height: 1.5;
    color: #94a3b8;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .card-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding: 0 20px 16px;
  }

  .chip {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 4px 10px;
    font-size: 12px;
    color: #64748b;
    background: rgba(0, 0, 0, 0.25);
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 6px;
    transition: color 0.2s, border-color 0.2s, background 0.2s;

    i {
      font-size: 12px;
    }

    &.on {
      color: #a7f3d0;
      background: rgba(0, 212, 170, 0.08);
      border-color: rgba(0, 212, 170, 0.2);
    }
  }

  .card-footer {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: auto;
    padding: 14px 16px;
    background: rgba(0, 0, 0, 0.2);
    border-top: 1px solid rgba(255, 255, 255, 0.05);
  }

  .btn-ghost {
    color: #94a3b8 !important;
    background: transparent !important;
    border-color: rgba(255, 255, 255, 0.12) !important;

    &:hover {
      color: #e2e8f0 !important;
      border-color: rgba(255, 255, 255, 0.22) !important;
      background: rgba(255, 255, 255, 0.04) !important;
    }
  }

  .btn-use {
    margin-left: auto;
    font-weight: 500;
  }

  .btn-delete {
    color: #f87171 !important;
    padding: 0 8px !important;

    &:hover {
      color: #fca5a5 !important;
    }
  }

  .page-tip {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    margin: 22px 0 0;
    padding: 12px 14px;
    font-size: 12px;
    line-height: 1.55;
    color: #64748b;
    background: rgba(0, 212, 170, 0.04);
    border: 1px dashed rgba(0, 212, 170, 0.15);
    border-radius: 8px;

    i {
      margin-top: 2px;
      color: #00d4aa;
    }
  }
}

@media (max-width: 768px) {
  .strategy-page .strategy-grid {
    grid-template-columns: 1fr;
  }
}
</style>
