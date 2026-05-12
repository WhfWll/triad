<template>
    <div>
        <div class="main-title"> 
            <router-link :underline="false" class="classA" :to="{ path: '/burpsuite' }"  >Burpsuite
            </router-link> 
            <label class="currentpagetitle" >   
                <span>新建任务</span>
            </label> 
        </div>
        <BannerBox tips="任务参数通过任务场景进行默认配置，用户可以在创建任务时对默认参数进行调整。" style="margin-bottom: 16px;">
            <el-button type="primary" size="small" @click="submithandle" > 执行扫描 </el-button> 
            <el-button  class="delCancel" size="small" @click="$router.push('/burpsuite')" > 取消 </el-button> 
        </BannerBox>
        <div class="createtask_box">
            <el-form :model="taskform2" ref="form" :rules="rules" style="height: 100%">
                <div class="basic_config">
                    <el-row :gutter="20">
                 
                        <el-col :span="12">
                            <el-form-item label="" prop="target"  class="taskNameClass">
                                <label class="dialog_item_label">测试范围<i class="is-required" style="float: right">*</i></label>
                                <el-input type="textarea"  v-model="taskform2.target"   size="small"
                                class="form_item_width"   placeholder="请输入测试范围" ></el-input>
                            </el-form-item> 
                            <el-form-item label="" prop="taskName"  class="taskNameClass">
                                <label class="dialog_item_label">任务名称<i class="is-required" style="float: right">*</i></label>
                                <el-input  v-model="taskform2.taskName"   size="small"
                                class="form_item_width"   placeholder="请输入任务名称" maxlength="50"></el-input>
                            </el-form-item> 
                           
                          
                        </el-col>
                    </el-row>
                </div>
            
            </el-form>
        </div>
    </div>
</template>
<style lang="less" scoped>
.delCancel{
    color: #f56c6c !important;
    border-color: #f56c6c !important;
}
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
import scene from '@/api/scene.js'
import { task } from '@/api/Tripartitetools.js'
import BannerBox from "@/components/BannerBox.vue";
export default{
    name:'createXray',
    components:{
        BannerBox
    },
    data(){
        return{
            flag: this.$route.query.flag, //1新建
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
                pentestpattern:1,
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
            },
            taskform2:{ 
                target:'',
                taskName:'',
            },
            isMore:false,
            rantypelist:[],
            rantimelist:[],
            templatelist:[],
            rules:{
                target: [
                    { required: true, message: "测试范围不能为空", trigger: "blur" },
                    // { validator: checkTarget, trigger: ['blur'] }
                ],
                taskName: [
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
                websiteLoginConfig:{}
            },
            isShowMore:false
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/burpsuite';
    },
    mounted(){
       
       
    },
    methods:{
        async getTaskEnum(){
            const res = await task.taskEnum();
            if(res.code == 200){
                this.rantypelist = res.data.executeType;
                this.rantimelist = res.data.runtimePeriod ||[];
                this.cyclePlanningType = res.data.cyclePlanningType;
                this.cyclePlanningMonthValue = res.data.cyclePlanningMonthValue;
                this.cyclePlanningWeekValue =  res.data.cyclePlanningWeekValue;
                this.agreementList = res.data.webLoginType;
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
                    if (item.isDefault == 1) {
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
                if(this.$route.query.flag == 1){
                    this.template_desc =res.data.describe;
                    let _config = res.data.config;

                    this.taskform.port_scan_type = _config.portScanConfig.portScanType;
                    this.taskform.scan_port = _config.portScanConfig.scanPort; 
                    this.config = _config;  
 
                }else{
                    this.template_desc =res.data.describe;
                }
         
            }
        },
        submithandle(){ //保存创建任务  
            this.$refs.form.validate(async (valid) => {
                if (valid) {  
                    // var json = this.saveparameter();  
                    const res = await task.taskSave2(this.taskform2);
 
                    if (res.code  == 200) {
                        this.$message({
                            message: '新建任务成功',
                            type: "success",
                        });
                        this.$router.push({
                            path: `/burpsuite`,
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
   
        targetinput(e){
            this.taskform.target =
                this.taskform.target.replace(/[^\S\r\n]/, "");
            this.taskform.taskname =
                this.taskform.target.substr(0, 20) +
                "_" +
                this.$commonjs.nowtime();
        },
        clickupload(){
            document.querySelector(".btnUploadID").click();
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
        showMoreconfig(){
            this.isShowMore = !this.isShowMore;
        },
    }
}
</script>