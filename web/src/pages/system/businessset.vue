
<!-- 系统设置---系统工具-- 业务设置-->
<template>
    <div style="height: 100%; overflow-y: auto; overflow-x: hidden;">
        <el-row :gutter="20">
<!--
            <el-col :span="12">
                <div class="title_left_line">
                    <label>TCP盲测平台</label> 
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            系统和自定义的TCP盲测平台只能同时使用一个。修改系统的TCP盲测平台配置将直接修改所有渗透节点上的盲测配置。使用自定义TCP盲测平台，系统将不再显示检测的盲测的漏洞结果。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block">
                    <el-form ref="tcpform" :model="tcpform" label-width="140px" class="sysform"  :rules="rules2">
                        <div style="margin-bottom:24px">
                            <el-radio-group v-model="tcpform.type" @change="changeTcpType">
                                <el-radio :label="1">使用系统TCP盲测平台</el-radio>
                                <el-radio :label="2">使用自定义TCP盲测平台</el-radio>
                            </el-radio-group>
                        </div> 
                        <el-form-item label="监听IP：" class="syswarnvalue" prop="ip" >
                            <el-input v-model="tcpform.ip" 
                                style="width:calc(100% - 190px)" v-if="tcpform.type == 2">
                            </el-input>   
                            <span v-if="tcpform.type == 1">每个渗透节点的IP地址为监听IP地址</span>
                        </el-form-item>
                        <el-form-item label="监听端口：" class="syswarnvalue" prop="port">
                            <el-input v-model="tcpform.port" :disabled="tcpform.type === 1"  
                                style="width:calc(100% - 190px)">
                            </el-input>  
                            
                        </el-form-item>
                    </el-form>
                    <el-button type="primary" class="div_blockbtn"  @click="saveTCP" >保存设置</el-button>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="title_left_line">
                    <label>HTTP盲测平台</label>
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            系统和自定义的HTTP盲测平台只能同时使用一个。修改系统的HTTP盲测平台配置将直接修改所有渗透节点上的盲测配置。使用自定义HTTP盲测平台，系统将不再显示检测的盲测的漏洞结果。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block">
                    <el-form ref="httpform" :model="httpform" label-width="140px" class="sysform" :rules="rules3">
                        <div style="margin-bottom:24px">
                            <el-radio-group v-model="httpform.type" @change="changeHTTPType">
                                <el-radio :label="1">使用系统HTTP盲测平台</el-radio>
                                <el-radio :label="2">使用自定义HTTP盲测平台</el-radio>
                            </el-radio-group>
                        </div>
                        <el-form-item label="监听IP：" class="syswarnvalue" prop="ip"  >
                            <el-input v-model="httpform.ip" style="width:calc(100% - 190px)" v-if="httpform.type == 2">
                            </el-input>
                            <span v-if="httpform.type == 1">每个渗透节点的IP地址为监听IP地址</span>
                        </el-form-item>
                        <el-form-item label="监听端口：" class="syswarnvalue" prop="port"  >
                            <el-input v-model="httpform.port" :disabled="httpform.type === 1" style="width:calc(100% - 190px)">
                            </el-input>
            
                        </el-form-item>
                    </el-form>
                    <el-button type="primary" class="div_blockbtn" @click="saveHTTP">保存设置</el-button>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="title_left_line">
                    <label>DNS盲测平台</label> 
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            系统和自定义的DNS盲测平台只能同时使用一个。系统的DNS盲测平台配置不可修改。使用自定义DNS盲测平台，系统将不再显示检测的盲测的漏洞结果。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block">
                    <el-form ref="dnsform" :model="dnsform"  label-width="150px" class="sysform" :rules="rules4">
                        <div style="margin-bottom:24px">
                            <el-radio-group v-model="dnsform.type" @change="changeDnsType">
                                <el-radio :label="1">使用系统DNS盲测平台</el-radio>
                                <el-radio :label="2">使用自定义DNS盲测平台</el-radio>
                            </el-radio-group>
                        </div>
                        <el-form-item label="解释域名：" class="syswarnvalue " prop="address">
                            <el-input v-model="dnsform.address2" class="innerinput"
                                :disabled="dnsform.type ==1" v-if="dnsform.type ==1">
                            </el-input>
                            <el-input v-model="dnsform.address" class="innerinput"
                                :disabled="dnsform.type ==1" v-else>
                            </el-input>
                        </el-form-item> 
                    </el-form>
                    <el-button type="primary" style="margin-top:54px" class="div_blockbtn" @click="saveDNS">保存设置</el-button>
                </div>
            </el-col>


