<template>
    <div>
        <!-- 新建流量分析 -->
        <div class="main-title"> 
            <router-link :underline="false" class="classA" :to="{ path: '/traffic' }"  >被动流量
            </router-link> 
            <label class="currentpagetitle"> 
                <span>新建流量分析</span> 
            </label> 
        </div>
        <BannerBox tips="流量分析任务是通过劫持分析用户流量，来分析是否存在安全漏洞。" style="margin-bottom: 16px;">
            <el-button type="primary" size="small" @click="submithandle" > 执行任务 </el-button> 
        </BannerBox>
        <div class="formlist context_box_bg mt16">
            <el-form :model="Addtrafficform" ref="form" :rules="rules" style="height: 100%">
                <div  style="position: relative">
                    <label class="dialog_item_label" style="vertical-align: top; width: 103px">测试目标</label>
                    <el-form-item label=" " prop="target" label-width="10px" style="display: inline-block; margin-right: 0">
                        <el-input type="textarea" :rows="4" @input="targetinput" v-model="Addtrafficform.target"
                             autocomplete="off" placeholder="测试目标不能为空" resize="none"
                            style="width: 320px; margin-bottom: 10px; "></el-input>
                        <div  >
                            <el-button type="primary" size="small" style="vertical-align: top; margin-right: 27px"
                                @click="clickupload()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32)">只能上传.txt或.xls或.xlsx格式文件</span>
                            <input type="file" class="btnUploadID" ref="upload" @change="changeuploaID($event)" style="display: none"
                                id="input-file-ID" />
                        </div>
                    </el-form-item>
                    <el-tooltip placement="right-start">
                        <div slot="content">
                            测试范围支持IP、IP段、域名、URL，多个不同目标用“换行”隔开；<br />
                            示例：<br />“192.168.0.127”、“192.168.0.10-127”、<br />“4dogs.cn”、“www.4dogs.cn”、“http://www.4dogs.cn/aqjc/”、<br />“192.168.0.127:8000”
                        </div>
                        <i class="iconfont icontishi icontsstyle" style="position: absolute; top: 18px"></i>
                    </el-tooltip> 
                </div>
                <div>
                    <el-form-item label="" prop="taskname" class="taskNameClass">
                        <label class="dialog_item_label">任务名称</label>
                        <el-input v-model="Addtrafficform.taskname" :class="tasknameError" size="small" style="width: 320px"
                            placeholder="请输入任务名称" maxlength="50"></el-input>
                    </el-form-item>
                </div>
                <div> 
                    <el-form-item>
                        <label class="dialog_item_label">劫持监听网卡</label>
                        <el-select v-model="Addtrafficform.networkcard" size="small" placeholder="请选择" style="width: 320px">
                            <el-option v-for="(item, index) in networkcardlist" :key="index" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item>
                        <label class="dialog_item_label">劫持监听端口</label>
                        <el-input v-model="Addtrafficform.networkport" size="small" style="width: 320px" placeholder="">
                        </el-input>
                    </el-form-item>
                    <el-form-item>
                        <label class="dialog_item_label">漏洞插件</label>
                        <!-- <el-select v-model="Addtrafficform.duration" size="small" placeholder="请选择" style="width: 320px">
                            <el-option v-for="(item, index) in durationlist" :key="index" :label="item.label" :value="item.value"></el-option>
                        </el-select> -->
                        <el-select size="small" 
                                v-model="Addtrafficform.vulname"
                                multiple
                                filterable
                                allow-create
                                default-first-option 
                                placeholder=""
                                class="form_item_width"  >
                                <el-option
                                    v-for="(item,i) in vulnamelist"
                                    :key="i"
                                    :label="item.label"
                                    :value="item.value">
                                </el-option>
                            </el-select>
                    </el-form-item>
                    <el-form-item>
                        <label class="dialog_item_label">任务时长</label>
                        <el-select v-model="Addtrafficform.duration" size="small" placeholder="请选择" style="width: 320px">
                            <el-option v-for="(item, index) in durationlist" :key="index" :label="item.label" :value="item.value"></el-option>
                        </el-select>
                    </el-form-item>
                </div>
                <div class="more_config">
                    <label  @click="showMoreconfig" style="cursor: pointer;">更多配置</label>
                    <i></i>
                </div>
                <div v-if=" isShowMore">
                    <el-form-item label="" prop="template" >
                        <label class="dialog_item_label">待测凭证</label>
                        <el-select  
                            v-model="Addtrafficform.waitCred_pattern" 
                            size="small" placeholder="请选择" style="width:100px;margin-right:22px"
                           >
                            <el-option v-for="(item, index) in credPatternlist" :key="index" :label="item.label"
                                :value="item.value"></el-option>
                        </el-select>  
                        <el-input   v-model="Addtrafficform.waitCred_value"   size="small" style="width:200px"
                              ></el-input>
                        <el-tooltip placement="right-start">
                            <div slot="content">
                                cookie/header名和cookie/header值用英文冒号隔开，不同组数据用英文分号隔开 <br /> 
                                如 Cookie: lang=zh-CN; PHPSESSID=jvfbua8qs3vb81rapv45d2134q
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                    </el-form-item>
                    <div>
                        <label for="" class="dialog_item_label">fuzz参数</label> 
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline"  >字符型</label>
                            <el-select size="small" 
                                v-model="Addtrafficform.fuzzParam.character"
                                multiple
                                filterable
                                allow-create
                                default-first-option 
                                placeholder=""
                                class="form_item_width"  >
                                <el-option
                                    v-for="(item,i) in characterlist"
                                    :key="i"
                                    :label="item"
                                    :value="item">
                                </el-option>
                            </el-select>
                            <el-tooltip placement="right-start">
                            <div slot="content">
                                对字符串类型的参数名进行fuzz，<br/>如 /vul/sqli/sqli_str.php?name=zhangsan 中的name
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                        </el-form-item>
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline"  >数字型</label>
                            <el-select size="small" 
                                v-model="Addtrafficform.fuzzParam.number"
                                multiple
                                filterable
                                allow-create
                                default-first-option 
                                placeholder=""
                                class="form_item_width"  >
                                <el-option
                                    v-for="(item,i) in numberlist"
                                    :key="i"
                                    :label="item"
                                    :value="item">
                                </el-option>
                            </el-select>
                            <el-tooltip placement="right-start">
                            <div slot="content">
                                对数字类型的参数名进行fuzz，<br/>如 /vul/sqli/sqli_id.php?id=1 中的id
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                        </el-form-item>
                    </div> 
                    <div>
                        <label for="" class="dialog_item_label">fuzz字典</label> 
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline" style="vertical-align: top;  " >字符型</label>
                            <el-input   v-model="Addtrafficform.fuzzDict.character" size="small" 
                                class="form_item_width"  rows="4"  
                                autocomplete="off" type="textarea" resize="none" placeholder=""></el-input>
                            <el-tooltip placement="right-start">
                                <div slot="content">
                                    使用该字典中的字符串进行fuzz，如 /vul/sqli/sqli_str.php?name=zhangsan 中的zhangsan
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <div style="padding-left: 114px;margin-top: 4px;margin-bottom: 4px;">
                            <el-button type="primary" size="mini" style="  margin-right: 12px"
                            @click="clickupload1()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px; ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadID1" ref="upload" @change="changeuploaID1($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline" style="vertical-align: top;  " >数字型</label>
                            <el-input   v-model="Addtrafficform.fuzzDict.number" size="small" 
                                class="form_item_width"  rows="4"  
                                autocomplete="off" type="textarea" resize="none" placeholder=""></el-input>
                            <el-tooltip placement="right-start">
                                <div slot="content">
                                    使用该字典中的字符串进行fuzz，如 /vul/sqli/sqli_id.php?id=1 中的 1
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <div style="padding-left: 114px;margin-top: 4px;">
                            <el-button type="primary" size="mini" style="  margin-right:12px"
                            @click="clickupload2()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px; ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadID2" ref="upload" @change="changeuploaID2($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                    </div>
                    <div>
                        <label for="" class="dialog_item_label" style="margin-bottom: 12px;">响应字典</label> 
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline" style="vertical-align: top;  " >json关键字</label>
                            <el-input   v-model="Addtrafficform.response.jsonKeyword" size="small" 
                                class="form_item_width"  rows="4"  
                                autocomplete="off" type="textarea" resize="none" placeholder="请输入json关键字"></el-input>
                            <el-tooltip placement="right-start">
                                <div slot="content">
                                    作为json响应报文的判断依据，用于提高准确率
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline"  >非json关键字</label>
                            <el-switch
                                v-model="Addtrafficform.response.noJsonSwitch"
                                class="elSwitch"  >
                            </el-switch> 
                        </el-form-item>
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px" v-show="Addtrafficform.response.noJsonSwitch"> 
                            <el-input   v-model="Addtrafficform.response.noJsonKeyword" size="small"  style="margin-left: 116px;"
                                class="form_item_width"  rows="4"  
                                autocomplete="off" type="textarea" resize="none" ></el-input>
                                <el-tooltip placement="right-start">
                                <div slot="content">
                                    作为非json响应报文的判断依据，用于提高准确率
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </el-form-item>
                    </div> 
                </div>
            </el-form>
        </div>

    </div>
