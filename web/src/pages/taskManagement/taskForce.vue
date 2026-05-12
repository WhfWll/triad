<template>
    <!-- 流量分析 -->
    <div>
        <div class="main-title  ">
            <label for="">任务组</label>
        </div>
        <div class="trafficlist context_box_bg">
            <div class="search-box">
                <div class="operationbutton">
                    <el-button type="primary" size="small" @click="btnNewTraffic">新建任务组</el-button>
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" style="padding-left:8px"
                        trigger="click" :visible-arrow="false" v-model="alldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消</el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDelete">确定</el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除</el-button>
                    </el-popover>
                </div>
                <div class="serach-condition"> 
                    <div class="search-text">
                        <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch" v-model="formData.search_field"
                            class="input-with-select" size="small" clearable>
                        </el-input>
                        <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                    </div>
                </div>
            </div>

            <el-table ref="myTable" :data="tableData" v-loading="Loading" @selection-change="handleSelectionChange"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" style="width: 100%">
                <el-table-column type="selection" width="55"> </el-table-column>
            
                <el-table-column prop="taskName" label="任务名称" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <el-link @click="btninfo(scope.row)">{{
                        scope.row.taskName
                        }}</el-link>
                    </template>
                </el-table-column> 
                <el-table-column prop="risk_level" label="任务风险"  >
                    
                    <template slot-scope="scope">
                        <span :class="[  
                            {'riskstyle risk_hight': scope.row.riskLevel == '2' },
                            {'riskstyle risk_middle': scope.row.riskLevel == '3' },
                            {'riskstyle risk_low': scope.row.riskLevel == '4' },
                            { 'riskstyle risk_nofind': scope.row.riskLevel == '5' }
                        ]"><i></i>{{ scope.row.riskLevelName }} </span>
                       
                    </template>
                </el-table-column>
                <el-table-column prop="" label="漏洞"> 
                    <template slot-scope="scope">
                        <span class="tag_status tag_danger bug_status">{{scope.row.vulNum[0]}}</span>
                        <span class="tag_status tag_warning bug_status">{{ scope.row.vulNum[1] }}</span>
                        <span class="tag_status tag_primary bug_status">{{ scope.row.vulNum[2] }}</span>
                        <span class="tag_status tag_success bug_status">{{ scope.row.vulNum[3] }}</span>
                    </template>
                </el-table-column>
                <el-table-column prop="updateTime" label="更新时间"> </el-table-column>
                <el-table-column prop="task_status" label="状态">
                    <template slot-scope="scope">
                        <!-- 2待开始 3运行中 4已结束 --> 
                        <div v-if="showOperateButton && rowId == scope.row.id  "> 
                            <!-- 待开始 2 -->
                            <div v-if="scope.row.status == 2">
                                <el-link class="link_primary" :underline="false" @click="btnTaskstop(scope.row.id)">结束</el-link> 
                                <el-popover placement="bottom" width="170" :visible-arrow="false" :ref="`popover_id-${scope.row.id}`"
                                    popper-class="delButton_popper">
                                    <p class="delText">
                                        <i class="el-icon-warning"></i>确定删除吗？
                                    </p>
                                    <div style="text-align: right; margin: 0">
                                        <el-button size="mini" class="delCancel" @click="btnCancelDel(scope)">取消
                                        </el-button>
                                        <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                    </div>
                                    <span slot="reference" style="cursor: pointer;vertical-align: bottom;">删除</span>
                                </el-popover>
                                <el-button>详情</el-button>
                            </div>
                            <!-- 运行中 3 -->
                            <div v-else-if="scope.row.status == 3"> 
                                <el-link class="link_primary" :underline="false" @click="btnTaskstop(scope.row.id)">结束</el-link>
                                <el-popover placement="bottom" width="170" :visible-arrow="false" :ref="`popover_id-${scope.row.id}`"
                                    popper-class="delButton_popper">
                                    <p class="delText">
                                        <i class="el-icon-warning"></i>确定删除吗？
                                    </p>
                                    <div style="text-align: right; margin: 0">
                                        <el-button size="mini" class="delCancel" @click="btnCancelDel(scope)">取消
                                        </el-button>
                                        <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                    </div>
                                    <span slot="reference" style="cursor: pointer;vertical-align: bottom;">删除</span>
                                </el-popover>
                            </div>
                            <!-- 已结束 4 -->
                            <div v-else-if="scope.row.status == 4"> 
                                <span style="margin:0 10px">详情</span>
                                <span>编辑</span>
                                <span style="margin:0 10px">新建任跳转过去</span>
                                <el-popover placement="bottom" width="170" :visible-arrow="false" :ref="`popover_id-${scope.row.id}`"
                                    popper-class="delButton_popper">
                                    <p class="delText">
                                        <i class="el-icon-warning"></i>确定删除吗？
                                    </p>
                                    <div style="text-align: right; margin: 0">
                                        <el-button size="mini" class="delCancel" @click="btnCancelDel(scope)">取消
                                        </el-button>
                                        <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                    </div>
                                    <span slot="reference" style="cursor: pointer;vertical-align: bottom;">删除</span>
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
            <el-pagination :page-size="pageSize" background layout="total, prev, pager, next, sizes, jumper" :total="totalpage"
                :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>
        </div>
    </div>
