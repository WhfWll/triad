<template>
  	<div > 
  		<div class="main-title  ">  
			生成报告
	  	</div> 
        <div class="reportbox"> 
            <div class="operationbox whitebg">
                <!-- <el-button type="primary"  size="small"  @click="submitCreateReport" >生成报告</el-button> -->
                <xzbutton 
                    type="primary" 
                    @click="submitCreateReport"  
                    size="small"  >生成报告</xzbutton>  

                <!-- <el-button type="primary"  size="small"  @click="reset" >重置</el-button> -->
                <xzbutton 
                    type="primary" 
                    @click="reset"  
                    size="small" style="margin-left:8px">重置</xzbutton>  
                <label for="">提示：</label>
                <span>任务报告是渗透任务中所有测试目标的综合分析报告，目标报告是单个目标的测试报告。</span>
            </div>
            <el-form :model="reportform" label-width="0"  status-icon   ref="reportform" :rules="reportrules"> 
                <el-row :gutter="15">
                    <el-col :span="12">
                    <div >
                        <div class="grid-title color3"  > 
                            <label for="">报告基础参数设置</label>
                            <img src="../../assets/rp.png" alt="">
                        </div>
                        <div class="grid-content" style="height:588px" >
                            <div style="display: inline-block;" v-if="!task_id">
                                <el-form-item
                                    label=""  > 
                                    <label class="dialog_item_label">任务记录</label>
                                    <el-select
                                        v-model="reportform.num"
                                        size="small" 
                                        placeholder="请选择"
                                        style=" width: 150px;" @change="changeTaskNum" >
                                        <el-option  label="最近5条" value="5" ></el-option>
                                        <el-option  label="最近10条" value="10" ></el-option>
                                        <el-option  label="最近20条" value="20" ></el-option>
                                        <el-option  label="最近50条" value="50" ></el-option>
                                        <el-option  label="最近100条" value="100" ></el-option>
                                        <el-option  label="全部" value="0" ></el-option>
                                    </el-select>   
                                </el-form-item> 
                            </div>
                            <!-- <div style="display: inline-block;margin-left:72px;">
                                <el-form-item
                                    label=" " > 
                                    <label class="dialog_item_label">测试主体</label>
                                    <el-input
                                        v-model="reportform.name"
                                        size="small"
                                        style="width:150px"
                                        placeholder="测试主体"
                                        maxlength="50" 
                                    ></el-input>
                                </el-form-item> 
                            </div>
                            <div style="display: inline-block;margin-left:72px;">
                                <el-form-item
                                    label=" " > 
                                    <label class="dialog_item_label">测试时间</label>
                                    <el-date-picker
                                        v-model="reportform.time"
                                        type="datetime"
                                        placeholder="测试时间">
                                    </el-date-picker>
                                </el-form-item> 
                            </div> -->
                            <el-form-item
                                label=" " class="bottomsmall"> 
                                <label class="dialog_item_label">报告格式</label>
                                <el-checkbox-group v-model="reportform.reportstyle"  style=" display: inline-block;line-height: 30px;" @change="changeReportForm">
                                    <el-checkbox label="html" name="html">HTML</el-checkbox>
                                    <el-checkbox label="pdf" name="pdf">PDF</el-checkbox>
                                    <el-checkbox label="word" name="word">Word</el-checkbox>
                                    <el-checkbox label="excel" name="excel">Excel</el-checkbox> 
                                    <el-checkbox label="csv" name="csv">CSV</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                            <el-form-item label=" " class="bottomsmall"> 
                                <label class="dialog_item_label">输出范围<i class="is-required" style="float:right;">*</i></label> 
                                <ul class="outputRange">
                                    <li v-for="(item,i) in outputRangelist" :key=i :class="[{'checked':item.isChecked==1},{'':item.isChecked==0}]" @click="changeChecked(i)">
                                        <i class="el-icon-check" ></i> 
                                        <label for="" >{{item.name}}</label>
                                    </li>
                                    <!-- <li class="checked" >
                                        <i class="el-icon-check"></i>
                                        <label for="">高危目标</label>
                                    </li>
                                    <li >
                                        <i class="el-icon-check"></i>
                                        <label for="">中危目标</label>
                                    </li>
                                    <li >
                                        <i class="el-icon-check"></i>
                                        <label for="">低危目标</label>
                                    </li>
                                    <li >
                                        <i class="el-icon-check"></i>
                                        <label for="">安全目标</label>
                                    </li> -->
                                </ul>
                            </el-form-item>
                        
                            <div class="invoice-list"> 
                                <div class="table_head">
                                    <span style=""></span> 
                                    <ul class="invoice-header"> 
                                        <li class="invoice-item task1">测试目标</li>
                                        <!-- <li class="invoice-item">目标风险</li> -->
                                        <li class="invoice-item task2">存活目标</li>
                                        <li class="invoice-item task3">更新时间</li> 
                                    </ul>
                                </div>
                                
                                <el-tree
                                    :props="props"
                                    :data="tableData"
                                    show-checkbox
                                    node-key="id" 
                                    ref="treeData"
                                    class="treeTable"
                                    :expand-on-click-node="false"
                                    @check-change="handleCheckChange">
                                    <!-- 使用自定义,需要加slot-scope,返回两个值,node是当前节点指定的对象
                                    data是当前节点的数据 -->
                                    <div class="custom-tree-node"  slot-scope="{ data }" > 
                                        <div class="total_info_box clearfix" v-if="data.name">
                                            <span class="table_info_item" :title="data.name "> {{data.name }} 任务</span>  
                                            <span class="table_info_item" > {{data.alive_number  }}</span>
                                            <span class="table_info_item"> {{data.update_time }}</span>
                                        </div>
                                        <span   class="table_info_node" v-else>
                                            <span class="table_info_item targeturl" :title="data.url">{{data.url}} 
                                                <!-- <el-tooltip class="item" effect="dark" :content="data.url" placement="bottom-start">
                                                    <span>{{data.url}}</span> 
                                                </el-tooltip> -->
                                            </span> 
                                            <span :class="[ 
                                                {'table_info_item riskstyle risk_hight': data.risk_level=='高危' } ,
                                                {'table_info_item riskstyle risk_middle': data.risk_level=='中危' },
                                                {'table_info_item riskstyle risk_low':data.risk_level =='低危' },
                                                {'table_info_item riskstyle risk_nofind':data.risk_level =='未发现' }]"><i></i>{{data.risk_level}}</span> 
 
                                        </span>
                                    </div>
                                </el-tree>
                            </div>
                        </div>
                        
                    </div>  
                    </el-col>
                    <el-col :span="12">
                        <div class="" style="margin-bottom:15px">
                            <div class="grid-title color1">
                                <el-checkbox v-model="taskreportmodel" class="task_checked"></el-checkbox>
                                <label for="">任务报告设置</label>
                                <img src="../../assets/rp.png" alt="">
                            </div>
                            <div class="grid-content h262"> 
                                <!-- <el-form-item
                                    prop ='name'
                                    label=" " > 
                                    <label class="dialog_item_label">报告标题 </label>
                                    <el-input
                                        v-model="reportform.name"
                                        size="small"
                                        style="width:50%"
                                        placeholder="报告标题"
                                        maxlength="50" 
                                        :disabled="!taskreportmodel"
                                    ></el-input>
                                </el-form-item>  -->
                                <el-form-item
                                    label=""  > 
                                    <label class="dialog_item_label">报告模板</label>
                                    <el-select
                                        v-model="reportform.template"
                                        size="small" 
                                        placeholder="请选择"
                                        :disabled="!taskreportmodel"
                                        style=" width: 50%;"  >
                                        <el-option
                                            v-for="(item,index) in templatelist"
                                            :key="index"
                                            :label="item.name"
                                            :value="item.id"
                                        ></el-option>
                                    </el-select>  
                                    <el-link type="primary" class="linkreport" :disabled="!taskreportmodel"  @click="taskReport">自定义报告</el-link> 
                                </el-form-item> 
                                <el-form-item
                                    label=" " > 
                                    <label class="dialog_item_label">安全评价</label>
                                    <el-input
                                        v-model="reportform.remarks"
                                        type="textarea"
                                        size="small"
                                        :rows="5" 
                                        resize="none"
                                        style=" width:70%;vertical-align: top;" 
                                        placeholder="安全评价" 
                                        maxlength="1000"
                                        :disabled="!taskreportmodel"
                                    ></el-input>
                                </el-form-item> 
                            </div>
                        </div>
                        <div class="">
                            <div class="grid-title color2">
                                <el-checkbox v-model="targetreportmodel" @change="changetargetreportmodel"  class="target_checked"></el-checkbox>
                                <label for="">目标报告设置</label>
                                <img src="../../assets/rp.png" alt="">
                            </div>
                            <div class="grid-content h262">
                                <el-form-item
                                    label=" " > 
                                    <label class="dialog_item_label">输出方式</label>
                                    <el-radio-group v-model="reportform.typeradio" style=" display: inline-block;l " :disabled="!targetreportmodel" @change="changeOutType">
                                        <el-radio :label="1">合并输出</el-radio>
                                        <el-radio :label="0">逐个输出</el-radio> 
                                    </el-radio-group>
                                </el-form-item>  
                                <!-- <el-form-item
                                    prop='targetname'
                                    label=" " > 
                                    <label class="dialog_item_label">报告标题 </label>
                                    <el-input
                                        v-model="reportform.targetname"
                                        size="small"
                                        style="width:50%"
                                        placeholder="报告标题"
                                        maxlength="50" 
                                        :disabled="!targetreportmodel || !reportform.typeradio"
                                    ></el-input> 
                                </el-form-item>  -->
                                <el-form-item
                                    label=""  > 
                                    <label class="dialog_item_label">报告模板</label>
                                    <el-select
                                        v-model="reportform.targettemplate"
                                        size="small" 
                                        placeholder="请选择报告模板"
                                        :disabled="!targetreportmodel"
                                        style=" width: 50%"  >
                                        <el-option
                                            v-for="(item,index) in targettemplatelist"
                                            :key="index"
                                            :label="item.name"
                                            :value="item.id"
                                        ></el-option>
                                    </el-select>  
                                    <el-link :disabled="!targetreportmodel" type="primary" class="linkreport" @click="targetReport">自定义报告</el-link> 
                                </el-form-item> 
                                <el-form-item
                                    label=" " prop="output"  > 
                                    <label class="dialog_item_label">输出漏洞</label>
                                    <el-checkbox-group v-model="reportform.output"  style=" display: inline-block;line-height: 20px;" :disabled="!targetreportmodel">
                                        <el-checkbox label="1" name="1">致命漏洞</el-checkbox>
                                        <el-checkbox label="2" name="2">高危漏洞</el-checkbox>
                                        <el-checkbox label="3" name="3">中危漏洞</el-checkbox>
                                        <el-checkbox label="4" name="4">低危漏洞</el-checkbox> 
                                    </el-checkbox-group>
                                </el-form-item>
                                <el-form-item
                                    label=" "  > 
                                    <label class="dialog_item_label">响应报文</label>
                                    <el-radio-group v-model="reportform.is_packet" style=" display: inline-block;line-height: 20px;" :disabled="!targetreportmodel">
                                        <el-radio :label="1">输出</el-radio>
                                        <el-radio :label="0">不输出</el-radio> 
                                    </el-radio-group>
                                </el-form-item>  
                                
                            </div>
                        </div>
                    </el-col> 
                </el-row>  
            </el-form>
        </div>

        <el-dialog 
            title="自定义目标报告" 
            :visible.sync="dialogVisibleTarget"  
            :before-close="cancelTargetReport" 
            class="createReportDialog"
            width='1184px' 
            :close-on-click-modal="false" 
            :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="submitTargetReport">保存</el-button>
                <el-button size="small" @click="cancelTargetReport">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form ref="targetReportform" :model="targetReportform" label-width="0"  >
                    <el-form-item
                        label=" " > 
                        <label class="dialog_item_label">封面标题</label>
                        <el-input
                            v-model="targetReportform.title"
                            size="small"
                            style="width:320px"
                            placeholder="封面标题"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item> 
                    <el-form-item
                        label=" " > 
                        <label class="dialog_item_label">测试单位</label>
                        <el-input
                            v-model="targetReportform.unit"
                            size="small"
                            style="width:320px"
                            placeholder="测试单位"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item>
                    <el-form-item
                        label=" " class="fengmianBg"> 
                        <label class="dialog_item_label">封面背景</label>
                        <img :src="targetReportform.bgImg" alt="" class="form-img" >
                        <el-button type="primary"  class="changeBgBtn" @click="clickupload2()">更换背景</el-button> 
            			<label style="color:rgba(72, 72, 102, 0.64);    line-height: 20px;">建议上传1200*1700px尺寸图片,文件小于10M,仅支持PNG、JPG与JPEG格式。</label>
            			<input type="file"  class="btnUploadID2"  ref='upload' @change="changeBg2($event)" style="display:none" id="input-file-ID" > 
                    </el-form-item> 
                    <el-form-item
                        label=" " > 
                        <label class="dialog_item_label" style="vertical-align: top;">报告目录</label>
                        <el-tree
                            class="catatree"
                            :data="targetreedata"
                            show-checkbox
                            default-expand-all
                            node-key="id"
                            ref="targettree"
                            highlight-current
                            style="display: inline-block;"
                            :props="targetdefaultProps">
                            </el-tree>
                    </el-form-item>
                </el-form>
            </div> 
        </el-dialog>

        <el-dialog 
            title="自定义任务报告" 
            :visible.sync="dialogVisibleTask"  
            :before-close="cancelTaskReport" 
            class="createReportDialog"
            width='1184px' 
            :close-on-click-modal="false" 
            :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="submitTaskReport">保存</el-button>
                <el-button size="small" @click="cancelTaskReport">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form ref="taskReportform" :model="taskReportform" label-width="0"  >
                    <el-form-item
                        label=" " > 
                        <label class="dialog_item_label">封面标题</label>
                        <el-input
                            v-model="taskReportform.title"
                            size="small"
                            style="width:320px"
                            placeholder="封面标题"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item> 
                    <el-form-item
                        label=" " > 
                        <label class="dialog_item_label">测试单位</label>
                        <el-input
                            v-model="taskReportform.unit"
                            size="small"
                            style="width:320px"
                            placeholder="测试单位"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item>
                    <el-form-item
                        label=" " class="fengmianBg"> 
                        <label class="dialog_item_label">封面背景</label>
                        <img :src="taskReportform.bgImg" alt="" class="form-img" >
                        <el-button type="primary"  class="changeBgBtn" @click="clickupload1()">更换背景</el-button> 
            			<label style="color:rgba(72, 72, 102, 0.64);    line-height: 20px;">建议上传1200*1700px尺寸图片,文件小于10M,仅支持PNG、JPG与JPEG格式。</label>
            			<input type="file"  class="btnUploadID1"  ref='upload' @change="changeBg1($event)" style="display:none" id="input-file-ID" > 
                    </el-form-item> 
                    <el-form-item
                        label=" " > 
                        <label class="dialog_item_label" style="vertical-align: top;">报告目录</label>
                        <el-tree
                            class="catatree"
                            :data="tasktreedata"
                            show-checkbox
                            default-expand-all
                            node-key="id"
                            ref="tasktree"
                            highlight-current
                            style="display: inline-block;"
                            :props="defaultProps">
                            </el-tree>
                    </el-form-item>
                </el-form>
            </div> 
        </el-dialog>
  	</div>
