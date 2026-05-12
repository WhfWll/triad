<template>
    <div>
        <div class="main-title"  >
            <label class="taskSceneBtn" >任务组</label>
        </div>
        <div class="taskgrouplist context_box_bg">
            <div class="search-box">
                <div class="operationbutton"> 
                <el-button type="primary" size="small" @click="btnNewTaskGroup"  >新建任务组</el-button>
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
                <!-- <div>
                    <el-date-picker v-model="formData.time" type="daterange" format="yyyy-MM-dd"
                    value-format="yyyy-MM-dd HH:mm:ss" :default-time="['00:00:00', '23:59:59']" range-separator="—"
                    start-placeholder="开始日期" end-placeholder="结束日期" size="small" clearable>
                    </el-date-picker>
                </div> -->
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
                :data="tableData" 
                @selection-change="handleSelectionChange"
                @cell-mouse-enter="mouseenter" 
                @cell-mouse-leave="mouseleave" 
                style="width: 100%"
            
                >
                    <el-table-column type="selection" width="55"> </el-table-column>

                    <el-table-column prop="name" label="任务组名称" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <el-link @click="btnTaskgroupinfo(scope.row)">{{
                            scope.row.name
                            }}</el-link>
                        </template>
                    </el-table-column>
                    <el-table-column   label="任务风险" show-overflow-tooltip> 
                        <template slot-scope="scope">
                            <span class="tag_status tag_danger bug_status">{{scope.row. high_num}}</span>
                            <span class="tag_status tag_warning bug_status">{{scope.row.middle_num}}</span>
                            <span class="tag_status tag_primary bug_status">{{scope.row.low_num}}</span>
                            <span class="tag_status tag_success bug_status">{{scope.row.safe_num}}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="update_time" label="更新时间"> </el-table-column>
                    <el-table-column prop="create_time" label="创建时间"> </el-table-column>
                    <el-table-column prop="statusName" label="状态">
                        <template slot-scope="scope">
                            <div v-if="showOperateButton && rowId == scope.row.id  "> 
                                <el-link class="link_primary" :underline="false" @click="btnUpdate(scope)">编辑</el-link>
                                <el-link class="link_primary" :underline="false" @click="btnCreateTask(scope)">新建任务</el-link>
                                <el-popover placement="bottom" width="170" :visible-arrow="false" :ref="`popover_id-${scope.row.id}`"
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
                            </div>
                            <div v-else> 
                                <span :class="[ 
                                    { 'tag_status tag_danger1': scope.row.status == 1 } ,
                                    { 'tag_status tag_warning': scope.row.status ==2 },
                                    { 'tag_status tag_primary': scope.row.status == 3 },
                                    { 'tag_status tag_success': scope.row.status ==4 },
                                { 'tag_status tag_danger': scope.row.status == 5 }]"><i></i>{{ scope.row.Status }}</span>
                                </div>
                        </template>
                    </el-table-column>
            </el-table> 
            <el-pagination :page-size="pageSize" background layout="total, prev, pager, next, sizes, jumper"
                :total="totalpage" :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>
        </div>
        <el-dialog :title="dialogTitle" :visible.sync="dialogaddFormVisible" :before-close="cancelform" width="900px;"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false"  >
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="groupform" label-width="0" status-icon ref="ruleForm" :rules="rules" >
                    <el-form-item label=" " prop="name">
                        <label class="dialog_item_label">任务组名称</label>
                        <el-input v-model="groupform.name" size="small" style="width:320px"  
                            ></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="unit">
                        <label class="dialog_item_label" style="vertical-align: top;">任务组描述</label>
                        <el-input type="textarea" v-model="groupform.describe" size="small" rows="4" resize="none" style="width:320px"  
                        ></el-input>
                    </el-form-item>
                    
                </el-form>
            </div>
        </el-dialog>



    </div>
