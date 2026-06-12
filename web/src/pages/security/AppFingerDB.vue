<template>
  <div class="security-container">
    <p class="page-intro">
      SOURCE 指纹知识库，覆盖多种应用类型（Web中间件、CMS、OA、数据库等）的组件指纹数据。
    </p>

    <div class="list_box">
      <div class="search-box">
        <div class="operationbutton">
          <span class="db-stat">共 <strong>{{ totalRecords }}</strong> 条指纹记录</span>
          <el-button type="primary" size="small" @click="handleImport" style="margin-left: 12px;">导入</el-button>
        </div>
        <div class="serach-condition">
          <div class="search-text">
            <el-input
              placeholder="搜索指纹名称"
              @keydown.enter.native="handleSearch"
              v-model="keyword"
              class="input-with-select"
              size="small"
              clearable
            />
            <el-button type="primary" size="small" @click="handleSearch">搜索</el-button>
          </div>
          <div>
            <el-button type="primary" size="small" @click="handleReset">重置</el-button>
          </div>
        </div>
      </div>

      <el-table v-loading="tableLoading" :data="tableData" style="width: 100%" class="myTable">
        <el-table-column prop="fingerName" label="指纹名称" width="400" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <el-link :underline="false" class="link_primary" @click="showDetail(scope.row)">{{ scope.row.fingerName }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="fingerClassEnum" label="指纹分类" width="240" />
        <el-table-column prop="level" label="分层" width="160" />
        <el-table-column prop="createTime" label="创建时间" width="140" />
        <el-table-column label="操作" width="150" fixed="right">
          <template slot-scope="scope">
            <span style="white-space: nowrap;">
              <el-button type="text" size="small" @click="showDetail(scope.row)" style="padding-left: 0; padding-right: 8px;">详情</el-button>
              <el-button type="text" size="small" @click="openTestDialog(scope.row)" style="padding-left: 0; padding-right: 0;">测试</el-button>
            </span>
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

    <el-dialog title="导入指纹" :visible.sync="importVisible" width="480px" custom-class="theme-dialog" top="15vh" :close-on-click-modal="false">
      <div class="import-dialog">
        <el-upload
          ref="uploadRef"
          :auto-upload="false"
          :limit="1"
          :on-change="handleFileChange"
          accept=".zip,.json,.nfp"
          drag
        >
          <i class="el-icon-upload"></i>
          <div class="el-upload__text">将文件拖到此处，或<em>点击选择文件</em></div>
          <div class="el-upload__tip" slot="tip">支持 .zip 压缩包（VulKit目录）、单个 fingerprint.json 或 .nfp（Nmap指纹）</div>
        </el-upload>
        <div v-if="selectedFile" class="import-file-info">
          已选择: <strong>{{ selectedFile.name }}</strong>
        </div>
        <div v-if="importResult" class="import-result">
          <p v-if="importResult.success > 0" class="import-success">
            <i class="el-icon-success"></i> {{ importResult.message || ('成功导入 ' + importResult.success + ' 条记录') }}
          </p>
          <p v-if="importResult.errors && importResult.errors.length" class="import-errors">
            <span v-for="(err, i) in importResult.errors" :key="i">{{ err }}</span>
          </p>
        </div>
      </div>
      <span slot="footer">
        <el-button @click="importVisible = false">关闭</el-button>
        <el-button type="primary" :loading="importLoading" @click="doImportFinger" :disabled="!selectedFile">开始导入</el-button>
      </span>
    </el-dialog>

    <el-dialog title="指纹详情" :visible.sync="detailVisible" width="780px" custom-class="theme-dialog" top="5vh">
      <div v-if="currentFinger" class="finger-detail">
        <div class="finger-detail-header">
          <h3>{{ currentFinger.fingerName }}</h3>
          <p class="finger-meta">
            <span class="finger-tag">{{ currentFinger.fingerClassEnum || '-' }}</span>
            <span v-if="currentFinger.level" class="finger-tag level-tag">L{{ currentFinger.level }}</span>
          </p>
        </div>

        <el-descriptions :column="2" border class="finger-info">
          <el-descriptions-item label="指纹名称">{{ currentFinger.fingerName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="应用名称">{{ currentFinger.appName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="指纹分类">{{ currentFinger.fingerClassEnum || '-' }}</el-descriptions-item>
          <el-descriptions-item label="版本">{{ currentFinger.version || '-' }}</el-descriptions-item>
        </el-descriptions>

        <div class="detail-section">
          <h4>指纹规则 (Rule)</h4>
          <pre class="rule-block">{{ currentFinger.rule || '暂无规则' }}</pre>
        </div>

        <div class="detail-section">
          <h4>指纹描述</h4>
          <p>{{ currentFinger.desc || '暂无描述' }}</p>
        </div>
      </div>
      <span slot="footer">
        <el-button @click="detailVisible = false">关闭</el-button>
      </span>
    </el-dialog>

    <el-dialog title="测试指纹" :visible.sync="testVisible" width="700px" custom-class="theme-dialog" top="10vh" :close-on-click-modal="false" @close="closeTestDialog">
      <el-form :model="testForm" label-width="100px">
        <el-form-item label="指纹名称">
          <el-input v-model="testForm.name" disabled />
        </el-form-item>
        <el-form-item label="指纹规则">
          <el-input type="textarea" :rows="3" v-model="testForm.rule" disabled />
        </el-form-item>
        <el-form-item label="测试地址">
          <el-input v-model="testForm.url" placeholder="请输入测试地址，如 http://example.com" />
        </el-form-item>
        <el-form-item label="测试日志">
          <el-input type="textarea" :rows="8" v-model="testForm.log" disabled placeholder="测试日志将在此显示..." />
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="closeTestDialog">关闭</el-button>
        <el-button type="primary" :loading="testLoading" @click="doFingerTest">测试</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { fingerprint } from '@/api/tool.js'

export default {
  name: 'AppFingerDB',
  data() {
    return {
      keyword: '',
      currentPage: 1,
      pageSize: 20,
      totalRecords: 0,
      tableData: [],
      tableLoading: false,
      detailVisible: false,
      currentFinger: null,
      importVisible: false,
      importLoading: false,
      selectedFile: null,
      importResult: null,
      testVisible: false,
      testLoading: false,
      testFingerId: 0,
      testForm: {
        name: '',
        rule: '',
        url: '',
        log: ''
      },
      testTimer: null,
      testCallId: ''
    }
  },
  mounted() {
    this.$store.state.activefirstMenu = '/appsec/finger-db'
    this.loadData()
  },
  methods: {
    async loadData() {
      this.tableLoading = true
      try {
        const params = {
          page: this.currentPage,
          size: this.pageSize,
          fingerName: this.keyword
        }
        const res = await fingerprint.getObjectData(params)
        if (res.code === 200) {
          this.tableData = res.data.list || []
          this.totalRecords = res.data.total || 0
        } else {
          this.tableData = []
          this.totalRecords = 0
        }
      } catch (e) {
        console.error('加载指纹列表失败:', e)
        this.tableData = []
        this.totalRecords = 0
      } finally {
        this.tableLoading = false
      }
    },
    handleSearch() {
      this.currentPage = 1
      this.loadData()
    },
    handleReset() {
      this.keyword = ''
      this.currentPage = 1
      this.loadData()
    },
    handleImport() {
      this.importVisible = true
      this.selectedFile = null
      this.importResult = null
    },
    handleFileChange(file, fileList) {
      this.selectedFile = fileList.length > 0 ? fileList[0].raw : null
      this.importResult = null
    },
    async doImportFinger() {
      if (!this.selectedFile) {
        this.$message.warning('请先选择文件')
        return
      }
      this.importLoading = true
      try {
        const res = await fingerprint.importFingerprint({ file: this.selectedFile })
        if (res.code === 200) {
          this.importResult = res.data
          if (res.data.success > 0) {
            this.$message.success(`导入成功，共 ${res.data.total} 条，成功 ${res.data.success} 条，跳过 ${res.data.skip} 条`)
            if (res.data.errors && res.data.errors.length === 0) {
              this.importVisible = false
              this.loadData()
            }
          }
        } else {
          this.$message({ message: res.msg || '导入失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '请求异常', type: 'error' })
      } finally {
        this.importLoading = false
      }
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
    async showDetail(row) {
      try {
        const res = await fingerprint.getInfo({ id: row.id })
        if (res.code === 200) {
          this.currentFinger = res.data
          this.detailVisible = true
        } else {
          this.$message({ message: res.msg || '获取详情失败', type: 'error' })
        }
      } catch (e) {
        this.$message({ message: '请求异常', type: 'error' })
      }
    },
    openTestDialog(row) {
      this.testVisible = true
      this.testForm.name = row.fingerName || ''
      this.testForm.rule = row.rule || ''
      this.testForm.url = ''
      this.testForm.log = ''
      this.testFingerId = row.id
      this.testCallId = ''
    },
    clearTestTimer() {
      if (this.testTimer) {
        clearTimeout(this.testTimer)
        this.testTimer = null
      }
    },
    closeTestDialog() {
      this.clearTestTimer()
      this.testVisible = false
      this.testFingerId = 0
      this.testCallId = ''
      this.testLoading = false
      this.testForm.log = ''
    },
    doFingerTest() {
      if (!this.testForm.url) {
        this.$message.error('请输入测试地址')
        return
      }
      this.clearTestTimer()
      this.testForm.log = ''
      this.testLoading = true
      const params = {
        fingerID: this.testFingerId,
        toolName: this.testForm.name,
        param: [{ key: 'root_url', value: this.testForm.url }]
      }
      fingerprint.testfinger(params).then((res) => {
        if (res.code === 200) {
          if (!res.data.token) {
            this.$message.error('测试任务启动失败：未返回 token')
            this.testLoading = false
            return
          }
          this.testCallId = res.data.token
          this.pollTestResult()
        } else {
          this.$message.error(res.msg || '测试失败')
          this.testLoading = false
        }
      }).catch(() => {
        this.$message.error('请求异常')
        this.testLoading = false
      })
    },
    pollTestResult() {
      if (!this.testCallId) return
      fingerprint.testfingerlog({ callId: this.testCallId }).then((res) => {
        if (res.code === 200) {
          if (res.data.result) {
            this.testForm.log += res.data.result
          } else if (!res.data.end) {
            this.testForm.log += '正在测试中...\n'
          }
          if (res.data.end) {
            this.clearTestTimer()
            this.testLoading = false
            this.testForm.log += '测试结束\n'
          } else {
            this.clearTestTimer()
            this.testTimer = setTimeout(() => {
              this.pollTestResult()
            }, 3000)
          }
        } else {
          this.clearTestTimer()
          this.testLoading = false
          this.$message.error(res.msg || '获取测试结果失败')
        }
      }).catch(() => {
        this.clearTestTimer()
        this.testLoading = false
        this.$message.error('请求异常')
      })
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
  max-width: 900px;
}

.db-stat {
  color: #94a3b8;
  font-size: 13px;
}

.finger-detail {
  color: #e2e8f0;
}

.finger-detail-header {
  margin-bottom: 20px;
}

.finger-detail-header h3 {
  margin: 0 0 8px;
  color: #00d4aa;
  font-size: 18px;
}

.finger-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
}

.finger-tag {
  display: inline-block;
  font-size: 12px;
  padding: 2px 10px;
  border-radius: 3px;
  background: rgba(0, 212, 170, 0.12);
  color: #00d4aa;
  font-weight: 500;
}

.level-tag {
  background: rgba(64, 158, 255, 0.12);
  color: #409eff;
}

.finger-info {
  margin-bottom: 20px;
}

.finger-info /deep/ .el-descriptions__label {
  color: #94a3b8;
  width: 90px;
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

.rule-block {
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  padding: 12px;
  color: #e2e8f0;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 250px;
  overflow-y: auto;
  margin: 0;
}

.finger-info /deep/ .el-descriptions {
  background: transparent;
}

.finger-info /deep/ table {
  background: rgba(15, 23, 42, 0.8);
  border-collapse: collapse;
}

.finger-info /deep/ .el-descriptions__cell,
.finger-info /deep/ .el-descriptions-item__label,
.finger-info /deep/ .el-descriptions-item__content {
  background: rgba(15, 23, 42, 0.9);
  border-color: rgba(255, 255, 255, 0.08);
  color: #94a3b8;
}

.finger-info /deep/ .el-descriptions-item__label {
  background: rgba(30, 41, 59, 0.9);
  color: #cbd5e1;
  border-color: rgba(255, 255, 255, 0.08);
}

.finger-info /deep/ .el-descriptions-item__content {
  background: rgba(15, 23, 42, 0.9);
  color: #94a3b8;
  border-color: rgba(255, 255, 255, 0.08);
}

.finger-info /deep/ table tr,
.finger-info /deep/ table td,
.finger-info /deep/ table th {
  background: rgba(15, 23, 42, 0.9);
}

.import-dialog {
  margin-bottom: 12px;

  /deep/ .el-upload {
    width: 100%;
  }

  /deep/ .el-upload-dragger {
    width: 100%;
    background: rgba(0, 0, 0, 0.35);
    border: 2px dashed rgba(255, 255, 255, 0.15);
    border-radius: 8px;
    transition: border-color 0.2s, background 0.2s;

    &:hover {
      border-color: rgba(0, 212, 170, 0.45);
      background: rgba(0, 212, 170, 0.06);
    }

    .el-icon-upload {
      color: #94a3b8;
      font-size: 48px;
      margin: 28px 0 12px;
      line-height: 1;
    }

    .el-upload__text {
      color: #94a3b8;
      font-size: 14px;

      em {
        color: #00d4aa;
        font-style: normal;
      }
    }
  }

  /deep/ .el-upload__tip {
    color: #64748b;
    font-size: 12px;
    margin-top: 10px;
    line-height: 1.5;
  }
}

.import-file-info {
  margin: 12px 0 0;
  color: #cbd5e1;
  font-size: 13px;
}

.import-result {
  margin: 12px 0 0;
  padding: 12px;
  border-radius: 4px;
  background: rgba(15, 23, 42, 0.8);
}

.import-success {
  margin: 0;
  color: #00d4aa;
  font-size: 13px;
  i {
    margin-right: 4px;
  }
}

.import-errors {
  margin: 8px 0 0;
  color: #f56565;
  font-size: 12px;
  span {
    display: block;
    line-height: 1.6;
  }
}
</style>