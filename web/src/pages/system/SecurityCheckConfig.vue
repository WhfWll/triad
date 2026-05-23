<!-- 系统配置 - 安全检查配置 -->
<template>
  <div class="security-check-config">
    <p class="page-intro">
      配置主机 / 应用 / 数据三大安全检查模块的运行参数。并发数表示<strong>同一任务内</strong>可同时扫描的目标数量上限；
      修改保存后对新发起的扫描任务立即生效（无需重启后端）。
    </p>

    <el-row :gutter="20">
      <el-col :span="12">
        <div class="title_left_line">
          <label>扫描并发配置</label>
          <el-tooltip class="item" effect="dark" placement="right">
            <div slot="content">
              对齐需求文档默认值：主机安全检查 10、应用安全检查 5、数据安全检查 5。<br />
              有效范围 1~50。
            </div>
            <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
          </el-tooltip>
        </div>
        <div class="div_block">
          <el-form ref="formRef" :model="form" label-width="150px" class="sysform">
            <el-form-item label="主机安全检查：" class="syswarnvalue">
              <el-input
                type="number"
                min="1"
                max="50"
                v-model.number="form.hostConcurrent"
                class="innerinput"
                style="width:calc(100% - 190px)"
              />
              <el-tooltip class="item" effect="dark" placement="right">
                <div slot="content">基线核查 / CVE 漏洞 / YARA 恶意代码，默认 10</div>
                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
              </el-tooltip>
            </el-form-item>
            <el-form-item label="应用安全检查：" class="syswarnvalue">
              <el-input
                type="number"
                min="1"
                max="50"
                v-model.number="form.appConcurrent"
                class="innerinput"
                style="width:calc(100% - 190px)"
              />
              <el-tooltip class="item" effect="dark" placement="right">
                <div slot="content">动态扫描 / 专项应用检测，默认 5</div>
                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
              </el-tooltip>
            </el-form-item>
            <el-form-item label="数据安全检查：" class="syswarnvalue">
              <el-input
                type="number"
                min="1"
                max="50"
                v-model.number="form.dataConcurrent"
                class="innerinput"
                style="width:calc(100% - 190px)"
              />
              <el-tooltip class="item" effect="dark" placement="right">
                <div slot="content">数据库基线 / 敏感数据发现，默认 5</div>
                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
              </el-tooltip>
            </el-form-item>
          </el-form>
          <el-button type="primary" class="div_blockbtn" :loading="saving" @click="save">保存设置</el-button>
        </div>
      </el-col>

      <el-col :span="12">
        <div class="title_left_line">
          <label>快捷入口</label>
        </div>
        <div class="div_block quick-links">
          <router-link class="quick-link" to="/hostsec/tasks">主机安全检查 · 任务管理</router-link>
          <router-link class="quick-link" to="/hostsec/rules">主机安全检查 · 检测规则</router-link>
          <router-link class="quick-link" to="/appsec/tasks">应用安全检查 · 任务管理</router-link>
          <router-link class="quick-link" to="/datasec/tasks">数据安全检查 · 任务管理</router-link>
          <router-link class="quick-link" to="/datasec/rules">数据安全检查 · 检测规则</router-link>
          <router-link class="quick-link" to="/datasec/targets">数据安全检查 · 目标库</router-link>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { securityCheckConfig } from '@/api/system.js'

export default {
  name: 'SecurityCheckConfig',
  data() {
    return {
      saving: false,
      form: {
        hostConcurrent: 10,
        appConcurrent: 5,
        dataConcurrent: 5,
      },
    }
  },
  created() {
    this.load()
  },
  methods: {
    async load() {
      try {
        const res = await securityCheckConfig.get()
        if (res.code === 200 && res.data) {
          this.form.hostConcurrent = res.data.hostConcurrent || 10
          this.form.appConcurrent = res.data.appConcurrent || 5
          this.form.dataConcurrent = res.data.dataConcurrent || 5
        } else if (res.msg) {
          this.$message({ message: res.msg, type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '加载配置失败', type: 'error' })
      }
    },
    async save() {
      const host = Number(this.form.hostConcurrent)
      const app = Number(this.form.appConcurrent)
      const data = Number(this.form.dataConcurrent)
      if (!host || host < 1 || host > 50 || !app || app < 1 || app > 50 || !data || data < 1 || data > 50) {
        this.$message({ message: '并发数须为 1~50 的整数', type: 'warning' })
        return
      }
      this.saving = true
      try {
        const res = await securityCheckConfig.save({
          hostConcurrent: host,
          appConcurrent: app,
          dataConcurrent: data,
        })
        if (res.code === 200) {
          this.$message({ message: '保存成功', type: 'success' })
        } else {
          this.$message({ message: res.msg || '保存失败', type: 'error' })
        }
      } finally {
        this.saving = false
      }
    },
  },
}
</script>

<style scoped lang="less">
@import '../security/css/appsec-tokens.less';

.security-check-config {
  height: 100%;
  overflow-x: hidden;
  overflow-y: auto;
  padding: 4px 0 24px;
}

.page-intro {
  margin: 0 0 20px;
  padding: 12px 16px;
  background: @appsec-accent-soft;
  border: 1px solid @appsec-border-default;
  border-radius: @appsec-radius-sm;
  color: @appsec-text-body;
  font-size: @appsec-font-size-sm;
  line-height: 1.6;

  strong {
    color: @appsec-text-primary;
    font-weight: @appsec-font-weight-medium;
  }
}

.title_left_line {
  margin-bottom: 16px;
  label {
    font-size: @appsec-font-size-base;
    font-weight: @appsec-font-weight-medium;
    color: @appsec-text-primary;
    margin-right: 8px;
  }
}

.div_block {
  .appsec-card();
  min-height: 120px;

  .div_blockbtn {
    margin-top: 8px;
  }
}

.quick-links {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quick-link {
  color: @appsec-accent;
  font-size: @appsec-font-size-base;
  text-decoration: none;
  transition: color 0.2s;
  &:hover {
    color: @appsec-accent-hover;
    text-decoration: underline;
  }
}

/deep/ .el-form-item__label {
  text-align: left;
  color: @appsec-text-body;
}
</style>
