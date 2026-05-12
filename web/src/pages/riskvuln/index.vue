<template>
    <div>
        <div class="main-title  ">   
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label class="taskSceneBtn" >漏洞风险</label>
            
        </div>  
         <div  class="box-table-list context_box_bg">
            <div class="search-box"> 
                <div class="operationbutton"> 
                     <el-button type="primary"  @click="handleStatus()" :disabled="ids.length<=0" >更改状态</el-button> 
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" style="padding-left:8px"
                        trigger="click" :visible-arrow="false" v-model="alldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                        <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消</el-button>
                        <el-button size="mini" type="primary" @click="btnMultiDeleteTask">确定</el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference" :disabled="!ids.length">删除</el-button>
                    </el-popover>
                </div> 
                <div class="serach-condition">
                    <div class="search-text">
                        <el-input placeholder="漏洞地址" @keydown.enter.native="handlesearch" v-model="param.location" class="input-with-select" size="small"
                         >
                        </el-input> 
                    </div>
                    <div class="search-text">
                        <el-input placeholder="漏洞名称" @keydown.enter.native="handlesearch" v-model="param.search" class="input-with-select" size="small"
                         >
                        </el-input>
                        <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                    </div>
                </div>
            </div>
            <el-table :data="tableData" style="width: 100%" class="myTable" @selection-change="handleSelectionChange"
                  @cell-mouse-enter="mouseenter" height="calc(100% - 102px)"
                        @cell-mouse-leave="mouseleave" >
                <el-table-column type="selection" width="55" align="center" />   
                <el-table-column   label="漏洞地址" min-width="20%">
                    <template slot-scope="scope" >
                        <el-tooltip
                            effect="dark"
                            :content="scope.row.location"
                            placement="top"
                            :open-delay="500"
                            popper-class="copyable-tooltip" >
                            <span class="cell-text">{{ scope.row.location }}</span>
                        </el-tooltip> 
                    </template>
                </el-table-column> 
                <el-table-column prop="name" label="漏洞名称" min-width="20%" show-overflow-tooltip>
                </el-table-column>
                <el-table-column prop="typeName" label="漏洞类型" min-width="15%" show-overflow-tooltip>  
                    <template slot="header" slot-scope="scope">
                        <span class="cursorPointer" @click="clickButton('漏洞类型')" :class="param.type ? 'active' : ''">漏洞类型<i class="iconfont iconshaixuan"></i></span>
                        <el-select popper-class="thSelect" v-model="param.type" clearable placeholder="漏洞类型" size="small" ref="vulTypeSelect" @change="handlesearch">
                            <el-option v-for="(item, i) in vulTypelist" :key="i" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                    <template #default="scope" >
                        <span>{{ scope.row.typeName }}</span>
                    </template>
                </el-table-column>
                <el-table-column   prop="riskName" label="漏洞等级" min-width="10%" show-overflow-tooltip>  
                    <template slot="header" slot-scope="scope">
                        <span class="cursorPointer" @click="clickButton('漏洞等级')" :class="param.risk ? 'active' : ''">漏洞等级<i class="iconfont iconshaixuan"></i></span>
                        <el-select popper-class="thSelect" v-model="param.risk" clearable placeholder="漏洞等级" size="small" ref="vulRiskSelect" @change="handlesearch">
                            <el-option v-for="(item, i) in vulRisklist" :key="i" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                    <template #default="scope"> 
                        <div class="risk_style">
                            <el-tag type="danger"  class="risk risk4" size="mini" v-if="scope.row.riskLevel==0"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                            <el-tag type="danger" class="risk risk0" size="mini" v-if="scope.row.riskLevel==1"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                            <el-tag type="warning" class="risk risk1" size="mini" v-if="scope.row.riskLevel==2"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                            <el-tag type="primary" class="risk risk2"  size="mini"  v-if="scope.row.riskLevel==3"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                            <el-tag type="success" class="risk risk3" size="mini"  v-if="scope.row.riskLevel==4"><i></i>{{scope.row.riskLevelStr}}</el-tag> 
                            <el-tag type="info"  class="risk risk4" size="mini"  v-if="scope.row.riskLevel==5"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                        </div>  
                    </template>
                </el-table-column> 
                
                <el-table-column   prop="verifyTypeName" label="验证方式" min-width="11%" show-overflow-tooltip> 
                    <template slot="header" slot-scope="scope">
                        <span class="cursorPointer" @click="clickButton('验证方式')" :class="param.verifyType ? 'active' : ''">验证方式<i class="iconfont iconshaixuan"></i></span>
                        <el-select popper-class="thSelect" v-model="param.verifyType" clearable placeholder="验证方式" size="small" ref="verifyTypeSelect" @change="handlesearch">
                            <el-option v-for="(item, i) in verifyTypelist" :key="i" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                </el-table-column>
                <el-table-column   prop="statusName" label="漏洞状态" min-width="10%" show-overflow-tooltip>    
                    <template slot="header" slot-scope="scope">
                        <span class="cursorPointer" @click="clickButton('漏洞状态')" :class="param.status ? 'active' : ''">漏洞状态<i class="iconfont iconshaixuan"></i></span>
                        <el-select popper-class="thSelect" v-model="param.status" clearable placeholder="漏洞状态" size="small" ref="statusSelect" @change="handlesearch">
                            <el-option v-for="(item, i) in statusList" :key="i" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                </el-table-column> 
                <!-- <el-table-column   prop="fixStatusStr" label="修复状态" show-overflow-tooltip>    
                </el-table-column>  -->
                <el-table-column prop="findTime" label="发现时间"  width="200">
                    <template #default="scope"> 
                         <div v-if="showOperateButton && rowId == scope.row.id">
                            <!-- <el-button  link type="primary" @click="btnvulninfo(scope.row)"  >详情</el-button>   -->
                            <el-link
                                :underline="false"
                                class="link_primary"
                                style="vertical-align: initial"
                                @click="btnvulninfo(scope.row)">详情</el-link>
                            <!-- <el-button size="mini"  link type="primary"  @click="handleStatus(scope.row)">更改状态</el-button> -->
                            <el-link
                                :underline="false"
                                class="link_primary"
                                style="vertical-align: initial"
                                @click="handleStatus(scope.row)">更改状态</el-link>
                            <!-- <el-button  link type="primary"  @click="delData(scope.row)">删除</el-button> -->
                            <el-popover placement="bottom" width="170" :visible-arrow="false" :ref="`popover_id-${scope.row.id}`"
                                popper-class="delButton_popper">
                                <p class="delText">
                                <i class="el-icon-warning"></i>确定删除吗？
                                </p>
                                <div style="text-align: right; margin: 0">
                                <el-button size="mini" class="delCancel" @click="btnCancelDel(scope)">取消</el-button>
                                <el-button size="mini" type="primary" @click="delData(scope)">确定</el-button>
                                </div>
                                <!-- <span slot="reference">删除</span> -->
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference">删除
                                </el-link>
                            </el-popover>
                         </div> 
                         <div v-else>
                            {{ scope.row.findTime }}
                         </div>
                    </template>
                </el-table-column> 
               
            </el-table> 

            <el-pagination  :page-size="param.size" background layout="total, prev, pager, next, sizes, jumper"
                :total="total" :current-page="param.page" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>
        </div>

         <el-dialog :visible.sync="dialogStatusVisible" title="更改状态" width="500px" :close-on-click-modal="false" custom-class="change-status-dialog">
            <el-form label-width="80px" style="padding: 10px 0;">
                <el-form-item label="选择状态" style="margin-bottom: 0;">
                    <el-select v-model="selectedStatus" size="small" style="width: 100%;">
                        <el-option v-for="(item,i) in headerStatusSelect" :key="i" :label="item.label" :value="item.value"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button size="small" @click="dialogStatusVisible = false">关闭</el-button>
                    <el-button size="small"  type="primary" @click="handleSave">保存</el-button>
                </div>
            </template>
        </el-dialog>
        <!-- 详情 -->
        <el-dialog
            :title="vulInfo.vul_name"
            :visible.sync="vuldialogVisible"
            width="1184px"
            class="buginfobox" 
            :close-on-click-modal="false" 
            :show-close="true" >
             <div class="dialog_b_btn">  
				<el-button size="small" @click="cancalvuldialogVisible">关闭</el-button>
            </div>
            <div class="buginfo_box" > 
                <div class="bugbasicinfo" style="background-color: #fff;padding: 24px;border: 1px solid #E8E8F5; "> 
                    <el-table
                        :data="infotabledata"
                        size='small'
                        style="width: 100%">
                        <el-table-column
                            prop="vul_typeName"
                            label="漏洞类型" > 
                        </el-table-column>
                         <el-table-column  prop="vul_statusName" label="漏洞状态"> 
                         </el-table-column>
                        <el-table-column
                            prop="vul_risk"
                            label="漏洞风险" >
                            <template slot-scope="scope">
                                <span   > 
                                   <span :class="[ 
                                   {'riskstyle risk_hight': scope.row.vul_risk_value == 1} ,
                                   {'riskstyle risk_middle': scope.row.vul_risk_value == 2},
                                   {'riskstyle risk_low':scope.row.vul_risk_value == 3 },
                                   {'riskstyle risk_nofind':scope.row.vul_risk_value == 4 }]"
                                        ><i></i>{{scope.row.vul_risk}}</span> 
                                </span> 
                            </template>
                        </el-table-column>
                        <el-table-column  prop="vul_targetUrl" label="所属资产"> 
                         </el-table-column> 
                    </el-table>
                </div> 
                <div class="bugotherinfo">
                    <div class="part_title"> 漏洞描述</div>
                    <div class="content" v-if="vulInfo.vul_description" > {{vulInfo.vul_description}}  </div> 
                </div>
                <div class="bugotherinfo" >
                    <div class="part_title">修复建议</div>
                    <div class="content" v-if="vulInfo.fixSuggest" >{{vulInfo.fixSuggest}}</div> 
                </div>
                <div class="bugotherinfo">
                    <div class="part_title">漏洞测试报文</div>
                    <div class="bugbasicinfo">
                        <el-row style="margin-top:10px" :gutter="20" >
                                <el-col :span="12">
                                    <div class="message requestpack">
                                        <label class="title_bg title_bg1">请求报文</label>
                                        <div> 
                                            <div class="packheight " style="height:100%">
                                                <div class="packinput" style="height:100%">
                                                    <div  style=" background: #fff;" >
                                                        <pre  >{{ verMsg.request }}</pre>
                                                    </div> 
                                                    
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </el-col>
                            <el-col :span="12">
                                <div class="message">
                                    <label class="title_bg title_bg2">响应报文</label>
                                    <div style="    background: #fff;"> 
                                        <div style="    background: #fff;">
                                            <pre >{{ verMsg.response }}</pre>
                                        </div>
                                    </div>
                                </div>
                            </el-col>
                        </el-row>
                    </div>
                </div>
                <div class="bugotherinfo">   
                     <div class="part_title">漏洞生命周期</div>
                    <el-timeline>
                        <el-timeline-item   placement="top" v-for="(item,i) in Lifecyclelist" :key="i" :timestamp="item.time">
                           {{ item.content }}
                        </el-timeline-item>
                        <!-- <el-timeline-item timestamp="2018/4/3" placement="top">
                            11
                        </el-timeline-item> -->
                    </el-timeline>
                </div>
            </div>
        </el-dialog>
    </div>
