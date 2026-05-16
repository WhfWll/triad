<template>
   <div>
        <!-- bas任务 -->
        <div class="main-title  ">  
            安全配置核查 
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
                <el-table-column prop="name" label="任务名称" :show-overflow-tooltip="true">
                </el-table-column> 
			    <el-table-column   
					prop="templateName"
					label="评估方案" :show-overflow-tooltip="true">    
				</el-table-column>    
                <el-table-column prop="riskLevelEnum" label="风险" width="120">
                    <template slot="header"> 
                        <span class="cursorPointer" @click="clickButton('风险')"
                            :class="(formData.riskLevel !== '' && formData.riskLevel !== 0) ? 'active' : ''">风险<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style="width:150px;" v-model="formData.riskLevel" clearable
                            size="small" ref="riskLevel" @change="handlesearch">
                            <el-option v-for="(item, index) in riskLevellist" :key="index" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                         <!-- {{ scope.row.riskLevelEnum }} -->
                         <span :class="[
                            { 'riskstyle risk_hight': scope.row.riskLevel == '1' },
                            { 'riskstyle risk_middle': scope.row.riskLevel == '2' },
                            { 'riskstyle risk_low': scope.row.riskLevel == '3' },
                            { 'riskstyle risk_nofind': scope.row.riskLevel == '4' }
                            ]"><i></i>{{ scope.row.riskLevelEnum }} </span>
                    </template>
                </el-table-column> 
                <el-table-column   
					prop="userName"
					label="创建者"  >    
				</el-table-column> 
				<el-table-column prop="createTime" label="创建时间"  > 
                </el-table-column> 
				<el-table-column   label="状态"  > 
					<template slot-scope="scope" >  
                        <div v-if="showEditFileNameButton && rowId == scope.row.id"> 
                            <el-link  :underline="false" class="link_primary" v-if="scope.row.status != 3"  @click="handleStop(scope.row)">结束</el-link>  
                            <el-link  :underline="false" class="link_primary"   @click="handleInfo(scope.row)">详情</el-link>  
                           
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
                             <span :class="[  
                                { 'tag_status tag_warning': scope.row.status ==1 },
                                { 'tag_status tag_primary': scope.row.status == 2 },
                                { 'tag_status tag_success': scope.row.status ==3 }]"><i></i>{{ scope.row.statusEnum }}</span>
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
    name:'bastask',
    data(){
        return{
            alldelvisible:false,
            multipleSelection:[],
            tableData:[],
            formData:{
                search:'',
                page:1,
                riskLevel:'',
            },
            pageSize:10,
            showEditFileNameButton:false,
            rowId:'',
            currentpage:1,
            totalpage:0,
            riskLevellist:[],
        }
    },
    created(){
        this.$store.state.activefirstMenu="/bastask"; 
    },
    mounted(){
        this.basruleenum();
        this.getData();
    },
    methods:{
        async basruleenum(){
            const res = await bas.basEnum();
           
            if(res.code == 200){ 
                this.riskLevellist = res.data.riskLevel; 
                this.riskLevellist.unshift(
                    { label: "全部", value: 0 }
                ) 
            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }  
             
        },
        async getData(){
            const res = await bas.basTasklist({
                page:this.formData.page,
                size:this.pageSize,
                search:this.formData.search,
                riskLevel:this.formData.riskLevel,
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
                path: `/createbas`,
                query: {
                    isAdd: 1,
                }
            })
        },
        async handleDel(scope){
            let params = {
                id:scope.row.id+''
            }
            const res = await bas.basTaskdel(params)
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
            const res = await bas.basTaskdel(params)
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
            this.formData.riskLevel='';
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
            path: `/bastaskDetail`,
                query: { 
                    id: row.id,
                    name: row.name,
                }
            })
        },
        async handleStop(row){ //结束
            const res = await bas.basTaskstop({
                id:row.id
            })
            if(res.code == 200){ 
                this.$message({
                    message:'结束成功',
                    type: 'success'
                }); 
                this.getData();
            }else{
                this.$message({
                    message:res.msg,
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
        clickButton(type) {
            switch (type) {
                case '风险':
                this.$refs.riskLevel.toggleMenu();
                break;
            }
        },
    }
}
</script>