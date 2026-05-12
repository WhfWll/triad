/* 

    渗透任务详情页
 */
<template>
    <div>
        <div class="main-title">
            <router-link :underline="false" class="classA" :to="{ path: '/vulScanTask' }" >漏扫任务
            </router-link>  
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark"  :content="task_name"  placement="bottom">
                    <span>{{task_name}}</span>
                </el-tooltip>
            </label>
        </div> 
        <div class="taskinfolist  ">
            <div class="iconlist">
                <ul>
                    <li class="status">
                        <!--   待触发：1 ；待执行：2；运行中：3；已完成：4 ；暂停中：5。 -->
                        <!-- <span class="tag_status tag_danger1" v-if="status == 1">{{ statusName }}</span> -->
                        <span class="tag_status tag_warning" v-if="status == 2">{{ statusName }}</span>
                        <span class="tag_status tag_primary" v-if="status == 3">{{ statusName }}</span>
                        <span class="tag_status tag_success" v-if="status == 4">{{ statusName }}</span>
                        <span class="tag_status tag_danger" v-if="status == 5">{{ statusName }}</span>
                    </li>
                </ul>
            </div>
         
            <el-tabs v-model="activeName" @tab-click="handleClick" >
                <el-tab-pane label="任务概览" name="tabs1">
                    <detailoverview @actLD="actLD" @act ='act' ref="overview" :task_name="task_name" :task_id=task_id ></detailoverview> 
                </el-tab-pane>
                <el-tab-pane label="测试目标" name="tabs2">
                    <targetlist  ref="targetlist" :task_id=task_id :task_name=task_name></targetlist> 
                </el-tab-pane>
                <el-tab-pane label="漏洞测试" name="tabs4"> 
                    <vulnmsg  ref="vulnmsg" :task_id=task_id></vulnmsg>
                </el-tab-pane> 
        
            </el-tabs>
        </div>



    </div> 
</template>
<style lang="less" scoped>
     .taskinfolist{ 
        // height: calc(100% - 39px);
        min-height: calc(100% - 39px);
        // overflow-y: auto;
        box-sizing: border-box;  
        position: relative; 
        /deep/ .el-tabs{
            height: 100%; 
            .el-tabs__content{
                height: calc(100% - 64px);
                overflow: inherit;
                >div{
                    height: 100%;
                    // box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
                    // border-radius: 4px;
                    // margin-bottom: 10px;
                }
            }
        } 
    }	
    .iconlist{
        position: absolute;
        left: 1010px;
        top: 14px;
        z-index: 9; 
    }
    .dynamicfeed{
        position: absolute;
        right:24px;
        top: 14px;
        z-index: 9; 
    }
    .iconlist,
    .dynamicfeed{
        ul{
            list-style: none;
            li{
                width: 20px;
                height: 20px;
                line-height: 20px;
                display: inline-block;
                margin-left:20px;
                i{
                    display: inline-block; 
                    cursor: pointer;
                }
                &.line{
                    span{
                        display: inline-block;
                        height: 16px;
                        width: 2px;
                        background-color: #D8D8D8;
                        vertical-align: middle;
                    }
                }
                &.run {
                    i {
                        color: #66C681;
                    }
                }
                &.status{
                    width: 52px;
                    vertical-align: middle;
                    span{
                        display: block;
                    }
                }
                &.pasuse,&.resume {
                    i {
                        color: #F9B640;
                    }
                }
                &.rungrey{ 
                    i{
                        color: rgba(72, 72, 102, 0.3200);
                        cursor: not-allowed;
                    }
                }
                &.stop {
                    i {
                        color: #F87D7D;
                    }
                }
                &.target{
                    margin-left: 4px;
                    i{
                        color: #4C7AE3;
                    }
                }
                &.attack{
                    i{
                        color: #F9B640;
                    }
                }
                &.loop{
                    i{
                        color: #F87D7D;
                    }
                }
            }
           
        }
    }
    /deep/ .el-dialog__header{
        padding: 15px 0;
        background: #4C7AE3; 
        .el-dialog__title{
            display: inline-block;
            font-size: 14px;
            // width: 104px;
            padding: 0 24px;
            text-align: center;
            border-left: 2px solid #fff;
            height: 32px;
            line-height: 32px;
            background:rgba(255,255,255,0.12);
            color: #fff; 
            width: auto;
        }
    }
    /deep/ .el-tabs__item{
        height: 48px;
        line-height: 48px;
        padding: 0 24px;
    }
    /deep/ .el-tabs__item.is-active{
        color: #4C7AE3;
        font-weight: 500;
    }
    /deep/ .el-tabs__nav-wrap{
        padding: 0 24px; 
    }
    /deep/ .el-tabs__nav-wrap::after{
        background: #E8E8F5;
        height: 0px;
    }
    /deep/ .el-tabs__header{
        margin: 0 0 15px;
        background: #fff;
        border-radius: 4px;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
        border: none;
        padding: 0 24px;
    } 
