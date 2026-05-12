<template>
    <div>
        <div class="main-title" > 
            <router-link :underline="false" class="classA" :to="{ path: '/vulScanTask' }"  >漏扫任务
            </router-link> 
            <label class="currentpagetitle" >   
                <span>创建任务</span>
            </label> 
        </div>
        
        <BannerBox ref="BannerBox" tips="执行漏洞扫描任务" style="margin-bottom: 16px;">
            <el-button type="primary" size="small" @click="submithandle" > 执行任务 </el-button> 
        </BannerBox>
        <div class="createtask_box">
            <el-form :model="taskform" ref="form" :rules="rules" style="height: 100%">
                <div class="basic_config">
                    <el-row :gutter="20">
                        <el-col :span="12">
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
                            
                        </el-col>
                    </el-row>

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


                    <el-form-item label="" prop="onlyPortScan">
                        <label class="dialog_item_label">仅资产探测</label>
                         <!-- 开关组件 -->
                        <el-switch
                        v-model="taskform.asset_detect"
                        active-text="开启"
                        inactive-text="关闭">
                    </el-switch>
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
    min-height: calc(100% - 39px);
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
    height: 0px;
    margin-bottom: 50px;
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
import scene from '@/api/scene.js'
import { task } from '@/api/task.js'
import { vulscan } from '@/api/vulscan.js'

import BannerBox from "@/components/BannerBox.vue"; 
import postPenetrationModule from "@/pages/sceneManagement/components/post_penetration_module.vue";
var XLSX = require('xlsx');
export default{
    name:'createtask',
    components:{
        BannerBox,
        postPenetrationModule,
    },
    data(){
        return{
            type:this.$route.query.type, //1：任务组列表，2任务组详情列表新建,3 资产组渗透任务
            group_id: this.$route.query.group_id,  
            group_name:this.$route.query.group_name,
            flag: this.$route.query.flag, //1新建，6新增目标
            task_id:this.$route.query.task_id,//任务id
            taskName:this.$route.query.taskName,//任务id
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
                dnode:'',
                node:'',
                proxy_server:'',
                proxy_name:'',
                vpnnode:'', 
                guess_is_open:false, 
                baseline_is_open:false, //基线
                baseline_type:1,
                baseline_scene:'',
                baseline_scene_name:'',
                asset_detect: false,
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
        this.$store.state.activefirstMenu = '/vulScanTask';
        }
        else{ //任务组内
            this.$store.state.activefirstMenu = '/taskgroup';
        }
        this.user_id = localStorage.getItem('user_id-par');

        if(this.type == 3){
            //资产组-渗透任务
            this.$store.state.activefirstMenu = '/assetManagement';
            let asset = localStorage.getItem('checkedasset');
            this.taskform.target = asset;
        }
    },
    mounted(){
        let that = this;
        this.getSceneEnum(); //场景枚举
        this.geTasktemplate(); //场景类型列表
        if(this.$route.query.flag &&this.$route.query.flag == 1){
           
        }else if(this.$route.query.flag &&this.$route.query.flag == 2){

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
                // // 执行时间段
                // if(this.$store.state.copyArr.executeJson&&this.$store.state.copyArr.executeJson.runtimePeriod){
                //     this.taskform.runtime = this.$store.state.copyArr.executeJson.runtimePeriod
                // }
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
                // this.rantimelist = res.data.runtimePeriod ||[];
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
                        var json = this.saveparameter();
                        const loading = this.$loading()
                        const res = await vulscan.taskSave(json);
                        loading.close();
                        
                        if (res.code  == 200) {
                            this.$message({
                                message: '保存任务成功',
                                type: "success",
                            });
                            this.$router.push({
                                path: `/vulScanTask`,
                            });
                            
                        } else {
                            this.$message({
                                message: res.msg,
                                type: "error",
                            });
                        } 
                    
                    
                }
            });
        },
      
        saveparameter(){ //获得保存参数
           
           var json = {};  
           json.userId = Number(this.user_id),
           json.target = this.taskform.target
               .split("\n")
               .join(","); 
           json.name = this.taskform.taskname;                   
           //端口   
           json.toScanPort = this.taskform.scan_port
           json.only_port_scan = this.taskform.asset_detect
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
                     
                    if(this.$route.query.flag &&this.$route.query.flag == 2){ //复制 
                       
                        this.agent = this.$store.state.copyArr.config.proxyConfig;
                      
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