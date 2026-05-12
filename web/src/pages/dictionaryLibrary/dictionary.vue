<template>
    <div>
        <div class="main-title">
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >字典库</label>
        </div>
        <div class="tollboxlist context_box_bg">
            
            <!-- selectPlaceholder="字典类型"
            :selectList="dictionariesType" 原来的下拉框-->
            <Operation 
                one="新建" 
                two="删除"
                three="搜索"
                four="重置"
                :delList="delList"
                @handleOneClick="handleOperationUpdate"
                @handleTwoClick="handleOperationDel"
                @handleThreeClick="handleOperationSearch"
                @handleFourClick="handleOperationReset"
                @handleSelectChange="handleSelectChange"
            ></Operation>
            <el-table :data="dictionaryData" style="width: 100%"  height="calc(100% - 102px)"
                @selection-change="handleChange" @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" v-loading="loading"  class="myTable" >
                <el-table-column type="selection" width="55" :selectable='checkboxT'  > </el-table-column>
                <el-table-column prop="name" label="字典名称">
                    <!-- <template slot-scope="scope"  > 
                        <span  @click="handleEdit(scope.row)">{{scope.row.name}}</span>
                        <el-link   :underline="false" >{{scope.row.name}}</el-link>  
                    </template> -->
                </el-table-column> 
                <el-table-column prop="serviceName" label="适用范围"></el-table-column>
                <el-table-column prop="typesName" label="字典类型">
                        <template slot-scope="scope" slot="header">
                            <span class="cursorPointer"  @click="clickButton('字典类型')" :class="service ? 'spanActive' : ''">字典类型<i class="iconfont iconshaixuan"></i></span>
                            <el-select popper-class="thSelect" v-model="service"  clearable placeholder="字典类型"  size="small" ref="selectListRef" @change="handleSelectChange">  
                                <el-option
                                    v-for="(item) in dictionariesType"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value"> 
                                </el-option>
                            </el-select>  
                        </template>
                    </el-table-column>
                <el-table-column label="默认字典">
                    <template slot-scope="scope">
                        <span v-if="![4,5].includes(scope.row.types)">{{scope.row.isDefaultName}}</span>
                        <span v-else>--</span>
                    </template>
                </el-table-column>
                <el-table-column prop="updateTime" label="添加时间">
                        <template slot-scope="scope">
                            <div v-if="showEditFileNameButton && rowId == scope.row.id">
                                <el-link  @click="handleEdit(scope.row)"  class="link_primary linkafter"  :underline="false" style="margin-right: 10px">详情</el-link>
                                    <el-popover 
                                    placement="bottom" 
                                    width="210" 
                                    :visible-arrow="false"
                                    :ref="`popover_id-${scope.row.id}`" 
                                    :disabled="scope.row.types  == 4 || scope.row.types == 5"
                                    popper-class="delButton_popper">
                                    <p class="delText"><i class="el-icon-warning"></i>确定设置为默认字典吗？</p>
                                    <div style="text-align: right; margin: 0">
                                        <el-button size="mini" class="delCancel"
                                        @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消
                                        </el-button>
                                        <el-button size="mini" type="primary" @click="btnDefault(scope)">确定
                                        </el-button>
                                    </div> 
                                    <!-- <el-link  slot="reference"  :underline="false"  class="link_info linkafter" >详情</el-link>  -->
                                    <el-link  slot="reference"  :underline="false"  class="link_primary linkafter" :disabled="scope.row.types  == 4 || scope.row.types == 5"  >默认</el-link> 
                                </el-popover>
                                <el-popover
                                    placement="bottom"
                                    width="170"   
                                    :visible-arrow="false"
                                    :ref="`popover-${scope.row.id}`"
                                    :disabled="scope.row.sources  == 1 || scope.row.isDefault == 1" 
                                    popper-class="delButton_popper">
                                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                    <div style="text-align: right; margin: 0">
                                        <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover-${scope.row.id}`].doClose()">取消</el-button>
                                        <el-button size="mini" type="primary"  @click="handleDel(scope)">确定</el-button>
                                    </div>  
                                    <el-link :underline="false" class="link_danger linkafter2" slot="reference"  :disabled="scope.row.sources  == 1 || scope.row.isDefault == 1" size="small">删除</el-link>
                                </el-popover>
                            </div> 
                            <div v-else >
                                <span>{{scope.row.updateTime}}</span>
                            </div>                   
                        </template> 
                </el-table-column>
            </el-table>
            <el-pagination
                :page-size="page_size" 
                background
                layout="total,  prev, pager, next, sizes,jumper"
                :total="total"
                :current-page="page"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                >
            </el-pagination> 
           
        </div>
        <el-dialog class="zidiandia" title="新建字典" :visible.sync="createDictionaries" v-if="createDictionaries"  :before-close="cancelform" width='1184px' :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="handleCreateSave">保存</el-button>
                <el-button size="small" @click="handleSaveCancel">取消</el-button>
            </div>
            <Edit
                :exhibitionList="exhibitionList"
                :serviceType="serviceType"
                :dictionariesType="dictionariesType"
                :pageType = "pageType"
                ref="createEditRef"
            ></Edit>
        </el-dialog>

         <el-dialog title="详情"  :visible.sync="editDictionaries"  :before-close="cancelform" width='1184px' :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="handleClickEdit" v-if="isUpdate">编辑</el-button>
                <el-button size="small" v-else >编辑中</el-button>
                <el-button size="small" @click="handleEditSave" v-if="!isUpdate">保存</el-button>
                <el-button size="small" @click="handleEditCancel">取消</el-button>
            </div>
            <Edit
                v-if="editShow"
                :show="isUpdate"
                :exhibitionList="exhibitionList"
                :serviceType="serviceType"
                :dictionariesType="dictionariesType"
                :editData="editData"
                :pageType = "pageType"
                ref="cancelEditRef"
            ></Edit>
        </el-dialog>
    </div>
</template>

<script>
import Operation from "../../components/Operation";
import Edit from "../../components/edit";
// import { dictionariesType} from "../../utils/dictionaries";
import { dictionary } from '@/api/tool.js'
export default {
    name: 'Dictionary',
    data:()=>({
        selectList:{
            type: Array,
            default: function () {
                return []
            }
        },
        showEditFileNameButton:false,
        rowId:'',
       dictionaryData:[] ,
       createDictionaries:false,
       editDictionaries:false,
       dictionariesType: [],
       serviceType: [],
       exhibitionList:{
            name:"字典名称",
            namePlaceholder:"字典名称",
            typeOne:"字典类型",
            typeTwo:"适用范围",
            file:"字典内容",
            uploadName:"上传字典文件",
            uploadType:"txt",
            remarks:""
       },
       pageType:'dictionaries',
       search:"",
       service:"",
       type:"",
       source:"",
       is_default:"",
       start_date:"",
       end_date:"",
       total:0,
       page:1,
       page_size:10,
       loading:false,
       delList:[],
       editData:{},
       editShow:true,
       alldelvisible:true,
       isUpdate:true,
       value:'',
    }),
    components:{
        Operation,
        Edit,
    },
    created(){
        this.$store.state.activefirstMenu="/dictionary";  
        this.getTypes()
    },
    mounted(){
        this.getMapList();
    },
    watch:{
        editData(){
            this.editShow = false;
            this.$nextTick(()=>{
                this.editShow = true;
            })
        }
    },
    methods:{
        getTypes () {
            dictionary.getServiceSelect().then(res => {
                if (res.code === 200) {
                    this.dictionariesType = res.data.types
                    this.dictionariesType.unshift({
                        label: '全部',
                        value: ''
                    })
                } else {
                    this.$message.error(res.msg)
                }
            })
        },
        handlesearch() {},
        checkboxT(row,index){ 
            if(row.sources === 1 || row.isDefault === 1){
                return 0;
            }else{
                return 1;
            }
    	},
        handleClickEdit(){ //点击编辑
            this.isUpdate = false; 
        }, 
        async handleDel(scope){ //单个目标删除
            let params = {
                dictIds:scope.row.id
            }
            const res = await dictionary.handleDel(params) 
            if(res.code === 200){
                this.$message.success('删除成功')    
                scope._self.$refs[`popover-${scope.row.id}`].doClose();
                    this.getMapList()
            }else{
                this.$message.error(res.msg)
            } 
        },
        async handleDelMuch(ids){ //多个目标删除
             let params = {
                dictIds:ids
            }
            const res = await dictionary.handleDel(params) 
             if(res.code === 200){
                this.$message.success('删除成功')
                this.getMapList()
            }else{
                this.$message.error(res.msg)
            } 
        },
        async btnDefault(scope){ //设置默认字典
            const res = await dictionary.defaultData({
                dictId: scope.row.id,
                types: scope.row.types,
                service: scope.row.service
            });
            if(res.code == 200){
                 this.$message({
                    message: '设置默认字典成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                this.getMapList(); 
            }else{
                 this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        //  获取需要编辑的字典数据
        handleEdit(row){
            let params = {
                dictId: row.id
            }
            dictionary.getDetailData(params).then(res => {
                if(!res.code === 200){
                   return this.$message.error(res.msg)
                }else{
                    this.editData = res.data;
                }
            })
            .catch(err=>{})
            .finally(_=>{
                    // let params = {}
                    // dictionary.getServiceSelect(params).then(res => {
                    //     if(res.success){
                            this.editDictionaries = true;
                    //         if(row.types == '路径字典'){ 
                    //             let res2 = [];
                    //             res.scan_dict_list.forEach(item => {
                    //                 res2.push({
                    //                     label: item.service_type,
                    //                     value: item.service_type_number
                    //                 })
                    //             }) 
                    //             this.serviceType = this.dictionariesType
                    //         }else{
                    //             let res2 = [];
                    //             res.secret_dict_list.forEach(item => {
                    //                 res2.push({
                    //                     label: item.service_type,
                    //                     value: item.service_type_number
                    //                 })
                    //             })
                    //              this.serviceType = res2;
                    //         }
                    //     }
                    // })
                    // .catch(e =>{})
                    // .finally(_=>{
                    //     this.loading = false;
                    // })
            })
        },
        //  选择删除多条数据
        handleChange(list){
            this.delList = list;
        },
        handleOperationUpdate(value){
             this.createDictionaries = true;
            // if (this.serviceType.length > 0) {
            //     this.createDictionaries = true;
            // } else {
            //     this.$ajax({
            //         url:`/dictionary/keyvalue/service/select/`,
            //         method:"GET",
            //         params:{
            //         }
            //     })
            //     .then(res => {
            //         const result = res.data;
            //         if(result.success){
            //             result.secret_dict_list.forEach(item => {
            //                 this.serviceType.push({
            //                     label: item.service_type,
            //                     value: item.service_type_number
            //                 })
            //             })
            //             this.createDictionaries = true;
            //         }
            //     })
            //     .catch(e =>{})
            //     .finally(_=>{
            //         this.loading = false;
            //     })
            // }
            
        },
        //  删除多条数据
        handleOperationDel(){
            let delList = []
            this.delList.forEach(el =>{
                delList.push(el.id)
            })
            this.handleDelMuch(delList.join(','))
        },
        handleOperationSearch(value){
            this.search = value;
            this.getMapList();
        },
        //  重置
        handleOperationReset(){
           this.search = "";
           this.typesName = "";
           this.source ="";
           this.isDefaultName = "";
           this.service = "";
           this.end_date = "";
           this.total = 0;
           this.page = 1;
           this.page_size = 10;
           this.getMapList();
        },
        currentpage(){},
        handleSelectChange(type){
            this.service = type
            this.getMapList();
        },
        cancelform(){},
        //  新增保存
        async handleCreateSave(){
            let valiDatas = this.$refs.createEditRef.handleEdit();
            if(valiDatas === null)
                return;
            const data = this.handleFormDate(valiDatas);

            // let params = {
            //     data
            // }
            this.$ajax({
                url:"/smart/tools/dictaddoredit",
                method:"POST",
                data
            })
            .then(res =>{
                const result = res?.data;
                this.$refs.createEditRef.handleClearFiles();
                if(result?.code === 200){
                    this.handleSaveCancel();
                    this.getMapList();
                    this.$message.success('保存成功');
                }else{
                    this.$message.error('保存失败')
                }
            })
            .catch(err =>{})
        },
         //  新增取消
        handleSaveCancel(){
            this.createDictionaries = false;
            this.$refs.createEditRef.handleClearFiles();
        },
        handleFormDate(data){
            const formDate = new FormData();
            formDate.append('name',data.name);
            formDate.append('types',data.dictionariesValue);
            // formDate.append('upload',data?.files[0] || '');
            // formDate.append('service', data.dictionariesValue === 4 ? 17 : data.serviceValue);
            formDate.append('service',data.serviceValue);
            // formDate.append('description',data.remarks ? data.remarks : data.name);
            formDate.append('content',data.remarks);
            formDate.append('sources',2);
            if(data.id){
                formDate.append('id',data.id);
            }
            return formDate;
        },
        //  编辑保存
        handleEditSave(){
            let valiDatas = this.$refs.cancelEditRef.handleEdit();
            if(valiDatas === null)
                return;
            const data = this.handleFormDate(valiDatas);
            this.$ajax({
                url:"/smart/tools/dictaddoredit",
                method:"POST",
                data
            })
            .then(res =>{
                const result = res.data;
                if(result.code === 200){
                    this.handleEditCancel();
                    this.getMapList();
                    this.$message.success('保存成功')
                }else{
                    this.$message.error('保存失败')
                }
            })
            .catch(e=>{})
        },
        //  编辑取消
        handleEditCancel(){
            this.editDictionaries = false; 
            // this.editData = {};
            this.$refs.cancelEditRef.handleClearFiles();
            this.isUpdate = true;
        },
        handleSizeChange(page_size){
            this.page_size = page_size;
            this.page = 1;
            this.getMapList();
        },
        handleCurrentChange(page){
            this.page = page;
            this.getMapList();
        },
        //  获取展示列表数据
        async getMapList(){
            this.loading = true;
            
            let params = {
                search:this.search,
                types:this.service,
                // types:this.types,
                // source:this.source,
                // is_default:this.is_default,
                // start_date:this.start_date,
                // end_date:this.end_date,
                size:this.page_size,
                page:this.page
            }
            const res = await dictionary.getMapList(params)

                // const result = res.data;
                this.loading = false;
                if(res.code === 200){
                    this.dictionaryData = res.data.list;
                    this.total = res.data.total;
                }
        
        },
        icons(h, { column }) {
            let that = this
            return h(
              "div",

              [
                h("span", column.label),
                h(
                  "i",
                  {
                    slot: "reference",
                    class: "iconfont iconshaixuan",
                    style:"color:rgba(72,72,102,0.32);margin-left:5px;vertical-align:initial",
                    on: {
                        click: function() {
                            that.clickButton(column);
                        }
                    }
                  },
                  ""
                )
              ]
            );
          },
        // 图标点击事件
        clickButton(type) {
            switch (type) {
                case '字典类型':
                this.$refs.selectListRef.toggleMenu();
                break;
            }
        },
        mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){ 
            let t = this.$refs['popover_id-' + row.id].showPopper;
            let t2 = this.$refs['popover-' + row.id].showPopper; 
            if (!t && !t2) {
                this.showEditFileNameButton = false;
                this.rowId = "";
            } 
        },
    }
};
</script>

<style lang="less" scoped>
.tollboxlist {
    padding: 24px;
    background: #fff;
    height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    border-radius: 4px;
}
::v-deep .el-dialog{
    background: #F7F7FB;
    border-radius: 4px;
}
.el-table .el-link::after{
    content: "";
    padding: 0;
    width: 0;
}
::v-deep .el-popconfirm__main{
    padding: 8px 0;
}
.el-popper{
    margin: 0;
}
::v-deep .popconfirm{
    margin: 220px !important;
}
// .el-link.el-link--danger {
//     color: #F35F28;
// }
/deep/ .zidiandia .el-form-item__label{
    text-align: left!important;;
}
/deep/ .zidiandia .el-form-item__label:after{
    left:-6px;
}
/deep/ .el-input__inner{
    height: 32px !important;
    line-height: 32px  !important;
}
/deep/ .el-input__suffix{
    top: -1px;
}
/deep/ .myTable{
        thead {
            .spanActive{
                color:#4C7AE3;
                i{
                    color: #4C7AE3;
                }
            }
            .cursorPointer{
                cursor: pointer;
            }
            .cell{
                line-height: 15px;
                >span{
                    position: absolute;
                }
            }
            .iconfont{
                color:rgba(72,72,102,0.32);
                margin-left:5px;
            }
            .el-select{
                height: 0;
                visibility: hidden;
                .el-input, .el-input__inner{
                    height: 0!important;
                }
            }
        }
    }
    .linkafter  {
        display: inline-block;
        border-right: 1px solid #E8E8F5;
        padding-right: 10px !important;
        height: 14px;
        line-height: 16px;
        padding-left: 0px;
    }
    .linkafter2{
        display: inline-block;
        border-right: 0;
        padding-right: 10px !important;
        height: 14px;
        line-height: 16px;
        padding-left: 10px;
    }

</style>