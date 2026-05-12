<template>
    <div>
        <!-- 待检测漏洞 -->
        <div class="tasktarget_box">
            <div class="target_list">
                <div class="search-box">
                    <div class="operationbutton"> 
                        <el-button type="primary" style="margin-right:8px" size="small" @click="testvultest()"  
                            :disabled="!multipleSelectionbug.length">测试</el-button>
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
                            <el-input placeholder="搜索漏洞名称与漏洞位置" @keydown.enter.native="handlesearchbug" v-model="formDatabug.keyword" class="input-with-select"
                                size="small" clearable> </el-input>
                            <el-button type="primary" size="small" @click="handlesearchbug">搜索</el-button>
                        </div>
                        <div>
                            <el-button type="primary" size="small" @click="handleResetbug">重置</el-button>
                        </div>
                    </div>
                </div>
                <el-table  :data="buglisttableData" tooltip-effect="dark" height="calc(100% - 102px)"
                    ref="myTable" style="width: 100%" @selection-change="handleSelectionChangebug"
                    @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                    <el-table-column type="selection" width="55">
                    </el-table-column>
                    <el-table-column prop="targetUrl" label="测试目标" :show-overflow-tooltip="true"> 
                    </el-table-column>
                    <el-table-column prop="name" label="漏洞名称" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <el-link @click="btnbuginfo(scope.row)">{{scope.row.name}}</el-link>
                        </template>
                    </el-table-column>
                    <el-table-column prop="typeName" label="漏洞类型" :show-overflow-tooltip="true">
                        <template slot-scope="scope" slot="header">
                            <div style="position: relative;height:10px;">
                                <span class="cursorPointer" @click="clickButton('漏洞类型')"
                                    :class="(formDatabug.type !== '' && formDatabug.type !== '0') ? 'active' : ''">漏洞类型<i
                                        class="iconfont iconshaixuan"></i>
                                </span>
                            </div>
                            <el-select popper-class="thSelect" style=" width:150px;" v-model="formDatabug.type"
                                size="small" ref="vulTypeSelect" @change="handlesearchbug">
                                <el-option v-for="(item,i) in vulTypelist" :key="i" :label="item.label" :value="item.value">
                                </el-option>
                            </el-select>
                        </template>
                    </el-table-column>
                    <el-table-column prop="riskName" label="漏洞风险">
                        <template slot-scope="scope"  slot="header">
                            <div style="position: relative;height:10px;">
                                <span class="cursorPointer" @click="clickButton('漏洞风险')"
                                    :class="(formDatabug.risk_level !== '' && formDatabug.risk_level !== '0') ? 'active' : ''">漏洞风险<i
                                        class="iconfont iconshaixuan"></i>
                                </span>
                            </div>
                            <el-select popper-class="thSelect" style=" width:150px;" v-model="formDatabug.risk_level"
                                size="small" ref="vulRiskSelect" @change="handlesearchbug">
                                <el-option v-for="(item,i) in vulrisklist" :key="i" :label="item.label" :value="item.value">
                                </el-option>
                            </el-select>
                        </template>
                        <template slot-scope="scope">
                            <!-- 1-致命/2-高危/3-中危/4-低危/5-信息 -->
                            <span :class="[ 
                            { 'riskstyle risk_hight': scope.row.risk == 1 } ,
                            { 'riskstyle risk_middle': scope.row.risk == 2 },
                            { 'riskstyle risk_low': scope.row.risk ==3 },
                            { 'riskstyle risk_nofind': scope.row.risk ==4 }]"><i></i>{{
                                scope.row.riskName}}</span>
                        </template>
                    </el-table-column>

                    <el-table-column prop="location" :show-overflow-tooltip="true" label="漏洞地址">

                    </el-table-column>
                    <el-table-column prop="statusName" :show-overflow-tooltip="true" label="漏洞状态"> 
                    </el-table-column>
                    <el-table-column prop="testStatusName" label="测试状态">   
                        <template slot-scope="scope">
                            <div v-if="showOperateButton && rowId == scope.row.id">
                                <el-link :underline="false" class="link_primary" @click="btnbuginfo(scope.row)"
                                    style="padding-left: 0px;">详情</el-link>

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
                            <div v-else>
                                <span :class="[
                                { 'tag_status tag_primary': scope.row.testStatus == 1 },
                                { 'tag_status tag_warning': scope.row.testStatus == 2 },
                                { 'tag_status tag_danger': scope.row.testStatus == 3 },
                                { 'tag_status tag_success': scope.row.testStatus == 4 }]">{{ scope.row.testStatusName }}
                                    </span>
                            </div>
                        </template>
                    </el-table-column> 

                </el-table>

                <el-pagination background @size-change="handleSizeChangebug" @current-change="handleCurrentChangebug"
                    :current-page="currentPage" :page-size="pageSize" layout="total,  prev, pager, next, sizes, jumper"
                    :total="total">
                </el-pagination>
            </div> 
            <!-- 漏洞详情 -->
            <vulnmsginfo v-model="sendVal" :vulninfo=vulninfo :task_id="task_id"   >
            </vulnmsginfo>
          
        </div>
    </div>

