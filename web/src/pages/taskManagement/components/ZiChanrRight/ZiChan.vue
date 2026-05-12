<template>
  <div>
    <div v-show="false" class="main-title">
      <label class="taskSceneBtn">Xray</label>
    </div>
    <div class="tasklist context_box_bg">
      <div class="search-box">
        <div class="operationbutton">
     
         
          <el-upload
            ref="excelUploadRef"
            style="display: inline-block; width: 0px"
            action="#"
            :show-file-list="false"
            :before-upload="beforeUpload"
            :http-request="upload"
            :on-preview="handlePreview"
            :on-remove="handleRemove"
            :before-remove="beforeRemove"
            :limit="1"
            :on-exceed="handleExceed"
            :file-list="fileList"
            accept=".json"
          >
            <el-button
               v-show="false"
              style="margin-right: 100px"
              type="primary"
              class="btn-blue"
              >上传任务</el-button
            >
          </el-upload>
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
              placeholder="搜索任务目标与系统"
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
            
           <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                   
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
          prop="assetName"
          label="资产"
          :show-overflow-tooltip="true"
        >
          <template slot-scope="scope">
            <el-link @click="btnTaskinfo(scope.row)">{{
              scope.row.assetName
            }}</el-link>
          </template>
        </el-table-column>

        <el-table-column prop="system" label="系统"> </el-table-column>
        <el-table-column prop="openPort" label="开放端口"> </el-table-column>
        <el-table-column prop="assetRiskName" label="资产风险">
            <template slot-scope="scope">
                <span :style="getRiskLevelColor(scope.row.assetRiskName)">
                  {{ scope.row.assetRiskName }}
                </span>
              </template>
         </el-table-column>
        <el-table-column prop="vulInfo" label="漏洞">
               <template slot-scope="scope">
                    <div style="display:flex;justify-content: space-between;">
                        <div>
                          <span style="color:red;margin-right:3px">{{scope.row.vulInfo.deadlyVul}}</span>
                          <span style="color:#f9ca84;margin-right:3px">{{scope.row.vulInfo.highRiskVul}}</span>
                          <span style="color:#1223ff;margin-right:3px">{{scope.row.vulInfo.mediumRiskVul}}</span>
                          <span style="color:#8ac23a">{{scope.row.vulInfo.lowRiskVul}}</span>
                        </div>
                          <span style="color:#00">{{scope.row.vulInfo.totalRiskVul}}</span>
                    </div>
                </template>
         </el-table-column>
        <el-table-column prop="assetRiskName" label="测试时间">
          <template slot-scope="scope">
            <div v-if="showOperateButton && rowId == scope.row.id">
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
              <span><i></i>{{ scope.row.testTime }}</span
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

    </div>
  </div>
</template>
<style lang="less" scoped>
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
import { traffic } from '@/api/assetManagement'; 


// import { task } from '@/api/task.js'
export default {
  name: 'x-ray',
  components: {
  },
  data () {
    return {
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
    // this.$store.state.activefirstMenu = '/x-ray'

  },
  mounted () {
    // this.getData()
  },

  methods: {
      getRiskLevelColor(levelRiskName) {
        const riskColors = {
          '高危': 'red',
          '中危': 'orange',
          '低危': 'yellow',
          '安全': 'green'
        };
        return { color: riskColors[levelRiskName] || 'defaultColor' }; // 如果没有匹配的风险等级，返回默认颜色
      },
    cancelform () {
      // alert(1)
    },
    beforeUpload (file) {
      const fileType = file.type
      if (fileType !== 'application/json') {
        this.$message.error('只能上传JSON文件')
        return false  // 阻止文件上传
      } else {
        // console.log(file, 'file---------------------------');
      }

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
      console.log(obj, 'obj');
      let that_ = this
      try {

        let formData = new FormData() //使用FormData来传输文件，就是form的升级版
        formData.append('file', obj.file)
        formData.append('taskName', obj.file.name)
        this.$ajax
          .post("/smart/tripartite/xrayupload", formData, {
            headers: { 'Content-Type': 'multipart/form-data' }  //设置请求头
          })
          .then((dt) => {
            if(dt.data.code == 200){
              this.$message({
                message: "上传成功",
                type: "success"
              })
              this.fileList = []
              this.getData()

            }else{
              this.$message({
                message: dt.data.msg,
                type: "error"
              })
              this.fileList = []

            }
          })


      } catch (error) {
        console.log(eror)
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
      return this.$confirm(`确定移除 ${file.name}？`)
    },

    //  新增取消
    handleSaveCancel () {
      this.createDictionaries = false
    },

    async getData () {
      try {
        const res = await traffic.trafficvulndel({
          page: this.formData.page_num,
          size: this.pageSize,
          search: this.formData.search_field,
          groupID:this.$store.state.groupID[0]||-1, 
        })
        if (res.code == 200) {
          this.tableData = res.data.assetsInfo
          console.log(this.tableData, 'this.tableData')
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
            this.formData.page_num = 1;
            this.formData.search_field = ''; 
            this.currentPage = 1;
            this.getData();
        },
    // handleReset () {
    //   this.formData.page_num = 1
    //   this.formData.search_field = ""
    //   this.formData.time = ""
    //   this.formData.task_status = ""
    //   this.formData.risk_level = 0
    //   this.pageSize = 10
    //   this.currentpage = 1
    //   this.getData()
    // },
    btnTaskinfo (row) { //详情

      // 跳转到概览
          this.$emit('gotoFirstTab',[row.assetID,2])


    },



    async btnMultiDeleteTask () {
      if (this.multipleSelection.length == 0) return
      let _ids = this.multipleSelection.map(item => item.assetID)

      const res = await traffic.trafficlog({
        ids: _ids.join(",")
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
      console.log(scope,'scopescopescope');
      const res = await traffic.trafficlog({ ids: scope.assetID + '' })

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
    async btnCopyTask (row) { //任务复 制
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
      this.multipleSelection = val
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

  }

}
</script>