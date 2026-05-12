<template>
    <div>
        <div class="main-title" v-if="!group_id && type != 3" > 
            <router-link :underline="false" class="classA" :to="{ path: '/task' }"  >渗透任务
            </router-link> 
            <label class="currentpagetitle" >   
                <span>创建任务</span>
            </label> 
        </div>
        <div class="main-title" v-else-if="group_id && type != 3" >
            <router-link :underline="false" class="classA" :to="{ path: '/taskgroup' }" >任务组
            </router-link>  
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark"  :content="group_name"  placement="bottom">
                    <span>{{group_name}}</span>
                </el-tooltip>
            </label>
        </div>
        <div class="main-title" v-if="type == 3" > 
            <router-link :underline="false" class="classA" :to="{ path: '/assettree' }"  >资产管理
            </router-link> 
            <label class="currentpagetitle" >   
                <span>渗透任务</span>
            </label> 
        </div>
        <BannerBox ref="BannerBox" tips="任务参数通过任务场景进行默认配置，用户可以在创建任务时对默认参数进行调整。" style="margin-bottom: 16px;">
            <el-button type="primary" size="small" @click="submithandle" > 执行任务 </el-button> 
        </BannerBox>
        <div class="createtask_box">
            <el-form :model="taskform" ref="form" :rules="rules" style="height: 100%">
                <div class="basic_config">
                    <el-row :gutter="20">
                        <el-col :span="12">
                            <el-form-item label="" prop="template"  >
                                <label class="dialog_item_label">任务场景</label>
                                <el-select 
                                    :disabled="Boolean($route.query.disabled)"
                                    v-model="taskform.template" 
                                    size="small" placeholder="请选择" class="form_item_width"  
                                    @change="gettemplateconfigbyid(taskform.template)">
                                    <el-option v-for="(item, index) in templatelist" :key="index" :label="item.name"
                                        :value="item.id"></el-option>
                                </el-select> 
                                <div class="template-desc" >{{ template_desc }}</div>
                            </el-form-item>
                            <div   style="position: relative">
                                <label class="dialog_item_label" style="vertical-align: top;  ">测试范围 <i class="is-required" style="float: right">*</i></label>
                                <el-form-item label="" prop="target" label-width="0" class="target_box" 
                                    style="display: inline-block; margin-right: 0">
                                    <el-input type="textarea" :rows="4" 
                                    @input="targetinput" v-model="taskform.target"
                                    autocomplete="off" placeholder="测试范围不能为空" resize="none" class="form_item_width"  
                                        style=" margin-bottom: 10px; "></el-input>
                                    <div >
                                        <el-button type="primary" size="small" style="vertical-align: top; margin-right: 27px"
                                        @click="clickupload()">导入</el-button>
                                        <span style="color: rgba(72, 72, 102, 0.32)">只能上传.txt或.xls或.xlsx格式文件</span>
                                        <input type="file" class="btnUploadID" ref="upload" @change="changeuploaID($event)"
                                            style="display: none" id="input-file-ID" />
                                    </div>
                                </el-form-item> 
                                <el-tooltip placement="right-start">
                                    <div slot="content">
                                    测试范围支持IP、IP段、域名、URL，多个不同目标用“换行”隔开；<br />
                                    示例：<br />“192.168.0.127”、“192.168.0.10-127”、<br />“baidu.com”、“www.baidu.com”、“http://www.baidu.com/aqjc/”、<br />“192.168.0.127:8000”
                                    </div>
                                    <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                                </el-tooltip>
                            </div>
                            <el-form-item label="" prop="taskname"  class="taskNameClass">
                                <label class="dialog_item_label">任务名称<i class="is-required" style="float: right">*</i></label>
                                <el-input :disabled="Boolean($route.query.disabled)" v-model="taskform.taskname"   size="small"
                                class="form_item_width"   placeholder="请输入任务名称" maxlength="50"></el-input>
                            </el-form-item> 
                            <el-form-item label="" prop="execute_type" >
                                <label class="dialog_item_label">任务优先级</label>
                                <el-select
                                :disabled="Boolean($route.query.disabled)"
                                v-model="taskform.weight" size="small" 
                                    placeholder="请选择" class="form_item_width" > 
                                    <el-option v-for="(item,i) in prioritylist" :key="i" :label="item.label" :value="item.value">
                                    </el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item label="" prop="execute_type">
                                <label class="dialog_item_label" style="vertical-align: top;">排除目标</label>
                                <el-input type="textarea" :rows="4" v-model="taskform.excludeTarget" autocomplete="off" placeholder="此处可以填写本次任务中不检测哪些目标" 
                                    resize="none" class="form_item_width" style="margin-bottom: 10px;">
                                </el-input>
                            </el-form-item>
                            <!-- 新增测试模式 -->
                            <el-form-item label="" prop="testMode">
                                <label class="dialog_item_label">测试模式</label>
                                <el-checkbox-group v-model="taskform.testMode" :disabled="Boolean($route.query.disabled)" style="display: inline-block; margin-left: 10px;">
                                    <el-checkbox label="1" style="margin-right: 20px;">
                                        原理验证
                                    </el-checkbox>
                                    <el-checkbox label="2">
                                        版本匹配
                                    </el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12"> 
                            <el-form-item label="" prop="execute_type" >
                                <label class="dialog_item_label">执行方式<i class="is-required" style="float: right">*</i></label>
                                <el-select
                                :disabled="Boolean($route.query.disabled)"
                                v-model="taskform.execute_type" size="small" 
                                    placeholder="请选择" class="form_item_width" >
                                    <!-- 1即时，2定时，3周期，4监控 -->
                                    <el-option v-for="(item,i) in rantypelist" :key="i" :label="item.label" :value="item.value">
                                    </el-option>
                                </el-select>
                                <div class="margin_left_width" v-if="taskform.execute_type == 2">
                                    <label class="execute_type_label">计划时间</label>
                                    <el-date-picker :clearable="false" v-model="taskform.time" type="datetime" size="small"
                                    format="yyyy-MM-dd HH:mm:ss" value-format="yyyy-MM-dd HH:mm:ss" placeholder="选择日期时间">
                                    </el-date-picker>
                                </div>
                                <div class="margin_left_width" v-if="taskform.execute_type == 3" style="margin-bottom: 8px">
                                    <div>
                                    <label class="execute_type_label">计划时间</label>
                                    <el-select v-model="taskform.execute_cycletype" size="small" placeholder="请选择"
                                        style="width: 126px; margin-right: 16px">
                                        <el-option v-for="(item,i) in cyclePlanningType" :key="i" :value="item.value" :label="item.label"></el-option> 
                                    </el-select>
                                    <el-select v-if="taskform.execute_cycletype == '2'"
                                        v-model="taskform.execute_cycletype_day" size="small" placeholder="请选择"
                                        style="width: 126px; margin-right: 16px">
                                        <el-option v-for="(item,i) in cyclePlanningMonthValue" :key="i" :label="item.label" :value="item.value">
                                        </el-option>
                                    </el-select>
                                    <el-select v-if="taskform.execute_cycletype == '1'"
                                        v-model="taskform.execute_cycletype_week" size="small" placeholder="请选择"
                                        style="width: 126px; margin-right: 16px">
                                        <el-option v-for="(item,i) in cyclePlanningWeekValue" :key="i" :value="item.value" :label="item.label"></el-option> 
                                    </el-select>
                                    <el-time-select v-model="taskform.execute_cycletype_starttime" size="small"
                                        :clearable="false" :picker-options="{
                                            start: '00:00',
                                            step: '00:01',
                                            end: '23:59',
                                        }" style="width: 120px !important" placeholder="选择时间">
                                    </el-time-select>
                                    </div>
                                    <div>
                                    <label class="execute_type_label">终止时间</label>
                                    <el-date-picker v-model="taskform.endtime" size="small" format="yyyy-MM-dd HH:mm:ss"
                                        value-format="yyyy-MM-dd HH:mm:ss" type="datetime" placeholder="选择日期时间" :clearable="false">
                                    </el-date-picker>
                                    </div>
                                </div>
                            </el-form-item>
                            <el-form-item label="" prop="runtime">
                                <label class="dialog_item_label">执行时间段<i class="is-required" style="float: right">*</i></label> 
                                <el-select v-model="taskform.runtime" size="small"  multiple class="sel_runtime form_item_width"
                                        placeholder="请选择执行时间段"  :disabled="Boolean($route.query.disabled)">
                                    <el-option v-for="item in rantimelist" :key="item.value" :label="item.label" :value="item.label" >
                                    </el-option>
                                    </el-select> 
                            </el-form-item>  
                            <!-- <el-form-item label="" prop="execute_type" >
                                <label class="dialog_item_label">测试强度</label>
                                <el-select
                                :disabled="Boolean($route.query.disabled)"
                                v-model="taskform.testIntensity" size="small" 
                                    placeholder="请选择" class="form_item_width" > 
                                    <el-option v-for="(item,i) in strengthlist" :key="i" :label="item.label" :value="item.value">
                                    </el-option>
                                </el-select>
                            </el-form-item> -->
                        </el-col>
                    </el-row>
                </div>
                <div class="more_config">
                    <label  @click="showMoreconfig" style="cursor: pointer;">更多配置</label>
                    <i></i>
                </div>
                <div v-if="isShowMore">
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label">端口扫描范围</label>
                        <el-select  v-model="taskform.port_scan_type" size="small" placeholder="请选择" class="form_item_width"  :disabled="Boolean($route.query.disabled)"
                           @change="changePortSacn">
                            <el-option v-for="(item, index) in portRange_list" :key="index" :label="item.label" :value="item.value"></el-option>
                        </el-select>
                    </el-form-item>
                        <el-form-item label="" prop="scan_port" >
                        <label class="dialog_item_label"></label>
                        <el-input   v-model="taskform.scan_port" size="small" style="width:820px;" rows="6"  :disabled="Boolean($route.query.disabled)" 
                        autocomplete="off" type="textarea" resize="none" placeholder="请输入扫描端口"></el-input>
                    </el-form-item>
                    <div>
                        <div class="dialog_item_label">网站登录凭证</div> 
                        <el-button type="primary" size="mini" style="vertical-align: middle;margin-right: 27px"  :disabled="Boolean($route.query.disabled)"
                                @click="clickloginadd()">新增登录</el-button>

                        <div class="div_width"   style="margin-top:16px;margin-bottom:16px;width:930px">
                            <el-table :data="login_conf" size="small" style="width: 100%">
                                <el-table-column  prop="target" label="登录地址">
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.target }}</span>    
                                        <el-input v-else v-model="scope.row.target" size="small" ></el-input>    
                                    </template>
                                </el-table-column>
                                <el-table-column prop="protocol"  label="凭证类型"> 
                                    <template slot-scope="scope"> 
                                        <el-select ref="uptopro" :disabled="!scope.row.dataShow"
                                            size="small" style="width: 120px"
                                            v-model="scope.row.verifyType">
                                            <el-option v-for="item in agreementList" :key="item.name" :label="item.label" :value="item.value">
                                            </el-option>
                                        </el-select>
                                    </template>
                                </el-table-column>
                                <el-table-column  prop="voucher" label="凭证"> 
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.verifyValue }}</span>    
                                        <el-input  type="textarea" :rows="1" v-else v-model="scope.row.verifyValue" size="small" resize="none"></el-input>    
                                    
                                    </template>
                                </el-table-column>
                                <el-table-column   prop="verifyStatus" label="状态" width="80"> 
                                    <template slot-scope="scope"  >
                                        <span :class="[ 
                                                { 'tag_status tag_danger1': scope.row.verifyStatus == 3 } , 
                                                { 'tag_status tag_primary': scope.row.verifyStatus == 2 },
                                                { 'tag_status tag_success': scope.row.verifyStatus == 1 } ]" style="width: 60px;">{{ scope.row.verifyStatusZh }}</span>
                                    </template>
                                </el-table-column>
                                <el-table-column  label="操作"  width="150">
                                    <template slot-scope="scope"  >
                                            <el-link :underline="false" @click="loginSave(scope)" 
                                                v-if="scope.row.dataShow" > 保存 </el-link>
                                            <el-link :underline="false" @click="loginUpdate(scope)" > 编辑 </el-link>
                                            <el-link :underline="false" @click="loginDelete(scope)" > 删除 </el-link>
                                        </template>
                                </el-table-column>
                            </el-table>
                        </div> 
                    </div> 
                    <div>
                        <div class="dialog_item_label">Local Storage</div> 
                        <el-button type="primary" size="mini" style="vertical-align: middle;margin-right: 27px"  :disabled="Boolean($route.query.disabled)"
                                @click="localstorageAdd()">新增</el-button>
                        <div class="div_width"   style="margin-top:16px;margin-bottom:16px;width:930px">
                            <el-table :data="localstorage_conf" size="small" style="width: 100%">
                                <el-table-column  prop="target" label="key">
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.key }}</span>    
                                        <el-input v-else v-model="scope.row.key" size="small" ></el-input>    
                                    </template>
                                </el-table-column> 
                                <el-table-column  prop="voucher" label="value"> 
                                    <template slot-scope="scope">
                                        <span v-if="!scope.row.dataShow">{{ scope.row.value }}</span>    
                                        <el-input  type="textarea" :rows="1" v-else 
                                        v-model="scope.row.value" size="small" resize="none"></el-input>    
                                    
                                    </template>
                                </el-table-column> 
                                <el-table-column  label="操作"  width="150">
                                    <template slot-scope="scope"  >
                                            <el-link :underline="false" @click="localstorageSave(scope)" 
                                                v-if="scope.row.dataShow" > 保存 </el-link>
                                            <el-link :underline="false" @click="localstorageUpdate(scope)" > 编辑 </el-link>
                                            <el-link :underline="false" @click="localstorageDelete(scope)" > 删除 </el-link>
                                        </template>
                                </el-table-column>
                            </el-table>
                        </div>
                    </div>
                    <el-form-item> 
                        <label for="" class="dialog_item_label" >漏洞利用</label>
                        <el-switch v-model="isUse" class="elSwitch"  :disabled="Boolean($route.query.disabled)"> </el-switch> 
                    </el-form-item>
                    <el-form-item> 
                        <label for="" class="dialog_item_label" >安全测试</label>
                        <el-switch v-model="isSafetyTesting" class="elSwitch"  :disabled="Boolean($route.query.disabled)"> </el-switch> 
                    </el-form-item>
                    <el-form-item> 
                        <label for="" class="dialog_item_label" >横向移动</label>
                        <el-switch v-model="isPostPenetration" class="elSwitch"  :disabled="Boolean($route.query.disabled)"> </el-switch> 
                        <div v-show="isPostPenetration">
                            <postPenetrationModule ref="postPenetration" ></postPenetrationModule>
                        </div>
                    </el-form-item>
                    <el-form-item> 
                        <label for="" class="dialog_item_label" >渗透模式</label>
                        <el-radio-group v-model="taskform.mode" @change="handleradio">
                            <el-radio :label="1">通用渗透</el-radio>
                            <el-radio :label="2">定向渗透</el-radio>
                        </el-radio-group>
                
                    </el-form-item>
                    <el-form-item v-if="taskform.mode == 2">
                        <label class="dialog_item_label_m">定向节点</label>
                        <el-select 
                            v-model="taskform.dnode" 
                            placeholder="请选择" 
                            multiple
                            collapse-tags
                            style="width: 720px"
                        >
                            <el-option
                                v-for="item in nodelist"
                                :key="item.id"
                                :label="item.name"
                                :value="item.id"
                            >
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item> 
                        <label for="" class="dialog_item_label" >代理配置</label>
                        <el-switch v-model="isAgent" class="elSwitch"  :disabled="Boolean($route.query.disabled)"> </el-switch> 
                        <div v-show="isAgent">
                            <agentModule ref="agentModule" ></agentModule>
                        </div>
                    </el-form-item>
                </div>
            </el-form>
        </div>
    </div>
