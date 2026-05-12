<template>
    <div>
        <div class="main-title"  >
            <router-link :underline="false" class="classA" :to="{ path: '/burpsuite' }"  >Burpsuite
            </router-link> 
            <label class="currentpagetitle" >   
                <span>详情</span>
            </label> 
        </div>
        <div class="tasklist context_box_bg">
      <div class="search-box">
        <div class="operationbutton"> 
       
        </div>
      
      </div>
      <el-table ref="myTable" 
      
      :data="tableData" v-loading="Loading" 
        @cell-mouse-enter="mouseenter" 
        @cell-mouse-leave="mouseleave" 
        style="width: 100%;padding:0 50px;"
      
        >

     


        <el-table-column prop="action" label="Action"  :show-overflow-tooltip="true"></el-table-column>
        <el-table-column prop="issueType" label="issue type"  :show-overflow-tooltip="true"></el-table-column>
        <el-table-column prop="host" label="host"  :show-overflow-tooltip="true"></el-table-column>
        <el-table-column prop="path" label="path"  :show-overflow-tooltip="true"></el-table-column>
        <el-table-column prop="insertionPoint" label="insertion point"  :show-overflow-tooltip="true"></el-table-column>
        <el-table-column prop="riskEnum" label="severity"  :show-overflow-tooltip="true"></el-table-column>
        <el-table-column prop="confidence" label="confidence"  :show-overflow-tooltip="true"></el-table-column>
        <el-table-column prop="time" label="time">
            <template slot-scope="scope">
    
              <div v-if="showOperateButton && rowId == scope.row.id  ">
            
                <!-- 已结束 4 -->
                <div  >
                  <el-link class="link_primary" :underline="false" @click="clickDetail(scope.row)">详情</el-link>
                </div>
          
              </div>
              <div v-else> 
                <span>{{ scope.row.time }}</span>
              </div>

            </template>
          </el-table-column>
   
       
     

      </el-table>
      <el-pagination :page-size="pageSize" background layout="total, prev, pager, next, sizes, jumper"
        :total="totalpage" :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
      </el-pagination>
    </div>
       <el-dialog :title=" detailAlert.issueType"  :visible.sync="dialogaddFormVisible"  width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click=" cancelform">关闭</el-button>
            </div>
            <div style="padding:24px">
                <div style="margin-top:15px;">
                   <span >漏洞风险：</span> <span style="margin-left:15px">{{detailAlert.riskEnum}}</span>
                </div>
                <div style="margin-top:15px;">
                   <span >主机：</span> <span style="margin-left:15px">{{detailAlert.host}}</span>
                </div>
                <div style="margin-top:15px;">
                   <span >路径：</span> <span style="margin-left:15px">{{detailAlert.path}}</span>
                </div>
                <div style="margin-top:15px;">
                   <span >漏洞描述：</span> <span style="margin-left:15px">{{detailAlert.desc}}</span>
                </div>
                <div style="margin-top:15px;">
                   <span >修复建议：</span> <span style="margin-left:15px" v-html="detailAlert.fix"></span>
                </div>
                <div style="margin-top:15px;">
                   <span >测试报文</span> 
                </div>
                <div style="margin-top:15px;display:flex" v-for="(itt,idd) in detailAlert.requestResponse" :key="idd">
                 
                  <div style="border:1px solid #f2f2f2;width:50%;height:400px;overflow-y:auto;overflow-x:hidden">
                     <div style="margin-bottom:20px">请求报文</div>
                    <pre>{{itt[0]}}</pre>
                  </div>
                  <div style="border:1px solid #f2f2f2;width:50%;height:400px;overflow-y:auto;overflow-x:hidden">
                    <div style="margin-bottom:20px">响应报文</div>
                    <pre>{{itt[1]}}</pre>
                  </div>
                  






                </div>
            </div>
        </el-dialog>
    <CreateReport :type="1" ref="CreateReport"  :dialogVisible = 'dialogVisible' @click="saveCreate()" @clearCreate="clearCreate()"></CreateReport>
  </div>
