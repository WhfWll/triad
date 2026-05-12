<template>
    <div>
        <div class="main-title">
             
             <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >发件配置管理</label>
        </div>
        <div class="gophishbox context_box_bg"> 
             <div class="search-box">
                <div class="operationbutton">
                    <xz-button type="primary" @click="AddSendprofile" size="small">新建配置</xz-button>
                   
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
                <el-table-column prop="interface_type" label="接口类型">
                </el-table-column>   
                <el-table-column prop="modified_date" label="修改时间" width="200">
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id"> 
                            <el-link class="link_primary" :underline="false"
                                @click="updatesendprofile(scope.row)" > 编辑</el-link>  
                             <el-link class="link_primary" :underline="false"
                                @click="copysendprofile(scope.row)" > 复制</el-link>  
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

        <el-dialog :title="sendprofileid!=0?'编辑发件配置':'创建发件配置'" :visible.sync="dialogaddFormVisible" 
        :before-close="cancelform" width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="sendTest" style="margin-right: 10px;">发送测试邮件</el-button>

                <el-button size="small" @click="submitForm">保存</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>

            <el-dialog
                title="发送测试邮件"
                :visible.sync="dialogVisible"
                width="50%"
                :show-close="false" 
                 :close-on-click-modal="false"
                 append-to-body
                :before-close="handleClose">
                <div>
                    <el-input v-model="sendtestform.first_name" size="small" style="width:150px;margin-right: 10px;" placeholder="名称" ></el-input>
                    <el-input v-model="sendtestform.email" size="small" style="width:200px;margin-right: 10px;" placeholder="邮箱地址" ></el-input>
                    <el-input v-model="sendtestform.position" size="small" style="width:150px;margin-right: 10px;" placeholder="职位"></el-input>
                </div>
                
                <span slot="footer" class="dialog-footer">
                    <el-button size="small" @click="handleClose">取 消</el-button>
                    <el-button size="small" type="primary" @click="handleSendTest">确 定</el-button>
                </span>
                </el-dialog> 

            <div style="padding:24px">
                <el-form :model="sendprofileform" label-width="0" status-icon ref="ruleFormadduser"  >
                    <el-row>
                        <el-col :span="12">
                             <el-form-item label=" " prop="username">
                                <label class="dialog_item_label">配置名称</label>
                                <el-input v-model="sendprofileform.name" size="small" style="width:320px" placeholder="请输入配置名称"
                                    maxlength="50"></el-input> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                              <el-form-item label=" " prop="username">
                                <label class="dialog_item_label">接口类型</label>
                                <el-input v-model="sendprofileform.interface_type" size="small" style="width:320px"  disabled
                                    maxlength="50"></el-input> 
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="12">
                            <el-form-item label=" " prop="username">
                                <label class="dialog_item_label">邮箱账号</label>
                                <el-input v-model="sendprofileform.host" size="small" style="width:320px"  placeholder="请输入邮箱账号"
                                    maxlength="50"></el-input> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                             <el-form-item label=" " prop="username">
                                <label class="dialog_item_label">发送邮件服务器</label>
                                <el-input v-model="sendprofileform.from_address" size="small" style="width:320px"  placeholder="请输入发送邮件服务器"
                                    maxlength="50"></el-input> 
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="12">
                             <el-form-item label=" " prop="username">
                                <label class="dialog_item_label">用户名</label>
                                <el-input v-model="sendprofileform.username" size="small" style="width:320px"  placeholder="请输入用户名"
                                    maxlength="50"></el-input> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                            <el-form-item label=" " prop="username">
                                <label class="dialog_item_label">密码</label>
                                <el-input v-model="sendprofileform.password" size="small" style="width:320px"  placeholder="请输入密码"
                                    maxlength="50"></el-input> 
                                <el-checkbox v-model="sendprofileform.ignore_cert_errors">忽略证书错误</el-checkbox>
                            </el-form-item>
                        </el-col>
                    </el-row> 
                    <div>
                        <div>
                              <div class="dialog_item_label">自定义邮件头</div> 
                        <el-button type="primary" size="mini" style="vertical-align: middle;float: right; "  
                                @click="clickadd()">添加</el-button>
                        </div>  
                         <div class="div_width"   style="margin-top:16px;margin-bottom:16px; ">
                            <el-table :data="login_conf" size="small" style="width: 100%">
                                <el-table-column  prop="target" label="配置名称">
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.key }}</span>    
                                        <el-input v-else v-model="scope.row.key" size="small" ></el-input>    
                                    </template>
                                </el-table-column>
                                <el-table-column prop="protocol"  label="配置值"> 
                                   <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.value }}</span>    
                                        <el-input v-else v-model="scope.row.value" size="small" ></el-input>    
                                    </template>
                                </el-table-column> 
                                <el-table-column  label="操作"  width="150">
                                    <template slot-scope="scope"  >
                                            <el-link :underline="false" @click="tbSave(scope)" 
                                                v-if="scope.row.dataShow" > 保存 </el-link>
                                            <el-link :underline="false" @click="tbUpdate(scope)" > 编辑 </el-link>
                                            <el-link :underline="false" @click="tbDelete(scope)" > 删除 </el-link>
                                        </template>
                                </el-table-column>
                            </el-table>
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
export default ({
    name: 'usermanagement',
    components: {
        XzButton,
        DelButton,
    },
    data() {
        return{
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
            sendprofileform:{
                name:'',
                interface_type:'SMTP',
                host:'',
                from_address:'',
                username:'',
                password:'',
                ignore_cert_errors:false,
            },
            login_conf:[],
            sendprofileid:0,
            dialogVisible:false,
            sendtestform:{
                first_name:'',
                email:'',
                position:'',
            }
        }
    },
    created(){
         this.$store.state.activefirstMenu = '/sendprofile';
    },
    mounted(){
        this.getData(); //获得列表数据
    },
    methods:{
        AddSendprofile(){
            this.dialogaddFormVisible = true;
        },
        async updatesendprofile(row){ //编辑 

            this.sendprofileid =row.id;
            this.dialogaddFormVisible = true;

            this.sendprofileform.name = row.name;
            this.sendprofileform.host = row.host;
            this.sendprofileform.from_address = row.from_address;
            this.sendprofileform.ignore_cert_errors = row.ignore_cert_errors;
            this.sendprofileform.interface_type = row.interface_type;
            this.sendprofileform.username = row.username;
            this.sendprofileform.password =row.password;

            
            const newArray = row.headers.map(item => ({
                ...item,
                dataShow: false
            }));
            this.login_conf = newArray;
 
        },
        //复制
        copysendprofile(row){
            this.sendprofileid=0;
            this.dialogaddFormVisible = true; 
            this.sendprofileform.name = 'copy of'+row.name;
            this.sendprofileform.host = row.host;
            this.sendprofileform.from_address = row.from_address;
            this.sendprofileform.ignore_cert_errors = row.ignore_cert_errors;
            this.sendprofileform.interface_type = row.interface_type;
            this.sendprofileform.username = row.username;
            this.sendprofileform.password =row.password;
 

            const newArray = row.headers.map(item => ({
                ...item,
                dataShow: false
            }));
            this.login_conf = newArray;

        },
        handlesearch(){
            this.getData();
        },
        handleReset(){
            this.formData.search = '';
            this.getData();
        },
        async getData() { 
            const dt = await gophish.profileall({
               page:this.page_num,
                size:this.pageSize,
                search:this.formData.search,
            });
            if (dt.code === 200) {
                this.tableData = dt.data.sendingProfile;
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
           
            const data = await gophish.profiledelete({
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
        //发送测试邮件
        async sendTest(){
            this.dialogVisible = true;
        },
        handleClose(){
            this.dialogVisible = false; 
            this.sendtestform.email = '';
            this.sendtestform.first_name = '';
            this.sendtestform.position = '';
        },
        //提交 发送测试邮件
        async handleSendTest(){
            var dt = await gophish.send_test_email({
                first_name:this.sendtestform.first_name,
                email:this.sendtestform.email,
                position:this.sendtestform.position,
                url:'',
                template:{},
                smtp:{
                    host:this.sendprofileform.host, 
                    from_address:this.sendprofileform.from_address,
                    ignore_cert_errors:this.sendprofileform.ignore_cert_errors,
                    username:this.sendprofileform.username,
                    password:this.sendprofileform.password,
                    headers:this.login_conf
                }
            })
            if (dt.code === 200) {
                this.$message({
                    message:'发送测试邮件成功',
                    type: 'success'
                });
                this.handleClose(); 
            } else {
                this.$message({
                    message: data.msg,
                    type: 'error'
                });
            }

        },
        async submitForm(){ 
            if(this.sendprofileid!=0){ //编辑
                const data = await gophish.profileupdate({
                    id:Number(this.sendprofileid),
                    name:this.sendprofileform.name,
                    host:this.sendprofileform.host,
                    interface_type:this.sendprofileform.interface_type,
                    from_address:this.sendprofileform.from_address,
                    ignore_cert_errors:this.sendprofileform.ignore_cert_errors,
                    username:this.sendprofileform.username,
                    password:this.sendprofileform.password,
                    headers:this.login_conf
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
                const data = await gophish.profilecreate({
                    name:this.sendprofileform.name,
                    host:this.sendprofileform.host,
                    interface_type:this.sendprofileform.interface_type,
                    from_address:this.sendprofileform.from_address,
                    ignore_cert_errors:this.sendprofileform.ignore_cert_errors,
                    username:this.sendprofileform.username,
                    password:this.sendprofileform.password,
                    headers:this.login_conf
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
            this.sendprofileform.name='';
            this.sendprofileform.host = '';
            this.sendprofileform.from_address = '';
            this.sendprofileform.ignore_cert_errors=false;
            this.sendprofileform.password='';
            this.sendprofileform.username = '';
            this.login_conf=[];
            this.sendprofileid = 0;
        },
        clickadd(){
             this.login_conf.push({
                dataShow: true,
            })
        },
        tbUpdate(scope){
             this.login_conf[scope.$index].dataShow = true;
        },
        tbDelete(scope){
            this.login_conf.splice(scope.$index, 1);
        },
        tbSave(scope){
            this.login_conf[scope.$index].dataShow = false;
        },
        clickupload(){
            document.querySelector(".btnUploadID").click();
        },
        changeuploaID(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1);
            if (fileSuffix.indexOf("csv") != -1) {
                
                var reader = new FileReader();
                reader.onload = function (e) {
                  
                    var data = e.target.result;
                    if (that.rABS) {
                        that.wb = XLSX.read(btoa(fixdata(data)), {
                            type: "base64",
                        });
                    } else {
                        that.wb = XLSX.read(data, {
                            type: "binary",
                        });
                    }  
                    let carData = XLSX.utils.sheet_to_json(
                        that.wb.Sheets[that.wb.SheetNames[0]],{header:1, }
                    ); 
                    let arr = []; 

                    console.log(carData)
                    for(var i=1;i<carData.length;i++){  

                        arr.push({
                            first_name:carData[i][0],
                            email:carData[i][2],
                            position:carData[i][3],
                        })
                    } 
                   that.login_conf = arr  
                };
                if (that.rABS) {
                    reader.readAsArrayBuffer(f);
                } else {
                    reader.readAsBinaryString(f);
                }
            } else{
                
            }
        },
    }
})
</script>