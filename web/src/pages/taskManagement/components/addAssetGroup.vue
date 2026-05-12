<template>
    <div>
        <el-dialog :title="title" :visible.sync="dialogaddFormVisible" :before-close="cancelform" width="1184px"
            :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false"  >
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="groupform" label-width="0" status-icon ref="assetgroup" :rules="rules1">
                    <el-form-item label=" " prop="name">
                        <label class="dialog_item_label">资产组名称</label>
                        <el-input v-model="groupform.name" size="small" style="width:320px" placeholder="资产组名称"
                            maxlength="50"></el-input>
                    </el-form-item>
                    <el-form-item label=" " prop="pid">
                        <label class="dialog_item_label">上级资产组</label> 
                        <selectTree v-model="selectedNode" :options="treeData" @input="getValue" style="width:320px"> </selectTree> 
                    </el-form-item>
                    <el-form-item label="" prop="range">
                        <label class="dialog_item_label" style="vertical-align:top;">资产组说明</label> 
                        <el-input type="textarea" :rows="10" class="txtareacontent" v-model="groupform.remarks" resize="none"
                            size="small" style="width:320px;display:inline-block;" autocomplete="off"
                            placeholder="资产组说明" ></el-input>
                    </el-form-item>
                </el-form>
            </div>
        </el-dialog>
    </div>
</template>
<style lang="less" scoped>

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
            title:'新建资产组',
            dialogaddFormVisible:false,
            groupform:{
                name:'',
                pid:'',
                remarks:'',
            },
            rules1:{
                name: [
                    { required: true, message: '资产组名称不能为空', trigger: 'blur' },  
                ],
            },
            selectedNode:null,
            treeData:[],
        }
    },
    watch: {
        value(newVal, oldVal) {
            this.dialogaddFormVisible = newVal;   
            this.getZCList();
            if(!this.id){
                this.title = '新建资产组';
            }else{
                this.title = '编辑资产组';
                this.getGroupinfo();
            }
        }
    },
    methods:{
        async getGroupinfo(){
            const res = await traffic.assetGroupDetail({
                    id: this.id
                })
                if (res.code == 200) {
                    let info = res.data;
                    this.groupform.name=info.name;
                    this.groupform.remarks=info.desc;
                    this.groupform.pid = info.pid;
                    this.selectedNode =  '1_'+info.pid;
                }
        },
        async getZCList () {
            const res = await traffic.assetGroupTree()
            if (res.code == 200) { 
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
                type: '1',
                children: item.items && item.items.length > 0 ? this.transformData(item.items) : []
            }));
        }, 
        getValue(value){ 
            this.groupform.pid = value.split('_')[1];  
        },
        submitForm(){
            let type = !this.id ? 1:0;
            console.log(this.selectedNode) 
            this.$refs.assetgroup.validate( async (valid) => {
                if (valid) { 
                    const res = await traffic.addassetgroup({
                        name:this.groupform.name,
                        pid:this.groupform.pid,
                        desc:this.groupform.remarks,
                        id:this.id
                    },type)
                    if(res.code ==200){
                        this.$message({
                            message: '保存资产组成功',
                            type: 'success'
                        })
                        this.dialogaddFormVisible = false; 
                        this.groupform.name=''; 
                        this.groupform.pid='';
                        this.groupform.remarks='';
                        this.selectedNode=null;
                        this.$emit("assetGroupCancel");
                        this.$emit("replaceData");
                    }else{
                        this.$message({
                            message: '保存资产组失败',
                            type: 'error'
                        })
                    }
                }
            });
        },
        cancelform(){
            this.dialogaddFormVisible = false; 
            this.groupform.name=''; 
            this.groupform.pid='';
            this.groupform.remarks='';
            this.selectedNode=null;
            this.$emit("assetGroupCancel");
        },
    }
}
</script>