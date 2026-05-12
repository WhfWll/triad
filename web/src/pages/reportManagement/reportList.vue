<template>
    <div>
        <div class="main-title  ">  
			
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label class="taskSceneBtn" >报告清单</label>
	  	</div> 
          <div class="reportlist context_box_bg">
            <div class="search-box"> 
				<div  class="operationbutton" >  
					<el-popover
                        popper-class="delButton_popper"
                        placement="bottom-start"
                        width="170"
                        style="padding-right:8px"
                        trigger="click" 
                        :visible-arrow="false"
                        v-model="alldelvisible" >
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="" >
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false" >取消</el-button>
                            <el-button size="mini" type="primary"  @click="AllDel">确定</el-button>
                        </div>  
                        <el-button type="warning"  size="small"  slot="reference" :disabled="!multipleSelection.length">删除</el-button> 
                    </el-popover>   
                    <!-- <el-button type="primary" size="small" @click="downword">下载任务word</el-button>
                    <el-button type="primary" size="small" @click="downwordtarget">下载目标word</el-button>
                    <el-button type="primary" size="small" @click="downhtml">下载任务HTMLzip</el-button>
                    <el-button type="primary" size="small" @click="downloadhtmlTarget">下载目标HTMLzip</el-button> -->
				</div>
				<div class="serach-condition" > 
					<div class="search-text">
						<el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch"  v-model="formData.search" class="input-with-select"  size="small" clearable > </el-input>
						<el-button type="primary"  size="small" @click="handlesearch">搜索</el-button> 
					</div>
					<div >
						<el-button type="primary"  size="small" @click="handleReset">重置</el-button> 
					</div>   
				</div>  
			</div>
            <el-table
                ref="multipleTable"
                :data="tableData" 
                tooltip-effect="dark" 
                style="width: 100%" height="calc(100% - 102px)"
                @selection-change="handleSelectionChange" @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column
                    type="selection"
                    width="55"   > 
                    </el-table-column> 
                <el-table-column
                    prop="name"
                    label="报告标题"  :show-overflow-tooltip="true">  
                    </el-table-column>
                <el-table-column
                    prop="typeName"
                    label="报告类型"  > 
                </el-table-column>
                <el-table-column
                    prop="formatName"
                    label="报告格式"  > 
                </el-table-column>
                <el-table-column
                    prop="createTime"
                    label="生成时间"  > 
                </el-table-column> 
                <el-table-column
                    prop="status"
                    label="状态" >
                        <template slot-scope="scope"> 
                            <div v-if="showEditFileNameButton && rowId == scope.row.id">
                                <el-link :underline="false" class="link_primary" v-show="true"  v-if="scope.row.status==3" @click="btnPreview(scope.row)">预览</el-link>
                                <el-link :underline="false" class="link_primary" @click="btnDownLoad(scope.row)" v-if="scope.row.status==3">下载</el-link>
                                <el-popover
                                    placement="bottom"
                                    width="170"   
                                    :visible-arrow="false"
                                    :ref="`popover_id-${scope.row.id}`"
                                    popper-class="delButton_popper" >
                                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                    <div style="text-align: right; margin: 0">
                                        <el-button size="mini" class="delCancel"  @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                        <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                    </div> 
                                    <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference" >删除</el-link>  
                                </el-popover> 
                            </div>
                            <div v-else > 
                                <span :class="[ 
                                    { 'tag_status tag_danger1': scope.row.status == 1 } ,
                                    { 'tag_status tag_primary': scope.row.status ==2 },
                                    { 'tag_status tag_success': scope.row.status == 3 } ]">{{scope.row.statusName}}</span>
                            </div>
                        </template>
                </el-table-column>    
            </el-table>
            <el-pagination 
                background
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                :current-page="currentPage" 
                :page-size="pageSize"
                layout=" total,  prev, pager, next, sizes,jumper"
                :total="total">
            </el-pagination>
            <div id="targetschartpie"></div>
        </div>
    </div>
</template>
<style scoped lang="less">
.reportlist{

    padding: 24px; 
	background: #fff; 
	height: calc(100% - 39px);
    box-sizing: border-box;
	box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
}
#targetschartpie{
    width: 200px;
    height: 200px;  
    opacity: 0;
}
</style>
<script>
import {report} from '@/api/report.js' 
import Docxtemplater from 'docxtemplater';
import JSZip from 'jszip';
import JSZip2 from 'jszip2';
import JSZipUtils from 'jszip-utils';
// import { saveAs } from 'file-saver';
var echarts = require('echarts');  
import FileSaver from "file-saver";
import {report_css} from './utils/report.js'
import {bootstrap_css} from './utils/bootstrap.js'
import {echarts_min} from './utils/echarts_min.js'
import {iconfont} from './utils/iconfont.js' 
// import {htmlTemplate} from './utils/htmlTaskTemp.html'
const expressions = require("angular-expressions");
var assign = require("lodash/assign");
import { export2Excel } from '@/utils/excelexport'; 

