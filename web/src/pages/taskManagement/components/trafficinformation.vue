<template>
    <div>
        <!-- 任务-综合信息 --> 
        <div class="part1">
            <div class="target_showbox">
                <el-row :gutter="16">
                    <targetinfostate :label="'高危目标'" :content="''" :flag=7 :activelayoutflag=activelayoutflag :iconfont="''"
                        :className="'high_target'" :targetNumber=risk_level[0]  >
                    </targetinfostate>
                    <targetinfostate :label="'中危目标'" :content="''" :flag=8 :activelayoutflag=activelayoutflag :iconfont="''"
                        :className="'middle_target'" :targetNumber=risk_level[1]   > </targetinfostate>
                    <targetinfostate :label="'低危目标'" :content="''" :flag=9 :activelayoutflag=activelayoutflag :iconfont="''"
                        :className="'low_target'" :targetNumber=risk_level[2]   >
                    </targetinfostate>
                    <targetinfostate :label="'安全目标'" :content="''" :flag=10 :activelayoutflag=activelayoutflag :iconfont="''"
                        :className="'safe_target'" :targetNumber=risk_level[3]  >
                    </targetinfostate>
                </el-row> 
                <el-row :gutter="16" style="margin-top:15px">
                    <targetinfostate :label="'致命漏洞'" :content="''" :flag=2 :activelayoutflag=activelayoutflag
                        :iconfont="'iconcunhuomubiao'" :className="'exist_target'" :targetNumber=vuln[0]   
                         > </targetinfostate> 
                    <targetinfostate :label="'高危漏洞'" :content="''" :flag=3 :activelayoutflag=activelayoutflag
                        :iconfont="'iconbucunhuomubiao'" :className="'no_exist_target'" :targetNumber=vuln[1]  
                        >
                    </targetinfostate>
                    <targetinfostate :label="'中危漏洞'" :content="''" :flag=6 :activelayoutflag=activelayoutflag
                        :iconfont="'iconweijiancemubiao'" :className="'no_detect_target'" :targetNumber=vuln[2]  
                         >
                    </targetinfostate>
                    <targetinfostate :label="'低危漏洞'" :content="''" :flag=1 :activelayoutflag=activelayoutflag
                        :iconfont="'iconmubiaozongshu'" :className="'target_number'" :targetNumber=vuln[3]  
                         > </targetinfostate>
                </el-row> 
            </div> 
        </div>
        <div class="part2">
            <div class="part_title" style="font-size:14px;">任务信息</div> 
            <el-descriptions title="" 
                :colon="false" 
                :column="4"  
                :contentClassName="'contentClassName'" >
                <el-descriptions-item label="任务名称">
                    <span  v-if="!isUpdate">{{ taskinfo.taskName }}</span>
                    <el-input v-model="Addtrafficform.taskname" size="small" v-else
                            placeholder="请输入任务名称" ></el-input></el-descriptions-item>
                <el-descriptions-item label="任务风险">
                    <span :class="[ 
                                { 'riskstyle risk_hight': taskinfo.riskLevel == 2 } ,
                                { 'riskstyle risk_middle': taskinfo.riskLevel == 3 },
                                { 'riskstyle risk_low': taskinfo.riskLevel == 4 },
                                { 'riskstyle risk_nofind': taskinfo.riskLevel == 5 }]"><i></i>
            
                            {{ taskinfo.riskLevelName }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="创建时间">{{ taskinfo.createTime }}</el-descriptions-item>
                <el-descriptions-item label="执行时间">{{ taskinfo.createTime }}</el-descriptions-item>
                <el-descriptions-item label="任务时长">
                    <span  v-if="!isUpdate">{{taskinfo.expireTimeName}}</span>
                    <el-select v-model="Addtrafficform.duration" size="small" placeholder="请选择"  v-else >
                            <el-option v-for="(item, index) in durationlist" :key="index" :label="item.label" :value="item.value"></el-option>
                        </el-select>
                </el-descriptions-item>
                <el-descriptions-item label="劫持监听网卡">
                    <span v-if="!isUpdate">
                        {{taskinfo.networkCard}}
                    </span>
                    <el-select v-model="Addtrafficform.networkcard" size="small" placeholder="请选择" v-else>
                            <el-option v-for="(item, index) in networkcardlist" :key="index" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                   
                </el-descriptions-item>
                <el-descriptions-item label="劫持监听端口">
                    <span v-if="!isUpdate">{{taskinfo.port}}</span>
                    <el-input v-model="Addtrafficform.networkport" size="small"  v-else placeholder="">
                        </el-input></el-descriptions-item>
                <el-descriptions-item label="提交者">{{taskinfo.userName}}</el-descriptions-item>
            </el-descriptions>
        </div>
        <div class="part2">
            <div class="part_title" style="font-size:14px;">检测目标 </div>
            <div class="info_param  info_paramheight220"  v-if="!isUpdate"> 
                <div class="targetTxt" style="margin-top:6px;">
                    <p class="targetcss"> {{target}}</p>
                </div>
            </div>
            <div class="info_param" v-else>
                <el-input type="textarea" :rows="4"  v-model="Addtrafficform.target"
                    autocomplete="off" placeholder="测试目标不能为空" resize="none" ></el-input>
            </div>
        </div>
        <div class="part2">
            <div class="part_title" style="font-size:14px;">漏洞插件 </div>
            <div class="info_param  info_paramheight220"  v-if="!isUpdate"> 
                <div class="targetTxt" style="margin-top:6px;">
                    <p class="targetcss"> {{taskinfo.vulConfig_zh}}</p>
                </div>
            </div>
            <div class="info_param" v-else>
                <el-select size="small" 
                    v-model="Addtrafficform.vulname"
                    multiple
                    filterable
                    allow-create
                    default-first-option 
                    placeholder=""
                    style="width: 100%;">
                    <el-option
                        v-for="(item,i) in vulnamelist"
                        :key="i"
                        :label="item.label"
                        :value="item.value">
                    </el-option>
                </el-select>
            </div>
        </div>
        <div class="part2">
            <div class="part_title" style="font-size:14px;">待测凭证 </div>
            <div class="info_param  info_paramheight220"  v-if="!isUpdate"> 
                <div class="targetTxt" style="margin-top:6px;">
                    <p class="targetcss" >
                        <span v-if="taskinfo.otherConfig.waitCred.pattern==1">cookie</span><span  v-else>header</span>
                        ：{{taskinfo.otherConfig.waitCred.value}}
                    </p> 
                </div>
            </div>
            <div class="info_param" v-else>
                <el-select  
                    v-model="Addtrafficform.waitCred_pattern" 
                    size="small" placeholder="请选择" style="width:100px;margin-right:20px"
                    >
                    <el-option v-for="(item, index) in credPatternlist" :key="index" :label="item.label"
                        :value="item.value"></el-option>
                </el-select>  
                <el-input   v-model="Addtrafficform.waitCred_value"   size="small" style="width:calc(100% - 120px)"
                        ></el-input>
            </div>
        </div>
        <div class="part2">
            <div class="part_title" style="font-size:14px;">fuzz参数 </div>
            <div class="info_param  info_paramheight220" v-if="!isUpdate"> 
                <div class="targetTxt" style="margin-top:6px;">
                    <p class="targetcss">字符型： {{taskinfo.otherConfig.FuzzParam.character}}</p>
                </div>
                <div class="targetTxt" style="margin-top:6px;">
                    <p class="targetcss">数字型： {{taskinfo.otherConfig.FuzzParam.number}}</p>
                </div>
            </div>
            <div class="info_param" v-else>
                <div style="margin: 5px 0">
                    <label class="dialog_item_label_m topline" style="margin-right: 20px;" >字符型</label>
                        <el-select size="small" 
                            v-model="Addtrafficform.fuzzParam.character"
                            multiple
                            filterable
                            allow-create
                            default-first-option 
                            placeholder=""
                             style="width: calc(100% - 60px);" >
                            <el-option
                                v-for="(item,i) in characterlist"
                                :key="i"
                                :label="item"
                                :value="item">
                            </el-option>
                        </el-select>
                </div>
                <div  style="margin: 5px 0">
                    <label class="dialog_item_label_m topline" style="margin-right: 20px;"  >数字型</label>
                    <el-select size="small" 
                        v-model="Addtrafficform.fuzzParam.number"
                        multiple
                        filterable
                        allow-create
                        default-first-option 
                        placeholder=""
                        style="width: calc(100% - 60px);" >
                        <el-option
                            v-for="(item,i) in numberlist"
                            :key="i"
                            :label="item"
                            :value="item">
                        </el-option>
                    </el-select>
                </div>
            </div>
        </div>
        <div class="part2">
            <div class="part_title" style="font-size:14px;">fuzz字典 </div>
            <div class="info_param  info_paramheight220" v-if="!isUpdate"> 
                <div class="targetTxt" style="margin-top:6px;">
                    <p class="targetcss">字符型： {{taskinfo.otherConfig.FuzzDict.character}}</p>
                </div>
                <div class="targetTxt" style="margin-top:6px;">
                    <p class="targetcss">数字型：{{taskinfo.otherConfig.FuzzDict.number}}</p>
                </div>
            </div>
            <div  class="info_param" v-else>
                <div style="margin: 5px 0">
                    <label class="dialog_item_label_m topline" style="vertical-align: top; margin-right: 20px; " >字符型</label>
                        <el-input   v-model="Addtrafficform.fuzzDict.character" size="small" 
                        style="width: calc(100% - 60px);" rows="4"  
                            autocomplete="off" type="textarea" resize="none" placeholder=""></el-input>
                </div>
                <div style="margin: 5px 0">
                    <label class="dialog_item_label_m topline" style="vertical-align: top; margin-right: 20px; " >数字型</label>
                    <el-input   v-model="Addtrafficform.fuzzDict.number" size="small" 
                            style="width: calc(100% - 60px);" rows="4"  
                                autocomplete="off" type="textarea" resize="none" placeholder=""></el-input>
                </div>
            </div>
        </div>
        <div class="part2">
            <div class="part_title" style="font-size:14px;">响应字典 </div>
            <div class="info_param  info_paramheight220" v-if="!isUpdate"> 
                <div class="targetTxt" style="margin-top:6px;">
                    <p class="targetcss">json关键字： {{taskinfo.otherConfig.Response.jsonKeyword}}</p>
                </div>
                <div class="targetTxt" style="margin-top:6px;" v-if="taskinfo.otherConfig.Response.noJsonSwitch==true">
                    <p class="targetcss">非json关键字： {{taskinfo.otherConfig.Response.noJsonKeyword}}</p>
                </div>
            </div>
            <div  class="info_param" v-else>
                <div style="margin: 5px 0"> 
                    <label class="dialog_item_label_m topline" style="vertical-align: top; margin-right: 20px; " >json关键字</label>
                    <el-input   v-model="Addtrafficform.response.jsonKeyword" size="small" 
                        rows="4"  style="width: calc(100% - 90px);"
                    autocomplete="off" type="textarea" resize="none" placeholder="请输入json关键字"></el-input>           
                </div>
                <div style="margin: 5px 0">
                    <label class="dialog_item_label_m topline" style=" margin-right: 20px; " >非json关键字</label>
                    <el-switch
                        v-model="Addtrafficform.response.noJsonSwitch"
                        class="elSwitch"  >
                    </el-switch> 
                </div>
                <div style="margin: 5px 0">
                    <el-input   v-model="Addtrafficform.response.noJsonKeyword" size="small"   
                    style="width: calc(100% - 90px); margin-left: 82px;"  rows="4"  
                        autocomplete="off" type="textarea" resize="none" ></el-input>
                </div>
            </div>
        </div>


    </div>
</template>
<style lang="less" scoped>
.el-descriptions{
    padding: 24px 13px;
}
/deep/ .el-descriptions .contentClassName{
    margin-left:20px;
    color:rgba(72, 72, 102, 0.64);
}
.target_showbox {
    background:#F7F7FB;
                // min-height: 200px;
                border-radius: 4px;
                // padding: 24px;
                box-sizing: border-box;
        
                .targetMsg_info {
                    min-height: 82px;
                    // max-height: 400px; 
                    background: #fff;
                    margin-top: 8px;
                    border-radius: 4px;
                    padding: 16px 24px;
                    box-sizing: border-box;
        
                    .targetdetail {
                        padding: 0 14px;
                        font-size: 13px;
                        color: rgba(72, 72, 102, 0.64);
                        line-height: 20px;
                        font-weight: 500;
                        max-height: 640px;
                        overflow-y: auto;
                        word-wrap: break-word;
                        word-break: normal;
                    }
        
                    .part_title {
                        font-size: 13px;
                    }
                }
                }
.part1 {
       // padding: 0 24px;
       padding-bottom: 15px;
       // border-bottom: 1px solid #E8E8F5;
   }
.part2 {
    padding: 24px 24px 0;
    background: #fff;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    border-radius: 4px;
    box-sizing: border-box;
    margin-bottom: 15px;
}

.part2 {
    .part_smalltitle {
        border-left: 3px solid #4C7AE3;
        margin-bottom: 13px;
    }

    label {
        font-family: SourceHanSansCN-Medium, SourceHanSansCN;
        font-size: 13px;
        font-weight: 500;
        color: rgba(72, 72, 102, 0.87);
    }

    span {
        color: rgba(72, 72, 102, 0.64);
    }
}


.part_title {
    font-size: 14px;
    // margin-bottom: 16px;
    font-weight: 500;
    border-left: 3px solid #4C7AE3;
    padding-left: 10px;
    height: 14px;
    line-height: 14px;
    color: rgba(72, 72, 102, 0.87);
}

.part2_sm_title {
    font-size: 13px;
    font-weight: 500;
    border-left: 3px solid #fff;
    padding-left: 10px;
    height: 14px;
    line-height: 16px;
    color: #484866;
}

.Smalltitle {
    padding-left: 15px;

    .heightauto {
        margin-top: 10px;
        color: #484866;
        width: 24.5%;
        display: inline-block;
    }

    .heightauto1 {
        margin-top: 10px;
        color: #484866;
        width: 100%;
        display: inline-block;
    }

    .portcss {
        margin-top: 8px;
        color: rgba(72, 72, 102, 0.64);
        word-wrap: break-word;
        line-height: 20px;
    }

    span {
        margin: 0 25px;
    }
}

.Smalltitle4 {
    padding-left: 15px;

    .heightauto {
        margin-top: 10px;
        color: #484866;
        width: 25%;
        display: inline-block;
    }
}

.Smalltitle2 {
    padding-left: 15px;

    .heightauto {
        margin-top: 10px;
        color: #484866;
    }

    span {
        margin: 0 25px;

    }
}
.info_param {
    border-bottom: 1px solid #E8E8F5;
}
.info_param {
    position: relative;
    // margin-bottom: 24px;
    border-radius: 4px;
    padding: 24px 13px;
    box-sizing: border-box;
    font-size: 13px;
    font-weight: 500;
    overflow: hidden;

    .targetbg {
        position: absolute;
        display: inline-block;
    }
 
    .targetTxt {
        font-size: 13px;
        color: #484866;
        margin-top: 15px;
        // padding-left: 12px;
        // overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        // cursor: pointer;
        // position: absolute;
        z-index: 1;
        display: inline-block;
        width: calc(100% - 24px);
        box-sizing: border-box;
    }

    .info_param_txt {
        margin-top: 16px;

        >div {
            margin-bottom: 8px;
            height: 18px;
            line-height: 18px;
            padding-left: 12px;

            label {
                display: inline-block;
                color: rgba(255, 255, 255, 0.7);
            }

            span {
                display: inline-block;
                color: rgba(255, 255, 255, 1);
            }
        }

        >div:last-child {
            margin-bottom: 0;
        }
    }

    .info_param_crawler>div label {
        width: 100px;
    }

    .info_param_guess>div label {
        width: 70px;
    }

    .info_param_login>div label {
        width: 86px;
    }

    .switch_box {
        float: right;
        width: 40px;
        height: 20px;
        line-height: 15px;
        border-radius: 20px;
        position: relative;

        label {
            position: absolute;
            font-size: 13px;
        }

        span {
            position: absolute;
            display: inline-block;
            width: 16px;
            height: 16px;
            border-radius: 50%;
        }

        &.switch_open {
            background: #fff;

            label {
                left: 4px;
                top: 3px;
            }

            span {
                right: 3px;
                top: 2px;
            }
        }

        &.switch_close {
            label {
                color: #fff;
                right: 5px;
                top: 2px;
            }

            span {
                left: 2px;
                top: 2px;
            }
        }
    }
}
.info_paramheight220 {
    border-bottom: 1px solid #E8E8F5;

    .targetcss {
        max-height: 156px;
        white-space: normal;
        word-wrap: break-word;
        overflow-y: auto;
        line-height: 20px;
        color: rgba(72, 72, 102, 0.64)
    }
}

.info_paramheight100 {
    height: 100px;
    border-bottom: 1px solid #E8E8F5;
}

.info_paramheight150 {
    height: 150px;
    border-bottom: 1px solid #E8E8F5;
}

.info_paramheight125 {
    height: 125px;
    border-bottom: 1px solid #E8E8F5;
}
</style>
<script>
import targetinfostate from "./Targetinfostate.vue";
import { traffic } from '@/api/traffic.js'
export default {
    name:'',
    components: {
        targetinfostate, 
    },
    props:{
        task_id:{}
    },
    data(){
        return{
            activelayoutflag: false, 
            taskinfo: {
                taskName: '', 
                userName: "", 
                createTime: '',
                expireTime: '', 
                expire:'', 
                networkCard:'', 
                port:'',
                riskLevel:'',
                riskLevelName:'',
                vulConfig:'',
                otherConfig:{
                    FuzzDict:{},
                    FuzzParam:{},
                    Response:{},
                    waitCred:{},
                }
            },
            target:'',
            risk_level:[],
            vuln:[],
            target_title:'',
            target_info:'',
            isUpdate:false,
            Addtrafficform:{
                target:'',
                taskname:'',
                duration:'',//时长
                networkcard:'',
                networkport:'',
                vulname:[],
                response:{
                    jsonKeyword:'',
                    noJsonSwitch:true,
                    noJsonKeyword:'',
                },
                fuzzParam:{
                    character:[],
                    number:[],
                },
                fuzzDict:{
                    character:'',
                    number:'',
                },
                waitCred_pattern:'',
                waitCred_value:'',
            },
            vulnamelist:[],
            credPatternlist:[],
            characterlist:[],
            numberlist:[],
            durationlist:[],
            networkcardlist:[],
        }
    },
    created() {
        this.$store.state.activefirstMenu = "/traffic";
    },
    mounted() { 
        this.getEnum();
    },
    methods: {
        handleupdate(falg){
            this.isUpdate = falg; 
            this.Addtrafficform.target = this.target;
            this.Addtrafficform.taskname =  this.taskinfo.taskName;
            this.Addtrafficform.duration = this.taskinfo.expireTime;
            this.Addtrafficform.networkcard = this.taskinfo.networkCard;
            this.Addtrafficform.networkport = this.taskinfo.port;

            this.Addtrafficform.vulname = this.taskinfo.vulConfig.split(',');
            // console.log(this.taskinfo.vulConfig)
            // this.Addtrafficform.fuzzParam.character = this.taskinfo.otherConfig.FuzzParam.character ;
            // this.Addtrafficform.fuzzParam.number = this.taskinfo.otherConfig.FuzzParam.number;

            this.Addtrafficform.fuzzDict.character = this.taskinfo.otherConfig.FuzzDict.character;
            this.Addtrafficform.fuzzDict.number = this.taskinfo.otherConfig.FuzzDict.number;

            this.Addtrafficform.response.jsonKeyword = this.taskinfo.otherConfig.Response.jsonKeyword;
            this.Addtrafficform.response.noJsonSwitch = this.taskinfo.otherConfig.Response.noJsonSwitch;
            this.Addtrafficform.response.noJsonKeyword = this.taskinfo.otherConfig.Response.noJsonKeyword;

            this.Addtrafficform.waitCred_pattern = this.taskinfo.otherConfig.waitCred.pattern;
            this.Addtrafficform.waitCred_value = this.taskinfo.otherConfig.waitCred.value;
           

        },
        async getEnum(){
            const res = await traffic.trafficEnum();
            if(res.code == 200){
                this.durationlist = res.data.expireTime;
                this.networkcardlist = res.data.networkCard;
                this.credPatternlist = res.data.credPattern;
                this.characterlist = res.data.fuzzParam.character == ''?[]:res.data.fuzzParam.character.split(',');
                this.numberlist =  res.data.fuzzParam.number == ''?[]:res.data.fuzzParam.number.split(',');
                this.vulnamelist = res.data.vulName;

                //默认值
                this.Addtrafficform.fuzzParam.character = this.characterlist;
                this.Addtrafficform.fuzzParam.number = this.numberlist;
                this.Addtrafficform.response.jsonKeyword = res.data.response.jsonKeyword;
                this.Addtrafficform.response.noJsonKeyword = res.data.response.noJsonKeyword;

                
                let result = this.vulnamelist.filter((item) => {
                     return item.isDefault
                })
                this.Addtrafficform.vulname  = result.map(item => item.value); 
               

            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
      
        async getData(){
            const res = await traffic.trafficbase({
                flowTaskId:this.task_id
            });
            if(res.code == 200){ 
                this.taskinfo = res.data; 
                this.target = res.data.target;
                this.risk_level = res.data.targetNum;
                this.vuln  = res.data.vulNum;

            }else{

            } 
        },
        getAllData(){
            this.Addtrafficform.userId = this.taskinfo.userId;
            return this.Addtrafficform;
        },
        
    },
}
</script>