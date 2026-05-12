<template>
  	<div >
        <div class="main-title  ">  
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >指纹库</label> 
	  	</div> 
        <!-- <div class="Statistics">
            <div class="StatisticsData">
                <div class="putaway objectstyle" v-if="!isShow" >
                    <div class="total" >
                        <i  class="iconfont iconloudongzongshu"></i>
                        <label class="" for="">指纹总数</label>
                        <span>{{objectTotal}}</span>
                    </div>
                    <div v-for="(item, i) in putawayres" :key="i"> 
                        <i :class="'vul'+item.fingerprint_type_name_num"></i>
                        <label for="">{{item.fingerprint_type_name}}</label>
                        <span class="of">{{item.fingerprint_type_of}}</span>
                        <span >{{item.fingerprint_type_sum}}</span> 
                    </div> 
                </div>
                <div class="open" v-else >
                    <div>
                        <div class="chart" id="pie"> </div>
                        <div class="otherObejct objectstyle">
                            <div v-for="(item,i) in objectlist" :key="i">
                                <i :class="'vul'+item.fingerprint_type_name_num"></i>
                                <label for="">{{item.fingerprint_type_name}}</label>
                                <span class="of">{{item.fingerprint_type_of}}</span>
                                <span >{{item.fingerprint_type_sum}}</span>
                            </div> 
                        </div>
                    </div>
                </div>
            </div>
            <div class="showhide">
                <div @click="showhide">
                    <div v-if="!isShow">
                        <label for="">展开</label>
                        <i class="iconfont iconxialaxianxing"></i>
                    </div> 
                    <div v-else>
                        <label for="">收起</label>
                        <i class="iconfont iconxialaxianxing" style="transform:rotate(180deg)"></i>
                    </div>
                </div>
            </div>
        </div> -->
        <div class="objectlist">
            <div class="search-box"> 
			    <div  class="operationbutton" >  
                    <!-- <el-button type="primary" @click="createFinger()" size="small" style="margin-right:10px;">新建</el-button> -->
                    <xzbutton  
                    type="primary" 
                    @click="createFinger"  
                    size="small"  >新建</xzbutton>  
					<!-- <el-popover
                        popper-class="delButton_popper"
                        placement="bottom-start"
                        width="170"
                        style="padding-right:8px"
                        trigger="click" 
                        :visible-arrow="false"
                        v-model="alldelvisible" >
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="" >
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false" >取消</el-button>
                            <el-button size="mini" type="primary"  @click="AllDel">确定</el-button>
                        </div>  
                        <el-button type="warning"  size="small"  slot="reference" :disabled="!multipleSelection.length">删除</el-button> 
                    </el-popover>  -->
                    <delbutton  
                    :width="170"  
                    @click="AllDel"  
                    :disabled="!multipleSelection.length" style="margin-left: 8px;"></delbutton> 
					
				</div>
				<div class="serach-condition" >
					<!-- <div  > 
						<el-select v-model="formData.type"  style=" width:150px;" clearable placeholder="指纹类型"  size="small">  
							<el-option
								v-for="(item,i) in vulobjectlist"
								:key="i"
								:label="item.name"
								:value="item.name_num"> 
							</el-option>
						</el-select>  
	                </div> -->
                    
					<div class="search-text">
						<el-input placeholder="请输入关键字"  v-model="formData.search_field" @keydown.enter.native="handlesearch" class="input-with-select"  size="small" clearable > </el-input>
						<!-- <el-button type="primary"  size="small" @click="handlesearch">搜索</el-button> -->
                        <xzbutton 
                        type="primary" 
                        @click="handlesearch" 
                        :disabled="false" 
                        size="small"  >搜索</xzbutton>
					</div>
					<div >
						<!-- <el-button type="primary"  size="small" @click="handleReset">重置</el-button> -->
                        <xzbutton 
							type="primary" 
							@click="handleReset" 
							:disabled="false" 
							size="small"  >重置</xzbutton>  
					</div>   
				</div>  
			</div>
			<el-table
				:data="tableData"  style="width: 100%" class="myTable"  height="calc(100% - 102px)"
				v-loading = "Loading"  @selection-change="handleSelectionChange"  @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
				<el-table-column   
					width="55" 
			      	type="selection" :selectable="checkSelectable">
			    </el-table-column> 
			    <el-table-column   
					prop="fingerName"
					label="指纹名称"   :show-overflow-tooltip="true">  
                    <template slot-scope="scope">
                        <!-- <div style="cursor:pointer">{{scope.row.name}}</div> -->
                        <div style="cursor:pointer"  @click="getInfo(scope.row)">{{scope.row.fingerName}}</div>
                    </template>
                </el-table-column>
				<el-table-column  
					prop="fingerClassEnum"
					label="指纹类型"  > 
                    <template slot-scope="scope" slot="header">
                        <span class="cursorPointer"  @click="clickButton('指纹类型')" :class="formData.type ? 'active' : ''">指纹类型<i class="iconfont iconshaixuan"></i></span>
                        <el-select popper-class="thSelect" v-model="formData.type"  clearable placeholder="指纹类型"  size="small" ref="vulobjectlistRef" @change="handlesearch">  
							<el-option
								v-for="(item,i) in vulobjectlist"
								:key="i"
								:label="item.label"
								:value="item.value"> 
							</el-option>
						</el-select>  
                    </template> 
				</el-table-column>
                <el-table-column   
					prop="version"
					label=" 版本"   >
				</el-table-column> 
				<el-table-column   
					prop="fingerTypeName"
					label="指纹分类"   >
				</el-table-column> 
                <el-table-column   
					prop="level"
					label="分层"   >
				</el-table-column>  
                <el-table-column  
					label="操作"  >
                    <template slot-scope="scope" > 
                        <div v-if="showEditFileNameButton && rowId == scope.row.id">
                                <!-- <el-link  :underline="false" class="link_primary"  @click="getInfo(scope.row)">详情</el-link> -->
                            <el-link  v-if="false" :underline="false" class="link_primary"  @click="handleInfo(scope.row)">验证</el-link>  
                            <el-link  :underline="false" class="link_danger linkafter" style="padding:0"  @click="getInfo(scope.row)" >详情</el-link>
                            <el-link   :underline="false" class="link_primary"  @click="handleTest(scope.row)">测试</el-link>  
                            
                            <el-popover   v-if="scope.row.source !== 1"
                                placement="bottom"
                                width="170"   
                                :ref="`popover_id-${scope.row.id}`"
                                popper-class="delButton_popper" >
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                    <el-button size="mini" type="primary" @click="handleDel(scope)">确定</el-button>
                                </div> 
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference">删除</el-link>  
                            </el-popover>
                            
                        </div>
                        <div v-else >
                            <span>{{scope.row.SourceName}}</span>
                        </div> 
					</template>

				</el-table-column> 
				<!-- <el-table-column prop="status" label="操作" min-width="10%"  >
				</el-table-column> -->
			</el-table>
			<el-pagination
				:page-size="pageSize" 
				background
				layout=" total,  prev, pager, next,sizes, jumper"
				:total="totalpage"
				:current-page="currentpage"
				@current-change = "currentchange"
				@size-change="handleSizeChange" >
			</el-pagination>
        </div>
        <!-- 验证弹窗 -->
        <el-dialog
            :title="curCheckRuleName"
            :visible.sync="vuldialogVisible"
            width="1184px"
            class="fingerValidate" 
            :close-on-click-modal="false" 
            :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="checkFinger">开始验证</el-button>
				<el-button size="small" @click="cancalvuldialogVisible">关闭</el-button>
            </div>
            <div class="buginfo_box" > 
                <div class="bugbasicinfo">
                    <el-form :model="checkFingerform" label-width="125px"  status-icon  ref="ruleFormadduser" :rules="rules1" > 
                        <el-form-item
                            label="指纹规则"
                            prop="role"  > 
                            <el-select v-model="checkFingerform.rule_name" @change="selectRule"  style=" width:320px;" clearable  size="small">  
                                <el-option
                                    v-for="(item,i) in validateRules"
                                    :key="i"
                                    :label="item.name"
                                    :value="item.name"> 
                                </el-option>
                            </el-select> 
                        </el-form-item>   
                        <el-form-item
                            label="规则内容"
                            prop="info" class="infobox" > 
                            <el-input 
                                type="textarea"
                                :disabled="true"
                                v-model="checkFingerform.rule_content"  
                                class="ruleinfo" 
                                autocomplete="off" 
                                placeholder="请输入规则内容" ></el-input> 
                        </el-form-item>   
                        <el-form-item
                            label="验证方式"
                            prop="role"  > 
                            <el-select v-model="checkFingerform.verify_mode" @change="changeVerifyModes"  style=" width:320px;" clearable  size="small">  
                                <el-option
                                    v-for="(item,i) in validateWays"
                                    :key="i"
                                    :label="item.name"
                                    :value="item.name_num"> 
                                </el-option>
                            </el-select> 
                        </el-form-item> 
                        <!-- <el-form-item
                            label="验证内容"
                            v-if="checkFingerform.rule_name.includes('图标') && Number(checkFingerform.verify_mode) === 2" class="infobox" > 
                            <el-input 
                                v-model="checkFingerform.verify_content"
                                disabled
                                type="textarea"  
                                class="ruleinfo" 
                                autocomplete="off" 
                                placeholder="请输入测试内容所在的URL" ></el-input>  -->
                        <!-- </el-form-item> -->
                        <el-form-item
                            label="验证内容"
                            prop="verify_content" class="infobox" > 
                            <el-input 
                                v-model="checkFingerform.verify_content"
                                type="textarea"  
                                class="ruleinfo" 
                                autocomplete="off" 
                                placeholder="请输入测试内容所在的URL" ></el-input> 
                        </el-form-item>  
                        <!-- <el-form-item class="uploadFormItem" label="" v-show="!(checkFingerform.rule_name.includes('正则') || checkFingerform.verify_mode === 1)"> 
                            <label class="dialog_item_label nolabelleft" style="vertical-align: top; margin-top: 7px; ">上传图标</label>
                            <input type="file"  class="btnUploadID"  id='upload'  accept="image/png,image/jpeg,image/jpg"> 
                            <el-button type="primary"  size="small" @click="clickupload()">上传图标</el-button>
                        </el-form-item> 
                        <el-form-item class="uploadFormItem"
                            label="" v-show="checkFingerform.rule_name.includes('正则') || checkFingerform.verify_mode === 1"> 
                            <label class="dialog_item_label nolabelleft" style="vertical-align: top; margin-top: 7px; ">上传图标</label>
                            <el-button type="primary" :disabled="true"  size="small">上传图标</el-button>
                        </el-form-item>     -->
                        <el-form-item
                            label="验证结果"
                            class="infobox" > 
                            <el-input
                                style="margin-top: 5px;" 
                                type="textarea"  
                                class="ruleinfo" 
                                v-model="checkFingerform.result"
                                autocomplete="off" 
                                placeholder="" ></el-input> 
                        </el-form-item>  
                </el-form>
                </div>   
            </div>
        </el-dialog>
        <!-- 详情弹窗 -->
          <!-- <el-dialog
            :title="curCheckRuleName"
            :visible.sync="vuldialogVisible1"
            width="1184px"
            class="fingerValidate" 
            :close-on-click-modal="false" 
            :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small"   @click="toEdit()">编辑</el-button>
                <el-button size="small" @click="cancalvuldialogVisible1">关闭</el-button>
            </div>
            <div class="buginfo_box" > 
                <div class="bugbasicinfo">
                    <el-form :model="detailFingerform" label-width="0"  status-icon  ref="ruleFormadduser"> 
                        <el-form-item
                        class="name"> 
                        <label class="dialog_item_label flexBet"><span>指纹名称</span><span>*</span></label>
                        <el-input 
                            style=" width:320px;"
                            v-model="detailFingerform.name"  
                            autocomplete="off" 
                            placeholder="请输入指纹名称" disabled>
                        </el-input> 
                        </el-form-item>  
                        <el-form-item
                            label=""> 
                            <label class="dialog_item_label flexBet"><span>指纹类型</span><span>*</span></label>
                            <el-select v-model="detailFingerform.rule_name" disabled  style=" width:320px;" clearable  size="small">  
                                <el-option                                 
                                    v-for="(item,i) in vulobjectlist"
                                    :key="i"
                                    :label="item.name"
                                    :value="item.name_num"> 
                                </el-option>
                            </el-select> 
                        </el-form-item>   
                        <el-form-item
                            label=" "
                            prop="info" class="infobox flexColumn" > 
                            <label class="dialog_item_label" style="height:16px;margin-right:12px">正则规则</label>
                            <div class="bugbasicinfo">
                                <el-table
                                    class="diatable"
                                    :data="detailFingerform.re_rule"
                                    size='small' 
                                    >
                                    <el-table-column
                                        prop="name"
                                        label="类型" >
                                    </el-table-column>
                                    <el-table-column
                                        prop="value"
                                        label="规则">
                                    </el-table-column>
                              </el-table>
                            </div>
                        </el-form-item>   
                        <el-form-item
                        label=" "
                        prop="info" class="infobox flexColumn" > 
                        <label class="dialog_item_label" style="height:16px;margin-right:12px;width:150px;">html内容meta</label>
                        <div class="bugbasicinfo">
                            <el-table
                                class="diatable"
                                :data="detailFingerform.meta"
                                size='small'
                                >
                                <el-table-column
                                    prop="name"
                                    label="key" >
                                </el-table-column>
                                <el-table-column
                                    label="value" 
                                    prop="value">
        
                                </el-table-column>
                          </el-table>
                        </div>
                    </el-form-item> 
                    <el-form-item
                        label=" "
                        prop="info" class="infobox flexColumn" > 
                        <label class="dialog_item_label" style="height:16px;margin-right:12px;width:150px;">html报文headers</label>
                        <div class="bugbasicinfo">
                            <el-table
                                class="diatable"
                                :data="detailFingerform.header"
                                size='small'
                                >
                                <el-table-column
                                    prop="name"
                                    label="key" >
                                </el-table-column>
                                <el-table-column
                                    label="value"
                                     prop="value">

                                </el-table-column>
                        </el-table>
                        </div>
                    </el-form-item>  
                    <el-form-item
                        label=""
                        class="infobox flexColumn" > 
                        <label class="dialog_item_label" style="margin-right:12px">指纹描述</label>
                        <el-input 
                            type="textarea"  
                            class="ruleinfo" 
                            v-model="detailFingerform.description"
                            autocomplete="off" 
                            placeholder="" ></el-input> 
                    </el-form-item>  
                </el-form>
                </div>   
            </div>
        </el-dialog> -->
        <fingerDialog v-if="showFingerDialog" :visible="showFingerDialog" 
            :vulobjectlist="vulobjectlist" 
            :dialogTitle="dialogTitle" 
            @refreshData="refreshData" 
            :editFingerInfo="editFingerInfo" 
            @handleClose="closeFingerDialog"/>
        <testfinger 
            :visible='textFinger' 
            :fingerInfo="fingerInfo"
            @handleClose="closeTestFingerDialog"
            ></testfinger>
    </div>