</template>
<style lang="less" scoped>
.taskgrouplist{
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
</style>
<script>
import { task_group } from '@/api/task.js'
export default {
    name:'taskgroup',
    data(){
        return{
            alldelvisible:false,
            multipleSelection:[],
            formData:{
                page:1,
                search:'',
            },
            tableData:[],
            pageSize:10,
            currentpage:1,
            totalpage:0,
            showOperateButton:false,
            rowId:'',
            dialogaddFormVisible:false,
            rules:{},
            groupform:{
                name:'',
                describe:'',
            },
            dialogTitle:'新建任务组',
            groupId:'',
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/taskgroup';
    },
    mounted(){
        this.getData();
    },
    methods:{
        async getData(){
            const res = await task_group.grouplist({
                    page: this.formData.page,
                    search: this.formData.search,  
                    size: this.pageSize, 
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
        },
        btnNewTaskGroup(){
            this.dialogTitle = '新建任务组';
            this.dialogaddFormVisible = true;
        }, 
        btnCreateTask(scope){ //新建任务
            this.$router.push({
                path: `/createtask`,
                query: { 
                    flag: 1, 
                    group_id:scope.row.id,
                    group_name:scope.row.name,
                    type:1, //任务组列表进去
                }
            });
        },
        async btnMultiDeleteTask(){
            if (this.multipleSelection.length == 0) return; 
              let _ids = this.multipleSelection.map(item => item.id); 
              
              const res = await	task_group.delGroup({
                taskIds: _ids.join(",")
              }); 
              if (res.code == 200) {
                this.$message({
                  message: "批量删除任务组成功",
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
        handlesearch(){
              //搜索
            this.formData.page = 1;
            this.getData();
            this.currentpage = 1;
        },
        handleReset(){
            this.formData.page = 1;
            this.formData.search = ""; 
            this.pageSize = 10;
            this.currentpage = 1;
            this.getData();
        },
        btnTaskgroupinfo(row){ //任务组详情
            this.$router.push({
                path: `/taskGroupDetail`,
                query: { 
                    id: row.id, 
                    name:row.name,
                }
            });
            localStorage.setItem("groupTab", "tabs1");
               
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
        async btnCancelDel(){
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
            
        },
        async btnDel(scope){//单个删除
            let params = {
                id: scope.row.id
            }
            const data = await task_group.delGroup(params);
            if (data.code === 200) {
                this.$message({
                    message: '删除任务组成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose(); 
                this.getData();
            } else {
                this.$message({
                    message: data.msg,
                    type: 'error'
                });
            }
        },
        btnUpdate(scope){//编辑
            this.dialogTitle = '编辑任务组';
            this.dialogaddFormVisible = true;
            this.groupform.name = scope.row.name;
            this.groupform.describe = scope.row.describe;
            this.groupId = scope.row.id
        },
        async submitForm(){ //保存任务组 
            if(this.groupId){ //编辑
                this.groupform.id = this.groupId;
                const data = await task_group.updateGroup(this.groupform);
                if (data.code === 200) {
                this.$message({
                    message: '保存任务组成功',
                    type: 'success'
                });
                this.dialogaddFormVisible = false;
                this.getData();
                this.groupform.name='';
                this.groupform.describe='';
                this.groupform.id='';
            } else {
                this.$message({
                    message: data.msg,
                    type: 'error'
                }); 
            }
            }else{
                const data = await task_group.createGroup(this.groupform);
                if (data.code === 200) {
                this.$message({
                    message: '保存任务组成功',
                    type: 'success'
                });
                this.dialogaddFormVisible = false;
                this.getData();
                this.groupform.name='';
                this.groupform.describe='';
                this.groupform.id='';
            } else {
                this.$message({
                    message: data.msg,
                    type: 'error'
                }); 
            }
            }
           
            
        },
        cancelform(){
            this.dialogaddFormVisible = false;
            this.groupform.name='';
            this.groupform.describe='';
        },
    },
}
</script>