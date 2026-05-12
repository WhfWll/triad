<template>
    <div > 
        <div class="totalbox">  
            <div class="tabsbox  ">
                <div class="planTime">
                    <label style="margin-right: 10px;">启用模型增强</label>
                    <el-switch
                        class="switchbtn"
                        active-color="#4c7ae3"
                        v-model="isShow"
                        inactive-color="#E8E8F5"
                        @change="fnSwitch">
                    </el-switch> 
                    <div data-v-2754877f="" class="tipsCont">
                        <label data-v-2754877f="" class="topTips">提示：</label>
                        <span data-v-2754877f="">请确保至少配置一个文本模型和一个图像模型，否则增强功能可能无法正常工作。</span>
                    </div>
                </div>
            </div>  
        </div> 
        <div class="totalbox boxlist">
            <div class="search-box">
                 <div class="operationbutton">
                    <el-button type="primary" @click="createAIset">新建模型</el-button>
                 </div>
            </div>

            <el-table :data="tableData"   @selection-change="handleSelectionChange"
                @cell-mouse-enter="mouseenter"
                @cell-mouse-leave="mouseleave" > 
                <el-table-column prop="modelName" label="模型名称" show-overflow-tooltip>
                    <template #default="scope">
                        <span>{{ scope.row.modelName }}<span class="tagstyle" v-if="scope.row.isDefault == 1">默认模型</span> </span>
                    </template>
                </el-table-column>
                <el-table-column prop="platform" label="平台类型" show-overflow-tooltip> </el-table-column>
                <el-table-column prop="apiUrl" label="API地址" show-overflow-tooltip> </el-table-column>
                <el-table-column prop="apiKey" label="API秘钥" show-overflow-tooltip> </el-table-column>
                <el-table-column prop="modelID" label="模型ID" show-overflow-tooltip> </el-table-column>
                <!-- <el-table-column prop="modelTypeDesc" label="模型类型" width="100" show-overflow-tooltip> </el-table-column> -->
                <el-table-column prop="name" label="模型类型" width="200" > 
                     <template #default="scope"> 
                         <div v-if="showOperateButton && rowId == scope.row.id">
                            <el-link
                                    :underline="false"
                                    class="link_primary"
                                    style="vertical-align: initial"
                                    @click="modelcheck(scope.row)"
                                    >测试</el-link>
                            <el-link
                                    :underline="false"
                                    class="link_primary"
                                    style="vertical-align: initial"
                                    @click="editModel(scope.row)"
                                    >编辑</el-link>
                            <el-link v-if="scope.row.isDefault!=1"
                                    :underline="false"
                                    class="link_primary"
                                    style="vertical-align: initial"
                                    @click="setDefaultModel(scope.row)"
                                    >默认</el-link>
                            <el-popover placement="bottom" width="170" :visible-arrow="false" :ref="`popover_id-${scope.row.id}`"
                                    popper-class="delButton_popper">
                                    <p class="delText">
                                    <i class="el-icon-warning"></i>确定删除吗？
                                    </p>
                                    <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="btnCancelDel(scope)">取消</el-button>
                                    <el-button size="mini" type="primary" @click="deleteModel(scope)">确定</el-button>
                                    </div>
                                    <!-- <span slot="reference">删除</span> -->
                                    <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference">删除
                                    </el-link>
                                </el-popover>
                            </div>
                            <div v-else>
                                {{ scope.row.modelTypeDesc}}
                            </div>
                    </template>
                </el-table-column>
            </el-table>

            <el-pagination  :page-size="param.size" background layout="total, prev, pager, next, sizes, jumper"
                :total="total" :current-page="param.page" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>

        </div>
       

          <!-- AI模型配置对话框 -->
        <el-dialog :visible.sync="modelDialogVisible" :title="modelDialogTitle" width="1000px"
         :close-on-click-modal="false" class="ai-model-dialog" :show-close="false">
            <div class="dialog_b_btn" >  
                <el-button  @click="submitModelForm" size="small"  >确定</el-button>
                <el-button @click="cancelModelDialog" size="small"  >取消</el-button> 
            </div>
            <el-form :model="modelForm" :rules="modelRules" ref="modelForm" 
                label-width="100px" style="padding: 20px 30px 15px 30px;"> 
                <el-form-item label="模型名称" prop="name">
                    <el-input v-model="modelForm.name" placeholder="请输入模型名称" size="medium" style="width: 320px;"></el-input>
                </el-form-item>
                <el-form-item label="模型类型" prop="modelType">
                    <el-select v-model="modelForm.modelType" placeholder="请选择模型类型"  style="width: 320px;" size="medium">
                        <el-option v-for="(item,i) in modelTypelist" :key="i" 
                        :label="item.label" :value="item.value"></el-option>
                    </el-select>
                </el-form-item> 
                <el-form-item label="平台类型" prop="platform">
                    <el-select v-model="modelForm.platform" placeholder="请选择平台类型" style="width: 320px" size="medium">
                         <el-option v-for="(item,i) in platformlist" :key="i" 
                            :label="item.label" :value="item.value"></el-option>
                    </el-select>
                </el-form-item> 
        
                <el-form-item label="API地址" prop="apiUrl">
                    <el-input v-model="modelForm.apiUrl" placeholder="请输入API地址" size="medium"  style="width: 320px;"></el-input>
                </el-form-item>
            
        
                <el-form-item label="API秘钥" prop="apiKey">
                    <el-input v-model="modelForm.apiKey" type="password" placeholder="请输入API秘钥" show-password size="medium"  style="width: 320px;"></el-input>
                </el-form-item>
            
                <el-form-item label="模型ID" prop="modelId">
                    <el-input v-model="modelForm.modelId" placeholder="请输入模型ID" size="medium"  style="width: 320px;"></el-input>
                </el-form-item>
            
                <!-- <el-form-item label="设为默认">
                    <el-switch v-model="modelForm.isDefault" active-color="#4c7ae3"></el-switch>
                </el-form-item>  -->
                
               
            </el-form>
        </el-dialog>

    </div>
