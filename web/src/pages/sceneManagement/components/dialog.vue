<template>
    <div @click="showtable()">
        <!-- 新建代理 -->
        <el-dialog :title="title" :visible.sync="dialogaddFormVisible" :before-close="cancelform" width="1184px" id="dialog"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false" class="newUserDialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">保存</el-button>
                <el-button size="small" @click="cancelform">关闭</el-button>
            </div>
            <div style="padding:24px">
<!-- 口令猜测 ----------------------------------------------------------------------------------------------->
                <el-form :model="commandForm" label-width="0" status-icon ref="commandFormRef" :rules="commandFormRules" v-if="type === 1">
                    <div> 
                        <el-form-item label="" style="display: inline-block; width: 29%">
                            <label class="dialog_item_label">猜测次数</label>
                            <el-select v-model="commandForm.guess_number" size="small" placeholder="请选择" style="width: 160px">
                                <el-option v-for="(item, index) in guessCount_list" :key="index" :label="item" :value="item"></el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item label=""  style="display: inline-block;width: 29%;  ">
                            <label class="dialog_item_label">猜测时间</label>
                            <el-select v-model="commandForm.guess_timeout" size="small" placeholder="请选择" style="width: 160px">
                                <el-option v-for="(item, index) in guessTime_list" :key="index" :label="item.label" :value="item.value"></el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item prop="pass_dict" label=""  style="display: inline-block; width: 29%;  ">
                            <label class="dialog_item_label">猜测速率</label>
                            <el-select v-model="commandForm.guess_rate" size="small" placeholder="请选择" style="width: 160px">
                                <el-option v-for="(item, index) in guessRate_list" :key="index" :label="item.label" :value="item.value"></el-option>
                            </el-select>
                        </el-form-item>
                    </div>
                    <el-table ref="multipleTable"  :data="wordguesstable" style="width: 100%; padding: 24px"  
                        class="context_box_bg" @selection-change="handleSelectionChange">
                        <el-table-column type="selection" width="55">
                        </el-table-column>
                        <el-table-column prop="service" label="服务">
                        </el-table-column>
                        <el-table-column prop="user_dict_list" label="账号字典">
                            <template slot-scope="scope">
                                <el-select v-model="scope.row.user_checked" placeholder="请选择" size="small" style="width: 144px"
                                    @change="handleChangeUserList($event, scope.row)">
                                    <el-option v-for="(item, index) in scope.row.user_dict_list" :key="index" :label="item" :value="item">
                                    </el-option>
                                </el-select>
                            </template>
                        </el-table-column>
                        <el-table-column prop="pass_dict_list" label="密码字典" show-overflow-tooltip>
                            <template slot-scope="scope">
                                <el-select v-model="scope.row.pass_checked" placeholder="请选择" size="small" style="width: 144px"
                                    @change="handleChangePassList($event, scope.row)">
                                    <el-option v-for="(item, index) in scope.row.pass_dict_list" :key="index" :label="item" :value="item">
                                    </el-option>
                                </el-select>
                            </template>
                        </el-table-column>
                    </el-table> 
                </el-form>
 <!-- 新增凭证 ------------------------------------------------------------------------------------------------------->
                <el-form :model="newVoucherForm" label-width="0" status-icon ref="newVoucherFormRef" :rules="newVoucherFormRules" v-if="type === 2">
                    <el-form-item label="" prop="testTarget"  >
                        <label class="dialog_item_label">测试目标</label>
                        <el-input v-model="newVoucherForm.testTarget" size="small" style="width:320px" placeholder="请输入测试目标"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label="" prop="agreement">
                        <label class="dialog_item_label">协议</label>
                        <el-select v-model="newVoucherForm.agreement" size="small" placeholder="请选择" style="width: 320px">
                            <el-option v-for="(item, index) in agreement_list" :key="index" :label="item.name" :value="item.id"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="" prop="voucherType">
                        <label class="dialog_item_label">凭证类型</label>
                        <el-select v-model="newVoucherForm.voucherType" size="small" placeholder="请选择" style="width: 320px">
                            <el-option v-for="(item, index) in voucherType_list" :key="index" :label="item.name" :value="item.id"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="" prop="cookie" >
                        <label class="dialog_item_label">Cookie</label>
                        <el-input v-model="newVoucherForm.cookie" size="small" style="width:320px" autocomplete="off" type="textarea"
                            placeholder="请输入凭证类型"></el-input>
                    </el-form-item>
                </el-form>


<!-- 子域名字典 ------------------------------------------------------------------------------------------->
                <el-form :model="subdomainCollectionForm" label-width="0" status-icon ref="subdomainCollectionFormRef" :rules="subdomainCollectionFormRules" v-if="type === 3">
                    <!-- <el-form-item
                        prop = 'subdomaincollect'
                        label=" " class="lastline"> 
                        <label class="dialog_item_label bottomerror">子域名收集</label>
                        <el-checkbox-group v-model="subdomainCollectionForm.subdomaincollect" style="display: inline-block;">
                            <el-checkbox label="子域名搜索"></el-checkbox>
                            <el-checkbox label="DNS查询"></el-checkbox>
                        </el-checkbox-group>
                    </el-form-item>  -->
                    <el-form-item label="" prop="subdomain_dict">
                        <label class="dialog_item_label">子域名字典</label>
                        <el-select v-model="subdomainCollectionForm.subdomain_dict" size="small" placeholder="请选择" style="width: 320px">
                            <el-option v-for="(item, index) in subdomainDictionary_list" :key="index" :label="item" :value="item"></el-option>
                        </el-select>
                    </el-form-item>
                </el-form>


