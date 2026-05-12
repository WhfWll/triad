<template>
    <div class="taskloginfo">
        <div class="toplogo"> 
            <span style="padding-left:12px;font-weight: 600;">现场漏洞验证工具集
            </span>
        </div>
        <div class="infos">
            <div>
                <span>{{ target }}</span>
            </div>
            <div>
                <ul class="loglist">
                    <li v-for="(item,index) in loglist" :key=index>
                        <span>[*]</span>
                        <span>[{{ item }}]</span>
                        
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>
<script> 
import bas from '@/api/bas.js'
export default ({
    name: 'bastargetlog',
    data() {
        return {
            id: this.$route.query.id,
            target: this.$route.query.target,
            loglist:[],
        }
    },
    created() {
        this.getloginfo();
    },
    methods: {
        async getloginfo(time) {
            const res = await bas.basTargetlog({
                id: this.id
            });
            if(res.code == 200){
                this.loglist = res.data.content;
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                }); 
            }
            
             
        },
    },
})

</script>
<style scoped lang="less">



.taskloginfo {
    width: 100%;
    height: 100%;
    background-color: #FFF;
    display: flex;
    flex-direction: column;

    .toplogo {
        width: 100%;
        height: 48px;
        line-height: 48px;
        background: #4C7AE3;
        box-shadow: 0px 2px 6px 1px rgba(0, 21, 41, 0.11999999731779099);
        border-radius: 0px 0px 0px 0px;

        span {
            color: #FFF;
            padding-left: 12px;
            font-weight: 600;
        }
    }

    .infos {
        padding: 24px;
        box-sizing: border-box; 
        color: rgba(72, 72, 102, 0.6400);
    }
    .loglist {
        list-style: none;
        margin-top: 24px;
        li {
            color: rgba(72, 72, 102, 0.6400);
            font-size: 13px;
        }
    }
}
</style>

