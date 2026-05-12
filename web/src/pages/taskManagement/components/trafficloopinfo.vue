<template>
    <!-- 被动流量 漏洞信息 -->
    <div class="tasktarget_box">
        <div class="target_list">
            <div class="search-box">
                <div class="operationbutton">
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" trigger="click"
                        :visible-arrow="false" v-model="bugalldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="bugalldelvisible = false">取消</el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDeleteTarget">确定</el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference"
                            :disabled="!multipleSelectionbug.length">删除</el-button>
                    </el-popover>
                </div>
                <div class="serach-condition">

                    <div class="search-text">
                        <el-input placeholder="搜索漏洞名称/漏洞位置/测试目标" @keydown.enter.native="handlesearchbug" v-model="formDatabug.keyword" class="input-with-select"
                            size="small" clearable> </el-input>
                        <el-button type="primary" size="small" @click="handlesearchbug">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleResetbug">重置</el-button> 
                    </div>
                </div>
            </div>
            <el-table  :data="buglisttableData" tooltip-effect="dark" v-loading="buglistcloading"
                ref="myTable" style="width: 100%" @selection-change="handleSelectionChangebug"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column type="selection" width="55">
                </el-table-column>
                <el-table-column prop="title" label="漏洞名称" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <el-link @click="btnbuginfo(scope.row)">{{scope.row.title}}</el-link>
                    </template>
                </el-table-column>
                <el-table-column prop="vul_risk_label" label="漏洞风险">
                    <template slot-scope="scope" slot="header">
                        <div style="position: relative;height:10px;">
                            <span class="cursorPointer" @click="clickButton('漏洞风险')"
                                :class="(formDatabug.risk_level !== '' && formDatabug.risk_level !== 0) ? 'active' : ''">漏洞风险<i
                                    class="iconfont iconshaixuan"></i>
                            </span>
                        </div>
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formDatabug.risk_level"
                            size="small" ref="vulSelect" @change="handlesearchbug">
                            <el-option v-for="(item,i) in vulrisklist" :key="i" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                        <span :class="[ 
                        { 'riskstyle risk_hight': scope.row.riskLevel == 1 } ,
                        { 'riskstyle risk_middle': scope.row.riskLevel == 2 },
                        { 'riskstyle risk_low': scope.row.riskLevel ==3 },
                        { 'riskstyle risk_nofind': scope.row.riskLevel ==4 }]"><i></i>{{scope.row.riskLevelName}}</span>
                    </template>
                </el-table-column>

                <el-table-column prop="host" :show-overflow-tooltip="true" label="漏洞位置">

                </el-table-column>
                <el-table-column prop="ip" label="测试目标">
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                          
                            <el-popover placement="bottom" width="170" :visible-arrow="false"
                                :ref="`popover-${scope.row.id}`" popper-class="delButton_popper">
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel"
                                        @click="scope._self.$refs[`popover-${scope.row.id}`].doClose()">取消
                                    </el-button>
                                    <el-button size="mini" type="primary" @click="btnbugdel(scope)">确定</el-button>
                                </div>
                                <el-link :underline="false" class="link_danger linkafter" slot="reference">删除</el-link>
                            </el-popover>
                        </div>
                        <div v-else>{{scope.row.ip}}  </div>
                    </template>
                </el-table-column>

            </el-table>

            <el-pagination background @size-change="handleSizeChangebug" @current-change="handleCurrentChangebug"
                :current-page="currentPage" :page-size="pageSize" layout="total,  prev, pager, next, sizes, jumper"
                :total="total">
            </el-pagination>
        </div>
      
        <!-- 漏洞详情 -->
        <singleloopinfo v-model="sendVal" :vuln_id=vuln_id :vulninfo=vulninfo :buglevel=buglevel :vulthreatlist=vulthreatlist
            :title="title"
            @saveData="handleSave()">
        </singleloopinfo>
       
       
    </div>
</template>

<script>
import singleloopinfo from './singleloopinfo.vue' 
import { traffic } from '@/api/traffic.js'

