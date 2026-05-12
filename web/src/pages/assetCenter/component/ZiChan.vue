<template>
  <div> 
    <div class="tasklist context_box_bg">
      <div class="search-box">
        <div  class="operationbutton"> 
          <el-button size="small" type="primary" @click="openAddDialog">新增</el-button>  
        <el-button   
            type="primary" 
            size="small"
            @click="showDialog"
            >导入 </el-button>

        <el-button   type="primary" size="small" @click="exportData">导出</el-button> 
        <!-- <el-button   type="primary" size="small" @click="selectCurrentPage">全选</el-button>
        <el-button   type="primary" size="small" @click="deselectCurrentPage">取消全选</el-button> -->
        <!-- <el-button
            style="margin-right: 10px;"
            type="primary"
            size="small"
            @click="btnMultiDeleteTask"
        >删除</el-button> -->

        <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" style="padding-left:8px"
            trigger="click" :visible-arrow="false" v-model="alldelvisible">
            <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
            <div style="text-align: right; margin: 0" class="">
              <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消</el-button>
              <el-button size="mini" type="primary" @click="btnMultiDeleteTask">确定</el-button>
            </div>
            <el-button type="warning" size="small" slot="reference"   style="margin-right: 10px;" :disabled="!multipleSelection.length">删除</el-button>
          </el-popover>
        <el-dropdown    size="small" split-button type="primary" @command="handleCommand"> 
            更多操作+
            <el-dropdown-menu slot="dropdown">
                <el-dropdown-item :disabled="!testIP"   command="a">渗透测试</el-dropdown-item> 
                 <el-dropdown-item   command="selectedall">全选</el-dropdown-item> 
                  <el-dropdown-item    command="noselectedall">取消全选</el-dropdown-item> 
            </el-dropdown-menu>
            
        </el-dropdown> 
 
        </div>
        <div class="serach-condition"> 
          <div class="search-text"  > 
            <!-- <label for="" style="color:#525966;font-size: 14px;">资产IP：</label> -->
            <el-input
              style="width: 140px;margin-right:10px;margin-bottom:5px"
              placeholder="资产IP"
              v-model="formData.assetIP"
              class="input-with-select"
              size="small"
              clearable
            >
            </el-input> 
            <el-input
            style="width: 140px;margin-right:10px;margin-bottom:5px"
              placeholder="端口"
              v-model="formData.port"
              class="input-with-select"
              size="small"
              clearable
            >
            </el-input>
            <el-input
            style="width: 140px;margin-right:10px;margin-bottom:5px"
              placeholder="服务"
              v-model="formData.service"
              class="input-with-select"
              size="small"
              clearable
            >
            </el-input>
            <el-input
            v-show="false"
            style="width: 140px;margin-right:10px;"
              placeholder="请输入搜索内容…"
              @keydown.enter="handlesearch"
              v-model="formData.search"
              class="input-with-select"
              size="small"
              clearable
            >
            </el-input>
            <el-button type="primary" size="small" @click="handlesearch"
              >搜索</el-button
            >

           <el-button type="primary" size="small" @click="handleReset">重置</el-button>
           <!-- <el-button v-if="!isshowHign" type="primary" size="small" @click="isshowHign=true">更多搜索</el-button>
           <el-button v-if="isshowHign" type="primary" size="small" @click="isshowHign=false; formData ={}">关闭更多搜索</el-button> -->


          </div>
          <div>
            <el-button
              type="primary"
              size="small"
              @click="handleReset"
              v-show="false"
              >重置</el-button
            >
          </div>
        </div>
      </div> 
      <el-table
        ref="myTable"
        :data="tableData"
        v-loading="Loading"
        @selection-change="handleSelectionChange"
        @cell-mouse-enter="mouseenter"
        @cell-mouse-leave="mouseleave"
        style="width: 100%"
         height="calc(100% - 112px)"
      >
        <el-table-column type="selection" width="40"  > </el-table-column>

        <el-table-column
          prop="assetName"
          label="资产IP"  width="120"
          :show-overflow-tooltip="true"
        >
          <template #default="scope">

             {{ scope.row.ip}}
          </template>
        </el-table-column>

        <!-- <el-table-column prop="assetType" label="资产类型" width="80"> </el-table-column> -->
      
        <el-table-column :show-overflow-tooltip="true" prop="openPort" label="开放端口"> </el-table-column>
        <el-table-column  prop="assetRiskName" label="风险 | 漏洞 " width="150">
            <template #default="scope">
                <span :style="getRiskLevelColor(scope.row.assetRiskName)">
                  {{ scope.row.assetRiskName }}
                </span> 
                <span style="color:red;margin-left:10px">{{scope.row.vulStatics.deadlyVul}}</span>
                <span style="color:#f9ca84;margin-left:10px">{{scope.row.vulStatics.highRiskVul}}</span>
                <span style="color:#1223ff;margin-left:10px">{{scope.row.vulStatics.mediumRiskVul}}</span>
                <span style="color:#8ac23a;margin-left:10px">{{scope.row.vulStatics.lowRiskVul}}</span>
              </template>
         </el-table-column>
        <el-table-column v-if="false" prop="vulInfo" label="漏洞">
               <template #default="scope">
                    <div style="display:flex;justify-content: space-between;">
                        <div>
                          <span style="color:red;margin-right:3px">{{scope.row.vulInfo.highRiskAsset}}</span>
                          <span style="color:#f9ca84;margin-right:3px">{{scope.row.vulInfo.middleRiskAsset}}</span>
                          <span style="color:#1223ff;margin-right:3px">{{scope.row.vulInfo.lowRiskAsset}}</span>
                          <span style="color:#8ac23a">{{scope.row.vulInfo.safeAsset}}</span>
                        </div>
                          <span style="color:#000">{{scope.row.vulInfo.totalRiskVul}}</span>
                    </div>
                </template>
         </el-table-column> 
        <el-table-column prop="assetGroupName" label="资产组"  :show-overflow-tooltip="true"> </el-table-column>
         <!-- <el-table-column   prop="isCloudHost" label="云主机" width="70"  >
               <template #default="scope">
                    <div style="display:flex;justify-content: space-between;"> 
                          <span style="color:#000">{{scope.row.isCloudHost ? '是':'否'}}</span>
                    </div>
                </template>
         </el-table-column> -->
        <el-table-column  :show-overflow-tooltip="true" prop="system" label="系统" > </el-table-column>
        <el-table-column :show-overflow-tooltip="true"  prop="testTime" label="更新时间"> </el-table-column>
        <el-table-column  prop="" label="操作" width="160" >
          <template #default="scope">

                <el-link 
                  class="link_primary"
                  :underline="false"
                  @click="openEditDialog(scope.row)"
                  >编辑</el-link
                >
                <el-link
                  class="link_primary"
                  :underline="false"
                  style="margin: 0 5px "
                  @click="btnTaskinfo(scope.row)"
                  >详情</el-link
                >
                <!-- <el-link
                style="margin-left: 0px;"
                  class="link_primary"
                  type="primary"
                  :underline="false"
                  @click="btnDel(scope.row)"
                  >删除</el-link
                > -->
                <el-popover placement="bottom" width="170" :visible-arrow="false" :ref="`popover_id-${scope.row.assetID}`"
                    popper-class="delButton_popper">
                    <p class="delText">
                      <i class="el-icon-warning"></i>确定删除吗？
                    </p>
                    <div style="text-align: right; margin: 0">
                      <el-button size="mini" class="delCancel" @click="btnCancelDel(scope)">取消</el-button>
                      <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                    </div>
                    <!-- <span slot="reference">删除</span> -->
                    <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference">删除
                    </el-link>
                  </el-popover>
          </template>
        </el-table-column>


      </el-table>
      <div class="pagination-container">
        <el-pagination
          :page-size="pageSize"
          background
          layout="total, prev, pager, next, sizes, jumper"
          :total="totalpage"
          :current-page="currentpage"
          @current-change="currentchange"
          @size-change="handleSizeChange"
          class="pagination"
        >
        </el-pagination>
  </div>
    </div>
        <!-- 新增/编辑对话框 -->
        <el-dialog @close="closeDia" destroy-on-close width="800px" :append-to-body="false"
            :title="dialogTitle" :visible.sync="dialogVisible"> 
            <el-form 
              label-width="120px" 
              label-position="left" 
              ref="formAdd" :model="formAdd"   
              
              :rules="rules" style="padding: 20px;">
                <el-form-item  label="资产IP/域名 " prop="ip">
                    <el-input  v-model="formAdd.ip"></el-input>
                </el-form-item>
                <!-- <el-form-item  label=" " >
                当前选择资产组：{{ fatherNodelabel }}
                </el-form-item> -->
                <el-form-item v-if="false" label="归属资产组 " > 
                    <el-tree
                        ref="treeRefs"
                        :data="filteredData" 
                        node-key="id"
                        :props="defaultProps"
                        :expand-on-click-node="false" >
                        <template
                        #default="{ node, data }">
                            <div class="tree-node" @mouseenter="node.hover = true" @mouseleave="node.hover = false">
                                <div style="color: #fff;position: relative;width:200px;" @click="fatherNode(node, data)"
                                >{{ node.label }}

                                </div
                                >
                            </div>
                        </template>
                    </el-tree>
                </el-form-item>
                <el-form-item  prop="selectedNode" label="归属资产组">
                <el-select style="width:100%" v-model="selectedNode" placeholder="请选择资产组" @visible-change="handleVisibleChange">
                <el-option
                    v-if="false"
                    :key="0"
                    :value="null"
                    label=""
                />
                <el-tree
                    ref="treeRefs"
                    :data="filteredData"
                    node-key="id"
                    :props="defaultProps"
                    :expand-on-click-node="false"
                    :default-expanded-keys="defaultExpandedKeys"
                    class="tree-select"
                >
                    <template #default="{ node, data }">
                    <el-option
                        :key="data.id"
                        :label="node.label"
                        :value="data.id"
                        @click="handleNodeClick(node, data)"
                    >
                        <span class="tree-node">{{ node.label }}</span>
                    </el-option>
                    </template>
                </el-tree>
                </el-select>
            </el-form-item>
        <el-form-item  label="操作系统 ">
          <el-input  v-model="formAdd.opSys"></el-input>
        </el-form-item>
        <el-form-item  label="资产名称 ">
          <el-input  v-model="formAdd.name" ></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="IP地址段 ">
          <el-input  v-model="formAdd.ipSegment"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="资产类型 ">
          <!-- <el-input   v-model="formAdd.assetType"></el-input> -->
          <el-select  clearable  style="width: 100%;margin-right:10px;margin-bottom:5px" v-model="formAdd.assetType" placeholder="资产类型">
            <el-option
              v-for="item in enumobj.assetType"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>

        <el-form-item v-if="false" label="基础软件名称 ">
          <el-input  v-model="formAdd.baseSoftwareName"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="基础硬件名称 ">
          <el-input  v-model="formAdd.baseHardwareName"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="基础软件版本 ">
          <el-input  v-model="formAdd.baseSoftwareVersion"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="资产用途 ">
          <el-input  v-model="formAdd.purpose"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="设备形态 ">
          <el-input  v-model="formAdd.equipmentForm"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="部署位置 ">
          <el-input  v-model="formAdd.deploymentLocation"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="等保级别 " style="">
          <el-select clearable style="width: 100%;height: 26px;"   v-model="formAdd.equalProtectionLevel" placeholder="等保级别">
            <el-option
              v-for="item in enumobj.fillingLevel"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>

        <el-form-item v-if="false" label="系统名称 ">
          <el-input  v-model="formAdd.systemName"></el-input>
        </el-form-item>
        <el-form-item label="管理员 ">
          <el-input  v-model="formAdd.systemAdmin"></el-input>
        </el-form-item>
         <el-form-item   label="部门 ">
          <el-input  v-model="formAdd.responsibleDepartment"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="运维负责人 ">
          <el-input  v-model="formAdd.systemOp"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="可信设备 " style="">
          <el-select clearable style="width: 100%;height: 26px;"   v-model="formAdd.trustLevel" placeholder="可信设备">
            <el-option
              v-for="item in enumobj.trustLevel"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="false" label="设备权重 " style="">
          <el-select clearable style="width: 100%;height: 26px;"   v-model="formAdd.deviceWeight" placeholder="设备权重">
            <el-option
              v-for="item in enumobj.deviceWeight"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="标签 ">
          <el-input  v-model="formAdd.tags"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="是否云主机 ">
          <el-select clearable style="width: 100%;height: 26px;" v-model="formAdd.isCloudHost" placeholder="请选择是否云主机">
            <el-option  label="是" :value="1" ></el-option>
            <el-option  label="否" :value="0" ></el-option>
          </el-select>
        </el-form-item>
         <el-form-item v-if="false" label="登录用户名">
          <el-input  v-model="formAdd.user"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="登录密码">
          <el-input type="password" v-model="formAdd.password"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="登录端口">
          <el-input type="number" v-model="formAdd.port"></el-input>
        </el-form-item>
        <el-form-item v-if="false" label="连接协议 ">
          <el-select   style="width: 100%;height: 26px;" v-model="formAdd.protocol" >
            <el-option v-for="(item,i) in enumobj.loginProtocol" :key="i" :label="item.label" :value="item.value" ></el-option>
             
          </el-select>
        </el-form-item> 
            </el-form>
            <template v-slot:footer>
                <el-button   @click="handleClose()">取消</el-button>
                <el-button type="primary" @click="handleSubmit()">确定</el-button>
            </template>
        </el-dialog>

        <el-dialog @close="closeDia" destroy-on-close width="600px"   title="批量修改资产组" :visible.sync="dialogVisibleWW">
            <el-form label-width="120px" label-position="left" ref="formAdd" :model="formAdd" :rules="rules1">
                <el-form-item  prop="selectedNode" label="归属资产组">
                    <!-- {{fatherNodeID}}--{{ selectedNode }} -->
                    <el-select style="width:100%" v-model="selectedNode" placeholder="请选择资产组" @visible-change="handleVisibleChange">

                        <el-tree
                            ref="treeRefs"
                            :data="filteredData"
                            node-key="id"
                            :props="defaultProps"
                            :expand-on-click-node="false"
                            :default-expanded-keys="defaultExpandedKeys"
                            class="tree-select" >
                            <template #default="{ node, data }">
                            <el-option
                                :key="data.id"
                                :label="node.label"
                                :value="data.id"
                                @click="handleNodeClick2(node, data)" >
                                <span class="tree-node">{{ node.label }}</span>
                            </el-option>
                            </template>
                        </el-tree>
                    </el-select>
                </el-form-item> 
            </el-form>
            <template #footer>
                <el-button type="primary" plain @click="dialogVisibleWW = false">取消</el-button>
                <el-button type="primary" @click="BatchAssetGroups">确定</el-button>
            </template>
        </el-dialog>
        <!-- 资产详情-->
        <el-dialog @close="closeDia" destroy-on-close width="1184px"  title="资产信息" :visible.sync="dialogVisible2"> 
          <RightIPZiChan ref="rightIPZiChan" :selectType="selectType" :detailOBJ= "detailOBJ" v-if="dialogVisible2" @tab-change="handleTabChange" />
        </el-dialog>
        <el-dialog :visible.sync="dialogVisibleD" title="导入资产">
            <div style="padding: 20px;">
                <p style="color: #f06c0b;">导入资产表中的IP不能为空，资产组为空时，资产将导入默认组</p>
                <div>
                <el-button style="margin: 20px 0" size="small" type="primary"  @click="downloadTemplate">下载资产模板</el-button>
                </div>
                <el-upload
                ref="excelUploadRef"
                style="display: inline-block; "
                action="#"
                :show-file-list="false"
                :before-upload="beforeUpload"
                :http-request="upload"
                :limit="111"
                :on-exceed="handleExceed"
                accept=".csv,.docx,.xlsx"
                >
                <el-button
                    style="margin-top: 10px"
                    type="primary"
                    class="btn-blue"
                    size="small"
                >
                    上传资产表
                </el-button>
                </el-upload>
            </div>
      <template #footer>
        <el-button @click="dialogVisibleD = false">关闭</el-button>
        <!-- <el-button type="primary" @click="dialogVisibleD = false">确定</el-button> -->
      </template>
    </el-dialog>
  </div>