</template>
<style lang="less" scoped>
    /deep/ .el-form-item.is-required:not(.is-no-asterisk)>.el-form-item__label:before{
        position: absolute;
        left: 86px;
    }
    .tipsCont  {
        display: inline-block;
        line-height: 31px;
        margin-left: 24px;
        font-size: 14px;
        span  {
            color: rgba(72, 72, 102, 0.64);
            font-size: 13px;
        }
        .topTips  {
            color: #4c7ae3;
            width: 40px;
        }
    } 
    .tagstyle{
        display: inline-block;
        font-size: 12px;
        margin-left: 8px;
        background-color: rgba(72, 72, 102, 0.14);
        border-radius: 4px;
        padding: 4px 8px;
    }
    .totalbox{
        padding: 24px 24px;
        box-sizing: border-box;
        background: #fff; 
        box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.12);
        border-radius: 4px 4px 4px 4px; 
        margin-bottom: 16px;
        
    }
    .tabsbox {
        border-radius: 4px;
    }
    /deep/ .el-form-item__label{
        position:relative;
        height: 16px;
        border-left: 3px solid #4c7ae3;
        line-height: 16px;
        padding-left: 10px;
        top: 12px;
    }
 
    .dialog_item_label {
        font-size: 14px;
        border-left: 3px solid #4c7ae3;
        padding-left: 8px;
        font-weight: 500;
        width: 113px;
        display: inline-block;
        line-height: 16px;
        box-sizing: border-box;
        color: rgba(72, 72, 102, 0.87);
    }
    .boxlist{
        height: calc(100% - 95px);
    }
