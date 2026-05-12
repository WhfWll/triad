<template>
    <el-dialog :title="bugform.name" :visible.sync="bugdialogVisible" width="1184px" class="buginfobox"
        :close-on-click-modal="false" :show-close="false">
        <div class="dialog_b_btn">
            <!-- <el-button size="small" @click="btnUpdatebug">{{ updateBugtxt }}</el-button> -->
            <!-- <el-button size="small" @click="saveUpdatebug" v-if="is_bugUpdate">保存</el-button> -->
            <el-button size="small" @click="cancalbugdialogVisible">关闭</el-button>
        </div>
        <div class="buginfo_box">
            <div class="bugbasicinfo">
                <el-table :data="bugbasicinfo" size='small' style="width: 100%">
                    <el-table-column prop="type" label="漏洞类型">
                        <template slot-scope="scope">
                            <span>{{ scope.row.typeName }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="riskName" label="漏洞风险">
                        <template slot-scope="scope">
                            <span v-if="!is_bugUpdate">
                                <span :class="[
                                    { 'riskstyle risk_hight': scope.row.risk == 1 },
                                    { 'riskstyle risk_middle': scope.row.risk== 2 },
                                    { 'riskstyle risk_low': scope.row.risk == 3 },
                                    { 'riskstyle risk_nofind': scope.row.risk == 4 }]"><i></i>{{
                                    scope.row.riskName }}</span>
                            </span>
                            <el-select v-model="updateinfo.risk" size="mini" v-if="is_bugUpdate">
                                <el-option v-for="(item, index) in buglevel" :key="index" :label="item.label"
                                    :value="item.value">
                                </el-option>
                            </el-select>
                        </template>
                    </el-table-column>
                    <el-table-column prop="vulNumber" label="漏洞编号">
                    </el-table-column>
                    
                    <el-table-column prop="cvss" label="漏洞评分">
                    </el-table-column>
                    <el-table-column prop="exploitImpact" label="漏洞影响">
                        <template slot-scope="scope">
                            <span v-if="!is_bugUpdate">{{ scope.row.exploitImpact }}</span>
                            <el-select v-model="updateinfo.use_impact_value" size="mini" v-if="is_bugUpdate">
                                <el-option v-for="(item, index) in vulthreatlist" :key="index" :label="item.label"
                                    :value="item.value">
                                </el-option>
                            </el-select>
                        </template>
                    </el-table-column>
                    <el-table-column prop="publishedTime" label="披露时间">
                    </el-table-column>
                    <el-table-column prop="status" label="状态">
                        <template slot-scope="scope">
                            <span :class="[
                            { 'tag_status tag_primary': vulninfo.status == 1 },
                            { 'tag_status tag_warning': vulninfo.status == 2 },
                            { 'tag_status tag_danger': vulninfo.status == 3 },
                            { 'tag_status tag_success': vulninfo.status == 4 }]">{{ vulninfo.statusName }}
                                </span>
                        </template>
                         
                    </el-table-column>
                </el-table>
            </div>
            <!-- <div class="bugotherinfo">
                <div class="part_title">漏洞类型</div>
                <div class="content">{{ bugform.typeName }}</div>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">漏洞风险</div>
                <div class="content">{{ bugform.riskName }}</div>
            </div> -->
            <!-- <div class="bugotherinfo">
                <div class="part_title">漏洞编号</div>
                <div class="content">{{ bugform.vulNumber }}</div>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">漏洞评分</div>
                <div class="content">{{ bugform.cvss }}</div>
            </div> -->
            <!-- <div class="bugotherinfo">
                <div class="part_title">披露时间</div>
                <div class="content">{{ bugform.publishedTime }}</div>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">漏洞影响</div>
                <div class="content">{{ bugform.exploitImpact }}</div>
            </div> -->
             <div class="bugotherinfo">
                <div class="part_title">漏洞描述</div>
                <div class="content" v-if="!is_bugUpdate"> {{ bugform.description }} </div>
                <el-input class="textarea" type="textarea" v-model="updateinfo.description" :row="5"
                    v-if="is_bugUpdate"></el-input>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">修复建议</div>
                <div class="content" v-if="!is_bugUpdate">{{ bugform.fixSuggest }}</div>
                <el-input class="textarea" type="textarea" v-model="updateinfo.fixSuggest" :row="5"
                    v-if="is_bugUpdate"></el-input>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">参考链接</div>
                <div class="content">{{ bugform.refUrl }}</div>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">漏洞位置</div>
                <div class="content" v-if="!is_bugUpdate">{{ bugform.location }}</div>
                <el-input class="textarea" type="textarea" v-model="updateinfo.location" :row="3"
                    v-if="is_bugUpdate"></el-input>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">漏洞结果</div>
                <div class="content">{{ bugform.vulResult }}</div>
            </div>
            
            <div class="bugotherinfo">
                <div class="part_title">漏洞测试报文</div>
                <div class="btnDiv">
                    <!-- <span class="spanBtn spanBtn1" @click="btnValidate" >验证</span>
                    <span class="spanBtn spanBtn2" @click="btnTest" v-if="requestpack">测试</span>
                    <span class="spanBtn spanBtn2 disabledSpan" v-else>测试</span> -->

                    <el-button  size="mini" type="primary" @click="btnValidate()" :disabled="bugform.vulParam==''">验证 </el-button>
                    <!-- <el-button  size="mini" type="primary" @click="btnTest()" :disabled="!requestpack&&verMsg.length!=1">测试 </el-button> -->
                    <el-button  size="mini" type="primary" @click="btnTest()" :disabled="!requestpack">测试 </el-button>
                </div>
                <div class="bugbasicinfo">
                    <!-- <el-table :data="bugmessage" size='small' row-key="id" :expand-row-keys="expands"
                        :show-header="false" @expand-change="handleExpandChange" @row-click="rowClick" -->
                    <el-table :data="bugmessage" size='small' row-key="id" :expand-row-keys="expands" 
                        :show-header="false"
                        style="width: 100%">
                        <el-table-column type="expand">
                            <template slot-scope="scope">
                                
                                <el-row style="margin-top:10px" :gutter="20" v-for="(item,index) in verMsg" :key="index">
                                    <el-col :span="12">
                                        <div class="message requestpack">
                                            <label class="title_bg title_bg1">请求报文</label>
                                            <div>
                                                <!-- v-model="item.request"  :value="highlightedRequest(item.request,item.payload)"-->
                                                <div class="packheight " style="height:100%">
                                                    <div class="packinput" style="height:100%">
                                                        <div @click="isShowInput = true" v-show="!isShowInput" style=" background: #fff;" >
                                                            <pre v-dompurify-html="highlightedRequest(item.request,item.payload) "></pre>
                                                        </div>
                                                        <el-input v-show="isShowInput" class="packtxt" type="textarea" :rows="7"
                                                            v-model="item.request"
                                                            @change = "highlightedRequest(item.request,item.payload)" 
                                                            resize="none"  >
                                                        </el-input>
                                                        <!-- <el-input v-show="!isShowInput" class="packtxt" type="textarea" :rows="7"
                                                            :value = "highlightedRequest(item.request,item.payload)" 
                                                            @focus="isShowInput= true"
                                                            resize="none"  >
                                                        </el-input> -->
                                                        
                                                       
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </el-col>
                                    <el-col :span="12">
                                        <div class="message">
                                            <label class="title_bg title_bg2">响应报文</label>
                                            <div style="    background: #fff;">
                                                <!-- <pre>{{ item.response }}</pre> -->
                                                   <div style="    background: #fff;">
                                                            <pre v-dompurify-html="highlightedRequest2(item.response,item.payload_success_flag) "></pre>
                                                     </div>
                                            </div>
                                        </div>
                                    </el-col>
                                </el-row>
                                <!-- <div class="packbtn">
                                    <el-button style="padding:9px 24px; " size="mini" type="primary"
                                        @click="btnValidate(scope)">发送请求
                                    </el-button>
                                </div> -->
                            </template>
                        </el-table-column>
                        <el-table-column>
                            <template slot-scope="scope">
                                报文{{ scope.$index + 1 }}
                            </template>
                        </el-table-column>

                    </el-table>
                </div>
            </div>

        </div>
    </el-dialog>
