<template>
    <!-- 任务组详情 -->
    <div>
        <div class="main-title">
            <router-link :underline="false" class="classA" :to="{ path: '/taskgroup' }" >任务组
            </router-link>  
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark"  :content="group_name"  placement="bottom">
                    <span>{{group_name}}</span>
                </el-tooltip>
            </label>
        </div> 
        <div class="groupinfolist">
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
                    <groupoverview :group_id=group_id></groupoverview>
                </el-tab-pane>
                <el-tab-pane label="渗透任务" name="tabs2"> 
                    <groupTask :group_id=group_id :group_name=group_name></groupTask>
                </el-tab-pane>
            </el-tabs>
        </div>

    </div>
</template>
<style scoped lang="less">
.groupinfolist{
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
                }
            }
        } 
}
.iconlist{
    position: absolute;
    right: 20px;
    top: 14px;
    z-index: 9; 
}
.iconlist {
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
import { task_group } from '@/api/task.js'
import groupoverview from './TaskGroupDetail_overview.vue';
import groupTask from './TaskList.vue'
export default {
    name:'',
    components:{
        groupoverview,
        groupTask
    },
    data(){
        return{
            group_id: this.$route.query.id,  
            group_name:this.$route.query.name,
            activeName:'tabs1',
            status:'',
            statusName:'',
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/taskgroup';
        this.getTaskStatus();

        var _tab =  !(localStorage.getItem('groupTab')) ? 'tabs1' :localStorage.getItem('groupTab') ;
        this.activeName = _tab;
    },
    mounted(){
       
    },
    methods:{
        handleClick(tabs){
            localStorage.setItem('groupTab', tabs.name); 
        },
        async getTaskStatus(){ //任务状态
            const res = await task_group.getTaskGroupStatus({
                group_id:this.group_id
            })
            if(res.code == 200){
                this.status = res.data.statusNumber;
                this.statusName = res.data.status;
            }
        }, 
    },
}
</script>