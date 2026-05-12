<template>
    <el-dialog 
        title="生成报告" 
        :visible.sync="dialogVisible" 
        width="1184px" 
        :close-on-click-modal="false"
        :validate-on-rule-change="false" 
        :show-close="false"  >
    <div class="dialog_title_tips">
        <span>提示：点击生成报告后，将在报告中心的清单页面下载报告。</span>
    </div>
    <div class="dialog_b_btn">
        <el-button size="small" @click="saveCreate">生成报告</el-button>
        <el-button size="small" @click="clearCreate">取消</el-button>
    </div>
    <div style="padding:24px">
        <div>
            <el-form ref="reportform" :model="reportform" label-width="0" class="clearfix" :rules="rules">
                <el-form-item label=" " style="margin-bottom:0">
                    <label class="dialog_item_label">报告格式：<i class="is-required" style="float: right">*</i></label> 
                    <el-radio-group v-model="reportform.format"  > 
                        <el-radio v-for="(item,i) in formatlist" 
                        :key="i" :label="item.value" value="item.value"  
                        >{{item.label}}</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item v-show="isShow" label=" " style="margin-bottom:0">
                    <label class="dialog_item_label">输出方式：<i class="is-required" style="float: right">*</i></label> 
                    <el-radio-group v-model="reportform.outputType"  > 
                        <el-radio v-for="(item,i) in [{label:'合并输出',value:1},{label:'逐个输出',value:2}]" 
                        :key="i" :label="item.value" value="item.value"  
                        >{{item.label}}</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="" style="margin-bottom:0" v-if="type == 1">
                    <label class="dialog_item_label">报告标题：<i class="is-required" style="float: right">*</i></label>
                    <el-input v-model="reportform.name"
                        size="small"
                        style="width:520px"
                        placeholder="报告标题" 
                    ></el-input>
                </el-form-item>

                <el-form-item label="" style="margin-bottom:0" v-if="type == 1">
                    <label class="dialog_item_label">报告封面：</label>
                    <img v-if="reportform.bgImg" :src="reportform.bgImg" class="form-img" alt="">
                    <img v-else src="../assets/images/templatebg.png" class="form-img" alt="">
                    <el-button type="primary"  class="changeBgBtn" @click="clickupload1()">更换背景</el-button> 
            		<label style="color:rgba(72, 72, 102, 0.64);line-height: 20px;margin-left:3px">建议上传1200*1700px尺寸图片,文件小于10M,仅支持PNG、JPG与JPEG格式。</label>
            		<input type="file"  class="btnUploadID1"  ref='upload' @change="changeBg1($event)" style="display:none" id="input-file-ID" > 
                </el-form-item>

                <el-form-item label="" style="margin-bottom:0">
                    <label class="dialog_item_label">报告目录：</label>
                    <el-tree
                        class="catatree"
                        :data="treedata"
                        show-checkbox
                        default-expand-all
                        node-key="value"
                        ref="tree"
                        highlight-current 
                        style="display: inline-block;vertical-align: text-top;"
                        @check="(click, checked)=>{handleCheckChange(click, checked)}"
                        :props="defaultProps">
                    </el-tree>
                </el-form-item>
            </el-form>
        </div>
    </div>
