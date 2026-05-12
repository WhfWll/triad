<template>
    <div>
        <div class="main-title">
            <router-link :underline="false" class="classA" :to="{ path: '/evaluationScheme' }" >评估方案
            </router-link>  
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark"  :content="title"  placement="bottom">
                    <span>{{title}}</span>
                </el-tooltip>
            </label>
        </div>  
        <BannerBox tips="" style="margin-bottom: 8px;">  
            <!-- <el-button type="primary" size="small" @click="btnSave" >保存方案</el-button>   -->
            <el-button type="primary" size="small" @click="btnUpdate" v-if="isAdd == 0 && !isUpdate">编辑方案</el-button>  
            <el-button type="primary" size="small" @click="btnSave" v-if=" isUpdate && isAdd != 0">保存方案</el-button> 
            <el-button type="primary" size="small" @click="btnSave" v-if="isAdd == 0 && isUpdate ">保存方案</el-button>
            <el-button   size="small" @click="cancalScheme">取消</el-button>  
        </BannerBox>
        <div class=" tastBtnCont">
            <div class="step_box" > 
                <div :class="class1" @click="btnStep(1)" >
                    <div> 
                        <i></i>
                        <label for="" class="s_label">通用参数</label> 
                    </div> 
                </div>
                <div :class="class2" @click="btnStep(2)"  > 
                    <div> 
                        <i></i>
                        <label for="" class="s_label">漏洞测试</label>
                    </div>
                  
                </div> 
            </div>
        </div>
        <div class="step_content_box">
            <el-form :model="sceneform" ref="sceneform" :rules="rules" style="height: 100%">
                <div v-show="class1 == 'active'" class="step"> 
                    <el-form-item label="" prop="template" >
                        <label class="dialog_item_label">方案名称 </label> 
                        <el-input  v-model="sceneform.name" style="width: 320px" :disabled="!isUpdate"  placeholder="请输入方案名称"></el-input> 
                    </el-form-item> 
                    <el-form-item>
                        <label class="dialog_item_label">方案描述</label>
                        <el-input type="textarea" :autosize="{ minRows: 3, maxRows: 5}" placeholder="请输入方案描述"
                            v-model="sceneform.describe"    :disabled="!isUpdate"
                            style="width: 720px; margin-bottom: 10px;vertical-align: text-top;" resize="none">
                        </el-input>
                    </el-form-item>
                </div>
                <div  v-show="class2 == 'active'" class="step">  
                    <div> 
                        <div  class="operationbutton" style="display:inline-block">  
                            <el-button type="primary" size="small" @click="btnSelect"  v-show="!isUpdate">查询已选</el-button> 
                        </div>
                        <div class="serach-condition" style="margin-bottom: 24px;">
                            <div class="search-text">
                                <el-input placeholder="请输入查找内容" @keydown.enter.native="handlesearch" v-model="form_search.search"
                                    class="input-with-select" size="small" clearable> </el-input>
                                <el-button type="primary"  size="small" @click="handlesearch">查询</el-button> 
                            </div>
                            <div>
                                <el-button type="primary"  size="small" @click="handleReset">重置</el-button> 
                            </div> 
                        </div>
                        <el-table ref="multipleTable" :data="tableData" tooltip-effect="dark" style="width: 100%" 
                            @select="handleSelect"
                            @select-all="handleSelectAll"
                            >
                            <el-table-column type="selection" width="55" :selectable='checkboxT' >
                            </el-table-column>
                            <el-table-column prop="nameZh" label="剧本名称"  > 
                                <template slot-scope="scope" > 
                                    <span v-if="scope.row.nameZh!=''">{{scope.row.nameZh}}</span>
                                    <span v-else>{{scope.row.name}}</span>
                                </template>
                            </el-table-column> 
                            <el-table-column prop="riskLevelEnum" label="影响级别"  > 
                            </el-table-column> 
                            <el-table-column prop="content" label="内容" show-overflow-tooltip>
                            </el-table-column>
                        </el-table>
                        <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange"
                            :current-page="currentpage" :page-size="pageSize" layout=" total,  prev, pager, next, sizes,jumper" :total="total">
                        </el-pagination>
                    </div>
                </div> 
                 
            </el-form>
        </div> 
    </div>
