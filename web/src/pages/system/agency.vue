<template>
    <div>
        <div class="main-title  ">
            代理管理
        </div>
        <div class="nodelist div-list">
            <div class="search-box">
                <div class="operationbutton">
                    <!-- <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" trigger="click"
                        :visible-arrow="false" v-model="alldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消</el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDeleteTemplate">确定</el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除
                        </el-button>
                    </el-popover> -->
                     <xz-button type="primary" @click="Addagency" size="small" style="margin-right:8px">新建</xz-button>
                    <del-button :width="170" @click="btnMultiDeleteTemplate" :disabled="!multipleSelection.length">
                    </del-button>
                </div>
                <div class="serach-condition">
                    <div class="search-text">
                        <el-input placeholder="搜索代理名称" @keydown.enter.native="handlesearch" v-model="search_item.search" class="input-with-select"
                            size="small" clearable> </el-input>
                        <xz-button type="primary" @click="handlesearch" :disabled="false" size="small">搜索</xz-button>
                    </div>
                    <div>
                        <xz-button type="primary" @click="handleReset" :disabled="false" size="small">重置</xz-button>
                    </div>
                </div>
            </div>


            <el-table ref="multipleTable" :data="tableData" tooltip-effect="dark" v-model="Loading" style="width: 100%"
                @selection-change="handleSelectionChange" @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column type="selection" width="55">
                </el-table-column>
                <el-table-column prop="proxy_name" label="代理名称">
                    <template slot-scope="scope">
                        <el-link @click="btninfo(scope.row)">{{scope.row.proxy_name}}</el-link>
                    </template>
                </el-table-column>
                <el-table-column prop="proxy_ip" label="代理IP">
                </el-table-column>
                <el-table-column prop="proxy_port" label="代理端口">
                </el-table-column>
                <el-table-column prop="proxy_status" label="状态">
                    <template slot-scope="scope" >      
                        <span class="tag_status tag_success" v-if="scope.row.proxy_status == 'online'">在线</span> 
                        <span class="tag_status tag_danger" v-if="scope.row.proxy_status == 'offline'">离线</span>
                    </template>
                </el-table-column>
                <el-table-column prop="proxy_scheme" label=" 代理协议">
                    <template slot-scope="scope" > 
						<div v-if="showEditFileNameButton && rowId == scope.row.id">
                            <!-- <el-link :underline="false" class="link_primary">启用</el-link> -->
							<el-popover placement="bottom" width="170" :visible-arrow="false"
								:ref="`popover_id-${scope.row.id}`"
								popper-class="delButton_popper" >
								<p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
								<div style="text-align: right; margin: 0">
									<el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button> 
									<el-button size="mini" type="primary" @click="handleDel(scope)">确定</el-button>
								</div> 
								<el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference" >删除</el-link>  
							</el-popover>   
						</div>
						<div v-else >
							{{ scope.row.proxy_scheme }}
						</div>
					</template>
                </el-table-column>
                <!-- <el-table-column prop="proxy_type" label="代理类型">11-10号罗要求注释 wm
                    <template slot-scope="scope" > 
						<div v-if="showEditFileNameButton && rowId == scope.row.id">
							<el-popover placement="bottom" width="170" :visible-arrow="false"
								:ref="`popover_id-${scope.row.id}`"
								popper-class="delButton_popper" >
								<p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
								<div style="text-align: right; margin: 0">
									<el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button> 
									<el-button size="mini" type="primary" @click="handleDel(scope)">确定</el-button>
								</div> 
								<el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference" >删除</el-link>  
							</el-popover>   
						</div>
						<div v-else >
							<span v-if="scope.row.proxy_type == 'static'">静态代理</span>
                            <span v-if="scope.row.proxy_type == 'dynamic'">动态代理</span>
						</div>
					</template>
                </el-table-column> -->

                <!-- <el-table-column label="状态">
                    <template slot-scope="scope">

                        <el-link :underline="false" class="link_primary">启用</el-link>

                        <el-popover placement="bottom" width="170" :visible-arrow="false"
                            :ref="`popover_id-${scope.$index}`" popper-class="delButton_popper">
                            <p class="delText">
                                <i class="el-icon-warning"></i>确定删除吗？
                            </p>
                            <div style="text-align: right; margin: 0">
                                <el-button size="mini" class="delCancel"
                                    @click="scope._self.$refs[`popover_id-${scope.$index}`].doClose()">取消
                                </el-button>
                                <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                            </div>
                            <span slot="reference">删除</span>
                        </el-popover>


                    </template>
                </el-table-column> -->
            </el-table>
            <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange"
                :current-page="currentPage" :page-size="pageSize" layout=" total,  prev, pager, next, sizes,jumper"
                :total="total">
            </el-pagination>


        </div>
        <!-- 新建代理 -->
        <el-dialog title="新建代理" :visible.sync="dialogaddFormVisible" :before-close="cancelform" width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">关闭</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="agencyform" label-width="0" status-icon ref="ruleForm" :rules="rules1">
                    <el-form-item label=" " prop="name">
                        <label class="dialog_item_label">代理名称</label>
                        <el-input v-model="agencyform.name" size="small" style="width:320px" autocomplete="off" placeholder="请输入代理名称">
                        </el-input>
                    </el-form-item>
                    <!-- <el-form-item label=" " prop="type"> 11-10号罗要求注释 wm
                        <label class="dialog_item_label">代理类型</label>
                        <el-select v-model="agencyform.type" size="small" placeholder="请选择" style="width: 320px">
                            <el-option v-for="(item, index) in type_list" :key="index" :label="item.name" :value="item.id"></el-option>
                        </el-select> 
                    </el-form-item> -->
                    <el-form-item label=" " prop="agreement">
                        <label class="dialog_item_label">协议</label>
                        <el-select v-model="agencyform.agreement" size="small" placeholder="请选择" style="width: 320px">
                            <el-option v-for="(item, index) in agreement_list" :key="index" :label="item.name" :value="item.id"></el-option>
                        </el-select>
                    
                    </el-form-item> 
                    <el-form-item label=" " prop="agencyip"  >
                        <label class="dialog_item_label">代理IP</label>
                        <el-input v-model="agencyform.agencyip" size="small" style="width:320px" placeholder="请输入代理IP"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label="" prop="agencyport" >
                        <label class="dialog_item_label">代理端口</label>
                        <el-input v-model="agencyform.agencyport" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入代理端口"></el-input>
                    </el-form-item>
                   
                    
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">账号</label>
                        <el-input v-model="agencyform.username" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入账号"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="password">
                        <label class="dialog_item_label">密码</label>
                        <el-input type="password" v-model="agencyform.password" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入密码"></el-input>
                    </el-form-item>
                </el-form>
            </div>
        </el-dialog>
    </div>