</template>
<style lang="less" scoped>
    .diatable{
    // margin-top: 20px;
     border:1px solid #E8E8F5;
     background-color: #FFFFFF;
     padding: 20px;
    /deep/ th{
        padding: 7px 0;
    }
    }
    .Statistics{
        background: #FFFFFF;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
        border-radius: 4px;
        padding: 16px 0;
        box-sizing: border-box;
        font-size: 0;
        .StatisticsData{
            width: calc(100% - 100px);
            display: inline-block;
            .putaway{ 
                >div.total{
                    i{
                        display: inline-block;
                        width: 14px;
                        height: 14px;
                        font-size: 14px;
                        color: #4C7AE3;
                    }
                    label{
                        display: inline-block;
                        font-size: 13px;
                        color: rgba(72, 72, 102, 0.64);
                    }
                    span{
                        display: inline-block;
                        font-size: 13px;
                        color: #4C7AE3;;
                    }
                }
            }
            .open{
                >div{
                    width: 100%;
                    .chart{
                        display: inline-block;
                        width: 20%;
                        height: 218px;
                    }
                    .otherObejct{
                        display: inline-block;
                        width: calc(100% - 20%);
                        vertical-align: top;
                        >div{
                            width: 25%;
                        }
                    }
                } 
            } 
        }
    }
    .objectstyle{
         >div{
            display: inline-block;
            width: 20%;
            font-size: 0;
            padding: 19px 0 19px  24px;
            box-sizing: border-box;
            i{
                display: inline-block;
                width: 12px;
                height: 8px; 
                border-radius: 2px;
                margin-right: 5px;
            }
            label{
                display: inline-block;
                font-size: 13px;
                color: rgba(72, 72, 102, 0.64);
                margin-right: 12px;
                width: 56px;
            }
            span{
                display: inline-block;
                font-size: 13px;
                color: rgba(72, 72, 102, 0.87);  
            }
            span.of{
                width: 40px;
                margin-right: 8px;  
            }
        }
         
    }
    .showhide{
        display: inline-block;
        width: 98px;
        height: 100%;
        color: #4C7AE3;
        >div{
            cursor: pointer;
            text-align: right;
            padding-right: 24px;
            box-sizing: border-box;
            label{
                font-size: 13px;
                display: inline-block;
                margin-right: 8px;
                cursor: pointer;
            }
            i{
                display: inline-block;
                font-size: 10px;
            }
        }
    }
    .vul1{
        background-color: #1855DF;
    }
    .vul2{
        background-color: #F86407;
    }
    .vul3{
        background-color: #FFB700;
    }
    .vul4{
        background-color: #00A65D;
    }
    .vul5{
        background-color: #09C1F7;
    }
    .vul6{
        background-color: #F74752;
    }
    .vul7{
        background-color: #5517C7;
    }
    .vul8{
        background-color: #F8858C;
    }
    .vul9{
        background-color: #FFA56D;
    }
    .vul10{
        background-color: #F9D743;
    }
    .vul11{
        background-color: #27D488;
    }
    .vul12{
        background-color: #7299F0;
    }
    .vul13{
        background-color: #9671DA;
    }
    .vul14{
        background-color: #86E4FF;
    }
    .objectlist{
        margin-top: 15px;
        padding: 24px; 
        background: #fff; 
        height: calc(100% - 39px);
        box-sizing: border-box;
	    box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
        border-radius: 4px;
    }
    .buginfo_box{
       background: #F7F7FB; 
       height:100%;
       overflow: auto;
    }
    .bugbasicinfo{
        padding: 24px;
        .el-form-item{
            margin-bottom: 14px;
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
            margin-bottom: 14px;
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
    .part_title{
        font-size: 14px; 
        margin-bottom:16px;
        font-weight: 800;
        border-left: 3px solid #4C7AE3;
        padding-left: 10px;
        height: 14px;
        line-height: 14px;
        color:rgba(72,72,102,0.87);
    }
    .errormsg{
        font-size: 12px;
        color: #F56C6C;
        padding-top: 4px;
    }
    .nolabelleft{
        border-left: none;
    }
    .flexColumn /deep/ .el-form-item__content{
        flex-direction: column;
        .dialog_item_label{
            margin-bottom: 15px;
        }
        .bugbasicinfo{
            padding: 0;
        }
    }
    .infobox /deep/ .el-form-item__content{
        display: flex;
        justify-content: center;
        margin-left: 0px;
    }
    .uploadFormItem{
        position: relative;
        /deep/ .el-button{
            position: absolute;
            left: 112px;
        }
    }
    .ruleinfo /deep/ .el-textarea__inner{
        width: 100%;
        height: 80px!important;
        background: #FFFFFF;
        border-radius: 2px;
        border: 1px solid #E8E8F5;
    }
    .isError /deep/ .el-textarea__inner{
        border: 1px solid  #F56C6C  !important;
    }
    /deep/ .el-table tr th:nth-child(1) .cell{
        padding-left: 32px;
    }
    /deep/ .el-table tr td:nth-child(1) .cell{
        padding-left: 32px;
    }
    /deep/ .el-button--small, .el-button--small.is-round {
        padding: 8px 8px;
        vertical-align: middle;
    }
    /deep/ .buginfo_box .el-input--small .el-input__inner {
        // height: 40px;
    }
    
    /deep/ .myTable{
        thead {
            .cursorPointer{
                cursor: pointer;
                &.active{
                    color:#4C7AE3;
                    i{
                        color:#4C7AE3;
                    }
                }
            }
            .cell{
                line-height: 15px;
                >span{
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
</style>
<script>
// var echarts = require('echarts');
import fingerDialog from '@/components/fingerDialog.vue'
import xzbutton from "@/components/XzButton.vue"; 
import delbutton from "@/components/DelButton.vue";
import testfinger from './components/testFinger.vue'
import $ from 'jquery'
import { fingerprint } from '@/api/tool.js'
export default ({
    name:'FingerprintLib',
    components: {
        fingerDialog,
        xzbutton,
		delbutton,
        testfinger
    },
    data(){
        var validateContent = (rule, value, callback) => {
            let strRegex = /((([A-Za-z]{3,9}:(?:\/\/)?)(?:[\-;:&=\+\$,\w]+@)?[A-Za-z0-9\.\-]+|(?:www\.|[\-;:&=\+\$,\w]+@)[A-Za-z0-9\.\-]+)((?:\/[\+~%\/\.\w\-_]*)?\??(?:[\-\+=&;%@\.\w_]*)#?(?:[\.\!\/\\\w]*))?)/
            var re=new RegExp(strRegex);
            if (!re.test(value)) {
                if(this.checkFingerform.rule_name.includes('正则') && this.checkFingerform.verify_mode === 2) {
                    callback();
                }
                return callback(new Error('网址格式不正确'));
            } else {
                callback();
            }
        };
        return{
            showEditFileNameButton:false,
            rowId:'',
            dialogTitle:'新建指纹',
            checkFingerform:{
                id: '',
            	rule_name:'',
            	rule_content:'',
            	verify_mode:'',
            	verify_content:'',
                // fingerprint_info: '',
                request_id: '',
                result: '',
                headers: [],
                meta: [],
                re_rule: [],
                ruleList: []
            },
            detailFingerform:{
                id: '',
            	name:'',
                rule_name:'',
            	description:'',
            	meta:[],
            	header:[],
                re_rule: [],
                type_zh: '',
            },
            // rules:{ 
            //     target:[
            //         { required: true, message: '任务目标不能为空', trigger: 'blur' },  
            //     ],  
            //     taskname:[
            //         { required: true, message: '任务名称不能为空', trigger: ['blur', 'change'] },
            //     ],
            //     port:[
            //         { required: true, validator: validatePort, trigger: ['blur'] },
            //     ],
            // },
            rules1:{ 
                verify_content: [
                    { required: true, message: '验证内容不能为空', trigger: 'blur' }
                    // { validator: validateContent, trigger: 'blur' } 
                ]
            },
            multipleSelection:[],
            alldelvisible:false,
            isShow:false,
            objectTotal:0,
            putawayres:[],
            objectlist:[],
            color:['#1855DF','#F86407','#FFB700','#00A65D','#09C1F7','#F74752','#5517C7','#F8858C','#FFA56D','#F9D743','#27D488','#7299F0','#9671DA','#86E4FF'],
            pageSize:10,
            totalpage:0,
			currentpage:1, 
            formData:{
                page_num:1,
                search:'',
                type:'',
                search_field:'',
            },
            vulobjectlist:[],
            validateRules:[],
            validateWays: [],
            Loading:false,
            tableData:[],
            vuldialogVisible:false,
            vuldialogVisible1:false,
            is_vulUpdate:false,
            curCheckRuleName:'',
            updateBugtxt:'编辑',
            vulbasicinfo:[],
            risklist:[],
            showerr:false,
            showFingerDialog: false,
            editFingerInfo: null,
            timeoutEvent: null,
            textFinger:false,
            fingerInfo:null,
        }
    },
    created() {
        this.$store.state.activefirstMenu="/fingerprint"; 
        this.pageSize = this.commonjs.pageSize;
    },
    mounted(){
        this.getVulObjectlist();
        // this.getObjectStatistics();
        this.getObjectData(); 
    },
    methods: {
        icons(h, { column }) {
            let that = this
            return h(
              "div",

              [
                h("span", column.label),
                h(
                  "i",
                  {
                    slot: "reference",
                    class: "iconfont iconshaixuan",
                    style:"color:rgba(72,72,102,0.32);margin-left:5px;vertical-align:initial",
                    on: {
                        click: function() {
                            that.clickButton(column);
                        }
                    }
                  },
                  ""
                )
              ]
            );
          },
        // 图标点击事件
        clickButton(type) {
            switch (type) {
                case '指纹类型':
                this.$refs.vulobjectlistRef.toggleMenu();
                break;
            }
        },
        showhide(){
            this.isShow = !this.isShow;
            if(this.isShow){
                
                this.$nextTick(() => {
                   this.echartpie( this.objectlist);
                //  this.echartpie([])
                }); 
            }
        },
        handleSelectionChange(val){
            this.multipleSelection = val
        },
        // 删除规则
        async handleDel (scope){ //删除 
            let params = {
                ids:[scope.row.id].join(',')
            }
            const res = await fingerprint.deleteFingerprint(params)
            if(res.code == 200){ 
                this.$message({
                    message:res.msg,
                    type: 'success'
                });
                this.$message({
                    message:'删除成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                this.getObjectData(); 
                
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        async AllDel(){ //批量删除
			if(this.multipleSelection.length == 0) return;
    		var ids = [];
    		for (var i = 0; i < this.multipleSelection.length; i++) {
    			ids.push(this.multipleSelection[i].id);
    		}
            let params = {
                ids:ids.join(',')
            }
            const res = await fingerprint.deleteFingerprint(params)
            if(res.code == 200){ 
                this.$message({
                    message:'删除成功',
                    type: 'success'
                });
                this.alldelvisible = false;
                this.getObjectData(); 
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
		},
        // async getObjectStatistics(){ //漏洞对象统计
        //     let params = {}
        //     const res = await fingerprint.getObjectStatistics(params)
        //     if(res.success){ 
        //         this.objectTotal = res.fingerprint_sum;
        //         this.objectlist = res.results;
        //         this.putawayres = res.results.slice(0,4);
        //     }else{
        //         this.$message({
        //             message:res.msg,
        //             type: 'error'
        //         });
        //     }
        // },
        echartpie(res){
            let dt = [];
            for(var i=0;i<res.length;i++){
                dt.push( {value: res[i].fingerprint_type_sum, name: res[i].fingerprint_type_name})
            }
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('pie')); 
            // 绘制图表 
            myChart.setOption({ 
                tooltip: {
                    trigger: 'item',
                    formatter: '{b}: {c} ({d}%)'
                }, 
                title:{
                    text:this.objectTotal, 
                    textStyle:{
                        color:'#4C7AE3',
                        fontSize:24,
                        fontWeight:'500'
                    },
                    subtext:'指纹总数(个)',
                    subtextStyle:{
                        color:'rgba(72, 72, 102, 0.64)',
                        fontSize:11,
                    },
                    left:'center',
                    top:'middle'
                },
                color:this.color,
                series: [
                     {
                        name: '',
                        type: 'pie',
                        radius: ['60%', '80%'],
                        data:[{value:0,name:''}],  
                        avoidLabelOverlap: true,
                        hoverAnimation:false,
                        label: {
                            show: false,
                            position: 'center'
                        }, 
                        labelLine: {
                            show: false
                        },
                        itemStyle: {
                            normal: {
                                color: 'rgba(0,0,0,.12)'
                            }
                        }
                    },
                    {
                        name: '漏洞对象',
                        type: 'pie',
                        radius: ['60%', '80%'],
                        avoidLabelOverlap: true,
                        hoverAnimation:false,
                        label: {
                            show: false,
                            position: 'center'
                        }, 
                        labelLine: {
                            show: false
                        },
                        data: dt
                    },
                   
                ]
            }); 
            window.onresize = myChart.resize;  //窗口大小变化自适应
        },
        async getVulObjectlist(){ //指纹类型下拉
            let params = {}
            const res = await fingerprint.getVulObjectlist(params)
            console.log(res,'指纹类型下拉');
            if(res.code == 200){  
                this.vulobjectlist = res.data.class;
                // this.vulobjectlist.unshift({
                //         name: '全部',
                //         name_num: ''
                //     })
            }else{
                // this.$message({
                //     message:res.msg,
                //     type: 'error'
                // });
            }
        },
        async getObjectData(){ //指纹列表
            let params = {
                fingerName:this.formData.search_field,
                class:this.formData.type,
                page:this.formData.page_num,
                size:this.pageSize
            }
            const res = await fingerprint.getObjectData(params) 
            if(res.code == 200){ 
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }else{
                // this.$message({
                //     message:res.msg,
                //     type: 'error'
                // });
            }  
        },
        handlesearch(){
            this.formData.page_num = 1;
			this.getObjectData();
			this.currentpage = 1;
        },
        handleReset(){ //重置
            this.formData.search_field = ''; 
			this.formData.page_num = 1;
			this.formData.type = '';
            this.currentpage = 1;
            this.pageSize=10;
			this.getObjectData();
        },
        currentchange(t){
            this.formData.page_num = t; 
            this.getObjectData();
            this.currentpage = t;
        },
        handleSizeChange(t){
            this.formData.page_num = 1;
            this.pageSize = t;
            this.getObjectData();
        },
        //获取详情
        async getInfo(row){ 
            this.dialogTitle = '指纹详情'
            this.showFingerDialog = true
            this.editFingerInfo =  row;
            // let params = {
            //     id:row.id
            // }
            // const res = await fingerprint.getInfo(params)
            // if(res.success){  
            //     this.showFingerDialog = true
            //     this.editFingerInfo =  res.results;
            // }else{
            //     this.$message({
            //         message:res.msg,
            //         type: 'error'
            //     });
            // }

            
        },
        // 验证
        async handleInfo(row){
            let id = row.id
            this.curCheckRuleName = row.name
            let params = {
                id:id
            }
            const res = await fingerprint.getInfo(params)
            if(res.success){
                this.vuldialogVisible = true;
                this.getvulrisk();
                this.validateRules = []
                let results = res.results
                this.checkFingerform.id = id
                this.checkFingerform.request_id = results.request_id
            //    this.checkFingerform.fingerprint_info = results.info
                this.checkFingerform.meta = results.meta
                this.checkFingerform.headers = results.headers
                this.checkFingerform.re_rule = results.re_rule
                results.re_rule.forEach((item, index) => {
                this.validateRules.push({
                    name: `正则规则${index+1}`,
                    value: item.value
                })
                })
                results.meta.forEach((item, index) => {
                this.validateRules.push({
                    name: `meta规则${index+1}`,
                    value: item.value
                })
                })
                results.headers.forEach((item, index) => {
                this.validateRules.push({
                    name: `headers规则${index+1}`,
                    value: item.value
                })
                })
            //    if (results.ico.is_open) {
            //        results.ico.rules.forEach(item => {
            //            this.validateRules.push({
            //                name: `图标_${item.name}`,
            //                name_num: `图标_${item.name}`,
            //                content: item.content
            //            })
            //        })
            //    }
            //    if (results.regex.is_open) {
            //        results.regex.rules.forEach(item => {
            //            this.validateRules.push({
            //                name: `正则_${item.name}`,
            //                name_num: `正则_${item.name}`,
            //                content: item.feature + '\r\n' + item.version
            //            })
            //        })
            //    }
            //    this.editFingerInfo =  dt.results;
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        selectRule (name) {
            let obj = this.validateRules.find(item => {
                return item.name === name
            })
            if (obj) {
                this.checkFingerform.rule_content = obj.value
                // if (this.checkFingerform.rule_name.includes('图标') && this.checkFingerform.verify_mode === 2) {
                //     this.checkFingerform.verify_content = ''
                //     console.log('this.$refs.validateForm', this.$refs.validateForm)
                //     this.$refs.ruleFormadduser.clearValidate('verify_content');
                // }
            }
        },
        async getvulrisk(){ //验证里的验证方式下拉
            let params = {}
            const res = await fingerprint.getvulrisk(params)
            if(res.success){  
                this.validateWays = res.results;
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        // async checkFinger(){ //验证指纹
        //     if (this.checkFingerform.rule_name.includes('图标') && Number(this.checkFingerform.verify_mode) === 2) {
        //         this.checkFingerform.upload_ico = document.querySelector('#upload').files[0]
        //     } else {
        //         this.checkFingerform.upload_ico =''
        //     }
        //     let formDate = new FormData()
        //     formDate.append('id', this.checkFingerform.id);
        //     formDate.append('rule_content', this.checkFingerform.rule_content);
        //     formDate.append('verify_mode', this.checkFingerform.verify_mode);
        //     formDate.append('verify_content', this.checkFingerform.verify_content);
        //     formDate.append('request_id', this.checkFingerform.request_id);
        //     // if(!this.showerr){
        //     const res = await fingerprint.checkFinger(formDate)
        //             if(res.success){
        //                 this.beginCheckFinger()                                    
        //             }else{
        //                 this.$message({
        //                     message:res.msg,
        //                     type: 'error'
        //                 });
        //             }   
        // },
        // async beginCheckFinger () {
        //     let params = {
        //         id : this.checkFingerform.id,
        //         request_id :this.checkFingerform.request_id
        //     }
        //     const res = await fingerprint.beginCheckFinger(params)
        //         if(res.success){ 
        //             if (res.result.length > 0) {
        //                 this.$message({
        //                     message:res.msg,
        //                     type: 'success'
        //                 });
        //                 this.checkFingerform.result = JSON.parse(res.result[0]).msg
        //             } else {
        //                 this.timeoutEvent = setTimeout(() => {
        //                     this.beginCheckFinger()
        //                 }, 3000)
        //             }
        //         }else{
        //             this.$message({
        //                 message:res.msg,
        //                 type: 'error'
        //             });
        //         }  
        // },
        checkFinger(){ //验证指纹
            //  var _v={ 
            //     id: this.checkFingerform.id,
            //     rule_name:this.checkFingerform.rule_name,
            //     rule_content:this.checkFingerform.rule_content ,
            //     verify_mode:this.checkFingerform.verify_mode ,
            //     verify_content:this.checkFingerform.verify_content,
            //     upload_ico:this.checkFingerform.upload_ico,
            //     request_id: this.checkFingerform.request_id,
            // }
            if (this.checkFingerform.rule_name.includes('图标') && Number(this.checkFingerform.verify_mode) === 2) {
                this.checkFingerform.upload_ico = document.querySelector('#upload').files[0]
            } else {
                this.checkFingerform.upload_ico =''
            }
            let formDate = new FormData()
            formDate.append('id', this.checkFingerform.id);
            // formDate.append('fingerprint_info', this.checkFingerform.fingerprint_info);
            // formDate.append('rule_name', this.checkFingerform.rule_name);
            formDate.append('rule_content', this.checkFingerform.rule_content);
            formDate.append('verify_mode', this.checkFingerform.verify_mode);
            formDate.append('verify_content', this.checkFingerform.verify_content);
            // formDate.append('upload_ico', this.checkFingerform.upload_ico);
            formDate.append('request_id', this.checkFingerform.request_id);
            // if(!this.showerr){
                this.$ajax({
                    method:'post',
                    url:'/tools/fingerprint/verify/',
                    data: formDate
                })
                .then(dt => {  
                    let res = dt.data;
                    if(res.success){
                        this.beginCheckFinger()                                    
                    }else{
                        this.$message({
                            message:res.msg,
                            type: 'error'
                        });
                    }   
                })
                .catch(data=>{ 
                });
            // } 
        },
        beginCheckFinger () {
          this.$ajax.get('/tools/fingerprint/verify/result/',{
                params: {
                    id : this.checkFingerform.id,
                    request_id :this.checkFingerform.request_id  
                }
            })
            .then(dt2 => {  
                let res2 = dt2.data;
                if(res2.success){ 
                    if (res2.result.length > 0) {
                        this.$message({
                            message:res2.msg,
                            type: 'success'
                        });
                        this.checkFingerform.result = JSON.parse(res2.result[0]).msg
                    } else {
                        this.timeoutEvent = setTimeout(() => {
                            this.beginCheckFinger()
                        }, 3000)
                    }
                }else{
                    this.$message({
                        message:res2.msg,
                        type: 'error'
                    });
                }   
            })
            .catch(data=>{ 
            });
        },
        cancalvuldialogVisible1(){
            this.vuldialogVisible1 = false;
        },
        cancalvuldialogVisible(){
            clearTimeout(this.timeoutEvent)
            this.vuldialogVisible = false;
            this.checkFingerform.id= '';
            this.checkFingerform.rule_name =''; 
            this.checkFingerform.rule_content = '';
            // this.checkFingerform.fingerprint_info = '';
            this.checkFingerform.verify_mode = '';
            this.checkFingerform.verify_content = ''; 
            this.checkFingerform.result = ''; 

            // let file = document.getElementById('upload')
            // file.outerHTML = file.outerHTML
            this.is_vulUpdate = false; 
            this.updateBugtxt = '编辑';
            this.showerr =false;
        },
        change_vul_description(){ //漏洞描述
            if(!this.editVulinfo.vul_description){
                this.showerr = true;
            }else{
                this.showerr = false;
            }
        },
        createFinger(){
            this.dialogTitle = '新建指纹'
            this.editFingerInfo = null
            this.showFingerDialog = true
           
        },
        handleTest(row){ //测试
            this.textFinger = true;
            this.fingerInfo = row; 
        },
        async toEdit (item) {
            this.dialogTitle = '编辑指纹'
            let params = {
                id:item.id
            }
            const res = await fingerprint.toEdit(params)
            // this.$ajax.get('/tools/fingerprint/detail/',{
            //     params: {
            //         id:item.id
            //     }
            // })
            // .then((res) => { 
            //     var dt = res.data; 
                if(res.success){  
                    this.showFingerDialog = true
                    this.editFingerInfo =  res.results;
                }else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }  
            // })
            // .catch((error) => {
            //     console.log(error);
            // })
        },
        closeFingerDialog () {
            this.showFingerDialog = false
        },
        closeTestFingerDialog(){
            this.textFinger = false;
        },
        refreshData () {
            this.showFingerDialog = false
            this.getObjectData();
            //  this.getObjectStatistics();
        },
        clickupload(){  
            document.querySelector('.btnUploadID').click();
        },
        // 上传图标更换
    // handleChangeIcon (e) {
    //         let that = this
    //         var file = e.target.files[0]
    //         that.checkFingerform.upload_ico = file
    //         // console.log('file', file)
    //         // console.log('that.generateHashForm.upload_ico', that.generateHashForm.upload_ico)
    // },
    changeVerifyModes (val) {
        // if (val === 1) {
        //     let file = document.getElementById('upload')
        //     file.outerHTML = file.outerHTML
        // } else if (val === 2) {
        //     this.checkFingerform.verify_content = ''
        //     if (this.checkFingerform.rule_name.includes('正则')) {
        //             this.$refs.ruleFormadduser.clearValidate('verify_content');
        //         }
        // }
    },
    checkSelectable(row) {
        return row.source !== 1
    },
    mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
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
            }
        },
  },
})
</script>