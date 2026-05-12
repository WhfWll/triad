
<template>
    <el-form :model="agentform" label-width="0" ref="agentformref" :rules="Rules">
        <el-form-item label="" prop="addr">
            <label class="dialog_item_label_m">代理IP<i class="is-required" style="float: right">*</i></label>
            <el-input  v-model="agentform.addr" 
                size="small" class="frame_width" placeholder="请输入代理IP" :disabled="!isUpdate" ></el-input>
        </el-form-item>
        <el-form-item label="" prop="port">
            <label class="dialog_item_label_m">代理端口<i class="is-required" style="float: right">*</i></label>
            <el-input   v-model="agentform.port" 
                size="small" class="frame_width" placeholder="请输入代理端口" :disabled="!isUpdate" ></el-input>
        </el-form-item>
        <el-form-item prop="proto">
            <label class="dialog_item_label_m">代理协议<i class="is-required" style="float: right">*</i></label>
            <el-select v-model="agentform.proto" filterable placeholder="请选择" class="frame_width" :disabled="!isUpdate">
                <el-option v-for="(item, index) in agreementlist" :key="index" :label="item.label" :value="item.value" ></el-option>
            </el-select>
            <div style="padding-left: 130px;">
                <el-checkbox v-model="agentform.isAuth">代理服务器认证</el-checkbox>
            </div>
            
        </el-form-item>
        <div v-if="agentform.isAuth">  
            <el-form-item label="" prop="username">
                <label class="dialog_item_label_m">账号</label>
                <el-input  v-model="agentform.username" 
                    size="small" class="frame_width" placeholder="请输入账号" :disabled="!isUpdate" ></el-input>
            </el-form-item>
            <el-form-item label="" prop="password">
                <label class="dialog_item_label_m">密码</label> 
                <el-input type="password" placeholder="请输入密码"  :disabled="!isUpdate"
                    v-model="agentform.password"   class="frame_width"  >
                </el-input>
            </el-form-item> 
        </div>
    </el-form>
</template>
<style lang="less" scoped>
.frame_width {
    width: 720px;
}
i.is-required {
    margin-right: -6px;
    color: #f56c6c;
    font-size: 12px; 
}
.dialog_item_label_m {
    display: inline-block;
    min-width: 100px;
    font-size: 14px;
    font-weight: 500;
    color: rgba(72, 72, 102, 0.87);
    line-height: 16px;
    margin-right: 16px;
    margin-left: 10px;
}
/deep/ .el-form-item__error{
    left: 125px;
}
/deep/ .el-form-item .el-form-item{
        margin-bottom: 14px;
    }
</style>
<script>
export default {
    data() {
        return {
            agentform: {   
                "addr":"",
                "port":"",
                "proto":'',
                "isAuth":false,
                "username":"",
                "password":""
            },
            agreementlist: [ ],
            isUpdate:false, 
            Rules:{ 
                addr:[
                    { required: true, message: '代理IP不能为空', trigger: 'blur' },  
                ],  
                port:[
                    { required: true, message: '代理端口不能为空', trigger: 'blur' },
                ],
                proto:[
                    { required: true, message: '请选择代理协议', trigger: 'blur' },
                ],
               
            },
        }
    },
    props: {  
    },
    methods: {
        getEnum(list,flag){ 
            this.agreementlist = list;
        },
        getIsUpdate(flag) {
            this.isUpdate = flag;
        },
        getConifg(_config) {   
            this.agentform = _config; 
            if( this.agentform.proto==0){
                this.agentform.proto='';
            }
        },
        getAllData() {  
            let _this = this; 
            // this.$refs.agentformref.validate(async (valid) => {
            //     if (valid) { 
            //         return _this.agentform; 
            //     }
            // }); 
            return _this.agentform; 
        },
    },
}
</script>