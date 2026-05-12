<template>
    <div>
        <div class="main-title" v-if="!group_id" >
            <label class="taskSceneBtn" >漏扫任务</label>
        </div>
        <div class="tasklist context_box_bg">
      <div class="search-box">
        <div class="operationbutton"> 
          <el-button type="primary" size="small" @click="btnNewTask" :disabled='!isShowShouQuan'>新建</el-button>
          <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" style="padding-left:8px"
            trigger="click" :visible-arrow="false" v-model="alldelvisible">
            <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
            <div style="text-align: right; margin: 0" class="">
              <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消</el-button>
              <el-button size="mini" type="primary" @click="btnMultiDeleteTask">确定</el-button>
            </div>
            <el-button type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除</el-button>
          </el-popover>
        </div>
        <div class="serach-condition">
          <div>
            <el-date-picker v-model="formData.time" type="daterange" format="yyyy-MM-dd"
              value-format="yyyy-MM-dd HH:mm:ss" :default-time="['00:00:00', '23:59:59']" range-separator="—"
              start-placeholder="开始日期" end-placeholder="结束日期" size="small" clearable>
            </el-date-picker>
          </div>
          <div class="search-text">
            <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch" v-model="formData.search_field" class="input-with-select" size="small"
              clearable>
            </el-input>
            <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
          </div>
          <div>
            <el-button type="primary" size="small" @click="handleReset">重置</el-button>
          </div>
        </div>
      </div>
      <el-table ref="myTable" 
      :data="tableData" v-loading="Loading" 
      @selection-change="handleSelectionChange"
        @cell-mouse-enter="mouseenter" 
        @cell-mouse-leave="mouseleave" 
        style="width: 100%"
      
        >
        <el-table-column type="selection" width="55"> </el-table-column>

        <el-table-column prop="taskName" label="任务名称" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <el-link @click="btnTaskinfo(scope.row)">{{
            scope.row.name
            }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="executeTypeName" label="任务类型">
          <template slot-scope="scope">
            {{ scope.row.typeName}}
          </template>
        </el-table-column>
        <el-table-column   label="任务风险" show-overflow-tooltip>
          <template slot="header">
            <!-- <div style="position: relative;height:10px;"> -->
              <span class="cursorPointer" @click="clickButton('任务风险')"
                :class="(formData.risk !== '' && formData.risk !== 0) ? 'active' : ''">任务风险<i
                  class="iconfont iconshaixuan"></i>
              </span>
              
            <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.risk" clearable
              size="small" ref="loglistRef" @change="handlesearch">
              <el-option v-for="(item, index) in risklevellist" :key="index" :label="item.label" :value="item.value">
              </el-option>
            </el-select>
          </template>
          <template slot-scope="scope">
            <span :class="[
              { 'riskstyle risk_hight': scope.row.risk == '1' },
              { 'riskstyle risk_middle': scope.row.risk == '2' },
              { 'riskstyle risk_low': scope.row.risk == '3' },
              { 'riskstyle risk_nofind': scope.row.risk == '4' }
            ]"><i></i>{{ scope.row.riskName }} </span>
            <!-- <i v-if="scope.row.task_risk_icon.verify_success"
              class="iconfont iconloudongyanzhengchenggong iconstyle iconstylecor3"></i>
            <i v-if="scope.row.task_risk_icon.used_success"
              class="iconfont iconloudongliyongchenggong iconstyle iconstylecor1"></i>
            <i v-if="scope.row.task_risk_icon.remote_control"
              class="iconfont iconloudongyuanchengkongzhi iconstyle iconstylecor2"></i> -->
          </template>
        </el-table-column>
        <el-table-column prop="" label="目标风险">
          <template slot="header">
            目标风险
            <el-tooltip class="item" effect="dark" placement="right">
              <div slot="content">
                从左往右依次为高危目标、中危目标、低危目标、安全目标
              </div>
              <i class="iconfont icontishi" style="position: absolute;top:0;left:66px"></i>
            </el-tooltip>
          </template>
          <template slot-scope="scope">
            <span class="tag_status tag_danger bug_status">{{scope.row.high}}</span>
            <span class="tag_status tag_warning bug_status">{{scope.row.middle}}</span>
            <span class="tag_status tag_primary bug_status">{{scope.row.low}}</span>
            <span class="tag_status tag_success bug_status">{{scope.row.safe}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="更新时间"> </el-table-column>
        <!-- <el-table-column prop="" label="任务进度"> -->
          <!-- <template slot-scope="scope"> -->
              <!-- <el-progress :text-inside="true" :stroke-width="12" :percentage="scope.row.progress"></el-progress> -->
          <!-- </template> -->
        <!-- </el-table-column> -->
        <el-table-column prop="statusName" label="状态">
          <template slot-scope="scope">
            <!-- 待触发  1 ；待执行  2；运行中   3；已完成   4 ；暂停中   5  -->
            <!--    
              ——待触发1：结束、删除、复制
              ——待执行2：结束、删除、复制、优先执行
              ——运行中3：暂停、结束、删除、复制
              ——已结束4：报告、删除、复制、验证测试
              ——暂停中5：开始、结束、删除、复制

              周期任务，每个操作后面加 终止周期
						-->
            <div v-if="showOperateButton && rowId == scope.row.id  ">
              <!-- 待触发 1 -->
              <div v-if="scope.row.status == 1">
               
              </div>
              <!-- 待执行 2 -->
              <div v-else-if="scope.row.status == 2">
                <el-link class="link_primary" :underline="false" @click="btnTaskstop(scope.row.id)">结束</el-link>
                
                <el-popover placement="bottom" width="170" :visible-arrow="false"
                    :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                    <div style="text-align: right; margin: 0">
                        <el-button size="mini" class="delCancel"
                            @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消
                        </el-button>
                        <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                    </div>
                    <el-link :underline="false" class="link_danger linkafter" style="padding:0"
                        slot="reference">删除</el-link>
                </el-popover>
                
              </div>
              <!-- 运行中 3 -->
              <div v-else-if="scope.row.status == 3">
                <el-link class="link_primary" :underline="false" @click="btnTaskstop(scope.row.id)">结束</el-link>
                <el-popover placement="bottom" width="80" trigger="click" popper-class="learnMore"
                  :ref="`popover_id-${scope.row.id}`" :visible-arrow="false" style="padding:0">
                  <ul class="operationbox">
                    <li>
                      <el-popover placement="bottom" width="170" :visible-arrow="false"
                        :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                        <p class="delText">
                          <i class="el-icon-warning"></i>确定删除吗？
                        </p>
                        <div style="text-align: right; margin: 0">
                          <el-button size="mini" class="delCancel" @click="btnCancelDel(scope)">取消</el-button>
                          <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                        </div>
                        <div slot="reference">删除</div>
                      </el-popover>
                    </li>
                  </ul>
                  <el-link :underline="false" class="link_info" slot="reference">更多</el-link>
                </el-popover>
              </div>
              <!-- 已结束 4 -->
              <div v-else-if="scope.row.status == 4">
                <el-link class="link_primary" :underline="false" @click="btnReport(scope.row)">报告</el-link>
                <el-link class="link_primary" :underline="false" @click="btnTaskstop(scope.row.id)">结束</el-link>
               
                <el-popover placement="bottom" width="170" :visible-arrow="false"
                    :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                    <div style="text-align: right; margin: 0">
                        <el-button size="mini" class="delCancel"
                            @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消
                        </el-button>
                        <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                    </div>
                    <el-link :underline="false" class="link_danger linkafter" style="padding:0"
                        slot="reference">删除</el-link>
                </el-popover>

                
              </div>
              <!-- 暂停中 5 -->
              <div v-else-if="scope.row.status == 5">
                <el-link class="link_primary" :underline="false" @click="btnTaskstop(scope.row.id)">结束</el-link>
                <el-popover placement="bottom" width="80" trigger="click" popper-class="learnMore"
                  :ref="`popover_id-${scope.row.id}`" :visible-arrow="false" style="padding:0">
                  <ul class="operationbox">
                    <li>
                      <el-popover placement="bottom" width="170" :visible-arrow="false"
                        :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                        <p class="delText">
                          <i class="el-icon-warning"></i>确定删除吗？
                        </p>
                        <div style="text-align: right; margin: 0">
                          <el-button size="mini" class="delCancel" @click="btnCancelDel(scope)">取消</el-button>
                          <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                        </div>
                        <div slot="reference">删除</div>
                      </el-popover>
                    </li>

                  </ul>
                  <el-link :underline="false" class="link_info" slot="reference">更多</el-link>
                </el-popover>
              </div>
            </div>
            <div v-else> 
              <span :class="[ 
                  { 'tag_status tag_danger1': scope.row.status == 1 } ,
                  { 'tag_status tag_warning': scope.row.status ==2 },
                  { 'tag_status tag_primary': scope.row.status == 3 },
                  { 'tag_status tag_success': scope.row.status ==4 },
              { 'tag_status tag_danger': scope.row.status == 5 }]"><i></i>{{ scope.row.statusName }}</span>
            </div>

          </template>
        </el-table-column>

      </el-table>
      <el-pagination :page-size="pageSize" background layout="total, prev, pager, next, sizes, jumper"
        :total="totalpage" :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
      </el-pagination>
    </div>
    <CreateReport :isShow='false' :type="5" ref="CreateReport"  :dialogVisible = 'dialogVisible' @click="saveCreate()" @clearCreate="clearCreate()"></CreateReport>
  </div>
