<template>
  	<div > 
  		<div class="main-title  ">  
			报告模板
	  	</div> 
	  	<div class="templatelist context_box_bg"> 
	  		<div class="search-box" >
				<div class="operationbutton"  > 
                    <!-- <el-button type="primary"  size="small" @click="addTaskreport()">新建任务报告模板</el-button> -->
                    <xzbutton 
                    type="primary" 
                    @click="addTaskreport" 
                    size="small">新建任务报告模板</xzbutton>

                    <!-- <el-button type="primary"  plain size="small" class="whitebtn" style="" @click="addTargetreport()">新建目标报告模板</el-button>  -->
                    <xzbutton 
                    type="blueBorderWhiteBg"
                    class="whitebtn nomargin"
                    @click="addTargetreport" 
                    size="small" style="margin-left: 8px;">新建目标报告模板</xzbutton>
<!-- 
                    <el-popover
                        popper-class="delButton_popper"
                        placement="bottom-start"
                        width="170" 
                        trigger="click" 
                        :visible-arrow="false"
                        v-model="alldelvisible" >
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="" >
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false" >取消</el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDeleteTemplate" >确定</el-button>
                        </div>  
                        <el-button type="warning"  size="small"  slot="reference" :disabled="!multipleSelection.length">删除</el-button> 
                    </el-popover>   -->
                    <delbutton 
                    :width="170"  
                    @click="btnMultiDeleteTemplate"  
                    :disabled="!multipleSelection.length"></delbutton> 
				</div> 
				<div class="serach-condition" > 
					<div class="search-text">
						<el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearch" v-model="search_item.search" class="input-with-select"  size="small" clearable > </el-input>
						<!-- <el-button type="primary"  size="small" @click="handlesearch">搜索</el-button> -->
                        <xzbutton 
                        type="primary" 
                        @click="handlesearch" 
                        :disabled="false" 
                        size="small"  >搜索</xzbutton>
					</div>
					<div >
						<!-- <el-button type="primary"  size="small" @click="handleReset">重置</el-button> -->
                        <xzbutton 
                        type="primary" 
                        @click="handleReset" 
                        :disabled="false" 
                        size="small"  >重置</xzbutton> 
					</div>  
				</div> 
			</div>
			<el-table
                ref="multipleTable"
                :data="tableData" 
                tooltip-effect="dark"
                v-model="Loading"
                style="width: 100%"
                @selection-change="handleSelectionChange"  @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column
                    type="selection"
                    width="55" :selectable='checkboxT'>
                    </el-table-column> 
                <el-table-column
                    prop="name"
                    label="模板名称" > 
                     <template slot-scope="scope"> 
                        <el-link  @click="fnDetails(scope.row)">{{scope.row.name}}</el-link> 
                    </template>
                </el-table-column>
                <el-table-column
                    prop="type"
                    label="模板类型" > 
                    <template slot-scope="scope"> 
                        <span v-if="scope.row.type == 1">任务报告</span>
                        <span v-if="scope.row.type == 2">目标报告</span>
                    </template>
                </el-table-column>
                <el-table-column
                	show-overflow-tooltip
                    :render-header="icons"
                    prop="default"
                    label="默认模板"  >
                    <template slot-scope="scope">
                        <span v-if="scope.row.default">是</span>
                        <span v-else>否</span>
                    </template>
                </el-table-column>
                <el-table-column
                    prop="create_time"
                    label="创建时间" >
                    </el-table-column>  
                <el-table-column
                    prop="user"
                    label="提交者" >
                    <template slot-scope="scope">
                        <div v-if="showEditFileNameButton && rowId == scope.row.id">
                            <el-popover
                                placement="bottom"
                                width="200" 
                                :visible-arrow="false"
                                :ref="`popover-${scope.row.id}`"
                                popper-class="delButton_popper" >
                                <p class="delText"><i class="el-icon-warning"></i>确定设置为默认模板吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover-${scope.row.id}`].doClose()" >取消</el-button>
                                    <el-button size="mini" type="primary" @click="btnDefaultTemplate(scope)">确定</el-button>
                                </div>  
                                <el-link :underline="false" class="link_primary linkafter" style="padding:0" slot="reference" >默认</el-link> 
                            </el-popover>  
                            <el-popover
                                placement="bottom"
                                width="170"   
                                :visible-arrow="false"
                                :disabled="scope.row.default?true:false"
                                :ref="`popover_id-${scope.row.id}`"
                                popper-class="delButton_popper" >
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel"  @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                    <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                </div> 
                                
                                <el-link :underline="false" class="link_danger linkafter2" slot="reference" :disabled="scope.row.default?true:false">删除</el-link> 
                            </el-popover> 
                        </div>   
                        <div v-else >
                            <span>{{scope.row.user}}</span>
                        </div>
                    </template>
                    </el-table-column> 
                <!-- <el-table-column label="操作" >
                    
                </el-table-column> -->
            </el-table>
            <el-pagination 
                background
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                :current-page="currentPage" 
                :page-size="pageSize"
                layout=" total,  prev, pager, next, sizes,jumper"
                :total="total">
            </el-pagination>
	  	</div>
	  	<taskdialog 
            v-model="sendVal"  
            :title="dialogtitle" 
			:flag="flag"
            :templateid="templateid"
            :templatename="templatename"
            @danger="clickSave()"></taskdialog>
  	</div>
</template>
<style lang="less" scoped>
.whitebtn{
    margin-right:10px;
    color: #4c7ae3 !important;;
    background:#fff!important;
    border-color: #4C7AE3;
}
.templatelist{
	padding: 24px; 
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
}
.linkafter {  
    display: inline-block;
    border-right: 1px solid #E8E8F5;
    padding-right: 10px !important;
    height: 14px;
    line-height: 16px;
    padding-left:10px ;
}
.linkafter2 {  
    display: inline-block;
    // border-right: 1px solid #E8E8F5;
    padding-right: 10px !important;
    height: 14px;
    line-height: 16px;
    padding-left:10px ;
}
.target_active{
        .msg_icon{ 
            background: rgba(255, 255, 255, 0.3);
            color: #fff;
        }
        .msg_label{
            color: rgba(255, 255, 255, 0.7) !important;
        }
        .msg_count{
            label{
                color: #fff !important;
            }
            i{
                color: #fff !important;
            } 
        } 
        .icontsstyle{
            color: rgba(255, 255, 255, 0.7) !important;;
        }
    } 
</style>
<script>    
import taskdialog from '../../components/taskdialog.vue' 
import xzbutton from "@/components/XzButton.vue"; 
import delbutton from "@/components/DelButton.vue";
import API from '@/api/report.js'
export default({
    name:'reporttemplate',
    components: {
        taskdialog,
        xzbutton,
		delbutton, 
    },
    data(){  
    	return{ 
            showEditFileNameButton:false,
            rowId:'', 
            rowId2 :'',
            sendVal: false,
            flag:5,
            dialogtitle:'',
            templateid:'',
            templatename:'',
            tableData:[],
            pageSize:10,
            total:0,
            currentPage:1,
            multipleSelection:[],
            search_item:{
                search_field:'',
                page:1,
            },
            Loading:false,
            alldelvisible:false,
            delvisible: false,  
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/reporttemplate"; 
        this.pageSize = this.commonjs.pageSize;
    },
    mounted:function(){
        this.getData();   
    },
    methods:{  
        checkboxT(row,index){ 
            if(row.default){
                return 0;
            }else{
                return 1;
            }
    	},
        	icons(h,{column}){
            const inReview = '任务报告和目标报告都只允许有一个默认报告模板，且默认模板不允许删除' 
			return h('div', {
					style: { 
							'padding-left':' 0 !important',
							'height': '16px',
							'line-height': '16px',
							'overflow': 'initial',
						}
					},
					[ h('span', column.label),
                    h('el-tooltip', {
                        props: {
                            placement: 'top'
                        }
                    }, [
                        h('div', {
                            slot: 'content',
                            style: {
                                // 'width':'100px',
                                whiteSpace: 'normal', 
                            }
                        }, inReview), 
                        h('i', {
                            class: 'iconfont icontishi',
                            style: 'color:rgba(72,72,102,0.32);margin-left:5px;vertical-align: initial;'
                        })
                    ],)
                ],
    　　　　 )
        },
        getData(){
            this.Loading = true;  ///task/template/
            let params = {
                search:this.search_item.search,  
                page:this.search_item.page,
                page_size:this.pageSize
            }
            API.getReportTemplateList(params).then(res => {
                if(!res.error){
                    this.Loading = false; 
                    this.tableData = res.data; 
                    this.total = res.count;
                }
				else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            })
            .catch(data=>{ });
        },
        handlesearch(){ 
            this.search_item.page = 1;
            this.currentPage = 1;
            this.getData();
        },
        handleReset(){  //重置
			this.search_item.page =1;
			this.search_item.search=''; 
            this.pageSize = 10;
            this.currentPage = 1;
			this.getData();
		},
    	handleSelectionChange(val){
            this.multipleSelection = val;
        },
        handleSizeChange(t){
            this.search_item.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handleCurrentChange(t){
            this.search_item.page = t;
            this.getData();
        },
        btnDel(scope){ //删除模板
            let params = {
                id:scope.row.id,
            }
            API.deleteReportTemplate(params).then(res => {
                if(res.success){
                    this.deldialogVisible = false;
                    this.$message({
                        message:'删除模板成功',
                        type: 'success'
                    });
                    scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                    this.getData();
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }
            })
            .catch(err=>{})
             
        },
        // btnupdate(row){  
        //     this.dialogtitle = row.template_name; 
        //     this.templateid = row.id;
        //     this.templatename = row.template_name;
        //     this.sendVal = true; 
        // },
        cancelPopover(id){
            this.$refs[id].doClose();
        }, 
        clickSave(){
            this.getData()
        },  
        btnDefaultTemplate(scope){ //默认模板
            let params = {
                id:scope.row.id
            }
            API.btnDefaultTemplate(params).then(dt => {
                if(dt.success){
                        this.$message({
                            message:'默认模板设置成功',
                            type: 'success'
                        });
                        scope._self.$refs[`popover-${scope.row.id}`].doClose()
                        this.getData();
                    }else{
                        this.$message({
                            message:dt.error,
                            type: 'error'
                        });
                    }
            })
            .catch(error=>{});
        },
        btnMultiDeleteTemplate(){ //批量删
            if(this.multipleSelection.length == 0) return;
    		let _ids = [];
    		for (var i = 0; i < this.multipleSelection.length; i++) {
    			_ids.push(this.multipleSelection[i].id);
			}
            let params = {
                id:_ids.join(',') 
            }
            API.deleteReportTemplate(params).then(res => {
                if(res.success){
					this.$message({
                        message:'删除模板成功',
                        type: 'success'
                    });
                    this.alldelvisible = false
                    this.getData();
                }else{
					this.$message({
                        message:dt.error,
                        type: 'error'
                    });
				}
            })
            .catch((error) => {
                console.log(error);
            })
        },
        addTaskreport:function(_id,_name){  
			this.$router.push({
                path: `/addtaskreport`,
                 query: { 
                	isUpdate:false,
                    template_id:0,
                    title:'新建任务报告模板'
                }
            });
        },
        addTargetreport:function(_id,_name){  
			this.$router.push({
                path: `/addtargetreport`,
                query: { 
                	isUpdate:false,
                    template_id:0,
                    title:'新建目标报告模板'
                }
            });
		},
        fnDetails(row){ //报告模板详情
            let _path = '';
            if(row.type == 2){ //目标
                _path = '/addtargetreport';
            }else{ //任务
                _path = '/addtaskreport';
            }
            this.$router.push({
                path: _path,
				query: {  
                    isUpdate:true,
                    template_id:row.id,
                    title:row.name
                }
			});
        },
        mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
            let t = this.$refs['popover_id-' + row.id].showPopper;
            let t2 = this.$refs['popover-' + row.id].showPopper;
            // console.log('leave', t, t2)
            if(!t && !t2){
                  this.showEditFileNameButton = false;
                  this.rowId = "";
            }
          
        },
    }
})
 
</script>
