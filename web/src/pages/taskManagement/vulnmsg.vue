<template>
    <!-- 漏洞信息 -->
    <div class="tasktarget_box">
        <div class="target_list">
            <div class="search-box">
                <div class="operationbutton"> 
                    <!-- <el-button type="primary" style="margin-right:8px" size="small" @click="updateStatus"
                        :disabled="!multipleSelectionbug.length">变更状态</el-button> -->
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
            <el-table  :data="buglisttableData" tooltip-effect="dark" v-loading="buglistcloading" height="calc(100% - 102px)"
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
                    <template slot-scope="scope" slot="header">
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

                <el-table-column prop="location" :show-overflow-tooltip="true" label="漏洞位置">

                </el-table-column>
                <el-table-column prop="status" label="状态"> 
                    <!-- <template slot-scope="scope" slot="header">
                        <div style="position: relative;height:10px;">
                            <span class="cursorPointer" @click="clickButton('状态')"
                                :class="(formDatabug.status !== '' && formDatabug.status !== 0) ? 'active' : ''">状态
                                <i class="iconfont iconshaixuan"></i>
                            </span>
                            <el-tooltip class="item" effect="dark" placement="right">
                                <div slot="content"> 
                                    验证成功，通过POC验证发现的漏洞；<br />
                                    利用成功，通过EXP利用发现的漏洞；<br /> 
                                </div>
                                <i class="iconfont icontishi" style="position: absolute;top: 6px;left:56px"></i>
                            </el-tooltip>
                        </div>
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formDatabug.status" clearable
                            size="small" ref="status" @change="handlesearchbug">
                            <el-option v-for="(item, index) in statusSellist" :key="index" :label="item.label"
                                :value="item.value">
                            </el-option>
                        </el-select>
                    </template> -->
                    <template slot-scope="scope" slot="header">
                        <div style="position: relative;height:10px;">
                            <span class="cursorPointer" @click="clickButton('状态')"
                                :class="(formDatabug.status !== '' && formDatabug.status !== '0') ? 'active' : ''">状态<i
                                    class="iconfont iconshaixuan"></i>
                            </span>
                        </div>
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formDatabug.status"
                            size="small" ref="status" @change="handlesearchbug">
                            <el-option v-for="(item,i) in statusTypelist" :key="i" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                    <!-- <template slot-scope="scope" slot="header">
                        <div style="position: relative;height:10px;">
                            <span class="cursorPointer" style="cursor:initial">状态</span>
                            <el-tooltip class="item" effect="dark" placement="right">
                                <div slot="content"> 
                                    验证成功，通过POC验证发现的漏洞；<br />
                                    利用成功，通过EXP利用发现的漏洞；<br /> 
                                </div>
                                <i class="iconfont icontishi" style="position: absolute;top: 6px;left:26px"></i>
                            </el-tooltip>
                        </div>
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formDatabug.status" clearable
                            size="small" ref="status" @change="handlesearchbug">
                        </el-select>
                    </template>  -->
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                            <!-- :disabled="scope.rowisSnapshot != 1"  -->
                            <el-link :underline="false" class="link_primary" :disabled="scope.row.isSnapshot != 1"  @click="btnSnap(scope.row)"
                                style="padding-left: 0px;">截图</el-link>
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
                            { 'tag_status tag_primary': scope.row.status == 1 },
                            { 'tag_status tag_warning': scope.row.status == 2 },
                            { 'tag_status tag_danger': scope.row.status == 3 },
                            { 'tag_status tag_success': scope.row.status == 4 }]">{{ scope.row.statusName }}
                                </span>
                        </div>
                    </template>
                </el-table-column>
                <!-- <el-table-column prop="time" label="测试时间"> 
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                            <el-link :underline="false" class="link_primary" @click="btnbuguse(scope.row)"
                                style="padding-left: 0px;" v-if="scope.row.status[0] == 3">利用</el-link>

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
                            <span >2022-12-12 12:00</span>
                        </div>
                    </template>
                </el-table-column> -->

            </el-table>

            <el-pagination background @size-change="handleSizeChangebug" @current-change="handleCurrentChangebug"
                :current-page="currentPage" :page-size="pageSize" layout="total,  prev, pager, next, sizes, jumper"
                :total="total">
            </el-pagination>
        </div>
        <el-dialog title="变更状态" :visible.sync="upadtestatusdialogVisible" width="1184px" :close-on-click-modal="false"
            :validate-on-rule-change="false" :show-close="false" class="updatestatus">
            <div class="dialog_b_btn">
                <el-button size="small" @click="saveUpdateStatus">确定</el-button>
                <el-button size="small" @click="upadtestatusdialogVisible = false">关闭</el-button>
            </div>
            <div style="padding:24px">
                <div>
                    <el-form ref="statusform" :model="statusform" label-width="0" class="clearfix"
                        style="text-align: center;">
                        <el-form-item label="" style="margin-bottom:0">
                            <el-select v-model="statusform.status" placeholder="请选择" class="selstatus"
                                style="width:360px;height:40px">
                                <el-option v-for="(item,index) in statuslist" :key="index" :label="item.label"
                                    :value="item.value">
                                </el-option>
                            </el-select>
                        </el-form-item>
                    </el-form>
                </div>
            </div>
        </el-dialog>
        <!-- 漏洞详情 -->
        <vulnmsginfo v-model="sendVal" :vulninfo=vulninfo :task_id="task_id" :buglevel=buglevel :vulthreatlist=vulthreatlist
            @saveData="handleSave()">
        </vulnmsginfo>
        <!-- 漏洞利用 -->
        <vulnused v-model="usedVal" :task_id="task_id" :useinfo=useinfo  >

        </vulnused>
        <el-dialog :title="testtitle+'验证'" :visible.sync="TestdialogVisible" width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="startTest" v-loading="yzloading">开始验证</el-button>
                <el-button size="small" @click="closeTest()">关闭</el-button>
            </div>
            <div style="padding:24px">
                <div style="padding-bottom:32px;color:rgba(72, 72, 102, 0.64);border-bottom:1px solid  #E8E8F5">
                    {{target_result}}
                </div>
                <div style="margin-top:26px">
                    <div class="controlbox" v-for="(item,index) in verify_result" :key=index>

                        <div v-if="item.str">
                            {{item.str}}
                        </div>
                        <el-row :gutter="20" v-else>
                            <el-col :span="12">
                                <div class="message">
                                    <label class="title_bg title_bg1">请求报文</label>
                                    <div>
                                        <pre>{{item.back[0]}}</pre>
                                    </div>
                                </div>
                            </el-col>
                            <el-col :span="12">
                                <div class="message">
                                    <label class="title_bg title_bg2">响应报文</label>
                                    <div>
                                        <pre>{{item.back[1]}}</pre>
                                    </div>
                                </div>
                            </el-col>

                        </el-row>
                    </div>
                </div>

            </div>
        </el-dialog>
        
       

        <el-dialog title="截图" :visible.sync="dialogTableVisible">
               <div class="dialog_b_btn">
                <el-button size="small" @click="dialogTableVisible= false">关闭</el-button>
            </div>
             <img width="100%" height="99.5%" :src="snapBase" alt="Base64 Image" />
        </el-dialog>
    </div>