</template>
<style  lang="less" scoped>


 
 

/deep/ .el-tree__empty-text{
    line-height: 40px;
    color: rgba(72,72,102,.32);
    font-weight: 500;
    font-size: 13px;
}
/deep/ .el-checkbox{
    margin-right: 16px;
}
 
/deep/ .el-link.el-link--primary:not(.is-disabled):hover{
    color: #4C7AE3;
}
.reportbox{
    padding-bottom: 24px;  
    min-height: calc(100% - 39px);
    box-sizing: border-box; 
}
/deep/ .el-form-item__content{
    line-height: 16px;
}
/deep/ .el-form-item__label{
    line-height: 24px;
}
/deep/ .createReportDialog .fengmianBg .el-form-item__content{
    display: flex;
    align-items: flex-end;
}
.bottomsmall {
    margin-bottom: 12px;
    /deep/ .el-form-item__content{
        line-height: 28px;
    }
}
/deep/ .el-input__icon{
    line-height: 30px;
}
.createreportbox{
	padding: 24px; 
	// background: #fff;
}	
.operationbox{
    padding: 24px;
    margin-bottom: 15px;
    label{
        margin-left: 16px;
        color: #4C7AE3;
        font-size: 13px;
    }
    span{
        margin-left: 4px;
        font-size: 13px;
        color: rgba(72, 72, 102, 0.64);
    }
}
.whitebg{
    background: #fff;
    border-radius:4px;
    box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
}
.grid-title{
    height: 48px; 
    border-radius: 4px 4px 0px 0px; 
    padding-left: 24px;
    box-sizing: border-box;
    /deep/ .el-checkbox{
        margin-right: 20px;
        .el-checkbox__inner{
            background: rgba(255, 255, 255, 0.1200);
             border-color:rgba(255, 255, 255, 0.7) ;
        }
        .el-checkbox__input:hover{
            border-color:rgba(255, 255, 255, 0.7) !important;
        }
        .el-checkbox__input.is-checked .el-checkbox__inner, 
        .el-checkbox__input.is-indeterminate .el-checkbox__inner{
            background: rgba(255, 255, 255, 0.1200) !important;
                border-color: rgba(255, 255, 255, 0.7000) !important;
        }
        .el-checkbox__input.is-focus .el-checkbox__inner{
            border-color:rgba(255, 255, 255, 0.7);
        }
    }
    label{
        display: inline-block;
        color: #fff;
        font-size: 14px;
        margin-top: 14px;
    }
    img{
        float:right;
        width: 120px;
        height: 48px;
    }
}
.is-required{
    margin-right:4px;
    color:#F56C6C;
    font-size: 12px; 
}
.color1{
    background: #468C76;
}
.color2{
    background: #6A468C;
}
.color3{
    background: #4C7AE3;
}
.changeBgBtn{
    margin: 0 16px;
    width: 80px;
    height: 32px;
    text-align: center;
    line-height: 7px;
    padding: 0px;
    border-radius: 2px;
}
.catatree{
    font-size: 13px;
    font-weight: 400;
    color: rgba(72, 72, 102, 0.64);
}
.grid-content{
    background: #fff;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    border-radius: 0px 0px 4px 4px ; 
    padding: 24px;
    box-sizing: border-box;
}
.h262{
    height: 262px;
} 
.linkreport{
    line-height: 20px;
    margin-left: 16px; 
}
.grid-content /deep/ .el-form-item {
    position: relative;
}
.grid-content /deep/ .el-form-item__label {
    position: absolute !important;
    left: 100px !important;
}
.grid-content /deep/ .el-form-item__error{
    margin-left: 112px;
}
.outputRange{
    display: inline-block;
    list-style: none;
    li{ 
        display: inline-block;
        height: 24px;
        width: 68px;
        text-align: center;
        line-height: 24px;
        border-radius: 24px;
        border: 1px solid #E8E8F5;
        margin-right: 16px;
        cursor: pointer;
        i{
            display: inline-block;
            height: 14px;
            width: 14px;
            font-size: 10px;
            text-align: center;
            line-height: 14px;
            background: rgba(72, 72, 102, 0.32);
            color: #fff;
            border-radius: 50%;
            margin-right: 9px;
            cursor: pointer;
        }
        label{
            font-size: 12px;
             cursor: pointer;
            color: rgba(72, 72, 102, 0.64);
        } 
    }
    li:last-child{
         margin-right: 0;
    }
    li.checked{
        background: #4C7AE3;
        color: #fff;
        
        i{
            background: #fff;
            color: #4C7AE3;
            font-weight: 600;
        }
        label{
            color: #fff;
        }
    }
}
/deep/.treeTable .el-tree-node__content{
    height: 40px;
    background-color: #fff; 
    border-bottom: 1px solid #E8E8F5;
}