</template>
<style lang="less" scoped>
@import './css/createscene.less';
</style>
<style lang="less" scoped>

</style>
<script>
import bas from '@/api/bas.js'
import BannerBox from "@/components/BannerBox.vue";
export default {
    name:'createscheme',
    components: { 
        BannerBox
    },
    data(){
        return{
            title:'新增方案',
            isAdd: this.$route.query.isAdd,//1,创建，0,编辑,2:复制
            id: this.$route.query.id,
            name:this.$route.query.name,
            class1:'active',
            class2:'',
            isUpdate:false,
            currentpage:1,
            pageSize:10,
            total:0,
            sceneform:{
                name:'',
                describe:'',
            },
            rules:{
                name:[
                    { required: true, message: '方案名称不能为空', trigger: 'blur' },
                ]
            },
            tableData:[],
            STEP:2,
            form_search:{
                page:1,
                serach:'',
            },
            pageSize:10,
            tableIds:[],
            vulnlist:[],
            firstIN: 1,
            ids:'',
            isShowSelect:false,
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/evaluationScheme';
    },
    mounted(){
        if (this.isAdd == 1){ //新增
            this.isUpdate = true;
            this.title = '新增方案'; 
        }else if(this.isAdd == 0){
            this.title = '编辑方案：' + this.name;
            this.isUpdate = false;  
            this.getInfo();

        } 
        else if(this.isAdd == 2){
            this.title = '复制方案：' + this.name;
            this.isUpdate = true;  
            this.getInfo();

        }
    },
    methods:{
        async getInfo(){
            const res = await bas.bastemplatebyid({
                id: this.id
            })
            if(res.code == 200){
                if(this.isAdd == 2){
                    this.sceneform.name = 'copy_'+res.data.name;
                }else{
                    this.sceneform.name = res.data.name;
                }
               
                this.sceneform.describe = res.data.desc; 
                this.tableIds = res.data.ruleIds; 
                this.ids = this.tableIds.join(',')
            }else{

            }
        },
        checkboxT(row, index) {
            if (!this.isUpdate) {
                return 0;
            } else {
                return 1;
            }
        },
        btnUpdate(){
            this.isUpdate = true;  
            this.ids = '';
            this.form_search.page=1;
            this.form_search.search = '';
            this.currentpage = 1;
            this.getrulelist();
        },
        btnStep(index){
            for (var i = 0; i < this.STEP; i++) {
                let _i = i + 1;
                if (_i == index) {
                    this['class' + index] = 'active';
                } else {
                    this['class' + _i] = 'finish';
                }
            }    
            if (this.class1 == 'active') { 
                
            }
            if (this.class2 == 'active') {
                this.getrulelist();
                
            }
           
        },
        async getrulelist(){
            let parm = {
                page: this.form_search.page,
                size: this.pageSize,
                search:this.form_search.search, 
            } 
            if(this.isShowSelect){
                parm.ids = this.ids
            }
            const res = await bas.getBasrule(parm);
            if(res.code == 200){
                this.tableData = res.data.list;
                this.total = res.data.total;  

                let defaultSelected = []; 

                this.tableData && this.tableData.forEach((item,i)=>{
                    if (this.tableIds.indexOf(item.id) !=-1){ //存在
                        defaultSelected.push(item)
                    }
                })   
                this.$nextTick(() => {
                    this.toggleSelection(defaultSelected)
                })
            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }
        },
        toggleSelection(rows) {   // 默认选中项
            if (rows) {
                rows.forEach(row => {
                    this.$refs.multipleTable.toggleRowSelection(row, true)
                })
            } else {
                this.$refs.multipleTable.clearSelection()
            }
        },
        handleSizeChange(t){
            this.form_search.page = 1;
            this.pageSize = t;
            this.getrulelist();
        },
        handleCurrentChange(t){
            this.form_search.page = t; 
            this.getrulelist();
            this.currentpage = t;
        },
        btnSave(){
            this.$refs.sceneform.validate(async (valid) => { 
                if (valid) {   
                    let _json={}; 
                    _json.name = this.sceneform.name;
                    _json.desc = this.sceneform.describe;  
                    _json.ruleIds =  this.tableIds ; 
                    if(this.isAdd == 0){
                        _json.id = Number(this.id);
                    } 
                    const res = await bas.saveScene(_json);
                    
                    if(res.code == 200){
                        this.$message({
                            message: "保存成功",
                            type: "success",
                        });
                        this.$router.push({
                            path: `/evaluationScheme`, 
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
        cancalScheme(){
            this.$router.push({
                path: `/evaluationScheme`, 
            });
        },
        handleSelectAll(val){
            if (this.firstIN === 1) { // 意思第一次点击不会执行hangleSelectAll里面的方法
                var v = this
                // remove
                if (val.length > 0) { 
                    for (const n in val) {
                        if (this.tableIds.indexOf(val[n].id) == -1) { //存在
                            this.tableIds.push(val[n].id)
                        }
                        
                    }
                }
                if (val.length === 0) { //取消 
                    for (var i = 0; i < v.tableData.length; i++) {
                        for (var j in this.tableIds) {
                            if (v.tableData[i].id === this.tableIds[j]) { 
                                v.vulnlist.splice(j, 1)
                                this.tableIds.splice(j, 1);
                                break
                            }
                        }
                    }

                }
                if (v.vulnlist.length === 0) {
                    for (const i in val) { 
                        v.tableIds.push(val[i].id);
                    }
                } else {
                    for (const i in val) {
                        let flag = false 
                        for (const j in v.tableIds) {
                            if (v.tableIds[j]  === val[i].id) {
                                flag = true
                                break
                            }
                        }
                        if (!flag) { 
                            v.tableIds.push(val[i].id);
                        }
                    }
                } 
            } 
        },
        handleSelect(val, row){
            
            if (this.firstIN === 1) { // 设置第一次进来才回触发事件
                /* 1 => add ; 0 => remove*/
                let flag = 0
                for (const i in val) {
                    if (row.id === val[i].id) {
                        flag = 1
                        break
                    }
                }
                if (flag === 1) {
                    // 如果判断当前为添加则将当前勾选数据push到指定数组中 
                    this.tableIds.push(row.id);
                } else {
                    // 否则从数组中删除当前行数据 
                    for (const i in this.tableIds) {
                        if (this.tableIds[i] === row.id) { 
                            this.tableIds.splice(i,1);
                        } 
                    } 
                }
            } 
        }, 
        handlesearch(){ 
            this.form_search.page = 1;
            this.currentpage = 1;
            this.getrulelist();
            
        },
        handleReset(){
            this.form_search.search='';
            this.form_search.page=1;
            this.pageSize = 10;
            this.getrulelist();
            this.currentpage = 1;
        },
        async btnSelect(){
            this.isShowSelect = true;
            this.form_search.page=1;
            this.form_search.search = '';
            this.currentpage = 1;
            this.getrulelist();
            // let parm = {
            //     page: this.form_search.page,
            //     size: this.pageSize,
            //     // search:this.form_search.search, 
            //     ids:this.ids
            // }  
            // const res = await bas.getBasrule(parm)
            // if(res.code == 200){
            //     this.tableData = res.data.list;
            //     this.total = res.data.total;  

            //     let defaultSelected = []; 

            //     this.tableData && this.tableData.forEach((item,i)=>{
            //         if (this.tableIds.indexOf(item.id) !=-1){ //存在
            //             defaultSelected.push(item)
            //         }
            //     })   
            //     this.$nextTick(() => {
            //         this.toggleSelection(defaultSelected)
            //     })
            // }else{
            //     this.$message({
            //         message: res.msg,
            //         type: "error",
            //     });
            // }
        },
    }
}
</script>