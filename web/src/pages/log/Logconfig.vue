//日志配置
<template>
    <div>
        <div class="main-title  ">
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i>
            <label class="taskSceneBtn" >日志配置</label>
            
        </div>
        <el-row :gutter="20">
            <el-col :span="24">
                <div class="totalbox">
                    <div class="tabsbox productinfo sysupgradebox systeminfolist">
                             <div class="planTime">
                                <label>周期备份</label>
                                <el-switch
                                    class="switchbtn"
                                    active-color="#4c7ae3"
                                    v-model="onlineUpdate.isopen"
                                    inactive-color="#E8E8F5"
                                    @change="fnReptileSwitch">
                                </el-switch>
                                <el-select v-model="value" placeholder="请选择" class="selectbox"> 
                                    <el-option
                                      v-for="item in cycle"
                                      :key="item.value"
                                      :label="item.label"
                                      :value="item.value">
                                    </el-option>
                                </el-select><span class="month">月</span>
                                <el-button type="primary" size="small"  @click="saveCongig()">保存配置</el-button>
                           </div>
                    </div>
                    <div class="tabsbox productinfo sysupgradebox systeminfolist">
                        <div class="bottomline">
                            <label>立即备份</label>
                            <el-button type="primary" size="small" @click="saveCopy()" style="margin-left:50px">执行备份</el-button>
                        </div>
                    </div>
                    <div class="tabsbox productinfo sysupgradebox systeminfolist">
                        <div class="bottomline">
                            <label>日志保留时间</label>
                            <el-input @input="inputChange" v-model="savaTime" placeholder="单位/天"  style="margin-left:30px;width:100px"></el-input><span style="margin-left:10px">天</span>
                            <el-button type="primary" size="small" @click="setLog()" style="margin-left:50px">保存</el-button>
                            
                        </div>
                    </div>
                </div>
            </el-col>
        </el-row>
        <el-row :gutter="20" style="margin-top:15px;height: calc(100% - 224px);">
            <el-col :span="24" style="height: 100%;">
                <div class="totalbox totalbox2">
                    <el-table class="apitable tablePadding" :data="TableData" style="width: 100%;"  height="calc(100% - 60px)"
                        @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" v-loading = "Loading">
                        <el-table-column prop="name" label="备份文件名称">
                        </el-table-column>
                        <el-table-column prop="createTime" label="备份时间">
                                <template slot-scope="scope" > 
                                    <div v-if="showEditFileNameButton && rowId == scope.row.id">
                                            <!-- <el-link  :underline="false" class="link_primary"  @click="getInfo(scope.row)">恢复</el-link> -->
                                        <el-link  :underline="false" class="link_primary"  @click="handleDown(scope.row)">下载</el-link>  
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
                                            <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference">删除</el-link>  
                                        </el-popover>
                                        <!-- <el-link :underline="false" class="link_danger linkafter" style="padding:0" disabled  v-else>删除</el-link>  -->
                                    </div>
                                    <div v-else >
                                        <span>{{scope.row.createTime}}</span>
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
                    @current-change = "currentchange"
                    @size-change="handleSizeChange" >
                </el-pagination>
                </div>
            </el-col>
        </el-row>
    </div>
</template>

