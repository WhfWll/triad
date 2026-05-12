//系统设置---系统工具---配置备份 页面代码
<template>
    <div>
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
                </div>
            </el-col>
        </el-row>
        <el-row :gutter="20" style="margin-top:16px;height: calc(100% - 153px);">
            <el-col :span="24" style="height: 100%;">
                <div class="totalbox" style="margin-bottom:10px">
                    <el-table class="apitable tablePadding" :data="TableData" style="width: 100%;" height="calc(100% - 50px)"
                     @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" v-loading = "Loading">
                        <el-table-column prop="name" label="备份文件名称">
                        </el-table-column>
                        <el-table-column prop="createTime" label="备份时间">
                        </el-table-column>
                        <el-table-column prop="" label="操作">
                            <template slot-scope="scope" > 
                                
                                <el-link  :underline="false" class="link_primary"  @click="getInfo(scope.row)">恢复</el-link>
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
import { config }from '@/api/system.js'
export default({
    name:'sysconfigcopy',
    components: {
  	},
    data(){ 
        //系统监控
    	return{
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
        this.$store.state.activefirstMenu="/systemsetting"; 
    },
    mounted:function(){
        this.getData();
        // this.getSelect();
        this.getConfig();
    },
    beforeDestroy(){
    },
    methods:{
        fnReptileSwitch(val){
            this.onlineUpdate.isopen = val;
        },  
        async getData(){//展示列表 
            let params = {
                page: this.currentpage,
                size: this.pageSize
            }
            const res = await config.getlogs(params);
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
            const res = await config.getconfiginfo();
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
        //恢复
        async getInfo(row){
            let params = {
                id: row.id
            }
            const res = await config.liveCopy(params);
            if(res.code === 200){  
                this.getData();
                this.$message.success('恢复成功')
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }

        },
        // async getSelect(){//获取下拉参数
        //     const res = await config.getselectinfo();
        //     if(res.success){  
        //         this.cycle = res.cycle;
                
        //     }else{
        //         this.$message({
        //             message:res.error,
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
            const res = await config.saveConfiginfo(params);
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
            const res = await config.saveCopyinfo();
            if(res.code === 200){  
                this.$message.success('备份成功');
                this.getData();
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        //下载
        async handleDown(row){
            let params = {
                id: row.id
            }
            const res = await config.download(params);
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
            const res = await config.handleDelinfo(params)
            if(res.code === 200){ 
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
        currentchange:function(t){
    		// this.page= 1; 
            this.currentpage = t;
            this.getData();
           
		},
		handleSizeChange(t){
			this.page_num = 1;
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
        height: 100%;
        padding: 24px  ;
        box-sizing: border-box;
        background: #fff;
        box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.1200);
        border-radius: 4px 4px 4px 4px;
        /deep/ .el-table{
           .cell{
            height: 15px;
            line-height: 15px;
           }
        }
    }
    .tabsbox {
        border-radius: 4px;
    }
    .systeminfolist{
    .planTime{
        height: 45px;
        line-height: 45px;
        display: flex;
        align-items: center;
    }
    .bottomline{
        display: flex;
        align-items: center;
        height:45px;
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

