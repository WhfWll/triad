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
            @click="btnsave('yak')"
            size="small" style="margin-right:8px;">生成Yak模板</xzbutton> 
            <xzbutton 
            type="primary" 
            @click="btnsave('python')"
            size="small">生成Python模板</xzbutton> 
            <label for="">提示：</label>
            <span>用于生成符合系统要求的漏洞检测脚本的模板文件</span>
        </div>
        <div class="addteportbox clearfix">
            <div class="context_box_bg">
               <el-form :model="vuln_info" status-icon  ref="rulesFrom"  :rules="rules" label-width="0" > 
                    <!-- <el-form-item
                        prop = ''
                        label=" " > 
                        <label class="dialog_item_label" style="margin-right: 30px;">系统已有漏洞</label>
                        <el-radio-group v-model="radio" @change="changeRadio">
                            <el-radio :label="3">已有漏洞</el-radio>
                            <el-radio :label="6">新增漏洞</el-radio>
                        </el-radio-group>
                        <div style="margin:10px 0 0 110px;" v-if="radio == 3">
                            <el-select
                                @change="getvulDetail(value)"
                                v-model="value"
                                filterable
                                remote
                                reserve-keyword
                                placeholder="请输入关键词"
                                :remote-method="remoteMethod"
                                :loading="loading" style="width:812px">
                                <el-option
                                v-for="item in options"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value">
                            </el-option>
                          </el-select>
                        </div>
                    </el-form-item>  -->

                    <div class="twolinebox">
                        <el-form-item
                            prop = 'vul_name'
                            label=" "  class="frontline"> 
                            <label class="dialog_item_label topline">漏洞名称</label>
                            <el-input
                                v-model="vuln_info.vul_name"
                                size="small"
                                style="width:300px"
                                :disabled="update"
                                placeholder="漏洞名称"
                                maxlength="50"  
                            ></el-input>
                        </el-form-item>
                        <!-- <div  > 
						<el-select v-model="vuln_info.risk_level"  style=" width:150px;" clearable placeholder="漏洞风险"  size="small" ref="vulSelect">  
							<el-option
								v-for="(item,i) in vulrisklist"
								:key="i"
								:label="item.label"
								:value="item.value"> 
							</el-option>
						</el-select>  
	                </div> -->
                        <el-form-item
                            prop = 'vul_risk'
                            label=" " style="display: inline-block;"> 
                            <label class="dialog_item_label topline">漏洞风险</label>
                            <el-select v-model="vuln_info.vul_risk"  style=" width:300px;" clearable placeholder="漏洞风险"  size="small" ref="vulSelect">  
                                <el-option
                                    v-for="(item,i) in vulrisklist"
                                    :key="i"
                                    :label="item.label"
                                    :value="item.value"> 
                                </el-option>
                            </el-select>  
                        </el-form-item>
                    </div>
                    
                    <div class="twolinebox">
                        <el-form-item
                            prop = 'vul_class'
                            label=" " class="frontline"> 
                            <label class="dialog_item_label topline">漏洞分类</label>
                            <el-select v-model="vuln_info.vul_class"  style=" width:300px;" clearable placeholder="漏洞分类"  size="small" ref="vulSelect">  
                                <el-option
                                    v-for="(item,i) in objectlist"
                                    :key="i"
                                    :label="item"
                                    :value="item"> 
                                </el-option>
                            </el-select>  
                        </el-form-item>

                        <el-form-item
                            prop = 'vul_type'
                            label=" " style="display: inline-block;"> 
                            <label class="dialog_item_label topline">漏洞类型</label>
                            <el-select v-model="vuln_info.vul_type"  style=" width:300px;" clearable placeholder="漏洞类型"  size="small" ref="vulSelect">  
                                <el-option
                                    v-for="(item,i) in typelist"
                                    :key="i"
                                    :label="item"
                                    :value="item"> 
                                </el-option>
                            </el-select> 
                        </el-form-item>
                    </div>

                    <div class="twolinebox">
                        <el-form-item
                            prop = 'use_impact'
                            label=" " class="frontline"> 
                            <label class="dialog_item_label topline">利用影响</label>
                            <el-select v-model="vuln_info.use_impact"  style=" width:300px;" clearable placeholder="利用影响"  size="small" ref="vulSelect">  
                                <el-option
                                    v-for="(item,i) in threatist"
                                    :key="i"
                                    :label="item"
                                    :value="item"> 
                                </el-option>
                            </el-select> 
                        </el-form-item>

                        <el-form-item
                            prop = 'cvss'
                            label=" " style="display: inline-block;"> 
                            <label class="dialog_item_label topline">CVSS评分</label>
                            <el-input
                                v-model="vuln_info.cvss"
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
                        label=" " prop = 'publish_time' class="twolinebox"> 
                        <label class="dialog_item_label topline">披露时间</label>
                            <div class="block" style="display: inline-block;">
                                <el-date-picker
                                v-model="vuln_info.publish_time"
                                value-format="yyyy-MM-dd"  
                                format="yyyy-MM-dd"
                                type="date"
                                placeholder="选择日期">
                                </el-date-picker>
                            </div>
                        </el-form-item>
                    </div>


                  <el-form-item
                    prop = 'vul_description'
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label">漏洞描述</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="3"
                            v-model="vuln_info.vul_description"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="1000"  
                        ></el-input>
                    </el-form-item>

                    <el-form-item
                    prop = 'repair_suggest'
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label">修复建议</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="3"
                            v-model="vuln_info.repair_suggest"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item>  
                    <!-- //新增 -->
                    <el-form-item
                    prop = 'vul_num'
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label" >漏洞编号</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="3"
                            v-model="vuln_info.vul_num"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item> 

                    <el-form-item
                    prop = 'affect_range'
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label">影响范围</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="1"
                            v-model="vuln_info.affect_range "
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item> 
                    <!-- //新增 -->
                    <el-form-item
                    prop = 'vul_analysis'
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label">漏洞分析说明</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="3"
                            v-model="vuln_info.vul_analysis"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item>  

                    <el-form-item
                    prop = 'reference_link'
                    label=" " style="display:block;"> 
                    <label class="dialog_item_label">参考链接</label>
                        <el-input
                            style="display:block;width:935px;"
                            type="textarea"
                            :rows="1"
                            v-model="vuln_info.reference_link"
                            size="small"
                            :disabled="update"
                            placeholder=""
                            maxlength="50"  
                        ></el-input>
                    </el-form-item>  
                    <el-form-item
                    prop = 'params'
                    label=" " class="lastline"> 
                    <label class="dialog_item_label bottomerror">脚本输入参数</label>
                    <el-checkbox-group v-model="vuln_info.params" style="display: inline-block;">
                        <el-checkbox label="漏洞地址"></el-checkbox>
                        <el-checkbox label="监听地址"></el-checkbox>
                        <el-checkbox label="账号字典"></el-checkbox>
                        <el-checkbox label="密码字典"></el-checkbox>
                        <el-checkbox label="资源文件"></el-checkbox>
                    </el-checkbox-group>
                    </el-form-item> 
                </el-form>
            </div>

        </div>
  	</div>