</style>
<script>
import { task } from '@/api/task.js'
import { vulscan } from '@/api/vulscan.js'

import detailoverview from './TaskDetail_overview.vue';

import targetlist from './VulScanTaskDetail_targetlist.vue';
import vulnmsg from './VulScanTaskDetail_vullist.vue';
export default {
    name:'',
    components:{
        detailoverview,
        targetlist,
        vulnmsg,
    },
    data(){
        return{
            task_id: this.$route.query.id,  
            task_name:this.$route.query.name,
            activeName:'tabs1', 
            status:'',
            statusName:'',
            subactiveName:'subtabs1',
        }
    },
    created(){

    },
    mounted(){
        var _tab =  !(localStorage.getItem('taskTab')) ? 'tabs1' :localStorage.getItem('taskTab') ;
        this.activeName = _tab;

        if(this.activeName == 'tabs2'){
            this.$refs.targetlist.getData(); 
        }

        if(this.activeName == 'tabs4'){
                this.$refs.vulnmsg.getSelectlist();
                this.$refs.vulnmsg.getData();
            }
         
        this.getTaskStatus();

    },
    methods:{
        //监听概述页面的目标风险点击事件跳转到测试目标
        act(da){
            this.activeName = da.tab;
            //改子组件的level
            this.$refs.targetlist.formData.riskLevel = da.level
            this.$refs.targetlist.getData();
        },
        //监听概述页面的漏洞风险点击事件跳转到测试目标
        actLD(da){ 
            this.activeName = da.tab;
            //改子组件的level
            this.$refs.vulnmsg.formDatabug.risk_level = da.level
            this.$refs.vulnmsg.getData();
        },
        handleClick(tabs){
            localStorage.setItem('taskTab', tabs.name);
           
            if(tabs.name == 'tabs2'){
                this.$refs.targetlist.getData();
            }
      
            if(tabs.name == 'tabs4'){
                this.$refs.vulnmsg.getSelectlist();
                this.$refs.vulnmsg.getData();
            }

        },
        getInitData(){

        },
        async getTaskStatus(){ //任务状态
            const res = await vulscan.taskState({
                id:this.task_id
            })
            if(res.code == 200){
                this.status = res.data.status;
                this.statusName = res.data.statusName;
            }
        }, 
        async TaskStop(){
            const dt = await task.taskchangestate({
                taskId:  this.task_id,
                operate:'stop'
              });
              if (dt.code == 200) {
                this.$message({
                  message: "结束任务成功",
                  type: "success"
                }); 
                this.getTaskStatus();
              } else {
                this.$message({
                  message: dt.msg,
                  type: "error"
                });
              }
        },
         
    },
    watch: {
        activeName(newValue, oldValue) {
            if(newValue=='tabs9'){
              
            }else if(newValue=='tabs2'){
              
            }
        }
    },
}
</script>