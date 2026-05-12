<template>
    <el-dialog :title="title" :visible.sync="bugdialogVisible" width="1184px" class="buginfobox"
        :close-on-click-modal="false" :show-close="false">
        <div class="dialog_b_btn">
            <!-- <el-button size="small" @click="btnUpdatebug">{{ updateBugtxt }}</el-button> -->
            <!-- <el-button size="small" @click="saveUpdatebug" v-if="is_bugUpdate">保存</el-button> -->
            <el-button size="small" @click="cancalbugdialogVisible">关闭</el-button>
        </div>
        <div class="buginfo_box">
<!-- 上table......................................................................................................... -->
            <div class="bugbasicinfo">
                <el-table :data="bugbasicinfo" size='small' style="width: 100%">
                    <el-table-column prop="ip" label="IP"> 
                    </el-table-column>
                    <el-table-column prop="host" label="Host">
                    </el-table-column>
                    <el-table-column prop="port" label="端口"> 
                    </el-table-column>

                    <el-table-column prop="riskLevel" label="风险等级">
                        <template slot-scope="scope">
                            <span v-if="!is_bugUpdate">
                                <span :class="[
                                { 'riskstyle risk_hight': scope.row.riskLevel == 1 },
                                { 'riskstyle risk_middle': scope.row.riskLevel == 2 },
                                { 'riskstyle risk_low': scope.row.riskLevel == 3 },
                                { 'riskstyle risk_nofind': scope.row.riskLevel == 4 }]"><i></i>{{ scope.row.riskLevelName }}</span>
                            </span> 
                        </template>
                    </el-table-column>
                    <el-table-column prop="createTime" label="发现时间">
                    </el-table-column>
                    <el-table-column prop="createTime" label="更新时间">
                    </el-table-column>


                </el-table>
            </div>
<!-- 下table.................................................................................................................................. -->
            <div class="bugbasicinfo" style="margin-top:24px">
                <el-table :data="bugbasicinfo" size='small' style="width: 100%">
                    <el-table-column prop="id" label="ID"> 
                    </el-table-column>
                    <el-table-column prop="token" label="反连 Tocken"> 
                    </el-table-column>
                    <el-table-column prop="hash" label="Hash">
                    </el-table-column>
                    <el-table-column prop="type" label="类型">
                    </el-table-column>
                    <el-table-column prop="origin" label="来源">
                    </el-table-column>
                    <el-table-column prop="status" label="验证状态"> 
                    </el-table-column>
        
                </el-table>
            </div>
    <!-- 最下信息展示.................................................................................................................................. -->
    <div class="lastbasicinfo" style="margin-top:24px">
        <div class="bugotherinfo">
            <div class="lastline"><span class="firsttit">URL:</span><span>{{buginfo.url}}</span></div>
            <div class="lastline" style="border-bottom: none;"><span class="firsttit">Parameter:</span><span>{{ buginfo.parameter }}</span></div>
            <el-collapse v-model="activeNames" @change="handleChange">
                <el-collapse-item title="详情:" name="3">
                    <div>{{ buginfo.detail }}</div>
                </el-collapse-item>
                <!-- <el-collapse-item title="请求:" name="1">
                   <div>{{ buginfo.request }}</div>
                </el-collapse-item>
                <el-collapse-item title="响应:" name="2">
                    <div>{{ buginfo.response }}</div>
                </el-collapse-item>
                 -->
            </el-collapse>
            <div>
                <!-- <div class="part_title">报文</div> -->
                <div class="btnDiv">   
                    <el-button  size="mini" type="primary" @click="btnTest()"  >测试 </el-button> 
                </div>
            </div>
            <el-row style="margin-top:10px" :gutter="20"  >
                <el-col :span="12">
                    <div class="message requestpack">
                        <label class="title_bg title_bg1">请求报文</label>
                        <div> 
                            <div class="packheight " style="height:100%">
                                <div class="packinput" style="height:100%">
                                    <!--    -->
                                    <div @click="isShowInput= true" v-show="!isShowInput" style=" background: #fff;">
                                        <pre v-dompurify-html="highlightedRequest(buginfo.request,buginfo.payload) "></pre>
                                    </div>
                                    <el-input v-show="isShowInput" class="packtxt" type="textarea" :rows="7"
                                        v-model=" buginfo.request"  
                                        @change = "highlightedRequest( buginfo.request,buginfo.payload)" 
                                        resize="none"  >
                                    </el-input>
                                    
                                </div>
                            </div>
                        </div>
                    </div>
                </el-col>
                <el-col :span="12">
                    <div class="message">
                        <label class="title_bg title_bg2">响应报文</label>
                        <div style="    background: #fff;">
                            <!-- <pre>{{ buginfo.response }}</pre> -->
                                <div style="    background: #fff;">
                                        <pre v-dompurify-html="highlightedRequest2(buginfo.response,buginfo.payload_success_flag) "></pre>
                                </div>
                        </div>
                    </div>
                </el-col>
            </el-row>

        </div>
    </div>        

        </div>
    </el-dialog>
