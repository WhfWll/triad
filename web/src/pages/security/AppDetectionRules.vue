<template>
  <div class="security-container">
    <p class="page-intro">
      应用安全检测规则统一从漏洞库读取，支持按分类、类型、风险、状态筛选，并可导入规则包。
    </p>

    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <span class="db-stat">共 <strong>{{ totalRecords }}</strong> 条应用检测规则</span>
          <el-button type="primary" size="small" icon="el-icon-upload2" @click="importVisible = true">导入规则包</el-button>
        </div>
        <div class="serach-condition">
          <div class="search-text">
            <el-input
              v-model="filters.libName"
              placeholder="搜索规则名称 / CVE / CNVD"
              class="input-with-select"
              size="small"
              clearable
              @keydown.enter.native="handleSearch"
            />
            <el-button type="primary" size="small" @click="handleSearch">搜索</el-button>
          </div>
          <div class="filter-line">
            <el-select v-model="filters.libClass" clearable size="small" placeholder="规则分类">
              <el-option v-for="item in enums.class" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-select v-model="filters.libType" clearable size="small" placeholder="漏洞类型">
              <el-option v-for="item in enums.type" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-select v-model="filters.libRisk" clearable size="small" placeholder="风险等级">
              <el-option v-for="item in enums.risk" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-select v-model="filters.status" clearable size="small" placeholder="规则状态">
              <el-option v-for="item in enums.status" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-button size="small" @click="handleReset">重置</el-button>
          </div>
        </div>
      </div>

      <el-table v-loading="tableLoading" :data="tableData" style="width: 100%" class="myTable">
        <el-table-column prop="name" label="规则名称" min-width="240" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="showDetail(scope.row)">{{ scope.row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="classEnum" label="规则分类" width="120" />
        <el-table-column prop="typeEnum" label="漏洞类型" width="140" />
        <el-table-column prop="riskName" label="风险等级" width="100">
          <template slot-scope="scope">
            <span :class="riskClass(scope.row.risk)">{{ scope.row.riskName || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="vulNum" label="漏洞编号" min-width="160" :show-overflow-tooltip="true" />
        <el-table-column prop="verifyType" label="验证方式" width="110" />
        <el-table-column prop="statusEnum" label="状态" width="90" />
        <el-table-column label="操作" width="140" fixed="right">
          <template slot-scope="scope">
            <el-button type="text" size="small" @click="showDetail(scope.row)">详情</el-button>
            <el-button type="text" size="small" @click="toggleStatus(scope.row)">
              {{ Number(scope.row.status) === 1 ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :page-size="pageSize"
        background
        layout="total, prev, pager, next, sizes, jumper"
        :total="totalRecords"
        :current-page="currentPage"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>

    <el-dialog title="规则详情" :visible.sync="detailVisible" width="780px" custom-class="theme-dialog" top="5vh">
      <div v-if="currentRule" class="rule-detail">
        <h3>{{ currentRule.name }}</h3>
        <p class="rule-meta">
          <span :class="riskClass(currentRule.risk)">{{ currentRule.riskName || '-' }}</span>
          <span v-if="currentRule.vulNum" class="rule-tag">{{ currentRule.vulNum }}</span>
          <span v-if="currentRule.verifyType" class="rule-tag">{{ currentRule.verifyType }}</span>
        </p>
        <el-descriptions :column="2" border class="rule-info">
          <el-descriptions-item label="规则分类">{{ currentRule.classEnum || '-' }}</el-descriptions-item>
          <el-descriptions-item label="漏洞类型">{{ currentRule.typeEnum || '-' }}</el-descriptions-item>
          <el-descriptions-item label="利用影响">{{ currentRule.exploitImpactEnum || '-' }}</el-descriptions-item>
          <el-descriptions-item label="CVSS">{{ currentRule.cvss || '-' }}</el-descriptions-item>
          <el-descriptions-item label="披露时间">{{ currentRule.publishTime || '-' }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ currentRule.statusEnum || '-' }}</el-descriptions-item>
        </el-descriptions>
        <div class="detail-section">
          <h4>规则描述</h4>
          <p>{{ currentRule.description || '暂无描述' }}</p>
        </div>
        <div class="detail-section">
          <h4>影响范围</h4>
          <p>{{ currentRule.affectRange || '暂无' }}</p>
        </div>
        <div class="detail-section">
          <h4>修复建议</h4>
          <p>{{ currentRule.fixSuggest || '暂无' }}</p>
        </div>
      </div>
      <span slot="footer">
        <el-button @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>

    <el-dialog title="导入应用检测规则包" :visible.sync="importVisible" width="500px" custom-class="theme-dialog" @closed="resetImport">
      <div class="custom-upload" @click="triggerFileInput" @drop.prevent="handleDrop" @dragover.prevent>
        <input ref="fileInput" type="file" accept=".zip,.json,.yak" style="display:none" @change="onNativeFileChange">
        <div class="custom-upload-dragger" :class="{ 'is-dragover': dragOver }">
          <i class="el-icon-upload"></i>
          <div class="el-upload__text" v-if="!selectedFile">将规则包拖到此处，或<em>点击选择</em></div>
          <div class="el-upload__text" v-else>已选择文件：<em>{{ selectedFile.name }}</em></div>
        </div>
        <div class="el-upload__tip">支持 .zip、.json、.yak 格式</div>
      </div>
      <span slot="footer">
        <el-button @click="importVisible = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="submitImport">开始导入</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { vulnerability } from '@/api/tool.js'

export default {
  name: 'AppDetectionRules',
  data() {
    return {
      filters: {
        libName: '',
        libClass: '',
        libType: '',
        libRisk: '',
        status: ''
      },
      enums: {
        class: [],
        type: [],
        risk: [],
        status: []
      },
      currentPage: 1,
      pageSize: 20,
      totalRecords: 0,
      tableData: [],
      tableLoading: false,
      detailVisible: false,
      currentRule: null,
      importVisible: false,
      importing: false,
      selectedFile: null,
      dragOver: false
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/rules'
    this.loadEnums()
    this.loadData()
  },
  methods: {
    async loadEnums() {
      try {
        const res = await vulnerability.getVulObjectlist()
        if (res.code === 200 && res.data) {
          this.enums = {
            class: res.data.class || [],
            type: res.data.type || [],
            risk: res.data.risk || [],
            status: res.data.status || []
          }
        }
      } catch (e) {}
    },
    async loadData() {
      this.tableLoading = true
      try {
        const params = {
          page: this.currentPage,
          size: this.pageSize,
          libName: this.filters.libName,
          libClass: this.filters.libClass,
          libType: this.filters.libType,
          libRisk: this.filters.libRisk,
          status: this.filters.status
        }
        const res = await vulnerability.getObjectData(params)
        if (res.code === 200) {
          this.tableData = (res.data && res.data.list) || []
          this.totalRecords = (res.data && res.data.total) || 0
        } else {
          this.$message({ message: res.msg || '查询失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '请求异常', type: 'error' })
      } finally {
        this.tableLoading = false
      }
    },
    handleSearch() {
      this.currentPage = 1
      this.loadData()
    },
    handleReset() {
      this.filters = { libName: '', libClass: '', libType: '', libRisk: '', status: '' }
      this.currentPage = 1
      this.loadData()
    },
    handleCurrentChange(page) {
      this.currentPage = page
      this.loadData()
    },
    handleSizeChange(size) {
      this.pageSize = size
      this.currentPage = 1
      this.loadData()
    },
    showDetail(row) {
      this.currentRule = row
      this.detailVisible = true
    },
    async toggleStatus(row) {
      const nextStatus = Number(row.status) === 1 ? 2 : 1
      const action = nextStatus === 1 ? '启用' : '禁用'
      try {
        await this.$confirm(`确认${action}该规则？`, '提示', { type: 'warning' })
        const res = await vulnerability.openVul({ libId: row.id, status: nextStatus })
        if (res.code === 200) {
          this.$message({ message: `${action}成功`, type: 'success' })
          this.loadData()
        } else {
          this.$message({ message: res.msg || `${action}失败`, type: 'error' })
        }
      } catch (e) {}
    },
    riskClass(risk) {
      if (risk === 0) return 'risk-critical'
      if (risk === 1) return 'risk-high'
      if (risk === 2) return 'risk-medium'
      if (risk === 3) return 'risk-low'
      return 'risk-info'
    },
    triggerFileInput() {
      this.$refs.fileInput && this.$refs.fileInput.click()
    },
    isImportFile(name) {
      const lower = (name || '').toLowerCase()
      return lower.endsWith('.zip') || lower.endsWith('.json') || lower.endsWith('.yak')
    },
    onNativeFileChange(e) {
      const file = e.target.files && e.target.files[0]
      if (!file) return
      if (!this.isImportFile(file.name)) {
        this.$message({ message: '仅支持 .zip、.json、.yak 格式', type: 'warning' })
        return
      }
      this.selectedFile = file
    },
    handleDrop(e) {
      this.dragOver = false
      const file = e.dataTransfer.files && e.dataTransfer.files[0]
      if (!file) return
      if (!this.isImportFile(file.name)) {
        this.$message({ message: '仅支持 .zip、.json、.yak 格式', type: 'warning' })
        return
      }
      this.selectedFile = file
    },
    resetImport() {
      this.selectedFile = null
      this.dragOver = false
      this.importing = false
      if (this.$refs.fileInput) this.$refs.fileInput.value = ''
    },
    async submitImport() {
      if (!this.selectedFile) {
        this.$message({ message: '请先选择要导入的规则包文件', type: 'warning' })
        return
      }
      this.importing = true
      try {
        const formData = new FormData()
        formData.append('file', this.selectedFile)
        const res = await vulnerability.importVulnVulKit(formData)
        const data = res.data || res
        if (data.code === 200) {
          this.$message({ message: '导入完成', type: 'success' })
          this.importVisible = false
          this.loadData()
        } else {
          this.$message({ message: data.msg || '导入失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '导入请求异常', type: 'error' })
      } finally {
        this.importing = false
      }
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
}

.operationbutton {
  display: flex;
  align-items: center;
  gap: 12px;
}

.db-stat {
  color: #94a3b8;
  font-size: 13px;
}

.filter-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-line .el-select {
  width: 132px;
}

.risk-critical,
.risk-high,
.risk-medium,
.risk-low,
.risk-info {
  display: inline-block;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 3px;
}

.risk-critical { color: #ff4d4f; background: rgba(255, 77, 79, 0.14); }
.risk-high { color: #f56c6c; background: rgba(245, 108, 108, 0.12); }
.risk-medium { color: #e6a23c; background: rgba(230, 162, 60, 0.12); }
.risk-low { color: #409eff; background: rgba(64, 158, 255, 0.12); }
.risk-info { color: #909399; background: rgba(144, 147, 153, 0.12); }

.rule-detail {
  color: #e2e8f0;
}

.rule-detail h3 {
  margin: 0 0 8px;
  color: #00d4aa;
  font-size: 18px;
}

.rule-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 18px;
}

.rule-tag {
  display: inline-block;
  font-size: 12px;
  padding: 2px 10px;
  border-radius: 3px;
  background: rgba(0, 212, 170, 0.12);
  color: #00d4aa;
}

.detail-section {
  margin-top: 16px;
}

.detail-section h4 {
  color: #00d4aa;
  margin: 0 0 10px;
  font-size: 15px;
}

.detail-section p {
  color: #cbd5e1;
  line-height: 1.7;
  font-size: 14px;
  margin: 0;
}

.custom-upload {
  cursor: pointer;
}

.custom-upload-dragger {
  background: rgba(0, 0, 0, 0.7);
  border: 2px dashed rgba(255, 255, 255, 0.15);
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  transition: all 0.3s;
}

.custom-upload-dragger:hover,
.custom-upload-dragger.is-dragover {
  border-color: #00d4aa;
}

.custom-upload .el-icon-upload {
  color: #94a3b8;
  font-size: 48px;
}

.custom-upload .el-upload__text {
  color: #94a3b8;
  font-size: 14px;
  margin-top: 12px;
}

.custom-upload .el-upload__text em {
  color: #00d4aa;
  font-style: normal;
}

.custom-upload .el-upload__tip {
  color: #64748b;
  font-size: 12px;
  margin-top: 8px;
  text-align: center;
}

/deep/ .rule-info .el-descriptions__cell,
/deep/ .rule-info .el-descriptions-item__label,
/deep/ .rule-info .el-descriptions-item__content {
  background: rgba(15, 23, 42, 0.9);
  border-color: rgba(255, 255, 255, 0.08);
  color: #94a3b8;
}

/deep/ .rule-info .el-descriptions-item__label {
  background: rgba(30, 41, 59, 0.9);
  color: #cbd5e1;
}
</style>