import ExcelJS from 'exceljs';
export default {
    name:'reportlist',
    data(){
        return{
            reportName:'',
            alldelvisible:false,
            total:0,
            currentPage:1,
            pageSize:10,
            showEditFileNameButton:false,
            rowId:'',
            multipleSelection:[],
            formData:{
                page:1,
                search:'',
            },
            tableData:[],
            targetdata:{},
            imgUrl:'',
            reportcontent:{ },
            hyAllData:[],
            reportcontent1:{},
            chartData:[], 
            titleNameArr:['第一季度持仓','第二季度持仓'],
            list:[[['代码','持仓','收益率','收益'],[1,2,3,4],[1,2,3,4]],[['代码','持仓','收益率','收益'],[5,6,5,6],[5,6,5,6]]],
            excelTitleName:'养鸡报告', 
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/reportlist';
    },
    mounted(){
        this.getData();
    },
    methods:{
        async getData(){
            const res = await report.reportList({
                page:this.formData.page,
                size:this.pageSize,
                search:this.formData.search
            })
            if(res.code == 200){
                this.tableData = res.data.list; 
                this.total = res.data.total;

            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                }); 
            }
        },
        handleCurrentChange(t){
            this.formData.page = t; 
            this.getData();
            this.currentPage = t;
        },
        handleSizeChange(t){
            this.formData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handlesearch(){
            this.formData.page = 1;
            this.getData();
            this.currentpage = 1;
        },
        handleReset(){
            this.formData.search = ''; 
			this.formData.page = 1;
            this.pageSize=10;
            this.currentpage = 1;
            this.getData();
        },
        handleSelectionChange(val){
            this.multipleSelection = val
        },
        async AllDel(){
            if(this.multipleSelection.length == 0) return;
    		let _ids = [];
    		for (var i = 0; i < this.multipleSelection.length; i++) {
    			_ids.push(this.multipleSelection[i].id);
			} 
            const res = await report.reportDel({
                reportIds:_ids.join(',')
            });
            if(res.code == 200){ 
                this.$message({
                    message:'删除报告成功',
                    type: 'success'
                });
                this.alldelvisible = false;
                this.getData();
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        async btnDel(scope){
            const res = await report.reportDel({
                reportIds:scope.row.id
            });
            if(res.code == 200){
                this.deldialogVisible = false;
                this.$message({
                    message:'删除报告成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                this.getData();
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
			
            let t = this.$refs['popover_id-' + row.id].showPopper;
            if(!t){
                  this.showEditFileNameButton = false;
            this.rowId = "";
            }
          
        }, 
        btnPreview(row){ //预留
            if(row.type ==1 || row.type ==3 || row.type ==5){ //任务 

                const routeData = this.$router.resolve({
                    path: '/taskreportview',
                    query: { id: row.id, }
                })
                window.open(routeData.href, '_blank')
 
               
            }
            if(row.type ==2 || row.type ==4 || row.type ==6){//目标  
                const routeData = this.$router.resolve({
                    path: '/targetreportview',
                    query: { id: row.id, }
                })
                window.open(routeData.href, '_blank')
            }
        },
        async btnDownLoad(row){ //下载
             console.log(row,'row');
            const res = await report.downLoadfile({
                reportId:row.id
            });
            // 临时接口
            // const res = await report.downLoadreporttemp({
            //     reportId:row.id
            // })
            if(res.code == 200){
                debugger
                var configJson = JSON.parse(res.data.configJson);
                if (configJson.cover) {
                    console.log(configJson.cover);
                    const regex = /background:\s*url\((data:image\/png;base64,[^\)]+)\)/;
                    var report_css_new = report_css.replace(regex, `background: url(${configJson.cover})`);
                }    
                else{
                    var report_css_new = report_css;
                }

                this.reportcontent = JSON.parse(res.data.content);
                this.reportName = res.data.name;
//                 this.hyAllData = [
//                     {
//     "reportId":136,
//     "reportCover":{
//         "title":"目标测试报告1111",
//         "createTime":"2023-09-22",
//         "backgroundImg":""
//     },
//     "catalogParent":[
//         {
//             "name":"1. 报告概要",
//             "id":"targetOverview",
//             "isShow":true,
//             "catalog":null
//         },
//         {
//             "name":"2. 资产信息",
//             "id":"assetInfo",
//             "isShow":true,
//             "catalog":null
//         },
//         {
//             "name":"3. 漏洞信息",
//             "id":"vulInfo",
//             "isShow":true,
//             "catalog":null
//         }
//     ],
//     "targetOverview":{
//         "target":"cheilmc.co.kr",
//         "risk":"低危",
//         "vulnStat":{
//             "total":16,
//             "deadlyNumber":0,
//             "highNumber":0,
//             "middleNumber":0,
//             "lowNumber":16
//         },
//         "vulnVerify":{
//             "verifySuccess":16,
//             "useSuccess":0,
//             "repairSuccess":0
//         },
//         "createDate":"2023-09-20至2023-09-20"
//     },
//     "assetInfo":{
//         "component":"http、https",
//         "service":"80/http、443/https",
//         "ipOrUrl":"cheilmc.co.kr",
//         "system":""
//     },
//     "vulInfo":[
//         {
//             "vulName":"将'Permissions-Policy header not implemented'翻译成中文为：\"'Permissions-Policy'头部尚未实现\"。",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"“Permissions-Policy”头部允许开发者有选择地启用和禁用各种浏览器功能和API。",
//             "res":"{\"pocname\":\"bd2ff718-596b-a1cd-3b08-6e61d4b3d8da\"}",
//             "fix":"近期，关于“Permissions-Policy”头部存在的一个漏洞引发了广泛的关注。这个漏洞使得开发者可以任意地启用和禁用各种浏览器功能和API，从而可能对用户的隐私和安全造成潜在风险。\n\n然而，我们可以采取一些措施来修复这个漏洞，以保护用户的利益和数据安全。首先，我们需要强调在使用“Permissions-Policy”头部时谨慎选择和启用功能和API。开发者应该仔细考虑每个功能的必要性和潜在风险，并只启用必要的功能。\n\n其次，我们可以通过更新和升级浏览器的版本来修复这个漏洞。浏览器厂商可以在新版本中修复该漏洞，并确保用户使用的浏览器具有最新的安全补丁和修复程序。因此，建议用户时刻关注并更新他们所使用的浏览器，以保持其安全性。\n\n另外，在网站和应用程序的开发过程中，开发者可以采取一些额外的安全措施来强化对漏洞的防范。这包括但不限于：对用户输入数据进行有效的验证和过滤、采用安全的编程实践、使用最新的安全框架和加密技术等。这些措施可以帮助开发者及早发现和修复潜在的漏洞。\n\n最后，对于已经被利用的漏洞，开发者应该尽快采取措施来修复并通知用户。及时披露漏洞和提供解决方案是至关重要的，以确保用户的信息和安全不受到进一步的威胁。\n\n综上所述，虽然“Permissions-Policy”头部存在一个漏洞，但我们可以通过谨慎选择和启用功能、升级浏览器、加强开发过程中的安全措施和及时修复已知漏洞来修复该漏洞。通过这些努力，可以更好地保护用户的隐私和数据安全。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"没有设置HttpOnly标志的Cookies",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"一个或多个Cookie没有设置HttpOnly标志。当使用HttpOnly标志设置Cookie时，它指示浏览器该Cookie只能由服务器访问，而不能由客户端脚本访问。对于会话Cookie来说，这是一种重要的安全保护措施。",
//             "res":"{\"pocname\":\"04b4867e-5297-23ff-c810-97acf15ed71a\"}",
//             "fix":"一个或多个Cookie没有设置HttpOnly标志会导致安全漏洞的出现，攻击者可以利用此漏洞来窃取用户的会话信息或伪造会话。因此，设置HttpOnly标志是一种重要的安全保护措施，可以提高应用程序的安全性。\n\n修复该漏洞的方法是为每个Cookie设置HttpOnly标志。在应用程序中，可以通过编程方式设置HttpOnly标志，确保在客户端脚本中无法访问和修改Cookie。在使用第三方库或框架时，要确保它们正确地实现了HttpOnly标志，并且该标志被正确地设置。\n\n除了设置HttpOnly标志之外，应用程序还应该使用HTTPS协议来保护会话信息的传输，以避免会话劫持或信息泄露。在应用程序中，还可以使用其他安全措施来保护会话和敏感信息，如使用CSRF令牌、检测异常登录等等。\n\n通过正确地设置HttpOnly标志和其他安全措施，可以最大程度地提高应用程序的安全性，保护用户的会话信息和个人数据，避免被黑客攻击。因此，开发者需要关注和修复此漏洞，确保应用程序的安全性和可靠性。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"可能找到虚拟主机",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"虚拟主机是一种在单个服务器（或服务器池）上托管多个域名（各自处理每个名称）的方法。这允许一个服务器共享其资源，例如内存和处理器周期，而无需所有提供的服务使用相同的主机名。当 Host 标头被操纵并测试各种常见的虚拟主机时，此 Web 服务器会做出不同的响应。这可能表明存在虚拟主机。",
//             "res":"{\"pocname\":\"aa32ac88-2c39-e6ac-07d8-3f02ad5c1b53\"}",
//             "fix":"为了修复上述描述的虚拟主机漏洞，我们可以采取以下措施。首先，服务器管理员应确保实施了正确的虚拟主机配置。这包括仔细设置各个虚拟主机的主机名和域名，并确保每个虚拟主机的配置是独立的，不会相互影响。\n\n其次，服务器应该对 Host 标头进行验证和过滤，以防止恶意操纵。可以通过限制 Host 标头的内容，仅允许合法的主机名和域名，来减少攻击面。此外，过滤掉特殊字符和符号，以防止利用输入验证漏洞。\n\n另外，服务器管理员还可以使用安全的虚拟主机软件和工具来帮助防止漏洞的发生。这些工具可以提供额外的安全层，例如入侵检测系统（IDS）和防火墙，在虚拟主机层面上检测和阻止任何恶意的请求。\n\n定期更新和升级服务器软件和操作系统也是非常重要的。漏洞通常与旧版本的软件相关，因此确保服务器上安装的软件和操作系统始终是最新的，以减少攻击者利用已知漏洞的可能性。\n\n此外，服务器管理员还应该定期进行安全审计和漏洞扫描，以发现任何可能的漏洞并及时修复。定期审查服务器配置和日志也有助于及早发现任何异常活动。\n\n最后，持续的员工培训和意识提高也是至关重要的。在员工中普及安全意识和最佳实践，教育他们如何避免潜在的安全风险，可以大大减少漏洞的发生。\n\n通过以上措施的综合应用，可以帮助修复和预防虚拟主机漏洞的发生，保护服务器和托管的域名免受潜在的攻击和侵犯。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"TLS/SSL证书即将过期。",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"您的服务器使用的一个TLS/SSL证书即将过期。\u003cbr/\u003e\u003cbr/\u003e一旦证书过期，大多数网络浏览器将向终端用户显示安全警告，要求他们手动确认您证书链的真实性。软件或自动系统可能会静默地拒绝连接到服务器。\u003cbr/\u003e\u003cbr/\u003e此警报不一定是由服务器（叶子）证书引起的，但可能是由中间证书触发的。请参考警报详细信息中的证书序列号来确定受影响的证书。",
//             "res":"{\"pocname\":\"029afcbb-3ec2-be3c-5cc3-29f5cfe016f4\"}",
//             "fix":"此漏洞需要立即采取措施进行修复，以避免可能的安全问题。修复措施包括更新证书，重新生成证书链，或更新中间证书。以下是一些可以采取的修复措施：\n\n1. 更新证书：从证书颁发机构（CA）获取最新的证书，并使用此证书替换受影响的证书。这将确保您的证书链不会过期，从而避免出现安全警告。\n\n2. 重新生成证书链：重新生成整个证书链，以确保其中的所有证书都是最新的，并且颁发者都是可信的。使用此修复措施的最终结果将是一个时间更长且更安全的证书链。\n\n3. 更新中间证书：如果中间证书引起了问题，您可以从CA处获取最新的中间证书，并使用此证书更新受影响的证书。这将修复中间证书问题，并确保您的证书链不会过期。\n\n总之，修复这个漏洞是非常重要的。在实施任何修复措施之前，请确保备份所有受影响的证书和证书链。如果您不确定如何执行修复，请咨询您的IT安全专家或者CA颁发机构，以获取正确的指导。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"未实现HTTP严格传输安全(HSTS)",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"\"HTTP严格传输安全性（HSTS）告知浏览器一个网站只能通过HTTPS方式访问。\n        \n检测到您的网络应用程序未实施HTTP严格传输安全性（HSTS），因为响应中缺少Strict Transport Security头部信息。\"",
//             "res":"{\"pocname\":\"34a6c791-c497-27d5-7272-6a968e9fdccb\"}",
//             "fix":"HTTP严格传输安全性（HSTS）被用来告知浏览器一个网站只能通过HTTPS方式访问，以帮助保护用户的敏感数据。如果您的网络应用程序未正确实施HSTS，攻击者可能利用该情况通过中间人攻击等方式将用户数据劫持或篡改。因此，及时修复HSTS漏洞至关重要。\n\n修复HSTS的方法主要包括两种。第一种，通过为响应添加Strict Transport Security头部信息来实现。在服务器设置中添加Strict-Transport-Security: max-age=31536000; includeSubDomains头部信息，其中max-age参数用于设置浏览器缓存HSTS策略的时间，而includeSubDomains参数则可以将HSTS策略应用于所有子域名。第二种方法是使用网站安全证书，将网站升级到HTTPS协议。在此过程中，服务器会自动将HSTS头部信息添加到所有信任的HTTPS通信中。\n\n无论哪种方法，修复HSTS漏洞都需要涉及到服务器端的配置和修改。因此，在实施修复措施之前，一定要进行充分的测试并备份相关数据。此外，及时禁用或删除不再使用的HSTS策略也是必要的，以确保服务器的安全性。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"未实施内容安全策略 (CSP)",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"'内容安全策略（Content Security Policy，CSP）是一种增加的安全层，有助于检测和缓解某些类型的攻击，包括跨站脚本（XSS）和数据注入攻击。\u003cbr/\u003e\u003cbr/\u003e\n\n可以通过添加 \u003cstrong\u003eContent-Security-Policy\u003c/strong\u003e 头来实现内容安全策略（CSP）。该头的值是一个包含策略指令的字符串，用于描述您的内容安全策略。要实现CSP，您应该为站点使用的所有资源类型定义允许的来源列表。例如，如果您有一个简单的站点需要本地加载脚本、样式表和图片以及来自jQuery库的CDN，CSP头应该如下所示：\n\n\u003cpre\u003e\u003ccode\u003e\nContent-Security-Policy:\n    default-src 'self';\n    script-src 'self' https://code.jquery.com;\n\u003c/code\u003e\u003c/pre\u003e\n\n\u003cbr/\u003e\u003cbr/\u003e\n\n检测到您的Web应用程序没有实现内容安全策略（CSP），因为响应中缺少CSP头。建议将内容安全策略（CSP）实现到您的Web应用程序中。'",
//             "res":"{\"pocname\":\"011055fc-94f1-ab96-56ac-53117c56fb4d\"}",
//             "fix":"该漏洞描述了网站安全中的内容安全策略（Content Security Policy，CSP）实现问题。为了增加网站的安全性，建议在网站中实现CSP。具体做法是添加Content-Security-Policy头，并为站点中使用的所有资源类型定义允许的来源列表。例如，对于一个简单的站点，应该定义本地加载脚本、样式表和图片，以及来自jQuery库的CDN允许的来源。这些来源被写入CSP头，作为策略指令的字符串。这样可以检测和缓解一些类型的攻击，包括跨站脚本（XSS）和数据注入攻击。如果您的Web应用程序没有实现内容安全策略（CSP），建议您立即实现。这样可以确保您的网站更加安全，令攻击者难以突破网站安全防御。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"\"无HTTP重定向\"",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"检测到您的Web应用程序使用HTTP协议，但不会自动重定向用户到HTTPS。",
//             "res":"{\"pocname\":\"3f6a8a0e-07f2-af81-54ff-61020299caeb\"}",
//             "fix":"该漏洞描述了一个常见的web应用程序安全问题，即应用程序使用HTTP协议而不是HTTPS协议。HTTP协议是一种不安全的传输协议，未加密的数据包会通过网络传递。这使得攻击者可以轻易地读取、修改和窃取这些数据，并对用户隐私造成威胁。因此，为了确保应用程序使用最高级别的安全措施，就需要使用HTTPS协议代替HTTP协议。HTTPS协议使用传输层安全（TLS）或安全套接字层（SSL）协议来加密数据传输。这样，只有受信任的接收方可以读取和使用传输数据。\n\n为了修复这个漏洞，应用程序开发人员需要将HTTP协议更改为HTTPS协议。他们可以在应用程序配置文件中为使用HTTPS协议的路径添加重定向规则，这将确保用户始终通过HTTPS协议访问应用程序。此外，开发人员应可以使用SSL证书或其他安全措施来保护传输过程中的数据。这可以通过将SSL证书绑定到应用程序域名上来实现，或使用基于云的安全服务来实现。这些安全措施将确保用户隐私的安全，并使攻击者无法查看或窃取数据。在应用程序修复这个漏洞后，网站管理员应该定期进行安全测试和漏洞扫描，以确保应用程序继续保持安全并且满足最高级别的安全标准。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"缺少、不一致或相互矛盾属性的Cookie。",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"“以下至少一个cookie属性会导致该cookie无效或与同一cookie的其他属性或使用环境不兼容。虽然这本身不是漏洞，但它很可能会导致应用程序的意外行为，进而可能引起二次安全问题。”",
//             "res":"{\"pocname\":\"fd1ec03e-319d-ae3b-425f-d8c1275619d8\"}",
//             "fix":"该漏洞描述了一种可能会导致应用程序出现意外行为的情况，但其本身并不是漏洞。 然而，由于该情况可能引起二次安全问题，因此开发人员和管理员需要注意这个问题，并在必要时采取修复措施。\n\n为了避免出现不兼容的Cookie属性导致的意外行为和二次安全问题，开发人员应该仔细检查应用中使用的Cookie属性，并确保它们与应用程序使用的其他部分兼容。\n\n一种解决方法是使用Cookie安全标志，这将强制浏览器只在安全连接上发送Cookie。另外，可以采用限制cookie范围的方法，在特定的域名或路径下使用Cookie，以此来限制Cookie的使用范围，从而减少意外行为和二次安全问题的发生。\n\n最重要的是，开发人员应该经常进行漏洞扫描和安全测试，以及对代码进行定期的审查，从而及时发现和修复可能存在的问题。这可以确保应用程序始终保持安全，并减少因错误配置和使用Cookie属性导致的潜在风险。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"未设置安全标志的Cookie",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"一个或多个cookie没有设置安全标志。当一个cookie被设置了安全标志后，它会指示浏览器该cookie只能在安全的SSL/TLS通道上访问。这是会话cookie重要的安全保护措施。",
//             "res":"{\"pocname\":\"575ba039-89b6-b7d7-eba6-553135c04553\"}",
//             "fix":"该漏洞的修复措施旨在确保会话cookie的安全保护措施得到有效实施。对于未设置安全标志的cookie，开发人员应立即对其进行修改，将其设置上安全标志标识。通过这种方法，开发人员可以消除会话cookie的任何潜在风险，防止恶意用户在传输过程中劫持用户身份，从而造成安全威胁。\n\n为确保实现安全标志标识的cookie准确性，开发人员还需要通过SSL/TLS访问控制来限制访问网络资源的用户身份。通过实施包括内容保护和访问控制的安全策略，开发人员可以确保会话cookie的安全性得到最大程度的保障。在进行此类操作的过程中，开发人员应仔细跟踪会话cookie的使用情况，及时发现和修复任何与设置安全标志标识相关的错误或漏洞。\n\n总之，对会话cookie设置安全标志标识是防范安全风险必不可少的一项措施。通过限制cookie访问范围，加强访问控制和内容保护等操作，开发人员可以为用户提供更加安全可靠的网络环境，保护客户的私人和机密信息。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"点击劫持：X-Frame-Options头部",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"点击劫持（用户界面篡改攻击，UI篡改攻击，UI篡改）是一种恶意技术，可以欺骗Web用户点击与其所认为点击不同的内容，从而可能泄露机密信息或在点击看似无害的网页时控制他们的计算机。服务器未返回具有DENY或SAMEORIGIN值的X-Frame-Options头，这意味着该网站可能面临点击劫持攻击的风险。 X-Frame-Options HTTP响应头可用于指示是否应允许浏览器在框架或iframe中呈现页面。站点可以使用此功能来避免点击劫持攻击，确保其内容不嵌入不受信任的站点。",
//             "res":"{\"pocname\":\"b8e2c082-44f1-cf0b-0b8e-0e0bb357e798\"}",
//             "fix":"点击劫持是一种恶意技术，可以控制用户的计算机或者泄漏机密信息。这种攻击方式会欺骗用户点击看似无害的网页，但是实际上它已经被篡改了。这种攻击技术对网站安全造成了严重威胁。为了防止点击劫持攻击，我们可以使用X-Frame-Options HTTP响应头来避免嵌入不信任的站点。如果服务器没有返回具有DENY或SAMEORIGIN值的X-Frame-Options头，就存在被点击劫持的风险。\n\nX-Frame-Options设置框架的漏洞，是通过消除非授权嵌入的方式来解决的。它的主要作用是避免网站被嵌到不受信任的站点。为了设置这个功能，我们可以在服务器上添加以下代码：\n\nX-Frame-Options: DENY\n\n或\n\nX-Frame-Options: SAMEORIGIN\n\n这样就可以避免网站面临点击劫持攻击的风险了。如果站点未修改这个设置，那么可能会面临许多安全问题。因此，我们应该要了解这个漏洞，并采取相应的措施来保护我们的网站和用户不受攻击威胁。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"Unencrypted connection的中文翻译为“非加密连接”。",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"这个扫描目标连接是未加密的。一个潜在的攻击者可以拦截和修改从这个站点发送和接收的数据。",
//             "res":"{\"pocname\":\"bca221d1-8581-3375-4097-66b0048ed088\"}",
//             "fix":"该漏洞是因为扫描目标连接未加密，导致攻击者可以拦截和修改该站点发送和接收的数据，因此需要采取一定的措施进行修复。通常，可以通过SSL（安全套接字层）协议来解决这个问题，SSL协议是一种加密协议，可以保护网络传输过程中的数据安全。\n\n对于站点管理者来说，他们需要考虑启用SSL证书以支持HTTPS协议，HTTPS协议相比HTTP协议是基于SSL加密的，在传输过程中可以将数据进行加密传输，确保数据的安全性。并且站点管理者可以通过定期的安全检查以及更新站点的安全软件，及时发现和修复可能存在的漏洞，保证站点的安全性。\n\n对于用户来说，他们需要注意防范网络钓鱼以及恶意软件的攻击，尽可能选择使用SSL加密的站点，同时多使用杀毒软件和防火墙等安全软件，增强自身网络安全意识。在使用公共网络时，为了避免个人账号和密码等敏感信息的泄露，建议不要随意使用未知的公共网络，并且尽可能不在公共网络中登录个人账号。\n\n综上所述，修复该漏洞的关键在于采取措施来加强站点的安全保障体系，包括启用SSL证书、定期更新安全软件、及时发现和修复可能存在的漏洞等。同时，用户也需要更加注重自身的网络安全，选择SSL加密的站点、使用安全软件，避免在公共网络中登录敏感信息，从而保障网络的安全。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         }
//     ]
// },


// {
//     "reportId":13886,
//     "reportCover":{
//         "title":"目标测试报告2222222",
//         "createTime":"2023-09-22",
//         "backgroundImg":""
//     },
//     "catalogParent":[
//         {
//             "name":"1. 报告概要88",
//             "id":"targetOverview",
//             "isShow":true,
//             "catalog":null
//         },
//         {
//             "name":"2. 资产信息88",
//             "id":"assetInfo",
//             "isShow":true,
//             "catalog":null
//         },
//         {
//             "name":"3. 漏洞信息88",
//             "id":"vulInfo",
//             "isShow":true,
//             "catalog":null
//         }
//     ],
//     "targetOverview":{
//         "target":"222222222cheilmc.co.kr",
//         "risk":"低危",
//         "vulnStat":{
//             "total":16,
//             "deadlyNumber":0,
//             "highNumber":0,
//             "middleNumber":0,
//             "lowNumber":16
//         },
//         "vulnVerify":{
//             "verifySuccess":16,
//             "useSuccess":0,
//             "repairSuccess":0
//         },
//         "createDate":"2023-09-20至2023-09-20"
//     },
//     "assetInfo":{
//         "component":"http、https",
//         "service":"80/http、443/https",
//         "ipOrUrl":"cheilmc.co.kr",
//         "system":""
//     },
//     "vulInfo":[
//         {
//             "vulName":"将'Permissions-Policy header not implemented'翻译成中文为：\"'Permissions-Policy'头部尚未实现\"。",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"“Permissions-Policy”头部允许开发者有选择地启用和禁用各种浏览器功能和API。",
//             "res":"{\"pocname\":\"bd2ff718-596b-a1cd-3b08-6e61d4b3d8da\"}",
//             "fix":"近期，关于“Permissions-Policy”头部存在的一个漏洞引发了广泛的关注。这个漏洞使得开发者可以任意地启用和禁用各种浏览器功能和API，从而可能对用户的隐私和安全造成潜在风险。\n\n然而，我们可以采取一些措施来修复这个漏洞，以保护用户的利益和数据安全。首先，我们需要强调在使用“Permissions-Policy”头部时谨慎选择和启用功能和API。开发者应该仔细考虑每个功能的必要性和潜在风险，并只启用必要的功能。\n\n其次，我们可以通过更新和升级浏览器的版本来修复这个漏洞。浏览器厂商可以在新版本中修复该漏洞，并确保用户使用的浏览器具有最新的安全补丁和修复程序。因此，建议用户时刻关注并更新他们所使用的浏览器，以保持其安全性。\n\n另外，在网站和应用程序的开发过程中，开发者可以采取一些额外的安全措施来强化对漏洞的防范。这包括但不限于：对用户输入数据进行有效的验证和过滤、采用安全的编程实践、使用最新的安全框架和加密技术等。这些措施可以帮助开发者及早发现和修复潜在的漏洞。\n\n最后，对于已经被利用的漏洞，开发者应该尽快采取措施来修复并通知用户。及时披露漏洞和提供解决方案是至关重要的，以确保用户的信息和安全不受到进一步的威胁。\n\n综上所述，虽然“Permissions-Policy”头部存在一个漏洞，但我们可以通过谨慎选择和启用功能、升级浏览器、加强开发过程中的安全措施和及时修复已知漏洞来修复该漏洞。通过这些努力，可以更好地保护用户的隐私和数据安全。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"没有设置HttpOnly标志的Cookies",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"一个或多个Cookie没有设置HttpOnly标志。当使用HttpOnly标志设置Cookie时，它指示浏览器该Cookie只能由服务器访问，而不能由客户端脚本访问。对于会话Cookie来说，这是一种重要的安全保护措施。",
//             "res":"{\"pocname\":\"04b4867e-5297-23ff-c810-97acf15ed71a\"}",
//             "fix":"一个或多个Cookie没有设置HttpOnly标志会导致安全漏洞的出现，攻击者可以利用此漏洞来窃取用户的会话信息或伪造会话。因此，设置HttpOnly标志是一种重要的安全保护措施，可以提高应用程序的安全性。\n\n修复该漏洞的方法是为每个Cookie设置HttpOnly标志。在应用程序中，可以通过编程方式设置HttpOnly标志，确保在客户端脚本中无法访问和修改Cookie。在使用第三方库或框架时，要确保它们正确地实现了HttpOnly标志，并且该标志被正确地设置。\n\n除了设置HttpOnly标志之外，应用程序还应该使用HTTPS协议来保护会话信息的传输，以避免会话劫持或信息泄露。在应用程序中，还可以使用其他安全措施来保护会话和敏感信息，如使用CSRF令牌、检测异常登录等等。\n\n通过正确地设置HttpOnly标志和其他安全措施，可以最大程度地提高应用程序的安全性，保护用户的会话信息和个人数据，避免被黑客攻击。因此，开发者需要关注和修复此漏洞，确保应用程序的安全性和可靠性。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"可能找到虚拟主机",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"虚拟主机是一种在单个服务器（或服务器池）上托管多个域名（各自处理每个名称）的方法。这允许一个服务器共享其资源，例如内存和处理器周期，而无需所有提供的服务使用相同的主机名。当 Host 标头被操纵并测试各种常见的虚拟主机时，此 Web 服务器会做出不同的响应。这可能表明存在虚拟主机。",
//             "res":"{\"pocname\":\"aa32ac88-2c39-e6ac-07d8-3f02ad5c1b53\"}",
//             "fix":"为了修复上述描述的虚拟主机漏洞，我们可以采取以下措施。首先，服务器管理员应确保实施了正确的虚拟主机配置。这包括仔细设置各个虚拟主机的主机名和域名，并确保每个虚拟主机的配置是独立的，不会相互影响。\n\n其次，服务器应该对 Host 标头进行验证和过滤，以防止恶意操纵。可以通过限制 Host 标头的内容，仅允许合法的主机名和域名，来减少攻击面。此外，过滤掉特殊字符和符号，以防止利用输入验证漏洞。\n\n另外，服务器管理员还可以使用安全的虚拟主机软件和工具来帮助防止漏洞的发生。这些工具可以提供额外的安全层，例如入侵检测系统（IDS）和防火墙，在虚拟主机层面上检测和阻止任何恶意的请求。\n\n定期更新和升级服务器软件和操作系统也是非常重要的。漏洞通常与旧版本的软件相关，因此确保服务器上安装的软件和操作系统始终是最新的，以减少攻击者利用已知漏洞的可能性。\n\n此外，服务器管理员还应该定期进行安全审计和漏洞扫描，以发现任何可能的漏洞并及时修复。定期审查服务器配置和日志也有助于及早发现任何异常活动。\n\n最后，持续的员工培训和意识提高也是至关重要的。在员工中普及安全意识和最佳实践，教育他们如何避免潜在的安全风险，可以大大减少漏洞的发生。\n\n通过以上措施的综合应用，可以帮助修复和预防虚拟主机漏洞的发生，保护服务器和托管的域名免受潜在的攻击和侵犯。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"TLS/SSL证书即将过期。",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"您的服务器使用的一个TLS/SSL证书即将过期。\u003cbr/\u003e\u003cbr/\u003e一旦证书过期，大多数网络浏览器将向终端用户显示安全警告，要求他们手动确认您证书链的真实性。软件或自动系统可能会静默地拒绝连接到服务器。\u003cbr/\u003e\u003cbr/\u003e此警报不一定是由服务器（叶子）证书引起的，但可能是由中间证书触发的。请参考警报详细信息中的证书序列号来确定受影响的证书。",
//             "res":"{\"pocname\":\"029afcbb-3ec2-be3c-5cc3-29f5cfe016f4\"}",
//             "fix":"此漏洞需要立即采取措施进行修复，以避免可能的安全问题。修复措施包括更新证书，重新生成证书链，或更新中间证书。以下是一些可以采取的修复措施：\n\n1. 更新证书：从证书颁发机构（CA）获取最新的证书，并使用此证书替换受影响的证书。这将确保您的证书链不会过期，从而避免出现安全警告。\n\n2. 重新生成证书链：重新生成整个证书链，以确保其中的所有证书都是最新的，并且颁发者都是可信的。使用此修复措施的最终结果将是一个时间更长且更安全的证书链。\n\n3. 更新中间证书：如果中间证书引起了问题，您可以从CA处获取最新的中间证书，并使用此证书更新受影响的证书。这将修复中间证书问题，并确保您的证书链不会过期。\n\n总之，修复这个漏洞是非常重要的。在实施任何修复措施之前，请确保备份所有受影响的证书和证书链。如果您不确定如何执行修复，请咨询您的IT安全专家或者CA颁发机构，以获取正确的指导。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"未实现HTTP严格传输安全(HSTS)",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"\"HTTP严格传输安全性（HSTS）告知浏览器一个网站只能通过HTTPS方式访问。\n        \n检测到您的网络应用程序未实施HTTP严格传输安全性（HSTS），因为响应中缺少Strict Transport Security头部信息。\"",
//             "res":"{\"pocname\":\"34a6c791-c497-27d5-7272-6a968e9fdccb\"}",
//             "fix":"HTTP严格传输安全性（HSTS）被用来告知浏览器一个网站只能通过HTTPS方式访问，以帮助保护用户的敏感数据。如果您的网络应用程序未正确实施HSTS，攻击者可能利用该情况通过中间人攻击等方式将用户数据劫持或篡改。因此，及时修复HSTS漏洞至关重要。\n\n修复HSTS的方法主要包括两种。第一种，通过为响应添加Strict Transport Security头部信息来实现。在服务器设置中添加Strict-Transport-Security: max-age=31536000; includeSubDomains头部信息，其中max-age参数用于设置浏览器缓存HSTS策略的时间，而includeSubDomains参数则可以将HSTS策略应用于所有子域名。第二种方法是使用网站安全证书，将网站升级到HTTPS协议。在此过程中，服务器会自动将HSTS头部信息添加到所有信任的HTTPS通信中。\n\n无论哪种方法，修复HSTS漏洞都需要涉及到服务器端的配置和修改。因此，在实施修复措施之前，一定要进行充分的测试并备份相关数据。此外，及时禁用或删除不再使用的HSTS策略也是必要的，以确保服务器的安全性。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"未实施内容安全策略 (CSP)",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"'内容安全策略（Content Security Policy，CSP）是一种增加的安全层，有助于检测和缓解某些类型的攻击，包括跨站脚本（XSS）和数据注入攻击。\u003cbr/\u003e\u003cbr/\u003e\n\n可以通过添加 \u003cstrong\u003eContent-Security-Policy\u003c/strong\u003e 头来实现内容安全策略（CSP）。该头的值是一个包含策略指令的字符串，用于描述您的内容安全策略。要实现CSP，您应该为站点使用的所有资源类型定义允许的来源列表。例如，如果您有一个简单的站点需要本地加载脚本、样式表和图片以及来自jQuery库的CDN，CSP头应该如下所示：\n\n\u003cpre\u003e\u003ccode\u003e\nContent-Security-Policy:\n    default-src 'self';\n    script-src 'self' https://code.jquery.com;\n\u003c/code\u003e\u003c/pre\u003e\n\n\u003cbr/\u003e\u003cbr/\u003e\n\n检测到您的Web应用程序没有实现内容安全策略（CSP），因为响应中缺少CSP头。建议将内容安全策略（CSP）实现到您的Web应用程序中。'",
//             "res":"{\"pocname\":\"011055fc-94f1-ab96-56ac-53117c56fb4d\"}",
//             "fix":"该漏洞描述了网站安全中的内容安全策略（Content Security Policy，CSP）实现问题。为了增加网站的安全性，建议在网站中实现CSP。具体做法是添加Content-Security-Policy头，并为站点中使用的所有资源类型定义允许的来源列表。例如，对于一个简单的站点，应该定义本地加载脚本、样式表和图片，以及来自jQuery库的CDN允许的来源。这些来源被写入CSP头，作为策略指令的字符串。这样可以检测和缓解一些类型的攻击，包括跨站脚本（XSS）和数据注入攻击。如果您的Web应用程序没有实现内容安全策略（CSP），建议您立即实现。这样可以确保您的网站更加安全，令攻击者难以突破网站安全防御。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"\"无HTTP重定向\"",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"检测到您的Web应用程序使用HTTP协议，但不会自动重定向用户到HTTPS。",
//             "res":"{\"pocname\":\"3f6a8a0e-07f2-af81-54ff-61020299caeb\"}",
//             "fix":"该漏洞描述了一个常见的web应用程序安全问题，即应用程序使用HTTP协议而不是HTTPS协议。HTTP协议是一种不安全的传输协议，未加密的数据包会通过网络传递。这使得攻击者可以轻易地读取、修改和窃取这些数据，并对用户隐私造成威胁。因此，为了确保应用程序使用最高级别的安全措施，就需要使用HTTPS协议代替HTTP协议。HTTPS协议使用传输层安全（TLS）或安全套接字层（SSL）协议来加密数据传输。这样，只有受信任的接收方可以读取和使用传输数据。\n\n为了修复这个漏洞，应用程序开发人员需要将HTTP协议更改为HTTPS协议。他们可以在应用程序配置文件中为使用HTTPS协议的路径添加重定向规则，这将确保用户始终通过HTTPS协议访问应用程序。此外，开发人员应可以使用SSL证书或其他安全措施来保护传输过程中的数据。这可以通过将SSL证书绑定到应用程序域名上来实现，或使用基于云的安全服务来实现。这些安全措施将确保用户隐私的安全，并使攻击者无法查看或窃取数据。在应用程序修复这个漏洞后，网站管理员应该定期进行安全测试和漏洞扫描，以确保应用程序继续保持安全并且满足最高级别的安全标准。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"缺少、不一致或相互矛盾属性的Cookie。",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"“以下至少一个cookie属性会导致该cookie无效或与同一cookie的其他属性或使用环境不兼容。虽然这本身不是漏洞，但它很可能会导致应用程序的意外行为，进而可能引起二次安全问题。”",
//             "res":"{\"pocname\":\"fd1ec03e-319d-ae3b-425f-d8c1275619d8\"}",
//             "fix":"该漏洞描述了一种可能会导致应用程序出现意外行为的情况，但其本身并不是漏洞。 然而，由于该情况可能引起二次安全问题，因此开发人员和管理员需要注意这个问题，并在必要时采取修复措施。\n\n为了避免出现不兼容的Cookie属性导致的意外行为和二次安全问题，开发人员应该仔细检查应用中使用的Cookie属性，并确保它们与应用程序使用的其他部分兼容。\n\n一种解决方法是使用Cookie安全标志，这将强制浏览器只在安全连接上发送Cookie。另外，可以采用限制cookie范围的方法，在特定的域名或路径下使用Cookie，以此来限制Cookie的使用范围，从而减少意外行为和二次安全问题的发生。\n\n最重要的是，开发人员应该经常进行漏洞扫描和安全测试，以及对代码进行定期的审查，从而及时发现和修复可能存在的问题。这可以确保应用程序始终保持安全，并减少因错误配置和使用Cookie属性导致的潜在风险。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"未设置安全标志的Cookie",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"一个或多个cookie没有设置安全标志。当一个cookie被设置了安全标志后，它会指示浏览器该cookie只能在安全的SSL/TLS通道上访问。这是会话cookie重要的安全保护措施。",
//             "res":"{\"pocname\":\"575ba039-89b6-b7d7-eba6-553135c04553\"}",
//             "fix":"该漏洞的修复措施旨在确保会话cookie的安全保护措施得到有效实施。对于未设置安全标志的cookie，开发人员应立即对其进行修改，将其设置上安全标志标识。通过这种方法，开发人员可以消除会话cookie的任何潜在风险，防止恶意用户在传输过程中劫持用户身份，从而造成安全威胁。\n\n为确保实现安全标志标识的cookie准确性，开发人员还需要通过SSL/TLS访问控制来限制访问网络资源的用户身份。通过实施包括内容保护和访问控制的安全策略，开发人员可以确保会话cookie的安全性得到最大程度的保障。在进行此类操作的过程中，开发人员应仔细跟踪会话cookie的使用情况，及时发现和修复任何与设置安全标志标识相关的错误或漏洞。\n\n总之，对会话cookie设置安全标志标识是防范安全风险必不可少的一项措施。通过限制cookie访问范围，加强访问控制和内容保护等操作，开发人员可以为用户提供更加安全可靠的网络环境，保护客户的私人和机密信息。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"点击劫持：X-Frame-Options头部",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":2,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"点击劫持（用户界面篡改攻击，UI篡改攻击，UI篡改）是一种恶意技术，可以欺骗Web用户点击与其所认为点击不同的内容，从而可能泄露机密信息或在点击看似无害的网页时控制他们的计算机。服务器未返回具有DENY或SAMEORIGIN值的X-Frame-Options头，这意味着该网站可能面临点击劫持攻击的风险。 X-Frame-Options HTTP响应头可用于指示是否应允许浏览器在框架或iframe中呈现页面。站点可以使用此功能来避免点击劫持攻击，确保其内容不嵌入不受信任的站点。",
//             "res":"{\"pocname\":\"b8e2c082-44f1-cf0b-0b8e-0e0bb357e798\"}",
//             "fix":"点击劫持是一种恶意技术，可以控制用户的计算机或者泄漏机密信息。这种攻击方式会欺骗用户点击看似无害的网页，但是实际上它已经被篡改了。这种攻击技术对网站安全造成了严重威胁。为了防止点击劫持攻击，我们可以使用X-Frame-Options HTTP响应头来避免嵌入不信任的站点。如果服务器没有返回具有DENY或SAMEORIGIN值的X-Frame-Options头，就存在被点击劫持的风险。\n\nX-Frame-Options设置框架的漏洞，是通过消除非授权嵌入的方式来解决的。它的主要作用是避免网站被嵌到不受信任的站点。为了设置这个功能，我们可以在服务器上添加以下代码：\n\nX-Frame-Options: DENY\n\n或\n\nX-Frame-Options: SAMEORIGIN\n\n这样就可以避免网站面临点击劫持攻击的风险了。如果站点未修改这个设置，那么可能会面临许多安全问题。因此，我们应该要了解这个漏洞，并采取相应的措施来保护我们的网站和用户不受攻击威胁。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         },
//         {
//             "vulName":"Unencrypted connection的中文翻译为“非加密连接”。",
//             "risk":"低危",
//             "vulStatus":"验证存在",
//             "number":1,
//             "type":"其他",
//             "cve":"",
//             "publishDate":"",
//             "describe":"这个扫描目标连接是未加密的。一个潜在的攻击者可以拦截和修改从这个站点发送和接收的数据。",
//             "res":"{\"pocname\":\"bca221d1-8581-3375-4097-66b0048ed088\"}",
//             "fix":"该漏洞是因为扫描目标连接未加密，导致攻击者可以拦截和修改该站点发送和接收的数据，因此需要采取一定的措施进行修复。通常，可以通过SSL（安全套接字层）协议来解决这个问题，SSL协议是一种加密协议，可以保护网络传输过程中的数据安全。\n\n对于站点管理者来说，他们需要考虑启用SSL证书以支持HTTPS协议，HTTPS协议相比HTTP协议是基于SSL加密的，在传输过程中可以将数据进行加密传输，确保数据的安全性。并且站点管理者可以通过定期的安全检查以及更新站点的安全软件，及时发现和修复可能存在的漏洞，保证站点的安全性。\n\n对于用户来说，他们需要注意防范网络钓鱼以及恶意软件的攻击，尽可能选择使用SSL加密的站点，同时多使用杀毒软件和防火墙等安全软件，增强自身网络安全意识。在使用公共网络时，为了避免个人账号和密码等敏感信息的泄露，建议不要随意使用未知的公共网络，并且尽可能不在公共网络中登录个人账号。\n\n综上所述，修复该漏洞的关键在于采取措施来加强站点的安全保障体系，包括启用SSL证书、定期更新安全软件、及时发现和修复可能存在的漏洞等。同时，用户也需要更加注重自身的网络安全，选择SSL加密的站点、使用安全软件，避免在公共网络中登录敏感信息，从而保障网络的安全。",
//             "affectRange":"",
//             "AffectTargets":"cheilmc.co.kr",
//             "location":"",
//             "link":""
//         }
//     ]
// }
// ];
                // this.hyAllData = JSON.parse(res.data.content);
                console.log(this.reportcontent);
                if(row.type ==1 || row.type ==3 || row.type ==5){ //任务--总数
                    if(row.format == 1){
                        //html
                        this.downhtml(row, report_css_new);
                    }
                    if(row.format == 2){
                        //word
                        this.downword(row);
                    }
                    if(row.format == 3){
                        //pdf
                        const routeData = this.$router.resolve({
                            path: '/taskreportview',
                            query: { id: row.id, type:'pdf'}
                        })
                        window.open(routeData.href, '_blank')
                    }
                    if(row.format == 4){
                        //excel
                        this.taskExcel();
                    }
                    if(row.format == 5){
                        //csv
                        this.taskCSV();
                    }
                }
                if(row.type ==2 || row.type ==4 || row.type ==6){//目标
                    if(row.format == 1){
                        //html
                        this.downloadhtmlTarget(row);
                    }
                    if(row.format == 2){
                        //word
                        this.downwordtarget(row);
                    }
                    if(row.format == 3){
                        //pdf
                        const routeData = this.$router.resolve({
                            path: '/targetreportview',
                            query: { id: row.id, type:'pdf'}
                        })
                        window.open(routeData.href, '_blank')
                    }
                    if(row.format == 4){ 
                        // excel
                        this.targetExcel() 
                    }
                    if(row.format == 5){
                        //csv
                        this.targetCSV();
                    }
                }
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            } 
            
        },

        base64DataURLToArrayBuffer(dataURL) {
            const base64Regex = /^data:image\/(png|jpg|jpeg|svg|svg\+xml);base64,/;
            if (!base64Regex.test(dataURL)) {
                return false;
            }
            const stringBase64 = dataURL.replace(base64Regex, "");
            let binaryString;
            if (typeof window !== "undefined") {
                binaryString = window.atob(stringBase64);
            } else {
                binaryString = new Buffer(stringBase64, "base64").toString("binary");
            }
            const len = binaryString.length;
            const bytes = new Uint8Array(len);
            for (let i = 0; i < len; i++) {
                const ascii = binaryString.charCodeAt(i);
                bytes[i] = ascii;
            }
            return bytes.buffer;
        },
        getImg(_data) {
            // 生成echarts图表
            let option = { 
                animation: false ,
                color:['#F87D7D','#F9B640','#4C7AE3','#15C53D'],
                tooltip: { 
                    trigger: 'item',
                    formatter: "{b}: {c} ({d}%)"
                },   
                series: [
                    {
                        name:'',
                        type:'pie', 
                        radius: ['60%', '80%'],
                        avoidLabelOverlap: false,
                        data: this.chartData,
                        label: {
                            normal: {
                                formatter: function(params, ticket, callback) {
                                    return params.name + '：\n ('+params.percent+'%)';
                                },
                                position: 'inner',
                                show : false
                            },
                        },

                        tooltip: {
                            trigger: 'item',
                            formatter: "{b}: {c} ({d}%)"
                        },
                    }
                ]
            } ; 
            let myEchart = echarts.init(document.getElementById("targetschartpie"));
            myEchart.setOption(option); // option是echarts图表
            // 获取echart图片
            this.imgUrl = myEchart.getDataURL({
                pixelRatio: 1,// 导出的图片分辨率比例，默认为 1。
                backgroundColor: '#fff',// 导出的图片背景色，默认使用 option 里的 backgroundColor
                type: 'png',//图片类型支持png和jpeg
            });  
        },
        downword(row){
            let _this = this;
            let targetRisk = this.reportcontent.targetRisk;
            this.chartData = [
                {name: ' 高危目标      ',value:2},
                {name: ' 中危目标      ',value:2},
                {name: ' 低危目标      ',value:2},
                {name: ' 安全目标      ',value:2}
            ] 
            // this.chartData = [
            //     {name: ' 高危目标      ',value:targetRisk.highNumber},
            //     {name: ' 中危目标      ',value:targetRisk.MiddleNumber},
            //     {name: ' 低危目标      ',value:targetRisk.lowNumber},
            //     {name: ' 安全目标      ',value:targetRisk.safeNumber}
            // ] 
            this.getImg(); 
            this.reportcontent.imgUrl = '';
            let ImageModule = require('docxtemplater-image-module');
             
            // 读取并获得模板文件的二进制内容
            // 模板位置，public文件夹下
            JSZipUtils.getBinaryContent('task.docx', (error, content) => {
                // 抛出异常
                if (error) throw error;

                let opts = {
                    centered: false,
                    fileType: "docx"
                }
                opts.getImage = function (tagValue) {
                    if (tagValue.size && tagValue.data) {
                    return _this.base64DataURLToArrayBuffer(tagValue.data);
                    }
                    return _this.base64DataURLToArrayBuffer(tagValue);
                };
                opts.getSize = function (_, tagValue) {
                    if (tagValue.size && tagValue.data) {
                        return tagValue.size;
                    }
                    return [200, 200];
                }; 
                this.reportcontent.imgUrl = _this.imgUrl;  
                // 创建一个JSZip实例，内容为模板的内容
                let zip = new JSZip2(content);  
  
                const expressionParser = require("./utils/expressions.js");
                expressionParser.filters.where = function (input, query) {
                    return input.filter(function (item) {
                        return expressions.compile(query)(item);
                    });
                };
                 
                // 创建并加载docxtemplater实例对象
                // 3.5版本可以下载图片
                let doc = new Docxtemplater();  
                doc.attachModule(new ImageModule(opts)); 
                doc.loadZip(zip).setOptions({ parser: expressionParser });
                // 设置模板变量的值
                doc.setData({
                    ...this.reportcontent, 
                });
                try {
                    // 用模板变量的值替换所有模板变量
                    doc.render();
                } catch (error) {
                // 抛出异常
                let e = {
                    message: error.message,
                    name: error.name,
                    stack: error.stack,
                    properties: error.properties
                };
                console.log(JSON.stringify({ error: e }));
                throw error;
                }
                // 生成一个代表docxtemplater对象的zip文件（不是一个真实的文件，而是在内存中的表示）
                let out = doc.getZip().generate({
                    type: 'blob',
                    mimeType: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
                });
                // 将目标文件对象保存为目标类型的文件，并命名
                saveAs(out, this.reportcontent.reportCover.title+'.docx');
            }); 
            
        },
        downwordtarget(){ 
            // 读取并获得模板文件的二进制内容
            // 模板位置，public文件夹下
            JSZipUtils.getBinaryContent('/target.docx', (error, content) => {
                // 抛出异常
                if (error) throw error; 
                const expressionParser = require("./utils/expressions.js");
                expressionParser.filters.where = function (input, query) {
                    return input.filter(function (item) {
                        return expressions.compile(query)(item);
                    });
                };
                // 创建一个JSZip实例，内容为模板的内容
                let zip = new JSZip2(content);
                // 创建并加载docxtemplater实例对象
                let doc = new Docxtemplater();  
               
                doc.loadZip(zip).setOptions({ parser: expressionParser });
                // 设置模板变量的值
                doc.setData({
                    ...this.reportcontent, 
                });
                try {
                    // 用模板变量的值替换所有模板变量
                    doc.render();
                } catch (error) {
                // 抛出异常
                let e = {
                    message: error.message,
                    name: error.name,
                    stack: error.stack,
                    properties: error.properties
                }; 
                throw error;
                }
                // 生成一个代表docxtemplater对象的zip文件（不是一个真实的文件，而是在内存中的表示）
                let out = doc.getZip().generate({
                    type: 'blob',
                    mimeType: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
                });
                // 将目标文件对象保存为目标类型的文件，并命名
                saveAs(out, this.reportcontent.reportCover.title+'.docx');
            }); 
        },
        async downhtml(row, report_css_new){
            let targetRisk = this.reportcontent.targetRisk;
            // this.chartData = [
            //     {name: ' 高危目标      ',value:2},
            //     {name: ' 中危目标      ',value:2 },
            //     {name: ' 低危目标      ',value:2 },
            //     {name: ' 安全目标      ',value:2}
            // ] 
            this.chartData = [
                {name: ' 高危目标      ',value:targetRisk.highNumber||0},
                {name: ' 中危目标      ',value:targetRisk.MiddleNumber||0},
                {name: ' 低危目标      ',value:targetRisk.lowNumber||0},
                {name: ' 安全目标      ',value:targetRisk.safeNumber||0}
            ] 
            let charsOptionStr = { 
                animation: false ,
                color:['#F87D7D','#F9B640','#4C7AE3','#15C53D'],
                tooltip: { 
                    trigger: 'item',
                    formatter: "{b}: {c} ({d}%)"
                },   
                legend: {
                    type: 'scroll',
                    orient: 'vertical',
                    textStyle: {
                        color: '#484866'          // 图例文字颜色
                    },
                    right: 100,
                    top: 120,
                    bottom: 50,
                    itemGap: 35,
                    itemWidth:6,
                    itemHeight:6,
                    icon: "circle",  // 使用回调函数  
                },    
                series: [
                    {
                        name:'',
                        type:'pie',  
                        radius: ['50%', '65%'],
                        avoidLabelOverlap: false,
                        data: this.chartData,
                        label: {
                            normal: {
                                formatter: function(params, ticket, callback) {
                                    return params.name + '：\n ('+params.percent+'%)';
                                },
                                position: 'inner',
                                show : false
                            },
                        },

                        tooltip: {
                            trigger: 'item',
                            formatter: "{b}: {c} ({d}%)"
                        },
                    }
                ]
            }; 
            const zip = new JSZip()  
            let foldercss = zip.folder("css");
            let folderjs = zip.folder("js"); 
            var promises = [];  
            promises.push(folderjs.file('echarts_min.js', echarts_min, { binary: true })); 
            promises.push(foldercss.file('report.css', report_css_new, { binary: true })); 
            promises.push(foldercss.file('bootstrap.css', bootstrap_css, { binary: true })); 
            promises.push(foldercss.file('iconfont.css', iconfont, { binary: true })); 
            
            let _htmlfilename = row.name+".html";
            let _zipname = row.name+".zip";
            console.log(_htmlfilename)
            // 任务模板
            
                
                // this.hyAllData.forEach(async(item,index)=>{
                    await import("./utils/htmlTaskTemp.js").then((modules) => { 
                        var htmlFiles = modules.htmlTemplate
                        htmlFiles = htmlFiles.replace(`%chart_tmp%`, JSON.stringify(charsOptionStr));
                        htmlFiles = htmlFiles.replace(`%data%`, JSON.stringify(this.reportcontent));
                        console.log(this.reportcontent,'this.reportcontent!!!!!!!!');
                         promises.push(zip.file('index.html', htmlFiles)) 
                          });
                // })
           
            Promise.all(promises).then(() => {
                zip.generateAsync({ 
                    type: "blob" ,
                    responseType: 'arraybuffer'
                }).then((content) => {
                    // 生成二进制流
                    
                    FileSaver.saveAs(content, _zipname ); // 利用file-saver保存文件  自定义文件名
                });
            }).catch((res) => {
                alert("文件压缩失败!!!");
            }); 

        },
        async downloadhtmlTarget(){ 
            const zip = new JSZip() 
            let foldercss = zip.folder("css");
            let folderjs = zip.folder("js"); 
            var promises = []; 

            promises.push(folderjs.file('echarts_min.js', echarts_min, { binary: true })); 
            promises.push(foldercss.file('report.css', report_css, { binary: true })); 
            promises.push(foldercss.file('bootstrap.css', bootstrap_css, { binary: true })); 
            promises.push(foldercss.file('iconfont.css', iconfont, { binary: true })); 

            function escapeHtml(unsafe) {
                return unsafe 
                    .replace(/</g, "&lt;")
                    .replace(/>/g, "&gt;") 
            }

        
            await import("./utils/htmlTargetTemp.js").then((modules) => { 
                let htmlFiles = modules.htmlTargetTemplate 
                .replace(`%data%`, escapeHtml(JSON.stringify(this.reportcontent)))
            
                ;   
                promises.push(zip.file("index.html", htmlFiles)) 
            });
           
            Promise.all(promises).then(() => { 
                zip.generateAsync({ 
                    type: "blob" ,
                    responseType: 'arraybuffer'
                }).then((content) => {
                    console.log(1)
                    // 生成二进制流
                    FileSaver.saveAs(content, this.reportName+".zip"); // 利用file-saver保存文件  自定义文件名
                });
            }).catch((res) => {
                alert("文件压缩失败");
            }); 



        },
        taskExcel(){ //任务导出excel
            const workbook = new ExcelJS.Workbook();
            let taskOverview = this.reportcontent.taskOverview;
            let targetDetails = this.reportcontent.targetDetails;
            let vulDetails  = this.reportcontent.vulDetails;
            this.reportcontent.catalogParent.forEach(item=>{
                if(item.id == 'taskOverview'){ //概要
                    const worksheet = workbook.addWorksheet(item.name); 
                    worksheet.addRow(['任务名称', '风险等级', '目标分布','漏洞分布','漏洞验证','任务场景','测试时间']);
                    worksheet.addRow([
                        taskOverview.taskName, 
                        taskOverview.taskRiskStr, 
                        '总目标：'+taskOverview.targetStat.total+'，存活目标：'+taskOverview.targetStat.liveTarget+'；高危目标：'+taskOverview.targetStat.HighTarget+'；中危目标'+taskOverview.targetStat.middleTarget+'；低危目标'+taskOverview.targetStat.lowTarget+'；安全目标：'+taskOverview.targetStat.safeTarget,
                        '漏洞总数：'+taskOverview.vulnStat.total+'；致命漏洞：'+taskOverview.vulnStat.deadlyNumber+'，高危漏洞：'+taskOverview.vulnStat.highNumber+'，中危漏洞：'+taskOverview.vulnStat.middleNumber+'，低危漏洞：'+taskOverview.vulnStat.lowNumber,
                        '验证成功：'+taskOverview.vulnVerify.verifySuccess+'，利用成功：'+taskOverview.vulnVerify.useSuccess+'，未验证：'+taskOverview.vulnVerify.repairSuccess,
                        taskOverview.templateName,
                        taskOverview.createDate
                    ]);
                }
                if(item.id == 'taskStat'){
                    let targetRisk = this.reportcontent.targetRisk;
                    let vulRisk = this.reportcontent.vulRisk;
                    let vulType = this.reportcontent.vulType;
                    let topVulRisk = this.reportcontent.topVulRisk;

                    item.catalog.forEach(child=>{
                        if(child.id == 'targetRisk' ){
                            const worksheet = workbook.addWorksheet(child.name); 
                            worksheet.addRow(['目标类型','数量','占比']);
                            worksheet.addRow(['高危目标',targetRisk.highNumber,targetRisk.highNumberRate]);
                            worksheet.addRow(['中危目标',targetRisk.MiddleNumber,targetRisk.MiddleNumberRate]);
                            worksheet.addRow(['低危目标',targetRisk.lowNumber,targetRisk.lowNumberRate]);
                            worksheet.addRow(['安全目标',targetRisk.safeNumber,targetRisk.safeNumberRate]);
                        }
                        if(child.id == 'vulRisk'){
                            const worksheet = workbook.addWorksheet(child.name); 
                            worksheet.addRow(['风险类型','验证存在','未验证','利用成功','漏洞总数','漏洞类型占比']);
                            vulRisk.forEach(risk=>{
                                worksheet.addRow([
                                    risk.riskType,
                                    risk.verifySuccess,
                                    risk.repairSuccess,
                                    risk.useSuccess,
                                    risk.total,
                                    risk.percent
                                ])
                            })
                        }
                        if(child.id == 'vulType'){
                            const worksheet = workbook.addWorksheet(child.name); 
                            worksheet.addRow(['漏洞类型','数量','占比','影响目标数量']);
                            vulType.forEach(type=>{
                                worksheet.addRow([
                                    type.vulnType,
                                    type.total,
                                    type.percent,
                                    type.targetNumber
                                ])
                            })
                        }
                        if(child.id == 'topVulRisk'){
                            const worksheet = workbook.addWorksheet(child.name);
                            worksheet.addRow(['漏洞名称','漏洞风险','出现次数','影响目标']); 
                            topVulRisk.forEach(top=>{
                                worksheet.addRow([
                                    top.vulName,
                                    top.risk,
                                    top.number,
                                    top.affectTargets
                                ])
                            })
                        } 
                    })
                }
                if(item.id == 'targetDetails'){
                    const worksheet = workbook.addWorksheet(item.name);
                    worksheet.addRow(['测试目标','目标风险','致命漏洞','高危漏洞','中危漏洞','低危漏洞','漏洞状态']); 
                    targetDetails.forEach(detail=>{
                        worksheet.addRow([
                            detail.target,
                            detail.risk,
                            detail.deadlyNumber,
                            detail.highNumber,
                            detail.middleNumber,
                            detail.lowNumber,
                            detail.vulStatus
                        ]);
                    })
                }
                if(item.id == 'vulDetails'){
                    const worksheet = workbook.addWorksheet(item.name);
                    worksheet.addRow(['漏洞名称', '漏洞风险', '漏洞状态','漏洞类型','漏洞编号','公开日期','漏洞描述','修复建议','影响范围','影响目标','参考链接','出现次数']);
                    if(vulDetails){
                        vulDetails.forEach(vulInfo=>{
                            worksheet.addRow([
                                vulInfo.vulName,
                                vulInfo.risk,
                                vulInfo.vulStatus,
                                vulInfo.type,
                                vulInfo.cve,
                                vulInfo.publishDate,
                                vulInfo.describe, 
                                vulInfo.fix,
                                vulInfo.affectRange,
                                vulInfo.AffectTargets, 
                                vulInfo.link,
                                vulInfo.number
                            ]);
                        })
                        
                    }
                }
            }) 
            workbook.xlsx.writeBuffer().then((buffer) => {
                const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
                saveAs(blob, this.reportcontent.reportCover.title+'.xlsx');
            }); 
        },
        targetExcel(){ //目标导出excel 
            const workbook = new ExcelJS.Workbook();
            let targetOverview = this.reportcontent.targetOverview;
            let assetInfo = this.reportcontent.assetInfo;
            let vulInfo = this.reportcontent.vulInfo; 

            this.reportcontent.catalogParent.forEach(item=>{ 
                if(item.id == 'targetOverview'){ //概要
                    const worksheet = workbook.addWorksheet(item.name); 
                    worksheet.addRow(['测试目标', '风险等级', '漏洞分布','漏洞状态','测试时间']);
                    worksheet.addRow([
                        targetOverview.target, 
                        targetOverview.risk, 
                        '漏洞总数：'+targetOverview.vulnStat.total+'；致命漏洞：'+targetOverview.vulnStat.deadlyNumber+'，高危漏洞：'+targetOverview.vulnStat.highNumber+'，中危漏洞：'+targetOverview.vulnStat.middleNumber+'，低危漏洞：'+targetOverview.vulnStat.lowNumber,
                        '验证成功：'+targetOverview.vulnVerify.verifySuccess+'，利用成功：'+targetOverview.vulnVerify.useSuccess+'，未验证：'+targetOverview.vulnVerify.repairSuccess
,
                        targetOverview.createDate
                    ]);
                }
                if(item.id == 'assetInfo'){ //资产信息
                    const worksheet = workbook.addWorksheet(item.name); 
                    worksheet.addRow(['组件/指纹', '服务', 'IP/域名','操作系统']);
                    if(assetInfo){ 
                        worksheet.addRow([assetInfo.component,  assetInfo.service, assetInfo.ipOrUrl,assetInfo.system]);
                    }
                }

                if(item.id == 'vulInfo'){ //漏洞信息
                    const worksheet = workbook.addWorksheet(item.name); 
                    worksheet.addRow(['漏洞名称', '漏洞风险', '漏洞状态','漏洞类型','漏洞编号','披露日期','漏洞描述','漏洞结果','修复建议','影响范围','漏洞位置','参考链接','playload','请求报文','响应报文']);
                    if(vulInfo){
                        vulInfo.forEach(vul=>{
                            worksheet.addRow([
                            vul.vulName,
                            vul.risk,
                            vul.vulStatus,
                            vul.type,
                            vul.cve,
                            vul.publishDate,
                            vul.describe,
                            vul.res,
                            vul.fix,
                            vul.affectRange,
                            vul.location,
                            vul.link,
                            vul.verMsg.payload,
                            vul.verMsg.request,
                            vul.verMsg.response,
                        ]);
                        })
                        
                    }
                }

            })  
            workbook.xlsx.writeBuffer().then((buffer) => {
                const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
                    
                saveAs(blob, this.reportcontent.reportCover.title+'.xlsx');
            }); 
        },
        taskCSV(){
            let data=[];
            let taskOverview = this.reportcontent.taskOverview;
            let targetDetails = this.reportcontent.targetDetails;
            let vulDetails  = this.reportcontent.vulDetails;
            this.reportcontent.catalogParent.forEach(item=>{
                if(item.id == 'taskOverview'){ //概要
                    data.push(['任务名称', '风险等级', '目标分布','漏洞分布','漏洞验证','任务场景','测试时间']);
                    data.push([
                        taskOverview.taskName, 
                        taskOverview.taskRiskStr, 
                        '总目标：'+taskOverview.targetStat.total+'，存活目标：'+taskOverview.targetStat.liveTarget+'；高危目标：'+taskOverview.targetStat.HighTarget+'；中危目标'+taskOverview.targetStat.middleTarget+'；低危目标'+taskOverview.targetStat.lowTarget+'；安全目标：'+taskOverview.targetStat.safeTarget,
                        '漏洞总数：'+taskOverview.vulnStat.total+'；致命漏洞：'+taskOverview.vulnStat.deadlyNumber+'，高危漏洞：'+taskOverview.vulnStat.highNumber+'，中危漏洞：'+taskOverview.vulnStat.middleNumber+'，低危漏洞：'+taskOverview.vulnStat.lowNumber,
                        '验证成功：'+taskOverview.vulnVerify.verifySuccess+'，利用成功：'+taskOverview.vulnVerify.useSuccess+'，未验证：'+taskOverview.vulnVerify.repairSuccess,
                        taskOverview.templateName,
                        taskOverview.createDate
                    ]);
                    data.push([]);
                }
                if(item.id == 'taskStat'){
                    let targetRisk = this.reportcontent.targetRisk;
                    let vulRisk = this.reportcontent.vulRisk;
                    let vulType = this.reportcontent.vulType;
                    let topVulRisk = this.reportcontent.topVulRisk; 
                    item.catalog.forEach(child=>{
                        if(child.id == 'targetRisk' ){ 
                            data.push(['目标类型','数量','占比']);
                            data.push(['高危目标',targetRisk.highNumber,targetRisk.highNumberRate]);
                            data.push(['中危目标',targetRisk.MiddleNumber,targetRisk.MiddleNumberRate]);
                            data.push(['低危目标',targetRisk.lowNumber,targetRisk.lowNumberRate]);
                            data.push(['安全目标',targetRisk.safeNumber,targetRisk.safeNumberRate]);
                        }
                        if(child.id == 'vulRisk'){  
                            data.push(['风险类型','验证存在','未验证','利用成功','漏洞总数','漏洞类型占比']);
                            vulRisk.forEach(risk=>{
                                data.push([
                                    risk.riskType,
                                    risk.verifySuccess,
                                    risk.repairSuccess,
                                    risk.useSuccess,
                                    risk.total,
                                    risk.percent
                                ])
                            })
                        }
                        if(child.id == 'vulType'){ 
                            data.push(['漏洞类型','数量','占比','影响目标数量']);
                            vulType.forEach(type=>{
                                data.push([
                                    type.vulnType,
                                    type.total,
                                    type.percent,
                                    type.targetNumber
                                ])
                            })
                        }
                        if(child.id == 'topVulRisk'){ 
                            data.push(['漏洞名称','漏洞风险','出现次数','影响目标']); 
                            topVulRisk.forEach(top=>{
                                data.push([
                                    top.vulName,
                                    top.risk,
                                    top.number,
                                    top.affectTargets
                                ])
                            })
                        } 
                        data.push([]);
                    });
                    data.push([]);
                }
                if(item.id == 'targetDetails'){ 
                    data.push(['测试目标','目标风险','致命漏洞','高危漏洞','中危漏洞','低危漏洞','漏洞状态']); 
                    targetDetails.forEach(detail=>{
                        data.push([
                            detail.target,
                            detail.risk,
                            detail.deadlyNumber,
                            detail.highNumber,
                            detail.middleNumber,
                            detail.lowNumber,
                            detail.vulStatus
                        ]);
                    })
                    data.push([]);
                }
                if(item.id == 'vulDetails'){ 
                    data.push(['漏洞名称', '漏洞风险', '漏洞状态','漏洞类型','漏洞编号','公开日期','漏洞描述','修复建议','影响范围','影响目标','参考链接','出现次数']);
                    if(vulDetails){
                        vulDetails.forEach(vulInfo=>{
                            data.push([
                                vulInfo.vulName,
                                vulInfo.risk,
                                vulInfo.vulStatus,
                                vulInfo.type,
                                vulInfo.cve,
                                vulInfo.publishDate,
                                vulInfo.describe, 
                                vulInfo.fix,
                                vulInfo.affectRange,
                                vulInfo.AffectTargets, 
                                vulInfo.link,
                                vulInfo.number
                            ]);
                        })
                        
                    }
                    data.push([]);
                }
            }); 
            const csv = data.join('\n');
 
            const blob = new Blob([`\ufeff${csv}`], { type: 'text/csv;charset=utf-8;' }); 
            saveAs(blob, this.reportcontent.reportCover.title+'.csv');
        },
        targetCSV(){
            let data=[];
            let targetOverview = this.reportcontent.targetOverview;
            let assetInfo = this.reportcontent.assetInfo;
            let vulInfo = this.reportcontent.vulInfo;

            this.reportcontent.catalogParent.forEach(item=>{ 
                if(item.id == 'targetOverview'){ //概要 
                    data.push(['测试目标', '风险等级', '漏洞分布','漏洞状态','测试时间']);
                    data.push([
                        targetOverview.target, 
                        targetOverview.risk, 
                        '漏洞总数：'+targetOverview.vulnStat.total+'；致命漏洞：'+targetOverview.vulnStat.deadlyNumber+'，高危漏洞：'+targetOverview.vulnStat.highNumber+'，中危漏洞：'+targetOverview.vulnStat.middleNumber+'，低危漏洞：'+targetOverview.vulnStat.lowNumber,
                        '验证成功：'+targetOverview.vulnVerify.verifySuccess+'，利用成功：'+targetOverview.vulnVerify.useSuccess+'，未验证：'+targetOverview.vulnVerify.repairSuccess
,
                        targetOverview.createDate
                    ]);
                    data.push([]);
                }
                if(item.id == 'assetInfo'){ //资产信息 
                    data.push(['组件/指纹', '服务', 'IP/域名','操作系统']);
                    if(assetInfo){ 
                        data.push([assetInfo.component,  assetInfo.service, assetInfo.ipOrUrl,assetInfo.system]);
                    }
                    data.push([]);
                }

                if(item.id == 'vulInfo'){ //漏洞信息 
                    data.push(['漏洞名称', '漏洞风险', '漏洞状态','漏洞类型','漏洞编号','披露日期','漏洞描述','漏洞结果','修复建议','影响范围','漏洞位置','参考链接','playload','请求报文','响应报文']);
                    if(vulInfo){
                        vulInfo.forEach(vul=>{
                            data.push([
                            vul.vulName,
                            vul.risk,
                            vul.vulStatus,
                            vul.type,
                            vul.cve,
                            vul.publishDate,
                            vul.describe,
                            vul.res,
                            vul.fix,
                            vul.affectRange,
                            vul.location,
                            vul.link,
                            vul.verMsg.payload,
                            vul.verMsg.request,
                            vul.verMsg.response,
                        ]);
                        })
                        
                    }
                    data.push([]);
                }

            })  

            const csv = data.join('\n');
 

            const blob = new Blob([`\ufeff${csv}`], { type: 'text/csv;charset=utf-8;' }); 
            saveAs(blob, this.reportcontent.reportCover.title+'.csv');
        },

    }
}
</script>