</template>
<style lang="less" scoped>
.lastbasicinfo{
     background: #fff;
     border: 1px solid rgba(232, 232, 245, 1);
     .lastline{
        height:50px;
        color: rgba(72, 72, 102, 0.87);
        padding-left:50px;
        line-height:50px;
        border-bottom: 1px solid rgba(232, 232, 245, 1);
        .firsttit{
            display: inline-block;
            margin-right: 50px;
            width:100px;
        }
     }
}
/deep/ .el-collapse-item__header{
    color: rgba(72, 72, 102, 0.87);
    padding-left:50px;
}
/deep/ .el-collapse-item__content{
    color: rgba(72, 72, 102, 0.87);
    padding-left:50px;
}
/deep/ .el-collapse-item__arrow{
    margin: 0 8px 0 -60px;
}

.title_bg {
    width: 84px;
    height: 32px;
    font-size: 13px;
    font-weight: 500;
}

.title_bg1 {
    background-color: rgba(243, 95, 40, 0.12) !important;
    border: 1px solid rgba(24, 144, 255, 0.08);
    color: #F35F28 !important;
    border-left: 3px solid #F35F28;

}

.title_bg2 {
    background-color: rgba(76, 122, 227, 0.12) !important;
    border: 1px solid rgba(24, 144, 255, 0.08);
    color: #4C7AE3 !important;
    border-left: 3px solid #4C7AE3;
}

 .bugbasicinfo {
     padding: 24px;
     background: #fff;
     border: 1px solid rgba(232, 232, 245, 1);
 }

 .buginfo_box {
     padding: 24px;
 }
 .title_bg {
    width: 84px;
    height: 32px;
    font-size: 13px;
    font-weight: 500;
}

.title_bg1 {
    background-color: rgba(243, 95, 40, 0.12) !important;
    border: 1px solid rgba(24, 144, 255, 0.08);
    color: #F35F28 !important;
    border-left: 3px solid #F35F28;

}

.title_bg2 {
    background-color: rgba(76, 122, 227, 0.12) !important;
    border: 1px solid rgba(24, 144, 255, 0.08);
    color: #4C7AE3 !important;
    border-left: 3px solid #4C7AE3;
}
 .packinput { 
     box-sizing: border-box !important;

     .packtxt {
         height: 100%;

         /deep/ textarea {
             height: 100%;
             border: 0 !important;
         }
     }
 }
  .requestpack>div {
     background: #fff !important;
     padding: 0 !important;
 }
 .message .title_bg {
     margin-bottom: 8px;
 }

 .message>label {
     display: inline-block;
     width: 80px;
     text-align: center;
     height: 26px;
     line-height: 26px;
     color: #fff;
     background-color: #4c7ae3;
     font-weight: bold;
     font-size: 12px;
 }

 .message>div {
     height: 253px;
     overflow-y: auto;
 }
 .message>div {
       // margin-bottom: 24px; 
       background: #F7F7FB;
       border-radius: 4px;
       border: 1px solid #E8E8F5;
       padding: 16px;
       box-sizing: border-box;
   }
