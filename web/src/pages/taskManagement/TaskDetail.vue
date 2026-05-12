/*  
    渗透任务详情页
 */
<template>
    <div>
        <div class="main-title">
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
      
            <router-link :underline="false" class="classA" :to="{ path: '/task' }" >渗透任务
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
                    <li class="pasuse" @click="TaskPause" title="任务暂停" v-if="status == 3">
                        <i class="iconfont iconzanting"></i>
                    </li>
                    <li class="resume" @click="TaskResume" title="任务开始" v-else-if="status == 5">
                        <i class="iconfont iconkaishi"></i>
                    </li>
                    <li class="rungrey" v-else>
                        <i class="iconfont iconkaishi"></i>
                    </li>
                    <li class="stop" @click="TaskStop" title="任务结束" v-if="status != 4">
                        <i class="iconfont icontingzhi"></i>
                    </li>
                    <li class="rungrey" v-else>
                        <i class="iconfont icontingzhi"></i>
                    </li>  
                </ul>
            </div>
            <div class="dynamicfeed">
                <ul>
                   
                    <li class="target" @click="AddTarget" title="新增目标">
                        <i class="iconfont iconxinzengmubiao"></i>
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
                <el-tab-pane label="信息收集" name="tabs3"> 
                    <informationGathering ref="information" :subactiveName="'subtabs1'" v-if="activeName === 'tabs3'"></informationGathering>
                    
                </el-tab-pane> 
                <el-tab-pane label="漏洞测试" name="tabs4"> 
                    <vulnmsg  ref="vulnmsg" :task_id=task_id></vulnmsg>
                </el-tab-pane> 
                <el-tab-pane label="漏洞取证" name="tabs6"> 
                    <obtain  ref="obtain" :task_id=task_id></obtain>
                </el-tab-pane>  
                <el-tab-pane label="远程会话" name="tabs8"> 
                    <remoteSession  ref="remoteSession" :task_id=task_id></remoteSession>
                </el-tab-pane>  
                <el-tab-pane label="待测漏洞" name="tabs7"> 
                    <tobetested  ref="tobetested" :task_id=task_id></tobetested>
                </el-tab-pane> 
                <el-tab-pane label="攻击拓扑" name="tabs9"> 
                    <Attacktopology v-if="activeName=='tabs9'" ref="Attacktopology" ></Attacktopology>
                </el-tab-pane> 
                <el-tab-pane label="测试日志" name="tabs5" > 
                    <loglist ref="loglist" :task_id=task_id ></loglist> 
                </el-tab-pane>
            </el-tabs>
        </div>



    </div> 
</template>
<style lang="less" scoped>
    /deep/ .taskinfolist{ 
        height: calc(100% - 39px);
        // min-height: calc(100% - 39px);
        // overflow-y: auto;
        box-sizing: border-box;  
        position: relative; 
        >.el-tabs{
            height: 100%; 
            >.el-tabs__content{
                height: calc(100% - 56px);
                // overflow: inherit;
                overflow-y: auto;
                overflow-x: hidden;
                >div{
                    height: 100%;
                    // box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
                    // border-radius: 4px;
                    // margin-bottom: 10px;
                    >div{
                        height: 100%;
                    }
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
import detailoverview from './TaskDetail_overview.vue';
import targetlist from './TaskDetail_targetlist.vue';
import vulnmsg from './vulnmsg.vue';
import obtain from './obtain.vue';
import  Attacktopology  from "./components/Attacktopology.vue";
import remoteSession from './remoteSession.vue';
import loglist from './TaskDetail_log.vue'
import informationGathering from "../sceneManagement/informationGathering.vue";
import tobetested from './TaskDetail_Tobetested.vue'
export default {
    name:'',
    components:{
        detailoverview,
        targetlist,
        vulnmsg,
        loglist,
        obtain,
        remoteSession,
        informationGathering,
        tobetested,
        Attacktopology
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
        if(this.activeName == 'tabs1'){
            this.$refs.overview.getData();
        }
        if(this.activeName == 'tabs2'){
            this.$refs.targetlist.getData(); 
        }
        if(this.activeName == 'tabs3'){  //信息收集 
        }
        if(this.activeName == 'tabs4'){
                this.$refs.vulnmsg.getSelectlist();
                this.$refs.vulnmsg.getData();
            }
        if( this.activeName  == 'tabs5'){
            this.$refs.loglist.getData();
        }
        if( this.activeName  == 'tabs7'){
            this.$refs.tobetested.getSelectlist();
            this.$refs.tobetested.getData();
        }
        if( this.activeName  == 'tabs9'){
            
            // this.$refs.Attacktopology.getData()

        }


        // this.status = localStorage.getItem('status');
        // this.statusName = localStorage.getItem('statusName');
        
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
            if(tabs.name == 'tabs1'){
                this.$refs.overview.getData();
            }
            if(tabs.name == 'tabs2'){
                this.$refs.targetlist.getData();
            }
            if(tabs.name == 'tabs3'){  //信息收集 
            }
            if(tabs.name == 'tabs4'){
                this.$refs.vulnmsg.getSelectlist();
                this.$refs.vulnmsg.getData();
            }
            if(tabs.name == 'tabs5'){
                this.$refs.loglist.getData();
            }
            if(tabs.name == 'tabs6'){
                this.$refs.obtain.getTableData();
            }
            if(tabs.name == 'tabs8'){
                this.$refs.remoteSession.getTableData();
            }
            if(tabs.name == 'tabs7'){
                this.$refs.tobetested.getSelectlist();
                this.$refs.tobetested.getData();
            }
        },
        getInitData(){

        },
        async getTaskStatus(){ //任务状态
            const res = await task.getTaskStatus({
                taskId:this.task_id
            })
            if(res.code == 200){
                this.status = res.data.status;
                this.statusName = res.data.statusName;
            }
        }, 
        async TaskPause(){ //暂停
            const dt = await task.taskchangestate({
                taskId: this.task_id,
                operate:'pause'
              });
              if (dt.code == 200) {
                this.$message({
                  message: "任务暂停",
                  type: "success"
                }); 
                this.getTaskStatus();
              } else {
                this.$message({
                  message: dt.msg,
                  type: "error"
                });
              }
              this.getTaskStatus() //重新获取任务状态
        },
        async TaskResume(){ //开始
            const dt = await task.taskchangestate({
                taskId: this.task_id,
                operate:'resume'
              });
              if (dt.code == 200) {
                this.$message({
                  message: "任务重新开始",
                  type: "success"
                }); 
                this.getTaskStatus();
              } else {
                this.$message({
                  message: dt.msg,
                  type: "error"
                });
              }

              this.getTaskStatus()//重新获取任务状态
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
        AddTarget(){ //新增目标
            this.$router.push({
                path: `/createtask`,
                query: { 
                    task_id: this.task_id,
                    flag: 6, 
                    taskName:this.task_name,
                    disabled:true,
                }
            });
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