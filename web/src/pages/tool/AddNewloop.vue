<template>
    <!-- 漏洞库新建漏洞弹窗代码 -->
  	<div class="mainbox">  
        <div class="main-title  ">   
            <router-link :underline="false" class="classA" :to="{ path: '/vulnerability' }" >漏洞库</router-link>  
            <label class="currentpagetitle">
                <span>{{title}}</span> 
            </label>
	  	</div> 
          
        <div class="operationbox whitebg context_box_bg">
            <!-- <el-button type="primary" size="small" @click="btnUpdateTemplate">生成模板</el-button>   -->
            <!-- <el-button type="primary"  size="small"  @click="submitform"  v-if="!update">保存漏洞</el-button> -->
            <xzbutton 
            type="primary" 
            @click="submitform"
            v-if="!update"  
            size="small">新建漏洞</xzbutton> 
            <label for="">提示：</label>
            <span>配置漏洞工具调用的相关条件，上传工具文件，即可将漏洞工具添加进行漏洞库</span>
        </div>
        <div class="addteportbox clearfix">
            <div class="context_box_bg">
               <el-form :model="taskfrom" status-icon  ref="rulesFrom"  :rules="rules" label-width="0" > 
                <!-- 系统已有漏洞部分 -->
                    <el-form-item
                        prop ='radio'
                        label=" " > 
                        <label class="dialog_item_label" style="margin-right: 30px;">工具调用条件</label>
                        <el-radio-group v-model="taskfrom.radio"  @change="changeRadio">
                            <el-radio label="service">服务</el-radio>
                            <el-radio label="web"> Web指纹</el-radio>
                        </el-radio-group>
                        <!-- v-if="taskfrom.radio == 1" -->
                        <div style="margin:10px 0 0 145px;">
                            <el-select
                                @change="getvulDetail(value)"
                                v-model="taskfrom.vul_relate"
                                filterable
                                remote
                                reserve-keyword
                                placeholder="请选择/输入漏洞关联服务（请选择/输入漏洞关联WEB指纹） "
                                :remote-method="remoteMethod"
                                :loading="loading" style="width:806px">
                                <el-option
                                v-for="item in options"
                                :key="item"
                                :label="item"
                                :value="item">
                            </el-option>
                          </el-select>
                        </div>
                    </el-form-item> 
                    <!-- 漏洞工具类型部分 -->
                    <el-form-item
                    prop ='radio2'
                    label=" " > 
                    <label class="dialog_item_label" style="margin-right: 30px;">漏洞工具类型</label>
                    <el-radio-group v-model="taskfrom.radio2">
                        <el-radio label="POC">POC</el-radio>
                        <el-radio label="EXP">EXP</el-radio>
                      </el-radio-group>
                    </el-form-item>

                    <el-form-item
                    prop ='radio3'
                    label=" " > 
                    <label class="dialog_item_label" style="margin-right: 30px;">工具编写语言</label>
                    <el-radio-group v-model="taskfrom.radio3">
                        <el-radio label="Python">Python</el-radio>
                        <el-radio label="Yak">Yak</el-radio>
                      </el-radio-group>
                    </el-form-item>
                    <!-- 上传工具文件部分  -->
                    <el-form-item prop="newUploadFile"  label="" class="fengmian">
                        <div class="uploadbox">
                            <div><label class="dialog_item_label" style="margin-right: 30px;">上传工具文件</label></div>
                            <!-- <el-button type="primary"  size="small" style="vertical-align: top;margin-right:27px"  @click="importFileClick()">上传文件</el-button> -->
                            <xzbutton 
                            type="primary" 
                            @click="importFileClick()"
                            style="vertical-align: top;margin-right:27px"  
                            size="small"  >上传文件</xzbutton> 
                            <input
                                type="file"
                                class="fileInput"
                                ref="upload"
                                @change="changeFiles($event)"
                                style="display:none"
                                id="fileUpload"/> 
                            <span v-for="(item, index) in fileList" :key="index" class="file">{{item.name}}</span>
                        </div>
                        <!-- <div v-for="(item,index) in fileList" :key="index" :class="['progressCont',{'default':item.success},{'errorProgre':item.success == false}]">
                            <div class="fileCont">
                                <i class="iconfont iconfujian"></i>
                                <span class="fileNameText" >{{ item.name }}</span> -->
                                <!-- <i class="iconfont iconquxiao" @click="clearFiles"></i> -->
                            <!-- </div> -->
                            <!-- <el-progress :stroke-width="2" :percentage="item.percentage" :format="progressText"></el-progress>
                            <span class="uploadStatus successColor"  v-if="item.success">上传成功!</span>
                            <span class="uploadStatus errorColor"  e-else>{{ item.error }}</span> -->
                        <!-- </div> -->
                    </el-form-item>
                </el-form>
            </div>

        </div>
  	</div>
</template>

