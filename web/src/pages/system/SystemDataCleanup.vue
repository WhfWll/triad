<!-- 系统配置 - 系统数据初始化 -->
<template>
  <div class="system-data-cleanup">
    <p class="page-intro">
      初始化系统运行数据，清理所有扫描任务、检测结果、安全报告和远程会话等运行时产生的数据。
      检测规则、漏洞库、安全策略等配置数据<strong>不受影响</strong>。
    </p>

    <div class="cleanup-cards">
      <div class="cleanup-card">
        <div class="card-header">
          <i class="iconfont iconrenwu"></i>
          <span class="card-title">将清理的数据</span>
        </div>
        <div class="card-body">
          <ul class="data-list to-clean">
            <li>所有扫描任务及任务目标</li>
            <li>所有检查结果（基线检查 / CVE 漏洞 / 恶意代码）</li>
            <li>所有检测结果、漏洞记录、证据数据</li>
            <li>所有安全报告及报告记录</li>
            <li>所有数据库安全检查结果</li>
            <li>所有敏感数据发现结果</li>
            <li>所有远程会话记录</li>
            <li>所有任务日志和日志详情</li>
          </ul>
        </div>
      </div>

      <div class="cleanup-card">
        <div class="card-header">
          <i class="iconfont iconbaocun"></i>
          <span class="card-title">将保留的数据</span>
        </div>
        <div class="card-body">
          <ul class="data-list keep">
            <li>用户与用户组</li>
            <li>安全配置核查规则库</li>
            <li>CVE 漏洞库</li>
            <li>YARA 恶意代码规则库</li>
            <li>数据安全检测规则库</li>
            <li>系统配置与策略</li>
            <li>资产信息与目标库</li>
            <li>审计日志</li>
          </ul>
        </div>
      </div>
    </div>

    <div class="action-section">
      <el-button
        type="danger"
        size="medium"
        :loading="cleaning"
        :disabled="cleaning"
        @click="handleCleanup"
      >
        <i class="el-icon-delete"></i>
        {{ cleaning ? '正在清理...' : '初始化系统数据' }}
      </el-button>
      <p class="action-hint">
        <i class="el-icon-warning-outline"></i>
        此操作不可撤销，请确认已备份重要数据后再执行
      </p>
    </div>
  </div>
</template>

<script>
import { system } from '@/api/system.js'

export default {
  name: 'SystemDataCleanup',
  data() {
    return {
      cleaning: false,
    }
  },
  methods: {
    handleCleanup() {
      this.$confirm(
        '确定要初始化系统运行数据吗？\n\n' +
        '此操作将永久删除所有扫描任务、检测结果、安全报告等运行数据，且不可恢复。\n\n' +
        '检测规则、漏洞库、安全策略等配置数据将保留不受影响。',
        '系统数据初始化',
        {
          confirmButtonText: '确定执行',
          cancelButtonText: '取消',
          type: 'warning',
          dangerouslyUseHTMLString: false,
          customClass: 'cleanup-confirm-dialog',
        }
      ).then(() => {
        this.executeCleanup()
      }).catch(() => {})
    },
    async executeCleanup() {
      this.cleaning = true
      try {
        const res = await system.cleanupRuntimeData()
        if (res.code === 200) {
          this.$message.success('系统运行数据清理完成')
        } else {
          this.$message.error(res.msg || '清理失败')
        }
      } catch (err) {
        this.$message.error('请求清理接口失败')
        console.error(err)
      } finally {
        this.cleaning = false
      }
    },
  }
}
</script>

<style scoped lang="less">
@import '../security/css/appsec-tokens.less';

.system-data-cleanup {
  height: 100%;
  overflow-y: auto;
  padding: 20px;
  box-sizing: border-box;
}

.page-intro {
  font-size: 14px;
  color: @appsec-text-body;
  line-height: 1.6;
  margin: 0 0 24px;
  padding: 16px 20px;
  background: @appsec-bg-surface;
  border-radius: @appsec-radius-md;
  border: 1px solid @appsec-border-default;
  strong {
    color: @appsec-accent;
    font-weight: 600;
  }
}

.cleanup-cards {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;
}

.cleanup-card {
  flex: 1;
  background: @appsec-bg-surface;
  border: 1px solid @appsec-border-default;
  border-radius: @appsec-radius-md;
  box-shadow: @appsec-shadow-card;
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  padding: 14px 20px;
  border-bottom: 1px solid @appsec-border-subtle;

  i {
    font-size: 16px;
    margin-right: 8px;
  }

  .card-title {
    font-size: 15px;
    font-weight: 600;
    color: @appsec-text-strong;
  }
}

.cleanup-card:first-child .card-header i {
  color: #f56c6c;
}

.cleanup-card:last-child .card-header i {
  color: #67c23a;
}

.card-body {
  padding: 12px 20px 16px;
}

.data-list {
  list-style: none;
  padding: 0;
  margin: 0;

  li {
    position: relative;
    padding: 6px 0 6px 20px;
    font-size: 13px;
    color: @appsec-text-body;
    line-height: 1.5;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 13px;
      width: 6px;
      height: 6px;
      border-radius: 50%;
    }
  }
}

.data-list.to-clean li::before {
  background: #f56c6c;
}

.data-list.keep li::before {
  background: #67c23a;
}

.action-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 0;
  gap: 12px;
}

.action-hint {
  font-size: 12px;
  color: @appsec-text-muted;
  margin: 0;

  i {
    margin-right: 4px;
    color: #e6a23c;
  }
}
</style>
