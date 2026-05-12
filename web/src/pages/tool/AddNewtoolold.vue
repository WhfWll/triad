<template>
    <!-- 漏洞库工具模板弹窗主体代码 -->
  	<div class="mainbox">  
        <div class="main-title  ">   
            <router-link :underline="false" class="classA" :to="{ path: '/vulnerability' }" >漏洞库</router-link>  
            <label class="currentpagetitle">
                <span>{{title}}</span> 
            </label>
	  	</div> 
          <!-- 保存模板按钮 -->
        <div class="operationbox whitebg context_box_bg">
            <!-- <el-button type="primary"  size="small"  @click='btnsave()'>生成模板</el-button> -->
            <xzbutton 
            type="primary" 
            @click="btnsave()"
            size="small" style="margin-right:8px;">生成Yak模板</xzbutton> 
            <xzbutton 
            type="primary" 
            @click="btnsave()"
            size="small"  >生成Python模板</xzbutton> 
            <label for="">提示：</label>
            <span>用于生成符合系统要求的漏洞检测脚本的模板文件</span>
        </div>
        <div class="addteportbox clearfix">
            <div class="context_box_bg">
               <el-form :model="vuln_info" status-icon  ref="rulesFrom"  :rules="rules" label-width="0" > 
                    <el-form-item
                        prop = ''
                        label=" " > 
                        <label class="dialog_item_label" style="margin-right: 30px;">系统已有漏洞</label>
                        <el-radio-group v-model="radio" @change="changeRadio">
                            <el-radio :label="3">已有漏洞</el-radio>
                            <el-radio :label="6">新增漏洞</el-radio>
                        </el-radio-group>
                        <div style="margin:10px 0 0 103px;" v-if="radio == 3">
                            <el-select
                                @change="getvulDetail(value)"
                                v-model="value"
                                filterable
                                remote
                                reserve-keyword
                                placeholder="请输入关键词"
                                :remote-method="remoteMethod"
                                :loading="loading" style="width:806px">
                                <el-option
                                v-for="item in options"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value">
                            </el-option>
                          </el-select>
                        </div>
                    </el-form-item> 

                    <div class="twolinebox">
                        <el-form-item
                            prop = 'name'
                            label=" " style="display: inline-block;margin-right:100px;"> 
                            <label class="dialog_item_label" style="margin-right:30px;">漏洞名称</label>
                            <el-input
                                v-model="vuln_info.name"
                                size="small"
                                style="width:300px"
                                :disabled="update"
                                placeholder=""
                                maxlength="50"  
                            ></el-input>
                        </el-form-item>
                        <el-form-item
                            prop = 'risk_level'
                            label=" " style="display: inline-block;"> 
                            <label class="dialog_item_label" style="margin-right:30px;">漏洞风险</label>
                            <el-input
                                v-model="vuln_info.risk_level"
                                size="small"
                                style="width:300px"
                                :disabled="update"
                                placeholder=""
                                maxlength="50"  
                            ></el-input>
                        </el-form-item>
                    </div>
                    
                    <div class="twolinebox">
                        <el-form-item
                            prop = 'object'
                            label=" " style="display: inline-block;margin-right:100px;"> 
                            <label class="dialog_item_label" style="margin-right:30px;">漏洞对象</label>
                            <el-input
                                v-model="vuln_info.object"
                                size="small"
                                style="width:300px"
                                :disabled="update"
                                placeholder=""
                                maxlength="50"  
                            ></el-input>
                        </el-form-item>

                        <el-form-item
                            prop = 'type'
                            label=" " style="display: inline-block;"> 
                            <label class="dialog_item_label" style="margin-right:30px;">漏洞类型</label>
                            <el-input
                                v-model="vuln_info.type"
                                size="small"
                                style="width:300px"
                                :disabled="update"
                                placeholder=""
                                maxlength="50"  
                            ></el-input>
                        </el-form-item>
                    </div>

                    <div class="twolinebox">
                        <el-form-item
                            prop = 'threat'
                            label=" " style="display: inline-block;margin-right:100px;"> 
                            <label class="dialog_item_label" style="margin-right:30px;">威胁程度</label>
                            <el-input
                                v-model="vuln_info.threat"
                                size="small"
                                style="width:300px"
                                :disabled="update"
                                placeholder=""
                                maxlength="50"  
                            ></el-input>
                        </el-form-item>

                        <el-form-item
                            prop = ''
                            label=" " style="display: inline-block;"> 
                            <label class="dialog_item_label" style="margin-right:30px;">VUL_ID</label>
                            <el-input
                                v-model="vuln_info.vulID"
                                size="small"
                                style="width:300px"
                                :disabled="update"
                                placeholder=""
                                maxlength="50"  
                            ></el-input>
                        </el-form-item> 
                    </div>
                    <div class="twolinebox">
                        <el-form-item
                        label=" " prop = 'find_time' class="twolinebox"> 
                        <label class="dialog_item_label" style="margin-right: 30px;">公开时间</label>
                            <div class="block" style="display: inline-block;">
                                <!-- <span class="demonstration">公开时间</span> -->
                                <el-date-picker
                                v-model="vuln_info.find_time"
                                type="date"
                                placeholder="选择日期">
                                </el-date-picker>
                            </div>
                        </el-form-item>
                    </div>


                  <el-form-item
                    prop = ''
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label" style="margin-right:30px;">漏洞描述</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="5"
                            v-model="vuln_info.desc"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="1000"  
                        ></el-input>
                    </el-form-item>

                    <el-form-item
                    prop = ''
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label" style="margin-right:30px;">修复建议</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="3"
                            v-model="vuln_info.solution"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item>  
                
                    <el-form-item
                    prop = ''
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label" style="margin-right:30px;">影响范围</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="1"
                            v-model="vuln_info.affect_version "
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item> 

                    <el-form-item
                    prop = ''
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label" style="margin-right:30px;">参考链接</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="1"
                            v-model="vuln_info.ref"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item>  