</template>
<style lang="less">
.thSelect {
    min-width: 150px;
}
</style>
<style lang="less" scoped>
    /deep/ .myTable{
    thead {
        .cursorPointer{
            cursor: pointer;
            &.active{
                color:#4C7AE3;
                i{
                    color: #4C7AE3;
                }
            }
        }
        .cell{
            line-height: 15px; 
            // height: 26px;
            .caret-wrapper {
                top: -9px;
            } 
            >span {
                position: absolute;
            }
        }
        .iconfont{
            color:rgba(72,72,102,0.32);
            margin-left:5px;
        }
        .el-select{
            height: 0;
            visibility: hidden;
            .el-input, .el-input__inner{
                height: 0!important;
            }
        }
    }
}
    .el-row{
        height:100%;
    }
/* 允许 tooltip 内容被选中和复制 */
.copyable-tooltip {
  user-select: text !important;
  pointer-events: auto !important;
}
.cell-text {
  display: inline-block;
  width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.serach-condition > div {
    display: inline-block;
    margin-right: 8px;
    vertical-align: middle;
}
.bugotherinfo{
    margin-top: 24px;
    .part_title{
        margin-bottom: 8px;
        font-size: 14px;
        border-left: 3px solid #4c7ae3;
        padding-left: 8px;
        font-weight: 500;
        width: 113px;
        display: inline-block;
        line-height: 16px;
        box-sizing: border-box;
        color: rgba(72, 72, 102, 0.87) !important;
    }
    .content{
        background:rgba(255,255,255,1);
        border-radius:2px;
        border:1px solid rgba(232,232,245,1);
        padding: 12px 16px;
        color: rgba(72, 72, 102, 0.64);
        font-size: 13px;
    }
}
.box-table-list{
     padding: 24px;
    background: #fff;
    height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
.buginfo_box{
    background: #F7F7FB; 
    height:100%;
    overflow: auto;
    padding: 24px ;
    box-sizing: border-box;
}
.bugbasicinfo{
    /* padding: 24px;
    background: #fff; */
    .el-form-item{
        margin-bottom: 14px;
        &.validateLog{
            position: relative;
            .showWaiting{
                position: absolute;
                top: 45px;
                left:10px;
                z-index: 2;
            }
        }
    }
    /deep/ .el-form-item__label{
        text-align: left;
        border-left: 3px solid #4C7AE3;
        line-height: 16px;
        padding-left: 8px;
        margin-top: 12px;
        position: relative;
        &:before{
            display: none;
        }
    }
    .flexBet{
        position: relative;
        span:nth-child(2){
            position: absolute;
            right: 6px;
            top: 3px;
            color: #F56C6C;
        }
    }
}
 .requestpack>div {
     background: #fff !important;
     padding: 0 !important;
 }
.message .title_bg {
     margin-bottom: 8px;
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
     height: 253px;
     overflow-y: auto;
 }
 .packbtn {
     height: 65px;
     box-sizing: border-box;
     padding: 16px;
     text-align: left;
     padding-left: 0;
 }
.message>div {
    // margin-bottom: 24px; 
    background: #F7F7FB;
    border-radius: 4px;
    border: 1px solid #E8E8F5;
    padding: 16px;
    box-sizing: border-box;
}
/deep/ .change-status-dialog {
    .el-dialog__body {
        padding: 10px 20px !important;
    }
}
</style>
<script>
import risk from '@/api/risk.js' 
export default {
    name:'riskindex',
    comments:{ 
    },
    data(){
        return{
            sendVal: false,
            tableData:[], 
            total:0,
            param:{
                search:'',
                location:'',
                page:1,
                size:10, 
                type: '',
                risk: '',
                verifyType: '',
                status: '',
            },
            vulTypelist: [],
            vulRisklist: [],
            verifyTypelist: [],
            statusList: [],
            dialogStatusVisible:false,
            selectedStatus:'',
            headerStatusSelect:[],
            ids:[],
            showOperateButton: false,
            rowId: '',
            vulInfo:{},
            vuldialogVisible:false,
            infotabledata:[],
            Lifecyclelist:[],
            verMsg:{},
            alldelvisible:false,
            operateIds: '',
        }
    },
    created(){
         this.$store.state.activefirstMenu = "/riskvuln";
    },
    mounted(){
        this.getData();
        this.getEnums();
    },
    methods:{
        async getEnums() {
            try {
                const res = await risk.riskvulnums();
                if (res.code == 200 && res.data) {
                    this.vulTypelist = res.data.vulType || [];
                    this.vulRisklist = res.data.vulRisk || [];
                    this.verifyTypelist = res.data.vulVerifyType || [];

                    if (res.data.vulRiskStatus) {
                        this.statusList = res.data.vulRiskStatus;
                        this.headerStatusSelect = res.data.vulRiskStatus;
                    }
                }
            } catch (e) {
                console.error(e);
            }
        },
        clickButton(type) {
            switch (type) {
                case '漏洞类型':
                    this.$refs.vulTypeSelect.toggleMenu();
                    break;
                case '漏洞等级':
                    this.$refs.vulRiskSelect.toggleMenu();
                    break;
                case '验证方式':
                    this.$refs.verifyTypeSelect.toggleMenu();
                    break;
                case '漏洞状态':
                    this.$refs.statusSelect.toggleMenu();
                    break;
            }
        },
        getData(){
            risk.riskvullist(this.param).then(res=>{
                if(res.code == 200){
                    this.tableData = res.data.list;
                    this.total = res.data.total;
                }
            })
        },
        // 更改状态
        handleStatus(row){  
            if(row){
                this.operateIds = row.id + '';
                this.selectedStatus = row.status;
            }else{
                this.operateIds = this.ids.join(',');
                this.selectedStatus = '';
            }
            this.dialogStatusVisible = true;
        },
        //保存状态
        handleSave(){
             if(!this.selectedStatus){
                 this.$message.warning('请选择状态');
                 return;
             }
             risk.riskvulnStatus({
                'ids':this.operateIds,
                status:this.selectedStatus
            }).then((res)=>{
                if(res.code == 200){
                    this.$message.success('更新状态成功');
                    this.dialogStatusVisible = false; 
                    this.selectedStatus = '';
                    this.getData();
                }
                else{
                    this.$message.error(res.msg);
                }
            })
        },
        btnCancelDel(scope){
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose(); 
        },
        copyText(text) {
            // 创建一个临时 input 元素
            const textarea = document.createElement('textarea');
            textarea.value = text;
            document.body.appendChild(textarea);
            textarea.select();
            try {
                const success = document.execCommand('copy');
                if (success) {
                this.$message.success('复制成功！');
                } else {
                this.$message.error('复制失败，请手动复制。');
                }
            } catch (err) {
                this.$message.error('复制出错，请手动复制。');
            }
            document.body.removeChild(textarea);
        },
        handlesearch(){
            this.param.page=1;
            this.getData();
        },
        handleReset(){
            this.param.page=1;
            this.param.size=10;
            this.param.search='';
            this.param.location = '';
            this.param.type = '';
            this.param.risk = '';
            this.param.verifyType = '';
            this.param.status = '';
            this.getData();
        },
        handleSelectionChange(selection){
            this.ids = selection.map(item => item.id) 
        },
        handleSizeChange(t){
            this.param.page = 1;
            this.param.size = t;
            this.getData();
        },
        currentchange(t){
            this.param.page = t;
            // this.param.size = 10;
            this.getData();
        },
         mouseenter (row, column, cell, event) {
            this.showOperateButton = true
            let _id = row.id
            this.rowId = _id//赋值行id，便于页面判断
            },
        mouseleave (row, colum, cell, event) {
            let t = this.$refs['popover_id-' + row.id] && this.$refs['popover_id-' + row.id].showPopper
            let t2 = this.$refs['popover-' + row.id] && this.$refs['popover-' + row.id].showPopper

            if (!t && !t2) {
                this.showOperateButton = false
                this.rowId = ""
            }

        },
        async btnvulninfo(row){
             this.infotabledata=[];
            const res = await risk.vulinfo({
                id:row.id,
            })
            if(res.code == 200){ 
                this.vuldialogVisible = true;
                let _info = res.data;
                this.vulInfo.vul_name = _info.name;
                var vuljosn= { 
                    vul_type:_info.vul_type, 
                    vul_typeName:_info.typeName,
                    vul_risk:_info.riskName,
                    vul_risk_value:_info.risk, 
                    vul_statusName:_info.statusName, 
                    vul_targetUrl:_info.targetUrl  ,
                    vul_fixStatusStr : _info.fixStatusStr  ,
                }; 
                this.infotabledata.push(vuljosn);
                this.vulInfo.vul_description = _info.description;
                this.vulInfo.fixSuggest =  _info.fixSuggest;
                this.vulInfo.vulLifecycle = _info.vulLifecycle;
                // this.Lifecyclelist = _info.vulLifecycle.lifecycle || [];

                const allLifecycles = _info.vulLifecycle.flatMap(item => item.lifecycle);
                // （从晚到早）
                allLifecycles.sort((a, b) => {
                    return new Date(b.time) - new Date(a.time);
                });
                 this.Lifecyclelist =allLifecycles;
                this.verMsg = _info.verMsg && _info.verMsg[0];
            }

        },
        async btnMultiDeleteTask(){
            if (this.ids.length == 0) return; 
            let _ids = this.ids; 
             const res = await	risk.vuldelbyid({
                ids: _ids.join(",")
              }); 
              if (res.code == 200) {
                this.$message({
                  message: "删除漏洞成功",
                  type: "success"
                }); 
                this.param.page = 1; 
                this.alldelvisible = false;
                this.getData();
              } else {
                this.$message({
                  message: res.msg,
                  type: "error"
                });
              }
        },
        async delData(scope){
             const res = await	risk.vuldelbyid({
                ids: scope.row.id+'',
              }); 
              if (res.code == 200) {
                this.$message({
                  message: "删除漏洞成功",
                  type: "success"
                }); 
                 scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
                this.getData();
              } else {
                this.$message({
                  message: res.msg,
                  type: "error"
                });
              }
        }, 
        cancalvuldialogVisible(){
            this.vuldialogVisible = false;
          
        },
    }
}
</script>