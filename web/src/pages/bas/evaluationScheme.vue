<template>
    <div>
        <!-- 评估方案 -->
        <div class="main-title  ">  
            评估方案 
	  	</div>
        <div class="list_box">
            <div class="search-box"> 
                <div  class="operationbutton" >  
                    <el-button type="primary" size="small" @click="btnCreate"  >新建</el-button>
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" style="padding-left:8px"
                        trigger="click" :visible-arrow="false" v-model="alldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消</el-button>
                            <el-button size="mini" type="primary" @click="AllDel">确定</el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除</el-button>
                    </el-popover> 
				</div>
                <div class="serach-condition" > 
					<div class="search-text">
						<el-input placeholder="搜索关键字" @keydown.enter.native="handlesearch" v-model="formData.search" class="input-with-select"  size="small" clearable > </el-input>
						<el-button type="primary"  size="small" @click="handlesearch">搜索</el-button> 
					</div> 
					<div >
						<el-button type="primary"  size="small" @click="handleReset">重置</el-button> 
					</div>   
				</div>  
            </div> 
			<el-table
				:data="tableData"  style="width: 100%"  class="myTable"   @selection-change="handleSelectionChange" 
				  @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" >
                  <el-table-column  
                    width="55" 
                    type="selection"  >
                </el-table-column>  
			    <el-table-column   
					prop="name"
					label="方案名称" :show-overflow-tooltip="true">    
				</el-table-column>    
                <el-table-column prop="desc" label="描述" :show-overflow-tooltip="true">
                </el-table-column> 
                <el-table-column prop="isDefaultEnum" label="默认方案" > 
                </el-table-column>
				<el-table-column prop="createTime" label="创建时间"  > 
					<template slot-scope="scope" >  
                        <div v-if="showEditFileNameButton && rowId == scope.row.id"> 
                            <el-link  :underline="false" class="link_primary"   @click="handleInfo(scope.row)">详情</el-link>  
                            <el-link  :underline="false" class="link_primary"  @click="handleCopy(scope.row)">复制</el-link>  
                            <el-link  :underline="false" class="link_primary" v-if="scope.row.isDefault==0"  @click="handleDefault(scope.row)">默认</el-link>  
                            <el-popover 
                                placement="bottom"
                                width="170"   
                                :ref="`popover_id-${scope.row.id}`"
                                popper-class="delButton_popper" >
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                    <el-button size="mini" type="primary" @click="handleDel(scope)">确定</el-button>
                                </div> 
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0"  slot="reference" >删除</el-link>  
                            </el-popover>
                        </div>
                        <div v-else >
                             {{scope.row.createTime}} 
                        </div> 
                    </template>
                   
				</el-table-column>
			</el-table> 
			<el-pagination
				:page-size="pageSize" 
				background
				layout=" total,  prev, pager, next,sizes, jumper"
				:total="totalpage"
				:current-page="currentpage"
				@current-change = "handlecurrentchange"
				@size-change="handleSizeChange" >
			</el-pagination>
        </div>
    </div>
</template>
<style lang="less" scoped>
@import './css/bas-list-page.less';
</style>
<script>
import bas from '@/api/bas.js'
export default {
    name:'evaluationScheme',
    data(){
        return{
            alldelvisible:false,
            multipleSelection:[],
            tableData:[],
            formData:{
                search:'',
                page:1,
            },
            pageSize:10,
            showEditFileNameButton:false,
            rowId:'',
            currentpage:1,
            totalpage:0,
        }
    },
    created(){
        this.$store.state.activefirstMenu="/evaluationScheme"; 
    },
    mounted(){
        this.getData();
    },
    methods:{
        async getData(){
            const res = await bas.getbastemplate({
                page:this.formData.page,
                size:this.pageSize,
                search:this.formData.search
            });
            if(res.code == 200){
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }else{ 
                this.$message({
                    message: res.msg,
                    type: "error",
                }); 
            }
        },
        btnCreate(){
            this.$router.push({
                path: `/createscheme`,
                query: {
                    isAdd: 1,
                }
            })
        },
        async handleDel(scope){
            let params = {
                id:scope.row.id
            }
            const res = await bas.bastemplatedel(params)
            if(res.code == 200){ 
                this.$message({
                    message:'删除成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                this.getData();
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        async AllDel(){
            if(this.multipleSelection.length == 0) return;
    		var ids = [];
    		for (var i = 0; i < this.multipleSelection.length; i++) {
    			ids.push(this.multipleSelection[i].id);
    		}
            let params = {
                id:ids.join(',')
            }
            const res = await bas.bastemplatedel(params)
            if(res.code == 200){ 
                this.$message({
                    message:'删除成功',
                    type: 'success'
                });
                this.alldelvisible = false;
                this.getData();
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        handlesearch(){
            this.formData.page = 1;
			this.getData();
			this.currentpage = 1;
        },
        handleReset(){
            this.formData.page = 1;
            this.formData.search='';
            this.pageSize=10;
			this.getData();
			this.currentpage = 1;
        },
        handlecurrentchange(t){
            this.formData.page = t; 
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t){
            this.formData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handleInfo(row){
            this.$router.push({
            path: `/createscheme`,
                query: {
                    isAdd: 0,
                    id: row.id,
                    name: row.name,
                }
            })
        },
        handleCopy(row){ //复制
            this.$router.push({
                path: `/createscheme`,
                query: {
                    isAdd: 2,
                    id: row.id,
                    name: row.name,
                }
            })
        },
        async handleDefault(row){ //默认
            const res = await bas.bastemplatesetdefault({
                id:row.id
            });
            if(res.code == 200){
                this.$message({
                    message:'设置默认成功',
                    type: 'success'
                }); 
                this.getData();
            }else{
                this.$message({
                    message:'设置默认失败',
                    type: 'error'
                });  
            }
        },
        handleSelectionChange(val){
            this.multipleSelection = val
        },
        mouseenter(row,colum,cell,event){  
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
            if (this.$refs['popover_id-' + row.id]) {
                let t = this.$refs['popover_id-' + row.id].showPopper;
                if(!t){
                    this.showEditFileNameButton = false;
                    this.rowId = "";
                }
            } else {
                this.showEditFileNameButton = false;
                this.rowId = "";
            }
        },
    }
}
</script>