<!-- //工具触发参数 是新版添加 没有接口 -->
                    <!-- <el-form-item
                    prop = ''
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label" style="margin-right:30px;">工具触发参数</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="1"
                            v-model="vuln_info.ref"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item>  -->

                    <el-form-item
                    prop = 'toolconfig'
                    label=" " > 
                    <label class="dialog_item_label" style="margin-bottom: 15px;display: block;">工具配置</label>
                    <label class="dialog_item_label" style="margin-right: 30px;border-left: none;">漏洞识别模式</label>
                    <el-radio-group v-model="radio2">
                        <el-radio :label="1">服务</el-radio>
                        <el-radio :label="2">Web指纹</el-radio>
                        <el-radio :label="3">设备指纹</el-radio>
                      </el-radio-group>
                      <div style="margin-left: 15px;display: inline-block;">
                        <el-input
                                v-model="select_value"
                                size="small"
                                style="width:487px"
                                :disabled="update"
                                placeholder=""
                                maxlength="50"  
                        ></el-input> 
                    </div>
                    </el-form-item> 
<!-- //脚本参数是新版添加 没有接口 -->
                    <!-- <el-form-item
                    prop = 'loopmodel'
                    label=" " > 
                    <label class="dialog_item_label" style="margin-right: 30px;border-left: none;">脚本参数</label>
                    <el-input
                        v-model="vuln_info.ref"
                        size="small"
                        style="width:300px"
                        :disabled="update"
                        placeholder=""
                        maxlength="50"  
                     ></el-input> 
                    </el-form-item> -->
                </el-form>
            </div>

        </div>
  	</div>
