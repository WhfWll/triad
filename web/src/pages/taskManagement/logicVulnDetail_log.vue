<template>
    <div class="loglist">
        <div  class="log_box">
            <div class="search-box"> 
                <div class="serach-condition">
                    <div class="search-text">
                        <el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearch"  v-model="formData.search" class="input-with-select" size="small"
                            clearable> </el-input>
                        <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                    </div>
                </div>
            </div>
        </div> 
        <el-table ref="logTable" :data="logtableData" tooltip-effect="dark" style="width: 100%"
            @selection-change="handleSelectionChange" 
            @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" >
            <!-- <el-table-column type="selection" width="55">
            </el-table-column>  -->
            <el-table-column prop="target_url" label="测试目标">
            </el-table-column>  
             <el-table-column prop="createTime" label="创建时间">
                </el-table-column>
            <el-table-column prop="startTime" label="开始时间">
            </el-table-column>
            <el-table-column prop="endTime" label="结束时间">
            </el-table-column>
             <el-table-column prop="isAliveName" label="存活状态">
            </el-table-column>
            <el-table-column prop="status" label="测试状态">
                <template slot-scope="scope">
                    <div v-if="showOperateButton && rowId == scope.row.targetId">
                        <el-link class="link_primary" :underline="false" @click="btnShow(scope.row)">详情</el-link>
                    </div>
                    <div v-else>
                        <span :class="[ 
                            { 'tag_status tag_danger1': scope.row.status == 1 } ,
                            { 'tag_status tag_warning': scope.row.status == 2 },
                            { 'tag_status tag_primary': scope.row.status == 3 },
                            { 'tag_status tag_success': scope.row.status == 4 },
                        { 'tag_status tag_danger': scope.row.status ==5 }]">{{
                            scope.row.statusName}}</span>
                    </div>
                </template>
            </el-table-column>
        </el-table>
        <el-pagination background @size-change="handleSizeChangelog" @current-change="handleCurrentChangelog"
            :current-page="currentPage" :page-size="pageSize" layout="total,  prev, pager, next, sizes, jumper"
            :total="total">
        </el-pagination>

        <el-dialog :title="dialogtitle" :visible.sync="dialogVisible" width="1184px" :close-on-click-modal="false"
            :validate-on-rule-change="false" :show-close="false" >
            <div class="dialog_b_btn"> 
                <el-button size="small" @click="dialogVisible = false">关闭</el-button>
            </div>
            <div style="padding:24px">
                <div>
                    <ul class="loginfolist">
                        <li v-for="(item,index) in loginfolist" :key=index>
                            <span>[*]</span>
                            <span>[{{ item.createTime }}]</span>
                            <span>{{ item.pocname }}：</span>
                            <span>{{ item.result }}</span>
                        </li>
                    </ul>
                </div>
            </div>
        </el-dialog>  
    </div>
</template>
<style lang="less" scoped>
.loglist{
    background: #fff;
    padding: 24px 24px;
    box-sizing: border-box; 
    box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
    border-radius: 4px;
}
</style>
<script>
import { logic } from '@/api/task.js'
export default {
    name:'',
    props:{ 
        task_id:{},
    },
    data(){
        return{
            LoginfoID:null,
            multipleSelection:[],
            alldelvisible:false,
            formData:{
                search:'',
                page:1,
            },
            showOperateButton:false,
            rowId:'',
            currentPage:1,
            total:0,
            pageSize:10,
            logtableData:[],
            dialogtitle:'',
            dialogVisible:false,
            loginfolist:[],
        }
    },
    created(){

    },
    watch: {
        dialogVisible(newValue, oldValue) {
            if(newValue == false){
                clearInterval(this.LoginfoID);
            }
        }
    },
    methods:{
        async getData(){
            const res  = await logic.getLoglist({
                "task_id": this.task_id, 
                "search":this.formData.search, 
                "page": this.formData.page, 
                "size": this.pageSize
            })
            if(res.code== 200){
                this.logtableData = res.data.list;
                this.total = res.data.total;
            }else{
                this.$message({
                    message: res.msg,
                    type: "error"
                });
            }

        },
        handlesearch(){
            this.formData.page = 1;
            this.currentPage = 1;
            this.getData();
        },
        handleReset(){
            this.formData.page = 1;
            this.formData.search = ''; 
            this.currentPage = 1;
            this.getData();
        },
        handleSizeChangelog(t){
            this.formData.page = 1;
            this.currentPage = 1;
            this.pageSize = t;
            this.getData();
        },
        handleCurrentChangelog(t){ 
            this.formData.page = t;
            this.currentPage = t;
            this.getData();
        },
        
        btnShow(row){
            this.dialogVisible = true;
            this.dialogtitle = row.target_url;
            this.getLoginfo(row.id);
          this.LoginfoID = setInterval(() => {
                this.getLoginfo(row.id);
            }, 5000);
        },
        async getLoginfo(_id){
            const res = await logic.loginfo({
                id:_id
            })
            if(res.code == 200){
                this.loginfolist = res.data.list
            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }
        },
        mouseenter(row, colum, cell, event) { 
            this.showOperateButton = true; 
            this.rowId = row.targetId ;  //赋值行id，便于页面判断 
        },
        mouseleave(row, colum, cell, event) { 
            this.showOperateButton = false; 
            this.rowId='';
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
    }
}
</script>