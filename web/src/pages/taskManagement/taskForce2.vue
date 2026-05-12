<template>
  <!-- 流量分析 -->
  <div>
    <div class="main-title">
      <label for="">资产管理</label>
    </div>
    <div class="trafficlist context_box_bg">
      <div class="search-box">
        <div class="serach-condition">
          <div class="search-text">
            <el-input
              style="width: 300px"
              placeholder="请输入关键字"
              @keydown.enter.native="handlesearch"
              v-model="formData.search_field"
              class="input-with-select"
              size="small"
              clearable
            >
            </el-input>
            <el-button type="primary" size="small" @click="handlesearch"
              >搜索</el-button
            >
          </div>
          <div>
            <el-button type="primary" size="small" @click="handleReset"
              >重置</el-button
            >
          </div>
          <div>
            <el-button type="primary" size="small" @click="headervisible = true"
              >高级搜索</el-button
            >
          </div>
          <div>
            <el-dropdown trigger="click" @command="handleCommand">
              <el-button type="primary">
                更多操作<i class="el-icon-plus el-icon--right"></i>
              </el-button>
              <el-dropdown-menu slot="dropdown">
                <!-- <el-dropdown-item command="penetrationtest" >渗透测试</el-dropdown-item> -->
                <!-- <el-dropdown-item  command="addAsset" >新增资产</el-dropdown-item> -->
                <!-- <el-dropdown-item command="addAssetGroup">新增资产组</el-dropdown-item> -->
                <el-dropdown-item command="findAsset">发现资产</el-dropdown-item> 
                <el-dropdown-item command = "importAsset">导入资产</el-dropdown-item>
                <el-dropdown-item command="syncAsset">资产同步</el-dropdown-item> 
                <!-- <el-dropdown-item command="selectedAsset">选择资产</el-dropdown-item> -->
                <!-- <el-dropdown-item command="cancelSelected">取消选择</el-dropdown-item> -->
                <!-- <el-dropdown-item command="assetdelete">批量删除</el-dropdown-item> -->
              </el-dropdown-menu>
            </el-dropdown>
          </div>
        </div>
      </div>
    </div>
    
    <div style="margin-top: 20px; display: flex; justify-content: space-between">
      <!-- 左 -->
      <div style="flex: 1;overflow:hidden">
        <el-tabs v-model="nameTab" type="border-card">
          <el-tab-pane  label="资产统计">
            <!-- <div style="text-align: right">查看更多</div> -->
            <LeftZiChan @clickSOnTOTab= 'clickSOnTOTab' :zichanObj="zichanObj"></LeftZiChan>
          </el-tab-pane>
          <el-tab-pane label="资产树">
            <RightZiChan ref="RightZiChan" @updateAsset="updateAsset" @updateAssetGroup="updateAssetGroup" @clicknode="hanldeclicknode"></RightZiChan>
          </el-tab-pane>
        </el-tabs>
      </div>
      <!-- 右 -->
      <div style="flex: 2; margin-left: 20px">
        <el-tabs v-model="activeTabName" type="border-card" @tab-click="handleClick" v-if="nameTab == 1&&$store.state.groupID[1] != 2">
          <el-tab-pane name="one"  label="概览"><Gailan ref="gailan" /></el-tab-pane>
          <el-tab-pane name="tw" label="子资产组"><ZIZIC ref="zizichanzu" @gotoFirstTab='gotoFirstTab' /></el-tab-pane>
          <el-tab-pane name="sa" label="资产"><ZIChan ref="zichan" @gotoFirstTab='gotoFirstTab' /></el-tab-pane>
          <el-tab-pane name="fou" label="资产漏洞"><ZIChanLD ref="zichanloudong" @gotoFirstTab='gotoFirstTab'/></el-tab-pane>
        </el-tabs>
        <div v-if="nameTab == 1&&$store.state.groupID[1] == 2"><RightIPZiChan /></div>
        <RightChar v-if="nameTab == 0"></RightChar>
      </div>
    </div>
    <!-- 弹窗 -->
    <el-dialog
      title="高级搜索"
      :visible.sync="headervisible"
      width="1184px"
      class="hashbox"
      :close-on-click-modal="false"
      :show-close="false"
      @close="closeHeaderDialog"
    >
      <div class="dialog_b_btn">
        <el-button size="small" @click="saveHeaderForm">搜索</el-button>
        <el-button size="small" @click="resetHeaderFormFields">重置</el-button>
        <el-button size="small" @click="handleCloseHeader"
          >隐藏高级搜索</el-button
        >
      </div>
      <div class="hash-mainDiv">
        <el-form
          style=" padding: 20px 10px"
          :model="generateHeaderForm"
          class="hashEl-form"
          label-width="120px"
          status-icon
          ref="headerForm"
          
          inline
        >
          <!-- 行 1 -->
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="IP/域名">
                <el-input
                  v-model="generateHeaderForm.ip"
                  placeholder="请输入IP/域名"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="操作系统">
                <el-input
                  v-model="generateHeaderForm.operateSystem"
                  placeholder="请输入操作系统"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="资产风险">
                <!-- 这里应根据你的数据结构使用适当的组件，例如el-select -->
                <el-select
                  class="typeSelect"
                  v-model="generateHeaderForm.riskLevel"
                  placeholder="请选择"
                  clearable
                  size="small"
                >
                  <el-option
                    v-for="item in zichanEnumObj.assetRisk"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  >
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 行 2 -->
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="资产名称">
                <el-input
                  v-model="generateHeaderForm.assetName"
                  placeholder="请输入资产名称"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="业务系统">
                <el-input
                  v-model="generateHeaderForm.businessSystem"
                  placeholder="请输入业务系统"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="责任部门">
                <el-input
                  v-model="generateHeaderForm.responsibleDepartment"
                  placeholder="请输入责任部门"
                ></el-input>
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 行 3 -->
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="端口">
                <el-input
                  v-model="generateHeaderForm.port"
                  placeholder="请输入端口"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="服务">
                <el-input
                  v-model="generateHeaderForm.service"
                  placeholder="请输入服务"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="组件">
                <el-input
                  v-model="generateHeaderForm.component"
                  placeholder="请输入组件"
                ></el-input>
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 行 4 -->
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="漏洞名称">
                <el-input
                  v-model="generateHeaderForm.vulName"
                  placeholder="请输入漏洞名称"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="标签">
                <el-input
                  v-model="generateHeaderForm.tags"
                  placeholder="请输入漏洞名称"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="备案等级">
                <!-- 这里应根据你的数据结构使用适当的组件，例如el-select -->
                <el-select
                  class="typeSelect"
                  v-model="generateHeaderForm.filingLevel"
                  placeholder="请选择"
                  clearable
                  size="small"
                >
                  <el-option
                    v-for="item in zichanEnumObj.assetFilingLevel"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  >
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <!-- 其他表单项 -->
          </el-row>
        </el-form>
      </div>
    </el-dialog>
    <!-- 新增资产 -->
    <addAsset ref="addAsset" :value="addAssetisshow" @AddassetCancel="AddassetCancel" :id="updataeseetid"></addAsset>
    <!-- 发现资产 -->
    <findAsset ref="findAsset" :value="findAssetisshow" @findCancel="findCancel"></findAsset>
    <!-- 导入资产 -->
    <syncAsset ref="syncAsset" :value="syncAssetisshow" @importCancel="importCancel"></syncAsset>
    <!-- 发现资产 -->
    <importAsset ref="importAsset" :value="importAssetisshow" @importCancel="importCancel"></importAsset>
    <!-- 新增资产组 -->
    <addAssetGroup ref="addAssetGroup" :value="addAssetGroupisshow" @assetGroupCancel="assetGroupCancel" :id="updateGroupid"></addAssetGroup>
  </div>