<el-col :span="12">
                <div class="title_left_line">
                    <label>icmp盲测平台</label> 
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            系统和自定义的icmp盲测平台只能同时使用一个。系统的icmp盲测平台配置不可修改。使用自定义icmp盲测平台，系统将不再显示检测的盲测的漏洞结果。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block">
                    <el-form ref="icmpform" :model="icmpform"  label-width="150px" class="sysform" :rules="rules4">
                        <div style="margin-bottom:24px">
                            <el-radio-group v-model="icmpform.type" @change="changeicmpType">
                                <el-radio :label="1">使用系统icmp盲测平台</el-radio>
                                <el-radio :label="2">使用自定义icmp盲测平台</el-radio>
                            </el-radio-group>
                        </div>
                        <el-form-item label="盲测地址：" class="syswarnvalue " prop="address">
                            <el-input v-model="icmpform.address2" class="innerinput"
                                :disabled="icmpform.type ==1" v-if="icmpform.type ==1">
                            </el-input>
                            <el-input v-model="icmpform.address" class="innerinput"
                                :disabled="icmpform.type ==1" v-else>
                            </el-input>
                        </el-form-item> 
                    </el-form>
                    <el-button type="primary" style="margin-top:54px" class="div_blockbtn" @click="saveicmp">保存设置</el-button>
                </div>
            </el-col>
