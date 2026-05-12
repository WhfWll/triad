//系统设置---系统工具---最后一个Tab  其他设置 页面代码
<template>
    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;">
        <el-row :gutter="20">
            <el-col :span="12">
                <div class="title_left_line">
                    <label>系统告警</label>
                    <span>
                        <el-switch v-model="otherconfig.system_threshold" class="elSwitch" @change="changeWarnstatus">
                        </el-switch>
                    </span>
                </div>
                <div class="div_block">
                    <el-form ref="warningform" :model="warningform" label-width="140px" class="sysform" :rules="setUp">
                        <el-form-item label="CPU告警阈值：" class="syswarnvalue" prop="cpu_threshold">
                            <el-input v-model="warningform.cpu_threshold" :disabled="!otherconfig.system_threshold"
                                style="width:calc(100% - 190px)">
                            </el-input> %
                            <el-tooltip placement="right">
                                <div slot="content">告警范围1%~100%</div>
                                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <el-form-item label="内存告警阈值：" class="syswarnvalue" prop="memory_threshold">
                            <el-input v-model="warningform.memory_threshold" :disabled="!otherconfig.system_threshold"
                                style="width:calc(100% - 190px)">
                            </el-input> %
                            <el-tooltip placement="right">
                                <div slot="content">告警范围1%~100%</div>
                                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <el-form-item label="硬盘告警阈值：" class="syswarnvalue" prop="disk_threshold">
                            <el-input v-model="warningform.disk_threshold" :disabled="!otherconfig.system_threshold"
                                style="width:calc(100% - 190px)">
                            </el-input> %
                            <el-tooltip placement="right">
                                <div slot="content">告警范围1%~100%</div>
                                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <el-form-item label="流量告警阈值：" class="syswarnvalue" prop="flow_threshold">
                            <el-input v-model="warningform.flow_threshold" :disabled="!otherconfig.system_threshold"
                                style="width:calc(100% - 190px)">
                            </el-input> Mbps
                            <el-tooltip placement="right">
                                <div slot="content">告警范围1~1000M</div>
                                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                            </el-tooltip>
                        </el-form-item>
                    </el-form>
                    <el-button type="primary" class="div_blockbtn" :disabled="!otherconfig.system_threshold"
                        @click="btnSaveWarning">保存设置</el-button>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="title_left_line">
                    <label>Syslog服务</label>
                    <span>
                        <el-switch v-model="otherconfig.syslogopen" class="elSwitch" @change="changeSysstatus">
                        </el-switch>
                    </span>
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            syslog服务用于配置syslog服务器地址和接收日志的类型，勾选不同的日志类型时，将向syslog服务器发送指定类型的日志。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block">
                    <el-form ref="servelform" :model="servelform" :rules="rules1" label-width="150px" class="sysform">
                        <el-form-item label="Syslog服务器地址：" class="syswarnvalue " prop="address">
                            <el-input v-model="servelform.address" class="innerinput"
                                :disabled="!otherconfig.syslogopen">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="Syslog服务器端口：" class="syswarnvalue" prop="port">
                            <el-input v-model="servelform.port" class="innerinput"
                                :disabled="!otherconfig.syslogopen">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="发送日志类型：" class="syswarnvalue" prop="log_type">
                            <el-checkbox-group v-model="servelform.log_type" :disabled="!otherconfig.syslogopen"
                                style="display: inline-block;width:290px;text-align: left;">
                                <el-checkbox :label="item.value"  v-for="item in sendLogTypeOptions" :key="item.value">{{item.label}}
                                </el-checkbox>
                            </el-checkbox-group>
                        </el-form-item>
                    </el-form>
                    <el-button type="primary" style="margin-top:54px" class="div_blockbtn"
                        :disabled="!otherconfig.syslogopen" @click="sysform3()">保存设置</el-button>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="title_left_line">
                    <label>系统访问白名单</label>
                    <span>
                        <el-switch v-model="otherconfig.whitelist" class="elSwitch" @change="changeWhitestatus()">
                        </el-switch>
                    </span>
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            系统访问白名单功能开启后，只有白名单内的IP才能访问系统，因此设置白名单IP时要考虑清晰，否则将导致不能访问系统。<br />
                            白名单IP支持IP、IP段格式，且须用换行隔开。<br />
                            示例：“192.168.0.127”、“192.168.0.10-127”、”192.168.0/24”<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block">
                    <div>
                        <label class="dialog_item_label" style="border-left: none;">白名单IP：</label>
                        <el-input v-model="whitelistinfo" type="textarea" size="small" :rows="10" resize="none"
                            style=" width:70%;vertical-align: top;" placeholder="" maxlength="1000"
                            :disabled="!otherconfig.whitelist"></el-input>
                    </div>
                    <el-button type="primary" style="margin-top:10px" class="div_blockbtn"
                        :disabled="!otherconfig.whitelist" @click="whiteform3()">保存设置</el-button>
                </div>

            </el-col>
            <el-col :span="12">
                <div class="title_left_line">
                    <label>邮箱配置</label>
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            邮箱配置是给系统设定邮箱账号，让系统能够通过邮件系统将测试报告发送给用户邮箱，邮箱使用SMTP协议发送邮件。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip> 
                </div>
                <div class="div_block">
                    <el-form ref="emailform" :model="emailform" :rules="rules2" label-width="140px"
                        style="text-align: left;" class="sysform">
                        <el-form-item label="邮箱服务器地址：" class="syswarnvalue" prop="address">
                            <el-input v-model="emailform.address" class="innerinput">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="邮箱服务器端口：" class="syswarnvalue" prop="port">
                            <el-input v-model="emailform.port" class="innerinput">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="邮箱账号：" prop="username" style="text-align: left;">
                            <el-input v-model="emailform.username" class="innerinput">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="邮箱密码：" class="syswarnvalue" prop="password">
                            <el-input v-model="emailform.password" type='password'  class="innerinput">
                            </el-input>
                        </el-form-item>
                    </el-form>
                    <el-button type="primary" class="div_blockbtn" @click="emailform4()">保存设置</el-button>
                </div>
            </el-col>
            <!-- <el-col :span="12">
                <div class="title_left_line">
                    <label>网络配置</label>
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            网络配置将修改系统的网络配置参数，配置不当可能影响系统的正常访问。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block" style="height:406px">
                    <el-form ref="networkform" :model="networkform" :rules="rules3" label-width="140px"
                        style="text-align: left;" class="sysform">
                        <el-form-item label="IP地址：" class="syswarnvalue" prop="address">
                            <el-input v-model="networkform.address" class="innerinput">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="子网掩码：" class="syswarnvalue" prop="netmask">
                            <el-input v-model="networkform.netmask" class="innerinput">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="默认网关：" prop="gateway" style="text-align: left;">
                            <el-input v-model="networkform.gateway" class="innerinput">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="DNS服务器：" class="syswarnvalue" prop="dns">
                            <el-input v-model="networkform.dns" class="innerinput">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="备用DNS服务器：" class="syswarnvalue" prop="spare_dns" label-width="130px" style="margin-left:10px;">
                            <el-input v-model="networkform.spare_dns" class="innerinput">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="业务端口：" class="syswarnvalue" prop="port">
                            <el-input v-model="networkform.port" class="innerinput">
                            </el-input>
                        </el-form-item>
                    </el-form>
                    <el-button type="primary" class="div_blockbtn" @click="networkform5()">保存设置</el-button>
                </div>
            </el-col> -->
            <!-- <el-col :span="12">
                <div class="title_left_line">
                    <label>路由配置</label>
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            路由配置将更新系统的路由表，配置不当可能会影响系统的正常访问。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block" style="height:406px;">
                     <xz-button type="primary" @click="Addluyou()" :disabled="false" size="small" style="margin-bottom:12px">新增路由</xz-button>
                        <el-table
                            class="diatable"
                            :data="Luyoudata"
                            size='small'
                            style="width: 100%;" :height="330">
                            <el-table-column
                                prop="ip"
                                label="目的网络" >
                                <template slot-scope="scope">
                                    <span v-if="!scope.row.canInput">{{scope.row.ip}}</span>
                                    <el-input v-else v-model="scope.row.ip"/>
                                </template>
                            </el-table-column>
                            <el-table-column                                     
                                label="子网掩码"  prop="netmask">
                                <template slot-scope="scope">
                                    <span v-if="!scope.row.canInput">{{scope.row.netmask}}</span>
                                    <el-input v-else v-model="scope.row.netmask"/>
                                </template>
                            </el-table-column>
                            <el-table-column
                                prop="gateway"
                                label="网关" >
                                <template slot-scope="scope">
                                    <span v-if="!scope.row.canInput">{{scope.row.gateway}}</span>
                                    <el-input v-else v-model="scope.row.gateway"/>
                                </template>
                            </el-table-column>
                            <el-table-column
                                    width="200"
                                label="操作">
                                    <template slot-scope="scope">
                                    <el-link  :underline="false" class="link_primary"  @click="Saveluyou(scope.row)" v-if="scope.row.canInput">保存</el-link>
                                    <el-popover
                                        placement="bottom"
                                        width="200"   
                                        :ref="`popover_id-${scope.$index}`"
                                        popper-class="delButton_popper" >
                                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                        <div style="text-align: right; margin: 0">
                                            <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.$index}`].doClose()">取消</el-button>
                                            <el-button size="mini" type="primary" @click="handleDelRule(scope)">确定</el-button>
                                        </div> 
                                        <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference" >删除</el-link>  
                                    </el-popover>
                                </template>
                            </el-table-column>
                        </el-table>
                </div>
            </el-col> -->
            
        </el-row>
    </div>
</template>

<script> 
import DelButton from "@/components/DelButton.vue";
import XzButton from "@/components/XzButton.vue";
import { otherset ,system} from '@/api/system.js'
import common from '@/utils/common.js'
import $ from 'jquery'
export default({
    name:'othersysset',
    components: {
        XzButton,
        DelButton
  	},
    data(){ 
        //  系统告警   相关
        let threshold = (rule, value, callback)=>{
            if (value < 0 || value > 100) return callback(new Error(`只能在0-100之间`));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()

        }
        let memory = (rule, value, callback)=>{
            if (value < 0 || value > 100) return callback(new Error(`只能在0-100之间`));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()
        }
        let disk = (rule, value, callback)=>{
            if (value < 0 || value > 100) return callback(new Error(`只能在0-100之间`));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()

        }
        let flow = (rule, value, callback)=>{
            if (value < 0 || value > 10000) return callback(new Error(`只能在0-10000之间`));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()

        }
        //其他设置 ip验证。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
        var validatePass2 = (rule, value, callback) => { 
            if (!value) {
              callback(new Error('请输入服务器地址'));
            } else {
                const re =
                    /^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])\.((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$/
                if (re.test(value)) {
                    callback();
                } else {
                    return callback(new Error('请输入正确的服务器地址'));
                }
            }
        };
         //其他设置 端口验证。
        var validatePass3 = (rule, value, callback) => { 
            if (!value) {
              callback(new Error('请输入服务器端口'));
            } else {
                const re = /^[0-9]*$/;
                if (re.test(value)) {
                    if ((value > 0 && value < 65535) || value === 0) {
                        callback();
                    } else {
                        return callback(new Error('请输入正确的服务器端口'));
                    }
                } else {
                    return callback(new Error('请输入正确的服务器端口'));
                }
            }
        };
    	return{
            AvailabilityScore:false,
            Luyoudata:[],
            //其他设置 。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
            warningform:{
                cpu_threshold: 10,
                disk_threshold: 0,
                flow_threshold:0,
                memory_threshold: 0,
            },
            otherconfig:{
                system_threshold:false,
                syslogopen:false,
                whitelist:false,
            },
            whitelistinfo:'',
            // sendLogTypeOptions: ['审计日志', '告警日志', '报错日志', '调试日志'],
            sendLogTypeOptions: [
                {
                    label: '审计日志',
                    value: '1'
                },
                {
                    label: '调试日志',
                    value: '2'
                },
                {
                    label: '告警日志',
                    value: '3'
                },
                {
                    label: '报错日志',
                    value: '4'
                }
            ],
            servelform:{
                log_type:[],
                address:'',
                port:'',  
            },
            emailform:{
                address:'',
                port:'',
                username:'',
                password:'',
            },
            networkform:{
                address:'',
                port:'',
                gateway:'',
                dns:'',
                spare_dns:'',
                netmask:''
            },
            securityform:{
                password_rank:'',
                password_cycle:'',
                login_timeout:'',
                user_limit:'',
                admin_limit:'',
                ban_time:'',
            },
            setUp:{
                // CPU告警阈值
                cpu_threshold:[
                    { required: true, message: 'CPU告警阈值不能为空', trigger: ['blur','change'] }, 
                    { validator:threshold, trigger: ['blur','change'] }
                ],
                // 内存告警阈值
                memory_threshold:[
                    { required: true, message: '内存告警阈值不能为空', trigger: ['blur','change'] }, 
                    { validator:memory, trigger: ['blur','change'] }
                ],
                //  硬盘报警
                disk_threshold:[
                    { required: true, message: '硬盘报警阈值不能为空', trigger: ['blur','change'] }, 
                    { validator:disk, trigger: ['blur','change'] }
                ],
                //  流量告警
                flow_threshold:[
                    { required: true, message: '流量告警阈值不能为空', trigger: ['blur','change'] }, 
                    { validator:flow, trigger: ['blur','change'] }
                ],
                
            },
            rules1: {
                address: [
                    { required: true, message: 'Syslog服务器地址不能为空', trigger: 'blur' }, 
                    { validator: validatePass2, trigger: 'blur' },
                ],
                port:[
                     { required: true, message: 'Syslog服务器端口不能为空', trigger: 'blur' },
                     { validator: validatePass3, trigger: 'blur' },
                ], 
                log_type:[
                     { required: true, message: '', trigger: 'blur' }
                ], 
            },
            rules2: {
                address: [
                    { required: true, message: '邮箱服务器地址不能为空', trigger: 'blur' }, 
                    // { validator: validatePass2, trigger: 'blur' },
                    // { max: 100, message: '邮箱服务器地址不正确', trigger: 'blur' },
                ],
                port:[
                     { required: true, message: '邮箱服务器端口不能为空', trigger: 'blur' },
                     { validator: validatePass3, trigger: 'blur' },
                ], 
                username: [
                    { required: true, message: '邮箱不能为空', trigger: 'blur' },
                    { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur'] }
                ],
                password: [
                    { required: true, message: '密码不能为空', trigger: 'blur' },
                    { min:6, max: 20, message: '请输入正确的密码', trigger: ['blur'] }
                ],
            },
            rules3: {
                address: [
                    { required: true, message: 'IP地址不能为空', trigger: 'blur' }, 
                    { validator: validatePass2, trigger: 'blur' },
                    // { max: 100, message: '邮箱服务器地址不正确', trigger: 'blur' },
                ],
                netmask:[
                     { required: true, message: '子网掩码不能为空', trigger: 'blur' },
                     { validator: validatePass2, trigger: 'blur' },
                ], 
                gateway: [
                    { required: true, message: '默认网关不能为空', trigger: 'blur' }, 
                    { validator: validatePass2, trigger: 'blur' },
                ],
                dns:[
                     { required: true, message: 'DNS服务器不能为空', trigger: 'blur' },
                     { validator: validatePass2, trigger: 'blur' },
                ], 
                // spare_dns: [
                //     { required: false, message: '备用DNS服务器不能为空', trigger: 'blur' }, 
                //     { validator: validatePass2, trigger: 'blur' },
                // ],
                port:[
                     { required: true, message: '端口不能为空', trigger: 'blur' },
                     { validator: validatePass3, trigger: 'blur' },
                ], 
                
            },
        //      // 。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
            
    	}
    }, 
    created:function(){
    },
    mounted:function(){
        this.get_email_config();  
        this.get_syslog_config();
        this.get_whitelist_config(); 
        this.getwarninginfo();
        this.get_network_config(); 
         this.getLuyouList(); 
     this.getscoredata();
    },
    beforeDestroy(){
    },
    methods:{

        //三个开关按钮
        changeWarnstatus(value){
            if (!value) {
                this.btnSaveWarning();
            }
        },
        changeSysstatus (value){
            if (!value) {
                this.sysform3();
            }
        },
        changeSysstatus2 (value){
                this.getscoredataSave();
        },
        changeWhitestatus (value){
            if(this.otherconfig.whitelist == false){
                if (!value) {
                this.whiteform3(); 
               }
            }
            
        },
        // 路由配置...............................................................................................
      async getLuyouList(){
            const res = await otherset.getLuyou();
            if(res.code === 200){  
                this.Luyoudata = res.data.list
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        validateIp (value, type) {
            if (!value) {
                let text = ''
                switch (type) {
                    case 'ip':
                        text = '请输入目的网络地址'
                        break
                    case 'netmask':
                        text = '请输入子网掩码地址'
                        break
                    case 'gateway':
                        text = '请输入网关地址'
                        break
                }
              return text
            } else {
                const re =
                    /^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])\.((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$/
                if (re.test(value)) {
                    return ''
                } else {
                    let text = ''
                    switch (type) {
                        case 'ip':
                            text = '请输入正确的目的网络地址'
                            break
                        case 'netmask':
                            text = '请输入正确的子网掩码地址'
                            break
                        case 'gateway':
                            text = '请输入正确的网关地址'
                            break
                    }
                    return text
                }
            }
        },
        //新增路由
        async Saveluyou(row){
            let valid1 = this.validateIp(row.ip, 'ip')
            if (valid1 !== '') {
                this.$message.warning(valid1)
                return false
            }
            let valid2 = this.validateIp(row.netmask, 'netmask')
            if (valid2 !== '') {
                this.$message.warning(valid2)
                return false
            }
            let valid3 = this.validateIp(row.gateway, 'gateway')
            if (valid3 !== '') {
                this.$message.warning(valid3)
                return false
            }
            let params = {
              ip:row.ip,
              netmask: row.netmask,
              gateway: row.gateway
            }
            const res = await otherset.addnewLuyou(params);
            if(res.code === 200){ 
                this.$message({
                    message:'保存成功',
                    type: 'success'
                });
                row.canInput = false
                this.getLuyouList(); 
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        // 删除路由
         async handleDelRule(scope){ //删除
         let row = scope.row
         let index = scope.$index
         if (row.canInput) {
            // 尚未保存的
            this.Luyoudata.splice(index, 1)
            scope._self.$refs[`popover_id-${index}`].doClose()
         } else {
            let params = {
              ip:row.ip,
              netmask: row.netmask,
              gateway: row.gateway
            }
            const res = await otherset.deleteLuyou(params);
            if(res.code === 200){ 
                    this.$message({
                        message:'删除成功',
                        type: 'success'
                    });
                scope._self.$refs[`popover_id-${index}`].doClose()
                this.getLuyouList(); 
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            } 
         }
            
       },
        //新增路由
        Addluyou(){
            this.Luyoudata.unshift({
                ip: '',
                gateway: '',
                netmask: '',
                canInput: true
            })
        },
        //系统告警接口。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
        async getwarninginfo(){ //获取 告警设置 信息
            
            const res = await otherset.getSysteminfo();
            if(res.code === 200){  
                this.otherconfig.system_threshold = res.data.isOpen === 1 ? true : false;
                this.warningform.cpu_threshold = res.data.cpuWarn; 
                this.warningform.disk_threshold =res.data.diskWarn; 
                this.warningform.flow_threshold =res.data.flowWarn; 
                this.warningform.memory_threshold =res.data.memoryWarn; 
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        async getscoredata(){ //可利用评分配置信息
            
            const res = await system.getusescoreinfo();
            if(res.code === 200){  
               console.log(res.data.isOpen,'可利用评分配置信息');
               this.AvailabilityScore = res.data.isOpen === 1 ? true : false;
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        async getscoredataSave(){ //可利用评分配置保存
            const res = await system.getusescoresaveSave({isOpen:this.AvailabilityScore ? 1 : 2});
            if(res.code === 200){  
                this.getscoredata()
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        btnSaveWarning(){ //保存告警设置 
            this.$refs.warningform.validate( async valid=>{
                if(!valid) return ;
           
                const res = await otherset.saveSystemWarn({
                            cpuWarn: this.warningform.cpu_threshold,
                            diskWarn: this.warningform.disk_threshold,
                            flowWarn:this.warningform.flow_threshold,
                            // flow_threshold:parseInt(this.warningform.flow_threshold)*1024,//8-30根据指示去掉二次处理wm
                            memoryWarn: this.warningform.memory_threshold,
                            isOpen: this.otherconfig.system_threshold ? 1 : 2      //打开传1  关闭传2
                        });
                if(res.code === 200){
                    if (this.otherconfig.system_threshold) {
                        this.$message({
                        message:'保存成功',
                        type: 'success'
                       });
                    }
                this.getwarninginfo();
                }else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }
            })
        },
        //系统访问白名单 获取接口。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
        async get_whitelist_config(){
                const res = await otherset.getWhite();
                if(res.code === 200){  
                    this.otherconfig.whitelist = res.data.isOpen === 1 ? true : false;
                    this.whitelistinfo =  res.data.ip;
                }else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }  
        },
        //系统访问白名单 配置保存接口
        async whiteform3(){
            let params = {
                // whitelist: JSON.stringify({
                    isOpen:this.otherconfig.whitelist ? 1: 2,
                    ip: this.whitelistinfo,
                // })
            }
            const res = await otherset.saveWhite(params);
            if(res.code === 200){ 
                // if (this.otherconfig.whitelist) {
                    this.$message({
                        message:'保存成功',
                        type: 'success'
                    });
                // }
                this.get_whitelist_config(); 
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }  
        },
         //syslog获取接口。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
         async get_syslog_config(){
                const res = await otherset.getSyslog();
                if(res.code === 200){  
                    this.otherconfig.syslogopen = res.data.isOpen === 1 ? true : false;
                    this.servelform.log_type = res.data.types.split(',')
                    this.servelform.port = res.data.port;
                    this.servelform.address = res.data.ip;
                }else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }  
        },
        //syslog 邮箱配置保存接口
        async sysform3(){

            if (this.servelform.log_type.length === 0) {
                this.$message.warning('请至少选择一项发送日志类型')
                return false
            }
            let params = {
                isOpen:this.otherconfig.syslogopen ? 1 : 2,
                ip: this.servelform.address,  
                port:this.servelform.port,
                types: this.servelform.log_type.join(',')
            }
            const res = await otherset.saveSyslog(params);
                if(res.code === 200){
                    if (this.otherconfig.syslogopen) {
                        this.$message({
                            message:'保存成功',
                            type: 'success'
                        });
                    }
                    
                    this.get_syslog_config();
                }else{
                    this.$message({
                        message:'保存失败',
                        type: 'error'
                    });
                }
        },
         //邮箱配置获取接口。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
         async get_email_config(){
                const res = await otherset.getEmail();
                if(res.code === 200){  
                    this.emailform.port = res.data.port;
                    this.emailform.address = res.data.address;
                    this.emailform.username = res.data.username;
                    // this.emailform.password = res.data.password;                 
                    this.emailform.password = this.commonjs.decryptCBC(res.data.password,this.commonjs.myKey);
                }else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }  
        },
        //其他设置 邮箱配置保存接口
        emailform4(){
            this.$refs.emailform.validate( async valid=>{
                if(!valid) return ;
           
                const res = await otherset.saveEmail({
                            address: this.emailform.address,  
                            port:this.emailform.port,
                            username:this.emailform.username,
                            // password:this.emailform.password,
                            password:this.$commonjs.encryptCBC(this.emailform.password, this.$commonjs.myKey)
                        });
                if(res.code === 200){
                    this.$message({
                        message:'保存成功',
                        type: 'success'
                    });
                    this.get_email_config(); 
                }else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }
            })
        },
    //。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
         //网络配置获取接口
         async get_network_config(){
                const res = await otherset.getNetwork();
                if(res.code === 200){   
                    this.networkform.address = res.data.ip;
                    this.networkform.netmask = res.data.mask;
                    this.networkform.gateway = res.data.gateway;
                    this.networkform.dns = res.data.dnsServer;
                    this.networkform.spare_dns = res.data.standbyDnsServer;
                    this.networkform.port = res.data.webPort;
                }else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }  
        },
        //网络配置保存接口
        networkform5(){
            this.$refs.networkform.validate( async valid=>{
                if(!valid) return ;
           
                const res = await otherset.saveNetwork({
                            ip:this.networkform.address,
                            mask:this.networkform.netmask,
                            gateway:this.networkform.gateway,
                            dnsServer:this.networkform.dns,
                            standbyDnsServer:this.networkform.spare_dns,
                            webPort:this.networkform.port,
                        });
                if(res.code === 200){
                    this.$message({
                        message:'保存成功',
                        type: 'success'
                    });
                    this.get_network_config(); 
                }else{
                    this.$message({
                        message:dt.msg,
                        type: 'error'
                    });
                }
            })
        },

    }
})
 
</script>

<style scoped lang="less">

    .innerinput{
        width:calc(100% - 120px);
    }
    .el-button--primary.is-disabled{
        background-color: rgba(76, 122, 227, .5) !important;
        border-color: rgba(76, 122, 227, .1) !important;
    }
        .title_left_line{
            font-size: 14px;
            margin-bottom: 16px;
            font-weight: 500;
            border-left: 3px solid #4C7AE3;
            padding-left: 10px;
            height: 14px;
            line-height: 14px;
            color: rgba(72, 72, 102, 0.87);
            >span{
                padding-left: 32px;
                font-size: 13px;
                color: rgba(72, 72, 102, 0.63);
            }
            margin-bottom: 16px;
        }
        .div_block{
            height: 298px;
            background-color: #fff;
            box-shadow: 0px 2px 4px 1px rgba(76, 122, 227, 0.11);
            border-radius: 4px;
            margin-bottom: 24px;
            padding: 24px;
            box-sizing: border-box;
            /deep/ .el-form-item__label{
                height:32px;
                line-height: 32px;
            }
            /deep/ .el-form-item__content{
                height:32px;
                line-height: 32px;
            }
            .div_blockbtn{
                margin-left: 0;
                width:82px;
                height:32px;
                line-height: 7px;
                padding-right: 40px;
                text-align: center;
                padding-left: 14px;
                margin-left: 10px;
            }
        }
        .updatetime{
            float: right!important;
            border-left: none!important;
        }
        /deep/ .el-checkbox__label{
            font-size: 13px;
        }
        /deep/ .el-form-item__label{
            /* color: rgba(72, 72, 102, 0.64);
            font-weight: 500;
            font-size: 13px; */
            font-size: 14px;
            font-weight: 500;
            color: rgba(72,72,102,0.87);
        }
        /deep/ .el-form-item__content{
            color: rgba(72, 72, 102, 0.64);
            font-weight: 500;
            font-size: 13px;
        }
        /deep/ .el-dialog{
            height: 192px !important;
        }
        .dialogtxt{
            text-align: center;
            margin-top: 55px;
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
        /deep/ .el-table .cell, .el-table th div, .el-table--border td:first-child .cell, .el-table--border th:first-child .cell{
            padding-left:32px;
        }
        
         
    .upgradeactive{
        display: block !important;
    }
    .systemupgrade{
        display: none;
    }
    .knowledgeupgrade{
        display: none;
    }
    .sysupgradebox{
        /* margin-bottom: 24px; */
        .spversion{
            label{
                display: inline-block;
                width: 150px;
            }
            span{
                display: inline-block;
                // width: 200px;
            }
        }
    }
    .context_box_bg{
        background: none;
    }
    /deep/ .el-tabs__header {
        margin: 0 0 15px;
        background: #fff;
        border-radius: 4px;
        box-shadow: 0px 2px 4px 0px rgb(76 122 227 / 12%);
        border: none;
    }
    .systembox{
        //
        // background: #fff;
        min-height: calc(100% - 39px);
        box-sizing: border-box;
        // padding: 24px;
        // box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
        /deep/ .el-tabs{
            // min-height: 700px;
        } 
        /deep/ .el-tabs__item.is-top.is-active{
            // background: #fff;
        }
        /deep/ .el-tabs__header{
            margin: 0 0 15px;
            background: #fff;
            border-radius: 4px;
            box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
            border: none;
        }
        /deep/ .el-tabs__content{
            // padding: 24px ; 
            // background: #fff;
            // min-height: 680px;
            margin-top: 24px;
        }
        .tabsbox{
            background: #fff;
        }
    } 
    /deep/ .el-switch.is-checked .el-switch__core{
        border-color: #4C7AE3;
            background-color: #4C7AE3;
    }
    .system{ 
        div {
            float:left;
            text-align: center;
        }
    }
    .iptw{
        width: 500px !important;
    }
    .systeminfolist{
        width: 88%;
        /* border-top: 1px solid #E8E8F5; */
        border-left: 1px solid #E8E8F5;
        border-right: 1px solid #e8e8f5;
        .planTime{
            height: 50px;
            line-height:50px;
            
        }
        .switchbtn{
            float: right;
            margin-top: 15px;
            margin-right: 20px;
        }
        .blueword{
            color:#4C7AE3;
            cursor: pointer;
        }
        .timeword{
            margin-left: 10px;
        }
        .system{
            // display: flex;
            // align-items: center;
            // justify-content: space-between;
            font-weight: 500;
            line-height: 30px;
            padding: 10px 20px;
            padding: 10px 20px;
        line-height: 30px;
            width: calc(100% - 160px);
        float: left;
        box-sizing: border-box;
            .system_data{
                background: #fff;
                color: rgba(72, 72, 102, 0.64);
                font-weight: 500;
    
            }
        }
    } 
    
    /deep/ .el-input--small .el-input__icon {
        line-height: 52px;
    }
    /deep/ .el-checkbox__label {
        font-size: 13px;
        color: rgba(72, 72, 102, 0.64);
    }
    /deep/ .el-form-item__label{
     text-align: left;
    }
    </style>
