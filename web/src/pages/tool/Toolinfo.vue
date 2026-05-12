<template>
  	<div > 
  		<!-- <div class="main-title  "> 
            <div class="goback" @click="goBack">
                <i class="el-icon-arrow-left"></i>
            </div>
            <label>工具详情</label> 
	  	</div>  -->
        <div class="main-title  ">   
            <router-link :underline="false" class="classA" :to="{ path: '/toolmanagement' }" >工具管理</router-link>  
            <label class="currentpagetitle">
                <span>工具详情</span> 
            </label>
	  	</div> 
	  	<div class="toolinfoboxlist"> 
	  		<div class="toolbasicinfo">
	  			<div>
	  				<strong>工具名称：</strong>
	  				<span>{{toolinfo.name}}</span>
	  			</div>
	  			<div>
	  				<strong>工具类型：</strong>
	  				<span>{{toolinfo.type[1]}}</span>
	  			</div>
                <div>
	  				<strong>工具对象：</strong>
	  				<span>{{toolinfo.object[1]}}</span>
	  			</div>
                <div>
	  				<strong>工具来源：</strong>
	  				<span>{{toolinfo.source[1]}}</span>
	  			</div>
                <div>
	  				<strong>工具状态：</strong>
	  				<span>{{toolinfo.status[1]}}</span>
	  			</div>
	  			<div>
	  				<strong>集成时间：</strong>
	  				<span>{{toolinfo.add_time}}</span>
	  			</div>
	  		</div>
	  		<div class="infodiv">
	  			<div class="part_title">工具介绍</div>
	  			<div class="info">
	  			 	{{toolinfo.detail}}
	  			</div>
	  		</div>
            <div class="infodiv">
                <div class="part_title">工具输入</div>
                <div  class="info" style="white-space:pre-line;">
                    {{toolinfo.input}}
                </div>
            </div>
      		<div class="infodiv">
      			<div class="part_title">工具输出</div>
      			<div  class="info" style="white-space:pre-line;">
      				{{toolinfo.output}}
      			</div>
      		</div>
            <!-- <div class="infodiv">
                <div class="part_title">脚本文件</div>
                <el-table
                    :data="tableData"
                    stripe
                    style="width: 100%">
                    <el-table-column
                      prop="script_name"
                      label="脚本名称"  :show-overflow-tooltip="true"
                      >
                    </el-table-column>
                    <el-table-column
                      prop="description"
                      label="脚本描述"  :show-overflow-tooltip="true">
                    </el-table-column> 
                    <el-table-column prop="status" label="操作"  >
                        <template slot-scope="scope" >
                            
                            <el-link type="primary" :underline="false" @click="handleDel(scope.row.script_id,scope.row.script_name)">删除</el-link>
                        </template>
                    </el-table-column>
                </el-table> 
            </div> -->
	  	</div> 
  	</div>
</template>
<style  scoped lang="less">
.toolinfoboxlist{
	padding: 24px; 
	background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
}	
.toolbasicinfo{
	background-color:#F7F7FB !important;
	border:1px solid #e2e5ed;
	padding: 16px 0;
}
 .toolbasicinfo >div{
 	min-width: 200px;
 	display: inline-block;
     font-size: 13px;
     line-height: 28px;
     font-weight: 500;
     padding: 0 24px ;
     strong{
         color: rgba(72, 72, 102, 0.87);
        font-weight: 500;
     }
     span{
         color: rgba(72, 72, 102, 0.64);
         font-weight: 500;
     }
 }
 .infodiv{
 	margin: 24px 0;
 }
  .infodiv > strong{
        display: block;
        margin: 20px 0 24px;
        font-size: 13px;
        color: rgba(72, 72, 102, 0.87);
        font-weight: 500;
  }
  .infodiv > .info{ 
      font-size: 13px;
      color: rgba(72, 72, 102, 0.64);
    padding-left: 14px;
  }
/deep/ .el-table td:not(.el-table-column--selection):first-child .cell, 
/deep/ .el-table th:not(.el-table-column--selection):first-child .cell{
    padding-left: 32px !important;
}
.part_title{
    font-size: 13px;
    margin-bottom:16px;
    font-weight: 500;
    border-left: 3px solid #4C7AE3;
    padding-left: 10px;
    height: 14px;
    line-height: 14px;
    color:rgba(72,72,102,0.87);
}
</style>
<script>  
export default({
    name:'toolinfo',
    data(){ 
    	return{  
    		id:this.$route.query.id,
            fpage_num : this.$route.query.page_num,
    		toolinfo:{
    			name:'',
    			type:'',
    			output:'',
    			input:'',
    			detail:'',
    			add_time:'',
    		},
            tableData:[],
            currentpage:1,
            totalpage:0,
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/toolmanagement";  
    },
    mounted:function(){ 
    	this.getDataInfo();  
    },
    methods:{ 
    	getDataInfo:function(){  
            this.$ajax.get('/tools/interfaces/detail/',{
                params: {
                    id:this.id, 
                }
            })
            .then((data) => { 
                var dt = data.data;
               if(!dt.error){  
 					this.toolinfo = dt ;
                   
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
    	goBack() {
	       // this.$router.go(-1);
            this.$router.push({
                path: `/toolmanagement`,
                query: { 
                    page_num: this.fpage_num, 
                }
            });
	    },
        currentchange(){

        }, 
        handleDel(id,name){ 
            this.$confirm('确定要删除脚本【'+name+'】吗？', '删除消息', {
                distinguishCancelAndClose: true,
                confirmButtonText: '确定',
                cancelButtonText: '取消',  
            }).then(() => {  
                this.$ajax({
                    method:'post',
                    url:'/tools/delete_script/',
                    data:this.qs.stringify({
                        script_id:id
                    })
                })
                .then(data => { 
                    var dt = data.data;  
                    if(dt.success){ 
                        this.$message({
                            message:'删除脚本成功',
                            type: 'success'
                        });
                        this.getDataInfo();
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
        }
        
    }
})
 
</script>