</style>
<script>
import { thresholdFreedmanDiaconis, tree } from 'd3';
import { toJSONString } from 'xe-utils';
import { system } from '@/api/system.js'
export default {
    data(){
        return{
            isShow:false,
            tableData:[],
            showOperateButton:false,
            rowId:'',
            param:{
                size:10,
                page:1,
            },
            ids:[],
            total:0,
            modelDialogTitle:'',
             // AI模型相关数据
            aiModelTableData: [],
            aiModelCurrentPage: 1,
            aiModelPageSize: 10,
            aiModelTotal: 0,
            aiModelLoading: false,
            modelDialogVisible: false,
            modelForm: {
                name: '',
                platform: 'OpenAI兼容模式',
                apiUrl: '',
                apiKey: '',
                modelId: '',
                // isDefault: false,
                modelType:0
            },
            modelFormRules: {
                name: [{ required: true, message: '请输入模型名称', trigger: 'blur' }],
                apiUrl: [{ required: true, message: '请输入API地址', trigger: 'blur' }],
                apiKey: [{ required: true, message: '请输入API密钥', trigger: 'blur' }],
                modelId: [{ required: true, message: '请输入模型ID', trigger: 'blur' }]
            },
            modelCurrentPage: 1,
            modelPageSize: 10, 
            modelRules: {
                name: [
                    { required: true, message: '请输入模型名称', trigger: 'blur' }
                ],
                platform: [
                    { required: true, message: '请选择平台类型', trigger: 'change' }
                ],
                apiUrl: [
                    { required: true, message: '请输入API地址', trigger: 'blur' },
                    { type: 'url', message: '请输入正确的URL格式', trigger: 'blur' }
                ],
                apiKey: [
                    { required: true, message: '请输入API秘钥', trigger: 'blur' }
                ],
                modelId: [
                    { required: true, message: '请输入模型ID', trigger: 'blur' }
                ],
                modelType:[
                     { required: true, message: '请选择模型类型', trigger: 'change' }
                ]
            },
            modelTypelist:[],
            platformlist:[],
        }
    },
    mounted(){
        
    },
    methods:{
        init(){
            this.getEnum();
            this.getData(); 
            this.getenhancedetail();
        },
        async fnSwitch(){
            console.log(11)
            const res = await system.enhanceedit({
                status:this.isShow?1:2
            });
            if(res.code === 200) {
                 this.$message.success('设置成功');
            } else {
                this.$message.error(res.msg );
            }
        },
      
        async getenhancedetail(){
            const res = await system.enhancedetail();
            if(res.code == 200){
                if(res.data.status == 2){ 
                    this.isShow = false;
                } 
                if(res.data.status == 1){ 
                    this.isShow = true;
                }
            }
        },
        async getEnum(){
            const res = await system.getAIenum( );
            if(res.code === 200) {
                this.modelTypelist = res.data.modelType
                this.platformlist = res.data.platform;

                
            } else {
                this.$message.error(res.msg);
            }
        },
        createAIset(){
            this.showAddModelDialog();

        },
        
        async getData() { 
            const res = await system.getAiModels({
                page: this.param.page,
                size: this.param.size
            });
            if(res.code === 200) {
             
                this.tableData = res.data.list;
                this.total = res.data.total;
            } else {
                this.$message.error(res.msg || '加载AI模型失败');
            }
           
        },

        showAddModelDialog() {
            this.modelDialogTitle = '新增大模型';
            this.isEditingModel = false;
            this.editingModelId = null;
            this.resetModelForm();
            this.modelDialogVisible = true;
        },

        editModel(row) {
            console.log(row)
            this.modelDialogTitle = '编辑大模型';
            this.isEditingModel = true;
            this.editingModelId = row.id;
            this.modelForm = {
                name: row.modelName,
                platform: row.platform,
                apiUrl: row.apiUrl,
                apiKey: row.apiKey,
                modelId: row.modelID,
                // isDefault: row.isDefault,
                modelType:row.modelType,
            };
            this.modelDialogVisible = true;
        },
     
        async modelcheck(row){//测试
            const res = await system.modelcheck({
                id:row.id
            });  
            if(res.code == 200){
                if(res.data.enabled== 1){
                    this.$message.success(res.data.enabledDesc);
                }
                else{
                    this.$message.error(res.data.enabledDesc);
                }
             
            }else{
                 this.$message.error('测试模型失败');
            }
        },
        async setDefaultModel(row) {
            try {
                // 调用API设置默认模型
                const res = await system.setDefaultAiModel({ id: row.id });
                if(res.code === 200) {
                    this.$message.success('设置默认模型成功');
                    this.getData();
                } else {
                    this.$message.error(res.message || '设置默认模型失败');
                }
            } catch (error) {
                console.error('设置默认模型失败:', error);
                this.$message.error('设置默认模型失败');
            }
        },

        async deleteModel( scope) {
            try {
                // 调用API删除模型
                const res = await system.deleteAiModels({ ids: scope.row.id+'' });
                if(res.code === 200) {
                    this.$message.success('删除成功');
                    this.getData();
                } else {
                    this.$message.error(res.message || '删除失败');
                }
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
            } catch (error) { 
                this.$message.error('删除失败');
            }
        },

        submitModelForm() {
            this.$refs.modelForm.validate(async (valid) => {
                if (valid) {
                    try {
                        // 准备接口参数，映射表单数据到接口格式
                        const apiParams = {
                            modelName: this.modelForm.name,
                            platform: this.modelForm.platform,
                            apiUrl: this.modelForm.apiUrl,
                            apiKey: this.modelForm.apiKey,
                            modelID: this.modelForm.modelId,
                            // isDefault: this.modelForm.isDefault ? 1 : 0,
                            modelType:this.modelForm.modelType
                        };

                        // 如果是编辑模式，添加id参数
                        if (this.isEditingModel) {
                            apiParams.id = this.editingModelId;
                        }

                        const res = await system.addOrUpdateAiModel(apiParams);
                        
                        if (res.code === 200) {
                            this.$message.success(this.isEditingModel ? '编辑成功' : '新增成功');
                            this.modelDialogVisible = false;
                            this.resetModelForm();
                            // 重新加载模型列表
                            this.getData();
                        } else {
                            this.$message.error(res.msg || (this.isEditingModel ? '编辑失败' : '新增失败'));
                        }
                    } catch (error) {
                        console.error('提交模型表单失败:', error);
                        this.$message.error(this.isEditingModel ? '编辑失败' : '新增失败');
                    }
                }
            });
        },

        cancelModelDialog() {
            this.modelDialogVisible = false;
            this.resetModelForm();
        },

        resetModelForm() {
            this.modelForm = {
                name: '',
                platform: '',
                apiUrl: '',
                apiKey: '',
                modelId: '',
                isDefault: false,
                modelType:'',
            };
            if (this.$refs.modelForm) {
                this.$refs.modelForm.resetFields();
            }
        },
 
        handleSelectionChange(selection){
            this.ids = selection.map(item => item.id) 
        },
        handleSizeChange(t){
            this.param.page = 1;
            this.param.size = t;
            getData();
        },
        currentchange(t){
            this.param.page = t;
            this.param.size = 10;
            this.getData();
        },
        delData(){

        },
        btnCancelDel(scope){
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose(); 
        },
        mouseenter (row, column, cell, event) {
            this.showOperateButton = true
            let _id = row.id
            this.rowId = _id//赋值行id，便于页面判断
            },
        mouseleave (row, colum, cell, event) {
            let t = this.$refs['popover_id-' + row.id] && this.$refs['popover_id-' + row.id].showPopper
            let t2 = this.$refs['popover-' + row.id] && this.$refs['popover-' + row.id].showPopper

            if (!t && !t2) {
                this.showOperateButton = false
                this.rowId = ""
            } 
        },
    }
}
</script>