</template>
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
            vulrisklist:[],
            objectlist:[],
            typelist:[],
            threatist:[],
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
    	  	vuln_info:{
                publish_time:'',
                vul_name:'',
                vul_risk:'',
                vul_num:'',
                reference_link:'',
                vul_class:'',
                use_impact:'',
                vul_type:'',
                cvss:'',
                affect_range:'',
                vul_description:'',
                repair_suggest:'',
                vul_analysis:'',
                params:[],

			},
            select_value:'',
            rules:{
                vul_name:[
                    { required: true, message: '漏洞名称不能为空', trigger: 'blur' }, 
                ],
                vul_risk:[
                    { required: true, message: '漏洞风险不能为空', trigger: 'blur' }, 
                ],
                vul_class:[
                    { required: true, message: '漏洞分类不能为空', trigger: 'blur' }, 
                ],
                vul_type:[
                    { required: true, message: '漏洞类型不能为空', trigger: 'blur' }, 
                ],
                use_impact:[
                    { required: true, message: '利用影响不能为空', trigger: 'blur' }, 
                ],
                vul_description:[
                    { required: true, message: '漏洞描述不能为空', trigger: 'blur' }, 
                ],
                repair_suggest:[
                    { required: true, message: '修复建议不能为空', trigger: 'blur' }, 
                ],
                params:[
                    { required: true, message: '至少选择一项脚本输入参数', trigger: 'blur' }, 
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
        this.getParams(); 
    },
    methods:{ 
        getParams(){
        this.$ajax.get('/tools/vul/params/',{
                params: {
                }
            })
            .then((res) => { 
                var dt = res.data; 
                if(dt.success){
                  this.vulrisklist = dt.data.risk;
                  this.objectlist = dt.data.class;
                  this.typelist = dt.data.type;
                  this.threatist = dt.data.impact;
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
    btnsave(language){ //保存
        this.$refs.rulesFrom.validate(valid => {
            if (valid) {
                this.$ajax({
                method:'post',
                url:'/tools/vul/script_download/',
                contentType: 'application/json',
                data: {
                    language: language,
                    vul_name:this.vuln_info.vul_name,
                    publish_time:this.vuln_info.publish_time,
                    vul_risk:this.vuln_info.vul_risk,
                    vul_num:this.vuln_info.vul_num,
                    reference_link:this.vuln_info.reference_link,
                    vul_class:this.vuln_info.vul_class,
                    use_impact:this.vuln_info.use_impact,
                    vul_type:this.vuln_info.vul_type,
                    cvss:this.vuln_info.cvss,
                    affect_range:this.vuln_info.affect_range,
                    vul_description:this.vuln_info.vul_description,
                    repair_suggest:this.vuln_info.repair_suggest,
                    vul_analysis:this.vuln_info.vul_analysis,
                    params: JSON.stringify({
                        target: this.vuln_info.params.includes('漏洞地址'),
                        listen_address: this.vuln_info.params.includes('监听地址'),
                        user_dict: this.vuln_info.params.includes('账号字典'),
                        password_dict: this.vuln_info.params.includes('密码字典'),
                        resource_file: this.vuln_info.params.includes('资源文件'),
                    })
                }
            })
            .then(res => {
                if (res.data.error) {
                    this.$message.error(res.data.error)
                } else {
                    const blob = new Blob([res.data],{
                        type:"application/octet-stream"
                    });
                    let filename = sessionStorage.filename;
                    // alert(filename)
                    const url =  window.URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.download = filename;
                    a.href = url;
                    a.click();
                    sessionStorage.filename = ''
                }
                
        
            })
            .catch(error=>{ });
            }
        })
        },
    //     remoteMethod(query) {
    //     if (query !== '') {
    //       this.loading = true;
    //       setTimeout(() => {
    //         this.$ajax.get('/tools/newTool/associateVuln/',{
    //             params: {
    //                 keyword: query
    //             }
    //         })
    //         .then((res) => { 
    //             var dt = res.data; 
    //             if(dt.success){
    //                 this.loading = false;
    //                 this.options = []
    //                 dt.vuln_list.forEach(item => {
    //                     this.options.push({
    //                         label: item, 
    //                         value: item
    //                     })
    //                 })
    //             }else{
    //                 this.$message({
    //                     message:dt.msg,
    //                     type: 'error'
    //                 });
    //             }  
    //         })
    //         .catch((error) => {
    //             console.log(error);
    //         })
    //       }, 200);
    //     } else {
    //       this.options = [];
    //     }
    //   },
    //   系统已有漏洞远程搜索实现
    //   getvulDetail(value){
    //     this.$ajax.get('/tools/newTool/vulnDetial/',{
    //             params: {
    //                 name: value
    //             }
    //         })
    //         .then((res) => { 
    //             var dt = res.data; 
    //             if(dt.success){
    //                 this.vuln_info = dt.vuln_info
    //                 this.select_value = dt.tool_info.select_value
    //                 this.radio2 = Number(dt.tool_info.select_mode)
    //             }else{
    //                 this.$message({
    //                     message:dt.msg,
    //                     type: 'error'
    //                 });
    //             }  
    //         })
    //         .catch((error) => {
    //             console.log(error);
    //         })
    //   },
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

<!-- 生成脚本模板样式 -->
<style scoped lang="less">
    .frontline{
        display: inline-block;
        margin-right:100px!important;
    }
    .topline{
        margin-right: 30px;
    }
    .el-date-editor.el-input, .el-date-editor.el-input__inner {
    width: 300px !important;
    }
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
        width:78px;
    }
    .bottomerror {
        margin-bottom: 15px;
        display: inline-block;
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
            left: 112px;
        }
        /deep/ .lastline .el-form-item__label  {
            position: absolute;
            left: 140px;
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
             border-radius: 4px;
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
                 border-radius: 4px;
                .context_box_bg{
                    padding: 24px 24px!important;
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
