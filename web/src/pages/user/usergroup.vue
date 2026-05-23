<template>
    <div>
        <div class="userbox context_box_bg">
            <div class="search-box">
                <div class="operationbutton">
                    <xz-button type="primary" @click="AddUsergroupgroup" size="small">新建</xz-button>
                    <del-button :width="170" @click="AllDel" style="margin-left: 8px;"
                        :disabled="!multipleSelection.length || (role === '2' && selectedID.includes(2)) || selectedID.includes(3) || selectedID.includes(4)">
                    </del-button>
                </div>
            </div>
            <el-table :data="tableData" style="width: 100%" v-loading="Loading" @selection-change="handleSelectionChange"  height="calc(100% - 102px)"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column type="selection" :selectable='checkboxT' width="55">
                </el-table-column>
                <el-table-column prop="name" label="用户组">
                </el-table-column>
                <el-table-column prop="number" label="组员数">
                </el-table-column>
                <el-table-column prop="pid_str" label="上级组">
                </el-table-column>

                <el-table-column prop="status_str" label="创建时间" width="400">
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                            <el-link class="link_primary" :underline="false" @click="updateUsergroup(scope.row)"
                                :disabled="role === '2' && Number(role) === scope.row.role && (scope.row.id !== Number(userID))">
                                编辑</el-link>
                            <el-link class="link_primary" :underline="false"
                                @click="updatepwd(scope.row.id, scope.row.username)"
                                :disabled="role === '2' && Number(role) === scope.row.role && (scope.row.id !== Number(userID))">
                                组员</el-link>
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
                            <!-- <el-link :underline="false" class="link_primary" style="padding:0">删除</el-link> -->
                        </div>
                        <div v-else>
                            <span class="">{{ scope.row.create_time }}</span>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination :page-size="10" background layout=" total,  prev, pager, next,sizes, jumper" :total="totalpage"
                :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>


        </div>
        <el-dialog title="新建用户组" :visible.sync="dialogaddFormVisible" :before-close="cancelform" width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="usergroupform" label-width="0" status-icon ref="ruleFormadduser" :rules="rules1">
                    <el-form-item label=" " prop="name">
                        <label class="dialog_item_label">组名</label>
                        <el-input v-model="usergroupform.name" size="small" style="width:320px" placeholder="请输入组名"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="pid">
                        <label class="dialog_item_label">上级组</label>
                        <el-cascader style="width:320px" v-model="usergroupform.pid" :options="groupOptions"
                            :props="props"></el-cascader>
                    </el-form-item>
                    <el-form-item label="" prop="range">
                        <label class="dialog_item_label" style="vertical-align:middle;">测试范围限制</label>
                        <el-switch v-model="usergroupform.range_open" active-color="#4C7AE3" inactive-color="#E4E7ED">
                        </el-switch>
                        <el-input type="textarea" :rows="10" class="txtareacontent" v-model="usergroupform.range"
                            size="small" style="width:320px;display:block;margin-left:114px; " autocomplete="off"
                            placeholder="请输入测试范围" v-if="usergroupform.range_open"></el-input>
                    </el-form-item>
                </el-form>
            </div>
        </el-dialog>
        <el-dialog title="编辑用户组" width="1184px" :visible.sync="dialogupdateFormVisible" :before-close="updatecancelform"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="udpatesubmitForm">确定</el-button>
                <el-button size="small" @click="updatecancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="usergroupform1" label-width="0" status-icon ref="ruleFormupdateuser" :rules="rules2">
                    <el-form-item label=" " prop="name">
                        <label class="dialog_item_label">组名</label>
                        <el-input :disabled="usergroupform1.role == 4 || usergroupform1.role == 3"
                            v-model="usergroupform1.name" size="small" style="width:320px" placeholder="请输入用户名"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="pid">
                        <label class="dialog_item_label">上级组</label>
                        <el-cascader @change = 'ceso' style="width:320px" v-model="usergroupform1.pid" :options="groupOptions"
                            :props="props"></el-cascader>
                    </el-form-item>
                    <el-form-item label="" prop="range">
                        <label class="dialog_item_label" style="vertical-align: middle;">测试范围限制</label>
                        <el-switch v-model="usergroupform1.range_open" active-color="#4C7AE3" inactive-color="#E4E7ED">
                        </el-switch>
                        <el-input type="textarea" :rows="10" class="txtareacontent" v-model="usergroupform1.range"
                            size="small" style="width:320px;display:block;margin-left:114px; " autocomplete="off"
                            placeholder="请输入测试范围" v-if="usergroupform1.range_open"></el-input>
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
    height: calc(100% - 39px);
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
export default ({
    name: 'usergroup',
    components: {
        XzButton,
        DelButton,
    },
    data() {
        return {
            props: { checkStrictly: true, emitPath: false, expandTrigger: 'hover' },
            value: false,
            valuedate: '2099-12-30',
            valuedate2: '2099-12-30',
            isMore: false,
            showOperateButton: false,
            rowId: '',
            multipleSelection: [],
            currentpage: 1,
            totalpage: 0,
            page_num: 1,
            tableData: [],
            Loading: false,
            usergroupform: {
                name: '',
                pid: 0,
                range_open: '',
                range: '',
            },
            usergroupform1: {
                name: '',
                pid: 0,
                range: '',
                range_open: ''
            },
            dialogaddFormVisible: false,
            dialogupdateFormVisible: false,
            dialogUsedate: false,
            rules1: {
                name: [
                    { required: true, message: '组名不能为空', trigger: 'blur' },
                    { max: 50, message: '组名最大长度不能超过50', trigger: 'blur' },
                ],
                pid: [
                    { required: true, message: '请选择上级', trigger: 'change' }
                ],
                // range: [
                //     { required: true, message: '测试范围不能为空', trigger: 'blur' },
                // ],
            },
            rules2: {
                name: [
                    { required: true, message: '组名不能为空', trigger: 'blur' },
                    { max: 50, message: '组名最大长度不能超过50', trigger: 'blur' },
                ],
                pid: [
                    { required: true, message: '请选择上级', trigger: 'change' }
                ],
                range: [
                    { required: true, message: '测试范围不能为空', trigger: 'blur' },
                ],
            },
            dialogVisiblepwd: false,
            adduser: true,
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
            groupOptions: []
        }
    },
    watch: {
        'usergroupform1.role': function (value) {
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
        this.$store.state.activefirstMenu = "/usergroup";
        this.role = this.commonjs.decryptCBC(localStorage.getItem('role'), this.commonjs.myKey);
        this.user = this.commonjs.decryptCBC(localStorage.getItem('user'), this.commonjs.myKey);
        this.userID = this.commonjs.decryptCBC(localStorage.getItem('user_id'), this.commonjs.myKey);
        this.pageSize = this.commonjs.pageSize;
    },
    mounted: function () {
        this.getData(); //获得列表数据
//         this.tableData =[ {
//     "id": 69,
//     "name": "1",
//     "number": 0,
//     "pid": 0,
//     "pid_str": "无",
//     "pid_arr": null,
//     "create_time": "2023-03-22 16:18:15",
//     "range_open": 0,
//     "range": ""
// }]
        this.totalpage = 1;
    },
    methods: {
        ceso(a){
            console.log(a,'aaaaaaaaaaaaaaaaa');
        },
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
            const dt = await user.usergroupList({
                page: this.page_num,
                size: this.pageSize,
            });
            if (dt.code === 200) {
                this.tableData = dt.data.list;
                this.totalpage = dt.data.total;
            } else {
                this.$message({
                    message: dt.msg || '错误，请重试',
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
            const data = await user.groupDel(params);
            if (data.code === 200) {
                this.$message({
                    message: data.msg || '成功！',
                    type: 'success'
                });
                this.alldelvisible = false;
                this.getData();
            } else {
                this.$message({
                    message: data.msg || '错误，请重试',
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
            const data = await user.groupDel(params);
            if (data.code === 200) {
                this.$message({
                    message: data.msg || '成功！' ,
                    type: 'success'
                });
                // scope._self.$refs[`popover_id-${scope.row.user_id}`].doClose();
                // scope._self.$refs[`popover-${scope.row.user_id}`].doClose()
                this.getData();
            } else {
                this.$message({
                    message: data.msg || '错误，请重试',
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
            this.dialogaddFormVisible = false;
            this.usergroupform.name = '';
            this.usergroupform.range = '';
            this.$refs.ruleFormadduser.resetFields();
        },
        updatecancelform() {
            this.usergroupform1.name = '';
            this.usergroupform1.range = '';
            this.dialogupdateFormVisible = false;
            this.$refs.ruleFormupdateuser.resetFields();
        },
        submitForm() {
            this.$refs.ruleFormadduser.validate(async (valid) => {
                if (valid) {
                    let params = {
                        id: this.usergroupform.id,
                        name: this.usergroupform.name,
                        pid: this.usergroupform.pid,
                        range: this.usergroupform.range,
                        range_open: this.usergroupform.range_open ? 1 : 0
                    }
                    const data = await user.addUsergroup(params);
                    if (data.code === 200) {
                        this.$message({
                            message: data.msg|| '成功！',
                            type: 'success'
                        });
                        this.dialogaddFormVisible = false;
                        this.getData();
                    } else {
                        this.$message({
                            message: data.msg||'错误，请重试',
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
                    let params = {
                        id: this.usergroupform1.id,
                        name: this.usergroupform1.name,
                        pid: this.usergroupform1.pid,
                        range: this.usergroupform1.range,
                        range_open: this.usergroupform1.range_open ? 1 : 0
                    }
                    const data = await user.updateUsergroup(params);
                    if (data.code === 200) {
                        this.$message({
                            message: data.msg || '成功！',
                            type: 'success'
                        });
                    } else {
                        this.$message({
                            message: data.msg || '错误，请重试',
                            type: 'error'
                        });
                    }
                    this.dialogupdateFormVisible = false;
                    this.getData();
                    // this.$refs.ruleFormuser.resetFields(); 
                }
            });
        },
        AddUsergroupgroup() {
            this.getupgroup()
            this.dialogaddFormVisible = true;
            this.usergroupform.name = '';
            this.usergroupform.pid = 0;
            this.usergroupform.range = '';
            this.usergroupform.range_open = false;
        },
                // 获取上级组列表
        async getupgroup() {
            const data = await user.getupgroup();
            if (data.code === 200) {
                this.groupOptions = data.data.list
                this.groupOptions.push({
                    label: '无',
                    value: '0'
                })
            } else {
                this.$message({
                    message: data.msg || '错误，请重试',
                    type: 'error'
                });
            }
        },
        async updateUsergroup(rows) {
            console.log(rows,'----------------')
            await this.getupgroup()
            this.usergroupform1.id = rows.id;
            this.usergroupform1.name = rows.name;
            this.usergroupform1.range = rows.range;
            this.usergroupform1.range_open = rows.range_open === 0 ? false : true;
            this.usergroupform1.pid = rows.pid === 0 ? '0' : String(rows.pid);
            
            // 如果 pid 为 0，说明是顶层组，不需要额外处理，直接显示为“无”（如果选项中有 0 对应“无”）
            // 如果 pid 不为 0，级联选择器会根据 pid 自动选中对应项（前提是 emitPath: false）
            
            // if (!this.usergroupform1.pid) {
            //     this.usergroupform1.pid = '无'
            // }
            console.log(rows,'22', this.usergroupform1.pid)
            // if(rows.range_open == 0){
            //     this.usergroupform1.range_open = false;
            // }else{
            //     this.usergroupform1.range_open = true;
            // }
            this.dialogupdateFormVisible = true;
        },
        // // 获取上级组列表
        // async getupgroup() {
        //     const data = await user.getupgroup();
        //     if (data.code === 200) {
        //         this.groupOptions = data.data.list
        //         console.log(this.groupOptions,'this.groupOptions');
        //         this.groupOptions.push({
        //             label: '无',
        //             value: 0
        //         })
        //     } else {
        //         this.$message({
        //             message: data.msg,
        //             type: 'error'
        //         });
        //     }
        // },
        //  updateUsergroup(rows) {
        //     console.log(rows,'==================')
        //     this.getupgroup()
        //     this.usergroupform1.id = rows.id;
        //     this.usergroupform1.name = rows.name;
        //     this.usergroupform1.range = rows.range;
        //     this.usergroupform1.range_open = rows.range_open === 0 ? false : true;
        //     if (rows.pid == 0) {
                
        //     }
        //     this.usergroupform1.pid =rows.pid_arr?rows.pid_arr.reverse():rows.pid_arr;
        //     // this.usergroupform1.pid =[rows.pid+'',rows.id+''] ;
        //     if (!this.usergroupform1.pid) {
        //         this.usergroupform1.pid = '无'
        //     }
        //     console.log(rows,'22', this.usergroupform1.pid)
        //     // if(rows.range_open == 0){
        //     //     this.usergroupform1.range_open = false;
        //     // }else{
        //     //     this.usergroupform1.range_open = true;
        //     // }
        //     this.dialogupdateFormVisible = true;
        // },

        updatepwd(id, username) {
            this.$router.push({
                name: 'groupmanagement', params: {
                    id: id
                }
            });
        },
        async setislogin(userid, is_active) { //设置禁用启用  
            const data = await user.setIsLoginUser({
                is_active: is_active ? 0 : 1,
                user_id: userid
            });
            if (data.success) {
                this.$message({
                    message: data.msg || '成功！',
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
        // 账号有效期..........................................................................
        useDate() {
            this.dialogUsedate = true;
        },
        usedatesubmit() {

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