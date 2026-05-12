<template>
     <div style="padding:24px 0px 0px 28px" class="formbox">
        <el-form :model="form" status-icon :rules="editRules" ref="editRef" label-width="100px" class="demo-ruleForm" label-position="right">
            <el-form-item :label="exhibitionList.name" prop="name">
                <el-input :disabled="show" type="text" size="small"  v-model="form.name" autocomplete="off" :placeholder="exhibitionList.namePlaceholder || '请输入'"></el-input>
            </el-form-item>
            <!-- 类型 -->
            <el-form-item :label="exhibitionList.typeOne" prop="dictionariesValue">
                <el-select :disabled="showNoUpdate" v-model="form.dictionariesValue" placeholder="请选择字典类型" clearable @change="changeDic">
                    <el-option
                        v-for="item in dictionariesType"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                    </el-option>
                </el-select>
            </el-form-item>
            <!-- 范围 -->
            <el-form-item :label="exhibitionList.typeTwo" prop="serviceValue" v-if="exhibitionList.typeTwo"> 
                <el-select :disabled="showNoUpdate" v-model="form.serviceValue" placeholder="请选择适用范围" clearable >

                    <el-option
                    v-for="item in serviceType"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value">
                    </el-option>
                </el-select>
            </el-form-item>

            <el-form-item  :label="exhibitionList.file" prop="files">
                <el-checkbox-group :disabled="show" v-model="form.files" v-show="false"></el-checkbox-group>
                <el-upload
                    action=""
                    ref="editUpload"
                    :auto-upload="false"
                    :on-change="handleChange"
                    :on-remove="handleRemoveFile"
                    :show-file-list="true"
                    :limit="1"
                    :accept="exhibitionList.uploadType"
                    :disabled="show"
                    >
                    <el-button size="small" type="primary">{{exhibitionList.uploadName}}</el-button>
                    <el-button size="small" type="primary" v-show="showNoUpdate" @click="downlaod">下载字典</el-button>
                    <div for="" style=" display: inline-block; width:800px;height: 34px;line-height: 18px;margin-left: 24px;color: rgba(72, 72, 102, 0.32);vertical-align: middle;" v-if="pageType =='dictionaries'">
                        <div style="line-height: 18px; text-align: left;">可填写字典内容或导入字典文件，文件限txt格式（大小不超过1MB），字典内容用换行隔开（不超过10万行）</div>
                        <div style="line-height: 18px; text-align: left;">字典文件和字典内容都有时，以字典内容为准</div>
                    </div>
                    <div for="" style=" display: inline-block; width:800px;height: 34px;line-height: 18px;margin-left: 24px;color: rgba(72, 72, 102, 0.32);vertical-align: middle;" v-if="pageType =='auxi'">
                        <div style="line-height: 18px; text-align: left;">上传文件的文件名必须为英文字符或者数字</div>
                    </div>
                </el-upload>
            </el-form-item>

            <el-form-item :label="exhibitionList.remarks" prop="remarks" v-if="exhibitionList.remarks">
                <el-input type="textarea" :disabled="show" :rows="10" v-model="form.remarks"  placeholder=""></el-input>
            </el-form-item>
            <el-form-item :label="exhibitionList.remarks" prop="remarks" v-else>
                <el-input type="textarea" :disabled="show" :rows="10" v-model="form.remarks" style="margin-left: 20px;"
                    placeholder=""></el-input>
            </el-form-item>
        </el-form>
     </div>
</template>

