<template>
    <div>
        <div class="main-title">
             
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >钓鱼网站管理</label>
        </div>
        <div class="gophishbox context_box_bg"> 
             <div class="search-box">
                <div class="operationbutton">
                    <xz-button type="primary" @click="AddTemplate" size="small">新建钓鱼网站</xz-button> 
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

        <el-dialog  :title="pageid!=0?'编辑钓鱼网站':'新建钓鱼网站'" 
            :visible.sync="dialogaddFormVisible" 
            :before-close="cancelform" width="1184px"
            :close-on-click-modal="false" 
            :validate-on-rule-change="false" 
            :show-close="false"  
         
            class="newUserDialog">
                <!--    
                @open="initEditor"   
                @close="destroyEditor"  
                @opened="onDialogOpened"    
                @closed="onDialogClosed" -->
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>

            <el-dialog
                title="导入网站"
                :visible.sync="dialogVisible"
                width="50%"
                :show-close="false" 
                 append-to-body
                :before-close="handleClose">
                <el-input v-model="import_site_url" size="small" style="width:100%"  placeholder="请输入网站地址" ></el-input>
                <span slot="footer" class="dialog-footer">
                    <el-button size="small" @click="handleClose">取 消</el-button>
                    <el-button size="small" type="primary" @click="handleImport">确 定</el-button>
                </span>
                </el-dialog> 
            <div style="padding:24px">
                <el-form :model="pageform" label-width="0" status-icon ref="ruleFormadduser"  > 
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">钓鱼网站名称</label>
                        <el-input v-model="pageform.name" size="small" style="width:350px" placeholder="请输入钓鱼网站名称"
                            maxlength="50"></el-input>
                        <div style="margin-left: 114px;">
                            <el-button type="primary" size="small" style="vertical-align: middle; margin-right: 16px"
                            @click="importsite()">导入网站</el-button> 
                        </div> 
                    </el-form-item> 
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">钓鱼网站内容</label>  
                        <el-tabs v-model="activeName" type="card" @tab-click="handleClick"> 
                            <el-tab-pane label="HTML" name="tabs1">  
                                  <ckeditor1 ref="editor1"  @change="changeeditor1"  id='editor1' :content="content"></ckeditor1>  
                            </el-tab-pane> 
                            <!-- <el-tab-pane label="Text" name="tabs2">
                                 <el-input type="textarea" v-model="pageform.text" size="small" :rows="10" resize="none"  ></el-input>
                            </el-tab-pane> -->
                        </el-tabs> 
                    </el-form-item>  
                    <el-form-item label=" " prop="capture_credentials">
                          <el-checkbox v-model="pageform.capture_credentials">捕获提交的数据</el-checkbox>
                    </el-form-item> 
                    <div v-if="pageform.capture_credentials">
                        <el-form-item label=" " prop="capture_passwords">
                            <el-checkbox v-model="pageform.capture_passwords">捕获密码</el-checkbox>
                        </el-form-item> 
                        <el-form-item label=" " prop="capture_passwords">
                            <label class="dialog_item_label">重定向到</label>
                            <el-input v-model="pageform.redirect_url" size="small" style="width:550px"  ></el-input>
                        </el-form-item> 
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
/* 控制编辑器最大高度 */
::v-deep .w-e-text-container {
  max-height: 400px !important;
}

/* 控制内容中的图片不要过大 */
::v-deep .w-e-text img {
  max-width: 100% !important;
  height: auto !important;
  border-radius: 4px;
}

/* 防止 dialog 被内容撑破 */
.editor-container {
  overflow: hidden;
  padding: 10px 20px;
}
</style>
<script>
 

import XzButton from "../../components/XzButton.vue";
import DelButton from "../../components/DelButton.vue";
import gophish from "@/api/gophish.js";
import _ from 'lodash' 
import E from 'wangeditor' 
import ckeditor1 from "./components/ckeditor.vue"