-->

            <el-col :span="12">
                <div class="title_left_line">
                    <label>任务并发配置</label> 
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            配置系统并发执行渗透测试任务的数量。
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block">
                    <el-form ref="concurrencyform" :model="concurrencyform" label-width="150px" class="sysform"> 
                        <el-form-item label="系统最大并发IP：" class="syswarnvalue " >
                            <el-input type="number" min=1
                                v-model="concurrencyform.ipcount" class="innerinput" style="width:calc(100% - 190px)" 
                            onKeypress="return(/^[1-9]*$/.test(String.fromCharCode(event.keyCode)))">
                            </el-input>
                            <el-tooltip class="item" effect="dark" placement="right">
                                <div slot="content">
                                    系统最大并发IP在出厂时已经确定。
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                            </el-tooltip>
                        </el-form-item>
                        <el-form-item label="系统并发任务：" class="syswarnvalue "  >
                            <el-input type="number"  min=1 v-model="concurrencyform.taskcount" class="innerinput" style="width:calc(100% - 190px)"  
                            onKeypress="return(/^[1-9]*$/.test(String.fromCharCode(event.keyCode)))">
                            </el-input>
                            <el-tooltip class="item" effect="dark" placement="right">
                                <div slot="content">
                                    多任务并发执行时，系统会根据并发任务数量将最大IP检测能力均衡到每个任务中，因此设置系统并发任务的数量不能小于最大并发数量。
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                            </el-tooltip>
                        </el-form-item>
                    </el-form>
                    <el-button type="primary" style="margin-top:10px" class="div_blockbtn" @click="saveconcurrency">保存设置</el-button>
                </div>

            </el-col>
            <el-col :span="12">
                <div class="title_left_line">
                    <label>远控监听</label>
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            系统和自定义的远控监听只能同时使用一个。修改系统的远控监听将修改所有渗透节点上的远控监听配置。使用自定义的远控监听配置，远控将反向连接到自定义的监听地址上。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block">
                    <div style="margin-bottom:24px">
                        <el-radio-group v-model="controlform.type" @change="changeControl">
                            <el-radio :label="1">使用系统远控监听</el-radio>
                            <el-radio :label="2">使用自定义远控监听</el-radio>
                        </el-radio-group>
                    </div>
                    <el-form ref="controlform2" :model="controlform"  label-width="140px"
                        style="text-align: left;" class="sysform" :rules="rules1" v-if="controlform.type == 1">
                        <el-form-item label="监听IP：" class="syswarnvalue">
                            <span>每个渗透节点的IP地址为监听IP地址</span>
                        </el-form-item>
                        <el-form-item label="监听端口：" class="syswarnvalue" prop="port">
                            <el-input v-model="controlform.port" :disabled="true" style="width:calc(100% - 190px)">  
                            </el-input>
                        </el-form-item>
                    </el-form>
                    <el-form ref="controlform2" :model="controlform"  label-width="140px"
                        style="text-align: left;" class="sysform" :rules="rules1"  v-if="controlform.type == 2">
                        <el-form-item label="监听IP：" class="syswarnvalue" prop="ip">
                            <el-input v-model="controlform.ip" style="width:calc(100% - 190px)">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="监听端口：" class="syswarnvalue" prop="port">
                            <el-input v-model="controlform.port" style="width:calc(100% - 190px)">  
                            </el-input>
                        </el-form-item>
                    </el-form>
                    <el-button type="primary" class="div_blockbtn" @click="savereverse">保存设置</el-button>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="title_left_line">
                    <label>测试目标黑白名单</label>
                    <span>
                        <!-- @change="targetipsave" -->
                        <el-switch v-model="targetip.is_Open" class="elSwitch" >
                        </el-switch>
                    </span>
                    <!-- <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            系统和自定义的远控监听只能同时使用一个。修改系统的远控监听将修改所有渗透节点上的远控监听配置。使用自定义的远控监听配置，远控将反向连接到自定义的监听地址上。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip> -->
                </div>
                <div class="div_block">
                    <el-form ref="targetip" :model="targetip"  label-width="140px"
                        style="text-align: left;" class="sysform ">
                        <div style="margin-bottom:24px"> 
                            <el-radio v-model="targetip.type" :label="2" style="margin-bottom: 12px;">测试目标黑名单</el-radio>
                            <el-input 
                                type="textarea"
                                resize="none"
                                v-model="targetip.blackList" >  
                            </el-input>
                            <el-radio v-model="targetip.type" :label="1" style="margin-bottom: 12px;margin-top:18px">测试目标白名单</el-radio>
                            <el-input 
                                type="textarea"
                                resize="none"
                                v-model="targetip.whiteList" >  
                            </el-input>
 
                        </div>
                         
                    </el-form>
                    <el-button type="primary" class="div_blockbtn" @click="targetipsave">保存设置</el-button>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="title_left_line">
                    <label>测试范围校验</label>
                    <el-tooltip class="item" effect="dark" placement="right">
                        <div slot="content">
                            开启后，系统将自动校验您输入的测试范围（如IP、域名、URL等）是否符合规范格式，避免因目标格式错误导致扫描失败。<br />
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;"></i>
                    </el-tooltip>
                </div>
                <div class="div_block" style="height:106px;">
                    <span>
                        <el-switch v-model="AvailabilityScore" class="elSwitch" @change="changeSysstatus2">
                        </el-switch>
                    </span>
                </div>
            </el-col>
            
            
        </el-row>
    </div>
</template>

<script>
import { businessset } from '@/api/system.js'
import { otherset ,system} from '@/api/system.js'