</template>
<!-- 生成脚本模板样式 -->
<style scoped lang="less">
/deep/ .twolinebox .el-form-item{
    position: relative;
}
 /deep/ .twolinebox .el-form-item__label {
    position: absolute;
    left: 75px;
}   
/deep/ .el-textarea__inner{
    padding: 10px 5px;
}
.twolinebox  .dialog_item_label {
    width:75px;
}
/deep/  .el-radio {
color: rgba(72, 72, 102, 0.64);
cursor: pointer;
margin-right: 30px;
}
/deep/ .el-input__icon {
    height: 100%;
    width: 25px;
    text-align: center;
    transition: all .3s;
    line-height: 32px;
}
/deep/ .el-form-item__content{
        line-height: 32px;;
    }
    .catatree{
        font-size: 13px;
        font-weight: 400;
        color:rgba(72, 72, 102, 0.64)
    }
    /deep/ .el-form-item__label {
        position: absolute;
        left: 140px;
    }
    /deep/ .el-form-item__error {
        left: 112px;
    }
    .dialog_item_label{
        font-size: 14px;
        color: rgba(72, 72, 102, 0.87);
        font-weight:700;
    }
    .search-box .search-text .input-with-select{
        vertical-align: inherit;
    }
    .search-box{
        margin-bottom: 0;
        padding-bottom: 0;
    }
    .operationbox{
        padding: 24px;
        margin-bottom: 15px;
         box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
        label{
            margin-left: 24px;
            color: #4C7AE3;
            font-size: 13px;
        }
        span{
            margin-left: 4px;
            font-size: 13px;
            color: rgba(72, 72, 102, 0.64);
        }
    }
    .serach-condition > div{
        margin-bottom: 16px;
    }
    .dialog_b_btn{
        position: absolute;
        top: 15px;
        right: 24px;
        font-size: 14px;
        button{
            color: #4C7AE3;
        } 
    }
     @media (max-width: 1440px) {
         
        /deep/ .el-dialog{
            height: calc(100% - 96px);
        }
    }
    @media  (min-width: 1440px) { 
        /deep/ .el-dialog{
            height: calc(100% - 176px);
        }
    }
     .addexperiencebox{  
        height:calc(100% - 130px);
        box-sizing: border-box;
        margin-bottom: 30px;
    }	
 .target_active{
        .msg_icon{ 
            background: rgba(255, 255, 255, 0.3);
            color: #fff;
        }
        .msg_label{
            color: rgba(255, 255, 255, 0.7) !important;
        }
        .msg_count{
            label{
                color: #fff !important;
            }
            i{
                color: #fff !important;
            } 
        } 
        .icontsstyle{
            color: rgba(255, 255, 255, 0.7) !important;;
        }
    }
    .is-required{
        margin-right:4px;
        color:#F56C6C;
        font-size: 12px; 
    } 
    .mainbox{
        display: flex;
        flex-direction: column;
        .operationbox{
            height: 78px;
            box-sizing: border-box;
        }
         .addteportbox{
             flex: 1;
             box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
            .context_box_bg{
                padding: 31px 24px!important;
                height: 100%;
                box-sizing: border-box;
                /deep/ .el-form{
                    // .el-form-item__label{
                    //     text-align: left;
                    //     height: 18px;
                    //     line-height: 18px;
                    //     margin-top:12px;
                    //     padding-left: 6px;
                    //     border-left: 3px solid #4C7AE3;
                    // }
                    .speciallabel .el-form-item__label{
                        margin-top:5px; 
                    }
                    .fengmian{
                        .el-form-item__content{
                            /* height: 140px; */
                            /* display: flex;
                            align-items: flex-end; */
                            .form-img{
                                width: 100px;
                                height: 140px;
                            }
                            .changeBgBtn{
                                /* margin: 0 16px; */
                                 width: 80px;
                                height: 32px;
                                text-align: center;
                                line-height: 7px;
                                padding: 0px;
                                border-radius: 2px;
                            }
                        }
                        
                    }
                    
                }
            }
        }
    }
   