</template>
<style lang="less" scoped>



.tag_status {
    height: 20px;
    border-radius: 12px;
    color: #fff;
    display: inline-block;
    width: 65px;
    text-align: center;
    line-height: 20px;
    font-size: 12px;
    margin-right: 4px;
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
 .part_title {
     font-size: 14px;
     margin-bottom: 16px;
     font-weight: 800;
     border-left: 3px solid #4C7AE3;
     padding-left: 10px;
     height: 14px;
     line-height: 14px;
     color: rgba(72, 72, 102, 0.89);
 }
 .bugotherinfo {
     margin-top: 32px;

     .part_title {
         margin-bottom: 8px;
     }

     .content {
         background: rgba(255, 255, 255, 1);
         border-radius: 2px;
         border: 1px solid rgba(232, 232, 245, 1);
         padding: 12px 16px;
         color: rgba(72, 72, 102, 0.64);
         font-size: 13px;
     }
     .btnDiv{
        display: flex;
        margin-top: 12px;
        margin-bottom: 16px;
        .spanBtn{
            width: 50px;
            height: 32px;
            border-radius: 2px;
            border: 1px solid #4C7AE3;
            font-size: 13px;
            box-sizing: border-box;
            display: flex;
            justify-content: center;
            align-items: center;
            cursor: pointer;
        }
        .spanBtn1{
            background-color: #4C7AE3;
            color: #FFF;
            margin-right: 8px;
        }
        .spanBtn2{
            color: #4C7AE3;
        }
        .disabledSpan{
            border: 1px solid #e2e9f3;
            color: #e2e9f3;
            cursor: not-allowed;
        }
     }
     :deep(.expanded){
        display: none;
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
 .packbtn {
     height: 65px;
     box-sizing: border-box;
     padding: 16px;
     text-align: left;
     padding-left: 0;
 }

 .packinput {
     // padding: 16px;
     // height: 185px;
     box-sizing: border-box !important;

     .packtxt {
         height: 100%;

         /deep/ textarea {
             height: 100%;
             border: 0 !important;
         }
     }
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
import {
    task,
} from '@/api/task.js';
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
        task_id: String,
        typeNameWang: String,
    },
    data(){
        return{
            isShowInput:false,
            verMsg:[],
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
            bugmessage: [],
            bugform: {
                typeName: '',
                riskName: '',
                vulNumber: '',
                cvss: '',
                publishedTime: '',
                exploitImpact: '',
                description: '',
                fixSuggest: '',
                refUrl: '',
                vulAddress: '',
                vulResult: '',
                pocname: '',
                vulParam: '',
            },
            // buglevel: [],
            expands: [],  // 要展开的行，数值的元素是row的key值
            responsepack: '',
            requestpack: '',
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
        escapeHtml(unsafe) {
            return unsafe 
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;") 
        },
        escapeRegExp(string) {
            return string.replace(/[.*+?^${}()|\[\]\\]/g, '\\$&');
        },
         highlightedRequest(request,payload) {

            this.requestpack = request
            if(payload == ''){
                return this.escapeHtml(request)
            } 
            
            console.log('原始 payload:', payload);
            console.log('请求报文:', request);
            
            // 1. 尝试直接匹配原始 payload
            let escapedPayload = this.escapeRegExp(payload);
            let replacestr = request.replace(new RegExp(`(${escapedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
            console.log('直接匹配结果:', replacestr !== request ? '成功' : '失败');
            
            // 2. 如果直接匹配失败，尝试 URL 编码后匹配（标准编码，空格为 %20）
            if (replacestr === request) {
                let encodedPayload = encodeURIComponent(payload);
                console.log('标准编码 payload:', encodedPayload);
                let escapedEncodedPayload = this.escapeRegExp(encodedPayload);
                replacestr = request.replace(new RegExp(`(${escapedEncodedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
                console.log('标准编码匹配结果:', replacestr !== request ? '成功' : '失败');
            }
            
            // 3. 如果仍然失败，尝试 URL 编码后匹配（表单编码，空格为 +）
            if (replacestr === request) {
                let encodedPayload = encodeURIComponent(payload).replace(/%20/g, '+');
                console.log('表单编码 payload:', encodedPayload);
                let escapedEncodedPayload = this.escapeRegExp(encodedPayload);
                replacestr = request.replace(new RegExp(`(${escapedEncodedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
                console.log('表单编码匹配结果:', replacestr !== request ? '成功' : '失败');
            }
            
            // 4. 如果仍然失败，尝试处理分号编码的情况
            if (replacestr === request) {
                let encodedPayloadWithSemicolon = encodeURIComponent(payload).replace(/%3B/g, ';');
                console.log('分号编码 payload:', encodedPayloadWithSemicolon);
                let escapedEncodedPayloadWithSemicolon = this.escapeRegExp(encodedPayloadWithSemicolon);
                replacestr = request.replace(new RegExp(`(${escapedEncodedPayloadWithSemicolon})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
                console.log('分号编码匹配结果:', replacestr !== request ? '成功' : '失败');
            }
            
            // 最后对结果进行 HTML 转义，确保安全显示
            return this.escapeHtml(replacestr)
        },
         highlightedRequest2(request,payload) {

            if(payload == ''){
                return this.escapeHtml(request)
            } 
            
            // 1. 尝试直接匹配原始 payload
            let escapedPayload = this.escapeRegExp(payload);
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
            
            // 4. 如果仍然失败，尝试处理分号编码的情况
            if (replacestr === request) {
                let encodedPayloadWithSemicolon = encodeURIComponent(payload).replace(/%3B/g, ';');
                let escapedEncodedPayloadWithSemicolon = this.escapeRegExp(encodedPayloadWithSemicolon);
                replacestr = request.replace(new RegExp(`(${escapedEncodedPayloadWithSemicolon})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
            }
                
            // 最后对结果进行 HTML 转义，确保安全显示
            return this.escapeHtml(replacestr)
        },
        showhide1() {
            this.isShow1 = !this.isShow1;
        },
        showhide2() {
            this.isShow2 = !this.isShow2;
        },
        async getVulninfo(){
                // taskVulId:this.vulninfo.taskVulId || this.vulninfo.vulID|| this.vulninfo.id

            let res = {}
            if(this.typeNameWang == '资产概览'){
                 res = await task.taskVulinfo({
                taskVulId: this.vulninfo.vulID
            })
            }else if(this.typeNameWang == '资产漏洞'){
                 res = await task.taskVulinfo({
                taskVulId: this.vulninfo.taskVulId
            })
            }else{
                 res = await task.taskVulinfo({
                taskVulId: this.vulninfo.id
            })
            }
            if (res.code == 200) {  
                this.bugdialogVisible = true;
                this.bugform = res.data
                this.bugbasicinfo = [];
                var bugjosn = {
                    typeName: res.data.typeName, 
                    riskName: res.data.riskName,
                    risk: res.data.risk,
                    vulNumber: res.data.vulNumber,
                    cvss: res.data.cvss,
                    exploitImpact: res.data.exploitImpact,
                    publishedTime: res.data.publishedTime,
                    status: res.data.status,
                    statusName: res.data.statusName
                };

                this.bugbasicinfo.push(bugjosn); 
                // this.updateinfo.check_vuln_id = res.data.id;
                // this.updateinfo.use_impact_label = res.data.use_impact_label;
                // this.bugform.check_vuln_id = res.data.id;
                // this.bugform.name = res.data.name;
                // this.bugform.pocname = res.data.pocname;
                // this.bugform.priority_description = res.data.priority_description;
                // this.bugform.description = res.data.description;
                // this.bugform.fix_suggest = res.data.fix_suggest;
                // this.bugform.vul_analysis = res.data.vul_analysis;
                // this.bugform.code = res.data.code;
                // this.bugform.vul_location = res.data.vul_location;
                // this.bugform.result = res.data.result;
                // this.bugform.status = res.data.status;
                // this.bugform.vul_num = res.data.vul_num;
                // this.bugform.affect_range = res.data.affect_range;
                this.bugmessage = [{
                    id: res.data.id,
                    verMsg: res.data.verMsg,
                    respVerMsg: res.data.respVerMsg
                }]
                this.expands = []
                this.expands.push(res.data.id)
                // this.requestpack = res.data.verMsg;
                // this.responsepack = res.data.respVerMsg;
               this.verMsg = res.data.verMsg
        //        this.verMsg = [
                
        //         {  //修改
        //     "request":"111111111112222222",
        //     "response":"22222222444444444",
        //     "payload":"2",
        //     "payload_success_flag":"4",
        //  },
        //         {  //修改
        //     "request":"111111111112222222",
        //     "response":"22222222444444444",
        //     "payload":"2",
        //     "payload_success_flag":"4",
        //  }
        //  ],
      
                console.log(111, this.expands)
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            } 
        },
        saveUpdatebug() { //保存编辑漏洞 
            var _j = {
                id: this.updateinfo.check_vuln_id,
                name: this.updateinfo.name,
                vul_risk: this.updateinfo.risk_lever_number,
                description: this.updateinfo.description,
                fix_suggest: this.updateinfo.fix_suggest,
                vul_location: this.updateinfo.vul_location,
                use_impact_value: this.updateinfo.use_impact_value,
            }

            this.$ajax({
                method: 'post',
                url: '/task/vul/update/',
                data: this.qs.stringify(_j)
            }) .then(dt => {
                    let res = dt.data;
                    if (res.success) {
                        this.$message({
                            message: '编辑漏洞信息成功',
                            type: 'success'
                        });
                        this.is_bugUpdate = false;
                        this.bugdialogVisible = false;
                        this.updateBugtxt = '编辑';
                        this.updateinfo.check_vuln_id = '',
                            this.updateinfo.name = '';
                        this.updateinfo.risk_lever = '';
                        this.updateinfo.cve = '';
                        this.updateinfo.cnvd = '';
                        this.updateinfo.cnnvd = '';
                        this.updateinfo.vuln_location = '';
                        this.updateinfo.ref_url = '';
                        this.updateinfo.fix_suggest = '';
                        this.updateinfo.detail = ''; 
                        this.$emit("saveData");
                    } else {
                        this.$message({
                            message: res.error,
                            type: 'error'
                        });
                    }
                })
                .catch(data => {
                });
        }, 
        // 验证
        async btnValidate(scope) {
            let params = {
                taskId: Number(this.task_id),
                taskVulId: this.vulninfo.id,
                // taskVulId: 6,
                pocname: this.bugform.pocname,
                vulParam: this.bugform.vulParam + ''
            }
            const res = await task.vulValidate(params)
            if (res.code === 200) {
                this.$message({
                    message: '漏洞验证成功，目标站点依然存在漏洞',
                    type: 'success'
                });
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        // 测试
        async btnTest(){
            // console.log(2, this.requestpack)
            this.isShowInput = false; 
            
            // 检查请求报文中是否包含message参数，如果缺失则添加
            let requestToSend = this.requestpack;
            if (requestToSend.includes('message=') && !requestToSend.includes('message=' + encodeURIComponent(this.verMsg[0].payload))) {
                // 如果message参数存在但值为空，替换为完整的payload
                requestToSend = requestToSend.replace(/message=([^&]*)/, 'message=' + encodeURIComponent(this.verMsg[0].payload));
            } else if (!requestToSend.includes('message=')) {
                // 如果message参数不存在，添加它
                if (requestToSend.includes('?')) {
                    requestToSend += '&message=' + encodeURIComponent(this.verMsg[0].payload);
                } else {
                    requestToSend += '?message=' + encodeURIComponent(this.verMsg[0].payload);
                }
            }
            
            let newpack= Base64.encode(requestToSend);
            console.log(newpack)
            // console.log(window.btoa(newpack))
            // return;
            let params = {
                taskVulId: this.vulninfo.id,
                // verMsg: window.btoa((modifiedString).substring(1, originalString.length - 1))

                // taskVulId: 6,
                // verMsg: JSON.stringify(this.requestpack)
                verMsg:newpack
                // verMsg: window.btoa(JSON.stringify(this.requestpack))
            }
            const res = await task.vulTest(params)
            if (res.code === 200) {
                this.$message({
                    message: '测试报文成功',
                    type: 'success'
                });
                // this.responsepack = window.atob(res.data.respVerMsg)
                this.verMsg[0].response =  Base64.decode(res.data.respVerMsg) 
            } else {
                this.$message({
                    message: "报文发送失败，目标站点可能不存活",
                    type: 'error'
                });
            }
        },
        rowClick(row, event, column) { //在<table>里，我们 已经设置row的key值设置为每行数据id：row-key="id"
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
            this.isShowInput=false;
            this.bugdialogVisible = false;
            this.is_bugUpdate = false;
            this.responsepack = '';
            this.updateBugtxt = '编辑';
            this.expands = [];
        },
        btnUpdatebug() { //漏洞编辑
            this.is_bugUpdate = true;
            this.updateBugtxt = '编辑中';
            this.updateinfo.id = this.bugform.id,
                this.updateinfo.name = this.bugform.name;
            this.updateinfo.risk_lever = this.bugbasicinfo[0].risk_lever;
            this.updateinfo.risk_lever_number = this.bugbasicinfo[0].risk_lever_number;
            this.updateinfo.use_impact_label = this.bugbasicinfo[0].use_impact_label;
            this.updateinfo.use_impact_value = this.bugbasicinfo[0].use_impact_value;
            this.updateinfo.vul_location = this.bugform.vul_location;
            this.updateinfo.fix_suggest = this.bugform.fix_suggest;
            this.updateinfo.description = this.bugform.description;

        },
        
        
    },
    computed: {
    sanitizedData() {
      return this.sanitizeHTML(this.rawData);
    }
  },
}
</script>