<template>
    <div>
        <el-dialog title="资产同步" :visible.sync="dialogform" 
            :before-close="handleCancel" width="984px" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="handleSave">同步</el-button>
                <el-button size="small" @click="handleCancel">关闭</el-button>
            </div>
            <div style="height: calc(100% - 50px);padding:24px;">
                <div style="margin-bottom:24px">
                    发现资产是从渗透任务中提取资产管理中不包含的资产，将同步到未分组资产组中。
                </div> 
                <el-table
                    ref="multipleTable"
                    :data="tableData"
                    tooltip-effect="dark"
                    style="width: 100%" height="90%"
                    @selection-change="handleSelectionChange">
                    <el-table-column
                        type="selection"
                        width="55">
                    </el-table-column>
                    <el-table-column label="资产URL" prop="url" >
                    </el-table-column> 
                </el-table>
            </div> 
        </el-dialog>
    </div>
</template>
<style scoped lang="less">
    .el-table{
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    }
</style>
<script>
import { traffic } from '@/api/assetManagement.js'
export default {
    name:'',
    props:{
        value: {}, 
    },
    data(){
        return{
            dialogform:false,
            tableData:[],
            multipleSelection:[],
             
        }
    },
    watch: {
        value(newVal, oldVal) {
            this.dialogform = newVal;  
            this.findData();
        }
    },
    created(){

    },
    mounted(){

    },
    methods:{
        async findData(){
            const res =await traffic.assetfinddiff(); 
            if(res.code == 200){
                this.tableData = res.data.list; 
            }
        },
        async handleSave(){
            console.log(this.multipleSelection)
            if (this.multipleSelection.length == 0) return; 
            let _ids = this.multipleSelection.map(item => item.id); 
            const res = await traffic.syncassetfinddiff({
                ids:_ids.join(',')
            });
            if(res.code == 200){
                this.$message({
                  message: "资产同步成功",
                  type: "success"
                }); 
                this.dialogform = false; 
                this.$emit("findCancel");
            }else{

            }
        },
        handleCancel(){
            this.dialogform = false; 
            this.$emit("findCancel");
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        
    }
}
</script>