</template>
<style lang="less" scoped>
/deep/ .el-progress-bar__outer{
  // height: 10px !important;
}
/deep/ .el-progress-bar__innerText{
  vertical-align: super !important;
  color: #fff  !important;
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
    import { task } from '@/api/task.js'
    import { vulscan } from '@/api/vulscan.js'

    import CreateReport from "@/components/CreateReport.vue";
    export default {  
        name:'task',
        components:{
          CreateReport
        },
        props:{ 
            group_id:{},
            group_name:{},
        },
        data(){
            return{
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
                    risk: 0
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
            this.$store.state.activefirstMenu = '/vulScanTask';
            
            if(localStorage.getItem('shouquan')== 'false'){
              this.isShowShouQuan = false;
            }else{
              this.isShowShouQuan = true;
            }
          this.timeIDFOrSHouQuan =  setInterval(() => {
              if(localStorage.getItem('shouquan')== 'false'){
              this.isShowShouQuan = false;
            }else{
              this.isShowShouQuan = true;
            }
            }, 60000);
        },
        mounted(){
            this.getTaskEnum();
            let _this = this;
            this.getData(); 
            this.timer = setInterval(function(){
              _this.getData(true);
            },5000)
        },
        beforeDestroy () { 
          // 清除定时器
          clearInterval(this.timeIDFOrSHouQuan);
          this.timer ? clearInterval(this.timer) : null;
        },
        methods:{
          clickButton(type) {
            switch (type) {
              case '任务风险':
                this.$refs.loglistRef.toggleMenu();
                break;
            }
          }, 
          async getTaskEnum(){
              const res = await task.taskEnum();
              if(res.code == 200){
                this.risklevellist = res.data.riskLevel;
                this.risklevellist.unshift(
                    { label: "全部", value: 0 }
                ) 
                 
              }else{

              }
          },
            async getData(notloading){
              let multipleSelection = [];
              if (notloading) {
                multipleSelection = this.multipleSelection;
              }
             try { 
                const res = await vulscan.taskList({
                    page: this.formData.page_num,
                    search: this.formData.search_field, 
                    risk: this.formData.risk === 0 ? '' : this.formData.risk,
                    size: this.pageSize,
                    startTime:this.formData.time == null? '': this.formData?.time[0],
                    endTime: this.formData.time == null? '': this.formData?.time[1],
                });
                if(res.code == 200){
                    this.tableData = res.data.list;
                    this.totalpage = res.data.total;
                    // 解决 刷新的时候，已经勾选的行，可以依旧勾选上
                    if (notloading && this.tableData) {
                      let ids = [];
                      multipleSelection.forEach(item => {
                        ids.push(item.id);
                      });
                      this.$nextTick(() => {
                        this.tableData && this.tableData.forEach(item => {
                          if (ids.includes(item.id)) {
                            this.$refs.myTable.toggleRowSelection(item, true);
                          }
                        });
                      });
                    } 
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
              if(!this.group_id){
                this.$router.push({
                    path: `/vulScanTaskCreate`,
                    query: { 
                        flag: 1, 
                    }
                });
              }else{
                this.$router.push({
                    path: `/createtask`,
                    query: { 
                        flag: 1, 
                        group_id:this.group_id,
                        group_name:this.group_name,
                        type:2, //任务组详情进去
                    }
                });
              }
                
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
                this.formData.risk = 0;
                this.pageSize = 10;
                this.currentpage = 1;
                this.getData();
            },
            btnTaskinfo(row){ //详情
              this.$router.push({
                    path: `/vulScanTaskDetail`,
                    query: { 
                        id: row.id, 
                        name:row.taskName,
                    }
                });
                localStorage.setItem("taskTab", "tabs1");

            },
            async btnTaskstop(taskid){ //结束

              const loading = this.$loading()
              const dt = await vulscan.taskStop({
                id: taskid,
              });
              loading.close();

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
            async btnMultiDeleteTask(){
              if (this.multipleSelection.length == 0) return; 
              let _ids = this.multipleSelection.map(item => item.id); 
              
              const res = await	vulscan.taskDelete({
                ids: _ids.join(",")
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
              const res = await vulscan.taskDelete({ ids: scope.row.id + '' })
       
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
            handleSelectionChange(val){
              this.multipleSelection = val;
            },
            mouseenter(row, colum, cell, event) {
                this.showOperateButton = true;
                this.rowId = row.id   //赋值行id，便于页面判断  
            },
            mouseleave(row, colum, cell, event) { 
              let t = this.$refs['popover_id-' + row.id].showPopper;
              if(!t){
                    this.showEditFileNameButton = false;
              this.rowId = "";
              }

                // if (row.executeType == 3 || row.status != 1){  
                //     if (!this.$refs['popover-' + row.id]) {
                //     this.showOperateButton = false;
                //     this.rowId = "";
                //     return;
                //     }else{
                //     let isShow = this.$refs['popover-' + row.id].showPopper;
                //     if (!isShow) {
                //         this.showOperateButton = false;
                //         this.rowId = "";
                //     }
                //     }  
                // }else{ 
                //     if (!this.$refs['popover_id-' + row.id]){
                //     this.showOperateButton = false;
                //     this.rowId = "";
                //     return;
                //     }else{
                //     let isShow = this.$refs['popover_id-' + row.id].showPopper;
                //     if (!isShow) {
                //         this.showOperateButton = false;
                //         this.rowId = "";
                //     }
                //     } 
                // } 
            },
            btnReport(row){ //点击报告 
              this.dialogVisible = true;
              this.$refs.CreateReport.getinit(row.id,row.name);
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