</style>
<!-- 生成脚本模板功能实现 -->
<script>  
import xzbutton from "@/components/XzButton.vue";
export default({
    name:'addnewtool',
    components: {
    	xzbutton,
  	},
    data(){  
    	return{
            showerr:false,
            options: [],
            value: [],
            list: [],
            loading: false,
            states: [],
            title:this.$route.query.title,
            isUpdate:this.$route.query.isUpdate == 'true'? true:false, //编辑模板
            template_id:this.$route.query.template_id,  //模板ID
            radio: 3,
            radio2: 1,
    	  	vuln_info:{
                find_time:'',
                name:'',
                risk_level:'',
                solution:'',
                ref:'',
                object:'',
                threat:'',
                type:'',
                vulID:'',
                affect_version:'',
                desc:''
			},
            select_value:'',
            rules:{
                // taskname:[
                //     { required: true, message: '系统已有漏洞不能为空', trigger: 'blur' }, 
                // ],
                name:[
                    { required: true, message: '漏洞名称不能为空', trigger: 'blur' }, 
                ],
                risk_level:[
                    { required: true, message: '漏洞风险不能为空', trigger: 'blur' }, 
                ],
                object:[
                    { required: true, message: '漏洞对象不能为空', trigger: 'blur' }, 
                ],
                type:[
                    { required: true, message: '漏洞类型不能为空', trigger: 'blur' }, 
                ],
                threat:[
                    { required: true, message: '威胁程度不能为空', trigger: 'blur' }, 
                ],
                find_time:[
                    { required: true, message: '时间不能为空', trigger: 'blur' }, 
                ],
                vulID:[
                    { required: true, message: 'VUL_ID不能为空', trigger: 'blur' }, 
                ],
            }, 
            Loading:false,
            update:false,
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/vulnerability"; 
        this.pageSize = this.commonjs.pageSize; 
 
    },
    mounted:function(){   
        this.list = this.states.map(item => {
            return { value: `value:${item}`, label: `label:${item}` };
        })
        
    },
    methods:{ 
        remoteMethod(query) {
        if (query !== '') {
          this.loading = true;
          setTimeout(() => {
            this.$ajax.get('/tools/newTool/associateVuln/',{
                params: {
                    keyword: query
                }
            })
            .then((res) => { 
                var dt = res.data; 
                if(dt.success){
                    this.loading = false;
                    this.options = []
                    dt.vuln_list.forEach(item => {
                        this.options.push({
                            label: item, 
                            value: item
                        })
                    })
                }else{
                    this.$message({
                        message:dt.msg,
                        type: 'error'
                    });
                }  
            })
            .catch((error) => {
                console.log(error);
            })
          }, 200);
        } else {
          this.options = [];
        }
      },
    //   系统已有漏洞远程搜索实现
      getvulDetail(value){
        this.$ajax.get('/tools/newTool/vulnDetial/',{
                params: {
                    name: value
                }
            })
            .then((res) => { 
                var dt = res.data; 
                if(dt.success){
                    this.vuln_info = dt.vuln_info
                    this.select_value = dt.tool_info.select_value
                    this.radio2 = Number(dt.tool_info.select_mode)
                }else{
                    this.$message({
                        message:dt.msg,
                        type: 'error'
                    });
                }  
            })
            .catch((error) => {
                console.log(error);
            })
      },
    //   生成脚本模板保存功能实现
    btnsave(){ //保存
            // this.$refs.rulesFrom.validate((valid) => {
            //     if (valid) {
                    
            this.$ajax({
                method:'post',
                url:'/tools/newTool/generateTemplates/',
                contentType: 'application/json',
                data: {
                    vuln_info: this.vuln_info,
                    tool_info:{
                        select_mode:this.radio2,
                        select_value:this.select_value
                    }
                }
            })
            .then(res => { 
                var dt = res.data; 
                let that = this;  
                if(dt.success){ 
                    const blob = new Blob([dt.file_content],{
                        type:"application/octet-stream"
                    });
                    let filename = dt.filename
                    const url =  window.URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.download = filename;
                    a.href = url;
                    a.click();                      
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }
        
            })
            .catch(error=>{ });
            //     }
            // });
        },
        clickupload(){  
            document.querySelector('.btnUploadID').click();
        },
        disabledFn(){  //禁用节点
            return this.update;
        },
        btnUpdateTemplate(){ //点击编辑模板
            this.update = false; 
        },
        //单选框切换
        changeRadio () {
                this.vuln_info = {}
                this.select_value =''
        }
        
    }
})
 
</script>
