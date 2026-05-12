<template>
  <div>
    <div class="main-title">
      <label class="taskSceneBtn">WIFI检测</label>
    </div>
    <div class="tasklist context_box_bg">
      <div class="search-box">
        <div class="operationbutton">
          <el-button type="primary" size="small" @click="btnNewTask"
            >新建</el-button
          >
          <!-- <el-button
            type="primary"
            size="small"
            @click="createDictionaries = true"
            >上传</el-button
          > -->
       
          <el-popover
            popper-class="delButton_popper"
            placement="bottom-start"
            width="170"
            style="padding-left: 8px"
            trigger="click"
            :visible-arrow="false"
            v-model="alldelvisible"
          >
            <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
            <div style="text-align: right; margin: 0" class="">
              <el-button
                size="mini"
                class="delCancel"
                @click="alldelvisible = false"
                >取消</el-button
              >
              <el-button size="mini" type="primary" @click="btnMultiDeleteTask"
                >确定</el-button
              >
            </div>
            <el-button
              type="warning"
              size="small"
              slot="reference"
              :disabled="!multipleSelection.length"
              >删除</el-button
            >
          </el-popover>
        </div>
        <div class="serach-condition">
          <div v-show="false">
            <el-date-picker
              v-model="formData.time"
              type="daterange"
              format="yyyy-MM-dd"
              value-format="yyyy-MM-dd HH:mm:ss"
              :default-time="['00:00:00', '23:59:59']"
              range-separator="—"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              size="small"
              clearable
            >
            </el-date-picker>
          </div>
          <div class="search-text">
            <el-input
              style="width: 240px"
              placeholder="搜索WiFi名称"
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
      >
        <el-table-column type="selection" width="55"> </el-table-column>

        <el-table-column
          prop="taskName"
          label="任务名称"
          width="450px"
          :show-overflow-tooltip="true"
        >
          <template slot-scope="scope">
            <el-link @click="btnTaskinfo(scope.row)">{{
              scope.row.taskName
            }}</el-link>
          </template>
        </el-table-column>

        <el-table-column prop="ssid" label="WIFI名称"> </el-table-column>
       
        <el-table-column prop="carrier" label="无线协议类型"> </el-table-column>
        <el-table-column prop="encrypt" label="加密方式"> </el-table-column>
        <el-table-column prop="channel" label="信道"> </el-table-column>
        <el-table-column prop="createTime" label="开始时间"> </el-table-column>
        <el-table-column prop="status" label="状态">
          <template slot-scope="scope">
            <div v-if="showOperateButton && rowId == scope.row.taskId">
              <!-- 已结束 4 -->
              <div>
                <el-link
                  class="link_primary"
                  :underline="false"
                  @click="btnTaskinfo(scope.row)"
                  >详情</el-link
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
              <span><i></i>{{ scope.row.status }}</span
              >
            </div>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        :page-size="pageSize"
        background
        layout="total, prev, pager, next, sizes, jumper"
        :total="totalpage"
        :current-page="currentpage"
        @current-change="currentchange"
        @size-change="handleSizeChange"
      >
      </el-pagination>
      <!-- <el-dialog
        class="zidiandia"
        title="上传"
        :visible.sync="createDictionaries"
        v-if="createDictionaries"
        :before-close="cancelform"
        width="1184px"
        :close-on-click-modal="false"
        :show-close="false"
      >
        <div class="dialog_b_btn">
          <el-button size="small" @click="handleCreateSave">保存</el-button>
          <el-button size="small" @click="handleSaveCancel">取消</el-button>
        </div>
        
      </el-dialog> -->
         <el-dialog title="新建WIFI检测任务"  :visible.sync="dialogaddFormVisible"  width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="handleCreateSave">保存</el-button>
                <el-button size="small" @click=" cancelform">关闭</el-button>
            </div>
            <div style="padding:24px;  ">
              <div style="display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px">任务名称</div>
                  <el-input clearable style="width:300px;margin-left:50px" v-model="dialogTaskName" placeholder="请输入任务名称"></el-input>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px">检测目标</div>
              </div>
              <!-- tabel -->
             <div style="margin-top:10px">
               <el-table  highlight-current-row class="addDialogtb"
                  ref="myTableAdd"
                  :data="WIFIlist"
                  style="width: 100%" @row-click="handleRowClick" @select="handleonecheck">
                  <el-table-column type="selection" width="55"> </el-table-column>
                  <el-table-column prop="ssid" label="SSID"> </el-table-column>
                  <el-table-column prop="lastSignalRssi" label="信号强度"> </el-table-column>
                  <el-table-column prop="ssidCryptsetEnum" label="加密类型"> </el-table-column>
                  <el-table-column prop="sourceMac" label="MAC"> </el-table-column>
                  <el-table-column prop="manuf" label="厂商"> </el-table-column>
                  <el-table-column prop="carrier" label="无线协议"> </el-table-column>
                  <!-- <el-table-column prop="status" label="">
                    <template slot-scope="scope">
                    
                        <div>
                          <el-link
                            class="link_primaryBywang"
                            :underline="false"
                            @click="xinJianTask(scope.row)"
                            >选择</el-link
                          >
                         
                        </div>
                    
                    </template>
                </el-table-column> -->
                </el-table>
             </div>
            </div>
        </el-dialog>



         <el-dialog :title="rowData.taskName"  :visible.sync="dialogaddFormVisible2"  width="800px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <!-- <el-button size="small" @click="handleCreateSave">保存</el-button> -->
                <el-button size="small" @click=" cancelform2">关闭</el-button>
            </div>
            <div style="padding:24px;  ">
              <div style="display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">任务名称</div>
                <span style="margin-left:40px">{{rowData.taskName}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">WIFI名称</div>
                <span style="margin-left:40px">{{rowData.ssid}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">无线协议类型</div>
                <span style="margin-left:40px">{{rowData.carrier}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">加密方式</div>
                <span style="margin-left:40px">{{rowData.encrypt}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">信道</div>
                <span style="margin-left:40px">{{rowData.channel}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">开始时间</div>
                <span style="margin-left:40px">{{rowData.createTime}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">状态</div>
                <span style="margin-left:40px">{{rowData.status}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">ssid</div>
                <span style="margin-left:40px">{{rowData.ssid}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">taskId</div>
                <span style="margin-left:40px">{{rowData.taskId}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">passwd</div>
                <span style="margin-left:40px">{{rowData.passwd}}</span>
              </div>
              <div style="margin-top:50px;display:flex;justify-content: flex-start; align-items: center;">
                <div style="border-left:3px solid #4c7ae3;padding-left:10px;width:120px">日志</div>
                
              </div>
                <div style="margin-left:100px" v-for="(item,index) in rowData.logList" :key="index">
                  <div>{{item.time}}</div>
                  <div style="margin-bottom:10px">{{item.log}}</div>
                  
                </div>
       
            </div>
        </el-dialog>
    </div>
  </div>
</template>
<style lang="less" scoped>
.addDialogtb /deep/ th .el-checkbox__inner{
  display: none;
}
.link_primaryBywang{
  color: #60607a;
}
.tasklist {
  padding: 24px;
  background: #fff;
  min-height: calc(100% - 39px);
  box-sizing: border-box;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
/deep/ thead {
  .cursorPointer {
    cursor: pointer;
    position: absolute;
    // top: 6px;
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
import { task } from '@/api/detectionTask.js'

// import { task } from '@/api/task.js'
export default {
  name: 'detectionTask',
  components: {
  },
  data () {
    return {
      dialogaddFormVisible2:false,
      rowData: {},//行数据
      targetMAC: '',
      WIFIlist:[],
      dialogTaskName:'',
      dialogaddFormVisible:false,
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
        risk_level: 0
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
    }
  },
  created () {
    this.getWIFIData()
    this.$store.state.activefirstMenu = '/detectionTask'

  },
  mounted () {
    // this.getTaskEnum();
    this.getData()
    let _this = this
    // this.timer = setInterval(function(){
    //   _this.getData();
    // },5000)

  },

  methods: {
    // 新建WiFi检测任务
    xinJianTask(row){
      this.targetMAC = row.sourceMac
      this.$message({
        message: "选择成功",
        type: "success"
      })
    },
    cancelform () {
      // alert(1)
      this.dialogaddFormVisible = false
    },
    cancelform2 () {
      // alert(1)
      this.dialogaddFormVisible2 = false
    },
   async handleCreateSave(){
      if(this.dialogTaskName == ''){
        this.$message({
          message: "请输入任务名称",
          type: "error"
        })
        return
      }
      if(this.targetMAC == ''){
        this.$message({
          message: "请选择检测目标",
          type: "error"
        })
        return
      }
       const res = await task.createWIFITask({
        taskName:this.dialogTaskName,
        sourceMac:this.targetMAC
       })
        if(res.code == 200){
          this.$message({
            message: "新建任务成功",
            type: "success"
          })
          this.dialogaddFormVisible = false
          this.getData()
          this.dialogTaskName = ''
          this.targetMAC = ''
          }

    },


    //  新增取消
    handleSaveCancel () {
      this.createDictionaries = false
    },

    async getWIFIData () {
      try {
        const res = await task.wifiList()
        if (res.code == 200) {
         console.log(res.data.list,'res-----getWIFIData');
         this.WIFIlist = res.data.list
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
    async getData () {
      try {
        const res = await task.taskList({
          page: this.formData.page_num,
          size: this.pageSize,
          search: this.formData.search_field,
          // riskLevel: this.formData.risk_level === 0 ? '' : this.formData.risk_level,
          // startTime:this.formData.time == null? '': this.formData?.time[0],
          // endTime: this.formData.time == null? '': this.formData?.time[1],
        })
        if (res.code == 200) {
          this.tableData = res.data.list
          this.totalpage = res.data.total
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
    btnNewTask () {
      this.dialogaddFormVisible = true
      // this.$router.push({
      //   path: `/createXray`,
      //   query: {
      //     flag: 1,
      //   }
      // })
    },
    handlesearch () {
      //搜索
      this.formData.page_num = 1
      this.getData()
      this.currentpage = 1
    },
    handleReset () {
      this.formData.page_num = 1
      this.formData.search_field = ""
      this.formData.time = ""
      this.formData.task_status = ""
      this.formData.risk_level = 0
      this.pageSize = 10
      this.currentpage = 1
      this.getData()
    },
    btnTaskinfo (row) { //详情
      this.rowData = row
      this.dialogaddFormVisible2 = true
      console.log(row, 'row');


    },



    async btnMultiDeleteTask () {
      if (this.multipleSelection.length == 0) return
      let _ids = this.multipleSelection.map(item => item.taskId)

      const res = await task.taskDel({
        taskIds: _ids.join(",")
      })
      if (res.code == 200) {
        this.$message({
          message: "删除任务成功",
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
    async btnDel (scope) {
      const res = await task.taskDel({ taskIds: scope.taskId + '' })

      if (res.code == 200) {
        this.$message({
          message: "删除任务成功",
          type: "success"
        })

        this.getData()
      } else {
        this.$message({
          message: res.msg,
          type: "error"
        })
      }
    },
    // 取消删除
    btnCancelDel (scope) {
      scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
      if (scope.row.executeType == 3 || scope.row.status != 1) {
        scope._self.$refs[`popover-${scope.row.id}`].doClose()
      }
    },
    async btnCopyTask (row) { //任务复制
      const dt = await task.taskCopy({
        taskId: row.id,
      })
      if (dt.code == 200) {
        this.$message({
          message: "复制任务成功",
          type: "success"
        })
        this.getData()
        this.$store.commit("updateCopyArr", dt.data)
        this.$router.push({
          path: `/createtask`,
          query: {
            flag: 2,
          }
        })
      } else {
        this.$message({
          message: dt.msg,
          type: "error"
        })
      }

    },
    handleSizeChange (t) {
      this.formData.page_num = 1
      this.pageSize = t
      this.getData()
    },
    currentchange (t) {
      this.formData.page_num = t
      this.getData()
      this.currentpage = t
    },
    handleSelectionChange (val) {
      console.log(val, 'val-------')
      this.multipleSelection = val
    },
    mouseenter (row, colum, cell, event) {
      this.showOperateButton = true
      this.rowId = row.taskId   //赋值行id，便于页面判断  
    },
    mouseleave (row, colum, cell, event) {
        if (!this.$refs['popover-' + row.taskId]) {
          this.showOperateButton = false
          this.rowId = ""
          return
        } else {
          let isShow = this.$refs['popover-' + row.taskId].showPopper
          if (!isShow) {
            this.showOperateButton = false
            this.rowId = ""
          }
        }
    },
    handleonecheck(selection, row){  
        this.$refs.myTableAdd.clearSelection();
        if(selection.length === 0 ){
            return;
        }
        if(row){
            this.multipleSelection = row;
            this.$refs.myTableAdd.toggleRowSelection(row,true); 
        }

        this.targetMAC = row.sourceMac
        this.$message({
          message: "选择成功",
          type: "success"
        })

        // this.dialog_title = row.vul_name;
        // this.dialog_desc = row.vul_description;
        // this.testform.vul_id = row.id;

        // $('.testloop').removeClass('closeAnimate').addClass('openAnimate'); 
    },
    handleRowClick(row){
        this.$refs.myTableAdd.clearSelection();
        this.$refs.myTableAdd.toggleRowSelection(row);

        this.targetMAC = row.sourceMac
        this.$message({
          message: "选择成功",
          type: "success"
        })
    },

  }

}
</script>