</style>
<script>
import { traffic } from '@/api/traffic.js'
import { logic } from '@/api/task.js'
import {Base64} from 'js-base64';
export default {
    name:'vulnmsginfo',
    components: {
        
    },
    props:{
        value: {}, // 注意此处获取的value对应的就是组件标签中的v-model
        vulninfo:{},
        buglevel:Array,
        vulthreatlist: Array,
        vuln_id:{},
        title:{}
    },
    data(){
        return{
            activeNames: [''],
            isShow1: false,
            isShow2: false,
            is_bugUpdate: false,
            bugdialogVisible:false,
            updateinfo: {
                check_vuln_id: '',
                name: '',
                object: '',
                type: '',
                risk_lever: '',
                cve: '',
                cnvd: '',
                cnnvd: '',
                description: '',
                fix_suggest: '',
                ref_url: '',
                vuln_location: '',
                risk_lever_number: '',
                status: [],
                use_impact_label: '',
                use_impact_value: '',
                use_impact_label: '',
                use_impact_value: '',
                priority_description: '',
                vul_analysis: '',
                code: '',
                affect_range: '',
            },
            buglisttableData: [],
            updateBugtxt: '编辑',
            bugbasicinfo: [],
            buginfo:{
                request:'',
                response:'',
                parameter:'',
                detail:'',
                payload:'',
                payload_success_flag:'',
            },
            bugmessage: [], 
            // buglevel: [],
            expands: [],  // 要展开的行，数值的元素是row的key值
            responsepack: '',
            requestpack: '',
            isShowInput:false,
        }
    },
    created() {
        
    },
    watch:{ 
        value(newVal, oldVal) { 
            // 监测value的变化，并赋值。
            if (newVal){
                this.getVulninfo(); 
                this.bugdialogVisible = newVal;
 
            } 
            
        },
        bugdialogVisible(val) {
            this.$emit("input", val); // 此处监测showMask目的为关闭弹窗时，重新更换value值，注意emit的事件一定要为input。
        }
    },
    mounted() {
        this.bugdialogVisible = this.value; // 在生命周期中，把获取的value值获取给bugdialogVisible
        
    },
    methods: {
        handleChange(val) {
            console.log(val);
        },
        showhide1() {
            this.isShow1 = !this.isShow1;
        },
        showhide2() {
            this.isShow2 = !this.isShow2;
        },
        async getVulninfo(){
            const res = await traffic.trafficVulnInfo({
                flowRiskId: this.vuln_id
            })
            if (res.code == 200) { 
                this.bugdialogVisible = true;
                this.bugbasicinfo = [];
                var bugjosn = {
                    ip:res.data.ip,
                    host:res.data.host,
                    port:res.data.port,
                    risk:res.data.risk,
                    riskLevel:res.data.riskLevel,
                    riskLevelName:res.data.riskLevelName,
                    createTime:res.data.createTime,
                    updateTime:res.data.updateTime,
                    id:res.data.id,
                    status:res.data.status,
                    type:res.data.type,
                    token:res.data.token,
                    origin:res.data.origin, 
                    hash:res.data.hash
                };

                this.bugbasicinfo.push(bugjosn); 
                this.buginfo.url = res.data.url;
                this.buginfo.request = res.data.request;
                this.buginfo.response = res.data.response;
                this.buginfo.parameter = res.data.parameter;
                this.buginfo.detail = res.data.detail;
                this.buginfo.payload = res.data.payload;
                this.buginfo.payload_success_flag = res.data.payload_success_flag;

                this.requestpack  = this.buginfo.request;
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            } 
        },
       
        rowClick(row, event, column) { //在<table>里，我们已经设置row的key值设置为每行数据id：row-key="id"
            Array.prototype.remove = function (val) {
                let index = this.indexOf(val);
                if (index > -1) {
                    this.splice(index, 1);
                }
            };
            if (this.expands.indexOf(row.id) < 0) {
                this.expands.push(row.id);
            } else {
                this.expands.remove(row.id);
            }
            if (this.expands[0]) {
                console.log(this.expands[0]);
            }

        },
        cancalbugdialogVisible() {
            this.isShowInput = false;
            this.bugdialogVisible = false;
            this.is_bugUpdate = false;
            this.responsepack = '';
            this.bugbasicinfo = [];
            // this.buginfo = {};
            
        },
        escapeRegExp(string) {
            return string.replace(/[.*+?^${}()|\[\]\\]/g, '\\$&');
        },
        // HTML 转义函数
        escapeHtml(unsafe) {
            return unsafe 
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;") 
        },
        highlightedRequest(request,payload) {
            // console.log(request,'请求参数00');
            if(!request) return;
            this.requestpack = request 
            
            // 对 payload 也进行 HTML 转义，确保与请求报文格式一致
            let safePayload = this.escapeHtml(payload);
            
            // 1. 尝试直接匹配（转义后的 payload）
            let escapedPayload = this.escapeRegExp(safePayload);
            let replacestr = request.replace(new RegExp(`(${escapedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
            
            // 2. 如果直接匹配失败，尝试 URL 编码后匹配（标准编码，空格为 %20）
            if (replacestr === request) {
                let encodedPayload = encodeURIComponent(payload);
                let escapedEncodedPayload = this.escapeRegExp(encodedPayload);
                replacestr = request.replace(new RegExp(`(${escapedEncodedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
            }
            
            // 3. 如果仍然失败，尝试 URL 编码后匹配（表单编码，空格为 +）
            if (replacestr === request) {
                let encodedPayload = encodeURIComponent(payload).replace(/%20/g, '+');
                let escapedEncodedPayload = this.escapeRegExp(encodedPayload);
                replacestr = request.replace(new RegExp(`(${escapedEncodedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
            }

            //  return request
            return replacestr
        },
         highlightedRequest2(request,payload) {
            if(!request) return;
            // 将包含 target 的部分替换为带有红色标记的内容
            let str = request;
                str=str.replace(/</g, "\&lt");
	        str=str.replace(/>/g, "\&gt");
            
            // 对 payload 也进行 HTML 转义，确保与请求报文格式一致
            let safePayload = this.escapeHtml(payload);
            
            // 1. 尝试直接匹配（转义后的 payload）
            let escapedPayload = this.escapeRegExp(safePayload);
            let replacestr = str.replace(new RegExp(`(${escapedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
            
            // 2. 如果直接匹配失败，尝试 URL 编码后匹配（标准编码，空格为 %20）
            if (replacestr === str) {
                let encodedPayload = encodeURIComponent(payload);
                let escapedEncodedPayload = this.escapeRegExp(encodedPayload);
                replacestr = str.replace(new RegExp(`(${escapedEncodedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
            }
            
            // 3. 如果仍然失败，尝试 URL 编码后匹配（表单编码，空格为 +）
            if (replacestr === str) {
                let encodedPayload = encodeURIComponent(payload).replace(/%20/g, '+');
                let escapedEncodedPayload = this.escapeRegExp(encodedPayload);
                replacestr = str.replace(new RegExp(`(${escapedEncodedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
            }
                
            return replacestr
        },
        async btnTest(){
            // console.log(this.requestpack);
            // console.log(window.btoa(this.requestpack))
            this.isShowInput = false
            let params = { 
                verMsg: Base64.encode(this.requestpack) 
            } 
            const res = await logic.vultest(params)
            if (res.code === 200) { 
                this.$message({
                    message: '测试报文成功',
                    type: 'success'
                }); 
                // window.atob(res.data.respVerMsg)
                //  Base64.decode(res.data.respVerMsg);  
                this.buginfo.response = Base64.decode(res.data.respVerMsg);  

            } else {
                this.dialogloading = false;
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        
    },
}
</script>