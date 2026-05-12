/* 
  ip详情组件
 */
<template>
  <div class="network-asset">
    <!-- 标题区域 -->
    <div class="header">
      <div class="ip-info">
        <span class="ip-address">{{ allGLObj.ip }}</span>
        <el-button v-if="allGLObj.ip" size="mini" style="margin-left: 20px" :type="buttonType">{{
          riskLevel
        }}</el-button>
      </div>
    </div>
    <!-- path -->
    <div style="color: #979797; margin-bottom: 20px">{{ allGLObj.path }}</div>
    <!-- 中间四个 -->
    <div class="row">
      <div class="column">
        <span class="icon">  <img style="width:40px" src="../../../assets/icon/assetManage/wuliduankou.png" alt="" srcset=""> </span>
        <div class="info">
          <div style="color:#ff6a00">端口</div>
          <div style="color:#ff6a00" class="ss">{{ allGLObj.portNum }}</div>
        </div>
      </div>
      <div class="column">
        <span class="icon"> <img style="width:40px" src="../../../assets/icon/assetManage/yuming.png" alt="" srcset=""></span>
        <div class="info">
          <div >域名</div>
          <div class="ss">{{ allGLObj.serviceNum }}</div>
        </div>
      </div>
      <div class="column">
        <span class="icon"> <img style="width:40px" src="../../../assets/icon/assetManage/zujian.png" alt="" srcset=""></span>
        <div class="info">
          <div style="color:#01c700">组件</div>
          <div style="color:#01c700" class="ss">{{ allGLObj.componentNum }}</div>
        </div>
      </div>
      <div class="column">
        <span class="icon"> <img style="width:40px" src="../../../assets/icon/assetManage/0.png" alt="" srcset=""></span>
        <div class="info">
          <div style="color:#d4237a">漏洞</div>
          <div style="color:#d4237a" class="ss">{{ allGLObj.vulNum }}</div>
        </div>
      </div>
    </div>
    <!-- 信息区域 -->
 <div class="card-section" style="">
  <div class="info-section">
    <div class="info-item">
      <span class="info-label">IP/域名：</span>
      <span class="info-content">{{ allGLObj.domain }}</span>
    </div>
    <div class="info-item">
      <span class="info-label">操作系统：</span>
      <span class="info-content">{{ allGLObj.opSystem }}</span>
    </div>
  </div>
  <div class="info-section">
    <div class="info-item">
      <span class="info-label">主机名：</span>
      <span class="info-content">{{ allGLObj.hostName }}</span>
    </div>
    <div class="info-item">
      <span class="info-label">硬件：</span>
      <span class="info-content">{{ allGLObj.hardware }}</span>
    </div>
  </div>
