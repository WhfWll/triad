<template>
    <!-- 二级 被动流量  -->
    <div class="tasktarget_box">
        <div class="target_list">
            <div class="search-box">
                <div class="operationbutton"> 
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" trigger="click"
                        :visible-arrow="false" v-model="bugalldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="bugalldelvisible = false">取消</el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDelete">确定</el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference"
                            :disabled="!multipleSelection.length">删除</el-button>
                    </el-popover>
                </div>
                <div class="serach-condition">

                    <div class="search-text">
                        <el-input placeholder="搜索URL、Title、IP" @keydown.enter.native="handlesearchbug" v-model="formData.keyword" class="input-with-select"
                            size="small" clearable> </el-input>
                        <el-button type="primary" size="small" @click="handlesearchbug">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleResetbug">重置</el-button> 
                    </div>
                </div>
            </div>
            <el-table  :data="tableData" tooltip-effect="dark" v-loading="buglistcloading"
                ref="myTable" style="width: 100%" @selection-change="handleSelectionChangebug"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column type="selection" width="55">
                </el-table-column>
                <el-table-column prop="method" label="方法" :show-overflow-tooltip="true"> 
                </el-table-column> 

                <el-table-column prop="respCode" :show-overflow-tooltip="true" label="状态码">
                </el-table-column>
                <el-table-column prop="url" :show-overflow-tooltip="true" label="URL">
                </el-table-column>
                <el-table-column prop="respTitle" :show-overflow-tooltip="true" label="Title">
                </el-table-column> 
                <el-table-column prop="tags" :show-overflow-tooltip="true" label="Tags">
                </el-table-column> 
                <el-table-column prop="ip" :show-overflow-tooltip="true" label="IP">
                </el-table-column>
                <el-table-column prop="respContentType" :show-overflow-tooltip="true" label="响应类型">
                </el-table-column>
 
                <el-table-column prop="createTime" label="请求时间"> 
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                            <el-link :underline="false" class="link_primary" @click="btninfo(scope.row)"
                                style="padding-left: 0px;" >详情</el-link> 
                        </div>
                        <div v-else>{{ scope.row.createTime }}
                        </div>
                    </template>
                </el-table-column>

            </el-table>

            <el-pagination background @size-change="handleSizeChangebug" @current-change="handleCurrentChangebug"
                :current-page="currentPage" :page-size="pageSize" layout="total,  prev, pager, next, sizes, jumper"
                :total="total">
            </el-pagination>
        </div>
        <el-dialog title="漏洞测试" :visible.sync="upadtestatusdialogVisible" width="1184px" :close-on-click-modal="false"
            :validate-on-rule-change="false"  :show-close="false" class="updatestatus">
            <div class="dialog_b_btn">
                <!-- <el-button size="small" @click="saveUpdateStatus">确定</el-button> -->
                <el-button size="small" @click="upadtestatusdialogVisible = false">关闭</el-button>
            </div>
           
            <el-row  style="height:100%; ">
                <el-col :span="12" style="margin-right:24px;height:100%"> 
                    <div class="message requestpack">
                        <label class="title_bg title_bg1">请求报文</label> 
                        <div style="  background: #fff;">
                            <pre>{{ requestpack }}</pre>
                        </div>
                    </div>
                </el-col>
                <el-col :span="12"  style="height:100%">
                    <div class="message">
                        <label class="title_bg title_bg2">响应报文</label>
                        <div style="   background: #fff;">
                            <pre>{{ responsepack }}</pre>
                        </div>
                    </div>
                </el-col>
            </el-row>
           
        </el-dialog>
    </div>
</template>

