<template>
    <!-- 漏洞信息 -->
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
                        <el-input placeholder="搜索漏洞名称与漏洞位置" @keydown.enter.native="handlesearchbug" v-model="formDatabug.keyword" class="input-with-select"
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
                <el-table-column prop="ip" label="测试目标" :show-overflow-tooltip="true" width="200px"> 
                </el-table-column>

                <el-table-column prop="port" :show-overflow-tooltip="true" label="端口" width="100px">
                </el-table-column>

                <el-table-column prop="name" label="漏洞名称" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <el-link @click="btnbuginfo(scope.row)">{{scope.row.name}}</el-link>
                    </template>
                </el-table-column>

                <el-table-column prop="cve" :show-overflow-tooltip="true" width="200px" label="cve">
                </el-table-column>

                <el-table-column prop="cwe" label="cwe类型" :show-overflow-tooltip="true" width="200px">
                    <template slot-scope="scope" slot="header">
                        <div style="position: relative;height:10px;">
                            <span class="cursorPointer">cwe类型</span>
                        </div>
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formDatabug.name"
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

                <el-table-column prop="createTime" label="发现时间"> 
                   
                    <template slot-scope="scope" slot="header">
                        <div style="position: relative;height:10px;">
                            <span class="cursorPointer" >发现时间</span>
                        </div>
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formDatabug.port"
                            size="small" ref="port" @change="handlesearchbug">
                            <el-option v-for="(item,i) in statusTypelist" :key="i" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                    
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                            <!-- :disabled="scope.rowisSnapshot != 1"  -->
                            
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
                            <span>{{ scope.row.createTime }}</span>
                        </div>
                    </template>
                </el-table-column>

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
        
      
    </div>
</template>

<script>
import vulnmsginfo from './Vulnmsginfo.vue'
import VulScanDatabaseList_detail from './VulScanDatabaseList_detail.vue'

import { task } from '@/api/task.js';
import { vulscan } from '@/api/vulscan.js';
import { vulnerability } from "@/api/tool.js";
export default {
    name: "vulnmsg",
    components:{
        vulnmsginfo,
        VulScanDatabaseList_detail
    },
    props:{
        target_id:{}, 
        task_id:{},
    },
    data: () => ({  
        usedVal:false,
        vulninfo:{},
        snapBase:'',
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
            const res = await vulscan.vulList({  
                page: this.formDatabug.page,
                size: this.pageSize,
                taskId: this.task_id,
                risk: this.formDatabug.risk_level,
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
        padding: 24px;
        background: #fff;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12); 
        box-sizing: border-box;
        margin-bottom: 15px;
        border-radius: 4px;
    }
</style>