<template>
    <div>
        <div class="main-title">
            <router-link :underline="false" class="classA" :to="{ path: '/logicvuln' }" >逻辑漏洞
            </router-link>  
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark"  :content="task_name"  placement="bottom">
                    <span>{{task_name}}</span>
                </el-tooltip>
            </label>
        </div> 
        <div class="taskinfolist  ">
            <el-tabs v-model="activeName" @tab-click="handleClick" > 
                <el-tab-pane label="测试目标" name="tabs1">
                    <targetlist  ref="targetlist" :task_id=task_id :task_name=task_name></targetlist> 
                </el-tab-pane>
                <el-tab-pane label="漏洞管理" name="tabs2"> 
                    <vulnmsg  ref="vulnmsg" :task_id=task_id></vulnmsg>
                </el-tab-pane> 
                <el-tab-pane label="被动流量" name="tabs3"> 
                    <passiveflux ref="passiveflux" :task_id=task_id >
                    </passiveflux> 
                </el-tab-pane>
                <el-tab-pane label="日志" name="tabs4" > 
                    <loglist ref="loglist" :task_id=task_id ></loglist> 
                </el-tab-pane>
            </el-tabs>
        </div>
    </div>
</template>
<style lang="less" scoped>
    .taskinfolist{  
        min-height: calc(100% - 39px); 
        box-sizing: border-box;  
        position: relative; 
        /deep/ .el-tabs{
            height: 100%; 
            .el-tabs__content{
                height: calc(100% - 64px);
                overflow: inherit;
                >div{
                    height: 100%; 
                }
            }
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
import targetlist from './logicVulnDetail_target.vue'
import vulnmsg from './logicVulnDetail_vuln.vue'
import loglist from './logicVulnDetail_log.vue'
import passiveflux from "./components/passiveflux.vue";
export default {
    name:'',
    props:{

    },
    components:{
        targetlist,
        vulnmsg,
        loglist,
        passiveflux
    },
    data(){
        return{
            task_id: this.$route.query.id,  
            task_name:this.$route.query.name,
            activeName:'',
        } 
    },
    created(){
        this.$store.state.activefirstMenu = '/logicvuln';
    },
    mounted(){
        var _tab =  !(localStorage.getItem('taskTab')) ? 'tabs1' :localStorage.getItem('taskTab') ;
        this.activeName = _tab;
        if(this.activeName == 'tabs1'){
            this.$refs.targetlist.getData();
        }
        if(this.activeName == 'tabs2'){
            this.$refs.vulnmsg.getData();
        }
        if(this.activeName == 'tabs3'){  
            this.$refs.passiveflux.getData();  
        }
        if(this.activeName == 'tabs4'){  
            this.$refs.loglist.getData();
        }
    },
    methods:{
        handleClick(tabs){
            if(tabs.name == 'tabs1'){
                this.$refs.targetlist.getData();
                console.log(1)
            }
            if(tabs.name == 'tabs2'){
                this.$refs.vulnmsg.getData();
                console.log(2)
            }
            if(tabs.name == 'tabs3'){  
                console.log(3)
                this.$refs.passiveflux.getData();  
            }
            if(tabs.name == 'tabs4'){  
                this.$refs.loglist.getData();
            }
        },
    }
}
</script>