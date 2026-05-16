<template>
    <div class="create-bas-task">
        <div class="main-title">
            <router-link :underline="false" class="classA" :to="{ path: '/bastask' }" >安全配置核查
            </router-link>  
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark"  :content="title"  placement="bottom">
                    <span>{{title}}</span>
                </el-tooltip>
            </label>
        </div>  
        <div class="list_box">
            <div class="search-box create-bas-task__toolbar">
                <div class="operationbutton">
                    <el-button type="primary" size="small" @click="btnSave">保存</el-button>
                    <el-button size="small" class="btn-cancel" @click="clearTask">取消</el-button>
                </div>
            </div>
            <el-form :model="taskform" ref="form" :rules="rules" class="create-bas-task__form">
                <el-form-item label="" prop="template">
                    <label class="dialog_item_label">评估方案</label>
                    <el-select 
                        :disabled="Boolean($route.query.disabled)"
                        v-model="taskform.template" 
                        style="width: 620px;" 
                        size="small" placeholder="请选择" class="form_item_width"  
                        @click="gettemplateconfigbyid()" >
                        <el-option v-for="(item, index) in templatelist" :key="index" :label="item.name"
                            :value="item.id"></el-option>
                    </el-select> 
                    <div class="template-desc">{{ template_desc }}</div>
                </el-form-item> 
                <div class="create-bas-task__targets">
                    <label class="dialog_item_label">任务目标</label>  
                    <div class="div_width">
                        <el-table
                            :data="targetnodeList"
                            size="small"
                            style="width: 100%"
                            class="myTable"
                            @selection-change="handleSelectionChange"
                            :max-height="500">
                            <el-table-column width="55" type="selection"></el-table-column>  
                            <el-table-column prop="ip" label="Agent IP"></el-table-column>  
                        </el-table>
                    </div>
                </div>
                <el-form-item label="" prop="taskname" class="taskNameClass">
                    <label class="dialog_item_label">任务名称</label>
                    <el-input :disabled="Boolean($route.query.disabled)" 
                        v-model="taskform.taskname"   
                        size="small"
                        class="form_item_width"  
                        style="width: 620px;" 
                        placeholder="请输入任务名称" maxlength="50"></el-input>
                </el-form-item> 
            </el-form>
        </div>
    </div>
</template>
<style lang="less" scoped>
@import './css/bas-list-page.less';

.create-bas-task__toolbar {
    margin-bottom: 8px;
}

.classA {
    color: #00d4aa;
    margin-right: 12px;
    text-decoration: none;
    &:hover {
        color: #33e4c4;
    }
}

.currentpagetitle {
    color: rgba(226, 232, 240, 0.87);
    font-weight: 600;
}

.dialog_item_label {
    display: block;
    color: #94a3b8;
    font-size: 14px;
    margin-bottom: 8px;
    font-weight: 500;
}

.template-desc {
    color: rgba(148, 163, 184, 0.85);
    font-size: 13px;
    margin-top: 10px;
    line-height: 1.5;
    max-width: 620px;
}

.create-bas-task__targets {
    margin-bottom: 20px;
}

.div_width {
    margin-top: 4px;
    margin-bottom: 8px;
    max-width: 730px;
}

/deep/ .create-bas-task__form .el-form-item {
    margin-bottom: 20px;
}

/deep/ .operationbutton .btn-cancel {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(0, 212, 170, 0.25);
    color: #94a3b8;
    &:hover,
    &:focus {
        background: rgba(0, 212, 170, 0.1);
        border-color: rgba(0, 212, 170, 0.45);
        color: #e2e8f0;
    }
}

/deep/ .form_item_width .el-input__inner {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(0, 212, 170, 0.2);
    color: #e2e8f0;
    border-radius: 4px;
}

/deep/ .form_item_width .el-input__inner::placeholder {
    color: rgba(148, 163, 184, 0.45);
}