</template>
<style lang="less" scoped>
.tastBtnCont {
    font-size: 13px;
}

.tastBtnCont .search-box {
    border-radius:4px;
    padding: 24px;
    min-height: 0;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
.operationbox label  {
    margin-left: 24px;
    color: #4c7ae3;
    font-size: 13px;
}
.topTips {
    color: #4c7ae3;
    width: 40px;
}

.tipsCont {
    float: left;
    /* margin-left:24px; */
    line-height: 31px;
}
.formlist{
    padding: 24px;
    box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
    border-radius: 4px;
}
/deep/ .el-textarea {
    vertical-align:auto;
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
.form_item_width{
    width: 320px;
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
</style>
<script> 
import { traffic } from '@/api/traffic.js';
import { node } from "@/api/system.js";
import BannerBox from "@/components/BannerBox.vue"; 
export default {
    name:'addtraffic',
    components:{
        BannerBox, 
    },
    data(){
        return{
            Addtrafficform:{
                target:'',
                taskname:'',
                stnode:'',
                networkcard:'',
                networkport:'',
                duration:'',
                waitCred_pattern:'',
                waitCred_value:'',
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
                vulname:[],
            },
            rules:{},
            durationlist:[],
            nodelist:[],
            networkcardlist:[],
            userID:'',
            tasknameError:'',
            isShowMore:false,
            credPatternlist:[],
            characterlist:[],
            numberlist:[],
            vulnamelist:[],
        }
    },
    created() {
        this.$store.state.activefirstMenu = "/traffic";
        this.userID = this.$commonjs.decryptCBC(localStorage.getItem('user_id'),this.$commonjs.myKey); 
    },
    mounted() {
        this.getEnum(); 
    },
    methods: {
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
      
    
        async submithandle(){
            let config={
                waitCred:{
                    pattern:Number(this.Addtrafficform.waitCred_pattern),
                    value:this.Addtrafficform.waitCred_value
                },
                fuzzParam:{
                    character:this.Addtrafficform.fuzzParam.character.join(','),
                    number:this.Addtrafficform.fuzzParam.number.join(','),
                }, 
                fuzzDict:{
                    character:this.Addtrafficform.fuzzDict.character,
                    number:this.Addtrafficform.fuzzDict.number,
                },
                response:{
                    jsonKeyword:this.Addtrafficform.response.jsonKeyword,
                    noJsonSwitch:this.Addtrafficform.response.noJsonSwitch,
                    noJsonKeyword:this.Addtrafficform.response.noJsonKeyword
                },
            }
            let param={
                taskName: this.Addtrafficform.taskname,
                networkCard:this.Addtrafficform.networkcard,
                port:this.Addtrafficform.networkport,
                expireTime:this.Addtrafficform.duration,
                targetUrl:this.Addtrafficform.target, 
                userId:Number(this.userID),
                otherConfig:JSON.stringify(config),
                vulConfig:this.Addtrafficform.vulname.join(',')
            };
            const res = await traffic.addTraffic(param);
            if(res.code == 200){
                this.$message({
                    message: '保存成功',
                    type: 'success'
                });
                this.$router.push({
                    path: `/traffic`,
                });
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        targetinput(e) {
            this.Addtrafficform.target =
                this.Addtrafficform.target.replace(/[^\S\r\n]/, "");
            this.Addtrafficform.taskname =
                this.Addtrafficform.target.substr(0, 20) +
                "_" +
                this.commonjs.nowtime();
        },
        clickupload() {
            document.querySelector(".btnUploadID").click();
        },
        changeuploaID(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return;
            // that.taskfile.name = f.name;
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
                    let carData = XLSX.utils.sheet_to_json(
                        that.wb.Sheets[that.wb.SheetNames[0]]
                    );
                    let arr = [];
                    for (var key in carData) {
                        for (var k in carData[key]) {
                            if (arr.indexOf(k) === -1) {
                                arr.push(k);
                            }
                            if (arr.indexOf(carData[key][k]) === -1) {
                                arr.push(carData[key][k]);
                            }
                        }
                    }
                    that.Addtrafficform.target = arr.join("\n");
                    let str = "";
                    arr.forEach((item) => {
                        str += item;
                    });
                    that.Addtrafficform.taskname =
                        str.substring(0, 20) + "_" + that.commonjs.nowtime();
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
                        that.Addtrafficform.target = reader.result;
                        that.Addtrafficform.taskname =
                            reader.result.substring(0, 20) + "_" + that.commonjs.nowtime();
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        showMoreconfig(){
            this.isShowMore  = !this.isShowMore;
        },
        clickupload1(){
            document.querySelector(".btnUploadID1").click();
        },
        changeuploaID1(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1); 
            if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) { 
                        that.Addtrafficform.fuzzDict.character = reader.result; 
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        clickupload2(){
            document.querySelector(".btnUploadID2").click();
        },
        changeuploaID2(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1); 
            if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) { 
                        that.Addtrafficform.fuzzDict.number = reader.result; 
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
    },
}
</script>