</template>
<style  scoped lang="less">
.div-list{ 
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.12);
    border-radius: 4px;
}
.info_box {
    padding: 24px;
}

.basicinfo {
    padding: 24px;
    background: #fff;
    border: 1px solid rgba(232, 232, 245, 1);
    margin-bottom: 24px;
}

.otherinfo {
    margin-top: 32px;
    >div{
        margin-bottom: 16px;
        label{
            margin-left: 6px;
            color: rgba(72, 72, 102, 0.87);
        }
        span{
            display: inline-block;
            color: rgba(72, 72, 102, 0.63);
            font-size: 13px;

        }
    }
    .bgcolorfff{
        background-color: #fff;
        border: 1px solid rgba(232, 232, 245, 1);
    }
    .part_title {
        margin-bottom: 16px;
        >label{
            margin-left: 0;
        }
        >span{
            margin-left: 32px;
        }
    }

    .content {
        // background: rgba(255, 255, 255, 1);
        border-radius: 2px;
        // border: 1px solid rgba(232, 232, 245, 1);
        padding: 12px 16px;
        color: rgba(72, 72, 102, 0.64);
        font-size: 13px;
    }
}
.minH{
    height: 50px;
}
.part_title {
    font-size: 14px;
    margin-bottom: 16px;
    font-weight: 500;
    border-left: 3px solid #4C7AE3;
    padding-left: 10px;
    height: 14px;
    line-height: 14px;
    color: rgba(72, 72, 102, 0.87);
}
.chartdiv{
    background-color: #fff;
    height: 300px;
    box-shadow: 0px 2px 4px 0px rgb(76 122 227 / 12%);
    >label{
        width: 100%;
        height: 60px;
        line-height: 60px;
        padding-left: 24px;
        font-weight: bold;
    }
    >div{
        height: calc(100% - 80px);
    }
}
/deep/ .el-dialog__body{
    background-color: rgb(247, 247, 251);
}
/deep/ .el-form-item__label {
    position: absolute;
    left: 127px;
}

 /deep/ .el-form-item__error {
    left: 112px;
}
</style>
<script>
import DelButton from "@/components/DelButton.vue";
import XzButton from "@/components/XzButton.vue";
import { proxy } from '@/api/system.js';
var echarts = require('echarts');
export default {
    name:"agency",
    components:{
        XzButton,
        DelButton
    },
    data(){
        var validatePwd =  (rule, value, callback) => {
            const reg = /^\S+$/; 
            if (reg.test(value)) {
                callback();
            } else {
                return callback(new Error('格式不正确(不能包含空格)'));
            }
        }; 
        return{
            showEditFileNameButton:false,
            rowId:'',
            Loading:false,
            alldelvisible:false,
            search_item: {
                page:1,
                search:'',
            },
            multipleSelection: [],
            tableData:[],
            currentPage:1,
            pageSize:10,
            total:0,
            updatetxt:'编辑',
            is_Update:false,
            dialogVisible:false,
            nodeform:{
                name:'',
            },
            agencyform:{
                id:'',
                type:'',
                name:'',
            	agencyport:'',
            	agencyip:'',
            	password:'',
                username:'',
            	agreement:'',
                
            },
            agreement_list: [
                {
                    id: 'Socks5',
                    name: 'Socks5',
                },
                {
                    id: 'HTTP',
                    name: 'HTTP',
                },
                {
                    id: 'HTTPS',
                    name: 'HTTPS',
                },
            ],
            type_list:[
                {
                    id: 'static',
                    name: '静态代理',
                },
                {
                    id: 'dynamic',
                    name: '动态代理',
                }

            ],
            basicinfo:[],
            dialogaddFormVisible:false,
            // dialogupdateFormVisible:false,
            rules1: {
                name: [
                    { required: true, message: '代理名称不能为空', trigger: 'blur' },
                ],
                // type: [
                //     { required: true, message: '请选择代理类型', trigger: 'change' },
                // ],
                agreement: [
                    { required: true, message: '请选择代理协议', trigger: 'change' },
                ],
                agencyip: [
                    { required: true, message: '代理IP不能为空', trigger: 'blur' }, 
                ],
                agencyport: [
                    { required: true, message: '代理端口不能为空', trigger: 'blur' },
                    { max: 50, message: '代理端口不能超过50', trigger: 'blur' },
                ], 
                username: [
                    { required: true, message: '账号不能为空', trigger: 'blur' }, 
                ],
                password: [
                    { required: true, message: '密码不能为空', trigger: 'blur' }, 
                    { validator: validatePwd, trigger: 'blur' }
                ],
            },

        }
        
    },
    created: function () {
        this.$store.state.activefirstMenu = "/agency"; 
    },
    mounted: function () {
        this.getData();
     
    },
    methods: {
        Addagency(){  
        	this.dialogaddFormVisible = true; 
            this.agencyform.username = '';
            this.agencyform.id = '';
            this.agencyform.password = '';
            this.agencyform.type = '';
            this.agencyform.name = '';
            this.agencyform.agencyip = '';
            this.agencyform.agreement = '';
            this.agencyform.agencyport = '';
        },
        cancelform(){
            this.agencyform.username = '';
            this.agencyform.id = '';
            this.agencyform.password = '';
            this.agencyform.type = '';
            this.agencyform.name = '';
            this.agencyform.agencyip = '';
            this.agencyform.agreement = '';
            this.agencyform.agencyport = '';
            this.dialogaddFormVisible = false; 
            this.$refs.ruleForm.resetFields(); 
        },
        submitForm(){  
            this.$refs.ruleForm.validate( async (valid) => {
                if (valid) { 
                    let _j={
                        proxy_name: this.agencyform.name,
                        proxy_ip: this.agencyform.agencyip,
                        proxy_port: this.agencyform.agencyport,
                        proxy_scheme: this.agencyform.agreement,
                        username: this.agencyform.username,
                        password: this.agencyform.password,
                        proxy_type:this.agencyform.type,
                    } 
                    if (!this.agencyform.id){ //添加
                        const data = await proxy.addProxy(_j); 
                        if (data.success) {
                            this.$message({
                                message: '新增代理成功',
                                type: 'success'
                            });
                        } else {
                            this.$message({
                                message: data.error,
                                type: 'error'
                            });
                        }   
                    }
                    else{
                        _j.id = this.agencyform.id;
                        const data = await proxy.editProxy(_j); 
                        if (data.success) {
                            this.$message({
                                message: '编辑代理成功',
                                type: 'success'
                            });
                        } else {
                            this.$message({
                                message: data.error,
                                type: 'error'
                            });
                        }   
                    }
                    
                    this.dialogaddFormVisible = false; 
                    this.getData(); 
                    this.agencyform.username = '';
                    this.agencyform.id = '';
                    this.agencyform.password = '';
                    this.agencyform.type = '';
                    this.agencyform.name = '';
                    this.agencyform.agencyip = '';
                    this.agencyform.agreement = '';
                    this.agencyform.agencyport = '';
                }
            }); 
            
        }, 
        btninfo(row) {

            this.agencyform.id = row.id;
            this.agencyform.name = row.proxy_name;
            this.agencyform.username = row.username;
            this.agencyform.type = row.proxy_type;
            this.agencyform.password = row.password;
            this.agencyform.agencyport = row.proxy_port;
            this.agencyform.agencyip = row.proxy_ip;
            this.agencyform.agreement = row.proxy_scheme;
            this.dialogaddFormVisible = true;
        },

        async getData(){
            this.Loading = true;  ///task/template/
            let params = {
                search: this.search_item.search,
                page: this.search_item.page,
                page_size: this.pageSize,
            }
            const res = await proxy.getData(params)

            this.Loading = false;
            this.tableData = res.data;
            this.total = res.count;
        },
        handlesearch(){
            this.search_item.page = 1;
            this.currentPage = 1;
            this.getData();
        },
        handleReset(){
            this.search_item.page = 1;
            this.search_item.search = '';
            this.pageSize = 10;
            this.currentPage = 1;
            this.getData();
        },
        async btnMultiDeleteTemplate(){ //批量删除
            if (this.multipleSelection.length == 0) return;
            var _ids = this.multipleSelection.map(item => item.id);
            const res = await proxy.MultDelete({
                ids: _ids.join(',')
            })
            if (res.success) {
                this.$message({
                    message: res.message,
                    type: 'success'
                });
                this.alldelvisible = false
                this.getData();
            } else {
                this.$message({
                    message: res.error,
                    type: 'error'
                });
            }
        },
        async btnDel(scope) { //单个删除
            const res = await proxy.MultDelete({
                ids: scope.row.id
            })
            if (res.success) {
                this.$message({
                    message: '删除节点成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.$index}`].doClose();
                this.getData(); 
            } else {
                this.$message({
                    message: res.error,
                    type: 'error'
                });
            }
            // document.querySelector('body').click()
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        handleCurrentChange(t){
            this.search_item.page = t;
            this.getData();
        },
        handleSizeChange(t){
            this.search_item.page = 1;
            this.pageSize = t;
            this.getData();
        },
        btnUpdate(){
            this.is_Update=true;
        },
        saveUpdate(){

        },
        handleDel:function(scope){ //删除 
				this.$ajax({
	                method:'delete',
	                url:'/systems/information/delete/',
	                data: {
	                    id:scope.row.id+''
	                } 
	            })
	            .then(data => { 
	                var dt = data.data;  
	                if(dt.success){ 
	                	this.$message({
	                        message:dt.msg,
	                        type: 'success'
						});
						scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
	                    this.getData();
	                }else{
	                    this.$message({
	                        message:dt.error,
	                        type: 'error'
	                    });
	                }  
	                
	            })
	            .catch(data=>{
	                console.log(data); //错误信息
	            }); 
    	},
        mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
			
            let t = this.$refs['popover_id-' + row.id].showPopper;
            if(!t){
                  this.showEditFileNameButton = false;
            this.rowId = "";
            }
          
        },
       
    },
}
</script>