</div>

    
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="端口列表" name="first">
        <!-- 表格区域 -->
        <el-table style="margin-bottom: 30px" border :data="portList" stripe>
          <el-table-column
            prop="port"
            label="端口"
            width="80"
          ></el-table-column>
          <el-table-column
            prop="protocol"
            label="协议"
            width="80"
          ></el-table-column>
          <el-table-column prop="service" label="服务"></el-table-column>
          <el-table-column width="260px" prop="assembly" label="组件/指纹">
            <template slot-scope="scope">
              <span
                v-for="(item, index) in scope.row.assembly"
                :key="index"
                class="tag"
              >
                {{ item }}
              </span>
            </template>
          </el-table-column>
          <el-table-column
            prop="title"
            label="首页标题"
            width="120"
          ></el-table-column>

          <el-table-column
            prop="testTime"
            label="测试时间"
            width="180"
          ></el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="管理信息" name="second">
        <el-form ref="form" :model="form" label-width="100px">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="资产名称">
                <el-input v-model="form.assetName"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="业务系统">
                <el-input v-model="form.businessName"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="资产类型">
                <el-input v-model="form.assetType"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="虚拟资产">
                <el-input v-model="form.virtualAsset"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="备案等级">
                <!-- <el-input :min="0" v-model="form.filingLevel"></el-input> -->
                <!-- //备案等级是下拉框 -->
                <el-select v-model="form.filingLevel" placeholder="请选择">
                  <el-option
                    v-for="item in enumList"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  ></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="归属地">
                <el-input v-model="form.belongingPlace"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="责任部门">
                <el-input v-model="form.responsibleDepartment"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="责任人">
                <el-input v-model="form.responsiblePerson"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="责任人邮箱">
                <el-input v-model="form.responsiblePersonEmail"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24" style="text-align: right; margin-bottom: 30px">
              <el-form-item label="资产标签">
                <el-input v-model="form.assetLabels"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24" style="text-align: right; margin-bottom: 30px">
              <el-button type="primary" @click="onSubmit">编辑</el-button>
              <el-button @click="onReset">重置</el-button>
            </el-col>
          </el-row>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <!-- 漏洞测试 -->
    <el-table
      @cell-mouse-enter="mouseenter"
      @cell-mouse-leave="mouseleave"
      border
      :data="LouDongList"
      stripe
    >
      <el-table-column
        prop="vulName"
        label="漏洞名称"
        width="180"
      ></el-table-column>
      <el-table-column prop="vulNameType" label="漏洞类型"></el-table-column>
      <el-table-column prop="riskName" label="风险"></el-table-column>
      <el-table-column width="240px" prop="vulLocation" label="漏洞地址">
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100"></el-table-column>

      <el-table-column prop="testTime" label="测试时间" width="145">
        <!-- 详情、更改状态、删除 -->
        <template slot-scope="scope">
          <div v-if="showOperateButton && rowId == scope.row.id">
            <div>
              <el-link
                class="link_primary"
                :underline="false"
                @click="btnTaskinfo(scope.row)"
                >详情</el-link
              >
              <el-link
                class="link_primary"
                :underline="false"
                @click="changeStatus(scope.row)"
                >更改状态</el-link
              >

              <el-link
                class="link_primary"
                type="danger"
                :underline="false"
                @click="btnDel(scope.row)"
                >删除</el-link
              >
            </div>
          </div>
          <div v-else>
            <span><i></i>{{ scope.row.testTime }}</span>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <!-- 漏洞详情 -->
    <vulnmsginfo
      :typeNameWang="typeNameWang"
      v-model="sendVal"
      :vulninfo="vulninfo"
      :task_id="task_id"
      @saveData="handleSave()"
    >
    </vulnmsginfo>
    <!-- 更改状态的弹框 -->
      <el-dialog
            title="更改状态"
            width="400px"
            :visible.sync="dialogVisible"
            class="rulebox" 
            :close-on-click-modal="false" 
            :show-close="false" >
            <div class="dialog_b_btn">  
                <el-button size="small" @click="chaggeCLickStatus" >确定</el-button>
				<el-button size="small" @click="dialogVisible = false">取消</el-button>
            </div>
            <div class="buginfo_box" > 
                <div class="bugbasicinfo">
                    <!-- 更改状态的下拉框 -->
                    <el-form style="margin-top:30px" :model="formStatus" label-width="80px">
                        <el-form-item label="更改状态">
                            <el-select v-model="formStatus.status" placeholder="请选择">
                                <el-option
                                    v-for="(item,index) in LDenumList"
                                    :key="index"
                                    :label="item"
                                    :value="index"
                                >
                                </el-option>
                            </el-select>
                        </el-form-item>
                    </el-form>
                </div>   
            </div>
        </el-dialog>
  </div>
</template>