import $ from 'jquery'
export default ({
    name: 'businessset',
    components: {
    },
    data() {
        return {
            AvailabilityScore: false,
            tcpform: {
                id:'',
                type:1,
                ip: '', 
                port:'6888',
            }, 
            httpform: {
                id: '',
                type: 1,
                ip: '',
                port: '8111',
            }, 
            whitelistinfo: '',
            dnsform: { 
                id: '',
                type:1,
                address: '', 
                address2: 'dnstunnel.run'
            },
            icmpform:{
                id: '',
                type:1,
                address: '', 
                address2: ''
            },
            concurrencyform:{
                id: '',
                ipcount:16,
                taskcount:4,
            },
            controlform:{
                id: '',
                type: 1,
                ip: '',
                port: '6666',
            },
            targetip:{
                type:1,
                blackList:'',
                whiteList:'',
                isOpen:2,  //1开启，2关闭
                is_Open:false, 
            },
            rules1: {
                ip: [
                    { required: true, message: '监听IP不能为空', trigger: 'blur' },
                ],
                port: [
                    { required: true, message: '监听端口不能为空', trigger: 'change' },
                ],
            },
            rules2: {
                ip: [
                    { required: true, message: '监听IP不能为空', trigger: 'blur' },
                ],
                port: [
                    { required: true, message: '监听端口不能为空', trigger: 'blur' },
                ],
            },
            rules3: {
                ip: [
                    { required: true, message: '监听IP不能为空', trigger: 'blur' },
                ],
                port: [
                    { required: true, message: '监听端口不能为空', trigger: 'blur' },
                ],
            },
            rules4: {
                address: [
                    { required: true, message: '解释域名不能为空', trigger: 'blur' },
                ],
            },
         
        }
    },
    created: function () {
        this.getTCP();
        this.getHTTP();
        this.getDNS();
        this.getconcurrency();
        this.getreverse();
        this.getTargetip();
    },
    mounted: function () {
     this.getscoredata();

    },
    methods: {
           changeSysstatus2 (value){
                this.getscoredataSave();
        },
            async getscoredata(){ //测试范围校验信息
            
            const res = await system.gettestscopeinfo();
            if(res.code === 200){  
               console.log(res.data.isOpen,'测试范围校验信息');
               this.AvailabilityScore = res.data.isOpen === 1 ? true : false;
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        async getscoredataSave(){ //测试范围校验保存
            const res = await system.gettestscopesave({isOpen:this.AvailabilityScore ? 1 : 2});
            if(res.code === 200){  
                this.getscoredata()
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        async getTCP(){
            const res = await businessset.getTCP();
            if(res.code == 200){
                // this.tcpform.id = res.data.id,
                this.tcpform.type = res.data.type;
                this.tcpform.ip = res.data.host;
                this.tcpform.port = res.data.port;
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async saveTCP(){
            if (this.tcpform.type == 1) {
                this.tcpform.ip = '';
            }
            const res = await businessset.saveTCP({
                id:this.tcpform.id,
                type:this.tcpform.type,
                host:this.tcpform.ip,
                port:this.tcpform.port,
            });
            if(res.code == 200){
                this.$message({
                    message: '保存成功',
                    type: 'success'
                });
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async getHTTP() {
            const res = await businessset.getHTTP();
            if (res.code == 200) {
                this.httpform.id = res.data.id,
                this.httpform.type = res.data.type;
                this.httpform.ip = res.data.host;
                this.httpform.port = res.data.port;
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async changeHTTPType(value) {
            if (value === 1) {
                this.httpform.port = '8111'
            } else {
                this.httpform.port = ''
                const res = await businessset.getHTTP();
                if (res.code == 200) {
                    this.httpform.port = res.data.port;
                } else {
                    this.$message({
                        message: res.msg,
                        type: 'error'
                    });
                }
            }
        },
        async saveHTTP(){
            if (this.httpform.type == 1) {
                this.httpform.ip = '';
            }
            const res = await businessset.saveHTTP({
                id: this.httpform.id,
                type: this.httpform.type,
                host: this.httpform.ip,
                port: this.httpform.port,
            });
            if (res.code == 200) {
                this.$message({
                    message: '保存成功',
                    type: 'success'
                });
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async getDNS(){
            const res = await businessset.getDNS();
            if (res.code == 200) {
                this.dnsform.id = res.data.id;
                this.dnsform.type = res.data.type;
                this.dnsform.address = res.data.domain;
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async saveDNS(){
            const res = await businessset.saveDNS({
                id: this.dnsform.id,
                type: this.dnsform.type,
                domain: this.dnsform.address,
            });
            if (res.code == 200) {
                this.$message({
                    message: '保存成功',
                    type: 'success'
                });
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async saveicmp(){

        },
        async getconcurrency(){
            const res = await businessset.getconcurrency();
            if (res.code == 200) {
                // this.concurrencyform.id = res.data.id;
                this.concurrencyform.ipcount = res.data.curIp;
                this.concurrencyform.taskcount = res.data.curTasks;
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async saveconcurrency(){
            const res = await businessset.saveconcurrency({ 
                curIp: this.concurrencyform.ipcount,
                curTasks: this.concurrencyform.taskcount,
            });
            if (res.code == 200) {
                this.$message({
                    message: '保存成功',
                    type: 'success'
                });
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async getreverse(){
            const res = await businessset.getreverse();
            if (res.code == 200) {
                this.controlform.id = res.data.id;
                this.controlform.type = res.data.reverseType;
                this.controlform.ip = res.data.reverseHost;
                this.controlform.port = res.data.reversePort;
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async savereverse() {
            if (this.controlform.type == 2) {
                this.$refs.controlform2.validate( async (valid) => {
                    if (valid) { 
                        const res = await businessset.savereverse({
                            id: this.controlform.id,
                            reverseType: this.controlform.type,
                            reverseHost: this.controlform.ip,
                            reversePort: this.controlform.port,
                        });
                        if (res.code == 200) {
                            this.$message({
                                message: '保存成功',
                                type: 'success'
                            });
                        } else {
                            this.$message({
                                message: res.msg,
                                type: 'error'
                            });
                        }
                    }
                });
            } else {
                this.controlform.ip = '';
                const res = await businessset.savereverse({
                    id: this.controlform.id,
                    reverseType: this.controlform.type,
                    reverseHost: this.controlform.ip,
                    reversePort: this.controlform.port,
                });
                if (res.code == 200) {
                    this.$message({
                        message: '保存成功',
                        type: 'success'
                    });
                } else {
                    this.$message({
                        message: res.msg,
                        type: 'error'
                    });
                }
            }
              
        },
        async changeTcpType (value) {
            if (value === 1) {
                this.tcpform.port = '6888'
            } else {
                this.tcpform.port = ''
                const res = await businessset.getTCP();
                if(res.code == 200){
                    this.tcpform.port = res.data.port;
                }else{
                    this.$message({
                        message: res.msg,
                        type: 'error'
                    });
                }
            } 
        },
     
        async changeDnsType (value) {
            if (value === 1) {
                this.dnsform.address2 = 'dnstunnel.run'
            } else {
            this.dnsform.address = ''
            const res = await businessset.getDNS();
            if (res.code == 200) {
                this.dnsform.address = res.data.domain;
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
            }
        },
        changeicmpType(value){
            if (value === 1) {
                this.icmpform.address2 = ''
            } else {
                this.icmpform.address = '' 
            }
        },
        async changeControl (value) {
            if (value === 1) {
                this.controlform.port = '6666'
                 this.$nextTick(()=>{
                    this.$refs.controlform.clearValidate(['ip']);
                })
            } else {
                const res = await businessset.getreverse();
                if (res.code == 200) {
                    this.controlform.ip = res.data.reverseHost;
                    this.controlform.port = res.data.reversePort;
                } else {
                    this.$message({
                        message: res.msg,
                        type: 'error'
                    });
                }
            }
        },
        async getTargetip(){
            const res = await businessset.getTargetIp();
            if (res.code == 200) { 
                 if(res.data.isOpen == 1){
                    this.targetip.is_Open =true;
                 }else{
                    this.targetip.is_Open =false;
                 }
                 this.targetip.type = res.data.type;
                 this.targetip.blackList = res.data.blackList;
                 this.targetip.whiteList = res.data.whiteList;
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async targetipsave(){
            let isOpen = this.targetip.is_Open?1:2;
            const res = await businessset.saveTargetIp({
                isOpen: isOpen,
                type: this.targetip.type,
                whiteList:  this.targetip.whiteList,
                blackList: this.targetip.blackList,
            });
            if (res.code == 200) {
                this.$message({
                    message: '保存成功',
                    type: 'success'
                });
                this.getTargetip();
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        
    }
})

</script>

<style scoped lang="less">
/deep/ .el-tooltip__popper.is-dark{
    width: auto;
}
.innerinput {
    width: calc(100% - 120px);
}

.el-button--primary.is-disabled {
    background-color: rgba(76, 122, 227, .5) !important;
    border-color: rgba(76, 122, 227, .1) !important;
}

.title_left_line {
    font-size: 14px;
    margin-bottom: 16px;
    font-weight: 500;
    border-left: 3px solid #4C7AE3;
    padding-left: 10px;
    height: 14px;
    line-height: 14px;
    color: rgba(72, 72, 102, 0.87);

    >span {
        padding-left: 32px;
        font-size: 13px;
        color: rgba(72, 72, 102, 0.63);
    }

    margin-bottom: 16px;
}

.div_block {
    height: 298px;
    background-color: #fff;
    box-shadow: 0px 2px 4px 1px rgba(76, 122, 227, 0.11);
    border-radius: 4px;
    margin-bottom: 24px;
    padding: 24px;
    box-sizing: border-box;

    // /deep/ .el-form-item__label {
    //     height: 32px;
    //     line-height: 32px;
    // }

    // /deep/ .el-form-item__content {
    //     height: 32px;
    //     line-height: 32px;
    // }

    .div_blockbtn {
        margin-left: 0;
        width: 82px;
        height: 32px;
        line-height: 7px;
        padding-right: 40px;
        text-align: center;
        padding-left: 14px;
        // margin-left: 10px;
    }
}

.updatetime {
    float: right !important;
    border-left: none !important;
}

/deep/ .el-checkbox__label {
    font-size: 13px;
}

/deep/ .el-form-item__label {
    /* color: rgba(72, 72, 102, 0.64);
            font-weight: 500;
            font-size: 13px; */
    font-size: 14px;
    font-weight: 500;
    color: rgba(72, 72, 102, 0.87);
}

/deep/ .el-form-item__content {
    color: rgba(72, 72, 102, 0.64);
    font-weight: 500;
    font-size: 13px;
}

/deep/ .el-dialog {
    height: 192px !important;
}

.dialogtxt {
    text-align: center;
    margin-top: 55px;
}

/deep/ .el-tabs__item {
    height: 48px;
    line-height: 48px;
    padding: 0 24px;
}

/deep/ .el-tabs__item.is-active {
    color: #4C7AE3;
    font-weight: 500;
}

/deep/ .el-tabs__nav-wrap {
    padding: 0 24px;
}

/deep/ .el-tabs__nav-wrap::after {
    background: #E8E8F5;
    height: 1px;
}

/deep/ .el-tabs__header {
    margin: 0 0 24px;
}

/deep/ .el-table .cell,
.el-table th div,
.el-table--border td:first-child .cell,
.el-table--border th:first-child .cell {
    padding-left: 32px;
}


.upgradeactive {
    display: block !important;
}

.systemupgrade {
    display: none;
}

.knowledgeupgrade {
    display: none;
}

.sysupgradebox {

    /* margin-bottom: 24px; */
    .spversion {
        label {
            display: inline-block;
            width: 150px;
        }

        span {
            display: inline-block;
            // width: 200px;
        }
    }
}

.context_box_bg {
    background: none;
}

/deep/ .el-tabs__header {
    margin: 0 0 15px;
    background: #fff;
    border-radius: 4px;
    box-shadow: 0px 2px 4px 0px rgb(76 122 227 / 12%);
    border: none;
}

.systembox {
    //
    // background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;

    // padding: 24px;
    // box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
    /deep/ .el-tabs {
        // min-height: 700px;
    }

    /deep/ .el-tabs__item.is-top.is-active {
        // background: #fff;
    }

    /deep/ .el-tabs__header {
        margin: 0 0 15px;
        background: #fff;
        border-radius: 4px;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
        border: none;
    }

    /deep/ .el-tabs__content {
        // padding: 24px ; 
        // background: #fff;
        // min-height: 680px;
        margin-top: 24px;
    }

    .tabsbox {
        background: #fff;
    }
}

/deep/ .el-switch.is-checked .el-switch__core {
    border-color: #4C7AE3;
    background-color: #4C7AE3;
}

.system {
    div {
        float: left;
        text-align: center;
    }
}

.iptw {
    width: 500px !important;
}

.systeminfolist {
    width: 88%;
    /* border-top: 1px solid #E8E8F5; */
    border-left: 1px solid #E8E8F5;
    border-right: 1px solid #e8e8f5;

    .planTime {
        height: 50px;
        line-height: 50px;

    }

    .switchbtn {
        float: right;
        margin-top: 15px;
        margin-right: 20px;
    }

    .blueword {
        color: #4C7AE3;
        cursor: pointer;
    }

    .timeword {
        margin-left: 10px;
    }

    .system {
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

        .system_data {
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

/deep/ .el-form-item__label {
    text-align: left;
}
</style>
