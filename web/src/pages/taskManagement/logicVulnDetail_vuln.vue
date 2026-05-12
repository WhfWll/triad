<template>
    <div class="taskVulnlist">
        <div class="search-box">
            <div class="operationbutton"> 
                <el-popover popper-class="delButton_popper" placement="bottom-start" width="170"
                    trigger="click" :visible-arrow="false" v-model="alldelvisible">
                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                    <div style="text-align: right; margin: 0" class="">
                        <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消
                        </el-button>
                        <el-button size="mini" type="primary" @click="btnMultiDelete">确定
                        </el-button>
                    </div>
                    <el-button  type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除
                    </el-button>
                </el-popover>  
            </div>
            <div class="serach-condition">
                <div class="search-text">
                    <el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearch"  v-model="fromData.search" class="input-with-select" size="small"
                        clearable> </el-input>
                    <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                </div>
                <div>
                    <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                </div>
            </div>
        </div>
        <el-table ref="targetTable" :data="tableData" tooltip-effect="dark" style="width: 100%"
        @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave"  
            @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55">
            </el-table-column>
              <el-table-column prop="targetUrl" label="测试目标">  
            </el-table-column> 
            <el-table-column prop="typeName" label="漏洞类型">
            </el-table-column> 
            <el-table-column prop="risk_level" label="漏洞风险"> 
                <template slot-scope="scope">
                    <span
                        :class="[ 
                        { 'riskstyle risk_hight': scope.row.risk == 1 } ,
                        { 'riskstyle risk_middle': scope.row.risk ==2 },
                        { 'riskstyle risk_low': scope.row.risk == 3 },
                        { 'riskstyle risk_nofind': scope.row.risk ==4 }]"><i></i>{{scope.row.riskName}}</span>
                </template>
            </el-table-column>  
             <el-table-column label="漏洞位置" prop="location"> 
            </el-table-column> 
            <el-table-column prop="statusName" label="发现时间" width="200" fixed="right"> 
                <template slot-scope="scope">  
                    <div v-if="showOperateButton && rowId == scope.row.id  ">  
                        <el-link class="link_primary" :underline="false" @click="btnShow(scope.row)">详情</el-link>
                        <el-popover placement="bottom" width="170" :visible-arrow="false"
                            :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                            <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                            <div style="text-align: right; margin: 0">
                                <el-button size="mini" class="delCancel"
                                    @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">
                                    取消</el-button>
                                <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                            </div>
                            <el-link :underline="false" class="link_info linkafter" style="padding:0" slot="reference">
                                删除 </el-link>
                        </el-popover>
                    </div>
                    <div v-else>
                            {{ scope.row.createTime }} 
                    </div>
                </template>

            </el-table-column>
        </el-table>
        <el-pagination :page-size="pageSize" background layout=" total, prev, pager, next, sizes, jumper"
            :total="totalpage" :current-page="currentpage" @current-change="currentchange"
            @size-change="handleSizeChange">
        </el-pagination>


        <el-dialog :title="bugform.name" :visible.sync="bugdialogVisible" width="1184px" class="buginfobox"
        :close-on-click-modal="false" :show-close="false">
        <div class="dialog_b_btn"> 
            <el-button size="small" @click="cancalbugdialogVisible">关闭</el-button>
        </div>
        <div class="buginfo_box" v-loading="dialogloading">
            <div class="bugbasicinfo">
                <el-table :data="bugbasicinfo" size='small' style="width: 100%; ">
                    <el-table-column prop="typeName" label="漏洞类型"> 
                    </el-table-column>
                    <el-table-column prop="riskName" label="漏洞风险"> 
                    </el-table-column> 
                </el-table>
            </div>
           
             <div class="bugotherinfo">
                <div class="part_title">漏洞描述</div>
                <div class="content" > {{ bugform.description }} </div>
                 
            </div>
            <div class="bugotherinfo">
                <div class="part_title">修复建议</div>
                <div class="content" >{{ bugform.fixSuggest }}</div>
                
            </div> 
            <div class="bugotherinfo">
                <div class="part_title">漏洞地址</div>
                <div class="content"  >{{ bugform.vulAddress }}</div>
              
            </div> 
            <div class="bugotherinfo">
                <div class="part_title">漏洞结果</div>
                <div class="content"  >{{ bugform.result }}</div> 
            </div> 
            <div class="bugotherinfo">
                <div class="part_title">攻击载荷</div>
                <div class="content"  >{{ bugform.payload }}</div> 
            </div> 
            <div class="bugotherinfo">
                <div class="part_title">报文</div>
                <div class="btnDiv">   
                    <el-button  size="mini" type="primary" @click="btnTest()" :disabled="!requestpack&&verMsg.length!=1">测试 </el-button> 
                </div>
                <div class="bugbasicinfo"> 
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
                                                <div class="packheight " style="height:100%">
                                                    <div class="packinput" style="height:100%">
                                                        <el-input v-show="isShowInput"  class="packtxt" type="textarea" :rows="7"
                                                            v-model="item.request"
                                                            @change = "highlightedRequest(item.request,bugform.payload)" 
                                                            resize="none"  >
                                                        </el-input> 
                                                        
                                                        <div @click="isShowInput= true" v-show="!isShowInput" style="    background: #fff;">
                                                            <pre v-dompurify-html="highlightedRequest(item.request,bugform.payload) "></pre>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </el-col>
                                    <el-col :span="12">
                                        <div class="message">
                                            <label class="title_bg title_bg2">响应报文</label>
                                            <div style="    background: #fff;"> 
                                                <div style="    background: #fff;">
                                                    <pre v-dompurify-html="highlightedRequest2(item.response,bugform.payload_success_flag) "></pre>
                                                </div>
                                            </div>
                                        </div>
                                    </el-col>
                                </el-row> 
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

    </div>
