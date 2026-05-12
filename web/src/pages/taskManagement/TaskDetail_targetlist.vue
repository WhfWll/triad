<template>
    <div class="taskTargetList">
        <div  class="tasktarget_box">
            <div class="search-box">
                <div class="operationbutton"> 
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170"
                        trigger="click" :visible-arrow="false" v-model="alldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消
                            </el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDeleteTarget">确定
                            </el-button>
                        </div>
                        <el-button  type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除
                        </el-button>
                    </el-popover> 
                    <el-button :disabled="!multipleSelection.length" style="margin-left:10px" type="primary" size="small" @click="btnReport(multipleSelection)">生成报告</el-button>
                    <!-- <h4 v-for="item in multipleSelection" :key="item">{{item}}</h4> -->
                </div>
                <div class="serach-condition">
                    <div class="search-text">
                        <el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearchTarget"  v-model="formData.target_url" class="input-with-select" size="small"
                            clearable> </el-input>
                        <el-button type="primary" size="small" @click="handlesearchTarget">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleResetTarget">重置</el-button>
                    </div>
                </div>
            </div>
        </div> 
        <el-table ref="targetTable" :data="tableData" tooltip-effect="dark" 
            style="width: 100%"
             height="calc(100% - 102px)"
            @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave"  
            @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55">
            </el-table-column>
              <el-table-column prop="targetUrl" label="测试目标"> 
                <template slot-scope="scope">
                    <el-link @click="btnTasktargetinfo(scope.row)">
                        
                        {{scope.row.targetUrl}}
                    </el-link>
                </template>
            </el-table-column> 
            <el-table-column prop="opSys" label="系统">
            </el-table-column>
            <el-table-column prop="openPort" label="开放端口">
                <template slot-scope="scope">
                    {{ scope.row.openPort.join(',') }}
                </template>
            </el-table-column> 
            <el-table-column prop="risk_level" label="目标风险">
                <template slot="header">
                    <!-- <div style="position: relative;height:10px;"> -->
                        <span class="cursorPointer" @click="clickButton('目标风险')"
                            :class="(formData.riskLevel !== '' && formData.riskLevel !== 0) ? 'active' : ''">目标风险<i
                                class="iconfont iconshaixuan"></i>
                        </span>
                    <!-- </div> -->
                    <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.riskLevel" clearable
                        size="small" ref="list" @change="handlesearchTarget">
                        <el-option v-for="(item,index) in risklevellist" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                    </el-select>
                </template>
                <template slot-scope="scope">
                    <span
                        :class="[ 
                        { 'riskstyle risk_hight': scope.row.riskLevel == 1 } ,
                        { 'riskstyle risk_middle': scope.row.riskLevel ==2 },
                        { 'riskstyle risk_low': scope.row.riskLevel == 3 },
                        { 'riskstyle risk_nofind': scope.row.riskLevel ==4 }]"><i></i>{{scope.row.riskLevelName}}</span>
                </template>
            </el-table-column>  
             <el-table-column label="漏洞数量" width="280">
                <template slot="header">
                    漏洞 
                </template>
                <template slot-scope="scope">
                    <span class="tag_status tag_danger bug_status" v-if="scope.row.vulNum[0] >99">99+</span>
                    <span class="tag_status tag_danger bug_status"
                        v-else>{{ scope.row.vulNum[0] }}</span>
                    <span class="tag_status tag_warning bug_status" v-if="scope.row.vulNum[1] >99">99+</span>
                    <span class="tag_status tag_warning bug_status"
                        v-else >{{scope.row.vulNum[1]}}</span>
                    <span class="tag_status tag_primary bug_status" v-if="scope.row.vulNum[2] >99">99+</span>
                    <span class="tag_status tag_primary bug_status"
                        v-else>{{scope.row.vulNum[2]}}</span>
                    <span class="tag_status tag_success bug_status" v-if="scope.row.vulNum[3] >99">99+</span>
                    <span class="tag_status tag_success bug_status"
                        v-else>{{scope.row.vulNum[3]}}</span> 
                </template>
            </el-table-column>
        
            <el-table-column prop="isAliveName" label="存活状态">
            </el-table-column>
            <el-table-column prop="statusName" label="状态" width="200" fixed="right"> 
                <template slot-scope="scope">
                    <div v-if="showOperateButton && rowId == scope.row.id  ">
                        <!-- 待开始、运行中、已结束 -->
                        <!-- //2022-8-15wm 罗确定只有 从路径图按钮跳转时  进入渗透路径页面 -->
                        <el-link class="link_primary" v-if="scope.row.status ==4" :underline="false" @click="btnReport(scope.row)">报告</el-link>
                        <el-link :underline="false" v-if="scope.row.status !=4" class="link_danger"
                            @click="targetStop(scope.row.id)">结束
                        </el-link>
                        <el-link :underline="false"   class="link_danger"
                            @click="btnTasktargetinfo(scope.row)">路径图
                        </el-link>
                        <el-popover placement="bottom" width="170" :visible-arrow="false"
                            :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                            <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                            <div style="text-align: right; margin: 0">
                                <el-button size="mini" class="delCancel"
                                    @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">
                                    取消</el-button>
                                <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                            </div>
                            <el-link :underline="false" class="link_info linkafter" style="padding:0" slot="reference">
                                删除 </el-link>
                        </el-popover>
                    </div>
                    <div v-else>
                        <span class="tag_status tag_danger1" v-if="scope.row.status == 1">{{
                        scope.row.statusName}}</span>
                        <span class="tag_status tag_warning" v-if="scope.row.status == 2">{{
                        scope.row.statusName}}</span>
                        <span class="tag_status tag_primary" v-if="scope.row.status == 3">{{
                        scope.row.statusName}}</span>
                        <span class="tag_status tag_success" v-if="scope.row.status == 4">{{
                        scope.row.statusName}}</span>
                        <span class="tag_status tag_danger" v-if="scope.row.status == 5">{{
                        scope.row.statusName}}</span>
                    </div>
                </template>

            </el-table-column>
        </el-table>
        <el-pagination :page-size="pageSize" background layout=" total, prev, pager, next, sizes, jumper"
            :total="totalpage" :current-page="currentpage" @current-change="currentchange"
            @size-change="handleSizeChangetarget">
        </el-pagination>
        <CreateReport :isShow="is_show" :type="2" ref="CreateReport"  :dialogVisible = 'dialogVisible' @click="saveCreate()" @clearCreate="clearCreate()"></CreateReport>
    </div>
