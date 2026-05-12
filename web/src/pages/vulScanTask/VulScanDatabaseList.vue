<template>
    <div>
        <div class="main-title  ">  
			CVE数据库管理
	  	</div> 
          <div class="tasktarget_box">
            <div class="target_list context_box_bg">
            <div class="search-box"> 
				<div class="serach-condition" > 
					<div class="search-text">
						<el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch"  v-model="formData.search" class="input-with-select"  size="small" clearable > </el-input>
						<el-button type="primary"  size="small" @click="handlesearch">搜索</el-button> 
					</div>
					<div >
						<el-button type="primary"  size="small" @click="handleReset">重置</el-button> 
					</div>   
				</div>  
			</div>
            <el-table
                ref="multipleTable"
                :data="tableData" 
                tooltip-effect="dark" 
                style="width: 100%"
                @selection-change="handleSelectionChange" @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column type="selection" width="55"> </el-table-column> 

                <el-table-column prop="cve" label="CVE编号"  :show-overflow-tooltip="true"> </el-table-column>
                
                <el-table-column prop="title_zh" label="标题"> </el-table-column>

                <el-table-column prop="cwe" label="CWE编号"> </el-table-column>

                <el-table-column prop="product" label="影响产品">
                    <template #default="{ row }">
                        <div class="custom-cell" :title="row.product">
                        {{ row.product }}
                        </div>
                    </template>
                </el-table-column> 

                <el-table-column prop="severity" label="漏洞级别">
                    <template slot-scope="scope">
                        <!-- 1-致命/2-高危/3-中危/4-低危/5-信息 -->
                        <span :class="[ 
                            { 'riskstyle risk_hight': scope.row.severity == 'CRITICAL' } ,
                            { 'riskstyle risk_middle': scope.row.severity == 'HIGH' },
                            { 'riskstyle risk_low': scope.row.severity == 'MEDIUM' },
                            { 'riskstyle risk_nofind': scope.row.severity == 'LOW' }]">
                            <i></i>
                            {{scope.row.severityName}}
                            <span v-if="scope.row.severityName">{{ scope.row.baseCvssv2Score }}</span>                        </span>
                    </template>
                </el-table-column> 

                <el-table-column prop="published_date" label="披露时间"> </el-table-column> 

                <el-table-column
                    prop="updated_at"
                    label="更新时间" >
                        <template slot-scope="scope"> 
                            <div v-if="showEditFileNameButton && rowId == scope.row.id">
                                <el-link :underline="false" class="link_primary" v-show="true" @click="btnDetail(scope.row)">详情</el-link>
                                <!-- <el-popover
                                    placement="bottom"
                                    width="170"   
                                    :visible-arrow="false"
                                    :ref="`popover_id-${scope.row.id}`"
                                    popper-class="delButton_popper" >
                                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                    <div style="text-align: right; margin: 0">
                                        <el-button size="mini" class="delCancel"  @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                        <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                    </div> 
                                    <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference" >删除</el-link>  
                                </el-popover>  -->
                            </div>
                            <div v-else > 
                                <span>{{scope.row.updated_at}}</span>
                            </div>
                        </template>
                </el-table-column>    
            </el-table>
            <el-pagination 
                background
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                :current-page="currentPage" 
                :page-size="pageSize"
                layout=" total,  prev, pager, next, sizes,jumper"
                :total="total">
            </el-pagination>
        </div>
        </div>

          <!-- 漏洞详情 -->
          <VulScanDatabaseList_detail v-model="sendVal" :vulninfo=vulninfo @saveData="handleSave()">
        </VulScanDatabaseList_detail>
    </div>
</template>
<style scoped lang="less">
.custom-cell {
    white-space: nowrap; /* 不换行 */
    overflow: hidden; /* 隐藏超出的内容 */
    text-overflow: ellipsis; /* 添加省略号 */
    max-width: 100px; /* 控制最大宽度 */
}

.link_danger:hover{
    color: rgba(72,72,102, 0.65)!important;
}

 .examplecode {
     color: #4C7AE3;
     padding-left: 10px;
     margin: 10px 0;
     font-style: italic;
     cursor: pointer;
 }