<!-- 端口扫描 -------------------------------------------------------------------------------------------->
                <el-form :model="portScanForm" label-width="0" status-icon ref="portScanFormRef" :rules="portScanFormRules" v-if="type === 4">
                    <el-form-item
                        prop ='tcpscan'
                        label=" " > 
                        <label class="dialog_item_label">TCP扫描</label>
                        <el-radio-group v-model="portScanForm.tcp_scan_type">
                            <el-radio :label="1" value="1">TCP-Connect</el-radio>
                            <el-radio :label="2" value="2">TCP SYC</el-radio>
                            <el-radio :label="3" value="3">TCP FIN</el-radio>
                        </el-radio-group>
                    </el-form-item>
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label">端口范围</label>
                        <el-select v-model="portScanForm.port_scan_type" size="small" placeholder="请选择" style="width: 320px" @change="changePortSacn">
                            <el-option v-for="(item, index) in portRange_list" :key="index" :label="item.name" :value="item.value"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="" prop="scan_port" >
                        <label class="dialog_item_label">扫描端口</label>
                        <el-input v-model="portScanForm.scan_port" size="small" style="width:720px;vertical-align: text-top;" rows="6" autocomplete="off" type="textarea" resize="none" placeholder="请输入扫描端口"></el-input>
                    </el-form-item>
                </el-form>



<!-- 动态爬虫-------------------------------------------------------------------------->
                <el-form class="dynamicCrawlerForm" :model="dynamicCrawlerForm" label-width="0" status-icon ref="dynamicCrawlerFormRef" :rules="dynamicCrawlerFormRules" v-if="type === 5">
                    <el-form-item class="formItem1" label="">
                        <label class="dialog_item_label"  style="display: inline-block;height:16px;margin-top: 11px;width:90px;">参数设置</label>
                        <div class="infobox" style="display:inline-block;">
                            <div class="threelinebox">
                                <el-form-item
                                    prop = 'depth'
                                    label=" "  style="display: inline-block;margin-right:40px"> 
                                    <label class="dialog_item_label topline">爬取深度</label>
                                        <el-select v-model="dynamicCrawlerForm.crawler_depth"  style=" width:180px;" clearable placeholder="爬取深度"  size="small" ref="vulSelect">  
                                            <el-option
                                                v-for="(item,i) in depthlist"
                                                :key="i"
                                                :label="item"
                                                :value="item"> 
                                            </el-option>
                                        </el-select> 
                                </el-form-item>
                                <!-- <el-form-item
                                    prop = 'width'
                                    label=" " style="display: inline-block;margin-right:40px"> 
                                    <label class="dialog_item_label topline">爬取宽度</label>
                                    <el-select v-model="dynamicCrawlerForm.width"  style=" width:160px;" clearable placeholder="爬取宽度"  size="small" ref="vulSelect">  
                                        <el-option
                                            v-for="(item,i) in widthlist"
                                            :key="i"
                                            :label="item.label"
                                            :value="item.value"> 
                                        </el-option>
                                    </el-select>  
                                </el-form-item> -->
                                <el-form-item prop="link_sum" label=" " style="display: inline-block;margin-right:40px" label-width="">
                                    <label class="dialog_item_label topline">链接总数</label>
                                    <el-select v-model="dynamicCrawlerForm.total_link"    size="small" placeholder="请选择"  style=" width:180px;" >
                                        <el-option v-for="(item, index) in link_sum" :key="index" :label="item" :value="item"></el-option>
                                    </el-select>
                                </el-form-item>
                                <el-form-item
                                    prop = 'scope'
                                    label=" " style="display: inline-block;margin-right:40px"> 
                                    <label class="dialog_item_label topline">爬取范围</label>
                                    <el-select v-model="dynamicCrawlerForm.crawl_range"  style=" width:180px;" clearable placeholder="爬取范围"  size="small" ref="vulSelect">  
                                        <el-option
                                            v-for="(item,i) in scopelist"
                                            :key="i"
                                            :label="item[1]"
                                            :value="item[0]"> 
                                        </el-option>
                                    </el-select>  
                                </el-form-item>
                            </div>
                            <div class="threelinebox">
                                <el-form-item
                                    prop = 'timeout'
                                    label=" "  style="display: inline-block;margin-right:40px"> 
                                    <label class="dialog_item_label topline">单链接超时</label>
                                        <el-select v-model="dynamicCrawlerForm.single_link_timeout"  style=" width:180px;" clearable placeholder="单链接超时"  size="small" ref="vulSelect">  
                                            <el-option
                                                v-for="(item,i) in timeoutlist"
                                                :key="i"
                                                :label="item"
                                                :value="item"> 
                                            </el-option>
                                        </el-select> 
                                </el-form-item>
                                <el-form-item
                                    prop = 'speed'
                                    label=" " style="display: inline-block;margin-right:40px"> 
                                    <label class="dialog_item_label topline">爬取速率</label>
                                    <el-select v-model="dynamicCrawlerForm.crawl_speed"  style=" width:180px;" clearable placeholder="爬取速率"  size="small" ref="vulSelect">  
                                        <el-option
                                            v-for="(item,i) in speedlist"
                                            :key="i"
                                            :label="item[1]"
                                            :value="item[0]"> 
                                        </el-option>
                                    </el-select>  
                                </el-form-item>
                                <el-form-item
                                    prop = 'url'
                                    label=" " style="display: inline-block;margin-right:40px"> 
                                    <label class="dialog_item_label topline">URL去重</label>
                                    <el-select v-model="dynamicCrawlerForm.url_removal"  style=" width:180px;" clearable placeholder="URL去重"  size="small" ref="vulSelect">  
                                    
                                        <el-option label="重度去重" value="1"></el-option>
                                        <el-option label="普通去重" value="2"></el-option>
                                        <el-option label="轻度去重" value="3"></el-option>
                                    </el-select>  
                                </el-form-item>
                                <el-form-item label="" style="display: inline-block;margin-right:40px" label-width="">
                                    <label class="dialog_item_label topline">爬取时长</label>
                                    <el-select v-model="dynamicCrawlerForm.url_time"  size="small" placeholder="请选择" style="width:180px"> 
                                        <el-option label="5分钟" value="5"></el-option>
                                        <el-option label="10分钟" value="10"></el-option>
                                        <el-option label="30分钟" value="30"></el-option>
                                        <el-option label="60分钟" value="60"></el-option>
                                        <el-option label="180分钟" value="180"></el-option>
                                        <el-option label="不限" value="0"></el-option>
                                    </el-select>
                                </el-form-item>
                            </div> 
                           <el-form-item
                            prop = 'analysetype'
                            label=" " style="display:block;"> 
                            <label class="dialog_item_label">分析文件类型</label>
                                <el-input
                                    class="infoinput"
                                    style=""
                                    v-model="dynamicCrawlerForm.analyze_file_type"
                                    size="small"
                                    placeholder=""
                                    maxlength="50"  
                                ></el-input>
                            </el-form-item> 
                            <el-form-item
                            prop = 'blackkey'
                            label=" " style="display:block;"> 
                            <label class="dialog_item_label">关键字黑名单</label>
                                <el-input
                                    class="infoinput"
                                    style=""
                                    v-model="dynamicCrawlerForm.black_key"
                                    size="small"
                                    placeholder=""
                                    maxlength="50"  
                                ></el-input>
                            </el-form-item> 
                            <el-form-item
                            prop = 'blackurl'
                            label=" " style="display:block;"> 
                            <label class="dialog_item_label">URL黑名单</label>
                                <el-input
                                    class="infoinput"
                                    style=""
                                    v-model="dynamicCrawlerForm.black_url"
                                    size="small"
                                    placeholder=""
                                    maxlength="50"  
                                ></el-input>
                            </el-form-item> 
                            <el-form-item
                                prop = 'whiteurl'
                                label=" " style="display:block;"> 
                                <label class="dialog_item_label">URL白名单</label>
                                    <el-input
                                        class="infoinput"
                                        style=""
                                        v-model="dynamicCrawlerForm.white_url"
                                        size="small"
                                        placeholder=""
                                        maxlength="50"  
                                    ></el-input>
    
                            </el-form-item> 
                            <el-form-item prop = ' '
                                label=" " style="display:block;"> 
                                <el-radio-group v-model="dynamicCrawlerForm.radio" @change="changeRadio" class="radiobox">
                                    <el-radio :label="2">默认header</el-radio>
                                    <el-radio :label="3"> 自定义header</el-radio>
                                </el-radio-group>
                            </el-form-item>
                            <!-- crawler_header -->
                            <div style="margin-bottom: 8px; width: 670px; margin-left: 10px;">
                                <el-input v-if="dynamicCrawlerForm.radio == 2" v-model="dynamicCrawlerForm.crawler_header" disabled
                                    size="small"></el-input>
                                <el-input v-if="dynamicCrawlerForm.radio == 3" type="textarea" rows="5" resize="none"
                                    v-model="dynamicCrawlerForm.crawler_header1" size="small"></el-input>
                            </div>
                            <el-form-item prop='blackurl' label=" " style="display:block;">
                                <label class="dialog_item_label" style="width:160px">HTTP弱口令爆破</label>
                                <el-switch v-model="dynamicCrawlerForm.http_weak_password_blasting.is_open" class="elSwitch"  >
                                </el-switch>
                            </el-form-item>
                            <div>   
                                <el-form-item prop='blackurl' label=" " style="display:inline-block; ">
                                    <label class="dialog_item_label">账号字典</label>
                                    <el-select v-model="dynamicCrawlerForm.http_weak_password_blasting.user_dict" placeholder="请选择" size="small" style="width: 144px" >
                                        <el-option v-for="(item, index) in user_dict_data" :key="index" :label="item" :value="item">
                                        </el-option>
                                    </el-select>
                                </el-form-item>
                                <el-form-item prop='blackurl' label=" " style="display:inline-block; ">
                                    <label class="dialog_item_label">口令字典</label>
                                    <el-select v-model="dynamicCrawlerForm.http_weak_password_blasting.password_dict" placeholder="请选择" size="small" style="width: 144px">
                                        <el-option v-for="(item, index) in pass_dict_data" :key="index" :label="item" :value="item">
                                        </el-option>
                                    </el-select>
                                </el-form-item>
                            </div>
                            <el-form-item prop="" label=""> 
                                <label class="dialog_item_label">验证码识别</label>
                                <el-radio-group v-model="dynamicCrawlerForm.http_weak_password_blasting.verification_code_detect"  >
                                    <el-radio :label="1" value="1">数字+字母</el-radio>
                                    <el-radio :label="2" value="2">简单算术</el-radio> 
                                </el-radio-group>
                            </el-form-item>
                            
                         </div>   
                    </el-form-item>
                </el-form>