</template>

<script>
import vulnmsginfo from './Vulnmsginfo.vue'
import vulnused from './vulnused.vue'
import {
    task, 
} from '@/api/task.js';
import { vulnerability } from "@/api/tool.js";
export default {
    name: "vulnmsg",
    components:{
        vulnmsginfo,
        vulnused
    },
    props:{
        target_id:{}, 
        task_id:{},
    },
    data: () => ({  
        dialogTableVisible:false,
        usedVal:false,
        vulninfo:{},
        snapBase:'',
        // dialogshow:false,
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
        // useresult:[],
        // usecmd:'',
        // bugusedialogVisible:false, 
        TestdialogVisible:false,
        yzloading:false,
        testtitle:'',
        target_result:'',
        verify_result:[],
        pocname:'',
        target_result_id:'', 
        // bugdialogVisible:false,
        currentPage:1,
        total:0,
        bugalldelvisible:false,
        // responsepack:'',
        // requestpack:'',
        levellist:[],
        statusSellist:[],
        rowId:'',
        showOperateButton:false,
        vulrisklist:[],
        vulTypelist:[],
        vulthreatlist: [],
        statusTypelist: [
            {
                label: '验证存在',
                value: 2
            },{
                label: '利用成功',
                value: 3
            },
        ]

    }),
    created () {
        // this.getSelectlist();
        // this.getData('');
    },
    methods: { 
        //漏洞风险类型列表
        async getSelectlist(){
            const res = await vulnerability.getVulObjectlist();

            if (res.code === 200) {
                
                this.vulrisklist = res.data.risk; 
                this.vulTypelist = res.data.type; 
                // this.vulthreatlist = res.data.impact;
                let buglevel = [];
                for (var i = 0; i < this.vulrisklist.length;i++){
                    if( i != 0){
                        buglevel.push(this.vulrisklist[i]) 
                    }
                } 
                this.buglevel = buglevel;
                // this.statusSellist = res.data.status; 
                // this.statusSellist.unshift({ label: '全部', value: 0 });

                // let status = [];
                // for (var i = 0; i < this.statusSellist.length; i++) {
                //     if (i != 0) {
                //         status.push(this.statusSellist[i])
                //     }
                // } 
                // this.statuslist = status;

            }
            else {
                this.$message({
                    message: res.msg,
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
        updateStatus(){
            this.upadtestatusdialogVisible = true;
        }, 
        handleSave(){
            this.getData();
        },
        async getData(target_ids, notloading) {

            let multipleSelection = notloading ? this.multipleSelectionbug : []; 
            // if (!notloading) {
            //     this.buglistcloading = true;
            // }
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
                dataType:1, //1-漏洞测试,2-待测漏洞
                // taskId: 26,
                // vul_risk: this.formDatabug.risk_level == 0 ? '': this.formDatabug.risk_level,
                type: this.formDatabug.type,
                risk: this.formDatabug.risk_level,
                // status: this.formDatabug.status == 0 ? '' :this.formDatabug.status,
                // target_id: _target_ids,
                search: this.formDatabug.keyword,
                status: this.formDatabug.status
            }) ;
            if (res.code === 200) {
                // this.buglistcloading = false;
                // this.buglisttableData = res.data.results;
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
            this.formDatabug.status='';
            this.formDatabug.risk_level = '';
            this.formDatabug.keyword='';
            this.pageSize=10;
            this.formDatabug.page = 1;
            this.getData();
        },
        handlesearchbug(){
            this.formDatabug.page = 1;
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
            const res = await task.taskVulDelete({
                taskVulIds: scope.row.id + '',
                taskId: this.task_id 
            });
            if (res.code === 200) {
                this.$message({
                    message: '删除漏洞信息成功',
                    type: 'success'
                });
                scope._self.$refs[`popover-${scope.row.id}`].doClose();
                this.getData();
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }     
        },
        async btnMultiDeleteTarget(){ //批量删除漏洞
            if(this.multipleSelectionbug.length == 0) return;
    		let _ids = [];
    		for (var i = 0; i < this.multipleSelectionbug.length; i++) {
                _ids.push(this.multipleSelectionbug[i].id);
			}
         
            const res = await task.taskVulDelete({
                taskVulIds: _ids.join(','),
                taskId: this.task_id,
            });
            if (res.code === 200) {
                this.$message({
                    message: res.msg,
                    type: 'success'
                }); 
                this.bugalldelvisible = false; 
                this.getData();
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }   
        },
        async saveUpdateStatus(){
            if(this.multipleSelectionbug.length == 0) return;
    		let _ids = [];
    		for (var i = 0; i < this.multipleSelectionbug.length; i++) {
                // let _is_repeat = this.multipleSelectionbug[i].is_repeat ? 1:0;
    			// _ids.push({target_result_id:this.multipleSelectionbug[i].target_result_id,is_repeat:_is_repeat});

                _ids.push(this.multipleSelectionbug[i].id);


			} 

            const res = await task.updateVulStatus({
                id: _ids.join(','), 
                vul_status: this.statusform.status
            })
            if (res.success) {
                this.$message({
                    message: '变更状态成功',
                    type: 'success'
                });
                this.upadtestatusdialogVisible = false;
                this.getData();
            } else {
                this.$message({
                    message: dt.error,
                    type: 'error'
                });
            } 
        },
        testhover(){
            this.tt = false;
        },
        btnbuginfo(row){  
            this.vulninfo = row;
            // this.dialogshow = true;
            this.sendVal = true;   
        }, 
        // 点击截图
        btnSnap(row){  
           let that = this
        //    漏洞截图（新接口）
                that.$ajax.get('/smart/task/getvulsnapshot',{
                    params: { 
                        taskVulId:row.id,
                        
                    }
                }).then(dt=>{
                    let res = dt.data;
                    console.log(res,'000000aaaaaa');
                    if(res.code == 200){  
                        this.snapBase = res.data.snapshot
                        console.log(this.snapBase,'截图');

                        // this.snapBase = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAABQAAAAMgCAYAAAB8mM/7AAAAAXNSR0IArs4c6QAAIABJREFUeJzs3XuwneP58PErGTS0W0vYBnFoECTtVhQhYiTqEJSEihCmTjUYNM50YitFHEoZhCmDoUQlQQgaQeLUOFYFUUE0JEiQTSIHOd3vH/nt9dqyk7XCr/O+vXw+M3smWfdz3+tea6/9z3ee51ltSiklAAAAAICU2v6/3gAAAAAA8J8jAAIAAABAYgIgAAAAACQmAAIAAABAYgIgAAAAACQmAAIAAABAYgIgAAAAACQmAAIAAABAYgIgAAAAACQmAAIAAABAYv+rAfD999+veszChQvjzTffjNmzZ6/w+qWUGD9+fCxcuHCpscmTJ0dTU9M33t9nn30Ws2bN+kZzm5qa4qOPPmrx2OLFi+Pxxx+Pp556arl7qqa1tb/q3XffXe7rrjZ/6tSpMW3atG88Xk0tn4nWTJ8+PaZOnbrM8c8++ywmTZoUpZRvNL/a+Lfd3zfx+uuvx1133VXz8e+880589tlnSz3+n9gbAAAA8N/rWwfAUko88sgjsfPOO0ffvn2XedyMGTPiyCOPjP322y/GjBkT8+fPr4yNHz8+2rRpU/nZZJNNWl3jlltuia222iqmT59eeWzs2LGx4447xsYbbxxrrrlmHH300S3Wrra/6dOnR2NjY6yxxhoxbty4FX5tv/3tb+OAAw6I4447Lvr37x/z5s2LOXPmRP/+/WO33XaLRx55ZJnvSTWtrd3shRdeiC5dukTHjh1jzTXXjEGDBq3Q/Pfeey+23XbbOP3002OPPfaI/fbbr0VMqja+PLV+JlrT1NQUe+yxR6yzzjrRoUOH2HfffVsE37lz50b//v1jv/32iwEDBsQWW2wRb7zxRs3zq43fcMMNLT6LzT9nnnlmTfOrmTt3bqyzzjot1n7vvffi4osvjp/85Cdx9dVX17TOa6+9FptuummMHj265tcGAAAAfEeVb+mTTz4pw4YNKx07dixdu3Zt9Zhnn3221NfXlz/96U9l8eLFS4336dOnnH322eWiiy4qF110URkzZsxSx7z11lslIkpElKlTp5ZSSlmwYEHp2rVrueGGG8r9999funfvXiKiXHLJJTXvb/To0eXCCy8sEVFGjRq1Qq/t+uuvL/X19eXLL78sCxYsKA0NDeXss88upZSyePHiEhFl4MCBy3n3lm15a8+aNav06dOnvPjii+X5558v22yzTYmI8vHHH9c0v5RSunbtWi644IJSSimzZ88uHTp0KEcddVTN48tTy2diWY477rhy//33l4kTJ5Zf//rXJSLKyJEjK+Pnnntu2XXXXSv/79evX+nYsWOZN29eTfOXN7548eLSuXPncvTRR5crr7yyXHvtteWaa64pEVFGjx5d0/rVXHPNNaVfv36Vz/ott9xSGevTp09N79e8efNKp06dSkSUu+++u+b3DgAAAPhu+tYBsNlBBx3Uarz48MMPS319fTn//PNbnffCCy+UTp06lffee2+Zay9YsKB07969nHHGGS0C4Lhx48qTTz5ZOe7jjz8uEVF69epV8/5KKeWVV15pNQAub+6MGTNKRJTGxsbKY9dee22JiPLhhx+WUkqpq6v7RgGw2tozZ84sM2bMqIzdddddJSIq72G1+W+//XaJiPL8889Xxv/4xz+WiCjvvvtu1fFaLe89b83ixYvL+++/X/n/lClTSkSUe++9t5RSyqJFi0qHDh3KFVdcUTnmueeeKxFRbr/99qrzq41PnTp1qWD26quvlrq6ujJ//vyq86v54osvSn19fRk3blxZtGjRUuOHHXZYTe/X2WefXc4+++wWAfDb7g0AAADIq6ZLgJuamuLII4+MLl26xAYbbBCXX375Use0a9eu1blHHXVUzJ07N4499th4/vnnl7ofXWNjY0ycODE23HDD2GWXXeLdd99dao3LLrss9t5779hhhx1aPN61a9fo3r175f9rrbVWNDQ0RENDQ837i4j43ve+t8yxZc194oknIiJaPFeXLl0iIuLll1+OiIhVV101Fi1aFI2NjbHOOuvE9ttvXxmLiHjwwQejf//+sdVWW8WOO+4YzzzzTE1r19XVxRprrFEZe+WVV+Kcc86JDTbYoKb5zfcl/Op94rbYYouIWHIfumrjy9t7tfctImL27Nlx0kknxbbbbhvbbrttXHPNNRER0aZNm+jQoUPluFdffTV23XXX2HvvvSMiYtKkSTFlypSYMmVK5ZhOnTpFxJJLYqvNrza+3nrrxT777NNiryNGjIh+/frFyiuvXHV+xJLLn6+99trYfvvto0uXLnH66afHjBkzIiJi8ODBMX369Nhxxx1jo402ijFjxrR4rjZt2kRExNChQ2PzzTePTTbZJIYMGdLimKeeeipeeeWVOOmkk5aaW21vAAAAwHdTTQHwgAMOiF122SVef/31GDBgQJx55pnx2muvVZ03b968ePjhh2PWrFlxyCGHxOGHHx7rrrtu3HjjjZVjzjvvvPjzn/8cvXr1iqeeeiq22mqrFpHwpZdeitGjR8fpp59e9flmzpwZ48ePj4MPPriWl/WtNMexr4a41VdfPSKWfCFJs0GDBkVTU1P84he/iBdeeCH22muvmDNnTkydOjX23Xff6Nu3b7z88suxaNGiOPnkk1do7aamprjsssti0KBBsccee9S8t5133jkiIh5//PHKePO94ubMmVN1fHl7r8WBBx4Ybdu2jeeeey6OO+64OPnkkyvRMiJi/vz5ce+990bfvn2jZ8+escoqq0RERMeOHaO+vj4efPDBpfb11S9wWdb8Wse/6s4774wDDjig5vmXX3553HHHHTFq1KgYNmxYXHHFFXHVVVdFRESfPn3izjvvjOOPPz6mTJkSPXv2XOoekc8++2zcfvvt0bt375g0aVIceuih8Y9//CMiIj7//PM48cQT46abboq2bVv/012R1wYAAAB8N1QNgM8991yMHTs22rVrF0OHDq2Eh3vuuafq4hMmTIiIiIsvvjjGjBkTr776anTt2jVOO+20SrjZYYcd4je/+U2MHDkyzj///Jg1a1ZcdtllEbEkNh177LFx8803x0orrVT1+a677ro466yzYuutt6567Lc1d+7ciIj4wQ9+UHms/M830n71G45POumkuOaaa+KOO+6IgQMHxvTp0+Ohhx6KUkp07Ngxtt5662jbtm1su+22lTPbal3773//e+WMwh49elTOwqs2f9NNN42Ghoa49dZb4+GHH45///vfccstt0RExPrrr191fHl7r+bpp5+OUaNGxRFHHBErrbRS9OnTJw499NBYbbXVKse89dZb8eijj0bEkjNEmz8Pbdu2jcMOOywmTpwYV199dUybNi0GDx4cEREbbrhh1fm1jn/1uPfffz969OhR0/zZs2fHWWedFf3794811lgjttxyyzjllFOiY8eOERGx6aabxiGHHBKDBw+OBx54ICIiTjnllBZrd+7cOe6777649NJLY9SoURERlfd+wIAB0djYGOutt94y399aXxsAAADw3VE1ADZ/w2pdXV2sttpq0alTpxg5cmRNlxZOnDgxIiJ+9rOfRUTEKqusEv37949Zs2bF008/3XIjbdtGY2NjdO/evfJtvAMHDowuXbrE5MmTY+zYsfHqq69GRMQzzzyz1KXCb7/9dgwbNiwaGxur7ut/w1prrRUR0eIbh5vDW/NlqRERP/zhDyv/3muvvSJiyaWsHTp0iHfeeScmTJgQ++23X4wcObIyv9a199lnnxgyZEg8+eSTERFx77331jz/9ttvj4aGhjj44IPj1FNPjREjRkRExFZbbVV1fHl7r+aVV16JiIif/vSnlb3ecccdsd1221WO6dKlS1x33XXx5ptvRocOHeKOO+6ojP3ud7+Lfv36xbnnnhu9e/eOv/71rxERlbMWq82vZbzZiBEj4sADD1zqEvFlzX/77bcjImKbbbapHHvllVfGEUccsdTa++67b5x//vkxYcKE+OKLLyqPr7766pXIvssuu0TEkqh33333xfPPPx/t27ePsWPHVv5+XnvttRaXldf62gAAAIDvjqoB8PPPP4+IiB//+Mexzz77VH5qOcuu+Z50X375ZeWx5jjy5ptvtjqnd+/elXvHjRo1Km6//fbo0aNH9OjRI84///yIiOjbt2/cdtttlTmzZ8+OQw89NG666aYWZ5L9JzXfb635/m4RER9//HFERGy22WatzmmOayuvvHIsWrQojj/++Lj88svjpptuimOOOeYbr929e/fYfffdK8G1lvkNDQ3x9NNPx8yZM+OCCy6IiIgTTzwxvv/971cdX97eq5k5c2ZE/N+zQ5vNmzdvqWPXXXfdOOGEE2L8+PGVx9q3bx9DhgyJmTNnxhNPPBGffPJJdOrUKXbaaaea5q/I+N133x2/+tWvlvlavj6/OeT985//rPraIiJ69eoVEbHMs1vbtWsX2223Xayyyirx5JNPxoQJEyp/CwcddFBERFxwwQVL3Q+wltcGAAAAfHdUDYAbb7xxREQ89NBDlcc+/fTTms6023LLLSOiZexpPrup+Uslvu5f//pX5cynsWPHxtSpUys/N9xwQ0REvPjii3HqqadGxJKz3A4//PA477zzKlFy7ty5MW3atKr7+zZ23nnnqKurq9yfrXnvEUtiaWumT58eERE///nPY9iwYXHDDTfEX/7yl6ivr//Wa7dr1y46d+78jeZfddVV0bVr17j00ktbXfvr48vbezXNz//VL8CYPHnyMp971VVXbXFG3Vfd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AEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCY4pitFAAAJb0lEQVQAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAAAAJCYAAgAAAEBiAiAAAPB/2rEDGQAAAIBB/tb3+AojAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYE4AAAAAAMCYAAQAAAGBMAAIAAADAmAAEAAAAgDEBCAAAAABjAhAAAAAAxgQgAAAAAIwJQAAAAAAYC9HtYf/tln7sAAAAAElFTkSuQmCC"
                        this.dialogTableVisible =true
                    }else{
                        this.dialogTableVisible =false
                        that.$message({
                            message:res.msg,
                            type: 'error'
                        });
                    }
                    
                }).catch(err=>{})
        }, 
        fnResult(){
            let that = this;
            this.timer = setInterval(function(){
                that.$ajax.get('/task/vul/verify/result/',{
                    params: { 
                        target_id:that.target_id,
                        target_result_id:that.target_result_id,
                    }
                }).then(dt=>{
                    let res = dt.data;
                    if(res.success){  
                        that.target_result = res.data.target_result;
                        // that.verify_result = res.data.verify_result;
                        for(var i=0; i< res.data.verify_result.length;i++){
                        var item =  res.data.verify_result[i];
                        if(that.isJSON(item)){ //报文
                            var back = JSON.parse(item);  
                            that.verify_result.push({'str':'success：'+back.success});
                            for( var key in back ){ 
                                if(key != 'request' && key != 'response' && key != 'success'){ 
                                    that.verify_result.push({'str':key+'：url：'+back[key].url})
                                } 
                            }
                            that.verify_result.push({'back':[back.request[0],back.response[0]]})
                        }else{
                            that.verify_result.push({'str':item})
                        }
                        if(item == '验证结束'){
                            clearInterval(that.timer);
                        }
                    }

                    }else{
                        that.$message({
                            message:res.error,
                            type: 'error'
                        });
                    }
                    
                }).catch(err=>{})
            },5000);
        },
       
        btnbuguse(row){ //利用
            // this.bugusedialogVisible = true;
            this.useinfo.pocname = row.pocname;
            this.useinfo.id = row.check_vuln_id;
            this.useinfo.title = row.name;
            this.useinfo.target = row.title; 
            this.useinfo.time = this.commonjs.nowtime1();  
            this.usedVal = true;
            // this.$nextTick(() =>{
            //     document.getElementById('useinput').focus();
            // })
        }, 
         
        btnverification(row){ //验证
            this.verify_result = [];
            this.testtitle = row.name; 
            this.TestdialogVisible = true;
            this.pocname = row.pocname;
            this.target_id = row.target_id;
            this.target_result_id = row.target_result_id;
            this.$ajax.get('/task/vul/verify/result/',{
                params: { 
                    target_id:this.target_id,
                    target_result_id:this.target_result_id,
                }
            }).then(dt=>{
                let res = dt.data;
                if(res.success){  
                    this.target_result = res.data.target_result;
                    // this.verify_result = res.data.verify_result; 
                    for(var i=0; i< res.data.verify_result.length;i++){
                        var item =  res.data.verify_result[i];
                        if(this.isJSON(item)){ //报文
                            var back = JSON.parse(item);  
                            this.verify_result.push({'str':'success：'+back.success});
                            for( var key in back ){ 
                                if(key != 'request' && key != 'response' && key != 'success'){ 
                                    this.verify_result.push({'str':key+'：url：'+back[key].url})
                                } 
                            }
                            this.verify_result.push({'back':[back.request[0],back.response[0]]})
                        }else{
                            this.verify_result.push({'str':item})
                        }
                    }
                    // console.log(this.verify_result);
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
                
            }).catch(err=>{})

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
        startTest(){ //开始验证
            this.yzloading = true;
            this.$ajax({
                method:'post',
                url:'/task/vul/verify/',
                data: {
                    target_id:this.target_id ,
                    target_result_id:this.target_result_id,
                    pocname:this.pocname
                } 
            })
            .then(dt =>{
                let res = dt.data;
                this.yzloading = false;
                if(res.success){   
                    this.$message({
                        message:res.msg,
                        type: 'success'
                    });
                    this.fnResult();
                    // this.$ajax.get('/task/target/vuln/verify/result/',{
                    //     params: { 
                    //         target_id:this.target_id,
                    //         target_result_id:this.target_result_id,
                    //     }
                    // }).then(dt=>{
                    //     let res = dt.data;
                    //     if(res.success){  
                    //         this.target_result = res.data.target_result;
                    //         this.verify_result = res.data.verify_result;

                    //     }else{
                    //         this.$message({
                    //             message:res.error,
                    //             type: 'error'
                    //         });
                    //     }
                        
                    // }).catch(err=>{})
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            }).catch(err=>{})
        },
        closeTest(){
            this.yzloading = false;
            this.TestdialogVisible = false; 
            clearInterval(this.timer);
        },
        startTestResult(){ //开始验证--显示结果
            this.$ajax.get('/task/vul/verify/result/',{
                params: { 
                    target_id:this.target_id,
                    target_result_id:this.target_result_id,
                }
            }).then(dt=>{
                let res = dt.data;
                if(res.success){  
                    this.target_result = res.data.target_result;
                    // this.verify_result = res.data.verify_result;
                    for(var i=0; i< res.data.verify_result.length;i++){
                        var item =  res.data.verify_result[i];
                        if(this.isJSON(item)){ //报文
                            var back = JSON.parse(item);  
                            this.verify_result.push({'str':'success：'+back.success});
                            for( var key in back ){ 
                                if(key != 'request' && key != 'response' && key != 'success'){ 
                                    this.verify_result.push({'str':key+'：url：'+back[key].url})
                                }  
                            }
                            this.verify_result.push({'back':[back.request[0],back.response[0]]})
                        }else{
                            this.verify_result.push({'str':item})
                        }

                        
                    }

                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
                
            }).catch(err=>{})
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
        height: 100%;
        padding: 24px;
        background: #fff;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12); 
        box-sizing: border-box;
        margin-bottom: 15px;
        border-radius: 4px;
    }
</style>