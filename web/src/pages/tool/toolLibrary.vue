<template>
    <div>
        <div class="objectlist">  
			<el-table
				:data="tableData"  style="width: 100%" class="myTable"  
				 @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave"> 
			    <el-table-column   
					prop="name"
					label="名称"   :show-overflow-tooltip="true">   
                </el-table-column>
				<el-table-column  
					prop="fileType"
					label="类型"  > 
                    
				</el-table-column>
				 
                <el-table-column  
					prop="createTime"
					label="创建时间"  >
                    <template slot-scope="scope" > 
                        <div v-if="showEditFileNameButton && rowId == scope.row.id">
                            <el-link  :underline="false" class="link_primary"  @click="btnDown(scope.row)">下载</el-link>  
                        </div>
                        <div v-else >
                            <span>{{scope.row.createTime}}</span>
                        </div> 
					</template>

				</el-table-column> 
				<!-- <el-table-column prop="status" label="操作" min-width="10%"  >
				</el-table-column> -->
			</el-table>
			<el-pagination
				:page-size="pageSize" 
				background
				layout=" total,  prev, pager, next,sizes, jumper"
				:total="totalpage"
				:current-page="currentpage"
				@current-change = "currentchange"
				@size-change="handleSizeChange" >
			</el-pagination>
        </div>
    </div>
</template>
<style lang="less" scoped>
</style>

<script>
import { auxiliarytool } from '@/api/tool.js'
import bas from '@/api/bas.js'
export default ({
    name:'toolLibrary',
    components: { 
    },
    data(){ 
        return{
            pageSize:10,
            totalpage:0,
            currentpage:1,
            tableData:[],
            showEditFileNameButton:false,
            rowId:0,
            formData:{
                page_num:1,
            }
        }
    },
    created() { 
        this.pageSize = this.commonjs.pageSize;
    },
    methods:{
        async getData(){
            const res = await auxiliarytool.toollist({
                page:this.formData.page_num,
                size:this.pageSize
            });
            if(res.code == 200){
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }

        },
        currentchange(t){
            this.formData.page_num = t; 
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t){
            this.formData.page_num = 1;
            this.pageSize = t;
            this.getData();
        },
        mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
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
            }
        },
        async btnDown(row){ //下载

            const res = await bas.downagent({
                getTempToken:true
            })
            if(res.code ==200){
                let linkhref = '';
                if(process.env.NODE_ENV == 'development'){
                    linkhref += 'http://'+process.env.VUE_APP_API_URL; 
                }else{
                    linkhref += window.location.protocol + "//" + window.location.hostname  ;
                }
                linkhref+='/api/smart/tools/toolfiledownload?temp_token='+res.data.tempToken+'&filePath='+row.filePath+'&name='+row.name;

                if(res.data.tempToken){
                    window.open(linkhref, '_blank')
                }
                else{
                    this.$message({
                        message: '下载内容有误',
                        type: 'error'
                    });
                }
            }
        },
    },
    
})
</script>
