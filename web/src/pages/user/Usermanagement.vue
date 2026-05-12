<template>
    <div>
        <div class="main-title  ">
            用户管理
            当前身份为：
            <span class="" v-if="role == '1'">普通用户</span>
            <span class="" v-if="role == '2'">管理员</span>
            <span class="" v-if="role == '3'">审计员</span>
            <span class="" v-if="role == '4'">超级管理员</span>
        </div>
        <div class="userbox context_box_bg">
            <div class="search-box">
                <div class="operationbutton">
                    <xz-button type="primary" @click="AddUser" size="small">新建</xz-button>
                    <del-button :width="170" @click="AllDel" style="margin-left: 8px;"
                        :disabled="!multipleSelection.length || (role === '2' && selectedID.includes(2)) || selectedID.includes(3) || selectedID.includes(4)">
                    </del-button>
                </div> 
            </div>
            <el-table :data="tableData" style="width: 100%" v-loading="Loading" @selection-change="handleSelectionChange"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
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
                <el-table-column prop="group_str" label="用户组">
                </el-table-column>
                <el-table-column prop="email" label="邮箱">
                </el-table-column>
                <el-table-column prop="account_expire" label="有效期">
                </el-table-column>
                <el-table-column prop="last_time" label="最后在线">
                </el-table-column>

                <el-table-column prop="status_str" label="状态" width="200">
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                            <el-link class="link_primary" :underline="false" @click="updateUser(scope.row)"
                                :disabled="(role === '2' && scope.row.role !== 1 && scope.row.role !== 2) || (role === '2' && Number(role) === scope.row.role && (scope.row.id !== Number(userID)))">
                                编辑</el-link>
                            <el-link class="link_primary" :underline="false"
                                @click="updatepwd(scope.row.id, scope.row.username)"
                                :disabled="(role === '2' && scope.row.role !== 1 && scope.row.role !== 2) || (role === '2' && Number(role) === scope.row.role && (scope.row.id !== Number(userID)))">
                                密码</el-link>

                            <el-popover placement="bottom" width="60" trigger="click" popper-class="learnMore"
                                :ref="`popover-${scope.row.id}`" :visible-arrow="false" style="padding:0"
                                v-if="!((role === '2' && scope.row.role == 2) || (role === '4' && scope.row.role == 3) || (role === '4' && scope.row.role == 4) || (role === '2' && Number(role) === scope.row.role && (scope.row.user_id !== Number(userID))) || (role === '2' && scope.row.role !== 1 && scope.row.role !== 2))">
                                <ul class="operationbox">
                                    <li :disabled="(role === '2' && scope.row.role == 2) || (role === '4' && scope.row.role == 3) || (role === '4' && scope.row.role == 4) || (role === '2' && Number(role) === scope.row.role && (scope.row.user_id !== Number(userID)))"
                                        v-if="scope.row.status === 1" @click="setislogin(scope.row.id, scope.row.status)">禁用
                                    </li>
                                    <li v-else @click="setislogin(scope.row.id, scope.row.status)">启用</li>
                                    <li @click="useDate(scope.row.id, scope.row.account_expire)">有效期</li>
                                    <li>
                                        <el-popover placement=" bottom" width="170" :visible-arrow="false"
                                            :disabled="(role === '2' && scope.row.role == 2) || (role === '4' && scope.row.role == 3) || (role === '4' && scope.row.role == 4) || (role === '2' && Number(role) === scope.row.role && (scope.row.user_id !== Number(userID)))"
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
                                            <span slot="reference"
                                                :disabled="(role === '2' && scope.row.role == 2) || (role === '4' && scope.row.role == 3) || (role === '4' && scope.row.role == 4) || (role === '2' && Number(role) === scope.row.role && (scope.row.user_id !== Number(userID)))">删除</span>
                                        </el-popover>
                                    </li>
                                </ul>
                                <el-link :underline="false" class="link_info" slot="reference">更多
                                </el-link>
                            </el-popover>

                        </div>
                        <div v-else>
                            <span class="">{{ scope.row.status_str }}</span>

                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination :page-size="10" background layout=" total,  prev, pager, next,sizes, jumper" :total="totalpage"
                :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
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
                    <el-form-item label=" " prop="role">
                        <label class="dialog_item_label">身份</label>
                        <el-select v-model="userform.role" size="small" style="width:320px" placeholder="请选择身份">
                            <el-option label="普通用户" value="1"></el-option>
                            <el-option label="管理员" value="2" v-if="role === '4'"></el-option>
                            <!-- <el-option label="审计员" value="3"></el-option> -->
                        </el-select>
                    </el-form-item>
                    <el-form-item label=" " prop="group_ids">
                        <label class="dialog_item_label">用户组</label>
                        <el-cascader style="width:320px" v-model="userform.group_ids" :options="groupOptions" :props="props"
                            placeholder="请选择用户组"></el-cascader>
                        <!-- <el-input v-model="userform.group_ids" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入部门"></el-input> -->
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
                    <el-form-item label=" " prop="email">
                        <label class="dialog_item_label">邮箱</label>
                        <el-input v-model="userform.email" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入邮箱"></el-input>
                    </el-form-item>

                    <el-form-item label=" " prop="account_expire">
                        <label class="dialog_item_label">账号有效期</label>
                        <div class="userblock">
                            <el-date-picker v-model="userform.account_expire" type="date" value-format="yyyy-MM-dd"
                                format="yyyy-MM-dd" placeholder="选择日期" :picker-options="expireTimeOption">
                            </el-date-picker>
                        </div>
                    </el-form-item>

                    <el-form-item label="">
                        <label class="dialog_item_label">部门</label>
                        <el-input v-model="userform.department" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入部门"></el-input>
                    </el-form-item>

                    <el-form-item label="">
                        <label class="dialog_item_label" style="vertical-align: top; margin-top: 7px; ">备注</label>
                        <el-input type="textarea" class="txtareacontent" v-model="userform.remark" size="small"
                            style="width:320px; " autocomplete="off" placeholder="请输入备注"></el-input>
                    </el-form-item>
                </el-form>
            </div>
        </el-dialog>
        <el-dialog title="修改用户" width="1184px" :visible.sync="dialogupdateFormVisible" :before-close="updatecancelform"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="udpatesubmitForm">确定</el-button>
                <el-button size="small" @click="updatecancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="userform1" label-width="0" status-icon ref="ruleFormupdateuser" :rules="rules2">
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">用户名</label>
                        <el-input :disabled="Number(userform1.role) == 4 || Number(userform1.role) == 3"
                            v-model="userform1.username" size="small" style="width:320px" placeholder="请输入用户名"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="role">
                        <label class="dialog_item_label">身份</label>
                        <el-select
                            :disabled="Number(userform1.role) == 3 && role === '4' || Number(userform1.role) == 4 && role === '4' || role === '2'"
                            v-model="userform1.role" size="small" style="width:320px" placeholder="请选择身份">
                            <el-option v-for="item in roleOptions" :key="item.value" :label="item.label"
                                :value="item.value"></el-option>
                            <!-- <el-option label="普通用户" value="1"></el-option>
                                <el-option label="管理员" value="2"></el-option>
                                <el-option label="审计员" value="3" v-if="role === '4' && userform1.role == 4"></el-option>
                                <el-option label="超级管理员" value="4" v-if="role === '4' && userform1.role == 4"></el-option> -->
                        </el-select>
                    </el-form-item>
                    <el-form-item label=" " prop="group_ids">
                        <label class="dialog_item_label">用户组</label>
                        <!-- <el-input v-model="userform1.group_ids" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入部门"></el-input> -->
                        <el-cascader style="width:320px" v-model="userform1.group_ids" :options="groupOptions"
                            :props="props" @change="change1"></el-cascader>
                    </el-form-item>
                    <el-form-item label=" " prop="email">
                        <label class="dialog_item_label">邮箱</label>
                        <el-input v-model="userform1.email" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入邮箱"></el-input>
                    </el-form-item>
                    <el-form-item label="">
                        <label class="dialog_item_label">部门</label>
                        <el-input v-model="userform1.department" size="small" style="width:320px" autocomplete="off"
                            placeholder="请输入部门"></el-input>
                    </el-form-item>
                    <el-form-item label="">
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
        <el-dialog title="账号有效期" :visible.sync="dialogUsedate" class="updatepwdbox newUserDialog"
            :before-close="pwdcancelform" width="984px" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="usedatesubmit">确定</el-button>
                <el-button size="small" @click="usedatecancel">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="timeForm" label-width="0" status-icon>
                    <el-form-item label=" " prop="username">
                        <label class="dialog_item_label">账号有效期</label>
                        <div class="userblock">
                            <el-date-picker v-model="timeForm.valuedate" type="date" value-format="yyyy-MM-dd"
                                placeholder="选择日期" :picker-options="expireTimeOption">
                            </el-date-picker>
                        </div>
                    </el-form-item>
                </el-form>
            </div>

        </el-dialog>
    </div>