<script>
import { traffic } from '@/api/assetManagement.js'
import vulnmsginfo from '@/pages/taskManagement/Vulnmsginfo.vue'
export default {
  components: {
    vulnmsginfo,
  },
  data () {
    return {
      formStatus: {
        status: ''
      },
      dialogVisible: false,
      typeNameWang: '资产概览',
      enumList: [],
      LDenumList: [],
      task_id: '',
      sendVal: false,
      vulninfo: {},
      showOperateButton: false,
      rowId: '',
      LouDongList: [],
      form: {

        assetName: '',
        businessName: '',
        assetType: '',
        virtualAsset: '',
        filingLevel: '',
        belongingPlace: '',
        responsibleDepartment: '',
        responsiblePerson: '',
        responsiblePersonEmail: '',
        assetLabels: '',

      },
      portList: [],
      activeName: 'first',
      allGLObj: {},
      riskLevel: '高危', // 这个值可能是 '高危', '中危', '低危', '安全',
      tableData: [
        // 每行数据的格式

        // 更多数据...
      ]
    }
  },
  methods: {
   async chaggeCLickStatus(row){
    const res = await traffic.danDuZIChanvulinfoeditstatus({
      id: this.vulninfo.id,
      status: this.formStatus.status
    })
    if (res.code == 200) {
      this.$message({
        message: '更改成功',
        type: 'success'
      })
      this.getLouDongTestList()
      this.dialogVisible = false
    } else {
      this.$message({
        message: res.msg,
        type: 'error'
      })
    }
    },
    // 备案等级枚举
    async BeiAnEnum () { 
      const res = await traffic.assetDengJiEnum()
      if (res.code == 200) {
       console.log('备案等级枚举', res.data)
        this.enumList = res.data.assetFilingLevel
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 漏洞状态枚举
    async LuDongZhuangTaiEnum () { 
      const res = await traffic.danDuZIChanvulenums()
      if (res.code == 200) {
       console.log('漏洞状态枚举', res.data)
        this.LDenumList = res.data.vulStatus
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    async btnDel (row) { //单个删除
      const res = await traffic.danDuZIChanvultestinfodel({
        ids: row.id,
      })
      if (res.code == 200) {
        this.$message({
          message: '删除成功',
          type: 'success'
        })
        this.getLouDongTestList()
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    btnTaskinfo (r) {
      console.log(r,'333333333333333');
      this.vulninfo = r
      this.task_id = r.taskID +''
      this.sendVal = true
    },
    changeStatus(row) {
      console.log(row, '更改状态')
      this.dialogVisible = true
      this.vulninfo = row
    },
    onSubmit () {
      console.log('提交表单:', this.form)
      this.getGLXXEditList()
      // 这里可以添加提交表单的逻辑
    },
    onReset () {
      this.$refs.form.resetFields()
    },
    handleClick (tab, event) {
      console.log(tab, event)
    },
    // 单独资产概览
    async getZCList () {
      if (this.$store.state.groupID[1] == 1) return
      const res = await traffic.danDuZIChanGL({
        assetID: this.$store.state.groupID[0] || 7
      })
      if (res.code == 200) {
        this.allGLObj = res.data
        this.riskLevel = this.allGLObj.riskName
        console.log(res, '-------------------------------------------')

      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 资产概览-管理信息详情查询
    async getGLXXList () {
      const res = await traffic.danDuZIChanmanageinfo({
        assetID: this.$store.state.groupID[0] || 7
      })
      if (res.code == 200) {
        this.form = res.data
        console.log(res, '---------------222----------------------------')
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 漏洞测试-列表
    async getLouDongTestList () {
      const res = await traffic.assetLDTestList({
        assetID: this.$store.state.groupID[0] || 7,
        page: 1,
        size: 12
      })
      console.log(res, '------999999999999999999999---------漏洞测试-列表 ----------------------------')

      if (res.code == 200) {
        this.LouDongList = res.data.assetVulTestListInfo
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 漏洞测试- 详情
    async getLouDongTestDetailList () {
      const res = await traffic.assetLDTestList({
        assetID: this.$store.state.groupID[0] || 7,
        page: 1,
        size: 12
      })

      if (res.code == 200) {
        this.LouDongList = res.data.assetVulTestListInfo
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 资产概览-管理信息详情修改 --todo
    async getGLXXEditList () {
      const res = await traffic.danDuZIChanmanageinfoupdate({ ...this.form, assetID: this.$store.state.groupID[0] || 7 })
      if (res.code == 200) {
        // 提示成功
        this.$message({
          message: res.msg || '修改成功',
          type: 'success'
        })
        this.getGLXXList()
        // 打印`122

      } else {
        this.$message({
          message: res.msg,
          type: 'error'

        })
      }
    },
    // 端口列表
    async getPortList () {
      const res = await traffic.danDuZIChanportlist({
        assetID: this.$store.state.groupID[0] || 7,
        page: 1,
        size: 15
      })
      if (res.code == 200) {
        this.portList = res.data.assetPortListInfo
        console.log(this.portList, 'this.portListthis.portListthis.portListthis.portList')
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    mouseenter (row, colum, cell, event) {
      this.showOperateButton = true
      this.rowId = row.id   //赋值行id，便于页面判断  
    },
    mouseleave (row, colum, cell, event) {
      if (row.executeType == 3 || row.status != 1) {
        if (!this.$refs['popover-' + row.id]) {
          this.showOperateButton = false
          this.rowId = ""
          return
        } else {
          let isShow = this.$refs['popover-' + row.id].showPopper
          if (!isShow) {
            this.showOperateButton = false
            this.rowId = ""
          }
        }
      } else {
        if (!this.$refs['popover_id-' + row.id]) {
          this.showOperateButton = false
          this.rowId = ""
          return
        } else {
          let isShow = this.$refs['popover_id-' + row.id].showPopper
          if (!isShow) {
            this.showOperateButton = false
            this.rowId = ""
          }
        }
      }
    },
  },
  computed: {
    buttonType () {
      switch (this.riskLevel) {
        case '高危':
          return 'danger'
        case '中危':
          return 'warning'
        case '低危':
          return 'info'
        case '安全':
          return 'success'
        case '未知':
          return 'info'
        default:
          return 'default' // 默认不设置类型，使用原始按钮样式
      }
    }
  },
  watch: {
    // 监听 Vuex store 中的 groupID
    '$store.state.groupID' (newVal, oldVal) {
      // 当 groupID 发生变化时，执行一些操作
      // newVal 是新的 groupID 值，oldVal 是旧的值
      if (newVal) {
        this.$nextTick(() => {
          this.getZCList()
          this.getPortList()
          this.getGLXXList()
          this.getLouDongTestList()
        })
      }
    }
  },
  created () {
    this.getZCList()
    this.getPortList()
    this.getGLXXList()
    this.getLouDongTestList()
    this.BeiAnEnum()
    this.LuDongZhuangTaiEnum()
  },
};
</script>



<style scoped>
.row {
  display: flex;
}
.column {
  flex: 1;
  margin: 2px;
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
}
.icon {
  font-size: 20px;
}
.info {
  text-align: left;
  padding: 0 20px;
  background: #f7f9fc;
}
.info .ss {
  font-size: 22px;
  font-weight: 700;
}
.info > div:first-child {
  font-weight: bold;
}
.info > div:last-child {
  color: #4c7ae3;
}
</style>
<style scoped>
.tag {
  background-color: #d4edda; /* 淡绿色背景 */
  color: #155724; /* 较深的绿色文本，用于对比 */
  padding: 0.25em 0.5em; /* 标签内边距 */
  margin: 2px; /* 标签间距 */
  border-radius: 0.25rem; /* 圆角边框 */
  display: inline-block; /* 行内块显示 */
  font-size: 0.875em; /* 字体大小 */
}
.network-asset {
  /* 整体布局样式 */
  background: #fff;
  padding: 20px;
}

.header {
  /* 标题区域样式 */
  display: flex;
  justify-content: space-between;
  padding-bottom: 16px;
  /* background-color: #f3f3f3; */
}

.ip-info {
  /* IP信息样式 */
  display: flex;
  align-items: center;
}

.ip-address {
  /* IP地址样式 */
  font-size: 20px;
  margin-right: 8px;
  font-weight: 700;
}
.statistics {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.statistics .stat-item {
  /* 统计信息单项样式 */

  margin-right: 24px;
  border: 1px solid red;
}

.stat-number {
  /* 统计数字样式 */
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.stat-text {
  /* 统计文本样式 */
  margin-left: 4px;
}

.info-section {
  /* 信息区域样式 */
  display: flex;
  flex-wrap: wrap;
  /* padding: 16px; */
  /* background-color: #fff; */
  /* margin: 20px 0; */
}

.info-item {
  /* 单个信息项样式 */
  flex: 1;
  
}

.info-label {
  /* 信息标签样式 */
  font-weight: bold;
  text-align-last: left;
}

.info-content {
  /* 信息内容样式 */
}

/* Element UI 表格样式调整可以在此处添加 */
</style>

<style scoped>
.info-section {
  background: linear-gradient(135deg, #f7f9fc 0%, #f7f9fc 100%);
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 8px 20px rgba(0,0,0,0.12);
  margin: 40px auto;
  border: 1px solid rgba(0,0,0,0.1);
}

.info-item {
  display: flex;
  align-items: center;
  /* margin-bottom: 15px; */
  position: relative;
}

.info-label {
  font-weight: 600;
  color: #4A5A6A;
  font-size: 14px;
  flex-basis: 100px;
}

.info-content {
  
  color: #6B7C93;
  background-color: #fff;
  border-radius: 5px;
  padding: 10px 15px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.05);
  flex-grow: 0.5;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  font-size: 13px;
}

.info-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 3px;
  height: 100%;
  background: linear-gradient(180deg, #4c7ae3 0%, #6a8ee6 100%);
  border-radius: 5px 0 0 5px;
}

.info-item:hover .info-content {
  transform: translateX(5px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

/* Optional: Add icons */
.info-item .info-icon {
  margin-right: 10px;
  color: #f7645b;
  font-size: 20px; /* Adjust size based on your icon set */
}

</style>


<style>
.row {
  display: flex;
  justify-content: space-around;
  padding: 20px;
  background: linear-gradient(135deg, #f7f9fc 0%, #f7f9fc 100%);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
  margin-bottom: 20px; /* 如果有多行，保持间距 */
}

.column {
  flex: 1;
  text-align: center; /* 确保内容居中 */
  margin: 0 10px; /* 给列之间一些间距 */
}

.icon {
  font-size: 24px; /* 调整图标大小 */
  display: inline-block;
  margin-bottom: 8px; /* 图标和文字之间的间距 */
}

.info {
  background: #FFFFFF;
  border-radius: 5px;
  padding: 10px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  transition: all 0.3s ease;
}

.info:hover {
  transform: translateY(-5px); /* 悬停效果 */
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
}

.info > div:first-child {
  color: #4c7ae3; /* 标题颜色 */
  font-weight: 600; /* 加粗 */
  margin-bottom: 5px; /* 与数值间的间距 */
}

.ss {
  font-size: 20px; /* 数值大小 */
  color: #333; /* 数值颜色 */
}
</style>


