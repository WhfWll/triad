<template>
    <!-- 流量详情 -->
    <div>
        <div class="main-title  ">
            <router-link :underline="false" class="classA" :to="{ path: '/traffic' }">被动流量</router-link>
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark" :content="task_name" placement="bottom">
                    <span> {{ task_name }}</span>
                </el-tooltip> 
            </label>
        </div>
        <div class="taskinfolist context_box_bg">
            <div class="iconlist">
                <ul>
                    <li style=" margin-right: 20px;" >
                        <span class="tag_status tag_primary" style="cursor: pointer" 
                            v-if="activeName=='tabs1' && !isUpdate && taskStatus_num == 3"  @click="handleupdate()" > 
                            编辑
                        </span>
                        <span class="tag_status tag_primary" style="cursor: pointer" 
                        v-if="isUpdate"  @click="handleupdatesage()"  > 
                            保存
                        </span>
                    </li>
                    <li class="status">
                        <!--   待触发：1 ；待执行：2；运行中：3；已完成：4 ；暂停中：5。 -->
                        <span class="tag_status tag_danger1" v-if="taskStatus_num == 1">{{ taskStatus }}</span>
                        <span class="tag_status tag_warning" v-if="taskStatus_num == 2">{{ taskStatus }}</span>
                        <span class="tag_status tag_primary" v-if="taskStatus_num == 3">{{ taskStatus }}</span>
                        <span class="tag_status tag_success" v-if="taskStatus_num == 4">{{ taskStatus }}</span>
                        <span class="tag_status tag_danger" v-if="taskStatus_num == 5">{{ taskStatus }}</span>
                    </li>
                    <!-- <li class="pasuse" @click="TaskPause" title="任务暂停" v-if="taskStatus_num == 3">
                        <i class="iconfont iconzanting"></i>
                    </li>
                    <li class="resume" @click="TaskResume" title="任务开始" v-else-if="taskStatus_num == 5">
                        <i class="iconfont iconkaishi"></i>
                    </li>
                    <li class="rungrey" v-else>
                        <i class="iconfont iconkaishi"></i>
                    </li> -->
                    <li class="stop" @click="TaskStop" title="任务结束" v-if="taskStatus_num != 4">
                        <i class="iconfont icontingzhi"></i>
                    </li>
                    <li class="rungrey" v-else>
                        <i class="iconfont icontingzhi"></i>
                    </li> 
                    <li @click="download">
                        <i class="iconfont iconxiazai"></i>
                    </li>
                </ul>
            </div>
            <el-tabs v-model="activeName" @tab-click="handleClick">
                <el-tab-pane label="概述" name="tabs1">
                    <trafficinformation ref="trafficinformation" :task_id=task_id >
                    </trafficinformation> 
                </el-tab-pane> 
                <el-tab-pane label="漏洞信息" name="tabs2">
                    <trafficloopinfo ref="trafficloopinfo" :trafficid=task_id >
                    </trafficloopinfo> 
                </el-tab-pane>
                <!-- <el-tab-pane label="被动流量" name="tabs3"> 
                    <passiveflux ref="passiveflux" :task_id=task_id >
                    </passiveflux> 
                </el-tab-pane> -->
                <el-tab-pane label="测试日志" name="tabs4">
                    <testlog ref="testlog" :task_id=task_id></testlog>
                </el-tab-pane>
            </el-tabs>
        </div>
    </div>
