<template>
    <div class="log_box">
        <div class="search-box">
            <div class="operationbutton">
                <el-button type="primary" style="margin-right:8px" size="small" @click="clearlog">清空</el-button>
            </div>
            <div class="serach-condition"> 
                <div class="search-text">
                    <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearchlog" v-model="formDatalog.search"
                        class="input-with-select" size="small" clearable> </el-input>
                    <el-button type="primary" size="small" @click="handlesearchlog">搜索
                    </el-button>
                </div>
                <div>
                    <el-button type="primary" size="small" @click="handleResetlog">重置
                    </el-button>
                </div>
            </div>
        </div> 
        <div class="loglist">  
            <el-table 
                ref="logTable" 
                :data="logtableData" 
                style="width: 100%" :show-header="false"
                 > 
                <el-table-column prop="content" label="日志内容">
                </el-table-column>  
            </el-table>
            <el-pagination background @size-change="handleSizeChangelog" @current-change="handleCurrentChangelog"
                :current-page="currentPage" :page-size="pageSize" layout="total,  prev, pager, next, sizes, jumper"
                :total="total">
            </el-pagination>

        </div>
        
    </div>
</template>
<style lang="less" scoped>
.loglist >div{
    color: rgba(72, 72, 102, 0.6400);
    font-size: 13px;
    margin-bottom: 8px;
} 
  
.log_box {
    padding: 24px 24px;
    box-sizing: border-box;
    position: relative;
    height: 100%;
    background: #fff;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    border-radius: 4px;

    .el-table__body-wrapper {
        height: calc(100% - 54px);
    }

    .link_primary {
        font-size: 13px;
    }
}
</style>
<script>
import { traffic } from '@/api/traffic.js'
export default {
    name: 'tasklog',
    props: {
        task_id: {},
    },
    data() {
        return {
            logloading: false,
            currentPage: 1,
            pageSize: 10,
            total: 0,
            logtableData: [],
            formDatalog: {
                search: '',
                page: 1, 
            },
            showOperateButton: false,
            rowId: '',
        }
    },
    mounted() {

    },
    methods: {
        async getData() {
            const res = await traffic.trafficlog({
                flowTaskId: this.task_id,
                search: this.formDatalog.search,
                page:this.formDatalog.page,
                size:this.pageSize,
            });
            if(res.code == 200){
                this.logtableData = res.data.list; 
                this.total = res.data.total;
            }
           
        },
        async clearlog(){
            const res = await traffic.clearTrafficlog({
                flowTaskId: this.task_id, 
            });
            if(res.code == 200){
                this.$message({
                    message:'清空日志成功',
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
        handlesearchlog() {
            this.formDatalog.page = 1;
            this.currentPage = 1;
            this.getData();
        },
        handleResetlog() {
            this.formDatalog.page = 1;
            this.formDatalog.search = '';
            this.currentpage = 1;
            this.getData();
        },
        handleSizeChangelog(t) {
            this.formDatalog.page = 1;
            this.currentpage = 1;
            this.pageSize = t;
            this.getData();
        },
        handleCurrentChangelog(t) {
            this.formDatalog.page = t;
            this.currentpage = t;
            this.getData();
        },
     
    },
}
</script>