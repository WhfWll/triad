<template>
  <!-- 资产详情 -->
  <div class="right-ip-zichan" style="padding: 20px;">
    <!-- ip/域名 -->
    <div class="ip-domain" style="margin-bottom:10px">
      <div class="ip-domain-item">
        <strong style="font-size:18px">{{ getDisplayData().ip || detailOBJ.manageInfo?.ip }}</strong>
         <!-- riskLevelStr -->
         <span v-show="detailOBJ.manageInfo?.riskLevelStr!='安全'" v-if="detailOBJ.manageInfo?.riskLevelStr" :style="{background:detailOBJ.manageInfo?.riskLevelStr=='高危'?'red':detailOBJ.manageInfo?.riskLevelStr=='中危'?'orange':detailOBJ.manageInfo?.riskLevelStr=='低危'?'#8ac23a':'#16c60c'}"  style="color:#fff;margin-left:10px;padding:3px 8px;border-radius:10px;">{{detailOBJ.manageInfo?.riskLevelStr||''}}</span>
           <span v-show="detailOBJ.manageInfo?.score !='0分'" v-if="detailOBJ.manageInfo?.score"  style="color:#fff;margin-left:10px;padding:3px 8px;border-radius:10px;background-color:red">{{detailOBJ.manageInfo?.score||''}}</span>
      </div>
    </div>
    <!-- 资产名称 -->
    <!-- 资产详情区域 -->
    <div class="asset-details">
      <el-row :gutter="20">
        <el-col :span="12">
          <div class="detail-item"><strong>IP/域名:</strong> {{ getDisplayData().ip || detailOBJ.manageInfo?.ip }}

          </div>
          <div class="detail-item"><strong>资产名称:</strong> {{ detailOBJ.manageInfo?.assetName||'' }} </div> 
          <div class="detail-item"><strong>资产组:</strong>{{ detailOBJ.manageInfo?.assetGroupName }}  </div>
          <div class="detail-item"><strong>部门:</strong> {{ detailOBJ.manageInfo?.responsibleDepartment }}</div> 
        </el-col>
        <el-col :span="12">
          <div class="detail-item"><strong>操作系统:</strong> {{ getDisplayData().opSys || detailOBJ.manageInfo?.opSys }}</div> 
          <div class="detail-item"><strong>系统管理员:</strong> {{ detailOBJ.manageInfo.systemAdmin }}</div>
          <div class="detail-item"><strong>标签:</strong> {{ detailOBJ.manageInfo.tags }}</div> 
         
        </el-col>
      </el-row>
    </div>

    <!-- Tab区域 -->
    <el-tabs v-model="activeTab" @tab-click="handleTabClick">
      <!-- 资产信息 -->
      <el-tab-pane label="资产信息" name="assetInfo">
        <el-table
          class="tabsList"
          :data="detailOBJ.assetInfo?.list || detailOBJ.assetInfo || []"
          style="width: 100%;"
          v-loading="assetLoading"
        >
          <el-table-column prop="port" label="端口" width="100"></el-table-column>
          <el-table-column prop="protocol" label="协议" width="100"></el-table-column>
          <el-table-column prop="service" label="服务" width="120"></el-table-column>
          <el-table-column prop="assembly" label="指纹" :show-overflow-tooltip="true"></el-table-column>
          <!-- <el-table-column prop="remark" label="首页标题" :show-overflow-tooltip="true"></el-table-column> -->
          <el-table-column prop="createTime" label="测试时间" width="180"></el-table-column>
        </el-table>

        <!-- 分页组件 -->
        <div class="pagination-container" v-if="assetPagination.total > 0">
          <el-pagination
            @size-change="(size) => handleSizeChange('asset', size)"
            @current-change="(page) => handleCurrentChange('asset', page)"
            :current-page="assetPagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="assetPagination.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="assetPagination.total"
          >
          </el-pagination>
        </div>
      </el-tab-pane>

      <!-- 漏洞信息 -->
      <el-tab-pane label="漏洞信息" name="vulInfo">
        <el-table
          class="tabsList"
          :data="detailOBJ.vulInfo?.list || []"
          style="width: 100%;"
          v-loading="vulLoading"
        >
          <el-table-column prop="name" label="漏洞名称" :show-overflow-tooltip="true"></el-table-column>
          <el-table-column prop="typeName" label="漏洞类型" width="120"></el-table-column>
          <el-table-column prop="riskName" label="风险" width="100">
            <template #default="scope">
              <!-- <span :style="getRiskLevelColor(scope.row.riskLevelStr)">
                {{ scope.row.riskLevelStr }}
              </span> -->
              <span :class="[
                { 'riskstyle risk_hight': scope.row.riskLevel == '1' },
                { 'riskstyle risk_middle': scope.row.riskLevel == '2' },
                { 'riskstyle risk_low': scope.row.riskLevel == '3' },
                { 'riskstyle risk_nofind': scope.row.riskLevel == '4' }
              ]"><i></i>{{ scope.row.riskLevelStr }} </span>
              </template>
          </el-table-column>
          <el-table-column prop="location" label="漏洞地址" :show-overflow-tooltip="true"></el-table-column>
          <el-table-column prop="verifyTypeName" label="验证方式" width="120"></el-table-column>
          <el-table-column prop="statusName" label="处置状态" width="100"></el-table-column>
          <el-table-column prop="findTime" label="测试时间" width="180"></el-table-column>
        </el-table>
        <!-- <div v-if="!(detailOBJ.vulInfo?.list || []).length && !vulLoading" class="empty-data">
          暂无漏洞信息
        </div> -->

        <!-- 分页组件 -->
        <div class="pagination-container" v-if="vulPagination.total > 0">
          <el-pagination
            @size-change="(size) => handleSizeChange('vul', size)"
            @current-change="(page) => handleCurrentChange('vul', page)"
            :current-page="vulPagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="vulPagination.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="vulPagination.total"
          >
          </el-pagination>
        </div>
      </el-tab-pane>

      <!-- 配置安全 -->
      <el-tab-pane label="配置安全" name="configInfo" v-if="false">
        <el-table
          class="tabsList"
          :data="detailOBJ.configInfo?.list || detailOBJ.configInfo || []"
          style="width: 100%;"
          v-loading="configLoading"
        >
          <el-table-column prop="itemCheckedName" label="核查项名称" :show-overflow-tooltip="true"></el-table-column>
          <el-table-column prop="itemCheckedDesc" label="核查项描述" :show-overflow-tooltip="true"></el-table-column>
          <el-table-column prop="checkRes" label="核查结果" width="120">
            <template #default="scope">
              <el-tag
                :type="getCheckResultType(scope.row.checkRes)"
                size="small"
              >
                {{ getCheckResultText(scope.row.checkRes) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="checkTime" label="完成时间" width="180"></el-table-column>
          <!-- 操作 -->
          <el-table-column prop="checkTime" label="操作" width="180">
            <template #default="scope">
              <el-button type="primary" size="small" @click="handleCheck(scope.row)">详情</el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页组件 -->
        <div class="pagination-container" v-if="configPagination.total > 0">
          <el-pagination
            @size-change="(size) => handleSizeChange('config', size)"
            @current-change="(page) => handleCurrentChange('config', page)"
            :current-page="configPagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="configPagination.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="configPagination.total"
          >
          </el-pagination>
        </div>
      </el-tab-pane>

      <!-- 管理信息 -->
      <el-tab-pane label="管理信息" name="manageInfo"  v-if="false">
        <div class="management-details">
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="detail-item"><strong>IP/域名:</strong> {{ detailOBJ.manageInfo?.ip || '-' }}</div>
              <div class="detail-item"><strong>操作系统:</strong> {{ detailOBJ.manageInfo?.opSys || '-' }}</div>
              <!-- <div class="detail-item"><strong>部署位置:</strong> {{ detailOBJ.manageInfo?.deploymentLocation || '-' }}</div> -->
              <!-- <div class="detail-item"><strong>运维负责人:</strong> {{ detailOBJ.manageInfo?.systemOp || '-' }}</div> -->
              <div class="detail-item"><strong>标签:</strong> {{ detailOBJ.manageInfo?.tags || '-' }}</div>
            </el-col>
            <el-col :span="8">
                <!-- <div class="detail-item"><strong>IP地址段:</strong> {{ detailOBJ.manageInfo?.ipSegment || '-' }}</div>
                <div class="detail-item"><strong>系统名称:</strong> {{ detailOBJ.manageInfo?.systemName || '-' }}</div>
              <div class="detail-item"><strong>系统用途:</strong> {{ detailOBJ.manageInfo?.purpose || '-' }}</div>
              <div class="detail-item"><strong>基础软件名称:</strong> {{ detailOBJ.manageInfo?.baseSoftwareName || '-' }}</div> -->
              <div class="detail-item"><strong>资产名称:</strong> {{ detailOBJ.manageInfo?.assetName || '-' }}</div>
              <!-- 设备权重 -->
              <div class="detail-item"><strong>设备权重:</strong> {{ detailOBJ.manageInfo?.deviceWeight || '-' }}</div>
              <!-- 可信设备 -->
              <div class="detail-item"><strong>可信设备:</strong> {{ detailOBJ.manageInfo?.trustLevel || '-' }}</div>

            </el-col>
            <el-col :span="8">
              <div class="detail-item"><strong>归属资产组:</strong> {{ detailOBJ.manageInfo?.assetGroupName || '-' }}</div>
              <div class="detail-item"><strong>设备形态:</strong> {{ detailOBJ.manageInfo?.equipmentForm || '-' }}</div>
              <!-- <div class="detail-item"><strong>基础软件版本:</strong> {{ detailOBJ.manageInfo?.baseSoftwareVersion || '-' }}</div> -->
              <!-- <div class="detail-item"><strong>基础硬件名称:</strong> {{ detailOBJ.manageInfo?.baseHardwareName || '-' /></div> -->
              <div class="detail-item"><strong>管理员:</strong> {{ detailOBJ.manageInfo?.systemAdmin || '-' }}</div>
            </el-col>
          </el-row>
          <el-row :gutter="20" style="margin-top: 20px;">
            <el-col :span="8">
              <div class="detail-item"><strong>等保级别:</strong> {{ fillingLevel }}</div>
              <!-- <div class="detail-item"><strong>是否云主机:</strong> {{ detailOBJ.manageInfo?.isCloudHost ? '是' : '否' }}</div> -->
            </el-col>
          </el-row>
           <el-row :gutter="20" style="margin-top: 20px;">
            <el-col :span="6">
              <div class="detail-item"><strong>登录用户名:</strong> 
                <!-- {{ detailOBJ.manageInfo?detailOBJ.manageInfo.connections[0]?detailOBJ.manageInfo.connections[0].username :'' :''}} -->
              </div>
             </el-col> 
             <el-col :span="6"> 
             <div class="detail-item"><strong>登录密码:</strong> 
              <!-- {{ detailOBJ.manageInfo?detailOBJ.manageInfo.connections[0]?detailOBJ.manageInfo.connections[0].password :'':'' }} -->
            </div>
               </el-col> 
             <el-col :span="6">  
             <div class="detail-item"><strong>登录端口:</strong> 
              <!-- {{ detailOBJ.manageInfo?detailOBJ.manageInfo.connections[0]?detailOBJ.manageInfo.connections[0].port :'':'' }} -->
            </div>
               </el-col> 
             <el-col :span="6">  
             <div class="detail-item"><strong>连接协议:</strong> 
              <!-- {{ detailOBJ.manageInfo?detailOBJ.manageInfo.connections[0]?detailOBJ.manageInfo.connections[0].protocol :'':'' }} -->
            </div>
            </el-col>
          </el-row>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 详情弹框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="60%"
      :top="'10vh'"
      :append-to-body="true"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <div v-if="Object.keys(dialogData).length === 0" class="no-data">
          暂无数据
        </div>
        <div v-else>
          <div class="config-detail-item" v-for="(value, key) in dialogData" :key="key">
           <div v-if="key!='checkDetail'">
            <!-- {{key  }} -->
           <div v-show="key !='id'&&key != 'checkResMsg'&&key != 'ip'&&key != 'itemCategoryName'">
            <div class="detail-label">{{ getFieldLabel(key) }}:</div>
            <div class="detail-value">{{ formatFieldValue(key, value) }}</div>
           </div>
           </div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'RightIPZiChan',
  props: {
    detailOBJ: {
      type: Object,
      default: () => ({})
    },
    selectType: {
      type: Number,
      default: 1
    }
  },
  data() {
    return {
      loading: false,
      dialogVisible: false,
      dialogTitle: '',
      dialogData: {},
      // 分页数据
      assetPagination: { page: 1, size: 10, total: 0 },
      vulPagination: { page: 1, size: 10, total: 0 },
      configPagination: { page: 1, size: 10, total: 0 },
      // loading状态
      assetLoading: false,
      vulLoading: false,
      configLoading: false
    };
  },
  computed: {
    activeTab: {
      get() {
        const tabMap = {
          1: 'assetInfo',
          2: 'vulInfo',
          3: 'configInfo',
          4: 'manageInfo'
        };
        return tabMap[this.selectType] || 'assetInfo';
      },
      set(value) {
        // 当tab改变时，通知父组件
        const typeMap = {
          'assetInfo': 1,
          'vulInfo': 2,
          'configInfo': 3,
          'manageInfo': 4
        };
        const newSelectType = typeMap[value];
        if (newSelectType && newSelectType !== this.selectType) {
          this.$emit('tab-change', newSelectType);
        }
      }
    },
    assetTypeLabel() {
      const assetTypes = [
        { "label": "OA系统", "value": 1 },
        { "label": "办公设备", "value": 2 },
        { "label": "网络设备", "value": 3 },
        { "label": "安全设备", "value": 4 },
        { "label": "视频监控", "value": 5 },
        { "label": "工控设备", "value": 6 },
        { "label": "IOT", "value": 7 },
        { "label": "业务系统", "value": 8 }
      ];
      const assetType = this.detailOBJ.manageInfo?.assetType || this.detailOBJ.assetType;
      const type = assetTypes.find(type => type.value === assetType);
      return type ? type.label : '未知类型';
    },
    fillingLevel() {
      const levels = [
        { "label": "无", "value": 1 },
        { "label": "等保一级", "value": 2 },
        { "label": "等保二级", "value": 3 },
        { "label": "等保三级", "value": 4 },
        { "label": "等保四级", "value": 5 },
        { "label": "等保五级", "value": 6 }
      ];
      const level = this.detailOBJ.manageInfo?.equalProtectionLevel || this.detailOBJ.equalProtectionLevel;
      const type = levels.find(type => type.value === level);
      return type ? type.label : '未知级别';
    }
  },
  methods: {
    handleCheck(row) {
      console.log('handleCheck called with:', row);
      this.dialogVisible = true;
      this.dialogTitle = `配置安全详情 `;
      this.dialogData = { ...row };
    },
    getFieldLabel(key) {
      const labelMap = {
        'itemCheckedName': '检测项名称',
        'itemCheckedDesc': '检测项描述',
        'checkRes': '检测结果',
        'checkTime': '测试时间',
        'createTime': '创建时间',
        'updateTime': '更新时间',
        'remark': '备注',
        'itemCategoryName': '检测项分类',
        'operator': '操作员',
        'checkDetail': '检查详情',
        'suggestion': '建议',
        'riskLevel': '风险等级'
      };
      return labelMap[key] || key;
    },
    formatFieldValue(key, value) {
      if (value === null || value === undefined) {
        return '-';
      }

      if (key === 'checkRes') {
        return this.getCheckResultText(value);
      }

      if (typeof value === 'object') {
        return JSON.stringify(value, null, 2);
      }

      return String(value);
    },
    handleTabClick(tab) {
      console.log('Tab changed:', tab);
      // 当切换tab时，通知父组件获取对应数据
      // this.emitTabChange(this.activeTab);
      // v-model会自动更新activeTab的值，不需要手动设置
      // this.activeTab = tab.name; // 移除这行，避免重复触发
    },
    // 页面大小改变
    handleSizeChange(type, newSize) {
      if (type === 'asset') {
        this.assetPagination.size = newSize;
        this.assetPagination.page = 1;
        this.emitTabChange('assetInfo');
      } else if (type === 'vul') {
        this.vulPagination.size = newSize;
        this.vulPagination.page = 1;
        this.emitTabChange('vulInfo');
      } else if (type === 'config') {
        this.configPagination.size = newSize;
        this.configPagination.page = 1;
        this.emitTabChange('configInfo');
      }
    },
    // 当前页改变
    handleCurrentChange(type, newPage) {
      if (type === 'asset') {
        this.assetPagination.page = newPage;
        this.emitTabChange('assetInfo');
      } else if (type === 'vul') {
        this.vulPagination.page = newPage;
        this.emitTabChange('vulInfo');
      } else if (type === 'config') {
        this.configPagination.page = newPage;
        this.emitTabChange('configInfo');
      }
    },
    // 发射tab改变事件
    emitTabChange(tabName) {
      const typeMap = {
        'assetInfo': 1,
        'vulInfo': 2,
        'configInfo': 3,
        'manageInfo': 4
      };
      const selectType = typeMap[tabName];

      // 获取对应的分页参数
      let pagination = { page: 1, size: 10 };
      if (tabName === 'assetInfo') {
        pagination = this.assetPagination;
      } else if (tabName === 'vulInfo') {
        pagination = this.vulPagination;
      } else if (tabName === 'configInfo') {
        pagination = this.configPagination;
      }

      this.$emit('tab-change', selectType, pagination.page, pagination.size);
    },
    // 更新分页总数
    updatePagination(data) {
      if (this.selectType === 1 && data.assetInfo) {
        this.assetPagination.total = data.assetInfo.total || 0;
      } else if (this.selectType === 2 && data.vulInfo) {
        this.vulPagination.total = data.vulInfo.total || 0;
      } else if (this.selectType === 3 && data.configInfo) {
        this.configPagination.total = data.configInfo.total || 0;
      }
    },
    // 设置loading状态
    setLoading(type, loading) {
      if (type === 1) {
        this.assetLoading = loading;
      } else if (type === 2) {
        this.vulLoading = loading;
      } else if (type === 3) {
        this.configLoading = loading;
      }
    },
    getRiskLevelColor(riskName) {
      const riskColors = {
        '致命': '#ff0000',
        '高危': '#ffa500',
        '中危': '#8ac23a',
        '低危': '#16c60c',
        '信息': '#000000'
      };
      return { color: riskColors[riskName] || '#666' };
    },
    getCheckResultType(checkRes) {
      const typeMap = {
        1: 'success', // 通过
        2: 'danger',  // 不通过
        3: 'warning'  // 错误
      };
      return typeMap[checkRes] || 'info';
    },
    getCheckResultText(checkRes) {
      const textMap = {
        1: '通过',
        2: '不通过',
        3: '错误'
      };
      return textMap[checkRes] || '未知';
    },
    getDisplayData() {
      // 优先从对应的数据源获取显示数据
      return this.detailOBJ.manageInfo || this.detailOBJ || {};
    }
  },
  watch: {
    selectType: {
      immediate: false, // 改为false，避免初始化时触发
      handler(newVal, oldVal) {
        console.log('selectType changed from', oldVal, 'to', newVal);
        // 不在这里触发事件，避免重复调用
      }
    },
    detailOBJ: {
      deep: true,
      immediate: true,
      handler(newVal) {
        if (newVal) {
          this.updatePagination(newVal);
        }
      }
    }
  },
  mounted() {
    if (this.detailOBJ) {
      this.updatePagination(this.detailOBJ);
    }
  }
};
</script>

<style scoped>
.detail-item {
  margin-bottom: 10px;
  font-size: 14px;
}

.tabsList {
  overflow-y: auto;
  max-height: 400px;
}

.empty-data {
  text-align: center;
  padding: 40px;
  color: #999;
  font-size: 14px;
}

.total-info {
  margin-top: 10px;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
  font-size: 14px;
  color: #666;
}

.management-details {
  padding: 20px 0;
}

.right-ip-zichan .asset-details {
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 4px;
  margin-bottom: 20px;
}

/* 弹框样式 */
.dialog-content {
  max-height: 60vh;
  overflow-y: auto;
}

.config-detail-item {
  margin-bottom: 16px;
  padding-bottom: 12px;
}

.config-detail-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.detail-label {
  font-weight: bold;
  color: #303133;
  font-size: 14px;
}

.detail-value {
  color: #606266;
  font-size: 14px;
  word-break: break-all;
  white-space: pre-wrap;
  line-height: 1.5;
}

.dialog-footer {
  text-align: right;
}

.no-data {
  text-align: center;
  padding: 40px;
  color: #999;
  font-size: 14px;
}

/* 确保弹框可见 */
:deep(.el-dialog) {
  z-index: 9999 !important;
}

:deep(.el-overlay) {
  z-index: 9998 !important;
}

/* 分页样式 */
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding: 16px 0;
}
</style>
