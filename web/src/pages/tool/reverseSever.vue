<template>
    <div>
         <div class="totalbox">  
            <div class="tabsbox  ">
                <div class="planTime">
                    <label style="margin-right: 10px;">反连地址</label>
                    <el-input v-model="form.host" size="small" style="width: 200px;margin-right: 10px;"></el-input>
                    <label style="margin-right: 10px;">反连端口</label>
                    <el-input v-model="form.port"  size="small" style="width: 200px;margin-right: 10px;"></el-input>
                    <el-button type="primary" @click="btnOpen">启动反连服务器</el-button>
                       <el-button type="primary"  style="float: right;" @click="btnStop">关闭反连服务器</el-button> 
                </div>
            </div>  
        </div>  
        <div class="totalbox boxlist">
            <div class="search-box"> 
                <el-button type="primary" @click="clearMsg" >清空</el-button>
                <span v-if="reverseURL" style="margin-left: 20px;">http反连地址：{{ reverseURL }}</span>
            </div>
             <el-table :data="tableData"  height="calc(100% - 102px)"
                @selection-change="handleSelectionChange"   
                > 
                <el-table-column prop="reverseType" label="反连类型" show-overflow-tooltip> </el-table-column>
                <el-table-column prop="remoteAddr" label="连接来源" show-overflow-tooltip> </el-table-column>
                <el-table-column prop="token" label="TOKEN" show-overflow-tooltip> </el-table-column>
                <el-table-column prop="response" label="响应" show-overflow-tooltip> </el-table-column> 
                 
            </el-table>
             <el-pagination  :page-size="param.size" background layout="total, prev, pager, next, sizes, jumper"
                :total="total" :current-page="param.page" @current-change="currentchange" @size-change="handleSizeChange">
            </el-pagination>

        </div>
    </div>
</template>
<style lang="less" scoped>
    .totalbox{
        padding: 24px ;
        box-sizing: border-box;
        background: #fff; 
        box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.12);
        border-radius: 4px 4px 4px 4px; 
        margin-bottom: 16px;
        
    }
    .tabsbox {
        border-radius: 4px;
    }
    .boxlist{
        padding: 20px;
        box-sizing: border-box;
        height: calc(100% - 150px);
    }
</style>
<script>
import { auxiliarytool } from '@/api/tool.js'
export default {
    name:'reverseSever',
    data(){
        return{
            form:{
                host:'',
                port:'',
            },
             param:{
                size:10,
                page:1,
            },
            total:0,
            ids:[],
            tableData:[],
            reverseURL:'',
            timer:null,
        }
    },
    mounted(){
       this.getData();
        this.reverseStatus();
        this.timer = setInterval(()=>{
            this.reverseStatus();
            this.getData();
        },5000)
    },
    beforeDestroy() { 
        clearInterval(this.timer);
        this.timeer = null;
    },
    methods:{
        async reverseStatus(){
            const res = await auxiliarytool.reversestatus()
            if(res.code == 200){
               if(res.data.status){
                    this.reverseURL = res.data.reverseUrl
               }else{
                    clearInterval(this.timer);
                    this.timeer = null;
                    this.reverseURL ='';
                    if(res.data.errInfo){
                        this.$message({
                            message:res.data.errInfo,
                            type: 'error'
                        });
                    }
                    
               }
            }else{
                this.$message({
                    message:'获取反连服务器状态失败',
                    type: 'error'
                });
            }
        },
        async btnOpen(){
            if(this.form.host == '' || this.form.port ==''){
                this.$message({
                    message:'反连服务器/反连端口 不能为空',
                    type: 'error'
                });
                return;
            }
             const res = await auxiliarytool.reversestart({
                host:this.form.host,
                port:this.form.port
            })
            if(res.code == 200){
                 this.$message({
                    message:'开启反连服务器成功',
                    type: 'success'
                });
                this.form.host = ''
                this.form.port = ''

                clearInterval(this.timer);
                this.reverseStatus();
                this.timer = setInterval(()=>{
                    this.reverseStatus();
                    this.getData();
                },5000)
            }else{
                 this.$message({
                    message:'开启反连服务器失败',
                    type: 'error'
                });
            }
        },
        async btnStop(){
            const res = await auxiliarytool.reversestop()
            if(res.code == 200){
                this.$message({
                    message:'关闭反连服务器成功',
                    type: 'success'
                });
                this.reverseURL ='';
            }else{
                this.$message({
                    message:'关闭反连服务器操作失败',
                    type: 'error'
                });
            }
        },
        async getData(){
            const res = await auxiliarytool.reversemessage({
                page:this.param.page,
                size:this.param.size
            })
            if(res.code == 200){
                this.tableData = res.data.list;
                this.total = res.data.total;
            }
        },
        async clearMsg(){ //清空
            const res = await auxiliarytool.reverseclear()
            if(res.code == 200){
                this.$message({
                    message:'清空信息成功',
                    type: 'success'
                });
            }else{
                this.$message({
                    message:'清空数据失败',
                    type: 'error'
                });
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
    }
}
</script>