</template>
<style lang="less" scoped>
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
        height: 100%;
        padding: 24px;
        background: #fff;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12); 
        box-sizing: border-box; 
        border-radius: 4px;
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
</style>
<script>
import vulnmsginfo from './Vulnmsginfo.vue';
import {
    task, 
} from '@/api/task.js';
import { vulnerability } from "@/api/tool.js";
export default {
    name:'',
    components:{
        vulnmsginfo, 
    },
    data(){
        return {
            sendVal: false,
            vulninfo:{},
            multipleSelectionbug:[],
            currentPage:1,
            total:0,
            pageSize:10,
            formDatabug:{
                risk_level:'', 
                type:'',
                search:'',
                keyword:'',
                page:1,
                status:'',
            },
            vulrisklist:[],
            vulTypelist:[],
            buglevel:[],
            showOperateButton:false,
            rowId:'',
            buglisttableData:[],
            bugalldelvisible:false,
        }
    },
    props:{ 
        task_id:{},
    },
    created(){ 
        // this.getSelectlist(); 
    },
    mounted(){ 
    },
    methods:{
        //漏洞风险类型列表
        async getSelectlist(){
            const res = await vulnerability.getVulObjectlist();

            if (res.code === 200) { 
                this.vulrisklist = res.data.risk; 
                this.vulTypelist = res.data.type; 
 

            }
            else {
                this.$message({
                    message: res.error,
                    type: 'error'
                });
            } 
        },
        clickButton(type) {
            switch (type) {
                case '漏洞类型':
                    this.$refs.vulTypeSelect.toggleMenu();
                    break;
                case '漏洞风险':
                    this.$refs.vulRiskSelect.toggleMenu(); 
                    break;
                case '状态':
                    this.$refs.status.toggleMenu();
                    break;
            }
        }, 
        async getData(target_ids, notloading) {

            let multipleSelection = notloading ? this.multipleSelectionbug : []; 
          
            let _target_ids = '';
            if (target_ids === undefined || target_ids ===''){ 
                _target_ids = '';
            }else{
                _target_ids = target_ids   ;
            }
            const res = await task.taskVulList({  
                page: this.formDatabug.page,
                size: this.pageSize,
                taskId: this.task_id,
                dataType:2, //1-漏洞测试,2-待测漏洞
                type: this.formDatabug.type,
                risk: this.formDatabug.risk_level, 
                search: this.formDatabug.keyword, 
            }) ;
            if (res.code === 200) { 
                this.buglisttableData = res.data.list;
                this.total = res.data.total;

                // 解决 刷新的时候，已经勾选的行，可以依旧勾选上
                if (notloading) {
                    let ids = [];
                    multipleSelection.forEach(item => {
                        ids.push(item.id);
                    });
                    this.$nextTick(() => {
                        this.buglisttableData.forEach(item => {
                            if (ids.includes(item.id)) {
                                this.$refs.myTable.toggleRowSelection(item, true);
                            }
                        });
                    });
                } 
            }
            else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
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
        btnMultiDeleteTarget(){

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
        btnbuginfo(row){  
            this.vulninfo = row; 
            this.sendVal = true;   
        }, 
        async testvultest(){ //测试
            if (this.multipleSelectionbug.length == 0) return; 
              let _ids = this.multipleSelectionbug.map(item => item.id); 

            const res = await task.testvultest({
                taskVulIds:_ids.join(","),
                taskId:this.task_id
            });
            if(res.code == 200){
                this.$message({
                  message: '测试成功',
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
    },
}
</script>