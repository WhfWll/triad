<!-- DNS log............................................................................................................. -->
<template>
  	<div > 
        <div class="auxiliarytool context_box_bg"> 
            <div class="basic">
                <div><strong>介绍：</strong>DNS Log功能用于辅助无回显漏洞盲测</div>
                <div><strong>使用帮助：</strong>默认使用系统的DNS Log盲测模块，也可以选择用户自定义的DNS Log盲测模块，使用自定义DNS Log时，系统将不能判断盲测漏洞是否存在</div>
            </div>
            <div style="margin-top:20px" class="search-box">
                <!-- <xzbutton 
                    type="primary" 
                    @click="clearall" 
                    :disabled="tableData.length == 0"
                    size="small">DNS配置</xzbutton>   -->
                <delbutton 
                    :width="170"  
                    :disabled="!multipleSelection.length" style="margin-left: 8px;"></delbutton> 
                <div class="serach-condition">
                    <div class="search-text">
                        <el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearch"  v-model="search_field" class="input-with-select"  size="small" clearable > </el-input>
                        <xzbutton 
                        type="primary" 
                        @click="handlesearch" 
                        :disabled="false" 
                        size="small"  >搜索</xzbutton>  
                    </div>
                    <div >
                        <xzbutton 
                        type="primary" 
                        @click="handleReset" 
                        :disabled="false" 
                        size="small"  >重置</xzbutton>  
                    </div>  
            
                </div>
            </div>
            <div>
                <el-table
                    @selection-change="handleSelectionChange"
                    :data="dnstableData" 
                    style="width: 100%">
                    <el-table-column
                        prop="url"
                        label="域名"
                        >
                    </el-table-column>
                    <el-table-column
                        prop="headers"
                        label="DNS类型">
                    </el-table-column>
                    <el-table-column
                        prop="post_data"
                        label="远端IP"
                        >
                    </el-table-column>
                    <el-table-column
                        prop="req_time"
                        label="时间戳" 
                        >
                    </el-table-column>
                </el-table>
                <el-pagination
                    :page-size="10" 
                    background
                    layout="total,  prev, pager, next, sizes,jumper"
                    :total="totalpage"
                    :current-page="currentpage"
                    @current-change = "currentchange"
                    >
                </el-pagination>
            </div>
<!--IP域名绑定............................................................................................................. -->
                <!-- <el-tab-pane label="IP域名绑定" name="IP域名绑定">
                    <div class="basic">
                        <div><strong>介绍：</strong>渗透资源模块是渗透过程中，被渗透目标需要从远程下载的文件</div>
                        <div><strong>使用帮助：</strong>编写渗透脚本时，复制需要执行文件的链接到脚本中，即可在执行脚本时，被测目标能够从链接地址下载渗透资源文件</div>
                    </div>
                    <div style="margin-top:20px" class="search-box">
                        <xzbutton 
                            type="primary" 
                            @click="clearall" 
                            :disabled="tableData.length == 0"
                            size="small">新增绑定</xzbutton>  
                        <delbutton 
                            :width="170"  
                            :disabled="!multipleSelection.length"></delbutton> 
                        <div class="serach-condition">
                            <div class="search-text">
                                <el-input placeholder="请输入关键字"  v-model="search_field" class="input-with-select"  size="small" clearable > </el-input>
                                <xzbutton 
                                type="primary" 
                                @click="handlesearch" 
                                :disabled="false" 
                                size="small"  >搜索</xzbutton>  
                            </div>
                            <div >
                                <xzbutton 
                                type="primary" 
                                @click="handleReset" 
                                :disabled="false" 
                                size="small"  >重置</xzbutton>  
                            </div>  
					
                        </div>
                    </div>
                    <div>
                        <el-table
                        @selection-change="handleSelectionChange"
                        :data="dnstableData" 
                        style="width: 100%">
                        <el-table-column
                            prop="url"
                            label="IP"
                            >
                        </el-table-column>
                        <el-table-column
                            prop="headers"
                            label="域名">
                        </el-table-column>
                        </el-table>
                        <el-pagination
                            :page-size="10" 
                            background
                            layout="total,  prev, pager, next, sizes,jumper"
                            :total="totalpage"
                            :current-page="currentpage"
                            @current-change = "currentchange"
                            >
                        </el-pagination>
                    </div>
                </el-tab-pane> -->
            <!-- </el-tabs>  -->
            
	  	</div>
	  	<!-- <el-dialog title="清空http记录" :visible.sync="dialogFormVisible"   width="640px" :before-close="cancelform" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div class="dialogtxt" >此操作将清空工具下的所有HTTP记录，确认执行吗？</div>
            
        </el-dialog> -->
  	</div>