export default ({
    name: 'gophishpage',
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
            pageform:{
                name:'',   
                html:'',
                capture_credentials:false,
                capture_passwords:false,
                redirect_url:'',
            }, 
            pageid:0,
            editor: null,
            content: '',
            attachments:[],
            dialogVisible:false,
            import_site_url:'',
            editorConfig: {
                language: 'zh-CN',  // 中文界面
                toolbar: [
                    // ['Bold', 'Italic', 'Underline'],  // 基础格式
                    // ['NumberedList', 'BulletedList'], // 列表
                    // ['Link', 'Unlink', 'Image'],      // 链接、图片
                    // ['Source'],// 源码模式          
                    ['Cut','Copy','Paste','PasteText','PasteFromWord','-','Print', 'SpellChecker', 'Scayt'],
                    ['Undo','Redo','-','Find','Replace','-','SelectAll','RemoveFormat'],
                    ['Form', 'Checkbox', 'Radio', 'TextField', 'Textarea', 'Select', 'Button', 'ImageButton', 'HiddenField'],
                    ['Source','-','Save','NewPage','Preview','-','Templates'],  
                    ['Bold','Italic','Underline','Strike','-','Subscript','Superscript'],
                    ['NumberedList','BulletedList','-','Outdent','Indent','Blockquote'],
                    ['JustifyLeft','JustifyCenter','JustifyRight','JustifyBlock'],
                    ['Link','Unlink','Anchor'],
                    ['Image','Flash','Table','HorizontalRule','Smiley','SpecialChar','PageBreak'],
                   
                    ['Styles','Format','Font','FontSize'],
                    ['TextColor','BGColor'] ,       
                   
                ],
                // 图片上传配置（后端接口）
                filebrowserImageUploadUrl: '/api/upload/image',
                filebrowserUploadMethod: 'form',
                removePlugins: ['about'], 
            },
            // editor: ClassicEditor  // 指向本地导入的编辑器
        }
    },
    created(){
         this.$store.state.activefirstMenu = '/gophishpage';
    },
    mounted(){
        this.getData(); //获得列表数据 
    },
    beforeDestroy() {
        this.onDialogClosed()
    },
    methods:{ 
        //富文本 
        changeeditor1(val) { 

            this.content  = val;
            console.log('changeeditor1'+val); 
            //在做其他的处理 

            console.log('content:'+this.content); 
           
        },
       // 弹出层打开时初始化编辑器
        initEditor() {
            // 确保 DOM 已渲染（延迟一小段时间，避免 el-dialog 动画未完成）
            this.$nextTick(() => {
                const editor = this.$refs.ckeditor.instance
                if (editor) {
                // 可在此处设置编辑器内容或其他操作
                editor.setData(this.content)
                }
            })
        },
        // 关闭弹出层时销毁编辑器（解决多次打开时的冲突）
        destroyEditor() {
            const editor = this.$refs.ckeditor?.instance
            if (editor) {
                editor.destroy()  // 销毁实例，避免内存泄漏
            }
        },
        // 保存富文本内容
        handleSave() {
            const content = this.content
            console.log('保存的内容:', content)
            // 调用接口提交...
            this.dialogVisible = false
        },
        AddTemplate(){
            this.dialogaddFormVisible = true;
             
        },
        // Dialog 完全打开后初始化编辑器
        onDialogOpened() {
            this.onDialogClosed() // 先销毁旧实例
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
            this.editor.txt.html(this.content )
        },

        // Dialog 关闭时销毁编辑器
        onDialogClosed() {
            if (this.editor) {
                this.editor.destroy()
                this.editor = null
            }
        }, 
        updateTemplate(row){ //编辑 
            this.pageid =row.id;
            this.dialogaddFormVisible = true;

            this.pageform.name = row.name;
            this.pageform.redirect_url = row.redirect_url;
            this.pageform.capture_credentials = row.capture_credentials; 
            this.pageform.capture_passwords = row.capture_passwords; 
            this.content = row.html; 
            
        },
        copyTemplate(row){  //复制
            this.dialogaddFormVisible = true; 
            this.pageform.name = 'copy of '+row.name;
            this.pageform.redirect_url = row.redirect_url;
            this.pageform.capture_credentials = row.capture_credentials; 
            this.pageform.capture_passwords = row.capture_passwords; 
            this.content = row.html; 
        },
        //导入网站 弹出层
        importsite(){
            this.dialogVisible = true;
        },
        handleClose(){
            this.dialogVisible = false;
            this.import_site_url = '';
        },
        //导入网站
        async handleImport(){ 
            const dt = await gophish.pageimport_site({
                url:this.import_site_url,
                include_resources:false,
            })
            if (dt.code === 200) {
                this.content = dt.data.html; 
                // this.editor.txt.html(this.content ) 
                this.dialogVisible = false;
                this.import_site_url = '';
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        handlesearch(){
            this.getData();
        },
        handleReset(){
            this.formData.search = '';
            this.getData();
        },
        async getData() { 
            const dt = await gophish.pageall({
                page:this.page_num,
                size:this.pageSize,
                search:this.formData.search,
            });
            if (dt.code === 200) {
                this.tableData = dt.data.landingPage;
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
           
            const data = await gophish.pagedelete({
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
            if(this.pageid!=0){ //编辑
                const data = await gophish.pageupdate({
                    id:Number(this.pageid),
                    name:this.pageform.name, 
                    html:this.content, 
                    capture_credentials:this.pageform.capture_credentials,
                    capture_passwords:this.pageform.capture_passwords,
                    redirect_url:this.pageform.redirect_url,
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
                const data = await gophish.pagecreate({
                    name:this.pageform.name, 
                    html:this.content, 
                    capture_credentials:this.pageform.capture_credentials,
                    capture_passwords:this.pageform.capture_passwords,
                    redirect_url:this.pageform.redirect_url,

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
            this.pageform.name='';
            this.pageform.redirect_url='';
            this.pageform.capture_credentials=false;
            this.pageform.capture_passwords=false;
            this.$refs.editor1.setData('');
            this.content='';
            this.pageform.html=''; 
            this.pageid = 0;
        },  
        handleClick(){

        },
    }
})
</script>