<script>
import {
    taskVul, 
} from '@/api/task.js';
// import { vulnerability } from "@/api/tool.js";
import { traffic } from '@/api/traffic.js'
export default {
    name: "passiveflux",
    props:{
        target_id:{}, 
        task_id:{},
    },
    data: () => ({  
        usedVal:false,
        vulninfo:{},
        dialogshow:false,
        sendVal: false,
        bugtypelist:[],
        buglevel:[],
        statuslist:[],
        upadtestatusdialogVisible:false,
        multipleSelection:[],
        buglistcloading:false,
        formData:{ 
            keyword:'',
            page:1, 
        },
        pageSize:10,
        is_bugUpdate:false, 
        tableData:[], 
        statusform:{
            status:'',
        }, 
        useinfo:{
            pocname:'',
            title:'',
            target:'',
            time:'',
        },
        TestdialogVisible:false,
        yzloading:false,
        testtitle:'',
        target_result:'',
        verify_result:[],
        pocname:'',
        target_result_id:'', 
        currentPage:1,
        total:0,
        bugalldelvisible:false,
        levellist:[],
        statusSellist:[],
        rowId:'',
        showOperateButton:false,
        vulrisklist:[],
        vulthreatlist: [],
        responsepack:'',
        requestpack:'',
    }),
    created() {
        
    },
    mounted() {
        
    },
    methods: { 
        async getData(){
             const res = await traffic.trafficlist({
                flowTaskId:this.task_id,
                page:this.formData.page,
                search:this.formData.keyword,
                size:this.pageSize
             });
             if(res.code == 200){
                 this.tableData = res.data.list;
                 this.total = res.data.total;

             }else{
                 this.$message({
                    message: res.msg,
                    type: 'error'
                });
             }
            
        },
        async btninfo(row){ //详情 
            const res = await traffic.gettrafficHeader({
                flowBaseId:row.id
            });
            if(res.code == 200){
                this.upadtestatusdialogVisible = true;
                this.requestpack = res.data.reqHeader;
                this.responsepack = res.data.respHeader;
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }


        },
        
        clickButton(type) {
            switch (type) {
                case '漏洞风险':
                    this.$refs.vulSelect.toggleMenu();
                    break;
                case '状态':
                    this.$refs.status.toggleMenu();
                    break;
            }
        }, 
        updateStatus(){
            this.upadtestatusdialogVisible = true;
        },  
        handleResetbug(){ //重置
            this.formData.type='';
            this.formData.risk_level = '';
            this.formData.keyword='';
            this.pageSize=10;
            this.formData.page = 1;
            this.getData();
        },
        handlesearchbug(){
            this.getData();
        },
        handleSizeChangebug(t){
            this.formData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handleCurrentChangebug(t){
            this.formData.page = t;
            this.getData();
        },
        handleSelectionChangebug(val){
            this.multipleSelection = val;
        },
        async btnMultiDelete(){ //批量删除
            if (this.multipleSelection.length == 0) return;
            var ids = [];
            for (var i = 0; i < this.multipleSelection.length; i++) {
                ids.push(this.multipleSelection[i].id);
            }
            const res = await traffic.trafficlistinfodel({
                passive_traffic_base_id: ids.join(',')
            })
            if(res.code == 200){
                this.$message({
                    message: '删除成功',
                    type: 'success'
                });
                this.bugalldelvisible = false;
                this.getData();
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        
        testhover(){
            this.tt = false;
        },
        // btnbuginfo(row){  
        //     this.vulninfo = row;
        //     this.upadtestatusdialogVisible = true;
        //     this.sendVal = true;   
        // }, 
        
        isJSON(str) {
            if (typeof str == 'string') {
                try {
                    var obj=JSON.parse(str);
                    if(typeof obj == 'object' && obj ){
                        return true;
                    }else{
                        return false;
                    }

                } catch(e) { 
                    return false;
                }
            } 
        },
        closeTest(){
            this.yzloading = false;
            this.TestdialogVisible = false; 
            clearInterval(this.timer);
        },
        mouseenter(row, colum, cell, event) {
            this.showOperateButton = true;
            this.rowId = row.id;  //赋值行id，便于页面判断 
        },
        mouseleave(row, colum, cell, event) {
           

            // let t = this.$refs['popover-' + row.id].showPopper;
            // if (!t) {
                this.showOperateButton = false;
                this.rowId = "";
            // }
        },
    },
};
</script>

<style lang="less" scoped>
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
.el-col-12 {
    width: 48%;
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
            // height: 192px  !important;
        }
        /deep/ .el-dialog{
            // height: auto !important;
        }
        /deep/ .el-dialog__body{
            padding:  24px  !important;
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
    .requestpack>div {
     background: #fff !important;
    //  padding: 0 !important;
 }
    .message{
        height: 100%;
        .title_bg {
            margin-bottom: 8px;
        }
    }
 
 .bugotherinfo {
     margin-top: 32px;

     .part_title {
         margin-bottom: 8px;
     }

     /* .content {
         background: rgba(255, 255, 255, 1);
         border-radius: 2px;
         border: 1px solid rgba(232, 232, 245, 1);
         padding: 12px 16px;
         color: rgba(72, 72, 102, 0.64);
         font-size: 13px;
     } */
 }
 .message>label {
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

 .message>div {
    //  height: 253px;
    height: calc(100% - 50px);
     overflow-y: auto;
     pre{
        white-space: pre-wrap;
        white-space: -moz-pre-wrap;
        white-space: -pre-wrap;
        white-space: -o-pre-wrap;
        word-wrap: break-word;
     }
 }
 .packbtn {
     height: 65px;
     box-sizing: border-box;
     padding: 16px;
     text-align: left;
     padding-left: 0;
 }

 .packinput {
     box-sizing: border-box !important;
     .packtxt {
         height: 100%;

         /deep/ textarea {
             height: 100%;
             border: 0 !important;
         }
     }
 }
   .message>div {
       background: #F7F7FB;
       border-radius: 4px;
       border: 1px solid #E8E8F5;
       padding: 16px;
       box-sizing: border-box;
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
     .target_list{ 
        padding: 24px;
        background: #fff;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12); 
        box-sizing: border-box;
        margin-bottom: 15px;
        border-radius: 4px;
    }
</style>