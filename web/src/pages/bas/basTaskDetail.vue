<template>
    <div> 
        <div class="main-title">
            <router-link :underline="false" class="classA" :to="{ path: '/bastask' }" >安全配置核查
            </router-link>  
            <label class="currentpagetitle">
                <el-tooltip class="item" effect="dark" :content="name+'详情-测试目标'"   placement="bottom">
                    <span>{{name}}详情-测试目标</span>
                </el-tooltip>
            </label>
        </div>  
        <div>
            <el-tabs v-model="activeName" @tab-click="handleClick" > 
                <el-tab-pane label="测试目标" name="tabs1">
                    <bastargetlist  ref="bastargetlist" :task_id=id :task_name=name></bastargetlist> 
                </el-tab-pane>
                <el-tab-pane label="漏洞测试" name="tabs2">
                    <basvultest  ref="overview" :task_name="name" :task_id=id ></basvultest> 
                </el-tab-pane>
            </el-tabs>
        </div>
    </div>
</template>
<style lang="less" scoped>

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
import bastargetlist from './basTaskDetail_target.vue';
import basvultest from './basTaskDetail_vulTest.vue'

export default {
    name:'bastaskDetail',
    components:{
        bastargetlist,
        basvultest
    },
    data(){
        return{
            id: this.$route.query.id,
            name:this.$route.query.name, 
            activeName:'tabs1'
        }
    },
    created(){
        this.$store.state.activefirstMenu="/bastask"; 
    },
    mounted(){
        this.$refs.bastargetlist.getData(); 
    },
    methods:{
        handleClick(){

        },
    }
}
</script>