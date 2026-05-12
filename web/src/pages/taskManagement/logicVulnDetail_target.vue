<template>
    <div class="taskTargetList">
        <div class="search-box">
            <!-- <div class="operationbutton"> 
                <el-popover popper-class="delButton_popper" placement="bottom-start" width="170"
                    trigger="click" :visible-arrow="false" v-model="alldelvisible">
                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                    <div style="text-align: right; margin: 0" class="">
                        <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消
                        </el-button>
                        <el-button size="mini" type="primary" @click="btnMultiDelete">确定
                        </el-button>
                    </div>
                    <el-button  type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除
                    </el-button>
                </el-popover>  
            </div>
            <div class="serach-condition">
                <div class="search-text">
                    <el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearch"  v-model="fromData.search" class="input-with-select" size="small"
                        clearable> </el-input>
                    <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                </div>
                <div>
                    <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                </div>
            </div> -->
        </div>
        <el-table ref="targetTable" :data="tableData" tooltip-effect="dark" style="width: 100%" 
        @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" >
            <!-- <el-table-column type="selection" width="55">
            </el-table-column> -->
              <el-table-column prop="targetUrl" label="测试目标">  
            </el-table-column> 
            <el-table-column prop="typeName" label="漏洞类型">
            </el-table-column> 
            <el-table-column prop="risk_level" label="目标风险"> 
                <template slot-scope="scope">
                    <span
                        :class="[ 
                        { 'riskstyle risk_hight': scope.row.risk == 1 } ,
                        { 'riskstyle risk_middle': scope.row.risk ==2 },
                        { 'riskstyle risk_low': scope.row.risk == 3 },
                        { 'riskstyle risk_nofind': scope.row.risk ==4 }]"><i></i>{{scope.row.riskName}}</span>
                </template>
            </el-table-column>  
             <el-table-column label="存活状态" prop="isAliveName"> 
            </el-table-column> 
            <el-table-column prop="statusName" label="状态" width="200" fixed="right"> 
                <template slot-scope="scope"> 
                    <div v-if="showOperateButton && rowId == scope.row.id  ">
                        <el-link class="link_primary" v-if="scope.row.status == 4" :underline="false" @click="btnReport(scope.row)">报告</el-link>
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
        <CreateReport :isShow="is_show" :type="4" ref="CreateReport"  :dialogVisible = 'dialogVisible' @click="saveCreate()" @clearCreate="clearCreate()"></CreateReport>
    </div>
</template>
<style lang="less" scoped>
    .taskTargetList{
        background: #fff;
        padding: 24px 24px;
        box-sizing: border-box; 
        box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
        border-radius: 4px;
    }
</style>
<script>
import { logic } from '@/api/task.js'
import CreateReport from "@/components/CreateReport.vue";
export default {
    name:'',
    props:{ 
        task_id:{},
        task_name:{},
    },
    components:{
          CreateReport
    },
    data(){
        return{
            alldelvisible:false,
            tableData:[],
            pageSize:10,
            totalpage:0,
            currentpage:1,
            fromData:{
                page:1,
                search:'',
            },
            multipleSelection:[],
            showOperateButton:false,
            rowId:'',
            dialogVisible:false,
            is_show:false,
        }
    },
    created(){

    },
    mounted(){

    },
    methods:{
        async getData(){
            const res = await logic.getTargetlist({
                page:this.fromData.page,
                size:this.pageSize,
                search:this.fromData.search,
                task_id:this.task_id
            })
            if(res.code == 200){
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }
        },
        handleSizeChangetarget(t){
            this.formData.page = 1;
            this.currentPage = 1;
            this.pageSize = t;
            this.getData();
        },
        currentchange(t){
            this.formData.page = t;
            this.currentPage = t;
            this.getData();
        },
        // btnMultiDelete(){
            
        // },
        handleReset(){
            this.fromData.page_num = 1;
            this.fromData.search= ""; 
            this.pageSize = 10;
            this.currentpage = 1;
            this.getData();
        },
        handlesearch(){
            //搜索
            this.fromData.page = 1;
            this.getData();
            this.currentpage = 1;
        },  
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        mouseenter(row){
            this.showOperateButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断 
        },
        mouseleave(row){
            this.showOperateButton = false;
            this.rowId = "";
            // if (!this.$refs['popover_id-' + row.id]) {
            //     this.showOperateButton = false;
            //     this.rowId = "";
            //     return;
            // } else {
            //     let isShow = this.$refs['popover_id-' + row.id].showPopper;
            //     if (!isShow) {
            //         this.showOperateButton = false;
            //         this.rowId = "";
            //     }
            // }
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
    }
}
</script>