</template>
<style lang="less" scoped>
.trafficlist{
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
/deep/ .el-table td.el-table__cell div{
    line-height: 20px;
}
</style>
<script>
import { traffic } from '@/api/traffic.js'
export default {
    name:'taskForce',
    data(){
        return{
            formData:{
                search_field:'',
                page_num:1,

            },
            tableData:[],
            multipleSelection:[],
            alldelvisible:false,
            Loading:false,
            pageSize:10,
            totalpage:0,
            currentpage:1,
            showOperateButton:false,
            rowId:'',
            userID:'',
        }
    },
    created() {
        this.$store.state.activefirstMenu = "/taskForce";
        this.userID = this.$commonjs.decryptCBC(localStorage.getItem('user_id'),this.$commonjs.myKey); 
    },
    mounted() {
        this.getData();
    },
    methods: { 
        async getData(){
            const res = await traffic.getData({
                search: this.formData.search_field,
                size:this.pageSize,
                page:this.formData.page_num
            });
            
            if(res.code == 200){
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            } 
        },
        btninfo(row){ //详情
            this.$router.push({
                path: `/trafficinfo`,
                query: {
                    id: row.id,
                    name: row.taskName,
                    // page_num: this.formData.page_num,
                    // risk_level: row.risk_level
                }
            });
            localStorage.setItem('trafficTab','tabs1');
        },
        btnNewTraffic(){
            this.$router.push({
                path: `/addtraffic`,
                query: {
                    // dialogtitle: "新建",
                    // taskid: row.id,
                    // flag: 1,
                    // page_sendVal: true,
                    // titleName: "流量分析"
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
            this.formData.task_status = "";
            this.formData.risk_level = 0;
            this.pageSize = 10;
            this.currentpage = 1;
            this.getData();
        },
        async btnMultiDelete(){
            if (this.multipleSelection.length == 0) return;
            var ids = [];
            for (var i = 0; i < this.multipleSelection.length; i++) {
                ids.push(this.multipleSelection[i].id);
            }
            const res = await traffic.delTraffic({
                flowTaskIds: ids.join(','),
            })
            if (res.code == 200) {
                this.$message({
                    message: '删除成功',
                    type: 'success'
                });
                this.getData();
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }

        },
        async btnDel(scope){ //单个删除
            const res = await traffic.delTraffic({
                flowTaskIds:scope.row.id,
            })
            if(res.code == 200){
                this.$message({
                    message: '删除成功',
                    type: 'success'
                });
                this.getData();
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        btnCancelDel(scope) {
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
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
        async btnTaskstop(id){ //结束

            const res = await traffic.trafficStatus({
                flowTaskId:id,
                operate:'stop',
                userId:this.userID
            })
            if(res.code == 200){
                this.$message({
                    message: '任务结束',
                    type: 'success'
                });
                this.getData();
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            } 
        },
        currentchange(t){
            this.formData.page_num = t;
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t){
            this.formData.page_num = 1;
            this.pageSize = t;
            this.getData();
        },
    },
}
</script>