<script> 
import log from '@/api/log.js'
// import jsFileDownload from 'js-file-download'
export default({
    name:'logconfig',
    components: {
  	},
    data(){ 
        //系统监控
    	return{
            savaTime:"",
            cycle: [
                {
                value: 1,
                label:1
                }, {
                value: 2,
                label: 2
                },
                 {
                value: 3,
                label: 3
                }, {
                value: 6,
                label: 6
                }, {
                value: 12,
                label: 12
                }
            ],
            value: 1,
            onlineUpdate: {
                plan_time2: 1,
                isopen:false
            },
            Loading:false,
            TableData:[],
            showEditFileNameButton:false,
            rowId:'',
            pageSize:10,
            totalpage:0,
			currentpage:1, 
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/logconfig"; 
    },
    mounted:function(){
        this.getData();
        // this.getSelect();
        this.getConfig();
        this.getexptime()
    },
    beforeDestroy(){
    },
    methods:{
    inputChange(val) {
      this.savaTime = this.savaTime.replace(/[^0-9.]/g, '')
    },
        fnReptileSwitch(val){
            this.onlineUpdate.isopen = val;
        },  
        async getData(){//展示列表 
            let params = {
                page: this.currentpage,
                size: this.pageSize
            }
            const res = await log.getlogs(params);

            if(res.code === 200){  
                this.Loading = false;
                this.TableData = res.data.list;
                this.totalpage = res.data.total;
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        async getConfig(){//进来先获取配置
            const res = await log.getconfiginfo();
            if(res.code === 200){  
                this.value = res.data.cycle;
                if(res.data.isOpen == 1){
                    this.onlineUpdate.isopen = true;
                }else{
                    this.onlineUpdate.isopen = false;
                }
                
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        // async getSelect(){//获取下拉参数
        //     const res = await log.getselectinfo();
        //     if(res.success){  
        //         this.cycle = res.cycle;
                
        //     }else{
        //         this.$message({
        //             message:res.msg,
        //             type: 'error'
        //         });
        //     }

        // },
        //保存配置
        async saveCongig(){
            let params = {
                cycle: this.value,
                isOpen: this.onlineUpdate.isopen ? 1 : 2
            }
            const res = await log.saveConfiginfo(params);
            if(res.code === 200){  
                this.getData();
                this.$message.success('保存配置成功')
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }

        },
        //执行备份
        async saveCopy(){
            const res = await log.saveCopyinfo();
            if(res.code === 200){  
                this.$message.success('备份成功')
                this.getData();
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        },
        //获取日志信息
        async getexptime(){
            const res = await log.getexptime();
            if(res.code === 200){  
                this.savaTime  = res.data.expirationTime
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        },
        setLog(){
            if (this.savaTime <0) {
                this.$message.error("日志保留时间不能为负")
                return
            }else{
                this.setexptime()
            }
        },
        //baocun日志信息
        async setexptime(){
            const res = await log.setexptime({
                expirationTime:this.savaTime
            });
            if(res.code === 200){  
                this.$message.success("保存日志信息成功")
                this.getexptime()
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        },
        //下载
        async handleDown(row){
            let params = {
                id: row.id
            }
            const res = await log.download(params); 

            if(res.code){  // 有code码返回，认为是失败
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }else{
                const blob = new Blob([res],{
                    type:"application/octet-stream"
                });
                const url =  window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                // a.download = row.name + '.csv';
                a.download = row.name + '.csv';
                a.href = url;
                a.click();
            }
        },
         // 删除规则
         async handleDel (scope){ //删除 
            let params = {
                id:scope.row.id
            }
            const res = await log.handleDelinfo(params)
            if(res.code === 200){ 
                this.$message({
                    message:'删除成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                this.getData();
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
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
    }
})
 
</script>

<style scoped lang="less">
    .selectbox{
        width:90px;
        margin:0 10px 0 40px;
    }
    .month{
        display: inline-block;
        color: rgba(72, 72, 102, 0.87);
        margin-right: 40px;
    }
    .totalbox{
        padding:24px;
        background: #fff;
        box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.1200);
        border-radius: 4px 4px 4px 4px;
        &.totalbox2{
            margin-bottom: 10px;
        }
        /deep/ .el-table{
           .cell{
            height: 15px;
            line-height: 15px;
           }
        }
    }
    .totalbox2{
        padding:20px;
        height: calc(100% - 39px  );
    }
    .tabsbox {
        border-radius: 4px;
    }
    .systeminfolist{
    .planTime{
        height: 40px;
        line-height: 40px;
        display: flex;
        align-items: center;
    }
    .bottomline{
        display: flex;
        align-items: center;
        height:40px;
    }
    .switchbtn{
        margin-left: 50px;
    }
    .blueword{
        color:#4C7AE3;
        cursor: pointer;
    }
    .timeword{
        margin-left: 10px;
    }
    .system{
      
        font-weight: 500;
        line-height: 30px;
        padding: 10px 20px;
        padding: 10px 20px;
        line-height: 30px;
        width: calc(100% - 160px);
        float: left;
        box-sizing: border-box;
        .system_data{
            background: #fff;
            color: rgba(72, 72, 102, 0.64);
            font-weight: 500;

        }
    }
} 
.systeminfolist>div {
    /* border-top: 1px solid #E8E8F5; */
    border-radius: none;
    overflow: hidden;
    // background: #f5f8fb;
    font-size: 13px;
}
.systeminfolist >div >label,
.systeminfolist >div >span{
    line-height: 30px;
    float: left;
    box-sizing: border-box;
 
}
.systeminfolist>div>label {
    height: 15px;
    line-height: 15px;
    border-left: 3px solid #4C7AE3;
    padding-left: 8px;
    display: inline-block;
    /* width: 160px; */
    text-align: center;
    /* border-right: 1px solid #E8E8F5; */
    /* background: #F7F7FB !important; */
    color: rgba(72, 72, 102, 0.87);
    font-weight: 500;
}
.systeminfolist>div>span {
    display: block;
    color: rgba(72, 72, 102, 0.64);
    font-weight: 500;
}
</style>