<script>
import { dictionary } from '@/api/tool.js'
export default {
    
    data(){
        
        // var checkName = (rule, value, callback) => {
        //     if (!value) {
        //       return callback(new Error('上传文件不能为空'));
        //     } else {
        //       const reg = /^[0-9a-zA-Z_]*$/; 
        //       if (reg.test(value)) {
        //         callback();
        //       } else {
        //         return callback(new Error('请输入至少2个汉字'));
        //       }
        //     }
        // };
        return {
            dictionariesType: [],
            serviceType: [],
            userOptions: [],
            domainOptions: [],
            pathOptions: [],
            form:{
                name:"",
                dictionariesValue:"",
                serviceValue:"",
                files:[],
                remarks:"",
            },
            showNoUpdate:false,
            fileList:[],
            editRules:{
                name: [
                    { required: true, message: '请输入字典名称', trigger: 'blur'},
                ],
                dictionariesValue:[
                    { required: true, message:"请选择字典类型",trigger:'change'}
                ],
                serviceValue:[
                    { required: true, message:"请选择适用范围",trigger:'change'}
                ]
            },
            // show:false,
        }
    },
    props:{
        exhibitionList:{
            type:Object,
            default:()=>{
                return {}
            }
        },
        pageType:String,
        // serviceType:Array,
        // dictionariesType:Array,
        editData:Object,
            default:()=>{
                return {}
        },
        show:{},
    },
    created(){
        if(this.editData !== undefined){
            this.form.name = this.editData.name || '';
            this.form.remarks = this.editData.content || '';
            this.form.serviceValue = this.editData.service  || '';
            this.form.dictionariesValue = this.editData.types   || '';
            // this.form.files = this.editData.upload || '';
            this.form.id = this.editData.id || '';
            // this.show = true;
            this.showNoUpdate = true;

            this.publicPath = this.publicPath;
        } else {
            this.form.dictionariesValue = ''
            this.form.serviceValue = ''
        }
        this.getTypes() // 获取字典类型
    },
    methods:{
        getTypes () {
            dictionary.getServiceSelect().then(res => {
                if (res.code === 200) {
                   this.dictionariesType = res.data.types
                   this.userOptions = res.data.service.weakPass
                    this.pathOptions = res.data.service.webPathScan
                    this.domainOptions = res.data.service.subdomainScan
                    
                    if (this.editData !== undefined){
                        if (this.form.dictionariesValue == 1 || this.form.dictionariesValue == 2) {
                            this.serviceType = this.userOptions; 
                        } else if (this.form.dictionariesValue == 4) {
                            this.serviceType = this.pathOptions
                        } else {
                            this.serviceType = this.domainOptions
                        }
                    }else{
                        if (this.form.dictionariesValue === 1 || this.form.dictionariesValue === 2) {
                            this.serviceType = this.userOptions;
                        } else if (this.form.dictionariesValue === 4) {
                            this.serviceType = this.pathOptions
                        } else {
                            this.serviceType = this.domainOptions
                        }
                    }
                    console.log(1, this.serviceType)
                } else {
                    this.$message.error(res.msg)
                }
            })
        },
        changeDic (val) {
            if (val === 1 || val === 2) {
                this.serviceType = this.userOptions
            } else if (val === 4) {
                this.serviceType = this.pathOptions
            } else if (val === 5) {
                this.serviceType = this.domainOptions
            }
        },
        handleChange(_,files){
            var that = this;
            for(let i=0; i<files.length;i++){
                let _file= files[i];
                let fileSuffix = _file.name.substr(_file.name.lastIndexOf(".") + 1);
                // let fileName = _file.name.substring(0, _file.name.lastIndexOf('.'))
                // let uid = _file.uid;
                // var reg = new RegExp("[\\u4E00-\\u9FFF]+","g");
// 　　            if(reg.test(fileName)){
//                     this.$message.error('的名称应该全部为英文字母或数字')
//                     // this.$refs.editUpload.value = '';
//                     this.$refs.editUpload.clearFiles();
//                     return false
//                 }else{
                    if (fileSuffix.indexOf("txt") != -1) {
                        var input = _.target;
                        var reader = new FileReader();
                        reader.onload = function () {
                            let fileMaxSize = 1024 * 1;//10M
                            let size = _file.size/1024;
                            if (size > fileMaxSize) {
                                that.$message.error('文件大小不能大于1M！');
                                return false;
                            }

                            if(that.pageType != 'auxi'){
                                if (reader.result) {
                                    that.form.remarks = reader.result;
                                }
                            }
                            
                        };
                        reader.readAsText(_.raw);
                  } else {
                    this.$message.error('只能上传txt文件')
                    // this.$refs.editUpload.value = '';
                    this.$refs.editUpload.clearFiles();
                    return false
                  }
                // }
            }
           this.form.files = files;
        },
        handleRemoveFile(_,files){
            this.form.files = files;
        },
        handleEdit(){
            let validatas = null;
            this.$refs['editRef'].validate((valid) => {
                if(valid){
                    validatas = this.handleFormFiles();
                }
            })
            return validatas;
        },
        handleFormFiles(){
            const form = this.form;
            if(typeof form.files === 'string'){
                return {
                    ...form,
                    // files:""
                }
            }else{
                return {
                    ...form,
                    // files:[form.files[0].raw]
                }
            }
        },
        handleClearFiles(){
            this.$refs.editRef.resetFields();
             this.form.dictionariesValue = ''
            this.form.serviceValue = ''
            this.$refs.editUpload.clearFiles();
        },
        downlaod(){
            let dictionariesLabel ='';
            this.dictionariesType.forEach((item)=>{
                if(item.value ==this.form.dictionariesValue ){
                    dictionariesLabel= item.label;
                }
            })
            let serviceLabel='';
            this.serviceType.forEach((item)=>{
                if(item.value ==this.form.serviceValue ){
                    serviceLabel= item.label;
                }
            })

            let filename = this.form.name+'_'+dictionariesLabel+'_'+serviceLabel+'.txt';
        
            const blob = new Blob([this.form.remarks],{
                type:"text/plain"
            });
            const url =  window.URL.createObjectURL(blob);
            const a = document.createElement('a'); 
            a.download = filename
            a.href = url;
            a.click();
        },
    }
};
</script>

<style lang="less" scoped>
.el-input,
.el-select{
    width: 320px;
    height: 32px;
}
.el-textarea{
    width:905px;
}

::v-deep .el-form-item__label{
    margin-right: 20px;
}
::v-deep .el-form-item__label:after{
    display: block;
    content: "";
    width: 3px;
    height: 14px;
    background: #4C7AE3;
    position: relative;
    left: 18px;
    top: 13px;
    float: left;
}
::v-deep .el-form-item__label:before {
    display: block;
    content: "\00a0";
    position: relative;
    left: 110px;
    top: 4px;
    float: left;
    margin: 0  !important;
}
::v-deep .el-form-item__error {
    left: 22px;
}
::v-deep .el-upload-list{
    width: 320px;
    margin-left: 20px;
}
textarea{
    resize: none;
}
</style>