<template>
    <div>
        <div class="main-title"> 
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >邮件模板管理</label>
        </div>
        <div class="gophishbox context_box_bg"> 
             <div class="search-box">
                <div class="operationbutton">
                    <xz-button type="primary" @click="AddTemplate" size="small">新建模板</xz-button> 
                </div> 
                <div class="serach-condition"> 
                    <div class="search-text">
                        <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select" size="small"
                        >
                        </el-input>
                        <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                    </div>
                </div>
            </div>
             <el-table :data="tableData" style="width: 100%"  @selection-change="handleSelectionChange"  height="calc(100% - 102px)"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave"> 
                <el-table-column prop="name" label="名称">
                </el-table-column>   
                <el-table-column prop="last_time" label="修改时间"  >
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id"> 
                            <el-link class="link_primary" :underline="false"
                                @click="updateTemplate(scope.row)" > 编辑</el-link> 
                            <el-link class="link_primary" :underline="false"
                                @click="copyTemplate(scope.row)" > 复制</el-link>  
                            <el-popover placement="bottom" width="170" :ref="`popover-${scope.row.id}`"
                                popper-class="delButton_popper">
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="fncancel(scope)">取消</el-button>
                                    <el-button size="mini" type="primary" @click="fnDel(scope)">确定</el-button>
                                </div>
                                <el-link :underline="false" class="link_primary" style="padding:0"
                                    slot="reference">删除</el-link>
                            </el-popover>  
                        </div>
                        <div v-else>
                            <span class="">{{ scope.row.modified_date }}</span> 
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination :page-size="10" background layout=" total,  prev, pager, next,sizes, jumper" :total="totalpage"
                :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>
        </div>

        <el-dialog  :title="templateid!=0?'编辑邮件模板':'新建邮件模板'" 
            :visible.sync="dialogaddFormVisible" 
            :before-close="cancelform" width="1184px"
            :close-on-click-modal="false" 
            :validate-on-rule-change="false" 
            :show-close="false"  
            class="newUserDialog">

            <!--  @opened="onDialogOpened"   
            @closed="onDialogClosed" -->

            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="templteform" label-width="0" status-icon ref="ruleFormadduser"  >
                    <el-row>
                        <el-col :span="12">
                            <el-form-item label=" " prop="username">
                                <label class="dialog_item_label">模板名称</label>
                                <el-input v-model="templteform.name" size="small" style="width:300px" placeholder="请输入模板名称"
                                    maxlength="50"></el-input>
                                <!-- <div style="margin-left: 114px;">
                                    <el-button type="primary" size="small" style="vertical-align: middle; margin-right: 16px"
                                    @click="clickupload()">批量导入邮箱账号</el-button>
                                    <input type="file" class="btnUploadID" ref="upload" @change="changeuploaID($event)"
                                        style="display: none" id="input-file-ID" />
                                    <el-link  class="link_primary" style="vertical-align: middle;" @click="downfile()">下载模板</el-link> 
                                </div>  -->
                            </el-form-item>
                        </el-col>
                        <el-col :span="12"> 
                            <el-form-item label=" " prop="username">
                                <label class="dialog_item_label">显示发件人</label> 
                                <el-input v-model="templteform.envelope_sender" size="small" style="width:300px" placeholder="<test@example.com>"
                                    maxlength="50"></el-input>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">邮件主题</label> 
                        <el-input v-model="templteform.subject" size="small" style="width:300px" placeholder="请输入邮件主题"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">邮件正文</label>  
                        <el-tabs v-model="activeName" type="card" @tab-click="handleClick"> 
                            <el-tab-pane label="HTML" name="tabs1">
                                <!-- <div ref="editorContainer" class="editor-container"></div> -->
                                  <ckeditor1 ref="editor1"  @change="changeeditor1"  id='editor1' :content="content"></ckeditor1>  
                            </el-tab-pane> 
                            <el-tab-pane label="Text" name="tabs2">
                                 <el-input type="textarea" v-model="templteform.text" size="small" :rows="10" resize="none"  ></el-input>
                            </el-tab-pane>
                        </el-tabs>

                    </el-form-item>
                    <div>
                        <div>
                            <div class="dialog_item_label">邮件附件</div>  
                            <el-button type="primary" size="mini" style="vertical-align: middle;  "  
                                @click="clickupload()">添加附件</el-button> 
                            <input type="file" class="btnUploadID" ref="upload" @change="changeuploaID($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                        <div >
                            <div v-for="(item,i) in attachments" :key="i" style="margin-top: 10px;" >
                                <span><i class="el-icon-paperclip"></i> {{item.name}} </span> 
                                <i class="el-icon-circle-close" title="删除" style="margin-left: 10px;" @click="delfile(i)"></i>
                            </div>
                        </div>
                    </div>

                </el-form>
            </div>
        </el-dialog>


    </div>