<!-- web路径爆破 -------------------------------------------------------------------------->
                <el-form :model="webPathForm" label-width="0" status-icon ref="webPathFormRef" :rules="webPathFormRules" v-if="type === 6">
                    <el-form-item class="formItem1" label="" >
                        <!-- <label class="dialog_item_label"  style="display: inline-block;height:16px;margin-top: 11px">参数设置</label> -->
                            <div class="infobox2">
                                <el-form-item label="" prop="guess_rate">
                                    <label class="dialog_item_label">猜测速率</label>
                                    <el-select v-model="webPathForm.guess_rate" size="small" placeholder="请选择" style="width: 320px">
                                        <el-option v-for="(item, index) in guessRate_list" :key="index" :label="item.label" :value="item.value"></el-option>
                                    </el-select>
                                </el-form-item>
                                <el-form-item   label="" prop="scan_dict" >
                                    <label class="dialog_item_label">路径字典</label>
                                    <el-select v-model="webPathForm.scan_dict" filterable placeholder="请选择" style="width: 720px" multiple >
                                        <el-option v-for="(item,index) in pathDict_list" :key="index" :label="item.label" :value="item.value">
                                        </el-option>
                                    </el-select>
                                </el-form-item>
                                <el-form-item label="" prop="title_black" >
                                    <label class="dialog_item_label">排除标题黑名单</label>
                                    <el-input v-model="webPathForm.title_black" size="small" style="width:720px" autocomplete="off"
                                        placeholder="请输入排除标题黑名单"></el-input>
                                </el-form-item>
                            </div>
                     </el-form-item>
                </el-form>