</template>
<style lang="less" scoped>
.createtask_box{
    padding: 24px;
    background: #fff;
    height: calc(100% - 129px);
    overflow-x: hidden;
    overflow-y: auto;
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
.form_item_width{
    width: 320px;
}
.margin_left_width{
    margin-left:110px;
}
.execute_type_label {
    padding-left: 8px;
    color: rgba(72, 72, 102, 0.64);
    font-size: 13px;
    margin-right: 32px;
    box-sizing: border-box;
}
.basic_config{
    position:relative;
}
.basic_config::before{
    left: 0;
    bottom: 0;
    width: 100%;
    height: 1px;
    content: "";
    position: absolute; 
    background: #E8E8F5;
    z-index: 1;
}
.more_config{
    margin:24px 0;
    label{ 
        color: #4C7AE3;
        font-size:13px;
    }
    i{
        display:inline-block;
        width:16px;
        height:16px;
        background:url(../../assets/images/show@2x.png);
        background-size: cover;
        vertical-align: middle;
    }
}
.dialog_item_label_m {
    display: inline-block;
    font-size: 14px;
    font-weight: 500;
    color: rgba(72, 72, 102, 0.87);
    padding-left: 11px;
    width: 104px; 
    line-height: 16px;
}
i.is-required {
    margin-right: 4px;
    color: #f56c6c;
    font-size: 12px; 
}
.template-desc{
    margin-left:113px;
    width:720px;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 1;
    word-break: break-all;
    line-height: 24px;
    color: rgba(72, 72, 102, 0.45);
    font-size: 12px;
    padding:0 4px;
    font-style:italic;
}
</style>
<script>
import { traffic } from '@/api/assetManagement.js'
import scene from '@/api/scene.js'
import { task,task_group } from '@/api/task.js'
import BannerBox from "@/components/BannerBox.vue"; 
import postPenetrationModule from "@/pages/sceneManagement/components/post_penetration_module.vue";
import agentModule from "./components/agentModule.vue"
var XLSX = require('xlsx');
export default{
    name:'createtask',
    components:{
        BannerBox,
        postPenetrationModule,
        agentModule
    },
    data(){
        return{
            type:this.$route.query.type, //1：任务组列表，2任务组详情列表新建,3 资产组渗透任务
            group_id: this.$route.query.group_id,  
            group_name:this.$route.query.group_name,
            flag: this.$route.query.flag, //1新建，6新增目标,3资产模块创建
            task_id:this.$route.query.task_id,//任务id
            taskName:this.$route.query.taskName,//任务id
            assetip:this.$route.query.TestObjectives, //资产传过来ip
            template_desc:'',
            taskform:{
                template:'', //任务场景 
                target:'',
                taskname:'',
                time:'',
                isVulntest:false, //快速漏洞检测
                execute_type:1,
                execute_cycletype:'',
                execute_cycletype_day:'',
                execute_cycletype_week:'',
                execute_cycletype_starttime:'',
                endtime:'',
                rantime:'', //执行时间段
                weight:'', //优先级
                testIntensity:'',//测试强度
                isCreateReport:false,
                reportType:[],
                report_task_template:'',
                report_target_template:'',
                reportemail:'',
                port_scan_type:100, 
                scan_port:'',
                vul_use_open:false,
                loginurl:'',
                cookie:'',
                headers:'',
                username:'',
                password:'',
                mode:1,
                dnode:[],
                node:'',
                proxy_server:'',
                proxy_name:'',
                vpnnode:'', 
                guess_is_open:false, 
                baseline_is_open:false, //基线
                baseline_type:1,
                baseline_scene:'',
                baseline_scene_name:'',
                excludeTarget:'',
                testMode:['1'], //测试模式，默认选中原理验证
            },
            isMore:false,
            rantypelist:[],
            rantimelist:[],
            templatelist:[],
            prioritylist:[], //优先级
            strengthlist:[], //强度
            rules:{
                target: [
                    { required: true, message: "测试范围不能为空", trigger: "blur" },
                    // { validator: checkTarget, trigger: ['blur'] }
                ],
                taskname: [
                {
                    required: true,
                    message: "任务名称不能为空",
                    trigger: ["blur", "change"],
                },
                ],
                execute_type:[
                    { required: true, message: "请选择执行方式", trigger: ["blur", "change"] },
                ],
                // runtime:[
                //     { required: true, message: "请选择执行时间段", trigger: ["blur", "change"] },
                // ]
            },
            portRange_list:[],
            portRangeValue:[],
            login_conf:[], 
            agreementList:[],
            cyclePlanningType:[],
            cyclePlanningMonthValue:[],
            cyclePlanningWeekValue:[], 
            config:{
                portScanConfig:{},
                websiteLoginConfig:{},
                testIntensity:0,
                webCrawlerConfig:{
                    localStorage:{
                        isOpen:true,
                        list:[]
                    }
                },
                lateralMove:{},
                safeTest:false,
                vulExploit:false,
                mode: {
                    mode: '',
                    distributeNodeId: ''
                },
                proxyConfig:{
                    "isOpen":false,
                    "addr":"",
                    "port":"",
                    "proto":1,
                    "isAuth":false,
                    "username":"",
                    "password":""
                }
            },
            isShowMore:false,
            localstorage_conf:[],
            isUse:false,
            isSafetyTesting:false,
            isPostPenetration:false,
            nodelist:[],
            isUpdate:false,
            isShow:false,
            lateralMoveJumpNum:[],
            postPenetration:{},
            user_id:0,
            isAgent:false,
            agent:{},
            agentlist:[],
        }
    },
    created(){
        // this.$store.state.activefirstMenu = '/task';
        if(!this.group_id){ //渗透任务
            this.$store.state.activefirstMenu = '/task';
        }
        else{ //任务组内
            this.$store.state.activefirstMenu = '/taskgroup';
        }
        this.user_id = localStorage.getItem('user_id-par');

        if(this.type == 3){
            //资产组-渗透任务
            this.$store.state.activefirstMenu = '/assettree';
            let asset = localStorage.getItem('checkedasset');
            this.taskform.target = this.assetip;
        }
        
    },
    mounted(){
        let that = this;
        this.getTaskEnum(); //任务枚举
        this.getSceneEnum(); //场景枚举
        this.geTasktemplate(); //场景类型列表
        if(this.$route.query.flag &&this.$route.query.flag == 1){
           
        }else if(this.$route.query.flag &&this.$route.query.flag == 2){
            // this.getTaskEnum();
            // this.getSceneEnum();
            // this.geTasktemplate();
            // console.log(this.$store.state.copyArr,'copyArr');
            setTimeout(() => {
                // 任务场景
                this.taskform.template = this.$store.state.copyArr.taskTemplateId
                this.gettemplateconfigbyid(this.taskform.template)
                // 任务名称
                this.taskform.taskname = this.$store.state.copyArr.taskName
                // 测试范围
                this.taskform.target = this.$store.state.copyArr.target
                // 执行方式
                this.taskform.execute_type = this.$store.state.copyArr.executeType
                // 优先级
                this.taskform.weight = this.$store.state.copyArr.weight;
                if(this.taskform.execute_type == 3){
                    //是周期任务赋值时间
                       if(this.$store.state.copyArr.executeJson&&this.$store.state.copyArr.executeJson.cyclePlanningType){
                            this.taskform.execute_cycletype = this.$store.state.copyArr.executeJson.cyclePlanningType
                            //this.taskform.execute_cycletype == 2是选择的月1是星期
                            if(this.taskform.execute_cycletype == 2){
                            this.taskform.execute_cycletype_day = this.$store.state.copyArr.executeJson.cyclePlanningValue

                            }else if(this.taskform.execute_cycletype == 1){
                                this.taskform.execute_cycletype_week = this.$store.state.copyArr.executeJson.cyclePlanningValue
                            }
                            this.taskform.execute_cycletype_starttime = this.$store.state.copyArr.executeJson.cyclePlanningHour
                            this.taskform.endtime = this.$store.state.copyArr.executeJson.endTime

                        }
                }else if(this.taskform.execute_type == 2){
                    //定时任务赋值时间
                    if(this.$store.state.copyArr.executeJson&&this.$store.state.copyArr.executeJson.startTime){
                        this.taskform.time = this.$store.state.copyArr.executeJson.startTime
                    }
                }
                // 执行时间段
                if(this.$store.state.copyArr.executeJson&&this.$store.state.copyArr.executeJson.runtimePeriod){
                    this.taskform.runtime = this.$store.state.copyArr.executeJson.runtimePeriod
                }
                //端口扫描范围
                if(this.$store.state.copyArr.config&&this.$store.state.copyArr.config.portScanConfig){
                    setTimeout(() => {
                        this.config = this.$store.state.copyArr.config
                        this.taskform.port_scan_type = this.$store.state.copyArr.config.portScanConfig.portScanType
                        // that.changePortSacn()
                        this.taskform.scan_port = this.$store.state.copyArr.config.portScanConfig.scanPort
                    }, 1000);
                }
                //网站登录凭证
                if(this.$store.state.copyArr.config&&this.$store.state.copyArr.config.websiteLoginConfig){
                   
                    this.$store.state.copyArr.config.websiteLoginConfig.list.forEach((item, i) => {
                    // if (i == scope.$index) {
                    //     item.dataShow = false;
                    //     item.verifyStatus = res.data.statusCode;
                    //     item.verifyStatusZh = res.data.status;
                    //     }
                    this.login_conf.push({
                        dataShow : false,
                        target : item.target,
                        verifyType : item.verifyType,
                        verifyValue : item.verifyValue,
                        verifyStatus : item.verifyStatus,
                        verifyStatusZh : item.verifyStatusZh
                    })
                    })


                    this.taskform.agreement = this.$store.state.copyArr.config.websiteLoginConfig.agreement
                }

                this.taskform.testIntensity = this.$store.state.copyArr.config.testIntensity;
                this.isUse = this.$store.state.copyArr.config.vulExploit;
                this.isSafetyTesting = this.$store.state.copyArr.config.safeTest;
                this.isPostPenetration = this.$store.state.copyArr.config.lateralMove.isOpen;
                
                // 添加测试模式的复制处理
                if(this.$store.state.copyArr.config.testMode) {
                    this.taskform.testMode = this.$store.state.copyArr.config.testMode.split(',');
                }

                this.postPenetration = this.$store.state.copyArr.config.lateralMove;  
            }, 500);
        
            // this.$nextTick(()=>{
            //     this.getTaskinfo();
            // })
        }
        else if(this.$route.query.flag &&this.$route.query.flag == 6){ //目标
            // this.getTaskEnum();
            // this.getSceneEnum();
            // this.geTasktemplate();
            this.$nextTick(()=>{
                this.getTaskinfo();
            })
            
        }else{
            // this.getTaskEnum();
            // this.getSceneEnum();
            // this.geTasktemplate();
        }
       
    },
    methods:{
        async getTaskEnum(){
            let _this = this;
            const res = await task.taskEnum();
            if(res.code == 200){
                this.rantypelist = res.data.executeType;
                this.rantimelist = res.data.runtimePeriod ||[];
                this.cyclePlanningType = res.data.cyclePlanningType;
                this.cyclePlanningMonthValue = res.data.cyclePlanningMonthValue;
                this.cyclePlanningWeekValue =  res.data.cyclePlanningWeekValue;
                this.agreementList = res.data.webLoginType;

                this.prioritylist = res.data.weight;
                this.strengthlist = res.data.testIntensity; 

                this.prioritylist.forEach(item =>{
                    if( item.default==1){
                        this.taskform.weight = item.value
                    }
                })
                this.strengthlist.forEach(item =>{
                    if( item.default==1){
                        this.taskform.testIntensity = item.value
                    }
                })
                
                this.agentlist = res.data.proxyProto
                // +++++++++默认选第一个时间
                if (typeof this.rantimelist == 'object' && this.rantimelist.length > 0 ) {
                    this.rantimelist.forEach((item, i) => {
                        this.taskform.runtime.push(item.label)
                    })
                    // this.taskform.runtime.push(this.rantimelist[0].label);
                    }
            }else{

            }
        },
        async getSceneEnum(){ //场景枚举
            let res = await scene.getSceneEnum();
            if(res.code == 200){ 
                this.portRange_list = res.data.portScan.portRange;
                this.portRangeValue = res.data.portScan.portRangeValue; 

                this.lateralMoveJumpNum = res.data.lateralMove;
  
                this.portRangeValue.forEach(item =>{
                    if(item.value == this.taskform.port_scan_type){ 
                        this.taskform.scan_port= item.label;
                    }
                })
            }
        },
        async geTasktemplate() { //任务场景  
            const res = await scene.sceneoptions();
            if (res.code == 200) { 
                this.templatelist = res.data.taskTemplate;
                this.templatelist.forEach(((item, i) => {
                    if (item.isDefault == 1 && this.flag!=2) {
                        this.taskform.template = item.id;
                        this.template_desc = item.describe;
                        this.gettemplateconfigbyid(item.id);
                    }
                }))
            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }
        },
        async gettemplateconfigbyid(_id){
            const res = await scene.getSceneinfo({
                taskTemplateId: _id
            })
            if (res.code == 200) {
                this.isShowMore = false;
                if(this.$route.query.flag == 1){
                    this.template_desc =res.data.describe;
                    let _config = res.data.config;
                    this.config = _config; 
                    // console.log(this.config)
                    this.taskform.port_scan_type = _config.portScanConfig.portScanType;
                    this.taskform.scan_port = _config.portScanConfig.scanPort; 
                   
                    this.isUse = res.data.config.vulExploit;
                    this.isSafetyTesting = res.data.config.safeTest;
                    this.isPostPenetration = res.data.config.lateralMove.isOpen;  
                    this.postPenetration = res.data.config.lateralMove;
                    
                    // 添加测试模式的场景配置处理
                    if(res.data.config.testMode) {
                        this.taskform.testMode = res.data.config.testMode.split(',');
                    } else {
                        this.taskform.testMode = ['1']; // 默认选中原理验证
                    }
                 
 
                }else{
                    this.template_desc =res.data.describe;
                    this.config = res.data.config;  
                }
         
            }else{
                
            }
        },
        async getTaskinfo(){ //获得任务详情
            const res = await task.taskCopy({
                taskId:this.task_id,
            })
            if(res.code == 200){ 
                let _this = this;
                this.taskform.template = res.data.taskTemplateId; 
                // this.template_desc =  res.data.
                this.taskform.target = '';
                this.taskform.execute_type = res.data.executeType; 
                this.taskform.taskname = res.data.taskName; 
                //执行时间段
                // let runtimePeriodarr = res.data.runtimePeriod.split(',');
                // this.taskform.runtime  = runtimePeriodarr;
                // 执行时间段
                if(res.data.executeJson&&res.data.executeJson.runtimePeriod){
                    this.taskform.runtime = res.data.executeJson.runtimePeriod
                }


                if(this.taskform.execute_type == 3){
                    //是周期任务赋值时间
                       if(res.data.executeJson&&res.data.executeJson.cyclePlanningType){
                            this.taskform.execute_cycletype = res.data.executeJson.cyclePlanningType
                            //this.taskform.execute_cycletype == 2是选择的月1是星期
                            if(this.taskform.execute_cycletype == 2){
                            this.taskform.execute_cycletype_day = res.data.executeJson.cyclePlanningValue

                            }else if(this.taskform.execute_cycletype == 1){
                                this.taskform.execute_cycletype_week = res.data.executeJson.cyclePlanningValue
                            }
                            this.taskform.execute_cycletype_starttime =res.data.executeJson.cyclePlanningHour
                            this.taskform.endtime = res.data.executeJson.endTime

                        }
                }else if(this.taskform.execute_type == 2){
                    //定时任务赋值时间
                    if(res.data.executeJson&&res.data.executeJson.startTime){
                        this.taskform.time =res.data.executeJson.startTime
                    }
                } 
                //端口扫描范围
                if(res.data.config&&res.data.config.portScanConfig){
                    setTimeout(() => { 
                        this.taskform.port_scan_type = res.data.config.portScanConfig.portScanType 
                        this.taskform.scan_port = res.data.config.portScanConfig.scanPort
                    }, 1000);
                }

                //网站登录凭证
                if(res.data.config&&res.data.config.websiteLoginConfig){
                   
                    res.data.config.websiteLoginConfig.list.forEach((item, i) => { 
                        this.login_conf.push({
                            dataShow : false,
                            target : item.target,
                            verifyType : item.verifyType,
                            verifyValue : item.verifyValue,
                            verifyStatus : item.verifyStatus,
                            verifyStatusZh : item.verifyStatusZh
                        })
                    })


                   this.taskform.agreement = res.data.config.websiteLoginConfig.agreement
               }
                this.taskform.weight = res.data.weight;
                this.taskform.testIntensity = res.data.config.testIntensity;
                this.isUse = res.data.config.vulExploit;
                this.isSafetyTesting = res.data.config.safeTest;
                this.isPostPenetration = res.data.config.lateralMove.isOpen;
                this.localstorage_conf =  res.data.config.webCrawlerConfig.localStorage.list ;
 
                this.postPenetration = res.data.config.lateralMove;
               //代理
               this.agent = res.data.config.proxyConfig;
               this.isAgent = res.data.config.proxyConfig.isOpen;

                this.config =  res.data.config;
                   
                
            }else{

            }
        },
        submithandle(){ //保存创建任务  
            this.$refs.form.validate(async (valid) => {
                if (valid) { 
                    if(this.flag == 6){ //新增目标
                        const res = await task.addtarget({
                            taskId:this.task_id,
                            userId:this.user_id ,
                            target:this.taskform.target.split("\n").join(","),
                        });
                        if (res.code  == 200) {
                            this.$message({
                                message: '新建目标成功',
                                type: "success",
                            });
                            this.$router.push({
                                path: `/taskDetail`,
                                query: { 
                                    id:this.task_id, 
                                    name:this.taskName,
                                }
                            });
                        } else {
                            this.$message({
                                message: res.msg,
                                type: "error",
                            });
                        } 

                    } else{
                        var json = this.saveparameter();

                        const loading = this.$loading()
                        const res = await task.taskSave(json);
                        loading.close();
                        
                        if (res.code  == 200) {
                            this.$message({
                                message: '保存任务成功',
                                type: "success",
                            });
                            if(!this.group_id ){ //渗透任务
                                if(this.type == 3){
                                    // this.assetpenetrationsync(res.data.task_id);
                                    // 资产管理接口异常，优化体验：跳转到任务列表查看执行状态
                                    this.$router.push({
                                        path: `/task`,
                                    });
                                }
                                else{
                                    this.$router.push({
                                        path: `/task`,
                                    });
                                }
                               
                            }else{ //任务组
                                this.groupbind(res.data.task_id);
                            }
                            
                        } else {
                            this.$message({
                                message: res.msg,
                                type: "error",
                            });
                        } 
                    }
                    
                }
            });
        },
        async groupbind(taskid){ //任务组与任务绑定
            const res = await task_group.groupBindTask({
                task_id:taskid,
                group_id:this.group_id,
            })
            if(res.code == 200){
                if(this.type == 1){
                    this.$router.push({
                        path: `/taskgroup`,
                    });
                }else if(this.type==2){
                    this.$router.push({
                        path: `/taskGroupDetail`,
                        query: { 
                            id: this.group_id, 
                            name:this.group_name,
                        }
                    });
                }
               
            }
        },
        async assetpenetrationsync(taskid){ //资产渗透
            const res = await traffic.assetpenetrationsync({
                taskId:taskid
            });
            if(res.code ==200){
                this.$router.push({
                    path: `/assetManagement`,
                });
            }
        },
        saveparameter(){ //获得保存参数
            var json = {};  
            json.userId = Number(this.user_id),
            json.taskTemplateId = this.taskform.template;
            json.target = this.taskform.target
               .split("\n")
               .join(","); 
            json.excludeTarget = this.taskform.excludeTarget
               .split("\n")
               .join(","); 
            json.taskName = this.taskform.taskname; 
            json.executeType = this.taskform.execute_type;
            json.weight = this.taskform.weight; 
        
            let executeJson = {};
            executeJson.runtimePeriod = this.taskform.runtime 
            //执行方式
            if (json.executeType == 2) {
               //定时
               //json.plan_time = this.taskform.time;
               executeJson.startTime = this.taskform.time
            }
           if (json.executeType == 3) {
               //周期 
               if (this.taskform.execute_cycletype == 1) {
                   //月 
                   executeJson.cyclePlanningValue =this.taskform.execute_cycletype_week ;
               }
               if (this.taskform.execute_cycletype == 2) {
                   //日 
                   executeJson.cyclePlanningValue =  this.taskform.execute_cycletype_day ;
               }  
               executeJson.cyclePlanningType = this.taskform.execute_cycletype; 
               executeJson.cyclePlanningHour = this.taskform.execute_cycletype_starttime;
               executeJson.endTime = this.taskform.endtime; 
           }
           json.executeJson = JSON.stringify(executeJson);
           //端口 
           this.config.portScanConfig.portScanType = this.taskform.port_scan_type ;
           this.config.portScanConfig.scanPort = this.taskform.scan_port;
 

           //登录凭证   
           this.config.websiteLoginConfig.isOpen=true;
           this.config.websiteLoginConfig.list = this.login_conf;

           this.config.testIntensity = this.taskform.testIntensity; //测试强度
           this.config.vulExploit=this.isUse;
           this.config.safeTest=this.isSafetyTesting;
           
           // 添加测试模式处理
           this.config.testMode = this.taskform.testMode.join(',');

           this.config.lateralMove.isOpen=this.isPostPenetration;
           this.postPenetration.targetNum = Number(this.postPenetration.targetNum)
           this.postPenetration.timeout = Number(this.postPenetration.timeout)
           this.config.lateralMove = this.postPenetration;

           //localstorage
           this.config.webCrawlerConfig.localStorage.isOpen = true;
           this.config.webCrawlerConfig.localStorage.list = this.localstorage_conf;
           this.config.mode.mode = this.taskform.mode
           if (this.taskform.mode === 2) {
            this.config.mode.distributeNodeId = this.taskform.dnode
           } else {
            this.config.mode.distributeNodeId = []
           }

           //代理 
        //    this.config.proxyConfig.isOpen = this.isAgent;
           if(this.isAgent){
                this.agent = this.$refs.agentModule.getAllData();
                this.config.proxyConfig = this.agent;
                this.config.proxyConfig.isOpen =  this.isAgent;
           }else{
                this.config.proxyConfig = this.agent;
                this.config.proxyConfig.isOpen  = false;
           }
           

           if(this.$route.query.flag &&this.$route.query.flag == 2){ //复制
                if(!this.isAgent){ 
                    this.config.proxyConfig = { 
                        isOpen:false,
                    };
                } else{
                    this.config.proxyConfig = { 
                        isOpen:true,
                        "addr":this.agent.Addr,
                        "port":this.agent.Port,
                        "proto":1,
                        "isAuth":false,
                    };
                    // if (!this.agent.isauth){
                    //     this.config.proxyConfig.isAuth = true;
                    //     this.config.proxyConfig.username = this.agent.username;
                    //     this.config.proxyConfig.password = this.agent.password;
                    // }
                }
           }
 
          
           json.config =  JSON.stringify(this.config);
           return json;

       }, 
        targetinput(e){
            if(this.flag!=6){
                this.taskform.target =
                    this.taskform.target.replace(/[^\S\r\n]/, "");
                this.taskform.taskname =
                    this.taskform.target.substr(0, 20) +
                    "_" +
                    this.$commonjs.nowtime();
            }
            
        },
        clickupload(){
            document.querySelector(".btnUploadID").click();
        },
        changeuploaID(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1);
            if (fileSuffix.indexOf("xls") != -1) {
                var reader = new FileReader();
                reader.onload = function (e) {
                  
                    var data = e.target.result;
                    if (that.rABS) {
                        that.wb = XLSX.read(btoa(fixdata(data)), {
                            type: "base64",
                        });
                    } else {
                        that.wb = XLSX.read(data, {
                            type: "binary",
                        });
                    } 
                    // let tt = that.wb.SheetNames[0];
                    // console.log(tt);

                    // let ttt = that.wb.Sheets[that.wb.SheetNames[0]];
                    // console.log(ttt)

                    let carData = XLSX.utils.sheet_to_json(
                        that.wb.Sheets[that.wb.SheetNames[0]],{header:1, }
                    ); 
                    let arr = []; 
                    for(var i=0;i<carData.length;i++){
                        let item = carData[i][0];
                        arr.push(item)
                    }

                    // for (var key in carData) {
                    //     for (var k in carData[key]) {
                    //         if (arr.indexOf(k) === -1) {
                    //             arr.push(k);
                    //         }
                    //         if (arr.indexOf(carData[key][k]) === -1) {
                    //             arr.push(carData[key][k]);
                    //         }
                    //     }
                    // }
                     
                    that.taskform.target = arr.join("\n");
                    let str = "";
                    arr.forEach((item) => {
                        str += item;
                    });
                    that.taskform.taskname =
                        str.substring(0, 20) + "_" + that.$commonjs.nowtime();
                };
                if (that.rABS) {
                    reader.readAsArrayBuffer(f);
                } else {
                    reader.readAsBinaryString(f);
                }
            } else if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) {
                        that.taskform.target = reader.result;
                        that.taskform.taskname =
                            reader.result.substring(0, 20) + "_" + that.$commonjs.nowtime();
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        changePortSacn(){
            this.portRangeValue.forEach(item =>{
                
                if(item.value == this.taskform.port_scan_type){ 
                    this.taskform.scan_port= item.label;
                    console.log(this.portRangeValue,"this.taskform.port_scan_type",this.taskform.scan_port);
                }
                
            })
            
        },
        clickloginadd(){
            this.login_conf.push({
                dataShow: true,
            })
        },
        async loginSave(scope) { 
            const res = await task.websitelogincheck({
                task_check_target: this.taskform.target,
                target: scope.row.target,
                verifyType: scope.row.verifyType,
                verifyValue: scope.row.verifyValue,
            })
            if (res.code == 200) {
                this.login_conf.forEach((item, i) => {
                    if (i == scope.$index) {
                        item.dataShow = false;
                        item.verifyStatus = res.data.statusCode;
                        item.verifyStatusZh = res.data.status;
                    }

                })
            } else {
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }

        },
        loginUpdate(scope) { //编辑
            this.$set(this.login_conf[scope.$index], 'dataShow', true)
        },
        loginDelete(scope) {  //登录删除
            this.login_conf.splice(scope.$index, 1);
        },
        showMoreconfig(){ //开关过更多配置
            this.isShowMore = !this.isShowMore;
            if(this.isShowMore){  
                this.$nextTick(()=>{ 
 
                    this.$refs.postPenetration.getEnum(this.lateralMoveJumpNum,1); 
                    this.$refs.postPenetration.getIsUpdate(true); 
                    this.$refs.postPenetration.getConifg(this.postPenetration); 
                    this.isPostPenetration = this.config.lateralMove.isOpen; 

                    this.taskform.port_scan_type = this.config.portScanConfig.portScanType;
                    this.taskform.scan_port = this.config.portScanConfig.scanPort; 
                   
                    this.isUse = this.config.vulExploit;
                    this.isSafetyTesting = this.config.safeTest; 
                    this.postPenetration = this.config.lateralMove; 
                    this.$refs.agentModule.getIsUpdate(true); 
                    this.$refs.agentModule.getConifg(this.agent); 
                    this.$refs.agentModule.getEnum(this.agentlist,1); 
                     
                    if(this.$route.query.flag &&this.$route.query.flag == 2){ //复制 
                       
                        this.agent = this.$store.state.copyArr.config.proxyConfig;
                      
                        this.$refs.agentModule.getConifg(this.agent) 
                        this.isAgent = this.$store.state.copyArr.config.proxyConfig.isOpen;
                    }
               }); 
            }
        },
        localstorageAdd(){ 
            this.localstorage_conf.push({
                dataShow: true,
            })
        },
        localstorageSave(scope){
            this.localstorage_conf.forEach((item, i) => {
                if (i == scope.$index) {
                    item.dataShow = false; 
                } 
            })
        },
        localstorageUpdate(scope){
            this.$set(this.localstorage_conf[scope.$index], 'dataShow', true)
        },
        localstorageDelete(scope){
            this.localstorage_conf.splice(scope.$index, 1);
        },
        async handleradio(e){
            if(e == 2){
                this.getNodelist();
            }
        },
        async getNodelist() {
            const res = await task.getNodeData();
            if (res.code === 200) {
                this.nodelist = res.data.list
            } else {
                this.$message.error(res.msg)
            }
         },
    }
}
</script>