<template>
  	<div > 
  		<div class="main-title  ">
	  		业务工具
	  	</div>
	  	<div class="tollboxlist context_box_bg">  
	  		<div class="search-box"> 
                <div class="operationbutton" > 
                    <!-- <el-button type="primary" @click="createTool" size="small">新建</el-button> --> 
                    <!-- <el-popover
                        popper-class="delButton_popper"
                        placement="bottom-start"
                        width="170" 
                        trigger="click" 
                        :visible-arrow="false"
                        v-model="alldelvisible" >
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="" >
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false" >取消</el-button>
                            <el-button size="mini" type="primary"  @click="AllDel">确定</el-button>
                        </div>  
                        <el-button type="warning"  size="small"  slot="reference" :disabled="!multipleSelection.length">删除</el-button> 
                    </el-popover>   -->
                </div>
                <div  class="serach-condition" >
                    <div > 
                        <el-select v-model="formData.type"  style=" width:150px;" placeholder="请选择工具类型" clearable size="small">  
                            <el-option
                                v-for="(item,index) in toolstype"
                                :key="index"
                                :label="item[1]"
                                :value="item[0]" 
                            >
                            </el-option>
                        </el-select> 
                    </div> 
                    <div class="search-text">
						<el-input placeholder="请输入关键字"   @keydown.enter.native="handlesearch"  v-model="formData.search_field" class="input-with-select"  size="small" clearable > </el-input>
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
				:data="tableData"  style="width: 100%"    
                @selection-change="handleSelectionChange"
                v-loading = "Loading"> 
                <!-- <el-table-column   
                  type="selection" width="55" :selectable='checkboxT'>
                </el-table-column> -->
				<el-table-column  
					prop="name"
					label="工具名称"  :show-overflow-tooltip="true"
					 > 
				</el-table-column> 
				<el-table-column
					prop="type"
					label="工具类型"
				 > 
                  <template slot-scope="scope" > 
                       <span  >{{scope.row.type[1]}}</span>
                  </template>
				</el-table-column>
                <el-table-column
                    prop="source"
                    label="脚本来源" > 
                 <template slot-scope="scope" >  
                        <span  >{{scope.row.source[1]}}</span> 
                    </template>
                </el-table-column>
				<el-table-column  
					prop="add_time"
					label="添加时间"  :show-overflow-tooltip="true">
				</el-table-column> 
				
				<el-table-column prop="status" label="操作"  >
					<template slot-scope="scope" >
						 
                        <!-- <el-link type="primary" :underline="false" v-if="scope.row.source[0] == '1'" disabled >编辑</el-link>  
                        <el-link type="primary" :underline="false" v-else @click="btnUpdate(scope.row)" >编辑</el-link>  
                        <el-link class="link_danger"  :underline="false" v-if="scope.row.source[0] == '1'" disabled >删除</el-link>  
                        <el-link class="link_danger"   :underline="false" v-else  @click="handleDel(scope.row.id,scope.row.name)" >删除</el-link>   -->
                        <el-link type="primary" :underline="false" @click="fnCheckinfo(scope.row.id)" >详情</el-link>  
					</template>
				</el-table-column>
			</el-table>
			<el-pagination
				:page-size="pageSize" 
				background
				layout="total,  prev, pager, next, sizes, jumper"
				:total="totalpage"
				:current-page="currentpage"
				@current-change = "currentchange"
                @size-change="handleSizeChange" >
			</el-pagination>

	  	</div>
	  	<el-dialog 
            :title="title" 
            :id="toolid" 
            :visible.sync="dialogFormVisible" 
            class="toolbox" 
            :before-close="canclDialog"  
            width='1184px' 
            :close-on-click-modal="false" 
            :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="submitDialog">确定</el-button>
                <el-button size="small" @click="canclDialog">取消</el-button>
            </div>
            <div style="padding:24px">
                <el-form :model="form" ref="ruleForm" label-width="100px" :rules="rules">
                    <div>
                        <el-form-item label="工具名称：" prop="name"  style="display: inline-block;width:50%"  >
                            <el-input v-model="form.name" autocomplete="off" placeholder="请输入工具名称..." maxlength="100"></el-input>
                        </el-form-item>
                        <el-form-item label="工具类型："  prop="type"  style="float:right;width:50%"  >
                            <el-select v-model="form.type" placeholder="请选择工具类型" style="width:100%">
                                <el-option
                                    v-for="(item,index) in toolstype"
                                    :key="index"
                                    :label="item[1]"
                                    :value="item[0]" 
                                >
                                </el-option>
                            </el-select>
                        </el-form-item>
                    </div>
                    <div class="clearfix">
                        <el-form-item label="测试强度：" prop = "rank" style="float:left;width:50%"  >
                            <el-select v-model="form.rank"   placeholder="请选择测试强度" style="width:100%">
                                <el-option
                                    v-for="(item,index) in level"
                                    :key="index"
                                    :label="item[1]"
                                    :value="item[0]">
                                    </el-option>
                            </el-select>
                        </el-form-item> 
                        <el-form-item label="工具对象：" style="float:left;width:50%">
                            <el-select v-model="form.object"   placeholder="请选择工具对象" style="width:100%">
                                <el-option
                                    v-for="(item,index) in objectlist"
                                    :key="index"
                                    :label="item[1]"
                                    :value="item[0]">
                                    </el-option>
                            </el-select>
                        </el-form-item> 
                    </div> 
                    <el-form-item label="工具描述："  >
                        <el-input type="textarea" v-model="form.desc"></el-input>
                    </el-form-item>
                    <div class="uploadbox" v-if="!toolid">
                        <label style="width:100px">脚本文件：</label>
                        <div class="uploadfile">
                            <el-link type="primary" @click="clickupload"><i class="el-icon-paperclip"></i>上传脚本</el-link>
                            <input type="file" name="" 
                                        class="btnUploadID"   
                                        @change="changeuploaID($event)" 
                                        style="display:none" 
                                        id="input-file-ID" 
                                    
                                        v-if="clearShow"   
                                        accept=".py" >
                            <p>目前仅支持python脚本文件</p>  
                        </div>  
                        <div style="padding-left: 100px;  float: left; height:40px; line-height: 40px;">
                                {{form.script_name}}
                        </div>
                    </div>
                    
                    <!-- <div>
                        <div class="filelist">
                            <div  v-for="(item,index) in uploadfileslist">
                                <label :title="item.script_name">{{item.script_name}}</label> 
                                <i class="iconfont iconshanchuwenjian" @click="clearfile(index)"></i>
                                <el-input  placeholder="请输入脚本描述..." v-model="item.description" class="filedesc"></el-input>
                            </div> 
                        </div>
                    </div> -->
                    
                </el-form>
            </div>
           
        </el-dialog>
  	</div>