<!-- 横向移动 -->
                <el-form :model="transverseForm" label-width="0" ref="transverseFormRef" v-if="type === 7">
                    <el-form-item label="" prop="guessRate">
                        <label class="dialog_item_label">后渗透数量</label>
                        <el-input v-model="transverseForm.number"  size="small" style="width: 320px" placeholder="请输入后渗透数量"  ></el-input>
                    </el-form-item>
                    <el-form-item>
                        <label class="dialog_item_label">渗透条数</label>
                        <el-select v-model="transverseForm.ttl" filterable placeholder="请选择" style="width: 320px" >
                            <el-option v-for="(item, index) in infiltrationHopslist" :key="index" :label="item" :value="item"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="" prop="excludedTitleBlacklist">
                        <label class="dialog_item_label">横向范围</label>
                        <!-- :autosize="{ minRows: 3, maxRows: 8}"  -->
                        <el-input type="textarea" placeholder="请输入内容" rows="5"
                            v-model="transverseForm.scope"  
                            style="width: 720px; margin-bottom: 10px;vertical-align: text-top;">
                        </el-input>
                    </el-form-item>
                </el-form>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { proxy } from '@/api/system.js';
export default {
    name:"myDialog",
    components:{
    }, 
    data(){
        var validatePwd =  (rule, value, callback) => {
            const reg = /^\S+$/; 
            if (reg.test(value)) {
                callback();
            } else {
                return callback(new Error('格式不正确(不能包含空格)'));
            }
        }; 
        return{
            title: '',
            is_Update:false, 
            commandForm: { // 口令猜测 
                guess_number: '', // 猜测次数
                guess_timeout: '', // 猜测时间
                guess_rate: '', // 猜测速率
                infos:[],
            },
            account_list: [], // 账号字典下拉列表
            password_list: [], // 密码字典下拉列表
            guessCount_list: ["0", "1", "100", "1000", "10000"], // 猜测次数下拉列表
            guessTime_list: [
                { value: 60, label: "1min" },
                { value: 180, label: "3min" },
                { value: 300, label: "5min" },
                { value: 600, label: "10min" },
                { value: 1800, label: "30min" },
                { value: 3600, label: "60min" },
                { value: 0, label: "不限" },
            ], // 猜测时间下拉列表
            guessRate_list: [
                { value: 1, label: "高速" },
                { value: 2, label: "中速" },
                { value: 3, label: "低速" },
            ], // 猜测速率下拉列表

            newVoucherForm: { // 新增凭证
                testTarget: '', // 测试目标
                agreement: '', // 协议
                voucherType: '', // 凭证类型
                cookie: '' // Cookie
            },
            agreement_list: [], // 协议下拉列表
            voucherType_list: [], // 凭证类型下拉列表

            subdomainCollectionForm: { // 子域名收集 
                is_open:true,
                subdomain_dict: '' // 子域名字典
            },
            subdomainDictionary_list: [], // 子域名字典下拉列表
            
            portScanForm: { // 端口扫描  
                port_scan_type: 4, // 端口范围
                scan_port: '', // 扫描端口
                tcp_scan_type: 1,
            },
            portRange_list: [
                {
                    name:'TOP10端口',
                    value:1
                },
                {
                    name: 'TOP20端口',
                    value: 2
                },
                {
                    name: 'TOP50端口',
                    value: 3
                },
                {
                    name: 'TOP100端口',
                    value: 4
                },
                {
                    name: 'TOP500端口',
                    value: 5
                },
                {
                    name: 'TOP1000端口',
                    value: 6
                },
                {
                    name: '全部端口',
                    value: 7
                },
                {
                    name: '自定义端口',
                    value: 8
                },
            ], // 端口范围下拉列表
            transverseForm:{
                number:'', //数量
                ttl:'',
                scope:'', //范围
            },
            infiltrationHopslist: ['0', '1', '2', '3', '4', '5', '6', '7', '8'],
            dynamicCrawlerForm: {// 动态爬虫 
                width:'',// 爬取宽度  
                crawler_depth: '', //深度
                total_link: '', //链接总数
                crawl_range: '', //爬取范围
                single_link_timeout:'', //单链接超时
                crawl_speed: '', //爬取速率 
                url_removal:2, //url去重,
                url_time:'0',
                analyze_file_type: '', //分析文件类型
                black_key: '', //关键字白名单
                black_url: '', //URL白名单 
                white_url: '',
                crawl_time: '',
                radio:2,
                crawler_header:
                    "{'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64; rv:68.0) Gecko/20100101 Firefox/68.0', 'Accept': '*/*', 'Accept-Language': 'zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2', 'Accept-Encoding': 'gzip, deflate', 'Connection': 'close'}",
                crawler_header1: "",
                http_weak_password_blasting:{
                    is_open:false,
                    password_dict:0,
                    user_dict:'',
                    verification_code_detect:'',
                }
            },
            user_dict_data:[ 'http', 'common'], //账号字典
            pass_dict_data: ['http', 'common'], //密码字典
            user_dict_list:[], 
            depthlist: ["1", "2", "3", "4", "5", "不限"], //爬取深度
            link_sum: [
                "100",
                "200",
                "500",
                "1000",
                "5000",
                "10000",
                "100000",
                "不限",
            ], //链接总数
            scopelist: [
                [1, " 爬取全域名"],
                [2, " 爬取目标URL与子目录"],
                [3, "爬取目标url"],
            ], //爬取范围
            timeoutlist: ["2", "5", "10", "20", "30", "60", "120"], //单链接超时
            speedlist: [
                [1, "高速"],
                [2, "中速"],
                [3, "低速"],
            ], 
            urllist: [], 
            usetimelist: [], 

            webPathForm: { // web路径猜测
                guess_rate: '', // 猜测速率
                scan_dict: [], // 路径字典
                title_black: '', // 排除标题黑名单
            },
            webGuessRate_list: [], // 猜测速率下拉列表
            pathDict_list: [
                 {
                    label:'PHP字典',
                    value:'php'
                    },
                    {
                    label:'ASP字典',
                    value:'asp'
                    },
                    {
                    label:'ASPX字典',
                    value:'aspx'
                    },
                    {
                    label:'JSP字典',
                    value:'jsp'
                    },
            ], // 路径字典下拉列表

            agencyform:{
                id:'',
                type:'',
                name:'',
            	agencyport:'',
            	agencyip:'',
            	password:'',
                username:'',
            	agreement:'',
                
            },
            basicinfo:[],
            commandFormRules: {
                account: [
                    { required: true, message: '请选择账号字典', trigger: 'change' },
                ],
                password: [
                    { required: true, message: '请选择密码字典', trigger: 'change' },
                ],
                replenishAccount: [
                    { required: true, message: '补充账号不能为空', trigger: 'blur' },
                ],
                replenishCommon: [
                    { required: true, message: '补充口令不能为空', trigger: 'blur' },
                ],
                guessCount: [
                    { required: true, message: '请选择猜测次数', trigger: 'change' },
                ],
                guessTime: [
                    { required: true, message: '请选择猜测时间', trigger: 'change' },
                ],
                guessRate: [
                    { required: true, message: '请选择猜测速率', trigger: 'change' },
                ],
            },
            newVoucherFormRules: {
                testTarget: [
                    { required: true, message: '测试目标不能为空', trigger: 'blur' },
                ],
                agreement: [
                    { required: true, message: '请选择协议', trigger: 'change' },
                ],
                voucherType: [
                    { required: true, message: '请选择凭证类型', trigger: 'change' },
                ],
                cookie: [
                    { required: true, message: 'Cookie不能为空', trigger: 'blur' },
                ]
            },
            subdomainCollectionFormRules: {
                subdomain_dict: [
                    { required: true, message: '请选择子域名字典', trigger: 'change' },
                ]
            },
            portScanFormRules: {
                // portRange: [
                //     { required: true, message: '请选择端口范围', trigger: 'change' },
                // ],
                // scanPort: [
                //     { required: true, message: '扫描端口不能为空', trigger: 'blur' },
                // ]
            },
            dynamicCrawlerFormRules: {},
            webPathFormRules: {
                guess_rate: [
                    { required: true, message: '请选择猜测速率', trigger: 'change' },
                ],
                scan_dict: [
                    { required: true, message: '请选择路径字典', trigger: 'change' },
                ],
                title_black: [
                    { required: true, message: '排除标题黑名单不能为空', trigger: 'blur' },
                ]
            },
            rules1: {
                name: [
                    { required: true, message: '代理名称不能为空', trigger: 'blur' },
                ],
                agreement: [
                    { required: true, message: '请选择代理协议', trigger: 'change' },
                ],
                agencyip: [
                    { required: true, message: '代理IP不能为空', trigger: 'blur' }, 
                ],
                agencyport: [
                    { required: true, message: '代理端口不能为空', trigger: 'blur' },
                    { max: 50, message: '代理端口不能超过50', trigger: 'blur' },
                ], 
                username: [
                    { required: true, message: '账号不能为空', trigger: 'blur' }, 
                ],
                password: [
                    { required: true, message: '密码不能为空', trigger: 'blur' }, 
                    { validator: validatePwd, trigger: 'blur' }
                ],
            },
            webPassWordCheckedData:[],
            wordguesstable:[],
            PasswordData:[],

        }
        
    },
    // props: ['dialogaddFormVisible','type','portconfig','word_guess'],
    props:{
        dialogaddFormVisible:{},
        type:{},
        portconfig:{},
        word_guess:{},
        web_crawler:{},
        web_path:{},
        post_penetration:{},
        subdomain:{}
    },
    created () {
        console.log(this.type)
        switch (this.type) {
            case 1:
                this.title = '口令猜测';
                this.commandForm = this.word_guess; 
                this.PasswordData = this.word_guess.infos;
                this.$nextTick(() => {
                    this.fnWebGetData(); 
                    this.showtable();
                });  
                break
            case 2:
                this.title = '新增凭证'
                break
            case 3:
                this.title = '子域名收集'
                this.getSubdomainList();
                this.subdomainCollectionForm = this.subdomain
                break
            case 4:
                this.title = '端口扫描';
                this.portScanForm = this.portconfig;
                
                // this.$nextTick(() => {
                //     this.changePortSacn();
                // });
                break
            case 5:
                this.title = '动态爬虫';
                this.dynamicCrawlerForm = this.web_crawler;
                break
            case 6:
                this.title = 'web路径爆破';
                this.webPathForm = this.web_path
                break
            case 7:
                this.title = '横向移动';
                this.transverseForm = this.post_penetration 
                break
        } 
        
    }, 
    mounted () {  
    },
    methods: {
        showtable(){ },
        changePortSacn(){ //端口范围选择事件 
            switch (this.portScanForm.port_scan_type) {
                case 1: // top10
                    this.portScanForm.scan_port='21,22,23,80,443,445,3306,8000,8080,8088';
                    break;
                case 2: //top20
                    this.portScanForm.scan_port ='21,22,23,80,443,445,3306,7000-7002,8000-8003,8080-8083,8088,9200';
                    break;
                case 3: //top50
                    this.portScanForm.scan_port ='21,22,23,80,88,106,110,111,113,119,135,139,143,144,179,199,389,427,1521,1630,1158,443,445,888,777,999,1070,1080,1090,3306,7000-7003,8000-8003,8008,8080-8083,8088,9000-9002,8090,9200,9300';
                    break;
                case 4: // top100
                    this.portScanForm.scan_port ='7,5555,9,13,21,22,23,25,26,37,53,79,80,81,88,106,110,111,113,119,135,139,143,144,179,199,389,427,443,444,445,465,513,514,515,543,544,548,554,587,631,646,873,888,990,993,995,1025,1026,1027,1028,1029,1080,1110,1433,1443,1720,1723,1755,1900,2000,2001,2049,2121,2181,2717,3000,3128,3306,3389,3986,4899,5000,5009,5051,5060,5101,5190,5357,5432,5631,5666,5800,5900,6000,6001,6646,7000,7001,7002,7003,7004,7005,7070,8000,8008,8009,8080,8081,8443,8888,9100,9999,10000,11211,32768,49152,49153,49154,49155,49156,49157,8088,9090,8090,8001,82,9080,8082,8089,9000,8002,89,8083,8200,90,8086,801,8011,8085,9001,9200,8100,8012,85,8084,8070,8091,8003,99,7777,8010,8028,8087,83,808,38888,8181,800,18080,8099,8899,86,8360,8300,8800,8180,3505,9002,8053,1000,7080,8989,28017,9060,8006,41516,880,8484,6677,8016,84,7200,9085,5555,8280,1980,8161,9091,7890,8060,6080,8880,8020,889,8881,9081,7007,8004,38501,1010,17,19,255,1024,1030,1041,1048,1049,1053,1054,1056,1064,1065,1801,2103,2107,2967,3001,3703,5001,5050,6004,8031,10010,10250,10255,6888,87,91,92,98,1081,1082,1118,1888,2008,2020,2100,2375,3008,6648,6868,7008,7071,7074,7078,7088,7680,7687,7688,8018,8030,8038,8042,8044,8046,8048,8069,8092,8093,8094,8095,8096,8097,8098,8101,8108,8118,8172,8222,8244,8258,8288,8448,8834,8838,8848,8858,8868,8879,8983,9008,9010,9043,9082,9083,9084,9086,9087,9088,9089,9092,9093,9094,9095,9096,9097,9098,9099,9443,9448,9800,9981,9986,9988,9998,10001,10002,10004,10008,12018,12443,14000,16080,18000,18001,18002,18004,18008,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018';
                    break;
                case 5: // top500
                    this.portScanForm.scan_port ='7,8,9,13,17,19,20,21,22,23,25,26,37,53,60,65,66,70,77,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,103,106,110,111,113,114,119,122,132,133,135,139,143,144,171,179,180,188,199,200,206,208,211,235,255,268,280,299,302,308,321,381,389,403,421,423,427,442,443,444,445,447,465,511,513,514,515,517,522,543,544,548,554,587,591,610,631,646,666,688,701,770,778,800,801,802,803,804,805,806,808,809,811,812,866,873,877,880,888,889,925,955,983,990,993,995,999,1000,1001,1005,1010,1024,1025,1026,1027,1028,1029,1030,1039,1041,1042,1048,1049,1053,1054,1056,1064,1065,1080,1081,1082,1085,1088,1100,1107,1108,1110,1118,1122,1123,1128,1158,1180,1182,1200,1212,1213,1234,1300,1301,1313,1356,1433,1443,1500,1550,1666,1680,1700,1720,1722,1723,1755,1790,1800,1801,1818,1863,1888,1900,1933,1949,1979,1980,1982,2000,2001,2005,2006,2007,2008,2009,2010,2011,2012,2013,2014,2015,2020,2046,2049,2051,2060,2070,2080,2093,2100,2103,2107,2110,2121,2125,2181,2222,2301,2348,2375,2382,2480,2490,2517,2521,2585,2717,2808,2886,2901,2967,3000,3001,3008,3012,3013,3030,3050,3080,3128,3216,3220,3306,3312,3333,3380,3389,3443,3456,3465,3503,3505,3535,3580,3588,3600,3606,3668,3690,3703,3721,3880,3938,3986,4000,4001,4016,4040,4300,4389,4430,4433,4440,5000,5001,5002,5003,5009,5013,5050,5600,5601,5631,5632,5644,5655,5656,5666,5678,5800,5811,5881,5887,5888,5898,5900,5902,5966,6000,6001,6002,6003,6004,6006,6010,6060,6080,6088,6090,6101,6118,6170,6180,6198,6226,6259,6379,6388,6886,6888,6889,6890,6920,6969,6988,7000,7001,7002,7003,7004,7005,7006,7007,7008,7009,7010,7011,7012,7017,7018,7020,7021,7022,7028,7031,7041,7048,7050,7060,7070,7071,7074,7078,7080,7081,7084,7086,7088,7094,7100,7101,7102,7108,7111,7117,7123,7129,7171,7180,7200,7201,7202,7215,7272,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8016,8018,8019,8020,8021,8022,8023,8024,8025,8026,8027,8028,8029,8030,8031,8032,8033,8035,8036,8037,8038,8039,8040,8041,8042,8043,8044,8045,8046,8048,8050,8051,8053,8055,8056,8057,8058,8060,8061,8062,8064,8065,8066,8069,8070,8071,8073,8077,8078,8079,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8102,8103,8104,8172,8180,8181,8182,8183,8184,8186,8188,8288,8300,8308,8322,8333,8341,8343,8360,8380,8580,8582,8585,8600,8601,8610,8649,8660,9200,9201';
                    break;
                case 6: // top1000
                    this.portScanForm.scan_port ='7,8,9,13,17,19,20,21,22,23,25,26,37,53,60,65,66,70,77,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,103,106,110,111,113,114,119,122,132,133,135,139,143,144,171,179,180,188,199,200,206,208,211,235,255,268,280,299,302,308,321,381,389,403,421,423,427,442,443,444,445,447,465,511,513,514,515,517,522,543,544,548,554,587,591,610,631,646,666,688,701,770,778,800,801,802,803,804,805,806,808,809,811,812,866,873,877,880,888,889,925,955,983,990,993,995,999,1000,1001,1005,1010,1024,1025,1026,1027,1028,1029,1030,1039,1041,1042,1048,1049,1053,1054,1056,1064,1065,1080,1081,1082,1085,1088,1100,1107,1108,1110,1118,1122,1123,1128,1158,1180,1182,1200,1212,1213,1234,1300,1301,1313,1356,1433,1443,1500,1550,1666,1680,1700,1720,1722,1723,1755,1790,1800,1801,1818,1863,1888,1900,1933,1949,1979,1980,1982,2000,2001,2005,2006,2007,2008,2009,2010,2011,2012,2013,2014,2015,2020,2046,2049,2051,2060,2070,2080,2093,2100,2103,2107,2110,2121,2125,2181,2222,2301,2348,2375,2382,2480,2490,2517,2521,2585,2717,2808,2886,2901,2967,3000,3001,3008,3012,3013,3030,3050,3080,3128,3216,3220,3306,3312,3333,3380,3389,3443,3456,3465,3503,3505,3535,3580,3588,3600,3606,3668,3690,3703,3721,3880,3938,3986,4000,4001,4016,4040,4300,4389,4430,4433,4440,4443,4567,4848,4850,4899,5000,5001,5002,5003,5009,5013,5050,5051,5060,5080,5081,5098,5100,5101,5155,5156,5190,5200,5201,5203,5233,5255,5256,5280,5357,5432,5544,5552,5555,5561,5600,5601,5631,5632,5644,5655,5656,5666,5678,5800,5811,5881,5887,5888,5898,5900,5902,5966,6000,6001,6002,6003,6004,6006,6010,6060,6080,6088,6090,6101,6118,6170,6180,6198,6226,6259,6379,6388,6443,6510,6543,6546,6565,6587,6600,6602,6603,6611,6646,6648,6666,6677,6680,6688,6699,6778,6789,6800,6801,6842,6868,6869,6879,6886,6888,6889,6890,6920,6969,6988,7000,7001,7002,7003,7004,7005,7006,7007,7008,7009,7010,7011,7012,7017,7018,7020,7021,7022,7028,7031,7041,7048,7050,7060,7070,7071,7074,7078,7080,7081,7084,7086,7088,7094,7100,7101,7102,7108,7111,7117,7123,7129,7171,7180,7200,7201,7202,7215,7272,7288,7321,7330,7380,7443,7500,7567,7680,7687,7688,7700,7702,7703,7709,7711,7713,7742,7776,7777,7778,7788,7791,7801,7856,7888,7890,7899,7900,7909,7915,7921,7925,7942,7943,7979,7999,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8016,8018,8019,8020,8021,8022,8023,8024,8025,8026,8027,8028,8029,8030,8031,8032,8033,8035,8036,8037,8038,8039,8040,8041,8042,8043,8044,8045,8046,8048,8050,8051,8053,8055,8056,8057,8058,8060,8061,8062,8064,8065,8066,8069,8070,8071,8073,8077,8078,8079,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8102,8103,8104,8108,8111,8112,8118,8119,8122,8123,8130,8133,8136,8144,8161,8168,8172,8180,8181,8182,8183,8184,8186,8188,8189,8190,8191,8192,8193,8196,8197,8200,8213,8220,8222,8244,8258,8260,8280,8282,8283,8288,8300,8308,8322,8333,8341,8343,8360,8380,8381,8382,8383,8384,8390,8399,8400,8401,8402,8443,8445,8448,8477,8480,8481,8484,8500,8567,8580,8582,8585,8600,8601,8610,8649,8660,8666,8680,8686,8688,8700,8710,8720,8735,8780,8781,8787,8788,8799,8800,8801,8802,8806,8808,8809,8810,8813,8822,8834,8838,8839,8844,8848,8858,8860,8864,8866,8868,8877,8879,8880,8881,8885,8886,8887,8888,8889,8890,8891,8892,8895,8898,8899,8900,8902,8910,8912,8913,8955,8956,8972,8980,8983,8987,8988,8989,8990,8991,8997,8999,9000,9001,9002,9003,9004,9005,9006,9007,9008,9009,9010,9011,9012,9013,9014,9015,9020,9022,9025,9030,9031,9036,9039,9043,9050,9053,9060,9061,9070,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9101,9110,9111,9112,9131,9180,9182,9190,9191,9200,9201,9212,9231,9300,9301,9302,9437,9443,9448,9494,9500,9504,9507,9527,9595,9666,9696,9704,9800,9845,9876,9888,9889,9898,9900,9901,9909,9910,9912,9914,9918,9919,9980,9981,9986,9988,9990,9991,9992,9995,9997,9998,9999,10000,10001,10002,10003,10004,10007,10008,10009,10010,10016,10021,10025,10038,10040,10051,10066,10068,10080,10082,10086,10087,10088,10089,10099,10118,10154,10250,10255,10333,11000,11001,11080,11158,11211,11324,11347,11362,11366,11372,11381,12001,12018,12333,12345,12443,12881,13333,13382,13988,14000,14007,15000,15004,15018,15580,15672,15693,15801,15888,16080,16788,17000,17003,17095,17777,18000,18001,18002,18004,18008,18060,18080,18081,18082,18085,18088,18090,18098,18103,18264,18801,18803,18880,18881,18888,19001,19010,19045,19080,19101,19244,20000,20001,20021,20022,20046,20052,20140,20151,20153,20720,20806,20808,21000,21080,21245,21501,21502,22222,22343,22580,23352,23454,25006,25024,27000,27017,28017,28018,28080,28099,28214,28280,28780,30000,30001,30058,30082,30088,30551,31188,31945,32766,32768,34440,38000,38080,38086,38443,38501,38517,38888,40000,40069,40080,40310,41516,42424,43651,45149,45177,47078,47088,47583,48080,49152,49153,49154,49155,49156,49157,49705,49960,50000,50030,50045,50060,50070,50075,50080,50090,50240,51106,55351,55858,57880,58000,58031,58060,58080,58898,59009,59777,59999,60010,60022,60101,60465,61081,61999,65000,65001,65055,65129,65486,65493,65533,65535';
                    break;
                case 7: // 全部
                    this.portScanForm.scan_port ='0-65535';
                    break;
                case 8: // 自定义
                    this.portScanForm.scan_port ='';
                    break;
            }
        },
        fnWebGetData() {
            //获取口令猜测数据
            this.$ajax
            .get("/task/task/get/user_password/dictionary/", {})
            .then((dt) => {
            let res = dt.data;
            if (res.success) {
                this.wordguesstable = res.result;
                this.wordguesstable.pop(); 
 
                //设置选中状态
                this.$nextTick(() => {
                    let _that = this;
                    let _checkedData = this.PasswordData;
                    let _data = this.wordguesstable;
                    for (let i = 0; i < _data.length; i++) {
                        for (let j = 0; j < _checkedData.length; j++) {
                            if (_data[i].service == _checkedData[j].service) {
                                _data[i].pass_checked = _checkedData[j].pass_dict;
                                _data[i].user_checked = _checkedData[j].user_dict;
                                if (_checkedData[j].pass_dict && _checkedData[j].user_dict) {
                                    _that.$refs.multipleTable.toggleRowSelection(_data[i], true);
                                }
                            }
                        }
                    }
                     
                })
            }
            })
            .catch((err) => {});
        },
        handleSelectionChange(val){ //口令猜测表格多选框
            this.commandForm.infos = val; 
            let _that = this;
            if (val) {
                val.forEach((v, index) => { 
                _that.commandForm.infos[index].user_dict = v.user_checked;
                _that.commandForm.infos[index].pass_dict = v.pass_checked;
                });
            }
 
        },
        handleChangeUserList(val, row) { //口令猜测用户字典
            row.user_checked = val;
            row.user_dict = val;
        },
        handleChangePassList(val, row) { //口令猜测密码字典
            row.pass_checked = val;
            row.pass_dict = val;
        },
         getSubdomainList() { //子域名字典
            this.$ajax
                .get("/task/task/get/subdomain/dictionary/", {
                    params: {},
                })
                .then((dt) => {
                    let res = dt.data;
                    if (res.success) {
                        this.subdomainDictionary_list = res.result.dict_list
                    } else {
                        this.$message.error(res.error)
                    } 
                })
                .catch((res) => { });
        },
        changeRadio () {
            
        },
        submitForm(){
            switch (this.type) {
                case 1: // 口令猜测
                    this.$refs.commandFormRef.validate( async (valid) => {
                        if (valid) {
                            this.$emit('wordguess', this.commandForm);
                            this.$emit('hanldeClose', false)
                        }
                    })
                    break
                case 2: // 新增凭证
                    this.$refs.newVoucherFormRef.validate( async (valid) => {
                        if (valid) {

                        }
                    })
                    break
                case 3: // 子域名收集
                    this.$refs.subdomainCollectionFormRef.validate( async (valid) => {
                        if (valid) {
                             this.$emit('subdomain_config', this.subdomainCollectionForm);
                            this.$emit('hanldeClose', false)
                        }
                    })
                    break
                case 4: // 端口扫描
                    this.$refs.portScanFormRef.validate( async (valid) => {
                        if (valid) {
                            this.$emit('portScanconfig', this.portScanForm);
                            this.$emit('hanldeClose', false)
                        }   
                    })
                    break
                case 5: // 动态爬虫
                    this.$refs.dynamicCrawlerFormRef.validate( async (valid) => {
                        if (valid) {
                            
                            this.$emit('web_crawler_config', this.dynamicCrawlerForm);
                            this.$emit('hanldeClose', false)
                        }
                    })
                    break
                case 6: // web路径猜测
                    this.$refs.webPathFormRef.validate( async (valid) => {
                        if (valid) {
                            this.$emit('web_path_config', this.webPathForm);
                            this.$emit('hanldeClose', false)
                            
                        }
                    })
                    break
                case 7: // 横向移动
                    this.$refs.transverseFormRef.validate( async (valid) => {
                        if (valid) {
                            
                             this.$emit('post_penetration_config',this.transverseForm );
                            this.$emit('hanldeClose', false)
                        }
                    })
                    break
            } 
            
        }, 
        cancelform(){
            this.$emit('hanldeClose', false)
            // this.dialogaddFormVisible = false;
        },
        btnUpdate(){
            this.is_Update=true;
        },
        saveUpdate(){

        },
       
        
       
    },
}
</script>
<style  scoped lang="less">
.el-form-item {
    margin-bottom: 8px;
}
.infobox{
    .infoinput{
        display:block;
        width:670px;
        margin-left:10px;
    }
    .radiobox{
        padding-left: 10px; 
    }
}
.dialog_item_label{
    font-size: 14px;
    border-left: 3px solid #4c7ae3;
    text-align: left;
    padding-left: 10px;
    font-weight: 500;
    display: inline-block;
    line-height: 16px;
    box-sizing: border-box;
    font-size: 13px;
    line-height: 14px;
    font-weight: 600;
    color: rgba(72,72,102,0.87);
}
.infobox .dialog_item_label{
    border-left: none;
    font-weight: 400;
    width: 95px;
    color: rgba(72,72,102,0.64);
}
.infobox2 .dialog_item_label{
    border-left: none;
    font-weight: 600;
    width: 115px;
    color: rgba(72,72,102,0.87);
}
/deep/ .infobox .el-radio__label{
    color: rgba(72,72,102,0.64)!important;  
}
/deep/ .el-form-item__content{
    margin-bottom: 10px;
}
/deep/ .el-dialog__body{
    background-color: rgb(247, 247, 251);
}
/deep/ .el-form-item__label {
    position: absolute;
    left: 127px;
}
 /deep/ .el-form-item__error {
    left: 112px;
}
/deep/ .dynamicCrawlerForm{
    .formItem1 > .el-form-item__content{
        display: flex;
    }  
}
/deep/ .el-radio__label {
    font-size: 14px;
    padding-left: 10px;
    color: rgba(72,72,102,0.64);
}
</style>