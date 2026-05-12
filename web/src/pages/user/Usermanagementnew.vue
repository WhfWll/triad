<template>
    <div>
        <div class="main-title  ">
            
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >用户管理</label>
        </div>
        <div class="userbox context_box_bg">
            <div class="search-box">
                <div class="operationbutton">
                    <xz-button type="primary" @click="AddUser" size="small">新建</xz-button>
                    <del-button :width="170" @click="AllDel" style="margin-left: 8px;"
                        :disabled="!multipleSelection.length  || (role === '2' && selectedID.includes(2)) || selectedID.includes(3) || selectedID.includes(4)">
                    </del-button>
                </div>
                <div class="serach-condition"> 
                    <div class="search-text">
                        <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select" size="small"
                        clearable>
                        </el-input>
                        <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                    </div>
                </div>
            </div>
            <el-table :data="tableData" style="width: 100%" v-loading="Loading"  height="calc(100% - 102px)"
                @selection-change="handleSelectionChange" @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column type="selection" :selectable='checkboxT' width="55">
                </el-table-column>
                <el-table-column prop="username" label="账号">
                </el-table-column>
                <el-table-column prop="role" label="身份">
                    <template slot-scope="scope">
                        <div>
                            <span class="" v-if="scope.row.role == '1'">普通用户</span>
                            <span class="" v-if="scope.row.role == '2'">管理员</span>
                            <span class="" v-if="scope.row.role == '3'">审计员</span>
                            <span class="" v-if="scope.row.role == '4'">超级管理员</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="groupStr" label="用户组">
                </el-table-column>
                <el-table-column prop="isAlive" label="在线状态">
                    <template slot-scope="scope">
                        <div>
                            <span class="" v-if="scope.row.isAlive == '1'">在线</span>
                            <span class="" v-if="scope.row.isAlive == '2'">离线</span>
                          
                        </div>
                    </template>
                </el-table-column>
                <!-- <el-table-column prop="email" label="邮箱">
                </el-table-column> -->
                <el-table-column prop="statusStr" label="账号状态">
                </el-table-column>
                <el-table-column prop="lastTime" label="最后在线时间">
                </el-table-column>
