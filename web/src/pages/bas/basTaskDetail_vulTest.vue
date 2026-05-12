<template>
    <div>
        <div class="statistics_box">
            <div class="attack_stat">
                <label for="" class="part_title">攻击统计</label>
                <div>
                    <label class="lbname"> <i></i>未拦截数量 </label>
                    <span class="lbvalue">{{status[0]}}</span>
                </div>
                <div>
                    <label class="lbname" > <i></i>已拦截数量 </label>
                    <span class="lbvalue">{{status[1]}}</span>
                </div> 
            </div>
            <div class="vuln_stat">
                <label for="" class="part_title">影响级别统计</label>
                <div>
                    <label class="lbname"> <i class="risk_deadly"></i>高危 </label>
                    <span class="lbvalue">{{riskLevel[0]}}</span>
                </div>
                <div>
                    <label class="lbname" > <i class="risk_hight"></i>中危 </label>
                    <span class="lbvalue">{{riskLevel[1]}}</span>
                </div>
                <div>
                    <label class="lbname" > <i class="risk_middle"></i>低危 </label>
                    <span class="lbvalue">{{riskLevel[2]}}</span>
                </div>
                <div>
                    <label class="lbname" > <i class="risk_low"></i>安全 </label>
                    <span class="lbvalue">{{riskLevel[3]}}</span>
                </div>
                
            </div>
        </div>
        <div class="list_box">
            <div class="search-box"> 
                <div  class="operationbutton" >   
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" style="padding-left:8px"
                        trigger="click" :visible-arrow="false" v-model="alldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消</el-button>
                            <el-button size="mini" type="primary" @click="AllDel">确定</el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除</el-button>
                    </el-popover> 
				</div>
                <div class="serach-condition" > 
					<div class="search-text">
						<el-input placeholder="搜索关键字" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select"  size="small" clearable > </el-input>
						<el-button type="primary"  size="small" @click="handlesearch">搜索</el-button> 
					</div> 
					<div >
						<el-button type="primary"  size="small" @click="handleReset">重置</el-button> 
					</div>   
				</div>  
            </div> 
            <el-table
				:data="tableData"  style="width: 100%"  class="myTable"   @selection-change="handleSelectionChange" 
				  @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" >
                  <el-table-column  
                    width="55" 
                    type="selection"  >
                </el-table-column>  
			    <el-table-column   
					prop="addr"
					label="攻击目标" width="180px">    
				</el-table-column>    
				<el-table-column prop="ruleName" label="剧本名称"  :show-overflow-tooltip="true" > 
                </el-table-column>  
                <el-table-column prop="attackModeName" label="攻击方式"  width="200px"> 
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('攻击方式')"
                            :class="(formData.attackType !== '' ) ? 'active' : ''">攻击方式<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.attackType" clearable
                        size="small" ref="attackMode" @change="handlesearch">
                        <el-option v-for="(item, index) in attackModelist" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                         {{ scope.row.attackModeName }}
                    </template>
                </el-table-column>    
                <el-table-column prop="riskLevelName" label="影响级别"  width="140px"> 
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('影响级别')"
                            :class="(formData.riskLevel !== '' && formData.riskLevel !== 0) ? 'active' : ''">影响级别<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.riskLevel" clearable
                        size="small" ref="class" @change="handlesearch">
                        <el-option v-for="(item, index) in riskLevellist" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                         {{ scope.row.riskLevelName }}
                    </template>
                </el-table-column>  
                <el-table-column prop="attackStageName" label="攻击阶段" width="140px" > 
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('攻击阶段')"
                            :class="(formData.attackStage !== '' && formData.attackStage !== 0) ? 'active' : ''">攻击阶段<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.attackStage" clearable
                        size="small" ref="attackStage" @change="handlesearch">
                        <el-option v-for="(item, index) in attackTypelist" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                         {{ scope.row.attackStageName }}
                    </template>
                </el-table-column>   
				<el-table-column prop="status" label="攻击结果"  width="140px"> 
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('攻击结果')"
                            :class="(formData.status !== '' ) ? 'active' : ''">攻击方式<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.status" clearable
                        size="small" ref="status" @change="handlesearch">
                        <el-option v-for="(item, index) in statuslist" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template> 
					<template slot-scope="scope" >  
                        <div v-if="showEditFileNameButton && rowId == scope.row.id"> 
                            <el-link  :underline="false" class="link_primary"  @click="handleInfo(scope.row)">详情</el-link>   
                            <el-popover 
                                placement="bottom"
                                width="170"   
                                :ref="`popover_id-${scope.row.id}`"
                                popper-class="delButton_popper" >
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                    <el-button size="mini" type="primary" @click="handleDel(scope)">确定</el-button>
                                </div> 
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0"  slot="reference" >删除</el-link>  
                            </el-popover>
                        </div>
                        <div v-else > 
                             <span :class="[  
                                { 'tag_status tag_warning': scope.row.status ==1 },
                                { 'tag_status tag_primary': scope.row.status == 2 },
                                { 'tag_status tag_success': scope.row.status ==3 }]"><i></i>{{ scope.row.statusName }}</span>
                        </div> 
                    </template>
                   
				</el-table-column>
			</el-table> 
			<el-pagination
				:page-size="pageSize" 
				background
				layout=" total,  prev, pager, next,sizes, jumper"
				:total="totalpage"
				:current-page="currentpage"
				@current-change = "handlecurrentchange"
				@size-change="handleSizeChange" >
			</el-pagination>
        </div>

         <!-- 详情 -->
         <el-dialog
            title="剧本详情"
            :visible.sync="vuldialogVisible"
            width="1184px"
            class="buginfobox" 
            :close-on-click-modal="false" 
            :show-close="false" >
            <div class="dialog_b_btn">  
                <!-- <el-button size="small" @click="btnEdit">{{updateBugtxt}}</el-button>
                <el-button size="small" @click="saveEdit" v-if="is_Update">保存</el-button> -->
				<el-button size="small" @click="cancalvuldialogVisible">关闭</el-button>
            </div>
            <div class="buginfo_box" > 
                <div class="bugbasicinfo"> 
                    <el-table
                        :data="vulbasicinfo"
                        size='small'
                        style="width: 100%">
                        <el-table-column
                            prop="nameZh"
                            label="剧本名称" > 
                        </el-table-column>
                        <el-table-column
                            prop="affectTarget"
                            label="剧本目标"> 
                            <template slot-scope="scope">
                                <span   v-if="!is_Update"  >{{scope.row.affectTarget}}  
                                </span> 
                                <el-input  v-if="is_Update"  v-model="editInfo.affectTarget"  size="mini"  ></el-input>
                            </template>
                        </el-table-column>
                        <!-- <el-table-column
                            prop="vul_risk"
                            label="攻击类型" >
                            <template slot-scope="scope">
                                <span v-if="!is_Update" > 
                                   <span >{{scope.row.classType}}</span> 
                                </span>
                                <el-select v-model="editInfo.classType"   size="mini" v-if="is_Update" >  
                                    <el-option
                                        v-for="(item,index) in ruleenum"
                                        :key="index"
                                        :label="item.label"
                                        :value="item.value">
                                        </el-option>
                                </el-select> 
                            </template>
                        </el-table-column>
                        <el-table-column
                            prop="attackStage"
                            label="攻击阶段">
                            <template slot-scope="scope">
                                <span   v-if="!is_Update"  >{{scope.row.attackStageEnum}}  
                                </span> 
                                <el-select v-model="editInfo.attackStage" size="mini" v-if="is_Update">
                                    <el-option
                                        v-for="(item,index) in attackTypelist"
                                        :key="index"
                                        :label="item.label"
                                        :value="item.value">
                                    </el-option>
                                </el-select>
                            </template>
                        </el-table-column>
                         -->
                        <el-table-column
                            prop="riskLevelEnum"
                            label="影响级别"> 
                            <template slot-scope="scope">
                                <span   v-if="!is_Update"  >{{scope.row.riskLevelEnum}}  
                                </span>
                                <el-select v-model="editInfo.riskLevel" size="mini" v-if="is_Update">
                                    <el-option
                                        v-for="(item,index) in riskLevellist"
                                        :key="index"
                                        :label="item.label"
                                        :value="item.value">
                                    </el-option>
                                </el-select>
                            </template>
                        </el-table-column>
                        <el-table-column
                            prop="affectScore"
                            label="影响评分">
                            <template slot-scope="scope">
                                <span   v-if="!is_Update"  >{{scope.row.affectScore}}  
                                </span> 
                                <el-input  v-if="is_Update"  v-model="editInfo.affectScore"  size="mini"  ></el-input>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
             
                <div class="bugotherinfo">
                    <div class="part_title">剧本描述</div>
                    <div class="content"  > {{scriptLibraryInfo.descriptionZh}}  </div>  
                </div>
                <div class="bugotherinfo">
                    <div class="part_title">CONTENT</div>
                    <div class="content"  > {{scriptLibraryInfo.content}}  </div>  
                </div>
                <div class="bugotherinfo" >
                    <div class="part_title">关联攻击方式</div>
                    <div class="content" v-if="!is_Update" v-html="scriptLibraryInfo.relationAttackMethod"></div>
                    <div :class="[{'isError':showerr}]">
                        <el-input class="textarea" type="textarea" 
                        v-model="editInfo.relationAttackMethod" @change="change_vul_description"  
                        size="mini" :row="10" v-if="is_Update" 
                        resize="none"></el-input>
                    </div>  
                </div>
                <div class="bugotherinfo" >
                    <div class="part_title">修复建议</div>
                    <div class="content" v-if="!is_Update">{{scriptLibraryInfo.fixSuggest}}</div>
                    <el-input class="textarea" type="textarea" v-model="editInfo.fixSuggest"  size="mini"  
                    resize="none" :row="3" v-if="is_Update" ></el-input>
                   
                </div> 
                <div class="bugotherinfo">
                    <div class="part_title">参考链接</div>
                    <div class="content" v-if="!is_Update">{{scriptLibraryInfo.refUrl}}</div>
                    <el-input class="textarea" type="textarea" v-model="editInfo.refUrl" resize="none" size="mini"  :row="3" v-if="is_Update" ></el-input>
                </div>
                
            </div>
        </el-dialog>


    </div>