/deep/ thead {
    .cursorPointer {
        cursor: pointer;
        position: absolute;
        top: 6px;
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
/deep/ .el-checkbox__input.is-checked .el-checkbox__inner,
/deep/ .el-checkbox__input.is-indeterminate .el-checkbox__inner{
    background-color: #4C7AE3  !important;
    border-color: #4C7AE3 !important;
}

.el-table td.el-table__cell div{
    line-height: 20px;
}
.el-link.el-link--default:hover{
    color: #4C7AE3;
}
 
.el-button--primary.is-disabled{
    background-color: rgba(76, 122, 227, .5) !important;
    border-color: rgba(76, 122, 227, .2) !important;
}
    .updatestatus{
        /deep/ .el-dialog__body{
            height: 192px  !important;
        }
        /deep/ .el-dialog{
            height: auto !important;
        }
        /deep/ .el-dialog__body{
            padding: 72px 152px  !important;
        }
    }
    .tag_status{
        width: auto;
        padding: 0 8px;
    }
  
    .title_bg{
        width: 84px;
        height: 32px;
        font-size: 13px;
        font-weight: 500;
    }
    .title_bg1{
        background-color: rgba(243, 95, 40, 0.12) !important;
        border: 1px solid rgba(24, 144, 255, 0.08);     
        color: #F35F28 !important;
        border-left:3px solid #F35F28;
       
    }
    .title_bg2{
        background-color: rgba(76, 122, 227, 0.12) !important;
        border: 1px solid rgba(24, 144, 255, 0.08);      
        color: #4C7AE3 !important;
        border-left:3px solid #4C7AE3;
    }
  
    .message >div{ 
        // margin-bottom: 24px; 
        background: #F7F7FB;
        border-radius: 4px;
        border: 1px solid #E8E8F5;
        padding: 16px;
        box-sizing: border-box;
    }
    
    .message .title_bg{
        margin-bottom: 8px; 
    }
    .message >label{
        display: inline-block;
        width: 80px;
        text-align: center;
        height: 26px;
        line-height: 26px;
        color: #fff;
        background-color: #4c7ae3;
        font-weight: bold;
        font-size: 12px;
    }
    .message >div{
        height: 253px;
        overflow-y: auto;
    }
    .delButton_popper{
        padding: 16px !important;
        .el-button--mini{
            padding: 5px 10px;
            border-radius: 2px;
        }
    }
    .delText{
        margin-bottom: 16px ;
        color:rgba(72,72,102,0.64);
        i{
            color: #F9B640;
            margin-right: 10px;
        }
    }
    .controlbox{
        margin-top: 16px;
        color: rgba(72, 72, 102, 0.64);
        .cmdresult{
            padding: 16px 0;
            word-wrap: break-word; 
            word-break: normal; 
        }
        
    } 
    .useinput{
        width: 90% !important;
        box-sizing: border-box;
       
        /deep/ .el-input__inner{
             border:none !important; 
             padding-left:0;
        }
    }
    .tasktarget_box{ 
        box-sizing: border-box;
        position: relative;
        height: 100%;
        // background: #fff;
        .el-table__body-wrapper{
            height: calc(100% - 54px);
        }
    }
    .target_Statistics{
        height: 144px;
        padding: 24px;
        background: #fff;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
        box-sizing: border-box;
        margin-bottom: 15px;
        border-radius: 4px;
    }
     .target_list{ 
        padding: 24px;
        background: #fff;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12); 
        box-sizing: border-box;
        margin-bottom: 15px;
        border-radius: 4px;
    }
</style>
<script>
import {vulscan} from '@/api/vulscan.js' 
import VulScanDatabaseList_detail from './VulScanDatabaseList_detail.vue'

export default {
    name:'cveDBList',
    components:{
        VulScanDatabaseList_detail,
    },
    data(){
        return{
            total:0,
            currentPage:1,
            pageSize:10,
            showEditFileNameButton:false,
            rowId:'',
            multipleSelection:[],
            formData:{
                page:1,
                search:'',
            },
            tableData:[],
            sendVal: false,
            vulninfo:{}
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/vulScanDatabase';
    },
    mounted(){
        this.getData();
    },
    methods:{
        async getData(){
            const res = await vulscan.cveList({
                page:this.formData.page,
                size:this.pageSize,
                search:this.formData.search
            })
            if(res.code == 200){
                this.tableData = res.data.list; 
                this.total = res.data.total;

            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                }); 
            }
        },
        handleSave(){
            this.getData();
        },
        handleCurrentChange(t){
            this.formData.page = t; 
            this.getData();
            this.currentPage = t;
        },
        handleSizeChange(t){
            this.formData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handlesearch(){
            this.formData.page = 1;
            this.getData();
            this.currentpage = 1;
        },
        handleReset(){
            this.formData.search = ''; 
			this.formData.page = 1;
            this.pageSize=10;
            this.currentpage = 1;
            this.getData();
        },
        handleSelectionChange(val){
            this.multipleSelection = val
        },
        async btnDetail(row){
            this.vulninfo = row;
            this.sendVal = true;
        },
        mouseenter(row,colum,cell,event){ 
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
			this.showEditFileNameButton = false;
        }, 
    }
}
</script>