</template>
<style lang="less" scoped>
:deep(.el-table__body){
    .cell{
        line-height: 20px;
    }
}
.taskTargetList{
    height: 100%;
    background: #fff;
    padding: 24px 24px;
    box-sizing: border-box; 
    box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
    border-radius: 4px;
}
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

    // >span {
    //   position: absolute;
    // }
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
</style>
<script>
import { task } from '@/api/task.js'
import CreateReport from "@/components/CreateReport.vue";
import login from '@/api/login';
export default {
    name:'targetlist',
    components:{
          CreateReport
    },
    props:{ 
        task_id:{},
        task_name:{},
    },
    data(){
        return{
            is_show:false,
            dialogVisible:false,
            alldelvisible:false,
            formData:{
                target_url:'',
                riskLevel:0,
                page:1
            },
            multipleSelection:[],
            pageSize:10,
            totalpage:0,
            tableData:[],
            risklevellist:[],
            currentpage:1,
            showOperateButton:false,
            rowId:'',
        }
    },
    created(){

    },
    mounted(){
        this.getTaskEnum();
        this.getData();
    },
    methods:{
        clickButton(type){
            switch (type) {
              case '目标风险':
                this.$refs.list.toggleMenu();
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
                this.$message({
                  message: res.msg,
                  type: "error"
                });
              }
          },
        async getData(){
            const res = await task.getTargetlist({
                "taskId": this.task_id,  
                "search":this.formData.target_url,  
                "riskLevel":this.formData.riskLevel, 
                "page":this.formData.page, 
                "size": this.pageSize 
            })
            if(res.code == 200){
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }else{
                this.$message({
                  message: res.msg,
                  type: "error"
                });
            }
        },
        handlesearchTarget(){
            this.formData.page = 1;
            this.currentPage = 1;
            this.getData();
        },
        handleResetTarget(){
            this.pageSize = 10;
            this.formData.page = 1;
            this.formData.target_url = ''; 
            this.formData.riskLevel = 0;
            this.currentpage = 1;
            this.getData();
        },

        async btnMultiDeleteTarget(){
            if (this.multipleSelection.length == 0) return; 
              let _ids = this.multipleSelection.map(item => item.id); 
              
              const res = await	task.targetdel({
                taskId:this.task_id,
                targetIds: _ids.join(",")
              }); 
              if (res.code == 200) {
                this.$message({
                  message: "删除任务成功",
                  type: "success"
                }); 
                this.currentpage = 1;
                this.formData.page = 1;
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
            const res = await task.targetdel({
                taskId:this.task_id,
                targetIds:scope.row.id
            })
            if(res.code == 200){
                this.$message({
                    message: "删除任务成功",
                    type: "success"
                  });
                  scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
              
                  this.getData();
            }else{
                this.$message({
                    message: res.msg,
                    type: "error"
                });
            }
        },
        async targetStop(_id){ //目标结束
            const res = await task.targetStop({
                targetId:_id,
                operate:'stop'
            });
            if(res.code == 200){
                this.$message({
                    message: "测试目标结束成功",
                    type: "success"
                });
                this.getData();
            }else{
                this.$message({
                    message: res.msg,
                    type: "error"
                });
            }
        },
        currentchange(t){
            this.formData.page = t;
            this.getData();
        },
        handleSizeChangetarget(t){
            this.formData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        mouseenter(row){
            this.showOperateButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断 
        },
        mouseleave(row){
            if (!this.$refs['popover_id-' + row.id]) {
                this.showOperateButton = false;
                this.rowId = "";
                return;
            } else {
                let isShow = this.$refs['popover_id-' + row.id].showPopper;
                if (!isShow) {
                    this.showOperateButton = false;
                    this.rowId = "";
                }
            }
        },
        btnReport(row){ //点击报告 
        console.log(row,'点击报告 点击报告 点击报告 点击报告 ',Object.prototype.toString.call(row));
        if (Object.prototype.toString.call(row)=='[object Array]' && row.length>1) {
            this.is_show = true;
        }else{
            this.is_show = false;
        }
            this.dialogVisible = true;
            let ids = ''
            let targetUrl = ''
            let batchConfigJson = {}
            if (Object.prototype.toString.call(row)=='[object Array]') {
                row.forEach(element => {
                    ids=ids+ element.id+','
                    targetUrl=targetUrl+ element.targetUrl+','
                    batchConfigJson[element.id] = element.targetUrl
                });
                ids = ids.substring(0, ids.length - 1)
                targetUrl = targetUrl.substring(0, targetUrl.length - 1)
                // console.log(ids,'===============',targetUrl);
                this.$refs.CreateReport.getinit(ids,targetUrl,batchConfigJson);
            }else{
                    batchConfigJson[row.id] = row.targetUrl
                    let idNu = row.id+''
             this.$refs.CreateReport.getinit(idNu,row.targetUrl,batchConfigJson);
            }
        },
        async saveCreate(params){ //生成成功，  
            this.dialogVisible = false; 
        },
        clearCreate(){
            this.dialogVisible = false;
        },
        btnTasktargetinfo(row){ //点击目标详情
            this.$router.push({
                path: '/targetDetail',
                query: {
                    id: row.id,
                    name: row.targetUrl, 
                    taskname: this.task_name,
                   
                    taskId: this.task_id,  
                }
            });
        },
    }
}
</script>