<!-- prop="status_str" label="账号状态" -->
                <el-table-column prop="accountExpire" label="账号有效期" width="200">
                     <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id  ">
                            <el-link class="link_primary" :underline="false" @click="updateUser(scope.row)"
                                :disabled="role === '2' && Number(role) === scope.row.role && (scope.row.id !== Number(userID))||(role =='2'&&scope.row.role=='4')">
                                编辑</el-link>
                            <el-link v-show="false"  class="link_primary" :underline="false"
                                @click="updatepwd(scope.row.id,scope.row.username)"
                                :disabled="role === '2' && Number(role) === scope.row.role && (scope.row.id !== Number(userID))">
                                密码</el-link>
                            <el-link v-if="scope.row.status==1" class="link_primary" :underline="false"
                             @click="setislogin(scope.row.id,scope.row.status)"
                             :disabled="role === '2' && Number(role) === scope.row.role||(role =='2'&&scope.row.role=='4')"
                             >禁用</el-link>
                            <el-link v-else 
                            @click="setislogin(scope.row.id,scope.row.status)" class="link_primary" :underline="false"
                            :disabled="role === '2' && Number(role) === scope.row.role ||(role =='2'&&scope.row.role=='4')"
                            >启用</el-link>
                            <el-popover placement="bottom" width="100" trigger="click" popper-class="learnMore"
                                :ref="`popover-${scope.row.id}`" :visible-arrow="false" style="padding:0"
                               v-show="!(role =='2'&&scope.row.role=='4')"
                                v-if="!((role === '2' && scope.row.role == 2) || (role === '4' && scope.row.role == 3) || (role === '4' && scope.row.role == 4) || (role === '2' && Number(role) === scope.row.role && (scope.row.id !== Number(userID))))">
                                <ul class="operationbox">
                                    <li v-if="role ==4 || role==2" @click="resetPas(scope.row.id)">重置密码</li>
                                    
                                    <li @click="useDate(scope.row.id,scope.row.accountExpire)">修改有效期</li>
                                    <li>
                                        <el-popover placement=" bottom" width="170" :visible-arrow="false"
                                            :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                                            <p class="delText">
                                                <i class="el-icon-warning"></i>确定删除吗？
                                            </p>
                                            <div style="text-align: right; margin: 0">
                                                <el-button size="mini" class="delCancel" @click="fncancel(scope)">取消
                                                </el-button>
                                                <el-button size="mini" type="primary" @click="fnDel(scope)">确定
                                                </el-button>
                                            </div>
                                            <span slot="reference">删除</span>
                                        </el-popover>
                                    </li>
                                </ul>
                                <el-link :underline="false" class="link_info" slot="reference"  >更多
                                </el-link>
                            </el-popover>

                        </div>
                        <div v-else>
                            {{scope.row.accountExpire}}
                            <!-- <span class="" v-if="scope.row.role == '1'">普通用户</span>
                            <span class="" v-if="scope.row.role == '2'">管理员</span>
                            <span class="" v-if="scope.row.role == '3'">审计员</span>
                            <span class="" v-if="scope.row.role == '4'">超级管理员</span> -->
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination :page-size="10" background layout=" total,  prev, pager, next,sizes, jumper"
                :total="totalpage" :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>


        </div>
        <el-dialog title="新建用户" :visible.sync="dialogaddFormVisible" :before-close="cancelform" width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="userform" label-width="0" status-icon ref="ruleFormadduser" :rules="rules1">
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">账号</label>
                        <el-input v-model="userform.username" size="small" style="width:320px" placeholder="请输入用户名"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label="" prop="role">
                        <label class="dialog_item_label">身份</label>
                        <el-select v-model="userform.role" size="small" style="width:320px" placeholder="请选择身份">
                            <el-option label="普通用户" value="1"></el-option>
                            <el-option label="管理员" value="2" v-if="role === '4'"></el-option>
                            <!-- <el-option label="审计员" value="3"></el-option> -->
                        </el-select>
                    </el-form-item>
                    <el-form-item label="" prop="groupIds">
                        <label class="dialog_item_label">用户组</label>
                        <el-cascader
                            v-model="userform.groupIds"
                            :options="groupOptions"
                            :props="{ multiple: true, emitPath: false, checkStrictly: true, expandTrigger: 'hover' }"
                            size="small"
                            style="width:320px"
                            placeholder="请选择用户组"
                            clearable
                            collapse-tags>
                        </el-cascader>
                    </el-form-item>
                    <el-form-item label=" " prop="password">
                        <label class="dialog_item_label">密码</label>
                        <el-input type="password" v-model="userform.password" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入新密码"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="repassword">
                        <label class="dialog_item_label">确认密码</label>
                        <el-input type="password" v-model="userform.repassword" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入确认密码"></el-input>
                    </el-form-item>
                        <el-form-item label=" " prop="accountExpireTime">
                        <label class="dialog_item_label">账号有效期</label>
                            <div class="userblock">
                                <el-date-picker
                                v-model="userform.accountExpireTime"
                                type="date"
                                format="yyyy-MM-dd" value-format="yyyy-MM-dd"
                                placeholder="选择日期" :picker-options="expireTimeOption">
                                </el-date-picker>
                            </div>
                    </el-form-item>
                    <el-form-item label=" " prop="email">
                        <label class="dialog_item_label">邮箱</label>
                        <el-input v-model="userform.email" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入邮箱"></el-input>
                    </el-form-item>

                

                    <el-form-item label="" prop="department">
                        <label class="dialog_item_label">部门</label>
                        <el-input v-model="userform.department" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入部门"></el-input>
                    </el-form-item>
                    <el-form-item label="" prop="remark">
                        <label class="dialog_item_label" style="vertical-align: top; margin-top: 7px; ">备注</label>
                        <el-input type="textarea" class="txtareacontent" v-model="userform.remark" size="small"
                            style="width:320px; " autocomplete="off" placeholder="请输入备注"></el-input>
                    </el-form-item>
                </el-form>
            </div>
        </el-dialog>
        <el-dialog title="编辑用户" width="1184px" :visible.sync="dialogupdateFormVisible" :before-close="updatecancelform"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="udpatesubmitForm">确定</el-button>
                <el-button size="small" @click="updatecancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="userform1" label-width="0" status-icon ref="ruleFormupdateuser" :rules="rules2">
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">账号</label>
                        <el-input disabled v-model="userform1.username"
                            size="small" style="width:320px" placeholder="请输入用户名" maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="role">
                        <label class="dialog_item_label">身份</label>
                        <el-select
                            :disabled="userform1.role == 3 && role === '4' || userform1.role == 4 && role === '4' || role === '2' "
                            v-model="userform1.role" size="small" style="width:320px" placeholder="请选择身份">
                            <el-option v-for="item in roleOptions" :key="item.value" :label="item.label"
                                :value="item.value"></el-option>
                            <!-- <el-option label="普通用户" value="1"></el-option>
                                <el-option label="管理员" value="2"></el-option>
                                <el-option label="审计员" value="3" v-if="role === '4' && userform1.role == 4"></el-option>
                                <el-option label="超级管理员" value="4" v-if="role === '4' && userform1.role == 4"></el-option> -->
                        </el-select>
                    </el-form-item>
                    <el-form-item label="" prop="groupIds">
                        <label class="dialog_item_label">用户组</label>
                        <el-cascader
                            v-model="userform1.groupIds"
                            :options="groupOptions"
                            :props="{ multiple: true, emitPath: false, checkStrictly: true, expandTrigger: 'hover' }"
                            size="small"
                            style="width:320px"
                            placeholder="请选择用户组"
                            clearable
                            collapse-tags>
                        </el-cascader>
                    </el-form-item>
                    <el-form-item label=" " prop="email">
                        <label class="dialog_item_label">邮箱</label>
                        <el-input v-model="userform1.email" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入邮箱"></el-input>
                    </el-form-item>
                    <el-form-item label="" prop="department">
                        <label class="dialog_item_label">部门</label>
                        <el-input v-model="userform1.department" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入部门"></el-input>
                    </el-form-item>
                    <el-form-item label="" prop="remark">
                        <label class="dialog_item_label" style="vertical-align: top; margin-top: 7px; ">备注</label>
                        <el-input type="textarea" class="txtareacontent" v-model="userform1.remark" size="small"
                            style="width:320px; " autocomplete="off" placeholder="请输入备注"></el-input>
                    </el-form-item>
                </el-form>
            </div>

        </el-dialog>

        <el-dialog title="修改密码" :visible.sync="dialogVisiblepwd" class="updatepwdbox newUserDialog"
            :before-close="pwdcancelform" width="1184px" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="pwdsubmitForm">确定</el-button>
                <el-button size="small" @click="pwdcancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="pwddata" label-width="0" status-icon ref="ruleFormPWD" :rules="rulespwd">
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">用户名</label>
                        <el-input v-model="pwddata.username" size="small" style="width:320px" placeholder="请输入用户名"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="password">
                        <label class="dialog_item_label">新密码</label>
                        <el-input type="password" v-model="pwddata.password" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入新密码"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="repassword">
                        <label class="dialog_item_label">确认密码</label>
                        <el-input type="password" v-model="pwddata.repassword" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入确认密码"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="admin_password">
                        <label class="dialog_item_label">管理员密码</label>
                        <el-input type="password" v-model="pwddata.admin_password" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入管理员密码"></el-input>
                    </el-form-item>
                </el-form>
            </div>
        </el-dialog>

        <!-- //有效期弹窗。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。 -->
        <el-dialog title="修改账号有效期" :visible.sync="dialogUsedate" class="updatepwdbox newUserDialog"
            :before-close="pwdcancelform" width="984px" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="usedatesubmit">修改有效期</el-button>
                <el-button size="small" @click="usedatecancel">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="pwddata" label-width="0" status-icon ref="ruleFormPWD" :rules="rulespwd">
                    <el-form-item label=" " prop="">
                        <label class="dialog_item_label">账号有效期</label>
                            <div class="userblock">
                                <el-date-picker
                                v-model="valuedate2"
                                type="date"
                                format="yyyy-MM-dd" value-format="yyyy-MM-dd"
                                placeholder="选择有效期截止日期" :picker-options="expireTimeOption" >
                                </el-date-picker>
                            </div>
                    </el-form-item>
                </el-form>
            </div>

        </el-dialog>
        <!-- //重置密码 -->
        <el-dialog title="重置密码" :visible.sync="dialogResetPass" class="updatepwdbox newUserDialog"
            :before-close="pwdcancelform2" width="984px" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="usedatesubmit2">重置密码</el-button>
                <el-button size="small" @click="usedatecancel2">取消</el-button>
            </div>
             <div style="padding:24px">
                <el-form :model="demoData" label-width="0" status-icon ref="ruleFormPWD" :rules="rulespwd" >
                 
                    <el-form-item label=" " prop="admin_password">
                        <label class="dialog_item_label">管理员密码</label>
                        <el-input min-length="8" max-length="60" type="password" v-model="demoData.admin_password" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入管理员密码"></el-input>
                    </el-form-item>
                </el-form>
            </div>

        </el-dialog>
    </div>