export default {
    name: "trafficloopinfo",
    components:{
        singleloopinfo,
        // vulnused
    },
    props:{ 
        trafficid:{},
    },
    data: () => ({  
        usedVal:false,
        vuln_id:'',
        title:'',
        vulninfo:{},
        dialogshow:false,
        sendVal: false,
        bugtypelist:[],
        buglevel:[],
        statuslist:[],
        upadtestatusdialogVisible:false,
        multipleSelectionbug:[],
        buglistcloading:false,
        formDatabug:{
            risk_level:'', 
            type:'',
            search:'',
            keyword:'',
            page:1,
            status:'',
        },
        pageSize:10,
        is_bugUpdate:false, 
        buglisttableData:[], 
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

    }),
    mounted() {
        // this.getSelectlist(); 
    },
    methods: { 
        async getData(){
            const res = await traffic.trafficVuln({
                page:this.formDatabug.page,
                search: this.formDatabug.keyword,
                size:this.pageSize,
                flowTaskId:this.trafficid,
                riskLevel:this.formDatabug.risk_level,
            });
            if(res.code == 200){
                this.buglisttableData = res.data.list;
                this.total = res.data.total;
            }

        },
        //漏洞风险类型列表
        async getEnum(){
            const res = await traffic.trafficEnum();
            if(res.code == 200){
                this.vulrisklist = res.data.vulRiskLevel; 
                this.vulrisklist.unshift({label:'全部',value:0});

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
        handleSave(){
            this.getData();
        },
        // async getbuglist(target_ids, notloading) {

        //     let multipleSelection = notloading ? this.multipleSelectionbug : []; 
        //     // if (!notloading) {
        //     //     this.buglistcloading = true;
        //     // }
        //     let _target_ids = '';
        //     if (target_ids === undefined || target_ids ===''){ 
        //         _target_ids = '';
        //     }else{
        //         _target_ids = target_ids   ;
        //     }
        //     const res = await taskVul.taskVulList({  
        //         page: this.formDatabug.page,
        //         page_size: this.pageSize,
        //         task_id: this.task_id,
        //         vul_risk: this.formDatabug.risk_level == 0 ? '': this.formDatabug.risk_level,
        //         status: this.formDatabug.status == 0 ? '' :this.formDatabug.status,
        //         target_id: _target_ids,
        //         search: this.formDatabug.keyword, 
        //     }) ;
        //     if (res.success) {
        //         // this.buglistcloading = false;
        //         this.buglisttableData = res.data.results;
        //         this.total = res.data.count;

        //         // 解决 刷新的时候，已经勾选的行，可以依旧勾选上
        //         if (notloading) {
        //             let ids = [];
        //             multipleSelection.forEach(item => {
        //                 ids.push(item.check_vul_id);
        //             });
        //             this.$nextTick(() => {
        //                 this.buglisttableData.forEach(item => {
        //                     if (ids.includes(item.check_vul_id)) {
        //                         this.$refs.myTable.toggleRowSelection(item, true);
        //                     }
        //                 });
        //             });
        //         } 
        //     }
        //     else {
        //         this.$message({
        //             message: res.error,
        //             type: 'error'
        //         });
        //     }
        // },
        handleResetbug(){ //重置
            this.formDatabug.type='';
            this.formDatabug.risk_level = '';
            this.formDatabug.keyword='';
            this.pageSize=10;
            this.formDatabug.page = 1;
            this.getData();
        },
        handlesearchbug(){
            this.getData();
        },
        handleSizeChangebug(t){
            this.formDatabug.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handleCurrentChangebug(t){
            this.formDatabug.page = t;
            this.getData();
        },
        handleSelectionChangebug(val){
            this.multipleSelectionbug = val;
        },
       
        async btnbugdel(scope){ //删除漏洞
            const res = await traffic.trafficvulndel({
                passive_traffic_risk_ids: scope.row.id + '', 
            });
            if (res.code == 200) {
                this.$message({
                    message: '删除漏洞信息成功',
                    type: 'success'
                });
                scope._self.$refs[`popover-${scope.row.id}`].doClose();
                this.getData();
            } else {
                this.$message({
                    message: res.error,
                    type: 'error'
                });
            }     
        },
        async btnMultiDeleteTarget(){ //批量删除漏洞
            if(this.multipleSelectionbug.length == 0) return;
    		let _ids = [];
            let vul_risks=[];
    		for (var i = 0; i < this.multipleSelectionbug.length; i++) {
                _ids.push(this.multipleSelectionbug[i].id); 
			}
         
            const res = await traffic.trafficvulndel({
                passive_traffic_risk_ids: _ids.join(','), 
            });
            if (res.code == 200) {
                this.$message({
                    message: res.msg,
                    type: 'success'
                }); 
                this.bugalldelvisible = false; 
                this.getData();
            } else {
                this.$message({
                    message: res.error,
                    type: 'error'
                });
            }   
        },
        startTest(){

        },
       
        testhover(){
            this.tt = false;
        },
        btnbuginfo(row){  
            this.vulninfo = row;
            this.vuln_id = row.id;
            this.dialogshow = true;
            this.sendVal = true;  
            this.title = row.title;
        }, 
  
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
           

            let t = this.$refs['popover-' + row.id].showPopper;
            if (!t) {
                this.showOperateButton = false;
                this.rowId = "";
            }
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
    // .bugbasicinfo{
    //     padding: 24px;
    //     background: #fff;
    //     border:1px solid rgba(232,232,245,1);
    // }
    // .buginfo_box{
    //     padding: 24px ;
    // }
    // .bugotherinfo{
    //     margin-top: 32px;
    //     .part_title{
    //         margin-bottom: 8px;
    //     }
    //     .content{
    //         background:rgba(255,255,255,1);
    //         border-radius:2px;
    //         border:1px solid rgba(232,232,245,1);
    //         padding: 12px 16px;
    //         color:rgba(72,72,102,0.64);
    //         font-size: 13px;
    //     }
    // }
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