</template>
<style lang="less" scoped>
  .demo-table-expand {
    font-size: 0;
  }
  .demo-table-expand label {
    width: 90px;
    color: #99a9bf;
  }
  .demo-table-expand .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
    width: 90%;
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
      color: #4C7AE3;

      i {
        color: #4C7AE3;
      }
    }
  }

  .cell {
    line-height: 15px;

    >span {
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
    import { task } from '@/api/Tripartitetools.js'
    import CreateReport from "@/components/CreateReport.vue";
    export default {  
        name:'detailXray',
        components:{
          CreateReport
        },
        data(){
            return{
              detailAlert:{},
                dialogaddFormVisible:false,
                jsonData: {
                total: 25,
                limit: 10,
                skip: 0,
                links: {
                  previous: undefined,
                  next: function () {},
                }},


              timeIDFOrSHouQuan:null,
                isShowShouQuan:true,
                alldelvisible:false,
                Loading:false,
                formData:{
                    page_num: 1,
                    task_status: "",
                    starttime: "",
                    stoptime: "",
                    search_field: "",
                    time: "",
                    risk_level: 0
                },
                tableData:[],
                multipleSelection:[],
                totalpage:0,
                currentpage:1,
                risklevellist:[],
                pageSize:10,
                showOperateButton:false,
                rowId:'',
                timer:null, 
                dialogVisible:false,
            }
        },
        created(){
            this.$store.state.activefirstMenu = '/burpsuite';
       
        },
        mounted(){
            this.getData();
           
        },
        beforeDestroy () { 
        },
        methods:{
          clickDetail(row){
            this.dialogaddFormVisible = true
            this.detailAlert = row
            console.log(this.detailAlert,'this.detailAlert');
          },
           cancelform(){
            this.dialogaddFormVisible = false
            this.detailAlert = {}
            console.log(this.detailAlert,'this.detailAlert---关闭');
           },
    
            async getData(){  
             try { 
                const res = await task.taskListDetail2({
                    page: this.formData.page_num,
                    size: this.pageSize,
                    burpsuiteId: this.$route.query.id,
                    search :this.$route.query.name      
                });
                if(res.code == 200){
                    this.tableData = res.data.list;
                    this.totalpage = res.data.total;
                }else{
                    this.$message({
                        message: res.msg,
                        type: "error"
                    });
                }
             } catch (error) {
              console.log(error)
             }
            },
            btnNewTask(){
                this.$router.push({
                    path: `/createXray`,
                    query: { 
                        flag: 1, 
                    }
                });
            },
            handlesearch(){
                 //搜索
                this.formData.page_num = 1;
                this.getData();
                this.currentpage = 1;
            },
            handleReset(){
                this.formData.page_num = 1;
                this.formData.search_field = "";
                this.formData.time = "";
                this.formData.task_status = "";
                this.formData.risk_level = 0;
                this.pageSize = 10;
                this.currentpage = 1;
                this.getData();
            },
            btnTaskinfo(row){ //详情
              this.$router.push({
                    path: `/detailXray`,
                    query: { 
                        id: row.id, 
                        name:row.taskName,
                    }
                });
                // localStorage.setItem("taskTab", "tabs1");


            },
            async btnResumeTask(taskid){  //开始
              const dt = await task.taskchangestate({
                taskId: taskid,
                operate:'resume'
              });
              if (dt.code == 200) {
                this.getData();
              } else {
                this.$message({
                  message: dt.msg,
                  type: "error"
                });
              }
            },
            async btnTaskstop(taskid){ //结束
              const dt = await task.taskchangestate({
                taskId: taskid,
                operate:'stop'
              });
              if (dt.code == 200) {
                this.$message({
                  message: "结束任务成功",
                  type: "success"
                }); 
                this.getData();
              } else {
                this.$message({
                  message: dt.msg,
                  type: "error"
                });
              }
            },
            async btnPauseTask(taskid){ //暂停
              const dt = await task.taskchangestate({
                taskId: taskid,
                operate:'pause'
              });
              if (dt.code == 200) {
                this.getData();
              } else {
                this.$message({
                  message: dt.msg,
                  type: "error"
                });
              }
            },
            async btnMultiDeleteTask(){
              if (this.multipleSelection.length == 0) return; 
              let _ids = this.multipleSelection.map(item => item.id); 
              
              const res = await	task.taskDel({
                taskIds: _ids.join(",")
              }); 
              if (res.code == 200) {
                this.$message({
                  message: "删除任务成功",
                  type: "success"
                }); 
                this.currentpage = 1;
                this.formData.page_num = 1;
                this.alldelvisible = false;
                this.getData();
              } else {
                this.$message({
                  message: res.msg,
                  type: "error"
                });
              }
       
            },
            async btnDel(scope){
              const res = await task.taskDel({ taskIds: scope.row.id + '' })
       
              if (res.code == 200) {
                  this.$message({
                    message: "删除任务成功",
                    type: "success"
                  });
                  scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
              
                  this.getData();
              } else {
                  this.$message({
                    message: res.msg,
                    type: "error"
                  });
              } 
            },
            // 取消删除
            btnCancelDel(scope){
              scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
              if (scope.row.executeType == 3 || scope.row.status !=1) { 
                scope._self.$refs[`popover-${scope.row.id}`].doClose();
              }
            },
            async btnCopyTask(row){ //任务复制
              const dt = await task.taskCopy({
                taskId: row.id, 
              });
              if (dt.code == 200) {
                this.$message({
                    message: "复制任务成功",
                    type: "success"
                  });
                this.getData();
                this.$store.commit("updateCopyArr", dt.data) 
                this.$router.push({
                    path: `/createtask`,
                    query: { 
                        flag: 2, 
                    }
                });
              } else {
                this.$message({
                  message: dt.msg,
                  type: "error"
                });
              }

            },
            handleSizeChange(t){
                this.formData.page_num = 1;
                this.pageSize = t;
                this.getData();
            },
            currentchange(t){
                this.formData.page_num = t;
                this.getData();
                this.currentpage = t;
            },
      
            mouseenter(row, colum, cell, event) {
                this.showOperateButton = true;
                this.rowId = row.id   //赋值行id，便于页面判断  
            },
            mouseleave(row, colum, cell, event) { 
                if (row.executeType == 3 || row.status != 1){  
                    if (!this.$refs['popover-' + row.id]) {
                    this.showOperateButton = false;
                    this.rowId = "";
                    return;
                    }else{
                    let isShow = this.$refs['popover-' + row.id].showPopper;
                    if (!isShow) {
                        this.showOperateButton = false;
                        this.rowId = "";
                    }
                    }  
                }else{ 
                    if (!this.$refs['popover_id-' + row.id]){
                    this.showOperateButton = false;
                    this.rowId = "";
                    return;
                    }else{
                    let isShow = this.$refs['popover_id-' + row.id].showPopper;
                    if (!isShow) {
                        this.showOperateButton = false;
                        this.rowId = "";
                    }
                    } 
                } 
            },
            btnReport(row){ //点击报告 
              this.dialogVisible = true;
              this.$refs.CreateReport.getinit(row.id,row.taskName);
            },
            async saveCreate(params){ //生成成功，  
              this.dialogVisible = false; 
            },
            clearCreate(){
              this.dialogVisible = false;
            },
        }

    }
</script>