</template>
<style scoped lang="less">
    /deep/ .el-table td:not(.el-table-column--selection):first-child .cell, 
    /deep/ .el-table th:not(.el-table-column--selection):first-child .cell{
        padding-left: 32px !important;
    }
.tollboxlist{
	padding: 24px; 
	background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
    border-radius: 4px;
}	 
.tollboxlist .el-tabs__content{
    min-height: 750px;
    // height: 100%;
}
.toolbox .el-dialog__body{
    max-height: 520px;
    overflow-y: auto;
}
.uploadbox{
    overflow:hidden;
}
.uploadbox > label{
    float: left;
    width:100px;
    text-align: center;
    height: 40px;
    line-height: 40px;

}
.uploadbox > label:before{
    content: '*';
    color: #F56C6C;
    margin-right: 4px;
}
.uploadbox >div.uploadfile{
    float: left;
    width: calc(100% - 110px);
    border-bottom:1px solid #dcdfe6;
    height: 40px;
    line-height: 40px;
}

.uploadbox >div.uploadfile p{
    display: inline-block;
    color: #909399;
    font-size: 12px;
}
.filelist{
    margin-left: 100px;
    width: calc(100% - 110px);
    margin-top: 10px;
}
.filelist > div{
    margin:5px 0;
}
.filelist > div >label{
    width: 150px;
    display: inline-block; 
    overflow:hidden;
    white-space:nowrap;
    text-overflow:ellipsis;
}
.filelist > div >i{ 
    cursor: pointer;
    display: inline-block;

}
.filelist > div > .filedesc{
    height: 40px;
    line-height: 40px;
    margin-left: 20px;
    display: inline-block;
    width: calc(100% - 198px); 
}
.dialog_b_btn{
    position: absolute;
    top: 15px;
    right: 24px;
    font-size: 14px;
    button{
        color: #4C7AE3;
    } 
}
    @media (max-width: 1440px) {
        
    /deep/ .el-dialog{
        height: calc(100% - 96px);
    }
}
@media  (min-width: 1440px) { 
    /deep/ .el-dialog{
        height: calc(100% - 176px);
    }
}
.dialog_item_label{
    font-size: 14px;
    border-left: 3px solid #4C7AE3;
    padding-left: 8px; 
    font-weight:500;
    width: 113px;
    display: inline-block;
    // height: 18px;
    line-height: 16px; 
    box-sizing: border-box;
} 
</style>
<script>  
import {encryptCBC,decryptCBC} from '@/commonFunction/des.js'
import xzbutton from "@/components/XzButton.vue";
export default({
    name:'toolmanagement',
    components: {
    	xzbutton,
  	},
    data(){ 
        var valiFile = (rule, value, callback) => {   
            if(this.addform.hasOwnProperty(rule.field)){
                if (!this.addform[rule.field]) {
                    callback(new Error('请上传脚本文件'));
                } else {
                    callback();
                }
            } 
            
        };
    	return{  
    		totalpage:0,
    		currentpage:this.$route.query.page_num ? parseInt(this.$route.query.page_num) :1,
    		formData:{
    			page_num:this.$route.query.page_num ? parseInt(this.$route.query.page_num) :1,
    			type:'', 
    			search_field:'',
    		},
    		status:[
				[
		            0,
		            "全部"
		        ],
		         
			],
            typelist:[],
			tableData:[], 
    		Loading:false,
            multipleSelection:[],
            dialogFormVisible:false,
            formLabelWidth:'100px',
            form:{
                name:'',
                type:'',
                object:'',
                introduce:'',
                inputstr:'',
                output:'',
                file:null,
                script_name:'',
            },
            files:[
                {
                    name:'',
                    desc:''
                }
            ] ,
            rules:{
                name:[
                    { required: true, message: '工具名称不能为空', trigger: 'blur' }, 
                ],
                type:[
                    {required: true, message: '请选择工具类型', trigger: 'change' }, 
                ],
                rank:[
                    {required: true, message: '请选择测试强度', trigger: 'change' }, 
                ],
                file:[
                   { required: true, validator: valiFile, trigger: "change" }, 
                ]
            }, 
            valiFile:valiFile,
            clearShow:true,
            uploadfileslist:[],
            toolid:'',
            title:'',
            source:[],
            toolstype:[],
            objectlist:[],
            toolstatus:[],
            level:[],
            myKey: '4dogs.cn',
            pageSize :10,
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/toolmanagement"; 
    },
    mounted:function(){  
     	this.getTooltype(); 
    	this.getData();
    },
    watch:{ 
        '$route': function(to, from){ // 路由改变时执行     
            this.currentpage = this.$route.query.page_num ? parseInt(this.$route.query.page_num) :1;
            this.formData.page_num = this.$route.query.page_num ? parseInt(this.$route.query.page_num) :1;
            this.pageSize = this.commonjs.pageSize;  
            this.getData(); 
        }
    },
    methods:{ 
        checkboxT(row,index){
            if(row.source[0]== '1'){
                return 0;
            }else{
                return 1;
            }
        },
    	getTooltype:function(){ 
            this.$ajax.get('/tools/interfaces/params/',{
                params: {}
            }).then((data) => { 
                var res = data.data;   
                if(!res.error){  
                    this.toolstype = res.type;
                    this.objectlist = res.object;
                    this.source = res.source;
                    this.toolstatus = res.status; 
                    this.level = res.rank; 
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
            })
            .catch(function(error){
                console.log(error);
            })
    	}, 
    	getData:function(){ 
            this.Loading = true; 
            this.$ajax.get('/tools/interfaces/list/',{
                params: {
                    page: this.formData.page_num, 
  					name:this.formData.search_field, 
                      type:this.formData.type, 
                      page_size:this.pageSize
                }
            })
            .then((data) => { 
                var dt = data.data;    
               
                this.Loading = false;
                this.tableData= dt.results;
                this.totalpage = dt.count;
 
            })
            .catch(function(error){
                console.log(error);
            })
        },
        handleReset(){
            this.formData.search_field = '';
            this.formData.type = '';
            this.formData.page_num = 1; 
            this.currentpage = 1;
            this.pageSize = 10;
            this.getData();
        },
        handlesearch(){
            this.formData.page_num = 1; 
            this.getData();
            this.currentpage = 1;
         

        }, 
    	currentchange:function(t){
    		this.formData.page_num = t; 
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t){
            this.formData.page_num = 1;
            this.pageSize = t;
            this.getData();
        },
    	fnCheckinfo:function(id){   //查看
            // this.$router.push({
            //     path: `/toolinfo/${id}`,
            // });

            this.$router.push({
                path: `/toolinfo`,
                query: { 
                    id: id, 
                    page_num:this.formData.page_num
                }
            });
         
        },
        AllDel(){
    	   if(this.multipleSelection.length == 0) return;
            var ids = [];
            for (var i = 0; i < this.multipleSelection.length; i++) {
                ids.push(this.multipleSelection[i].id);
            }
            this.handleDel(ids.join(','));
        },
        createTool(){
            this.dialogFormVisible = true;
            this.toolid = '';
            this.title='创建工具';
            this.rules  = {
                name:[
                    { required: true, message: '工具名称不能为空', trigger: 'blur' }, 
                ],
                type:[
                    {required: true, message: '请选择工具类型', trigger: 'change' }, 
                ],
                rank:[
                    {required: true, message: '请选择测试强度', trigger: 'change' }, 
                ],
                file:[
                   { required: true, validator: this.valiFile, trigger: "change" }, 
                ]
            };
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        submitDialog(id){ 
            if(!this.toolid){
                this.submitcreateTool();
            }
            else{
                this.submitupdateTool();
            }
        },
        submitcreateTool(){
            let formData = new FormData(); 
            this.$refs.ruleForm.validate((valid) => { 
                if (valid) {   
                    let formData = new FormData(); 
                    formData.append('name',this.form.name); 
                    formData.append('type',this.form.type); 
                    formData.append('object',this.form.object); 
                    formData.append('rank',this.form.rank); 
                    formData.append('detail',this.form.desc); 
                    formData.append('user', decryptCBC(localStorage.getItem('user_id'),this.myKey)); 
                    formData.append('status',0); 
                    formData.append('file',this.form.file); 
                    let config = {
                        headers: {
                            'Content-Type': 'multipart/form-data'
                        }
                    }
                    this.$ajax.post('/tools/interfaces/create/',
                        formData, config
                    ).then(data => { 
                        var dt = data.data;  
                        if(dt.success){ 
                            this.$message({
                                message:'新增工具成功',
                                type: 'success'
                            });
                            this.dialogFormVisible = false; 
                            this.form.name = '';
                            this.form.type='';
                            this.form.object='';
                            this.form.rank='';
                            this.form.file=null;
                            this.form.desc='';
                            this.form.script_name = '';
                            document.getElementById('input-file-ID').value = null;
                            this.$refs['ruleForm'].resetFields(); 
                            
                            this.getData();
                        }else{
                            this.$message({
                                message:dt.error,
                                type: 'error'
                            });
                        }  
                        
                    })
                    .catch(data=>{
                        console.log(data); //错误信息
                    });

                }
            });
        },
        submitupdateTool(){
            let formData = new FormData(); 
            this.$refs.ruleForm.validate((valid) => { 
                if (valid) {   
                    this.$ajax({
                        method:'patch',
                        url:'/tools/interfaces/update/',
                        data: {
                        id:this.toolid ,
                            name:this.form.name,
                            type:this.form.type,
                            object:this.form.object,
                            rank:this.form.rank, 
                            user:this.form.user, 
                            status:this.form.status,
                        } 
                    })
                    .then(data => { 
                        var dt = data.data;  
                         if(!dt.error){  
                            this.dialogFormVisible = false;
                            this.$message({
                                message:'修改工具成功',
                                type: 'success'
                            });
                            this.getData();
                            this.form.name = '';
                            this.form.type='';
                            this.form.object='';
                            this.form.rank=''; 
                            this.form.desc=''; 
                            this.$refs['ruleForm'].resetFields(); 
                        }else{
                            this.$message({
                                message:dt.error,
                                type: 'error'
                            });
                        }  
                        
                    })
                    .catch(data=>{
                        console.log(data); //错误信息
                    }); 

                }
            });
        },
        canclDialog(){
            this.dialogFormVisible = false;

            this.$refs.ruleForm.resetFields(); 
            this.form.name = '';
            this.form.type = '';
            this.form.introduce = '';
            this.form.inputstr = '';
            this.form.output = '';
            this.uploadfileslist=[];
        },
        clickupload(){   
            document.querySelector('.btnUploadID').click();
        },
        changeuploaID:function(e){  
            let deviceFile = e.target.files;  
            for(let i=0;i<deviceFile.length;i++){  
                this.form.file = deviceFile[i];
                this.form.script_name = deviceFile[i].name 
            }  
        },
        // clearfile(index){
        //     this.uploadfileslist.splice(index,1);
        //     document.getElementById('input-file-ID').value = null
        // },
        handleDel(id,name){
            var des = name ? '确定要删除选择的工具【'+name+'】吗？' :'确定要删除选择的工具吗？' ;
            this.$confirm(des, '删除消息', {
                distinguishCancelAndClose: true,
                confirmButtonText: '确定',
                cancelButtonText: '取消',  
            }).then(() => {  
                this.$ajax({
                    method:'delete',
                    url:'/tools/interfaces/delete/',
                    data: {
                        id:id
                    } 
                })
                .then(data => { 
                    var dt = data.data;  
                    if(dt.success){ 
                        this.$message({
                            message:'删除工具成功',
                            type: 'success'
                        });
                        this.getData();
                    }else{
                        this.$message({
                            message:dt.error,
                            type: 'error'
                        });
                    }  
                    
                })
                .catch(data=>{
                    console.log(data); //错误信息
                });
            }).catch(action => {
                      
            }); 
        },
        btnUpdate(row){
            this.dialogFormVisible = true;
            this.toolid = row.id;
            this.title='编辑工具';

            this.form.name = row.name;
            this.form.type= row.type[0];
            this.form.object=row.object[0];
            this.form.rank= row.rank[0];
       
            this.form.desc=row.desc;
            this.form.script_name = '';

            this.rules = {
                name:[
                    { required: true, message: '工具名称不能为空', trigger: 'blur' }, 
                ],
                type:[
                    {required: true, message: '请选择工具类型', trigger: 'change' }, 
                ],
                rank:[
                    {required: true, message: '请选择测试强度', trigger: 'change' }, 
                ], 
            }; 
        },
       
    }
})
 
</script>
