<template>
    <!-- 剧本库 -->
    <div class="outContainer" > 
        <div class="main-title  ">  
            剧本库 
	  	</div>
        <div class="list_box">
            <div class="search-box"> 
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
				:data="tableData"  style="width: 100%"  class="myTable"  @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" >
			    <el-table-column   
					prop="nameZh"
					label="剧本名称" :show-overflow-tooltip="true"> 
                    <template slot-scope="scope" > 
                        <span v-if="scope.row.nameZh!=''">{{scope.row.nameZh}}</span>
                        <span v-else>{{scope.row.name}}</span>
                    </template>
				</el-table-column>   
                <el-table-column prop="attackModeEnum" label="分类"  >
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('分类')"
                            :class="(formData.classType !== '' && formData.classType !== 0) ? 'active' : ''">分类<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.classType" clearable
                        size="small" ref="classtype" @change="handlesearch">
                        <el-option v-for="(item, index) in classtypelist" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                         {{ scope.row.classType }}
                    </template>
                </el-table-column> 
                <el-table-column prop="attackModeEnum" label="攻击方式"  >
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('攻击方式')"
                            :class="(formData.attackType !== '' ) ? 'active' : ''">攻击方式<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.attackType" clearable
                        size="small" ref="attackMode" @change="handlesearch">
                        <el-option v-for="(item, index) in attackTypelist" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                         {{ scope.row.attackModeEnum }}
                    </template>
                </el-table-column> 
                <el-table-column prop="classEnum" label="影响级别" show-overflow-tooltip >
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('影响级别')"
                            :class="(formData.riskLevel !== '' && formData.riskLevel !== 0) ? 'active' : ''">影响级别<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.riskLevel" clearable
                        size="small" ref="class" @change="handlesearch">
                        <el-option v-for="(item, index) in riskLevellistAll" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                         {{ scope.row.riskLevelEnum }}
                    </template>
                </el-table-column>   
                <el-table-column prop="cve" label="攻击阶段"  >  
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('攻击阶段')"
                            :class="(formData.attackStage !== '' && formData.attackStage !== 0) ? 'active' : ''">攻击阶段<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.attackStage" clearable
                        size="small" ref="attackStage" @change="handlesearch">
                        <el-option v-for="(item, index) in attackStagelistAll" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope" >  
                        <div v-if="showEditFileNameButton && rowId == scope.row.id"> 
                            <el-link  :underline="false" class="link_primary"   @click="handleInfo(scope.row)">详情</el-link>   
                        </div>
                        <div v-else >
                             {{scope.row.attackStageEnum}} 
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
  
        <div class="add_pop" id="infodialog">
            <div class="test_dialog">
                <div class="close" @click="closeDialog">
                    <i class="iconfont iconquxiao"></i>
                </div>
                <div class="el-dialog__header">
                    <div class="dialog__title clearfix">
                        <span class="title_name"> {{dialogtitle}} </span>
                    </div>
                </div>
                <div class="dialog_b_btn"> 
                    <el-button size="small" @click="closeDialog">取消</el-button>
                </div>
                <div class="el-dialog__body">
                    {{ dialogcontent }}
                </div>
            </div>
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
                <el-button size="small" @click="btnEdit">{{updateBugtxt}}</el-button>
                <el-button size="small" @click="saveEdit" v-if="is_Update">保存</el-button>
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
                        <el-table-column
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
                                        v-for="(item,index) in attackStagelist"
                                        :key="index"
                                        :label="item.label"
                                        :value="item.value">
                                    </el-option>
                                </el-select>
                            </template>
                        </el-table-column>
                        
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
    .list_box{
        background: #FFFFFF;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
        padding: 24px;
        border-radius: 4px;
        box-sizing: border-box;
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
    .test_dialog {
        width: 100%;
        height: 100%;
        position: relative;

        .close {
            cursor: pointer;
            width: 54px;
            height: 54px;
            text-align: center;
            line-height: 54px;
            position: absolute;
            left: -54px;
            top: 70px;
            background: #4c7ae3;
            color: #fff;

            i {
                font-size: 22px;
            }
        }

        /deep/ .el-dialog__header {
            border-radius: 4px 0px 0px 0px;
        }

        .dialog__title {
            font-size: 14px;

            box-sizing: border-box;

            .title_name {
                display: inline-block;
                border-left: 2px solid #fff;
                background: hsla(0, 0%, 100%, 0.12);
                color: #fff;
                height: 32px;
                line-height: 32px;
                text-align: center;
                padding: 0 24px;
            }

            .title_item {
                // float: left;
                // width: 200px;
                margin: 8px 24px 8px 0;

                .micon {
                    color: rgba(255, 255, 255, 0.7);
                    vertical-align: top;
                    height: 20px;
                    line-height: 19px;
                    display: inline-block;
                }

                .lbname {
                    display: inline-block;
                    // width: 54px;
                    margin-right: 18px;
                    margin-left: 8px;
                    color: rgba(255, 255, 255, 0.7);
                    height: 20px;
                    line-height: 20px;
                    vertical-align: middle;
                }

                .spvalue {
                    // width: calc(100% - 100px) !important;
                    display: inline-block;
                    overflow: hidden;
                    white-space: nowrap;
                    text-overflow: ellipsis;
                    color: #fff;
                    vertical-align: text-top;
                    height: 20px;
                    line-height: 20px;
                    position: relative;
                    vertical-align: middle;
                }

                .name {
                    display: block;
                    font-size: 14px;
                    width: 100%;
                    margin-bottom: 10px;
                    border-left: 3px solid #4c7ae3;
                    padding-left: 10px;
                    padding-right: 10px;
                    box-sizing: border-box;
                    color: #fff;

                    span {
                        display: block;
                        float: right;
                        font-size: 12px;
                        border-radius: 10px;
                        color: #fff;
                        padding: 2px 8px;
                    }

                    .spcor1 {
                        background: #09c1f7;
                    }

                    .spcor2 {
                        background: #15c53d;
                    }

                    .spcor3 {
                        background: #f35f28;
                    }
                }
            }
        } 
    }
    .outContainer {
        position: relative;
        min-height: calc(100% - 40px);
        overflow: hidden;
        padding: 0 !important;
        margin: 20px 15px;
        height: auto !important;
    }
    .add_pop{
        position: absolute;
        // right: 0;
        right: -635px;
        top: 50px;
        width: 580px;
        height: 100%;
        background: #fff;
        z-index: 88;
        box-shadow: 0px 8px 32px 1px rgba(76, 122, 227, 0.12);
        border-radius: 4px 0px 0px 0px;
    }
    .openAnimate {
        right: 0px;
        box-shadow: 0px 8px 32px 0px rgba(76, 122, 227, 0.12);
        animation: openAnimate 0.5s;
        -webkit-animation: openAnimate 0.5s;
        animation-fill-mode: forwards;
    }

    .closeAnimate {
        right: -640px;
        animation: closeAnimate 0.5s;
        -webkit-animation: closeAnimate 0.5s;
        animation-fill-mode: forwards;
    }
    @keyframes openAnimate {
        0% {
            right: -635px;
        }

        100% {
            right: 0px;
        }
    }

    @keyframes closeAnimate {
        0% {
            right: 0px;
        }

        100% {
            right: -635px;
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
    name:'scriptLibrary',
    data(){
        return{
            formData:{
                search:'',
                page:1,
                riskLevel:'',
                attackStage:'',
                attackType:'',
                classType:'',
            },
            pageSize:10,
            currentpage:1,
            totalpage:0,
            tableData:[],
            showEditFileNameButton:false,
            rowId:'', 
            ruleenum:[], 
            dialogtitle:'',
            dialogcontent:'',
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
            attackStagelistAll:[],
            attackStagelist:[],
            riskLevellistAll:[],
            riskLevellist:[],
            attackTypelist:[],
            classtypelist:[],
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/scriptLibrary';
    },
    mounted(){
        this.basruleenum();
        this.getData();
    },
    methods:{
        async basruleenum(){
            const res = await bas.basruleenum();
           
            if(res.code == 200){
                let newall = [{ label: "全部", value: '' }];
                this.ruleenum = res.data.class;
                this.riskLevellist = res.data.riskLevel;
                this.riskLevellistAll = newall.concat(res.data.riskLevel)
               
                this.attackStagelist = res.data.attackStage;
                this.attackStagelistAll = newall.concat(res.data.attackStage); 
               
                this.classtypelist = newall.concat(res.data.class);

                this.attackTypelist = newall.concat(res.data.attackType);
                 
            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
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
                case '分类':
                    this.$refs.classtype.toggleMenu();
                    break;
            }
          },
        async getData(){
            const res = await bas.getBasrule({
                page:this.formData.page,
                size:this.pageSize,
                search:this.formData.search,
                riskLevel:this.formData.riskLevel,
                attackStage:this.formData.attackStage,
                attackType:this.formData.attackType,
                classType:this.formData.classType,
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
        handlesearch(){
            this.formData.page = 1;
			this.getData();
			this.currentpage = 1;
        },
        handleReset(){
            this.formData.page = 1;
            this.formData.search='';
            this.pageSize=10;
            this.formData.riskLevel = '';
            this.formData.attackType='';
            this.formData.classType='';
            this.formData.attackStage='';
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
        mouseenter(row,colum,cell,event){  
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
             
            this.showEditFileNameButton = false;
            this.rowId = "";
             
        },
        handleInfo(row){
            // let obj = document.getElementById('infodialog');
            // obj.classList.add("openAnimate");
            // obj.classList.remove('closeAnimate');
            // this.dialogtitle = row.nameZh;
            // this.dialogcontent = JSON.stringify(row);
            console.log(row)
            this.vuldialogVisible = true;
            this.vulbasicinfo = [];
            this.vulbasicinfo.push({
                nameZh:row.nameZh,
                affectTarget:row.affectTarget, 
                classType:row.classType,
                attackStage:row.attackStage,
                attackStageEnum:row.attackStageEnum,
                affectScore:row.affectScore, 
                riskLevelEnum:row.riskLevelEnum,
                riskLevel:row.riskLevel,
            })
            this.scriptLibraryInfo.id = row.id;
            this.scriptLibraryInfo.attackMode = row.attackMode;
            this.scriptLibraryInfo.descriptionZh = row.descriptionZh;
            this.scriptLibraryInfo.refUrl = row.refUrl;
            this.scriptLibraryInfo.fixSuggest = row.fixSuggest;
            this.scriptLibraryInfo.relationAttackMethod = row.relationAttackMethod; 
            this.scriptLibraryInfo.content = row.content;


        },
        closeDialog(){
            // let obj = document.getElementById('infodialog');
            // obj.classList.add("closeAnimate");
            // obj.classList.remove('openAnimate');
            // this.dialogtitle =''; 
            // this.dialogcontent =''
        },
        cancalvuldialogVisible(){
            this.vuldialogVisible = false;
            this.is_Update = false;
            this.updateBugtxt = '编辑';
            this.showerr =false;
        },
        btnEdit(){ 
            this.is_Update = true;
            this.updateBugtxt = '编辑中'; 
            // this.editInfo.vul_id= this.scriptLibraryInfo.vul_id;   
            this.editInfo.riskLevel = this.vulbasicinfo[0].riskLevel;  
            this.editInfo.attackStage = this.vulbasicinfo[0].attackStage; 
            this.editInfo.affectTarget =  this.vulbasicinfo[0].affectTarget; 
            this.editInfo.classType = this.vulbasicinfo[0].classType;
            this.editInfo.affectScore = this.vulbasicinfo[0].affectScore;
            this.editInfo.descriptionZh = this.scriptLibraryInfo.descriptionZh;
            this.editInfo.refUrl = this.scriptLibraryInfo.refUrl;
            this.editInfo.fixSuggest = this.scriptLibraryInfo.fixSuggest; 
            this.editInfo.relationAttackMethod =  this.scriptLibraryInfo.relationAttackMethod; 
            this.editInfo.attackMode = this.scriptLibraryInfo.attackMode; 
            this.showerr =false; 
        },
        async saveEdit(){
            const res = await bas.basruleedit({
                "id": this.scriptLibraryInfo.id,
                "affectTarget": this.editInfo.affectTarget,
                "attackMode": this.editInfo.attackMode,
                "attackStage": this.editInfo.attackStage,
                "riskLevel": this.editInfo.riskLevel,
                "affectScore":  this.editInfo.affectScore,
                "relationAttackMethod":  this.editInfo.relationAttackMethod,
                "fixSuggest":  this.editInfo.fixSuggest,
                "refUrl":  this.editInfo.refUrl
            });
            if(res.code == 200){
                this.$message({
                    message:'编辑成功',
                    type: 'success'
                });
                this.vuldialogVisible = false;
                this.is_Update = false;
                this.updateBugtxt = '编辑';
                this.showerr =false;
                this.getData();
            }else{
                this.$message({
                    message:'编辑规则失败',
                    type: 'error'
                });
            }
        },
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