</template>
<style lang="less" scoped>
:deep(.el-table-column--selection .cell){
    padding-left: 12px !important;
}
:deep(.el-button){
        vertical-align: top !important;
}
.pagination-container {
  display: flex;
  justify-content: flex-end; /* 居中对齐 */
  // margin-top: 20px;
}
 ::v-deep .el-input__wrapper{
  height: 26px;
}
.operation-buttons {
  position: relative;
}
.operation-buttons .el-link {
  margin-right: 8px;
}
.show-buttons .el-link {
  display: inline-block;
}
.tasklist {
  padding:0 0 20px 4px;
  height: calc(100% );
  width: 100%;
  box-sizing: border-box;
  // box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
 thead {
  .cursorPointer {
    cursor: pointer;
    position: absolute;
    &.active {
      color: #4c7ae3;

      i {
        color: #4c7ae3;
      }
    }
  }

  .cell {
    line-height: 15px;

    > span {
      position: absolute;
    }
  }

  .iconfont {
    color: rgba(72, 72, 102, 0.32);
    margin-left: 5px;
  }

  .el-select {
    height: 0;
    visibility: hidden;

    .el-input,
    .el-input__inner {
      height: 0 !important;
    }
  }
}
</style>
<script> 
import asset from '@/api/asset.js';  
import axios from 'axios';
import * as  XLSX from  'xlsx';
import RightIPZiChan from './RightIPZiChan.vue'
export default {
  name: 'x-ray',
  components: {
    RightIPZiChan,
  },
  data () {
    return {
      formKey: 0,
      dialogVisibleWW:false,
      isQuanXuan:false,
      defaultExpandedKeys: [], // 存放要默认展开的节点ID
      selectedNode: '',
      selectedNode2: '',
      dialogVisible2:false,
      dialogVisibleD:false,
      fatherNodeID :'',
      fatherNodelabel:"",
        allids:[],
        ttredata:[],
        defaultProps : {
        children: 'children',
        label: 'label',
      },
      filteredData:[],
      isshowHign:false,
      formAdd: {
        id:0,
        ip: '',
        assetGroupID: '', 
        systemAdmin:'', 
        name:'',
        opSys:'',
        responsibleDepartment:'',
        tags:'',
        selectedNode:'',
      }, 
      enumobj:{},
      editRowData:{},
      fileList: [],
      createDictionaries: false,
      timeIDFOrSHouQuan: null,
      alldelvisible: false,
      Loading: false,
      formData: {
        page_num: 1,
        task_status: "",
        starttime: "",
        stoptime: "",
        search_field: "",
        time: "",
        risk_level: 0,
        isCloudHost:'',
      },
      uploadFormData: {},
      tableData: [],
      multipleSelection: [],
      totalpage: 0,
      currentpage: 1,
      risklevellist: [],
      pageSize: 10,
      showOperateButton: false,
      rowId: '',
      timer: null,
      dialogVisible: false,
      dialogTitle:'',
      treedata:[],
      exportIDS:'',
      testIP:'',
      detailOBJ:{},
      selectType:1,
      rememberselect:false,
      rememberselectArr:[],
      currentAssetRow: null, // 存储当前查看详情的资产行数据
      rules: {
        ip: [
          { required: true, message: '资产IP/域名不能为空', trigger: 'blur' }
          // { pattern: /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/, message: '请输入有效的IP地址', trigger: 'blur' }
        ],
        selectedNode: [
          { required: true, message: '归属资产组不能为空', trigger: 'change' }
        ]
      },
      rules1: {
        ip: [
          { required: true, message: '资产IP/域名不能为空', trigger: 'blur' }
          // { pattern: /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/, message: '请输入有效的IP地址', trigger: 'blur' }
        ],
        selectedNode: [
          { required: true, message: '归属资产组不能为空', trigger: 'change' }
        ]
      },
    }
  },
  created () {
    // this.$store.state.activefirstMenu = '/x-ray'
    this.getData3()
    this.emumOptions()
    this.getZCList()
  },

  mounted() {
    let _this = this; 
    // 初次加载时获取数据
    this.getData();

    // 使用自定义事件监听 groupID 变化
    window.addEventListener("groupIDChanged", function (e) {
        _this.getData();
    });
},


  methods: {
    gotoTest2(){
         console.log(2)
        let _this = this
        this.$router.push({
            path: `/createBaselineTask`,
            query: {
              TestObjectives: _this.testIP,
            }
          })
      },
    gotoTest(){
        console.log(1)
      let _this = this
      this.$router.push({
          path: `/createtask`,
          query: {
            TestObjectives: _this.testIP,
            type:3,
          }
        })
    },
    handleNodeClick(node, data) {
      this.selectedNode = data.id;
    },
    handleNodeClick2(node, data) {
      this.selectedNode2 = data.id;
    },
    handleVisibleChange(visible) {
      if (visible) {
        this.$refs.treeRefs.updateKeyChildren(null, this.filteredData);
      }
    },
    showDialog() {
      this.dialogVisibleD = true;
    },
   async BatchAssetGroups(){
    if(!this.selectedNode2) return this.$message({
                    message: '请选择资产组',
                    type: 'error'
                });
      const arrMu = this.multipleSelection.map(item=>item.assetID)
      console.log(this.selectedNode2,arrMu);

      const assetid = Number(this.selectedNode2.split('_')[1]) 
      try {
        const res = await asset.updateManyAssetGroups({
          "assetIds":arrMu, //选择的资产id,int数组，必须
          "assetGroupId":assetid//所属资产组id,int，必须
        })
        // const res = await asset.updateManyAssetGroups(formData);
        console.log(res,'rrrrrrr');
      this.getData3()
        this.dialogVisibleWW = false;
      } catch (error) {
        console.log(error);
        this.dialogVisibleWW = false;
      }
    },
    downloadTemplate() {
      const url = '/AssetManagementTemplate.xlsx';
      const link = document.createElement('a');
      link.href = url;
      link.download = 'AssetManagementTemplate.xlsx';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
     fatherNode (node, data)  {
      let id = data.id.split('_')[1];
      this.fatherNodeID = id
      this.fatherNodelabel = data.label
    },
     transformData (data, level = 1)  {
      if (!Array.isArray(data)) {
        return []; // 返回空数组作为兜底
      }
      let unknownIdCounter = 0;
      return data.map(item => ({
        id: typeof item === 'object' && item.type && item.id ? `${item.type}_${item.id}` : `unknown_${unknownIdCounter++}`,
        label: typeof item === 'object' && item.name ? item.name : '未命名',
        type: typeof item === 'object' && item.type ? item.type : '未知类型',
        level: level,
        children: item.items && Array.isArray(item.items) ? this.transformData(item.items, level + 1) : [],
        hover: false, // 用于控制按钮显示
      }));
    },
     getTreeAssetID  (list)  {
      for (var i = 0; i < list.length; i++) {
        this.allids.push(list[i].id);
        if (list[i].children && list[i].children.length > 0) {
          this.getTreeAssetID(list[i].children);
        }
      }
    },
    async  getZCList () {
      const res = await asset.trafficlistinfodel();
      if (res.code == 200) {
        this.ttredata = this.transformData(res.data.list);
        this.filteredData = this.ttredata;
        this.defaultExpandedKeys = this.getAllNodeIds(this.filteredData);
        this.getTreeAssetID(this.ttredata);
      } else {
         this.$message({
          message: "获取资源失败",
          type: "error"
        })
        
      }
    },
    getAllNodeIds(data) {
      let ids = [];
      for (let item of data) {
        ids.push(item.id);
        if (item.children) {
          ids = ids.concat(this.getAllNodeIds(item.children));
        }
      }
      return ids;
    },
    async emumOptions(){
      const res = await asset.enumApi()
      if(res.code == 200){
        this.enumobj = res.data
      }
    },
    async openAddDialog  ()  {
      this.formKey++; // 改变 key，强制重建 el-form

      this.dialogTitle = '新增资产'; 
      // await fetchEnumOptions();
      this.dialogVisible = true; 
      // localStorage.getItem("groupID") 
      this.formAdd.selectedNode = localStorage.getItem("groupID")

      this.selectedNode = '';
      // this.selectedNode = '1_'+localStorage.getItem("groupID");

      this.getZCList();

    },
    handleClose(){

      // 刷新当前页面
      this.$router.go(0)
      this.$refs.formAdd.resetFields()
      this.dialogVisible = false
      this.formAdd.selectedNode ='';
      this.formAdd.assetGroupID = '';
      this.formAdd.id =0;
      this.formAdd.ip = '';
      this.formAdd.name= '';
      this.formAdd.opSys = '';
      this.formAdd.systemAdmin = '';
      this.formAdd.responsibleDepartment =  '';
      this.formAdd.tags = '';
      this.selectedNode = '';
    },
    //保存 资产
  async handleSubmit(){ 
    try { 
      this.$refs.formAdd.validate( async(valid) => {
      if (valid) {
        // 提交表单逻辑
      
        if(this.dialogTitle == '新增资产'){
          this.formAdd.assetGroupID =  this.fatherNodeID  -0
          // this.formAdd.assetType =  this.formAdd.assetType  -0
          // this.formAdd.equalProtectionLevel =  this.formAdd.equalProtectionLevel  -0
          // this.formAdd.port = Number(this.formAdd.port)
          // this.formAdd.protocol = Number(this.formAdd.protocol)
          // this.formAdd.ip = this.formAdd.ip;
          // this.formAdd.name=this.formAdd.name;
          // this.formAdd.opSys = this.formAdd.opSys;
          // this.formAdd.systemAdmin = this.formAdd.systemAdmin;
          // this.formAdd.responsibleDepartment = this.formAdd.responsibleDepartment;
          // this.formAdd.tags = this.formAdd.tags;
            const res = await asset.addasset(this.formAdd,1)
              if (res.code == 200) {
                 this.$message({
                  message: "新增资产成功",
                  type: "success"
                })
                this.getData()
                this.dialogVisible = false;
                // 通知树组件去请求
                // 生成随机数1-10万之间
                //   let randomNum = Math.floor(Math.random() * 100000) + 1;
                // localStorage.setItem("refeshTree",randomNum)
              }else{
                    this.$message.error(res.msg);  
              }
        }else{
            // console.log( this.formAdd.assetGroupID,'bianji-----',this.fatherNodeID);
            this.formAdd.assetGroupID = this.fatherNodeID? this.fatherNodeID  -0 :this.formAdd.assetGroupID
            // this.formAdd.assetType =  this.formAdd.assetType  -0
            // this.formAdd.equalProtectionLevel =  this.formAdd.equalProtectionLevel  -0
            // this.formAdd.trustLevel =  this.formAdd.trustLevel  -0
            // this.formAdd.deviceWeight =  this.formAdd.deviceWeight  -0
            // this.formAdd.id = this.formAdd.id;
            // // this.formAdd.port = Number(this.formAdd.port);
            // // this.formAdd.protocol = Number(this.formAdd.protocol)
            // this.formAdd.ip = this.formAdd.ip;
            // this.formAdd.name = this.formAdd.name;
            // this.formAdd.opSys = this.formAdd.opSys;
            // this.formAdd.systemAdmin = this.formAdd.systemAdmin;
            // this.formAdd.responsibleDepartment = this.formAdd.responsibleDepartment;
            // this.formAdd.tags = this.formAdd.tags;
            // console.log('bianji save -----', this.formAdd)
            const res = await asset.addasset(this.formAdd,2)
            if (res.code == 200) {
               this.$message({
                message: "编辑成功",
                type: "success"
              })
              this.getData()
              this.dialogVisible = false;
              // 刷新当前页面
              this.$router.go(0)
              // 通知树组件去请求
              //  let randomNum = Math.floor(Math.random() * 100000) + 1;
              // localStorage.setItem("refeshTree",randomNum)
            }else{
              this.$message({
                message: res.msg,
                type: "error"
              })
            }
          }
          // this.dialogVisible = false;
          // 刷新
          // this.$router.go(0)
      } else {
        // 弹窗提示：表单验证失败
         this.$message.error('请输入必填项');  
        return false;
      }
    });



  } catch (error) {
    console.log(error);
    this.dialogVisible = false;

  }
},
closeDia(){
  this.formAdd = {}
  // 重置详情页面相关数据
  if (this.dialogVisible2 === false) {
    this.selectType = 1;
    this.currentAssetRow = null;
    this.detailOBJ = {};
  }
},
selectCurrentPage() {
    this.isQuanXuan = true
    this.$refs.myTable.clearSelection();
    this.tableData.forEach(row => {
      this.$refs.myTable.toggleRowSelection(row, true);
    });
    this.getAllDataIP()
  },
  deselectCurrentPage() {
    this.isQuanXuan = false
    this.$refs.myTable.clearSelection();
    this.rememberselect = false
    this.rememberselectArr = []
    this.multipleSelection = []
    this.exportIDS = ''
    this.exportIPs = ''
    this.testIP = ''
  },
    async openEditDialog  (data) {
      try {
        const res = await asset.assetDetail({id:data.assetID,selectType:4})
        this.dialogTitle = '编辑资产';
        console.log(this.formAdd,'--------',res.data)
        this.dialogVisible = true;
        // this.editRowData = data;
        // this.formAdd = res.data.manageInfo

        this.formAdd.id =  data.assetID
        this.formAdd.selectedNode = res.data.manageInfo.assetGroupID;
        this.formAdd.assetGroupID = res.data.manageInfo.assetGroupID;
        this.selectedNode = '1_'+  this.formAdd.selectedNode ;

        // this.formAdd.trustLevel =  res.data.manageInfo.trustLevel
        this.formAdd.name =  res.data.manageInfo.assetName; 
        this.formAdd.opSys = res.data.manageInfo.opSys;
        this.formAdd.systemAdmin = res.data.manageInfo.systemAdmin;
        this.formAdd.responsibleDepartment = res.data.manageInfo.responsibleDepartment;
        this.formAdd.tags = res.data.manageInfo.tags;
        this.formAdd.ip = res.data.manageInfo.ip;
        this.formAdd.systemAdmin = res.data.manageInfo.systemAdmin
        


        console.log(this.formAdd)
        } catch (error) {
          console.log(error);
        }
    },

    handleStorageChange(event) {
      if (event.key === this.storageKey) {

        this.storageValue = event.newValue;
        this.getData()
      }
    },
      getRiskLevelColor(levelRiskName) {
        const riskColors = {
          '高危': '#ff0000',
          '中危': '#ffa500',
          '低危': '#1223fe',
          '安全': '#8ac23a',
          '未知': '#808080'
        };
        return { color: riskColors[levelRiskName] || 'defaultColor' }; // 如果没有匹配的风险等级，返回默认颜色
      },
    cancelform () {
      // alert(1)
    },
    beforeUpload(file) {
      // 获取文件名并转换为小写以忽略大小写差异
      // const isJson = file.name.toLowerCase().endsWith('.csv');
      // const isDocx = file.name.toLowerCase().endsWith('.docx');

      // // 检查文件类型
      // if (!isJson && !isDocx) {
      //   this.$message.error('只能上传 .csv 文件');
      //   return false;
      // }

      // // 可以根据需要增加更多的检查，例如文件大小等
      // return true;
    },

    handleFormDate (data) {
      console.log(data, 'data-111')
      const formDate = new FormData()
      formDate.append('taskName', data.name)
      formDate.append('file', data)
      formDate.append('file2', '11111')
      return formDate
    },
    // --上传
    async upload (obj) {
      console.log(obj, 'obj750----');
      let that_ = this
      let formData = new FormData() //使
      formData.append('file', obj.file)
      try { 
          asset.importAsset({
            file:obj.file
          }).then((dt) => {
              // 处理响应
            if (dt.code === 200) {
              this.$message.success('上传成功');
              that_.dialogVisibleD = false;
              that_.getData()
            } else {
              this.$message.error(`上传失败: ${dt.msg}`);
            }
          }) 

      } catch (error) {
        // 处理错误
        this.$message.error('上传失败');
        console.error(error);
      }finally {
        if (this.$refs.excelUploadRef) {
          this.$refs.excelUploadRef.clearFiles();
        }
      }

    },

  async exportData() {
    try {
      // 调用 API 获取要导出的数据
      const response = await asset.exportAsset({
        // 传递所需的参数，例如搜索条件等
        ...this.formData,
        ids:this.isQuanXuan? 'all' : this.exportIDS

      });

      if (response.code === 200) {
        // 将数组数据转换为工作簿
        const worksheet = XLSX.utils.json_to_sheet(response.data.assetsInfo);
        const workbook = XLSX.utils.book_new();
        XLSX.utils.book_append_sheet(workbook, worksheet, 'Sheet1');

        // 将工作簿转换为二进制数据
        const excelBuffer = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });

        // 创建 Blob 对象
        const blob = new Blob([excelBuffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });

        // 创建下载链接
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;

        // 设置下载文件名
        link.download = 'exported_data.xlsx';

        // 触发下载
        document.body.appendChild(link);
        link.click();

        // 移除下载链接
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);

         this.$message.success('数据导出成功');
      } else {
         this.$message.error(`导出失败: ${response.msg}`);
      }
    } catch (error) {
      console.error(error);
    }
  },
  async exportData1(){
    try {
      // 调用 API 获取要导出的数据
      const response = await asset.exportAsset({
        // 传递所需的参数，例如搜索条件等
        ...this.formData,
        ids:this.isQuanXuan? 'all' : this.exportIDS,
        exportType:2 //特殊导出
      });

      if (response.code === 200) {
        // 将数组数据转换为工作簿
        const worksheet = XLSX.utils.json_to_sheet(response.data.assetsInfo);
        const workbook = XLSX.utils.book_new();
        XLSX.utils.book_append_sheet(workbook, worksheet, 'Sheet1');

        // 将工作簿转换为二进制数据
        const excelBuffer = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });

        // 创建 Blob 对象
        const blob = new Blob([excelBuffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });

        // 创建下载链接
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;

        // 设置下载文件名
        link.download = 'exported_data.xlsx';

        // 触发下载
        document.body.appendChild(link);
        link.click();

        // 移除下载链接
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);

         this.$message.success('数据导出成功');
      } else {
         this.$message.error(`导出失败: ${response.msg}`);
      }
    } catch (error) {
      console.error(error);
    }
  },

    handleRemove (file, fileList) {
      console.log(file, fileList)
    },
    handlePreview (file) {
      console.log(file)
    },
    handleExceed (files, fileList) {
       this.$message.warning(`当前限制选择 3 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`)
    },
    beforeRemove (file, fileList) {
      return confirm(`确定移除 ${file.name}？`)
    },

    //  新增取消
    handleSaveCancel () {
      this.createDictionaries = false
    },

    async getData () {
        const id =  localStorage.getItem("groupID")
      try {
        const res = await asset.trafficvulndel({
          page: this.formData.page_num,
          size: this.pageSize,
          search: this.formData.search_field,
          groupID:id? id:'',
          ...this.formData,
        })
        if (res.code == 200) {
          this.tableData = res.data.assetsInfo 
          this.totalpage = res.data.count
        } else {
           this.$message({
            message: res.msg,
            type: "error"
          })
        }
      } catch (error) {
        console.log(error)
      }
    },
    async getAllDataIP () {
      const id =  localStorage.getItem("groupID")
      try {
        const res = await asset.getAllDataIP({
          search: this.formData.search_field,
          groupID:id? id:'',
          ...this.formData,

        })
        if (res.code == 200) {
         console.log(res.data, 'res.data')
        //  返回值：res.data.assetIDs和res.data.assetIPs
        // 给删除、导出、更多操作用
        this.exportIDS = res.data.assetIDs
        this.exportIPs = res.data.assetIPs
        this.multipleSelection = res.data.assetIPs
        this.isQuanXuan = true
        this.testIP = res.data.assetIPs.join(',')
        this.rememberselect = true
        this.rememberselectArr = res.data.assetIPs
        } else {
           this.$message({
            message: res.msg,
            type: "error"
          })
        }
      } catch (error) {
        console.log(error)
      }
    },

    async getData3 () {
      try {
        const res = await asset.trafficvulndel({
          page: 1,
          size: 10,
        })
        if (res.code == 200) {
          this.tableData = res.data.assetsInfo
          this.totalpage = res.data.count
        } else {
           this.$message({
            message: res.msg,
            type: "error"
          })
        }
      } catch (error) {
        console.log(error)
      }
    },

    handlesearch () {
      //搜索
      this.formData.page_num = 1
      this.getData()
      this.currentpage = 1
    },
         handleReset(){
          this.formData ={};

            this.formData.page_num = 1;
            this.currentPage = 1;
            localStorage.removeItem('groupID')
            this.getData();
        },

  async btnTaskinfo (row) { //详情
    // 详情接口
    try {
      this.currentAssetRow = row; // 存储当前行数据
      this.selectType = 1; // 确保有默认值，避免undefined
      const res = await asset.assetDetail({
        id: row.assetID,
        selectType: this.selectType,
        page: 1,
        size: 10
      })
      this.dialogVisible2 = true // 详情页面
      this.detailOBJ = res.data
    } catch (error) {
      console.log(error);
    }
    },

    async handleTabChange(newSelectType, page = 1, size = 10) {
      // 处理子组件的tab切换事件
      console.log('handleTabChange called with:', { newSelectType, page, size, currentType: this.selectType });

      this.selectType = newSelectType;
      if (this.currentAssetRow && newSelectType) {
        try {
          // 设置loading状态
          this.$refs.rightIPZiChan && this.$refs.rightIPZiChan.setLoading(newSelectType, true);

          const res = await asset.assetDetail({
            id: this.currentAssetRow.assetID,
            selectType: newSelectType,
            page: page,
            size: size
          });
          this.detailOBJ = res.data;

          // 更新分页信息
          this.$refs.rightIPZiChan && this.$refs.rightIPZiChan.updatePagination(res.data);
        } catch (error) {
          console.log(error);
        } finally {
          // 关闭loading状态
          this.$refs.rightIPZiChan && this.$refs.rightIPZiChan.setLoading(newSelectType, false);
        }
      }
    },



    async btnMultiDeleteTask () {
      if (this.multipleSelection.length == 0) return
      let _ids = this.multipleSelection.map(item => item.assetID)

      let msg = '你确定要删除所选数据？';

      if(_ids.length===0){
           this.$message({
              message: "请选择要删除的数据。",
              type: "success"
            })
          return
      }
     
      const res = await asset.deleteAsset({
          assetIds: _ids.join(",")
        })
        if (res.code == 200) {
            this.$message({
            message: "删除资产成功",
            type: "success"
          })
          this.currentpage = 1
          this.formData.page_num = 1
          this.alldelvisible = false
          this.getData()
        } else {
            this.$message({
            message: res.msg,
            type: "error"
          })
        }


     
    },
    // 取消删除
    btnCancelDel(scope){
      scope._self.$refs[`popover_id-${scope.row.assetID}`].doClose();
     
    },
    async btnDel (scope) { 
        const res = await asset.deleteAsset({ assetIds: scope.row.assetID + '' })

        if (res.code == 200) {
            this.$message({
            message: "删除资产成功",
            type: "success"
          })
          scope._self.$refs[`popover_id-${scope.row.assetID}`].doClose();
          this.getData();
        } else {
          
          this.$message({
              message: res.msg,
              type: 'error'
          });
        }  
    },
    // 取消删除 
    handleSizeChange (t) {
      this.formData.page_num = 1
      this.pageSize = t
      this.getData()
    },
    currentchange (t) {
      this.formData.page_num = t
      this.getData()
      this.currentpage = t
      if(this.rememberselect){
       setTimeout(() => {
        this.getAllDataIP()
        // 要全部选中
        console.log(this.rememberselectArr,'this.rememberselectArr');

        this.multipleSelection = this.rememberselectArr
        this.isQuanXuan = true
        this.testIP = this.rememberselectArr.join(',')
        // this.exportIDS = this.rememberselectArr.map(item => item.assetID).join(",")
        this.exportIPs = this.rememberselectArr.join(",")
        this.$refs.myTable.clearSelection();
        this.tableData.forEach(row => {
          this.$refs.myTable.toggleRowSelection(row, true);
        });
       }, 600);
      }
    },
    handleSelectionChange (val) {
      this.multipleSelection = val
      this.exportIDS = val.map(item => item.assetID).join(",")
      this.testIP = val.map(item => item.ip).join(",")


      console.log(val,'资产IP资产IP',this.multipleSelection);
    },
    mouseenter (row, colum, cell, event) {
      this.showOperateButton = true
      this.rowId = row.id   //赋值行id，便于页面判断
    },
    mouseleave (row, colum, cell, event) {
      this.showOperateButton = false;
      this.rowId = '';
      // if (row.executeType == 3 || row.status != 1) {
      //   if (!this.$refs['popover-' + row.id]) {
      //     this.showOperateButton = false
      //     this.rowId = ""
      //     return
      //   } else {
      //     let isShow = this.$refs['popover-' + row.id].showPopper
      //     if (!isShow) {
      //       this.showOperateButton = false
      //       this.rowId = ""
      //     }
      //   }
      // } else {
      //   if (!this.$refs['popover_id-' + row.id]) {
      //     this.showOperateButton = false
      //     this.rowId = ""
      //     return
      //   } else {
      //     let isShow = this.$refs['popover_id-' + row.id].showPopper
      //     if (!isShow) {
      //       this.showOperateButton = false
      //       this.rowId = ""
      //     }
      //   }
      // }
    },
    handleCommand(command){
      console.log(command)

        if(command == 'a'){
            this.gotoTest();
        }
        if(command == 'selectedall'){
          this.selectCurrentPage();
        }
        if(command == 'noselectedall'){
          this.deselectCurrentPage();
        }
    },

  },
  watch: {
    id(newVal) {
      if (newVal) {
        this.getData()
      }
    },
    selectedNode(newVal) {
      if (newVal) {
        // this.formAdd.selectedNode = newVal;

        this.formAdd.selectedNode = newVal.split('_')[1];
        if(typeof this.selectedNode == 'string'&&this.selectedNode.includes('_') ){
          this.fatherNodeID = this.selectedNode.split('_')[1];
        }else{
          this.fatherNodeID = newVal;
        }
      }
    }
  },
}
</script>