</template>
<style scoped lang="less">
.gophishbox {
    padding: 24px;
    background: #fff;
    height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
</style>
<script>
import XzButton from "../../components/XzButton.vue";
import DelButton from "../../components/DelButton.vue";
import gophish from "@/api/gophish.js";
import _ from 'lodash'
var XLSX = require('xlsx');
import E from 'wangeditor'

import ckeditor1 from "./components/ckeditor.vue"
export default ({
    name: 'usermanagement',
    components: {
        XzButton,
        DelButton,
        ckeditor1,
    },
    data() {
        return{
            activeName:'tabs1',
            multipleSelection: [],
            currentpage: 1,
            totalpage: 0,
            page_num: 1,
            pageSize:10,
            tableData: [],
            isMore: false,
            showOperateButton: false,
            rowId:'',
            formData:{
                search:'',
            },
            dialogaddFormVisible:false,
            templteform:{
                name:'',
                envelope_sender:'',
                subject:'',
                text:'',
            },
            login_conf:[],
            templateid:0,
            editor: null,
            content: '',
            attachments:[],
        }
    },
    created(){
         this.$store.state.activefirstMenu = '/template';
    },
    mounted(){
        this.getData(); //获得列表数据
       
    },
    methods:{ 
        //富文本 
        changeeditor1(val) { 

            this.content  = val;
            console.log('changeeditor1'+val); 
            //在做其他的处理 

            console.log('content:'+this.content); 
           
        },
        AddTemplate(){
            this.dialogaddFormVisible = true; 
        },
        // Dialog 完全打开后初始化编辑器
        onDialogOpened() {
            // 避免重复初始化
            if (this.editor) return

            this.editor = new E(this.$refs.editorContainer)

            // 配置编辑器
            this.editor.config.height = 250
            this.editor.config.zIndex = 1000
            this.editor.config.uploadImgShowBase64 = false // 关闭 base64，建议走服务器上传
            this.editor.config.uploadImgServer = '/api/upload' // 图片上传接口
            this.editor.config.onchange = (html) => {
                this.content = html
            }

            // 创建编辑器
            this.editor.create()

            // 可选：设置默认内容
            this.editor.txt.html(this.content || '<p>请输入邮件内容...</p>')
        },

        // Dialog 关闭时销毁编辑器
        onDialogClosed() {
            if (this.editor) {
                this.editor.destroy()
                this.editor = null
            }
        }, 
        updateTemplate(row){ //编辑 
            this.templateid =row.id;
            this.dialogaddFormVisible = true;

            this.templteform.name = row.name;
            this.templteform.subject = row.subject;
            this.templteform.envelope_sender = row.envelope_sender;
            this.attachments = row.attachments||[];
            this.content = row.html;
            this.templteform.text = row.text;
 
        },
        copyTemplate(row){  //复制
            this.dialogaddFormVisible = true; 
            this.templteform.name = 'copy of '+row.name;
            this.templteform.subject = row.subject;
            this.templteform.envelope_sender = row.envelope_sender;
            this.attachments = row.attachments || [];
            this.content = row.html;
            this.templteform.text = row.text;
        },
        // 下载模板
        downfile(){
             window.open('/group_template.csv');
        },
        handlesearch(){
            this.getData();
        },
        handleReset(){
            this.formData.search = '';
            this.getData();
        },
        async getData() { 
            const dt = await gophish.templateall({
                page:this.page_num,
                size:this.pageSize,
                search:this.formData.search,
            });
            if (dt.code === 200) {
                this.tableData = dt.data.templates;
                this.totalpage = dt.data.total;
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        handleSelectionChange(val) {
            this.multipleSelection = val;
            this.selectedID = []
            this.multipleSelection.forEach(item => {
                this.selectedID.push(item.role)
            })
        },
        currentchange(t) {
            this.page_num = t;
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t) {
            this.page_num = 1;
            this.pageSize = t;
            this.getData();
        },
        mouseenter(row, colum, cell, event) {
            this.showOperateButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event) {
            if (!this.$refs['popover-' + row.id]) {
                this.showOperateButton = false;
                this.rowId = "";
                return;
            } else {
                let isShow = this.$refs['popover-' + row.id].showPopper;
                if (!isShow) {
                    this.showOperateButton = false;
                    this.rowId = "";
                }

            } 
        },
        fncancel(scope) {
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
            scope._self.$refs[`popover-${scope.row.id}`].doClose()
        },
        async fnDel(scope) { //单个删除
           
            const data = await gophish.templatedelete({
                id:scope.row.id
            });
            if (data.code === 200) {
                this.$message({
                    message: data.msg || '删除成功！' ,
                    type: 'success'
                });
               
                this.getData();
            } else {
                this.$message({
                    message: data.msg || '错误，请重试',
                    type: 'error'
                });
            }

        },
        async submitForm(){ 
            this.content = this.$refs.editor1.getData(); 
            if(this.templateid!=0){ //编辑
                const data = await gophish.templateupdate({
                    id:Number(this.templateid),
                    name:this.templteform.name,
                    subject:this.templteform.subject,
                    envelope_sender:this.templteform.envelope_sender,
                    html:this.content,
                    text:this.templteform.text,
                    attachments:this.attachments
                });
                if (data.code === 200) {
                    this.$message({
                        message:'编辑保存成功',
                        type: 'success'
                    });
                    this.dialogaddFormVisible = false;
                    this.getData();
                    this.cancelform();
                } else {
                    this.$message({
                        message: data.msg,
                        type: 'error'
                    });
                }
            }else{ //新增
                const data = await gophish.templatecreate({
                    name:this.templteform.name,
                    subject:this.templteform.subject,
                    envelope_sender:this.templteform.envelope_sender,
                    html:this.content,
                    text:this.templteform.text,
                    attachments:this.attachments
                });
                if (data.code === 200) {
                    this.$message({
                        message:'创建保存成功',
                        type: 'success'
                    });
                    this.dialogaddFormVisible = false;
                    this.getData();
                    this.cancelform();
                } else {
                    this.$message({
                        message: data.msg,
                        type: 'error'
                    });
                }
            }

           
        },
        cancelform(){
            this.dialogaddFormVisible = false;
            this.templteform.name='';
            this.templteform.subject='';
            this.templteform.envelope_sender='';
             this.$refs.editor1.setData('');
            this.content='';
            this.templteform.text='';
            this.attachments=[];
            this.templateid = 0;
        }, 
        clickupload(){
            document.querySelector(".btnUploadID").click();
        },
        changeuploaID(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
       
            var reader = new FileReader();
            reader.onload = function (e) { 
                var data = e.target.result;
 
                var _type = data.split(',')[0].split(';')[0].split(':')[1];
                var base64str = data.split(',')[1];

                that.attachments.push({
                    name:f.name,
                    type:_type,
                    content:base64str
                })
                 
                
            };
            reader.readAsDataURL(f);  
        },
        delfile(i){
            this.attachments.splice(i,1);
        },
        handleClick(){

        },
    }
})
</script>