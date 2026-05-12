<template>
    <div>
        <el-dialog title="发现资产" :visible.sync="dialogform" 
            :before-close="handleCancel" width="984px" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="handleStartDetect">开始</el-button>
                <el-button size="small" @click="handleCancel">关闭</el-button>
            </div>
            <div style="height: calc(100% - 50px);padding:24px;">
                <div style="margin-bottom:24px">
                    发现资产是从渗透任务中提取资产管理中不包含的资产，将同步到未分组资产组中。
                </div>
                <!-- 添加表单 -->
                <el-form :model="formData" ref="formRef" :rules="rules" label-width="120px">
                    <el-form-item label="IP探测范围" prop="ipDetectRange">
                        <el-input type="textarea" v-model="formData.ipDetectRange" placeholder="请输入IP探测范围"></el-input>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" size="small" @click="handleImport">导入目标</el-button>
                        <span class="tip-text">只能上传 txt或csv或xls或xlsx格式文件</span>
                    </el-form-item>
                    <el-form-item label="IP排除范围" prop="ipExcludeRange">
                        <el-input type="textarea" v-model="formData.ipExcludeRange" placeholder="请输入IP排除范围"></el-input>
                    </el-form-item>
                    <el-form-item label="端口扫描范围" prop="portScanRange">
                        <el-select v-model="formData.portScanRange" placeholder="请选择端口扫描范围" style="width: 100%" @change="handlePortRangeChange">
                            <el-option label="自定义端口" value="custom"></el-option>
                            <el-option label="top10端口" value="top10"></el-option>
                            <el-option label="top100端口" value="top100"></el-option>
                            <el-option label="top500端口" value="top500"></el-option>
                            <el-option label="top1000端口" value="top1000"></el-option>
                            <el-option label="全部端口" value="all"></el-option>
                        </el-select>
                        <el-input 
                            type="textarea"
                            v-model="formData.portDetail"
                            :placeholder="getPortPlaceholder"
                            :rows="3"
                            style="margin-top: 10px;">
                        </el-input>
                    </el-form-item>
                </el-form>
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
            formData: {
                taskName: '',
                ipDetectRange: '',
                ipExcludeRange: '',
                portScanRange: '',
                customPort: '',
                executeType: 'immediate'
            },
            rules: {
                taskName: [
                    { required: true, message: '任务名称不能为空', trigger: 'blur' }
                ],
                ipDetectRange: [
                    { required: true, message: 'IP探测范围不能为空', trigger: 'blur' }
                ],
                portScanRange: [
                    { required: true, message: '请选择端口扫描范围', trigger: 'change' }
                ]
            },
            portRangeMap: {
                top10: '20,21,22,23,25,80,110,443,1433,3306',
                top100: '20-25,80-89,110,443,1433,3306,3389,5432,6379,8080-8089',
                top500: '1-1000',
                top1000: '1-2000',
                all: '1-65535'
            }
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
        async handleStartDetect(){
            // 表单验证
            try {
                await this.$refs.formRef.validate();
            } catch (error) {
                return;
            }
            const params = {
                ipRange: this.formData.ipDetectRange,
                excludeRange: this.formData.ipExcludeRange || '',
                portRange: this.formData.portScanRange === 'custom' ? 
                    this.formData.portDetail : 
                    this.portRangeMap[this.formData.portScanRange]
            };
            const res = await traffic.assetDetect(params);
            if(res.code == 200){
                this.$message({
                    message: "开始探测资产",
                    type: "success"
                }); 
                this.dialogform = false; 
                this.$emit("findCancel", true); // 传递一个参数表示需要开始获取进度
            } else {
                this.$message({
                    message: res.msg || "探测失败",
                    type: "error"
                });
            }
        },
        handleCancel(){
            this.dialogform = false; 
            this.$emit("findCancel");
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        handleImport() {
            // 实现文件导入逻辑
        },
        async handlePortRangeChange(value) {
            if (value === 'custom') {
                this.formData.portDetail = '';
                return;
            }
            this.formData.portDetail = this.portRangeMap[value] || '';
        }
    },
    computed: {
        getPortPlaceholder() {
            return this.formData.portScanRange === 'custom' 
                ? '请输入自定义扫描端口，多个端口用逗号分隔，端口范围用-连接，如：80,443,1000-2000' 
                : '端口范围';
        }
    }
}
</script>