/deep/ .treeTable .el-tree-node__children .el-tree-node__content{
   
    background-color: #F7F7FB;  
}
.invoice-list {

    /deep/ .el-tree-node__content{
        cursor: auto;
    }
    // border: 1px solid #ebeef5;
    margin-top: 10px;
    
    .treeTable{
        max-height: 370px;
        overflow-y: auto;
    }
    .table_head{
        background-color: #F7F7FB;  
        border-bottom: 1px solid #E8E8F5;
    }
    .table_head > span{
        display: inline-block;
        width:46px; 
    }
    .invoice-header {
        display: inline-block;
        list-style: none;
       
        width: calc(100% - 46px);
        .invoice-item {
            display: inline-block;
            padding: 8px;
            padding-right: 0;
            flex: 1;
            // border-left: 1px solid #ebeef5;
            padding-left: 10px;
            font-size: 13px;
            color: rgba(72, 72, 102, 0.87);
            box-sizing: border-box;
        }
        .task1{ 
            width: 54%;
            padding: 8px 0 8px 10px;
        }
        .task2{
            width: 21%;
        }
        .task3{
            width: 25%;
        }
    }
    .el-tree-node__content {
        background: #f2f2f2;
        height: 40px;
    }
    .el-tree-node__children {
        .el-tree-node__content {
        background: #fff;
        border-bottom: 1px solid #E8E8F5;
        }
    }
    .custom-tree-node {
        width: 100%;
        height: 100%;
        .total_info_box { 
            // border-bottom: 1px solid #E8E8F5;
            line-height: 40px;
            // display: flex;
            span{
                display: inline-block;
                // flex: 1;
                // float: left;
                font-size: 12px;
                // margin: 0 20px;
                i{
                    display: inline-block;
                    margin-right: 3px;
                }
                 padding-left: 10px;
                box-sizing: border-box;
                color: rgba(72, 72, 102, 0.64);
                vertical-align: bottom;
            }
            span:nth-child(1){  
                width: 54%;
                overflow: hidden;
                text-overflow:ellipsis;
                white-space: nowrap;
            }
            span:nth-child(2){
                width: 21%;
            }
            span:nth-child(3){
                width: 25%;
            }
        }
        .table_info_node {
            // display: flex;
            // height: 100%;
            height: 40px;
            line-height: 40px;
            display: block;
            .table_info_item {
                // flex: 1;
                display: inline-block;
                // height: 100%;
                height: 40px;
                // border-bottom: 1px solid #E8E8F5;
                padding-left: 10px;
                line-height: 40px;
                font-size: 13px;
                // width: 25%;
                vertical-align: bottom;
                color: rgba(72, 72, 102, 0.64);
                box-sizing: border-box;
                cursor: auto;
            }
            .targeturl{
                width: 75%;
                padding-right:10px ;
                overflow: hidden;
                text-overflow:ellipsis;
                white-space: nowrap;
            }
            .riskstyle{
                 width: 25%;
                 padding-left: 6px;
            }
             
        }
    }
}
.el-checkbox{
    color: rgba(72, 72, 102, 0.64);
}
.el-radio{
    color: rgba(72, 72, 102, 0.64);
}
.form-img{
    width: 100px;
    height: 140px;
}
 