/deep/ .form_item_width .el-input.is-disabled .el-input__inner {
    background: rgba(255, 255, 255, 0.03);
    color: rgba(148, 163, 184, 0.55);
}
</style>
<script>
import bas from '@/api/bas.js'
export default {
    name:'createbastask',
    data(){
        return{
            title:'新增任务',
            isAdd: this.$route.query.isAdd,//1,创建，0,编辑,2:复制
            id: this.$route.query.id,
            name:this.$route.query.name, 
            isUpdate:false,  
            taskform:{
                template:'',
                taskname:'', 
            },
            templatelist:[],
            target_conf:[],
            targetnodeList:[],
            template_desc:'',
            rules:{},
            multipleSelection:[],
        }
    },
    created(){
        this.$store.state.activefirstMenu="/bastask"; 
    },
    mounted(){
        this.getTemplist();
        this.getSelectlist();
        // if (this.isAdd == 1){ //新增
        //     this.isUpdate = true;
        //     this.title = '新增方案'; 
        // }else{
        //     this.title = '编辑方案：' + this.name;
        //     this.isUpdate = false;  
           
        // } 
    },
    methods:{
        // async getInfo(){
        //     const res = await bas.bastemplatebyid({
        //         id: this.id
        //     })
        //     if(res.code == 200){
        //         this.sceneform.name = res.data.name;
        //         this.sceneform.describe = res.data.desc; 
        //         this.tableIds = res.data.ruleIds; 
        //     }else{

        //     }
        // },

        async getTemplist(){
            const res = await bas.getbastemplate({
                page:1,
                size:100000,
                search:''
            });
            if(res.code == 200){ 
                const list = res.data && Array.isArray(res.data.list) ? res.data.list : [];
                this.templatelist = list;
                // 设置默认值
                list.forEach(item => {
                    if (item.isDefault == 1) {
                        this.taskform.template = item.id;
                    }
                });
            }else{ 
                this.$message({
                    message: res.msg,
                    type: "error",
                }); 
            }
        },
        async getSelectlist(){
            const res = await bas.getSelectNodelist();
            if(res.code == 200){
                this.targetnodeList = res.data && Array.isArray(res.data.list) ? res.data.list : [];
            }else{

            }
        },
        async gettemplateconfigbyid(){
            const res = await bas.bastemplatebyid({
                id: this.taskform.template
            })
            if (res.code == 200) { 
                this.template_desc =res.data.desc;  
            }else{
                
            }
        },
        btnUpdate(){
            this.isUpdate = true;  
        },
           
        btnSave(){ 
            if(this.multipleSelection.length == 0) return;
    		var ids = [];
    		for (var i = 0; i < this.multipleSelection.length; i++) {
    			ids.push(this.multipleSelection[i].id);
    		} 
            this.$refs.form.validate(async (valid) => { 
                if (valid) {   
                    let _json={}; 
                    _json.name = this.taskform.taskname;
                    _json.basTemplateId = this.taskform.template;
                    _json.basNodeIds = ids;
        
                    const res = await bas.basTaskcreate(_json);
                    
                    if(res.code == 200){
                        this.$message({
                            message: "保存成功",
                            type: "success",
                        });
                        this.$router.push({
                            path: `/bastask`, 
                        });
                    }else{
                        this.$message({
                            message: res.msg,
                            type: "error",
                        });
                    }
                }
            });
        },
        clearTask(){
            this.$router.push({
                path: `/bastask`, 
            });
        },
        targetAdd(){ 
            this.target_conf.push({
                dataShow: true,
            })
        },
        async targetSave(scope){ 
            if(!scope.row.url) return;
            let arrurl = scope.row.url.split(':');
             
            const res = await bas.basagentisonline({
                host:arrurl[0],
                port:arrurl[1],
            }) ;
            if(res.code == 200){
                this.target_conf.forEach((item, i) => {
                    if (i == scope.$index) {
                        item.dataShow = false; 
                        item.status = '在线'
                    } 
                })
            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }
            
        },
        targetUpdate(scope){
            this.$set(this.target_conf[scope.$index], 'dataShow', true)
        },
        targetDelete(scope){
            this.target_conf.splice(scope.$index, 1);
        },
        handleSelectionChange(val){
            this.multipleSelection = val
        },
     
    }
}
</script>