</template>
<style lang="less" scoped>
.trafficlist {
  padding: 24px;
  background: #fff;
  min-height: calc(100% - 39px);
  box-sizing: border-box;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
/deep/ .el-table td.el-table__cell div {
  line-height: 20px;
}
/deep/ .el-tabs--border-card>.el-tabs__header .el-tabs__item.is-active {
    color: #4c7ae3 !important;
}
</style>
<script>
import { traffic } from '@/api/assetManagement.js'
import LeftZiChan from './components/LeftZiChan.vue'
import RightZiChan from './components/RightZiChan.vue'
import RightIPZiChan from './components/RightIPZiChan.vue'
import RightChar from './components/RightCharts.vue'
import Gailan from './components/ZiChanrRight/gailan.vue'
import ZIZIC from './components/ZiChanrRight/ZIZiChanG.vue'
import ZIChan from './components/ZiChanrRight/ZiChan.vue'
import ZIChanLD from './components/ZiChanrRight/ZiChanLD.vue'
import addAsset from './components/Addasset.vue'
import findAsset from './components/FindAsset.vue'
import syncAsset from './components/SyncAsset.vue'
import importAsset from './components/importAsset.vue'
import addAssetGroup from './components/addAssetGroup.vue'
export default {
  name: 'assetManagement',
  components: {
    LeftZiChan,
    RightZiChan,
    RightChar,
    Gailan,
    ZIZIC,
    ZIChan,
    ZIChanLD,
    RightIPZiChan,
    addAsset,
    findAsset,
    syncAsset,
    importAsset,
    addAssetGroup
  },
  data () { 
    return {
      activeTabName: 'one', // 默认激活第一个tab
      generateHeaderForm: {},
      headervisible: false,

      zichanList: [],
      zichanObj: {},
      nameTab: "",
      formData: {
        search_field: '',
        page_num: 1,

      },
      zichanEnumObj: {},
      tableData: [],
      multipleSelection: [],
      alldelvisible: false,
      Loading: false,
      pageSize: 10,
      totalpage: 0,
      currentpage: 1,
      showOperateButton: false,
      rowId: '',
      userID: '',
      addAssetisshow:false,
      findAssetisshow:false,
      allCheckedAsset:[],
      dropdowndisabled:true,
      importAssetisshow:false, 
      syncAssetisshow:false, 
      addAssetGroupisshow:false,
      selgroupids:[],
      selassetids:[],
      updataeseetid:'',
      updateGroupid:'',
    }
  },
  created () {
    this.$store.state.activefirstMenu = "/assetManagement"
    // this.userID = this.$commonjs.decryptCBC(localStorage.getItem('user_id'),this.$commonjs.myKey); 
  },
  mounted () {
    this.getZC()
    this.getZCEnum()
    console.log(this.$store.state.groupID[0])
    // this.$store.commit('setgroupID',0)
  },
  methods: {
    clickSOnTOTab (name) {
      this.nameTab = name
    },
    gotoFirstTab(groupID) {
      this.activeTabName = 'one'; // 设置为第一个tab的name,重新设置groupID
      this.$store.commit('setgroupID',groupID)
    },
    hanldeclicknode(){
        this.activeTabName = 'one'; 
    },
    // 高级搜索
  async  saveHeaderForm () {
    try {
        console.log(this.generateHeaderForm,'高级搜索');
        const res  = await traffic.trafficlistinfodel(this.generateHeaderForm)
        if (res.code == 200) {
        this.$refs.RightZiChan.data = this.$refs.RightZiChan.transformData(res.data.list)
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
        this.handleCloseHeader()
    } catch (error) {
        console.log(error);
    }
    },
    
    closeHeaderDialog () {
       this.resetHeaderFormFields();
    },
    handleCloseHeader(){
        this.headervisible = false
        this.resetHeaderFormFields();
    },
  // 弹窗关闭时重置表单字段
  resetHeaderFormFields() {
    Object.keys(this.generateHeaderForm).forEach(key => {
      if (typeof this.generateHeaderForm[key] === 'string') {
        this.generateHeaderForm[key] = ''; // 将字符串字段设置为空字符串
      } else if (typeof this.generateHeaderForm[key] === 'number') {
         if (key === 'riskLevel' || key === 'filingLevel') {
          this.generateHeaderForm[key] = null; // 或者你的下拉框默认值
        } else {
          this.generateHeaderForm[key] = 0; // 其他数字字段设置为0
        }
      }
    });
  },
  //右边 tab 切换
  handleClick(tab){ 
    let _name = tab.name
    if(_name == 'one'){ //概览
      this.$refs.gailan.getData();
    }
    if(_name == 'tw'){ //子资产组
      this.$refs.zizichanzu.getData();
    }
    if(_name == 'sa'){ //资产
      this.$refs.zichan.getData();
    }
    if(_name == 'fou'){ //资产漏洞
      this.$refs.zichanloudong.getData();
    }
  },
    // 枚举
    async getZCEnum () {
      const res = await traffic.trafficEnum()
      if (res.code == 200) {
        this.zichanEnumObj = res.data
        console.log(res, '枚举枚举枚举')
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    async getZC () {
      const res = await traffic.getStatus()
      if (res.code == 200) {
        this.zichanObj = res.data
        console.log(res, '3333')
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 搜锁-0普通走索
    async getData () {
      const res = await traffic.trafficlistinfodel({
        search: this.formData.search_field,
      })
      console.log(this.$refs.RightZiChan, 'RightZiChan')
      if (res.code == 200) {
        this.$refs.RightZiChan.data = this.$refs.RightZiChan.transformData(res.data.list)
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    btninfo (row) { //详情
      this.$router.push({
        path: `/trafficinfo`,
        query: {
          id: row.id,
          name: row.taskName,
          // page_num: this.formData.page_num,
          // risk_level: row.risk_level
        }
      })
      localStorage.setItem('trafficTab', 'tabs1')
    },

    handlesearch () {
      if(this.formData.search_field){
        this.getData()
      }else{
        this.$message.error("搜索条件不能为空")
      }

    },
    handleReset () {
      this.formData.search_field = ""
      this.getData()
    },
    async btnMultiDelete () {
      if (this.multipleSelection.length == 0) return
      var ids = []
      for (var i = 0; i < this.multipleSelection.length; i++) {
        ids.push(this.multipleSelection[i].id)
      }
      const res = await traffic.delTraffic({
        flowTaskIds: ids.join(','),
      })
      if (res.code == 200) {
        this.$message({
          message: '删除成功',
          type: 'success'
        })
        this.getData()
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }

    },
    async btnDel (scope) { //单个删除
      const res = await traffic.delTraffic({
        flowTaskIds: scope.row.id,
      })
      if (res.code == 200) {
        this.$message({
          message: '删除成功',
          type: 'success'
        })
        this.getData()
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    btnCancelDel (scope) {
      scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
    },
    handleSelectionChange (val) {
      this.multipleSelection = val
    },
    mouseenter (row) {
      this.showOperateButton = true
      this.rowId = row.id   //赋值行id，便于页面判断 
    },
    mouseleave (row) {

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

    },
    async btnTaskstop (id) { //结束

      const res = await traffic.trafficStatus({
        flowTaskId: id,
        operate: 'stop',
        userId: this.userID
      })
      if (res.code == 200) {
        this.$message({
          message: '任务结束',
          type: 'success'
        })
        this.getData()
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    currentchange (t) {
      this.formData.page_num = t
      this.getData()
      this.currentpage = t
    },
    handleSizeChange (t) {
      this.formData.page_num = 1
      this.pageSize = t
      this.getData()
    },
    // 新增资产
    handleCommand(command){
      if(command == 'addAsset'){
        this.addAssetisshow = true; 
      }
      if(command == 'addAssetGroup'){ //资产组
        this.addAssetGroupisshow = true; 
      }
      if(command == 'findAsset'){
        this.findAssetisshow = true; 
      }
      if(command == 'importAsset'){  
        this.importAssetisshow = true;
      }
      if(command == 'syncAsset'){  
        this.syncAssetisshow = true;
      }
      if(command == 'penetrationtest'){
        let checkedNodes = this.$refs.RightZiChan.checkedNodes(); 
 
        this.allCheckedAsset = [];  
        this.getTreeAssetLabel(checkedNodes); 
        if(this.allCheckedAsset.length>0){
          this.$router.push({
              path: `/createtask`,
              query: { 
                  flag: 1,  
                  type:3, //资产组 渗透任务
              }
          });
          localStorage.setItem("checkedasset", this.allCheckedAsset.join(','));
        } 
      }
      if(command =='selectedAsset'){ //选择资产
        this.$refs.RightZiChan.selectedAsset();
      }
      if(command == 'cancelSelected'){ //取消选择
        this.$refs.RightZiChan.cancelSelected();
        // this.$refs.tree.setCheckedKeys([]);
      }
      if(command == 'assetdelete'){ //批量删除
          let checkedNodes = this.$refs.RightZiChan.checkedNodes();
          this.selgroupids=[];
          this.selassetids=[];
          this.getTreeAssetID(checkedNodes); 
           
          this.handleassetdel();
      }
    }, 
    getTreeAssetLabel(list){ 
      for(var i=0;i<list.length;i++){
        if(list[i].type==2){
            this.allCheckedAsset.push(list[i].label)
        } 
      } 
    },
    getTreeAssetID(list){  
      for(var i=0;i<list.length;i++){
        let item = list[i];
        if(item.type == 1){
          this.selgroupids.push(item.id.split('_')[1])
        }else{
          this.selassetids.push(item.id.split('_')[1])
        } 
      }
    },
    findCancel(needProgress){
      this.findAssetisshow = false;
      this.$refs.RightZiChan.replaceData();
      if (needProgress) {
        this.$refs.RightZiChan.getTaskProgress();
      }
    },
    AddassetCancel(){
      this.addAssetisshow = false;
      this.$refs.RightZiChan.replaceData();
    },
    importCancel(){
      this.importAssetisshow = false;
      this.$refs.RightZiChan.replaceData();
    },
    assetGroupCancel(){
      this.addAssetGroupisshow = false;
      this.$refs.RightZiChan.replaceData();
    },
    async handleassetdel(){ 
      const res = await traffic.assetdelete({
        assetIds:this.selassetids.join(','),
        groupIds:this.selgroupids.join(',')
      });
      if(res.code ==200){
        this.$message({
          message: '批量删除资产组资产成功',
          type: 'success'
        })
        this.$refs.RightZiChan.replaceData();
      }else{
        this.$message({
          message: '批量删除资产组资产失败',
          type: 'erroe'
        })
      }
    },
    updateAsset(id){ //编辑资产
      
      this.updataeseetid = id;
      this.addAssetisshow = true;   
    },
    updateAssetGroup(id){ //编辑资产组
      this.updateGroupid = id;
      this.addAssetGroupisshow = true; 

    },
  },
}
</script>