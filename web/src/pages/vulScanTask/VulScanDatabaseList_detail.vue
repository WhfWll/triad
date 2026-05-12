<template>
    <el-dialog :title="bugform.title" :visible.sync="bugdialogVisible" width="1184px" class="buginfobox"
        :close-on-click-modal="false" :show-close="false">
        <div class="dialog_b_btn">
            <el-button size="small" @click="cancalbugdialogVisible">关闭</el-button>
        </div>
        <div class="buginfo_box">
            <div class="bugbasicinfo">
                <el-table :data="bugbasicinfo" size='small' style="width: 100%">
                    <el-table-column prop="cve" label="漏洞编号">
                    </el-table-column>

                    <el-table-column prop="severityName" label="漏洞风险">
                        <template slot-scope="scope">
                            <!-- 1-致命/2-高危/3-中危/4-低危/5-信息 -->
                            <span :class="[ 
                                { 'riskstyle risk_hight': scope.row.severity == 'CRITICAL' } ,
                                { 'riskstyle risk_middle': scope.row.severity == 'HIGH' },
                                { 'riskstyle risk_low': scope.row.severity == 'MEDIUM' },
                                { 'riskstyle risk_nofind': scope.row.severity == 'LOW' }]">
                                <i></i>
                                {{scope.row.severityName}}
                                <span v-if="scope.row.severityName">{{ scope.row.base_cvs_sv2_score }}</span>                        </span>
                        </template>
                    </el-table-column>

                    <el-table-column prop="cwe" label="漏洞类型"></el-table-column>
                    
                    <el-table-column prop="vendor" label="厂商"></el-table-column>

                    <el-table-column prop="product" label="产品"></el-table-column>

                    <el-table-column prop="published_date" label="披露时间"></el-table-column>
                    
                </el-table>
            </div>
            
            <div class="bugotherinfo">
                <div class="part_title">漏洞描述</div>
                <div class="content" v-if="!is_bugUpdate"> {{ bugform.description_main_zh }} </div>
                <el-input class="textarea" type="textarea" v-model="updateinfo.description_main_zh" :row="5"
                    v-if="is_bugUpdate"></el-input>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">修复建议</div>
                <div class="content" v-if="!is_bugUpdate">{{ bugform.solution }}</div>
                <el-input class="textarea" type="textarea" v-model="updateinfo.fixSuggest" :row="5"
                    v-if="is_bugUpdate"></el-input>
            </div>
            <div class="bugotherinfo">
                <div class="part_title">参考链接</div>
                <pre class="content">{{ bugform.references }}</pre>
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
import { vulscan } from '@/api/vulscan.js';

export default {
    name:'VulScanDatabaseList_detail',
    components: {
        
    },
    props:{
        value: {}, // 注意此处获取的value对应的就是组件标签中的v-model
        vulninfo:{},
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
                title: 'CVE详情',
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
        showhide1() {
            this.isShow1 = !this.isShow1;
        },
        showhide2() {
            this.isShow2 = !this.isShow2;
        },
        async getVulninfo(){
            let res = {}
            res = await vulscan.cveDetail({
                id: this.vulninfo.id
            })
            
            if (res.code == 200) {  
                this.bugdialogVisible = true;
                this.bugform = res.data
                this.bugform.title = "CVE详情"
                this.bugbasicinfo = [];

                var bugjson = {
                    cwe: res.data.cwe, 
                    severityName: res.data.severityName,
                    severity: res.data.severity,
                    cve: res.data.cve,
                    cvss_version: res.data.cvss_version,
                    exploitImpact: res.data.name,
                    publish_date: res.data.publish_date,
                    vendor: res.data.vendor,
                    product: res.data.product,
                    published_date: res.data.published_date,
                    base_cvs_sv2_score: res.data.base_cvs_sv2_score,
                };
                this.bugbasicinfo.push(bugjson); 
             
                this.bugmessage = [{
                    id: res.data.id,
                    verMsg: res.data.name,
                    respVerMsg: res.data.name
                }]
                this.expands = []
                this.expands.push(res.data.id)
                // this.requestpack = res.data.verMsg;
                // this.responsepack = res.data.respVerMsg;
            //    this.verMsg = res.data.name
      
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