<script>  
import xzbutton from "@/components/XzButton.vue";
export default({
    name:'addnewloop',
    components: {
    	xzbutton,
  	},
    data(){  
    	return{ 
            options: [],
            value: [],
            list: [],
            uploadFileName: '',
            fileList:[],
            fileStatusData:[],
            repeatNames:[],
            percentage: -1,
            title:this.$route.query.title,
            taskfrom:{
                radio: 'service',
                radio2: 'POC',
                radio3: 'Python',
                vul_relate: '',
                newUploadFile: []
			},
            rules:{
                radio:[
                    { required: true, message: '已有漏洞不能为空', trigger: 'change' }, 
                ],
                radio2:[
                    { required: true, message: '漏洞工具类型不能为空', trigger: 'change' }, 
                ],
                radio3:[
                    { required: true, message: '漏洞工具类型不能为空', trigger: 'change' }, 
                ],
                newUploadFile:[
                    { required: true, message: '已有漏洞不能为空', trigger: 'change' }, 
                ],
                
            }, 
            // isUpdate:this.$route.query.isUpdate == 'true'? true:false, //编辑模板
            // template_id:this.$route.query.template_id,  //模板ID
            
            loading:false,
            update:false,
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/vulnerability"; 
        this.pageSize = this.commonjs.pageSize; 
 
    },
    mounted:function(){   
        // this.list = this.states.map(item => {
        //     return { value: `value:${item}`, label: `label:${item}` };
        // })       
    },
    methods:{
        // 系统已有漏洞按钮切换
      changeRadio () {
                this.vuln_info = {}
                this.select_value =''
        },
        remoteMethod(query) {
        if (query !== '') {
          this.loading = true;
          setTimeout(() => {
            this.$ajax.get('/tools/vul/relation_params/',{
                params: {
                    type: this.taskfrom.radio,
                    search: query
                }
            })
            .then((res) => { 
                var dt = res.data; 
                if(dt.success){
                    this.loading = false;
                    this.options = dt.data
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
        getvulDetail(value){
        // this.$ajax.get('/tools/vul/relation_params/',{
        //         params: {
        //             type: '',
        //             search: value
        //         }
        //     })
        //     .then((res) => { 
        //         var dt = res.data; 
        //         if(dt.success){
        //             this.vuln_info = dt.vuln_info
        //             this.select_value = dt.tool_info.select_value
        //             this.radio2 = Number(dt.tool_info.select_mode)
        //         }else{
        //             this.$message({
        //                 message:dt.msg,
        //                 type: 'error'
        //             });
        //         }  
        //     })
        //     .catch((error) => {
        //         console.log(error);
        //     })
      },

              // 导入文件上传功能实现
      importFileClick() {
            document.querySelector(".fileInput").click();
        },
        changeFiles(e){
            console.log('ee', e)
            this.repeatNames = []
			let _files = e.target.files;
            console.log('e.target.files', e.target.files)
			// let _arr = [];
			// this.taskfrom.newUploadFile = [];
            this.fileList = []
			_files.forEach((item,index) =>{
				const fileType = item.name.substring(item.name.lastIndexOf(".") + 1);
				let isFileType = false ;
				if(fileType == 'yak' || fileType == 'python'){
					isFileType = true;
				}
				const isLt100M = item.size / 1024 / 1024 < 100;

				if (!isFileType) {
					this.$message.error('上传文件只能是python/yak格式!');
					return false;
				}
				// if (!isLt100M) {
				// 	this.$message.error('上传文件大小不能超过 100MB!');
				// 	return false;
				// }
				
				// _arr.push(item)
                this.fileList.push(item);
			});
            
			// _arr.forEach((i,n) =>{
			// 	if(this.repeatNames.indexOf(i.name) == -1){
			// 		let _len = this.fileList.length;

			// 		if(_len>=20){
			// 			this.$message.warning(`上传文件个数不能超过20`);//(`当前限制选择 20 个文件，本次选择了`+_len+`个文件`);
			// 			var file = document.getElementById('fileUpload');
			// 			file.value = ''; //file的value值只能设置为空字符串
			// 			return false;
			// 		}else{
			// 			this.repeatNames.push(i.name);
			// 			this.fileList.push(i);
			// 			this.fileStatusData.push(i); 

			// 			this.taskfrom.newUploadFile.push(i)
			// 		}
            //         this.uploadFileName = this.taskfrom.newUploadFile[0].name
            //         e.target.value = ''
			// 	}
			// })
			// this.importReoprtForm.reports = this.fileStatusData;
		},
        // 显示上传进度
        // progressText(percentage) {
		// 	return percentage === 100 ? '' : `${percentage}%`;
		// },
        // 清空文件
        clearFiles () {
            this.percentage = -1
            this.fileStatusData = []
            this.uploadFileName = ''
            this.repeatNames = []
        },
        //从文件列表移除文件
        // handleRemove(file, fileList) {
        //     console.log(file, fileList);
        // },
        // handlePreview(file) {
        //     console.log(file);
        // },

        submitform:function(){ // 确定上传功能实现
            let _formData = new FormData();
            this.fileList.forEach(function (file) {// 因为要上传多个文件，所以需要遍历
                _formData.append('file', file); 
            })
            // _formData.append('type', this.taskfrom.radio);
            _formData.append('vul_relation', this.taskfrom.vul_relate); 
            _formData.append('verify_type', this.taskfrom.radio2); 
            // _formData.append('language', this.taskfrom.radio3); 
            let config = {
                headers:{
                    'Content-Type':'multipart/form-data'
                }
            };
            this.$ajax.post(
                '/tools/vul/create/',
                _formData,
                 config
            ).
            then((res) =>{
                let dt = res.data;
                if(dt.success){
                    this.fileStatusData.forEach((item,index) =>{
                        item.percentage = 100;
                        item.error = '';
                        item.success = true
                    })
                    this.$message.success(dt.msg)
                    this.$router.push({
                        path: `/vulnerability`,
                        query: { 
                        }
                    });
                } else {
                    this.fileStatusData.forEach((item,index) =>{
                        item.error = '上传失败';
                        item.success = false
                        item.percentage = 80;
                    })
                    // this.$router.push({
                    //     path: `/vulnerability`,
                    //     query: { 
                    //     }
                    // });
                    this.$forceUpdate();
                    // this.percentage = 0
                    this.$message.error(dt.error)
                }
            }).catch((error) =>{
                console.log(error)
            })
        },    
        
    }
})
//  新建漏洞弹窗样式部分
</script>
<style scoped lang="less">
    /deep/  .el-radio {
     color: rgba(72, 72, 102, 0.64);
     cursor: pointer;
     margin-right: 30px;
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
                             flex-direction: column;
                         }
                         .tishiword{
                             font-size: 14px;
                             color: rgba(72, 72, 102, 0.87);
                         }
                         .el-upload-list__item:first-child {
                             margin-top: 20px;
                         }
                         .el-form-item__content{
                             display: flex;
                             align-items: flex-start;
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
     .uploadbox{
        display: flex;
        justify-content: left;
        align-items: center;
        .file{
            color: rgba(72, 72, 102, 0.64);
        }
    }
    .upload-demo{
        display: inline-block;
    }
    .explainword{
        position: absolute;
        top: 22px;
        left: 120px;
        font-size: 13px;
        font-family: SourceHanSansCN-Medium, SourceHanSansCN;
        font-weight: 500;
        color:#fff;
    }
    .uploadFileDiv{
        margin-left: 113px;
        margin-top: 24px;
        .uploadFileName{
            width: 49%;
            /* display:flex; */
            justify-content:space-between;
            .successStatus{
                color: #BBBBC8;
                width:307px;
            }
            .failStatus{
                width:307px;
                color: #F35F28;
            }
            .el-icon-close{
                padding-left: 10px;
                cursor: pointer;
            }
        }
    }
    .progressCont{
        width: 320px;
        margin-left: 20px;
        padding-top:16px;
        line-height:18px!important;
        position:relative;
        .uploadStatus{
            position: absolute;
            top:26px;
            right:-38px;
        }
    }
    /deep/.el-progress__text{
        display:none;
    }
    /deep/.el-upload-list{
        display:none;
    }
    .default /deep/.el-progress-bar__inner{
        background-color: rgba(76, 122, 227, 1);//rgba(243, 95, 40, 1);
    }
    .default /deep/span{
        color:rgba(76, 122, 227, 1)!important;
    }
    .errorProgre /deep/.el-progress-bar__inner{
        background-color: rgba(243, 95, 40, 1);
    }
    .errorProgre /deep/span{
        color: rgba(243, 95, 40, 1)!important;
    }
    .successColor{
        color:rgba(76, 122, 227, 1);
    }
    .errorColor{
        color:rgba(243, 95, 40, 1)!important;
    }
    .fileCont{
        height:18px;
        color:rgba(72, 72, 102, 0.64);
        font-size:13px;
        // overflow:hidden;
        .fileNameText{
            padding-left:6px;
            display:block;
            width:220px;	
            height:18px;
            line-height:18px;
            overflow:hidden;
            text-overflow:ellipsis;
            white-space:nowrap;
            float:left;
            color:rgba(76, 122, 227, 1);
            // color:rgba(243, 95, 40, 1);
        }
        i{
            display:block;
            float:left;
        }
        .iconquxiao{
            margin-left:16px;
            cursor:pointer;
            // z-index: 99;
        }
    }
    /deep/ .failProgress{
        .el-progress-bar__outer{
            background-color: #F35F28;
        }
        .el-progress__text{
            color: #F35F28;
        }
    }
    /deep/ .successProgress{
        .el-progress__text{
            color: #4C7AE3;
        }
    }
 </style>