</style>
<script>   
import xzbutton from "@/components/XzButton.vue"; 
import API from '@/api/report.js'
export default({
    name:'createreport', 
    components: {
    	xzbutton, 
  	},
    data(){  
    	return{  
            taskreportmodel:true,
            targetreportmodel:false,
            task_id:this.$route.query.task_id, //任务ID
            reportform:{
                num:'5',
                name:'自动化渗透测试综合报告',
                targetname:'自动化渗透测试多目标报告',
                typeradio:1,
                reportstyle:['html'],
                output:['1','2','3'],
                is_packet:0,
                targettemplate:'',
                template:'',
            },
            reportrules:{
                name: [
                    { required: true, message: '任务报告标题不能为空', trigger: 'blur' }, 
                ], 
                targetname:[
                     { required: true, message: '目标报告标题不能为空', trigger: 'blur' }, 
                ],
                output : [
                        { type: 'array', required: true, message: '请至少选择一个输出漏洞', trigger: 'change' }
                ],
              
            },
            // outputRangelist:[[4,'高危目标',1],[3,'中危目标',0],[2,'低危目标',0],[1,'安全目标',0]], //输出范围
            outputRangelist:[
                {
                    id:1,
                    name:'高危',
                    isChecked:1,
                },
                 {
                    id:2,
                    name:'中危',
                    isChecked:1,
                },
                 {
                    id:3,
                    name:'低危',
                    isChecked:1,
                },
                 {
                    id:4,
                    name:'安全',
                    isChecked:0,
                },
            ],
            templatelist:[],
            targettemplatelist:[],
            checkList:[],
            radio:1, 
            reportstyle:[],
            props: {
                label: 'id', // 需要指定的节点渲染对象属性
                children: 'target_list' // 指定的子级
            },
            tableData: [], // tree组件渲染的数据
            dialogVisibleTarget:false,
            dialogVisibleTask:false,
            targetReportform:{
                name:'自定义目标报告',
                title:'',
                unit:'',
                cover:'',
                // bgImg:require('../../assets/templatebgsm.png'),
                bgImg:'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAANAAAAEjCAIAAADMtxggAAAACXBIWXMAAAuJAAALiQE3ycutAAAF+mlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4gPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iQWRvYmUgWE1QIENvcmUgNS42LWMxNDUgNzkuMTYzNDk5LCAyMDE4LzA4LzEzLTE2OjQwOjIyICAgICAgICAiPiA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtbG5zOnhtcE1NPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvbW0vIiB4bWxuczpzdEV2dD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL3NUeXBlL1Jlc291cmNlRXZlbnQjIiB4bWxuczpkYz0iaHR0cDovL3B1cmwub3JnL2RjL2VsZW1lbnRzLzEuMS8iIHhtbG5zOnBob3Rvc2hvcD0iaHR0cDovL25zLmFkb2JlLmNvbS9waG90b3Nob3AvMS4wLyIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgQ0MgMjAxOSAoV2luZG93cykiIHhtcDpDcmVhdGVEYXRlPSIyMDIxLTA4LTA0VDE2OjIwOjA2KzA4OjAwIiB4bXA6TWV0YWRhdGFEYXRlPSIyMDIxLTA4LTA0VDE2OjIwOjA2KzA4OjAwIiB4bXA6TW9kaWZ5RGF0ZT0iMjAyMS0wOC0wNFQxNjoyMDowNiswODowMCIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDo0ZGNiMGUxNy00YTBmLWNjNGMtYjYzYi0xOTkxNDQyOGM5OGIiIHhtcE1NOkRvY3VtZW50SUQ9ImFkb2JlOmRvY2lkOnBob3Rvc2hvcDo3MDc3ODBjNy1jMjRjLTIyNGMtYTZiYS1iZDc1NDdjMDZhYzUiIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDpkODgyNDJhMi00NTczLTgyNDMtOGQ3NS0zMjBlNWE5NzRjNWMiIGRjOmZvcm1hdD0iaW1hZ2UvcG5nIiBwaG90b3Nob3A6Q29sb3JNb2RlPSIzIiBwaG90b3Nob3A6SUNDUHJvZmlsZT0ic1JHQiBJRUM2MTk2Ni0yLjEiPiA8eG1wTU06SGlzdG9yeT4gPHJkZjpTZXE+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJjcmVhdGVkIiBzdEV2dDppbnN0YW5jZUlEPSJ4bXAuaWlkOmQ4ODI0MmEyLTQ1NzMtODI0My04ZDc1LTMyMGU1YTk3NGM1YyIgc3RFdnQ6d2hlbj0iMjAyMS0wOC0wNFQxNjoyMDowNiswODowMCIgc3RFdnQ6c29mdHdhcmVBZ2VudD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTkgKFdpbmRvd3MpIi8+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJzYXZlZCIgc3RFdnQ6aW5zdGFuY2VJRD0ieG1wLmlpZDo0ZGNiMGUxNy00YTBmLWNjNGMtYjYzYi0xOTkxNDQyOGM5OGIiIHN0RXZ0OndoZW49IjIwMjEtMDgtMDRUMTY6MjA6MDYrMDg6MDAiIHN0RXZ0OnNvZnR3YXJlQWdlbnQ9IkFkb2JlIFBob3Rvc2hvcCBDQyAyMDE5IChXaW5kb3dzKSIgc3RFdnQ6Y2hhbmdlZD0iLyIvPiA8L3JkZjpTZXE+IDwveG1wTU06SGlzdG9yeT4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz4h6fgpAAAgZklEQVR4nO2dSW9cWZbf/+fe+16MjCBFioPmTClTyqxuZLlQWWUXjHbb6F7YsL3wwoaHhQED/gT2wlvDGy/8DexFG27AUxnVBtyobvRQ5aou2OUasjqHzlRmKjWLEiVOMbzx3uPFizmCgyjGU5A8Pwhk6L0XL4LBH8+999yJmm2GIOSFEd2EPFGv+w0IZwsRTsgVEU7IFRFOyBURTsgVEU7IFRFOyBURTsgVEU7IFRFOyBURTsgVEU7IFRFOyBURTsgVEU7IFRFO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',
            },
            taskReportform:{
                name:'自定义任务报告',
                title:'',
                unit:'',
                cover:'',
                // bgImg:require('../../assets/templatebgsm.png'),
                bgImg:'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAANAAAAEjCAIAAADMtxggAAAACXBIWXMAAAuJAAALiQE3ycutAAAF+mlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4gPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iQWRvYmUgWE1QIENvcmUgNS42LWMxNDUgNzkuMTYzNDk5LCAyMDE4LzA4LzEzLTE2OjQwOjIyICAgICAgICAiPiA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtbG5zOnhtcE1NPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvbW0vIiB4bWxuczpzdEV2dD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL3NUeXBlL1Jlc291cmNlRXZlbnQjIiB4bWxuczpkYz0iaHR0cDovL3B1cmwub3JnL2RjL2VsZW1lbnRzLzEuMS8iIHhtbG5zOnBob3Rvc2hvcD0iaHR0cDovL25zLmFkb2JlLmNvbS9waG90b3Nob3AvMS4wLyIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgQ0MgMjAxOSAoV2luZG93cykiIHhtcDpDcmVhdGVEYXRlPSIyMDIxLTA4LTA0VDE2OjIwOjA2KzA4OjAwIiB4bXA6TWV0YWRhdGFEYXRlPSIyMDIxLTA4LTA0VDE2OjIwOjA2KzA4OjAwIiB4bXA6TW9kaWZ5RGF0ZT0iMjAyMS0wOC0wNFQxNjoyMDowNiswODowMCIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDo0ZGNiMGUxNy00YTBmLWNjNGMtYjYzYi0xOTkxNDQyOGM5OGIiIHhtcE1NOkRvY3VtZW50SUQ9ImFkb2JlOmRvY2lkOnBob3Rvc2hvcDo3MDc3ODBjNy1jMjRjLTIyNGMtYTZiYS1iZDc1NDdjMDZhYzUiIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDpkODgyNDJhMi00NTczLTgyNDMtOGQ3NS0zMjBlNWE5NzRjNWMiIGRjOmZvcm1hdD0iaW1hZ2UvcG5nIiBwaG90b3Nob3A6Q29sb3JNb2RlPSIzIiBwaG90b3Nob3A6SUNDUHJvZmlsZT0ic1JHQiBJRUM2MTk2Ni0yLjEiPiA8eG1wTU06SGlzdG9yeT4gPHJkZjpTZXE+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJjcmVhdGVkIiBzdEV2dDppbnN0YW5jZUlEPSJ4bXAuaWlkOmQ4ODI0MmEyLTQ1NzMtODI0My04ZDc1LTMyMGU1YTk3NGM1YyIgc3RFdnQ6d2hlbj0iMjAyMS0wOC0wNFQxNjoyMDowNiswODowMCIgc3RFdnQ6c29mdHdhcmVBZ2VudD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTkgKFdpbmRvd3MpIi8+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJzYXZlZCIgc3RFdnQ6aW5zdGFuY2VJRD0ieG1wLmlpZDo0ZGNiMGUxNy00YTBmLWNjNGMtYjYzYi0xOTkxNDQyOGM5OGIiIHN0RXZ0OndoZW49IjIwMjEtMDgtMDRUMTY6MjA6MDYrMDg6MDAiIHN0RXZ0OnNvZnR3YXJlQWdlbnQ9IkFkb2JlIFBob3Rvc2hvcCBDQyAyMDE5IChXaW5kb3dzKSIgc3RFdnQ6Y2hhbmdlZD0iLyIvPiA8L3JkZjpTZXE+IDwveG1wTU06SGlzdG9yeT4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz4h6fgpAAAgZklEQVR4nO2dSW9cWZbf/+fe+16MjCBFioPmTClTyqxuZLlQWWUXjHbb6F7YsL3wwoaHhQED/gT2wlvDGy/8DexFG27AUxnVBtyobvRQ5aou2OUasjqHzlRmKjWLEiVOMbzx3uPFizmCgyjGU5A8Pwhk6L0XL4LBH8+999yJmm2GIOSFEd2EPFGv+w0IZwsRTsgVEU7IFRFOyBURTsgVEU7IFRFOyBURTsgVEU7IFRFOyBURTsgVEU7IFRFOyBURTsgVEU7IFRFO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',
            },
            tasktreedata:[
            {
                    id:1,
                    label:'任务概述',
                    key:'task_info',
                    level:1
                },
                {
                    id:2,
                    label:'信息统计',
                    key:'info_statistics',
                    level:1,
                    children:[
                        { id:3,label:'目标统计',key:'target'},
                        { id:4,label:'漏洞风险统计',key:'vuln_risk'},
                        { id:5,label:'漏洞类型统计',key:'vuln_type'},
                        { id:6,label:'服务类型统计',key:'service'}
                    ]
                },{
                    id:7,
                    label:'目标详情',
                    key:'target_detail',
                    level:1
                },{
                    id:8,
                    label:'漏洞详情',
                    key:'vuln_info', //is_open:0/1
                    level:1,
                    children:[
                        { id:9,label:'已发现漏洞',key:'find_vuln'},
                        // { id:12,label:'已修复漏洞',key:'repair_vuln'},
                    ],
                },{
                    id:10,
                    label:'脆弱账号',
                    key:'fragile_account',
                    level:1
                }
            ],
            targetreedata:[
            {
                    id:1,
                    label:'报告摘要', 
                    key:'task_info',
                },
                {
                    id:2,
                    label:'测试范围',
                    key:'coverage'
                },
                {
                    id:3,
                    label:'资产信息',
                    key:'assets'
                },
                {
                    id:4,
                    label:'漏洞信息',
                    key:'vuln_info', //is_open:0/1
                    children:[
                        {
                            id:5,
                            label:'检测存在漏洞',
                            key:'find_vuln'
                        }
                    ],
                },
                {
                    id:7,
                    label:'攻击面',
                    key:'attack'
                },
            ],
            targetdefaultProps:{
                children: 'children',
                label: 'label'
            },
            defaultProps: {
                children: 'children',
                label: 'label'
            },
    	}
    }, 
    created:function(){ 
        this.$store.state.activefirstMenu="/createreport"; 
       

        // this.getReportTask(); //任务列表
        // this.getTemplateList(1); //报告模板 参数，1：任务，2:目标
        // this.getTemplateList(2);
    },
    mounted:function(){  
    },
    methods:{  
        getReportTask(){ 
            let that  = this;
            //http://192.168.0.171:8080/report/report/task/target/?risk_level=1&page_size=200
            var risk_level_arr = [];

            for(var i=0;i<this.outputRangelist.length;i++){
                if(this.outputRangelist[i].isChecked == 1){
                    risk_level_arr.push(this.outputRangelist[i].id)
                }
            }
            let params = {
                page_size:this.reportform.num,
                risk_level:risk_level_arr.join(','),
                task_id:this.task_id
            }
            API.getTargetList(params).then(res => {
                if(res.success){ 
                    this.tableData = res.data;
                    // this.$nextTick(() => {
                    //     that.$refs.treeData.setCheckedNodes(res.data); 
                    // });
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            }).catch(data=>{ 
                console.log(data);
            });
        },
        getTemplateList(_type){ //获得模板列表
            let params = {
                type:_type
            }
            API.getTemplateList(params).then(res => {
                if(res.success){ 
                    if(_type == 1){ //任务
                        this.templatelist = res.data;
                        for (var i = 0; i < this.templatelist.length; i++) {
                            if (this.templatelist[i].default) {
                                this.reportform.template = this.templatelist[i].id;  
                                break;
                            }
                        }
                    }else{//目标
                        this.targettemplatelist = res.data;
                        for (var i = 0; i < this.targettemplatelist.length; i++) {
                            if (this.targettemplatelist[i].default) {
                                this.reportform.targettemplate = this.targettemplatelist[i].id; 
                                break;
                            }
                        }
                    }       
                } else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            })
            .catch(data=>{ });
        },
        // getSupplierPayInvoice () {
        //     this.tableData = [{
        //         id: 13, 
        //         name:'192.168.0.111', 
        //         alive_number:1,
        //         update_time: "2021-04-20 10:52", 
        //         risk_level: "未发现", 
        //         span:true,
        //         target_list: [ 
        //             {
        //                 url: "192.168.0.63", 
        //                 risk_level: "高危", 
        //                 id: 16617
        //             }, {
        //                 url: "192.168.0.163", 
        //                 risk_level: "低危", 
        //                 id: 16618
        //             }
        //         ]
        //     }, {
        //         id: 14, 
        //         alive_number:2,
        //          span:true,
        //         name:'192.168.1.111-192..168.1.222', 
        //         update_time: "2021-04-20 10:52",  
        //         risk_level: "未发现", 
        //         target_list: [
                
        //         ]
        //     }]

        // },
        changeTaskNum(){ //选择任务条目
            this.getReportTask();
        },
        changeReportForm(){ //选择报告格式
            console.log(this.reportform.reportstyle);
        },
        changeOutType(e){ //选择输出方式
            if(!e){
                this.reportform.targetname = '地址+渗透测试报告';
            }else if(e){
                this.reportform.targetname = '自动化渗透测试多目标报告';
            }
        },
    	handleCheckChange(){
            let selected = this.$refs.treeData.getCheckedNodes()
            // console.log(33, selected)
        },
        targetReport(){
            this.fnGetCustomTemplateConfig(2);
            // this.dialogVisibleTarget = true;
        },
        submitTargetReport(){ //保存自定义目标报告
            var submitjson = {
                name:this.targetReportform.name,
                title:this.targetReportform.title,
                unit:this.targetReportform.unit,
                cover:this.targetReportform.bgImg,
                type:2,
            };
            let treedata = this.$refs.targettree.getCheckedNodes();
            console.log(treedata);
            for(var i=0;i<treedata.length;i++){
                let item = treedata[i];
                if(item.children){
                    let key = item.key;
                    submitjson[key] =  1; 
                }else{
                    let key = item.key;
                    submitjson[key] = 1;
                }
                
            }
            // console.log(submitjson);
            API.submitTargetReport(submitjson).then(res => {
                if(res.success){ 
                    this.$message({
                        message:'自定义目标报告添加成功',
                        type: 'success'
                    });
                    this.dialogVisibleTarget = false; 
                    this.getTemplateList(2); 
                    setTimeout(() => {
                        this.reportform.targettemplate = res.id; 
                    }, 1000);
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            })
            .catch(error=>{ });
        },
        cancelTargetReport(){
            this.dialogVisibleTarget = false;
        },
        clickupload(){},
        changeBg(){},
        taskReport(){
            this.fnGetCustomTemplateConfig(1);
            // this.dialogVisibleTask = true;
        },
        cancelTaskReport(){
            this.dialogVisibleTask = false;
        },
        submitTaskReport(){ //保存自定义任务报告
            var submitjson = {
                name:this.taskReportform.name,
                title:this.taskReportform.title,
                unit:this.taskReportform.unit,
                cover:this.taskReportform.bgImg,
                type:1,
            };
            let treedata = this.$refs.tasktree.getCheckedNodes();
            // console.log(treedata);
            for(var i=0;i<treedata.length;i++){
                let item = treedata[i];
                if(item.children){
                    let key = item.key;
                    submitjson[key] =  1; 
                }else{
                    let key = item.key;
                    submitjson[key] = 1;
                }
                
            }
            // console.log(submitjson);
            API.submitTargetReport(submitjson).then(res => {
                let that = this;  
                if(res.success){ 
                    this.$message({
                        message:'自定义任务报告添加成功',
                        type: 'success'
                    });
                    this.dialogVisibleTask = false;  
                    this.getTemplateList(1);
                    setTimeout(() => {
                        that.reportform.template = res.id; 
                    }, 1000);
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            })
            .catch(error=>{ });

            // let newtreedata=[];
            // for(var i=0;i<treedata.length;i++){
            //     if(treedata[i].level){
            //         newtreedata.push(treedata[i]);
            //     }
            // }
            // console.log(newtreedata)
            // for(var i=0;i<newtreedata.length;i++){
            //     let item = newtreedata[i];
            //     if(item.children){ 
            //         let key = item.key;
            //         submitjson[key] = {
            //             is_open:1
            //         }; 
            //         for(var j=0;j<item.children.length;j++){
            //             var _j = item.children[j]; 
            //             submitjson[key][_j.key] = 1;
            //         }
            //     }else{
            //         let key = item.key;
            //         submitjson[key] = 1;
            //     }
            // }
            // console.log(submitjson);


        },
        fnGetCustomTemplateConfig(type){ //获得自定义模板配置
            var name = type == 1? '自定义任务报告':'自定义目标报告';
            let params = {
                template_name:name
            }
            API.fnGetCustomTemplateConfig(params).then(res => {
                if(res.success){  
                    if(type == 1){ //任务
                        this.dialogVisibleTask = true;  
                        this.taskReportform.title = res.data.title;
                        this.taskReportform.unit = res.data.detect_unit; 
                        this.taskReportform.bgImg = res.data.cover_data ? res.data.cover_data : this.taskReportform.bgImg;
                        // console.log(this.taskReportfrom);
                        // this.$refs.tasktree.getCheckedNodes();

                        var mulu =  res.data;
 
                        var dataarr=[];
                        for( var key in mulu ){
                            if(key == 'id' && key == 'name'){
                                continue;
                            }
                            var value = mulu[key]; 
                            for(var i=0;i<this.tasktreedata.length;i++){
                                var item = this.tasktreedata[i];
                                if(item.key == key && value==1)
                                {
                                    dataarr.push({
                                        id:item.id,
                                        label:item.label,
                                        key:item.key,
                                    })
                                }
                                if(item.children){
                                    for(var j = 0; j<item.children.length;j++){
                                        var child = item.children[j];
                                        if(child.key == key && value==1){
                                             dataarr.push({
                                                id:child.id,
                                                label:child.label,
                                                key:child.key,
                                            })
                                        }
                                    }
                                }
                            } 
                        } 
                        var that = this;
                        this.$nextTick(() => {
                            that.$refs.tasktree.setCheckedNodes( dataarr);
                        }); 

                    }else{
                        this.dialogVisibleTarget = true;  
                        this.targetReportform.title = res.data.title;
                        this.targetReportform.unit = res.data.detect_unit;
                        this.targetReportform.bgImg = res.data.cover_data ? res.data.cover_data : this.targetReportform.bgImg;

                        var mulu =  res.data;
 
                        var dataarr=[];
                        for( var key in mulu ){
                            if(key == 'id' && key == 'name'){
                                continue;
                            }
                            var value = mulu[key]; 
                            for(var i=0;i<this.targetreedata.length;i++){
                                var item = this.targetreedata[i];
                                if(item.key == key && value==1)
                                {
                                    dataarr.push({
                                        id:item.id,
                                        label:item.label,
                                        key:item.key,
                                    })
                                }
                                if(item.children){
                                    for(var j = 0; j<item.children.length;j++){
                                        var child = item.children[j];
                                        if(child.key == key && value==1){
                                             dataarr.push({
                                                id:child.id,
                                                label:child.label,
                                                key:child.key,
                                            })
                                        }
                                    }
                                }
                            } 
                        }
 
                        var that = this;
                        this.$nextTick(() => {
                            that.$refs.targettree.setCheckedNodes( dataarr);
                        }); 
                    } 
                }
				else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            })
            .catch(data=>{ });
        },
        changeChecked(i){  //选择输出范围
            if(this.outputRangelist[i].isChecked == 0){ 
                        this.$set(this.outputRangelist[i],'isChecked',1) 
            }else{ 
                this.$set(this.outputRangelist[i],'isChecked',0) 
            }
            // console.log(this.outputRangelist) 

            this.getReportTask();
        },
        changetargetreportmodel(){
            // console.log(this.targetreportmodel);
            // let that = this;
            // if(this.targetreportmodel){ //选择了目标设置
            //     this.$nextTick(() => {
            //         that.reportrules.output = [
            //             { type: 'array', required: true, message: '请至少选择一个输出漏洞', trigger: 'change' }
            //         ]

            //         console.log(that.reportrules)
            //     });   
            // }else{
                // delete that.reportrules.output;
                // this.reportform.output = ['1'];

            if(!this.targetreportmodel){
                let _field = this.$refs.reportform.fields ; /*当然，你可以打印一下fields*/
                _field.map(i => {
                    if(i.prop === 'output'){  //通过prop属性值相同来判断是哪个输入框，比如：要移除prop为'user'
                        i.resetField() 
                        return false ;
                    } 
                })
            }
                
            // }
            //  console.log(that.reportrules)

        },
        submitCreateReport(){ //生成报告
            var submitjson = {
                form:this.reportform.reportstyle.join(','), 
            };
            console.log(submitjson);
            // 任务设置
            let tasksetting={
                is_generate:this.taskreportmodel,
                title:this.reportform.name,
                report_template:this.reportform.template,
                estimate:this.reportform.remarks
            };
            
            //目标设置
            if(!this.targetreportmodel){ //选择了目标设置
                delete this.reportrules.output;
            }


            let targetsetting={
                is_generate:this.targetreportmodel,
                report_template:this.reportform.targettemplate,
                title:this.reportform.targetname,
                output:this.reportform.output.join(','),
                is_packet:this.reportform.is_packet,
                is_merge:!!this.reportform.typeradio
            };
 
            // 输出范围
            let ranglist = [];
            let seletedranglist = this.$refs.treeData.getCheckedNodes();
            let seletedranglist2 = this.$refs.treeData.getHalfCheckedNodes();
            this.tableData.forEach(item => {
                if (item.is_task) {
                    let parentNode = {
                        task_id: '',
                        target_ids: []
                    }
                    let count = 0
                    // 过滤被选中的父类
                    seletedranglist.forEach(item2 => {
                        if (item2.is_task && item2.id === item.id) {
                            count++
                            parentNode.task_id = item.id
                        }
                    })
                    // 从半选中筛查父类有没有选中
                    if (count === 0) {
                        seletedranglist2.forEach(item2 => {
                            if (item2.is_task && item2.id === item.id) {
                                parentNode.task_id = item.id
                            }
                        })
                    }
                    // 被选中的子类
                    item.target_list.forEach(item3 => {
                        seletedranglist.forEach(item2 => {
                            if (!item2.is_task && item2.id === item3.id) {
                                parentNode.target_ids.push(item3.id)
                            }
                        })
                    })
                    if (parentNode.task_id) {
                        ranglist.push(parentNode)
                    }
                }
            })
            // console.log('ranglist', ranglist)
            // console.log('this.$refs.treeData', this.$refs.treeData.getHalfCheckedNodes())
            // console.log(seletedranglist);
            // let ranglist = [];
            // for(var i=0;i<seletedranglist.length;i++){
            //     var item = seletedranglist[i];
            //     let task = {};
            //     debugger
            //     if(item.is_task){ //任务
            //         task.task_id = item.id;
            //         // let target= [];
            //         ranglist.push(task);
            //         // if(seletedranglist[i].target_list){ //有子集
            //         //     // for(var j = 0;j<seletedranglist[i].target_list.length;j++){ 
            //         //     //     target.push(seletedranglist[i].target_list[j].id)
            //         //     // }
            //         //     // task.target_id = target;
            //         //     ranglist.push(task);
            //         // }else{
            //         //     // task.target_id = seletedranglist[i].id;
            //         //     // ranglist.push(task);
            //         // }
            //     }else{
            //             task.target_id = seletedranglist[i].id;
            //             ranglist.push(task);
            //     } 
            // }
            
            // console.log(ranglist);
            
             submitjson.range = JSON.stringify(ranglist);
            submitjson.task_setting = JSON.stringify(tasksetting);
            submitjson.target_setting = JSON.stringify(targetsetting);

            // console.log(submitjson);

            // return;
            this.$refs.reportform.validate((valid) => {
                if (valid) {
                    API.submitCreateReport(submitjson).then(res => {
                        if (res.success){ 
                            this.$message({
                                message:'生成报告成功',
                                type: 'success'
                            });
                            
                        }else{
                            this.$message({
                                message: res.error,
                                type: 'error'
                            });
                        }
                    })
                    .catch(error=>{ });
                }
            });
        },
        reset(){ //重置
            // this.reportform = {
            //     num:'5',
            //     name:'自动化渗透测试综合报告',
            //     targetname:'自动化渗透测试多目标报告',
            //     typeradio:1,
            //     reportstyle:['html'],
            //     output:['1','2','3'],
            //     is_packet: 0,
            //     remarks: '' // 安全评价
            // };
            this.reportform.num = '5'
            this.reportform.name = '自动化渗透测试综合报告'
            this.reportform.targetname = '自动化渗透测试多目标报告'
            this.reportform.typeradio = 1
            this.reportform.reportstyle = ['html'];
            this.reportform.output = ['1','2','3'];
            this.reportform.is_packet = 0;
            this.reportform.remarks = ''
            this.getReportTask();
            this.getTemplateList(1); //参数，1：任务，2:目标
            this.getTemplateList(2);
            this.targetreportmodel = false
            this.$set(this.outputRangelist[0],'isChecked',1) 
            this.$set(this.outputRangelist[1],'isChecked',1) 
            this.$set(this.outputRangelist[2],'isChecked',1) 
            this.$set(this.outputRangelist[3],'isChecked',0) 
            // this.reportform.template = this.templatelist[0].id; 
            // this.reportform.targettemplate = this.targettemplatelist[0].id;
        },
        clickupload1(){  
            document.querySelector('.btnUploadID1').click();
        },
        changeBg1(e){
            // console.log(e)
            let that = this
            // let deviceFile = e.target.files;  
            // for(let i=0;i<deviceFile.length;i++){ 
            //     console.log(deviceFile)
            //     this.addform.file = deviceFile[i];
            //     this.addform.filename = deviceFile[i].name 
            // } 
            var id = e.srcElement.id
            var file = e.target.files[0]
            var reads = new FileReader()
            reads.readAsDataURL(file)
            reads.onload = function (e) {
                let fileMaxSize = 1024 * 10;//10M
                let size = file.size/1024;
                if (size > fileMaxSize) {
                    that.$message.error('文件大小不能大于10M！');
                    return false;
                }
                that.taskReportform.bgImg = this.result
            }
            
        },
        
        clickupload2(){  
            document.querySelector('.btnUploadID2').click();
        },
        changeBg2(e){
            // console.log(e)
            let that = this
            // let deviceFile = e.target.files;  
            // for(let i=0;i<deviceFile.length;i++){ 
            //     console.log(deviceFile)
            //     this.addform.file = deviceFile[i];
            //     this.addform.filename = deviceFile[i].name 
            // } 
            var id = e.srcElement.id
            var file = e.target.files[0]
            var reads = new FileReader()
            reads.readAsDataURL(file)
            reads.onload = function (e) {
                let fileMaxSize = 1024 * 10;//10M
                let size = file.size/1024;
                if (size > fileMaxSize) {
                    that.$message.error('文件大小不能大于10M！');
                    return false;
                }
                // that.taskfrom.bgImg = this.result
                that.targetReportform.bgImg = this.result;
            }
            
        },
       
        
    }
})
 
</script>