</template>
<style lang="less" scoped>
.taskVulnlist{
    background: #fff;
    padding: 24px 24px;
    box-sizing: border-box; 
    box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
    border-radius: 4px;
}
.bugbasicinfo {
    padding: 24px;
    background: #fff;
    border: 1px solid #e8e8f5;
}
.buginfo_box {
     padding: 24px;
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
       // margin-bottom: 24px; 
       background: #F7F7FB;
       border-radius: 4px;
       border: 1px solid #E8E8F5;
       padding: 16px;
       box-sizing: border-box;
   }
   .message>div {
     height: 253px;
     overflow-y: auto;
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
</style>
<script>
import { logic } from '@/api/task.js'
import {Base64} from 'js-base64';
export default {
    name:'',
    props:{ 
        task_id:{},
        task_name:{},
    },
    data(){
        return{
            alldelvisible:false,
            tableData:[],
            pageSize:10,
            totalpage:0,
            currentpage:1,
            fromData:{
                search:'',
                page:1,
            },
            multipleSelection:[],
            showOperateButton:false,
            bugform:{
                payload:'',
            },
            bugdialogVisible:false,
            bugbasicinfo:[],
            bugmessage:[],
            verMsg:[],
            isShowInput:false,
            requestpack:'',
            expands:[],
            dialogloading:false,
            payload_success_flag:'',
        }
    },
    created(){

    },
    mounted(){

    },
    methods:{
        async getData(){
            const res = await logic.getVulnlist({
                page:this.fromData.page,
                size:this.pageSize,
                search:this.fromData.search,
                task_id:this.task_id
            });
            if(res.code == 200){
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }
        },
        handleReset(){
            this.fromData.page_num = 1;
            this.fromData.search= ""; 
            this.pageSize = 10;
            this.currentpage = 1;
            this.getData();
        },
        handlesearch(){
            //搜索
            this.fromData.page = 1;
            this.getData();
            this.currentpage = 1;
        }, 
        handleSizeChange(t){
            this.fromData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        currentchange(t){
            this.fromData.page = t;
            this.getData();
            this.currentpage = t;
        },
        async btnMultiDelete(){
            if (this.multipleSelection.length == 0) return; 
              let _ids = this.multipleSelection.map(item => item.id); 
              const res = await logic.vuldel({
                ids: _ids.join(','), 
            });
            if (res.code === 200) {
                this.$message({
                    message: '删除漏洞信息成功',
                    type: 'success'
                });
                this.currentpage = 1;
                this.formData.page = 1;
                this.alldelvisible = false;
                this.getData();
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }     

        },
        async btnDel(scope){
            const res = await logic.vuldel({
                ids: scope.row.id + '', 
            });
            if (res.code === 200) {
                this.$message({
                    message: '删除漏洞信息成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
                this.getData();
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }     
        },
        mouseenter(row, colum, cell, event) { 
            this.showOperateButton = true; 
            this.rowId = row.id ;  //赋值行id，便于页面判断 
        },
        mouseleave(row, colum, cell, event) { 
            if (!this.$refs['popover_id-' + row.id]) {
                this.showOperateButton = false;
                this.rowId = "";
                return;
            } else {
                let isShow = this.$refs['popover_id-' + row.id].showPopper;
                if (!isShow) {
                    this.showOperateButton = false;
                    this.rowId = "";
                }
            }
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        async btnShow(row){ //详情
            const res = await logic.vulnDetail({
                id:row.id
            });
            if(res.code == 200){ 
                this.bugbasicinfo=[];
                this.bugdialogVisible = true;
                this.bugform.id = row.id;
                this.bugform.name = res.data.name;
                let json={
                    typeName:row.typeName,
                    riskName:res.data.riskName,
                    port:'',
                    findtime:'',

                }
                this.bugbasicinfo.push(json);
                this.bugform.description = res.data.description;
                this.bugform.fixSuggest = res.data.fixSuggest;
                this.bugform.vulAddress = res.data.location;
                this.bugform.result = res.data.result;
                this.bugform.payload = res.data.payload;
                this.bugmessage = [{
                    id: res.data.id,
                    verMsg: res.data.verMsg,
                    // respVerMsg: res.data.respVerMsg
                }]
                this.expands = []
                this.expands.push(res.data.id)
                this.verMsg = res.data.verMsg;
                this.requestpack = res.data.verMsg[0].request;
               
                this.bugform.payload_success_flag = res.data.payload_success_flag;

            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        cancalbugdialogVisible(){
            this.bugdialogVisible = false;
            this.dialogloading = false;
            this.isShowInput = false;
        },
        async btnTest(){ 
            this.dialogloading = true;
            this.isShowInput = false; 
            
            // 检查请求报文中是否包含message参数，如果缺失则添加
            let requestToSend = this.requestpack;
            if (requestToSend.includes('message=') && !requestToSend.includes('message=' + encodeURIComponent(this.bugform.payload))) {
                // 如果message参数存在但值为空，替换为完整的payload
                requestToSend = requestToSend.replace(/message=([^&]*)/, 'message=' + encodeURIComponent(this.bugform.payload));
            } else if (!requestToSend.includes('message=')) {
                // 如果message参数不存在，添加它
                if (requestToSend.includes('?')) {
                    requestToSend += '&message=' + encodeURIComponent(this.bugform.payload);
                } else {
                    requestToSend += '?message=' + encodeURIComponent(this.bugform.payload);
                }
            }
            
            let params = {
                // taskVulId: this.bugform.id, 
                verMsg:  Base64.encode(requestToSend) 
            } 
            const res = await logic.vultest(params)
            if (res.code === 200) {
                this.dialogloading = false;
                this.$message({
                    message: '测试报文成功',
                    type: 'success'
                }); 
              
                let str = res.data.respVerMsg;  
                // window.atob(str);  
                this.verMsg[0].response = Base64.decode(str);  
            } else {
                this.dialogloading = false;
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
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
        highlightedRequest(request,payload){ 
            this.requestpack = request;
            if(!request) return;
            let replacestr = request;
            // 将包含 "大3" 的部分替换为带有红色标记的内容
            // 将包含 target 的部分替换为带有红色标记的内容
            if(payload!=""){
                // 1. 尝试直接匹配原始 payload
                let escapedPayload = this.escapeRegExp(payload);
                replacestr = request.replace(new RegExp(`(${escapedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
                
                // 2. 尝试 URL 编码后匹配（标准编码，空格为 %20）
                if (replacestr === request) {
                    let encodedPayload = encodeURIComponent(payload);
                    let escapedEncodedPayload = this.escapeRegExp(encodedPayload);
                    replacestr = request.replace(new RegExp(`(${escapedEncodedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
                }
                
                // 3. 如果标准编码匹配失败，尝试表单编码后匹配（空格为 +）
                if (replacestr === request) {
                    let formEncodedPayload = encodeURIComponent(payload).replace(/%20/g, '+');
                    let escapedFormEncodedPayload = this.escapeRegExp(formEncodedPayload);
                    replacestr = request.replace(new RegExp(`(${escapedFormEncodedPayload})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
                }
                
                // 4. 尝试处理分号编码的情况
                if (replacestr === request) {
                    let encodedPayloadWithSemicolon = encodeURIComponent(payload).replace(/%3B/g, ';');
                    let escapedEncodedPayloadWithSemicolon = this.escapeRegExp(encodedPayloadWithSemicolon);
                    replacestr = request.replace(new RegExp(`(${escapedEncodedPayloadWithSemicolon})`, 'g'), '<span style="color: red;font-weight:700">$1</span>');
                }
            }
            
            return replacestr
        },
        highlightedRequest2(request,payload){
            if(!request) return;
            let str = request;
            str=str.replace(/</g, "&lt;");
            str=str.replace(/>/g, "&gt;");
            
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
    }
}
</script>