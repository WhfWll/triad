<template>
    <!-- 逻辑漏洞 -->
    <div>
        <div class="main-title" >
            <label class="taskSceneBtn" >逻辑漏洞</label>
        </div>
        <div class="tasklist context_box_bg">
            <div class="search-box">
                <div class="operationbutton"> 
                <el-button type="primary" size="small" @click="btnNewTask" >新建</el-button>
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
                <div class="search-text">
                    <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select" size="small"
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
                <el-table-column prop="name" label="任务名称" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <el-link @click="btnTaskinfo(scope.row)">{{
                        scope.row.name
                        }}</el-link>
                    </template>
                </el-table-column>
                <el-table-column prop="typeName" label="检测类型">
                 
                </el-table-column>
                <el-table-column   label="任务风险" show-overflow-tooltip> 
                    <template slot-scope="scope">
                        <span :class="[
                        { 'riskstyle risk_hight': scope.row.risk == '1' },
                        { 'riskstyle risk_middle': scope.row.risk == '2' },
                        { 'riskstyle risk_low': scope.row.risk == '3' },
                        { 'riskstyle risk_nofind': scope.row.risk == '4' }
                        ]"><i></i>{{ scope.row.riskName }} </span> 
                    </template>
                </el-table-column>
                <el-table-column prop="updateTime" label="更新时间"> </el-table-column>
                <el-table-column prop="status" label="状态">
                    <template slot-scope="scope"> 
                        <div v-if="showOperateButton && rowId == scope.row.id  "> 
                                <el-link class="link_primary" :underline="false"  v-if="scope.row.status == 4" @click="btnReport(scope.row)">报告</el-link>
                                <el-link class="link_primary" :underline="false" v-if="scope.row.status != 4" @click="btnTaskstop(scope.row)">结束</el-link>
                                <el-link class="link_primary" :underline="false" @click="btnCopyTask(scope.row)">复制</el-link>
                                <el-popover placement="bottom" width="170" :visible-arrow="false" :ref="`popover_id-${scope.row.id}`"
                                popper-class="delButton_popper">
                                    <p class="delText">
                                    <i class="el-icon-warning"></i>确定删除吗？
                                    </p>
                                    <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                    <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                </div> 
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference">删除
                                </el-link>
                            </el-popover> 
                        </div>
                        <div v-else> 
                            <span :class="[ 
                                { 'tag_status tag_danger1': scope.row.status == 1 } ,
                                { 'tag_status tag_warning': scope.row.status == 2 },
                                { 'tag_status tag_primary': scope.row.status == 3 },
                                { 'tag_status tag_success': scope.row.status == 4 }]"><i></i>{{ scope.row.statusName }}</span>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination :page-size="pageSize" background layout="total, prev, pager, next, sizes, jumper"
                :total="totalpage" :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>
        </div>
        <CreateReport :isShow='false' :type="3" ref="CreateReport"   
            :dialogVisible = 'dialogVisible' @click="saveCreate()" 
            @clearCreate="clearCreate()"></CreateReport>
    </div>
</template>
<style scoped lang="less">
.tasklist{
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
</style>
<script>
import { logic } from '@/api/task.js'
import CreateReport from "@/components/CreateReport.vue";
export default {
    name:'logicvuln',
    components:{
        CreateReport
    },
    data(){
        return{
            alldelvisible:false,
            formData:{
                search:'',
                page:1,
            },
            pageSize:10,
            multipleSelection:[],
            Loading:false,
            tableData:[],
            currentpage:1,
            totalpage:0,
            showOperateButton:false,
            rowId:'',
            dialogVisible:false,
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/logicvuln';
    },
    mounted(){
        this.getData();
    },
    methods:{
        async getData(){
            const res = await logic.getTasklist({
                page:this.formData.page,
                size:this.pageSize,
                search:this.formData.search
            })
            if(res.code == 200){
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }
            else{

            }
        },
        handleReset(){
            this.formData.page_num = 1;
            this.formData.search= ""; 
            this.pageSize = 10;
            this.currentpage = 1;
            this.getData();
        },
        handlesearch(){
            //搜索
            this.formData.page = 1;
            this.getData();
            this.currentpage = 1;
        }, 
        btnNewTask(){
            this.$router.push({
                path: `/createlogicvuln`,
                query: { 
                    flag: 1, 
                }
            });
        },
        btnTaskinfo(row){
            this.$router.push({
                path: `/logicvulnDetail`,
                query: { 
                    id: row.id, 
                    name:row.name,
                }
            });  
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        mouseenter(row, colum, cell, event){
            this.showOperateButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断   
        },
        mouseleave(row, colum, cell, event){
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
        handleSizeChange(t){
            this.formData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        currentchange(t){
            this.formData.page = t;
            this.getData();
            this.currentpage = t;
        }, 
        async btnDel(scope){ 
            const res = await logic.taskDel({ ids: scope.row.id + '' })
       
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
        async btnMultiDeleteTask(){
            if (this.multipleSelection.length == 0) return; 
              let _ids = this.multipleSelection.map(item => item.id); 
              
              const res = await	logic.taskDel({
                ids: _ids.join(",")
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
        async btnTaskstop(row){
            const res = await logic.logicTaskStop({
                id:row.id
            })
            if (res.code == 200) {
                this.$message({
                  message: "结束任务成功",
                  type: "success"
                }); 
                this.getData();
              } else {
                this.$message({
                  message: res.msg,
                  type: "error"
                });
              }
        },
        async btnCopyTask(row){ //复制
            this.$router.push({
                path: `/createlogicvuln`,
                query: { 
                    flag: 2, 
                    id: row.id,  
                }
            }); 
        },
        btnReport(row){ //报告
            this.dialogVisible = true;
            this.$refs.CreateReport.getinit(row.id,row.name);
        },
        saveCreate(params){ //生成成功，  
            this.dialogVisible = false; 
        },
        clearCreate(){
            this.dialogVisible = false;
        },
    }
}
</script>