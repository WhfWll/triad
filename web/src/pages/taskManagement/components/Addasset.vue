<template>
    <div>
        <el-dialog :title="title" :visible.sync="dialogform" 
            :before-close="handleCancel" width="984px" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="handleSave">保存</el-button>
                <el-button size="small" @click="handleCancel">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="assetform" label-width="0" status-icon ref="assetform" :rules="rules">
                    <el-form-item label=" " prop="ip">
                        <label class="dialog_item_label">IP/域名</label>
                        <el-input min-length="8" max-length="60"  v-model="assetform.ip" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入IP/域名"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="">
                        <label class="dialog_item_label">操作系统</label>
                        <el-input min-length="8" max-length="60"  v-model="assetform.opSys" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入I操作系统"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="">
                        <label class="dialog_item_label">硬件</label>
                        <el-input min-length="8" max-length="60"  v-model="assetform.hardware" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入硬件"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="">
                        <label class="dialog_item_label">主机名</label>
                        <el-input min-length="8" max-length="60"  v-model="assetform.hostname" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入主机名"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="">
                        <label class="dialog_item_label">归属资产组</label> 
                        <selectTree v-model="selectedNode" :options="treeData" @input="getValue" style="width:320px"> </selectTree>  
                    </el-form-item>
                    <el-form-item label=" " prop="">
                        <label class="dialog_item_label">资产名称</label>
                        <el-input min-length="8" max-length="60"  v-model="assetform.name" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入资产名称"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="">
                        <label class="dialog_item_label">资产标签</label>
                        <el-input min-length="8" max-length="60"  v-model="assetform.tags" size="small" style="width:320px"
                            autocomplete="off" placeholder="请输入资产标签"></el-input>
                    </el-form-item>
                </el-form>
            </div> 
        </el-dialog>
    </div>
</template>
<style scoped lang="less">

</style>
<script>
import { traffic } from '@/api/assetManagement.js'
import selectTree from '@/components/selectTree1.vue'
export default {
    name:'',
    props:{
        value: {}, 
        id:{},
    },
    components: {
        selectTree
    },
    data(){
        return{
            title:'添加资产',
            propsmultiple: { multiple: false },
            dialogform:false,
            assetform:{
                ip:'',
                opSys:"",
                name:'',
                groupId:'',
                hardware:'',
                hostname:'',
                tags:'',
            }, 
            multiple: true, 
            selectedNode:null,
            data:[],
            treeData: [],
            rules:{
                ip: [
                    { required: true, message: 'ip/域名不能为空', trigger: 'blur' },  
                ],
            },
            assetgrouplist:[], 
            newdata:[],
            assetids:[], 
        }
    },
    watch: {
        value(newVal, oldVal) {
            this.dialogform = newVal; 
            this.getZCList();  
            if(!this.id){
                this.title = '新增资产'
            }else{
                this.title = '编辑资产'
                this.getAssetinfo();
            }
        }
    },
    created(){
     
    },
    mounted(){

    },
    methods:{ 
        async getAssetinfo(){
            const res = await traffic.assetDetail({
                id: this.id
                })
                if (res.code == 200) {
                    let info = res.data;
                    this.assetform.ip=info.ip;
                    this.assetform.opSys=info.opSys;
                    this.assetform.name=info.name; 
                    this.assetform.hostname=info.hostname;
                    this.assetform.hardware=info.hardware; 
                    this.assetform.tags = info.tags; 
                    this.assetform.groupId = info.groupId;
                    
                    this.selectedNode =  '1_'+info.groupId;
                }
        },
        // 树状图接口
        async getZCList () {
            const res = await traffic.assetGroupTree()
            if (res.code == 200) { 
                this.newdata=[]; 
                this.treeData = this.transformData(res.data.list) 
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                })
            }
        },
        // 转换函数
        transformData(data) {
            return data.map(item => ({
                id: '1_'+item.id,
                gid:item.id,
                label: item.name,
                type: 1,
                children: item.items && item.items.length > 0 ? this.transformData(item.items) : []
            }));
        },
        getValue(value){ 
            this.assetform.groupId = value.split('_')[1];
        },
        handleSave(){
            let params = {
                ip:this.assetform.ip,
                opSys:this.assetform.opSys,
                name:this.assetform.name,
                groupId:this.assetform.groupId,
                hardware:this.assetform.hardware,
                hostname:this.assetform.hostname,
                tags:this.assetform.tags,
                id:this.id,
            }
            console.log(params);
         
            let type = !this.id ? 1:0;
            
            this.$refs.assetform.validate( async (valid) => {
                if (valid) { 
                    // type:0 编辑，1 添加
                    const res = await traffic.addasset(params,type);
                    if(res.code ==200){
                        this.$message({
                            message: '保持资产成功',
                            type: 'success'
                        })
                        this.dialogform = false; 
                        this.assetform.ip='';
                        this.assetform.opSys='';
                        this.assetform.name='';
                        this.assetform.groupId='';
                        this.assetform.hardware='';
                        this.assetform.tags='';
                        this.selectedNode=null;
                        this.$emit("AddassetCancel");
                        this.$emit("replaceData");
                    }else{
                        this.$message({
                            message: '保持资产失败',
                            type: 'error'
                        })
                    }
                }
            });
        },
        handleCancel(){
            this.dialogform = false; 
            this.$emit("AddassetCancel");
        },
    }
}
</script>