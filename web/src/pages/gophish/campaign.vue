<template>
    <div>
        <div class="main-title  ">
             
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >钓鱼活动管理</label>
        </div>
        <div class="gophishbox context_box_bg">
            <div class="search-box">
                <div class="operationbutton">
                    <xz-button type="primary" @click="AddCampaign" size="small">新建钓鱼活动</xz-button>
                  
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
            <el-table :data="tableData" style="width: 100%"   height="calc(100% - 102px)"
                @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">  
                <el-table-column prop="name" label="钓鱼活动名称">
                </el-table-column>  
                <el-table-column prop="created_date" label="创建时间">
                </el-table-column> 
                <el-table-column prop="status" label="状态"  >
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('状态')"
                            :class="formData.activeStatus == '' ? '' : 'active'">状态 <i class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" 
                            v-model="formData.activeStatus"
                            size="small" ref="statusRef" @change="handlesearch">
                            <el-option v-for="(item, index) in statuslist" :key="index" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select> 
                    </template>
                    <template slot-scope="scope">
                        <div v-if="showOperateButton && rowId == scope.row.id">
                            <el-link class="link_primary" :underline="false" @click="handleinfo(scope.row)"
                                 >
                                详情</el-link>
                            <el-link class="link_primary" :underline="false" v-if="scope.row.status !='Completed'"
                                @click="handleEnd(scope.row)"
                                 >
                                结束</el-link>
                            <el-link class="link_primary" :underline="false"
                                @click="handleCopy(scope.row)"
                                 >
                                复制</el-link>
                        
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
                            <span class="tag_status tag_success1" v-if="scope.row.status == 'Emails Sent'"><i></i>已发送</span>
                            <span class="tag_status tag_primary" v-if="scope.row.status == 'In progress'"><i></i>进行中</span>
                            <span class="tag_status tag_success" v-if="scope.row.status == 'Completed'"><i></i>已完成</span>
                            <span class="tag_status tag_info" v-if="scope.row.status == 'Queued'"><i></i>排队中</span>

                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination :page-size="10" background layout=" total,  prev, pager, next,sizes, jumper" :total="totalpage"
                :current-page="currentpage" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>
 
        </div>
         <el-dialog :title="campaignid!=0?'编辑钓鱼活动':'创建钓鱼活动'" :visible.sync="dialogaddFormVisible" 
        :before-close="cancelform" width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">  
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
                <el-form :model="campaignform" ref="form" label-width="0" status-icon  :rules="rules" >
                    <el-row>
                        <el-col :span="12">
                             <el-form-item label="" prop="name">
                                <label class="dialog_item_label">事件名称<i class="is-required" style="float: right">*</i></label>
                                <el-input v-model="campaignform.name" size="small" style="width:380px" placeholder="请输入事件名称"
                                    maxlength="50"></el-input> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                              <el-form-item label="" prop="template">
                                <label class="dialog_item_label">邮件模板<i class="is-required" style="float: right">*</i></label> 
                                <el-select v-model="campaignform.template" size="small"  style="width:380px"  placeholder="请选择邮件模板">
                                    <el-option
                                    v-for="(item,i) in templatelist"
                                    :key="i"
                                    :label="item.name"
                                    :value="item.name">
                                    </el-option>
                                </el-select>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="12">
                            <el-form-item label="" prop="page">
                                <label class="dialog_item_label">钓鱼页面<i class="is-required" style="float: right">*</i></label>  
                                <el-select v-model="campaignform.page" size="small"  style="width:380px"  placeholder="请选择钓鱼页面">
                                    <el-option
                                    v-for="(item,i) in pagelist"
                                    :key="i"
                                    :label="item.name"
                                    :value="item.name">
                                    </el-option>
                                </el-select>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                             <el-form-item label="" prop="url">
                                <label class="dialog_item_label">钓鱼URL<i class="is-required" style="float: right">*</i></label>
                                <el-input v-model="campaignform.url" size="small" style="width:380px"  placeholder="请输入钓鱼URL"  ></el-input> 
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="12">
                             <el-form-item label="" prop="launch_date">
                                <label class="dialog_item_label">启动时间<i class="is-required" style="float: right">*</i></label> 
                                <el-date-picker
                                    v-model="campaignform.launch_date"
                                    type="datetime"
                                    size="small"
                                    style="width:380px"
                                    :default-value="defaultValue" 
                                    :picker-options="pickerOptions"
                                    :clearable=false
                                    placeholder="选择启动时间">
                                    </el-date-picker>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                            <el-form-item label=" " prop="send_by_date">
                                <label class="dialog_item_label">截止时间(可选)</label> 
                                 <el-date-picker
                                    v-model="campaignform.send_by_date"
                                    type="datetime"
                                     size="small"
                                     style="width:380px"
                                    placeholder="选择发送截止时间">
                                    </el-date-picker>
                            </el-form-item>
                        </el-col>
                    </el-row> 
                     <el-row>
                        <el-col :span="12">
                             <el-form-item label="" prop="smtp">
                                <label class="dialog_item_label">邮箱配置<i class="is-required" style="float: right">*</i></label> 
                                <el-select v-model="campaignform.smtp" size="small"  style="width:380px"  placeholder="请选择邮箱配置">
                                    <el-option
                                    v-for="(item,i) in smtplist"
                                    :key="i"
                                    :label="item.name"
                                    :value="item.name">
                                    </el-option>
                                </el-select>
                                <el-link type="primary" :underline="false" style="margin-left: 113px;"  @click="sendTest">发送测试邮件</el-link>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                            <el-form-item label="" prop="groups">
                                <label class="dialog_item_label">钓鱼目标<i class="is-required" style="float: right">*</i></label>  
                                <el-select v-model="campaignform.groups" size="small" multiple  style="width:380px"  placeholder="请选择钓鱼目标组">
                                    <el-option
                                    v-for="(item,i) in groupslist"
                                    :key="i"
                                    :label="item.name"
                                    :value="item.name">
                                    </el-option>
                                </el-select>

                            </el-form-item>
                        </el-col>
                    </el-row> 
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
:deep(.el-date-editor){
    input{
        width: 380px;
    }
}
.tag_status{
    width: 80px !important;
}
.cell {
    line-height: 15px; 
    >span {
        position: absolute;
    }
} 
i.is-required {
    margin-right: 4px;
    color: #f56c6c;
    font-size: 12px; 
}
:deep(thead)  {
    .cursorPointer {
        cursor: pointer; 
        position: absolute; 
        &.active {
            color: #4C7AE3; 
            i {
                color: #4C7AE3;
            }
        }
    } 
    .cell {
        line-height: 15px; 
        >span {
            position: absolute;
        }
    }

    .iconfont {
        color: rgba(72, 72, 102, 0.32);
        margin-left: 5px;
    }

    .el-select {
        height: 0;
        visibility: hidden;

        .el-input,
        .el-input__inner {
            height: 0 !important;
        }
    }
}
</style>
<script>
import XzButton from "../../components/XzButton.vue";
import DelButton from "../../components/DelButton.vue";
import user from "@/api/user.js";
import _ from 'lodash'
import gophish from "@/api/gophish.js";
import { isDateSame } from "xe-utils";
export default ({
    name: 'usermanagement',
    components: {
        XzButton,
        DelButton,
    },
    data() { 
        return {  
            isMore: false,
            showOperateButton: false,
            rowId: '',
            multipleSelection: [],
            currentpage: 1,
            totalpage: 0,
            page_num: 1,
            tableData: [], 
            Loading: false,
            formData:{
                search:'',
                activeStatus:'',
            },  
            pageSize: 10, 
            campaignform:{
                name:'',
                template:'',
                url:'',
                page:'',
                smtp:'',
                launch_date:new Date(),
                send_by_date:'',
                groups:[]
            },
            defaultValue: new Date(),  // 默认打开时定位到今天
            pickerOptions: {
                // 禁用日期：只允许选择今天及以后
                disabledDate: (time) => {
                    const today = new Date();
                    today.setHours(0, 0, 0, 0);
                    return time.getTime() < today.getTime();
                    },
                    // 控制时间选择：如果是“今天”，则时间只能选当前时间之后
                    shortcuts: [], // 可选：添加快捷选项
                    onPick: ({ maxDate, minDate }) => {}, // 可选
            },
            dialogaddFormVisible:false,
            campaignid:0,
            templatelist:[],
            pagelist:[],
            smtplist:[],
            groupslist:[],
            dialogVisible:false,
            sendtestform:{
                first_name:'',
                email:'',
                position:'',
            },
            rules:{
                name: [
                    { required: true, message: "时间名称不能为空", trigger: "blur" }, 
                ],
                template:[
                    { required: true, message: "请选择邮件模板", trigger: "change" }, 
                ],
                page:[
                    { required: true, message: "请选择钓鱼页面", trigger: "change" }, 
                ],
                url:[
                    { required: true, message: "钓鱼URL不能为空", trigger: "blur" }, 
                ],
                launch_date:[
                    { required: true, message: "请选择启动时间", trigger: "change" }, 
                ],
                smtp:[
                     { required: true, message: "请选择邮箱配置", trigger: "change" }, 
                ],
                groups:[
                    { required: true, message: "请选择钓鱼目标", trigger: "change" },
                ]
            },
            statuslist:[
                {
                    label:'全部',
                    value:'',
                },
                {
                    label:'已发送',
                    value:'Emails Sent',
                },
                {
                    label:'进行中',
                    value:'In progress',
                },
                {
                    label:'已完成',
                    value:'Completed',
                },
                {
                    label:'排队中',
                    value:'Queued',
                },
            ]
        }
    }, 
    watch: {
        // 当用户切换日期时，动态更新时间禁用规则
        selectedDateTime() {
        this.pickerOptions = this.getPickerOptions();
        },
    },
    created: function () {
        this.$store.state.activefirstMenu = "/campaign"; 
    },
    mounted: function () {
        this.getData(); //获得列表数据
        // 初始化 pickerOptions 的时间限制逻辑
        this.pickerOptions = this.getPickerOptions();
    },
    methods: {
        clickButton(type) {
            switch (type) {
              case '状态':
                this.$refs.statusRef.toggleMenu(); 
                break;
            }
          },
        getPickerOptions() {
            const self = this;
            return {
                disabledDate(time) {
                // 禁用今天之前的日期
                const today = new Date();
                today.setHours(0, 0, 0, 0);
                return time.getTime() < today.getTime();
                },
                disabledHours() {
                const selectedDate = self.selectedDateTime;
                const now = new Date();
                // 如果选择的是“今天”，则禁用已经过去的小时
                if (self.isSameDay(selectedDate, now)) {
                    const hours = [];
                    const currentHour = now.getHours();
                    for (let i = 0; i < currentHour; i++) {
                    hours.push(i);
                    }
                    return hours;
                }
                // 不是今天，则不禁用任何小时
                return [];
                },
                disabledMinutes(hour) {
                    const selectedDate = self.selectedDateTime;
                    const now = new Date();
                    // 如果是“今天”且小时是当前小时，则禁用已过去的分钟
                    if (
                        self.isSameDay(selectedDate, now) &&
                        hour === now.getHours()
                    ) {
                        const minutes = [];
                        const currentMinute = now.getMinutes();
                        for (let i = 0; i < currentMinute; i++) {
                        minutes.push(i);
                        }
                        return minutes;
                    }
                    return [];
                },
                disabledSeconds(hour, minute) {
                    const selectedDate = self.selectedDateTime;
                    const now = new Date();
                    // 如果是“今天”且是当前小时+分钟，则禁用已过去的秒
                    if (
                        self.isSameDay(selectedDate, now) &&
                        hour === now.getHours() &&
                        minute === now.getMinutes()
                    ) {
                        const seconds = [];
                        const currentSecond = now.getSeconds();
                        for (let i = 0; i < currentSecond; i++) {
                        seconds.push(i);
                        }
                        return seconds;
                    }
                    return [];
                    },
            };
        },
        isSameDay(date1, date2) {
            return (
                date1 &&
                date2 &&
                date1.getFullYear() === date2.getFullYear() &&
                date1.getMonth() === date2.getMonth() &&
                date1.getDate() === date2.getDate()
            );
        }, 
        handlesearch(){
            this.getData();
        },
        handleReset(){
            this.formData.search = '';
            this.formData.activeStatus='';
            this.getData();
        }, 
        async getData() {
            this.Loading = false;
            const dt = await gophish.campaignall({
                page:this.page_num,
                size:this.pageSize,
                search:this.formData.search,
                activeStatus:this.formData.activeStatus,
            });
            console.log(this.formData.activeStatus)
            if (dt.code === 200) {
                this.tableData = dt.data.campaigns;
                this.totalpage = dt.data.total;
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        }, 
        handleinfo(row){ //详情
            this.$router.push({
                path: `/campaignDetail`,
                query: { 
                    id: row.id, 
                    name:row.name,
                }
            });
        },
        // 结束
        async handleEnd(row){ 
            var data = await gophish.campaigncomplete({
                id:row.id
            })
            if (data.code === 200) {
                 this.$message({
                    message: data.msg || '操作成功！' ,
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
        //复制
        async handleCopy(row){
            this.getSelecteData()
            this.dialogaddFormVisible = true; 
  
            var dt = await gophish.campaigndetail({
                id:row.id
            })
            if(dt.code == 200){ 
                this.campaignform.name = 'Copy of '+dt.data.name;
                this.campaignform.launch_date = dt.data.launch_date;
                this.campaignform.send_by_date = dt.data.send_by_date;
                this.campaignform.smtp = dt.data.smtp.name
                this.campaignform.template =  dt.data.template.name;
                this.campaignform.url = dt.data.url;
                this.campaignform.page = dt.data.page.name
            }
        },
        fncancel(scope) {
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
            scope._self.$refs[`popover-${scope.row.id}`].doClose()
        },
        async fnDel(scope) { //单个删除 
            const data = await gophish.campaigndelete({
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
                    message: data.msg,
                    type: 'error'
                });
            }

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
            this.campaignform.name='';
            this.campaignform.template='';
            this.campaignform.url='';
            this.campaignform.page='';
            this.campaignform.smtp='';
            this.campaignform.launch_date='';
            this.campaignform.send_by_date='';
            this.campaignform.groups=[];
        }, 
        async submitForm() { 
             this.$refs.form.validate(async (valid) => {
                if (valid) { 
                    const data = await gophish.campaigncreate({
                        name:this.campaignform.name,
                        template:{
                            name:this.campaignform.template,
                        },
                        url:this.campaignform.url,
                        page:{
                            name:this.campaignform.page,
                        },
                        smtp:{
                            name:this.campaignform.smtp,
                        },
                        launch_date:this.campaignform.launch_date,
                        send_by_date:this.campaignform.send_by_date,
                        groups:this.campaignform.groups.map(item => ({ name: item })) 
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
            }); 
        },
        
        AddCampaign() {
            this.getSelecteData()
            this.dialogaddFormVisible = true; 
        },
        // 获取下拉组列表
        getSelecteData() {
            this.gettemplatelist(); // 邮件模板
            this.getpagelist();//钓鱼页面
            this.getsmtplist();// 邮箱配置
            this.getgroupslist(); // 钓鱼目标组
        },
        // 邮件模板
        async gettemplatelist(){
            const dt = await gophish.templateall({
                page:1,
                size:10000,
                search:'',
            });
            if (dt.code === 200) {
                this.templatelist = dt.data.templates; 
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        // 钓鱼页面
        async getpagelist(){
            const dt = await gophish.pageall({
                page:1,
                size:10000,
                search:'',
            });
            if (dt.code === 200) {
                this.pagelist = dt.data.landingPage; 
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        //邮箱配置
        async getsmtplist(){
            const dt = await gophish.profileall({
                page:1,
                size:10000,
                search:'',
            });
            if (dt.code === 200) {
                this.smtplist = dt.data.sendingProfile; 
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        // 钓鱼目标组
        async getgroupslist(){
            const dt = await gophish.groupall({
                page:1,
                size:10000,
                search:'',
            });
            if (dt.code === 200) {
                this.groupslist = dt.data.groups; 
            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            } 
        },
        updateUser(rows) { 
            this.dialogupdateFormVisible = true;
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
         //发送测试邮件
        async sendTest(){
            if(this.campaignform.smtp&&this.campaignform.page&&this.campaignform.template){ 
                this.dialogVisible = true; 
            }else{
                 this.$message({
                    message: '请补充完善表单信息',
                    type: 'warning'
                });
            }
        },
        handleClose(){
            this.dialogVisible = false; 
            this.sendtestform.email = '';
            this.sendtestform.first_name = '';
            this.sendtestform.position = '';
        },
        //提交 发送测试邮件
        async handleSendTest(){
            if(this.campaignform.smtp&&this.campaignform.page&&this.campaignform.template){ 
                var dt = await gophish.send_test_email({
                    first_name:this.sendtestform.first_name,
                    email:this.sendtestform.email,
                    position:this.sendtestform.position,
                    url:'',
                    template:{},
                    smtp:{
                        name:this.campaignform.smtp,
                    },
                    page:{
                        name:this.campaignform.page,
                    },
                    template:{
                        name:this.campaignform.template,
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
                        message: dt.msg,
                        type: 'error'
                    });
                }
            }else{
                 this.$message({
                    message: '请补充完善表单信息',
                    type: 'warning'
                });
            }  
        },

    }
})

</script>