</el-dialog>
</template>
<style lang="less" scoped>
/deep/ .dialog_title_tips{
    position: absolute;
    top: 20px;
    left: 124px;
    font-size: 14px;
    color: #fff;
    font-weight: 400
}
i.is-required  {
    margin-right: 4px;
    color: #f56c6c;
    font-size: 12px; 
}
.form-img {
    width: 100px;
    height: 140px;
}
.changeBgBtn {
    margin-left: 10px; 
    margin-bottom: 20px;
}
</style>
<script>
import {report} from '@/api/report.js' 
export default {
    name:'',
    data(){
        return{
            // dialogVisible:false,
            reportform:{
                format:'',
                name:'',
                outputType:2,
                bgImg:''
            },
            rules:{
                reportstyle: [
                    { required: true, message: '请选择至少一个报告格式', trigger: 'blur' },  
                ],
                title:[
                    { required: true, message: '报告标题不能为空', trigger: 'blur' }, 
                ],
            },
            treedata:[ ],
            defaultProps: {
                children: 'items',
                label: 'label',
                disabled: false
            },
            formatlist:[],
            taskTreedt:[{"value":"taskOverview","label":"任务概述","items":[]},{"value":"taskStat","label":"信息统计","items":[{"value":"targetRisk","label":"目标风险统计","items":[]},{"value":"vulRisk","label":"漏洞风险统计","items":[]},{"value":"vulType","label":"漏洞类型统计","items":[]},{"value":"topVulRisk","label":"Top 危险漏洞","items":[]}]},{"value":"targetDetails","label":"目标详情","items":[]},{"value":"vulDetails","label":"漏洞详情","items":[]}],
            targetTreedt:[{"value":"targetOverview","label":"报告概要","items":[]},{"value":"assetInfo","label":"资产信息","items":[]},{"value":"vulInfo","label":"漏洞信息","items":[]}],
            user_id:'',
            halfCheckedKeys:[],
            AllNodedata:[], //全部选中 
            objId:'',
            name:'',

            batchConfigJson:{}, 
        }
    },
    props:{
        title:{},
        type:{}, //1任务，2目标 ，3逻辑漏洞任务，4：逻辑漏洞目标
        dialogVisible:{},
        isShow:{},
         
    },
    mounted(){
        this.user_id = this.$commonjs.decryptCBC(localStorage.getItem('user_id'), this.$commonjs.myKey);   
        // this.getreportEnum();
       
        
    },
    methods:{
        getinit(_id,_name,batchConfigJson){
            this.getreportEnum();
            this.objId = _id;
            this.name = _name;
            this.batchConfigJson = batchConfigJson;
            if(this.type == 1 || this.type == 3 || this.type == 5){
                this.reportform.name = this.name+' 综述报告';
                this.treedata = this.taskTreedt
            }else if(this.type == 2 || this.type == 4 || this.type == 6){
                this.reportform.name = this.name+' 目标测试报告';
                this.treedata = this.targetTreedt
            }
            this.$nextTick(() => 
                this.checkTheNode(this.treedata)
            )
        },
        async getreportEnum(){
            const res = await report.reportEnum();
            if(res.code == 200){
                this.formatlist = res.data.format;

                this.taskTreedt = res.data.taskContent;
                this.targetTreedt = res.data.targetContent;

                this.reportform.format = this.formatlist[0].value;
            }
        },
        checkTheNode(val){ //设置选中的节点
            //获取传进来的
            if (val.length !== 0) {
                val.forEach((item) => {
                    //直接循环
                    // if (item.is_show == true) {
                        //判断当前树是否被选中 选中的话就push到 全局变量 arr 里面
                        this.AllNodedata.push(item.value); //把id push到里面
                    // }
                    if (item.items != null) {
                        //这里是判断下级是否还有数据 下级数据是用 children 包起来的
                        this.checkTheNode(item.items);
                    }
                });
                //循环以后直接选中 el-tree
                this.$refs.tree.setCheckedKeys(this.AllNodedata);
            } 
        },
        handleCheckChange(data, checked,) { 
            console.log(checked)
            this.halfCheckedKeys = checked.halfCheckedKeys; //半选 
        }, 
        async saveCreate(e){ 
            let params = {
                name:this.reportform.name,
                type:this.type,
                format:this.reportform.format,
                userId:Number(this.user_id),
                configJson:'',
                outputType: this.reportform.outputType,
                objIDName:JSON.stringify(this.batchConfigJson)
            };

            let treedatakey = this.$refs.tree.getCheckedKeys();
            let halfCheckedKeys = this.$refs.tree.getHalfCheckedKeys();
            let checkkey = treedatakey.concat(halfCheckedKeys);//全选 半选 合并
            let content = {}; 
            for(var i=0;i<this.AllNodedata.length;i++){
                let item = this.AllNodedata[i];
                if(checkkey.indexOf(item) !=-1 ){
                    content[item] =1;
                }else{
                    content[item] =0;
                }
            }
            
            let configJson = {
                objId: this.objId ,
                content:content 
            }
            if (this.reportform.bgImg != ''){
                configJson['cover'] = this.reportform.bgImg
            }
            params.configJson = JSON.stringify(configJson); 
            console.log(configJson);
   
            const res = await report.createReport(params)
            if(res.code == 200){
                this.$message({
                      message: params.name+'已生成，请到报告中心下载',
                      type: "success",
                  });
                  this.$emit('click') 
              }else{
                  this.$message({
                      message: res.msg,
                      type: "error",
                  });
              }
        },
        clearCreate(e){
            this.$emit('clearCreate', e)
        },
        clickupload1(){  
            document.querySelector('.btnUploadID1').click();
        },
        changeBg1(e){
            let that = this;
            var id = e.srcElement.id;
            var file = e.target.files[0];
            var reads = new FileReader();
            reads.readAsDataURL(file);
            reads.onload = function (e) {
                let fileMaxSize = 1024 * 10;//10M
                let size = file.size/1024;
                if (size > fileMaxSize) {
                    that.$message.error('文件大小不能大于10M！');
                    return false;
                }
                that.reportform.bgImg = this.result;
            }
        },
    }
}
</script>