</template>
<style lang="less" scoped>
.part_title {
     font-size: 14px;
     margin-bottom: 16px;
     font-weight: 800;
     border-left: 3px solid #4C7AE3;
     padding-left: 10px;
     height: 14px;
     line-height: 14px;
     color: rgba(72, 72, 102, 0.89);
 }
.list_box{
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
.bugbasicinfo {
    padding: 24px;
    background: #fff;
    border: 1px solid #e8e8f5;
}
.statistics_box{
    padding: 24px;
    background: #fff;
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    margin-bottom: 12px;
    .attack_stat{
        margin-bottom: 24px;
    }
    .attack_stat,
    .vuln_stat{
        label{
            display: inline-block;
            font-size: 14px;
            margin-right: 32px;
        }
        >div{
            display: inline-block;
            margin-right: 24px;
        }
        .lbname{
            i{
                display: inline-block;
                width: 6px;
                height: 6px;
                border-radius: 50%;
                
                margin-right: 8px;
                vertical-align: middle;
            }
            color: #484866;
            font-size: 13px; 
            margin-right: 16px;
        }
        .lbvalue{
            display:inline-block;
            font-size: 13px;
        }
    } 
}
.risk_deadly {
    background: #F87D7D;
}
.risk_hight {
    background: #FDC665;
}
.risk_middle {
    background: #4c7ae3;
}
.risk_low {
    background: #65c680;
}
/deep/ thead {
    .cursorPointer {
        cursor: pointer; 
        position: absolute;
        // top: 6px;
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
    .bugotherinfo{
        margin-top: 32px;
        .part_title{
            margin-bottom: 8px;
        }
        .content{
            background:rgba(255,255,255,1);
            border-radius:2px;
            border:1px solid rgba(232,232,245,1);
            padding: 12px 16px;
            color:rgba(72,72,102,0.64);
            font-size: 13px;
        }
    }
</style>
<script>
import bas from '@/api/bas.js'
export default {
    data(){
        return{
            alldelvisible:false,
            multipleSelection:[],
            tableData:[],
            formData:{
                search:'',
                riskLevel:'',
                page:1,
                attackStage:'',
                attackType:'',
                status:'',
            },
            pageSize:10,
            showEditFileNameButton:false,
            rowId:'',
            currentpage:1,
            totalpage:0,
            status:[],
            riskLevel:[],
            riskLevellist:[],
            attackTypelist:[],
            ruleenum:[],
            updateBugtxt:'编辑',
            is_Update:false,
            scriptLibraryInfo:{
                attackMode:'',
                descriptionZh:'',
                refUrl:'',
                fixSuggest:'',
                relationAttackMethod:'',
                content:'',
            },
            editInfo:{
                riskLevel:'',
                riskLevelEnum:'',
                descriptionZh:'',
                refUrl:'',
                fixSuggest:'',
                relationAttackMethod:'', 
                attackMode:'',
                affectTarget:'',
                attackStage:'', 
                affectScore:'',
            },
            vuldialogVisible:false,
            showerr:false,
            vulbasicinfo:[],
            isShow2:false,
            showhide2:false,
            attackModelist:[],
            statuslist:[],
        }
    },
    props:{ 
        task_id:{},
    },
    created(){

    },
    mounted(){
        this.basruleenum();
        this.getbasvulstat();
        this.getData();
    },
    methods:{
        async basruleenum(){
            const res = await bas.basruleenum();
           
            if(res.code == 200){
                this.ruleenum = res.data.class;
                this.riskLevellist = res.data.riskLevel;
                this.riskLevellist.unshift(
                    { label: "全部", value: 0 }
                ) 
                 
                this.attackTypelist = res.data.attackStage;
                this.attackTypelist.unshift(
                    { label: "全部", value: 0 }
                ) 
                this.attackModelist = res.data.attackType;
                this.attackModelist.unshift(
                    { label: "全部", value: '' }
                ) 
                this.statuslist =  res.data.status;
                this.statuslist.unshift(
                    { label: "全部", value: '' }
                ) 

            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }  
             
        },
        async getbasvulstat(){
            const res = await bas.basvulstat({
                basTaskId:Number(this.task_id)
            });
            if(res.code == 200){
                this.status = res.data.status;
                this.riskLevel = res.data.riskLevel;
            }
        },
        async getData(){
            const res = await bas.basvullist({
                page:this.formData.page,
                size:this.pageSize,
                search:this.formData.search,
                basTaskId:this.task_id,
                riskLevel:this.formData.riskLevel,
                attackStage:this.formData.attackStage,
                attackMode:this.formData.attackType,
                status:this.formData.status,
            });
            if(res.code == 200){
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }else{ 
                this.$message({
                    message: res.msg,
                    type: "error",
                }); 
            }
        }, 
        async handleDel(scope){
            let params = {
                basVulIds:scope.row.id+'',
                basTaskId:this.task_id
            }
            const res = await bas.basvuldel(params)
            if(res.code == 200){ 
                this.$message({
                    message:'删除成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                this.getData();
                this.getbasvulstat();
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        async AllDel(){
            if(this.multipleSelection.length == 0) return;
    		var ids = [];
    		for (var i = 0; i < this.multipleSelection.length; i++) {
    			ids.push(this.multipleSelection[i].id);
    		}
            let params = {
                basVulIds:ids.join(','),
                basTaskId:this.task_id
            }
            const res = await bas.basvuldel(params)
            if(res.code == 200){ 
                this.$message({
                    message:'删除成功',
                    type: 'success'
                });
                this.alldelvisible = false;
                this.getData();
                this.getbasvulstat();
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        handlesearch(){
            this.formData.page = 1;
			this.getData();
			this.currentpage = 1;
        },
        handleReset(){
            this.formData.page = 1;
            this.formData.search='';
            this.pageSize=10;
			this.getData();
			this.currentpage = 1;
        },
        handlecurrentchange(t){
            this.formData.page = t; 
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t){
            this.formData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handleShowlod(row){ //日志 
            const routeData = this.$router.resolve({
                path: '/bastargetlog',
                query: { id: row.id, target: row.addr }
            })
            window.open(routeData.href, '_blank')


        },
        handleSelectionChange(val){
            this.multipleSelection = val
        },
        mouseenter(row,colum,cell,event){  
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
            if (this.$refs['popover_id-' + row.id]) {
                let t = this.$refs['popover_id-' + row.id].showPopper;
                if(!t){
                    this.showEditFileNameButton = false;
                    this.rowId = "";
                }
            } else {
                this.showEditFileNameButton = false;
                this.rowId = "";
            }
        },
        clickButton(type) {
            switch (type) {
                case '影响级别':
                    this.$refs.class.toggleMenu();
                    break;
                case '攻击阶段':
                    this.$refs.attackStage.toggleMenu();
                    break;
                case '攻击方式':
                    this.$refs.attackMode.toggleMenu();
                    break;
                case '攻击结果':
                    this.$refs.status.toggleMenu();
                    break;
            }
        },
        async handleInfo(row){  
            const res = await bas.basruleinfo({
                basRuleId:row.ruleId
            });
            if(res.code == 200){
               
                let dt = res.data;
                this.vuldialogVisible = true;
                this.vulbasicinfo = [];
                 
                this.vulbasicinfo.push({
                    nameZh:dt.nameZh,
                    affectTarget:dt.affectTarget, 
                    classType:dt.classType,
                    attackStage:dt.attackStage,
                    attackStageEnum:dt.attackStageEnum,
                    affectScore:dt.affectScore, 
                    riskLevelEnum:dt.riskLevelEnum,
                    riskLevel:dt.riskLevel,
                })
                this.scriptLibraryInfo.id = dt.id;
                this.scriptLibraryInfo.attackMode = dt.attackMode;
                this.scriptLibraryInfo.descriptionZh = dt.descriptionZh;
                this.scriptLibraryInfo.refUrl = dt.refUrl;
                this.scriptLibraryInfo.fixSuggest = dt.fixSuggest;
                this.scriptLibraryInfo.relationAttackMethod = dt.relationAttackMethod; 
                this.scriptLibraryInfo.content = dt.content;
            }else{

            } 
        }, 
        cancalvuldialogVisible(){
            this.vuldialogVisible = false;
            this.is_Update = false;
            this.updateBugtxt = '编辑';
            this.showerr =false;
        },
        // btnEdit(){ 
        //     this.is_Update = true;
        //     this.updateBugtxt = '编辑中';  
        //     this.editInfo.riskLevel = this.vulbasicinfo[0].riskLevel;  
        //     this.editInfo.attackStage = this.vulbasicinfo[0].attackStage; 
        //     this.editInfo.affectTarget =  this.vulbasicinfo[0].affectTarget; 
        //     this.editInfo.classType = this.vulbasicinfo[0].classType;
        //     this.editInfo.affectScore = this.vulbasicinfo[0].affectScore;
        //     this.editInfo.descriptionZh = this.scriptLibraryInfo.descriptionZh;
        //     this.editInfo.refUrl = this.scriptLibraryInfo.refUrl;
        //     this.editInfo.fixSuggest = this.scriptLibraryInfo.fixSuggest; 
        //     this.editInfo.relationAttackMethod =  this.scriptLibraryInfo.relationAttackMethod; 
        //     this.editInfo.attackMode = this.scriptLibraryInfo.attackMode; 
        //     this.showerr =false; 
        // },
        // async saveEdit(){
        //     const res = await bas.basruleedit({
        //         "id": this.scriptLibraryInfo.id,
        //         "affectTarget": this.editInfo.affectTarget,
        //         "attackMode": this.editInfo.attackMode,
        //         "attackStage": this.editInfo.attackStage,
        //         "riskLevel": this.editInfo.riskLevel,
        //         "affectScore":  this.editInfo.affectScore,
        //         "relationAttackMethod":  this.editInfo.relationAttackMethod,
        //         "fixSuggest":  this.editInfo.fixSuggest,
        //         "refUrl":  this.editInfo.refUrl
        //     });
        //     if(res.code == 200){
        //         this.$message({
        //             message:'编辑成功',
        //             type: 'success'
        //         });
        //         this.vuldialogVisible = false;
        //         this.is_Update = false;
        //         this.updateBugtxt = '编辑';
        //         this.showerr =false;
        //         this.getbasvulstat();
        //     }else{
        //         this.$message({
        //             message:'编辑规则失败',
        //             type: 'error'
        //         });
        //     }
        // },
        change_vul_description(){ //漏洞描述
            if(!this.editInfo.descriptionZh){
                this.showerr = true;
            }else{
                this.showerr = false;
            }
        },
    }
}
</script>