</template>
<style scoped lang="less">

.u-txt{
    display: inline;
}
.u-btn {
    display: none;

}

.el-table__body tr:hover {

    .u-btn {

        display: inline;

    }
        .u-txt{
            display: none;
        }

}

.learnMore {
    min-width: 100px !important;
    padding: 16px 0 !important;
}

.userbox{
	padding: 24px; 
    background: #fff;
    height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
}	
.txtareacontent /deep/ textarea{
	resize: none !important;
}
.newUserDialog /deep/ .el-form-item__label{
    position: absolute;
    left: 127px;
}
.newUserDialog /deep/ .el-form-item__error{
    left: 112px;
}
.popover_block{
    list-style: none;
    li{
        text-align: center;
        margin-top: 8px;
    }
} 
.userblock{
    display: inline-block;
}
/deep/ .userblock .el-date-editor.el-input, .el-date-editor.el-input__inner {
    width: 320px!important;  
}
/deep/ .el-popconfirm{
    .el-popconfirm__main{ 
            margin-bottom: 16px;
            color: rgba(72, 72, 102, .64);
   
    }
}
</style>
<script>   
import XzButton from "../../components/XzButton.vue"; 
import DelButton from "../../components/DelButton.vue"; 
import user from "@/api/user.js"; 
export default ({
    name:'usermanagement',
    components: {
    	XzButton,
		DelButton, 
  	},
    data(){ 
        var validatePass2 = (rule, value, callback) => {
            if (value === '') {
                callback(new Error('请输入确认密码'));
            } else if (value !== this.userform.password) {
                callback(new Error('两次输入密码不一致!'));
            } else {
                callback();
            }
        }; 
        var validatePass1 = (rule, value, callback) => {
            if (value === '') {
                callback(new Error('请输入确认密码'));
            } else if (value !== this.pwddata.password) {
                callback(new Error('两次输入密码不一致!'));
            } else {
                callback();
            }
        };
        var validatePwd =  (rule, value, callback) => {
            const reg = /^\S+$/; 
            if (reg.test(value)) {
                callback();
            } else {
                return callback(new Error('格式不正确(不能包含空格)'));
            }
        }; 
    	return{  
            formData:{
                search:'',
            },
            account_expireZ:'',
            resetPasID:'',
            resetPassword:'',
            demoData:{
                admin_password:'',
            },
            dialogResetPass:false,
            updateExpID:'',//修改有效期的id
            valuedate: '2099-12-30',
            valuedate2: '2099-12-30',
            isMore:false,
            showOperateButton:false,
            rowId:'',
    		multipleSelection:[], 
            currentpage:1,
            totalpage:0,
            page_num:1,
            tableData:[],
            Loading:false,
            groupOptions: [],
            userform:{
            	username:'',
            	role:'1',
                groupIds:[],
            	password:'',
            	repassword:'',
            	email:'',
            	department:'',
            	remark:'',
                accountExpireTime:''
            },
            userform1:{
                id:'',
            	username:'',
            	role:'1', 
                groupIds:[],
            	email:'',
            	department:'',
            	remark:'',
            }, 
            dialogaddFormVisible:false,
            dialogupdateFormVisible:false,
            dialogUsedate:false,
            rules1:{
                username: [
                    { required: true, message: '用户名不能为空', trigger: 'blur' }, 
                    { max: 50, message: '用户名最大长度不能超过50', trigger: 'blur' },
                ],
                role:[
                     { required: true, message: '请选择身份', trigger: 'change' }
                ],
                password: [
                    { required: true, message: '新密码不能为空', trigger: 'blur' },
                     { min: 8, max: 20, message: '密码长度8-20字符', trigger: 'blur' },
                    { validator: validatePwd, trigger: 'blur' } 
                ],
                repassword: [
                    { required: true, message: '确认密码不能为空', trigger: 'blur' }, 
                     { min: 8, max: 20, message: '密码长度8-20字符', trigger: 'blur' },
                    { validator: validatePass2, trigger: 'blur' }
                ],
                email: [
                    { required: true, message: '邮箱不能为空', trigger: 'blur' },
                    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur'] }
                ],
                accountExpireTime: [
                    { required: true, message: '账号有效期不能为空', trigger: 'blur' },
                ],
            },
            rules2: {
                username: [
                    { required: true, message: '用户名不能为空', trigger: 'blur' }, 
                    { max: 50, message: '用户名最大长度不能超过50', trigger: 'blur' },
                ],
                role:[
                     { required: true, message: '请选择身份', trigger: 'change' }
                ], 
                email: [
                    { required: true, message: '邮箱不能为空', trigger: 'blur' },
                    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
                ],
            },
            dialogVisiblepwd:false,
            pwddata:{
                id:'',
            	username:'', 
            	password:'',
            	repassword:'',
            	'admin_password':'',
            },
            adduser:true,
            validatePass:validatePass2,
            rulespwd:{
                username: [
                    { required: true, message: '用户名不能为空', trigger: 'blur' }, 
                ], 
                password: [
                    { required: true, message: '新密码不能为空', trigger: 'blur' }, 
                ],
                repassword: [
                    { required: true, message: '确认密码不能为空', trigger: 'blur' }, 
                    { validator: validatePass1, trigger: 'blur' }
                ],
                'admin_password':[
                    { required: true, message: '管理员密码不能为空', trigger: 'blur' }, 
                ]
            },
            is_audit:false, //是否审计员
            alldelvisible:false,
            user:'',
            role:1,  //1  普通用户   2  管理员   3  审核员
            selectedID: [], // 多选被选中的role数组
            userID: '',
            roleOptions: [
                {
                    value: 1,
                    label: '普通用户'
                },
                {
                    value: 2,
                    label: '管理员'
                },
                {
                    value: 3,
                    label: '审计员'
                },
                {
                    value: 4,
                    label: '超级管理员'
                }
            ],
            columns:[{
                    prop: "username",
                    label: "用户名"
                },{
                    prop:'email',
                    label:'邮箱',
                }, {
                prop: 'last_login',
                    label: '最后登录时间',
                }
            ],
            isSHow:false,
            looked:false,
            pageSize:10,
            //设置失效日期今天之前的日期不可选
            expireTimeOption: {
                disabledDate(date) {
                    //disabledDate 文档上：设置禁用状态，参数为当前日期，要求返回 Boolean
                    return date.getTime() < Date.now() - 24 * 60 * 60 * 1000;
                    //   return time.getTime() > Date.now();
                },
            },
    	}
    },
    watch: {
        'userform1.role': function(value) {
            if (value === 1 || value === 2) {
                this.roleOptions = [
                {
                    value: 1,
                    label: '普通用户'
                },
                {
                    value: 2,
                    label: '管理员'
                }
                ]
            } else {
                this.roleOptions = [
                {
                    value: 1,
                    label: '普通用户'
                },
                {
                    value: 2,
                    label: '管理员'
                },
                {
                    value: 3,
                    label: '审计员'
                },
                {
                    value: 4,
                    label: '超级管理员'
                }
            ]
            }
        }
    },
    created:function(){
        this.$store.state.activefirstMenu="/usermanagement"; 
        this.role = this.commonjs.decryptCBC(localStorage.getItem('role'),this.commonjs.myKey); 
        this.user = this.commonjs.decryptCBC(localStorage.getItem('user'),this.commonjs.myKey); 
        this.userID = this.commonjs.decryptCBC(localStorage.getItem('user_id'),this.commonjs.myKey); 
        // this.userID = this.commonjs.decryptCBC(localStorage.getItem('id'),this.commonjs.myKey); 
        this.pageSize = this.commonjs.pageSize;
    },  
    mounted:function(){  
        this.getData(); //获得列表数据
        this.getGroupList();
    },
    methods:{
        async getGroupList() {
            const data = await user.getupgroup();
            if (data.code === 200) {
                this.groupOptions = this.cleanChildren(data.data.list);
            }
        },
        cleanChildren(list) {
            list.forEach(item => {
                if (item.children === null || item.children.length === 0) {
                    delete item.children;
                } else {
                    this.cleanChildren(item.children);
                }
            });
            return list;
        },
        checkboxT(row,index){
            if(row.role== '3' || row.role== '4'){
                return 0;
            }else{
                return 1;
            }
        },
    	goBack() {
	       this.$router.go(-1);
	    }, 
        async getData(){
            this.Loading = true;   
            const data = await user.userList({
                page:this.page_num,
                size:this.pageSize,  
                search:this.formData.search
            }); 
            this.Loading = false;
            this.tableData = data.data.list;
            this.totalpage = data.data.total ; 
        }, 
        handlesearch(){
              //搜索
            this.page_num = 1;
            this.getData();
            this.currentpage = 1;
        },
        handleReset(){
            this.page_num = 1;
            this.formData.search = ""; 
            this.pageSize = 10;
            this.currentpage = 1;
            this.getData();
        },
        async AllDel(){
            console.log(this.multipleSelection,'this.multipleSelection');
        	if(this.multipleSelection.length == 0) return;  
            var ids = this.multipleSelection.map(item => item.id); 
            const data = await user.userMultDelete({userIds:ids.join(',')});
            if(data.code == 200){
                this.$message({
                    message:data.msg || '删除成功',
                    type: 'success'
                });
                this.alldelvisible = false;
                this.getData();
            }else{
                this.$message({
                    message:data.msg,
                    type: 'error'
                });
            } 


        },
        fncancel(scope){
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
            scope._self.$refs[`popover-${scope.row.id}`].doClose()
        },
        async fnDel(scope){ //单个删除    
        console.log(scope,'scope');
            const data = await user.userDel({userIds:scope.row.id+''}); 
            if(data.code == 200){
                this.$message({
                    message:data.msg||'删除成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
                scope._self.$refs[`popover-${scope.row.id}`].doClose()
                this.getData();
            }else{
                this.$message({
                    message:data.msg,
                    type: 'error'
                });
            } 
	         
    	},
        handleSelectionChange(val){
        	this.multipleSelection = val;
            this.selectedID = []
            this.multipleSelection.forEach(item => {
                this.selectedID.push(item.role)
            })
        },
        currentchange(t){
        	this.page_num = t; 
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t) {
            this.page_num = 1;
            this.pageSize = t;
            this.getData();
        },
        cancelform(){
            this.userform.username = '';
            this.userform.role = '1';
            this.userform.groupIds = [];
            this.userform.password = '';
            this.userform.repassword = '';
            this.userform.email = '';
            this.userform.department = '';
            this.userform.remark = '';
            this.userform.accountExpireTime = '';
            this.dialogaddFormVisible = false; 
            this.$refs.ruleFormadduser.resetFields(); 
        },
        updatecancelform(){
            this.userform1.username = '';
            this.userform1.role = 1; 
            this.userform1.groupIds = [];
            this.userform1.email = '';
            this.userform1.department = '';
            this.userform1.remark = '';
            this.dialogupdateFormVisible = false; 
            this.$refs.ruleFormupdateuser.resetFields(); 
        },
        submitForm(){ 
            this.$refs.ruleFormadduser.validate( async (valid) => {
                if (valid) { 
                 this.userform.password = this.$commonjs.encryptCBC( this.userform.password, this.$commonjs.myKey);
                 this.userform.repassword= this.$commonjs.encryptCBC(this.userform.repassword, this.$commonjs.myKey);
                 
                 let params = { ...this.userform };
                 params.id = 0;
                 if (Array.isArray(params.groupIds)) {
                     params.groupIds = params.groupIds.join(',');
                 }
                    
                    const data = await  user.addUser(params); 
                    // console.log(data,'添加用户data');
                    if(data.code == 200){ 
                        this.$message({
                            message:data.msg || '添加成功',
                            type: 'success'
                        });
                        this.dialogaddFormVisible = false;
                        this.getData();
                    }else{
                        this.$message({
                            message:data.msg,
                            type: 'error'
                        });
                    }   
                   
                    // this.$refs.ruleFormuser.resetFields();  
                }
            }); 
            
        }, 
        udpatesubmitForm(){
             this.$refs.ruleFormupdateuser.validate(async (valid) => {
                if (valid) {
                    let userform1 = _.cloneDeep(this.userform1);
                    if (Array.isArray(userform1.groupIds)) {
                        userform1.groupIds = userform1.groupIds.join(',');
                    }
                    const data = await user.updateUser(userform1);
                    if(data.code == 200){
                        this.$message({
                            message:data.msg || '修改成功',
                            type: 'success'
                        });
                        this.dialogupdateFormVisible = false;
                        this.getData();
                    }else{
                        this.$message({
                            message:data.msg,
                            type: 'error'
                        });
                    }   
                    
                    // this.$refs.ruleFormuser.resetFields(); 
                }
            });
        },
        AddUser(){  
        	this.dialogaddFormVisible = true; 
            this.userform.username = '';
            this.userform.role = '1';
            this.userform.groupIds = [];
            this.userform.password = '';
            this.userform.repassword = '';
            this.userform.email = '';
            this.userform.department = '';
            this.userform.remark = '';
            this.userform.accountExpireTime = '';
        },
        async updateUser(rows){

            this.userform1.id = rows.id;
            this.userform1.username = rows.username;
            this.userform1.email = rows.email;
            this.userform1.department = rows.department;
            this.userform1.remark = rows.remark;
            this.userform1.role = Number(rows.role)
            // 优先使用列表中已有的 groupIds 数据
            if (rows.groupIds && Array.isArray(rows.groupIds)) {
                this.userform1.groupIds = rows.groupIds;
            } else if (rows.groupIds && typeof rows.groupIds === 'string') {
                this.userform1.groupIds = rows.groupIds.split(',').filter(id => id).map(Number);
            } else {
                this.userform1.groupIds = [];
            }
            
            // console.log(this.userform);
            if(rows.role == '3'){
                this.is_audit = true;
            }else{
                this.is_audit = false;
            }

            // 获取用户详情以填充用户组（如果列表数据不全，可以从这里补充，但不应覆盖为空）
            const res = await user.userInfo({ id: rows.id });
            if (res.code === 200 && res.data) {
                 // 假设后端返回的是逗号分隔的字符串或数组，这里做兼容处理
                 let gIds = res.data.groupIds || res.data.group_ids;

                 // 如果没有 direct IDs，检查是否有 groups 对象数组
                 if (!gIds && res.data.groups && Array.isArray(res.data.groups)) {
                     gIds = res.data.groups.map(g => g.id);
                 }

                 if (gIds) {
                     if (typeof gIds === 'string') {
                         this.userform1.groupIds = gIds.split(',').filter(id => id).map(Number); // 转为数字数组
                     } else if (Array.isArray(gIds)) {
                         this.userform1.groupIds = gIds;
                     }
                 } 
            }

        	this.dialogupdateFormVisible = true; 
        },
        pwdcancelform(){ 
            this.pwddata.id = '';
            this.pwddata.username = '';
            this.pwddata.password = '';
            this.pwddata.repassword = '';
            this.pwddata.admin_password = ''; 
        	this.dialogVisiblepwd = false;
            this.$refs.ruleFormPWD.resetFields(); 
        },
        pwdcancelform2(){ 
           
        },
        pwdsubmitForm(){ //保存修改密码
            this.$refs.ruleFormPWD.validate( async (valid) => {
                if (valid) {
                    const data = await user.updatePwd(this.pwddata); 
                    if(data.code == 200){
                        this.$message({
                            message:data.msg || '修改成功',
                            type: 'success'
                        });
                    }else{
                        this.$message({
                            message:data.msg,
                            type: 'error'
                        });
                    }    
                    this.dialogVisiblepwd = false;
                    this.getData();
                    this.pwddata.id = '';
                    this.pwddata.username = '';
                    this.pwddata.password = '';
                    this.pwddata.repassword = '';
                    this.pwddata.admin_password  = ''; 
                    // this.$refs.ruleFormPWD.resetFields();  
                }
            });
        },
        updatepwd(id,username){
            this.pwddata.username = username;
            this.pwddata.id = id;
        	this.dialogVisiblepwd = true;
        },
        async setislogin(userid,is_active){ //设置禁用启用  
            const data = await user.setIsLoginUser({
                status:is_active==1 ? 5:1,
                userId : userid 
            });  
            if(data.code == 200){
                this.$message({
                    message:data.msg || '设置成功',
                    type: 'success'
                });
                this.getData();
            }else{
                this.$message({
                    message:data.msg,
                    type: 'error'
                });
            }             
            
        },
// 账号有效期..........................................................................
        useDate(id,accountExpire){
            this.dialogUsedate = true;
            this.updateExpID = id;
            this.account_expireZ = accountExpire
            this.valuedate2 =this.account_expireZ
        },
        async  usedatesubmit(){
             const res =  await user.updateExp({
                userId: this.updateExpID ,
                accountExpireTime:this.valuedate2
            })
            if(res.code == 200){
                this.$message({
                    message:res.msg || '设置成功',
                    type: 'success'
                });
                this.getData();
                 this.dialogUsedate = false;
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
                 this.dialogUsedate = false;
            }
        },
        usedatecancel(){
            this.dialogUsedate = false;
        },
//--重置密码--------------------------------------------------------------------------------------------
       resetPas(id){
        this.dialogResetPass = true;
        this.resetPasID = id;
       },
         async usedatesubmit2(){
           console.log('871-重置密码');
          const res = await user.resetPwd({
            // userId:this.$commonjs.encryptCBC(this.resetPasID, this.$commonjs.myKey),
            userId:this.resetPasID,
            password:this.$commonjs.encryptCBC(this.demoData.admin_password, this.$commonjs.myKey)
            // password:this.demoData.admin_password
          });
          if(res.code == 200){
                this.$message({
                 message:res.msg || '重置成功',
                 type: 'success'
                });
                this.dialogResetPass = false;
                this.demoData.admin_password = '';
                this.getData();
          }else{
                this.$message({
                 message:res.msg,
                 type: 'error'
                });
                this.dialogResetPass = false;
          }
         },
           usedatecancel2(){
            this.dialogResetPass = false;
            this.demoData.admin_password = '';
        },
       mouseenter(row,colum,cell,event){   
            this.showOperateButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
            if (!this.$refs['popover-' + row.id]){
                this.showOperateButton = false;
                this.rowId = "";
                return;
            }else{
                let isShow = this.$refs['popover-' + row.id].showPopper; 
                if (!isShow) {
                    this.showOperateButton = false;
                    this.rowId = "";
                }

            }

           
        },
    
    }
})
 
</script>