</template>
<style lang="less" scoped>
    .context_box_bg {
        background: none;
    }
    /deep/ .el-tabs__item {
        height: 48px;
        line-height: 48px;
        padding: 0 24px;
    }
    /deep/ .el-tabs__item.is-active {
        color: #4C7AE3;
        font-weight: 500;
    }

    /deep/ .el-tabs__nav-wrap {
        padding: 0 24px;
    }

    /deep/ .el-tabs__nav-wrap::after {
        background: #E8E8F5;
        height: 1px;
    }

    /deep/ .el-tabs__header {
        margin: 0 0 15px;
        background: #fff;
        border-radius: 4px;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
        border: none;
    }
    .taskinfolist {
       // height: calc(100% - 39px);
       min-height: calc(100% - 39px);
       overflow-y: auto;
       box-sizing: border-box;
       position: relative;

       /deep/ .el-tabs {
           height: 100%;

           .el-tabs__content {
               height: calc(100% - 64px);

               >div {
                   height: 100%;
               }
           }
       }

   }  
   .iconlist {
        position: absolute;
        right: 24px;
        top: 10px;
        z-index: 9;
        ul {
            list-style: none; 
            li {
                width: 20px;
                height: 20px;
                line-height: 20px;
                display: inline-block;
                margin-left: 20px;

                i {
                    display: inline-block;
                    cursor: pointer;
                }

                &.line {
                    span {
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

                &.status {
                    width: 52px;
                    vertical-align: middle;

                    span {
                        display: block;
                    }
                }

                &.pasuse,
                &.resume {
                    i {
                        color: #F9B640;
                    }
                }

                &.rungrey {
                    i {
                        color: rgba(72, 72, 102, 0.3200);
                        cursor: not-allowed;
                    }
                }

                &.stop {
                    i {
                        color: #F87D7D;
                    }
                }

                &.target {
                    margin-left: 4px;

                    i {
                        color: #4C7AE3;
                    }
                }

                &.attack {
                    i {
                        color: #F9B640;
                    }
                }

                &.loop {
                    i {
                        color: #F87D7D;
                    }
                }
            } 
        }
    }
 /deep/ .el-dialog__header {
     padding: 15px 0;
     background: #4C7AE3;

     .el-dialog__title {
         display: inline-block;
         font-size: 14px;
         // width: 104px;
         padding: 0 24px;
         text-align: center;
         border-left: 2px solid #fff;
         height: 32px;
         line-height: 32px;
         background: rgba(255, 255, 255, 0.12);
         color: #fff;
         width: auto;
     }
 }
</style>
<script>
import trafficinformation from "./components/trafficinformation.vue";
import trafficloopinfo from "./components/trafficloopinfo.vue";
import passiveflux from "./components/passiveflux.vue";
import testlog from "./components/testlog.vue";
import { traffic } from '@/api/traffic.js'
import jsFileDownload from 'js-file-download'
export default {
    name:'trafficinfo',
    components: {
        trafficinformation,
        trafficloopinfo,
        passiveflux,
        testlog
    },
    data(){
        return{
            taskStatus_num:0,
            taskStatus:'',
            task_id: this.$route.query.id,
            task_name: this.$route.query.name,
            activeName:'tabs1',
            userID:'',
            isUpdate:false,
        }
    },
    created() {
        this.$store.state.activefirstMenu = "/traffic"; 
    },
    mounted() {
        var _tab = !(localStorage.getItem('trafficTab')) ? 'tabs1' : localStorage.getItem('trafficTab');
        this.activeName = _tab;
        this.userID = this.$commonjs.decryptCBC(localStorage.getItem('user_id'),this.$commonjs.myKey); 
        this.getData(this.activeName);
        this.getStatus();
    },
    methods: {
        async getStatus(){
            const res = await traffic.getStatus({
                flowTaskId:this.task_id,
            });
            if(res.code == 200){
                this.taskStatus_num = res.data.status;
                this.taskStatus = res.data.statusName;
            }

        },
        async TaskPause(){ //暂停
            const res = await traffic.trafficStatus({
                flowTaskId:this.task_id,
                operate:'pause',
                userId:this.userID
            })
            if(res.code == 200){
                this.$message({
                    message: '任务暂停',
                    type: 'success'
                });
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async TaskResume(){ //开始
            const res = await traffic.trafficStatus({
                flowTaskId:this.task_id,
                operate:'resume',
                userId:this.userID
            })
            if(res.code == 200){
                this.$message({
                    message: '任务重新开始',
                    type: 'success'
                });
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async TaskStop(){ //结束
            const res = await traffic.trafficStatus({
                flowTaskId:this.task_id,
                operate:'stop',
                userId:this.userID
            })
            if(res.code == 200){
                this.$message({
                    message: '任务结束',
                    type: 'success'
                });
                this.getStatus();
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        handleClick(tab){
            localStorage.setItem('trafficTab', tab.name);
            this.getData(tab.name); 
        },
        async getData(tabname) {  
            if (tabname == 'tabs1') { //综合信息 
                this.$refs.trafficinformation.getData();  
            }
            else if (tabname == 'tabs2'){
                this.$refs.trafficloopinfo.getEnum(); 
                this.$refs.trafficloopinfo.getData(); 
                 
            }
            else if (tabname == 'tabs3'){
                 this.$refs.passiveflux.getData();  
            } else if (tabname == 'tabs4'){
                this.$refs.testlog.getData();
            }
        },
        async download(){
            const res = await traffic.downloadhttps({
                flowTaskId: this.task_id, 
            });
            if(res.code == 200){
                jsFileDownload(res.data.cert,'系统https证书.crt');
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
           

        },
        handleupdate(){
            this.isUpdate = true;
            this.$refs.trafficinformation.handleupdate(this.isUpdate)
        },
        async handleupdatesage(){ //保存编辑
            let param =  this.$refs.trafficinformation.getAllData();
            console.log(param);
            let config={
                waitCred:{
                    pattern:Number(param.waitCred_pattern),
                    value:param.waitCred_value
                },
                fuzzParam:{
                    character:param.fuzzParam.character.join(','),
                    number:param.fuzzParam.number.join(','),
                }, 
                fuzzDict:{
                    character:param.fuzzDict.character,
                    number:param.fuzzDict.number,
                },
                response:{
                    jsonKeyword:param.response.jsonKeyword,
                    noJsonSwitch:param.response.noJsonSwitch,
                    noJsonKeyword:param.response.noJsonKeyword
                },
            }
            let newparam = {
                flowTaskId:this.task_id, 
                taskName:param.taskname,
                expireTime:param.duration,
                networkCard:param.networkcard,
                port:param.networkport,
                targetUrl:param.target, 
                vulConfig:param.vulname.join(','),
                userId:Number(param.userId),
                otherConfig:JSON.stringify(config),
            }
            const res = await traffic.flowtaskedit(newparam);
            if(res.code == 200){
                this.isUpdate = false;
                this.$refs.trafficinformation.getData();  
                this.$refs.trafficinformation.handleupdate(this.isUpdate)
                this.$message({
                    message: '编辑任务成功',
                    type: 'success'
                });
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
    },
}
</script>
 