</template>
<style lang="less" scoped>
   /deep/  .el-table{
        tr{
            height: 57px;
        }
        .el-button--small{
            padding: 0;
        }
    }
   /deep/ .el-dialog{
        height: 192px !important;
    }
    .upload_dialog /deep/ .el-dialog{
        height: calc(100% - 96px) !important;
    }
   
    .dialogtxt{
        text-align: center;
        margin-top: 55px;
    }
    /deep/ .el-table td:not(.el-table-column--selection):first-child .cell, 
    /deep/ .el-table th:not(.el-table-column--selection):first-child .cell{
        padding-left: 32px !important;
    }
    /deep/ .el-tabs__item{
        height: 48px;
        line-height: 48px;
        padding: 0 24px;
    }
    /deep/ .el-tabs__item.is-active{
        color: #4C7AE3;
        font-weight: 500;
    }
    /deep/ .el-tabs__nav-wrap{
        padding: 0 24px; 
    }
    /deep/ .el-tabs__nav-wrap::after{
        background: #E8E8F5;
        height: 1px;
    }
    /deep/ .el-tabs__header{
        margin: 0 0 24px;
    }
    .auxiliarytool{ 
       background: #fff;
        min-height: calc(100% - 39px);
        box-sizing: border-box;
        // padding: 24px;
        
         /deep/ .el-tabs__header{
            margin:0;
        }
        /deep/ .el-tabs__content{
            padding: 24px ; 
            background: #fff;
            // min-height: 680px;
        }
    }	
    .auxiliarytool .el-tabs__item.is-top.is-active{
        background: #fff;
    }
    .auxiliarytool .el-tabs__header{
        margin:0;
    }
    .auxiliarytool .el-tabs__content{
        padding:20px 30px; 
        background: #fff;
        min-height: 500px;
        box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
        border-radius: 4px;
    }
    .auxiliarytool .tabsbox{
        background: #fff;
    }
    .auxiliarytool .basic{
        background: #F7F7FB !important;;
        border:1px solid #e2e5ed;
        border-left: 2px solid #4c7ae3;
        padding: 5px 10px;
    }
    .auxiliarytool .basic > div{
        color: rgba(72, 72, 102, 0.64);
        font-weight: 500;
        margin:10px 8px;
        font-size: 13px;
    }
    .auxiliarytool .basic > div strong{
        display: inline-block;
        width: 80px;
        color: #4c7ae3;
        
    }
    .auxiliarytool .el-tabs__content{
        min-height: 650px;
    } 
    .auxiliarytool .el-date-editor.el-input, 
    .auxiliarytool .el-date-editor.el-input__inner{
        width: 100% !important;
    }
    .scriptbox{ 
        border:1px solid #ebebeb;
        
    }
    .scriptbox >strong{
        display: inline-block;
        width: 100%;
        background: #f2f3f9;
        padding: 10px 20px;
        box-sizing: border-box;
        font-size: 14px;
    }
    .scriptbox >div{
        width: 100%;
        padding: 10px 20px;
        box-sizing: border-box;
        overflow-y: auto;
        overflow-x: hidden;
    }
    .Buttonbox{
        text-align: center;
    }
    .Buttonbox >div{ 
        display: inline-block;
        margin: 5px 0;
    }
    .slelecteddata {
        padding: 10px 5px;
    }
    .slelecteddata > div{
        background: #f0f2f5;
        border-radius: 4px;
        padding:5px 10px;
        font-size: 12px;
        color: #606266;
        margin: 5px 0;
    }
    .slelecteddata > div > span{
        display: inline-block;
        color: #4c7ae3;
        
    }
    .el-icon-error{
        vertical-align: middle;
        font-size: 14px;
        cursor: pointer;
    }
    /deep/ .shentoudia .el-form-item__label{
        text-align: left!important;;
    }
    /deep/ .shentoudia .el-form-item__label:after{
        left:-6px;
    }
</style>
<script>  
import About from "@/components/About.vue";
import Operation from "@/components/Operation";
import xzbutton from "@/components/XzButton.vue";
import delbutton from "@/components/DelButton.vue";
import { auxiliarytool } from '@/api/tool.js'
export default ({
    name:'dnslog',
    components:{
        About,
        Operation,
        delbutton,
        xzbutton
    },
    data(){  
    	return{
            multipleSelection: [],
            page:1,
            total:0,
            page_size:10,
            loading:false,  
    		activeName:'httplog',
            tableData:[],
            search_field:'',
            dnstableData: [],
            totalpage:0,
            currentpage:1,
            createInfiltration:false,
            pageType:'auxi',
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/auxiliarytool"; 
    },
    mounted:function(){   

    },
    methods:{  
        //  取消
        handlecancel(){
            this.$refs.edit.handleClearFiles();
            this.$refs.edit.form.remarks = ''
            this.createInfiltration = false;    
        },
         handleSelectionChange(val){
            this.multipleSelection = val;
        },
        clearall(){
            this.dialogFormVisible = true;
        },
        // handleCurrentChange(t){
        //     this.search_item.page = t;
        //     this.getData();
        // },
        // handleSizeChange(t){
        //     this.search_item.page = 1;
        //     this.pageSize = t;
        //     this.getData();
        // },
        handleReset(){
            this.search_field = '';
            this.page_num = 1;
            this.getHttploglist();
        },
        handlesearch(){
            this.page_num = 1; 
            this.getHttploglist();
            this.currentpage = 1;
        },
        currentchange(t){
            this.page_num = t; 
            this.getHttploglist();
            this.currentpage = t;
        },
        cancelform(){
            this.dialogFormVisible = false;
        },
        //  提交表单
        saveAndSendData(){
            let valiDatas = this.$refs.edit.handleEdit();
            if(valiDatas === null)
                return;
            const data = this.handleFormDate(valiDatas);
            // TODO: 接口调用
            // console.log(validatas);
            let config = {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            }
            this.$ajax({
                url:"/tools/assists/penetration_resources/create/",
                method:"POST",
                data
            })
            .then((data) => {
                let dt = data.data
                if (dt.success) {
                    this.$message.success(dt.msg);
                    this.createInfiltration = false;
                    this.getPenetrationResources();
                    this.$refs.edit.handleClearFiles();
                    this.$refs.edit.form.remarks = ''

                } else {
                    this.$message.error(dt.msg)
                }
            })
            .catch((data)=>{
                console.log(data); //错误信息
            });
            // TODO: 调用完成 重置form表单
            // this.$refs.edit.$refs.editRef.resetFields();
        },
    }
})
 
</script>