</template>
<style scoped lang="less">
.u-txt {
    display: inline;
}

.u-btn {
    display: none;

}

.el-table__body tr:hover {

    .u-btn {

        display: inline;

    }

    .u-txt {
        display: none;
    }

}

.learnMore {
    min-width: 100px !important;
    padding: 16px 0 !important;
}

.userbox {
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}

.txtareacontent /deep/ textarea {
    resize: none !important;
}

.newUserDialog /deep/ .el-form-item__label {
    position: absolute;
    left: 127px;
}

.newUserDialog /deep/ .el-form-item__error {
    left: 112px;
}

.popover_block {
    list-style: none;

    li {
        text-align: center;
        margin-top: 8px;
    }
}

.userblock {
    display: inline-block;
}

/deep/ .userblock .el-date-editor.el-input,
.el-date-editor.el-input__inner {
    width: 320px !important;
}

/deep/ .el-popconfirm {
    .el-popconfirm__main {
        margin-bottom: 16px;
        color: rgba(72, 72, 102, .64);

    }
}
</style>
<script>
import XzButton from "../../components/XzButton.vue";
import DelButton from "../../components/DelButton.vue";
import user from "@/api/user.js";
import _ from 'lodash'
export default ({
    name: 'usermanagement',
    components: {
        XzButton,
        DelButton,
    },
    data() {
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
        var validatePwd = (rule, value, callback) => {
            const reg = /^\S+$/;
            if (reg.test(value)) {
                callback();
            } else {
                return callback(new Error('格式不正确(不能包含空格)'));
            }
        };
        return {
            props: { multiple: true },
            groupOptions: [],
            value: [],
            options: [{}],
            isMore: false,
            showOperateButton: false,
            rowId: '',
            multipleSelection: [],
            currentpage: 1,
            totalpage: 0,
            page_num: 1,
            tableData: [],
            Loading: false,
            userform: {
                username: '',
                role: '1',
                password: '',
                repassword: '',
                email: '',
                department: '',
                remark: '',
                account_expire: '2099-12-31',
                group_ids: []
            },
            userform1: {
                id: '',
                username: '',
                role: '',
                email: '',
                department: '',
                remark: '',
                group_ids: []
            },
            dialogaddFormVisible: false,
            dialogupdateFormVisible: false,
            dialogUsedate: false,
            rules1: {
                username: [
                    { required: true, message: '用户名不能为空', trigger: 'blur' },
                    { max: 50, message: '用户名最大长度不能超过50', trigger: 'blur' },
                ],
                role: [
                    { required: true, message: '请选择身份', trigger: 'change' }
                ],
                account_expire: [
                    { type: 'string', required: true, message: '请选择账号有效期', trigger: 'change' }
                ],
                password: [
                    { required: true, message: '新密码不能为空', trigger: 'blur' },
                    { min: 8, max: 20, message: '密码长度8-20字符', trigger: 'blur' },
                    { validator: validatePwd, trigger: 'blur' }
                ],
                repassword: [
                    { required: true, message: '确认密码不能为空', trigger: 'blur' },
                    { min: 8, max: 20, message: '密码长度8-20字符', trigger: 'blur' },
                    { validator: this.validatePass, trigger: 'blur' }
                ],
                email: [
                    { required: true, message: '邮箱不能为空', trigger: 'blur' },
                    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur'] }
                ],
                group_ids: [
                    { required: true, message: '备注不能为空', trigger: 'blur' },
                ]
            },
            rules2: {
                username: [
                    { required: true, message: '用户名不能为空', trigger: 'blur' },
                    { max: 50, message: '用户名最大长度不能超过50', trigger: 'blur' },
                ],
                role: [
                    { required: true, message: '请选择身份', trigger: 'change' }
                ],
                email: [
                    { required: true, message: '邮箱不能为空', trigger: 'blur' },
                    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
                ],
                group_ids: [
                    { required: true, message: '备注不能为空', trigger: 'blur' },
                ]
            },
            dialogVisiblepwd: false,
            pwddata: {
                id: '',
                username: '',
                password: '',
                repassword: '',
                admin_password: '',
            },
            timeForm: {
                user_id: '',
                valuedate: ''
            },
            adduser: true,
            validatePass: validatePass2,
            rulespwd: {
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
                'admin_password': [
                    { required: true, message: '管理员密码不能为空', trigger: 'blur' },
                ]
            },
            is_audit: false, //是否审计员
            alldelvisible: false,
            user: '',
            role: 1,  //1  普通用户   2  管理员   3  审核员
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
            columns: [{
                prop: "username",
                label: "用户名"
            }, {
                prop: 'email',
                label: '邮箱',
            }, {
                prop: 'last_login',
                label: '最后登录时间',
            }
            ],
            isSHow: false,
            looked: false,
            pageSize: 10,
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
        'userform1.role': function (value) {
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
    created: function () {
        this.$store.state.activefirstMenu = "/usermanagement";
        this.role = this.commonjs.decryptCBC(localStorage.getItem('role'), this.commonjs.myKey);
        this.user = this.commonjs.decryptCBC(localStorage.getItem('user'), this.commonjs.myKey);
        this.userID = this.commonjs.decryptCBC(localStorage.getItem('user_id'), this.commonjs.myKey);
        this.pageSize = this.commonjs.pageSize;
    },
    mounted: function () {
        // this.getData(); //获得列表数据
    },
    methods: {
        checkboxT(row, index) {
            if (row.role == '3' || row.role == '4') {
                return 0;
            } else {
                return 1;
            }
        },
        goBack() {
            this.$router.go(-1);
        },
        async getData() {
            this.Loading = false;
            const dt = await user.userList({
                page: this.page_num,
                size: this.pageSize,
                search: this.formData.search,  
            });
            if (dt.code === 200) {
                this.tableData = dt.data.list;
                this.totalpage = dt.data.total;
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        
        async AllDel() {
            if (this.multipleSelection.length == 0) return;
            var ids = this.multipleSelection.map(item => item.id);
            let params = {
                ids: ids.join(','),
                status: 0
            }
            const data = await user.userDel(params);
            if (data.code === 200) {
                this.$message({
                    message: data.msg,
                    type: 'success'
                });
                this.alldelvisible = false;
                this.getData();
            } else {
                this.$message({
                    message: data.msg,
                    type: 'error'
                });
            }


        },
        fncancel(scope) {
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
            scope._self.$refs[`popover-${scope.row.id}`].doClose()
        },
        async fnDel(scope) { //单个删除
            let params = {
                ids: scope.row.id,
                status: 0
            }
            const data = await user.userDel(params);
            if (data.code === 200) {
                this.$message({
                    message: data.msg,
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
                scope._self.$refs[`popover-${scope.row.id}`].doClose()
                this.getData();
            } else {
                this.$message({
                    message: data.msg,
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
        cancelform() {
            this.userform.username = '';
            this.userform.role = '1';
            this.userform.password = '';
            this.userform.repassword = '';
            this.userform.email = '';
            this.userform.department = '';
            this.userform.remark = '';
            this.dialogaddFormVisible = false;
            this.$refs.ruleFormadduser.resetFields();
        },
        updatecancelform() {
            this.userform1.username = '';
            this.userform1.role = '1';
            this.userform1.email = '';
            this.userform1.department = '';
            this.userform1.remark = '';
            this.dialogupdateFormVisible = false;
            this.$refs.ruleFormupdateuser.resetFields();
        },
        submitForm() {
            this.$refs.ruleFormadduser.validate(async (valid) => {
                if (valid) {
                    let userform = _.cloneDeep(this.userform)
                    userform.group_ids = userform.group_ids.join(',')
                    const data = await user.addUser(userform);
                    if (data.code === 200) {
                        this.$message({
                            message: data.msg,
                            type: 'success'
                        });
                        this.dialogaddFormVisible = false;
                        this.getData();
                    } else {
                        this.$message({
                            message: data.msg,
                            type: 'error'
                        });
                    }
                    // this.$refs.ruleFormuser.resetFields();  
                }
            });

        },
        udpatesubmitForm() {
            this.$refs.ruleFormupdateuser.validate(async (valid) => {
                if (valid) {
                    let userform1 = _.cloneDeep(this.userform1)
                    userform1.group_ids = userform1.group_ids.join(',')
                    const data = await user.updateUser(userform1);
                    if (data.code === 200) {
                        this.$message({
                            message: data.msg,
                            type: 'success'
                        });
                        this.dialogupdateFormVisible = false;
                        this.getData();
                    } else {
                        this.$message({
                            message: data.msg,
                            type: 'error'
                        });
                    }
                    
                    // this.$refs.ruleFormuser.resetFields(); 
                }
            });
        },
        AddUser() {
            this.getupgroup()
            this.dialogaddFormVisible = true;
            this.userform.username = '';
            this.userform.role = '1';
            this.userform.password = '';
            this.userform.repassword = '';
            this.userform.email = '';
            this.userform.department = '';
            this.userform.remark = '';
        },
        // 获取下拉组列表
        async getupgroup() {
            const data = await user.getupgroup();
            if (data.code === 200) {
                this.groupOptions = data.data.list
                this.groupOptions.push({
                    label: '无',
                    value: '无'
                })
            } else {
                this.$message({
                    message: data.msg,
                    type: 'error'
                });
            }
        },
        updateUser(rows) {
            this.getupgroup()
            console.log('rows', rows)
            this.userform1.id = rows.id;
            this.userform1.username = rows.username;
            this.userform1.email = rows.email;
            this.userform1.department = rows.department;
            console.log(this.userform1.group_ids)
            this.userform1.remark = rows.remark;
            this.userform1.role = rows.role;
            this.userform1.group_ids = rows.group_arr;
            // this.userform1.group_ids = [['14', '33']];
            if (this.userform1.group_ids.length === 0) {
                this.userform1.group_ids = [['无']]
            }

            this.dialogupdateFormVisible = true;
        },
        pwdcancelform() {
            this.pwddata.id = '';
            this.pwddata.username = '';
            this.pwddata.password = '';
            this.pwddata.repassword = '';
            this.pwddata.admin_password = '';
            this.dialogVisiblepwd = false;
            this.$refs.ruleFormPWD.resetFields();
        },
        change1 (val) {
            console.log(22, val)
        },
        pwdsubmitForm() { //保存修改密码
            this.$refs.ruleFormPWD.validate(async (valid) => {
                if (valid) {
                    const data = await user.updatePwd(this.pwddata);
                    if (data.code === 200) {
                        this.$message({
                            message: data.msg,
                            type: 'success'
                        });
                        this.dialogVisiblepwd = false;
                        this.getData();
                        this.pwddata.id = '';
                        this.pwddata.username = '';
                        this.pwddata.password = '';
                        this.pwddata.repassword = '';
                        this.pwddata.admin_password = '';
                        // this.$refs.ruleFormPWD.resetFields(); 
                    } else {
                        this.$message({
                            message: data.msg,
                            type: 'error'
                        });
                    }

                }
            });
        },
        updatepwd(id, username) {
            this.pwddata.username = username;
            this.pwddata.id = id;
            this.dialogVisiblepwd = true;
        },
        async setislogin(userid, status) { //设置禁用启用  
            let params = {
                ids: userid,
                status: status === 1 ? 5 : 1
            }
            const data = await user.userDel(params);
            if (data.code === 200) {
                this.$message({
                    message: data.msg,
                    type: 'success'
                });
                this.getData();
            } else {
                this.$message({
                    message: data.msg,
                    type: 'error'
                });
            }

        },
        // 账号有效期..........................................................................
        useDate(id, time) {
            this.timeForm.user_id = id
            this.timeForm.valuedate = time
            this.dialogUsedate = true;
        },
        async usedatesubmit() {
            let params = {
                id: this.timeForm.user_id,
                account_expire: this.timeForm.valuedate,
            }
            const dt = await user.saveYouxiaoqi(params);
            if (dt.code === 200) {
                this.$message.success(dt.msg)
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            }
        },
        usedatecancel() {
            this.dialogUsedate = false;
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

    }
})

</script>